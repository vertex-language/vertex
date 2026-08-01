package ast

import "github.com/vertex-language/vertex/token"

// Receiver is `( Identifier : ReceiverType )` (A.6.1). Type may be wrapped in an
// *OwnershipType for the mut/var/shared forms, and carries a TypeParameterList
// as an *IndexExpr for a method on a generic type — A.7.6 makes that list bind
// the names rather than introduce fresh ones.
type Receiver struct {
	Lparen token.Pos
	Name   *Ident
	Colon  token.Pos
	Type   Expr
	Rparen token.Pos
}

func (r *Receiver) Pos() token.Pos { return r.Lparen }
func (r *Receiver) End() token.Pos { return r.Rparen + 1 }

// FuncDecl is a FunctionDeclaration and also an InitializerDeclaration and
// DeinitializerDeclaration (A.6.1, A.6.4).
//
// Those get no node of their own because they have no distinct shape: A.1.3
// makes `init` and `deinit` ContextualKeywords that are ordinary method names
// in a receiver declaration, so they scan as IDENT and parse as this. Whether a
// given FuncDecl is an initializer is a question about its Name and Recv,
// answered by the analyzer.
type FuncDecl struct {
	Doc        *CommentGroup
	Recv       *Receiver // nil for a free function
	Name       *Ident
	TypeParams *TypeParamList // nil; A.7.6 forbids one on a method
	Type       *FuncType
	Body       *BlockStmt // nil only in error recovery
}

// Field is one entry of a struct or class body (A.6.2).
type Field struct {
	Doc     *CommentGroup
	Name    *Ident
	Colon   token.Pos
	Type    Expr
	Assign  token.Pos
	Default Expr // A.6.2: evaluated at construction for any omitted field
	Comment *CommentGroup
}

func (f *Field) Pos() token.Pos { return f.Name.Pos() }
func (f *Field) End() token.Pos {
	if f.Default != nil {
		return f.Default.End()
	}
	return f.Type.End()
}

// RecordDecl is both StructDeclaration and ClassDeclaration (A.6.2, A.6.3).
//
// One node, because A.6.3 says a class is byte-for-byte identical in layout to
// a struct and differs only in its member and method model. The shapes are the
// same; Kw carries the distinction and every consumer that cares reads it.
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

// Variant is one entry of an enum body (A.6.5): a unit variant, a payload
// variant `Name(T, U)`, or a unit variant with an explicit discriminant.
//
// An explicit discriminant on a payload variant parses and is rejected (A.14).
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

// TypeAliasDecl is `type Name[params] = Target` (A.6.6). A Target of
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

// ConstraintElem is one element of a constraint body (A.7.2). Exactly one field
// is non-nil.
//
// Set holds a TypeSet or a ConstraintName undifferentiated, because A.7.2 says
// a single identifier parses as both and is resolved by what the name denotes.
// A union is a *BinaryExpr with Op OR; a `~T` term is a *UnaryExpr with TILDE.
type ConstraintElem struct {
	Set    Expr
	Method *MethodReq
}

// MethodReq is a MethodRequirement (A.7.2). Satisfied by any type declaring a
// matching receiver method; monomorphization lowers the call directly, so this
// introduces no interface value and no vtable.
type MethodReq struct {
	Doc    *CommentGroup
	Func   token.Pos
	Name   *Ident
	Params *ParamList
	Arrow  token.Pos
	Result Expr
}

func (m *MethodReq) Pos() token.Pos { return m.Func }
func (m *MethodReq) End() token.Pos {
	if m.Result != nil {
		return m.Result.End()
	}
	return m.Params.End()
}

type ConstraintDecl struct {
	Doc        *CommentGroup
	Constraint token.Pos
	Name       *Ident
	Lbrace     token.Pos
	Elems      []*ConstraintElem
	Rbrace     token.Pos
}

// VarDecl is `let`/`var` with initializers, and bare `var Binding` (A.5.1).
// Also a TopLevelDeclaration, where A.2 requires a compile-time-evaluable
// initializer.
type VarDecl struct {
	Doc      *CommentGroup
	KwPos    token.Pos
	Kw       token.Kind // LET or VAR
	Bindings []*Binding
	Assign   token.Pos // NoPos for bare `var x: T`
	Values   []Expr    // owning positions: any may be a *TransferExpr
	Comment  *CommentGroup
}

// Binding is one entry of a BindingList (A.5.1).
type Binding struct {
	Name  *Ident // may be a BlankIdentifier
	Colon token.Pos
	Type  Expr // nil when inferred; required for bare `var x: T`
}

func (b *Binding) Pos() token.Pos { return b.Name.Pos() }
func (b *Binding) End() token.Pos {
	if b.Type != nil {
		return b.Type.End()
	}
	return b.Name.End()
}

