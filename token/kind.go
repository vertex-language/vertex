// Package token defines the lexical vocabulary of Vertex: token kinds,
// contextual-keyword identity, source positions, and diagnostics.
//
// It imports nothing from scanner, ast, or parser, and it never will
// (compiler_frontend.md §2). Diagnostic lives here rather than in scanner
// because both scanner and parser emit them (§3).
package token

// Kind classifies a token. It is deliberately uint8: Token must stay
// pointer-free and small (§3).
//
// Kinds are assigned in contiguous ranges bounded by unexported sentinels so
// that IsLiteral, IsOperator, and IsReserved are range comparisons rather than
// switches (§3). Do not reorder these constants without updating the sentinels;
// the ranges are the API.
//
// Only ReservedWord (vertex_grammar.md A.2) gets a Kind. Contextual keywords
// and strict reserved words scan as IDENT and carry identity in Ctx — see
// ctx.go for why the strict words live there too.
type Kind uint8

const (
	INVALID Kind = iota // an unrecognized or malformed token
	EOF                 // end of input; zero-width, positioned at len(src)
	COMMENT             // MultiLineComment, SingleLineComment, or HashbangComment

	literalBeg
	IDENT           // IdentifierName that is not a ReservedWord
	PRIVATE_IDENT   // #foo
	NUMBER          // DecimalLiteral or NonDecimalIntegerLiteral, raw spelling
	BIGINT          // ...n
	STRING          // "..." or '...'
	REGEX           // /body/flags
	TEMPLATE        // `...`     NoSubstitutionTemplate
	TEMPLATE_HEAD   // `...${
	TEMPLATE_MIDDLE // }...${
	TEMPLATE_TAIL   // }...`
	literalEnd

	operatorBeg

	// Grouping and punctuation.
	LBRACE       // {
	RBRACE       // }
	LPAREN       // (
	RPAREN       // )
	LBRACK       // [
	RBRACK       // ]
	SEMI         // ;
	COMMA        // ,
	PERIOD       // .
	ELLIPSIS     // ...
	COLON        // :
	QUESTION     // ?
	QUESTION_DOT // ?.   OptionalChainingPunctuator; A.3 requires lookahead ∉ DecimalDigit
	ARROW        // =>
	AT           // @
	HASH         // #    bare; PrivateIdentifier is one PRIVATE_IDENT token

	// Comparison.
	LT         // 
	GT         // >
	LEQ        // <=
	GEQ        // >=    never emitted by the scanner; see §4.2 and JoinGT
	EQL        // ==
	NEQ        // !=
	STRICT_EQL // ===
	STRICT_NEQ // !==

	// Arithmetic, bitwise, logical.
	NOT      // !
	TILDE    // ~
	ADD      // +
	SUB      // -
	MUL      // *
	QUO      // /
	REM      // %
	EXP      // **
	INC      // ++
	DEC      // --
	SHL      // 
	SHR      // >>    never emitted by the scanner
	USHR     // >>>   never emitted by the scanner
	AND      // &
	OR       // |
	XOR      // ^
	LAND     // &&
	LOR      // ||
	COALESCE // ??

	// Assignment.
	ASSIGN          // =
	ADD_ASSIGN      // +=
	SUB_ASSIGN      // -=
	MUL_ASSIGN      // *=
	QUO_ASSIGN      // /=
	REM_ASSIGN      // %=
	EXP_ASSIGN      // **=
	SHL_ASSIGN      // <<=
	SHR_ASSIGN      // >>=   never emitted by the scanner
	USHR_ASSIGN     // >>>=  never emitted by the scanner
	AND_ASSIGN      // &=
	OR_ASSIGN       // |=
	XOR_ASSIGN      // ^=
	LAND_ASSIGN     // &&=
	LOR_ASSIGN      // ||=
	COALESCE_ASSIGN // ??=

	operatorEnd

	reservedBeg
	AWAIT
	BREAK
	CASE
	CATCH
	CLASS
	CONST
	CONTINUE
	DEBUGGER
	DEFAULT
	DELETE
	DO
	ELSE
	ENUM
	EXPORT
	EXTENDS
	FALSE
	FINALLY
	FOR
	FUNCTION
	IF
	IMPORT
	IN
	INSTANCEOF
	NEW
	NULL
	RETURN
	SUPER
	SWITCH
	THIS
	THROW
	TRUE
	TRY
	TYPEOF
	VAR
	VOID
	WHILE
	WITH
	YIELD
	reservedEnd
)

