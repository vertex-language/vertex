package scanner

import "github.com/vertex-language/vertex/token"

// scanNumber finds the extent of a NumericLiteral (A.4) and validates its
// shape. It does not decode: `1_024` yields the bytes as written, with no
// value, no width, and no separator stripping (§4.6, §8.3). Decoding belongs
// to a phase that knows the target type.
//
// Validation still happens here, because the extent and the shape are the same
// walk — a malformed separator changes where the literal ends.
func (s *scanner) scanNumber(start int) {
	kind := token.NUMBER
	var flags token.Flags

	if s.src[s.off] == '.' {
		s.off++
		flags |= s.scanDigits(10)
		flags |= s.scanExponent()
		s.finishNumber(start, kind, flags)
		return
	}

	if s.src[s.off] == '0' {
		s.off++
		switch c := s.at(0); {
		case c == 'b' || c == 'B':
			s.off++
			flags |= s.scanRadix(start, 2, "binary")
			kind, flags = s.maybeBigInt(kind, flags)
			s.finishNumber(start, kind, flags)
			return
		case c == 'o' || c == 'O':
			s.off++
			flags |= s.scanRadix(start, 8, "octal")
			kind, flags = s.maybeBigInt(kind, flags)
			s.finishNumber(start, kind, flags)
			return
		case c == 'x' || c == 'X':
			s.off++
			flags |= s.scanRadix(start, 16, "hexadecimal")
			kind, flags = s.maybeBigInt(kind, flags)
			s.finishNumber(start, kind, flags)
			return
		case c == 'n':
			s.off++ // DecimalBigIntegerLiteral: `0n`
			s.finishNumber(start, token.BIGINT, flags)
			return
		case isDecimalDigit(c):
			// A.4 has no LegacyOctalIntegerLiteral. `0123` is not a literal.
			for isDecimalDigit(s.at(0)) {
				s.off++
			}
			s.error(start, s.off, "leading zeros are not permitted in numeric literals")
			s.finishNumber(start, kind, flags)
			return
		}
		// Plain `0`, possibly followed by a fraction or exponent.
	} else {
		flags |= s.scanDigits(10)
	}

	if s.at(0) == 'n' {
		s.off++
		s.finishNumber(start, token.BIGINT, flags)
		return
	}

	fractional := false
	if s.at(0) == '.' {
		fractional = true
		s.off++
		flags |= s.scanDigitsOpt(10)
	}
	if e := s.scanExponent(); e != 0 || s.hadExponent {
		fractional = true
		flags |= e
	}
	if fractional && s.at(0) == 'n' {
		s.off++
		s.error(start, s.off, "bigint literal must be an integer")
		s.finishNumber(start, token.BIGINT, flags)
		return
	}
	s.finishNumber(start, kind, flags)
}

func (s *scanner) maybeBigInt(kind token.Kind, flags token.Flags) (token.Kind, token.Flags) {
	if s.at(0) == 'n' {
		s.off++
		return token.BIGINT, flags
	}
	return kind, flags
}

// finishNumber emits, after checking that the literal is not immediately
// followed by an identifier character. `3in` and `0x1p` would otherwise scan
// as two tokens and produce a confusing parse error instead of a lexical one.
func (s *scanner) finishNumber(start int, kind token.Kind, flags token.Flags) {
	if c := s.at(0); isIdentStartByte(c) || isDecimalDigit(c) || c == '\\' || c >= 0x80 {
		bad := s.off
		s.scanIdentTail()
		s.error(bad, s.off, "identifier cannot immediately follow a numeric literal")
	}
	s.emit(kind, token.CtxNone, start, flags)
}

// hadExponent is set by scanExponent so that `1e5n` is caught. It is reset on
// entry; the scanner is single-threaded per file.
func (s *scanner) scanExponent() token.Flags {
	s.hadExponent = false
	if c := s.at(0); c != 'e' && c != 'E' {
		return 0
	}
	mark := s.off
	s.off++
	if c := s.at(0); c == '+' || c == '-' {
		s.off++
	}
	if !isDecimalDigit(s.at(0)) {
		s.error(mark, s.off, "exponent has no digits")
		return 0
	}
	s.hadExponent = true
	return s.scanDigits(10)
}

// scanRadix handles the digits after 0b / 0o / 0x, where at least one digit is
// required (A.4 uses BinaryDigits[+Sep] with no _opt).
func (s *scanner) scanRadix(start, base int, name string) token.Flags {
	before := s.off
	flags := s.scanDigitsOpt(base)
	if s.off == before {
		s.error(start, s.off, name+" literal has no digits")
	}
	return flags
}

func (s *scanner) scanDigits(base int) token.Flags { return s.scanDigitsOpt(base) }

// scanDigitsOpt consumes digits and NumericLiteralSeparators, diagnosing a
// leading, trailing, or doubled separator while still consuming it — the raw
// spelling must cover the whole literal or the span is wrong.
func (s *scanner) scanDigitsOpt(base int) token.Flags {
	var flags token.Flags
	lastWasSep := false
	first := true

	for s.off < len(s.src) {
		c := s.src[s.off]
		if c == '_' {
			flags |= token.HasEscape // "raw spelling needs decoding"
			if first {
				s.error(s.off, s.off+1, "numeric separator cannot appear at the start of a digit sequence")
			} else if lastWasSep {
				s.error(s.off, s.off+1, "numeric separator cannot appear twice in a row")
			}
			lastWasSep = true
			first = false
			s.off++
			continue
		}
		if !digitInBase(c, base) {
			break
		}
		lastWasSep = false
		first = false
		s.off++
	}
	if lastWasSep {
		s.error(s.off-1, s.off, "numeric separator cannot appear at the end of a digit sequence")
	}
	return flags
}

func digitInBase(c byte, base int) bool {
	switch base {
	case 2:
		return c == '0' || c == '1'
	case 8:
		return c >= '0' && c <= '7'
	case 10:
		return isDecimalDigit(c)
	case 16:
		return isHexDigit(c)
	}
	return false
}