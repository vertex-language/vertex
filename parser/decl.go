package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
)

// Declarations — sections C.1, D, H, I.

// parseDeclStmt handles the declarations introduced by a reserved word.
func (p *parser) parseDeclStmt() ast.Stmt {
	var decorators []*ast.Decorator
	if p.at(token.AT) {
		decorators = p.parseDecorators()
	}

	switch p.kind() {
	case token.FUNCTION:
		return p.parseFunction(decorators, token.NoPos, ast.AccelNone, false)
	case token.CLASS:
		return p.parseClass(decorators, token.NoPos)
	case token.CONST:
		// `const enum` (H) versus a const binding (C.1).
		if p.peek(1).Kind == token.ENUM {
			return p.parseEnum(p.next().Pos)
		}
		return p.parseVarDecl()
	case token.ENUM:
		return p.parseEnum(token.NoPos)
	case token.IMPORT:
		return p.parseImportDecl()
	case token.EXPORT:
		return p.parseExportDecl()
	}

	if len(decorators) > 0 {
		p.errorAt(decorators[0].Pos(), decorators[len(decorators)-1].End(),
			"a decorator must be followed by a class, function, or member declaration")
	}
	return p.parseExprStmt()
}

// tryContextualDecl handles the declarations introduced by a contextual
// keyword, which the grammar does not disambiguate for us.
//
// ExpressionStatement's lookahead restriction (C) names only `{`, `function`,
// `async function`, `class`, and `let [`. It says nothing about `struct`,
// `interface`, `type`, `enum`, `namespace`, `declare`, `module`, `global`,
// `abstract`, `kernel`, `graph`, `using`, or `accessor` — every one of which is
// a plain identifier that could begin an expression. So the lookahead tests
// below are the parser's own policy, not the grammar's. See the review note.
//
// Returns nil with the cursor unmoved when the word is just an identifier.
func (p *parser) tryContextualDecl() ast.Stmt {
	switch p.cur().Ctx {
	case token.CtxStruct:
		// `struct [no LineTerminator here] BindingIdentifier` (E.1, L). A line
		// break forces the identifier reading, which is what keeps `struct`
		// usable as a name (§4.3).
		if p.peek(1).Kind == token.IDENT && !p.peek(1).NLBefore() {
			return p.parseStruct(nil)
		}

	case token.CtxKernel, token.CtxGraph:
		// AcceleratedFunctionModifier [no LineTerminator here] function (D.4, L).
		if p.peek(1).Kind == token.FUNCTION && !p.peek(1).NLBefore() {
			accel := ast.AccelKernel
			if p.cur().Ctx == token.CtxGraph {
				accel = ast.AccelGraph
			}
			accelPos := p.next().Pos
			return p.parseFunctionAccel(nil, accelPos, accel)
		}

	case token.CtxAsync:
		if p.peek(1).Kind == token.FUNCTION && !p.peek(1).NLBefore() {
			asyncPos := p.next().Pos
			return p.parseFunction(nil, asyncPos, ast.AccelNone, true)
		}

	case token.CtxInterface:
		if p.peek(1).Kind == token.IDENT {
			return p.parseInterface()
		}

	case token.CtxType:
		// `type X =` and `type X<T> =`. Without the `=` this is an expression.
		if p.peek(1).Kind == token.IDENT &&
			(p.peek(2).Kind == token.ASSIGN || p.peek(2).Kind == token.LT) {
			return p.parseTypeAlias()
		}

	case token.CtxNamespace:
		if p.peek(1).Kind == token.IDENT {
			return p.parseNamespace()
		}

	case token.CtxDeclare:
		if p.atAmbientStart(1) {
			return p.parseAmbient()
		}

	case token.CtxModule:
		if p.peek(1).Kind == token.STRING {
			return p.parseModuleDecl()
		}

	case token.CtxGlobal:
		if p.peek(1).Kind == token.LBRACE {
			return p.parseModuleDecl()
		}

	case token.CtxAbstract:
		if p.peek(1).Kind == token.CLASS {
			abstractPos := p.next().Pos
			return p.parseClass(nil, abstractPos)
		}

	case token.CtxLet:
		switch p.peek(1).Kind {
		case token.IDENT, token.LBRACE:
			return p.parseVarDecl()
		case token.LBRACK:
			// `let [` is excluded from ExpressionStatement, so it must be a
			// declaration here.
			return p.parseVarDecl()
		}

	case token.CtxUsing:
		// `using [no LineTerminator here] [lookahead ≠ await]` (C.1).
		if p.peek(1).Kind == token.IDENT && !p.peek(1).NLBefore() &&
			p.peek(1).Ctx != token.CtxAwait {
			return p.parseVarDecl()
		}
	}
	return nil
}

