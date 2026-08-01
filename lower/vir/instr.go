package vir

import (
	"github.com/vertex-language/vertex/lower/hir"
	ir "github.com/vertex-language/vvm/ir/vir"
)

// instr.go is the table. hir's Op vocabulary was chosen as a subset of
// vir's §4 opcodes spelled identically, so almost every instruction is one
// map lookup plus an operand walk.
//
// The table is explicit rather than a name lookup through
// vir.ParseOpcode(op.String()) on purpose: hir spells its pointer
// arithmetic ops "field.ptr" and "index.ptr" — the *instruction* form,
// suffix included — while vir's closed opcode vocabulary spells them
// "field" and "index". A name-based bridge would work for sixty opcodes
// and silently fail for two.
var opcodes = map[hir.Op]ir.Opcode{
	hir.OpAdd:  ir.OpAdd,
	hir.OpSub:  ir.OpSub,
	hir.OpMul:  ir.OpMul,
	hir.OpUDiv: ir.OpUDiv,
	hir.OpSDiv: ir.OpSDiv,
	hir.OpURem: ir.OpURem,
	hir.OpSRem: ir.OpSRem,
	hir.OpNeg:  ir.OpNeg,
	hir.OpAbs:  ir.OpAbs,

	hir.OpAnd:  ir.OpAnd,
	hir.OpOr:   ir.OpOr,
	hir.OpXor:  ir.OpXor,
	hir.OpNot:  ir.OpNot,
	hir.OpShl:  ir.OpShl,
	hir.OpLShr: ir.OpLShr,
	hir.OpAShr: ir.OpAShr,

	hir.OpEq:  ir.OpEq,
	hir.OpNe:  ir.OpNe,
	hir.OpSlt: ir.OpSlt,
	hir.OpSgt: ir.OpSgt,
	hir.OpSle: ir.OpSle,
	hir.OpSge: ir.OpSge,
	hir.OpUlt: ir.OpUlt,
	hir.OpUgt: ir.OpUgt,
	hir.OpUle: ir.OpUle,
	hir.OpUge: ir.OpUge,
	hir.OpLt:  ir.OpLt,
	hir.OpGt:  ir.OpGt,
	hir.OpLe:  ir.OpLe,
	hir.OpGe:  ir.OpGe,

	hir.OpSelect: ir.OpSelect,

	hir.OpLoad:    ir.OpLoad,
	hir.OpStore:   ir.OpStore,
	hir.OpMemcopy: ir.OpMemcopy,
	hir.OpMemmove: ir.OpMemmove,
	hir.OpMemset:  ir.OpMemset,

	hir.OpTrunc:      ir.OpTrunc,
	hir.OpSext:       ir.OpSext,
	hir.OpZext:       ir.OpZext,
	hir.OpFdemote:    ir.OpFdemote,
	hir.OpFpromote:   ir.OpFpromote,
	hir.OpBitcast:    ir.OpBitcast,
	hir.OpSfromint:   ir.OpSfromint,
	hir.OpUfromint:   ir.OpUfromint,
	hir.OpStointSat:  ir.OpStointSat,
	hir.OpUtointSat:  ir.OpUtointSat,

	hir.OpSMin: ir.OpSMin,
	hir.OpSMax: ir.OpSMax,
	hir.OpUMin: ir.OpUMin,
	hir.OpUMax: ir.OpUMax,

	// OpAlloca, OpFieldPtr, OpIndexPtr and OpCall are absent deliberately:
	// each has a dedicated builder entry point below, because vir's form
	// carries something an Emit(result, op, suffix, args...) call can't —
	// an align clause, a struct/field name pair, an element type operand,
	// or a callee.
}

// comparisons are the opcodes whose result type is i1 rather than their
// suffix (opTable's ruleBool). hir's Instr.Type holds the *result* type
// for these, but vir's suffix names the *operand* type — `eq.ptr` yields
// i1 — so the suffix is read off the first argument instead.
var comparisons = map[hir.Op]bool{
	hir.OpEq: true, hir.OpNe: true,
	hir.OpSlt: true, hir.OpSgt: true, hir.OpSle: true, hir.OpSge: true,
	hir.OpUlt: true, hir.OpUgt: true, hir.OpUle: true, hir.OpUge: true,
	hir.OpLt: true, hir.OpGt: true, hir.OpLe: true, hir.OpGe: true,
}

