package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// ------------------------------------------------------------------- scopes

func (fb *funcBuilder) openScope(k scopeKind) {
	fb.top = &scope{parent: fb.top, kind: k}
}

func (fb *funcBuilder) closeScope() {
	fb.top = fb.top.parent
}

func (fb *funcBuilder) bind(obj types.Object, b *binding) {
	fb.bindings()[obj] = b
}

// bindings is a flat object->binding map for the whole function. Scoping was
// already settled by the analyzer — two locals with one spelling are two
// objects — so hir needs no name lookup, only identity.
func (fb *funcBuilder) bindings() map[types.Object]*binding {
	if fb.binds == nil {
		fb.binds = map[types.Object]*binding{}
	}
	return fb.binds
}

func (fb *funcBuilder) lookup(obj types.Object) *binding {
	return fb.binds[obj]
}

// epilogueTo emits teardown for every scope from the innermost out to and
// including `stop`, or to the function scope when stop is nil.
//
// Order within one scope: deferred calls first, in reverse registration
// order, then locals in reverse declaration order. Neither source document
// resolves the interleaving at a shared exit edge; this is the pin.
//
// Everything is emitted before *every* terminator leaving the scope —
// fall-through, return, break, continue. With no unwinder the edge set is
// finite and static, which is what makes duplication rather than a runtime
// list the right shape.
func (fb *funcBuilder) epilogueTo(stop *scope) {
	for s := fb.top; s != nil; s = s.parent {
		for i := len(s.defers) - 1; i >= 0; i-- {
			d := s.defers[i]
			fb.seq.add(&Instr{
				Op: OpCall, Type: Void,
				Callee: d.callee, Module: d.module, Args: d.args,
			})
		}
		for i := len(s.locals) - 1; i >= 0; i-- {
			b := s.locals[i]
			if b.dead {
				// Transferred away: already dead, not destroyed again.
				continue
			}
			fb.dropBinding(b)
		}
		if s == stop || (stop == nil && s.kind == scopeFunc) {
			return
		}
	}
}

// ------------------------------------------------------------------ statements

func (fb *funcBuilder) block(b *ast.BlockStmt) {
	if b == nil {
		return
	}
	fb.openScope(scopeBlock)
	for _, s := range b.List {
		fb.stmt(s)
		if fb.terminated() {
			break
		}
	}
	if !fb.terminated() {
		fb.epilogueTo(fb.top)
	}
	fb.closeScope()
}

func (fb *funcBuilder) stmt(s ast.Stmt) {
	switch x := s.(type) {
	case nil, *ast.BadStmt:
		return

	case *ast.BlockStmt:
		fb.block(x)

	case *ast.DeclStmt:
		if vd, ok := x.Decl.(*ast.VarDecl); ok {
			fb.varDecl(vd)
		}

	case *ast.ExprStmt:
		fb.expr(x.X)

	case *ast.AssignStmt:
		fb.assign(x)

	case *ast.IfStmt:
		fb.ifStmt(x)

	case *ast.WhileStmt:
		fb.whileStmt(x)

	case *ast.ForStmt:
		fb.forStmt(x)

	case *ast.SwitchStmt:
		fb.switchStmt(x)

	case *ast.SelectStmt:
		// select is entirely a runtime construct: a descriptor array plus
		// one of two entry points, chosen statically by whether the
		// statement is awaited. Both need async.
		fb.l.todoAt(x.Pos(), "select")

	case *ast.ReturnStmt:
		fb.returnStmt(x)

	case *ast.DeferStmt:
		fb.deferStmt(x)

	case *ast.BranchStmt:
		switch x.Tok {
		case token.BREAK:
			fb.epilogueTo(fb.enclosingLoop())
			fb.seq.add(&Break{})
		case token.CONTINUE:
			fb.epilogueTo(fb.enclosingLoop())
			fb.seq.add(&Continue{})
		case token.FALLTHROUGH:
			// vir's switch cases do not fall through, so this needs the
			// following clause's body duplicated into this one.
			fb.l.todoAt(x.Pos(), "fallthrough")
		}
	}
}

func (fb *funcBuilder) enclosingLoop() *scope {
	for s := fb.top; s != nil; s = s.parent {
		if s.kind == scopeLoop {
			return s
		}
	}
	return fb.top
}

