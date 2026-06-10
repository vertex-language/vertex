package compiler

import (
	"fmt"

	cir "github.com/vertex-language/ir/c"
)

// ─── Lowerer struct — add pointerReceiverMethods ──────────────────────────────

type Lowerer struct {
	diags        *Diagnostics
	mod          *cir.Module
	arrStruct    *cir.StructType // arrays_Array — see runtime/arrays/arrays.vs
	arrStructPtr cir.Type        // *arrays_Array
	pkg          string

	structTypes            map[string]*cir.StructType
	classTypes             map[string]*cir.StructType
	classDecls             map[string]*ClassDecl
	enumTypes              map[string]*VEnum
	nativeClasses          map[string]bool
	userFuncs              map[string]bool
	pointerReceiverMethods map[string]bool

	testEntryFunc  string
	testEntryVType VType

	tempSeq int

	optionalTypes map[string]cir.Type
}

func NewLowerer(diags *Diagnostics, mod *cir.Module) *Lowerer {
	arrStruct, arrStructPtr := setupArraysRuntime(mod)
	setupVtxMaps(mod)
	return &Lowerer{
		diags:                  diags,
		mod:                    mod,
		arrStruct:              arrStruct,
		arrStructPtr:           arrStructPtr,
		structTypes:            make(map[string]*cir.StructType),
		classTypes:             make(map[string]*cir.StructType),
		classDecls:             make(map[string]*ClassDecl),
		enumTypes:              make(map[string]*VEnum),
		nativeClasses:          make(map[string]bool),
		userFuncs:              make(map[string]bool),
		pointerReceiverMethods: make(map[string]bool),
		optionalTypes:          make(map[string]cir.Type),
	}
}

// cFuncName returns the C-level name for a Vertex free function.
// "main" is always kept as-is — it is the C entry point.
// All other user functions are prefixed with the package name so Vertex
// symbols can never collide with C stdlib or linker-reserved names.
func (l *Lowerer) cFuncName(name string) string {
	if name == "main" || l.pkg == "" {
		return name
	}
	return l.pkg + "_" + name
}

// NEW: isPointerVType reports whether the type naturally lowers to a C pointer.
// If true, we can safely use NULL for .none and skip generating a wrapper struct.
func (l *Lowerer) isPointerVType(vt VType) bool {
	switch vt.(type) {
	case *VClass, *VDynArray, *VString, *VPointer, *VMap:
		return true
	}
	return false
}

// NEW: wrapOptional implicitly wraps a value into an Optional compound literal if needed.
func (l *Lowerer) wrapOptional(targetVT VType, valVT VType, val cir.Expr) cir.Expr {
	optVT, isOptTarget := targetVT.(*VOptional)
	_, isOptVal := valVT.(*VOptional) // FIX: Check if the value is already an optional
	_, isNilVal := valVT.(*VNil)
	
	if isOptTarget {
		// FIX: If the RHS is already an optional, don't double-wrap it!
		if isOptVal {
			return val
		}

		if l.isPointerVType(optVT.Elem) {
			// Pointer type: no struct needed, just return the value or NULL
			if isNilVal { return cir.NullPtr() }
			return val
		}

		// Value type: We must emit a compound struct literal
		optCT := l.vtypeToCIRFallback(targetVT)

		if isNilVal {
			// .none -> { .has_value = false }
			return cir.CompoundLit(optCT, cir.InitStruct(
				cir.FieldInit{Field: "has_value", Value: cir.BoolLit(false)},
			))
		}

		// .some(val) -> { .has_value = true, .value = val }
		return cir.CompoundLit(optCT, cir.InitStruct(
			cir.FieldInit{Field: "has_value", Value: cir.BoolLit(true)},
			cir.FieldInit{Field: "value", Value: val},
		))
	}
	return val
}

// cTypeName returns the C-level name for a Vertex-defined struct / class / enum.
func (l *Lowerer) cTypeName(name string) string {
	if l.pkg == "" {
		return name
	}
	return l.pkg + "_" + name
}

// cMethodName returns the C-level name for a receiver or associated method.
// e.g. pkg="main", type="Dog", method="bark" → "main_Dog__bark"
func (l *Lowerer) cMethodName(typeName, methodName string) string {
	return l.cTypeName(typeName) + "__" + methodName
}

// LowerFile drives both lowering passes over file.
func (l *Lowerer) LowerFile(file *File) {
	l.pkg = file.Package // captured once; used by all cFuncName/cTypeName calls
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
	st := cir.Struct(l.cTypeName(d.Name), fields...)
	l.mod.RegisterType(st)
	l.structTypes[d.Name] = st // map key stays original name for all internal lookups
}

