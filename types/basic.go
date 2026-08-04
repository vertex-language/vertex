package types

// BasicKind enumerates the PredeclaredTypeNames, the predeclared tensor element
// type names, and the untyped kinds a literal carries before it lands.
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

	// Tensor element types. Ordinary identifiers in the same implicit scope as
	// the type names, but with their own legality rule: §2.3 ⊢ legal "only
	// inside an `npu` body". They are Basics because they are scalar element
	// types with their own widths, and IsTensorElem is what gates them.
	BF16
	FP8E4M3
	FP8E5M2
	Int4

	// Untyped kinds. §4.1 ⊢ "a literal has no type until it lands; where a
	// destination type exists the literal takes it". These never appear in a
	// declared signature; they exist only between a literal and its
	// destination.
	UntypedBool
	UntypedInt
	UntypedFloat
	UntypedChar
	UntypedString

	// UntypedNil is not a general value. §10 ⊢ "there is no optional type, no
	// propagation operator, and no general `nil`. `nil` belongs to
	// `typed_ptr T` and to nothing else."
	UntypedNil
)

// BasicInfo is the property bit set of a *Basic.
//
// The flags are named Info* rather than Is* because predicates.go exports
// IsInteger, IsFloat, IsNumeric, IsString, IsOrdered, and IsUntyped as
// functions over Type, and Go gives constants and functions one namespace per
// package. The functions are the API the analyzer and lower call; these flags
// are the representation those functions read.
type BasicInfo uint

const (
	InfoBoolean BasicInfo = 1 << iota
	InfoInteger
	InfoUnsigned
	InfoFloat
	InfoChar
	InfoString
	InfoUntyped

	// InfoTensorElem marks the four names of §2.3's tensor-element family. It
	// is carried alongside the family bit (bf16 and the fp8 pair are floats,
	// int4 is an integer) so arithmetic inside an npu body reads the same
	// predicates as anywhere else, while ConvertibleTo can still single them
	// out for §4.2's constructor-spelling exception.
	InfoTensorElem

	InfoNumeric = InfoInteger | InfoFloat

	// §3.5 ⊢ Ordered is "numerics and `string`". `char` is deliberately absent:
	// it is comparable but not ordered, so `'a' < 'b'` is an error.
	InfoOrdered = InfoNumeric | InfoString
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
// type identity — which is what makes §2.3's byte rule fall out. ⊢ "`byte` is
// an alias for `uint8`, not a distinct type", so there is no `byte` entry here
// and the universe scope binds both spellings to Typ[Uint8].
//
// The cost is that a diagnostic about a value the user wrote as `byte` says
// "uint8". That is the correct trade: a second object to preserve the spelling
// would make Identical lie.
var Typ = []*Basic{
	Invalid: {Invalid, 0, "invalid type"},

	Bool: {Bool, InfoBoolean, "bool"},

	// §2.3 ⊢ "`int` and `uint` are the target's pointer width and are distinct
	// types from `int64`/`uint64` even where the widths agree." Distinctness is
	// this entry; the width is Sizes.WordSize (sizes.go).
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

	// §2.3 ⊢ "`char` is one Unicode scalar value." It is its own kind, not an
	// alias for int32: grammar.md keeps char_value and string_value separate,
	// and §4.1 ⊢ 'A' and "A" "never interconvert implicitly".
	Char:   {Char, InfoChar, "char"},
	String: {String, InfoString, "string"},

	BF16:    {BF16, InfoFloat | InfoTensorElem, "bf16"},
	FP8E4M3: {FP8E4M3, InfoFloat | InfoTensorElem, "fp8e4m3"},
	FP8E5M2: {FP8E5M2, InfoFloat | InfoTensorElem, "fp8e5m2"},
	Int4:    {Int4, InfoInteger | InfoTensorElem, "int4"},

	UntypedBool:   {UntypedBool, InfoBoolean | InfoUntyped, "untyped bool"},
	UntypedInt:    {UntypedInt, InfoInteger | InfoUntyped, "untyped int"},
	UntypedFloat:  {UntypedFloat, InfoFloat | InfoUntyped, "untyped float"},
	UntypedChar:   {UntypedChar, InfoChar | InfoUntyped, "untyped char"},
	UntypedString: {UntypedString, InfoString | InfoUntyped, "untyped string"},
	UntypedNil:    {UntypedNil, InfoUntyped, "untyped nil"},
}

// predeclaredTypes are the PredeclaredTypeNames, in the order grammar.md lists
// them. `byte` and `uint8` share an entry deliberately; see Typ.
var predeclaredTypes = map[string]*Basic{
	"int": Typ[Int], "int8": Typ[Int8], "int16": Typ[Int16],
	"int32": Typ[Int32], "int64": Typ[Int64],
	"uint": Typ[Uint], "uint8": Typ[Uint8], "uint16": Typ[Uint16],
	"uint32": Typ[Uint32], "uint64": Typ[Uint64],
	"byte":    Typ[Uint8],
	"float32": Typ[Float32], "float64": Typ[Float64],
	"bool": Typ[Bool], "char": Typ[Char], "string": Typ[String],
}

// predeclaredTensorElems are the predeclared tensor element type names. They
// are a separate table from predeclaredTypes because they are a separate family
// with a separate legality rule (§2.3), not because the scanner tells them
// apart — it does not, and neither does Universe, which holds both.
var predeclaredTensorElems = map[string]*Basic{
	"bf16": Typ[BF16], "fp8e4m3": Typ[FP8E4M3],
	"fp8e5m2": Typ[FP8E5M2], "int4": Typ[Int4],
}

// LookupPredeclared returns the basic type for a PredeclaredTypeName, or nil.
func LookupPredeclared(name string) *Basic { return predeclaredTypes[name] }

// LookupTensorElem returns the basic type for a predeclared tensor element type
// name, or nil. The analyzer consults this to decide whether a resolved name
// carries §2.3's npu-body restriction.
func LookupTensorElem(name string) *Basic { return predeclaredTensorElems[name] }

// Default returns the type an untyped constant takes when it reaches a position
// that imposes none. §4.1 makes this rare — most literals land somewhere typed
// — but a bare `let x = 1` needs an answer.
//
// UntypedNil has no default: §10 gives it no type of its own, so it stays
// untyped and the analyzer rejects a destination that is not a typed_ptr.
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