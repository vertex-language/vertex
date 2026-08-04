package types

// Underlying strips one layer of naming. §3.1's `~T` is defined over it, and it
// is what every structural predicate below reaches for first.
func Underlying(t Type) Type {
	if t == nil {
		return Typ[Invalid]
	}
	return t.Underlying()
}

// AsBasic returns t's underlying *Basic, or nil.
func AsBasic(t Type) *Basic {
	b, _ := Underlying(t).(*Basic)
	return b
}

func AsNamed(t Type) *Named {
	n, _ := t.(*Named)
	return n
}

func AsStruct(t Type) *Struct {
	s, _ := Underlying(t).(*Struct)
	return s
}

func AsEnum(t Type) *Enum {
	e, _ := Underlying(t).(*Enum)
	return e
}

func AsSignature(t Type) *Signature {
	s, _ := Underlying(t).(*Signature)
	return s
}

func AsVector(t Type) *Vector {
	v, _ := Underlying(t).(*Vector)
	return v
}

func AsTensor(t Type) *Tensor {
	x, _ := Underlying(t).(*Tensor)
	return x
}

func IsInvalid(t Type) bool {
	b := AsBasic(t)
	return b != nil && b.kind == Invalid
}

// IsUntyped does not go through Underlying: an untyped constant never reaches a
// Named, so a non-*Basic answers false via is()'s nil guard.
func IsUntyped(t Type) bool {
	b, _ := t.(*Basic)
	return b.is(InfoUntyped)
}

func IsBool(t Type) bool    { return AsBasic(t).is(InfoBoolean) }
func IsInteger(t Type) bool { return AsBasic(t).is(InfoInteger) }
func IsFloat(t Type) bool   { return AsBasic(t).is(InfoFloat) }
func IsNumeric(t Type) bool { return AsBasic(t).is(InfoNumeric) }
func IsString(t Type) bool  { return AsBasic(t).is(InfoString) }
func IsChar(t Type) bool    { return AsBasic(t).is(InfoChar) }

// IsOrdered backs `< <= > >=`. §3.5 ⊢ Ordered is "numerics and `string`", so
// `char` answers false: it is comparable but not ordered.
func IsOrdered(t Type) bool { return AsBasic(t).is(InfoOrdered) }

// IsTensorElem reports one of §2.3's tensor element type names, which are
// "legal only inside an `npu` body" and, per grammar.md, never in a signature.
func IsTensorElem(t Type) bool { return AsBasic(t).is(InfoTensorElem) }

// IsPredicate reports a lane predicate. §5.1 ⊢ it "may not be an `if`
// condition, a `&&` operand, a field, or a channel element"; those four
// rejections read this.
func IsPredicate(t Type) bool {
	_, ok := Underlying(t).(*Predicate)
	return ok
}

// IsIndirection reports whether t breaks a recursive type.
//
// §3.4 ⊢ "a type may not contain itself by value, directly or through a cycle
// of struct/class/array/tuple fields... Break the cycle with an indirection,
// all of which are one word: `unique T`, `shared T`, `weak T`, `typed_ptr T`,
// `[]T`, `map[K]V`, or `chan T`." That list is exactly this predicate, and it is
// also exactly the set sizes.go gives one word.
func IsIndirection(t Type) bool {
	switch Underlying(t).(type) {
	case *Ownership, *Pointer, *Slice, *Map, *Chan:
		return true
	}
	return false
}

// IsIdentityComparable reports whether `===` and `!==` apply.
//
// §3.5 ⊢ `===` "asks whether two bindings name the same object and is legal on
// classes only", and ⊢ "`==` on a `typed_ptr` already compares addresses, so
// `===` does not apply to one."
func IsIdentityComparable(t Type) bool {
	s := AsStruct(t)
	return s != nil && s.Class()
}

// InteriorAssignable reports whether a value of type t grants assignability to
// what it holds, regardless of how the base binding was declared.
//
// §6.2 ⊢ an AssignTarget is assignable when it is "a field of an assignable
// value, or of any class or `shared`/`unique` handle" or "an element of an
// assignable `[N]T`, or of any `[]T` or `map[K]V`". Those "or of any" clauses
// are this predicate; Var.Assignable is the rest of the list. `weak` is absent:
// it is upgraded before it is read (§3.3), so there is nothing to assign
// through.
func InteriorAssignable(t Type) bool {
	switch u := Underlying(t).(type) {
	case *Struct:
		return u.Class()
	case *Ownership:
		return u.kind == Unique || u.kind == Shared
	case *Slice, *Map:
		return true
	}
	return false
}

