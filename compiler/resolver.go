package compiler

import (
	"fmt"
	"strconv"
	"strings"
)

// ─────────────────────────────────────────────────────────────────────────────
// Resolver — name resolution and type inference pass
//
// Two sub-passes per file:
//   1. collectDecls – populates the package scope with all top-level names.
//   2. resolveDecls – resolves types and walks function bodies.
// ─────────────────────────────────────────────────────────────────────────────

type Resolver struct {
	diags *Diagnostics
	pkg   *Scope 
	file  *File

	// Generics tracking
	genericFuncs   map[string]*FuncDecl
	genericStructs map[string]*StructDecl
	instFuncs      map[string]bool
	instStructs    map[string]bool
}

func NewResolver(diags *Diagnostics, pkgScope *Scope) *Resolver {
	return &Resolver{
		diags:          diags,
		pkg:            pkgScope,
		genericFuncs:   make(map[string]*FuncDecl),
		genericStructs: make(map[string]*StructDecl),
		instFuncs:      make(map[string]bool),
		instStructs:    make(map[string]bool),
	}
}

// ResolveFile resolves all declarations in file and returns the annotated *File.
func (r *Resolver) ResolveFile(file *File) *File {
	r.file = file
	r.collectDecls(file)
	r.resolveDecls(file)
	return file
}

// ─── Pass 1: collect top-level names ─────────────────────────────────────────

func (r *Resolver) collectDecls(file *File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *StructDecl:
			if len(d.TypeParams) > 0 {
				r.genericStructs[d.Name] = d
				continue // Do NOT put templates in standard symbol table
			}
			r.pkg.Define(&Symbol{
				Name: d.Name, Kind: SymStruct,
				Type: &VStruct{Name: d.Name, Decl: d},
				Decl: d,
			})
		case *ClassDecl:
			r.pkg.Define(&Symbol{
				Name: d.Name, Kind: SymClass,
				Type: &VClass{Name: d.Name, Decl: d},
				Decl: d,
			})
		case *EnumDecl:
			rawType := VType(&VInt{Bits: 32, Signed: true})
			if d.RawType != nil {
				rawType = r.resolveTypeExpr(d.RawType, r.pkg)
			}
			enumType := &VEnum{Name: d.Name, RawType: rawType, Decl: d}
			r.pkg.Define(&Symbol{Name: d.Name, Kind: SymEnum, Type: enumType, Decl: d})
			for _, c := range d.Cases {
				r.pkg.Define(&Symbol{Name: d.Name + "." + c.Name, Kind: SymEnumCase, Type: enumType, Decl: d})
				r.pkg.Define(&Symbol{Name: c.Name, Kind: SymEnumCase, Type: enumType, Decl: d})
			}
		case *TypeAliasDecl:
			t := r.resolveTypeExpr(d.Type, r.pkg)
			r.pkg.Define(&Symbol{
				Name: d.Name, Kind: SymTypeAlias,
				Type: &VTypeAlias{Name: d.Name, Underlying: t},
				Decl: d, IsConst: true,
			})
		case *FuncDecl:
			if len(d.TypeParams) > 0 {
				r.genericFuncs[d.Name] = d
				continue // Do NOT put templates in standard symbol table
			}
			if d.Receiver == nil {
				r.pkg.Define(&Symbol{Name: d.Name, Kind: SymFunc, Decl: d})
			}
		case *VarDecl:
			for _, name := range d.Binding.Names {
				r.pkg.Define(&Symbol{Name: name, Kind: SymVar, Decl: d, IsConst: d.IsLet})
			}
		}
	}
}

// ─── Pass 2: resolve types and bodies ────────────────────────────────────────

func (r *Resolver) resolveDecls(file *File) {
	// Crucial: Iterating by index rather than range allows us to resolve 
	// newly appended generic clones dynamically at the end of the file!
	for i := 0; i < len(file.Decls); i++ {
		decl := file.Decls[i]
		switch d := decl.(type) {
		case *FuncDecl:
			if len(d.TypeParams) == 0 {
				r.resolveFunc(d)
			}
		case *VarDecl:
			r.resolveVarDeclGlobal(d)
		case *StructDecl:
			if len(d.TypeParams) == 0 {
				r.resolveStruct(d)
			}
		case *ClassDecl:
			r.resolveClass(d)
		}
	}
}

