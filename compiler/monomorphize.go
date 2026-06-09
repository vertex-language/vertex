package compiler

import "strings"

// instantiateFunc creates a concrete copy of a generic function.
func (r *Resolver) instantiateFunc(gen *FuncDecl, resolvedArgs []VType, astArgs []TypeExpr) string {
	nameExt := ""
	for _, vt := range resolvedArgs {
		nameExt += "_" + sanitizeMangle(vt.String())
	}
	concreteName := gen.Name + nameExt

	if r.instFuncs[concreteName] {
		return concreteName
	}
	r.instFuncs[concreteName] = true

	subst := make(map[string]TypeExpr)
	for i, tp := range gen.TypeParams {
		subst[tp] = astArgs[i]
	}

	cloned := cloneFuncDecl(gen, subst)
	cloned.Name = concreteName
	cloned.TypeParams = nil // Strip generic params!

	// Append the new concrete function to the file so it gets resolved and lowered.
	r.file.Decls = append(r.file.Decls, cloned)
	r.pkg.Define(&Symbol{Name: concreteName, Kind: SymFunc, Decl: cloned})

	return concreteName
}

// instantiateStruct creates a concrete copy of a generic struct.
func (r *Resolver) instantiateStruct(gen *StructDecl, resolvedArgs []VType, astArgs []TypeExpr) string {
	nameExt := ""
	for _, vt := range resolvedArgs {
		nameExt += "_" + sanitizeMangle(vt.String())
	}
	concreteName := gen.Name + nameExt

	if r.instStructs[concreteName] {
		return concreteName
	}
	r.instStructs[concreteName] = true

	subst := make(map[string]TypeExpr)
	for i, tp := range gen.TypeParams {
		subst[tp] = astArgs[i]
	}

	cloned := cloneStructDecl(gen, subst)
	cloned.Name = concreteName
	cloned.TypeParams = nil

	r.file.Decls = append(r.file.Decls, cloned)
	r.pkg.Define(&Symbol{
		Name: concreteName, Kind: SymStruct,
		Type: &VStruct{Name: concreteName, Decl: cloned},
		Decl: cloned,
	})

	return concreteName
}

