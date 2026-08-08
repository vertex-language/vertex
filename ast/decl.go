package ast

import "github.com/vertex-language/vertex/token"

// Declarations — sections C.1, D, E, H, I, J.

// VarKind distinguishes the binding forms in C.1.
type VarKind uint8

const (
	VarVar VarKind = iota
	VarLet
	VarConst
	VarUsing
	VarAwaitUsing
)

type (
	// VarDecl is VariableStatement, LexicalDeclaration, UsingDeclaration, and
	// AwaitUsingDeclaration (C.1). AwaitPos is set only for VarAwaitUsing.
	VarDecl struct {
		AwaitPos token.Pos
		KindPos  token.Pos
		Kind     VarKind
		List     []*Binding // len >= 1
		Semi     token.Pos
	}

	// Binding is one LexicalBinding, VariableDeclaration, or UsingBinding
	// (C.1), and also one AmbientBinding (I).
	//
	// Name and Pattern are exclusive: exactly one is non-nil.
	Binding struct {
		Name     *Ident
		Pattern  Expr // *ObjectPattern or *ArrayPattern
		Definite token.Pos // DefiniteAssignmentAssertion `!`, NoPos if absent
		Type     TypeExpr
		Init     Expr
	}

	// FuncDecl is FunctionDeclaration, GeneratorDeclaration,
	// AsyncFunctionDeclaration, AsyncGeneratorDeclaration (D), and
	// AcceleratedFunctionDeclaration (D.4).
	//
	// One node for plain, kernel, and graph functions, matching D.4's single
	// production. They differ in ways that aren't grammatical.
	//
	// Async and Gen are recorded even though D.4 admits neither on an
	// accelerated function, so `kernel async function f() {}` can be rejected
	// by name rather than by parse failure (§5.3, §8.4).
	FuncDecl struct {
		Decorators []*Decorator
		Accel      AccelKind
		AccelPos   token.Pos
		FuncPos    token.Pos
		Name       *Ident // nil for the export-default and expression forms
		TypeParams *TypeParamList
		Params     *ParamList
		Result     TypeExpr   // may be *TypePredicate
		Body       *BlockStmt // nil ⇒ signature only, the `;` form
		Semi       token.Pos
		Async, Gen bool
	}

	// ClassDecl is ClassDeclaration and the shared shape behind ClassExpr (E).
	ClassDecl struct {
		Decorators  []*Decorator
		AbstractPos token.Pos
		ClassPos    token.Pos
		Name        *Ident // nil for anonymous
		TypeParams  *TypeParamList
		Extends     *HeritageClause
		Implements  *HeritageClause
		Lbrace      token.Pos
		Members     []Decl
		Rbrace      token.Pos
	}

	// StructDecl is StructDeclaration (E.1).
	//
	// Extends is present even though E.1 has no ClassExtendsClause: §6.3
	// requires `struct S extends B {}` to parse into a real StructDecl
	// carrying its heritage clause, so the eventual message can name the
	// construct instead of saying "unexpected token `extends`". A valid
	// program leaves it nil. See the note in the package README.
	//
	// Body is nil for the ambient form, `declare struct S;` (I). Members live
	// inside Body, so the nil case is unambiguous.
	StructDecl struct {
		Decorators []*Decorator
		StructPos  token.Pos // the contextual `struct` identifier
		Name       *Ident
		TypeParams *TypeParamList
		Extends    *HeritageClause // parse-first-reject-later; nil when valid
		Implements *HeritageClause
		Body       *StructBody // nil ⇒ ambient form
		Semi       token.Pos
	}

	// StructBody holds struct members in source order, which is layout order
	// (§5.3). Never sorted, canonicalized, or deduped (§5.4).
	StructBody struct {
		Lbrace  token.Pos
		Members []Decl
		Rbrace  token.Pos
	}

	// HeritageClause is ClassExtendsClause, ImplementsClause, or
	// InterfaceExtendsClause (E, H). Types has len >= 1 in a valid clause.
	HeritageClause struct {
		KeywordPos token.Pos
		Keyword    token.Kind // EXTENDS, or IDENT for the contextual `implements`
		Types      []TypeExpr
	}

	// FieldDefinition (E), also a StructElement (E.1).
	FieldDecl struct {
		Decorators  []*Decorator
		Mods        Modifiers
		AccessorPos token.Pos // `accessor`, NoPos if absent
		Name        Node      // *Ident, *BasicLit, *ComputedKey, or *PrivateIdent
		Optional    token.Pos // OptionalMarker `?`
		Definite    token.Pos // DefiniteAssignmentAssertion `!`
		Type        TypeExpr
		Init        Expr
		Semi        token.Pos
	}

	// MethodDecl is MethodDefinition (D.3) in a class, struct, or object
	// literal: plain, generator, async, async generator, get, and set.
	//
	// It implements both Decl and Expr: as a class or struct member it is a
	// Decl, and inside an object literal it lives in ObjectLit.Props, which is
	// typed []Expr (see the comment on that field in expr.go) — one node
	// shape covers MethodDefinition wherever the grammar allows it, matching
	// how a class/struct member and an object-literal method share this
	// exact production (D.3).
	MethodDecl struct {
		Decorators []*Decorator
		Mods       Modifiers
		AsyncPos   token.Pos
		StarPos    token.Pos // generator
		Accessor   token.Ctx // CtxGet, CtxSet, or CtxNone
		AccessorPos token.Pos
		Name       Node
		Optional   token.Pos // MethodOptionalMarker, only under [+Optional]
		TypeParams *TypeParamList
		Params     *ParamList
		Result     TypeExpr
		Body       *BlockStmt
	}

	// CtorDecl is ConstructorDeclaration (E). Body is nil for the `;` form.
	CtorDecl struct {
		Decorators []*Decorator
		Mods       Modifiers
		CtorPos    token.Pos
		Params     *ParamList
		Body       *BlockStmt
		Semi       token.Pos
	}

	// StaticBlockDecl is ClassStaticBlock (E).
	StaticBlockDecl struct {
		StaticPos token.Pos
		Body      *BlockStmt
	}

	// InterfaceDecl is InterfaceDeclaration (H).
	InterfaceDecl struct {
		IfacePos   token.Pos
		Name       *Ident
		TypeParams *TypeParamList
		Extends    *HeritageClause
		Body       *ObjectType
	}

	// TypeAliasDecl is TypeAliasDeclaration (H).
	TypeAliasDecl struct {
		TypePos    token.Pos
		Name       *Ident
		TypeParams *TypeParamList
		Assign     token.Pos
		Type       TypeExpr // may be *TypePredicate
		Semi       token.Pos
	}

	// EnumDecl is EnumDeclaration (H).
	EnumDecl struct {
		ConstPos   token.Pos // NoPos unless `const enum`
		EnumPos    token.Pos
		Name       *Ident
		Underlying TypeExpr // EnumUnderlyingType
		Lbrace     token.Pos
		Members    []*EnumMember
		Rbrace     token.Pos
	}

	// EnumMember is EnumMember (H). Value is an AssignmentExpression, not
	// folded to a constant — §1 forbids folding.
	EnumMember struct {
		Name  Node // *Ident or *BasicLit (StringLiteral)
		Assign token.Pos
		Value Expr
	}

	// NamespaceDecl is NamespaceDeclaration and AmbientNamespaceDeclaration
	// (H, I).
	NamespaceDecl struct {
		NsPos  token.Pos
		Name   Node // *Ident or *QualifiedName (IdentifierPath)
		Lbrace token.Pos
		Items  []Stmt
		Rbrace token.Pos
	}

	// ModuleDecl is AmbientModuleDeclaration and AmbientGlobalAugmentation
	// (I). Name is nil for `global`.
	ModuleDecl struct {
		KeywordPos token.Pos
		KeywordEnd token.Pos
		IsGlobal   bool
		Name       *BasicLit // StringLiteral
		Lbrace     token.Pos // NoPos for `module "x";`
		Items      []Stmt
		Rbrace     token.Pos
		Semi       token.Pos
	}

	// AmbientDecl is `declare X` (I). Inner is the declaration it wraps.
	//
	// The wrapper is a real node rather than a Modifiers bit because
	// AmbientDeclaration is its own production and its contents are a
	// restricted grammar — the `declare` in `declare class` is not the
	// ClassElementModifier `declare`.
	AmbientDecl struct {
		DeclarePos token.Pos
		Inner      Decl
	}
)

