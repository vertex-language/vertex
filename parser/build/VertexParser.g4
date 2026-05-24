// VertexParser.g4
// Vertex Language — Specification 1.9
//
// Backend responsibilities (not encoded here):
//   • Type checking, inference, and coercion
//   • Qualifier restrictions (.await only in async, .try only in Result fn, …)
//   • defer / break / continue placement rules
//   • Switch exhaustiveness for enums
//   • mut call-site validation (var binding + & prefix)
//   • init / deinit naming rules; weak only on .new() instances
//   • Type aliases at package level only
//   • nativeFuncDecl only valid inside a native class (one with ': parentName')
//   • ELLIPSIS only valid as the final nativeParam
//   • Struct literals may not appear directly as the condition of if / for /
//     switch — wrap in parentheses to disambiguate (§25)

parser grammar VertexParser;
options { tokenVocab = VertexLexer; }


// ═══════════════════════════════════════════════════════════════════
// File — §31
// ═══════════════════════════════════════════════════════════════════

file
    : packageDecl? buildDecl? importDecl* topLevelDecl* EOF
    ;

packageDecl
    : PACKAGE ID
    ;

buildDecl
    : BUILD buildTag (COMMA buildTag)*
    ;

buildTag
    : ID
    ;

importDecl
    : IMPORT STRING_LIT
    | IMPORT LPAREN STRING_LIT+ RPAREN
    ;


// ═══════════════════════════════════════════════════════════════════
// Top-level declarations
// ═══════════════════════════════════════════════════════════════════

topLevelDecl
    : funcDecl
    | structDecl
    | classDecl
    | enumDecl
    | typeAliasDecl
    | varDeclStmt
    ;


// ─── §19 — Function declaration ─────────────────────────────────────

funcDecl
    : FUNC ID genericParams? LPAREN paramList? RPAREN funcQualifier? returnType? block
    ;

genericParams
    : LT typeParam (COMMA typeParam)* GT
    ;

typeParam
    : ID
    ;

paramList
    : param (COMMA param)*
    ;

param
    : ID COLON MUT? type
    ;

funcQualifier
    : ASYNC
    | THREAD
    | PROCESS
    | GPU
    ;

returnType
    : ARROW type
    ;


// ─── §25 — Struct ───────────────────────────────────────────────────
//
// Declaration — fields are let or var, all must be provided at init site.
// Instantiation uses brace-literal syntax: TypeName{field: value, …}
// Trailing commas are valid. Partial initialization is a compile error.
// All field labels are required — positional initialization is not supported.
// Fields may appear in any order inside the literal.
// Struct literals may NOT appear directly as the condition of if / for /
// switch — backend enforces this; wrap in parentheses to disambiguate.

structDecl
    : STRUCT ID genericParams? LBRACE structField* RBRACE
    ;

structField
    : (LET | VAR) ID COLON type
    ;

// Brace-literal instantiation — Point{x: 3, y: 4}
// An empty struct literal is written as TypeName{} — the inner list is optional.
structLiteralExpr
    : ID LBRACE (structFieldInit (COMMA structFieldInit)* COMMA?)? RBRACE
    ;

structFieldInit
    : ID COLON expr
    ;


// ─── §28 — Class ────────────────────────────────────────────────────
//
// A class with ': parentName' is a native class — a zero-size compile-time
// dispatch surface bound to the import whose namespace matches parentName.
// The parentName is the last path segment of the import:
//
//   import "lib/sdl2"       →  namespace "sdl2"   →  class SDL2  : sdl2  { … }
//   import "linux/syscalls" →  namespace "syscalls" →  class Sys  : syscalls { … }
//   import "gpu/cuda"       →  namespace "cuda"    →  class Cuda : cuda   { … }
//
// Native class members are method declarations without bodies; the backend
// owns the implementation.  Regular classes have only fields.

classDecl
    : CLASS ID (COLON ID)? LBRACE classMember* RBRACE
    ;

classMember
    : classField
    | nativeFuncDecl
    ;

classField
    : (LET | VAR) ID COLON type
    ;

