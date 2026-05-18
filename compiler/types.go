package compiler

import "github.com/vertex-language/wasm-compiler/wasm"

// VType is the interface every Vertex type implements.
type VType interface {
	IsVoid() bool
	WasmType() (wasm.ValType, bool)
	Size() uint32  // byte size in linear memory
	String() string
}

// ── Primitive types ───────────────────────────────────────────────────────────

type PrimitiveKind int

const (
	KindVoid   PrimitiveKind = iota
	KindBool                 // i32 0/1
	KindInt                  // i32 — Vertex's default integer
	KindInt64                // i64
	KindUInt                 // i32 unsigned
	KindFloat                // f32
	KindDouble               // f64
	KindString               // i32 pointer to NUL-terminated bytes
)

type PrimitiveType struct{ Kind PrimitiveKind }

func (t *PrimitiveType) IsVoid() bool { return t.Kind == KindVoid }

func (t *PrimitiveType) WasmType() (wasm.ValType, bool) {
	switch t.Kind {
	case KindVoid:
		return 0, false
	case KindBool, KindInt, KindUInt, KindString:
		return wasm.I32, true
	case KindInt64:
		return wasm.I64, true
	case KindFloat:
		return wasm.F32, true
	case KindDouble:
		return wasm.F64, true
	}
	return wasm.I32, true
}

func (t *PrimitiveType) Size() uint32 {
	switch t.Kind {
	case KindVoid:
		return 0
	case KindBool:
		return 1
	case KindInt, KindUInt, KindFloat, KindString:
		return 4
	case KindInt64, KindDouble:
		return 8
	}
	return 4
}

func (t *PrimitiveType) String() string {
	switch t.Kind {
	case KindVoid:
		return "Void"
	case KindBool:
		return "Bool"
	case KindInt:
		return "Int"
	case KindInt64:
		return "Int64"
	case KindUInt:
		return "UInt"
	case KindFloat:
		return "Float"
	case KindDouble:
		return "Double"
	case KindString:
		return "String"
	}
	return "unknown"
}

// ── Struct / class types ──────────────────────────────────────────────────────

type StructType struct {
	Name   string
	Fields []StructField
	IsRef  bool // true = class (heap reference), false = struct (value)
}

type StructField struct {
	Name   string
	Type   VType
	Offset uint32
}

func (t *StructType) IsVoid() bool                   { return false }
func (t *StructType) WasmType() (wasm.ValType, bool) { return wasm.I32, true }
func (t *StructType) String() string                 { return t.Name }
func (t *StructType) Size() uint32 {
	if t.IsRef {
		return 4 // class instance = 32-bit heap pointer
	}
	var s uint32
	for _, f := range t.Fields {
		s += f.Type.Size()
	}
	return s
}

// FieldOffset returns the byte offset of named field, or -1.
func (t *StructType) FieldOffset(name string) (VType, uint32, bool) {
	for _, f := range t.Fields {
		if f.Name == name {
			return f.Type, f.Offset, true
		}
	}
	return nil, 0, false
}

// ── Optional ──────────────────────────────────────────────────────────────────

type OptionalType struct{ Inner VType }

func (t *OptionalType) IsVoid() bool                   { return false }
func (t *OptionalType) WasmType() (wasm.ValType, bool) { return wasm.I32, true }
func (t *OptionalType) Size() uint32                   { return 4 }
func (t *OptionalType) String() string                 { return t.Inner.String() + "?" }

// ── Builtin type map ──────────────────────────────────────────────────────────

var builtinTypes = map[string]VType{
	"Void":    &PrimitiveType{Kind: KindVoid},
	"Bool":    &PrimitiveType{Kind: KindBool},
	"Int":     &PrimitiveType{Kind: KindInt},
	"Int32":   &PrimitiveType{Kind: KindInt},
	"Int64":   &PrimitiveType{Kind: KindInt64},
	"UInt":    &PrimitiveType{Kind: KindUInt},
	"UInt32":  &PrimitiveType{Kind: KindUInt},
	"Float":   &PrimitiveType{Kind: KindFloat},
	"Float32": &PrimitiveType{Kind: KindFloat},
	"Double":  &PrimitiveType{Kind: KindDouble},
	"Float64": &PrimitiveType{Kind: KindDouble},
	"String":  &PrimitiveType{Kind: KindString},
}