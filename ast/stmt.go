package ast

import "github.com/vertex-language/vertex/token"

// Statements — vertex_grammar.md section C.

type (
	// BlockStmt is Block (C).
	BlockStmt struct {
		Lbrace token.Pos
		List   []Stmt
		Rbrace token.Pos
	}

	// EmptyStmt is a lone `;` (C).
	EmptyStmt struct{ Semi token.Pos }

	// ExprStmt is ExpressionStatement (C). Semi is NoPos when the terminator
	// was inserted (§6.2); the node's End is then the expression's End, so no
	// span covers a character that isn't there.
	ExprStmt struct {
		X    Expr
		Semi token.Pos
	}

	// IfStmt is IfStatement (C.3). Else is nil, a *BlockStmt, or another Stmt.
	IfStmt struct {
		If     token.Pos
		Lparen token.Pos
		Cond   Expr
		Rparen token.Pos
		Body   Stmt
		ElsePos token.Pos // NoPos if no else
		Else   Stmt
	}

	// DoWhileStmt is DoWhileStatement (C.3).
	DoWhileStmt struct {
		Do     token.Pos
		Body   Stmt
		While  token.Pos
		Cond   Expr
		Rparen token.Pos
		Semi   token.Pos // ASI-eligible
	}

	// WhileStmt is WhileStatement (C.3).
	WhileStmt struct {
		While  token.Pos
		Cond   Expr
		Rparen token.Pos
		Body   Stmt
	}

	// ForStmt is the three-clause ForStatement (C.3). Init is a Stmt because
	// it may be a VarDecl; any of the three may be nil.
	ForStmt struct {
		For    token.Pos
		Init   Stmt
		Cond   Expr
		Post   Expr
		Rparen token.Pos
		Body   Stmt
	}

	// ForInStmt is the `in` form of ForInOfStatement (C.3). Left is a
	// *VarDecl or an Expr used as a target.
	ForInStmt struct {
		For    token.Pos
		Left   Node
		In     token.Pos
		Right  Expr
		Rparen token.Pos
		Body   Stmt
	}

	// ForOfStmt is the `of` form, including `for await` (C.3).
	ForOfStmt struct {
		For      token.Pos
		AwaitPos token.Pos // NoPos unless `for await`
		Left     Node
		OfPos    token.Pos
		Right    Expr
		Rparen   token.Pos
		Body     Stmt
	}

	// BranchStmt is ContinueStatement and BreakStatement (C.4). Label is nil
	// for the bare form; the [no LineTerminator here] restriction is enforced
	// at parse time, so a label here was on the same line.
	BranchStmt struct {
		TokPos token.Pos
		TokEnd token.Pos
		Tok    token.Kind // CONTINUE or BREAK
		Label  *Ident
		Semi   token.Pos
	}

	// ReturnStmt is ReturnStatement (C.4).
	ReturnStmt struct {
		Return    token.Pos
		ReturnEnd token.Pos
		Result    Expr
		Semi      token.Pos
	}

	// LabeledStmt is LabelledStatement (C.4).
	LabeledStmt struct {
		Label *Ident
		Colon token.Pos
		Stmt  Stmt
	}

	// ThrowStmt is ThrowStatement (C.4).
	ThrowStmt struct {
		Throw token.Pos
		X     Expr
		Semi  token.Pos
	}

	// TryStmt is TryStatement (C.4). Exactly one of Catch and Finally may be
	// nil; both nil parses, per §6.3, and is rejected by name later.
	TryStmt struct {
		Try     token.Pos
		Body    *BlockStmt
		Catch   *CatchClause
		Finally *FinallyClause
	}

	// CatchClause is Catch (C.4). Param is nil for the bare `catch {}` form.
	// CatchType is the `: any` / `: unknown` annotation, which is the only
	// annotation a catch parameter admits.
	CatchClause struct {
		Catch     token.Pos
		Lparen    token.Pos
		Param     Expr // *Ident, *ObjectPattern, or *ArrayPattern
		CatchType TypeExpr
		Rparen    token.Pos
		Body      *BlockStmt
	}

	// FinallyClause is Finally (C.4).
	FinallyClause struct {
		Finally token.Pos
		Body    *BlockStmt
	}

	// SwitchStmt is SwitchStatement (C.3).
	SwitchStmt struct {
		Switch token.Pos
		Tag    Expr
		Lbrace token.Pos
		Cases  []*CaseClause
		Rbrace token.Pos
	}

	// CaseClause is CaseClause or DefaultClause (C.3). Cond is nil for
	// `default`. More than one default parses and is rejected later.
	CaseClause struct {
		Case  token.Pos
		Cond  Expr
		Colon token.Pos
		Body  []Stmt
	}

	// DebuggerStmt is DebuggerStatement (C).
	DebuggerStmt struct {
		Debugger    token.Pos
		DebuggerEnd token.Pos
		Semi        token.Pos
	}
)

