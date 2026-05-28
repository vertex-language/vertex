// codegen_stmt.go
package compiler

import (
	"strconv"

	"github.com/vertex-language/compiler/wasm"
	"github.com/vertex-language/vertex/parser"
)

type funcGen struct {
	cg     *CodeGen
	sig    *FuncSig
	scope  *Scope
	body   *wasm.FunctionBody

	localCount  uint32
	localTypes  []wasm.ValType
	paramCount  uint32

	frameAllocs map[string]int32
	loopDepth   int
}

func (fg *funcGen) newLocal(t Type) uint32 {
	idx := fg.localCount
	fg.localCount++
	vt := ToWasmVal(t.Wasm())
	if t.Wasm() == WasmNone {
		vt = wasm.I32
	}
	fg.localTypes = append(fg.localTypes, vt)
	return idx
}

func (fg *funcGen) declareLocals(body *wasm.FunctionBody) {
	if fg.localCount <= fg.paramCount {
		return
	}
	currentType := fg.localTypes[fg.paramCount]
	var count uint32

	for i := fg.paramCount; i < fg.localCount; i++ {
		vt := fg.localTypes[i]
		if vt == currentType {
			count++
		} else {
			body.AddLocals(count, currentType)
			currentType = vt
			count = 1
		}
	}
	if count > 0 {
		body.AddLocals(count, currentType)
	}
}

// ── Struct-literal helpers ────────────────────────────────────────────────────

func structLiteralOf(expr parser.IExprContext) parser.IStructLiteralExprContext {
	if expr == nil {
		return nil
	}
	p := expr.Primary()
	if p == nil {
		return nil
	}
	return p.StructLiteralExpr()
}

func inferStructType(declared Type, expr parser.IExprContext, scope *Scope) Type {
	if declared != nil && (declared.Kind() == KindStruct || declared.Kind() == KindClass) {
		return declared
	}
	if sl := structLiteralOf(expr); sl != nil {
		name := sl.ID().GetText()
		if sym := scope.Lookup(name); sym != nil && sym.Kind == SymType {
			return sym.Type
		}
	}
	return nil
}

// ── Pre-scan ──────────────────────────────────────────────────────────────────

func (fg *funcGen) preScanBlock(ctx parser.IBlockContext) {
	if ctx == nil {
		return
	}
	for _, stmt := range ctx.AllStmt() {
		if vd := stmt.VarDeclStmt(); vd != nil {
			fg.preScanVarDecl(vd)
		}
		if is := stmt.IfStmt(); is != nil {
			fg.preScanBlock(is.Block())
			for _, eic := range is.AllElseIfClause() {
				fg.preScanBlock(eic.Block())
			}
			if ec := is.ElseClause(); ec != nil {
				fg.preScanBlock(ec.Block())
			}
		}
		if ws := stmt.WhileStmt(); ws != nil {
			fg.preScanBlock(ws.Block())
		}
		if fi := stmt.ForInStmt(); fi != nil {
			fg.preScanBlock(fi.Block())
		}
		if sw := stmt.SwitchStmt(); sw != nil {
			for _, sc := range sw.AllSwitchCase() {
				for _, s := range sc.AllStmt() {
					if vd := s.VarDeclStmt(); vd != nil {
						fg.preScanVarDecl(vd)
					}
				}
			}
			if dc := sw.DefaultCase(); dc != nil {
				for _, s := range dc.AllStmt() {
					if vd := s.VarDeclStmt(); vd != nil {
						fg.preScanVarDecl(vd)
					}
				}
			}
		}
	}
}

