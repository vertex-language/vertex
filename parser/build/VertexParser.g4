// ============================================================
// VertexParser.g4  –  Vertex 6 Syntactic Grammar (ANTLR4)
//
// Depends on: VertexLexer.g4
// ============================================================

parser grammar VertexParser;
options { 
    tokenVocab = VertexLexer; 
    superClass = VertexParserBase;
}

// ─────────────────────────────────────────────────────────────
// TOP LEVEL
// ─────────────────────────────────────────────────────────────
topLevel
    : statements? EOF
    ;

// ─────────────────────────────────────────────────────────────
// STATEMENTS
// ─────────────────────────────────────────────────────────────
statements
    : statement+
    ;

statement
    : expression         SEMICOLON?   # expressionStatement
    | declaration        SEMICOLON?   # declarationStatement
    | loopStatement      SEMICOLON?   # loopStmt
    | branchStatement    SEMICOLON?   # branchStmt
    | labeledStatement   SEMICOLON?   # labeledStmt
    | controlTransfer    SEMICOLON?   # controlTransferStmt
    | deferStatement     SEMICOLON?   # deferStmt
    | doStatement        SEMICOLON?   # doStmt
    | compilerControl    SEMICOLON?   # compilerControlStmt
    | macroDeclaration   SEMICOLON?   # macroDeclStmt
    ;

// ─── Loop statements ──────────────────────────────────────────
loopStatement
    : forInStatement
    | whileStatement
    | repeatWhileStatement
    | forStatement
    ;

forInStatement
    : FOR CASE? pattern IN expression whereClause? codeBlock
    ;

whileStatement
    : WHILE conditionList codeBlock
    ;

conditionList
    : condition (COMMA condition)*
    ;

condition
    : expression
    | availabilityCondition
    | caseCondition
    | optionalBindingCondition
    ;

caseCondition
    : CASE pattern initializer
    ;

optionalBindingCondition
    : (LET | VAR) pattern initializer
    ;

repeatWhileStatement
    : REPEAT codeBlock WHILE expression
    ;

forStatement
    : FOR LPAREN expression? SEMICOLON expression? SEMICOLON expression? RPAREN codeBlock
    ;

// ─── Branch statements ────────────────────────────────────────
branchStatement
    : ifStatement
    | guardStatement
    | switchStatement
    | tryStatement
    ;

ifStatement
    : IF conditionList codeBlock elseClause?
    ;

elseClause
    : ELSE codeBlock
    | ELSE ifStatement
    ;

guardStatement
    : GUARD conditionList ELSE codeBlock
    ;

switchStatement
    : SWITCH expression LBRACE switchCases? RBRACE
    ;

switchCases
    : switchCase+
    ;

switchCase
    : caseLabel statements
    | defaultLabel statements
    ;

caseLabel
    : CASE caseItemList COLON
    ;

caseItemList
    : caseItem (COMMA caseItem)*
    ;

caseItem
    : pattern whereClause?
    ;

defaultLabel
    : DEFAULT COLON
    ;

whereClause
    : WHERE expression
    ;

// ─── Labeled statement ────────────────────────────────────────
labeledStatement
    : labelName COLON (loopStatement | ifStatement | switchStatement | doStatement)
    ;

labelName : identifier ;

// ─── Control transfer ─────────────────────────────────────────
controlTransfer
    : BREAK    labelName?                                         # breakStatement
    | CONTINUE labelName?                                         # continueStatement
    | FALLTHROUGH                                                 # fallthroughStatement
    | RETURN   expression?                                        # returnStatement
    | THROW    expression                                         # throwStatement
    ;

// ─── Defer ────────────────────────────────────────────────────
deferStatement
    : DEFER codeBlock
    ;

// ─── Do / catch ───────────────────────────────────────────────
doStatement
    : DO codeBlock catchClause*
    ;

catchClause
    : CATCH catchPatternList? codeBlock
    ;

