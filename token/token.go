package token

// Token is a lexeme: a classification and a span. Sixteen bytes' worth of
// fields declared, twelve after Go's layout — no pointers, so a []Token is
// invisible to the GC however large (§3). That matters because the parser
// buffers whole files (§4.1).
//
// Literal text is not stored. Recover it with File.Slice; the scanner
// allocates no strings.
type Token struct {
	Kind  Kind  // classification
	Ctx   Ctx   // contextual-keyword id; 0 unless Kind == IDENT
	Flags Flags // NLBefore, HasEscape, Unterminated
	Pos   Pos   // first byte
	End   Pos   // one past the last byte
}

// Flags records lexical facts that have no home in Kind.
type Flags uint8

const (
	// NLBefore marks that at least one LineTerminator appeared between the
	// previous token and this one. A newline is never a token of its own
	// (§4.3): whether it ends a statement is a parsing question, answered by
	// expectSemi (§6.2) and the [no LineTerminator here] restrictions in
	// grammar L.
	//
	// A comment spanning a newline sets this on the next real token, whether
	// or not comments are being retained.
	NLBefore Flags = 1 << iota

	// HasEscape marks an IdentifierName containing a UnicodeEscapeSequence, a
	// StringLiteral or Template containing any EscapeSequence, or a
	// NumericLiteral containing a NumericLiteralSeparator. Two consumers:
	// keyword lookup must be skipped for escaped identifiers, and a later
	// phase knows the raw spelling needs decoding rather than a direct compare.
	HasEscape

	// Unterminated marks a string, template, comment, or regex that ran to
	// end of input or end of line. The token still has an exact span, so
	// recovery (§6.4) can proceed and every node keeps a non-zero span (§1).
	Unterminated
)

// Has reports whether all of mask is set.
func (f Flags) Has(mask Flags) bool { return f&mask == mask }

func (t Token) NLBefore() bool     { return t.Flags&NLBefore != 0 }
func (t Token) HasEscape() bool    { return t.Flags&HasEscape != 0 }
func (t Token) Unterminated() bool { return t.Flags&Unterminated != 0 }

// Len is the token's width in bytes. Zero only for EOF.
func (t Token) Len() int {
	if t.Pos == NoPos || t.End == NoPos {
		return 0
	}
	return int(t.End - t.Pos)
}

// Adjacent reports whether u begins exactly where t ends, with no intervening
// whitespace or comment. This is the test the expression parser uses when
// joining a run of GT tokens (§4.2), which is why `a > > b` stays an error.
func (t Token) Adjacent(u Token) bool { return t.End == u.Pos && t.End != NoPos }

// IsContextual reports whether t is an identifier spelling the given
// contextual keyword. Whether it *declares* anything is a parse question.
func (t Token) IsContextual(c Ctx) bool { return t.Kind == IDENT && t.Ctx == c }

func (t Token) String() string {
	if t.Ctx != CtxNone {
		return "IDENT(" + t.Ctx.String() + ")"
	}
	return t.Kind.String()
}