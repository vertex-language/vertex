package hir

import (
	"sort"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// info.go is hir's query layer over types: every question the rest of the
// package asks about a checked object rather than about a type's shape.
// bridge.go classifies types; this file classifies *objects* — what a
// parameter's convention is, where a global landed, which unit declares a
// function, what a marker says.
//
// The split matters for the same reason bridge.go exists: a change to how
// types spells parameter modes or function markers lands in one of these
// two files and nowhere else.

// ---------------------------------------------------------------------------
// markers and conventions
// ---------------------------------------------------------------------------

// markerOf renders a Func's A.6.1 FunctionMarker. The marker lives on the
// Signature because A.4.2 makes it part of the callee's contract, so a
// synthesized Func (which has no Origin) carries none by construction.
func (l *lowerer) markerOf(f *Func) string {
	if f == nil {
		return ""
	}
	fn, ok := f.Origin.(*types.Func)
	if !ok {
		return ""
	}
	sig := fn.Signature()
	if sig == nil {
		return ""
	}
	return sig.Marker().String()
}

// suspends reports whether fn's own signature carries A.6.1's suspend
// marker. It answers the same question as markerOf, but against a checked
// *types.Func directly rather than a lowered *hir.Func with an Origin —
// which is what buildEntry has in hand for main/the test function before
// any *hir.Func shell exists for it.
func (l *lowerer) suspends(fn *types.Func) bool {
	if fn == nil {
		return false
	}
	sig := fn.Signature()
	if sig == nil {
		return false
	}
	return sig.Marker().String() == "async"
}

// isMutParam reports A.3.2's `mut` convention: exclusive, non-owning,
// mutating. It lowers to a pointer parameter the callee writes through.
func (l *lowerer) isMutParam(obj types.Object) bool {
	v, ok := obj.(*types.Var)
	return ok && v.Mode() == types.ModeMut
}

// isOwningParam reports A.3.2's `var` convention. It answers only that the
// position is owning — whether this call passes the caller's original or a
// deep copy is decided at the call site by the marker (A.4.6), which is
// owningExpr's question, not this one's.
func (l *lowerer) isOwningParam(obj types.Object) bool {
	v, ok := obj.(*types.Var)
	return ok && v.Owning()
}

// consumesReceiver reports whether a method's receiver is `var`-typed. A
// receiver position has no argument slot to carry a marker, so A.6.1's
// single exception applies: the consumption is unconditional and there is
// no bare form that copies.
func (l *lowerer) consumesReceiver(fn *types.Func) bool {
	sig := fn.Signature()
	return sig != nil && sig.Recv() != nil && sig.Recv().Owning()
}

// needsAddressableReceiver reports whether the receiver argument must be a
// pointer. Two reasons, and they are unrelated: a `mut` receiver is
// literally the pointer parameter (A.3.2), and an aggregate receiver
// travels byval, which at the vir level is also a ptr — see decl.go's param.
func (l *lowerer) needsAddressableReceiver(fn *types.Func) bool {
	sig := fn.Signature()
	if sig == nil || sig.Recv() == nil {
		return false
	}
	if sig.Recv().Mode() == types.ModeMut {
		return true
	}
	return IsAggregate(l.hirType(sig.Recv().Type()))
}

// ---------------------------------------------------------------------------
// globals
// ---------------------------------------------------------------------------

// globalBinding is how a top-level object is reached from a function body.
// A.5.1's let/var split survives to module scope: a scalar `let` folded to a
// vir const yields a direct value and occupies no storage (§6.2), while a
// `var` is real module-level storage reached through its address.
type globalBinding struct {
	name    string
	mod     *Module
	typ     Type
	isConst bool
}

func (l *lowerer) bindGlobal(obj types.Object, name string, t Type, isConst bool) {
	if obj == nil {
		return
	}
	if l.globals == nil {
		l.globals = map[types.Object]*globalBinding{}
	}
	l.globals[obj] = &globalBinding{
		name: name, mod: l.currentModule(), typ: t, isConst: isConst,
	}
}

func (l *lowerer) globalFor(obj types.Object) *globalBinding {
	if obj == nil || l.globals == nil {
		return nil
	}
	return l.globals[obj]
}

// ---------------------------------------------------------------------------
// externs
// ---------------------------------------------------------------------------

// bindExtern records that obj (the checked object behind one `declare`
// block member) is reached by calling ef directly — never by scheduling
// it through the monomorphization worklist, which has no notion of a
// foreign body. See lower.go's externs field doc for the failure this
// prevents.
func (l *lowerer) bindExtern(obj types.Object, ef *ExternFunc) {
	if obj == nil {
		return
	}
	if l.externs == nil {
		l.externs = map[types.Object]*ExternFunc{}
	}
	l.externs[obj] = ef
}

func (l *lowerer) externFor(obj types.Object) *ExternFunc {
	if obj == nil || l.externs == nil {
		return nil
	}
	return l.externs[obj]
}

// ---------------------------------------------------------------------------
// enums
// ---------------------------------------------------------------------------

// variantTag resolves a variant name to its discriminant. A.6.5 already
// settled the value at construction — continuing from the previous variant
// when unwritten — so this is a lookup and never a computation.
func (l *lowerer) variantTag(t types.Type, name string) (int64, bool) {
	e, ok := l.underlying(t).(*types.Enum)
	if !ok {
		return 0, false
	}
	v := e.LookupVariant(name)
	if v == nil {
		return 0, false
	}
	return v.Value, true
}

// ---------------------------------------------------------------------------
// test mode
// ---------------------------------------------------------------------------

// testResultType is the type the test wrapper renders through fmt, or nil
// for a bare test function. A.12.2's table picks the verb from this type,
// and a test returning nothing reports through exit status alone — which is
// the entire mechanism behind a test with no Expected.
func (l *lowerer) testResultType(fn *types.Func) types.Type {
	if fn == nil {
		return nil
	}
	sig := fn.Signature()
	if sig == nil {
		return nil
	}
	res := sig.Results()
	if res == nil || res.Len() != 1 {
		return nil
	}
	return res.At(0).Type()
}

// ---------------------------------------------------------------------------
// declaration lookup
// ---------------------------------------------------------------------------

// unitFor finds the loaded unit declaring fn. Monomorphization needs it
// because an instantiation's body is read from the *declaring* package's
// Info, never from the call site's.
func (l *lowerer) unitFor(fn *types.Func) *Unit {
	if fn != nil && fn.Pkg() != nil {
		path := fn.Pkg().Path()
		for _, u := range l.units {
			if u.Path == path {
				return u
			}
			if u.Types != nil && u.Types.Path() == path {
				return u
			}
		}
	}
	if l.cur != nil {
		return l.cur.unit
	}
	return nil
}

// findFuncDecl recovers the syntax for a checked function. The link runs
// through Info.Defs rather than by name: a method and a free function may
// share a spelling, and Defs is keyed by the declaring identifier, which is
// unique.
func findFuncDecl(u *Unit, fn *types.Func) *ast.FuncDecl {
	if u == nil || fn == nil || u.Info == nil {
		return nil
	}
	for _, f := range u.Files {
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Name == nil {
				continue
			}
			if u.Info.Defs[fd.Name] == types.Object(fn) {
				return fd
			}
		}
	}
	return nil
}

