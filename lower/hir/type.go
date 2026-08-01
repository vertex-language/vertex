package hir

import (
	"fmt"
	"strings"

	"github.com/vertex-language/vertex/types"
)

// Type is hir's type representation: vir's type vocabulary, nothing more.
// Every Vertex type has already been decided into one of these by the time
// a Func body exists — there is no ownership qualifier, no generic
// parameter, and no constraint left to interpret.
type Type interface {
	String() string
	typeNode()
}

type IntType struct{ Bits int }   // i1, i8, i16, i32, i64, i128
type FloatType struct{ Bits int } // f16, f32, f64
type PtrType struct{}             // untyped; width is the target's usize
type VoidType struct{}
type StructType struct{ Def *Struct }
type ArrayType struct {
	Elem Type
	Len  int64
}

func (IntType) typeNode()    {}
func (FloatType) typeNode()  {}
func (PtrType) typeNode()    {}
func (VoidType) typeNode()   {}
func (StructType) typeNode() {}
func (ArrayType) typeNode()  {}

func (t IntType) String() string    { return fmt.Sprintf("i%d", t.Bits) }
func (t FloatType) String() string  { return fmt.Sprintf("f%d", t.Bits) }
func (PtrType) String() string      { return "ptr" }
func (VoidType) String() string     { return "void" }
func (t StructType) String() string { return "struct " + t.Def.Name }
func (t ArrayType) String() string  { return fmt.Sprintf("array[%s, %d]", t.Elem, t.Len) }

var (
	I1   = IntType{1}
	I8   = IntType{8}
	I16  = IntType{16}
	I32  = IntType{32}
	I64  = IntType{64}
	I128 = IntType{128}
	F32  = FloatType{32}
	F64  = FloatType{64}
	Ptr  = PtrType{}
	Void = VoidType{}
)

// IsAggregate reports whether t is memory-only. A Value of aggregate type
// is a pointer at the vir level — see the package doc comment.
func IsAggregate(t Type) bool {
	switch t.(type) {
	case StructType, ArrayType:
		return true
	}
	return false
}

func IsInt(t Type) bool   { _, ok := t.(IntType); return ok }
func IsFloat(t Type) bool { _, ok := t.(FloatType); return ok }
func IsPtr(t Type) bool   { _, ok := t.(PtrType); return ok }
func IsVoid(t Type) bool  { _, ok := t.(VoidType); return ok }

func TypeEqual(a, b Type) bool {
	switch x := a.(type) {
	case IntType:
		y, ok := b.(IntType)
		return ok && x.Bits == y.Bits
	case FloatType:
		y, ok := b.(FloatType)
		return ok && x.Bits == y.Bits
	case PtrType:
		_, ok := b.(PtrType)
		return ok
	case VoidType:
		_, ok := b.(VoidType)
		return ok
	case StructType:
		y, ok := b.(StructType)
		return ok && x.Def == y.Def
	case ArrayType:
		y, ok := b.(ArrayType)
		return ok && x.Len == y.Len && TypeEqual(x.Elem, y.Elem)
	}
	return false
}

// ---------------------------------------------------------------------------
// types.Type -> hir.Type
// ---------------------------------------------------------------------------

// typeLowerer maps checked Vertex types onto hir types, memoized by the
// rendered spelling of the source type so a diamond of references to the
// same struct lands on one *Struct.
//
// The fat-type shapes below are the language's own layout commitments,
// written once here:
//
//	string    {ptr, len}                 immutable, two words (A.9.4)
//	[]T       {ptr, len, cap}            three words (A.3.1)
//	map[K]V   ptr                        opaque handle into builtins/map
//	chan T    ptr                        refcounted handle (A.10.1)
//	unique T  ptr                        one word, deep copy (A.9.4)
//	shared T  ptr                        one word, atomic increment
//	weak T    ptr                        observes a shared allocation
//	func(...) {code, env}                two words; the one-word non-capturing
//	                                     form is what crosses a boundary (A.3.4)
//	abstract  ptr                        foreign handle, interior invisible
//	(A, B)    struct                     parentheses are the type's shape
//	enum      unit-only: its discriminant integer; payload: {tag, bytes}
type typeLowerer struct {
	l     *lowerer
	cache map[string]Type
	// synthesized headers are declared per-module, since a struct produces
	// no symbol and so cannot collide across modules.
	headers map[*Module]map[string]*Struct
	word    int64
}

