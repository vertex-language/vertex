package compiler

import (
	cir "github.com/vertex-language/ir/c"
)

// vtxArrayTypes holds the __vtx_array struct and pointer type so the lowerer
// can reference them without recreating them on every call.
type vtxArrayTypes struct {
	VtxArray    *cir.StructType
	VtxArrayPtr cir.Type
}

func newVtxArrayTypes() *vtxArrayTypes {
	vtxArray := cir.Struct("vtx_array",
		cir.Field("data",      cir.VoidPtr),
		cir.Field("len",       cir.UInt32),
		cir.Field("capacity",  cir.UInt32),
		cir.Field("elem_size", cir.UInt32),
	)
	return &vtxArrayTypes{
		VtxArray:    vtxArray,
		VtxArrayPtr: cir.Ptr(vtxArray),
	}
}

// setupVtxArrays registers the vtx_array struct, standard C externs, and the
// Vertex array runtime externs on the module.  Matches the shape of the old
// setupGLib so the lowerer only needs to swap symbol names.
func setupVtxArrays(mod *cir.Module, at *vtxArrayTypes) {
	mod.Include("<stdint.h>")
	mod.Include("<stdbool.h>")
	mod.Include("<string.h>")
	mod.Include("<stdlib.h>")

	mod.RegisterType(at.VtxArray)

	arrPtr := at.VtxArrayPtr

	// ── Standard C ────────────────────────────────────────────────────────────
	mod.Extern("malloc",  cir.Returns(cir.VoidPtr), cir.Param("size", cir.UIntSize))
	mod.Extern("realloc", cir.Returns(cir.VoidPtr), cir.Param("ptr", cir.VoidPtr), cir.Param("size", cir.UIntSize))
	mod.Extern("free",    cir.Returns(cir.Void),    cir.Param("ptr", cir.VoidPtr))
	mod.Extern("memset",
		cir.Returns(cir.VoidPtr),
		cir.Param("s", cir.VoidPtr),
		cir.Param("c", cir.Int32),
		cir.Param("n", cir.UIntSize),
	)
	mod.Extern("memcpy",
		cir.Returns(cir.VoidPtr),
		cir.Param("dst", cir.VoidPtr),
		cir.Param("src", cir.ConstPtr(cir.Void)),
		cir.Param("n", cir.UIntSize),
	)
	mod.Extern("memmove",
		cir.Returns(cir.VoidPtr),
		cir.Param("dst", cir.VoidPtr),
		cir.Param("src", cir.ConstPtr(cir.Void)),
		cir.Param("n", cir.UIntSize),
	)
	mod.Extern("strlen",
		cir.Returns(cir.UIntSize),
		cir.Param("s", cir.ConstPtr(cir.Char)),
	)
	mod.Extern("qsort",
		cir.Returns(cir.Void),
		cir.Param("base",   cir.VoidPtr),
		cir.Param("nmemb",  cir.UIntSize),
		cir.Param("size",   cir.UIntSize),
		cir.Param("compar", cir.VoidPtr),
	)

	// ── Vertex array runtime ──────────────────────────────────────────────────
	mod.Extern("v_array_new",
		cir.Returns(arrPtr),
		cir.Param("elem_size", cir.UInt32),
	)
	mod.Extern("v_array_new_cap",
		cir.Returns(arrPtr),
		cir.Param("elem_size", cir.UInt32),
		cir.Param("capacity",  cir.UInt32),
	)
	mod.Extern("v_array_free",
		cir.Returns(cir.Void),
		cir.Param("arr", arrPtr),
	)
	mod.Extern("v_array_push",
		cir.Returns(cir.Void),
		cir.Param("arr",  arrPtr),
		cir.Param("elem", cir.ConstPtr(cir.Void)),
	)
	mod.Extern("v_array_unshift",
		cir.Returns(cir.Void),
		cir.Param("arr",  arrPtr),
		cir.Param("elem", cir.ConstPtr(cir.Void)),
	)
	mod.Extern("v_array_remove_index",
		cir.Returns(cir.Void),
		cir.Param("arr",   arrPtr),
		cir.Param("index", cir.UInt32),
	)
	mod.Extern("v_array_sort",
		cir.Returns(cir.Void),
		cir.Param("arr", arrPtr),
		cir.Param("cmp", cir.VoidPtr),
	)
}