func (p *parser) atAmbientStart(n int) bool {
	switch t := p.peek(n); t.Kind {
	case token.VAR, token.CONST, token.CLASS, token.FUNCTION, token.ENUM:
		return true
	case token.IDENT:
		switch t.Ctx {
		case token.CtxLet, token.CtxNamespace, token.CtxModule, token.CtxGlobal,
			token.CtxStruct, token.CtxAbstract:
			return true
		}
	}
	return false
}

// --- bindings (C.1, C.2) ----------------------------------------------------

func (p *parser) parseVarDecl() ast.Stmt {
	d := &ast.VarDecl{}
	switch {
	case p.at(token.VAR):
		d.Kind, d.KindPos = ast.VarVar, p.next().Pos
	case p.at(token.CONST):
		d.Kind, d.KindPos = ast.VarConst, p.next().Pos
	case p.at(token.AWAIT):
		// AwaitUsingDeclaration: `await [no LT] using [no LT] ...` (C.1, L).
		d.AwaitPos = p.next().Pos
		p.noLineTerminator("`using`")
		d.Kind, d.KindPos = ast.VarAwaitUsing, p.expectCtx(token.CtxUsing)
		p.noLineTerminator("the binding list")
	case p.atCtx(token.CtxUsing):
		d.Kind, d.KindPos = ast.VarUsing, p.next().Pos
		p.noLineTerminator("the binding list")
	default:
		d.Kind, d.KindPos = ast.VarLet, p.next().Pos
	}

	// UsingBinding requires an initializer; LexicalBinding does not (C.1).
	requireInit := d.Kind == ast.VarUsing || d.Kind == ast.VarAwaitUsing
	for {
		before := p.i
		d.List = append(d.List, p.parseBinding(allowIn, requireInit))
		if !p.got(token.COMMA) {
			break
		}
		if !p.advanced(before) {
			break
		}
	}
	d.Semi = p.expectSemi()
	return d
}

func (p *parser) parseBinding(in inFlag, requireInit bool) *ast.Binding {
	b := &ast.Binding{}
	switch {
	case p.at(token.LBRACE), p.at(token.LBRACK):
		b.Pattern = p.parseBindingTarget()
	default:
		b.Name = p.parseIdent()
		// DefiniteAssignmentAssertion (C.1), only on the identifier form.
		if p.at(token.NOT) {
			b.Definite = p.next().Pos
		}
	}
	b.Type = p.parseTypeAnnotation()

	if p.at(token.ASSIGN) {
		p.next()
		b.Init = p.parseAssign(in)
	} else if b.Pattern != nil || requireInit {
		// `BindingPattern TypeAnnotation_opt Initializer` has no _opt on the
		// initializer, and neither does UsingBinding.
		p.errorAt(b.Pos(), b.End(), "this binding requires an initializer")
	}
	return b
}

// parseBindingTarget is BindingIdentifier | BindingPattern (C.2), also used for
// catch parameters and parameter names.
func (p *parser) parseBindingTarget() ast.Expr {
	switch p.kind() {
	case token.LBRACE:
		return p.parseObjectBindingPattern()
	case token.LBRACK:
		return p.parseArrayBindingPattern()
	default:
		return p.parseIdent()
	}
}

