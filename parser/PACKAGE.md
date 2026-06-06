

CONSTANTS

const (
	VertexLexerLET                  = 1
	VertexLexerVAR                  = 2
	VertexLexerFUNC                 = 3
	VertexLexerIF                   = 4
	VertexLexerELSE                 = 5
	VertexLexerFOR                  = 6
	VertexLexerIN                   = 7
	VertexLexerWHILE                = 8
	VertexLexerSWITCH               = 9
	VertexLexerCASE                 = 10
	VertexLexerDEFAULT              = 11
	VertexLexerRETURN               = 12
	VertexLexerBREAK                = 13
	VertexLexerCONTINUE             = 14
	VertexLexerFALLTHROUGH          = 15
	VertexLexerDEFER                = 16
	VertexLexerSTRUCT               = 17
	VertexLexerCLASS                = 18
	VertexLexerENUM                 = 19
	VertexLexerTYPE                 = 20
	VertexLexerIMPORT               = 21
	VertexLexerPACKAGE              = 22
	VertexLexerBUILD                = 23
	VertexLexerASM                  = 24
	VertexLexerWEAK                 = 25
	VertexLexerCONST_KW             = 26
	VertexLexerREINTERPRET          = 27
	VertexLexerASYNC                = 28
	VertexLexerTHREAD               = 29
	VertexLexerPROCESS              = 30
	VertexLexerGPU                  = 31
	VertexLexerTEST                 = 32
	VertexLexerCHAN                 = 33
	VertexLexerMAP                  = 34
	VertexLexerINOUT                = 35
	VertexLexerOUT_KW               = 36
	VertexLexerCLOBBER              = 37
	VertexLexerRESULT               = 38
	VertexLexerOK                   = 39
	VertexLexerERR_KW               = 40
	VertexLexerEXPECTED             = 41
	VertexLexerTRUE                 = 42
	VertexLexerFALSE                = 43
	VertexLexerNIL                  = 44
	VertexLexerELLIPSIS             = 45
	VertexLexerHALF_OPEN            = 46
	VertexLexerIDENTITY_EQ          = 47
	VertexLexerIDENTITY_NEQ         = 48
	VertexLexerOVERFLOW_ADD         = 49
	VertexLexerOVERFLOW_SUB         = 50
	VertexLexerOVERFLOW_MUL         = 51
	VertexLexerNIL_COALESCE         = 52
	VertexLexerARROW                = 53
	VertexLexerLSHIFT               = 54
	VertexLexerRSHIFT               = 55
	VertexLexerLEQ                  = 56
	VertexLexerGEQ                  = 57
	VertexLexerEQ                   = 58
	VertexLexerNEQ                  = 59
	VertexLexerLOGICAL_AND          = 60
	VertexLexerLOGICAL_OR           = 61
	VertexLexerPLUS_ASSIGN          = 62
	VertexLexerMINUS_ASSIGN         = 63
	VertexLexerSTAR_ASSIGN          = 64
	VertexLexerDIV_ASSIGN           = 65
	VertexLexerMOD_ASSIGN           = 66
	VertexLexerPLUS                 = 67
	VertexLexerMINUS                = 68
	VertexLexerSTAR                 = 69
	VertexLexerSLASH                = 70
	VertexLexerPERCENT              = 71
	VertexLexerTILDE                = 72
	VertexLexerAMP                  = 73
	VertexLexerPIPE                 = 74
	VertexLexerCARET                = 75
	VertexLexerBANG                 = 76
	VertexLexerLT                   = 77
	VertexLexerGT                   = 78
	VertexLexerASSIGN               = 79
	VertexLexerQUESTION             = 80
	VertexLexerCOLON                = 81
	VertexLexerLPAREN               = 82
	VertexLexerRPAREN               = 83
	VertexLexerLBRACE               = 84
	VertexLexerRBRACE               = 85
	VertexLexerLBRACKET             = 86
	VertexLexerRBRACKET             = 87
	VertexLexerDOT                  = 88
	VertexLexerCOMMA                = 89
	VertexLexerHEX_FLOAT_LIT        = 90
	VertexLexerHEX_INT_LIT          = 91
	VertexLexerOCT_INT_LIT          = 92
	VertexLexerBIN_INT_LIT          = 93
	VertexLexerDEC_FLOAT_LIT        = 94
	VertexLexerDEC_INT_LIT          = 95
	VertexLexerCHAR_LIT             = 96
	VertexLexerSTRING_LIT           = 97
	VertexLexerMULTILINE_STRING_LIT = 98
	VertexLexerIDENTIFIER           = 99
	VertexLexerWS                   = 100
	VertexLexerLINE_COMMENT         = 101
)
    VertexLexer tokens.

const (
	VertexParserEOF                  = antlr.TokenEOF
	VertexParserLET                  = 1
	VertexParserVAR                  = 2
	VertexParserFUNC                 = 3
	VertexParserIF                   = 4
	VertexParserELSE                 = 5
	VertexParserFOR                  = 6
	VertexParserIN                   = 7
	VertexParserWHILE                = 8
	VertexParserSWITCH               = 9
	VertexParserCASE                 = 10
	VertexParserDEFAULT              = 11
	VertexParserRETURN               = 12
	VertexParserBREAK                = 13
	VertexParserCONTINUE             = 14
	VertexParserFALLTHROUGH          = 15
	VertexParserDEFER                = 16
	VertexParserSTRUCT               = 17
	VertexParserCLASS                = 18
	VertexParserENUM                 = 19
	VertexParserTYPE                 = 20
	VertexParserIMPORT               = 21
	VertexParserPACKAGE              = 22
	VertexParserBUILD                = 23
	VertexParserASM                  = 24
	VertexParserWEAK                 = 25
	VertexParserCONST_KW             = 26
	VertexParserREINTERPRET          = 27
	VertexParserASYNC                = 28
	VertexParserTHREAD               = 29
	VertexParserPROCESS              = 30
	VertexParserGPU                  = 31
	VertexParserTEST                 = 32
	VertexParserCHAN                 = 33
	VertexParserMAP                  = 34
	VertexParserINOUT                = 35
	VertexParserOUT_KW               = 36
	VertexParserCLOBBER              = 37
	VertexParserRESULT               = 38
	VertexParserOK                   = 39
	VertexParserERR_KW               = 40
	VertexParserEXPECTED             = 41
	VertexParserTRUE                 = 42
	VertexParserFALSE                = 43
	VertexParserNIL                  = 44
	VertexParserELLIPSIS             = 45
	VertexParserHALF_OPEN            = 46
	VertexParserIDENTITY_EQ          = 47
	VertexParserIDENTITY_NEQ         = 48
	VertexParserOVERFLOW_ADD         = 49
	VertexParserOVERFLOW_SUB         = 50
	VertexParserOVERFLOW_MUL         = 51
	VertexParserNIL_COALESCE         = 52
	VertexParserARROW                = 53
	VertexParserLSHIFT               = 54
	VertexParserRSHIFT               = 55
	VertexParserLEQ                  = 56
	VertexParserGEQ                  = 57
	VertexParserEQ                   = 58
	VertexParserNEQ                  = 59
	VertexParserLOGICAL_AND          = 60
	VertexParserLOGICAL_OR           = 61
	VertexParserPLUS_ASSIGN          = 62
	VertexParserMINUS_ASSIGN         = 63
	VertexParserSTAR_ASSIGN          = 64
	VertexParserDIV_ASSIGN           = 65
	VertexParserMOD_ASSIGN           = 66
	VertexParserPLUS                 = 67
	VertexParserMINUS                = 68
	VertexParserSTAR                 = 69
	VertexParserSLASH                = 70
	VertexParserPERCENT              = 71
	VertexParserTILDE                = 72
	VertexParserAMP                  = 73
	VertexParserPIPE                 = 74
	VertexParserCARET                = 75
	VertexParserBANG                 = 76
	VertexParserLT                   = 77
	VertexParserGT                   = 78
	VertexParserASSIGN               = 79
	VertexParserQUESTION             = 80
	VertexParserCOLON                = 81
	VertexParserLPAREN               = 82
	VertexParserRPAREN               = 83
	VertexParserLBRACE               = 84
	VertexParserRBRACE               = 85
	VertexParserLBRACKET             = 86
	VertexParserRBRACKET             = 87
	VertexParserDOT                  = 88
	VertexParserCOMMA                = 89
	VertexParserHEX_FLOAT_LIT        = 90
	VertexParserHEX_INT_LIT          = 91
	VertexParserOCT_INT_LIT          = 92
	VertexParserBIN_INT_LIT          = 93
	VertexParserDEC_FLOAT_LIT        = 94
	VertexParserDEC_INT_LIT          = 95
	VertexParserCHAR_LIT             = 96
	VertexParserSTRING_LIT           = 97
	VertexParserMULTILINE_STRING_LIT = 98
	VertexParserIDENTIFIER           = 99
	VertexParserWS                   = 100
	VertexParserLINE_COMMENT         = 101
)
    VertexParser tokens.

