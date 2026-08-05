// Package hir is where every decision is made.
//
// Given the checked package graph and its *types.Info, hir produces a
// *Program: monomorphic, ownership-explicit, control-flow-flattened, with
// every builtin call named. lower/vir is mechanical afterward by
// construction — if it ever needs a type switch on "is this owning" or "was
// this transferred", that logic belongs here.
//
// hir does not import vvm and never sees a target triple. Layout is the one
// target-shaped fact it consumes, and it arrives as a *types.Sizes on
// Config. It imports builtins for names.go's ABI constants and features.go's
// FeatureSet only — never a module constructor, so what decides *which* call
// to emit cannot see *how* the callee is built.
//
// Three representation rules run through everything below:
//
//  1. Aggregates are pointers. vir makes struct and array memory-only, so a
//     Value whose Type is aggregate *is* a ptr at the vir level. Aggregate
//     parameters carry ByVal, aggregate results carry SRet, aggregate
//     assignment is a memcopy. vector and the lane predicate are the
//     exceptions: register-class types passed by value.
//  2. `let` is a value, `var` is a slot. semantics.md §6.1 already says a
//     `let` may be a register, an SSA value, or folded away entirely.
//     Making that literal sidesteps every Join-Convention definite-
//     assignment subtlety around mutation across branches, and it is why
//     only a `var` can reach a `mut` parameter. See forcesSlot.
//  3. One instruction representation, two control-flow shapes. A Func
//     carries structured Body until Flatten runs and flat Blocks afterward,
//     but both hold the same *Instr values. Flatten moves instructions; it
//     never rewrites them.
package hir

import (
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/token"
)

// Program is one whole-program lowering. Monomorphization is seeded from
// main (or the one test function under ModeTest) and walks the call graph,
// so hir consumes the entire checked graph at once even though lower/vir
// emits one module per Vertex package.
type Program struct {
	Modules []*Module // topological order; dependencies first
	Root    *Module   // holds the entry shim
	Entry   *Func     // the entry shim itself

	// Features is this package's answer to "where is the feature set
	// computed". Every emitted builtin call records its feature at the call
	// site, so the set can never disagree with the calls actually emitted.
	Features builtins.FeatureSet
}

// Module is one Vertex package. Name becomes vir's `module`, Path its
// `namespace` — semantics.md §1.3's "the path is a locator, the declared
// name is the qualifier", carried through to the linker unchanged.
type Module struct {
	Name string
	Path string

	Structs []*Struct
	Consts  []*Const
	Globals []*Global
	Links   []*Link
	Externs []*ExternGroup
	Imports []string // vir import paths, in first-reference order
	Funcs   []*Func

	imports map[string]bool
	structs map[string]*Struct // by hir name, for interning
}

func newModule(name, path string) *Module {
	return &Module{
		Name:    name,
		Path:    path,
		imports: map[string]bool{},
		structs: map[string]*Struct{},
	}
}

// Import records a cross-module dependency. Idempotent, and order-stable so
// a build is byte-reproducible.
func (m *Module) Import(path string) {
	if path == "" || m.imports[path] {
		return
	}
	m.imports[path] = true
	m.Imports = append(m.Imports, path)
}

// Struct is a lowered aggregate: a Vertex struct or class, a tuple, a
// boundary tuple, a synthesized header (_Vstr, _Vslice, _Vvec), or an async
// frame. Fields are in declaration order with no reordering — semantics.md
// §7.2 makes order observable through destruction order, and interop
// assumes it.
type Struct struct {
	Name   string
	Module *Module
	Fields []StructField
	Export bool

	Size  int64
	Align int64
}

type StructField struct {
	Name   string
	Type   *Type
	Offset int64
}

func (s *Struct) FieldIndex(name string) int {
	for i, f := range s.Fields {
		if f.Name == name {
			return i
		}
	}
	return -1
}

// Const is a compile-time scalar. Vertex's top-level `let` folds to one when
// its type is scalar; an aggregate one becomes a Global instead, since vir
// consts are scalars only.
type Const struct {
	Name   string
	Type   *Type
	Value  Value
	Export bool
}

