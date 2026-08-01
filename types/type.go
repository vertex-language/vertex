// Package types defines the Vertex type representation and the predicates over
// it that the analyzer and lower both need.
//
// It is deliberately not part of the analyzer. A.15's invariant is that "every
// value has a statically known layout, and every cost is decided at compile
// time" — which means lower needs sizes, alignments, and ownership discipline
// without needing the checker that produced them. Keeping the representation in
// its own package is what makes that possible.
//
// ast records shape; this package records meaning. The dependency runs one way:
// types may import ast (Info's maps are keyed by ast nodes), ast may never
// import types. No analysis result is ever written back onto a syntax node.
package types

// Type is the representation of a Vertex type.
//
// Underlying strips exactly one layer of naming: a Named's underlying is what
// it was declared as, and every other type is its own underlying. A.7.3's `~T`
// is defined in terms of this, and A.6.6's transparent alias resolves before a
// Type is ever constructed, so an alias never appears here.
type Type interface {
	Underlying() Type
	String() string
}

// Constraint is deliberately NOT a Type. A.7.2 ⊢ "a constraint is never a value
// type and is legal only in a [...] position", and A.14 lists `var c: Ordered`
// among the rejected forms. Making it a Type would put the rejection in a
// predicate instead of in the type system's shape.
//
// See constraint.go.

// Mode is a parameter or receiver access convention (A.3.2).
//
// `unique`, `shared`, and `weak` are real types — A.3.2 says they "may appear
// anywhere a Type may, including as type arguments". `mut` and `var` are not:
// they are "legal only in a parameter or receiver position". Two different
// things wearing one syntactic hat, so they get two different representations.
//
// The stacking rules of A.3.2 then fall out rather than needing a special case:
// `mut shared T` is a ModeMut Var whose type is *Ownership, and `mut var T` is
// unrepresentable because a Var carries one Mode.
type Mode uint8

const (
	ModeNone Mode = iota // by-value, non-owning
	ModeMut              // exclusive, non-owning, mutating; lowers to a pointer
	ModeVar              // owning; copy-or-move decided at the call site (A.4.6)
)

func (m Mode) String() string {
	switch m {
	case ModeMut:
		return "mut"
	case ModeVar:
		return "var"
	}
	return ""
}

// Marker is a FunctionMarker (A.6.1). It lives on Signature because A.4.2 makes
// it part of the callee's contract: "a launch prefix is legal only when the
// callee carries the matching function marker", checked at both ends.
type Marker uint8

const (
	MarkerNone Marker = iota
	MarkerAsync
	MarkerGPU
	MarkerNPU
	MarkerTest
)

func (m Marker) String() string {
	switch m {
	case MarkerAsync:
		return "async"
	case MarkerGPU:
		return "gpu"
	case MarkerNPU:
		return "npu"
	case MarkerTest:
		return "test"
	}
	return ""
}