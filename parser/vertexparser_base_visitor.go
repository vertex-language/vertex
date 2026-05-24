// Code generated from VertexParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // VertexParser

import "github.com/antlr4-go/antlr/v4"

type BaseVertexParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseVertexParserVisitor) VisitFile(ctx *FileContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPackageDecl(ctx *PackageDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBuildDecl(ctx *BuildDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBuildTag(ctx *BuildTagContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitImportDecl(ctx *ImportDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFuncDecl(ctx *FuncDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGenericParams(ctx *GenericParamsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeParam(ctx *TypeParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitParamList(ctx *ParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitParam(ctx *ParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFuncQualifier(ctx *FuncQualifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitReturnType(ctx *ReturnTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructDecl(ctx *StructDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructField(ctx *StructFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructLiteralExpr(ctx *StructLiteralExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructFieldInit(ctx *StructFieldInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClassDecl(ctx *ClassDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClassMember(ctx *ClassMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClassField(ctx *ClassFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitNativeFuncDecl(ctx *NativeFuncDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitNativeParamList(ctx *NativeParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitNativeParam(ctx *NativeParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumDecl(ctx *EnumDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumRawType(ctx *EnumRawTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumCaseDecl(ctx *EnumCaseDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumCase(ctx *EnumCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeAliasDecl(ctx *TypeAliasDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStmt(ctx *StmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitVarDeclStmt(ctx *VarDeclStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBindingKw(ctx *BindingKwContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleBind(ctx *TupleBindContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAssignStmt(ctx *AssignStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCompoundAssignStmt(ctx *CompoundAssignStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCompoundOp(ctx *CompoundOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLvalue(ctx *LvalueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIfStmt(ctx *IfStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitElseIfClause(ctx *ElseIfClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitElseClause(ctx *ElseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIfCondition(ctx *IfConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSwitchStmt(ctx *SwitchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDefaultCase(ctx *DefaultCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCasePatternList(ctx *CasePatternListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitCasePattern(ctx *CasePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitForInStmt(ctx *ForInStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWhileStmt(ctx *WhileStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBreakStmt(ctx *BreakStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitContinueStmt(ctx *ContinueStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFallthroughStmt(ctx *FallthroughStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitReturnStmt(ctx *ReturnStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDeferStmt(ctx *DeferStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExprStmt(ctx *ExprStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPostfixName(ctx *PostfixNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPrimary(ctx *PrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArgList(ctx *ArgListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArg(ctx *ArgContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArrayLiteralExpr(ctx *ArrayLiteralExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArrayConstructExpr(ctx *ArrayConstructExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExprList(ctx *ExprListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDictLiteralExpr(ctx *DictLiteralExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDictEntry(ctx *DictEntryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleExpr(ctx *TupleExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleElement(ctx *TupleElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAnonFuncExpr(ctx *AnonFuncExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitType(ctx *TypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleTypeElem(ctx *TupleTypeElemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFuncTypeParamList(ctx *FuncTypeParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFuncTypeParam(ctx *FuncTypeParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{} {
	return v.VisitChildren(ctx)
}
