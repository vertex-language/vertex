package compiler

import (
	"fmt"
	"strconv"

	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/wasm-compiler/wasm"
)

// ExprResult carries the result type of a compiled expression.
type ExprResult struct {
	Type      VType
	IsFuncRef bool // true = identifier resolved to a function; no value was pushed
}

var voidResult = ExprResult{Type: &PrimitiveType{Kind: KindVoid}}

// ── Top-level expression ──────────────────────────────────────────────────────

func (fc *FuncCompiler) compileExpression(ctx parser.IExpressionContext) (ExprResult, error) {
	if ctx == nil {
		return voidResult, nil
	}

	// ── Fast-path: detect assignment BEFORE emitting any LHS instructions ────
	// If we compiled the LHS first (local.get) and then tried to do local.set,
	// we'd leave a dangling value on the stack.
	if binExprs := ctx.BinaryExpressions(); binExprs != nil {
		all := binExprs.AllBinaryExpression()
		if len(all) == 1 {
			if assign, ok := all[0].(*parser.AssignmentExprContext); ok {
				return fc.compileAssignment(ctx.PrefixExpression(), assign)
			}
		}
	}

	// Normal expression: compile prefix then fold in binary ops via precedence climbing.
	lhs, err := fc.compilePrefixExpression(ctx.PrefixExpression())
	if err != nil {
		return voidResult, err
	}

	if binExprs := ctx.BinaryExpressions(); binExprs != nil {
		all := binExprs.AllBinaryExpression()
		idx := 0
		lhs, err = fc.compilePrecedenceClimbing(0, lhs, all, &idx)
		if err != nil {
			return voidResult, err
		}
	}
	return lhs, nil
}

// ── Assignment ────────────────────────────────────────────────────────────────

// compileAssignment handles `target = expr` without emitting a speculative
// local.get for the target.
func (fc *FuncCompiler) compileAssignment(
	lhsPfx parser.IPrefixExpressionContext,
	assign *parser.AssignmentExprContext,
) (ExprResult, error) {
	bare, ok := lhsPfx.(*parser.BarePostfixContext)
	if !ok {
		return voidResult, fmt.Errorf("unsupported assignment target (expected identifier)")
	}
	postfix := bare.PostfixExpression()

	// ── Simple local: name = expr ─────────────────────────────────────────────
	if len(postfix.AllPostfixSuffix()) == 0 {
		primary := postfix.PrimaryExpression()
		ident, ok := primary.(*parser.IdentifierExprContext)
		if !ok {
			return voidResult, fmt.Errorf("unsupported assignment target")
		}
		name := ident.Identifier().GetText()
		sym, found := fc.scope.lookup(name)
		if !found {
			return voidResult, fmt.Errorf("use of unresolved identifier %q", name)
		}
		if !sym.IsMut {
			return voidResult, fmt.Errorf("cannot assign to immutable binding %q", name)
		}

		// Compile RHS — its value lands on the stack.
		if _, err := fc.compileExpression(assign.Expression()); err != nil {
			return voidResult, err
		}
		fc.body.LocalSet(sym.WasmIdx)
		return voidResult, nil
	}

	// ── Struct field store: base.field = expr ─────────────────────────────────
	// Compile the base pointer, then do a memory store at the field offset.
	suffixes := postfix.AllPostfixSuffix()
	lastSuffix := suffixes[len(suffixes)-1]
	ems := lastSuffix.ExplicitMemberSuffix()
	if ems == nil {
		return voidResult, fmt.Errorf("unsupported compound assignment target")
	}
	fieldName := ""
	if id := ems.Identifier(); id != nil {
		fieldName = id.GetText()
	}

	// Compile the base expression (everything before the last member suffix).
	// For now we only support one level of indirection.
	baseRes, err := fc.compilePrimaryExpression(postfix.PrimaryExpression())
	if err != nil {
		return voidResult, err
	}
	// Apply intermediate suffixes (all but last).
	for _, s := range suffixes[:len(suffixes)-1] {
		baseRes, err = fc.compilePostfixSuffix(baseRes, s, postfix.PrimaryExpression())
		if err != nil {
			return voidResult, err
		}
	}

	st, isStruct := baseRes.Type.(*StructType)
	if !isStruct {
		return voidResult, fmt.Errorf("member assignment on non-struct type %v", baseRes.Type)
	}
	ftype, offset, found := st.FieldOffset(fieldName)
	if !found {
		return voidResult, fmt.Errorf("type %s has no field %q", st.Name, fieldName)
	}

	// Stack now: [basePtr]. Compile RHS so stack is [basePtr, value].
	if _, err := fc.compileExpression(assign.Expression()); err != nil {
		return voidResult, err
	}

	switch wt, _ := ftype.WasmType(); wt {
	case wasm.I32:
		fc.body.I32Store(2, offset)
	case wasm.I64:
		fc.body.I64Store(3, offset)
	case wasm.F32:
		fc.body.F32Store(2, offset)
	case wasm.F64:
		fc.body.F64Store(3, offset)
	default:
		fc.body.I32Store(2, offset)
	}
	return voidResult, nil
}

