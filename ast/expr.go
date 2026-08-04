package ast

import "github.com/vertex-language/vertex/token"

// -------------------------------------------------------------- primaries

// BasicLit is an int_lit, float_lit, char_lit, string_lit, or one of the
// reserved literal keywords. Value is the raw source spelling — escapes
// unresolved, digit separators intact — because a formatter needs the original
// and the analyzer does its own decoding.
type BasicLit struct {
	ValuePos token.Pos
	Kind     token.Kind // INT, FLOAT, CHAR, STRING, TRUE, FALSE, NIL
	Value    string
}

// NamespaceExpr is a NamespaceName: async, gpu, or npu. It appears only as the
// operand of a Selector. `chan` is not one — it has its own constructor
// production — so the two never compete.
//
// The lookahead that separates `npu Dot(a, b)` from `npu.Dot(a, b)` is a parser
// decision and leaves no trace: the namespace reading is simply a SelectorExpr
// over this node.
type NamespaceExpr struct {
	KwPos token.Pos
	Kw    token.Kind // ASYNC, GPU, NPU
}

// ParenExpr is `( Expression )`, and also the parenthesized-type alternative of
// Type. One node, since a parenthesized single type is that type.
type ParenExpr struct {
	Lparen token.Pos
	X      Expr
	Rparen token.Pos
}

// TupleExpr is both TupleType and TupleLit.
//
// One node, because the two are the same shape once types are Exprs: a
// TupleElem is `[identifier ":"] Type` and a TupleElemValue is
// `[identifier ":"] OwningExpr`. A named element is a *KeyValueExpr.
//
// Elems is never empty — a tuple has at least one element, and there is no unit
// type. A single element requires TrailingComma; without it the construct is a
// ParenExpr instead, which is why this is the one trailing comma the tree
// records.
type TupleExpr struct {
	Lparen        token.Pos
	Elems         []Expr // len >= 1
	TrailingComma bool
	Rparen        token.Pos
}

// ArrayLit is `[ ElementList ]`. Elements are owning positions, so any may be a
// *TransferExpr.
type ArrayLit struct {
	Lbrack token.Pos
	Elems  []Expr
	Rbrack token.Pos
}

// CompositeLit is `LiteralType LiteralValue`, where LiteralType is a TypeName
// with optional TypeArgs. Type is never nil — a bare `{...}` is a MapLit.
//
// The punctuation is load-bearing: a composite literal constructs a struct, and
// a class is constructed by calling an initializer. The reader tells the two
// apart from the syntax alone, so the tree keeps them distinct too.
type CompositeLit struct {
	Type   Expr   // *Ident, *SelectorExpr, or *IndexExpr
	Lbrace token.Pos
	Elems  []Expr // *KeyValueExpr; the key is an identifier
	Rbrace token.Pos
}

// MapLit is a braced literal with no type prefix. Its keys are arbitrary
// expressions, unlike a CompositeLit's field names.
type MapLit struct {
	Lbrace token.Pos
	Elems  []Expr // *KeyValueExpr
	Rbrace token.Pos
}

// KeyValueExpr covers every `X : Y` pair — a FieldValue, a map KeyValue, a
// named tuple element, and a named Argument.
type KeyValueExpr struct {
	Key   Expr
	Colon token.Pos
	Value Expr
}

// EnumShorthand is `.identifier` or `.identifier(args)` in expression position.
// Legal only where the enum type is fixed by context, which is a static rule.
//
// In Pattern position a leading `.` is never this; see EnumPattern.
type EnumShorthand struct {
	Dot    token.Pos
	Name   *Ident
	Lparen token.Pos // NoPos when there is no argument list
	Args   []Expr
	Rparen token.Pos
}

// FuncLit is `func Signature Block`. It begins with all enclosing parse context
// cleared and re-establishes it from its own marker, so a closure written
// inside an async body does not inherit that body's context.
type FuncLit struct {
	Type *FuncType
	Body *BlockStmt
}

