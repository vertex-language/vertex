package analyzer

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// There are two entry points into this file, not one, and the split is
// load-bearing. A constraint is never a value type and is legal only in a
// bracket position, and types.Constraint deliberately does not implement
// types.Type. So a constraint position calls constraintExpr and a type position
// calls typ, and `var c: Ordered` falls out of typ as an ordinary diagnostic
// rather than needing a predicate someone must remember to call.

// typ converts a type expression to a types.Type, recording the result.
func (c *Checker) typ(e ast.Expr) types.Type {
	t := c.typInternal(e)
	c.info.RecordType(e, types.NewTypeAndValue(types.TypeMode, t, nil))
	return t
}

func (c *Checker) typInternal(e ast.Expr) types.Type {
	switch x := e.(type) {
	case nil, *ast.BadExpr:
		// The parser already reported. Do not report again.
		return c.invalid()

	case *ast.ParenExpr:
		// A parenthesized single type is that type — one node serves both.
		return c.typ(x.X)

	case *ast.Ident:
		return c.typeFromName(x, x)

	case *ast.SelectorExpr:
		return c.qualifiedType(x)

	case *ast.IndexExpr:
		return c.instantiate(x)

	case *ast.OwnershipType:
		return c.ownershipType(x)

	case *ast.ArrayType:
		if x.Len == nil {
			return types.NewSlice(c.typ(x.Elem))
		}
		n, ok := c.constInt(x.Len)
		if !ok {
			return c.invalid()
		}
		// §3.1 ⊢ `[N]T` carries N in its identity, which is why the length is
		// resolved here and never kept as an expression.
		return types.NewArray(c.typ(x.Elem), n)

	case *ast.MapType:
		key := c.typ(x.Key)
		if !types.IsInvalid(key) && !types.IsComparable(key) {
			c.errorExpr(x.Key, diag.NonComparableKey, types.TypeString(key))
		}
		return types.NewMap(key, c.typ(x.Value))

	case *ast.ChanType:
		return types.NewChan(c.typ(x.Elem))

	case *ast.PointerType:
		// §3.2 ⊢ a typed_ptr may not be the direct base of another; the nested
		// form is written parenthesized. The parser recurses unguarded so the
		// unparenthesized form parses and is named here.
		if _, nested := x.Elem.(*ast.PointerType); nested {
			c.errorExpr(x, diag.NestedPointer)
		}
		return types.NewPointer(c.typ(x.Elem))

	case *ast.FuncType:
		// A bare FunctionType has no receiver and no Expected result: that form
		// reaches the grammar only through a declaration.
		return c.signature(nil, x, false)

	case *ast.TupleExpr:
		return c.tupleType(x)

	case *ast.TensorType:
		// Legal only inside an npu-marked function body or that function's own
		// signature. The parser accepts it everywhere on purpose, so the
		// rejection is here and can name the construct.
		if !c.ctx.npu {
			c.errorExpr(x, diag.TensorOutsideNpu)
		}
		return c.tensorType(x)

	case *ast.VectorType:
		return c.vectorType(x)

	case *ast.AbstractType:
		// Legal only as an alias target. Reaching it through typ means it was
		// written inline; the alias path never calls this.
		c.errorExpr(x, diag.AbstractInline)
		return c.invalid()

	case *ast.UnaryExpr:
		if x.Op == token.TILDE {
			// `~T` is underlying-type inside a TypeSet and bitwise-NOT
			// everywhere else. The parser produces one node for both, so the
			// position is what separates them.
			c.errorExpr(x, diag.TildeOutsideSet)
			return c.invalid()
		}
	}

	c.errorExpr(e, diag.NotAType, exprString(e))
	return c.invalid()
}

// typeFromName resolves an identifier in type position.
//
// Three of the parser's punts land here at once: a predeclared type name is an
// ordinary identifier, a constraint name looks identical to a type name, and a
// generic name without arguments is an error rather than a type.
func (c *Checker) typeFromName(id *ast.Ident, at ast.Node) types.Type {
	obj := c.lookup(id)
	if obj == nil {
		return c.invalid()
	}

	tn, ok := obj.(*types.TypeName)
	if !ok {
		c.errorExpr(at, diag.NotAType, id.Name)
		return c.invalid()
	}
	if tn.IsConstraint() {
		c.errorExpr(at, diag.ConstraintAsType, id.Name)
		return c.invalid()
	}

	// Phase 2 may reach a name whose own type is not yet resolved. Resolving it
	// on demand is what makes order-independence work, and the resolving stack
	// is what turns a self-reference into a diagnostic instead of a hang.
	c.objDecl(tn)

	t := tn.Type()
	if t == nil {
		return c.invalid()
	}

	// §2.3's tensor element names carry their own legality rule, and it is a
	// rule about position rather than shape. Two of them:
	//
	//   - a bare tensor element name is legal only inside an npu body;
	//   - and never in a signature at all.
	//
	// The element slot of a TensorType is exempt from both, since `tensor[bf16,
	// 4]` is precisely the form the npu signature rule licenses. That reading
	// is this implementation's, not a stated one.
	if types.IsTensorElem(t) && !c.inTensorElem {
		switch {
		case c.inSignature:
			c.errorExpr(at, diag.TensorElemInSignature, id.Name)
		case !c.ctx.npu:
			c.errorExpr(at, diag.TensorElemOutsideNpu, id.Name)
		}
	}

	// A generic name used bare is an error: inference works from value
	// arguments, and a type position has none.
	if n := types.AsNamed(t); n != nil && len(n.TypeParams()) > 0 && len(n.TypeArgs()) == 0 {
		c.errorExpr(at, diag.WrongTypeArgCount, id.Name, len(n.TypeParams()), 0)
		return c.invalid()
	}
	return t
}

