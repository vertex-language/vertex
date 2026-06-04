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

func (v *BaseVertexParserVisitor) VisitImportDecl(ctx *ImportDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeAliasDecl(ctx *TypeAliasDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFuncDecl(ctx *FuncDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitReceiver(ctx *ReceiverContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitGenericParams(ctx *GenericParamsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFuncQualifier(ctx *FuncQualifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitParamList(ctx *ParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitParam(ctx *ParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitVariadicParam(ctx *VariadicParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructDecl(ctx *StructDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructFieldDecl(ctx *StructFieldDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClassDecl(ctx *ClassDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitClassMember(ctx *ClassMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitQualifiedIdent(ctx *QualifiedIdentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumDecl(ctx *EnumDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumCaseDecl(ctx *EnumCaseDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumCaseItem(ctx *EnumCaseItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitEnumRawValue(ctx *EnumRawValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStmt(ctx *StmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitVarDecl(ctx *VarDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBindingPattern(ctx *BindingPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIfStmt(ctx *IfStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitIfCondition(ctx *IfConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitWhileStmt(ctx *WhileStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitForInStmt(ctx *ForInStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSwitchStmt(ctx *SwitchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitSwitchPattern(ctx *SwitchPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitReturnStmt(ctx *ReturnStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitDeferStmt(ctx *DeferStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExprOrAssignStmt(ctx *ExprOrAssignStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAssignOp(ctx *AssignOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArgList(ctx *ArgListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitArg(ctx *ArgContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructLiteralFields(ctx *StructLiteralFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitStructLiteralField(ctx *StructLiteralFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitMapLiteralFields(ctx *MapLiteralFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitMapLiteralField(ctx *MapLiteralFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAnonFuncExpr(ctx *AnonFuncExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAsmExpr(ctx *AsmExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAsmBody(ctx *AsmBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAsmInstr(ctx *AsmInstrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAsmConstraint(ctx *AsmConstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitAsmClobberDecl(ctx *AsmClobberDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTypeExpr(ctx *TypeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitFuncTypeParams(ctx *FuncTypeParamsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitBaseType(ctx *BaseTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleTypeElems(ctx *TupleTypeElemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitTupleTypeElem(ctx *TupleTypeElemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVertexParserVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}
