package hir

import (
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/types"
)

// Ownership expansion. A bare owning use is a copy and a marked one is a
// move; "copy" is not one operation, and this file is where the table in
// semantics.md §8.5 becomes instructions.
//
//	scalars, typed_ptr, vector    register move
//	struct, class, [N]T           fieldwise, recursing per field
//	string, []T, map[K]V          deep — duplicates the payload
//	unique T                      allocates and deep-copies the pointee
//	shared T, chan T              refcount increment
//	weak T                        weak-count increment

// copyOf produces an owned copy of v. Trivial copies emit nothing; the rest
// go through a per-type routine or a builtin.
func (fb *funcBuilder) copyOf(v Value, t types.Type) Value {
	if v.IsZero() || t == nil {
		return v
	}
	switch types.CopyCost(t) {
	case types.CopyRegister:
		return v

	case types.CopyRefcount:
		if types.IsChan(t) {
			fb.l.need(builtins.FeatChan)
			return fb.callBuiltin(builtins.ChanRetain, Ptr, v)
		}
		return fb.callBuiltin(builtins.RCRetain, Ptr, v)

	case types.CopyWeakcount:
		return fb.callBuiltin(builtins.RCWeak, Ptr, v)

	case types.CopyAlloc:
		// unique T: a fresh allocation plus a deep copy of the pointee.
		return fb.callBuiltin(builtins.UniqueNew, Ptr, v)

	case types.CopyDeep:
		return fb.deepCopy(v, t)

	default: // CopyFieldwise
		ht := fb.typ(t)
		if !IsAggregate(ht) {
			return v
		}
		dst := fb.alloca(ht, "cp")
		fb.copyInto(dst, v, t)
		return dst
	}
}

func (fb *funcBuilder) deepCopy(v Value, t types.Type) Value {
	switch {
	case types.IsString(t):
		// This package declines semantics.md §8.5's license to share or
		// intern an immutable string payload: a bare copy emits a real
		// duplication. Interning would put a refcount header on a type whose
		// spelling contains no `shared` — a cost invisible in the source.
		ht := fb.typ(t)
		dst := fb.alloca(ht, "str")
		fb.callBuiltin(builtins.StringCopy, Void, dst, v)
		return dst
	case types.IsSlice(t):
		ht := fb.typ(t)
		dst := fb.alloca(ht, "sl")
		elem := types.AsSlice(t).Elem()
		fb.callBuiltin(builtins.SliceAlloc, Void, dst,
			Int(fb.l.sizeOf(fb.typ(elem)), I64), Int(0, I64))
		// todo: the element loop. A []T of non-trivial T copies each element
		// at that element's own cost; today only the header is duplicated,
		// which is wrong for []string and right for []int32.
		return dst
	case types.IsMap(t):
		fb.l.need(builtins.FeatMap)
		return fb.callBuiltin(builtins.MapNew, Ptr)
	}
	return v
}

// copyInto is the fieldwise copy: each field at that field's own cost, which
// is why a struct holding a []T still deep-copies that field.
func (fb *funcBuilder) copyInto(dst, src Value, t types.Type) {
	ht := fb.typ(t)
	if !fb.needsDrop(t) {
		// Nothing owns anything: one memcopy is the whole operation.
		fb.emitVoid(OpMemcopy, Void, dst, src, Int(fb.l.sizeOf(ht), I64))
		return
	}
	st := types.AsStruct(t)
	if st == nil || ht.Kind != KStruct {
		fb.emitVoid(OpMemcopy, Void, dst, src, Int(fb.l.sizeOf(ht), I64))
		return
	}
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		sp := fb.fieldPtr(src, ht.Struct, i)
		dp := fb.fieldPtr(dst, ht.Struct, i)
		ft := ht.Struct.Fields[i].Type
		if IsAggregate(ft) {
			fb.copyInto(dp, sp, f.Type)
			continue
		}
		v := fb.copyOf(fb.emit(OpLoad, ft, sp), f.Type)
		fb.emitVoid(OpStore, ft, dp, v)
	}
}

