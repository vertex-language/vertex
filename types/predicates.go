package types

// Underlying strips one layer of naming. A.7.3's `~T` is defined over this,
// and it is what every structural predicate below reaches for first.
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

func IsInvalid(t Type) bool {
	b := AsBasic(t)
	return b != nil && b.kind == Invalid
}

func IsUntyped(t Type) bool {
	b, _ := t.(*Basic)
	return b.is(IsUntyped)
}

func IsBool(t Type) bool    { return AsBasic(t).is(IsBoolean) }
func IsInteger(t Type) bool { return AsBasic(t).is(IsInteger) }
func IsFloat(t Type) bool   { return AsBasic(t).is(IsFloat) }
func IsNumeric(t Type) bool { return AsBasic(t).is(IsNumeric) }
func IsString(t Type) bool  { return AsBasic(t).is(IsString) }
func IsOrdered(t Type) bool { return AsBasic(t).is(IsOrdered) }

// --------------------------------------------------------------- identity

// Identical reports type identity.
//
// Named types are compared by object, not by shape: A.6.6 ⊢ "two abstract
// aliases never unify, however identical their provenance", and the same holds
// for two structs with the same fields. A transparent alias never reaches here
// as a Named, because the resolver substituted its target.
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
		// byte and uint8 share one *Basic (see basic.go), so A.1.4's "no
		// conversion is required or permitted between them" falls out of the
		// pointer comparison above and needs no case here.
		return ok && a.kind == b.kind

	case *Named:
		b, ok := y.(*Named)
		if !ok || a.obj != b.obj {
			return false
		}
		// A.7.5: two instantiations of one generic are identical only when
		// their type arguments are.
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

	case *Tuple:
		b, ok := y.(*Tuple)
		return ok && identicalTuple(a, b)

	case *Signature:
		b, ok := y.(*Signature)
		if !ok {
			return false
		}
		// A.3.4 ⊢ "a FunctionType names parameter types only; parameter names
		// belong to declarations, not to types", so names are not compared —
		// but Mode is, because a `mut T` parameter lowers to a pointer and a
		// bare one does not.
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

// IsComparable reports whether `==` and `!=` are defined on t (A.7.4).
//
// A.7.4 gives only "every type supporting == / !=", which is circular, so the
// set is pinned here:
//
//   - scalars and string compare by value;
//   - a fixed array is comparable when its element type is, since it is inline
//     storage with no header (A.3.1);
//   - a tuple likewise, elementwise;
//   - a struct or class is comparable when every field is — note this is `==`,
//     the "same bytes?" question, and is unrelated to A.4.5's `===` storage
//     identity, which is legal on classes regardless;
//   - a unit-only enum is its discriminant integer (A.6.5) and is comparable; a
//     payload enum is comparable when every payload type is;
//   - a slice and a map are NOT comparable — a slice is a {ptr,len,cap} header
//     over shared storage and comparing headers would answer a question nobody
//     asked, while comparing contents is an O(n) walk the language never
//     performs implicitly;
//   - typed_ptr compares as an address, which is why A.14 rejects *ordering*
//     one against nil while leaving equality alone;
//   - unique/shared/weak, chan, and abstract are not comparable: the first
//     three would have to choose between handle and pointee, and an abstract
//     handle "has no nil and never participates in a null comparison" (A.3.3).
//
// This is the predicate `map[K]V` requires of K (A.3.1) and the one `comparable`
// admits (A.7.4).
func IsComparable(t Type) bool {
	switch u := Underlying(t).(type) {
	case *Basic:
		return u.kind != Invalid && u.kind != UntypedNil

	case *Pointer:
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

	case *Enum:
		for _, v := range u.variants {
			for _, p := range v.Payload {
				if !IsComparable(p) {
					return false
				}
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
// Vertex has no implicit conversions between distinct types. A.1.5.2 is
// explicit that 'A' and "A" "never interconvert implicitly", and A.4.4 makes
// every widening an `as`. So the rule reduces to identity, plus one case: an
// untyped constant takes its destination's type when it is representable there
// (A.1.5.1).
//
// Ownership is not an assignability question. Whether a `var`-typed parameter
// receives the caller's original or a deep copy is decided by the presence of
// the marker at the call site (A.4.6), never by the type — so a bare T argument
// is assignable to a `var T` parameter and the marker's absence means copy.
func AssignableTo(v, t Type) bool {
	if IsInvalid(v) || IsInvalid(t) {
		return true // an error was already reported; do not cascade
	}
	if Identical(v, t) {
		return true
	}
	if IsUntyped(v) {
		return untypedAssignable(v, t)
	}
	// A.5.2 ⊢ nil is the map-erase operand and the typed_ptr absence value; it
	// has no type of its own.
	if b, ok := v.(*Basic); ok && b.kind == UntypedNil {
		_, isPtr := Underlying(t).(*Pointer)
		return isPtr
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
		return dst.is(IsBoolean)
	case UntypedInt:
		// An untyped integer reaches a float destination; the representability
		// check is the caller's, since it needs the Value and not just the type.
		return dst.is(IsNumeric)
	case UntypedFloat:
		return dst.is(IsFloat)
	case UntypedChar:
		return dst.is(IsChar)
	case UntypedString:
		return dst.is(IsString)
	}
	return false
}

// ------------------------------------------------------------ conversion

// ConvertibleTo reports whether `v as T` is legal (A.4.4).
//
// A.4.4 ⊢ `as` never touches memory: between pointer types it is a static
// reinterpretation and always legal; between numerics it is a width-selected
// instruction; on an enum it is a tag read. There is no dynamic cast, because
// there is no runtime type information for one to consult.
func ConvertibleTo(v, t Type) bool {
	if IsInvalid(v) || IsInvalid(t) {
		return true
	}
	if AssignableTo(v, t) {
		return true
	}

	uv, ut := Underlying(v), Underlying(t)

	// A.4.4: abstract → typed_ptr is legal only for a memory-flat family, and
	// there is no cast to abstract in any direction or family.
	if _, ok := ut.(*Abstract); ok {
		return false
	}
	if a, ok := uv.(*Abstract); ok {
		if _, isPtr := ut.(*Pointer); isPtr {
			return a.family == FamilyMemoryFlat
		}
		return false
	}

	// Numeric ↔ numeric, and char ↔ integer. A.1.5.2 keeps char and string
	// apart implicitly; an explicit `as` between them is still not a width
	// conversion and stays rejected.
	bv, bt := AsBasic(v), AsBasic(t)
	if bv != nil && bt != nil {
		switch {
		case bv.is(IsNumeric) && bt.is(IsNumeric):
			return true
		case bv.is(IsChar) && bt.is(IsInteger), bv.is(IsInteger) && bt.is(IsChar):
			return true
		}
	}

	// A.6.5 ⊢ a unit-only enum *is* its discriminant integer, so the cast is a
	// tag read. A payload enum has no such reading.
	if e, ok := uv.(*Enum); ok && e.UnitOnly() && bt.is(IsInteger) {
		return true
	}
	if e, ok := ut.(*Enum); ok && e.UnitOnly() && bv.is(IsInteger) {
		return true
	}

	// Pointer ↔ pointer is a static reinterpretation.
	if _, ok := uv.(*Pointer); ok {
		_, ok2 := ut.(*Pointer)
		return ok2
	}
	return false
}