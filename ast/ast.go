// Package ast defines the Vertex syntax tree. It is the parser's sole output
// and the analyzer's sole input.
//
// The tree records shape, never meaning. Several constructs are resolved by
// what an operand denotes rather than by syntactic form — an Index whose
// operand denotes a generic declaration is a TypeArgs list, `&` is address-of
// on a value and dereference on a typed_ptr, `~` is bitwise-NOT in an
// expression and underlying-type in a TypeSetTerm, and a constraint element
// that is a single identifier is both a one-term TypeSet and a constraint name.
// The parser cannot know which, and does not have to: each is a static rule
// checked over an already-parsed tree. This package is that tree.
//
// The corollary is that types are Exprs. A TypeName is an *Ident, one bracket
// node serves Index and TypeArgs alike, and a tuple type is a tuple literal.
// Conversion to a type representation happens in the analyzer.
//
// Two things deliberately leave no trace here:
//
//   - Terminator significance. Whether a line terminator ends a statement
//     depends on the innermost enclosing bracketing construct, which the parser
//     resolves as it goes; no node records it.
//   - Trailing commas, with one exception. TupleExpr keeps its own, because a
//     one-element tuple is distinguished from a parenthesized expression by
//     nothing else. Everywhere a trailing comma is optional and inert —
//     TypeArgs, Parameters, ArrayLit, LiteralValue — it is not recorded.
//
// Source syntax that is licensed by context does survive, because it is
// syntax: the `var` transfer marker is the entire difference between a move and
// a deep copy, so it is a node.
package ast

import "github.com/vertex-language/vertex/token"

type Node interface {
	Pos() token.Pos // first character of the node
	End() token.Pos // one past the last character
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Decl interface {
	Node
	declNode()
}

// ---------------------------------------------------------------- comments

// Comment is one line comment or one general comment.
type Comment struct {
	Slash token.Pos
	Text  string // including the opening delimiter, excluding any trailing line terminator
}

func (c *Comment) Pos() token.Pos { return c.Slash }
func (c *Comment) End() token.Pos { return c.Slash + token.Pos(len(c.Text)) }

// IsGeneral reports whether c is a `/* */` comment rather than a `//` one.
func (c *Comment) IsGeneral() bool {
	return len(c.Text) > 1 && c.Text[1] == '*'
}

// ActsAsTerminator reports whether c acts like a line terminator rather than
// like white space. A general comment containing no line terminator is white
// space; every other comment, line comments included, is a terminator.
func (c *Comment) ActsAsTerminator() bool {
	if !c.IsGeneral() {
		return true
	}
	for i := 0; i < len(c.Text); i++ {
		if c.Text[i] == '\n' || c.Text[i] == '\r' {
			return true
		}
	}
	return false
}

// CommentGroup is a run of comments with no other token and no blank line
// between them.
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() token.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() token.Pos { return g.List[len(g.List)-1].End() }

// ------------------------------------------------------------- identifiers

// Ident is an identifier token. That includes the blank identifier `_`, every
// contextual keyword used as a name, every predeclared type, tensor-element,
// and constraint name, and every reserved builtin name. None of those are
// distinguished lexically; the analyzer resolves them against the implicit
// outermost scope.
type Ident struct {
	NamePos token.Pos
	Name    string
}

func (x *Ident) Pos() token.Pos { return x.NamePos }
func (x *Ident) End() token.Pos { return x.NamePos + token.Pos(len(x.Name)) }
func (x *Ident) exprNode()      {}

// IsBlank reports whether x is the blank identifier, which introduces no
// binding. Which positions accept it is a static rule.
func (x *Ident) IsBlank() bool { return x != nil && x.Name == "_" }

func (x *Ident) String() string {
	if x == nil {
		return "<nil>"
	}
	return x.Name
}

// ----------------------------------------------------------- shared parts

// Param is one ParameterDecl.
//
// Name is nil where the parameter is written as a bare type, which is what a
// FunctionType's parameter list produces. That names must be either all present
// or all absent within one list is a static rule, so a mixed list parses.
//
// Ellipsis marks a variadic parameter. That it must be last and that there may
// be at most one are both static rules.
type Param struct {
	Doc      *CommentGroup
	Name     *Ident    // nil for a bare type
	Colon    token.Pos // NoPos when Name is nil
	Ellipsis token.Pos // position of `...`, else NoPos
	Type     Expr
}

func (p *Param) Pos() token.Pos {
	switch {
	case p.Name != nil:
		return p.Name.Pos()
	case p.Ellipsis.IsValid():
		return p.Ellipsis
	}
	return p.Type.Pos()
}
func (p *Param) End() token.Pos { return p.Type.End() }

// ParamList is Parameters.
type ParamList struct {
	Lparen token.Pos
	List   []*Param
	Rparen token.Pos
}

func (l *ParamList) Pos() token.Pos { return l.Lparen }
func (l *ParamList) End() token.Pos { return l.Rparen + 1 }

// TypeParam is one TypeParamDecl.
//
// Constraint is nil both for a bare name, which is constrained by `any`, and
// for a name in a group whose constraint appears on a later entry —
// `[A, B: Number]` constrains both. The parser does not distribute the trailing
// constraint backward; that is performed over an already-parsed list, so the
// written form survives for a formatter to reproduce.
type TypeParam struct {
	Name       *Ident
	Colon      token.Pos
	Constraint Expr // nil; a TypeSet
}

func (p *TypeParam) Pos() token.Pos { return p.Name.Pos() }
func (p *TypeParam) End() token.Pos {
	if p.Constraint != nil {
		return p.Constraint.End()
	}
	return p.Name.End()
}

// TypeParamList is TypeParameters.
type TypeParamList struct {
	Lbrack token.Pos
	List   []*TypeParam
	Rbrack token.Pos
}

func (l *TypeParamList) Pos() token.Pos { return l.Lbrack }
func (l *TypeParamList) End() token.Pos { return l.Rbrack + 1 }

// Marker is one FunctionMarker. `test` is a contextual keyword and scans as an
// identifier, so Kind alone cannot carry it; Name is authoritative.
type Marker struct {
	MarkerPos token.Pos
	Kind      token.Kind // ASYNC, GPU, NPU, or IDENT for `test`
	Name      string     // "async", "gpu", "npu", "test"
}

func (m *Marker) Pos() token.Pos { return m.MarkerPos }
func (m *Marker) End() token.Pos { return m.MarkerPos + token.Pos(len(m.Name)) }