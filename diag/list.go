package diag

import (
	"sort"
	"strings"

	"github.com/vertex-language/vertex/token"
)

// Reporter is what a phase accepts. scanner, parser, and analyzer take one
// rather than a *List so an LSP can stream diagnostics and a batch build can
// collect them, without either phase knowing which.
type Reporter interface {
	Report(*Diagnostic)
}

// ReporterFunc adapts a function to Reporter.
type ReporterFunc func(*Diagnostic)

func (f ReporterFunc) Report(d *Diagnostic) { f(d) }

// List collects diagnostics. The zero value is ready to use.
type List struct {
	items  []*Diagnostic
	errors int
}

func (l *List) Report(d *Diagnostic) {
	if d == nil {
		return
	}
	l.items = append(l.items, d)
	if d.Sev == Error {
		l.errors++
	}
}

// Add is Report over a freshly constructed diagnostic, returning it so a caller
// can chain notes and fix-its onto an already-recorded diagnostic.
func (l *List) Add(code Code, pos, end token.Pos, args ...any) *Diagnostic {
	d := New(code, pos, end, args...)
	l.Report(d)
	return d
}

func (l *List) Items() []*Diagnostic { return l.items }
func (l *List) Len() int             { return len(l.items) }
func (l *List) NumErrors() int       { return l.errors }
func (l *List) HasErrors() bool      { return l.errors > 0 }

// Sort orders by position, then by code, stably. Positions are FileSet-global,
// so this also groups by file in AddFile order.
func (l *List) Sort() {
	sort.SliceStable(l.items, func(i, j int) bool {
		a, b := l.items[i], l.items[j]
		if a.Pos != b.Pos {
			return a.Pos < b.Pos
		}
		return a.Code < b.Code
	})
}

// Dedup removes diagnostics identical in code, position, and message. Parser
// recovery routinely re-reports at a resync point; this is cheaper than making
// every recovery path idempotent.
func (l *List) Dedup() {
	seen := make(map[[3]any]bool, len(l.items))
	out := l.items[:0]
	errs := 0
	for _, d := range l.items {
		k := [3]any{d.Code, d.Pos, d.Message}
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, d)
		if d.Sev == Error {
			errs++
		}
	}
	l.items, l.errors = out, errs
}

// Truncate keeps the first n diagnostics. Callers that truncate should say so
// in their output; this type does not synthesize a "too many errors" entry,
// because that text would be an unregistered message.
func (l *List) Truncate(n int) {
	if n < 0 || n >= len(l.items) {
		return
	}
	l.items = l.items[:n]
	l.errors = 0
	for _, d := range l.items {
		if d.Sev == Error {
			l.errors++
		}
	}
}

// Err returns a non-nil error if the list holds any Error-severity diagnostic.
// Idiomatic use: `if err := list.Err(); err != nil { ... }`.
func (l *List) Err() error {
	if !l.HasErrors() {
		return nil
	}
	return l
}

// Error renders messages only, one per line, so a *List can be returned as a
// plain error without dragging in the renderer or a FileSet.
func (l *List) Error() string {
	switch len(l.items) {
	case 0:
		return "no diagnostics"
	case 1:
		return l.items[0].Message
	}
	var b strings.Builder
	for i, d := range l.items {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(d.Message)
	}
	return b.String()
}