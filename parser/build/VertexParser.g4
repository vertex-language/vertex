// VertexParser.g4
// Vertex Language — Syntactic Grammar, Specification 2.1
//
// ─────────────────────────────────────────────────────────────────────────────
// OPERATOR PRECEDENCE  (high → low, per §17)
//
//   postfix   .method()  .field  [i]  ()            ← listed first in expr
//   prefix    - ! ~   &  (address-of)
//   binary    << >>
//             * / % &*
//             + - &+ &-
//             ... ..
//             ??                                    (right-associative)
//             == != < > <= >= === !==
//             &&
//             ||
//             ? :                                    (right-associative)
//   statement = += -= *= /= %=                       ← NOT part of expr
//
// ─────────────────────────────────────────────────────────────────────────────
// INTENTIONAL OVER-PERMISSIVENESS (backend narrows these)
//
//   • struct/class/enum declarations are allowed inside function bodies.
//     The spec only prohibits nesting them inside each other, not in functions.
//   • asm() is accepted anywhere; backend rejects it outside intrinsics packages.
//   • Struct literals are accepted in if/for/switch conditions; backend rejects.
//   • typeAliasDecl is accepted inside blocks; backend restricts to package scope.
//   • Assignments and compound assignments share exprOrAssignStmt; the backend
//     validates that the left-hand expression is an addressable lvalue.
//   • reinterpret<T>(expr) accepts any expr as its argument; backend validates
//     that T is a pointer type and expr is a pointer or addressable value.
// ─────────────────────────────────────────────────────────────────────────────

parser grammar VertexParser;
options { tokenVocab = VertexLexer; }


// ════════════════════════════════════════════════════════════════════════════
// §45–46, §33  FILE STRUCTURE
//
// Source file layout (enforced by backend where ordering matters):
//   package declaration  →  build tag(s)  →  imports  →  declarations
//
// packageDecl is marked optional here so the parser produces a useful error
// message when it is missing, rather than a generic syntax failure.
// ════════════════════════════════════════════════════════════════════════════

file
    : packageDecl? buildDecl* importDecl* topLevelDecl* EOF
    ;

// §46  Exactly one per file; must be a valid identifier.
packageDecl
    : PACKAGE IDENTIFIER
    ;

// §45  One 'build <tag>' per line; multiple tags in a file are all required.
buildDecl
    : BUILD IDENTIFIER
    ;

// §33  Single or grouped imports. Grouped form is newline-separated (no commas).
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

// §5  Type alias — backend restricts to package scope.
typeAliasDecl
    : TYPE IDENTIFIER ASSIGN typeExpr
    ;

// ════════════════════════════════════════════════════════════════════════════
// §21, §28, §32, §36, §39–41  FUNCTION DECLARATIONS
//
// Named functions are top-level only. Anonymous functions are expressions
// (anonFuncExpr inside expr). Associated functions (§28) are top-level
// functions distinguished by a receiver argument.
//
// Shape:
//   func  [receiver]  name  [<T, U>]  ( params )  [qualifier]  [-> type]  { body }
// ════════════════════════════════════════════════════════════════════════════

funcDecl
    : FUNC receiver? IDENTIFIER genericParams?
      LPAREN paramList? RPAREN
      funcQualifier? (ARROW typeExpr)?
      block
    ;

// §28  Value receiver (p: T) copies; pointer receiver (p: *T) aliases caller.
// Backend auto-inserts & at call sites for pointer receivers.
// Receiver type must be a plain or pointer-to-named-type (backend validates).
receiver
    : LPAREN IDENTIFIER COLON typeExpr RPAREN
    ;

// §32  Unconstrained type parameters only in 2.1. 'where T:' is deferred.
genericParams
    : LT IDENTIFIER (COMMA IDENTIFIER)* GT
    ;

// §36, §39–41  Qualifier sits between the closing ')' and the return arrow.
funcQualifier
    : ASYNC
    | THREAD
    | PROCESS
    | GPU
    ;

// §21  Parameter list. A variadic parameter, if present, must be last.
// Both named functions and anonymous functions share this rule.
// classMember native method signatures also reuse it, gaining variadic
// support for free (backend validates C ABI compatibility).
paramList
    : param (COMMA param)* (COMMA variadicParam)? COMMA?
    | variadicParam COMMA?
    ;