// ChanConstructor is `chan [ Type ] ( [ Expression ] )`, the only expression
// form of chan. Cap is the optional capacity.
//
// It gets a node because nothing else produces this shape: a ChanType writes no
// brackets, so there is no bracket node to reuse and no reading to defer.
type ChanConstructor struct {
	Chan   token.Pos
	Lbrack token.Pos
	Elem   Expr
	Rbrack token.Pos
	Lparen token.Pos
	Cap    Expr // nil when omitted
	Rparen token.Pos
}

// HeapConstructor is `unique(x)`, `shared(x)`, or `weak(x)`.
//
// It gets a node because without one the shape collides: `unique (T)` is an
// OwnershipType over a ParenExpr, and once types are Exprs that is
// indistinguishable from the constructor. The keyword spelling is what makes
// this a constructor rather than a call over a reserved name, so the tree says
// which reading the parser took.
type HeapConstructor struct {
	KwPos  token.Pos
	Kw     token.Kind // UNIQUE, SHARED, WEAK
	Lparen token.Pos
	X      Expr
	Rparen token.Pos
}

// --------------------------------------------------------------- postfix

type SelectorExpr struct {
	X   Expr
	Dot token.Pos
	Sel *Ident
}

// TupleIndexExpr is positional tuple access, `t.0`. Chains compose: `t.0.0`.
//
// The scanner has already decided this is an index rather than a float, under
// the one restriction to longest-match scanning, so the digits arrive as their
// own token. Text is the raw spelling; that it must be decimal and free of `_`
// is a static rule, and decoding belongs with every other literal's decoding, in
// the analyzer.
type TupleIndexExpr struct {
	X        Expr
	Dot      token.Pos
	IndexPos token.Pos
	Text     string
}

// IndexExpr is both Index and TypeArgs: `a[i]`, `a[low..high]`, and
// `Stack[int32]` are one node.
//
// Which reading applies is settled by what the operand denotes, not by shape,
// so the parser records the brackets and stops. Indices holds one expression
// for the Index reading and a TypeList for the TypeArgs reading; a slice is the
// Index reading whose single entry is a *BinaryExpr with Op DOTDOT.
//
// The same node carries a receiver type's TypeParameters and an OperandName's
// TypeArgs, which are the same brackets in another position.
type IndexExpr struct {
	X       Expr
	Lbrack  token.Pos
	Indices []Expr // len >= 1
	Rbrack  token.Pos
}

// CallExpr is `PrimaryExpr Arguments`, and with it three forms the grammar
// names separately but does not shape differently.
//
// A TypeOperatorCall — sizeof, alignof, reinterpret — is this node with an
// *Ident callee. Those are the only calls taking a Type in argument position,
// and the parser recognizes them by name, which is sound because a reserved
// builtin name may not be shadowed. The type argument leaves its own trace: a
// `[3]int32` first argument is an *ArrayType where an expression would have
// been an *ArrayLit.
//
// A VectorCall is this node with a *VectorType callee, which is what marks it
// as one; no ordinary call reading applies. An ExpectedType is this node with
// an *Ident callee spelled `Expected`, since that name and `error` are ordinary
// identifiers. Arity and argument shape are static rules in every case.
type CallExpr struct {
	Fun    Expr
	Lparen token.Pos
	Args   []Expr // *KeyValueExpr for a named argument; owning positions
	Rparen token.Pos
}

// LaunchExpr is a launch prefix applied to a call: thread, async, gpu, or npu.
// A prefix modifies scheduling only, never the callee's signature. Config is
// written only on gpu.
type LaunchExpr struct {
	KwPos  token.Pos
	Kw     token.Kind // THREAD, ASYNC, GPU, NPU
	Config *LaunchConfig
	Call   Expr // *CallExpr
}

// LaunchConfig is `( "blocks" : E , "threads" : E )`. Fixed arity and fixed
// names, so it is not a general argument list and the names are not recorded.
type LaunchConfig struct {
	Lparen  token.Pos
	Blocks  Expr
	Threads Expr
	Rparen  token.Pos
}

