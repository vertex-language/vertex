// Package ast defines the Vertex syntax tree.
//
// The tree records shape, never meaning. Several Vertex constructs are resolved
// by what an operand denotes rather than by syntactic form — a[i] versus
// Stack[int32] (A.3.6), &x as address-of versus dereference (A.4.4), a lone
// identifier in a constraint body (A.7.2). The parser cannot know which, and
// A.0.2 says it does not have to: "Nothing else in Vertex is context-sensitive.
// Every other rejection is a static rule checked over an already-parsed tree."
// This package is that tree.
//
// The corollary is that types are Exprs. TypeName is an Ident, [...] is
// index/slice/instantiate/type-argument at once, and a tuple type is a tuple
// literal. Conversion to a type representation happens in the analyzer.
//
// Context parameters (Await, Npu, Own, Lit) leave no trace here — they are
// parser state. Source syntax licensed by them does survive: A.4.6's `var`
// marker is the entire difference between move and deep copy, so it is a node.
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

// Comment is one // or /* */ comment.
type Comment struct {
	Slash token.Pos
	Text  string // including the opening delimiter, excluding any trailing newline
}

func (c *Comment) Pos() token.Pos { return c.Slash }
func (c *Comment) End() token.Pos { return c.Slash + token.Pos(len(c.Text)) }

// IsLineTerminator reports whether this comment spans a LineTerminator, which
// per A.1.1 makes it one for statement-termination purposes (A.0.6).
func (c *Comment) IsLineTerminator() bool {
	for i := 0; i < len(c.Text); i++ {
		if c.Text[i] == '\n' || c.Text[i] == '\r' {
			return true
		}
	}
	return false
}

// CommentGroup is a run of comments with no other tokens and no blank lines
// between them.
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() token.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() token.Pos { return g.List[len(g.List)-1].End() }

// ------------------------------------------------------------------ ident

// Ident is an identifier, including a BlankIdentifier, a PredeclaredTypeName,
// a ReservedBuiltinName, and a ContextualKeyword used as a name. None of those
// are distinguished lexically (A.1.3, A.1.4); the analyzer resolves them
// against the implicit outermost scope.
type Ident struct {
	NamePos token.Pos
	Name    string
}

func (x *Ident) Pos() token.Pos { return x.NamePos }
func (x *Ident) End() token.Pos { return x.NamePos + token.Pos(len(x.Name)) }
func (x *Ident) exprNode()      {}

// IsBlank reports whether x is the BlankIdentifier `_` (A.1.2). It never
// introduces a usable binding.
func (x *Ident) IsBlank() bool { return x != nil && x.Name == "_" }

func (x *Ident) String() string {
	if x == nil {
		return "<nil>"
	}
	return x.Name
}

// ----------------------------------------------------------- shared parts

// Param is one entry in a ParamList.
//
// Name is nil inside a FunctionType (A.3.4: a function type names parameter
// types only). Ellipsis marks a variadic parameter; A.6.1 requires it last and
// permits at most one, both checked statically.
type Param struct {
	Doc      *CommentGroup
	Name     *Ident    // may be nil
	Colon    token.Pos // NoPos when Name is nil
	Ellipsis token.Pos // position of `...`, else NoPos
	Type     Expr
}

func (p *Param) Pos() token.Pos {
	if p.Name != nil {
		return p.Name.Pos()
	}
	if p.Ellipsis.IsValid() {
		return p.Ellipsis
	}
	return p.Type.Pos()
}
func (p *Param) End() token.Pos { return p.Type.End() }

type ParamList struct {
	Lparen token.Pos
	List   []*Param
	Rparen token.Pos
}

func (l *ParamList) Pos() token.Pos { return l.Lparen }
func (l *ParamList) End() token.Pos { return l.Rparen + 1 }

// TypeParam is one entry in a TypeParameterList (A.7.1).
//
// Constraint is nil both for a bare name (constraint `any`) and for a name in a
// group whose constraint appears on a later entry — A.7.1's "[A, B: Number]
// constrains both". The parser does not distribute; the analyzer does, because
// distributing would erase the written form vfmt needs.
type TypeParam struct {
	Name       *Ident // may be a BlankIdentifier
	Colon      token.Pos
	Constraint Expr // nil
}

func (p *TypeParam) Pos() token.Pos { return p.Name.Pos() }
func (p *TypeParam) End() token.Pos {
	if p.Constraint != nil {
		return p.Constraint.End()
	}
	return p.Name.End()
}

type TypeParamList struct {
	Lbrack token.Pos
	List   []*TypeParam
	Rbrack token.Pos
}

func (l *TypeParamList) Pos() token.Pos { return l.Lbrack }
func (l *TypeParamList) End() token.Pos { return l.Rbrack + 1 }

// Marker is a FunctionMarker (A.6.1). `test` is a ContextualKeyword and scans
// as IDENT, so Kind alone cannot carry it; Name is authoritative.
type Marker struct {
	MarkerPos token.Pos
	Kind      token.Kind // ASYNC, GPU, NPU, or IDENT for `test`
	Name      string     // "async", "gpu", "npu", "test"
}

func (m *Marker) Pos() token.Pos { return m.MarkerPos }
func (m *Marker) End() token.Pos { return m.MarkerPos + token.Pos(len(m.Name)) }