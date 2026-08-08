package ast

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/vertex-language/vertex/token"
)

// Fdump writes a debug rendering of n to w.
//
// The format is not stable and is not part of this package's contract (§5.5).
// §7's golden tests consume it anyway — that is fine, because those goldens
// are regenerated, not hand-written, and a format change shows up as a diff
// across every fixture at once rather than as a silent behavior change.
//
// It is public rather than test-only because §7's tests need it and a
// test-only dumper drifts from the node set it is meant to cover.
//
// Output is deterministic: struct fields render in declaration order and slices
// in source order. Nothing here iterates a Go map (§7).
func Fdump(w io.Writer, n Node) error {
	d := &dumper{w: w}
	d.node(reflect.ValueOf(n), 0)
	return d.err
}

type dumper struct {
	w   io.Writer
	err error
}

func (d *dumper) printf(format string, args ...any) {
	if d.err != nil {
		return
	}
	_, d.err = fmt.Fprintf(d.w, format, args...)
}

var (
	posType  = reflect.TypeOf(token.NoPos)
	nodeType = reflect.TypeOf((*Node)(nil)).Elem()
)

func (d *dumper) node(v reflect.Value, depth int) {
	if d.err != nil {
		return
	}
	ind := strings.Repeat("  ", depth)

	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			d.printf("nil")
			return
		}
		d.node(v.Elem(), depth)

	case reflect.Ptr:
		if v.IsNil() {
			d.printf("nil")
			return
		}
		if n, ok := v.Interface().(Node); ok {
			d.printf("%s {  // %d-%d\n", v.Type().Elem().Name(), n.Pos(), n.End())
		} else {
			d.printf("%s {\n", v.Type().Elem().Name())
		}
		d.fields(v.Elem(), depth+1)
		d.printf("%s}", ind)

	case reflect.Struct:
		d.printf("%s {\n", v.Type().Name())
		d.fields(v, depth+1)
		d.printf("%s}", ind)

	case reflect.Slice:
		if v.IsNil() {
			d.printf("nil")
			return
		}
		if v.Len() == 0 {
			d.printf("[]")
			return
		}
		d.printf("[\n")
		for i := 0; i < v.Len(); i++ {
			d.printf("%s  ", ind)
			d.node(v.Index(i), depth+1)
			d.printf("\n")
		}
		d.printf("%s]", ind)

	default:
		if v.Type() == posType {
			p := token.Pos(v.Uint())
			if p == token.NoPos {
				d.printf("-")
			} else {
				d.printf("%d", p)
			}
			return
		}
		d.printf("%v", v.Interface())
	}
}

func (d *dumper) fields(v reflect.Value, depth int) {
	if v.Kind() != reflect.Struct {
		return
	}
	ind := strings.Repeat("  ", depth)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue // unexported
		}
		d.printf("%s%s: ", ind, f.Name)
		d.node(v.Field(i), depth)
		d.printf("\n")
	}
}