func (fg *funcGen) preScanVarDecl(ctx parser.IVarDeclStmtContext) {
	if ctx.ID() == nil {
		return
	}
	name := ctx.ID().GetText()

	var declared Type
	if ctx.Type_() != nil {
		declared = ResolveType(ctx.Type_(), fg.cg.scope)
	}

	t := inferStructType(declared, ctx.Expr(), fg.cg.scope)
	if t != nil {
		size := SizeOf(t)
		if size <= 0 {
			return
		}
		fg.frameAllocs[name] = fg.cg.allocFrame(size)
		return
	}

	if ac := arrayConstructOf(ctx.Expr()); ac != nil {
		size := resolveArrayConstructSize(ac)
		if size > 0 {
			fg.frameAllocs[name] = fg.cg.allocFrame(size)
		}
	}
}

// ── Block ─────────────────────────────────────────────────────────────────────

func (fg *funcGen) genBlock(ctx parser.IBlockContext) {
	for _, stmt := range ctx.AllStmt() {
		fg.genStmt(stmt)
	}
}

// ── Statement dispatch ────────────────────────────────────────────────────────

func (fg *funcGen) genStmt(ctx parser.IStmtContext) {
	switch {
	case ctx.VarDeclStmt() != nil:
		fg.genVarDecl(ctx.VarDeclStmt())
	case ctx.AssignStmt() != nil:
		fg.genAssign(ctx.AssignStmt())
	case ctx.CompoundAssignStmt() != nil:
		fg.genCompoundAssign(ctx.CompoundAssignStmt())
	case ctx.IfStmt() != nil:
		fg.genIf(ctx.IfStmt())
	case ctx.WhileStmt() != nil:
		fg.genWhile(ctx.WhileStmt())
	case ctx.ForInStmt() != nil:
		fg.genForIn(ctx.ForInStmt())
	case ctx.SwitchStmt() != nil:
		fg.genSwitch(ctx.SwitchStmt())
	case ctx.ReturnStmt() != nil:
		fg.genReturn(ctx.ReturnStmt())
	case ctx.DeferStmt() != nil:
		fg.genDefer(ctx.DeferStmt())
	case ctx.BreakStmt() != nil:
		fg.body.Br(1)
	case ctx.ContinueStmt() != nil:
		fg.body.Br(0)
	case ctx.ExprStmt() != nil:
		t := fg.genExpr(ctx.ExprStmt().Expr())
		if t != nil && t.Kind() != KindVoid {
			fg.body.Drop()
		}
	}
}

// ── var / let ─────────────────────────────────────────────────────────────────

func (fg *funcGen) genVarDecl(ctx parser.IVarDeclStmtContext) {
	if ctx.ID() == nil {
		return
	}
	name := ctx.ID().GetText()
	mutable := ctx.BindingKw() != nil && ctx.BindingKw().VAR() != nil

	var declType Type
	if ctx.Type_() != nil {
		declType = ResolveType(ctx.Type_(), fg.cg.scope)
	}

	if _, preAlloced := fg.frameAllocs[name]; !preAlloced {
		t := inferStructType(declType, ctx.Expr(), fg.cg.scope)
		if t != nil {
			size := SizeOf(t)
			if size > 0 {
				fg.frameAllocs[name] = fg.cg.allocFrame(size)
				if declType == nil || declType.Kind() == KindVoid {
					declType = t
				}
			}
		}
	}

	// ── fixed-size array in frame ─────────────────────────────────────────
	if addr, ok := fg.frameAllocs[name]; ok {
		if ac := arrayConstructOf(ctx.Expr()); ac != nil {
			size := resolveArrayConstructSize(ac)
			elemType := ResolveType(ac.Type_(), fg.cg.scope)

			// zero-fill word by word (correct for both short and long form)
			for off := 0; off < size; off += 4 {
				fg.body.I32Const(addr + int32(off))
				fg.body.I32Const(0)
				fg.body.I32Store(2, 0)
			}

			ptrLocal := fg.newLocal(Int)
			fg.body.I32Const(addr)
			fg.body.LocalSet(ptrLocal)

			fg.scope.Define(&Symbol{
				Name:     name,
				Kind:     SymVar,
				Type:     &ArrayType{Elem: elemType},
				LocalIdx: ptrLocal,
				Mutable:  mutable,
				IsFrame:  true,
				FrameOff: addr,
			})
			return
		}
	}

	// ── struct / class variable stored in linear memory ───────────────────
	if addr, ok := fg.frameAllocs[name]; ok {
		if declType == nil || declType.Kind() == KindVoid {
			declType = inferStructType(nil, ctx.Expr(), fg.cg.scope)
		}

		ptrLocal := fg.newLocal(Int)
		fg.body.I32Const(addr)
		fg.body.LocalSet(ptrLocal)

		if ctx.Expr() != nil {
			if sl := structLiteralOf(ctx.Expr()); sl != nil {
				fg.genStructLitFields(addr, sl, declType)
			} else {
				srcType := fg.genExpr(ctx.Expr())
				fg.emitStructMemcopy(addr, srcType)
			}
		}

		fg.scope.Define(&Symbol{
			Name:     name,
			Kind:     SymVar,
			Type:     declType,
			LocalIdx: ptrLocal,
			Mutable:  mutable,
			IsFrame:  true,
			FrameOff: addr,
		})
		return
	}

	// ── regular scalar local ──────────────────────────────────────────────
	valType := Type(Int)
	if declType != nil && declType.Kind() != KindVoid {
		valType = declType
	}

	localIdx := fg.newLocal(valType)
	fg.scope.Define(&Symbol{
		Name:     name,
		Kind:     SymVar,
		Type:     valType,
		LocalIdx: localIdx,
		Mutable:  mutable,
	})

	if ctx.Expr() != nil {
		fg.genExpr(ctx.Expr())
		fg.body.LocalSet(localIdx)
	}
}