// --------------------------------------------------------------- identity

// Identical reports type identity, which §4.1 makes the whole of
// assignability: ⊢ "a value is assignable to a destination when their types are
// identical. There is no subtyping, no coercion, and no promotion."
//
// Named types are compared by declaring object, not by shape (§3.1's
// nominality). A transparent alias never reaches here as a Named, because the
// resolver substituted its target.
func Identical(x, y Type) bool {
	if x == y {
		return true
	}
	if x == nil || y == nil {
		return false
	}

	switch a := x.(type) {
	case *Basic:
		b, ok := y.(*Basic)
		// byte and uint8 share one *Basic (see basic.go), so §2.3's alias rule
		// falls out of the pointer comparison above and needs no case here.
		return ok && a.kind == b.kind

	case *Named:
		b, ok := y.(*Named)
		if !ok || a.obj != b.obj {
			return false
		}
		// Two instantiations of one generic are identical only when their type
		// arguments are.
		if len(a.typeArgs) != len(b.typeArgs) {
			return false
		}
		for i := range a.typeArgs {
			if !Identical(a.typeArgs[i], b.typeArgs[i]) {
				return false
			}
		}
		return true

	case *Abstract:
		b, ok := y.(*Abstract)
		return ok && a.obj == b.obj

	case *TypeParam:
		b, ok := y.(*TypeParam)
		return ok && a == b

	case *Ownership:
		b, ok := y.(*Ownership)
		return ok && a.kind == b.kind && Identical(a.elem, b.elem)

	case *Array:
		b, ok := y.(*Array)
		return ok && a.len == b.len && Identical(a.elem, b.elem)

	case *Slice:
		b, ok := y.(*Slice)
		return ok && Identical(a.elem, b.elem)

	case *Map:
		b, ok := y.(*Map)
		return ok && Identical(a.key, b.key) && Identical(a.elem, b.elem)

	case *Chan:
		b, ok := y.(*Chan)
		return ok && Identical(a.elem, b.elem)

	case *Pointer:
		b, ok := y.(*Pointer)
		return ok && Identical(a.elem, b.elem)

	case *Vector:
		b, ok := y.(*Vector)
		return ok && a.lanes == b.lanes && Identical(a.elem, b.elem)

	case *Predicate:
		b, ok := y.(*Predicate)
		return ok && a.lanes == b.lanes

	case *Tuple:
		b, ok := y.(*Tuple)
		return ok && identicalTuple(a, b)

	case *Signature:
		b, ok := y.(*Signature)
		if !ok {
			return false
		}
		// §3.1 ⊢ "two `func` types are the same type when their parameter
		// types, marker, and result agree — parameter names are not part of
		// the type." Mode is part of it, because a `mut T` parameter is a
		// different convention from a bare one.
		//
		// Expected is not compared: it is not a Type, and a FunctionType can
		// never carry one.
		return a.variadic == b.variadic && a.marker == b.marker &&
			identicalTuple(a.params, b.params) &&
			identicalTuple(a.results, b.results)

	case *Tensor:
		b, ok := y.(*Tensor)
		if !ok || len(a.shape) != len(b.shape) || !Identical(a.elem, b.elem) {
			return false
		}
		for i := range a.shape {
			if a.shape[i] != b.shape[i] {
				return false
			}
		}
		return true

	case *Struct:
		b, ok := y.(*Struct)
		if !ok || a.class != b.class || len(a.fields) != len(b.fields) {
			return false
		}
		for i, f := range a.fields {
			g := b.fields[i]
			if f.Name != g.Name || !Identical(f.Type, g.Type) {
				return false
			}
		}
		return true
	}
	return false
}

func identicalTuple(a, b *Tuple) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		x, y := a.At(i), b.At(i)
		if x.mode != y.mode || !Identical(x.typ, y.typ) {
			return false
		}
	}
	return true
}

// ------------------------------------------------------------ comparable

