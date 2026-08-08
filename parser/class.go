package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
)

// Classes, structs, members, and decorators — sections E, E.1, D.3, F.

// modifierScope bounds which ClassElementModifiers are accepted at a site.
// The set is the same words everywhere (E), so the difference is only which
// ones get a diagnostic.
type modifierScope uint8

const (
	memberModifiers modifierScope = iota
	paramModifiers
)

// parseModifiers collects a ClassElementModifier sequence (E).
//
// The grammar derives *any* sequence, including duplicates and orders the
// language rejects, so this loop accepts anything and records positions.
// Pointing at the offending word needs per-modifier positions; the bitset
// alone can't do it (§5.2).
func (p *parser) parseModifiers(scope modifierScope) ast.Modifiers {
	var m ast.Modifiers
	for p.at(token.IDENT) {
		bit, ok := modifierBit(p.cur().Ctx)
		if !ok {
			break
		}
		// A modifier word is a modifier only if a member still follows.
		// Otherwise `readonly` is the member's own name, which the grammar
		// permits and which real code uses.
		if !p.modifierFollows() {
			break
		}
		t := p.next()
		if m.Set&bit != 0 {
			p.errorf(t, "duplicate `%s` modifier", t.Ctx)
		}
		if scope == paramModifiers && bit&(ast.ModStatic|ast.ModAbstract|ast.ModDeclare) != 0 {
			p.errorf(t, "`%s` is not allowed on a parameter", t.Ctx)
		}
		m.Set |= bit
		m.List = append(m.List, ast.ModifierTok{Bit: bit, Ctx: t.Ctx, Pos: t.Pos, End: t.End})
	}
	return m
}

func modifierBit(c token.Ctx) (ast.ModifierSet, bool) {
	switch c {
	case token.CtxDeclare:
		return ast.ModDeclare, true
	case token.CtxStatic:
		return ast.ModStatic, true
	case token.CtxAbstract:
		return ast.ModAbstract, true
	case token.CtxOverride:
		return ast.ModOverride, true
	case token.CtxReadonly:
		return ast.ModReadonly, true
	case token.CtxPublic:
		return ast.ModPublic, true
	case token.CtxProtected:
		return ast.ModProtected, true
	case token.CtxPrivate:
		return ast.ModPrivate, true
	}
	return 0, false
}

// modifierFollows reports whether the token after a candidate modifier
// continues a member declaration rather than terminating one.
func (p *parser) modifierFollows() bool {
	switch t := p.peek(1); t.Kind {
	case token.IDENT, token.PRIVATE_IDENT, token.STRING, token.NUMBER, token.BIGINT,
		token.LBRACK, token.MUL, token.AT:
		return true
	case token.SEMI, token.ASSIGN, token.LPAREN, token.LT, token.COLON,
		token.QUESTION, token.NOT, token.RBRACE, token.COMMA, token.RPAREN:
		// `readonly;`, `readonly = 1`, `readonly()`, `readonly: T` — the word
		// is the member.
		return false
	default:
		return t.Kind.IsReserved()
	}
}

// --- decorators (F) ---------------------------------------------------------

func (p *parser) parseDecorators() []*ast.Decorator {
	var out []*ast.Decorator
	for p.at(token.AT) {
		at := p.next().Pos
		out = append(out, &ast.Decorator{At: at, X: p.parseDecoratorExpr()})
	}
	return out
}