catchPatternList
    : catchPattern (COMMA catchPattern)*
    ;

catchPattern
    : pattern whereClause?
    ;

tryStatement
    : TRY codeBlock catchClause*
    ;

// ─── Compiler control ─────────────────────────────────────────
compilerControl
    : conditionalCompilationBlock
    | lineControlStatement
    | diagnosticStatement
    ;

conditionalCompilationBlock
    : POUND_IF compilationCondition statements?
      (POUND_ELSEIF compilationCondition statements?)*
      (POUND_ELSE statements?)?
      POUND_ENDIF
    ;

compilationCondition
    : platformCondition
    | identifier
    | booleanLiteral
    | LPAREN compilationCondition RPAREN
    | OPERATOR compilationCondition
    | compilationCondition OPERATOR compilationCondition
    ;

platformCondition
    : OS_KW     LPAREN identifier RPAREN
    | ARCH_KW  LPAREN identifier RPAREN
    | VERTEX_KW LPAREN OPERATOR decimalVersion RPAREN
    | COMPILER_KW LPAREN OPERATOR decimalVersion RPAREN
    | CANIMPORT_KW LPAREN identifier (COMMA VERSION_KW COLON decimalVersion)? RPAREN
    ;

decimalVersion
    : INTEGER_LITERAL (DOT INTEGER_LITERAL (DOT INTEGER_LITERAL)?)?
    ;

lineControlStatement
    : POUND_SOURCE_LOCATION
      LPAREN
        FILE_KEYWORD COLON STRING_LITERAL COMMA
        LINE_KEYWORD COLON INTEGER_LITERAL
      RPAREN
    | POUND_SOURCE_LOCATION LPAREN RPAREN
    ;

diagnosticStatement
    : OPERATOR STRING_LITERAL
    ;

// ─────────────────────────────────────────────────────────────
// CODE BLOCK
// ─────────────────────────────────────────────────────────────
codeBlock
    : LBRACE statements? RBRACE
    ;

// ─────────────────────────────────────────────────────────────
// DECLARATIONS
// ─────────────────────────────────────────────────────────────
declaration
    : importDeclaration
    | constantDeclaration
    | variableDeclaration
    | typealiasDeclaration
    | functionDeclaration
    | enumDeclaration
    | structDeclaration
    | classDeclaration
    | actorDeclaration
    | protocolDeclaration
    | initializerDeclaration
    | deinitializerDeclaration
    | extensionDeclaration
    | subscriptDeclaration
    ;

// ─── Import ───────────────────────────────────────────────────
importDeclaration
    : attributes? IMPORT importSpec
    | attributes? IMPORT LPAREN importSpec* RPAREN
    ;

importSpec
    : importAlias? STRING_LITERAL
    ;

importAlias
    : identifier
    | DOT
    | UNDERSCORE
    ;

// ─── Constants ────────────────────────────────────────────────
constantDeclaration
    : attributes? declarationModifiers? LET patternInitializerList
    ;

patternInitializerList
    : patternInitializer (COMMA patternInitializer)*
    ;

patternInitializer
    : pattern initializer?
    ;

initializer
    : ASSIGN expression
    ;

// ─── Variables ────────────────────────────────────────────────
variableDeclaration
    : variableDeclarationHead patternInitializerList
    | variableDeclarationHead variableName typeAnnotation codeBlock
    | variableDeclarationHead variableName typeAnnotation getterSetterBlock
    | variableDeclarationHead variableName typeAnnotation getterSetterKeywordBlock
    | variableDeclarationHead variableName typeAnnotation initializer? willSetDidSetBlock
    ;

variableDeclarationHead
    : attributes? declarationModifiers? VAR
    ;

variableName
    : identifier
    ;

getterSetterBlock
    : LBRACE getterClause setterClause? RBRACE
    | LBRACE setterClause getterClause RBRACE
    ;

getterClause
    : attributes? mutationModifier? GET_KW codeBlock
    ;

