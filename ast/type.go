package ast

import "github.com/vertex-language/vertex/token"

// Types — vertex_grammar.md section G.
//
// A separate hierarchy, entered at TypeAnnotation, ReturnTypeAnnotation,
// TypeArguments, heritage clauses, and the `< Type >` assertion (§5.1). The
// parser always knows which side it is on.
//
// The genuinely ambiguous cases are few, and each gets one node covering both
// readings, resolved later:
//
//	Foo<N>      type or const-generic value?   → TypeRef,  on what N binds to
//	Foo<3>      literal type or const arg?     → LiteralType, against the params
//	mutating T  passing mode or type operator? → TypeOp{Mutating}, on position

type (
	// TypeRef is TypeReference (G.1). Name is an *Ident or *QualifiedName.
	//
	// This node also covers the const-generic reading of `Foo<N>`: whether N
	// is a type or a value follows from what it binds to, which is not a parse
	// question.
	TypeRef struct {
		Name Node
		Args *TypeArgList
	}

	// PredefinedType is a PredefinedType name (G.1). Recognizing it is a
	// lexical table lookup, not name resolution — `string` and `number` are
	// contextual keywords carrying a token.Ctx, while `void` and `null` are
	// reserved and carry a token.Kind, so both fields exist and one is set.
	PredefinedType struct {
		NamePos token.Pos
		NameEnd token.Pos
		Ctx     token.Ctx
		Kind    token.Kind
	}

	// LiteralType is LiteralType (G.1), including the negated numeric forms.
	//
	// Also the node for a const-generic argument: `Foo<3>` is this, and which
	// reading applies is settled against the parameter list later.
	LiteralType struct {
		MinusPos token.Pos // NoPos unless `- NumericLiteral`
		Value    *BasicLit
	}

	// ThisType is `this` in type position (G.1).
	ThisType struct {
		ThisPos token.Pos
		ThisEnd token.Pos
	}

	// ParenType is ParenthesizedType (G.1). Retained for the same reason
	// ParenExpr is (§5.4): `(A | B)[]` and `A | B[]` are different types.
	ParenType struct {
		Lparen token.Pos
		X      TypeExpr
		Rparen token.Pos
	}

	// ArrayType is `T[]` (G.1 PostfixType). Distinct from IndexedAccessType,
	// which has an index.
	ArrayType struct {
		Elem   TypeExpr
		Lbrack token.Pos
		Rbrack token.Pos
	}

	// IndexedAccessType is `T[K]` (G.1 PostfixType).
	IndexedAccessType struct {
		X      TypeExpr
		Lbrack token.Pos
		Index  TypeExpr
		Rbrack token.Pos
	}

	// UnionType is UnionType (G.1). LeadPos is the optional leading `|`.
	UnionType struct {
		LeadPos token.Pos
		Types   []TypeExpr
	}

	// IntersectionType is IntersectionType (G.1).
	IntersectionType struct {
		LeadPos token.Pos
		Types   []TypeExpr
	}

	// TypeOp is keyof, unique, readonly, or mutating (G.1 TypeOperator).
	//
	// `mutating T` is ambiguous between a passing mode and a type operator and
	// gets this node either way; position decides later.
	TypeOp struct {
		OpPos token.Pos
		Op    token.Ctx // CtxKeyof, CtxUnique, CtxReadonly, CtxMutating
		X     TypeExpr
	}

	// InferType is `infer T` with an optional constraint (G.1).
	InferType struct {
		InferPos   token.Pos
		Name       *Ident
		ExtendsPos token.Pos
		Constraint TypeExpr
	}

	// CondType is ConditionalType (G.1).
	CondType struct {
		Check      TypeExpr
		ExtendsPos token.Pos
		Extends    TypeExpr
		Quest      token.Pos
		Then       TypeExpr
		Colon      token.Pos
		Else       TypeExpr
	}

	// FuncType is FunctionType (G.4). Result may be a *TypePredicate.
	FuncType struct {
		TypeParams *TypeParamList
		Params     *ParamTypeList
		Arrow      token.Pos
		Result     TypeExpr
	}

	// CtorType is ConstructorType (G.4).
	CtorType struct {
		AbstractPos token.Pos
		NewPos      token.Pos
		TypeParams  *TypeParamList
		Params      *ParamTypeList
		Arrow       token.Pos
		Result      TypeExpr
	}

	// ObjectType is ObjectType (G.3), also the body of an InterfaceDecl.
	ObjectType struct {
		Lbrace  token.Pos
		Members []TypeExpr // the seven TypeMember forms
		Rbrace  token.Pos
	}

	// MappedType is MappedType (G.6).
	MappedType struct {
		Lbrace       token.Pos
		ReadonlyPos  token.Pos
		ReadonlySign token.Kind // ADD, SUB, or INVALID for a bare `readonly`
		Lbrack       token.Pos
		Name         *Ident
		InPos        token.Pos
		Constraint   TypeExpr
		AsPos        token.Pos
		As           TypeExpr
		Rbrack       token.Pos
		OptionalPos  token.Pos
		OptionalSign token.Kind
		Type         TypeExpr
		Rbrace       token.Pos
	}

	// TupleType is TupleType (G.5).
	TupleType struct {
		Lbrack token.Pos
		Elems  []TypeExpr // *TupleElem
		Rbrack token.Pos
	}

	// TupleElem is TupleElementType (G.5), covering the named, optional, and
	// rest forms in one node.
	TupleElem struct {
		Ellipsis    token.Pos
		Name        *Ident
		NameColon   token.Pos
		OptionalPos token.Pos
		Type        TypeExpr
	}

	// TypeQuery is `typeof x` and `typeof import(...)` in type position (G.1).
	TypeQuery struct {
		TypeofPos token.Pos
		Name      Node // *Ident, *QualifiedName, or *ImportType
		Args      *TypeArgList
	}

	// ImportType is ImportType (G.1).
	ImportType struct {
		ImportPos token.Pos
		Lparen    token.Pos
		Path      *BasicLit
		With      *WithClause
		Rparen    token.Pos
		Qualifier Node // *Ident or *QualifiedName after the dot
		Args      *TypeArgList
	}

	// TemplateLiteralType is TemplateLiteralType (G.7).
	//
	// The tokens are exactly the ones in A.6 — the grammar has no
	// TemplateTypeHead precisely so there is one way to tokenize `` `hello ${ ``
	// (§4.5) — so this node mirrors TemplateLit with TypeExpr substitutions.
	TemplateLiteralType struct {
		Quasis []*TemplateElem
		Types  []TypeExpr // len == len(Quasis)-1
	}

	// TypePredicate is TypePredicate (G.4).
	//
	// It is a TypeExpr even though G.1 keeps `Type` and `TypePredicate`
	// distinct, because ReturnTypeAnnotation admits either and FuncDecl.Result
	// has to hold both. Positions where only a Type is legal reject it later.
	TypePredicate struct {
		AssertsPos token.Pos
		Param      Node // *Ident or *ThisType
		IsPos      token.Pos
		Type       TypeExpr // nil for the bare `asserts x` form
	}
)

