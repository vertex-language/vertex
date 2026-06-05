package compiler

import (
	"fmt"

	cir "github.com/vertex-language/ir/c"
)

// ─────────────────────────────────────────────────────────────────────────────
// Lowerer — converts a resolved AST into ir/c builder calls.
//
// Passes within the Lowerer:
//   1. registerTypes   — emit struct/class/enum typedef definitions.
//   2. lowerFunctions  — emit function definitions (including receivers).
// ─────────────────────────────────────────────────────────────────────────────

// Lowerer holds all mutable state for a single module lowering.
type Lowerer struct {
	diags   *Diagnostics
	mod     *cir.Module
	gt      *glibTypes

	structTypes   map[string]*cir.StructType
	classTypes    map[string]*cir.StructType
	enumTypes     map[string]cir.Type
	nativeClasses map[string]bool

	// testEntryFunc, when non-empty, identifies the single test function being
	// compiled. lowerFunctions skips all other test-qualified functions and
	// injects a main() that calls this one.
	testEntryFunc string

	tempSeq int
}

func NewLowerer(diags *Diagnostics, mod *cir.Module) *Lowerer {
    gt := newGlibTypes()
    setupGLib(mod, gt)
    return &Lowerer{
        diags:         diags,
        mod:           mod,
        gt:            gt,
        structTypes:   make(map[string]*cir.StructType),
        classTypes:    make(map[string]*cir.StructType),
        enumTypes:     make(map[string]cir.Type),
        nativeClasses: make(map[string]bool),
    }
}

// LowerFile drives both lowering passes over file.
func (l *Lowerer) LowerFile(file *File) {
	l.registerTypes(file)
	l.lowerFunctions(file)
}

func (l *Lowerer) tempName() string {
	l.tempSeq++
	return fmt.Sprintf("_t%d", l.tempSeq)
}

// ─── Pass 1: type registration ────────────────────────────────────────────────

func (l *Lowerer) registerTypes(file *File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *StructDecl:
			l.registerStruct(d)
		case *ClassDecl:
			l.registerClass(d)
		case *EnumDecl:
			l.registerEnum(d)
		}
	}
}

func (l *Lowerer) registerStruct(d *StructDecl) {
	var fields []cir.FieldDef
	for _, f := range d.Fields {
		ft := l.resolvedFieldCIRType(f.Type)
		if ft != nil {
			fields = append(fields, cir.Field(f.Name, ft))
		}
	}
	st := cir.Struct(d.Name, fields...)
	l.mod.RegisterType(st)
	l.structTypes[d.Name] = st
}

func (l *Lowerer) registerClass(d *ClassDecl) {
    // Native-interface class (class Foo : pkg): register C externs, not a struct.
    if d.BaseName != "" {
        l.nativeClasses[d.Name] = true
        for _, m := range d.Members {
            if m.IsField {
                continue
            }
            l.registerNativeMethod(m)
        }
        return
    }
    // Regular class: emit a struct definition.
    var fields []cir.FieldDef
    for _, m := range d.Members {
        if !m.IsField {
            continue
        }
        ft := l.resolvedFieldCIRType(m.Type)
        if ft != nil {
            fields = append(fields, cir.Field(m.Name, ft))
        }
    }
    st := cir.Struct(d.Name, fields...)
    l.mod.RegisterType(st)
    l.classTypes[d.Name] = st
}

// registerNativeMethod forwards a native-interface method signature as a C extern.
// If the last parameter is variadic ("..."), cir.Variadic is appended so the
// emitted C declaration gets the correct "..." suffix (e.g. printf).
func (l *Lowerer) registerNativeMethod(m *ClassMember) {
    if l.mod.LookupExtern(m.Name) != nil {
        return
    }
    retCIR := cir.Type(cir.Void)
    if m.RetType != nil {
        if ct := l.vtypeToCIR(l.resolveTypeExprVType(m.RetType)); ct != nil {
            retCIR = ct
        }
    }
    opts := []cir.FuncOpt{cir.Returns(retCIR)}
    for _, p := range m.Params {
        pt := l.resolveTypeExprVType(p.Type)
        ct := l.vtypeToCIR(pt)
        if ct == nil {
            ct = cir.VoidPtr
        }
        // Always emit the named param. For a variadic param (fmt: ...*const char)
        // this adds  Param("fmt", ConstPtr(Char))  and then appends  Variadic,
        // producing the correct C signature:  const char* fmt, ...
        opts = append(opts, cir.Param(p.Name, ct))
        if p.IsVariadic {
            opts = append(opts, cir.Variadic)
        }
    }
    l.mod.Extern(m.Name, opts...)
}

func (l *Lowerer) registerEnum(d *EnumDecl) {
	// Enums lower to C enums (int32_t).
	// ir/c doesn't have a native enum builder, so we use the EnumDecl
	// and emit the typedef via a type alias registered as Int32.
	// The actual C enum text comes from EmitC. For now we just track the type.
	l.enumTypes[d.Name] = cir.Int32
}

// resolvedFieldCIRType resolves a field's TypeExpr to a CIR type,
// handling structs that may be referenced before their definition is registered.
func (l *Lowerer) resolvedFieldCIRType(te TypeExpr) cir.Type {
	if te == nil {
		return cir.VoidPtr
	}
	switch t := te.(type) {
	case *NamedTypeExpr:
		if bt, ok := BuiltinTypes[t.Name]; ok {
			if ct := bt.CIRType(); ct != nil {
				return ct
			}
		}
		// May be a struct or class.
		if st, ok := l.structTypes[t.Name]; ok {
			return st
		}
		if st, ok := l.classTypes[t.Name]; ok {
			return cir.Ptr(st)
		}
		return cir.VoidPtr
	case *PointerTypeExpr:
		elem := l.resolvedFieldCIRType(t.Elem)
		if elem == nil {
			elem = cir.Void
		}
		if t.IsConst {
			return cir.ConstPtr(elem)
		}
		return cir.Ptr(elem)
	case *OptionalTypeExpr:
		return cir.Ptr(l.resolvedFieldCIRType(t.Elem))
	case *ArrayTypeExpr:
		return l.gt.GArrayPtr
	}
	return cir.VoidPtr
}

// ─── Pass 2: function lowering ────────────────────────────────────────────────

func (l *Lowerer) lowerFunctions(file *File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *FuncDecl:
			// In test mode, skip every test function except the active entry.
			if l.testEntryFunc != "" &&
				d.Qualifier == FuncQualTest &&
				d.Name != l.testEntryFunc {
				continue
			}
			l.lowerFuncDecl(d)
		case *VarDecl:
			l.lowerGlobalVar(d)
		}
	}
	// Inject a standalone main() when compiling a single test function.
	if l.testEntryFunc != "" {
		l.injectTestMain()
	}
}

