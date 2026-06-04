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
