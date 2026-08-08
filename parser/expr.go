package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
)

// Expressions — vertex_grammar.md section B.

// inFlag is the [In] parameter threaded through B.4 and B.5. It controls
// exactly one thing: whether `in` is a RelationalExpression operator. The
// three-clause `for` head clears it.
type inFlag bool

const (
	allowIn inFlag = true
	noIn    inFlag = false
)

// parseExpr is Expression (B.5): the comma operator over
// AssignmentExpressions.
func (p *parser) parseExpr(in inFlag) ast.Expr {
	defer p.trace("Expr")()

	x := p.parseAssign(in)
	if !p.at(token.COMMA) {
		return x
	}
	seq := &ast.SeqExpr{Exprs: []ast.Expr{x}}
	for p.got(token.COMMA) {
		seq.Exprs = append(seq.Exprs, p.parseAssign(in))
	}
	return seq
}

// parseAssign is AssignmentExpression (B.5).
func (p *parser) parseAssign(in inFlag) ast.Expr {
	defer p.trace("AssignmentExpression")()
	if !p.enter() {
		return p.badExpr()
	}
	defer p.leave()

	// YieldExpression is gated on [+Yield] in the grammar, but the parser does
	// not thread Yield/Await as flags: §1 makes [+Await] unconditional at file
	// scope, and a `yield` outside a generator is an early error, not a parse
	// failure (§6.3). So both are parsed wherever they appear and rejected by
	// name later.
	if p.at(token.YIELD) {
		return p.parseYield(in)
	}

	// Arrow functions must be tried before the conditional chain: their heads
	// are prefixes of ParenthesizedExpression and of a type assertion.
	if x, ok := p.tryArrow(in); ok {
		return x
	}

	x := p.parseConditional(in)

	if op := p.kind(); isAssignOp(op) {
		t := p.next()
		return &ast.AssignExpr{
			Lhs: p.toTarget(x, op), OpPos: t.Pos, OpEnd: t.End, Op: op,
			Rhs: p.parseAssign(in),
		}
	}
	return x
}

func isAssignOp(k token.Kind) bool {
	switch k {
	case token.ASSIGN, token.ADD_ASSIGN, token.SUB_ASSIGN, token.MUL_ASSIGN,
		token.QUO_ASSIGN, token.REM_ASSIGN, token.EXP_ASSIGN, token.SHL_ASSIGN,
		token.SHR_ASSIGN, token.USHR_ASSIGN, token.AND_ASSIGN, token.OR_ASSIGN,
		token.XOR_ASSIGN, token.LAND_ASSIGN, token.LOR_ASSIGN, token.COALESCE_ASSIGN:
		return true
	}
	return false
}

// toTarget reinterprets an ObjectLit or ArrayLit on the left of `=` as a
// destructuring pattern (B.6).
//
// This is the grammar's only cover (K) and the only in-parser reinterpretation
// (§5.1). It runs only for simple assignment: `[a] += 1` has no pattern
// reading, and reinterpreting there would invent a node the grammar does not
// derive.
func (p *parser) toTarget(x ast.Expr, op token.Kind) ast.Expr {
	if op != token.ASSIGN {
		return x
	}
	switch x := x.(type) {
	case *ast.ObjectLit:
		return p.objectToPattern(x)
	case *ast.ArrayLit:
		return p.arrayToPattern(x)
	}
	return x
}

func (p *parser) objectToPattern(o *ast.ObjectLit) ast.Expr {
	pat := &ast.ObjectPattern{Lbrace: o.Lbrace, Rbrace: o.Rbrace}
	for _, prop := range o.Props {
		switch prop := prop.(type) {
		case *ast.PropertyDef:
			if prop.IsCover {
				// `{ a = 1 }` — the CoverInitializedName finally becomes legal.
				pat.Props = append(pat.Props, &ast.PropertyPattern{
					Value: &ast.AssignPattern{Lhs: prop.Key, Assign: token.NoPos, Rhs: prop.Value},
				})
				continue
			}
			pat.Props = append(pat.Props, &ast.PropertyPattern{
				Key: prop.Key, Colon: prop.Colon, Value: p.toNestedPattern(prop.Value),
			})
		case *ast.SpreadElem:
			pat.Props = append(pat.Props, &ast.RestElem{Ellipsis: prop.Ellipsis, X: prop.X})
		case *ast.Ident:
			pat.Props = append(pat.Props, &ast.PropertyPattern{Value: prop})
		default:
			p.errorAt(prop.Pos(), prop.End(), "not a valid assignment target")
			pat.Props = append(pat.Props, prop)
		}
	}
	return pat
}

