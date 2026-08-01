package diag

import (
	"bytes"
	"fmt"
	"io"

	"github.com/vertex-language/vertex/token"
)

// Renderer produces human-facing output. It is deliberately separate from
// Diagnostic: the normative comparison in A.12.2 is against Text(), never
// against anything this file produces.
type Renderer struct {
	Fset *token.FileSet

	// Source supplies a file's bytes for excerpting. May be nil, in which case
	// diagnostics render without a source line or caret.
	Source func(*token.File) []byte

	// Color enables ANSI escapes.
	Color bool
}

const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiRed    = "\x1b[1;31m"
	ansiYellow = "\x1b[1;33m"
	ansiCyan   = "\x1b[1;36m"
)

func (r *Renderer) color(code, s string) string {
	if !r.Color {
		return s
	}
	return code + s + ansiReset
}

func (r *Renderer) Render(w io.Writer, d *Diagnostic) error {
	sevColor := ansiRed
	if d.Sev == Warning {
		sevColor = ansiYellow
	}

	head := fmt.Sprintf("%s[%s]", d.Sev, d.Code)
	if _, err := fmt.Fprintf(w, "%s: %s\n",
		r.color(sevColor, head), r.color(ansiBold, d.Message)); err != nil {
		return err
	}

	pos := r.Fset.Position(d.Pos)
	if pos.IsValid() {
		if _, err := fmt.Fprintf(w, "  --> %s\n", pos); err != nil {
			return err
		}
	}

	start, end := d.Span()
	if err := r.excerpt(w, start, end, sevColor); err != nil {
		return err
	}

	for _, n := range d.Notes {
		np := r.Fset.Position(n.Pos)
		if _, err := fmt.Fprintf(w, "  %s %s\n", r.color(ansiCyan, "note:"), n.Message); err != nil {
			return err
		}
		if np.IsValid() {
			if _, err := fmt.Fprintf(w, "  --> %s\n", np); err != nil {
				return err
			}
		}
		nEnd := n.End
		if !nEnd.IsValid() || nEnd <= n.Pos {
			nEnd = n.Pos + 1
		}
		if err := r.excerpt(w, n.Pos, nEnd, ansiCyan); err != nil {
			return err
		}
	}

	for _, f := range d.Fixits {
		suggestion := f.Text
		if suggestion == "" {
			suggestion = "(delete)"
		}
		if _, err := fmt.Fprintf(w, "  %s %s: %s\n",
			r.color(ansiCyan, "help:"), f.Message, suggestion); err != nil {
			return err
		}
	}

	return nil
}

func (r *Renderer) RenderList(w io.Writer, l *List) error {
	for i, d := range l.Items() {
		if i > 0 {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
		if err := r.Render(w, d); err != nil {
			return err
		}
	}
	return nil
}

// excerpt writes the source line containing start, with a caret run under
// [start, end) clipped to that line.
func (r *Renderer) excerpt(w io.Writer, start, end token.Pos, hl string) error {
	if r.Source == nil || r.Fset == nil {
		return nil
	}
	f := r.Fset.File(start)
	if f == nil {
		return nil
	}
	src := r.Source(f)
	if src == nil {
		return nil
	}

	pos := f.Position(start)
	if !pos.IsValid() || pos.Offset > len(src) {
		return nil
	}

	lineStart := pos.Offset - (pos.Column - 1)
	if lineStart < 0 {
		lineStart = 0
	}
	lineEnd := lineStart
	for lineEnd < len(src) && src[lineEnd] != '\n' && src[lineEnd] != '\r' {
		lineEnd++
	}
	line := src[lineStart:lineEnd]

	width := 1
	if endOff := int(end) - f.Base(); endOff > pos.Offset {
		width = endOff - pos.Offset
		if pos.Offset+width > lineEnd {
			width = lineEnd - pos.Offset // clip a multi-line span to this line
		}
		if width < 1 {
			width = 1
		}
	}

	gutter := fmt.Sprintf("%d", pos.Line)
	pad := bytes.Repeat([]byte{' '}, len(gutter))

	if _, err := fmt.Fprintf(w, "%s | %s\n", gutter, line); err != nil {
		return err
	}

	// Build the caret run by echoing the line prefix with tabs preserved and
	// everything else blanked. This aligns under any tab width without the
	// renderer needing to know one.
	prefix := line[:pos.Offset-lineStart]
	indent := make([]byte, len(prefix))
	for i, c := range prefix {
		if c == '\t' {
			indent[i] = '\t'
		} else {
			indent[i] = ' '
		}
	}
	carets := bytes.Repeat([]byte{'^'}, width)

	_, err := fmt.Fprintf(w, "%s | %s%s\n", pad, indent, r.color(hl, string(carets)))
	return err
}