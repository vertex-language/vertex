package hir

import (
	"strings"

	"github.com/vertex-language/vertex/types"
)

// Type is hir's own type representation. It is not types.Type: by the time a
// body is lowered every type is concrete (monomorphization ran first) and
// every aggregate has a name and a layout, neither of which types.Type
// carries. It is also not vir.Type, because hir must not import vvm.
type Type struct {
	Kind   Kind
	Bits   int     // Int, Float
	Signed bool    // Int
	Elem   *Type   // Array, Vector, Predicate
	Len    int64   // Array, Vector, Predicate
	Struct *Struct // Struct
}

type Kind uint8

const (
	KVoid Kind = iota
	KInt
	KFloat
	KPtr
	KStruct
	KArray
	KVector
	KPredicate // vec[i1, N]; a value-only type, never in memory
)

// The scalars. bool is i8 because i1 has no ABI-agreed memory
// representation; comparisons yield i1 and widen at the store. char is i32.
var (
	Void = &Type{Kind: KVoid}
	I1   = &Type{Kind: KInt, Bits: 1, Signed: false}
	I8   = &Type{Kind: KInt, Bits: 8, Signed: true}
	I16  = &Type{Kind: KInt, Bits: 16, Signed: true}
	I32  = &Type{Kind: KInt, Bits: 32, Signed: true}
	I64  = &Type{Kind: KInt, Bits: 64, Signed: true}
	U8   = &Type{Kind: KInt, Bits: 8}
	U16  = &Type{Kind: KInt, Bits: 16}
	U32  = &Type{Kind: KInt, Bits: 32}
	U64  = &Type{Kind: KInt, Bits: 64}
	F32  = &Type{Kind: KFloat, Bits: 32}
	F64  = &Type{Kind: KFloat, Bits: 64}
	Ptr  = &Type{Kind: KPtr}

	Bool = I8
	Char = I32
)

func ArrayOf(elem *Type, n int64) *Type  { return &Type{Kind: KArray, Elem: elem, Len: n} }
func VectorOf(elem *Type, n int64) *Type { return &Type{Kind: KVector, Elem: elem, Len: n} }
func StructOf(s *Struct) *Type           { return &Type{Kind: KStruct, Struct: s} }

// IsAggregate reports whether t is memory-only, and therefore *is* a ptr at
// the vir level. Vector and the lane predicate are deliberately absent —
// they are register-class types (lowering.md §4.3, §15).
func IsAggregate(t *Type) bool {
	return t != nil && (t.Kind == KStruct || t.Kind == KArray)
}

func IsInt(t *Type) bool   { return t != nil && t.Kind == KInt }
func IsFloat(t *Type) bool { return t != nil && t.Kind == KFloat }
func IsPtr(t *Type) bool   { return t != nil && t.Kind == KPtr }
func IsVoid(t *Type) bool  { return t == nil || t.Kind == KVoid }

func (t *Type) String() string {
	switch t.Kind {
	case KVoid:
		return "void"
	case KInt:
		return "i" + itoa(t.Bits)
	case KFloat:
		return "f" + itoa(t.Bits)
	case KPtr:
		return "ptr"
	case KStruct:
		return "struct " + t.Struct.Name
	case KArray:
		return "array[" + t.Elem.String() + ", " + itoa(int(t.Len)) + "]"
	case KVector:
		return "vec[" + t.Elem.String() + ", " + itoa(int(t.Len)) + "]"
	case KPredicate:
		return "vec[i1, " + itoa(int(t.Len)) + "]"
	}
	return "?"
}

// ------------------------------------------------------- from types.Type