// --- Type members (G.3) -----------------------------------------------------
//
// These are TypeExprs so they can sit in ObjectType.Members, and grammar E
// references MethodSignature, GetAccessorSignature, SetAccessorSignature, and
// IndexSignature directly as ClassElements — so a class body holds them too,
// wrapped by the parser in the modifiers it saw. One set of nodes, three
// homes: object types, interface bodies, and ambient class bodies.

type (
	// PropertySig is PropertySignature (G.3).
	PropertySig struct {
		ReadonlyPos token.Pos
		Name        Node
		OptionalPos token.Pos
		Type        TypeExpr
		Sep         token.Pos // TypeMemberSeparator, ASI-eligible
	}

	// MethodSig is MethodSignature (G.3).
	MethodSig struct {
		Name        Node
		OptionalPos token.Pos
		TypeParams  *TypeParamList
		Params      *ParamTypeList
		Result      TypeExpr
		Sep         token.Pos
	}

	// CallSig is CallSignature (G.3).
	CallSig struct {
		TypeParams *TypeParamList
		Params     *ParamTypeList
		Result     TypeExpr
		Sep        token.Pos
	}

	// ConstructSig is ConstructSignature (G.3).
	ConstructSig struct {
		AbstractPos token.Pos
		NewPos      token.Pos
		TypeParams  *TypeParamList
		Params      *ParamTypeList
		Result      TypeExpr
		Sep         token.Pos
	}

	// IndexSig is IndexSignature (G.3).
	IndexSig struct {
		ReadonlyPos token.Pos
		Lbrack      token.Pos
		Name        *Ident
		Key         TypeExpr
		Rbrack      token.Pos
		Type        TypeExpr
		Sep         token.Pos
	}

	// AccessorSig is GetAccessorSignature and SetAccessorSignature (G.3).
	AccessorSig struct {
		KwPos  token.Pos
		Kw     token.Ctx // CtxGet or CtxSet
		Name   Node
		Lparen token.Pos
		Param  *ParamType // set form only
		Rparen token.Pos
		Result TypeExpr // get form only
		Sep    token.Pos
	}
)

