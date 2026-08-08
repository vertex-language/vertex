package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
)

// Statements — vertex_grammar.md section C.

// parseStmt is Statement | Declaration (C StatementListItem).
func (p *parser) parseStmt() ast.Stmt {
	defer p.trace("Statement")()
	if !p.enter() {
		return &ast.BadStmt{From: p.pos(), To: p.end()}
	}
	defer p.leave()

	switch p.kind() {
	case token.LBRACE:
		return p.parseBlock()
	case token.SEMI:
		return &ast.EmptyStmt{Semi: p.next().Pos}
	case token.VAR:
		return p.parseVarDecl()
	case token.IF:
		return p.parseIf()
	case token.DO:
		return p.parseDoWhile()
	case token.WHILE:
		return p.parseWhile()
	case token.FOR:
		return p.parseFor()
	case token.SWITCH:
		return p.parseSwitch()
	case token.CONTINUE, token.BREAK:
		return p.parseBranch()
	case token.RETURN:
		return p.parseReturn()
	case token.THROW:
		return p.parseThrow()
	case token.TRY:
		return p.parseTry()
	case token.DEBUGGER:
		t := p.next()
		return &ast.DebuggerStmt{Debugger: t.Pos, DebuggerEnd: t.End, Semi: p.expectSemi()}
	case token.FUNCTION, token.CLASS, token.CONST, token.IMPORT, token.EXPORT, token.AT:
		return p.parseDeclStmt()
	case token.IDENT:
		if d := p.tryContextualDecl(); d != nil {
			return d
		}
		// A labelled statement: `LabelIdentifier : LabelledItem` (C.4).
		if p.peek(1).Kind == token.COLON {
			label := p.parseIdent()
			colon := p.next().Pos
			return &ast.LabeledStmt{Label: label, Colon: colon, Stmt: p.parseStmt()}
		}
	}
	return p.parseExprStmt()
}

// parseExprStmt is ExpressionStatement (C).
//
// The grammar's lookahead restriction is `∉ { {, function, async function,
// class, let [ }`. Every one of those is dispatched above, so reaching here
// means the restriction is satisfied — except for `let [`, which is checked
// here because `let` is an IDENT (see ctx.go) and could be an expression.
func (p *parser) parseExprStmt() ast.Stmt {
	if p.atCtx(token.CtxLet) && p.peek(1).Kind == token.LBRACK {
		p.errorf(p.cur(), "an expression statement cannot begin with `let [`")
	}
	x := p.parseExpr(allowIn)
	return &ast.ExprStmt{X: x, Semi: p.expectSemi()}
}

func (p *parser) parseBlock() *ast.BlockStmt {
	b := &ast.BlockStmt{Lbrace: p.expect(token.LBRACE)}
	b.List = p.parseStmtList(token.RBRACE)
	b.Rbrace = p.expect(token.RBRACE)
	return b
}

func (p *parser) parseStmtList(stop token.Kind) []ast.Stmt {
	var out []ast.Stmt
	for !p.at(stop) && !p.atEOF() {
		before := p.i
		s := p.parseStmt()
		if s != nil {
			out = append(out, s)
		}
		if !p.advanced(before) {
			if !p.advanceTo(syncStmt) {
				break
			}
		}
	}
	return out
}

func (p *parser) parseIf() ast.Stmt {
	s := &ast.IfStmt{If: p.next().Pos}
	s.Lparen = p.expect(token.LPAREN)
	s.Cond = p.parseExpr(allowIn)
	s.Rparen = p.expect(token.RPAREN)
	s.Body = p.parseStmt()
	// `if (...) Statement [lookahead ≠ else]` versus the else form (C.3).
	if p.at(token.ELSE) {
		s.ElsePos = p.next().Pos
		s.Else = p.parseStmt()
	}
	return s
}

func (p *parser) parseDoWhile() ast.Stmt {
	s := &ast.DoWhileStmt{Do: p.next().Pos}
	s.Body = p.parseStmt()
	s.While = p.expect(token.WHILE)
	p.expect(token.LPAREN)
	s.Cond = p.parseExpr(allowIn)
	s.Rparen = p.expect(token.RPAREN)
	// DoWhileStatement is in TerminatedByASI (L).
	s.Semi = p.expectSemi()
	return s
}

func (p *parser) parseWhile() ast.Stmt {
	s := &ast.WhileStmt{While: p.next().Pos}
	p.expect(token.LPAREN)
	s.Cond = p.parseExpr(allowIn)
	s.Rparen = p.expect(token.RPAREN)
	s.Body = p.parseStmt()
	return s
}