func sanitizeMangle(s string) string {
	s = strings.ReplaceAll(s, "*", "ptr_")
	s = strings.ReplaceAll(s, "[]", "arr_")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

// ─── AST Deep Cloning functions ──────────────────────────────────────────────

func cloneFuncDecl(d *FuncDecl, subst map[string]TypeExpr) *FuncDecl {
	c := &FuncDecl{
		Pos: d.Pos, Name: d.Name, Qualifier: d.Qualifier,
		RetType: cloneTypeExpr(d.RetType, subst),
		Body:    cloneBlock(d.Body, subst),
	}
	if d.Receiver != nil {
		c.Receiver = &Receiver{Pos: d.Receiver.Pos, Name: d.Receiver.Name, IsPtr: d.Receiver.IsPtr, Type: cloneTypeExpr(d.Receiver.Type, subst)}
	}
	for _, p := range d.Params {
		c.Params = append(c.Params, &Param{Pos: p.Pos, Name: p.Name, IsVariadic: p.IsVariadic, Type: cloneTypeExpr(p.Type, subst)})
	}
	return c
}

func cloneStructDecl(d *StructDecl, subst map[string]TypeExpr) *StructDecl {
	c := &StructDecl{Pos: d.Pos, Name: d.Name}
	for _, f := range d.Fields {
		c.Fields = append(c.Fields, &FieldDecl{Pos: f.Pos, Name: f.Name, Type: cloneTypeExpr(f.Type, subst)})
	}
	return c
}

func cloneTypeExpr(t TypeExpr, subst map[string]TypeExpr) TypeExpr {
	if t == nil {
		return nil
	}
	switch te := t.(type) {
	case *NamedTypeExpr:
		// Substitute the type argument!
		if rep, ok := subst[te.Name]; ok {
			return cloneTypeExpr(rep, nil) 
		}
		c := &NamedTypeExpr{Pos: te.Pos, Pkg: te.Pkg, Name: te.Name}
		for _, arg := range te.TypeArgs {
			c.TypeArgs = append(c.TypeArgs, cloneTypeExpr(arg, subst))
		}
		return c
	case *PointerTypeExpr:
		return &PointerTypeExpr{Pos: te.Pos, IsConst: te.IsConst, Optional: te.Optional, Elem: cloneTypeExpr(te.Elem, subst)}
	case *ArrayTypeExpr:
		return &ArrayTypeExpr{Pos: te.Pos, Elem: cloneTypeExpr(te.Elem, subst)}
	case *OptionalTypeExpr:
		return &OptionalTypeExpr{Pos: te.Pos, Elem: cloneTypeExpr(te.Elem, subst)}
	case *FuncTypeExpr:
		c := &FuncTypeExpr{Pos: te.Pos, RetType: cloneTypeExpr(te.RetType, subst)}
		for _, p := range te.Params {
			c.Params = append(c.Params, cloneTypeExpr(p, subst))
		}
		return c
	case *TupleTypeExpr:
		c := &TupleTypeExpr{Pos: te.Pos, Labels: append([]string{}, te.Labels...)}
		for _, el := range te.Elems {
			c.Elems = append(c.Elems, cloneTypeExpr(el, subst))
		}
		return c
	case *ResultTypeExpr:
		return &ResultTypeExpr{Pos: te.Pos, Ok: cloneTypeExpr(te.Ok, subst), Err: cloneTypeExpr(te.Err, subst)}
	case *ChanTypeExpr:
		return &ChanTypeExpr{Pos: te.Pos, Elem: cloneTypeExpr(te.Elem, subst)}
	case *ExpectedTypeExpr:
		return &ExpectedTypeExpr{Pos: te.Pos, ReturnType: cloneTypeExpr(te.ReturnType, subst), Channel: te.Channel, Value: te.Value}
	}
	return nil
}

func cloneExpr(expr Expr, subst map[string]TypeExpr) Expr {
	if expr == nil { return nil }
	switch e := expr.(type) {
	case *IntLitExpr: return &IntLitExpr{exprBase: exprBase{Pos: e.Pos}, Value: e.Value, IsUnsigned: e.IsUnsigned}
	case *FloatLitExpr: return &FloatLitExpr{exprBase: exprBase{Pos: e.Pos}, Value: e.Value, Is32Bit: e.Is32Bit}
	case *BoolLitExpr: return &BoolLitExpr{exprBase: exprBase{Pos: e.Pos}, Value: e.Value}
	case *CharLitExpr: return &CharLitExpr{exprBase: exprBase{Pos: e.Pos}, Value: e.Value}
	case *StringLitExpr: return &StringLitExpr{exprBase: exprBase{Pos: e.Pos}, Value: e.Value}
	case *NilLitExpr: return &NilLitExpr{exprBase: exprBase{Pos: e.Pos}}
	case *IdentExpr: return &IdentExpr{exprBase: exprBase{Pos: e.Pos}, Name: e.Name}
	case *DotEnumExpr: return &DotEnumExpr{exprBase: exprBase{Pos: e.Pos}, Case: e.Case}
	case *BinaryExpr: return &BinaryExpr{exprBase: exprBase{Pos: e.Pos}, Op: e.Op, Left: cloneExpr(e.Left, subst), Right: cloneExpr(e.Right, subst)}
	case *UnaryExpr: return &UnaryExpr{exprBase: exprBase{Pos: e.Pos}, Op: e.Op, Operand: cloneExpr(e.Operand, subst)}
	case *TernaryExpr: return &TernaryExpr{exprBase: exprBase{Pos: e.Pos}, Cond: cloneExpr(e.Cond, subst), Then: cloneExpr(e.Then, subst), Else: cloneExpr(e.Else, subst)}
	case *CallExpr:
		c := &CallExpr{exprBase: exprBase{Pos: e.Pos}, Func: cloneExpr(e.Func, subst)}
		for _, a := range e.TypeArgs { c.TypeArgs = append(c.TypeArgs, cloneTypeExpr(a, subst)) }
		for _, a := range e.Args { c.Args = append(c.Args, &Arg{Pos: a.Pos, Label: a.Label, Value: cloneExpr(a.Value, subst)}) }
		return c
	case *MethodCallExpr:
		c := &MethodCallExpr{exprBase: exprBase{Pos: e.Pos}, Recv: cloneExpr(e.Recv, subst), Method: e.Method}
		for _, a := range e.Args { c.Args = append(c.Args, &Arg{Pos: a.Pos, Label: a.Label, Value: cloneExpr(a.Value, subst)}) }
		return c
	case *FieldExpr: return &FieldExpr{exprBase: exprBase{Pos: e.Pos}, Recv: cloneExpr(e.Recv, subst), Field: e.Field}
	case *IndexExpr: return &IndexExpr{exprBase: exprBase{Pos: e.Pos}, Recv: cloneExpr(e.Recv, subst), Index: cloneExpr(e.Index, subst)}
	case *StructLitExpr:
		c := &StructLitExpr{exprBase: exprBase{Pos: e.Pos}, TypeName: e.TypeName}
		for _, a := range e.TypeArgs { c.TypeArgs = append(c.TypeArgs, cloneTypeExpr(a, subst)) }
		for _, f := range e.Fields { c.Fields = append(c.Fields, &StructLitField{Pos: f.Pos, Name: f.Name, Value: cloneExpr(f.Value, subst)}) }
		return c
	case *ArrayLitExpr:
		c := &ArrayLitExpr{exprBase: exprBase{Pos: e.Pos}}
		for _, el := range e.Elems { c.Elems = append(c.Elems, cloneExpr(el, subst)) }
		return c
	case *ArrayCtorExpr:
		c := &ArrayCtorExpr{exprBase: exprBase{Pos: e.Pos}, ElemTypeName: e.ElemTypeName}
		if rep, ok := subst[c.ElemTypeName]; ok {
			if nte, ok2 := rep.(*NamedTypeExpr); ok2 {
				c.ElemTypeName = nte.Name
			}
		}
		for _, a := range e.Args { c.Args = append(c.Args, &Arg{Pos: a.Pos, Label: a.Label, Value: cloneExpr(a.Value, subst)}) }
		return c
	case *MapLitExpr:
		c := &MapLitExpr{exprBase: exprBase{Pos: e.Pos}}
		for _, f := range e.Fields { c.Fields = append(c.Fields, &MapLitField{Pos: f.Pos, Key: cloneExpr(f.Key, subst), Value: cloneExpr(f.Value, subst)}) }
		return c
	case *TypeConvExpr: return &TypeConvExpr{exprBase: exprBase{Pos: e.Pos}, TargetType: cloneTypeExpr(e.TargetType, subst), Value: cloneExpr(e.Value, subst)}
	case *CastExpr: return &CastExpr{exprBase: exprBase{Pos: e.Pos}, TargetType: cloneTypeExpr(e.TargetType, subst), Value: cloneExpr(e.Value, subst)}
	case *ResultExpr: return &ResultExpr{exprBase: exprBase{Pos: e.Pos}, IsOk: e.IsOk, Value: cloneExpr(e.Value, subst)}
	case *TupleLitExpr:
		c := &TupleLitExpr{exprBase: exprBase{Pos: e.Pos}}
		for _, el := range e.Elems { c.Elems = append(c.Elems, cloneExpr(el, subst)) }
		return c
	case *AnonFuncExpr:
		c := &AnonFuncExpr{exprBase: exprBase{Pos: e.Pos}, Qualifier: e.Qualifier, RetType: cloneTypeExpr(e.RetType, subst), Body: cloneBlock(e.Body, subst)}
		for _, p := range e.Params { c.Params = append(c.Params, &Param{Pos: p.Pos, Name: p.Name, IsVariadic: p.IsVariadic, Type: cloneTypeExpr(p.Type, subst)}) }
		return c
	}
	return nil
}

func cloneStmt(stmt Stmt, subst map[string]TypeExpr) Stmt {
	if stmt == nil { return nil }
	switch s := stmt.(type) {
	case *BlockStmt: return cloneBlock(s, subst)
	case *LocalDeclStmt: return &LocalDeclStmt{Pos: s.Pos, Decl: cloneVarDecl(s.Decl, subst)}
	case *IfStmt:
		c := &IfStmt{Pos: s.Pos, Then: cloneBlock(s.Then, subst), Else: cloneStmt(s.Else, subst)}
		if letCond, ok := s.Cond.(*IfLetCond); ok {
			c.Cond = &IfLetCond{Pos: letCond.Pos, Name: letCond.Name, Expr: cloneExpr(letCond.Expr, subst)}
		} else if exprCond, ok := s.Cond.(*IfExprCond); ok {
			c.Cond = &IfExprCond{Pos: exprCond.Pos, Expr: cloneExpr(exprCond.Expr, subst)}
		}
		return c
	case *WhileStmt: return &WhileStmt{Pos: s.Pos, Cond: cloneExpr(s.Cond, subst), Body: cloneBlock(s.Body, subst)}
	case *ForInStmt: return &ForInStmt{Pos: s.Pos, Var: s.Var, Iter: cloneExpr(s.Iter, subst), Body: cloneBlock(s.Body, subst)}
	case *SwitchStmt:
		c := &SwitchStmt{Pos: s.Pos, Subj: cloneExpr(s.Subj, subst)}
		for _, sc := range s.Cases {
			cc := &SwitchCase{Pos: sc.Pos, IsDefault: sc.IsDefault}
			for _, p := range sc.Patterns {
				switch pat := p.(type) {
				case *ExprPattern: cc.Patterns = append(cc.Patterns, &ExprPattern{Pos: pat.Pos, Expr: cloneExpr(pat.Expr, subst)})
				case *EnumShortPattern: cc.Patterns = append(cc.Patterns, &EnumShortPattern{Pos: pat.Pos, Case: pat.Case})
				case *ResultOkPattern: cc.Patterns = append(cc.Patterns, &ResultOkPattern{Pos: pat.Pos, Bind: pat.Bind})
				case *ResultErrPattern: cc.Patterns = append(cc.Patterns, &ResultErrPattern{Pos: pat.Pos, Bind: pat.Bind})
				}
			}
			for _, b := range sc.Body { cc.Body = append(cc.Body, cloneStmt(b, subst)) }
			c.Cases = append(c.Cases, cc)
		}
		return c
	case *ReturnStmt: return &ReturnStmt{Pos: s.Pos, Value: cloneExpr(s.Value, subst)}
	case *DeferStmt: return &DeferStmt{Pos: s.Pos, Call: cloneExpr(s.Call, subst)}
	case *BreakStmt: return &BreakStmt{Pos: s.Pos}
	case *ContinueStmt: return &ContinueStmt{Pos: s.Pos}
	case *FallthroughStmt: return &FallthroughStmt{Pos: s.Pos}
	case *AssignStmt: return &AssignStmt{Pos: s.Pos, LHS: cloneExpr(s.LHS, subst), Op: s.Op, RHS: cloneExpr(s.RHS, subst)}
	case *ExprStmt: return &ExprStmt{Pos: s.Pos, Expr: cloneExpr(s.Expr, subst)}
	}
	return nil
}

func cloneBlock(b *BlockStmt, subst map[string]TypeExpr) *BlockStmt {
	if b == nil { return nil }
	c := &BlockStmt{Pos: b.Pos}
	for _, s := range b.Stmts {
		c.Stmts = append(c.Stmts, cloneStmt(s, subst))
	}
	return c
}

func cloneVarDecl(d *VarDecl, subst map[string]TypeExpr) *VarDecl {
	c := &VarDecl{Pos: d.Pos, IsLet: d.IsLet, IsWeak: d.IsWeak, TypeHint: cloneTypeExpr(d.TypeHint, subst), Value: cloneExpr(d.Value, subst)}
	c.Binding = &BindingPattern{Pos: d.Binding.Pos, Names: append([]string{}, d.Binding.Names...)}
	return c
}