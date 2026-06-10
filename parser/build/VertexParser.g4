// VertexParser.g4
// Vertex Language — Syntactic Grammar, Specification 2.1
//
// ─────────────────────────────────────────────────────────────────────────────
// OPERATOR PRECEDENCE  (high → low, per §17)
//
//   postfix   .method<T>()  .method()  .field  [i]  ()
//   prefix    - ! ~   &  (address-of)
//   binary    as T                                   (§6.1 cast / reinterpret)
//             << >>
//             * / % &*
//             + - &+ &-
//             & ^ |
//             ... ..<
//             ??                                     (right-associative)
//             == != < > <= >= === !==
//             &&
//             ||
//             ? :                                    (right-associative)
//   statement = += -= *= /= %=
//
// ─────────────────────────────────────────────────────────────────────────────
// ARRAY TYPES  (§24)
//
//   [T]      — dynamic heap array. var required; push/pop valid.
//   [T; N]   — fixed stack array. Size N is a compile-time expr and is part
//              of the type — [uint8; 1024] and [uint8; 512] are distinct types.
//              Zero-fill is implicit when no initializer is provided.
//
//   Declaration forms:
//     var buf: [uint8; 1024]              — fixed, zero-fill, no = required
//     var buf: [uint8; 1024] = [1, 2, …]  — fixed, explicit initializer
//     var x:   [int32] = []              — dynamic, empty
//     var x  = [1, 2, 3]                — dynamic (var + literal)
//     let x  = [1, 2, 3]                — fixed immutable (let + literal)
//
//   varDecl alt 1 handles all [T; N] annotations (with or without =).
//   varDecl alt 2 handles all other declarations.
//   The SEMI token inside the brackets disambiguates the two alternatives.
//
// ─────────────────────────────────────────────────────────────────────────────
// GENERIC SYNTAX  (§32)
//
//   Declaration:     func identity<T>(value: T) -> T
//                    struct Box<T> { value: T }
//   Call:            identity<int32>(value: 42)
//   Method call:     obj.method<int32>(args)
//   Struct literal:  Box<int32>{value: 99}
//   Qualified:       yourpackage.Box<string>{value: "world"}
//   Type position:   -> Box<int32>   var x: Box<string>
//
// ─────────────────────────────────────────────────────────────────────────────
// INTENTIONAL OVER-PERMISSIVENESS (backend narrows these)
//
//   • struct/class/enum declarations are allowed inside function bodies.
//   • asm() is accepted anywhere; backend rejects it outside intrinsics packages.
//   • Struct literals are accepted in if/for/switch conditions; backend rejects.
//   • typeAliasDecl is accepted inside blocks; backend restricts to package scope.
//   • Assignments and compound assignments share exprOrAssignStmt; the backend
//     validates that the left-hand expression is an addressable lvalue.
//   • expr AS typeExpr accepts any expr on the left and any typeExpr on the
//     right; the backend validates the conversion/reinterpretation.
//   • Fixed array initializers [v1, v2, …] are array literal exprs; the backend
//     validates element count matches N in [T; N].
// ─────────────────────────────────────────────────────────────────────────────

parser grammar VertexParser;
options { tokenVocab = VertexLexer; }


// ════════════════════════════════════════════════════════════════════════════
// §45–46, §33  FILE STRUCTURE
// ════════════════════════════════════════════════════════════════════════════

file
    : packageDecl? buildDecl* importDecl* topLevelDecl* EOF
    ;

packageDecl
    : PACKAGE IDENTIFIER
    ;

buildDecl
    : BUILD IDENTIFIER
    | BUILD TEST
    ;

importDecl
    : IMPORT STRING_LIT
    | IMPORT LPAREN STRING_LIT+ RPAREN
    ;


// ════════════════════════════════════════════════════════════════════════════
// TOP-LEVEL DECLARATIONS
// ════════════════════════════════════════════════════════════════════════════

