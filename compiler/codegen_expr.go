package compiler

import (
	"strconv"
	"strings"

	"github.com/vertex-language/compiler/wasm"
	"github.com/vertex-language/vertex/parser"
)

// genExpr emits code for an expression and returns its Vertex type.
func (fg *funcGen) genExpr(ctx parser.IExprContext) Type {
	if ctx == nil {
		return Void
	}

	exprs := ctx.AllExpr()

	// ── Primary leaf ────────────────────────────────────────────────────────
	if primary := ctx.Primary(); primary != nil &&
		len(exprs) == 0 &&
		ctx.ArgList() == nil &&
		ctx.DOT() == nil &&
		ctx.LBRACKET() == nil {
		return fg.genPrimary(primary)
	}

	// ── Type cast: PrimitiveType(expr)  e.g. int(x), uint16(port) ──────────
	if ctx.PrimitiveType() != nil && len(exprs) == 1 {
		srcType := fg.genExpr(exprs[0])
		return fg.emitCast(ctx.PrimitiveType(), srcType)
	}

	// ── Unary ───────────────────────────────────────────────────────────────
	if len(exprs) == 1 &&
		ctx.ArgList() == nil &&
		ctx.DOT() == nil &&
		ctx.LBRACKET() == nil {

		switch {
		case ctx.MINUS() != nil:
			t := fg.genExpr(exprs[0])
			fg.emitNeg(t)
			return t
		case ctx.BANG() != nil:
			fg.genExpr(exprs[0])
			fg.body.I32Eqz()
			return Bool
		case ctx.TILDE() != nil:
			t := fg.genExpr(exprs[0])
			fg.body.I32Const(-1)
			fg.body.I32Xor()
			return t
		case ctx.AMP() != nil:
			t := fg.genExpr(exprs[0])
			return &PointerType{Elem: t, Mutable: true}
		}
	}

	// ── Postfix: expr.name  /  expr.name(args) ──────────────────────────────
	if ctx.DOT() != nil && ctx.PostfixName() != nil && len(exprs) == 1 {
		field := ctx.PostfixName().GetText()
		return fg.genPostfix(exprs[0], field, ctx.ArgList())
	}

	// ── Function call: expr(args), no dot ───────────────────────────────────
	if ctx.ArgList() != nil && ctx.DOT() == nil && len(exprs) == 1 {
		return fg.genCall(exprs[0], ctx.ArgList())
	}

	// ── Subscript: expr[expr] ────────────────────────────────────────────────
	if ctx.LBRACKET() != nil && len(exprs) == 2 {
		return fg.genSubscript(exprs[0], exprs[1])
	}

	// ── Binary operators ─────────────────────────────────────────────────────
	if len(exprs) == 2 {
		return fg.genBinary(ctx, exprs[0], exprs[1])
	}

	// ── Ternary: cond ? a : b ────────────────────────────────────────────────
	if ctx.QUESTION() != nil && ctx.COLON() != nil && len(exprs) == 3 {
		return fg.genTernary(exprs[0], exprs[1], exprs[2])
	}

	// ── Nil-coalesce: a ?? b ─────────────────────────────────────────────────
	if ctx.NIL_COALESCE() != nil && len(exprs) == 2 {
		t := fg.genExpr(exprs[0])
		fg.body.Drop()
		fg.genExpr(exprs[1])
		return t
	}

	return Void
}

// ── Primary ──────────────────────────────────────────────────────────────────

func (fg *funcGen) genPrimary(ctx parser.IPrimaryContext) Type {
	switch {
	case ctx.Literal() != nil:
		return fg.genPrimaryLiteral(ctx.Literal())

	// StructLiteralExpr must be checked BEFORE the bare ID case because
	// the grammar rule `structLiteralExpr : ID LBRACE …` also starts with an
	// identifier. ANTLR places the struct literal in its own sub-context, so
	// PrimaryContext.ID() is nil when StructLiteralExpr() is non-nil.
	case ctx.StructLiteralExpr() != nil:
		return fg.genStructLiteralExpr(ctx.StructLiteralExpr())

	case ctx.ID() != nil:
		return fg.genIdent(ctx.ID().GetText())

	case ctx.DOT() != nil && ctx.ID() != nil:
		enumCase := ctx.ID().GetText()
		fg.body.I32Const(fg.resolveEnumCaseVal(enumCase))
		return Int

	case ctx.LPAREN() != nil && ctx.Expr() != nil:
		return fg.genExpr(ctx.Expr())

	case ctx.TupleExpr() != nil:
		fg.body.I32Const(0)
		return Void

	case ctx.ArrayLiteralExpr() != nil:
		return fg.genArrayLit(ctx.ArrayLiteralExpr())

	case ctx.ArrayConstructExpr() != nil:
		fg.body.I32Const(0)
		return &ArrayType{Elem: Void}

	case ctx.AnonFuncExpr() != nil:
		fg.body.I32Const(0)
		return &FuncType{}
	}
	return Void
}