func (p *parser) arrayToPattern(a *ast.ArrayLit) ast.Expr {
	pat := &ast.ArrayPattern{Lbrack: a.Lbrack, Rbrack: a.Rbrack}
	for _, elt := range a.Elts {
		switch elt := elt.(type) {
		case *ast.SpreadElem:
			pat.Elts = append(pat.Elts, &ast.RestElem{Ellipsis: elt.Ellipsis, X: elt.X})
		case *ast.Elision:
			pat.Elts = append(pat.Elts, elt)
		default:
			pat.Elts = append(pat.Elts, p.toNestedPattern(elt))
		}
	}
	return pat
}

func (p *parser) toNestedPattern(x ast.Expr) ast.Expr {
	switch x := x.(type) {
	case *ast.ObjectLit:
		return p.objectToPattern(x)
	case *ast.ArrayLit:
		return p.arrayToPattern(x)
	case *ast.AssignExpr:
		if x.Op == token.ASSIGN {
			return &ast.AssignPattern{Lhs: p.toNestedPattern(x.Lhs), Assign: x.OpPos, Rhs: x.Rhs}
		}
	}
	return x
}

func (p *parser) parseYield(in inFlag) ast.Expr {
	t := p.next()
	y := &ast.YieldExpr{YieldPos: t.Pos, YieldEnd: t.End}

	// `yield [no LineTerminator here] AssignmentExpression` and the `*` form.
	// A line break makes it a bare yield, which is the whole point of the
	// restriction.
	if p.nlBefore() {
		return y
	}
	if p.at(token.MUL) {
		p.next()
		y.Delegate = true
		y.X = p.parseAssign(in)
		return y
	}
	if p.startsExpr() {
		y.X = p.parseAssign(in)
	}
	return y
}

// parseConditional is ConditionalExpression (B.4).
func (p *parser) parseConditional(in inFlag) ast.Expr {
	cond := p.parseBinary(in, 1)
	if !p.at(token.QUESTION) {
		return cond
	}
	q := p.next().Pos
	// The consequent is [+In] unconditionally, per B.4.
	then := p.parseAssign(allowIn)
	colon := p.expect(token.COLON)
	return &ast.CondExpr{Cond: cond, Quest: q, Then: then, Colon: colon, Else: p.parseAssign(in)}
}

// parseBinary climbs the binary chain in B.4 using the precedence table on
// token.Kind (§3, §6).
//
// Mixing `??` with `||` or `&&` is ungrammatical but parses here and is
// rejected by name later (§6.3), so the message can say what to parenthesize.
func (p *parser) parseBinary(in inFlag, minPrec int) ast.Expr {
	if !p.enter() {
		return p.badExpr()
	}
	defer p.leave()

	x := p.parseUnary()

	for {
		op, opPos, opEnd, width := p.binaryOp(in)
		prec := op.Precedence()
		if prec < minPrec {
			break
		}
		for n := 0; n < width; n++ {
			p.next()
		}

		// ** is right-associative and its left operand is an
		// UpdateExpression, not a UnaryExpression (B.4). Both are exceptions
		// to the plain climb.
		next := prec + 1
		if op == token.EXP {
			next = prec
		}
		y := p.parseBinary(in, next)
		x = &ast.BinaryExpr{X: x, OpPos: opPos, OpEnd: opEnd, Op: op, Y: y}

		if op == token.COALESCE || op == token.LOR || op == token.LAND {
			p.checkCoalesceMix(x.(*ast.BinaryExpr))
		}
	}
	return x
}