topLevelDecl
    : funcDecl
    | structDecl
    | classDecl
    | enumDecl
    | typeAliasDecl
    | varDecl
    ;

typeAliasDecl
    : TYPE IDENTIFIER ASSIGN typeExpr
    ;


// ════════════════════════════════════════════════════════════════════════════
// §21, §28, §32, §36, §39–41  FUNCTION DECLARATIONS
// ════════════════════════════════════════════════════════════════════════════

funcDecl
    : FUNC receiver? IDENTIFIER genericParams?
      LPAREN paramList? RPAREN
      funcQualifier? (ARROW typeExpr)?
      block
    ;

receiver
    : LPAREN IDENTIFIER COLON typeExpr RPAREN
    ;

genericParams
    : LT IDENTIFIER (COMMA IDENTIFIER)* GT
    ;

funcQualifier
    : ASYNC
    | THREAD
    | PROCESS
    | GPU
    | TEST
    ;

paramList
    : param (COMMA param)* (COMMA variadicParam)? COMMA?
    | variadicParam COMMA?
    ;

param
    : IDENTIFIER COLON typeExpr
    ;

variadicParam
    : IDENTIFIER COLON ELLIPSIS typeExpr
    ;


// ════════════════════════════════════════════════════════════════════════════
// §27  STRUCT DECLARATIONS
// ════════════════════════════════════════════════════════════════════════════

structDecl
    : STRUCT IDENTIFIER genericParams? LBRACE structFieldDecl* RBRACE
    ;

structFieldDecl
    : IDENTIFIER COLON typeExpr
    ;


// ════════════════════════════════════════════════════════════════════════════
// §30, §44  CLASS DECLARATIONS
// ════════════════════════════════════════════════════════════════════════════

classDecl
    : CLASS IDENTIFIER (COLON qualifiedIdent)?
      LBRACE classMember* RBRACE
    ;

classMember
    : IDENTIFIER COLON typeExpr
    | FUNC IDENTIFIER LPAREN paramList? RPAREN (ARROW typeExpr)?
    ;

qualifiedIdent
    : IDENTIFIER (DOT IDENTIFIER)*
    ;


// ════════════════════════════════════════════════════════════════════════════
// §29  ENUM DECLARATIONS
// ════════════════════════════════════════════════════════════════════════════

enumDecl
    : ENUM IDENTIFIER (COLON typeExpr)? LBRACE enumCaseDecl+ RBRACE
    ;

enumCaseDecl
    : CASE enumCaseItem (COMMA enumCaseItem)*
    ;

enumCaseItem
    : IDENTIFIER (ASSIGN enumRawValue)?
    ;

enumRawValue
    : MINUS? (DEC_INT_LIT | HEX_INT_LIT | OCT_INT_LIT | BIN_INT_LIT)
    | STRING_LIT
    ;


// ════════════════════════════════════════════════════════════════════════════
// STATEMENTS  §18–§31
// ════════════════════════════════════════════════════════════════════════════

block
    : LBRACE stmt* RBRACE
    ;

stmt
    : varDecl
    | ifStmt
    | whileStmt
    | forInStmt
    | switchStmt
    | returnStmt
    | BREAK
    | CONTINUE
    | FALLTHROUGH
    | deferStmt
    | structDecl
    | classDecl
    | enumDecl
    | typeAliasDecl
    | exprOrAssignStmt
    ;

// ── §24  Variable declarations ────────────────────────────────────────────────
//
// Alt 1 — fixed array annotation: var buf: [T; N]  or  var buf: [T; N] = [...]
//   The SEMI inside the brackets is the distinguishing token. The initializer
//   is optional — omitting it implies zero-fill. When present, the backend
//   validates that the literal contains exactly N elements.
//
// Alt 2 — standard form: var x = expr  or  var x: T = expr
//   Covers dynamic arrays (var x: [T] = []), scalars, structs, and everything
//   else. Always requires ASSIGN.
//
varDecl
    : WEAK? (LET | VAR) bindingPattern COLON LBRACKET typeExpr SEMI expr RBRACKET (ASSIGN expr)?
    | WEAK? (LET | VAR) bindingPattern (COLON typeExpr)? ASSIGN expr
    ;

