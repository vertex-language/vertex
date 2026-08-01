package ast

import "github.com/vertex-language/vertex/token"

// BlockStmt is `{ StatementList }` (A.5).
//
// A.0.6 makes this brace newline-significant, unlike the braces of a
// CompositeLit, MapLit, FieldList, or DeclareBody. The scanner cannot tell them
// apart, which is why it records line breaks as Token.NLBefore and never
// suppresses them.
type BlockStmt struct {
	Lbrace token.Pos
	List   []Stmt
	Rbrace token.Pos
}

// DeclStmt wraps a declaration appearing in statement position. In practice a
// *VarDecl, since A.5 lists VariableDeclaration among the statements and A.2
// lists it among the top-level declarations.
type DeclStmt struct {
	Decl Decl
}

// ExprStmt is A.5.9's ExpressionStatement. The grammar excludes a bare
// CompositeLiteral or MapLiteral here, which is what keeps `{...}` unambiguous
// against a BlockStmt.
type ExprStmt struct {
	X Expr
}

// AssignStmt is both forms of A.5.2. Op is ASSIGN for the list form and a
// CompoundAssignOperator otherwise; the compound form has exactly one target
// and one value.
//
// A target may be `&x` (an *UnaryExpr, dereference-on-the-write-side) or a
// BlankIdentifier. Assignment is a statement and never an expression, so there
// is no `=` inside a condition anywhere in this tree.
type AssignStmt struct {
	Targets []Expr
	OpPos   token.Pos
	Op      token.Kind
	Values  []Expr // owning positions: any may be a *TransferExpr
}

// IfStmt has no initializer clause. A.5.4 is explicit that the error-checking
// idiom is two statements and that the verbosity is intentional.
type IfStmt struct {
	If   token.Pos
	Cond Expr // parsed under [~Lit]
	Body *BlockStmt
	Else Stmt // *BlockStmt or *IfStmt; nil if absent
}

// WhileStmt is the only loop primitive (A.5.5).
type WhileStmt struct {
	While token.Pos
	Cond  Expr
	Body  *BlockStmt
}

// ForStmt is `for IterationBinding in Expr Block` (A.5.6).
//
// Mode is ILLEGAL for the bare (shared-access) form, MUT for exclusive access,
// or VAR for the consuming form. A.5.6 puts the marker on the binding rather
// than the iterable because what transfers is each element, one per iteration.
//
// Names holds one name, or two for the index/value and key/value forms.
type ForStmt struct {
	For     token.Pos
	ModePos token.Pos
	Mode    token.Kind // ILLEGAL, MUT, VAR
	Names   []*Ident   // len 1 or 2; entries may be BlankIdentifier
	In      token.Pos
	X       Expr
	Body    *BlockStmt
}

type SwitchStmt struct {
	Switch token.Pos
	Tag    Expr
	Lbrace token.Pos
	Cases  []*CaseClause
	Rbrace token.Pos
}

// CaseClause is `case PatternList :` or `default :` (A.5.7).
// Patterns == nil marks the default clause; A.5.7 permits at most one.
type CaseClause struct {
	Case     token.Pos // position of `case` or `default`
	Patterns []Expr    // nil for default
	Colon    token.Pos
	Body     []Stmt
}

// EnumPattern is `.Name` or `.Name(bindings)` in case position (A.5.7).
//
// Distinct from EnumShorthand because the payload entries are binding names,
// not expressions, and A.5.7 makes them views into the payload rather than
// copies.
type EnumPattern struct {
	Dot    token.Pos
	Name   *Ident
	Lparen token.Pos // NoPos when there is no payload list
	Binds  []*Ident  // entries may be BlankIdentifier
	Rparen token.Pos
}

type ReturnStmt struct {
	Return  token.Pos
	Results []Expr // bare comma list, never parenthesized (A.5.3)
}

// DeferStmt takes a call and nothing else (A.5.8). Arguments are evaluated at
// registration; only the call is postponed.
type DeferStmt struct {
	Defer token.Pos
	Call  *CallExpr
}

// BranchStmt is `break`, `continue`, or `fallthrough`. There are no loop labels
// (A.5.9), so it carries none.
type BranchStmt struct {
	TokPos token.Pos
	Tok    token.Kind // BREAK, CONTINUE, FALLTHROUGH
}

