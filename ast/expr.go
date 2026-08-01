package ast

import "github.com/vertex-language/vertex/token"

// ------------------------------------------------------------- primaries

// BasicLit is a NumericLiteral, StringLiteral, CharLiteral, or a
// ReservedLiteralKeyword (true/false/nil), which A.1.3 makes Literals
// syntactically. Value is the raw source spelling, never unescaped — vfmt needs
// the original, and the analyzer does its own decoding.
type BasicLit struct {
	ValuePos token.Pos
	Kind     token.Kind // INT, FLOAT, CHAR, STRING, TRUE, FALSE, NIL
	Value    string
}

// NamespaceExpr is one of the four keyword namespaces of A.4.1: async, gpu,
// npu, chan. It appears only as the X of a SelectorExpr (async.Readable) or, for
// chan, as the head of a construction (chan[float32](64)) — the lookahead
// restrictions in A.4.2 are what keep these distinct from launch prefixes.
type NamespaceExpr struct {
	KwPos token.Pos
	Kw    token.Kind // ASYNC, GPU, NPU, CHAN
}

type ParenExpr struct {
	Lparen token.Pos
	X      Expr
	Rparen token.Pos
}

// TupleExpr is TupleType, TupleLiteral, and UnitType (A.3.1, A.4.7).
//
// One node, because the three are the same shape: TupleElement is
// `Type | Identifier : Type` and TupleElementValue is
// `OwningExpression | Identifier : OwningExpression`. With types represented as
// Exprs there is nothing left to tell them apart, and sizeof((int32, bool))
// reaches the expression parser where the type parser never runs.
//
// A named element is a *KeyValueExpr. Empty Elems with no trailing comma is the
// unit type `()`. A single element requires TrailingComma (A.4.7); `(x)` without
// one is a ParenExpr.
type TupleExpr struct {
	Lparen        token.Pos
	Elems         []Expr
	TrailingComma bool
	Rparen        token.Pos
}

// ArrayLit is `[ ElementList ]` (A.4.7). Elements are owning positions, so any
// may be a *TransferExpr.
type ArrayLit struct {
	Lbrack token.Pos
	Elems  []Expr
	Rbrack token.Pos
}

// CompositeLit is `TypeName LiteralBody` or `InstantiatedType LiteralBody`
// (A.4.7). Type is never nil — a bare `{...}` is a MapLit.
//
// A.4.7 makes the punctuation load-bearing: a struct is built with this node, a
// class by calling its init (A.6.4). The reader tells the storage discipline
// from the syntax alone, so the tree keeps them distinct too.
type CompositeLit struct {
	Type   Expr // *Ident, *SelectorExpr, or *IndexExpr
	Lbrace token.Pos
	Elems  []Expr // *KeyValueExpr; A.4.7 requires the key to be an Identifier
	Rbrace token.Pos
}

// MapLit is a braced literal with no type prefix (A.4.7). Its keys are
// arbitrary Expressions, unlike a CompositeLit's field names.
type MapLit struct {
	Lbrace token.Pos
	Elems  []Expr // *KeyValueExpr
	Rbrace token.Pos
}

// KeyValueExpr covers every `X : Y` pair: a composite-literal field, a map
// entry, a named tuple element, and a named call argument (A.4.3, A.4.7).
type KeyValueExpr struct {
	Key   Expr
	Colon token.Pos
	Value Expr
}

// EnumShorthand is `.Name` or `.Name(args)` (A.4.1). Legal only where the enum
// type is inferable from context, which is a static rule.
type EnumShorthand struct {
	Dot    token.Pos
	Name   *Ident
	Lparen token.Pos // NoPos when there is no argument list
	Args   []Expr
	Rparen token.Pos
}

// FuncLit is a FunctionExpression (A.4.1). A.0.3: it begins with all four
// context parameters cleared and re-sets them from its own Marker, so an
// anonymous closure inside an async body may not await unless itself marked.
type FuncLit struct {
	Type *FuncType
	Body *BlockStmt
}

