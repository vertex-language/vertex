// Package scanner turns Vertex source into a token buffer.
//
// It records structure and interprets none of it (compiler_frontend.md §4):
// contextual keywords stay identifiers, literals keep raw spelling, and a
// newline is a flag rather than a token. The whole file is tokenized up front
// (§4.1) — speculation in the parser needs an immutable buffer and O(1)
// rollback.
package scanner

import (
	"github.com/vertex-language/vertex/token"
)

// Scan tokenizes file. It always returns a token slice ending in EOF, and
// diagnostics sorted by position.
//
// Comments are always emitted as COMMENT tokens. Whether they survive into a
// tree is the parser's ParseComments decision (§8.1); a scanner that dropped
// them could not serve §8.5's highlighters, and ParseFileTokens must accept a
// buffer from either caller.
func Scan(file *token.File) ([]token.Token, []token.Diagnostic) {
	s := &scanner{
		file: file,
		src:  file.Size2Bytes(),
	}
	s.run()
	token.SortDiagnostics(s.diags)
	return s.toks, s.diags
}

type scanner struct {
	file *token.File
	src  []byte
	off  int

	toks  []token.Token
	diags []token.Diagnostic

	stack []frame // §4.4, §4.5 — one stack, three jobs; see context.go

	// nlPending carries "a LineTerminator appeared" forward to the next token.
	// It is deliberately *not* cleared by a comment: `x ⏎ //c ⏎ y` must leave
	// NLBefore on y, because the parser skips comments and §6.2's expectSemi
	// would otherwise never see the break.
	nlPending bool

	// prev is the previous *significant* token — comments excluded. The `/`
	// decision (§4.4) and the brace classification both read it.
	prev    token.Token
	hasPrev bool

	// Set when the corresponding closer is emitted, so that a following `/`
	// can ask what the matching opener was.
	prevParenWasControl bool
	prevBraceWasBlock   bool
}

func (s *scanner) run() {
	// Estimate one token per six bytes; resizing a []Token is cheap but
	// pointless when the ratio is this predictable.
	s.toks = make([]token.Token, 0, len(s.src)/6+16)

	if len(s.src) >= 2 && s.src[0] == '#' && s.src[1] == '!' {
		start := s.off
		s.off += 2
		s.skipToLineEnd()
		s.emit(token.COMMENT, token.CtxNone, start, 0)
	}

	for {
		s.skipSpace()
		if s.off >= len(s.src) {
			break
		}
		s.scanToken()
	}

	// An explicit EOF token. Not in §8.8's listed surface, but the parser
	// would otherwise bounds-check every lookahead, and §8.5's consumers can
	// ignore it. Zero width, positioned at Size().
	start := s.off
	s.emit(token.EOF, token.CtxNone, start, 0)
}

// scanToken dispatches on the first byte. Every path advances s.off, so the
// loop in run cannot spin.
func (s *scanner) scanToken() {
	start := s.off
	c := s.src[s.off]

	switch {
	case c == '"' || c == '\'':
		s.scanString(start, c)
		return
	case c == '`':
		s.off++
		kind, flags := s.scanTemplateBody(start)
		s.emit(kind, token.CtxNone, start, flags)
		return
	case c >= '0' && c <= '9':
		s.scanNumber(start)
		return
	case c == '.':
		if s.off+1 < len(s.src) && isDecimalDigit(s.src[s.off+1]) {
			s.scanNumber(start) // `.5` — DecimalLiteral's second production
			return
		}
	case c == '/':
		s.scanSlash(start)
		return
	case c == '#':
		s.off++
		if s.off < len(s.src) && isIdentStartByte(s.src[s.off]) {
			flags := s.scanIdentTail()
			s.emit(token.PRIVATE_IDENT, token.CtxNone, start, flags)
			return
		}
		s.emit(token.HASH, token.CtxNone, start, 0)
		return
	}

	if isIdentStartByte(c) || c >= 0x80 || c == '\\' {
		if s.scanIdent(start) {
			return
		}
		// Fell through: not an identifier start after all (a stray high byte).
		s.error(start, s.off, "invalid character in source")
		s.emit(token.INVALID, token.CtxNone, start, 0)
		return
	}

	s.scanPunct(start)
}

// ---------------------------------------------------------------------------
// whitespace

