package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// Call-site resolution, package-scope references, and the two expression
// forms whose lowering is a builtin call rather than an opcode.
//
// Everything here is reached from expr.go and nothing here is reached from
// anywhere else, which is why it is one file: the question "what does this
// call actually name" is asked once, in callTarget, and answered once.

// ------------------------------------------------------------ callee resolution

// calleeObject resolves a call's callee to the function object it names, or
// nil when the callee is a value rather than a name. Every indirect form —
// a closure, a stored function, a foreign class method — lands on nil, and
// callTarget turns that into the one todo.
func (fb *funcBuilder) calleeObject(x *ast.CallExpr) *types.Func {
	return fb.funcOf(x.Fun)
}

func (fb *funcBuilder) funcOf(e ast.Expr) *types.Func {
	switch f := unparen(e).(type) {
	case *ast.Ident:
		fn, _ := fb.l.info.ObjectOf(f).(*types.Func)
		return fn

	case *ast.SelectorExpr:
		// A method value, a package-qualified function, and a foreign
		// function are one shape here: the object hangs off Sel either way,
		// and which it was is a question the signature answers later.
		fn, _ := fb.l.info.ObjectOf(f.Sel).(*types.Func)
		return fn

	case *ast.IndexExpr:
		// Explicit type arguments. ast cannot tell TypeArgs from Index, so
		// the operand is what carries the name in both readings — and only
		// the TypeArgs reading ever reaches a callee position.
		return fb.funcOf(f.X)
	}
	return nil
}

// calleeSignature is the signature the call site is typed against. It is
// read from the object rather than from Info.TypeOf(x.Fun) so that a method
// and a free function answer the same way.
func (fb *funcBuilder) calleeSignature(x *ast.CallExpr) *types.Signature {
	if obj := fb.calleeObject(x); obj != nil {
		return obj.Signature()
	}
	return nil
}

// typeArgsFor is the concrete type arguments this call instantiates with.
//
// genericParams reports parameters only for a method on a generic type, so
// the arguments come from the receiver's own TypeArgs — `s.push(v)` where s
// is Stack[int32] instantiates push with [int32], and there is no inference
// to perform because the receiver already carries the answer. An explicit
// bracket list is accepted where one is written.
//
// fb.sub composes on descent: a receiver spelled Stack[T] inside a body
// being lowered for T=int32 must substitute before the args are read, or the
// worklist keys on a type parameter and mono.go's guard trips.
func (fb *funcBuilder) typeArgsFor(x *ast.CallExpr, obj *types.Func) []types.Type {
	params := genericParams(obj)
	if len(params) == 0 {
		return nil
	}

	if ix, ok := unparen(x.Fun).(*ast.IndexExpr); ok {
		out := make([]types.Type, 0, len(ix.Indices))
		for _, a := range ix.Indices {
			out = append(out, fb.sub.apply(fb.l.info.TypeOf(a)))
		}
		if len(out) == len(params) {
			return out
		}
	}

	if sel, ok := unparen(x.Fun).(*ast.SelectorExpr); ok {
		recv := fb.sub.apply(fb.l.info.TypeOf(sel.X))
		if n := types.AsNamed(recv); n != nil {
			if args := n.TypeArgs(); len(args) == len(params) {
				return args
			}
		}
	}

	fb.l.todoAt(x.Pos(), "generic call whose type arguments are neither written nor carried by the receiver")
	return nil
}

// orderArgs resolves named arguments into positional order. Named and
// positional are not mixed — the analyzer guaranteed that — so this is a
// permutation and never a merge, and it leaves no trace in the emitted call.
func (fb *funcBuilder) orderArgs(x *ast.CallExpr, sig *types.Signature) []ast.Expr {
	named := false
	for _, a := range x.Args {
		if _, ok := a.(*ast.KeyValueExpr); ok {
			named = true
			break
		}
	}
	if !named || sig == nil {
		return x.Args
	}

	n := sig.Params().Len()
	out := make([]ast.Expr, n)
	pos := 0
	for _, a := range x.Args {
		kv, ok := a.(*ast.KeyValueExpr)
		if !ok {
			if pos < n {
				out[pos] = a
			}
			pos++
			continue
		}
		key, _ := kv.Key.(*ast.Ident)
		if key == nil {
			fb.l.bugAt(kv.Pos(), "named argument with a non-identifier key")
		}
		idx := -1
		for i := 0; i < n; i++ {
			if sig.Params().At(i).Name() == key.Name {
				idx = i
				break
			}
		}
		if idx < 0 {
			fb.l.bugAt(kv.Pos(), "named argument names an absent parameter "+key.Name)
		}
		out[idx] = kv.Value
	}
	for i, e := range out {
		if e == nil {
			fb.l.bugAt(x.Pos(), "call leaves parameter "+sig.Params().At(i).Name()+" unset")
		}
	}
	return out
}

// --------------------------------------------------------- package-scope refs