// AwaitExpr is `await X`. It parses unconditionally; whether the enclosing body
// licenses it is a static rule.
type AwaitExpr struct {
	Await token.Pos
	X     Expr
}

// -------------------------------------------------------------- operators

// UnaryExpr covers the three unary_op spellings `-`, `!`, `~`, and also `&`,
// which is not a unary_op but derives through PointerPrimary and so binds
// tighter than a selector: `&p.add(1)` is `(&p).add(1)`.
//
// Two of the four are deliberately unresolved here. `&` is address-of on a
// value and dereference on a typed_ptr, read from the operand's statically
// written type. `~` is bitwise-NOT in an expression and underlying-type in a
// TypeSetTerm. Neither is distinguished syntactically, so both are this node
// and the analyzer decides.
type UnaryExpr struct {
	OpPos token.Pos
	Op    token.Kind // SUB, NOT, TILDE, AND
	X     Expr
}

// BinaryExpr is one binary_op applied to two operands. The seven precedence
// levels are a property of the operator, read from token.Kind.Prec, not a
// nesting of tree shapes.
//
// DOTDOT is this node too. It is non-associative, which the parser rejects
// rather than folding — `a..b..c` is a compile error either way it might have
// been read.
type BinaryExpr struct {
	X     Expr
	OpPos token.Pos
	Op    token.Kind
	Y     Expr
}

// CastExpr is `x as T`. Left-associative and binding tighter than every binary
// operator, and written as its own production rather than as a binary_op
// because its right operand is a Type.
type CastExpr struct {
	X    Expr
	As   token.Pos
	Type Expr
}

// TransferExpr is the ownership marker, `var` applied to a UnaryExpr.
//
// It is a node rather than a flag because the marker is one production serving
// six owning positions, and one node covers all of them without six flags. Its
// presence is the entire difference between a move and a deep copy, so it must
// survive parsing as syntax and must never be normalized away.
//
// Target is written as a full UnaryExpr so that `var f(a)` and `var items[0]`
// parse. That the operand must be a binding or a field path is a static rule;
// there is no TransferTarget production, and nothing else in the grammar
// competes for this text.
type TransferExpr struct {
	Var    token.Pos
	Target Expr
}

// ------------------------------------------------------------------ types

// OwnershipType is `mut T`, `var T`, `unique T`, `shared T`, or `weak T`.
// Qualifiers do not stack; the recursion is unguarded so a stacked form parses
// and can be diagnosed as itself.
type OwnershipType struct {
	KwPos token.Pos
	Kw    token.Kind // MUT, VAR, UNIQUE, SHARED, WEAK
	X     Expr
}

// ArrayType is both ArrayType and SliceType: `[N]T` and `[]T`. Len == nil is
// the slice form.
type ArrayType struct {
	Lbrack token.Pos
	Len    Expr // nil for a slice
	Rbrack token.Pos
	Elem   Expr
}

type MapType struct {
	Map    token.Pos
	Lbrack token.Pos
	Key    Expr
	Rbrack token.Pos
	Value  Expr
}

// FuncType is a FunctionType, and also the signature half of a FuncDecl and a
// FuncLit. In a bare FunctionType every Param has a nil Name.
//
// Markers holds every FunctionMarker written. A signature carries at most one,
// but the repetition is written so that more than one parses, so all of them
// are kept and the extras are rejected later.
//
// Result carries `-> Type`, or an ExpectedType on a declaration, where it is a
// *CallExpr. Omitting it is the void form; there is no `void` type name.
type FuncType struct {
	Func    token.Pos
	Params  *ParamList
	Markers []*Marker
	Arrow   token.Pos
	Result  Expr
}

// ChanType is `chan T`. A channel type carries no direction.
type ChanType struct {
	Chan token.Pos
	Elem Expr
}

// PointerType is `typed_ptr T`. One may not be the direct base of another, so a
// nested form is written with parentheses and arrives with Elem a *ParenExpr.
type PointerType struct {
	Kw   token.Pos
	Elem Expr
}