bindingPattern
    : IDENTIFIER
    | LPAREN IDENTIFIER (COMMA IDENTIFIER)+ RPAREN
    ;

ifStmt
    : IF ifCondition block (ELSE (block | ifStmt))?
    ;

ifCondition
    : LET IDENTIFIER ASSIGN expr
    | expr
    ;

whileStmt
    : WHILE expr block
    ;

forInStmt
    : FOR IDENTIFIER IN expr block
    ;

switchStmt
    : SWITCH expr LBRACE switchCase* RBRACE
    ;

switchCase
    : CASE switchPattern (COMMA switchPattern)* COLON stmt+
    | DEFAULT COLON stmt+
    ;

switchPattern
    : expr
    | DOT IDENTIFIER
    | (OK | ERR_KW) LPAREN LET IDENTIFIER RPAREN
    ;

returnStmt
    : RETURN expr?
    ;

deferStmt
    : DEFER expr
    ;

exprOrAssignStmt
    : expr (assignOp expr)?
    ;

assignOp
    : ASSIGN
    | PLUS_ASSIGN | MINUS_ASSIGN | STAR_ASSIGN | DIV_ASSIGN | MOD_ASSIGN
    ;


// ════════════════════════════════════════════════════════════════════════════
// EXPRESSIONS  §7–§17, §24–§26, §34–§35, §38
// ════════════════════════════════════════════════════════════════════════════

expr
    // ── Postfix ───────────────────────────────────────────────────────────────
    : expr DOT IDENTIFIER LT typeExpr (COMMA typeExpr)* GT LPAREN argList? RPAREN
    | expr DOT IDENTIFIER LPAREN argList? RPAREN
    | expr DOT IDENTIFIER
    | expr LBRACKET expr RBRACKET
    | expr LPAREN argList? RPAREN

    // ── Cast / reinterpret (§6.1) ─────────────────────────────────────────────
    | expr AS typeExpr

    // ── Binary operators (high → low) ─────────────────────────────────────────
    | expr (LSHIFT | RSHIFT) expr
    | expr (STAR | SLASH | PERCENT | OVERFLOW_MUL) expr
    | expr (PLUS | MINUS | OVERFLOW_ADD | OVERFLOW_SUB) expr
    | expr AMP   expr
    | expr CARET expr
    | expr PIPE  expr
    | expr (ELLIPSIS | HALF_OPEN) expr
    | <assoc=right> expr NIL_COALESCE expr
    | expr (EQ | NEQ | LT | GT | LEQ | GEQ
           | IDENTITY_EQ | IDENTITY_NEQ) expr
    | expr LOGICAL_AND expr
    | expr LOGICAL_OR  expr
    | <assoc=right> expr QUESTION expr COLON expr

    // ── Prefix / unary ────────────────────────────────────────────────────────
    | (MINUS | BANG | TILDE) expr
    | AMP expr

    // ── Primary expressions ───────────────────────────────────────────────────
    | qualifiedIdent LT typeExpr (COMMA typeExpr)* GT LPAREN argList? RPAREN
    | qualifiedIdent LT typeExpr (COMMA typeExpr)* GT LBRACE structLiteralFields? RBRACE
    | qualifiedIdent LBRACE structLiteralFields? RBRACE
    | LBRACE mapLiteralFields? RBRACE
    | MAP LBRACKET typeExpr RBRACKET typeExpr LPAREN argList? RPAREN
    | literal
    | IDENTIFIER
    | DOT IDENTIFIER
    | LPAREN expr RPAREN
    | LPAREN expr (COMMA expr)+ RPAREN
    | LPAREN RPAREN
    | LBRACKET (expr (COMMA expr)* COMMA?)? RBRACKET    // array literal: [], [1,2,3]
    | anonFuncExpr
    | RESULT LPAREN (OK | ERR_KW) COMMA expr RPAREN
    | asmExpr
    ;

