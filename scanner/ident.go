package scanner

import (
	"unicode"
	"unicode/utf8"
)

// isIDStart reports whether ch has the Unicode ID_Start property (A.1.2).
//
// ID_Start is derived as L + Nl + Other_ID_Start, minus Pattern_Syntax and
// Pattern_White_Space. Go's unicode package exposes the pieces but not the
// derived property, so it is assembled here.
//
// This should become a generated table pinned to a specific Unicode version.
// The property set changes between versions, and a language whose identifier
// rules drift with the host toolchain's stdlib has made its grammar depend on
// something outside its own specification.
func isIDStart(ch rune) bool {
	return unicode.In(ch, unicode.L, unicode.Nl, unicode.Other_ID_Start) &&
		!unicode.In(ch, unicode.Pattern_Syntax, unicode.Pattern_White_Space)
}

// isIDContinue reports whether ch has the Unicode ID_Continue property.
func isIDContinue(ch rune) bool {
	return unicode.In(ch,
		unicode.L, unicode.Nl, unicode.Other_ID_Start,
		unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc,
		unicode.Other_ID_Continue,
	) && !unicode.In(ch, unicode.Pattern_Syntax, unicode.Pattern_White_Space)
}

// isIdentStart matches IdentifierStart: UnicodeIDStart or '_' (A.1.2).
func isIdentStart(ch rune) bool {
	if ch < utf8.RuneSelf {
		return ch == '_' ||
			'a' <= ch && ch <= 'z' ||
			'A' <= ch && ch <= 'Z'
	}
	return ch != eof && isIDStart(ch)
}

// isIdentPart matches IdentifierPart: UnicodeIDContinue or '_' (A.1.2).
// '$' is deliberately absent — A.1.2 ⊢ makes it a diagnosed error, not a
// silently-unmatched character.
func isIdentPart(ch rune) bool {
	if ch < utf8.RuneSelf {
		return ch == '_' ||
			'a' <= ch && ch <= 'z' ||
			'A' <= ch && ch <= 'Z' ||
			'0' <= ch && ch <= '9'
	}
	return ch != eof && isIDContinue(ch)
}

// scanIdent consumes an IdentifierName. Keyword classification happens in the
// caller via token.Lookup; the BlankIdentifier '_' falls out as an ordinary
// IDENT, which is what ast.Ident.IsBlank later tests.
func (s *Scanner) scanIdent() string {
	offs := s.offset
	for isIdentPart(s.ch) {
		s.next()
	}
	return string(s.src[offs:s.offset])
}