// -------------------------------------------------------------- postfix

type SelectorExpr struct {
	X   Expr
	Dot token.Pos
	Sel *Ident
}

// TupleIndexExpr is positional tuple access, `t.0` (A.4.3). Chains compose:
// `t.0.0`.
//
// Lexing note: maximal munch scans `t.0.0` as IDENT, PERIOD, FLOAT("0.0"). The
// parser must split a FLOAT of the form <digits>.<digits> when it appears
// immediately after a selector dot. Index holds the decoded value; Text holds
// the source spelling of this component only.
type TupleIndexExpr struct {
	X        Expr
	Dot      token.Pos
	IndexPos token.Pos
	Index    int
	Text     string
}

// IndexExpr is `x[a]`, `x[a..b]`, and `Stack[int32]` — index, slice, generic
// instantiation, and type-argument list, all one node (A.3.6).
//
// A.3.6 calls this the language's general compile-time configuration slot; the
// same node therefore also carries `new[T]`, `npu.Quantize[T]`, and `chan[T]`.
// A slice is an IndexExpr whose single Index is a BinaryExpr with Op DOTDOT.
type IndexExpr struct {
	X       Expr
	Lbrack  token.Pos
	Indices []Expr // len >= 1
	Rbrack  token.Pos
}

// CallExpr is `f(a, b)` and `f[T](a)` (the latter with Fun == *IndexExpr).
//
// Builtins get no node of their own: A.4.8 says they are ordinary call syntax
// over reserved names, and A.1.4 keeps those names unshadowable. sizeof(Type)
// works because types are Exprs; `align:` and `zero:` are ordinary named
// arguments. Arity and type-argument shape are static rules.
//
// ExpectedType (A.12.2) is also a CallExpr — `Expected` is a plain identifier
// and so is `error`.
type CallExpr struct {
	Fun    Expr
	Lparen token.Pos
	Args   []Expr // *KeyValueExpr for a named argument; *TransferExpr where Own is set
	Rparen token.Pos
}

// LaunchExpr is a launch prefix applied to a call (A.4.2): thread, async, gpu,
// npu. A.4.2 is explicit that these modify scheduling, never the callee's
// signature. Config is legal only on gpu.
//
// The [lookahead != .] restriction that separates `npu Dot(a,b)` from
// `npu.Dot(a,b)` is a parser decision and leaves no trace: a namespace access
// produces SelectorExpr{X: NamespaceExpr}.
type LaunchExpr struct {
	KwPos  token.Pos
	Kw     token.Kind // THREAD, ASYNC, GPU, NPU
	Config *LaunchConfig
	Call   Expr // *CallExpr
}

// LaunchConfig is `(blocks: E, threads: E)` (A.4.2). Fixed arity and fixed
// names, so it is not a general argument list.
type LaunchConfig struct {
	Lparen  token.Pos
	Blocks  Expr
	Threads Expr
	Rparen  token.Pos
}

type AwaitExpr struct {
	Await token.Pos
	X     Expr
}

// -------------------------------------------------------------- operators

// UnaryExpr covers `-`, `!`, `~`, and `&`.
//
// Two of those are deliberately unresolved here. `&` is address-of on a value
// and dereference on a typed_ptr, keyed on the operand's statically written
// type (A.4.4). `~` is bitwise-NOT in an expression and underlying-type in a
// type-set element (A.7.3). Both are one node; the analyzer decides.
type UnaryExpr struct {
	OpPos token.Pos
	Op    token.Kind
	X     Expr
}

// BinaryExpr collapses A.4.5's cascade. Precedence comes from token.Kind.Prec()
// and A.13; the cascade nonterminals are grammar-writing devices, not shapes.
//
// DOTDOT is a BinaryExpr too. A.4.5 makes it non-associative, which the parser
// rejects rather than folding — `a..b..c` is a compile error.
type BinaryExpr struct {
	X     Expr
	OpPos token.Pos
	Op    token.Kind
	Y     Expr
}

