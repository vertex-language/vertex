package analyzer

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// ----------------------------------------------- phase 3: walk the bodies

// checkBodies resolves every identifier inside a function body to an object.
//
// This pass records Uses and opens scopes; it computes no types for
// expressions. That split is deliberate — resolution answers "what does this
// name denote", and several of the parser's ambiguities (A.3.6's a[i] vs.
// Stack[int32], A.4.4's & as address-of vs. dereference) are answerable from
// the answer alone, before any expression typing exists to complicate them.
func (c *Checker) checkBodies() {
	for _, f := range c.files {
		fileScope := c.fileScopeOf(f)
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}
			c.scope = fileScope
			c.funcBody(fd)
		}
	}
	c.scope = nil
}

func (c *Checker) funcBody(d *ast.FuncDecl) {
	sig := c.signatureOf(d)

	// A.6.1 ⊢ npu "sets [+Npu], which licenses tensor types and the npu.
	// namespace and simultaneously restricts the body." This is the one A.0.2
	// context parameter the parser did not track, so it is tracked here.
	savedNpu := c.npu
	c.npu = sig != nil && sig.Marker() == types.MarkerNPU
	defer func() { c.npu = savedNpu }()

	c.openScope(d, "function "+d.Name.Name)
	defer c.closeScope()

	// A.0.3 ⊢ propagation stops at a function boundary, and the receiver and
	// parameters belong to the body's own scope rather than to the file's.
	if sig != nil {
		if r := sig.Recv(); r != nil {
			c.scope.Insert(r)
		}
		for i := 0; i < sig.Params().Len(); i++ {
			p := sig.Params().At(i)
			if p.Name() != "" && p.Name() != "_" {
				c.scope.Insert(p)
			}
		}
	}
	c.blockBody(d.Body)
}

func (c *Checker) signatureOf(d *ast.FuncDecl) *types.Signature {
	if obj := c.info.Defs[d.Name]; obj != nil {
		if f, ok := obj.(*types.Func); ok {
			c.objDecl(f)
			return f.Signature()
		}
	}
	return nil
}

// blockBody walks a block's statements in the already-open scope. It exists
// separately from block() so a function body does not nest one scope inside
// another for no reason.
func (c *Checker) blockBody(b *ast.BlockStmt) {
	if b == nil {
		return
	}
	for _, s := range b.List {
		c.stmt(s)
	}
}

func (c *Checker) block(b *ast.BlockStmt) {
	if b == nil {
		return
	}
	c.openScope(b, "block")
	defer c.closeScope()
	c.blockBody(b)
}

func (c *Checker) stmt(s ast.Stmt) {
	switch x := s.(type) {
	case nil, *ast.BadStmt:
		return

	case *ast.BlockStmt:
		c.block(x)

	case *ast.DeclStmt:
		if vd, ok := x.Decl.(*ast.VarDecl); ok {
			c.localVarDecl(vd)
		}

	case *ast.ExprStmt:
		c.expr(x.X)

	case *ast.AssignStmt:
		// A.5.2 ⊢ assignment is a statement, never an expression, so both
		// sides are ordinary expressions here and there is no binding to
		// introduce.
		for _, t := range x.Targets {
			c.expr(t)
		}
		for _, v := range x.Values {
			c.expr(v)
		}

	case *ast.IfStmt:
		c.expr(x.Cond)
		c.block(x.Body)
		c.stmt(x.Else)

	case *ast.WhileStmt:
		c.expr(x.Cond)
		c.block(x.Body)

	case *ast.ForStmt:
		c.forStmt(x)

	case *ast.SwitchStmt:
		c.expr(x.Tag)
		for _, cl := range x.Cases {
			c.caseClause(cl)
		}

	case *ast.SelectStmt:
		for _, cl := range x.Cases {
			c.selectClause(cl)
		}

	case *ast.ReturnStmt:
		for _, r := range x.Results {
			c.expr(r)
		}

	case *ast.DeferStmt:
		// A.5.8 ⊢ "its arguments are evaluated at registration; only the call
		// itself is postponed", so the call is resolved here like any other.
		c.expr(x.Call)
		if c.info.Defers != nil {
			c.info.Defers[x] = c.scope
		}

	case *ast.BranchStmt:
		// A.5.9 ⊢ there are no loop labels, so nothing to resolve.
	}
}