// IsComparable reports whether `==` and `!=` are defined on t, and is both the
// predicate `comparable` admits and the one a map key must satisfy.
//
// §3.5's table:
//
//	comparable  numerics, bool, char, string, typed_ptr T, enums, and any
//	            struct, class, tuple, or [N]T whose every component is comparable
//	neither     []T, map[K]V, chan T, func types, vector, tensor, abstract,
//	            and the three heap handles
//
// Two readings are worth stating. First, the table qualifies struct, class,
// tuple, and [N]T componentwise but lists "enums" flat, so a payload enum is
// comparable here whatever its payloads hold; that is the table read literally
// rather than an oversight assumed. Second, a slice and a map are excluded even
// though their elements might be comparable, because comparing headers answers
// a question nobody asked and comparing contents is an O(n) walk the language
// never performs implicitly.
func IsComparable(t Type) bool {
	switch u := Underlying(t).(type) {
	case *Basic:
		return u.kind != Invalid && u.kind != UntypedNil

	case *Pointer:
		return true

	case *Enum:
		return true

	case *Array:
		return IsComparable(u.elem)

	case *Tuple:
		for i := 0; i < u.Len(); i++ {
			if !IsComparable(u.At(i).typ) {
				return false
			}
		}
		return true

	case *Struct:
		for _, f := range u.fields {
			if !IsComparable(f.Type) {
				return false
			}
		}
		return true

	case *TypeParam:
		return u.constraint == Comparable ||
			(u.constraint != nil && u.constraint.impliesComparable())
	}
	return false
}

func (c *Constraint) impliesComparable() bool {
	if c == Comparable {
		return true
	}
	for _, e := range c.embeds {
		if e.impliesComparable() {
			return true
		}
	}
	// A type set of comparable terms implies comparability, which is what lets
	// `[K: Ordered]` serve as a map key without also writing `comparable`.
	if len(c.terms) == 0 {
		return false
	}
	for _, term := range c.terms {
		if !IsComparable(term.Type) {
			return false
		}
	}
	return true
}

// ------------------------------------------------------------ assignment

// AssignableTo reports whether a value of type v may be assigned to type t.
//
// §4.1 ⊢ "a value is assignable to a destination when their types are
// identical. There is no subtyping, no coercion, and no promotion", relaxed by
// exactly two things: untyped literals, and §4.3's two implicit conversions —
// which are not general assignability and live in their own predicates below.
//
// Ownership is not an assignability question. §8.1 ⊢ whether a `var T`
// parameter receives the caller's original or a deep copy "is decided at the
// call site by the presence or absence of the `var` marker, never by the
// declaration", so a bare T argument is assignable to a `var T` parameter and
// the marker's absence means copy.
//
// The representability half of the untyped rule is not here: it needs the
// Value, not just the type, and it needs a target — see Sizes.Representable.
func AssignableTo(v, t Type) bool {
	if IsInvalid(v) || IsInvalid(t) {
		return true // already diagnosed; do not cascade
	}
	if Identical(v, t) {
		return true
	}
	if b, ok := v.(*Basic); ok && b.kind == UntypedNil {
		// §10 ⊢ "`nil` belongs to `typed_ptr T` and to nothing else."
		_, isPtr := Underlying(t).(*Pointer)
		return isPtr
	}
	if IsUntyped(v) {
		return untypedAssignable(v, t)
	}
	return false
}

func untypedAssignable(v, t Type) bool {
	src, _ := v.(*Basic)
	dst := AsBasic(t)
	if src == nil || dst == nil {
		return false
	}
	switch src.kind {
	case UntypedBool:
		return dst.is(InfoBoolean)
	case UntypedInt:
		// An untyped integer reaches a float destination; representability is
		// the caller's check, since it needs the Value and the target width.
		return dst.is(InfoNumeric)
	case UntypedFloat:
		return dst.is(InfoFloat)
	case UntypedChar:
		// §4.1 ⊢ 'A' and "A" "never interconvert implicitly", and char is not
		// an integer: an untyped char reaches a char destination and nothing
		// else.
		return dst.is(InfoChar)
	case UntypedString:
		return dst.is(InfoString)
	}
	return false
}

// ------------------------------------------------- the two implicit conversions

