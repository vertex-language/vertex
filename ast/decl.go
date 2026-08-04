package ast

import "github.com/vertex-language/vertex/token"

// Receiver is `( identifier : ReceiverType )`.
//
// Type may be wrapped in an *OwnershipType for the mut, var, and shared forms,
// and carries a TypeParameters list as an *IndexExpr for a method on a generic
// type. That list re-declares the receiver type's existing names rather than
// introducing fresh ones.
type Receiver struct {
	Lparen token.Pos
	Name   *Ident
	Colon  token.Pos
	Type   Expr
	Rparen token.Pos
}

func (r *Receiver) Pos() token.Pos { return r.Lparen }
func (r *Receiver) End() token.Pos { return r.Rparen + 1 }

// FuncDecl is both FunctionDecl and MethodDecl, and with them the initializer
// and deinitializer forms.
//
// Those get no shape of their own: `init` and `deinit` are contextual keywords
// that are ordinary method names in a receiver declaration, so they arrive as
// identifiers and land in Name like any other. Whether a given FuncDecl is one
// is a question about its Name and Recv, answered by the analyzer.
//
// TypeParams is parsed on a method too, where it is rejected, so that the
// diagnostic can point a caret at the bracket list rather than report a syntax
// error.
type FuncDecl struct {
	Doc        *CommentGroup
	Recv       *Receiver // nil for a free function
	Name       *Ident
	TypeParams *TypeParamList
	Type       *FuncType
	Body       *BlockStmt // nil only in error recovery
}

// Field is one FieldDecl. It serves a struct body, a class body, and a foreign
// class body, which are one production.
//
// A field list is newline-separated juxtaposition rather than a comma list, so
// the enclosing brace is terminator-significant and two fields on one line do
// not parse. Default is evaluated at construction for any omitted field.
type Field struct {
	Doc     *CommentGroup
	Name    *Ident
	Colon   token.Pos
	Type    Expr
	Assign  token.Pos
	Default Expr
	Comment *CommentGroup
}

func (f *Field) Pos() token.Pos { return f.Name.Pos() }
func (f *Field) End() token.Pos {
	if f.Default != nil {
		return f.Default.End()
	}
	return f.Type.End()
}

// RecordDecl is both StructDecl and ClassDecl.
//
// One node, because a class is byte-for-byte identical in layout to a struct
// and differs only in its member and method model. The shapes are the same; Kw
// carries the distinction and every consumer that cares reads it.
type RecordDecl struct {
	Doc        *CommentGroup
	KwPos      token.Pos
	Kw         token.Kind // STRUCT or CLASS
	Name       *Ident
	TypeParams *TypeParamList
	Lbrace     token.Pos
	Fields     []*Field
	Rbrace     token.Pos
}

// Variant is one enum variant: a unit variant, a payload variant `Name(T, U)`,
// or a unit variant with an explicit discriminant.
//
// Both suffixes are accepted on any variant, so an explicit discriminant on a
// payload variant parses and can be diagnosed as itself.
type Variant struct {
	Doc     *CommentGroup
	Name    *Ident
	Lparen  token.Pos
	Payload []Expr // types
	Rparen  token.Pos
	Assign  token.Pos
	Value   Expr
	Comment *CommentGroup
}

func (v *Variant) Pos() token.Pos { return v.Name.Pos() }
func (v *Variant) End() token.Pos {
	switch {
	case v.Value != nil:
		return v.Value.End()
	case v.Rparen.IsValid():
		return v.Rparen + 1
	}
	return v.Name.End()
}

// EnumDecl is an enum declaration. Its body is a comma-separated variant list
// and is not terminator-significant, which is what lets a variant list span
// lines.
type EnumDecl struct {
	Doc        *CommentGroup
	Enum       token.Pos
	Name       *Ident
	TypeParams *TypeParamList
	Colon      token.Pos
	Discrim    Expr // DiscriminantType; nil if absent
	Lbrace     token.Pos
	Variants   []*Variant
	Rbrace     token.Pos
}

// TypeAliasDecl is `type Name[params] = AliasTarget`. A Target of
// *AbstractType makes the alias nominal and opaque; anything else is
// transparent.
type TypeAliasDecl struct {
	Doc        *CommentGroup
	Type       token.Pos
	Name       *Ident
	TypeParams *TypeParamList
	Assign     token.Pos
	Target     Expr
}

