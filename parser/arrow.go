package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
)

// Arrow functions (D.2) — the second and third speculation sites in §6.1.
//
// ArrowFormalParameters has no parenthesized cover, so telling `(a: number) =>
// a` from `(a)` means scanning to the matching `)` and past it. With a buffered
// token slice a checkpoint is an integer index and rollback is free (§4.1).

// tryArrow attempts every arrow-function head at the cursor. It returns false
// with the cursor unmoved if none applies.
func (p *parser) tryArrow(in inFlag) (ast.Expr, bool) {
	switch {
	case p.at(token.IDENT) && p.cur().Ctx == token.CtxAsync:
		if x, ok := p.tryAsyncArrow(in); ok {
			return x, true
		}
		// Fall through: `async` may just be an identifier.
		if x, ok := p.tryPlainArrow(in); ok {
			return x, true
		}
		return nil, false

	case p.at(token.IDENT):
		return p.tryPlainArrow(in)

	case p.at(token.LPAREN):
		return p.tryParenArrow(in, token.NoPos)

	case p.at(token.LT):
		// Generic arrow head (§6.1 site 3), speculated ahead of the same `(`
		// decision. This is also the reading that competes with the `< Type >`
		// assertion in B.4 — see the note in expr.go's parseUnary.
		return p.tryGenericArrow(in, token.NoPos)
	}
	return nil, false
}

// tryPlainArrow handles `x => ...`, the one head with no parentheses.
//
// The bare-identifier and parenthesized forms are not normalized into one: a
// bare head has no parentheses to point at, and inventing them would be
// synthesis (§1).
func (p *parser) tryPlainArrow(in inFlag) (ast.Expr, bool) {
	if p.peek(1).Kind != token.ARROW || p.peek(1).NLBefore() {
		return nil, false
	}
	id := p.parseIdent()
	arrow := p.next().Pos
	return &ast.ArrowFunc{Ident: id, Arrow: arrow, Body: p.parseConciseBody(in, false)}, true
}

func (p *parser) tryAsyncArrow(in inFlag) (ast.Expr, bool) {
	// `async [no LineTerminator here] AsyncArrowBindingIdentifier => ...`
	if p.peek(1).Kind == token.IDENT && !p.peek(1).NLBefore() &&
		p.peek(2).Kind == token.ARROW && !p.peek(2).NLBefore() {
		asyncPos := p.next().Pos
		id := p.parseIdent()
		arrow := p.next().Pos
		return &ast.ArrowFunc{AsyncPos: asyncPos, Ident: id, Arrow: arrow,
			Body: p.parseConciseBody(in, true)}, true
	}

	// AsyncArrowHead: `async TypeParameters_opt ( ... ) ReturnType_opt =>`
	next := p.peek(1)
	if next.NLBefore() {
		return nil, false
	}
	if next.Kind != token.LPAREN && next.Kind != token.LT {
		return nil, false
	}

	var out ast.Expr
	ok := p.speculate(func() bool {
		asyncPos := p.next().Pos
		var x ast.Expr
		var got bool
		if p.at(token.LT) {
			x, got = p.tryGenericArrow(in, asyncPos)
		} else {
			x, got = p.tryParenArrow(in, asyncPos)
		}
		if !got {
			return false
		}
		out = x
		return true
	})
	return out, ok
}

// tryParenArrow speculates ArrowFormalParameters and falls back to
// ParenthesizedExpression (§6.1 site 2).
func (p *parser) tryParenArrow(in inFlag, asyncPos token.Pos) (ast.Expr, bool) {
	// Cheap rejection first: skip the balanced group and look for `=>` or a
	// return-type annotation. This is not required for correctness — the
	// speculation below would reject too — but it keeps the common
	// `(a + b) * c` case from building and discarding a parameter list.
	if !p.arrowFollows() {
		return nil, false
	}

	var out ast.Expr
	ok := p.speculate(func() bool {
		fn := &ast.ArrowFunc{AsyncPos: asyncPos}
		fn.Params = p.parseParamList()
		if fn.Params.Rparen == token.NoPos {
			return false
		}
		if p.at(token.COLON) {
			p.next()
			fn.Result = p.parseTypeOrPredicate()
		}
		if !p.at(token.ARROW) || p.nlBefore() {
			// ArrowParameters [no LineTerminator here] => (L).
			return false
		}
		fn.Arrow = p.next().Pos
		fn.Body = p.parseConciseBody(in, asyncPos != token.NoPos)
		out = fn
		return true
	})
	return out, ok
}

func (p *parser) tryGenericArrow(in inFlag, asyncPos token.Pos) (ast.Expr, bool) {
	var out ast.Expr
	ok := p.speculate(func() bool {
		fn := &ast.ArrowFunc{AsyncPos: asyncPos}
		fn.TypeParams = p.parseTypeParams()
		if fn.TypeParams == nil || fn.TypeParams.Rangle == token.NoPos {
			return false
		}
		if !p.at(token.LPAREN) {
			return false
		}
		fn.Params = p.parseParamList()
		if fn.Params.Rparen == token.NoPos {
			return false
		}
		if p.at(token.COLON) {
			p.next()
			fn.Result = p.parseTypeOrPredicate()
		}
		if !p.at(token.ARROW) || p.nlBefore() {
			return false
		}
		fn.Arrow = p.next().Pos
		fn.Body = p.parseConciseBody(in, asyncPos != token.NoPos)
		out = fn
		return true
	})
	return out, ok
}

// arrowFollows skips the balanced parenthesized group at the cursor and reports
// whether an arrow head could follow. Pure lookahead over the immutable token
// buffer — no allocation, no diagnostics.
func (p *parser) arrowFollows() bool {
	depth := 0
	for j := 0; ; j++ {
		t := p.peek(j)
		switch t.Kind {
		case token.LPAREN:
			depth++
		case token.RPAREN:
			depth--
			if depth == 0 {
				nx := p.peek(j + 1)
				return nx.Kind == token.ARROW || nx.Kind == token.COLON
			}
		case token.EOF:
			return false
		}
	}
}

// parseConciseBody is ConciseBody / AsyncConciseBody (D.2).
func (p *parser) parseConciseBody(in inFlag, async bool) ast.Node {
	if p.at(token.LBRACE) {
		return p.parseBlock()
	}
	// [lookahead ≠ {] ExpressionBody. ExpressionBody is an
	// AssignmentExpression, so a comma at this level belongs to the enclosing
	// list rather than to the body.
	return p.parseAssign(in)
}