func (l *Lowerer) registerClass(d *ClassDecl) {
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
	st := cir.Struct(l.cTypeName(d.Name), fields...)
	l.mod.RegisterType(st)
	l.classTypes[d.Name] = st
	l.classDecls[d.Name] = d
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

// Update registerEnum to store the VEnum with its AST declaration
func (l *Lowerer) registerEnum(d *EnumDecl) {
	// Reconstruct the VEnum so we have it for explicit type hint resolution
	l.enumTypes[d.Name] = &VEnum{Name: d.Name, Decl: d}
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
		if st, ok := l.structTypes[t.Name]; ok {
			return st
		}
		if st, ok := l.classTypes[t.Name]; ok {
			return cir.Ptr(st)
		}
		return cir.VoidPtr
	case *FixedArrayTypeExpr:
		elem := l.resolvedFieldCIRType(t.Elem)
		if elem == nil {
			elem = cir.Void
		}
		if t.Size != nil {
			if il, ok := t.Size.(*IntLitExpr); ok && il.Value > 0 {
				return cir.Array(elem, int(il.Value))
			}
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
		return l.arrStructPtr
	}
	return cir.VoidPtr
}

// ─── Pass 2: function lowering ────────────────────────────────────────────────

func (l *Lowerer) lowerFunctions(file *File) {
	hasInit := make(map[string]bool)
	hasDeinit := make(map[string]bool)

	for _, decl := range file.Decls {
		if fn, ok := decl.(*FuncDecl); ok {
			l.userFuncs[fn.Name] = true

			if fn.Receiver != nil {
				recvTypeName := extractTypeName(fn.Receiver.Type)
				if fn.Name == "init" {
					hasInit[recvTypeName] = true
				} else if fn.Name == "deinit" {
					hasDeinit[recvTypeName] = true
				}
				// Track which methods have pointer receivers so call sites
				// can pass &recv instead of recv.
				if fn.Receiver.IsPtr {
					l.pointerReceiverMethods[recvTypeName+"__"+fn.Name] = true
				}
			}
		}
	}

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ClassDecl:
			if l.nativeClasses[d.Name] {
				continue
			}
			if !hasInit[d.Name] {
				l.lowerDefaultClassInit(d)
			}
			if !hasDeinit[d.Name] {
				l.lowerDefaultClassDeinit(d)
			}
		case *EnumDecl:
			l.lowerEnumHelpers(d)
		case *FuncDecl:
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

	if l.testEntryFunc != "" {
		l.injectTestMain()
	}
}

// lowerDefaultClassInit automatically synthesizes a constructor that assigns all fields in order.
func (l *Lowerer) lowerDefaultClassInit(d *ClassDecl) {
	cName := l.cMethodName(d.Name, "init")
	st := l.classTypes[d.Name]
	ptrType := cir.Ptr(st)

	opts := []cir.FuncOpt{cir.Returns(cir.Void), cir.Param("this", ptrType)}

	var fields []*ClassMember
	for _, m := range d.Members {
		if m.IsField {
			fields = append(fields, m)
			ft := l.resolvedFieldCIRType(m.Type)
			if ft == nil { ft = cir.VoidPtr }
			opts = append(opts, cir.Param(m.Name, ft))
		}
	}

	def := l.mod.Func(cName, opts...)
	def.Body(func(b *cir.Builder) {
		this := b.Param("this")
		for _, f := range fields {
			ft := l.resolvedFieldCIRType(f.Type)
			if ft == nil { ft = cir.VoidPtr }
			val := b.Param(f.Name)
            
			// FIXED: Use StructStore instead of Assign(GetField)
			b.StructStore(cir.Deref(this, st), val, cir.Step(st, f.Name, ft))
		}
		b.Return()
	})
}

// lowerDefaultClassDeinit automatically synthesizes a destructor that cleans up known heap fields.
func (l *Lowerer) lowerDefaultClassDeinit(d *ClassDecl) {
	cName := l.cMethodName(d.Name, "deinit")
	st := l.classTypes[d.Name]
	ptrType := cir.Ptr(st)

	opts := []cir.FuncOpt{cir.Returns(cir.Void), cir.Param("this", ptrType)}
	def := l.mod.Func(cName, opts...)

	def.Body(func(b *cir.Builder) {
		this := b.Param("this")
		for _, m := range d.Members {
			if !m.IsField {
				continue
			}
			vt := l.resolveTypeExprVType(m.Type)

			if _, ok := vt.(*VDynArray); ok {
				arrPtrType := l.arrStructPtr
				fieldVal := b.StructLoad(cir.Deref(this, st), cir.Step(st, m.Name, arrPtrType))
				b.If(b.Neq(fieldVal, cir.NullPtr()), cir.B(func(b *cir.Builder) {
					b.Stmt(b.Call("arrays_free", fieldVal))
				}))
			}
			// strings are const char* — not owned, no cleanup needed
		}
		b.Return()
	})
}

// NEW: Generates the getter and constructor C helpers for an enum's raw values.
func (l *Lowerer) lowerEnumHelpers(d *EnumDecl) {
	var rawTypeVT VType = &VInt{Bits: 32, Signed: true}
	if d.RawType != nil {
		rawTypeVT = l.resolveTypeExprVType(d.RawType)
	}
	ct := l.vtypeToCIR(rawTypeVT)
	if ct == nil {
		ct = l.vtypeToCIRFallback(rawTypeVT)
	}

	// 1. rawValue helper: pkg_EnumName_rawValue(int32) -> raw_type
	def1 := l.mod.Func(l.cFuncName(d.Name+"_rawValue"), cir.Returns(ct), cir.Param("val", cir.Int32))
	def1.Body(func(b *cir.Builder) {
		val := b.Param("val")
		var cases []cir.SwitchCase
		var currentInt int64 = 0
		for _, c := range d.Cases {
			var retVal cir.Expr
			if c.RawValue != nil {
				if il, ok := c.RawValue.(*IntLitExpr); ok {
					currentInt = il.Value
				}
			}
			
			if _, isStr := rawTypeVT.(*VString); isStr {
				if c.RawValue != nil {
					if sl, ok := c.RawValue.(*StringLitExpr); ok {
						retVal = l.mod.StringLit(l.tempName(), sl.Value)
					}
				} else {
					retVal = l.mod.StringLit(l.tempName(), c.Name) // default to case name string
				}
			} else {
				retVal = b.Cast(ct, cir.IntLit(currentInt))
			}
			
			cases = append(cases, cir.Case(cir.IntLit(currentInt), cir.B(func(b *cir.Builder) {
				b.ReturnVal(retVal)
			})))
			currentInt++
		}
		
		cases = append(cases, cir.Default(cir.B(func(b *cir.Builder) {
			if _, isStr := rawTypeVT.(*VString); isStr {
				b.ReturnVal(cir.NullPtr())
			} else {
				b.ReturnVal(b.Cast(ct, cir.IntLit(0)))
			}
		})))
		b.Switch(val, cases...)
	})

	// 2. from_rawValue helper: pkg_EnumName_from_rawValue(raw_type) -> opt_Enum
	enumVT := l.resolveTypeExprVType(&NamedTypeExpr{Name: d.Name})
	optVT := &VOptional{Elem: enumVT}
	optCT := l.vtypeToCIRFallback(optVT) // Resolve the struct opt_Enum type

	def2 := l.mod.Func(l.cFuncName(d.Name+"_from_rawValue"), cir.Returns(optCT), cir.Param("raw", ct))
	def2.Body(func(b *cir.Builder) {
		raw := b.Param("raw")
		_, isStr := rawTypeVT.(*VString)
		if isStr {
			if l.mod.LookupExtern("strcmp") == nil {
				l.mod.Extern("strcmp", cir.Returns(cir.Int32), cir.Param("s1", cir.ConstPtr(cir.Char)), cir.Param("s2", cir.ConstPtr(cir.Char)))
			}
		}

		var currentInt int64 = 0
		var conds []struct {
			cond cir.Expr
			val  int64
		}
		
		for _, c := range d.Cases {
			if c.RawValue != nil {
				if il, ok := c.RawValue.(*IntLitExpr); ok {
					currentInt = il.Value
				}
			}
			
			if isStr {
				var strVal string
				if c.RawValue != nil {
					if sl, ok := c.RawValue.(*StringLitExpr); ok {
						strVal = sl.Value
					}
				} else {
					strVal = c.Name
				}
				lit := l.mod.StringLit(l.tempName(), strVal)
				cond := b.Eq(b.Call("strcmp", raw, lit), cir.IntLit(0))
				conds = append(conds, struct{cond cir.Expr; val int64}{cond, currentInt})
			} else {
				cond := b.Eq(raw, b.Cast(ct, cir.IntLit(currentInt)))
				conds = append(conds, struct{cond cir.Expr; val int64}{cond, currentInt})
			}
			currentInt++
		}

		var build func(i int) *cir.Block
		build = func(i int) *cir.Block {
			if i >= len(conds) {
				return cir.B(func(b *cir.Builder) { 
					// Return .none (has_value = false)
					b.ReturnVal(cir.CompoundLit(optCT, cir.InitStruct(
						cir.FieldInit{Field: "has_value", Value: cir.BoolLit(false)},
					))) 
				})
			}
			return cir.B(func(b *cir.Builder) {
				b.IfElse(conds[i].cond, cir.B(func(b *cir.Builder) {
					// Return .some (has_value = true, value = EnumCase)
					b.ReturnVal(cir.CompoundLit(optCT, cir.InitStruct(
						cir.FieldInit{Field: "has_value", Value: cir.BoolLit(true)},
						cir.FieldInit{Field: "value", Value: b.Cast(cir.Int32, cir.IntLit(conds[i].val))},
					)))
				}), build(i+1))
			})
		}
		b.Inline(build(0))
	})
}

func (l *Lowerer) lowerGlobalVar(d *VarDecl) {
	if len(d.Binding.Names) == 0 {
		return
	}
	// Value may be nil for: var buf: [T; N]  (fixed array without initializer)
	var vt VType
	if d.Value != nil {
		vt = d.Value.GetVType()
	}
	if d.TypeHint != nil {
		vt = l.resolveTypeExprVType(d.TypeHint)
	}
	if vt == nil {
		return
	}
	ct := l.vtypeToCIR(vt)
	if ct == nil {
		ct = cir.VoidPtr
	}
	name := d.Binding.Names[0]
	if d.IsLet && d.Value != nil {
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

// ─── lowerFuncDecl — mark pointer receiver in funcCtx ────────────────────────

func (l *Lowerer) lowerFuncDecl(fn *FuncDecl) {
	if fn.Qualifier == FuncQualTest {
		l.lowerTestFuncDecl(fn)
		return
	}
	if fn.Qualifier != FuncQualNone {
		l.diags.Warnf(fn.Pos, "function qualifier %v is not yet supported; ignoring", fn.Qualifier)
	}

	var cName string
	if fn.Receiver != nil {
		recvTypeName := extractTypeName(fn.Receiver.Type)
		cName = l.cMethodName(recvTypeName, fn.Name)
	} else {
		cName = l.cFuncName(fn.Name)
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
		// When the receiver is *T the CIR param is a pointer; record it so
		// field-access and assignment paths can dereference it correctly.
		if fn.Receiver.IsPtr {
			fc.ptrParams[fn.Receiver.Name] = true
		}
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

// ─── funcCtx — add ptrParams ──────────────────────────────────────────────────

type funcCtx struct {
	params    map[string]bool
	ptrParams map[string]bool // receiver params declared as *T (IsPtr=true)
	locals    map[string]cir.Expr
	defers    []func(*cir.Builder)
}

func newFuncCtx() *funcCtx {
	return &funcCtx{
		params:    make(map[string]bool),
		ptrParams: make(map[string]bool),
		locals:    make(map[string]cir.Expr),
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
		if st.Value != nil {
			// 1. Evaluate the return expression BEFORE defers
			valExpr := l.lowerExpr(b, st.Value, fc)
			
			// Get the CIR type for the temporary variable
			vt := st.Value.GetVType()
			ct := l.vtypeToCIR(vt)
			if ct == nil {
				ct = l.vtypeToCIRFallback(vt)
			}
			
			// 2. Store it in a safe temporary local
			tmp := b.Local(l.tempName(), ct)
			b.Assign(tmp, valExpr)
			
			// 3. NOW it is safe to emit the defers (freeing memory)
			fc.emitDefers(b)
			
			// 4. Return the safely stored value
			b.ReturnVal(tmp)
		} else {
			fc.emitDefers(b)
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
	// Value may be nil for: var buf: [T; N]  (fixed array without initializer)
	var vtype VType
	if d.Value != nil {
		vtype = d.Value.GetVType()
	}
	if d.TypeHint != nil {
		vtype = l.resolveTypeExprVType(d.TypeHint)
	}
	if vtype == nil {
		vtype = &VUnknown{}
	}

	// ── Fixed array with no initializer: var buf: [T; N] ─────────────────────
	// Declare the local and emit memset for zero-fill, then return early.
	if d.Value == nil {
		for _, name := range d.Binding.Names {
			ref := l.declareLocal(b, name, vtype, d.IsLet)
			fc.locals[name] = ref
			if fa, isFixed := vtype.(*VFixedArray); isFixed && fa.Size > 0 {
				elemCIR := l.vtypeToCIR(fa.Elem)
				if elemCIR == nil {
					elemCIR = l.vtypeToCIRFallback(fa.Elem)
				}
				arrSize := b.Mul(cir.UIntLit(uint64(fa.Size)), b.SizeOf(elemCIR))
				b.Stmt(b.Call("memset", ref, cir.IntLit(0), arrSize))
			}
		}
		return
	}

	// NOTE: With the updated grammar, [T] is always dynamic and [T; N] is always
	// fixed. The TypeHint (already resolved above) is authoritative — do NOT
	// coerce VDynArray to VFixedArray based on the value being an array literal.

	for i, name := range d.Binding.Names {
		_ = i
		ref := l.declareLocal(b, name, vtype, d.IsLet)
		fc.locals[name] = ref

		if arrLit, ok := d.Value.(*ArrayLitExpr); ok {
			// ── Dynamic array from literal ─────────────────────────────────────
			// var x = [e1, e2, ...]       → vtype is VDynArray (inferred by resolver)
			// var x: [T] = [e1, e2, ...]  → vtype is VDynArray (annotated)
			if da, isDyn := vtype.(*VDynArray); isDyn {
				elemCIR := l.vtypeToCIR(da.Elem)
				if elemCIR == nil {
					elemCIR = l.vtypeToCIRFallback(da.Elem)
				}
				elemSize := b.Cast(cir.UInt32, b.SizeOf(elemCIR))
				b.Assign(ref, b.Call("arrays_newWithCapacity", elemSize,
					b.Cast(cir.UInt32, cir.UIntLit(uint64(len(arrLit.Elems))))))
				for _, elem := range arrLit.Elems {
					val := l.lowerExpr(b, elem, fc)
					tmp := b.Local(l.tempName(), elemCIR)
					b.Assign(tmp, b.Cast(elemCIR, val))
					b.Stmt(b.Call("arrays_push", ref, b.AddrOf(tmp)))
				}
				continue
			}
			// ── Fixed array from literal ───────────────────────────────────────
			// let arr = [e1, e2, ...]          → vtype is VFixedArray (inferred)
			// var arr: [T; N] = [e1, e2, ...]  → vtype is VFixedArray (annotated)
			if fa, isFixed := vtype.(*VFixedArray); isFixed {
				elemCIR := l.vtypeToCIR(fa.Elem)
				if elemCIR == nil {
					elemCIR = l.vtypeToCIRFallback(fa.Elem)
				}
				for j, elem := range arrLit.Elems {
					val := l.lowerExpr(b, elem, fc)
					if elemCIR != nil {
						val = b.Cast(elemCIR, val)
					}
					b.Assign(b.Index(ref, cir.IntLit(int64(j)), elemCIR), val)
				}
				continue
			}
		}

		initExpr := l.lowerExpr(b, d.Value, fc)

		if initExpr != nil {
			initExpr = l.wrapOptional(vtype, d.Value.GetVType(), initExpr)

			if _, isOpt := vtype.(*VOptional); !isOpt {
				if _, isStruct := vtype.(*VStruct); !isStruct {
					if targetCT := l.vtypeToCIR(vtype); targetCT != nil {
						initExpr = b.Cast(targetCT, initExpr)
					}
				}
			}
			b.Assign(ref, initExpr)
		}
	}
}

func (l *Lowerer) vtypeToCIRFallback(vt VType) cir.Type {
	switch t := vt.(type) {
	case *VFixedArray:
		if t.Size <= 0 {
			return cir.VoidPtr
		}
		elemCT := l.vtypeToCIR(t.Elem)
		if elemCT == nil {
			elemCT = l.vtypeToCIRFallback(t.Elem)
		}
		if elemCT == nil {
			return cir.VoidPtr
		}
		return cir.Array(elemCT, t.Size)
	case *VDynArray:
		return l.arrStructPtr
	case *VString:
		return cir.ConstPtr(cir.Char)
	case *VStruct:
		if st, ok := l.structTypes[t.Name]; ok {
			return st
		}
	case *VMap:
		return cir.VoidPtr
	case *VClass:
		if st, ok := l.classTypes[t.Name]; ok {
			return cir.Ptr(st)
		}
	case *VEnum:
		return cir.Int32
	case *VOptional:
		if l.isPointerVType(t.Elem) {
			inner := l.vtypeToCIR(t.Elem)
			if inner == nil {
				inner = l.vtypeToCIRFallback(t.Elem)
			}
			return inner
		}
		typeName := vtypeName(t.Elem)
		if cached, ok := l.optionalTypes[typeName]; ok {
			return cached
		}
		elemCT := l.vtypeToCIR(t.Elem)
		if elemCT == nil {
			elemCT = l.vtypeToCIRFallback(t.Elem)
		}
		structName := l.cTypeName("opt_" + typeName)
		st := cir.Struct(structName,
			cir.Field("has_value", cir.Bool),
			cir.Field("value", elemCT),
		)
		l.mod.RegisterType(st)
		l.optionalTypes[typeName] = st
		return st
	case *VResult:
		return cir.VoidPtr
	case *VInt:
		ct := t.CIRType()
		if ct == nil {
			panic(fmt.Sprintf("vtypeToCIRFallback: VInt{%d,%v} has no CIR type", t.Bits, t.Signed))
		}
		return ct
	case *VFloat:
		ct := t.CIRType()
		if ct == nil {
			panic(fmt.Sprintf("vtypeToCIRFallback: VFloat{%d} has no CIR type", t.Bits))
		}
		return ct
	case *VBool:
		return cir.Bool
	case *VChar:
		return cir.Char
	}
	return cir.VoidPtr
}

func (l *Lowerer) declareLocal(b *cir.Builder, name string, vt VType, isLet bool) cir.Expr {
	ct := l.vtypeToCIR(vt)
	if ct == nil {
		ct = l.vtypeToCIRFallback(vt)
	}
	ref := b.Local(name, ct)
	return ref
}

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
		return l.arrStructPtr
	case *VString:
		return cir.ConstPtr(cir.Char)
	case *VTypeAlias:
		return l.vtypeToCIR(t.Underlying)
	}
	return vt.CIRType()
}

// lowerAssign lowers an assignment statement.
// When the LHS is a value-type struct field chain, StructStore is emitted
// so the MIR never sees a StructLoadExpr on the write side of an AssignStmt.
func (l *Lowerer) lowerAssign(b *cir.Builder, st *AssignStmt, fc *funcCtx) {
	// ── Struct and Class field write: LHS is a field chain ─────────────────────
	if fe, ok := st.LHS.(*FieldExpr); ok {
		if baseAST, steps, chainOK := l.collectStructChain(fe); chainOK {
			baseExpr := l.lowerExpr(b, baseAST, fc)
			rhs := l.lowerExpr(b, st.RHS, fc)

			// Class receivers are always pointers — dereference before StructStore.
			if vc, isClass := baseAST.GetVType().(*VClass); isClass {
				baseExpr = cir.Deref(baseExpr, l.classTypes[vc.Name])
			}
			// Pointer receiver struct params (func (c: *Counter) …) are also
			// lowered as pointers in CIR even though VType shows the bare struct.
			if id, ok2 := baseAST.(*IdentExpr); ok2 && fc.ptrParams[id.Name] {
				if vs, ok3 := baseAST.GetVType().(*VStruct); ok3 {
					if structST, ok4 := l.structTypes[vs.Name]; ok4 {
						baseExpr = cir.Deref(baseExpr, structST)
					}
				}
			}

			switch st.Op {
			case OpAssign:
				b.StructStore(baseExpr, rhs, steps...)
			default:
				cur := b.StructLoad(baseExpr, steps...)
				var newVal cir.Expr
				switch st.Op {
				case OpAddAssign:
					newVal = b.Add(cur, rhs)
				case OpSubAssign:
					newVal = b.Sub(cur, rhs)
				case OpMulAssign:
					newVal = b.Mul(cur, rhs)
				case OpDivAssign:
					newVal = b.Div(cur, rhs)
				case OpModAssign:
					newVal = b.Mod(cur, rhs)
				default:
					newVal = rhs
				}
				b.StructStore(baseExpr, newVal, steps...)
			}
			return
		}
	}

	// ── Original pointer / scalar / class assignment logic ────────────────────
	lhs := l.lowerExpr(b, st.LHS, fc)
	rhs := l.lowerExpr(b, st.RHS, fc)
	if lhs == nil || rhs == nil {
		return
	}

	if ptr, isPtr := st.LHS.GetVType().(*VPointer); isPtr {
		elemCIR := l.vtypeToCIR(ptr.Elem)
		if elemCIR == nil {
			elemCIR = cir.Int32
		}
		storeLHS := cir.Deref(lhs, elemCIR)
		if _, rhsIsPtr := st.RHS.GetVType().(*VPointer); rhsIsPtr {
			if id, ok := st.RHS.(*IdentExpr); ok && fc.params[id.Name] {
				rhs = cir.Deref(rhs, elemCIR)
			}
		}
		rhs = l.wrapOptional(ptr.Elem, st.RHS.GetVType(), rhs)
		switch st.Op {
		case OpAssign:
			b.Assign(storeLHS, rhs)
		case OpAddAssign:
			b.Assign(storeLHS, b.Add(cir.Deref(lhs, elemCIR), rhs))
		case OpSubAssign:
			b.Assign(storeLHS, b.Sub(cir.Deref(lhs, elemCIR), rhs))
		case OpMulAssign:
			b.Assign(storeLHS, b.Mul(cir.Deref(lhs, elemCIR), rhs))
		case OpDivAssign:
			b.Assign(storeLHS, b.Div(cir.Deref(lhs, elemCIR), rhs))
		case OpModAssign:
			b.Assign(storeLHS, b.Mod(cir.Deref(lhs, elemCIR), rhs))
		}
		return
	}

	rhs = l.wrapOptional(st.LHS.GetVType(), st.RHS.GetVType(), rhs)
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
		val := l.lowerExpr(b, c.Expr, fc)
		valVT := c.Expr.GetVType()
		optVT, isOpt := valVT.(*VOptional)

		if !isOpt || l.isPointerVType(optVT.Elem) {
			// Pointer optional or non-optional: standard != NULL check.
			tmp := b.Local(l.tempName(), l.vtypeToCIRFallback(valVT))
			b.Assign(tmp, val)
			fc.locals[c.Name] = tmp
			cond = b.Neq(tmp, cir.NullPtr())
		} else {
			// Value-type struct optional — use the opt_T struct type for correct offsets.
			optCT := l.vtypeToCIRFallback(valVT)
			optSt, _ := optCT.(*cir.StructType)
			tmp := b.Local(l.tempName(), optCT)
			b.Assign(tmp, val)

			elemCT := l.vtypeToCIRFallback(optVT.Elem)
			boundLocal := b.Local(c.Name, elemCT)
			b.Assign(boundLocal, b.DotField(tmp, optSt, "value", elemCT))
			fc.locals[c.Name] = boundLocal

			cond = b.DotField(tmp, optSt, "has_value", cir.Bool)
		}
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
	// Range literal: lo..<hi or lo...hi.
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
		post := &cir.AssignExpr{
			LHS: iRef,
			RHS: b.Add(iRef, b.Cast(cir.Int32, cir.IntLit(1))),
		}
		b.For(nil, cond, post, cir.B(func(b *cir.Builder) {
			l.lowerBlock(b, st.Body, fc)
		}))
		return
	}

	iterType := st.Iter.GetVType()
	iter := l.lowerExpr(b, st.Iter, fc)
	if iter == nil {
		return
	}
	switch it := iterType.(type) {
	case *VFixedArray:
		if it.Size <= 0 {
			l.diags.Warnf(st.Pos, "for-in over fixed array with unknown size; skipping")
			return
		}
		elemCIR := l.vtypeToCIR(it.Elem)
		if elemCIR == nil {
			elemCIR = l.vtypeToCIRFallback(it.Elem)
		}
		iRef := b.Local(l.tempName(), cir.UInt32)
		b.Assign(iRef, b.Cast(cir.UInt32, cir.UIntLit(0)))
		post := &cir.AssignExpr{
			LHS: iRef,
			RHS: b.Add(iRef, b.Cast(cir.UInt32, cir.UIntLit(1))),
		}
		limit := b.Cast(cir.UInt32, cir.UIntLit(uint64(it.Size)))
		b.For(nil, b.Lt(iRef, limit), post, cir.B(func(b *cir.Builder) {
			elemRef := b.Local(st.Var, elemCIR)
			fc.locals[st.Var] = elemRef
			b.Assign(elemRef, b.Index(iter, iRef, elemCIR))
			l.lowerBlock(b, st.Body, fc)
		}))

	case *VDynArray:
		elemCIR := l.vtypeToCIR(it.Elem)
		if elemCIR == nil {
			elemCIR = l.vtypeToCIRFallback(it.Elem)
		}
		as := l.arrStruct
		iRef := b.Local(l.tempName(), cir.UInt32)
		b.Assign(iRef, b.Cast(cir.UInt32, cir.UIntLit(0)))
		arrLen := b.GetField(iter, as, "len", cir.UInt32)
		b.While(b.Lt(iRef, arrLen), cir.B(func(b *cir.Builder) {
			dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(iter, as, "data", cir.VoidPtr))
			elemRef := b.Local(st.Var, elemCIR)
			fc.locals[st.Var] = elemRef
			b.Assign(elemRef, b.Index(dataPtr, iRef, elemCIR))
			l.lowerBlock(b, st.Body, fc)
			b.Assign(iRef, b.Add(iRef, b.Cast(cir.UInt32, cir.UIntLit(1))))
		}))

	default:
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
		return cir.IntLit(int64(e.Value))
	case *StringLitExpr:
		return l.mod.StringLit(l.tempName(), e.Value)
	case *NilLitExpr:
		return cir.NullPtr()

	case *IdentExpr:
		return fc.varRef(b, e.Name)

	case *DotEnumExpr:
		return l.lowerDotEnum(b, e)

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

	case *CastExpr:
		// Handles all 'expr as typeExpr' forms: pointer↔pointer reinterpret,
		// integer widening/truncation, float↔int conversion, pointer↔integer.
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

func (l *Lowerer) lowerDotEnum(b *cir.Builder, e *DotEnumExpr) cir.Expr {
	return l.enumCaseExpr(b, e.GetVType(), e.Case)
}

func (l *Lowerer) enumCaseExpr(b *cir.Builder, subjType VType, caseName string) cir.Expr {
	if ev, ok := subjType.(*VEnum); ok {
		var currentInt int64 = 0
		for _, c := range ev.Decl.Cases {
			if c.RawValue != nil {
				if il, ok := c.RawValue.(*IntLitExpr); ok {
					currentInt = il.Value
				}
			}
			if c.Name == caseName {
				return cir.IntLit(currentInt)
			}
			currentInt++
		}
	}
	return cir.IntLit(0)
}

func (l *Lowerer) lowerBinaryExpr(b *cir.Builder, e *BinaryExpr, fc *funcCtx) cir.Expr {
	if e.Op == BinNilCoalesce {
		var collect func(expr Expr) []Expr
		collect = func(expr Expr) []Expr {
			if bin, ok := expr.(*BinaryExpr); ok && bin.Op == BinNilCoalesce {
				return append(collect(bin.Left), bin.Right)
			}
			return []Expr{expr}
		}
		operands := collect(e)

		finalCT := l.vtypeToCIR(e.GetVType())
		if finalCT == nil {
			finalCT = l.vtypeToCIRFallback(e.GetVType())
		}
		result := b.Local(l.tempName(), finalCT)

		var build func(idx int) *cir.Block
		build = func(idx int) *cir.Block {
			if idx == len(operands)-1 {
				return cir.B(func(b *cir.Builder) {
					val := l.lowerExpr(b, operands[idx], fc)
					b.Assign(result, b.Cast(finalCT, val))
				})
			}
			return cir.B(func(b *cir.Builder) {
				val := l.lowerExpr(b, operands[idx], fc)
				vt := operands[idx].GetVType()
				optVT, isOpt := vt.(*VOptional)

				if isOpt && !l.isPointerVType(optVT.Elem) {
					// Value-type optional — use the opt_T struct type for correct field offsets.
					optCIR := l.vtypeToCIRFallback(vt)
					optSt, _ := optCIR.(*cir.StructType)
					tmp := b.Local(l.tempName(), optCIR)
					b.Assign(tmp, val)

					hasValueField := b.DotField(tmp, optSt, "has_value", cir.Bool)
					cond := b.Eq(b.Cast(cir.Int32, hasValueField), cir.IntLit(1))
					thenBlk := cir.B(func(b *cir.Builder) {
						unwrapped := b.DotField(tmp, optSt, "value", finalCT)
						b.Assign(result, b.Cast(finalCT, unwrapped))
					})
					b.IfElse(cond, thenBlk, build(idx+1))
				} else {
					// Pointer-type optional or plain fallback.
					tmp := b.Local(l.tempName(), l.vtypeToCIRFallback(vt))
					b.Assign(tmp, val)
					cond := b.Neq(tmp, cir.NullPtr())
					thenBlk := cir.B(func(b *cir.Builder) {
						b.Assign(result, b.Cast(finalCT, tmp))
					})
					b.IfElse(cond, thenBlk, build(idx+1))
				}
			})
		}
		b.Inline(build(0))
		return result
	}

	left := l.lowerExpr(b, e.Left, fc)
	right := l.lowerExpr(b, e.Right, fc)
	switch e.Op {
	case BinAdd:         return b.Add(left, right)
	case BinSub:         return b.Sub(left, right)
	case BinMul:         return b.Mul(left, right)
	case BinDiv:         return b.Div(left, right)
	case BinMod:         return b.Mod(left, right)
	case BinShl:         return b.Shl(left, right)
	case BinShr:         return b.Shr(left, right)
	case BinBitAnd:      return b.And(left, right)
	case BinBitXor:      return b.Xor(left, right)
	case BinBitOr:       return b.Or(left, right)
	case BinEq:          return b.Eq(left, right)
	case BinNeq:         return b.Neq(left, right)
	case BinLt:          return b.Lt(left, right)
	case BinLte:         return b.Lte(left, right)
	case BinGt:          return b.Gt(left, right)
	case BinGte:         return b.Gte(left, right)
	case BinAnd:         return b.LogAnd(left, right)
	case BinOr:          return b.LogOr(left, right)
	case BinOverflowAdd: return b.Add(b.Cast(cir.UInt32, left), b.Cast(cir.UInt32, right))
	case BinOverflowSub: return b.Sub(b.Cast(cir.UInt32, left), b.Cast(cir.UInt32, right))
	case BinOverflowMul: return b.Mul(b.Cast(cir.UInt32, left), b.Cast(cir.UInt32, right))
	case BinIdentityEq:  return b.Eq(left, right)
	case BinIdentityNeq: return b.Neq(left, right)
	}
	return left
}

func (l *Lowerer) lowerUnaryExpr(b *cir.Builder, e *UnaryExpr, fc *funcCtx) cir.Expr {
	switch e.Op {
	case UnAddrOf:
		// The Vertex grammar gives 'as' higher precedence than '&', so the
		// user's "&x as *T" is parsed as "&(x as *T)" — i.e.:
		//   UnaryExpr{UnAddrOf, CastExpr{Value: x, TargetType: *T}}
		//
		// Taking the address of a cast result is not a valid C lvalue.
		// When the cast target is a pointer type, flip the order:
		//   &(x as *T)  →  (*T)(&x)
		if cast, ok := e.Operand.(*CastExpr); ok {
			targetVT := l.resolveTypeExprVType(cast.TargetType)
			if _, isPtr := targetVT.(*VPointer); isPtr {
				inner := l.lowerExpr(b, cast.Value, fc)
				addr := b.AddrOf(inner)
				ct := l.vtypeToCIR(targetVT)
				if ct == nil {
					ct = l.vtypeToCIRFallback(targetVT)
				}
				return b.Cast(ct, addr)
			}
		}
		op := l.lowerExpr(b, e.Operand, fc)
		return b.AddrOf(op)
	}

	op := l.lowerExpr(b, e.Operand, fc)
	switch e.Op {
	case UnNeg:
		return b.Neg(op)
	case UnNot:
		return b.Not(op)
	case UnBitNot:
		return b.BitNot(op)
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

// callReturnCIRType resolves a VType to its CIR representation for use as a
// call return type. Falls back to Void when the type cannot be mapped.
func (l *Lowerer) callReturnCIRType(vt VType) cir.Type {
	if vt == nil {
		return cir.Void
	}
	ct := l.vtypeToCIR(vt)
	if ct == nil {
		ct = l.vtypeToCIRFallback(vt)
	}
	if ct == nil {
		return cir.Void
	}
	return ct
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
		if _, isClass := l.classTypes[id.Name]; isClass {
			return l.lowerClassInstantiate(b, id.Name, e.Args, fc)
		}
		if l.nativeClasses[id.Name] {
			return l.lowerClassInstantiate(b, id.Name, e.Args, fc)
		}

		// Check for Enum init(rawValue:) constructor
		if opt, isOpt := e.GetVType().(*VOptional); isOpt {
			if ve, isEnum := opt.Elem.(*VEnum); isEnum && len(e.Args) == 1 && e.Args[0].Label == "rawValue" {
				argVal := l.lowerExpr(b, e.Args[0].Value, fc)
				ce := b.Call(l.cFuncName(ve.Name+"_from_rawValue"), argVal)
				if callNode, ok := ce.(*cir.CallExpr); ok {
					// Stamp the explicit struct opt_T instead of a pointer!
					callNode.Type = l.vtypeToCIRFallback(opt) 
				}
				return ce
			}
		}
	}

	fn := l.lowerExpr(b, e.Func, fc)
	var args []cir.Expr
	for _, a := range e.Args {
		args = append(args, l.lowerExpr(b, a.Value, fc))
	}

	retCIR := l.callReturnCIRType(e.GetVType())

	if id, ok := e.Func.(*IdentExpr); ok {
		cName := id.Name
		if l.userFuncs[id.Name] {
			cName = l.cFuncName(id.Name)
		}
		ce := b.Call(cName, args...)
		if callNode, ok := ce.(*cir.CallExpr); ok {
			callNode.Type = retCIR
		}
		return ce
	}
	ce := b.CallPtr(fn, args...)
	if callNode, ok := ce.(*cir.CallPtrExpr); ok {
		callNode.Type = retCIR
	}
	return ce
}

func (l *Lowerer) lowerClassInstantiate(b *cir.Builder, className string, args []*Arg, fc *funcCtx) cir.Expr {
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
	b.Stmt(b.Call(l.cMethodName(className, "init"), initArgs...))
	return tmp
}

func (l *Lowerer) lowerMethodCall(b *cir.Builder, e *MethodCallExpr, fc *funcCtx) cir.Expr {
	recv := l.lowerExpr(b, e.Recv, fc)
	recvType := e.Recv.GetVType()
	retCIR := l.callReturnCIRType(e.GetVType())

	var result cir.Expr
	switch rt := recvType.(type) {
	case *VDynArray:
		result = l.lowerDynArrayMethod(b, recv, rt, e.Method, e.Args, fc)
	case *VString:
		result = l.lowerStringMethod(b, recv, rt, e.Method, e.Args, fc)
	case *VClass:
		result = l.lowerClassMethod(b, recv, rt.Name, e.Method, e.Args, fc)
	case *VStruct:
		result = l.lowerStructMethod(b, recv, rt.Name, e.Method, e.Args, fc)
	default:
		if typeName := vtypeBaseName(recvType); typeName != "" {
			var args []cir.Expr
			args = append(args, recv)
			for _, a := range e.Args {
				args = append(args, l.lowerExpr(b, a.Value, fc))
			}
			result = b.Call(l.cMethodName(typeName, e.Method), args...)
		} else {
			l.diags.Warnf(e.Pos, "method %q on type %v not lowered", e.Method, recvType)
			return cir.NullPtr()
		}
	}

	// Stamp the resolved return type onto the CIR call node so the backend
	// can route float returns through XMM registers correctly.
	if result != nil {
		if ce, ok := result.(*cir.CallExpr); ok {
			ce.Type = retCIR
		}
	}
	return result
}

func (l *Lowerer) lowerDynArrayMethod(
	b *cir.Builder, recv cir.Expr, rt *VDynArray,
	method string, args []*Arg, fc *funcCtx,
) cir.Expr {
	elemCIR := l.vtypeToCIR(rt.Elem)
	if elemCIR == nil {
		elemCIR = l.vtypeToCIRFallback(rt.Elem)
	}
	elemSize := b.Cast(cir.UInt32, b.SizeOf(elemCIR))
	as := l.arrStruct

	switch method {
	case "push":
		val := l.lowerExpr(b, args[0].Value, fc)
		tmp := b.Local(l.tempName(), elemCIR)
		b.Assign(tmp, val)
		b.Stmt(b.Call("arrays_push", recv, b.AddrOf(tmp)))
		return nil

	case "unshift":
		val := l.lowerExpr(b, args[0].Value, fc)
		tmp := b.Local(l.tempName(), elemCIR)
		b.Assign(tmp, val)
		b.Stmt(b.Call("arrays_unshift", recv, b.AddrOf(tmp)))
		return nil

	case "pop":
		optVT := &VOptional{Elem: rt.Elem}
		optCT := l.vtypeToCIRFallback(optVT)
		result := b.Local(l.tempName(), optCT)
		b.Assign(result, cir.CompoundLit(optCT, cir.InitStruct(
			cir.FieldInit{Field: "has_value", Value: cir.BoolLit(false)},
		)))
		lenVal := b.GetField(recv, as, "len", cir.UInt32)
		b.If(b.Gt(lenVal, cir.UIntLit(0)),
			cir.B(func(b *cir.Builder) {
				newLen := b.Sub(lenVal, cir.UIntLit(1))
				dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
				// capture element BEFORE decrementing len
				elemLocal := b.Local(l.tempName(), elemCIR)
				b.Assign(elemLocal, b.Index(dataPtr, newLen, elemCIR))
				b.StructStore(cir.Deref(recv, as), newLen, cir.Step(as, "len", cir.UInt32))
				b.Assign(result, cir.CompoundLit(optCT, cir.InitStruct(
					cir.FieldInit{Field: "has_value", Value: cir.BoolLit(true)},
					cir.FieldInit{Field: "value", Value: elemLocal},
				)))
			}),
		)
		return result

	case "shift":
		optVT := &VOptional{Elem: rt.Elem}
		optCT := l.vtypeToCIRFallback(optVT)
		result := b.Local(l.tempName(), optCT)
		b.Assign(result, cir.CompoundLit(optCT, cir.InitStruct(
			cir.FieldInit{Field: "has_value", Value: cir.BoolLit(false)},
		)))
		lenVal := b.GetField(recv, as, "len", cir.UInt32)
		b.If(b.Gt(lenVal, cir.UIntLit(0)),
			cir.B(func(b *cir.Builder) {
				dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
				// capture element BEFORE removeAt shifts the array
				elemLocal := b.Local(l.tempName(), elemCIR)
				b.Assign(elemLocal, b.Index(dataPtr, cir.UIntLit(0), elemCIR))
				b.Stmt(b.Call("arrays_removeAt", recv, cir.UIntLit(0)))
				b.Assign(result, cir.CompoundLit(optCT, cir.InitStruct(
					cir.FieldInit{Field: "has_value", Value: cir.BoolLit(true)},
					cir.FieldInit{Field: "value", Value: elemLocal},
				)))
			}),
		)
		return result

	case "delete":
		b.Stmt(b.Call("arrays_free", recv))
		return nil

	case "fill":
		val := l.lowerExpr(b, args[0].Value, fc)
		byteSize := b.Mul(b.Cast(cir.UIntSize, b.GetField(recv, as, "len", cir.UInt32)), b.SizeOf(elemCIR))
		b.Stmt(b.Call("memset", b.GetField(recv, as, "data", cir.VoidPtr), val, byteSize))
		return nil

	case "reverse":
		lo := b.Local(l.tempName(), cir.UInt32)
		hi := b.Local(l.tempName(), cir.UInt32)
		tmp := b.Local(l.tempName(), elemCIR)
		b.Assign(lo, cir.UIntLit(0))
		b.Assign(hi, b.Sub(b.GetField(recv, as, "len", cir.UInt32), cir.UIntLit(1)))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
		b.While(b.Lt(lo, hi),
			cir.B(func(b *cir.Builder) {
				b.Assign(tmp, b.Index(dataPtr, lo, elemCIR))
				b.Assign(b.Index(dataPtr, lo, elemCIR), b.Index(dataPtr, hi, elemCIR))
				b.Assign(b.Index(dataPtr, hi, elemCIR), tmp)
				b.Assign(lo, b.Add(lo, cir.UIntLit(1)))
				b.Assign(hi, b.Sub(hi, cir.UIntLit(1)))
			}),
		)
		return nil

	case "sort":
		cmpFn := l.lowerExpr(b, args[0].Value, fc)
		b.Stmt(b.Call("arrays_sort", recv, cmpFn))
		return nil

	case "indexOf":
		val := l.lowerExpr(b, args[0].Value, fc)
		result := b.Local(l.tempName(), cir.Int32)
		b.Assign(result, cir.IntLit(-1))
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
		b.While(b.Lt(i, b.GetField(recv, as, "len", cir.UInt32)),
			cir.B(func(b *cir.Builder) {
				b.If(b.Eq(b.Index(dataPtr, i, elemCIR), val),
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
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
		b.While(b.Lt(i, b.GetField(recv, as, "len", cir.UInt32)),
			cir.B(func(b *cir.Builder) {
				b.If(b.Eq(b.Index(dataPtr, i, elemCIR), val),
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
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
		b.While(b.Lt(i, b.GetField(recv, as, "len", cir.UInt32)),
			cir.B(func(b *cir.Builder) {
				elem := b.AddrOf(b.Index(dataPtr, i, elemCIR))
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
		out := b.Local(l.tempName(), l.arrStructPtr)
		b.Assign(out, b.Call("arrays_newWithCapacity", elemSize, b.GetField(recv, as, "len", cir.UInt32)))
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
		b.While(b.Lt(i, b.GetField(recv, as, "len", cir.UInt32)),
			cir.B(func(b *cir.Builder) {
				mapped := b.Local(l.tempName(), elemCIR)
				b.Assign(mapped, b.CallPtr(mapFn, b.Index(dataPtr, i, elemCIR)))
				b.Stmt(b.Call("arrays_push", out, b.AddrOf(mapped)))
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return out

	case "filter":
		filterFn := l.lowerExpr(b, args[0].Value, fc)
		out := b.Local(l.tempName(), l.arrStructPtr)
		b.Assign(out, b.Call("arrays_new", elemSize))
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
		b.While(b.Lt(i, b.GetField(recv, as, "len", cir.UInt32)),
			cir.B(func(b *cir.Builder) {
				elem := b.Index(dataPtr, i, elemCIR)
				b.If(b.CallPtr(filterFn, elem),
					cir.B(func(b *cir.Builder) {
						tmp := b.Local(l.tempName(), elemCIR)
						b.Assign(tmp, elem)
						b.Stmt(b.Call("arrays_push", out, b.AddrOf(tmp)))
					}),
				)
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return out

	case "slice":
		start := l.lowerExpr(b, args[0].Value, fc)
		end := l.lowerExpr(b, args[1].Value, fc)
		out := b.Local(l.tempName(), l.arrStructPtr)
		b.Assign(out, b.Call("arrays_new", elemSize))
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, b.Cast(cir.UInt32, start))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
		b.While(
			b.LogAnd(b.Lt(i, b.Cast(cir.UInt32, end)), b.Lt(i, b.GetField(recv, as, "len", cir.UInt32))),
			cir.B(func(b *cir.Builder) {
				tmp := b.Local(l.tempName(), elemCIR)
				b.Assign(tmp, b.Index(dataPtr, i, elemCIR))
				b.Stmt(b.Call("arrays_push", out, b.AddrOf(tmp)))
				b.Assign(i, b.Add(i, cir.UIntLit(1)))
			}),
		)
		return out

	case "concat":
		other := l.lowerExpr(b, args[0].Value, fc)
		out := b.Local(l.tempName(), l.arrStructPtr)
		totalLen := b.Add(b.GetField(recv, as, "len", cir.UInt32), b.GetField(other, as, "len", cir.UInt32))
		b.Assign(out, b.Call("arrays_newWithCapacity", elemSize, totalLen))
		iRef := b.Local(l.tempName(), cir.UInt32)
		b.Assign(iRef, cir.UIntLit(0))
		recvData := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
		b.While(b.Lt(iRef, b.GetField(recv, as, "len", cir.UInt32)),
			cir.B(func(b *cir.Builder) {
				tmp := b.Local(l.tempName(), elemCIR)
				b.Assign(tmp, b.Index(recvData, iRef, elemCIR))
				b.Stmt(b.Call("arrays_push", out, b.AddrOf(tmp)))
				b.Assign(iRef, b.Add(iRef, cir.UIntLit(1)))
			}),
		)
		b.Assign(iRef, cir.UIntLit(0))
		otherData := b.Cast(cir.Ptr(elemCIR), b.GetField(other, as, "data", cir.VoidPtr))
		b.While(b.Lt(iRef, b.GetField(other, as, "len", cir.UInt32)),
			cir.B(func(b *cir.Builder) {
				tmp := b.Local(l.tempName(), elemCIR)
				b.Assign(tmp, b.Index(otherData, iRef, elemCIR))
				b.Stmt(b.Call("arrays_push", out, b.AddrOf(tmp)))
				b.Assign(iRef, b.Add(iRef, cir.UIntLit(1)))
			}),
		)
		return out

	case "forEach":
		forEachFn := l.lowerExpr(b, args[0].Value, fc)
		i := b.Local(l.tempName(), cir.UInt32)
		b.Assign(i, cir.UIntLit(0))
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, as, "data", cir.VoidPtr))
		b.While(b.Lt(i, b.GetField(recv, as, "len", cir.UInt32)),
			cir.B(func(b *cir.Builder) {
				b.Stmt(b.CallPtr(forEachFn, b.Index(dataPtr, i, elemCIR)))
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
	// Strings are now const char* — no heap management at the language level.
	l.diags.Warnf(Pos{}, "string method %q not supported on immutable strings", method)
	return cir.NullPtr()
}

func (l *Lowerer) lowerClassMethod(
	b *cir.Builder, recv cir.Expr, className, method string, args []*Arg, fc *funcCtx,
) cir.Expr {
	// Native-interface class: plain C call — no prefix, no receiver arg.
	if l.nativeClasses[className] {
		var callArgs []cir.Expr
		for _, a := range args {
			callArgs = append(callArgs, l.lowerExpr(b, a.Value, fc))
		}
		return b.Call(method, callArgs...)
	}
	switch method {
	case "delete":
		b.Stmt(b.Call(l.cMethodName(className, "deinit"), recv))
		b.Stmt(b.Call("free", recv))
		return nil
	case "new":
		return recv
	}
	var callArgs []cir.Expr
	callArgs = append(callArgs, recv)
	for _, a := range args {
		callArgs = append(callArgs, l.lowerExpr(b, a.Value, fc))
	}
	return b.Call(l.cMethodName(className, method), callArgs...)
}

func (l *Lowerer) lowerStructMethod(
	b *cir.Builder, recv cir.Expr, structName, method string, args []*Arg, fc *funcCtx,
) cir.Expr {
	var callArgs []cir.Expr
	// Pointer receiver methods expect *T — take the address of the local.
	// Value receiver methods expect T  — pass the value directly.
	if l.pointerReceiverMethods[structName+"__"+method] {
		callArgs = append(callArgs, b.AddrOf(recv))
	} else {
		callArgs = append(callArgs, recv)
	}
	for _, a := range args {
		callArgs = append(callArgs, l.lowerExpr(b, a.Value, fc))
	}
	return b.Call(l.cMethodName(structName, method), callArgs...)
}

// ─── Struct chain helpers ─────────────────────────────────────────────────────

// collectStructChain walks a FieldExpr tree upward and, when every level
// accesses a registered value-type struct (not a class pointer), returns:
//   - baseAST: the root AST expression (the local variable or param)
//   - steps:   ordered slice of FieldStep, outermost last
//   - ok:      false if any level is not a value-type struct (caller should fallback)
//
// Example for e = FieldExpr{Recv: FieldExpr{Recv: IdentExpr{"l"}, Field:"start"}, Field:"x"}
// where l:Line, start:Vec2, x:float32:
//   baseAST = IdentExpr{"l"}
//   steps   = [Step(LineType,"start",Vec2Type), Step(Vec2Type,"x",Float32)]
func (l *Lowerer) collectStructChain(e *FieldExpr) (baseAST Expr, steps []cir.FieldStep, ok bool) {
    var reversed []cir.FieldStep
    curr := e

    for {
        recvVT  := curr.Recv.GetVType()
        fieldVT := curr.GetVType()

        structST, isStruct := l.lookupStructCIRType(recvVT)
        if !isStruct {
            return nil, nil, false
        }

        // Always go through vtypeToCIR first; it returns the package-level
        // singleton (cir.Float32, cir.Int32, etc.) via vt.CIRType().
        // vtypeToCIRFallback must only be reached for aggregate types
        // (structs, arrays) where vtypeToCIR legitimately returns nil.
        fieldCT := l.vtypeToCIR(fieldVT)
        if fieldCT == nil {
            fieldCT = l.vtypeToCIRFallback(fieldVT)
        }
        if fieldCT == nil {
            // Unresolvable type — bail out and let the caller use the
            // GetField fallback path rather than emitting a broken chain.
            return nil, nil, false
        }

        reversed = append(reversed, cir.Step(structST, curr.Field, fieldCT))

        inner, innerIsField := curr.Recv.(*FieldExpr)
        if !innerIsField {
            baseAST = curr.Recv
            break
        }
        if _, ok2 := l.lookupStructCIRType(inner.Recv.GetVType()); !ok2 {
            baseAST = curr.Recv
            break
        }
        curr = inner
    }

    steps = make([]cir.FieldStep, len(reversed))
    for i, s := range reversed {
        steps[len(reversed)-1-i] = s
    }

    // Final sanity: the last step's FieldType is what becomes StructLoadExpr.Type.
    // It must be a scalar cir type recognised by the MIR's typeClass switch,
    // not a struct or pointer.
    if len(steps) > 0 {
        last := steps[len(steps)-1].FieldType
        if _, isStruct := last.(*cir.StructType); isStruct {
            // The chain ends on a struct field, not a scalar.
            // StructLoadExpr would carry a struct type, which is valid but
            // uncommon; flag it so the MIR author knows to handle it.
            _ = last // acceptable — MIR handles struct-typed loads as memcpy
        }
    }

    return baseAST, steps, true
}

// lookupStructCIRType returns the CIR StructType for a VStruct, or (nil, false)
// for anything else.  Classes are excluded because they use pointer access (->)
// not value access (.).
func (l *Lowerer) lookupStructCIRType(vt VType) (*cir.StructType, bool) {
	if vs, ok := vt.(*VStruct); ok {
		if st, ok2 := l.structTypes[vs.Name]; ok2 {
			return st, true
		}
	}
	// ADDED: Treat classes like structs so they can use StructLoad/StructStore
	if vc, ok := vt.(*VClass); ok {
		if st, ok2 := l.classTypes[vc.Name]; ok2 {
			return st, true
		}
	}
	return nil, false
}

// lowerFieldExpr lowers a field access expression.
// Value-type struct chains (l.start.x) use StructLoad; everything else falls
// back to the existing GetField / DotField paths.
func (l *Lowerer) lowerFieldExpr(b *cir.Builder, e *FieldExpr, fc *funcCtx) cir.Expr {
	recvType := e.Recv.GetVType()

	if ve, ok := recvType.(*VEnum); ok {
		if e.Field == "rawValue" {
			recv := l.lowerExpr(b, e.Recv, fc)
			return b.Call(l.cFuncName(ve.Name+"_rawValue"), recv)
		}
		return l.enumCaseExpr(b, ve, e.Field)
	}

	if baseAST, steps, ok := l.collectStructChain(e); ok {
		baseExpr := l.lowerExpr(b, baseAST, fc)
		if vc, isClass := baseAST.GetVType().(*VClass); isClass {
			baseExpr = cir.Deref(baseExpr, l.classTypes[vc.Name])
		}
		if id, ok2 := baseAST.(*IdentExpr); ok2 && fc.ptrParams[id.Name] {
			if vs, ok3 := baseAST.GetVType().(*VStruct); ok3 {
				if structST, ok4 := l.structTypes[vs.Name]; ok4 {
					baseExpr = cir.Deref(baseExpr, structST)
				}
			}
		}
		return b.StructLoad(baseExpr, steps...)
	}

	recv := l.lowerExpr(b, e.Recv, fc)

	fieldVT := e.GetVType()
	fieldCT := l.vtypeToCIR(fieldVT)
	if fieldCT == nil {
		fieldCT = l.vtypeToCIRFallback(fieldVT)
	}

	switch rt := recvType.(type) {
	case *VDynArray:
		if e.Field == "length" {
			return b.Cast(cir.Int32, b.GetField(recv, l.arrStruct, "len", cir.UInt32))
		}
	case *VString:
		if e.Field == "length" {
			return b.Cast(cir.Int32, b.Call("strlen", recv))
		}
	case *VClass:
		if st, ok := l.classTypes[rt.Name]; ok {
			return b.StructLoad(cir.Deref(recv, st), cir.Step(st, e.Field, fieldCT))
		}
	}

	return b.GetField(recv, nil, e.Field, fieldCT)
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
		dataPtr := b.Cast(cir.Ptr(elemCIR), b.GetField(recv, l.arrStruct, "data", cir.VoidPtr))
		return b.Index(dataPtr, idx, elemCIR)
	case *VFixedArray:
		elemCIR := l.vtypeToCIR(rt.Elem)
		if elemCIR == nil {
			elemCIR = l.vtypeToCIRFallback(rt.Elem)
		}
		return b.Index(recv, idx, elemCIR)
	}
	return b.Index(recv, idx, cir.VoidPtr)
}

func (l *Lowerer) lowerStructLit(b *cir.Builder, e *StructLitExpr, fc *funcCtx) cir.Expr {
	// ── Class literal: Counter{value: 10} → malloc + init(this, fields…) ─────
	if st, ok := l.classTypes[e.TypeName]; ok {
		ptrType := cir.Ptr(st)
		tmp := b.Local(l.tempName(), ptrType)
		b.Assign(tmp, b.Cast(ptrType, b.Call("malloc", b.SizeOf(st))))

		// Lower all field values up-front, keyed by name.
		fieldVals := make(map[string]cir.Expr, len(e.Fields))
		for _, f := range e.Fields {
			fieldVals[f.Name] = l.lowerExpr(b, f.Value, fc)
		}

		// Emit init args in declaration order so they match the
		// auto-generated init(this, field0, field1, …) signature.
		initArgs := []cir.Expr{tmp}
		if cd, ok2 := l.classDecls[e.TypeName]; ok2 {
			for _, m := range cd.Members {
				if !m.IsField {
					continue
				}
				if v, ok3 := fieldVals[m.Name]; ok3 {
					initArgs = append(initArgs, v)
				} else {
					initArgs = append(initArgs, cir.NullPtr())
				}
			}
		} else {
			// classDecls unavailable (native class): use literal order as fallback.
			for _, f := range e.Fields {
				initArgs = append(initArgs, fieldVals[f.Name])
			}
		}

		b.Stmt(b.Call(l.cMethodName(e.TypeName, "init"), initArgs...))
		return tmp
	}

	// ── Struct literal (value type) ───────────────────────────────────────────
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
	vt := e.GetVType()
	fa, ok := vt.(*VFixedArray)
	if !ok {
		return cir.NullPtr()
	}
	elemCIR := l.vtypeToCIR(fa.Elem)
	if elemCIR == nil {
		elemCIR = l.vtypeToCIRFallback(fa.Elem)
	}
	// Allocate a temp local and assign each element by index.
	// cir.CompoundLit is never used for arrays: C does not allow compound
	// literals of array type as assignment RHS, and MIR rejects them.
	arrType := cir.Array(elemCIR, len(e.Elems))
	tmp := b.Local(l.tempName(), arrType)
	for i, el := range e.Elems {
		val := l.lowerExpr(b, el, fc)
		if elemCIR != nil {
			val = b.Cast(elemCIR, val)
		}
		b.Assign(b.Index(tmp, cir.IntLit(int64(i)), elemCIR), val)
	}
	return tmp
}

func (l *Lowerer) lowerMapLit(b *cir.Builder, e *MapLitExpr, fc *funcCtx) cir.Expr {
	l.diags.Warnf(e.Pos, "map literals are a stub; full map support is TBD")
	result := b.Local(l.tempName(), cir.VoidPtr)
	b.Assign(result, b.Call("v_map_new", cir.UIntLit(8), cir.UIntLit(8)))
	return result
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

	for _, a := range e.Args {
		if a.Label == "capacity" {
			cap := l.lowerExpr(b, a.Value, fc)
			return b.Call("arrays_newWithCapacity", elemSize, b.Cast(cir.UInt32, cap))
		}
	}
	return b.Call("arrays_new", elemSize)
}

// Update resolveTypeExprVType to handle Enums
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
		if ev, ok := l.enumTypes[t.Name]; ok {
			return ev
		}
		return &VUnknown{Name: t.Name}
	case *FixedArrayTypeExpr:
		elem := l.resolveTypeExprVType(t.Elem)
		size := -1
		if t.Size != nil {
			if il, ok := t.Size.(*IntLitExpr); ok {
				size = int(il.Value)
			}
		}
		return &VFixedArray{Elem: elem, Size: size}
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
	var retCIR cir.Type = cir.Int32
	var retVT VType = &VInt{Bits: 32, Signed: true}

	if fn.RetType != nil {
		if exp, ok := fn.RetType.(*ExpectedTypeExpr); ok && exp.ReturnType != nil {
			vt := l.resolveTypeExprVType(exp.ReturnType)
			if ct := l.vtypeToCIR(vt); ct != nil {
				retCIR = ct
				retVT = vt
			}
		}
	}

	l.testEntryVType = retVT

	def := l.mod.Func(l.cFuncName(fn.Name), cir.Returns(retCIR))
	if fn.Body == nil {
		return
	}

	fc := newFuncCtx()
	for _, p := range fn.Params {
		fc.params[p.Name] = true
	}

	def.Body(func(b *cir.Builder) {
		l.lowerBlock(b, fn.Body, fc)
		if !blockEndsWithReturn(fn.Body) {
			fc.emitDefers(b)
			b.Return()
		}
	})
}

// injectTestMain emits:
func (l *Lowerer) injectTestMain() {
	l.ensurePrintf()
	// Use the prefixed C name when calling the test entry function.
	entryCName := l.cFuncName(l.testEntryFunc)

	vt := l.testEntryVType
	if vt == nil {
		vt = &VInt{Bits: 32, Signed: true}
	}
	fmtStr := l.printfFormatFor(vt) + "\n"

	retCIR := l.vtypeToCIR(vt)
	if retCIR == nil {
		retCIR = cir.Int32
	}

	def := l.mod.Func("main", cir.Returns(cir.Int32))
	def.Body(func(b *cir.Builder) {
		tmp := b.Local(l.tempName(), retCIR)

		callExpr := b.Call(entryCName)
		if ce, ok := callExpr.(*cir.CallExpr); ok {
			ce.Type = retCIR
		}

		b.Assign(tmp, callExpr)
		fmtLit := l.mod.StringLit(l.tempName(), fmtStr)
		b.Stmt(b.Call("printf", fmtLit, tmp))
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
	// Unwrap optional — ?? and if-let always yield the inner type at runtime.
	if opt, ok := vt.(*VOptional); ok {
		return l.printfFormatFor(opt.Elem)
	}
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
		return "%f" // FIX: Changed from "%g" to "%f" to preserve trailing zeros for the test runner
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