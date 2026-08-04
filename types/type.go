// Package types defines the Vertex type representation, the compile-time
// constant representation, and the predicates over both that the analyzer and
// lower share.
//
// It is deliberately not the analyzer. §3.4 ⊢ every type has a size known at
// compile time, and lower needs sizes, alignments, and §8's ownership
// discipline without needing the checker that produced them. Keeping the
// representation in its own package is what makes that possible.
//
// ast records shape; this package records meaning. The dependency runs one
// way: types imports ast (Info's side tables are keyed by ast nodes) and token
// (positions, and the one home of the reserved-builtin set). ast never imports
// types, and no analysis result is ever written back onto a syntax node.
//
// Citation convention: a bare § is semantics.md, which owns every rule
// grammar.md defers to as a "static rule". CamelCase names are grammar.md
// productions. Where both documents are silent — enum layout, closure
// representation, the vector ABI — the comment says so rather than dressing an
// implementation choice as a rule.
package types

// Type is the representation of a Vertex type.
//
// Underlying strips exactly one layer of naming: a Named's underlying is what
// its declaration named, and every other type is its own underlying. §3.1's
// `~T` type-set term is defined over this, and §3.1's transparent alias
// resolves before a Type is ever constructed, so an alias never appears here.
type Type interface {
	Underlying() Type
	String() string
}

// Two things in this package are deliberately NOT Types:
//
//   - Constraint (constraint.go). §9 ⊢ a ConstraintDecl "is never a value
//     type"; §3.2 ⊢ a constraint name is legal "only in a `[`…`]` position".
//     Keeping it out of the Type hierarchy makes `var c: Ordered` a shape error
//     the checker cannot forget to raise, instead of a predicate someone has to
//     remember to call.
//   - Expected (composite.go). §3.2 ⊢ `Expected(…)` is legal "only as a
//     FunctionDecl/MethodDecl result, in a `build test` file". It hangs off
//     Signature rather than sitting in the result Tuple, so it cannot reach a
//     field, a binding, or a FunctionType.
//
// A range is a third case with no representation at all: §5.2 ⊢ `a..b` "has no
// type and cannot be bound, returned, passed, or stored". It is carried as
// RangeMode in info.go, never as a Type.

// Mode is a parameter or receiver access convention.
//
// §3.2 ⊢ `mut T` and `var T` are legal in "a parameter or receiver only", while
// `unique T`, `shared T`, and `weak T` are legal "anywhere a Type is". Two
// different things wearing one syntactic hat, so they get two different
// representations: the first pair is a Mode on a Var, the second is an
// *Ownership.
//
// The split is what makes ReceiverType's three-way choice fall out.
// `(x: shared T)` is a ModeNone Var over an *Ownership; `(x: mut T)` is a
// ModeMut Var over a bare T. §3.2's "qualifiers do not stack" then needs no
// special case in this package at all — a Var carries one Mode, so `mut var T`
// is unrepresentable, and `mut shared T` is representable and rejected by the
// analyzer against the stacking rule.
type Mode uint8

const (
	ModeNone Mode = iota // shared, read-only, non-owning
	ModeMut              // exclusive; §7.2 ⊢ the caller's binding must be `var`
	ModeVar              // owning; copy-or-move decided at the call site (§8.2)
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

// Marker is a FunctionMarker. It lives on Signature because §7.4 ⊢ "the marker
// is part of the function's type, so it is checked at the declaration and again
// at every call. The marker must agree at both ends."
//
// That is also why §4.1 ⊢ "a `func(int32)` is not assignable to a
// `func(int32) async`" needs no rule of its own here: Identical compares
// markers, and assignability reduces to identity.
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