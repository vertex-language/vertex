package diag

import "fmt"

// Code identifies a diagnostic rule. It is the stable handle: renumbering one
// is a breaking change to the language's test corpus, because A.12.2's
// Expected(error, "...") compares rendered message text as specification.
//
// Numeric values are assigned explicitly, never by iota, so inserting a code
// cannot shift another. Ranges follow A.14's categories:
//
//	0xxx  internal
//	1xxx  lexical                       (A.1)
//	2xxx  syntactic                     (A.2–A.5)
//	3xxx  declarations and names        (A.6)
//	4xxx  types                         (A.3)
//	5xxx  ownership and exclusivity     (A.9)
//	6xxx  generics                      (A.7)
//	7xxx  pointers and memory           (A.4.8)
//	8xxx  interop                       (A.8)
//	9xxx  concurrency, devices, testing (A.10–A.12)
//
// Gaps within a range are deliberate and reserved. Ranges 5xxx, 6xxx, 8xxx, and
// 9xxx remain almost entirely unpopulated because the analyzer passes that
// raise them do not exist yet; A.14 is the inventory they will be drawn from.
type Code int

// ------------------------------------------------------------------ internal

const (
	// Internal marks a broken compiler invariant. Never produced by valid or
	// invalid source; if a user sees one, it is a bug in the compiler.
	Internal Code = 1
)

// ------------------------------------------------------------------- lexical

// Lexical diagnostics (A.1). The scanner is the only producer.
const (
	IllegalCharacter Code = 1001
	IllegalUTF8      Code = 1002
	DollarInIdent    Code = 1003 // A.1.2 ⊢ `$` is not an identifier character

	UnterminatedComment Code = 1010 // A.1.1 ⊢ does not nest; the first */ closes

	UnterminatedString    Code = 1020
	NewlineInString       Code = 1021 // DoubleStringCharacter excludes LineTerminator
	UnterminatedRawString Code = 1022
	UnterminatedChar      Code = 1023
	EmptyChar             Code = 1024
	CharTooLong           Code = 1025 // A.1.5.2 ⊢ exactly one Unicode scalar value
	InvalidEscape         Code = 1026
	InvalidHexEscape      Code = 1027 // \xHH
	InvalidUnicodeEscape  Code = 1028 // \u{HexDigits}
	UnicodeEscapeRange    Code = 1029 // > 10FFFF, or a surrogate

	SeparatorLeads       Code = 1040 // A.1.5.1 ⊢ `_` may not lead a digit run
	SeparatorTrails      Code = 1041 // ...nor trail one
	SeparatorDoubled     Code = 1042 // ...nor be doubled
	EmptyDigits          Code = 1043 // `0x` with no digits
	DigitOutOfRange      Code = 1044 // `0b2`, `0o9`
	HexFloatNoExponent   Code = 1045 // A.1.5.1 ⊢ `0xC.3` is not a literal
	MissingExponentDigit Code = 1046 // `1e`, `0x1p`
	NumberJoinedToIdent  Code = 1047 // `123abc`
)

// ----------------------------------------------------------------- syntactic

// Syntactic diagnostics (A.2–A.5). The parser is the only producer.
//
// This range is deliberately small. Of the four context parameters in A.0.2,
// only Lit is enforced during parsing; Await, Npu, and Own name constructs that
// A.14 lists as parsing and being rejected later, so their diagnostics are
// static-rule codes raised by the analyzer, not syntax errors raised here.
const (
	ExpectedToken   Code = 2001
	ExpectedStmtEnd Code = 2002
	ExpectedExpr    Code = 2003
	ExpectedType    Code = 2004
	ExpectedIdent   Code = 2005
	ExpectedDecl    Code = 2006
	ExpectedCase    Code = 2007
	ExpectedCall    Code = 2008 // A.4.2 — a launch prefix modifies a call

	RangeNotAssociative Code = 2010 // A.4.5 ⊢ `a..b..c` is a compile error
	LiteralInHeader     Code = 2011 // A.4.7 — parenthesize it
	TupleNeedsComma     Code = 2012 // A.4.7 ⊢ one-element tuple literal

	MissingPackageClause Code = 2020 // A.2 ⊢ mandatory, first non-comment
	MisplacedBuildClause Code = 2021 // A.2.2 — the second construct in the file
	UnknownBuildTag      Code = 2022 // A.2.2 ⊢ an error, not a silent exclusion
	ImportAfterDecl      Code = 2023 // A.2 — imports precede declarations
	DeferNotCall         Code = 2024 // A.5.8 ⊢ defer takes a call and nothing else
)

