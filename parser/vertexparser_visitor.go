// Code generated from VertexParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // VertexParser

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by VertexParser.
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
