package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/types"
)

// builtinCall lowers a reserved builtin. Each one's shape is arity- and
// type-argument-specific rather than expressible as a signature, which is
// why the checker special-cases each by Id and so does this.
func (fb *funcBuilder) builtinCall(b *types.Builtin, x *ast.CallExpr) Value {
	switch b.Id() {
	case types.SizeofId:
		t := fb.typ(fb.l.info.TypeOf(x.Args[0]))
		return Int(fb.l.sizeOf(t), I64)

	case types.AlignofId:
		t := fb.typ(fb.l.info.TypeOf(x.Args[0]))
		return Int(fb.l.alignOf(t), I64)

	case types.ReinterpretId:
		// Static reinterpretation, never a read. Emits nothing.
		return fb.expr(x.Args[1])

	case types.AddrId:
		return fb.address(x.Args[0])

	case types.PanicId:
		msg := fb.expr(x.Args[0])
		hdr := fb.l.header(fb.fn.Module, "_Vstr", field("p", Ptr), field("len", I64))
		p := fb.emit(OpLoad, Ptr, fb.fieldPtr(msg, hdr, 0))
		n := fb.emit(OpLoad, I64, fb.fieldPtr(msg, hdr, 1))
		fb.callBuiltin(builtins.Panic, Void, p, n)
		// panic does not return: deferred calls do not run, no deinit runs,
		// and there is no catch, recover, or unwind. noreturn plus
		// unreachable is exactly that, and it discharges "every path
		// returns" for free.
		fb.seq.add(&UnreachableStmt{})
		return Value{}

	case types.NewId:
		return fb.builtinNew(x)

	case types.DeleteId:
		fb.callBuiltin(builtins.MemFree, Void, fb.expr(x.Args[0]))
		return Value{}

	case types.CopyId:
		// memmove — always overlap-safe, with no unsafe variant to select.
		dst := fb.expr(x.Args[0])
		src := fb.expr(x.Args[1])
		n := fb.convert(fb.expr(x.Args[2]), I64)
		fb.emitVoid(OpMemmove, Void, dst, src, n)
		return Value{}

	case types.ZeroId:
		p := fb.expr(x.Args[0])
		n := fb.convert(fb.expr(x.Args[1]), I64)
		fb.emitVoid(OpMemset, Void, p, Int(0, U8), n)
		return Value{}

	case types.MinId, types.MaxId, types.ClampId:
		// todo: the scalar forms are smin/smax/umin/umax plus a compose for
		// clamp; the vector forms are the same opcodes over vec. Neither is
		// wired.
		fb.l.todoAt(x.Pos(), "min/max/clamp")

	case types.UpgradeId:
		// The increment-if-nonzero loop. It reports failure through the
		// boundary tuple, which is why a zero weak needs no separate check.
		out := fb.alloca(fb.typ(fb.l.info.TypeOf(x)), "up")
		fb.callBuiltin(builtins.RCUpgrade, Void, fb.expr(x.Args[0]), out)
		return out

	case types.DropId:
		// The release, emitted at the drop call instead of at scope exit.
		v := fb.expr(x.Args[0])
		fb.dropValue(v, fb.l.info.TypeOf(x.Args[0]))
		if b := fb.bindingOf(x.Args[0]); b != nil {
			b.dead = true
		}
		return Value{}

	case types.ResizeId, types.BlendId:
		fb.l.todoAt(x.Pos(), "resize / blend")

	case types.TransferId:
		// Reserved and bound to nothing, purely so a call spelled either way
		// diagnoses as a misspelled ownership marker. Reaching lowering means
		// the analyzer let one through.
		fb.l.bugAt(x.Pos(), "transfer() reached lowering")
	}
	fb.l.bugAt(x.Pos(), "unlowerable builtin")
	return Value{}
}

// builtinNew is the fallible allocation. It builds a boundary tuple:
// (typed_ptr T, string), zero-and-message on failure.
//
// Overflow of count * sizeof(T) is an allocation failure, not a trap — the
// one place a `*` overflow does not reach the trapping path, because count
// is caller data off a wire.
func (fb *funcBuilder) builtinNew(x *ast.CallExpr) Value {
	fb.l.need(builtins.FeatMemory)

	elem := fb.typ(fb.l.info.TypeOf(x))
	size := Int(fb.l.sizeOf(elem), I64)
	count := Int(1, I64)
	if len(x.Args) > 0 {
		count = fb.convert(fb.expr(x.Args[0]), I64)
	}

	out := fb.alloca(fb.typ(fb.l.info.TypeOf(x)), "new")
	fb.zero(out, fb.typ(fb.l.info.TypeOf(x)))

	n := fb.emit(OpMul, I64, count, size)
	over := fb.emit(OpUMulO, I1, count, size)

	p := fb.callBuiltin(builtins.MemAllocate, Ptr, n)
	isNull := fb.emit(OpEq, I1, p, Null())
	bad := fb.emit(OpOr, I1, over, isNull)

	branch := &If{Cond: bad, Then: &Seq{}, Else: &Seq{}}
	outer := fb.seq

	fb.seq = branch.Then
	// todo: write "out of memory" into the tuple's string slot. The message
	// is a literal, so this is a stringHeader plus two stores; it is left out
	// only because the tuple's field indices depend on the boundary-tuple
	// shape, which decl.go does not name yet.
	fb.seq = branch.Else
	fb.emitVoid(OpStore, Ptr, out, p)

	fb.seq = outer
	fb.seq.add(branch)
	return out
}