package vir

import (
	"github.com/vertex-language/vertex/lower/hir"
	ir "github.com/vertex-language/vvm/ir/vir"
)

// value.go maps hir's operand and type vocabularies onto vir's. hir's type
// set was chosen as vir's type set, so this is a rename with two wrinkles
// worth knowing about, both marked below.

func (l *lowerer) value(hm *hir.Module, v hir.Value) ir.Operand {
	switch v.Kind {
	case hir.ValRef:
		return ir.Ident(v.Name)

	case hir.ValInt:
		// Wrinkle 1: §1.1 makes a finite float literal carry a '.' or an
		// exponent, so an integer-shaped value at a float type has to be
		// spelled as a float. hir produces this where a zero or a literal
		// widened into a float position.
		if hir.IsFloat(v.Type) {
			return ir.FloatLiteral(float64(v.Int))
		}
		return ir.IntLiteral(v.Int)

	case hir.ValFloat:
		return ir.FloatLiteral(v.Float)

	case hir.ValBool:
		return ir.BoolLiteral(v.Bool)

	case hir.ValNull:
		return ir.NullLiteral()

	case hir.ValGlobal:
		// Wrinkle 2: vir has no address-of instruction. A global's name in
		// operand position *is* its address — which is why hir spells a
		// function address as a global initialized `addr <fn>` and loads
		// it, rather than reaching for an operator that doesn't exist.
		return l.qualified(hm, v)

	case hir.ValConst:
		return l.qualified(hm, v)
	}
	l.errorf(0, "internal: no vir operand for value kind %d (%q)", v.Kind, v.Name)
	return ir.Ident("_")
}

func (l *lowerer) qualified(hm *hir.Module, v hir.Value) ir.Operand {
	if v.Module != "" && v.Module != hm.Name {
		return ir.QualifiedIdent(v.Module, v.Name)
	}
	return ir.Ident(v.Name)
}

func (l *lowerer) typ(hm *hir.Module, t hir.Type) ir.Type {
	switch x := t.(type) {
	case nil:
		return ir.Void
	case hir.IntType:
		return ir.IntType{Bits: x.Bits}
	case hir.FloatType:
		return ir.FloatType{Bits: x.Bits}
	case hir.PtrType:
		return ir.Ptr
	case hir.VoidType:
		return ir.Void
	case hir.ArrayType:
		return ir.ArrayType{Elem: l.typ(hm, x.Elem), Len: int(x.Len)}
	case hir.StructType:
		if x.Def == nil {
			l.errorf(0, "internal: struct type with no definition")
			return ir.Void
		}
		st := ir.StructType{Name: x.Def.Name}
		// Cross-module struct identity is nominal per-origin (§7.4), so a
		// shape borrowed from another module names its origin rather than
		// relying on the spelling matching.
		if x.Def.Module != nil && x.Def.Module != hm {
			st.Import = x.Def.Module.Name
		}
		return st
	}
	l.errorf(0, "internal: no vir type for %T", t)
	return ir.Void
}