// parseFor handles ForStatement and ForInOfStatement (C.3), which share a head
// up to the token after the initializer.
func (p *parser) parseFor() ast.Stmt {
	forPos := p.next().Pos

	awaitPos := token.NoPos
	if p.at(token.AWAIT) {
		awaitPos = p.next().Pos
	}
	p.expect(token.LPAREN)

	// Empty initializer: only the three-clause form.
	if p.at(token.SEMI) {
		return p.finishForStmt(forPos, nil)
	}

	var left ast.Node
	if p.atForDeclaration() {
		// ForDeclaration and `var ForBinding` (C.3). Parsed with [~In] so that
		// `in` is the loop keyword rather than an operator.
		left = p.parseForDeclaration()
	} else {
		left = p.parseExpr(noIn)
	}

	switch {
	case p.at(token.IN) && awaitPos == token.NoPos:
		s := &ast.ForInStmt{For: forPos, Left: left, In: p.next().Pos}
		s.Right = p.parseExpr(allowIn)
		s.Rparen = p.expect(token.RPAREN)
		s.Body = p.parseStmt()
		return s

	case p.atCtx(token.CtxOf):
		s := &ast.ForOfStmt{For: forPos, AwaitPos: awaitPos, Left: left, OfPos: p.next().Pos}
		s.Right = p.parseAssign(allowIn)
		s.Rparen = p.expect(token.RPAREN)
		s.Body = p.parseStmt()
		return s
	}

	if awaitPos != token.NoPos {
		p.errorAt(awaitPos, awaitPos+5, "`for await` requires an `of` clause")
	}
	return p.finishForStmt(forPos, left)
}

// finishForStmt completes the three-clause form.
//
// The initializer already consumed its own `;` when it was a
// LexicalDeclaration or VariableStatement — C.3's third production writes only
// two semicolons for that reason.
func (p *parser) finishForStmt(forPos token.Pos, init ast.Node) ast.Stmt {
	s := &ast.ForStmt{For: forPos}
	switch init := init.(type) {
	case nil:
		p.expect(token.SEMI)
	case *ast.VarDecl:
		s.Init = init
		if init.Semi == token.NoPos {
			p.expect(token.SEMI)
		}
	case ast.Expr:
		s.Init = &ast.ExprStmt{X: init, Semi: p.expect(token.SEMI)}
	}
	if !p.at(token.SEMI) {
		s.Cond = p.parseExpr(allowIn)
	}
	p.expect(token.SEMI)
	if !p.at(token.RPAREN) {
		s.Post = p.parseExpr(allowIn)
	}
	s.Rparen = p.expect(token.RPAREN)
	s.Body = p.parseStmt()
	return s
}

func (p *parser) atForDeclaration() bool {
	switch p.kind() {
	case token.VAR, token.CONST:
		return true
	case token.AWAIT:
		return p.peek(1).Kind == token.IDENT && p.peek(1).Ctx == token.CtxUsing
	case token.IDENT:
		switch p.cur().Ctx {
		case token.CtxLet:
			// `let` is only a declaration keyword when a binding follows.
			switch p.peek(1).Kind {
			case token.IDENT, token.LBRACK, token.LBRACE:
				return true
			}
		case token.CtxUsing:
			// `using [no LineTerminator here] [lookahead ≠ await]` (C.3).
			//
			// The lookahead restriction needs no explicit check against
			// `await`: `await` is a ReservedWord (token.AWAIT), never an
			// IDENT, so the Kind == IDENT test below already excludes it.
			// There is no token.CtxAwait — `await` never scans as a
			// contextual identifier, so no Ctx value for it exists.
			return !p.peek(1).NLBefore() && p.peek(1).Kind == token.IDENT
		}
	}
	return false
}

func (p *parser) parseForDeclaration() ast.Node {
	d := &ast.VarDecl{}
	switch {
	case p.at(token.VAR):
		d.Kind, d.KindPos = ast.VarVar, p.next().Pos
	case p.at(token.CONST):
		d.Kind, d.KindPos = ast.VarConst, p.next().Pos
	case p.at(token.AWAIT):
		d.AwaitPos = p.next().Pos
		p.noLineTerminator("`using`")
		d.Kind, d.KindPos = ast.VarAwaitUsing, p.expectCtx(token.CtxUsing)
	case p.atCtx(token.CtxUsing):
		d.Kind, d.KindPos = ast.VarUsing, p.next().Pos
	default:
		d.Kind, d.KindPos = ast.VarLet, p.next().Pos
	}
	d.List = append(d.List, p.parseBinding(noIn, false))
	// No `;` here: the caller decides whether this was a for-in/of head.
	d.Semi = token.NoPos
	return d
}