func newTypeLowerer(l *lowerer) *typeLowerer {
	w := l.conf.PointerSize
	if w == 0 {
		w = 8
	}
	return &typeLowerer{l: l, cache: map[string]Type{}, headers: map[*Module]map[string]*Struct{}, word: w}
}

// lower converts a checked type to its hir shape, in the module currently
// being emitted into (synthesized headers are declared there).
func (tl *typeLowerer) lower(m *Module, t types.Type) Type {
	if t == nil {
		return Void
	}
	t = tl.l.subst(t)
	key := m.Path + "\x00" + types.TypeString(t)
	if got, ok := tl.cache[key]; ok {
		return got
	}
	out := tl.build(m, t)
	tl.cache[key] = out
	return out
}

func (tl *typeLowerer) build(m *Module, t types.Type) Type {
	switch tl.l.classify(t) {
	case kInvalid:
		return Void
	case kUnit:
		return Void
	case kBool:
		// vir's i1 is the boolean type; comparisons already yield it, and
		// a bool field is one byte of storage per types.Sizes.
		return I1
	case kChar:
		return I32 // A.1.5.2: one Unicode scalar value, held in 4 bytes
	case kInt:
		return IntType{Bits: tl.l.intBits(t)}
	case kFloat:
		return FloatType{Bits: tl.l.floatBits(t)}
	case kString:
		return StructType{tl.header(m, "vx_string", []Field{
			{Name: "ptr", Type: Ptr}, {Name: "len", Type: I64},
		})}
	case kSlice:
		return StructType{tl.header(m, "vx_slice", []Field{
			{Name: "ptr", Type: Ptr}, {Name: "len", Type: I64}, {Name: "cap", Type: I64},
		})}
	case kArray:
		elem, n := tl.l.arrayParts(t)
		return ArrayType{Elem: tl.lower(m, elem), Len: n}
	case kMap, kChan, kPointer, kAbstract, kUnique, kShared, kWeak:
		return Ptr
	case kFunc:
		// A.3.4: two words, {code, env}. The non-capturing one-word form is
		// narrowed per-expression at a boundary, not by the type.
		return StructType{tl.header(m, "vx_closure", []Field{
			{Name: "code", Type: Ptr}, {Name: "env", Type: Ptr},
		})}
	case kTuple:
		return StructType{tl.tuple(m, t)}
	case kStruct:
		return StructType{tl.named(m, t)}
	case kEnum:
		return tl.enum(m, t)
	case kTensor:
		tl.l.errorf(0, "todo: tensor types have no host lowering — a device-marked body is out of scope for VIR emission")
		return Void
	}
	tl.l.errorf(0, "internal: no hir lowering for type %s", types.TypeString(t))
	return Void
}

// header declares (or reuses) a synthesized layout struct in m. Fields are
// laid out with natural alignment, matching A.6.2's declaration-order rule.
func (tl *typeLowerer) header(m *Module, name string, fields []Field) *Struct {
	byMod := tl.headers[m]
	if byMod == nil {
		byMod = map[string]*Struct{}
		tl.headers[m] = byMod
	}
	if s, ok := byMod[name]; ok {
		return s
	}
	s := &Struct{Name: m.uniqueName(name), Module: m}
	tl.layout(s, fields)
	byMod[name] = s
	m.Structs = append(m.Structs, s)
	return s
}

// layout assigns offsets in declaration order with ABI padding and no
// reordering (A.6.2) — a reader predicts layout straight from source order.
func (tl *typeLowerer) layout(s *Struct, fields []Field) {
	var off, max int64 = 0, 1
	out := make([]Field, len(fields))
	for i, f := range fields {
		a := tl.align(f.Type)
		if a > max {
			max = a
		}
		if r := off % a; r != 0 {
			off += a - r
		}
		f.Offset = off
		out[i] = f
		off += tl.size(f.Type)
	}
	if r := off % max; r != 0 {
		off += max - r
	}
	s.Fields, s.Size, s.Align = out, off, max
}