func (e *emitter) instr(in *hir.Instr) {
	e.loc(in.Pos)

	switch in.Op {
	case hir.OpCall:
		e.call(in)
		return

	case hir.OpAlloca:
		// alloca's align is a clause, not an operand (§2.3).
		if len(in.Args) != 1 {
			e.l.errorf(in.Pos, "internal: alloca takes one size operand, got %d", len(in.Args))
			return
		}
		e.fb.Alloca(in.Name, e.l.value(e.hm, in.Args[0]), in.Align)
		return

	case hir.OpFieldPtr:
		if in.Field == nil || len(in.Args) != 1 {
			e.l.errorf(in.Pos, "internal: field.ptr without a struct/field pair")
			return
		}
		e.fb.FieldPointer(in.Name, e.l.value(e.hm, in.Args[0]), in.Field.Struct.Name, in.Field.Field)
		return

	case hir.OpIndexPtr:
		if in.Elem == nil || len(in.Args) != 2 {
			e.l.errorf(in.Pos, "internal: index.ptr without an element type")
			return
		}
		e.fb.IndexPointer(in.Name, e.l.value(e.hm, in.Args[0]),
			e.l.typ(e.hm, in.Elem), e.l.value(e.hm, in.Args[1]))
		return
	}

	op, ok := opcodes[in.Op]
	if !ok {
		e.l.errorf(in.Pos, "internal: hir opcode %s has no vir spelling", in.Op)
		return
	}
	e.fb.Emit(in.Name, op, e.suffix(in), e.operands(in.Args)...)
}

// suffix picks the instruction's type suffix. Three cases, and they are
// the whole rule: a comparison names its operand type, a void-typed
// instruction has no suffix at all (`memcopy dst, src, n`), and everything
// else names its own result type.
func (e *emitter) suffix(in *hir.Instr) ir.Type {
	if comparisons[in.Op] {
		if len(in.Args) == 0 || in.Args[0].Type == nil {
			e.l.errorf(in.Pos, "internal: comparison %s has no operand to take a suffix from", in.Op)
			return nil
		}
		return e.l.typ(e.hm, in.Args[0].Type)
	}
	if in.Type == nil || hir.IsVoid(in.Type) {
		return nil
	}
	return e.l.typ(e.hm, in.Type)
}

// call routes the three callee shapes. A qualified call is emitted as
// `module.symbol` and erased by vvm's importer Rewrite before cpu/lower
// ever sees it.
//
// This builds the ir.Instruction directly rather than going through
// builder.go's Call/CallImported/CallIndirect convenience wrappers, which
// all hardcode Suffix to nil. That was invisible for a locally-resolved
// call — cpu/lower/<arch> derived the result type from its own function
// table instead — but a rewritten cross-module call (importer.Rewrite
// erases the qualified ident into a bare mangled symbol, per its own
// per-kind summary) has no local declaration for a backend to consult, and
// needs Suffix to carry the checked result type hir already computed
// (hir's builder.go: `Instr{..., Type: f.Result, ...}`) — the same
// contract every other ruleSuffix opcode already relies on. Setting it
// unconditionally, for every call shape, keeps this one rule instead of a
// special case for the cross-module path alone.
func (e *emitter) call(in *hir.Instr) {
	c := in.Call
	if c == nil {
		e.l.errorf(in.Pos, "internal: call instruction with no callee")
		return
	}
	args := e.operands(in.Args)
	suffix := e.suffix(in)

	switch {
	case c.Indirect != nil:
		if c.Sig == "" {
			e.l.todo(in.Pos, "indirect call with no fnsig — vir types call.<fnsig> against a declared signature, and nothing declares one yet")
			return
		}
		e.fb.EmitInstruction(ir.Instruction{
			Result: in.Name, Op: ir.OpCall, Suffix: suffix, Sig: c.Sig,
			Args: append([]ir.Operand{e.l.value(e.hm, *c.Indirect)}, args...),
		})
	case c.Module != "":
		e.fb.EmitInstruction(ir.Instruction{
			Result: in.Name, Op: ir.OpCall, Suffix: suffix,
			Args: append([]ir.Operand{ir.QualifiedIdent(c.Module, c.Name)}, args...),
		})
	default:
		e.fb.EmitInstruction(ir.Instruction{
			Result: in.Name, Op: ir.OpCall, Suffix: suffix,
			Args: append([]ir.Operand{ir.Ident(c.Name)}, args...),
		})
	}
}

func (e *emitter) operands(vs []hir.Value) []ir.Operand {
	if len(vs) == 0 {
		return nil
	}
	out := make([]ir.Operand, 0, len(vs))
	for _, v := range vs {
		out = append(out, e.l.value(e.hm, v))
	}
	return out
}