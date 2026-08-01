package token

import "strconv"

type Kind int

const (
	INVALID Kind = iota
	EOF
	COMMENT

	literalBeg
	IDENT  // foo, _, int32, sizeof, build, framework, test
	INT    // 12, 0xC3, 1_000, 0b1010
	FLOAT  // 1.5, 1e9, 0x1.8p3
	CHAR   // 'A', '\u{1F600}'
	STRING // "abc", `raw`
	literalEnd

	// ReservedLiteralKeyword (A.1.3): Literals syntactically, reserved lexically.
	// Separate from literalBeg..literalEnd because they carry no Lit text and
	// are keyword-looked-up, not scanned as literals.
	reservedLitBeg
	TRUE
	FALSE
	NIL
	reservedLitEnd

	operatorBeg
	LPAREN   // (
	RPAREN   // )
	LBRACK   // [
	RBRACK   // ]
	LBRACE   // {
	RBRACE   // }
	COMMA    // ,
	PERIOD   // .
	DOTDOT   // ..
	ELLIPSIS // ...
	COLON    // :
	ARROW    // ->

	ASSIGN     // =
	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=
	AND_ASSIGN // &=
	OR_ASSIGN  // |=
	XOR_ASSIGN // ^=
	SHL_ASSIGN // <<=
	SHR_ASSIGN // >>=

	ADD   // +
	SUB   // -
	MUL   // *
	QUO   // /
	REM   // %
	TILDE // ~   bitwise-NOT (A.4.4) or underlying-type (A.7.3)
	AND   // &   address-of / dereference / bitwise-AND (A.1.6)
	OR    // |
	XOR   // ^
	SHL   // 
	SHR   // >>

	WRAP_ADD // &+
	WRAP_SUB // &-
	WRAP_MUL // &*

	EQL      // ==
	NEQ      // !=
	LSS      // 
	GTR      // >
	LEQ      // <=
	GEQ      // >=
	IDENTICAL     // ===
	NOT_IDENTICAL // !==

	LAND // &&
	LOR  // ||
	NOT  // !
	operatorEnd

	keywordBeg
	ABSTRACT
	AS
	ASYNC
	AWAIT
	BREAK
	CASE
	CHAN
	CLASS
	CONSTRAINT
	CONTINUE
	DECLARE
	DEFAULT
	DEFER
	ELSE
	ENUM
	FALLTHROUGH
	FOR
	FUNC
	GPU
	IF
	IMPORT
	IN
	LET
	MAP
	MUT
	NPU
	PACKAGE
	RETURN
	SELECT
	SHARED
	STRUCT
	SWITCH
	TENSOR
	THREAD
	TYPE
	TYPED_PTR
	UNIQUE
	VAR
	WEAK
	WHILE
	keywordEnd
)

var kinds = [...]string{
	INVALID: "INVALID",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	TRUE:  "true",
	FALSE: "false",
	NIL:   "nil",

	LPAREN: "(", RPAREN: ")", LBRACK: "[", RBRACK: "]",
	LBRACE: "{", RBRACE: "}", COMMA: ",", PERIOD: ".",
	DOTDOT: "..", ELLIPSIS: "...", COLON: ":", ARROW: "->",

	ASSIGN: "=",
	ADD_ASSIGN: "+=", SUB_ASSIGN: "-=", MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=", REM_ASSIGN: "%=", AND_ASSIGN: "&=",
	OR_ASSIGN: "|=", XOR_ASSIGN: "^=", SHL_ASSIGN: "<<=", SHR_ASSIGN: ">>=",

	ADD: "+", SUB: "-", MUL: "*", QUO: "/", REM: "%",
	TILDE: "~", AND: "&", OR: "|", XOR: "^", SHL: "<<", SHR: ">>",
	WRAP_ADD: "&+", WRAP_SUB: "&-", WRAP_MUL: "&*",

	EQL: "==", NEQ: "!=", LSS: "<", GTR: ">", LEQ: "<=", GEQ: ">=",
	IDENTICAL: "===", NOT_IDENTICAL: "!==",
	LAND: "&&", LOR: "||", NOT: "!",

	ABSTRACT: "abstract", AS: "as", ASYNC: "async", AWAIT: "await",
	BREAK: "break", CASE: "case", CHAN: "chan", CLASS: "class",
	CONSTRAINT: "constraint", CONTINUE: "continue", DECLARE: "declare",
	DEFAULT: "default", DEFER: "defer", ELSE: "else", ENUM: "enum",
	FALLTHROUGH: "fallthrough", FOR: "for", FUNC: "func", GPU: "gpu",
	IF: "if", IMPORT: "import", IN: "in", LET: "let", MAP: "map",
	MUT: "mut", NPU: "npu", PACKAGE: "package", RETURN: "return",
	SELECT: "select", SHARED: "shared", STRUCT: "struct", SWITCH: "switch",
	TENSOR: "tensor", THREAD: "thread", TYPE: "type", TYPED_PTR: "typed_ptr",
	UNIQUE: "unique", VAR: "var", WEAK: "weak", WHILE: "while",
}