// Restore resolveFunc and resolveStruct to normal! Generic definitions never get here.
func (r *Resolver) resolveFunc(fn *FuncDecl) {
	fnScope := NewScope(r.pkg)

	if fn.Receiver != nil {
		recvType := r.resolveTypeExpr(fn.Receiver.Type, fnScope)
		if pt, ok := recvType.(*VPointer); ok {
			fn.Receiver.IsPtr = true
			recvType = pt.Elem
		}
		fn.Receiver.Type = typeExprFor(recvType)
		fnScope.Define(&Symbol{Name: fn.Receiver.Name, Kind: SymParam, Type: recvType})
	}

	for _, p := range fn.Params {
		pt := r.resolveTypeExpr(p.Type, fnScope)
		fnScope.Define(&Symbol{Name: p.Name, Kind: SymParam, Type: pt})
	}

	var retVType VType = &VVoid{}
	if fn.RetType != nil {
		retVType = r.resolveTypeExpr(fn.RetType, fnScope)
	}

	if fn.Body != nil {
		r.resolveBlock(fn.Body, fnScope, retVType)
	}
}

func (r *Resolver) resolveStruct(d *StructDecl) {
	for _, f := range d.Fields {
		r.resolveTypeExpr(f.Type, r.pkg)
	}
}

func (r *Resolver) resolveClass(d *ClassDecl) {
	for _, m := range d.Members {
		if m.IsField {
			r.resolveTypeExpr(m.Type, r.pkg)
		}
	}
}

func (r *Resolver) resolveVarDeclGlobal(d *VarDecl) {
	scope := r.pkg
	// Value may be nil for: var buf: [T; N]  (no initializer required)
	var vtype VType
	if d.Value != nil {
		vtype = r.resolveExpr(d.Value, scope)
	} else {
		vtype = &VVoid{}
	}
	if d.TypeHint != nil {
		vtype = r.resolveTypeExpr(d.TypeHint, scope)
	}
	if s, ok := vtype.(*VString); ok {
		s.Mutable = !d.IsLet
	}
	for _, name := range d.Binding.Names {
		if sym, ok := scope.LookupLocal(name); ok {
			sym.Type = vtype
			sym.IsConst = d.IsLet
		}
	}
}

// ─── Block / statement resolution ────────────────────────────────────────────

func (r *Resolver) resolveBlock(blk *BlockStmt, scope *Scope, retType VType) {
	for _, s := range blk.Stmts {
		r.resolveStmt(s, scope, retType)
	}
}