// ------------------------------------------------- declarations and names

// Declarations and names (A.6). The analyzer's resolve pass is the producer.
//
// Every code here answers "what does this name denote", which is the question
// resolution exists to answer and the one several of A.0.2's deferred
// ambiguities reduce to: A.3.6's Stack[int32] versus a[i], A.7.2's single
// identifier as both a TypeSet and a ConstraintName.
const (
	UndeclaredName       Code = 3001
	DuplicateDeclaration Code = 3002
	ShadowedBuiltin      Code = 3003 // A.1.4 ⊢ a ReservedBuiltinName may not be shadowed
	ReservedAsMember     Code = 3004 // A.1.4 ✗ func (w: var Widget) new() { }
	NotAType             Code = 3005
	ConstraintAsType     Code = 3006 // A.7.2 ✗ var c: Ordered
	NotAConstraint       Code = 3007
	TypeCycle            Code = 3008
	BlankNotUsable       Code = 3009 // A.1.2 ⊢ `_` never introduces a binding

	NotGeneric         Code = 3020 // A.3.6 — brackets on a non-generic name
	WrongTypeArgCount  Code = 3021
	MethodTypeParams   Code = 3022 // A.7.6 ⊢ a method may not declare its own
	DuplicateTypeParam Code = 3023 // A.7.1 ⊢ names must be unique within a list

	DuplicateField    Code = 3030
	DuplicateVariant  Code = 3031
	PayloadDiscrim    Code = 3032 // A.6.5 ⊢ explicit discriminant on a unit variant only
	DiscrimNoType     Code = 3033 // A.6.5 ⊢ ...and only with a declared DiscriminantType
	InitializerResult Code = 3034 // A.8.3 ⊢ an initializer must return its enclosing type

	VariadicNotLast Code = 3040 // A.6.1 ⊢ a variadic parameter must be last
	MultipleVariadic Code = 3041 // A.6.1 ⊢ ...and there may be at most one
	IterationBinding Code = 3042 // A.5.6 — `mut a, b` has no production
)

// ------------------------------------------------------------------- types

// Types (A.3). Shape rules checked over an already-resolved type.
const (
	StackedOwnership   Code = 4001 // A.3.2 ⊢ qualifiers do not stack
	MutOutsidePosition Code = 4002 // A.3.2 ⊢ mut/var only in a parameter or receiver
	TildeOutsideSet    Code = 4003 // A.7.3 ✗ type X = ~int
	AbstractInline     Code = 4004 // A.3.3 ⊢ abstract only as an alias target
	ArrayLenNotConst   Code = 4005 // A.3.1 ⊢ ArrayLength must be compile-time
	ArrayLenNegative   Code = 4006
	NonComparableKey   Code = 4007 // A.3.1 ⊢ map[K]V requires K comparable
	TensorOutsideNpu   Code = 4008 // A.3.5 ⊢ grammatical only under [+Npu]
	NotRepresentable   Code = 4009 // A.1.5.1 ⊢ never a silent truncation
	InfiniteType       Code = 4010 // a type that directly contains itself
)

// ---------------------------------------------------------------- ownership

