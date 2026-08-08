package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
)

// Types — vertex_grammar.md section G.
//
// A separate sublanguage with its own node hierarchy (§5.1), entered at
// TypeAnnotation, ReturnTypeAnnotation, TypeArguments, heritage clauses, and
// the `< Type >` assertion. The parser always knows which side it is on, so
// nothing here re-derives "am I in a type?".
//
// The closing `>` is one token by construction: the scanner never merges `>`
// (§4.2), so this file does nothing special to split one.

// parseTypeAnnotation is TypeAnnotation (G.1): `: Type`.
func (p *parser) parseTypeAnnotation() ast.TypeExpr {
	if !p.got(token.COLON) {
		return nil
	}
	return p.parseType()
}

// parseTypeOrPredicate is ReturnTypeAnnotation's payload (G.1), which admits
// either a Type or a TypePredicate.
func (p *parser) parseTypeOrPredicate() ast.TypeExpr {
	if t := p.tryTypePredicate(); t != nil {
		return t
	}
	return p.parseType()
}

func (p *parser) tryTypePredicate() ast.TypeExpr {
	var out ast.TypeExpr
	p.speculate(func() bool {
		pred := &ast.TypePredicate{}
		if p.atCtx(token.CtxAsserts) && !p.peek(1).NLBefore() {
			pred.AssertsPos = p.next().Pos
		}
		switch {
		case p.at(token.THIS):
			t := p.next()
			pred.Param = &ast.ThisType{ThisPos: t.Pos, ThisEnd: t.End}
		case p.at(token.IDENT):
			pred.Param = p.parseIdent()
		default:
			return false
		}
		if p.atCtx(token.CtxIs) {
			pred.IsPos = p.next().Pos
			pred.Type = p.parseType()
			out = pred
			return true
		}
		// `asserts x` with no `is` is the bare form. Without `asserts` and
		// without `is`, this is an ordinary TypeReference, not a predicate.
		if pred.AssertsPos != token.NoPos {
			out = pred
			return true
		}
		return false
	})
	return out
}

// parseType is Type (G.1): FunctionType, ConstructorType, or ConditionalType.
func (p *parser) parseType() ast.TypeExpr {
	defer p.trace("Type")()
	if !p.enter() {
		return p.badType()
	}
	defer p.leave()

	if t := p.tryFunctionType(); t != nil {
		return t
	}
	if p.at(token.NEW) || (p.atCtx(token.CtxAbstract) && p.peek(1).Kind == token.NEW) {
		return p.parseCtorType()
	}
	return p.parseConditionalType()
}

func (p *parser) tryFunctionType() ast.TypeExpr {
	if !p.at(token.LPAREN) && !p.at(token.LT) {
		return nil
	}
	var out ast.TypeExpr
	p.speculate(func() bool {
		ft := &ast.FuncType{}
		if p.at(token.LT) {
			ft.TypeParams = p.parseTypeParams()
			if ft.TypeParams == nil || ft.TypeParams.Rangle == token.NoPos {
				return false
			}
		}
		if !p.at(token.LPAREN) {
			return false
		}
		ft.Params = p.parseParamTypeList()
		if ft.Params.Rparen == token.NoPos {
			return false
		}
		if !p.at(token.ARROW) {
			return false
		}
		ft.Arrow = p.next().Pos
		ft.Result = p.parseTypeOrPredicate()
		out = ft
		return true
	})
	return out
}

func (p *parser) parseCtorType() ast.TypeExpr {
	ct := &ast.CtorType{AbstractPos: token.NoPos}
	if p.atCtx(token.CtxAbstract) {
		ct.AbstractPos = p.next().Pos
	}
	ct.NewPos = p.expect(token.NEW)
	if p.at(token.LT) {
		ct.TypeParams = p.parseTypeParams()
	}
	ct.Params = p.parseParamTypeList()
	ct.Arrow = p.expect(token.ARROW)
	ct.Result = p.parseType()
	return ct
}

// parseConditionalType is ConditionalType (G.1).
//
// The `extends` here carries a [no ConditionalType here] restriction, so the
// ExtendsType is parsed as a UnionType and a nested conditional must be
// parenthesized.
func (p *parser) parseConditionalType() ast.TypeExpr {
	check := p.parseUnionType()
	if !p.at(token.EXTENDS) {
		return check
	}
	c := &ast.CondType{Check: check, ExtendsPos: p.next().Pos}
	c.Extends = p.parseUnionType()
	c.Quest = p.expect(token.QUESTION)
	c.Then = p.parseType()
	c.Colon = p.expect(token.COLON)
	c.Else = p.parseType()
	return c
}

