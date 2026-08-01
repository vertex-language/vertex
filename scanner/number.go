package scanner

import (
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
)

// digitVal returns the value of ch as a hexadecimal digit, or 16.
func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch-'a') + 10
	case 'A' <= ch && ch <= 'F':
		return int(ch-'A') + 10
	}
	return 16
}

func baseName(base int) string {
	switch base {
	case 2:
		return "binary"
	case 8:
		return "octal"
	case 16:
		return "hexadecimal"
	}
	return "decimal"
}

// scanDigits consumes a digit run in the given base, enforcing A.1.5.1's
// separator rules: '_' may not lead a run, trail it, or be doubled. It returns
// the number of digits consumed, separators excluded.
//
// A decimal digit out of range for a narrow base is consumed rather than
// terminating the run, so 0b1234 produces one diagnostic per bad digit instead
// of a literal followed by a spurious second literal.
func (s *Scanner) scanDigits(base int) int {
	count := 0
	started := false
	prevSep := false
	sepOffset := -1

	for {
		switch {
		case s.ch == '_':
			switch {
			case !started:
				s.error(diag.SeparatorLeads, s.offset)
			case prevSep:
				s.error(diag.SeparatorDoubled, s.offset)
			}
			sepOffset = s.offset
			prevSep = true
			started = true
			s.next()

		case digitVal(s.ch) < base:
			count++
			prevSep = false
			started = true
			s.next()

		case base < 10 && '0' <= s.ch && s.ch <= '9':
			s.error(diag.DigitOutOfRange, s.offset, s.ch, baseName(base))
			count++
			prevSep = false
			started = true
			s.next()

		default:
			if prevSep {
				s.error(diag.SeparatorTrails, sepOffset)
			}
			return count
		}
	}
}

// scanNumber consumes a NumericLiteral (A.1.5.1) and reports whether it is an
// integer or a float.
//
// Two decisions here are load-bearing beyond this function:
//
//   - A fractional part must be non-empty, so `1.` is not a literal. That is
//     what makes `1..5` scan as INT DOTDOT INT rather than requiring lookahead
//     surgery.
//   - A hexadecimal float requires its binary exponent (A.1.5.1 ⊢), so `0xC.3`
//     is diagnosed as an incomplete float rather than silently accepted or
//     split into three tokens.
func (s *Scanner) scanNumber() (token.Kind, string) {
	offs := s.offset
	kind := token.INT
	base := 10
	prefix := ""

	if s.ch == '0' {
		switch s.peek() {
		case 'b', 'B':
			prefix, base = "0b", 2
			s.next()
			s.next()
		case 'o', 'O':
			prefix, base = "0o", 8
			s.next()
			s.next()
		case 'x', 'X':
			prefix, base = "0x", 16
			s.next()
			s.next()
		}
	}

	switch base {
	case 16:
		n := s.scanDigits(16)
		hasFrac := false
		if s.ch == '.' && s.peek() != '.' {
			s.next()
			// HexDigitsWithSeparatorsopt — the fraction may be empty here,
			// unlike the decimal case.
			s.scanDigits(16)
			hasFrac = true
		}
		if n == 0 && !hasFrac {
			s.error(diag.EmptyDigits, offs, "'"+prefix+"'")
		}
		switch {
		case s.ch == 'p' || s.ch == 'P':
			kind = token.FLOAT
			s.next()
			if s.ch == '+' || s.ch == '-' {
				s.next()
			}
			if s.scanDigits(10) == 0 {
				s.error(diag.MissingExponentDigit, s.offset)
			}
		case hasFrac:
			s.error(diag.HexFloatNoExponent, offs)
			kind = token.FLOAT
		}

	case 2, 8:
		if s.scanDigits(base) == 0 {
			s.error(diag.EmptyDigits, offs, "'"+prefix+"'")
		}

	default:
		s.scanDigits(10)
		if s.ch == '.' && '0' <= s.peek() && s.peek() <= '9' {
			kind = token.FLOAT
			s.next()
			s.scanDigits(10)
		}
		if s.ch == 'e' || s.ch == 'E' {
			kind = token.FLOAT
			s.next()
			if s.ch == '+' || s.ch == '-' {
				s.next()
			}
			if s.scanDigits(10) == 0 {
				s.error(diag.MissingExponentDigit, s.offset)
			}
		}
	}

	if isIdentPart(s.ch) {
		bad := s.offset
		for isIdentPart(s.ch) {
			s.next()
		}
		s.errorSpan(diag.NumberJoinedToIdent, bad, s.offset)
	}

	return kind, string(s.src[offs:s.offset])
}

// scanTupleIndex consumes the digit run of a positional tuple access (A.4.3).
//
// It is entered only when the previous token was PERIOD, where a digit can mean
// nothing else: EnumShorthand takes an identifier, and Vertex has no
// leading-dot float literals. Scanning here rather than splitting a FLOAT in
// the parser is what makes `t.0.0` fall out — the second dot ends this run and
// the rule fires again.
func (s *Scanner) scanTupleIndex() string {
	offs := s.offset
	for '0' <= s.ch && s.ch <= '9' {
		s.next()
	}
	if isIdentPart(s.ch) {
		bad := s.offset
		for isIdentPart(s.ch) {
			s.next()
		}
		s.errorSpan(diag.NumberJoinedToIdent, bad, s.offset)
	}
	return string(s.src[offs:s.offset])
}