package analyzer

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// typ converts a type expression to a types.Type.
//
// There are two entry points into this file, not one, and the split is
// load-bearing. A.7.2 ⊢ "a constraint is never a value type and is legal only
// in a [...] position", and types.Constraint deliberately does not implement
// types.Type. So a constraint position calls constraintExpr and a type position
// calls typ, and A.14's `var c: Ordered` falls out as a diagnostic from typ
// rather than needing a predicate someone must remember to call.
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
		return c.typInternal(x.X)

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
		n, ok := c.arrayLen(x.Len)
		if !ok {
			return c.invalid()
		}
		return types.NewArray(c.typ(x.Elem), n)

	case *ast.MapType:
		key := c.typ(x.Key)
		// A.3.1 ⊢ map[K]V requires K to satisfy `comparable` (A.7.4).
		if !types.IsInvalid(key) && !types.IsComparable(key) {
			c.errorExpr(x.Key, diag.NonComparableKey, types.TypeString(key))
		}
		return types.NewMap(key, c.typ(x.Value))

	case *ast.ChanType:
		return types.NewChan(c.typ(x.Elem))

	case *ast.PointerType:
		return types.NewPointer(c.typ(x.Elem))

	case *ast.FuncType:
		// A.3.4 ⊢ a bare FunctionType has no receiver; the marker comes from
		// the node itself, which signature reads.
		return c.signature(nil, x)

	case *ast.TupleExpr:
		return c.tupleType(x)

	case *ast.TensorType:
		// A.3.5 ⊢ TensorType is grammatical only under [+Npu]. The parser
		// accepts it everywhere on purpose (A.14 lists it among the forms that
		// parse and are rejected later), so the rejection is here.
		if !c.npu {
			c.errorExpr(x, diag.TensorOutsideNpu)
		}
		return c.tensorType(x)

	case *ast.AbstractType:
		// A.3.3 ⊢ `abstract` appears only on the right-hand side of a
		// TypeAliasDeclaration. Reaching it through typ means it was written
		// inline; the alias path never calls this.
		c.errorExpr(x, diag.AbstractInline)
		return c.invalid()

	case *ast.UnaryExpr:
		if x.Op == token.TILDE {
			// A.7.3 ✗ `type X = ~int` — `~T` is valid only inside a type set.
			// The parser produces the same node for both readings, so the
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
// This is where three of the parser's punts land at once: a PredeclaredTypeName
// is an ordinary identifier (A.1.4), a ConstraintName looks identical to a
// TypeName (A.7.2), and a generic name without arguments is an error rather
// than a type.
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

	// A.7.2 ⊢ a constraint "is never a value type and is legal only in a [...]
	// position." A.14 lists `var c: Ordered` among the rejected forms.
	if tn.IsConstraint() {
		c.errorExpr(at, diag.ConstraintAsType, id.Name)
		return c.invalid()
	}

	// Phase 2 may reach a name whose own type is not yet resolved. Resolving it
	// on demand is what makes A.2's order-independence work, and the resolving
	// stack is what turns a self-reference into a diagnostic instead of a hang.
	c.objDecl(tn)

	t := tn.Type()
	if t == nil {
		return c.invalid()
	}

	// A generic name used bare is an error: A.7.5 ⊢ inference works from value
	// arguments, and a type position has none.
	if n := types.AsNamed(t); n != nil && len(n.TypeParams()) > 0 && len(n.TypeArgs()) == 0 {
		c.errorExpr(at, diag.WrongTypeArgCount, id.Name, len(n.TypeParams()), 0)
		return c.invalid()
	}
	return t
}

// qualifiedType resolves `pkg.Type` (A.2.3's QualifiedTypeName).
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

	// A.2.3 ⊢ "the imported package's declared name is the qualifier under
	// which its symbols are reached; the import path is a locator, not a name."
	member := pkgName.Imported().Scope().Lookup(x.Sel.Name)
	if member == nil {
		c.errorExpr(x.Sel, diag.UndeclaredName, x.Sel.Name)
		return c.invalid()
	}
	c.info.RecordUse(x.Sel, member)

	tn, ok := member.(*types.TypeName)
	if !ok || tn.IsConstraint() {
		c.errorExpr(x, diag.NotAType, x.Sel.Name)
		return c.invalid()
	}
	return tn.Type()
}

// instantiate resolves `Stack[int32]` in type position.
//
// A.3.6 ⊢ this shape and `a[i]` "share bracket syntax and are distinguished by
// whether the operand names a generic declaration. This is the one syntactic
// overlap resolved by the operand's meaning rather than by shape." In type
// position only the instantiation reading is grammatical, so the ambiguity does
// not arise here — it arises in expression position, which typecheck handles.
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

	// A.7.5 ⊢ "constraint satisfaction is checked once per instantiation, at
	// the instantiation site, never at runtime."
	for i, arg := range args {
		if types.IsInvalid(arg) {
			continue
		}
		if cst := params[i].Constraint(); cst != nil && !cst.Satisfies(arg) {
			c.errorExpr(x.Indices[i], diag.NotAConstraint, types.TypeString(arg))
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

// ownershipType resolves A.3.2's five qualifiers.
//
// Only unique/shared/weak produce a Type. A.3.2 ⊢ mut and var are "legal only
// in a parameter or receiver position", so reaching them through typ means they
// were written somewhere else — types.Mode is where the legal ones live, set by
// paramVar below.
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
		c.errorExpr(x, diag.MutOutsidePosition, x.Kw.String())
		return c.typ(x.X)
	default:
		return c.invalid()
	}

	// A.3.2 ⊢ "qualifiers do not stack: mut var T, mut mut T, and
	// shared unique T are compile errors." The parser recurses unguarded
	// (A.14 lists the stacked form as parsing), so the check is here.
	if inner, ok := stripParens(x.X).(*ast.OwnershipType); ok {
		c.errorExpr(x, diag.StackedOwnership,
			x.Kw.String()+" "+inner.Kw.String())
	}
	return types.NewOwnership(kind, c.typ(x.X))
}