// ── Precedence Climbing (Operator Precedence) ─────────────────────────────────

// precedenceMap defines the binding power of binary operators.
// Higher numbers = tighter binding (evaluated first).
var precedenceMap = map[string]int{
	"*": 7, "/": 7, "%": 7,
	"+": 6, "-": 6,
	"<<": 5, ">>": 5,
	"<": 4, "<=": 4, ">": 4, ">=": 4,
	"==": 3, "!=": 3,
	"&": 2,
	"^": 1,
	"|": 0,
	"&&": -1,
	"||": -2,
}

func getOpInfo(ctx parser.IBinaryExpressionContext) (string, int, bool) {
	if b, ok := ctx.(*parser.BinaryOpContext); ok {
		op := b.BinaryOperator().GetText()
		if prec, found := precedenceMap[op]; found {
			return op, prec, true
		}
		return op, 0, true
	}
	return "", 0, false
}

// compilePrecedenceClimbing dynamically builds the AST structure by looking ahead
// at operator priority, emitting instructions in the correct mathematical order.
func (fc *FuncCompiler) compilePrecedenceClimbing(
	minPrec int,
	lhs ExprResult,
	bins []parser.IBinaryExpressionContext,
	idx *int,
) (ExprResult, error) {

	for *idx < len(bins) {
		binCtx := bins[*idx]

		opStr, prec, isBinOp := getOpInfo(binCtx)

		// If it's a binary operator but its precedence is lower than what we are
		// currently climbing, break out and let the previous caller handle it.
		if isBinOp && prec < minPrec {
			break
		}

		// Consume this operator
		*idx++

		switch b := binCtx.(type) {
		case *parser.BinaryOpContext:
			// Compile the immediate right-hand prefix expression
			rhs, err := fc.compilePrefixExpression(b.PrefixExpression())
			if err != nil {
				return voidResult, err
			}

			// Peek ahead: as long as the NEXT operator binds tighter, recursively climb
			for *idx < len(bins) {
				nextCtx := bins[*idx]
				_, nextPrec, nextIsBin := getOpInfo(nextCtx)
				if !nextIsBin {
					break // Not a standard binary op, handle below
				}
				if nextPrec > prec {
					rhs, err = fc.compilePrecedenceClimbing(prec+1, rhs, bins, idx)
					if err != nil {
						return voidResult, err
					}
				} else {
					break
				}
			}

			// Both sides are compiled and sitting on the stack, now emit the operation instruction
			lhs, err = fc.emitBinaryOp(opStr, lhs, rhs)
			if err != nil {
				return voidResult, err
			}

		case *parser.AssignmentExprContext:
			if _, err := fc.compileExpression(b.Expression()); err != nil {
				return voidResult, err
			}
			lhs = voidResult

		case *parser.TernaryExprContext:
			rhs, err := fc.compileExpression(b.Expression())
			if err != nil {
				return voidResult, err
			}
			_ = rhs

		case *parser.TypeCastExprContext:
			// no-op for now
		}
	}

	return lhs, nil
}

// ── Binary expression emitter ─────────────────────────────────────────────────

