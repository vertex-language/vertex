package ast

import "github.com/vertex-language/vertex/token"

// BlockStmt is `{ StatementList }`.
//
// This brace is terminator-significant, unlike the braces of a LiteralValue,
// MapLit, field list, or declare body. That distinction is a parsing question
// resolved as the brace is opened, and the tree records only which construct
// the parser built.
type BlockStmt struct {
	Lbrace token.Pos
	List   []Stmt
	Rbrace token.Pos
}

// DeclStmt wraps a declaration appearing in statement position — in practice a
// *VarDecl, which is both a Statement and a TopLevelDecl.
type DeclStmt struct {
	Decl Decl
}

// ExprStmt is an expression in statement position. A bare CompositeLit or
// MapLit is excluded here, which is what keeps `{...}` unambiguous against a
// Block.
type ExprStmt struct {
	X Expr
}

// AssignStmt is both forms of assignment. Op is ASSIGN for the list form and a
// compound assign_op otherwise, in which case there is exactly one target and
// one value.
//
// A target is a bare PrimaryExpr, so a dereference-write `&p = 99` is a
// *UnaryExpr and the blank identifier is an ordinary *Ident; neither needs its
// own shape. Which shapes are assignable is a static rule. Assignment is a
// statement and never an expression, so no `=` appears inside a condition
// anywhere in this tree.
type AssignStmt struct {
	Targets []Expr
	OpPos   token.Pos
	Op      token.Kind
	Values  []Expr // owning positions: any may be a *TransferExpr
}

// IfStmt has no initializer clause; the two-statement error-checking idiom is
// intentional.
type IfStmt struct {
	If   token.Pos
	Cond Expr
	Body *BlockStmt
	Else Stmt // *BlockStmt or *IfStmt; nil if absent
}

// WhileStmt is the only loop primitive.
type WhileStmt struct {
	While token.Pos
	Cond  Expr
	Body  *BlockStmt
}

// ForStmt is `for IterationBinding in Expression Block`.
//
// Mode is INVALID for the bare form, MUT for exclusive access, or VAR for the
// consuming form. The marker sits on the binding rather than on the iterable,
// because what transfers is each element, one per iteration.
//
// Names holds one name, or two for the index/value and key/value forms. The
// marker and the two-name form do not combine, but both are parsed together
// here so the combination can be diagnosed as itself rather than as a syntax
// error at the comma.
type ForStmt struct {
	For     token.Pos
	ModePos token.Pos
	Mode    token.Kind // INVALID, MUT, VAR
	Names   []*Ident   // len 1 or 2
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

// CaseClause is `case PatternList :` or `default :`. Patterns == nil marks the
// default clause; that there is at most one is a static rule.
type CaseClause struct {
	Case     token.Pos // position of `case` or `default`
	Patterns []Expr    // nil for default
	Colon    token.Pos
	Body     []Stmt
}

// EnumPattern is `.identifier` or `.identifier(bindings)` in Pattern position,
// where a leading dot is always this and never an EnumShorthand reached through
// Expression.
//
// The distinction is not cosmetic: the payload entries are binding names rather
// than expressions, and they are views into the payload rather than copies.
type EnumPattern struct {
	Dot    token.Pos
	Name   *Ident
	Lparen token.Pos // NoPos when there is no payload list
	Binds  []*Ident
	Rparen token.Pos
}

// ReturnStmt is a bare comma list, never parenthesized: parentheses construct a
// tuple, bare commas unbuild one.
type ReturnStmt struct {
	Return  token.Pos
	Results []Expr
}

// DeferStmt takes a call and nothing else.
type DeferStmt struct {
	Defer token.Pos
	Call  *CallExpr
}

// BranchStmt is `break`, `continue`, or `fallthrough`. There are no loop
// labels, so it carries none.
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

// SelectClause is one clause of a select, covering all three ChannelCase forms.
//
// Kw is LET or VAR for the declaring form, which introduces bindings scoped to
// this clause's body; INVALID otherwise. Targets is non-nil for the assigning
// form over pre-declared targets. Exactly one of Bindings and Targets is
// non-empty, and both are empty for the bare form.
//
// Op is the channel operation, optionally wrapped in an *AwaitExpr. That one
// select is entirely bare or entirely awaited, and which calls are admissible
// here at all, are static rules. Op == nil marks the default clause.
type SelectClause struct {
	Case token.Pos

	KwPos    token.Pos
	Kw       token.Kind // LET, VAR, or INVALID
	Bindings []*Binding
	Targets  []Expr

	Assign token.Pos
	Op     Expr // *CallExpr or *AwaitExpr wrapping one; nil for default

	Colon token.Pos
	Body  []Stmt
}

// BadStmt marks an unparseable statement span; see BadExpr.
type BadStmt struct {
	From, To token.Pos
}

// -------------------------------------------------------------- positions

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
	return s.Return + token.Pos(len(token.RETURN.Spelling()))
}

func (s *BranchStmt) End() token.Pos {
	return s.TokPos + token.Pos(len(s.Tok.Spelling()))
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

// EnumPattern is an Expr because it occupies a Pattern slot, whose other
// alternative is an Expression.
func (*EnumPattern) exprNode() {}