// typeParamsOf recovers a generic function's type parameters, in
// declaration order.
//
// types records no parameter list on Func or Signature — a *TypeParam
// carries its own Index (A.7.1's position in the TypeParameterList), so the
// list is reconstructed by collecting every parameter the signature
// mentions and sorting on that index. A parameter mentioned only in the
// body is invisible here, but such a parameter is also uninferable at a
// call site, so it cannot reach a substitution either way.
func typeParamsOf(fn *types.Func) []*types.TypeParam {
	if fn == nil {
		return nil
	}
	sig := fn.Signature()
	if sig == nil {
		return nil
	}
	seen := map[*types.TypeParam]bool{}
	visited := map[types.Type]bool{}
	var out []*types.TypeParam
	if r := sig.Recv(); r != nil {
		collectTypeParams(r.Type(), seen, visited, &out)
	}
	collectTypeParams(sig.Params(), seen, visited, &out)
	collectTypeParams(sig.Results(), seen, visited, &out)
	sort.Slice(out, func(i, j int) bool { return out[i].Index() < out[j].Index() })
	return out
}

// collectTypeParams walks t gathering every *TypeParam it mentions. visited
// guards a self-referential Named — `struct Node { next: typed_ptr Node }`
// is legal and reaches itself.
func collectTypeParams(t types.Type, seen map[*types.TypeParam]bool, visited map[types.Type]bool, out *[]*types.TypeParam) {
	if t == nil || visited[t] {
		return
	}
	visited[t] = true

	switch x := t.(type) {
	case *types.TypeParam:
		if !seen[x] {
			seen[x] = true
			*out = append(*out, x)
		}
	case *types.Ownership:
		collectTypeParams(x.Elem(), seen, visited, out)
	case *types.Array:
		collectTypeParams(x.Elem(), seen, visited, out)
	case *types.Slice:
		collectTypeParams(x.Elem(), seen, visited, out)
	case *types.Chan:
		collectTypeParams(x.Elem(), seen, visited, out)
	case *types.Pointer:
		collectTypeParams(x.Elem(), seen, visited, out)
	case *types.Tensor:
		collectTypeParams(x.Elem(), seen, visited, out)
	case *types.Map:
		collectTypeParams(x.Key(), seen, visited, out)
		collectTypeParams(x.Elem(), seen, visited, out)
	case *types.Tuple:
		for i := 0; i < x.Len(); i++ {
			collectTypeParams(x.At(i).Type(), seen, visited, out)
		}
	case *types.Signature:
		if r := x.Recv(); r != nil {
			collectTypeParams(r.Type(), seen, visited, out)
		}
		collectTypeParams(x.Params(), seen, visited, out)
		collectTypeParams(x.Results(), seen, visited, out)
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			collectTypeParams(x.Field(i).Type, seen, visited, out)
		}
	case *types.Enum:
		for i := 0; i < x.NumVariants(); i++ {
			for _, p := range x.Variant(i).Payload {
				collectTypeParams(p, seen, visited, out)
			}
		}
	case *types.Named:
		for _, a := range x.TypeArgs() {
			collectTypeParams(a, seen, visited, out)
		}
		collectTypeParams(x.Underlying(), seen, visited, out)
	}
}

