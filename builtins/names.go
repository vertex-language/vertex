// Package builtins is the ABI between the compiler and the vir modules it
// manufactures to make the language's own features work.
//
// This file is item (1) of the package's three responsibilities (overview
// §5.5): the ABI as constants. lower/hir imports *only* this file's contents
// and features.go — never a module constructor — so that what decides which
// call to emit cannot see how the callee is built (invariant 10).
package builtins

// Namespace is the vir `namespace` every builtin module declares. It is the
// single place the spelling is decided (§5.1, point 2): nothing forces
// "builtins", and changing it here changes every mangled symbol at once.
const Namespace = "builtins"

// Tier 0 — platform. The only modules carrying target/link/extern.
const (
	ModuleMemory  = "memory"
	ModuleThread  = "thread"
	ModulePoll    = "poll"
	ModuleSync    = "sync"
	ModuleIO      = "io"
	ModuleConsole = "console"
	ModuleTime    = "time"
	ModuleProcess = "process"
)

// Tier 1 — portable. Pure compute; reach the platform only by importing
// tier 0.
const (
	ModuleRC      = "rc"
	ModuleSlice   = "slice"
	ModuleMap     = "map"
	ModuleString  = "string"
	ModuleChan    = "chan"
	ModuleReactor = "reactor"
	ModulePanic   = "panic"
	ModuleFmt     = "fmt"
)

// Symbol is one builtin entry point: the module that owns it and the
// function name within it. lower/hir emits `import "<Module>"` plus a
// qualified `call <Module>.<Func>`; vvm's importer resolves and rewrites it
// to the mangled extern symbol.
type Symbol struct {
	Module string
	Func   string
}

func (s Symbol) String() string { return s.Module + "." + s.Func }

// The heap door (§5.2, tier 0). Everything that allocates goes through
// exactly these three.
var (
	MemAllocate = Symbol{ModuleMemory, "allocate"} // (size i64) -> ptr, null on failure
	MemFree     = Symbol{ModuleMemory, "free"}     // (p ptr) -> void
	MemResize   = Symbol{ModuleMemory, "resize"}   // (p ptr, size i64) -> ptr
	MemAlignedAllocate = Symbol{ModuleMemory, "allocate_aligned"}
	MemZero     = Symbol{ModuleMemory, "zero"} // (p ptr, len i64) -> void
)

// Reference counting: `shared T` retain/release, `weak T`, and upgrade's
// CAS-against-zero (A.4.8's "a race the type system cannot statically win").
var (
	RCNew     = Symbol{ModuleRC, "new"}     // (payload_size i64) -> ptr
	RCRetain  = Symbol{ModuleRC, "retain"}  // (h ptr) -> ptr
	RCRelease = Symbol{ModuleRC, "release"} // (h ptr, deinit ptr) -> i1  (true: payload died)
	RCWeak    = Symbol{ModuleRC, "weak"}    // (h ptr) -> ptr
	RCWeakDrop = Symbol{ModuleRC, "weak_drop"}
	RCUpgrade = Symbol{ModuleRC, "upgrade"} // (w ptr, out ptr) -> ptr, "" | msg
	RCPayload = Symbol{ModuleRC, "payload"} // (h ptr) -> ptr — control word skipped
)

// unique T: one word, but a bare copy is deep (A.9.4's one cost cliff).
var (
	UniqueNew  = Symbol{ModuleMemory, "unique_new"}
	UniqueFree = Symbol{ModuleMemory, "unique_free"}
)

// []T — growth policy, the sole implicit-allocation exception (A.3.1).
var (
	SliceGrow   = Symbol{ModuleSlice, "grow"}   // (hdr ptr, elem_size i64, want i64) -> void
	SliceAlloc  = Symbol{ModuleSlice, "alloc"}  // (hdr ptr, elem_size i64, n i64) -> void
	SliceFree   = Symbol{ModuleSlice, "free"}   // (hdr ptr) -> void
	SliceOOB    = Symbol{ModuleSlice, "out_of_bounds"} // (i i64, len i64) -> noreturn
)

