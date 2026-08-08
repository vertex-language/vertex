package scanner

import "github.com/vertex-language/vertex/token"

// One stack serves three purposes that §4.4 and §4.5 describe separately but
// that cannot be separated: each `(` records whether it followed a control
// keyword, each `{` records block vs. object literal, and each `${` opens a
// template substitution. They share a stack because the `}` in `` `a${b}c` ``
// must not be read as closing a block, and that decision needs the same
// depth counter as the other two.
type frameKind uint8

const (
	frameParen        frameKind = iota // (
	frameControlParen                  // ( following if / while / for / switch / catch
	frameBrack                         // [
	frameBlock                         // { opening a Block
	frameObject                        // { opening an ObjectLiteral
	frameSubst                         // ${ inside a template
)

type frame struct {
	kind frameKind
	pos  int // opener offset, for the unmatched-bracket diagnostic
}

func (s *scanner) push(k frameKind, pos int) { s.stack = append(s.stack, frame{k, pos}) }

func (s *scanner) top() (frame, bool) {
	if len(s.stack) == 0 {
		return frame{}, false
	}
	return s.stack[len(s.stack)-1], true
}

// popIf removes the top frame when it is one of want. A mismatch leaves the
// stack alone: the closer is still emitted, the parser's recovery (§6.4)
// handles the structure, and §4.4's "one recoverable diagnostic, never a
// cascade" is preserved because the stack does not unwind past a real opener.
func (s *scanner) popIf(want ...frameKind) (frame, bool) {
	f, ok := s.top()
	if !ok {
		return frame{}, false
	}
	for _, w := range want {
		if f.kind == w {
			s.stack = s.stack[:len(s.stack)-1]
			return f, true
		}
	}
	return frame{}, false
}

// controlParen reports whether a `(` at this position follows one of the
// keywords whose parenthesized head is a control-flow head rather than a call
// or a parameter list. It is the difference between `if (x) /re/` (regex) and
// `f(x) / 2` (division).
func (s *scanner) controlParen() bool {
	if !s.hasPrev {
		return false
	}
	switch s.prev.Kind {
	case token.IF, token.WHILE, token.FOR, token.SWITCH, token.CATCH, token.WITH:
		return true
	}
	return false
}

// regexAllowed implements the previous-significant-token rule (§4.4).
//
// The question is whether a value has just been produced. After a value, `/`
// divides; otherwise it opens a RegularExpressionLiteral.
func (s *scanner) regexAllowed() bool {
	if !s.hasPrev {
		return true
	}
	switch k := s.prev.Kind; k {
	// A template that has just *opened* a substitution leaves us in
	// expression position: `` `a${/re/}` `` is a regex.
	case token.TEMPLATE_HEAD, token.TEMPLATE_MIDDLE:
		return true
	case token.TEMPLATE, token.TEMPLATE_TAIL:
		return false

	case token.RPAREN:
		return s.prevParenWasControl
	case token.RBRACE:
		return s.prevBraceWasBlock
	case token.RBRACK:
		return false

	// Postfix update produced a value.
	case token.INC, token.DEC:
		return false

	// Reserved words that denote values.
	case token.THIS, token.SUPER, token.TRUE, token.FALSE, token.NULL:
		return false

	default:
		switch {
		case k.IsReserved():
			return true // return /re/, typeof /re/, case /re/, ...
		case k.IsLiteral():
			return false // IDENT is in this range, and so are all literals
		case k.IsOperator():
			return true
		}
	}
	return true
}

// braceOpensBlock classifies a `{`.
//
// This is a heuristic and is documented as one: §4.4 accepts that
// misclassification costs one recoverable diagnostic. The blast radius is
// narrow — the classification is read back only when a `/` immediately
// follows the matching `}`, so getting it wrong turns one regex into one
// division and nothing else.
func (s *scanner) braceOpensBlock() bool {
	if !s.hasPrev {
		return true
	}
	switch k := s.prev.Kind; k {
	// Statement-position closers and separators.
	case token.SEMI, token.LBRACE, token.RBRACE, token.RPAREN, token.ARROW, token.RBRACK:
		return true
	case token.ELSE, token.DO, token.TRY, token.FINALLY:
		return true

	case token.COLON:
		// `{a: {b: 1}}` is an object; `case x: {` and `label: {` are blocks.
		// The enclosing frame decides: a colon inside an object literal or an
		// argument list is a value position.
		if f, ok := s.top(); ok {
			return !(f.kind == frameObject || f.kind == frameParen || f.kind == frameBrack)
		}
		return true

	// Value positions: an operator or a keyword expecting an operand.
	case token.RETURN, token.TYPEOF, token.VOID, token.DELETE, token.NEW,
		token.IN, token.INSTANCEOF, token.CASE, token.YIELD, token.AWAIT,
		token.DEFAULT, token.EXTENDS:
		return false

	default:
		if k.IsOperator() {
			return false // ( [ , = ? && ... all expect a value
		}
		return true
	}
}