// --- Parameter and argument lists (G.2, G.4) --------------------------------

type (
	// ParamTypeList is ParameterTypeList (G.4). Distinct from ParamList (D.1):
	// a type-position parameter takes no decorators, no modifiers, and no
	// initializer, and conflating them would mean fields that are always nil
	// on one side.
	ParamTypeList struct {
		Lparen token.Pos
		This   *ThisParam
		List   []*ParamType
		Rest   *RestElem
		Rparen token.Pos
	}

	// ParamType is ParameterType (G.4).
	ParamType struct {
		Name        Expr // *Ident, *ObjectPattern, or *ArrayPattern
		OptionalPos token.Pos
		Type        TypeExpr
	}

	// TypeParamList is TypeParameters (G.2).
	TypeParamList struct {
		Langle token.Pos
		List   []*TypeParam
		Rangle token.Pos
	}

	// TypeParam is TypeParameter (G.2).
	//
	// Constraint and Type are exclusive: G.2 admits either `T extends U` or
	// `T: U`, and both spellings reach the same later phase.
	TypeParam struct {
		Mods       []TypeParamMod
		Name       *Ident
		ExtendsPos token.Pos
		Constraint TypeExpr
		ColonPos   token.Pos
		Type       TypeExpr
		Assign     token.Pos
		Default    TypeExpr
	}

	// TypeParamMod is TypeParameterModifier (G.2): const, in, out.
	TypeParamMod struct {
		Pos token.Pos
		End token.Pos
		Tok token.Kind // CONST, IN, or IDENT for the contextual `out`
	}

	// TypeArgList is TypeArguments (G.2).
	//
	// The closing `>` is one token by construction: the scanner never merges
	// `>` (§4.2), so the type parser does nothing special here.
	TypeArgList struct {
		Langle token.Pos
		List   []TypeExpr
		Rangle token.Pos
	}
)

// --- spans ------------------------------------------------------------------

