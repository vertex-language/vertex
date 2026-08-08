package scanner

import (
	"unicode"
	"unicode/utf8"

	"github.com/vertex-language/vertex/token"
)

// scanIdent scans an IdentifierName (A.2) and classifies it. Returns false if
// the byte at the cursor does not in fact start an identifier.
//
// Classification is one call to token.LookupIdent: a ReservedWord becomes its
// own Kind, a ContextualKeyword becomes IDENT plus a Ctx, everything else is
// a bare IDENT. Predeclared type names (int32, usize, float64) and every
// PredefinedType member are ordinary identifiers here (§4.6) — the scanner
// does not know they are types.
func (s *scanner) scanIdent(start int) bool {
	flags, ok := s.scanIdentStart()
	if !ok {
		return false
	}
	flags |= s.scanIdentTail()

	kind := token.IDENT
	ctx := token.CtxNone
	if flags&token.HasEscape == 0 {
		// An escaped spelling never matches a keyword: `\u0069f` is an
		// identifier named "if", not the IF token. This is the whole reason
		// HasEscape exists on identifiers.
		kind, ctx = token.LookupIdent(s.src[start:s.off])
	}
	s.emit(kind, ctx, start, flags)
	return true
}

func (s *scanner) scanIdentStart() (token.Flags, bool) {
	if s.off >= len(s.src) {
		return 0, false
	}
	c := s.src[s.off]
	switch {
	case c == '\\':
		if !s.scanUnicodeEscape() {
			return token.HasEscape, false
		}
		return token.HasEscape, true
	case isIdentStartByte(c):
		s.off++
		return 0, true
	case c >= 0x80:
		r, size := utf8.DecodeRune(s.src[s.off:])
		if isIdentStartRune(r) {
			s.off += size
			return 0, true
		}
	}
	return 0, false
}

// scanIdentTail consumes IdentifierPart characters. Split out because
// PRIVATE_IDENT (`#foo`) reuses it after the `#`.
func (s *scanner) scanIdentTail() token.Flags {
	var flags token.Flags
	for s.off < len(s.src) {
		c := s.src[s.off]
		switch {
		case isIdentPartByte(c):
			s.off++
		case c == '\\':
			if !s.scanUnicodeEscape() {
				return flags | token.HasEscape
			}
			flags |= token.HasEscape
		case c >= 0x80:
			r, size := utf8.DecodeRune(s.src[s.off:])
			if !isIdentPartRune(r) {
				return flags
			}
			s.off += size
		default:
			return flags
		}
	}
	return flags
}

// scanUnicodeEscape consumes `\uXXXX` or `\u{...}`. The code point is not
// decoded and not validated against IdentifierPart — that is a later phase's
// job, and rejecting here would mean the scanner knows what an identifier
// means rather than where it ends.
func (s *scanner) scanUnicodeEscape() bool {
	start := s.off
	s.off++ // backslash
	if s.at(0) != 'u' {
		s.error(start, s.off, "expected unicode escape sequence in identifier")
		return false
	}
	s.off++
	if s.at(0) == '{' {
		s.off++
		n := 0
		for s.off < len(s.src) && isHexDigit(s.src[s.off]) {
			s.off++
			n++
		}
		if n == 0 || s.at(0) != '}' {
			s.error(start, s.off, "malformed unicode code point escape")
			return false
		}
		s.off++
		return true
	}
	for i := 0; i < 4; i++ {
		if !isHexDigit(s.at(0)) {
			s.error(start, s.off, "unicode escape requires four hex digits")
			return false
		}
		s.off++
	}
	return true
}

func isIdentStartByte(c byte) bool {
	return c == '$' || c == '_' ||
		(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentPartByte(c byte) bool {
	return isIdentStartByte(c) || (c >= '0' && c <= '9')
}

func isIdentStartRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.In(r, unicode.Nl, unicode.Other_ID_Start)
}

func isIdentPartRune(r rune) bool {
	if isIdentStartRune(r) {
		return true
	}
	if r == 0x200C || r == 0x200D { // ZWNJ, ZWJ
		return true
	}
	return unicode.In(r, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc, unicode.Other_ID_Continue)
}

func isDecimalDigit(c byte) bool { return c >= '0' && c <= '9' }

func isHexDigit(c byte) bool {
	return isDecimalDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}