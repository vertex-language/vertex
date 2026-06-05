package compiler

import (
	"fmt"
	"strings"

	cir "github.com/vertex-language/ir/c"
)

// ─────────────────────────────────────────────────────────────────────────────
// VType — internal Vertex type hierarchy
// ─────────────────────────────────────────────────────────────────────────────

// VType is the internal representation of a resolved Vertex type.
type VType interface {
	vtypeNode()
	String() string
	// CIRType returns the ir/c Type for this VType.
	// Returns nil for types the lowerer must handle specially.
	CIRType() cir.Type
	Equal(VType) bool
}

// ─── Primitives ───────────────────────────────────────────────────────────────

type VInt struct {
	Bits   int  // 8 | 16 | 32 | 64
	Signed bool
}

type VFloat struct{ Bits int } // 32 | 64

type VBool struct{}
type VChar struct{}
type VVoid struct{}

// VString represents both immutable (let) and mutable (var) strings.
type VString struct{ Mutable bool }

// VNil is the type of the nil literal.
type VNil struct{}

// VUnknown is used during resolution for forward-referenced or unresolved names.
type VUnknown struct{ Name string }

// ─── Composite ────────────────────────────────────────────────────────────────

type VPointer struct {
	Elem    VType
	IsConst bool
}

type VOptional struct{ Elem VType }

// VFixedArray is a stack-allocated C array.
type VFixedArray struct {
	Elem VType
	Size int // -1 if size is unknown at resolution time
}

// VDynArray is a heap-allocated GLib GArray.
type VDynArray struct{ Elem VType }

type VTuple struct {
	Elems  []VType
	Labels []string // parallel; "" for unlabelled
}

type VFunc struct {
	Params []VType
	Return VType
}

type VResult struct {
	Ok  VType
	Err VType
}

type VChan struct{ Elem VType }

type VMap struct {
	Key   VType
	Value VType
}

// VExpected is the resolved type of an Expected(channel, value) test annotation.
// It carries no runtime representation; it is only used by the compiler to
// route test-function lowering and to store the expected output string.
type VExpected struct {
	Channel string // "stdout" | "exitCode"
	Value   string // expected output, e.g. "15"
}

// ─── Named user types ─────────────────────────────────────────────────────────

type VStruct struct {
	Name string
	Decl *StructDecl
}

type VClass struct {
	Name string
	Decl *ClassDecl
}

type VEnum struct {
	Name    string
	RawType VType // e.g. VInt{32, true}
	Decl    *EnumDecl
}

type VTypeAlias struct {
	Name       string
	Underlying VType
}

// VRange is the type of a range expression: lo..<hi or lo...hi.
// It is distinct from VDynArray — a range is not a GLib array.
type VRange struct{ Elem VType }

func (*VRange) vtypeNode()        {}
func (*VRange) String() string    { return "range" }
func (*VRange) CIRType() cir.Type { return nil }
func (*VRange) Equal(o VType) bool {
    u, ok := o.(*VRange)
    return ok && u.Elem.Equal(u.Elem)
}

// ─── vtypeNode markers ────────────────────────────────────────────────────────

func (*VInt) vtypeNode()       {}
func (*VFloat) vtypeNode()     {}
func (*VBool) vtypeNode()      {}
func (*VChar) vtypeNode()      {}
func (*VVoid) vtypeNode()      {}
func (*VString) vtypeNode()    {}
func (*VNil) vtypeNode()       {}
func (*VUnknown) vtypeNode()   {}
func (*VPointer) vtypeNode()   {}
func (*VOptional) vtypeNode()  {}
func (*VFixedArray) vtypeNode() {}
func (*VDynArray) vtypeNode()  {}
func (*VTuple) vtypeNode()     {}
func (*VFunc) vtypeNode()      {}
func (*VResult) vtypeNode()    {}
func (*VChan) vtypeNode()      {}
func (*VMap) vtypeNode()       {}
func (*VExpected) vtypeNode()  {}
func (*VStruct) vtypeNode()    {}
func (*VClass) vtypeNode()     {}
func (*VEnum) vtypeNode()      {}
func (*VTypeAlias) vtypeNode() {}

// ─── String representations ───────────────────────────────────────────────────