setterClause
    : attributes? mutationModifier? SET_KW setterName? codeBlock
    ;

setterName
    : LPAREN identifier RPAREN
    ;

getterSetterKeywordBlock
    : LBRACE getterKeywordClause setterKeywordClause? RBRACE
    | LBRACE setterKeywordClause getterKeywordClause RBRACE
    ;

getterKeywordClause
    : attributes? mutationModifier? GET_KW
    ;

setterKeywordClause
    : attributes? mutationModifier? SET_KW
    ;

willSetDidSetBlock
    : LBRACE willSetClause didSetClause? RBRACE
    | LBRACE didSetClause willSetClause? RBRACE
    ;

willSetClause
    : attributes? WILLSET_KW setterName? codeBlock
    ;

didSetClause
    : attributes? DIDSET_KW setterName? codeBlock
    ;

// ─── Typealias ────────────────────────────────────────────────
typealiasDeclaration
    : attributes? accessLevelModifier? TYPEALIAS identifier
      genericParameterClause? ASSIGN type
    ;

// ─── Functions ────────────────────────────────────────────────
functionDeclaration
    : functionHead identifier genericParameterClause?
      functionSignature
      genericWhereClause?
      functionBody?
    ;

functionHead
    : attributes? declarationModifiers? FUNC
    ;

functionSignature
    : parameterClause asyncModifier? THROWS? functionResult?
    ;

asyncModifier
    : ASYNC_KW
    ;

functionResult
    : ARROW attributes? type
    ;

parameterClause
    : LPAREN RPAREN
    | LPAREN parameterList RPAREN
    ;

parameterList
    : parameter (COMMA parameter)*
    ;

parameter
    : externalParameterName? localParameterName typeAnnotation defaultArgumentClause?
    | externalParameterName? localParameterName typeAnnotation ELLIPSIS
    ;

externalParameterName
    : identifier
    | UNDERSCORE
    ;

localParameterName
    : identifier
    | UNDERSCORE
    ;

defaultArgumentClause
    : ASSIGN expression
    ;

functionBody
    : codeBlock
    ;

// ─── Enumerations ─────────────────────────────────────────────
enumDeclaration
    : attributes? accessLevelModifier? ENUM identifier
      genericParameterClause?
      typeInheritanceClause?
      genericWhereClause?
      LBRACE enumMembers RBRACE
    ;

enumMembers
    : enumMember*
    ;

enumMember
    : declaration
    | unionStyleEnumCaseClause
    | rawValueStyleEnumCaseClause
    | compilerControl
    ;

unionStyleEnumCaseClause
    : attributes? CASE unionStyleEnumCaseList
    ;

unionStyleEnumCaseList
    : unionStyleEnumCase (COMMA unionStyleEnumCase)*
    ;

unionStyleEnumCase
    : identifier tupleType?
    ;

rawValueStyleEnumCaseClause
    : attributes? CASE rawValueStyleEnumCaseList
    ;

rawValueStyleEnumCaseList
    : rawValueStyleEnumCase (COMMA rawValueStyleEnumCase)*
    ;

rawValueStyleEnumCase
    : identifier rawValueAssignment?
    ;

rawValueAssignment
    : ASSIGN rawValueLiteral
    ;

rawValueLiteral
    : numericLiteral
    | STRING_LITERAL
    | booleanLiteral
    ;

// ─── Structs ──────────────────────────────────────────────────
structDeclaration
    : attributes? accessLevelModifier? STRUCT identifier
      genericParameterClause?
      typeInheritanceClause?
      genericWhereClause? structBody
    ;

structBody
    : LBRACE structMember* RBRACE
    ;

structMember
    : declaration
    | compilerControl
    ;

// ─── Classes ──────────────────────────────────────────────────
classDeclaration
    : attributes? accessLevelModifier? FINAL_KW? CLASS identifier
      genericParameterClause?
      typeInheritanceClause? genericWhereClause?
      classBody
    ;

classBody
    : LBRACE classMember* RBRACE
    ;