func (r *Resolver) resolveStmt(s Stmt, scope *Scope, retType VType) {
	switch st := s.(type) {
	case *LocalDeclStmt:
		r.resolveLocalDecl(st.Decl, scope)
	case *IfStmt:
		r.resolveIfStmt(st, scope, retType)
	case *WhileStmt:
		inner := NewScope(scope)
		r.resolveExpr(st.Cond, scope)
		r.resolveBlock(st.Body, inner, retType)
	case *ForInStmt:
		inner := NewScope(scope)
		iterType := r.resolveExpr(st.Iter, scope)
		var elemType VType = &VUnknown{Name: "element"}
		switch it := iterType.(type) {
		case *VDynArray:
			elemType = it.Elem
		case *VFixedArray:
			elemType = it.Elem
		case *VRange:
			elemType = it.Elem
		}
		inner.Define(&Symbol{Name: st.Var, Kind: SymVar, Type: elemType})
		r.resolveBlock(st.Body, inner, retType)
	case *SwitchStmt:
		subjType := r.resolveExpr(st.Subj, scope) // CHANGED: capture subjType
		for _, c := range st.Cases {
			inner := NewScope(scope)
			for _, p := range c.Patterns {
				switch pat := p.(type) {
				case *ExprPattern:
					r.resolveExpr(pat.Expr, scope)
					// NEW: Propagate switch subject type down to enum shorthands
					if dot, ok := pat.Expr.(*DotEnumExpr); ok {
						if _, isUnknown := dot.GetVType().(*VUnknown); isUnknown {
							dot.SetVType(subjType)
						}
					}
				case *ResultOkPattern:
					inner.Define(&Symbol{Name: pat.Bind, Kind: SymVar})
				case *ResultErrPattern:
					inner.Define(&Symbol{Name: pat.Bind, Kind: SymVar})
				}
			}
			for _, stmt := range c.Body {
				r.resolveStmt(stmt, inner, retType)
			}
		}
	case *ReturnStmt:
		if st.Value != nil {
			valType := r.resolveExpr(st.Value, scope)
			// NEW: Propagate expected return type down to enum shorthands
			if dot, ok := st.Value.(*DotEnumExpr); ok {
				if _, isUnknown := dot.GetVType().(*VUnknown); isUnknown {
					dot.SetVType(retType)
					valType = retType
				}
			}
			// Guard: Prevent returning values from void functions
			if _, isVoid := retType.(*VVoid); isVoid {
				r.diags.Errorf(st.Pos, "void function cannot return a value")
			} else if _, isTestExp := retType.(*VExpected); !isTestExp && !valType.Equal(retType) {
				// Guard: Ensure the returned value matches the declared return type,
				// UNLESS it's a test function with an Expected(...) annotation.
				r.diags.Errorf(st.Pos, "type mismatch: expected %s, got %s", retType, valType)
			}
		} else {
			// Guard: Prevent empty returns in non-void functions
			if _, isVoid := retType.(*VVoid); !isVoid {
				r.diags.Errorf(st.Pos, "missing return value, expected %s", retType)
			}
		}
	case *DeferStmt:
		r.resolveExpr(st.Call, scope)
	case *AssignStmt:
		lType := r.resolveExpr(st.LHS, scope) // CHANGED: capture LHS type
		r.resolveExpr(st.RHS, scope)
		// NEW: Propagate LHS type down to enum shorthands on RHS
		if dot, ok := st.RHS.(*DotEnumExpr); ok {
			if _, isUnknown := dot.GetVType().(*VUnknown); isUnknown {
				dot.SetVType(lType)
			}
		}
	case *ExprStmt:
		r.resolveExpr(st.Expr, scope)
	case *BlockStmt:
		r.resolveBlock(st, NewScope(scope), retType)
	}
}

func (r *Resolver) resolveLocalDecl(d *VarDecl, scope *Scope) {
	// Value may be nil for fixed-array declarations without initializer: var buf: [T; N]
	var valType VType
	if d.Value != nil {
		valType = r.resolveExpr(d.Value, scope)
	} else {
		valType = &VVoid{}
	}

	if d.TypeHint != nil {
		hinted := r.resolveTypeExpr(d.TypeHint, scope)
		// With the updated grammar, [T; N] and [T] are unambiguous:
		//   FixedArrayTypeExpr → always VFixedArray (stack)
		//   ArrayTypeExpr      → always VDynArray   (heap)
		// The old VDynArray → VFixedArray coercion is no longer correct; remove it.

		if d.Value != nil {
			if dot, ok := d.Value.(*DotEnumExpr); ok {
				dot.SetVType(hinted)
			}
		}

		valType = hinted
	} else if d.Value != nil {
		// No explicit annotation — infer from binding kind:
		//   let + array literal → fixed stack array  (keep VFixedArray)
		//   var + array literal → dynamic heap array (convert to VDynArray)
		if !d.IsLet {
			if fa, isFixed := valType.(*VFixedArray); isFixed {
				if _, isArr := d.Value.(*ArrayLitExpr); isArr {
					valType = &VDynArray{Elem: fa.Elem}
				}
			}
		}
	}

	if s, ok := valType.(*VString); ok {
		s.Mutable = !d.IsLet
	}

	for _, name := range d.Binding.Names {
		scope.Define(&Symbol{
			Name:    name,
			Kind:    SymVar,
			Type:    valType,
			Decl:    d,
			IsConst: d.IsLet,
		})
	}
}

