package compiler

import (
	"fmt"
	"strings"
)

// Severity classifies a compiler message.
type Severity int

const (
	SevError   Severity = iota
	SevWarning
	SevNote
)

// Diagnostic is a single compiler message with an optional source position.
type Diagnostic struct {
	Pos      Pos
	Severity Severity
	Message  string
}

func (d *Diagnostic) String() string {
	label := "error"
	if d.Severity == SevWarning {
		label = "warning"
	} else if d.Severity == SevNote {
		label = "note"
	}
	if d.Pos.IsValid() {
		return fmt.Sprintf("%s:%d:%d: %s: %s",
			d.Pos.File, d.Pos.Line, d.Pos.Column, label, d.Message)
	}
	return fmt.Sprintf("%s: %s", label, d.Message)
}

// Diagnostics accumulates compiler messages.
type Diagnostics struct {
	items []*Diagnostic
}

func NewDiagnostics() *Diagnostics { return &Diagnostics{} }

func (d *Diagnostics) Reset() { d.items = nil }

func (d *Diagnostics) add(pos Pos, sev Severity, format string, args ...interface{}) {
	d.items = append(d.items, &Diagnostic{
		Pos:      pos,
		Severity: sev,
		Message:  fmt.Sprintf(format, args...),
	})
}

func (d *Diagnostics) Errorf(pos Pos, format string, args ...interface{}) {
	d.add(pos, SevError, format, args...)
}

func (d *Diagnostics) Warnf(pos Pos, format string, args ...interface{}) {
	d.add(pos, SevWarning, format, args...)
}

func (d *Diagnostics) Notef(pos Pos, format string, args ...interface{}) {
	d.add(pos, SevNote, format, args...)
}

func (d *Diagnostics) HasErrors() bool {
	for _, item := range d.items {
		if item.Severity == SevError {
			return true
		}
	}
	return false
}

// Error returns all diagnostics formatted as a single error, or nil.
func (d *Diagnostics) Error() error {
	if !d.HasErrors() {
		return nil
	}
	var sb strings.Builder
	for _, item := range d.items {
		sb.WriteString(item.String())
		sb.WriteByte('\n')
	}
	return fmt.Errorf("%s", strings.TrimRight(sb.String(), "\n"))
}

// All returns the collected diagnostics slice.
func (d *Diagnostics) All() []*Diagnostic { return d.items }