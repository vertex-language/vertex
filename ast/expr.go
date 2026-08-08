package ast

import "github.com/vertex-language/vertex/token"

// Expressions — vertex_grammar.md section B.

type (
	// BasicLit is a Literal (B.1): NullLiteral, BooleanLiteral,
	// NumericLiteral, StringLiteral, and RegularExpressionLiteral.
	//
	// The value is not here and never will be. `1_024` is five bytes of raw
	// spelling, `0b1010` is six; decoding needs the target width and belongs
	// to a phase that knows it (§4.6, §8.3). Recover the text with
	// file.Between(lit.Pos(), lit.End()).
	BasicLit struct {
		Kind     token.Kind // NUMBER, BIGINT, STRING, REGEX, TRUE, FALSE, NULL
		ValuePos token.Pos
		ValueEnd token.Pos
		HasEscape bool // separators or escape sequences present in the spelling
	}

	// TemplateLit is a TemplateLiteral (A.6). Quasis and Exprs interleave:
	// len(Quasis) == len(Exprs)+1, always.
	//
	// A tagged template is TaggedTemplateExpr wrapping this one.
	TemplateLit struct {
		Quasis []*TemplateElem // len >= 1
		Exprs  []Expr          // len == len(Quasis)-1
	}

	// TemplateElem is one NoSubstitutionTemplate, TemplateHead,
	// TemplateMiddle, or TemplateTail token. Raw spelling, uncooked: whether a
	// NotEscapeSequence is legal depends on tagging, which is not this node's
	// business.
	TemplateElem struct {
		Kind  token.Kind
		Start token.Pos
		Stop  token.Pos
	}

	// ArrayLit is an ArrayLiteral (B.2). Elts may contain *Elision.
	ArrayLit struct {
		Lbrack token.Pos
		Elts   []Expr
		Rbrack token.Pos
	}

	// Elision is a hole in an array literal (B.2). It is a real node rather
	// than a nil slot: a hole is meaningful, and §5.4's "omit what's inert"
	// covers trailing commas, not these.
	Elision struct {
		Comma token.Pos
	}

	// ObjectLit is an ObjectLiteral (B.2).
	//
	// CoverInit marks that at least one PropertyDef used CoverInitializedName
	// (grammar K) — `{ a = 1 }`. That form is legal only after reinterpretation
	// to an ObjectAssignmentPattern (B.6), which is the grammar's only cover
	// and the only in-parser reinterpretation (§5.1). If this flag survives to
	// a node that was never reinterpreted, a later phase rejects it by name.
	ObjectLit struct {
		Lbrace    token.Pos
		Props     []Expr // *PropertyDef, *SpreadElem, *Ident (shorthand), *MethodDef
		Rbrace    token.Pos
		CoverInit bool
	}

	// PropertyDef is `key: value`, `key = init` (the cover form), or a
	// computed key (B.2).
	PropertyDef struct {
		Key      Expr // *Ident, *BasicLit, or *ComputedKey
		Colon    token.Pos // NoPos for the cover form
		Value    Expr
		IsCover  bool // CoverInitializedName; Colon is NoPos and Value is the initializer
	}

	// ComputedKey is ComputedPropertyName (B.2), `[expr]`.
	ComputedKey struct {
		Lbrack token.Pos
		X      Expr
		Rbrack token.Pos
	}

	// SpreadElem is `... expr` in an array literal, object literal, or
	// argument list (B.2, B.3).
	SpreadElem struct {
		Ellipsis token.Pos
		X        Expr
	}

	// ParenExpr is retained rather than folded away (§5.4):
	// `(makeBox<boolean>)(true)` does not read the same without it. Consumers
	// that don't care use Unparen.
	ParenExpr struct {
		Lparen token.Pos
		X      Expr
		Rparen token.Pos
	}

	// FuncExpr is FunctionExpression, GeneratorExpression,
	// AsyncFunctionExpression, or AsyncGeneratorExpression (D).
	//
	// It holds a *FuncDecl rather than repeating its fields, mirroring the way
	// ClassExpression and ClassDeclaration share ClassTail (E). Decl.Name is
	// nil for the anonymous forms.
	FuncExpr struct {
		Fn *FuncDecl
	}

	// ClassExpr is ClassExpression (E), sharing ClassDecl for the same reason.
	ClassExpr struct {
		Class *ClassDecl
	}

	// ArrowFunc is ArrowFunction or AsyncArrowFunction (D.2).
	//
	// Params is nil when the head was a bare BindingIdentifier, in which case
	// Ident is set. The two forms are not normalized into one: a bare
	// identifier head has no parentheses to point at, and inventing them would
	// be synthesis (§1).
	ArrowFunc struct {
		AsyncPos   token.Pos // NoPos if not async
		Ident      *Ident    // set iff Params == nil
		TypeParams *TypeParamList
		Params     *ParamList
		Result     TypeExpr // may be *TypePredicate
		Arrow      token.Pos
		Body       Node // Expr (ExpressionBody) or *BlockStmt
	}

	// MemberExpr is `x.name`, `x?.name`, `x.#name`, and `x?.#name` (B.3).
	MemberExpr struct {
		X        Expr
		Optional bool  // reached through ?.
		Dot      token.Pos
		Sel      Node  // *Ident or *PrivateIdent
	}

	// IndexExpr is `x[i]` and `x?.[i]` (B.3).
	IndexExpr struct {
		X        Expr
		Optional bool
		Lbrack   token.Pos
		Index    Expr
		Rbrack   token.Pos
	}

	// CallExpr is CallExpression, SuperCall, and the optional-chain call forms
	// (B.3). TypeArgs is set for `f<T>(x)`.
	CallExpr struct {
		Fun      Expr
		Optional bool
		TypeArgs *TypeArgList
		Lparen   token.Pos
		Args     []Expr // may contain *SpreadElem
		Rparen   token.Pos
	}

	// NewExpr is `new C(...)` and the bare `new C` form (B.3), where Lparen is
	// NoPos and Args is nil.
	NewExpr struct {
		NewPos   token.Pos
		Callee   Expr
		TypeArgs *TypeArgList
		Lparen   token.Pos
		Args     []Expr
		Rparen   token.Pos
	}

	// TaggedTemplateExpr is `tag`...`` (B.3).
	TaggedTemplateExpr struct {
		Tag      Expr
		TypeArgs *TypeArgList
		Template *TemplateLit
	}

	// SuperExpr is the `super` keyword in SuperProperty and SuperCall (B.3).
	// It is a leaf; the property or call wrapping it is a MemberExpr,
	// IndexExpr, or CallExpr.
	SuperExpr struct {
		SuperPos token.Pos
		SuperEnd token.Pos
	}

	// ThisExpr is `this` (B.1).
	ThisExpr struct {
		ThisPos token.Pos
		ThisEnd token.Pos
	}

	// MetaProp is NewTarget or ImportMeta (B.3).
	MetaProp struct {
		MetaPos token.Pos
		Meta    token.Kind // NEW or IMPORT
		Prop    *Ident     // target or meta
	}

	// ImportCall is ImportCall (B.3), including the `import.defer(...)` and
	// `import.source(...)` phase forms.
	ImportCall struct {
		ImportPos token.Pos
		Phase     ImportPhase
		PhasePos  token.Pos // NoPos when Phase == PhaseEval
		Lparen    token.Pos
		Args      []Expr // 1 or 2; a second is the options object
		Rparen    token.Pos
	}

	// InstantiationExpr is `f<T>` with no call following (B.3) — an
	// instantiation expression, committed only when TypeArguments closes and
	// the next token is in InstantiationFollowSet. See §6.1.
	InstantiationExpr struct {
		X        Expr
		TypeArgs *TypeArgList
	}

	// NonNullExpr is the postfix `!` (B.3), under a [no LineTerminator here]
	// restriction.
	NonNullExpr struct {
		X    Expr
		Bang token.Pos
		BangEnd token.Pos
	}

	// TypeAssertExpr is the prefix form `<Type> x` (B.4).
	//
	// There is exactly one reading of a prefix `<` because there is one goal
	// symbol and no per-file mode (§1). Nothing competes for this position.
	TypeAssertExpr struct {
		Langle token.Pos
		Type   TypeExpr
		Rangle token.Pos
		X      Expr
	}

	// AsExpr is `x as T`, `x as const`, and `x satisfies T` (B.4). One node
	// for three because they differ only in the keyword and in what a later
	// phase does with them.
	AsExpr struct {
		X       Expr
		OpPos   token.Pos
		Op      token.Ctx // CtxAs or CtxSatisfies
		IsConst bool      // `as const`; Type is nil
		Type    TypeExpr
		ConstEnd token.Pos // set iff IsConst
	}

	// UnaryExpr is delete, void, typeof, +, -, ~, ! (B.4).
	UnaryExpr struct {
		OpPos token.Pos
		Op    token.Kind
		X     Expr
	}

	// UpdateExpr is ++ and --, prefix or postfix (B.4).
	UpdateExpr struct {
		OpPos   token.Pos
		OpEnd   token.Pos
		Op      token.Kind
		Prefix  bool
		X       Expr
	}

	// AwaitExpr is AwaitExpression (B.5).
	AwaitExpr struct {
		AwaitPos token.Pos
		X        Expr
	}

	// YieldExpr is YieldExpression (B.5). X is nil for a bare `yield`.
	YieldExpr struct {
		YieldPos token.Pos
		YieldEnd token.Pos
		Delegate bool // yield *
		X        Expr
	}

	// BinaryExpr covers the whole binary chain in B.4, including `in`,
	// `instanceof`, the logical operators, and `??`.
	//
	// Op is the joined kind for a `>` run: the scanner emits single GT tokens
	// (§4.2) and the expression parser rejoins them with token.JoinGT before
	// building this node. OpEnd therefore covers the whole run.
	//
	// Mixing `??` with `||` or `&&` without parentheses is ungrammatical in
	// B.4 but parses here, per §6.3, and is rejected by name later.
	BinaryExpr struct {
		X     Expr
		OpPos token.Pos
		OpEnd token.Pos
		Op    token.Kind
		Y     Expr
	}

	// AssignExpr is simple and compound assignment, including &&=, ||=, ??=
	// (B.5). Lhs may be a pattern after reinterpretation (B.6).
	AssignExpr struct {
		Lhs   Expr
		OpPos token.Pos
		OpEnd token.Pos
		Op    token.Kind
		Rhs   Expr
	}

	// CondExpr is ConditionalExpression (B.4).
	CondExpr struct {
		Cond  Expr
		Quest token.Pos
		Then  Expr
		Colon token.Pos
		Else  Expr
	}

	// SeqExpr is the comma operator (B.5). len(Exprs) >= 2.
	SeqExpr struct {
		Exprs []Expr
	}
)