param
    : IDENTIFIER COLON typeExpr
    ;

// §21.1  Variadic parameter — Go-style '...' prefix on the element type.
// Must be the final parameter. Inside the function body the parameter
// is iterable (for-in). Backend maps it to a C variadic on native
// class methods and to a slice-backed sequence on regular functions.
variadicParam
    : IDENTIFIER COLON ELLIPSIS typeExpr
    ;

// ════════════════════════════════════════════════════════════════════════════
// §27  STRUCT DECLARATIONS
//
// Struct bodies contain field declarations only. Methods are top-level
// associated functions with a value or pointer receiver (§28). The 'methods
// inside structs' feature is listed as removed in 2.1.
// ════════════════════════════════════════════════════════════════════════════

structDecl
    : STRUCT IDENTIFIER LBRACE structFieldDecl* RBRACE
    ;

structFieldDecl
    : IDENTIFIER COLON typeExpr
    ;

// ════════════════════════════════════════════════════════════════════════════
// §30, §44  CLASS DECLARATIONS
//
// Regular classes: field declarations only. init/deinit and all other methods
// are top-level associated functions using receiver syntax.
//
// Native-interface classes (§44): additionally contain bodyless method
// signatures. The ': qualifiedIdent' suffix names the native module/vtable.
// Inheritance is removed in 2.1; the ':' suffix is for native binding only.
// ════════════════════════════════════════════════════════════════════════════

classDecl
    : CLASS IDENTIFIER (COLON qualifiedIdent)?
      LBRACE classMember* RBRACE
    ;

// Fields and native method signatures are syntactically unambiguous:
// fields start with IDENTIFIER COLON, native methods start with FUNC.
// Native method signatures reuse paramList, so variadic params are
// supported automatically (e.g. func printf(fmt: ...*const char)).
classMember
    : IDENTIFIER COLON typeExpr                                       // field
    | FUNC IDENTIFIER LPAREN paramList? RPAREN (ARROW typeExpr)?      // native method (no body)
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

// Raw values are integer or string literals (§29). Integer values may be
// negative. Backend validates the raw value type matches the enum's declared
// type and that int raw values are sequential / non-duplicate.
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
    : varDecl               // §2  let / var binding (also weak let §30.1)
    | ifStmt                // §18
    | whileStmt             // §22
    | forInStmt             // §23
    | switchStmt            // §19
    | returnStmt            // §21
    | BREAK                 // §20 — backend: must be inside loop or switch
    | CONTINUE              // §20 — backend: must be inside loop, not switch
    | FALLTHROUGH           // §19 — backend: must be inside switch case
    | deferStmt             // §31 — backend: must be inside function body
    | structDecl            // spec only prohibits nesting inside struct/class
    | classDecl             // spec only prohibits nesting inside class/struct
    | enumDecl              // spec only prohibits nesting inside struct/class
    | typeAliasDecl         // backend: rejects if not at package scope
    | exprOrAssignStmt      // expression statement, assignment, or compound assignment
    ;

// ── §2, §26, §30.1  Variable declaration ─────────────────────────────────────
//
//   let x = 10
//   var y: int32 = 20
//   let a: Animal? = nil
//   let (lo, hi) = minMax(values: nums)    — tuple destructuring
//   weak let b = a                         — non-owning ref (backend validates)
//
// Backend rules:
//   • WEAK is only valid with LET, not VAR.
//   • WEAK is only valid when the RHS is a .new() ref-counted instance.
//   • Type annotation, if present, must be compatible with the inferred type.
varDecl
    : WEAK? (LET | VAR) bindingPattern (COLON typeExpr)? ASSIGN expr
    ;

bindingPattern
    : IDENTIFIER                                         // simple: let x = …
    | LPAREN IDENTIFIER (COMMA IDENTIFIER)+ RPAREN       // tuple: let (a, b) = …
    ;

// ── §18  If / else if / else ──────────────────────────────────────────────────
ifStmt
    : IF ifCondition block (ELSE (block | ifStmt))?
    ;

// if-let safely unwraps T? or binds the Ok value of a Result(T,E).
// Backend resolves which semantics apply based on the type of the RHS expr.
ifCondition
    : LET IDENTIFIER ASSIGN expr   // if-let optional/Result binding
    | expr                          // boolean expression
    ;