// genStructLiteralExpr allocates a temporary frame slot for an inline struct
// literal (one that appears in expression position rather than directly in a
// var decl), initialises all fields, and leaves the base address on the stack.
//
// Example:   if (Point{x: 1, y: 2} == p) { … }
func (fg *funcGen) genStructLiteralExpr(ctx parser.IStructLiteralExprContext) Type {
	typeName := ctx.ID().GetText()
	sym := fg.cg.scope.Lookup(typeName)
	if sym == nil || sym.Kind != SymType {
		fg.body.I32Const(0)
		return Void
	}
	t := sym.Type
	size := SizeOf(t)
	if size <= 0 {
		fg.body.I32Const(0)
		return t
	}

	// Allocate a static frame slot for this temporary.
	addr := fg.cg.allocFrame(size)
	fg.genStructLitFields(addr, ctx, t)

	// Leave the base address on the wasm stack.
	fg.body.I32Const(addr)
	return t
}

func (fg *funcGen) genPrimaryLiteral(ctx parser.ILiteralContext) Type {
	switch {
	case ctx.DEC_INT_LIT() != nil:
		v, _ := strconv.ParseInt(strings.ReplaceAll(ctx.DEC_INT_LIT().GetText(), "_", ""), 10, 64)
		fg.body.I32Const(int32(v))
		return Int
	case ctx.HEX_INT_LIT() != nil:
		s := strings.ReplaceAll(ctx.HEX_INT_LIT().GetText()[2:], "_", "")
		v, _ := strconv.ParseInt(s, 16, 64)
		fg.body.I32Const(int32(v))
		return Int
	case ctx.BIN_INT_LIT() != nil:
		s := ctx.BIN_INT_LIT().GetText()[2:]
		v, _ := strconv.ParseInt(s, 2, 64)
		fg.body.I32Const(int32(v))
		return Int
	case ctx.OCT_INT_LIT() != nil:
		s := ctx.OCT_INT_LIT().GetText()[2:]
		v, _ := strconv.ParseInt(s, 8, 64)
		fg.body.I32Const(int32(v))
		return Int
	case ctx.DEC_FLOAT_LIT() != nil:
		v, _ := strconv.ParseFloat(ctx.DEC_FLOAT_LIT().GetText(), 32)
		fg.body.F32Const(float32(v))
		return Float
	case ctx.STRING_LIT() != nil:
		s := stripQuotes(ctx.STRING_LIT().GetText())
		off := fg.cg.internString(s)
		fg.body.I32Const(off)
		return StrType
	case ctx.RAW_STRING_LIT() != nil:
		s := strings.Trim(ctx.RAW_STRING_LIT().GetText(), "`")
		off := fg.cg.internString(s)
		fg.body.I32Const(off)
		return StrType
	case ctx.TRUE() != nil:
		fg.body.I32Const(1)
		return Bool
	case ctx.FALSE() != nil:
		fg.body.I32Const(0)
		return Bool
	case ctx.NIL() != nil:
		fg.body.I32Const(0)
		return &OptionalType{Elem: Void}
	}
	return Void
}