// localVarDecl introduces A.5.1's bindings.
//
// A.5.1 ⊢ "let is immutable and not guaranteed to be addressable — it may be a
// register, an SSA value, or folded away entirely. var is mutable and owns a
// real stack slot for its whole lifetime." Mutable is what A.3.2's "a mut
// argument must be an addressable var binding" later reads.
func (c *Checker) localVarDecl(d *ast.VarDecl) {
	// Initializers are resolved first: a binding is not in scope inside its own
	// initializer, so `let x = x` refers to an outer x or is undeclared.
	for _, v := range d.Values {
		c.expr(v)
	}
	for _, b := range d.Bindings {
		var t types.Type
		if b.Type != nil {
			t = c.typ(b.Type)
		}
		obj := types.NewVar(b.Name.Pos(), c.pkg, b.Name.Name, t)
		obj.SetMutable(d.Kw == token.VAR)
		c.declare(c.scope, b.Name, obj)
	}
}

// forStmt resolves A.5.6's single loop shape.
//
// The mode marker sits on the binding rather than the iterable because, as
// A.5.6 puts it, "what is being transferred is each element, one per iteration,
// and the marker names what moves." The parser accepts `mut a, b` even though
// no production combines them (A.14), so that combination is rejected here.
//
// The bare (shared-access) form carries token.INVALID — the zero Kind — since
// there is no marker token to record.
func (c *Checker) forStmt(x *ast.ForStmt) {
	c.expr(x.X)

	c.openScope(x, "for")
	defer c.closeScope()

	if x.Mode != token.INVALID && len(x.Names) > 1 {
		c.errorAt(x.ModePos, x.ModePos, diag.MutOutsidePosition, x.Mode.String())
	}
	for _, n := range x.Names {
		obj := types.NewVar(n.Pos(), c.pkg, n.Name, nil)
		// The `mut` form iterates by exclusive access and the `var` form
		// consumes; both give the binding a real slot. The bare form is a
		// shared view and does not.
		obj.SetMutable(x.Mode == token.MUT || x.Mode == token.VAR)
		c.declare(c.scope, n, obj)
	}
	c.blockBody(x.Body)
}

func (c *Checker) caseClause(cl *ast.CaseClause) {
	c.openScope(cl, "case")
	defer c.closeScope()

	for _, p := range cl.Patterns {
		c.pattern(p)
	}
	for _, s := range cl.Body {
		c.stmt(s)
	}
}

// pattern resolves a Pattern (A.5.7). An EnumPattern's payload entries are
// binding names rather than expressions — A.5.7 ⊢ they are "a view into the
// payload, not a copy" — which is why they are declared here rather than
// resolved as uses.
func (c *Checker) pattern(p ast.Expr) {
	ep, ok := p.(*ast.EnumPattern)
	if !ok {
		c.expr(p)
		return
	}
	for _, b := range ep.Binds {
		obj := types.NewVar(b.Pos(), c.pkg, b.Name, nil)
		c.declare(c.scope, b, obj)
	}
}

func (c *Checker) selectClause(cl *ast.SelectClause) {
	c.openScope(cl, "select case")
	defer c.closeScope()

	for _, t := range cl.Targets {
		c.expr(t)
	}
	c.expr(cl.Op)
	for _, s := range cl.Body {
		c.stmt(s)
	}
}

// ------------------------------------------------------------ expressions