// ── §22  While ────────────────────────────────────────────────────────────────
whileStmt
    : WHILE expr block
    ;

// ── §23  For-in ───────────────────────────────────────────────────────────────
// The loop variable is a single identifier. Tuple for-loop destructuring
// ('for (a, b) in pairs') is listed as deferred in 2.1.
forInStmt
    : FOR IDENTIFIER IN expr block
    ;

// ── §19  Switch ───────────────────────────────────────────────────────────────
switchStmt
    : SWITCH expr LBRACE switchCase* RBRACE
    ;

// Backend enforces:
//   • At least one stmt in each case body (empty body without fallthrough is an error).
//   • Exactly one fallthrough per case at most, and only as the final statement.
//   • Exhaustiveness for enum switches (default not required when all cases covered).
switchCase
    : CASE switchPattern (COMMA switchPattern)* COLON stmt+
    | DEFAULT COLON stmt+
    ;

// Pattern forms:
//   expr             — literal or qualified enum value (case 0, case Direction.north)
//   DOT IDENTIFIER   — enum shorthand (case .north)
//   Ok(let id)       — Result Ok binding; id is bound as the success value
//   Err(let id)      — Result Err binding; id is bound as the error value
switchPattern
    : expr
    | DOT IDENTIFIER
    | (OK | ERR_KW) LPAREN LET IDENTIFIER RPAREN
    ;

// ── §21  Return ───────────────────────────────────────────────────────────────
// Void functions omit the expression. The parser is greedy: 'return' followed
// by a parseable expression always consumes it. Use an explicit closing '}' or
// the next keyword to terminate a void return.
returnStmt
    : RETURN expr?
    ;

// ── §31  Defer ────────────────────────────────────────────────────────────────
// Grammar accepts any expression; backend requires it to be a direct call.
// The anonymous-function IIFE form 'defer func() { … }()' is a call
// expression and parses naturally without special-casing.
// Multiple defers run LIFO. Backend rejects defer at file scope.
deferStmt
    : DEFER expr
    ;

// ── Expression and assignment statements ──────────────────────────────────────
// Assignments are statements, not expressions — they have no value and cannot
// be nested inside expressions (no 'a = (b = c)' style).
// The LHS is parsed as a full expr for uniformity; backend validates it is
// an addressable lvalue (variable, field access, or subscript).
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
// ANTLR4 left-recursive rule rewriting encodes precedence via alternative
// ordering: alternative 1 = highest precedence.
//
// Left-recursive alternatives (starting with 'expr') form postfix and binary
// operators. Non-left-recursive alternatives are prefix operators or primary
// expressions.
//
// ── Struct-literal vs. block ambiguity (§27) ─────────────────────────────────
// 'TypeName { field: val, … }' looks like an identifier followed by a block.
// The IDENTIFIER LBRACE alternative for struct literals is listed BEFORE the
// bare IDENTIFIER alternative. ANTLR4's LL(*) lookahead resolves the choice
// by examining the first token(s) inside '{':
//
//   IDENTIFIER COLON   → struct literal field  → parse as struct literal
//   keyword / other    → statement             → parse as bare IDENTIFIER,
//                                                let enclosing stmt rule
//                                                consume the block
//
// An empty '{}' on an identifier is parsed as a zero-field struct literal.
// Backend rejects it if the struct type has required fields.
//
// ── reinterpret<T>(expr) ambiguity note ──────────────────────────────────────
// Because REINTERPRET is a dedicated keyword, the parser knows that any LT
// immediately following it opens a type-argument list, not a comparison.
// No semantic predicate or lookahead conflict arises.
// ════════════════════════════════════════════════════════════════════════════

