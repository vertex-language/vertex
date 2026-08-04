// Package token defines the lexical vocabulary and source-position machinery
// shared by every stage of the Vertex toolchain. It depends on nothing else.
//
// There is no statement terminator token and no automatic semicolon insertion.
// A terminator is a run of line terminators, and a run of any length is one
// terminator, so line structure reaches the parser as a single flag on the
// following token rather than as a token of its own.
//
// A NEWLINE token would not work. Whether a line terminator is a terminator
// depends on the innermost enclosing bracketing construct — significant inside
// a Block, ordinary white space inside a LiteralValue — and whether a given `{`
// opens one or the other is a parsing question. So the scanner emits every line
// terminator, interprets none, and the parser decides.
package token

// Token is one lexeme.
type Token struct {
	Kind Kind
	Pos  Pos

	// Lit is the raw source text for IDENT, INT, FLOAT, CHAR, STRING, COMMENT,
	// and INVALID. Raw, not decoded: escape sequences and digit separators
	// survive exactly as written, because a formatter needs the original
	// spelling and decoding belongs to a later phase. Empty for every kind
	// whose spelling is fixed.
	Lit string

	// NLBefore reports whether one or more line terminators separate this token
	// from the previous one. A general comment containing a line terminator
	// sets it too, since such a comment acts like a line terminator.
	//
	// One bool suffices for the whole terminator rule: a run of terminators is
	// one terminator, so no consumer ever needs the count.
	NLBefore bool
}

func (t Token) Is(k Kind) bool { return t.Kind == k }

// Text returns the token's source text: its Lit where it has one, its Kind's
// fixed spelling otherwise.
func (t Token) Text() string {
	if t.Lit != "" {
		return t.Lit
	}
	return t.Kind.Spelling()
}

// End is one past the token's last character, the convention used throughout
// the tree for node extents.
func (t Token) End() Pos { return t.Pos + Pos(len(t.Text())) }

// IsBlank reports whether t is the blank identifier `_`, which is an ordinary
// identifier token that introduces no binding.
func (t Token) IsBlank() bool { return t.Kind == IDENT && t.Lit == "_" }

// Contextual keyword spellings. Each scans as IDENT — this package mints no
// Kind for any of them — and is recognized only at the production named here.
const (
	CtxBuild     = "build"     // BuildClause
	CtxTest      = "test"      // BuildTag, and FunctionMarker
	CtxInit      = "init"      // MethodDecl by name, and ForeignInitDecl
	CtxDeinit    = "deinit"    // MethodDecl by name
	CtxFramework = "framework" // DeclareDecl
	CtxModule    = "module"    // DeclareDecl
	CtxBlocks    = "blocks"    // LaunchConfig
	CtxThreads   = "threads"   // LaunchConfig
	CtxExpected  = "Expected"  // ExpectedType
	CtxError     = "error"     // ExpectedType
)

// The remaining contextual keywords are the BuildTag spellings; BuildTag.String
// and LookupBuildTag are their one home.

// IsCtx reports whether t is the identifier `name`. Call sites pass a Ctx
// constant rather than comparing Lit directly.
func (t Token) IsCtx(name string) bool {
	return t.Kind == IDENT && t.Lit == name
}

// Type-operator names. These are the only call forms taking a Type in argument
// position, and a parser recognizes them by name. That is sound only because a
// reserved builtin name may not be shadowed or declared as a member, method,
// field, local, or parameter — so the guarantee and the names live together.
const (
	BuiltinSizeof      = "sizeof"
	BuiltinAlignof     = "alignof"
	BuiltinReinterpret = "reinterpret"
)

// IsTypeOperator reports whether name heads a TypeOperatorCall.
func IsTypeOperator(name string) bool {
	switch name {
	case BuiltinSizeof, BuiltinAlignof, BuiltinReinterpret:
		return true
	}
	return false
}

// reservedBuiltins is the full set of names pre-bound in the implicit outermost
// scope that may not be shadowed. They are ordinary identifiers lexically;
// the set lives here so the non-shadowing guarantee has one home.
//
// `transfer` is reserved and bound to nothing, so that x.transfer() diagnoses
// as a misspelled ownership marker rather than as an unknown name.
var reservedBuiltins = map[string]bool{
	"new": true, "delete": true, "resize": true, "copy": true,
	"zero": true, "addr": true, "sizeof": true, "alignof": true,
	"reinterpret": true, "upgrade": true, "drop": true, "panic": true,
	"blend": true, "min": true, "max": true, "clamp": true,
	"transfer": true,
}

// IsReservedBuiltin reports whether name may not be shadowed or declared.
func IsReservedBuiltin(name string) bool { return reservedBuiltins[name] }