func (t *TypeRef) Pos() token.Pos           { return t.Name.Pos() }
func (t *TypeRef) End() token.Pos           { return endOf(t.Name.End(), t.Args) }
func (t *PredefinedType) Pos() token.Pos    { return t.NamePos }
func (t *PredefinedType) End() token.Pos    { return t.NameEnd }
func (t *LiteralType) Pos() token.Pos       { return spanOf(t.Value.Pos(), nil) }
func (t *LiteralType) End() token.Pos       { return t.Value.End() }
func (t *ThisType) Pos() token.Pos          { return t.ThisPos }
func (t *ThisType) End() token.Pos          { return t.ThisEnd }
func (t *ParenType) Pos() token.Pos         { return t.Lparen }
func (t *ParenType) End() token.Pos         { return t.Rparen + 1 }
func (t *ArrayType) Pos() token.Pos         { return t.Elem.Pos() }
func (t *ArrayType) End() token.Pos         { return t.Rbrack + 1 }
func (t *IndexedAccessType) Pos() token.Pos { return t.X.Pos() }
func (t *IndexedAccessType) End() token.Pos { return t.Rbrack + 1 }

func (t *UnionType) Pos() token.Pos {
	if t.LeadPos != token.NoPos {
		return t.LeadPos
	}
	return t.Types[0].Pos()
}
func (t *UnionType) End() token.Pos { return t.Types[len(t.Types)-1].End() }
func (t *IntersectionType) Pos() token.Pos {
	if t.LeadPos != token.NoPos {
		return t.LeadPos
	}
	return t.Types[0].Pos()
}
func (t *IntersectionType) End() token.Pos { return t.Types[len(t.Types)-1].End() }

func (t *TypeOp) Pos() token.Pos     { return t.OpPos }
func (t *TypeOp) End() token.Pos     { return t.X.End() }
func (t *InferType) Pos() token.Pos  { return t.InferPos }
func (t *InferType) End() token.Pos  { return endOf(t.Name.End(), t.Constraint) }
func (t *CondType) Pos() token.Pos   { return t.Check.Pos() }
func (t *CondType) End() token.Pos   { return t.Else.End() }
func (t *FuncType) Pos() token.Pos   { return spanOf(t.Params.Pos(), t.TypeParams) }
func (t *FuncType) End() token.Pos   { return t.Result.End() }
func (t *CtorType) Pos() token.Pos {
	if t.AbstractPos != token.NoPos {
		return t.AbstractPos
	}
	return t.NewPos
}
func (t *CtorType) End() token.Pos            { return t.Result.End() }
func (t *ObjectType) Pos() token.Pos          { return t.Lbrace }
func (t *ObjectType) End() token.Pos          { return t.Rbrace + 1 }
func (t *MappedType) Pos() token.Pos          { return t.Lbrace }
func (t *MappedType) End() token.Pos          { return t.Rbrace + 1 }
func (t *TupleType) Pos() token.Pos           { return t.Lbrack }
func (t *TupleType) End() token.Pos           { return t.Rbrack + 1 }
func (t *TupleElem) Pos() token.Pos {
	if t.Ellipsis != token.NoPos {
		return t.Ellipsis
	}
	return spanOf(t.Type.Pos(), t.Name)
}
func (t *TupleElem) End() token.Pos { return t.Type.End() }

func (t *TypeQuery) Pos() token.Pos  { return t.TypeofPos }
func (t *TypeQuery) End() token.Pos  { return endOf(t.Name.End(), t.Args) }
func (t *ImportType) Pos() token.Pos { return t.ImportPos }
func (t *ImportType) End() token.Pos {
	return endOf(t.Rparen+1, t.Qualifier, t.Args)
}
func (t *TemplateLiteralType) Pos() token.Pos { return t.Quasis[0].Pos() }
func (t *TemplateLiteralType) End() token.Pos { return t.Quasis[len(t.Quasis)-1].End() }

func (t *TypePredicate) Pos() token.Pos {
	if t.AssertsPos != token.NoPos {
		return t.AssertsPos
	}
	return t.Param.Pos()
}
func (t *TypePredicate) End() token.Pos { return endOf(t.Param.End(), t.Type) }

func sigEnd(sep token.Pos, last token.Pos) token.Pos {
	if sep != token.NoPos {
		return sep + 1
	}
	return last
}

