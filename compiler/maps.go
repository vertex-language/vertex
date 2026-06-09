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