// parseDecoratorExpr is the three Decorator forms (F). Deliberately not
// parseExpr: F admits only a member expression, a parenthesized expression, or
// a call, so `@a + b` must not parse.
func (p *parser) parseDecoratorExpr() ast.Expr {
	if p.at(token.LPAREN) {
		lp := p.next().Pos
		x := p.parseExpr(allowIn)
		return &ast.ParenExpr{Lparen: lp, X: x, Rparen: p.expect(token.RPAREN)}
	}

	var x ast.Expr
	if p.at(token.PRIVATE_IDENT) {
		t := p.next()
		x = &ast.PrivateIdent{HashPos: t.Pos, NameEnd: t.End}
	} else {
		x = p.parseIdent()
	}
	for p.at(token.PERIOD) {
		dot := p.next().Pos
		x = &ast.MemberExpr{X: x, Dot: dot, Sel: p.parseIdentName()}
	}

	var typeArgs *ast.TypeArgList
	if p.at(token.LT) {
		p.speculate(func() bool {
			a := p.parseTypeArgs()
			if a == nil || a.Rangle == token.NoPos || !p.at(token.LPAREN) {
				return false
			}
			typeArgs = a
			return true
		})
	}
	if p.at(token.LPAREN) {
		lp := p.next().Pos
		call := &ast.CallExpr{Fun: x, TypeArgs: typeArgs, Lparen: lp}
		call.Args = p.parseArgs()
		call.Rparen = p.expect(token.RPAREN)
		return call
	}
	return x
}

// --- classes (E) ------------------------------------------------------------

func (p *parser) parseClass(decorators []*ast.Decorator, abstractPos token.Pos) *ast.ClassDecl {
	defer p.trace("ClassDeclaration")()

	d := &ast.ClassDecl{Decorators: decorators, AbstractPos: abstractPos}
	d.ClassPos = p.expect(token.CLASS)
	if p.at(token.IDENT) {
		d.Name = p.parseIdent()
	}
	if p.at(token.LT) {
		d.TypeParams = p.parseTypeParams()
	}
	if p.at(token.EXTENDS) {
		d.Extends = p.parseHeritage(token.EXTENDS)
	}
	if p.atCtx(token.CtxImplements) {
		d.Implements = p.parseHeritage(token.IDENT)
	}
	d.Lbrace = p.expect(token.LBRACE)
	d.Members = p.parseMembers(false)
	d.Rbrace = p.expect(token.RBRACE)
	return d
}

// parseHeritage is ClassExtendsClause, ImplementsClause, and
// InterfaceExtendsClause (E, H).
//
// ClassExtendsClause takes a LeftHandSideExpression with optional
// TypeArguments, while the implements and interface-extends clauses take a
// TypeReference list. Both land in a HeritageClause holding TypeExprs; the
// extends form wraps its expression in a TypeRef when it is a plain name and
// diagnoses anything else, so that `class C extends (mixin(B))` is reported
// rather than silently reshaped.
func (p *parser) parseHeritage(kw token.Kind) *ast.HeritageClause {
	t := p.next()
	h := &ast.HeritageClause{KeywordPos: t.Pos, Keyword: kw}
	for {
		before := p.i
		h.Types = append(h.Types, p.parseHeritageType())
		if !p.got(token.COMMA) {
			break
		}
		if !p.advanced(before) {
			break
		}
	}
	if len(h.Types) == 0 {
		h.Types = append(h.Types, p.badType())
	}
	return h
}

func (p *parser) parseHeritageType() ast.TypeExpr {
	if !p.at(token.IDENT) {
		p.errorf(p.cur(), "expected a type name, found %s", p.describe(p.cur()))
		return p.badType()
	}
	ref := &ast.TypeRef{Name: p.parseQualifiedName()}
	if p.at(token.LT) && !p.nlBefore() {
		ref.Args = p.parseTypeArgs()
	}
	return ref
}

// parseMembers is ClassBody (E) and StructBody (E.1).
//
// Source order is preserved and never sorted, canonicalized, or deduped
// (§5.4) — for a struct, source order is layout order (§5.3).
func (p *parser) parseMembers(isStruct bool) []ast.Decl {
	var out []ast.Decl
	for !p.at(token.RBRACE) && !p.atEOF() {
		before := p.i
		m := p.parseMember(isStruct)
		if m != nil {
			out = append(out, m)
		}
		if !p.advanced(before) {
			if !p.advanceTo(syncMember) {
				break
			}
			// A Bad node holds the slot so member offsets do not silently
			// shift (§5.4).
			out = append(out, &ast.BadDecl{From: p.toks[before].Pos, To: p.pos()})
		}
	}
	return out
}

