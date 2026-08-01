// builtins/panic.go
package builtins

import vir "github.com/vertex-language/vvm/ir/vir"

// Panic returns the "panic" builtins module for t: currently just oom,
// which hir/builtin.go's builtinNew/builtinBox call on an allocation
// failure this package made on the program's behalf (A.10.1's split
// between new/resize failing politely and unique/shared panicking).
//
// oom writes a fixed message to stderr via libc write(2)/WriteFile, then
// traps — vir's own §5.3 trap terminator, not a returning call, so a
// caller never needs an unreachable check beyond what hir already emits.
func Panic(t vir.Target) *vir.Module {
	m := vir.NewModule("panic")
	m.SetNamespace("builtins")
	m.SetTarget(t.Arch, t.OS, t.ABI, t.Tiers...)

	msg := "vertex: out of memory\n"

	switch t.OS {
	case "windows":
		m.DeclareLink(vir.LinkShared, "kernel32")
		k32 := m.DeclareExternGroup("kernel32")
		k32.DeclareFunction("GetStdHandle",
			[]vir.Param{{Name: "nStdHandle", Type: vir.I32}}, vir.Ptr)
		k32.DeclareFunction("WriteFile", []vir.Param{
			{Name: "hFile", Type: vir.Ptr},
			{Name: "lpBuffer", Type: vir.Ptr},
			{Name: "nNumberOfBytesToWrite", Type: vir.I32},
			{Name: "lpNumberOfBytesWritten", Type: vir.Ptr},
			{Name: "lpOverlapped", Type: vir.Ptr},
		}, vir.I32)

		g := m.DeclareGlobal("oom_msg",
			vir.ArrayType{Elem: vir.I8, Len: int(len(msg))},
			vir.InitByteString{Data: []byte(msg)})
		_ = g

		fb := m.DeclareFunction("oom", nil, vir.Void, true)
		// STD_ERROR_HANDLE = -12 (0xFFFFFFF4 as i32)
		h := fb.Call("h", "GetStdHandle", vir.IntLiteral(-12))
		written := fb.Alloca("written", vir.IntLiteral(4), 4)
		fb.Call("ok", "WriteFile", h, vir.Ident("oom_msg"),
			vir.IntLiteral(int64(len(msg))), written, vir.NullLiteral())
		fb.Trap()

	default: // linux, darwin, and every other libc-backed target
		libName := "c"
		if t.OS == "macos" || t.OS == "ios" || t.OS == "watchos" ||
			t.OS == "tvos" || t.OS == "visionos" {
			libName = "System"
		}
		m.DeclareLink(vir.LinkShared, libName)
		lib := m.DeclareExternGroup(libName)
		lib.DeclareFunction("write", []vir.Param{
			{Name: "fd", Type: vir.I32},
			{Name: "buf", Type: vir.Ptr},
			{Name: "count", Type: vir.I64},
		}, vir.I64)

		m.DeclareGlobal("oom_msg",
			vir.ArrayType{Elem: vir.I8, Len: len(msg)},
			vir.InitByteString{Data: []byte(msg)})

		fb := m.DeclareFunction("oom", nil, vir.Void, true)
		fb.Call("n", "write", vir.IntLiteral(2), vir.Ident("oom_msg"),
			vir.IntLiteral(int64(len(msg))))
		fb.Trap()
	}

	return m
}