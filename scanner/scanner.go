// Package scanner converts Vertex source bytes into tokens.
//
// The scanner emits every line terminator and interprets none. A line
// terminator is a terminator if and only if the innermost enclosing bracketing
// construct is one whose production writes `terminator`, and whether a given
// `{` opens a Block or a LiteralValue is a parsing question. So line structure
// is only ever recorded — on token.Token.NLBefore — and the parser decides what
// it means. There is no SEMICOLON kind, no synthesized NEWLINE, and no
// automatic semicolon insertion.
//
// A run of consecutive line terminators is one terminator, so one bool carries
// the whole rule: no consumer ever needs the count.
//
// A line terminator inside a raw string is part of that string's value and must
// not end a statement, though it does advance the line table. Keeping those two
// effects separate is why consumeLineTerminator and skipSpace are separate
// calls below.
//
// The scanner recognizes keywords and the three reserved literal keywords
// through token.Lookup, and nothing else. Contextual keywords, predeclared type
// and tensor-element and constraint names, and reserved builtin names are all
// ordinary identifiers here; they are pre-bound in an implicit outermost scope
// that the scanner must not know about.
package scanner

import (
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
)

// Mode controls optional scanner behavior.
type Mode uint

const (
	// ScanComments returns COMMENT tokens instead of discarding them. A
	// comment participates in line structure either way; the mode controls
	// only whether the token is surfaced.
	ScanComments Mode = 1 << iota
)

const eof = -1

const (
	lineSep = 0x2028 // <LS>
	paraSep = 0x2029 // <PS>
	nbsp    = 0x00A0
	bom     = 0xFEFF // ordinary white space wherever it appears, leading BOM included
)

type Scanner struct {
	file *token.File
	src  []byte
	rep  diag.Reporter
	mode Mode

	ch       rune // current character; eof at end
	offset   int  // byte offset of ch
	rdOffset int  // byte offset just past ch

	// nlPending records that a line terminator has been seen since the last
	// token was produced. It is consumed by the next token's NLBefore.
	nlPending bool

	// prev and prev2 are the kinds of the last two non-comment tokens. Their
	// only use is the one restriction to longest-match scanning: a float_lit
	// may not begin immediately after a `.` whose own preceding token can end
	// an operand. Both halves of that condition are needed, which is why two
	// kinds are kept rather than one.
	prev  token.Kind
	prev2 token.Kind

	ErrorCount int
}

// Init prepares s to scan src, which must be the contents of file.
func (s *Scanner) Init(file *token.File, src []byte, rep diag.Reporter, mode Mode) {
	if file.Size() != len(src) {
		panic("scanner: file size does not match src length")
	}
	s.file = file
	s.src = src
	s.rep = rep
	s.mode = mode

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.nlPending = false
	s.prev = token.INVALID
	s.prev2 = token.INVALID
	s.ErrorCount = 0

	s.next()
}

func (s *Scanner) error(code diag.Code, offset int, args ...any) {
	s.ErrorCount++
	if s.rep != nil {
		s.rep.Report(diag.At(code, s.file.Pos(offset), args...))
	}
}

func (s *Scanner) errorSpan(code diag.Code, from, to int, args ...any) {
	s.ErrorCount++
	if s.rep != nil {
		s.rep.Report(diag.New(code, s.file.Pos(from), s.file.Pos(to), args...))
	}
}

// next advances to the following character. It does not touch the line table:
// line terminators are consumed explicitly by consumeLineTerminator, because
// <CR><LF> is one terminator over two bytes and <LS>/<PS> are three each.
//
// NUL is reported here rather than at the dispatch below, so that one appearing
// inside a string or comment is still rejected — a compiler must reject the NUL
// character in source text, not merely in code.
func (s *Scanner) next() {
	if s.rdOffset >= len(s.src) {
		s.offset = len(s.src)
		s.ch = eof
		return
	}
	s.offset = s.rdOffset
	if b := s.src[s.rdOffset]; b < utf8.RuneSelf {
		s.rdOffset++
		if b == 0 {
			s.error(diag.NulCharacter, s.offset)
		}
		s.ch = rune(b)
		return
	}
	r, w := utf8.DecodeRune(s.src[s.rdOffset:])
	s.rdOffset += w
	if r == utf8.RuneError && w == 1 {
		s.error(diag.IllegalUTF8, s.offset)
	}
	s.ch = r
}

// peek returns the byte after the current character without advancing. Used
// only for ASCII lookahead, where a byte and a rune coincide.
func (s *Scanner) peek() byte {
	if s.rdOffset < len(s.src) {
		return s.src[s.rdOffset]
	}
	return 0
}

func isLineTerminator(ch rune) bool {
	return ch == '\n' || ch == '\r' || ch == lineSep || ch == paraSep
}

