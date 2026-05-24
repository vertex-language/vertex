

CONSTANTS

const (
	VertexLexerBLOCK_COMMENT   = 1
	VertexLexerLINE_COMMENT    = 2
	VertexLexerBREAK           = 3
	VertexLexerCASE            = 4
	VertexLexerCONTINUE        = 5
	VertexLexerDEFAULT         = 6
	VertexLexerDEFER           = 7
	VertexLexerELSE            = 8
	VertexLexerFALLTHROUGH     = 9
	VertexLexerFOR             = 10
	VertexLexerIF              = 11
	VertexLexerIN              = 12
	VertexLexerRETURN          = 13
	VertexLexerSWITCH          = 14
	VertexLexerWHILE           = 15
	VertexLexerBUILD           = 16
	VertexLexerCLASS           = 17
	VertexLexerENUM            = 18
	VertexLexerFUNC            = 19
	VertexLexerIMPORT          = 20
	VertexLexerLET             = 21
	VertexLexerPACKAGE         = 22
	VertexLexerSTRUCT          = 23
	VertexLexerTYPE            = 24
	VertexLexerVAR             = 25
	VertexLexerWEAK            = 26
	VertexLexerMUT             = 27
	VertexLexerASYNC           = 28
	VertexLexerTHREAD          = 29
	VertexLexerPROCESS         = 30
	VertexLexerGPU             = 31
	VertexLexerANY             = 32
	VertexLexerCHANNEL         = 33
	VertexLexerOPAQUE          = 34
	VertexLexerRESULT          = 35
	VertexLexerOK              = 36
	VertexLexerERR             = 37
	VertexLexerFALSE           = 38
	VertexLexerNIL             = 39
	VertexLexerTRUE            = 40
	VertexLexerINT             = 41
	VertexLexerINT8            = 42
	VertexLexerINT16           = 43
	VertexLexerINT32           = 44
	VertexLexerINT64           = 45
	VertexLexerUINT            = 46
	VertexLexerUINT8           = 47
	VertexLexerUINT16          = 48
	VertexLexerUINT32          = 49
	VertexLexerUINT64          = 50
	VertexLexerFLOAT           = 51
	VertexLexerDOUBLE          = 52
	VertexLexerBOOL            = 53
	VertexLexerSTRING          = 54
	VertexLexerCHAR            = 55
	VertexLexerVOID            = 56
	VertexLexerOVERFLOW_ADD    = 57
	VertexLexerOVERFLOW_SUB    = 58
	VertexLexerOVERFLOW_MUL    = 59
	VertexLexerPLUS_ASSIGN     = 60
	VertexLexerMINUS_ASSIGN    = 61
	VertexLexerSTAR_ASSIGN     = 62
	VertexLexerSLASH_ASSIGN    = 63
	VertexLexerMOD_ASSIGN      = 64
	VertexLexerIDENTITY_EQ     = 65
	VertexLexerIDENTITY_NEQ    = 66
	VertexLexerEQ              = 67
	VertexLexerNEQ             = 68
	VertexLexerGTE             = 69
	VertexLexerLTE             = 70
	VertexLexerLSHIFT          = 71
	VertexLexerRSHIFT          = 72
	VertexLexerAND             = 73
	VertexLexerOR              = 74
	VertexLexerARROW           = 75
	VertexLexerNIL_COALESCE    = 76
	VertexLexerHALF_OPEN_RANGE = 77
	VertexLexerELLIPSIS        = 78
	VertexLexerPLUS            = 79
	VertexLexerMINUS           = 80
	VertexLexerSTAR            = 81
	VertexLexerSLASH           = 82
	VertexLexerPERCENT         = 83
	VertexLexerAMP             = 84
	VertexLexerPIPE            = 85
	VertexLexerCARET           = 86
	VertexLexerTILDE           = 87
	VertexLexerBANG            = 88
	VertexLexerGT              = 89
	VertexLexerLT              = 90
	VertexLexerASSIGN          = 91
	VertexLexerQUESTION        = 92
	VertexLexerLPAREN          = 93
	VertexLexerRPAREN          = 94
	VertexLexerLBRACE          = 95
	VertexLexerRBRACE          = 96
	VertexLexerLBRACKET        = 97
	VertexLexerRBRACKET        = 98
	VertexLexerCOLON           = 99
	VertexLexerCOMMA           = 100
	VertexLexerDOT             = 101
	VertexLexerSEMICOLON       = 102
	VertexLexerHEX_FLOAT_LIT   = 103
	VertexLexerHEX_INT_LIT     = 104
	VertexLexerBIN_INT_LIT     = 105
	VertexLexerOCT_INT_LIT     = 106
	VertexLexerDEC_FLOAT_LIT   = 107
	VertexLexerDEC_INT_LIT     = 108
	VertexLexerSTRING_LIT      = 109
	VertexLexerRAW_STRING_LIT  = 110
	VertexLexerID              = 111
	VertexLexerWS              = 112
)
    VertexLexer tokens.

const (
	VertexParserEOF             = antlr.TokenEOF
	VertexParserBLOCK_COMMENT   = 1
	VertexParserLINE_COMMENT    = 2
	VertexParserBREAK           = 3
	VertexParserCASE            = 4
	VertexParserCONTINUE        = 5
	VertexParserDEFAULT         = 6
	VertexParserDEFER           = 7
	VertexParserELSE            = 8
	VertexParserFALLTHROUGH     = 9
	VertexParserFOR             = 10
	VertexParserIF              = 11
	VertexParserIN              = 12
	VertexParserRETURN          = 13
	VertexParserSWITCH          = 14
	VertexParserWHILE           = 15
	VertexParserBUILD           = 16
	VertexParserCLASS           = 17
	VertexParserENUM            = 18
	VertexParserFUNC            = 19
	VertexParserIMPORT          = 20
	VertexParserLET             = 21
	VertexParserPACKAGE         = 22
	VertexParserSTRUCT          = 23
	VertexParserTYPE            = 24
	VertexParserVAR             = 25
	VertexParserWEAK            = 26
	VertexParserMUT             = 27
	VertexParserASYNC           = 28
	VertexParserTHREAD          = 29
	VertexParserPROCESS         = 30
	VertexParserGPU             = 31
	VertexParserANY             = 32
	VertexParserCHANNEL         = 33
	VertexParserOPAQUE          = 34
	VertexParserRESULT          = 35
	VertexParserOK              = 36
	VertexParserERR             = 37
	VertexParserFALSE           = 38
	VertexParserNIL             = 39
	VertexParserTRUE            = 40
	VertexParserINT             = 41
	VertexParserINT8            = 42
	VertexParserINT16           = 43
	VertexParserINT32           = 44
	VertexParserINT64           = 45
	VertexParserUINT            = 46
	VertexParserUINT8           = 47
	VertexParserUINT16          = 48
	VertexParserUINT32          = 49
	VertexParserUINT64          = 50
	VertexParserFLOAT           = 51
	VertexParserDOUBLE          = 52
	VertexParserBOOL            = 53
	VertexParserSTRING          = 54
	VertexParserCHAR            = 55
	VertexParserVOID            = 56
	VertexParserOVERFLOW_ADD    = 57
	VertexParserOVERFLOW_SUB    = 58
	VertexParserOVERFLOW_MUL    = 59
	VertexParserPLUS_ASSIGN     = 60
	VertexParserMINUS_ASSIGN    = 61
	VertexParserSTAR_ASSIGN     = 62
	VertexParserSLASH_ASSIGN    = 63
	VertexParserMOD_ASSIGN      = 64
	VertexParserIDENTITY_EQ     = 65
	VertexParserIDENTITY_NEQ    = 66
	VertexParserEQ              = 67
	VertexParserNEQ             = 68
	VertexParserGTE             = 69
	VertexParserLTE             = 70
	VertexParserLSHIFT          = 71
	VertexParserRSHIFT          = 72
	VertexParserAND             = 73
	VertexParserOR              = 74
	VertexParserARROW           = 75
	VertexParserNIL_COALESCE    = 76
	VertexParserHALF_OPEN_RANGE = 77
	VertexParserELLIPSIS        = 78
	VertexParserPLUS            = 79
	VertexParserMINUS           = 80
	VertexParserSTAR            = 81
	VertexParserSLASH           = 82
	VertexParserPERCENT         = 83
	VertexParserAMP             = 84
	VertexParserPIPE            = 85
	VertexParserCARET           = 86
	VertexParserTILDE           = 87
	VertexParserBANG            = 88
	VertexParserGT              = 89
	VertexParserLT              = 90
	VertexParserASSIGN          = 91
	VertexParserQUESTION        = 92
	VertexParserLPAREN          = 93
	VertexParserRPAREN          = 94
	VertexParserLBRACE          = 95
	VertexParserRBRACE          = 96
	VertexParserLBRACKET        = 97
	VertexParserRBRACKET        = 98
	VertexParserCOLON           = 99
	VertexParserCOMMA           = 100
	VertexParserDOT             = 101
	VertexParserSEMICOLON       = 102
	VertexParserHEX_FLOAT_LIT   = 103
	VertexParserHEX_INT_LIT     = 104
	VertexParserBIN_INT_LIT     = 105
	VertexParserOCT_INT_LIT     = 106
	VertexParserDEC_FLOAT_LIT   = 107
	VertexParserDEC_INT_LIT     = 108
	VertexParserSTRING_LIT      = 109
	VertexParserRAW_STRING_LIT  = 110
	VertexParserID              = 111
	VertexParserWS              = 112
)
    VertexParser tokens.

const (
	VertexParserRULE_file               = 0
	VertexParserRULE_packageDecl        = 1
	VertexParserRULE_buildDecl          = 2
	VertexParserRULE_buildTag           = 3
	VertexParserRULE_importDecl         = 4
	VertexParserRULE_topLevelDecl       = 5
	VertexParserRULE_funcDecl           = 6
	VertexParserRULE_genericParams      = 7
	VertexParserRULE_typeParam          = 8
	VertexParserRULE_paramList          = 9
	VertexParserRULE_param              = 10
	VertexParserRULE_funcQualifier      = 11
	VertexParserRULE_returnType         = 12
	VertexParserRULE_structDecl         = 13
	VertexParserRULE_structField        = 14
	VertexParserRULE_structLiteralExpr  = 15
	VertexParserRULE_structFieldInit    = 16
	VertexParserRULE_classDecl          = 17
	VertexParserRULE_classMember        = 18
	VertexParserRULE_classField         = 19
	VertexParserRULE_nativeFuncDecl     = 20
	VertexParserRULE_nativeParamList    = 21
	VertexParserRULE_nativeParam        = 22
	VertexParserRULE_enumDecl           = 23
	VertexParserRULE_enumRawType        = 24
	VertexParserRULE_enumCaseDecl       = 25
	VertexParserRULE_enumCase           = 26
	VertexParserRULE_typeAliasDecl      = 27
	VertexParserRULE_block              = 28
	VertexParserRULE_stmt               = 29
	VertexParserRULE_varDeclStmt        = 30
	VertexParserRULE_bindingKw          = 31
	VertexParserRULE_tupleBind          = 32
	VertexParserRULE_assignStmt         = 33
	VertexParserRULE_compoundAssignStmt = 34
	VertexParserRULE_compoundOp         = 35
	VertexParserRULE_lvalue             = 36
	VertexParserRULE_ifStmt             = 37
	VertexParserRULE_elseIfClause       = 38
	VertexParserRULE_elseClause         = 39
	VertexParserRULE_ifCondition        = 40
	VertexParserRULE_switchStmt         = 41
	VertexParserRULE_switchCase         = 42
	VertexParserRULE_defaultCase        = 43
	VertexParserRULE_casePatternList    = 44
	VertexParserRULE_casePattern        = 45
	VertexParserRULE_forInStmt          = 46
	VertexParserRULE_whileStmt          = 47
	VertexParserRULE_breakStmt          = 48
	VertexParserRULE_continueStmt       = 49
	VertexParserRULE_fallthroughStmt    = 50
	VertexParserRULE_returnStmt         = 51
	VertexParserRULE_deferStmt          = 52
	VertexParserRULE_exprStmt           = 53
	VertexParserRULE_expr               = 54
	VertexParserRULE_postfixName        = 55
	VertexParserRULE_primary            = 56
	VertexParserRULE_argList            = 57
	VertexParserRULE_arg                = 58
	VertexParserRULE_literal            = 59
	VertexParserRULE_arrayLiteralExpr   = 60
	VertexParserRULE_arrayConstructExpr = 61
	VertexParserRULE_exprList           = 62
	VertexParserRULE_dictLiteralExpr    = 63
	VertexParserRULE_dictEntry          = 64
	VertexParserRULE_tupleExpr          = 65
	VertexParserRULE_tupleElement       = 66
	VertexParserRULE_anonFuncExpr       = 67
	VertexParserRULE_type               = 68
	VertexParserRULE_tupleTypeElem      = 69
	VertexParserRULE_funcTypeParamList  = 70
	VertexParserRULE_funcTypeParam      = 71
	VertexParserRULE_primitiveType      = 72
)
    VertexParser rules.


VARIABLES

var VertexLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}
var VertexParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

FUNCTIONS