expr
    // ── Postfix (highest precedence) ─────────────────────────────────────────
    // Method call is listed before field access. When the parser sees
    // 'expr DOT IDENTIFIER', LL(*) lookahead checks for LPAREN:
    //   LPAREN follows  → method call wins (longer alternative)
    //   anything else   → field access
    //
    // This uniform rule covers all postfix method calls including the
    // execution postfixes (.await, .spawn, .fork, .dispatch, .new, .delete,
    // .try) and type-level intrinsics (.channel) — they are all just method
    // calls whose names happen to be reserved by the backend.
    : expr DOT IDENTIFIER LPAREN argList? RPAREN   // method call:  a.method(args)
    | expr DOT IDENTIFIER                           // field / property: a.field
    | expr LBRACKET expr RBRACKET                   // subscript:    a[i], map["key"]
    | expr LPAREN argList? RPAREN                   // call:         f(args)

    // ── Binary operators (high → low per §17) ────────────────────────────────
    | expr (LSHIFT | RSHIFT) expr                            // §9  shift
    | expr (STAR | SLASH | PERCENT | OVERFLOW_MUL) expr      // §7, §10 multiplicative
    | expr (PLUS | MINUS | OVERFLOW_ADD | OVERFLOW_SUB) expr // §7, §10 additive
    | expr (ELLIPSIS | HALF_OPEN) expr                       // §13 range
    | <assoc=right> expr NIL_COALESCE expr                   // §15 ?? right-assoc
    | expr (EQ | NEQ | LT | GT | LEQ | GEQ
           | IDENTITY_EQ | IDENTITY_NEQ) expr                // §11, §16 comparison/identity
    | expr LOGICAL_AND expr                                  // §12 &&
    | expr LOGICAL_OR  expr                                  // §12 ||
    | <assoc=right> expr QUESTION expr COLON expr            // §14 ternary right-assoc

    // ── Prefix / unary ────────────────────────────────────────────────────────
    // These are non-left-recursive, so ANTLR4 effectively gives them higher
    // precedence than any binary alternative above.
    | (MINUS | BANG | TILDE) expr   // §7 unary -, §12 !, §9 ~
    | AMP expr                       // §4 address-of &x; backend validates lvalue

    // ── Primary expressions ───────────────────────────────────────────────────
    // Struct literal before bare IDENTIFIER (see ambiguity note above).
    | IDENTIFIER LBRACE structLiteralFields? RBRACE                  // §27 struct literal
    | LBRACE mapLiteralFields? RBRACE                                // §25 map literal {"k": v}
    | MAP LBRACKET typeExpr RBRACKET typeExpr LPAREN argList? RPAREN // §25 empty map allocation
    | literal                                                        // §1  all literal forms
    | IDENTIFIER                                                     // variable, type name
    | DOT IDENTIFIER                                                 // §29 enum shorthand .caseName
    | LPAREN expr RPAREN                                             // grouping (not a tuple)
    | LPAREN expr (COMMA expr)+ RPAREN                               // §37 tuple literal (a, b, …)
    | LPAREN RPAREN                                                  // §37 empty tuple / void ()
    | LBRACKET (expr (COMMA expr)* COMMA?)? RBRACKET                 // §24 array literal [a, b, …]
    | anonFuncExpr                                                   // §35 anonymous function
    | RESULT LPAREN (OK | ERR_KW) COMMA expr RPAREN                  // §38.3 Result(Ok,v)/Result(Err,e)
    | REINTERPRET LT typeExpr GT LPAREN expr RPAREN                  // raw pointer reinterpretation
    | asmExpr                                                        // §47 inline assembly
    ;

// ── Argument list (§21) ───────────────────────────────────────────────────────
// Supports both positional and labeled (keyword-argument) styles.
// Labels are optional per argument; backend validates call-site conventions.
// Trailing comma is accepted for multiline call formatting.
argList
    : arg (COMMA arg)* COMMA?
    ;

arg
    : (IDENTIFIER COLON)? expr
    ;

// ── Struct literal fields (§27) ───────────────────────────────────────────────
// All field labels are required (positional initialization is not supported).
// Trailing comma is accepted for multiline formatting.
// Backend validates all fields of the type are provided.
structLiteralFields
    : structLiteralField (COMMA structLiteralField)* COMMA?
    ;

structLiteralField
    : IDENTIFIER COLON expr
    ;

// ── Map literal fields (§25) ──────────────────────────────────────────────────
// Trailing comma is accepted for multiline formatting.
mapLiteralFields
    : mapLiteralField (COMMA mapLiteralField)* COMMA?
    ;

mapLiteralField
    : expr COLON expr
    ;

