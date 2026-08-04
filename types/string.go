package types

import (
	"fmt"
	"strings"
)

// TypeString renders t in source syntax.
//
// This is what a diag template's %s receives, so it is normative text: a change
// here changes what an `Expected(error, "...")` result matches. diag's registry
// exists for the same reason, and the two must move together.
func TypeString(t Type) string {
	var b strings.Builder
	writeType(&b, t)
	return b.String()
}

func writeType(b *strings.Builder, t Type) {
	switch x := t.(type) {
	case nil:
		b.WriteString("<nil>")

	case *Basic:
		b.WriteString(x.name)

	case *Named:
		if x.obj != nil {
			b.WriteString(x.obj.name)
		}
		if len(x.typeArgs) > 0 {
			b.WriteByte('[')
			for i, a := range x.typeArgs {
				if i > 0 {
					b.WriteString(", ")
				}
				writeType(b, a)
			}
			b.WriteByte(']')
		}

	case *Abstract:
		if x.obj != nil {
			b.WriteString(x.obj.name)
			return
		}
		b.WriteString("abstract")

	case *TypeParam:
		b.WriteString(x.name)

	case *Ownership:
		b.WriteString(x.kind.String())
		b.WriteByte(' ')
		writeType(b, x.elem)

	case *Array:
		fmt.Fprintf(b, "[%d]", x.len)
		writeType(b, x.elem)

	case *Slice:
		b.WriteString("[]")
		writeType(b, x.elem)

	case *Map:
		b.WriteString("map[")
		writeType(b, x.key)
		b.WriteByte(']')
		writeType(b, x.elem)

	case *Chan:
		b.WriteString("chan ")
		writeType(b, x.elem)

	case *Pointer:
		b.WriteString("typed_ptr ")
		// A typed_ptr may not be the direct base of another, so a nested one is
		// written parenthesized and must render that way to round-trip.
		if _, nested := x.elem.(*Pointer); nested {
			b.WriteByte('(')
			writeType(b, x.elem)
			b.WriteByte(')')
			return
		}
		writeType(b, x.elem)

	case *Vector:
		b.WriteString("vector[")
		writeType(b, x.elem)
		fmt.Fprintf(b, ", %d]", x.lanes)

	case *Predicate:
		// §5.1 ⊢ a lane predicate "has no source spelling". A diagnostic still
		// has to name it, so this is a description rather than syntax, and it
		// is written so a reader cannot mistake it for something writable.
		fmt.Fprintf(b, "lane predicate of %d lanes", x.lanes)

	case *Tuple:
		writeTuple(b, x)

	case *Signature:
		writeSignature(b, x)

	case *Tensor:
		b.WriteString("tensor[")
		writeType(b, x.elem)
		for _, d := range x.shape {
			fmt.Fprintf(b, ", %d", d)
		}
		b.WriteByte(']')

	case *Struct:
		kw := "struct"
		if x.class {
			kw = "class"
		}
		b.WriteString(kw)
		b.WriteString(" {")
		for i, f := range x.fields {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(f.Name)
			b.WriteString(": ")
			writeType(b, f.Type)
		}
		b.WriteByte('}')

	case *Enum:
		b.WriteString("enum {")
		for i, v := range x.variants {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(v.Name)
		}
		b.WriteByte('}')

	default:
		fmt.Fprintf(b, "<unknown %T>", t)
	}
}

func writeTuple(b *strings.Builder, t *Tuple) {
	b.WriteByte('(')
	for i := 0; i < t.Len(); i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		v := t.At(i)
		if v.name != "" {
			b.WriteString(v.name)
			b.WriteString(": ")
		}
		if m := v.mode.String(); m != "" {
			b.WriteString(m)
			b.WriteByte(' ')
		}
		writeType(b, v.typ)
	}
	b.WriteByte(')')
}

