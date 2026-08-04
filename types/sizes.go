package types

import (
	"math/big"

	"github.com/vertex-language/vertex/token"
)

// Sizes answers layout questions for one target.
//
// §2.3 ⊢ "`int` and `uint` are the target's pointer width and are distinct
// types from `int64`/`uint64` even where the widths agree", so WordSize drives
// both the pointer width and the width of those two types. Everything that
// depends on it — Sizeof, Alignof, and Representable — is therefore a method
// here rather than a free function: a caller cannot ask a layout question
// without saying which target it is asking about.
type Sizes struct {
	WordSize int64 // pointer, and `int`/`uint`, width in bytes
	MaxAlign int64 // maximum alignment of a scalar or aggregate
}

var (
	Sizes64 = &Sizes{WordSize: 8, MaxAlign: 8}
	Sizes32 = &Sizes{WordSize: 4, MaxAlign: 8}
)

// SizesFor returns the layout rules for a build tag. §1.2 makes the tag a
// whole-file selection rule, so one compilation has one answer.
func SizesFor(tag token.BuildTag) *Sizes {
	switch tag {
	case token.TagJS, token.TagWasm:
		return Sizes32
	}
	return Sizes64
}

// intWidth reports a sized integer kind's width in bits and its signedness.
func (s *Sizes) intWidth(k BasicKind) (bits uint, signed, ok bool) {
	switch k {
	case Int4:
		return 4, true, true
	case Int8:
		return 8, true, true
	case Int16:
		return 16, true, true
	case Int32:
		return 32, true, true
	case Int64:
		return 64, true, true
	case Uint8:
		return 8, false, true
	case Uint16:
		return 16, false, true
	case Uint32:
		return 32, false, true
	case Uint64:
		return 64, false, true
	case Int:
		return uint(s.WordSize) * 8, true, true
	case Uint:
		return uint(s.WordSize) * 8, false, true
	}
	return 0, false, false
}

// intRange returns the inclusive bounds of a sized integer kind. It is computed
// rather than tabled because the bounds of `int` and `uint` are the target's.
func (s *Sizes) intRange(k BasicKind) (lo, hi *big.Int, ok bool) {
	bits, signed, ok := s.intWidth(k)
	if !ok {
		return nil, nil, false
	}
	if signed {
		h := new(big.Int).Lsh(big.NewInt(1), bits-1)
		l := new(big.Int).Neg(h)
		return l, h.Sub(h, big.NewInt(1)), true
	}
	h := new(big.Int).Lsh(big.NewInt(1), bits)
	return big.NewInt(0), h.Sub(h, big.NewInt(1)), true
}

// Alignof returns t's alignment in bytes.
func (s *Sizes) Alignof(t Type) int64 {
	switch u := Underlying(t).(type) {
	case *Basic:
		if a := s.basicSize(u.kind); a > 0 {
			if a > s.MaxAlign {
				return s.MaxAlign
			}
			return a
		}
		return 1

	case *Array:
		// Inline storage with no header, so it aligns as its element does.
		return s.Alignof(u.elem)

	case *Vector:
		// The sources fix no vector ABI. A lane-wise load wants the whole
		// vector naturally aligned, so this aligns to its own size rather than
		// clamping to MaxAlign — an implementation choice, not a stated rule.
		return roundPow2(s.Sizeof(u))

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

	case *Predicate:
		return 1
	}
	return s.WordSize
}

// Sizeof returns t's size in bytes.
//
// §3.4 ⊢ "every type has a size known at compile time", and it fixes one layout
// fact directly: the seven indirections are "all one word". Everything else
// below — the string header, the closure representation, the enum layout — is
// this implementation's choice, because neither source document fixes it.
func (s *Sizes) Sizeof(t Type) int64 {
	switch u := Underlying(t).(type) {
	case *Basic:
		return s.basicSize(u.kind)

	case *Array:
		if u.len <= 0 {
			return 0
		}
		return u.len * s.arrayStride(u.elem)

	// §3.4's one-word indirections.
	case *Slice, *Map, *Chan, *Pointer, *Ownership:
		return s.WordSize

	case *Abstract:
		// A handle, one word, like the indirections it sits beside.
		return s.WordSize

	case *Signature:
		// A function value is at most {code, env}. The sources fix no
		// representation and §7.3's capture-by-value is a property of the
		// expression rather than of the type, so this is the conservative
		// answer and lower narrows it from what produced the value.
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

	case *Vector:
		return u.lanes * s.Sizeof(u.elem)

	case *Tensor:
		// Element widths are taken in bits, because int4 is sub-byte and only
		// ever exists packed inside one of these.
		bits := u.Elems() * s.elemBits(u.elem)
		return align(bits, 8) / 8

	case *Predicate:
		// §5.1 keeps it out of every storage position, so it never has a
		// layout to report. A backend may give the lane mask whatever
		// representation it likes.
		return 0
	}
	return 0
}