func InitEmptyAnonFuncExprContext(p *AnonFuncExprContext)
func InitEmptyArgContext(p *ArgContext)
func InitEmptyArgListContext(p *ArgListContext)
func InitEmptyArrayConstructExprContext(p *ArrayConstructExprContext)
func InitEmptyArrayLiteralExprContext(p *ArrayLiteralExprContext)
func InitEmptyAssignStmtContext(p *AssignStmtContext)
func InitEmptyBindingKwContext(p *BindingKwContext)
func InitEmptyBlockContext(p *BlockContext)
func InitEmptyBreakStmtContext(p *BreakStmtContext)
func InitEmptyBuildDeclContext(p *BuildDeclContext)
func InitEmptyBuildTagContext(p *BuildTagContext)
func InitEmptyCasePatternContext(p *CasePatternContext)
func InitEmptyCasePatternListContext(p *CasePatternListContext)
func InitEmptyClassDeclContext(p *ClassDeclContext)
func InitEmptyClassFieldContext(p *ClassFieldContext)
func InitEmptyClassMemberContext(p *ClassMemberContext)
func InitEmptyCompoundAssignStmtContext(p *CompoundAssignStmtContext)
func InitEmptyCompoundOpContext(p *CompoundOpContext)
func InitEmptyContinueStmtContext(p *ContinueStmtContext)
func InitEmptyDefaultCaseContext(p *DefaultCaseContext)
func InitEmptyDeferStmtContext(p *DeferStmtContext)
func InitEmptyDictEntryContext(p *DictEntryContext)
func InitEmptyDictLiteralExprContext(p *DictLiteralExprContext)
func InitEmptyElseClauseContext(p *ElseClauseContext)
func InitEmptyElseIfClauseContext(p *ElseIfClauseContext)
func InitEmptyEnumCaseContext(p *EnumCaseContext)
func InitEmptyEnumCaseDeclContext(p *EnumCaseDeclContext)
func InitEmptyEnumDeclContext(p *EnumDeclContext)
func InitEmptyEnumRawTypeContext(p *EnumRawTypeContext)
func InitEmptyExprContext(p *ExprContext)
func InitEmptyExprListContext(p *ExprListContext)
func InitEmptyExprStmtContext(p *ExprStmtContext)
func InitEmptyFallthroughStmtContext(p *FallthroughStmtContext)
func InitEmptyFileContext(p *FileContext)
func InitEmptyForInStmtContext(p *ForInStmtContext)
func InitEmptyFuncDeclContext(p *FuncDeclContext)
func InitEmptyFuncQualifierContext(p *FuncQualifierContext)
func InitEmptyFuncTypeParamContext(p *FuncTypeParamContext)
func InitEmptyFuncTypeParamListContext(p *FuncTypeParamListContext)
func InitEmptyGenericParamsContext(p *GenericParamsContext)
func InitEmptyIfConditionContext(p *IfConditionContext)
func InitEmptyIfStmtContext(p *IfStmtContext)
func InitEmptyImportDeclContext(p *ImportDeclContext)
func InitEmptyLiteralContext(p *LiteralContext)
func InitEmptyLvalueContext(p *LvalueContext)
func InitEmptyNativeFuncDeclContext(p *NativeFuncDeclContext)
func InitEmptyNativeParamContext(p *NativeParamContext)
func InitEmptyNativeParamListContext(p *NativeParamListContext)
func InitEmptyPackageDeclContext(p *PackageDeclContext)
func InitEmptyParamContext(p *ParamContext)
func InitEmptyParamListContext(p *ParamListContext)
func InitEmptyPostfixNameContext(p *PostfixNameContext)
func InitEmptyPrimaryContext(p *PrimaryContext)
func InitEmptyPrimitiveTypeContext(p *PrimitiveTypeContext)
func InitEmptyReturnStmtContext(p *ReturnStmtContext)
func InitEmptyReturnTypeContext(p *ReturnTypeContext)
func InitEmptyStmtContext(p *StmtContext)
func InitEmptyStructDeclContext(p *StructDeclContext)
func InitEmptyStructFieldContext(p *StructFieldContext)
func InitEmptyStructFieldInitContext(p *StructFieldInitContext)
func InitEmptyStructLiteralExprContext(p *StructLiteralExprContext)
func InitEmptySwitchCaseContext(p *SwitchCaseContext)
func InitEmptySwitchStmtContext(p *SwitchStmtContext)
func InitEmptyTopLevelDeclContext(p *TopLevelDeclContext)
func InitEmptyTupleBindContext(p *TupleBindContext)
func InitEmptyTupleElementContext(p *TupleElementContext)
func InitEmptyTupleExprContext(p *TupleExprContext)
func InitEmptyTupleTypeElemContext(p *TupleTypeElemContext)
func InitEmptyTypeAliasDeclContext(p *TypeAliasDeclContext)
func InitEmptyTypeContext(p *TypeContext)
func InitEmptyTypeParamContext(p *TypeParamContext)
func InitEmptyVarDeclStmtContext(p *VarDeclStmtContext)
func InitEmptyWhileStmtContext(p *WhileStmtContext)
func VertexLexerInit()
    VertexLexerInit initializes any static state used to implement VertexLexer.
    By default the static state used to implement the lexer is lazily
    initialized during the first call to NewVertexLexer(). You can call this
    function if you wish to initialize the static state ahead of time.

func VertexParserInit()
    VertexParserInit initializes any static state used to implement
    VertexParser. By default the static state used to implement the parser is
    lazily initialized during the first call to NewVertexParser(). You can call
    this function if you wish to initialize the static state ahead of time.


TYPES

type AnonFuncExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAnonFuncExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnonFuncExprContext

func NewEmptyAnonFuncExprContext() *AnonFuncExprContext

func (s *AnonFuncExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AnonFuncExprContext) Block() IBlockContext

func (s *AnonFuncExprContext) FUNC() antlr.TerminalNode

func (s *AnonFuncExprContext) GetParser() antlr.Parser

func (s *AnonFuncExprContext) GetRuleContext() antlr.RuleContext

func (*AnonFuncExprContext) IsAnonFuncExprContext()

func (s *AnonFuncExprContext) LPAREN() antlr.TerminalNode

func (s *AnonFuncExprContext) ParamList() IParamListContext

func (s *AnonFuncExprContext) RPAREN() antlr.TerminalNode

func (s *AnonFuncExprContext) ReturnType() IReturnTypeContext