// ImportPhase distinguishes the ImportCall forms in B.3.
type ImportPhase uint8

const (
	PhaseEval ImportPhase = iota
	PhaseDefer
	PhaseSource
)

// --- Destructuring patterns (B.6, C.2) ------------------------------------
//
// One family serves both. ObjectBindingPattern (C.2) parses directly into
// these nodes, since declaration position is known without a cover.
// ObjectAssignmentPattern (B.6) arrives by reinterpreting an ObjectLit, which
// §5.1 identifies as the only in-parser reinterpretation.
//
// They are Exprs because DestructuringAssignmentTarget is a
// LeftHandSideExpression (B.6), and a binding position that admits a pattern
// admits nothing an Expr could not hold.

type (
	// ObjectPattern is ObjectBindingPattern (C.2) or a reinterpreted
	// ObjectAssignmentPattern (B.6).
	ObjectPattern struct {
		Lbrace token.Pos
		Props  []Expr // *PropertyPattern or *RestElem
		Rbrace token.Pos
	}

	// ArrayPattern is ArrayBindingPattern or ArrayAssignmentPattern. Elts may
	// contain *Elision and a trailing *RestElem.
	ArrayPattern struct {
		Lbrack token.Pos
		Elts   []Expr
		Rbrack token.Pos
	}

	// PropertyPattern is one BindingProperty or AssignmentProperty. Key is nil
	// for the shorthand form, where Value is the *Ident or *AssignPattern.
	PropertyPattern struct {
		Key   Expr // nil for shorthand
		Colon token.Pos
		Value Expr
	}

	// AssignPattern is a target with a default: `a = 1` in a binding or
	// assignment pattern. Distinct from AssignExpr, which is an operator.
	AssignPattern struct {
		Lhs   Expr
		Assign token.Pos
		Rhs   Expr
	}

	// RestElem is BindingRestElement, BindingRestProperty,
	// AssignmentRestElement, or AssignmentRestProperty. Type is set only in a
	// FunctionRestParameter (D.1).
	RestElem struct {
		Ellipsis token.Pos
		X        Expr
		Type     TypeExpr
	}
)

