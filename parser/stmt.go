package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
)

// parseBlockStmt parses `{ StatementList }`.
//
// The depth reset is the crux of the terminator rule. Every brace that holds a
// comma-separated list treats a line break as white space; this one holds
// statements, which do end at line terminators. Saving and zeroing depth is
// what lets one counter serve both.
func (p *parser) parseBlockStmt() *ast.BlockStmt {
	b := &ast.BlockStmt{}
	saved := p.enterTerminated()

	b.Lbrace = p.expect(token.LBRACE)
	b.List = p.parseStmtList(func() bool { return p.at(token.RBRACE) })
	b.Rbrace = p.expect(token.RBRACE)

	p.leave(saved)
	return b
}

// parseStmtList parses `[ terminator ] { Statement terminator }`. The leading
// optional terminator needs no code: a run of line terminators is one
// terminator, and it reaches the parser as a flag on the token that follows.
func (p *parser) parseStmtList(stop func() bool) []ast.Stmt {
	var list []ast.Stmt
	for !stop() && !p.at(token.EOF) {
		before := p.tok.Pos
		list = append(list, p.parseStmt())
		p.expectTerminator()
		if p.stalled(before) {
			continue
		}
	}
	return list
}

func (p *parser) parseStmt() ast.Stmt {
	switch p.tok.Kind {
	case token.LBRACE:
		return p.parseBlockStmt()

	case token.LET, token.VAR:
		return &ast.DeclStmt{Decl: p.parseVarDecl(p.leadComment)}

	case token.IF:
		return p.parseIfStmt()
	case token.WHILE:
		return p.parseWhileStmt()
	case token.FOR:
		return p.parseForStmt()
	case token.SWITCH:
		return p.parseSwitchStmt()
	case token.SELECT:
		return p.parseSelectStmt()

	case token.RETURN:
		return p.parseReturnStmt()
	case token.DEFER:
		return p.parseDeferStmt()

	case token.BREAK, token.CONTINUE, token.FALLTHROUGH:
		s := &ast.BranchStmt{TokPos: p.tok.Pos, Tok: p.tok.Kind}
		p.advanceToken()
		return s

	case token.ELSE:
		// An `else` reached as a statement start means a line terminator
		// already ended the if statement it belonged to.
		p.errorHere(diag.ExpectedToken, "a statement", p.describe(p.tok))
		bad := &ast.BadStmt{From: p.tok.Pos, To: p.tok.End()}
		p.advance(stmtStart)
		return bad

	case token.RBRACE, token.EOF:
		return &ast.BadStmt{From: p.tok.Pos, To: p.tok.Pos}
	}

	return p.parseSimpleStmt()
}

// parseSimpleStmt parses an assignment or an expression statement.
//
// Assignment is a statement and never an expression, so the decision is made
// after the target list is parsed rather than by lookahead — and there is no
// `=` inside any condition for it to be confused with. An assignment target is
// a bare PrimaryExpr, so a dereference-write and the blank identifier both
// arrive as ordinary nodes; which shapes are assignable is a static rule.
func (p *parser) parseSimpleStmt() ast.Stmt {
	first := p.parseExpr()

	if p.continues() && p.tok.Kind.IsCompoundAssign() {
		op, opPos := p.tok.Kind, p.tok.Pos
		p.advanceToken()
		val := p.parseExpr()
		return &ast.AssignStmt{
			Targets: []ast.Expr{first},
			OpPos:   opPos, Op: op,
			Values: []ast.Expr{val},
		}
	}

	targets := []ast.Expr{first}
	for p.continues() && p.got(token.COMMA) {
		targets = append(targets, p.parseExpr())
	}

	if p.continues() && p.at(token.ASSIGN) {
		opPos := p.tok.Pos
		p.advanceToken()
		var values []ast.Expr
		for {
			values = append(values, p.parseExpr())
			if !p.got(token.COMMA) {
				break
			}
		}
		return &ast.AssignStmt{Targets: targets, OpPos: opPos, Op: token.ASSIGN, Values: values}
	}

	if len(targets) > 1 {
		p.errorAt(diag.ExpectedToken, targets[1].Pos(), "'='", "a comma-separated expression list")
	}
	// An ExpressionStatement may not be a bare composite or map literal. A
	// statement that opens with `{` is already a Block by dispatch above; the
	// remaining case is a typed literal, which needs no parse decision and is
	// left to the analyzer.
	return &ast.ExprStmt{X: first}
}