func writeSignature(b *strings.Builder, s *Signature) {
	b.WriteString("func")
	if s.recv != nil {
		b.WriteString(" (")
		b.WriteString(s.recv.name)
		b.WriteString(": ")
		if m := s.recv.mode.String(); m != "" {
			b.WriteString(m)
			b.WriteByte(' ')
		}
		writeType(b, s.recv.typ)
		b.WriteByte(')')
	}

	b.WriteByte('(')
	for i := 0; i < s.params.Len(); i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		v := s.params.At(i)
		if v.name != "" {
			b.WriteString(v.name)
			b.WriteString(": ")
		}
		if s.variadic && i == s.params.Len()-1 {
			b.WriteString("...")
		}
		if m := v.mode.String(); m != "" {
			b.WriteString(m)
			b.WriteByte(' ')
		}
		writeType(b, v.typ)
	}
	b.WriteByte(')')

	if m := s.marker.String(); m != "" {
		b.WriteByte(' ')
		b.WriteString(m)
	}

	if s.expected != nil {
		b.WriteString(" -> ")
		writeExpected(b, s.expected)
		return
	}

	// §7.1 ⊢ omitting the result is the void form; there is no `void` type
	// name, so a void signature renders with nothing after the parameters.
	if s.results.Len() == 0 {
		return
	}
	b.WriteString(" -> ")
	if s.results.Len() == 1 && s.results.At(0).name == "" {
		writeType(b, s.results.At(0).typ)
		return
	}
	writeTuple(b, s.results)
}

func writeExpected(b *strings.Builder, e *Expected) {
	b.WriteString("Expected(")
	if e.typ != nil {
		writeType(b, e.typ)
	} else {
		b.WriteString("error")
	}
	if e.hasMsg {
		fmt.Fprintf(b, ", %q", e.msg)
	}
	b.WriteByte(')')
}

// ExpectedString renders an ExpectedType for a diagnostic. It is separate from
// TypeString because Expected is not a Type.
func ExpectedString(e *Expected) string {
	if e == nil {
		return ""
	}
	var b strings.Builder
	writeExpected(&b, e)
	return b.String()
}

// ConstraintString renders a constraint for a diagnostic. Also separate,
// because a Constraint is not a Type either.
func ConstraintString(c *Constraint) string {
	switch {
	case c == nil:
		return "<nil>"
	case c.obj != nil:
		return c.obj.name
	case c.IsAny():
		return "any"
	}

	var b strings.Builder
	for i, t := range c.terms {
		if i > 0 {
			b.WriteString(" | ")
		}
		if t.Tilde {
			b.WriteByte('~')
		}
		writeType(&b, t.Type)
	}
	for _, e := range c.embeds {
		if b.Len() > 0 {
			b.WriteString("; ")
		}
		b.WriteString(ConstraintString(e))
	}
	for _, m := range c.methods {
		if b.Len() > 0 {
			b.WriteString("; ")
		}
		b.WriteString("func ")
		b.WriteString(m.name)
		b.WriteString(strings.TrimPrefix(TypeString(m.typ), "func"))
	}
	return b.String()
}

// ObjectString renders an object for a diagnostic: its kind, name, and type.
func ObjectString(obj Object) string {
	switch o := obj.(type) {
	case *Var:
		return fmt.Sprintf("variable %s of type %s", o.name, TypeString(o.typ))
	case *Const:
		return fmt.Sprintf("constant %s of type %s", o.name, TypeString(o.typ))
	case *Func:
		return fmt.Sprintf("func %s%s", o.name, strings.TrimPrefix(TypeString(o.typ), "func"))
	case *TypeName:
		if o.IsConstraint() {
			return fmt.Sprintf("constraint %s", o.name)
		}
		return fmt.Sprintf("type %s", o.name)
	case *Builtin:
		return fmt.Sprintf("builtin %s", o.name)
	case *PkgName:
		return fmt.Sprintf("package %s", o.name)
	}
	if obj == nil {
		return "<nil>"
	}
	return obj.Name()
}

func (o *Ownership) String() string { return TypeString(o) }
func (a *Array) String() string     { return TypeString(a) }
func (s *Slice) String() string     { return TypeString(s) }
func (m *Map) String() string       { return TypeString(m) }
func (t *Tuple) String() string     { return TypeString(t) }
func (c *Chan) String() string      { return TypeString(c) }
func (p *Pointer) String() string   { return TypeString(p) }
func (v *Vector) String() string    { return TypeString(v) }
func (p *Predicate) String() string { return TypeString(p) }
func (s *Signature) String() string { return TypeString(s) }
func (t *Tensor) String() string    { return TypeString(t) }
func (n *Named) String() string     { return TypeString(n) }
func (s *Struct) String() string    { return TypeString(s) }
func (e *Enum) String() string      { return TypeString(e) }
func (a *Abstract) String() string  { return TypeString(a) }
func (t *TypeParam) String() string { return TypeString(t) }
func (e *Expected) String() string  { return ExpectedString(e) }
func (c *Constraint) String() string { return ConstraintString(c) }