package ast

import "fmt"

type Visitor interface {
	// Visit is called for each node. Returning a non-nil visitor walks the
	// node's children with it, then calls Visit(nil) on it to signal the end
	// of that subtree.
	Visit(node Node) Visitor
}

func walkExprs(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkStmts(v Visitor, list []Stmt) {
	for _, s := range list {
		Walk(v, s)
	}
}

func walkIdents(v Visitor, list []*Ident) {
	for _, x := range list {
		Walk(v, x)
	}
}

func maybe(v Visitor, n Node) {
	// A typed-nil in an interface is not == nil, so callers pass concrete
	// pointers and this guards the common optional-field case.
	switch x := n.(type) {
	case nil:
		return
	case Expr:
		if x == nil {
			return
		}
	}
	Walk(v, n)
}

func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *Comment:
		// leaf

	case *CommentGroup:
		for _, c := range n.List {
			Walk(v, c)
		}

	case *Ident, *BasicLit, *NamespaceExpr, *AbstractType,
		*BadExpr, *BadStmt, *BadDecl, *BranchStmt, *BuildClause, *Marker:
		// leaves

	// ---- shared parts

	case *Param:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		Walk(v, n.Type)

	case *ParamList:
		for _, p := range n.List {
			Walk(v, p)
		}

	case *TypeParam:
		Walk(v, n.Name)
		if n.Constraint != nil {
			Walk(v, n.Constraint)
		}

	case *TypeParamList:
		for _, p := range n.List {
			Walk(v, p)
		}

	// ---- expressions

	case *ParenExpr:
		Walk(v, n.X)

	case *TupleExpr:
		walkExprs(v, n.Elems)

	case *ArrayLit:
		walkExprs(v, n.Elems)

	case *CompositeLit:
		Walk(v, n.Type)
		walkExprs(v, n.Elems)

	case *MapLit:
		walkExprs(v, n.Elems)

	case *KeyValueExpr:
		Walk(v, n.Key)
		Walk(v, n.Value)

	case *EnumShorthand:
		Walk(v, n.Name)
		walkExprs(v, n.Args)

	case *FuncLit:
		Walk(v, n.Type)
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *SelectorExpr:
		Walk(v, n.X)
		Walk(v, n.Sel)

	case *TupleIndexExpr:
		Walk(v, n.X)

	case *IndexExpr:
		Walk(v, n.X)
		walkExprs(v, n.Indices)

	case *CallExpr:
		Walk(v, n.Fun)
		walkExprs(v, n.Args)

	case *LaunchExpr:
		if n.Config != nil {
			Walk(v, n.Config)
		}
		Walk(v, n.Call)

	case *LaunchConfig:
		Walk(v, n.Blocks)
		Walk(v, n.Threads)

	case *AwaitExpr:
		Walk(v, n.X)

	case *UnaryExpr:
		Walk(v, n.X)

	case *BinaryExpr:
		Walk(v, n.X)
		Walk(v, n.Y)

	case *CastExpr:
		Walk(v, n.X)
		Walk(v, n.Type)

	case *TransferExpr:
		Walk(v, n.Target)

	// ---- types

	case *OwnershipType:
		Walk(v, n.X)

	case *ArrayType:
		if n.Len != nil {
			Walk(v, n.Len)
		}
		Walk(v, n.Elem)

	case *MapType:
		Walk(v, n.Key)
		Walk(v, n.Value)

	case *FuncType:
		Walk(v, n.Params)
		if n.Marker != nil {
			Walk(v, n.Marker)
		}
		if n.Result != nil {
			Walk(v, n.Result)
		}

	case *ChanType:
		Walk(v, n.Elem)

	case *PointerType:
		Walk(v, n.Elem)

	case *TensorType:
		Walk(v, n.Elem)
		walkExprs(v, n.Shape)

	// ---- statements

	case *BlockStmt:
		walkStmts(v, n.List)

	case *DeclStmt:
		Walk(v, n.Decl)

	case *ExprStmt:
		Walk(v, n.X)

	case *AssignStmt:
		walkExprs(v, n.Targets)
		walkExprs(v, n.Values)

	case *IfStmt:
		Walk(v, n.Cond)
		Walk(v, n.Body)
		if n.Else != nil {
			Walk(v, n.Else)
		}

	case *WhileStmt:
		Walk(v, n.Cond)
		Walk(v, n.Body)

	case *ForStmt:
		walkIdents(v, n.Names)
		Walk(v, n.X)
		Walk(v, n.Body)

	case *SwitchStmt:
		Walk(v, n.Tag)
		for _, c := range n.Cases {
			Walk(v, c)
		}

	case *CaseClause:
		walkExprs(v, n.Patterns)
		walkStmts(v, n.Body)

	case *EnumPattern:
		Walk(v, n.Name)
		walkIdents(v, n.Binds)

	case *ReturnStmt:
		walkExprs(v, n.Results)

	case *DeferStmt:
		Walk(v, n.Call)

	case *SelectStmt:
		for _, c := range n.Cases {
			Walk(v, c)
		}

	case *SelectClause:
		walkExprs(v, n.Targets)
		if n.Op != nil {
			Walk(v, n.Op)
		}
		walkStmts(v, n.Body)

	// ---- declarations

	case *Receiver:
		Walk(v, n.Name)
		Walk(v, n.Type)

	case *FuncDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		if n.Recv != nil {
			Walk(v, n.Recv)
		}
		Walk(v, n.Name)
		if n.TypeParams != nil {
			Walk(v, n.TypeParams)
		}
		Walk(v, n.Type)
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *Field:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		Walk(v, n.Name)
		Walk(v, n.Type)
		if n.Default != nil {
			Walk(v, n.Default)
		}

	case *RecordDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		Walk(v, n.Name)
		if n.TypeParams != nil {
			Walk(v, n.TypeParams)
		}
		for _, f := range n.Fields {
			Walk(v, f)
		}

	case *Variant:
		Walk(v, n.Name)
		walkExprs(v, n.Payload)
		if n.Value != nil {
			Walk(v, n.Value)
		}

	case *EnumDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		Walk(v, n.Name)
		if n.TypeParams != nil {
			Walk(v, n.TypeParams)
		}
		if n.Discrim != nil {
			Walk(v, n.Discrim)
		}
		for _, x := range n.Variants {
			Walk(v, x)
		}

	case *TypeAliasDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		Walk(v, n.Name)
		if n.TypeParams != nil {
			Walk(v, n.TypeParams)
		}
		Walk(v, n.Target)

	case *ConstraintElem:
		if n.Set != nil {
			Walk(v, n.Set)
		}
		if n.Method != nil {
			Walk(v, n.Method)
		}

	case *MethodReq:
		Walk(v, n.Name)
		Walk(v, n.Params)
		if n.Result != nil {
			Walk(v, n.Result)
		}

	case *ConstraintDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		Walk(v, n.Name)
		for _, e := range n.Elems {
			Walk(v, e)
		}

	case *Binding:
		Walk(v, n.Name)
		if n.Type != nil {
			Walk(v, n.Type)
		}

	case *VarDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		for _, b := range n.Bindings {
			Walk(v, b)
		}
		walkExprs(v, n.Values)

	case *ImportDecl:
		for _, p := range n.Paths {
			Walk(v, p)
		}

	case *VariantTag:
		for _, t := range n.Tags {
			Walk(v, t)
		}

	case *DeclareDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		if n.Variant != nil {
			Walk(v, n.Variant)
		}
		Walk(v, n.Path)
		for _, m := range n.Members {
			Walk(v, m)
		}

	case *ForeignFunc:
		walkIdents(v, n.Modifiers)
		if n.Name != nil {
			Walk(v, n.Name)
		}
		Walk(v, n.Params)
		if n.Result != nil {
			Walk(v, n.Result)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ForeignClass:
		walkIdents(v, n.Modifiers)
		Walk(v, n.Name)
		for _, m := range n.Members {
			Walk(v, m)
		}

	case *ForeignField:
		Walk(v, n.Name)
		Walk(v, n.Type)

	// ---- containers

	case *File:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		Walk(v, n.Name)
		if n.Build != nil {
			Walk(v, n.Build)
		}
		for _, d := range n.Imports {
			Walk(v, d)
		}
		for _, d := range n.Decls {
			Walk(v, d)
		}

	case *Package:
		for _, f := range n.Files {
			Walk(v, f)
		}

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	v.Visit(nil)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect walks node in depth-first order, calling f for each node. If f returns
// false, the node's children are skipped. f is also called with nil after a
// subtree's children have been walked.
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}