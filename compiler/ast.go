package compiler

// ─────────────────────────────────────────────────────────────────────────────
// Source positions
// ─────────────────────────────────────────────────────────────────────────────

// Pos is a source location used for error reporting.
type Pos struct {
	File   string
	Line   int
	Column int
}

func (p Pos) IsValid() bool { return p.Line > 0 }

// ─────────────────────────────────────────────────────────────────────────────
// Node — root of the AST hierarchy
// ─────────────────────────────────────────────────────────────────────────────

type Node interface {
	nodePos() Pos
}

// ─────────────────────────────────────────────────────────────────────────────
// File
// ─────────────────────────────────────────────────────────────────────────────

type File struct {
	Pos     Pos
	Package string
	Imports []*ImportDecl
	Decls   []Decl
}

func (f *File) nodePos() Pos { return f.Pos }

type ImportDecl struct {
	Pos  Pos
	Path string // raw path string, e.g. "core/io"
}

func (d *ImportDecl) nodePos() Pos { return d.Pos }

// ─────────────────────────────────────────────────────────────────────────────
// Top-level declarations
// ─────────────────────────────────────────────────────────────────────────────

type Decl interface {
	Node
	declNode()
}

// FuncDecl covers both free functions and associated functions (receivers).
type FuncDecl struct {
	Pos        Pos
	Name       string
	Receiver   *Receiver  // nil for free functions
	TypeParams []string   // generic param names, e.g. ["T"]
	Params     []*Param
	Qualifier  FuncQual
	RetType    TypeExpr   // nil → void
	Body       *BlockStmt
}

type Receiver struct {
	Pos   Pos
	Name  string
	Type  TypeExpr
	IsPtr bool // *T receiver (computed by resolver)
}

type Param struct {
	Pos        Pos
	Name       string
	Type       TypeExpr
	IsVariadic bool // true for the trailing '...' parameter
}

type FuncQual int

const (
	FuncQualNone    FuncQual = iota
	FuncQualAsync            // deferred
	FuncQualThread           // deferred
	FuncQualProcess          // deferred
	FuncQualGPU              // deferred
	FuncQualTest             // test function: return → printf, gets its own main()
)

type StructDecl struct {
	Pos    Pos
	Name   string
	Fields []*FieldDecl
}

type FieldDecl struct {
	Pos  Pos
	Name string
	Type TypeExpr
}

type ClassDecl struct {
	Pos      Pos
	Name     string
	BaseName string // native-interface base; "" for regular classes
	Members  []*ClassMember
}

type ClassMember struct {
	Pos     Pos
	IsField bool
	Name    string
	Type    TypeExpr  // fields
	Params  []*Param  // method signature members
	RetType TypeExpr  // method signature members
}

type EnumDecl struct {
	Pos     Pos
	Name    string
	RawType TypeExpr // optional underlying type
	Cases   []*EnumCase
}

type EnumCase struct {
	Pos      Pos
	Name     string
	RawValue Expr // optional
}

type TypeAliasDecl struct {
	Pos  Pos
	Name string
	Type TypeExpr
}

type VarDecl struct {
	Pos      Pos
	IsLet    bool
	IsWeak   bool
	Binding  *BindingPattern
	TypeHint TypeExpr // optional explicit annotation
	Value    Expr
}

type BindingPattern struct {
	Pos   Pos
	Names []string // len==1 simple, >1 tuple destructure
}

// Decl marker implementations.
func (*FuncDecl) declNode()      {}
func (*StructDecl) declNode()    {}
func (*ClassDecl) declNode()     {}
func (*EnumDecl) declNode()      {}
func (*TypeAliasDecl) declNode() {}
func (*VarDecl) declNode()       {}

func (d *FuncDecl) nodePos() Pos      { return d.Pos }
func (d *StructDecl) nodePos() Pos    { return d.Pos }
func (d *ClassDecl) nodePos() Pos     { return d.Pos }
func (d *EnumDecl) nodePos() Pos      { return d.Pos }
func (d *TypeAliasDecl) nodePos() Pos { return d.Pos }
func (d *VarDecl) nodePos() Pos       { return d.Pos }

