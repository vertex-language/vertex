// Package diag is the shared diagnostic currency for the Vertex toolchain.
// Scanner, parser, analyzer, and vir all produce *Diagnostic values through
// this package rather than formatting their own error strings.
//
// A rejection has one shape everywhere: a stable Code, a span, a message
// rendered from a registered template, and optionally secondary spans and
// machine-applicable edits. The template registry is what keeps one rule's
// message identical across every site that raises it — which is what makes the
// error form of grammar.md's ExpectedType, whose second operand is a message
// string, stable enough to serve as specification.
//
// Rendering is deliberately a separate concern. The normative comparison is
// always against Diagnostic.Text(), never against anything Renderer produces,
// so a source excerpt or a column number can never leak into it.
//
// diag depends on token for Pos and FileSet, and on nothing else.
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

	// Message is the rendered template. This is the normative text that the
	// error form of ExpectedType compares against.
	Message string

	// Notes are secondary spans. A use-after-transfer needs two: where the
	// binding died, and where it was subsequently read.
	Notes []Note

	// Fixits are machine-applicable edits. Present from day one because
	// grammar.md already commits to one by name — `transfer` is reserved and
	// bound to nothing so that x.transfer() diagnoses as a misspelled
	// ownership marker — and threading a new field through every constructor
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
// is what makes an expected-message comparison stable enough to be
// specification.
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

// AtToken is New over a token's full extent. Token.End() derives that extent
// from Lit where a token carries one and from the Kind's fixed spelling
// otherwise, so this is correct for keywords and punctuation too.
func AtToken(code Code, t token.Token, args ...any) *Diagnostic {
	return New(code, t.Pos, t.End(), args...)
}

// AtNode is New over any node's extent. It takes the two methods rather than an
// ast.Node so that diag does not depend on ast — the dependency runs the other
// way, and a cycle here would be structural.
func AtNode(code Code, pos, end token.Pos, args ...any) *Diagnostic {
	return New(code, pos, end, args...)
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

// WithInsert is WithFixit over an empty span.
func (d *Diagnostic) WithInsert(pos token.Pos, text, format string, args ...any) *Diagnostic {
	return d.WithFixit(pos, pos, text, format, args...)
}

// WithDelete is WithFixit with empty replacement text.
func (d *Diagnostic) WithDelete(pos, end token.Pos, format string, args ...any) *Diagnostic {
	return d.WithFixit(pos, end, "", format, args...)
}

// Text returns the message alone, with no position, severity, code, notes, or
// caret. This is the string an expected-message comparison tests against;
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