// varDecl introduces a local. `let` is a value, `var` is a slot — rule 2 of
// the package doc — unless something forces a slot anyway.
func (fb *funcBuilder) varDecl(d *ast.VarDecl) {
	// Initializers are lowered first: a binding is not in scope inside its
	// own initializer.
	vals := make([]Value, 0, len(d.Values))
	for _, v := range d.Values {
		vals = append(vals, fb.owning(v))
	}

	for i, bnd := range d.Bindings {
		obj := fb.l.info.Defs[bnd.Name]
		if obj == nil || bnd.Name.IsBlank() {
			continue // evaluated and discarded
		}
		t := fb.typ(obj.Type())

		var init Value
		if i < len(vals) {
			init = vals[i]
		}

		b := &binding{obj: obj, typ: t, vt: obj.Type()}
		switch {
		case fb.forcesSlot(obj, t):
			b.slot = true
			b.val = fb.alloca(t, "s")
			if init.IsZero() {
				// A var with a type and no initializer is that type's zero
				// value. Every type has one, so there is no
				// definite-assignment analysis anywhere in the language.
				fb.zero(b.val, t)
			} else {
				fb.storeInto(b.val, init, t)
			}
		default:
			if init.IsZero() {
				init = fb.zeroValue(t)
			}
			b.val = init
		}

		fb.bind(obj, b)
		if fb.needsDrop(obj.Type()) {
			fb.top.locals = append(fb.top.locals, b)
		}
	}
}

// forcesSlot is lowering.md §5's list. The join convention carries a
// reassigned name across blocks without memory, so a slot is needed only
// when an *address* is — unobservable except through `===`, which is why
// `===` is on the list.
//
// todo: three of the six entries need information hir does not have yet —
// "passed to a mut parameter", "operand of addr", and "captured by a closure
// that outlives it" are all uses, discovered after the declaration. Today
// every `var` gets a slot, which is correct and wasteful; narrowing it wants
// a use-collecting pre-pass over the body.
func (fb *funcBuilder) forcesSlot(obj types.Object, t *Type) bool {
	if IsAggregate(t) {
		return true
	}
	v, ok := obj.(*types.Var)
	return ok && v.Mutable()
}

func (fb *funcBuilder) assign(s *ast.AssignStmt) {
	if s.Op != token.ASSIGN {
		// The compound form takes exactly one target and one value.
		lhs := fb.expr(s.Targets[0])
		rhs := fb.expr(s.Values[0])
		op := binaryOpFor(compoundBase(s.Op), lhs.Type)
		res := fb.checkedBinary(op, lhs, rhs, s.OpPos)
		fb.assignTo(s.Targets[0], res)
		return
	}

	// Right-hand sides are all evaluated before any target is written: a
	// swap `a, b = b, a` must not observe a half-written pair.
	vals := make([]Value, 0, len(s.Values))
	for _, v := range s.Values {
		vals = append(vals, fb.owning(v))
	}
	for i, t := range s.Targets {
		if i < len(vals) {
			fb.assignTo(t, vals[i])
		}
	}
}

func (fb *funcBuilder) assignTo(target ast.Expr, v Value) {
	if id, ok := target.(*ast.Ident); ok {
		if id.IsBlank() {
			return // evaluated and discarded
		}
		obj := fb.l.info.ObjectOf(id)
		b := fb.lookup(obj)
		if b == nil {
			fb.l.bugAt(id.Pos(), "assignment to an unbound name "+id.Name)
		}
		if b.slot {
			// The old value dies here. A struct field holding a []T is not
			// leaked by an overwrite.
			fb.dropInPlace(b.val, b.vt)
			fb.storeInto(b.val, v, b.typ)
			return
		}
		// Join Convention: reassigning the name is the merge. No memory
		// involved, no phi to reconstruct.
		fb.seq.add(&Instr{Result: b.val.Name, Op: OpBitcast, Type: b.typ, Args: []Value{v}})
		return
	}

	// A field path, an index, or a dereference: all resolve to an address.
	p := fb.address(target)
	t := fb.typ(fb.l.info.TypeOf(target))
	fb.dropInPlace(p, fb.l.info.TypeOf(target))
	fb.storeInto(p, v, t)
}