classMember
    : declaration
    | compilerControl
    ;

// ─── Actors ───────────────────────────────────────────────────
actorDeclaration
    : attributes? accessLevelModifier? ACTOR_KW identifier
      genericParameterClause?
      typeInheritanceClause?
      genericWhereClause? actorBody
    ;

actorBody
    : LBRACE actorMember* RBRACE
    ;

actorMember
    : declaration
    | compilerControl
    ;

// ─── Protocols ────────────────────────────────────────────────
protocolDeclaration
    : attributes? accessLevelModifier? PROTOCOL identifier
      primaryAssociatedTypeClause?
      typeInheritanceClause?
      genericWhereClause? protocolBody
    ;

primaryAssociatedTypeClause
    : LT primaryAssociatedTypeList GT
    ;

primaryAssociatedTypeList
    : identifier (COMMA identifier)*
    ;

protocolBody
    : LBRACE protocolMember* RBRACE
    ;

protocolMember
    : protocolMemberDeclaration
    | compilerControl
    ;

protocolMemberDeclaration
    : protocolPropertyDeclaration
    | protocolMethodDeclaration
    | protocolInitializerDeclaration
    | protocolSubscriptDeclaration
    | protocolAssociatedTypeDeclaration
    | typealiasDeclaration
    ;

protocolPropertyDeclaration
    : variableDeclarationHead variableName typeAnnotation getterSetterKeywordBlock
    ;

protocolMethodDeclaration
    : functionHead identifier genericParameterClause?
      functionSignature
      genericWhereClause?
    ;

protocolInitializerDeclaration
    : attributes? declarationModifiers? INIT QUESTION_POSTFIX?
      genericParameterClause?
      parameterClause
      THROWS?
      genericWhereClause?
    ;

protocolSubscriptDeclaration
    : subscriptHead subscriptResult genericWhereClause? getterSetterKeywordBlock
    ;

protocolAssociatedTypeDeclaration
    : attributes? accessLevelModifier? ASSOCIATEDTYPE identifier
      typeInheritanceClause?
      (ASSIGN type)?
      genericWhereClause?
    ;

// ─── Initializers ─────────────────────────────────────────────
initializerDeclaration
    : attributes? declarationModifiers?
      INIT QUESTION_POSTFIX?
      genericParameterClause?
      parameterClause
      asyncModifier? THROWS?
      genericWhereClause?
      initializerBody
    ;

initializerBody
    : codeBlock
    ;

// ─── Deinitializers ───────────────────────────────────────────
deinitializerDeclaration
    : attributes? DEINIT codeBlock
    ;

// ─── Extensions ───────────────────────────────────────────────
extensionDeclaration
    : attributes? accessLevelModifier? EXTENSION typeIdentifier
      (COLON typeInheritanceList)?
      genericWhereClause? extensionBody
    ;

extensionBody
    : LBRACE extensionMember* RBRACE
    ;

extensionMember
    : declaration
    | compilerControl
    ;

// ─── Subscripts ───────────────────────────────────────────────
subscriptDeclaration
    : subscriptHead subscriptResult
      genericWhereClause?
      ( getterSetterBlock
      | getterSetterKeywordBlock
      )
    ;

subscriptHead
    : attributes? declarationModifiers? SUBSCRIPT
      genericParameterClause?
      parameterClause
    ;

subscriptResult
    : ARROW attributes? type
    ;

// ─── Macros ───────────────────────────────────────────────────
macroDeclaration
    : attributes? declarationModifiers?
      MACRO_KW identifier genericParameterClause?
      parameterClause? (ARROW type)?
      (ASSIGN expression)?
      genericWhereClause?
    ;

// ─────────────────────────────────────────────────────────────
// DECLARATION MODIFIERS
// ─────────────────────────────────────────────────────────────
declarationModifiers
    : declarationModifier+
    ;

