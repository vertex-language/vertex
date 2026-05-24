// VertexLexer.g4
// Vertex Language — Specification 1.9
//
// Tokenisation only. All semantic constraints are backend responsibilities.

lexer grammar VertexLexer;

// ═══════════════════════════════════════════════════════════════════
// Comments (consumed before any other rule)
// ═══════════════════════════════════════════════════════════════════

BLOCK_COMMENT : '/*' .*? '*/'  -> skip ;
LINE_COMMENT  : '//' ~[\r\n]*  -> skip ;

// ═══════════════════════════════════════════════════════════════════
// Keywords — must appear before ID so maximal-munch picks them up
// ═══════════════════════════════════════════════════════════════════

// Control flow
BREAK       : 'break'       ;
CASE        : 'case'        ;
CONTINUE    : 'continue'    ;
DEFAULT     : 'default'     ;
DEFER       : 'defer'       ;
ELSE        : 'else'        ;
FALLTHROUGH : 'fallthrough' ;
FOR         : 'for'         ;
IF          : 'if'          ;
IN          : 'in'          ;
RETURN      : 'return'      ;
SWITCH      : 'switch'      ;
WHILE       : 'while'       ;

// Declarations
BUILD       : 'build'       ;
CLASS       : 'class'       ;
ENUM        : 'enum'        ;
FUNC        : 'func'        ;
IMPORT      : 'import'      ;
LET         : 'let'         ;
PACKAGE     : 'package'     ;
STRUCT      : 'struct'      ;
TYPE        : 'type'        ;
VAR         : 'var'         ;
WEAK        : 'weak'        ;

// Modifiers / qualifiers
MUT         : 'mut'         ;
ASYNC       : 'async'       ;
THREAD      : 'thread'      ;
PROCESS     : 'process'     ;
GPU         : 'gpu'         ;

// Pointer / FFI
ANY         : 'any'         ;
CHANNEL     : 'channel'     ;
OPAQUE      : 'opaque'      ;

// Result type vocabulary — reserved; backend enforces usage context
RESULT      : 'Result'      ;
OK          : 'Ok'          ;
ERR         : 'Err'         ;

// Value keywords
FALSE       : 'false'       ;
NIL         : 'nil'         ;
TRUE        : 'true'        ;

// ─── Built-in primitive types ──────────────────────────────────────

INT     : 'int'    ;
INT8    : 'int8'   ;
INT16   : 'int16'  ;
INT32   : 'int32'  ;
INT64   : 'int64'  ;
UINT    : 'uint'   ;
UINT8   : 'uint8'  ;
UINT16  : 'uint16' ;
UINT32  : 'uint32' ;
UINT64  : 'uint64' ;
FLOAT   : 'float'  ;
DOUBLE  : 'double' ;
BOOL    : 'bool'   ;
STRING  : 'string' ;
CHAR    : 'char'   ;
VOID    : 'void'   ;

// ═══════════════════════════════════════════════════════════════════
// Operators — multi-character tokens must precede their single-char
// prefixes so ANTLR's maximal-munch picks the right token.
// ═══════════════════════════════════════════════════════════════════

// §8 — Overflow operators (before AMP / PLUS / MINUS / STAR)
OVERFLOW_ADD : '&+' ;
OVERFLOW_SUB : '&-' ;
OVERFLOW_MUL : '&*' ;

// §6 — Compound assignment (before ASSIGN and arithmetic singles)
PLUS_ASSIGN  : '+=' ;
MINUS_ASSIGN : '-=' ;
STAR_ASSIGN  : '*=' ;
SLASH_ASSIGN : '/=' ;
MOD_ASSIGN   : '%=' ;

// §14 — Identity operators (before EQ / NEQ)
IDENTITY_EQ  : '===' ;
IDENTITY_NEQ : '!==' ;

// §9 — Equality (before ASSIGN / BANG)
EQ  : '==' ;
NEQ : '!=' ;

// §9 — Ordered comparison (before ASSIGN)
GTE : '>=' ;
LTE : '<=' ;

// §7 — Bit-shift (before GT / LT)
LSHIFT : '<<' ;
RSHIFT : '>>' ;