func arrayConstructOf(expr parser.IExprContext) parser.IArrayConstructExprContext {
	if expr == nil {
		return nil
	}
	p := expr.Primary()
	if p == nil {
		return nil
	}
	return p.ArrayConstructExpr()
}

// resolveArrayConstructSize returns the total byte size of the array.
//
// Handles all four construction forms (§22.1, §22.2):
//
//	[T]()                        — empty growable, no frame allocation
//	[T](n)                       — short form: fixed, n elements, zero-filled
//	[T](capacity: n)             — growable with hint, no frame allocation
//	[T](repeating: v, count: n)  — long form: fixed, n elements
func resolveArrayConstructSize(ac parser.IArrayConstructExprContext) int {
	elemType := ResolveType(ac.Type_(), nil)
	elemSize := SizeOf(elemType)
	if elemSize <= 0 {
		elemSize = 1
	}
	ids := ac.AllID()
	exprs := ac.AllExpr()

	// ← 2.0: short form [T](n) — single unlabelled count expression,
	// always zero-fills (no repeating value).
	if len(ids) == 0 && len(exprs) == 1 {
		if n := literalIntExpr(exprs[0]); n > 0 {
			return elemSize * n
		}
	}

	// Long form [T](repeating: v, count: n) — find the "count" label.
	for i, id := range ids {
		if id.GetText() == "count" && i < len(exprs) {
			if n := literalIntExpr(exprs[i]); n > 0 {
				return elemSize * n
			}
		}
	}
	return 0
}

func literalIntExpr(expr parser.IExprContext) int {
	if expr == nil {
		return 0
	}
	p := expr.Primary()
	if p == nil || p.Literal() == nil {
		return 0
	}
	lit := p.Literal()
	if lit.DEC_INT_LIT() != nil {
		v, err := strconv.Atoi(lit.DEC_INT_LIT().GetText())
		if err == nil {
			return v
		}
	}
	return 0
}

