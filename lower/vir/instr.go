// instr.go
package vir

import (
	"github.com/vertex-language/vertex/lower/hir"
	vir "github.com/vertex-language/vvm/ir/vir"
)

// The opcode table.
//
// hir.Op is a subset of vir's §4 opcodes spelled identically, so this is a
// table. It is an explicit table rather than vir.ParseOpcode(op.String())
// because hir spells its pointer arithmetic field.ptr / index.ptr — the
// instruction form, suffix included — while vir's closed vocabulary spells
// them field / index. A name-based bridge works for sixty opcodes and fails
// silently for two.
var opTable = map[hir.Op]vir.Opcode{
	hir.OpAdd:  vir.OpAdd,
	hir.OpSub:  vir.OpSub,
	hir.OpMul:  vir.OpMul,
	hir.OpNeg:  vir.OpNeg,
	hir.OpSDiv: vir.OpSDiv,
	hir.OpUDiv: vir.OpUDiv,
	hir.OpSRem: vir.OpSRem,
	hir.OpURem: vir.OpURem,

	hir.OpSAddO: vir.OpSAddO,
	hir.OpUAddO: vir.OpUAddO,
	hir.OpSSubO: vir.OpSSubO,
	hir.OpUSubO: vir.OpUSubO,
	hir.OpSMulO: vir.OpSMulO,
	hir.OpUMulO: vir.OpUMulO,

	hir.OpAnd:  vir.OpAnd,
	hir.OpOr:   vir.OpOr,
	hir.OpXor:  vir.OpXor,
	hir.OpNot:  vir.OpNot,
	hir.OpShl:  vir.OpShl,
	hir.OpLShr: vir.OpLShr,
	hir.OpAShr: vir.OpAShr,

	hir.OpEq:  vir.OpEq,
	hir.OpNe:  vir.OpNe,
	hir.OpSlt: vir.OpSlt,
	hir.OpSgt: vir.OpSgt,
	hir.OpSle: vir.OpSle,
	hir.OpSge: vir.OpSge,
	hir.OpUlt: vir.OpUlt,
	hir.OpUgt: vir.OpUgt,
	hir.OpUle: vir.OpUle,
	hir.OpUge: vir.OpUge,
	// The float comparison family drops the f: hir disambiguates by
	// mnemonic, vir by the operand type in the suffix.
	hir.OpFlt: vir.OpLt,
	hir.OpFgt: vir.OpGt,
	hir.OpFle: vir.OpLe,
	hir.OpFge: vir.OpGe,

	hir.OpSelect: vir.OpSelect,

	hir.OpAlloca:  vir.OpAlloca,
	hir.OpLoad:    vir.OpLoad,
	hir.OpStore:   vir.OpStore,
	hir.OpMemcopy: vir.OpMemcopy,
	hir.OpMemmove: vir.OpMemmove,
	hir.OpMemset:  vir.OpMemset,
	hir.OpFieldPtr: vir.OpField, // the two that a name-based bridge would miss
	hir.OpIndexPtr: vir.OpIndex,

	hir.OpAtomicLoad:  vir.OpAtomicLoad,
	hir.OpAtomicStore: vir.OpAtomicStore,
	hir.OpAtomicAdd:   vir.OpAtomicAdd,
	hir.OpAtomicSub:   vir.OpAtomicSub,
	hir.OpCmpxchg:     vir.OpCmpxchg,

	hir.OpTrunc:    vir.OpTrunc,
	hir.OpSext:     vir.OpSext,
	hir.OpZext:     vir.OpZext,
	hir.OpFdemote:  vir.OpFdemote,
	hir.OpFpromote: vir.OpFpromote,
	hir.OpBitcast:  vir.OpBitcast,
	hir.OpSfromint: vir.OpSfromint,
	hir.OpUfromint: vir.OpUfromint,
	hir.OpStoint:   vir.OpStoint,
	hir.OpUtoint:   vir.OpUtoint,

	hir.OpCall: vir.OpCall,
}