// binaryOp identifies the infix operator at the cursor, joining a run of
// adjacent GT tokens (§4.2). width is how many tokens to consume.
//
// The scanner under-munches `>` so that `Array<Box<int32>>` comes apart in
// type context; joining is the expression parser's job and this is where it
// happens. Adjacency is the whitespace test, so `a > > b` never joins.
func (p *parser) binaryOp(in inFlag) (op token.Kind, pos, end token.Pos, width int) {
	t := p.cur()

	switch t.Kind {
	case token.GT:
		gts, assign, w := 1, false, 1
		for {
			nx := p.peek(w)
			if !p.peek(w - 1).Adjacent(nx) {
				break
			}
			if nx.Kind == token.GT && gts < 3 && !assign {
				gts++
				w++
				continue
			}
			if nx.Kind == token.ASSIGN {
				assign = true
				w++
			}
			break
		}
		joined := token.JoinGT(gts, assign)
		return joined, t.Pos, p.peek(w - 1).End, w

	case token.IN:
		if in == noIn {
			return token.INVALID, token.NoPos, token.NoPos, 0
		}
		return token.IN, t.Pos, t.End, 1

	case token.IDENT:
		// `as`, `as const`, and `satisfies` bind at relational precedence but
		// are IDENT tokens, so they are not in the precedence table (§3). They
		// are handled by parsePostfixTypeOps instead, which runs below unary.
		return token.INVALID, token.NoPos, token.NoPos, 0
	}

	if t.Kind.IsBinaryOperator() {
		return t.Kind, t.Pos, t.End, 1
	}
	return token.INVALID, token.NoPos, token.NoPos, 0
}

func (p *parser) checkCoalesceMix(x *ast.BinaryExpr) {
	mixes := func(y ast.Expr) bool {
		b, ok := y.(*ast.BinaryExpr)
		if !ok {
			return false
		}
		switch x.Op {
		case token.COALESCE:
			return b.Op == token.LOR || b.Op == token.LAND
		default:
			return b.Op == token.COALESCE
		}
	}
	if mixes(x.X) || mixes(x.Y) {
		p.errorAt(x.OpPos, x.OpEnd,
			"`??` cannot be mixed with `||` or `&&` without parentheses")
	}
}

// parseUnary is UnaryExpression (B.4).
func (p *parser) parseUnary() ast.Expr {
	if !p.enter() {
		return p.badExpr()
	}
	defer p.leave()

	switch k := p.kind(); k {
	case token.DELETE, token.VOID, token.TYPEOF, token.ADD, token.SUB,
		token.TILDE, token.NOT:
		t := p.next()
		return &ast.UnaryExpr{OpPos: t.Pos, Op: k, X: p.parseUnary()}

	case token.AWAIT:
		t := p.next()
		return &ast.AwaitExpr{AwaitPos: t.Pos, X: p.parseUnary()}

	case token.INC, token.DEC:
		t := p.next()
		return &ast.UpdateExpr{OpPos: t.Pos, OpEnd: t.End, Op: k, Prefix: true, X: p.parseUnary()}

	case token.LT:
		// `< Type > UnaryExpression` (B.4).
		//
		// §1 claims this production owns the prefix-`<` position with no
		// competing reading. It does not: a generic arrow head is also `<`
		// (D.2, via ArrowFormalParameters), and §6.1 item 3 says as much. The
		// generic-arrow attempt already ran in tryArrow before we got here, so
		// reaching this point means it failed — which makes the assertion
		// reading correct by elimination rather than by uniqueness.
		return p.parseTypeAssert()
	}

	x := p.parsePostfix(p.parseLHS())
	return x
}

func (p *parser) parseTypeAssert() ast.Expr {
	langle := p.next().Pos
	t := p.parseType()
	rangle := p.expect(token.GT)
	return &ast.TypeAssertExpr{Langle: langle, Type: t, Rangle: rangle, X: p.parseUnary()}
}

// parsePostfix applies the postfix operators that bind tighter than any binary
// operator: `++`, `--`, `!`, `as`, and `satisfies`.
func (p *parser) parsePostfix(x ast.Expr) ast.Expr {
	for {
		switch {
		case (p.at(token.INC) || p.at(token.DEC)) && !p.nlBefore():
			// LeftHandSideExpression [no LineTerminator here] ++ / -- (L).
			t := p.next()
			x = &ast.UpdateExpr{OpPos: t.Pos, OpEnd: t.End, Op: t.Kind, X: x}

		case p.atCtx(token.CtxAs) && !p.nlBefore():
			// RelationalExpression [no LineTerminator here] as Type / as const
			t := p.next()
			as := &ast.AsExpr{X: x, OpPos: t.Pos, Op: token.CtxAs}
			if p.at(token.CONST) {
				c := p.next()
				as.IsConst = true
				as.ConstEnd = c.End
			} else {
				as.Type = p.parseType()
			}
			x = as

		case p.atCtx(token.CtxSatisfies) && !p.nlBefore():
			t := p.next()
			x = &ast.AsExpr{X: x, OpPos: t.Pos, Op: token.CtxSatisfies, Type: p.parseType()}

		default:
			return x
		}
	}
}

