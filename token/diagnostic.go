package token

import "sort"

// Diagnostic lives in token, not in scanner or parser, because both emit them
// and both would otherwise need a shared package above token (§3).
//
// It carries no severity and no code. A cap on reported errors and any
// grouping or coloring are presentation decisions and belong to whatever
// renders them; diagnostic slices are returned sorted by Pos and are never
// truncated.
type Diagnostic struct {
	Pos Pos    // first byte of the offending construct
	End Pos    // NoPos when the diagnostic is a point, not a span
	Msg string // one line, no trailing punctuation, no position prefix
}

// IsSpan reports whether d covers a range rather than a point.
func (d Diagnostic) IsSpan() bool { return d.End != NoPos && d.End > d.Pos }

// SortDiagnostics orders ds by Pos, then End, then Msg. The tiebreakers are
// what make output reproducible when two diagnostics share a position — §7's
// determinism guarantee is untestable without them.
//
// Ordering across files is not defined here and cannot be: a bare-offset Pos
// has no meaning outside its unit (§8.2). Callers holding several units key on
// (unit index in the caller's input list, Pos).
func SortDiagnostics(ds []Diagnostic) {
	sort.SliceStable(ds, func(i, j int) bool {
		a, b := ds[i], ds[j]
		if a.Pos != b.Pos {
			return a.Pos < b.Pos
		}
		if a.End != b.End {
			return a.End < b.End
		}
		return a.Msg < b.Msg
	})
}