func (tl *typeLowerer) size(t Type) int64 {
	switch x := t.(type) {
	case IntType:
		if x.Bits <= 8 {
			return 1
		}
		return int64((x.Bits + 7) / 8)
	case FloatType:
		return int64(x.Bits / 8)
	case PtrType:
		return tl.word
	case VoidType:
		return 0
	case StructType:
		return x.Def.Size
	case ArrayType:
		return x.Len * tl.size(x.Elem) // arrays have no inter-element padding
	}
	return 0
}

func (tl *typeLowerer) align(t Type) int64 {
	switch x := t.(type) {
	case StructType:
		return x.Def.Align
	case ArrayType:
		return tl.align(x.Elem)
	}
	if n := tl.size(t); n > 0 {
		return n
	}
	return 1
}

// Sizeof and Alignof are the queries lower/vir, own.go, and the builtin
// allocation sites all share; nothing recomputes layout independently.
func (tl *typeLowerer) Sizeof(t Type) int64  { return tl.size(t) }
func (tl *typeLowerer) Alignof(t Type) int64 { return tl.align(t) }

func (tl *typeLowerer) named(m *Module, t types.Type) *Struct {
	name, fields := tl.l.structParts(t)
	byMod := tl.headers[m]
	if byMod == nil {
		byMod = map[string]*Struct{}
		tl.headers[m] = byMod
	}
	key := "n:" + types.TypeString(t)
	if s, ok := byMod[key]; ok {
		return s
	}
	s := &Struct{Name: m.uniqueName(name), Module: m, Origin: t}
	byMod[key] = s
	m.Structs = append(m.Structs, s)

	// Bind the *Struct before lowering fields, so a self-referential field
	// (next: typed_ptr Node) reaches its own enclosing type without looping
	// — the same trick analyzer's recordDecl uses.
	hf := make([]Field, 0, len(fields))
	for _, f := range fields {
		hf = append(hf, Field{Name: sanitize(f.Name), Type: tl.lower(m, f.Type)})
	}
	tl.layout(s, hf)
	return s
}

func (tl *typeLowerer) tuple(m *Module, t types.Type) *Struct {
	elems := tl.l.tupleElems(t)
	var sb strings.Builder
	sb.WriteString("vx_tuple")
	hf := make([]Field, 0, len(elems))
	for i, e := range elems {
		sb.WriteString("_")
		sb.WriteString(sanitize(types.TypeString(e.Type)))
		name := e.Name
		if name == "" {
			name = "_" + itoa(i)
		}
		hf = append(hf, Field{Name: sanitize(name), Type: tl.lower(m, e.Type)})
	}
	return tl.header(m, sb.String(), hf)
}

// enum implements A.6.5 directly: a unit-only enum *is* its discriminant
// integer, so it lowers to that integer and nothing else. A payload enum is
// a tag plus opaque storage sized to the largest variant; the tag tells a
// copy or teardown routine which interpretation to walk.
func (tl *typeLowerer) enum(m *Module, t types.Type) Type {
	disc, unitOnly, payload := tl.l.enumParts(t)
	tag := IntType{Bits: tl.l.intBits(disc)}
	if unitOnly {
		return tag
	}
	var maxSize, maxAlign int64 = 0, 1
	for _, v := range payload {
		for _, ft := range v.Types {
			h := tl.lower(m, ft)
			if s := tl.size(h); s > maxSize {
				maxSize = s
			}
			if a := tl.align(h); a > maxAlign {
				maxAlign = a
			}
		}
	}
	name := "vx_enum_" + sanitize(types.TypeString(t))
	return StructType{tl.header(m, name, []Field{
		{Name: "tag", Type: tag},
		{Name: "payload", Type: ArrayType{Elem: I8, Len: maxSize}},
	})}
}