// isWhiteSpace matches the white_space set: tab, VT, FF, space, NBSP, ZWNBSP,
// and any code point in category Zs.
func isWhiteSpace(ch rune) bool {
	switch ch {
	case '\t', '\v', '\f', ' ', nbsp, bom:
		return true
	}
	return ch >= utf8.RuneSelf && unicode.Is(unicode.Zs, ch)
}

// consumeLineTerminator consumes one full line_terminator sequence and records
// the new line. It deliberately does not set nlPending: a terminator inside a
// raw string advances the line table without ending a statement.
func (s *Scanner) consumeLineTerminator() {
	if s.ch == '\r' {
		s.next()
		if s.ch == '\n' {
			s.next()
		}
	} else {
		s.next()
	}
	s.file.AddLine(s.offset)
}

func (s *Scanner) skipSpace() {
	for {
		switch {
		case isLineTerminator(s.ch):
			s.nlPending = true
			s.consumeLineTerminator()
		case isWhiteSpace(s.ch):
			s.next()
		default:
			return
		}
	}
}

// token builds a token carrying the pending line-break flag and clears it.
//
// COMMENT does not update the kind history, so a comment sitting between a `.`
// and a digit does not disturb the selector-dot restriction.
func (s *Scanner) token(kind token.Kind, pos token.Pos, lit string) token.Token {
	t := token.Token{Kind: kind, Pos: pos, Lit: lit, NLBefore: s.nlPending}
	s.nlPending = false
	if kind != token.COMMENT {
		s.prev2, s.prev = s.prev, kind
	}
	return t
}

// afterSelectorDot reports whether a numeric literal beginning here falls under
// the one restriction to longest-match scanning: the previous token is a `.`
// and the one before it can end an operand. Where it holds, the scanner
// produces an int_lit and any `.` inside the would-be float scans separately,
// which is what makes a chain like `t.0.0` fall out.
func (s *Scanner) afterSelectorDot() bool {
	return s.prev == token.PERIOD && s.prev2.EndsOperand()
}

