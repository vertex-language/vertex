package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// funcBuilder lowers one function body. It owns the current statement
// sequence, the binding environment, and the scope stack that makes
// defer/deinit epilogue expansion possible.
type funcBuilder struct {
	l    *lowerer
	fn   *Func
	decl *ast.FuncDecl

	seq    *Seq
	scopes []*scope
	env    map[types.Object]*binding
}

// binding is how a source object is reached from lowered code.
//
// A.5.1 is the whole rule: `let` is immutable and "not guaranteed to be
// addressable — it may be a register, an SSA value, or folded away
// entirely," while `var` "is mutable and owns a real stack slot for its
// whole lifetime." So a let becomes a plain named value and a var becomes
// an alloca, which is also why only a var may be passed to a mut parameter:
// a let may not physically exist anywhere to point at.
type binding struct {
	Value  Value // the value itself (let) or the slot pointer (var)
	Type   Type  // the logical type of the binding
	Slot   bool  // Value is a pointer to storage
	Owning bool  // teardown must be emitted where liveness ends
	Src    types.Type
	Dead   bool // transferred: A.6.4's "teardown simply not emitted"
}

type scopeKind uint8

const (
	scopeBlock scopeKind = iota
	scopeLoop
	scopeFunc
)

// scope collects what has to run on every exit edge out of it. A.5.8
// requires deferred calls in reverse registration order on *every* exit
// edge, and vir has no unwinder — so "every exit edge" is a finite,
// statically known set, and a defer costs exactly the call it defers.
type scope struct {
	kind     scopeKind
	defers   []*Instr
	bindings []*binding
}

func newFuncBuilder(l *lowerer, fn *Func, decl *ast.FuncDecl) *funcBuilder {
	b := &funcBuilder{l: l, fn: fn, decl: decl, seq: &Seq{}, env: map[types.Object]*binding{}}
	fn.Body = b.seq
	b.push(scopeFunc)
	b.bindParams()
	return b
}

func (b *funcBuilder) info() *types.Info { return b.l.cur.unit.Info }
func (b *funcBuilder) mod() *Module      { return b.fn.Module }

func (b *funcBuilder) bindParams() {
	if b.decl == nil {
		return
	}
	i := 0
	if b.fn.Result == Void && len(b.fn.Params) > 0 && b.fn.Params[0].SRet != nil {
		i = 1
	}
	if b.decl.Recv != nil && b.decl.Recv.Name != nil {
		if obj := b.info().Defs[b.decl.Recv.Name]; obj != nil && i < len(b.fn.Params) {
			p := b.fn.Params[i]
			b.env[obj] = &binding{
				Value: Ref(p.Name, p.Type), Type: p.Type,
				Slot: p.ByVal != nil, Src: obj.Type(),
			}
		}
		i++
	}
	for _, ap := range b.decl.Type.Params.List {
		if ap.Name == nil || i >= len(b.fn.Params) {
			i++
			continue
		}
		obj := b.info().Defs[ap.Name]
		p := b.fn.Params[i]
		i++
		if obj == nil {
			continue
		}
		// A `mut` parameter is literally the pointer parameter (A.8.5,
		// A.3.2): the callee writes through it into the caller's slot.
		mut := b.l.isMutParam(obj)
		b.env[obj] = &binding{
			Value: Ref(p.Name, p.Type), Type: p.Type,
			Slot:   mut || p.ByVal != nil,
			Owning: b.l.isOwningParam(obj),
			Src:    obj.Type(),
		}
	}
}

// ---------------------------------------------------------------------------
// emission
// ---------------------------------------------------------------------------

// emit appends one instruction to the current straight-line run, opening a
// new run if the tail of the sequence is a control-flow statement.
func (b *funcBuilder) emit(in *Instr) Value {
	var run *Instrs
	if n := len(b.seq.List); n > 0 {
		run, _ = b.seq.List[n-1].(*Instrs)
	}
	if run == nil {
		run = &Instrs{}
		b.seq.add(run)
	}
	run.List = append(run.List, in)
	return in.Result()
}

func (b *funcBuilder) op(pos token.Pos, op Op, t Type, args ...Value) Value {
	return b.emit(&Instr{Pos: pos, Name: b.fn.fresh(op.String()), Op: op, Type: t, Args: args})
}

func (b *funcBuilder) opVoid(pos token.Pos, op Op, t Type, args ...Value) {
	b.emit(&Instr{Pos: pos, Op: op, Type: t, Args: args})
}

// call emits a direct call to a function in this program.
func (b *funcBuilder) call(pos token.Pos, f *Func, args ...Value) Value {
	c := &Callee{Name: f.Name}
	if f.Module != b.mod() {
		c.Module = f.Module.Name
		b.mod().AddImport(f.Module.Name)
	}
	name := ""
	if !IsVoid(f.Result) {
		name = b.fn.fresh("r")
	}
	return b.emit(&Instr{Pos: pos, Name: name, Op: OpCall, Type: f.Result, Args: args, Call: c})
}

// callExtern emits a call to a foreign entry point declared by a declare
// block. It is bare — no module qualifier — because a declare block is a
// linkage boundary, not a namespace: its symbols are injected into the
// file's current package (A.8.1).
func (b *funcBuilder) callExtern(pos token.Pos, name string, result Type, args ...Value) Value {
	rn := ""
	if !IsVoid(result) {
		rn = b.fn.fresh("r")
	}
	return b.emit(&Instr{Pos: pos, Name: rn, Op: OpCall, Type: result,
		Args: args, Call: &Callee{Name: name}})
}