// varsOf wraps bare types as unnamed tuple elements. A.3.1's TupleElement
// may be named, so a Tuple holds Vars; a tuple built from an expression list
// has nothing to name them with.
func varsOf(ts []types.Type) []*types.Var {
	out := make([]*types.Var, 0, len(ts))
	for _, t := range ts {
		out = append(out, types.NewVar(token.NoPos, nil, "", t))
	}
	return out
}

// ---------------------------------------------------------------------------
// teardown
// ---------------------------------------------------------------------------

// instanceOf resolves a *types.Func to its monomorphized shell, scheduling
// the instantiation if this is the first reference. Only for functions with
// no type arguments — a declared deinit is never generic in its own right.
func (l *lowerer) instanceOf(fn *types.Func) *Func {
	if f := l.work.lookup(fn, nil); f != nil {
		return f
	}
	u := l.unitFor(fn)
	if u == nil {
		return nil
	}
	depth := 0
	if l.cur != nil {
		depth = l.cur.depth + 1
	}
	return l.work.enqueue(u, fn, nil, depth)
}

// emitUserDeinit is the body of a synthesized per-type teardown routine.
//
// A.6.4 fixes the order: the declared deinit runs first, on an object whose
// fields are all still live, and the per-field teardown follows in reverse
// declaration order. A type with no declared deinit still needs the second
// half, which is why this runs for every owning composite and not only for
// the ones that wrote one.
func (l *lowerer) emitUserDeinit(b *funcBuilder, t types.Type, p Value) {
	if n, ok := l.subst(t).(*types.Named); ok {
		if m := n.LookupMethod(token.CtxDeinit); m != nil {
			if target := l.instanceOf(m); target != nil {
				b.call(0, target, p)
			}
		}
	}

	switch l.classify(t) {
	case kUnique:
		// Deep: the pointee's own teardown, then the block itself.
		inner := l.elem(t)
		h := b.load(0, Ptr, p)
		l.deinitAt(b, inner, h)
		b.callBuiltin(0, symUniqueFree, Void, h)

	case kStruct:
		st, ok := l.hirType(t).(StructType)
		if !ok {
			return
		}
		_, fields := l.structParts(t)
		for i := len(fields) - 1; i >= 0; i-- {
			if i >= len(st.Def.Fields) || !l.classify(fields[i].Type).owning() {
				continue
			}
			l.deinitAt(b, fields[i].Type, b.fieldPtr(0, st.Def, p, st.Def.Fields[i].Name))
		}

	case kEnum:
		l.todo(0, "enum deinit routine — switch on the tag, tear down the live variant only")
	}
}

// deinitAt tears down the value stored at addr. An aggregate value *is* its
// address, so only a scalar is loaded first — the same rule the package doc
// comment states once.
func (l *lowerer) deinitAt(b *funcBuilder, t types.Type, addr Value) {
	if t == nil || !l.classify(t).owning() {
		return
	}
	ht := l.hirType(t)
	v := addr
	if !IsAggregate(ht) {
		v = b.load(0, ht, addr)
	}
	l.own.emitDeinit(b, 0, &binding{Value: v, Type: ht, Src: t, Owning: true})
}

// ---------------------------------------------------------------------------
// field defaults
// ---------------------------------------------------------------------------

// emitFieldDefaults handles A.6.2's "a field default is evaluated at
// construction for any field the literal omits."
//
// types.Field records only *that* a default exists, not its expression, so
// there is nothing here to lower yet: the expression lives on the declaring
// ast.RecordDecl and would have to be lowered under the declaring package's
// Info rather than the construction site's. That is a real pass, not a
// lookup, so this reports rather than silently zeroing — compositeLit's
// memset already covers the no-default case honestly, and a defaulted field
// silently reading zero is exactly the bug worth refusing to ship.
func (l *lowerer) emitFieldDefaults(b *funcBuilder, pos token.Pos, t types.Type, st StructType, slot Value, written map[string]bool) {
	s, ok := l.underlying(t).(*types.Struct)
	if !ok {
		return
	}
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		if !f.HasDefault || written[sanitize(f.Name)] {
			continue
		}
		l.todo(pos, "field default for %s.%s — types records that one exists, not its expression",
			types.TypeString(t), f.Name)
	}
}