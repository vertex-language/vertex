package compiler

import (
	cir "github.com/vertex-language/ir/c"
)

// setupVtxMaps registers the Vertex map runtime externs on the module.
// Maps use void* for both key and value — full typed implementation is TBD.
func setupVtxMaps(mod *cir.Module) {
	mod.Extern("v_map_new",
		cir.Returns(cir.VoidPtr),
		cir.Param("key_size", cir.UInt32),
		cir.Param("val_size", cir.UInt32),
	)
	mod.Extern("v_map_free",
		cir.Returns(cir.Void),
		cir.Param("map", cir.VoidPtr),
	)
	mod.Extern("v_map_insert",
		cir.Returns(cir.Void),
		cir.Param("map", cir.VoidPtr),
		cir.Param("key", cir.ConstPtr(cir.Void)),
		cir.Param("val", cir.ConstPtr(cir.Void)),
	)
	mod.Extern("v_map_lookup",
		cir.Returns(cir.VoidPtr),
		cir.Param("map", cir.VoidPtr),
		cir.Param("key", cir.ConstPtr(cir.Void)),
	)
	mod.Extern("v_map_remove",
		cir.Returns(cir.Bool),
		cir.Param("map", cir.VoidPtr),
		cir.Param("key", cir.ConstPtr(cir.Void)),
	)
	mod.Extern("v_map_destroy",
		cir.Returns(cir.Void),
		cir.Param("map", cir.VoidPtr),
	)
}

// setupArraysRuntime registers the arrays_Array struct and the arrays package
// function externs. The canonical source of the ABI is runtime/arrays/arrays.vs;
// names here must stay in sync with that file's package-prefixed C symbols.
func setupArraysRuntime(mod *cir.Module) (*cir.StructType, cir.Type) {
	mod.Include("<stdint.h>")
	mod.Include("<stdbool.h>")
	mod.Include("<string.h>")
	mod.Include("<stdlib.h>")

	// Mirrors the `class Array` in runtime/arrays/arrays.vs after lowering
	// with package prefix "arrays".
	arrStruct := cir.Struct("arrays_Array",
		cir.Field("data",     cir.VoidPtr),
		cir.Field("len",      cir.UInt32),
		cir.Field("capacity", cir.UInt32),
		cir.Field("elemSize", cir.UInt32),
	)
	mod.RegisterType(arrStruct)
	arrStructPtr := cir.Ptr(arrStruct)

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

	// ── arrays package runtime ────────────────────────────────────────────────
	mod.Extern("arrays_new",
		cir.Returns(arrStructPtr),
		cir.Param("elemSize", cir.UInt32),
	)
	mod.Extern("arrays_newWithCapacity",
		cir.Returns(arrStructPtr),
		cir.Param("elemSize", cir.UInt32),
		cir.Param("capacity", cir.UInt32),
	)
	mod.Extern("arrays_free",
		cir.Returns(cir.Void),
		cir.Param("arr", arrStructPtr),
	)
	mod.Extern("arrays_push",
		cir.Returns(cir.Void),
		cir.Param("arr",  arrStructPtr),
		cir.Param("elem", cir.ConstPtr(cir.Void)),
	)
	mod.Extern("arrays_unshift",
		cir.Returns(cir.Void),
		cir.Param("arr",  arrStructPtr),
		cir.Param("elem", cir.ConstPtr(cir.Void)),
	)
	mod.Extern("arrays_removeAt",
		cir.Returns(cir.Void),
		cir.Param("arr",   arrStructPtr),
		cir.Param("index", cir.UInt32),
	)
	mod.Extern("arrays_sort",
		cir.Returns(cir.Void),
		cir.Param("arr", arrStructPtr),
		cir.Param("cmp", cir.VoidPtr),
	)

	return arrStruct, arrStructPtr
}