// --- spans ------------------------------------------------------------------

func (x *BasicLit) Pos() token.Pos     { return x.ValuePos }
func (x *BasicLit) End() token.Pos     { return x.ValueEnd }
func (x *TemplateElem) Pos() token.Pos { return x.Start }
func (x *TemplateElem) End() token.Pos { return x.Stop }
func (x *TemplateLit) Pos() token.Pos  { return x.Quasis[0].Pos() }
func (x *TemplateLit) End() token.Pos  { return x.Quasis[len(x.Quasis)-1].End() }
func (x *ArrayLit) Pos() token.Pos     { return x.Lbrack }
func (x *ArrayLit) End() token.Pos     { return x.Rbrack + 1 }
func (x *Elision) Pos() token.Pos      { return x.Comma }
func (x *Elision) End() token.Pos      { return x.Comma + 1 }
func (x *ObjectLit) Pos() token.Pos    { return x.Lbrace }
func (x *ObjectLit) End() token.Pos    { return x.Rbrace + 1 }
func (x *PropertyDef) Pos() token.Pos  { return x.Key.Pos() }
func (x *PropertyDef) End() token.Pos  { return endOf(x.Key.End(), x.Value) }
func (x *ComputedKey) Pos() token.Pos  { return x.Lbrack }
func (x *ComputedKey) End() token.Pos  { return x.Rbrack + 1 }
func (x *SpreadElem) Pos() token.Pos   { return x.Ellipsis }
func (x *SpreadElem) End() token.Pos   { return x.X.End() }
func (x *ParenExpr) Pos() token.Pos    { return x.Lparen }
func (x *ParenExpr) End() token.Pos    { return x.Rparen + 1 }
func (x *FuncExpr) Pos() token.Pos     { return x.Fn.Pos() }
func (x *FuncExpr) End() token.Pos     { return x.Fn.End() }
func (x *ClassExpr) Pos() token.Pos    { return x.Class.Pos() }
func (x *ClassExpr) End() token.Pos    { return x.Class.End() }