func (s *AnonFuncExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ArgContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgContext

func NewEmptyArgContext() *ArgContext

func (s *ArgContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArgContext) COLON() antlr.TerminalNode

func (s *ArgContext) Expr() IExprContext

func (s *ArgContext) GetParser() antlr.Parser

func (s *ArgContext) GetRuleContext() antlr.RuleContext

func (s *ArgContext) ID() antlr.TerminalNode

func (*ArgContext) IsArgContext()

func (s *ArgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ArgListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArgListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgListContext

func NewEmptyArgListContext() *ArgListContext

func (s *ArgListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArgListContext) AllArg() []IArgContext

func (s *ArgListContext) AllCOMMA() []antlr.TerminalNode

func (s *ArgListContext) Arg(i int) IArgContext

func (s *ArgListContext) COMMA(i int) antlr.TerminalNode

func (s *ArgListContext) GetParser() antlr.Parser

func (s *ArgListContext) GetRuleContext() antlr.RuleContext

func (*ArgListContext) IsArgListContext()

func (s *ArgListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ArrayConstructExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArrayConstructExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayConstructExprContext

func NewEmptyArrayConstructExprContext() *ArrayConstructExprContext

func (s *ArrayConstructExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArrayConstructExprContext) AllCOLON() []antlr.TerminalNode

func (s *ArrayConstructExprContext) AllExpr() []IExprContext

func (s *ArrayConstructExprContext) AllID() []antlr.TerminalNode

func (s *ArrayConstructExprContext) COLON(i int) antlr.TerminalNode

func (s *ArrayConstructExprContext) COMMA() antlr.TerminalNode

func (s *ArrayConstructExprContext) Expr(i int) IExprContext

func (s *ArrayConstructExprContext) GetParser() antlr.Parser

func (s *ArrayConstructExprContext) GetRuleContext() antlr.RuleContext

func (s *ArrayConstructExprContext) ID(i int) antlr.TerminalNode

func (*ArrayConstructExprContext) IsArrayConstructExprContext()

func (s *ArrayConstructExprContext) LBRACKET() antlr.TerminalNode

func (s *ArrayConstructExprContext) LPAREN() antlr.TerminalNode

func (s *ArrayConstructExprContext) RBRACKET() antlr.TerminalNode

func (s *ArrayConstructExprContext) RPAREN() antlr.TerminalNode

func (s *ArrayConstructExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ArrayConstructExprContext) Type_() ITypeContext

type ArrayLiteralExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArrayLiteralExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayLiteralExprContext

func NewEmptyArrayLiteralExprContext() *ArrayLiteralExprContext

func (s *ArrayLiteralExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArrayLiteralExprContext) COMMA() antlr.TerminalNode

func (s *ArrayLiteralExprContext) ExprList() IExprListContext

func (s *ArrayLiteralExprContext) GetParser() antlr.Parser

func (s *ArrayLiteralExprContext) GetRuleContext() antlr.RuleContext

func (*ArrayLiteralExprContext) IsArrayLiteralExprContext()

func (s *ArrayLiteralExprContext) LBRACKET() antlr.TerminalNode

func (s *ArrayLiteralExprContext) RBRACKET() antlr.TerminalNode

func (s *ArrayLiteralExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AssignStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAssignStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignStmtContext

func NewEmptyAssignStmtContext() *AssignStmtContext

func (s *AssignStmtContext) ASSIGN() antlr.TerminalNode

func (s *AssignStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AssignStmtContext) Expr() IExprContext

func (s *AssignStmtContext) GetParser() antlr.Parser

func (s *AssignStmtContext) GetRuleContext() antlr.RuleContext

func (*AssignStmtContext) IsAssignStmtContext()

func (s *AssignStmtContext) Lvalue() ILvalueContext

func (s *AssignStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BaseVertexParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseVertexParserVisitor) VisitAnonFuncExpr(ctx *AnonFuncExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitArg(ctx *ArgContext) interface{}

func (v *BaseVertexParserVisitor) VisitArgList(ctx *ArgListContext) interface{}

func (v *BaseVertexParserVisitor) VisitArrayConstructExpr(ctx *ArrayConstructExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitArrayLiteralExpr(ctx *ArrayLiteralExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitAssignStmt(ctx *AssignStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitBindingKw(ctx *BindingKwContext) interface{}

func (v *BaseVertexParserVisitor) VisitBlock(ctx *BlockContext) interface{}

func (v *BaseVertexParserVisitor) VisitBreakStmt(ctx *BreakStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitBuildDecl(ctx *BuildDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitBuildTag(ctx *BuildTagContext) interface{}

func (v *BaseVertexParserVisitor) VisitCasePattern(ctx *CasePatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitCasePatternList(ctx *CasePatternListContext) interface{}

func (v *BaseVertexParserVisitor) VisitClassDecl(ctx *ClassDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitClassField(ctx *ClassFieldContext) interface{}

func (v *BaseVertexParserVisitor) VisitClassMember(ctx *ClassMemberContext) interface{}

func (v *BaseVertexParserVisitor) VisitCompoundAssignStmt(ctx *CompoundAssignStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitCompoundOp(ctx *CompoundOpContext) interface{}

func (v *BaseVertexParserVisitor) VisitContinueStmt(ctx *ContinueStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitDefaultCase(ctx *DefaultCaseContext) interface{}

func (v *BaseVertexParserVisitor) VisitDeferStmt(ctx *DeferStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitDictEntry(ctx *DictEntryContext) interface{}

func (v *BaseVertexParserVisitor) VisitDictLiteralExpr(ctx *DictLiteralExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitElseClause(ctx *ElseClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitElseIfClause(ctx *ElseIfClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumCase(ctx *EnumCaseContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumCaseDecl(ctx *EnumCaseDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumDecl(ctx *EnumDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumRawType(ctx *EnumRawTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitExpr(ctx *ExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitExprList(ctx *ExprListContext) interface{}

func (v *BaseVertexParserVisitor) VisitExprStmt(ctx *ExprStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitFallthroughStmt(ctx *FallthroughStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitFile(ctx *FileContext) interface{}

func (v *BaseVertexParserVisitor) VisitForInStmt(ctx *ForInStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitFuncDecl(ctx *FuncDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitFuncQualifier(ctx *FuncQualifierContext) interface{}

func (v *BaseVertexParserVisitor) VisitFuncTypeParam(ctx *FuncTypeParamContext) interface{}

func (v *BaseVertexParserVisitor) VisitFuncTypeParamList(ctx *FuncTypeParamListContext) interface{}

func (v *BaseVertexParserVisitor) VisitGenericParams(ctx *GenericParamsContext) interface{}

func (v *BaseVertexParserVisitor) VisitIfCondition(ctx *IfConditionContext) interface{}

func (v *BaseVertexParserVisitor) VisitIfStmt(ctx *IfStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitImportDecl(ctx *ImportDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitLiteral(ctx *LiteralContext) interface{}

func (v *BaseVertexParserVisitor) VisitLvalue(ctx *LvalueContext) interface{}

func (v *BaseVertexParserVisitor) VisitNativeFuncDecl(ctx *NativeFuncDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitNativeParam(ctx *NativeParamContext) interface{}

func (v *BaseVertexParserVisitor) VisitNativeParamList(ctx *NativeParamListContext) interface{}

func (v *BaseVertexParserVisitor) VisitPackageDecl(ctx *PackageDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitParam(ctx *ParamContext) interface{}

func (v *BaseVertexParserVisitor) VisitParamList(ctx *ParamListContext) interface{}

func (v *BaseVertexParserVisitor) VisitPostfixName(ctx *PostfixNameContext) interface{}

func (v *BaseVertexParserVisitor) VisitPrimary(ctx *PrimaryContext) interface{}

func (v *BaseVertexParserVisitor) VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitReturnStmt(ctx *ReturnStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitReturnType(ctx *ReturnTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitStmt(ctx *StmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructDecl(ctx *StructDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructField(ctx *StructFieldContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructFieldInit(ctx *StructFieldInitContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructLiteralExpr(ctx *StructLiteralExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{}

func (v *BaseVertexParserVisitor) VisitSwitchStmt(ctx *SwitchStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleBind(ctx *TupleBindContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleElement(ctx *TupleElementContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleExpr(ctx *TupleExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleTypeElem(ctx *TupleTypeElemContext) interface{}

func (v *BaseVertexParserVisitor) VisitType(ctx *TypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeAliasDecl(ctx *TypeAliasDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeParam(ctx *TypeParamContext) interface{}

func (v *BaseVertexParserVisitor) VisitVarDeclStmt(ctx *VarDeclStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitWhileStmt(ctx *WhileStmtContext) interface{}

type BindingKwContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBindingKwContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BindingKwContext

func NewEmptyBindingKwContext() *BindingKwContext

func (s *BindingKwContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BindingKwContext) GetParser() antlr.Parser

func (s *BindingKwContext) GetRuleContext() antlr.RuleContext

func (*BindingKwContext) IsBindingKwContext()

func (s *BindingKwContext) LET() antlr.TerminalNode

func (s *BindingKwContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *BindingKwContext) VAR() antlr.TerminalNode

func (s *BindingKwContext) WEAK() antlr.TerminalNode

type BlockContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BlockContext

func NewEmptyBlockContext() *BlockContext

func (s *BlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BlockContext) AllStmt() []IStmtContext

func (s *BlockContext) GetParser() antlr.Parser

func (s *BlockContext) GetRuleContext() antlr.RuleContext

func (*BlockContext) IsBlockContext()

func (s *BlockContext) LBRACE() antlr.TerminalNode

func (s *BlockContext) RBRACE() antlr.TerminalNode

func (s *BlockContext) Stmt(i int) IStmtContext

func (s *BlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BreakStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBreakStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BreakStmtContext

func NewEmptyBreakStmtContext() *BreakStmtContext

func (s *BreakStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BreakStmtContext) BREAK() antlr.TerminalNode

func (s *BreakStmtContext) GetParser() antlr.Parser

func (s *BreakStmtContext) GetRuleContext() antlr.RuleContext

func (*BreakStmtContext) IsBreakStmtContext()

func (s *BreakStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BuildDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBuildDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BuildDeclContext

func NewEmptyBuildDeclContext() *BuildDeclContext

func (s *BuildDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BuildDeclContext) AllBuildTag() []IBuildTagContext

func (s *BuildDeclContext) AllCOMMA() []antlr.TerminalNode

func (s *BuildDeclContext) BUILD() antlr.TerminalNode

func (s *BuildDeclContext) BuildTag(i int) IBuildTagContext

func (s *BuildDeclContext) COMMA(i int) antlr.TerminalNode

func (s *BuildDeclContext) GetParser() antlr.Parser

func (s *BuildDeclContext) GetRuleContext() antlr.RuleContext

func (*BuildDeclContext) IsBuildDeclContext()

func (s *BuildDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BuildTagContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBuildTagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BuildTagContext

func NewEmptyBuildTagContext() *BuildTagContext

func (s *BuildTagContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BuildTagContext) GetParser() antlr.Parser

func (s *BuildTagContext) GetRuleContext() antlr.RuleContext

func (s *BuildTagContext) ID() antlr.TerminalNode

func (*BuildTagContext) IsBuildTagContext()

func (s *BuildTagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CasePatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCasePatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CasePatternContext

func NewEmptyCasePatternContext() *CasePatternContext

func (s *CasePatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CasePatternContext) DOT() antlr.TerminalNode

func (s *CasePatternContext) ERR() antlr.TerminalNode

func (s *CasePatternContext) GetParser() antlr.Parser

func (s *CasePatternContext) GetRuleContext() antlr.RuleContext

func (s *CasePatternContext) ID() antlr.TerminalNode

func (*CasePatternContext) IsCasePatternContext()

func (s *CasePatternContext) LET() antlr.TerminalNode

func (s *CasePatternContext) LPAREN() antlr.TerminalNode

func (s *CasePatternContext) Literal() ILiteralContext

func (s *CasePatternContext) OK() antlr.TerminalNode

func (s *CasePatternContext) RPAREN() antlr.TerminalNode

func (s *CasePatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CasePatternListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCasePatternListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CasePatternListContext

func NewEmptyCasePatternListContext() *CasePatternListContext

func (s *CasePatternListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CasePatternListContext) AllCOMMA() []antlr.TerminalNode

func (s *CasePatternListContext) AllCasePattern() []ICasePatternContext

func (s *CasePatternListContext) COMMA(i int) antlr.TerminalNode

func (s *CasePatternListContext) CasePattern(i int) ICasePatternContext

func (s *CasePatternListContext) GetParser() antlr.Parser

func (s *CasePatternListContext) GetRuleContext() antlr.RuleContext

func (*CasePatternListContext) IsCasePatternListContext()

func (s *CasePatternListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClassDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClassDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassDeclContext

func NewEmptyClassDeclContext() *ClassDeclContext

func (s *ClassDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClassDeclContext) AllClassMember() []IClassMemberContext

func (s *ClassDeclContext) AllID() []antlr.TerminalNode

func (s *ClassDeclContext) CLASS() antlr.TerminalNode

func (s *ClassDeclContext) COLON() antlr.TerminalNode

func (s *ClassDeclContext) ClassMember(i int) IClassMemberContext

func (s *ClassDeclContext) GetParser() antlr.Parser

func (s *ClassDeclContext) GetRuleContext() antlr.RuleContext

func (s *ClassDeclContext) ID(i int) antlr.TerminalNode

func (*ClassDeclContext) IsClassDeclContext()

func (s *ClassDeclContext) LBRACE() antlr.TerminalNode

func (s *ClassDeclContext) RBRACE() antlr.TerminalNode

func (s *ClassDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClassFieldContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClassFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassFieldContext

func NewEmptyClassFieldContext() *ClassFieldContext

func (s *ClassFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClassFieldContext) COLON() antlr.TerminalNode

func (s *ClassFieldContext) GetParser() antlr.Parser

func (s *ClassFieldContext) GetRuleContext() antlr.RuleContext

func (s *ClassFieldContext) ID() antlr.TerminalNode

func (*ClassFieldContext) IsClassFieldContext()

func (s *ClassFieldContext) LET() antlr.TerminalNode

func (s *ClassFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ClassFieldContext) Type_() ITypeContext

func (s *ClassFieldContext) VAR() antlr.TerminalNode

type ClassMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClassMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassMemberContext

func NewEmptyClassMemberContext() *ClassMemberContext

func (s *ClassMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClassMemberContext) ClassField() IClassFieldContext

func (s *ClassMemberContext) GetParser() antlr.Parser

func (s *ClassMemberContext) GetRuleContext() antlr.RuleContext

func (*ClassMemberContext) IsClassMemberContext()

func (s *ClassMemberContext) NativeFuncDecl() INativeFuncDeclContext

func (s *ClassMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CompoundAssignStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCompoundAssignStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompoundAssignStmtContext

func NewEmptyCompoundAssignStmtContext() *CompoundAssignStmtContext

func (s *CompoundAssignStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CompoundAssignStmtContext) CompoundOp() ICompoundOpContext

func (s *CompoundAssignStmtContext) Expr() IExprContext

func (s *CompoundAssignStmtContext) GetParser() antlr.Parser

func (s *CompoundAssignStmtContext) GetRuleContext() antlr.RuleContext

func (*CompoundAssignStmtContext) IsCompoundAssignStmtContext()

func (s *CompoundAssignStmtContext) Lvalue() ILvalueContext

func (s *CompoundAssignStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CompoundOpContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCompoundOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompoundOpContext

func NewEmptyCompoundOpContext() *CompoundOpContext

func (s *CompoundOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CompoundOpContext) GetParser() antlr.Parser

func (s *CompoundOpContext) GetRuleContext() antlr.RuleContext

func (*CompoundOpContext) IsCompoundOpContext()

func (s *CompoundOpContext) MINUS_ASSIGN() antlr.TerminalNode

func (s *CompoundOpContext) MOD_ASSIGN() antlr.TerminalNode

func (s *CompoundOpContext) PLUS_ASSIGN() antlr.TerminalNode

func (s *CompoundOpContext) SLASH_ASSIGN() antlr.TerminalNode

func (s *CompoundOpContext) STAR_ASSIGN() antlr.TerminalNode

func (s *CompoundOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ContinueStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewContinueStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContinueStmtContext

func NewEmptyContinueStmtContext() *ContinueStmtContext

func (s *ContinueStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ContinueStmtContext) CONTINUE() antlr.TerminalNode

func (s *ContinueStmtContext) GetParser() antlr.Parser

func (s *ContinueStmtContext) GetRuleContext() antlr.RuleContext

func (*ContinueStmtContext) IsContinueStmtContext()

func (s *ContinueStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DefaultCaseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDefaultCaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DefaultCaseContext

func NewEmptyDefaultCaseContext() *DefaultCaseContext

func (s *DefaultCaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DefaultCaseContext) AllStmt() []IStmtContext

func (s *DefaultCaseContext) COLON() antlr.TerminalNode

func (s *DefaultCaseContext) DEFAULT() antlr.TerminalNode

func (s *DefaultCaseContext) GetParser() antlr.Parser

func (s *DefaultCaseContext) GetRuleContext() antlr.RuleContext

func (*DefaultCaseContext) IsDefaultCaseContext()

func (s *DefaultCaseContext) Stmt(i int) IStmtContext

func (s *DefaultCaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DeferStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDeferStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeferStmtContext

func NewEmptyDeferStmtContext() *DeferStmtContext

func (s *DeferStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeferStmtContext) DEFER() antlr.TerminalNode

func (s *DeferStmtContext) Expr() IExprContext

func (s *DeferStmtContext) GetParser() antlr.Parser

func (s *DeferStmtContext) GetRuleContext() antlr.RuleContext

func (*DeferStmtContext) IsDeferStmtContext()

func (s *DeferStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DictEntryContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDictEntryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DictEntryContext

func NewEmptyDictEntryContext() *DictEntryContext

func (s *DictEntryContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DictEntryContext) AllExpr() []IExprContext

func (s *DictEntryContext) COLON() antlr.TerminalNode

func (s *DictEntryContext) Expr(i int) IExprContext

func (s *DictEntryContext) GetParser() antlr.Parser

func (s *DictEntryContext) GetRuleContext() antlr.RuleContext

func (*DictEntryContext) IsDictEntryContext()

func (s *DictEntryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DictLiteralExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDictLiteralExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DictLiteralExprContext

func NewEmptyDictLiteralExprContext() *DictLiteralExprContext

func (s *DictLiteralExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DictLiteralExprContext) AllCOMMA() []antlr.TerminalNode

func (s *DictLiteralExprContext) AllDictEntry() []IDictEntryContext

func (s *DictLiteralExprContext) COMMA(i int) antlr.TerminalNode

func (s *DictLiteralExprContext) DictEntry(i int) IDictEntryContext

func (s *DictLiteralExprContext) GetParser() antlr.Parser

func (s *DictLiteralExprContext) GetRuleContext() antlr.RuleContext

func (*DictLiteralExprContext) IsDictLiteralExprContext()

func (s *DictLiteralExprContext) LBRACKET() antlr.TerminalNode

func (s *DictLiteralExprContext) RBRACKET() antlr.TerminalNode

func (s *DictLiteralExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ElseClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewElseClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElseClauseContext

func NewEmptyElseClauseContext() *ElseClauseContext

func (s *ElseClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ElseClauseContext) Block() IBlockContext

func (s *ElseClauseContext) ELSE() antlr.TerminalNode

func (s *ElseClauseContext) GetParser() antlr.Parser

func (s *ElseClauseContext) GetRuleContext() antlr.RuleContext

func (*ElseClauseContext) IsElseClauseContext()

func (s *ElseClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ElseIfClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewElseIfClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElseIfClauseContext

func NewEmptyElseIfClauseContext() *ElseIfClauseContext

func (s *ElseIfClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ElseIfClauseContext) Block() IBlockContext

func (s *ElseIfClauseContext) ELSE() antlr.TerminalNode

func (s *ElseIfClauseContext) GetParser() antlr.Parser

func (s *ElseIfClauseContext) GetRuleContext() antlr.RuleContext

func (s *ElseIfClauseContext) IF() antlr.TerminalNode

func (s *ElseIfClauseContext) IfCondition() IIfConditionContext

func (*ElseIfClauseContext) IsElseIfClauseContext()

func (s *ElseIfClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type EnumCaseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumCaseContext() *EnumCaseContext

func NewEnumCaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumCaseContext

func (s *EnumCaseContext) ASSIGN() antlr.TerminalNode

func (s *EnumCaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumCaseContext) GetParser() antlr.Parser

func (s *EnumCaseContext) GetRuleContext() antlr.RuleContext

func (s *EnumCaseContext) ID() antlr.TerminalNode

func (*EnumCaseContext) IsEnumCaseContext()

func (s *EnumCaseContext) Literal() ILiteralContext

func (s *EnumCaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type EnumCaseDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumCaseDeclContext() *EnumCaseDeclContext

func NewEnumCaseDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumCaseDeclContext

func (s *EnumCaseDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumCaseDeclContext) AllCOMMA() []antlr.TerminalNode

func (s *EnumCaseDeclContext) AllEnumCase() []IEnumCaseContext

func (s *EnumCaseDeclContext) CASE() antlr.TerminalNode

func (s *EnumCaseDeclContext) COMMA(i int) antlr.TerminalNode

func (s *EnumCaseDeclContext) EnumCase(i int) IEnumCaseContext

func (s *EnumCaseDeclContext) GetParser() antlr.Parser

func (s *EnumCaseDeclContext) GetRuleContext() antlr.RuleContext

func (*EnumCaseDeclContext) IsEnumCaseDeclContext()

func (s *EnumCaseDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type EnumDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumDeclContext() *EnumDeclContext

func NewEnumDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumDeclContext

func (s *EnumDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumDeclContext) AllEnumCaseDecl() []IEnumCaseDeclContext

func (s *EnumDeclContext) COLON() antlr.TerminalNode

func (s *EnumDeclContext) ENUM() antlr.TerminalNode

func (s *EnumDeclContext) EnumCaseDecl(i int) IEnumCaseDeclContext

func (s *EnumDeclContext) EnumRawType() IEnumRawTypeContext

func (s *EnumDeclContext) GetParser() antlr.Parser

func (s *EnumDeclContext) GetRuleContext() antlr.RuleContext

func (s *EnumDeclContext) ID() antlr.TerminalNode

func (*EnumDeclContext) IsEnumDeclContext()

func (s *EnumDeclContext) LBRACE() antlr.TerminalNode

func (s *EnumDeclContext) RBRACE() antlr.TerminalNode

func (s *EnumDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type EnumRawTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumRawTypeContext() *EnumRawTypeContext

func NewEnumRawTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumRawTypeContext

func (s *EnumRawTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumRawTypeContext) GetParser() antlr.Parser

func (s *EnumRawTypeContext) GetRuleContext() antlr.RuleContext

func (s *EnumRawTypeContext) INT() antlr.TerminalNode

func (*EnumRawTypeContext) IsEnumRawTypeContext()

func (s *EnumRawTypeContext) STRING() antlr.TerminalNode

func (s *EnumRawTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExprContext() *ExprContext

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext

func (s *ExprContext) AMP() antlr.TerminalNode

func (s *ExprContext) AND() antlr.TerminalNode

func (s *ExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExprContext) AllExpr() []IExprContext

func (s *ExprContext) ArgList() IArgListContext

func (s *ExprContext) BANG() antlr.TerminalNode

func (s *ExprContext) CARET() antlr.TerminalNode

func (s *ExprContext) COLON() antlr.TerminalNode

func (s *ExprContext) COMMA() antlr.TerminalNode

func (s *ExprContext) DOT() antlr.TerminalNode

func (s *ExprContext) ELLIPSIS() antlr.TerminalNode

func (s *ExprContext) EQ() antlr.TerminalNode

func (s *ExprContext) ERR() antlr.TerminalNode

func (s *ExprContext) Expr(i int) IExprContext

func (s *ExprContext) GT() antlr.TerminalNode

func (s *ExprContext) GTE() antlr.TerminalNode

func (s *ExprContext) GetParser() antlr.Parser

func (s *ExprContext) GetRuleContext() antlr.RuleContext

func (s *ExprContext) HALF_OPEN_RANGE() antlr.TerminalNode

func (s *ExprContext) IDENTITY_EQ() antlr.TerminalNode

func (s *ExprContext) IDENTITY_NEQ() antlr.TerminalNode

func (*ExprContext) IsExprContext()

func (s *ExprContext) LBRACKET() antlr.TerminalNode

func (s *ExprContext) LPAREN() antlr.TerminalNode

func (s *ExprContext) LSHIFT() antlr.TerminalNode

func (s *ExprContext) LT() antlr.TerminalNode

func (s *ExprContext) LTE() antlr.TerminalNode

func (s *ExprContext) MINUS() antlr.TerminalNode

func (s *ExprContext) NEQ() antlr.TerminalNode

func (s *ExprContext) NIL_COALESCE() antlr.TerminalNode

func (s *ExprContext) OK() antlr.TerminalNode

func (s *ExprContext) OR() antlr.TerminalNode

func (s *ExprContext) OVERFLOW_ADD() antlr.TerminalNode

func (s *ExprContext) OVERFLOW_MUL() antlr.TerminalNode

func (s *ExprContext) OVERFLOW_SUB() antlr.TerminalNode

func (s *ExprContext) PERCENT() antlr.TerminalNode

func (s *ExprContext) PIPE() antlr.TerminalNode

func (s *ExprContext) PLUS() antlr.TerminalNode

func (s *ExprContext) PostfixName() IPostfixNameContext

func (s *ExprContext) Primary() IPrimaryContext

func (s *ExprContext) PrimitiveType() IPrimitiveTypeContext

func (s *ExprContext) QUESTION() antlr.TerminalNode

func (s *ExprContext) RBRACKET() antlr.TerminalNode

func (s *ExprContext) RESULT() antlr.TerminalNode

func (s *ExprContext) RPAREN() antlr.TerminalNode

func (s *ExprContext) RSHIFT() antlr.TerminalNode

func (s *ExprContext) SLASH() antlr.TerminalNode

func (s *ExprContext) STAR() antlr.TerminalNode

func (s *ExprContext) TILDE() antlr.TerminalNode

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExprListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExprListContext() *ExprListContext

func NewExprListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprListContext

func (s *ExprListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExprListContext) AllCOMMA() []antlr.TerminalNode

func (s *ExprListContext) AllExpr() []IExprContext

func (s *ExprListContext) COMMA(i int) antlr.TerminalNode

func (s *ExprListContext) Expr(i int) IExprContext

func (s *ExprListContext) GetParser() antlr.Parser

func (s *ExprListContext) GetRuleContext() antlr.RuleContext

func (*ExprListContext) IsExprListContext()

func (s *ExprListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExprStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExprStmtContext() *ExprStmtContext

func NewExprStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprStmtContext

func (s *ExprStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExprStmtContext) Expr() IExprContext

func (s *ExprStmtContext) GetParser() antlr.Parser

func (s *ExprStmtContext) GetRuleContext() antlr.RuleContext

func (*ExprStmtContext) IsExprStmtContext()

func (s *ExprStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FallthroughStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFallthroughStmtContext() *FallthroughStmtContext

func NewFallthroughStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FallthroughStmtContext

func (s *FallthroughStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FallthroughStmtContext) FALLTHROUGH() antlr.TerminalNode

func (s *FallthroughStmtContext) GetParser() antlr.Parser

func (s *FallthroughStmtContext) GetRuleContext() antlr.RuleContext

func (*FallthroughStmtContext) IsFallthroughStmtContext()

func (s *FallthroughStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FileContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFileContext() *FileContext

func NewFileContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FileContext

func (s *FileContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FileContext) AllImportDecl() []IImportDeclContext

func (s *FileContext) AllTopLevelDecl() []ITopLevelDeclContext

func (s *FileContext) BuildDecl() IBuildDeclContext

func (s *FileContext) EOF() antlr.TerminalNode

func (s *FileContext) GetParser() antlr.Parser

func (s *FileContext) GetRuleContext() antlr.RuleContext

func (s *FileContext) ImportDecl(i int) IImportDeclContext

func (*FileContext) IsFileContext()

func (s *FileContext) PackageDecl() IPackageDeclContext

func (s *FileContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FileContext) TopLevelDecl(i int) ITopLevelDeclContext

type ForInStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyForInStmtContext() *ForInStmtContext

func NewForInStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForInStmtContext

func (s *ForInStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ForInStmtContext) Block() IBlockContext

func (s *ForInStmtContext) Expr() IExprContext

func (s *ForInStmtContext) FOR() antlr.TerminalNode

func (s *ForInStmtContext) GetParser() antlr.Parser

func (s *ForInStmtContext) GetRuleContext() antlr.RuleContext

func (s *ForInStmtContext) ID() antlr.TerminalNode

func (s *ForInStmtContext) IN() antlr.TerminalNode

func (*ForInStmtContext) IsForInStmtContext()

func (s *ForInStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FuncDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFuncDeclContext() *FuncDeclContext

func NewFuncDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncDeclContext

func (s *FuncDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FuncDeclContext) Block() IBlockContext

func (s *FuncDeclContext) FUNC() antlr.TerminalNode

func (s *FuncDeclContext) FuncQualifier() IFuncQualifierContext

func (s *FuncDeclContext) GenericParams() IGenericParamsContext

func (s *FuncDeclContext) GetParser() antlr.Parser

func (s *FuncDeclContext) GetRuleContext() antlr.RuleContext

func (s *FuncDeclContext) ID() antlr.TerminalNode

func (*FuncDeclContext) IsFuncDeclContext()

func (s *FuncDeclContext) LPAREN() antlr.TerminalNode

func (s *FuncDeclContext) ParamList() IParamListContext

func (s *FuncDeclContext) RPAREN() antlr.TerminalNode

func (s *FuncDeclContext) ReturnType() IReturnTypeContext

func (s *FuncDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FuncQualifierContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFuncQualifierContext() *FuncQualifierContext

func NewFuncQualifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncQualifierContext

func (s *FuncQualifierContext) ASYNC() antlr.TerminalNode

func (s *FuncQualifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FuncQualifierContext) GPU() antlr.TerminalNode

func (s *FuncQualifierContext) GetParser() antlr.Parser

func (s *FuncQualifierContext) GetRuleContext() antlr.RuleContext

func (*FuncQualifierContext) IsFuncQualifierContext()

func (s *FuncQualifierContext) PROCESS() antlr.TerminalNode

func (s *FuncQualifierContext) THREAD() antlr.TerminalNode

func (s *FuncQualifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FuncTypeParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFuncTypeParamContext() *FuncTypeParamContext

func NewFuncTypeParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncTypeParamContext

func (s *FuncTypeParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FuncTypeParamContext) GetParser() antlr.Parser

func (s *FuncTypeParamContext) GetRuleContext() antlr.RuleContext

func (*FuncTypeParamContext) IsFuncTypeParamContext()

func (s *FuncTypeParamContext) MUT() antlr.TerminalNode

func (s *FuncTypeParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FuncTypeParamContext) Type_() ITypeContext

type FuncTypeParamListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFuncTypeParamListContext() *FuncTypeParamListContext

func NewFuncTypeParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncTypeParamListContext

func (s *FuncTypeParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FuncTypeParamListContext) AllCOMMA() []antlr.TerminalNode

func (s *FuncTypeParamListContext) AllFuncTypeParam() []IFuncTypeParamContext

func (s *FuncTypeParamListContext) COMMA(i int) antlr.TerminalNode

func (s *FuncTypeParamListContext) FuncTypeParam(i int) IFuncTypeParamContext

func (s *FuncTypeParamListContext) GetParser() antlr.Parser

func (s *FuncTypeParamListContext) GetRuleContext() antlr.RuleContext

func (*FuncTypeParamListContext) IsFuncTypeParamListContext()

func (s *FuncTypeParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type GenericParamsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericParamsContext() *GenericParamsContext

func NewGenericParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericParamsContext

func (s *GenericParamsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericParamsContext) AllCOMMA() []antlr.TerminalNode

func (s *GenericParamsContext) AllTypeParam() []ITypeParamContext

func (s *GenericParamsContext) COMMA(i int) antlr.TerminalNode

func (s *GenericParamsContext) GT() antlr.TerminalNode

func (s *GenericParamsContext) GetParser() antlr.Parser

func (s *GenericParamsContext) GetRuleContext() antlr.RuleContext

func (*GenericParamsContext) IsGenericParamsContext()

func (s *GenericParamsContext) LT() antlr.TerminalNode

func (s *GenericParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *GenericParamsContext) TypeParam(i int) ITypeParamContext

type IAnonFuncExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNC() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	Block() IBlockContext
	ParamList() IParamListContext
	ReturnType() IReturnTypeContext

	// IsAnonFuncExprContext differentiates from other interfaces.
	IsAnonFuncExprContext()
}
    IAnonFuncExprContext is an interface to support dynamic dispatch.

type IArgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expr() IExprContext
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsArgContext differentiates from other interfaces.
	IsArgContext()
}
    IArgContext is an interface to support dynamic dispatch.

type IArgListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllArg() []IArgContext
	Arg(i int) IArgContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsArgListContext differentiates from other interfaces.
	IsArgListContext()
}
    IArgListContext is an interface to support dynamic dispatch.

type IArrayConstructExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	Type_() ITypeContext
	RBRACKET() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	AllCOLON() []antlr.TerminalNode
	COLON(i int) antlr.TerminalNode
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	COMMA() antlr.TerminalNode

	// IsArrayConstructExprContext differentiates from other interfaces.
	IsArrayConstructExprContext()
}
    IArrayConstructExprContext is an interface to support dynamic dispatch.

type IArrayLiteralExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	RBRACKET() antlr.TerminalNode
	ExprList() IExprListContext
	COMMA() antlr.TerminalNode

	// IsArrayLiteralExprContext differentiates from other interfaces.
	IsArrayLiteralExprContext()
}
    IArrayLiteralExprContext is an interface to support dynamic dispatch.

type IAssignStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Lvalue() ILvalueContext
	ASSIGN() antlr.TerminalNode
	Expr() IExprContext

	// IsAssignStmtContext differentiates from other interfaces.
	IsAssignStmtContext()
}
    IAssignStmtContext is an interface to support dynamic dispatch.

type IBindingKwContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WEAK() antlr.TerminalNode
	LET() antlr.TerminalNode
	VAR() antlr.TerminalNode

	// IsBindingKwContext differentiates from other interfaces.
	IsBindingKwContext()
}
    IBindingKwContext is an interface to support dynamic dispatch.

type IBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllStmt() []IStmtContext
	Stmt(i int) IStmtContext

	// IsBlockContext differentiates from other interfaces.
	IsBlockContext()
}
    IBlockContext is an interface to support dynamic dispatch.

type IBreakStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BREAK() antlr.TerminalNode

	// IsBreakStmtContext differentiates from other interfaces.
	IsBreakStmtContext()
}
    IBreakStmtContext is an interface to support dynamic dispatch.

type IBuildDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BUILD() antlr.TerminalNode
	AllBuildTag() []IBuildTagContext
	BuildTag(i int) IBuildTagContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsBuildDeclContext differentiates from other interfaces.
	IsBuildDeclContext()
}
    IBuildDeclContext is an interface to support dynamic dispatch.

type IBuildTagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsBuildTagContext differentiates from other interfaces.
	IsBuildTagContext()
}
    IBuildTagContext is an interface to support dynamic dispatch.

type ICasePatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Literal() ILiteralContext
	DOT() antlr.TerminalNode
	ID() antlr.TerminalNode
	OK() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	LET() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ERR() antlr.TerminalNode

	// IsCasePatternContext differentiates from other interfaces.
	IsCasePatternContext()
}
    ICasePatternContext is an interface to support dynamic dispatch.

type ICasePatternListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCasePattern() []ICasePatternContext
	CasePattern(i int) ICasePatternContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsCasePatternListContext differentiates from other interfaces.
	IsCasePatternListContext()
}
    ICasePatternListContext is an interface to support dynamic dispatch.

type IClassDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CLASS() antlr.TerminalNode
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	COLON() antlr.TerminalNode
	AllClassMember() []IClassMemberContext
	ClassMember(i int) IClassMemberContext

	// IsClassDeclContext differentiates from other interfaces.
	IsClassDeclContext()
}
    IClassDeclContext is an interface to support dynamic dispatch.

type IClassFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext
	LET() antlr.TerminalNode
	VAR() antlr.TerminalNode

	// IsClassFieldContext differentiates from other interfaces.
	IsClassFieldContext()
}
    IClassFieldContext is an interface to support dynamic dispatch.

type IClassMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ClassField() IClassFieldContext
	NativeFuncDecl() INativeFuncDeclContext

	// IsClassMemberContext differentiates from other interfaces.
	IsClassMemberContext()
}
    IClassMemberContext is an interface to support dynamic dispatch.

type ICompoundAssignStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Lvalue() ILvalueContext
	CompoundOp() ICompoundOpContext
	Expr() IExprContext

	// IsCompoundAssignStmtContext differentiates from other interfaces.
	IsCompoundAssignStmtContext()
}
    ICompoundAssignStmtContext is an interface to support dynamic dispatch.

type ICompoundOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PLUS_ASSIGN() antlr.TerminalNode
	MINUS_ASSIGN() antlr.TerminalNode
	STAR_ASSIGN() antlr.TerminalNode
	SLASH_ASSIGN() antlr.TerminalNode
	MOD_ASSIGN() antlr.TerminalNode

	// IsCompoundOpContext differentiates from other interfaces.
	IsCompoundOpContext()
}
    ICompoundOpContext is an interface to support dynamic dispatch.

type IContinueStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CONTINUE() antlr.TerminalNode

	// IsContinueStmtContext differentiates from other interfaces.
	IsContinueStmtContext()
}
    IContinueStmtContext is an interface to support dynamic dispatch.

type IDefaultCaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEFAULT() antlr.TerminalNode
	COLON() antlr.TerminalNode
	AllStmt() []IStmtContext
	Stmt(i int) IStmtContext

	// IsDefaultCaseContext differentiates from other interfaces.
	IsDefaultCaseContext()
}
    IDefaultCaseContext is an interface to support dynamic dispatch.

type IDeferStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEFER() antlr.TerminalNode
	Expr() IExprContext

	// IsDeferStmtContext differentiates from other interfaces.
	IsDeferStmtContext()
}
    IDeferStmtContext is an interface to support dynamic dispatch.

type IDictEntryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	COLON() antlr.TerminalNode

	// IsDictEntryContext differentiates from other interfaces.
	IsDictEntryContext()
}
    IDictEntryContext is an interface to support dynamic dispatch.

type IDictLiteralExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	AllDictEntry() []IDictEntryContext
	DictEntry(i int) IDictEntryContext
	RBRACKET() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsDictLiteralExprContext differentiates from other interfaces.
	IsDictLiteralExprContext()
}
    IDictLiteralExprContext is an interface to support dynamic dispatch.

type IElseClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ELSE() antlr.TerminalNode
	Block() IBlockContext

	// IsElseClauseContext differentiates from other interfaces.
	IsElseClauseContext()
}
    IElseClauseContext is an interface to support dynamic dispatch.

type IElseIfClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ELSE() antlr.TerminalNode
	IF() antlr.TerminalNode
	IfCondition() IIfConditionContext
	Block() IBlockContext

	// IsElseIfClauseContext differentiates from other interfaces.
	IsElseIfClauseContext()
}
    IElseIfClauseContext is an interface to support dynamic dispatch.

type IEnumCaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	Literal() ILiteralContext

	// IsEnumCaseContext differentiates from other interfaces.
	IsEnumCaseContext()
}
    IEnumCaseContext is an interface to support dynamic dispatch.

type IEnumCaseDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	AllEnumCase() []IEnumCaseContext
	EnumCase(i int) IEnumCaseContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsEnumCaseDeclContext differentiates from other interfaces.
	IsEnumCaseDeclContext()
}
    IEnumCaseDeclContext is an interface to support dynamic dispatch.

type IEnumDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ENUM() antlr.TerminalNode
	ID() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	COLON() antlr.TerminalNode
	EnumRawType() IEnumRawTypeContext
	AllEnumCaseDecl() []IEnumCaseDeclContext
	EnumCaseDecl(i int) IEnumCaseDeclContext

	// IsEnumDeclContext differentiates from other interfaces.
	IsEnumDeclContext()
}
    IEnumDeclContext is an interface to support dynamic dispatch.

type IEnumRawTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT() antlr.TerminalNode
	STRING() antlr.TerminalNode

	// IsEnumRawTypeContext differentiates from other interfaces.
	IsEnumRawTypeContext()
}
    IEnumRawTypeContext is an interface to support dynamic dispatch.

type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MINUS() antlr.TerminalNode
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	BANG() antlr.TerminalNode
	TILDE() antlr.TerminalNode
	AMP() antlr.TerminalNode
	PrimitiveType() IPrimitiveTypeContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	RESULT() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	OK() antlr.TerminalNode
	ERR() antlr.TerminalNode
	Primary() IPrimaryContext
	LSHIFT() antlr.TerminalNode
	RSHIFT() antlr.TerminalNode
	CARET() antlr.TerminalNode
	PIPE() antlr.TerminalNode
	STAR() antlr.TerminalNode
	SLASH() antlr.TerminalNode
	PERCENT() antlr.TerminalNode
	OVERFLOW_MUL() antlr.TerminalNode
	PLUS() antlr.TerminalNode
	OVERFLOW_ADD() antlr.TerminalNode
	OVERFLOW_SUB() antlr.TerminalNode
	ELLIPSIS() antlr.TerminalNode
	HALF_OPEN_RANGE() antlr.TerminalNode
	NIL_COALESCE() antlr.TerminalNode
	EQ() antlr.TerminalNode
	NEQ() antlr.TerminalNode
	LT() antlr.TerminalNode
	GT() antlr.TerminalNode
	LTE() antlr.TerminalNode
	GTE() antlr.TerminalNode
	IDENTITY_EQ() antlr.TerminalNode
	IDENTITY_NEQ() antlr.TerminalNode
	AND() antlr.TerminalNode
	OR() antlr.TerminalNode
	QUESTION() antlr.TerminalNode
	COLON() antlr.TerminalNode
	DOT() antlr.TerminalNode
	PostfixName() IPostfixNameContext
	ArgList() IArgListContext
	LBRACKET() antlr.TerminalNode
	RBRACKET() antlr.TerminalNode

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}
    IExprContext is an interface to support dynamic dispatch.

type IExprListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsExprListContext differentiates from other interfaces.
	IsExprListContext()
}
    IExprListContext is an interface to support dynamic dispatch.

type IExprStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expr() IExprContext

	// IsExprStmtContext differentiates from other interfaces.
	IsExprStmtContext()
}
    IExprStmtContext is an interface to support dynamic dispatch.

type IFallthroughStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FALLTHROUGH() antlr.TerminalNode

	// IsFallthroughStmtContext differentiates from other interfaces.
	IsFallthroughStmtContext()
}
    IFallthroughStmtContext is an interface to support dynamic dispatch.

type IFileContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	PackageDecl() IPackageDeclContext
	BuildDecl() IBuildDeclContext
	AllImportDecl() []IImportDeclContext
	ImportDecl(i int) IImportDeclContext
	AllTopLevelDecl() []ITopLevelDeclContext
	TopLevelDecl(i int) ITopLevelDeclContext

	// IsFileContext differentiates from other interfaces.
	IsFileContext()
}
    IFileContext is an interface to support dynamic dispatch.

type IForInStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FOR() antlr.TerminalNode
	ID() antlr.TerminalNode
	IN() antlr.TerminalNode
	Expr() IExprContext
	Block() IBlockContext

	// IsForInStmtContext differentiates from other interfaces.
	IsForInStmtContext()
}
    IForInStmtContext is an interface to support dynamic dispatch.

type IFuncDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNC() antlr.TerminalNode
	ID() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	Block() IBlockContext
	GenericParams() IGenericParamsContext
	ParamList() IParamListContext
	FuncQualifier() IFuncQualifierContext
	ReturnType() IReturnTypeContext

	// IsFuncDeclContext differentiates from other interfaces.
	IsFuncDeclContext()
}
    IFuncDeclContext is an interface to support dynamic dispatch.

type IFuncQualifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASYNC() antlr.TerminalNode
	THREAD() antlr.TerminalNode
	PROCESS() antlr.TerminalNode
	GPU() antlr.TerminalNode

	// IsFuncQualifierContext differentiates from other interfaces.
	IsFuncQualifierContext()
}
    IFuncQualifierContext is an interface to support dynamic dispatch.

type IFuncTypeParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Type_() ITypeContext
	MUT() antlr.TerminalNode

	// IsFuncTypeParamContext differentiates from other interfaces.
	IsFuncTypeParamContext()
}
    IFuncTypeParamContext is an interface to support dynamic dispatch.

type IFuncTypeParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllFuncTypeParam() []IFuncTypeParamContext
	FuncTypeParam(i int) IFuncTypeParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFuncTypeParamListContext differentiates from other interfaces.
	IsFuncTypeParamListContext()
}
    IFuncTypeParamListContext is an interface to support dynamic dispatch.

type IGenericParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LT() antlr.TerminalNode
	AllTypeParam() []ITypeParamContext
	TypeParam(i int) ITypeParamContext
	GT() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsGenericParamsContext differentiates from other interfaces.
	IsGenericParamsContext()
}
    IGenericParamsContext is an interface to support dynamic dispatch.

type IIfConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LET() antlr.TerminalNode
	ID() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	Expr() IExprContext

	// IsIfConditionContext differentiates from other interfaces.
	IsIfConditionContext()
}
    IIfConditionContext is an interface to support dynamic dispatch.

type IIfStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IF() antlr.TerminalNode
	IfCondition() IIfConditionContext
	Block() IBlockContext
	AllElseIfClause() []IElseIfClauseContext
	ElseIfClause(i int) IElseIfClauseContext
	ElseClause() IElseClauseContext

	// IsIfStmtContext differentiates from other interfaces.
	IsIfStmtContext()
}
    IIfStmtContext is an interface to support dynamic dispatch.

type IImportDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IMPORT() antlr.TerminalNode
	AllSTRING_LIT() []antlr.TerminalNode
	STRING_LIT(i int) antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode

	// IsImportDeclContext differentiates from other interfaces.
	IsImportDeclContext()
}
    IImportDeclContext is an interface to support dynamic dispatch.

type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEC_INT_LIT() antlr.TerminalNode
	HEX_INT_LIT() antlr.TerminalNode
	BIN_INT_LIT() antlr.TerminalNode
	OCT_INT_LIT() antlr.TerminalNode
	DEC_FLOAT_LIT() antlr.TerminalNode
	HEX_FLOAT_LIT() antlr.TerminalNode
	STRING_LIT() antlr.TerminalNode
	RAW_STRING_LIT() antlr.TerminalNode
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode
	NIL() antlr.TerminalNode

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}
    ILiteralContext is an interface to support dynamic dispatch.

type ILvalueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	Lvalue() ILvalueContext
	DOT() antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	Expr() IExprContext
	RBRACKET() antlr.TerminalNode

	// IsLvalueContext differentiates from other interfaces.
	IsLvalueContext()
}
    ILvalueContext is an interface to support dynamic dispatch.

type INativeFuncDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNC() antlr.TerminalNode
	ID() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	NativeParamList() INativeParamListContext
	ReturnType() IReturnTypeContext

	// IsNativeFuncDeclContext differentiates from other interfaces.
	IsNativeFuncDeclContext()
}
    INativeFuncDeclContext is an interface to support dynamic dispatch.

type INativeParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext
	MUT() antlr.TerminalNode
	ELLIPSIS() antlr.TerminalNode

	// IsNativeParamContext differentiates from other interfaces.
	IsNativeParamContext()
}
    INativeParamContext is an interface to support dynamic dispatch.

type INativeParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNativeParam() []INativeParamContext
	NativeParam(i int) INativeParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsNativeParamListContext differentiates from other interfaces.
	IsNativeParamListContext()
}
    INativeParamListContext is an interface to support dynamic dispatch.

type IPackageDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PACKAGE() antlr.TerminalNode
	ID() antlr.TerminalNode

	// IsPackageDeclContext differentiates from other interfaces.
	IsPackageDeclContext()
}
    IPackageDeclContext is an interface to support dynamic dispatch.

type IParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext
	MUT() antlr.TerminalNode

	// IsParamContext differentiates from other interfaces.
	IsParamContext()
}
    IParamContext is an interface to support dynamic dispatch.

type IParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllParam() []IParamContext
	Param(i int) IParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsParamListContext differentiates from other interfaces.
	IsParamListContext()
}
    IParamListContext is an interface to support dynamic dispatch.

type IPostfixNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	ANY() antlr.TerminalNode

	// IsPostfixNameContext differentiates from other interfaces.
	IsPostfixNameContext()
}
    IPostfixNameContext is an interface to support dynamic dispatch.

type IPrimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Literal() ILiteralContext
	StructLiteralExpr() IStructLiteralExprContext
	ID() antlr.TerminalNode
	GPU() antlr.TerminalNode
	ASYNC() antlr.TerminalNode
	THREAD() antlr.TerminalNode
	PROCESS() antlr.TerminalNode
	DOT() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	Expr() IExprContext
	TupleExpr() ITupleExprContext
	ArrayLiteralExpr() IArrayLiteralExprContext
	ArrayConstructExpr() IArrayConstructExprContext
	DictLiteralExpr() IDictLiteralExprContext
	AnonFuncExpr() IAnonFuncExprContext

	// IsPrimaryContext differentiates from other interfaces.
	IsPrimaryContext()
}
    IPrimaryContext is an interface to support dynamic dispatch.

type IPrimitiveTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT() antlr.TerminalNode
	INT8() antlr.TerminalNode
	INT16() antlr.TerminalNode
	INT32() antlr.TerminalNode
	INT64() antlr.TerminalNode
	UINT() antlr.TerminalNode
	UINT8() antlr.TerminalNode
	UINT16() antlr.TerminalNode
	UINT32() antlr.TerminalNode
	UINT64() antlr.TerminalNode
	FLOAT() antlr.TerminalNode
	DOUBLE() antlr.TerminalNode
	BOOL() antlr.TerminalNode
	STRING() antlr.TerminalNode
	CHAR() antlr.TerminalNode
	VOID() antlr.TerminalNode

	// IsPrimitiveTypeContext differentiates from other interfaces.
	IsPrimitiveTypeContext()
}
    IPrimitiveTypeContext is an interface to support dynamic dispatch.

type IReturnStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RETURN() antlr.TerminalNode
	Expr() IExprContext

	// IsReturnStmtContext differentiates from other interfaces.
	IsReturnStmtContext()
}
    IReturnStmtContext is an interface to support dynamic dispatch.

type IReturnTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ARROW() antlr.TerminalNode
	Type_() ITypeContext

	// IsReturnTypeContext differentiates from other interfaces.
	IsReturnTypeContext()
}
    IReturnTypeContext is an interface to support dynamic dispatch.

type IStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VarDeclStmt() IVarDeclStmtContext
	AssignStmt() IAssignStmtContext
	CompoundAssignStmt() ICompoundAssignStmtContext
	IfStmt() IIfStmtContext
	SwitchStmt() ISwitchStmtContext
	ForInStmt() IForInStmtContext
	WhileStmt() IWhileStmtContext
	BreakStmt() IBreakStmtContext
	ContinueStmt() IContinueStmtContext
	FallthroughStmt() IFallthroughStmtContext
	ReturnStmt() IReturnStmtContext
	DeferStmt() IDeferStmtContext
	ExprStmt() IExprStmtContext

	// IsStmtContext differentiates from other interfaces.
	IsStmtContext()
}
    IStmtContext is an interface to support dynamic dispatch.

type IStructDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRUCT() antlr.TerminalNode
	ID() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	GenericParams() IGenericParamsContext
	AllStructField() []IStructFieldContext
	StructField(i int) IStructFieldContext

	// IsStructDeclContext differentiates from other interfaces.
	IsStructDeclContext()
}
    IStructDeclContext is an interface to support dynamic dispatch.

type IStructFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext
	LET() antlr.TerminalNode
	VAR() antlr.TerminalNode

	// IsStructFieldContext differentiates from other interfaces.
	IsStructFieldContext()
}
    IStructFieldContext is an interface to support dynamic dispatch.

type IStructFieldInitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Expr() IExprContext

	// IsStructFieldInitContext differentiates from other interfaces.
	IsStructFieldInitContext()
}
    IStructFieldInitContext is an interface to support dynamic dispatch.

type IStructLiteralExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllStructFieldInit() []IStructFieldInitContext
	StructFieldInit(i int) IStructFieldInitContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsStructLiteralExprContext differentiates from other interfaces.
	IsStructLiteralExprContext()
}
    IStructLiteralExprContext is an interface to support dynamic dispatch.

type ISwitchCaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	CasePatternList() ICasePatternListContext
	COLON() antlr.TerminalNode
	AllStmt() []IStmtContext
	Stmt(i int) IStmtContext

	// IsSwitchCaseContext differentiates from other interfaces.
	IsSwitchCaseContext()
}
    ISwitchCaseContext is an interface to support dynamic dispatch.

type ISwitchStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SWITCH() antlr.TerminalNode
	Expr() IExprContext
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllSwitchCase() []ISwitchCaseContext
	SwitchCase(i int) ISwitchCaseContext
	DefaultCase() IDefaultCaseContext

	// IsSwitchStmtContext differentiates from other interfaces.
	IsSwitchStmtContext()
}
    ISwitchStmtContext is an interface to support dynamic dispatch.

type ITopLevelDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FuncDecl() IFuncDeclContext
	StructDecl() IStructDeclContext
	ClassDecl() IClassDeclContext
	EnumDecl() IEnumDeclContext
	TypeAliasDecl() ITypeAliasDeclContext
	VarDeclStmt() IVarDeclStmtContext

	// IsTopLevelDeclContext differentiates from other interfaces.
	IsTopLevelDeclContext()
}
    ITopLevelDeclContext is an interface to support dynamic dispatch.

type ITupleBindContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTupleBindContext differentiates from other interfaces.
	IsTupleBindContext()
}
    ITupleBindContext is an interface to support dynamic dispatch.

type ITupleElementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expr() IExprContext
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsTupleElementContext differentiates from other interfaces.
	IsTupleElementContext()
}
    ITupleElementContext is an interface to support dynamic dispatch.

type ITupleExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	AllTupleElement() []ITupleElementContext
	TupleElement(i int) ITupleElementContext
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTupleExprContext differentiates from other interfaces.
	IsTupleExprContext()
}
    ITupleExprContext is an interface to support dynamic dispatch.

type ITupleTypeElemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Type_() ITypeContext
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsTupleTypeElemContext differentiates from other interfaces.
	IsTupleTypeElemContext()
}
    ITupleTypeElemContext is an interface to support dynamic dispatch.

type ITypeAliasDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE() antlr.TerminalNode
	ID() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	Type_() ITypeContext

	// IsTypeAliasDeclContext differentiates from other interfaces.
	IsTypeAliasDeclContext()
}
    ITypeAliasDeclContext is an interface to support dynamic dispatch.

type ITypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PrimitiveType() IPrimitiveTypeContext
	ID() antlr.TerminalNode
	LT() antlr.TerminalNode
	AllType_() []ITypeContext
	Type_(i int) ITypeContext
	GT() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	RBRACKET() antlr.TerminalNode
	COLON() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AllTupleTypeElem() []ITupleTypeElemContext
	TupleTypeElem(i int) ITupleTypeElemContext
	RPAREN() antlr.TerminalNode
	FUNC() antlr.TerminalNode
	FuncTypeParamList() IFuncTypeParamListContext
	ARROW() antlr.TerminalNode
	CHANNEL() antlr.TerminalNode
	ANY() antlr.TerminalNode
	MUT() antlr.TerminalNode
	VOID() antlr.TerminalNode
	OPAQUE() antlr.TerminalNode
	RESULT() antlr.TerminalNode
	QUESTION() antlr.TerminalNode

	// IsTypeContext differentiates from other interfaces.
	IsTypeContext()
}
    ITypeContext is an interface to support dynamic dispatch.

type ITypeParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsTypeParamContext differentiates from other interfaces.
	IsTypeParamContext()
}
    ITypeParamContext is an interface to support dynamic dispatch.

type IVarDeclStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BindingKw() IBindingKwContext
	TupleBind() ITupleBindContext
	ASSIGN() antlr.TerminalNode
	Expr() IExprContext
	COLON() antlr.TerminalNode
	Type_() ITypeContext
	ID() antlr.TerminalNode

	// IsVarDeclStmtContext differentiates from other interfaces.
	IsVarDeclStmtContext()
}
    IVarDeclStmtContext is an interface to support dynamic dispatch.

type IWhileStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHILE() antlr.TerminalNode
	Expr() IExprContext
	Block() IBlockContext

	// IsWhileStmtContext differentiates from other interfaces.
	IsWhileStmtContext()
}
    IWhileStmtContext is an interface to support dynamic dispatch.

type IfConditionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyIfConditionContext() *IfConditionContext

func NewIfConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfConditionContext

func (s *IfConditionContext) ASSIGN() antlr.TerminalNode

func (s *IfConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IfConditionContext) Expr() IExprContext

func (s *IfConditionContext) GetParser() antlr.Parser

func (s *IfConditionContext) GetRuleContext() antlr.RuleContext

func (s *IfConditionContext) ID() antlr.TerminalNode

func (*IfConditionContext) IsIfConditionContext()

func (s *IfConditionContext) LET() antlr.TerminalNode

func (s *IfConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type IfStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyIfStmtContext() *IfStmtContext

func NewIfStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfStmtContext

func (s *IfStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IfStmtContext) AllElseIfClause() []IElseIfClauseContext

func (s *IfStmtContext) Block() IBlockContext

func (s *IfStmtContext) ElseClause() IElseClauseContext

func (s *IfStmtContext) ElseIfClause(i int) IElseIfClauseContext

func (s *IfStmtContext) GetParser() antlr.Parser

func (s *IfStmtContext) GetRuleContext() antlr.RuleContext

func (s *IfStmtContext) IF() antlr.TerminalNode

func (s *IfStmtContext) IfCondition() IIfConditionContext

func (*IfStmtContext) IsIfStmtContext()

func (s *IfStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ImportDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyImportDeclContext() *ImportDeclContext

func NewImportDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportDeclContext

func (s *ImportDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ImportDeclContext) AllSTRING_LIT() []antlr.TerminalNode

func (s *ImportDeclContext) GetParser() antlr.Parser

func (s *ImportDeclContext) GetRuleContext() antlr.RuleContext

func (s *ImportDeclContext) IMPORT() antlr.TerminalNode

func (*ImportDeclContext) IsImportDeclContext()

func (s *ImportDeclContext) LPAREN() antlr.TerminalNode

func (s *ImportDeclContext) RPAREN() antlr.TerminalNode

func (s *ImportDeclContext) STRING_LIT(i int) antlr.TerminalNode

func (s *ImportDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type LiteralContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLiteralContext() *LiteralContext

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext

func (s *LiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LiteralContext) BIN_INT_LIT() antlr.TerminalNode

func (s *LiteralContext) DEC_FLOAT_LIT() antlr.TerminalNode

func (s *LiteralContext) DEC_INT_LIT() antlr.TerminalNode

func (s *LiteralContext) FALSE() antlr.TerminalNode

func (s *LiteralContext) GetParser() antlr.Parser

func (s *LiteralContext) GetRuleContext() antlr.RuleContext

func (s *LiteralContext) HEX_FLOAT_LIT() antlr.TerminalNode

func (s *LiteralContext) HEX_INT_LIT() antlr.TerminalNode

func (*LiteralContext) IsLiteralContext()

func (s *LiteralContext) NIL() antlr.TerminalNode

func (s *LiteralContext) OCT_INT_LIT() antlr.TerminalNode

func (s *LiteralContext) RAW_STRING_LIT() antlr.TerminalNode

func (s *LiteralContext) STRING_LIT() antlr.TerminalNode

func (s *LiteralContext) TRUE() antlr.TerminalNode

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type LvalueContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLvalueContext() *LvalueContext

func NewLvalueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LvalueContext

func (s *LvalueContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LvalueContext) DOT() antlr.TerminalNode

func (s *LvalueContext) Expr() IExprContext

func (s *LvalueContext) GetParser() antlr.Parser

func (s *LvalueContext) GetRuleContext() antlr.RuleContext

func (s *LvalueContext) ID() antlr.TerminalNode

func (*LvalueContext) IsLvalueContext()

func (s *LvalueContext) LBRACKET() antlr.TerminalNode

func (s *LvalueContext) Lvalue() ILvalueContext

func (s *LvalueContext) RBRACKET() antlr.TerminalNode

func (s *LvalueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type NativeFuncDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyNativeFuncDeclContext() *NativeFuncDeclContext

func NewNativeFuncDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NativeFuncDeclContext

func (s *NativeFuncDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *NativeFuncDeclContext) FUNC() antlr.TerminalNode

func (s *NativeFuncDeclContext) GetParser() antlr.Parser

func (s *NativeFuncDeclContext) GetRuleContext() antlr.RuleContext

func (s *NativeFuncDeclContext) ID() antlr.TerminalNode

func (*NativeFuncDeclContext) IsNativeFuncDeclContext()

func (s *NativeFuncDeclContext) LPAREN() antlr.TerminalNode

func (s *NativeFuncDeclContext) NativeParamList() INativeParamListContext

func (s *NativeFuncDeclContext) RPAREN() antlr.TerminalNode

func (s *NativeFuncDeclContext) ReturnType() IReturnTypeContext

func (s *NativeFuncDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type NativeParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyNativeParamContext() *NativeParamContext

func NewNativeParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NativeParamContext

func (s *NativeParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *NativeParamContext) COLON() antlr.TerminalNode

func (s *NativeParamContext) ELLIPSIS() antlr.TerminalNode

func (s *NativeParamContext) GetParser() antlr.Parser

func (s *NativeParamContext) GetRuleContext() antlr.RuleContext

func (s *NativeParamContext) ID() antlr.TerminalNode

func (*NativeParamContext) IsNativeParamContext()

func (s *NativeParamContext) MUT() antlr.TerminalNode

func (s *NativeParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *NativeParamContext) Type_() ITypeContext

type NativeParamListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyNativeParamListContext() *NativeParamListContext

func NewNativeParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NativeParamListContext

func (s *NativeParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *NativeParamListContext) AllCOMMA() []antlr.TerminalNode

func (s *NativeParamListContext) AllNativeParam() []INativeParamContext

func (s *NativeParamListContext) COMMA(i int) antlr.TerminalNode

func (s *NativeParamListContext) GetParser() antlr.Parser

func (s *NativeParamListContext) GetRuleContext() antlr.RuleContext

func (*NativeParamListContext) IsNativeParamListContext()

func (s *NativeParamListContext) NativeParam(i int) INativeParamContext

func (s *NativeParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PackageDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPackageDeclContext() *PackageDeclContext

func NewPackageDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PackageDeclContext

func (s *PackageDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PackageDeclContext) GetParser() antlr.Parser

func (s *PackageDeclContext) GetRuleContext() antlr.RuleContext

func (s *PackageDeclContext) ID() antlr.TerminalNode

func (*PackageDeclContext) IsPackageDeclContext()

func (s *PackageDeclContext) PACKAGE() antlr.TerminalNode

func (s *PackageDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyParamContext() *ParamContext

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext

func (s *ParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParamContext) COLON() antlr.TerminalNode

func (s *ParamContext) GetParser() antlr.Parser

func (s *ParamContext) GetRuleContext() antlr.RuleContext

func (s *ParamContext) ID() antlr.TerminalNode

func (*ParamContext) IsParamContext()

func (s *ParamContext) MUT() antlr.TerminalNode

func (s *ParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ParamContext) Type_() ITypeContext

type ParamListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyParamListContext() *ParamListContext

func NewParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamListContext

func (s *ParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParamListContext) AllCOMMA() []antlr.TerminalNode

func (s *ParamListContext) AllParam() []IParamContext

func (s *ParamListContext) COMMA(i int) antlr.TerminalNode

func (s *ParamListContext) GetParser() antlr.Parser

func (s *ParamListContext) GetRuleContext() antlr.RuleContext

func (*ParamListContext) IsParamListContext()

func (s *ParamListContext) Param(i int) IParamContext

func (s *ParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PostfixNameContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPostfixNameContext() *PostfixNameContext

func NewPostfixNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PostfixNameContext

func (s *PostfixNameContext) ANY() antlr.TerminalNode

func (s *PostfixNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PostfixNameContext) GetParser() antlr.Parser

func (s *PostfixNameContext) GetRuleContext() antlr.RuleContext

func (s *PostfixNameContext) ID() antlr.TerminalNode

func (*PostfixNameContext) IsPostfixNameContext()

func (s *PostfixNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PrimaryContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPrimaryContext() *PrimaryContext

func NewPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryContext

func (s *PrimaryContext) ASYNC() antlr.TerminalNode

func (s *PrimaryContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PrimaryContext) AnonFuncExpr() IAnonFuncExprContext

func (s *PrimaryContext) ArrayConstructExpr() IArrayConstructExprContext

func (s *PrimaryContext) ArrayLiteralExpr() IArrayLiteralExprContext

func (s *PrimaryContext) DOT() antlr.TerminalNode

func (s *PrimaryContext) DictLiteralExpr() IDictLiteralExprContext

func (s *PrimaryContext) Expr() IExprContext

func (s *PrimaryContext) GPU() antlr.TerminalNode

func (s *PrimaryContext) GetParser() antlr.Parser

func (s *PrimaryContext) GetRuleContext() antlr.RuleContext

func (s *PrimaryContext) ID() antlr.TerminalNode

func (*PrimaryContext) IsPrimaryContext()

func (s *PrimaryContext) LPAREN() antlr.TerminalNode

func (s *PrimaryContext) Literal() ILiteralContext

func (s *PrimaryContext) PROCESS() antlr.TerminalNode

func (s *PrimaryContext) RPAREN() antlr.TerminalNode

func (s *PrimaryContext) StructLiteralExpr() IStructLiteralExprContext

func (s *PrimaryContext) THREAD() antlr.TerminalNode

func (s *PrimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *PrimaryContext) TupleExpr() ITupleExprContext

type PrimitiveTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPrimitiveTypeContext() *PrimitiveTypeContext

func NewPrimitiveTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimitiveTypeContext

func (s *PrimitiveTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PrimitiveTypeContext) BOOL() antlr.TerminalNode

func (s *PrimitiveTypeContext) CHAR() antlr.TerminalNode

func (s *PrimitiveTypeContext) DOUBLE() antlr.TerminalNode

func (s *PrimitiveTypeContext) FLOAT() antlr.TerminalNode

func (s *PrimitiveTypeContext) GetParser() antlr.Parser

func (s *PrimitiveTypeContext) GetRuleContext() antlr.RuleContext

func (s *PrimitiveTypeContext) INT() antlr.TerminalNode

func (s *PrimitiveTypeContext) INT16() antlr.TerminalNode

func (s *PrimitiveTypeContext) INT32() antlr.TerminalNode

func (s *PrimitiveTypeContext) INT64() antlr.TerminalNode

func (s *PrimitiveTypeContext) INT8() antlr.TerminalNode

func (*PrimitiveTypeContext) IsPrimitiveTypeContext()

func (s *PrimitiveTypeContext) STRING() antlr.TerminalNode

func (s *PrimitiveTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *PrimitiveTypeContext) UINT() antlr.TerminalNode

func (s *PrimitiveTypeContext) UINT16() antlr.TerminalNode

func (s *PrimitiveTypeContext) UINT32() antlr.TerminalNode

func (s *PrimitiveTypeContext) UINT64() antlr.TerminalNode

func (s *PrimitiveTypeContext) UINT8() antlr.TerminalNode

func (s *PrimitiveTypeContext) VOID() antlr.TerminalNode

type ReturnStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyReturnStmtContext() *ReturnStmtContext

func NewReturnStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReturnStmtContext

func (s *ReturnStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ReturnStmtContext) Expr() IExprContext

func (s *ReturnStmtContext) GetParser() antlr.Parser

func (s *ReturnStmtContext) GetRuleContext() antlr.RuleContext

func (*ReturnStmtContext) IsReturnStmtContext()

func (s *ReturnStmtContext) RETURN() antlr.TerminalNode

func (s *ReturnStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ReturnTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyReturnTypeContext() *ReturnTypeContext

func NewReturnTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReturnTypeContext

func (s *ReturnTypeContext) ARROW() antlr.TerminalNode

func (s *ReturnTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ReturnTypeContext) GetParser() antlr.Parser

func (s *ReturnTypeContext) GetRuleContext() antlr.RuleContext

func (*ReturnTypeContext) IsReturnTypeContext()

func (s *ReturnTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ReturnTypeContext) Type_() ITypeContext

type StmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStmtContext() *StmtContext

func NewStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtContext

func (s *StmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StmtContext) AssignStmt() IAssignStmtContext

func (s *StmtContext) BreakStmt() IBreakStmtContext

func (s *StmtContext) CompoundAssignStmt() ICompoundAssignStmtContext

func (s *StmtContext) ContinueStmt() IContinueStmtContext

func (s *StmtContext) DeferStmt() IDeferStmtContext

func (s *StmtContext) ExprStmt() IExprStmtContext

func (s *StmtContext) FallthroughStmt() IFallthroughStmtContext

func (s *StmtContext) ForInStmt() IForInStmtContext

func (s *StmtContext) GetParser() antlr.Parser

func (s *StmtContext) GetRuleContext() antlr.RuleContext

func (s *StmtContext) IfStmt() IIfStmtContext

func (*StmtContext) IsStmtContext()

func (s *StmtContext) ReturnStmt() IReturnStmtContext

func (s *StmtContext) SwitchStmt() ISwitchStmtContext

func (s *StmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *StmtContext) VarDeclStmt() IVarDeclStmtContext

func (s *StmtContext) WhileStmt() IWhileStmtContext

type StructDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructDeclContext() *StructDeclContext

func NewStructDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructDeclContext

func (s *StructDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructDeclContext) AllStructField() []IStructFieldContext

func (s *StructDeclContext) GenericParams() IGenericParamsContext

func (s *StructDeclContext) GetParser() antlr.Parser

func (s *StructDeclContext) GetRuleContext() antlr.RuleContext

func (s *StructDeclContext) ID() antlr.TerminalNode

func (*StructDeclContext) IsStructDeclContext()

func (s *StructDeclContext) LBRACE() antlr.TerminalNode

func (s *StructDeclContext) RBRACE() antlr.TerminalNode

func (s *StructDeclContext) STRUCT() antlr.TerminalNode

func (s *StructDeclContext) StructField(i int) IStructFieldContext

func (s *StructDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type StructFieldContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructFieldContext() *StructFieldContext

func NewStructFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructFieldContext

func (s *StructFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructFieldContext) COLON() antlr.TerminalNode

func (s *StructFieldContext) GetParser() antlr.Parser

func (s *StructFieldContext) GetRuleContext() antlr.RuleContext

func (s *StructFieldContext) ID() antlr.TerminalNode

func (*StructFieldContext) IsStructFieldContext()

func (s *StructFieldContext) LET() antlr.TerminalNode

func (s *StructFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *StructFieldContext) Type_() ITypeContext

func (s *StructFieldContext) VAR() antlr.TerminalNode

type StructFieldInitContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructFieldInitContext() *StructFieldInitContext

func NewStructFieldInitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructFieldInitContext

func (s *StructFieldInitContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructFieldInitContext) COLON() antlr.TerminalNode

func (s *StructFieldInitContext) Expr() IExprContext

func (s *StructFieldInitContext) GetParser() antlr.Parser

func (s *StructFieldInitContext) GetRuleContext() antlr.RuleContext

func (s *StructFieldInitContext) ID() antlr.TerminalNode

func (*StructFieldInitContext) IsStructFieldInitContext()

func (s *StructFieldInitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type StructLiteralExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructLiteralExprContext() *StructLiteralExprContext

func NewStructLiteralExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructLiteralExprContext

func (s *StructLiteralExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructLiteralExprContext) AllCOMMA() []antlr.TerminalNode

func (s *StructLiteralExprContext) AllStructFieldInit() []IStructFieldInitContext

func (s *StructLiteralExprContext) COMMA(i int) antlr.TerminalNode

func (s *StructLiteralExprContext) GetParser() antlr.Parser

func (s *StructLiteralExprContext) GetRuleContext() antlr.RuleContext

func (s *StructLiteralExprContext) ID() antlr.TerminalNode

func (*StructLiteralExprContext) IsStructLiteralExprContext()

func (s *StructLiteralExprContext) LBRACE() antlr.TerminalNode

func (s *StructLiteralExprContext) RBRACE() antlr.TerminalNode

func (s *StructLiteralExprContext) StructFieldInit(i int) IStructFieldInitContext

func (s *StructLiteralExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SwitchCaseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchCaseContext() *SwitchCaseContext

func NewSwitchCaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchCaseContext

func (s *SwitchCaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchCaseContext) AllStmt() []IStmtContext

func (s *SwitchCaseContext) CASE() antlr.TerminalNode

func (s *SwitchCaseContext) COLON() antlr.TerminalNode

func (s *SwitchCaseContext) CasePatternList() ICasePatternListContext

func (s *SwitchCaseContext) GetParser() antlr.Parser

func (s *SwitchCaseContext) GetRuleContext() antlr.RuleContext

func (*SwitchCaseContext) IsSwitchCaseContext()

func (s *SwitchCaseContext) Stmt(i int) IStmtContext

func (s *SwitchCaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SwitchStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchStmtContext() *SwitchStmtContext

func NewSwitchStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchStmtContext

func (s *SwitchStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchStmtContext) AllSwitchCase() []ISwitchCaseContext

func (s *SwitchStmtContext) DefaultCase() IDefaultCaseContext

func (s *SwitchStmtContext) Expr() IExprContext

func (s *SwitchStmtContext) GetParser() antlr.Parser

func (s *SwitchStmtContext) GetRuleContext() antlr.RuleContext

func (*SwitchStmtContext) IsSwitchStmtContext()

func (s *SwitchStmtContext) LBRACE() antlr.TerminalNode

func (s *SwitchStmtContext) RBRACE() antlr.TerminalNode

func (s *SwitchStmtContext) SWITCH() antlr.TerminalNode

func (s *SwitchStmtContext) SwitchCase(i int) ISwitchCaseContext

func (s *SwitchStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TopLevelDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTopLevelDeclContext() *TopLevelDeclContext

func NewTopLevelDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopLevelDeclContext

func (s *TopLevelDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TopLevelDeclContext) ClassDecl() IClassDeclContext

func (s *TopLevelDeclContext) EnumDecl() IEnumDeclContext

func (s *TopLevelDeclContext) FuncDecl() IFuncDeclContext

func (s *TopLevelDeclContext) GetParser() antlr.Parser

func (s *TopLevelDeclContext) GetRuleContext() antlr.RuleContext

func (*TopLevelDeclContext) IsTopLevelDeclContext()

func (s *TopLevelDeclContext) StructDecl() IStructDeclContext

func (s *TopLevelDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TopLevelDeclContext) TypeAliasDecl() ITypeAliasDeclContext

func (s *TopLevelDeclContext) VarDeclStmt() IVarDeclStmtContext

type TupleBindContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleBindContext() *TupleBindContext

func NewTupleBindContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleBindContext

func (s *TupleBindContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleBindContext) AllCOMMA() []antlr.TerminalNode

func (s *TupleBindContext) AllID() []antlr.TerminalNode

func (s *TupleBindContext) COMMA(i int) antlr.TerminalNode

func (s *TupleBindContext) GetParser() antlr.Parser

func (s *TupleBindContext) GetRuleContext() antlr.RuleContext

func (s *TupleBindContext) ID(i int) antlr.TerminalNode

func (*TupleBindContext) IsTupleBindContext()

func (s *TupleBindContext) LPAREN() antlr.TerminalNode

func (s *TupleBindContext) RPAREN() antlr.TerminalNode

func (s *TupleBindContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TupleElementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleElementContext() *TupleElementContext

func NewTupleElementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleElementContext

func (s *TupleElementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleElementContext) COLON() antlr.TerminalNode

func (s *TupleElementContext) Expr() IExprContext

func (s *TupleElementContext) GetParser() antlr.Parser

func (s *TupleElementContext) GetRuleContext() antlr.RuleContext

func (s *TupleElementContext) ID() antlr.TerminalNode

func (*TupleElementContext) IsTupleElementContext()

func (s *TupleElementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TupleExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleExprContext() *TupleExprContext

func NewTupleExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleExprContext

func (s *TupleExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleExprContext) AllCOMMA() []antlr.TerminalNode

func (s *TupleExprContext) AllTupleElement() []ITupleElementContext

func (s *TupleExprContext) COMMA(i int) antlr.TerminalNode

func (s *TupleExprContext) GetParser() antlr.Parser

func (s *TupleExprContext) GetRuleContext() antlr.RuleContext

func (*TupleExprContext) IsTupleExprContext()

func (s *TupleExprContext) LPAREN() antlr.TerminalNode

func (s *TupleExprContext) RPAREN() antlr.TerminalNode

func (s *TupleExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleExprContext) TupleElement(i int) ITupleElementContext

type TupleTypeElemContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleTypeElemContext() *TupleTypeElemContext

func NewTupleTypeElemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleTypeElemContext

func (s *TupleTypeElemContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleTypeElemContext) COLON() antlr.TerminalNode

func (s *TupleTypeElemContext) GetParser() antlr.Parser

func (s *TupleTypeElemContext) GetRuleContext() antlr.RuleContext

func (s *TupleTypeElemContext) ID() antlr.TerminalNode

func (*TupleTypeElemContext) IsTupleTypeElemContext()

func (s *TupleTypeElemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleTypeElemContext) Type_() ITypeContext

type TypeAliasDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeAliasDeclContext() *TypeAliasDeclContext

func NewTypeAliasDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeAliasDeclContext

func (s *TypeAliasDeclContext) ASSIGN() antlr.TerminalNode

func (s *TypeAliasDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeAliasDeclContext) GetParser() antlr.Parser

func (s *TypeAliasDeclContext) GetRuleContext() antlr.RuleContext

func (s *TypeAliasDeclContext) ID() antlr.TerminalNode

func (*TypeAliasDeclContext) IsTypeAliasDeclContext()

func (s *TypeAliasDeclContext) TYPE() antlr.TerminalNode

func (s *TypeAliasDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeAliasDeclContext) Type_() ITypeContext

type TypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeContext() *TypeContext

func NewTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeContext

func (s *TypeContext) ANY() antlr.TerminalNode

func (s *TypeContext) ARROW() antlr.TerminalNode

func (s *TypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeContext) AllCOMMA() []antlr.TerminalNode

func (s *TypeContext) AllTupleTypeElem() []ITupleTypeElemContext

func (s *TypeContext) AllType_() []ITypeContext

func (s *TypeContext) CHANNEL() antlr.TerminalNode

func (s *TypeContext) COLON() antlr.TerminalNode

func (s *TypeContext) COMMA(i int) antlr.TerminalNode

func (s *TypeContext) FUNC() antlr.TerminalNode

func (s *TypeContext) FuncTypeParamList() IFuncTypeParamListContext

func (s *TypeContext) GT() antlr.TerminalNode

func (s *TypeContext) GetParser() antlr.Parser

func (s *TypeContext) GetRuleContext() antlr.RuleContext

func (s *TypeContext) ID() antlr.TerminalNode

func (*TypeContext) IsTypeContext()

func (s *TypeContext) LBRACKET() antlr.TerminalNode

func (s *TypeContext) LPAREN() antlr.TerminalNode

func (s *TypeContext) LT() antlr.TerminalNode

func (s *TypeContext) MUT() antlr.TerminalNode

func (s *TypeContext) OPAQUE() antlr.TerminalNode

func (s *TypeContext) PrimitiveType() IPrimitiveTypeContext

func (s *TypeContext) QUESTION() antlr.TerminalNode

func (s *TypeContext) RBRACKET() antlr.TerminalNode

func (s *TypeContext) RESULT() antlr.TerminalNode

func (s *TypeContext) RPAREN() antlr.TerminalNode

func (s *TypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeContext) TupleTypeElem(i int) ITupleTypeElemContext

func (s *TypeContext) Type_(i int) ITypeContext

func (s *TypeContext) VOID() antlr.TerminalNode

type TypeParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeParamContext() *TypeParamContext

func NewTypeParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeParamContext

func (s *TypeParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeParamContext) GetParser() antlr.Parser

func (s *TypeParamContext) GetRuleContext() antlr.RuleContext

func (s *TypeParamContext) ID() antlr.TerminalNode

func (*TypeParamContext) IsTypeParamContext()

func (s *TypeParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type VarDeclStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyVarDeclStmtContext() *VarDeclStmtContext

func NewVarDeclStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarDeclStmtContext

func (s *VarDeclStmtContext) ASSIGN() antlr.TerminalNode

func (s *VarDeclStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *VarDeclStmtContext) BindingKw() IBindingKwContext

func (s *VarDeclStmtContext) COLON() antlr.TerminalNode

func (s *VarDeclStmtContext) Expr() IExprContext

func (s *VarDeclStmtContext) GetParser() antlr.Parser

func (s *VarDeclStmtContext) GetRuleContext() antlr.RuleContext

func (s *VarDeclStmtContext) ID() antlr.TerminalNode

func (*VarDeclStmtContext) IsVarDeclStmtContext()

func (s *VarDeclStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *VarDeclStmtContext) TupleBind() ITupleBindContext

func (s *VarDeclStmtContext) Type_() ITypeContext

type VertexLexer struct {
	*antlr.BaseLexer

	// Has unexported fields.
}

func NewVertexLexer(input antlr.CharStream) *VertexLexer
    NewVertexLexer produces a new lexer instance for the optional input
    antlr.CharStream.

type VertexParser struct {
	*antlr.BaseParser
}

func NewVertexParser(input antlr.TokenStream) *VertexParser
    NewVertexParser produces a new parser instance for the optional input
    antlr.TokenStream.

func (p *VertexParser) AnonFuncExpr() (localctx IAnonFuncExprContext)

func (p *VertexParser) Arg() (localctx IArgContext)

func (p *VertexParser) ArgList() (localctx IArgListContext)

func (p *VertexParser) ArrayConstructExpr() (localctx IArrayConstructExprContext)

func (p *VertexParser) ArrayLiteralExpr() (localctx IArrayLiteralExprContext)

func (p *VertexParser) AssignStmt() (localctx IAssignStmtContext)

func (p *VertexParser) BindingKw() (localctx IBindingKwContext)

func (p *VertexParser) Block() (localctx IBlockContext)

func (p *VertexParser) BreakStmt() (localctx IBreakStmtContext)

func (p *VertexParser) BuildDecl() (localctx IBuildDeclContext)

func (p *VertexParser) BuildTag() (localctx IBuildTagContext)

func (p *VertexParser) CasePattern() (localctx ICasePatternContext)

func (p *VertexParser) CasePatternList() (localctx ICasePatternListContext)

func (p *VertexParser) ClassDecl() (localctx IClassDeclContext)

func (p *VertexParser) ClassField() (localctx IClassFieldContext)

func (p *VertexParser) ClassMember() (localctx IClassMemberContext)

func (p *VertexParser) CompoundAssignStmt() (localctx ICompoundAssignStmtContext)

func (p *VertexParser) CompoundOp() (localctx ICompoundOpContext)

func (p *VertexParser) ContinueStmt() (localctx IContinueStmtContext)

func (p *VertexParser) DefaultCase() (localctx IDefaultCaseContext)

func (p *VertexParser) DeferStmt() (localctx IDeferStmtContext)

func (p *VertexParser) DictEntry() (localctx IDictEntryContext)

func (p *VertexParser) DictLiteralExpr() (localctx IDictLiteralExprContext)

func (p *VertexParser) ElseClause() (localctx IElseClauseContext)

func (p *VertexParser) ElseIfClause() (localctx IElseIfClauseContext)

func (p *VertexParser) EnumCase() (localctx IEnumCaseContext)

func (p *VertexParser) EnumCaseDecl() (localctx IEnumCaseDeclContext)

func (p *VertexParser) EnumDecl() (localctx IEnumDeclContext)

func (p *VertexParser) EnumRawType() (localctx IEnumRawTypeContext)

func (p *VertexParser) Expr() (localctx IExprContext)

func (p *VertexParser) ExprList() (localctx IExprListContext)

func (p *VertexParser) ExprStmt() (localctx IExprStmtContext)

func (p *VertexParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) FallthroughStmt() (localctx IFallthroughStmtContext)

func (p *VertexParser) File() (localctx IFileContext)

func (p *VertexParser) ForInStmt() (localctx IForInStmtContext)

func (p *VertexParser) FuncDecl() (localctx IFuncDeclContext)

func (p *VertexParser) FuncQualifier() (localctx IFuncQualifierContext)

func (p *VertexParser) FuncTypeParam() (localctx IFuncTypeParamContext)

func (p *VertexParser) FuncTypeParamList() (localctx IFuncTypeParamListContext)

func (p *VertexParser) GenericParams() (localctx IGenericParamsContext)

func (p *VertexParser) IfCondition() (localctx IIfConditionContext)

func (p *VertexParser) IfStmt() (localctx IIfStmtContext)

func (p *VertexParser) ImportDecl() (localctx IImportDeclContext)

func (p *VertexParser) Literal() (localctx ILiteralContext)

func (p *VertexParser) Lvalue() (localctx ILvalueContext)

func (p *VertexParser) Lvalue_Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) NativeFuncDecl() (localctx INativeFuncDeclContext)

func (p *VertexParser) NativeParam() (localctx INativeParamContext)

func (p *VertexParser) NativeParamList() (localctx INativeParamListContext)

func (p *VertexParser) PackageDecl() (localctx IPackageDeclContext)

func (p *VertexParser) Param() (localctx IParamContext)

func (p *VertexParser) ParamList() (localctx IParamListContext)

func (p *VertexParser) PostfixName() (localctx IPostfixNameContext)

func (p *VertexParser) Primary() (localctx IPrimaryContext)

func (p *VertexParser) PrimitiveType() (localctx IPrimitiveTypeContext)

func (p *VertexParser) ReturnStmt() (localctx IReturnStmtContext)

func (p *VertexParser) ReturnType() (localctx IReturnTypeContext)

func (p *VertexParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool

func (p *VertexParser) Stmt() (localctx IStmtContext)

func (p *VertexParser) StructDecl() (localctx IStructDeclContext)

func (p *VertexParser) StructField() (localctx IStructFieldContext)

func (p *VertexParser) StructFieldInit() (localctx IStructFieldInitContext)

func (p *VertexParser) StructLiteralExpr() (localctx IStructLiteralExprContext)

func (p *VertexParser) SwitchCase() (localctx ISwitchCaseContext)

func (p *VertexParser) SwitchStmt() (localctx ISwitchStmtContext)

func (p *VertexParser) TopLevelDecl() (localctx ITopLevelDeclContext)

func (p *VertexParser) TupleBind() (localctx ITupleBindContext)

func (p *VertexParser) TupleElement() (localctx ITupleElementContext)

func (p *VertexParser) TupleExpr() (localctx ITupleExprContext)

func (p *VertexParser) TupleTypeElem() (localctx ITupleTypeElemContext)

func (p *VertexParser) TypeAliasDecl() (localctx ITypeAliasDeclContext)

func (p *VertexParser) TypeParam() (localctx ITypeParamContext)

func (p *VertexParser) Type_() (localctx ITypeContext)

func (p *VertexParser) Type__Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) VarDeclStmt() (localctx IVarDeclStmtContext)

func (p *VertexParser) WhileStmt() (localctx IWhileStmtContext)

type VertexParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by VertexParser#file.
	VisitFile(ctx *FileContext) interface{}

	// Visit a parse tree produced by VertexParser#packageDecl.
	VisitPackageDecl(ctx *PackageDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#buildDecl.
	VisitBuildDecl(ctx *BuildDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#buildTag.
	VisitBuildTag(ctx *BuildTagContext) interface{}

	// Visit a parse tree produced by VertexParser#importDecl.
	VisitImportDecl(ctx *ImportDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#topLevelDecl.
	VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#funcDecl.
	VisitFuncDecl(ctx *FuncDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#genericParams.
	VisitGenericParams(ctx *GenericParamsContext) interface{}

	// Visit a parse tree produced by VertexParser#typeParam.
	VisitTypeParam(ctx *TypeParamContext) interface{}

	// Visit a parse tree produced by VertexParser#paramList.
	VisitParamList(ctx *ParamListContext) interface{}

	// Visit a parse tree produced by VertexParser#param.
	VisitParam(ctx *ParamContext) interface{}

	// Visit a parse tree produced by VertexParser#funcQualifier.
	VisitFuncQualifier(ctx *FuncQualifierContext) interface{}

	// Visit a parse tree produced by VertexParser#returnType.
	VisitReturnType(ctx *ReturnTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#structDecl.
	VisitStructDecl(ctx *StructDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#structField.
	VisitStructField(ctx *StructFieldContext) interface{}

	// Visit a parse tree produced by VertexParser#structLiteralExpr.
	VisitStructLiteralExpr(ctx *StructLiteralExprContext) interface{}

	// Visit a parse tree produced by VertexParser#structFieldInit.
	VisitStructFieldInit(ctx *StructFieldInitContext) interface{}

	// Visit a parse tree produced by VertexParser#classDecl.
	VisitClassDecl(ctx *ClassDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#classMember.
	VisitClassMember(ctx *ClassMemberContext) interface{}

	// Visit a parse tree produced by VertexParser#classField.
	VisitClassField(ctx *ClassFieldContext) interface{}

	// Visit a parse tree produced by VertexParser#nativeFuncDecl.
	VisitNativeFuncDecl(ctx *NativeFuncDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#nativeParamList.
	VisitNativeParamList(ctx *NativeParamListContext) interface{}

	// Visit a parse tree produced by VertexParser#nativeParam.
	VisitNativeParam(ctx *NativeParamContext) interface{}

	// Visit a parse tree produced by VertexParser#enumDecl.
	VisitEnumDecl(ctx *EnumDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#enumRawType.
	VisitEnumRawType(ctx *EnumRawTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#enumCaseDecl.
	VisitEnumCaseDecl(ctx *EnumCaseDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#enumCase.
	VisitEnumCase(ctx *EnumCaseContext) interface{}

	// Visit a parse tree produced by VertexParser#typeAliasDecl.
	VisitTypeAliasDecl(ctx *TypeAliasDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by VertexParser#stmt.
	VisitStmt(ctx *StmtContext) interface{}

	// Visit a parse tree produced by VertexParser#varDeclStmt.
	VisitVarDeclStmt(ctx *VarDeclStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#bindingKw.
	VisitBindingKw(ctx *BindingKwContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleBind.
	VisitTupleBind(ctx *TupleBindContext) interface{}

	// Visit a parse tree produced by VertexParser#assignStmt.
	VisitAssignStmt(ctx *AssignStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#compoundAssignStmt.
	VisitCompoundAssignStmt(ctx *CompoundAssignStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#compoundOp.
	VisitCompoundOp(ctx *CompoundOpContext) interface{}

	// Visit a parse tree produced by VertexParser#lvalue.
	VisitLvalue(ctx *LvalueContext) interface{}

	// Visit a parse tree produced by VertexParser#ifStmt.
	VisitIfStmt(ctx *IfStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#elseIfClause.
	VisitElseIfClause(ctx *ElseIfClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#elseClause.
	VisitElseClause(ctx *ElseClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#ifCondition.
	VisitIfCondition(ctx *IfConditionContext) interface{}

	// Visit a parse tree produced by VertexParser#switchStmt.
	VisitSwitchStmt(ctx *SwitchStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#switchCase.
	VisitSwitchCase(ctx *SwitchCaseContext) interface{}

	// Visit a parse tree produced by VertexParser#defaultCase.
	VisitDefaultCase(ctx *DefaultCaseContext) interface{}

	// Visit a parse tree produced by VertexParser#casePatternList.
	VisitCasePatternList(ctx *CasePatternListContext) interface{}

	// Visit a parse tree produced by VertexParser#casePattern.
	VisitCasePattern(ctx *CasePatternContext) interface{}

	// Visit a parse tree produced by VertexParser#forInStmt.
	VisitForInStmt(ctx *ForInStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#whileStmt.
	VisitWhileStmt(ctx *WhileStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#breakStmt.
	VisitBreakStmt(ctx *BreakStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#continueStmt.
	VisitContinueStmt(ctx *ContinueStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#fallthroughStmt.
	VisitFallthroughStmt(ctx *FallthroughStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#returnStmt.
	VisitReturnStmt(ctx *ReturnStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#deferStmt.
	VisitDeferStmt(ctx *DeferStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#exprStmt.
	VisitExprStmt(ctx *ExprStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by VertexParser#postfixName.
	VisitPostfixName(ctx *PostfixNameContext) interface{}

	// Visit a parse tree produced by VertexParser#primary.
	VisitPrimary(ctx *PrimaryContext) interface{}

	// Visit a parse tree produced by VertexParser#argList.
	VisitArgList(ctx *ArgListContext) interface{}

	// Visit a parse tree produced by VertexParser#arg.
	VisitArg(ctx *ArgContext) interface{}

	// Visit a parse tree produced by VertexParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}

	// Visit a parse tree produced by VertexParser#arrayLiteralExpr.
	VisitArrayLiteralExpr(ctx *ArrayLiteralExprContext) interface{}

	// Visit a parse tree produced by VertexParser#arrayConstructExpr.
	VisitArrayConstructExpr(ctx *ArrayConstructExprContext) interface{}

	// Visit a parse tree produced by VertexParser#exprList.
	VisitExprList(ctx *ExprListContext) interface{}

	// Visit a parse tree produced by VertexParser#dictLiteralExpr.
	VisitDictLiteralExpr(ctx *DictLiteralExprContext) interface{}

	// Visit a parse tree produced by VertexParser#dictEntry.
	VisitDictEntry(ctx *DictEntryContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleExpr.
	VisitTupleExpr(ctx *TupleExprContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleElement.
	VisitTupleElement(ctx *TupleElementContext) interface{}

	// Visit a parse tree produced by VertexParser#anonFuncExpr.
	VisitAnonFuncExpr(ctx *AnonFuncExprContext) interface{}

	// Visit a parse tree produced by VertexParser#type.
	VisitType(ctx *TypeContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleTypeElem.
	VisitTupleTypeElem(ctx *TupleTypeElemContext) interface{}

	// Visit a parse tree produced by VertexParser#funcTypeParamList.
	VisitFuncTypeParamList(ctx *FuncTypeParamListContext) interface{}

	// Visit a parse tree produced by VertexParser#funcTypeParam.
	VisitFuncTypeParam(ctx *FuncTypeParamContext) interface{}

	// Visit a parse tree produced by VertexParser#primitiveType.
	VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{}
}
    A complete Visitor for a parse tree produced by VertexParser.

type WhileStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyWhileStmtContext() *WhileStmtContext

func NewWhileStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhileStmtContext

func (s *WhileStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *WhileStmtContext) Block() IBlockContext

func (s *WhileStmtContext) Expr() IExprContext

func (s *WhileStmtContext) GetParser() antlr.Parser

func (s *WhileStmtContext) GetRuleContext() antlr.RuleContext

func (*WhileStmtContext) IsWhileStmtContext()

func (s *WhileStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *WhileStmtContext) WHILE() antlr.TerminalNode

