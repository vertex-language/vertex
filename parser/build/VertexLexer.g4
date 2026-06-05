// VertexLexer.g4
// Vertex Language — Lexical Grammar, Specification 2.1
//
// Design notes:
//   • Newlines are plain whitespace — Vertex has no automatic semicolon
//     insertion, so no NEWLINE token is produced. Statement boundaries are
//     determined by the syntactic structure of the parser rules.
//   • Scalar type names (int, int32, float, bool, char, string, void …) are
//     identifiers, not keywords. This keeps the keyword set small and lets
//     users shadow them in theory (backend may warn).
//   • Concurrency postfix names (await, spawn, fork, dispatch, channel, new,
//     delete, try) are identifiers. They are method names resolved by the
//     backend, never reserved syntax.
//   • Result, Ok, Err are promoted to keywords so the parser can reference
//     them unambiguously in Result-construction and switch-pattern rules
//     without semantic predicates. Expected follows the same reasoning for
//     the same reason — it appears in typeExpr position after '->' and would
//     be ambiguous as a plain identifier there.
//   • test is promoted to a keyword so funcQualifier can reference it as a
//     token (matching async, thread, process, gpu). Because 'build test' is
//     the canonical test-file tag, buildDecl explicitly accepts TEST there.
//   • const is a keyword, but is only grammatically valid in pointer type
//     position (*const T). The backend rejects it elsewhere.
//   • inout, out, clobber are reserved for asm() constraint syntax (§47).
//     They would be very unusual variable names in systems code, and the
//     asm() form is confined to intrinsics packages anyway.
//   • reinterpret is a keyword so that 'reinterpret<T>(expr)' is unambiguous
//     at the parser level — the LT/GT tokens cannot be mistaken for
//     comparison operators when immediately preceded by this token.
//   • char values use single-quote literals ('A', '\n') — a distinct token  // ← UPDATED
//     from double-quoted STRING_LIT. The backend validates that exactly one  // ← UPDATED
//     Unicode code unit is represented. string continues to use "…".        // ← UPDATED

lexer grammar VertexLexer;

// ── Control-flow keywords ─────────────────────────────────────────────────────
LET         : 'let' ;
VAR         : 'var' ;
FUNC        : 'func' ;
IF          : 'if' ;
ELSE        : 'else' ;
FOR         : 'for' ;
IN          : 'in' ;          // also reused as the asm 'in' constraint keyword
WHILE       : 'while' ;
SWITCH      : 'switch' ;
CASE        : 'case' ;
DEFAULT     : 'default' ;
RETURN      : 'return' ;
BREAK       : 'break' ;
CONTINUE    : 'continue' ;
FALLTHROUGH : 'fallthrough' ;
DEFER       : 'defer' ;

// ── Declaration keywords ──────────────────────────────────────────────────────
STRUCT  : 'struct' ;
CLASS   : 'class' ;
ENUM    : 'enum' ;
TYPE    : 'type' ;
IMPORT  : 'import' ;
PACKAGE : 'package' ;
BUILD   : 'build' ;
ASM     : 'asm' ;

// ── Modifier keywords ─────────────────────────────────────────────────────────
WEAK     : 'weak' ;
CONST_KW : 'const' ;          // only valid in *const T — context enforced by backend

// ── Cast keywords ─────────────────────────────────────────────────────────────
// reinterpret<T>(expr) — raw pointer reinterpretation. Promoted to a keyword
// so the parser can unambiguously treat the immediately following '<' as the
// opening of a type-argument list rather than a less-than comparison.
// Backend validates: T must be a pointer type; expr must be a pointer or
// addressable value. Zero runtime cost — type annotation only.
REINTERPRET : 'reinterpret' ;

// ── Concurrency qualifiers (§36, §39–§41) ─────────────────────────────────────
ASYNC   : 'async' ;
THREAD  : 'thread' ;
PROCESS : 'process' ;
GPU     : 'gpu' ;

// ── Testing qualifier (compiler_testing §4.1) ─────────────────────────────────
// Promoted to a keyword so funcQualifier can reference it as a token,
// matching the pattern of the concurrency qualifiers above.
// 'build test' is handled in buildDecl by explicitly accepting TEST there.
TEST : 'test' ;

// ── Channel type keyword (§42) ────────────────────────────────────────────────
CHAN : 'chan' ;

// ── Map type keyword (§25) ────────────────────────────────────────────────────
MAP : 'map' ;

// ── Inline-assembly constraint keywords (§47) ─────────────────────────────────
// 'in' is reused from the control-flow keywords above (same token, IN).
// 'out' and 'inout' would be ambiguous as plain identifiers inside asm()
// bodies, so they are reserved. 'clobber' completes the set.
INOUT   : 'inout' ;
OUT_KW  : 'out' ;
CLOBBER : 'clobber' ;

// ── Result / error-handling names (§38.3) ─────────────────────────────────────
// Expected follows the same promotion rationale as Result: it appears in
// typeExpr position after '->' and must be unambiguous to the parser.
// Channel names (stdout, exitCode) are left as plain identifiers — they
// only appear inside Expected(...) where the backend validates them.
RESULT   : 'Result' ;
OK       : 'Ok' ;
ERR_KW   : 'Err' ;
EXPECTED : 'Expected' ;

// ── Boolean / nil literals (§1) ───────────────────────────────────────────────
TRUE  : 'true' ;
FALSE : 'false' ;
NIL   : 'nil' ;

