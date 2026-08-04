package diag

import "fmt"

// Code identifies a diagnostic rule. It is the stable handle: renumbering one
// is a breaking change to the language's test corpus, because the error form of
// grammar.md's ExpectedType — `Expected(error, "...")` — compares rendered
// message text as specification.
//
// Numeric values are assigned explicitly, never by iota, so inserting a code
// cannot shift another. Ranges follow grammar.md's own section order:
//
//	0xxx  internal
//	1xxx  lexical elements
//	2xxx  syntax errors
//	3xxx  types
//	4xxx  expressions
//	5xxx  statements
//	6xxx  declarations and names
//	7xxx  generics
//	8xxx  declare blocks
//	9xxx  test result types
//
// Gaps within a range are deliberate and reserved. Only 1xxx and 2xxx are
// densely populated today, because the scanner and the parser are the only
// phases that exist; every other range holds the forms grammar.md names as
// parsing-and-then-rejected, waiting for the analyzer pass that raises them.
//
// grammar.md defines syntax only. Where it says "static rule," it means the
// form derives and is checked after parsing — those are the codes outside 1xxx
// and 2xxx. Rules that grammar.md does not mention at all belong to
// semantics.md and are not represented here.
type Code int

// ------------------------------------------------------------------ internal

const (
	// Internal marks a broken compiler invariant. Never produced by valid or
	// invalid source; if a user sees one, it is a bug in the compiler.
	Internal Code = 1
)

// ----------------------------------------------------------------- lexical

// Lexical elements. The scanner is the only producer.
const (
	IllegalCharacter Code = 1001
	IllegalUTF8      Code = 1002
	NulCharacter     Code = 1003 // ⊢ a compiler must reject U+0000 in source text
	DollarInIdent    Code = 1004 // ⊢ `$` is not an identifier character in any position

	UnterminatedComment Code = 1010 // ⊢ general comments do not nest; the first */ closes

	UnterminatedString    Code = 1020
	NewlineInString       Code = 1021 // ⊢ a line terminator may not appear in an interpreted string
	UnterminatedRawString Code = 1022
	UnterminatedChar      Code = 1023
	EmptyChar             Code = 1024
	CharTooLong           Code = 1025 // ⊢ exactly one Unicode scalar value
	InvalidEscape         Code = 1026
	InvalidHexEscape      Code = 1027 // \xHH
	InvalidUnicodeEscape  Code = 1028 // \u{HexDigits}
	UnicodeEscapeRange    Code = 1029 // > 10FFFF, or a surrogate

	SeparatorLeads       Code = 1040 // ⊢ `_` may not lead a digit run
	SeparatorTrails      Code = 1041 // ...nor trail one
	SeparatorDoubled     Code = 1042 // ...nor be doubled
	EmptyDigits          Code = 1043 // `0x` with no digits
	DigitOutOfRange      Code = 1044 // `0b2`, `0o9`
	HexFloatNoExponent   Code = 1045 // ⊢ a hexadecimal float requires its binary exponent
	MissingExponentDigit Code = 1046 // `1e`, `0x1p`
	NumberJoinedToIdent  Code = 1047 // `123abc`
)

// ------------------------------------------------------------------ syntax

// Syntax errors. The parser is the only producer, and every code here is
// decidable from the tree the parser has in hand — no resolution required.
//
// The range is deliberately small. Most of what grammar.md forbids is written
// so that it parses and is rejected afterwards, precisely so the diagnostic can
// name the construct instead of pointing at a token.
const (
	ExpectedToken      Code = 2001
	ExpectedTerminator Code = 2002
	ExpectedExpr       Code = 2003
	ExpectedType       Code = 2004
	ExpectedIdent      Code = 2005
	ExpectedDecl       Code = 2006
	ExpectedCase       Code = 2007
	ExpectedCall       Code = 2008 // ⊢ CallExpr: "a call and nothing else"

	RangeNotAssociative Code = 2010 // ⊢ `a..b..c` is a compile error
	LiteralInHeader     Code = 2011 // ⊢ parenthesize a literal used in a header
	EmptyTuple          Code = 2012 // ⊢ a tuple has at least one element; there is no unit type
	DuplicateDefault    Code = 2013 // ⊢ at most one `default` clause

	MissingPackageClause Code = 2020 // ⊢ mandatory, first non-comment construct
	MisplacedBuildClause Code = 2021 // ⊢ the second construct in the file
	UnknownBuildTag      Code = 2022 // ⊢ an error, not a silent exclusion
	ImportAfterDecl      Code = 2023 // ⊢ all imports precede all declarations
)