// TensorType is `tensor[T, dims...]`. Legal only inside an npu-marked function
// body or that function's own signature; elsewhere it parses and is rejected.
type TensorType struct {
	Tensor token.Pos
	Lbrack token.Pos
	Elem   Expr
	Shape  []Expr // int_lits, len >= 1
	Rbrack token.Pos
}

// VectorType is `vector[T, N]`. It is legal wherever a Type is; where it may
// actually appear is a static rule. As the callee of a CallExpr it makes that
// call a VectorCall.
type VectorType struct {
	Vector token.Pos
	Lbrack token.Pos
	Elem   Expr
	Comma  token.Pos
	Len    Expr // int_lit
	Rbrack token.Pos
}

// AbstractType is the bare `abstract`, legal only as an alias target.
type AbstractType struct {
	Abstract token.Pos
}

// BadExpr marks a span the parser could not make sense of. It exists so that
// recovery yields a tree the analyzer can still walk; the analyzer skips it
// silently, a diagnostic having already been reported at parse time.
type BadExpr struct {
	From, To token.Pos
}

// -------------------------------------------------------------- positions

func (x *BasicLit) Pos() token.Pos        { return x.ValuePos }
func (x *NamespaceExpr) Pos() token.Pos   { return x.KwPos }
func (x *ParenExpr) Pos() token.Pos       { return x.Lparen }
func (x *TupleExpr) Pos() token.Pos       { return x.Lparen }
func (x *ArrayLit) Pos() token.Pos        { return x.Lbrack }
func (x *CompositeLit) Pos() token.Pos    { return x.Type.Pos() }
func (x *MapLit) Pos() token.Pos          { return x.Lbrace }
func (x *KeyValueExpr) Pos() token.Pos    { return x.Key.Pos() }
func (x *EnumShorthand) Pos() token.Pos   { return x.Dot }
func (x *FuncLit) Pos() token.Pos         { return x.Type.Pos() }
func (x *ChanConstructor) Pos() token.Pos { return x.Chan }
func (x *HeapConstructor) Pos() token.Pos { return x.KwPos }
func (x *SelectorExpr) Pos() token.Pos    { return x.X.Pos() }
func (x *TupleIndexExpr) Pos() token.Pos  { return x.X.Pos() }
func (x *IndexExpr) Pos() token.Pos       { return x.X.Pos() }
func (x *CallExpr) Pos() token.Pos        { return x.Fun.Pos() }
func (x *LaunchExpr) Pos() token.Pos      { return x.KwPos }
func (x *LaunchConfig) Pos() token.Pos    { return x.Lparen }
func (x *AwaitExpr) Pos() token.Pos       { return x.Await }
func (x *UnaryExpr) Pos() token.Pos       { return x.OpPos }
func (x *BinaryExpr) Pos() token.Pos      { return x.X.Pos() }
func (x *CastExpr) Pos() token.Pos        { return x.X.Pos() }
func (x *TransferExpr) Pos() token.Pos    { return x.Var }
func (x *OwnershipType) Pos() token.Pos   { return x.KwPos }
func (x *ArrayType) Pos() token.Pos       { return x.Lbrack }
func (x *MapType) Pos() token.Pos         { return x.Map }
func (x *FuncType) Pos() token.Pos        { return x.Func }
func (x *ChanType) Pos() token.Pos        { return x.Chan }
func (x *PointerType) Pos() token.Pos     { return x.Kw }
func (x *TensorType) Pos() token.Pos      { return x.Tensor }
func (x *VectorType) Pos() token.Pos      { return x.Vector }
func (x *AbstractType) Pos() token.Pos    { return x.Abstract }
func (x *BadExpr) Pos() token.Pos         { return x.From }