// AccelKind is the AcceleratedFunctionModifier (D.4).
type AccelKind uint8

const (
	AccelNone AccelKind = iota
	AccelKernel
	AccelGraph
)

// --- Parameters (D.1) -------------------------------------------------------

type (
	// ParamList is FormalParameters (D.1).
	ParamList struct {
		Lparen token.Pos
		This   *ThisParam // ThisParameter, always first when present
		List   []*Param
		Rest   *RestElem  // FunctionRestParameter, always last when present
		Rparen token.Pos
	}

	// ThisParam is `this: T` (D.1). Not a Param: it takes no decorators, no
	// modifiers, no initializer, and no optional marker.
	ThisParam struct {
		ThisPos token.Pos
		Type    TypeExpr
	}

	// Param is FormalParameter (D.1).
	//
	// Mods carries the parameter-property modifiers (AccessibilityModifier,
	// override, readonly). They are the same words as ClassElementModifier and
	// share the bitset; whether they are legal here depends on the enclosing
	// function being a constructor, which is not a parse question.
	Param struct {
		Decorators []*Decorator
		Mods       Modifiers
		Name       Expr // *Ident, *ObjectPattern, or *ArrayPattern
		Optional   token.Pos
		Type       TypeExpr
		Init       Expr
	}
)

// --- Decorators (F) ---------------------------------------------------------