func (t *PropertySig) Pos() token.Pos {
	if t.ReadonlyPos != token.NoPos {
		return t.ReadonlyPos
	}
	return t.Name.Pos()
}
func (t *PropertySig) End() token.Pos {
	last := t.Name.End()
	if t.OptionalPos != token.NoPos && t.OptionalPos+1 > last {
		last = t.OptionalPos + 1
	}
	return sigEnd(t.Sep, endOf(last, t.Type))
}
func (t *MethodSig) Pos() token.Pos { return t.Name.Pos() }
func (t *MethodSig) End() token.Pos {
	return sigEnd(t.Sep, endOf(t.Params.End(), t.Result))
}
func (t *CallSig) Pos() token.Pos { return spanOf(t.Params.Pos(), t.TypeParams) }
func (t *CallSig) End() token.Pos { return sigEnd(t.Sep, endOf(t.Params.End(), t.Result)) }
func (t *ConstructSig) Pos() token.Pos {
	if t.AbstractPos != token.NoPos {
		return t.AbstractPos
	}
	return t.NewPos
}
func (t *ConstructSig) End() token.Pos { return sigEnd(t.Sep, endOf(t.Params.End(), t.Result)) }
func (t *IndexSig) Pos() token.Pos {
	if t.ReadonlyPos != token.NoPos {
		return t.ReadonlyPos
	}
	return t.Lbrack
}
func (t *IndexSig) End() token.Pos    { return sigEnd(t.Sep, t.Type.End()) }
func (t *AccessorSig) Pos() token.Pos { return t.KwPos }
func (t *AccessorSig) End() token.Pos {
	return sigEnd(t.Sep, endOf(t.Rparen+1, t.Result))
}

func (t *ParamTypeList) Pos() token.Pos { return t.Lparen }
func (t *ParamTypeList) End() token.Pos { return t.Rparen + 1 }
func (t *ParamType) Pos() token.Pos     { return t.Name.Pos() }
func (t *ParamType) End() token.Pos {
	last := t.Name.End()
	if t.OptionalPos != token.NoPos && t.OptionalPos+1 > last {
		last = t.OptionalPos + 1
	}
	return endOf(last, t.Type)
}
func (t *TypeParamList) Pos() token.Pos { return t.Langle }
func (t *TypeParamList) End() token.Pos { return t.Rangle + 1 }
func (t *TypeParam) Pos() token.Pos {
	if n := len(t.Mods); n > 0 {
		return t.Mods[0].Pos
	}
	return t.Name.Pos()
}
func (t *TypeParam) End() token.Pos {
	return endOf(t.Name.End(), t.Constraint, t.Type, t.Default)
}
func (t *TypeArgList) Pos() token.Pos { return t.Langle }
func (t *TypeArgList) End() token.Pos { return t.Rangle + 1 }

func (*TypeRef) typeNode()             {}
func (*PredefinedType) typeNode()      {}
func (*LiteralType) typeNode()         {}
func (*ThisType) typeNode()            {}
func (*ParenType) typeNode()           {}
func (*ArrayType) typeNode()           {}
func (*IndexedAccessType) typeNode()   {}
func (*UnionType) typeNode()           {}
func (*IntersectionType) typeNode()    {}
func (*TypeOp) typeNode()              {}
func (*InferType) typeNode()           {}
func (*CondType) typeNode()            {}
func (*FuncType) typeNode()            {}
func (*CtorType) typeNode()            {}
func (*ObjectType) typeNode()          {}
func (*MappedType) typeNode()          {}
func (*TupleType) typeNode()           {}
func (*TupleElem) typeNode()           {}
func (*TypeQuery) typeNode()           {}
func (*ImportType) typeNode()          {}
func (*TemplateLiteralType) typeNode() {}
func (*TypePredicate) typeNode()       {}
func (*PropertySig) typeNode()         {}
func (*MethodSig) typeNode()           {}
func (*CallSig) typeNode()             {}
func (*ConstructSig) typeNode()        {}
func (*IndexSig) typeNode()            {}
func (*AccessorSig) typeNode()         {}