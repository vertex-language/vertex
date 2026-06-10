// VertexLexer.g4
// Vertex Language — Lexical Grammar, Specification 2.1
//
// Design notes:
//   • Newlines are plain whitespace — Vertex has no automatic semicolon
//     insertion, so no NEWLINE token is produced. Statement boundaries are
//     determined by the syntactic structure of the parser rules.
//   • Scalar type names (int, int32, float32, float64, bool, char, string, void …) are
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
//   • as is promoted to a keyword so 'expr as typeExpr' is unambiguous in
//     binary position. It is the unified cast and reinterpretation operator
//     (§17.5), replacing the former reinterpret<T>(expr) prefix form. The
//     backend determines the operation from the operand types at compile time:
//     pointer↔pointer is a no-op reinterpretation; numeric combinations are
//     widening, truncation, or float↔int conversion as appropriate.
//   • char values use single-quote literals ('A', '\n') — a distinct token
//     from double-quoted STRING_LIT. The backend validates that exactly one
//     Unicode code unit is represented. string continues to use "…".
//   • SEMI (';') is used exclusively as the size separator inside fixed array
//     type expressions: [T; N]. It is not a statement terminator — Vertex
//     uses no semicolon insertion and no statement-ending semicolons.

lexer grammar VertexLexer;

// ── Control-flow keywords ─────────────────────────────────────────────────────
LET         : 'let' ;
VAR         : 'var' ;
FUNC        : 'func' ;
IF          : 'if' ;
ELSE        : 'else' ;
FOR         : 'for' ;
IN          : 'in' ;
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
CONST_KW : 'const' ;

// ── Cast keyword (§6.1) ───────────────────────────────────────────────────────
AS : 'as' ;

// ── Concurrency qualifiers (§36, §39–§41) ─────────────────────────────────────
ASYNC   : 'async' ;
THREAD  : 'thread' ;
PROCESS : 'process' ;
GPU     : 'gpu' ;

// ── Testing qualifier (§48) ───────────────────────────────────────────────────
TEST : 'test' ;

// ── Channel type keyword (§42) ────────────────────────────────────────────────
CHAN : 'chan' ;

// ── Map type keyword (§25) ────────────────────────────────────────────────────
MAP : 'map' ;

// ── Inline-assembly constraint keywords (§47) ─────────────────────────────────
INOUT   : 'inout' ;
OUT_KW  : 'out' ;
CLOBBER : 'clobber' ;

// ── Result / error-handling names (§38.3) ─────────────────────────────────────
RESULT   : 'Result' ;
OK       : 'Ok' ;
ERR_KW   : 'Err' ;
EXPECTED : 'Expected' ;

// ── Boolean / nil literals (§1) ───────────────────────────────────────────────
TRUE  : 'true' ;
FALSE : 'false' ;
NIL   : 'nil' ;

// ── Multi-character operators ──────────────────────────────────────────────────
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
SEMI     : ';' ;    // fixed array size separator: [T; N] — not a statement terminator

// ── Numeric literals (§1) ──────────────────────────────────────────────────────
HEX_FLOAT_LIT : '0' [xX] HEX_DIGIT+ ('.' HEX_DIGIT+)? [pP] [+\-]? DEC_DIGIT+ ;
HEX_INT_LIT   : '0' [xX] HEX_DIGIT (HEX_DIGIT | '_')* ;
OCT_INT_LIT   : '0' [oO] OCT_DIGIT (OCT_DIGIT | '_')* ;
BIN_INT_LIT   : '0' [bB] BIN_DIGIT (BIN_DIGIT | '_')* ;

DEC_FLOAT_LIT : DEC_SEQ '.' DEC_SEQ DEC_EXP?
              | DEC_SEQ DEC_EXP
              ;

DEC_INT_LIT   : DEC_DIGIT (DEC_DIGIT | '_')* ;

// ── String and character literals (§1, §3) ─────────────────────────────────────
CHAR_LIT            : '\'' CHAR_CHAR '\'' ;
STRING_LIT          : '"'  STR_CHAR*  '"' ;
MULTILINE_STRING_LIT: '`'  .*?        '`' ;

// ── Identifiers ────────────────────────────────────────────────────────────────
IDENTIFIER : ID_START ID_CONT* ;

// ── Whitespace and comments ────────────────────────────────────────────────────
WS           : [ \t\r\n]+ -> skip ;
LINE_COMMENT : '//' ~[\r\n]* -> skip ;

// ── Fragments ──────────────────────────────────────────────────────────────────
fragment HEX_DIGIT : [0-9a-fA-F] ;
fragment OCT_DIGIT : [0-7] ;
fragment BIN_DIGIT : [01] ;
fragment DEC_DIGIT : [0-9] ;
fragment DEC_SEQ   : DEC_DIGIT (DEC_DIGIT | '_')* ;
fragment DEC_EXP   : [eE] [+\-]? DEC_DIGIT+ ;
fragment CHAR_CHAR : ~['\\\r\n] | '\\' [nrtbf'\\] ;
fragment STR_CHAR  : ~["\\\r\n] | '\\' [nrtbf"\\] ;
fragment ID_START  : [a-zA-Z_] ;
fragment ID_CONT   : [a-zA-Z0-9_] ;