// Decorator is `@ expr` (F). X is a *Ident, *QualifiedName, *ParenExpr, or
// *CallExpr, matching the three DecoratorMemberExpression forms.
type Decorator struct {
	At token.Pos
	X  Expr
}

// --- Imports and exports (J) ------------------------------------------------

type (
	// ImportDecl is ImportDeclaration (J.1).
	//
	// TypePos marks `import type`. Phase covers `import defer` and
	// `import source`. Specs is nil for a bare side-effect import.
	ImportDecl struct {
		ImportPos token.Pos
		TypePos   token.Pos
		Phase     ImportPhase
		PhasePos  token.Pos
		Default   *Ident        // ImportedDefaultBinding
		Namespace *Ident        // NameSpaceImport, `* as x`
		NamedPos  token.Pos     // `{`, NoPos if no NamedImports
		Named     []*ImportSpec
		NamedEnd  token.Pos
		FromPos   token.Pos
		Path      *BasicLit
		With      *WithClause
		Semi      token.Pos
	}

	// ImportSpec is ImportSpecifier (J.1). Name is nil for the plain form.
	ImportSpec struct {
		TypePos token.Pos
		Name    Node // ModuleExportName: *Ident or *BasicLit
		AsPos   token.Pos
		Local   *Ident
	}

	// ImportEqualsDecl is ImportEqualsDeclaration (J.1).
	ImportEqualsDecl struct {
		ImportPos  token.Pos
		Name       *Ident
		Assign     token.Pos
		Entity     Node      // *Ident or *QualifiedName
		RequirePos token.Pos // NoPos unless the require form
		Path       *BasicLit
		Rparen     token.Pos
		Semi       token.Pos
	}

	// ExportDecl is ExportDeclaration (J.2). Exactly one of Decl, Star,
	// Named, Default, and Assign describes the form.
	ExportDecl struct {
		ExportPos  token.Pos
		TypePos    token.Pos // `export type`
		DefaultPos token.Pos // `export default`
		StarPos    token.Pos // `export *`
		StarAs     Node      // ModuleExportName after `* as`
		NamedPos   token.Pos
		Named      []*ExportSpec
		NamedEnd   token.Pos
		Decl       Stmt      // the declaration or statement form
		Value      Expr      // export default <expr>
		AssignPos  token.Pos // `export = E`
		Entity     Node      // E in `export = E`
		NamespacePos token.Pos // `export as namespace X`
		Namespace  *Ident
		FromPos    token.Pos
		Path       *BasicLit
		With       *WithClause
		Semi       token.Pos
	}

	// ExportSpec is ExportSpecifier (J.2).
	ExportSpec struct {
		TypePos token.Pos
		Name    Node
		AsPos   token.Pos
		As      Node
	}

	// WithClause is WithClause (J.1) and the attribute list inside ImportType
	// (G.1).
	WithClause struct {
		WithPos token.Pos
		Lbrace  token.Pos
		List    []*ImportAttr
		Rbrace  token.Pos
	}

	// ImportAttr is ImportAttribute (J.1).
	ImportAttr struct {
		Key   Node // *Ident or *BasicLit
		Colon token.Pos
		Value *BasicLit
	}
)

// --- spans ------------------------------------------------------------------

func (d *VarDecl) Pos() token.Pos {
	if d.AwaitPos != token.NoPos {
		return d.AwaitPos
	}
	return d.KindPos
}
func (d *VarDecl) End() token.Pos {
	return semiEnd(d.Semi, d.List[len(d.List)-1].End())
}

func (d *Binding) Pos() token.Pos {
	if d.Name != nil {
		return d.Name.Pos()
	}
	return d.Pattern.Pos()
}
func (d *Binding) End() token.Pos {
	last := d.Pos()
	if d.Name != nil {
		last = d.Name.End()
	} else {
		last = d.Pattern.End()
	}
	if d.Definite != token.NoPos && d.Definite+1 > last {
		last = d.Definite + 1
	}
	return endOf(last, d.Type, d.Init)
}

