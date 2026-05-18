

CONSTANTS

const (
	VertexLexerWS                       = 1
	VertexLexerNEWLINE                  = 2
	VertexLexerLINE_COMMENT             = 3
	VertexLexerBLOCK_COMMENT            = 4
	VertexLexerASSOCIATEDTYPE           = 5
	VertexLexerCLASS                    = 6
	VertexLexerDEINIT                   = 7
	VertexLexerENUM                     = 8
	VertexLexerEXTENSION                = 9
	VertexLexerFUNC                     = 10
	VertexLexerIMPORT                   = 11
	VertexLexerINIT                     = 12
	VertexLexerINOUT                    = 13
	VertexLexerLET                      = 14
	VertexLexerOPEN                     = 15
	VertexLexerPRIVATE                  = 16
	VertexLexerPROTOCOL                 = 17
	VertexLexerPUBLIC                   = 18
	VertexLexerSTATIC                   = 19
	VertexLexerSTRUCT                   = 20
	VertexLexerSUBSCRIPT                = 21
	VertexLexerTYPEALIAS                = 22
	VertexLexerVAR                      = 23
	VertexLexerINTERNAL                 = 24
	VertexLexerBREAK                    = 25
	VertexLexerCASE                     = 26
	VertexLexerCATCH                    = 27
	VertexLexerCONTINUE                 = 28
	VertexLexerDEFAULT                  = 29
	VertexLexerDEFER                    = 30
	VertexLexerDO                       = 31
	VertexLexerELSE                     = 32
	VertexLexerFALLTHROUGH              = 33
	VertexLexerFOR                      = 34
	VertexLexerGUARD                    = 35
	VertexLexerIF                       = 36
	VertexLexerIN                       = 37
	VertexLexerREPEAT                   = 38
	VertexLexerRETURN                   = 39
	VertexLexerTHROW                    = 40
	VertexLexerSWITCH                   = 41
	VertexLexerWHERE                    = 42
	VertexLexerWHILE                    = 43
	VertexLexerANY                      = 44
	VertexLexerAS                       = 45
	VertexLexerAWAIT                    = 46
	VertexLexerFALSE                    = 47
	VertexLexerIS                       = 48
	VertexLexerNIL                      = 49
	VertexLexerSELF                     = 50
	VertexLexerSELF_UPPER               = 51
	VertexLexerSUPER                    = 52
	VertexLexerTHROWS                   = 53
	VertexLexerTRUE                     = 54
	VertexLexerTRY                      = 55
	VertexLexerCONSUME                  = 56
	VertexLexerCOPY                     = 57
	VertexLexerDISCARD                  = 58
	VertexLexerBORROW                   = 59
	VertexLexerUNDERSCORE               = 60
	VertexLexerOS_KW                    = 61
	VertexLexerARCH_KW                  = 62
	VertexLexerVERTEX_KW                = 63
	VertexLexerCOMPILER_KW              = 64
	VertexLexerCANIMPORT_KW             = 65
	VertexLexerVERSION_KW               = 66
	VertexLexerFILE_KEYWORD             = 67
	VertexLexerLINE_KEYWORD             = 68
	VertexLexerGET_KW                   = 69
	VertexLexerSET_KW                   = 70
	VertexLexerWILLSET_KW               = 71
	VertexLexerDIDSET_KW                = 72
	VertexLexerASYNC_KW                 = 73
	VertexLexerFINAL_KW                 = 74
	VertexLexerACTOR_KW                 = 75
	VertexLexerPREFIX_KW                = 76
	VertexLexerPOSTFIX_KW               = 77
	VertexLexerMACRO_KW                 = 78
	VertexLexerDYNAMIC_KW               = 79
	VertexLexerLAZY_KW                  = 80
	VertexLexerOPTIONAL_KW              = 81
	VertexLexerOVERRIDE_KW              = 82
	VertexLexerREQUIRED_KW              = 83
	VertexLexerUNOWNED_KW               = 84
	VertexLexerWEAK_KW                  = 85
	VertexLexerNONISOLATED_KW           = 86
	VertexLexerMUTATING_KW              = 87
	VertexLexerNONMUTATING_KW           = 88
	VertexLexerSOME_KW                  = 89
	VertexLexerTYPE_KW                  = 90
	VertexLexerPROTOCOL_KW              = 91
	VertexLexerCONSUMING_KW             = 92
	VertexLexerBORROWING_KW             = 93
	VertexLexerSENDING_KW               = 94
	VertexLexerPOUND_AVAILABLE          = 95
	VertexLexerPOUND_UNAVAILABLE        = 96
	VertexLexerPOUND_IF                 = 97
	VertexLexerPOUND_ELSEIF             = 98
	VertexLexerPOUND_ELSE               = 99
	VertexLexerPOUND_ENDIF              = 100
	VertexLexerPOUND_SOURCE_LOCATION    = 101
	VertexLexerPOUND_FILE               = 102
	VertexLexerPOUND_FILEID             = 103
	VertexLexerPOUND_FILEPATH           = 104
	VertexLexerPOUND_LINE               = 105
	VertexLexerPOUND_COLUMN             = 106
	VertexLexerPOUND_FUNCTION           = 107
	VertexLexerLPAREN                   = 108
	VertexLexerRPAREN                   = 109
	VertexLexerLBRACE                   = 110
	VertexLexerRBRACE                   = 111
	VertexLexerLBRACKET                 = 112
	VertexLexerRBRACKET                 = 113
	VertexLexerDOT                      = 114
	VertexLexerCOMMA                    = 115
	VertexLexerCOLON                    = 116
	VertexLexerSEMICOLON                = 117
	VertexLexerASSIGN                   = 118
	VertexLexerAT                       = 119
	VertexLexerHASH                     = 120
	VertexLexerAMPERSAND                = 121
	VertexLexerARROW                    = 122
	VertexLexerBACKTICK                 = 123
	VertexLexerELLIPSIS                 = 124
	VertexLexerRANGE_HALF_OPEN          = 125
	VertexLexerLT                       = 126
	VertexLexerGT                       = 127
	VertexLexerEXCLAIM_POSTFIX          = 128
	VertexLexerQUESTION_POSTFIX         = 129
	VertexLexerINTEGER_LITERAL          = 130
	VertexLexerFLOAT_LITERAL            = 131
	VertexLexerSTRING_LITERAL           = 132
	VertexLexerMULTILINE_STRING_LITERAL = 133
	VertexLexerEXTENDED_STRING_LITERAL  = 134
	VertexLexerIDENTIFIER               = 135
	VertexLexerOPERATOR                 = 136
)
    VertexLexer tokens.

const (
	VertexLexerBLOCK_COMMENT_CHANNEL = 2
	VertexLexerLINE_COMMENT_CHANNEL  = 3
)
    VertexLexer escapedChannels.

const (
	VertexParserEOF                      = antlr.TokenEOF
	VertexParserWS                       = 1
	VertexParserNEWLINE                  = 2
	VertexParserLINE_COMMENT             = 3
	VertexParserBLOCK_COMMENT            = 4
	VertexParserASSOCIATEDTYPE           = 5
	VertexParserCLASS                    = 6
	VertexParserDEINIT                   = 7
	VertexParserENUM                     = 8
	VertexParserEXTENSION                = 9
	VertexParserFUNC                     = 10
	VertexParserIMPORT                   = 11
	VertexParserINIT                     = 12
	VertexParserINOUT                    = 13
	VertexParserLET                      = 14
	VertexParserOPEN                     = 15
	VertexParserPRIVATE                  = 16
	VertexParserPROTOCOL                 = 17
	VertexParserPUBLIC                   = 18
	VertexParserSTATIC                   = 19
	VertexParserSTRUCT                   = 20
	VertexParserSUBSCRIPT                = 21
	VertexParserTYPEALIAS                = 22
	VertexParserVAR                      = 23
	VertexParserINTERNAL                 = 24
	VertexParserBREAK                    = 25
	VertexParserCASE                     = 26
	VertexParserCATCH                    = 27
	VertexParserCONTINUE                 = 28
	VertexParserDEFAULT                  = 29
	VertexParserDEFER                    = 30
	VertexParserDO                       = 31
	VertexParserELSE                     = 32
	VertexParserFALLTHROUGH              = 33
	VertexParserFOR                      = 34
	VertexParserGUARD                    = 35
	VertexParserIF                       = 36
	VertexParserIN                       = 37
	VertexParserREPEAT                   = 38
	VertexParserRETURN                   = 39
	VertexParserTHROW                    = 40
	VertexParserSWITCH                   = 41
	VertexParserWHERE                    = 42
	VertexParserWHILE                    = 43
	VertexParserANY                      = 44
	VertexParserAS                       = 45
	VertexParserAWAIT                    = 46
	VertexParserFALSE                    = 47
	VertexParserIS                       = 48
	VertexParserNIL                      = 49
	VertexParserSELF                     = 50
	VertexParserSELF_UPPER               = 51
	VertexParserSUPER                    = 52
	VertexParserTHROWS                   = 53
	VertexParserTRUE                     = 54
	VertexParserTRY                      = 55
	VertexParserCONSUME                  = 56
	VertexParserCOPY                     = 57
	VertexParserDISCARD                  = 58
	VertexParserBORROW                   = 59
	VertexParserUNDERSCORE               = 60
	VertexParserOS_KW                    = 61
	VertexParserARCH_KW                  = 62
	VertexParserVERTEX_KW                = 63
	VertexParserCOMPILER_KW              = 64
	VertexParserCANIMPORT_KW             = 65
	VertexParserVERSION_KW               = 66
	VertexParserFILE_KEYWORD             = 67
	VertexParserLINE_KEYWORD             = 68
	VertexParserGET_KW                   = 69
	VertexParserSET_KW                   = 70
	VertexParserWILLSET_KW               = 71
	VertexParserDIDSET_KW                = 72
	VertexParserASYNC_KW                 = 73
	VertexParserFINAL_KW                 = 74
	VertexParserACTOR_KW                 = 75
	VertexParserPREFIX_KW                = 76
	VertexParserPOSTFIX_KW               = 77
	VertexParserMACRO_KW                 = 78
	VertexParserDYNAMIC_KW               = 79
	VertexParserLAZY_KW                  = 80
	VertexParserOPTIONAL_KW              = 81
	VertexParserOVERRIDE_KW              = 82
	VertexParserREQUIRED_KW              = 83
	VertexParserUNOWNED_KW               = 84
	VertexParserWEAK_KW                  = 85
	VertexParserNONISOLATED_KW           = 86
	VertexParserMUTATING_KW              = 87
	VertexParserNONMUTATING_KW           = 88
	VertexParserSOME_KW                  = 89
	VertexParserTYPE_KW                  = 90
	VertexParserPROTOCOL_KW              = 91
	VertexParserCONSUMING_KW             = 92
	VertexParserBORROWING_KW             = 93
	VertexParserSENDING_KW               = 94
	VertexParserPOUND_AVAILABLE          = 95
	VertexParserPOUND_UNAVAILABLE        = 96
	VertexParserPOUND_IF                 = 97
	VertexParserPOUND_ELSEIF             = 98
	VertexParserPOUND_ELSE               = 99
	VertexParserPOUND_ENDIF              = 100
	VertexParserPOUND_SOURCE_LOCATION    = 101
	VertexParserPOUND_FILE               = 102
	VertexParserPOUND_FILEID             = 103
	VertexParserPOUND_FILEPATH           = 104
	VertexParserPOUND_LINE               = 105
	VertexParserPOUND_COLUMN             = 106
	VertexParserPOUND_FUNCTION           = 107
	VertexParserLPAREN                   = 108
	VertexParserRPAREN                   = 109
	VertexParserLBRACE                   = 110
	VertexParserRBRACE                   = 111
	VertexParserLBRACKET                 = 112
	VertexParserRBRACKET                 = 113
	VertexParserDOT                      = 114
	VertexParserCOMMA                    = 115
	VertexParserCOLON                    = 116
	VertexParserSEMICOLON                = 117
	VertexParserASSIGN                   = 118
	VertexParserAT                       = 119
	VertexParserHASH                     = 120
	VertexParserAMPERSAND                = 121
	VertexParserARROW                    = 122
	VertexParserBACKTICK                 = 123
	VertexParserELLIPSIS                 = 124
	VertexParserRANGE_HALF_OPEN          = 125
	VertexParserLT                       = 126
	VertexParserGT                       = 127
	VertexParserEXCLAIM_POSTFIX          = 128
	VertexParserQUESTION_POSTFIX         = 129
	VertexParserINTEGER_LITERAL          = 130
	VertexParserFLOAT_LITERAL            = 131
	VertexParserSTRING_LITERAL           = 132
	VertexParserMULTILINE_STRING_LITERAL = 133
	VertexParserEXTENDED_STRING_LITERAL  = 134
	VertexParserIDENTIFIER               = 135
	VertexParserOPERATOR                 = 136
)
    VertexParser tokens.

const (
	VertexParserRULE_topLevel                          = 0
	VertexParserRULE_statements                        = 1
	VertexParserRULE_statement                         = 2
	VertexParserRULE_loopStatement                     = 3
	VertexParserRULE_forInStatement                    = 4
	VertexParserRULE_whileStatement                    = 5
	VertexParserRULE_conditionList                     = 6
	VertexParserRULE_condition                         = 7
	VertexParserRULE_caseCondition                     = 8
	VertexParserRULE_optionalBindingCondition          = 9
	VertexParserRULE_repeatWhileStatement              = 10
	VertexParserRULE_forStatement                      = 11
	VertexParserRULE_branchStatement                   = 12
	VertexParserRULE_ifStatement                       = 13
	VertexParserRULE_elseClause                        = 14
	VertexParserRULE_guardStatement                    = 15
	VertexParserRULE_switchStatement                   = 16
	VertexParserRULE_switchCases                       = 17
	VertexParserRULE_switchCase                        = 18
	VertexParserRULE_caseLabel                         = 19
	VertexParserRULE_caseItemList                      = 20
	VertexParserRULE_caseItem                          = 21
	VertexParserRULE_defaultLabel                      = 22
	VertexParserRULE_whereClause                       = 23
	VertexParserRULE_labeledStatement                  = 24
	VertexParserRULE_labelName                         = 25
	VertexParserRULE_controlTransfer                   = 26
	VertexParserRULE_deferStatement                    = 27
	VertexParserRULE_doStatement                       = 28
	VertexParserRULE_catchClause                       = 29
	VertexParserRULE_catchPatternList                  = 30
	VertexParserRULE_catchPattern                      = 31
	VertexParserRULE_tryStatement                      = 32
	VertexParserRULE_compilerControl                   = 33
	VertexParserRULE_conditionalCompilationBlock       = 34
	VertexParserRULE_compilationCondition              = 35
	VertexParserRULE_platformCondition                 = 36
	VertexParserRULE_decimalVersion                    = 37
	VertexParserRULE_lineControlStatement              = 38
	VertexParserRULE_diagnosticStatement               = 39
	VertexParserRULE_codeBlock                         = 40
	VertexParserRULE_declaration                       = 41
	VertexParserRULE_importDeclaration                 = 42
	VertexParserRULE_importSpec                        = 43
	VertexParserRULE_importAlias                       = 44
	VertexParserRULE_constantDeclaration               = 45
	VertexParserRULE_patternInitializerList            = 46
	VertexParserRULE_patternInitializer                = 47
	VertexParserRULE_initializer                       = 48
	VertexParserRULE_variableDeclaration               = 49
	VertexParserRULE_variableDeclarationHead           = 50
	VertexParserRULE_variableName                      = 51
	VertexParserRULE_getterSetterBlock                 = 52
	VertexParserRULE_getterClause                      = 53
	VertexParserRULE_setterClause                      = 54
	VertexParserRULE_setterName                        = 55
	VertexParserRULE_getterSetterKeywordBlock          = 56
	VertexParserRULE_getterKeywordClause               = 57
	VertexParserRULE_setterKeywordClause               = 58
	VertexParserRULE_willSetDidSetBlock                = 59
	VertexParserRULE_willSetClause                     = 60
	VertexParserRULE_didSetClause                      = 61
	VertexParserRULE_typealiasDeclaration              = 62
	VertexParserRULE_functionDeclaration               = 63
	VertexParserRULE_functionHead                      = 64
	VertexParserRULE_functionSignature                 = 65
	VertexParserRULE_asyncModifier                     = 66
	VertexParserRULE_functionResult                    = 67
	VertexParserRULE_parameterClause                   = 68
	VertexParserRULE_parameterList                     = 69
	VertexParserRULE_parameter                         = 70
	VertexParserRULE_externalParameterName             = 71
	VertexParserRULE_localParameterName                = 72
	VertexParserRULE_defaultArgumentClause             = 73
	VertexParserRULE_functionBody                      = 74
	VertexParserRULE_enumDeclaration                   = 75
	VertexParserRULE_enumMembers                       = 76
	VertexParserRULE_enumMember                        = 77
	VertexParserRULE_unionStyleEnumCaseClause          = 78
	VertexParserRULE_unionStyleEnumCaseList            = 79
	VertexParserRULE_unionStyleEnumCase                = 80
	VertexParserRULE_rawValueStyleEnumCaseClause       = 81
	VertexParserRULE_rawValueStyleEnumCaseList         = 82
	VertexParserRULE_rawValueStyleEnumCase             = 83
	VertexParserRULE_rawValueAssignment                = 84
	VertexParserRULE_rawValueLiteral                   = 85
	VertexParserRULE_structDeclaration                 = 86
	VertexParserRULE_structBody                        = 87
	VertexParserRULE_structMember                      = 88
	VertexParserRULE_classDeclaration                  = 89
	VertexParserRULE_classBody                         = 90
	VertexParserRULE_classMember                       = 91
	VertexParserRULE_actorDeclaration                  = 92
	VertexParserRULE_actorBody                         = 93
	VertexParserRULE_actorMember                       = 94
	VertexParserRULE_protocolDeclaration               = 95
	VertexParserRULE_primaryAssociatedTypeClause       = 96
	VertexParserRULE_primaryAssociatedTypeList         = 97
	VertexParserRULE_protocolBody                      = 98
	VertexParserRULE_protocolMember                    = 99
	VertexParserRULE_protocolMemberDeclaration         = 100
	VertexParserRULE_protocolPropertyDeclaration       = 101
	VertexParserRULE_protocolMethodDeclaration         = 102
	VertexParserRULE_protocolInitializerDeclaration    = 103
	VertexParserRULE_protocolSubscriptDeclaration      = 104
	VertexParserRULE_protocolAssociatedTypeDeclaration = 105
	VertexParserRULE_initializerDeclaration            = 106
	VertexParserRULE_initializerBody                   = 107
	VertexParserRULE_deinitializerDeclaration          = 108
	VertexParserRULE_extensionDeclaration              = 109
	VertexParserRULE_extensionBody                     = 110
	VertexParserRULE_extensionMember                   = 111
	VertexParserRULE_subscriptDeclaration              = 112
	VertexParserRULE_subscriptHead                     = 113
	VertexParserRULE_subscriptResult                   = 114
	VertexParserRULE_macroDeclaration                  = 115
	VertexParserRULE_declarationModifiers              = 116
	VertexParserRULE_declarationModifier               = 117
	VertexParserRULE_accessLevelModifier               = 118
	VertexParserRULE_mutationModifier                  = 119
	VertexParserRULE_expression                        = 120
	VertexParserRULE_tryOperator                       = 121
	VertexParserRULE_awaitOperator                     = 122
	VertexParserRULE_binaryExpressions                 = 123
	VertexParserRULE_binaryExpression                  = 124
	VertexParserRULE_binaryOperator                    = 125
	VertexParserRULE_assignmentOperator                = 126
	VertexParserRULE_conditionalOperator               = 127
	VertexParserRULE_typeCastingOperator               = 128
	VertexParserRULE_prefixExpression                  = 129
	VertexParserRULE_inOutExpression                   = 130
	VertexParserRULE_prefixOperator                    = 131
	VertexParserRULE_postfixExpression                 = 132
	VertexParserRULE_postfixSuffix                     = 133
	VertexParserRULE_postfixOperator                   = 134
	VertexParserRULE_forcedValueSuffix                 = 135
	VertexParserRULE_optionalChainingLiteral           = 136
	VertexParserRULE_functionCallSuffix                = 137
	VertexParserRULE_functionCallArgumentClause        = 138
	VertexParserRULE_functionCallArgumentList          = 139
	VertexParserRULE_functionCallArgument              = 140
	VertexParserRULE_operator_                         = 141
	VertexParserRULE_trailingClosures                  = 142
	VertexParserRULE_labeledTrailingClosure            = 143
	VertexParserRULE_initializerSuffix                 = 144
	VertexParserRULE_argumentNames                     = 145
	VertexParserRULE_explicitMemberSuffix              = 146
	VertexParserRULE_postfixSelfSuffix                 = 147
	VertexParserRULE_subscriptSuffix                   = 148
	VertexParserRULE_primaryExpression                 = 149
	VertexParserRULE_literalExpression                 = 150
	VertexParserRULE_poundFileExpression               = 151
	VertexParserRULE_literal                           = 152
	VertexParserRULE_numericLiteral                    = 153
	VertexParserRULE_booleanLiteral                    = 154
	VertexParserRULE_arrayLiteral                      = 155
	VertexParserRULE_arrayLiteralItems                 = 156
	VertexParserRULE_dictionaryLiteral                 = 157
	VertexParserRULE_dictionaryLiteralItems            = 158
	VertexParserRULE_dictionaryLiteralItem             = 159
	VertexParserRULE_selfExpression                    = 160
	VertexParserRULE_superExpression                   = 161
	VertexParserRULE_closureExpression                 = 162
	VertexParserRULE_closureSignature                  = 163
	VertexParserRULE_captureList                       = 164
	VertexParserRULE_captureListItems                  = 165
	VertexParserRULE_captureListItem                   = 166
	VertexParserRULE_captureSpecifier                  = 167
	VertexParserRULE_closureParameterClause            = 168
	VertexParserRULE_closureParameterList              = 169
	VertexParserRULE_closureParameter                  = 170
	VertexParserRULE_identifierList                    = 171
	VertexParserRULE_parenthesizedExpression           = 172
	VertexParserRULE_tupleExpression                   = 173
	VertexParserRULE_tupleElementList                  = 174
	VertexParserRULE_tupleElement                      = 175
	VertexParserRULE_implicitMemberExpression          = 176
	VertexParserRULE_wildcardExpression                = 177
	VertexParserRULE_macroExpansionExpression          = 178
	VertexParserRULE_type                              = 179
	VertexParserRULE_typeAnnotationHead                = 180
	VertexParserRULE_functionType                      = 181
	VertexParserRULE_functionTypeArgumentClause        = 182
	VertexParserRULE_functionTypeArgumentList          = 183
	VertexParserRULE_functionTypeArgument              = 184
	VertexParserRULE_argumentLabel                     = 185
	VertexParserRULE_arrayType                         = 186
	VertexParserRULE_dictionaryType                    = 187
	VertexParserRULE_typeIdentifier                    = 188
	VertexParserRULE_tupleType                         = 189
	VertexParserRULE_tupleTypeElementList              = 190
	VertexParserRULE_tupleTypeElement                  = 191
	VertexParserRULE_protocolCompositionType           = 192
	VertexParserRULE_protocolCompositionTypeElement    = 193
	VertexParserRULE_suppressedType                    = 194
	VertexParserRULE_existentialType                   = 195
	VertexParserRULE_opaqueType                        = 196
	VertexParserRULE_selfType                          = 197
	VertexParserRULE_typeAnnotation                    = 198
	VertexParserRULE_typeInheritanceClause             = 199
	VertexParserRULE_typeInheritanceList               = 200
	VertexParserRULE_genericParameterClause            = 201
	VertexParserRULE_genericParameterList              = 202
	VertexParserRULE_genericParameter                  = 203
	VertexParserRULE_genericWhereClause                = 204
	VertexParserRULE_requirementList                   = 205
	VertexParserRULE_requirement                       = 206
	VertexParserRULE_conformanceRequirement            = 207
	VertexParserRULE_sameTypeRequirement               = 208
	VertexParserRULE_layoutConstraintRequirement       = 209
	VertexParserRULE_genericArgumentClause             = 210
	VertexParserRULE_genericArgumentList               = 211
	VertexParserRULE_genericArgument                   = 212
	VertexParserRULE_pattern                           = 213
	VertexParserRULE_wildcardPattern                   = 214
	VertexParserRULE_identifierPattern                 = 215
	VertexParserRULE_valueBindingPattern               = 216
	VertexParserRULE_tuplePattern                      = 217
	VertexParserRULE_tuplePatternElementList           = 218
	VertexParserRULE_tuplePatternElement               = 219
	VertexParserRULE_enumCasePattern                   = 220
	VertexParserRULE_optionalPattern                   = 221
	VertexParserRULE_expressionPattern                 = 222
	VertexParserRULE_attributes                        = 223
	VertexParserRULE_attribute                         = 224
	VertexParserRULE_attributeArguments                = 225
	VertexParserRULE_attributeArgumentList             = 226
	VertexParserRULE_attributeArgument                 = 227
	VertexParserRULE_availabilityCondition             = 228
	VertexParserRULE_availabilityArguments             = 229
	VertexParserRULE_availabilityArgument              = 230
	VertexParserRULE_identifier                        = 231
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

func InitEmptyAccessLevelModifierContext(p *AccessLevelModifierContext)
func InitEmptyActorBodyContext(p *ActorBodyContext)
func InitEmptyActorDeclarationContext(p *ActorDeclarationContext)
func InitEmptyActorMemberContext(p *ActorMemberContext)
func InitEmptyArgumentLabelContext(p *ArgumentLabelContext)
func InitEmptyArgumentNamesContext(p *ArgumentNamesContext)
func InitEmptyArrayLiteralContext(p *ArrayLiteralContext)
func InitEmptyArrayLiteralItemsContext(p *ArrayLiteralItemsContext)
func InitEmptyArrayTypeContext(p *ArrayTypeContext)
func InitEmptyAssignmentOperatorContext(p *AssignmentOperatorContext)
func InitEmptyAsyncModifierContext(p *AsyncModifierContext)
func InitEmptyAttributeArgumentContext(p *AttributeArgumentContext)
func InitEmptyAttributeArgumentListContext(p *AttributeArgumentListContext)
func InitEmptyAttributeArgumentsContext(p *AttributeArgumentsContext)
func InitEmptyAttributeContext(p *AttributeContext)
func InitEmptyAttributesContext(p *AttributesContext)
func InitEmptyAvailabilityArgumentContext(p *AvailabilityArgumentContext)
func InitEmptyAvailabilityArgumentsContext(p *AvailabilityArgumentsContext)
func InitEmptyAvailabilityConditionContext(p *AvailabilityConditionContext)
func InitEmptyAwaitOperatorContext(p *AwaitOperatorContext)
func InitEmptyBinaryExpressionContext(p *BinaryExpressionContext)
func InitEmptyBinaryExpressionsContext(p *BinaryExpressionsContext)
func InitEmptyBinaryOperatorContext(p *BinaryOperatorContext)
func InitEmptyBooleanLiteralContext(p *BooleanLiteralContext)
func InitEmptyBranchStatementContext(p *BranchStatementContext)
func InitEmptyCaptureListContext(p *CaptureListContext)
func InitEmptyCaptureListItemContext(p *CaptureListItemContext)
func InitEmptyCaptureListItemsContext(p *CaptureListItemsContext)
func InitEmptyCaptureSpecifierContext(p *CaptureSpecifierContext)
func InitEmptyCaseConditionContext(p *CaseConditionContext)
func InitEmptyCaseItemContext(p *CaseItemContext)
func InitEmptyCaseItemListContext(p *CaseItemListContext)
func InitEmptyCaseLabelContext(p *CaseLabelContext)
func InitEmptyCatchClauseContext(p *CatchClauseContext)
func InitEmptyCatchPatternContext(p *CatchPatternContext)
func InitEmptyCatchPatternListContext(p *CatchPatternListContext)
func InitEmptyClassBodyContext(p *ClassBodyContext)
func InitEmptyClassDeclarationContext(p *ClassDeclarationContext)
func InitEmptyClassMemberContext(p *ClassMemberContext)
func InitEmptyClosureExpressionContext(p *ClosureExpressionContext)
func InitEmptyClosureParameterClauseContext(p *ClosureParameterClauseContext)
func InitEmptyClosureParameterContext(p *ClosureParameterContext)
func InitEmptyClosureParameterListContext(p *ClosureParameterListContext)
func InitEmptyClosureSignatureContext(p *ClosureSignatureContext)
func InitEmptyCodeBlockContext(p *CodeBlockContext)
func InitEmptyCompilationConditionContext(p *CompilationConditionContext)
func InitEmptyCompilerControlContext(p *CompilerControlContext)
func InitEmptyConditionContext(p *ConditionContext)
func InitEmptyConditionListContext(p *ConditionListContext)
func InitEmptyConditionalCompilationBlockContext(p *ConditionalCompilationBlockContext)
func InitEmptyConditionalOperatorContext(p *ConditionalOperatorContext)
func InitEmptyConformanceRequirementContext(p *ConformanceRequirementContext)
func InitEmptyConstantDeclarationContext(p *ConstantDeclarationContext)
func InitEmptyControlTransferContext(p *ControlTransferContext)
func InitEmptyDecimalVersionContext(p *DecimalVersionContext)
func InitEmptyDeclarationContext(p *DeclarationContext)
func InitEmptyDeclarationModifierContext(p *DeclarationModifierContext)
func InitEmptyDeclarationModifiersContext(p *DeclarationModifiersContext)
func InitEmptyDefaultArgumentClauseContext(p *DefaultArgumentClauseContext)
func InitEmptyDefaultLabelContext(p *DefaultLabelContext)
func InitEmptyDeferStatementContext(p *DeferStatementContext)
func InitEmptyDeinitializerDeclarationContext(p *DeinitializerDeclarationContext)
func InitEmptyDiagnosticStatementContext(p *DiagnosticStatementContext)
func InitEmptyDictionaryLiteralContext(p *DictionaryLiteralContext)
func InitEmptyDictionaryLiteralItemContext(p *DictionaryLiteralItemContext)
func InitEmptyDictionaryLiteralItemsContext(p *DictionaryLiteralItemsContext)
func InitEmptyDictionaryTypeContext(p *DictionaryTypeContext)
func InitEmptyDidSetClauseContext(p *DidSetClauseContext)
func InitEmptyDoStatementContext(p *DoStatementContext)
func InitEmptyElseClauseContext(p *ElseClauseContext)
func InitEmptyEnumCasePatternContext(p *EnumCasePatternContext)
func InitEmptyEnumDeclarationContext(p *EnumDeclarationContext)
func InitEmptyEnumMemberContext(p *EnumMemberContext)
func InitEmptyEnumMembersContext(p *EnumMembersContext)
func InitEmptyExistentialTypeContext(p *ExistentialTypeContext)
func InitEmptyExplicitMemberSuffixContext(p *ExplicitMemberSuffixContext)
func InitEmptyExpressionContext(p *ExpressionContext)
func InitEmptyExpressionPatternContext(p *ExpressionPatternContext)
func InitEmptyExtensionBodyContext(p *ExtensionBodyContext)
func InitEmptyExtensionDeclarationContext(p *ExtensionDeclarationContext)
func InitEmptyExtensionMemberContext(p *ExtensionMemberContext)
func InitEmptyExternalParameterNameContext(p *ExternalParameterNameContext)
func InitEmptyForInStatementContext(p *ForInStatementContext)
func InitEmptyForStatementContext(p *ForStatementContext)
func InitEmptyForcedValueSuffixContext(p *ForcedValueSuffixContext)
func InitEmptyFunctionBodyContext(p *FunctionBodyContext)
func InitEmptyFunctionCallArgumentClauseContext(p *FunctionCallArgumentClauseContext)
func InitEmptyFunctionCallArgumentContext(p *FunctionCallArgumentContext)
func InitEmptyFunctionCallArgumentListContext(p *FunctionCallArgumentListContext)
func InitEmptyFunctionCallSuffixContext(p *FunctionCallSuffixContext)
func InitEmptyFunctionDeclarationContext(p *FunctionDeclarationContext)
func InitEmptyFunctionHeadContext(p *FunctionHeadContext)
func InitEmptyFunctionResultContext(p *FunctionResultContext)
func InitEmptyFunctionSignatureContext(p *FunctionSignatureContext)
func InitEmptyFunctionTypeArgumentClauseContext(p *FunctionTypeArgumentClauseContext)
func InitEmptyFunctionTypeArgumentContext(p *FunctionTypeArgumentContext)
func InitEmptyFunctionTypeArgumentListContext(p *FunctionTypeArgumentListContext)
func InitEmptyFunctionTypeContext(p *FunctionTypeContext)
func InitEmptyGenericArgumentClauseContext(p *GenericArgumentClauseContext)
func InitEmptyGenericArgumentContext(p *GenericArgumentContext)
func InitEmptyGenericArgumentListContext(p *GenericArgumentListContext)
func InitEmptyGenericParameterClauseContext(p *GenericParameterClauseContext)
func InitEmptyGenericParameterContext(p *GenericParameterContext)
func InitEmptyGenericParameterListContext(p *GenericParameterListContext)
func InitEmptyGenericWhereClauseContext(p *GenericWhereClauseContext)
func InitEmptyGetterClauseContext(p *GetterClauseContext)
func InitEmptyGetterKeywordClauseContext(p *GetterKeywordClauseContext)
func InitEmptyGetterSetterBlockContext(p *GetterSetterBlockContext)
func InitEmptyGetterSetterKeywordBlockContext(p *GetterSetterKeywordBlockContext)
func InitEmptyGuardStatementContext(p *GuardStatementContext)
func InitEmptyIdentifierContext(p *IdentifierContext)
func InitEmptyIdentifierListContext(p *IdentifierListContext)
func InitEmptyIdentifierPatternContext(p *IdentifierPatternContext)
func InitEmptyIfStatementContext(p *IfStatementContext)
func InitEmptyImplicitMemberExpressionContext(p *ImplicitMemberExpressionContext)
func InitEmptyImportAliasContext(p *ImportAliasContext)
func InitEmptyImportDeclarationContext(p *ImportDeclarationContext)
func InitEmptyImportSpecContext(p *ImportSpecContext)
func InitEmptyInOutExpressionContext(p *InOutExpressionContext)
func InitEmptyInitializerBodyContext(p *InitializerBodyContext)
func InitEmptyInitializerContext(p *InitializerContext)
func InitEmptyInitializerDeclarationContext(p *InitializerDeclarationContext)
func InitEmptyInitializerSuffixContext(p *InitializerSuffixContext)
func InitEmptyLabelNameContext(p *LabelNameContext)
func InitEmptyLabeledStatementContext(p *LabeledStatementContext)
func InitEmptyLabeledTrailingClosureContext(p *LabeledTrailingClosureContext)
func InitEmptyLayoutConstraintRequirementContext(p *LayoutConstraintRequirementContext)
func InitEmptyLineControlStatementContext(p *LineControlStatementContext)
func InitEmptyLiteralContext(p *LiteralContext)
func InitEmptyLiteralExpressionContext(p *LiteralExpressionContext)
func InitEmptyLocalParameterNameContext(p *LocalParameterNameContext)
func InitEmptyLoopStatementContext(p *LoopStatementContext)
func InitEmptyMacroDeclarationContext(p *MacroDeclarationContext)
func InitEmptyMacroExpansionExpressionContext(p *MacroExpansionExpressionContext)
func InitEmptyMutationModifierContext(p *MutationModifierContext)
func InitEmptyNumericLiteralContext(p *NumericLiteralContext)
func InitEmptyOpaqueTypeContext(p *OpaqueTypeContext)
func InitEmptyOperator_Context(p *Operator_Context)
func InitEmptyOptionalBindingConditionContext(p *OptionalBindingConditionContext)
func InitEmptyOptionalChainingLiteralContext(p *OptionalChainingLiteralContext)
func InitEmptyOptionalPatternContext(p *OptionalPatternContext)
func InitEmptyParameterClauseContext(p *ParameterClauseContext)
func InitEmptyParameterContext(p *ParameterContext)
func InitEmptyParameterListContext(p *ParameterListContext)
func InitEmptyParenthesizedExpressionContext(p *ParenthesizedExpressionContext)
func InitEmptyPatternContext(p *PatternContext)
func InitEmptyPatternInitializerContext(p *PatternInitializerContext)
func InitEmptyPatternInitializerListContext(p *PatternInitializerListContext)
func InitEmptyPlatformConditionContext(p *PlatformConditionContext)
func InitEmptyPostfixExpressionContext(p *PostfixExpressionContext)
func InitEmptyPostfixOperatorContext(p *PostfixOperatorContext)
func InitEmptyPostfixSelfSuffixContext(p *PostfixSelfSuffixContext)
func InitEmptyPostfixSuffixContext(p *PostfixSuffixContext)
func InitEmptyPoundFileExpressionContext(p *PoundFileExpressionContext)
func InitEmptyPrefixExpressionContext(p *PrefixExpressionContext)
func InitEmptyPrefixOperatorContext(p *PrefixOperatorContext)
func InitEmptyPrimaryAssociatedTypeClauseContext(p *PrimaryAssociatedTypeClauseContext)
func InitEmptyPrimaryAssociatedTypeListContext(p *PrimaryAssociatedTypeListContext)
func InitEmptyPrimaryExpressionContext(p *PrimaryExpressionContext)
func InitEmptyProtocolAssociatedTypeDeclarationContext(p *ProtocolAssociatedTypeDeclarationContext)
func InitEmptyProtocolBodyContext(p *ProtocolBodyContext)
func InitEmptyProtocolCompositionTypeContext(p *ProtocolCompositionTypeContext)
func InitEmptyProtocolCompositionTypeElementContext(p *ProtocolCompositionTypeElementContext)
func InitEmptyProtocolDeclarationContext(p *ProtocolDeclarationContext)
func InitEmptyProtocolInitializerDeclarationContext(p *ProtocolInitializerDeclarationContext)
func InitEmptyProtocolMemberContext(p *ProtocolMemberContext)
func InitEmptyProtocolMemberDeclarationContext(p *ProtocolMemberDeclarationContext)
func InitEmptyProtocolMethodDeclarationContext(p *ProtocolMethodDeclarationContext)
func InitEmptyProtocolPropertyDeclarationContext(p *ProtocolPropertyDeclarationContext)
func InitEmptyProtocolSubscriptDeclarationContext(p *ProtocolSubscriptDeclarationContext)
func InitEmptyRawValueAssignmentContext(p *RawValueAssignmentContext)
func InitEmptyRawValueLiteralContext(p *RawValueLiteralContext)
func InitEmptyRawValueStyleEnumCaseClauseContext(p *RawValueStyleEnumCaseClauseContext)
func InitEmptyRawValueStyleEnumCaseContext(p *RawValueStyleEnumCaseContext)
func InitEmptyRawValueStyleEnumCaseListContext(p *RawValueStyleEnumCaseListContext)
func InitEmptyRepeatWhileStatementContext(p *RepeatWhileStatementContext)
func InitEmptyRequirementContext(p *RequirementContext)
func InitEmptyRequirementListContext(p *RequirementListContext)
func InitEmptySameTypeRequirementContext(p *SameTypeRequirementContext)
func InitEmptySelfExpressionContext(p *SelfExpressionContext)
func InitEmptySelfTypeContext(p *SelfTypeContext)
func InitEmptySetterClauseContext(p *SetterClauseContext)
func InitEmptySetterKeywordClauseContext(p *SetterKeywordClauseContext)
func InitEmptySetterNameContext(p *SetterNameContext)
func InitEmptyStatementContext(p *StatementContext)
func InitEmptyStatementsContext(p *StatementsContext)
func InitEmptyStructBodyContext(p *StructBodyContext)
func InitEmptyStructDeclarationContext(p *StructDeclarationContext)
func InitEmptyStructMemberContext(p *StructMemberContext)
func InitEmptySubscriptDeclarationContext(p *SubscriptDeclarationContext)
func InitEmptySubscriptHeadContext(p *SubscriptHeadContext)
func InitEmptySubscriptResultContext(p *SubscriptResultContext)
func InitEmptySubscriptSuffixContext(p *SubscriptSuffixContext)
func InitEmptySuperExpressionContext(p *SuperExpressionContext)
func InitEmptySuppressedTypeContext(p *SuppressedTypeContext)
func InitEmptySwitchCaseContext(p *SwitchCaseContext)
func InitEmptySwitchCasesContext(p *SwitchCasesContext)
func InitEmptySwitchStatementContext(p *SwitchStatementContext)
func InitEmptyTopLevelContext(p *TopLevelContext)
func InitEmptyTrailingClosuresContext(p *TrailingClosuresContext)
func InitEmptyTryOperatorContext(p *TryOperatorContext)
func InitEmptyTryStatementContext(p *TryStatementContext)
func InitEmptyTupleElementContext(p *TupleElementContext)
func InitEmptyTupleElementListContext(p *TupleElementListContext)
func InitEmptyTupleExpressionContext(p *TupleExpressionContext)
func InitEmptyTuplePatternContext(p *TuplePatternContext)
func InitEmptyTuplePatternElementContext(p *TuplePatternElementContext)
func InitEmptyTuplePatternElementListContext(p *TuplePatternElementListContext)
func InitEmptyTupleTypeContext(p *TupleTypeContext)
func InitEmptyTupleTypeElementContext(p *TupleTypeElementContext)
func InitEmptyTupleTypeElementListContext(p *TupleTypeElementListContext)
func InitEmptyTypeAnnotationContext(p *TypeAnnotationContext)
func InitEmptyTypeAnnotationHeadContext(p *TypeAnnotationHeadContext)
func InitEmptyTypeCastingOperatorContext(p *TypeCastingOperatorContext)
func InitEmptyTypeContext(p *TypeContext)
func InitEmptyTypeIdentifierContext(p *TypeIdentifierContext)
func InitEmptyTypeInheritanceClauseContext(p *TypeInheritanceClauseContext)
func InitEmptyTypeInheritanceListContext(p *TypeInheritanceListContext)
func InitEmptyTypealiasDeclarationContext(p *TypealiasDeclarationContext)
func InitEmptyUnionStyleEnumCaseClauseContext(p *UnionStyleEnumCaseClauseContext)
func InitEmptyUnionStyleEnumCaseContext(p *UnionStyleEnumCaseContext)
func InitEmptyUnionStyleEnumCaseListContext(p *UnionStyleEnumCaseListContext)
func InitEmptyValueBindingPatternContext(p *ValueBindingPatternContext)
func InitEmptyVariableDeclarationContext(p *VariableDeclarationContext)
func InitEmptyVariableDeclarationHeadContext(p *VariableDeclarationHeadContext)
func InitEmptyVariableNameContext(p *VariableNameContext)
func InitEmptyWhereClauseContext(p *WhereClauseContext)
func InitEmptyWhileStatementContext(p *WhileStatementContext)
func InitEmptyWildcardExpressionContext(p *WildcardExpressionContext)
func InitEmptyWildcardPatternContext(p *WildcardPatternContext)
func InitEmptyWillSetClauseContext(p *WillSetClauseContext)
func InitEmptyWillSetDidSetBlockContext(p *WillSetDidSetBlockContext)
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

type AccessLevelModifierContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAccessLevelModifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AccessLevelModifierContext

func NewEmptyAccessLevelModifierContext() *AccessLevelModifierContext

func (s *AccessLevelModifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AccessLevelModifierContext) GetParser() antlr.Parser

func (s *AccessLevelModifierContext) GetRuleContext() antlr.RuleContext

func (s *AccessLevelModifierContext) INTERNAL() antlr.TerminalNode

func (*AccessLevelModifierContext) IsAccessLevelModifierContext()

func (s *AccessLevelModifierContext) LPAREN() antlr.TerminalNode

func (s *AccessLevelModifierContext) OPEN() antlr.TerminalNode

func (s *AccessLevelModifierContext) PRIVATE() antlr.TerminalNode

func (s *AccessLevelModifierContext) PUBLIC() antlr.TerminalNode

func (s *AccessLevelModifierContext) RPAREN() antlr.TerminalNode

func (s *AccessLevelModifierContext) SET_KW() antlr.TerminalNode

func (s *AccessLevelModifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ActorBodyContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewActorBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ActorBodyContext

func NewEmptyActorBodyContext() *ActorBodyContext

func (s *ActorBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ActorBodyContext) ActorMember(i int) IActorMemberContext

func (s *ActorBodyContext) AllActorMember() []IActorMemberContext

func (s *ActorBodyContext) GetParser() antlr.Parser

func (s *ActorBodyContext) GetRuleContext() antlr.RuleContext

func (*ActorBodyContext) IsActorBodyContext()

func (s *ActorBodyContext) LBRACE() antlr.TerminalNode

func (s *ActorBodyContext) RBRACE() antlr.TerminalNode

func (s *ActorBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ActorDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewActorDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ActorDeclarationContext

func NewEmptyActorDeclarationContext() *ActorDeclarationContext

func (s *ActorDeclarationContext) ACTOR_KW() antlr.TerminalNode

func (s *ActorDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ActorDeclarationContext) AccessLevelModifier() IAccessLevelModifierContext

func (s *ActorDeclarationContext) ActorBody() IActorBodyContext

func (s *ActorDeclarationContext) Attributes() IAttributesContext

func (s *ActorDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *ActorDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *ActorDeclarationContext) GetParser() antlr.Parser

func (s *ActorDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *ActorDeclarationContext) Identifier() IIdentifierContext

func (*ActorDeclarationContext) IsActorDeclarationContext()

func (s *ActorDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ActorDeclarationContext) TypeInheritanceClause() ITypeInheritanceClauseContext

type ActorMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewActorMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ActorMemberContext

func NewEmptyActorMemberContext() *ActorMemberContext

func (s *ActorMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ActorMemberContext) CompilerControl() ICompilerControlContext

func (s *ActorMemberContext) Declaration() IDeclarationContext

func (s *ActorMemberContext) GetParser() antlr.Parser

func (s *ActorMemberContext) GetRuleContext() antlr.RuleContext

func (*ActorMemberContext) IsActorMemberContext()

func (s *ActorMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AnnotatedTypeContext struct {
	TypeContext
}

func NewAnnotatedTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AnnotatedTypeContext

func (s *AnnotatedTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AnnotatedTypeContext) GetRuleContext() antlr.RuleContext

func (s *AnnotatedTypeContext) TypeAnnotationHead() ITypeAnnotationHeadContext

func (s *AnnotatedTypeContext) Type_() ITypeContext

type ArgumentLabelContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArgumentLabelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentLabelContext

func NewEmptyArgumentLabelContext() *ArgumentLabelContext

func (s *ArgumentLabelContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArgumentLabelContext) COLON() antlr.TerminalNode

func (s *ArgumentLabelContext) GetParser() antlr.Parser

func (s *ArgumentLabelContext) GetRuleContext() antlr.RuleContext

func (s *ArgumentLabelContext) Identifier() IIdentifierContext

func (*ArgumentLabelContext) IsArgumentLabelContext()

func (s *ArgumentLabelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ArgumentLabelContext) UNDERSCORE() antlr.TerminalNode

type ArgumentNamesContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArgumentNamesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentNamesContext

func NewEmptyArgumentNamesContext() *ArgumentNamesContext

func (s *ArgumentNamesContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArgumentNamesContext) AllCOLON() []antlr.TerminalNode

func (s *ArgumentNamesContext) AllIdentifier() []IIdentifierContext

func (s *ArgumentNamesContext) COLON(i int) antlr.TerminalNode

func (s *ArgumentNamesContext) GetParser() antlr.Parser

func (s *ArgumentNamesContext) GetRuleContext() antlr.RuleContext

func (s *ArgumentNamesContext) Identifier(i int) IIdentifierContext

func (*ArgumentNamesContext) IsArgumentNamesContext()

func (s *ArgumentNamesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ArrTypeContext struct {
	TypeContext
}

func NewArrTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ArrTypeContext

func (s *ArrTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArrTypeContext) ArrayType() IArrayTypeContext

func (s *ArrTypeContext) GetRuleContext() antlr.RuleContext

type ArrayLiteralContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArrayLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayLiteralContext

func NewEmptyArrayLiteralContext() *ArrayLiteralContext

func (s *ArrayLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArrayLiteralContext) ArrayLiteralItems() IArrayLiteralItemsContext

func (s *ArrayLiteralContext) GetParser() antlr.Parser

func (s *ArrayLiteralContext) GetRuleContext() antlr.RuleContext

func (*ArrayLiteralContext) IsArrayLiteralContext()

func (s *ArrayLiteralContext) LBRACKET() antlr.TerminalNode

func (s *ArrayLiteralContext) RBRACKET() antlr.TerminalNode

func (s *ArrayLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ArrayLiteralItemsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArrayLiteralItemsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayLiteralItemsContext

func NewEmptyArrayLiteralItemsContext() *ArrayLiteralItemsContext

func (s *ArrayLiteralItemsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArrayLiteralItemsContext) AllCOMMA() []antlr.TerminalNode

func (s *ArrayLiteralItemsContext) AllExpression() []IExpressionContext

func (s *ArrayLiteralItemsContext) COMMA(i int) antlr.TerminalNode

func (s *ArrayLiteralItemsContext) Expression(i int) IExpressionContext

func (s *ArrayLiteralItemsContext) GetParser() antlr.Parser

func (s *ArrayLiteralItemsContext) GetRuleContext() antlr.RuleContext

func (*ArrayLiteralItemsContext) IsArrayLiteralItemsContext()

func (s *ArrayLiteralItemsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ArrayTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArrayTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayTypeContext

func NewEmptyArrayTypeContext() *ArrayTypeContext

func (s *ArrayTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArrayTypeContext) GetParser() antlr.Parser

func (s *ArrayTypeContext) GetRuleContext() antlr.RuleContext

func (*ArrayTypeContext) IsArrayTypeContext()

func (s *ArrayTypeContext) LBRACKET() antlr.TerminalNode

func (s *ArrayTypeContext) RBRACKET() antlr.TerminalNode

func (s *ArrayTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ArrayTypeContext) Type_() ITypeContext

type AsPatContext struct {
	PatternContext
}

func NewAsPatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AsPatContext

func (s *AsPatContext) AS() antlr.TerminalNode

func (s *AsPatContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AsPatContext) GetRuleContext() antlr.RuleContext

func (s *AsPatContext) Pattern() IPatternContext

func (s *AsPatContext) Type_() ITypeContext

type AssignmentExprContext struct {
	BinaryExpressionContext
}

func NewAssignmentExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AssignmentExprContext

func (s *AssignmentExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AssignmentExprContext) AssignmentOperator() IAssignmentOperatorContext

func (s *AssignmentExprContext) Expression() IExpressionContext

func (s *AssignmentExprContext) GetRuleContext() antlr.RuleContext

func (s *AssignmentExprContext) TryOperator() ITryOperatorContext

type AssignmentOperatorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAssignmentOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentOperatorContext

func NewEmptyAssignmentOperatorContext() *AssignmentOperatorContext

func (s *AssignmentOperatorContext) ASSIGN() antlr.TerminalNode

func (s *AssignmentOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AssignmentOperatorContext) GetParser() antlr.Parser

func (s *AssignmentOperatorContext) GetRuleContext() antlr.RuleContext

func (*AssignmentOperatorContext) IsAssignmentOperatorContext()

func (s *AssignmentOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AsyncModifierContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAsyncModifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AsyncModifierContext

func NewEmptyAsyncModifierContext() *AsyncModifierContext

func (s *AsyncModifierContext) ASYNC_KW() antlr.TerminalNode

func (s *AsyncModifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AsyncModifierContext) GetParser() antlr.Parser

func (s *AsyncModifierContext) GetRuleContext() antlr.RuleContext

func (*AsyncModifierContext) IsAsyncModifierContext()

func (s *AsyncModifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AttributeArgumentContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAttributeArgumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributeArgumentContext

func NewEmptyAttributeArgumentContext() *AttributeArgumentContext

func (s *AttributeArgumentContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AttributeArgumentContext) COLON() antlr.TerminalNode

func (s *AttributeArgumentContext) Expression() IExpressionContext

func (s *AttributeArgumentContext) GetParser() antlr.Parser

func (s *AttributeArgumentContext) GetRuleContext() antlr.RuleContext

func (s *AttributeArgumentContext) Identifier() IIdentifierContext

func (*AttributeArgumentContext) IsAttributeArgumentContext()

func (s *AttributeArgumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *AttributeArgumentContext) Type_() ITypeContext

type AttributeArgumentListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAttributeArgumentListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributeArgumentListContext

func NewEmptyAttributeArgumentListContext() *AttributeArgumentListContext

func (s *AttributeArgumentListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AttributeArgumentListContext) AllAttributeArgument() []IAttributeArgumentContext

func (s *AttributeArgumentListContext) AllCOMMA() []antlr.TerminalNode

func (s *AttributeArgumentListContext) AttributeArgument(i int) IAttributeArgumentContext

func (s *AttributeArgumentListContext) COMMA(i int) antlr.TerminalNode

func (s *AttributeArgumentListContext) GetParser() antlr.Parser

func (s *AttributeArgumentListContext) GetRuleContext() antlr.RuleContext

func (*AttributeArgumentListContext) IsAttributeArgumentListContext()

func (s *AttributeArgumentListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AttributeArgumentsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAttributeArgumentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributeArgumentsContext

func NewEmptyAttributeArgumentsContext() *AttributeArgumentsContext

func (s *AttributeArgumentsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AttributeArgumentsContext) AttributeArgumentList() IAttributeArgumentListContext

func (s *AttributeArgumentsContext) GetParser() antlr.Parser

func (s *AttributeArgumentsContext) GetRuleContext() antlr.RuleContext

func (*AttributeArgumentsContext) IsAttributeArgumentsContext()

func (s *AttributeArgumentsContext) LPAREN() antlr.TerminalNode

func (s *AttributeArgumentsContext) RPAREN() antlr.TerminalNode

func (s *AttributeArgumentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AttributeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAttributeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributeContext

func NewEmptyAttributeContext() *AttributeContext

func (s *AttributeContext) AT() antlr.TerminalNode

func (s *AttributeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AttributeContext) AttributeArguments() IAttributeArgumentsContext

func (s *AttributeContext) GetParser() antlr.Parser

func (s *AttributeContext) GetRuleContext() antlr.RuleContext

func (s *AttributeContext) Identifier() IIdentifierContext

func (*AttributeContext) IsAttributeContext()

func (s *AttributeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AttributesContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAttributesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributesContext

func NewEmptyAttributesContext() *AttributesContext

func (s *AttributesContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AttributesContext) AllAttribute() []IAttributeContext

func (s *AttributesContext) Attribute(i int) IAttributeContext

func (s *AttributesContext) GetParser() antlr.Parser

func (s *AttributesContext) GetRuleContext() antlr.RuleContext

func (*AttributesContext) IsAttributesContext()

func (s *AttributesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AvailabilityArgumentContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAvailabilityArgumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AvailabilityArgumentContext

func NewEmptyAvailabilityArgumentContext() *AvailabilityArgumentContext

func (s *AvailabilityArgumentContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AvailabilityArgumentContext) DecimalVersion() IDecimalVersionContext

func (s *AvailabilityArgumentContext) GetParser() antlr.Parser

func (s *AvailabilityArgumentContext) GetRuleContext() antlr.RuleContext

func (s *AvailabilityArgumentContext) Identifier() IIdentifierContext

func (*AvailabilityArgumentContext) IsAvailabilityArgumentContext()

func (s *AvailabilityArgumentContext) OPERATOR() antlr.TerminalNode

func (s *AvailabilityArgumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AvailabilityArgumentsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAvailabilityArgumentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AvailabilityArgumentsContext

func NewEmptyAvailabilityArgumentsContext() *AvailabilityArgumentsContext

func (s *AvailabilityArgumentsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AvailabilityArgumentsContext) AllAvailabilityArgument() []IAvailabilityArgumentContext

func (s *AvailabilityArgumentsContext) AllCOMMA() []antlr.TerminalNode

func (s *AvailabilityArgumentsContext) AvailabilityArgument(i int) IAvailabilityArgumentContext

func (s *AvailabilityArgumentsContext) COMMA(i int) antlr.TerminalNode

func (s *AvailabilityArgumentsContext) GetParser() antlr.Parser

func (s *AvailabilityArgumentsContext) GetRuleContext() antlr.RuleContext

func (*AvailabilityArgumentsContext) IsAvailabilityArgumentsContext()

func (s *AvailabilityArgumentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AvailabilityConditionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAvailabilityConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AvailabilityConditionContext

func NewEmptyAvailabilityConditionContext() *AvailabilityConditionContext

func (s *AvailabilityConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AvailabilityConditionContext) AvailabilityArguments() IAvailabilityArgumentsContext

func (s *AvailabilityConditionContext) GetParser() antlr.Parser

func (s *AvailabilityConditionContext) GetRuleContext() antlr.RuleContext

func (*AvailabilityConditionContext) IsAvailabilityConditionContext()

func (s *AvailabilityConditionContext) LPAREN() antlr.TerminalNode

func (s *AvailabilityConditionContext) POUND_AVAILABLE() antlr.TerminalNode

func (s *AvailabilityConditionContext) POUND_UNAVAILABLE() antlr.TerminalNode

func (s *AvailabilityConditionContext) RPAREN() antlr.TerminalNode

func (s *AvailabilityConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AwaitOperatorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAwaitOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AwaitOperatorContext

func NewEmptyAwaitOperatorContext() *AwaitOperatorContext

func (s *AwaitOperatorContext) AWAIT() antlr.TerminalNode

func (s *AwaitOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AwaitOperatorContext) GetParser() antlr.Parser

func (s *AwaitOperatorContext) GetRuleContext() antlr.RuleContext

func (*AwaitOperatorContext) IsAwaitOperatorContext()

func (s *AwaitOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BarePostfixContext struct {
	PrefixExpressionContext
}

func NewBarePostfixContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BarePostfixContext

func (s *BarePostfixContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BarePostfixContext) GetRuleContext() antlr.RuleContext

func (s *BarePostfixContext) PostfixExpression() IPostfixExpressionContext

type BaseVertexParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseVertexParserVisitor) VisitAccessLevelModifier(ctx *AccessLevelModifierContext) interface{}

func (v *BaseVertexParserVisitor) VisitActorBody(ctx *ActorBodyContext) interface{}

func (v *BaseVertexParserVisitor) VisitActorDeclaration(ctx *ActorDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitActorMember(ctx *ActorMemberContext) interface{}

func (v *BaseVertexParserVisitor) VisitAnnotatedType(ctx *AnnotatedTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitArgumentLabel(ctx *ArgumentLabelContext) interface{}

func (v *BaseVertexParserVisitor) VisitArgumentNames(ctx *ArgumentNamesContext) interface{}

func (v *BaseVertexParserVisitor) VisitArrType(ctx *ArrTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitArrayLiteral(ctx *ArrayLiteralContext) interface{}

func (v *BaseVertexParserVisitor) VisitArrayLiteralItems(ctx *ArrayLiteralItemsContext) interface{}

func (v *BaseVertexParserVisitor) VisitArrayType(ctx *ArrayTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitAsPat(ctx *AsPatContext) interface{}

func (v *BaseVertexParserVisitor) VisitAssignmentExpr(ctx *AssignmentExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitAssignmentOperator(ctx *AssignmentOperatorContext) interface{}

func (v *BaseVertexParserVisitor) VisitAsyncModifier(ctx *AsyncModifierContext) interface{}

func (v *BaseVertexParserVisitor) VisitAttribute(ctx *AttributeContext) interface{}

func (v *BaseVertexParserVisitor) VisitAttributeArgument(ctx *AttributeArgumentContext) interface{}

func (v *BaseVertexParserVisitor) VisitAttributeArgumentList(ctx *AttributeArgumentListContext) interface{}

func (v *BaseVertexParserVisitor) VisitAttributeArguments(ctx *AttributeArgumentsContext) interface{}

func (v *BaseVertexParserVisitor) VisitAttributes(ctx *AttributesContext) interface{}

func (v *BaseVertexParserVisitor) VisitAvailabilityArgument(ctx *AvailabilityArgumentContext) interface{}

func (v *BaseVertexParserVisitor) VisitAvailabilityArguments(ctx *AvailabilityArgumentsContext) interface{}

func (v *BaseVertexParserVisitor) VisitAvailabilityCondition(ctx *AvailabilityConditionContext) interface{}

func (v *BaseVertexParserVisitor) VisitAwaitOperator(ctx *AwaitOperatorContext) interface{}

func (v *BaseVertexParserVisitor) VisitBarePostfix(ctx *BarePostfixContext) interface{}

func (v *BaseVertexParserVisitor) VisitBinaryExpressions(ctx *BinaryExpressionsContext) interface{}

func (v *BaseVertexParserVisitor) VisitBinaryOp(ctx *BinaryOpContext) interface{}

func (v *BaseVertexParserVisitor) VisitBinaryOperator(ctx *BinaryOperatorContext) interface{}

func (v *BaseVertexParserVisitor) VisitBindingPat(ctx *BindingPatContext) interface{}

func (v *BaseVertexParserVisitor) VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{}

func (v *BaseVertexParserVisitor) VisitBranchStatement(ctx *BranchStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitBranchStmt(ctx *BranchStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitBreakStatement(ctx *BreakStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitCaptureList(ctx *CaptureListContext) interface{}

func (v *BaseVertexParserVisitor) VisitCaptureListItem(ctx *CaptureListItemContext) interface{}

func (v *BaseVertexParserVisitor) VisitCaptureListItems(ctx *CaptureListItemsContext) interface{}

func (v *BaseVertexParserVisitor) VisitCaptureSpecifier(ctx *CaptureSpecifierContext) interface{}

func (v *BaseVertexParserVisitor) VisitCaseCondition(ctx *CaseConditionContext) interface{}

func (v *BaseVertexParserVisitor) VisitCaseItem(ctx *CaseItemContext) interface{}

func (v *BaseVertexParserVisitor) VisitCaseItemList(ctx *CaseItemListContext) interface{}

func (v *BaseVertexParserVisitor) VisitCaseLabel(ctx *CaseLabelContext) interface{}

func (v *BaseVertexParserVisitor) VisitCatchClause(ctx *CatchClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitCatchPattern(ctx *CatchPatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitCatchPatternList(ctx *CatchPatternListContext) interface{}

func (v *BaseVertexParserVisitor) VisitClassBody(ctx *ClassBodyContext) interface{}

func (v *BaseVertexParserVisitor) VisitClassDeclaration(ctx *ClassDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitClassMember(ctx *ClassMemberContext) interface{}

func (v *BaseVertexParserVisitor) VisitClosureExpr(ctx *ClosureExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitClosureExpression(ctx *ClosureExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitClosureParameter(ctx *ClosureParameterContext) interface{}

func (v *BaseVertexParserVisitor) VisitClosureParameterClause(ctx *ClosureParameterClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitClosureParameterList(ctx *ClosureParameterListContext) interface{}

func (v *BaseVertexParserVisitor) VisitClosureSignature(ctx *ClosureSignatureContext) interface{}

func (v *BaseVertexParserVisitor) VisitCodeBlock(ctx *CodeBlockContext) interface{}

func (v *BaseVertexParserVisitor) VisitCompilationCondition(ctx *CompilationConditionContext) interface{}

func (v *BaseVertexParserVisitor) VisitCompilerControl(ctx *CompilerControlContext) interface{}

func (v *BaseVertexParserVisitor) VisitCompilerControlStmt(ctx *CompilerControlStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitCondition(ctx *ConditionContext) interface{}

func (v *BaseVertexParserVisitor) VisitConditionList(ctx *ConditionListContext) interface{}

func (v *BaseVertexParserVisitor) VisitConditionalCompilationBlock(ctx *ConditionalCompilationBlockContext) interface{}

func (v *BaseVertexParserVisitor) VisitConditionalOperator(ctx *ConditionalOperatorContext) interface{}

func (v *BaseVertexParserVisitor) VisitConformanceRequirement(ctx *ConformanceRequirementContext) interface{}

func (v *BaseVertexParserVisitor) VisitConstantDeclaration(ctx *ConstantDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitContinueStatement(ctx *ContinueStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitControlTransferStmt(ctx *ControlTransferStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitDecimalVersion(ctx *DecimalVersionContext) interface{}

func (v *BaseVertexParserVisitor) VisitDeclaration(ctx *DeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitDeclarationModifier(ctx *DeclarationModifierContext) interface{}

func (v *BaseVertexParserVisitor) VisitDeclarationModifiers(ctx *DeclarationModifiersContext) interface{}

func (v *BaseVertexParserVisitor) VisitDeclarationStatement(ctx *DeclarationStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitDefaultArgumentClause(ctx *DefaultArgumentClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitDefaultLabel(ctx *DefaultLabelContext) interface{}

func (v *BaseVertexParserVisitor) VisitDeferStatement(ctx *DeferStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitDeferStmt(ctx *DeferStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitDeinitializerDeclaration(ctx *DeinitializerDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitDiagnosticStatement(ctx *DiagnosticStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitDictType(ctx *DictTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitDictionaryLiteral(ctx *DictionaryLiteralContext) interface{}

func (v *BaseVertexParserVisitor) VisitDictionaryLiteralItem(ctx *DictionaryLiteralItemContext) interface{}

func (v *BaseVertexParserVisitor) VisitDictionaryLiteralItems(ctx *DictionaryLiteralItemsContext) interface{}

func (v *BaseVertexParserVisitor) VisitDictionaryType(ctx *DictionaryTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitDidSetClause(ctx *DidSetClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitDoStatement(ctx *DoStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitDoStmt(ctx *DoStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitElseClause(ctx *ElseClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumCasePat(ctx *EnumCasePatContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumCasePattern(ctx *EnumCasePatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumDeclaration(ctx *EnumDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumMember(ctx *EnumMemberContext) interface{}

func (v *BaseVertexParserVisitor) VisitEnumMembers(ctx *EnumMembersContext) interface{}

func (v *BaseVertexParserVisitor) VisitExistType(ctx *ExistTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitExistentialType(ctx *ExistentialTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitExplicitMemberSuffix(ctx *ExplicitMemberSuffixContext) interface{}

func (v *BaseVertexParserVisitor) VisitExprPat(ctx *ExprPatContext) interface{}

func (v *BaseVertexParserVisitor) VisitExpression(ctx *ExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitExpressionPattern(ctx *ExpressionPatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitExtensionBody(ctx *ExtensionBodyContext) interface{}

func (v *BaseVertexParserVisitor) VisitExtensionDeclaration(ctx *ExtensionDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitExtensionMember(ctx *ExtensionMemberContext) interface{}

func (v *BaseVertexParserVisitor) VisitExternalParameterName(ctx *ExternalParameterNameContext) interface{}

func (v *BaseVertexParserVisitor) VisitFallthroughStatement(ctx *FallthroughStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitForInStatement(ctx *ForInStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitForStatement(ctx *ForStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitForcedValueSuffix(ctx *ForcedValueSuffixContext) interface{}

func (v *BaseVertexParserVisitor) VisitFuncType(ctx *FuncTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionBody(ctx *FunctionBodyContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionCallArgument(ctx *FunctionCallArgumentContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionCallArgumentClause(ctx *FunctionCallArgumentClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionCallArgumentList(ctx *FunctionCallArgumentListContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionCallSuffix(ctx *FunctionCallSuffixContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionDeclaration(ctx *FunctionDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionHead(ctx *FunctionHeadContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionResult(ctx *FunctionResultContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionSignature(ctx *FunctionSignatureContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionType(ctx *FunctionTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionTypeArgument(ctx *FunctionTypeArgumentContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionTypeArgumentClause(ctx *FunctionTypeArgumentClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitFunctionTypeArgumentList(ctx *FunctionTypeArgumentListContext) interface{}

func (v *BaseVertexParserVisitor) VisitGenericArgument(ctx *GenericArgumentContext) interface{}

func (v *BaseVertexParserVisitor) VisitGenericArgumentClause(ctx *GenericArgumentClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitGenericArgumentList(ctx *GenericArgumentListContext) interface{}

func (v *BaseVertexParserVisitor) VisitGenericParameter(ctx *GenericParameterContext) interface{}

func (v *BaseVertexParserVisitor) VisitGenericParameterClause(ctx *GenericParameterClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitGenericParameterList(ctx *GenericParameterListContext) interface{}

func (v *BaseVertexParserVisitor) VisitGenericWhereClause(ctx *GenericWhereClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitGetterClause(ctx *GetterClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitGetterKeywordClause(ctx *GetterKeywordClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitGetterSetterBlock(ctx *GetterSetterBlockContext) interface{}

func (v *BaseVertexParserVisitor) VisitGetterSetterKeywordBlock(ctx *GetterSetterKeywordBlockContext) interface{}

func (v *BaseVertexParserVisitor) VisitGuardStatement(ctx *GuardStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitIdentPat(ctx *IdentPatContext) interface{}

func (v *BaseVertexParserVisitor) VisitIdentifier(ctx *IdentifierContext) interface{}

func (v *BaseVertexParserVisitor) VisitIdentifierExpr(ctx *IdentifierExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitIdentifierList(ctx *IdentifierListContext) interface{}

func (v *BaseVertexParserVisitor) VisitIdentifierPattern(ctx *IdentifierPatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitIfStatement(ctx *IfStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitImplicitMemberExpr(ctx *ImplicitMemberExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitImplicitMemberExpression(ctx *ImplicitMemberExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitImportAlias(ctx *ImportAliasContext) interface{}

func (v *BaseVertexParserVisitor) VisitImportDeclaration(ctx *ImportDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitImportSpec(ctx *ImportSpecContext) interface{}

func (v *BaseVertexParserVisitor) VisitInOutExpression(ctx *InOutExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitInitializer(ctx *InitializerContext) interface{}

func (v *BaseVertexParserVisitor) VisitInitializerBody(ctx *InitializerBodyContext) interface{}

func (v *BaseVertexParserVisitor) VisitInitializerDeclaration(ctx *InitializerDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitInitializerSuffix(ctx *InitializerSuffixContext) interface{}

func (v *BaseVertexParserVisitor) VisitInout(ctx *InoutContext) interface{}

func (v *BaseVertexParserVisitor) VisitIsPat(ctx *IsPatContext) interface{}

func (v *BaseVertexParserVisitor) VisitLabelName(ctx *LabelNameContext) interface{}

func (v *BaseVertexParserVisitor) VisitLabeledStatement(ctx *LabeledStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitLabeledStmt(ctx *LabeledStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitLabeledTrailingClosure(ctx *LabeledTrailingClosureContext) interface{}

func (v *BaseVertexParserVisitor) VisitLayoutConstraintRequirement(ctx *LayoutConstraintRequirementContext) interface{}

func (v *BaseVertexParserVisitor) VisitLineControlStatement(ctx *LineControlStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitLitExpr(ctx *LitExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitLiteral(ctx *LiteralContext) interface{}

func (v *BaseVertexParserVisitor) VisitLiteralExpression(ctx *LiteralExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitLocalParameterName(ctx *LocalParameterNameContext) interface{}

func (v *BaseVertexParserVisitor) VisitLoopStatement(ctx *LoopStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitLoopStmt(ctx *LoopStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitMacroDeclStmt(ctx *MacroDeclStmtContext) interface{}

func (v *BaseVertexParserVisitor) VisitMacroDeclaration(ctx *MacroDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitMacroExpansionExpression(ctx *MacroExpansionExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitMacroExpr(ctx *MacroExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitMetatypeType_(ctx *MetatypeType_Context) interface{}

func (v *BaseVertexParserVisitor) VisitMutationModifier(ctx *MutationModifierContext) interface{}

func (v *BaseVertexParserVisitor) VisitNamedType(ctx *NamedTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitNumericLiteral(ctx *NumericLiteralContext) interface{}

func (v *BaseVertexParserVisitor) VisitOpaqueType(ctx *OpaqueTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitOpaqueType_(ctx *OpaqueType_Context) interface{}

func (v *BaseVertexParserVisitor) VisitOperator_(ctx *Operator_Context) interface{}

func (v *BaseVertexParserVisitor) VisitOptPat(ctx *OptPatContext) interface{}

func (v *BaseVertexParserVisitor) VisitOptType(ctx *OptTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitOptionalBindingCondition(ctx *OptionalBindingConditionContext) interface{}

func (v *BaseVertexParserVisitor) VisitOptionalChainingLiteral(ctx *OptionalChainingLiteralContext) interface{}

func (v *BaseVertexParserVisitor) VisitOptionalPattern(ctx *OptionalPatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitParameter(ctx *ParameterContext) interface{}

func (v *BaseVertexParserVisitor) VisitParameterClause(ctx *ParameterClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitParameterList(ctx *ParameterListContext) interface{}

func (v *BaseVertexParserVisitor) VisitParenExpr(ctx *ParenExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitParenthesizedExpression(ctx *ParenthesizedExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitPatternInitializer(ctx *PatternInitializerContext) interface{}

func (v *BaseVertexParserVisitor) VisitPatternInitializerList(ctx *PatternInitializerListContext) interface{}

func (v *BaseVertexParserVisitor) VisitPlatformCondition(ctx *PlatformConditionContext) interface{}

func (v *BaseVertexParserVisitor) VisitPostfixExpression(ctx *PostfixExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitPostfixOperator(ctx *PostfixOperatorContext) interface{}

func (v *BaseVertexParserVisitor) VisitPostfixSelfSuffix(ctx *PostfixSelfSuffixContext) interface{}

func (v *BaseVertexParserVisitor) VisitPostfixSuffix(ctx *PostfixSuffixContext) interface{}

func (v *BaseVertexParserVisitor) VisitPoundFileExpression(ctx *PoundFileExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitPrefixOp(ctx *PrefixOpContext) interface{}

func (v *BaseVertexParserVisitor) VisitPrefixOperator(ctx *PrefixOperatorContext) interface{}

func (v *BaseVertexParserVisitor) VisitPrimaryAssociatedTypeClause(ctx *PrimaryAssociatedTypeClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitPrimaryAssociatedTypeList(ctx *PrimaryAssociatedTypeListContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtoCompType(ctx *ProtoCompTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolAssociatedTypeDeclaration(ctx *ProtocolAssociatedTypeDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolBody(ctx *ProtocolBodyContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolCompositionType(ctx *ProtocolCompositionTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolCompositionTypeElement(ctx *ProtocolCompositionTypeElementContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolDeclaration(ctx *ProtocolDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolInitializerDeclaration(ctx *ProtocolInitializerDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolMember(ctx *ProtocolMemberContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolMemberDeclaration(ctx *ProtocolMemberDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolMethodDeclaration(ctx *ProtocolMethodDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolPropertyDeclaration(ctx *ProtocolPropertyDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitProtocolSubscriptDeclaration(ctx *ProtocolSubscriptDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitRawValueAssignment(ctx *RawValueAssignmentContext) interface{}

func (v *BaseVertexParserVisitor) VisitRawValueLiteral(ctx *RawValueLiteralContext) interface{}

func (v *BaseVertexParserVisitor) VisitRawValueStyleEnumCase(ctx *RawValueStyleEnumCaseContext) interface{}

func (v *BaseVertexParserVisitor) VisitRawValueStyleEnumCaseClause(ctx *RawValueStyleEnumCaseClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitRawValueStyleEnumCaseList(ctx *RawValueStyleEnumCaseListContext) interface{}

func (v *BaseVertexParserVisitor) VisitRepeatWhileStatement(ctx *RepeatWhileStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitRequirement(ctx *RequirementContext) interface{}

func (v *BaseVertexParserVisitor) VisitRequirementList(ctx *RequirementListContext) interface{}

func (v *BaseVertexParserVisitor) VisitReturnStatement(ctx *ReturnStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitSameTypeRequirement(ctx *SameTypeRequirementContext) interface{}

func (v *BaseVertexParserVisitor) VisitSelfExpr(ctx *SelfExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitSelfExpression(ctx *SelfExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitSelfType(ctx *SelfTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitSelfType_(ctx *SelfType_Context) interface{}

func (v *BaseVertexParserVisitor) VisitSetterClause(ctx *SetterClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitSetterKeywordClause(ctx *SetterKeywordClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitSetterName(ctx *SetterNameContext) interface{}

func (v *BaseVertexParserVisitor) VisitStatements(ctx *StatementsContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructBody(ctx *StructBodyContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructDeclaration(ctx *StructDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitStructMember(ctx *StructMemberContext) interface{}

func (v *BaseVertexParserVisitor) VisitSubscriptDeclaration(ctx *SubscriptDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitSubscriptHead(ctx *SubscriptHeadContext) interface{}

func (v *BaseVertexParserVisitor) VisitSubscriptResult(ctx *SubscriptResultContext) interface{}

func (v *BaseVertexParserVisitor) VisitSubscriptSuffix(ctx *SubscriptSuffixContext) interface{}

func (v *BaseVertexParserVisitor) VisitSuperExpr(ctx *SuperExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitSuperExpression(ctx *SuperExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitSuppressedType(ctx *SuppressedTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{}

func (v *BaseVertexParserVisitor) VisitSwitchCases(ctx *SwitchCasesContext) interface{}

func (v *BaseVertexParserVisitor) VisitSwitchStatement(ctx *SwitchStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitTernaryExpr(ctx *TernaryExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitThrowStatement(ctx *ThrowStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitTopLevel(ctx *TopLevelContext) interface{}

func (v *BaseVertexParserVisitor) VisitTrailingClosures(ctx *TrailingClosuresContext) interface{}

func (v *BaseVertexParserVisitor) VisitTryOperator(ctx *TryOperatorContext) interface{}

func (v *BaseVertexParserVisitor) VisitTryStatement(ctx *TryStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupType(ctx *TupTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleElement(ctx *TupleElementContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleElementList(ctx *TupleElementListContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleExpr(ctx *TupleExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleExpression(ctx *TupleExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitTuplePat(ctx *TuplePatContext) interface{}

func (v *BaseVertexParserVisitor) VisitTuplePattern(ctx *TuplePatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitTuplePatternElement(ctx *TuplePatternElementContext) interface{}

func (v *BaseVertexParserVisitor) VisitTuplePatternElementList(ctx *TuplePatternElementListContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleType(ctx *TupleTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleTypeElement(ctx *TupleTypeElementContext) interface{}

func (v *BaseVertexParserVisitor) VisitTupleTypeElementList(ctx *TupleTypeElementListContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeAnnotation(ctx *TypeAnnotationContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeAnnotationHead(ctx *TypeAnnotationHeadContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeCastExpr(ctx *TypeCastExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeCastingOperator(ctx *TypeCastingOperatorContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeIdentifier(ctx *TypeIdentifierContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeInheritanceClause(ctx *TypeInheritanceClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypeInheritanceList(ctx *TypeInheritanceListContext) interface{}

func (v *BaseVertexParserVisitor) VisitTypealiasDeclaration(ctx *TypealiasDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitUnionStyleEnumCase(ctx *UnionStyleEnumCaseContext) interface{}

func (v *BaseVertexParserVisitor) VisitUnionStyleEnumCaseClause(ctx *UnionStyleEnumCaseClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitUnionStyleEnumCaseList(ctx *UnionStyleEnumCaseListContext) interface{}

func (v *BaseVertexParserVisitor) VisitValueBindingPattern(ctx *ValueBindingPatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{}

func (v *BaseVertexParserVisitor) VisitVariableDeclarationHead(ctx *VariableDeclarationHeadContext) interface{}

func (v *BaseVertexParserVisitor) VisitVariableName(ctx *VariableNameContext) interface{}

func (v *BaseVertexParserVisitor) VisitVariadicType(ctx *VariadicTypeContext) interface{}

func (v *BaseVertexParserVisitor) VisitWhereClause(ctx *WhereClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitWhileStatement(ctx *WhileStatementContext) interface{}

func (v *BaseVertexParserVisitor) VisitWildcardExpr(ctx *WildcardExprContext) interface{}

func (v *BaseVertexParserVisitor) VisitWildcardExpression(ctx *WildcardExpressionContext) interface{}

func (v *BaseVertexParserVisitor) VisitWildcardPat(ctx *WildcardPatContext) interface{}

func (v *BaseVertexParserVisitor) VisitWildcardPattern(ctx *WildcardPatternContext) interface{}

func (v *BaseVertexParserVisitor) VisitWillSetClause(ctx *WillSetClauseContext) interface{}

func (v *BaseVertexParserVisitor) VisitWillSetDidSetBlock(ctx *WillSetDidSetBlockContext) interface{}

type BinaryExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBinaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryExpressionContext

func NewEmptyBinaryExpressionContext() *BinaryExpressionContext

func (s *BinaryExpressionContext) CopyAll(ctx *BinaryExpressionContext)

func (s *BinaryExpressionContext) GetParser() antlr.Parser

func (s *BinaryExpressionContext) GetRuleContext() antlr.RuleContext

func (*BinaryExpressionContext) IsBinaryExpressionContext()

func (s *BinaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BinaryExpressionsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBinaryExpressionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryExpressionsContext

func NewEmptyBinaryExpressionsContext() *BinaryExpressionsContext

func (s *BinaryExpressionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BinaryExpressionsContext) AllBinaryExpression() []IBinaryExpressionContext

func (s *BinaryExpressionsContext) BinaryExpression(i int) IBinaryExpressionContext

func (s *BinaryExpressionsContext) GetParser() antlr.Parser

func (s *BinaryExpressionsContext) GetRuleContext() antlr.RuleContext

func (*BinaryExpressionsContext) IsBinaryExpressionsContext()

func (s *BinaryExpressionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BinaryOpContext struct {
	BinaryExpressionContext
}

func NewBinaryOpContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BinaryOpContext

func (s *BinaryOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BinaryOpContext) BinaryOperator() IBinaryOperatorContext

func (s *BinaryOpContext) GetRuleContext() antlr.RuleContext

func (s *BinaryOpContext) PrefixExpression() IPrefixExpressionContext

type BinaryOperatorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBinaryOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryOperatorContext

func NewEmptyBinaryOperatorContext() *BinaryOperatorContext

func (s *BinaryOperatorContext) AMPERSAND() antlr.TerminalNode

func (s *BinaryOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BinaryOperatorContext) DOT() antlr.TerminalNode

func (s *BinaryOperatorContext) GT() antlr.TerminalNode

func (s *BinaryOperatorContext) GetParser() antlr.Parser

func (s *BinaryOperatorContext) GetRuleContext() antlr.RuleContext

func (*BinaryOperatorContext) IsBinaryOperatorContext()

func (s *BinaryOperatorContext) LT() antlr.TerminalNode

func (s *BinaryOperatorContext) OPERATOR() antlr.TerminalNode

func (s *BinaryOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BindingPatContext struct {
	PatternContext
}

func NewBindingPatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BindingPatContext

func (s *BindingPatContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BindingPatContext) GetRuleContext() antlr.RuleContext

func (s *BindingPatContext) ValueBindingPattern() IValueBindingPatternContext

type BooleanLiteralContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBooleanLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanLiteralContext

func NewEmptyBooleanLiteralContext() *BooleanLiteralContext

func (s *BooleanLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BooleanLiteralContext) FALSE() antlr.TerminalNode

func (s *BooleanLiteralContext) GetParser() antlr.Parser

func (s *BooleanLiteralContext) GetRuleContext() antlr.RuleContext

func (*BooleanLiteralContext) IsBooleanLiteralContext()

func (s *BooleanLiteralContext) TRUE() antlr.TerminalNode

func (s *BooleanLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BranchStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBranchStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BranchStatementContext

func NewEmptyBranchStatementContext() *BranchStatementContext

func (s *BranchStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BranchStatementContext) GetParser() antlr.Parser

func (s *BranchStatementContext) GetRuleContext() antlr.RuleContext

func (s *BranchStatementContext) GuardStatement() IGuardStatementContext

func (s *BranchStatementContext) IfStatement() IIfStatementContext

func (*BranchStatementContext) IsBranchStatementContext()

func (s *BranchStatementContext) SwitchStatement() ISwitchStatementContext

func (s *BranchStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *BranchStatementContext) TryStatement() ITryStatementContext

type BranchStmtContext struct {
	StatementContext
}

func NewBranchStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BranchStmtContext

func (s *BranchStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BranchStmtContext) BranchStatement() IBranchStatementContext

func (s *BranchStmtContext) GetRuleContext() antlr.RuleContext

func (s *BranchStmtContext) SEMICOLON() antlr.TerminalNode

type BreakStatementContext struct {
	ControlTransferContext
}

func NewBreakStatementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BreakStatementContext

func (s *BreakStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BreakStatementContext) BREAK() antlr.TerminalNode

func (s *BreakStatementContext) GetRuleContext() antlr.RuleContext

func (s *BreakStatementContext) LabelName() ILabelNameContext

type CaptureListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCaptureListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaptureListContext

func NewEmptyCaptureListContext() *CaptureListContext

func (s *CaptureListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CaptureListContext) CaptureListItems() ICaptureListItemsContext

func (s *CaptureListContext) GetParser() antlr.Parser

func (s *CaptureListContext) GetRuleContext() antlr.RuleContext

func (*CaptureListContext) IsCaptureListContext()

func (s *CaptureListContext) LBRACKET() antlr.TerminalNode

func (s *CaptureListContext) RBRACKET() antlr.TerminalNode

func (s *CaptureListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CaptureListItemContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCaptureListItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaptureListItemContext

func NewEmptyCaptureListItemContext() *CaptureListItemContext

func (s *CaptureListItemContext) ASSIGN() antlr.TerminalNode

func (s *CaptureListItemContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CaptureListItemContext) CaptureSpecifier() ICaptureSpecifierContext

func (s *CaptureListItemContext) Expression() IExpressionContext

func (s *CaptureListItemContext) GetParser() antlr.Parser

func (s *CaptureListItemContext) GetRuleContext() antlr.RuleContext

func (s *CaptureListItemContext) Identifier() IIdentifierContext

func (*CaptureListItemContext) IsCaptureListItemContext()

func (s *CaptureListItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CaptureListItemsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCaptureListItemsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaptureListItemsContext

func NewEmptyCaptureListItemsContext() *CaptureListItemsContext

func (s *CaptureListItemsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CaptureListItemsContext) AllCOMMA() []antlr.TerminalNode

func (s *CaptureListItemsContext) AllCaptureListItem() []ICaptureListItemContext

func (s *CaptureListItemsContext) COMMA(i int) antlr.TerminalNode

func (s *CaptureListItemsContext) CaptureListItem(i int) ICaptureListItemContext

func (s *CaptureListItemsContext) GetParser() antlr.Parser

func (s *CaptureListItemsContext) GetRuleContext() antlr.RuleContext

func (*CaptureListItemsContext) IsCaptureListItemsContext()

func (s *CaptureListItemsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CaptureSpecifierContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCaptureSpecifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaptureSpecifierContext

func NewEmptyCaptureSpecifierContext() *CaptureSpecifierContext

func (s *CaptureSpecifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CaptureSpecifierContext) GetParser() antlr.Parser

func (s *CaptureSpecifierContext) GetRuleContext() antlr.RuleContext

func (*CaptureSpecifierContext) IsCaptureSpecifierContext()

func (s *CaptureSpecifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *CaptureSpecifierContext) UNOWNED_KW() antlr.TerminalNode

func (s *CaptureSpecifierContext) WEAK_KW() antlr.TerminalNode

type CaseConditionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCaseConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaseConditionContext

func NewEmptyCaseConditionContext() *CaseConditionContext

func (s *CaseConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CaseConditionContext) CASE() antlr.TerminalNode

func (s *CaseConditionContext) GetParser() antlr.Parser

func (s *CaseConditionContext) GetRuleContext() antlr.RuleContext

func (s *CaseConditionContext) Initializer() IInitializerContext

func (*CaseConditionContext) IsCaseConditionContext()

func (s *CaseConditionContext) Pattern() IPatternContext

func (s *CaseConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CaseItemContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCaseItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaseItemContext

func NewEmptyCaseItemContext() *CaseItemContext

func (s *CaseItemContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CaseItemContext) GetParser() antlr.Parser

func (s *CaseItemContext) GetRuleContext() antlr.RuleContext

func (*CaseItemContext) IsCaseItemContext()

func (s *CaseItemContext) Pattern() IPatternContext

func (s *CaseItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *CaseItemContext) WhereClause() IWhereClauseContext

type CaseItemListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCaseItemListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaseItemListContext

func NewEmptyCaseItemListContext() *CaseItemListContext

func (s *CaseItemListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CaseItemListContext) AllCOMMA() []antlr.TerminalNode

func (s *CaseItemListContext) AllCaseItem() []ICaseItemContext

func (s *CaseItemListContext) COMMA(i int) antlr.TerminalNode

func (s *CaseItemListContext) CaseItem(i int) ICaseItemContext

func (s *CaseItemListContext) GetParser() antlr.Parser

func (s *CaseItemListContext) GetRuleContext() antlr.RuleContext

func (*CaseItemListContext) IsCaseItemListContext()

func (s *CaseItemListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CaseLabelContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCaseLabelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaseLabelContext

func NewEmptyCaseLabelContext() *CaseLabelContext

func (s *CaseLabelContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CaseLabelContext) CASE() antlr.TerminalNode

func (s *CaseLabelContext) COLON() antlr.TerminalNode

func (s *CaseLabelContext) CaseItemList() ICaseItemListContext

func (s *CaseLabelContext) GetParser() antlr.Parser

func (s *CaseLabelContext) GetRuleContext() antlr.RuleContext

func (*CaseLabelContext) IsCaseLabelContext()

func (s *CaseLabelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CatchClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCatchClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CatchClauseContext

func NewEmptyCatchClauseContext() *CatchClauseContext

func (s *CatchClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CatchClauseContext) CATCH() antlr.TerminalNode

func (s *CatchClauseContext) CatchPatternList() ICatchPatternListContext

func (s *CatchClauseContext) CodeBlock() ICodeBlockContext

func (s *CatchClauseContext) GetParser() antlr.Parser

func (s *CatchClauseContext) GetRuleContext() antlr.RuleContext

func (*CatchClauseContext) IsCatchClauseContext()

func (s *CatchClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CatchPatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCatchPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CatchPatternContext

func NewEmptyCatchPatternContext() *CatchPatternContext

func (s *CatchPatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CatchPatternContext) GetParser() antlr.Parser

func (s *CatchPatternContext) GetRuleContext() antlr.RuleContext

func (*CatchPatternContext) IsCatchPatternContext()

func (s *CatchPatternContext) Pattern() IPatternContext

func (s *CatchPatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *CatchPatternContext) WhereClause() IWhereClauseContext

type CatchPatternListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCatchPatternListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CatchPatternListContext

func NewEmptyCatchPatternListContext() *CatchPatternListContext

func (s *CatchPatternListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CatchPatternListContext) AllCOMMA() []antlr.TerminalNode

func (s *CatchPatternListContext) AllCatchPattern() []ICatchPatternContext

func (s *CatchPatternListContext) COMMA(i int) antlr.TerminalNode

func (s *CatchPatternListContext) CatchPattern(i int) ICatchPatternContext

func (s *CatchPatternListContext) GetParser() antlr.Parser

func (s *CatchPatternListContext) GetRuleContext() antlr.RuleContext

func (*CatchPatternListContext) IsCatchPatternListContext()

func (s *CatchPatternListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClassBodyContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClassBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassBodyContext

func NewEmptyClassBodyContext() *ClassBodyContext

func (s *ClassBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClassBodyContext) AllClassMember() []IClassMemberContext

func (s *ClassBodyContext) ClassMember(i int) IClassMemberContext

func (s *ClassBodyContext) GetParser() antlr.Parser

func (s *ClassBodyContext) GetRuleContext() antlr.RuleContext

func (*ClassBodyContext) IsClassBodyContext()

func (s *ClassBodyContext) LBRACE() antlr.TerminalNode

func (s *ClassBodyContext) RBRACE() antlr.TerminalNode

func (s *ClassBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClassDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClassDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassDeclarationContext

func NewEmptyClassDeclarationContext() *ClassDeclarationContext

func (s *ClassDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClassDeclarationContext) AccessLevelModifier() IAccessLevelModifierContext

func (s *ClassDeclarationContext) Attributes() IAttributesContext

func (s *ClassDeclarationContext) CLASS() antlr.TerminalNode

func (s *ClassDeclarationContext) ClassBody() IClassBodyContext

func (s *ClassDeclarationContext) FINAL_KW() antlr.TerminalNode

func (s *ClassDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *ClassDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *ClassDeclarationContext) GetParser() antlr.Parser

func (s *ClassDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *ClassDeclarationContext) Identifier() IIdentifierContext

func (*ClassDeclarationContext) IsClassDeclarationContext()

func (s *ClassDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ClassDeclarationContext) TypeInheritanceClause() ITypeInheritanceClauseContext

type ClassMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClassMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClassMemberContext

func NewEmptyClassMemberContext() *ClassMemberContext

func (s *ClassMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClassMemberContext) CompilerControl() ICompilerControlContext

func (s *ClassMemberContext) Declaration() IDeclarationContext

func (s *ClassMemberContext) GetParser() antlr.Parser

func (s *ClassMemberContext) GetRuleContext() antlr.RuleContext

func (*ClassMemberContext) IsClassMemberContext()

func (s *ClassMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClosureExprContext struct {
	PrimaryExpressionContext
}

func NewClosureExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ClosureExprContext

func (s *ClosureExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClosureExprContext) ClosureExpression() IClosureExpressionContext

func (s *ClosureExprContext) GetRuleContext() antlr.RuleContext

type ClosureExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClosureExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClosureExpressionContext

func NewEmptyClosureExpressionContext() *ClosureExpressionContext

func (s *ClosureExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClosureExpressionContext) ClosureSignature() IClosureSignatureContext

func (s *ClosureExpressionContext) GetParser() antlr.Parser

func (s *ClosureExpressionContext) GetRuleContext() antlr.RuleContext

func (*ClosureExpressionContext) IsClosureExpressionContext()

func (s *ClosureExpressionContext) LBRACE() antlr.TerminalNode

func (s *ClosureExpressionContext) RBRACE() antlr.TerminalNode

func (s *ClosureExpressionContext) Statements() IStatementsContext

func (s *ClosureExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClosureParameterClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClosureParameterClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClosureParameterClauseContext

func NewEmptyClosureParameterClauseContext() *ClosureParameterClauseContext

func (s *ClosureParameterClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClosureParameterClauseContext) ClosureParameterList() IClosureParameterListContext

func (s *ClosureParameterClauseContext) GetParser() antlr.Parser

func (s *ClosureParameterClauseContext) GetRuleContext() antlr.RuleContext

func (s *ClosureParameterClauseContext) IdentifierList() IIdentifierListContext

func (*ClosureParameterClauseContext) IsClosureParameterClauseContext()

func (s *ClosureParameterClauseContext) LPAREN() antlr.TerminalNode

func (s *ClosureParameterClauseContext) RPAREN() antlr.TerminalNode

func (s *ClosureParameterClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClosureParameterContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClosureParameterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClosureParameterContext

func NewEmptyClosureParameterContext() *ClosureParameterContext

func (s *ClosureParameterContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClosureParameterContext) ELLIPSIS() antlr.TerminalNode

func (s *ClosureParameterContext) GetParser() antlr.Parser

func (s *ClosureParameterContext) GetRuleContext() antlr.RuleContext

func (s *ClosureParameterContext) Identifier() IIdentifierContext

func (*ClosureParameterContext) IsClosureParameterContext()

func (s *ClosureParameterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ClosureParameterContext) TypeAnnotation() ITypeAnnotationContext

type ClosureParameterListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClosureParameterListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClosureParameterListContext

func NewEmptyClosureParameterListContext() *ClosureParameterListContext

func (s *ClosureParameterListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClosureParameterListContext) AllCOMMA() []antlr.TerminalNode

func (s *ClosureParameterListContext) AllClosureParameter() []IClosureParameterContext

func (s *ClosureParameterListContext) COMMA(i int) antlr.TerminalNode

func (s *ClosureParameterListContext) ClosureParameter(i int) IClosureParameterContext

func (s *ClosureParameterListContext) GetParser() antlr.Parser

func (s *ClosureParameterListContext) GetRuleContext() antlr.RuleContext

func (*ClosureParameterListContext) IsClosureParameterListContext()

func (s *ClosureParameterListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ClosureSignatureContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewClosureSignatureContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClosureSignatureContext

func NewEmptyClosureSignatureContext() *ClosureSignatureContext

func (s *ClosureSignatureContext) ASYNC_KW() antlr.TerminalNode

func (s *ClosureSignatureContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ClosureSignatureContext) CaptureList() ICaptureListContext

func (s *ClosureSignatureContext) ClosureParameterClause() IClosureParameterClauseContext

func (s *ClosureSignatureContext) FunctionResult() IFunctionResultContext

func (s *ClosureSignatureContext) GetParser() antlr.Parser

func (s *ClosureSignatureContext) GetRuleContext() antlr.RuleContext

func (s *ClosureSignatureContext) IN() antlr.TerminalNode

func (*ClosureSignatureContext) IsClosureSignatureContext()

func (s *ClosureSignatureContext) THROWS() antlr.TerminalNode

func (s *ClosureSignatureContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CodeBlockContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCodeBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CodeBlockContext

func NewEmptyCodeBlockContext() *CodeBlockContext

func (s *CodeBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CodeBlockContext) GetParser() antlr.Parser

func (s *CodeBlockContext) GetRuleContext() antlr.RuleContext

func (*CodeBlockContext) IsCodeBlockContext()

func (s *CodeBlockContext) LBRACE() antlr.TerminalNode

func (s *CodeBlockContext) RBRACE() antlr.TerminalNode

func (s *CodeBlockContext) Statements() IStatementsContext

func (s *CodeBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CompilationConditionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCompilationConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompilationConditionContext

func NewEmptyCompilationConditionContext() *CompilationConditionContext

func (s *CompilationConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CompilationConditionContext) AllCompilationCondition() []ICompilationConditionContext

func (s *CompilationConditionContext) BooleanLiteral() IBooleanLiteralContext

func (s *CompilationConditionContext) CompilationCondition(i int) ICompilationConditionContext

func (s *CompilationConditionContext) GetParser() antlr.Parser

func (s *CompilationConditionContext) GetRuleContext() antlr.RuleContext

func (s *CompilationConditionContext) Identifier() IIdentifierContext

func (*CompilationConditionContext) IsCompilationConditionContext()

func (s *CompilationConditionContext) LPAREN() antlr.TerminalNode

func (s *CompilationConditionContext) OPERATOR() antlr.TerminalNode

func (s *CompilationConditionContext) PlatformCondition() IPlatformConditionContext

func (s *CompilationConditionContext) RPAREN() antlr.TerminalNode

func (s *CompilationConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CompilerControlContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCompilerControlContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompilerControlContext

func NewEmptyCompilerControlContext() *CompilerControlContext

func (s *CompilerControlContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CompilerControlContext) ConditionalCompilationBlock() IConditionalCompilationBlockContext

func (s *CompilerControlContext) DiagnosticStatement() IDiagnosticStatementContext

func (s *CompilerControlContext) GetParser() antlr.Parser

func (s *CompilerControlContext) GetRuleContext() antlr.RuleContext

func (*CompilerControlContext) IsCompilerControlContext()

func (s *CompilerControlContext) LineControlStatement() ILineControlStatementContext

func (s *CompilerControlContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CompilerControlStmtContext struct {
	StatementContext
}

func NewCompilerControlStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CompilerControlStmtContext

func (s *CompilerControlStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CompilerControlStmtContext) CompilerControl() ICompilerControlContext

func (s *CompilerControlStmtContext) GetRuleContext() antlr.RuleContext

func (s *CompilerControlStmtContext) SEMICOLON() antlr.TerminalNode

type ConditionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionContext

func NewEmptyConditionContext() *ConditionContext

func (s *ConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ConditionContext) AvailabilityCondition() IAvailabilityConditionContext

func (s *ConditionContext) CaseCondition() ICaseConditionContext

func (s *ConditionContext) Expression() IExpressionContext

func (s *ConditionContext) GetParser() antlr.Parser

func (s *ConditionContext) GetRuleContext() antlr.RuleContext

func (*ConditionContext) IsConditionContext()

func (s *ConditionContext) OptionalBindingCondition() IOptionalBindingConditionContext

func (s *ConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ConditionListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewConditionListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionListContext

func NewEmptyConditionListContext() *ConditionListContext

func (s *ConditionListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ConditionListContext) AllCOMMA() []antlr.TerminalNode

func (s *ConditionListContext) AllCondition() []IConditionContext

func (s *ConditionListContext) COMMA(i int) antlr.TerminalNode

func (s *ConditionListContext) Condition(i int) IConditionContext

func (s *ConditionListContext) GetParser() antlr.Parser

func (s *ConditionListContext) GetRuleContext() antlr.RuleContext

func (*ConditionListContext) IsConditionListContext()

func (s *ConditionListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ConditionalCompilationBlockContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewConditionalCompilationBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalCompilationBlockContext

func NewEmptyConditionalCompilationBlockContext() *ConditionalCompilationBlockContext

func (s *ConditionalCompilationBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ConditionalCompilationBlockContext) AllCompilationCondition() []ICompilationConditionContext

func (s *ConditionalCompilationBlockContext) AllPOUND_ELSEIF() []antlr.TerminalNode

func (s *ConditionalCompilationBlockContext) AllStatements() []IStatementsContext

func (s *ConditionalCompilationBlockContext) CompilationCondition(i int) ICompilationConditionContext

func (s *ConditionalCompilationBlockContext) GetParser() antlr.Parser

func (s *ConditionalCompilationBlockContext) GetRuleContext() antlr.RuleContext

func (*ConditionalCompilationBlockContext) IsConditionalCompilationBlockContext()

func (s *ConditionalCompilationBlockContext) POUND_ELSE() antlr.TerminalNode

func (s *ConditionalCompilationBlockContext) POUND_ELSEIF(i int) antlr.TerminalNode

func (s *ConditionalCompilationBlockContext) POUND_ENDIF() antlr.TerminalNode

func (s *ConditionalCompilationBlockContext) POUND_IF() antlr.TerminalNode

func (s *ConditionalCompilationBlockContext) Statements(i int) IStatementsContext

func (s *ConditionalCompilationBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ConditionalOperatorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewConditionalOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalOperatorContext

func NewEmptyConditionalOperatorContext() *ConditionalOperatorContext

func (s *ConditionalOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ConditionalOperatorContext) COLON() antlr.TerminalNode

func (s *ConditionalOperatorContext) Expression() IExpressionContext

func (s *ConditionalOperatorContext) GetParser() antlr.Parser

func (s *ConditionalOperatorContext) GetRuleContext() antlr.RuleContext

func (*ConditionalOperatorContext) IsConditionalOperatorContext()

func (s *ConditionalOperatorContext) QUESTION_POSTFIX() antlr.TerminalNode

func (s *ConditionalOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ConformanceRequirementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewConformanceRequirementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConformanceRequirementContext

func NewEmptyConformanceRequirementContext() *ConformanceRequirementContext

func (s *ConformanceRequirementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ConformanceRequirementContext) AllTypeIdentifier() []ITypeIdentifierContext

func (s *ConformanceRequirementContext) COLON() antlr.TerminalNode

func (s *ConformanceRequirementContext) GetParser() antlr.Parser

func (s *ConformanceRequirementContext) GetRuleContext() antlr.RuleContext

func (*ConformanceRequirementContext) IsConformanceRequirementContext()

func (s *ConformanceRequirementContext) ProtocolCompositionType() IProtocolCompositionTypeContext

func (s *ConformanceRequirementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ConformanceRequirementContext) TypeIdentifier(i int) ITypeIdentifierContext

type ConstantDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewConstantDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstantDeclarationContext

func NewEmptyConstantDeclarationContext() *ConstantDeclarationContext

func (s *ConstantDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ConstantDeclarationContext) Attributes() IAttributesContext

func (s *ConstantDeclarationContext) DeclarationModifiers() IDeclarationModifiersContext

func (s *ConstantDeclarationContext) GetParser() antlr.Parser

func (s *ConstantDeclarationContext) GetRuleContext() antlr.RuleContext

func (*ConstantDeclarationContext) IsConstantDeclarationContext()

func (s *ConstantDeclarationContext) LET() antlr.TerminalNode

func (s *ConstantDeclarationContext) PatternInitializerList() IPatternInitializerListContext

func (s *ConstantDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ContinueStatementContext struct {
	ControlTransferContext
}

func NewContinueStatementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ContinueStatementContext

func (s *ContinueStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ContinueStatementContext) CONTINUE() antlr.TerminalNode

func (s *ContinueStatementContext) GetRuleContext() antlr.RuleContext

func (s *ContinueStatementContext) LabelName() ILabelNameContext

type ControlTransferContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewControlTransferContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ControlTransferContext

func NewEmptyControlTransferContext() *ControlTransferContext

func (s *ControlTransferContext) CopyAll(ctx *ControlTransferContext)

func (s *ControlTransferContext) GetParser() antlr.Parser

func (s *ControlTransferContext) GetRuleContext() antlr.RuleContext

func (*ControlTransferContext) IsControlTransferContext()

func (s *ControlTransferContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ControlTransferStmtContext struct {
	StatementContext
}

func NewControlTransferStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ControlTransferStmtContext

func (s *ControlTransferStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ControlTransferStmtContext) ControlTransfer() IControlTransferContext

func (s *ControlTransferStmtContext) GetRuleContext() antlr.RuleContext

func (s *ControlTransferStmtContext) SEMICOLON() antlr.TerminalNode

type DecimalVersionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDecimalVersionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DecimalVersionContext

func NewEmptyDecimalVersionContext() *DecimalVersionContext

func (s *DecimalVersionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DecimalVersionContext) AllDOT() []antlr.TerminalNode

func (s *DecimalVersionContext) AllINTEGER_LITERAL() []antlr.TerminalNode

func (s *DecimalVersionContext) DOT(i int) antlr.TerminalNode

func (s *DecimalVersionContext) GetParser() antlr.Parser

func (s *DecimalVersionContext) GetRuleContext() antlr.RuleContext

func (s *DecimalVersionContext) INTEGER_LITERAL(i int) antlr.TerminalNode

func (*DecimalVersionContext) IsDecimalVersionContext()

func (s *DecimalVersionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclarationContext

func NewEmptyDeclarationContext() *DeclarationContext

func (s *DeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeclarationContext) ActorDeclaration() IActorDeclarationContext

func (s *DeclarationContext) ClassDeclaration() IClassDeclarationContext

func (s *DeclarationContext) ConstantDeclaration() IConstantDeclarationContext

func (s *DeclarationContext) DeinitializerDeclaration() IDeinitializerDeclarationContext

func (s *DeclarationContext) EnumDeclaration() IEnumDeclarationContext

func (s *DeclarationContext) ExtensionDeclaration() IExtensionDeclarationContext

func (s *DeclarationContext) FunctionDeclaration() IFunctionDeclarationContext

func (s *DeclarationContext) GetParser() antlr.Parser

func (s *DeclarationContext) GetRuleContext() antlr.RuleContext

func (s *DeclarationContext) ImportDeclaration() IImportDeclarationContext

func (s *DeclarationContext) InitializerDeclaration() IInitializerDeclarationContext

func (*DeclarationContext) IsDeclarationContext()

func (s *DeclarationContext) ProtocolDeclaration() IProtocolDeclarationContext

func (s *DeclarationContext) StructDeclaration() IStructDeclarationContext

func (s *DeclarationContext) SubscriptDeclaration() ISubscriptDeclarationContext

func (s *DeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *DeclarationContext) TypealiasDeclaration() ITypealiasDeclarationContext

func (s *DeclarationContext) VariableDeclaration() IVariableDeclarationContext

type DeclarationModifierContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDeclarationModifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclarationModifierContext

func NewEmptyDeclarationModifierContext() *DeclarationModifierContext

func (s *DeclarationModifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeclarationModifierContext) AccessLevelModifier() IAccessLevelModifierContext

func (s *DeclarationModifierContext) CLASS() antlr.TerminalNode

func (s *DeclarationModifierContext) FINAL_KW() antlr.TerminalNode

func (s *DeclarationModifierContext) GetParser() antlr.Parser

func (s *DeclarationModifierContext) GetRuleContext() antlr.RuleContext

func (*DeclarationModifierContext) IsDeclarationModifierContext()

func (s *DeclarationModifierContext) LAZY_KW() antlr.TerminalNode

func (s *DeclarationModifierContext) MutationModifier() IMutationModifierContext

func (s *DeclarationModifierContext) NONISOLATED_KW() antlr.TerminalNode

func (s *DeclarationModifierContext) OPTIONAL_KW() antlr.TerminalNode

func (s *DeclarationModifierContext) OVERRIDE_KW() antlr.TerminalNode

func (s *DeclarationModifierContext) POSTFIX_KW() antlr.TerminalNode

func (s *DeclarationModifierContext) PREFIX_KW() antlr.TerminalNode

func (s *DeclarationModifierContext) REQUIRED_KW() antlr.TerminalNode

func (s *DeclarationModifierContext) STATIC() antlr.TerminalNode

func (s *DeclarationModifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *DeclarationModifierContext) UNOWNED_KW() antlr.TerminalNode

func (s *DeclarationModifierContext) WEAK_KW() antlr.TerminalNode

type DeclarationModifiersContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDeclarationModifiersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclarationModifiersContext

func NewEmptyDeclarationModifiersContext() *DeclarationModifiersContext

func (s *DeclarationModifiersContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeclarationModifiersContext) AllDeclarationModifier() []IDeclarationModifierContext

func (s *DeclarationModifiersContext) DeclarationModifier(i int) IDeclarationModifierContext

func (s *DeclarationModifiersContext) GetParser() antlr.Parser

func (s *DeclarationModifiersContext) GetRuleContext() antlr.RuleContext

func (*DeclarationModifiersContext) IsDeclarationModifiersContext()

func (s *DeclarationModifiersContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DeclarationStatementContext struct {
	StatementContext
}

func NewDeclarationStatementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DeclarationStatementContext

func (s *DeclarationStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeclarationStatementContext) Declaration() IDeclarationContext

func (s *DeclarationStatementContext) GetRuleContext() antlr.RuleContext

func (s *DeclarationStatementContext) SEMICOLON() antlr.TerminalNode

type DefaultArgumentClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDefaultArgumentClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DefaultArgumentClauseContext

func NewEmptyDefaultArgumentClauseContext() *DefaultArgumentClauseContext

func (s *DefaultArgumentClauseContext) ASSIGN() antlr.TerminalNode

func (s *DefaultArgumentClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DefaultArgumentClauseContext) Expression() IExpressionContext

func (s *DefaultArgumentClauseContext) GetParser() antlr.Parser

func (s *DefaultArgumentClauseContext) GetRuleContext() antlr.RuleContext

func (*DefaultArgumentClauseContext) IsDefaultArgumentClauseContext()

func (s *DefaultArgumentClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DefaultLabelContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDefaultLabelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DefaultLabelContext

func NewEmptyDefaultLabelContext() *DefaultLabelContext

func (s *DefaultLabelContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DefaultLabelContext) COLON() antlr.TerminalNode

func (s *DefaultLabelContext) DEFAULT() antlr.TerminalNode

func (s *DefaultLabelContext) GetParser() antlr.Parser

func (s *DefaultLabelContext) GetRuleContext() antlr.RuleContext

func (*DefaultLabelContext) IsDefaultLabelContext()

func (s *DefaultLabelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DeferStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDeferStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeferStatementContext

func NewEmptyDeferStatementContext() *DeferStatementContext

func (s *DeferStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeferStatementContext) CodeBlock() ICodeBlockContext

func (s *DeferStatementContext) DEFER() antlr.TerminalNode

func (s *DeferStatementContext) GetParser() antlr.Parser

func (s *DeferStatementContext) GetRuleContext() antlr.RuleContext

func (*DeferStatementContext) IsDeferStatementContext()

func (s *DeferStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DeferStmtContext struct {
	StatementContext
}

func NewDeferStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DeferStmtContext

func (s *DeferStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeferStmtContext) DeferStatement() IDeferStatementContext

func (s *DeferStmtContext) GetRuleContext() antlr.RuleContext

func (s *DeferStmtContext) SEMICOLON() antlr.TerminalNode

type DeinitializerDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDeinitializerDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeinitializerDeclarationContext

func NewEmptyDeinitializerDeclarationContext() *DeinitializerDeclarationContext

func (s *DeinitializerDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeinitializerDeclarationContext) Attributes() IAttributesContext

func (s *DeinitializerDeclarationContext) CodeBlock() ICodeBlockContext

func (s *DeinitializerDeclarationContext) DEINIT() antlr.TerminalNode

func (s *DeinitializerDeclarationContext) GetParser() antlr.Parser

func (s *DeinitializerDeclarationContext) GetRuleContext() antlr.RuleContext

func (*DeinitializerDeclarationContext) IsDeinitializerDeclarationContext()

func (s *DeinitializerDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DiagnosticStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDiagnosticStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DiagnosticStatementContext

func NewEmptyDiagnosticStatementContext() *DiagnosticStatementContext

func (s *DiagnosticStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DiagnosticStatementContext) GetParser() antlr.Parser

func (s *DiagnosticStatementContext) GetRuleContext() antlr.RuleContext

func (*DiagnosticStatementContext) IsDiagnosticStatementContext()

func (s *DiagnosticStatementContext) OPERATOR() antlr.TerminalNode

func (s *DiagnosticStatementContext) STRING_LITERAL() antlr.TerminalNode

func (s *DiagnosticStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DictTypeContext struct {
	TypeContext
}

func NewDictTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DictTypeContext

func (s *DictTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DictTypeContext) DictionaryType() IDictionaryTypeContext

func (s *DictTypeContext) GetRuleContext() antlr.RuleContext

type DictionaryLiteralContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDictionaryLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DictionaryLiteralContext

func NewEmptyDictionaryLiteralContext() *DictionaryLiteralContext

func (s *DictionaryLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DictionaryLiteralContext) COLON() antlr.TerminalNode

func (s *DictionaryLiteralContext) DictionaryLiteralItems() IDictionaryLiteralItemsContext

func (s *DictionaryLiteralContext) GetParser() antlr.Parser

func (s *DictionaryLiteralContext) GetRuleContext() antlr.RuleContext

func (*DictionaryLiteralContext) IsDictionaryLiteralContext()

func (s *DictionaryLiteralContext) LBRACKET() antlr.TerminalNode

func (s *DictionaryLiteralContext) RBRACKET() antlr.TerminalNode

func (s *DictionaryLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DictionaryLiteralItemContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDictionaryLiteralItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DictionaryLiteralItemContext

func NewEmptyDictionaryLiteralItemContext() *DictionaryLiteralItemContext

func (s *DictionaryLiteralItemContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DictionaryLiteralItemContext) AllExpression() []IExpressionContext

func (s *DictionaryLiteralItemContext) COLON() antlr.TerminalNode

func (s *DictionaryLiteralItemContext) Expression(i int) IExpressionContext

func (s *DictionaryLiteralItemContext) GetParser() antlr.Parser

func (s *DictionaryLiteralItemContext) GetRuleContext() antlr.RuleContext

func (*DictionaryLiteralItemContext) IsDictionaryLiteralItemContext()

func (s *DictionaryLiteralItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DictionaryLiteralItemsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDictionaryLiteralItemsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DictionaryLiteralItemsContext

func NewEmptyDictionaryLiteralItemsContext() *DictionaryLiteralItemsContext

func (s *DictionaryLiteralItemsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DictionaryLiteralItemsContext) AllCOMMA() []antlr.TerminalNode

func (s *DictionaryLiteralItemsContext) AllDictionaryLiteralItem() []IDictionaryLiteralItemContext

func (s *DictionaryLiteralItemsContext) COMMA(i int) antlr.TerminalNode

func (s *DictionaryLiteralItemsContext) DictionaryLiteralItem(i int) IDictionaryLiteralItemContext

func (s *DictionaryLiteralItemsContext) GetParser() antlr.Parser

func (s *DictionaryLiteralItemsContext) GetRuleContext() antlr.RuleContext

func (*DictionaryLiteralItemsContext) IsDictionaryLiteralItemsContext()

func (s *DictionaryLiteralItemsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DictionaryTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDictionaryTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DictionaryTypeContext

func NewEmptyDictionaryTypeContext() *DictionaryTypeContext

func (s *DictionaryTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DictionaryTypeContext) AllType_() []ITypeContext

func (s *DictionaryTypeContext) COLON() antlr.TerminalNode

func (s *DictionaryTypeContext) GetParser() antlr.Parser

func (s *DictionaryTypeContext) GetRuleContext() antlr.RuleContext

func (*DictionaryTypeContext) IsDictionaryTypeContext()

func (s *DictionaryTypeContext) LBRACKET() antlr.TerminalNode

func (s *DictionaryTypeContext) RBRACKET() antlr.TerminalNode

func (s *DictionaryTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *DictionaryTypeContext) Type_(i int) ITypeContext

type DidSetClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDidSetClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DidSetClauseContext

func NewEmptyDidSetClauseContext() *DidSetClauseContext

func (s *DidSetClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DidSetClauseContext) Attributes() IAttributesContext

func (s *DidSetClauseContext) CodeBlock() ICodeBlockContext

func (s *DidSetClauseContext) DIDSET_KW() antlr.TerminalNode

func (s *DidSetClauseContext) GetParser() antlr.Parser

func (s *DidSetClauseContext) GetRuleContext() antlr.RuleContext

func (*DidSetClauseContext) IsDidSetClauseContext()

func (s *DidSetClauseContext) SetterName() ISetterNameContext

func (s *DidSetClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DoStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDoStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DoStatementContext

func NewEmptyDoStatementContext() *DoStatementContext

func (s *DoStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DoStatementContext) AllCatchClause() []ICatchClauseContext

func (s *DoStatementContext) CatchClause(i int) ICatchClauseContext

func (s *DoStatementContext) CodeBlock() ICodeBlockContext

func (s *DoStatementContext) DO() antlr.TerminalNode

func (s *DoStatementContext) GetParser() antlr.Parser

func (s *DoStatementContext) GetRuleContext() antlr.RuleContext

func (*DoStatementContext) IsDoStatementContext()

func (s *DoStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DoStmtContext struct {
	StatementContext
}

func NewDoStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DoStmtContext

func (s *DoStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DoStmtContext) DoStatement() IDoStatementContext

func (s *DoStmtContext) GetRuleContext() antlr.RuleContext

func (s *DoStmtContext) SEMICOLON() antlr.TerminalNode

type ElseClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewElseClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElseClauseContext

func NewEmptyElseClauseContext() *ElseClauseContext

func (s *ElseClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ElseClauseContext) CodeBlock() ICodeBlockContext

func (s *ElseClauseContext) ELSE() antlr.TerminalNode

func (s *ElseClauseContext) GetParser() antlr.Parser

func (s *ElseClauseContext) GetRuleContext() antlr.RuleContext

func (s *ElseClauseContext) IfStatement() IIfStatementContext

func (*ElseClauseContext) IsElseClauseContext()

func (s *ElseClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type EnumCasePatContext struct {
	PatternContext
}

func NewEnumCasePatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EnumCasePatContext

func (s *EnumCasePatContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumCasePatContext) EnumCasePattern() IEnumCasePatternContext

func (s *EnumCasePatContext) GetRuleContext() antlr.RuleContext

type EnumCasePatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumCasePatternContext() *EnumCasePatternContext

func NewEnumCasePatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumCasePatternContext

func (s *EnumCasePatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumCasePatternContext) DOT() antlr.TerminalNode

func (s *EnumCasePatternContext) GetParser() antlr.Parser

func (s *EnumCasePatternContext) GetRuleContext() antlr.RuleContext

func (s *EnumCasePatternContext) Identifier() IIdentifierContext

func (*EnumCasePatternContext) IsEnumCasePatternContext()

func (s *EnumCasePatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *EnumCasePatternContext) TuplePattern() ITuplePatternContext

func (s *EnumCasePatternContext) TypeIdentifier() ITypeIdentifierContext

type EnumDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumDeclarationContext() *EnumDeclarationContext

func NewEnumDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumDeclarationContext

func (s *EnumDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumDeclarationContext) AccessLevelModifier() IAccessLevelModifierContext

func (s *EnumDeclarationContext) Attributes() IAttributesContext

func (s *EnumDeclarationContext) ENUM() antlr.TerminalNode

func (s *EnumDeclarationContext) EnumMembers() IEnumMembersContext

func (s *EnumDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *EnumDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *EnumDeclarationContext) GetParser() antlr.Parser

func (s *EnumDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *EnumDeclarationContext) Identifier() IIdentifierContext

func (*EnumDeclarationContext) IsEnumDeclarationContext()

func (s *EnumDeclarationContext) LBRACE() antlr.TerminalNode

func (s *EnumDeclarationContext) RBRACE() antlr.TerminalNode

func (s *EnumDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *EnumDeclarationContext) TypeInheritanceClause() ITypeInheritanceClauseContext

type EnumMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumMemberContext() *EnumMemberContext

func NewEnumMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumMemberContext

func (s *EnumMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumMemberContext) CompilerControl() ICompilerControlContext

func (s *EnumMemberContext) Declaration() IDeclarationContext

func (s *EnumMemberContext) GetParser() antlr.Parser

func (s *EnumMemberContext) GetRuleContext() antlr.RuleContext

func (*EnumMemberContext) IsEnumMemberContext()

func (s *EnumMemberContext) RawValueStyleEnumCaseClause() IRawValueStyleEnumCaseClauseContext

func (s *EnumMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *EnumMemberContext) UnionStyleEnumCaseClause() IUnionStyleEnumCaseClauseContext

type EnumMembersContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumMembersContext() *EnumMembersContext

func NewEnumMembersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumMembersContext

func (s *EnumMembersContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumMembersContext) AllEnumMember() []IEnumMemberContext

func (s *EnumMembersContext) EnumMember(i int) IEnumMemberContext

func (s *EnumMembersContext) GetParser() antlr.Parser

func (s *EnumMembersContext) GetRuleContext() antlr.RuleContext

func (*EnumMembersContext) IsEnumMembersContext()

func (s *EnumMembersContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExistTypeContext struct {
	TypeContext
}

func NewExistTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExistTypeContext

func (s *ExistTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExistTypeContext) ExistentialType() IExistentialTypeContext

func (s *ExistTypeContext) GetRuleContext() antlr.RuleContext

type ExistentialTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExistentialTypeContext() *ExistentialTypeContext

func NewExistentialTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExistentialTypeContext

func (s *ExistentialTypeContext) ANY() antlr.TerminalNode

func (s *ExistentialTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExistentialTypeContext) GetParser() antlr.Parser

func (s *ExistentialTypeContext) GetRuleContext() antlr.RuleContext

func (*ExistentialTypeContext) IsExistentialTypeContext()

func (s *ExistentialTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ExistentialTypeContext) Type_() ITypeContext

type ExplicitMemberSuffixContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExplicitMemberSuffixContext() *ExplicitMemberSuffixContext

func NewExplicitMemberSuffixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExplicitMemberSuffixContext

func (s *ExplicitMemberSuffixContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExplicitMemberSuffixContext) DOT() antlr.TerminalNode

func (s *ExplicitMemberSuffixContext) GenericArgumentClause() IGenericArgumentClauseContext

func (s *ExplicitMemberSuffixContext) GetParser() antlr.Parser

func (s *ExplicitMemberSuffixContext) GetRuleContext() antlr.RuleContext

func (s *ExplicitMemberSuffixContext) INTEGER_LITERAL() antlr.TerminalNode

func (s *ExplicitMemberSuffixContext) Identifier() IIdentifierContext

func (*ExplicitMemberSuffixContext) IsExplicitMemberSuffixContext()

func (s *ExplicitMemberSuffixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExprPatContext struct {
	PatternContext
}

func NewExprPatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprPatContext

func (s *ExprPatContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExprPatContext) ExpressionPattern() IExpressionPatternContext

func (s *ExprPatContext) GetRuleContext() antlr.RuleContext

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExpressionContext() *ExpressionContext

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext

func (s *ExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExpressionContext) AwaitOperator() IAwaitOperatorContext

func (s *ExpressionContext) BinaryExpressions() IBinaryExpressionsContext

func (s *ExpressionContext) GetParser() antlr.Parser

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext

func (*ExpressionContext) IsExpressionContext()

func (s *ExpressionContext) PrefixExpression() IPrefixExpressionContext

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ExpressionContext) TryOperator() ITryOperatorContext

type ExpressionPatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExpressionPatternContext() *ExpressionPatternContext

func NewExpressionPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionPatternContext

func (s *ExpressionPatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExpressionPatternContext) Expression() IExpressionContext

func (s *ExpressionPatternContext) GetParser() antlr.Parser

func (s *ExpressionPatternContext) GetRuleContext() antlr.RuleContext

func (*ExpressionPatternContext) IsExpressionPatternContext()

func (s *ExpressionPatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExpressionStatementContext struct {
	StatementContext
}

func NewExpressionStatementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExpressionStatementContext

func (s *ExpressionStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExpressionStatementContext) Expression() IExpressionContext

func (s *ExpressionStatementContext) GetRuleContext() antlr.RuleContext

func (s *ExpressionStatementContext) SEMICOLON() antlr.TerminalNode

type ExtensionBodyContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExtensionBodyContext() *ExtensionBodyContext

func NewExtensionBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExtensionBodyContext

func (s *ExtensionBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExtensionBodyContext) AllExtensionMember() []IExtensionMemberContext

func (s *ExtensionBodyContext) ExtensionMember(i int) IExtensionMemberContext

func (s *ExtensionBodyContext) GetParser() antlr.Parser

func (s *ExtensionBodyContext) GetRuleContext() antlr.RuleContext

func (*ExtensionBodyContext) IsExtensionBodyContext()

func (s *ExtensionBodyContext) LBRACE() antlr.TerminalNode

func (s *ExtensionBodyContext) RBRACE() antlr.TerminalNode

func (s *ExtensionBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExtensionDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExtensionDeclarationContext() *ExtensionDeclarationContext

func NewExtensionDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExtensionDeclarationContext

func (s *ExtensionDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExtensionDeclarationContext) AccessLevelModifier() IAccessLevelModifierContext

func (s *ExtensionDeclarationContext) Attributes() IAttributesContext

func (s *ExtensionDeclarationContext) COLON() antlr.TerminalNode

func (s *ExtensionDeclarationContext) EXTENSION() antlr.TerminalNode

func (s *ExtensionDeclarationContext) ExtensionBody() IExtensionBodyContext

func (s *ExtensionDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *ExtensionDeclarationContext) GetParser() antlr.Parser

func (s *ExtensionDeclarationContext) GetRuleContext() antlr.RuleContext

func (*ExtensionDeclarationContext) IsExtensionDeclarationContext()

func (s *ExtensionDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ExtensionDeclarationContext) TypeIdentifier() ITypeIdentifierContext

func (s *ExtensionDeclarationContext) TypeInheritanceList() ITypeInheritanceListContext

type ExtensionMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExtensionMemberContext() *ExtensionMemberContext

func NewExtensionMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExtensionMemberContext

func (s *ExtensionMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExtensionMemberContext) CompilerControl() ICompilerControlContext

func (s *ExtensionMemberContext) Declaration() IDeclarationContext

func (s *ExtensionMemberContext) GetParser() antlr.Parser

func (s *ExtensionMemberContext) GetRuleContext() antlr.RuleContext

func (*ExtensionMemberContext) IsExtensionMemberContext()

func (s *ExtensionMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternalParameterNameContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternalParameterNameContext() *ExternalParameterNameContext

func NewExternalParameterNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternalParameterNameContext

func (s *ExternalParameterNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternalParameterNameContext) GetParser() antlr.Parser

func (s *ExternalParameterNameContext) GetRuleContext() antlr.RuleContext

func (s *ExternalParameterNameContext) Identifier() IIdentifierContext

func (*ExternalParameterNameContext) IsExternalParameterNameContext()

func (s *ExternalParameterNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ExternalParameterNameContext) UNDERSCORE() antlr.TerminalNode

type FallthroughStatementContext struct {
	ControlTransferContext
}

func NewFallthroughStatementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FallthroughStatementContext

func (s *FallthroughStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FallthroughStatementContext) FALLTHROUGH() antlr.TerminalNode

func (s *FallthroughStatementContext) GetRuleContext() antlr.RuleContext

type ForInStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyForInStatementContext() *ForInStatementContext

func NewForInStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForInStatementContext

func (s *ForInStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ForInStatementContext) CASE() antlr.TerminalNode

func (s *ForInStatementContext) CodeBlock() ICodeBlockContext

func (s *ForInStatementContext) Expression() IExpressionContext

func (s *ForInStatementContext) FOR() antlr.TerminalNode

func (s *ForInStatementContext) GetParser() antlr.Parser

func (s *ForInStatementContext) GetRuleContext() antlr.RuleContext

func (s *ForInStatementContext) IN() antlr.TerminalNode

func (*ForInStatementContext) IsForInStatementContext()

func (s *ForInStatementContext) Pattern() IPatternContext

func (s *ForInStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ForInStatementContext) WhereClause() IWhereClauseContext

type ForStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyForStatementContext() *ForStatementContext

func NewForStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForStatementContext

func (s *ForStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ForStatementContext) AllExpression() []IExpressionContext

func (s *ForStatementContext) AllSEMICOLON() []antlr.TerminalNode

func (s *ForStatementContext) CodeBlock() ICodeBlockContext

func (s *ForStatementContext) Expression(i int) IExpressionContext

func (s *ForStatementContext) FOR() antlr.TerminalNode

func (s *ForStatementContext) GetParser() antlr.Parser

func (s *ForStatementContext) GetRuleContext() antlr.RuleContext

func (*ForStatementContext) IsForStatementContext()

func (s *ForStatementContext) LPAREN() antlr.TerminalNode

func (s *ForStatementContext) RPAREN() antlr.TerminalNode

func (s *ForStatementContext) SEMICOLON(i int) antlr.TerminalNode

func (s *ForStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ForcedValueSuffixContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyForcedValueSuffixContext() *ForcedValueSuffixContext

func NewForcedValueSuffixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForcedValueSuffixContext

func (s *ForcedValueSuffixContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ForcedValueSuffixContext) EXCLAIM_POSTFIX() antlr.TerminalNode

func (s *ForcedValueSuffixContext) GetParser() antlr.Parser

func (s *ForcedValueSuffixContext) GetRuleContext() antlr.RuleContext

func (*ForcedValueSuffixContext) IsForcedValueSuffixContext()

func (s *ForcedValueSuffixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FuncTypeContext struct {
	TypeContext
}

func NewFuncTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FuncTypeContext

func (s *FuncTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FuncTypeContext) FunctionType() IFunctionTypeContext

func (s *FuncTypeContext) GetRuleContext() antlr.RuleContext

type FunctionBodyContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionBodyContext() *FunctionBodyContext

func NewFunctionBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionBodyContext

func (s *FunctionBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionBodyContext) CodeBlock() ICodeBlockContext

func (s *FunctionBodyContext) GetParser() antlr.Parser

func (s *FunctionBodyContext) GetRuleContext() antlr.RuleContext

func (*FunctionBodyContext) IsFunctionBodyContext()

func (s *FunctionBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FunctionCallArgumentClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionCallArgumentClauseContext() *FunctionCallArgumentClauseContext

func NewFunctionCallArgumentClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallArgumentClauseContext

func (s *FunctionCallArgumentClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionCallArgumentClauseContext) FunctionCallArgumentList() IFunctionCallArgumentListContext

func (s *FunctionCallArgumentClauseContext) GetParser() antlr.Parser

func (s *FunctionCallArgumentClauseContext) GetRuleContext() antlr.RuleContext

func (*FunctionCallArgumentClauseContext) IsFunctionCallArgumentClauseContext()

func (s *FunctionCallArgumentClauseContext) LPAREN() antlr.TerminalNode

func (s *FunctionCallArgumentClauseContext) RPAREN() antlr.TerminalNode

func (s *FunctionCallArgumentClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FunctionCallArgumentContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionCallArgumentContext() *FunctionCallArgumentContext

func NewFunctionCallArgumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallArgumentContext

func (s *FunctionCallArgumentContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionCallArgumentContext) COLON() antlr.TerminalNode

func (s *FunctionCallArgumentContext) Expression() IExpressionContext

func (s *FunctionCallArgumentContext) GetParser() antlr.Parser

func (s *FunctionCallArgumentContext) GetRuleContext() antlr.RuleContext

func (s *FunctionCallArgumentContext) Identifier() IIdentifierContext

func (*FunctionCallArgumentContext) IsFunctionCallArgumentContext()

func (s *FunctionCallArgumentContext) Operator_() IOperator_Context

func (s *FunctionCallArgumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FunctionCallArgumentContext) UNDERSCORE() antlr.TerminalNode

type FunctionCallArgumentListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionCallArgumentListContext() *FunctionCallArgumentListContext

func NewFunctionCallArgumentListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallArgumentListContext

func (s *FunctionCallArgumentListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionCallArgumentListContext) AllCOMMA() []antlr.TerminalNode

func (s *FunctionCallArgumentListContext) AllFunctionCallArgument() []IFunctionCallArgumentContext

func (s *FunctionCallArgumentListContext) COMMA(i int) antlr.TerminalNode

func (s *FunctionCallArgumentListContext) FunctionCallArgument(i int) IFunctionCallArgumentContext

func (s *FunctionCallArgumentListContext) GetParser() antlr.Parser

func (s *FunctionCallArgumentListContext) GetRuleContext() antlr.RuleContext

func (*FunctionCallArgumentListContext) IsFunctionCallArgumentListContext()

func (s *FunctionCallArgumentListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FunctionCallSuffixContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionCallSuffixContext() *FunctionCallSuffixContext

func NewFunctionCallSuffixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallSuffixContext

func (s *FunctionCallSuffixContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionCallSuffixContext) FunctionCallArgumentClause() IFunctionCallArgumentClauseContext

func (s *FunctionCallSuffixContext) GetParser() antlr.Parser

func (s *FunctionCallSuffixContext) GetRuleContext() antlr.RuleContext

func (*FunctionCallSuffixContext) IsFunctionCallSuffixContext()

func (s *FunctionCallSuffixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FunctionCallSuffixContext) TrailingClosures() ITrailingClosuresContext

type FunctionDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionDeclarationContext() *FunctionDeclarationContext

func NewFunctionDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionDeclarationContext

func (s *FunctionDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionDeclarationContext) FunctionBody() IFunctionBodyContext

func (s *FunctionDeclarationContext) FunctionHead() IFunctionHeadContext

func (s *FunctionDeclarationContext) FunctionSignature() IFunctionSignatureContext

func (s *FunctionDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *FunctionDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *FunctionDeclarationContext) GetParser() antlr.Parser

func (s *FunctionDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *FunctionDeclarationContext) Identifier() IIdentifierContext

func (*FunctionDeclarationContext) IsFunctionDeclarationContext()

func (s *FunctionDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FunctionHeadContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionHeadContext() *FunctionHeadContext

func NewFunctionHeadContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionHeadContext

func (s *FunctionHeadContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionHeadContext) Attributes() IAttributesContext

func (s *FunctionHeadContext) DeclarationModifiers() IDeclarationModifiersContext

func (s *FunctionHeadContext) FUNC() antlr.TerminalNode

func (s *FunctionHeadContext) GetParser() antlr.Parser

func (s *FunctionHeadContext) GetRuleContext() antlr.RuleContext

func (*FunctionHeadContext) IsFunctionHeadContext()

func (s *FunctionHeadContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FunctionResultContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionResultContext() *FunctionResultContext

func NewFunctionResultContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionResultContext

func (s *FunctionResultContext) ARROW() antlr.TerminalNode

func (s *FunctionResultContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionResultContext) Attributes() IAttributesContext

func (s *FunctionResultContext) GetParser() antlr.Parser

func (s *FunctionResultContext) GetRuleContext() antlr.RuleContext

func (*FunctionResultContext) IsFunctionResultContext()

func (s *FunctionResultContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FunctionResultContext) Type_() ITypeContext

type FunctionSignatureContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionSignatureContext() *FunctionSignatureContext

func NewFunctionSignatureContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionSignatureContext

func (s *FunctionSignatureContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionSignatureContext) AsyncModifier() IAsyncModifierContext

func (s *FunctionSignatureContext) FunctionResult() IFunctionResultContext

func (s *FunctionSignatureContext) GetParser() antlr.Parser

func (s *FunctionSignatureContext) GetRuleContext() antlr.RuleContext

func (*FunctionSignatureContext) IsFunctionSignatureContext()

func (s *FunctionSignatureContext) ParameterClause() IParameterClauseContext

func (s *FunctionSignatureContext) THROWS() antlr.TerminalNode

func (s *FunctionSignatureContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FunctionTypeArgumentClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionTypeArgumentClauseContext() *FunctionTypeArgumentClauseContext

func NewFunctionTypeArgumentClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionTypeArgumentClauseContext

func (s *FunctionTypeArgumentClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionTypeArgumentClauseContext) ELLIPSIS() antlr.TerminalNode

func (s *FunctionTypeArgumentClauseContext) FunctionTypeArgumentList() IFunctionTypeArgumentListContext

func (s *FunctionTypeArgumentClauseContext) GetParser() antlr.Parser

func (s *FunctionTypeArgumentClauseContext) GetRuleContext() antlr.RuleContext

func (*FunctionTypeArgumentClauseContext) IsFunctionTypeArgumentClauseContext()

func (s *FunctionTypeArgumentClauseContext) LPAREN() antlr.TerminalNode

func (s *FunctionTypeArgumentClauseContext) RPAREN() antlr.TerminalNode

func (s *FunctionTypeArgumentClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FunctionTypeArgumentContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionTypeArgumentContext() *FunctionTypeArgumentContext

func NewFunctionTypeArgumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionTypeArgumentContext

func (s *FunctionTypeArgumentContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionTypeArgumentContext) ArgumentLabel() IArgumentLabelContext

func (s *FunctionTypeArgumentContext) Attributes() IAttributesContext

func (s *FunctionTypeArgumentContext) GetParser() antlr.Parser

func (s *FunctionTypeArgumentContext) GetRuleContext() antlr.RuleContext

func (s *FunctionTypeArgumentContext) INOUT() antlr.TerminalNode

func (*FunctionTypeArgumentContext) IsFunctionTypeArgumentContext()

func (s *FunctionTypeArgumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FunctionTypeArgumentContext) TypeAnnotation() ITypeAnnotationContext

func (s *FunctionTypeArgumentContext) Type_() ITypeContext

type FunctionTypeArgumentListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionTypeArgumentListContext() *FunctionTypeArgumentListContext

func NewFunctionTypeArgumentListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionTypeArgumentListContext

func (s *FunctionTypeArgumentListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionTypeArgumentListContext) AllCOMMA() []antlr.TerminalNode

func (s *FunctionTypeArgumentListContext) AllFunctionTypeArgument() []IFunctionTypeArgumentContext

func (s *FunctionTypeArgumentListContext) COMMA(i int) antlr.TerminalNode

func (s *FunctionTypeArgumentListContext) FunctionTypeArgument(i int) IFunctionTypeArgumentContext

func (s *FunctionTypeArgumentListContext) GetParser() antlr.Parser

func (s *FunctionTypeArgumentListContext) GetRuleContext() antlr.RuleContext

func (*FunctionTypeArgumentListContext) IsFunctionTypeArgumentListContext()

func (s *FunctionTypeArgumentListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FunctionTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionTypeContext() *FunctionTypeContext

func NewFunctionTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionTypeContext

func (s *FunctionTypeContext) ARROW() antlr.TerminalNode

func (s *FunctionTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionTypeContext) AsyncModifier() IAsyncModifierContext

func (s *FunctionTypeContext) Attributes() IAttributesContext

func (s *FunctionTypeContext) FunctionTypeArgumentClause() IFunctionTypeArgumentClauseContext

func (s *FunctionTypeContext) GetParser() antlr.Parser

func (s *FunctionTypeContext) GetRuleContext() antlr.RuleContext

func (*FunctionTypeContext) IsFunctionTypeContext()

func (s *FunctionTypeContext) THROWS() antlr.TerminalNode

func (s *FunctionTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FunctionTypeContext) Type_() ITypeContext

type GenericArgumentClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericArgumentClauseContext() *GenericArgumentClauseContext

func NewGenericArgumentClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericArgumentClauseContext

func (s *GenericArgumentClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericArgumentClauseContext) GT() antlr.TerminalNode

func (s *GenericArgumentClauseContext) GenericArgumentList() IGenericArgumentListContext

func (s *GenericArgumentClauseContext) GetParser() antlr.Parser

func (s *GenericArgumentClauseContext) GetRuleContext() antlr.RuleContext

func (*GenericArgumentClauseContext) IsGenericArgumentClauseContext()

func (s *GenericArgumentClauseContext) LT() antlr.TerminalNode

func (s *GenericArgumentClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type GenericArgumentContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericArgumentContext() *GenericArgumentContext

func NewGenericArgumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericArgumentContext

func (s *GenericArgumentContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericArgumentContext) GetParser() antlr.Parser

func (s *GenericArgumentContext) GetRuleContext() antlr.RuleContext

func (*GenericArgumentContext) IsGenericArgumentContext()

func (s *GenericArgumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *GenericArgumentContext) Type_() ITypeContext

type GenericArgumentListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericArgumentListContext() *GenericArgumentListContext

func NewGenericArgumentListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericArgumentListContext

func (s *GenericArgumentListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericArgumentListContext) AllCOMMA() []antlr.TerminalNode

func (s *GenericArgumentListContext) AllGenericArgument() []IGenericArgumentContext

func (s *GenericArgumentListContext) COMMA(i int) antlr.TerminalNode

func (s *GenericArgumentListContext) GenericArgument(i int) IGenericArgumentContext

func (s *GenericArgumentListContext) GetParser() antlr.Parser

func (s *GenericArgumentListContext) GetRuleContext() antlr.RuleContext

func (*GenericArgumentListContext) IsGenericArgumentListContext()

func (s *GenericArgumentListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type GenericParameterClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericParameterClauseContext() *GenericParameterClauseContext

func NewGenericParameterClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericParameterClauseContext

func (s *GenericParameterClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericParameterClauseContext) GT() antlr.TerminalNode

func (s *GenericParameterClauseContext) GenericParameterList() IGenericParameterListContext

func (s *GenericParameterClauseContext) GetParser() antlr.Parser

func (s *GenericParameterClauseContext) GetRuleContext() antlr.RuleContext

func (*GenericParameterClauseContext) IsGenericParameterClauseContext()

func (s *GenericParameterClauseContext) LT() antlr.TerminalNode

func (s *GenericParameterClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type GenericParameterContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericParameterContext() *GenericParameterContext

func NewGenericParameterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericParameterContext

func (s *GenericParameterContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericParameterContext) COLON() antlr.TerminalNode

func (s *GenericParameterContext) GetParser() antlr.Parser

func (s *GenericParameterContext) GetRuleContext() antlr.RuleContext

func (s *GenericParameterContext) Identifier() IIdentifierContext

func (*GenericParameterContext) IsGenericParameterContext()

func (s *GenericParameterContext) ProtocolCompositionType() IProtocolCompositionTypeContext

func (s *GenericParameterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *GenericParameterContext) TypeIdentifier() ITypeIdentifierContext

type GenericParameterListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericParameterListContext() *GenericParameterListContext

func NewGenericParameterListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericParameterListContext

func (s *GenericParameterListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericParameterListContext) AllCOMMA() []antlr.TerminalNode

func (s *GenericParameterListContext) AllGenericParameter() []IGenericParameterContext

func (s *GenericParameterListContext) COMMA(i int) antlr.TerminalNode

func (s *GenericParameterListContext) GenericParameter(i int) IGenericParameterContext

func (s *GenericParameterListContext) GetParser() antlr.Parser

func (s *GenericParameterListContext) GetRuleContext() antlr.RuleContext

func (*GenericParameterListContext) IsGenericParameterListContext()

func (s *GenericParameterListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type GenericWhereClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericWhereClauseContext() *GenericWhereClauseContext

func NewGenericWhereClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericWhereClauseContext

func (s *GenericWhereClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericWhereClauseContext) GetParser() antlr.Parser

func (s *GenericWhereClauseContext) GetRuleContext() antlr.RuleContext

func (*GenericWhereClauseContext) IsGenericWhereClauseContext()

func (s *GenericWhereClauseContext) RequirementList() IRequirementListContext

func (s *GenericWhereClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *GenericWhereClauseContext) WHERE() antlr.TerminalNode

type GetterClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGetterClauseContext() *GetterClauseContext

func NewGetterClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GetterClauseContext

func (s *GetterClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GetterClauseContext) Attributes() IAttributesContext

func (s *GetterClauseContext) CodeBlock() ICodeBlockContext

func (s *GetterClauseContext) GET_KW() antlr.TerminalNode

func (s *GetterClauseContext) GetParser() antlr.Parser

func (s *GetterClauseContext) GetRuleContext() antlr.RuleContext

func (*GetterClauseContext) IsGetterClauseContext()

func (s *GetterClauseContext) MutationModifier() IMutationModifierContext

func (s *GetterClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type GetterKeywordClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGetterKeywordClauseContext() *GetterKeywordClauseContext

func NewGetterKeywordClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GetterKeywordClauseContext

func (s *GetterKeywordClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GetterKeywordClauseContext) Attributes() IAttributesContext

func (s *GetterKeywordClauseContext) GET_KW() antlr.TerminalNode

func (s *GetterKeywordClauseContext) GetParser() antlr.Parser

func (s *GetterKeywordClauseContext) GetRuleContext() antlr.RuleContext

func (*GetterKeywordClauseContext) IsGetterKeywordClauseContext()

func (s *GetterKeywordClauseContext) MutationModifier() IMutationModifierContext

func (s *GetterKeywordClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type GetterSetterBlockContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGetterSetterBlockContext() *GetterSetterBlockContext

func NewGetterSetterBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GetterSetterBlockContext

func (s *GetterSetterBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GetterSetterBlockContext) GetParser() antlr.Parser

func (s *GetterSetterBlockContext) GetRuleContext() antlr.RuleContext

func (s *GetterSetterBlockContext) GetterClause() IGetterClauseContext

func (*GetterSetterBlockContext) IsGetterSetterBlockContext()

func (s *GetterSetterBlockContext) LBRACE() antlr.TerminalNode

func (s *GetterSetterBlockContext) RBRACE() antlr.TerminalNode

func (s *GetterSetterBlockContext) SetterClause() ISetterClauseContext

func (s *GetterSetterBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type GetterSetterKeywordBlockContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGetterSetterKeywordBlockContext() *GetterSetterKeywordBlockContext

func NewGetterSetterKeywordBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GetterSetterKeywordBlockContext

func (s *GetterSetterKeywordBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GetterSetterKeywordBlockContext) GetParser() antlr.Parser

func (s *GetterSetterKeywordBlockContext) GetRuleContext() antlr.RuleContext

func (s *GetterSetterKeywordBlockContext) GetterKeywordClause() IGetterKeywordClauseContext

func (*GetterSetterKeywordBlockContext) IsGetterSetterKeywordBlockContext()

func (s *GetterSetterKeywordBlockContext) LBRACE() antlr.TerminalNode

func (s *GetterSetterKeywordBlockContext) RBRACE() antlr.TerminalNode

func (s *GetterSetterKeywordBlockContext) SetterKeywordClause() ISetterKeywordClauseContext

func (s *GetterSetterKeywordBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type GuardStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGuardStatementContext() *GuardStatementContext

func NewGuardStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GuardStatementContext

func (s *GuardStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GuardStatementContext) CodeBlock() ICodeBlockContext

func (s *GuardStatementContext) ConditionList() IConditionListContext

func (s *GuardStatementContext) ELSE() antlr.TerminalNode

func (s *GuardStatementContext) GUARD() antlr.TerminalNode

func (s *GuardStatementContext) GetParser() antlr.Parser

func (s *GuardStatementContext) GetRuleContext() antlr.RuleContext

func (*GuardStatementContext) IsGuardStatementContext()

func (s *GuardStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type IAccessLevelModifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PRIVATE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	SET_KW() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	INTERNAL() antlr.TerminalNode
	PUBLIC() antlr.TerminalNode
	OPEN() antlr.TerminalNode

	// IsAccessLevelModifierContext differentiates from other interfaces.
	IsAccessLevelModifierContext()
}
    IAccessLevelModifierContext is an interface to support dynamic dispatch.

type IActorBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllActorMember() []IActorMemberContext
	ActorMember(i int) IActorMemberContext

	// IsActorBodyContext differentiates from other interfaces.
	IsActorBodyContext()
}
    IActorBodyContext is an interface to support dynamic dispatch.

type IActorDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ACTOR_KW() antlr.TerminalNode
	Identifier() IIdentifierContext
	ActorBody() IActorBodyContext
	Attributes() IAttributesContext
	AccessLevelModifier() IAccessLevelModifierContext
	GenericParameterClause() IGenericParameterClauseContext
	TypeInheritanceClause() ITypeInheritanceClauseContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsActorDeclarationContext differentiates from other interfaces.
	IsActorDeclarationContext()
}
    IActorDeclarationContext is an interface to support dynamic dispatch.

type IActorMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Declaration() IDeclarationContext
	CompilerControl() ICompilerControlContext

	// IsActorMemberContext differentiates from other interfaces.
	IsActorMemberContext()
}
    IActorMemberContext is an interface to support dynamic dispatch.

type IArgumentLabelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	COLON() antlr.TerminalNode
	UNDERSCORE() antlr.TerminalNode

	// IsArgumentLabelContext differentiates from other interfaces.
	IsArgumentLabelContext()
}
    IArgumentLabelContext is an interface to support dynamic dispatch.

type IArgumentNamesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIdentifier() []IIdentifierContext
	Identifier(i int) IIdentifierContext
	AllCOLON() []antlr.TerminalNode
	COLON(i int) antlr.TerminalNode

	// IsArgumentNamesContext differentiates from other interfaces.
	IsArgumentNamesContext()
}
    IArgumentNamesContext is an interface to support dynamic dispatch.

type IArrayLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	RBRACKET() antlr.TerminalNode
	ArrayLiteralItems() IArrayLiteralItemsContext

	// IsArrayLiteralContext differentiates from other interfaces.
	IsArrayLiteralContext()
}
    IArrayLiteralContext is an interface to support dynamic dispatch.

type IArrayLiteralItemsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsArrayLiteralItemsContext differentiates from other interfaces.
	IsArrayLiteralItemsContext()
}
    IArrayLiteralItemsContext is an interface to support dynamic dispatch.

type IArrayTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	Type_() ITypeContext
	RBRACKET() antlr.TerminalNode

	// IsArrayTypeContext differentiates from other interfaces.
	IsArrayTypeContext()
}
    IArrayTypeContext is an interface to support dynamic dispatch.

type IAssignmentOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASSIGN() antlr.TerminalNode

	// IsAssignmentOperatorContext differentiates from other interfaces.
	IsAssignmentOperatorContext()
}
    IAssignmentOperatorContext is an interface to support dynamic dispatch.

type IAsyncModifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASYNC_KW() antlr.TerminalNode

	// IsAsyncModifierContext differentiates from other interfaces.
	IsAsyncModifierContext()
}
    IAsyncModifierContext is an interface to support dynamic dispatch.

type IAttributeArgumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	Identifier() IIdentifierContext
	COLON() antlr.TerminalNode
	Type_() ITypeContext

	// IsAttributeArgumentContext differentiates from other interfaces.
	IsAttributeArgumentContext()
}
    IAttributeArgumentContext is an interface to support dynamic dispatch.

type IAttributeArgumentListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAttributeArgument() []IAttributeArgumentContext
	AttributeArgument(i int) IAttributeArgumentContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsAttributeArgumentListContext differentiates from other interfaces.
	IsAttributeArgumentListContext()
}
    IAttributeArgumentListContext is an interface to support dynamic dispatch.

type IAttributeArgumentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AttributeArgumentList() IAttributeArgumentListContext

	// IsAttributeArgumentsContext differentiates from other interfaces.
	IsAttributeArgumentsContext()
}
    IAttributeArgumentsContext is an interface to support dynamic dispatch.

type IAttributeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AT() antlr.TerminalNode
	Identifier() IIdentifierContext
	AttributeArguments() IAttributeArgumentsContext

	// IsAttributeContext differentiates from other interfaces.
	IsAttributeContext()
}
    IAttributeContext is an interface to support dynamic dispatch.

type IAttributesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAttribute() []IAttributeContext
	Attribute(i int) IAttributeContext

	// IsAttributesContext differentiates from other interfaces.
	IsAttributesContext()
}
    IAttributesContext is an interface to support dynamic dispatch.

type IAvailabilityArgumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	DecimalVersion() IDecimalVersionContext
	OPERATOR() antlr.TerminalNode

	// IsAvailabilityArgumentContext differentiates from other interfaces.
	IsAvailabilityArgumentContext()
}
    IAvailabilityArgumentContext is an interface to support dynamic dispatch.

type IAvailabilityArgumentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAvailabilityArgument() []IAvailabilityArgumentContext
	AvailabilityArgument(i int) IAvailabilityArgumentContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsAvailabilityArgumentsContext differentiates from other interfaces.
	IsAvailabilityArgumentsContext()
}
    IAvailabilityArgumentsContext is an interface to support dynamic dispatch.

type IAvailabilityConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	POUND_AVAILABLE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AvailabilityArguments() IAvailabilityArgumentsContext
	RPAREN() antlr.TerminalNode
	POUND_UNAVAILABLE() antlr.TerminalNode

	// IsAvailabilityConditionContext differentiates from other interfaces.
	IsAvailabilityConditionContext()
}
    IAvailabilityConditionContext is an interface to support dynamic dispatch.

type IAwaitOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AWAIT() antlr.TerminalNode

	// IsAwaitOperatorContext differentiates from other interfaces.
	IsAwaitOperatorContext()
}
    IAwaitOperatorContext is an interface to support dynamic dispatch.

type IBinaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsBinaryExpressionContext differentiates from other interfaces.
	IsBinaryExpressionContext()
}
    IBinaryExpressionContext is an interface to support dynamic dispatch.

type IBinaryExpressionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllBinaryExpression() []IBinaryExpressionContext
	BinaryExpression(i int) IBinaryExpressionContext

	// IsBinaryExpressionsContext differentiates from other interfaces.
	IsBinaryExpressionsContext()
}
    IBinaryExpressionsContext is an interface to support dynamic dispatch.

type IBinaryOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATOR() antlr.TerminalNode
	AMPERSAND() antlr.TerminalNode
	DOT() antlr.TerminalNode
	LT() antlr.TerminalNode
	GT() antlr.TerminalNode

	// IsBinaryOperatorContext differentiates from other interfaces.
	IsBinaryOperatorContext()
}
    IBinaryOperatorContext is an interface to support dynamic dispatch.

type IBooleanLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode

	// IsBooleanLiteralContext differentiates from other interfaces.
	IsBooleanLiteralContext()
}
    IBooleanLiteralContext is an interface to support dynamic dispatch.

type IBranchStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IfStatement() IIfStatementContext
	GuardStatement() IGuardStatementContext
	SwitchStatement() ISwitchStatementContext
	TryStatement() ITryStatementContext

	// IsBranchStatementContext differentiates from other interfaces.
	IsBranchStatementContext()
}
    IBranchStatementContext is an interface to support dynamic dispatch.

type ICaptureListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	CaptureListItems() ICaptureListItemsContext
	RBRACKET() antlr.TerminalNode

	// IsCaptureListContext differentiates from other interfaces.
	IsCaptureListContext()
}
    ICaptureListContext is an interface to support dynamic dispatch.

type ICaptureListItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	CaptureSpecifier() ICaptureSpecifierContext
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext

	// IsCaptureListItemContext differentiates from other interfaces.
	IsCaptureListItemContext()
}
    ICaptureListItemContext is an interface to support dynamic dispatch.

type ICaptureListItemsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCaptureListItem() []ICaptureListItemContext
	CaptureListItem(i int) ICaptureListItemContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsCaptureListItemsContext differentiates from other interfaces.
	IsCaptureListItemsContext()
}
    ICaptureListItemsContext is an interface to support dynamic dispatch.

type ICaptureSpecifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WEAK_KW() antlr.TerminalNode
	UNOWNED_KW() antlr.TerminalNode

	// IsCaptureSpecifierContext differentiates from other interfaces.
	IsCaptureSpecifierContext()
}
    ICaptureSpecifierContext is an interface to support dynamic dispatch.

type ICaseConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	Pattern() IPatternContext
	Initializer() IInitializerContext

	// IsCaseConditionContext differentiates from other interfaces.
	IsCaseConditionContext()
}
    ICaseConditionContext is an interface to support dynamic dispatch.

type ICaseItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Pattern() IPatternContext
	WhereClause() IWhereClauseContext

	// IsCaseItemContext differentiates from other interfaces.
	IsCaseItemContext()
}
    ICaseItemContext is an interface to support dynamic dispatch.

type ICaseItemListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCaseItem() []ICaseItemContext
	CaseItem(i int) ICaseItemContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsCaseItemListContext differentiates from other interfaces.
	IsCaseItemListContext()
}
    ICaseItemListContext is an interface to support dynamic dispatch.

type ICaseLabelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	CaseItemList() ICaseItemListContext
	COLON() antlr.TerminalNode

	// IsCaseLabelContext differentiates from other interfaces.
	IsCaseLabelContext()
}
    ICaseLabelContext is an interface to support dynamic dispatch.

type ICatchClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CATCH() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	CatchPatternList() ICatchPatternListContext

	// IsCatchClauseContext differentiates from other interfaces.
	IsCatchClauseContext()
}
    ICatchClauseContext is an interface to support dynamic dispatch.

type ICatchPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Pattern() IPatternContext
	WhereClause() IWhereClauseContext

	// IsCatchPatternContext differentiates from other interfaces.
	IsCatchPatternContext()
}
    ICatchPatternContext is an interface to support dynamic dispatch.

type ICatchPatternListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCatchPattern() []ICatchPatternContext
	CatchPattern(i int) ICatchPatternContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsCatchPatternListContext differentiates from other interfaces.
	IsCatchPatternListContext()
}
    ICatchPatternListContext is an interface to support dynamic dispatch.

type IClassBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllClassMember() []IClassMemberContext
	ClassMember(i int) IClassMemberContext

	// IsClassBodyContext differentiates from other interfaces.
	IsClassBodyContext()
}
    IClassBodyContext is an interface to support dynamic dispatch.

type IClassDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CLASS() antlr.TerminalNode
	Identifier() IIdentifierContext
	ClassBody() IClassBodyContext
	Attributes() IAttributesContext
	AccessLevelModifier() IAccessLevelModifierContext
	FINAL_KW() antlr.TerminalNode
	GenericParameterClause() IGenericParameterClauseContext
	TypeInheritanceClause() ITypeInheritanceClauseContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsClassDeclarationContext differentiates from other interfaces.
	IsClassDeclarationContext()
}
    IClassDeclarationContext is an interface to support dynamic dispatch.

type IClassMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Declaration() IDeclarationContext
	CompilerControl() ICompilerControlContext

	// IsClassMemberContext differentiates from other interfaces.
	IsClassMemberContext()
}
    IClassMemberContext is an interface to support dynamic dispatch.

type IClosureExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	ClosureSignature() IClosureSignatureContext
	Statements() IStatementsContext

	// IsClosureExpressionContext differentiates from other interfaces.
	IsClosureExpressionContext()
}
    IClosureExpressionContext is an interface to support dynamic dispatch.

type IClosureParameterClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ClosureParameterList() IClosureParameterListContext
	IdentifierList() IIdentifierListContext

	// IsClosureParameterClauseContext differentiates from other interfaces.
	IsClosureParameterClauseContext()
}
    IClosureParameterClauseContext is an interface to support dynamic dispatch.

type IClosureParameterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	TypeAnnotation() ITypeAnnotationContext
	ELLIPSIS() antlr.TerminalNode

	// IsClosureParameterContext differentiates from other interfaces.
	IsClosureParameterContext()
}
    IClosureParameterContext is an interface to support dynamic dispatch.

type IClosureParameterListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllClosureParameter() []IClosureParameterContext
	ClosureParameter(i int) IClosureParameterContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsClosureParameterListContext differentiates from other interfaces.
	IsClosureParameterListContext()
}
    IClosureParameterListContext is an interface to support dynamic dispatch.

type IClosureSignatureContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ClosureParameterClause() IClosureParameterClauseContext
	IN() antlr.TerminalNode
	CaptureList() ICaptureListContext
	ASYNC_KW() antlr.TerminalNode
	THROWS() antlr.TerminalNode
	FunctionResult() IFunctionResultContext

	// IsClosureSignatureContext differentiates from other interfaces.
	IsClosureSignatureContext()
}
    IClosureSignatureContext is an interface to support dynamic dispatch.

type ICodeBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	Statements() IStatementsContext

	// IsCodeBlockContext differentiates from other interfaces.
	IsCodeBlockContext()
}
    ICodeBlockContext is an interface to support dynamic dispatch.

type ICompilationConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PlatformCondition() IPlatformConditionContext
	Identifier() IIdentifierContext
	BooleanLiteral() IBooleanLiteralContext
	LPAREN() antlr.TerminalNode
	AllCompilationCondition() []ICompilationConditionContext
	CompilationCondition(i int) ICompilationConditionContext
	RPAREN() antlr.TerminalNode
	OPERATOR() antlr.TerminalNode

	// IsCompilationConditionContext differentiates from other interfaces.
	IsCompilationConditionContext()
}
    ICompilationConditionContext is an interface to support dynamic dispatch.

type ICompilerControlContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ConditionalCompilationBlock() IConditionalCompilationBlockContext
	LineControlStatement() ILineControlStatementContext
	DiagnosticStatement() IDiagnosticStatementContext

	// IsCompilerControlContext differentiates from other interfaces.
	IsCompilerControlContext()
}
    ICompilerControlContext is an interface to support dynamic dispatch.

type IConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	AvailabilityCondition() IAvailabilityConditionContext
	CaseCondition() ICaseConditionContext
	OptionalBindingCondition() IOptionalBindingConditionContext

	// IsConditionContext differentiates from other interfaces.
	IsConditionContext()
}
    IConditionContext is an interface to support dynamic dispatch.

type IConditionListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCondition() []IConditionContext
	Condition(i int) IConditionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsConditionListContext differentiates from other interfaces.
	IsConditionListContext()
}
    IConditionListContext is an interface to support dynamic dispatch.

type IConditionalCompilationBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	POUND_IF() antlr.TerminalNode
	AllCompilationCondition() []ICompilationConditionContext
	CompilationCondition(i int) ICompilationConditionContext
	POUND_ENDIF() antlr.TerminalNode
	AllStatements() []IStatementsContext
	Statements(i int) IStatementsContext
	AllPOUND_ELSEIF() []antlr.TerminalNode
	POUND_ELSEIF(i int) antlr.TerminalNode
	POUND_ELSE() antlr.TerminalNode

	// IsConditionalCompilationBlockContext differentiates from other interfaces.
	IsConditionalCompilationBlockContext()
}
    IConditionalCompilationBlockContext is an interface to support dynamic
    dispatch.

type IConditionalOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QUESTION_POSTFIX() antlr.TerminalNode
	Expression() IExpressionContext
	COLON() antlr.TerminalNode

	// IsConditionalOperatorContext differentiates from other interfaces.
	IsConditionalOperatorContext()
}
    IConditionalOperatorContext is an interface to support dynamic dispatch.

type IConformanceRequirementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTypeIdentifier() []ITypeIdentifierContext
	TypeIdentifier(i int) ITypeIdentifierContext
	COLON() antlr.TerminalNode
	ProtocolCompositionType() IProtocolCompositionTypeContext

	// IsConformanceRequirementContext differentiates from other interfaces.
	IsConformanceRequirementContext()
}
    IConformanceRequirementContext is an interface to support dynamic dispatch.

type IConstantDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LET() antlr.TerminalNode
	PatternInitializerList() IPatternInitializerListContext
	Attributes() IAttributesContext
	DeclarationModifiers() IDeclarationModifiersContext

	// IsConstantDeclarationContext differentiates from other interfaces.
	IsConstantDeclarationContext()
}
    IConstantDeclarationContext is an interface to support dynamic dispatch.

type IControlTransferContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsControlTransferContext differentiates from other interfaces.
	IsControlTransferContext()
}
    IControlTransferContext is an interface to support dynamic dispatch.

type IDecimalVersionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllINTEGER_LITERAL() []antlr.TerminalNode
	INTEGER_LITERAL(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsDecimalVersionContext differentiates from other interfaces.
	IsDecimalVersionContext()
}
    IDecimalVersionContext is an interface to support dynamic dispatch.

type IDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ImportDeclaration() IImportDeclarationContext
	ConstantDeclaration() IConstantDeclarationContext
	VariableDeclaration() IVariableDeclarationContext
	TypealiasDeclaration() ITypealiasDeclarationContext
	FunctionDeclaration() IFunctionDeclarationContext
	EnumDeclaration() IEnumDeclarationContext
	StructDeclaration() IStructDeclarationContext
	ClassDeclaration() IClassDeclarationContext
	ActorDeclaration() IActorDeclarationContext
	ProtocolDeclaration() IProtocolDeclarationContext
	InitializerDeclaration() IInitializerDeclarationContext
	DeinitializerDeclaration() IDeinitializerDeclarationContext
	ExtensionDeclaration() IExtensionDeclarationContext
	SubscriptDeclaration() ISubscriptDeclarationContext

	// IsDeclarationContext differentiates from other interfaces.
	IsDeclarationContext()
}
    IDeclarationContext is an interface to support dynamic dispatch.

type IDeclarationModifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AccessLevelModifier() IAccessLevelModifierContext
	MutationModifier() IMutationModifierContext
	CLASS() antlr.TerminalNode
	FINAL_KW() antlr.TerminalNode
	LAZY_KW() antlr.TerminalNode
	OPTIONAL_KW() antlr.TerminalNode
	OVERRIDE_KW() antlr.TerminalNode
	POSTFIX_KW() antlr.TerminalNode
	PREFIX_KW() antlr.TerminalNode
	REQUIRED_KW() antlr.TerminalNode
	STATIC() antlr.TerminalNode
	UNOWNED_KW() antlr.TerminalNode
	WEAK_KW() antlr.TerminalNode
	NONISOLATED_KW() antlr.TerminalNode

	// IsDeclarationModifierContext differentiates from other interfaces.
	IsDeclarationModifierContext()
}
    IDeclarationModifierContext is an interface to support dynamic dispatch.

type IDeclarationModifiersContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDeclarationModifier() []IDeclarationModifierContext
	DeclarationModifier(i int) IDeclarationModifierContext

	// IsDeclarationModifiersContext differentiates from other interfaces.
	IsDeclarationModifiersContext()
}
    IDeclarationModifiersContext is an interface to support dynamic dispatch.

type IDefaultArgumentClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext

	// IsDefaultArgumentClauseContext differentiates from other interfaces.
	IsDefaultArgumentClauseContext()
}
    IDefaultArgumentClauseContext is an interface to support dynamic dispatch.

type IDefaultLabelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEFAULT() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsDefaultLabelContext differentiates from other interfaces.
	IsDefaultLabelContext()
}
    IDefaultLabelContext is an interface to support dynamic dispatch.

type IDeferStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEFER() antlr.TerminalNode
	CodeBlock() ICodeBlockContext

	// IsDeferStatementContext differentiates from other interfaces.
	IsDeferStatementContext()
}
    IDeferStatementContext is an interface to support dynamic dispatch.

type IDeinitializerDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEINIT() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	Attributes() IAttributesContext

	// IsDeinitializerDeclarationContext differentiates from other interfaces.
	IsDeinitializerDeclarationContext()
}
    IDeinitializerDeclarationContext is an interface to support dynamic
    dispatch.

type IDiagnosticStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATOR() antlr.TerminalNode
	STRING_LITERAL() antlr.TerminalNode

	// IsDiagnosticStatementContext differentiates from other interfaces.
	IsDiagnosticStatementContext()
}
    IDiagnosticStatementContext is an interface to support dynamic dispatch.

type IDictionaryLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	DictionaryLiteralItems() IDictionaryLiteralItemsContext
	RBRACKET() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsDictionaryLiteralContext differentiates from other interfaces.
	IsDictionaryLiteralContext()
}
    IDictionaryLiteralContext is an interface to support dynamic dispatch.

type IDictionaryLiteralItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	COLON() antlr.TerminalNode

	// IsDictionaryLiteralItemContext differentiates from other interfaces.
	IsDictionaryLiteralItemContext()
}
    IDictionaryLiteralItemContext is an interface to support dynamic dispatch.

type IDictionaryLiteralItemsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDictionaryLiteralItem() []IDictionaryLiteralItemContext
	DictionaryLiteralItem(i int) IDictionaryLiteralItemContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsDictionaryLiteralItemsContext differentiates from other interfaces.
	IsDictionaryLiteralItemsContext()
}
    IDictionaryLiteralItemsContext is an interface to support dynamic dispatch.

type IDictionaryTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	AllType_() []ITypeContext
	Type_(i int) ITypeContext
	COLON() antlr.TerminalNode
	RBRACKET() antlr.TerminalNode

	// IsDictionaryTypeContext differentiates from other interfaces.
	IsDictionaryTypeContext()
}
    IDictionaryTypeContext is an interface to support dynamic dispatch.

type IDidSetClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DIDSET_KW() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	Attributes() IAttributesContext
	SetterName() ISetterNameContext

	// IsDidSetClauseContext differentiates from other interfaces.
	IsDidSetClauseContext()
}
    IDidSetClauseContext is an interface to support dynamic dispatch.

type IDoStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DO() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	AllCatchClause() []ICatchClauseContext
	CatchClause(i int) ICatchClauseContext

	// IsDoStatementContext differentiates from other interfaces.
	IsDoStatementContext()
}
    IDoStatementContext is an interface to support dynamic dispatch.

type IElseClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ELSE() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	IfStatement() IIfStatementContext

	// IsElseClauseContext differentiates from other interfaces.
	IsElseClauseContext()
}
    IElseClauseContext is an interface to support dynamic dispatch.

type IEnumCasePatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOT() antlr.TerminalNode
	Identifier() IIdentifierContext
	TypeIdentifier() ITypeIdentifierContext
	TuplePattern() ITuplePatternContext

	// IsEnumCasePatternContext differentiates from other interfaces.
	IsEnumCasePatternContext()
}
    IEnumCasePatternContext is an interface to support dynamic dispatch.

type IEnumDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ENUM() antlr.TerminalNode
	Identifier() IIdentifierContext
	LBRACE() antlr.TerminalNode
	EnumMembers() IEnumMembersContext
	RBRACE() antlr.TerminalNode
	Attributes() IAttributesContext
	AccessLevelModifier() IAccessLevelModifierContext
	GenericParameterClause() IGenericParameterClauseContext
	TypeInheritanceClause() ITypeInheritanceClauseContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsEnumDeclarationContext differentiates from other interfaces.
	IsEnumDeclarationContext()
}
    IEnumDeclarationContext is an interface to support dynamic dispatch.

type IEnumMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Declaration() IDeclarationContext
	UnionStyleEnumCaseClause() IUnionStyleEnumCaseClauseContext
	RawValueStyleEnumCaseClause() IRawValueStyleEnumCaseClauseContext
	CompilerControl() ICompilerControlContext

	// IsEnumMemberContext differentiates from other interfaces.
	IsEnumMemberContext()
}
    IEnumMemberContext is an interface to support dynamic dispatch.

type IEnumMembersContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllEnumMember() []IEnumMemberContext
	EnumMember(i int) IEnumMemberContext

	// IsEnumMembersContext differentiates from other interfaces.
	IsEnumMembersContext()
}
    IEnumMembersContext is an interface to support dynamic dispatch.

type IExistentialTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ANY() antlr.TerminalNode
	Type_() ITypeContext

	// IsExistentialTypeContext differentiates from other interfaces.
	IsExistentialTypeContext()
}
    IExistentialTypeContext is an interface to support dynamic dispatch.

type IExplicitMemberSuffixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOT() antlr.TerminalNode
	Identifier() IIdentifierContext
	GenericArgumentClause() IGenericArgumentClauseContext
	INTEGER_LITERAL() antlr.TerminalNode

	// IsExplicitMemberSuffixContext differentiates from other interfaces.
	IsExplicitMemberSuffixContext()
}
    IExplicitMemberSuffixContext is an interface to support dynamic dispatch.

type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PrefixExpression() IPrefixExpressionContext
	TryOperator() ITryOperatorContext
	AwaitOperator() IAwaitOperatorContext
	BinaryExpressions() IBinaryExpressionsContext

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}
    IExpressionContext is an interface to support dynamic dispatch.

type IExpressionPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext

	// IsExpressionPatternContext differentiates from other interfaces.
	IsExpressionPatternContext()
}
    IExpressionPatternContext is an interface to support dynamic dispatch.

type IExtensionBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllExtensionMember() []IExtensionMemberContext
	ExtensionMember(i int) IExtensionMemberContext

	// IsExtensionBodyContext differentiates from other interfaces.
	IsExtensionBodyContext()
}
    IExtensionBodyContext is an interface to support dynamic dispatch.

type IExtensionDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXTENSION() antlr.TerminalNode
	TypeIdentifier() ITypeIdentifierContext
	ExtensionBody() IExtensionBodyContext
	Attributes() IAttributesContext
	AccessLevelModifier() IAccessLevelModifierContext
	COLON() antlr.TerminalNode
	TypeInheritanceList() ITypeInheritanceListContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsExtensionDeclarationContext differentiates from other interfaces.
	IsExtensionDeclarationContext()
}
    IExtensionDeclarationContext is an interface to support dynamic dispatch.

type IExtensionMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Declaration() IDeclarationContext
	CompilerControl() ICompilerControlContext

	// IsExtensionMemberContext differentiates from other interfaces.
	IsExtensionMemberContext()
}
    IExtensionMemberContext is an interface to support dynamic dispatch.

type IExternalParameterNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	UNDERSCORE() antlr.TerminalNode

	// IsExternalParameterNameContext differentiates from other interfaces.
	IsExternalParameterNameContext()
}
    IExternalParameterNameContext is an interface to support dynamic dispatch.

type IForInStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FOR() antlr.TerminalNode
	Pattern() IPatternContext
	IN() antlr.TerminalNode
	Expression() IExpressionContext
	CodeBlock() ICodeBlockContext
	CASE() antlr.TerminalNode
	WhereClause() IWhereClauseContext

	// IsForInStatementContext differentiates from other interfaces.
	IsForInStatementContext()
}
    IForInStatementContext is an interface to support dynamic dispatch.

type IForStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FOR() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AllSEMICOLON() []antlr.TerminalNode
	SEMICOLON(i int) antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext

	// IsForStatementContext differentiates from other interfaces.
	IsForStatementContext()
}
    IForStatementContext is an interface to support dynamic dispatch.

type IForcedValueSuffixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXCLAIM_POSTFIX() antlr.TerminalNode

	// IsForcedValueSuffixContext differentiates from other interfaces.
	IsForcedValueSuffixContext()
}
    IForcedValueSuffixContext is an interface to support dynamic dispatch.

type IFunctionBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CodeBlock() ICodeBlockContext

	// IsFunctionBodyContext differentiates from other interfaces.
	IsFunctionBodyContext()
}
    IFunctionBodyContext is an interface to support dynamic dispatch.

type IFunctionCallArgumentClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	FunctionCallArgumentList() IFunctionCallArgumentListContext

	// IsFunctionCallArgumentClauseContext differentiates from other interfaces.
	IsFunctionCallArgumentClauseContext()
}
    IFunctionCallArgumentClauseContext is an interface to support dynamic
    dispatch.

type IFunctionCallArgumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	Identifier() IIdentifierContext
	COLON() antlr.TerminalNode
	UNDERSCORE() antlr.TerminalNode
	Operator_() IOperator_Context

	// IsFunctionCallArgumentContext differentiates from other interfaces.
	IsFunctionCallArgumentContext()
}
    IFunctionCallArgumentContext is an interface to support dynamic dispatch.

type IFunctionCallArgumentListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllFunctionCallArgument() []IFunctionCallArgumentContext
	FunctionCallArgument(i int) IFunctionCallArgumentContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFunctionCallArgumentListContext differentiates from other interfaces.
	IsFunctionCallArgumentListContext()
}
    IFunctionCallArgumentListContext is an interface to support dynamic
    dispatch.

type IFunctionCallSuffixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FunctionCallArgumentClause() IFunctionCallArgumentClauseContext
	TrailingClosures() ITrailingClosuresContext

	// IsFunctionCallSuffixContext differentiates from other interfaces.
	IsFunctionCallSuffixContext()
}
    IFunctionCallSuffixContext is an interface to support dynamic dispatch.

type IFunctionDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FunctionHead() IFunctionHeadContext
	Identifier() IIdentifierContext
	FunctionSignature() IFunctionSignatureContext
	GenericParameterClause() IGenericParameterClauseContext
	GenericWhereClause() IGenericWhereClauseContext
	FunctionBody() IFunctionBodyContext

	// IsFunctionDeclarationContext differentiates from other interfaces.
	IsFunctionDeclarationContext()
}
    IFunctionDeclarationContext is an interface to support dynamic dispatch.

type IFunctionHeadContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNC() antlr.TerminalNode
	Attributes() IAttributesContext
	DeclarationModifiers() IDeclarationModifiersContext

	// IsFunctionHeadContext differentiates from other interfaces.
	IsFunctionHeadContext()
}
    IFunctionHeadContext is an interface to support dynamic dispatch.

type IFunctionResultContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ARROW() antlr.TerminalNode
	Type_() ITypeContext
	Attributes() IAttributesContext

	// IsFunctionResultContext differentiates from other interfaces.
	IsFunctionResultContext()
}
    IFunctionResultContext is an interface to support dynamic dispatch.

type IFunctionSignatureContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ParameterClause() IParameterClauseContext
	AsyncModifier() IAsyncModifierContext
	THROWS() antlr.TerminalNode
	FunctionResult() IFunctionResultContext

	// IsFunctionSignatureContext differentiates from other interfaces.
	IsFunctionSignatureContext()
}
    IFunctionSignatureContext is an interface to support dynamic dispatch.

type IFunctionTypeArgumentClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	FunctionTypeArgumentList() IFunctionTypeArgumentListContext
	ELLIPSIS() antlr.TerminalNode

	// IsFunctionTypeArgumentClauseContext differentiates from other interfaces.
	IsFunctionTypeArgumentClauseContext()
}
    IFunctionTypeArgumentClauseContext is an interface to support dynamic
    dispatch.

type IFunctionTypeArgumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Type_() ITypeContext
	Attributes() IAttributesContext
	INOUT() antlr.TerminalNode
	ArgumentLabel() IArgumentLabelContext
	TypeAnnotation() ITypeAnnotationContext

	// IsFunctionTypeArgumentContext differentiates from other interfaces.
	IsFunctionTypeArgumentContext()
}
    IFunctionTypeArgumentContext is an interface to support dynamic dispatch.

type IFunctionTypeArgumentListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllFunctionTypeArgument() []IFunctionTypeArgumentContext
	FunctionTypeArgument(i int) IFunctionTypeArgumentContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFunctionTypeArgumentListContext differentiates from other interfaces.
	IsFunctionTypeArgumentListContext()
}
    IFunctionTypeArgumentListContext is an interface to support dynamic
    dispatch.

type IFunctionTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FunctionTypeArgumentClause() IFunctionTypeArgumentClauseContext
	ARROW() antlr.TerminalNode
	Type_() ITypeContext
	Attributes() IAttributesContext
	AsyncModifier() IAsyncModifierContext
	THROWS() antlr.TerminalNode

	// IsFunctionTypeContext differentiates from other interfaces.
	IsFunctionTypeContext()
}
    IFunctionTypeContext is an interface to support dynamic dispatch.

type IGenericArgumentClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LT() antlr.TerminalNode
	GenericArgumentList() IGenericArgumentListContext
	GT() antlr.TerminalNode

	// IsGenericArgumentClauseContext differentiates from other interfaces.
	IsGenericArgumentClauseContext()
}
    IGenericArgumentClauseContext is an interface to support dynamic dispatch.

type IGenericArgumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Type_() ITypeContext

	// IsGenericArgumentContext differentiates from other interfaces.
	IsGenericArgumentContext()
}
    IGenericArgumentContext is an interface to support dynamic dispatch.

type IGenericArgumentListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllGenericArgument() []IGenericArgumentContext
	GenericArgument(i int) IGenericArgumentContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsGenericArgumentListContext differentiates from other interfaces.
	IsGenericArgumentListContext()
}
    IGenericArgumentListContext is an interface to support dynamic dispatch.

type IGenericParameterClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LT() antlr.TerminalNode
	GenericParameterList() IGenericParameterListContext
	GT() antlr.TerminalNode

	// IsGenericParameterClauseContext differentiates from other interfaces.
	IsGenericParameterClauseContext()
}
    IGenericParameterClauseContext is an interface to support dynamic dispatch.

type IGenericParameterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	COLON() antlr.TerminalNode
	TypeIdentifier() ITypeIdentifierContext
	ProtocolCompositionType() IProtocolCompositionTypeContext

	// IsGenericParameterContext differentiates from other interfaces.
	IsGenericParameterContext()
}
    IGenericParameterContext is an interface to support dynamic dispatch.

type IGenericParameterListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllGenericParameter() []IGenericParameterContext
	GenericParameter(i int) IGenericParameterContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsGenericParameterListContext differentiates from other interfaces.
	IsGenericParameterListContext()
}
    IGenericParameterListContext is an interface to support dynamic dispatch.

type IGenericWhereClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHERE() antlr.TerminalNode
	RequirementList() IRequirementListContext

	// IsGenericWhereClauseContext differentiates from other interfaces.
	IsGenericWhereClauseContext()
}
    IGenericWhereClauseContext is an interface to support dynamic dispatch.

type IGetterClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	GET_KW() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	Attributes() IAttributesContext
	MutationModifier() IMutationModifierContext

	// IsGetterClauseContext differentiates from other interfaces.
	IsGetterClauseContext()
}
    IGetterClauseContext is an interface to support dynamic dispatch.

type IGetterKeywordClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	GET_KW() antlr.TerminalNode
	Attributes() IAttributesContext
	MutationModifier() IMutationModifierContext

	// IsGetterKeywordClauseContext differentiates from other interfaces.
	IsGetterKeywordClauseContext()
}
    IGetterKeywordClauseContext is an interface to support dynamic dispatch.

type IGetterSetterBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	GetterClause() IGetterClauseContext
	RBRACE() antlr.TerminalNode
	SetterClause() ISetterClauseContext

	// IsGetterSetterBlockContext differentiates from other interfaces.
	IsGetterSetterBlockContext()
}
    IGetterSetterBlockContext is an interface to support dynamic dispatch.

type IGetterSetterKeywordBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	GetterKeywordClause() IGetterKeywordClauseContext
	RBRACE() antlr.TerminalNode
	SetterKeywordClause() ISetterKeywordClauseContext

	// IsGetterSetterKeywordBlockContext differentiates from other interfaces.
	IsGetterSetterKeywordBlockContext()
}
    IGetterSetterKeywordBlockContext is an interface to support dynamic
    dispatch.

type IGuardStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	GUARD() antlr.TerminalNode
	ConditionList() IConditionListContext
	ELSE() antlr.TerminalNode
	CodeBlock() ICodeBlockContext

	// IsGuardStatementContext differentiates from other interfaces.
	IsGuardStatementContext()
}
    IGuardStatementContext is an interface to support dynamic dispatch.

type IIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	OS_KW() antlr.TerminalNode
	ARCH_KW() antlr.TerminalNode
	VERTEX_KW() antlr.TerminalNode
	COMPILER_KW() antlr.TerminalNode
	CANIMPORT_KW() antlr.TerminalNode
	VERSION_KW() antlr.TerminalNode
	FILE_KEYWORD() antlr.TerminalNode
	LINE_KEYWORD() antlr.TerminalNode
	GET_KW() antlr.TerminalNode
	SET_KW() antlr.TerminalNode
	WILLSET_KW() antlr.TerminalNode
	DIDSET_KW() antlr.TerminalNode
	ASYNC_KW() antlr.TerminalNode
	ACTOR_KW() antlr.TerminalNode
	PREFIX_KW() antlr.TerminalNode
	POSTFIX_KW() antlr.TerminalNode
	MACRO_KW() antlr.TerminalNode
	DYNAMIC_KW() antlr.TerminalNode
	FINAL_KW() antlr.TerminalNode
	LAZY_KW() antlr.TerminalNode
	OPTIONAL_KW() antlr.TerminalNode
	OVERRIDE_KW() antlr.TerminalNode
	REQUIRED_KW() antlr.TerminalNode
	UNOWNED_KW() antlr.TerminalNode
	WEAK_KW() antlr.TerminalNode
	NONISOLATED_KW() antlr.TerminalNode
	MUTATING_KW() antlr.TerminalNode
	NONMUTATING_KW() antlr.TerminalNode
	SOME_KW() antlr.TerminalNode
	TYPE_KW() antlr.TerminalNode
	PROTOCOL_KW() antlr.TerminalNode
	CONSUMING_KW() antlr.TerminalNode
	BORROWING_KW() antlr.TerminalNode
	SENDING_KW() antlr.TerminalNode

	// IsIdentifierContext differentiates from other interfaces.
	IsIdentifierContext()
}
    IIdentifierContext is an interface to support dynamic dispatch.

type IIdentifierListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIdentifier() []IIdentifierContext
	Identifier(i int) IIdentifierContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsIdentifierListContext differentiates from other interfaces.
	IsIdentifierListContext()
}
    IIdentifierListContext is an interface to support dynamic dispatch.

type IIdentifierPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext

	// IsIdentifierPatternContext differentiates from other interfaces.
	IsIdentifierPatternContext()
}
    IIdentifierPatternContext is an interface to support dynamic dispatch.

type IIfStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IF() antlr.TerminalNode
	ConditionList() IConditionListContext
	CodeBlock() ICodeBlockContext
	ElseClause() IElseClauseContext

	// IsIfStatementContext differentiates from other interfaces.
	IsIfStatementContext()
}
    IIfStatementContext is an interface to support dynamic dispatch.

type IImplicitMemberExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOT() antlr.TerminalNode
	Identifier() IIdentifierContext
	GenericArgumentClause() IGenericArgumentClauseContext

	// IsImplicitMemberExpressionContext differentiates from other interfaces.
	IsImplicitMemberExpressionContext()
}
    IImplicitMemberExpressionContext is an interface to support dynamic
    dispatch.

type IImportAliasContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	DOT() antlr.TerminalNode
	UNDERSCORE() antlr.TerminalNode

	// IsImportAliasContext differentiates from other interfaces.
	IsImportAliasContext()
}
    IImportAliasContext is an interface to support dynamic dispatch.

type IImportDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IMPORT() antlr.TerminalNode
	AllImportSpec() []IImportSpecContext
	ImportSpec(i int) IImportSpecContext
	Attributes() IAttributesContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode

	// IsImportDeclarationContext differentiates from other interfaces.
	IsImportDeclarationContext()
}
    IImportDeclarationContext is an interface to support dynamic dispatch.

type IImportSpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING_LITERAL() antlr.TerminalNode
	ImportAlias() IImportAliasContext

	// IsImportSpecContext differentiates from other interfaces.
	IsImportSpecContext()
}
    IImportSpecContext is an interface to support dynamic dispatch.

type IInOutExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AMPERSAND() antlr.TerminalNode
	Identifier() IIdentifierContext

	// IsInOutExpressionContext differentiates from other interfaces.
	IsInOutExpressionContext()
}
    IInOutExpressionContext is an interface to support dynamic dispatch.

type IInitializerBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CodeBlock() ICodeBlockContext

	// IsInitializerBodyContext differentiates from other interfaces.
	IsInitializerBodyContext()
}
    IInitializerBodyContext is an interface to support dynamic dispatch.

type IInitializerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext

	// IsInitializerContext differentiates from other interfaces.
	IsInitializerContext()
}
    IInitializerContext is an interface to support dynamic dispatch.

type IInitializerDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INIT() antlr.TerminalNode
	ParameterClause() IParameterClauseContext
	InitializerBody() IInitializerBodyContext
	Attributes() IAttributesContext
	DeclarationModifiers() IDeclarationModifiersContext
	QUESTION_POSTFIX() antlr.TerminalNode
	GenericParameterClause() IGenericParameterClauseContext
	AsyncModifier() IAsyncModifierContext
	THROWS() antlr.TerminalNode
	GenericWhereClause() IGenericWhereClauseContext

	// IsInitializerDeclarationContext differentiates from other interfaces.
	IsInitializerDeclarationContext()
}
    IInitializerDeclarationContext is an interface to support dynamic dispatch.

type IInitializerSuffixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOT() antlr.TerminalNode
	INIT() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ArgumentNames() IArgumentNamesContext
	RPAREN() antlr.TerminalNode

	// IsInitializerSuffixContext differentiates from other interfaces.
	IsInitializerSuffixContext()
}
    IInitializerSuffixContext is an interface to support dynamic dispatch.

type ILabelNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext

	// IsLabelNameContext differentiates from other interfaces.
	IsLabelNameContext()
}
    ILabelNameContext is an interface to support dynamic dispatch.

type ILabeledStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LabelName() ILabelNameContext
	COLON() antlr.TerminalNode
	LoopStatement() ILoopStatementContext
	IfStatement() IIfStatementContext
	SwitchStatement() ISwitchStatementContext
	DoStatement() IDoStatementContext

	// IsLabeledStatementContext differentiates from other interfaces.
	IsLabeledStatementContext()
}
    ILabeledStatementContext is an interface to support dynamic dispatch.

type ILabeledTrailingClosureContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	COLON() antlr.TerminalNode
	ClosureExpression() IClosureExpressionContext

	// IsLabeledTrailingClosureContext differentiates from other interfaces.
	IsLabeledTrailingClosureContext()
}
    ILabeledTrailingClosureContext is an interface to support dynamic dispatch.

type ILayoutConstraintRequirementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TypeIdentifier() ITypeIdentifierContext
	COLON() antlr.TerminalNode
	Identifier() IIdentifierContext

	// IsLayoutConstraintRequirementContext differentiates from other interfaces.
	IsLayoutConstraintRequirementContext()
}
    ILayoutConstraintRequirementContext is an interface to support dynamic
    dispatch.

type ILineControlStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	POUND_SOURCE_LOCATION() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	FILE_KEYWORD() antlr.TerminalNode
	AllCOLON() []antlr.TerminalNode
	COLON(i int) antlr.TerminalNode
	STRING_LITERAL() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	LINE_KEYWORD() antlr.TerminalNode
	INTEGER_LITERAL() antlr.TerminalNode
	RPAREN() antlr.TerminalNode

	// IsLineControlStatementContext differentiates from other interfaces.
	IsLineControlStatementContext()
}
    ILineControlStatementContext is an interface to support dynamic dispatch.

type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NumericLiteral() INumericLiteralContext
	STRING_LITERAL() antlr.TerminalNode
	MULTILINE_STRING_LITERAL() antlr.TerminalNode
	EXTENDED_STRING_LITERAL() antlr.TerminalNode
	BooleanLiteral() IBooleanLiteralContext
	NIL() antlr.TerminalNode

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}
    ILiteralContext is an interface to support dynamic dispatch.

type ILiteralExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Literal() ILiteralContext
	ArrayLiteral() IArrayLiteralContext
	DictionaryLiteral() IDictionaryLiteralContext
	PoundFileExpression() IPoundFileExpressionContext

	// IsLiteralExpressionContext differentiates from other interfaces.
	IsLiteralExpressionContext()
}
    ILiteralExpressionContext is an interface to support dynamic dispatch.

type ILocalParameterNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	UNDERSCORE() antlr.TerminalNode

	// IsLocalParameterNameContext differentiates from other interfaces.
	IsLocalParameterNameContext()
}
    ILocalParameterNameContext is an interface to support dynamic dispatch.

type ILoopStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ForInStatement() IForInStatementContext
	WhileStatement() IWhileStatementContext
	RepeatWhileStatement() IRepeatWhileStatementContext
	ForStatement() IForStatementContext

	// IsLoopStatementContext differentiates from other interfaces.
	IsLoopStatementContext()
}
    ILoopStatementContext is an interface to support dynamic dispatch.

type IMacroDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MACRO_KW() antlr.TerminalNode
	Identifier() IIdentifierContext
	Attributes() IAttributesContext
	DeclarationModifiers() IDeclarationModifiersContext
	GenericParameterClause() IGenericParameterClauseContext
	ParameterClause() IParameterClauseContext
	ARROW() antlr.TerminalNode
	Type_() ITypeContext
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsMacroDeclarationContext differentiates from other interfaces.
	IsMacroDeclarationContext()
}
    IMacroDeclarationContext is an interface to support dynamic dispatch.

type IMacroExpansionExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HASH() antlr.TerminalNode
	Identifier() IIdentifierContext
	GenericArgumentClause() IGenericArgumentClauseContext
	FunctionCallArgumentClause() IFunctionCallArgumentClauseContext
	TrailingClosures() ITrailingClosuresContext

	// IsMacroExpansionExpressionContext differentiates from other interfaces.
	IsMacroExpansionExpressionContext()
}
    IMacroExpansionExpressionContext is an interface to support dynamic
    dispatch.

type IMutationModifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MUTATING_KW() antlr.TerminalNode
	NONMUTATING_KW() antlr.TerminalNode

	// IsMutationModifierContext differentiates from other interfaces.
	IsMutationModifierContext()
}
    IMutationModifierContext is an interface to support dynamic dispatch.

type INumericLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INTEGER_LITERAL() antlr.TerminalNode
	OPERATOR() antlr.TerminalNode
	FLOAT_LITERAL() antlr.TerminalNode

	// IsNumericLiteralContext differentiates from other interfaces.
	IsNumericLiteralContext()
}
    INumericLiteralContext is an interface to support dynamic dispatch.

type IOpaqueTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SOME_KW() antlr.TerminalNode
	Type_() ITypeContext

	// IsOpaqueTypeContext differentiates from other interfaces.
	IsOpaqueTypeContext()
}
    IOpaqueTypeContext is an interface to support dynamic dispatch.

type IOperator_Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATOR() antlr.TerminalNode
	DOT() antlr.TerminalNode

	// IsOperator_Context differentiates from other interfaces.
	IsOperator_Context()
}
    IOperator_Context is an interface to support dynamic dispatch.

type IOptionalBindingConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Pattern() IPatternContext
	Initializer() IInitializerContext
	LET() antlr.TerminalNode
	VAR() antlr.TerminalNode

	// IsOptionalBindingConditionContext differentiates from other interfaces.
	IsOptionalBindingConditionContext()
}
    IOptionalBindingConditionContext is an interface to support dynamic
    dispatch.

type IOptionalChainingLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QUESTION_POSTFIX() antlr.TerminalNode

	// IsOptionalChainingLiteralContext differentiates from other interfaces.
	IsOptionalChainingLiteralContext()
}
    IOptionalChainingLiteralContext is an interface to support dynamic dispatch.

type IOptionalPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IdentifierPattern() IIdentifierPatternContext
	QUESTION_POSTFIX() antlr.TerminalNode

	// IsOptionalPatternContext differentiates from other interfaces.
	IsOptionalPatternContext()
}
    IOptionalPatternContext is an interface to support dynamic dispatch.

type IParameterClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ParameterList() IParameterListContext

	// IsParameterClauseContext differentiates from other interfaces.
	IsParameterClauseContext()
}
    IParameterClauseContext is an interface to support dynamic dispatch.

type IParameterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LocalParameterName() ILocalParameterNameContext
	TypeAnnotation() ITypeAnnotationContext
	ExternalParameterName() IExternalParameterNameContext
	DefaultArgumentClause() IDefaultArgumentClauseContext
	ELLIPSIS() antlr.TerminalNode

	// IsParameterContext differentiates from other interfaces.
	IsParameterContext()
}
    IParameterContext is an interface to support dynamic dispatch.

type IParameterListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllParameter() []IParameterContext
	Parameter(i int) IParameterContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsParameterListContext differentiates from other interfaces.
	IsParameterListContext()
}
    IParameterListContext is an interface to support dynamic dispatch.

type IParenthesizedExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN() antlr.TerminalNode

	// IsParenthesizedExpressionContext differentiates from other interfaces.
	IsParenthesizedExpressionContext()
}
    IParenthesizedExpressionContext is an interface to support dynamic dispatch.

type IPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPatternContext differentiates from other interfaces.
	IsPatternContext()
}
    IPatternContext is an interface to support dynamic dispatch.

type IPatternInitializerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Pattern() IPatternContext
	Initializer() IInitializerContext

	// IsPatternInitializerContext differentiates from other interfaces.
	IsPatternInitializerContext()
}
    IPatternInitializerContext is an interface to support dynamic dispatch.

type IPatternInitializerListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPatternInitializer() []IPatternInitializerContext
	PatternInitializer(i int) IPatternInitializerContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsPatternInitializerListContext differentiates from other interfaces.
	IsPatternInitializerListContext()
}
    IPatternInitializerListContext is an interface to support dynamic dispatch.

type IPlatformConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OS_KW() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	Identifier() IIdentifierContext
	RPAREN() antlr.TerminalNode
	ARCH_KW() antlr.TerminalNode
	VERTEX_KW() antlr.TerminalNode
	OPERATOR() antlr.TerminalNode
	DecimalVersion() IDecimalVersionContext
	COMPILER_KW() antlr.TerminalNode
	CANIMPORT_KW() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	VERSION_KW() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsPlatformConditionContext differentiates from other interfaces.
	IsPlatformConditionContext()
}
    IPlatformConditionContext is an interface to support dynamic dispatch.

type IPostfixExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PrimaryExpression() IPrimaryExpressionContext
	AllPostfixSuffix() []IPostfixSuffixContext
	PostfixSuffix(i int) IPostfixSuffixContext

	// IsPostfixExpressionContext differentiates from other interfaces.
	IsPostfixExpressionContext()
}
    IPostfixExpressionContext is an interface to support dynamic dispatch.

type IPostfixOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATOR() antlr.TerminalNode

	// IsPostfixOperatorContext differentiates from other interfaces.
	IsPostfixOperatorContext()
}
    IPostfixOperatorContext is an interface to support dynamic dispatch.

type IPostfixSelfSuffixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOT() antlr.TerminalNode
	SELF() antlr.TerminalNode

	// IsPostfixSelfSuffixContext differentiates from other interfaces.
	IsPostfixSelfSuffixContext()
}
    IPostfixSelfSuffixContext is an interface to support dynamic dispatch.

type IPostfixSuffixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PostfixOperator() IPostfixOperatorContext
	FunctionCallSuffix() IFunctionCallSuffixContext
	InitializerSuffix() IInitializerSuffixContext
	ExplicitMemberSuffix() IExplicitMemberSuffixContext
	PostfixSelfSuffix() IPostfixSelfSuffixContext
	SubscriptSuffix() ISubscriptSuffixContext
	ForcedValueSuffix() IForcedValueSuffixContext
	OptionalChainingLiteral() IOptionalChainingLiteralContext

	// IsPostfixSuffixContext differentiates from other interfaces.
	IsPostfixSuffixContext()
}
    IPostfixSuffixContext is an interface to support dynamic dispatch.

type IPoundFileExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	POUND_FILE() antlr.TerminalNode
	POUND_FILEID() antlr.TerminalNode
	POUND_FILEPATH() antlr.TerminalNode
	POUND_LINE() antlr.TerminalNode
	POUND_COLUMN() antlr.TerminalNode
	POUND_FUNCTION() antlr.TerminalNode

	// IsPoundFileExpressionContext differentiates from other interfaces.
	IsPoundFileExpressionContext()
}
    IPoundFileExpressionContext is an interface to support dynamic dispatch.

type IPrefixExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPrefixExpressionContext differentiates from other interfaces.
	IsPrefixExpressionContext()
}
    IPrefixExpressionContext is an interface to support dynamic dispatch.

type IPrefixOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATOR() antlr.TerminalNode

	// IsPrefixOperatorContext differentiates from other interfaces.
	IsPrefixOperatorContext()
}
    IPrefixOperatorContext is an interface to support dynamic dispatch.

type IPrimaryAssociatedTypeClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LT() antlr.TerminalNode
	PrimaryAssociatedTypeList() IPrimaryAssociatedTypeListContext
	GT() antlr.TerminalNode

	// IsPrimaryAssociatedTypeClauseContext differentiates from other interfaces.
	IsPrimaryAssociatedTypeClauseContext()
}
    IPrimaryAssociatedTypeClauseContext is an interface to support dynamic
    dispatch.

type IPrimaryAssociatedTypeListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIdentifier() []IIdentifierContext
	Identifier(i int) IIdentifierContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsPrimaryAssociatedTypeListContext differentiates from other interfaces.
	IsPrimaryAssociatedTypeListContext()
}
    IPrimaryAssociatedTypeListContext is an interface to support dynamic
    dispatch.

type IPrimaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPrimaryExpressionContext differentiates from other interfaces.
	IsPrimaryExpressionContext()
}
    IPrimaryExpressionContext is an interface to support dynamic dispatch.

type IProtocolAssociatedTypeDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASSOCIATEDTYPE() antlr.TerminalNode
	Identifier() IIdentifierContext
	Attributes() IAttributesContext
	AccessLevelModifier() IAccessLevelModifierContext
	TypeInheritanceClause() ITypeInheritanceClauseContext
	ASSIGN() antlr.TerminalNode
	Type_() ITypeContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsProtocolAssociatedTypeDeclarationContext differentiates from other interfaces.
	IsProtocolAssociatedTypeDeclarationContext()
}
    IProtocolAssociatedTypeDeclarationContext is an interface to support dynamic
    dispatch.

type IProtocolBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllProtocolMember() []IProtocolMemberContext
	ProtocolMember(i int) IProtocolMemberContext

	// IsProtocolBodyContext differentiates from other interfaces.
	IsProtocolBodyContext()
}
    IProtocolBodyContext is an interface to support dynamic dispatch.

type IProtocolCompositionTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllProtocolCompositionTypeElement() []IProtocolCompositionTypeElementContext
	ProtocolCompositionTypeElement(i int) IProtocolCompositionTypeElementContext
	AllAMPERSAND() []antlr.TerminalNode
	AMPERSAND(i int) antlr.TerminalNode
	ANY() antlr.TerminalNode

	// IsProtocolCompositionTypeContext differentiates from other interfaces.
	IsProtocolCompositionTypeContext()
}
    IProtocolCompositionTypeContext is an interface to support dynamic dispatch.

type IProtocolCompositionTypeElementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TypeIdentifier() ITypeIdentifierContext
	SuppressedType() ISuppressedTypeContext

	// IsProtocolCompositionTypeElementContext differentiates from other interfaces.
	IsProtocolCompositionTypeElementContext()
}
    IProtocolCompositionTypeElementContext is an interface to support dynamic
    dispatch.

type IProtocolDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PROTOCOL() antlr.TerminalNode
	Identifier() IIdentifierContext
	ProtocolBody() IProtocolBodyContext
	Attributes() IAttributesContext
	AccessLevelModifier() IAccessLevelModifierContext
	PrimaryAssociatedTypeClause() IPrimaryAssociatedTypeClauseContext
	TypeInheritanceClause() ITypeInheritanceClauseContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsProtocolDeclarationContext differentiates from other interfaces.
	IsProtocolDeclarationContext()
}
    IProtocolDeclarationContext is an interface to support dynamic dispatch.

type IProtocolInitializerDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INIT() antlr.TerminalNode
	ParameterClause() IParameterClauseContext
	Attributes() IAttributesContext
	DeclarationModifiers() IDeclarationModifiersContext
	QUESTION_POSTFIX() antlr.TerminalNode
	GenericParameterClause() IGenericParameterClauseContext
	THROWS() antlr.TerminalNode
	GenericWhereClause() IGenericWhereClauseContext

	// IsProtocolInitializerDeclarationContext differentiates from other interfaces.
	IsProtocolInitializerDeclarationContext()
}
    IProtocolInitializerDeclarationContext is an interface to support dynamic
    dispatch.

type IProtocolMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ProtocolMemberDeclaration() IProtocolMemberDeclarationContext
	CompilerControl() ICompilerControlContext

	// IsProtocolMemberContext differentiates from other interfaces.
	IsProtocolMemberContext()
}
    IProtocolMemberContext is an interface to support dynamic dispatch.

type IProtocolMemberDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ProtocolPropertyDeclaration() IProtocolPropertyDeclarationContext
	ProtocolMethodDeclaration() IProtocolMethodDeclarationContext
	ProtocolInitializerDeclaration() IProtocolInitializerDeclarationContext
	ProtocolSubscriptDeclaration() IProtocolSubscriptDeclarationContext
	ProtocolAssociatedTypeDeclaration() IProtocolAssociatedTypeDeclarationContext
	TypealiasDeclaration() ITypealiasDeclarationContext

	// IsProtocolMemberDeclarationContext differentiates from other interfaces.
	IsProtocolMemberDeclarationContext()
}
    IProtocolMemberDeclarationContext is an interface to support dynamic
    dispatch.

type IProtocolMethodDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FunctionHead() IFunctionHeadContext
	Identifier() IIdentifierContext
	FunctionSignature() IFunctionSignatureContext
	GenericParameterClause() IGenericParameterClauseContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsProtocolMethodDeclarationContext differentiates from other interfaces.
	IsProtocolMethodDeclarationContext()
}
    IProtocolMethodDeclarationContext is an interface to support dynamic
    dispatch.

type IProtocolPropertyDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VariableDeclarationHead() IVariableDeclarationHeadContext
	VariableName() IVariableNameContext
	TypeAnnotation() ITypeAnnotationContext
	GetterSetterKeywordBlock() IGetterSetterKeywordBlockContext

	// IsProtocolPropertyDeclarationContext differentiates from other interfaces.
	IsProtocolPropertyDeclarationContext()
}
    IProtocolPropertyDeclarationContext is an interface to support dynamic
    dispatch.

type IProtocolSubscriptDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SubscriptHead() ISubscriptHeadContext
	SubscriptResult() ISubscriptResultContext
	GetterSetterKeywordBlock() IGetterSetterKeywordBlockContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsProtocolSubscriptDeclarationContext differentiates from other interfaces.
	IsProtocolSubscriptDeclarationContext()
}
    IProtocolSubscriptDeclarationContext is an interface to support dynamic
    dispatch.

type IRawValueAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASSIGN() antlr.TerminalNode
	RawValueLiteral() IRawValueLiteralContext

	// IsRawValueAssignmentContext differentiates from other interfaces.
	IsRawValueAssignmentContext()
}
    IRawValueAssignmentContext is an interface to support dynamic dispatch.

type IRawValueLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NumericLiteral() INumericLiteralContext
	STRING_LITERAL() antlr.TerminalNode
	BooleanLiteral() IBooleanLiteralContext

	// IsRawValueLiteralContext differentiates from other interfaces.
	IsRawValueLiteralContext()
}
    IRawValueLiteralContext is an interface to support dynamic dispatch.

type IRawValueStyleEnumCaseClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	RawValueStyleEnumCaseList() IRawValueStyleEnumCaseListContext
	Attributes() IAttributesContext

	// IsRawValueStyleEnumCaseClauseContext differentiates from other interfaces.
	IsRawValueStyleEnumCaseClauseContext()
}
    IRawValueStyleEnumCaseClauseContext is an interface to support dynamic
    dispatch.

type IRawValueStyleEnumCaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	RawValueAssignment() IRawValueAssignmentContext

	// IsRawValueStyleEnumCaseContext differentiates from other interfaces.
	IsRawValueStyleEnumCaseContext()
}
    IRawValueStyleEnumCaseContext is an interface to support dynamic dispatch.

type IRawValueStyleEnumCaseListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRawValueStyleEnumCase() []IRawValueStyleEnumCaseContext
	RawValueStyleEnumCase(i int) IRawValueStyleEnumCaseContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsRawValueStyleEnumCaseListContext differentiates from other interfaces.
	IsRawValueStyleEnumCaseListContext()
}
    IRawValueStyleEnumCaseListContext is an interface to support dynamic
    dispatch.

type IRepeatWhileStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	REPEAT() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	WHILE() antlr.TerminalNode
	Expression() IExpressionContext

	// IsRepeatWhileStatementContext differentiates from other interfaces.
	IsRepeatWhileStatementContext()
}
    IRepeatWhileStatementContext is an interface to support dynamic dispatch.

type IRequirementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ConformanceRequirement() IConformanceRequirementContext
	SameTypeRequirement() ISameTypeRequirementContext
	LayoutConstraintRequirement() ILayoutConstraintRequirementContext

	// IsRequirementContext differentiates from other interfaces.
	IsRequirementContext()
}
    IRequirementContext is an interface to support dynamic dispatch.

type IRequirementListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRequirement() []IRequirementContext
	Requirement(i int) IRequirementContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsRequirementListContext differentiates from other interfaces.
	IsRequirementListContext()
}
    IRequirementListContext is an interface to support dynamic dispatch.

type ISameTypeRequirementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TypeIdentifier() ITypeIdentifierContext
	ASSIGN() antlr.TerminalNode
	Type_() ITypeContext

	// IsSameTypeRequirementContext differentiates from other interfaces.
	IsSameTypeRequirementContext()
}
    ISameTypeRequirementContext is an interface to support dynamic dispatch.

type ISelfExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SELF() antlr.TerminalNode
	DOT() antlr.TerminalNode
	Identifier() IIdentifierContext
	LBRACKET() antlr.TerminalNode
	FunctionCallArgumentList() IFunctionCallArgumentListContext
	RBRACKET() antlr.TerminalNode
	INIT() antlr.TerminalNode

	// IsSelfExpressionContext differentiates from other interfaces.
	IsSelfExpressionContext()
}
    ISelfExpressionContext is an interface to support dynamic dispatch.

type ISelfTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SELF_UPPER() antlr.TerminalNode

	// IsSelfTypeContext differentiates from other interfaces.
	IsSelfTypeContext()
}
    ISelfTypeContext is an interface to support dynamic dispatch.

type ISetterClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SET_KW() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	Attributes() IAttributesContext
	MutationModifier() IMutationModifierContext
	SetterName() ISetterNameContext

	// IsSetterClauseContext differentiates from other interfaces.
	IsSetterClauseContext()
}
    ISetterClauseContext is an interface to support dynamic dispatch.

type ISetterKeywordClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SET_KW() antlr.TerminalNode
	Attributes() IAttributesContext
	MutationModifier() IMutationModifierContext

	// IsSetterKeywordClauseContext differentiates from other interfaces.
	IsSetterKeywordClauseContext()
}
    ISetterKeywordClauseContext is an interface to support dynamic dispatch.

type ISetterNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	Identifier() IIdentifierContext
	RPAREN() antlr.TerminalNode

	// IsSetterNameContext differentiates from other interfaces.
	IsSetterNameContext()
}
    ISetterNameContext is an interface to support dynamic dispatch.

type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}
    IStatementContext is an interface to support dynamic dispatch.

type IStatementsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllStatement() []IStatementContext
	Statement(i int) IStatementContext

	// IsStatementsContext differentiates from other interfaces.
	IsStatementsContext()
}
    IStatementsContext is an interface to support dynamic dispatch.

type IStructBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllStructMember() []IStructMemberContext
	StructMember(i int) IStructMemberContext

	// IsStructBodyContext differentiates from other interfaces.
	IsStructBodyContext()
}
    IStructBodyContext is an interface to support dynamic dispatch.

type IStructDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRUCT() antlr.TerminalNode
	Identifier() IIdentifierContext
	StructBody() IStructBodyContext
	Attributes() IAttributesContext
	AccessLevelModifier() IAccessLevelModifierContext
	GenericParameterClause() IGenericParameterClauseContext
	TypeInheritanceClause() ITypeInheritanceClauseContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsStructDeclarationContext differentiates from other interfaces.
	IsStructDeclarationContext()
}
    IStructDeclarationContext is an interface to support dynamic dispatch.

type IStructMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Declaration() IDeclarationContext
	CompilerControl() ICompilerControlContext

	// IsStructMemberContext differentiates from other interfaces.
	IsStructMemberContext()
}
    IStructMemberContext is an interface to support dynamic dispatch.

type ISubscriptDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SubscriptHead() ISubscriptHeadContext
	SubscriptResult() ISubscriptResultContext
	GetterSetterBlock() IGetterSetterBlockContext
	GetterSetterKeywordBlock() IGetterSetterKeywordBlockContext
	GenericWhereClause() IGenericWhereClauseContext

	// IsSubscriptDeclarationContext differentiates from other interfaces.
	IsSubscriptDeclarationContext()
}
    ISubscriptDeclarationContext is an interface to support dynamic dispatch.

type ISubscriptHeadContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SUBSCRIPT() antlr.TerminalNode
	ParameterClause() IParameterClauseContext
	Attributes() IAttributesContext
	DeclarationModifiers() IDeclarationModifiersContext
	GenericParameterClause() IGenericParameterClauseContext

	// IsSubscriptHeadContext differentiates from other interfaces.
	IsSubscriptHeadContext()
}
    ISubscriptHeadContext is an interface to support dynamic dispatch.

type ISubscriptResultContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ARROW() antlr.TerminalNode
	Type_() ITypeContext
	Attributes() IAttributesContext

	// IsSubscriptResultContext differentiates from other interfaces.
	IsSubscriptResultContext()
}
    ISubscriptResultContext is an interface to support dynamic dispatch.

type ISubscriptSuffixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	FunctionCallArgumentList() IFunctionCallArgumentListContext
	RBRACKET() antlr.TerminalNode

	// IsSubscriptSuffixContext differentiates from other interfaces.
	IsSubscriptSuffixContext()
}
    ISubscriptSuffixContext is an interface to support dynamic dispatch.

type ISuperExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SUPER() antlr.TerminalNode
	DOT() antlr.TerminalNode
	Identifier() IIdentifierContext
	LBRACKET() antlr.TerminalNode
	FunctionCallArgumentList() IFunctionCallArgumentListContext
	RBRACKET() antlr.TerminalNode
	INIT() antlr.TerminalNode

	// IsSuperExpressionContext differentiates from other interfaces.
	IsSuperExpressionContext()
}
    ISuperExpressionContext is an interface to support dynamic dispatch.

type ISuppressedTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATOR() antlr.TerminalNode
	TypeIdentifier() ITypeIdentifierContext

	// IsSuppressedTypeContext differentiates from other interfaces.
	IsSuppressedTypeContext()
}
    ISuppressedTypeContext is an interface to support dynamic dispatch.

type ISwitchCaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CaseLabel() ICaseLabelContext
	Statements() IStatementsContext
	DefaultLabel() IDefaultLabelContext

	// IsSwitchCaseContext differentiates from other interfaces.
	IsSwitchCaseContext()
}
    ISwitchCaseContext is an interface to support dynamic dispatch.

type ISwitchCasesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSwitchCase() []ISwitchCaseContext
	SwitchCase(i int) ISwitchCaseContext

	// IsSwitchCasesContext differentiates from other interfaces.
	IsSwitchCasesContext()
}
    ISwitchCasesContext is an interface to support dynamic dispatch.

type ISwitchStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SWITCH() antlr.TerminalNode
	Expression() IExpressionContext
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	SwitchCases() ISwitchCasesContext

	// IsSwitchStatementContext differentiates from other interfaces.
	IsSwitchStatementContext()
}
    ISwitchStatementContext is an interface to support dynamic dispatch.

type ITopLevelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	Statements() IStatementsContext

	// IsTopLevelContext differentiates from other interfaces.
	IsTopLevelContext()
}
    ITopLevelContext is an interface to support dynamic dispatch.

type ITrailingClosuresContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ClosureExpression() IClosureExpressionContext
	AllLabeledTrailingClosure() []ILabeledTrailingClosureContext
	LabeledTrailingClosure(i int) ILabeledTrailingClosureContext

	// IsTrailingClosuresContext differentiates from other interfaces.
	IsTrailingClosuresContext()
}
    ITrailingClosuresContext is an interface to support dynamic dispatch.

type ITryOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRY() antlr.TerminalNode
	EXCLAIM_POSTFIX() antlr.TerminalNode
	QUESTION_POSTFIX() antlr.TerminalNode

	// IsTryOperatorContext differentiates from other interfaces.
	IsTryOperatorContext()
}
    ITryOperatorContext is an interface to support dynamic dispatch.

type ITryStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRY() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	AllCatchClause() []ICatchClauseContext
	CatchClause(i int) ICatchClauseContext

	// IsTryStatementContext differentiates from other interfaces.
	IsTryStatementContext()
}
    ITryStatementContext is an interface to support dynamic dispatch.

type ITupleElementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	Identifier() IIdentifierContext
	COLON() antlr.TerminalNode

	// IsTupleElementContext differentiates from other interfaces.
	IsTupleElementContext()
}
    ITupleElementContext is an interface to support dynamic dispatch.

type ITupleElementListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTupleElement() []ITupleElementContext
	TupleElement(i int) ITupleElementContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTupleElementListContext differentiates from other interfaces.
	IsTupleElementListContext()
}
    ITupleElementListContext is an interface to support dynamic dispatch.

type ITupleExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	TupleElementList() ITupleElementListContext

	// IsTupleExpressionContext differentiates from other interfaces.
	IsTupleExpressionContext()
}
    ITupleExpressionContext is an interface to support dynamic dispatch.

type ITuplePatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	TuplePatternElementList() ITuplePatternElementListContext

	// IsTuplePatternContext differentiates from other interfaces.
	IsTuplePatternContext()
}
    ITuplePatternContext is an interface to support dynamic dispatch.

type ITuplePatternElementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Pattern() IPatternContext
	Identifier() IIdentifierContext
	COLON() antlr.TerminalNode

	// IsTuplePatternElementContext differentiates from other interfaces.
	IsTuplePatternElementContext()
}
    ITuplePatternElementContext is an interface to support dynamic dispatch.

type ITuplePatternElementListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTuplePatternElement() []ITuplePatternElementContext
	TuplePatternElement(i int) ITuplePatternElementContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTuplePatternElementListContext differentiates from other interfaces.
	IsTuplePatternElementListContext()
}
    ITuplePatternElementListContext is an interface to support dynamic dispatch.

type ITupleTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	TupleTypeElementList() ITupleTypeElementListContext

	// IsTupleTypeContext differentiates from other interfaces.
	IsTupleTypeContext()
}
    ITupleTypeContext is an interface to support dynamic dispatch.

type ITupleTypeElementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	TypeAnnotation() ITypeAnnotationContext
	Type_() ITypeContext

	// IsTupleTypeElementContext differentiates from other interfaces.
	IsTupleTypeElementContext()
}
    ITupleTypeElementContext is an interface to support dynamic dispatch.

type ITupleTypeElementListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTupleTypeElement() []ITupleTypeElementContext
	TupleTypeElement(i int) ITupleTypeElementContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTupleTypeElementListContext differentiates from other interfaces.
	IsTupleTypeElementListContext()
}
    ITupleTypeElementListContext is an interface to support dynamic dispatch.

type ITypeAnnotationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COLON() antlr.TerminalNode
	Type_() ITypeContext
	Attributes() IAttributesContext
	INOUT() antlr.TerminalNode

	// IsTypeAnnotationContext differentiates from other interfaces.
	IsTypeAnnotationContext()
}
    ITypeAnnotationContext is an interface to support dynamic dispatch.

type ITypeAnnotationHeadContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INOUT() antlr.TerminalNode
	Attributes() IAttributesContext

	// IsTypeAnnotationHeadContext differentiates from other interfaces.
	IsTypeAnnotationHeadContext()
}
    ITypeAnnotationHeadContext is an interface to support dynamic dispatch.

type ITypeCastingOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IS() antlr.TerminalNode
	Type_() ITypeContext
	AS() antlr.TerminalNode
	EXCLAIM_POSTFIX() antlr.TerminalNode
	QUESTION_POSTFIX() antlr.TerminalNode

	// IsTypeCastingOperatorContext differentiates from other interfaces.
	IsTypeCastingOperatorContext()
}
    ITypeCastingOperatorContext is an interface to support dynamic dispatch.

type ITypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsTypeContext differentiates from other interfaces.
	IsTypeContext()
}
    ITypeContext is an interface to support dynamic dispatch.

type ITypeIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	GenericArgumentClause() IGenericArgumentClauseContext
	DOT() antlr.TerminalNode
	TypeIdentifier() ITypeIdentifierContext

	// IsTypeIdentifierContext differentiates from other interfaces.
	IsTypeIdentifierContext()
}
    ITypeIdentifierContext is an interface to support dynamic dispatch.

type ITypeInheritanceClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COLON() antlr.TerminalNode
	TypeInheritanceList() ITypeInheritanceListContext

	// IsTypeInheritanceClauseContext differentiates from other interfaces.
	IsTypeInheritanceClauseContext()
}
    ITypeInheritanceClauseContext is an interface to support dynamic dispatch.

type ITypeInheritanceListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTypeIdentifier() []ITypeIdentifierContext
	TypeIdentifier(i int) ITypeIdentifierContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTypeInheritanceListContext differentiates from other interfaces.
	IsTypeInheritanceListContext()
}
    ITypeInheritanceListContext is an interface to support dynamic dispatch.

type ITypealiasDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPEALIAS() antlr.TerminalNode
	Identifier() IIdentifierContext
	ASSIGN() antlr.TerminalNode
	Type_() ITypeContext
	Attributes() IAttributesContext
	AccessLevelModifier() IAccessLevelModifierContext
	GenericParameterClause() IGenericParameterClauseContext

	// IsTypealiasDeclarationContext differentiates from other interfaces.
	IsTypealiasDeclarationContext()
}
    ITypealiasDeclarationContext is an interface to support dynamic dispatch.

type IUnionStyleEnumCaseClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	UnionStyleEnumCaseList() IUnionStyleEnumCaseListContext
	Attributes() IAttributesContext

	// IsUnionStyleEnumCaseClauseContext differentiates from other interfaces.
	IsUnionStyleEnumCaseClauseContext()
}
    IUnionStyleEnumCaseClauseContext is an interface to support dynamic
    dispatch.

type IUnionStyleEnumCaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	TupleType() ITupleTypeContext

	// IsUnionStyleEnumCaseContext differentiates from other interfaces.
	IsUnionStyleEnumCaseContext()
}
    IUnionStyleEnumCaseContext is an interface to support dynamic dispatch.

type IUnionStyleEnumCaseListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllUnionStyleEnumCase() []IUnionStyleEnumCaseContext
	UnionStyleEnumCase(i int) IUnionStyleEnumCaseContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsUnionStyleEnumCaseListContext differentiates from other interfaces.
	IsUnionStyleEnumCaseListContext()
}
    IUnionStyleEnumCaseListContext is an interface to support dynamic dispatch.

type IValueBindingPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Pattern() IPatternContext
	LET() antlr.TerminalNode
	VAR() antlr.TerminalNode
	INOUT() antlr.TerminalNode

	// IsValueBindingPatternContext differentiates from other interfaces.
	IsValueBindingPatternContext()
}
    IValueBindingPatternContext is an interface to support dynamic dispatch.

type IVariableDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VariableDeclarationHead() IVariableDeclarationHeadContext
	PatternInitializerList() IPatternInitializerListContext
	VariableName() IVariableNameContext
	TypeAnnotation() ITypeAnnotationContext
	CodeBlock() ICodeBlockContext
	GetterSetterBlock() IGetterSetterBlockContext
	GetterSetterKeywordBlock() IGetterSetterKeywordBlockContext
	WillSetDidSetBlock() IWillSetDidSetBlockContext
	Initializer() IInitializerContext

	// IsVariableDeclarationContext differentiates from other interfaces.
	IsVariableDeclarationContext()
}
    IVariableDeclarationContext is an interface to support dynamic dispatch.

type IVariableDeclarationHeadContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VAR() antlr.TerminalNode
	Attributes() IAttributesContext
	DeclarationModifiers() IDeclarationModifiersContext

	// IsVariableDeclarationHeadContext differentiates from other interfaces.
	IsVariableDeclarationHeadContext()
}
    IVariableDeclarationHeadContext is an interface to support dynamic dispatch.

type IVariableNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext

	// IsVariableNameContext differentiates from other interfaces.
	IsVariableNameContext()
}
    IVariableNameContext is an interface to support dynamic dispatch.

type IWhereClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHERE() antlr.TerminalNode
	Expression() IExpressionContext

	// IsWhereClauseContext differentiates from other interfaces.
	IsWhereClauseContext()
}
    IWhereClauseContext is an interface to support dynamic dispatch.

type IWhileStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHILE() antlr.TerminalNode
	ConditionList() IConditionListContext
	CodeBlock() ICodeBlockContext

	// IsWhileStatementContext differentiates from other interfaces.
	IsWhileStatementContext()
}
    IWhileStatementContext is an interface to support dynamic dispatch.

type IWildcardExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UNDERSCORE() antlr.TerminalNode

	// IsWildcardExpressionContext differentiates from other interfaces.
	IsWildcardExpressionContext()
}
    IWildcardExpressionContext is an interface to support dynamic dispatch.

type IWildcardPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UNDERSCORE() antlr.TerminalNode

	// IsWildcardPatternContext differentiates from other interfaces.
	IsWildcardPatternContext()
}
    IWildcardPatternContext is an interface to support dynamic dispatch.

type IWillSetClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WILLSET_KW() antlr.TerminalNode
	CodeBlock() ICodeBlockContext
	Attributes() IAttributesContext
	SetterName() ISetterNameContext

	// IsWillSetClauseContext differentiates from other interfaces.
	IsWillSetClauseContext()
}
    IWillSetClauseContext is an interface to support dynamic dispatch.

type IWillSetDidSetBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	WillSetClause() IWillSetClauseContext
	RBRACE() antlr.TerminalNode
	DidSetClause() IDidSetClauseContext

	// IsWillSetDidSetBlockContext differentiates from other interfaces.
	IsWillSetDidSetBlockContext()
}
    IWillSetDidSetBlockContext is an interface to support dynamic dispatch.

type IdentPatContext struct {
	PatternContext
}

func NewIdentPatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentPatContext

func (s *IdentPatContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IdentPatContext) GetRuleContext() antlr.RuleContext

func (s *IdentPatContext) IdentifierPattern() IIdentifierPatternContext

func (s *IdentPatContext) TypeAnnotation() ITypeAnnotationContext

type IdentifierContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyIdentifierContext() *IdentifierContext

func NewIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifierContext

func (s *IdentifierContext) ACTOR_KW() antlr.TerminalNode

func (s *IdentifierContext) ARCH_KW() antlr.TerminalNode

func (s *IdentifierContext) ASYNC_KW() antlr.TerminalNode

func (s *IdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IdentifierContext) BORROWING_KW() antlr.TerminalNode

func (s *IdentifierContext) CANIMPORT_KW() antlr.TerminalNode

func (s *IdentifierContext) COMPILER_KW() antlr.TerminalNode

func (s *IdentifierContext) CONSUMING_KW() antlr.TerminalNode

func (s *IdentifierContext) DIDSET_KW() antlr.TerminalNode

func (s *IdentifierContext) DYNAMIC_KW() antlr.TerminalNode

func (s *IdentifierContext) FILE_KEYWORD() antlr.TerminalNode

func (s *IdentifierContext) FINAL_KW() antlr.TerminalNode

func (s *IdentifierContext) GET_KW() antlr.TerminalNode

func (s *IdentifierContext) GetParser() antlr.Parser

func (s *IdentifierContext) GetRuleContext() antlr.RuleContext

func (s *IdentifierContext) IDENTIFIER() antlr.TerminalNode

func (*IdentifierContext) IsIdentifierContext()

func (s *IdentifierContext) LAZY_KW() antlr.TerminalNode

func (s *IdentifierContext) LINE_KEYWORD() antlr.TerminalNode

func (s *IdentifierContext) MACRO_KW() antlr.TerminalNode

func (s *IdentifierContext) MUTATING_KW() antlr.TerminalNode

func (s *IdentifierContext) NONISOLATED_KW() antlr.TerminalNode

func (s *IdentifierContext) NONMUTATING_KW() antlr.TerminalNode

func (s *IdentifierContext) OPTIONAL_KW() antlr.TerminalNode

func (s *IdentifierContext) OS_KW() antlr.TerminalNode

func (s *IdentifierContext) OVERRIDE_KW() antlr.TerminalNode

func (s *IdentifierContext) POSTFIX_KW() antlr.TerminalNode

func (s *IdentifierContext) PREFIX_KW() antlr.TerminalNode

func (s *IdentifierContext) PROTOCOL_KW() antlr.TerminalNode

func (s *IdentifierContext) REQUIRED_KW() antlr.TerminalNode

func (s *IdentifierContext) SENDING_KW() antlr.TerminalNode

func (s *IdentifierContext) SET_KW() antlr.TerminalNode

func (s *IdentifierContext) SOME_KW() antlr.TerminalNode

func (s *IdentifierContext) TYPE_KW() antlr.TerminalNode

func (s *IdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *IdentifierContext) UNOWNED_KW() antlr.TerminalNode

func (s *IdentifierContext) VERSION_KW() antlr.TerminalNode

func (s *IdentifierContext) VERTEX_KW() antlr.TerminalNode

func (s *IdentifierContext) WEAK_KW() antlr.TerminalNode

func (s *IdentifierContext) WILLSET_KW() antlr.TerminalNode

type IdentifierExprContext struct {
	PrimaryExpressionContext
}

func NewIdentifierExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentifierExprContext

func (s *IdentifierExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IdentifierExprContext) GenericArgumentClause() IGenericArgumentClauseContext

func (s *IdentifierExprContext) GetRuleContext() antlr.RuleContext

func (s *IdentifierExprContext) Identifier() IIdentifierContext

type IdentifierListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyIdentifierListContext() *IdentifierListContext

func NewIdentifierListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifierListContext

func (s *IdentifierListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IdentifierListContext) AllCOMMA() []antlr.TerminalNode

func (s *IdentifierListContext) AllIdentifier() []IIdentifierContext

func (s *IdentifierListContext) COMMA(i int) antlr.TerminalNode

func (s *IdentifierListContext) GetParser() antlr.Parser

func (s *IdentifierListContext) GetRuleContext() antlr.RuleContext

func (s *IdentifierListContext) Identifier(i int) IIdentifierContext

func (*IdentifierListContext) IsIdentifierListContext()

func (s *IdentifierListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type IdentifierPatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyIdentifierPatternContext() *IdentifierPatternContext

func NewIdentifierPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifierPatternContext

func (s *IdentifierPatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IdentifierPatternContext) GetParser() antlr.Parser

func (s *IdentifierPatternContext) GetRuleContext() antlr.RuleContext

func (s *IdentifierPatternContext) Identifier() IIdentifierContext

func (*IdentifierPatternContext) IsIdentifierPatternContext()

func (s *IdentifierPatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type IfStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyIfStatementContext() *IfStatementContext

func NewIfStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfStatementContext

func (s *IfStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IfStatementContext) CodeBlock() ICodeBlockContext

func (s *IfStatementContext) ConditionList() IConditionListContext

func (s *IfStatementContext) ElseClause() IElseClauseContext

func (s *IfStatementContext) GetParser() antlr.Parser

func (s *IfStatementContext) GetRuleContext() antlr.RuleContext

func (s *IfStatementContext) IF() antlr.TerminalNode

func (*IfStatementContext) IsIfStatementContext()

func (s *IfStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ImplicitMemberExprContext struct {
	PrimaryExpressionContext
}

func NewImplicitMemberExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ImplicitMemberExprContext

func (s *ImplicitMemberExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ImplicitMemberExprContext) GetRuleContext() antlr.RuleContext

func (s *ImplicitMemberExprContext) ImplicitMemberExpression() IImplicitMemberExpressionContext

type ImplicitMemberExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyImplicitMemberExpressionContext() *ImplicitMemberExpressionContext

func NewImplicitMemberExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImplicitMemberExpressionContext

func (s *ImplicitMemberExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ImplicitMemberExpressionContext) DOT() antlr.TerminalNode

func (s *ImplicitMemberExpressionContext) GenericArgumentClause() IGenericArgumentClauseContext

func (s *ImplicitMemberExpressionContext) GetParser() antlr.Parser

func (s *ImplicitMemberExpressionContext) GetRuleContext() antlr.RuleContext

func (s *ImplicitMemberExpressionContext) Identifier() IIdentifierContext

func (*ImplicitMemberExpressionContext) IsImplicitMemberExpressionContext()

func (s *ImplicitMemberExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ImportAliasContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyImportAliasContext() *ImportAliasContext

func NewImportAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportAliasContext

func (s *ImportAliasContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ImportAliasContext) DOT() antlr.TerminalNode

func (s *ImportAliasContext) GetParser() antlr.Parser

func (s *ImportAliasContext) GetRuleContext() antlr.RuleContext

func (s *ImportAliasContext) Identifier() IIdentifierContext

func (*ImportAliasContext) IsImportAliasContext()

func (s *ImportAliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ImportAliasContext) UNDERSCORE() antlr.TerminalNode

type ImportDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyImportDeclarationContext() *ImportDeclarationContext

func NewImportDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportDeclarationContext

func (s *ImportDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ImportDeclarationContext) AllImportSpec() []IImportSpecContext

func (s *ImportDeclarationContext) Attributes() IAttributesContext

func (s *ImportDeclarationContext) GetParser() antlr.Parser

func (s *ImportDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *ImportDeclarationContext) IMPORT() antlr.TerminalNode

func (s *ImportDeclarationContext) ImportSpec(i int) IImportSpecContext

func (*ImportDeclarationContext) IsImportDeclarationContext()

func (s *ImportDeclarationContext) LPAREN() antlr.TerminalNode

func (s *ImportDeclarationContext) RPAREN() antlr.TerminalNode

func (s *ImportDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ImportSpecContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyImportSpecContext() *ImportSpecContext

func NewImportSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportSpecContext

func (s *ImportSpecContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ImportSpecContext) GetParser() antlr.Parser

func (s *ImportSpecContext) GetRuleContext() antlr.RuleContext

func (s *ImportSpecContext) ImportAlias() IImportAliasContext

func (*ImportSpecContext) IsImportSpecContext()

func (s *ImportSpecContext) STRING_LITERAL() antlr.TerminalNode

func (s *ImportSpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type InOutExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyInOutExpressionContext() *InOutExpressionContext

func NewInOutExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InOutExpressionContext

func (s *InOutExpressionContext) AMPERSAND() antlr.TerminalNode

func (s *InOutExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *InOutExpressionContext) GetParser() antlr.Parser

func (s *InOutExpressionContext) GetRuleContext() antlr.RuleContext

func (s *InOutExpressionContext) Identifier() IIdentifierContext

func (*InOutExpressionContext) IsInOutExpressionContext()

func (s *InOutExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type InitializerBodyContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyInitializerBodyContext() *InitializerBodyContext

func NewInitializerBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InitializerBodyContext

func (s *InitializerBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *InitializerBodyContext) CodeBlock() ICodeBlockContext

func (s *InitializerBodyContext) GetParser() antlr.Parser

func (s *InitializerBodyContext) GetRuleContext() antlr.RuleContext

func (*InitializerBodyContext) IsInitializerBodyContext()

func (s *InitializerBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type InitializerContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyInitializerContext() *InitializerContext

func NewInitializerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InitializerContext

func (s *InitializerContext) ASSIGN() antlr.TerminalNode

func (s *InitializerContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *InitializerContext) Expression() IExpressionContext

func (s *InitializerContext) GetParser() antlr.Parser

func (s *InitializerContext) GetRuleContext() antlr.RuleContext

func (*InitializerContext) IsInitializerContext()

func (s *InitializerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type InitializerDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyInitializerDeclarationContext() *InitializerDeclarationContext

func NewInitializerDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InitializerDeclarationContext

func (s *InitializerDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *InitializerDeclarationContext) AsyncModifier() IAsyncModifierContext

func (s *InitializerDeclarationContext) Attributes() IAttributesContext

func (s *InitializerDeclarationContext) DeclarationModifiers() IDeclarationModifiersContext

func (s *InitializerDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *InitializerDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *InitializerDeclarationContext) GetParser() antlr.Parser

func (s *InitializerDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *InitializerDeclarationContext) INIT() antlr.TerminalNode

func (s *InitializerDeclarationContext) InitializerBody() IInitializerBodyContext

func (*InitializerDeclarationContext) IsInitializerDeclarationContext()

func (s *InitializerDeclarationContext) ParameterClause() IParameterClauseContext

func (s *InitializerDeclarationContext) QUESTION_POSTFIX() antlr.TerminalNode

func (s *InitializerDeclarationContext) THROWS() antlr.TerminalNode

func (s *InitializerDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type InitializerSuffixContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyInitializerSuffixContext() *InitializerSuffixContext

func NewInitializerSuffixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InitializerSuffixContext

func (s *InitializerSuffixContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *InitializerSuffixContext) ArgumentNames() IArgumentNamesContext

func (s *InitializerSuffixContext) DOT() antlr.TerminalNode

func (s *InitializerSuffixContext) GetParser() antlr.Parser

func (s *InitializerSuffixContext) GetRuleContext() antlr.RuleContext

func (s *InitializerSuffixContext) INIT() antlr.TerminalNode

func (*InitializerSuffixContext) IsInitializerSuffixContext()

func (s *InitializerSuffixContext) LPAREN() antlr.TerminalNode

func (s *InitializerSuffixContext) RPAREN() antlr.TerminalNode

func (s *InitializerSuffixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type InoutContext struct {
	PrefixExpressionContext
}

func NewInoutContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InoutContext

func (s *InoutContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *InoutContext) GetRuleContext() antlr.RuleContext

func (s *InoutContext) InOutExpression() IInOutExpressionContext

type IsPatContext struct {
	PatternContext
}

func NewIsPatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IsPatContext

func (s *IsPatContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IsPatContext) GetRuleContext() antlr.RuleContext

func (s *IsPatContext) IS() antlr.TerminalNode

func (s *IsPatContext) Type_() ITypeContext

type LabelNameContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLabelNameContext() *LabelNameContext

func NewLabelNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabelNameContext

func (s *LabelNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LabelNameContext) GetParser() antlr.Parser

func (s *LabelNameContext) GetRuleContext() antlr.RuleContext

func (s *LabelNameContext) Identifier() IIdentifierContext

func (*LabelNameContext) IsLabelNameContext()

func (s *LabelNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type LabeledStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLabeledStatementContext() *LabeledStatementContext

func NewLabeledStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabeledStatementContext

func (s *LabeledStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LabeledStatementContext) COLON() antlr.TerminalNode

func (s *LabeledStatementContext) DoStatement() IDoStatementContext

func (s *LabeledStatementContext) GetParser() antlr.Parser

func (s *LabeledStatementContext) GetRuleContext() antlr.RuleContext

func (s *LabeledStatementContext) IfStatement() IIfStatementContext

func (*LabeledStatementContext) IsLabeledStatementContext()

func (s *LabeledStatementContext) LabelName() ILabelNameContext

func (s *LabeledStatementContext) LoopStatement() ILoopStatementContext

func (s *LabeledStatementContext) SwitchStatement() ISwitchStatementContext

func (s *LabeledStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type LabeledStmtContext struct {
	StatementContext
}

func NewLabeledStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LabeledStmtContext

func (s *LabeledStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LabeledStmtContext) GetRuleContext() antlr.RuleContext

func (s *LabeledStmtContext) LabeledStatement() ILabeledStatementContext

func (s *LabeledStmtContext) SEMICOLON() antlr.TerminalNode

type LabeledTrailingClosureContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLabeledTrailingClosureContext() *LabeledTrailingClosureContext

func NewLabeledTrailingClosureContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabeledTrailingClosureContext

func (s *LabeledTrailingClosureContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LabeledTrailingClosureContext) COLON() antlr.TerminalNode

func (s *LabeledTrailingClosureContext) ClosureExpression() IClosureExpressionContext

func (s *LabeledTrailingClosureContext) GetParser() antlr.Parser

func (s *LabeledTrailingClosureContext) GetRuleContext() antlr.RuleContext

func (s *LabeledTrailingClosureContext) Identifier() IIdentifierContext

func (*LabeledTrailingClosureContext) IsLabeledTrailingClosureContext()

func (s *LabeledTrailingClosureContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type LayoutConstraintRequirementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLayoutConstraintRequirementContext() *LayoutConstraintRequirementContext

func NewLayoutConstraintRequirementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LayoutConstraintRequirementContext

func (s *LayoutConstraintRequirementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LayoutConstraintRequirementContext) COLON() antlr.TerminalNode

func (s *LayoutConstraintRequirementContext) GetParser() antlr.Parser

func (s *LayoutConstraintRequirementContext) GetRuleContext() antlr.RuleContext

func (s *LayoutConstraintRequirementContext) Identifier() IIdentifierContext

func (*LayoutConstraintRequirementContext) IsLayoutConstraintRequirementContext()

func (s *LayoutConstraintRequirementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *LayoutConstraintRequirementContext) TypeIdentifier() ITypeIdentifierContext

type LineControlStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLineControlStatementContext() *LineControlStatementContext

func NewLineControlStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LineControlStatementContext

func (s *LineControlStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LineControlStatementContext) AllCOLON() []antlr.TerminalNode

func (s *LineControlStatementContext) COLON(i int) antlr.TerminalNode

func (s *LineControlStatementContext) COMMA() antlr.TerminalNode

func (s *LineControlStatementContext) FILE_KEYWORD() antlr.TerminalNode

func (s *LineControlStatementContext) GetParser() antlr.Parser

func (s *LineControlStatementContext) GetRuleContext() antlr.RuleContext

func (s *LineControlStatementContext) INTEGER_LITERAL() antlr.TerminalNode

func (*LineControlStatementContext) IsLineControlStatementContext()

func (s *LineControlStatementContext) LINE_KEYWORD() antlr.TerminalNode

func (s *LineControlStatementContext) LPAREN() antlr.TerminalNode

func (s *LineControlStatementContext) POUND_SOURCE_LOCATION() antlr.TerminalNode

func (s *LineControlStatementContext) RPAREN() antlr.TerminalNode

func (s *LineControlStatementContext) STRING_LITERAL() antlr.TerminalNode

func (s *LineControlStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type LitExprContext struct {
	PrimaryExpressionContext
}

func NewLitExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LitExprContext

func (s *LitExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LitExprContext) GetRuleContext() antlr.RuleContext

func (s *LitExprContext) LiteralExpression() ILiteralExpressionContext

type LiteralContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLiteralContext() *LiteralContext

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext

func (s *LiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LiteralContext) BooleanLiteral() IBooleanLiteralContext

func (s *LiteralContext) EXTENDED_STRING_LITERAL() antlr.TerminalNode

func (s *LiteralContext) GetParser() antlr.Parser

func (s *LiteralContext) GetRuleContext() antlr.RuleContext

func (*LiteralContext) IsLiteralContext()

func (s *LiteralContext) MULTILINE_STRING_LITERAL() antlr.TerminalNode

func (s *LiteralContext) NIL() antlr.TerminalNode

func (s *LiteralContext) NumericLiteral() INumericLiteralContext

func (s *LiteralContext) STRING_LITERAL() antlr.TerminalNode

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type LiteralExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLiteralExpressionContext() *LiteralExpressionContext

func NewLiteralExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralExpressionContext

func (s *LiteralExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LiteralExpressionContext) ArrayLiteral() IArrayLiteralContext

func (s *LiteralExpressionContext) DictionaryLiteral() IDictionaryLiteralContext

func (s *LiteralExpressionContext) GetParser() antlr.Parser

func (s *LiteralExpressionContext) GetRuleContext() antlr.RuleContext

func (*LiteralExpressionContext) IsLiteralExpressionContext()

func (s *LiteralExpressionContext) Literal() ILiteralContext

func (s *LiteralExpressionContext) PoundFileExpression() IPoundFileExpressionContext

func (s *LiteralExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type LocalParameterNameContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLocalParameterNameContext() *LocalParameterNameContext

func NewLocalParameterNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LocalParameterNameContext

func (s *LocalParameterNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LocalParameterNameContext) GetParser() antlr.Parser

func (s *LocalParameterNameContext) GetRuleContext() antlr.RuleContext

func (s *LocalParameterNameContext) Identifier() IIdentifierContext

func (*LocalParameterNameContext) IsLocalParameterNameContext()

func (s *LocalParameterNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *LocalParameterNameContext) UNDERSCORE() antlr.TerminalNode

type LoopStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLoopStatementContext() *LoopStatementContext

func NewLoopStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LoopStatementContext

func (s *LoopStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LoopStatementContext) ForInStatement() IForInStatementContext

func (s *LoopStatementContext) ForStatement() IForStatementContext

func (s *LoopStatementContext) GetParser() antlr.Parser

func (s *LoopStatementContext) GetRuleContext() antlr.RuleContext

func (*LoopStatementContext) IsLoopStatementContext()

func (s *LoopStatementContext) RepeatWhileStatement() IRepeatWhileStatementContext

func (s *LoopStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *LoopStatementContext) WhileStatement() IWhileStatementContext

type LoopStmtContext struct {
	StatementContext
}

func NewLoopStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LoopStmtContext

func (s *LoopStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LoopStmtContext) GetRuleContext() antlr.RuleContext

func (s *LoopStmtContext) LoopStatement() ILoopStatementContext

func (s *LoopStmtContext) SEMICOLON() antlr.TerminalNode

type MacroDeclStmtContext struct {
	StatementContext
}

func NewMacroDeclStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MacroDeclStmtContext

func (s *MacroDeclStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MacroDeclStmtContext) GetRuleContext() antlr.RuleContext

func (s *MacroDeclStmtContext) MacroDeclaration() IMacroDeclarationContext

func (s *MacroDeclStmtContext) SEMICOLON() antlr.TerminalNode

type MacroDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyMacroDeclarationContext() *MacroDeclarationContext

func NewMacroDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MacroDeclarationContext

func (s *MacroDeclarationContext) ARROW() antlr.TerminalNode

func (s *MacroDeclarationContext) ASSIGN() antlr.TerminalNode

func (s *MacroDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MacroDeclarationContext) Attributes() IAttributesContext

func (s *MacroDeclarationContext) DeclarationModifiers() IDeclarationModifiersContext

func (s *MacroDeclarationContext) Expression() IExpressionContext

func (s *MacroDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *MacroDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *MacroDeclarationContext) GetParser() antlr.Parser

func (s *MacroDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *MacroDeclarationContext) Identifier() IIdentifierContext

func (*MacroDeclarationContext) IsMacroDeclarationContext()

func (s *MacroDeclarationContext) MACRO_KW() antlr.TerminalNode

func (s *MacroDeclarationContext) ParameterClause() IParameterClauseContext

func (s *MacroDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *MacroDeclarationContext) Type_() ITypeContext

type MacroExpansionExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyMacroExpansionExpressionContext() *MacroExpansionExpressionContext

func NewMacroExpansionExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MacroExpansionExpressionContext

func (s *MacroExpansionExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MacroExpansionExpressionContext) FunctionCallArgumentClause() IFunctionCallArgumentClauseContext

func (s *MacroExpansionExpressionContext) GenericArgumentClause() IGenericArgumentClauseContext

func (s *MacroExpansionExpressionContext) GetParser() antlr.Parser

func (s *MacroExpansionExpressionContext) GetRuleContext() antlr.RuleContext

func (s *MacroExpansionExpressionContext) HASH() antlr.TerminalNode

func (s *MacroExpansionExpressionContext) Identifier() IIdentifierContext

func (*MacroExpansionExpressionContext) IsMacroExpansionExpressionContext()

func (s *MacroExpansionExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *MacroExpansionExpressionContext) TrailingClosures() ITrailingClosuresContext

type MacroExprContext struct {
	PrimaryExpressionContext
}

func NewMacroExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MacroExprContext

func (s *MacroExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MacroExprContext) GetRuleContext() antlr.RuleContext

func (s *MacroExprContext) MacroExpansionExpression() IMacroExpansionExpressionContext

type MetatypeType_Context struct {
	TypeContext
}

func NewMetatypeType_Context(parser antlr.Parser, ctx antlr.ParserRuleContext) *MetatypeType_Context

func (s *MetatypeType_Context) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MetatypeType_Context) DOT() antlr.TerminalNode

func (s *MetatypeType_Context) GetRuleContext() antlr.RuleContext

func (s *MetatypeType_Context) PROTOCOL_KW() antlr.TerminalNode

func (s *MetatypeType_Context) TYPE_KW() antlr.TerminalNode

func (s *MetatypeType_Context) Type_() ITypeContext

type MutationModifierContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyMutationModifierContext() *MutationModifierContext

func NewMutationModifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MutationModifierContext

func (s *MutationModifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MutationModifierContext) GetParser() antlr.Parser

func (s *MutationModifierContext) GetRuleContext() antlr.RuleContext

func (*MutationModifierContext) IsMutationModifierContext()

func (s *MutationModifierContext) MUTATING_KW() antlr.TerminalNode

func (s *MutationModifierContext) NONMUTATING_KW() antlr.TerminalNode

func (s *MutationModifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type NamedTypeContext struct {
	TypeContext
}

func NewNamedTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NamedTypeContext

func (s *NamedTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *NamedTypeContext) GetRuleContext() antlr.RuleContext

func (s *NamedTypeContext) TypeIdentifier() ITypeIdentifierContext

type NumericLiteralContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyNumericLiteralContext() *NumericLiteralContext

func NewNumericLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumericLiteralContext

func (s *NumericLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *NumericLiteralContext) FLOAT_LITERAL() antlr.TerminalNode

func (s *NumericLiteralContext) GetParser() antlr.Parser

func (s *NumericLiteralContext) GetRuleContext() antlr.RuleContext

func (s *NumericLiteralContext) INTEGER_LITERAL() antlr.TerminalNode

func (*NumericLiteralContext) IsNumericLiteralContext()

func (s *NumericLiteralContext) OPERATOR() antlr.TerminalNode

func (s *NumericLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type OpaqueTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyOpaqueTypeContext() *OpaqueTypeContext

func NewOpaqueTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OpaqueTypeContext

func (s *OpaqueTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *OpaqueTypeContext) GetParser() antlr.Parser

func (s *OpaqueTypeContext) GetRuleContext() antlr.RuleContext

func (*OpaqueTypeContext) IsOpaqueTypeContext()

func (s *OpaqueTypeContext) SOME_KW() antlr.TerminalNode

func (s *OpaqueTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *OpaqueTypeContext) Type_() ITypeContext

type OpaqueType_Context struct {
	TypeContext
}

func NewOpaqueType_Context(parser antlr.Parser, ctx antlr.ParserRuleContext) *OpaqueType_Context

func (s *OpaqueType_Context) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *OpaqueType_Context) GetRuleContext() antlr.RuleContext

func (s *OpaqueType_Context) OpaqueType() IOpaqueTypeContext

type Operator_Context struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyOperator_Context() *Operator_Context

func NewOperator_Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Operator_Context

func (s *Operator_Context) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *Operator_Context) DOT() antlr.TerminalNode

func (s *Operator_Context) GetParser() antlr.Parser

func (s *Operator_Context) GetRuleContext() antlr.RuleContext

func (*Operator_Context) IsOperator_Context()

func (s *Operator_Context) OPERATOR() antlr.TerminalNode

func (s *Operator_Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type OptPatContext struct {
	PatternContext
}

func NewOptPatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *OptPatContext

func (s *OptPatContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *OptPatContext) GetRuleContext() antlr.RuleContext

func (s *OptPatContext) OptionalPattern() IOptionalPatternContext

type OptTypeContext struct {
	TypeContext
}

func NewOptTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *OptTypeContext

func (s *OptTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *OptTypeContext) GetRuleContext() antlr.RuleContext

func (s *OptTypeContext) QUESTION_POSTFIX() antlr.TerminalNode

func (s *OptTypeContext) Type_() ITypeContext

type OptionalBindingConditionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyOptionalBindingConditionContext() *OptionalBindingConditionContext

func NewOptionalBindingConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptionalBindingConditionContext

func (s *OptionalBindingConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *OptionalBindingConditionContext) GetParser() antlr.Parser

func (s *OptionalBindingConditionContext) GetRuleContext() antlr.RuleContext

func (s *OptionalBindingConditionContext) Initializer() IInitializerContext

func (*OptionalBindingConditionContext) IsOptionalBindingConditionContext()

func (s *OptionalBindingConditionContext) LET() antlr.TerminalNode

func (s *OptionalBindingConditionContext) Pattern() IPatternContext

func (s *OptionalBindingConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *OptionalBindingConditionContext) VAR() antlr.TerminalNode

type OptionalChainingLiteralContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyOptionalChainingLiteralContext() *OptionalChainingLiteralContext

func NewOptionalChainingLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptionalChainingLiteralContext

func (s *OptionalChainingLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *OptionalChainingLiteralContext) GetParser() antlr.Parser

func (s *OptionalChainingLiteralContext) GetRuleContext() antlr.RuleContext

func (*OptionalChainingLiteralContext) IsOptionalChainingLiteralContext()

func (s *OptionalChainingLiteralContext) QUESTION_POSTFIX() antlr.TerminalNode

func (s *OptionalChainingLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type OptionalPatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyOptionalPatternContext() *OptionalPatternContext

func NewOptionalPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptionalPatternContext

func (s *OptionalPatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *OptionalPatternContext) GetParser() antlr.Parser

func (s *OptionalPatternContext) GetRuleContext() antlr.RuleContext

func (s *OptionalPatternContext) IdentifierPattern() IIdentifierPatternContext

func (*OptionalPatternContext) IsOptionalPatternContext()

func (s *OptionalPatternContext) QUESTION_POSTFIX() antlr.TerminalNode

func (s *OptionalPatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ParameterClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyParameterClauseContext() *ParameterClauseContext

func NewParameterClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParameterClauseContext

func (s *ParameterClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParameterClauseContext) GetParser() antlr.Parser

func (s *ParameterClauseContext) GetRuleContext() antlr.RuleContext

func (*ParameterClauseContext) IsParameterClauseContext()

func (s *ParameterClauseContext) LPAREN() antlr.TerminalNode

func (s *ParameterClauseContext) ParameterList() IParameterListContext

func (s *ParameterClauseContext) RPAREN() antlr.TerminalNode

func (s *ParameterClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ParameterContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyParameterContext() *ParameterContext

func NewParameterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParameterContext

func (s *ParameterContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParameterContext) DefaultArgumentClause() IDefaultArgumentClauseContext

func (s *ParameterContext) ELLIPSIS() antlr.TerminalNode

func (s *ParameterContext) ExternalParameterName() IExternalParameterNameContext

func (s *ParameterContext) GetParser() antlr.Parser

func (s *ParameterContext) GetRuleContext() antlr.RuleContext

func (*ParameterContext) IsParameterContext()

func (s *ParameterContext) LocalParameterName() ILocalParameterNameContext

func (s *ParameterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ParameterContext) TypeAnnotation() ITypeAnnotationContext

type ParameterListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyParameterListContext() *ParameterListContext

func NewParameterListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParameterListContext

func (s *ParameterListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParameterListContext) AllCOMMA() []antlr.TerminalNode

func (s *ParameterListContext) AllParameter() []IParameterContext

func (s *ParameterListContext) COMMA(i int) antlr.TerminalNode

func (s *ParameterListContext) GetParser() antlr.Parser

func (s *ParameterListContext) GetRuleContext() antlr.RuleContext

func (*ParameterListContext) IsParameterListContext()

func (s *ParameterListContext) Parameter(i int) IParameterContext

func (s *ParameterListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ParenExprContext struct {
	PrimaryExpressionContext
}

func NewParenExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenExprContext

func (s *ParenExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParenExprContext) GetRuleContext() antlr.RuleContext

func (s *ParenExprContext) ParenthesizedExpression() IParenthesizedExpressionContext

type ParenthesizedExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyParenthesizedExpressionContext() *ParenthesizedExpressionContext

func NewParenthesizedExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParenthesizedExpressionContext

func (s *ParenthesizedExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParenthesizedExpressionContext) Expression() IExpressionContext

func (s *ParenthesizedExpressionContext) GetParser() antlr.Parser

func (s *ParenthesizedExpressionContext) GetRuleContext() antlr.RuleContext

func (*ParenthesizedExpressionContext) IsParenthesizedExpressionContext()

func (s *ParenthesizedExpressionContext) LPAREN() antlr.TerminalNode

func (s *ParenthesizedExpressionContext) RPAREN() antlr.TerminalNode

func (s *ParenthesizedExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPatternContext() *PatternContext

func NewPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PatternContext

func (s *PatternContext) CopyAll(ctx *PatternContext)

func (s *PatternContext) GetParser() antlr.Parser

func (s *PatternContext) GetRuleContext() antlr.RuleContext

func (*PatternContext) IsPatternContext()

func (s *PatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PatternInitializerContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPatternInitializerContext() *PatternInitializerContext

func NewPatternInitializerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PatternInitializerContext

func (s *PatternInitializerContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PatternInitializerContext) GetParser() antlr.Parser

func (s *PatternInitializerContext) GetRuleContext() antlr.RuleContext

func (s *PatternInitializerContext) Initializer() IInitializerContext

func (*PatternInitializerContext) IsPatternInitializerContext()

func (s *PatternInitializerContext) Pattern() IPatternContext

func (s *PatternInitializerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PatternInitializerListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPatternInitializerListContext() *PatternInitializerListContext

func NewPatternInitializerListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PatternInitializerListContext

func (s *PatternInitializerListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PatternInitializerListContext) AllCOMMA() []antlr.TerminalNode

func (s *PatternInitializerListContext) AllPatternInitializer() []IPatternInitializerContext

func (s *PatternInitializerListContext) COMMA(i int) antlr.TerminalNode

func (s *PatternInitializerListContext) GetParser() antlr.Parser

func (s *PatternInitializerListContext) GetRuleContext() antlr.RuleContext

func (*PatternInitializerListContext) IsPatternInitializerListContext()

func (s *PatternInitializerListContext) PatternInitializer(i int) IPatternInitializerContext

func (s *PatternInitializerListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PlatformConditionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPlatformConditionContext() *PlatformConditionContext

func NewPlatformConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PlatformConditionContext

func (s *PlatformConditionContext) ARCH_KW() antlr.TerminalNode

func (s *PlatformConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PlatformConditionContext) CANIMPORT_KW() antlr.TerminalNode

func (s *PlatformConditionContext) COLON() antlr.TerminalNode

func (s *PlatformConditionContext) COMMA() antlr.TerminalNode

func (s *PlatformConditionContext) COMPILER_KW() antlr.TerminalNode

func (s *PlatformConditionContext) DecimalVersion() IDecimalVersionContext

func (s *PlatformConditionContext) GetParser() antlr.Parser

func (s *PlatformConditionContext) GetRuleContext() antlr.RuleContext

func (s *PlatformConditionContext) Identifier() IIdentifierContext

func (*PlatformConditionContext) IsPlatformConditionContext()

func (s *PlatformConditionContext) LPAREN() antlr.TerminalNode

func (s *PlatformConditionContext) OPERATOR() antlr.TerminalNode

func (s *PlatformConditionContext) OS_KW() antlr.TerminalNode

func (s *PlatformConditionContext) RPAREN() antlr.TerminalNode

func (s *PlatformConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *PlatformConditionContext) VERSION_KW() antlr.TerminalNode

func (s *PlatformConditionContext) VERTEX_KW() antlr.TerminalNode

type PostfixExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPostfixExpressionContext() *PostfixExpressionContext

func NewPostfixExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PostfixExpressionContext

func (s *PostfixExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PostfixExpressionContext) AllPostfixSuffix() []IPostfixSuffixContext

func (s *PostfixExpressionContext) GetParser() antlr.Parser

func (s *PostfixExpressionContext) GetRuleContext() antlr.RuleContext

func (*PostfixExpressionContext) IsPostfixExpressionContext()

func (s *PostfixExpressionContext) PostfixSuffix(i int) IPostfixSuffixContext

func (s *PostfixExpressionContext) PrimaryExpression() IPrimaryExpressionContext

func (s *PostfixExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PostfixOperatorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPostfixOperatorContext() *PostfixOperatorContext

func NewPostfixOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PostfixOperatorContext

func (s *PostfixOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PostfixOperatorContext) GetParser() antlr.Parser

func (s *PostfixOperatorContext) GetRuleContext() antlr.RuleContext

func (*PostfixOperatorContext) IsPostfixOperatorContext()

func (s *PostfixOperatorContext) OPERATOR() antlr.TerminalNode

func (s *PostfixOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PostfixSelfSuffixContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPostfixSelfSuffixContext() *PostfixSelfSuffixContext

func NewPostfixSelfSuffixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PostfixSelfSuffixContext

func (s *PostfixSelfSuffixContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PostfixSelfSuffixContext) DOT() antlr.TerminalNode

func (s *PostfixSelfSuffixContext) GetParser() antlr.Parser

func (s *PostfixSelfSuffixContext) GetRuleContext() antlr.RuleContext

func (*PostfixSelfSuffixContext) IsPostfixSelfSuffixContext()

func (s *PostfixSelfSuffixContext) SELF() antlr.TerminalNode

func (s *PostfixSelfSuffixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PostfixSuffixContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPostfixSuffixContext() *PostfixSuffixContext

func NewPostfixSuffixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PostfixSuffixContext

func (s *PostfixSuffixContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PostfixSuffixContext) ExplicitMemberSuffix() IExplicitMemberSuffixContext

func (s *PostfixSuffixContext) ForcedValueSuffix() IForcedValueSuffixContext

func (s *PostfixSuffixContext) FunctionCallSuffix() IFunctionCallSuffixContext

func (s *PostfixSuffixContext) GetParser() antlr.Parser

func (s *PostfixSuffixContext) GetRuleContext() antlr.RuleContext

func (s *PostfixSuffixContext) InitializerSuffix() IInitializerSuffixContext

func (*PostfixSuffixContext) IsPostfixSuffixContext()

func (s *PostfixSuffixContext) OptionalChainingLiteral() IOptionalChainingLiteralContext

func (s *PostfixSuffixContext) PostfixOperator() IPostfixOperatorContext

func (s *PostfixSuffixContext) PostfixSelfSuffix() IPostfixSelfSuffixContext

func (s *PostfixSuffixContext) SubscriptSuffix() ISubscriptSuffixContext

func (s *PostfixSuffixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PoundFileExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPoundFileExpressionContext() *PoundFileExpressionContext

func NewPoundFileExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PoundFileExpressionContext

func (s *PoundFileExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PoundFileExpressionContext) GetParser() antlr.Parser

func (s *PoundFileExpressionContext) GetRuleContext() antlr.RuleContext

func (*PoundFileExpressionContext) IsPoundFileExpressionContext()

func (s *PoundFileExpressionContext) POUND_COLUMN() antlr.TerminalNode

func (s *PoundFileExpressionContext) POUND_FILE() antlr.TerminalNode

func (s *PoundFileExpressionContext) POUND_FILEID() antlr.TerminalNode

func (s *PoundFileExpressionContext) POUND_FILEPATH() antlr.TerminalNode

func (s *PoundFileExpressionContext) POUND_FUNCTION() antlr.TerminalNode

func (s *PoundFileExpressionContext) POUND_LINE() antlr.TerminalNode

func (s *PoundFileExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PrefixExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPrefixExpressionContext() *PrefixExpressionContext

func NewPrefixExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrefixExpressionContext

func (s *PrefixExpressionContext) CopyAll(ctx *PrefixExpressionContext)

func (s *PrefixExpressionContext) GetParser() antlr.Parser

func (s *PrefixExpressionContext) GetRuleContext() antlr.RuleContext

func (*PrefixExpressionContext) IsPrefixExpressionContext()

func (s *PrefixExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PrefixOpContext struct {
	PrefixExpressionContext
}

func NewPrefixOpContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PrefixOpContext

func (s *PrefixOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PrefixOpContext) GetRuleContext() antlr.RuleContext

func (s *PrefixOpContext) PostfixExpression() IPostfixExpressionContext

func (s *PrefixOpContext) PrefixOperator() IPrefixOperatorContext

type PrefixOperatorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPrefixOperatorContext() *PrefixOperatorContext

func NewPrefixOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrefixOperatorContext

func (s *PrefixOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PrefixOperatorContext) GetParser() antlr.Parser

func (s *PrefixOperatorContext) GetRuleContext() antlr.RuleContext

func (*PrefixOperatorContext) IsPrefixOperatorContext()

func (s *PrefixOperatorContext) OPERATOR() antlr.TerminalNode

func (s *PrefixOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PrimaryAssociatedTypeClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPrimaryAssociatedTypeClauseContext() *PrimaryAssociatedTypeClauseContext

func NewPrimaryAssociatedTypeClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryAssociatedTypeClauseContext

func (s *PrimaryAssociatedTypeClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PrimaryAssociatedTypeClauseContext) GT() antlr.TerminalNode

func (s *PrimaryAssociatedTypeClauseContext) GetParser() antlr.Parser

func (s *PrimaryAssociatedTypeClauseContext) GetRuleContext() antlr.RuleContext

func (*PrimaryAssociatedTypeClauseContext) IsPrimaryAssociatedTypeClauseContext()

func (s *PrimaryAssociatedTypeClauseContext) LT() antlr.TerminalNode

func (s *PrimaryAssociatedTypeClauseContext) PrimaryAssociatedTypeList() IPrimaryAssociatedTypeListContext

func (s *PrimaryAssociatedTypeClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PrimaryAssociatedTypeListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPrimaryAssociatedTypeListContext() *PrimaryAssociatedTypeListContext

func NewPrimaryAssociatedTypeListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryAssociatedTypeListContext

func (s *PrimaryAssociatedTypeListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PrimaryAssociatedTypeListContext) AllCOMMA() []antlr.TerminalNode

func (s *PrimaryAssociatedTypeListContext) AllIdentifier() []IIdentifierContext

func (s *PrimaryAssociatedTypeListContext) COMMA(i int) antlr.TerminalNode

func (s *PrimaryAssociatedTypeListContext) GetParser() antlr.Parser

func (s *PrimaryAssociatedTypeListContext) GetRuleContext() antlr.RuleContext

func (s *PrimaryAssociatedTypeListContext) Identifier(i int) IIdentifierContext

func (*PrimaryAssociatedTypeListContext) IsPrimaryAssociatedTypeListContext()

func (s *PrimaryAssociatedTypeListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PrimaryExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPrimaryExpressionContext() *PrimaryExpressionContext

func NewPrimaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryExpressionContext

func (s *PrimaryExpressionContext) CopyAll(ctx *PrimaryExpressionContext)

func (s *PrimaryExpressionContext) GetParser() antlr.Parser

func (s *PrimaryExpressionContext) GetRuleContext() antlr.RuleContext

func (*PrimaryExpressionContext) IsPrimaryExpressionContext()

func (s *PrimaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ProtoCompTypeContext struct {
	TypeContext
}

func NewProtoCompTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ProtoCompTypeContext

func (s *ProtoCompTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtoCompTypeContext) GetRuleContext() antlr.RuleContext

func (s *ProtoCompTypeContext) ProtocolCompositionType() IProtocolCompositionTypeContext

type ProtocolAssociatedTypeDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolAssociatedTypeDeclarationContext() *ProtocolAssociatedTypeDeclarationContext

func NewProtocolAssociatedTypeDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolAssociatedTypeDeclarationContext

func (s *ProtocolAssociatedTypeDeclarationContext) ASSIGN() antlr.TerminalNode

func (s *ProtocolAssociatedTypeDeclarationContext) ASSOCIATEDTYPE() antlr.TerminalNode

func (s *ProtocolAssociatedTypeDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolAssociatedTypeDeclarationContext) AccessLevelModifier() IAccessLevelModifierContext

func (s *ProtocolAssociatedTypeDeclarationContext) Attributes() IAttributesContext

func (s *ProtocolAssociatedTypeDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *ProtocolAssociatedTypeDeclarationContext) GetParser() antlr.Parser

func (s *ProtocolAssociatedTypeDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *ProtocolAssociatedTypeDeclarationContext) Identifier() IIdentifierContext

func (*ProtocolAssociatedTypeDeclarationContext) IsProtocolAssociatedTypeDeclarationContext()

func (s *ProtocolAssociatedTypeDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ProtocolAssociatedTypeDeclarationContext) TypeInheritanceClause() ITypeInheritanceClauseContext

func (s *ProtocolAssociatedTypeDeclarationContext) Type_() ITypeContext

type ProtocolBodyContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolBodyContext() *ProtocolBodyContext

func NewProtocolBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolBodyContext

func (s *ProtocolBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolBodyContext) AllProtocolMember() []IProtocolMemberContext

func (s *ProtocolBodyContext) GetParser() antlr.Parser

func (s *ProtocolBodyContext) GetRuleContext() antlr.RuleContext

func (*ProtocolBodyContext) IsProtocolBodyContext()

func (s *ProtocolBodyContext) LBRACE() antlr.TerminalNode

func (s *ProtocolBodyContext) ProtocolMember(i int) IProtocolMemberContext

func (s *ProtocolBodyContext) RBRACE() antlr.TerminalNode

func (s *ProtocolBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ProtocolCompositionTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolCompositionTypeContext() *ProtocolCompositionTypeContext

func NewProtocolCompositionTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolCompositionTypeContext

func (s *ProtocolCompositionTypeContext) AMPERSAND(i int) antlr.TerminalNode

func (s *ProtocolCompositionTypeContext) ANY() antlr.TerminalNode

func (s *ProtocolCompositionTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolCompositionTypeContext) AllAMPERSAND() []antlr.TerminalNode

func (s *ProtocolCompositionTypeContext) AllProtocolCompositionTypeElement() []IProtocolCompositionTypeElementContext

func (s *ProtocolCompositionTypeContext) GetParser() antlr.Parser

func (s *ProtocolCompositionTypeContext) GetRuleContext() antlr.RuleContext

func (*ProtocolCompositionTypeContext) IsProtocolCompositionTypeContext()

func (s *ProtocolCompositionTypeContext) ProtocolCompositionTypeElement(i int) IProtocolCompositionTypeElementContext

func (s *ProtocolCompositionTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ProtocolCompositionTypeElementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolCompositionTypeElementContext() *ProtocolCompositionTypeElementContext

func NewProtocolCompositionTypeElementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolCompositionTypeElementContext

func (s *ProtocolCompositionTypeElementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolCompositionTypeElementContext) GetParser() antlr.Parser

func (s *ProtocolCompositionTypeElementContext) GetRuleContext() antlr.RuleContext

func (*ProtocolCompositionTypeElementContext) IsProtocolCompositionTypeElementContext()

func (s *ProtocolCompositionTypeElementContext) SuppressedType() ISuppressedTypeContext

func (s *ProtocolCompositionTypeElementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ProtocolCompositionTypeElementContext) TypeIdentifier() ITypeIdentifierContext

type ProtocolDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolDeclarationContext() *ProtocolDeclarationContext

func NewProtocolDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolDeclarationContext

func (s *ProtocolDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolDeclarationContext) AccessLevelModifier() IAccessLevelModifierContext

func (s *ProtocolDeclarationContext) Attributes() IAttributesContext

func (s *ProtocolDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *ProtocolDeclarationContext) GetParser() antlr.Parser

func (s *ProtocolDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *ProtocolDeclarationContext) Identifier() IIdentifierContext

func (*ProtocolDeclarationContext) IsProtocolDeclarationContext()

func (s *ProtocolDeclarationContext) PROTOCOL() antlr.TerminalNode

func (s *ProtocolDeclarationContext) PrimaryAssociatedTypeClause() IPrimaryAssociatedTypeClauseContext

func (s *ProtocolDeclarationContext) ProtocolBody() IProtocolBodyContext

func (s *ProtocolDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ProtocolDeclarationContext) TypeInheritanceClause() ITypeInheritanceClauseContext

type ProtocolInitializerDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolInitializerDeclarationContext() *ProtocolInitializerDeclarationContext

func NewProtocolInitializerDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolInitializerDeclarationContext

func (s *ProtocolInitializerDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolInitializerDeclarationContext) Attributes() IAttributesContext

func (s *ProtocolInitializerDeclarationContext) DeclarationModifiers() IDeclarationModifiersContext

func (s *ProtocolInitializerDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *ProtocolInitializerDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *ProtocolInitializerDeclarationContext) GetParser() antlr.Parser

func (s *ProtocolInitializerDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *ProtocolInitializerDeclarationContext) INIT() antlr.TerminalNode

func (*ProtocolInitializerDeclarationContext) IsProtocolInitializerDeclarationContext()

func (s *ProtocolInitializerDeclarationContext) ParameterClause() IParameterClauseContext

func (s *ProtocolInitializerDeclarationContext) QUESTION_POSTFIX() antlr.TerminalNode

func (s *ProtocolInitializerDeclarationContext) THROWS() antlr.TerminalNode

func (s *ProtocolInitializerDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ProtocolMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolMemberContext() *ProtocolMemberContext

func NewProtocolMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolMemberContext

func (s *ProtocolMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolMemberContext) CompilerControl() ICompilerControlContext

func (s *ProtocolMemberContext) GetParser() antlr.Parser

func (s *ProtocolMemberContext) GetRuleContext() antlr.RuleContext

func (*ProtocolMemberContext) IsProtocolMemberContext()

func (s *ProtocolMemberContext) ProtocolMemberDeclaration() IProtocolMemberDeclarationContext

func (s *ProtocolMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ProtocolMemberDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolMemberDeclarationContext() *ProtocolMemberDeclarationContext

func NewProtocolMemberDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolMemberDeclarationContext

func (s *ProtocolMemberDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolMemberDeclarationContext) GetParser() antlr.Parser

func (s *ProtocolMemberDeclarationContext) GetRuleContext() antlr.RuleContext

func (*ProtocolMemberDeclarationContext) IsProtocolMemberDeclarationContext()

func (s *ProtocolMemberDeclarationContext) ProtocolAssociatedTypeDeclaration() IProtocolAssociatedTypeDeclarationContext

func (s *ProtocolMemberDeclarationContext) ProtocolInitializerDeclaration() IProtocolInitializerDeclarationContext

func (s *ProtocolMemberDeclarationContext) ProtocolMethodDeclaration() IProtocolMethodDeclarationContext

func (s *ProtocolMemberDeclarationContext) ProtocolPropertyDeclaration() IProtocolPropertyDeclarationContext

func (s *ProtocolMemberDeclarationContext) ProtocolSubscriptDeclaration() IProtocolSubscriptDeclarationContext

func (s *ProtocolMemberDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ProtocolMemberDeclarationContext) TypealiasDeclaration() ITypealiasDeclarationContext

type ProtocolMethodDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolMethodDeclarationContext() *ProtocolMethodDeclarationContext

func NewProtocolMethodDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolMethodDeclarationContext

func (s *ProtocolMethodDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolMethodDeclarationContext) FunctionHead() IFunctionHeadContext

func (s *ProtocolMethodDeclarationContext) FunctionSignature() IFunctionSignatureContext

func (s *ProtocolMethodDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *ProtocolMethodDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *ProtocolMethodDeclarationContext) GetParser() antlr.Parser

func (s *ProtocolMethodDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *ProtocolMethodDeclarationContext) Identifier() IIdentifierContext

func (*ProtocolMethodDeclarationContext) IsProtocolMethodDeclarationContext()

func (s *ProtocolMethodDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ProtocolPropertyDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolPropertyDeclarationContext() *ProtocolPropertyDeclarationContext

func NewProtocolPropertyDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolPropertyDeclarationContext

func (s *ProtocolPropertyDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolPropertyDeclarationContext) GetParser() antlr.Parser

func (s *ProtocolPropertyDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *ProtocolPropertyDeclarationContext) GetterSetterKeywordBlock() IGetterSetterKeywordBlockContext

func (*ProtocolPropertyDeclarationContext) IsProtocolPropertyDeclarationContext()

func (s *ProtocolPropertyDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ProtocolPropertyDeclarationContext) TypeAnnotation() ITypeAnnotationContext

func (s *ProtocolPropertyDeclarationContext) VariableDeclarationHead() IVariableDeclarationHeadContext

func (s *ProtocolPropertyDeclarationContext) VariableName() IVariableNameContext

type ProtocolSubscriptDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyProtocolSubscriptDeclarationContext() *ProtocolSubscriptDeclarationContext

func NewProtocolSubscriptDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProtocolSubscriptDeclarationContext

func (s *ProtocolSubscriptDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProtocolSubscriptDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *ProtocolSubscriptDeclarationContext) GetParser() antlr.Parser

func (s *ProtocolSubscriptDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *ProtocolSubscriptDeclarationContext) GetterSetterKeywordBlock() IGetterSetterKeywordBlockContext

func (*ProtocolSubscriptDeclarationContext) IsProtocolSubscriptDeclarationContext()

func (s *ProtocolSubscriptDeclarationContext) SubscriptHead() ISubscriptHeadContext

func (s *ProtocolSubscriptDeclarationContext) SubscriptResult() ISubscriptResultContext

func (s *ProtocolSubscriptDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type RawValueAssignmentContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyRawValueAssignmentContext() *RawValueAssignmentContext

func NewRawValueAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RawValueAssignmentContext

func (s *RawValueAssignmentContext) ASSIGN() antlr.TerminalNode

func (s *RawValueAssignmentContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RawValueAssignmentContext) GetParser() antlr.Parser

func (s *RawValueAssignmentContext) GetRuleContext() antlr.RuleContext

func (*RawValueAssignmentContext) IsRawValueAssignmentContext()

func (s *RawValueAssignmentContext) RawValueLiteral() IRawValueLiteralContext

func (s *RawValueAssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type RawValueLiteralContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyRawValueLiteralContext() *RawValueLiteralContext

func NewRawValueLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RawValueLiteralContext

func (s *RawValueLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RawValueLiteralContext) BooleanLiteral() IBooleanLiteralContext

func (s *RawValueLiteralContext) GetParser() antlr.Parser

func (s *RawValueLiteralContext) GetRuleContext() antlr.RuleContext

func (*RawValueLiteralContext) IsRawValueLiteralContext()

func (s *RawValueLiteralContext) NumericLiteral() INumericLiteralContext

func (s *RawValueLiteralContext) STRING_LITERAL() antlr.TerminalNode

func (s *RawValueLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type RawValueStyleEnumCaseClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyRawValueStyleEnumCaseClauseContext() *RawValueStyleEnumCaseClauseContext

func NewRawValueStyleEnumCaseClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RawValueStyleEnumCaseClauseContext

func (s *RawValueStyleEnumCaseClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RawValueStyleEnumCaseClauseContext) Attributes() IAttributesContext

func (s *RawValueStyleEnumCaseClauseContext) CASE() antlr.TerminalNode

func (s *RawValueStyleEnumCaseClauseContext) GetParser() antlr.Parser

func (s *RawValueStyleEnumCaseClauseContext) GetRuleContext() antlr.RuleContext

func (*RawValueStyleEnumCaseClauseContext) IsRawValueStyleEnumCaseClauseContext()

func (s *RawValueStyleEnumCaseClauseContext) RawValueStyleEnumCaseList() IRawValueStyleEnumCaseListContext

func (s *RawValueStyleEnumCaseClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type RawValueStyleEnumCaseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyRawValueStyleEnumCaseContext() *RawValueStyleEnumCaseContext

func NewRawValueStyleEnumCaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RawValueStyleEnumCaseContext

func (s *RawValueStyleEnumCaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RawValueStyleEnumCaseContext) GetParser() antlr.Parser

func (s *RawValueStyleEnumCaseContext) GetRuleContext() antlr.RuleContext

func (s *RawValueStyleEnumCaseContext) Identifier() IIdentifierContext

func (*RawValueStyleEnumCaseContext) IsRawValueStyleEnumCaseContext()

func (s *RawValueStyleEnumCaseContext) RawValueAssignment() IRawValueAssignmentContext

func (s *RawValueStyleEnumCaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type RawValueStyleEnumCaseListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyRawValueStyleEnumCaseListContext() *RawValueStyleEnumCaseListContext

func NewRawValueStyleEnumCaseListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RawValueStyleEnumCaseListContext

func (s *RawValueStyleEnumCaseListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RawValueStyleEnumCaseListContext) AllCOMMA() []antlr.TerminalNode

func (s *RawValueStyleEnumCaseListContext) AllRawValueStyleEnumCase() []IRawValueStyleEnumCaseContext

func (s *RawValueStyleEnumCaseListContext) COMMA(i int) antlr.TerminalNode

func (s *RawValueStyleEnumCaseListContext) GetParser() antlr.Parser

func (s *RawValueStyleEnumCaseListContext) GetRuleContext() antlr.RuleContext

func (*RawValueStyleEnumCaseListContext) IsRawValueStyleEnumCaseListContext()

func (s *RawValueStyleEnumCaseListContext) RawValueStyleEnumCase(i int) IRawValueStyleEnumCaseContext

func (s *RawValueStyleEnumCaseListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type RepeatWhileStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyRepeatWhileStatementContext() *RepeatWhileStatementContext

func NewRepeatWhileStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RepeatWhileStatementContext

func (s *RepeatWhileStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RepeatWhileStatementContext) CodeBlock() ICodeBlockContext

func (s *RepeatWhileStatementContext) Expression() IExpressionContext

func (s *RepeatWhileStatementContext) GetParser() antlr.Parser

func (s *RepeatWhileStatementContext) GetRuleContext() antlr.RuleContext

func (*RepeatWhileStatementContext) IsRepeatWhileStatementContext()

func (s *RepeatWhileStatementContext) REPEAT() antlr.TerminalNode

func (s *RepeatWhileStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *RepeatWhileStatementContext) WHILE() antlr.TerminalNode

type RequirementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyRequirementContext() *RequirementContext

func NewRequirementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RequirementContext

func (s *RequirementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RequirementContext) ConformanceRequirement() IConformanceRequirementContext

func (s *RequirementContext) GetParser() antlr.Parser

func (s *RequirementContext) GetRuleContext() antlr.RuleContext

func (*RequirementContext) IsRequirementContext()

func (s *RequirementContext) LayoutConstraintRequirement() ILayoutConstraintRequirementContext

func (s *RequirementContext) SameTypeRequirement() ISameTypeRequirementContext

func (s *RequirementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type RequirementListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyRequirementListContext() *RequirementListContext

func NewRequirementListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RequirementListContext

func (s *RequirementListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RequirementListContext) AllCOMMA() []antlr.TerminalNode

func (s *RequirementListContext) AllRequirement() []IRequirementContext

func (s *RequirementListContext) COMMA(i int) antlr.TerminalNode

func (s *RequirementListContext) GetParser() antlr.Parser

func (s *RequirementListContext) GetRuleContext() antlr.RuleContext

func (*RequirementListContext) IsRequirementListContext()

func (s *RequirementListContext) Requirement(i int) IRequirementContext

func (s *RequirementListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ReturnStatementContext struct {
	ControlTransferContext
}

func NewReturnStatementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ReturnStatementContext

func (s *ReturnStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ReturnStatementContext) Expression() IExpressionContext

func (s *ReturnStatementContext) GetRuleContext() antlr.RuleContext

func (s *ReturnStatementContext) RETURN() antlr.TerminalNode

type SameTypeRequirementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySameTypeRequirementContext() *SameTypeRequirementContext

func NewSameTypeRequirementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SameTypeRequirementContext

func (s *SameTypeRequirementContext) ASSIGN() antlr.TerminalNode

func (s *SameTypeRequirementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SameTypeRequirementContext) GetParser() antlr.Parser

func (s *SameTypeRequirementContext) GetRuleContext() antlr.RuleContext

func (*SameTypeRequirementContext) IsSameTypeRequirementContext()

func (s *SameTypeRequirementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *SameTypeRequirementContext) TypeIdentifier() ITypeIdentifierContext

func (s *SameTypeRequirementContext) Type_() ITypeContext

type SelfExprContext struct {
	PrimaryExpressionContext
}

func NewSelfExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelfExprContext

func (s *SelfExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SelfExprContext) GetRuleContext() antlr.RuleContext

func (s *SelfExprContext) SelfExpression() ISelfExpressionContext

type SelfExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySelfExpressionContext() *SelfExpressionContext

func NewSelfExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelfExpressionContext

func (s *SelfExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SelfExpressionContext) DOT() antlr.TerminalNode

func (s *SelfExpressionContext) FunctionCallArgumentList() IFunctionCallArgumentListContext

func (s *SelfExpressionContext) GetParser() antlr.Parser

func (s *SelfExpressionContext) GetRuleContext() antlr.RuleContext

func (s *SelfExpressionContext) INIT() antlr.TerminalNode

func (s *SelfExpressionContext) Identifier() IIdentifierContext

func (*SelfExpressionContext) IsSelfExpressionContext()

func (s *SelfExpressionContext) LBRACKET() antlr.TerminalNode

func (s *SelfExpressionContext) RBRACKET() antlr.TerminalNode

func (s *SelfExpressionContext) SELF() antlr.TerminalNode

func (s *SelfExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SelfTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySelfTypeContext() *SelfTypeContext

func NewSelfTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelfTypeContext

func (s *SelfTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SelfTypeContext) GetParser() antlr.Parser

func (s *SelfTypeContext) GetRuleContext() antlr.RuleContext

func (*SelfTypeContext) IsSelfTypeContext()

func (s *SelfTypeContext) SELF_UPPER() antlr.TerminalNode

func (s *SelfTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SelfType_Context struct {
	TypeContext
}

func NewSelfType_Context(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelfType_Context

func (s *SelfType_Context) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SelfType_Context) GetRuleContext() antlr.RuleContext

func (s *SelfType_Context) SelfType() ISelfTypeContext

type SetterClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySetterClauseContext() *SetterClauseContext

func NewSetterClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetterClauseContext

func (s *SetterClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SetterClauseContext) Attributes() IAttributesContext

func (s *SetterClauseContext) CodeBlock() ICodeBlockContext

func (s *SetterClauseContext) GetParser() antlr.Parser

func (s *SetterClauseContext) GetRuleContext() antlr.RuleContext

func (*SetterClauseContext) IsSetterClauseContext()

func (s *SetterClauseContext) MutationModifier() IMutationModifierContext

func (s *SetterClauseContext) SET_KW() antlr.TerminalNode

func (s *SetterClauseContext) SetterName() ISetterNameContext

func (s *SetterClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SetterKeywordClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySetterKeywordClauseContext() *SetterKeywordClauseContext

func NewSetterKeywordClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetterKeywordClauseContext

func (s *SetterKeywordClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SetterKeywordClauseContext) Attributes() IAttributesContext

func (s *SetterKeywordClauseContext) GetParser() antlr.Parser

func (s *SetterKeywordClauseContext) GetRuleContext() antlr.RuleContext

func (*SetterKeywordClauseContext) IsSetterKeywordClauseContext()

func (s *SetterKeywordClauseContext) MutationModifier() IMutationModifierContext

func (s *SetterKeywordClauseContext) SET_KW() antlr.TerminalNode

func (s *SetterKeywordClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SetterNameContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySetterNameContext() *SetterNameContext

func NewSetterNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetterNameContext

func (s *SetterNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SetterNameContext) GetParser() antlr.Parser

func (s *SetterNameContext) GetRuleContext() antlr.RuleContext

func (s *SetterNameContext) Identifier() IIdentifierContext

func (*SetterNameContext) IsSetterNameContext()

func (s *SetterNameContext) LPAREN() antlr.TerminalNode

func (s *SetterNameContext) RPAREN() antlr.TerminalNode

func (s *SetterNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type StatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStatementContext() *StatementContext

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext

func (s *StatementContext) CopyAll(ctx *StatementContext)

func (s *StatementContext) GetParser() antlr.Parser

func (s *StatementContext) GetRuleContext() antlr.RuleContext

func (*StatementContext) IsStatementContext()

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type StatementsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStatementsContext() *StatementsContext

func NewStatementsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementsContext

func (s *StatementsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StatementsContext) AllStatement() []IStatementContext

func (s *StatementsContext) GetParser() antlr.Parser

func (s *StatementsContext) GetRuleContext() antlr.RuleContext

func (*StatementsContext) IsStatementsContext()

func (s *StatementsContext) Statement(i int) IStatementContext

func (s *StatementsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type StructBodyContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructBodyContext() *StructBodyContext

func NewStructBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructBodyContext

func (s *StructBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructBodyContext) AllStructMember() []IStructMemberContext

func (s *StructBodyContext) GetParser() antlr.Parser

func (s *StructBodyContext) GetRuleContext() antlr.RuleContext

func (*StructBodyContext) IsStructBodyContext()

func (s *StructBodyContext) LBRACE() antlr.TerminalNode

func (s *StructBodyContext) RBRACE() antlr.TerminalNode

func (s *StructBodyContext) StructMember(i int) IStructMemberContext

func (s *StructBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type StructDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructDeclarationContext() *StructDeclarationContext

func NewStructDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructDeclarationContext

func (s *StructDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructDeclarationContext) AccessLevelModifier() IAccessLevelModifierContext

func (s *StructDeclarationContext) Attributes() IAttributesContext

func (s *StructDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *StructDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *StructDeclarationContext) GetParser() antlr.Parser

func (s *StructDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *StructDeclarationContext) Identifier() IIdentifierContext

func (*StructDeclarationContext) IsStructDeclarationContext()

func (s *StructDeclarationContext) STRUCT() antlr.TerminalNode

func (s *StructDeclarationContext) StructBody() IStructBodyContext

func (s *StructDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *StructDeclarationContext) TypeInheritanceClause() ITypeInheritanceClauseContext

type StructMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStructMemberContext() *StructMemberContext

func NewStructMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructMemberContext

func (s *StructMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StructMemberContext) CompilerControl() ICompilerControlContext

func (s *StructMemberContext) Declaration() IDeclarationContext

func (s *StructMemberContext) GetParser() antlr.Parser

func (s *StructMemberContext) GetRuleContext() antlr.RuleContext

func (*StructMemberContext) IsStructMemberContext()

func (s *StructMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SubscriptDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySubscriptDeclarationContext() *SubscriptDeclarationContext

func NewSubscriptDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubscriptDeclarationContext

func (s *SubscriptDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SubscriptDeclarationContext) GenericWhereClause() IGenericWhereClauseContext

func (s *SubscriptDeclarationContext) GetParser() antlr.Parser

func (s *SubscriptDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *SubscriptDeclarationContext) GetterSetterBlock() IGetterSetterBlockContext

func (s *SubscriptDeclarationContext) GetterSetterKeywordBlock() IGetterSetterKeywordBlockContext

func (*SubscriptDeclarationContext) IsSubscriptDeclarationContext()

func (s *SubscriptDeclarationContext) SubscriptHead() ISubscriptHeadContext

func (s *SubscriptDeclarationContext) SubscriptResult() ISubscriptResultContext

func (s *SubscriptDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SubscriptHeadContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySubscriptHeadContext() *SubscriptHeadContext

func NewSubscriptHeadContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubscriptHeadContext

func (s *SubscriptHeadContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SubscriptHeadContext) Attributes() IAttributesContext

func (s *SubscriptHeadContext) DeclarationModifiers() IDeclarationModifiersContext

func (s *SubscriptHeadContext) GenericParameterClause() IGenericParameterClauseContext

func (s *SubscriptHeadContext) GetParser() antlr.Parser

func (s *SubscriptHeadContext) GetRuleContext() antlr.RuleContext

func (*SubscriptHeadContext) IsSubscriptHeadContext()

func (s *SubscriptHeadContext) ParameterClause() IParameterClauseContext

func (s *SubscriptHeadContext) SUBSCRIPT() antlr.TerminalNode

func (s *SubscriptHeadContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SubscriptResultContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySubscriptResultContext() *SubscriptResultContext

func NewSubscriptResultContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubscriptResultContext

func (s *SubscriptResultContext) ARROW() antlr.TerminalNode

func (s *SubscriptResultContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SubscriptResultContext) Attributes() IAttributesContext

func (s *SubscriptResultContext) GetParser() antlr.Parser

func (s *SubscriptResultContext) GetRuleContext() antlr.RuleContext

func (*SubscriptResultContext) IsSubscriptResultContext()

func (s *SubscriptResultContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *SubscriptResultContext) Type_() ITypeContext

type SubscriptSuffixContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySubscriptSuffixContext() *SubscriptSuffixContext

func NewSubscriptSuffixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubscriptSuffixContext

func (s *SubscriptSuffixContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SubscriptSuffixContext) FunctionCallArgumentList() IFunctionCallArgumentListContext

func (s *SubscriptSuffixContext) GetParser() antlr.Parser

func (s *SubscriptSuffixContext) GetRuleContext() antlr.RuleContext

func (*SubscriptSuffixContext) IsSubscriptSuffixContext()

func (s *SubscriptSuffixContext) LBRACKET() antlr.TerminalNode

func (s *SubscriptSuffixContext) RBRACKET() antlr.TerminalNode

func (s *SubscriptSuffixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SuperExprContext struct {
	PrimaryExpressionContext
}

func NewSuperExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SuperExprContext

func (s *SuperExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SuperExprContext) GetRuleContext() antlr.RuleContext

func (s *SuperExprContext) SuperExpression() ISuperExpressionContext

type SuperExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySuperExpressionContext() *SuperExpressionContext

func NewSuperExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SuperExpressionContext

func (s *SuperExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SuperExpressionContext) DOT() antlr.TerminalNode

func (s *SuperExpressionContext) FunctionCallArgumentList() IFunctionCallArgumentListContext

func (s *SuperExpressionContext) GetParser() antlr.Parser

func (s *SuperExpressionContext) GetRuleContext() antlr.RuleContext

func (s *SuperExpressionContext) INIT() antlr.TerminalNode

func (s *SuperExpressionContext) Identifier() IIdentifierContext

func (*SuperExpressionContext) IsSuperExpressionContext()

func (s *SuperExpressionContext) LBRACKET() antlr.TerminalNode

func (s *SuperExpressionContext) RBRACKET() antlr.TerminalNode

func (s *SuperExpressionContext) SUPER() antlr.TerminalNode

func (s *SuperExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SuppressedTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySuppressedTypeContext() *SuppressedTypeContext

func NewSuppressedTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SuppressedTypeContext

func (s *SuppressedTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SuppressedTypeContext) GetParser() antlr.Parser

func (s *SuppressedTypeContext) GetRuleContext() antlr.RuleContext

func (*SuppressedTypeContext) IsSuppressedTypeContext()

func (s *SuppressedTypeContext) OPERATOR() antlr.TerminalNode

func (s *SuppressedTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *SuppressedTypeContext) TypeIdentifier() ITypeIdentifierContext

type SwitchCaseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchCaseContext() *SwitchCaseContext

func NewSwitchCaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchCaseContext

func (s *SwitchCaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchCaseContext) CaseLabel() ICaseLabelContext

func (s *SwitchCaseContext) DefaultLabel() IDefaultLabelContext

func (s *SwitchCaseContext) GetParser() antlr.Parser

func (s *SwitchCaseContext) GetRuleContext() antlr.RuleContext

func (*SwitchCaseContext) IsSwitchCaseContext()

func (s *SwitchCaseContext) Statements() IStatementsContext

func (s *SwitchCaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SwitchCasesContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchCasesContext() *SwitchCasesContext

func NewSwitchCasesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchCasesContext

func (s *SwitchCasesContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchCasesContext) AllSwitchCase() []ISwitchCaseContext

func (s *SwitchCasesContext) GetParser() antlr.Parser

func (s *SwitchCasesContext) GetRuleContext() antlr.RuleContext

func (*SwitchCasesContext) IsSwitchCasesContext()

func (s *SwitchCasesContext) SwitchCase(i int) ISwitchCaseContext

func (s *SwitchCasesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SwitchStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchStatementContext() *SwitchStatementContext

func NewSwitchStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchStatementContext

func (s *SwitchStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchStatementContext) Expression() IExpressionContext

func (s *SwitchStatementContext) GetParser() antlr.Parser

func (s *SwitchStatementContext) GetRuleContext() antlr.RuleContext

func (*SwitchStatementContext) IsSwitchStatementContext()

func (s *SwitchStatementContext) LBRACE() antlr.TerminalNode

func (s *SwitchStatementContext) RBRACE() antlr.TerminalNode

func (s *SwitchStatementContext) SWITCH() antlr.TerminalNode

func (s *SwitchStatementContext) SwitchCases() ISwitchCasesContext

func (s *SwitchStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TernaryExprContext struct {
	BinaryExpressionContext
}

func NewTernaryExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TernaryExprContext

func (s *TernaryExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TernaryExprContext) ConditionalOperator() IConditionalOperatorContext

func (s *TernaryExprContext) Expression() IExpressionContext

func (s *TernaryExprContext) GetRuleContext() antlr.RuleContext

func (s *TernaryExprContext) TryOperator() ITryOperatorContext

type ThrowStatementContext struct {
	ControlTransferContext
}

func NewThrowStatementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ThrowStatementContext

func (s *ThrowStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ThrowStatementContext) Expression() IExpressionContext

func (s *ThrowStatementContext) GetRuleContext() antlr.RuleContext

func (s *ThrowStatementContext) THROW() antlr.TerminalNode

type TopLevelContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTopLevelContext() *TopLevelContext

func NewTopLevelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopLevelContext

func (s *TopLevelContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TopLevelContext) EOF() antlr.TerminalNode

func (s *TopLevelContext) GetParser() antlr.Parser

func (s *TopLevelContext) GetRuleContext() antlr.RuleContext

func (*TopLevelContext) IsTopLevelContext()

func (s *TopLevelContext) Statements() IStatementsContext

func (s *TopLevelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TrailingClosuresContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTrailingClosuresContext() *TrailingClosuresContext

func NewTrailingClosuresContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TrailingClosuresContext

func (s *TrailingClosuresContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TrailingClosuresContext) AllLabeledTrailingClosure() []ILabeledTrailingClosureContext

func (s *TrailingClosuresContext) ClosureExpression() IClosureExpressionContext

func (s *TrailingClosuresContext) GetParser() antlr.Parser

func (s *TrailingClosuresContext) GetRuleContext() antlr.RuleContext

func (*TrailingClosuresContext) IsTrailingClosuresContext()

func (s *TrailingClosuresContext) LabeledTrailingClosure(i int) ILabeledTrailingClosureContext

func (s *TrailingClosuresContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TryOperatorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTryOperatorContext() *TryOperatorContext

func NewTryOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TryOperatorContext

func (s *TryOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TryOperatorContext) EXCLAIM_POSTFIX() antlr.TerminalNode

func (s *TryOperatorContext) GetParser() antlr.Parser

func (s *TryOperatorContext) GetRuleContext() antlr.RuleContext

func (*TryOperatorContext) IsTryOperatorContext()

func (s *TryOperatorContext) QUESTION_POSTFIX() antlr.TerminalNode

func (s *TryOperatorContext) TRY() antlr.TerminalNode

func (s *TryOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TryStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTryStatementContext() *TryStatementContext

func NewTryStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TryStatementContext

func (s *TryStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TryStatementContext) AllCatchClause() []ICatchClauseContext

func (s *TryStatementContext) CatchClause(i int) ICatchClauseContext

func (s *TryStatementContext) CodeBlock() ICodeBlockContext

func (s *TryStatementContext) GetParser() antlr.Parser

func (s *TryStatementContext) GetRuleContext() antlr.RuleContext

func (*TryStatementContext) IsTryStatementContext()

func (s *TryStatementContext) TRY() antlr.TerminalNode

func (s *TryStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TupTypeContext struct {
	TypeContext
}

func NewTupTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TupTypeContext

func (s *TupTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupTypeContext) GetRuleContext() antlr.RuleContext

func (s *TupTypeContext) TupleType() ITupleTypeContext

type TupleElementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleElementContext() *TupleElementContext

func NewTupleElementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleElementContext

func (s *TupleElementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleElementContext) COLON() antlr.TerminalNode

func (s *TupleElementContext) Expression() IExpressionContext

func (s *TupleElementContext) GetParser() antlr.Parser

func (s *TupleElementContext) GetRuleContext() antlr.RuleContext

func (s *TupleElementContext) Identifier() IIdentifierContext

func (*TupleElementContext) IsTupleElementContext()

func (s *TupleElementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TupleElementListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleElementListContext() *TupleElementListContext

func NewTupleElementListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleElementListContext

func (s *TupleElementListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleElementListContext) AllCOMMA() []antlr.TerminalNode

func (s *TupleElementListContext) AllTupleElement() []ITupleElementContext

func (s *TupleElementListContext) COMMA(i int) antlr.TerminalNode

func (s *TupleElementListContext) GetParser() antlr.Parser

func (s *TupleElementListContext) GetRuleContext() antlr.RuleContext

func (*TupleElementListContext) IsTupleElementListContext()

func (s *TupleElementListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleElementListContext) TupleElement(i int) ITupleElementContext

type TupleExprContext struct {
	PrimaryExpressionContext
}

func NewTupleExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TupleExprContext

func (s *TupleExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleExprContext) GetRuleContext() antlr.RuleContext

func (s *TupleExprContext) TupleExpression() ITupleExpressionContext

type TupleExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleExpressionContext() *TupleExpressionContext

func NewTupleExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleExpressionContext

func (s *TupleExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleExpressionContext) GetParser() antlr.Parser

func (s *TupleExpressionContext) GetRuleContext() antlr.RuleContext

func (*TupleExpressionContext) IsTupleExpressionContext()

func (s *TupleExpressionContext) LPAREN() antlr.TerminalNode

func (s *TupleExpressionContext) RPAREN() antlr.TerminalNode

func (s *TupleExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleExpressionContext) TupleElementList() ITupleElementListContext

type TuplePatContext struct {
	PatternContext
}

func NewTuplePatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TuplePatContext

func (s *TuplePatContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TuplePatContext) GetRuleContext() antlr.RuleContext

func (s *TuplePatContext) TuplePattern() ITuplePatternContext

func (s *TuplePatContext) TypeAnnotation() ITypeAnnotationContext

type TuplePatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTuplePatternContext() *TuplePatternContext

func NewTuplePatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TuplePatternContext

func (s *TuplePatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TuplePatternContext) GetParser() antlr.Parser

func (s *TuplePatternContext) GetRuleContext() antlr.RuleContext

func (*TuplePatternContext) IsTuplePatternContext()

func (s *TuplePatternContext) LPAREN() antlr.TerminalNode

func (s *TuplePatternContext) RPAREN() antlr.TerminalNode

func (s *TuplePatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TuplePatternContext) TuplePatternElementList() ITuplePatternElementListContext

type TuplePatternElementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTuplePatternElementContext() *TuplePatternElementContext

func NewTuplePatternElementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TuplePatternElementContext

func (s *TuplePatternElementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TuplePatternElementContext) COLON() antlr.TerminalNode

func (s *TuplePatternElementContext) GetParser() antlr.Parser

func (s *TuplePatternElementContext) GetRuleContext() antlr.RuleContext

func (s *TuplePatternElementContext) Identifier() IIdentifierContext

func (*TuplePatternElementContext) IsTuplePatternElementContext()

func (s *TuplePatternElementContext) Pattern() IPatternContext

func (s *TuplePatternElementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TuplePatternElementListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTuplePatternElementListContext() *TuplePatternElementListContext

func NewTuplePatternElementListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TuplePatternElementListContext

func (s *TuplePatternElementListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TuplePatternElementListContext) AllCOMMA() []antlr.TerminalNode

func (s *TuplePatternElementListContext) AllTuplePatternElement() []ITuplePatternElementContext

func (s *TuplePatternElementListContext) COMMA(i int) antlr.TerminalNode

func (s *TuplePatternElementListContext) GetParser() antlr.Parser

func (s *TuplePatternElementListContext) GetRuleContext() antlr.RuleContext

func (*TuplePatternElementListContext) IsTuplePatternElementListContext()

func (s *TuplePatternElementListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TuplePatternElementListContext) TuplePatternElement(i int) ITuplePatternElementContext

type TupleTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleTypeContext() *TupleTypeContext

func NewTupleTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleTypeContext

func (s *TupleTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleTypeContext) GetParser() antlr.Parser

func (s *TupleTypeContext) GetRuleContext() antlr.RuleContext

func (*TupleTypeContext) IsTupleTypeContext()

func (s *TupleTypeContext) LPAREN() antlr.TerminalNode

func (s *TupleTypeContext) RPAREN() antlr.TerminalNode

func (s *TupleTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleTypeContext) TupleTypeElementList() ITupleTypeElementListContext

type TupleTypeElementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleTypeElementContext() *TupleTypeElementContext

func NewTupleTypeElementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleTypeElementContext

func (s *TupleTypeElementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleTypeElementContext) GetParser() antlr.Parser

func (s *TupleTypeElementContext) GetRuleContext() antlr.RuleContext

func (s *TupleTypeElementContext) Identifier() IIdentifierContext

func (*TupleTypeElementContext) IsTupleTypeElementContext()

func (s *TupleTypeElementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleTypeElementContext) TypeAnnotation() ITypeAnnotationContext

func (s *TupleTypeElementContext) Type_() ITypeContext

type TupleTypeElementListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleTypeElementListContext() *TupleTypeElementListContext

func NewTupleTypeElementListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleTypeElementListContext

func (s *TupleTypeElementListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleTypeElementListContext) AllCOMMA() []antlr.TerminalNode

func (s *TupleTypeElementListContext) AllTupleTypeElement() []ITupleTypeElementContext

func (s *TupleTypeElementListContext) COMMA(i int) antlr.TerminalNode

func (s *TupleTypeElementListContext) GetParser() antlr.Parser

func (s *TupleTypeElementListContext) GetRuleContext() antlr.RuleContext

func (*TupleTypeElementListContext) IsTupleTypeElementListContext()

func (s *TupleTypeElementListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleTypeElementListContext) TupleTypeElement(i int) ITupleTypeElementContext

type TypeAnnotationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeAnnotationContext() *TypeAnnotationContext

func NewTypeAnnotationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeAnnotationContext

func (s *TypeAnnotationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeAnnotationContext) Attributes() IAttributesContext

func (s *TypeAnnotationContext) COLON() antlr.TerminalNode

func (s *TypeAnnotationContext) GetParser() antlr.Parser

func (s *TypeAnnotationContext) GetRuleContext() antlr.RuleContext

func (s *TypeAnnotationContext) INOUT() antlr.TerminalNode

func (*TypeAnnotationContext) IsTypeAnnotationContext()

func (s *TypeAnnotationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeAnnotationContext) Type_() ITypeContext

type TypeAnnotationHeadContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeAnnotationHeadContext() *TypeAnnotationHeadContext

func NewTypeAnnotationHeadContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeAnnotationHeadContext

func (s *TypeAnnotationHeadContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeAnnotationHeadContext) Attributes() IAttributesContext

func (s *TypeAnnotationHeadContext) GetParser() antlr.Parser

func (s *TypeAnnotationHeadContext) GetRuleContext() antlr.RuleContext

func (s *TypeAnnotationHeadContext) INOUT() antlr.TerminalNode

func (*TypeAnnotationHeadContext) IsTypeAnnotationHeadContext()

func (s *TypeAnnotationHeadContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TypeCastExprContext struct {
	BinaryExpressionContext
}

func NewTypeCastExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TypeCastExprContext

func (s *TypeCastExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeCastExprContext) GetRuleContext() antlr.RuleContext

func (s *TypeCastExprContext) TypeCastingOperator() ITypeCastingOperatorContext

type TypeCastingOperatorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeCastingOperatorContext() *TypeCastingOperatorContext

func NewTypeCastingOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeCastingOperatorContext

func (s *TypeCastingOperatorContext) AS() antlr.TerminalNode

func (s *TypeCastingOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeCastingOperatorContext) EXCLAIM_POSTFIX() antlr.TerminalNode

func (s *TypeCastingOperatorContext) GetParser() antlr.Parser

func (s *TypeCastingOperatorContext) GetRuleContext() antlr.RuleContext

func (s *TypeCastingOperatorContext) IS() antlr.TerminalNode

func (*TypeCastingOperatorContext) IsTypeCastingOperatorContext()

func (s *TypeCastingOperatorContext) QUESTION_POSTFIX() antlr.TerminalNode

func (s *TypeCastingOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeCastingOperatorContext) Type_() ITypeContext

type TypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeContext() *TypeContext

func NewTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeContext

func (s *TypeContext) CopyAll(ctx *TypeContext)

func (s *TypeContext) GetParser() antlr.Parser

func (s *TypeContext) GetRuleContext() antlr.RuleContext

func (*TypeContext) IsTypeContext()

func (s *TypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TypeIdentifierContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeIdentifierContext() *TypeIdentifierContext

func NewTypeIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeIdentifierContext

func (s *TypeIdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeIdentifierContext) DOT() antlr.TerminalNode

func (s *TypeIdentifierContext) GenericArgumentClause() IGenericArgumentClauseContext

func (s *TypeIdentifierContext) GetParser() antlr.Parser

func (s *TypeIdentifierContext) GetRuleContext() antlr.RuleContext

func (s *TypeIdentifierContext) Identifier() IIdentifierContext

func (*TypeIdentifierContext) IsTypeIdentifierContext()

func (s *TypeIdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeIdentifierContext) TypeIdentifier() ITypeIdentifierContext

type TypeInheritanceClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeInheritanceClauseContext() *TypeInheritanceClauseContext

func NewTypeInheritanceClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeInheritanceClauseContext

func (s *TypeInheritanceClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeInheritanceClauseContext) COLON() antlr.TerminalNode

func (s *TypeInheritanceClauseContext) GetParser() antlr.Parser

func (s *TypeInheritanceClauseContext) GetRuleContext() antlr.RuleContext

func (*TypeInheritanceClauseContext) IsTypeInheritanceClauseContext()

func (s *TypeInheritanceClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeInheritanceClauseContext) TypeInheritanceList() ITypeInheritanceListContext

type TypeInheritanceListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeInheritanceListContext() *TypeInheritanceListContext

func NewTypeInheritanceListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeInheritanceListContext

func (s *TypeInheritanceListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeInheritanceListContext) AllCOMMA() []antlr.TerminalNode

func (s *TypeInheritanceListContext) AllTypeIdentifier() []ITypeIdentifierContext

func (s *TypeInheritanceListContext) COMMA(i int) antlr.TerminalNode

func (s *TypeInheritanceListContext) GetParser() antlr.Parser

func (s *TypeInheritanceListContext) GetRuleContext() antlr.RuleContext

func (*TypeInheritanceListContext) IsTypeInheritanceListContext()

func (s *TypeInheritanceListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeInheritanceListContext) TypeIdentifier(i int) ITypeIdentifierContext

type TypealiasDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypealiasDeclarationContext() *TypealiasDeclarationContext

func NewTypealiasDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypealiasDeclarationContext

func (s *TypealiasDeclarationContext) ASSIGN() antlr.TerminalNode

func (s *TypealiasDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypealiasDeclarationContext) AccessLevelModifier() IAccessLevelModifierContext

func (s *TypealiasDeclarationContext) Attributes() IAttributesContext

func (s *TypealiasDeclarationContext) GenericParameterClause() IGenericParameterClauseContext

func (s *TypealiasDeclarationContext) GetParser() antlr.Parser

func (s *TypealiasDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *TypealiasDeclarationContext) Identifier() IIdentifierContext

func (*TypealiasDeclarationContext) IsTypealiasDeclarationContext()

func (s *TypealiasDeclarationContext) TYPEALIAS() antlr.TerminalNode

func (s *TypealiasDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypealiasDeclarationContext) Type_() ITypeContext

type UnionStyleEnumCaseClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyUnionStyleEnumCaseClauseContext() *UnionStyleEnumCaseClauseContext

func NewUnionStyleEnumCaseClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionStyleEnumCaseClauseContext

func (s *UnionStyleEnumCaseClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *UnionStyleEnumCaseClauseContext) Attributes() IAttributesContext

func (s *UnionStyleEnumCaseClauseContext) CASE() antlr.TerminalNode

func (s *UnionStyleEnumCaseClauseContext) GetParser() antlr.Parser

func (s *UnionStyleEnumCaseClauseContext) GetRuleContext() antlr.RuleContext

func (*UnionStyleEnumCaseClauseContext) IsUnionStyleEnumCaseClauseContext()

func (s *UnionStyleEnumCaseClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *UnionStyleEnumCaseClauseContext) UnionStyleEnumCaseList() IUnionStyleEnumCaseListContext

type UnionStyleEnumCaseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyUnionStyleEnumCaseContext() *UnionStyleEnumCaseContext

func NewUnionStyleEnumCaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionStyleEnumCaseContext

func (s *UnionStyleEnumCaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *UnionStyleEnumCaseContext) GetParser() antlr.Parser

func (s *UnionStyleEnumCaseContext) GetRuleContext() antlr.RuleContext

func (s *UnionStyleEnumCaseContext) Identifier() IIdentifierContext

func (*UnionStyleEnumCaseContext) IsUnionStyleEnumCaseContext()

func (s *UnionStyleEnumCaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *UnionStyleEnumCaseContext) TupleType() ITupleTypeContext

type UnionStyleEnumCaseListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyUnionStyleEnumCaseListContext() *UnionStyleEnumCaseListContext

func NewUnionStyleEnumCaseListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionStyleEnumCaseListContext

func (s *UnionStyleEnumCaseListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *UnionStyleEnumCaseListContext) AllCOMMA() []antlr.TerminalNode

func (s *UnionStyleEnumCaseListContext) AllUnionStyleEnumCase() []IUnionStyleEnumCaseContext

func (s *UnionStyleEnumCaseListContext) COMMA(i int) antlr.TerminalNode

func (s *UnionStyleEnumCaseListContext) GetParser() antlr.Parser

func (s *UnionStyleEnumCaseListContext) GetRuleContext() antlr.RuleContext

func (*UnionStyleEnumCaseListContext) IsUnionStyleEnumCaseListContext()

func (s *UnionStyleEnumCaseListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *UnionStyleEnumCaseListContext) UnionStyleEnumCase(i int) IUnionStyleEnumCaseContext

type ValueBindingPatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyValueBindingPatternContext() *ValueBindingPatternContext

func NewValueBindingPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueBindingPatternContext

func (s *ValueBindingPatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ValueBindingPatternContext) GetParser() antlr.Parser

func (s *ValueBindingPatternContext) GetRuleContext() antlr.RuleContext

func (s *ValueBindingPatternContext) INOUT() antlr.TerminalNode

func (*ValueBindingPatternContext) IsValueBindingPatternContext()

func (s *ValueBindingPatternContext) LET() antlr.TerminalNode

func (s *ValueBindingPatternContext) Pattern() IPatternContext

func (s *ValueBindingPatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ValueBindingPatternContext) VAR() antlr.TerminalNode

type VariableDeclarationContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyVariableDeclarationContext() *VariableDeclarationContext

func NewVariableDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableDeclarationContext

func (s *VariableDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *VariableDeclarationContext) CodeBlock() ICodeBlockContext

func (s *VariableDeclarationContext) GetParser() antlr.Parser

func (s *VariableDeclarationContext) GetRuleContext() antlr.RuleContext

func (s *VariableDeclarationContext) GetterSetterBlock() IGetterSetterBlockContext

func (s *VariableDeclarationContext) GetterSetterKeywordBlock() IGetterSetterKeywordBlockContext

func (s *VariableDeclarationContext) Initializer() IInitializerContext

func (*VariableDeclarationContext) IsVariableDeclarationContext()

func (s *VariableDeclarationContext) PatternInitializerList() IPatternInitializerListContext

func (s *VariableDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *VariableDeclarationContext) TypeAnnotation() ITypeAnnotationContext

func (s *VariableDeclarationContext) VariableDeclarationHead() IVariableDeclarationHeadContext

func (s *VariableDeclarationContext) VariableName() IVariableNameContext

func (s *VariableDeclarationContext) WillSetDidSetBlock() IWillSetDidSetBlockContext

type VariableDeclarationHeadContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyVariableDeclarationHeadContext() *VariableDeclarationHeadContext

func NewVariableDeclarationHeadContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableDeclarationHeadContext

func (s *VariableDeclarationHeadContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *VariableDeclarationHeadContext) Attributes() IAttributesContext

func (s *VariableDeclarationHeadContext) DeclarationModifiers() IDeclarationModifiersContext

func (s *VariableDeclarationHeadContext) GetParser() antlr.Parser

func (s *VariableDeclarationHeadContext) GetRuleContext() antlr.RuleContext

func (*VariableDeclarationHeadContext) IsVariableDeclarationHeadContext()

func (s *VariableDeclarationHeadContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *VariableDeclarationHeadContext) VAR() antlr.TerminalNode

type VariableNameContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyVariableNameContext() *VariableNameContext

func NewVariableNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableNameContext

func (s *VariableNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *VariableNameContext) GetParser() antlr.Parser

func (s *VariableNameContext) GetRuleContext() antlr.RuleContext

func (s *VariableNameContext) Identifier() IIdentifierContext

func (*VariableNameContext) IsVariableNameContext()

func (s *VariableNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type VariadicTypeContext struct {
	TypeContext
}

func NewVariadicTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *VariadicTypeContext

func (s *VariadicTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *VariadicTypeContext) ELLIPSIS() antlr.TerminalNode

func (s *VariadicTypeContext) GetRuleContext() antlr.RuleContext

func (s *VariadicTypeContext) Type_() ITypeContext

type VertexLexer struct {
	VertexLexerBase

	// Has unexported fields.
}

func NewVertexLexer(input antlr.CharStream) *VertexLexer
    NewVertexLexer produces a new lexer instance for the optional input
    antlr.CharStream.

type VertexLexerBase struct {
	*antlr.BaseLexer
}
    VertexLexerBase is the superclass for the ANTLR4-generated VertexLexer.
    The grammar declares `superClass = VertexLexerBase`. Add any custom lexer
    helper methods here as the language grows.

type VertexParser struct {
	VertexParserBase
}

func NewVertexParser(input antlr.TokenStream) *VertexParser
    NewVertexParser produces a new parser instance for the optional input
    antlr.TokenStream.

func (p *VertexParser) AccessLevelModifier() (localctx IAccessLevelModifierContext)

func (p *VertexParser) ActorBody() (localctx IActorBodyContext)

func (p *VertexParser) ActorDeclaration() (localctx IActorDeclarationContext)

func (p *VertexParser) ActorMember() (localctx IActorMemberContext)

func (p *VertexParser) ArgumentLabel() (localctx IArgumentLabelContext)

func (p *VertexParser) ArgumentNames() (localctx IArgumentNamesContext)

func (p *VertexParser) ArrayLiteral() (localctx IArrayLiteralContext)

func (p *VertexParser) ArrayLiteralItems() (localctx IArrayLiteralItemsContext)

func (p *VertexParser) ArrayType() (localctx IArrayTypeContext)

func (p *VertexParser) AssignmentOperator() (localctx IAssignmentOperatorContext)

func (p *VertexParser) AsyncModifier() (localctx IAsyncModifierContext)

func (p *VertexParser) Attribute() (localctx IAttributeContext)

func (p *VertexParser) AttributeArgument() (localctx IAttributeArgumentContext)

func (p *VertexParser) AttributeArgumentList() (localctx IAttributeArgumentListContext)

func (p *VertexParser) AttributeArguments() (localctx IAttributeArgumentsContext)

func (p *VertexParser) Attributes() (localctx IAttributesContext)

func (p *VertexParser) AvailabilityArgument() (localctx IAvailabilityArgumentContext)

func (p *VertexParser) AvailabilityArguments() (localctx IAvailabilityArgumentsContext)

func (p *VertexParser) AvailabilityCondition() (localctx IAvailabilityConditionContext)

func (p *VertexParser) AwaitOperator() (localctx IAwaitOperatorContext)

func (p *VertexParser) BinaryExpression() (localctx IBinaryExpressionContext)

func (p *VertexParser) BinaryExpression_Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) BinaryExpressions() (localctx IBinaryExpressionsContext)

func (p *VertexParser) BinaryOperator() (localctx IBinaryOperatorContext)

func (p *VertexParser) BooleanLiteral() (localctx IBooleanLiteralContext)

func (p *VertexParser) BranchStatement() (localctx IBranchStatementContext)

func (p *VertexParser) CaptureList() (localctx ICaptureListContext)

func (p *VertexParser) CaptureListItem() (localctx ICaptureListItemContext)

func (p *VertexParser) CaptureListItems() (localctx ICaptureListItemsContext)

func (p *VertexParser) CaptureSpecifier() (localctx ICaptureSpecifierContext)

func (p *VertexParser) CaseCondition() (localctx ICaseConditionContext)

func (p *VertexParser) CaseItem() (localctx ICaseItemContext)

func (p *VertexParser) CaseItemList() (localctx ICaseItemListContext)

func (p *VertexParser) CaseLabel() (localctx ICaseLabelContext)

func (p *VertexParser) CatchClause() (localctx ICatchClauseContext)

func (p *VertexParser) CatchPattern() (localctx ICatchPatternContext)

func (p *VertexParser) CatchPatternList() (localctx ICatchPatternListContext)

func (p *VertexParser) ClassBody() (localctx IClassBodyContext)

func (p *VertexParser) ClassDeclaration() (localctx IClassDeclarationContext)

func (p *VertexParser) ClassMember() (localctx IClassMemberContext)

func (p *VertexParser) ClosureExpression() (localctx IClosureExpressionContext)

func (p *VertexParser) ClosureParameter() (localctx IClosureParameterContext)

func (p *VertexParser) ClosureParameterClause() (localctx IClosureParameterClauseContext)

func (p *VertexParser) ClosureParameterList() (localctx IClosureParameterListContext)

func (p *VertexParser) ClosureSignature() (localctx IClosureSignatureContext)

func (p *VertexParser) CodeBlock() (localctx ICodeBlockContext)

func (p *VertexParser) CompilationCondition() (localctx ICompilationConditionContext)

func (p *VertexParser) CompilationCondition_Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) CompilerControl() (localctx ICompilerControlContext)

func (p *VertexParser) Condition() (localctx IConditionContext)

func (p *VertexParser) ConditionList() (localctx IConditionListContext)

func (p *VertexParser) ConditionalCompilationBlock() (localctx IConditionalCompilationBlockContext)

func (p *VertexParser) ConditionalOperator() (localctx IConditionalOperatorContext)

func (p *VertexParser) ConformanceRequirement() (localctx IConformanceRequirementContext)

func (p *VertexParser) ConstantDeclaration() (localctx IConstantDeclarationContext)

func (p *VertexParser) ControlTransfer() (localctx IControlTransferContext)

func (p *VertexParser) DecimalVersion() (localctx IDecimalVersionContext)

func (p *VertexParser) Declaration() (localctx IDeclarationContext)

func (p *VertexParser) DeclarationModifier() (localctx IDeclarationModifierContext)

func (p *VertexParser) DeclarationModifiers() (localctx IDeclarationModifiersContext)

func (p *VertexParser) DefaultArgumentClause() (localctx IDefaultArgumentClauseContext)

func (p *VertexParser) DefaultLabel() (localctx IDefaultLabelContext)

func (p *VertexParser) DeferStatement() (localctx IDeferStatementContext)

func (p *VertexParser) DeinitializerDeclaration() (localctx IDeinitializerDeclarationContext)

func (p *VertexParser) DiagnosticStatement() (localctx IDiagnosticStatementContext)

func (p *VertexParser) DictionaryLiteral() (localctx IDictionaryLiteralContext)

func (p *VertexParser) DictionaryLiteralItem() (localctx IDictionaryLiteralItemContext)

func (p *VertexParser) DictionaryLiteralItems() (localctx IDictionaryLiteralItemsContext)

func (p *VertexParser) DictionaryType() (localctx IDictionaryTypeContext)

func (p *VertexParser) DidSetClause() (localctx IDidSetClauseContext)

func (p *VertexParser) DoStatement() (localctx IDoStatementContext)

func (p *VertexParser) ElseClause() (localctx IElseClauseContext)

func (p *VertexParser) EnumCasePattern() (localctx IEnumCasePatternContext)

func (p *VertexParser) EnumDeclaration() (localctx IEnumDeclarationContext)

func (p *VertexParser) EnumMember() (localctx IEnumMemberContext)

func (p *VertexParser) EnumMembers() (localctx IEnumMembersContext)

func (p *VertexParser) ExistentialType() (localctx IExistentialTypeContext)

func (p *VertexParser) ExplicitMemberSuffix() (localctx IExplicitMemberSuffixContext)

func (p *VertexParser) Expression() (localctx IExpressionContext)

func (p *VertexParser) ExpressionPattern() (localctx IExpressionPatternContext)

func (p *VertexParser) ExtensionBody() (localctx IExtensionBodyContext)

func (p *VertexParser) ExtensionDeclaration() (localctx IExtensionDeclarationContext)

func (p *VertexParser) ExtensionMember() (localctx IExtensionMemberContext)

func (p *VertexParser) ExternalParameterName() (localctx IExternalParameterNameContext)

func (p *VertexParser) ForInStatement() (localctx IForInStatementContext)

func (p *VertexParser) ForStatement() (localctx IForStatementContext)

func (p *VertexParser) ForcedValueSuffix() (localctx IForcedValueSuffixContext)

func (p *VertexParser) FunctionBody() (localctx IFunctionBodyContext)

func (p *VertexParser) FunctionCallArgument() (localctx IFunctionCallArgumentContext)

func (p *VertexParser) FunctionCallArgumentClause() (localctx IFunctionCallArgumentClauseContext)

func (p *VertexParser) FunctionCallArgumentList() (localctx IFunctionCallArgumentListContext)

func (p *VertexParser) FunctionCallSuffix() (localctx IFunctionCallSuffixContext)

func (p *VertexParser) FunctionDeclaration() (localctx IFunctionDeclarationContext)

func (p *VertexParser) FunctionHead() (localctx IFunctionHeadContext)

func (p *VertexParser) FunctionResult() (localctx IFunctionResultContext)

func (p *VertexParser) FunctionSignature() (localctx IFunctionSignatureContext)

func (p *VertexParser) FunctionType() (localctx IFunctionTypeContext)

func (p *VertexParser) FunctionTypeArgument() (localctx IFunctionTypeArgumentContext)

func (p *VertexParser) FunctionTypeArgumentClause() (localctx IFunctionTypeArgumentClauseContext)

func (p *VertexParser) FunctionTypeArgumentList() (localctx IFunctionTypeArgumentListContext)

func (p *VertexParser) GenericArgument() (localctx IGenericArgumentContext)

func (p *VertexParser) GenericArgumentClause() (localctx IGenericArgumentClauseContext)

func (p *VertexParser) GenericArgumentList() (localctx IGenericArgumentListContext)

func (p *VertexParser) GenericParameter() (localctx IGenericParameterContext)

func (p *VertexParser) GenericParameterClause() (localctx IGenericParameterClauseContext)

func (p *VertexParser) GenericParameterList() (localctx IGenericParameterListContext)

func (p *VertexParser) GenericWhereClause() (localctx IGenericWhereClauseContext)

func (p *VertexParser) GetterClause() (localctx IGetterClauseContext)

func (p *VertexParser) GetterKeywordClause() (localctx IGetterKeywordClauseContext)

func (p *VertexParser) GetterSetterBlock() (localctx IGetterSetterBlockContext)

func (p *VertexParser) GetterSetterKeywordBlock() (localctx IGetterSetterKeywordBlockContext)

func (p *VertexParser) GuardStatement() (localctx IGuardStatementContext)

func (p *VertexParser) Identifier() (localctx IIdentifierContext)

func (p *VertexParser) IdentifierList() (localctx IIdentifierListContext)

func (p *VertexParser) IdentifierPattern() (localctx IIdentifierPatternContext)

func (p *VertexParser) IfStatement() (localctx IIfStatementContext)

func (p *VertexParser) ImplicitMemberExpression() (localctx IImplicitMemberExpressionContext)

func (p *VertexParser) ImportAlias() (localctx IImportAliasContext)

func (p *VertexParser) ImportDeclaration() (localctx IImportDeclarationContext)

func (p *VertexParser) ImportSpec() (localctx IImportSpecContext)

func (p *VertexParser) InOutExpression() (localctx IInOutExpressionContext)

func (p *VertexParser) Initializer() (localctx IInitializerContext)

func (p *VertexParser) InitializerBody() (localctx IInitializerBodyContext)

func (p *VertexParser) InitializerDeclaration() (localctx IInitializerDeclarationContext)

func (p *VertexParser) InitializerSuffix() (localctx IInitializerSuffixContext)

func (p *VertexParser) LabelName() (localctx ILabelNameContext)

func (p *VertexParser) LabeledStatement() (localctx ILabeledStatementContext)

func (p *VertexParser) LabeledTrailingClosure() (localctx ILabeledTrailingClosureContext)

func (p *VertexParser) LayoutConstraintRequirement() (localctx ILayoutConstraintRequirementContext)

func (p *VertexParser) LineControlStatement() (localctx ILineControlStatementContext)

func (p *VertexParser) Literal() (localctx ILiteralContext)

func (p *VertexParser) LiteralExpression() (localctx ILiteralExpressionContext)

func (p *VertexParser) LocalParameterName() (localctx ILocalParameterNameContext)

func (p *VertexParser) LoopStatement() (localctx ILoopStatementContext)

func (p *VertexParser) MacroDeclaration() (localctx IMacroDeclarationContext)

func (p *VertexParser) MacroExpansionExpression() (localctx IMacroExpansionExpressionContext)

func (p *VertexParser) MutationModifier() (localctx IMutationModifierContext)

func (p *VertexParser) NumericLiteral() (localctx INumericLiteralContext)

func (p *VertexParser) OpaqueType() (localctx IOpaqueTypeContext)

func (p *VertexParser) Operator_() (localctx IOperator_Context)

func (p *VertexParser) OptionalBindingCondition() (localctx IOptionalBindingConditionContext)

func (p *VertexParser) OptionalChainingLiteral() (localctx IOptionalChainingLiteralContext)

func (p *VertexParser) OptionalPattern() (localctx IOptionalPatternContext)

func (p *VertexParser) Parameter() (localctx IParameterContext)

func (p *VertexParser) ParameterClause() (localctx IParameterClauseContext)

func (p *VertexParser) ParameterList() (localctx IParameterListContext)

func (p *VertexParser) ParenthesizedExpression() (localctx IParenthesizedExpressionContext)

func (p *VertexParser) Pattern() (localctx IPatternContext)

func (p *VertexParser) PatternInitializer() (localctx IPatternInitializerContext)

func (p *VertexParser) PatternInitializerList() (localctx IPatternInitializerListContext)

func (p *VertexParser) Pattern_Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) PlatformCondition() (localctx IPlatformConditionContext)

func (p *VertexParser) PostfixExpression() (localctx IPostfixExpressionContext)

func (p *VertexParser) PostfixOperator() (localctx IPostfixOperatorContext)

func (p *VertexParser) PostfixOperator_Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) PostfixSelfSuffix() (localctx IPostfixSelfSuffixContext)

func (p *VertexParser) PostfixSuffix() (localctx IPostfixSuffixContext)

func (p *VertexParser) PoundFileExpression() (localctx IPoundFileExpressionContext)

func (p *VertexParser) PrefixExpression() (localctx IPrefixExpressionContext)

func (p *VertexParser) PrefixExpression_Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) PrefixOperator() (localctx IPrefixOperatorContext)

func (p *VertexParser) PrimaryAssociatedTypeClause() (localctx IPrimaryAssociatedTypeClauseContext)

func (p *VertexParser) PrimaryAssociatedTypeList() (localctx IPrimaryAssociatedTypeListContext)

func (p *VertexParser) PrimaryExpression() (localctx IPrimaryExpressionContext)

func (p *VertexParser) ProtocolAssociatedTypeDeclaration() (localctx IProtocolAssociatedTypeDeclarationContext)

func (p *VertexParser) ProtocolBody() (localctx IProtocolBodyContext)

func (p *VertexParser) ProtocolCompositionType() (localctx IProtocolCompositionTypeContext)

func (p *VertexParser) ProtocolCompositionTypeElement() (localctx IProtocolCompositionTypeElementContext)

func (p *VertexParser) ProtocolDeclaration() (localctx IProtocolDeclarationContext)

func (p *VertexParser) ProtocolInitializerDeclaration() (localctx IProtocolInitializerDeclarationContext)

func (p *VertexParser) ProtocolMember() (localctx IProtocolMemberContext)

func (p *VertexParser) ProtocolMemberDeclaration() (localctx IProtocolMemberDeclarationContext)

func (p *VertexParser) ProtocolMethodDeclaration() (localctx IProtocolMethodDeclarationContext)

func (p *VertexParser) ProtocolPropertyDeclaration() (localctx IProtocolPropertyDeclarationContext)

func (p *VertexParser) ProtocolSubscriptDeclaration() (localctx IProtocolSubscriptDeclarationContext)

func (p *VertexParser) RawValueAssignment() (localctx IRawValueAssignmentContext)

func (p *VertexParser) RawValueLiteral() (localctx IRawValueLiteralContext)

func (p *VertexParser) RawValueStyleEnumCase() (localctx IRawValueStyleEnumCaseContext)

func (p *VertexParser) RawValueStyleEnumCaseClause() (localctx IRawValueStyleEnumCaseClauseContext)

func (p *VertexParser) RawValueStyleEnumCaseList() (localctx IRawValueStyleEnumCaseListContext)

func (p *VertexParser) RepeatWhileStatement() (localctx IRepeatWhileStatementContext)

func (p *VertexParser) Requirement() (localctx IRequirementContext)

func (p *VertexParser) RequirementList() (localctx IRequirementListContext)

func (p *VertexParser) SameTypeRequirement() (localctx ISameTypeRequirementContext)

func (p *VertexParser) SelfExpression() (localctx ISelfExpressionContext)

func (p *VertexParser) SelfType() (localctx ISelfTypeContext)

func (p *VertexParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool

func (p *VertexParser) SetterClause() (localctx ISetterClauseContext)

func (p *VertexParser) SetterKeywordClause() (localctx ISetterKeywordClauseContext)

func (p *VertexParser) SetterName() (localctx ISetterNameContext)

func (p *VertexParser) Statement() (localctx IStatementContext)

func (p *VertexParser) Statements() (localctx IStatementsContext)

func (p *VertexParser) StructBody() (localctx IStructBodyContext)

func (p *VertexParser) StructDeclaration() (localctx IStructDeclarationContext)

func (p *VertexParser) StructMember() (localctx IStructMemberContext)

func (p *VertexParser) SubscriptDeclaration() (localctx ISubscriptDeclarationContext)

func (p *VertexParser) SubscriptHead() (localctx ISubscriptHeadContext)

func (p *VertexParser) SubscriptResult() (localctx ISubscriptResultContext)

func (p *VertexParser) SubscriptSuffix() (localctx ISubscriptSuffixContext)

func (p *VertexParser) SuperExpression() (localctx ISuperExpressionContext)

func (p *VertexParser) SuppressedType() (localctx ISuppressedTypeContext)

func (p *VertexParser) SwitchCase() (localctx ISwitchCaseContext)

func (p *VertexParser) SwitchCases() (localctx ISwitchCasesContext)

func (p *VertexParser) SwitchStatement() (localctx ISwitchStatementContext)

func (p *VertexParser) TopLevel() (localctx ITopLevelContext)

func (p *VertexParser) TrailingClosures() (localctx ITrailingClosuresContext)

func (p *VertexParser) TryOperator() (localctx ITryOperatorContext)

func (p *VertexParser) TryStatement() (localctx ITryStatementContext)

func (p *VertexParser) TupleElement() (localctx ITupleElementContext)

func (p *VertexParser) TupleElementList() (localctx ITupleElementListContext)

func (p *VertexParser) TupleExpression() (localctx ITupleExpressionContext)

func (p *VertexParser) TuplePattern() (localctx ITuplePatternContext)

func (p *VertexParser) TuplePatternElement() (localctx ITuplePatternElementContext)

func (p *VertexParser) TuplePatternElementList() (localctx ITuplePatternElementListContext)

func (p *VertexParser) TupleType() (localctx ITupleTypeContext)

func (p *VertexParser) TupleTypeElement() (localctx ITupleTypeElementContext)

func (p *VertexParser) TupleTypeElementList() (localctx ITupleTypeElementListContext)

func (p *VertexParser) TypeAnnotation() (localctx ITypeAnnotationContext)

func (p *VertexParser) TypeAnnotationHead() (localctx ITypeAnnotationHeadContext)

func (p *VertexParser) TypeCastingOperator() (localctx ITypeCastingOperatorContext)

func (p *VertexParser) TypeIdentifier() (localctx ITypeIdentifierContext)

func (p *VertexParser) TypeInheritanceClause() (localctx ITypeInheritanceClauseContext)

func (p *VertexParser) TypeInheritanceList() (localctx ITypeInheritanceListContext)

func (p *VertexParser) Type_() (localctx ITypeContext)

func (p *VertexParser) Type__Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *VertexParser) TypealiasDeclaration() (localctx ITypealiasDeclarationContext)

func (p *VertexParser) UnionStyleEnumCase() (localctx IUnionStyleEnumCaseContext)

func (p *VertexParser) UnionStyleEnumCaseClause() (localctx IUnionStyleEnumCaseClauseContext)

func (p *VertexParser) UnionStyleEnumCaseList() (localctx IUnionStyleEnumCaseListContext)

func (p *VertexParser) ValueBindingPattern() (localctx IValueBindingPatternContext)

func (p *VertexParser) VariableDeclaration() (localctx IVariableDeclarationContext)

func (p *VertexParser) VariableDeclarationHead() (localctx IVariableDeclarationHeadContext)

func (p *VertexParser) VariableName() (localctx IVariableNameContext)

func (p *VertexParser) WhereClause() (localctx IWhereClauseContext)

func (p *VertexParser) WhileStatement() (localctx IWhileStatementContext)

func (p *VertexParser) WildcardExpression() (localctx IWildcardExpressionContext)

func (p *VertexParser) WildcardPattern() (localctx IWildcardPatternContext)

func (p *VertexParser) WillSetClause() (localctx IWillSetClauseContext)

func (p *VertexParser) WillSetDidSetBlock() (localctx IWillSetDidSetBlockContext)

type VertexParserBase struct {
	*antlr.BaseParser
}
    VertexParserBase is the superclass for the ANTLR4-generated VertexParser.
    The grammar declares `superClass = VertexParserBase`; ANTLR embeds this
    struct into the generated parser so every predicate method below is
    accessible as p.isBinaryOp() etc. inside grammar actions.

type VertexParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by VertexParser#topLevel.
	VisitTopLevel(ctx *TopLevelContext) interface{}

	// Visit a parse tree produced by VertexParser#statements.
	VisitStatements(ctx *StatementsContext) interface{}

	// Visit a parse tree produced by VertexParser#expressionStatement.
	VisitExpressionStatement(ctx *ExpressionStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#declarationStatement.
	VisitDeclarationStatement(ctx *DeclarationStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#loopStmt.
	VisitLoopStmt(ctx *LoopStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#branchStmt.
	VisitBranchStmt(ctx *BranchStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#labeledStmt.
	VisitLabeledStmt(ctx *LabeledStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#controlTransferStmt.
	VisitControlTransferStmt(ctx *ControlTransferStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#deferStmt.
	VisitDeferStmt(ctx *DeferStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#doStmt.
	VisitDoStmt(ctx *DoStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#compilerControlStmt.
	VisitCompilerControlStmt(ctx *CompilerControlStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#macroDeclStmt.
	VisitMacroDeclStmt(ctx *MacroDeclStmtContext) interface{}

	// Visit a parse tree produced by VertexParser#loopStatement.
	VisitLoopStatement(ctx *LoopStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#forInStatement.
	VisitForInStatement(ctx *ForInStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#whileStatement.
	VisitWhileStatement(ctx *WhileStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#conditionList.
	VisitConditionList(ctx *ConditionListContext) interface{}

	// Visit a parse tree produced by VertexParser#condition.
	VisitCondition(ctx *ConditionContext) interface{}

	// Visit a parse tree produced by VertexParser#caseCondition.
	VisitCaseCondition(ctx *CaseConditionContext) interface{}

	// Visit a parse tree produced by VertexParser#optionalBindingCondition.
	VisitOptionalBindingCondition(ctx *OptionalBindingConditionContext) interface{}

	// Visit a parse tree produced by VertexParser#repeatWhileStatement.
	VisitRepeatWhileStatement(ctx *RepeatWhileStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#forStatement.
	VisitForStatement(ctx *ForStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#branchStatement.
	VisitBranchStatement(ctx *BranchStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#ifStatement.
	VisitIfStatement(ctx *IfStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#elseClause.
	VisitElseClause(ctx *ElseClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#guardStatement.
	VisitGuardStatement(ctx *GuardStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#switchStatement.
	VisitSwitchStatement(ctx *SwitchStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#switchCases.
	VisitSwitchCases(ctx *SwitchCasesContext) interface{}

	// Visit a parse tree produced by VertexParser#switchCase.
	VisitSwitchCase(ctx *SwitchCaseContext) interface{}

	// Visit a parse tree produced by VertexParser#caseLabel.
	VisitCaseLabel(ctx *CaseLabelContext) interface{}

	// Visit a parse tree produced by VertexParser#caseItemList.
	VisitCaseItemList(ctx *CaseItemListContext) interface{}

	// Visit a parse tree produced by VertexParser#caseItem.
	VisitCaseItem(ctx *CaseItemContext) interface{}

	// Visit a parse tree produced by VertexParser#defaultLabel.
	VisitDefaultLabel(ctx *DefaultLabelContext) interface{}

	// Visit a parse tree produced by VertexParser#whereClause.
	VisitWhereClause(ctx *WhereClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#labeledStatement.
	VisitLabeledStatement(ctx *LabeledStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#labelName.
	VisitLabelName(ctx *LabelNameContext) interface{}

	// Visit a parse tree produced by VertexParser#breakStatement.
	VisitBreakStatement(ctx *BreakStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#continueStatement.
	VisitContinueStatement(ctx *ContinueStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#fallthroughStatement.
	VisitFallthroughStatement(ctx *FallthroughStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#returnStatement.
	VisitReturnStatement(ctx *ReturnStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#throwStatement.
	VisitThrowStatement(ctx *ThrowStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#deferStatement.
	VisitDeferStatement(ctx *DeferStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#doStatement.
	VisitDoStatement(ctx *DoStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#catchClause.
	VisitCatchClause(ctx *CatchClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#catchPatternList.
	VisitCatchPatternList(ctx *CatchPatternListContext) interface{}

	// Visit a parse tree produced by VertexParser#catchPattern.
	VisitCatchPattern(ctx *CatchPatternContext) interface{}

	// Visit a parse tree produced by VertexParser#tryStatement.
	VisitTryStatement(ctx *TryStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#compilerControl.
	VisitCompilerControl(ctx *CompilerControlContext) interface{}

	// Visit a parse tree produced by VertexParser#conditionalCompilationBlock.
	VisitConditionalCompilationBlock(ctx *ConditionalCompilationBlockContext) interface{}

	// Visit a parse tree produced by VertexParser#compilationCondition.
	VisitCompilationCondition(ctx *CompilationConditionContext) interface{}

	// Visit a parse tree produced by VertexParser#platformCondition.
	VisitPlatformCondition(ctx *PlatformConditionContext) interface{}

	// Visit a parse tree produced by VertexParser#decimalVersion.
	VisitDecimalVersion(ctx *DecimalVersionContext) interface{}

	// Visit a parse tree produced by VertexParser#lineControlStatement.
	VisitLineControlStatement(ctx *LineControlStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#diagnosticStatement.
	VisitDiagnosticStatement(ctx *DiagnosticStatementContext) interface{}

	// Visit a parse tree produced by VertexParser#codeBlock.
	VisitCodeBlock(ctx *CodeBlockContext) interface{}

	// Visit a parse tree produced by VertexParser#declaration.
	VisitDeclaration(ctx *DeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#importDeclaration.
	VisitImportDeclaration(ctx *ImportDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#importSpec.
	VisitImportSpec(ctx *ImportSpecContext) interface{}

	// Visit a parse tree produced by VertexParser#importAlias.
	VisitImportAlias(ctx *ImportAliasContext) interface{}

	// Visit a parse tree produced by VertexParser#constantDeclaration.
	VisitConstantDeclaration(ctx *ConstantDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#patternInitializerList.
	VisitPatternInitializerList(ctx *PatternInitializerListContext) interface{}

	// Visit a parse tree produced by VertexParser#patternInitializer.
	VisitPatternInitializer(ctx *PatternInitializerContext) interface{}

	// Visit a parse tree produced by VertexParser#initializer.
	VisitInitializer(ctx *InitializerContext) interface{}

	// Visit a parse tree produced by VertexParser#variableDeclaration.
	VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#variableDeclarationHead.
	VisitVariableDeclarationHead(ctx *VariableDeclarationHeadContext) interface{}

	// Visit a parse tree produced by VertexParser#variableName.
	VisitVariableName(ctx *VariableNameContext) interface{}

	// Visit a parse tree produced by VertexParser#getterSetterBlock.
	VisitGetterSetterBlock(ctx *GetterSetterBlockContext) interface{}

	// Visit a parse tree produced by VertexParser#getterClause.
	VisitGetterClause(ctx *GetterClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#setterClause.
	VisitSetterClause(ctx *SetterClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#setterName.
	VisitSetterName(ctx *SetterNameContext) interface{}

	// Visit a parse tree produced by VertexParser#getterSetterKeywordBlock.
	VisitGetterSetterKeywordBlock(ctx *GetterSetterKeywordBlockContext) interface{}

	// Visit a parse tree produced by VertexParser#getterKeywordClause.
	VisitGetterKeywordClause(ctx *GetterKeywordClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#setterKeywordClause.
	VisitSetterKeywordClause(ctx *SetterKeywordClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#willSetDidSetBlock.
	VisitWillSetDidSetBlock(ctx *WillSetDidSetBlockContext) interface{}

	// Visit a parse tree produced by VertexParser#willSetClause.
	VisitWillSetClause(ctx *WillSetClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#didSetClause.
	VisitDidSetClause(ctx *DidSetClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#typealiasDeclaration.
	VisitTypealiasDeclaration(ctx *TypealiasDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#functionDeclaration.
	VisitFunctionDeclaration(ctx *FunctionDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#functionHead.
	VisitFunctionHead(ctx *FunctionHeadContext) interface{}

	// Visit a parse tree produced by VertexParser#functionSignature.
	VisitFunctionSignature(ctx *FunctionSignatureContext) interface{}

	// Visit a parse tree produced by VertexParser#asyncModifier.
	VisitAsyncModifier(ctx *AsyncModifierContext) interface{}

	// Visit a parse tree produced by VertexParser#functionResult.
	VisitFunctionResult(ctx *FunctionResultContext) interface{}

	// Visit a parse tree produced by VertexParser#parameterClause.
	VisitParameterClause(ctx *ParameterClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#parameterList.
	VisitParameterList(ctx *ParameterListContext) interface{}

	// Visit a parse tree produced by VertexParser#parameter.
	VisitParameter(ctx *ParameterContext) interface{}

	// Visit a parse tree produced by VertexParser#externalParameterName.
	VisitExternalParameterName(ctx *ExternalParameterNameContext) interface{}

	// Visit a parse tree produced by VertexParser#localParameterName.
	VisitLocalParameterName(ctx *LocalParameterNameContext) interface{}

	// Visit a parse tree produced by VertexParser#defaultArgumentClause.
	VisitDefaultArgumentClause(ctx *DefaultArgumentClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#functionBody.
	VisitFunctionBody(ctx *FunctionBodyContext) interface{}

	// Visit a parse tree produced by VertexParser#enumDeclaration.
	VisitEnumDeclaration(ctx *EnumDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#enumMembers.
	VisitEnumMembers(ctx *EnumMembersContext) interface{}

	// Visit a parse tree produced by VertexParser#enumMember.
	VisitEnumMember(ctx *EnumMemberContext) interface{}

	// Visit a parse tree produced by VertexParser#unionStyleEnumCaseClause.
	VisitUnionStyleEnumCaseClause(ctx *UnionStyleEnumCaseClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#unionStyleEnumCaseList.
	VisitUnionStyleEnumCaseList(ctx *UnionStyleEnumCaseListContext) interface{}

	// Visit a parse tree produced by VertexParser#unionStyleEnumCase.
	VisitUnionStyleEnumCase(ctx *UnionStyleEnumCaseContext) interface{}

	// Visit a parse tree produced by VertexParser#rawValueStyleEnumCaseClause.
	VisitRawValueStyleEnumCaseClause(ctx *RawValueStyleEnumCaseClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#rawValueStyleEnumCaseList.
	VisitRawValueStyleEnumCaseList(ctx *RawValueStyleEnumCaseListContext) interface{}

	// Visit a parse tree produced by VertexParser#rawValueStyleEnumCase.
	VisitRawValueStyleEnumCase(ctx *RawValueStyleEnumCaseContext) interface{}

	// Visit a parse tree produced by VertexParser#rawValueAssignment.
	VisitRawValueAssignment(ctx *RawValueAssignmentContext) interface{}

	// Visit a parse tree produced by VertexParser#rawValueLiteral.
	VisitRawValueLiteral(ctx *RawValueLiteralContext) interface{}

	// Visit a parse tree produced by VertexParser#structDeclaration.
	VisitStructDeclaration(ctx *StructDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#structBody.
	VisitStructBody(ctx *StructBodyContext) interface{}

	// Visit a parse tree produced by VertexParser#structMember.
	VisitStructMember(ctx *StructMemberContext) interface{}

	// Visit a parse tree produced by VertexParser#classDeclaration.
	VisitClassDeclaration(ctx *ClassDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#classBody.
	VisitClassBody(ctx *ClassBodyContext) interface{}

	// Visit a parse tree produced by VertexParser#classMember.
	VisitClassMember(ctx *ClassMemberContext) interface{}

	// Visit a parse tree produced by VertexParser#actorDeclaration.
	VisitActorDeclaration(ctx *ActorDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#actorBody.
	VisitActorBody(ctx *ActorBodyContext) interface{}

	// Visit a parse tree produced by VertexParser#actorMember.
	VisitActorMember(ctx *ActorMemberContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolDeclaration.
	VisitProtocolDeclaration(ctx *ProtocolDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#primaryAssociatedTypeClause.
	VisitPrimaryAssociatedTypeClause(ctx *PrimaryAssociatedTypeClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#primaryAssociatedTypeList.
	VisitPrimaryAssociatedTypeList(ctx *PrimaryAssociatedTypeListContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolBody.
	VisitProtocolBody(ctx *ProtocolBodyContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolMember.
	VisitProtocolMember(ctx *ProtocolMemberContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolMemberDeclaration.
	VisitProtocolMemberDeclaration(ctx *ProtocolMemberDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolPropertyDeclaration.
	VisitProtocolPropertyDeclaration(ctx *ProtocolPropertyDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolMethodDeclaration.
	VisitProtocolMethodDeclaration(ctx *ProtocolMethodDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolInitializerDeclaration.
	VisitProtocolInitializerDeclaration(ctx *ProtocolInitializerDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolSubscriptDeclaration.
	VisitProtocolSubscriptDeclaration(ctx *ProtocolSubscriptDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolAssociatedTypeDeclaration.
	VisitProtocolAssociatedTypeDeclaration(ctx *ProtocolAssociatedTypeDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#initializerDeclaration.
	VisitInitializerDeclaration(ctx *InitializerDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#initializerBody.
	VisitInitializerBody(ctx *InitializerBodyContext) interface{}

	// Visit a parse tree produced by VertexParser#deinitializerDeclaration.
	VisitDeinitializerDeclaration(ctx *DeinitializerDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#extensionDeclaration.
	VisitExtensionDeclaration(ctx *ExtensionDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#extensionBody.
	VisitExtensionBody(ctx *ExtensionBodyContext) interface{}

	// Visit a parse tree produced by VertexParser#extensionMember.
	VisitExtensionMember(ctx *ExtensionMemberContext) interface{}

	// Visit a parse tree produced by VertexParser#subscriptDeclaration.
	VisitSubscriptDeclaration(ctx *SubscriptDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#subscriptHead.
	VisitSubscriptHead(ctx *SubscriptHeadContext) interface{}

	// Visit a parse tree produced by VertexParser#subscriptResult.
	VisitSubscriptResult(ctx *SubscriptResultContext) interface{}

	// Visit a parse tree produced by VertexParser#macroDeclaration.
	VisitMacroDeclaration(ctx *MacroDeclarationContext) interface{}

	// Visit a parse tree produced by VertexParser#declarationModifiers.
	VisitDeclarationModifiers(ctx *DeclarationModifiersContext) interface{}

	// Visit a parse tree produced by VertexParser#declarationModifier.
	VisitDeclarationModifier(ctx *DeclarationModifierContext) interface{}

	// Visit a parse tree produced by VertexParser#accessLevelModifier.
	VisitAccessLevelModifier(ctx *AccessLevelModifierContext) interface{}

	// Visit a parse tree produced by VertexParser#mutationModifier.
	VisitMutationModifier(ctx *MutationModifierContext) interface{}

	// Visit a parse tree produced by VertexParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#tryOperator.
	VisitTryOperator(ctx *TryOperatorContext) interface{}

	// Visit a parse tree produced by VertexParser#awaitOperator.
	VisitAwaitOperator(ctx *AwaitOperatorContext) interface{}

	// Visit a parse tree produced by VertexParser#binaryExpressions.
	VisitBinaryExpressions(ctx *BinaryExpressionsContext) interface{}

	// Visit a parse tree produced by VertexParser#binaryOp.
	VisitBinaryOp(ctx *BinaryOpContext) interface{}

	// Visit a parse tree produced by VertexParser#assignmentExpr.
	VisitAssignmentExpr(ctx *AssignmentExprContext) interface{}

	// Visit a parse tree produced by VertexParser#ternaryExpr.
	VisitTernaryExpr(ctx *TernaryExprContext) interface{}

	// Visit a parse tree produced by VertexParser#typeCastExpr.
	VisitTypeCastExpr(ctx *TypeCastExprContext) interface{}

	// Visit a parse tree produced by VertexParser#binaryOperator.
	VisitBinaryOperator(ctx *BinaryOperatorContext) interface{}

	// Visit a parse tree produced by VertexParser#assignmentOperator.
	VisitAssignmentOperator(ctx *AssignmentOperatorContext) interface{}

	// Visit a parse tree produced by VertexParser#conditionalOperator.
	VisitConditionalOperator(ctx *ConditionalOperatorContext) interface{}

	// Visit a parse tree produced by VertexParser#typeCastingOperator.
	VisitTypeCastingOperator(ctx *TypeCastingOperatorContext) interface{}

	// Visit a parse tree produced by VertexParser#prefixOp.
	VisitPrefixOp(ctx *PrefixOpContext) interface{}

	// Visit a parse tree produced by VertexParser#inout.
	VisitInout(ctx *InoutContext) interface{}

	// Visit a parse tree produced by VertexParser#barePostfix.
	VisitBarePostfix(ctx *BarePostfixContext) interface{}

	// Visit a parse tree produced by VertexParser#inOutExpression.
	VisitInOutExpression(ctx *InOutExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#prefixOperator.
	VisitPrefixOperator(ctx *PrefixOperatorContext) interface{}

	// Visit a parse tree produced by VertexParser#postfixExpression.
	VisitPostfixExpression(ctx *PostfixExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#postfixSuffix.
	VisitPostfixSuffix(ctx *PostfixSuffixContext) interface{}

	// Visit a parse tree produced by VertexParser#postfixOperator.
	VisitPostfixOperator(ctx *PostfixOperatorContext) interface{}

	// Visit a parse tree produced by VertexParser#forcedValueSuffix.
	VisitForcedValueSuffix(ctx *ForcedValueSuffixContext) interface{}

	// Visit a parse tree produced by VertexParser#optionalChainingLiteral.
	VisitOptionalChainingLiteral(ctx *OptionalChainingLiteralContext) interface{}

	// Visit a parse tree produced by VertexParser#functionCallSuffix.
	VisitFunctionCallSuffix(ctx *FunctionCallSuffixContext) interface{}

	// Visit a parse tree produced by VertexParser#functionCallArgumentClause.
	VisitFunctionCallArgumentClause(ctx *FunctionCallArgumentClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#functionCallArgumentList.
	VisitFunctionCallArgumentList(ctx *FunctionCallArgumentListContext) interface{}

	// Visit a parse tree produced by VertexParser#functionCallArgument.
	VisitFunctionCallArgument(ctx *FunctionCallArgumentContext) interface{}

	// Visit a parse tree produced by VertexParser#operator_.
	VisitOperator_(ctx *Operator_Context) interface{}

	// Visit a parse tree produced by VertexParser#trailingClosures.
	VisitTrailingClosures(ctx *TrailingClosuresContext) interface{}

	// Visit a parse tree produced by VertexParser#labeledTrailingClosure.
	VisitLabeledTrailingClosure(ctx *LabeledTrailingClosureContext) interface{}

	// Visit a parse tree produced by VertexParser#initializerSuffix.
	VisitInitializerSuffix(ctx *InitializerSuffixContext) interface{}

	// Visit a parse tree produced by VertexParser#argumentNames.
	VisitArgumentNames(ctx *ArgumentNamesContext) interface{}

	// Visit a parse tree produced by VertexParser#explicitMemberSuffix.
	VisitExplicitMemberSuffix(ctx *ExplicitMemberSuffixContext) interface{}

	// Visit a parse tree produced by VertexParser#postfixSelfSuffix.
	VisitPostfixSelfSuffix(ctx *PostfixSelfSuffixContext) interface{}

	// Visit a parse tree produced by VertexParser#subscriptSuffix.
	VisitSubscriptSuffix(ctx *SubscriptSuffixContext) interface{}

	// Visit a parse tree produced by VertexParser#identifierExpr.
	VisitIdentifierExpr(ctx *IdentifierExprContext) interface{}

	// Visit a parse tree produced by VertexParser#litExpr.
	VisitLitExpr(ctx *LitExprContext) interface{}

	// Visit a parse tree produced by VertexParser#selfExpr.
	VisitSelfExpr(ctx *SelfExprContext) interface{}

	// Visit a parse tree produced by VertexParser#superExpr.
	VisitSuperExpr(ctx *SuperExprContext) interface{}

	// Visit a parse tree produced by VertexParser#closureExpr.
	VisitClosureExpr(ctx *ClosureExprContext) interface{}

	// Visit a parse tree produced by VertexParser#parenExpr.
	VisitParenExpr(ctx *ParenExprContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleExpr.
	VisitTupleExpr(ctx *TupleExprContext) interface{}

	// Visit a parse tree produced by VertexParser#implicitMemberExpr.
	VisitImplicitMemberExpr(ctx *ImplicitMemberExprContext) interface{}

	// Visit a parse tree produced by VertexParser#wildcardExpr.
	VisitWildcardExpr(ctx *WildcardExprContext) interface{}

	// Visit a parse tree produced by VertexParser#macroExpr.
	VisitMacroExpr(ctx *MacroExprContext) interface{}

	// Visit a parse tree produced by VertexParser#literalExpression.
	VisitLiteralExpression(ctx *LiteralExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#poundFileExpression.
	VisitPoundFileExpression(ctx *PoundFileExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}

	// Visit a parse tree produced by VertexParser#numericLiteral.
	VisitNumericLiteral(ctx *NumericLiteralContext) interface{}

	// Visit a parse tree produced by VertexParser#booleanLiteral.
	VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{}

	// Visit a parse tree produced by VertexParser#arrayLiteral.
	VisitArrayLiteral(ctx *ArrayLiteralContext) interface{}

	// Visit a parse tree produced by VertexParser#arrayLiteralItems.
	VisitArrayLiteralItems(ctx *ArrayLiteralItemsContext) interface{}

	// Visit a parse tree produced by VertexParser#dictionaryLiteral.
	VisitDictionaryLiteral(ctx *DictionaryLiteralContext) interface{}

	// Visit a parse tree produced by VertexParser#dictionaryLiteralItems.
	VisitDictionaryLiteralItems(ctx *DictionaryLiteralItemsContext) interface{}

	// Visit a parse tree produced by VertexParser#dictionaryLiteralItem.
	VisitDictionaryLiteralItem(ctx *DictionaryLiteralItemContext) interface{}

	// Visit a parse tree produced by VertexParser#selfExpression.
	VisitSelfExpression(ctx *SelfExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#superExpression.
	VisitSuperExpression(ctx *SuperExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#closureExpression.
	VisitClosureExpression(ctx *ClosureExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#closureSignature.
	VisitClosureSignature(ctx *ClosureSignatureContext) interface{}

	// Visit a parse tree produced by VertexParser#captureList.
	VisitCaptureList(ctx *CaptureListContext) interface{}

	// Visit a parse tree produced by VertexParser#captureListItems.
	VisitCaptureListItems(ctx *CaptureListItemsContext) interface{}

	// Visit a parse tree produced by VertexParser#captureListItem.
	VisitCaptureListItem(ctx *CaptureListItemContext) interface{}

	// Visit a parse tree produced by VertexParser#captureSpecifier.
	VisitCaptureSpecifier(ctx *CaptureSpecifierContext) interface{}

	// Visit a parse tree produced by VertexParser#closureParameterClause.
	VisitClosureParameterClause(ctx *ClosureParameterClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#closureParameterList.
	VisitClosureParameterList(ctx *ClosureParameterListContext) interface{}

	// Visit a parse tree produced by VertexParser#closureParameter.
	VisitClosureParameter(ctx *ClosureParameterContext) interface{}

	// Visit a parse tree produced by VertexParser#identifierList.
	VisitIdentifierList(ctx *IdentifierListContext) interface{}

	// Visit a parse tree produced by VertexParser#parenthesizedExpression.
	VisitParenthesizedExpression(ctx *ParenthesizedExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleExpression.
	VisitTupleExpression(ctx *TupleExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleElementList.
	VisitTupleElementList(ctx *TupleElementListContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleElement.
	VisitTupleElement(ctx *TupleElementContext) interface{}

	// Visit a parse tree produced by VertexParser#implicitMemberExpression.
	VisitImplicitMemberExpression(ctx *ImplicitMemberExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#wildcardExpression.
	VisitWildcardExpression(ctx *WildcardExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#macroExpansionExpression.
	VisitMacroExpansionExpression(ctx *MacroExpansionExpressionContext) interface{}

	// Visit a parse tree produced by VertexParser#variadicType.
	VisitVariadicType(ctx *VariadicTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#tupType.
	VisitTupType(ctx *TupTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#selfType_.
	VisitSelfType_(ctx *SelfType_Context) interface{}

	// Visit a parse tree produced by VertexParser#arrType.
	VisitArrType(ctx *ArrTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#opaqueType_.
	VisitOpaqueType_(ctx *OpaqueType_Context) interface{}

	// Visit a parse tree produced by VertexParser#dictType.
	VisitDictType(ctx *DictTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#namedType.
	VisitNamedType(ctx *NamedTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#optType.
	VisitOptType(ctx *OptTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#protoCompType.
	VisitProtoCompType(ctx *ProtoCompTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#existType.
	VisitExistType(ctx *ExistTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#annotatedType.
	VisitAnnotatedType(ctx *AnnotatedTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#metatypeType_.
	VisitMetatypeType_(ctx *MetatypeType_Context) interface{}

	// Visit a parse tree produced by VertexParser#funcType.
	VisitFuncType(ctx *FuncTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#typeAnnotationHead.
	VisitTypeAnnotationHead(ctx *TypeAnnotationHeadContext) interface{}

	// Visit a parse tree produced by VertexParser#functionType.
	VisitFunctionType(ctx *FunctionTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#functionTypeArgumentClause.
	VisitFunctionTypeArgumentClause(ctx *FunctionTypeArgumentClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#functionTypeArgumentList.
	VisitFunctionTypeArgumentList(ctx *FunctionTypeArgumentListContext) interface{}

	// Visit a parse tree produced by VertexParser#functionTypeArgument.
	VisitFunctionTypeArgument(ctx *FunctionTypeArgumentContext) interface{}

	// Visit a parse tree produced by VertexParser#argumentLabel.
	VisitArgumentLabel(ctx *ArgumentLabelContext) interface{}

	// Visit a parse tree produced by VertexParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#dictionaryType.
	VisitDictionaryType(ctx *DictionaryTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#typeIdentifier.
	VisitTypeIdentifier(ctx *TypeIdentifierContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleType.
	VisitTupleType(ctx *TupleTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleTypeElementList.
	VisitTupleTypeElementList(ctx *TupleTypeElementListContext) interface{}

	// Visit a parse tree produced by VertexParser#tupleTypeElement.
	VisitTupleTypeElement(ctx *TupleTypeElementContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolCompositionType.
	VisitProtocolCompositionType(ctx *ProtocolCompositionTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#protocolCompositionTypeElement.
	VisitProtocolCompositionTypeElement(ctx *ProtocolCompositionTypeElementContext) interface{}

	// Visit a parse tree produced by VertexParser#suppressedType.
	VisitSuppressedType(ctx *SuppressedTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#existentialType.
	VisitExistentialType(ctx *ExistentialTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#opaqueType.
	VisitOpaqueType(ctx *OpaqueTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#selfType.
	VisitSelfType(ctx *SelfTypeContext) interface{}

	// Visit a parse tree produced by VertexParser#typeAnnotation.
	VisitTypeAnnotation(ctx *TypeAnnotationContext) interface{}

	// Visit a parse tree produced by VertexParser#typeInheritanceClause.
	VisitTypeInheritanceClause(ctx *TypeInheritanceClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#typeInheritanceList.
	VisitTypeInheritanceList(ctx *TypeInheritanceListContext) interface{}

	// Visit a parse tree produced by VertexParser#genericParameterClause.
	VisitGenericParameterClause(ctx *GenericParameterClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#genericParameterList.
	VisitGenericParameterList(ctx *GenericParameterListContext) interface{}

	// Visit a parse tree produced by VertexParser#genericParameter.
	VisitGenericParameter(ctx *GenericParameterContext) interface{}

	// Visit a parse tree produced by VertexParser#genericWhereClause.
	VisitGenericWhereClause(ctx *GenericWhereClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#requirementList.
	VisitRequirementList(ctx *RequirementListContext) interface{}

	// Visit a parse tree produced by VertexParser#requirement.
	VisitRequirement(ctx *RequirementContext) interface{}

	// Visit a parse tree produced by VertexParser#conformanceRequirement.
	VisitConformanceRequirement(ctx *ConformanceRequirementContext) interface{}

	// Visit a parse tree produced by VertexParser#sameTypeRequirement.
	VisitSameTypeRequirement(ctx *SameTypeRequirementContext) interface{}

	// Visit a parse tree produced by VertexParser#layoutConstraintRequirement.
	VisitLayoutConstraintRequirement(ctx *LayoutConstraintRequirementContext) interface{}

	// Visit a parse tree produced by VertexParser#genericArgumentClause.
	VisitGenericArgumentClause(ctx *GenericArgumentClauseContext) interface{}

	// Visit a parse tree produced by VertexParser#genericArgumentList.
	VisitGenericArgumentList(ctx *GenericArgumentListContext) interface{}

	// Visit a parse tree produced by VertexParser#genericArgument.
	VisitGenericArgument(ctx *GenericArgumentContext) interface{}

	// Visit a parse tree produced by VertexParser#bindingPat.
	VisitBindingPat(ctx *BindingPatContext) interface{}

	// Visit a parse tree produced by VertexParser#isPat.
	VisitIsPat(ctx *IsPatContext) interface{}

	// Visit a parse tree produced by VertexParser#exprPat.
	VisitExprPat(ctx *ExprPatContext) interface{}

	// Visit a parse tree produced by VertexParser#enumCasePat.
	VisitEnumCasePat(ctx *EnumCasePatContext) interface{}

	// Visit a parse tree produced by VertexParser#tuplePat.
	VisitTuplePat(ctx *TuplePatContext) interface{}

	// Visit a parse tree produced by VertexParser#optPat.
	VisitOptPat(ctx *OptPatContext) interface{}

	// Visit a parse tree produced by VertexParser#identPat.
	VisitIdentPat(ctx *IdentPatContext) interface{}

	// Visit a parse tree produced by VertexParser#asPat.
	VisitAsPat(ctx *AsPatContext) interface{}

	// Visit a parse tree produced by VertexParser#wildcardPat.
	VisitWildcardPat(ctx *WildcardPatContext) interface{}

	// Visit a parse tree produced by VertexParser#wildcardPattern.
	VisitWildcardPattern(ctx *WildcardPatternContext) interface{}

	// Visit a parse tree produced by VertexParser#identifierPattern.
	VisitIdentifierPattern(ctx *IdentifierPatternContext) interface{}

	// Visit a parse tree produced by VertexParser#valueBindingPattern.
	VisitValueBindingPattern(ctx *ValueBindingPatternContext) interface{}

	// Visit a parse tree produced by VertexParser#tuplePattern.
	VisitTuplePattern(ctx *TuplePatternContext) interface{}

	// Visit a parse tree produced by VertexParser#tuplePatternElementList.
	VisitTuplePatternElementList(ctx *TuplePatternElementListContext) interface{}

	// Visit a parse tree produced by VertexParser#tuplePatternElement.
	VisitTuplePatternElement(ctx *TuplePatternElementContext) interface{}

	// Visit a parse tree produced by VertexParser#enumCasePattern.
	VisitEnumCasePattern(ctx *EnumCasePatternContext) interface{}

	// Visit a parse tree produced by VertexParser#optionalPattern.
	VisitOptionalPattern(ctx *OptionalPatternContext) interface{}

	// Visit a parse tree produced by VertexParser#expressionPattern.
	VisitExpressionPattern(ctx *ExpressionPatternContext) interface{}

	// Visit a parse tree produced by VertexParser#attributes.
	VisitAttributes(ctx *AttributesContext) interface{}

	// Visit a parse tree produced by VertexParser#attribute.
	VisitAttribute(ctx *AttributeContext) interface{}

	// Visit a parse tree produced by VertexParser#attributeArguments.
	VisitAttributeArguments(ctx *AttributeArgumentsContext) interface{}

	// Visit a parse tree produced by VertexParser#attributeArgumentList.
	VisitAttributeArgumentList(ctx *AttributeArgumentListContext) interface{}

	// Visit a parse tree produced by VertexParser#attributeArgument.
	VisitAttributeArgument(ctx *AttributeArgumentContext) interface{}

	// Visit a parse tree produced by VertexParser#availabilityCondition.
	VisitAvailabilityCondition(ctx *AvailabilityConditionContext) interface{}

	// Visit a parse tree produced by VertexParser#availabilityArguments.
	VisitAvailabilityArguments(ctx *AvailabilityArgumentsContext) interface{}

	// Visit a parse tree produced by VertexParser#availabilityArgument.
	VisitAvailabilityArgument(ctx *AvailabilityArgumentContext) interface{}

	// Visit a parse tree produced by VertexParser#identifier.
	VisitIdentifier(ctx *IdentifierContext) interface{}
}
    A complete Visitor for a parse tree produced by VertexParser.

type WhereClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyWhereClauseContext() *WhereClauseContext

func NewWhereClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereClauseContext

func (s *WhereClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *WhereClauseContext) Expression() IExpressionContext

func (s *WhereClauseContext) GetParser() antlr.Parser

func (s *WhereClauseContext) GetRuleContext() antlr.RuleContext

func (*WhereClauseContext) IsWhereClauseContext()

func (s *WhereClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *WhereClauseContext) WHERE() antlr.TerminalNode

type WhileStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyWhileStatementContext() *WhileStatementContext

func NewWhileStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhileStatementContext

func (s *WhileStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *WhileStatementContext) CodeBlock() ICodeBlockContext

func (s *WhileStatementContext) ConditionList() IConditionListContext

func (s *WhileStatementContext) GetParser() antlr.Parser

func (s *WhileStatementContext) GetRuleContext() antlr.RuleContext

func (*WhileStatementContext) IsWhileStatementContext()

func (s *WhileStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *WhileStatementContext) WHILE() antlr.TerminalNode

type WildcardExprContext struct {
	PrimaryExpressionContext
}

func NewWildcardExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *WildcardExprContext

func (s *WildcardExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *WildcardExprContext) GetRuleContext() antlr.RuleContext

func (s *WildcardExprContext) WildcardExpression() IWildcardExpressionContext

type WildcardExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyWildcardExpressionContext() *WildcardExpressionContext

func NewWildcardExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WildcardExpressionContext

func (s *WildcardExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *WildcardExpressionContext) GetParser() antlr.Parser

func (s *WildcardExpressionContext) GetRuleContext() antlr.RuleContext

func (*WildcardExpressionContext) IsWildcardExpressionContext()

func (s *WildcardExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *WildcardExpressionContext) UNDERSCORE() antlr.TerminalNode

type WildcardPatContext struct {
	PatternContext
}

func NewWildcardPatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *WildcardPatContext

func (s *WildcardPatContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *WildcardPatContext) GetRuleContext() antlr.RuleContext

func (s *WildcardPatContext) TypeAnnotation() ITypeAnnotationContext

func (s *WildcardPatContext) WildcardPattern() IWildcardPatternContext

type WildcardPatternContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyWildcardPatternContext() *WildcardPatternContext

func NewWildcardPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WildcardPatternContext

func (s *WildcardPatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *WildcardPatternContext) GetParser() antlr.Parser

func (s *WildcardPatternContext) GetRuleContext() antlr.RuleContext

func (*WildcardPatternContext) IsWildcardPatternContext()

func (s *WildcardPatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *WildcardPatternContext) UNDERSCORE() antlr.TerminalNode

type WillSetClauseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyWillSetClauseContext() *WillSetClauseContext

func NewWillSetClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WillSetClauseContext

func (s *WillSetClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *WillSetClauseContext) Attributes() IAttributesContext

func (s *WillSetClauseContext) CodeBlock() ICodeBlockContext

func (s *WillSetClauseContext) GetParser() antlr.Parser

func (s *WillSetClauseContext) GetRuleContext() antlr.RuleContext

func (*WillSetClauseContext) IsWillSetClauseContext()

func (s *WillSetClauseContext) SetterName() ISetterNameContext

func (s *WillSetClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *WillSetClauseContext) WILLSET_KW() antlr.TerminalNode

type WillSetDidSetBlockContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyWillSetDidSetBlockContext() *WillSetDidSetBlockContext

func NewWillSetDidSetBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WillSetDidSetBlockContext

func (s *WillSetDidSetBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *WillSetDidSetBlockContext) DidSetClause() IDidSetClauseContext

func (s *WillSetDidSetBlockContext) GetParser() antlr.Parser

func (s *WillSetDidSetBlockContext) GetRuleContext() antlr.RuleContext

func (*WillSetDidSetBlockContext) IsWillSetDidSetBlockContext()

func (s *WillSetDidSetBlockContext) LBRACE() antlr.TerminalNode

func (s *WillSetDidSetBlockContext) RBRACE() antlr.TerminalNode

func (s *WillSetDidSetBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *WillSetDidSetBlockContext) WillSetClause() IWillSetClauseContext

