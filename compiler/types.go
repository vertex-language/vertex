package compiler

import "fmt"

type Kind int

const (
	KindInvalid Kind = iota
	KindVoid
	KindBool
	KindInt
	KindInt8
	KindInt16
	KindInt32
	KindInt64
	KindUint
	KindUint8
	KindUint16
	KindUint32
	KindUint64
	KindFloat
	KindDouble
	KindString
	KindChar
	KindArray
	KindDict
	KindOptional
	KindResult
	KindTuple
	KindFunc
	KindPointer // any T  /  mut any T
	KindOpaque
	KindChannel
	KindStruct
	KindClass
	KindEnum
	KindNamed // unresolved forward reference
)

type WasmKind int

const (
	WasmI32 WasmKind = iota
	WasmI64
	WasmF32
	WasmF64
	WasmNone
)

type Type interface {
	Kind() Kind
	Wasm() WasmKind
	String() string
}

// ── Primitive ────────────────────────────────────────────────────────────────

type primitive struct{ k Kind }

var (
	Void    Type = &primitive{KindVoid}
	Bool    Type = &primitive{KindBool}
	Int     Type = &primitive{KindInt}
	Int8    Type = &primitive{KindInt8}
	Int16   Type = &primitive{KindInt16}
	Int32   Type = &primitive{KindInt32}
	Int64   Type = &primitive{KindInt64}
	Uint    Type = &primitive{KindUint}
	Uint8   Type = &primitive{KindUint8}
	Uint16  Type = &primitive{KindUint16}
	Uint32  Type = &primitive{KindUint32}
	Uint64  Type = &primitive{KindUint64}
	Float   Type = &primitive{KindFloat}
	Double  Type = &primitive{KindDouble}
	StrType Type = &primitive{KindString}
	Char    Type = &primitive{KindChar}
	Opaque  Type = &primitive{KindOpaque}
)

func (p *primitive) Kind() Kind { return p.k }
func (p *primitive) Wasm() WasmKind {
	switch p.k {
	case KindVoid:
		return WasmNone
	case KindInt64, KindUint64:
		return WasmI64
	case KindFloat:
		return WasmF32
	case KindDouble:
		return WasmF64
	default:
		return WasmI32
	}
}
func (p *primitive) String() string {
	names := map[Kind]string{
		KindVoid: "void", KindBool: "bool",
		KindInt: "int", KindInt8: "int8", KindInt16: "int16",
		KindInt32: "int32", KindInt64: "int64",
		KindUint: "uint", KindUint8: "uint8", KindUint16: "uint16",
		KindUint32: "uint32", KindUint64: "uint64",
		KindFloat: "float", KindDouble: "double",
		KindString: "string", KindChar: "char", KindOpaque: "opaque",
	}
	if s, ok := names[p.k]; ok {
		return s
	}
	return "?"
}

// ── Pointer ──────────────────────────────────────────────────────────────────

type PointerType struct {
	Elem    Type
	Mutable bool
}

func (p *PointerType) Kind() Kind      { return KindPointer }
func (p *PointerType) Wasm() WasmKind { return WasmI32 }
func (p *PointerType) String() string {
	if p.Mutable {
		return "mut any " + p.Elem.String()
	}
	return "any " + p.Elem.String()
}

// ── Array ────────────────────────────────────────────────────────────────────

type ArrayType struct{ Elem Type }

func (a *ArrayType) Kind() Kind      { return KindArray }
func (a *ArrayType) Wasm() WasmKind { return WasmI32 }
func (a *ArrayType) String() string  { return "[" + a.Elem.String() + "]" }

// ── Optional ─────────────────────────────────────────────────────────────────

type OptionalType struct{ Elem Type }

func (o *OptionalType) Kind() Kind      { return KindOptional }
func (o *OptionalType) Wasm() WasmKind { return o.Elem.Wasm() }
func (o *OptionalType) String() string  { return o.Elem.String() + "?" }

// ── Named (forward ref) ──────────────────────────────────────────────────────

type NamedType struct{ Name string }

func (n *NamedType) Kind() Kind      { return KindNamed }
func (n *NamedType) Wasm() WasmKind { return WasmI32 }
func (n *NamedType) String() string  { return n.Name }

// ── Struct ───────────────────────────────────────────────────────────────────

type StructField struct {
	Name    string
	Type    Type
	Mutable bool
	Offset  int
}

type StructType struct {
	Name   string
	Fields []*StructField
	Size   int
}

