package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// funcBuilder lowers one function body. It owns the scope stack, and
// therefore the epilogue machinery — see the known deviation in lower.go.
type funcBuilder struct {
	l    *lowerer
	fn   *Func
	sub  *subst
	seq  *Seq   // where instructions currently land
	top  *scope // innermost open scope
	sret string // the sret parameter's name, "" for a thin or void result
}

// scope is one lexical block's worth of teardown obligation.
//
// A binding lives here in declaration order; teardown runs in reverse. A
// transferred binding is marked dead and omitted — no flag is set and none
// is read, which is the whole reason conditional transfer is a compile
// error upstream.
type scope struct {
	parent  *scope
	kind    scopeKind
	locals  []*binding
	defers  []*deferred
}

type scopeKind uint8

const (
	scopeBlock scopeKind = iota
	scopeLoop
	scopeFunc
)

type binding struct {
	obj  types.Object
	val  Value // the named value, or the slot pointer when slot is true
	typ  *Type
	vt   types.Type
	slot bool
	dead bool
}

type deferred struct {
	callee string
	module string
	args   []Value
}

func (l *lowerer) lowerFuncDecl(in *instance, d *ast.FuncDecl) {
	sig := in.obj.Signature()
	if sig == nil {
		l.bugAt(d.Pos(), "function reached lowering with no signature")
	}
	if m := sig.Marker(); m == types.MarkerGPU || m == types.MarkerNPU {
		// The launch call site validates; the kernel body has nowhere to go.
		// There is no lower/gvir counterpart and vvm has no host-device
		// story to receive one.
		l.todoAt(d.Pos(), "gpu/npu function body")
	}
	if sig.Marker() == types.MarkerAsync {
		l.todoAt(d.Pos(), "async function body (see async.go)")
	}

	fb := &funcBuilder{
		l:   l,
		fn:  in.fn,
		sub: &subst{params: genericParams(in.obj), args: in.args},
		seq: &Seq{},
	}
	fb.signature(sig, d)
	fb.openScope(scopeFunc)

	// The receiver and the parameters belong to the body's own scope. A
	// parameter is not torn down by the callee unless it is owning: shared
	// and mut positions leave the value with the caller.
	fb.bindParams(sig, d)

	fb.block(d.Body)

	// A void function's fall-through exit still runs its epilogue, and a
	// value one already emitted its own at each return. Emitting here
	// unconditionally is safe: a body ending in return leaves nothing in
	// scope to tear down twice, because returnStmt closes the scopes as it
	// goes.
	if !fb.terminated() {
		fb.epilogueTo(nil)
		fb.seq.add(&ReturnStmt{})
	}
	fb.closeScope()
	in.fn.Body = fb.seq
}

// signature fills Params/Result. Three conventions, and the whole of the
// call-side cost:
//
//	x: T        shared    thin -> value, aggregate -> ptr
//	x: mut T    exclusive ptr to the caller's slot
//	x: var T    owning    thin -> value, aggregate -> byval
func (fb *funcBuilder) signature(sig *types.Signature, d *ast.FuncDecl) {
	if r := sig.Recv(); r != nil {
		fb.fn.Params = append(fb.fn.Params, fb.param(r))
	}

	// An aggregate result is written through a synthetic first parameter.
	// This applies to boundary tuples too — (int32, string) included.
	if sig.Results().Len() > 0 {
		rt := fb.resultType(sig)
		if IsAggregate(rt) {
			fb.sret = "ret"
			fb.fn.SRet = rt
			fb.fn.Params = append(fb.fn.Params, &Param{
				Name: "ret", Type: Ptr, SRet: rt.Struct,
			})
		} else {
			fb.fn.Result = rt
		}
	}

	for i := 0; i < sig.Params().Len(); i++ {
		fb.fn.Params = append(fb.fn.Params, fb.param(sig.Params().At(i)))
	}
	_ = d
}

func (fb *funcBuilder) resultType(sig *types.Signature) *Type {
	res := sig.Results()
	if res.Len() == 1 {
		return fb.typ(res.At(0).Type())
	}
	// A multi-slot result is a tuple, and a boundary tuple is the common
	// case. Both are one struct returned by sret.
	vars := make([]*types.Var, res.Len())
	for i := range vars {
		vars[i] = res.At(i)
	}
	return fb.l.typ(types.NewTuple(vars...), fb.fn.Module)
}

func (fb *funcBuilder) param(v *types.Var) *Param {
	t := fb.typ(v.Type())
	p := &Param{Name: v.Name(), Type: t}
	if p.Name == "" || p.Name == "_" {
		p.Name = fb.fn.fresh("p")
	}
	switch {
	case v.Mode() == types.ModeMut:
		// A mut parameter is a pointer to the caller's slot, which is why
		// the caller's binding must be a var.
		p.Type = Ptr
	case IsAggregate(t) && v.Mode() == types.ModeVar:
		// Owning aggregates use byval: vir requires the caller's object to
		// stay live and unmutated for the call, which is exactly true after
		// a transfer.
		p.ByVal = t.Struct
		p.Type = Ptr
	case IsAggregate(t):
		// Shared aggregates pass as a bare ptr. vir has no noalias or
		// per-parameter readonly, so exclusivity is discharged in the front
		// end and is not transmitted — the optimizer cannot exploit it.
		p.Type = Ptr
	}
	return p
}