// parseLHS is LeftHandSideExpression (B.3): NewExpression, CallExpression, and
// OptionalExpression, which share a left-recursive suffix loop.
func (p *parser) parseLHS() ast.Expr {
	var x ast.Expr

	switch {
	case p.at(token.NEW):
		x = p.parseNew()
	case p.at(token.SUPER):
		t := p.next()
		x = &ast.SuperExpr{SuperPos: t.Pos, SuperEnd: t.End}
	case p.at(token.IMPORT):
		x = p.parseImportExpr()
	default:
		x = p.parsePrimary()
	}
	return p.parseSuffixes(x)
}

// parseSuffixes is the member / call / optional-chain loop (B.3).
func (p *parser) parseSuffixes(x ast.Expr) ast.Expr {
	for {
		switch {
		case p.at(token.PERIOD):
			dot := p.next().Pos
			x = &ast.MemberExpr{X: x, Dot: dot, Sel: p.parseMemberName()}

		case p.at(token.QUESTION_DOT):
			dot := p.next().Pos
			switch {
			case p.at(token.LPAREN):
				x = p.finishCall(x, nil, true)
			case p.at(token.LBRACK):
				x = p.finishIndex(x, true)
			default:
				x = &ast.MemberExpr{X: x, Optional: true, Dot: dot, Sel: p.parseMemberName()}
			}

		case p.at(token.LBRACK):
			x = p.finishIndex(x, false)

		case p.at(token.LPAREN):
			x = p.finishCall(x, nil, false)

		case p.at(token.TEMPLATE) || p.at(token.TEMPLATE_HEAD):
			x = &ast.TaggedTemplateExpr{Tag: x, Template: p.parseTemplate()}

		case p.at(token.NOT) && !p.nlBefore():
			// MemberExpression / CallExpression / OptionalChain
			// [no LineTerminator here] ! (B.3, L).
			t := p.next()
			x = &ast.NonNullExpr{X: x, Bang: t.Pos, BangEnd: t.End}

		case p.at(token.LT) && !p.nlBefore():
			// Instantiation expressions (§6.1 site 1). Speculate
			// TypeArguments; commit only if it closes *and* the next token is
			// in InstantiationFollowSet. This stands in for a cover
			// nonterminal — the grammar states the rule declaratively and the
			// parser pays for it here.
			var args *ast.TypeArgList
			ok := p.speculate(func() bool {
				args = p.parseTypeArgs()
				if args == nil || args.Rangle == token.NoPos {
					return false
				}
				return p.at(token.LPAREN) || p.at(token.TEMPLATE) ||
					p.at(token.TEMPLATE_HEAD) || p.inInstantiationFollowSet()
			})
			if !ok {
				return x
			}
			switch {
			case p.at(token.LPAREN):
				x = p.finishCall(x, args, false)
			case p.at(token.TEMPLATE) || p.at(token.TEMPLATE_HEAD):
				x = &ast.TaggedTemplateExpr{Tag: x, TypeArgs: args, Template: p.parseTemplate()}
			default:
				x = &ast.InstantiationExpr{X: x, TypeArgs: args}
			}

		default:
			return x
		}
	}
}

// inInstantiationFollowSet implements the named set in B.3.
//
// The `>` entries are single GT tokens here, since the scanner never merges
// them — a `>>` in the source is two members of this set in a row, and neither
// is in it.
func (p *parser) inInstantiationFollowSet() bool {
	switch p.kind() {
	case token.LPAREN, token.TEMPLATE, token.TEMPLATE_HEAD, token.RPAREN,
		token.RBRACK, token.RBRACE, token.COLON, token.SEMI, token.COMMA,
		token.QUESTION, token.ASSIGN, token.EQL, token.STRICT_EQL, token.NEQ,
		token.STRICT_NEQ, token.OR, token.AND, token.LOR, token.LAND,
		token.COALESCE, token.EOF:
		return true
	}
	return false
}