declarationModifier
    : accessLevelModifier
    | mutationModifier
    | CLASS | FINAL_KW
    | LAZY_KW | OPTIONAL_KW | OVERRIDE_KW
    | POSTFIX_KW | PREFIX_KW | REQUIRED_KW | STATIC
    | UNOWNED_KW | WEAK_KW
    | NONISOLATED_KW
    ;

accessLevelModifier
    : PRIVATE   (LPAREN SET_KW RPAREN)?
    | INTERNAL  (LPAREN SET_KW RPAREN)?
    | PUBLIC    (LPAREN SET_KW RPAREN)?
    | OPEN      (LPAREN SET_KW RPAREN)?
    ;

mutationModifier
    : MUTATING_KW
    | NONMUTATING_KW
    ;

// ─────────────────────────────────────────────────────────────
// EXPRESSIONS
// ─────────────────────────────────────────────────────────────
expression
    : tryOperator? awaitOperator? prefixExpression binaryExpressions?
    ;

tryOperator
    : TRY EXCLAIM_POSTFIX?
    | TRY QUESTION_POSTFIX?
    ;

awaitOperator
    : AWAIT
    ;

binaryExpressions
    : binaryExpression+
    ;

// ── KEY FIX ──────────────────────────────────────────────────
// The `{p.isBinaryOp()}?` predicate is placed at the left edge of the
// `binaryOp` alternative so ANTLR4 hoists it into the prediction expression.
// When the upcoming OPERATOR token has symmetric whitespace (both sides or
// neither), isBinaryOp() returns true and the token is consumed as a binary
// operator continuing the current expression.
// When whitespace is asymmetric
// (right-side only → prefix context), isBinaryOp() returns false, ANTLR does
// not take this alternative, binaryExpressions ends, and the operator is left
// to open a new prefixOp statement — which is the correct behaviour.
//
// Without this predicate, ANTLR4's SLL prediction resolved the ambiguity in
// the wrong direction, splitting  `a + b`  and  `x = x * y`  into two
// statements (the bug visible in the generated WAT).
// ─────────────────────────────────────────────────────────────
binaryExpression
    : {p.isBinaryOp()}? binaryOperator prefixExpression       # binaryOp
    | assignmentOperator tryOperator? expression               # assignmentExpr
    | conditionalOperator tryOperator? expression              # ternaryExpr
    | typeCastingOperator                                      # typeCastExpr
    ;

binaryOperator
    : OPERATOR
    | AMPERSAND
    | DOT OPERATOR
    | LT
    | GT
    ;

assignmentOperator
    : ASSIGN
    ;

conditionalOperator
    : QUESTION_POSTFIX expression COLON
    ;

typeCastingOperator
    : IS type
    | AS type
    | AS EXCLAIM_POSTFIX type
    | AS QUESTION_POSTFIX type
    ;

// ── KEY FIX ──────────────────────────────────────────────────
// Guard prefixOp symmetrically: only match when the operator actually has
// prefix whitespace (left WS, no right WS), or no left context at all.
// This stops  `+expr`  from matching when the `+` is binary, providing a
// second layer of defence on top of the binaryOp predicate above.
// ─────────────────────────────────────────────────────────────
prefixExpression
    : {p.isPrefixOp()}? prefixOperator postfixExpression   # prefixOp
    | inOutExpression                                       # inout
    | postfixExpression                                     # barePostfix
    ;

inOutExpression
    : AMPERSAND identifier
    ;

prefixOperator
    : OPERATOR
    ;

postfixExpression
    : primaryExpression postfixSuffix*
    ;

postfixSuffix
    : postfixOperator
    | functionCallSuffix
    | initializerSuffix
    | explicitMemberSuffix
    | postfixSelfSuffix
    | subscriptSuffix
    | forcedValueSuffix
    | optionalChainingLiteral
    ;

postfixOperator         : {p.isPostfixOp()}? OPERATOR ;
forcedValueSuffix       : EXCLAIM_POSTFIX ;
optionalChainingLiteral : QUESTION_POSTFIX ;

