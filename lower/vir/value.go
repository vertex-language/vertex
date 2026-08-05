// value.go
package vir

import (
	"github.com/vertex-language/vertex/lower/hir"
	vir "github.com/vertex-language/vvm/ir/vir"
)

// hir.Type -> vir.Type and hir.Value -> vir.Operand. Both are total: every
// shape hir can build has exactly one vir spelling, and anything that
// reaches the default arm is a malformed program rather than a gap.

func (ml *moduleLowerer) typ(t *hir.Type) vir.Type {
	// hir spells a void result as a nil *Type on Func.Result.
	if t == nil {
		return vir.Void
	}
	switch t.Kind {
	case hir.KVoid:
		return vir.Void
	case hir.KInt:
		return vir.IntType{Bits: t.Bits}
	case hir.KFloat:
		return vir.FloatType{Bits: t.Bits}
	case hir.KPtr:
		return vir.Ptr
	case hir.KStruct:
		return ml.structType(t.Struct)
	case hir.KArray:
		return vir.ArrayType{Elem: ml.typ(t.Elem), Len: int(t.Len)}
	case hir.KVector:
		return vir.VecType{Elem: ml.typ(t.Elem), Len: int(t.Len)}
	case hir.KPredicate:
		// The lane predicate is vec[i1, N] and nothing else. hir keeps it a
		// distinct Kind because it never reaches storage; vir has no such
		// distinction to keep.
		return vir.VecType{Elem: vir.I1, Len: int(t.Len)}
	}
	ml.bug("type with no vir spelling: " + t.String())
	return vir.Void
}

// structType names a struct, qualifying it when the shape belongs to
// another module. vir treats cross-module struct identity as nominal
// per-origin, so the Import string is load-bearing rather than decorative —
// two StructTypes with the same Name and different Import are not Equal,
// which is exactly the property byval/sret comparison needs.
func (ml *moduleLowerer) structType(s *hir.Struct) vir.Type {
	if s == nil {
		ml.bug("struct type carrying no shape")
	}
	if s.Module == nil || s.Module == ml.src {
		return vir.StructType{Name: s.Name}
	}
	// hir interns a foreign shape without recording a dependency on the
	// module that owns it — it is reaching for a layout, not a symbol. vir
	// needs the import declared for the qualifier to resolve.
	ml.needImport(s.Module.Path)
	return vir.StructType{Name: s.Name, Import: s.Module.Path}
}

func (ml *moduleLowerer) operand(v hir.Value) vir.Operand {
	switch v.Kind {
	case hir.VName:
		return vir.Ident(v.Name)

	case hir.VInt:
		// An i1 literal spells true/false. It is the only place vir's bool
		// operand form can come from: hir's `bool` is i8, and an i1 only
		// ever appears as a comparison result or a branch condition.
		if it, ok := ml.typ(v.Type).(vir.IntType); ok && it.Bits == 1 {
			return vir.BoolLiteral(v.Int != 0)
		}
		return vir.IntLiteral(v.Int)

	case hir.VFloat:
		return vir.FloatLiteral(v.Flt)

	case hir.VNull:
		return vir.NullLiteral()

	case hir.VGlobal:
		// A global is an ordinary ident in operand position; the section it
		// was declared in is what distinguishes it.
		return vir.Ident(v.Name)

	case hir.VType:
		return vir.TypeOperand(ml.typ(v.Type))

	case hir.VFuncAddr:
		// vir has no address-of instruction and no `addr` operand form —
		// only InitAddressOf, inside a global initializer. A function
		// address in instruction position therefore has no spelling at all.
		// See the README's second upstream problem: this is the same hole
		// seen from the other side, and papering over it is not this
		// package's decision.
		ml.bug("function address in operand position has no vir spelling: " + v.Name)
	}
	ml.bug("empty operand reached lowering")
	return vir.Operand{}
}

// constInit projects hir's global-init grammar onto vir's. The two are the
// same grammar; hir narrowed the checker's constant expressions down to it
// already, so this is a rename and nothing else.
func (ml *moduleLowerer) constInit(c hir.ConstInit) vir.ConstInit {
	switch x := c.(type) {
	case nil:
		return vir.InitZero{}
	case hir.InitZero:
		return vir.InitZero{}
	case hir.InitScalar:
		return vir.InitLiteral{Value: ml.operand(x.Value)}
	case hir.InitBytes:
		return vir.InitByteString{Data: x.Data}
	case hir.InitAddrOf:
		// vir §6.2 permits `addr ident` for earlier functions and globals,
		// but §2.1 puts the whole global section before the whole fn
		// section, so this can never name an earlier *function*. Emitted as
		// hir built it; see the README.
		return vir.InitAddressOf{Name: x.Name}
	case hir.InitAggregate:
		elems := make([]vir.ConstInit, 0, len(x.Elems))
		for _, e := range x.Elems {
			elems = append(elems, ml.constInit(e))
		}
		return vir.InitAggregate{Elems: elems}
	}
	ml.bug("global initializer with no vir spelling")
	return vir.InitZero{}
}

// params translates a parameter list. ByVal and SRet carry struct *names*
// in vir, and hir carries the shapes themselves — the two are the same
// declaration seen from either side of the boundary.
func (ml *moduleLowerer) params(ps []*hir.Param) []vir.Param {
	out := make([]vir.Param, 0, len(ps))
	for _, p := range ps {
		vp := vir.Param{Name: p.Name, Type: ml.typ(p.Type)}
		if p.ByVal != nil {
			vp.ByVal = p.ByVal.Name
		}
		if p.SRet != nil {
			vp.SRet = p.SRet.Name
		}
		out = append(out, vp)
	}
	return out
}