func (fg *funcGen) genStructLitFields(addr int32, ctx parser.IStructLiteralExprContext, t Type) {
	st, ok := t.(*StructType)
	if !ok {
		if ct, ok2 := t.(*ClassType); ok2 {
			fg.genClassLitFields(addr, ctx, ct)
		}
		return
	}

	for _, init := range ctx.AllStructFieldInit() {
		fieldName := init.ID().GetText()
		f := st.Field(fieldName)
		if f == nil {
			continue
		}
		if f.Type.Kind() == KindStruct {
			if nestedSL := structLiteralOf(init.Expr()); nestedSL != nil {
				fg.genStructLitFields(addr+int32(f.Offset), nestedSL, f.Type)
				continue
			}
		}
		fg.body.I32Const(addr)
		fg.genExpr(init.Expr())
		fg.emitStore(f.Type, uint32(f.Offset))
	}
}

func (fg *funcGen) genClassLitFields(addr int32, ctx parser.IStructLiteralExprContext, ct *ClassType) {
	for _, init := range ctx.AllStructFieldInit() {
		fieldName := init.ID().GetText()
		var f *StructField
		for _, sf := range ct.Fields {
			if sf.Name == fieldName {
				f = sf
				break
			}
		}
		if f == nil {
			continue
		}
		fg.body.I32Const(addr)
		fg.genExpr(init.Expr())
		fg.emitStore(f.Type, uint32(f.Offset))
	}
}

// ── Assignment ────────────────────────────────────────────────────────────────

func (fg *funcGen) genAssign(ctx parser.IAssignStmtContext) {
	lv := ctx.Lvalue()

	if lv.DOT() != nil || lv.LBRACKET() != nil {
		fieldType := fg.genLvalueAddr(lv)
		fg.genExpr(ctx.Expr())
		fg.emitStore(fieldType, 0)
		return
	}

	name := lv.ID().GetText()
	sym := fg.scope.Lookup(name)
	if sym == nil {
		return
	}

	if sym.IsFrame {
		if sl := structLiteralOf(ctx.Expr()); sl != nil {
			fg.genStructLitFields(sym.FrameOff, sl, sym.Type)
		} else {
			srcType := fg.genExpr(ctx.Expr())
			fg.emitStructMemcopy(sym.FrameOff, srcType)
		}
		return
	}

	fg.genExpr(ctx.Expr())
	fg.body.LocalSet(sym.LocalIdx)
}

func lvalueName(ctx parser.ILvalueContext) string {
	if ctx.DOT() == nil && ctx.LBRACKET() == nil && ctx.ID() != nil {
		return ctx.ID().GetText()
	}
	return ""
}

// ── Compound assignment ───────────────────────────────────────────────────────

func (fg *funcGen) genCompoundAssign(ctx parser.ICompoundAssignStmtContext) {
	lv := ctx.Lvalue()

	if lv.DOT() != nil || lv.LBRACKET() != nil {
		fieldType := fg.genLvalueAddr(lv)
		addrLocal := fg.newLocal(Int)
		fg.body.LocalSet(addrLocal)

		fg.body.LocalGet(addrLocal)
		fg.emitLoad(fieldType, 0)
		fg.genExpr(ctx.Expr())
		fg.applyCompoundOp(ctx.CompoundOp(), fieldType)

		valLocal := fg.newLocal(fieldType)
		fg.body.LocalSet(valLocal)
		fg.body.LocalGet(addrLocal)
		fg.body.LocalGet(valLocal)
		fg.emitStore(fieldType, 0)
		return
	}

	name := lvalueName(lv)
	sym := fg.scope.Lookup(name)
	if sym == nil {
		return
	}
	fg.body.LocalGet(sym.LocalIdx)
	fg.genExpr(ctx.Expr())
	fg.applyCompoundOp(ctx.CompoundOp(), sym.Type)
	fg.body.LocalSet(sym.LocalIdx)
}

func (fg *funcGen) applyCompoundOp(op parser.ICompoundOpContext, t Type) {
	switch {
	case op.PLUS_ASSIGN() != nil:
		fg.emitAdd(t)
	case op.MINUS_ASSIGN() != nil:
		fg.emitSub(t)
	case op.STAR_ASSIGN() != nil:
		fg.emitMul(t)
	case op.SLASH_ASSIGN() != nil:
		fg.emitDiv(t)
	case op.MOD_ASSIGN() != nil:
		fg.emitRem(t)
	}
}

