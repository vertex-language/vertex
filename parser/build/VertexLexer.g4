// ============================================================
// VertexLexer.g4  –  Vertex 6 Lexical Grammar (ANTLR4)
//
// Reference: The Vertex Programming Language 6.3, §Lexical Structure
// ============================================================

lexer grammar VertexLexer;
options {
    superClass = VertexLexerBase;
}

// ─────────────────────────────────────────────────────────────
// CHANNELS
// ─────────────────────────────────────────────────────────────
channels { BLOCK_COMMENT_CHANNEL, LINE_COMMENT_CHANNEL }

// ─────────────────────────────────────────────────────────────
// WHITESPACE AND COMMENTS
// ─────────────────────────────────────────────────────────────
WS      : [ \t\u000B\u000C\u0000]+ -> channel(HIDDEN) ;
NEWLINE : ('\r\n' | '\r' | '\n') -> channel(HIDDEN) ;

LINE_COMMENT  : '//' ~[\r\n]* -> channel(LINE_COMMENT_CHANNEL) ;
BLOCK_COMMENT : '/*' (BLOCK_COMMENT | .)*? '*/' -> channel(BLOCK_COMMENT_CHANNEL) ;

// ─────────────────────────────────────────────────────────────
// KEYWORDS – Declarations
// ─────────────────────────────────────────────────────────────
ASSOCIATEDTYPE : 'associatedtype' ;
CLASS          : 'class' ;
DEINIT         : 'deinit' ;
ENUM           : 'enum' ;
EXTENSION      : 'extension' ;
FUNC           : 'func' ;
IMPORT         : 'import' ;
INIT           : 'init' ;
INOUT          : 'inout' ;
LET            : 'let' ;
OPEN           : 'open' ;
PRIVATE        : 'private' ;
PROTOCOL       : 'protocol' ;
PUBLIC         : 'public' ;
STATIC         : 'static' ;
STRUCT         : 'struct' ;
SUBSCRIPT      : 'subscript' ;
TYPEALIAS      : 'typealias' ;
VAR            : 'var' ;
INTERNAL       : 'internal' ;

// ─────────────────────────────────────────────────────────────
// KEYWORDS – Statements
// ─────────────────────────────────────────────────────────────
BREAK       : 'break' ;
CASE        : 'case' ;
CATCH       : 'catch' ;
CONTINUE    : 'continue' ;
DEFAULT     : 'default' ;
DEFER       : 'defer' ;
DO          : 'do' ;
ELSE        : 'else' ;
FALLTHROUGH : 'fallthrough' ;
FOR         : 'for' ;
GUARD       : 'guard' ;
IF          : 'if' ;
IN          : 'in' ;
REPEAT      : 'repeat' ;
RETURN      : 'return' ;
THROW       : 'throw' ;
SWITCH      : 'switch' ;
WHERE       : 'where' ;
WHILE       : 'while' ;

// ─────────────────────────────────────────────────────────────
// KEYWORDS – Expressions and Types
// ─────────────────────────────────────────────────────────────
ANY          : 'Any' ;
AS           : 'as' ;
AWAIT        : 'await' ;
FALSE        : 'false' ;
IS           : 'is' ;
NIL          : 'nil' ;
SELF         : 'self' ;
SELF_UPPER   : 'Self' ;
SUPER        : 'super' ;
THROWS       : 'throws' ;
TRUE         : 'true' ;
TRY          : 'try' ;
CONSUME      : 'consume' ;
COPY         : 'copy' ;
DISCARD      : 'discard' ;
BORROW       : 'borrow' ;

// ─────────────────────────────────────────────────────────────
// PATTERNS
// ─────────────────────────────────────────────────────────────
UNDERSCORE : '_' ;

