package analyzer

import (
	"strings"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// ----------------------------------------------- phase 3: walk the bodies

// checkBodies resolves every identifier inside a function body to an object.
//
// This pass records Uses and opens scopes; it computes no expression types.
// That split is deliberate — resolution answers "what does this name denote",
// and several of the parser's deferred ambiguities are answerable from that
// answer alone, before any expression typing exists to complicate them.
func (c *Checker) checkBodies() {
	for _, f := range c.files {
		fileScope := c.fileScopeOf(f)
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}
			c.funcBody(fd, fileScope)
		}
	}
	c.scope = nil
}

func (c *Checker) funcBody(d *ast.FuncDecl, fileScope *types.Scope) {
	obj := c.funcObj[d]
	var sig *types.Signature
	if obj != nil {
		c.objDecl(obj)
		sig = obj.Signature()
	}

	// The marker is what licenses `tensor` and `await`, and it is read from the
	// signature rather than from the tree because the marker is part of the
	// type. §1.4 makes `main` the one non-async function in which `await` is
	// legal, which is why IsEntry is consulted here and not only for its shape.
	saved := c.ctx
	c.ctx = bodyCtx{}
	if sig != nil {
		c.ctx.npu = sig.Marker() == types.MarkerNPU
		c.ctx.async = sig.Marker() == types.MarkerAsync
	}
	if obj != nil && obj.IsEntry() {
		c.ctx.async = true
	}
	defer func() { c.ctx = saved }()

	c.scope = fileScope
	c.openFuncScope(d, "function "+d.Name.Name)
	defer c.closeScope()

	// The receiver and the parameters belong to the body's own scope, not the
	// file's.
	if sig != nil {
		if r := sig.Recv(); r != nil && r.Name() != "" && r.Name() != "_" {
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

// blockBody walks a block's statements in the already-open scope. It exists
// separately from block so a function body does not nest one scope inside
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
	case nil, *ast.BadStmt, *ast.BranchStmt:
		// There are no loop labels, so a branch carries nothing to resolve.
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
		c.assignStmt(x)

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
		c.selectStmt(x)

	case *ast.ReturnStmt:
		for _, r := range x.Results {
			c.expr(r)
		}

	case *ast.DeferStmt:
		// Arguments are evaluated at registration; only the call is postponed,
		// so the call resolves here like any other.
		c.expr(x.Call)
		// Deferred calls run on every exit edge of the enclosing *function*,
		// which is the scope they are grouped by.
		c.info.RecordDefer(x, c.scope.FuncScope())
	}
}

// assignStmt resolves both assignment forms.
//
// Assignment is a statement and never an expression, so both sides are ordinary
// expressions and there is no binding to introduce. Which PrimaryExpr shapes
// are assignable is a static rule, and the shape half of it is decidable here:
// whether the *binding* underneath is assignable needs types and is not.
func (c *Checker) assignStmt(x *ast.AssignStmt) {
	for _, t := range x.Targets {
		if !assignableShape(t) {
			c.errorExpr(t, diag.NotAssignable, exprString(t))
		}
		if id, ok := stripParens(t).(*ast.Ident); ok && id.IsBlank() {
			// `_` names nothing, but it is a legal assignment target: the value
			// is evaluated and discarded. Resolving it would report a read.
			c.recordDef(id, nil)
			continue
		}
		c.expr(t)
	}
	for _, v := range x.Values {
		c.owning(v)
	}
}

// assignableShape reports whether e has the shape of an AssignTarget: a bare
// PrimaryExpr. A dereference-write `&p = 99` is a UnaryExpr and the blank
// identifier is an ordinary Ident; neither needs its own shape, and neither is
// excluded here.
func assignableShape(e ast.Expr) bool {
	switch x := stripParens(e).(type) {
	case *ast.Ident, *ast.SelectorExpr, *ast.IndexExpr, *ast.TupleIndexExpr:
		return true
	case *ast.UnaryExpr:
		return x.Op == token.AND
	}
	return false
}

// localVarDecl introduces a block-scoped binding.
//
// `let` fixes the binding and is not guaranteed to be addressable — it may be a
// register, an SSA value, or folded away entirely. `var` may be rebound and is
// what anything taking exclusive access or transferring requires. Mutable is
// what the `mut`-argument rule later reads.
func (c *Checker) localVarDecl(d *ast.VarDecl) {
	// Initializers are resolved first: a binding is not in scope inside its own
	// initializer, so `let x = x` refers to an outer x or is undeclared.
	for _, v := range d.Values {
		c.owning(v)
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

// forStmt resolves the single loop shape over an iterable.
//
// The mode marker sits on the binding rather than on the iterable, because what
// transfers is each element, one per iteration. The marker and the two-name
// form do not combine, but both parse together so the combination is diagnosed
// as itself rather than as a syntax error at the comma.
func (c *Checker) forStmt(x *ast.ForStmt) {
	c.expr(x.X)

	c.openScope(x, "for")
	defer c.closeScope()

	if x.Mode != token.INVALID && len(x.Names) > 1 {
		c.errorAt(x.ModePos, x.ModePos+token.Pos(len(x.Mode.Spelling())),
			diag.IterationBindingMode, x.Mode.Spelling())
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
	// A clause's statement list is its own scope, which is what lets a pattern
	// binding be visible in the body and nowhere else.
	c.openScope(cl, "case")
	defer c.closeScope()

	for _, p := range cl.Patterns {
		c.pattern(p)
	}
	for _, s := range cl.Body {
		c.stmt(s)
	}
}

// pattern resolves one Pattern. In pattern position a leading dot is always an
// enum pattern and never an enum shorthand reached through Expression: the
// payload entries are binding names rather than expressions, and they are views
// into the payload rather than copies — which is why they are declared here
// rather than resolved as uses.
func (c *Checker) pattern(p ast.Expr) {
	ep, ok := p.(*ast.EnumPattern)
	if !ok {
		c.expr(p)
		return
	}
	for _, b := range ep.Binds {
		c.declare(c.scope, b, types.NewVar(b.Pos(), c.pkg, b.Name, nil))
	}
}

// selectStmt resolves a select and checks the one shape rule decidable without
// types: a select is entirely bare or entirely awaited. Which calls are
// admissible in a ChannelCase is a question about their types and is not.
func (c *Checker) selectStmt(x *ast.SelectStmt) {
	awaited, bare := 0, 0
	for _, cl := range x.Cases {
		if cl.Op == nil {
			continue // the default clause
		}
		if _, ok := cl.Op.(*ast.AwaitExpr); ok {
			awaited++
		} else {
			bare++
		}
	}
	if awaited > 0 && bare > 0 {
		c.errorAt(x.Select, x.Select+token.Pos(len(token.SELECT.Spelling())),
			diag.SelectMixedAwait)
	}

	for _, cl := range x.Cases {
		c.selectClause(cl)
	}
}

// selectClause resolves one clause, covering all three channel-case forms. The
// declaring form introduces bindings scoped to this clause's body; the
// assigning form writes to pre-declared targets; the bare form has neither.
func (c *Checker) selectClause(cl *ast.SelectClause) {
	c.openScope(cl, "select case")
	defer c.closeScope()

	// The operation is resolved before the bindings, for the same reason a
	// local declaration's initializer is: a binding is not in scope inside the
	// expression that produces its value.
	for _, t := range cl.Targets {
		c.expr(t)
	}
	if cl.Op != nil {
		// An awaited channel operation is licensed by the enclosing body just
		// as any other await is.
		c.expr(cl.Op)
	}
	for _, b := range cl.Bindings {
		var t types.Type
		if b.Type != nil {
			t = c.typ(b.Type)
		}
		obj := types.NewVar(b.Name.Pos(), c.pkg, b.Name.Name, t)
		obj.SetMutable(cl.Kw == token.VAR)
		c.declare(c.scope, b.Name, obj)
	}
	for _, s := range cl.Body {
		c.stmt(s)
	}
}

// ------------------------------------------------------------ expressions

// owning resolves an expression in one of the six owning positions: a `let` or
// `var` initializer, an assignment value, a call argument, an array literal
// element, a composite literal field value, and a tuple element. Those are
// exactly the positions ast marks as owning, and exactly where the transfer
// marker is legal.
//
// Everywhere else, expr reports the marker as out of position. That is the
// whole mechanism: one node, two contexts, and the difference between a move
// and a deep copy riding on which one it was written in.
func (c *Checker) owning(e ast.Expr) {
	tr, ok := e.(*ast.TransferExpr)
	if !ok {
		c.expr(e)
		return
	}
	c.transferExpr(tr)
}

// transferExpr checks a marker in a position that admits one.
//
// The operand must be a binding or a field path. The parser writes it as a full
// unary expression so `var f(a)` and `var items[0]` parse and arrive here as
// real nodes rather than as syntax errors, which is what lets the diagnostic
// name the rule.
func (c *Checker) transferExpr(x *ast.TransferExpr) {
	c.expr(x.Target)

	base := transferBase(x.Target)
	if base == nil {
		c.errorExpr(x, diag.TransferComputed)
		return
	}
	// The binding that dies here is what lower reads to emit a move instead of
	// a copy, so the marker is recorded against it rather than merely accepted.
	if v, ok := c.objectOf(base).(*types.Var); ok {
		c.info.RecordTransfer(x, v)
	}
}

// transferBase returns the identifier at the root of a binding or field path,
// or nil if the operand is neither. An index is not a field path.
func transferBase(e ast.Expr) *ast.Ident {
	switch x := stripParens(e).(type) {
	case *ast.Ident:
		return x
	case *ast.SelectorExpr:
		return transferBase(x.X)
	}
	return nil
}

// expr resolves every identifier in an expression to an object.
//
// It computes no types. What it does do is answer the questions the parser
// deferred that turn purely on what a name denotes — which is most of them.
func (c *Checker) expr(e ast.Expr) {
	switch x := e.(type) {
	case nil, *ast.BadExpr, *ast.BasicLit, *ast.NamespaceExpr:
		return

	case *ast.Ident:
		c.lookup(x)

	case *ast.ParenExpr:
		c.expr(x.X)

	case *ast.SelectorExpr:
		c.selectorExpr(x)

	case *ast.TupleIndexExpr:
		// A tuple index must be written in decimal with no separator. The
		// scanner produced an ordinary int_lit under the one restriction to
		// longest-match scanning; the spelling rule is over that token's text.
		if !isTupleIndexSpelling(x.Text) {
			c.errorAt(x.IndexPos, x.End(), diag.TupleIndexSpelling)
		}
		c.expr(x.X)

	case *ast.IndexExpr:
		// Index and TypeArgs are one node, settled by what the operand denotes
		// rather than by shape. Resolution is exactly what tells them apart.
		c.expr(x.X)
		if c.denotesType(x.X) {
			c.typ(x)
			return
		}
		for _, ix := range x.Indices {
			c.expr(ix)
		}

	case *ast.CallExpr:
		c.callExpr(x)

	case *ast.LaunchExpr:
		if x.Config != nil {
			c.expr(x.Config.Blocks)
			c.expr(x.Config.Threads)
		}
		c.expr(x.Call)

	case *ast.AwaitExpr:
		// It parses unconditionally; whether the enclosing body licenses it is
		// this phase's question.
		if !c.ctx.async {
			c.errorAt(x.Await, x.Await+token.Pos(len(token.AWAIT.Spelling())),
				diag.AwaitOutsideAsync)
		}
		c.expr(x.X)

	case *ast.UnaryExpr:
		// `&` is address-of on a value and dereference on a typed_ptr, and `~`
		// is bitwise-NOT here and underlying-type in a TypeSetTerm. Both are
		// read off the operand's type, so neither is settled in this pass.
		c.expr(x.X)

	case *ast.BinaryExpr:
		c.expr(x.X)
		c.expr(x.Y)

	case *ast.CastExpr:
		c.expr(x.X)
		c.typ(x.Type)

	case *ast.TransferExpr:
		// Reached outside the six owning positions.
		c.errorExpr(x, diag.TransferOutsideOwning)
		c.expr(x.Target)

	case *ast.KeyValueExpr:
		c.expr(x.Key)
		c.expr(x.Value)

	case *ast.TupleExpr:
		for _, el := range x.Elems {
			if kv, ok := el.(*ast.KeyValueExpr); ok {
				// A named tuple element's key labels the element; it is not a
				// binding to resolve.
				c.owning(kv.Value)
				continue
			}
			c.owning(el)
		}

	case *ast.ArrayLit:
		for _, el := range x.Elems {
			c.owning(el)
		}

	case *ast.CompositeLit:
		c.typ(x.Type)
		for _, el := range x.Elems {
			// A FieldValue's key is an identifier naming a field, so it is not
			// resolved as an ordinary use.
			if kv, ok := el.(*ast.KeyValueExpr); ok {
				c.owning(kv.Value)
				continue
			}
			c.expr(el)
		}

	case *ast.MapLit:
		for _, el := range x.Elems {
			// A map literal's keys are arbitrary expressions, unlike a
			// composite literal's field names.
			c.expr(el)
		}

	case *ast.EnumShorthand:
		// Legal only where the enum type is fixed by context, which needs
		// types. The name is left for whatever types expressions.
		for _, a := range x.Args {
			c.owning(a)
		}

	case *ast.ChanConstructor:
		c.typ(x.Elem)
		c.expr(x.Cap)

	case *ast.HeapConstructor:
		c.expr(x.X)

	case *ast.FuncLit:
		c.funcLit(x)

	case *ast.EnumPattern:
		// Reached only if a pattern was walked as an expression, which
		// pattern() prevents. Resolve the payload names as uses rather than
		// silently dropping them.
		for _, b := range x.Binds {
			c.lookup(b)
		}

	case *ast.OwnershipType, *ast.ArrayType, *ast.MapType, *ast.ChanType,
		*ast.PointerType, *ast.FuncType, *ast.TensorType, *ast.VectorType,
		*ast.AbstractType:
		// A type reached in expression position — a type argument, a sizeof
		// operand, the callee of a vector call. The parser produces the same
		// nodes either way.
		c.typ(e)
	}
}

// selectorExpr resolves a selector's base. A base may be a package qualifier, a
// keyword namespace, or an ordinary value; only the first is resolvable without
// types, so the member is left for whatever types expressions.
func (c *Checker) selectorExpr(x *ast.SelectorExpr) {
	if _, ok := x.X.(*ast.NamespaceExpr); ok {
		// The namespace member set is closed and is not declarable,
		// shadowable, or extensible, so there is no scope to look in.
		return
	}
	c.expr(x.X)

	id, ok := x.X.(*ast.Ident)
	if !ok {
		return
	}
	pn, ok := c.objectOf(id).(*types.PkgName)
	if !ok {
		return
	}
	m := pn.Imported().Scope().Lookup(x.Sel.Name)
	if m == nil {
		c.errorExpr(x.Sel, diag.UndeclaredName, x.Sel.Name)
		return
	}
	c.recordUse(x.Sel, m)
	c.info.RecordSelection(x, &types.Selection{Kind: types.PackageMember, Obj: m})
}

// callExpr resolves a call and raises the two rules decidable from its shape.
func (c *Checker) callExpr(x *ast.CallExpr) {
	// `transfer` is reserved and bound to nothing, precisely so that a call
	// spelled either way diagnoses as a misspelled ownership marker rather than
	// as an unknown name. The fix-it names the real syntax.
	if name, at := calleeName(x.Fun); name == "transfer" {
		d := diag.New(diag.TransferNotCallable, at.Pos(), at.End())
		if len(x.Args) == 1 {
			d.WithFixit(x.Pos(), x.End(), "var "+exprString(x.Args[0]),
				"use the 'var' prefix")
		}
		c.report(d)
	} else {
		c.expr(x.Fun)
	}

	// A call takes either named or positional arguments, not both. Both shapes
	// land in one slice, so the mixing is checked rather than parsed against.
	named, positional := 0, 0
	for _, a := range x.Args {
		if _, ok := a.(*ast.KeyValueExpr); ok {
			named++
		} else {
			positional++
		}
	}
	if named > 0 && positional > 0 {
		c.errorExpr(x, diag.MixedArguments)
	}

	// sizeof, alignof, and reinterpret take a Type in argument position. The
	// parser already parsed the first argument as one by recognizing the name,
	// which is sound only because a reserved builtin may not be shadowed — so
	// this resolves it as a type rather than as an expression.
	typeFirst := false
	if id, ok := x.Fun.(*ast.Ident); ok {
		typeFirst = token.IsTypeOperator(id.Name)
	}

	for i, a := range x.Args {
		if typeFirst && i == 0 {
			c.typ(a)
			continue
		}
		if kv, ok := a.(*ast.KeyValueExpr); ok {
			// A named argument's key names a parameter, not a binding.
			c.owning(kv.Value)
			continue
		}
		c.owning(a)
	}
}

// calleeName returns the spelling a call's callee ends in, whether written as a
// bare name or as a member, and the node to point at.
func calleeName(fun ast.Expr) (string, ast.Node) {
	switch x := stripParens(fun).(type) {
	case *ast.Ident:
		return x.Name, x
	case *ast.SelectorExpr:
		return x.Sel.Name, x.Sel
	}
	return "", fun
}

// funcLit walks a function literal's body.
//
// A literal begins with all enclosing parse context cleared and re-establishes
// it from its own marker: a closure written inside an async body may not await
// unless it is itself marked async. Clearing the context here is that rule.
func (c *Checker) funcLit(x *ast.FuncLit) {
	saved := c.ctx
	c.ctx = bodyCtx{}
	if x.Type != nil {
		for _, m := range x.Type.Markers {
			switch markerFor(m) {
			case types.MarkerNPU:
				c.ctx.npu = true
			case types.MarkerAsync:
				c.ctx.async = true
			}
		}
	}
	defer func() { c.ctx = saved }()

	c.openFuncScope(x, "function literal")
	defer c.closeScope()

	if x.Type != nil {
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
	}
	c.blockBody(x.Body)
}

// denotesType reports whether an expression names a type rather than a value.
// It is the Index/TypeArgs fork, and it is answerable from resolution alone.
func (c *Checker) denotesType(e ast.Expr) bool {
	switch x := stripParens(e).(type) {
	case *ast.Ident:
		tn, ok := c.objectOf(x).(*types.TypeName)
		return ok && !tn.IsConstraint()
	case *ast.SelectorExpr:
		tn, ok := c.objectOf(x.Sel).(*types.TypeName)
		return ok && !tn.IsConstraint()
	}
	return false
}

func isTupleIndexSpelling(s string) bool {
	if s == "" || strings.ContainsRune(s, '_') {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}