// ── Lvalue address resolution ─────────────────────────────────────────────────

func (fg *funcGen) genLvalueAddr(ctx parser.ILvalueContext) Type {
	if ctx.DOT() != nil {
		baseType := fg.genLvalueAddr(ctx.Lvalue())
		return fg.resolveFieldAddr(baseType, ctx.ID().GetText())
	}
	if ctx.LBRACKET() != nil {
		baseType := fg.genLvalueAddr(ctx.Lvalue())
		elemType := Void
		if at, ok := baseType.(*ArrayType); ok {
			elemType = at.Elem
		}
		fg.genExpr(ctx.Expr())
		fg.body.I32Const(int32(SizeOf(elemType)))
		fg.body.I32Mul()
		fg.body.I32Add()
		return elemType
	}
	name := ctx.ID().GetText()
	sym := fg.scope.Lookup(name)
	if sym == nil {
		fg.body.I32Const(0)
		return Void
	}
	if sym.IsFrame {
		fg.body.I32Const(sym.FrameOff)
		return sym.Type
	}
	fg.body.LocalGet(sym.LocalIdx)
	return sym.Type
}

func (fg *funcGen) resolveFieldAddr(base Type, fieldName string) Type {
	var fields []*StructField
	switch v := base.(type) {
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
			return f.Type
		}
	}
	fg.body.Drop()
	fg.body.I32Const(0)
	return Void
}

// ── Load / store helpers ──────────────────────────────────────────────────────

func (fg *funcGen) emitLoad(t Type, offset uint32) {
	switch t.Wasm() {
	case WasmI64:
		fg.body.I64Load(3, offset)
	case WasmF32:
		fg.body.F32Load(2, offset)
	case WasmF64:
		fg.body.F64Load(3, offset)
	default:
		switch SizeOf(t) {
		case 1:
			fg.body.I32Load8U(0, offset)
		case 2:
			fg.body.I32Load16U(1, offset)
		default:
			fg.body.I32Load(2, offset)
		}
	}
}

func (fg *funcGen) emitStore(t Type, offset uint32) {
	switch t.Wasm() {
	case WasmI64:
		fg.body.I64Store(3, offset)
	case WasmF32:
		fg.body.F32Store(2, offset)
	case WasmF64:
		fg.body.F64Store(3, offset)
	default:
		switch SizeOf(t) {
		case 1:
			fg.body.I32Store8(0, offset)
		case 2:
			fg.body.I32Store16(1, offset)
		default:
			fg.body.I32Store(2, offset)
		}
	}
}

func (fg *funcGen) emitStructMemcopy(dstAddr int32, srcType Type) {
	size := SizeOf(srcType)
	if size <= 0 {
		fg.body.Drop()
		return
	}
	srcLocal := fg.newLocal(Int)
	fg.body.LocalSet(srcLocal)

	for off := 0; off < size; off += 4 {
		remaining := size - off
		fg.body.I32Const(dstAddr + int32(off))
		fg.body.LocalGet(srcLocal)
		switch {
		case remaining >= 4:
			fg.body.I32Load(2, uint32(off))
			fg.body.I32Store(2, 0)
		case remaining >= 2:
			fg.body.I32Load16U(1, uint32(off))
			fg.body.I32Store16(1, 0)
		default:
			fg.body.I32Load8U(0, uint32(off))
			fg.body.I32Store8(0, 0)
		}
	}
}

// ── If / else ─────────────────────────────────────────────────────────────────

func (fg *funcGen) genIf(ctx parser.IIfStmtContext) {
	fg.genIfCond(ctx.IfCondition())
	fg.body.If(wasm.BlockEmpty{})
	fg.genBlock(ctx.Block())

	for _, eic := range ctx.AllElseIfClause() {
		fg.body.Else()
		fg.genIfCond(eic.IfCondition())
		fg.body.If(wasm.BlockEmpty{})
		fg.genBlock(eic.Block())
		fg.body.End()
	}

	if ec := ctx.ElseClause(); ec != nil {
		fg.body.Else()
		fg.genBlock(ec.Block())
	}

	fg.body.End()
}