var (
	MapNew    = Symbol{ModuleMap, "new"}
	MapGet    = Symbol{ModuleMap, "get"}
	MapSet    = Symbol{ModuleMap, "set"}
	MapErase  = Symbol{ModuleMap, "erase"} // A.5.2: nil into a subscript
	MapLen    = Symbol{ModuleMap, "len"}
	MapFree   = Symbol{ModuleMap, "free"}
	MapIterInit = Symbol{ModuleMap, "iter_init"}
	MapIterNext = Symbol{ModuleMap, "iter_next"}
)

var (
	StringCopy    = Symbol{ModuleString, "copy"}
	StringFree    = Symbol{ModuleString, "free"}
	StringConcat  = Symbol{ModuleString, "concat"}
	StringCompare = Symbol{ModuleString, "compare"}
	StringLen     = Symbol{ModuleString, "len"}
	// A.5.6: decode UTF-8 into char scalars at variable stride.
	StringDecode = Symbol{ModuleString, "decode"}
	// A.8.5: NUL-terminated marshalling, manufactured only at a declare
	// boundary.
	StringMarshal   = Symbol{ModuleString, "marshal"}
	StringUnmarshal = Symbol{ModuleString, "unmarshal"}
)

var (
	ChanNew         = Symbol{ModuleChan, "new"}
	ChanSend        = Symbol{ModuleChan, "send"}
	ChanReceive     = Symbol{ModuleChan, "receive"}
	ChanTrySend     = Symbol{ModuleChan, "try_send"}
	ChanTryReceive  = Symbol{ModuleChan, "try_receive"}
	ChanClose       = Symbol{ModuleChan, "close"}
	ChanRetain      = Symbol{ModuleChan, "retain"}
	ChanRelease     = Symbol{ModuleChan, "release"}
	ChanSelect      = Symbol{ModuleChan, "select"}
	ChanAwaitReceive = Symbol{ModuleChan, "await_receive"}
)

// The reactor drives hir-synthesized state machines through the task ABI
// (overview §3): {poll fnptr, drop fnptr, state ptr}.
var (
	ReactorStart    = Symbol{ModuleReactor, "start"}
	ReactorShutdown = Symbol{ModuleReactor, "shutdown"}
	ReactorSpawn    = Symbol{ModuleReactor, "spawn"} // (poll ptr, drop ptr, state ptr) -> ptr
	ReactorDrive    = Symbol{ModuleReactor, "drive"} // run until the given task completes
)

var (
	ThreadSpawn  = Symbol{ModuleThread, "spawn"}
	ThreadJoin   = Symbol{ModuleThread, "join"}
	ThreadDetach = Symbol{ModuleThread, "detach"}
	ThreadYield  = Symbol{ModuleThread, "yield"}
)

// The panic floor (§5.2). Nothing on this path may allocate, since the OOM
// path routes through it.
var (
	Panic     = Symbol{ModulePanic, "panic"}      // (msg ptr, len i64) -> noreturn
	PanicOOM  = Symbol{ModulePanic, "oom"}        // -> noreturn
	PanicTrap = Symbol{ModulePanic, "trap"}       // fmt-free fallback -> noreturn
)

// fmt is normative, not convenient: A.12.2 defines an expected-value literal
// as the exact bytes these produce.
var (
	FmtInt    = Symbol{ModuleFmt, "render_int"}    // %d
	FmtUint   = Symbol{ModuleFmt, "render_uint"}   // %u
	FmtFloat  = Symbol{ModuleFmt, "render_float"}  // %f
	FmtString = Symbol{ModuleFmt, "render_string"} // %s
	FmtBool   = Symbol{ModuleFmt, "render_bool"}   // %d over 1/0
	FmtChar   = Symbol{ModuleFmt, "render_char"}
)

var (
	ConsoleOut = Symbol{ModuleConsole, "write_stdout"}
	ConsoleErr = Symbol{ModuleConsole, "write_stderr"}
)

var (
	ProcessExit = Symbol{ModuleProcess, "exit"}
	ProcessArgv = Symbol{ModuleProcess, "argv"}
	ProcessEnvp = Symbol{ModuleProcess, "envp"}
)