// ─────────────────────────────────────────────────────────────────────────────
// Type expressions (syntactic; resolved by Resolver into VType)
// ─────────────────────────────────────────────────────────────────────────────

type TypeExpr interface {
	Node
	typeExprNode()
}

type NamedTypeExpr struct {
	Pos  Pos
	Pkg  string // qualifier, e.g. "core" in "core.Context"; "" if none
	Name string // e.g. "int32", "Animal"
}

type PointerTypeExpr struct {
	Pos      Pos
	IsConst  bool
	Elem     TypeExpr
	Optional bool // *T?
}

type ArrayTypeExpr struct {
	Pos  Pos
	Elem TypeExpr // [T]
}

type OptionalTypeExpr struct {
	Pos  Pos
	Elem TypeExpr // T?
}

type FuncTypeExpr struct {
	Pos     Pos
	Params  []TypeExpr
	RetType TypeExpr
}

type TupleTypeExpr struct {
	Pos    Pos
	Elems  []TypeExpr
	Labels []string // parallel; "" for unlabelled
}

type ResultTypeExpr struct {
	Pos Pos
	Ok  TypeExpr
	Err TypeExpr
}

type ChanTypeExpr struct {
	Pos  Pos
	Elem TypeExpr
}

// ExpectedTypeExpr is the return-type annotation on a test-qualified function:
//
//	func test_add() test -> Expected(stdout, "15")
//
// Channel is the named output channel ("stdout", "exitCode").
// Value is the expected output string that the test runner will compare against.
type ExpectedTypeExpr struct {
	Pos     Pos
	Channel string // e.g. "stdout"
	Value   string // e.g. "15"
}

func (*NamedTypeExpr) typeExprNode()    {}
func (*PointerTypeExpr) typeExprNode()  {}
func (*ArrayTypeExpr) typeExprNode()    {}
func (*OptionalTypeExpr) typeExprNode() {}
func (*FuncTypeExpr) typeExprNode()     {}
func (*TupleTypeExpr) typeExprNode()    {}
func (*ResultTypeExpr) typeExprNode()   {}
func (*ChanTypeExpr) typeExprNode()     {}
func (*ExpectedTypeExpr) typeExprNode() {}

func (t *NamedTypeExpr) nodePos() Pos    { return t.Pos }
func (t *PointerTypeExpr) nodePos() Pos  { return t.Pos }
func (t *ArrayTypeExpr) nodePos() Pos    { return t.Pos }
func (t *OptionalTypeExpr) nodePos() Pos { return t.Pos }
func (t *FuncTypeExpr) nodePos() Pos     { return t.Pos }
func (t *TupleTypeExpr) nodePos() Pos    { return t.Pos }
func (t *ResultTypeExpr) nodePos() Pos   { return t.Pos }
func (t *ChanTypeExpr) nodePos() Pos     { return t.Pos }
func (t *ExpectedTypeExpr) nodePos() Pos { return t.Pos }

// ─────────────────────────────────────────────────────────────────────────────
// Statements
// ─────────────────────────────────────────────────────────────────────────────

type Stmt interface {
	Node
	stmtNode()
}

type BlockStmt struct {
	Pos   Pos
	Stmts []Stmt
}

type LocalDeclStmt struct {
	Pos  Pos
	Decl *VarDecl
}

type IfStmt struct {
	Pos  Pos
	Cond IfCond
	Then *BlockStmt
	Else Stmt // *BlockStmt | *IfStmt | nil
}

type IfCond interface {
	Node
	ifCondNode()
}

type IfLetCond struct {
	Pos  Pos
	Name string
	Expr Expr
}

type IfExprCond struct {
	Pos  Pos
	Expr Expr
}

type WhileStmt struct {
	Pos  Pos
	Cond Expr
	Body *BlockStmt
}

type ForInStmt struct {
	Pos  Pos
	Var  string
	Iter Expr
	Body *BlockStmt
}

type SwitchStmt struct {
	Pos   Pos
	Subj  Expr
	Cases []*SwitchCase
}