func (p *parser) parseMember(isStruct bool) ast.Decl {
	if p.at(token.SEMI) {
		// An empty ClassElement (E). Not a node — a stray `;` is inert
		// punctuation (§5.4).
		p.next()
		return nil
	}

	var decorators []*ast.Decorator
	if p.at(token.AT) {
		decorators = p.parseDecorators()
	}

	// ClassStaticBlock: `static { ... }` (E). Must be checked before the
	// modifier loop, which would otherwise eat the `static`.
	if p.atCtx(token.CtxStatic) && p.peek(1).Kind == token.LBRACE {
		staticPos := p.next().Pos
		return &ast.StaticBlockDecl{StaticPos: staticPos, Body: p.parseBlock()}
	}

	mods := p.parseModifiers(memberModifiers)

	// ConstructorDeclaration (E).
	if p.atCtx(token.CtxConstructor) && p.peek(1).Kind == token.LPAREN {
		c := &ast.CtorDecl{Decorators: decorators, Mods: mods, CtorPos: p.next().Pos}
		c.Params = p.parseParamList()
		if p.at(token.LBRACE) {
			c.Body = p.parseBlock()
		} else {
			c.Semi = p.expectSemi()
		}
		return c
	}

	// IndexSignature as a ClassElement (E). Reuses the type-member node.
	if p.at(token.LBRACK) && p.peek(1).Kind == token.IDENT && p.peek(2).Kind == token.COLON {
		sig := p.parseTypeMember()
		return &ast.FieldDecl{Decorators: decorators, Mods: mods, Name: sig.(ast.Node),
			Semi: token.NoPos}
	}

	asyncPos, starPos := token.NoPos, token.NoPos
	if p.atCtx(token.CtxAsync) && p.modifierFollows() && !p.peek(1).NLBefore() {
		// `async [no LineTerminator here] ClassElementName` (D.3, L).
		asyncPos = p.next().Pos
	}
	if p.at(token.MUL) {
		starPos = p.next().Pos
	}

	// Accessors (D.3). `get` / `set` are modifiers only when a name follows.
	if (p.atCtx(token.CtxGet) || p.atCtx(token.CtxSet)) && p.modifierFollows() {
		t := p.next()
		m := &ast.MethodDecl{Decorators: decorators, Mods: mods,
			Accessor: t.Ctx, AccessorPos: t.Pos, Name: p.parseClassElementName()}
		m.Params = p.parseParamList()
		if p.at(token.COLON) {
			p.next()
			m.Result = p.parseTypeOrPredicate()
		}
		if p.at(token.LBRACE) {
			m.Body = p.parseBlock()
		}
		return m
	}

	// FieldDefinition's `accessor` form (E).
	accessorPos := token.NoPos
	if p.atCtx(token.CtxAccessor) && p.modifierFollows() && !p.peek(1).NLBefore() {
		// `accessor [no LineTerminator here] ClassElementName` (L).
		accessorPos = p.next().Pos
	}

	name := p.parseClassElementName()

	optional := token.NoPos
	if p.at(token.QUESTION) {
		optional = p.next().Pos
	}

	// A method if a parameter list or type parameters follow.
	if p.at(token.LPAREN) || p.at(token.LT) {
		m := &ast.MethodDecl{Decorators: decorators, Mods: mods, AsyncPos: asyncPos,
			StarPos: starPos, Name: name, Optional: optional}
		if p.at(token.LT) {
			m.TypeParams = p.parseTypeParams()
		}
		m.Params = p.parseParamList()
		if p.at(token.COLON) {
			p.next()
			m.Result = p.parseTypeOrPredicate()
		}
		if p.at(token.LBRACE) {
			m.Body = p.parseBlock()
		} else {
			// MethodSignature as a ClassElement (E), terminated by ASI.
			p.expectSemi()
		}
		return m
	}

	f := &ast.FieldDecl{Decorators: decorators, Mods: mods, AccessorPos: accessorPos,
		Name: name, Optional: optional}
	if p.at(token.NOT) {
		f.Definite = p.next().Pos
	}
	f.Type = p.parseTypeAnnotation()
	if p.at(token.ASSIGN) {
		p.next()
		f.Init = p.parseAssign(allowIn)
	}
	// FieldDefinition is in TerminatedByASI (L).
	f.Semi = p.expectSemi()
	return f
}

