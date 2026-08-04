package token

import "strconv"

// Kind classifies a token.
//
// Values are assigned by iota in declaration order, and related kinds are
// grouped into contiguous ranges bounded by unexported sentinels, so every
// classification predicate below is a pair of comparisons rather than a switch.
// Renumbering is safe as long as a group's sentinels move with it.
type Kind int

const (
	INVALID Kind = iota
	EOF
	COMMENT

	// IDENT is every identifier: the blank identifier `_`, every contextual
	// keyword, every predeclared type / tensor-element / constraint name, and
	// every reserved builtin name. None of those are distinguished lexically,
	// and the scanner does not know them.
	IDENT

	// Scanned literals. Each carries its raw source spelling in Token.Lit.
	literalBeg
	INT    // 42, 1_000, 0b1010, 0o600, 0xBadFace
	FLOAT  // 1.5, 1e9, 6.674_28e-11, 0x1.8p3
	CHAR   // 'A', '\n', '\u{1F600}'
	STRING // "abc", `raw`
	literalEnd

	// Reserved literal keywords: literals syntactically, reserved lexically.
	// Their spelling is fixed, so they carry no Lit — which is why they sit in
	// their own range rather than with the scanned literals.
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

	assignBeg
	ASSIGN // =
	compoundAssignBeg
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
	compoundAssignEnd
	assignEnd

	ADD   // +
	SUB   // -
	MUL   // *
	QUO   // /
	REM   // %
	TILDE // ~   bitwise-NOT, or underlying-type in a TypeSetTerm
	AND   // &   address-of, dereference, or bitwise-AND
	OR    // |
	XOR   // ^
	SHL   // 
	SHR   // >>

	WRAP_ADD // &+
	WRAP_SUB // &-
	WRAP_MUL // &*

	EQL           // ==
	NEQ           // !=
	LSS           // 
	GTR           // >
	LEQ           // <=
	GEQ           // >=
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
	VECTOR
	WEAK
	WHILE
	keywordEnd
)

// names is indexed by Kind. Operators, keywords, and reserved literal keywords
// hold their source spelling; everything else holds a category name.
var names = [...]string{
	INVALID: "INVALID",
	EOF:     "EOF",
	COMMENT: "COMMENT",
	IDENT:   "IDENT",

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

	ASSIGN:     "=",
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
	UNIQUE: "unique", VAR: "var", VECTOR: "vector", WEAK: "weak",
	WHILE: "while",
}

// Spelling returns k's fixed source text, or "" if k has none.
//
// Only operators, keywords, and reserved literal keywords have one. IDENT and
// the scanned literals vary per token and carry their text in Token.Lit;
// INVALID, EOF, and COMMENT are categories, not lexemes.
func (k Kind) Spelling() string {
	switch {
	case k.IsOperator(), k.IsKeyword(), reservedLitBeg < k && k < reservedLitEnd:
		return names[k]
	}
	return ""
}

// String renders k for diagnostics: its spelling where it has one, its category
// name otherwise.
func (k Kind) String() string {
	if k >= 0 && int(k) < len(names) && names[k] != "" {
		return names[k]
	}
	return "Kind(" + strconv.Itoa(int(k)) + ")"
}

var keywords map[string]Kind

func init() {
	keywords = make(map[string]Kind, (keywordEnd-keywordBeg)+3)
	for k := keywordBeg + 1; k < keywordEnd; k++ {
		keywords[names[k]] = k
	}
	// Reserved literal keywords are reserved lexically, so they share the
	// table despite living in a separate Kind range.
	keywords["true"] = TRUE
	keywords["false"] = FALSE
	keywords["nil"] = NIL
}

// Lookup maps an identifier spelling to its keyword Kind, or IDENT.
//
// Contextual keywords are deliberately absent: each is an ordinary identifier
// everywhere except the single production that names it, so baking one in here
// would make it reserved unconditionally. Predeclared type, tensor-element, and
// constraint names, and reserved builtin names, are absent for a different
// reason — they are identifiers pre-bound in an implicit outermost scope, and
// the scanner must not know them at all.
func Lookup(ident string) Kind {
	if k, ok := keywords[ident]; ok {
		return k
	}
	return IDENT
}

// IsLiteral reports whether k is a BasicLit kind: a scanned literal or a
// reserved literal keyword. IDENT is not one — an identifier is an operand
// name, not a literal.
func (k Kind) IsLiteral() bool {
	return literalBeg < k && k < literalEnd ||
		reservedLitBeg < k && k < reservedLitEnd
}

// HasLit reports whether a token of kind k carries raw source text in Lit.
func (k Kind) HasLit() bool {
	return k == IDENT || k == COMMENT || k == INVALID ||
		literalBeg < k && k < literalEnd
}

func (k Kind) IsOperator() bool { return operatorBeg < k && k < operatorEnd }
func (k Kind) IsKeyword() bool  { return keywordBeg < k && k < keywordEnd }

// IsAssign reports whether k is `=` or a compound assignment operator.
func (k Kind) IsAssign() bool { return assignBeg < k && k < assignEnd }

// IsCompoundAssign reports whether k is one of the ten assign_op spellings,
// excluding plain `=`. The compound form takes exactly one target and one
// value; the plain form takes lists.
func (k Kind) IsCompoundAssign() bool {
	return compoundAssignBeg < k && k < compoundAssignEnd
}

// IsUnaryOp reports whether k is a unary_op. `&` is not one — it derives
// through PointerPrimary and binds tighter than a selector — and neither
// `await` nor `var` is, each being its own UnaryExpr alternative.
func (k Kind) IsUnaryOp() bool {
	switch k {
	case SUB, NOT, TILDE:
		return true
	}
	return false
}

// Precedence ladder. Higher binds tighter.
//
// Binary operators occupy the seven levels the grammar lists. CastPrec places
// `as` above all of them, but Prec never returns it: `as` is written as
// CastExpr rather than as a binary_op because its right operand is a Type, so
// a precedence-climbing loop cannot consume it and must not try.
const (
	LowestPrec  = 0 // not a binary operator
	CastPrec    = 8 // `as`
	UnaryPrec   = 9
	HighestPrec = 10
)

// Prec returns k's binary precedence, or LowestPrec if k is not a binary_op.
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
	}
	return LowestPrec
}

// IsBinaryOp reports whether k may join two expressions.
func (k Kind) IsBinaryOp() bool { return k.Prec() > LowestPrec }

// IsNonAssociative reports whether k may not be folded in either direction.
// `..` is the only such operator: a..b..c is a compile error, and a precedence
// table alone cannot say so.
func (k Kind) IsNonAssociative() bool { return k == DOTDOT }

// EndsOperand reports whether a token of kind k can close an operand, which is
// the condition on the one restriction to longest-match scanning.
//
// A float_lit may not begin immediately after a `.` whose own preceding token
// has this property; there the scanner produces an int_lit and the inner `.`
// scans separately, which is what makes `t.0.0` yield two TupleIndex chains
// instead of `t` `.` `0.0`.
//
// The set is exactly the one written in the rule: identifier, `)`, `]`, `}`,
// int_lit, float_lit, string_lit. char_lit and true/false/nil are absent, and
// the scanner must not generalize the set — the restriction is narrow on
// purpose, and `1.5` stays a float because its preceding token is not a `.`.
func (k Kind) EndsOperand() bool {
	switch k {
	case IDENT, RPAREN, RBRACK, RBRACE, INT, FLOAT, STRING:
		return true
	}
	return false
}