argList
    : arg (COMMA arg)* COMMA?
    ;

arg
    : (IDENTIFIER COLON)? expr
    ;

structLiteralFields
    : structLiteralField (COMMA structLiteralField)* COMMA?
    ;

structLiteralField
    : IDENTIFIER COLON expr
    ;

mapLiteralFields
    : mapLiteralField (COMMA mapLiteralField)* COMMA?
    ;

mapLiteralField
    : expr COLON expr
    ;

anonFuncExpr
    : FUNC LPAREN paramList? RPAREN funcQualifier? (ARROW typeExpr)? block
    ;

asmExpr
    : ASM LPAREN asmBody RPAREN
    ;

asmBody
    : asmInstr (COMMA asmInstr)*
      (COMMA asmConstraint)*
      (COMMA asmClobberDecl)?
    ;

asmInstr
    : STRING_LIT
    ;

asmConstraint
    : IN     LPAREN STRING_LIT RPAREN IDENTIFIER
    | INOUT  LPAREN STRING_LIT RPAREN IDENTIFIER
    | OUT_KW LPAREN STRING_LIT RPAREN
    ;

asmClobberDecl
    : CLOBBER LPAREN STRING_LIT (COMMA STRING_LIT)* RPAREN
    ;


// ════════════════════════════════════════════════════════════════════════════
// TYPE EXPRESSIONS  §3–§4, §24, §34, §37, §38.3, §42
//
// [T; N] is listed before [T] so ANTLR4 tries the longer match first.
// After LBRACKET typeExpr, the SEMI token unambiguously selects [T; N]
// over [T]. N is an expr — the backend validates it is a compile-time
// integer literal.
// ════════════════════════════════════════════════════════════════════════════

typeExpr
    : STAR CONST_KW? typeExpr QUESTION?                              // §4   *T, *const T, *T?
    | LBRACKET typeExpr SEMI expr RBRACKET                           // §24  [T; N] fixed array
    | LBRACKET typeExpr RBRACKET                                     // §24  [T]   dynamic array
    | MAP LBRACKET typeExpr RBRACKET typeExpr                        // §25  map[K]V
    | CHAN typeExpr                                                   // §42  chan T
    | FUNC LPAREN funcTypeParams? RPAREN (ARROW typeExpr)?           // §34  func(T)->T
    | LPAREN tupleTypeElems? RPAREN                                  // §37  (T,T) tuple / ()
    | RESULT LPAREN typeExpr COMMA typeExpr RPAREN                   // §38.3 Result(T,E)
    | EXPECTED LPAREN typeExpr COMMA STRING_LIT RPAREN               // §48  compiler testing
    | baseType LT typeExpr (COMMA typeExpr)* GT QUESTION?            // §32  generic: Box<T>
    | baseType QUESTION?                                             // named type T, or T?
    ;

funcTypeParams
    : typeExpr (COMMA typeExpr)* COMMA?
    ;

baseType
    : IDENTIFIER (DOT IDENTIFIER)*
    ;

tupleTypeElems
    : tupleTypeElem (COMMA tupleTypeElem)* COMMA?
    ;

tupleTypeElem
    : (IDENTIFIER COLON)? typeExpr
    ;


// ════════════════════════════════════════════════════════════════════════════
// §1  LITERALS
// ════════════════════════════════════════════════════════════════════════════

literal
    : DEC_INT_LIT
    | HEX_INT_LIT
    | OCT_INT_LIT
    | BIN_INT_LIT
    | DEC_FLOAT_LIT
    | HEX_FLOAT_LIT
    | CHAR_LIT
    | STRING_LIT
    | MULTILINE_STRING_LIT
    | TRUE
    | FALSE
    | NIL
    ;