package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
)

// parseBlockStmt parses `{ StatementList }`.
//
// The depth reset is the crux of A.0.6. Every other brace in the language holds
// a comma-separated list where a line break is whitespace; this one holds
// statements, which do end at line terminators. Saving and zeroing depth is
// what lets one flag serve both.
func (p *parser) parseBlockStmt() *ast.BlockStmt {
	b := &ast.BlockStmt{}
	savedDepth, savedLit := p.depth, p.noLit
	p.depth, p.noLit = 0, false

	b.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		b.List = append(b.List, p.parseStmt())
	}
	b.Rbrace = p.expect(token.RBRACE)

	p.depth, p.noLit = savedDepth, savedLit
	return b
}

func (p *parser) parseStmt() ast.Stmt {
	switch p.tok.Kind {
	case token.LBRACE:
		return p.parseBlockStmt()

	case token.LET, token.VAR:
		s := &ast.DeclStmt{Decl: p.parseVarDecl()}
		p.expectStmtEnd()
		return s

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
		p.expectStmtEnd()
		return s

	case token.RBRACE, token.EOF:
		return &ast.BadStmt{From: p.tok.Pos, To: p.tok.Pos}
	}

	return p.parseSimpleStmt()
}

// parseSimpleStmt parses an assignment or an expression statement.
//
// A.5.2 ⊢ assignment is a statement and never an expression, so the decision is
// made after the target list is parsed rather than by lookahead — and there is
// no `=` inside a condition anywhere for it to be confused with.
func (p *parser) parseSimpleStmt() ast.Stmt {
	first := p.parseExpr()

	if p.continues() && p.tok.Kind.IsCompoundAssign() {
		op, opPos := p.tok.Kind, p.tok.Pos
		p.advanceToken()
		val := p.parseExpr()
		p.expectStmtEnd()
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
		p.expectStmtEnd()
		return &ast.AssignStmt{Targets: targets, OpPos: opPos, Op: token.ASSIGN, Values: values}
	}

	if len(targets) > 1 {
		p.errorAt(diag.ExpectedToken, targets[1].Pos(), "'='", "a comma-separated expression list")
	}
	p.expectStmtEnd()
	return &ast.ExprStmt{X: first}
}

// parseIfStmt parses A.5.4. There is no initializer clause: the error-checking
// idiom is a destructuring let followed by a plain if, and its verbosity is
// intentional.
func (p *parser) parseIfStmt() *ast.IfStmt {
	s := &ast.IfStmt{If: p.expect(token.IF)}
	s.Cond = p.parseHeaderExpr("if")
	s.Body = p.parseBlockStmt()

	if p.at(token.ELSE) {
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
		return s
	}
	p.expectStmtEnd()
	return s
}

func (p *parser) parseWhileStmt() *ast.WhileStmt {
	s := &ast.WhileStmt{While: p.expect(token.WHILE)}
	s.Cond = p.parseHeaderExpr("while")
	s.Body = p.parseBlockStmt()
	p.expectStmtEnd()
	return s
}

// parseForStmt parses A.5.6's single loop shape.
//
// The grammar's IterationBinding alternatives do not combine — `mut a, b` has
// no production — so a mode marker and a two-name form are parsed together here
// and their combination left to the analyzer, which can say so plainly rather
// than reporting a syntax error at the comma.
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
	p.expectStmtEnd()
	return s
}

func (p *parser) parseSwitchStmt() *ast.SwitchStmt {
	s := &ast.SwitchStmt{Switch: p.expect(token.SWITCH)}
	s.Tag = p.parseHeaderExpr("switch")

	savedDepth, savedLit := p.depth, p.noLit
	p.depth, p.noLit = 0, false

	s.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		s.Cases = append(s.Cases, p.parseCaseClause())
	}
	s.Rbrace = p.expect(token.RBRACE)

	p.depth, p.noLit = savedDepth, savedLit
	p.expectStmtEnd()
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
		// Patterns stays nil, which is what marks the default clause. A.5.7
		// permits at most one; a second is a static rejection.
		p.advanceToken()
	default:
		p.errorHere(diag.ExpectedCase, p.describe(p.tok))
		p.advance(map[token.Kind]bool{token.CASE: true, token.DEFAULT: true, token.RBRACE: true})
		return c
	}

	c.Colon = p.expect(token.COLON)
	for !p.at(token.CASE) && !p.at(token.DEFAULT) && !p.at(token.RBRACE) && !p.at(token.EOF) {
		c.Body = append(c.Body, p.parseStmt())
	}
	return c
}

// parsePattern parses a Pattern (A.5.7). An EnumPattern's payload entries are
// binding names rather than expressions, which is what separates it from the
// otherwise identical EnumShorthand.
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
			pat.Binds = append(pat.Binds, p.expectIdent())
			if !p.got(token.COMMA) {
				break
			}
		}
		pat.Rparen = p.close(token.RPAREN)
	}
	return pat
}

// parseSelectStmt parses A.10.2.
//
// Every case must be a channel receive; nothing else is legal in case position.
// The check is left to the analyzer so the diagnostic can quote the offending
// expression, but the shape parsed here is deliberately narrow.
func (p *parser) parseSelectStmt() *ast.SelectStmt {
	s := &ast.SelectStmt{Select: p.expect(token.SELECT)}

	savedDepth, savedLit := p.depth, p.noLit
	p.depth, p.noLit = 0, false

	s.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		s.Cases = append(s.Cases, p.parseSelectClause())
	}
	s.Rbrace = p.expect(token.RBRACE)

	p.depth, p.noLit = savedDepth, savedLit
	p.expectStmtEnd()
	return s
}

func (p *parser) parseSelectClause() *ast.SelectClause {
	c := &ast.SelectClause{Case: p.tok.Pos}

	switch {
	case p.at(token.DEFAULT):
		p.advanceToken()
	case p.at(token.CASE):
		p.advanceToken()
		x := p.parseExpr()
		if p.at(token.ASSIGN) {
			c.Targets = append(c.Targets, x)
			c.Assign = p.tok.Pos
			p.advanceToken()
			c.Op = p.parseExpr()
		} else if p.at(token.COMMA) {
			c.Targets = append(c.Targets, x)
			for p.got(token.COMMA) {
				c.Targets = append(c.Targets, p.parseExpr())
			}
			c.Assign = p.expect(token.ASSIGN)
			c.Op = p.parseExpr()
		} else {
			c.Op = x
		}
	default:
		p.errorHere(diag.ExpectedCase, p.describe(p.tok))
		p.advance(map[token.Kind]bool{token.CASE: true, token.DEFAULT: true, token.RBRACE: true})
		return c
	}

	c.Colon = p.expect(token.COLON)
	for !p.at(token.CASE) && !p.at(token.DEFAULT) && !p.at(token.RBRACE) && !p.at(token.EOF) {
		c.Body = append(c.Body, p.parseStmt())
	}
	return c
}

// parseReturnStmt parses A.5.3. A multi-value return is a bare comma list with
// no wrapping parentheses: parentheses construct a tuple, bare commas unbuild
// one.
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
	p.expectStmtEnd()
	return s
}

// parseDeferStmt parses A.5.8 ⊢ defer takes a call and nothing else.
func (p *parser) parseDeferStmt() ast.Stmt {
	pos := p.expect(token.DEFER)
	x := p.parseExpr()
	p.expectStmtEnd()

	call, ok := x.(*ast.CallExpr)
	if !ok {
		p.errorAt(diag.DeferNotCall, x.Pos())
		return &ast.BadStmt{From: pos, To: x.End()}
	}
	return &ast.DeferStmt{Defer: pos, Call: call}
}