func (s *StructType) Kind() Kind      { return KindStruct }
func (s *StructType) Wasm() WasmKind { return WasmI32 }
func (s *StructType) String() string  { return s.Name }
func (s *StructType) Field(name string) *StructField {
	for _, f := range s.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// ── Class ────────────────────────────────────────────────────────────────────

type ClassType struct {
	Name   string
	Fields []*StructField
	Size   int

	// Native = true when the class has a ': parentName' and is a zero-size
	// compile-time dispatch surface. The backend removes it entirely.
	Native bool
	// Parent is the namespace from the import path, e.g. "sdl2", "syscalls".
	Parent string
}

func (c *ClassType) Kind() Kind      { return KindClass }
func (c *ClassType) Wasm() WasmKind { return WasmI32 }
func (c *ClassType) String() string  { return c.Name }

// ── Enum ─────────────────────────────────────────────────────────────────────

type EnumCase struct {
	Name   string
	IntVal int64
	StrVal string
}

type EnumType struct {
	Name    string
	RawKind Kind
	Cases   []*EnumCase
}

func (e *EnumType) Kind() Kind      { return KindEnum }
func (e *EnumType) Wasm() WasmKind { return WasmI32 }
func (e *EnumType) String() string  { return e.Name }

// ── FuncSig ──────────────────────────────────────────────────────────────────

type FuncQual int

const (
	QualNone    FuncQual = iota
	QualAsync
	QualThread
	QualProcess
	QualGPU
)

type FuncSig struct {
	Params []Type
	Ret    Type
	Qual   FuncQual
	Muts   []bool
}

type FuncType struct{ Sig *FuncSig }

func (f *FuncType) Kind() Kind      { return KindFunc }
func (f *FuncType) Wasm() WasmKind { return WasmI32 }
func (f *FuncType) String() string  { return "func(...)" }

// ── Result ───────────────────────────────────────────────────────────────────

type ResultType struct {
	Ok  Type
	Err Type
}

func (r *ResultType) Kind() Kind      { return KindResult }
func (r *ResultType) Wasm() WasmKind { return WasmI32 }
func (r *ResultType) String() string {
	return fmt.Sprintf("Result(%s, %s)", r.Ok, r.Err)
}

// ── Channel ──────────────────────────────────────────────────────────────────

type ChannelType struct{ Elem Type }

func (c *ChannelType) Kind() Kind      { return KindChannel }
func (c *ChannelType) Wasm() WasmKind { return WasmI32 }
func (c *ChannelType) String() string  { return "channel " + c.Elem.String() }

// ── Layout helpers ───────────────────────────────────────────────────────────

// SizeOf returns the byte size of a type in linear memory.
// Struct and class types return their actual computed size so that nested
// value-type fields are laid out correctly by LayoutStruct.
func SizeOf(t Type) int {
	switch t.Kind() {
	case KindVoid:
		return 0
	case KindBool, KindInt8, KindUint8:
		return 1
	case KindInt16, KindUint16:
		return 2
	case KindInt64, KindUint64, KindDouble:
		return 8
	case KindStruct:
		if st, ok := t.(*StructType); ok && st.Size > 0 {
			return st.Size
		}
		return 4
	case KindClass:
		if ct, ok := t.(*ClassType); ok && ct.Size > 0 {
			return ct.Size
		}
		return 4
	default:
		return 4
	}
}

// AlignOf returns the required alignment of a type.
// For structs the alignment is the maximum alignment of any field.
func AlignOf(t Type) int {
	switch t.Kind() {
	case KindBool, KindInt8, KindUint8:
		return 1
	case KindInt16, KindUint16:
		return 2
	case KindInt64, KindUint64, KindDouble:
		return 8
	case KindStruct:
		if st, ok := t.(*StructType); ok {
			max := 1
			for _, f := range st.Fields {
				if a := AlignOf(f.Type); a > max {
					max = a
				}
			}
			return max
		}
		return 4
	default:
		return 4
	}
}

func LayoutStruct(fields []*StructField) int {
	off := 0
	maxAlign := 1
	for _, f := range fields {
		a := AlignOf(f.Type)
		if a > maxAlign {
			maxAlign = a
		}
		if rem := off % a; rem != 0 {
			off += a - rem
		}
		f.Offset = off
		off += SizeOf(f.Type)
	}
	if rem := off % maxAlign; rem != 0 {
		off += maxAlign - rem
	}
	return off
}