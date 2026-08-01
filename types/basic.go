package types

// BasicKind enumerates A.1.4's PredeclaredTypeNames plus the untyped kinds a
// literal carries before it reaches a typed position (A.1.5.1).
type BasicKind int

const (
	Invalid BasicKind = iota

	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	Char
	String

	// Untyped kinds. A.1.5.1 ⊢ "An integer literal is untyped until it reaches
	// a typed position, where it takes that position's type." These never
	// appear in a declared signature; they exist only between the literal and
	// its destination.
	UntypedBool
	UntypedInt
	UntypedFloat
	UntypedChar
	UntypedString

	// UntypedNil is not a general value. A.5.2 ⊢ nil "is not a general value
	// and has no type of its own" — it appears only against typed_ptr and as
	// the map-erase operand.
	UntypedNil
)

// BasicInfo is the property bit set of a *Basic.
//
// The flags are named Info* rather than Is* because predicates.go exports
// IsInteger, IsFloat, IsNumeric, IsString, IsOrdered, and IsUntyped as
// functions over Type, and Go gives constants and functions one namespace per
// package. The functions are the API analyzer and lower call; these flags are
// the internal representation those functions read, so the flags are what
// carries the distinguishing prefix.
type BasicInfo uint

const (
	InfoBoolean BasicInfo = 1 << iota
	InfoInteger
	InfoUnsigned
	InfoFloat
	InfoChar
	InfoString
	InfoUntyped

	InfoNumeric = InfoInteger | InfoFloat
	InfoOrdered = InfoInteger | InfoFloat | InfoString | InfoChar
)

type Basic struct {
	kind BasicKind
	info BasicInfo
	name string
}

func (b *Basic) Kind() BasicKind  { return b.kind }
func (b *Basic) Info() BasicInfo  { return b.info }
func (b *Basic) Name() string     { return b.name }
func (b *Basic) Underlying() Type { return b }
func (b *Basic) String() string   { return b.name }

func (b *Basic) is(i BasicInfo) bool { return b != nil && b.info&i != 0 }

// Typ holds the singleton for each BasicKind, indexed by kind.
//
// There is exactly one object per kind, so identity comparison on *Basic is
// type identity — which is what makes A.1.4's byte/uint8 rule fall out. The
// spec says they "denote the same type; no conversion is required or permitted
// between them, in either direction", so `byte` is not a distinct entry here.
// The universe scope binds both spellings to Typ[Uint8].
//
// The cost is that a diagnostic about a value the user wrote as `byte` says
// "uint8". That is the correct trade: inventing a second object to preserve the
// spelling would make Identical() lie.
var Typ = []*Basic{
	Invalid: {Invalid, 0, "invalid type"},

	Bool:    {Bool, InfoBoolean, "bool"},
	Int:     {Int, InfoInteger, "int"},
	Int8:    {Int8, InfoInteger, "int8"},
	Int16:   {Int16, InfoInteger, "int16"},
	Int32:   {Int32, InfoInteger, "int32"},
	Int64:   {Int64, InfoInteger, "int64"},
	Uint:    {Uint, InfoInteger | InfoUnsigned, "uint"},
	Uint8:   {Uint8, InfoInteger | InfoUnsigned, "uint8"},
	Uint16:  {Uint16, InfoInteger | InfoUnsigned, "uint16"},
	Uint32:  {Uint32, InfoInteger | InfoUnsigned, "uint32"},
	Uint64:  {Uint64, InfoInteger | InfoUnsigned, "uint64"},
	Float32: {Float32, InfoFloat, "float32"},
	Float64: {Float64, InfoFloat, "float64"},

	// A.1.5.2 ⊢ a CharLiteral "denotes exactly one Unicode scalar value, held
	// in 4 bytes. 'A' and "A" are different types and never interconvert
	// implicitly." Char is therefore its own kind, not an alias for int32.
	Char:   {Char, InfoChar, "char"},
	String: {String, InfoString, "string"},

	UntypedBool:   {UntypedBool, InfoBoolean | InfoUntyped, "untyped bool"},
	UntypedInt:    {UntypedInt, InfoInteger | InfoUntyped, "untyped int"},
	UntypedFloat:  {UntypedFloat, InfoFloat | InfoUntyped, "untyped float"},
	UntypedChar:   {UntypedChar, InfoChar | InfoUntyped, "untyped char"},
	UntypedString: {UntypedString, InfoString | InfoUntyped, "untyped string"},
	UntypedNil:    {UntypedNil, InfoUntyped, "untyped nil"},
}

// Predeclared names, in the order A.1.4 lists them. `byte` and `uint8` share an
// entry deliberately; see Typ.
var predeclared = map[string]*Basic{
	"bool":    Typ[Bool],
	"int":     Typ[Int],
	"int8":    Typ[Int8],
	"int16":   Typ[Int16],
	"int32":   Typ[Int32],
	"int64":   Typ[Int64],
	"uint":    Typ[Uint],
	"uint8":   Typ[Uint8],
	"uint16":  Typ[Uint16],
	"uint32":  Typ[Uint32],
	"uint64":  Typ[Uint64],
	"byte":    Typ[Uint8],
	"float32": Typ[Float32],
	"float64": Typ[Float64],
	"char":    Typ[Char],
	"string":  Typ[String],
}

// LookupPredeclared returns the basic type for a PredeclaredTypeName, or nil.
// These are ordinary identifiers pre-bound in an implicit scope (A.1.4), so the
// resolver consults this rather than the scanner recognizing them.
func LookupPredeclared(name string) *Basic { return predeclared[name] }

// Default returns the type an untyped constant takes when it reaches a position
// that imposes none. A.1.5.1 makes this rare — most literals land in a typed
// position — but a bare `let x = 1` needs an answer.
func Default(t Type) Type {
	if b, ok := t.(*Basic); ok {
		switch b.kind {
		case UntypedBool:
			return Typ[Bool]
		case UntypedInt:
			return Typ[Int]
		case UntypedFloat:
			return Typ[Float64]
		case UntypedChar:
			return Typ[Char]
		case UntypedString:
			return Typ[String]
		}
	}
	return t
}