// Native method declaration — no body.
nativeFuncDecl
    : FUNC ID LPAREN nativeParamList? RPAREN returnType?
    ;

nativeParamList
    : nativeParam (COMMA nativeParam)*
    ;

// ELLIPSIS is the C variadic tail; must be last and only in native classes.
nativeParam
    : ID COLON MUT? type
    | ELLIPSIS
    ;


// ─── §27 — Enum ─────────────────────────────────────────────────────

enumDecl
    : ENUM ID (COLON enumRawType)? LBRACE enumCaseDecl+ RBRACE
    ;

enumRawType
    : INT
    | STRING
    ;

enumCaseDecl
    : CASE enumCase (COMMA enumCase)*
    ;

enumCase
    : ID (ASSIGN literal)?
    ;


// ─── §4 FFI — Type alias ────────────────────────────────────────────

typeAliasDecl
    : TYPE ID ASSIGN type
    ;


// ═══════════════════════════════════════════════════════════════════
// Statements
// ═══════════════════════════════════════════════════════════════════

block
    : LBRACE stmt* RBRACE
    ;

stmt
    : varDeclStmt
    | assignStmt
    | compoundAssignStmt
    | ifStmt
    | switchStmt
    | forInStmt
    | whileStmt
    | breakStmt
    | continueStmt
    | fallthroughStmt
    | returnStmt
    | deferStmt
    | exprStmt
    ;


// ─── Variable declaration ────────────────────────────────────────────

varDeclStmt
    : bindingKw tupleBind             ASSIGN expr
    | bindingKw tupleBind COLON type  ASSIGN expr
    | bindingKw ID        COLON type  ASSIGN expr
    | bindingKw ID                    ASSIGN expr
    | bindingKw ID        COLON type
    ;

bindingKw
    : WEAK LET
    | LET
    | VAR
    ;

tupleBind
    : LPAREN ID (COMMA ID)+ RPAREN
    ;


// ─── Assignment ───────────────────────────────────────────────────────

assignStmt
    : lvalue ASSIGN expr
    ;

compoundAssignStmt
    : lvalue compoundOp expr
    ;

compoundOp
    : PLUS_ASSIGN
    | MINUS_ASSIGN
    | STAR_ASSIGN
    | SLASH_ASSIGN
    | MOD_ASSIGN
    ;

lvalue
    : ID
    | lvalue DOT ID
    | lvalue LBRACKET expr RBRACKET
    ;


// ─── If / Else ────────────────────────────────────────────────────────

ifStmt
    : IF ifCondition block elseIfClause* elseClause?
    ;

elseIfClause
    : ELSE IF ifCondition block
    ;

elseClause
    : ELSE block
    ;

ifCondition
    : LET ID ASSIGN expr
    | expr
    ;


// ─── Switch ───────────────────────────────────────────────────────────

switchStmt
    : SWITCH expr LBRACE switchCase* defaultCase? RBRACE
    ;

switchCase
    : CASE casePatternList COLON stmt+
    ;

defaultCase
    : DEFAULT COLON stmt+
    ;

casePatternList
    : casePattern (COMMA casePattern)*
    ;

casePattern
    : literal
    | DOT ID
    | OK  LPAREN LET ID RPAREN
    | ERR LPAREN LET ID RPAREN
    ;


// ─── For-in ───────────────────────────────────────────────────────────

forInStmt
    : FOR ID IN expr block
    ;


// ─── While ────────────────────────────────────────────────────────────

whileStmt
    : WHILE expr block
    ;


// ─── Jump statements ──────────────────────────────────────────────────

breakStmt       : BREAK       ;
continueStmt    : CONTINUE    ;
fallthroughStmt : FALLTHROUGH ;

returnStmt
    : RETURN expr?
    ;


// ─── Defer ────────────────────────────────────────────────────────────

deferStmt
    : DEFER expr
    ;


// ─── Expression statement ─────────────────────────────────────────────

exprStmt
    : expr
    ;


// ═══════════════════════════════════════════════════════════════════
// Expressions
// ═══════════════════════════════════════════════════════════════════

