package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// builtinSymbol is builtins.Symbol with its feature already resolved, so a
// call site names one thing and the feature set stays honest.
type builtinSymbol = builtins.Symbol

var (
	symMemAllocate = builtins.MemAllocate
	symMemFree     = builtins.MemFree
	symMemResize   = builtins.MemResize
	symMemZero     = builtins.MemZero

	symRCNew     = builtins.RCNew
	symRCRetain  = builtins.RCRetain
	symRCRelease = builtins.RCRelease
	symRCWeak    = builtins.RCWeak
	symRCUpgrade = builtins.RCUpgrade

	symUniqueNew  = builtins.UniqueNew
	symUniqueFree = builtins.UniqueFree

	symSliceFree = builtins.SliceFree
	symSliceOOB  = builtins.SliceOOB

	symStringCopy    = builtins.StringCopy
	symStringFree    = builtins.StringFree
	symStringCompare = builtins.StringCompare

	symMapFree  = builtins.MapFree
	symChanRelease = builtins.ChanRelease
	symPanicOOM = builtins.PanicOOM
)

const builtinFeatureString = builtins.FeatString

// builtinCall lowers A.4.8's builtin call forms. Builtins get no node of
// their own — they are ordinary call syntax over reserved, unshadowable
// names — so dispatch is by the object's spelling.
func (b *funcBuilder) builtinCall(x *ast.CallExpr, bi *types.Builtin) Value {
	pos := x.Pos()
	switch bi.Name() {
	case "sizeof":
		t := b.info().TypeOf(x.Args[0])
		return IntVal(I64, b.l.types.Sizeof(b.l.hirType(t)))

	case "alignof":
		t := b.info().TypeOf(x.Args[0])
		return IntVal(I64, b.l.types.Alignof(b.l.hirType(t)))

	case "reinterpret":
		// A bit-cast between value types of identical byte size — never a
		// pointer cast, which is what `as` is for.
		to := b.l.hirType(b.info().TypeOf(x.Args[0]))
		return b.op(pos, OpBitcast, to, b.expr(x.Args[1]))

	case "addr":
		// addr accepts an addressable typed_ptr operand only; on any other
		// type &x is already the address.
		return b.addressOf(x.Args[0])

	case "new":
		return b.builtinNew(x)

	case "delete":
		b.callBuiltin(pos, symMemFree, Void, b.expr(x.Args[0]))
		return Value{}

	case "resize":
		b.l.todo(pos, "resize — invalidates its input on success, leaves it valid on failure (A.4.8)")
		return Value{}

	case "copy":
		// copy is always overlap-safe; there is deliberately no
		// overlap-unsafe variant.
		b.opVoid(pos, OpMemmove, Void, b.expr(x.Args[0]), b.expr(x.Args[1]), b.expr(x.Args[2]))
		return Value{}

	case "zero":
		b.opVoid(pos, OpMemset, Void, b.expr(x.Args[0]), IntVal(I8, 0), b.expr(x.Args[1]))
		return Value{}

	case "unique":
		return b.builtinBox(x, symUniqueNew)

	case "shared":
		return b.builtinBox(x, symRCNew)

	case "weak":
		return b.callBuiltin(pos, symRCWeak, Ptr, b.expr(x.Args[0]))

	case "upgrade":
		// Returns (shared T, string): the boundary-tuple convention applied
		// to a race the type system cannot statically win.
		st := b.l.hirType(b.info().TypeOf(x)).(StructType)
		dst := b.alloca(pos, st)
		b.callBuiltin(pos, symRCUpgrade, Void, dst, b.expr(x.Args[0]))
		return Value{Kind: ValRef, Name: dst.Name, Type: st}

	case "drop":
		// Explicitly ends a transferred binding's lifetime without
		// emitting its teardown.
		if bd := b.bindingFor(x.Args[0]); bd != nil {
			bd.Dead = true
		}
		return Value{}
	}
	b.l.todo(pos, "builtin %s", bi.Name())
	return Value{}
}