const (
	VertexParserRULE_file                = 0
	VertexParserRULE_packageDecl         = 1
	VertexParserRULE_buildDecl           = 2
	VertexParserRULE_importDecl          = 3
	VertexParserRULE_topLevelDecl        = 4
	VertexParserRULE_typeAliasDecl       = 5
	VertexParserRULE_funcDecl            = 6
	VertexParserRULE_receiver            = 7
	VertexParserRULE_genericParams       = 8
	VertexParserRULE_funcQualifier       = 9
	VertexParserRULE_paramList           = 10
	VertexParserRULE_param               = 11
	VertexParserRULE_variadicParam       = 12
	VertexParserRULE_structDecl          = 13
	VertexParserRULE_structFieldDecl     = 14
	VertexParserRULE_classDecl           = 15
	VertexParserRULE_classMember         = 16
	VertexParserRULE_qualifiedIdent      = 17
	VertexParserRULE_enumDecl            = 18
	VertexParserRULE_enumCaseDecl        = 19
	VertexParserRULE_enumCaseItem        = 20
	VertexParserRULE_enumRawValue        = 21
	VertexParserRULE_block               = 22
	VertexParserRULE_stmt                = 23
	VertexParserRULE_varDecl             = 24
	VertexParserRULE_bindingPattern      = 25
	VertexParserRULE_ifStmt              = 26
	VertexParserRULE_ifCondition         = 27
	VertexParserRULE_whileStmt           = 28
	VertexParserRULE_forInStmt           = 29
	VertexParserRULE_switchStmt          = 30
	VertexParserRULE_switchCase          = 31
	VertexParserRULE_switchPattern       = 32
	VertexParserRULE_returnStmt          = 33
	VertexParserRULE_deferStmt           = 34
	VertexParserRULE_exprOrAssignStmt    = 35
	VertexParserRULE_assignOp            = 36
	VertexParserRULE_expr                = 37
	VertexParserRULE_argList             = 38
	VertexParserRULE_arg                 = 39
	VertexParserRULE_structLiteralFields = 40
	VertexParserRULE_structLiteralField  = 41
	VertexParserRULE_mapLiteralFields    = 42
	VertexParserRULE_mapLiteralField     = 43
	VertexParserRULE_anonFuncExpr        = 44
	VertexParserRULE_asmExpr             = 45
	VertexParserRULE_asmBody             = 46
	VertexParserRULE_asmInstr            = 47
	VertexParserRULE_asmConstraint       = 48
	VertexParserRULE_asmClobberDecl      = 49
	VertexParserRULE_typeExpr            = 50
	VertexParserRULE_funcTypeParams      = 51
	VertexParserRULE_baseType            = 52
	VertexParserRULE_tupleTypeElems      = 53
	VertexParserRULE_tupleTypeElem       = 54
	VertexParserRULE_literal             = 55
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
func InitEmptyAsmBodyContext(p *AsmBodyContext)
func InitEmptyAsmClobberDeclContext(p *AsmClobberDeclContext)
func InitEmptyAsmConstraintContext(p *AsmConstraintContext)
func InitEmptyAsmExprContext(p *AsmExprContext)
func InitEmptyAsmInstrContext(p *AsmInstrContext)
func InitEmptyAssignOpContext(p *AssignOpContext)
func InitEmptyBaseTypeContext(p *BaseTypeContext)
func InitEmptyBindingPatternContext(p *BindingPatternContext)
func InitEmptyBlockContext(p *BlockContext)
func InitEmptyBuildDeclContext(p *BuildDeclContext)
func InitEmptyClassDeclContext(p *ClassDeclContext)
func InitEmptyClassMemberContext(p *ClassMemberContext)
func InitEmptyDeferStmtContext(p *DeferStmtContext)
func InitEmptyEnumCaseDeclContext(p *EnumCaseDeclContext)
func InitEmptyEnumCaseItemContext(p *EnumCaseItemContext)
func InitEmptyEnumDeclContext(p *EnumDeclContext)
func InitEmptyEnumRawValueContext(p *EnumRawValueContext)
func InitEmptyExprContext(p *ExprContext)
func InitEmptyExprOrAssignStmtContext(p *ExprOrAssignStmtContext)
func InitEmptyFileContext(p *FileContext)
func InitEmptyForInStmtContext(p *ForInStmtContext)
func InitEmptyFuncDeclContext(p *FuncDeclContext)
func InitEmptyFuncQualifierContext(p *FuncQualifierContext)
func InitEmptyFuncTypeParamsContext(p *FuncTypeParamsContext)
func InitEmptyGenericParamsContext(p *GenericParamsContext)
func InitEmptyIfConditionContext(p *IfConditionContext)
func InitEmptyIfStmtContext(p *IfStmtContext)
func InitEmptyImportDeclContext(p *ImportDeclContext)
func InitEmptyLiteralContext(p *LiteralContext)
func InitEmptyMapLiteralFieldContext(p *MapLiteralFieldContext)
func InitEmptyMapLiteralFieldsContext(p *MapLiteralFieldsContext)
func InitEmptyPackageDeclContext(p *PackageDeclContext)
func InitEmptyParamContext(p *ParamContext)
func InitEmptyParamListContext(p *ParamListContext)
func InitEmptyQualifiedIdentContext(p *QualifiedIdentContext)
func InitEmptyReceiverContext(p *ReceiverContext)
func InitEmptyReturnStmtContext(p *ReturnStmtContext)
func InitEmptyStmtContext(p *StmtContext)
func InitEmptyStructDeclContext(p *StructDeclContext)
func InitEmptyStructFieldDeclContext(p *StructFieldDeclContext)
func InitEmptyStructLiteralFieldContext(p *StructLiteralFieldContext)
func InitEmptyStructLiteralFieldsContext(p *StructLiteralFieldsContext)
func InitEmptySwitchCaseContext(p *SwitchCaseContext)
func InitEmptySwitchPatternContext(p *SwitchPatternContext)
func InitEmptySwitchStmtContext(p *SwitchStmtContext)
func InitEmptyTopLevelDeclContext(p *TopLevelDeclContext)
func InitEmptyTupleTypeElemContext(p *TupleTypeElemContext)
func InitEmptyTupleTypeElemsContext(p *TupleTypeElemsContext)
func InitEmptyTypeAliasDeclContext(p *TypeAliasDeclContext)
func InitEmptyTypeExprContext(p *TypeExprContext)
func InitEmptyVarDeclContext(p *VarDeclContext)
func InitEmptyVariadicParamContext(p *VariadicParamContext)
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

func (s *AnonFuncExprContext) ARROW() antlr.TerminalNode

func (s *AnonFuncExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AnonFuncExprContext) Block() IBlockContext

func (s *AnonFuncExprContext) FUNC() antlr.TerminalNode

func (s *AnonFuncExprContext) FuncQualifier() IFuncQualifierContext

func (s *AnonFuncExprContext) GetParser() antlr.Parser

func (s *AnonFuncExprContext) GetRuleContext() antlr.RuleContext

func (*AnonFuncExprContext) IsAnonFuncExprContext()

func (s *AnonFuncExprContext) LPAREN() antlr.TerminalNode

func (s *AnonFuncExprContext) ParamList() IParamListContext

func (s *AnonFuncExprContext) RPAREN() antlr.TerminalNode

func (s *AnonFuncExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *AnonFuncExprContext) TypeExpr() ITypeExprContext

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

func (s *ArgContext) IDENTIFIER() antlr.TerminalNode

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

type AsmBodyContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAsmBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AsmBodyContext

func NewEmptyAsmBodyContext() *AsmBodyContext

func (s *AsmBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AsmBodyContext) AllAsmConstraint() []IAsmConstraintContext

func (s *AsmBodyContext) AllAsmInstr() []IAsmInstrContext

func (s *AsmBodyContext) AllCOMMA() []antlr.TerminalNode

func (s *AsmBodyContext) AsmClobberDecl() IAsmClobberDeclContext

func (s *AsmBodyContext) AsmConstraint(i int) IAsmConstraintContext

func (s *AsmBodyContext) AsmInstr(i int) IAsmInstrContext

func (s *AsmBodyContext) COMMA(i int) antlr.TerminalNode

func (s *AsmBodyContext) GetParser() antlr.Parser

func (s *AsmBodyContext) GetRuleContext() antlr.RuleContext

func (*AsmBodyContext) IsAsmBodyContext()

func (s *AsmBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AsmClobberDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAsmClobberDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AsmClobberDeclContext

func NewEmptyAsmClobberDeclContext() *AsmClobberDeclContext

func (s *AsmClobberDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AsmClobberDeclContext) AllCOMMA() []antlr.TerminalNode

func (s *AsmClobberDeclContext) AllSTRING_LIT() []antlr.TerminalNode

func (s *AsmClobberDeclContext) CLOBBER() antlr.TerminalNode

func (s *AsmClobberDeclContext) COMMA(i int) antlr.TerminalNode

func (s *AsmClobberDeclContext) GetParser() antlr.Parser

func (s *AsmClobberDeclContext) GetRuleContext() antlr.RuleContext

func (*AsmClobberDeclContext) IsAsmClobberDeclContext()

func (s *AsmClobberDeclContext) LPAREN() antlr.TerminalNode

func (s *AsmClobberDeclContext) RPAREN() antlr.TerminalNode

func (s *AsmClobberDeclContext) STRING_LIT(i int) antlr.TerminalNode

func (s *AsmClobberDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AsmConstraintContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAsmConstraintContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AsmConstraintContext

func NewEmptyAsmConstraintContext() *AsmConstraintContext

func (s *AsmConstraintContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AsmConstraintContext) GetParser() antlr.Parser

func (s *AsmConstraintContext) GetRuleContext() antlr.RuleContext

func (s *AsmConstraintContext) IDENTIFIER() antlr.TerminalNode

func (s *AsmConstraintContext) IN() antlr.TerminalNode

func (s *AsmConstraintContext) INOUT() antlr.TerminalNode

func (*AsmConstraintContext) IsAsmConstraintContext()

func (s *AsmConstraintContext) LPAREN() antlr.TerminalNode

func (s *AsmConstraintContext) OUT_KW() antlr.TerminalNode

func (s *AsmConstraintContext) RPAREN() antlr.TerminalNode

func (s *AsmConstraintContext) STRING_LIT() antlr.TerminalNode

func (s *AsmConstraintContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AsmExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAsmExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AsmExprContext

func NewEmptyAsmExprContext() *AsmExprContext

func (s *AsmExprContext) ASM() antlr.TerminalNode

func (s *AsmExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AsmExprContext) AsmBody() IAsmBodyContext

func (s *AsmExprContext) GetParser() antlr.Parser

func (s *AsmExprContext) GetRuleContext() antlr.RuleContext

func (*AsmExprContext) IsAsmExprContext()

func (s *AsmExprContext) LPAREN() antlr.TerminalNode

func (s *AsmExprContext) RPAREN() antlr.TerminalNode

func (s *AsmExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AsmInstrContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAsmInstrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AsmInstrContext

func NewEmptyAsmInstrContext() *AsmInstrContext

func (s *AsmInstrContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AsmInstrContext) GetParser() antlr.Parser

func (s *AsmInstrContext) GetRuleContext() antlr.RuleContext

func (*AsmInstrContext) IsAsmInstrContext()

func (s *AsmInstrContext) STRING_LIT() antlr.TerminalNode

func (s *AsmInstrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AssignOpContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAssignOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignOpContext

func NewEmptyAssignOpContext() *AssignOpContext

func (s *AssignOpContext) ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AssignOpContext) DIV_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) GetParser() antlr.Parser

func (s *AssignOpContext) GetRuleContext() antlr.RuleContext

func (*AssignOpContext) IsAssignOpContext()

func (s *AssignOpContext) MINUS_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) MOD_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) PLUS_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) STAR_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BaseTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBaseTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BaseTypeContext

func NewEmptyBaseTypeContext() *BaseTypeContext

func (s *BaseTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BaseTypeContext) AllDOT() []antlr.TerminalNode

func (s *BaseTypeContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *BaseTypeContext) DOT(i int) antlr.TerminalNode

func (s *BaseTypeContext) GetParser() antlr.Parser

func (s *BaseTypeContext) GetRuleContext() antlr.RuleContext

func (s *BaseTypeContext) IDENTIFIER(i int) antlr.TerminalNode

func (*BaseTypeContext) IsBaseTypeContext()

func (s *BaseTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BaseVertexParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseVertexParserVisitor) VisitAnonFuncExpr(ctx *AnonFuncExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitArg(ctx *ArgContext) interface{}

func (v *BaseVertexParserVisitor) VisitArgList(ctx *ArgListContext) interface{}

func (v *BaseVertexParserVisitor) VisitAsmBody(ctx *AsmBodyContext) interface{}

func (v *BaseVertexParserVisitor) VisitAsmClobberDecl(ctx *AsmClobberDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitAsmConstraint(ctx *AsmConstraintContext) interface{}

func (v *BaseVertexParserVisitor) VisitAsmExpr(ctx *AsmExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitAsmInstr(ctx *AsmInstrContext) interface{}

func (v *BaseVertexParserVisitor) VisitAssignOp(ctx *AssignOpContext) interface{}

func (v *BaseVertexParserVisitor) VisitBaseType(ctx *BaseTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitBindingPattern(ctx *BindingPatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitBlock(ctx *BlockContext) interface{}

func (v *BaseVertexParserVisitor) VisitBuildDecl(ctx *BuildDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitClassDecl(ctx *ClassDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitClassMember(ctx *ClassMemberContext) interface{}

func (v *BaseVertexParserVisitor) VisitDeferStmt(ctx *DeferStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumCaseDecl(ctx *EnumCaseDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumCaseItem(ctx *EnumCaseItemContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumDecl(ctx *EnumDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumRawValue(ctx *EnumRawValueContext) interface{}

func (v *BaseVertexParserVisitor) VisitExpr(ctx *ExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitExprOrAssignStmt(ctx *ExprOrAssignStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitFile(ctx *FileContext) interface{}

func (v *BaseVertexParserVisitor) VisitForInStmt(ctx *ForInStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitFuncDecl(ctx *FuncDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitFuncQualifier(ctx *FuncQualifierContext) interface{}

func (v *BaseVertexParserVisitor) VisitFuncTypeParams(ctx *FuncTypeParamsContext) interface{}

func (v *BaseVertexParserVisitor) VisitGenericParams(ctx *GenericParamsContext) interface{}

func (v *BaseVertexParserVisitor) VisitIfCondition(ctx *IfConditionContext) interface{}

func (v *BaseVertexParserVisitor) VisitIfStmt(ctx *IfStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitImportDecl(ctx *ImportDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitLiteral(ctx *LiteralContext) interface{}

func (v *BaseVertexParserVisitor) VisitMapLiteralField(ctx *MapLiteralFieldContext) interface{}

func (v *BaseVertexParserVisitor) VisitMapLiteralFields(ctx *MapLiteralFieldsContext) interface{}

func (v *BaseVertexParserVisitor) VisitPackageDecl(ctx *PackageDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitParam(ctx *ParamContext) interface{}

func (v *BaseVertexParserVisitor) VisitParamList(ctx *ParamListContext) interface{}

func (v *BaseVertexParserVisitor) VisitQualifiedIdent(ctx *QualifiedIdentContext) interface{}

func (v *BaseVertexParserVisitor) VisitReceiver(ctx *ReceiverContext) interface{}

func (v *BaseVertexParserVisitor) VisitReturnStmt(ctx *ReturnStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitStmt(ctx *StmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructDecl(ctx *StructDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructFieldDecl(ctx *StructFieldDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructLiteralField(ctx *StructLiteralFieldContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructLiteralFields(ctx *StructLiteralFieldsContext) interface{}

func (v *BaseVertexParserVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{}

func (v *BaseVertexParserVisitor) VisitSwitchPattern(ctx *SwitchPatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitSwitchStmt(ctx *SwitchStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleTypeElem(ctx *TupleTypeElemContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleTypeElems(ctx *TupleTypeElemsContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeAliasDecl(ctx *TypeAliasDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeExpr(ctx *TypeExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitVarDecl(ctx *VarDeclContext) interface{}

func (v *BaseVertexParserVisitor) VisitVariadicParam(ctx *VariadicParamContext) interface{}

func (v *BaseVertexParserVisitor) VisitWhileStmt(ctx *WhileStmtContext) interface{}

type BindingPatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBindingPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BindingPatternContext

func NewEmptyBindingPatternContext() *BindingPatternContext

func (s *BindingPatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BindingPatternContext) AllCOMMA() []antlr.TerminalNode

func (s *BindingPatternContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *BindingPatternContext) COMMA(i int) antlr.TerminalNode

func (s *BindingPatternContext) GetParser() antlr.Parser

func (s *BindingPatternContext) GetRuleContext() antlr.RuleContext

func (s *BindingPatternContext) IDENTIFIER(i int) antlr.TerminalNode

func (*BindingPatternContext) IsBindingPatternContext()

func (s *BindingPatternContext) LPAREN() antlr.TerminalNode

func (s *BindingPatternContext) RPAREN() antlr.TerminalNode

func (s *BindingPatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

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

type BuildDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBuildDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BuildDeclContext

func NewEmptyBuildDeclContext() *BuildDeclContext

func (s *BuildDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BuildDeclContext) BUILD() antlr.TerminalNode

func (s *BuildDeclContext) GetParser() antlr.Parser

func (s *BuildDeclContext) GetRuleContext() antlr.RuleContext

func (s *BuildDeclContext) IDENTIFIER() antlr.TerminalNode

func (*BuildDeclContext) IsBuildDeclContext()

func (s *BuildDeclContext) TEST() antlr.TerminalNode

func (s *BuildDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClassDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClassDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassDeclContext

func NewEmptyClassDeclContext() *ClassDeclContext

func (s *ClassDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClassDeclContext) AllClassMember() []IClassMemberContext

func (s *ClassDeclContext) CLASS() antlr.TerminalNode

func (s *ClassDeclContext) COLON() antlr.TerminalNode

func (s *ClassDeclContext) ClassMember(i int) IClassMemberContext

func (s *ClassDeclContext) GetParser() antlr.Parser

func (s *ClassDeclContext) GetRuleContext() antlr.RuleContext

func (s *ClassDeclContext) IDENTIFIER() antlr.TerminalNode

func (*ClassDeclContext) IsClassDeclContext()

func (s *ClassDeclContext) LBRACE() antlr.TerminalNode

func (s *ClassDeclContext) QualifiedIdent() IQualifiedIdentContext

func (s *ClassDeclContext) RBRACE() antlr.TerminalNode

func (s *ClassDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClassMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClassMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassMemberContext

func NewEmptyClassMemberContext() *ClassMemberContext

func (s *ClassMemberContext) ARROW() antlr.TerminalNode

func (s *ClassMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClassMemberContext) COLON() antlr.TerminalNode

func (s *ClassMemberContext) FUNC() antlr.TerminalNode

func (s *ClassMemberContext) GetParser() antlr.Parser

func (s *ClassMemberContext) GetRuleContext() antlr.RuleContext

func (s *ClassMemberContext) IDENTIFIER() antlr.TerminalNode

func (*ClassMemberContext) IsClassMemberContext()

func (s *ClassMemberContext) LPAREN() antlr.TerminalNode

func (s *ClassMemberContext) ParamList() IParamListContext

func (s *ClassMemberContext) RPAREN() antlr.TerminalNode

func (s *ClassMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ClassMemberContext) TypeExpr() ITypeExprContext

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

type EnumCaseDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumCaseDeclContext() *EnumCaseDeclContext

func NewEnumCaseDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumCaseDeclContext

func (s *EnumCaseDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumCaseDeclContext) AllCOMMA() []antlr.TerminalNode

func (s *EnumCaseDeclContext) AllEnumCaseItem() []IEnumCaseItemContext

func (s *EnumCaseDeclContext) CASE() antlr.TerminalNode

func (s *EnumCaseDeclContext) COMMA(i int) antlr.TerminalNode

func (s *EnumCaseDeclContext) EnumCaseItem(i int) IEnumCaseItemContext

func (s *EnumCaseDeclContext) GetParser() antlr.Parser

func (s *EnumCaseDeclContext) GetRuleContext() antlr.RuleContext

func (*EnumCaseDeclContext) IsEnumCaseDeclContext()

func (s *EnumCaseDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type EnumCaseItemContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumCaseItemContext() *EnumCaseItemContext

func NewEnumCaseItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumCaseItemContext

func (s *EnumCaseItemContext) ASSIGN() antlr.TerminalNode

func (s *EnumCaseItemContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumCaseItemContext) EnumRawValue() IEnumRawValueContext

func (s *EnumCaseItemContext) GetParser() antlr.Parser

func (s *EnumCaseItemContext) GetRuleContext() antlr.RuleContext

func (s *EnumCaseItemContext) IDENTIFIER() antlr.TerminalNode

func (*EnumCaseItemContext) IsEnumCaseItemContext()

func (s *EnumCaseItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

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

func (s *EnumDeclContext) GetParser() antlr.Parser

func (s *EnumDeclContext) GetRuleContext() antlr.RuleContext

func (s *EnumDeclContext) IDENTIFIER() antlr.TerminalNode

func (*EnumDeclContext) IsEnumDeclContext()

func (s *EnumDeclContext) LBRACE() antlr.TerminalNode

func (s *EnumDeclContext) RBRACE() antlr.TerminalNode

func (s *EnumDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *EnumDeclContext) TypeExpr() ITypeExprContext

type EnumRawValueContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumRawValueContext() *EnumRawValueContext

func NewEnumRawValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumRawValueContext

func (s *EnumRawValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumRawValueContext) BIN_INT_LIT() antlr.TerminalNode

func (s *EnumRawValueContext) DEC_INT_LIT() antlr.TerminalNode

func (s *EnumRawValueContext) GetParser() antlr.Parser

func (s *EnumRawValueContext) GetRuleContext() antlr.RuleContext

func (s *EnumRawValueContext) HEX_INT_LIT() antlr.TerminalNode

func (*EnumRawValueContext) IsEnumRawValueContext()

func (s *EnumRawValueContext) MINUS() antlr.TerminalNode

func (s *EnumRawValueContext) OCT_INT_LIT() antlr.TerminalNode

func (s *EnumRawValueContext) STRING_LIT() antlr.TerminalNode

func (s *EnumRawValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExprContext() *ExprContext

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext

func (s *ExprContext) AMP() antlr.TerminalNode

func (s *ExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExprContext) AllCOMMA() []antlr.TerminalNode

func (s *ExprContext) AllExpr() []IExprContext

func (s *ExprContext) AllTypeExpr() []ITypeExprContext

func (s *ExprContext) AnonFuncExpr() IAnonFuncExprContext

func (s *ExprContext) ArgList() IArgListContext

func (s *ExprContext) AsmExpr() IAsmExprContext

func (s *ExprContext) BANG() antlr.TerminalNode

func (s *ExprContext) CARET() antlr.TerminalNode

func (s *ExprContext) COLON() antlr.TerminalNode

func (s *ExprContext) COMMA(i int) antlr.TerminalNode

func (s *ExprContext) DOT() antlr.TerminalNode

func (s *ExprContext) ELLIPSIS() antlr.TerminalNode

func (s *ExprContext) EQ() antlr.TerminalNode

func (s *ExprContext) ERR_KW() antlr.TerminalNode

func (s *ExprContext) Expr(i int) IExprContext

func (s *ExprContext) GEQ() antlr.TerminalNode

func (s *ExprContext) GT() antlr.TerminalNode

func (s *ExprContext) GetParser() antlr.Parser

func (s *ExprContext) GetRuleContext() antlr.RuleContext

func (s *ExprContext) HALF_OPEN() antlr.TerminalNode

func (s *ExprContext) IDENTIFIER() antlr.TerminalNode

func (s *ExprContext) IDENTITY_EQ() antlr.TerminalNode

func (s *ExprContext) IDENTITY_NEQ() antlr.TerminalNode

func (*ExprContext) IsExprContext()

func (s *ExprContext) LBRACE() antlr.TerminalNode

func (s *ExprContext) LBRACKET() antlr.TerminalNode

func (s *ExprContext) LEQ() antlr.TerminalNode

func (s *ExprContext) LOGICAL_AND() antlr.TerminalNode

func (s *ExprContext) LOGICAL_OR() antlr.TerminalNode

func (s *ExprContext) LPAREN() antlr.TerminalNode

func (s *ExprContext) LSHIFT() antlr.TerminalNode

func (s *ExprContext) LT() antlr.TerminalNode

func (s *ExprContext) Literal() ILiteralContext

func (s *ExprContext) MAP() antlr.TerminalNode

func (s *ExprContext) MINUS() antlr.TerminalNode

func (s *ExprContext) MapLiteralFields() IMapLiteralFieldsContext

func (s *ExprContext) NEQ() antlr.TerminalNode

func (s *ExprContext) NIL_COALESCE() antlr.TerminalNode

func (s *ExprContext) OK() antlr.TerminalNode

func (s *ExprContext) OVERFLOW_ADD() antlr.TerminalNode

func (s *ExprContext) OVERFLOW_MUL() antlr.TerminalNode

func (s *ExprContext) OVERFLOW_SUB() antlr.TerminalNode

func (s *ExprContext) PERCENT() antlr.TerminalNode

func (s *ExprContext) PIPE() antlr.TerminalNode

func (s *ExprContext) PLUS() antlr.TerminalNode

func (s *ExprContext) QUESTION() antlr.TerminalNode

func (s *ExprContext) RBRACE() antlr.TerminalNode

func (s *ExprContext) RBRACKET() antlr.TerminalNode

func (s *ExprContext) REINTERPRET() antlr.TerminalNode

func (s *ExprContext) RESULT() antlr.TerminalNode

func (s *ExprContext) RPAREN() antlr.TerminalNode

func (s *ExprContext) RSHIFT() antlr.TerminalNode

func (s *ExprContext) SLASH() antlr.TerminalNode

func (s *ExprContext) STAR() antlr.TerminalNode

func (s *ExprContext) StructLiteralFields() IStructLiteralFieldsContext

func (s *ExprContext) TILDE() antlr.TerminalNode

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ExprContext) TypeExpr(i int) ITypeExprContext

type ExprOrAssignStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExprOrAssignStmtContext() *ExprOrAssignStmtContext

func NewExprOrAssignStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprOrAssignStmtContext

func (s *ExprOrAssignStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExprOrAssignStmtContext) AllExpr() []IExprContext

func (s *ExprOrAssignStmtContext) AssignOp() IAssignOpContext

func (s *ExprOrAssignStmtContext) Expr(i int) IExprContext

func (s *ExprOrAssignStmtContext) GetParser() antlr.Parser

func (s *ExprOrAssignStmtContext) GetRuleContext() antlr.RuleContext

func (*ExprOrAssignStmtContext) IsExprOrAssignStmtContext()

func (s *ExprOrAssignStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FileContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFileContext() *FileContext

func NewFileContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FileContext

func (s *FileContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FileContext) AllBuildDecl() []IBuildDeclContext

func (s *FileContext) AllImportDecl() []IImportDeclContext

func (s *FileContext) AllTopLevelDecl() []ITopLevelDeclContext

func (s *FileContext) BuildDecl(i int) IBuildDeclContext

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

func (s *ForInStmtContext) IDENTIFIER() antlr.TerminalNode

func (s *ForInStmtContext) IN() antlr.TerminalNode

func (*ForInStmtContext) IsForInStmtContext()

func (s *ForInStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FuncDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFuncDeclContext() *FuncDeclContext

func NewFuncDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncDeclContext

func (s *FuncDeclContext) ARROW() antlr.TerminalNode

func (s *FuncDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FuncDeclContext) Block() IBlockContext

func (s *FuncDeclContext) FUNC() antlr.TerminalNode

func (s *FuncDeclContext) FuncQualifier() IFuncQualifierContext

func (s *FuncDeclContext) GenericParams() IGenericParamsContext

func (s *FuncDeclContext) GetParser() antlr.Parser

func (s *FuncDeclContext) GetRuleContext() antlr.RuleContext

func (s *FuncDeclContext) IDENTIFIER() antlr.TerminalNode

func (*FuncDeclContext) IsFuncDeclContext()

func (s *FuncDeclContext) LPAREN() antlr.TerminalNode

func (s *FuncDeclContext) ParamList() IParamListContext

func (s *FuncDeclContext) RPAREN() antlr.TerminalNode

func (s *FuncDeclContext) Receiver() IReceiverContext

func (s *FuncDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FuncDeclContext) TypeExpr() ITypeExprContext

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

func (s *FuncQualifierContext) TEST() antlr.TerminalNode

func (s *FuncQualifierContext) THREAD() antlr.TerminalNode

func (s *FuncQualifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FuncTypeParamsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFuncTypeParamsContext() *FuncTypeParamsContext

func NewFuncTypeParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncTypeParamsContext

func (s *FuncTypeParamsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FuncTypeParamsContext) AllCOMMA() []antlr.TerminalNode

func (s *FuncTypeParamsContext) AllTypeExpr() []ITypeExprContext

func (s *FuncTypeParamsContext) COMMA(i int) antlr.TerminalNode

func (s *FuncTypeParamsContext) GetParser() antlr.Parser

func (s *FuncTypeParamsContext) GetRuleContext() antlr.RuleContext

func (*FuncTypeParamsContext) IsFuncTypeParamsContext()

func (s *FuncTypeParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FuncTypeParamsContext) TypeExpr(i int) ITypeExprContext

type GenericParamsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericParamsContext() *GenericParamsContext

func NewGenericParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericParamsContext

func (s *GenericParamsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericParamsContext) AllCOMMA() []antlr.TerminalNode

func (s *GenericParamsContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *GenericParamsContext) COMMA(i int) antlr.TerminalNode

func (s *GenericParamsContext) GT() antlr.TerminalNode

func (s *GenericParamsContext) GetParser() antlr.Parser

func (s *GenericParamsContext) GetRuleContext() antlr.RuleContext

func (s *GenericParamsContext) IDENTIFIER(i int) antlr.TerminalNode

func (*GenericParamsContext) IsGenericParamsContext()

func (s *GenericParamsContext) LT() antlr.TerminalNode

func (s *GenericParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

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
	FuncQualifier() IFuncQualifierContext
	ARROW() antlr.TerminalNode
	TypeExpr() ITypeExprContext

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
	IDENTIFIER() antlr.TerminalNode
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

type IAsmBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAsmInstr() []IAsmInstrContext
	AsmInstr(i int) IAsmInstrContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	AllAsmConstraint() []IAsmConstraintContext
	AsmConstraint(i int) IAsmConstraintContext
	AsmClobberDecl() IAsmClobberDeclContext

	// IsAsmBodyContext differentiates from other interfaces.
	IsAsmBodyContext()
}
    IAsmBodyContext is an interface to support dynamic dispatch.

type IAsmClobberDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CLOBBER() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AllSTRING_LIT() []antlr.TerminalNode
	STRING_LIT(i int) antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsAsmClobberDeclContext differentiates from other interfaces.
	IsAsmClobberDeclContext()
}
    IAsmClobberDeclContext is an interface to support dynamic dispatch.

type IAsmConstraintContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IN() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	STRING_LIT() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	INOUT() antlr.TerminalNode
	OUT_KW() antlr.TerminalNode

	// IsAsmConstraintContext differentiates from other interfaces.
	IsAsmConstraintContext()
}
    IAsmConstraintContext is an interface to support dynamic dispatch.

type IAsmExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASM() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AsmBody() IAsmBodyContext
	RPAREN() antlr.TerminalNode

	// IsAsmExprContext differentiates from other interfaces.
	IsAsmExprContext()
}
    IAsmExprContext is an interface to support dynamic dispatch.

type IAsmInstrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING_LIT() antlr.TerminalNode

	// IsAsmInstrContext differentiates from other interfaces.
	IsAsmInstrContext()
}
    IAsmInstrContext is an interface to support dynamic dispatch.

type IAssignOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASSIGN() antlr.TerminalNode
	PLUS_ASSIGN() antlr.TerminalNode
	MINUS_ASSIGN() antlr.TerminalNode
	STAR_ASSIGN() antlr.TerminalNode
	DIV_ASSIGN() antlr.TerminalNode
	MOD_ASSIGN() antlr.TerminalNode

	// IsAssignOpContext differentiates from other interfaces.
	IsAssignOpContext()
}
    IAssignOpContext is an interface to support dynamic dispatch.

type IBaseTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsBaseTypeContext differentiates from other interfaces.
	IsBaseTypeContext()
}
    IBaseTypeContext is an interface to support dynamic dispatch.

type IBindingPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsBindingPatternContext differentiates from other interfaces.
	IsBindingPatternContext()
}
    IBindingPatternContext is an interface to support dynamic dispatch.

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

type IBuildDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BUILD() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	TEST() antlr.TerminalNode

	// IsBuildDeclContext differentiates from other interfaces.
	IsBuildDeclContext()
}
    IBuildDeclContext is an interface to support dynamic dispatch.

type IClassDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CLASS() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	COLON() antlr.TerminalNode
	QualifiedIdent() IQualifiedIdentContext
	AllClassMember() []IClassMemberContext
	ClassMember(i int) IClassMemberContext

	// IsClassDeclContext differentiates from other interfaces.
	IsClassDeclContext()
}
    IClassDeclContext is an interface to support dynamic dispatch.

type IClassMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	FUNC() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ParamList() IParamListContext
	ARROW() antlr.TerminalNode

	// IsClassMemberContext differentiates from other interfaces.
	IsClassMemberContext()
}
    IClassMemberContext is an interface to support dynamic dispatch.

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

type IEnumCaseDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	AllEnumCaseItem() []IEnumCaseItemContext
	EnumCaseItem(i int) IEnumCaseItemContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsEnumCaseDeclContext differentiates from other interfaces.
	IsEnumCaseDeclContext()
}
    IEnumCaseDeclContext is an interface to support dynamic dispatch.

type IEnumCaseItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	EnumRawValue() IEnumRawValueContext

	// IsEnumCaseItemContext differentiates from other interfaces.
	IsEnumCaseItemContext()
}
    IEnumCaseItemContext is an interface to support dynamic dispatch.

type IEnumDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ENUM() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	COLON() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	AllEnumCaseDecl() []IEnumCaseDeclContext
	EnumCaseDecl(i int) IEnumCaseDeclContext

	// IsEnumDeclContext differentiates from other interfaces.
	IsEnumDeclContext()
}
    IEnumDeclContext is an interface to support dynamic dispatch.

type IEnumRawValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEC_INT_LIT() antlr.TerminalNode
	HEX_INT_LIT() antlr.TerminalNode
	OCT_INT_LIT() antlr.TerminalNode
	BIN_INT_LIT() antlr.TerminalNode
	MINUS() antlr.TerminalNode
	STRING_LIT() antlr.TerminalNode

	// IsEnumRawValueContext differentiates from other interfaces.
	IsEnumRawValueContext()
}
    IEnumRawValueContext is an interface to support dynamic dispatch.

type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	MINUS() antlr.TerminalNode
	BANG() antlr.TerminalNode
	TILDE() antlr.TerminalNode
	AMP() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	StructLiteralFields() IStructLiteralFieldsContext
	MapLiteralFields() IMapLiteralFieldsContext
	MAP() antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	AllTypeExpr() []ITypeExprContext
	TypeExpr(i int) ITypeExprContext
	RBRACKET() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ArgList() IArgListContext
	Literal() ILiteralContext
	DOT() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	AnonFuncExpr() IAnonFuncExprContext
	RESULT() antlr.TerminalNode
	OK() antlr.TerminalNode
	ERR_KW() antlr.TerminalNode
	REINTERPRET() antlr.TerminalNode
	LT() antlr.TerminalNode
	GT() antlr.TerminalNode
	AsmExpr() IAsmExprContext
	LSHIFT() antlr.TerminalNode
	RSHIFT() antlr.TerminalNode
	STAR() antlr.TerminalNode
	SLASH() antlr.TerminalNode
	PERCENT() antlr.TerminalNode
	OVERFLOW_MUL() antlr.TerminalNode
	PLUS() antlr.TerminalNode
	OVERFLOW_ADD() antlr.TerminalNode
	OVERFLOW_SUB() antlr.TerminalNode
	CARET() antlr.TerminalNode
	PIPE() antlr.TerminalNode
	ELLIPSIS() antlr.TerminalNode
	HALF_OPEN() antlr.TerminalNode
	NIL_COALESCE() antlr.TerminalNode
	EQ() antlr.TerminalNode
	NEQ() antlr.TerminalNode
	LEQ() antlr.TerminalNode
	GEQ() antlr.TerminalNode
	IDENTITY_EQ() antlr.TerminalNode
	IDENTITY_NEQ() antlr.TerminalNode
	LOGICAL_AND() antlr.TerminalNode
	LOGICAL_OR() antlr.TerminalNode
	QUESTION() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}
    IExprContext is an interface to support dynamic dispatch.

type IExprOrAssignStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	AssignOp() IAssignOpContext

	// IsExprOrAssignStmtContext differentiates from other interfaces.
	IsExprOrAssignStmtContext()
}
    IExprOrAssignStmtContext is an interface to support dynamic dispatch.

type IFileContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	PackageDecl() IPackageDeclContext
	AllBuildDecl() []IBuildDeclContext
	BuildDecl(i int) IBuildDeclContext
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
	IDENTIFIER() antlr.TerminalNode
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
	IDENTIFIER() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	Block() IBlockContext
	Receiver() IReceiverContext
	GenericParams() IGenericParamsContext
	ParamList() IParamListContext
	FuncQualifier() IFuncQualifierContext
	ARROW() antlr.TerminalNode
	TypeExpr() ITypeExprContext

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
	TEST() antlr.TerminalNode

	// IsFuncQualifierContext differentiates from other interfaces.
	IsFuncQualifierContext()
}
    IFuncQualifierContext is an interface to support dynamic dispatch.

type IFuncTypeParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTypeExpr() []ITypeExprContext
	TypeExpr(i int) ITypeExprContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFuncTypeParamsContext differentiates from other interfaces.
	IsFuncTypeParamsContext()
}
    IFuncTypeParamsContext is an interface to support dynamic dispatch.

type IGenericParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LT() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
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
	IDENTIFIER() antlr.TerminalNode
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
	AllBlock() []IBlockContext
	Block(i int) IBlockContext
	ELSE() antlr.TerminalNode
	IfStmt() IIfStmtContext

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
	OCT_INT_LIT() antlr.TerminalNode
	BIN_INT_LIT() antlr.TerminalNode
	DEC_FLOAT_LIT() antlr.TerminalNode
	HEX_FLOAT_LIT() antlr.TerminalNode
	CHAR_LIT() antlr.TerminalNode
	STRING_LIT() antlr.TerminalNode
	MULTILINE_STRING_LIT() antlr.TerminalNode
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode
	NIL() antlr.TerminalNode

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}
    ILiteralContext is an interface to support dynamic dispatch.

type IMapLiteralFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	COLON() antlr.TerminalNode

	// IsMapLiteralFieldContext differentiates from other interfaces.
	IsMapLiteralFieldContext()
}
    IMapLiteralFieldContext is an interface to support dynamic dispatch.

type IMapLiteralFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllMapLiteralField() []IMapLiteralFieldContext
	MapLiteralField(i int) IMapLiteralFieldContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsMapLiteralFieldsContext differentiates from other interfaces.
	IsMapLiteralFieldsContext()
}
    IMapLiteralFieldsContext is an interface to support dynamic dispatch.

type IPackageDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PACKAGE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsPackageDeclContext differentiates from other interfaces.
	IsPackageDeclContext()
}
    IPackageDeclContext is an interface to support dynamic dispatch.

type IParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	TypeExpr() ITypeExprContext

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
	VariadicParam() IVariadicParamContext

	// IsParamListContext differentiates from other interfaces.
	IsParamListContext()
}
    IParamListContext is an interface to support dynamic dispatch.

type IQualifiedIdentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsQualifiedIdentContext differentiates from other interfaces.
	IsQualifiedIdentContext()
}
    IQualifiedIdentContext is an interface to support dynamic dispatch.

type IReceiverContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	RPAREN() antlr.TerminalNode

	// IsReceiverContext differentiates from other interfaces.
	IsReceiverContext()
}
    IReceiverContext is an interface to support dynamic dispatch.

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

type IStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VarDecl() IVarDeclContext
	IfStmt() IIfStmtContext
	WhileStmt() IWhileStmtContext
	ForInStmt() IForInStmtContext
	SwitchStmt() ISwitchStmtContext
	ReturnStmt() IReturnStmtContext
	BREAK() antlr.TerminalNode
	CONTINUE() antlr.TerminalNode
	FALLTHROUGH() antlr.TerminalNode
	DeferStmt() IDeferStmtContext
	StructDecl() IStructDeclContext
	ClassDecl() IClassDeclContext
	EnumDecl() IEnumDeclContext
	TypeAliasDecl() ITypeAliasDeclContext
	ExprOrAssignStmt() IExprOrAssignStmtContext

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
	IDENTIFIER() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllStructFieldDecl() []IStructFieldDeclContext
	StructFieldDecl(i int) IStructFieldDeclContext

	// IsStructDeclContext differentiates from other interfaces.
	IsStructDeclContext()
}
    IStructDeclContext is an interface to support dynamic dispatch.

type IStructFieldDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	TypeExpr() ITypeExprContext

	// IsStructFieldDeclContext differentiates from other interfaces.
	IsStructFieldDeclContext()
}
    IStructFieldDeclContext is an interface to support dynamic dispatch.

type IStructLiteralFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Expr() IExprContext

	// IsStructLiteralFieldContext differentiates from other interfaces.
	IsStructLiteralFieldContext()
}
    IStructLiteralFieldContext is an interface to support dynamic dispatch.

type IStructLiteralFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllStructLiteralField() []IStructLiteralFieldContext
	StructLiteralField(i int) IStructLiteralFieldContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsStructLiteralFieldsContext differentiates from other interfaces.
	IsStructLiteralFieldsContext()
}
    IStructLiteralFieldsContext is an interface to support dynamic dispatch.

type ISwitchCaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	AllSwitchPattern() []ISwitchPatternContext
	SwitchPattern(i int) ISwitchPatternContext
	COLON() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	AllStmt() []IStmtContext
	Stmt(i int) IStmtContext
	DEFAULT() antlr.TerminalNode

	// IsSwitchCaseContext differentiates from other interfaces.
	IsSwitchCaseContext()
}
    ISwitchCaseContext is an interface to support dynamic dispatch.

type ISwitchPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expr() IExprContext
	DOT() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	LET() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	OK() antlr.TerminalNode
	ERR_KW() antlr.TerminalNode

	// IsSwitchPatternContext differentiates from other interfaces.
	IsSwitchPatternContext()
}
    ISwitchPatternContext is an interface to support dynamic dispatch.

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
	VarDecl() IVarDeclContext

	// IsTopLevelDeclContext differentiates from other interfaces.
	IsTopLevelDeclContext()
}
    ITopLevelDeclContext is an interface to support dynamic dispatch.

type ITupleTypeElemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TypeExpr() ITypeExprContext
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsTupleTypeElemContext differentiates from other interfaces.
	IsTupleTypeElemContext()
}
    ITupleTypeElemContext is an interface to support dynamic dispatch.

type ITupleTypeElemsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTupleTypeElem() []ITupleTypeElemContext
	TupleTypeElem(i int) ITupleTypeElemContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTupleTypeElemsContext differentiates from other interfaces.
	IsTupleTypeElemsContext()
}
    ITupleTypeElemsContext is an interface to support dynamic dispatch.

type ITypeAliasDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	TypeExpr() ITypeExprContext

	// IsTypeAliasDeclContext differentiates from other interfaces.
	IsTypeAliasDeclContext()
}
    ITypeAliasDeclContext is an interface to support dynamic dispatch.

type ITypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STAR() antlr.TerminalNode
	AllTypeExpr() []ITypeExprContext
	TypeExpr(i int) ITypeExprContext
	CONST_KW() antlr.TerminalNode
	QUESTION() antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	RBRACKET() antlr.TerminalNode
	MAP() antlr.TerminalNode
	CHAN() antlr.TerminalNode
	FUNC() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	FuncTypeParams() IFuncTypeParamsContext
	ARROW() antlr.TerminalNode
	TupleTypeElems() ITupleTypeElemsContext
	RESULT() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	EXPECTED() antlr.TerminalNode
	STRING_LIT() antlr.TerminalNode
	BaseType() IBaseTypeContext

	// IsTypeExprContext differentiates from other interfaces.
	IsTypeExprContext()
}
    ITypeExprContext is an interface to support dynamic dispatch.

type IVarDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BindingPattern() IBindingPatternContext
	ASSIGN() antlr.TerminalNode
	Expr() IExprContext
	LET() antlr.TerminalNode
	VAR() antlr.TerminalNode
	WEAK() antlr.TerminalNode
	COLON() antlr.TerminalNode
	TypeExpr() ITypeExprContext

	// IsVarDeclContext differentiates from other interfaces.
	IsVarDeclContext()
}
    IVarDeclContext is an interface to support dynamic dispatch.

type IVariadicParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	ELLIPSIS() antlr.TerminalNode
	TypeExpr() ITypeExprContext

	// IsVariadicParamContext differentiates from other interfaces.
	IsVariadicParamContext()
}
    IVariadicParamContext is an interface to support dynamic dispatch.

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

func (s *IfConditionContext) IDENTIFIER() antlr.TerminalNode

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

func (s *IfStmtContext) AllBlock() []IBlockContext

func (s *IfStmtContext) Block(i int) IBlockContext

func (s *IfStmtContext) ELSE() antlr.TerminalNode

func (s *IfStmtContext) GetParser() antlr.Parser

func (s *IfStmtContext) GetRuleContext() antlr.RuleContext

func (s *IfStmtContext) IF() antlr.TerminalNode

func (s *IfStmtContext) IfCondition() IIfConditionContext

func (s *IfStmtContext) IfStmt() IIfStmtContext

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

func (s *LiteralContext) CHAR_LIT() antlr.TerminalNode

func (s *LiteralContext) DEC_FLOAT_LIT() antlr.TerminalNode

func (s *LiteralContext) DEC_INT_LIT() antlr.TerminalNode

func (s *LiteralContext) FALSE() antlr.TerminalNode

func (s *LiteralContext) GetParser() antlr.Parser

func (s *LiteralContext) GetRuleContext() antlr.RuleContext

func (s *LiteralContext) HEX_FLOAT_LIT() antlr.TerminalNode

func (s *LiteralContext) HEX_INT_LIT() antlr.TerminalNode

func (*LiteralContext) IsLiteralContext()

func (s *LiteralContext) MULTILINE_STRING_LIT() antlr.TerminalNode

func (s *LiteralContext) NIL() antlr.TerminalNode

func (s *LiteralContext) OCT_INT_LIT() antlr.TerminalNode

func (s *LiteralContext) STRING_LIT() antlr.TerminalNode

func (s *LiteralContext) TRUE() antlr.TerminalNode

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type MapLiteralFieldContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyMapLiteralFieldContext() *MapLiteralFieldContext

func NewMapLiteralFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapLiteralFieldContext

func (s *MapLiteralFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MapLiteralFieldContext) AllExpr() []IExprContext

func (s *MapLiteralFieldContext) COLON() antlr.TerminalNode

func (s *MapLiteralFieldContext) Expr(i int) IExprContext

func (s *MapLiteralFieldContext) GetParser() antlr.Parser

func (s *MapLiteralFieldContext) GetRuleContext() antlr.RuleContext

func (*MapLiteralFieldContext) IsMapLiteralFieldContext()

func (s *MapLiteralFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type MapLiteralFieldsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyMapLiteralFieldsContext() *MapLiteralFieldsContext

func NewMapLiteralFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapLiteralFieldsContext

func (s *MapLiteralFieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MapLiteralFieldsContext) AllCOMMA() []antlr.TerminalNode

func (s *MapLiteralFieldsContext) AllMapLiteralField() []IMapLiteralFieldContext

func (s *MapLiteralFieldsContext) COMMA(i int) antlr.TerminalNode

func (s *MapLiteralFieldsContext) GetParser() antlr.Parser

func (s *MapLiteralFieldsContext) GetRuleContext() antlr.RuleContext

func (*MapLiteralFieldsContext) IsMapLiteralFieldsContext()

func (s *MapLiteralFieldsContext) MapLiteralField(i int) IMapLiteralFieldContext

func (s *MapLiteralFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PackageDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPackageDeclContext() *PackageDeclContext

func NewPackageDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PackageDeclContext

func (s *PackageDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PackageDeclContext) GetParser() antlr.Parser

func (s *PackageDeclContext) GetRuleContext() antlr.RuleContext

func (s *PackageDeclContext) IDENTIFIER() antlr.TerminalNode

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

func (s *ParamContext) IDENTIFIER() antlr.TerminalNode

func (*ParamContext) IsParamContext()

func (s *ParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ParamContext) TypeExpr() ITypeExprContext

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

func (s *ParamListContext) VariadicParam() IVariadicParamContext

type QualifiedIdentContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyQualifiedIdentContext() *QualifiedIdentContext

func NewQualifiedIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QualifiedIdentContext

func (s *QualifiedIdentContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *QualifiedIdentContext) AllDOT() []antlr.TerminalNode

func (s *QualifiedIdentContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *QualifiedIdentContext) DOT(i int) antlr.TerminalNode

func (s *QualifiedIdentContext) GetParser() antlr.Parser

func (s *QualifiedIdentContext) GetRuleContext() antlr.RuleContext

func (s *QualifiedIdentContext) IDENTIFIER(i int) antlr.TerminalNode

func (*QualifiedIdentContext) IsQualifiedIdentContext()

func (s *QualifiedIdentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ReceiverContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyReceiverContext() *ReceiverContext

func NewReceiverContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReceiverContext

func (s *ReceiverContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ReceiverContext) COLON() antlr.TerminalNode

func (s *ReceiverContext) GetParser() antlr.Parser

func (s *ReceiverContext) GetRuleContext() antlr.RuleContext

func (s *ReceiverContext) IDENTIFIER() antlr.TerminalNode

func (*ReceiverContext) IsReceiverContext()

func (s *ReceiverContext) LPAREN() antlr.TerminalNode

func (s *ReceiverContext) RPAREN() antlr.TerminalNode

func (s *ReceiverContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ReceiverContext) TypeExpr() ITypeExprContext

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

type StmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStmtContext() *StmtContext

func NewStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtContext

func (s *StmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StmtContext) BREAK() antlr.TerminalNode

func (s *StmtContext) CONTINUE() antlr.TerminalNode

func (s *StmtContext) ClassDecl() IClassDeclContext

func (s *StmtContext) DeferStmt() IDeferStmtContext

func (s *StmtContext) EnumDecl() IEnumDeclContext

func (s *StmtContext) ExprOrAssignStmt() IExprOrAssignStmtContext

func (s *StmtContext) FALLTHROUGH() antlr.TerminalNode

func (s *StmtContext) ForInStmt() IForInStmtContext

func (s *StmtContext) GetParser() antlr.Parser

func (s *StmtContext) GetRuleContext() antlr.RuleContext

func (s *StmtContext) IfStmt() IIfStmtContext

func (*StmtContext) IsStmtContext()

func (s *StmtContext) ReturnStmt() IReturnStmtContext

func (s *StmtContext) StructDecl() IStructDeclContext

func (s *StmtContext) SwitchStmt() ISwitchStmtContext

func (s *StmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *StmtContext) TypeAliasDecl() ITypeAliasDeclContext

func (s *StmtContext) VarDecl() IVarDeclContext

func (s *StmtContext) WhileStmt() IWhileStmtContext

type StructDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructDeclContext() *StructDeclContext

func NewStructDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructDeclContext

func (s *StructDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructDeclContext) AllStructFieldDecl() []IStructFieldDeclContext

func (s *StructDeclContext) GetParser() antlr.Parser

func (s *StructDeclContext) GetRuleContext() antlr.RuleContext

func (s *StructDeclContext) IDENTIFIER() antlr.TerminalNode

func (*StructDeclContext) IsStructDeclContext()

func (s *StructDeclContext) LBRACE() antlr.TerminalNode

func (s *StructDeclContext) RBRACE() antlr.TerminalNode

func (s *StructDeclContext) STRUCT() antlr.TerminalNode

func (s *StructDeclContext) StructFieldDecl(i int) IStructFieldDeclContext

func (s *StructDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type StructFieldDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructFieldDeclContext() *StructFieldDeclContext

func NewStructFieldDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructFieldDeclContext

func (s *StructFieldDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructFieldDeclContext) COLON() antlr.TerminalNode

func (s *StructFieldDeclContext) GetParser() antlr.Parser

func (s *StructFieldDeclContext) GetRuleContext() antlr.RuleContext

func (s *StructFieldDeclContext) IDENTIFIER() antlr.TerminalNode

func (*StructFieldDeclContext) IsStructFieldDeclContext()

func (s *StructFieldDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *StructFieldDeclContext) TypeExpr() ITypeExprContext

type StructLiteralFieldContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructLiteralFieldContext() *StructLiteralFieldContext

func NewStructLiteralFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructLiteralFieldContext

func (s *StructLiteralFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructLiteralFieldContext) COLON() antlr.TerminalNode

func (s *StructLiteralFieldContext) Expr() IExprContext

func (s *StructLiteralFieldContext) GetParser() antlr.Parser

func (s *StructLiteralFieldContext) GetRuleContext() antlr.RuleContext

func (s *StructLiteralFieldContext) IDENTIFIER() antlr.TerminalNode

func (*StructLiteralFieldContext) IsStructLiteralFieldContext()

func (s *StructLiteralFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type StructLiteralFieldsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructLiteralFieldsContext() *StructLiteralFieldsContext

func NewStructLiteralFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructLiteralFieldsContext

func (s *StructLiteralFieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructLiteralFieldsContext) AllCOMMA() []antlr.TerminalNode

func (s *StructLiteralFieldsContext) AllStructLiteralField() []IStructLiteralFieldContext

func (s *StructLiteralFieldsContext) COMMA(i int) antlr.TerminalNode

func (s *StructLiteralFieldsContext) GetParser() antlr.Parser

func (s *StructLiteralFieldsContext) GetRuleContext() antlr.RuleContext

func (*StructLiteralFieldsContext) IsStructLiteralFieldsContext()

func (s *StructLiteralFieldsContext) StructLiteralField(i int) IStructLiteralFieldContext

func (s *StructLiteralFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SwitchCaseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchCaseContext() *SwitchCaseContext

func NewSwitchCaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchCaseContext

func (s *SwitchCaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchCaseContext) AllCOMMA() []antlr.TerminalNode

func (s *SwitchCaseContext) AllStmt() []IStmtContext

func (s *SwitchCaseContext) AllSwitchPattern() []ISwitchPatternContext

func (s *SwitchCaseContext) CASE() antlr.TerminalNode

func (s *SwitchCaseContext) COLON() antlr.TerminalNode

func (s *SwitchCaseContext) COMMA(i int) antlr.TerminalNode

func (s *SwitchCaseContext) DEFAULT() antlr.TerminalNode

func (s *SwitchCaseContext) GetParser() antlr.Parser

func (s *SwitchCaseContext) GetRuleContext() antlr.RuleContext

func (*SwitchCaseContext) IsSwitchCaseContext()

func (s *SwitchCaseContext) Stmt(i int) IStmtContext

func (s *SwitchCaseContext) SwitchPattern(i int) ISwitchPatternContext

func (s *SwitchCaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SwitchPatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchPatternContext() *SwitchPatternContext

func NewSwitchPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchPatternContext

func (s *SwitchPatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchPatternContext) DOT() antlr.TerminalNode

func (s *SwitchPatternContext) ERR_KW() antlr.TerminalNode

func (s *SwitchPatternContext) Expr() IExprContext

func (s *SwitchPatternContext) GetParser() antlr.Parser

func (s *SwitchPatternContext) GetRuleContext() antlr.RuleContext

func (s *SwitchPatternContext) IDENTIFIER() antlr.TerminalNode

func (*SwitchPatternContext) IsSwitchPatternContext()

func (s *SwitchPatternContext) LET() antlr.TerminalNode

func (s *SwitchPatternContext) LPAREN() antlr.TerminalNode

func (s *SwitchPatternContext) OK() antlr.TerminalNode

func (s *SwitchPatternContext) RPAREN() antlr.TerminalNode

func (s *SwitchPatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SwitchStmtContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchStmtContext() *SwitchStmtContext

func NewSwitchStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchStmtContext

func (s *SwitchStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchStmtContext) AllSwitchCase() []ISwitchCaseContext

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

func (s *TopLevelDeclContext) VarDecl() IVarDeclContext

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

func (s *TupleTypeElemContext) IDENTIFIER() antlr.TerminalNode

func (*TupleTypeElemContext) IsTupleTypeElemContext()

func (s *TupleTypeElemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleTypeElemContext) TypeExpr() ITypeExprContext

type TupleTypeElemsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleTypeElemsContext() *TupleTypeElemsContext

func NewTupleTypeElemsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleTypeElemsContext

func (s *TupleTypeElemsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleTypeElemsContext) AllCOMMA() []antlr.TerminalNode

func (s *TupleTypeElemsContext) AllTupleTypeElem() []ITupleTypeElemContext

func (s *TupleTypeElemsContext) COMMA(i int) antlr.TerminalNode

func (s *TupleTypeElemsContext) GetParser() antlr.Parser

func (s *TupleTypeElemsContext) GetRuleContext() antlr.RuleContext

func (*TupleTypeElemsContext) IsTupleTypeElemsContext()

func (s *TupleTypeElemsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleTypeElemsContext) TupleTypeElem(i int) ITupleTypeElemContext

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

func (s *TypeAliasDeclContext) IDENTIFIER() antlr.TerminalNode

func (*TypeAliasDeclContext) IsTypeAliasDeclContext()

func (s *TypeAliasDeclContext) TYPE() antlr.TerminalNode

func (s *TypeAliasDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeAliasDeclContext) TypeExpr() ITypeExprContext

type TypeExprContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeExprContext() *TypeExprContext

func NewTypeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeExprContext

func (s *TypeExprContext) ARROW() antlr.TerminalNode

func (s *TypeExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeExprContext) AllTypeExpr() []ITypeExprContext

func (s *TypeExprContext) BaseType() IBaseTypeContext

func (s *TypeExprContext) CHAN() antlr.TerminalNode

func (s *TypeExprContext) COMMA() antlr.TerminalNode

func (s *TypeExprContext) CONST_KW() antlr.TerminalNode

func (s *TypeExprContext) EXPECTED() antlr.TerminalNode

func (s *TypeExprContext) FUNC() antlr.TerminalNode

func (s *TypeExprContext) FuncTypeParams() IFuncTypeParamsContext

func (s *TypeExprContext) GetParser() antlr.Parser

func (s *TypeExprContext) GetRuleContext() antlr.RuleContext

func (*TypeExprContext) IsTypeExprContext()

func (s *TypeExprContext) LBRACKET() antlr.TerminalNode

func (s *TypeExprContext) LPAREN() antlr.TerminalNode

func (s *TypeExprContext) MAP() antlr.TerminalNode

func (s *TypeExprContext) QUESTION() antlr.TerminalNode

func (s *TypeExprContext) RBRACKET() antlr.TerminalNode

func (s *TypeExprContext) RESULT() antlr.TerminalNode

func (s *TypeExprContext) RPAREN() antlr.TerminalNode

func (s *TypeExprContext) STAR() antlr.TerminalNode

func (s *TypeExprContext) STRING_LIT() antlr.TerminalNode

func (s *TypeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeExprContext) TupleTypeElems() ITupleTypeElemsContext

func (s *TypeExprContext) TypeExpr(i int) ITypeExprContext

type VarDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyVarDeclContext() *VarDeclContext

func NewVarDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarDeclContext

func (s *VarDeclContext) ASSIGN() antlr.TerminalNode

func (s *VarDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *VarDeclContext) BindingPattern() IBindingPatternContext

func (s *VarDeclContext) COLON() antlr.TerminalNode

func (s *VarDeclContext) Expr() IExprContext

func (s *VarDeclContext) GetParser() antlr.Parser

func (s *VarDeclContext) GetRuleContext() antlr.RuleContext

func (*VarDeclContext) IsVarDeclContext()

func (s *VarDeclContext) LET() antlr.TerminalNode

func (s *VarDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *VarDeclContext) TypeExpr() ITypeExprContext

func (s *VarDeclContext) VAR() antlr.TerminalNode

func (s *VarDeclContext) WEAK() antlr.TerminalNode

type VariadicParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyVariadicParamContext() *VariadicParamContext

func NewVariadicParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariadicParamContext

func (s *VariadicParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *VariadicParamContext) COLON() antlr.TerminalNode

func (s *VariadicParamContext) ELLIPSIS() antlr.TerminalNode

func (s *VariadicParamContext) GetParser() antlr.Parser

func (s *VariadicParamContext) GetRuleContext() antlr.RuleContext

func (s *VariadicParamContext) IDENTIFIER() antlr.TerminalNode

func (*VariadicParamContext) IsVariadicParamContext()

func (s *VariadicParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *VariadicParamContext) TypeExpr() ITypeExprContext

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

func (p *VertexParser) AsmBody() (localctx IAsmBodyContext)

func (p *VertexParser) AsmClobberDecl() (localctx IAsmClobberDeclContext)

func (p *VertexParser) AsmConstraint() (localctx IAsmConstraintContext)

func (p *VertexParser) AsmExpr() (localctx IAsmExprContext)

func (p *VertexParser) AsmInstr() (localctx IAsmInstrContext)

func (p *VertexParser) AssignOp() (localctx IAssignOpContext)

func (p *VertexParser) BaseType() (localctx IBaseTypeContext)

func (p *VertexParser) BindingPattern() (localctx IBindingPatternContext)

func (p *VertexParser) Block() (localctx IBlockContext)

func (p *VertexParser) BuildDecl() (localctx IBuildDeclContext)

func (p *VertexParser) ClassDecl() (localctx IClassDeclContext)

func (p *VertexParser) ClassMember() (localctx IClassMemberContext)

func (p *VertexParser) DeferStmt() (localctx IDeferStmtContext)

func (p *VertexParser) EnumCaseDecl() (localctx IEnumCaseDeclContext)

func (p *VertexParser) EnumCaseItem() (localctx IEnumCaseItemContext)

func (p *VertexParser) EnumDecl() (localctx IEnumDeclContext)

func (p *VertexParser) EnumRawValue() (localctx IEnumRawValueContext)

func (p *VertexParser) Expr() (localctx IExprContext)

func (p *VertexParser) ExprOrAssignStmt() (localctx IExprOrAssignStmtContext)

func (p *VertexParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) File() (localctx IFileContext)

func (p *VertexParser) ForInStmt() (localctx IForInStmtContext)

func (p *VertexParser) FuncDecl() (localctx IFuncDeclContext)

func (p *VertexParser) FuncQualifier() (localctx IFuncQualifierContext)

func (p *VertexParser) FuncTypeParams() (localctx IFuncTypeParamsContext)

func (p *VertexParser) GenericParams() (localctx IGenericParamsContext)

func (p *VertexParser) IfCondition() (localctx IIfConditionContext)

func (p *VertexParser) IfStmt() (localctx IIfStmtContext)

func (p *VertexParser) ImportDecl() (localctx IImportDeclContext)

func (p *VertexParser) Literal() (localctx ILiteralContext)

func (p *VertexParser) MapLiteralField() (localctx IMapLiteralFieldContext)

func (p *VertexParser) MapLiteralFields() (localctx IMapLiteralFieldsContext)

func (p *VertexParser) PackageDecl() (localctx IPackageDeclContext)

func (p *VertexParser) Param() (localctx IParamContext)

func (p *VertexParser) ParamList() (localctx IParamListContext)

func (p *VertexParser) QualifiedIdent() (localctx IQualifiedIdentContext)

func (p *VertexParser) Receiver() (localctx IReceiverContext)

func (p *VertexParser) ReturnStmt() (localctx IReturnStmtContext)

func (p *VertexParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool

func (p *VertexParser) Stmt() (localctx IStmtContext)

func (p *VertexParser) StructDecl() (localctx IStructDeclContext)

func (p *VertexParser) StructFieldDecl() (localctx IStructFieldDeclContext)

func (p *VertexParser) StructLiteralField() (localctx IStructLiteralFieldContext)

func (p *VertexParser) StructLiteralFields() (localctx IStructLiteralFieldsContext)

func (p *VertexParser) SwitchCase() (localctx ISwitchCaseContext)

func (p *VertexParser) SwitchPattern() (localctx ISwitchPatternContext)

func (p *VertexParser) SwitchStmt() (localctx ISwitchStmtContext)

func (p *VertexParser) TopLevelDecl() (localctx ITopLevelDeclContext)

func (p *VertexParser) TupleTypeElem() (localctx ITupleTypeElemContext)

func (p *VertexParser) TupleTypeElems() (localctx ITupleTypeElemsContext)

func (p *VertexParser) TypeAliasDecl() (localctx ITypeAliasDeclContext)

func (p *VertexParser) TypeExpr() (localctx ITypeExprContext)

func (p *VertexParser) VarDecl() (localctx IVarDeclContext)

func (p *VertexParser) VariadicParam() (localctx IVariadicParamContext)

func (p *VertexParser) WhileStmt() (localctx IWhileStmtContext)

type VertexParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by VertexParser#file.
	VisitFile(ctx *FileContext) interface{}

	// Visit a parse tree produced by VertexParser#packageDecl.
	VisitPackageDecl(ctx *PackageDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#buildDecl.
	VisitBuildDecl(ctx *BuildDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#importDecl.
	VisitImportDecl(ctx *ImportDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#topLevelDecl.
	VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#typeAliasDecl.
	VisitTypeAliasDecl(ctx *TypeAliasDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#funcDecl.
	VisitFuncDecl(ctx *FuncDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#receiver.
	VisitReceiver(ctx *ReceiverContext) interface{}

	// Visit a parse tree produced by VertexParser#genericParams.
	VisitGenericParams(ctx *GenericParamsContext) interface{}

	// Visit a parse tree produced by VertexParser#funcQualifier.
	VisitFuncQualifier(ctx *FuncQualifierContext) interface{}

	// Visit a parse tree produced by VertexParser#paramList.
	VisitParamList(ctx *ParamListContext) interface{}

	// Visit a parse tree produced by VertexParser#param.
	VisitParam(ctx *ParamContext) interface{}

	// Visit a parse tree produced by VertexParser#variadicParam.
	VisitVariadicParam(ctx *VariadicParamContext) interface{}

	// Visit a parse tree produced by VertexParser#structDecl.
	VisitStructDecl(ctx *StructDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#structFieldDecl.
	VisitStructFieldDecl(ctx *StructFieldDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#classDecl.
	VisitClassDecl(ctx *ClassDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#classMember.
	VisitClassMember(ctx *ClassMemberContext) interface{}

	// Visit a parse tree produced by VertexParser#qualifiedIdent.
	VisitQualifiedIdent(ctx *QualifiedIdentContext) interface{}

	// Visit a parse tree produced by VertexParser#enumDecl.
	VisitEnumDecl(ctx *EnumDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#enumCaseDecl.
	VisitEnumCaseDecl(ctx *EnumCaseDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#enumCaseItem.
	VisitEnumCaseItem(ctx *EnumCaseItemContext) interface{}

	// Visit a parse tree produced by VertexParser#enumRawValue.
	VisitEnumRawValue(ctx *EnumRawValueContext) interface{}

	// Visit a parse tree produced by VertexParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by VertexParser#stmt.
	VisitStmt(ctx *StmtContext) interface{}

	// Visit a parse tree produced by VertexParser#varDecl.
	VisitVarDecl(ctx *VarDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#bindingPattern.
	VisitBindingPattern(ctx *BindingPatternContext) interface{}

	// Visit a parse tree produced by VertexParser#ifStmt.
	VisitIfStmt(ctx *IfStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#ifCondition.
	VisitIfCondition(ctx *IfConditionContext) interface{}

	// Visit a parse tree produced by VertexParser#whileStmt.
	VisitWhileStmt(ctx *WhileStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#forInStmt.
	VisitForInStmt(ctx *ForInStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#switchStmt.
	VisitSwitchStmt(ctx *SwitchStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#switchCase.
	VisitSwitchCase(ctx *SwitchCaseContext) interface{}

	// Visit a parse tree produced by VertexParser#switchPattern.
	VisitSwitchPattern(ctx *SwitchPatternContext) interface{}

	// Visit a parse tree produced by VertexParser#returnStmt.
	VisitReturnStmt(ctx *ReturnStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#deferStmt.
	VisitDeferStmt(ctx *DeferStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#exprOrAssignStmt.
	VisitExprOrAssignStmt(ctx *ExprOrAssignStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#assignOp.
	VisitAssignOp(ctx *AssignOpContext) interface{}

	// Visit a parse tree produced by VertexParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by VertexParser#argList.
	VisitArgList(ctx *ArgListContext) interface{}

	// Visit a parse tree produced by VertexParser#arg.
	VisitArg(ctx *ArgContext) interface{}

	// Visit a parse tree produced by VertexParser#structLiteralFields.
	VisitStructLiteralFields(ctx *StructLiteralFieldsContext) interface{}

	// Visit a parse tree produced by VertexParser#structLiteralField.
	VisitStructLiteralField(ctx *StructLiteralFieldContext) interface{}

	// Visit a parse tree produced by VertexParser#mapLiteralFields.
	VisitMapLiteralFields(ctx *MapLiteralFieldsContext) interface{}

	// Visit a parse tree produced by VertexParser#mapLiteralField.
	VisitMapLiteralField(ctx *MapLiteralFieldContext) interface{}

	// Visit a parse tree produced by VertexParser#anonFuncExpr.
	VisitAnonFuncExpr(ctx *AnonFuncExprContext) interface{}

	// Visit a parse tree produced by VertexParser#asmExpr.
	VisitAsmExpr(ctx *AsmExprContext) interface{}

	// Visit a parse tree produced by VertexParser#asmBody.
	VisitAsmBody(ctx *AsmBodyContext) interface{}

	// Visit a parse tree produced by VertexParser#asmInstr.
	VisitAsmInstr(ctx *AsmInstrContext) interface{}

	// Visit a parse tree produced by VertexParser#asmConstraint.
	VisitAsmConstraint(ctx *AsmConstraintContext) interface{}

	// Visit a parse tree produced by VertexParser#asmClobberDecl.
	VisitAsmClobberDecl(ctx *AsmClobberDeclContext) interface{}

	// Visit a parse tree produced by VertexParser#typeExpr.
	VisitTypeExpr(ctx *TypeExprContext) interface{}

	// Visit a parse tree produced by VertexParser#funcTypeParams.
	VisitFuncTypeParams(ctx *FuncTypeParamsContext) interface{}

	// Visit a parse tree produced by VertexParser#baseType.
	VisitBaseType(ctx *BaseTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleTypeElems.
	VisitTupleTypeElems(ctx *TupleTypeElemsContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleTypeElem.
	VisitTupleTypeElem(ctx *TupleTypeElemContext) interface{}

	// Visit a parse tree produced by VertexParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}
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

