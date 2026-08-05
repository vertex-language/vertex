package hir

// Op is hir's instruction vocabulary: a subset of vir's §4 opcodes, spelled
// as hir spells them.
//
// Two spellings deliberately differ from vir's, and that difference is why
// lower/vir carries an explicit table rather than bridging by name:
// hir writes `field.ptr` and `index.ptr` — the instruction form, suffix
// included — where vir's closed vocabulary writes `field` and `index`. A
// name-based bridge works for sixty opcodes and fails silently for two.
type Op uint16

const (
	OpInvalid Op = iota

	// Arithmetic. Every one of these is emitted with its overflow check
	// alongside, except where the source wrote a wrapping operator — vir
	// wraps and Vertex traps, and expr.go is where that gap is paid for.
	OpAdd
	OpSub
	OpMul
	OpNeg
	OpSDiv
	OpUDiv
	OpSRem
	OpURem

	// Overflow predicates: yield i1, not a result.
	OpSAddO
	OpUAddO
	OpSSubO
	OpUSubO
	OpSMulO
	OpUMulO

	// Bitwise and shifts.
	OpAnd
	OpOr
	OpXor
	OpNot
	OpShl
	OpLShr
	OpAShr

	// Comparisons. The vir suffix names the *operand* type, not the result,
	// so lower/vir reads it off Args[0] for this family.
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
	OpFlt
	OpFgt
	OpFle
	OpFge

	OpSelect

	// Memory.
	OpAlloca
	OpLoad
	OpStore
	OpMemcopy
	OpMemmove
	OpMemset
	OpFieldPtr // vir: field
	OpIndexPtr // vir: index

	// Atomics — shared/weak refcounts are inline, not a runtime call.
	OpAtomicLoad
	OpAtomicStore
	OpAtomicAdd
	OpAtomicSub
	OpCmpxchg

	// Conversions. Destination-explicit: Instr.Type is the destination.
	OpTrunc
	OpSext
	OpZext
	OpFdemote
	OpFpromote
	OpBitcast
	OpSfromint
	OpUfromint
	OpStoint
	OpUtoint

	// Calls.
	OpCall

	opCount
)

var opNames = [opCount]string{
	OpAdd: "add", OpSub: "sub", OpMul: "mul", OpNeg: "neg",
	OpSDiv: "sdiv", OpUDiv: "udiv", OpSRem: "srem", OpURem: "urem",

	OpSAddO: "saddo", OpUAddO: "uaddo", OpSSubO: "ssubo", OpUSubO: "usubo",
	OpSMulO: "smulo", OpUMulO: "umulo",

	OpAnd: "and", OpOr: "or", OpXor: "xor", OpNot: "not",
	OpShl: "shl", OpLShr: "lshr", OpAShr: "ashr",

	OpEq: "eq", OpNe: "ne",
	OpSlt: "slt", OpSgt: "sgt", OpSle: "sle", OpSge: "sge",
	OpUlt: "ult", OpUgt: "ugt", OpUle: "ule", OpUge: "uge",
	OpFlt: "lt", OpFgt: "gt", OpFle: "le", OpFge: "ge",

	OpSelect: "select",

	OpAlloca: "alloca", OpLoad: "load", OpStore: "store",
	OpMemcopy: "memcopy", OpMemmove: "memmove", OpMemset: "memset",
	OpFieldPtr: "field.ptr", OpIndexPtr: "index.ptr",

	OpAtomicLoad: "atomic_load", OpAtomicStore: "atomic_store",
	OpAtomicAdd: "atomic_add", OpAtomicSub: "atomic_sub", OpCmpxchg: "cmpxchg",

	OpTrunc: "trunc", OpSext: "sext", OpZext: "zext",
	OpFdemote: "fdemote", OpFpromote: "fpromote", OpBitcast: "bitcast",
	OpSfromint: "sfromint", OpUfromint: "ufromint",
	OpStoint: "stoint", OpUtoint: "utoint",

	OpCall: "call",
}

func (o Op) String() string {
	if int(o) > 0 && int(o) < len(opNames) && opNames[o] != "" {
		return opNames[o]
	}
	return "<op?>"
}

// IsComparison reports whether o's vir suffix names its operand type rather
// than its result type.
func (o Op) IsComparison() bool {
	switch o {
	case OpEq, OpNe, OpSlt, OpSgt, OpSle, OpSge,
		OpUlt, OpUgt, OpUle, OpUge, OpFlt, OpFgt, OpFle, OpFge:
		return true
	}
	return false
}

// IsVoid reports whether o produces no value, and therefore takes no vir
// type suffix.
func (o Op) IsVoid() bool {
	switch o {
	case OpStore, OpMemcopy, OpMemmove, OpMemset, OpAtomicStore:
		return true
	}
	return false
}