func (r *Resolver) resolveIfStmt(st *IfStmt, scope *Scope, retType VType) {
	inner := NewScope(scope)
	switch cond := st.Cond.(type) {
	case *IfLetCond:
		bound := r.resolveExpr(cond.Expr, scope)
		if opt, ok := bound.(*VOptional); ok {
			bound = opt.Elem
		}
		inner.Define(&Symbol{Name: cond.Name, Kind: SymVar, Type: bound})
	case *IfExprCond:
		r.resolveExpr(cond.Expr, scope)
	}
	r.resolveBlock(st.Then, inner, retType)
	if st.Else != nil {
		r.resolveStmt(st.Else, NewScope(scope), retType)
	}
}

// ─── Expression resolution ────────────────────────────────────────────────────

func (r *Resolver) resolveExpr(expr Expr, scope *Scope) VType {
	if expr == nil {
		return &VVoid{}
	}
	var t VType
	switch e := expr.(type) {
	case *IntLitExpr:
		if e.IsUnsigned {
			t = &VInt{Bits: 32, Signed: false}
		} else {
			t = &VInt{Bits: 32, Signed: true}
		}
	case *FloatLitExpr:
		if e.Is32Bit {
			t = &VFloat{Bits: 32}
		} else {
			t = &VFloat{Bits: 64}
		}
	case *BoolLitExpr:
		t = &VBool{}
	case *CharLitExpr:
		t = &VChar{}
	case *StringLitExpr:
		t = &VString{Mutable: false}
	case *NilLitExpr:
		t = &VNil{}

	case *IdentExpr:
		t = r.resolveIdent(e, scope)

	case *DotEnumExpr:
		t = &VUnknown{Name: "." + e.Case}

	case *UnaryExpr:
		inner := r.resolveExpr(e.Operand, scope)
		switch e.Op {
		case UnAddrOf:
			t = &VPointer{Elem: inner}
		case UnNot:
			t = &VBool{}
		default:
			t = inner
		}

	case *BinaryExpr:
		l := r.resolveExpr(e.Left, scope)
		rt := r.resolveExpr(e.Right, scope)

		// Context inference for DotEnumExpr in comparisons.
		if dot, ok := e.Left.(*DotEnumExpr); ok {
			if _, isUnknown := dot.GetVType().(*VUnknown); isUnknown {
				dot.SetVType(rt)
			}
		}
		if dot, ok := e.Right.(*DotEnumExpr); ok {
			if _, isUnknown := dot.GetVType().(*VUnknown); isUnknown {
				dot.SetVType(l)
			}
		}

		switch e.Op {
		case BinEq, BinNeq, BinLt, BinLte, BinGt, BinGte,
			BinAnd, BinOr, BinIdentityEq, BinIdentityNeq:
			t = &VBool{}
		case BinRangeHalfOpen, BinRangeClosed:
			t = &VRange{Elem: l}
		case BinNilCoalesce:
			if opt, ok := l.(*VOptional); ok {
				t = opt.Elem
			} else {
				t = l
			}
		default:
			t = l
		}

	case *TernaryExpr:
		r.resolveExpr(e.Cond, scope)
		t = r.resolveExpr(e.Then, scope)
		r.resolveExpr(e.Else, scope)

	case *CallExpr:
		t = r.resolveCallExpr(e, scope)

	case *MethodCallExpr:
		t = r.resolveMethodCall(e, scope)

	case *FieldExpr:
		t = r.resolveFieldExpr(e, scope)

	case *IndexExpr:
		t = r.resolveIndexExpr(e, scope)

	case *StructLitExpr:
		t = r.resolveStructLit(e, scope)

	case *ArrayLitExpr:
		var elemType VType = &VInt{Bits: 32, Signed: true}
		if len(e.Elems) > 0 {
			elemType = r.resolveExpr(e.Elems[0], scope)
			for _, el := range e.Elems[1:] {
				r.resolveExpr(el, scope)
			}
		}
		t = &VFixedArray{Elem: elemType, Size: len(e.Elems)}

	case *ArrayCtorExpr:
		t = r.resolveArrayCtor(e, scope)

	case *MapLitExpr:
		var keyType VType = &VUnknown{Name: "key"}
		var valType VType = &VUnknown{Name: "value"}
		if len(e.Fields) > 0 {
			keyType = r.resolveExpr(e.Fields[0].Key, scope)
			valType = r.resolveExpr(e.Fields[0].Value, scope)
			for _, f := range e.Fields[1:] {
				r.resolveExpr(f.Key, scope)
				r.resolveExpr(f.Value, scope)
			}
		}
		t = &VMap{Key: keyType, Value: valType}

	case *TupleLitExpr:
		var elems []VType
		for _, el := range e.Elems {
			elems = append(elems, r.resolveExpr(el, scope))
		}
		t = &VTuple{Elems: elems}

	case *ResultExpr:
		inner := r.resolveExpr(e.Value, scope)
		if e.IsOk {
			t = &VResult{Ok: inner, Err: &VVoid{}}
		} else {
			t = &VResult{Ok: &VVoid{}, Err: inner}
		}

	case *TypeConvExpr:
		r.resolveExpr(e.Value, scope)
		t = r.resolveTypeExpr(e.TargetType, scope)

	case *CastExpr:
		r.resolveExpr(e.Value, scope)
		t = r.resolveTypeExpr(e.TargetType, scope)

	case *AnonFuncExpr:
		t = &VFunc{}
	}

	if t == nil {
		t = &VUnknown{Name: "?"}
	}
	expr.SetVType(t)
	return t
}

