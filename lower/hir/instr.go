package hir

import "github.com/vertex-language/vertex/token"

// Op is hir's instruction vocabulary. It is deliberately a subset of vir's
// §4 opcodes, spelled identically, so lower/vir's mapping is a table rather
// than a translation. Anything vir can express that Vertex source cannot
// reach (atomics, vectors, syscalls, varargs cursors) is absent: builtins
// reach those through hand-built vir modules, not through this enum.
type Op uint8

const (
	OpInvalid Op = iota

	OpAdd
	OpSub
	OpMul
	OpUDiv
	OpSDiv
	OpURem
	OpSRem
	OpNeg
	OpAbs

	OpAnd
	OpOr
	OpXor
	OpNot
	OpShl
	OpLShr
	OpAShr

	OpEq
	OpNe
	OpSlt
	OpSgt
	OpSle
	OpSge
	OpUlt
	OpUgt
	OpUle
	OpUge
	OpLt // float comparisons
	OpGt
	OpLe
	OpGe

	OpSelect

	OpAlloca
	OpLoad
	OpStore
	OpMemcopy
	OpMemmove
	OpMemset
	OpFieldPtr
	OpIndexPtr

	OpTrunc
	OpSext
	OpZext
	OpFdemote
	OpFpromote
	OpBitcast
	OpSfromint
	OpUfromint
	OpStointSat
	OpUtointSat

	OpSMin
	OpSMax
	OpUMin
	OpUMax

	OpCall
)

var opNames = map[Op]string{
	OpAdd: "add", OpSub: "sub", OpMul: "mul", OpUDiv: "udiv", OpSDiv: "sdiv",
	OpURem: "urem", OpSRem: "srem", OpNeg: "neg", OpAbs: "abs",
	OpAnd: "and", OpOr: "or", OpXor: "xor", OpNot: "not",
	OpShl: "shl", OpLShr: "lshr", OpAShr: "ashr",
	OpEq: "eq", OpNe: "ne", OpSlt: "slt", OpSgt: "sgt", OpSle: "sle", OpSge: "sge",
	OpUlt: "ult", OpUgt: "ugt", OpUle: "ule", OpUge: "uge",
	OpLt: "lt", OpGt: "gt", OpLe: "le", OpGe: "ge",
	OpSelect: "select",
	OpAlloca: "alloca", OpLoad: "load", OpStore: "store",
	OpMemcopy: "memcopy", OpMemmove: "memmove", OpMemset: "memset",
	OpFieldPtr: "field.ptr", OpIndexPtr: "index.ptr",
	OpTrunc: "trunc", OpSext: "sext", OpZext: "zext",
	OpFdemote: "fdemote", OpFpromote: "fpromote", OpBitcast: "bitcast",
	OpSfromint: "sfromint", OpUfromint: "ufromint",
	OpStointSat: "stoint_sat", OpUtointSat: "utoint_sat",
	OpSMin: "smin", OpSMax: "smax", OpUMin: "umin", OpUMax: "umax",
	OpCall: "call",
}

func (o Op) String() string {
	if s, ok := opNames[o]; ok {
		return s
	}
	return "<invalid op>"
}

// ValueKind discriminates an operand.
type ValueKind uint8

const (
	ValInvalid ValueKind = iota
	ValRef               // a named value or a binding
	ValInt
	ValFloat
	ValBool
	ValNull
	ValGlobal // the address of a module-level global
	ValConst  // a named compile-time constant
)

// Value is one operand. Type is the *logical* hir type: for an aggregate it
// is the struct or array, while the operand itself is a pointer to storage.
type Value struct {
	Kind   ValueKind
	Name   string
	Module string // qualifier for a cross-module const/global reference
	Int    int64
	Float  float64
	Bool   bool
	Type   Type
}

func Ref(name string, t Type) Value     { return Value{Kind: ValRef, Name: name, Type: t} }
func IntVal(t Type, v int64) Value      { return Value{Kind: ValInt, Int: v, Type: t} }
func FloatVal(t Type, v float64) Value  { return Value{Kind: ValFloat, Float: v, Type: t} }
func BoolVal(v bool) Value              { return Value{Kind: ValBool, Bool: v, Type: I1} }
func NullVal() Value                    { return Value{Kind: ValNull, Type: Ptr} }
func GlobalVal(name string, t Type) Value {
	return Value{Kind: ValGlobal, Name: name, Type: t}
}