func (fb *funcBuilder) ifStmt(s *ast.IfStmt) {
	cond := fb.expr(s.Cond)
	n := &If{Cond: fb.narrowBool(cond), Then: &Seq{}}

	outer := fb.seq
	fb.seq = n.Then
	fb.block(s.Body)
	if s.Else != nil {
		n.Else = &Seq{}
		fb.seq = n.Else
		fb.stmt(s.Else)
	}
	fb.seq = outer
	fb.seq.add(n)
}

func (fb *funcBuilder) whileStmt(s *ast.WhileStmt) {
	loop := &Loop{Body: &Seq{}}
	outer := fb.seq
	fb.seq = loop.Body
	fb.openScope(scopeLoop)

	// The head test is inside the loop: `while c { … }` is
	// `loop { if !c { break }; … }`. One shape for every loop primitive,
	// which is what lets flatten.go stay short.
	cond := fb.narrowBool(fb.expr(s.Cond))
	brk := &If{Cond: cond, Then: &Seq{}, Else: &Seq{}}
	brk.Else.add(&Break{})
	fb.seq.add(brk)

	inner := fb.seq
	fb.seq = brk.Then
	fb.block(s.Body)
	fb.seq = inner

	fb.closeScope()
	fb.seq = outer
	fb.seq.add(loop)
}

// forStmt lowers the single loop shape over an iterable.
//
// Only the range and fixed-array/slice forms are here. The mode marker sits
// on the binding rather than the iterable, because what transfers is each
// element, one per iteration.
func (fb *funcBuilder) forStmt(s *ast.ForStmt) {
	it := fb.l.info.TypeOf(s.X)

	if rng, ok := s.X.(*ast.BinaryExpr); ok && rng.Op == token.DOTDOT {
		fb.forRange(s, rng)
		return
	}
	switch types.Underlying(it).(type) {
	case *types.Array, *types.Slice:
		fb.forIndexed(s, it)
	default:
		// map iteration needs the runtime cursor; string iteration needs
		// UTF-8 decode at variable stride.
		fb.l.todoAt(s.Pos(), "for over "+types.TypeString(it))
	}
}

// forRange is one counter and nothing else. A range is never materialized —
// it has no type and cannot be bound, returned, passed, or stored, so it
// becomes the counter's bounds.
func (fb *funcBuilder) forRange(s *ast.ForStmt, rng *ast.BinaryExpr) {
	lo := fb.expr(rng.X)
	hi := fb.expr(rng.Y)
	t := lo.Type

	i := fb.alloca(t, "i")
	fb.emitVoid(OpStore, t, i, lo)

	loop := &Loop{Body: &Seq{}}
	outer := fb.seq
	fb.seq = loop.Body
	fb.openScope(scopeLoop)

	cur := fb.emit(OpLoad, t, i)
	// Always exclusive of the upper bound, and empty when lo >= hi.
	cmp := fb.emit(cmpOp(OpSlt, t), I1, cur, hi)
	brk := &If{Cond: cmp, Then: &Seq{}, Else: &Seq{}}
	brk.Else.add(&Break{})
	fb.seq.add(brk)

	inner := fb.seq
	fb.seq = brk.Then
	if obj := fb.l.info.Defs[s.Names[0]]; obj != nil && !s.Names[0].IsBlank() {
		fb.bind(obj, &binding{obj: obj, val: cur, typ: t, vt: obj.Type()})
	}
	fb.block(s.Body)
	next := fb.emit(OpAdd, t, cur, Int(1, t))
	fb.emitVoid(OpStore, t, i, next)
	fb.seq = inner

	fb.closeScope()
	fb.seq = outer
	fb.seq.add(loop)
}

// forIndexed is counter plus index.ptr, for [N]T and []T.
//
// todo: the consuming form (`for var x in xs`) moves each element out and
// emits no per-element drop, leaving the container's teardown to drop the
// unvisited tail [i, len) — which is exactly what a break out of a consuming
// loop leaves behind. Not implemented; a consuming for currently lowers as
// the shared form and double-drops.
func (fb *funcBuilder) forIndexed(s *ast.ForStmt, it types.Type) {
	if s.Mode == token.VAR {
		fb.l.todoAt(s.Pos(), "consuming for loop")
	}
	fb.l.todoAt(s.Pos(), "for over an array or slice")
	_ = it
}

