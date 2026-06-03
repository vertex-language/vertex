package compiler

import (
	cir "github.com/vertex-language/ir/c"
)

// ─────────────────────────────────────────────────────────────────────────────
// GLib type constants and extern declarations
//
// These model GLib's runtime: GArray, GHashTable, GString.
// They will be replaced with Vertex's own builtin/* runtime in a future
// release; GLib is used now as a reference implementation to demonstrate the
// compiler's lowering strategy.
// ─────────────────────────────────────────────────────────────────────────────

// glibTypes holds pre-built ir/c struct type objects so the lowerer can
// reference them without recreating them each time.
type glibTypes struct {
	GArray     *cir.StructType
	GArrayPtr  cir.Type
	GString    *cir.StructType
	GStringPtr cir.Type
}

// newGlibTypes constructs the GLib struct definitions.
func newGlibTypes() *glibTypes {
	gArray := cir.Struct("_GArray",
		cir.Field("data", cir.VoidPtr),
		cir.Field("len", cir.UInt32),
	)
	gString := cir.Struct("_GString",
		cir.Field("str", cir.Ptr(cir.Char)),
		cir.Field("len", cir.UIntSize),
		cir.Field("allocated_len", cir.UIntSize),
	)
	return &glibTypes{
		GArray:     gArray,
		GArrayPtr:  cir.Ptr(gArray),
		GString:    gString,
		GStringPtr: cir.Ptr(gString),
	}
}

// setupModule adds GLib includes, type registrations, and extern function
// declarations to the ir/c Module.
func setupGLib(mod *cir.Module, gt *glibTypes) {
	mod.Include("<stdint.h>")
	mod.Include("<stdbool.h>")
	mod.Include("<string.h>")
	mod.Include("<stdlib.h>")
	mod.Include("<glib.h>")

	// Register struct layouts so EmitC() emits forward declarations.
	mod.RegisterType(gt.GArray)
	mod.RegisterType(gt.GString)

	// ── Standard C externs ────────────────────────────────────────────────────
	mod.Extern("malloc", cir.Returns(cir.VoidPtr), cir.Param("size", cir.UIntSize))
	mod.Extern("free", cir.Returns(cir.Void), cir.Param("ptr", cir.VoidPtr))
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
	mod.Extern("sizeof_dummy") // placeholder; sizeof is handled by b.SizeOf

	// ── GArray externs ────────────────────────────────────────────────────────
	gArrayPtrType := gt.GArrayPtr

	mod.Extern("g_array_new",
		cir.Returns(gArrayPtrType),
		cir.Param("zero_terminated", cir.Bool),
		cir.Param("clear_", cir.Bool),
		cir.Param("element_size", cir.UInt32),
	)
	mod.Extern("g_array_sized_new",
		cir.Returns(gArrayPtrType),
		cir.Param("zero_terminated", cir.Bool),
		cir.Param("clear_", cir.Bool),
		cir.Param("element_size", cir.UInt32),
		cir.Param("reserved_size", cir.UInt32),
	)
	mod.Extern("g_array_free",
		cir.Returns(cir.VoidPtr),
		cir.Param("array", gArrayPtrType),
		cir.Param("free_segment", cir.Bool),
	)
	mod.Extern("g_array_append_vals",
		cir.Returns(gArrayPtrType),
		cir.Param("array", gArrayPtrType),
		cir.Param("data", cir.ConstPtr(cir.Void)),
		cir.Param("len", cir.UInt32),
	)
	mod.Extern("g_array_prepend_vals",
		cir.Returns(gArrayPtrType),
		cir.Param("array", gArrayPtrType),
		cir.Param("data", cir.ConstPtr(cir.Void)),
		cir.Param("len", cir.UInt32),
	)
	mod.Extern("g_array_remove_index",
		cir.Returns(gArrayPtrType),
		cir.Param("array", gArrayPtrType),
		cir.Param("index_", cir.UInt32),
	)
	mod.Extern("g_array_sort",
		cir.Returns(cir.Void),
		cir.Param("array", gArrayPtrType),
		cir.Param("compare_func", cir.VoidPtr),
	)

	// ── GHashTable externs ────────────────────────────────────────────────────
	// GHashTable is opaque; use void* for the pointer type.
	mod.Extern("g_hash_table_new",
		cir.Returns(cir.VoidPtr),
		cir.Param("hash_func", cir.VoidPtr),
		cir.Param("key_equal_func", cir.VoidPtr),
	)
	mod.Extern("g_hash_table_insert",
		cir.Returns(cir.Bool),
		cir.Param("hash_table", cir.VoidPtr),
		cir.Param("key", cir.VoidPtr),
		cir.Param("value", cir.VoidPtr),
	)
	mod.Extern("g_hash_table_lookup",
		cir.Returns(cir.VoidPtr),
		cir.Param("hash_table", cir.VoidPtr),
		cir.Param("key", cir.ConstPtr(cir.Void)),
	)
	mod.Extern("g_hash_table_remove",
		cir.Returns(cir.Bool),
		cir.Param("hash_table", cir.VoidPtr),
		cir.Param("key", cir.ConstPtr(cir.Void)),
	)
	mod.Extern("g_hash_table_destroy",
		cir.Returns(cir.Void),
		cir.Param("hash_table", cir.VoidPtr),
	)

	// ── GString externs ───────────────────────────────────────────────────────
	gStrPtrType := gt.GStringPtr

	mod.Extern("g_string_new",
		cir.Returns(gStrPtrType),
		cir.Param("init", cir.ConstPtr(cir.Char)),
	)
	mod.Extern("g_string_free",
		cir.Returns(cir.Ptr(cir.Char)),
		cir.Param("string", gStrPtrType),
		cir.Param("free_segment", cir.Bool),
	)
	mod.Extern("g_string_append",
		cir.Returns(gStrPtrType),
		cir.Param("string", gStrPtrType),
		cir.Param("val", cir.ConstPtr(cir.Char)),
	)
	mod.Extern("g_string_append_len",
		cir.Returns(gStrPtrType),
		cir.Param("string", gStrPtrType),
		cir.Param("val", cir.ConstPtr(cir.Char)),
		cir.Param("len", cir.UIntSize),
	)
}