func (r *Resolver) resolveIdent(e *IdentExpr, scope *Scope) VType {
	sym, ok := scope.Lookup(e.Name)
	if !ok {
		if bt, ok2 := BuiltinTypes[e.Name]; ok2 {
			return bt
		}
		r.diags.Errorf(e.Pos, "undefined identifier %q", e.Name)
		return &VUnknown{Name: e.Name}
	}
	if sym.Type != nil {
		return sym.Type
	}
	return &VUnknown{Name: e.Name}
}

func (r *Resolver) resolveCallExpr(e *CallExpr, scope *Scope) VType {
	if id, ok := e.Func.(*IdentExpr); ok {
		
		// NEW: Intercept Generic Function Calls!
		if len(e.TypeArgs) > 0 {
			if genFn, exists := r.genericFuncs[id.Name]; exists {
				var resolvedArgs []VType
				for _, ta := range e.TypeArgs {
					resolvedArgs = append(resolvedArgs, r.resolveTypeExpr(ta, scope))
				}
				// Clone and append to AST
				concreteName := r.instantiateFunc(genFn, resolvedArgs, e.TypeArgs)
				
				// Transform AST node into a normal, concrete CallExpr
				id.Name = concreteName
				e.TypeArgs = nil
			} else {
				r.diags.Errorf(e.Pos, "unknown generic function %q", id.Name)
			}
		}
		
		if bt, isBT := BuiltinTypes[id.Name]; isBT {
			conv := &TypeConvExpr{
				exprBase:   exprBase{Pos: id.Pos},
				TargetType: &NamedTypeExpr{Pos: id.Pos, Name: id.Name},
				Value:      e.Args[0].Value,
			}
			if len(e.Args) > 0 {
				r.resolveExpr(e.Args[0].Value, scope)
			}
			conv.SetVType(bt)
			e.SetVType(bt)
			return bt
		}
		if sym, ok2 := scope.Lookup(id.Name); ok2 {
			if sym.Kind == SymClass {
				for _, a := range e.Args {
					r.resolveExpr(a.Value, scope)
				}
				return sym.Type
			}
			if sym.Kind == SymEnum {
				if len(e.Args) == 1 && e.Args[0].Label == "rawValue" {
					r.resolveExpr(e.Args[0].Value, scope)
					return &VOptional{Elem: sym.Type}
				}
			}
		}
	}
	r.resolveExpr(e.Func, scope)
	for _, a := range e.Args {
		r.resolveExpr(a.Value, scope)
	}
	if id, ok := e.Func.(*IdentExpr); ok {
		if sym, ok2 := scope.Lookup(id.Name); ok2 {
			if fn, ok3 := sym.Decl.(*FuncDecl); ok3 {
				posIdx := 0
				for _, a := range e.Args {
					if dot, ok4 := a.Value.(*DotEnumExpr); ok4 {
						if _, isUnknown := dot.GetVType().(*VUnknown); isUnknown {
							var matched *Param
							if a.Label != "" {
								for _, p := range fn.Params {
									if p.Name == a.Label {
										matched = p
										break
									}
								}
							} else if posIdx < len(fn.Params) {
								matched = fn.Params[posIdx]
							}
							if matched != nil {
								dot.SetVType(r.resolveTypeExpr(matched.Type, scope))
							}
						}
					}
					if a.Label == "" {
						posIdx++
					}
				}
				if fn.RetType != nil {
					return r.resolveTypeExpr(fn.RetType, scope)
				}
			}
		}
	}
	return &VUnknown{Name: "call"}
}