type SwitchCase struct {
	Pos       Pos
	IsDefault bool
	Patterns  []SwitchPattern
	Body      []Stmt
}

type SwitchPattern interface {
	Node
	switchPatternNode()
}

type ExprPattern struct {
	Pos  Pos
	Expr Expr
}

type EnumShortPattern struct {
	Pos  Pos
	Case string // e.g. "north" from ".north"
}

type ResultOkPattern struct {
	Pos  Pos
	Bind string
}

type ResultErrPattern struct {
	Pos  Pos
	Bind string
}

type ReturnStmt struct {
	Pos   Pos
	Value Expr // nil for void return
}

type DeferStmt struct {
	Pos  Pos
	Call Expr
}

type BreakStmt struct{ Pos Pos }
type ContinueStmt struct{ Pos Pos }
type FallthroughStmt struct{ Pos Pos }

type AssignStmt struct {
	Pos Pos
	LHS Expr
	Op  AssignOp
	RHS Expr
}

type AssignOp int

const (
	OpAssign    AssignOp = iota
	OpAddAssign          // +=
	OpSubAssign          // -=
	OpMulAssign          // *=
	OpDivAssign          // /=
	OpModAssign          // %=
)

type ExprStmt struct {
	Pos  Pos
	Expr Expr
}

// Stmt markers.
func (*BlockStmt) stmtNode()       {}
func (*LocalDeclStmt) stmtNode()   {}
func (*IfStmt) stmtNode()          {}
func (*WhileStmt) stmtNode()       {}
func (*ForInStmt) stmtNode()       {}
func (*SwitchStmt) stmtNode()      {}
func (*ReturnStmt) stmtNode()      {}
func (*DeferStmt) stmtNode()       {}
func (*BreakStmt) stmtNode()       {}
func (*ContinueStmt) stmtNode()    {}
func (*FallthroughStmt) stmtNode() {}
func (*AssignStmt) stmtNode()      {}
func (*ExprStmt) stmtNode()        {}

func (s *BlockStmt) nodePos() Pos       { return s.Pos }
func (s *LocalDeclStmt) nodePos() Pos   { return s.Pos }
func (s *IfStmt) nodePos() Pos          { return s.Pos }
func (s *WhileStmt) nodePos() Pos       { return s.Pos }
func (s *ForInStmt) nodePos() Pos       { return s.Pos }
func (s *SwitchStmt) nodePos() Pos      { return s.Pos }
func (s *ReturnStmt) nodePos() Pos      { return s.Pos }
func (s *DeferStmt) nodePos() Pos       { return s.Pos }
func (s *BreakStmt) nodePos() Pos       { return s.Pos }
func (s *ContinueStmt) nodePos() Pos    { return s.Pos }
func (s *FallthroughStmt) nodePos() Pos { return s.Pos }
func (s *AssignStmt) nodePos() Pos      { return s.Pos }
func (s *ExprStmt) nodePos() Pos        { return s.Pos }

func (*IfLetCond) ifCondNode()  {}
func (*IfExprCond) ifCondNode() {}

func (c *IfLetCond) nodePos() Pos  { return c.Pos }
func (c *IfExprCond) nodePos() Pos { return c.Pos }

func (*ExprPattern) switchPatternNode()      {}
func (*EnumShortPattern) switchPatternNode() {}
func (*ResultOkPattern) switchPatternNode()  {}
func (*ResultErrPattern) switchPatternNode() {}

func (p *ExprPattern) nodePos() Pos      { return p.Pos }
func (p *EnumShortPattern) nodePos() Pos { return p.Pos }
func (p *ResultOkPattern) nodePos() Pos  { return p.Pos }
func (p *ResultErrPattern) nodePos() Pos { return p.Pos }

// ─────────────────────────────────────────────────────────────────────────────
// Expressions
// ─────────────────────────────────────────────────────────────────────────────

// Expr is the expression interface. The vtype field is populated by the Resolver.
type Expr interface {
	Node
	exprNode()
	GetVType() VType
	SetVType(VType)
}

// exprBase is embedded by every expression node.
type exprBase struct {
	Pos   Pos
	vtype VType
}