func (fg *funcGen) genIfCond(ctx parser.IIfConditionContext) {
	if ctx == nil {
		return
	}
	fg.genExpr(ctx.Expr())
}

// ── While ─────────────────────────────────────────────────────────────────────

func (fg *funcGen) genWhile(ctx parser.IWhileStmtContext) {
	fg.body.Block(wasm.BlockEmpty{})
	fg.body.Loop(wasm.BlockEmpty{})
	fg.loopDepth++

	fg.genExpr(ctx.Expr())
	fg.body.I32Eqz()
	fg.body.BrIf(1)

	fg.genBlock(ctx.Block())
	fg.body.Br(0)

	fg.loopDepth--
	fg.body.End()
	fg.body.End()
}

// ── For-in ────────────────────────────────────────────────────────────────────

func (fg *funcGen) genForIn(ctx parser.IForInStmtContext) {
	iterName := ctx.ID().GetText()

	loopLocal := fg.newLocal(Int)
	counterLocal := fg.newLocal(Int)

	iterExpr := ctx.Expr()
	start, limit := fg.extractRange(iterExpr)
	if start != nil && limit != nil {
		fg.genExpr(start)
		fg.body.LocalSet(counterLocal)
		fg.genExpr(limit)
		fg.body.LocalSet(loopLocal)
	} else {
		fg.body.I32Const(0)
		fg.body.LocalSet(counterLocal)
		fg.genExpr(iterExpr)
		fg.body.LocalSet(loopLocal)
	}

	fg.scope.Define(&Symbol{
		Name:     iterName,
		Kind:     SymVar,
		Type:     Int,
		LocalIdx: counterLocal,
	})

	fg.body.Block(wasm.BlockEmpty{})
	fg.body.Loop(wasm.BlockEmpty{})
	fg.loopDepth++

	fg.body.LocalGet(counterLocal)
	fg.body.LocalGet(loopLocal)
	fg.body.I32GeS()
	fg.body.BrIf(1)

	fg.genBlock(ctx.Block())

	fg.body.LocalGet(counterLocal)
	fg.body.I32Const(1)
	fg.body.I32Add()
	fg.body.LocalSet(counterLocal)

	fg.body.Br(0)
	fg.loopDepth--
	fg.body.End()
	fg.body.End()
}

func (fg *funcGen) extractRange(ctx parser.IExprContext) (parser.IExprContext, parser.IExprContext) {
	exprs := ctx.AllExpr()
	if len(exprs) == 2 {
		if ctx.HALF_OPEN_RANGE() != nil || ctx.ELLIPSIS() != nil {
			return exprs[0], exprs[1]
		}
	}
	return nil, nil
}

// ── Switch ────────────────────────────────────────────────────────────────────

func (fg *funcGen) genSwitch(ctx parser.ISwitchStmtContext) {
	switchLocal := fg.newLocal(Int)
	fg.genExpr(ctx.Expr())
	fg.body.LocalSet(switchLocal)

	cases := ctx.AllSwitchCase()
	for i, sc := range cases {
		patterns := sc.CasePatternList().AllCasePattern()
		for j, pat := range patterns {
			fg.body.LocalGet(switchLocal)
			fg.genCasePattern(pat)
			fg.body.I32Eq()
			if j > 0 {
				fg.body.I32Or()
			}
		}
		if i == 0 {
			fg.body.If(wasm.BlockEmpty{})
		} else {
			fg.body.Else()
			fg.body.If(wasm.BlockEmpty{})
		}
		for _, s := range sc.AllStmt() {
			fg.genStmt(s)
		}
	}

	if dc := ctx.DefaultCase(); dc != nil {
		fg.body.Else()
		for _, s := range dc.AllStmt() {
			fg.genStmt(s)
		}
	}

	for range cases {
		fg.body.End()
	}
}