// expr resolves every identifier in an expression to an object.
//
// It computes no types. What it does do is answer the questions the parser
// deferred that turn purely on what a name denotes — which is most of them.
func (c *Checker) expr(e ast.Expr) {
	switch x := e.(type) {
	case nil, *ast.BadExpr, *ast.BasicLit, *ast.NamespaceExpr, *ast.AbstractType:
		return

	case *ast.Ident:
		c.lookup(x)

	case *ast.ParenExpr:
		c.expr(x.X)

	case *ast.SelectorExpr:
		// A selector's base may be a package qualifier, a namespace, or an
		// ordinary value. Only the first is resolvable without types, so the
		// rest is typecheck's; here the base is resolved and the member left.
		c.expr(x.X)
		if id, ok := x.X.(*ast.Ident); ok {
			if obj := c.info.Uses[id]; obj != nil {
				if pn, ok := obj.(*types.PkgName); ok {
					if m := pn.Imported().Scope().Lookup(x.Sel.Name); m != nil {
						c.info.RecordUse(x.Sel, m)
						c.info.RecordSelection(x, &types.Selection{
							Kind: types.PackageMember, Obj: m,
						})
					} else {
						c.errorExpr(x.Sel, diag.UndeclaredName, x.Sel.Name)
					}
				}
			}
		}

	case *ast.TupleIndexExpr:
		c.expr(x.X)

	case *ast.IndexExpr:
		// A.3.6's ambiguity. If X denotes a type, the brackets are a type
		// argument list; otherwise they are an index or a slice. Resolution is
		// exactly what tells them apart, which is why the annex says it is
		// "resolved by the operand's meaning rather than by shape."
		c.expr(x.X)
		if c.denotesType(x.X) {
			c.typ(x)
			return
		}
		for _, ix := range x.Indices {
			c.expr(ix)
		}

	case *ast.CallExpr:
		c.expr(x.Fun)
		c.callArgs(x)

	case *ast.LaunchExpr:
		if x.Config != nil {
			c.expr(x.Config.Blocks)
			c.expr(x.Config.Threads)
		}
		c.expr(x.Call)

	case *ast.AwaitExpr:
		c.expr(x.X)

	case *ast.UnaryExpr:
		c.expr(x.X)

	case *ast.BinaryExpr:
		c.expr(x.X)
		c.expr(x.Y)

	case *ast.CastExpr:
		c.expr(x.X)
		c.typ(x.Type)

	case *ast.TransferExpr:
		// A.4.6 ⊢ "the marker takes a binding or a field path and nothing
		// else." A.14 ✗ `let y = var pick(a, b)`. The parser accepts any unary
		// expression so this can name the rule; whether the position is an
		// owning one (A.9.1) is the ownership pass's question, not this one's.
		c.expr(x.Target)

	case *ast.KeyValueExpr:
		c.expr(x.Key)
		c.expr(x.Value)

	case *ast.TupleExpr:
		for _, el := range x.Elems {
			c.expr(el)
		}

	case *ast.ArrayLit:
		for _, el := range x.Elems {
			c.expr(el)
		}

	case *ast.CompositeLit:
		c.typ(x.Type)
		for _, el := range x.Elems {
			// A.4.7 ⊢ a FieldValue's key "must be an Identifier" naming a
			// field, so it is not resolved as an ordinary use.
			if kv, ok := el.(*ast.KeyValueExpr); ok {
				c.expr(kv.Value)
				continue
			}
			c.expr(el)
		}

	case *ast.MapLit:
		for _, el := range x.Elems {
			c.expr(el)
		}

	case *ast.EnumShorthand:
		// A.4.1 ⊢ legal only where the enum type is inferable from context,
		// which needs types. The name is left for typecheck.
		for _, a := range x.Args {
			c.expr(a)
		}

	case *ast.FuncLit:
		c.funcLit(x)

	case *ast.OwnershipType, *ast.ArrayType, *ast.MapType, *ast.ChanType,
		*ast.PointerType, *ast.FuncType, *ast.TensorType:
		// A type reached in expression position — a type argument or a sizeof
		// operand. The parser produces the same nodes either way.
		c.typ(e)
	}
}

// callArgs resolves a call's arguments, with one special case.
//
// A.4.8 ⊢ sizeof, alignof, and reinterpret take a Type in argument position.
// The parser already parsed the first argument as a type by recognizing the
// name — sound only because A.1.4 forbids shadowing a ReservedBuiltinName — so
// this resolves it as one rather than as an expression.
func (c *Checker) callArgs(x *ast.CallExpr) {
	typeFirst := false
	if id, ok := x.Fun.(*ast.Ident); ok {
		switch id.Name {
		case "sizeof", "alignof", "reinterpret":
			typeFirst = true
		}
	}
	for i, a := range x.Args {
		if typeFirst && i == 0 {
			c.typ(a)
			continue
		}
		if kv, ok := a.(*ast.KeyValueExpr); ok {
			// A.4.3's named argument. The key names a parameter, not a
			// binding, so only the value is resolved here.
			c.expr(kv.Value)
			continue
		}
		c.expr(a)
	}
}

// funcLit walks a FunctionExpression body (A.4.1).
//
// A.0.3 ⊢ propagation stops at a function boundary: a literal "begins with all
// four parameters cleared, and re-sets them from its own marker list — an
// anonymous closure written inside an async body may not await unless it is
// itself marked async." Clearing npu here is that rule.
func (c *Checker) funcLit(x *ast.FuncLit) {
	savedNpu := c.npu
	c.npu = x.Type != nil && x.Type.Marker != nil && x.Type.Marker.Name == "npu"
	defer func() { c.npu = savedNpu }()

	c.openScope(x, "function literal")
	defer c.closeScope()

	params, _ := c.paramList(x.Type.Params)
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		if p.Name() != "" && p.Name() != "_" {
			c.scope.Insert(p)
		}
	}
	if x.Type.Result != nil {
		c.typ(x.Type.Result)
	}
	c.blockBody(x.Body)
}

// denotesType reports whether an expression names a type rather than a value.
// It is the A.3.6 fork, and it is answerable from resolution alone.
func (c *Checker) denotesType(e ast.Expr) bool {
	switch x := stripParens(e).(type) {
	case *ast.Ident:
		obj := c.info.Uses[x]
		if obj == nil {
			return false
		}
		tn, ok := obj.(*types.TypeName)
		return ok && !tn.IsConstraint()
	case *ast.SelectorExpr:
		obj := c.info.Uses[x.Sel]
		_, ok := obj.(*types.TypeName)
		return ok
	}
	return false
}