func (d *FuncDecl) Pos() token.Pos {
	if n := len(d.Decorators); n > 0 {
		return d.Decorators[0].Pos()
	}
	if d.AccelPos != token.NoPos {
		return d.AccelPos
	}
	return d.FuncPos
}
func (d *FuncDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	last := endOf(d.Params.End(), d.Result)
	return semiEnd(d.Semi, last)
}

func (d *ClassDecl) Pos() token.Pos {
	if n := len(d.Decorators); n > 0 {
		return d.Decorators[0].Pos()
	}
	if d.AbstractPos != token.NoPos {
		return d.AbstractPos
	}
	return d.ClassPos
}
func (d *ClassDecl) End() token.Pos { return d.Rbrace + 1 }

func (d *StructDecl) Pos() token.Pos {
	if n := len(d.Decorators); n > 0 {
		return d.Decorators[0].Pos()
	}
	return d.StructPos
}
func (d *StructDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return semiEnd(d.Semi, d.Name.End())
}

func (d *StructBody) Pos() token.Pos     { return d.Lbrace }
func (d *StructBody) End() token.Pos     { return d.Rbrace + 1 }
func (d *HeritageClause) Pos() token.Pos { return d.KeywordPos }
func (d *HeritageClause) End() token.Pos { return d.Types[len(d.Types)-1].End() }

func (d *FieldDecl) Pos() token.Pos {
	if n := len(d.Decorators); n > 0 {
		return d.Decorators[0].Pos()
	}
	if p := d.Mods.Pos(); p != token.NoPos {
		return p
	}
	if d.AccessorPos != token.NoPos {
		return d.AccessorPos
	}
	return d.Name.Pos()
}
func (d *FieldDecl) End() token.Pos {
	last := d.Name.End()
	for _, p := range []token.Pos{d.Optional, d.Definite} {
		if p != token.NoPos && p+1 > last {
			last = p + 1
		}
	}
	return semiEnd(d.Semi, endOf(last, d.Type, d.Init))
}

func (d *MethodDecl) Pos() token.Pos {
	if n := len(d.Decorators); n > 0 {
		return d.Decorators[0].Pos()
	}
	if p := d.Mods.Pos(); p != token.NoPos {
		return p
	}
	for _, p := range []token.Pos{d.AsyncPos, d.AccessorPos, d.StarPos} {
		if p != token.NoPos {
			return p
		}
	}
	return d.Name.Pos()
}
func (d *MethodDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return endOf(d.Params.End(), d.Result)
}

func (d *CtorDecl) Pos() token.Pos {
	if n := len(d.Decorators); n > 0 {
		return d.Decorators[0].Pos()
	}
	if p := d.Mods.Pos(); p != token.NoPos {
		return p
	}
	return d.CtorPos
}
func (d *CtorDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return semiEnd(d.Semi, d.Params.End())
}

func (d *StaticBlockDecl) Pos() token.Pos { return d.StaticPos }
func (d *StaticBlockDecl) End() token.Pos { return d.Body.End() }
func (d *InterfaceDecl) Pos() token.Pos   { return d.IfacePos }
func (d *InterfaceDecl) End() token.Pos   { return d.Body.End() }
func (d *TypeAliasDecl) Pos() token.Pos   { return d.TypePos }
func (d *TypeAliasDecl) End() token.Pos   { return semiEnd(d.Semi, d.Type.End()) }

func (d *EnumDecl) Pos() token.Pos {
	if d.ConstPos != token.NoPos {
		return d.ConstPos
	}
	return d.EnumPos
}
func (d *EnumDecl) End() token.Pos      { return d.Rbrace + 1 }
func (d *EnumMember) Pos() token.Pos    { return d.Name.Pos() }
func (d *EnumMember) End() token.Pos    { return endOf(d.Name.End(), d.Value) }
func (d *NamespaceDecl) Pos() token.Pos { return d.NsPos }
func (d *NamespaceDecl) End() token.Pos { return d.Rbrace + 1 }
func (d *ModuleDecl) Pos() token.Pos    { return d.KeywordPos }
func (d *ModuleDecl) End() token.Pos {
	if d.Rbrace != token.NoPos {
		return d.Rbrace + 1
	}
	return semiEnd(d.Semi, d.Name.End())
}
func (d *AmbientDecl) Pos() token.Pos { return d.DeclarePos }
func (d *AmbientDecl) End() token.Pos { return d.Inner.End() }