// parseObjectBindingPattern is ObjectBindingPattern (C.2).
//
// Declaration position is known here, so this parses straight into pattern
// nodes with no cover and no reinterpretation — unlike the assignment side
// (B.6), which arrives through objectToPattern.
func (p *parser) parseObjectBindingPattern() ast.Expr {
	pat := &ast.ObjectPattern{Lbrace: p.next().Pos}
	for !p.at(token.RBRACE) && !p.atEOF() {
		before := p.i
		if p.at(token.ELLIPSIS) {
			t := p.next()
			// BindingRestProperty takes only an identifier, never a pattern.
			pat.Props = append(pat.Props, &ast.RestElem{Ellipsis: t.Pos, X: p.parseIdent()})
		} else {
			pat.Props = append(pat.Props, p.parseBindingProperty())
		}
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	pat.Rbrace = p.expect(token.RBRACE)
	return pat
}

func (p *parser) parseBindingProperty() ast.Expr {
	// SingleNameBinding versus `PropertyName : BindingElement` (C.2).
	if (p.at(token.IDENT) || p.kind().IsReserved()) && p.peek(1).Kind != token.COLON {
		id := p.parseIdent()
		if p.at(token.ASSIGN) {
			assign := p.next().Pos
			return &ast.PropertyPattern{Value: &ast.AssignPattern{
				Lhs: id, Assign: assign, Rhs: p.parseAssign(allowIn)}}
		}
		return &ast.PropertyPattern{Value: id}
	}
	key := p.parsePropertyName()
	colon := p.expect(token.COLON)
	return &ast.PropertyPattern{Key: key, Colon: colon, Value: p.parseBindingElement()}
}

func (p *parser) parseArrayBindingPattern() ast.Expr {
	pat := &ast.ArrayPattern{Lbrack: p.next().Pos}
	for !p.at(token.RBRACK) && !p.atEOF() {
		before := p.i
		switch {
		case p.at(token.COMMA):
			pat.Elts = append(pat.Elts, &ast.Elision{Comma: p.next().Pos})
			continue
		case p.at(token.ELLIPSIS):
			t := p.next()
			// BindingRestElement takes an identifier or a pattern (C.2).
			pat.Elts = append(pat.Elts, &ast.RestElem{Ellipsis: t.Pos, X: p.parseBindingTarget()})
		default:
			pat.Elts = append(pat.Elts, p.parseBindingElement())
		}
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	pat.Rbrack = p.expect(token.RBRACK)
	return pat
}

func (p *parser) parseBindingElement() ast.Expr {
	target := p.parseBindingTarget()
	if p.at(token.ASSIGN) {
		assign := p.next().Pos
		return &ast.AssignPattern{Lhs: target, Assign: assign, Rhs: p.parseAssign(allowIn)}
	}
	return target
}

// --- functions (D, D.1, D.4) ------------------------------------------------

func (p *parser) parseFunction(decorators []*ast.Decorator, asyncPos token.Pos, accel ast.AccelKind, async bool) *ast.FuncDecl {
	defer p.trace("FunctionDeclaration")()

	fn := &ast.FuncDecl{Decorators: decorators, Accel: accel, Async: async}
	if asyncPos != token.NoPos {
		fn.AccelPos = token.NoPos
	}
	fn.FuncPos = p.expect(token.FUNCTION)
	if p.at(token.MUL) {
		p.next()
		fn.Gen = true
	}
	if p.at(token.IDENT) {
		fn.Name = p.parseIdent()
	}
	if p.at(token.LT) {
		fn.TypeParams = p.parseTypeParams()
	}
	fn.Params = p.parseParamList()
	if p.at(token.COLON) {
		p.next()
		fn.Result = p.parseTypeOrPredicate()
	}
	if p.at(token.LBRACE) {
		fn.Body = p.parseBlock()
	} else {
		// The signature-only form, `function f(): void;` (D).
		fn.Semi = p.expectSemi()
	}
	return fn
}

// parseFunctionAccel is AcceleratedFunctionDeclaration (D.4).
//
// One FuncDecl for plain, kernel, and graph functions, matching D.4's single
// production. Async and Gen are recorded even though D.4 admits neither, so
// `kernel async function f() {}` is rejected by name rather than by parse
// failure (§5.3) — the check below is exactly that, and it runs here rather
// than in a later phase only because the positions are at hand.
func (p *parser) parseFunctionAccel(decorators []*ast.Decorator, accelPos token.Pos, accel ast.AccelKind) *ast.FuncDecl {
	fn := p.parseFunction(decorators, token.NoPos, accel, false)
	fn.AccelPos = accelPos
	if fn.Async || fn.Gen {
		what := "a generator"
		if fn.Async {
			what = "an async function"
		}
		p.errorAt(accelPos, fn.FuncPos, "an accelerated function cannot be %s", what)
	}
	return fn
}

// parseParamList is FormalParameters (D.1).
func (p *parser) parseParamList() *ast.ParamList {
	pl := &ast.ParamList{Lparen: p.expect(token.LPAREN)}
	for !p.at(token.RPAREN) && !p.atEOF() {
		before := p.i
		switch {
		case p.at(token.THIS) && pl.This == nil && len(pl.List) == 0:
			// ThisParameter is always first (D.1).
			t := p.next()
			pl.This = &ast.ThisParam{ThisPos: t.Pos, Type: p.parseTypeAnnotation()}
		case p.at(token.ELLIPSIS):
			t := p.next()
			r := &ast.RestElem{Ellipsis: t.Pos, X: p.parseBindingTarget()}
			r.Type = p.parseTypeAnnotation()
			pl.Rest = r
		default:
			pl.List = append(pl.List, p.parseParam())
		}
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	if p.at(token.RPAREN) {
		pl.Rparen = p.next().Pos
	} else {
		p.errorf(p.cur(), "expected `)` to close the parameter list, found %s", p.describe(p.cur()))
	}
	return pl
}

func (p *parser) parseParam() *ast.Param {
	prm := &ast.Param{}
	if p.at(token.AT) {
		prm.Decorators = p.parseDecorators()
	}
	prm.Mods = p.parseModifiers(paramModifiers)
	prm.Name = p.parseBindingTarget()
	if p.at(token.QUESTION) {
		prm.Optional = p.next().Pos
	}
	prm.Type = p.parseTypeAnnotation()
	if p.at(token.ASSIGN) {
		p.next()
		prm.Init = p.parseAssign(allowIn)
	}
	return prm
}

// --- interfaces, aliases, enums, namespaces (H) -----------------------------

func (p *parser) parseInterface() ast.Stmt {
	d := &ast.InterfaceDecl{IfacePos: p.next().Pos}
	d.Name = p.parseIdent()
	if p.at(token.LT) {
		d.TypeParams = p.parseTypeParams()
	}
	if p.at(token.EXTENDS) {
		d.Extends = p.parseHeritage(token.EXTENDS)
	}
	d.Body = p.parseObjectType()
	return d
}

func (p *parser) parseTypeAlias() ast.Stmt {
	d := &ast.TypeAliasDecl{TypePos: p.next().Pos}
	d.Name = p.parseIdent()
	if p.at(token.LT) {
		d.TypeParams = p.parseTypeParams()
	}
	d.Assign = p.expect(token.ASSIGN)
	d.Type = p.parseTypeOrPredicate()
	d.Semi = p.expectSemi()
	return d
}

func (p *parser) parseEnum(constPos token.Pos) ast.Stmt {
	d := &ast.EnumDecl{ConstPos: constPos, EnumPos: p.expect(token.ENUM)}
	d.Name = p.parseIdent()
	if p.at(token.COLON) {
		p.next()
		// EnumUnderlyingType is a TypeReference only (H).
		d.Underlying = p.parseType()
	}
	d.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.atEOF() {
		before := p.i
		m := &ast.EnumMember{}
		if p.at(token.STRING) {
			t := p.next()
			m.Name = &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
		} else {
			m.Name = p.parseIdentName()
		}
		if p.at(token.ASSIGN) {
			m.Assign = p.next().Pos
			// Not folded to a constant — §1 forbids folding.
			m.Value = p.parseAssign(allowIn)
		}
		d.Members = append(d.Members, m)
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	d.Rbrace = p.expect(token.RBRACE)
	return d
}

func (p *parser) parseNamespace() ast.Stmt {
	d := &ast.NamespaceDecl{NsPos: p.next().Pos}
	// IdentifierPath: dotted names declare nested namespaces (H).
	var name ast.Node = p.parseIdent()
	for p.at(token.PERIOD) {
		p.next()
		name = &ast.QualifiedName{X: name, Sel: p.parseIdent()}
	}
	d.Name = name
	d.Lbrace = p.expect(token.LBRACE)
	d.Items = p.parseModuleItems(token.RBRACE)
	d.Rbrace = p.expect(token.RBRACE)
	return d
}

func (p *parser) parseModuleDecl() ast.Stmt {
	t := p.next()
	d := &ast.ModuleDecl{KeywordPos: t.Pos, KeywordEnd: t.End, IsGlobal: t.Ctx == token.CtxGlobal}
	if !d.IsGlobal {
		if p.at(token.STRING) {
			s := p.next()
			d.Name = &ast.BasicLit{Kind: s.Kind, ValuePos: s.Pos, ValueEnd: s.End, HasEscape: s.HasEscape()}
		} else {
			p.errorf(p.cur(), "expected a module name string")
		}
	}
	if p.at(token.LBRACE) {
		d.Lbrace = p.next().Pos
		d.Items = p.parseModuleItems(token.RBRACE)
		d.Rbrace = p.expect(token.RBRACE)
		return d
	}
	// `module "x";` with no body (I).
	d.Semi = p.expectSemi()
	return d
}

// --- ambient declarations (I) -----------------------------------------------

func (p *parser) parseAmbient() ast.Stmt {
	d := &ast.AmbientDecl{DeclarePos: p.next().Pos}

	switch {
	case p.at(token.VAR), p.at(token.CONST), p.atCtx(token.CtxLet):
		d.Inner = p.parseVarDecl().(ast.Decl)
	case p.at(token.FUNCTION):
		d.Inner = p.parseFunction(nil, token.NoPos, ast.AccelNone, false)
	case p.at(token.CLASS):
		d.Inner = p.parseClass(nil, token.NoPos)
	case p.atCtx(token.CtxAbstract):
		abstractPos := p.next().Pos
		d.Inner = p.parseClass(nil, abstractPos)
	case p.at(token.ENUM):
		d.Inner = p.parseEnum(token.NoPos).(ast.Decl)
	case p.at(token.CONST) && p.peek(1).Kind == token.ENUM:
		d.Inner = p.parseEnum(p.next().Pos).(ast.Decl)
	case p.atCtx(token.CtxNamespace):
		d.Inner = p.parseNamespace().(ast.Decl)
	case p.atCtx(token.CtxModule), p.atCtx(token.CtxGlobal):
		d.Inner = p.parseModuleDecl().(ast.Decl)
	case p.atCtx(token.CtxStruct):
		// AmbientStructDeclaration: `struct [no LT] BindingIdentifier ;` (I).
		// Body stays nil, which is how §5.3 marks the ambient form.
		d.Inner = p.parseStruct(nil)
	default:
		p.errorf(p.cur(), "expected a declaration after `declare`, found %s", p.describe(p.cur()))
		d.Inner = &ast.BadDecl{From: d.DeclarePos, To: p.end()}
	}
	return d
}