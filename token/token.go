package token

// Token is one lexeme.
//
// There is no SEMICOLON kind and no synthesized NEWLINE kind. Vertex has no
// statement terminator (A.0.6), so line-structure reaches the parser as a flag
// on the following token rather than as a token of its own.
//
// The reason a NEWLINE token does not work: A.0.6 makes a line break inside
// { } ordinary whitespace, but statements also live inside { } blocks and end
// at a LineTerminator. Whether a given { opens a Block (newline-significant) or
// a CompositeLiteral / MapLiteral / FieldList / DeclareBody (not) is exactly
// the [+Lit] question the parser tracks — the scanner cannot know it, so the
// scanner must never suppress. It records; the parser decides.
type Token struct {
	Kind Kind
	Pos  Pos

	// Lit is the raw source text for IDENT, INT, FLOAT, CHAR, STRING, COMMENT.
	// Raw, not unescaped — vfmt needs the original spelling. Empty otherwise.
	Lit string

	// NLBefore is true if a LineTerminator separates this token from the
	// previous one. A MultiLineComment spanning a LineTerminator sets it too,
	// per A.1.1's rule that such a comment is itself a LineTerminator.
	//
	// This is also what a [no LineTerminator here] restriction reads.
	NLBefore bool
}

func (t Token) Is(k Kind) bool { return t.Kind == k }

func (t Token) End() Pos { return t.Pos + Pos(len(t.Lit)) }

// ContextualKeyword spellings (A.1.3). These scan as IDENT; the parser
// string-compares at the single production that names each.
const (
	CtxBuild     = "build"     // second line-initial token of a file (A.2.2)
	CtxDeinit    = "deinit"    // method name, or declare-block modifier (A.8.3)
	CtxError     = "error"     // inside Expected only (A.12.2)
	CtxFramework = "framework" // immediately after declare (A.8.1)
	CtxInit      = "init"      // method name, or declare-block modifier (A.8.3)
	CtxModule    = "module"    // immediately after declare (A.8.1)
	CtxTest      = "test"      // function-marker position (A.6.1)
)

// IsCtx reports whether t is the identifier `name`. Use with the Ctx*
// constants: tok.IsCtx(token.CtxFramework).
func (t Token) IsCtx(name string) bool {
	return t.Kind == IDENT && t.Lit == name
}

// IsBlank reports whether t is the BlankIdentifier `_` (A.1.2).
func (t Token) IsBlank() bool { return t.Kind == IDENT && t.Lit == "_" }