// parseIfStmt parses an if statement. There is no initializer clause: the
// two-statement error-checking idiom is intentional.
//
// The else clause is attached only when no line terminator precedes it, since
// inside a block a line terminator ends the statement.
func (p *parser) parseIfStmt() *ast.IfStmt {
	s := &ast.IfStmt{If: p.expect(token.IF)}
	s.Cond = p.parseHeaderExpr("if")
	s.Body = p.parseBlockStmt()

	if p.continues() && p.at(token.ELSE) {
		p.advanceToken()
		switch {
		case p.at(token.IF):
			s.Else = p.parseIfStmt()
		case p.at(token.LBRACE):
			s.Else = p.parseBlockStmt()
		default:
			p.errorHere(diag.ExpectedToken, "'{' or 'if'", p.describe(p.tok))
			s.Else = &ast.BadStmt{From: p.tok.Pos, To: p.tok.End()}
			p.advance(stmtStart)
		}
	}
	return s
}

// parseWhileStmt parses the only loop primitive.
func (p *parser) parseWhileStmt() *ast.WhileStmt {
	s := &ast.WhileStmt{While: p.expect(token.WHILE)}
	s.Cond = p.parseHeaderExpr("while")
	s.Body = p.parseBlockStmt()
	return s
}

// parseForStmt parses `for IterationBinding in Expression Block`.
//
// The mode marker sits on the binding rather than on the iterable, because what
// transfers is each element, one per iteration. The marker and the two-name
// form do not combine, but both are parsed together here so the combination is
// diagnosed as itself rather than as a syntax error at the comma.
func (p *parser) parseForStmt() *ast.ForStmt {
	s := &ast.ForStmt{For: p.expect(token.FOR)}

	if p.at(token.MUT) || p.at(token.VAR) {
		s.ModePos, s.Mode = p.tok.Pos, p.tok.Kind
		p.advanceToken()
	}

	s.Names = append(s.Names, p.expectIdent())
	if p.got(token.COMMA) {
		s.Names = append(s.Names, p.expectIdent())
	}

	s.In = p.expect(token.IN)
	s.X = p.parseHeaderExpr("for")
	s.Body = p.parseBlockStmt()
	return s
}

// parseSwitchStmt parses a switch. Its body and each clause's statement list
// are terminator-significant.
func (p *parser) parseSwitchStmt() *ast.SwitchStmt {
	s := &ast.SwitchStmt{Switch: p.expect(token.SWITCH)}
	s.Tag = p.parseHeaderExpr("switch")

	saved := p.enterTerminated()
	s.Lbrace = p.expect(token.LBRACE)

	seenDefault := false
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		before := p.tok.Pos
		isDefault := p.at(token.DEFAULT)
		if isDefault && seenDefault {
			p.errorHere(diag.DuplicateDefault, "switch")
		}
		seenDefault = seenDefault || isDefault
		s.Cases = append(s.Cases, p.parseCaseClause())
		if p.stalled(before) {
			continue
		}
	}
	s.Rbrace = p.expect(token.RBRACE)
	p.leave(saved)
	return s
}

func (p *parser) parseCaseClause() *ast.CaseClause {
	c := &ast.CaseClause{Case: p.tok.Pos}

	switch {
	case p.at(token.CASE):
		p.advanceToken()
		for {
			c.Patterns = append(c.Patterns, p.parsePattern())
			if !p.got(token.COMMA) {
				break
			}
		}
	case p.at(token.DEFAULT):
		// Patterns stays nil, which is what marks the default clause.
		p.advanceToken()
	default:
		p.errorHere(diag.ExpectedCase, p.describe(p.tok))
		p.advance(clauseStart)
		return c
	}

	c.Colon = p.expect(token.COLON)
	c.Body = p.parseStmtList(func() bool {
		return p.at(token.CASE) || p.at(token.DEFAULT) || p.at(token.RBRACE)
	})
	return c
}

// parsePattern parses one Pattern. In pattern position a leading `.` is always
// an enum pattern, never an enum shorthand reached through Expression: the
// payload entries are binding names rather than expressions, and they are views
// into the payload rather than copies.
func (p *parser) parsePattern() ast.Expr {
	if !p.at(token.PERIOD) {
		saved := p.noLit
		p.noLit = true
		x := p.parseExpr()
		p.noLit = saved
		return x
	}

	pat := &ast.EnumPattern{Dot: p.tok.Pos}
	p.advanceToken()
	pat.Name = p.expectIdent()

	if p.at(token.LPAREN) {
		pat.Lparen = p.open(token.LPAREN)
		for !p.at(token.RPAREN) && !p.at(token.EOF) {
			before := p.tok.Pos
			pat.Binds = append(pat.Binds, p.expectIdent())
			if !p.got(token.COMMA) {
				break
			}
			if p.stalled(before) {
				continue
			}
		}
		pat.Rparen = p.close(token.RPAREN)
	}
	return pat
}