func (l *Lowerer) lowerGlobalVar(d *VarDecl) {
	if len(d.Binding.Names) == 0 {
		return
	}
	vt := d.Value.GetVType()
	if vt == nil {
		return
	}
	ct := l.vtypeToCIR(vt)
	if ct == nil {
		ct = cir.VoidPtr
	}
	name := d.Binding.Names[0]
	// Only constant literals can be global initialisers in C.
	if d.IsLet {
		switch v := d.Value.(type) {
		case *IntLitExpr:
			l.mod.Global(name, ct, cir.IntLit(v.Value))
		case *FloatLitExpr:
			l.mod.Global(name, ct, cir.FloatLit(v.Value))
		case *BoolLitExpr:
			l.mod.Global(name, ct, cir.BoolLit(v.Value))
		case *StringLitExpr:
			l.mod.StringLit(name, v.Value)
		default:
			l.mod.GlobalZero(name, ct)
		}
	} else {
		l.mod.GlobalZero(name, ct)
	}
}

// lowerFuncDecl emits one ir/c FuncDef for a Vertex function or method.
func (l *Lowerer) lowerFuncDecl(fn *FuncDecl) {
	// Test-qualified functions use their own lowering path.
	if fn.Qualifier == FuncQualTest {
		l.lowerTestFuncDecl(fn)
		return
	}
	if fn.Qualifier != FuncQualNone {
		l.diags.Warnf(fn.Pos, "function qualifier %v is not yet supported; ignoring", fn.Qualifier)
	}

	cName := fn.Name
	if fn.Receiver != nil {
		recvTypeName := extractTypeName(fn.Receiver.Type)
		cName = recvTypeName + "__" + fn.Name
	}

	retType := cir.Void
	if fn.RetType != nil {
		vrt := l.resolveTypeExprVType(fn.RetType)
		if ct := l.vtypeToCIR(vrt); ct != nil {
			retType = ct
		}
	}

	opts := []cir.FuncOpt{cir.Returns(retType)}

	if fn.Receiver != nil {
		recvCType := l.receiverCIRType(fn.Receiver)
		opts = append(opts, cir.Param(fn.Receiver.Name, recvCType))
	}

	for _, p := range fn.Params {
		if p.IsVariadic {
			opts = append(opts, cir.Variadic)
			continue
		}
		pt := l.resolveTypeExprVType(p.Type)
		ct := l.vtypeToCIR(pt)
		if ct == nil {
			ct = cir.VoidPtr
		}
		opts = append(opts, cir.Param(p.Name, ct))
	}

	def := l.mod.Func(cName, opts...)

	if fn.Body == nil {
		return
	}

	fc := newFuncCtx()
	if fn.Receiver != nil {
		fc.params[fn.Receiver.Name] = true
	}
	for _, p := range fn.Params {
		fc.params[p.Name] = true
	}

	def.Body(func(b *cir.Builder) {
		l.lowerBlock(b, fn.Body, fc)
		if retType == cir.Void && !blockEndsWithReturn(fn.Body) {
			fc.emitDefers(b)
			b.Return()
		}
	})
}

func (l *Lowerer) receiverCIRType(recv *Receiver) cir.Type {
	typeName := extractTypeName(recv.Type)
	if st, ok := l.classTypes[typeName]; ok {
		return cir.Ptr(st)
	}
	if st, ok := l.structTypes[typeName]; ok {
		if recv.IsPtr {
			return cir.Ptr(st)
		}
		return st
	}
	return cir.VoidPtr
}

// ─── Block and statement lowering ─────────────────────────────────────────────

// funcCtx carries per-function state during lowering.
type funcCtx struct {
	params     map[string]bool
	locals     map[string]cir.Expr
	defers     []func(*cir.Builder)
	isTestFunc bool // when true, return-expr → printf + void return
}

func newFuncCtx() *funcCtx {
	return &funcCtx{
		params: make(map[string]bool),
		locals: make(map[string]cir.Expr),
	}
}

func (fc *funcCtx) pushDefer(fn func(*cir.Builder)) {
	fc.defers = append(fc.defers, fn)
}

func (fc *funcCtx) emitDefers(b *cir.Builder) {
	for i := len(fc.defers) - 1; i >= 0; i-- {
		fc.defers[i](b)
	}
}

func (fc *funcCtx) varRef(b *cir.Builder, name string) cir.Expr {
	if fc.params[name] {
		return b.Param(name)
	}
	if ref, ok := fc.locals[name]; ok {
		return ref
	}
	// Fallback: assume it's a param (handles forward references).
	return b.Param(name)
}

func (l *Lowerer) lowerBlock(b *cir.Builder, blk *BlockStmt, fc *funcCtx) {
	for _, s := range blk.Stmts {
		l.lowerStmt(b, s, fc)
	}
}

func (l *Lowerer) lowerStmt(b *cir.Builder, s Stmt, fc *funcCtx) {
	switch st := s.(type) {
	case *LocalDeclStmt:
		l.lowerLocalDecl(b, st.Decl, fc)
	case *ExprStmt:
		if e := l.lowerExpr(b, st.Expr, fc); e != nil {
			b.Stmt(e)
		}
	case *AssignStmt:
		l.lowerAssign(b, st, fc)
	case *ReturnStmt:
		fc.emitDefers(b)
		if st.Value != nil {
			val := l.lowerExpr(b, st.Value, fc)
			if fc.isTestFunc {
				// Test mode: write the value to stdout then void-return.
				// The format specifier is derived from the expression's resolved type.
				fmtStr := l.printfFormatFor(st.Value.GetVType()) + "\n"
				fmtLit := l.mod.StringLit(l.tempName(), fmtStr)
				b.Stmt(b.Call("printf", fmtLit, val))
				b.Return()
			} else {
				b.ReturnVal(val)
			}
		} else {
			b.Return()
		}
	case *DeferStmt:
		// Capture the call expression for deferred emission.
		call := st.Call
		fc.pushDefer(func(b *cir.Builder) {
			if e := l.lowerExpr(b, call, fc); e != nil {
				b.Stmt(e)
			}
		})
	case *IfStmt:
		l.lowerIfStmt(b, st, fc)
	case *WhileStmt:
		cond := l.lowerExpr(b, st.Cond, fc)
		body := cir.B(func(b *cir.Builder) { l.lowerBlock(b, st.Body, fc) })
		b.While(cond, body)
	case *ForInStmt:
		l.lowerForIn(b, st, fc)
	case *SwitchStmt:
		l.lowerSwitch(b, st, fc)
	case *BreakStmt:
		b.Break()
	case *ContinueStmt:
		b.Continue()
	case *FallthroughStmt:
		// In ir/c, fallthrough is modelled with CaseFallthrough — handled in switch.
	case *BlockStmt:
		l.lowerBlock(b, st, fc)
	}
}