func (x *BasicLit) End() token.Pos      { return x.ValuePos + token.Pos(len(x.Value)) }
func (x *ParenExpr) End() token.Pos     { return x.Rparen + 1 }
func (x *TupleExpr) End() token.Pos     { return x.Rparen + 1 }
func (x *ArrayLit) End() token.Pos      { return x.Rbrack + 1 }
func (x *CompositeLit) End() token.Pos  { return x.Rbrace + 1 }
func (x *MapLit) End() token.Pos        { return x.Rbrace + 1 }
func (x *KeyValueExpr) End() token.Pos  { return x.Value.End() }
func (x *FuncLit) End() token.Pos       { return x.Body.End() }
func (x *ChanConstructor) End() token.Pos { return x.Rparen + 1 }
func (x *HeapConstructor) End() token.Pos { return x.Rparen + 1 }
func (x *SelectorExpr) End() token.Pos  { return x.Sel.End() }
func (x *IndexExpr) End() token.Pos     { return x.Rbrack + 1 }
func (x *CallExpr) End() token.Pos      { return x.Rparen + 1 }
func (x *LaunchExpr) End() token.Pos    { return x.Call.End() }
func (x *LaunchConfig) End() token.Pos  { return x.Rparen + 1 }
func (x *AwaitExpr) End() token.Pos     { return x.X.End() }
func (x *UnaryExpr) End() token.Pos     { return x.X.End() }
func (x *BinaryExpr) End() token.Pos    { return x.Y.End() }
func (x *CastExpr) End() token.Pos      { return x.Type.End() }
func (x *TransferExpr) End() token.Pos  { return x.Target.End() }
func (x *OwnershipType) End() token.Pos { return x.X.End() }
func (x *ArrayType) End() token.Pos     { return x.Elem.End() }
func (x *MapType) End() token.Pos       { return x.Value.End() }
func (x *ChanType) End() token.Pos      { return x.Elem.End() }
func (x *PointerType) End() token.Pos   { return x.Elem.End() }
func (x *TensorType) End() token.Pos    { return x.Rbrack + 1 }
func (x *VectorType) End() token.Pos    { return x.Rbrack + 1 }
func (x *BadExpr) End() token.Pos       { return x.To }

func (x *NamespaceExpr) End() token.Pos {
	return x.KwPos + token.Pos(len(x.Kw.Spelling()))
}

func (x *AbstractType) End() token.Pos {
	return x.Abstract + token.Pos(len(token.ABSTRACT.Spelling()))
}

func (x *TupleIndexExpr) End() token.Pos {
	return x.IndexPos + token.Pos(len(x.Text))
}

func (x *EnumShorthand) End() token.Pos {
	if x.Rparen.IsValid() {
		return x.Rparen + 1
	}
	return x.Name.End()
}

func (x *FuncType) End() token.Pos {
	if x.Result != nil {
		return x.Result.End()
	}
	if n := len(x.Markers); n > 0 {
		return x.Markers[n-1].End()
	}
	return x.Params.End()
}

func (*BasicLit) exprNode()        {}
func (*NamespaceExpr) exprNode()   {}
func (*ParenExpr) exprNode()       {}
func (*TupleExpr) exprNode()       {}
func (*ArrayLit) exprNode()        {}
func (*CompositeLit) exprNode()    {}
func (*MapLit) exprNode()          {}
func (*KeyValueExpr) exprNode()    {}
func (*EnumShorthand) exprNode()   {}
func (*FuncLit) exprNode()         {}
func (*ChanConstructor) exprNode() {}
func (*HeapConstructor) exprNode() {}
func (*SelectorExpr) exprNode()    {}
func (*TupleIndexExpr) exprNode()  {}
func (*IndexExpr) exprNode()       {}
func (*CallExpr) exprNode()        {}
func (*LaunchExpr) exprNode()      {}
func (*AwaitExpr) exprNode()       {}
func (*UnaryExpr) exprNode()       {}
func (*BinaryExpr) exprNode()      {}
func (*CastExpr) exprNode()        {}
func (*TransferExpr) exprNode()    {}
func (*OwnershipType) exprNode()   {}
func (*ArrayType) exprNode()       {}
func (*MapType) exprNode()         {}
func (*FuncType) exprNode()        {}
func (*ChanType) exprNode()        {}
func (*PointerType) exprNode()     {}
func (*TensorType) exprNode()      {}
func (*VectorType) exprNode()      {}
func (*AbstractType) exprNode()    {}
func (*BadExpr) exprNode()         {}