// ─────────────────────────────────────────────────────────────
// HOISTED SOFT KEYWORDS / CONTEXT-SENSITIVE TOKENS
// Required here to satisfy split lexer/parser rules.
// ─────────────────────────────────────────────────────────────
OS_KW                : 'os' ;
ARCH_KW              : 'arch' ;
VERTEX_KW            : 'vertex' ;
COMPILER_KW          : 'compiler' ;
CANIMPORT_KW         : 'canImport' ;
VERSION_KW           : 'version' ;
FILE_KEYWORD         : 'file' ;
LINE_KEYWORD         : 'line' ;
GET_KW               : 'get' ;
SET_KW               : 'set' ;
WILLSET_KW           : 'willSet' ;
DIDSET_KW            : 'didSet' ;
ASYNC_KW             : 'async' ;
FINAL_KW             : 'final' ;
ACTOR_KW             : 'actor' ;
PREFIX_KW            : 'prefix' ;
POSTFIX_KW           : 'postfix' ;
MACRO_KW             : 'macro' ;
DYNAMIC_KW           : 'dynamic' ;
LAZY_KW              : 'lazy' ;
OPTIONAL_KW          : 'optional' ;
OVERRIDE_KW          : 'override' ;
REQUIRED_KW          : 'required' ;
UNOWNED_KW           : 'unowned' ;
WEAK_KW              : 'weak' ;
NONISOLATED_KW       : 'nonisolated' ;
MUTATING_KW          : 'mutating' ;
NONMUTATING_KW       : 'nonmutating' ;
SOME_KW              : 'some' ;
TYPE_KW              : 'Type' ;
PROTOCOL_KW          : 'Protocol' ;
CONSUMING_KW         : 'consuming' ;
BORROWING_KW         : 'borrowing' ;
SENDING_KW           : 'sending' ;

// ─────────────────────────────────────────────────────────────
// POUND KEYWORDS
// ─────────────────────────────────────────────────────────────
POUND_AVAILABLE       : '#available' ;
POUND_UNAVAILABLE     : '#unavailable' ;
POUND_IF              : '#if' ;
POUND_ELSEIF          : '#elseif' ;
POUND_ELSE            : '#else' ;
POUND_ENDIF           : '#endif' ;
POUND_SOURCE_LOCATION : '#sourceLocation' ;
POUND_FILE            : '#file' ;
POUND_FILEID          : '#fileID' ;
POUND_FILEPATH        : '#filePath' ;
POUND_LINE            : '#line' ;
POUND_COLUMN          : '#column' ;
POUND_FUNCTION        : '#function' ;

// ─────────────────────────────────────────────────────────────
// PUNCTUATION
// ─────────────────────────────────────────────────────────────
LPAREN    : '(' ;
RPAREN    : ')' ;
LBRACE    : '{' ;
RBRACE    : '}' ;
LBRACKET  : '[' ;
RBRACKET  : ']' ;
DOT       : '.' ;
COMMA     : ',' ;
COLON     : ':' ;
SEMICOLON : ';' ;
ASSIGN    : '=' ;
AT        : '@' ;
HASH      : '#' ;
AMPERSAND : '&' ;
ARROW     : '->' ;
BACKTICK  : '`' ;
ELLIPSIS  : '...' ;
RANGE_HALF_OPEN : '..<' ;
LT        : '<' ;
GT        : '>' ;

EXCLAIM_POSTFIX : '!'  ; 
QUESTION_POSTFIX: '?'  ;

// ─────────────────────────────────────────────────────────────
// LITERALS – Numeric
// ─────────────────────────────────────────────────────────────
INTEGER_LITERAL
    : DecimalLiteral
    | BinaryLiteral
    | OctalLiteral
    | HexadecimalLiteral
    ;

FLOAT_LITERAL
    : DecimalLiteral '.' DecimalLiteral DecimalExponent?
    | DecimalLiteral DecimalExponent
    | '0x' HexadecimalDigits ('.' HexadecimalDigits)? HexadecimalExponent
    ;

fragment DecimalLiteral          : [0-9] ([0-9_]* [0-9])? ;
fragment BinaryLiteral           : '0b' [01] ([01_]* [01])? ;
fragment OctalLiteral            : '0o' [0-7] ([0-7_]* [0-7])? ;
fragment HexadecimalLiteral      : '0x' HexadecimalDigits ;
fragment HexadecimalDigits       : [0-9a-fA-F] ([0-9a-fA-F_]* [0-9a-fA-F])? ;
fragment DecimalExponent         : [eE] [+\-]? DecimalLiteral ;
fragment HexadecimalExponent     : [pP] [+\-]? DecimalLiteral ;

// ─────────────────────────────────────────────────────────────
// LITERALS – String
// ─────────────────────────────────────────────────────────────
STRING_LITERAL           : '"' StringChar* '"' ;
MULTILINE_STRING_LITERAL : '"""' MultilineStringChar*? '"""' ;