func (p *parser) parseUnionType() ast.TypeExpr {
	lead := token.NoPos
	if p.at(token.OR) {
		lead = p.next().Pos
	}
	first := p.parseIntersectionType()
	if !p.at(token.OR) && lead == token.NoPos {
		return first
	}
	u := &ast.UnionType{LeadPos: lead, Types: []ast.TypeExpr{first}}
	for p.got(token.OR) {
		u.Types = append(u.Types, p.parseIntersectionType())
	}
	if len(u.Types) == 1 && lead == token.NoPos {
		return u.Types[0]
	}
	return u
}

func (p *parser) parseIntersectionType() ast.TypeExpr {
	lead := token.NoPos
	if p.at(token.AND) {
		lead = p.next().Pos
	}
	first := p.parseTypeOperator()
	if !p.at(token.AND) && lead == token.NoPos {
		return first
	}
	it := &ast.IntersectionType{LeadPos: lead, Types: []ast.TypeExpr{first}}
	for p.got(token.AND) {
		it.Types = append(it.Types, p.parseTypeOperator())
	}
	if len(it.Types) == 1 && lead == token.NoPos {
		return it.Types[0]
	}
	return it
}

// parseTypeOperator is TypeOperator (G.1): keyof, unique, readonly, mutating,
// and infer.
func (p *parser) parseTypeOperator() ast.TypeExpr {
	if p.at(token.IDENT) {
		switch c := p.cur().Ctx; c {
		case token.CtxKeyof, token.CtxUnique, token.CtxReadonly, token.CtxMutating:
			// `mutating T` is ambiguous between a passing mode and a type
			// operator and gets this node either way; position decides later
			// (§5.1).
			if p.typeFollows(1) {
				t := p.next()
				return &ast.TypeOp{OpPos: t.Pos, Op: c, X: p.parseTypeOperator()}
			}
		case token.CtxInfer:
			if p.peek(1).Kind == token.IDENT {
				inferPos := p.next().Pos
				it := &ast.InferType{InferPos: inferPos, Name: p.parseIdent()}
				if p.at(token.EXTENDS) {
					// InferConstraint also carries [no ConditionalType here].
					it.ExtendsPos = p.next().Pos
					it.Constraint = p.parseUnionType()
				}
				return it
			}
		}
	}
	return p.parsePostfixType()
}

// typeFollows reports whether a type can begin n tokens ahead, so that
// `readonly` as a property name is not read as an operator.
func (p *parser) typeFollows(n int) bool {
	switch t := p.peek(n); t.Kind {
	case token.IDENT, token.LBRACK, token.LBRACE, token.LPAREN, token.LT,
		token.STRING, token.NUMBER, token.BIGINT, token.TEMPLATE, token.TEMPLATE_HEAD,
		token.THIS, token.TYPEOF, token.IMPORT, token.NEW, token.VOID, token.NULL,
		token.TRUE, token.FALSE, token.SUB, token.OR, token.AND:
		return true
	}
	return false
}

// parsePostfixType is PostfixType (G.1): `T[]` and `T[K]`, both under a
// [no LineTerminator here] restriction (L).
func (p *parser) parsePostfixType() ast.TypeExpr {
	t := p.parsePrimaryType()
	for p.at(token.LBRACK) && !p.nlBefore() {
		lb := p.next().Pos
		if p.at(token.RBRACK) {
			t = &ast.ArrayType{Elem: t, Lbrack: lb, Rbrack: p.next().Pos}
			continue
		}
		idx := p.parseType()
		t = &ast.IndexedAccessType{X: t, Lbrack: lb, Index: idx, Rbrack: p.expect(token.RBRACK)}
	}
	return t
}