func (r *Resolver) resolveMethodCall(e *MethodCallExpr, scope *Scope) VType {
	recvType := r.resolveExpr(e.Recv, scope)
	for _, a := range e.Args {
		r.resolveExpr(a.Value, scope)
	}
	switch e.Method {
	case "push", "unshift", "fill", "sort", "reverse", "forEach",
		"delete", "close", "send":
		return &VVoid{}
	case "pop", "shift":
		if da, ok := recvType.(*VDynArray); ok {
			return &VOptional{Elem: da.Elem}
		}
	case "indexOf":
		return &VInt{Bits: 32, Signed: true}
	case "includes":
		return &VBool{}
	case "find":
		if da, ok := recvType.(*VDynArray); ok {
			return &VOptional{Elem: da.Elem}
		}
	case "map":
		return &VDynArray{Elem: &VUnknown{Name: "mapped"}}
	case "filter":
		return recvType
	case "slice":
		return recvType
	case "concat":
		return recvType
	case "receive", "tryReceive":
		if ch, ok := recvType.(*VChan); ok {
			return &VOptional{Elem: ch.Elem}
		}
	case "trySend":
		return &VBool{}
	case "new":
		return recvType
	}

	if r.file != nil {
		var typeName string
		switch rt := recvType.(type) {
		case *VStruct:
			typeName = rt.Name
		case *VClass:
			typeName = rt.Name
		case *VEnum:
			typeName = rt.Name
		}
		if typeName != "" {
			// Pass 1: user-defined receiver methods (FuncDecl with a receiver).
			for _, decl := range r.file.Decls {
				fn, ok := decl.(*FuncDecl)
				if !ok || fn.Receiver == nil || fn.Name != e.Method {
					continue
				}
				if extractTypeName(fn.Receiver.Type) == typeName {
					if fn.RetType != nil {
						return r.resolveTypeExpr(fn.RetType, scope)
					}
					return &VVoid{}
				}
			}

			// Pass 2: native class member method signatures (ClassDecl.Members).
			// This handles calls like libc.socket(...) where "socket" is declared
			// as a ClassMember, not as a FuncDecl with a receiver.
			for _, decl := range r.file.Decls {
				cl, ok := decl.(*ClassDecl)
				if !ok || cl.Name != typeName {
					continue
				}
				for _, m := range cl.Members {
					if m.IsField || m.Name != e.Method {
						continue
					}
					if m.RetType != nil {
						return r.resolveTypeExpr(m.RetType, scope)
					}
					return &VVoid{}
				}
			}
		}
	}

	return &VUnknown{Name: e.Method}
}