func (l *Lowerer) lowerLocalDecl(b *cir.Builder, d *VarDecl, fc *funcCtx) {
	vtype := d.Value.GetVType()
	if d.TypeHint != nil {
		vtype = l.resolveTypeExprVType(d.TypeHint)
	}
	if vtype == nil {
		vtype = &VUnknown{}
	}

	for i, name := range d.Binding.Names {
		// For tuple destructuring, skip elements beyond what the value provides.
		_ = i
		ref := l.declareLocal(b, name, vtype, d.IsLet)
		fc.locals[name] = ref

		// Emit initialiser.
		initExpr := l.lowerExpr(b, d.Value, fc)
		if initExpr != nil {
			// NEW: Safely cast the initializer to the target variable's type
			if targetCT := l.vtypeToCIR(vtype); targetCT != nil {
				initExpr = b.Cast(targetCT, initExpr)
			}
			b.Assign(ref, initExpr)
		}
	}
}

func (l *Lowerer) declareLocal(b *cir.Builder, name string, vt VType, isLet bool) cir.Expr {
	ct := l.vtypeToCIR(vt)
	if ct == nil {
		ct = l.vtypeToCIRFallback(vt)
	}
	ref := b.Local(name, ct)
	return ref
}

func (l *Lowerer) vtypeToCIRFallback(vt VType) cir.Type {
	switch t := vt.(type) {
	case *VDynArray:
		return l.gt.GArrayPtr
	case *VString:
		if t.Mutable {
			return l.gt.GStringPtr
		}
		return cir.ConstPtr(cir.Char)
	case *VStruct:
		if st, ok := l.structTypes[t.Name]; ok {
			return st
		}
	case *VMap:
		return cir.VoidPtr // GHashTable is opaque void*
	case *VClass:
		if st, ok := l.classTypes[t.Name]; ok {
			return cir.Ptr(st)
		}
	case *VEnum:
		return cir.Int32
	case *VOptional:
		inner := l.vtypeToCIR(t.Elem)
		if inner == nil {
			inner = l.vtypeToCIRFallback(t.Elem)
		}
		return cir.Ptr(inner)
	case *VResult:
		return cir.VoidPtr // simplified
	}
	return cir.VoidPtr
}

func (l *Lowerer) lowerAssign(b *cir.Builder, st *AssignStmt, fc *funcCtx) {
	lhs := l.lowerExpr(b, st.LHS, fc)
	rhs := l.lowerExpr(b, st.RHS, fc)
	if lhs == nil || rhs == nil {
		return
	}
	switch st.Op {
	case OpAssign:
		b.Assign(lhs, rhs)
	case OpAddAssign:
		b.Assign(lhs, b.Add(lhs, rhs))
	case OpSubAssign:
		b.Assign(lhs, b.Sub(lhs, rhs))
	case OpMulAssign:
		b.Assign(lhs, b.Mul(lhs, rhs))
	case OpDivAssign:
		b.Assign(lhs, b.Div(lhs, rhs))
	case OpModAssign:
		b.Assign(lhs, b.Mod(lhs, rhs))
	}
}

func (l *Lowerer) lowerIfStmt(b *cir.Builder, st *IfStmt, fc *funcCtx) {
	var cond cir.Expr
	switch c := st.Cond.(type) {
	case *IfExprCond:
		cond = l.lowerExpr(b, c.Expr, fc)
	case *IfLetCond:
		// Evaluate the optional expression.
		val := l.lowerExpr(b, c.Expr, fc)
		// Store into a temp, bind it.
		tmp := b.Local(l.tempName(), l.vtypeToCIRFallback(c.Expr.GetVType()))
		b.Assign(tmp, val)
		fc.locals[c.Name] = tmp
		cond = b.Neq(tmp, cir.NullPtr())
	}
	if cond == nil {
		cond = cir.BoolLit(true)
	}
	thenBlk := cir.B(func(b *cir.Builder) { l.lowerBlock(b, st.Then, fc) })
	if st.Else == nil {
		b.If(cond, thenBlk)
	} else {
		elseBlk := cir.B(func(b *cir.Builder) { l.lowerStmt(b, st.Else, fc) })
		b.IfElse(cond, thenBlk, elseBlk)
	}
}

func (l *Lowerer) lowerForIn(b *cir.Builder, st *ForInStmt, fc *funcCtx) {
	iterType := st.Iter.GetVType()
	iter := l.lowerExpr(b, st.Iter, fc)
	if iter == nil {
		return
	}
	switch it := iterType.(type) {
	case *VDynArray:
		elemCIR := l.vtypeToCIR(it.Elem)
		if elemCIR == nil {
			elemCIR = l.vtypeToCIRFallback(it.Elem)
		}
		iRef := b.Local(l.tempName(), cir.UInt32)
		b.Assign(iRef, cir.UIntLit(0))
		arrLen := b.GetField(iter, "len")
		b.While(b.Lt(iRef, arrLen), cir.B(func(b *cir.Builder) {
			dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(iter, "data"))
			elemRef := b.Local(st.Var, elemCIR)
			fc.locals[st.Var] = elemRef
			b.Assign(elemRef, b.Index(dataPtr, iRef))
			l.lowerBlock(b, st.Body, fc)
			b.Assign(iRef, b.Add(iRef, cir.UIntLit(1)))
		}))
	default:
		// Range iteration: for i in lo..<hi or lo...hi
		if be, ok := st.Iter.(*BinaryExpr); ok &&
			(be.Op == BinRangeHalfOpen || be.Op == BinRangeClosed) {
			lo := l.lowerExpr(b, be.Left, fc)
			hi := l.lowerExpr(b, be.Right, fc)
			iRef := b.Local(st.Var, cir.Int32)
			fc.locals[st.Var] = iRef
			b.Assign(iRef, lo)
			var cond cir.Expr
			if be.Op == BinRangeHalfOpen {
				cond = b.Lt(iRef, hi)
			} else {
				cond = b.Lte(iRef, hi)
			}
			b.While(cond, cir.B(func(b *cir.Builder) {
				l.lowerBlock(b, st.Body, fc)
				b.Assign(iRef, b.Add(iRef, cir.IntLit(1)))
			}))
			return
		}
		l.diags.Warnf(st.Pos, "for-in over unsupported type %v; skipping", iterType)
	}
}

func (l *Lowerer) lowerSwitch(b *cir.Builder, st *SwitchStmt, fc *funcCtx) {
	subj := l.lowerExpr(b, st.Subj, fc)
	if subj == nil {
		return
	}
	// Check if any case uses string comparison (string switch → if-else chain).
	if _, isStr := st.Subj.GetVType().(*VString); isStr {
		l.lowerStringSwitch(b, st, subj, fc)
		return
	}
	// Integer / enum switch.
	var cases []cir.SwitchCase
	for _, c := range st.Cases {
		body := cir.B(func(b *cir.Builder) {
			for _, s := range c.Body {
				l.lowerStmt(b, s, fc)
			}
		})
		if c.IsDefault {
			cases = append(cases, cir.Default(body))
		} else {
			for _, pat := range c.Patterns {
				var val cir.Expr
				switch p := pat.(type) {
				case *ExprPattern:
					val = l.lowerExpr(b, p.Expr, fc)
				case *EnumShortPattern:
					val = l.enumCaseExpr(b, st.Subj.GetVType(), p.Case)
				}
				if val != nil {
					if isFallthroughCase(c) {
						cases = append(cases, cir.CaseFallthrough(val, body))
					} else {
						cases = append(cases, cir.Case(val, body))
					}
				}
			}
		}
	}
	b.Switch(subj, cases...)
}