// CastExpr is `x as T` (A.4.4). Left-associative, binds tighter than every
// binary operator, and never touches memory.
type CastExpr struct {
	X    Expr
	As   token.Pos
	Type Expr
}

// TransferExpr is the ownership marker: `var target` in an owning position
// (A.4.6, A.9.1).
//
// It is a node rather than a bool because A.9.1 lists six distinct owning
// positions, and one node covers all of them without six flags. Its presence is
// the entire difference between move and deep copy, so it must survive parsing
// as syntax and must never be normalized away.
//
// Target is grammatically an Identifier or a selector chain (TransferTarget).
// Anything else parses into this node and is rejected statically — A.14 lists
// "var on a computed expression" as a rejected form, which means it parses.
type TransferExpr struct {
	Var    token.Pos
	Target Expr
}

// ------------------------------------------------------------------ types

// OwnershipType is `mut T`, `var T`, `unique T`, `shared T`, `weak T` (A.3.2).
// Qualifiers do not stack, but a stacked form parses and is rejected (A.14).
type OwnershipType struct {
	KwPos token.Pos
	Kw    token.Kind // MUT, VAR, UNIQUE, SHARED, WEAK
	X     Expr
}

// ArrayType is `[N]T` and `[]T` (A.3.1). Len == nil is the slice form.
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

// FuncType is a FunctionType (A.3.4) and also the signature half of a FuncDecl
// and FuncLit. In a bare FunctionType every Param has a nil Name.
//
// Result carries `-> Type` or `-> ExpectedType`; the latter is a *CallExpr.
// Omitting it is the void form — A.3.4 has no `void` type name.
type FuncType struct {
	Func   token.Pos
	Params *ParamList
	Marker *Marker // nil, or at most one (A.6.1)
	Arrow  token.Pos
	Result Expr
}

type ChanType struct {
	Chan token.Pos
	Elem Expr
}

// PointerType is `typed_ptr T` (A.3.3). Nesting requires parentheses, so the
// nested form has Elem == *ParenExpr.
type PointerType struct {
	Kw   token.Pos
	Elem Expr
}

// TensorType is `tensor[T, dims...]` (A.3.5). Grammatical only under [+Npu];
// a tensor outside an npu body parses and is rejected (A.14).
type TensorType struct {
	Tensor token.Pos
	Lbrack token.Pos
	Elem   Expr
	Shape  []Expr // IntegerLiterals
	Rbrack token.Pos
}

// AbstractType is the bare `abstract` of A.3.3, legal only as an alias target.
type AbstractType struct {
	Abstract token.Pos
}

// BadExpr marks a span the parser could not make sense of. It exists so
// recovery produces a tree the analyzer can still walk; the analyzer skips it
// silently, because a diagnostic was already reported at parse time.
type BadExpr struct {
	From, To token.Pos
}

// -------------------------------------------------------------- positions

