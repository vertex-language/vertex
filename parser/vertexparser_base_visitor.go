// Code generated from VertexParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // VertexParser

import "github.com/antlr4-go/antlr/v4"

type BaseVertexParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseVertexParserVisitor) VisitTopLevel(ctx *TopLevelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStatements(ctx *StatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDeclarationStatement(ctx *DeclarationStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLoopStmt(ctx *LoopStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBranchStmt(ctx *BranchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLabeledStmt(ctx *LabeledStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitControlTransferStmt(ctx *ControlTransferStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDeferStmt(ctx *DeferStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDoStmt(ctx *DoStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCompilerControlStmt(ctx *CompilerControlStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitMacroDeclStmt(ctx *MacroDeclStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLoopStatement(ctx *LoopStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitForInStatement(ctx *ForInStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWhileStatement(ctx *WhileStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitConditionList(ctx *ConditionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCondition(ctx *ConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCaseCondition(ctx *CaseConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitOptionalBindingCondition(ctx *OptionalBindingConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitRepeatWhileStatement(ctx *RepeatWhileStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitForStatement(ctx *ForStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBranchStatement(ctx *BranchStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIfStatement(ctx *IfStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitElseClause(ctx *ElseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGuardStatement(ctx *GuardStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSwitchStatement(ctx *SwitchStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSwitchCases(ctx *SwitchCasesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCaseLabel(ctx *CaseLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCaseItemList(ctx *CaseItemListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCaseItem(ctx *CaseItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDefaultLabel(ctx *DefaultLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWhereClause(ctx *WhereClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLabeledStatement(ctx *LabeledStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLabelName(ctx *LabelNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBreakStatement(ctx *BreakStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitContinueStatement(ctx *ContinueStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFallthroughStatement(ctx *FallthroughStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitReturnStatement(ctx *ReturnStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitThrowStatement(ctx *ThrowStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDeferStatement(ctx *DeferStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDoStatement(ctx *DoStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCatchClause(ctx *CatchClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCatchPatternList(ctx *CatchPatternListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCatchPattern(ctx *CatchPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTryStatement(ctx *TryStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCompilerControl(ctx *CompilerControlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitConditionalCompilationBlock(ctx *ConditionalCompilationBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCompilationCondition(ctx *CompilationConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPlatformCondition(ctx *PlatformConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDecimalVersion(ctx *DecimalVersionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLineControlStatement(ctx *LineControlStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDiagnosticStatement(ctx *DiagnosticStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCodeBlock(ctx *CodeBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDeclaration(ctx *DeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitImportDeclaration(ctx *ImportDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitImportSpec(ctx *ImportSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitImportAlias(ctx *ImportAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitConstantDeclaration(ctx *ConstantDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPatternInitializerList(ctx *PatternInitializerListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPatternInitializer(ctx *PatternInitializerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitInitializer(ctx *InitializerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitVariableDeclarationHead(ctx *VariableDeclarationHeadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitVariableName(ctx *VariableNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGetterSetterBlock(ctx *GetterSetterBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGetterClause(ctx *GetterClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSetterClause(ctx *SetterClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSetterName(ctx *SetterNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGetterSetterKeywordBlock(ctx *GetterSetterKeywordBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGetterKeywordClause(ctx *GetterKeywordClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSetterKeywordClause(ctx *SetterKeywordClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWillSetDidSetBlock(ctx *WillSetDidSetBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWillSetClause(ctx *WillSetClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDidSetClause(ctx *DidSetClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypealiasDeclaration(ctx *TypealiasDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionDeclaration(ctx *FunctionDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionHead(ctx *FunctionHeadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionSignature(ctx *FunctionSignatureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAsyncModifier(ctx *AsyncModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionResult(ctx *FunctionResultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitParameterClause(ctx *ParameterClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitParameterList(ctx *ParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitParameter(ctx *ParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExternalParameterName(ctx *ExternalParameterNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLocalParameterName(ctx *LocalParameterNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDefaultArgumentClause(ctx *DefaultArgumentClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionBody(ctx *FunctionBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumDeclaration(ctx *EnumDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumMembers(ctx *EnumMembersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumMember(ctx *EnumMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitUnionStyleEnumCaseClause(ctx *UnionStyleEnumCaseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitUnionStyleEnumCaseList(ctx *UnionStyleEnumCaseListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitUnionStyleEnumCase(ctx *UnionStyleEnumCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitRawValueStyleEnumCaseClause(ctx *RawValueStyleEnumCaseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitRawValueStyleEnumCaseList(ctx *RawValueStyleEnumCaseListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitRawValueStyleEnumCase(ctx *RawValueStyleEnumCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitRawValueAssignment(ctx *RawValueAssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitRawValueLiteral(ctx *RawValueLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructDeclaration(ctx *StructDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructBody(ctx *StructBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructMember(ctx *StructMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClassDeclaration(ctx *ClassDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClassBody(ctx *ClassBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClassMember(ctx *ClassMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitActorDeclaration(ctx *ActorDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitActorBody(ctx *ActorBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitActorMember(ctx *ActorMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolDeclaration(ctx *ProtocolDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPrimaryAssociatedTypeClause(ctx *PrimaryAssociatedTypeClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPrimaryAssociatedTypeList(ctx *PrimaryAssociatedTypeListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolBody(ctx *ProtocolBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolMember(ctx *ProtocolMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolMemberDeclaration(ctx *ProtocolMemberDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolPropertyDeclaration(ctx *ProtocolPropertyDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolMethodDeclaration(ctx *ProtocolMethodDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolInitializerDeclaration(ctx *ProtocolInitializerDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolSubscriptDeclaration(ctx *ProtocolSubscriptDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolAssociatedTypeDeclaration(ctx *ProtocolAssociatedTypeDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitInitializerDeclaration(ctx *InitializerDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitInitializerBody(ctx *InitializerBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDeinitializerDeclaration(ctx *DeinitializerDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExtensionDeclaration(ctx *ExtensionDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExtensionBody(ctx *ExtensionBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExtensionMember(ctx *ExtensionMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSubscriptDeclaration(ctx *SubscriptDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSubscriptHead(ctx *SubscriptHeadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSubscriptResult(ctx *SubscriptResultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitMacroDeclaration(ctx *MacroDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDeclarationModifiers(ctx *DeclarationModifiersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDeclarationModifier(ctx *DeclarationModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAccessLevelModifier(ctx *AccessLevelModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitMutationModifier(ctx *MutationModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTryOperator(ctx *TryOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAwaitOperator(ctx *AwaitOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBinaryExpressions(ctx *BinaryExpressionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBinaryOp(ctx *BinaryOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAssignmentExpr(ctx *AssignmentExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTernaryExpr(ctx *TernaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeCastExpr(ctx *TypeCastExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBinaryOperator(ctx *BinaryOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAssignmentOperator(ctx *AssignmentOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitConditionalOperator(ctx *ConditionalOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeCastingOperator(ctx *TypeCastingOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPrefixOp(ctx *PrefixOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitInout(ctx *InoutContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBarePostfix(ctx *BarePostfixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitInOutExpression(ctx *InOutExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPrefixOperator(ctx *PrefixOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPostfixExpression(ctx *PostfixExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPostfixSuffix(ctx *PostfixSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPostfixOperator(ctx *PostfixOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitForcedValueSuffix(ctx *ForcedValueSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitOptionalChainingLiteral(ctx *OptionalChainingLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionCallSuffix(ctx *FunctionCallSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionCallArgumentClause(ctx *FunctionCallArgumentClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionCallArgumentList(ctx *FunctionCallArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionCallArgument(ctx *FunctionCallArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitOperator_(ctx *Operator_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTrailingClosures(ctx *TrailingClosuresContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLabeledTrailingClosure(ctx *LabeledTrailingClosureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitInitializerSuffix(ctx *InitializerSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArgumentNames(ctx *ArgumentNamesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExplicitMemberSuffix(ctx *ExplicitMemberSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPostfixSelfSuffix(ctx *PostfixSelfSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSubscriptSuffix(ctx *SubscriptSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIdentifierExpr(ctx *IdentifierExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLitExpr(ctx *LitExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSelfExpr(ctx *SelfExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSuperExpr(ctx *SuperExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClosureExpr(ctx *ClosureExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitParenExpr(ctx *ParenExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleExpr(ctx *TupleExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitImplicitMemberExpr(ctx *ImplicitMemberExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWildcardExpr(ctx *WildcardExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitMacroExpr(ctx *MacroExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLiteralExpression(ctx *LiteralExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPoundFileExpression(ctx *PoundFileExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitNumericLiteral(ctx *NumericLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArrayLiteral(ctx *ArrayLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArrayLiteralItems(ctx *ArrayLiteralItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDictionaryLiteral(ctx *DictionaryLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDictionaryLiteralItems(ctx *DictionaryLiteralItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDictionaryLiteralItem(ctx *DictionaryLiteralItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSelfExpression(ctx *SelfExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSuperExpression(ctx *SuperExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClosureExpression(ctx *ClosureExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClosureSignature(ctx *ClosureSignatureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCaptureList(ctx *CaptureListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCaptureListItems(ctx *CaptureListItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCaptureListItem(ctx *CaptureListItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCaptureSpecifier(ctx *CaptureSpecifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClosureParameterClause(ctx *ClosureParameterClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClosureParameterList(ctx *ClosureParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClosureParameter(ctx *ClosureParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIdentifierList(ctx *IdentifierListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitParenthesizedExpression(ctx *ParenthesizedExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleExpression(ctx *TupleExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleElementList(ctx *TupleElementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleElement(ctx *TupleElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitImplicitMemberExpression(ctx *ImplicitMemberExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWildcardExpression(ctx *WildcardExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitMacroExpansionExpression(ctx *MacroExpansionExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitVariadicType(ctx *VariadicTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupType(ctx *TupTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSelfType_(ctx *SelfType_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArrType(ctx *ArrTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitOpaqueType_(ctx *OpaqueType_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDictType(ctx *DictTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitNamedType(ctx *NamedTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitOptType(ctx *OptTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtoCompType(ctx *ProtoCompTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExistType(ctx *ExistTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAnnotatedType(ctx *AnnotatedTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitMetatypeType_(ctx *MetatypeType_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFuncType(ctx *FuncTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeAnnotationHead(ctx *TypeAnnotationHeadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionType(ctx *FunctionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionTypeArgumentClause(ctx *FunctionTypeArgumentClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionTypeArgumentList(ctx *FunctionTypeArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFunctionTypeArgument(ctx *FunctionTypeArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArgumentLabel(ctx *ArgumentLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDictionaryType(ctx *DictionaryTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeIdentifier(ctx *TypeIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleType(ctx *TupleTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleTypeElementList(ctx *TupleTypeElementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleTypeElement(ctx *TupleTypeElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolCompositionType(ctx *ProtocolCompositionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitProtocolCompositionTypeElement(ctx *ProtocolCompositionTypeElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSuppressedType(ctx *SuppressedTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExistentialType(ctx *ExistentialTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitOpaqueType(ctx *OpaqueTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSelfType(ctx *SelfTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeAnnotation(ctx *TypeAnnotationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeInheritanceClause(ctx *TypeInheritanceClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeInheritanceList(ctx *TypeInheritanceListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGenericParameterClause(ctx *GenericParameterClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGenericParameterList(ctx *GenericParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGenericParameter(ctx *GenericParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGenericWhereClause(ctx *GenericWhereClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitRequirementList(ctx *RequirementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitRequirement(ctx *RequirementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitConformanceRequirement(ctx *ConformanceRequirementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSameTypeRequirement(ctx *SameTypeRequirementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLayoutConstraintRequirement(ctx *LayoutConstraintRequirementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGenericArgumentClause(ctx *GenericArgumentClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGenericArgumentList(ctx *GenericArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGenericArgument(ctx *GenericArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBindingPat(ctx *BindingPatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIsPat(ctx *IsPatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExprPat(ctx *ExprPatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumCasePat(ctx *EnumCasePatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTuplePat(ctx *TuplePatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitOptPat(ctx *OptPatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIdentPat(ctx *IdentPatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAsPat(ctx *AsPatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWildcardPat(ctx *WildcardPatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWildcardPattern(ctx *WildcardPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIdentifierPattern(ctx *IdentifierPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitValueBindingPattern(ctx *ValueBindingPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTuplePattern(ctx *TuplePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTuplePatternElementList(ctx *TuplePatternElementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTuplePatternElement(ctx *TuplePatternElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumCasePattern(ctx *EnumCasePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitOptionalPattern(ctx *OptionalPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExpressionPattern(ctx *ExpressionPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAttributes(ctx *AttributesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAttribute(ctx *AttributeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAttributeArguments(ctx *AttributeArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAttributeArgumentList(ctx *AttributeArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAttributeArgument(ctx *AttributeArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAvailabilityCondition(ctx *AvailabilityConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAvailabilityArguments(ctx *AvailabilityArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAvailabilityArgument(ctx *AvailabilityArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIdentifier(ctx *IdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}
