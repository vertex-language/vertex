package types

// Unit is `()` (A.3.1): zero bytes, one value. It is an empty Tuple rather than
// its own kind, so a void return and a unit value are the same object and
// nothing has to special-case one against the other.
var Unit = NewTuple()

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

// Ownership is `unique T`, `shared T`, or `weak T` (A.3.2) — all "ordinary
// one-word value types" that "may appear anywhere a Type may". `mut` and `var`
// are not here; they are a Mode on a Var (see type.go).
//
// A.3.2 ⊢ there is no unique→weak path: a unique block carries no control word
// for a weak reference to inspect. That is a predicate, not a shape, so it
// lives in predicates.go.
type Ownership struct {
	kind OwnKind
	elem Type
}

func NewOwnership(k OwnKind, elem Type) *Ownership { return &Ownership{k, elem} }

func (o *Ownership) Kind() OwnKind   { return o.kind }
func (o *Ownership) Elem() Type      { return o.elem }
func (o *Ownership) Underlying() Type { return o }

// ------------------------------------------------------------- array, slice

// Array is `[N]T` (A.3.1): inline storage of N × sizeof(T) with no header and
// no pointer. Len is resolved from the ArrayLength constant before construction.
type Array struct {
	len  int64
	elem Type
}

func NewArray(elem Type, n int64) *Array { return &Array{n, elem} }

func (a *Array) Len() int64        { return a.len }
func (a *Array) Elem() Type        { return a.elem }
func (a *Array) Underlying() Type  { return a }

// Slice is `[]T` (A.3.1): a three-word {ptr, len, cap} header over an
// implicitly heap-allocated block — the one implicit allocation in the language.
type Slice struct{ elem Type }

func NewSlice(elem Type) *Slice   { return &Slice{elem} }
func (s *Slice) Elem() Type       { return s.elem }
func (s *Slice) Underlying() Type { return s }

// Map is `map[K]V` (A.3.1). A.3.1 ⊢ K must satisfy `comparable`, checked at
// construction by the analyzer, not here.
type Map struct{ key, elem Type }

func NewMap(key, elem Type) *Map { return &Map{key, elem} }

func (m *Map) Key() Type        { return m.key }
func (m *Map) Elem() Type       { return m.elem }
func (m *Map) Underlying() Type { return m }

// ------------------------------------------------------------------- tuple

// Tuple is a TupleType (A.3.1) and also a Signature's result list.
//
// Elements may be named — A.3.1's TupleElement is `Type | Identifier : Type` —
// so a Tuple holds Vars rather than bare Types. An unnamed element is a Var
// with an empty name.
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

// IsUnit reports whether t is `()` — zero bytes, one value (A.3.1).
func (t *Tuple) IsUnit() bool { return t.Len() == 0 }

// ---------------------------------------------------------------- channel

// Chan is `chan T` (A.3.5): an implicitly heap-resident refcounted handle.
// Copying it bumps the count and is never a deep copy of buffered contents.
//
// There is no direction on a channel type. A.4.2 says a `thread` or `async`
// prefix "evaluates to a receive-only chan T", but that is a property of the
// launch expression's result, not a second type — the whole point of A.10.1 is
// that both sigils reduce to the same handle so no adapter is needed.
type Chan struct{ elem Type }

func NewChan(elem Type) *Chan   { return &Chan{elem} }
func (c *Chan) Elem() Type      { return c.elem }
func (c *Chan) Underlying() Type { return c }

// ---------------------------------------------------------------- pointer

// Pointer is `typed_ptr T` (A.3.3): the raw, last-resort pointer. No ownership
// tracking, no refcount, no teardown ever emitted. A.9.3 ⊢ it is the one type
// where exclusivity is a convention rather than a proof.
type Pointer struct{ elem Type }

func NewPointer(elem Type) *Pointer { return &Pointer{elem} }
func (p *Pointer) Elem() Type       { return p.elem }
func (p *Pointer) Underlying() Type { return p }

// --------------------------------------------------------------- signature

// Signature is a FunctionType (A.3.4) and the type of every function, method,
// initializer, and function literal.
//
// Recv is nil for a free function. Marker is part of the type because A.4.2
// makes it part of the callee's contract, checked at both the definition and
// the launch site.
type Signature struct {
	recv     *Var // nil for a free function; carries the receiver's Mode
	params   *Tuple
	results  *Tuple // Unit for the void form — A.3.4 has no `void` type name
	variadic bool
	marker   Marker
}

func NewSignature(recv *Var, params, results *Tuple, variadic bool, m Marker) *Signature {
	if results == nil {
		results = Unit
	}
	return &Signature{recv, params, results, variadic, m}
}

func (s *Signature) Recv() *Var        { return s.recv }
func (s *Signature) Params() *Tuple    { return s.params }
func (s *Signature) Results() *Tuple   { return s.results }
func (s *Signature) Variadic() bool    { return s.variadic }
func (s *Signature) Marker() Marker    { return s.marker }
func (s *Signature) Underlying() Type  { return s }

// Capturing reports whether a value of this signature type may carry an
// environment. It is unknowable from the type alone — A.3.4 ⊢ "a non-capturing
// function value is one word; a capturing closure is two words {code, env}" —
// so Sizeof is conservative and A.8.6's boundary check is a property of the
// *expression*, decided by the analyzer, not of the type.

// ----------------------------------------------------------------- tensor

// Tensor is `tensor[T, dims...]` (A.3.5), grammatical only under [+Npu].
//
// Shape entries are compile-time integer literals, so they are resolved to
// int64 here rather than kept as expressions. A.3.5 ⊢ signature-eligible
// element types are float32 and int8; the narrow body-only kinds (bf16, fp8,
// int4) are legal on a local binding only, which is a static rule over a
// declaration and not a distinction this type records.
type Tensor struct {
	elem  Type
	shape []int64
}

func NewTensor(elem Type, shape []int64) *Tensor { return &Tensor{elem, shape} }

func (t *Tensor) Elem() Type       { return t.elem }
func (t *Tensor) Shape() []int64   { return t.shape }
func (t *Tensor) Rank() int        { return len(t.shape) }
func (t *Tensor) Underlying() Type { return t }