type SelectStmt struct {
	Select token.Pos
	Lbrace token.Pos
	Cases  []*SelectClause
	Rbrace token.Pos
}

// SelectClause is one case of A.10.2. Targets is non-nil for the
// `x = ch.receive()` form. Op is the channel operation, optionally wrapped in an
// *AwaitExpr — A.10.2 requires a single select to be entirely bare or entirely
// awaited, checked statically.
//
// Op == nil marks the default clause.
type SelectClause struct {
	Case    token.Pos
	Targets []Expr
	Assign  token.Pos
	Op      Expr // *CallExpr or *AwaitExpr wrapping one; nil for default
	Colon   token.Pos
	Body    []Stmt
}

type BadStmt struct {
	From, To token.Pos
}

func (s *BlockStmt) Pos() token.Pos    { return s.Lbrace }
func (s *DeclStmt) Pos() token.Pos     { return s.Decl.Pos() }
func (s *ExprStmt) Pos() token.Pos     { return s.X.Pos() }
func (s *AssignStmt) Pos() token.Pos   { return s.Targets[0].Pos() }
func (s *IfStmt) Pos() token.Pos       { return s.If }
func (s *WhileStmt) Pos() token.Pos    { return s.While }
func (s *ForStmt) Pos() token.Pos      { return s.For }
func (s *SwitchStmt) Pos() token.Pos   { return s.Switch }
func (s *CaseClause) Pos() token.Pos   { return s.Case }
func (s *EnumPattern) Pos() token.Pos  { return s.Dot }
func (s *ReturnStmt) Pos() token.Pos   { return s.Return }
func (s *DeferStmt) Pos() token.Pos    { return s.Defer }
func (s *BranchStmt) Pos() token.Pos   { return s.TokPos }
func (s *SelectStmt) Pos() token.Pos   { return s.Select }
func (s *SelectClause) Pos() token.Pos { return s.Case }
func (s *BadStmt) Pos() token.Pos      { return s.From }

func (s *BlockStmt) End() token.Pos  { return s.Rbrace + 1 }
func (s *DeclStmt) End() token.Pos   { return s.Decl.End() }
func (s *ExprStmt) End() token.Pos   { return s.X.End() }
func (s *AssignStmt) End() token.Pos { return s.Values[len(s.Values)-1].End() }
func (s *WhileStmt) End() token.Pos  { return s.Body.End() }
func (s *ForStmt) End() token.Pos    { return s.Body.End() }
func (s *SwitchStmt) End() token.Pos { return s.Rbrace + 1 }
func (s *DeferStmt) End() token.Pos  { return s.Call.End() }
func (s *SelectStmt) End() token.Pos { return s.Rbrace + 1 }
func (s *BadStmt) End() token.Pos    { return s.To }

func (s *IfStmt) End() token.Pos {
	if s.Else != nil {
		return s.Else.End()
	}
	return s.Body.End()
}

func (s *CaseClause) End() token.Pos {
	if n := len(s.Body); n > 0 {
		return s.Body[n-1].End()
	}
	return s.Colon + 1
}

func (s *SelectClause) End() token.Pos {
	if n := len(s.Body); n > 0 {
		return s.Body[n-1].End()
	}
	return s.Colon + 1
}

func (s *EnumPattern) End() token.Pos {
	if s.Rparen.IsValid() {
		return s.Rparen + 1
	}
	return s.Name.End()
}

func (s *ReturnStmt) End() token.Pos {
	if n := len(s.Results); n > 0 {
		return s.Results[n-1].End()
	}
	return s.Return + 6
}

func (s *BranchStmt) End() token.Pos {
	return s.TokPos + token.Pos(len(s.Tok.String()))
}

func (*BlockStmt) stmtNode()  {}
func (*DeclStmt) stmtNode()   {}
func (*ExprStmt) stmtNode()   {}
func (*AssignStmt) stmtNode() {}
func (*IfStmt) stmtNode()     {}
func (*WhileStmt) stmtNode()  {}
func (*ForStmt) stmtNode()    {}
func (*SwitchStmt) stmtNode() {}
func (*ReturnStmt) stmtNode() {}
func (*DeferStmt) stmtNode()  {}
func (*BranchStmt) stmtNode() {}
func (*SelectStmt) stmtNode() {}
func (*BadStmt) stmtNode()    {}

func (*EnumPattern) exprNode() {} // appears in Pattern position (A.5.7)