package types

// ---------------------------------------------------------------- ownership

type OwnKind uint8

const (
	Unique OwnKind = iota
	Shared
	Weak
)

func (k OwnKind) String() string {
	switch k {
	case Shared:
		return "shared"
	case Weak:
		return "weak"
	}
	return "unique"
}

// Ownership is `unique T`, `shared T`, or `weak T` — §3.2 ⊢ legal "anywhere a
// Type is", which is what separates them from `mut` and `var` (a Mode on a Var;
// see type.go).
//
// grammar.md ⊢ these three "are keywords, not reserved builtin names; they get
// the HeapConstructor production" — so they are not in Universe's builtin set,
// and scope.go must not put them there.
type Ownership struct {
	kind OwnKind
	elem Type
}

func NewOwnership(k OwnKind, elem Type) *Ownership { return &Ownership{k, elem} }

func (o *Ownership) Kind() OwnKind    { return o.kind }
func (o *Ownership) Elem() Type       { return o.elem }
func (o *Ownership) Underlying() Type { return o }

// ------------------------------------------------------------- array, slice

// Array is `[N]T`. §5.3 ⊢ an ArrayLength "must be a non-negative integer" and a
// constant, so it is resolved to an int64 before construction rather than kept
// as an expression. §3.1 ⊢ "`[N]T` carries `N` in its identity; `[8]int32` and
// `[16]int32` are unrelated."
type Array struct {
	len  int64
	elem Type
}

func NewArray(elem Type, n int64) *Array { return &Array{n, elem} }

func (a *Array) Len() int64       { return a.len }
func (a *Array) Elem() Type       { return a.elem }
func (a *Array) Underlying() Type { return a }

// Slice is `[]T`. §3.4 lists it among the indirections that are "all one word",
// which is what makes it usable to break a recursive type.
type Slice struct{ elem Type }

func NewSlice(elem Type) *Slice   { return &Slice{elem} }
func (s *Slice) Elem() Type       { return s.elem }
func (s *Slice) Underlying() Type { return s }

// Map is `map[K]V`. K must satisfy `comparable` (IsComparable, predicates.go);
// the analyzer checks it at construction and raises NonComparableKey.
type Map struct{ key, elem Type }

func NewMap(key, elem Type) *Map { return &Map{key, elem} }

func (m *Map) Key() Type        { return m.key }
func (m *Map) Elem() Type       { return m.elem }
func (m *Map) Underlying() Type { return m }

// ------------------------------------------------------------------- tuple

// Tuple is a TupleType, and also a Signature's result list.
//
// Elements may be named — TupleElem is `[ identifier ":" ] Type` — so a Tuple
// holds Vars rather than bare Types. An unnamed element is a Var with an empty
// name; §3.1 ⊢ parameter names are not part of a func type, and the same holds
// for a tuple element's name, so Identical ignores them.
//
// grammar.md ⊢ "a tuple has at least one element; there is no unit type." An
// empty Tuple is therefore never a TupleType. It appears in exactly one place —
// a Signature's results for the void form — and Signature.IsVoid is how to ask
// about it. Nothing should hand an empty Tuple to a position expecting a Type.
type Tuple struct{ vars []*Var }

func NewTuple(x ...*Var) *Tuple { return &Tuple{x} }

func (t *Tuple) Len() int {
	if t == nil {
		return 0
	}
	return len(t.vars)
}

func (t *Tuple) At(i int) *Var    { return t.vars[i] }
func (t *Tuple) Underlying() Type { return t }

// ---------------------------------------------------------------- channel

// Chan is `chan T`. §3.4 makes it one word; §3.3 ⊢ its zero value is "a closed,
// empty channel", so a declaration without an initializer is still well
// defined and there is no definite-assignment analysis to write.
//
// A channel type carries no direction (grammar.md, *Channel and pointer
// types*). §7.4 ⊢ `thread` and `async` "hand back a `chan T`" — a property of
// the launch expression's result, not a second type.
type Chan struct{ elem Type }

func NewChan(elem Type) *Chan    { return &Chan{elem} }
func (c *Chan) Elem() Type       { return c.elem }
func (c *Chan) Underlying() Type { return c }

// ---------------------------------------------------------------- pointer

// Pointer is `typed_ptr T`: §10's third tier, where "`typed_ptr` misuse only"
// is undefined behaviour and §5.5 ⊢ operations "check nothing".
//
// §8.4 ⊢ it is "the one type these rules do not reach. Two copies are two
// unchecked aliases, and exclusivity there is convention rather than proof."
// §3.2 ⊢ it may never be the direct base of another PointerType and never a
// receiver type; both are static rules over a declaration, not shapes this type
// records.
type Pointer struct{ elem Type }

func NewPointer(elem Type) *Pointer { return &Pointer{elem} }
func (p *Pointer) Elem() Type       { return p.elem }
func (p *Pointer) Underlying() Type { return p }

// ----------------------------------------------------------------- vector

// Vector is `vector[T, N]`. §5.3 ⊢ the lane count is one of the three positions
// requiring a bare literal token, so it is an int64 here and never an
// expression.
//
// §3.2 ⊢ legal "anywhere except a `gpu`/`npu` body or signature, a foreign
// boundary, or a map key" — all static rules over a position, none of them a
// property of this shape.
type Vector struct {
	elem  Type
	lanes int64
}

func NewVector(elem Type, lanes int64) *Vector { return &Vector{elem, lanes} }

func (v *Vector) Elem() Type       { return v.elem }
func (v *Vector) Lanes() int64     { return v.lanes }
func (v *Vector) Underlying() Type { return v }

