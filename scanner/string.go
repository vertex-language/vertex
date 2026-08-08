package scanner

import "github.com/vertex-language/vertex/token"

// scanString handles A.5. Escapes are recorded, never decoded: HasEscape tells
// a later phase the raw spelling needs work, and that phase knows the target
// encoding.
//
// An unterminated string still gets an exact span and a real token, so the
// parser's recovery (§6.4) has something to work with and §1's non-zero span
// invariant holds.
func (s *scanner) scanString(start int, quote byte) {
	s.off++ // opening quote
	var flags token.Flags

	for {
		if s.off >= len(s.src) {
			flags |= token.Unterminated
			s.error(start, s.off, "unterminated string literal")
			break
		}
		c := s.src[s.off]

		if c == quote {
			s.off++
			break
		}

		if c == '\\' {
			flags |= token.HasEscape
			s.off++
			if s.off >= len(s.src) {
				continue // reported as unterminated on the next pass
			}
			// LineContinuation: a backslash followed by a
			// LineTerminatorSequence is legal and does not end the string.
			if s.consumeLineTerminator() {
				continue
			}
			s.off++ // the escaped character; EscapeSequence is not validated here
			continue
		}

		// A raw LineTerminator ends a string. <LS> and <PS> do not — A.5 lists
		// them as DoubleStringCharacter alternatives.
		if c == '\n' || c == '\r' {
			flags |= token.Unterminated
			s.error(start, s.off, "unterminated string literal")
			break
		}
		s.off++
	}

	s.emit(token.STRING, token.CtxNone, start, flags)
}

// consumeLineTerminator consumes one LineTerminatorSequence at the cursor and
// reports whether it did. CRLF counts as one.
func (s *scanner) consumeLineTerminator() bool {
	if s.off >= len(s.src) {
		return false
	}
	switch s.src[s.off] {
	case '\n':
		s.off++
		return true
	case '\r':
		s.off++
		if s.off < len(s.src) && s.src[s.off] == '\n' {
			s.off++
		}
		return true
	case 0xE2:
		if s.off+2 < len(s.src) && s.src[s.off+1] == 0x80 &&
			(s.src[s.off+2] == 0xA8 || s.src[s.off+2] == 0xA9) {
			s.off += 3
			return true
		}
	}
	return false
}