// emitBinaryOp emits the WASM instruction(s) for a binary operator.
// lhs value is already on the stack; this method assumes rhs is on stack and emits op.
func (fc *FuncCompiler) emitBinaryOp(
	op string, lhs, rhs ExprResult,
) (ExprResult, error) {
	resType := lhs.Type
	if resType == nil {
		resType = rhs.Type
	}

	pt, isPrim := resType.(*PrimitiveType)
	if !isPrim {
		return voidResult, fmt.Errorf("binary op %q not supported on type %v", op, resType)
	}

	switch op {
	case "+":
		switch pt.Kind {
		case KindInt, KindBool, KindUInt:
			fc.body.I32Add()
		case KindInt64:
			fc.body.I64Add()
		case KindFloat:
			fc.body.F32Add()
		case KindDouble:
			fc.body.F64Add()
		}
		return lhs, nil
	case "-":
		switch pt.Kind {
		case KindInt, KindBool, KindUInt:
			fc.body.I32Sub()
		case KindInt64:
			fc.body.I64Sub()
		case KindFloat:
			fc.body.F32Sub()
		case KindDouble:
			fc.body.F64Sub()
		}
		return lhs, nil
	case "*":
		switch pt.Kind {
		case KindInt, KindBool, KindUInt:
			fc.body.I32Mul()
		case KindInt64:
			fc.body.I64Mul()
		case KindFloat:
			fc.body.F32Mul()
		case KindDouble:
			fc.body.F64Mul()
		}
		return lhs, nil
	case "/":
		switch pt.Kind {
		case KindInt:
			fc.body.I32DivS()
		case KindUInt:
			fc.body.I32DivU()
		case KindInt64:
			fc.body.I64DivS()
		case KindFloat:
			fc.body.F32Div()
		case KindDouble:
			fc.body.F64Div()
		}
		return lhs, nil
	case "%":
		switch pt.Kind {
		case KindInt:
			fc.body.I32RemS()
		case KindUInt:
			fc.body.I32RemU()
		case KindInt64:
			fc.body.I64RemS()
		}
		return lhs, nil
	case "&":
		fc.body.I32And()
		return lhs, nil
	case "|":
		fc.body.I32Or()
		return lhs, nil
	case "^":
		fc.body.I32Xor()
		return lhs, nil
	case "<<":
		fc.body.I32Shl()
		return lhs, nil
	case ">>":
		fc.body.I32ShrS()
		return lhs, nil
	case "==":
		switch pt.Kind {
		case KindFloat:
			fc.body.F32Eq()
		case KindDouble:
			fc.body.F64Eq()
		default:
			fc.body.I32Eq()
		}
		return ExprResult{Type: &PrimitiveType{Kind: KindBool}}, nil
	case "!=":
		switch pt.Kind {
		case KindFloat:
			fc.body.F32Ne()
		case KindDouble:
			fc.body.F64Ne()
		default:
			fc.body.I32Ne()
		}
		return ExprResult{Type: &PrimitiveType{Kind: KindBool}}, nil
	case "<":
		switch pt.Kind {
		case KindFloat:
			fc.body.F32Lt()
		case KindDouble:
			fc.body.F64Lt()
		case KindUInt:
			fc.body.I32LtU()
		default:
			fc.body.I32LtS()
		}
		return ExprResult{Type: &PrimitiveType{Kind: KindBool}}, nil
	case ">":
		switch pt.Kind {
		case KindFloat:
			fc.body.F32Gt()
		case KindDouble:
			fc.body.F64Gt()
		case KindUInt:
			fc.body.I32GtU()
		default:
			fc.body.I32GtS()
		}
		return ExprResult{Type: &PrimitiveType{Kind: KindBool}}, nil
	case "<=":
		switch pt.Kind {
		case KindFloat:
			fc.body.F32Le()
		case KindDouble:
			fc.body.F64Le()
		case KindUInt:
			fc.body.I32LeU()
		default:
			fc.body.I32LeS()
		}
		return ExprResult{Type: &PrimitiveType{Kind: KindBool}}, nil
	case ">=":
		switch pt.Kind {
		case KindFloat:
			fc.body.F32Ge()
		case KindDouble:
			fc.body.F64Ge()
		case KindUInt:
			fc.body.I32GeU()
		default:
			fc.body.I32GeS()
		}
		return ExprResult{Type: &PrimitiveType{Kind: KindBool}}, nil
	case "&&":
		fc.body.I32And()
		return ExprResult{Type: &PrimitiveType{Kind: KindBool}}, nil
	case "||":
		fc.body.I32Or()
		return ExprResult{Type: &PrimitiveType{Kind: KindBool}}, nil
	}

	return voidResult, fmt.Errorf("unsupported binary operator %q", op)
}