func (e *exprBase) nodePos() Pos     { return e.Pos }
func (e *exprBase) exprNode()        {}
func (e *exprBase) GetVType() VType  { return e.vtype }
func (e *exprBase) SetVType(t VType) { e.vtype = t }

// ── Literals ──────────────────────────────────────────────────────────────────

type IntLitExpr struct {
	exprBase
	Value      int64
	IsUnsigned bool
}

type FloatLitExpr struct {
	exprBase
	Value   float64
	Is32Bit bool
}

// ReinterpretExpr is reinterpret<T>(expr) — zero-cost raw pointer reinterpretation.
type ReinterpretExpr struct {
	exprBase
	TargetType TypeExpr
	Value      Expr
}

type BoolLitExpr struct {
	exprBase
	Value bool
}

type StringLitExpr struct {
	exprBase
	Value string // unescaped content (backslash sequences already processed)
}

type NilLitExpr struct{ exprBase }

// ── Primary expressions ───────────────────────────────────────────────────────

type IdentExpr struct {
	exprBase
	Name string
}

type DotEnumExpr struct {
	exprBase
	Case string // e.g. "north" from ".north"
}

type StructLitExpr struct {
	exprBase
	TypeName string
	Fields   []*StructLitField
}

type StructLitField struct {
	Pos   Pos
	Name  string
	Value Expr
}

// ArrayLitExpr is a literal array: [a, b, c].
type ArrayLitExpr struct {
	exprBase
	Elems []Expr
}

// ArrayCtorExpr is the [T]() / [T](capacity:N) / [T](N) family.
type ArrayCtorExpr struct {
	exprBase
	ElemTypeName string
	Args         []*Arg
}

type TupleLitExpr struct {
	exprBase
	Elems []Expr
}

type ResultExpr struct {
	exprBase
	IsOk  bool // true → Ok, false → Err
	Value Expr
}

// ── Operators ─────────────────────────────────────────────────────────────────

type BinaryExpr struct {
	exprBase
	Op    BinOp
	Left  Expr
	Right Expr
}

type BinOp int

const (
	BinAdd           BinOp = iota
	BinSub
	BinMul
	BinDiv
	BinMod
	BinShl
	BinShr
	BinEq
	BinNeq
	BinLt
	BinLte
	BinGt
	BinGte
	BinAnd           // &&
	BinOr            // ||
	BinNilCoalesce   // ??
	BinOverflowAdd   // &+
	BinOverflowSub   // &-
	BinOverflowMul   // &*
	BinRangeHalfOpen // ..
	BinRangeClosed   // ...
	BinIdentityEq    // ===
	BinIdentityNeq   // !==
)

type UnaryExpr struct {
	exprBase
	Op      UnOp
	Operand Expr
}

type UnOp int

const (
	UnNeg    UnOp = iota // -
	UnNot                // !
	UnBitNot             // ~
	UnAddrOf             // &
)

type TernaryExpr struct {
	exprBase
	Cond Expr
	Then Expr
	Else Expr
}

// ── Calls and access ──────────────────────────────────────────────────────────

type CallExpr struct {
	exprBase
	Func Expr
	Args []*Arg
}

type MethodCallExpr struct {
	exprBase
	Recv   Expr
	Method string
	Args   []*Arg
}

type Arg struct {
	Pos   Pos
	Label string // "" if positional
	Value Expr
}

type FieldExpr struct {
	exprBase
	Recv  Expr
	Field string
}

type IndexExpr struct {
	exprBase
	Recv  Expr
	Index Expr
}

// TypeConvExpr is a type conversion: float(x), int8(y), etc.
type TypeConvExpr struct {
	exprBase
	TargetType TypeExpr
	Value      Expr
}

// AnonFuncExpr is an anonymous function value.
type AnonFuncExpr struct {
	exprBase
	Params    []*Param
	Qualifier FuncQual
	RetType   TypeExpr
	Body      *BlockStmt
}

// MapLitExpr represents a map literal: {"key": value, ...}
type MapLitExpr struct {
	exprBase
	Fields []*MapLitField
}

type MapLitField struct {
	Pos   Pos
	Key   Expr
	Value Expr
}