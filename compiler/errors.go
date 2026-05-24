package compiler

import (
	"fmt"
	"strings"
)

// Pos is a source position.
type Pos struct {
	File string
	Line int
	Col  int
}

func (p Pos) String() string {
	if p.File == "" {
		return fmt.Sprintf("%d:%d", p.Line, p.Col)
	}
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Col)
}

// CompileError is a single diagnostic.
type CompileError struct {
	Pos     Pos
	Message string
}

func (e *CompileError) Error() string { return fmt.Sprintf("%s: %s", e.Pos, e.Message) }

// ErrorList accumulates diagnostics.
type ErrorList []*CompileError

func (el ErrorList) Error() string {
	ss := make([]string, len(el))
	for i, e := range el {
		ss[i] = e.Error()
	}
	return strings.Join(ss, "\n")
}

func (el *ErrorList) add(pos Pos, format string, args ...any) {
	*el = append(*el, &CompileError{Pos: pos, Message: fmt.Sprintf(format, args...)})
}

func (el ErrorList) err() error {
	if len(el) == 0 {
		return nil
	}
	return el
}