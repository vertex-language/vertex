// Package ast defines the syntax tree for Vertex.
//
// It records shape, never meaning (compiler_frontend.md §5). Nodes carry
// spans, not strings; text comes back from the token.File that produced them
// (§8.3). Nothing here resolves a name, decodes a literal, folds a constant,
// or knows what a type means.
//
// ast imports token and nothing else, and knows nothing about who built it
// (§2). Any import of scanner or parser from this package is a structural
// mistake.
package ast

import "github.com/vertex-language/vertex/token"

// Node is the root of every hierarchy.
//
// Every node computes its own span from its children, storing a Pos only where
// a leading keyword or contextual identifier makes it necessary (§5.4). End()
// is exact, not approximate — §1's runtime traps depend on it, so a node whose
// end is "roughly the closing brace" is a bug.
//
// Every node has a non-zero span (§1). A node whose children are all missing
// is a Bad* node, which stores its own bounds.
type Node interface {
	Pos() token.Pos
	End() token.Pos
}

// The four hierarchies. Section G is a separate sublanguage, not a subset of
// section B (§5.1): ConditionalType, MappedType, TupleType, TypePredicate,
// infer, TypeQuery, ImportType, TemplateLiteralType, and the keyof / unique /
// readonly / mutating operators have no expression counterpart at all.
// Merging them would mean a dozen node types illegal in every expression
// position and a later phase re-deriving "am I in a type?" at every step.
type (
	Decl     interface {
		Node
		declNode()
	}
	Stmt interface {
		Node
		stmtNode()
	}
	Expr interface {
		Node
		exprNode()
	}
	TypeExpr interface {
		Node
		typeNode()
	}
)

// Declaration nodes implement both Decl and Stmt.
//
// StatementListItem (C) is Statement | Declaration, and ModuleItem (J) adds
// imports and exports. Wrapping a Decl in a synthesized DeclStmt would violate
// §1's "the parser synthesizes nothing — no inserted nodes", so the marker
// methods overlap instead. Placement rules (an ImportDeclaration inside a
// block, an ExportDeclaration in a Script) are early errors checked by name
// later, per §6.3.

// ---------------------------------------------------------------------------
// Names

// Ident is an IdentifierName (A.2), including one that spells a contextual
// keyword. Ctx carries the contextual identity; whether a `struct` spelling
// declares anything is a parse question already answered by the time this node
// exists.
//
// There is no Obj field and there never will be. Go's go/parser resolved
// identifiers during the parse for fifteen years and now deprecates it,
// because the Ident→Object relation isn't computable without types (§1, §9).
type Ident struct {
	NamePos token.Pos
	NameEnd token.Pos
	Ctx     token.Ctx // CtxNone for an ordinary name
	Escaped bool      // spelling contained a UnicodeEscapeSequence
}

func (x *Ident) Pos() token.Pos { return x.NamePos }
func (x *Ident) End() token.Pos { return x.NameEnd }

// PrivateIdent is `#name` (A.2). Distinct from Ident because the positions it
// is legal in are disjoint.
type PrivateIdent struct {
	HashPos token.Pos
	NameEnd token.Pos
}

func (x *PrivateIdent) Pos() token.Pos { return x.HashPos }
func (x *PrivateIdent) End() token.Pos { return x.NameEnd }

// QualifiedName is a dotted name: EntityName (J.1), NamespaceName, TypeName,
// TypeQueryName (G.1), and IdentifierPath (H) all have this shape.
//
// One node for all of them because they differ only in what may follow, which
// the enclosing node already records.
type QualifiedName struct {
	X   Node   // *Ident or *QualifiedName
	Sel *Ident // after the dot
}

func (x *QualifiedName) Pos() token.Pos { return x.X.Pos() }
func (x *QualifiedName) End() token.Pos { return x.Sel.End() }

// ---------------------------------------------------------------------------
// Comments

// Comment is one MultiLineComment, SingleLineComment, or HashbangComment
// (A.1). Retained only under parser.ParseComments; otherwise the scanner's
// COMMENT tokens are dropped and only NLBefore survives (§4.3, §8.1).
type Comment struct {
	Slash token.Pos
	Tail  token.Pos
}

func (c *Comment) Pos() token.Pos { return c.Slash }
func (c *Comment) End() token.Pos { return c.Tail }

// CommentGroup is a run of comments with no blank line between them.
type CommentGroup struct {
	List []*Comment // len > 0
}

func (g *CommentGroup) Pos() token.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() token.Pos { return g.List[len(g.List)-1].End() }

// ---------------------------------------------------------------------------
// Bad nodes

// Bad* nodes hold slots (§5.4). A BadDecl inside a struct body occupies a
// field position so offsets don't silently shift, and a BadExpr keeps a binary
// operator's operand count right. Each stores its own bounds because it has no
// children to derive them from, and each spans at least one token so §1's
// non-zero span invariant holds.
type (
	BadExpr struct{ From, To token.Pos }
	BadStmt struct{ From, To token.Pos }
	BadDecl struct{ From, To token.Pos }
	BadType struct{ From, To token.Pos }
)

func (x *BadExpr) Pos() token.Pos { return x.From }
func (x *BadExpr) End() token.Pos { return x.To }
func (x *BadStmt) Pos() token.Pos { return x.From }
func (x *BadStmt) End() token.Pos { return x.To }
func (x *BadDecl) Pos() token.Pos { return x.From }
func (x *BadDecl) End() token.Pos { return x.To }
func (x *BadType) Pos() token.Pos { return x.From }
func (x *BadType) End() token.Pos { return x.To }

func (*BadExpr) exprNode() {}
func (*BadStmt) stmtNode() {}
func (*BadDecl) declNode() {}
func (*BadDecl) stmtNode() {}
func (*BadType) typeNode() {}
func (*Ident) exprNode()   {}
func (*PrivateIdent) exprNode() {}

// ---------------------------------------------------------------------------
// span helpers

// spanOf returns the span of the first non-nil node, used where a production
// has several optional leading parts (decorators, then modifiers, then a
// keyword).
func spanOf(fallback token.Pos, cands ...Node) token.Pos {
	for _, n := range cands {
		if !isNil(n) {
			return n.Pos()
		}
	}
	return fallback
}

func endOf(fallback token.Pos, cands ...Node) token.Pos {
	for i := len(cands) - 1; i >= 0; i-- {
		if !isNil(cands[i]) {
			return cands[i].End()
		}
	}
	return fallback
}