// qualifiedType resolves `pkg.Type`.
func (c *Checker) qualifiedType(x *ast.SelectorExpr) types.Type {
	id, ok := x.X.(*ast.Ident)
	if !ok {
		c.errorExpr(x, diag.NotAType, exprString(x))
		return c.invalid()
	}
	obj := c.lookup(id)
	if obj == nil {
		return c.invalid()
	}
	pkgName, ok := obj.(*types.PkgName)
	if !ok {
		c.errorExpr(x, diag.NotAType, exprString(x))
		return c.invalid()
	}

	// §1.3 ⊢ the qualifier is the imported package's own name; the path is a
	// locator, not a name.
	member := pkgName.Imported().Scope().Lookup(x.Sel.Name)
	if member == nil {
		c.errorExpr(x.Sel, diag.UndeclaredName, x.Sel.Name)
		return c.invalid()
	}
	c.recordUse(x.Sel, member)
	c.info.RecordSelection(x, &types.Selection{Kind: types.PackageMember, Obj: member})

	tn, ok := member.(*types.TypeName)
	if !ok {
		c.errorExpr(x, diag.NotAType, x.Sel.Name)
		return c.invalid()
	}
	if tn.IsConstraint() {
		c.errorExpr(x, diag.ConstraintAsType, x.Sel.Name)
		return c.invalid()
	}
	if t := tn.Type(); t != nil {
		return t
	}
	return c.invalid()
}

// instantiate resolves `Stack[int32]` in type position.
//
// Index and TypeArgs are one node, distinguished by whether the operand denotes
// a generic declaration. In type position only the instantiation reading is
// grammatical, so the ambiguity does not arise here — it arises in expression
// position, which resolve.go settles by the same question.
func (c *Checker) instantiate(x *ast.IndexExpr) types.Type {
	base := c.typInternal(x.X)
	if types.IsInvalid(base) {
		return base
	}

	named := types.AsNamed(base)
	if named == nil || len(named.TypeParams()) == 0 {
		c.errorExpr(x.X, diag.NotGeneric, exprString(x.X))
		return c.invalid()
	}

	params := named.TypeParams()
	args := make([]types.Type, 0, len(x.Indices))
	for _, ix := range x.Indices {
		args = append(args, c.typ(ix))
	}

	if len(args) != len(params) {
		c.errorExpr(x, diag.WrongTypeArgCount, exprString(x.X), len(params), len(args))
		return c.invalid()
	}

	// §9.2 ⊢ constraint satisfaction is checked per instantiation, at the
	// instantiation site, and a failure is a compile error there. There is no
	// runtime counterpart.
	for i, arg := range args {
		if types.IsInvalid(arg) {
			continue
		}
		cst := params[i].Constraint()
		if cst != nil && !cst.Satisfies(arg) {
			c.errorExpr(x.Indices[i], diag.ConstraintNotSatisfied,
				types.TypeString(arg), types.ConstraintString(cst))
		}
	}

	inst := types.NewNamed(named.Obj(), named.Underlying())
	inst.SetTypeParams(params)
	inst.SetTypeArgs(args)

	if id, ok := x.X.(*ast.Ident); ok {
		c.info.RecordInstance(id, types.Instance{TypeArgs: args, Type: inst})
	}
	return inst
}

// ownershipType resolves the five qualifiers.
//
// Only unique/shared/weak produce a Type. §3.2 ⊢ mut and var are legal in a
// parameter or receiver position only, so reaching them through typ means they
// were written somewhere else — types.Mode is where the legal ones live, set by
// splitMode.
func (c *Checker) ownershipType(x *ast.OwnershipType) types.Type {
	var kind types.OwnKind
	switch x.Kw {
	case token.UNIQUE:
		kind = types.Unique
	case token.SHARED:
		kind = types.Shared
	case token.WEAK:
		kind = types.Weak
	case token.MUT, token.VAR:
		c.errorExpr(x, diag.MutOutsidePosition, x.Kw.Spelling())
		return c.typ(x.X)
	default:
		return c.invalid()
	}

	// §3.2 ⊢ qualifiers do not stack. The parser recurses unguarded so the
	// stacked form parses and is diagnosed as itself.
	if inner, ok := stripParens(x.X).(*ast.OwnershipType); ok {
		c.errorExpr(x, diag.StackedOwnership,
			x.Kw.Spelling()+" "+inner.Kw.Spelling())
	}
	return types.NewOwnership(kind, c.typ(x.X))
}