// Scan returns the next token. At end of input it returns EOF indefinitely,
// with NLBefore reflecting any trailing line terminator.
func (s *Scanner) Scan() token.Token {
	for {
		s.skipSpace()

		offs := s.offset
		pos := s.file.Pos(offs)
		ch := s.ch

		switch {
		case ch == eof:
			return s.token(token.EOF, pos, "")

		case ch == 0:
			// Already reported by next(); consumed here so one NUL does not
			// stall the scan.
			s.next()
			return s.token(token.INVALID, pos, string(s.src[offs:s.offset]))

		case isIdentStart(ch) || ch == '$':
			lit := s.scanIdent()
			// Lookup maps the reserved keywords and true/false/nil. Everything
			// else is an ordinary identifier and keeps its spelling.
			kind := token.Lookup(lit)
			if kind != token.IDENT {
				lit = ""
			}
			return s.token(kind, pos, lit)

		case '0' <= ch && ch <= '9':
			if s.afterSelectorDot() {
				return s.token(token.INT, pos, s.scanTupleIndex())
			}
			kind, lit := s.scanNumber()
			return s.token(kind, pos, lit)

		case ch == '"':
			return s.token(token.STRING, pos, s.scanString())

		case ch == '`':
			return s.token(token.STRING, pos, s.scanRawString())

		case ch == '\'':
			return s.token(token.CHAR, pos, s.scanChar())
		}

		// Operators, punctuation, and comments. The longest match wins.
		switch ch {
		case '/':
			if p := s.peek(); p == '/' || p == '*' {
				lit, actsAsTerminator := s.scanComment()
				t := s.token(token.COMMENT, pos, lit)
				if actsAsTerminator {
					// A general comment containing no line terminator is white
					// space; every other comment, line comments included, acts
					// like a line terminator.
					s.nlPending = true
				}
				if s.mode&ScanComments != 0 {
					return t
				}
				continue
			}
			return s.token(s.munch2('=', token.QUO_ASSIGN, token.QUO), pos, "")

		case '(':
			s.next()
			return s.token(token.LPAREN, pos, "")
		case ')':
			s.next()
			return s.token(token.RPAREN, pos, "")
		case '[':
			s.next()
			return s.token(token.LBRACK, pos, "")
		case ']':
			s.next()
			return s.token(token.RBRACK, pos, "")
		case '{':
			s.next()
			return s.token(token.LBRACE, pos, "")
		case '}':
			s.next()
			return s.token(token.RBRACE, pos, "")
		case ',':
			s.next()
			return s.token(token.COMMA, pos, "")
		case ':':
			s.next()
			return s.token(token.COLON, pos, "")
		case '~':
			s.next()
			return s.token(token.TILDE, pos, "")

		case '.':
			s.next()
			if s.ch == '.' {
				s.next()
				if s.ch == '.' {
					s.next()
					return s.token(token.ELLIPSIS, pos, "")
				}
				return s.token(token.DOTDOT, pos, "")
			}
			return s.token(token.PERIOD, pos, "")

		case '+':
			return s.token(s.munch2('=', token.ADD_ASSIGN, token.ADD), pos, "")
		case '*':
			return s.token(s.munch2('=', token.MUL_ASSIGN, token.MUL), pos, "")
		case '%':
			return s.token(s.munch2('=', token.REM_ASSIGN, token.REM), pos, "")
		case '^':
			return s.token(s.munch2('=', token.XOR_ASSIGN, token.XOR), pos, "")

		case '-':
			s.next()
			switch s.ch {
			case '>':
				s.next()
				return s.token(token.ARROW, pos, "")
			case '=':
				s.next()
				return s.token(token.SUB_ASSIGN, pos, "")
			}
			return s.token(token.SUB, pos, "")

		case '&':
			s.next()
			switch s.ch {
			case '&':
				s.next()
				return s.token(token.LAND, pos, "")
			case '+':
				s.next()
				return s.token(token.WRAP_ADD, pos, "")
			case '-':
				s.next()
				return s.token(token.WRAP_SUB, pos, "")
			case '*':
				s.next()
				return s.token(token.WRAP_MUL, pos, "")
			case '=':
				s.next()
				return s.token(token.AND_ASSIGN, pos, "")
			}
			return s.token(token.AND, pos, "")

		case '|':
			s.next()
			switch s.ch {
			case '|':
				s.next()
				return s.token(token.LOR, pos, "")
			case '=':
				s.next()
				return s.token(token.OR_ASSIGN, pos, "")
			}
			return s.token(token.OR, pos, "")

		case '<':
			s.next()
			switch s.ch {
			case '<':
				s.next()
				if s.ch == '=' {
					s.next()
					return s.token(token.SHL_ASSIGN, pos, "")
				}
				return s.token(token.SHL, pos, "")
			case '=':
				s.next()
				return s.token(token.LEQ, pos, "")
			}
			return s.token(token.LSS, pos, "")

		case '>':
			s.next()
			switch s.ch {
			case '>':
				s.next()
				if s.ch == '=' {
					s.next()
					return s.token(token.SHR_ASSIGN, pos, "")
				}
				return s.token(token.SHR, pos, "")
			case '=':
				s.next()
				return s.token(token.GEQ, pos, "")
			}
			return s.token(token.GTR, pos, "")

		case '=':
			s.next()
			if s.ch == '=' {
				s.next()
				if s.ch == '=' {
					s.next()
					return s.token(token.IDENTICAL, pos, "")
				}
				return s.token(token.EQL, pos, "")
			}
			return s.token(token.ASSIGN, pos, "")

		case '!':
			s.next()
			if s.ch == '=' {
				s.next()
				if s.ch == '=' {
					s.next()
					return s.token(token.NOT_IDENTICAL, pos, "")
				}
				return s.token(token.NEQ, pos, "")
			}
			return s.token(token.NOT, pos, "")
		}

		// Nothing matched. Report, consume, and keep going, so one stray
		// character does not abort the file.
		s.next()
		s.error(diag.IllegalCharacter, offs, strconv.QuoteRune(ch))
		return s.token(token.INVALID, pos, string(s.src[offs:s.offset]))
	}
}

// munch2 consumes the current character and, if the next one is want, that too.
func (s *Scanner) munch2(want byte, both, single token.Kind) token.Kind {
	s.next()
	if s.ch == rune(want) {
		s.next()
		return both
	}
	return single
}

// scanComment consumes a line comment or a general comment. actsAsTerminator
// reports whether the comment behaves as a line terminator rather than as white
// space, which every comment does except a general comment spanning no line.
func (s *Scanner) scanComment() (lit string, actsAsTerminator bool) {
	offs := s.offset
	s.next() // '/'

	if s.ch == '/' {
		s.next()
		for s.ch != eof && !isLineTerminator(s.ch) {
			s.next()
		}
		return string(s.src[offs:s.offset]), true
	}

	s.next() // '*'
	terminated := false
	for s.ch != eof {
		if isLineTerminator(s.ch) {
			actsAsTerminator = true
			s.consumeLineTerminator()
			continue
		}
		ch := s.ch
		s.next()
		// General comments do not nest; the first */ closes.
		if ch == '*' && s.ch == '/' {
			s.next()
			terminated = true
			break
		}
	}
	if !terminated {
		s.error(diag.UnterminatedComment, offs)
	}
	return string(s.src[offs:s.offset]), actsAsTerminator
}

// Tokenize scans src to completion, including the trailing EOF. Convenience
// for tests and tooling; the parser drives Scan directly.
func Tokenize(file *token.File, src []byte, rep diag.Reporter, mode Mode) []token.Token {
	var s Scanner
	s.Init(file, src, rep, mode)
	var out []token.Token
	for {
		t := s.Scan()
		out = append(out, t)
		if t.Kind == token.EOF {
			return out
		}
	}
}