// callBuiltin emits a qualified call into a builtin module and records the
// feature. lower/hir references builtin symbols through builtins' constants
// and never through string literals, so there is one place to grep and one
// place the build breaks when a signature changes.
func (b *funcBuilder) callBuiltin(pos token.Pos, s builtinSymbol, result Type, args ...Value) Value {
	b.l.needSymbol(s)
	b.mod().AddImport(s.Module)
	name := ""
	if !IsVoid(result) {
		name = b.fn.fresh(s.Func)
	}
	return b.emit(&Instr{Pos: pos, Name: name, Op: OpCall, Type: result, Args: args,
		Call: &Callee{Module: s.Module, Name: s.Func}})
}

// alloca creates a fresh stack slot. Slots are per-execution, so one in a
// loop body accumulates per iteration — every allocation site here is
// therefore hoisted to the function's first run, never emitted inside a
// Loop.
func (b *funcBuilder) alloca(pos token.Pos, t Type) Value {
	size := b.l.types.Sizeof(t)
	align := int(b.l.types.Alignof(t))
	in := &Instr{Pos: pos, Name: b.fn.fresh("slot"), Op: OpAlloca, Type: Ptr,
		Args: []Value{IntVal(I64, size)}, Align: align}
	// Hoist: prepend to the function's first straight-line run.
	root := b.fn.Body
	var first *Instrs
	if len(root.List) > 0 {
		first, _ = root.List[0].(*Instrs)
	}
	if first == nil {
		first = &Instrs{}
		root.List = append([]Stmt{first}, root.List...)
	}
	first.List = append([]*Instr{in}, first.List...)
	return in.Result()
}

func (b *funcBuilder) load(pos token.Pos, t Type, p Value) Value {
	if IsAggregate(t) {
		// An aggregate value *is* its address; there is nothing to load.
		return Value{Kind: ValRef, Name: p.Name, Type: t}
	}
	return b.op(pos, OpLoad, t, p)
}

func (b *funcBuilder) store(pos token.Pos, t Type, p, v Value) {
	if IsAggregate(t) {
		b.opVoid(pos, OpMemcopy, Void, p, v, IntVal(I64, b.l.types.Sizeof(t)))
		return
	}
	b.opVoid(pos, OpStore, t, p, v)
}

func (b *funcBuilder) fieldPtr(pos token.Pos, s *Struct, p Value, field string) Value {
	return b.emit(&Instr{Pos: pos, Name: b.fn.fresh(field), Op: OpFieldPtr, Type: Ptr,
		Args: []Value{p}, Field: &FieldRef{Struct: s, Field: field}})
}

func (b *funcBuilder) indexPtr(pos token.Pos, elem Type, p, idx Value) Value {
	return b.emit(&Instr{Pos: pos, Name: b.fn.fresh("elem"), Op: OpIndexPtr, Type: Ptr,
		Args: []Value{p, idx}, Elem: elem})
}

// loadField and storeField are the two shapes every fat-type access uses.
func (b *funcBuilder) loadField(pos token.Pos, s *Struct, p Value, name string) Value {
	f, _ := s.Field(name)
	return b.load(pos, f.Type, b.fieldPtr(pos, s, p, name))
}

func (b *funcBuilder) storeField(pos token.Pos, s *Struct, p Value, name string, v Value) {
	f, _ := s.Field(name)
	b.store(pos, f.Type, b.fieldPtr(pos, s, p, name), v)
}

// ---------------------------------------------------------------------------
// control flow and scopes
// ---------------------------------------------------------------------------

// into runs f with a fresh sequence as the emission target and returns it.
func (b *funcBuilder) into(f func()) *Seq {
	saved := b.seq
	b.seq = &Seq{}
	f()
	out := b.seq
	b.seq = saved
	return out
}

func (b *funcBuilder) push(k scopeKind) *scope {
	s := &scope{kind: k}
	b.scopes = append(b.scopes, s)
	return s
}

// pop closes the innermost scope, emitting its fall-through epilogue.
func (b *funcBuilder) pop() {
	n := len(b.scopes) - 1
	b.epilogue(b.scopes[n])
	b.scopes = b.scopes[:n]
}

// epilogue emits one scope's exit code: deferred calls first, in reverse
// registration order (A.5.8), then per-binding teardown in reverse
// declaration order (A.6.4). A transferred binding simply has its teardown
// not emitted — no flag is set and none is checked, which is exactly why
// conditional transfer is a compile error.
func (b *funcBuilder) epilogue(s *scope) {
	for i := len(s.defers) - 1; i >= 0; i-- {
		cp := *s.defers[i]
		b.emit(&cp)
	}
	for i := len(s.bindings) - 1; i >= 0; i-- {
		bd := s.bindings[i]
		if bd.Dead || !bd.Owning {
			continue
		}
		b.l.own.emitDeinit(b, 0, bd)
	}
}

// unwind emits the epilogues an exit edge crosses, innermost first, without
// closing the scopes — the same scope is still live for the fall-through
// path. This is what makes a return inside two nested blocks inside a loop
// run all three epilogues in the right order.
func (b *funcBuilder) unwind(upto scopeKind) {
	for i := len(b.scopes) - 1; i >= 0; i-- {
		b.epilogue(b.scopes[i])
		if b.scopes[i].kind == upto {
			return
		}
	}
}

func (b *funcBuilder) declareBinding(obj types.Object, bd *binding) {
	b.env[obj] = bd
	s := b.scopes[len(b.scopes)-1]
	s.bindings = append(s.bindings, bd)
}

func (b *funcBuilder) lookup(obj types.Object) *binding { return b.env[obj] }

// body lowers a function body and closes the function scope.
func (b *funcBuilder) body(block *ast.BlockStmt) {
	b.stmtList(block.List)
	b.unwind(scopeFunc)
	b.seq.add(&Return{})
}