func (r *Resolver) resolveFieldExpr(e *FieldExpr, scope *Scope) VType {
	recvType := r.resolveExpr(e.Recv, scope)
	switch rt := recvType.(type) {
	case *VDynArray:
		if e.Field == "length" {
			return &VInt{Bits: 32, Signed: false}
		}
	case *VStruct:
		if rt.Decl != nil {
			for _, f := range rt.Decl.Fields {
				if f.Name == e.Field {
					return r.resolveTypeExpr(f.Type, scope)
				}
			}
		}
	case *VClass:
		if rt.Decl != nil {
			for _, m := range rt.Decl.Members {
				if m.IsField && m.Name == e.Field {
					return r.resolveTypeExpr(m.Type, scope)
				}
			}
		}
	case *VString:
		if e.Field == "length" {
			return &VInt{Bits: 64, Signed: false}
		}
	case *VEnum:
		// NEW: Resolve the .rawValue property wrapper
		if e.Field == "rawValue" {
			return rt.RawType
		}
		if rt.Decl != nil {
			for _, c := range rt.Decl.Cases {
				if c.Name == e.Field {
					return rt
				}
			}
		}
	}
	return &VUnknown{Name: e.Field}
}

func (r *Resolver) resolveIndexExpr(e *IndexExpr, scope *Scope) VType {
	recvType := r.resolveExpr(e.Recv, scope)
	r.resolveExpr(e.Index, scope)
	switch rt := recvType.(type) {
	case *VDynArray:
		return rt.Elem
	case *VFixedArray:
		return rt.Elem
	}
	return &VUnknown{Name: "index"}
}

func (r *Resolver) resolveStructLit(e *StructLitExpr, scope *Scope) VType {
	// NEW: Intercept Generic Struct Literals
	if len(e.TypeArgs) > 0 {
		if genSt, exists := r.genericStructs[e.TypeName]; exists {
			var resolvedArgs []VType
			for _, ta := range e.TypeArgs {
				resolvedArgs = append(resolvedArgs, r.resolveTypeExpr(ta, scope))
			}
			concreteName := r.instantiateStruct(genSt, resolvedArgs, e.TypeArgs)
			e.TypeName = concreteName
			e.TypeArgs = nil
		} else {
			r.diags.Errorf(e.Pos, "unknown generic struct %q", e.TypeName)
		}
	}

	sym, ok := scope.Lookup(e.TypeName)
	for _, f := range e.Fields {
		r.resolveExpr(f.Value, scope)
	}
	if !ok {
		r.diags.Errorf(e.Pos, "unknown type %q in struct literal", e.TypeName)
		return &VUnknown{Name: e.TypeName}
	}
	return sym.Type
}

func (r *Resolver) resolveArrayCtor(e *ArrayCtorExpr, scope *Scope) VType {
	elemType, ok := BuiltinTypes[e.ElemTypeName]
	if !ok {
		if sym, ok2 := scope.Lookup(e.ElemTypeName); ok2 {
			elemType = sym.Type
		} else {
			r.diags.Errorf(e.Pos, "unknown element type %q in array constructor", e.ElemTypeName)
			elemType = &VUnknown{Name: e.ElemTypeName}
		}
	}
	for _, a := range e.Args {
		r.resolveExpr(a.Value, scope)
	}

	if len(e.Args) == 1 && e.Args[0].Label == "" {
		if il, ok := e.Args[0].Value.(*IntLitExpr); ok {
			return &VFixedArray{Elem: elemType, Size: int(il.Value)}
		}
		return &VFixedArray{Elem: elemType, Size: -1}
	}
	hasRepeating := false
	hasCount := false
	countSize := -1
	for _, a := range e.Args {
		switch a.Label {
		case "repeating":
			hasRepeating = true
		case "count":
			hasCount = true
			if il, ok := a.Value.(*IntLitExpr); ok {
				countSize = int(il.Value)
			}
		}
	}
	if hasRepeating && hasCount {
		return &VFixedArray{Elem: elemType, Size: countSize}
	}
	return &VDynArray{Elem: elemType}
}

// ─── Type expression resolution ───────────────────────────────────────────────