// Ownership and exclusivity (A.9). Listed now only because A.1.4 commits to
// this fix-it by name, and the Fixit path needed a real user before the scanner
// landed. The rest of A.14's ownership section belongs to the analyzer.
const (
	TransferMethodRemoved Code = 5001 // A.1.4 ✗ x.transfer()
)

// ------------------------------------------------------- pointers and memory

// Pointers and memory (A.4.8). Same rationale as above: A.4.8 commits to the
// addr → & fix-it by name.
const (
	AddrOnNonPointer Code = 7001 // A.4.8 ⊢ addr accepts a typed_ptr operand only
	AddrOnTemporary  Code = 7002 // A.4.8 ⊢ ...and only an addressable one
)

// ---------------------------------------------------------------- interop

// Interop (A.8). Only the forms the resolve pass can already see: A.8.3's ✗
// members parse (ast.ForeignFunc keeps Modifiers and Body for exactly this)
// and are rejected the moment a declare block is collected.
const (
	ForeignModifier Code = 8001 // A.8.3 ✗ visibility modifiers are banned
	ForeignBody     Code = 8002 // A.8.3 ✗ declarations cannot have bodies
	ForeignField    Code = 8003 // A.8.3 ✗ fields describe foreign-side layout
	FrameworkTag    Code = 8004 // A.8.2 ✗ framework blocks take no variant tag
	DeclareNoTag    Code = 8005 // A.8.1 ⊢ legal only in a file with a BuildClause
	NoFrameworks    Code = 8006 // A.8.1 ⊢ ...and only where the platform has them
)

// ------------------------------------------------------------------ registry

type info struct {
	tmpl string
	sev  Severity
}

