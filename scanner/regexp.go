package scanner

import "github.com/vertex-language/vertex/token"

// scanRegex scans a RegularExpressionLiteral (A.7) once §4.4's
// previous-significant-token rule has decided that `/` opens one.
//
// The body is delimited, not parsed: Pattern and everything under it (A.7) is
// a separate grammar that a later phase compiles. All this needs to know is
// where the closing `/` is, which means tracking character classes — a `/`
// inside `[...]` is literal.
//
// If the regex production is ever dropped, this file and §4.4's bracket
// classification go with it, and nothing else changes.
func (s *scanner) scanRegex(start int) {
	s.off++ // opening slash
	var flags token.Flags
	inClass := false

	for {
		if s.off >= len(s.src) {
			flags |= token.Unterminated
			s.error(start, s.off, "unterminated regular expression literal")
			break
		}

		c := s.src[s.off]

		// A LineTerminator ends the literal unterminated: RegularExpressionChar
		// excludes them, and recovering at the line boundary keeps the
		// misclassification cost to one diagnostic rather than a cascade.
		if c == '\n' || c == '\r' ||
			(c == 0xE2 && s.off+2 < len(s.src) && s.src[s.off+1] == 0x80 &&
				(s.src[s.off+2] == 0xA8 || s.src[s.off+2] == 0xA9)) {
			flags |= token.Unterminated
			s.error(start, s.off, "unterminated regular expression literal")
			break
		}

		switch {
		case c == '\\':
			flags |= token.HasEscape
			s.off++
			if s.off < len(s.src) {
				s.off++
			}
			continue
		case c == '[':
			inClass = true
		case c == ']':
			inClass = false
		case c == '/' && !inClass:
			s.off++
			s.scanRegexFlags()
			s.emit(token.REGEX, token.CtxNone, start, flags)
			return
		}
		s.off++
	}

	s.emit(token.REGEX, token.CtxNone, start, flags)
}

// scanRegexFlags consumes RegularExpressionFlags: any run of IdentifierPart
// characters. Which letters are valid, and whether one repeats, is a question
// for whoever compiles the pattern.
func (s *scanner) scanRegexFlags() {
	for s.off < len(s.src) {
		c := s.src[s.off]
		if isIdentPartByte(c) {
			s.off++
			continue
		}
		if c >= 0x80 {
			before := s.off
			s.scanIdentTail()
			if s.off == before {
				return
			}
			continue
		}
		return
	}
}