var kindNames = [...]string{
	INVALID: "INVALID",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:           "IDENT",
	PRIVATE_IDENT:   "PRIVATE_IDENT",
	NUMBER:          "NUMBER",
	BIGINT:          "BIGINT",
	STRING:          "STRING",
	REGEX:           "REGEX",
	TEMPLATE:        "TEMPLATE",
	TEMPLATE_HEAD:   "TEMPLATE_HEAD",
	TEMPLATE_MIDDLE: "TEMPLATE_MIDDLE",
	TEMPLATE_TAIL:   "TEMPLATE_TAIL",

	LBRACE:       "{",
	RBRACE:       "}",
	LPAREN:       "(",
	RPAREN:       ")",
	LBRACK:       "[",
	RBRACK:       "]",
	SEMI:         ";",
	COMMA:        ",",
	PERIOD:       ".",
	ELLIPSIS:     "...",
	COLON:        ":",
	QUESTION:     "?",
	QUESTION_DOT: "?.",
	ARROW:        "=>",
	AT:           "@",
	HASH:         "#",

	LT:         "<",
	GT:         ">",
	LEQ:        "<=",
	GEQ:        ">=",
	EQL:        "==",
	NEQ:        "!=",
	STRICT_EQL: "===",
	STRICT_NEQ: "!==",

	NOT:      "!",
	TILDE:    "~",
	ADD:      "+",
	SUB:      "-",
	MUL:      "*",
	QUO:      "/",
	REM:      "%",
	EXP:      "**",
	INC:      "++",
	DEC:      "--",
	SHL:      "<<",
	SHR:      ">>",
	USHR:     ">>>",
	AND:      "&",
	OR:       "|",
	XOR:      "^",
	LAND:     "&&",
	LOR:      "||",
	COALESCE: "??",

	ASSIGN:          "=",
	ADD_ASSIGN:      "+=",
	SUB_ASSIGN:      "-=",
	MUL_ASSIGN:      "*=",
	QUO_ASSIGN:      "/=",
	REM_ASSIGN:      "%=",
	EXP_ASSIGN:      "**=",
	SHL_ASSIGN:      "<<=",
	SHR_ASSIGN:      ">>=",
	USHR_ASSIGN:     ">>>=",
	AND_ASSIGN:      "&=",
	OR_ASSIGN:       "|=",
	XOR_ASSIGN:      "^=",
	LAND_ASSIGN:     "&&=",
	LOR_ASSIGN:      "||=",
	COALESCE_ASSIGN: "??=",

	AWAIT:      "await",
	BREAK:      "break",
	CASE:       "case",
	CATCH:      "catch",
	CLASS:      "class",
	CONST:      "const",
	CONTINUE:   "continue",
	DEBUGGER:   "debugger",
	DEFAULT:    "default",
	DELETE:     "delete",
	DO:         "do",
	ELSE:       "else",
	ENUM:       "enum",
	EXPORT:     "export",
	EXTENDS:    "extends",
	FALSE:      "false",
	FINALLY:    "finally",
	FOR:        "for",
	FUNCTION:   "function",
	IF:         "if",
	IMPORT:     "import",
	IN:         "in",
	INSTANCEOF: "instanceof",
	NEW:        "new",
	NULL:       "null",
	RETURN:     "return",
	SUPER:      "super",
	SWITCH:     "switch",
	THIS:       "this",
	THROW:      "throw",
	TRUE:       "true",
	TRY:        "try",
	TYPEOF:     "typeof",
	VAR:        "var",
	VOID:       "void",
	WHILE:      "while",
	WITH:       "with",
	YIELD:      "yield",
}

// String returns the operator spelling for operators and reserved words, and
// an uppercase category name otherwise. Golden token dumps (§7) depend on this
// being stable, so treat changes as fixture-breaking.
func (k Kind) String() string {
	if int(k) < len(kindNames) {
		if s := kindNames[k]; s != "" {
			return s
		}
	}
	return "Kind(" + itoa(int(k)) + ")"
}