func (r *Resolver) resolveTypeExpr(te TypeExpr, scope *Scope) VType {
	if te == nil {
		return &VVoid{}
	}
	switch t := te.(type) {
	case *NamedTypeExpr:
		return r.resolveNamedType(t, scope)
	case *PointerTypeExpr:
		elem := r.resolveTypeExpr(t.Elem, scope)
		pt := &VPointer{Elem: elem, IsConst: t.IsConst}
		if t.Optional {
			return &VOptional{Elem: pt}
		}
		return pt
	case *FixedArrayTypeExpr:
		elem := r.resolveTypeExpr(t.Elem, scope)
		size := -1
		if t.Size != nil {
			r.resolveExpr(t.Size, scope) // walk for completeness
			if il, ok := t.Size.(*IntLitExpr); ok {
				size = int(il.Value)
			}
		}
		return &VFixedArray{Elem: elem, Size: size}
	case *ArrayTypeExpr:
		elem := r.resolveTypeExpr(t.Elem, scope)
		return &VDynArray{Elem: elem}
	case *OptionalTypeExpr:
		elem := r.resolveTypeExpr(t.Elem, scope)
		return &VOptional{Elem: elem}
	case *FuncTypeExpr:
		vf := &VFunc{}
		for _, p := range t.Params {
			vf.Params = append(vf.Params, r.resolveTypeExpr(p, scope))
		}
		vf.Return = r.resolveTypeExpr(t.RetType, scope)
		return vf
	case *TupleTypeExpr:
		vt := &VTuple{Labels: t.Labels}
		for _, e := range t.Elems {
			vt.Elems = append(vt.Elems, r.resolveTypeExpr(e, scope))
		}
		return vt
	case *ResultTypeExpr:
		return &VResult{
			Ok:  r.resolveTypeExpr(t.Ok, scope),
			Err: r.resolveTypeExpr(t.Err, scope),
		}
	case *ChanTypeExpr:
		return &VChan{Elem: r.resolveTypeExpr(t.Elem, scope)}
	case *ExpectedTypeExpr:
		var retVType VType
		if t.ReturnType != nil {
			retVType = r.resolveTypeExpr(t.ReturnType, scope)
		}
		return &VExpected{
			Channel:    t.Channel,
			Value:      t.Value,
			ReturnType: retVType,
		}
	}
	return &VVoid{}
}

func (r *Resolver) resolveNamedType(t *NamedTypeExpr, scope *Scope) VType {
	name := t.Name
	if t.Pkg != "" {
		name = t.Pkg + "." + t.Name
	}
	
	// NEW: Intercept Generic Struct instances!
	if len(t.TypeArgs) > 0 {
		if genSt, exists := r.genericStructs[name]; exists {
			var resolvedArgs []VType
			for _, ta := range t.TypeArgs {
				resolvedArgs = append(resolvedArgs, r.resolveTypeExpr(ta, scope))
			}
			// Clone and append to AST
			concreteName := r.instantiateStruct(genSt, resolvedArgs, t.TypeArgs)
			
			// Rewrite this node to point to the concrete struct
			t.Name = concreteName
			t.TypeArgs = nil
			name = concreteName
		} else {
			r.diags.Errorf(t.Pos, "unknown generic struct %q", name)
		}
	}

	if bt, ok := BuiltinTypes[name]; ok {
		return bt
	}
	if sym, ok := scope.Lookup(name); ok {
		if sym.Type != nil {
			return sym.Type
		}
	}
	r.diags.Errorf(t.Pos, "unknown type %q", name)
	return &VUnknown{Name: name}
}

// ─── Utility ──────────────────────────────────────────────────────────────────

func typeExprFor(vt VType) TypeExpr {
	return &NamedTypeExpr{Name: vtypeName(vt)}
}

func vtypeName(vt VType) string {
	if vt == nil {
		return "void"
	}
	return vt.String()
}

func intVal(e Expr) (int64, bool) {
	if il, ok := e.(*IntLitExpr); ok {
		return il.Value, true
	}
	return 0, false
}

var _ = fmt.Sprintf
var _ = strconv.Itoa
var _ = strings.Join