// typ lowers a checked type. Every call goes through here, so every struct
// the program needs is declared exactly once per owning module and every
// layout question is answered by conf.Sizes and nothing else.
//
// The `in` module is where a synthesized aggregate lands. A named type goes
// to its own declaring module; a header or tuple goes to the module that
// first needed it, which is why lower/vir declares them per-module and why
// the cross-package byval/sret question in its README is still open.
func (l *lowerer) typ(t types.Type, in *Module) *Type {
	switch u := types.Underlying(t).(type) {
	case nil:
		return l.bug("nil type reached lowering")

	case *types.Basic:
		return l.basic(u)

	case *types.Named:
		// Unreachable: Underlying already stripped one layer of naming.
		return l.bug("named type survived Underlying")

	case *types.Struct:
		return StructOf(l.structFor(t, in))

	case *types.Enum:
		if u.UnitOnly() {
			// A unit enum is its discriminant integer, full stop.
			return l.basic(u.Discriminant())
		}
		return StructOf(l.enumStruct(t, u, in))

	case *types.Array:
		return ArrayOf(l.typ(u.Elem(), in), u.Len())

	case *types.Tuple:
		return StructOf(l.tupleStruct(u, in))

	case *types.Vector:
		return VectorOf(l.typ(u.Elem(), in), u.Lanes())

	case *types.Predicate:
		return &Type{Kind: KPredicate, Elem: I1, Len: u.Lanes()}

	case *types.Slice:
		// {ptr, len, cap}. A view and a dynamic slice share one shape here;
		// the view simply never grows. lowering.md §4.2 splits them, and
		// splitting is a size optimization we have not taken.
		return StructOf(l.header(in, "_Vvec", field("p", Ptr), field("len", I64), field("cap", I64)))

	case *types.Map, *types.Chan:
		// One ptr to a builtins-side table or channel core. The element type
		// reaches the runtime as sizes and function pointers, not as a shape.
		return Ptr

	case *types.Pointer, *types.Abstract:
		return Ptr

	case *types.Ownership:
		// unique/shared/weak are all one word. Which one decides the copy and
		// drop cost, and that is own.go's question, read off the types.Type
		// rather than off this.
		return Ptr

	case *types.Signature:
		// {code, env}. A non-capturing function still gets the pair; §12.2's
		// thunk is what fills the code word.
		return StructOf(l.header(in, "_Vfn", field("code", Ptr), field("env", Ptr)))

	case *types.TypeParam:
		return l.bug("type parameter survived monomorphization: " + types.TypeString(t))
	}
	return l.bug("unlowerable type " + types.TypeString(t))
}

func (l *lowerer) basic(b *types.Basic) *Type {
	switch b.Kind() {
	case types.Bool:
		return Bool
	case types.Int8:
		return I8
	case types.Int16:
		return I16
	case types.Int32:
		return I32
	case types.Int64:
		return I64
	case types.Uint8:
		return U8
	case types.Uint16:
		return U16
	case types.Uint32:
		return U32
	case types.Uint64:
		return U64
	case types.Int:
		// semantics.md §2.3: the target's pointer width, and a distinct type
		// from int64 even where the widths agree. The distinctness mattered
		// in the checker; here only the width survives.
		return intOfBits(int(l.conf.Sizes.WordSize)*8, true)
	case types.Uint:
		return intOfBits(int(l.conf.Sizes.WordSize)*8, false)
	case types.Float32:
		return F32
	case types.Float64:
		return F64
	case types.Char:
		return Char
	case types.String:
		return StructOf(l.header(l.mod, "_Vstr", field("p", Ptr), field("len", I64)))
	case types.BF16, types.FP8E4M3, types.FP8E5M2, types.Int4:
		// todo: tensor element types only exist packed inside a tensor, and
		// tensors only exist in an npu body, which has no CPU lowering.
		return l.todo("tensor element type " + b.Name())
	}
	return l.bug("unlowerable basic " + b.Name())
}

func intOfBits(bits int, signed bool) *Type {
	switch {
	case bits == 32 && signed:
		return I32
	case bits == 32:
		return U32
	case bits == 64 && signed:
		return I64
	}
	return U64
}

func field(name string, t *Type) StructField { return StructField{Name: name, Type: t} }

// header declares a synthesized aggregate once per module. The name is the
// key: _Vstr means the same shape everywhere, so re-deriving it is free and
// interning it keeps the emitted text stable.
func (l *lowerer) header(in *Module, name string, fields ...StructField) *Struct {
	if s, ok := in.structs[name]; ok {
		return s
	}
	s := &Struct{Name: name, Module: in, Fields: fields}
	l.layout(s)
	in.structs[name] = s
	in.Structs = append(in.Structs, s)
	return s
}