func (p *parser) parseMemberName() ast.Node {
	if p.at(token.PRIVATE_IDENT) {
		t := p.next()
		return &ast.PrivateIdent{HashPos: t.Pos, NameEnd: t.End}
	}
	return p.parseIdentName()
}

func (p *parser) finishIndex(x ast.Expr, optional bool) ast.Expr {
	lb := p.next().Pos
	idx := p.parseExpr(allowIn)
	return &ast.IndexExpr{X: x, Optional: optional, Lbrack: lb, Index: idx, Rbrack: p.expect(token.RBRACK)}
}

func (p *parser) finishCall(fun ast.Expr, args *ast.TypeArgList, optional bool) ast.Expr {
	lp := p.next().Pos
	call := &ast.CallExpr{Fun: fun, Optional: optional, TypeArgs: args, Lparen: lp}
	call.Args = p.parseArgs()
	call.Rparen = p.expect(token.RPAREN)
	return call
}

// parseArgs is ArgumentList (B.3), including the trailing comma, which leaves
// no trace (§5.4).
func (p *parser) parseArgs() []ast.Expr {
	var out []ast.Expr
	for !p.at(token.RPAREN) && !p.atEOF() {
		before := p.i
		if p.at(token.ELLIPSIS) {
			t := p.next()
			out = append(out, &ast.SpreadElem{Ellipsis: t.Pos, X: p.parseAssign(allowIn)})
		} else {
			out = append(out, p.parseAssign(allowIn))
		}
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	return out
}

func (p *parser) parseNew() ast.Expr {
	newPos := p.next().Pos

	// NewTarget: `new . target` (B.3).
	if p.at(token.PERIOD) {
		p.next()
		return &ast.MetaProp{MetaPos: newPos, Meta: token.NEW, Prop: p.parseIdentName()}
	}

	// `new NewExpression` with no arguments is its own production, so the
	// callee is parsed without the call suffix and Lparen stays NoPos.
	callee := p.parseMemberOnly()
	n := &ast.NewExpr{NewPos: newPos, Callee: callee}

	if p.at(token.LT) {
		p.speculate(func() bool {
			args := p.parseTypeArgs()
			if args == nil || args.Rangle == token.NoPos || !p.at(token.LPAREN) {
				return false
			}
			n.TypeArgs = args
			return true
		})
	}
	if p.at(token.LPAREN) {
		n.Lparen = p.next().Pos
		n.Args = p.parseArgs()
		n.Rparen = p.expect(token.RPAREN)
	}
	return n
}

// parseMemberOnly parses a MemberExpression without consuming a call, for the
// callee of `new`.
func (p *parser) parseMemberOnly() ast.Expr {
	var x ast.Expr
	if p.at(token.NEW) {
		x = p.parseNew()
	} else {
		x = p.parsePrimary()
	}
	for {
		switch {
		case p.at(token.PERIOD):
			dot := p.next().Pos
			x = &ast.MemberExpr{X: x, Dot: dot, Sel: p.parseMemberName()}
		case p.at(token.LBRACK):
			x = p.finishIndex(x, false)
		default:
			return x
		}
	}
}

// parseImportExpr covers ImportCall and ImportMeta (B.3).
func (p *parser) parseImportExpr() ast.Expr {
	importPos := p.next().Pos

	if p.at(token.PERIOD) {
		p.next()
		switch {
		case p.atCtx(token.CtxDefer), p.atCtx(token.CtxSource):
			t := p.next()
			phase := ast.PhaseDefer
			if t.Ctx == token.CtxSource {
				phase = ast.PhaseSource
			}
			c := &ast.ImportCall{ImportPos: importPos, Phase: phase, PhasePos: t.Pos}
			c.Lparen = p.expect(token.LPAREN)
			c.Args = p.parseArgs()
			c.Rparen = p.expect(token.RPAREN)
			return c
		default:
			return &ast.MetaProp{MetaPos: importPos, Meta: token.IMPORT, Prop: p.parseIdentName()}
		}
	}

	c := &ast.ImportCall{ImportPos: importPos, Phase: ast.PhaseEval}
	c.Lparen = p.expect(token.LPAREN)
	c.Args = p.parseArgs()
	c.Rparen = p.expect(token.RPAREN)
	return c
}

// parsePrimary is PrimaryExpression (B.1).
func (p *parser) parsePrimary() ast.Expr {
	defer p.trace("PrimaryExpression")()

	t := p.cur()
	switch t.Kind {
	case token.THIS:
		p.next()
		return &ast.ThisExpr{ThisPos: t.Pos, ThisEnd: t.End}

	case token.NUMBER, token.BIGINT, token.STRING, token.REGEX:
		p.next()
		return &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}

	case token.TRUE, token.FALSE, token.NULL:
		p.next()
		return &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End}

	case token.TEMPLATE, token.TEMPLATE_HEAD:
		return p.parseTemplate()

	case token.LBRACK:
		return p.parseArrayLit()

	case token.LBRACE:
		return p.parseObjectLit()

	case token.LPAREN:
		lp := p.next().Pos
		x := p.parseExpr(allowIn)
		return &ast.ParenExpr{Lparen: lp, X: x, Rparen: p.expect(token.RPAREN)}

	case token.FUNCTION:
		return &ast.FuncExpr{Fn: p.parseFunction(nil, token.NoPos, ast.AccelNone, false)}

	case token.CLASS:
		return &ast.ClassExpr{Class: p.parseClass(nil, token.NoPos)}

	case token.IDENT:
		switch t.Ctx {
		case token.CtxAsync:
			// `async function` — the [no LineTerminator here] restriction (L)
			// means a break forces the identifier reading.
			if p.peek(1).Kind == token.FUNCTION && !p.peek(1).NLBefore() {
				asyncPos := p.next().Pos
				return &ast.FuncExpr{Fn: p.parseFunction(nil, asyncPos, ast.AccelNone, true)}
			}
		case token.CtxAbstract:
			if p.peek(1).Kind == token.CLASS {
				abstractPos := p.next().Pos
				return &ast.ClassExpr{Class: p.parseClass(nil, abstractPos)}
			}
		}
		return p.parseIdent()

	case token.PRIVATE_IDENT:
		// Only legal as `#x in obj` (B.4). Parsed here and rejected by name
		// later if it appears anywhere else.
		p.next()
		return &ast.PrivateIdent{HashPos: t.Pos, NameEnd: t.End}
	}

	p.errorf(t, "expected an expression, found %s", p.describe(t))
	return p.badExpr()
}