functionCallSuffix
    : functionCallArgumentClause trailingClosures?
    | trailingClosures
    ;

functionCallArgumentClause
    : LPAREN RPAREN
    | LPAREN functionCallArgumentList RPAREN
    ;

functionCallArgumentList
    : functionCallArgument (COMMA functionCallArgument)*
    ;

functionCallArgument
    : expression
    | identifier COLON expression
    | UNDERSCORE COLON expression
    | (identifier | UNDERSCORE) COLON operator_
    ;

operator_
    : OPERATOR
    | DOT OPERATOR
    ;

trailingClosures
    : closureExpression labeledTrailingClosure*
    ;

labeledTrailingClosure
    : identifier COLON closureExpression
    ;

initializerSuffix
    : DOT INIT (LPAREN argumentNames RPAREN)?
    ;

argumentNames
    : (identifier COLON)+
    ;

explicitMemberSuffix
    : DOT identifier genericArgumentClause?
    | DOT INTEGER_LITERAL
    ;

postfixSelfSuffix
    : DOT SELF
    ;

subscriptSuffix
    : LBRACKET functionCallArgumentList RBRACKET
    ;

// ─── Primary Expressions ──────────────────────────────────────
primaryExpression
    : identifier genericArgumentClause?   # identifierExpr
    | literalExpression                   # litExpr
    | selfExpression                      # selfExpr
    | superExpression                     # superExpr
    | closureExpression                   # closureExpr
    | parenthesizedExpression             # parenExpr
    | tupleExpression                     # tupleExpr
    | implicitMemberExpression            # implicitMemberExpr
    | wildcardExpression                  # wildcardExpr
    | macroExpansionExpression            # macroExpr
    ;

literalExpression
    : literal
    | arrayLiteral
    | dictionaryLiteral
    | poundFileExpression
    ;

poundFileExpression
    : POUND_FILE | POUND_FILEID | POUND_FILEPATH
    | POUND_LINE | POUND_COLUMN | POUND_FUNCTION
    ;

literal
    : numericLiteral
    | STRING_LITERAL
    | MULTILINE_STRING_LITERAL
    | EXTENDED_STRING_LITERAL
    | booleanLiteral
    | NIL
    ;

numericLiteral
    : OPERATOR? INTEGER_LITERAL
    | OPERATOR? FLOAT_LITERAL
    ;

booleanLiteral
    : TRUE | FALSE
    ;

arrayLiteral
    : LBRACKET RBRACKET
    | LBRACKET arrayLiteralItems RBRACKET
    ;

arrayLiteralItems
    : expression (COMMA expression)* COMMA?
    ;

dictionaryLiteral
    : LBRACKET dictionaryLiteralItems RBRACKET
    | LBRACKET COLON RBRACKET
    ;

dictionaryLiteralItems
    : dictionaryLiteralItem (COMMA dictionaryLiteralItem)* COMMA?
    ;

dictionaryLiteralItem
    : expression COLON expression
    ;

selfExpression
    : SELF
    | SELF DOT identifier
    | SELF LBRACKET functionCallArgumentList RBRACKET
    | SELF DOT INIT
    ;

superExpression
    : SUPER DOT identifier
    | SUPER LBRACKET functionCallArgumentList RBRACKET
    | SUPER DOT INIT
    ;

closureExpression
    : LBRACE closureSignature? statements? RBRACE
    ;

closureSignature
    : captureList? closureParameterClause (ASYNC_KW)? THROWS? functionResult? IN
    | captureList IN
    ;

captureList
    : LBRACKET captureListItems RBRACKET
    ;

captureListItems
    : captureListItem (COMMA captureListItem)*
    ;

captureListItem
    : captureSpecifier? identifier (ASSIGN expression)?
    ;

captureSpecifier
    : WEAK_KW
    | UNOWNED_KW
    ;

closureParameterClause
    : LPAREN RPAREN
    | LPAREN closureParameterList RPAREN
    | identifierList
    ;

