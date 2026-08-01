// builtins/memory.go
package builtins

import vir "github.com/vertex-language/vvm/ir/vir"

// Memory returns the "memory" builtins module for t: libc malloc/free on
// Linux and darwin, HeapAlloc/HeapFree on Windows. Replaces the three
// hand-maintained memory_linux.vir / memory_windows.vir / memory_darwin.vir
// text sources with one Go function branching on target — the two must stay
// in lockstep, and a compiled branch can't drift the way three text files
// edited independently can.
func Memory(t vir.Target) *vir.Module {
	m := vir.NewModule("memory")
	m.SetNamespace("builtins")
	m.SetTarget(t.Arch, t.OS, t.ABI, t.Tiers...)

	switch t.OS {
	case "windows":
		m.DeclareLink(vir.LinkShared, "kernel32")
		k32 := m.DeclareExternGroup("kernel32")
		k32.DeclareFunction("GetProcessHeap", nil, vir.Ptr)
		k32.DeclareFunction("HeapAlloc", []vir.Param{
			{Name: "hHeap", Type: vir.Ptr},
			{Name: "dwFlags", Type: vir.I32},
			{Name: "dwBytes", Type: vir.I64},
		}, vir.Ptr)
		k32.DeclareFunction("HeapFree", []vir.Param{
			{Name: "hHeap", Type: vir.Ptr},
			{Name: "dwFlags", Type: vir.I32},
			{Name: "lpMem", Type: vir.Ptr},
		}, vir.I32)

		alloc := m.DeclareFunction("allocate",
			[]vir.Param{{Name: "size", Type: vir.I64}}, vir.Ptr, true)
		heap := alloc.Call("heap", "GetProcessHeap")
		p := alloc.Call("p", "HeapAlloc", heap, vir.IntLiteral(0), vir.Ident("size"))
		alloc.Return(p)

		free := m.DeclareFunction("free",
			[]vir.Param{{Name: "p", Type: vir.Ptr}}, vir.Void, true)
		heap2 := free.Call("heap", "GetProcessHeap")
		free.Call("ok", "HeapFree", heap2, vir.IntLiteral(0), vir.Ident("p"))
		free.Return()

	default: // linux and every other ELF target: libc malloc/free
		libName := "c"
		if t.OS == "macos" || t.OS == "ios" || t.OS == "watchos" ||
			t.OS == "tvos" || t.OS == "visionos" {
			// "System", not "libSystem.B.dylib" — linkdeps.go's
			// resolveMachOLinkDependencies special-cases this exact short
			// name, and entrypoint.go's linksLibC recognizes it as "a real
			// C runtime is available."
			libName = "System"
		}
		m.DeclareLink(vir.LinkShared, libName)
		lib := m.DeclareExternGroup(libName)
		lib.DeclareFunction("libc_malloc",
			[]vir.Param{{Name: "size", Type: vir.I64}}, vir.Ptr)
		lib.DeclareFunction("libc_free",
			[]vir.Param{{Name: "p", Type: vir.Ptr}}, vir.Void)

		alloc := m.DeclareFunction("allocate",
			[]vir.Param{{Name: "size", Type: vir.I64}}, vir.Ptr, true)
		p := alloc.Call("p", "libc_malloc", vir.Ident("size"))
		alloc.Return(p)

		free := m.DeclareFunction("free",
			[]vir.Param{{Name: "p", Type: vir.Ptr}}, vir.Void, true)
		free.Call("", "libc_free", vir.Ident("p"))
		free.Return()
	}

	return m
}