// ── Multi-character operators ──────────────────────────────────────────────────
// Listed before their single-character prefixes so maximal-munch wins:
//   '...' before '.', '.<' before '.', etc.
ELLIPSIS     : '...' ;
HALF_OPEN    : '..<' ;
IDENTITY_EQ  : '===' ;
IDENTITY_NEQ : '!==' ;
OVERFLOW_ADD : '&+' ;
OVERFLOW_SUB : '&-' ;
OVERFLOW_MUL : '&*' ;
NIL_COALESCE : '??' ;
ARROW        : '->' ;
LSHIFT       : '<<' ;
RSHIFT       : '>>' ;
LEQ          : '<=' ;
GEQ          : '>=' ;
EQ           : '==' ;
NEQ          : '!=' ;
LOGICAL_AND  : '&&' ;
LOGICAL_OR   : '||' ;
PLUS_ASSIGN  : '+=' ;
MINUS_ASSIGN : '-=' ;
STAR_ASSIGN  : '*=' ;
DIV_ASSIGN   : '/=' ;
MOD_ASSIGN   : '%=' ;

// ── Single-character operators ─────────────────────────────────────────────────
PLUS     : '+' ;
MINUS    : '-' ;
STAR     : '*' ;
SLASH    : '/' ;
PERCENT  : '%' ;
TILDE    : '~' ;
AMP      : '&' ;
PIPE     : '|' ;
CARET    : '^' ;
BANG     : '!' ;
LT       : '<' ;
GT       : '>' ;
ASSIGN   : '=' ;
QUESTION : '?' ;
COLON    : ':' ;

// ── Delimiters ─────────────────────────────────────────────────────────────────
LPAREN   : '(' ;
RPAREN   : ')' ;
LBRACE   : '{' ;
RBRACE   : '}' ;
LBRACKET : '[' ;
RBRACKET : ']' ;
DOT      : '.' ;
COMMA    : ',' ;

// ── Numeric literals (§1) ──────────────────────────────────────────────────────
//
// Ordering enforces maximal-munch when two rules share a prefix:
//   HEX_FLOAT_LIT before HEX_INT_LIT  — '0xFp2' is a float, '0xFF' is an int
//   DEC_FLOAT_LIT before DEC_INT_LIT  — '3.14' is a float, '42' is an int
//
// Negative numeric literals do not exist at the token level. The unary MINUS
// operator in the expression grammar handles '-1000', '-3.14', etc.
//
// Underscore separators (1_000_000, 0xFF_FF) are consumed here and ignored
// by the compiler backend. They carry no semantic meaning.

HEX_FLOAT_LIT : '0' [xX] HEX_DIGIT+ ('.' HEX_DIGIT+)? [pP] [+\-]? DEC_DIGIT+ ;
HEX_INT_LIT   : '0' [xX] HEX_DIGIT (HEX_DIGIT | '_')* ;
OCT_INT_LIT   : '0' [oO] OCT_DIGIT (OCT_DIGIT | '_')* ;
BIN_INT_LIT   : '0' [bB] BIN_DIGIT (BIN_DIGIT | '_')* ;

// Decimal float: either 'digits.digits[exp]' or 'digits exp' (no bare dot).
// A bare '3.' is NOT a float literal; it tokenises as DEC_INT_LIT DOT.
DEC_FLOAT_LIT : DEC_SEQ '.' DEC_SEQ DEC_EXP?
              | DEC_SEQ DEC_EXP
              ;

DEC_INT_LIT   : DEC_DIGIT (DEC_DIGIT | '_')* ;

// ── String and character literals (§1, §3) ─────────────────────────────────────  // ← UPDATED
//
// char values use single-quote syntax ('A', '\n'). The lexer enforces exactly  // ← UPDATED
// one CHAR_CHAR (a raw code unit or a recognised escape sequence) between the  // ← UPDATED
// delimiters, so invalid forms like '' or 'AB' are rejected at lex time.       // ← UPDATED
// The backend additionally validates that the code unit fits the declared       // ← UPDATED
// char type width (e.g. ASCII-range for *const char in native interop).        // ← UPDATED
//
// string values continue to use double-quote syntax ("hello"). Both may        // ← UPDATED
// appear in the same source file without ambiguity — their opening             // ← UPDATED
// delimiters are distinct characters.                                          // ← UPDATED
//
// Multiline strings use backtick delimiters (§3). Content is verbatim:
// no escape sequences are processed, and no indentation is stripped.

CHAR_LIT            : '\'' CHAR_CHAR '\'' ;                                      // ← NEW
STRING_LIT          : '"'  STR_CHAR*  '"' ;
MULTILINE_STRING_LIT: '`'  .*?        '`' ;

// ── Identifiers ────────────────────────────────────────────────────────────────
// Placed after all keyword rules. ANTLR lexes the longest keyword match
// before falling through to IDENTIFIER, so 'let' is LET, not IDENTIFIER.
IDENTIFIER : ID_START ID_CONT* ;

// ── Whitespace and comments ────────────────────────────────────────────────────
// Newlines skipped: no statement-terminator token is produced.
WS           : [ \t\r\n]+ -> skip ;
LINE_COMMENT : '//' ~[\r\n]* -> skip ;

// ── Fragments ──────────────────────────────────────────────────────────────────
fragment HEX_DIGIT : [0-9a-fA-F] ;
fragment OCT_DIGIT : [0-7] ;
fragment BIN_DIGIT : [01] ;
fragment DEC_DIGIT : [0-9] ;
fragment DEC_SEQ   : DEC_DIGIT (DEC_DIGIT | '_')* ;
fragment DEC_EXP   : [eE] [+\-]? DEC_DIGIT+ ;
fragment CHAR_CHAR : ~['\\\r\n] | '\\' [nrtbf'\\] ;                              // ← NEW  single code unit or escape
fragment STR_CHAR  : ~["\\\r\n] | '\\' [nrtbf"\\] ;
fragment ID_START  : [a-zA-Z_] ;
fragment ID_CONT   : [a-zA-Z0-9_] ;