// ── §35, §35.1  Anonymous function expression ─────────────────────────────────
// Identical to a named funcDecl minus the name and receiver.
// The execution qualifier sits between ')' and '->' — same position as named
// functions, introducing no new grammar rule (§35.1 states this explicitly).
//
// The IIFE pattern 'func(params) qualifier -> T { body }(args).postfix()' is
// parsed naturally: anonFuncExpr is a primary expr, followed by a call postfix
// '(args)', followed by further method-call postfixes.
anonFuncExpr
    : FUNC LPAREN paramList? RPAREN funcQualifier? (ARROW typeExpr)? block
    ;

// ── §47  Inline assembly ──────────────────────────────────────────────────────
// Backend rejects asm() outside a 'build intrinsics' function body.
// A void asm has no out/inout constraints. A returning asm maps out/inout
// constraints to the return type in declaration order.
// 'in' reuses the IN keyword (for-in loop keyword, same token).
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

// in("reg") param    — register loaded with param on entry, not an output
// inout("reg") param — register seeded on entry; exit value contributes to return
// out("reg")         — exit value contributes to return; undefined on entry
asmConstraint
    : IN     LPAREN STRING_LIT RPAREN IDENTIFIER   // in constraint
    | INOUT  LPAREN STRING_LIT RPAREN IDENTIFIER   // inout constraint
    | OUT_KW LPAREN STRING_LIT RPAREN              // out constraint (no param)
    ;

asmClobberDecl
    : CLOBBER LPAREN STRING_LIT (COMMA STRING_LIT)* RPAREN
    ;

// ════════════════════════════════════════════════════════════════════════════
// TYPE EXPRESSIONS  §3–§4, §24, §34, §37, §38.3, §42
//
// All scalar type names (int, int32, float, float64, bool, char, string, void,
// etc.) are plain identifiers matched by baseType. They are not keywords.
//
// Pointer vs. scalar optional (§4, §26):
//   *T   starts with STAR   → pointer rule handles it  (including *T? nullable)
//   T?   starts with IDENT  → baseType QUESTION?       (scalar optional)
// These are syntactically disjoint, so no ambiguity arises.
//
// This rule is right-recursive (recurses on the right), which ANTLR4 handles
// without rewriting.
// ════════════════════════════════════════════════════════════════════════════

typeExpr
    : STAR CONST_KW? typeExpr QUESTION?                         // §4  *T, *const T, *T?, *const T?
    | LBRACKET typeExpr RBRACKET                                // §24 [T] array
    | MAP LBRACKET typeExpr RBRACKET typeExpr                   // §25 map[K]V map
    | CHAN typeExpr                                             // §42 chan T channel
    | FUNC LPAREN funcTypeParams? RPAREN (ARROW typeExpr)?      // §34 func(T…)->T function type
    | LPAREN tupleTypeElems? RPAREN                             // §37 (T,T) tuple / () void
    | RESULT LPAREN typeExpr COMMA typeExpr RPAREN              // §38.3 Result(T, E)
    | baseType QUESTION?                                        // named type T, or optional T?
    ;

// §34  Function type parameters are types only — no parameter names.
funcTypeParams
    : typeExpr (COMMA typeExpr)* COMMA?
    ;

// Named type, possibly qualified (e.g. 'rtp.Packet', 'canvas.Context').
baseType
    : IDENTIFIER (DOT IDENTIFIER)*
    ;

// §37  Tuple type elements. Labels are optional per element.
// 'Mixed tuple element labels' is deferred in 2.1; backend validates
// that all elements either all have labels or all are unlabeled.
tupleTypeElems
    : tupleTypeElem (COMMA tupleTypeElem)* COMMA?
    ;

tupleTypeElem
    : (IDENTIFIER COLON)? typeExpr
    ;


// ════════════════════════════════════════════════════════════════════════════
// §1  LITERALS
// ════════════════════════════════════════════════════════════════════════════
// Negative numeric literals do not exist as tokens. Unary MINUS in the expr
// rule produces negative values at compile time.
// char and string share STRING_LIT; distinction is resolved by the backend
// based on the declared type of the binding.
literal
    : DEC_INT_LIT
    | HEX_INT_LIT
    | OCT_INT_LIT
    | BIN_INT_LIT
    | DEC_FLOAT_LIT
    | HEX_FLOAT_LIT
    | STRING_LIT
    | MULTILINE_STRING_LIT
    | TRUE
    | FALSE
    | NIL
    ;