expr
    : expr DOT postfixName LPAREN argList? RPAREN
    | expr DOT postfixName
    | expr LBRACKET expr RBRACKET
    | expr LPAREN argList? RPAREN

    | expr (LSHIFT | RSHIFT) expr

    | expr AMP   expr
    | expr CARET expr
    | expr PIPE  expr

    | expr (STAR | SLASH | PERCENT | OVERFLOW_MUL) expr

    | expr (PLUS | MINUS | OVERFLOW_ADD | OVERFLOW_SUB) expr

    | expr (ELLIPSIS | HALF_OPEN_RANGE) expr

    | <assoc=right> expr NIL_COALESCE expr

    | expr (EQ | NEQ | LT | GT | LTE | GTE | IDENTITY_EQ | IDENTITY_NEQ) expr

    | expr AND expr

    | expr OR expr

    | <assoc=right> expr QUESTION expr COLON expr

    | MINUS  expr
    | BANG   expr
    | TILDE  expr
    | AMP    expr

    | primitiveType LPAREN expr RPAREN

    | RESULT LPAREN (OK | ERR) COMMA expr RPAREN

    | primary
    ;


postfixName
    : ID
    | ANY
    ;


primary
    : literal
    | structLiteralExpr
    | ID
    | GPU
    | ASYNC
    | THREAD
    | PROCESS
    | DOT ID
    | LPAREN RPAREN
    | LPAREN expr RPAREN
    | tupleExpr
    | arrayLiteralExpr
    | arrayConstructExpr
    | dictLiteralExpr
    | anonFuncExpr
    ;


argList
    : arg (COMMA arg)*
    ;

arg
    : (ID COLON)? expr
    ;


literal
    : DEC_INT_LIT
    | HEX_INT_LIT
    | BIN_INT_LIT
    | OCT_INT_LIT
    | DEC_FLOAT_LIT
    | HEX_FLOAT_LIT
    | STRING_LIT
    | RAW_STRING_LIT
    | TRUE
    | FALSE
    | NIL
    ;


arrayLiteralExpr
    : LBRACKET (exprList COMMA?)? RBRACKET
    ;

arrayConstructExpr
    : LBRACKET type RBRACKET LPAREN (ID COLON expr COMMA ID COLON expr)? RPAREN
    ;

exprList
    : expr (COMMA expr)*
    ;

dictLiteralExpr
    : LBRACKET dictEntry (COMMA dictEntry)* COMMA? RBRACKET
    ;

dictEntry
    : expr COLON expr
    ;


tupleExpr
    : LPAREN tupleElement (COMMA tupleElement)+ RPAREN
    ;

tupleElement
    : (ID COLON)? expr
    ;


anonFuncExpr
    : FUNC LPAREN paramList? RPAREN returnType? block
    ;


// ═══════════════════════════════════════════════════════════════════
// Types
// ═══════════════════════════════════════════════════════════════════

type
    : primitiveType
    | ID (LT type (COMMA type)* GT)?
    | LBRACKET type RBRACKET
    | LBRACKET type COLON type RBRACKET
    | LPAREN tupleTypeElem (COMMA tupleTypeElem)+ RPAREN
    | LPAREN RPAREN
    | FUNC LPAREN funcTypeParamList? RPAREN (ARROW type)?
    | CHANNEL type
    | ANY      type
    | MUT ANY  type
    | ANY      VOID
    | MUT ANY  VOID
    | ANY      OPAQUE
    | MUT ANY  OPAQUE
    | RESULT LPAREN type COMMA type RPAREN
    | type QUESTION
    ;

tupleTypeElem
    : (ID COLON)? type
    ;

funcTypeParamList
    : funcTypeParam (COMMA funcTypeParam)*
    ;

funcTypeParam
    : MUT? type
    ;


primitiveType
    : INT    | INT8   | INT16  | INT32  | INT64
    | UINT   | UINT8  | UINT16 | UINT32 | UINT64
    | FLOAT  | DOUBLE
    | BOOL   | STRING | CHAR   | VOID
    ;