// Predicate is the result of a vector comparison.
//
// §5.1 ⊢ it "has no source spelling and may not be an `if` condition, a `&&`
// operand, a field, or a channel element." It gets a Type anyway because the
// comparison has to evaluate to something the checker can carry to `blend`; the
// prohibitions are read off IsPredicate at each of those four positions.
//
// Because it cannot be a field or a channel element, it never reaches storage,
// which is why Sizeof answers zero rather than inventing a lane-mask layout the
// sources do not fix.
type Predicate struct{ lanes int64 }

func NewPredicate(lanes int64) *Predicate { return &Predicate{lanes} }

func (p *Predicate) Lanes() int64      { return p.lanes }
func (p *Predicate) Underlying() Type  { return p }

// --------------------------------------------------------------- signature

// Signature is a FunctionType and the type of every function, method,
// initializer, deinitializer, function literal, foreign declaration, and
// MethodRequirement.
//
// Recv is nil for a free function. Marker is part of the type because §7.4
// makes it part of the callee's contract, checked at both ends.
//
// results is empty for the void form: §7.1 ⊢ "omitting the result is the void
// form; there is no `void` type and no unit type." There is deliberately no
// Unit value in this package to reach for instead.
//
// expected is the ExpectedType result and is non-nil only where results is
// empty. It is a separate field rather than an element of results because §3.2
// ⊢ `Expected(…)` is legal "only as a FunctionDecl/MethodDecl result" — putting
// it in the tuple would make it reachable wherever a Type is.
type Signature struct {
	recv     *Var // nil for a free function; carries the receiver's Mode
	params   *Tuple
	results  *Tuple
	variadic bool
	marker   Marker
	expected *Expected
}

func NewSignature(recv *Var, params, results *Tuple, variadic bool, m Marker) *Signature {
	if params == nil {
		params = NewTuple()
	}
	if results == nil {
		results = NewTuple()
	}
	return &Signature{recv: recv, params: params, results: results, variadic: variadic, marker: m}
}

func (s *Signature) Recv() *Var       { return s.recv }
func (s *Signature) Params() *Tuple   { return s.params }
func (s *Signature) Results() *Tuple  { return s.results }
func (s *Signature) Variadic() bool   { return s.variadic }
func (s *Signature) Marker() Marker   { return s.marker }
func (s *Signature) Underlying() Type { return s }

// IsVoid reports the void form: no result slot at all. This is the question to
// ask; there is no unit type to compare against.
func (s *Signature) IsVoid() bool { return s.results.Len() == 0 && s.expected == nil }

func (s *Signature) Expected() *Expected { return s.expected }

// SetExpected attaches an ExpectedType result. It is a setter because the
// analyzer resolves a signature's shape before it can know the file's build tag
// licenses one, and because §7.4's test-function shape is checked afterwards.
func (s *Signature) SetExpected(e *Expected) { s.expected = e }

// ---------------------------------------------------------------- expected

// Expected is an ExpectedType: the result form of a test function.
//
// It is not a Type, for the same reason Constraint is not: §3.2 ⊢ it is legal
// "only as a FunctionDecl/MethodDecl result, in a `build test` file", and §1.2
// makes `build test` a requirement rather than a permission. Keeping it out of
// the Type hierarchy is what makes
// `var f: func() -> Expected(int32, "5")` unrepresentable rather than merely
// rejected.
//
// Msg is normative text. It is compared against a Diagnostic's Text(), which is
// diag's whole reason for a template registry — so a change to a message
// template changes what an Expected result matches.
type Expected struct {
	typ    Type   // nil for the `error` form
	msg    string // "" when omitted, which only the error form permits
	hasMsg bool
}

// NewExpectedValue is `Expected(TypeName, string_lit)`. Both operands are
// written; there is no message-free value form.
func NewExpectedValue(t Type, msg string) *Expected {
	return &Expected{typ: t, msg: msg, hasMsg: true}
}

// NewExpectedError is `Expected(error [, string_lit])`. `error` is an ordinary
// identifier recognized only in this production, which is why it leaves no
// TypeName behind and Type answers nil.
func NewExpectedError(msg string, hasMsg bool) *Expected {
	return &Expected{msg: msg, hasMsg: hasMsg}
}

func (e *Expected) IsError() bool { return e != nil && e.typ == nil }
func (e *Expected) Type() Type    { return e.typ }
func (e *Expected) Msg() string   { return e.msg }
func (e *Expected) HasMsg() bool  { return e.hasMsg }

// ----------------------------------------------------------------- tensor

// Tensor is `tensor[T, dims...]`. §3.2 ⊢ legal "inside an `npu` body or that
// function's own signature", which is a static rule over a position.
//
// grammar.md's ShapeList is `int_lit { "," int_lit }` and §5.3 makes each entry
// one of the three bare-literal-token positions, so the shape is resolved to
// int64 here and a folded `-1` cannot reach it.
type Tensor struct {
	elem  Type
	shape []int64
}

func NewTensor(elem Type, shape []int64) *Tensor { return &Tensor{elem, shape} }

func (t *Tensor) Elem() Type       { return t.elem }
func (t *Tensor) Shape() []int64   { return t.shape }
func (t *Tensor) Rank() int        { return len(t.shape) }
func (t *Tensor) Underlying() Type { return t }

// Elems is the total element count, which is what Sizeof multiplies by the
// element width and what §4.3's `[N]T ↔ tensor[T, N]` launch conversion matches
// against an array length.
func (t *Tensor) Elems() int64 {
	n := int64(1)
	for _, d := range t.shape {
		n *= d
	}
	return n
}