func (p *parser) parsePrimaryType() ast.TypeExpr {
	t := p.cur()

	switch t.Kind {
	case token.LPAREN:
		lp := p.next().Pos
		inner := p.parseType()
		return &ast.ParenType{Lparen: lp, X: inner, Rparen: p.expect(token.RPAREN)}

	case token.LBRACE:
		if p.atMappedType() {
			return p.parseMappedType()
		}
		return p.parseObjectType()

	case token.LBRACK:
		return p.parseTupleType()

	case token.THIS:
		p.next()
		return &ast.ThisType{ThisPos: t.Pos, ThisEnd: t.End}

	case token.TYPEOF:
		p.next()
		q := &ast.TypeQuery{TypeofPos: t.Pos}
		if p.at(token.IMPORT) {
			q.Name = p.parseImportType()
		} else {
			q.Name = p.parseQualifiedName()
			if p.at(token.LT) {
				q.Args = p.parseTypeArgs()
			}
		}
		return q

	case token.IMPORT:
		return p.parseImportType()

	case token.STRING, token.NUMBER, token.BIGINT:
		p.next()
		return &ast.LiteralType{Value: &ast.BasicLit{
			Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}}

	case token.SUB:
		// `- NumericLiteral` and `- BigIntLiteral` (G.1).
		minus := p.next().Pos
		lit := p.cur()
		if lit.Kind != token.NUMBER && lit.Kind != token.BIGINT {
			p.errorf(lit, "expected a numeric literal after `-` in a type")
			return &ast.BadType{From: minus, To: lit.End}
		}
		p.next()
		return &ast.LiteralType{MinusPos: minus, Value: &ast.BasicLit{
			Kind: lit.Kind, ValuePos: lit.Pos, ValueEnd: lit.End}}

	case token.TRUE, token.FALSE, token.NULL:
		p.next()
		return &ast.LiteralType{Value: &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End}}

	case token.VOID:
		// PredefinedType includes void and null, which are ReservedWords and
		// so carry a Kind rather than a Ctx (G.1, §3).
		p.next()
		return &ast.PredefinedType{NamePos: t.Pos, NameEnd: t.End, Kind: t.Kind}

	case token.TEMPLATE, token.TEMPLATE_HEAD:
		return p.parseTemplateLiteralType()

	case token.IDENT:
		if t.Ctx.IsPredefinedType() {
			p.next()
			return &ast.PredefinedType{NamePos: t.Pos, NameEnd: t.End, Ctx: t.Ctx}
		}
		ref := &ast.TypeRef{Name: p.parseQualifiedName()}
		if p.at(token.LT) && !p.nlBefore() {
			// TypeName [no LineTerminator here] TypeArguments (G.1, L).
			ref.Args = p.parseTypeArgs()
		}
		return ref
	}

	p.errorf(t, "expected a type, found %s", p.describe(t))
	return p.badType()
}

func (p *parser) parseQualifiedName() ast.Node {
	var n ast.Node = p.parseIdent()
	for p.at(token.PERIOD) {
		p.next()
		n = &ast.QualifiedName{X: n, Sel: p.parseIdentName()}
	}
	return n
}

func (p *parser) parseImportType() ast.TypeExpr {
	it := &ast.ImportType{ImportPos: p.next().Pos}
	it.Lparen = p.expect(token.LPAREN)
	if p.at(token.STRING) {
		t := p.next()
		it.Path = &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
	} else {
		p.errorf(p.cur(), "expected a module specifier string")
	}
	if p.at(token.COMMA) {
		// ImportTypeAttributes: `, { with : { ... } }` (G.1). The outer braces
		// belong to this production, not to the WithClause.
		p.next()
		p.expect(token.LBRACE)
		it.With = p.parseWithClause()
		p.expect(token.RBRACE)
	}
	it.Rparen = p.expect(token.RPAREN)
	if p.got(token.PERIOD) {
		it.Qualifier = p.parseQualifiedName()
	}
	if p.at(token.LT) {
		it.Args = p.parseTypeArgs()
	}
	return it
}

func (p *parser) parseTemplateLiteralType() ast.TypeExpr {
	// Type-position templates reuse the same lexical tokens as A.6 — the
	// grammar has no TemplateTypeHead precisely so there is one way to
	// tokenize `` `hello ${ `` (§4.5) — so only this parse differs.
	lit := &ast.TemplateLiteralType{}
	t := p.next()
	lit.Quasis = append(lit.Quasis, &ast.TemplateElem{Kind: t.Kind, Start: t.Pos, Stop: t.End})
	if t.Kind == token.TEMPLATE {
		return lit
	}
	for {
		lit.Types = append(lit.Types, p.parseType())
		q := p.cur()
		if q.Kind != token.TEMPLATE_MIDDLE && q.Kind != token.TEMPLATE_TAIL {
			p.errorf(q, "expected a template continuation in a type")
			lit.Quasis = append(lit.Quasis, &ast.TemplateElem{
				Kind: token.TEMPLATE_TAIL, Start: q.Pos, Stop: q.Pos + 1})
			return lit
		}
		p.next()
		lit.Quasis = append(lit.Quasis, &ast.TemplateElem{Kind: q.Kind, Start: q.Pos, Stop: q.End})
		if q.Kind == token.TEMPLATE_TAIL {
			return lit
		}
	}
}

