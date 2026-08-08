package scanner

import "github.com/vertex-language/vertex/token"

// scanPunct handles A.3 Punctuator, DivPunctuator, and RightBracePunctuator.
//
// Maximal munch everywhere except `>`, which is always emitted as a single
// character (§4.2). See the flagged note about `<`.
func (s *scanner) scanPunct(start int) {
	c := s.src[s.off]
	s.off++

	var kind token.Kind
	switch c {
	case '{':
		block := s.braceOpensBlock()
		if block {
			s.push(frameBlock, start)
		} else {
			s.push(frameObject, start)
		}
		s.emit(token.LBRACE, token.CtxNone, start, 0)
		return

	case '}':
		// A `}` that closes a template substitution resumes template scanning
		// rather than closing a block (§4.5). The frame decides.
		if f, ok := s.top(); ok && f.kind == frameSubst {
			s.stack = s.stack[:len(s.stack)-1]
			kind, flags := s.scanTemplateBody(start)
			s.emit(kind, token.CtxNone, start, flags)
			return
		}
		f, ok := s.popIf(frameBlock, frameObject)
		s.prevBraceWasBlock = !ok || f.kind == frameBlock
		s.emit(token.RBRACE, token.CtxNone, start, 0)
		return

	case '(':
		if s.controlParen() {
			s.push(frameControlParen, start)
		} else {
			s.push(frameParen, start)
		}
		s.emit(token.LPAREN, token.CtxNone, start, 0)
		return

	case ')':
		f, ok := s.popIf(frameParen, frameControlParen)
		s.prevParenWasControl = ok && f.kind == frameControlParen
		s.emit(token.RPAREN, token.CtxNone, start, 0)
		return

	case '[':
		s.push(frameBrack, start)
		s.emit(token.LBRACK, token.CtxNone, start, 0)
		return

	case ']':
		s.popIf(frameBrack)
		s.emit(token.RBRACK, token.CtxNone, start, 0)
		return

	case ';':
		kind = token.SEMI
	case ',':
		kind = token.COMMA
	case '@':
		kind = token.AT
	case '~':
		kind = token.TILDE

	case '.':
		if s.has("..") {
			s.off += 2
			kind = token.ELLIPSIS
		} else {
			kind = token.PERIOD
		}

	case ':':
		kind = token.COLON

	case '?':
		switch {
		case s.at(0) == '?':
			s.off++
			if s.at(0) == '=' {
				s.off++
				kind = token.COALESCE_ASSIGN
			} else {
				kind = token.COALESCE
			}
		case s.at(0) == '.' && !isDecimalDigit(s.at(1)):
			// A.3: `?.` requires lookahead ∉ DecimalDigit, so `a?.5:b` stays
			// a ConditionalExpression.
			s.off++
			kind = token.QUESTION_DOT
		default:
			kind = token.QUESTION
		}

	case '=':
		switch {
		case s.at(0) == '=' && s.at(1) == '=':
			s.off += 2
			kind = token.STRICT_EQL
		case s.at(0) == '=':
			s.off++
			kind = token.EQL
		case s.at(0) == '>':
			s.off++
			kind = token.ARROW
		default:
			kind = token.ASSIGN
		}

	case '!':
		switch {
		case s.at(0) == '=' && s.at(1) == '=':
			s.off += 2
			kind = token.STRICT_NEQ
		case s.at(0) == '=':
			s.off++
			kind = token.NEQ
		default:
			kind = token.NOT
		}

	case '+':
		kind = s.oneOf(token.ADD, '+', token.INC, '=', token.ADD_ASSIGN)
	case '-':
		kind = s.oneOf(token.SUB, '-', token.DEC, '=', token.SUB_ASSIGN)
	case '&':
		kind = s.logical(token.AND, token.LAND, token.AND_ASSIGN, token.LAND_ASSIGN, '&')
	case '|':
		kind = s.logical(token.OR, token.LOR, token.OR_ASSIGN, token.LOR_ASSIGN, '|')

	case '^':
		if s.at(0) == '=' {
			s.off++
			kind = token.XOR_ASSIGN
		} else {
			kind = token.XOR
		}
	case '%':
		if s.at(0) == '=' {
			s.off++
			kind = token.REM_ASSIGN
		} else {
			kind = token.REM
		}
	case '/':
		if s.at(0) == '=' {
			s.off++
			kind = token.QUO_ASSIGN
		} else {
			kind = token.QUO
		}

	case '*':
		switch {
		case s.at(0) == '*' && s.at(1) == '=':
			s.off += 2
			kind = token.EXP_ASSIGN
		case s.at(0) == '*':
			s.off++
			kind = token.EXP
		case s.at(0) == '=':
			s.off++
			kind = token.MUL_ASSIGN
		default:
			kind = token.MUL
		}

	case '<':
		// MERGED — see the flagged note. `Foo<<T>() => void>` is grammatical
		// (G.4 FunctionType inside G.2 TypeArguments) and tokenizes wrong here.
		switch {
		case s.at(0) == '<' && s.at(1) == '=':
			s.off += 2
			kind = token.SHL_ASSIGN
		case s.at(0) == '<':
			s.off++
			kind = token.SHL
		case s.at(0) == '=':
			s.off++
			kind = token.LEQ
		default:
			kind = token.LT
		}

	case '>':
		// §4.2: never merged. `Array<FixedArray<int32, 4>>` and
		// `let a: Box<T>= init` both need this to come apart in type context,
		// and splitting later would mean mutating the buffer mid-speculation.
		// token.JoinGT reassembles the run in the expression parser.
		kind = token.GT

	default:
		s.error(start, s.off, "unexpected character")
		kind = token.INVALID
	}

	s.emit(kind, token.CtxNone, start, 0)
}

func (s *scanner) oneOf(base token.Kind, c2 byte, k2 token.Kind, c3 byte, k3 token.Kind) token.Kind {
	switch s.at(0) {
	case c2:
		s.off++
		return k2
	case c3:
		s.off++
		return k3
	}
	return base
}

func (s *scanner) logical(base, dbl, asgn, dblAsgn token.Kind, c byte) token.Kind {
	switch {
	case s.at(0) == c && s.at(1) == '=':
		s.off += 2
		return dblAsgn
	case s.at(0) == c:
		s.off++
		return dbl
	case s.at(0) == '=':
		s.off++
		return asgn
	}
	return base
}

// at returns the byte n positions ahead of the cursor, or 0 past the end. Zero
// is safe as a sentinel: a NUL byte in source is rejected as an unexpected
// character before any of these lookaheads run.
func (s *scanner) at(n int) byte {
	if s.off+n < len(s.src) {
		return s.src[s.off+n]
	}
	return 0
}

func (s *scanner) has(lit string) bool {
	if s.off+len(lit) > len(s.src) {
		return false
	}
	return string(s.src[s.off:s.off+len(lit)]) == lit
}