// ── Prefix expression ─────────────────────────────────────────────────────────

func (fc *FuncCompiler) compilePrefixExpression(
	ctx parser.IPrefixExpressionContext,
) (ExprResult, error) {
	switch p := ctx.(type) {
	case *parser.BarePostfixContext:
		return fc.compilePostfixExpression(p.PostfixExpression())
	case *parser.PrefixOpContext:
		res, err := fc.compilePostfixExpression(p.PostfixExpression())
		if err != nil {
			return voidResult, err
		}
		return fc.emitPrefixOp(p.PrefixOperator().GetText(), res)
	case *parser.InoutContext:
		name := p.InOutExpression().Identifier().GetText()
		sym, ok := fc.scope.lookup(name)
		if !ok {
			return voidResult, fmt.Errorf("use of unresolved identifier %q", name)
		}
		fc.body.LocalGet(sym.WasmIdx)
		return ExprResult{Type: sym.Type}, nil
	}
	return voidResult, fmt.Errorf("unsupported prefix expression %T", ctx)
}

func (fc *FuncCompiler) emitPrefixOp(op string, res ExprResult) (ExprResult, error) {
	pt, isPrim := res.Type.(*PrimitiveType)
	if !isPrim {
		return res, nil
	}
	switch op {
	case "-":
		switch pt.Kind {
		case KindInt, KindUInt:
			tmp := fc.allocLocal(wasm.I32)
			fc.body.LocalSet(tmp)
			fc.body.I32Const(0)
			fc.body.LocalGet(tmp)
			fc.body.I32Sub()
		case KindFloat:
			fc.body.F32Neg()
		case KindDouble:
			fc.body.F64Neg()
		}
	case "!":
		fc.body.I32Eqz()
	}
	return res, nil
}

// ── Postfix expression ────────────────────────────────────────────────────────

func (fc *FuncCompiler) compilePostfixExpression(
	ctx parser.IPostfixExpressionContext,
) (ExprResult, error) {
	res, err := fc.compilePrimaryExpression(ctx.PrimaryExpression())
	if err != nil {
		return voidResult, err
	}
	for _, suffix := range ctx.AllPostfixSuffix() {
		res, err = fc.compilePostfixSuffix(res, suffix, ctx.PrimaryExpression())
		if err != nil {
			return voidResult, err
		}
	}
	return res, nil
}

func (fc *FuncCompiler) compilePostfixSuffix(
	base ExprResult,
	suffix parser.IPostfixSuffixContext,
	primaryCtx parser.IPrimaryExpressionContext,
) (ExprResult, error) {
	if fcs := suffix.FunctionCallSuffix(); fcs != nil {
		return fc.compileFunctionCallSuffix(base, fcs, primaryCtx)
	}
	if ems := suffix.ExplicitMemberSuffix(); ems != nil {
		return fc.compileMemberAccess(base, ems)
	}
	if suffix.OptionalChainingLiteral() != nil || suffix.ForcedValueSuffix() != nil {
		return base, nil
	}
	if ss := suffix.SubscriptSuffix(); ss != nil {
		return fc.compileSubscript(base, ss)
	}
	return base, nil
}

// ── Primary expression ────────────────────────────────────────────────────────

func (fc *FuncCompiler) compilePrimaryExpression(
	ctx parser.IPrimaryExpressionContext,
) (ExprResult, error) {
	switch p := ctx.(type) {
	case *parser.LitExprContext:
		return fc.compileLiteralExpression(p.LiteralExpression())
	case *parser.IdentifierExprContext:
		return fc.compileIdentifier(p)
	case *parser.ParenExprContext:
		return fc.compileExpression(p.ParenthesizedExpression().Expression())
	case *parser.TupleExprContext:
		te := p.TupleExpression()
		if tl := te.TupleElementList(); tl != nil {
			if all := tl.AllTupleElement(); len(all) > 0 {
				return fc.compileExpression(all[0].Expression())
			}
		}
		return voidResult, nil
	case *parser.SelfExprContext:
		fc.body.I32Const(0)
		return ExprResult{Type: &PrimitiveType{Kind: KindInt}}, nil
	case *parser.WildcardExprContext:
		return voidResult, nil
	case *parser.ClosureExprContext:
		fc.body.I32Const(0)
		return ExprResult{Type: &PrimitiveType{Kind: KindInt}}, nil
	case *parser.MacroExprContext:
		return fc.compileMacroExpansion(p.MacroExpansionExpression())
	}
	return voidResult, fmt.Errorf("unsupported primary expression %T", ctx)
}