// ConstraintElem is one element of a constraint body. Exactly one field is
// non-nil.
//
// Set holds a TypeSet and a constraint name undifferentiated, because a single
// identifier parses as both and is resolved by what the name denotes. A union
// is a *BinaryExpr with Op OR; a `~T` term is a *UnaryExpr with TILDE.
type ConstraintElem struct {
	Set    Expr
	Method *MethodReq
}

// MethodReq is a MethodRequirement. It takes a full Signature, so a constraint
// can require a marked method.
type MethodReq struct {
	Doc    *CommentGroup
	Func   token.Pos
	Name   *Ident
	Type   *FuncType
}

func (m *MethodReq) Pos() token.Pos { return m.Func }
func (m *MethodReq) End() token.Pos { return m.Type.End() }

// ConstraintDecl is a constraint declaration. There are no interfaces; a
// constraint is its own declaration form and is legal only in a bracket
// position, which is a static rule. Multiple elements form an intersection.
type ConstraintDecl struct {
	Doc        *CommentGroup
	Constraint token.Pos
	Name       *Ident
	Lbrace     token.Pos
	Elems      []*ConstraintElem
	Rbrace     token.Pos
}

// VarDecl is `let`/`var` with initializers, and bare `var Binding`.
//
// The bare form covers all three initializer-free spellings uniformly, `var w`
// among them: statement-leading `var` is always a declaration, so a bare
// transfer marker outside an owning position lands on a real declaration node
// and is diagnosed as itself rather than as a syntax error.
//
// This is also a TopLevelDecl, where the initializer must be
// compile-time-evaluable and the bare form is rejected. Both are static rules.
type VarDecl struct {
	Doc      *CommentGroup
	KwPos    token.Pos
	Kw       token.Kind // LET or VAR
	Bindings []*Binding
	Assign   token.Pos // NoPos for the bare form
	Values   []Expr    // owning positions: any may be a *TransferExpr
	Comment  *CommentGroup
}

// Binding is one entry of a BindingList.
type Binding struct {
	Name  *Ident
	Colon token.Pos
	Type  Expr // nil when inferred
}

func (b *Binding) Pos() token.Pos { return b.Name.Pos() }
func (b *Binding) End() token.Pos {
	if b.Type != nil {
		return b.Type.End()
	}
	return b.Name.End()
}

// ImportDecl is `import "path"` or `import ( ... )`. There is no aliasing form,
// no dot-import, and no blank import, so there is nothing to record but paths —
// the qualifier comes from the imported package's own package clause, and the
// path is a locator, not a name.
type ImportDecl struct {
	Doc    *CommentGroup
	Import token.Pos
	Lparen token.Pos // NoPos for the single-path form
	Paths  []*BasicLit
	Rparen token.Pos
}

// --------------------------------------------------------- declare blocks

// ForeignMember is a member of a declare body or a foreign class body.
//
// The declare body admits a foreign function, a foreign class, and a nested
// declare declaration; the foreign class body admits a foreign function, a
// foreign initializer, and a field. A nested declare and a field each parse and
// are rejected, so both *DeclareDecl and *Field implement this.
type ForeignMember interface {
	Node
	foreignMember()
}

// DeclareDecl is `declare framework "S" { }` or `declare module ["tag"] "S" { }`.
// Kind is the contextual keyword, which scans as an identifier.
//
// The variant tag is hoisted out of the module form, so a tagged framework
// block parses and is rejected with a message about `declare framework` rather
// than as a syntax error at the bracket. The tag set is closed; membership is a
// static rule.
type DeclareDecl struct {
	Doc     *CommentGroup
	Declare token.Pos
	KindPos token.Pos
	Kind    string // "framework" or "module"
	Variant *VariantTag
	Path    *BasicLit
	Lbrace  token.Pos
	Members []ForeignMember
	Rbrace  token.Pos
}

// VariantTag is the bracketed tag set.
type VariantTag struct {
	Lbrack token.Pos
	Tags   []*BasicLit
	Rbrack token.Pos
}

func (v *VariantTag) Pos() token.Pos { return v.Lbrack }
func (v *VariantTag) End() token.Pos { return v.Rbrack + 1 }