// genIdent loads a named symbol onto the wasm stack.
func (fg *funcGen) genIdent(name string) Type {
	sym := fg.scope.Lookup(name)
	if sym == nil {
		fg.body.I32Const(0)
		return Void
	}
	switch sym.Kind {
	case SymVar, SymParam:
		if sym.IsFrame {
			// Struct variable: push its base address.
			fg.body.I32Const(sym.FrameOff)
		} else {
			fg.body.LocalGet(sym.LocalIdx)
		}
		return sym.Type
	case SymFunc, SymNative:
		fg.body.I32Const(int32(sym.FuncIdx))
		return sym.Type
	case SymType:
		// Type name used in expression position (e.g. namespace for native call).
		fg.body.I32Const(0)
		return sym.Type
	}
	fg.body.I32Const(0)
	return Void
}

// ── Postfix operations ────────────────────────────────────────────────────────

func (fg *funcGen) genPostfix(recv parser.IExprContext, field string, args parser.IArgListContext) Type {
	switch field {
	case "sizeof":
		t := fg.resolveExprType(recv)
		fg.body.Drop()
		fg.body.I32Const(int32(SizeOf(t)))
		return Int

	case "alignof":
		t := fg.resolveExprType(recv)
		fg.body.Drop()
		fg.body.I32Const(int32(AlignOf(t)))
		return Int

	case "any":
		recvType := fg.genExpr(recv)
		return &PointerType{Elem: recvType, Mutable: false}

	case "byteSize":
		fg.genExpr(recv)
		fg.body.Drop()
		fg.body.I32Const(0)
		return Int

	case "len":
		fg.genExpr(recv)
		fg.body.Drop()
		fg.body.I32Const(0)
		return Int

	case "new":
		fg.genExpr(recv)
		fg.body.Drop()
		fg.body.I32Const(0)
		return Void

	case "delete":
		fg.genExpr(recv)
		fg.body.Drop()
		return Void

	case "await", "try", "rawValue", "dispatch", "spawn", "fork":
		return fg.genExpr(recv)

	default:
		recvType := fg.genExpr(recv)
		if args != nil {
			return fg.genMethodCall(recvType, field, args)
		}
		return fg.genFieldLoad(recvType, field)
	}
}

// genFieldLoad loads a struct/class field; the base address is already on the stack.
func (fg *funcGen) genFieldLoad(t Type, fieldName string) Type {
	var fields []*StructField
	switch v := t.(type) {
	case *StructType:
		fields = v.Fields
	case *ClassType:
		fields = v.Fields
	default:
		fg.body.Drop()
		fg.body.I32Const(0)
		return Void
	}
	for _, f := range fields {
		if f.Name == fieldName {
			if f.Offset != 0 {
				fg.body.I32Const(int32(f.Offset))
				fg.body.I32Add()
			}
			fg.emitLoad(f.Type, 0)
			return f.Type
		}
	}
	fg.body.Drop()
	fg.body.I32Const(0)
	return Void
}

// genMethodCall emits a method / associated-function call on a value.
//
// Receiver handling:
//   - SymNative  (native class method): the receiver is a dummy type-name
//     identifier (pushes 0); drop it before forwarding real arguments.
//   - SymFunc    (associated function §26): the receiver value IS the first
//     argument and must stay on the stack.
func (fg *funcGen) genMethodCall(recvType Type, method string, args parser.IArgListContext) Type {
	// Resolve symbol: prefer "TypeName.method", fall back to bare "method".
	typeName := ""
	switch v := recvType.(type) {
	case *StructType:
		typeName = v.Name
	case *ClassType:
		typeName = v.Name
	case *EnumType:
		typeName = v.Name
	}

	var sym *Symbol
	if typeName != "" {
		sym = fg.cg.scope.Lookup(typeName + "." + method)
	}
	if sym == nil {
		sym = fg.cg.scope.Lookup(method)
	}
	if sym == nil {
		// Unknown method: discard receiver and all arguments.
		fg.body.Drop()
		for _, arg := range args.AllArg() {
			fg.genExpr(arg.Expr())
			fg.body.Drop()
		}
		fg.body.I32Const(0)
		return Void
	}

	switch sym.Kind {
	case SymNative:
		// Native class method: receiver is a namespace placeholder — drop it.
		fg.body.Drop()
	case SymFunc:
		// Associated function (§26): receiver is the first argument — keep it.
		// For Void receivers (rare edge case), drop since there's nothing to pass.
		if recvType != nil && recvType.Kind() == KindVoid {
			fg.body.Drop()
		}
	default:
		fg.body.Drop()
	}

	for _, arg := range args.AllArg() {
		fg.genExpr(arg.Expr())
	}
	fg.body.Call(sym.FuncIdx)

	if ft, ok := sym.Type.(*FuncType); ok && ft.Sig != nil {
		return ft.Sig.Ret
	}
	return Void
}