func (l *Lowerer) lowerStringSwitch(b *cir.Builder, st *SwitchStmt, subj cir.Expr, fc *funcCtx) {
	// Emit as if-else chain using strcmp.
	// First ensure strcmp is declared.
	if l.mod.LookupExtern("strcmp") == nil {
		l.mod.Extern("strcmp",
			cir.Returns(cir.Int32),
			cir.Param("s1", cir.ConstPtr(cir.Char)),
			cir.Param("s2", cir.ConstPtr(cir.Char)),
		)
	}
	var chains []struct {
		cond cir.Expr
		body *cir.Block
	}
	var defaultBody *cir.Block

	for _, c := range st.Cases {
		body := cir.B(func(b *cir.Builder) {
			for _, s := range c.Body {
				l.lowerStmt(b, s, fc)
			}
		})
		if c.IsDefault {
			defaultBody = body
			continue
		}
		for _, pat := range c.Patterns {
			if ep, ok := pat.(*ExprPattern); ok {
				lit := l.lowerExpr(b, ep.Expr, fc)
				cond := b.Eq(b.Call("strcmp", subj, lit), cir.IntLit(0))
				chains = append(chains, struct {
					cond cir.Expr
					body *cir.Block
				}{cond, body})
			}
		}
	}
	if len(chains) == 0 {
		if defaultBody != nil {
			b.Inline(defaultBody)
		}
		return
	}
	// Build nested if-else.
	var build func(i int) *cir.Block
	build = func(i int) *cir.Block {
		if i >= len(chains) {
			if defaultBody != nil {
				return defaultBody
			}
			return cir.B(func(*cir.Builder) {})
		}
		ch := chains[i]
		return cir.B(func(b *cir.Builder) {
			b.IfElse(ch.cond, ch.body, build(i+1))
		})
	}
	b.Inline(build(0))
}

// ─── Expression lowering ──────────────────────────────────────────────────────

func (l *Lowerer) lowerExpr(b *cir.Builder, expr Expr, fc *funcCtx) cir.Expr {
	if expr == nil {
		return cir.NullPtr()
	}
	switch e := expr.(type) {
	case *IntLitExpr:
		if e.IsUnsigned {
			return cir.UIntLit(uint64(e.Value))
		}
		return cir.IntLit(e.Value)
	case *FloatLitExpr:
		if vt, ok := e.GetVType().(*VFloat); ok && vt.Bits == 32 {
			return cir.FloatLit32(e.Value)
		}
		return cir.FloatLit(e.Value)
	case *BoolLitExpr:
		return cir.BoolLit(e.Value)
	case *CharLitExpr:
		// A char literal is a single code point; emit it as an integer constant
		// so it is compatible with C's char arithmetic and comparisons.
		return cir.IntLit(int64(e.Value))
	case *StringLitExpr:
		return l.mod.StringLit(l.tempName(), e.Value)
	case *NilLitExpr:
		return cir.NullPtr()

	case *IdentExpr:
		return fc.varRef(b, e.Name)

	case *DotEnumExpr:
		return l.lowerDotEnum(e)

	case *BinaryExpr:
		return l.lowerBinaryExpr(b, e, fc)

	case *UnaryExpr:
		return l.lowerUnaryExpr(b, e, fc)

	case *TernaryExpr:
		return l.lowerTernary(b, e, fc)

	case *CallExpr:
		return l.lowerCallExpr(b, e, fc)

	case *MethodCallExpr:
		return l.lowerMethodCall(b, e, fc)

	case *FieldExpr:
		return l.lowerFieldExpr(b, e, fc)

	case *IndexExpr:
		return l.lowerIndexExpr(b, e, fc)

	case *StructLitExpr:
		return l.lowerStructLit(b, e, fc)

	case *ArrayLitExpr:
		return l.lowerArrayLit(b, e, fc)

	case *ArrayCtorExpr:
		return l.lowerArrayCtor(b, e, fc)

	case *MapLitExpr:
		return l.lowerMapLit(b, e, fc)

	case *TypeConvExpr:
		inner := l.lowerExpr(b, e.Value, fc)
		targetVT := l.resolveTypeExprVType(e.TargetType)
		if ct := l.vtypeToCIR(targetVT); ct != nil {
			return b.Cast(ct, inner)
		}
		return inner

	case *ReinterpretExpr:
		inner := l.lowerExpr(b, e.Value, fc)
		targetVT := l.resolveTypeExprVType(e.TargetType)
		ct := l.vtypeToCIR(targetVT)
		if ct == nil {
			ct = l.vtypeToCIRFallback(targetVT)
		}
		return b.Cast(ct, inner)

	case *ResultExpr:
		l.diags.Warnf(e.Pos, "Result type lowering is not fully supported yet")
		return l.lowerExpr(b, e.Value, fc)

	case *TupleLitExpr:
		l.diags.Warnf(e.Pos, "tuple literals are not yet lowered")
		return cir.NullPtr()
	}
	return cir.NullPtr()
}

func (l *Lowerer) lowerDotEnum(e *DotEnumExpr) cir.Expr {
	// .caseName → integer literal (enum case value).
	// Without a specific enum type in context, emit 0 as a placeholder.
	// The resolver would need to propagate the enum type for full correctness.
	return cir.IntLit(0)
}

func (l *Lowerer) enumCaseExpr(b *cir.Builder, subjType VType, caseName string) cir.Expr {
	if ev, ok := subjType.(*VEnum); ok {
		// Find the case index.
		for i, c := range ev.Decl.Cases {
			if c.Name == caseName {
				if c.RawValue != nil {
					if il, ok := c.RawValue.(*IntLitExpr); ok {
						return cir.IntLit(il.Value)
					}
				}
				return cir.IntLit(int64(i))
			}
		}
	}
	return cir.IntLit(0)
}