func (x *ArrowFunc) Pos() token.Pos {
	if x.AsyncPos != token.NoPos {
		return x.AsyncPos
	}
	return spanOf(x.Arrow, x.Ident, x.TypeParams, x.Params)
}
func (x *ArrowFunc) End() token.Pos { return x.Body.End() }

func (x *MemberExpr) Pos() token.Pos { return x.X.Pos() }
func (x *MemberExpr) End() token.Pos { return x.Sel.End() }
func (x *IndexExpr) Pos() token.Pos  { return x.X.Pos() }
func (x *IndexExpr) End() token.Pos  { return x.Rbrack + 1 }
func (x *CallExpr) Pos() token.Pos   { return x.Fun.Pos() }
func (x *CallExpr) End() token.Pos   { return x.Rparen + 1 }
func (x *NewExpr) Pos() token.Pos    { return x.NewPos }
func (x *NewExpr) End() token.Pos {
	if x.Rparen != token.NoPos {
		return x.Rparen + 1
	}
	return endOf(x.Callee.End(), x.TypeArgs)
}
func (x *TaggedTemplateExpr) Pos() token.Pos { return x.Tag.Pos() }
func (x *TaggedTemplateExpr) End() token.Pos { return x.Template.End() }
func (x *SuperExpr) Pos() token.Pos          { return x.SuperPos }
func (x *SuperExpr) End() token.Pos          { return x.SuperEnd }
func (x *ThisExpr) Pos() token.Pos           { return x.ThisPos }
func (x *ThisExpr) End() token.Pos           { return x.ThisEnd }
func (x *MetaProp) Pos() token.Pos           { return x.MetaPos }
func (x *MetaProp) End() token.Pos           { return x.Prop.End() }
func (x *ImportCall) Pos() token.Pos         { return x.ImportPos }
func (x *ImportCall) End() token.Pos         { return x.Rparen + 1 }
func (x *InstantiationExpr) Pos() token.Pos  { return x.X.Pos() }
func (x *InstantiationExpr) End() token.Pos  { return x.TypeArgs.End() }
func (x *NonNullExpr) Pos() token.Pos        { return x.X.Pos() }
func (x *NonNullExpr) End() token.Pos        { return x.BangEnd }
func (x *TypeAssertExpr) Pos() token.Pos     { return x.Langle }
func (x *TypeAssertExpr) End() token.Pos     { return x.X.End() }
func (x *AsExpr) Pos() token.Pos             { return x.X.Pos() }
func (x *AsExpr) End() token.Pos {
	if x.IsConst {
		return x.ConstEnd
	}
	return x.Type.End()
}
func (x *UnaryExpr) Pos() token.Pos { return x.OpPos }
func (x *UnaryExpr) End() token.Pos { return x.X.End() }
func (x *UpdateExpr) Pos() token.Pos {
	if x.Prefix {
		return x.OpPos
	}
	return x.X.Pos()
}
func (x *UpdateExpr) End() token.Pos {
	if x.Prefix {
		return x.X.End()
	}
	return x.OpEnd
}
func (x *AwaitExpr) Pos() token.Pos { return x.AwaitPos }
func (x *AwaitExpr) End() token.Pos { return x.X.End() }
func (x *YieldExpr) Pos() token.Pos { return x.YieldPos }
func (x *YieldExpr) End() token.Pos { return endOf(x.YieldEnd, x.X) }
func (x *BinaryExpr) Pos() token.Pos { return x.X.Pos() }
func (x *BinaryExpr) End() token.Pos { return x.Y.End() }
func (x *AssignExpr) Pos() token.Pos { return x.Lhs.Pos() }
func (x *AssignExpr) End() token.Pos { return x.Rhs.End() }
func (x *CondExpr) Pos() token.Pos   { return x.Cond.Pos() }
func (x *CondExpr) End() token.Pos   { return x.Else.End() }
func (x *SeqExpr) Pos() token.Pos    { return x.Exprs[0].Pos() }
func (x *SeqExpr) End() token.Pos    { return x.Exprs[len(x.Exprs)-1].End() }