// ── Direct function call ──────────────────────────────────────────────────────

func (fg *funcGen) genCall(funcExpr parser.IExprContext, args parser.IArgListContext) Type {
	name := fg.exprName(funcExpr)

	sym := fg.cg.scope.Lookup(name)
	if sym == nil {
		for _, arg := range args.AllArg() {
			fg.genExpr(arg.Expr())
			fg.body.Drop()
		}
		fg.body.I32Const(0)
		return Void
	}

	for _, arg := range args.AllArg() {
		fg.genExpr(arg.Expr())
	}
	fg.body.Call(sym.FuncIdx)

	if ft, ok := sym.Type.(*FuncType); ok && ft.Sig != nil {
		return ft.Sig.Ret
	}
	return Void
}

func (fg *funcGen) exprName(ctx parser.IExprContext) string {
	if primary := ctx.Primary(); primary != nil {
		if primary.ID() != nil {
			return primary.ID().GetText()
		}
	}
	return ""
}

func (fg *funcGen) resolveExprType(ctx parser.IExprContext) Type {
	if primary := ctx.Primary(); primary != nil {
		if primary.ID() != nil {
			name := primary.ID().GetText()
			if sym := fg.scope.Lookup(name); sym != nil {
				return sym.Type
			}
		}
	}
	return Void
}

func (fg *funcGen) resolveEnumCaseVal(caseName string) int32 {
	for _, sym := range fg.cg.scope.entries {
		if et, ok := sym.Type.(*EnumType); ok {
			for _, c := range et.Cases {
				if c.Name == caseName {
					return int32(c.IntVal)
				}
			}
		}
	}
	return 0
}

// ── Array literal ─────────────────────────────────────────────────────────────

func (fg *funcGen) genArrayLit(ctx parser.IArrayLiteralExprContext) Type {
	elems := ctx.ExprList()
	if elems == nil {
		fg.body.I32Const(0)
		return &ArrayType{Elem: Void}
	}
	for _, e := range elems.AllExpr() {
		fg.genExpr(e)
		fg.body.Drop()
	}
	fg.body.I32Const(0)
	return &ArrayType{Elem: Int}
}

// ── Subscript ─────────────────────────────────────────────────────────────────

func (fg *funcGen) genSubscript(arr, idx parser.IExprContext) Type {
	arrType := fg.genExpr(arr)
	fg.genExpr(idx)
	elemType := Void
	elemSize := 4
	if at, ok := arrType.(*ArrayType); ok {
		elemType = at.Elem
		elemSize = SizeOf(elemType)
	}
	fg.body.I32Const(int32(elemSize))
	fg.body.I32Mul()
	fg.body.I32Add()
	fg.emitLoad(elemType, 0)
	return elemType
}

// ── Binary operators ──────────────────────────────────────────────────────────

func (fg *funcGen) genBinary(
	ctx parser.IExprContext,
	left, right parser.IExprContext,
) Type {
	ltype := fg.genExpr(left)
	fg.genExpr(right)

	switch {
	case ctx.PLUS() != nil:
		fg.emitAdd(ltype)
		return ltype
	case ctx.MINUS() != nil:
		fg.emitSub(ltype)
		return ltype
	case ctx.STAR() != nil:
		fg.emitMul(ltype)
		return ltype
	case ctx.SLASH() != nil:
		fg.emitDiv(ltype)
		return ltype
	case ctx.PERCENT() != nil:
		fg.emitRem(ltype)
		return ltype
	case ctx.OVERFLOW_ADD() != nil:
		fg.emitAdd(ltype)
		return ltype
	case ctx.OVERFLOW_SUB() != nil:
		fg.emitSub(ltype)
		return ltype
	case ctx.OVERFLOW_MUL() != nil:
		fg.emitMul(ltype)
		return ltype
	case ctx.AMP() != nil:
		fg.body.I32And()
		return ltype
	case ctx.PIPE() != nil:
		fg.body.I32Or()
		return ltype
	case ctx.CARET() != nil:
		fg.body.I32Xor()
		return ltype
	case ctx.LSHIFT() != nil:
		fg.body.I32Shl()
		return ltype
	case ctx.RSHIFT() != nil:
		fg.body.I32ShrS()
		return ltype
	case ctx.EQ() != nil:
		fg.emitEq(ltype)
		return Bool
	case ctx.NEQ() != nil:
		fg.emitEq(ltype)
		fg.body.I32Eqz()
		return Bool
	case ctx.LT() != nil:
		fg.emitLt(ltype)
		return Bool
	case ctx.GT() != nil:
		fg.emitGt(ltype)
		return Bool
	case ctx.LTE() != nil:
		fg.emitLe(ltype)
		return Bool
	case ctx.GTE() != nil:
		fg.emitGe(ltype)
		return Bool
	case ctx.AND() != nil:
		fg.body.I32And()
		return Bool
	case ctx.OR() != nil:
		fg.body.I32Or()
		return Bool
	case ctx.IDENTITY_EQ() != nil:
		fg.body.I32Eq()
		return Bool
	case ctx.IDENTITY_NEQ() != nil:
		fg.body.I32Ne()
		return Bool
	}
	return ltype
}