// registry maps each Code to its message template and severity.
//
// The template is the normative text. Changing one changes what
// Expected(error, "...") tests match, and is therefore a language change rather
// than an implementation detail — A.12.2 is explicit that this is the point.
//
// Templates are written to read the same whether or not a span is shown, since
// Diagnostic.Text() strips every position. That rules out deictic phrasing
// ("this literal", "here") in favour of naming the construct.
var registry = map[Code]info{
	Internal: {"internal compiler error: %s", Error},

	// lexical
	IllegalCharacter: {"illegal character %s in source", Error},
	IllegalUTF8:      {"invalid UTF-8 encoding", Error},
	DollarInIdent:    {"'$' is not an identifier character in Vertex", Error},

	UnterminatedComment: {"comment is not terminated before end of file", Error},

	UnterminatedString:    {"string literal is not terminated", Error},
	NewlineInString:       {"string literal is not terminated before end of line", Error},
	UnterminatedRawString: {"raw string literal is not terminated", Error},
	UnterminatedChar:      {"character literal is not terminated", Error},
	EmptyChar:             {"character literal is empty", Error},
	CharTooLong:           {"character literal holds more than one Unicode scalar value", Error},
	InvalidEscape:         {"unknown escape sequence '\\%c'", Error},
	InvalidHexEscape:      {"'\\x' requires exactly two hexadecimal digits", Error},
	InvalidUnicodeEscape:  {"'\\u' requires a braced hexadecimal scalar, as in \\u{1F600}", Error},
	UnicodeEscapeRange:    {"%s is not a Unicode scalar value", Error},

	SeparatorLeads:       {"'_' may not lead a digit run", Error},
	SeparatorTrails:      {"'_' may not trail a digit run", Error},
	SeparatorDoubled:     {"'_' may not be doubled in a digit run", Error},
	EmptyDigits:          {"%s requires at least one digit", Error},
	DigitOutOfRange:      {"digit '%c' is out of range for a %s literal", Error},
	HexFloatNoExponent:   {"hexadecimal float requires a binary exponent, as in 0xC.3p0", Error},
	MissingExponentDigit: {"exponent requires at least one digit", Error},
	NumberJoinedToIdent:  {"identifier characters may not directly follow a numeric literal", Error},

	// syntactic
	ExpectedToken:   {"expected %s, found %s", Error},
	ExpectedStmtEnd: {"expected end of statement, found %s", Error},
	ExpectedExpr:    {"expected an expression, found %s", Error},
	ExpectedType:    {"expected a type, found %s", Error},
	ExpectedIdent:   {"expected an identifier, found %s", Error},
	ExpectedDecl:    {"expected a declaration, found %s", Error},
	ExpectedCase:    {"expected 'case' or 'default', found %s", Error},
	ExpectedCall:    {"the '%s' prefix must be applied to a call", Error},

	RangeNotAssociative: {"'..' is not associative; parenthesize one of the ranges", Error},
	LiteralInHeader:     {"a composite literal may not appear unparenthesized in a %s header", Error},
	TupleNeedsComma:     {"a one-element tuple literal requires its trailing comma", Error},

	MissingPackageClause: {"file must begin with a package clause", Error},
	MisplacedBuildClause: {"a build clause must immediately follow the package clause", Error},
	UnknownBuildTag:      {"unknown build tag %q", Error},
	ImportAfterDecl:      {"imports must precede all declarations", Error},
	DeferNotCall:         {"'defer' takes a call and nothing else", Error},

	// declarations and names
	UndeclaredName:       {"undeclared name %s", Error},
	DuplicateDeclaration: {"%s is already declared in this scope", Error},
	ShadowedBuiltin:      {"%s is a reserved builtin name and may not be redeclared", Error},
	ReservedAsMember:     {"%s is a reserved builtin name, not a declarable member", Error},
	NotAType:             {"%s is not a type", Error},
	ConstraintAsType:     {"%s is a constraint, not a type; a constraint is legal only in a '[...]' position", Error},
	NotAConstraint:       {"%s does not satisfy the constraint", Error},
	TypeCycle:            {"type %s refers to itself through its own definition", Error},
	BlankNotUsable:       {"'_' introduces no binding and cannot be referred to", Error},

	NotGeneric:         {"%s is not generic and takes no type arguments", Error},
	WrongTypeArgCount:  {"%s takes %d type argument(s), but %d were given", Error},
	MethodTypeParams:   {"a method may not declare its own type parameters; %s is generic over its receiver", Error},
	DuplicateTypeParam: {"type parameter %s is already declared in this list", Error},

	DuplicateField:    {"field %s is already declared in %s", Error},
	DuplicateVariant:  {"variant %s is already declared in %s", Error},
	PayloadDiscrim:    {"an explicit discriminant is legal only on a unit variant", Error},
	DiscrimNoType:     {"an explicit discriminant requires a declared discriminant type", Error},
	InitializerResult: {"an initializer must return its enclosing type %s", Error},

	VariadicNotLast:  {"a variadic parameter must be last in its parameter list", Error},
	MultipleVariadic: {"a function may declare at most one variadic parameter", Error},
	IterationBinding: {"'%s' does not combine with a two-name iteration binding", Error},

	// types
	StackedOwnership:   {"ownership qualifiers do not stack: %s", Error},
	MutOutsidePosition: {"'%s' is legal only in a parameter or receiver position", Error},
	TildeOutsideSet:    {"'~' is legal only inside a type set", Error},
	AbstractInline:     {"'abstract' is legal only as the target of a type alias", Error},
	ArrayLenNotConst:   {"array length must be a compile-time constant", Error},
	ArrayLenNegative:   {"array length must not be negative", Error},
	NonComparableKey:   {"%s is not comparable and cannot be a map key", Error},
	TensorOutsideNpu:   {"'tensor' is legal only inside an npu-marked function", Error},
	NotRepresentable:   {"%s is not representable by type %s", Error},
	InfiniteType:       {"type %s has no finite size; it contains itself directly", Error},

	// ownership
	TransferMethodRemoved: {"'.transfer()' is not a method; ownership transfer is spelled with the 'var' prefix", Error},

	// pointers and memory
	AddrOnNonPointer: {"'addr' accepts a typed_ptr operand only; '&' is already the address of %s", Error},
	AddrOnTemporary:  {"'addr' requires an addressable operand: a var binding or a field path", Error},

	// interop
	ForeignModifier: {"a visibility modifier is not permitted on a foreign declaration", Error},
	ForeignBody:     {"a foreign declaration cannot have a body", Error},
	ForeignField:    {"a declare block describes call shape only; %s declares foreign-side layout", Error},
	FrameworkTag:    {"a framework block takes no variant tag", Error},
	DeclareNoTag:    {"a declare block requires the file to carry a build clause", Error},
	NoFrameworks:    {"build tag %q has no notion of a bundled framework", Error},
}