func (fg *funcGen) genCasePattern(ctx parser.ICasePatternContext) {
	if ctx.Literal() != nil {
		fg.genLiteral(ctx.Literal())
	} else {
		fg.body.I32Const(0)
	}
}

// ── Return ────────────────────────────────────────────────────────────────────

func (fg *funcGen) genReturn(ctx parser.IReturnStmtContext) {
	if ctx.Expr() != nil {
		fg.genExpr(ctx.Expr())
	}
	fg.body.Return()
}

// ── Defer ─────────────────────────────────────────────────────────────────────

func (fg *funcGen) genDefer(ctx parser.IDeferStmtContext) {
	if ctx.Expr() != nil {
		t := fg.genExpr(ctx.Expr())
		if t != nil && t.Kind() != KindVoid {
			fg.body.Drop()
		}
	}
}

// ── Arithmetic helpers ────────────────────────────────────────────────────────

func (fg *funcGen) emitAdd(t Type) {
	if t != nil && (t.Kind() == KindInt64 || t.Kind() == KindUint64) {
		fg.body.I64Add()
	} else if t != nil && t.Kind() == KindFloat {
		fg.body.F32Add()
	} else if t != nil && t.Kind() == KindDouble {
		fg.body.F64Add()
	} else {
		fg.body.I32Add()
	}
}

func (fg *funcGen) emitSub(t Type) {
	if t != nil && (t.Kind() == KindInt64 || t.Kind() == KindUint64) {
		fg.body.I64Sub()
	} else {
		fg.body.I32Sub()
	}
}

func (fg *funcGen) emitMul(t Type) {
	if t != nil && (t.Kind() == KindInt64 || t.Kind() == KindUint64) {
		fg.body.I64Mul()
	} else {
		fg.body.I32Mul()
	}
}

func (fg *funcGen) emitDiv(t Type) {
	if t != nil && t.Kind() == KindUint64 {
		fg.body.I64DivU()
	} else if t != nil && t.Kind() == KindInt64 {
		fg.body.I64DivS()
	} else {
		fg.body.I32DivS()
	}
}

func (fg *funcGen) emitRem(t Type) {
	if t != nil && (t.Kind() == KindInt64 || t.Kind() == KindUint64) {
		fg.body.I64RemS()
	} else {
		fg.body.I32RemS()
	}
}

func (fg *funcGen) genLiteral(ctx parser.ILiteralContext) {
	switch {
	case ctx.DEC_INT_LIT() != nil:
		v, _ := strconv.ParseInt(ctx.DEC_INT_LIT().GetText(), 10, 64)
		fg.body.I32Const(int32(v))
	case ctx.HEX_INT_LIT() != nil:
		s := ctx.HEX_INT_LIT().GetText()[2:]
		v, _ := strconv.ParseInt(s, 16, 64)
		fg.body.I32Const(int32(v))
	case ctx.BIN_INT_LIT() != nil:
		s := ctx.BIN_INT_LIT().GetText()[2:]
		v, _ := strconv.ParseInt(s, 2, 64)
		fg.body.I32Const(int32(v))
	case ctx.DEC_FLOAT_LIT() != nil:
		v, _ := strconv.ParseFloat(ctx.DEC_FLOAT_LIT().GetText(), 32)
		fg.body.F32Const(float32(v))
	case ctx.STRING_LIT() != nil:
		s := stripQuotes(ctx.STRING_LIT().GetText())
		off := fg.cg.internString(s)
		fg.body.I32Const(off)
	case ctx.TRUE() != nil:
		fg.body.I32Const(1)
	case ctx.FALSE() != nil:
		fg.body.I32Const(0)
	case ctx.NIL() != nil:
		fg.body.I32Const(0)
	}
}