// ImportDecl is `import "path"` or `import ( ... )` (A.2.3). There is no
// aliasing form, no dot-import, and no blank import, so there is nothing to
// record but paths — the qualifier comes from the imported package's own
// PackageClause.
type ImportDecl struct {
	Doc    *CommentGroup
	Import token.Pos
	Lparen token.Pos // NoPos for the single-path form
	Paths  []*BasicLit
	Rparen token.Pos
}

// ---------------------------------------------------------- declare blocks

// DeclareDecl is `declare framework "S" { }` or `declare module ["tag"] "S" { }`
// (A.8.1). Kind is the ContextualKeyword, which scans as IDENT.
//
// A variant tag on a framework block parses and is rejected (A.8.2).
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

// VariantTag is the closed bracketed tag set of A.8.2. It reuses the same
// compile-time-configuration bracket as generic instantiation and array length.
type VariantTag struct {
	Lbrack token.Pos
	Tags   []*BasicLit
	Rbrack token.Pos
}

func (v *VariantTag) Pos() token.Pos { return v.Lbrack }
func (v *VariantTag) End() token.Pos { return v.Rbrack + 1 }

// ForeignMember is a member of a declare block.
type ForeignMember interface {
	Node
	foreignMember()
}

// ForeignFunc is a ForeignFunctionDeclaration or ForeignInitializerDeclaration
// (A.8.3).
//
// Init marks the `init` prefix modifier, which A.8.3 is explicit is a modifier
// on func and not a function name. Name is nil for the unnamed initializer form
// that bare `Type(...)` construction resolves to.
//
// Modifiers and Body exist for the error forms. A.0.5 makes rejected forms part
// of the grammar, so `private init func() -> Bad` and a foreign declaration
// with a body must parse in order to be diagnosed as themselves rather than as
// a syntax error.
type ForeignFunc struct {
	Doc       *CommentGroup
	Modifiers []*Ident   // A.8.3 ✗ visibility modifiers are banned
	Init      token.Pos  // NoPos when absent
	Func      token.Pos
	Name      *Ident // nil for the unnamed initializer
	Params    *ParamList
	Arrow     token.Pos
	Result    Expr
	Body      *BlockStmt // A.8.3 ✗ declarations cannot have bodies
}

// ForeignClass is a ForeignClassDeclaration (A.8.3). Fields are banned — a
// declare block describes call shape only — but parse, so Members may hold a
// *ForeignField for diagnosis.
type ForeignClass struct {
	Doc       *CommentGroup
	Modifiers []*Ident
	Class     token.Pos
	Name      *Ident
	Lbrace    token.Pos
	Members   []ForeignMember
	Rbrace    token.Pos
}

// ForeignField exists only to be rejected (A.8.3 ✗ fields describe layout).
type ForeignField struct {
	Doc   *CommentGroup
	Name  *Ident
	Colon token.Pos
	Type  Expr
}

func (d *ForeignField) Pos() token.Pos { return d.Name.Pos() }
func (d *ForeignField) End() token.Pos { return d.Type.End() }

type BadDecl struct {
	From, To token.Pos
}

// -------------------------------------------------------------- positions

func (d *FuncDecl) Pos() token.Pos {
	if d.Recv != nil {
		return d.Recv.Pos()
	}
	return d.Type.Func
}
func (d *RecordDecl) Pos() token.Pos     { return d.KwPos }
func (d *EnumDecl) Pos() token.Pos       { return d.Enum }
func (d *TypeAliasDecl) Pos() token.Pos  { return d.Type }
func (d *ConstraintDecl) Pos() token.Pos { return d.Constraint }
func (d *VarDecl) Pos() token.Pos        { return d.KwPos }
func (d *ImportDecl) Pos() token.Pos     { return d.Import }
func (d *DeclareDecl) Pos() token.Pos    { return d.Declare }
func (d *ForeignClass) Pos() token.Pos   { return d.Class }
func (d *BadDecl) Pos() token.Pos        { return d.From }

// A.7.2 guarantees exactly one of Set/Method is non-nil, so the element's
// extent is whichever one is present.
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

func (d *ForeignFunc) Pos() token.Pos {
	if len(d.Modifiers) > 0 {
		return d.Modifiers[0].Pos()
	}
	if d.Init.IsValid() {
		return d.Init
	}
	return d.Func
}

func (d *FuncDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Type.End()
}
func (d *RecordDecl) End() token.Pos     { return d.Rbrace + 1 }
func (d *EnumDecl) End() token.Pos       { return d.Rbrace + 1 }
func (d *TypeAliasDecl) End() token.Pos  { return d.Target.End() }
func (d *ConstraintDecl) End() token.Pos { return d.Rbrace + 1 }
func (d *DeclareDecl) End() token.Pos    { return d.Rbrace + 1 }
func (d *ForeignClass) End() token.Pos   { return d.Rbrace + 1 }
func (d *BadDecl) End() token.Pos        { return d.To }

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

func (d *ForeignFunc) End() token.Pos {
	switch {
	case d.Body != nil:
		return d.Body.End()
	case d.Result != nil:
		return d.Result.End()
	}
	return d.Params.End()
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
func (*ForeignField) foreignMember() {}