func (fb *funcBuilder) bindParams(sig *types.Signature, d *ast.FuncDecl) {
	i := 0
	if r := sig.Recv(); r != nil {
		fb.bindParam(r, fb.fn.Params[i])
		i++
	}
	if fb.sret != "" {
		i++
	}
	for j := 0; j < sig.Params().Len(); j++ {
		if i < len(fb.fn.Params) {
			fb.bindParam(sig.Params().At(j), fb.fn.Params[i])
		}
		i++
	}
	_ = d
}

func (fb *funcBuilder) bindParam(v *types.Var, p *Param) {
	if v.Name() == "" || v.Name() == "_" {
		return
	}
	t := fb.typ(v.Type())
	b := &binding{
		obj:  v,
		val:  Name(p.Name, p.Type),
		typ:  t,
		vt:   v.Type(),
		slot: IsAggregate(t) || v.Mode() == types.ModeMut,
	}
	// Only an owning parameter is torn down by the callee. Shared and mut
	// leave the value with the caller.
	if v.Mode() == types.ModeVar {
		fb.top.locals = append(fb.top.locals, b)
	}
	fb.bind(v, b)
}

// ------------------------------------------------------------------ emission

func (fb *funcBuilder) emit(op Op, t *Type, args ...Value) Value {
	name := ""
	if !op.IsVoid() && !IsVoid(t) {
		name = fb.fn.fresh("v")
	}
	in := &Instr{Result: name, Op: op, Type: t, Args: args}
	fb.seq.add(in)
	if name == "" {
		return Value{}
	}
	return Name(name, t)
}

func (fb *funcBuilder) emitVoid(op Op, t *Type, args ...Value) {
	fb.seq.add(&Instr{Op: op, Type: t, Args: args})
}

// alloca hoists every slot into the entry block. vir allocas are
// per-execution and accumulate per loop iteration, so a slot written inside
// a loop body is allocated once, before the loop, and reused — sound because
// a loop-body local's teardown runs on every back edge, leaving the slot
// dead at the top of each iteration. Vertex has no dynamically sized local,
// so no alloca ever depends on a runtime value.
func (fb *funcBuilder) alloca(t *Type, hint string) Value {
	name := fb.fn.fresh(hint)
	fb.fn.Allocas = append(fb.fn.Allocas, &Instr{
		Result: name,
		Op:     OpAlloca,
		Type:   Ptr,
		Args:   []Value{Int(fb.l.sizeOf(t), I64)},
		Align:  int(fb.l.alignOf(t)),
	})
	return Name(name, Ptr)
}

// call emits a direct call within this module.
func (fb *funcBuilder) call(callee string, result *Type, args ...Value) Value {
	name := ""
	if !IsVoid(result) {
		name = fb.fn.fresh("r")
	}
	fb.seq.add(&Instr{Result: name, Op: OpCall, Type: result, Callee: callee, Args: args})
	if name == "" {
		return Value{}
	}
	return Name(name, result)
}

// callBuiltin emits a call into a builtins module, recording the import and
// the feature at the site that needs them.
//
// The Symbol value is the only thing hir ever names: never a "rc"/"retain"
// string pair. That is what keeps this package from seeing how the callee
// module is built.
func (fb *funcBuilder) callBuiltin(sym builtins.Symbol, result *Type, args ...Value) Value {
	fb.fn.Module.Import(sym.ImportPath())
	if f, ok := builtins.FeatureFor(sym); ok {
		fb.l.need(f)
	}
	name := ""
	if !IsVoid(result) {
		name = fb.fn.fresh("r")
	}
	fb.seq.add(&Instr{
		Result: name, Op: OpCall, Type: result,
		Module: sym.Module, Callee: sym.Func, Args: args,
	})
	if name == "" {
		return Value{}
	}
	return Name(name, result)
}

// callAcross emits a call into another Vertex package.
func (fb *funcBuilder) callAcross(m *Module, callee string, result *Type, args ...Value) Value {
	fb.fn.Module.Import(m.Path)
	name := ""
	if !IsVoid(result) {
		name = fb.fn.fresh("r")
	}
	fb.seq.add(&Instr{
		Result: name, Op: OpCall, Type: result,
		Module: m.Name, Callee: callee, Args: args,
	})
	if name == "" {
		return Value{}
	}
	return Name(name, result)
}

func (fb *funcBuilder) typ(t types.Type) *Type {
	return fb.l.typ(fb.sub.apply(t), fb.fn.Module)
}

func (fb *funcBuilder) pos(n ast.Node) token.Pos {
	if n == nil {
		return token.NoPos
	}
	return n.Pos()
}

func (fb *funcBuilder) terminated() bool {
	if len(fb.seq.List) == 0 {
		return false
	}
	switch fb.seq.List[len(fb.seq.List)-1].(type) {
	case *ReturnStmt, *Break, *Continue, *TrapStmt, *UnreachableStmt:
		return true
	}
	return false
}