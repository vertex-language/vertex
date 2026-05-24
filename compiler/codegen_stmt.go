package compiler

import (
	"strconv"

	"github.com/vertex-language/compiler/wasm"
	"github.com/vertex-language/vertex/parser"
)

// funcGen holds the mutable state for generating a single function body.
type funcGen struct {
	cg     *CodeGen
	sig    *FuncSig
	scope  *Scope
	body   *wasm.FunctionBody

	// Local variable tracking
	localCount uint32
	localTypes []wasm.ValType // wasm type of each local (params first, then vars)
	paramCount uint32

	// Struct frame: name → linear-memory address
	frameAllocs map[string]int32

	// Loop label depth for break/continue
	loopDepth int
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

// structLiteralOf extracts the IStructLiteralExprContext from an expression
// if it is a direct struct-literal primary, otherwise returns nil.
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

// inferStructType returns the struct/class type when:
//   - declared is already a struct/class, OR
//   - expr is a struct literal whose type name resolves in scope.
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

// preScanBlock walks a block recursively, pre-allocating linear-memory frame
// slots for every struct/class variable declaration it encounters.
func (fg *funcGen) preScanBlock(ctx parser.IBlockContext) {
	if ctx == nil {
		return
	}
	for _, stmt := range ctx.AllStmt() {
		if vd := stmt.VarDeclStmt(); vd != nil {
			fg.preScanVarDecl(vd)
		}
		// Recurse into control-flow blocks so nested struct vars get frame slots.
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
				// SwitchCase doesn't have its own block; stmts are inline.
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

	// Resolve declared type; if absent try to infer from a struct literal.
	var declared Type
	if ctx.Type_() != nil {
		declared = ResolveType(ctx.Type_(), fg.cg.scope)
	}
	t := inferStructType(declared, ctx.Expr(), fg.cg.scope)
	if t == nil {
		return
	}
	size := SizeOf(t)
	if size <= 0 {
		return
	}
	addr := fg.cg.allocFrame(size)
	fg.frameAllocs[name] = addr
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

	// Lazily allocate a frame slot if the pre-scan missed this variable
	// (e.g. it lives inside a deeply nested block that preScanBlock now covers,
	// but add a safety net here anyway).
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

	// ── Struct / class variable stored in linear memory ───────────────────
	if addr, ok := fg.frameAllocs[name]; ok {
		// Resolve the concrete type if it was inferred from the literal.
		if declType == nil || declType.Kind() == KindVoid {
			declType = inferStructType(nil, ctx.Expr(), fg.cg.scope)
		}

		// Store the frame base address in a wasm local for field-access codegen.
		ptrLocal := fg.newLocal(Int)
		fg.body.I32Const(addr)
		fg.body.LocalSet(ptrLocal)

		if ctx.Expr() != nil {
			if sl := structLiteralOf(ctx.Expr()); sl != nil {
				// Struct literal: initialise each field individually.
				fg.genStructLitFields(addr, sl, declType)
			} else {
				// Any other expression must produce a struct address; memcopy.
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

	// ── Regular scalar local ───────────────────────────────────────────────
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

// genStructLitFields stores each named field of a struct literal into linear
// memory starting at addr. Nested struct literals are handled recursively.
//
// WebAssembly store instruction order:  [addr] [value] i32.store
// The base address is encoded as an i32.const; the field offset is the
// instruction's static immediate, so we never need to compute addr+offset
// at runtime.
func (fg *funcGen) genStructLitFields(addr int32, ctx parser.IStructLiteralExprContext, t Type) {
	st, ok := t.(*StructType)
	if !ok {
		// Also handle ClassType (non-native classes with value fields)
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

		// Nested struct literal: recurse into the sub-region.
		if f.Type.Kind() == KindStruct {
			if nestedSL := structLiteralOf(init.Expr()); nestedSL != nil {
				fg.genStructLitFields(addr+int32(f.Offset), nestedSL, f.Type)
				continue
			}
		}

		// Scalar or non-literal expression: push addr, push value, store.
		fg.body.I32Const(addr)
		fg.genExpr(init.Expr())
		fg.emitStore(f.Type, uint32(f.Offset))
	}
}

// genClassLitFields is the same as genStructLitFields but for ClassType.
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

	// ── Struct field or array element via lvalue chain: base.field = expr ──
	if lv.DOT() != nil || lv.LBRACKET() != nil {
		fieldType := fg.genLvalueAddr(lv) // pushes full address onto stack
		fg.genExpr(ctx.Expr())            // pushes value
		fg.emitStore(fieldType, 0)
		return
	}

	// ── Simple identifier ───────────────────────────────────────────────────
	name := lv.ID().GetText()
	sym := fg.scope.Lookup(name)
	if sym == nil {
		return
	}

	if sym.IsFrame {
		// Struct variable: either re-init from a literal or memcopy.
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

	// ── Struct field compound assign: base.field op= expr ──────────────────
	if lv.DOT() != nil || lv.LBRACKET() != nil {
		fieldType := fg.genLvalueAddr(lv)
		addrLocal := fg.newLocal(Int)
		fg.body.LocalSet(addrLocal)

		// Load current value.
		fg.body.LocalGet(addrLocal)
		fg.emitLoad(fieldType, 0)

		// Evaluate RHS.
		fg.genExpr(ctx.Expr())

		// Apply operator.
		fg.applyCompoundOp(ctx.CompoundOp(), fieldType)

		// Store result: need [addr][val] — save result, load addr, load result.
		valLocal := fg.newLocal(fieldType)
		fg.body.LocalSet(valLocal)
		fg.body.LocalGet(addrLocal)
		fg.body.LocalGet(valLocal)
		fg.emitStore(fieldType, 0)
		return
	}

	// ── Simple variable ─────────────────────────────────────────────────────
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

// genLvalueAddr pushes the byte address of the lvalue onto the wasm stack and
// returns the Type of the location (used to select the right store instruction).
//
// Only lvalues rooted in frame-allocated struct variables are addressable.
func (fg *funcGen) genLvalueAddr(ctx parser.ILvalueContext) Type {
	if ctx.DOT() != nil {
		// Recurse: get address of base, then add field offset.
		baseType := fg.genLvalueAddr(ctx.Lvalue())
		return fg.resolveFieldAddr(baseType, ctx.ID().GetText())
	}
	if ctx.LBRACKET() != nil {
		// array[index]: base_addr + index * elemSize
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
	// Base case: simple identifier.
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
	// Non-frame scalar — not truly addressable, but surface the value so
	// callers can handle it (e.g. a pointer passed to a function).
	fg.body.LocalGet(sym.LocalIdx)
	return sym.Type
}

// resolveFieldAddr consumes the base address on the wasm stack, adds the
// named field's byte offset, and returns the field's type.
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

// emitLoad emits the wasm load instruction for type t.
// The memory address must already be on the stack; offset is a static immediate.
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

// emitStore emits the wasm store instruction for type t.
// Stack state expected: [..., i32_addr, value]. offset is a static immediate.
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

// emitStructMemcopy copies SizeOf(srcType) bytes from the address currently on
// the wasm stack into the known static destination address dstAddr.
// Copies 4 bytes at a time (LayoutStruct always pads to word alignment).
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