// globalRef is a reference to package-scope storage. It is always an
// *address*: vir globals are memory, and a scalar global's use is a load
// that identValue emits and address() deliberately does not. A top-level
// scalar `let` never reaches here — decl.go folded it to a vir const, and
// identValue answers it from types.Const before asking for a binding.
func (fb *funcBuilder) globalRef(obj types.Object) Value {
	if obj == nil {
		fb.l.bug("global reference to a nil object")
	}
	return GlobalRef(fb.l.qualify(obj.Name()), Ptr)
}

// ------------------------------------------------------------------ strings

// stringBinary lowers an operator whose operands are strings. Both operands
// are _Vstr headers, so every one of these is a builtin call: there is no
// opcode over a two-word aggregate.
//
// todo: concatenation's result is a fresh allocation with no binding and no
// scope entry, so a discarded `a + b` leaks. It wants the temporary
// teardown that the epilogue lift for async.go has to introduce anyway —
// recorded rather than half-built here.
func (fb *funcBuilder) stringBinary(x *ast.BinaryExpr, lhs, rhs Value) Value {
	fb.l.need(builtins.FeatString)
	hdr := fb.l.header(fb.fn.Module, "_Vstr", field("p", Ptr), field("len", I64))

	switch x.Op {
	case token.ADD:
		out := fb.alloca(StructOf(hdr), "cat")
		fb.callBuiltin(builtins.StringConcat, Void, out, lhs, rhs)
		return out

	case token.IDENTICAL, token.NOT_IDENTICAL:
		// `===` on a string is payload identity, not content equality —
		// the one comparison that must not call compare.
		a := fb.emit(OpLoad, Ptr, fb.fieldPtr(lhs, hdr, 0))
		b := fb.emit(OpLoad, Ptr, fb.fieldPtr(rhs, hdr, 0))
		if x.Op == token.IDENTICAL {
			return fb.emit(OpEq, I1, a, b)
		}
		return fb.emit(OpNe, I1, a, b)
	}

	// compare yields negative / zero / positive, so every ordering is that
	// result against zero.
	c := fb.callBuiltin(builtins.StringCompare, I32, lhs, rhs)
	zero := Int(0, I32)
	var op Op
	switch x.Op {
	case token.EQL:
		op = OpEq
	case token.NEQ:
		op = OpNe
	case token.LSS:
		op = OpSlt
	case token.GTR:
		op = OpSgt
	case token.LEQ:
		op = OpSle
	case token.GEQ:
		op = OpSge
	default:
		fb.l.bugAt(x.OpPos, "unlowerable string operator")
	}
	return fb.emit(op, I1, c, zero)
}

// ------------------------------------------------------- heap constructors

// heapConstructor is `unique(x)`, `shared(x)`, and `weak(x)`.
//
// All three lower to one word. unique and shared allocate, and their failure
// panics rather than reporting: only new and resize fail politely, because
// only they were asked for an amount the program chose. weak allocates
// nothing — it increments the weak count on an existing control block.
func (fb *funcBuilder) heapConstructor(x *ast.HeapConstructor) Value {
	inner := fb.l.info.TypeOf(x.X)
	it := fb.typ(inner)
	size := Int(fb.l.sizeOf(it), I64)

	switch x.Kw {
	case token.UNIQUE:
		fb.l.need(builtins.FeatMemory)
		v := fb.owning(x.X)
		p := fb.callBuiltin(builtins.MemAllocate, Ptr, size)
		fb.oomUnless(p)
		fb.storeInto(p, v, it)
		return p

	case token.SHARED:
		fb.l.need(builtins.FeatRC)
		v := fb.owning(x.X)
		h := fb.callBuiltin(builtins.RCNew, Ptr, size)
		fb.oomUnless(h)
		// The handle is the control block; the payload sits past it, and
		// only rc knows by how much.
		payload := fb.callBuiltin(builtins.RCPayload, Ptr, h)
		fb.storeInto(payload, v, it)
		return h

	case token.WEAK:
		// The operand is an existing shared handle and is not consumed:
		// a weak reference observes, it does not own.
		fb.l.need(builtins.FeatRC)
		return fb.callBuiltin(builtins.RCWeak, Ptr, fb.expr(x.X))
	}

	fb.l.bugAt(x.Pos(), "unlowerable heap constructor")
	return Value{}
}

// oomUnless panics when p is null. The panic module's oom path allocates
// nothing, which is what makes it callable from the failure of an
// allocation.
func (fb *funcBuilder) oomUnless(p Value) {
	isNull := fb.emit(OpEq, I1, p, Null())
	n := &If{Cond: isNull, Then: &Seq{}}

	outer := fb.seq
	fb.seq = n.Then
	fb.callBuiltin(builtins.PanicOOM, Void)
	fb.seq.add(&UnreachableStmt{})
	fb.seq = outer
	fb.seq.add(n)
}

func unparen(e ast.Expr) ast.Expr {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.X
	}
}