// tupleType resolves a TupleType. A named element arrives as a KeyValueExpr,
// since TupleElem is `[ identifier ":" ] Type`.
//
// grammar.md ⊢ a tuple has at least one element and there is no unit type, so
// an empty one is not a TupleType at all — the parser already reported it, and
// there is nothing here to build.
func (c *Checker) tupleType(x *ast.TupleExpr) types.Type {
	if len(x.Elems) == 0 {
		return c.invalid()
	}
	vars := make([]*types.Var, 0, len(x.Elems))
	for _, e := range x.Elems {
		if kv, ok := e.(*ast.KeyValueExpr); ok {
			name := ""
			if id, ok := kv.Key.(*ast.Ident); ok {
				name = id.Name
			}
			vars = append(vars, types.NewVar(kv.Pos(), c.pkg, name, c.typ(kv.Value)))
			continue
		}
		vars = append(vars, types.NewVar(e.Pos(), c.pkg, "", c.typ(e)))
	}
	return types.NewTuple(vars...)
}

func (c *Checker) tensorType(x *ast.TensorType) types.Type {
	saved := c.inTensorElem
	c.inTensorElem = true
	elem := c.typ(x.Elem)
	c.inTensorElem = saved

	shape := make([]int64, 0, len(x.Shape))
	for _, d := range x.Shape {
		n, ok := c.constInt(d)
		if !ok {
			return c.invalid()
		}
		shape = append(shape, n)
	}
	return types.NewTensor(elem, shape)
}

func (c *Checker) vectorType(x *ast.VectorType) types.Type {
	elem := c.typ(x.Elem)
	n, ok := c.constInt(x.Len)
	if !ok {
		return c.invalid()
	}
	return types.NewVector(elem, n)
}

// ------------------------------------------------------------- constraints

// constraintExpr converts a constraint element to a *Constraint.
//
// A single identifier here is the fork grammar.md calls out: it parses as both a
// one-term TypeSet and a constraint name, and is resolved by what the name
// denotes. That resolution is a scope lookup, and it happens below.
func (c *Checker) constraintExpr(e ast.Expr) *types.Constraint {
	if e == nil {
		// A bare TypeParamName is constrained by `any`.
		return types.Any
	}

	var terms []types.Term
	var embeds []*types.Constraint

	var walk func(ast.Expr)
	walk = func(x ast.Expr) {
		switch t := stripParens(x).(type) {
		case *ast.BinaryExpr:
			if t.Op == token.OR {
				// `|` within one element is a union: the set is every listed
				// type.
				walk(t.X)
				walk(t.Y)
				return
			}

		case *ast.UnaryExpr:
			if t.Op == token.TILDE {
				// `~T` admits every type whose underlying type is T; a bare T
				// admits only T exactly.
				terms = append(terms, types.Term{Tilde: true, Type: c.typ(t.X)})
				return
			}

		case *ast.Ident:
			obj := c.lookup(t)
			if obj == nil {
				return
			}
			tn, ok := obj.(*types.TypeName)
			if !ok {
				c.errorExpr(t, diag.NotAConstraint, t.Name)
				return
			}
			if tn.IsConstraint() {
				// A bare constraint name embeds that constraint's set.
				c.objDecl(tn)
				embeds = append(embeds, tn.Constraint())
				return
			}
			terms = append(terms, types.Term{Type: c.typeFromName(t, t)})
			return
		}
		terms = append(terms, types.Term{Type: c.typ(x)})
	}
	walk(e)

	if len(terms) == 0 && len(embeds) == 1 {
		return embeds[0]
	}
	return types.NewConstraint(nil, terms, nil, embeds)
}

// ---------------------------------------------------------------- helpers

func stripParens(e ast.Expr) ast.Expr {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.X
	}
}

// exprString renders an expression for a diagnostic argument. It is
// deliberately shallow — a full printer belongs to the printer package, and a
// diagnostic wants the name the reader wrote, not a reconstruction.
func exprString(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return exprString(x.X) + "." + x.Sel.Name
	case *ast.IndexExpr:
		return exprString(x.X) + "[...]"
	case *ast.TupleIndexExpr:
		return exprString(x.X) + "." + x.Text
	case *ast.CallExpr:
		return exprString(x.Fun) + "(...)"
	case *ast.ParenExpr:
		return exprString(x.X)
	case *ast.BasicLit:
		return x.Value
	}
	return "expression"
}