func (p *parser) parseTupleType() ast.TypeExpr {
	tt := &ast.TupleType{Lbrack: p.next().Pos}
	for !p.at(token.RBRACK) && !p.atEOF() {
		before := p.i
		tt.Elems = append(tt.Elems, p.parseTupleElem())
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	tt.Rbrack = p.expect(token.RBRACK)
	return tt
}

// parseTupleElem is TupleElementType (G.5), covering the named, optional, and
// rest forms in one node.
func (p *parser) parseTupleElem() ast.TypeExpr {
	e := &ast.TupleElem{}
	if p.at(token.ELLIPSIS) {
		e.Ellipsis = p.next().Pos
	}
	// A named element is `name : Type` or `name ? : Type`. Distinguishing it
	// from a bare TypeReference needs two tokens of lookahead.
	if p.at(token.IDENT) {
		if p.peek(1).Kind == token.COLON {
			e.Name = p.parseIdent()
			e.NameColon = p.next().Pos
		} else if p.peek(1).Kind == token.QUESTION && p.peek(2).Kind == token.COLON {
			e.Name = p.parseIdent()
			e.OptionalPos = p.next().Pos
			e.NameColon = p.next().Pos
		}
	}
	e.Type = p.parseType()
	if e.Name == nil && p.at(token.QUESTION) {
		e.OptionalPos = p.next().Pos
	}
	return e
}

// atMappedType distinguishes MappedType (G.6) from ObjectType with an
// IndexSignature (G.3). The difference is `in` after the parameter name.
func (p *parser) atMappedType() bool {
	j := 1
	for {
		t := p.peek(j)
		if t.Kind == token.IDENT && (t.Ctx == token.CtxReadonly) {
			j++
			continue
		}
		if t.Kind == token.ADD || t.Kind == token.SUB {
			// `+readonly` / `-readonly`
			if p.peek(j+1).Kind == token.IDENT && p.peek(j+1).Ctx == token.CtxReadonly {
				j += 2
				continue
			}
		}
		break
	}
	if p.peek(j).Kind != token.LBRACK || p.peek(j+1).Kind != token.IDENT {
		return false
	}
	return p.peek(j + 2).Kind == token.IN
}

func (p *parser) parseMappedType() ast.TypeExpr {
	m := &ast.MappedType{Lbrace: p.next().Pos, ReadonlySign: token.INVALID, OptionalSign: token.INVALID}

	if p.at(token.ADD) || p.at(token.SUB) {
		m.ReadonlySign = p.kind()
		m.ReadonlyPos = p.next().Pos
		p.expectCtx(token.CtxReadonly)
	} else if p.atCtx(token.CtxReadonly) {
		m.ReadonlyPos = p.next().Pos
	}

	m.Lbrack = p.expect(token.LBRACK)
	m.Name = p.parseIdent()
	m.InPos = p.expect(token.IN)
	m.Constraint = p.parseType()
	if p.atCtx(token.CtxAs) {
		m.AsPos = p.next().Pos
		m.As = p.parseType()
	}
	m.Rbrack = p.expect(token.RBRACK)

	if p.at(token.ADD) || p.at(token.SUB) {
		m.OptionalSign = p.kind()
		m.OptionalPos = p.next().Pos
		p.expect(token.QUESTION)
	} else if p.at(token.QUESTION) {
		m.OptionalPos = p.next().Pos
	}

	m.Type = p.parseTypeAnnotation()
	// TypeMemberSeparator_opt closes the body (G.6).
	if p.at(token.SEMI) || p.at(token.COMMA) {
		p.next()
	}
	m.Rbrace = p.expect(token.RBRACE)
	return m
}

func (p *parser) parseObjectType() *ast.ObjectType {
	o := &ast.ObjectType{Lbrace: p.expect(token.LBRACE)}
	for !p.at(token.RBRACE) && !p.atEOF() {
		before := p.i
		m := p.parseTypeMember()
		if m != nil {
			o.Members = append(o.Members, m)
		}
		if !p.advanced(before) {
			if !p.advanceTo(syncTypeMember) {
				break
			}
		}
	}
	o.Rbrace = p.expect(token.RBRACE)
	return o
}

// parseTypeMember is TypeMember (G.3): the seven signature forms.
//
// Every one of them is in TerminatedByASI (grammar L), so each ends with
// expectMemberSep rather than a required punctuator.
func (p *parser) parseTypeMember() ast.TypeExpr {
	switch {
	case p.at(token.LPAREN) || p.at(token.LT):
		sig := &ast.CallSig{}
		if p.at(token.LT) {
			sig.TypeParams = p.parseTypeParams()
		}
		sig.Params = p.parseParamTypeList()
		if p.at(token.COLON) {
			p.next()
			sig.Result = p.parseTypeOrPredicate()
		}
		sig.Sep = p.expectMemberSep()
		return sig

	case p.at(token.NEW) || (p.atCtx(token.CtxAbstract) && p.peek(1).Kind == token.NEW):
		sig := &ast.ConstructSig{}
		if p.atCtx(token.CtxAbstract) {
			sig.AbstractPos = p.next().Pos
		}
		sig.NewPos = p.next().Pos
		if p.at(token.LT) {
			sig.TypeParams = p.parseTypeParams()
		}
		sig.Params = p.parseParamTypeList()
		if p.at(token.COLON) {
			p.next()
			sig.Result = p.parseTypeOrPredicate()
		}
		sig.Sep = p.expectMemberSep()
		return sig
	}

	readonly := token.NoPos
	if p.atCtx(token.CtxReadonly) && !p.memberNameFollows(1) {
		// `readonly` is a modifier only if something that isn't a member
		// terminator follows it.
	}
	if p.atCtx(token.CtxReadonly) && p.memberNameFollows(1) {
		readonly = p.next().Pos
	}

	if p.at(token.LBRACK) && p.peek(1).Kind == token.IDENT && p.peek(2).Kind == token.COLON {
		sig := &ast.IndexSig{ReadonlyPos: readonly, Lbrack: p.next().Pos}
		sig.Name = p.parseIdent()
		p.expect(token.COLON)
		sig.Key = p.parseUnionType()
		sig.Rbrack = p.expect(token.RBRACK)
		sig.Type = p.parseTypeAnnotation()
		sig.Sep = p.expectMemberSep()
		return sig
	}

	if (p.atCtx(token.CtxGet) || p.atCtx(token.CtxSet)) && p.memberNameFollows(1) {
		t := p.next()
		sig := &ast.AccessorSig{KwPos: t.Pos, Kw: t.Ctx, Name: p.parsePropertyName()}
		sig.Lparen = p.expect(token.LPAREN)
		if t.Ctx == token.CtxSet && !p.at(token.RPAREN) {
			sig.Param = p.parseParamType()
		}
		sig.Rparen = p.expect(token.RPAREN)
		if p.at(token.COLON) {
			p.next()
			sig.Result = p.parseTypeOrPredicate()
		}
		sig.Sep = p.expectMemberSep()
		return sig
	}

	name := p.parsePropertyName()
	optional := token.NoPos
	if p.at(token.QUESTION) {
		optional = p.next().Pos
	}

	if p.at(token.LPAREN) || p.at(token.LT) {
		sig := &ast.MethodSig{Name: name, OptionalPos: optional}
		if p.at(token.LT) {
			sig.TypeParams = p.parseTypeParams()
		}
		sig.Params = p.parseParamTypeList()
		if p.at(token.COLON) {
			p.next()
			sig.Result = p.parseTypeOrPredicate()
		}
		sig.Sep = p.expectMemberSep()
		return sig
	}

	sig := &ast.PropertySig{ReadonlyPos: readonly, Name: name, OptionalPos: optional}
	sig.Type = p.parseTypeAnnotation()
	sig.Sep = p.expectMemberSep()
	return sig
}

// memberNameFollows reports whether the token n ahead can start a member name,
// which is how a contextual modifier is told from a property called by the
// same word.
func (p *parser) memberNameFollows(n int) bool {
	switch t := p.peek(n); t.Kind {
	case token.IDENT, token.STRING, token.NUMBER, token.BIGINT, token.LBRACK,
		token.PRIVATE_IDENT, token.MUL:
		return true
	default:
		return t.Kind.IsReserved()
	}
}

// --- parameters and arguments ----------------------------------------------

// parseTypeParams is TypeParameters (G.2).
func (p *parser) parseTypeParams() *ast.TypeParamList {
	if !p.at(token.LT) {
		return nil
	}
	tp := &ast.TypeParamList{Langle: p.next().Pos}
	for !p.at(token.GT) && !p.atEOF() {
		before := p.i
		tp.List = append(tp.List, p.parseTypeParam())
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	if p.at(token.GT) {
		tp.Rangle = p.next().Pos
	} else {
		p.errorf(p.cur(), "expected `>` to close type parameters, found %s", p.describe(p.cur()))
	}
	return tp
}

func (p *parser) parseTypeParam() *ast.TypeParam {
	tp := &ast.TypeParam{}
	for {
		t := p.cur()
		switch {
		case t.Kind == token.CONST:
		case t.Kind == token.IN:
		case t.Kind == token.IDENT && t.Ctx == token.CtxOut:
		default:
			goto done
		}
		// A modifier only if a name still follows; `<in>` alone is not valid.
		if p.peek(1).Kind != token.IDENT && p.peek(1).Kind != token.CONST &&
			p.peek(1).Kind != token.IN {
			goto done
		}
		p.next()
		tp.Mods = append(tp.Mods, ast.TypeParamMod{Pos: t.Pos, End: t.End, Tok: t.Kind})
	}
done:
	tp.Name = p.parseIdent()

	// G.2 admits `T extends U` or `T : U`. Both spellings reach the same later
	// phase; the parser records which was written.
	switch {
	case p.at(token.EXTENDS):
		tp.ExtendsPos = p.next().Pos
		tp.Constraint = p.parseType()
	case p.at(token.COLON):
		tp.ColonPos = p.next().Pos
		tp.Type = p.parseType()
	}
	if p.at(token.ASSIGN) {
		tp.Assign = p.next().Pos
		tp.Default = p.parseType()
	}
	return tp
}

// parseTypeArgs is TypeArguments (G.2).
//
// This is called from speculation (§6.1 site 1) as well as directly, so it
// must not error out on a failure to close — the caller checks Rangle and
// rolls back, and a diagnostic emitted here would be truncated anyway.
func (p *parser) parseTypeArgs() *ast.TypeArgList {
	if !p.at(token.LT) {
		return nil
	}
	ta := &ast.TypeArgList{Langle: p.next().Pos}
	for !p.at(token.GT) && !p.atEOF() {
		before := p.i
		ta.List = append(ta.List, p.parseType())
		if !p.got(token.COMMA) {
			break
		}
		if !p.advanced(before) {
			return ta
		}
	}
	if p.at(token.GT) {
		ta.Rangle = p.next().Pos
	}
	return ta
}

// parseParamTypeList is ParameterTypeList (G.4).
//
// Distinct from parseParamList (D.1): a type-position parameter takes no
// decorators, no modifiers, and no initializer, and sharing the function would
// mean accepting them here and rejecting them later for no gain.
func (p *parser) parseParamTypeList() *ast.ParamTypeList {
	pl := &ast.ParamTypeList{Lparen: p.expect(token.LPAREN)}
	for !p.at(token.RPAREN) && !p.atEOF() {
		before := p.i
		switch {
		case p.at(token.THIS):
			t := p.next()
			pl.This = &ast.ThisParam{ThisPos: t.Pos, Type: p.parseTypeAnnotation()}
		case p.at(token.ELLIPSIS):
			t := p.next()
			r := &ast.RestElem{Ellipsis: t.Pos, X: p.parseBindingTarget()}
			r.Type = p.parseTypeAnnotation()
			pl.Rest = r
		default:
			pl.List = append(pl.List, p.parseParamType())
		}
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	if p.at(token.RPAREN) {
		pl.Rparen = p.next().Pos
	}
	return pl
}

func (p *parser) parseParamType() *ast.ParamType {
	pt := &ast.ParamType{Name: p.parseBindingTarget()}
	if p.at(token.QUESTION) {
		pt.OptionalPos = p.next().Pos
	}
	pt.Type = p.parseTypeAnnotation()
	return pt
}

func (p *parser) badType() ast.TypeExpr {
	t := p.cur()
	end := t.End
	if end == t.Pos {
		end = t.Pos + 1
	}
	return &ast.BadType{From: t.Pos, To: end}
}