EXTENDED_STRING_LITERAL
    : '#"'   ExtendedStringChar*? '"#'
    | '##"'  ExtendedStringChar*? '"##'
    | '###"' ExtendedStringChar*? '"###'
    ;

fragment StringChar
    : ~["\\\r\n]
    | '\\' EscapeSeq
    ;

fragment MultilineStringChar
    : ~[\\"]
    | '\\' EscapeSeq
    | '"' ~["]
    | '""' ~["]
    ;

fragment ExtendedStringChar
    : ~["]
    | '"' ~[#]
    ;

fragment EscapeSeq
    : '0'
    | '\\'
    | 't'
    | 'n'
    | 'r'
    | '"'
    | 'u{' HexDigit+ '}'
    ;

fragment HexDigit : [0-9a-fA-F] ;

// ─────────────────────────────────────────────────────────────
// IDENTIFIERS
// ─────────────────────────────────────────────────────────────
IDENTIFIER
    : IdentHead IdentChar*
    | '`' IdentHead IdentChar* '`'
    | '$' [0-9]+
    ;

fragment IdentHead
    : [a-zA-Z_]
    | '\u00A8' | '\u00AA' | '\u00AD' | '\u00AF'
    | [\u00B2-\u00B5] | [\u00B7-\u00BA]
    | [\u00BC-\u00BE] | [\u00C0-\u00D6]
    | [\u00D8-\u00F6] | [\u00F8-\u00FF]
    | [\u0100-\u02FF] | [\u0370-\u167F]
    | [\u1681-\u180D] | [\u180F-\u1DBF]
    | [\u1E00-\u1FFF]
    | [\u200B-\u200D] | [\u202A-\u202E]
    | [\u203F-\u2040] | '\u2054'
    | [\u2060-\u206F] | [\u2070-\u20CF]
    | [\u2100-\u218F] | [\u2460-\u24FF]
    | [\u2776-\u2793] | [\u2C00-\u2DFF]
    | [\u2E80-\u2FFF] | [\u3004-\u3007]
    | [\u3021-\u302F] | [\u3031-\u303F]
    | [\u3040-\uD7FF] | [\uF900-\uFD3D]
    | [\uFD40-\uFDCF] | [\uFDF0-\uFE1F]
    | [\uFE30-\uFE44] | [\uFE47-\uFFFD]
    ;

fragment IdentChar
    : IdentHead
    | [0-9]
    | '\u0300'..'\u036F'
    | '\u1DC0'..'\u1DFF'
    | '\u20D0'..'\u20FF'
    | '\uFE20'..'\uFE2F'
    ;

// ─────────────────────────────────────────────────────────────
// OPERATORS
// ─────────────────────────────────────────────────────────────
OPERATOR
    : OperatorHead OperatorChar*
    | DotOperatorHead DotOperatorChar+
    ;

fragment OperatorHead
    : [/=\-+!*%<>&|^~?]
    | '\u00A1'..'\u00A7'
    | '\u00A9' | '\u00AB' | '\u00AC' | '\u00AE'
    | '\u00B0'..'\u00B1' | '\u00B6' | '\u00BB'
    | '\u00BF' | '\u00D7' | '\u00F7'
    | '\u2016'..'\u2017' | '\u2020'..'\u2027'
    | '\u2030'..'\u203E'
    | '\u2041'..'\u2053'
    | '\u2055'..'\u205E'
    | '\u2190'..'\u23FF'
    | '\u2500'..'\u2775'
    | '\u2794'..'\u2BFF'
    | '\u2E00'..'\u2E7F'
    | '\u3001'..'\u3003'
    | '\u3008'..'\u3020'
    | '\u3030'
    ;

fragment OperatorChar
    : OperatorHead
    | '\u0300'..'\u036F'
    | '\u1DC0'..'\u1DFF'
    | '\u20D0'..'\u20FF'
    | '\uFE00'..'\uFE0F'
    | '\uFE20'..'\uFE2F'
    | '\u{E0100}'..'\u{E01EF}'
    ;

fragment DotOperatorHead : '..' ;
fragment DotOperatorChar : '.' | OperatorChar ;