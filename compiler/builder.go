package compiler

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"
	"github.com/vertex-language/vertex/parser"
)

// ─────────────────────────────────────────────────────────────────────────────
// ASTBuilder — converts the ANTLR CST into our AST.
// ─────────────────────────────────────────────────────────────────────────────

type ASTBuilder struct {
	filename string
	diags    *Diagnostics
}

func newASTBuilder(filename string, diags *Diagnostics) *ASTBuilder {
	return &ASTBuilder{filename: filename, diags: diags}
}

// pos converts an ANTLR token to a source Pos.
func (b *ASTBuilder) pos(tok antlr.Token) Pos {
	if tok == nil {
		return Pos{File: b.filename}
	}
	return Pos{File: b.filename, Line: tok.GetLine(), Column: tok.GetColumn() + 1}
}

func (b *ASTBuilder) ctxPos(ctx antlr.ParserRuleContext) Pos {
	return b.pos(ctx.GetStart())
}

// ─────────────────────────────────────────────────────────────────────────────
// File
// ─────────────────────────────────────────────────────────────────────────────

func (b *ASTBuilder) BuildFile(ctx parser.IFileContext) *File {
	file := &File{Pos: b.ctxPos(ctx)}
	if pkg := ctx.PackageDecl(); pkg != nil {
		file.Package = pkg.IDENTIFIER().GetText()
	}
	for _, imp := range ctx.AllImportDecl() {
		file.Imports = append(file.Imports, b.buildImport(imp))
	}
	for _, tld := range ctx.AllTopLevelDecl() {
		if d := b.buildTopLevelDecl(tld); d != nil {
			file.Decls = append(file.Decls, d)
		}
	}
	return file
}

func (b *ASTBuilder) buildImport(ctx parser.IImportDeclContext) *ImportDecl {
	pos := b.ctxPos(ctx)
	// Grouped import: import ( "a" "b" ) → emit multiple ImportDecl.
	// For simplicity we return the first; the caller loops AllImportDecl.
	toks := ctx.AllSTRING_LIT()
	if len(toks) == 0 {
		return &ImportDecl{Pos: pos}
	}
	return &ImportDecl{Pos: pos, Path: unquote(toks[0].GetText())}
}