func (l *Lowerer) lowerBinaryExpr(b *cir.Builder, e *BinaryExpr, fc *funcCtx) cir.Expr {
	left := l.lowerExpr(b, e.Left, fc)
	right := l.lowerExpr(b, e.Right, fc)
	switch e.Op {
	case BinAdd:
		return b.Add(left, right)
	case BinSub:
		return b.Sub(left, right)
	case BinMul:
		return b.Mul(left, right)
	case BinDiv:
		return b.Div(left, right)
	case BinMod:
		return b.Mod(left, right)
	case BinShl:
		return b.Shl(left, right)
	case BinShr:
		return b.Shr(left, right)
	case BinBitAnd:
		return b.And(left, right)
	case BinBitXor:
		return b.Xor(left, right)
	case BinBitOr:
		return b.Or(left, right)
	case BinEq:
		return b.Eq(left, right)
	case BinNeq:
		return b.Neq(left, right)
	case BinLt:
		return b.Lt(left, right)
	case BinLte:
		return b.Lte(left, right)
	case BinGt:
		return b.Gt(left, right)
	case BinGte:
		return b.Gte(left, right)
	case BinAnd:
		return b.LogAnd(left, right)
	case BinOr:
		return b.LogOr(left, right)
	case BinNilCoalesce:
		// left ?? right  →  tmp = left; result = (tmp != NULL) ? tmp : right
		ct := l.vtypeToCIRFallback(e.Left.GetVType())
		tmp := b.Local(l.tempName(), ct)
		result := b.Local(l.tempName(), ct)
		b.Assign(tmp, left)
		b.IfElse(b.Neq(tmp, cir.NullPtr()),
			cir.B(func(b *cir.Builder) { b.Assign(result, tmp) }),
			cir.B(func(b *cir.Builder) { b.Assign(result, right) }),
		)
		return result
	case BinOverflowAdd:
		return b.Add(b.Cast(cir.UInt32, left), b.Cast(cir.UInt32, right))
	case BinOverflowSub:
		return b.Sub(b.Cast(cir.UInt32, left), b.Cast(cir.UInt32, right))
	case BinOverflowMul:
		return b.Mul(b.Cast(cir.UInt32, left), b.Cast(cir.UInt32, right))
	case BinIdentityEq:
		return b.Eq(left, right)
	case BinIdentityNeq:
		return b.Neq(left, right)
	}
	return left
}

func (l *Lowerer) lowerUnaryExpr(b *cir.Builder, e *UnaryExpr, fc *funcCtx) cir.Expr {
	op := l.lowerExpr(b, e.Operand, fc)
	switch e.Op {
	case UnNeg:
		return b.Neg(op)
	case UnNot:
		return b.Not(op)
	case UnBitNot:
		return b.BitNot(op)
	case UnAddrOf:
		return b.AddrOf(op)
	}
	return op
}

// lowerTernary emits a ternary as a temp + if-else.
func (l *Lowerer) lowerTernary(b *cir.Builder, e *TernaryExpr, fc *funcCtx) cir.Expr {
	vt := e.GetVType()
	ct := l.vtypeToCIR(vt)
	if ct == nil {
		ct = l.vtypeToCIRFallback(vt)
	}
	tmp := b.Local(l.tempName(), ct)
	cond := l.lowerExpr(b, e.Cond, fc)
	then := l.lowerExpr(b, e.Then, fc)
	els := l.lowerExpr(b, e.Else, fc)
	b.IfElse(cond,
		cir.B(func(b *cir.Builder) { b.Assign(tmp, then) }),
		cir.B(func(b *cir.Builder) { b.Assign(tmp, els) }),
	)
	return tmp
}

func (l *Lowerer) lowerCallExpr(b *cir.Builder, e *CallExpr, fc *funcCtx) cir.Expr {
    if id, ok := e.Func.(*IdentExpr); ok {
        if bt, isBT := BuiltinTypes[id.Name]; isBT {
            ct := bt.CIRType()
            if ct != nil && len(e.Args) == 1 {
                inner := l.lowerExpr(b, e.Args[0].Value, fc)
                return b.Cast(ct, inner)
            }
        }
        // Regular or native-interface class instantiation.
        if _, isClass := l.classTypes[id.Name]; isClass {
            return l.lowerClassInstantiate(b, id.Name, e.Args, fc)
        }
        if l.nativeClasses[id.Name] {
            return l.lowerClassInstantiate(b, id.Name, e.Args, fc)
        }
    }

    fn := l.lowerExpr(b, e.Func, fc)
    var args []cir.Expr
    for _, a := range e.Args {
        args = append(args, l.lowerExpr(b, a.Value, fc))
    }
    if id, ok := e.Func.(*IdentExpr); ok {
        return b.Call(id.Name, args...)
    }
    return b.CallPtr(fn, args...)
}

func (l *Lowerer) lowerClassInstantiate(b *cir.Builder, className string, args []*Arg, fc *funcCtx) cir.Expr {
    // Native-interface class: no underlying struct, no malloc/init.
    // The variable is a void* null used only to drive method dispatch.
    if l.nativeClasses[className] {
        return cir.NullPtr()
    }
    st := l.classTypes[className]
    ptrType := cir.Ptr(st)
    tmp := b.Local(l.tempName(), ptrType)
    b.Assign(tmp, b.Cast(ptrType, b.Call("malloc", b.SizeOf(st))))
    var initArgs []cir.Expr
    initArgs = append(initArgs, tmp)
    for _, a := range args {
        initArgs = append(initArgs, l.lowerExpr(b, a.Value, fc))
    }
    b.Stmt(b.Call(className+"__init", initArgs...))
    return tmp
}

func (l *Lowerer) lowerMethodCall(b *cir.Builder, e *MethodCallExpr, fc *funcCtx) cir.Expr {
	recv := l.lowerExpr(b, e.Recv, fc)
	recvType := e.Recv.GetVType()

	switch rt := recvType.(type) {
	case *VDynArray:
		return l.lowerDynArrayMethod(b, recv, rt, e.Method, e.Args, fc)
	case *VString:
		return l.lowerStringMethod(b, recv, rt, e.Method, e.Args, fc)
	case *VClass:
		return l.lowerClassMethod(b, recv, rt.Name, e.Method, e.Args, fc)
	case *VStruct:
		return l.lowerStructMethod(b, recv, rt.Name, e.Method, e.Args, fc)
	}

	// Generic method call: ClassName__method(recv, args...)
	if typeName := vtypeBaseName(recvType); typeName != "" {
		cName := typeName + "__" + e.Method
		var args []cir.Expr
		args = append(args, recv)
		for _, a := range e.Args {
			args = append(args, l.lowerExpr(b, a.Value, fc))
		}
		return b.Call(cName, args...)
	}

	l.diags.Warnf(e.Pos, "method %q on type %v not lowered", e.Method, recvType)
	return cir.NullPtr()
}