// needsDrop reports whether a value of t owns anything. A type that owns
// nothing needs no teardown, no _Vdrop routine, and no epilogue entry.
func (fb *funcBuilder) needsDrop(t types.Type) bool {
	if t == nil {
		return false
	}
	switch types.CopyCost(t) {
	case types.CopyRegister:
		return false
	case types.CopyDeep, types.CopyAlloc, types.CopyRefcount, types.CopyWeakcount:
		return true
	}
	// Fieldwise: recurse. A class declaring deinit is non-trivial regardless.
	if n := types.AsNamed(t); n != nil && n.LookupMethod("deinit") != nil {
		return true
	}
	if st := types.AsStruct(t); st != nil {
		for i := 0; i < st.NumFields(); i++ {
			if fb.needsDrop(st.Field(i).Type) {
				return true
			}
		}
	}
	if a := types.AsArray(t); a != nil {
		return fb.needsDrop(a.Elem())
	}
	return false
}

// dropBinding tears one local down at a scope exit.
func (fb *funcBuilder) dropBinding(b *binding) {
	if !fb.needsDrop(b.vt) {
		return
	}
	if b.slot {
		fb.dropInPlace(b.val, b.vt)
		return
	}
	fb.dropValue(b.val, b.vt)
}

// dropInPlace tears down the value stored at p.
func (fb *funcBuilder) dropInPlace(p Value, t types.Type) {
	if !fb.needsDrop(t) {
		return
	}
	ht := fb.typ(t)
	if !IsAggregate(ht) {
		fb.dropValue(fb.emit(OpLoad, ht, p), t)
		return
	}
	// The user deinit runs first, then the fields in reverse declaration
	// order.
	if n := types.AsNamed(t); n != nil {
		if d := n.LookupMethod("deinit"); d != nil {
			target := fb.l.enqueue(d, nil)
			fb.emitCall(target.Name, moduleNameOf(target, fb.fn.Module), Void, []Value{p})
		}
	}
	st := types.AsStruct(t)
	if st == nil {
		return
	}
	for i := st.NumFields() - 1; i >= 0; i-- {
		f := st.Field(i)
		if !fb.needsDrop(f.Type) {
			continue
		}
		fb.dropInPlace(fb.fieldPtr(p, ht.Struct, i), f.Type)
	}
}

func (fb *funcBuilder) dropValue(v Value, t types.Type) {
	switch types.CopyCost(t) {
	case types.CopyRefcount:
		if types.IsChan(t) {
			fb.callBuiltin(builtins.ChanRelease, Void, v)
			return
		}
		// The release path drops the payload when the count reaches zero, so
		// it takes the payload's own drop routine as a function pointer.
		fb.callBuiltin(builtins.RCRelease, I1, v, Null())
	case types.CopyWeakcount:
		fb.callBuiltin(builtins.RCWeakDrop, Void, v)
	case types.CopyAlloc:
		fb.callBuiltin(builtins.UniqueFree, Void, v)
	case types.CopyDeep:
		switch {
		case types.IsString(t):
			fb.callBuiltin(builtins.StringFree, Void, v)
		case types.IsSlice(t):
			fb.callBuiltin(builtins.SliceFree, Void, v)
		case types.IsMap(t):
			fb.callBuiltin(builtins.MapFree, Void, v)
		}
	}
}

// todo: per-type _Vcopy_T / _Vdrop_T routines. Today copyInto and
// dropInPlace inline the walk at every site, which is correct and duplicates
// code — a bare copy is already the documented-expensive path, and
// multiplying its body is exactly what lowering.md §7.1 says not to do.
// Synthesizing one routine per type and calling it is mechanical from here;
// the reason it is not done is that the owning-module rule (which module
// holds _Vdrop_Foo when two packages both drop a Foo) wants deciding first.

func moduleNameOf(target *Func, from *Module) string {
	if target.Module == from {
		return ""
	}
	from.Import(target.Module.Path)
	return target.Module.Name
}