// builtinNew allocates count * sizeof(T) bytes and returns
// (typed_ptr T, string). The block is zeroed by default; `zero: false` opts
// out and is a claim nothing checks. An allocation failure — including a
// byte extent that overflows, and a non-power-of-two alignment — is a null
// pointer and a non-empty string, not a distinct diagnostic.
func (b *funcBuilder) builtinNew(x *ast.CallExpr) Value {
	pos := x.Pos()
	elem := b.l.elem(b.info().TypeOf(x))
	esz := b.l.types.Sizeof(b.l.hirType(elem))

	count := IntVal(I64, 1)
	zero := true
	var align Value
	for i, a := range x.Args {
		kv, ok := a.(*ast.KeyValueExpr)
		if !ok {
			if i == 0 {
				count = b.expr(a)
			}
			continue
		}
		key, _ := kv.Key.(*ast.Ident)
		switch {
		case key == nil:
		case key.Name == "align":
			align = b.expr(kv.Value)
		case key.Name == "zero":
			if tv, ok := b.info().Types[kv.Value]; ok && tv.Value != nil {
				zero, _ = types.BoolVal_(tv.Value)
			}
		}
	}

	size := b.op(pos, OpMul, I64, count, IntVal(I64, esz))
	var p Value
	if align.Valid() {
		p = b.callBuiltin(pos, builtins.MemAlignedAllocate, Ptr, size, align)
	} else {
		p = b.callBuiltin(pos, symMemAllocate, Ptr, size)
	}
	if zero {
		n := &If{Cond: b.op(pos, OpNe, I1, p, NullVal())}
		n.Then = b.into(func() {
			b.opVoid(pos, OpMemset, Void, p, IntVal(I8, 0), size)
		})
		b.seq.add(n)
	}

	// The result is the standard error tuple: on failure the handle is
	// zeroed and the string carries a message.
	st, ok := b.l.hirType(b.info().TypeOf(x)).(StructType)
	if !ok {
		return p
	}
	slot := b.alloca(pos, st)
	b.storeField(pos, st.Def, slot, st.Def.Fields[0].Name, p)
	return Value{Kind: ValRef, Name: slot.Name, Type: st}
}

// builtinBox is unique(...) / shared(...) — the only two heap doors. Each
// constructs a wrapper around a value, so the copy/transfer rule does not
// apply to its operand: the operand is moved in unconditionally, exactly as
// into any constructor.
func (b *funcBuilder) builtinBox(x *ast.CallExpr, sym builtins.Symbol) Value {
	pos := x.Pos()
	inner := b.info().TypeOf(x.Args[0])
	hi := b.l.hirType(inner)
	h := b.callBuiltin(pos, sym, Ptr, IntVal(I64, b.l.types.Sizeof(hi)))

	oom := &If{Cond: b.op(pos, OpEq, I1, h, NullVal())}
	oom.Then = b.into(func() {
		// A.10.1 fixes the split: new/resize fail politely, while these
		// panic, matching native array allocation.
		b.callBuiltin(pos, symPanicOOM, Void)
		b.seq.add(&Unreachable{})
	})
	b.seq.add(oom)

	v := b.expr(x.Args[0])
	if bd := b.bindingFor(x.Args[0]); bd != nil {
		bd.Dead = true
	}
	b.store(pos, hi, h, v)
	return h
}

// callBuiltin emits a qualified call into a builtin module and records the
// feature. lower/hir references builtin symbols through builtins' constants
// and never through string literals, so there is one place to grep and one
// place the build breaks when a signature changes.
func (b *funcBuilder) callBuiltin(pos token.Pos, s builtinSymbol, result Type, args ...Value) Value {
	b.l.needSymbol(s)
	// AddImport takes the full "namespace/module" path (§7.3) — every
	// builtin module declares Namespace ("builtins", names.go), so the
	// bare module name alone would leave vvm's importer with nothing to
	// resolve against ("import \"memory\" does not resolve to any known
	// module"). The call's own Callee.Module below stays the bare module
	// name deliberately: §2.3's qualified-ident grammar is
	// `ident "." ident`, never namespace-qualified, so "memory.allocate"
	// is the correct call-site spelling even though "builtins/memory" is
	// what the import line needs to say to make that name resolvable.
	b.mod().AddImport(s.ImportPath())
	name := ""
	if !IsVoid(result) {
		name = b.fn.fresh(s.Func)
	}
	return b.emit(&Instr{Pos: pos, Name: name, Op: OpCall, Type: result, Args: args,
		Call: &Callee{Module: s.Module, Name: s.Func}})
}