func (x *BasicLit) Pos() token.Pos      { return x.ValuePos }
func (x *NamespaceExpr) Pos() token.Pos { return x.KwPos }
func (x *ParenExpr) Pos() token.Pos     { return x.Lparen }
func (x *TupleExpr) Pos() token.Pos     { return x.Lparen }
func (x *ArrayLit) Pos() token.Pos      { return x.Lbrack }
func (x *CompositeLit) Pos() token.Pos  { return x.Type.Pos() }
func (x *MapLit) Pos() token.Pos        { return x.Lbrace }
func (x *KeyValueExpr) Pos() token.Pos  { return x.Key.Pos() }
func (x *EnumShorthand) Pos() token.Pos { return x.Dot }
func (x *FuncLit) Pos() token.Pos       { return x.Type.Pos() }
func (x *SelectorExpr) Pos() token.Pos  { return x.X.Pos() }
func (x *TupleIndexExpr) Pos() token.Pos { return x.X.Pos() }
func (x *IndexExpr) Pos() token.Pos     { return x.X.Pos() }
func (x *CallExpr) Pos() token.Pos      { return x.Fun.Pos() }
func (x *LaunchExpr) Pos() token.Pos    { return x.KwPos }
func (x *LaunchConfig) Pos() token.Pos  { return x.Lparen }
func (x *AwaitExpr) Pos() token.Pos     { return x.Await }
func (x *UnaryExpr) Pos() token.Pos     { return x.OpPos }
func (x *BinaryExpr) Pos() token.Pos    { return x.X.Pos() }
func (x *CastExpr) Pos() token.Pos      { return x.X.Pos() }
func (x *TransferExpr) Pos() token.Pos  { return x.Var }
func (x *OwnershipType) Pos() token.Pos { return x.KwPos }
func (x *ArrayType) Pos() token.Pos     { return x.Lbrack }
func (x *MapType) Pos() token.Pos       { return x.Map }
func (x *FuncType) Pos() token.Pos      { return x.Func }
func (x *ChanType) Pos() token.Pos      { return x.Chan }
func (x *PointerType) Pos() token.Pos   { return x.Kw }
func (x *TensorType) Pos() token.Pos    { return x.Tensor }
func (x *AbstractType) Pos() token.Pos  { return x.Abstract }
func (x *BadExpr) Pos() token.Pos       { return x.From }

func (x *BasicLit) End() token.Pos      { return x.ValuePos + token.Pos(len(x.Value)) }
func (x *NamespaceExpr) End() token.Pos { return x.KwPos + token.Pos(len(x.Kw.String())) }
func (x *ParenExpr) End() token.Pos     { return x.Rparen + 1 }
func (x *TupleExpr) End() token.Pos     { return x.Rparen + 1 }
func (x *ArrayLit) End() token.Pos      { return x.Rbrack + 1 }
func (x *CompositeLit) End() token.Pos  { return x.Rbrace + 1 }
func (x *MapLit) End() token.Pos        { return x.Rbrace + 1 }
func (x *KeyValueExpr) End() token.Pos  { return x.Value.End() }
func (x *FuncLit) End() token.Pos       { return x.Body.End() }
func (x *SelectorExpr) End() token.Pos  { return x.Sel.End() }
func (x *TupleIndexExpr) End() token.Pos { return x.IndexPos + token.Pos(len(x.Text)) }
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
func (x *AbstractType) End() token.Pos  { return x.Abstract + 8 }
func (x *BadExpr) End() token.Pos       { return x.To }

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
	if x.Marker != nil {
		return x.Marker.End()
	}
	return x.Params.End()
}

func (*BasicLit) exprNode()       {}
func (*NamespaceExpr) exprNode()  {}
func (*ParenExpr) exprNode()      {}
func (*TupleExpr) exprNode()      {}
func (*ArrayLit) exprNode()       {}
func (*CompositeLit) exprNode()   {}
func (*MapLit) exprNode()         {}
func (*KeyValueExpr) exprNode()   {}
func (*EnumShorthand) exprNode()  {}
func (*FuncLit) exprNode()        {}
func (*SelectorExpr) exprNode()   {}
func (*TupleIndexExpr) exprNode() {}
func (*IndexExpr) exprNode()      {}
func (*CallExpr) exprNode()       {}
func (*LaunchExpr) exprNode()     {}
func (*AwaitExpr) exprNode()      {}
func (*UnaryExpr) exprNode()      {}
func (*BinaryExpr) exprNode()     {}
func (*CastExpr) exprNode()       {}
func (*TransferExpr) exprNode()   {}
func (*OwnershipType) exprNode()  {}
func (*ArrayType) exprNode()      {}
func (*MapType) exprNode()        {}
func (*FuncType) exprNode()       {}
func (*ChanType) exprNode()       {}
func (*PointerType) exprNode()    {}
func (*TensorType) exprNode()     {}
func (*AbstractType) exprNode()   {}
func (*BadExpr) exprNode()        {}