// parseClassElementName is ClassElementName (E): a PropertyName or a
// PrivateIdentifier.
func (p *parser) parseClassElementName() ast.Node {
	if p.at(token.PRIVATE_IDENT) {
		t := p.next()
		return &ast.PrivateIdent{HashPos: t.Pos, NameEnd: t.End}
	}
	return p.parsePropertyName()
}

// tryObjectMethod parses a MethodDefinition inside an object literal (B.2 →
// D.3). Returns nil with the cursor unmoved when the property is not a method.
func (p *parser) tryObjectMethod() ast.Expr {
	m := &ast.MethodDecl{}
	start := p.mark()

	if p.atCtx(token.CtxAsync) && p.modifierFollows() && !p.peek(1).NLBefore() {
		m.AsyncPos = p.next().Pos
	}
	if p.at(token.MUL) {
		m.StarPos = p.next().Pos
	}
	if (p.atCtx(token.CtxGet) || p.atCtx(token.CtxSet)) && p.modifierFollows() {
		t := p.next()
		m.Accessor, m.AccessorPos = t.Ctx, t.Pos
	}

	if m.AsyncPos == token.NoPos && m.StarPos == token.NoPos && m.Accessor == token.CtxNone {
		// No leading marker: only a method if `(` or `<` follows the name.
		if !p.methodNameThenParen() {
			p.reset(start)
			return nil
		}
	}

	m.Name = p.parsePropertyName()
	if p.at(token.LT) {
		m.TypeParams = p.parseTypeParams()
	}
	if !p.at(token.LPAREN) {
		p.reset(start)
		return nil
	}
	m.Params = p.parseParamList()
	if p.at(token.COLON) {
		p.next()
		m.Result = p.parseTypeOrPredicate()
	}
	m.Body = p.parseBlock()
	return m
}

func (p *parser) methodNameThenParen() bool {
	j := 1
	if p.at(token.LBRACK) {
		depth := 0
		for {
			switch p.peek(j - 1).Kind {
			case token.LBRACK:
				depth++
			case token.RBRACK:
				depth--
				if depth == 0 {
					goto check
				}
			case token.EOF:
				return false
			}
			j++
		}
	}
check:
	k := p.peek(j).Kind
	return k == token.LPAREN || k == token.LT
}

// --- structs (E.1) ----------------------------------------------------------

func (p *parser) parseStruct(decorators []*ast.Decorator) *ast.StructDecl {
	defer p.trace("StructDeclaration")()

	d := &ast.StructDecl{Decorators: decorators}
	// The contextual `struct` identifier (§5.3).
	d.StructPos = p.next().Pos
	p.noLineTerminator("the struct name")
	d.Name = p.parseIdent()
	if p.at(token.LT) {
		d.TypeParams = p.parseTypeParams()
	}

	// E.1 has no ClassExtendsClause, but §6.3 requires `struct S extends B {}`
	// to parse into a real StructDecl carrying its heritage clause so the
	// message can name the construct instead of saying "unexpected token
	// `extends`". A valid program leaves Extends nil.
	if p.at(token.EXTENDS) {
		d.Extends = p.parseHeritage(token.EXTENDS)
		p.errorAt(d.Extends.Pos(), d.Extends.End(),
			"a struct cannot extend another type; use `implements` to declare conformance")
	}
	if p.atCtx(token.CtxImplements) {
		d.Implements = p.parseHeritage(token.IDENT)
	}

	if p.at(token.LBRACE) {
		body := &ast.StructBody{Lbrace: p.next().Pos}
		body.Members = p.parseMembers(true)
		body.Rbrace = p.expect(token.RBRACE)
		d.Body = body
		return d
	}
	// AmbientStructDeclaration: `struct S ;` with no body (I).
	d.Semi = p.expectSemi()
	return d
}