func (p *parser) parseIdent() *ast.Ident {
	t := p.cur()
	if t.Kind != token.IDENT {
		p.errorf(t, "expected an identifier, found %s", p.describe(t))
		return &ast.Ident{NamePos: t.Pos, NameEnd: t.Pos + 1}
	}
	p.next()
	id := p.arena.newIdent()
	*id = ast.Ident{NamePos: t.Pos, NameEnd: t.End, Ctx: t.Ctx, Escaped: t.HasEscape()}
	return id
}

// parseIdentName accepts any IdentifierName, including reserved words, for
// positions where A.2 says IdentifierName rather than Identifier: after a dot,
// as a LiteralPropertyName, as an ImportSpecifier name.
func (p *parser) parseIdentName() *ast.Ident {
	t := p.cur()
	if t.Kind != token.IDENT && !t.Kind.IsReserved() {
		p.errorf(t, "expected a property name, found %s", p.describe(t))
		return &ast.Ident{NamePos: t.Pos, NameEnd: t.Pos + 1}
	}
	p.next()
	id := p.arena.newIdent()
	*id = ast.Ident{NamePos: t.Pos, NameEnd: t.End, Ctx: t.Ctx, Escaped: t.HasEscape()}
	return id
}

func (p *parser) parseTemplate() *ast.TemplateLit {
	lit := &ast.TemplateLit{}
	t := p.next()
	lit.Quasis = append(lit.Quasis, &ast.TemplateElem{Kind: t.Kind, Start: t.Pos, Stop: t.End})
	if t.Kind == token.TEMPLATE {
		return lit
	}
	for {
		lit.Exprs = append(lit.Exprs, p.parseExpr(allowIn))
		q := p.cur()
		if q.Kind != token.TEMPLATE_MIDDLE && q.Kind != token.TEMPLATE_TAIL {
			p.errorf(q, "expected a template continuation, found %s", p.describe(q))
			// Keep the quasi/expr invariant: len(Quasis) == len(Exprs)+1.
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

func (p *parser) parseArrayLit() ast.Expr {
	lb := p.next().Pos
	a := &ast.ArrayLit{Lbrack: lb}
	for !p.at(token.RBRACK) && !p.atEOF() {
		before := p.i
		switch {
		case p.at(token.COMMA):
			// An Elision is a real node, not a nil slot: a hole is meaningful.
			a.Elts = append(a.Elts, &ast.Elision{Comma: p.next().Pos})
			continue
		case p.at(token.ELLIPSIS):
			t := p.next()
			a.Elts = append(a.Elts, &ast.SpreadElem{Ellipsis: t.Pos, X: p.parseAssign(allowIn)})
		default:
			a.Elts = append(a.Elts, p.parseAssign(allowIn))
		}
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	a.Rbrack = p.expect(token.RBRACK)
	return a
}

func (p *parser) parseObjectLit() ast.Expr {
	lb := p.next().Pos
	o := &ast.ObjectLit{Lbrace: lb}
	for !p.at(token.RBRACE) && !p.atEOF() {
		before := p.i
		o.Props = append(o.Props, p.parseObjectProp(o))
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	o.Rbrace = p.expect(token.RBRACE)
	return o
}

func (p *parser) parseObjectProp(o *ast.ObjectLit) ast.Expr {
	if p.at(token.ELLIPSIS) {
		t := p.next()
		return &ast.SpreadElem{Ellipsis: t.Pos, X: p.parseAssign(allowIn)}
	}

	// A method definition in an object literal (B.2 → D.3).
	if m := p.tryObjectMethod(); m != nil {
		return m
	}

	key := p.parsePropertyName()

	switch {
	case p.at(token.COLON):
		colon := p.next().Pos
		return &ast.PropertyDef{Key: key, Colon: colon, Value: p.parseAssign(allowIn)}

	case p.at(token.ASSIGN):
		// CoverInitializedName (K). Legal only after reinterpretation to an
		// ObjectAssignmentPattern; the flag rides along so a node that was
		// never reinterpreted can be rejected by name.
		p.next()
		o.CoverInit = true
		return &ast.PropertyDef{Key: key, Colon: token.NoPos, IsCover: true, Value: p.parseAssign(allowIn)}

	default:
		// Shorthand. Must be an Identifier, not any IdentifierName.
		if id, ok := key.(*ast.Ident); ok {
			return id
		}
		p.errorAt(key.Pos(), key.End(), "expected `:` after this property name")
		return key
	}
}

func (p *parser) parsePropertyName() ast.Expr {
	switch {
	case p.at(token.LBRACK):
		lb := p.next().Pos
		x := p.parseAssign(allowIn)
		return &ast.ComputedKey{Lbrack: lb, X: x, Rbrack: p.expect(token.RBRACK)}
	case p.at(token.STRING), p.at(token.NUMBER), p.at(token.BIGINT):
		t := p.next()
		return &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
	default:
		return p.parseIdentName()
	}
}

// startsExpr reports whether the current token can begin an expression. Used
// where the grammar has an optional expression (`return`, `yield`) and ASI
// decides.
func (p *parser) startsExpr() bool {
	switch k := p.kind(); k {
	case token.IDENT, token.PRIVATE_IDENT, token.NUMBER, token.BIGINT, token.STRING,
		token.REGEX, token.TEMPLATE, token.TEMPLATE_HEAD, token.LPAREN, token.LBRACK,
		token.LBRACE, token.LT, token.NOT, token.TILDE, token.ADD, token.SUB,
		token.INC, token.DEC, token.THIS, token.SUPER, token.NEW, token.IMPORT,
		token.FUNCTION, token.CLASS, token.TRUE, token.FALSE, token.NULL,
		token.TYPEOF, token.VOID, token.DELETE, token.AWAIT, token.YIELD:
		return true
	}
	return false
}

func (p *parser) badExpr() ast.Expr {
	t := p.cur()
	// A Bad node holds a slot and must still have a non-zero span (§1, §5.4).
	end := t.End
	if end == t.Pos {
		end = t.Pos + 1
	}
	return &ast.BadExpr{From: t.Pos, To: end}
}