// §10 — Logical (before AMP / PIPE)
AND : '&&' ;
OR  : '||' ;

// Function return-type separator (before MINUS)
ARROW : '->' ;

// §13 — Nil-coalescing (before QUESTION)
NIL_COALESCE : '??' ;

// §11 — Range operators — both start with '.'; must precede DOT.
HALF_OPEN_RANGE : '..<' ;
ELLIPSIS        : '...' ;

// ─── Single-character operators ────────────────────────────────────
PLUS    : '+' ;
MINUS   : '-' ;
STAR    : '*' ;
SLASH   : '/' ;
PERCENT : '%' ;
// AMP is bitwise-AND infix AND the mut-param address-of prefix (&var).
// The parser distinguishes by syntactic position.
AMP     : '&' ;
PIPE    : '|' ;
CARET   : '^' ;
TILDE   : '~' ;
BANG    : '!' ;
GT      : '>' ;
LT      : '<' ;
ASSIGN  : '=' ;
QUESTION: '?' ;

// ═══════════════════════════════════════════════════════════════════
// Punctuation
// ═══════════════════════════════════════════════════════════════════

LPAREN    : '(' ;
RPAREN    : ')' ;
LBRACE    : '{' ;
RBRACE    : '}' ;
LBRACKET  : '[' ;
RBRACKET  : ']' ;
COLON     : ':' ;
COMMA     : ',' ;
DOT       : '.' ;
// Semicolons are silently discarded; they are not statement terminators.
SEMICOLON : ';' -> skip ;

// ═══════════════════════════════════════════════════════════════════
// Numeric literals — more specific / longer patterns first
// ═══════════════════════════════════════════════════════════════════

// §1 — Hex float   0xFp2   0xFp-2   0xC.3p0
HEX_FLOAT_LIT
    : '0' [xX] HEX_DIGITS ('.' HEX_DIGITS)? [pP] [+\-]? DEC_DIGITS
    ;

// §1 — Hex integer   0xFF   0xBadFace   0x0123_4567_89ab_cdef
HEX_INT_LIT
    : '0' [xX] HEX_DIGITS
    ;

// §1 — Binary integer   0b101010
BIN_INT_LIT
    : '0' [bB] [01] [01_]*
    ;

// §1 — Octal integer   0o52
OCT_INT_LIT
    : '0' [oO] [0-7] [0-7_]*
    ;

// §1 — Decimal float   3.14   1_000.000_1   1.25e2   1.25E-2
// Must precede DEC_INT_LIT so "3.14" is not lexed as "3" then ".14".
DEC_FLOAT_LIT
    : DEC_DIGITS '.' DEC_DIGITS ([eE] [+\-]? DEC_DIGITS)?
    | DEC_DIGITS [eE] [+\-]? DEC_DIGITS
    ;

// §1 — Decimal integer   42   1_000_000
DEC_INT_LIT
    : DEC_DIGITS
    ;

// ═══════════════════════════════════════════════════════════════════
// String literals
// ═══════════════════════════════════════════════════════════════════

// §1 — Double-quoted, single-line.  Used for both string and char values;
// the type system (backend) distinguishes them by annotation.
STRING_LIT
    : '"' (~["\\\r\n] | '\\' .)* '"'
    ;

// §3 — Backtick multiline string.  No indentation stripping (backend note).
RAW_STRING_LIT
    : '`' .*? '`'
    ;

// ═══════════════════════════════════════════════════════════════════
// Identifier — after all keywords
// ═══════════════════════════════════════════════════════════════════

ID : [a-zA-Z_] [a-zA-Z0-9_]* ;

// ═══════════════════════════════════════════════════════════════════
// Whitespace — Vertex is newline-insensitive; all WS skipped
// ═══════════════════════════════════════════════════════════════════

WS : [ \t\r\n]+ -> skip ;

// ═══════════════════════════════════════════════════════════════════
// Fragments
// ═══════════════════════════════════════════════════════════════════

fragment DEC_DIGITS : [0-9] [0-9_]* ;
fragment HEX_DIGITS : [0-9a-fA-F] [0-9a-fA-F_]* ;