func (fl *funcLowerer) instr(in *hir.Instr) {
	switch in.Op {
	case hir.OpCall:
		fl.call(in)
		return
	case hir.OpAlloca:
		fl.alloca(in)
		return
	case hir.OpIndexPtr:
		fl.indexPtr(in)
		return
	}

	op, ok := opTable[in.Op]
	if !ok {
		fl.ml.bug("no vir opcode for hir op " + in.Op.String())
	}
	fl.b.Emit(in.Result, op, fl.suffix(in), fl.operands(in.Args)...)
}

// suffix applies the two rules worth knowing.
//
//  1. Comparisons name their operand type, not their result type. vir's
//     eq/slt/ult family yields i1 from a suffix naming what is being
//     compared, while hir.Instr.Type holds the result — so the suffix comes
//     off Args[0] for these.
//  2. A void-typed instruction has no suffix. memcopy dst, src, n takes
//     none; store.i32 p, v does, and hir passes the real type there — which
//     is why this keys on hir's Type being Void rather than on the opcode
//     being void-resulting.
func (fl *funcLowerer) suffix(in *hir.Instr) vir.Type {
	if in.Op.IsComparison() {
		if len(in.Args) == 0 {
			fl.ml.bug("comparison with no operands in " + fl.fn.Name)
		}
		return fl.ml.typ(in.Args[0].Type)
	}
	if hir.IsVoid(in.Type) {
		return nil
	}
	return fl.ml.typ(in.Type)
}

func (fl *funcLowerer) operands(args []hir.Value) []vir.Operand {
	out := make([]vir.Operand, 0, len(args))
	for _, a := range args {
		out = append(out, fl.ml.operand(a))
	}
	return out
}

// alloca carries an align clause, which the plain Emit path has nowhere to
// put. Every hir alloca is a byte count plus an alignment: Vertex has no
// dynamically sized local, so the size operand is always a constant.
func (fl *funcLowerer) alloca(in *hir.Instr) {
	fl.b.EmitInstruction(vir.Instruction{
		Result: in.Result,
		Op:     vir.OpAlloca,
		Suffix: vir.Ptr,
		Args:   fl.operands(in.Args),
		Align:  in.Align,
	})
}

// indexPtr bridges the one arity mismatch between the two vocabularies.
//
// vir's index takes (base, elem type, index) and scales by the element's
// stride. hir scales in the front end and emits (base, byte offset) — the
// bounds check and the multiply are already separate instructions above it,
// because the check has to happen whether or not the scale folds. So the
// element type this names is i8: the offset is a byte count, and a stride
// of one leaves the arithmetic hir already did untouched.
func (fl *funcLowerer) indexPtr(in *hir.Instr) {
	if len(in.Args) != 2 {
		fl.ml.bug("index.ptr with unexpected arity in " + fl.fn.Name)
	}
	fl.b.Emit(in.Result, vir.OpIndex, vir.Ptr,
		fl.ml.operand(in.Args[0]),
		vir.TypeOperand(vir.I8),
		fl.ml.operand(in.Args[1]),
	)
}

// call takes no suffix: vir derives a call's result type from the callee's
// declaration rather than from the site.
//
// Module is the owning module name for a cross-module call — a Vertex
// package or a builtins module — and "" for a call within this module. The
// import it needs was recorded by hir at the site that made it.
func (fl *funcLowerer) call(in *hir.Instr) {
	if in.Sig != "" {
		// An indirect call. Nothing upstream produces one — hir todos on
		// every call through a function value — but the spelling is here so
		// that landing fnsig upstream does not also mean landing it here.
		if len(in.Args) == 0 {
			fl.ml.bug("indirect call with no function pointer in " + fl.fn.Name)
		}
		ops := fl.operands(in.Args)
		fl.b.EmitInstruction(vir.Instruction{
			Result: in.Result,
			Op:     vir.OpCall,
			Sig:    in.Sig,
			Args:   ops,
		})
		return
	}
	if in.Callee == "" {
		fl.ml.bug("call with no callee in " + fl.fn.Name)
	}

	callee := vir.Ident(in.Callee)
	if in.Module != "" {
		callee = vir.QualifiedIdent(in.Module, in.Callee)
	}
	args := append([]vir.Operand{callee}, fl.operands(in.Args)...)
	fl.b.Emit(in.Result, vir.OpCall, nil, args...)
}