// ------------------------------------------------------------------- types

// Types, and the signatures that are types. Shape rules checked over an
// already-resolved type.
const (
	NotAType         Code = 3001
	TypeCycle        Code = 3002
	InfiniteType     Code = 3003
	NotRepresentable Code = 3004

	StackedOwnership   Code = 3010 // ⊢ ownership qualifiers do not stack
	MutOutsidePosition Code = 3011 // ⊢ mut/var only in a parameter or receiver position
	NestedPointer      Code = 3012 // ⊢ a typed_ptr may not be the direct base of another
	AbstractInline     Code = 3013 // ⊢ abstract only as an alias target

	ArrayLenNotConst Code = 3020
	ArrayLenNegative Code = 3021
	NonComparableKey Code = 3022

	TensorOutsideNpu    Code = 3030 // ⊢ tensor only inside an npu-marked function
	TensorElemOutsideNpu Code = 3031 // ⊢ ...and its element type names are body-only
	TensorElemInSignature Code = 3032 // ⊢ ...never in a signature

	MixedParamNames  Code = 3040 // ⊢ names all present or all absent
	VariadicNotLast  Code = 3041 // ⊢ a variadic parameter must be last
	MultipleVariadic Code = 3042 // ⊢ ...and there may be at most one
	MultipleMarkers  Code = 3043 // ⊢ a signature carries at most one FunctionMarker
)

// ------------------------------------------------------------- expressions

// Expressions, including the ownership marker and the owning positions it may
// occupy.
const (
	TransferNotCallable  Code = 4001 // ⊢ `transfer` is bound to nothing, on purpose
	TransferComputed     Code = 4002 // ⊢ the operand must be a binding or a field path
	TransferOutsideOwning Code = 4003 // ⊢ six owning positions, and no others

	MixedArguments Code = 4010 // ⊢ arguments are positional or named, not both

	AwaitOutsideAsync   Code = 4020 // ⊢ parses unconditionally; licensed by the body
	TupleIndexSpelling  Code = 4021 // ⊢ a decimal_lit containing no `_`
	EnumShorthandNoType Code = 4022 // ⊢ legal only where the enum type is fixed by context
)

// -------------------------------------------------------------- statements

// Statements.
const (
	NotAssignable Code = 5001 // ⊢ which PrimaryExpr shapes are assignable is a static rule

	IterationBindingMode Code = 5010 // ⊢ `mut a, b` does not combine

	SelectCaseNotChannelOp Code = 5020 // ⊢ which calls are admissible in a ChannelCase
	SelectMixedAwait       Code = 5021 // ⊢ entirely bare or entirely awaited
)

// ------------------------------------------- declarations and names

// Declarations and names. Every code here answers "what does this name
// denote", which is the question resolution exists to answer and the one
// several of grammar.md's deferred ambiguities reduce to: an Index whose
// operand denotes a generic declaration, a single identifier that is both a
// one-term TypeSet and a constraint name.
const (
	UndeclaredName       Code = 6001
	DuplicateDeclaration Code = 6002
	ShadowedBuiltin      Code = 6003 // ⊢ a reserved builtin name may not be shadowed
	ReservedAsMember     Code = 6004 // ⊢ ...nor declared as a member, method, field, local, or parameter

	BlankNotUsable  Code = 6010 // ⊢ `_` introduces no binding
	BlankNotAllowed Code = 6011 // ⊢ which positions accept it is a static rule

	DuplicateField   Code = 6020
	DuplicateVariant Code = 6021
	PayloadDiscrim   Code = 6022 // ⊢ explicit discriminant on a unit variant only
	DiscrimNoType    Code = 6023 // ⊢ ...and only with a declared DiscriminantType

	TopLevelVarNotConst Code = 6030 // ⊢ a top-level initializer is compile-time-evaluable
	TopLevelBareVar     Code = 6031 // ⊢ ...and the bare form is rejected there
)

