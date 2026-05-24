// Code generated from VertexParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // VertexParser

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by VertexParser.
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