// ── Identifier ────────────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileIdentifier(p *parser.IdentifierExprContext) (ExprResult, error) {
	name := p.Identifier().GetText()

	// Local variable — push its value.
	if sym, ok := fc.scope.lookup(name); ok {
		fc.body.LocalGet(sym.WasmIdx)
		return ExprResult{Type: sym.Type}, nil
	}

	// Known function — mark as funcref so the call suffix skips the Drop.
	// We do NOT push anything; Call(idx) is emitted by compileFunctionCallSuffix.
	if info, ok := fc.c.funcMap[name]; ok {
		return ExprResult{Type: info.Ret, IsFuncRef: true}, nil
	}

	return voidResult, fmt.Errorf("use of unresolved identifier %q", name)
}

// ── Function calls ────────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileFunctionCallSuffix(
	base ExprResult,
	fcs parser.IFunctionCallSuffixContext,
	primaryCtx parser.IPrimaryExpressionContext,
) (ExprResult, error) {
	calleeName := ""
	if id, ok := primaryCtx.(*parser.IdentifierExprContext); ok {
		calleeName = id.Identifier().GetText()
	}

	info, ok := fc.c.funcMap[calleeName]
	if !ok {
		return voidResult, fmt.Errorf("call to unknown function %q", calleeName)
	}

	// Only Drop if the identifier pushed a speculative value (old path).
	// With IsFuncRef, nothing was pushed, so no Drop needed.
	if _, isIdent := primaryCtx.(*parser.IdentifierExprContext); isIdent && !base.IsFuncRef {
		fc.body.Drop()
	}

	// Compile arguments.
	if fca := fcs.FunctionCallArgumentClause(); fca != nil {
		if al := fca.FunctionCallArgumentList(); al != nil {
			for _, arg := range al.AllFunctionCallArgument() {
				if _, err := fc.compileExpression(arg.Expression()); err != nil {
					return voidResult, err
				}
			}
		}
	}

	fc.body.Call(info.WasmIdx)
	return ExprResult{Type: info.Ret}, nil
}

// ── Member access ─────────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileMemberAccess(
	base ExprResult,
	ems parser.IExplicitMemberSuffixContext,
) (ExprResult, error) {
	fieldName := ""
	if id := ems.Identifier(); id != nil {
		fieldName = id.GetText()
	}
	st, isStruct := base.Type.(*StructType)
	if !isStruct {
		return voidResult, fmt.Errorf("member access on non-struct type %v", base.Type)
	}
	ftype, offset, found := st.FieldOffset(fieldName)
	if !found {
		return voidResult, fmt.Errorf("type %s has no field %q", st.Name, fieldName)
	}
	switch wt, _ := ftype.WasmType(); wt {
	case wasm.I32:
		fc.body.I32Load(2, offset)
	case wasm.I64:
		fc.body.I64Load(3, offset)
	case wasm.F32:
		fc.body.F32Load(2, offset)
	case wasm.F64:
		fc.body.F64Load(3, offset)
	default:
		fc.body.I32Load(2, offset)
	}
	return ExprResult{Type: ftype}, nil
}

// ── Subscript ─────────────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileSubscript(
	base ExprResult,
	ss parser.ISubscriptSuffixContext,
) (ExprResult, error) {
	al := ss.FunctionCallArgumentList()
	if al == nil || len(al.AllFunctionCallArgument()) == 0 {
		return base, nil
	}
	if _, err := fc.compileExpression(al.AllFunctionCallArgument()[0].Expression()); err != nil {
		return voidResult, err
	}
	fc.body.I32Const(4)
	fc.body.I32Mul()
	fc.body.I32Add()
	fc.body.I32Load(2, 0)
	return ExprResult{Type: &PrimitiveType{Kind: KindInt}}, nil
}

