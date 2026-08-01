package types

// intBits fixes the width of `int` and `uint`.
//
// A.1.4 lists them as distinct from int64/uint64 but never pins a width. They
// are fixed at 64 bits on every build tag, including js and wasm, rather than
// tracking the pointer width:
//
//   - A.15's invariant is that "every value has a statically known layout".
//     A width that varies by target makes `sizeof(int)` a per-target answer and
//     makes a struct's layout non-portable across a `build` split, which is
//     exactly the class of bug the build tag is supposed to contain.
//   - A.1.5.1 ⊢ "a literal that does not fit the destination type is a compile
//     error." A target-varying width would make that a target-varying error,
//     so a file could compile under one tag and fail under another for a reason
//     invisible in the source.
//
// wasm32 therefore pays a two-instruction cost on `int` arithmetic. That is the
// right side of the trade for a language whose whole premise is that cost is
// visible and layout is fixed.
const intBits = 64

// Sizes answers layout questions for a target. Only the pointer width varies;
// see intBits.
type Sizes struct {
	WordSize int64 // pointer and word width in bytes
	MaxAlign int64 // maximum alignment of any type
}

var (
	Sizes64 = &Sizes{WordSize: 8, MaxAlign: 8}
	Sizes32 = &Sizes{WordSize: 4, MaxAlign: 8}
)

// SizesFor returns the layout rules for a build tag (A.2.2).
func SizesFor(tag string) *Sizes {
	switch tag {
	case "js", "wasm":
		return Sizes32
	}
	return Sizes64
}

// Alignof returns t's alignment in bytes.
func (s *Sizes) Alignof(t Type) int64 {
	switch u := Underlying(t).(type) {
	case *Basic:
		if a := basicSize(u.kind, s.WordSize); a > 0 {
			if a > s.MaxAlign {
				return s.MaxAlign
			}
			return a
		}
		return 1

	case *Array:
		// A.3.1 ⊢ [N]T is inline storage with no header, so it aligns as T.
		return s.Alignof(u.elem)

	case *Struct:
		var max int64 = 1
		for _, f := range u.fields {
			if a := s.Alignof(f.Type); a > max {
				max = a
			}
		}
		return max

	case *Tuple:
		var max int64 = 1
		for i := 0; i < u.Len(); i++ {
			if a := s.Alignof(u.At(i).typ); a > max {
				max = a
			}
		}
		return max

	case *Enum:
		max := s.Alignof(u.discrim)
		for _, v := range u.variants {
			for _, p := range v.Payload {
				if a := s.Alignof(p); a > max {
					max = a
				}
			}
		}
		return max
	}
	return s.WordSize
}

// Sizeof returns t's size in bytes.
func (s *Sizes) Sizeof(t Type) int64 {
	switch u := Underlying(t).(type) {
	case *Basic:
		return basicSize(u.kind, s.WordSize)

	case *Array:
		if u.len == 0 {
			return 0
		}
		return u.len * s.arrayStride(u.elem)

	case *Slice:
		// A.3.1 ⊢ []T is a three-word {ptr, len, cap} header.
		return 3 * s.WordSize

	case *Map:
		return s.WordSize

	case *Chan:
		// A.3.5 ⊢ an implicitly heap-resident refcounted handle.
		return s.WordSize

	case *Pointer:
		return s.WordSize

	case *Ownership:
		// A.3.2 ⊢ unique, shared, and weak "are ordinary one-word value types".
		return s.WordSize

	case *Abstract:
		return s.WordSize

	case *Signature:
		// A.3.4 ⊢ "a non-capturing function value is one word — a bare code
		// pointer. A capturing closure is two words, {code, env}." The type
		// alone does not say which, so this is the conservative answer; lower
		// narrows it from the expression that produced the value.
		return 2 * s.WordSize

	case *Struct:
		return s.structSize(u.fields)

	case *Tuple:
		fields := make([]*Field, u.Len())
		for i := range fields {
			v := u.At(i)
			fields[i] = &Field{Name: v.name, Type: v.typ}
		}
		return s.structSize(fields)

	case *Enum:
		return s.enumSize(u)

	case *Tensor:
		n := int64(1)
		for _, d := range u.shape {
			n *= d
		}
		return n * s.Sizeof(u.elem)
	}
	return 0
}

// Offsetsof returns each field's byte offset. A.6.2 ⊢ "fields laid out in
// declaration order with ABI padding" — there is no reordering, so a reader
// can predict the layout from the source.
func (s *Sizes) Offsetsof(fields []*Field) []int64 {
	offs := make([]int64, len(fields))
	var cur int64
	for i, f := range fields {
		a := s.Alignof(f.Type)
		cur = align(cur, a)
		offs[i] = cur
		cur += s.Sizeof(f.Type)
	}
	return offs
}

func (s *Sizes) structSize(fields []*Field) int64 {
	if len(fields) == 0 {
		return 0
	}
	offs := s.Offsetsof(fields)
	last := fields[len(fields)-1]
	size := offs[len(offs)-1] + s.Sizeof(last.Type)

	var maxAlign int64 = 1
	for _, f := range fields {
		if a := s.Alignof(f.Type); a > maxAlign {
			maxAlign = a
		}
	}
	return align(size, maxAlign)
}

// enumSize implements A.6.5: a unit-only enum *is* its discriminant integer; a
// payload enum is "a tagged union sized to the largest variant plus the tag".
func (s *Sizes) enumSize(e *Enum) int64 {
	if e.UnitOnly() {
		return s.Sizeof(e.discrim)
	}

	tagSize := s.Sizeof(e.discrim)
	var payloadAlign int64 = 1
	var payloadSize int64

	for _, v := range e.variants {
		if len(v.Payload) == 0 {
			continue
		}
		fields := make([]*Field, len(v.Payload))
		for i, p := range v.Payload {
			fields[i] = &Field{Type: p}
		}
		if sz := s.structSize(fields); sz > payloadSize {
			payloadSize = sz
		}
		for _, p := range v.Payload {
			if a := s.Alignof(p); a > payloadAlign {
				payloadAlign = a
			}
		}
	}

	total := align(tagSize, payloadAlign) + payloadSize
	maxAlign := payloadAlign
	if a := s.Alignof(e.discrim); a > maxAlign {
		maxAlign = a
	}
	return align(total, maxAlign)
}

// arrayStride is the per-element extent inside a fixed array, which is the
// element size rounded up to its own alignment.
func (s *Sizes) arrayStride(elem Type) int64 {
	return align(s.Sizeof(elem), s.Alignof(elem))
}

func align(x, a int64) int64 {
	if a <= 1 {
		return x
	}
	return (x + a - 1) &^ (a - 1)
}

func basicSize(k BasicKind, word int64) int64 {
	switch k {
	case Bool, Int8, Uint8:
		return 1
	case Int16, Uint16:
		return 2
	case Int32, Uint32, Float32:
		return 4
	case Int64, Uint64, Float64:
		return 8
	case Int, Uint:
		return intBits / 8
	case Char:
		// A.1.5.2 ⊢ "exactly one Unicode scalar value, held in 4 bytes."
		return 4
	case String:
		// A.1.5.2 ⊢ UTF-8 bytes with a length and no NUL terminator: {ptr, len}.
		return 2 * word
	}
	return 0
}