func (t *VInt) String() string {
	prefix := "int"
	if !t.Signed {
		prefix = "uint"
	}
	return fmt.Sprintf("%s%d", prefix, t.Bits)
}
func (t *VFloat) String() string     { return fmt.Sprintf("float%d", t.Bits) }
func (*VBool) String() string        { return "bool" }
func (*VChar) String() string        { return "char" }
func (*VVoid) String() string        { return "void" }
func (t *VString) String() string    { return "string" }
func (*VNil) String() string         { return "nil" }
func (t *VUnknown) String() string   { return "<unknown:" + t.Name + ">" }
func (t *VPointer) String() string {
	if t.IsConst {
		return "*const " + t.Elem.String()
	}
	return "*" + t.Elem.String()
}
func (t *VOptional) String() string   { return t.Elem.String() + "?" }
func (t *VFixedArray) String() string { return fmt.Sprintf("[%s](%d)", t.Elem, t.Size) }
func (t *VDynArray) String() string   { return "[" + t.Elem.String() + "]" }
func (t *VTuple) String() string {
	parts := make([]string, len(t.Elems))
	for i, e := range t.Elems {
		if i < len(t.Labels) && t.Labels[i] != "" {
			parts[i] = t.Labels[i] + ": " + e.String()
		} else {
			parts[i] = e.String()
		}
	}
	return "(" + strings.Join(parts, ", ") + ")"
}
func (t *VResult) String() string    { return fmt.Sprintf("Result(%s, %s)", t.Ok, t.Err) }
func (t *VChan) String() string      { return "chan " + t.Elem.String() }
func (t *VMap) String() string       { return fmt.Sprintf("map[%s]%s", t.Key, t.Value) }
func (t *VExpected) String() string  { return fmt.Sprintf("Expected(%s,%q)", t.Channel, t.Value) }
func (t *VStruct) String() string    { return t.Name }
func (t *VClass) String() string     { return t.Name }
func (t *VEnum) String() string      { return t.Name }
func (t *VTypeAlias) String() string { return t.Name }
func (t *VFunc) String() string      { return "func" }

// ─── CIRType — mapping to ir/c types ─────────────────────────────────────────

func (t *VInt) CIRType() cir.Type {
	if t.Signed {
		switch t.Bits {
		case 8:
			return cir.Int8
		case 16:
			return cir.Int16
		case 32:
			return cir.Int32
		case 64:
			return cir.Int64
		}
	} else {
		switch t.Bits {
		case 8:
			return cir.UInt8
		case 16:
			return cir.UInt16
		case 32:
			return cir.UInt32
		case 64:
			return cir.UInt64
		}
	}
	return cir.Int32
}

func (t *VFloat) CIRType() cir.Type {
	if t.Bits == 32 {
		return cir.Float32
	}
	return cir.Float64
}

func (*VBool) CIRType() cir.Type    { return cir.Bool }
func (*VChar) CIRType() cir.Type    { return cir.Char }
func (*VVoid) CIRType() cir.Type    { return cir.Void }
func (*VNil) CIRType() cir.Type     { return cir.VoidPtr }
func (*VUnknown) CIRType() cir.Type { return cir.VoidPtr }

func (t *VString) CIRType() cir.Type {
	if t.Mutable {
		return nil // lowerer emits GString* specially
	}
	return cir.ConstPtr(cir.Char)
}

func (t *VPointer) CIRType() cir.Type {
	if t.IsConst {
		return cir.ConstPtr(t.Elem.CIRType())
	}
	return cir.Ptr(t.Elem.CIRType())
}

func (t *VOptional) CIRType() cir.Type {
	if inner := t.Elem.CIRType(); inner != nil {
		return cir.Ptr(inner)
	}
	return cir.VoidPtr
}

func (t *VFixedArray) CIRType() cir.Type {
	if t.Size <= 0 {
		return nil
	}
	elem := t.Elem.CIRType()
	if elem == nil {
		return nil
	}
	return cir.Array(elem, t.Size)
}

// Dynamic arrays, classes, maps, tuples, results, channels → lowerer handles.
func (*VDynArray) CIRType() cir.Type   { return nil }
func (*VTuple) CIRType() cir.Type      { return nil }
func (*VFunc) CIRType() cir.Type       { return nil }
func (*VResult) CIRType() cir.Type     { return nil }
func (*VChan) CIRType() cir.Type       { return nil }
func (*VMap) CIRType() cir.Type        { return nil } // lowerer handles as GHashTable
func (*VStruct) CIRType() cir.Type     { return nil } // lowerer looks up cached *StructType
func (*VClass) CIRType() cir.Type      { return nil }
func (*VEnum) CIRType() cir.Type       { return cir.Int32 }

// VExpected carries no runtime value — it is a compile-time annotation only.
func (*VExpected) CIRType() cir.Type { return cir.Void }

