package hir

import (
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// ownership synthesizes and dispatches the per-type routines that "copy"
// and "teardown" actually mean. types.Info.Transfers answers move-vs-copy
// at each owning position, but copy is not one operation — it is the table
// in overview §3, and this file is that table's right-hand column.
//
// A.9.4 prices the difference: a bare copy of an owning fat type duplicates
// header *and* payload, O(data); a transfer copies the header and stops,
// O(1). unique T is the one cost cliff hidden behind a thin type, which is
// why its marker is visible rather than inferred.
type ownership struct {
	l     *lowerer
	copies map[string]*Func
	deinits map[string]*Func
}

func newOwnership(l *lowerer) *ownership {
	return &ownership{l: l, copies: map[string]*Func{}, deinits: map[string]*Func{}}
}

// emitCopy produces a copy of v at an owning position where no transfer
// marker was written.
func (o *ownership) emitCopy(b *funcBuilder, pos token.Pos, t types.Type, v Value) Value {
	switch o.l.classify(t) {
	case kString:
		// A.9.4 licenses sharing or interning an immutable string payload.
		// Declining the license is the safe default, and this is the
		// recorded choice: a string copy is a real duplication.
		o.l.need(builtins.FeatString)
		st := o.l.hirType(t).(StructType)
		dst := b.alloca(pos, st)
		b.callBuiltin(pos, symStringCopy, Void, dst, v)
		return Value{Kind: ValRef, Name: dst.Name, Type: st}

	case kShared, kChan:
		// Always cheap: an atomic increment, regardless of marker.
		s := symRCRetain
		if o.l.classify(t) == kChan {
			s = builtins.ChanRetain
		}
		return b.callBuiltin(pos, s, Ptr, v)

	case kWeak:
		return b.callBuiltin(pos, builtins.RCWeak, Ptr, v)

	case kSlice, kUnique, kMap, kStruct, kEnum:
		f := o.copyRoutine(t)
		if f == nil {
			return v
		}
		ht := o.l.hirType(t)
		dst := b.alloca(pos, ht)
		b.call(pos, f, dst, v)
		return Value{Kind: ValRef, Name: dst.Name, Type: ht}
	}
	return v // scalar, mut, non-owning view: a bit copy
}

// emitDeinit runs a binding's teardown where its liveness ends. Fields go
// in reverse declaration order, locals in reverse declaration order.
func (o *ownership) emitDeinit(b *funcBuilder, pos token.Pos, bd *binding) {
	switch o.l.classify(bd.Src) {
	case kString:
		b.callBuiltin(pos, symStringFree, Void, bd.Value)
	case kSlice:
		b.callBuiltin(pos, symSliceFree, Void, bd.Value)
	case kMap:
		b.callBuiltin(pos, symMapFree, Void, bd.Value)
	case kChan:
		b.callBuiltin(pos, symChanRelease, Void, bd.Value)
	case kShared:
		// release takes the payload's deinit so the last owner runs it.
		b.callBuiltin(pos, symRCRelease, I1, bd.Value, o.deinitAddr(b, pos, o.l.elem(bd.Src)))
	case kWeak:
		b.callBuiltin(pos, builtins.RCWeakDrop, Void, bd.Value)
	case kUnique, kStruct, kEnum:
		if f := o.deinitRoutine(bd.Src); f != nil {
			b.call(pos, f, bd.Value)
		}
	}
}

// copyRoutine synthesizes the per-type copy. Tier-2 symbols have exactly
// one owning module: the module declaring the type, or for a synthesized
// container routine, the module declaring the element type. Two packages
// copying the same []Foo reference the routine, they do not each emit it.
func (o *ownership) copyRoutine(t types.Type) *Func {
	key := types.TypeString(t)
	if f, ok := o.copies[key]; ok {
		return f
	}
	mod := o.owningModule(t)
	ht := o.l.types.lower(mod, t)
	f := &Func{
		Name: mod.uniqueName("copy__" + sanitize(key)), Module: mod, Kind: FuncCopy,
		Params: []*Param{{Name: "dst", Type: Ptr}, {Name: "src", Type: Ptr}},
		Result: Void, Export: true,
	}
	mod.Funcs = append(mod.Funcs, f)
	o.copies[key] = f

	prev := o.l.cur
	o.l.cur = &instance{unit: prev.unit, mod: mod}
	defer func() { o.l.cur = prev }()

	b := newFuncBuilder(o.l, f, nil)
	dst, src := Ref("dst", Ptr), Ref("src", Ptr)
	o.copyInto(b, ht, t, dst, src)
	b.seq.add(&Return{})
	return f
}

// copyInto is the recursion. A struct recurses field by field; an enum
// recurses into the live variant only — the tag tells the copy routine
// which interpretation to walk.
func (o *ownership) copyInto(b *funcBuilder, ht Type, t types.Type, dst, src Value) {
	switch o.l.classify(t) {
	case kSlice:
		// []T must genuinely duplicate: header plus payload, then an
		// element loop for an owning element type.
		o.l.need(builtins.FeatSlice)
		st := ht.(StructType)
		n := b.loadField(0, st.Def, src, "len")
		esz := o.l.types.Sizeof(o.l.hirType(o.l.elem(t)))
		b.callBuiltin(0, builtins.SliceAlloc, Void, dst, IntVal(I64, esz), n)
		o.elementLoop(b, t, dst, src, n)

	case kUnique:
		// Deep: walks and duplicates the pointee.
		inner := o.l.elem(t)
		hi := o.l.hirType(inner)
		h := b.callBuiltin(0, symUniqueNew, Ptr, IntVal(I64, o.l.types.Sizeof(hi)))
		o.copyInto(b, hi, inner, h, src)
		b.store(0, Ptr, dst, h)

	case kStruct:
		st := ht.(StructType)
		_, fields := o.l.structParts(t)
		b.opVoid(0, OpMemcopy, Void, dst, src, IntVal(I64, st.Def.Size))
		for i, f := range st.Def.Fields {
			if i >= len(fields) || !o.l.classify(fields[i].Type).owning() {
				continue
			}
			o.copyInto(b, f.Type, fields[i].Type,
				b.fieldPtr(0, st.Def, dst, f.Name), b.fieldPtr(0, st.Def, src, f.Name))
		}

	case kEnum:
		o.l.todo(0, "enum copy routine — switch on the tag, recurse into the live variant only")

	default:
		b.store(0, ht, dst, b.load(0, ht, src))
	}
}

func (o *ownership) elementLoop(b *funcBuilder, t types.Type, dst, src, n Value) {
	elem := o.l.elem(t)
	if !o.l.classify(elem).owning() {
		return // the bulk copy already duplicated plain elements
	}
	he := o.l.hirType(elem)
	i := b.alloca(0, I64)
	b.store(0, I64, i, IntVal(I64, 0))
	loop := &Loop{}
	loop.Body = b.into(func() {
		b.push(scopeLoop)
		cur := b.load(0, I64, i)
		b.seq.add(&If{Cond: b.op(0, OpUge, I1, cur, n), Then: &Seq{List: []Stmt{&Break{}}}})
		o.copyInto(b, he, elem, b.indexPtr(0, he, dst, cur), b.indexPtr(0, he, src, cur))
		b.store(0, I64, i, b.op(0, OpAdd, I64, cur, IntVal(I64, 1)))
		b.pop()
	})
	b.seq.add(loop)
}

func (o *ownership) deinitRoutine(t types.Type) *Func {
	key := types.TypeString(t)
	if f, ok := o.deinits[key]; ok {
		return f
	}
	mod := o.owningModule(t)
	f := &Func{
		Name: mod.uniqueName("deinit__" + sanitize(key)), Module: mod, Kind: FuncDeinit,
		Params: []*Param{{Name: "p", Type: Ptr}}, Result: Void, Export: true,
	}
	mod.Funcs = append(mod.Funcs, f)
	o.deinits[key] = f

	prev := o.l.cur
	o.l.cur = &instance{unit: prev.unit, mod: mod}
	defer func() { o.l.cur = prev }()

	b := newFuncBuilder(o.l, f, nil)
	o.l.emitUserDeinit(b, t, Ref("p", Ptr)) // the declared deinit, if any
	b.seq.add(&Return{})
	return f
}

// deinitAddr yields a function pointer to a type's teardown, for the
// builtins that must run it themselves (rc.release's last-owner path).
func (o *ownership) deinitAddr(b *funcBuilder, pos token.Pos, t types.Type) Value {
	f := o.deinitRoutine(t)
	if f == nil {
		return NullVal()
	}
	name := "addr_" + f.Name
	if findGlobal(b.mod(), name) == nil {
		b.mod().Globals = append(b.mod().Globals, &Global{
			Name: b.mod().uniqueName(name), Type: Ptr, Init: InitAddr{Name: f.Name},
		})
	}
	return b.load(pos, Ptr, GlobalVal(name, Ptr))
}

// owningModule picks the one module a synthesized symbol lives in: the
// module declaring the generic function, or for a per-type routine, the
// module declaring the type. vir has a strict flat namespace and
// declare-before-use, so two packages both copying a []Foo would otherwise
// produce a duplicate symbol rather than a silently-merged one.
func (o *ownership) owningModule(t types.Type) *Module {
	if n, ok := o.l.subst(t).(*types.Named); ok && n.Obj() != nil && n.Obj().Pkg() != nil {
		if m := o.l.modules[n.Obj().Pkg().Path()]; m != nil {
			return m
		}
	}
	if e := o.l.elem(t); e != nil {
		if n, ok := o.l.subst(e).(*types.Named); ok && n.Obj() != nil && n.Obj().Pkg() != nil {
			if m := o.l.modules[n.Obj().Pkg().Path()]; m != nil {
				return m
			}
		}
	}
	// A routine over nothing but predeclared types has no declaring
	// package; it belongs to the root module.
	return o.l.prog.Modules[0]
}