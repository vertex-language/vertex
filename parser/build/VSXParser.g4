// VSXParser.g4
// Vertex Source XML (VSX) — Parser
// Standalone. No Vertex grammar dependency.
//
// vsxExpr covers every expression form React's AssignmentExpression
// slot allows inside {} — confirmed from the JSX spec and React docs:
//
//   variables / field access        {state.status}
//   arithmetic                      {count + 1}
//   string concat                   {"Hello " + name}
//   comparison                      {count > 0}
//   logical &&                      {isOpen && <Panel />}
//   logical ||                      {val || "default"}
//   ternary                         {ok ? <A /> : <B />}
//   function call                   {formatDate(today)}
//   method call + chaining          {state.messages.map(...)}
//   anonymous func (Vertex form)    {func(msg) { ... }}
//   array literal                   {[1, 2, 3]}
//   object literal                  {{key: val}}
//   new expression                  {new Foo()}
//   spread child                    {...items}
//   nested element / fragment       {<Row />}  {<></>}
//   null / true / false             {null}  {true}
//   numeric / string literal        {42}  {"hello"}
//   unary ! -                       {!isOpen}  {-count}
//   grouped                         {(a + b) * c}

parser grammar VSXParser;
options { tokenVocab = VSXLexer; }


// ═══════════════════════════════════════════════════════════════════
// Top
// ═══════════════════════════════════════════════════════════════════

vsxPrimary
    : vsxFragment
    | vsxElement
    ;


// ═══════════════════════════════════════════════════════════════════
// Fragment  —  < > children? < / >
// ═══════════════════════════════════════════════════════════════════

vsxFragment
    : VSX_LT VSX_GT
      vsxChildren?
      VSX_LT VSX_SLASH VSX_GT
    ;


// ═══════════════════════════════════════════════════════════════════
// Element
// ═══════════════════════════════════════════════════════════════════

vsxElement
    : vsxSelfClosingElement
    | vsxOpeningElement vsxChildren? vsxClosingElement
    ;

vsxSelfClosingElement
    : VSX_LT vsxElementName vsxAttributes? VSX_SLASH VSX_GT
    ;

vsxOpeningElement
    : VSX_LT vsxElementName vsxAttributes? VSX_GT
    ;

vsxClosingElement
    : VSX_LT VSX_SLASH vsxElementName VSX_GT
    ;


// ═══════════════════════════════════════════════════════════════════
// Element name  —  Name | NS:Name | Member.Expr
// ═══════════════════════════════════════════════════════════════════

vsxElementName
    : vsxMemberExpression
    | vsxNamespacedName
    | vsxIdentifier
    ;

vsxIdentifier
    : VSX_IDENT
    ;

vsxNamespacedName
    : vsxIdentifier VSX_COLON vsxIdentifier
    ;

vsxMemberExpression
    : vsxIdentifier (VSX_DOT vsxIdentifier)+
    ;


// ═══════════════════════════════════════════════════════════════════
// Attributes
// ═══════════════════════════════════════════════════════════════════

vsxAttributes
    : vsxAttributeItem+
    ;

vsxAttributeItem
    : vsxSpreadAttribute
    | vsxAttribute
    ;

vsxSpreadAttribute
    : VSX_LBRACE VSX_SPREAD vsxExpr VSX_RBRACE
    ;

vsxAttribute
    : vsxAttributeName vsxAttributeInitializer?
    ;

vsxAttributeName
    : vsxNamespacedName
    | vsxIdentifier
    ;

vsxAttributeInitializer
    : VSX_ASSIGN vsxAttributeValue
    ;

vsxAttributeValue
    : VSX_DOUBLE_STRING
    | VSX_SINGLE_STRING
    | VSX_LBRACE vsxExpr VSX_RBRACE
    | vsxElement
    | vsxFragment
    ;


// ═══════════════════════════════════════════════════════════════════
// Children
// ═══════════════════════════════════════════════════════════════════

vsxChildren
    : vsxChild+
    ;

vsxChild
    : VSX_TEXT
    | vsxElement
    | vsxFragment
    | VSX_LBRACE vsxChildExpression? VSX_RBRACE
    ;

vsxChildExpression
    : VSX_SPREAD vsxExpr       // {...items}
    | vsxExpr                  // {value}
    ;


// ═══════════════════════════════════════════════════════════════════
// vsxExpr — full AssignmentExpression-equivalent
//
// Precedence from lowest to highest (ANTLR resolves left-recursion):
//
//   ternary         a ? b : c
//   logical or      a || b
//   logical and     a && b
//   equality        a == b   a != b
//   relational      a < b    a > b   a <= b   a >= b
//   additive        a + b    a - b
//   multiplicative  a * b    a / b   a % b
//   unary           !a   -a
//   postfix call    a(args)  a[i]  a.b
//   primary         literal  ident  anon-func  array  object
//                   nested element/fragment  grouped
// ═══════════════════════════════════════════════════════════════════