func (c *Checker) tupleType(x *ast.TupleExpr) types.Type {
	// A.3.1 ⊢ `()` is the unit type: zero bytes, one value.
	if len(x.Elems) == 0 {
		return types.Unit
	}
	vars := make([]*types.Var, 0, len(x.Elems))
	for _, e := range x.Elems {
		// A.3.1's TupleElement is `Type | Identifier : Type`, so a named
		// element arrives as a KeyValueExpr.
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
	elem := c.typ(x.Elem)
	shape := make([]int64, 0, len(x.Shape))
	for _, d := range x.Shape {
		n, ok := c.arrayLen(d)
		if !ok {
			return c.invalid()
		}
		shape = append(shape, n)
	}
	return types.NewTensor(elem, shape)
}

// arrayLen evaluates an ArrayLength (A.3.1) or a tensor shape entry (A.3.5).
//
// A.3.1 ⊢ "ArrayLength must be a compile-time constant." Only a literal and a
// constant identifier are handled here; a constant *expression* needs the full
// evaluator, which belongs to typecheck. Anything else is diagnosed rather than
// silently accepted, so widening this later cannot change a passing program.
func (c *Checker) arrayLen(e ast.Expr) (int64, bool) {
	switch x := stripParens(e).(type) {
	case *ast.BasicLit:
		if x.Kind != token.INT {
			break
		}
		v, ok := parseIntLit(x.Value)
		if !ok {
			break
		}
		if v < 0 {
			c.errorExpr(e, diag.ArrayLenNegative)
			return 0, false
		}
		return v, true

	case *ast.Ident:
		obj := c.lookup(x)
		if obj == nil {
			return 0, false
		}
		cn, ok := obj.(*types.Const)
		if !ok {
			break
		}
		if n, ok := types.Int64Val(cn.Val()); ok {
			if n < 0 {
				c.errorExpr(e, diag.ArrayLenNegative)
				return 0, false
			}
			return n, true
		}
	}
	c.errorExpr(e, diag.ArrayLenNotConst)
	return 0, false
}

// ------------------------------------------------------------- constraints

// constraintExpr converts a ConstraintExpression (A.7.3) to a *Constraint.
//
// It is the second entry point into this file. A single identifier here is the
// case A.7.2 calls out: it "parses as both a TypeSet of one term and a
// ConstraintName; it is resolved by what the name denotes." That resolution is
// a scope lookup, and it happens below.
func (c *Checker) constraintExpr(e ast.Expr) *types.Constraint {
	if e == nil {
		// A.7.1 ⊢ "a bare name is constraint `any`: [T] means [T: any]."
		return types.Any
	}

	var terms []types.Term
	var embeds []*types.Constraint

	var walk func(ast.Expr)
	walk = func(x ast.Expr) {
		switch t := stripParens(x).(type) {
		case *ast.BinaryExpr:
			if t.Op == token.OR {
				// A.7.3 ⊢ `|` is union: the type set is every listed type.
				walk(t.X)
				walk(t.Y)
				return
			}
		case *ast.UnaryExpr:
			if t.Op == token.TILDE {
				// A.7.3 ⊢ `~T` admits T and every type whose underlying type
				// is T, so an alias to float32 still satisfies ~float32.
				terms = append(terms, types.Term{Tilde: true, Type: c.typ(t.X)})
				return
			}
		case *ast.Ident:
			// The A.7.2 fork. A constraint name embeds its set; anything else
			// is a one-term type set.
			if obj := c.lookup(t); obj != nil {
				if tn, ok := obj.(*types.TypeName); ok && tn.IsConstraint() {
					embeds = append(embeds, tn.Constraint())
					return
				}
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
	case *ast.ParenExpr:
		return exprString(x.X)
	}
	return "expression"
}

// parseIntLit decodes an IntegerLiteral's source spelling (A.1.5.1), including
// its base prefix and `_` separators. The scanner already validated the shape,
// so this only has to read it.
func parseIntLit(s string) (int64, bool) {
	base := int64(10)
	if len(s) > 2 && s[0] == '0' {
		switch s[1] {
		case 'b', 'B':
			base, s = 2, s[2:]
		case 'o', 'O':
			base, s = 8, s[2:]
		case 'x', 'X':
			base, s = 16, s[2:]
		}
	}
	var n int64
	any := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '_' {
			continue
		}
		var d int64
		switch {
		case '0' <= ch && ch <= '9':
			d = int64(ch - '0')
		case 'a' <= ch && ch <= 'f':
			d = int64(ch-'a') + 10
		case 'A' <= ch && ch <= 'F':
			d = int64(ch-'A') + 10
		default:
			return 0, false
		}
		if d >= base {
			return 0, false
		}
		n = n*base + d
		any = true
	}
	return n, any
}