// ForeignFunc is both ForeignFuncDecl and ForeignInitDecl.
//
// Init marks the `init` prefix modifier, which is a modifier on func and not a
// function name; Name is nil for the unnamed initializer form that bare
// `Type(...)` construction resolves to.
//
// Body exists for the rejected form. A declare block describes call shapes
// only, so a block on a foreign declaration must parse in order to be diagnosed
// as itself. A marker on one is rejected the same way, and needs no field here
// because Type already carries every marker written.
type ForeignFunc struct {
	Doc    *CommentGroup
	Init   token.Pos // NoPos when absent
	Func   token.Pos
	Name   *Ident // nil for the unnamed initializer
	Type   *FuncType
	Body   *BlockStmt
}

// ForeignClass is a foreign class declaration. Its members are foreign
// functions, foreign initializers, and — parsed only to be rejected — fields.
type ForeignClass struct {
	Doc     *CommentGroup
	Class   token.Pos
	Name    *Ident
	Lbrace  token.Pos
	Members []ForeignMember
	Rbrace  token.Pos
}

// BadDecl marks an unparseable declaration span; see BadExpr.
type BadDecl struct {
	From, To token.Pos
}

// -------------------------------------------------------------- positions

func (d *RecordDecl) Pos() token.Pos     { return d.KwPos }
func (d *EnumDecl) Pos() token.Pos       { return d.Enum }
func (d *TypeAliasDecl) Pos() token.Pos  { return d.Type }
func (d *ConstraintDecl) Pos() token.Pos { return d.Constraint }
func (d *VarDecl) Pos() token.Pos        { return d.KwPos }
func (d *ImportDecl) Pos() token.Pos     { return d.Import }
func (d *DeclareDecl) Pos() token.Pos    { return d.Declare }
func (d *ForeignClass) Pos() token.Pos   { return d.Class }
func (d *BadDecl) Pos() token.Pos        { return d.From }

func (d *FuncDecl) Pos() token.Pos {
	if d.Recv != nil {
		return d.Recv.Pos()
	}
	return d.Type.Func
}

func (d *ForeignFunc) Pos() token.Pos {
	if d.Init.IsValid() {
		return d.Init
	}
	return d.Func
}

// Exactly one of Set and Method is non-nil, so the element's extent is
// whichever one is present.
func (e *ConstraintElem) Pos() token.Pos {
	if e.Set != nil {
		return e.Set.Pos()
	}
	return e.Method.Pos()
}

func (e *ConstraintElem) End() token.Pos {
	if e.Set != nil {
		return e.Set.End()
	}
	return e.Method.End()
}

func (d *RecordDecl) End() token.Pos     { return d.Rbrace + 1 }
func (d *EnumDecl) End() token.Pos       { return d.Rbrace + 1 }
func (d *TypeAliasDecl) End() token.Pos  { return d.Target.End() }
func (d *ConstraintDecl) End() token.Pos { return d.Rbrace + 1 }
func (d *DeclareDecl) End() token.Pos    { return d.Rbrace + 1 }
func (d *ForeignClass) End() token.Pos   { return d.Rbrace + 1 }
func (d *BadDecl) End() token.Pos        { return d.To }

func (d *FuncDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Type.End()
}

func (d *ForeignFunc) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Type.End()
}

func (d *VarDecl) End() token.Pos {
	if n := len(d.Values); n > 0 {
		return d.Values[n-1].End()
	}
	return d.Bindings[len(d.Bindings)-1].End()
}

func (d *ImportDecl) End() token.Pos {
	if d.Rparen.IsValid() {
		return d.Rparen + 1
	}
	return d.Paths[len(d.Paths)-1].End()
}

func (*FuncDecl) declNode()       {}
func (*RecordDecl) declNode()     {}
func (*EnumDecl) declNode()       {}
func (*TypeAliasDecl) declNode()  {}
func (*ConstraintDecl) declNode() {}
func (*VarDecl) declNode()        {}
func (*ImportDecl) declNode()     {}
func (*DeclareDecl) declNode()    {}
func (*BadDecl) declNode()        {}

func (*ForeignFunc) foreignMember()  {}
func (*ForeignClass) foreignMember() {}
func (*DeclareDecl) foreignMember()  {}
func (*Field) foreignMember()        {}