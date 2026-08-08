package ast

// Visitor is the primitive traversal interface. Visit is called for each node;
// returning a non-nil Visitor descends into that node's children with it, and
// returning nil prunes.
type Visitor interface{ Visit(Node) Visitor }

// Walk traverses n in source order, depth first.
//
// It panics on a node type it does not know (§5.4). That is a tripwire, not a
// defect: the traversal switch cannot fall behind new node types without a
// test failing. It also fires in *your* code if you hand it a node from a
// newer ast than the one you compiled against, which is the intended outcome
// (§8.4).
func Walk(v Visitor, n Node) {
	if isNil(n) {
		return
	}
	if v = v.Visit(n); v == nil {
		return
	}

	switch n := n.(type) {
	// --- leaves ------------------------------------------------------------
	case *Ident, *PrivateIdent, *BasicLit, *TemplateElem, *Comment,
		*Elision, *SuperExpr, *ThisExpr, *ThisType, *PredefinedType,
		*EmptyStmt, *DebuggerStmt, *BadExpr, *BadStmt, *BadDecl, *BadType:
		// no children

	case *CommentGroup:
		for _, c := range n.List {
			Walk(v, c)
		}
	case *QualifiedName:
		Walk(v, n.X)
		Walk(v, n.Sel)

	// --- file --------------------------------------------------------------
	case *File:
		walkStmts(v, n.Items)
		for _, g := range n.Comments {
			Walk(v, g)
		}

	// --- expressions -------------------------------------------------------
	case *TemplateLit:
		walkInterleaved(v, n.Quasis, n.Exprs)
	case *ArrayLit:
		walkExprs(v, n.Elts)
	case *ObjectLit:
		walkExprs(v, n.Props)
	case *PropertyDef:
		Walk(v, n.Key)
		Walk(v, n.Value)
	case *ComputedKey:
		Walk(v, n.X)
	case *SpreadElem:
		Walk(v, n.X)
	case *ParenExpr:
		Walk(v, n.X)
	case *FuncExpr:
		Walk(v, n.Fn)
	case *ClassExpr:
		Walk(v, n.Class)
	case *ArrowFunc:
		Walk(v, n.Ident)
		Walk(v, n.TypeParams)
		Walk(v, n.Params)
		Walk(v, n.Result)
		Walk(v, n.Body)
	case *MemberExpr:
		Walk(v, n.X)
		Walk(v, n.Sel)
	case *IndexExpr:
		Walk(v, n.X)
		Walk(v, n.Index)
	case *CallExpr:
		Walk(v, n.Fun)
		Walk(v, n.TypeArgs)
		walkExprs(v, n.Args)
	case *NewExpr:
		Walk(v, n.Callee)
		Walk(v, n.TypeArgs)
		walkExprs(v, n.Args)
	case *TaggedTemplateExpr:
		Walk(v, n.Tag)
		Walk(v, n.TypeArgs)
		Walk(v, n.Template)
	case *MetaProp:
		Walk(v, n.Prop)
	case *ImportCall:
		walkExprs(v, n.Args)
	case *InstantiationExpr:
		Walk(v, n.X)
		Walk(v, n.TypeArgs)
	case *NonNullExpr:
		Walk(v, n.X)
	case *TypeAssertExpr:
		Walk(v, n.Type)
		Walk(v, n.X)
	case *AsExpr:
		Walk(v, n.X)
		Walk(v, n.Type)
	case *UnaryExpr:
		Walk(v, n.X)
	case *UpdateExpr:
		Walk(v, n.X)
	case *AwaitExpr:
		Walk(v, n.X)
	case *YieldExpr:
		Walk(v, n.X)
	case *BinaryExpr:
		Walk(v, n.X)
		Walk(v, n.Y)
	case *AssignExpr:
		Walk(v, n.Lhs)
		Walk(v, n.Rhs)
	case *CondExpr:
		Walk(v, n.Cond)
		Walk(v, n.Then)
		Walk(v, n.Else)
	case *SeqExpr:
		walkExprs(v, n.Exprs)

	case *ObjectPattern:
		walkExprs(v, n.Props)
	case *ArrayPattern:
		walkExprs(v, n.Elts)
	case *PropertyPattern:
		Walk(v, n.Key)
		Walk(v, n.Value)
	case *AssignPattern:
		Walk(v, n.Lhs)
		Walk(v, n.Rhs)
	case *RestElem:
		Walk(v, n.X)
		Walk(v, n.Type)

	// --- statements --------------------------------------------------------
	case *BlockStmt:
		walkStmts(v, n.List)
	case *ExprStmt:
		Walk(v, n.X)
	case *IfStmt:
		Walk(v, n.Cond)
		Walk(v, n.Body)
		Walk(v, n.Else)
	case *DoWhileStmt:
		Walk(v, n.Body)
		Walk(v, n.Cond)
	case *WhileStmt:
		Walk(v, n.Cond)
		Walk(v, n.Body)
	case *ForStmt:
		Walk(v, n.Init)
		Walk(v, n.Cond)
		Walk(v, n.Post)
		Walk(v, n.Body)
	case *ForInStmt:
		Walk(v, n.Left)
		Walk(v, n.Right)
		Walk(v, n.Body)
	case *ForOfStmt:
		Walk(v, n.Left)
		Walk(v, n.Right)
		Walk(v, n.Body)
	case *BranchStmt:
		Walk(v, n.Label)
	case *ReturnStmt:
		Walk(v, n.Result)
	case *LabeledStmt:
		Walk(v, n.Label)
		Walk(v, n.Stmt)
	case *ThrowStmt:
		Walk(v, n.X)
	case *TryStmt:
		Walk(v, n.Body)
		Walk(v, n.Catch)
		Walk(v, n.Finally)
	case *CatchClause:
		Walk(v, n.Param)
		Walk(v, n.CatchType)
		Walk(v, n.Body)
	case *FinallyClause:
		Walk(v, n.Body)
	case *SwitchStmt:
		Walk(v, n.Tag)
		for _, c := range n.Cases {
			Walk(v, c)
		}
	case *CaseClause:
		Walk(v, n.Cond)
		walkStmts(v, n.Body)

	// --- declarations ------------------------------------------------------
	case *VarDecl:
		for _, b := range n.List {
			Walk(v, b)
		}
	case *Binding:
		Walk(v, n.Name)
		Walk(v, n.Pattern)
		Walk(v, n.Type)
		Walk(v, n.Init)
	case *FuncDecl:
		walkDecorators(v, n.Decorators)
		Walk(v, n.Name)
		Walk(v, n.TypeParams)
		Walk(v, n.Params)
		Walk(v, n.Result)
		Walk(v, n.Body)
	case *ClassDecl:
		walkDecorators(v, n.Decorators)
		Walk(v, n.Name)
		Walk(v, n.TypeParams)
		Walk(v, n.Extends)
		Walk(v, n.Implements)
		walkDecls(v, n.Members)
	case *StructDecl:
		walkDecorators(v, n.Decorators)
		Walk(v, n.Name)
		Walk(v, n.TypeParams)
		Walk(v, n.Extends) // nil in a valid program; see decl.go
		Walk(v, n.Implements)
		Walk(v, n.Body)
	case *StructBody:
		walkDecls(v, n.Members)
	case *HeritageClause:
		walkTypes(v, n.Types)
	case *FieldDecl:
		walkDecorators(v, n.Decorators)
		Walk(v, n.Name)
		Walk(v, n.Type)
		Walk(v, n.Init)
	case *MethodDecl:
		walkDecorators(v, n.Decorators)
		Walk(v, n.Name)
		Walk(v, n.TypeParams)
		Walk(v, n.Params)
		Walk(v, n.Result)
		Walk(v, n.Body)
	case *CtorDecl:
		walkDecorators(v, n.Decorators)
		Walk(v, n.Params)
		Walk(v, n.Body)
	case *StaticBlockDecl:
		Walk(v, n.Body)
	case *InterfaceDecl:
		Walk(v, n.Name)
		Walk(v, n.TypeParams)
		Walk(v, n.Extends)
		Walk(v, n.Body)
	case *TypeAliasDecl:
		Walk(v, n.Name)
		Walk(v, n.TypeParams)
		Walk(v, n.Type)
	case *EnumDecl:
		Walk(v, n.Name)
		Walk(v, n.Underlying)
		for _, m := range n.Members {
			Walk(v, m)
		}
	case *EnumMember:
		Walk(v, n.Name)
		Walk(v, n.Value)
	case *NamespaceDecl:
		Walk(v, n.Name)
		walkStmts(v, n.Items)
	case *ModuleDecl:
		Walk(v, n.Name)
		walkStmts(v, n.Items)
	case *AmbientDecl:
		Walk(v, n.Inner)
	case *Decorator:
		Walk(v, n.X)
	case *ParamList:
		Walk(v, n.This)
		for _, p := range n.List {
			Walk(v, p)
		}
		Walk(v, n.Rest)
	case *ThisParam:
		Walk(v, n.Type)
	case *Param:
		walkDecorators(v, n.Decorators)
		Walk(v, n.Name)
		Walk(v, n.Type)
		Walk(v, n.Init)

	case *ImportDecl:
		Walk(v, n.Default)
		Walk(v, n.Namespace)
		for _, s := range n.Named {
			Walk(v, s)
		}
		Walk(v, n.Path)
		Walk(v, n.With)
	case *ImportSpec:
		Walk(v, n.Name)
		Walk(v, n.Local)
	case *ImportEqualsDecl:
		Walk(v, n.Name)
		Walk(v, n.Entity)
		Walk(v, n.Path)
	case *ExportDecl:
		Walk(v, n.StarAs)
		for _, s := range n.Named {
			Walk(v, s)
		}
		Walk(v, n.Decl)
		Walk(v, n.Value)
		Walk(v, n.Entity)
		Walk(v, n.Namespace)
		Walk(v, n.Path)
		Walk(v, n.With)
	case *ExportSpec:
		Walk(v, n.Name)
		Walk(v, n.As)
	case *WithClause:
		for _, a := range n.List {
			Walk(v, a)
		}
	case *ImportAttr:
		Walk(v, n.Key)
		Walk(v, n.Value)

	// --- types -------------------------------------------------------------
	case *TypeRef:
		Walk(v, n.Name)
		Walk(v, n.Args)
	case *LiteralType:
		Walk(v, n.Value)
	case *ParenType:
		Walk(v, n.X)
	case *ArrayType:
		Walk(v, n.Elem)
	case *IndexedAccessType:
		Walk(v, n.X)
		Walk(v, n.Index)
	case *UnionType:
		walkTypes(v, n.Types)
	case *IntersectionType:
		walkTypes(v, n.Types)
	case *TypeOp:
		Walk(v, n.X)
	case *InferType:
		Walk(v, n.Name)
		Walk(v, n.Constraint)
	case *CondType:
		Walk(v, n.Check)
		Walk(v, n.Extends)
		Walk(v, n.Then)
		Walk(v, n.Else)
	case *FuncType:
		Walk(v, n.TypeParams)
		Walk(v, n.Params)
		Walk(v, n.Result)
	case *CtorType:
		Walk(v, n.TypeParams)
		Walk(v, n.Params)
		Walk(v, n.Result)
	case *ObjectType:
		walkTypes(v, n.Members)
	case *MappedType:
		Walk(v, n.Name)
		Walk(v, n.Constraint)
		Walk(v, n.As)
		Walk(v, n.Type)
	case *TupleType:
		walkTypes(v, n.Elems)
	case *TupleElem:
		Walk(v, n.Name)
		Walk(v, n.Type)
	case *TypeQuery:
		Walk(v, n.Name)
		Walk(v, n.Args)
	case *ImportType:
		Walk(v, n.Path)
		Walk(v, n.With)
		Walk(v, n.Qualifier)
		Walk(v, n.Args)
	case *TemplateLiteralType:
		for i, q := range n.Quasis {
			Walk(v, q)
			if i < len(n.Types) {
				Walk(v, n.Types[i])
			}
		}
	case *TypePredicate:
		Walk(v, n.Param)
		Walk(v, n.Type)

	case *PropertySig:
		Walk(v, n.Name)
		Walk(v, n.Type)
	case *MethodSig:
		Walk(v, n.Name)
		Walk(v, n.TypeParams)
		Walk(v, n.Params)
		Walk(v, n.Result)
	case *CallSig:
		Walk(v, n.TypeParams)
		Walk(v, n.Params)
		Walk(v, n.Result)
	case *ConstructSig:
		Walk(v, n.TypeParams)
		Walk(v, n.Params)
		Walk(v, n.Result)
	case *IndexSig:
		Walk(v, n.Name)
		Walk(v, n.Key)
		Walk(v, n.Type)
	case *AccessorSig:
		Walk(v, n.Name)
		Walk(v, n.Param)
		Walk(v, n.Result)

	case *ParamTypeList:
		Walk(v, n.This)
		for _, p := range n.List {
			Walk(v, p)
		}
		Walk(v, n.Rest)
	case *ParamType:
		Walk(v, n.Name)
		Walk(v, n.Type)
	case *TypeParamList:
		for _, p := range n.List {
			Walk(v, p)
		}
	case *TypeParam:
		Walk(v, n.Name)
		Walk(v, n.Constraint)
		Walk(v, n.Type)
		Walk(v, n.Default)
	case *TypeArgList:
		walkTypes(v, n.List)

	default:
		panic("ast.Walk: unknown node type")
	}

	v.Visit(nil)
}