func (s *scanner) skipSpace() {
	for s.off < len(s.src) {
		switch c := s.src[s.off]; c {
		case ' ', '\t', '\v', '\f':
			s.off++
		case '\n':
			s.off++
			s.nlPending = true
		case '\r':
			s.off++
			if s.off < len(s.src) && s.src[s.off] == '\n' {
				s.off++ // CRLF is one LineTerminatorSequence
			}
			s.nlPending = true
		case 0xC2: // U+00A0 NBSP
			if s.off+1 < len(s.src) && s.src[s.off+1] == 0xA0 {
				s.off += 2
				continue
			}
			return
		case 0xE2: // U+2028 LS, U+2029 PS
			if s.off+2 < len(s.src) && s.src[s.off+1] == 0x80 &&
				(s.src[s.off+2] == 0xA8 || s.src[s.off+2] == 0xA9) {
				s.off += 3
				s.nlPending = true
				continue
			}
			return
		case 0xEF: // U+FEFF ZWNBSP
			if s.off+2 < len(s.src) && s.src[s.off+1] == 0xBB && s.src[s.off+2] == 0xBF {
				s.off += 3
				continue
			}
			return
		case 0xE3: // U+3000 ideographic space
			if s.off+2 < len(s.src) && s.src[s.off+1] == 0x80 && s.src[s.off+2] == 0x80 {
				s.off += 3
				continue
			}
			return
		default:
			return
		}
	}
}

func (s *scanner) skipToLineEnd() {
	for s.off < len(s.src) {
		c := s.src[s.off]
		if c == '\n' || c == '\r' {
			return
		}
		if c == 0xE2 && s.off+2 < len(s.src) && s.src[s.off+1] == 0x80 &&
			(s.src[s.off+2] == 0xA8 || s.src[s.off+2] == 0xA9) {
			return
		}
		s.off++
	}
}

// ---------------------------------------------------------------------------
// comments

// scanSlash resolves the one context-dependent production in the lexical
// grammar (§4.4): `/` opens a comment, a regex, or a division.
func (s *scanner) scanSlash(start int) {
	if s.off+1 < len(s.src) {
		switch s.src[s.off+1] {
		case '/':
			s.off += 2
			s.skipToLineEnd()
			s.emit(token.COMMENT, token.CtxNone, start, 0)
			return
		case '*':
			s.scanBlockComment(start)
			return
		}
	}
	if s.regexAllowed() {
		s.scanRegex(start)
		return
	}
	s.scanPunct(start) // QUO or QUO_ASSIGN
}

func (s *scanner) scanBlockComment(start int) {
	s.off += 2
	spansLine := false
	closed := false
	for s.off < len(s.src) {
		c := s.src[s.off]
		if c == '*' && s.off+1 < len(s.src) && s.src[s.off+1] == '/' {
			s.off += 2
			closed = true
			break
		}
		if c == '\n' || c == '\r' {
			spansLine = true
		} else if c == 0xE2 && s.off+2 < len(s.src) && s.src[s.off+1] == 0x80 &&
			(s.src[s.off+2] == 0xA8 || s.src[s.off+2] == 0xA9) {
			spansLine = true
		}
		s.off++
	}
	var flags token.Flags
	if !closed {
		flags |= token.Unterminated
		s.error(start, s.off, "unterminated block comment")
	}
	s.emit(token.COMMENT, token.CtxNone, start, flags)
	if spansLine {
		// §4.3: a comment spanning a newline sets NLBefore on the next real
		// token, whether or not comments are retained.
		s.nlPending = true
	}
}

// ---------------------------------------------------------------------------
// emit

func (s *scanner) emit(kind token.Kind, ctx token.Ctx, start int, flags token.Flags) {
	if s.nlPending {
		flags |= token.NLBefore
	}
	t := token.Token{
		Kind:  kind,
		Ctx:   ctx,
		Flags: flags,
		Pos:   s.file.PosAt(start),
		End:   s.file.PosAt(s.off),
	}
	s.toks = append(s.toks, t)

	if kind == token.COMMENT {
		// nlPending survives a comment; see the field comment.
		return
	}
	s.nlPending = false
	s.prev = t
	s.hasPrev = true
}

func (s *scanner) error(start, end int, msg string) {
	p := s.file.PosAt(start)
	e := token.NoPos
	if end > start {
		e = s.file.PosAt(end)
	}
	s.diags = append(s.diags, token.Diagnostic{Pos: p, End: e, Msg: msg})
}