package ast

import "reflect"

// isNil reports whether n is a nil interface or a nil pointer inside a
// non-nil interface.
//
// This exists because optional children are typed pointers assigned into
// interface-typed fields: a *TypeExpr field holding (*TypeRef)(nil) is not ==
// nil, and a span computed from it would panic. Every optional child in this
// package is checked through here.
//
// The reflect call is confined to span computation and Walk, both of which run
// once per node at most. If it ever shows up in a profile, replace it with a
// type switch over the concrete node types — the set is closed.
func isNil(n Node) bool {
	if n == nil {
		return true
	}
	v := reflect.ValueOf(n)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map:
		return v.IsNil()
	}
	return false
}