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

// scanDigits consumes one digit run in the given base, enforcing the separator
// rules: '_' may appear between successive digits, and may not lead a run,
// trail one, or be doubled. It returns the number of digits consumed,
// separators excluded.
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

// scanNumber consumes a NumericLiteral and reports whether it is an integer or
// a float. Each base has its own tail because the three do not share a shape:
// only hexadecimal has a binary-exponent float form, only decimal has a bare
// exponent, and the two define a fractional part differently. There is no
// prefix-free octal form, so 0600 is the decimal integer 600.
//
// Two decisions here are load-bearing beyond this function:
//
//   - A decimal fractional part, if a point is written, must be non-empty, so
//     `1.` is not a literal. That is what makes `1..5` scan as
//     int_lit ".." int_lit without lookahead surgery elsewhere.
//   - A hexadecimal float requires its binary exponent, so `0xC.3` is
//     diagnosed as an incomplete float rather than silently accepted or split
//     into three tokens.
//
// A `.` immediately followed by another `.` is never consumed as a point, which
// is what keeps a range operator intact after any base.
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
		if n == 0 {
			s.error(diag.EmptyDigits, offs, "'"+prefix+"'")
		}
		hasPoint := false
		if s.ch == '.' && s.peek() != '.' {
			s.next()
			// The hexadecimal fraction may be empty; the exponent below is
			// what makes the literal well formed, not this digit run.
			s.scanDigits(16)
			hasPoint = true
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
		case hasPoint:
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

	s.rejectJoinedIdent()
	return kind, string(s.src[offs:s.offset])
}

// scanTupleIndex consumes the digit run of a positional tuple access.
//
// It is entered only where the selector-dot restriction applies, and its whole
// job is to produce an int_lit where longest match would have produced a
// float: no point is consumed and no exponent is recognized, so the next `.`
// ends this run and the rule fires again, which is what makes a chain compose.
//
// Separators are enforced here as they are in any other digit run, because the
// token produced is an ordinary int_lit. That a tuple index must additionally
// be written in decimal with no '_' is a static rule over the spelling this
// returns, not a lexical one.
func (s *Scanner) scanTupleIndex() string {
	offs := s.offset
	s.scanDigits(10)
	s.rejectJoinedIdent()
	return string(s.src[offs:s.offset])
}

// rejectJoinedIdent consumes any identifier characters directly following a
// numeric literal and reports them as one span, so `123abc` is one diagnosed
// token rather than an int_lit followed by an identifier.
func (s *Scanner) rejectJoinedIdent() {
	if !isIdentPart(s.ch) {
		return
	}
	bad := s.offset
	for isIdentPart(s.ch) {
		s.next()
	}
	s.errorSpan(diag.NumberJoinedToIdent, bad, s.offset)
}