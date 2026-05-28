// VSXLexer.g4
// Vertex Source XML (VSX) — Lexer
// Standalone. No Vertex grammar dependency.

lexer grammar VSXLexer;

// ═══════════════════════════════════════════════════════════════════
// Whitespace
// ═══════════════════════════════════════════════════════════════════

VSX_WS : [ \t\r\n]+ -> skip ;

// ═══════════════════════════════════════════════════════════════════
// Punctuation
// ═══════════════════════════════════════════════════════════════════

VSX_LPAREN    : '('  ;
VSX_RPAREN    : ')'  ;
VSX_LBRACKET  : '['  ;
VSX_RBRACKET  : ']'  ;
VSX_LBRACE    : '{'  ;
VSX_RBRACE    : '}'  ;
VSX_LT        : '<'  ;
VSX_GT        : '>'  ;
VSX_SLASH     : '/'  ;
VSX_ASSIGN    : '='  ;
VSX_COLON     : ':'  ;
VSX_DOT       : '.'  ;
VSX_COMMA     : ','  ;
VSX_QUESTION  : '?'  ;
VSX_BANG      : '!'  ;
VSX_AMP       : '&'  ;
VSX_PIPE      : '|'  ;
VSX_SPREAD    : '...' ;

// ═══════════════════════════════════════════════════════════════════
// Operators
// ═══════════════════════════════════════════════════════════════════

VSX_AND       : '&&' ;
VSX_OR        : '||' ;
VSX_EQ        : '==' ;
VSX_NEQ       : '!=' ;
VSX_LTE       : '<=' ;
VSX_GTE       : '>=' ;
VSX_PLUS      : '+'  ;
VSX_MINUS     : '-'  ;
VSX_STAR      : '*'  ;
VSX_PERCENT   : '%'  ;
VSX_ARROW_FN  : '=>' ;  // anonymous func arrow form

// ═══════════════════════════════════════════════════════════════════
// Keywords needed inside expr slots
// ═══════════════════════════════════════════════════════════════════

VSX_NULL      : 'null'      ;
VSX_TRUE      : 'true'      ;
VSX_FALSE     : 'false'     ;
VSX_FUNC      : 'func'      ;   // Vertex anon func keyword
VSX_RETURN    : 'return'    ;
VSX_NEW       : 'new'       ;

// ═══════════════════════════════════════════════════════════════════
// Literals
// ═══════════════════════════════════════════════════════════════════

VSX_NUMBER
    : [0-9]+ ('.' [0-9]+)?
    | '0x' [0-9a-fA-F]+
    ;

// Attribute string values — JSX spec: no backslash escapes, HTML refs allowed
VSX_DOUBLE_STRING : '"'  (~["{<>&] | HTMLCharRef)* '"'  ;
VSX_SINGLE_STRING : '\'' (~['{<>&] | HTMLCharRef)* '\'' ;

// Template-style string (backtick) for Vertex raw strings inside expr slots
VSX_BACKTICK_STRING : '`' .*? '`' ;

// ═══════════════════════════════════════════════════════════════════
// JSX text — raw character data between tags
// Excludes { } < > as those switch context
// ═══════════════════════════════════════════════════════════════════

VSX_TEXT : (VSXTextChar | HTMLCharRef)+ ;

// ═══════════════════════════════════════════════════════════════════
// Identifier
// ═══════════════════════════════════════════════════════════════════

VSX_IDENT : VSXIdentStart VSXIdentPart* ('-' VSXIdentPart+)* ;

// ═══════════════════════════════════════════════════════════════════
// Fragments
// ═══════════════════════════════════════════════════════════════════

fragment VSXIdentStart      : [a-zA-Z_$\u00C0-\uFFFF]           ;
fragment VSXIdentPart       : [a-zA-Z0-9_$\u00C0-\uFFFF]        ;
fragment VSXTextChar        : ~[{}<>]                             ;

fragment HTMLCharRef
    : '&' '#' [0-9]+ ';'
    | '&' '#' [xX] [0-9a-fA-F]+ ';'
    | '&' [a-zA-Z] [a-zA-Z0-9]* ';'
    ;