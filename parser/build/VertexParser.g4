// VertexParser.g4
// Vertex Language — Syntactic Grammar, Specification 2.1
//
// ─────────────────────────────────────────────────────────────────────────────
// OPERATOR PRECEDENCE  (high → low, per §17)
//
//   postfix   .method<T>()  .method()  .field  [i]  ()
//   prefix    - ! ~   &  (address-of)
//   binary    as T                                   (§17.5 cast / reinterpret,
//                                                    highest-precedence binary op)
//             << >>
//             * / % &*
//             + - &+ &-
//             & ^ |
//             ... ..
//             ??                                    (right-associative)
//             == != < > <= >= === !==
//             &&
//             ||
//             ? :                                    (right-associative)
//   statement = += -= *= /= %=
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
//   Generic calls and struct literals are anchored to qualifiedIdent (not
//   arbitrary expr) in the seed alternatives — this eliminates the classic
//   angle-bracket ambiguity with comparison operators.  a < b > (c) parses
//   as comparisons; identity<int32>(v) parses as a generic call because
//   'identity' is a qualifiedIdent seed, not a binary-expr continuation.
//
//   Generic method calls (expr DOT IDENTIFIER LT…GT LPAREN) are a left-
//   recursive postfix alternative and must precede the plain method-call
//   alternative so ANTLR4 tries the longer match first.
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
//   • Expected(channel, string) accepts any IDENTIFIER as the channel name;
//     backend validates it is a known channel (stdout).
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

// §32  Unconstrained type parameters only in 2.1.
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
//
// genericParams? added — supports struct Box<T> { value: T }
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

varDecl
    : WEAK? (LET | VAR) bindingPattern (COLON typeExpr)? ASSIGN expr
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
//
// Generic call and struct-literal seeds use qualifiedIdent rather than expr
// to eliminate angle-bracket ambiguity with comparison operators (see file
// header note).  Generic method calls use a left-recursive postfix form and
// are listed before the plain method-call alternative.
// ════════════════════════════════════════════════════════════════════════════

expr
    // ── Postfix (highest precedence) ─────────────────────────────────────────
    : expr DOT IDENTIFIER LT typeExpr (COMMA typeExpr)* GT LPAREN argList? RPAREN  // generic method call: a.method<T>(args)
    | expr DOT IDENTIFIER LPAREN argList? RPAREN                                   // method call:         a.method(args)
    | expr DOT IDENTIFIER                                                           // field access:        a.field
    | expr LBRACKET expr RBRACKET                                                   // subscript:           a[i]
    | expr LPAREN argList? RPAREN                                                   // call:                f(args)

    // ── Cast / reinterpret (§17.5) ────────────────────────────────────────────
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
    // Generic forms before plain forms — ANTLR4 tries these first.
    // qualifiedIdent covers both simple names (Box) and qualified names
    // (yourpackage.Box), keeping all three generic primary forms uniform.
    | qualifiedIdent LT typeExpr (COMMA typeExpr)* GT LPAREN argList? RPAREN       // generic call:         identity<int32>(args)
    | qualifiedIdent LT typeExpr (COMMA typeExpr)* GT LBRACE structLiteralFields? RBRACE  // generic struct lit:   Box<int32>{f: v}
    | qualifiedIdent LBRACE structLiteralFields? RBRACE                            // struct literal:       Box{f: v} / pkg.Box{f: v}
    | LBRACE mapLiteralFields? RBRACE                                              // map literal:          {"k": v}
    | MAP LBRACKET typeExpr RBRACKET typeExpr LPAREN argList? RPAREN               // map alloc:            map[K]V()
    | literal                                                                      // §1 literals
    | IDENTIFIER                                                                   // variable / type name
    | DOT IDENTIFIER                                                               // enum shorthand:       .caseName
    | LPAREN expr RPAREN                                                           // grouping
    | LPAREN expr (COMMA expr)+ RPAREN                                             // tuple literal:        (a, b, …)
    | LPAREN RPAREN                                                                // empty tuple / void:   ()
    | LBRACKET (expr (COMMA expr)* COMMA?)? RBRACKET                               // array literal:        [a, b, …]
    | anonFuncExpr                                                                 // anonymous function
    | RESULT LPAREN (OK | ERR_KW) COMMA expr RPAREN                               // Result(Ok,v) / Result(Err,e)
    | asmExpr                                                                      // inline assembly
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
// Generic instantiation added: baseType LT typeExpr (COMMA typeExpr)* GT
// Listed before plain baseType so ANTLR4 tries the longer match first.
// Covers all type positions: return types, annotations, struct fields, etc.
//   func foo() -> Box<int32>
//   var x: Box<string> = ...
//   struct Pair<A, B> { first: A   second: B }
// ════════════════════════════════════════════════════════════════════════════

typeExpr
    : STAR CONST_KW? typeExpr QUESTION?                              // §4   *T, *const T, *T?
    | LBRACKET typeExpr RBRACKET                                     // §24  [T] array
    | MAP LBRACKET typeExpr RBRACKET typeExpr                        // §25  map[K]V
    | CHAN typeExpr                                                   // §42  chan T
    | FUNC LPAREN funcTypeParams? RPAREN (ARROW typeExpr)?           // §34  func(T)->T
    | LPAREN tupleTypeElems? RPAREN                                  // §37  (T,T) tuple / ()
    | RESULT LPAREN typeExpr COMMA typeExpr RPAREN                   // §38.3 Result(T,E)
    | EXPECTED LPAREN typeExpr COMMA STRING_LIT RPAREN               // compiler_testing §4.2
    | baseType LT typeExpr (COMMA typeExpr)* GT QUESTION?            // §32  generic: Box<T>, Box<T>?
    | baseType QUESTION?                                             // named type T, or optional T?
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