vsxExpr
    // ternary — right-associative
    : <assoc=right> vsxExpr VSX_QUESTION vsxExpr VSX_COLON vsxExpr

    // logical
    | vsxExpr VSX_OR  vsxExpr
    | vsxExpr VSX_AND vsxExpr

    // equality
    | vsxExpr VSX_EQ  vsxExpr
    | vsxExpr VSX_NEQ vsxExpr

    // relational  (use tokens not raw < > to avoid tag ambiguity)
    | vsxExpr VSX_LT  vsxExpr
    | vsxExpr VSX_GT  vsxExpr
    | vsxExpr VSX_LTE vsxExpr
    | vsxExpr VSX_GTE vsxExpr

    // additive
    | vsxExpr VSX_PLUS  vsxExpr
    | vsxExpr VSX_MINUS vsxExpr

    // multiplicative
    | vsxExpr VSX_STAR    vsxExpr
    | vsxExpr VSX_SLASH   vsxExpr
    | vsxExpr VSX_PERCENT vsxExpr

    // unary prefix
    | VSX_BANG  vsxExpr
    | VSX_MINUS vsxExpr

    // postfix — method call, function call, index, field access
    | vsxExpr VSX_DOT VSX_IDENT VSX_LPAREN vsxArgList? VSX_RPAREN  // a.method(...)
    | vsxExpr VSX_DOT VSX_IDENT                                     // a.field
    | vsxExpr VSX_LPAREN vsxArgList? VSX_RPAREN                     // call(...)
    | vsxExpr VSX_LBRACKET vsxExpr VSX_RBRACKET                     // a[i]

    // new expression
    | VSX_NEW vsxExpr VSX_LPAREN vsxArgList? VSX_RPAREN

    // primary
    | vsxPrimaryExpr
    ;


vsxPrimaryExpr
    : vsxLiteral                                               // 42  "hi"  true  null
    | VSX_IDENT                                                // variable
    | vsxAnonFunc                                              // func(x) { ... }
    | vsxArrayLiteral                                          // [a, b, c]
    | vsxObjectLiteral                                         // {key: val}
    | vsxElement                                               // <Comp />
    | vsxFragment                                              // <></>
    | VSX_LPAREN vsxExpr VSX_RPAREN                            // (grouped)
    ;


// ─── Literals ───────────────────────────────────────────────────────

vsxLiteral
    : VSX_NUMBER
    | VSX_DOUBLE_STRING
    | VSX_SINGLE_STRING
    | VSX_BACKTICK_STRING
    | VSX_TRUE
    | VSX_FALSE
    | VSX_NULL
    ;


// ─── Anonymous function — Vertex form: func(params) { body } ────────
// Body is an opaque block; VSX does not parse Vertex statements.
// The host compiler owns the block contents.

vsxAnonFunc
    : VSX_FUNC VSX_LPAREN vsxParamList? VSX_RPAREN vsxOpaqueBlock
    ;

vsxParamList
    : VSX_IDENT (VSX_COMMA VSX_IDENT)*
    ;

// Opaque block — balanced braces, any content.
// VSX tracks depth only; Vertex parser owns what's inside.
vsxOpaqueBlock
    : VSX_LBRACE vsxOpaqueContent* VSX_RBRACE
    ;

vsxOpaqueContent
    : vsxOpaqueBlock                  // nested braces stay balanced
    | vsxElement                      // nested VSX elements are still parsed
    | vsxFragment
    | ~(VSX_LBRACE | VSX_RBRACE)      // any non-brace token passes through
    ;


// ─── Array literal ──────────────────────────────────────────────────

vsxArrayLiteral
    : VSX_LBRACKET (vsxExpr (VSX_COMMA vsxExpr)* VSX_COMMA?)? VSX_RBRACKET
    ;


// ─── Object literal — {{ key: val }} in JSX double-curly style ──────

vsxObjectLiteral
    : VSX_LBRACE (vsxObjectEntry (VSX_COMMA vsxObjectEntry)* VSX_COMMA?)? VSX_RBRACE
    ;

vsxObjectEntry
    : vsxObjectKey VSX_COLON vsxExpr
    | VSX_SPREAD vsxExpr
    ;

vsxObjectKey
    : VSX_IDENT
    | VSX_DOUBLE_STRING
    | VSX_SINGLE_STRING
    | VSX_LBRACKET vsxExpr VSX_RBRACKET    // computed key [expr]: val
    ;


// ─── Argument list ──────────────────────────────────────────────────

vsxArgList
    : vsxArg (VSX_COMMA vsxArg)*
    ;

vsxArg
    : VSX_SPREAD vsxExpr    // spread arg  f(...arr)
    | vsxExpr
    ;