func (l *Lowerer) lowerDynArrayMethod(
	b *cir.Builder, recv cir.Expr, rt *VDynArray,
	method string, args []*Arg, fc *funcCtx,
) cir.Expr {
	elemCIR := l.vtypeToCIR(rt.Elem)
	if elemCIR == nil {
		elemCIR = l.vtypeToCIRFallback(rt.Elem)
	}
	elemSize := b.SizeOf(elemCIR)

	switch method {
	case "push":
		val := l.lowerExpr(b, args[0].Value, fc)
		tmp := b.Local(l.tempName(), elemCIR)
		b.Assign(tmp, val)
		b.Stmt(b.Call("g_array_append_vals", recv, b.AddrOf(tmp), cir.UIntLit(1)))
		return nil
	case "unshift":
		val := l.lowerExpr(b, args[0].Value, fc)
		tmp := b.Local(l.tempName(), elemCIR)
		b.Assign(tmp, val)
		b.Stmt(b.Call("g_array_prepend_vals", recv, b.AddrOf(tmp), cir.UIntLit(1)))
		return nil
	case "pop":
		tmp := b.Local(l.tempName(), elemCIR)
		lenRef := b.GetField(recv, "len")
		b.If(b.Gt(lenRef, cir.UIntLit(0)),
			cir.B(func(b *cir.Builder) {
				dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
				idx := b.Sub(lenRef, cir.UIntLit(1))
				b.Assign(tmp, b.Index(dataPtr, idx))
				b.Stmt(b.Call("g_array_remove_index", recv, idx))
			}),
		)
		return tmp
	case "shift":
		tmp := b.Local(l.tempName(), elemCIR)
		lenRef := b.GetField(recv, "len")
		b.If(b.Gt(lenRef, cir.UIntLit(0)),
			cir.B(func(b *cir.Builder) {
				dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
				b.Assign(tmp, b.Index(dataPtr, cir.UIntLit(0)))
				b.Stmt(b.Call("g_array_remove_index", recv, cir.UIntLit(0)))
			}),
		)
		return tmp
	case "delete":
		b.Stmt(b.Call("g_array_free", recv, cir.BoolLit(true)))
		return nil
	case "fill":
		val := l.lowerExpr(b, args[0].Value, fc)
		byteSize := b.Mul(b.Cast(cir.UIntSize, b.GetField(recv, "len")), b.SizeOf(elemCIR))
		b.Stmt(b.Call("memset", b.GetField(recv, "data"), val, byteSize))
		return nil
	case "reverse":
		lo := b.Local(l.tempName(), cir.UInt32)
		hi := b.Local(l.tempName(), cir.UInt32)
		tmp2 := b.Local(l.tempName(), elemCIR)
		b.Assign(lo, cir.UIntLit(0))
		b.Assign(hi, b.Sub(b.GetField(recv, "len"), cir.UIntLit(1)))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
		b.While(b.Lt(lo, hi),
			cir.B(func(b *cir.Builder) {
				b.Assign(tmp2, b.Index(dataPtr, lo))
				b.Assign(b.Index(dataPtr, lo), b.Index(dataPtr, hi))
				b.Assign(b.Index(dataPtr, hi), tmp2)
				b.Assign(lo, b.Add(lo, cir.UIntLit(1)))
				b.Assign(hi, b.Sub(hi, cir.UIntLit(1)))
			}),
		)
		return nil
	case "sort":
		cmpFn := l.lowerExpr(b, args[0].Value, fc)
		b.Stmt(b.Call("g_array_sort", recv, cmpFn))
		return nil
	case "indexOf":
		val := l.lowerExpr(b, args[0].Value, fc)
		result := b.Local(l.tempName(), cir.Int32)
		b.Assign(result, cir.IntLit(-1))
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
		b.While(b.Lt(i, b.GetField(recv, "len")),
			cir.B(func(b *cir.Builder) {
				b.If(b.Eq(b.Index(dataPtr, i), val),
					cir.B(func(b *cir.Builder) {
						b.Assign(result, b.Cast(cir.Int32, i))
						b.Break()
					}),
				)
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return result
	case "includes":
		val := l.lowerExpr(b, args[0].Value, fc)
		found := b.Local(l.tempName(), cir.Bool)
		b.Assign(found, cir.BoolLit(false))
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
		b.While(b.Lt(i, b.GetField(recv, "len")),
			cir.B(func(b *cir.Builder) {
				b.If(b.Eq(b.Index(dataPtr, i), val),
					cir.B(func(b *cir.Builder) {
						b.Assign(found, cir.BoolLit(true))
						b.Break()
					}),
				)
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return found
	case "find":
		cmpFn := l.lowerExpr(b, args[0].Value, fc)
		resultPtr := b.Local(l.tempName(), cir.Ptr(elemCIR))
		b.Assign(resultPtr, cir.NullPtr())
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
		b.While(b.Lt(i, b.GetField(recv, "len")),
			cir.B(func(b *cir.Builder) {
				elem := b.AddrOf(b.Index(dataPtr, i))
				b.If(b.CallPtr(cmpFn, b.Deref(elem)),
					cir.B(func(b *cir.Builder) {
						b.Assign(resultPtr, elem)
						b.Break()
					}),
				)
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return resultPtr
	case "map":
		mapFn := l.lowerExpr(b, args[0].Value, fc)
		out := b.Local(l.tempName(), l.gt.GArrayPtr)
		b.Assign(out, b.Call("g_array_sized_new",
			cir.BoolLit(false), cir.BoolLit(false),
			elemSize, b.GetField(recv, "len")))
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
		b.While(b.Lt(i, b.GetField(recv, "len")),
			cir.B(func(b *cir.Builder) {
				mapped := b.Local(l.tempName(), elemCIR)
				b.Assign(mapped, b.CallPtr(mapFn, b.Index(dataPtr, i)))
				b.Stmt(b.Call("g_array_append_vals", out, b.AddrOf(mapped), cir.UIntLit(1)))
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return out
	case "filter":
		filterFn := l.lowerExpr(b, args[0].Value, fc)
		out := b.Local(l.tempName(), l.gt.GArrayPtr)
		b.Assign(out, b.Call("g_array_new",
			cir.BoolLit(false), cir.BoolLit(false), elemSize))
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
		b.While(b.Lt(i, b.GetField(recv, "len")),
			cir.B(func(b *cir.Builder) {
				elem := b.Index(dataPtr, i)
				b.If(b.CallPtr(filterFn, elem),
					cir.B(func(b *cir.Builder) {
						tmp := b.Local(l.tempName(), elemCIR)
						b.Assign(tmp, elem)
						b.Stmt(b.Call("g_array_append_vals", out, b.AddrOf(tmp), cir.UIntLit(1)))
					}),
				)
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return out
	case "slice":
		start := l.lowerExpr(b, args[0].Value, fc)
		end := l.lowerExpr(b, args[1].Value, fc)
		out := b.Local(l.tempName(), l.gt.GArrayPtr)
		b.Assign(out, b.Call("g_array_new",
			cir.BoolLit(false), cir.BoolLit(false), elemSize))
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, b.Cast(cir.UInt32, start))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
		b.While(b.LogAnd(b.Lt(i, b.Cast(cir.UInt32, end)), b.Lt(i, b.GetField(recv, "len"))),
			cir.B(func(b *cir.Builder) {
				tmp := b.Local(l.tempName(), elemCIR)
				b.Assign(tmp, b.Index(dataPtr, i))
				b.Stmt(b.Call("g_array_append_vals", out, b.AddrOf(tmp), cir.UIntLit(1)))
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return out
	case "concat":
		other := l.lowerExpr(b, args[0].Value, fc)
		out := b.Local(l.tempName(), l.gt.GArrayPtr)
		totalLen := b.Add(b.GetField(recv, "len"), b.GetField(other, "len"))
		b.Assign(out, b.Call("g_array_sized_new",
			cir.BoolLit(false), cir.BoolLit(false), elemSize, totalLen))
		b.Stmt(b.Call("g_array_append_vals",
			out, b.GetField(recv, "data"), b.GetField(recv, "len")))
		b.Stmt(b.Call("g_array_append_vals",
			out, b.GetField(other, "data"), b.GetField(other, "len")))
		return out
	case "forEach":
		forEachFn := l.lowerExpr(b, args[0].Value, fc)
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
		b.While(b.Lt(i, b.GetField(recv, "len")),
			cir.B(func(b *cir.Builder) {
				b.Stmt(b.CallPtr(forEachFn, b.Index(dataPtr, i)))
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return nil
	}
	l.diags.Warnf(Pos{}, "array method %q not supported", method)
	return cir.NullPtr()
}

// lowerDynArrayMethod needs the Pos from the MethodCallExpr.
// We thread it via the helper below.
var _ = (*MethodCallExpr)(nil) // suppress unused lint

func (l *Lowerer) lowerStringMethod(
	b *cir.Builder, recv cir.Expr, rt *VString,
	method string, args []*Arg, fc *funcCtx,
) cir.Expr {
	switch method {
	case "delete":
		if rt.Mutable {
			b.Stmt(b.Call("g_string_free", recv, cir.BoolLit(true)))
		}
		return nil
	}
	return cir.NullPtr()
}

func (l *Lowerer) lowerClassMethod(
    b *cir.Builder, recv cir.Expr, className, method string, args []*Arg, fc *funcCtx,
) cir.Expr {
    // Native-interface class: emit a plain C call — no ClassName__ prefix, no receiver arg.
    if l.nativeClasses[className] {
        var callArgs []cir.Expr
        for _, a := range args {
            callArgs = append(callArgs, l.lowerExpr(b, a.Value, fc))
        }
        return b.Call(method, callArgs...)
    }
    switch method {
    case "delete":
        b.Stmt(b.Call(className+"__deinit", recv))
        b.Stmt(b.Call("free", recv))
        return nil
    case "new":
        return recv
    }
    cFuncName := className + "__" + method
    var callArgs []cir.Expr
    callArgs = append(callArgs, recv)
    for _, a := range args {
        callArgs = append(callArgs, l.lowerExpr(b, a.Value, fc))
    }
    return b.Call(cFuncName, callArgs...)
}

func (l *Lowerer) lowerStructMethod(
	b *cir.Builder, recv cir.Expr, structName, method string, args []*Arg, fc *funcCtx,
) cir.Expr {
	cFuncName := structName + "__" + method
	var callArgs []cir.Expr
	callArgs = append(callArgs, recv)
	for _, a := range args {
		callArgs = append(callArgs, l.lowerExpr(b, a.Value, fc))
	}
	return b.Call(cFuncName, callArgs...)
}

func (l *Lowerer) lowerFieldExpr(b *cir.Builder, e *FieldExpr, fc *funcCtx) cir.Expr {
	recv := l.lowerExpr(b, e.Recv, fc)
	recvType := e.Recv.GetVType()

	switch rt := recvType.(type) {
	case *VDynArray:
		if e.Field == "length" {
			return b.GetField(recv, "len")
		}
	case *VString:
		if e.Field == "length" && rt.Mutable {
			return b.GetField(recv, "len")
		}
	case *VClass:
		return b.GetField(recv, e.Field) // ptr->field
	case *VStruct:
		return b.DotField(recv, e.Field) // val.field
	}
	return b.GetField(recv, e.Field) // fallback
}

func (l *Lowerer) lowerIndexExpr(b *cir.Builder, e *IndexExpr, fc *funcCtx) cir.Expr {
	recv := l.lowerExpr(b, e.Recv, fc)
	idx := l.lowerExpr(b, e.Index, fc)
	recvType := e.Recv.GetVType()

	switch rt := recvType.(type) {
	case *VDynArray:
		elemCIR := l.vtypeToCIR(rt.Elem)
		if elemCIR == nil {
			elemCIR = l.vtypeToCIRFallback(rt.Elem)
		}
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, "data"))
		return b.Index(dataPtr, idx)
	case *VFixedArray:
		return b.Index(recv, idx)
	}
	return b.Index(recv, idx)
}

func (l *Lowerer) lowerStructLit(b *cir.Builder, e *StructLitExpr, fc *funcCtx) cir.Expr {
	st, ok := l.structTypes[e.TypeName]
	if !ok {
		l.diags.Errorf(e.Pos, "struct %q not registered", e.TypeName)
		return cir.NullPtr()
	}
	var fields []cir.FieldInit
	for _, f := range e.Fields {
		fields = append(fields, cir.FieldInit{
			Field: f.Name,
			Value: l.lowerExpr(b, f.Value, fc),
		})
	}
	return cir.CompoundLit(st, cir.InitStruct(fields...))
}

func (l *Lowerer) lowerArrayLit(b *cir.Builder, e *ArrayLitExpr, fc *funcCtx) cir.Expr {
	if len(e.Elems) == 0 {
		return cir.NullPtr()
	}
	// Fixed array literal: const T arr[] = { ... }
	vt := e.GetVType()
	fa, ok := vt.(*VFixedArray)
	if !ok {
		return cir.NullPtr()
	}
	elemCIR := l.vtypeToCIR(fa.Elem)
	if elemCIR == nil {
		elemCIR = l.vtypeToCIRFallback(fa.Elem)
	}
	var elems []cir.Expr
	for _, el := range e.Elems {
		elems = append(elems, l.lowerExpr(b, el, fc))
	}
	arrType := cir.Array(elemCIR, len(elems))
	return cir.CompoundLit(arrType, cir.InitArray(elems...))
}

func (l *Lowerer) lowerMapLit(b *cir.Builder, e *MapLitExpr, fc *funcCtx) cir.Expr {
	// Emit equivalent of g_hash_table_new(NULL, NULL)
	hashTable := b.Local(l.tempName(), cir.VoidPtr)
	b.Assign(hashTable, b.Call("g_hash_table_new", cir.NullPtr(), cir.NullPtr()))
	
	for _, f := range e.Fields {
		k := l.lowerExpr(b, f.Key, fc)
		v := l.lowerExpr(b, f.Value, fc)
		
		// Note: GLib requires pointers. If you are inserting raw ints, 
		// you might need to use GINT_TO_POINTER macros or allocate memory here depending 
		// on how deep you want your GLib integration to go for scalars.
		b.Stmt(b.Call("g_hash_table_insert", hashTable, k, v))
	}
	return hashTable
}

func (l *Lowerer) lowerArrayCtor(b *cir.Builder, e *ArrayCtorExpr, fc *funcCtx) cir.Expr {
	vt := e.GetVType()
	switch t := vt.(type) {
	case *VFixedArray:
		return l.lowerFixedArrayCtor(b, e, t, fc)
	case *VDynArray:
		return l.lowerDynArrayCtor(b, e, t, fc)
	}
	return cir.NullPtr()
}

func (l *Lowerer) lowerFixedArrayCtor(b *cir.Builder, e *ArrayCtorExpr, t *VFixedArray, fc *funcCtx) cir.Expr {
	elemCIR := l.vtypeToCIR(t.Elem)
	if elemCIR == nil {
		elemCIR = l.vtypeToCIRFallback(t.Elem)
	}
	size := t.Size
	if size <= 0 {
		size = 1 // fallback
	}
	arrRef := b.Local(l.tempName(), cir.Array(elemCIR, size))
	// Determine fill value.
	fillVal := cir.IntLit(0) // default zero
	for _, a := range e.Args {
		if a.Label == "repeating" {
			fillVal = l.lowerExpr(b, a.Value, fc)
		}
	}
	arrSize := b.Mul(cir.UIntLit(uint64(size)), b.SizeOf(elemCIR))
	b.Stmt(b.Call("memset", arrRef, fillVal, arrSize))
	return arrRef
}

func (l *Lowerer) lowerDynArrayCtor(b *cir.Builder, e *ArrayCtorExpr, t *VDynArray, fc *funcCtx) cir.Expr {
	elemCIR := l.vtypeToCIR(t.Elem)
	if elemCIR == nil {
		elemCIR = l.vtypeToCIRFallback(t.Elem)
	}
	elemSize := b.Cast(cir.UInt32, b.SizeOf(elemCIR))

	// Check for capacity argument.
	for _, a := range e.Args {
		if a.Label == "capacity" {
			cap := l.lowerExpr(b, a.Value, fc)
			return b.Call("g_array_sized_new",
				cir.BoolLit(false), cir.BoolLit(false),
				elemSize, b.Cast(cir.UInt32, cap),
			)
		}
	}
	return b.Call("g_array_new",
		cir.BoolLit(false), cir.BoolLit(false), elemSize)
}

// ─── Type utilities ───────────────────────────────────────────────────────────

// vtypeToCIR returns the CIR type for vt, or nil if the lowerer must handle it.
func (l *Lowerer) vtypeToCIR(vt VType) cir.Type {
	if vt == nil {
		return nil
	}
	switch t := vt.(type) {
	case *VStruct:
		if st, ok := l.structTypes[t.Name]; ok {
			return st
		}
		return nil
	case *VClass:
		if st, ok := l.classTypes[t.Name]; ok {
			return cir.Ptr(st)
		}
		return nil
	case *VEnum:
		return cir.Int32
	case *VDynArray:
		return l.gt.GArrayPtr
	case *VString:
		if t.Mutable {
			return l.gt.GStringPtr
		}
		return cir.ConstPtr(cir.Char)
	case *VTypeAlias:
		return l.vtypeToCIR(t.Underlying)
	}
	return vt.CIRType()
}

// resolveTypeExprVType converts a syntactic TypeExpr to a VType by name lookup.
func (l *Lowerer) resolveTypeExprVType(te TypeExpr) VType {
	if te == nil {
		return &VVoid{}
	}
	switch t := te.(type) {
	case *NamedTypeExpr:
		if bt, ok := BuiltinTypes[t.Name]; ok {
			return bt
		}
		if _, ok := l.structTypes[t.Name]; ok {
			return &VStruct{Name: t.Name}
		}
		if _, ok := l.classTypes[t.Name]; ok {
			return &VClass{Name: t.Name}
		}
		return &VUnknown{Name: t.Name}
	case *PointerTypeExpr:
		return &VPointer{Elem: l.resolveTypeExprVType(t.Elem), IsConst: t.IsConst}
	case *ArrayTypeExpr:
		return &VDynArray{Elem: l.resolveTypeExprVType(t.Elem)}
	case *OptionalTypeExpr:
		return &VOptional{Elem: l.resolveTypeExprVType(t.Elem)}
	}
	return &VVoid{}
}

// lowerTestFuncDecl emits a void function whose return statements are rewritten
// to printf calls so the value appears on stdout for the test runner to capture.
func (l *Lowerer) lowerTestFuncDecl(fn *FuncDecl) {
	l.ensurePrintf()

	def := l.mod.Func(fn.Name, cir.Returns(cir.Void))
	if fn.Body == nil {
		return
	}

	fc := newFuncCtx()
	fc.isTestFunc = true
	for _, p := range fn.Params {
		fc.params[p.Name] = true
	}

	def.Body(func(b *cir.Builder) {
		l.lowerBlock(b, fn.Body, fc)
		// Each ReturnStmt already emits a void return via lowerStmt.
		// Only add a fallback return for functions that fall off the end
		// without any return statement at all.
		if !blockEndsWithReturn(fn.Body) {
			fc.emitDefers(b)
			b.Return()
		}
	})
}

// injectTestMain emits:
//
//	int main() { <testEntryFunc>(); return 0; }
func (l *Lowerer) injectTestMain() {
	entryName := l.testEntryFunc
	def := l.mod.Func("main", cir.Returns(cir.Int32))
	def.Body(func(b *cir.Builder) {
		b.Stmt(b.Call(entryName))
		b.ReturnVal(cir.IntLit(0))
	})
}

// ensurePrintf registers printf as a C extern and adds <stdio.h> to the module
// if they are not already present.
func (l *Lowerer) ensurePrintf() {
	if l.mod.LookupExtern("printf") != nil {
		return
	}
	l.mod.Include("<stdio.h>")
	l.mod.Extern("printf",
		cir.Returns(cir.Int32),
		cir.Param("fmt", cir.ConstPtr(cir.Char)),
		cir.Variadic,
	)
}

// printfFormatFor returns a printf format specifier appropriate for vt.
func (l *Lowerer) printfFormatFor(vt VType) string {
	switch t := vt.(type) {
	case *VInt:
		if t.Bits == 64 {
			if t.Signed {
				return "%lld"
			}
			return "%llu"
		}
		if t.Signed {
			return "%d"
		}
		return "%u"
	case *VFloat:
		return "%g"
	case *VBool:
		return "%d"
	case *VChar:
		return "%c"
	case *VString:
		return "%s"
	default:
		return "%d"
	}
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func vtypeBaseName(vt VType) string {
	switch t := vt.(type) {
	case *VStruct:
		return t.Name
	case *VClass:
		return t.Name
	case *VEnum:
		return t.Name
	}
	return ""
}

// extractTypeName pulls the bare type name from a TypeExpr.
func extractTypeName(te TypeExpr) string {
	if te == nil {
		return ""
	}
	switch t := te.(type) {
	case *NamedTypeExpr:
		return t.Name
	case *PointerTypeExpr:
		return extractTypeName(t.Elem)
	}
	return ""
}

// blockEndsWithReturn reports whether the last statement of blk is a return.
func blockEndsWithReturn(blk *BlockStmt) bool {
	if blk == nil || len(blk.Stmts) == 0 {
		return false
	}
	_, ok := blk.Stmts[len(blk.Stmts)-1].(*ReturnStmt)
	return ok
}

// isFallthroughCase reports whether a switch case ends with fallthrough.
func isFallthroughCase(c *SwitchCase) bool {
	if len(c.Body) == 0 {
		return false
	}
	_, ok := c.Body[len(c.Body)-1].(*FallthroughStmt)
	return ok
}

// Silence unused import warnings for helpers used in some branches only.
var _ = fmt.Sprintf