// Offsetsof returns each field's byte offset. Fields are laid out in
// declaration order with padding and no reordering, so a reader can predict the
// layout from the source; §7.2 makes declaration order observable anyway,
// through destruction order.
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

// enumSize lays a unit-only enum out as its discriminant and a payload enum as
// a tag plus the largest variant. The sources fix no enum layout; this is the
// implementation's choice and is not normative. What is normative is only that
// the size exists at compile time (§3.4) and that the zero value is the first
// variant with any payload zeroed (§3.3).
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

// arrayStride is the per-element extent inside a fixed array: the element size
// rounded up to its own alignment.
func (s *Sizes) arrayStride(elem Type) int64 {
	return align(s.Sizeof(elem), s.Alignof(elem))
}

// elemBits is a tensor element's width in bits. int4 has no byte size of its
// own, which is why tensor sizing goes through here and Sizeof(Typ[Int4])
// answers zero.
func (s *Sizes) elemBits(elem Type) int64 {
	if b := AsBasic(elem); b != nil && b.kind == Int4 {
		return 4
	}
	return s.Sizeof(elem) * 8
}

func (s *Sizes) basicSize(k BasicKind) int64 {
	switch k {
	case Bool, Int8, Uint8, FP8E4M3, FP8E5M2:
		return 1
	case Int16, Uint16, BF16:
		return 2
	case Int32, Uint32, Float32:
		return 4
	case Int64, Uint64, Float64:
		return 8
	case Int, Uint:
		return s.WordSize
	case Char:
		// One Unicode scalar value needs 21 bits; four bytes is the natural
		// carrier, and the sources fix no width.
		return 4
	case String:
		// grammar.md ⊢ "a string is UTF-8 bytes with a length and no NUL
		// terminator", so a length is stored: {ptr, len}.
		return 2 * s.WordSize
	case Int4:
		// Sub-byte, and only meaningful packed inside a tensor. See elemBits.
		return 0
	}
	return 0
}

func align(x, a int64) int64 {
	if a <= 1 {
		return x
	}
	return (x + a - 1) &^ (a - 1)
}

func roundPow2(x int64) int64 {
	p := int64(1)
	for p < x {
		p <<= 1
	}
	return p
}

// ------------------------------------------------------------------- cost

// CopyKind is what a bare copy of a value costs.
//
// §8.5 prices every copy and every transfer. The transfer column is uniformly
// O(1), so only the bare copy needs a classification:
//
//	scalars, typed_ptr        register move
//	struct, class, [N]T       fieldwise copy
//	string, []T, map[K]V      deep-copies the payload
//	unique T                  allocates and deep-copies the pointee
//	shared T, chan T          refcount increment
//	weak T                    weak-count increment
//
// §8.5 ⊢ "under generics the cost is fixed by the concrete type at
// instantiation, so a lint on large owned types fires per instantiation, not
// per declaration" — which is why this takes a Type and not a declaration.
type CopyKind uint8

const (
	CopyRegister CopyKind = iota
	CopyFieldwise
	CopyDeep
	CopyAlloc
	CopyRefcount
	CopyWeakcount
)

// CopyCost classifies a bare copy of t.
//
// CopyFieldwise is recursive: a struct holding a `[]T` still deep-copies that
// field, because the fieldwise copy copies each field at that field's own cost.
func CopyCost(t Type) CopyKind {
	switch u := Underlying(t).(type) {
	case *Basic:
		if u.is(InfoString) {
			return CopyDeep
		}
		return CopyRegister

	case *Pointer, *Signature, *Abstract, *Predicate:
		return CopyRegister

	case *Slice, *Map:
		return CopyDeep

	case *Chan:
		return CopyRefcount

	case *Ownership:
		switch u.kind {
		case Unique:
			return CopyAlloc
		case Weak:
			return CopyWeakcount
		}
		return CopyRefcount

	case *Array, *Struct, *Tuple, *Enum, *Vector, *Tensor:
		return CopyFieldwise
	}
	return CopyFieldwise
}