// declaredCodes lists every Code constant in this file.
//
// It is maintained by hand and exists so a conformance test can assert three
// properties the compiler cannot: that every declared code has a registry
// entry, that no two codes share a numeric value, and that no two share a
// template. Without the list there is nothing to enumerate — Go has no
// reflection over untyped constants — and a code declared but never registered
// would surface as a silent Internal at the first call site that used it.
var declaredCodes = []Code{
	Internal,

	IllegalCharacter, IllegalUTF8, DollarInIdent,
	UnterminatedComment,
	UnterminatedString, NewlineInString, UnterminatedRawString,
	UnterminatedChar, EmptyChar, CharTooLong,
	InvalidEscape, InvalidHexEscape, InvalidUnicodeEscape, UnicodeEscapeRange,
	SeparatorLeads, SeparatorTrails, SeparatorDoubled,
	EmptyDigits, DigitOutOfRange, HexFloatNoExponent,
	MissingExponentDigit, NumberJoinedToIdent,

	ExpectedToken, ExpectedStmtEnd, ExpectedExpr, ExpectedType,
	ExpectedIdent, ExpectedDecl, ExpectedCase, ExpectedCall,
	RangeNotAssociative, LiteralInHeader, TupleNeedsComma,
	MissingPackageClause, MisplacedBuildClause, UnknownBuildTag,
	ImportAfterDecl, DeferNotCall,

	UndeclaredName, DuplicateDeclaration, ShadowedBuiltin, ReservedAsMember,
	NotAType, ConstraintAsType, NotAConstraint, TypeCycle, BlankNotUsable,
	NotGeneric, WrongTypeArgCount, MethodTypeParams, DuplicateTypeParam,
	DuplicateField, DuplicateVariant, PayloadDiscrim, DiscrimNoType,
	InitializerResult,
	VariadicNotLast, MultipleVariadic, IterationBinding,

	StackedOwnership, MutOutsidePosition, TildeOutsideSet, AbstractInline,
	ArrayLenNotConst, ArrayLenNegative, NonComparableKey, TensorOutsideNpu,
	NotRepresentable, InfiniteType,

	TransferMethodRemoved,

	AddrOnNonPointer, AddrOnTemporary,

	ForeignModifier, ForeignBody, ForeignField, FrameworkTag,
	DeclareNoTag, NoFrameworks,
}

// String renders the stable identifier, e.g. "V1003". This is what a test that
// cares about which rule fired should pin, rather than the message text.
func (c Code) String() string { return fmt.Sprintf("V%04d", int(c)) }

// Severity reports the code's registered severity, defaulting to Error for an
// unregistered code so that a registry gap cannot downgrade a rejection into a
// warning.
func (c Code) Severity() Severity {
	if m, ok := registry[c]; ok {
		return m.sev
	}
	return Error
}

// Template exposes the raw format string. Intended for the test runner and for
// tooling that enumerates the normative surface, not for call sites — a call
// site that formats its own text is exactly what the registry exists to prevent.
func (c Code) Template() string { return registry[c].tmpl }

// Registered reports whether c has a registry entry.
func (c Code) Registered() bool {
	_, ok := registry[c]
	return ok
}

// Codes returns every declared code in declaration order.
func Codes() []Code {
	out := make([]Code, len(declaredCodes))
	copy(out, declaredCodes)
	return out
}