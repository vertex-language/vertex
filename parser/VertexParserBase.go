package parser

import "github.com/antlr4-go/antlr/v4"

// VertexParserBase is the superclass for the ANTLR4-generated VertexParser.
// The grammar declares `superClass = VertexParserBase`; ANTLR embeds this
// struct into the generated parser so every predicate method below is
// accessible as p.isBinaryOp() etc. inside grammar actions.
type VertexParserBase struct {
	*antlr.BaseParser
}

// ── Operator classification ───────────────────────────────────────────────────
//
// Swift (and Vertex) classify operators by the whitespace surrounding them:
//
//   left WS only  →  prefix   (e.g.  `return -x`)
//   right WS only →  postfix  (rare)
//   both or neither → binary  (e.g.  `a + b`  or  `a+b`)
//
// Without these predicates, ANTLR4's SLL prediction resolves the ambiguity
// between binaryOp and a new prefixOp statement in the wrong direction,
// splitting  `fibonacci(n-1) + fibonacci(n-2)`  into two statements.

// isBinaryOp returns true when the upcoming OPERATOR token should be treated
// as a binary (infix) operator.  Guards the `binaryOp` alternative in
// binaryExpression.
func (p *VertexParserBase) isBinaryOp() bool {
	left, right := p.operatorWS()
	// Binary: symmetric whitespace (both sides or neither side).
	return left == right
}

// isPrefixOp returns true when the upcoming OPERATOR token is a prefix
// operator.  Guards the `prefixOp` alternative in prefixExpression.
func (p *VertexParserBase) isPrefixOp() bool {
	left, right := p.operatorWS()
	// Prefix: whitespace on left but NOT right.
	// Also treat "no left token at all" (beginning of stream / block) as prefix.
	return left && !right
}

// isPostfixOp returns true when the upcoming OPERATOR token is a postfix
// operator (whitespace on right but not left).
func (p *VertexParserBase) isPostfixOp() bool {
	left, right := p.operatorWS()
	return !left && right
}

// operatorWS inspects the hidden-channel tokens immediately surrounding the
// upcoming OPERATOR token and reports whether there is whitespace on each
// side.  Uses CommonTokenStream.GetHiddenTokensToLeft / Right which return
// every off-default-channel token between the query index and the nearest
// on-channel token in that direction.
func (p *VertexParserBase) operatorWS() (leftWS, rightWS bool) {
	cts, ok := p.GetTokenStream().(*antlr.CommonTokenStream)
	if !ok {
		// Fallback: if we cannot inspect the stream, assume binary so we
		// don't silently drop parts of expressions.
		return true, true
	}

	tok := p.GetCurrentToken()
	if tok == nil {
		return true, true
	}
	idx := tok.GetTokenIndex()

	// Fill ensures all tokens up to (and slightly past) the current position
	// are loaded so the hidden-token queries below are reliable.
	cts.Fill()

	leftWS  = len(cts.GetHiddenTokensToLeft(idx, -1)) > 0
	rightWS = len(cts.GetHiddenTokensToRight(idx, -1)) > 0

	// Edge case: idx == 0 means the operator is literally the very first
	// token in the file — treat as prefix (no left context at all).
	if idx == 0 {
		leftWS = false
	}
	return
}