// ---------------------------------------------------------------- generics

// Generics and constraints.
const (
	NotGeneric        Code = 7001 // ⊢ brackets on a non-generic name
	WrongTypeArgCount Code = 7002
	MethodTypeParams  Code = 7003 // ⊢ a method may not declare its own
	DuplicateTypeParam Code = 7004

	TildeOutsideSet      Code = 7010 // ⊢ `~` outside a TypeSet is an error
	ConstraintAsType     Code = 7011 // ⊢ legal only in a `[` ... `]` position
	NotAConstraint       Code = 7012
	ConstraintNotSatisfied Code = 7013
)

// ---------------------------------------------------------- declare blocks

// Declare blocks. A declare block describes call shapes only; grammar.md lists
// four forms that parse there and are rejected, so the diagnostic can name the
// construct. Each has a code.
const (
	ForeignBody    Code = 8001 // ⊢ a Block on a ForeignFuncDecl
	ForeignField   Code = 8002 // ⊢ a FieldDecl in a ForeignClassDecl
	NestedDeclare  Code = 8003 // ⊢ a nested DeclareDecl
	ForeignMarker  Code = 8004 // ⊢ a FunctionMarker on a ForeignFuncDecl

	FrameworkTag     Code = 8010 // ⊢ a framework block takes no VariantTag
	UnknownVariantTag Code = 8011 // ⊢ the tag set is closed

	ForeignInitResult Code = 8020 // ⊢ an initializer returns its enclosing type
)

// ------------------------------------------------------ test result types

// Test result types. An ExpectedType reaches the grammar only through
// DeclResult, which is syntax; that it is further restricted to a file built
// under the `test` tag is the static rule this range holds.
const (
	ExpectedOutsideTest Code = 9001
)

// ------------------------------------------------------------------ registry

type info struct {
	tmpl string
	sev  Severity
}