func (fb *funcBuilder) switchStmt(s *ast.SwitchStmt) {
	tag := fb.expr(s.Tag)

	// An enum subject dispatches on the tag, not the value.
	if e := types.AsEnum(fb.l.info.TypeOf(s.Tag)); e != nil && !e.UnitOnly() {
		p := fb.address(s.Tag)
		tag = fb.emit(OpLoad, fb.l.basic(e.Discriminant()), p)
	}

	n := &SwitchStmt{Value: tag, Default: &Seq{}}
	outer := fb.seq
	sawDefault := false

	for _, cl := range s.Cases {
		body := &Seq{}
		fb.seq = body
		fb.openScope(scopeBlock)
		for _, st := range cl.Body {
			fb.stmt(st)
			if fb.terminated() {
				break
			}
		}
		if !fb.terminated() {
			fb.epilogueTo(fb.top)
		}
		fb.closeScope()

		if cl.Patterns == nil {
			n.Default = body
			sawDefault = true
			continue
		}
		for _, p := range cl.Patterns {
			v, ok := fb.patternValue(p)
			if !ok {
				fb.l.todoAt(p.Pos(), "non-constant switch pattern")
			}
			n.Cases = append(n.Cases, SwitchCase{Value: v, Body: body})
		}
	}

	if !sawDefault {
		// An exhaustive enum switch emits no default edge of its own, but
		// vir's switch terminator requires a default label — so it targets a
		// block ending in unreachable.
		n.Default.add(&UnreachableStmt{})
	}

	fb.seq = outer
	fb.seq.add(n)
}

func (fb *funcBuilder) patternValue(p ast.Expr) (int64, bool) {
	if ep, ok := p.(*ast.EnumPattern); ok {
		if len(ep.Binds) > 0 {
			// Payload bindings are field.ptr plus a reinterpretation into
			// the payload — views, not copies, and not assignable through.
			return 0, false
		}
		e := types.AsEnum(fb.l.info.TypeOf(p))
		if e == nil {
			return 0, false
		}
		if v := e.LookupVariant(ep.Name.Name); v != nil {
			return v.Value, true
		}
		return 0, false
	}
	if tv, ok := fb.l.info.Types[p]; ok && tv.Value != nil {
		return types.Int64Val(tv.Value)
	}
	return 0, false
}

func (fb *funcBuilder) returnStmt(s *ast.ReturnStmt) {
	switch {
	case len(s.Results) == 0:
		fb.epilogueTo(nil)
		fb.seq.add(&ReturnStmt{})

	case fb.sret != "":
		// An aggregate result is written through the sret pointer, then the
		// epilogue runs, then a bare return. Failure costs exactly what
		// success costs: the same stores into the same slot.
		dst := Name(fb.sret, Ptr)
		if len(s.Results) == 1 {
			v := fb.owning(s.Results[0])
			fb.storeInto(dst, v, fb.fn.SRet)
		} else {
			st := fb.fn.SRet.Struct
			for i, r := range s.Results {
				v := fb.owning(r)
				p := fb.fieldPtr(dst, st, i)
				fb.storeInto(p, v, st.Fields[i].Type)
			}
		}
		fb.epilogueTo(nil)
		fb.seq.add(&ReturnStmt{})

	default:
		v := fb.owning(s.Results[0])
		// The epilogue runs after the value is computed and before the
		// return, which is why a returned local is not dropped: owning() on
		// a transferred binding already marked it dead, and a returned
		// non-transferred local is copied.
		fb.epilogueTo(nil)
		fb.seq.add(&ReturnStmt{Value: &v})
	}
}

// deferStmt registers a call. Its callee and arguments are evaluated at the
// defer statement, so registration is static: no runtime defer list, no
// mask, just duplication of the deferred call onto each exit edge with the
// evaluated arguments held in slots.
func (fb *funcBuilder) deferStmt(s *ast.DeferStmt) {
	callee, mod, args := fb.callTarget(s.Call)
	held := make([]Value, 0, len(args))
	for _, a := range args {
		slot := fb.alloca(a.Type, "d")
		fb.emitVoid(OpStore, a.Type, slot, a)
		held = append(held, fb.emit(OpLoad, a.Type, slot))
	}
	// A deferred call runs when the enclosing *function* returns, not the
	// enclosing block — so it registers against the function scope.
	fn := fb.top
	for fn.parent != nil {
		fn = fn.parent
	}
	fn.defers = append(fn.defers, &deferred{callee: callee, module: mod, args: held})
}