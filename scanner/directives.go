package scanner

import (
	"regexp"

	"github.com/vertex-language/vertex/token"
)

// Directives collects comments matching re and returns them as expected
// diagnostics, keyed by the 1-based line of the comment.
//
// This ships in scanner rather than in a test package (§7) so that fixtures
// outside the repo can use it. Each match is positioned at the *preceding*
// token, which is what makes `let x: = 1; // ERROR "expected type"` point at
// the `=` rather than at the comment.
//
// If re has a capturing group, group 1 becomes Msg; otherwise the whole match
// does.
//
// The result is a map, per §8.8's signature, but §7 forbids iterating a map to
// produce output. Callers that render or compare must sort the keys.
func Directives(file *token.File, re *regexp.Regexp) map[int][]token.Diagnostic {
	toks, _ := Scan(file)
	out := make(map[int][]token.Diagnostic)

	var prev token.Token
	var hasPrev bool

	for _, t := range toks {
		if t.Kind != token.COMMENT {
			if t.Kind != token.EOF {
				prev = t
				hasPrev = true
			}
			continue
		}

		text := file.Slice(t)
		m := re.FindSubmatchIndex(text)
		if m == nil {
			continue
		}

		lo, hi := m[0], m[1]
		if len(m) >= 4 && m[2] >= 0 {
			lo, hi = m[2], m[3]
		}

		pos := t.Pos
		end := token.NoPos
		if hasPrev {
			pos, end = prev.Pos, prev.End
		}

		line := file.Position(t.Pos).Line
		out[line] = append(out[line], token.Diagnostic{
			Pos: pos,
			End: end,
			Msg: string(text[lo:hi]),
		})
	}
	return out
}