// parseSelectStmt parses a select. Its body and each clause's statement list
// are terminator-significant.
func (p *parser) parseSelectStmt() *ast.SelectStmt {
	s := &ast.SelectStmt{Select: p.expect(token.SELECT)}

	saved := p.enterTerminated()
	s.Lbrace = p.expect(token.LBRACE)

	seenDefault := false
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		before := p.tok.Pos
		isDefault := p.at(token.DEFAULT)
		if isDefault && seenDefault {
			p.errorHere(diag.DuplicateDefault, "select")
		}
		seenDefault = seenDefault || isDefault
		s.Cases = append(s.Cases, p.parseSelectClause())
		if p.stalled(before) {
			continue
		}
	}
	s.Rbrace = p.expect(token.RBRACE)
	p.leave(saved)
	return s
}

// parseSelectClause parses one clause, covering all three channel-case forms.
//
// The declaring form introduces bindings scoped to this clause's body; the
// assigning form writes to pre-declared targets; the bare form has neither.
// Which calls are admissible in this position, and the rule that one select is
// entirely bare or entirely awaited, are static rules — the shape checked here
// is only "a call, optionally awaited".
func (p *parser) parseSelectClause() *ast.SelectClause {
	c := &ast.SelectClause{Case: p.tok.Pos}

	switch {
	case p.at(token.DEFAULT):
		p.advanceToken()

	case p.at(token.CASE):
		p.advanceToken()

		switch {
		case p.at(token.LET), p.at(token.VAR):
			c.KwPos, c.Kw = p.tok.Pos, p.tok.Kind
			p.advanceToken()
			for {
				c.Bindings = append(c.Bindings, p.parseBinding())
				if !p.got(token.COMMA) {
					break
				}
			}
			c.Assign = p.expect(token.ASSIGN)
			c.Op = p.parseChannelOp()

		default:
			x := p.parseExpr()
			switch {
			case p.at(token.ASSIGN):
				c.Targets = append(c.Targets, x)
				c.Assign = p.tok.Pos
				p.advanceToken()
				c.Op = p.parseChannelOp()
			case p.at(token.COMMA):
				c.Targets = append(c.Targets, x)
				for p.got(token.COMMA) {
					c.Targets = append(c.Targets, p.parseExpr())
				}
				c.Assign = p.expect(token.ASSIGN)
				c.Op = p.parseChannelOp()
			default:
				c.Op = p.requireCall(x, "select")
			}
		}

	default:
		p.errorHere(diag.ExpectedCase, p.describe(p.tok))
		p.advance(clauseStart)
		return c
	}

	c.Colon = p.expect(token.COLON)
	c.Body = p.parseStmtList(func() bool {
		return p.at(token.CASE) || p.at(token.DEFAULT) || p.at(token.RBRACE)
	})
	return c
}

// parseChannelOp parses `CallExpr` or `"await" CallExpr`.
func (p *parser) parseChannelOp() ast.Expr {
	return p.requireCall(p.parseExpr(), "select")
}

// requireCall enforces the named "a call and nothing else" restriction shared
// by a launch prefix, `defer`, and a channel operation. An awaited call
// satisfies it where the position admits one.
func (p *parser) requireCall(x ast.Expr, kw string) ast.Expr {
	inner := x
	if a, ok := inner.(*ast.AwaitExpr); ok {
		inner = a.X
	}
	if _, ok := inner.(*ast.CallExpr); !ok {
		p.errorSpan(diag.ExpectedCall, x.Pos(), x.End(), kw)
	}
	return x
}

// parseReturnStmt parses a return. A multi-value return is a bare comma list,
// never parenthesized: parentheses construct a tuple, bare commas unbuild one.
func (p *parser) parseReturnStmt() *ast.ReturnStmt {
	s := &ast.ReturnStmt{Return: p.expect(token.RETURN)}

	if p.continues() && !p.at(token.RBRACE) && !p.at(token.EOF) {
		for {
			s.Results = append(s.Results, p.parseExpr())
			if !p.got(token.COMMA) {
				break
			}
		}
	}
	return s
}

// parseDeferStmt parses `defer CallExpr`. Arguments are evaluated at
// registration; only the call is postponed.
func (p *parser) parseDeferStmt() ast.Stmt {
	pos := p.expect(token.DEFER)
	x := p.parseExpr()

	call, ok := x.(*ast.CallExpr)
	if !ok {
		p.errorSpan(diag.ExpectedCall, x.Pos(), x.End(), "defer")
		return &ast.BadStmt{From: pos, To: x.End()}
	}
	return &ast.DeferStmt{Defer: pos, Call: call}
}