closureParameterList
    : closureParameter (COMMA closureParameter)*
    ;

closureParameter
    : identifier typeAnnotation?
    | identifier typeAnnotation ELLIPSIS
    ;

identifierList
    : identifier (COMMA identifier)*
    ;

parenthesizedExpression
    : LPAREN expression RPAREN
    ;

tupleExpression
    : LPAREN RPAREN
    | LPAREN tupleElementList RPAREN
    ;

tupleElementList
    : tupleElement (COMMA tupleElement)*
    ;

tupleElement
    : expression
    | identifier COLON expression
    ;

implicitMemberExpression
    : DOT identifier genericArgumentClause?
    ;

wildcardExpression
    : UNDERSCORE
    ;

macroExpansionExpression
    : HASH identifier genericArgumentClause? functionCallArgumentClause? trailingClosures?
    ;

// ─────────────────────────────────────────────────────────────
// TYPES
// ─────────────────────────────────────────────────────────────
type
    : typeAnnotationHead type                      # annotatedType
    | functionType                                 # funcType
    | arrayType                                    # arrType
    | dictionaryType                               # dictType
    | typeIdentifier                               # namedType
    | tupleType                                    # tupType
    | protocolCompositionType                      # protoCompType
    | existentialType                              # existType
    | opaqueType                                   # opaqueType_
    | selfType                                     # selfType_
    | type QUESTION_POSTFIX                        # optType
    | type DOT (TYPE_KW | PROTOCOL_KW)             # metatypeType_
    | type ELLIPSIS                                # variadicType
    ;

typeAnnotationHead
    : attributes? INOUT
    ;

functionType
    : attributes? functionTypeArgumentClause
      asyncModifier?
      THROWS?
      ARROW type
    ;

functionTypeArgumentClause
    : LPAREN RPAREN
    | LPAREN functionTypeArgumentList ELLIPSIS? RPAREN
    ;

functionTypeArgumentList
    : functionTypeArgument (COMMA functionTypeArgument)*
    ;

functionTypeArgument
    : attributes? INOUT? type
    | argumentLabel typeAnnotation
    ;

argumentLabel
    : identifier COLON
    | UNDERSCORE COLON
    ;

arrayType
    : LBRACKET type RBRACKET
    ;

dictionaryType
    : LBRACKET type COLON type RBRACKET
    ;

typeIdentifier
    : identifier genericArgumentClause? (DOT typeIdentifier)?
    ;

tupleType
    : LPAREN RPAREN
    | LPAREN tupleTypeElementList RPAREN
    ;

tupleTypeElementList
    : tupleTypeElement (COMMA tupleTypeElement)*
    ;

tupleTypeElement
    : identifier typeAnnotation
    | type
    ;

protocolCompositionType
    : protocolCompositionTypeElement (AMPERSAND protocolCompositionTypeElement)+
    | ANY AMPERSAND protocolCompositionTypeElement (AMPERSAND protocolCompositionTypeElement)*
    ;

protocolCompositionTypeElement
    : typeIdentifier
    | suppressedType
    ;

suppressedType
    : OPERATOR typeIdentifier
    ;

existentialType
    : ANY type
    ;

opaqueType
    : SOME_KW type
    ;

selfType
    : SELF_UPPER
    ;

typeAnnotation
    : COLON attributes? INOUT? type
    ;

typeInheritanceClause
    : COLON typeInheritanceList
    ;

typeInheritanceList
    : typeIdentifier (COMMA typeIdentifier)*
    ;

// ─────────────────────────────────────────────────────────────
// GENERIC PARAMETERS AND ARGUMENTS
// ─────────────────────────────────────────────────────────────
genericParameterClause
    : LT genericParameterList GT
    ;

genericParameterList
    : genericParameter (COMMA genericParameter)*
    ;

genericParameter
    : identifier
    | identifier COLON typeIdentifier
    | identifier COLON protocolCompositionType
    ;