// IsLiteral reports whether k is a literal-ish token. IDENT is inside this
// range, matching go/token; "literal" here means "carries source text that
// must be read back with File.Slice", not "denotes a value".
func (k Kind) IsLiteral() bool { return literalBeg < k && k < literalEnd }

// IsOperator reports whether k is a Punctuator, DivPunctuator, or
// RightBracePunctuator (A.3).
func (k Kind) IsOperator() bool { return operatorBeg < k && k < operatorEnd }

// IsReserved reports whether k is a ReservedWord (A.2). StrictReservedWord is
// not included: those scan as IDENT with a Ctx (see ctx.go).
func (k Kind) IsReserved() bool { return reservedBeg < k && k < reservedEnd }

// Precedence levels for the binary chain in vertex_grammar.md B.4.
const (
	LowestPrec  = 0  // non-binary
	UnaryPrec   = 12 // above every binary operator
	HighestPrec = 13
)

var precedences = [...]int{
	// CoalesceExpression and LogicalORExpression share a level. B.4 forbids
	// mixing ?? with || / && in one chain, which a precedence table cannot
	// express; the parser records the mix and rejects it later (§6.3).
	COALESCE: 1,
	LOR:      1,
	LAND:     2,

	OR:  3,
	XOR: 4,
	AND: 5,

	EQL: 6, NEQ: 6, STRICT_EQL: 6, STRICT_NEQ: 6,

	// RelationalExpression. `as`, `as const`, and `satisfies` also bind here,
	// but they are IDENT tokens carrying CtxAs / CtxSatisfies and are handled
	// at this level by the parser, not by this table.
	LT: 7, GT: 7, LEQ: 7, GEQ: 7, INSTANCEOF: 7,
	IN: 7, // only when [+In]; the guard is the parser's

	SHL: 8, SHR: 8, USHR: 8,

	ADD: 9, SUB: 9,

	MUL: 10, QUO: 10, REM: 10,

	// Right-associative, and its left operand is UpdateExpression rather than
	// UnaryExpression (B.4). Both are special cases in the climb.
	EXP: 11,
}

// Precedence returns the binary precedence of k, or LowestPrec if k is not a
// binary operator. Exposed here so a formatter can decide which parentheses
// are redundant without keeping a second copy of the table (§3).
//
// The joined forms >=, >>, >>>, >>=, and >>>= have precedences even though the
// scanner never emits them (§4.2). The expression parser joins a run of GT
// tokens with JoinGT and then asks this table for the result. Deleting them as
// dead would break the climb.
func (k Kind) Precedence() int {
	if int(k) < len(precedences) {
		return precedences[k]
	}
	return LowestPrec
}

// IsBinaryOperator reports whether k can appear as an infix operator.
func (k Kind) IsBinaryOperator() bool { return k.Precedence() > LowestPrec }

// JoinGT maps a run of adjacent GT tokens, optionally followed by an adjacent
// ASSIGN, to the operator kind it spells. The scanner deliberately
// under-munches `>` so that Array<Box<int32>> comes apart in type context
// (§4.2); the expression parser calls this after checking adjacency with
// tok[i].End == tok[i+1].Pos.
//
// Adjacency is the whitespace test, so `a > > b` never reaches this function
// with gts == 2.
func JoinGT(gts int, trailingAssign bool) Kind {
	switch gts {
	case 1:
		if trailingAssign {
			return GEQ
		}
		return GT
	case 2:
		if trailingAssign {
			return SHR_ASSIGN
		}
		return SHR
	case 3:
		if trailingAssign {
			return USHR_ASSIGN
		}
		return USHR
	}
	return INVALID
}

// ScannerEmits reports whether the scanner can produce k directly. It is false
// for exactly the five joined `>` forms. Used by the scanner's golden tests
// (§7) to assert that no fixture ever contains one.
func ScannerEmits(k Kind) bool {
	switch k {
	case GEQ, SHR, USHR, SHR_ASSIGN, USHR_ASSIGN:
		return false
	}
	return true
}

// itoa avoids importing strconv into the lowest package for one debug path.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [8]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}