func (k Kind) String() string {
	if k >= 0 && int(k) < len(kinds) && kinds[k] != "" {
		return kinds[k]
	}
	return "Kind(" + strconv.Itoa(int(k)) + ")"
}

var keywords map[string]Kind

func init() {
	keywords = make(map[string]Kind, (keywordEnd-keywordBeg)+3)
	for k := keywordBeg + 1; k < keywordEnd; k++ {
		keywords[kinds[k]] = k
	}
	// ReservedLiteralKeyword: reserved lexically, so they go in the same table.
	keywords["true"] = TRUE
	keywords["false"] = FALSE
	keywords["nil"] = NIL
}

// Lookup maps an identifier string to its keyword Kind, or IDENT.
//
// ContextualKeywords are deliberately absent: A.1.3 makes them identifiers
// everywhere except the one production that names each. PredeclaredTypeName
// and ReservedBuiltinName (A.1.4) are absent for a different reason — they are
// ordinary identifiers pre-bound in an implicit scope, and the scanner must
// not know them.
func Lookup(ident string) Kind {
	if k, ok := keywords[ident]; ok {
		return k
	}
	return IDENT
}

func (k Kind) IsLiteral() bool {
	return literalBeg < k && k < literalEnd || reservedLitBeg < k && k < reservedLitEnd
}
func (k Kind) IsOperator() bool { return operatorBeg < k && k < operatorEnd }
func (k Kind) IsKeyword() bool  { return keywordBeg < k && k < keywordEnd }

// Precedence levels. Higher binds tighter; A.13 numbers them the other way.
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 9
	HighestPrec = 10
)

// Prec returns the binary precedence of k, or LowestPrec.
//
// DOTDOT is listed here at level 4 but is NON-ASSOCIATIVE (A.4.5): the parser
// must reject a..b..c rather than fold it. The precedence-climbing loop needs
// a special case, not a table entry difference.
func (k Kind) Prec() int {
	switch k {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, GTR, LEQ, GEQ, IDENTICAL, NOT_IDENTICAL:
		return 3
	case DOTDOT:
		return 4
	case ADD, SUB, OR, XOR, WRAP_ADD, WRAP_SUB:
		return 5
	case MUL, QUO, REM, AND, WRAP_MUL:
		return 6
	case SHL, SHR:
		return 7
	case AS:
		return 8 // RHS is a Type, not an Expr — parser handles specially
	}
	return LowestPrec
}

// IsCompoundAssign reports whether k is a CompoundAssignOperator (A.5.2).
func (k Kind) IsCompoundAssign() bool {
	switch k {
	case ADD_ASSIGN, SUB_ASSIGN, MUL_ASSIGN, QUO_ASSIGN, REM_ASSIGN,
		AND_ASSIGN, OR_ASSIGN, XOR_ASSIGN, SHL_ASSIGN, SHR_ASSIGN:
		return true
	}
	return false
}