func (x *ObjectPattern) Pos() token.Pos   { return x.Lbrace }
func (x *ObjectPattern) End() token.Pos   { return x.Rbrace + 1 }
func (x *ArrayPattern) Pos() token.Pos    { return x.Lbrack }
func (x *ArrayPattern) End() token.Pos    { return x.Rbrack + 1 }
func (x *PropertyPattern) Pos() token.Pos { return spanOf(x.Value.Pos(), x.Key) }
func (x *PropertyPattern) End() token.Pos { return x.Value.End() }
func (x *AssignPattern) Pos() token.Pos   { return x.Lhs.Pos() }
func (x *AssignPattern) End() token.Pos   { return x.Rhs.End() }
func (x *RestElem) Pos() token.Pos        { return x.Ellipsis }
func (x *RestElem) End() token.Pos        { return endOf(x.X.End(), x.Type) }

func (*BasicLit) exprNode()           {}
func (*TemplateLit) exprNode()        {}
func (*ArrayLit) exprNode()           {}
func (*Elision) exprNode()            {}
func (*ObjectLit) exprNode()          {}
func (*PropertyDef) exprNode()        {}
func (*ComputedKey) exprNode()        {}
func (*SpreadElem) exprNode()         {}
func (*ParenExpr) exprNode()          {}
func (*FuncExpr) exprNode()           {}
func (*ClassExpr) exprNode()          {}
func (*ArrowFunc) exprNode()          {}
func (*MemberExpr) exprNode()         {}
func (*IndexExpr) exprNode()          {}
func (*CallExpr) exprNode()           {}
func (*NewExpr) exprNode()            {}
func (*TaggedTemplateExpr) exprNode() {}
func (*SuperExpr) exprNode()          {}
func (*ThisExpr) exprNode()           {}
func (*MetaProp) exprNode()           {}
func (*ImportCall) exprNode()         {}
func (*InstantiationExpr) exprNode()  {}
func (*NonNullExpr) exprNode()        {}
func (*TypeAssertExpr) exprNode()     {}
func (*AsExpr) exprNode()             {}
func (*UnaryExpr) exprNode()          {}
func (*UpdateExpr) exprNode()         {}
func (*AwaitExpr) exprNode()          {}
func (*YieldExpr) exprNode()          {}
func (*BinaryExpr) exprNode()         {}
func (*AssignExpr) exprNode()         {}
func (*CondExpr) exprNode()           {}
func (*SeqExpr) exprNode()            {}
func (*ObjectPattern) exprNode()      {}
func (*ArrayPattern) exprNode()       {}
func (*PropertyPattern) exprNode()    {}
func (*AssignPattern) exprNode()      {}
func (*RestElem) exprNode()           {}