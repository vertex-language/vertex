package scanner

import (
	"unicode"
	"unicode/utf8"

	"github.com/vertex-language/vertex/diag"
)

// isIDStart reports whether ch has the Unicode ID_Start property.
//
// ID_Start is derived as L + Nl + Other_ID_Start, minus Pattern_Syntax and
// Pattern_White_Space. Go's unicode package exposes the pieces but not the
// derived property, so it is assembled here.
//
// This should become a generated table pinned to a stated Unicode version. The
// property set changes between versions, and an implementation must not derive
// these from a host toolchain whose version may drift — a language whose
// identifier rules move with its compiler's stdlib has made its grammar depend
// on something outside its own specification.
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

// isIdentStart matches the first factor of `identifier`: unicode_id_start or
// '_'.
func isIdentStart(ch rune) bool {
	if ch < utf8.RuneSelf {
		return ch == '_' ||
			'a' <= ch && ch <= 'z' ||
			'A' <= ch && ch <= 'Z'
	}
	return ch != eof && isIDStart(ch)
}

// isIdentPart matches the repeated factor of `identifier`: unicode_id_continue
// or '_'.
//
// '$' is absent by construction. It is not an identifier character in any
// position, and it is a diagnosed error rather than a silently unmatched
// character, so scanIdent handles it explicitly instead.
func isIdentPart(ch rune) bool {
	if ch < utf8.RuneSelf {
		return ch == '_' ||
			'a' <= ch && ch <= 'z' ||
			'A' <= ch && ch <= 'Z' ||
			'0' <= ch && ch <= '9'
	}
	return ch != eof && isIDContinue(ch)
}

// scanIdent consumes an IdentifierName.
//
// A '$' anywhere in the run is consumed along with it and reported once over
// the span it occupies, so one pasted foreign name yields one diagnostic and
// one identifier token rather than a cascade of fragments. Keyword
// classification happens in the caller through token.Lookup; the blank
// identifier '_' falls out as an ordinary identifier, which is what
// ast.Ident.IsBlank tests later.
func (s *Scanner) scanIdent() string {
	offs := s.offset
	dollar := -1

	for {
		switch {
		case isIdentPart(s.ch):
			s.next()
		case s.ch == '$':
			if dollar < 0 {
				dollar = s.offset
			}
			s.next()
		default:
			if dollar >= 0 {
				s.errorSpan(diag.DollarInIdent, dollar, s.offset)
			}
			return string(s.src[offs:s.offset])
		}
	}
}