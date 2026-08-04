package scanner

import (
	"github.com/vertex-language/vertex/diag"
)

// scanString consumes an interpreted string literal. The returned text is the
// raw source spelling, both quotes and every escape exactly as written: a
// formatter needs the original, and decoding belongs to a later phase.
//
// A `'` is legal unescaped here, since string_value excludes only `"` and `\`.
func (s *Scanner) scanString() string {
	offs := s.offset
	s.next() // opening quote

	for {
		if s.ch == eof {
			s.error(diag.UnterminatedString, offs)
			break
		}
		// A line terminator may not appear in an interpreted string. It is
		// left unconsumed so that it still ends the statement and recovery
		// resumes on the next line rather than swallowing it.
		if isLineTerminator(s.ch) {
			s.error(diag.NewlineInString, offs)
			break
		}
		ch := s.ch
		s.next()
		if ch == '"' {
			break
		}
		if ch == '\\' {
			s.scanEscape()
		}
	}
	return string(s.src[offs:s.offset])
}

// scanRawString consumes a raw string literal. It recognizes no escape
// sequence, and every line terminator it spans is part of its value.
//
// Those terminators advance the line table but do not set nlPending — they are
// string content, not statement structure.
func (s *Scanner) scanRawString() string {
	offs := s.offset
	s.next() // opening back quote

	for {
		if s.ch == eof {
			s.error(diag.UnterminatedRawString, offs)
			break
		}
		if isLineTerminator(s.ch) {
			s.consumeLineTerminator()
			continue
		}
		ch := s.ch
		s.next()
		if ch == '`' {
			break
		}
	}
	return string(s.src[offs:s.offset])
}

// scanChar consumes a character literal, which denotes exactly one Unicode
// scalar value. The count below is in source units, so an escape counts once
// and a multi-byte rune counts once.
//
// A `"` is legal unescaped here, since char_value excludes only `'` and `\`.
func (s *Scanner) scanChar() string {
	offs := s.offset
	s.next() // opening quote

	n := 0
	terminated := false
	for {
		if s.ch == eof || isLineTerminator(s.ch) {
			s.error(diag.UnterminatedChar, offs)
			break
		}
		ch := s.ch
		s.next()
		if ch == '\'' {
			terminated = true
			break
		}
		n++
		if ch == '\\' {
			s.scanEscape()
		}
	}

	if terminated {
		switch {
		case n == 0:
			s.errorSpan(diag.EmptyChar, offs, s.offset)
		case n > 1:
			s.errorSpan(diag.CharTooLong, offs, s.offset)
		}
	}
	return string(s.src[offs:s.offset])
}

// scanEscape consumes an escape sequence. It is entered with the backslash
// already consumed, so s.ch is the character following it.
func (s *Scanner) scanEscape() {
	backslash := s.offset - 1

	switch s.ch {
	case '\'', '"', '\\', 'n', 'r', 't', 'v', 'b', 'f', '0':
		s.next()

	case 'x':
		s.next()
		for i := 0; i < 2; i++ {
			if digitVal(s.ch) >= 16 {
				s.errorSpan(diag.InvalidHexEscape, backslash, s.offset)
				return
			}
			s.next()
		}

	case 'u':
		s.next()
		if s.ch != '{' {
			s.errorSpan(diag.InvalidUnicodeEscape, backslash, s.offset)
			return
		}
		s.next()

		digits := s.offset
		var v rune
		n := 0
		for digitVal(s.ch) < 16 {
			if v <= 0x10FFFF {
				v = v*16 + rune(digitVal(s.ch))
			}
			n++
			s.next()
		}
		text := string(s.src[digits:s.offset])

		if n == 0 || s.ch != '}' {
			s.errorSpan(diag.InvalidUnicodeEscape, backslash, s.offset)
			return
		}
		s.next() // '}'

		// A surrogate is a code point but not a scalar value, and the escape
		// must denote a scalar value. Both failures share one code.
		if v > 0x10FFFF || (0xD800 <= v && v < 0xE000) {
			s.errorSpan(diag.UnicodeEscapeRange, backslash, s.offset, "U+"+text)
		}

	default:
		if s.ch == eof || isLineTerminator(s.ch) {
			return // the enclosing literal reports its own unterminated error
		}
		s.errorSpan(diag.InvalidEscape, backslash, s.offset, s.ch)
		s.next()
	}
}