func (p *parser) parseSwitch() ast.Stmt {
	s := &ast.SwitchStmt{Switch: p.next().Pos}
	p.expect(token.LPAREN)
	s.Tag = p.parseExpr(allowIn)
	p.expect(token.RPAREN)
	s.Lbrace = p.expect(token.LBRACE)

	for !p.at(token.RBRACE) && !p.atEOF() {
		before := p.i
		c := &ast.CaseClause{}
		switch {
		case p.at(token.CASE):
			c.Case = p.next().Pos
			c.Cond = p.parseExpr(allowIn)
		case p.at(token.DEFAULT):
			// More than one default parses and is rejected by name later
			// (§6.3), so the message can point at the duplicate.
			c.Case = p.next().Pos
		default:
			p.errorf(p.cur(), "expected `case` or `default`, found %s", p.describe(p.cur()))
			if !p.advanceTo(syncStmt) {
				break
			}
			continue
		}
		c.Colon = p.expect(token.COLON)
		for !p.at(token.CASE) && !p.at(token.DEFAULT) && !p.at(token.RBRACE) && !p.atEOF() {
			inner := p.i
			c.Body = append(c.Body, p.parseStmt())
			if !p.advanced(inner) {
				if !p.advanceTo(syncStmt) {
					break
				}
			}
		}
		s.Cases = append(s.Cases, c)
		p.advanced(before)
	}
	s.Rbrace = p.expect(token.RBRACE)
	return s
}

func (p *parser) parseBranch() ast.Stmt {
	t := p.next()
	s := &ast.BranchStmt{TokPos: t.Pos, TokEnd: t.End, Tok: t.Kind}
	// `continue / break [no LineTerminator here] LabelIdentifier` (C.4, L).
	if p.at(token.IDENT) && !p.nlBefore() {
		s.Label = p.parseIdent()
	}
	s.Semi = p.expectSemi()
	return s
}

func (p *parser) parseReturn() ast.Stmt {
	t := p.next()
	s := &ast.ReturnStmt{Return: t.Pos, ReturnEnd: t.End}
	// `return [no LineTerminator here] Expression` (C.4, L). A break makes it
	// a bare return, which is the classic ASI hazard and why the restriction
	// exists.
	if !p.nlBefore() && p.startsExpr() {
		s.Result = p.parseExpr(allowIn)
	}
	s.Semi = p.expectSemi()
	return s
}

func (p *parser) parseThrow() ast.Stmt {
	throwPos := p.next().Pos
	// `throw [no LineTerminator here] Expression` — unlike return, the
	// expression is required, so a break is an error rather than a shorter
	// statement.
	if p.nlBefore() {
		p.errorf(p.cur(), "no line break allowed between `throw` and its expression")
	}
	s := &ast.ThrowStmt{Throw: throwPos, X: p.parseExpr(allowIn)}
	s.Semi = p.expectSemi()
	return s
}

func (p *parser) parseTry() ast.Stmt {
	s := &ast.TryStmt{Try: p.next().Pos}
	s.Body = p.parseBlock()

	if p.at(token.CATCH) {
		c := &ast.CatchClause{Catch: p.next().Pos}
		if p.at(token.LPAREN) {
			c.Lparen = p.next().Pos
			c.Param = p.parseBindingTarget()
			if p.at(token.COLON) {
				// CatchTypeAnnotation is `: any` or `: unknown` only (C.4).
				// Anything else parses and is rejected by name, so the message
				// can say which annotations are allowed.
				p.next()
				c.CatchType = p.parseType()
				p.checkCatchType(c.CatchType)
			}
			c.Rparen = p.expect(token.RPAREN)
		}
		c.Body = p.parseBlock()
		s.Catch = c
	}
	if p.at(token.FINALLY) {
		f := &ast.FinallyClause{Finally: p.next().Pos}
		f.Body = p.parseBlock()
		s.Finally = f
	}
	// try with neither catch nor finally parses, per §6.3.
	if s.Catch == nil && s.Finally == nil {
		p.errorAt(s.Try, s.Body.End(), "`try` requires a `catch` or `finally` clause")
	}
	return s
}

func (p *parser) checkCatchType(t ast.TypeExpr) {
	pt, ok := ast.UnparenType(t).(*ast.PredefinedType)
	if ok && (pt.Ctx == token.CtxAny || pt.Ctx == token.CtxUnknown) {
		return
	}
	p.errorAt(t.Pos(), t.End(), "a catch parameter may only be annotated `any` or `unknown`")
}