// semiEnd returns the end of a statement whose terminator is ASI-eligible
// (§6.2). When the semicolon was inserted rather than written, the span stops
// at the last real token.
func semiEnd(semi token.Pos, last token.Pos) token.Pos {
	if semi != token.NoPos {
		return semi + 1
	}
	return last
}

func (s *BlockStmt) Pos() token.Pos   { return s.Lbrace }
func (s *BlockStmt) End() token.Pos   { return s.Rbrace + 1 }
func (s *EmptyStmt) Pos() token.Pos   { return s.Semi }
func (s *EmptyStmt) End() token.Pos   { return s.Semi + 1 }
func (s *ExprStmt) Pos() token.Pos    { return s.X.Pos() }
func (s *ExprStmt) End() token.Pos    { return semiEnd(s.Semi, s.X.End()) }
func (s *IfStmt) Pos() token.Pos      { return s.If }
func (s *IfStmt) End() token.Pos      { return endOf(s.Body.End(), s.Else) }
func (s *DoWhileStmt) Pos() token.Pos { return s.Do }
func (s *DoWhileStmt) End() token.Pos { return semiEnd(s.Semi, s.Rparen+1) }
func (s *WhileStmt) Pos() token.Pos   { return s.While }
func (s *WhileStmt) End() token.Pos   { return s.Body.End() }
func (s *ForStmt) Pos() token.Pos     { return s.For }
func (s *ForStmt) End() token.Pos     { return s.Body.End() }
func (s *ForInStmt) Pos() token.Pos   { return s.For }
func (s *ForInStmt) End() token.Pos   { return s.Body.End() }
func (s *ForOfStmt) Pos() token.Pos   { return s.For }
func (s *ForOfStmt) End() token.Pos   { return s.Body.End() }

func (s *BranchStmt) Pos() token.Pos { return s.TokPos }
func (s *BranchStmt) End() token.Pos {
	last := s.TokEnd
	if s.Label != nil {
		last = s.Label.End()
	}
	return semiEnd(s.Semi, last)
}
func (s *ReturnStmt) Pos() token.Pos { return s.Return }
func (s *ReturnStmt) End() token.Pos {
	last := s.ReturnEnd
	if !isNil(s.Result) {
		last = s.Result.End()
	}
	return semiEnd(s.Semi, last)
}
func (s *LabeledStmt) Pos() token.Pos    { return s.Label.Pos() }
func (s *LabeledStmt) End() token.Pos    { return s.Stmt.End() }
func (s *ThrowStmt) Pos() token.Pos      { return s.Throw }
func (s *ThrowStmt) End() token.Pos      { return semiEnd(s.Semi, s.X.End()) }
func (s *TryStmt) Pos() token.Pos        { return s.Try }
func (s *TryStmt) End() token.Pos        { return endOf(s.Body.End(), s.Catch, s.Finally) }
func (s *CatchClause) Pos() token.Pos    { return s.Catch }
func (s *CatchClause) End() token.Pos    { return s.Body.End() }
func (s *FinallyClause) Pos() token.Pos  { return s.Finally }
func (s *FinallyClause) End() token.Pos  { return s.Body.End() }
func (s *SwitchStmt) Pos() token.Pos     { return s.Switch }
func (s *SwitchStmt) End() token.Pos     { return s.Rbrace + 1 }
func (s *CaseClause) Pos() token.Pos     { return s.Case }
func (s *CaseClause) End() token.Pos {
	if n := len(s.Body); n > 0 {
		return s.Body[n-1].End()
	}
	return s.Colon + 1
}
func (s *DebuggerStmt) Pos() token.Pos { return s.Debugger }
func (s *DebuggerStmt) End() token.Pos { return semiEnd(s.Semi, s.DebuggerEnd) }

func (*BlockStmt) stmtNode()     {}
func (*EmptyStmt) stmtNode()     {}
func (*ExprStmt) stmtNode()      {}
func (*IfStmt) stmtNode()        {}
func (*DoWhileStmt) stmtNode()   {}
func (*WhileStmt) stmtNode()     {}
func (*ForStmt) stmtNode()       {}
func (*ForInStmt) stmtNode()     {}
func (*ForOfStmt) stmtNode()     {}
func (*BranchStmt) stmtNode()    {}
func (*ReturnStmt) stmtNode()    {}
func (*LabeledStmt) stmtNode()   {}
func (*ThrowStmt) stmtNode()     {}
func (*TryStmt) stmtNode()       {}
func (*SwitchStmt) stmtNode()    {}
func (*DebuggerStmt) stmtNode()  {}