func (t *VTypeAlias) CIRType() cir.Type { return t.Underlying.CIRType() }

// ─── Equal ────────────────────────────────────────────────────────────────────

func (t *VInt) Equal(o VType) bool {
	u, ok := o.(*VInt)
	return ok && t.Bits == u.Bits && t.Signed == u.Signed
}
func (t *VFloat) Equal(o VType) bool {
	u, ok := o.(*VFloat)
	return ok && t.Bits == u.Bits
}
func (*VBool) Equal(o VType) bool    { _, ok := o.(*VBool); return ok }
func (*VChar) Equal(o VType) bool    { _, ok := o.(*VChar); return ok }
func (*VVoid) Equal(o VType) bool    { _, ok := o.(*VVoid); return ok }
func (*VNil) Equal(o VType) bool     { _, ok := o.(*VNil); return ok }
func (*VUnknown) Equal(o VType) bool { return false }
func (t *VString) Equal(o VType) bool {
	u, ok := o.(*VString)
	return ok && t.Mutable == u.Mutable
}
func (t *VPointer) Equal(o VType) bool {
	u, ok := o.(*VPointer)
	return ok && t.IsConst == u.IsConst && t.Elem.Equal(u.Elem)
}
func (t *VOptional) Equal(o VType) bool {
	u, ok := o.(*VOptional)
	return ok && t.Elem.Equal(u.Elem)
}
func (t *VFixedArray) Equal(o VType) bool {
	u, ok := o.(*VFixedArray)
	return ok && t.Size == u.Size && t.Elem.Equal(u.Elem)
}
func (t *VDynArray) Equal(o VType) bool {
	u, ok := o.(*VDynArray)
	return ok && t.Elem.Equal(u.Elem)
}
func (t *VMap) Equal(o VType) bool {
	u, ok := o.(*VMap)
	return ok && t.Key.Equal(u.Key) && t.Value.Equal(u.Value)
}
func (t *VExpected) Equal(o VType) bool {
	u, ok := o.(*VExpected)
	return ok && t.Channel == u.Channel && t.Value == u.Value
}
func (t *VStruct) Equal(o VType) bool {
	u, ok := o.(*VStruct)
	return ok && t.Name == u.Name
}
func (t *VClass) Equal(o VType) bool {
	u, ok := o.(*VClass)
	return ok && t.Name == u.Name
}
func (t *VEnum) Equal(o VType) bool {
	u, ok := o.(*VEnum)
	return ok && t.Name == u.Name
}
func (t *VTypeAlias) Equal(o VType) bool {
	u, ok := o.(*VTypeAlias)
	return ok && t.Name == u.Name
}
func (t *VResult) Equal(o VType) bool {
	u, ok := o.(*VResult)
	return ok && t.Ok.Equal(u.Ok) && t.Err.Equal(u.Err)
}
func (t *VTuple) Equal(o VType) bool {
	u, ok := o.(*VTuple)
	if !ok || len(t.Elems) != len(u.Elems) {
		return false
	}
	for i := range t.Elems {
		if !t.Elems[i].Equal(u.Elems[i]) {
			return false
		}
	}
	return true
}
func (t *VFunc) Equal(o VType) bool { return false } // structural equality not needed yet
func (*VChan) Equal(o VType) bool   { return false }

// ─── Built-in type name table ─────────────────────────────────────────────────

// BuiltinTypes maps Vertex built-in type names to their VType.
var BuiltinTypes = map[string]VType{
	"int":     &VInt{Bits: 32, Signed: true},
	"int8":    &VInt{Bits: 8, Signed: true},
	"int16":   &VInt{Bits: 16, Signed: true},
	"int32":   &VInt{Bits: 32, Signed: true},
	"int64":   &VInt{Bits: 64, Signed: true},
	"uint":    &VInt{Bits: 32, Signed: false},
	"uint8":   &VInt{Bits: 8, Signed: false},
	"uint16":  &VInt{Bits: 16, Signed: false},
	"uint32":  &VInt{Bits: 32, Signed: false},
	"uint64":  &VInt{Bits: 64, Signed: false},
	"float":   &VFloat{Bits: 32},
	"float32": &VFloat{Bits: 32},
	"float64": &VFloat{Bits: 64},
	"bool":    &VBool{},
	"char":    &VChar{},
	"string":  &VString{},
	"void":    &VVoid{},
}

// IsBuiltinType reports whether name is a built-in scalar type.
func IsBuiltinType(name string) bool {
	_, ok := BuiltinTypes[name]
	return ok
}