func (p *ParamList) Pos() token.Pos { return p.Lparen }
func (p *ParamList) End() token.Pos { return p.Rparen + 1 }
func (p *ThisParam) Pos() token.Pos { return p.ThisPos }
func (p *ThisParam) End() token.Pos { return p.Type.End() }
func (p *Param) Pos() token.Pos {
	if n := len(p.Decorators); n > 0 {
		return p.Decorators[0].Pos()
	}
	if q := p.Mods.Pos(); q != token.NoPos {
		return q
	}
	return p.Name.Pos()
}
func (p *Param) End() token.Pos {
	last := p.Name.End()
	if p.Optional != token.NoPos && p.Optional+1 > last {
		last = p.Optional + 1
	}
	return endOf(last, p.Type, p.Init)
}

func (d *Decorator) Pos() token.Pos { return d.At }
func (d *Decorator) End() token.Pos { return d.X.End() }

func (d *ImportDecl) Pos() token.Pos { return d.ImportPos }
func (d *ImportDecl) End() token.Pos {
	last := d.ImportPos + 6
	if d.Path != nil {
		last = d.Path.End()
	}
	return semiEnd(d.Semi, endOf(last, d.With))
}
func (d *ImportSpec) Pos() token.Pos {
	if d.TypePos != token.NoPos {
		return d.TypePos
	}
	return spanOf(d.Local.Pos(), d.Name)
}
func (d *ImportSpec) End() token.Pos { return d.Local.End() }

func (d *ImportEqualsDecl) Pos() token.Pos { return d.ImportPos }
func (d *ImportEqualsDecl) End() token.Pos {
	if d.Rparen != token.NoPos {
		return semiEnd(d.Semi, d.Rparen+1)
	}
	return semiEnd(d.Semi, d.Entity.End())
}

func (d *ExportDecl) Pos() token.Pos { return d.ExportPos }
func (d *ExportDecl) End() token.Pos {
	last := d.ExportPos + 6
	for _, n := range []Node{d.Decl, d.Value, d.Entity, d.Namespace, d.Path, d.With} {
		if !isNil(n) && n.End() > last {
			last = n.End()
		}
	}
	if d.NamedEnd+1 > last {
		last = d.NamedEnd + 1
	}
	if !isNil(d.Decl) {
		return last // a declaration carries its own terminator
	}
	return semiEnd(d.Semi, last)
}
func (d *ExportSpec) Pos() token.Pos {
	if d.TypePos != token.NoPos {
		return d.TypePos
	}
	return d.Name.Pos()
}
func (d *ExportSpec) End() token.Pos { return endOf(d.Name.End(), d.As) }

func (d *WithClause) Pos() token.Pos { return d.WithPos }
func (d *WithClause) End() token.Pos { return d.Rbrace + 1 }
func (d *ImportAttr) Pos() token.Pos { return d.Key.Pos() }
func (d *ImportAttr) End() token.Pos { return d.Value.End() }

// Declarations are also statements; see the note in node.go.
func (*VarDecl) declNode()          {}
func (*VarDecl) stmtNode()          {}
func (*FuncDecl) declNode()         {}
func (*FuncDecl) stmtNode()         {}
func (*ClassDecl) declNode()        {}
func (*ClassDecl) stmtNode()        {}
func (*StructDecl) declNode()       {}
func (*StructDecl) stmtNode()       {}
func (*InterfaceDecl) declNode()    {}
func (*InterfaceDecl) stmtNode()    {}
func (*TypeAliasDecl) declNode()    {}
func (*TypeAliasDecl) stmtNode()    {}
func (*EnumDecl) declNode()         {}
func (*EnumDecl) stmtNode()         {}
func (*NamespaceDecl) declNode()    {}
func (*NamespaceDecl) stmtNode()    {}
func (*ModuleDecl) declNode()       {}
func (*ModuleDecl) stmtNode()       {}
func (*AmbientDecl) declNode()      {}
func (*AmbientDecl) stmtNode()      {}
func (*ImportDecl) declNode()       {}
func (*ImportDecl) stmtNode()       {}
func (*ImportEqualsDecl) declNode() {}
func (*ImportEqualsDecl) stmtNode() {}
func (*ExportDecl) declNode()       {}
func (*ExportDecl) stmtNode()       {}

func (*FieldDecl) declNode()       {}
func (*MethodDecl) declNode()      {}
func (*MethodDecl) exprNode()      {} // also usable as ObjectLit.Props[i] (D.3)
func (*CtorDecl) declNode()        {}
func (*StaticBlockDecl) declNode() {}