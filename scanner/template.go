package scanner

import "github.com/vertex-language/vertex/token"

// scanTemplateBody scans TemplateCharacters and the closing delimiter,
// starting from just past the opening `` ` `` or `}`. It returns the kind of
// the token that ends at the cursor.
//
// `${` pushes a frame and yields TEMPLATE_HEAD or TEMPLATE_MIDDLE; the
// matching `}` is routed back here by scanPunct (§4.5). An explicit stack, not
// a mode callback.
//
// Type-position templates (G.7) reuse these tokens exactly — the grammar has
// no TemplateTypeHead precisely so there is one way to tokenize `` `hello ${ ``
// — so this function does not know or care which side of the language it is on.
func (s *scanner) scanTemplateBody(start int) (token.Kind, token.Flags) {
	opening := s.src[start] // '`' for head/full, '}' for middle/tail
	var flags token.Flags

	for {
		if s.off >= len(s.src) {
			flags |= token.Unterminated
			s.error(start, s.off, "unterminated template literal")
			if opening == '`' {
				return token.TEMPLATE, flags
			}
			return token.TEMPLATE_TAIL, flags
		}

		switch c := s.src[s.off]; c {
		case '`':
			s.off++
			if opening == '`' {
				return token.TEMPLATE, flags
			}
			return token.TEMPLATE_TAIL, flags

		case '$':
			if s.off+1 < len(s.src) && s.src[s.off+1] == '{' {
				s.off += 2
				s.push(frameSubst, s.off-2)
				if opening == '`' {
					return token.TEMPLATE_HEAD, flags
				}
				return token.TEMPLATE_MIDDLE, flags
			}
			s.off++ // TemplateCharacter: `$` with lookahead ≠ {

		case '\\':
			flags |= token.HasEscape
			s.off++
			if s.consumeLineTerminator() {
				continue // LineContinuation
			}
			if s.off < len(s.src) {
				s.off++
			}
			// A.6 admits `\ NotEscapeSequence` as a TemplateCharacter, so a
			// malformed escape is not a lexical error here. It is only illegal
			// in the cooked value of an *untagged* template, which is a check
			// that needs to know about tagging — a parser question.

		case '\r':
			// LineTerminatorSequence is a TemplateCharacter; normalization of
			// CRLF belongs to the cooking phase, so consume and keep raw.
			s.off++
			if s.off < len(s.src) && s.src[s.off] == '\n' {
				s.off++
			}

		default:
			s.off++
		}
	}
}