// registry maps each Code to its message template and severity.
//
// The template is the normative text. Changing one changes what the error form
// of ExpectedType matches, and is therefore a language change rather than an
// implementation detail.
//
// Templates are written to read the same whether or not a span is shown, since
// Diagnostic.Text() strips every position. That rules out deictic phrasing
// ("this literal", "here") in favour of naming the construct.
var registry = map[Code]info{
	Internal: {"internal compiler error: %s", Error},

	// lexical
	IllegalCharacter: {"illegal character %s in source", Error},
	IllegalUTF8:      {"source text is not valid UTF-8", Error},
	NulCharacter:     {"the NUL character is not permitted in source text", Error},
	DollarInIdent:    {"'$' is not an identifier character in Vertex", Error},

	UnterminatedComment: {"general comment is not terminated before end of file", Error},

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

	// syntax
	ExpectedToken:      {"expected %s, found %s", Error},
	ExpectedTerminator: {"expected a line terminator, found %s", Error},
	ExpectedExpr:       {"expected an expression, found %s", Error},
	ExpectedType:       {"expected a type, found %s", Error},
	ExpectedIdent:      {"expected an identifier, found %s", Error},
	ExpectedDecl:       {"expected a declaration, found %s", Error},
	ExpectedCase:       {"expected 'case' or 'default', found %s", Error},
	ExpectedCall:       {"'%s' takes a call and nothing else", Error},

	RangeNotAssociative: {"'..' is not associative; parenthesize one of the ranges", Error},
	LiteralInHeader:     {"a composite literal may not appear unparenthesized in a '%s' header", Error},
	EmptyTuple:          {"a tuple has at least one element; Vertex has no unit type", Error},
	DuplicateDefault:    {"a '%s' statement may have at most one 'default' clause", Error},

	MissingPackageClause: {"file must begin with a package clause", Error},
	MisplacedBuildClause: {"a build clause must immediately follow the package clause", Error},
	UnknownBuildTag:      {"unknown build tag %q", Error},
	ImportAfterDecl:      {"imports must precede all declarations", Error},

	// types
	NotAType:         {"%s is not a type", Error},
	TypeCycle:        {"type %s refers to itself through its own definition", Error},
	InfiniteType:     {"type %s has no finite size; it contains itself directly", Error},
	NotRepresentable: {"%s is not representable by type %s", Error},

	StackedOwnership:   {"ownership qualifiers do not stack: %s", Error},
	MutOutsidePosition: {"'%s' is legal only in a parameter or receiver position", Error},
	NestedPointer:      {"'typed_ptr' may not be the direct base of another; parenthesize the inner type", Error},
	AbstractInline:     {"'abstract' is legal only as the target of a type alias", Error},

	ArrayLenNotConst: {"array length must be a compile-time constant", Error},
	ArrayLenNegative: {"array length must not be negative", Error},
	NonComparableKey: {"%s is not comparable and cannot be a map key", Error},

	TensorOutsideNpu:      {"'tensor' is legal only inside an npu-marked function", Error},
	TensorElemOutsideNpu:  {"%s is a tensor element type and is legal only inside an npu-marked function body", Error},
	TensorElemInSignature: {"%s is a tensor element type and may not appear in a signature", Error},

	MixedParamNames:  {"parameter names must be either all present or all absent", Error},
	VariadicNotLast:  {"a variadic parameter must be last in its parameter list", Error},
	MultipleVariadic: {"a function may declare at most one variadic parameter", Error},
	MultipleMarkers:  {"a signature carries at most one function marker", Error},

	// expressions
	TransferNotCallable:   {"'transfer' is not callable; ownership transfer is spelled with the 'var' prefix", Error},
	TransferComputed:      {"'var' requires a binding or a field path, not a computed expression", Error},
	TransferOutsideOwning: {"'var' is the transfer marker and is legal only in an owning position", Error},

	MixedArguments: {"a call takes either named or positional arguments, not both", Error},

	AwaitOutsideAsync:   {"'await' is legal only inside an async-marked function body", Error},
	TupleIndexSpelling:  {"a tuple index must be written in decimal without '_'", Error},
	EnumShorthandNoType: {"the enum type of a leading-dot shorthand is not fixed by context", Error},

	// statements
	NotAssignable: {"%s is not assignable", Error},

	IterationBindingMode: {"'%s' does not combine with a two-name iteration binding", Error},

	SelectCaseNotChannelOp: {"a 'select' case requires a channel operation", Error},
	SelectMixedAwait:       {"a 'select' statement is either entirely bare or entirely awaited", Error},

	// declarations and names
	UndeclaredName:       {"undeclared name %s", Error},
	DuplicateDeclaration: {"%s is already declared in this scope", Error},
	ShadowedBuiltin:      {"%s is a reserved builtin name and may not be redeclared", Error},
	ReservedAsMember:     {"%s is a reserved builtin name, not a declarable member", Error},

	BlankNotUsable:  {"'_' introduces no binding and cannot be referred to", Error},
	BlankNotAllowed: {"the blank identifier is not permitted as %s", Error},

	DuplicateField:   {"field %s is already declared in %s", Error},
	DuplicateVariant: {"variant %s is already declared in %s", Error},
	PayloadDiscrim:   {"an explicit discriminant is legal only on a unit variant", Error},
	DiscrimNoType:    {"an explicit discriminant requires a declared discriminant type", Error},

	TopLevelVarNotConst: {"a top-level variable's initializer must be compile-time-evaluable", Error},
	TopLevelBareVar:     {"the bare 'var' declaration form is not permitted at the top level", Error},

	// generics
	NotGeneric:         {"%s is not generic and takes no type arguments", Error},
	WrongTypeArgCount:  {"%s takes %d type argument(s), but %d were given", Error},
	MethodTypeParams:   {"a method may not declare its own type parameters; %s is generic over its receiver", Error},
	DuplicateTypeParam: {"type parameter %s is already declared in this list", Error},

	TildeOutsideSet:        {"'~' is legal only inside a type set", Error},
	ConstraintAsType:       {"%s is a constraint, not a type; a constraint is legal only in a '[' ']' position", Error},
	NotAConstraint:         {"%s is not a constraint", Error},
	ConstraintNotSatisfied: {"%s does not satisfy the constraint %s", Error},

	// declare blocks
	ForeignBody:   {"a foreign declaration cannot have a body", Error},
	ForeignField:  {"a declare block describes call shape only; %s declares foreign-side layout", Error},
	NestedDeclare: {"a declare block may not contain another declare block", Error},
	ForeignMarker: {"a function marker is not permitted on a foreign declaration", Error},

	FrameworkTag:      {"a framework block takes no variant tag", Error},
	UnknownVariantTag: {"unknown variant tag %q", Error},

	ForeignInitResult: {"a foreign initializer must return its enclosing type %s", Error},

	// test result types
	ExpectedOutsideTest: {"an Expected result requires the file to carry a 'build test' clause", Error},
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

	IllegalCharacter, IllegalUTF8, NulCharacter, DollarInIdent,
	UnterminatedComment,
	UnterminatedString, NewlineInString, UnterminatedRawString,
	UnterminatedChar, EmptyChar, CharTooLong,
	InvalidEscape, InvalidHexEscape, InvalidUnicodeEscape, UnicodeEscapeRange,
	SeparatorLeads, SeparatorTrails, SeparatorDoubled,
	EmptyDigits, DigitOutOfRange, HexFloatNoExponent,
	MissingExponentDigit, NumberJoinedToIdent,

	ExpectedToken, ExpectedTerminator, ExpectedExpr, ExpectedType,
	ExpectedIdent, ExpectedDecl, ExpectedCase, ExpectedCall,
	RangeNotAssociative, LiteralInHeader, EmptyTuple, DuplicateDefault,
	MissingPackageClause, MisplacedBuildClause, UnknownBuildTag,
	ImportAfterDecl,

	NotAType, TypeCycle, InfiniteType, NotRepresentable,
	StackedOwnership, MutOutsidePosition, NestedPointer, AbstractInline,
	ArrayLenNotConst, ArrayLenNegative, NonComparableKey,
	TensorOutsideNpu, TensorElemOutsideNpu, TensorElemInSignature,
	MixedParamNames, VariadicNotLast, MultipleVariadic, MultipleMarkers,

	TransferNotCallable, TransferComputed, TransferOutsideOwning,
	MixedArguments,
	AwaitOutsideAsync, TupleIndexSpelling, EnumShorthandNoType,

	NotAssignable,
	IterationBindingMode,
	SelectCaseNotChannelOp, SelectMixedAwait,

	UndeclaredName, DuplicateDeclaration, ShadowedBuiltin, ReservedAsMember,
	BlankNotUsable, BlankNotAllowed,
	DuplicateField, DuplicateVariant, PayloadDiscrim, DiscrimNoType,
	TopLevelVarNotConst, TopLevelBareVar,

	NotGeneric, WrongTypeArgCount, MethodTypeParams, DuplicateTypeParam,
	TildeOutsideSet, ConstraintAsType, NotAConstraint, ConstraintNotSatisfied,

	ForeignBody, ForeignField, NestedDeclare, ForeignMarker,
	FrameworkTag, UnknownVariantTag, ForeignInitResult,

	ExpectedOutsideTest,
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
// site that formats its own text is exactly what the registry exists to
// prevent.
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