// Global is module-level storage: a top-level `var`, a string literal's
// bytes, or a lazily-filled selector cache.
//
// vir's global init form is narrower than semantics.md §5.3's constant
// expressions — literal, zero, addr, or an aggregate of those, with no
// arithmetic and no const references — so decl.go folds initializers down to
// exactly those forms. There is no static-initialization-order problem
// because there is no initialization-time code.
type Global struct {
	Name   string
	Type   *Type
	Init   ConstInit
	Export bool
	TLS    bool
	Align  int
}

type ConstInit interface{ constInit() }

type (
	InitZero      struct{}
	InitScalar    struct{ Value Value }
	InitBytes     struct{ Data []byte }
	InitAddrOf    struct{ Name string }
	InitAggregate struct{ Elems []ConstInit }
)

func (InitZero) constInit()      {}
func (InitScalar) constInit()    {}
func (InitBytes) constInit()     {}
func (InitAddrOf) constInit()    {}
func (InitAggregate) constInit() {}

// Link and ExternGroup exist only for declare blocks. A module lowered from
// ordinary Vertex source carries neither — the invariant is checkable by
// reading the emitted text.
type Link struct {
	Kind string // "static" | "shared" | "framework"
	Name string
}

type ExternGroup struct {
	Dependency string
	Functions  []*ExternFunc
}

type ExternFunc struct {
	Name     string
	Params   []*Param
	Result   *Type
	Variadic bool
	NoReturn bool
}

// Func is one lowered function: a monomorphic instance, a method, a
// synthesized _Vcopy/_Vdrop routine, an epilogue-expanded body, or the entry
// shim.
type Func struct {
	Name   string
	Module *Module
	Pos    token.Pos

	Params []*Param
	Result *Type // nil for void; aggregates go through SRet instead
	SRet   *Type // non-nil when the result is an aggregate

	Export   bool
	Entry    bool
	NoReturn bool
	Variadic bool // foreign C variadics only; Vertex variadics are slices

	// Body is the structured tree; Blocks is what Flatten produces from it.
	// Exactly one is non-nil after lowering completes.
	Body   *Seq
	Blocks []*FlatBlock

	// Allocas are hoisted to the entry block. vir allocas are per-execution
	// and accumulate per loop iteration, so a slot written inside a loop is
	// allocated once, before the loop, and reused — sound because a
	// loop-body local's teardown runs on every back edge.
	Allocas []*Instr

	names int // counter for generated value names
}

type Param struct {
	Name  string
	Type  *Type
	ByVal *Struct // set when the parameter is an owning aggregate
	SRet  *Struct // set on the synthetic first parameter of an sret function
}

// fresh mints an unused value name. vir idents are [A-Za-z_][A-Za-z0-9_]*,
// and the Join Convention keys on the name, so uniqueness within a function
// is the whole requirement.
func (f *Func) fresh(hint string) string {
	f.names++
	if hint == "" {
		hint = "t"
	}
	return hint + itoa(f.names)
}

// ---------------------------------------------------------------- statements

// Seq is a run of structured statements. Stmt is closed: everything a Vertex
// body can do reaches Flatten as one of these.
type Seq struct{ List []Stmt }

func (s *Seq) add(x Stmt) {
	if x != nil {
		s.List = append(s.List, x)
	}
}

type Stmt interface{ stmtNode() }

type (
	// If is the only conditional shape. `&&`/`||` are already branches by
	// the time they get here — semantics.md §5.1 has no truthiness, so the
	// operand is a real bool.
	If struct {
		Cond Value
		Then *Seq
		Else *Seq // may be nil
	}

	// Loop is unconditional; every Vertex loop shape (while, for over a
	// range, an array, a map) is lowered to a Loop whose head tests and
	// Breaks. There are no loop labels, so Break/Continue always name the
	// innermost one.
	Loop struct {
		Body *Seq
	}

	// SwitchStmt dispatches on an integer or an enum tag. Cases are dense or
	// sparse; jump-table-vs-compare-chain is cpu/lower's decision, not this
	// package's, so both spell one vir switch.
	SwitchStmt struct {
		Value   Value
		Cases   []SwitchCase
		Default *Seq // never nil; an exhaustive enum switch gets Unreachable
	}

	Break    struct{}
	Continue struct{}

	// ReturnStmt carries the value for a thin result. An aggregate result
	// was already stored through the SRet pointer, so Value is nil there.
	ReturnStmt struct{ Value *Value }

	TrapStmt        struct{}
	UnreachableStmt struct{}
)

