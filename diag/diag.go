package diag

import (
	"fmt"

	"github.com/vertex-language/vertex/token"
)

type Severity uint8

const (
	Error Severity = iota
	Warning
)

func (s Severity) String() string {
	switch s {
	case Warning:
		return "warning"
	default:
		return "error"
	}
}

// Diagnostic is the single currency for a rejection, produced identically by
// scanner, parser, analyzer, and vir.
type Diagnostic struct {
	Code Code
	Sev  Severity

	// Pos is where the diagnostic points. End bounds the underlined span; if
	// it is NoPos or does not exceed Pos, one column is underlined.
	Pos token.Pos
	End token.Pos

	// Message is the rendered template. This is the normative text that
	// Expected(error, "...") compares against (A.12.2).
	Message string

	// Notes are secondary spans. A use-after-transfer (A.9.2) needs two:
	// where the binding died, and where it was subsequently read.
	Notes []Note

	// Fixits are machine-applicable edits. Present from day one because the
	// spec already commits to two by name — A.1.4's `.transfer()` and A.4.8's
	// addr(x) -> &x — and threading a new field through every constructor
	// later is worse than carrying an empty slice now.
	Fixits []Fixit
}

type Note struct {
	Pos     token.Pos
	End     token.Pos
	Message string
}

// Fixit is a single replacement over [Pos, End). An empty Text is a deletion;
// an empty span is an insertion.
type Fixit struct {
	Message string
	Pos     token.Pos
	End     token.Pos
	Text    string
}

// New builds a diagnostic from a code and a span, formatting the registered
// template with args.
//
// Call sites must not format their own text. The template registry is what
// keeps one rule's message identical across every site that raises it, which
// is what makes Expected(error, "...") stable enough to be specification.
func New(code Code, pos, end token.Pos, args ...any) *Diagnostic {
	m, ok := registry[code]
	if !ok {
		return &Diagnostic{
			Code:    Internal,
			Sev:     Error,
			Pos:     pos,
			End:     end,
			Message: fmt.Sprintf("unregistered diagnostic code %s", code),
		}
	}
	msg := m.tmpl
	if len(args) > 0 {
		msg = fmt.Sprintf(m.tmpl, args...)
	}
	return &Diagnostic{Code: code, Sev: m.sev, Pos: pos, End: end, Message: msg}
}

// At is New for a zero-width position.
func At(code Code, pos token.Pos, args ...any) *Diagnostic {
	return New(code, pos, token.NoPos, args...)
}

// AtToken is New over a token's full extent.
func AtToken(code Code, t token.Token, args ...any) *Diagnostic {
	return New(code, t.Pos, t.End(), args...)
}

func (d *Diagnostic) WithNote(pos, end token.Pos, format string, args ...any) *Diagnostic {
	d.Notes = append(d.Notes, Note{Pos: pos, End: end, Message: fmt.Sprintf(format, args...)})
	return d
}

// WithFixit attaches a replacement over [pos, end).
//
//	d.WithFixit(call.Pos(), call.End(), "var "+name, "use the 'var' prefix")
func (d *Diagnostic) WithFixit(pos, end token.Pos, text, format string, args ...any) *Diagnostic {
	d.Fixits = append(d.Fixits, Fixit{
		Message: fmt.Sprintf(format, args...),
		Pos:     pos, End: end, Text: text,
	})
	return d
}

// Text returns the message alone, with no position, severity, code, notes, or
// caret. This is the string an Expected(error, "...") test compares against;
// keeping it separate from Render is what stops a source excerpt or a column
// number from leaking into a normative comparison.
func (d *Diagnostic) Text() string { return d.Message }

func (d *Diagnostic) Error() string { return d.Message }

// Span reports the underlined extent, normalizing an absent or inverted End.
func (d *Diagnostic) Span() (pos, end token.Pos) {
	if d.End.IsValid() && d.End > d.Pos {
		return d.Pos, d.End
	}
	return d.Pos, d.Pos + 1
}