// structFor lowers a declared struct or class. One path for both: a class is
// byte-for-byte identical in layout to a struct and differs only in its
// member and method model, so nothing here branches on Class().
func (l *lowerer) structFor(t types.Type, in *Module) *Struct {
	name := l.typeName(t)
	owner := l.ownerOf(t, in)
	if s, ok := owner.structs[name]; ok {
		return s
	}
	st := types.AsStruct(t)
	if st == nil {
		return nil
	}

	// Declared before the fields are lowered, so a self-referential field
	// reaches its enclosing type without looping. This is the same two-step
	// types.Named itself uses, and it is why lower/vir has to re-sort the
	// struct section.
	s := &Struct{Name: name, Module: owner, Export: true}
	owner.structs[name] = s
	owner.Structs = append(owner.Structs, s)

	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		s.Fields = append(s.Fields, field(f.Name, l.typ(f.Type, owner)))
	}
	l.layout(s)
	return s
}

// enumStruct lays a payload enum out as tag + opaque payload bytes. Case
// bindings are a field.ptr plus a reinterpretation into the payload — views,
// not copies. The sources fix no enum layout; this is the implementation's
// choice, matching types.Sizes.enumSize so the two cannot disagree.
func (l *lowerer) enumStruct(t types.Type, e *types.Enum, in *Module) *Struct {
	name := l.typeName(t)
	owner := l.ownerOf(t, in)
	if s, ok := owner.structs[name]; ok {
		return s
	}
	tag := l.basic(e.Discriminant())
	total := l.conf.Sizes.Sizeof(t)
	tagSize := l.conf.Sizes.Sizeof(e.Discriminant())

	s := &Struct{
		Name:   name,
		Module: owner,
		Export: true,
		Fields: []StructField{
			field("tag", tag),
			field("payload", ArrayOf(U8, total-tagSize)),
		},
	}
	l.layout(s)
	owner.structs[name] = s
	owner.Structs = append(owner.Structs, s)
	return s
}

func (l *lowerer) tupleStruct(tp *types.Tuple, in *Module) *Struct {
	var parts []string
	fields := make([]StructField, 0, tp.Len())
	for i := 0; i < tp.Len(); i++ {
		ft := l.typ(tp.At(i).Type(), in)
		parts = append(parts, mangleType(ft))
		fields = append(fields, field("f"+itoa(i), ft))
	}
	return l.header(in, "_Vt_"+strings.Join(parts, "_"), fields...)
}

// layout fills offsets, size, and alignment from types.Sizes. Declaration
// order with padding and no reordering: a reader predicts the layout from
// the source, and interop assumes it.
func (l *lowerer) layout(s *Struct) {
	var cur, maxAlign int64 = 0, 1
	for i := range s.Fields {
		a := l.alignOf(s.Fields[i].Type)
		cur = align(cur, a)
		s.Fields[i].Offset = cur
		cur += l.sizeOf(s.Fields[i].Type)
		if a > maxAlign {
			maxAlign = a
		}
	}
	s.Size = align(cur, maxAlign)
	s.Align = maxAlign
}

func (l *lowerer) sizeOf(t *Type) int64 {
	switch t.Kind {
	case KVoid:
		return 0
	case KInt:
		if t.Bits <= 8 {
			return 1
		}
		return int64(t.Bits) / 8
	case KFloat:
		return int64(t.Bits) / 8
	case KPtr:
		return l.conf.Sizes.WordSize
	case KStruct:
		return t.Struct.Size
	case KArray:
		return t.Len * align(l.sizeOf(t.Elem), l.alignOf(t.Elem))
	case KVector:
		return t.Len * l.sizeOf(t.Elem)
	case KPredicate:
		// Never reaches storage, so it never has a layout to report.
		return 0
	}
	return 0
}

func (l *lowerer) alignOf(t *Type) int64 {
	switch t.Kind {
	case KStruct:
		return t.Struct.Align
	case KArray:
		return l.alignOf(t.Elem)
	case KVector:
		return roundPow2(l.sizeOf(t))
	case KPredicate:
		return 1
	}
	sz := l.sizeOf(t)
	if sz > l.conf.Sizes.MaxAlign {
		return l.conf.Sizes.MaxAlign
	}
	if sz == 0 {
		return 1
	}
	return sz
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