func (fg *funcGen) emitNeg(t Type) {
	switch {
	case t != nil && (t.Kind() == KindInt64 || t.Kind() == KindUint64):
		fg.body.I64Const(-1)
		fg.body.I64Mul()
	case t != nil && t.Kind() == KindFloat:
		fg.body.F32Neg()
	case t != nil && t.Kind() == KindDouble:
		fg.body.F64Neg()
	default:
		fg.body.I32Const(-1)
		fg.body.I32Mul()
	}
}

func (fg *funcGen) emitEq(t Type) {
	if t != nil && (t.Kind() == KindInt64 || t.Kind() == KindUint64) {
		fg.body.I64Eq()
	} else {
		fg.body.I32Eq()
	}
}

func (fg *funcGen) emitLt(t Type) {
	if t != nil && t.Kind() == KindUint {
		fg.body.I32LtU()
	} else {
		fg.body.I32LtS()
	}
}

func (fg *funcGen) emitGt(t Type) {
	if t != nil && t.Kind() == KindUint {
		fg.body.I32GtU()
	} else {
		fg.body.I32GtS()
	}
}

func (fg *funcGen) emitLe(t Type) {
	if t != nil && t.Kind() == KindUint {
		fg.body.I32LeU()
	} else {
		fg.body.I32LeS()
	}
}

func (fg *funcGen) emitGe(t Type) {
	if t != nil && t.Kind() == KindUint {
		fg.body.I32GeU()
	} else {
		fg.body.I32GeS()
	}
}

// ── Ternary ───────────────────────────────────────────────────────────────────

func (fg *funcGen) genTernary(cond, then, els parser.IExprContext) Type {
	fg.genExpr(cond)
	fg.body.If(wasm.BlockVal{Val: wasm.I32})
	t := fg.genExpr(then)
	fg.body.Else()
	fg.genExpr(els)
	fg.body.End()
	return t
}

// ── Type cast ─────────────────────────────────────────────────────────────────

func (fg *funcGen) emitCast(ctx parser.IPrimitiveTypeContext, src Type) Type {
	target := resolvePrimitive(ctx)
	switch {
	case target.Kind() == KindFloat && src.Wasm() == WasmI32:
		fg.body.F32ConvertI32S()
	case target.Kind() == KindDouble && src.Wasm() == WasmI32:
		fg.body.F64ConvertI32S()
	case target.Kind() == KindInt && src.Wasm() == WasmF32:
		fg.body.I32TruncF32S()
	case target.Kind() == KindInt && src.Wasm() == WasmF64:
		fg.body.I32TruncF64S()
	case target.Kind() == KindInt64 && src.Wasm() == WasmI32:
		fg.body.I64ExtendI32S()
	case target.Kind() == KindInt && src.Wasm() == WasmI64:
		fg.body.I32WrapI64()
	}
	return target
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func identName(ctx parser.IExprContext) string {
	if p := ctx.Primary(); p != nil && p.ID() != nil {
		return p.ID().GetText()
	}
	return ""
}

var _ = strings.TrimSpace