// LaunchConvertible reports §4.3's first implicit conversion: ⊢ at "a
// `gpu`/`npu` launch site", `[N]T ↔ tensor[T, N]`, "element type and shape
// matching exactly", bounded by "the launch expression itself".
//
// It is deliberately not part of AssignableTo. §4.3 ⊢ "neither reaches inside a
// body, and no third case is added anywhere", so the analyzer calls this at the
// launch site and nowhere else.
func LaunchConvertible(v, t Type) bool {
	if arrayTensorMatch(v, t) || arrayTensorMatch(t, v) {
		return true
	}
	return false
}

func arrayTensorMatch(a, b Type) bool {
	arr, ok := Underlying(a).(*Array)
	if !ok {
		return false
	}
	ten, ok := Underlying(b).(*Tensor)
	if !ok {
		return false
	}
	return Identical(arr.elem, ten.elem) && arr.len == ten.Elems()
}

// PointerCastElidable reports §4.3's second implicit conversion: ⊢ "a pointer
// cast with a written destination", `typed_ptr T` → `typed_ptr U` "with `as`
// elided", bounded by "both sides pointer types".
func PointerCastElidable(v, t Type) bool {
	if _, ok := Underlying(v).(*Pointer); !ok {
		return false
	}
	_, ok := Underlying(t).(*Pointer)
	return ok
}

// ------------------------------------------------------------ conversion

// ConvertibleTo reports whether `v as T` is legal.
//
// §4.2's table is the whole of it: integer→integer, float→integer,
// integer→float, enum→its discriminant type, `typed_ptr T`→`typed_ptr U`,
// `typed_ptr T`↔integer, and `abstract`→`typed_ptr T` for a memory-flat family.
// ⊢ "every width, signedness, or representation change between values is
// written", and there is no dynamic cast, because there is no runtime type
// information for one to consult.
//
// Three absences are deliberate, and each is a rule rather than an omission:
//
//   - enum→discriminant is one-way. ⊢ "there is no `n as Status`."
//   - The tensor element types take the constructor spelling instead. §4.2 ⊢
//     "`bf16(val)` is the form there, and `val as bf16` is not" — see
//     TensorElemConvertible.
//   - char↔integer is not in the table. §4.2 says "every... change between
//     values is written" and then enumerates; this reads the enumeration as
//     closed, so a code point is reached through a builtin rather than a cast.
func ConvertibleTo(v, t Type) bool {
	if IsInvalid(v) || IsInvalid(t) {
		return true
	}
	if AssignableTo(v, t) {
		return true
	}

	uv, ut := Underlying(v), Underlying(t)
	bv, bt := AsBasic(v), AsBasic(t)

	// abstract: one direction, one family, and nothing converts to one.
	if _, ok := ut.(*Abstract); ok {
		return false
	}
	if a, ok := uv.(*Abstract); ok {
		if _, isPtr := ut.(*Pointer); isPtr {
			return a.family == FamilyMemoryFlat
		}
		return false
	}

	// The tensor element types are outside `as` in either direction.
	if bv.is(InfoTensorElem) || bt.is(InfoTensorElem) {
		return false
	}

	// Numeric ↔ numeric. char is not numeric, so it is excluded here.
	if bv.is(InfoNumeric) && bt.is(InfoNumeric) {
		return true
	}

	// enum → its discriminant type, one-way.
	if e, ok := uv.(*Enum); ok {
		return bt != nil && Identical(e.discrim, bt)
	}

	// typed_ptr ↔ typed_ptr, and typed_ptr ↔ integer as an address value.
	if _, ok := uv.(*Pointer); ok {
		if _, ok := ut.(*Pointer); ok {
			return true
		}
		return bt.is(InfoInteger)
	}
	if _, ok := ut.(*Pointer); ok {
		return bv.is(InfoInteger)
	}

	return false
}

// TensorElemConvertible reports whether `T(val)` is the legal constructor form
// for a tensor element type.
//
// §4.2 ⊢ "the predeclared numeric types do not take the constructor spelling —
// write `i as float32`, not `float32(i)`. The tensor element types are the
// single exception: `bf16(val)` is the form there." So this is the only place a
// type name is callable, and it is separate from ConvertibleTo because the two
// spellings are not interchangeable in either direction.
//
// The reverse — reading a bf16 back out as a float32 — is not spelled by §4.2
// either way, so it is not licensed here.
func TensorElemConvertible(v, t Type) bool {
	if IsInvalid(v) || IsInvalid(t) {
		return true
	}
	return IsTensorElem(t) && IsNumeric(v) && !IsTensorElem(v)
}