// ── Literals ──────────────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileLiteralExpression(
	ctx parser.ILiteralExpressionContext,
) (ExprResult, error) {
	if lit := ctx.Literal(); lit != nil {
		return fc.compileLiteral(lit)
	}
	if ctx.ArrayLiteral() != nil {
		fc.body.I32Const(0)
		return ExprResult{Type: &PrimitiveType{Kind: KindInt}}, nil
	}
	if ctx.DictionaryLiteral() != nil {
		fc.body.I32Const(0)
		return ExprResult{Type: &PrimitiveType{Kind: KindInt}}, nil
	}
	return voidResult, fmt.Errorf("unsupported literal expression")
}

func (fc *FuncCompiler) compileLiteral(ctx parser.ILiteralContext) (ExprResult, error) {
	if num := ctx.NumericLiteral(); num != nil {
		return fc.compileNumericLiteral(num)
	}
	if bl := ctx.BooleanLiteral(); bl != nil {
		if bl.TRUE() != nil {
			fc.body.I32Const(1)
		} else {
			fc.body.I32Const(0)
		}
		return ExprResult{Type: &PrimitiveType{Kind: KindBool}}, nil
	}
	if ctx.NIL() != nil {
		fc.body.I32Const(0)
		return ExprResult{Type: &PrimitiveType{Kind: KindInt}}, nil
	}
	if s := ctx.STRING_LITERAL(); s != nil {
		offset := fc.c.internString(s.GetText())
		fc.body.I32Const(int32(offset))
		return ExprResult{Type: &PrimitiveType{Kind: KindString}}, nil
	}
	if s := ctx.MULTILINE_STRING_LITERAL(); s != nil {
		offset := fc.c.internString(s.GetText())
		fc.body.I32Const(int32(offset))
		return ExprResult{Type: &PrimitiveType{Kind: KindString}}, nil
	}
	return voidResult, fmt.Errorf("unsupported literal")
}

func (fc *FuncCompiler) compileNumericLiteral(
	ctx parser.INumericLiteralContext,
) (ExprResult, error) {
	text := ctx.GetText()
	sign := int32(1)
	if len(text) > 0 && text[0] == '-' {
		sign = -1
		text = text[1:]
	} else if len(text) > 0 && text[0] == '+' {
		text = text[1:]
	}

	if ctx.FLOAT_LITERAL() != nil {
		f, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return voidResult, fmt.Errorf("invalid float literal %q: %w", text, err)
		}
		fc.body.F64Const(f * float64(sign))
		return ExprResult{Type: &PrimitiveType{Kind: KindDouble}}, nil
	}

	base := 10
	if len(text) >= 2 {
		switch text[1] {
		case 'x', 'X':
			base, text = 16, text[2:]
		case 'o', 'O':
			base, text = 8, text[2:]
		case 'b', 'B':
			base, text = 2, text[2:]
		}
	}
	cleaned := ""
	for _, ch := range text {
		if ch != '_' {
			cleaned += string(ch)
		}
	}
	v, err := strconv.ParseInt(cleaned, base, 64)
	if err != nil {
		uv, err2 := strconv.ParseUint(cleaned, base, 64)
		if err2 != nil {
			return voidResult, fmt.Errorf("invalid integer literal %q: %w", ctx.GetText(), err)
		}
		v = int64(uv)
	}
	v *= int64(sign)
	fc.body.I32Const(int32(v))
	return ExprResult{Type: &PrimitiveType{Kind: KindInt}}, nil
}

// ── Macro expansion ───────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileMacroExpansion(
	ctx parser.IMacroExpansionExpressionContext,
) (ExprResult, error) {
	name := ctx.Identifier().GetText()
	switch name {
	case "assert":
		if fca := ctx.FunctionCallArgumentClause(); fca != nil {
			if al := fca.FunctionCallArgumentList(); al != nil {
				if all := al.AllFunctionCallArgument(); len(all) > 0 {
					if _, err := fc.compileExpression(all[0].Expression()); err != nil {
						return voidResult, err
					}
					fc.body.I32Eqz()
					fc.body.If(wasm.BlockEmpty{})
					fc.body.Unreachable()
					fc.body.End()
				}
			}
		}
		return voidResult, nil
	}
	fc.body.I32Const(0)
	return ExprResult{Type: &PrimitiveType{Kind: KindInt}}, nil
}