genericWhereClause
    : WHERE requirementList
    ;

requirementList
    : requirement (COMMA requirement)*
    ;

requirement
    : conformanceRequirement
    | sameTypeRequirement
    | layoutConstraintRequirement
    ;

conformanceRequirement
    : typeIdentifier COLON typeIdentifier
    | typeIdentifier COLON protocolCompositionType
    ;

sameTypeRequirement
    : typeIdentifier ASSIGN type
    ;

layoutConstraintRequirement
    : typeIdentifier COLON identifier
    ;

genericArgumentClause
    : LT genericArgumentList GT
    ;

genericArgumentList
    : genericArgument (COMMA genericArgument)*
    ;

genericArgument
    : type
    ;

// ─────────────────────────────────────────────────────────────
// PATTERNS
// ─────────────────────────────────────────────────────────────
pattern
    : wildcardPattern typeAnnotation?              # wildcardPat
    | identifierPattern typeAnnotation?            # identPat
    | valueBindingPattern                          # bindingPat
    | tuplePattern typeAnnotation?                 # tuplePat
    | enumCasePattern                              # enumCasePat
    | optionalPattern                              # optPat
    | IS type                                      # isPat
    | expressionPattern                            # exprPat
    | pattern AS type                              # asPat
    ;

wildcardPattern      : UNDERSCORE ;
identifierPattern    : identifier ;
valueBindingPattern  : (LET | VAR | INOUT) pattern ;

tuplePattern
    : LPAREN tuplePatternElementList? RPAREN
    ;

tuplePatternElementList
    : tuplePatternElement (COMMA tuplePatternElement)*
    ;

tuplePatternElement
    : pattern
    | identifier COLON pattern
    ;

enumCasePattern
    : typeIdentifier? DOT identifier tuplePattern?
    ;

optionalPattern
    : identifierPattern QUESTION_POSTFIX
    ;

expressionPattern
    : expression
    ;

// ─────────────────────────────────────────────────────────────
// ATTRIBUTES
// ─────────────────────────────────────────────────────────────
attributes
    : attribute+
    ;

attribute
    : AT identifier attributeArguments?
    ;

attributeArguments
    : LPAREN attributeArgumentList? RPAREN
    ;

attributeArgumentList
    : attributeArgument (COMMA attributeArgument)*
    ;

attributeArgument
    : expression
    | identifier COLON expression
    | type
    ;

// ─────────────────────────────────────────────────────────────
// AVAILABILITY CONDITIONS
// ─────────────────────────────────────────────────────────────
availabilityCondition
    : POUND_AVAILABLE   LPAREN availabilityArguments RPAREN
    | POUND_UNAVAILABLE LPAREN availabilityArguments RPAREN
    ;

availabilityArguments
    : availabilityArgument (COMMA availabilityArgument)*
    ;

availabilityArgument
    : identifier decimalVersion
    | OPERATOR
    ;

// ─────────────────────────────────────────────────────────────
// IDENTIFIER RECOVERY (Fallback logic for Soft Keywords)
// ─────────────────────────────────────────────────────────────
identifier
    : IDENTIFIER
    | OS_KW | ARCH_KW | VERTEX_KW | COMPILER_KW | CANIMPORT_KW | VERSION_KW
    | FILE_KEYWORD | LINE_KEYWORD | GET_KW | SET_KW | WILLSET_KW | DIDSET_KW
    | ASYNC_KW | ACTOR_KW | PREFIX_KW | POSTFIX_KW
    | MACRO_KW | DYNAMIC_KW | FINAL_KW | LAZY_KW | OPTIONAL_KW | OVERRIDE_KW | REQUIRED_KW
    | UNOWNED_KW | WEAK_KW | NONISOLATED_KW
    | MUTATING_KW | NONMUTATING_KW
    | SOME_KW | TYPE_KW | PROTOCOL_KW | CONSUMING_KW | BORROWING_KW | SENDING_KW
    ;