func (b *ASTBuilder) buildTopLevelDecl(ctx parser.ITopLevelDeclContext) Decl {
	switch {
	case ctx.FuncDecl() != nil:
		return b.buildFuncDecl(ctx.FuncDecl())
	case ctx.StructDecl() != nil:
		return b.buildStructDecl(ctx.StructDecl())
	case ctx.ClassDecl() != nil:
		return b.buildClassDecl(ctx.ClassDecl())
	case ctx.EnumDecl() != nil:
		return b.buildEnumDecl(ctx.EnumDecl())
	case ctx.TypeAliasDecl() != nil:
		return b.buildTypeAlias(ctx.TypeAliasDecl())
	case ctx.VarDecl() != nil:
		return b.buildVarDecl(ctx.VarDecl())
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Functions
// ─────────────────────────────────────────────────────────────────────────────

func (b *ASTBuilder) buildFuncDecl(ctx parser.IFuncDeclContext) *FuncDecl {
	pos := b.ctxPos(ctx)
	fn := &FuncDecl{
		Pos:  pos,
		Name: ctx.IDENTIFIER().GetText(),
	}
	if ctx.Receiver() != nil {
		fn.Receiver = b.buildReceiver(ctx.Receiver())
	}
	if ctx.GenericParams() != nil {
		for _, id := range ctx.GenericParams().AllIDENTIFIER() {
			fn.TypeParams = append(fn.TypeParams, id.GetText())
		}
	}
	if ctx.ParamList() != nil {
		fn.Params = b.buildParamList(ctx.ParamList())
	}
	if ctx.FuncQualifier() != nil {
		fn.Qualifier = b.buildFuncQual(ctx.FuncQualifier())
	}
	if ctx.TypeExpr() != nil {
		fn.RetType = b.buildTypeExpr(ctx.TypeExpr())
	}
	fn.Body = b.buildBlock(ctx.Block())
	return fn
}

func (b *ASTBuilder) buildReceiver(ctx parser.IReceiverContext) *Receiver {
	return &Receiver{
		Pos:  b.ctxPos(ctx),
		Name: ctx.IDENTIFIER().GetText(),
		Type: b.buildTypeExpr(ctx.TypeExpr()),
	}
}

func (b *ASTBuilder) buildFuncQual(ctx parser.IFuncQualifierContext) FuncQual {
	switch {
	case ctx.ASYNC() != nil:
		return FuncQualAsync
	case ctx.THREAD() != nil:
		return FuncQualThread
	case ctx.PROCESS() != nil:
		return FuncQualProcess
	case ctx.GPU() != nil:
		return FuncQualGPU
	}
	return FuncQualNone
}

func (b *ASTBuilder) buildParamList(ctx parser.IParamListContext) []*Param {
	var params []*Param
	for _, p := range ctx.AllParam() {
		params = append(params, &Param{
			Pos:  b.ctxPos(p),
			Name: p.IDENTIFIER().GetText(),
			Type: b.buildTypeExpr(p.TypeExpr()),
		})
	}
	return params
}

// ─────────────────────────────────────────────────────────────────────────────
// Struct / Class / Enum / TypeAlias / Var
// ─────────────────────────────────────────────────────────────────────────────

func (b *ASTBuilder) buildStructDecl(ctx parser.IStructDeclContext) *StructDecl {
	s := &StructDecl{Pos: b.ctxPos(ctx), Name: ctx.IDENTIFIER().GetText()}
	for _, f := range ctx.AllStructFieldDecl() {
		s.Fields = append(s.Fields, &FieldDecl{
			Pos:  b.ctxPos(f),
			Name: f.IDENTIFIER().GetText(),
			Type: b.buildTypeExpr(f.TypeExpr()),
		})
	}
	return s
}

func (b *ASTBuilder) buildClassDecl(ctx parser.IClassDeclContext) *ClassDecl {
	cl := &ClassDecl{Pos: b.ctxPos(ctx), Name: ctx.IDENTIFIER().GetText()}
	if ctx.QualifiedIdent() != nil {
		parts := make([]string, 0, len(ctx.QualifiedIdent().AllIDENTIFIER()))
		for _, id := range ctx.QualifiedIdent().AllIDENTIFIER() {
			parts = append(parts, id.GetText())
		}
		cl.BaseName = strings.Join(parts, ".")
	}
	for _, m := range ctx.AllClassMember() {
		mem := &ClassMember{Pos: b.ctxPos(m), Name: m.IDENTIFIER().GetText()}
		if m.FUNC() != nil {
			// Native method signature
			mem.IsField = false
			if m.ParamList() != nil {
				mem.Params = b.buildParamList(m.ParamList())
			}
			if m.TypeExpr() != nil {
				mem.RetType = b.buildTypeExpr(m.TypeExpr())
			}
		} else {
			// Field
			mem.IsField = true
			mem.Type = b.buildTypeExpr(m.TypeExpr())
		}
		cl.Members = append(cl.Members, mem)
	}
	return cl
}

func (b *ASTBuilder) buildEnumDecl(ctx parser.IEnumDeclContext) *EnumDecl {
	e := &EnumDecl{Pos: b.ctxPos(ctx), Name: ctx.IDENTIFIER().GetText()}
	if ctx.TypeExpr() != nil {
		e.RawType = b.buildTypeExpr(ctx.TypeExpr())
	}
	for _, cd := range ctx.AllEnumCaseDecl() {
		for _, ci := range cd.AllEnumCaseItem() {
			c := &EnumCase{Pos: b.ctxPos(ci), Name: ci.IDENTIFIER().GetText()}
			if ci.EnumRawValue() != nil {
				c.RawValue = b.buildEnumRawValue(ci.EnumRawValue())
			}
			e.Cases = append(e.Cases, c)
		}
	}
	return e
}

func (b *ASTBuilder) buildEnumRawValue(ctx parser.IEnumRawValueContext) Expr {
	pos := b.ctxPos(ctx)
	neg := ctx.MINUS() != nil
	switch {
	case ctx.DEC_INT_LIT() != nil:
		v := parseDecInt(ctx.DEC_INT_LIT().GetText())
		if neg {
			v = -v
		}
		return &IntLitExpr{exprBase: exprBase{Pos: pos}, Value: v}
	case ctx.HEX_INT_LIT() != nil:
		v := parseHexInt(ctx.HEX_INT_LIT().GetText())
		if neg {
			v = -v
		}
		return &IntLitExpr{exprBase: exprBase{Pos: pos}, Value: v}
	case ctx.STRING_LIT() != nil:
		return &StringLitExpr{exprBase: exprBase{Pos: pos}, Value: unquote(ctx.STRING_LIT().GetText())}
	}
	return &IntLitExpr{exprBase: exprBase{Pos: pos}}
}

func (b *ASTBuilder) buildTypeAlias(ctx parser.ITypeAliasDeclContext) *TypeAliasDecl {
	return &TypeAliasDecl{
		Pos:  b.ctxPos(ctx),
		Name: ctx.IDENTIFIER().GetText(),
		Type: b.buildTypeExpr(ctx.TypeExpr()),
	}
}

func (b *ASTBuilder) buildVarDecl(ctx parser.IVarDeclContext) *VarDecl {
	pos := b.ctxPos(ctx)
	vd := &VarDecl{
		Pos:    pos,
		IsLet:  ctx.LET() != nil,
		IsWeak: ctx.WEAK() != nil,
	}
	vd.Binding = b.buildBindingPattern(ctx.BindingPattern())
	if ctx.TypeExpr() != nil {
		vd.TypeHint = b.buildTypeExpr(ctx.TypeExpr())
	}
	vd.Value = b.buildExpr(ctx.Expr())
	return vd
}

func (b *ASTBuilder) buildBindingPattern(ctx parser.IBindingPatternContext) *BindingPattern {
	bp := &BindingPattern{Pos: b.ctxPos(ctx)}
	for _, id := range ctx.AllIDENTIFIER() {
		bp.Names = append(bp.Names, id.GetText())
	}
	return bp
}

// ─────────────────────────────────────────────────────────────────────────────
// Type expressions
// ─────────────────────────────────────────────────────────────────────────────

func (b *ASTBuilder) buildTypeExpr(ctx parser.ITypeExprContext) TypeExpr {
	if ctx == nil {
		return nil
	}
	pos := b.ctxPos(ctx)
	subExprs := ctx.AllTypeExpr()

	switch {
	case ctx.STAR() != nil:
		// *T, *const T, *T?, *const T?
		var elem TypeExpr
		if len(subExprs) > 0 {
			elem = b.buildTypeExpr(subExprs[0])
		}
		return &PointerTypeExpr{
			Pos:      pos,
			IsConst:  ctx.CONST_KW() != nil,
			Elem:     elem,
			Optional: ctx.QUESTION() != nil,
		}
	case ctx.LBRACKET() != nil && ctx.RBRACKET() != nil:
		// [T]
		var elem TypeExpr
		if len(subExprs) > 0 {
			elem = b.buildTypeExpr(subExprs[0])
		}
		return &ArrayTypeExpr{Pos: pos, Elem: elem}
	case ctx.CHAN() != nil:
		var elem TypeExpr
		if len(subExprs) > 0 {
			elem = b.buildTypeExpr(subExprs[0])
		}
		return &ChanTypeExpr{Pos: pos, Elem: elem}
	case ctx.FUNC() != nil:
		te := &FuncTypeExpr{Pos: pos}
		if ctx.FuncTypeParams() != nil {
			for _, p := range ctx.FuncTypeParams().AllTypeExpr() {
				te.Params = append(te.Params, b.buildTypeExpr(p))
			}
		}
		if ctx.ARROW() != nil && len(subExprs) > 0 {
			te.RetType = b.buildTypeExpr(subExprs[len(subExprs)-1])
		}
		return te
	case ctx.LPAREN() != nil && ctx.RPAREN() != nil:
		// () or (T,T) tuple
		tt := &TupleTypeExpr{Pos: pos}
		if ctx.TupleTypeElems() != nil {
			for _, te := range ctx.TupleTypeElems().AllTupleTypeElem() {
				var label string
				if te.IDENTIFIER() != nil && te.COLON() != nil {
					label = te.IDENTIFIER().GetText()
				}
				tt.Labels = append(tt.Labels, label)
				tt.Elems = append(tt.Elems, b.buildTypeExpr(te.TypeExpr()))
			}
		}
		return tt
	case ctx.RESULT() != nil:
		// Result(T, E)
		rt := &ResultTypeExpr{Pos: pos}
		if len(subExprs) >= 2 {
			rt.Ok = b.buildTypeExpr(subExprs[0])
			rt.Err = b.buildTypeExpr(subExprs[1])
		}
		return rt
	case ctx.BaseType() != nil:
		// Named type, possibly optional
		bt := ctx.BaseType()
		var parts []string
		for _, id := range bt.AllIDENTIFIER() {
			parts = append(parts, id.GetText())
		}
		var pkg, name string
		if len(parts) >= 2 {
			pkg = strings.Join(parts[:len(parts)-1], ".")
			name = parts[len(parts)-1]
		} else if len(parts) == 1 {
			name = parts[0]
		}
		named := &NamedTypeExpr{Pos: pos, Pkg: pkg, Name: name}
		if ctx.QUESTION() != nil {
			return &OptionalTypeExpr{Pos: pos, Elem: named}
		}
		return named
	}
	return &NamedTypeExpr{Pos: pos, Name: "void"}
}

// ─────────────────────────────────────────────────────────────────────────────
// Statements
// ─────────────────────────────────────────────────────────────────────────────

func (b *ASTBuilder) buildBlock(ctx parser.IBlockContext) *BlockStmt {
	if ctx == nil {
		return &BlockStmt{}
	}
	blk := &BlockStmt{Pos: b.ctxPos(ctx)}
	for _, s := range ctx.AllStmt() {
		if st := b.buildStmt(s); st != nil {
			blk.Stmts = append(blk.Stmts, st)
		}
	}
	return blk
}

func (b *ASTBuilder) buildStmt(ctx parser.IStmtContext) Stmt {
	pos := b.ctxPos(ctx)
	switch {
	case ctx.VarDecl() != nil:
		return &LocalDeclStmt{Pos: pos, Decl: b.buildVarDecl(ctx.VarDecl())}
	case ctx.IfStmt() != nil:
		return b.buildIfStmt(ctx.IfStmt())
	case ctx.WhileStmt() != nil:
		return b.buildWhileStmt(ctx.WhileStmt())
	case ctx.ForInStmt() != nil:
		return b.buildForInStmt(ctx.ForInStmt())
	case ctx.SwitchStmt() != nil:
		return b.buildSwitchStmt(ctx.SwitchStmt())
	case ctx.ReturnStmt() != nil:
		rs := ctx.ReturnStmt()
		ret := &ReturnStmt{Pos: b.ctxPos(rs)}
		if rs.Expr() != nil {
			ret.Value = b.buildExpr(rs.Expr())
		}
		return ret
	case ctx.DeferStmt() != nil:
		ds := ctx.DeferStmt()
		return &DeferStmt{Pos: b.ctxPos(ds), Call: b.buildExpr(ds.Expr())}
	case ctx.BREAK() != nil:
		return &BreakStmt{Pos: pos}
	case ctx.CONTINUE() != nil:
		return &ContinueStmt{Pos: pos}
	case ctx.FALLTHROUGH() != nil:
		return &FallthroughStmt{Pos: pos}
	case ctx.ExprOrAssignStmt() != nil:
		return b.buildExprOrAssign(ctx.ExprOrAssignStmt())
	case ctx.StructDecl() != nil:
		// struct/class/enum inside a function body — wrap as a local type decl.
		// Lowerer treats these as type registrations.
		return &ExprStmt{Pos: pos, Expr: &NilLitExpr{exprBase: exprBase{Pos: pos}}}
	}
	return nil
}

func (b *ASTBuilder) buildIfStmt(ctx parser.IIfStmtContext) *IfStmt {
	s := &IfStmt{Pos: b.ctxPos(ctx)}
	ic := ctx.IfCondition()
	if ic.LET() != nil {
		s.Cond = &IfLetCond{
			Pos:  b.ctxPos(ic),
			Name: ic.IDENTIFIER().GetText(),
			Expr: b.buildExpr(ic.Expr()),
		}
	} else {
		s.Cond = &IfExprCond{Pos: b.ctxPos(ic), Expr: b.buildExpr(ic.Expr())}
	}
	blocks := ctx.AllBlock()
	if len(blocks) > 0 {
		s.Then = b.buildBlock(blocks[0])
	}
	if ctx.ELSE() != nil {
		if ctx.IfStmt() != nil {
			s.Else = b.buildIfStmt(ctx.IfStmt())
		} else if len(blocks) > 1 {
			s.Else = b.buildBlock(blocks[1])
		}
	}
	return s
}

func (b *ASTBuilder) buildWhileStmt(ctx parser.IWhileStmtContext) *WhileStmt {
	return &WhileStmt{
		Pos:  b.ctxPos(ctx),
		Cond: b.buildExpr(ctx.Expr()),
		Body: b.buildBlock(ctx.Block()),
	}
}

func (b *ASTBuilder) buildForInStmt(ctx parser.IForInStmtContext) *ForInStmt {
	return &ForInStmt{
		Pos:  b.ctxPos(ctx),
		Var:  ctx.IDENTIFIER().GetText(),
		Iter: b.buildExpr(ctx.Expr()),
		Body: b.buildBlock(ctx.Block()),
	}
}

func (b *ASTBuilder) buildSwitchStmt(ctx parser.ISwitchStmtContext) *SwitchStmt {
	s := &SwitchStmt{Pos: b.ctxPos(ctx), Subj: b.buildExpr(ctx.Expr())}
	for _, c := range ctx.AllSwitchCase() {
		sc := &SwitchCase{Pos: b.ctxPos(c), IsDefault: c.DEFAULT() != nil}
		for _, p := range c.AllSwitchPattern() {
			sc.Patterns = append(sc.Patterns, b.buildSwitchPattern(p))
		}
		for _, st := range c.AllStmt() {
			if built := b.buildStmt(st); built != nil {
				sc.Body = append(sc.Body, built)
			}
		}
		s.Cases = append(s.Cases, sc)
	}
	return s
}

func (b *ASTBuilder) buildSwitchPattern(ctx parser.ISwitchPatternContext) SwitchPattern {
	pos := b.ctxPos(ctx)
	switch {
	case ctx.DOT() != nil && ctx.IDENTIFIER() != nil:
		return &EnumShortPattern{Pos: pos, Case: ctx.IDENTIFIER().GetText()}
	case ctx.OK() != nil:
		return &ResultOkPattern{Pos: pos, Bind: ctx.IDENTIFIER().GetText()}
	case ctx.ERR_KW() != nil:
		return &ResultErrPattern{Pos: pos, Bind: ctx.IDENTIFIER().GetText()}
	default:
		return &ExprPattern{Pos: pos, Expr: b.buildExpr(ctx.Expr())}
	}
}

func (b *ASTBuilder) buildExprOrAssign(ctx parser.IExprOrAssignStmtContext) Stmt {
	pos := b.ctxPos(ctx)
	exprs := ctx.AllExpr()
	if len(exprs) == 1 {
		return &ExprStmt{Pos: pos, Expr: b.buildExpr(exprs[0])}
	}
	if len(exprs) == 2 {
		op := OpAssign
		if ctx.AssignOp() != nil {
			op = b.buildAssignOp(ctx.AssignOp())
		}
		return &AssignStmt{
			Pos: pos,
			LHS: b.buildExpr(exprs[0]),
			Op:  op,
			RHS: b.buildExpr(exprs[1]),
		}
	}
	return &ExprStmt{Pos: pos, Expr: &NilLitExpr{exprBase: exprBase{Pos: pos}}}
}

func (b *ASTBuilder) buildAssignOp(ctx parser.IAssignOpContext) AssignOp {
	switch {
	case ctx.PLUS_ASSIGN() != nil:
		return OpAddAssign
	case ctx.MINUS_ASSIGN() != nil:
		return OpSubAssign
	case ctx.STAR_ASSIGN() != nil:
		return OpMulAssign
	case ctx.DIV_ASSIGN() != nil:
		return OpDivAssign
	case ctx.MOD_ASSIGN() != nil:
		return OpModAssign
	}
	return OpAssign
}

// ─────────────────────────────────────────────────────────────────────────────
// Expressions — disambiguation of the left-recursive expr rule.
//
// Strategy: inspect which terminal nodes are present AND the first token of the
// context. The first-token check distinguishes prefix/primary forms (where the
// leading token is the operator or delimiter) from postfix/binary forms (where
// the leading token is the start of the left operand expression).
// ─────────────────────────────────────────────────────────────────────────────

func (b *ASTBuilder) buildExpr(ctx parser.IExprContext) Expr {
	if ctx == nil {
		return &NilLitExpr{}
	}
	pos := b.ctxPos(ctx)
	exprs := ctx.AllExpr()
	nExprs := len(exprs)
	startTok := ctx.GetStart().GetTokenType()

	// ── Leaf / unique primary forms ──────────────────────────────────────────
	if ctx.AnonFuncExpr() != nil {
		return b.buildAnonFunc(ctx.AnonFuncExpr())
	}
	if ctx.Literal() != nil {
		return b.buildLiteral(pos, ctx.Literal())
	}
	if ctx.RESULT() != nil {
		return b.buildResultExpr(pos, ctx)
	}
	if ctx.AsmExpr() != nil {
		b.diags.Warnf(pos, "inline asm is not supported; expression replaced with nil")
		return &NilLitExpr{exprBase: exprBase{Pos: pos}}
	}

	// ── Postfix: leading child is an expr (start token is NOT an operator) ───
	isPrefix := startTok == parser.VertexLexerMINUS ||
		startTok == parser.VertexLexerBANG ||
		startTok == parser.VertexLexerTILDE ||
		startTok == parser.VertexLexerAMP ||
		startTok == parser.VertexLexerLPAREN ||
		startTok == parser.VertexLexerLBRACKET ||
		startTok == parser.VertexLexerDOT

	if !isPrefix && nExprs >= 1 {
		// Method call: expr.method(args)
		if ctx.DOT() != nil && ctx.IDENTIFIER() != nil && ctx.LPAREN() != nil {
			return &MethodCallExpr{
				exprBase: exprBase{Pos: pos},
				Recv:     b.buildExpr(exprs[0]),
				Method:   ctx.IDENTIFIER().GetText(),
				Args:     b.buildArgList(ctx.ArgList()),
			}
		}
		// Field access: expr.field
		if ctx.DOT() != nil && ctx.IDENTIFIER() != nil {
			return &FieldExpr{
				exprBase: exprBase{Pos: pos},
				Recv:     b.buildExpr(exprs[0]),
				Field:    ctx.IDENTIFIER().GetText(),
			}
		}
		// Subscript: expr[expr]
		if ctx.LBRACKET() != nil && ctx.RBRACKET() != nil && nExprs == 2 {
			return &IndexExpr{
				exprBase: exprBase{Pos: pos},
				Recv:     b.buildExpr(exprs[0]),
				Index:    b.buildExpr(exprs[1]),
			}
		}
		// Function call: expr(args)
		if ctx.LPAREN() != nil && ctx.RPAREN() != nil && nExprs == 1 {
			inner := b.buildExpr(exprs[0])
			if arrLit, ok := inner.(*ArrayLitExpr); ok && len(arrLit.Elems) == 1 {
				if id, ok2 := arrLit.Elems[0].(*IdentExpr); ok2 {
					return &ArrayCtorExpr{
						exprBase:     exprBase{Pos: pos},
						ElemTypeName: id.Name,
						Args:         b.buildArgList(ctx.ArgList()),
					}
				}
			}
			return &CallExpr{
				exprBase: exprBase{Pos: pos},
				Func:     inner,
				Args:     b.buildArgList(ctx.ArgList()),
			}
		}
	}

	// ── Binary operators (exactly two child exprs) ────────────────────────────
	if nExprs == 2 {
		if op, ok := b.detectBinOp(ctx); ok {
			return &BinaryExpr{
				exprBase: exprBase{Pos: pos},
				Op:       op,
				Left:     b.buildExpr(exprs[0]),
				Right:    b.buildExpr(exprs[1]),
			}
		}
	}

	// ── Ternary (three child exprs) ───────────────────────────────────────────
	if nExprs == 3 && ctx.QUESTION() != nil && ctx.COLON() != nil {
		return &TernaryExpr{
			exprBase: exprBase{Pos: pos},
			Cond:     b.buildExpr(exprs[0]),
			Then:     b.buildExpr(exprs[1]),
			Else:     b.buildExpr(exprs[2]),
		}
	}

	// ── Prefix unary ──────────────────────────────────────────────────────────
	if nExprs == 1 && (startTok == parser.VertexLexerMINUS ||
		startTok == parser.VertexLexerBANG ||
		startTok == parser.VertexLexerTILDE) {
		op := UnNeg
		if startTok == parser.VertexLexerBANG {
			op = UnNot
		} else if startTok == parser.VertexLexerTILDE {
			op = UnBitNot
		}
		return &UnaryExpr{exprBase: exprBase{Pos: pos}, Op: op, Operand: b.buildExpr(exprs[0])}
	}
	if nExprs == 1 && startTok == parser.VertexLexerAMP {
		return &UnaryExpr{exprBase: exprBase{Pos: pos}, Op: UnAddrOf, Operand: b.buildExpr(exprs[0])}
	}

	// ── Parenthesised / tuple ─────────────────────────────────────────────────
	if startTok == parser.VertexLexerLPAREN && ctx.LPAREN() != nil && ctx.RPAREN() != nil {
		if nExprs == 0 {
			return &TupleLitExpr{exprBase: exprBase{Pos: pos}}
		}
		if nExprs == 1 && len(ctx.AllCOMMA()) == 0 {
			return b.buildExpr(exprs[0]) // grouping
		}
		var elems []Expr
		for _, e := range exprs {
			elems = append(elems, b.buildExpr(e))
		}
		return &TupleLitExpr{exprBase: exprBase{Pos: pos}, Elems: elems}
	}

	// ── Array literal [a, b, …] ───────────────────────────────────────────────
	if startTok == parser.VertexLexerLBRACKET && ctx.LBRACKET() != nil {
		var elems []Expr
		for _, e := range exprs {
			elems = append(elems, b.buildExpr(e))
		}
		return &ArrayLitExpr{exprBase: exprBase{Pos: pos}, Elems: elems}
	}

	// ── Map literal {"k": v} ──────────────────────────────────────────────────
	if startTok == parser.VertexLexerLBRACE && ctx.LBRACE() != nil {
		return &MapLitExpr{
			exprBase: exprBase{Pos: pos},
			Fields:   b.buildMapLitFields(ctx.MapLiteralFields()),
		}
	}

	// ── Enum shorthand .CaseName ──────────────────────────────────────────────
	if startTok == parser.VertexLexerDOT && ctx.IDENTIFIER() != nil && nExprs == 0 {
		return &DotEnumExpr{exprBase: exprBase{Pos: pos}, Case: ctx.IDENTIFIER().GetText()}
	}

	// ── Identifier, possibly followed by struct literal ───────────────────────
	if ctx.IDENTIFIER() != nil && nExprs == 0 {
		name := ctx.IDENTIFIER().GetText()
		if ctx.LBRACE() != nil {
			return &StructLitExpr{
				exprBase: exprBase{Pos: pos},
				TypeName: name,
				Fields:   b.buildStructLitFields(ctx.StructLiteralFields()),
			}
		}
		return &IdentExpr{exprBase: exprBase{Pos: pos}, Name: name}
	}

	b.diags.Errorf(pos, "unrecognised expression form (start token=%d, nExprs=%d)", startTok, nExprs)
	return &NilLitExpr{exprBase: exprBase{Pos: pos}}
}

// Add this helper function below buildArgList
func (b *ASTBuilder) buildMapLitFields(ctx parser.IMapLiteralFieldsContext) []*MapLitField {
	if ctx == nil {
		return nil
	}
	var fields []*MapLitField
	for _, f := range ctx.AllMapLiteralField() {
		exprs := f.AllExpr()
		if len(exprs) == 2 {
			fields = append(fields, &MapLitField{
				Pos:   b.ctxPos(f),
				Key:   b.buildExpr(exprs[0]),
				Value: b.buildExpr(exprs[1]),
			})
		}
	}
	return fields
}

func (b *ASTBuilder) detectBinOp(ctx parser.IExprContext) (BinOp, bool) {
	switch {
	case ctx.PLUS() != nil:
		return BinAdd, true
	case ctx.MINUS() != nil:
		return BinSub, true
	case ctx.STAR() != nil:
		return BinMul, true
	case ctx.SLASH() != nil:
		return BinDiv, true
	case ctx.PERCENT() != nil:
		return BinMod, true
	case ctx.LSHIFT() != nil:
		return BinShl, true
	case ctx.RSHIFT() != nil:
		return BinShr, true
	case ctx.EQ() != nil:
		return BinEq, true
	case ctx.NEQ() != nil:
		return BinNeq, true
	case ctx.LT() != nil:
		return BinLt, true
	case ctx.GT() != nil:
		return BinGt, true
	case ctx.LEQ() != nil:
		return BinLte, true
	case ctx.GEQ() != nil:
		return BinGte, true
	case ctx.LOGICAL_AND() != nil:
		return BinAnd, true
	case ctx.LOGICAL_OR() != nil:
		return BinOr, true
	case ctx.NIL_COALESCE() != nil:
		return BinNilCoalesce, true
	case ctx.OVERFLOW_ADD() != nil:
		return BinOverflowAdd, true
	case ctx.OVERFLOW_SUB() != nil:
		return BinOverflowSub, true
	case ctx.OVERFLOW_MUL() != nil:
		return BinOverflowMul, true
	case ctx.ELLIPSIS() != nil:
		return BinRangeClosed, true
	case ctx.HALF_OPEN() != nil:
		return BinRangeHalfOpen, true
	case ctx.IDENTITY_EQ() != nil:
		return BinIdentityEq, true
	case ctx.IDENTITY_NEQ() != nil:
		return BinIdentityNeq, true
	}
	return 0, false
}

func (b *ASTBuilder) buildLiteral(pos Pos, ctx parser.ILiteralContext) Expr {
	switch {
	case ctx.TRUE() != nil:
		return &BoolLitExpr{exprBase: exprBase{Pos: pos}, Value: true}
	case ctx.FALSE() != nil:
		return &BoolLitExpr{exprBase: exprBase{Pos: pos}, Value: false}
	case ctx.NIL() != nil:
		return &NilLitExpr{exprBase: exprBase{Pos: pos}}
	case ctx.DEC_INT_LIT() != nil:
		return &IntLitExpr{exprBase: exprBase{Pos: pos}, Value: parseDecInt(ctx.DEC_INT_LIT().GetText())}
	case ctx.HEX_INT_LIT() != nil:
		return &IntLitExpr{exprBase: exprBase{Pos: pos}, Value: parseHexInt(ctx.HEX_INT_LIT().GetText())}
	case ctx.OCT_INT_LIT() != nil:
		return &IntLitExpr{exprBase: exprBase{Pos: pos}, Value: parseOctInt(ctx.OCT_INT_LIT().GetText())}
	case ctx.BIN_INT_LIT() != nil:
		return &IntLitExpr{exprBase: exprBase{Pos: pos}, Value: parseBinInt(ctx.BIN_INT_LIT().GetText())}
	case ctx.DEC_FLOAT_LIT() != nil:
		return &FloatLitExpr{exprBase: exprBase{Pos: pos}, Value: parseFloat(ctx.DEC_FLOAT_LIT().GetText())}
	case ctx.HEX_FLOAT_LIT() != nil:
		return &FloatLitExpr{exprBase: exprBase{Pos: pos}, Value: parseFloat(ctx.HEX_FLOAT_LIT().GetText())}
	case ctx.STRING_LIT() != nil:
		return &StringLitExpr{exprBase: exprBase{Pos: pos}, Value: unquote(ctx.STRING_LIT().GetText())}
	case ctx.MULTILINE_STRING_LIT() != nil:
		raw := ctx.MULTILINE_STRING_LIT().GetText()
		// Strip surrounding backticks.
		if len(raw) >= 2 {
			raw = raw[1 : len(raw)-1]
		}
		return &StringLitExpr{exprBase: exprBase{Pos: pos}, Value: raw}
	}
	return &NilLitExpr{exprBase: exprBase{Pos: pos}}
}

func (b *ASTBuilder) buildResultExpr(pos Pos, ctx parser.IExprContext) Expr {
	isOk := ctx.OK() != nil
	var val Expr
	exprs := ctx.AllExpr()
	if len(exprs) > 0 {
		val = b.buildExpr(exprs[0])
	}
	return &ResultExpr{exprBase: exprBase{Pos: pos}, IsOk: isOk, Value: val}
}

func (b *ASTBuilder) buildStructLitFields(ctx parser.IStructLiteralFieldsContext) []*StructLitField {
	if ctx == nil {
		return nil
	}
	var fields []*StructLitField
	for _, f := range ctx.AllStructLiteralField() {
		fields = append(fields, &StructLitField{
			Pos:   b.ctxPos(f),
			Name:  f.IDENTIFIER().GetText(),
			Value: b.buildExpr(f.Expr()),
		})
	}
	return fields
}

func (b *ASTBuilder) buildArgList(ctx parser.IArgListContext) []*Arg {
	if ctx == nil {
		return nil
	}
	var args []*Arg
	for _, a := range ctx.AllArg() {
		arg := &Arg{Pos: b.ctxPos(a), Value: b.buildExpr(a.Expr())}
		if a.IDENTIFIER() != nil && a.COLON() != nil {
			arg.Label = a.IDENTIFIER().GetText()
		}
		args = append(args, arg)
	}
	return args
}

func (b *ASTBuilder) buildAnonFunc(ctx parser.IAnonFuncExprContext) *AnonFuncExpr {
	fn := &AnonFuncExpr{exprBase: exprBase{Pos: b.ctxPos(ctx)}}
	if ctx.ParamList() != nil {
		fn.Params = b.buildParamList(ctx.ParamList())
	}
	if ctx.FuncQualifier() != nil {
		fn.Qualifier = b.buildFuncQual(ctx.FuncQualifier())
	}
	if ctx.TypeExpr() != nil {
		fn.RetType = b.buildTypeExpr(ctx.TypeExpr())
	}
	fn.Body = b.buildBlock(ctx.Block())
	return fn
}

// ─────────────────────────────────────────────────────────────────────────────
// Numeric / string helpers
// ─────────────────────────────────────────────────────────────────────────────

func stripUnderscores(s string) string {
	return strings.ReplaceAll(s, "_", "")
}

func parseDecInt(s string) int64 {
	v, _ := strconv.ParseInt(stripUnderscores(s), 10, 64)
	return v
}

func parseHexInt(s string) int64 {
	// Remove 0x/0X prefix.
	s = stripUnderscores(s)
	if len(s) > 2 {
		s = s[2:]
	}
	v, _ := strconv.ParseInt(s, 16, 64)
	return v
}

func parseOctInt(s string) int64 {
	s = stripUnderscores(s)
	if len(s) > 2 {
		s = s[2:]
	}
	v, _ := strconv.ParseInt(s, 8, 64)
	return v
}

func parseBinInt(s string) int64 {
	s = stripUnderscores(s)
	if len(s) > 2 {
		s = s[2:]
	}
	v, _ := strconv.ParseInt(s, 2, 64)
	return v
}

func parseFloat(s string) float64 {
	s = stripUnderscores(s)
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// unquote strips surrounding double-quotes and expands standard escape sequences.
func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	var buf strings.Builder
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\\' && i+1 < len(s) {
			i++
			switch s[i] {
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case 't':
				buf.WriteByte('\t')
			case 'b':
				buf.WriteByte('\b')
			case 'f':
				buf.WriteByte('\f')
			case '"':
				buf.WriteByte('"')
			case '\\':
				buf.WriteByte('\\')
			default:
				buf.WriteByte('\\')
				buf.WriteByte(s[i])
			}
			i++
		} else {
			buf.WriteRune(r)
			i += size
		}
	}
	return buf.String()
}

// isUpper reports whether a name is exported (starts with an upper-case letter).
// Currently unused but kept for future visibility rules.
func isUpper(name string) bool {
	if name == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}