func walkExprs(v Visitor, xs []Expr) {
	for _, x := range xs {
		Walk(v, x)
	}
}
func walkStmts(v Visitor, xs []Stmt) {
	for _, x := range xs {
		Walk(v, x)
	}
}
func walkDecls(v Visitor, xs []Decl) {
	for _, x := range xs {
		Walk(v, x)
	}
}
func walkTypes(v Visitor, xs []TypeExpr) {
	for _, x := range xs {
		Walk(v, x)
	}
}
func walkDecorators(v Visitor, xs []*Decorator) {
	for _, x := range xs {
		Walk(v, x)
	}
}

// walkInterleaved visits quasis and substitutions in source order, which for a
// template means quasi, expr, quasi, expr, ..., quasi.
func walkInterleaved(v Visitor, quasis []*TemplateElem, exprs []Expr) {
	for i, q := range quasis {
		Walk(v, q)
		if i < len(exprs) {
			Walk(v, exprs[i])
		}
	}
}

type inspector func(Node) bool

func (f inspector) Visit(n Node) Visitor {
	if n != nil && f(n) {
		return f
	}
	return nil
}

// Inspect walks n, calling f for each node. f returns whether to descend.
//
// It exists because Walk requires declaring a type with a method just to run a
// closure, and nearly every consumer wants the closure (§5.5).
func Inspect(n Node, f func(Node) bool) { Walk(inspector(f), n) }

// Unparen strips *ParenExpr wrappers.
//
// It exists because §5.4 retains ParenExpr, so every consumer that doesn't
// care about parentheses would otherwise write this itself.
func Unparen(x Expr) Expr {
	for {
		p, ok := x.(*ParenExpr)
		if !ok {
			return x
		}
		x = p.X
	}
}

// UnparenType is the same for the type hierarchy, which has its own
// ParenType (§5.1).
func UnparenType(t TypeExpr) TypeExpr {
	for {
		p, ok := t.(*ParenType)
		if !ok {
			return t
		}
		t = p.X
	}
}