func (v Value) Valid() bool { return v.Kind != ValInvalid }

// Callee names a call target. Exactly one shape applies: a local function, a
// qualified cross-module function (which lower/vir renders as
// `module.symbol` and vvm's importer rewrites away before codegen), or an
// indirect call through a function pointer against a named signature.
type Callee struct {
	Module   string // "" for same-module or extern
	Name     string
	Indirect *Value
	Sig      string // signature name, indirect calls only
}

// FieldRef names the struct and field a field.ptr computes into.
type FieldRef struct {
	Struct *Struct
	Field  string
}

// Instr is one instruction. Name binds the result under the Join
// Convention: a name's first assignment fixes its type permanently, and
// values merge across blocks by same-name assignment — there are no phi
// nodes on either side of the tree.
type Instr struct {
	Pos  token.Pos
	Name string // "" when the instruction produces no value
	Op   Op
	Type Type // result type / instruction type suffix
	Args []Value

	Call  *Callee   // OpCall
	Field *FieldRef // OpFieldPtr
	Elem  Type      // OpIndexPtr: the element type the arithmetic strides by
	Align int
}

// Result returns the Value this instruction binds, or an invalid Value.
func (i *Instr) Result() Value {
	if i.Name == "" {
		return Value{}
	}
	return Ref(i.Name, i.Type)
}

// ---------------------------------------------------------------------------
// Structured form
// ---------------------------------------------------------------------------

// Stmt is the structured control-flow shape a Func carries until Flatten
// runs. Structured statements are what let the async split see scopes and
// what make defer/deinit epilogue expansion CFG surgery on a tree rather
// than on a graph (overview §3, pass order).
type Stmt interface{ stmtNode() }

// Seq is an ordered run of statements.
type Seq struct{ List []Stmt }

// Instrs is a straight-line run. There is exactly one instruction
// representation across the whole pipeline: Flatten moves these into
// Blocks, it does not rewrite them.
type Instrs struct{ List []*Instr }

type If struct {
	Cond Value
	Then *Seq
	Else *Seq // may be nil
}

// Loop is the only loop shape, matching A.5.5's "while is the only loop
// primitive." The condition lives at the head of Body as an If whose Then
// is a Break, so Continue re-evaluates it for free.
type Loop struct{ Body *Seq }

type SwitchCase struct {
	Values []int64
	Body   *Seq
}

type Switch struct {
	Tag     Value
	Cases   []SwitchCase
	Default *Seq // may be nil
}

type Break struct{}
type Continue struct{}
type Return struct{ Value *Value }
type Trap struct{}
type Unreachable struct{}

func (*Seq) stmtNode()         {}
func (*Instrs) stmtNode()      {}
func (*If) stmtNode()          {}
func (*Loop) stmtNode()        {}
func (*Switch) stmtNode()      {}
func (*Break) stmtNode()       {}
func (*Continue) stmtNode()    {}
func (*Return) stmtNode()      {}
func (*Trap) stmtNode()        {}
func (*Unreachable) stmtNode() {}

func (s *Seq) add(x Stmt) { s.List = append(s.List, x) }

// ---------------------------------------------------------------------------
// Flat form
// ---------------------------------------------------------------------------

// Block is one labeled sequence ending in exactly one terminator. The entry
// block carries an empty Label: vir's entry block is implicit, unlabeled,
// and unbranchable-to.
type Block struct {
	Label string
	Instr []*Instr
	Term  Terminator
}

type Terminator interface{ termNode() }

type TermBranch struct{ Label string }
type TermBranchIf struct {
	Cond       Value
	Then, Else string
}
type TermCase struct {
	Value int64
	Label string
}
type TermSwitch struct {
	Value   Value
	Default string
	Cases   []TermCase
}
type TermReturn struct{ Value *Value }
type TermTrap struct{}
type TermUnreachable struct{}

func (TermBranch) termNode()      {}
func (TermBranchIf) termNode()    {}
func (TermSwitch) termNode()      {}
func (TermReturn) termNode()      {}
func (TermTrap) termNode()        {}
func (TermUnreachable) termNode() {}