type SwitchCase struct {
	Value int64
	Body  *Seq
}

func (*Instr) stmtNode()           {}
func (*If) stmtNode()              {}
func (*Loop) stmtNode()            {}
func (*SwitchStmt) stmtNode()      {}
func (*Break) stmtNode()           {}
func (*Continue) stmtNode()        {}
func (*ReturnStmt) stmtNode()      {}
func (*TrapStmt) stmtNode()        {}
func (*UnreachableStmt) stmtNode() {}

// --------------------------------------------------------------- instructions

// Instr is one instruction. It is a Stmt too, which is what lets Flatten
// move instructions between shapes without rewriting them.
//
// Type is the *result* type. lower/vir derives the vir suffix from it,
// except for the comparison family, where the suffix names the operand type
// and is read off Args[0].
type Instr struct {
	Result string
	Op     Op
	Type   *Type
	Args   []Value
	Align  int
	Pos    token.Pos

	// Callee is set for OpCall. Module is the owning module name for a
	// cross-module call — a Vertex package or a builtins module — and "" for
	// a call within this module. Sig names a fnsig for an indirect call and
	// is not yet produced by anything.
	Callee string
	Module string
	Sig    string
}

// Value is an operand.
type Value struct {
	Kind ValueKind
	Name string  // VName, VGlobal, VFuncAddr
	Int  int64   // VInt
	Flt  float64 // VFloat
	Type *Type
}

type ValueKind uint8

const (
	VNone ValueKind = iota
	VName
	VInt
	VFloat
	VNull
	VGlobal
	VFuncAddr // `addr f` — hir's only spelling for a function address
	VType     // a type in operand position, for index.ptr
)

func Name(n string, t *Type) Value    { return Value{Kind: VName, Name: n, Type: t} }
func Int(v int64, t *Type) Value      { return Value{Kind: VInt, Int: v, Type: t} }
func Float(v float64, t *Type) Value  { return Value{Kind: VFloat, Flt: v, Type: t} }
func Null() Value                     { return Value{Kind: VNull, Type: Ptr} }
func GlobalRef(n string, t *Type) Value { return Value{Kind: VGlobal, Name: n, Type: t} }
func FuncAddr(n string) Value         { return Value{Kind: VFuncAddr, Name: n, Type: Ptr} }
func TypeVal(t *Type) Value           { return Value{Kind: VType, Type: t} }

func (v Value) IsZero() bool { return v.Kind == VNone }

// ---------------------------------------------------------- flat control flow

// FlatBlock is what Flatten produces. The entry block is Blocks[0] and its
// Label is "" — vir's entry block is implicit, unlabeled, and unbranchable-to.
type FlatBlock struct {
	Label string
	Lines []*Instr
	Term  Term
}

type Term interface{ termNode() }

type (
	Br   struct{ Label string }
	BrIf struct {
		Cond       Value
		Then, Else string
	}
	SwitchTerm struct {
		Value   Value
		Default string
		Cases   []SwitchTermCase
	}
	Ret         struct{ Value *Value }
	TrapTerm    struct{}
	UnreachTerm struct{}
)

type SwitchTermCase struct {
	Value int64
	Label string
}

func (Br) termNode()          {}
func (BrIf) termNode()        {}
func (SwitchTerm) termNode()  {}
func (Ret) termNode()         {}
func (TrapTerm) termNode()    {}
func (UnreachTerm) termNode() {}

// itoa avoids pulling strconv into every file for one use.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}