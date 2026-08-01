# Vertex Compiler Architecture Overview

This document covers `vertex`, the compiler front end: source bytes in, an unverified `*vir.Module` out. It also specifies **where the builtins live** — the code that powers heap allocation, reference counting, channels, threading, async, and OS access — because that is the one question the pipeline diagram alone can't answer.

Verification, machine codegen, linking, and device codegen belong to `vvm` (`github.com/vertex-language/vvm`), a separate module this repo imports only at the last step.

The governing constraint through every boundary below is the one `ast`'s doc comment states: keep **shape** and **meaning** apart, and once meaning is decided, keep **decision** and **emission** apart.

---

## 1. The whole pipeline, end to end

```text
                        vertex (this repo)                              vvm
  .vs source
     │ scanner.Scan
     ▼
  token.Token stream
     │ parser.ParseFile / ParseDir
     ▼
  *ast.File / *ast.Package          shape only — ast records no meaning
     │ importer.Load
     │   ├─ parser.ParseDir (per package)
     │   ├─ recurse into imports, post-order
     │   └─ analyzer.Checker.Files
     ▼
  *types.Package + *types.Info      meaning, in side tables
     │ lower/hir
     ▼
  *hir.Program                      monomorphic, ownership-explicit, CFG-shaped,
     │                              builtin calls named
     │ lower/vir
     ▼
  *vir.Module   ────────────────────────────────┐
                                                │
  builtins.Modules(triple, needs)               │   assembled into one build set
    → []*vir.Module, built in memory ───────────┤   by the driver
                                                ▼
                                    importer.ResolveImports
                                    → ir/verify.Verify (each)
                                    → CheckReferences → Rewrite
                                    → cpu/lower/<arch> → object
                                    → objectwriter → linker
                                                ▼
                                          native binary
```

Nothing downstream of `analyzer` mutates `ast` or `types`. Each stage produces a *new* representation rather than annotating an old one — `hir.Program` is a fresh tree, `vir.Module` is a fresh tree. This matches the toolchain's existing habit (`types.Info`, not `ast` mutation) of keeping derived facts out of the thing they were derived from.

---

## 2. Front end: source to checked types

Already documented per-package; this is only their role in the flow.

| Package | Produces | Depends on |
|---|---|---|
| `token` | `Pos`, `FileSet`, `Kind`, `BuildTag` | — |
| `diag` | `*diag.Diagnostic`, `Code` registry | `token` |
| `ast` | `*ast.File` — shape, never meaning | `token` |
| `scanner` | `token.Token` stream | `token`, `diag` |
| `parser` | `*ast.File` / `*ast.Package` | `ast`, `diag`, `scanner`, `token` |
| `types` | `Type`, `Object`, `Info` side tables | `ast` only |
| `analyzer` | `*types.Package` + `*types.Info` | `ast`, `token`, `diag`, `types` |
| `importer` | checked `*Package` graph (`*Result`) | `analyzer`, `ast`, `diag`, `parser`, `token`, `types` |

`importer.Load` is the single entry point a driver calls: root import paths in, every transitively-needed checked package out, in dependency order (`Result.Order`). Everything below consumes that result — a `*types.Package`, its `*types.Info`, and the `*ast.File`s it was checked from.

Two facts from this half matter downstream and are easy to lose:

* **The build tag is not a target triple.** A.2.2 gives `linux | windows | darwin | js | wasm | test` — no architecture. It decides what *parses* (`build test` licenses the `test` marker and `Expected(...)`). The `(arch, os, abi)` triple VIR needs is a driver-level concern, resolved in §6.
* **`build test` is not a platform, and not merely an overlay.** It changes the grammar *and* the unit of compilation: under `build test` the driver compiles one test function at a time, in isolation, into its own binary. It still needs a real platform underneath. See §6.2.

---

## 3. `lower/hir` — where every decision is made

Given the checked package graph and its `*types.Info`, `hir` produces a `*hir.Program`: monomorphic, ownership-explicit, control-flow-flattened, with nothing left that a later phase has to reinterpret.

`lower/hir` does not depend on `vvm` at all, and never sees a target triple.

**Whole-program, not per-package.** Monomorphization is seeded from `main` and walks the call graph, so `hir` consumes the entire checked graph at once and emits one `*hir.Program`. `lower/vir` then emits one `*vir.Module` per originating Vertex package, preserving provenance for the namespace/module mapping in §4.

**Generic monomorphization.** `types.Info.Instances` records type arguments at each instantiation site, but only as the analyzer saw them while checking a generic body *once*, generically — it is not a concrete call graph. `hir` builds its own worklist, seeded from every externally reachable root — `main` and exported functions in a normal build, or, under `build test`, the single test function this build exists for plus its synthesized wrapper — composing substitutions as it descends and memoizing by `(Func, ConcreteTypeArgs)` so a diamond of instantiations is built once. A.7.6 makes non-terminating instantiation a compile error; the worklist carries its own depth guard rather than assuming the analyzer enforced it upstream.

Exported *generic* functions are not roots: there are no concrete type arguments to seed with. Only reachability from `main` (or the one test function) instantiates anything.

**Ownership made explicit.** `types.Info.Transfers` answers move-vs-copy at each owning position (A.4.6, A.9.1), but "copy" is not one operation:

| Type | Bare copy is | Becomes |
|---|---|---|
| scalar, `mut`, non-owning view | a bit copy | inline instructions |
| `string` | may be shared/interned — observably deep | builtin call |
| `[]T` | must genuinely duplicate (header + payload) | builtin call + synthesized element loop |
| `unique T` | deep — walks and duplicates the pointee (A.9.4) | synthesized per-type routine |
| `shared T` | an atomic increment | builtin call |
| struct/class/enum with owning fields | recurses field-by-field, or variant-only for an enum | synthesized per-type routine |

This table is the first place builtins appear, and §5 is about the right-hand column.

**`defer` and implicit `deinit` expansion.** A.5.8 requires every deferred call to run in reverse registration order on *every* exit edge, and VIR has no unwinder — so "every exit edge" becomes literal duplicated code at each `return`/`break`/`continue`/fall-through. Same for A.6.4's implicit per-binding teardown (skipped where `types.Info.Transfers` shows the binding left owning). This is CFG surgery, far easier while control flow is still structured, so it happens here, before flattening.

*Suspend edges:* a `return Pending` from the state-machine split is explicitly decoupled from this pass. A suspend edge is not a scope exit; epilogues bypass these points to preserve live state.

**Control-flow flattening.** Structured statements become VIR's Join Convention shape (§4.3): no phi nodes, values merge by same-name reassignment, every block ends in exactly one terminator. This is what makes `lower/vir` legitimately mechanical afterward. Note VIR's `switch` terminator takes a uniform operand/label list regardless of density — jump-table-vs-compare-chain is `cpu/lower/<arch>`'s decision, not this package's.

**`async`/`await` state-machine desugaring.** `async`-marked functions are rewritten into a `poll` loop plus a synthesized payload enum of suspended states. `await` points become state boundaries yielding `Pending`. A localized liveness pass extracts only variables surviving across an `await` into the enum, to minimize payload size.

* **Composition:** bare `await someAsyncFn(x)` is driven inline in the caller's poll loop, merging state machines without touching the reactor or allocating a channel.
* **Launch prefix bypass:** `async f(a, b)` ignores this transformation and generates the channel handshake that dispatches an ordinary call to the reactor.

Each state machine also gets a synthesized **drop routine**: if a task is cancelled, or its handle dies mid-flight, or the reactor is torn down at exit, the payload enum still holds live owning variables whose `deinit`s must run. That routine is tier 2 (§5.2), alongside the state machine itself.

The state machines are synthesized here. **The thing that drives them is not** — see §5.

**Selector and index resolution.** `types.Info.Selections` already recorded what every `SelectorExpr` denotes and `types.Sizes.Offsetsof` knows field offsets; `hir` turns those into direct `field.ptr`/`index.ptr` and direct calls. There are no vtables to build (A.6.3, A.7.2) — every call, including a generic one post-monomorphization, is direct by construction.

The one exception is the **task ABI**: a portable builtin reactor cannot name user state machines, so it receives `{poll fnptr, drop fnptr, state ptr}` and calls through it via VIR's `fnsig`-typed indirect call. That is an interface between tier 1 and tier 2, not a vtable, and it is the only indirect dispatch the compiler emits.

**Vertex-level variadics.** A.6.1's `...T` lowers to a stack-local fixed array plus a two-word slice view, entirely within `hir`. This is unrelated to VIR's `valist` mechanism, which exists for *reading* heterogeneous C varargs inside a callee. Vertex source never defines such a callee, so `hir` never emits `va_start`. It may still emit a `call` to a variadic extern `fnsig`; only the callee side needs a cursor.

**Globals.** A.2 requires a top-level `VariableDeclaration` to have a compile-time-evaluable initializer, and VIR's `global` init form is narrower still — literal, `zero`, `addr ident`, or an aggregate of those, with no arithmetic and no `const` references. `hir` folds initializers down to those forms; there is no static-initialization-order problem to have because there is no initialization-time code.

**`gpu`/`npu` bodies.** `hir` can represent a device-marked function (the marker and its constraints are enforced by `analyzer`), but there is no `lower/gvir` counterpart and `vvm` lists "no host↔device story" as its largest gap. A device-marked function is out of scope for VIR emission: the launch *call site* validates, but the kernel body has nowhere to go.

### Pass order

1. Monomorphization — everything below is type-dependent, so nothing can run before concrete types exist.
2. `async`/`await` state-machine split — over still-structured control flow, so the split can see scopes.
3. `defer` / `deinit` epilogue expansion — after the split, so each epilogue copy lands in a known state rather than having to be re-partitioned.
4. Ownership expansion — copy/drop calls and per-type routine synthesis.
5. Control-flow flattening.

Step 2 before step 3 is a real choice, not an accident: expanding epilogues first would duplicate them across exit edges that the state split then has to cut apart and reassign.

### What `hir` deliberately does not do

No register allocation, no instruction selection, no ABI classification beyond what VIR's grammar already exposes (`byval[S]`/`sret[S]`), no object files or relocations, no calling-convention register save areas. Anything genuinely target-dependent stays out — including, critically, *which platform's builtins are linked*.

---

## 4. `lower/vir` — the mechanical translator

Walks a `*hir.Program` and emits `*vir.Module`s through `vvm`'s append-only builder API. If this package ever needs a type switch on "is this an owning type" or "was this transferred," that logic belongs in `hir`.

* Struct/enum/const/global declarations translate near-mechanically from `types.Named`/`types.Struct`/`types.Enum` plus `types.Sizes`.
* Function bodies walk `hir.Block`s and emit instructions close to 1:1.
* A Vertex package maps onto VIR's namespace/module split exactly: **`namespace` is the import path, `module` is the `PackageClause` name** — which is precisely A.2.3's "the path is a locator, the declared name is the qualifier," carried through to the linker unchanged.
* Debug `loc` lines come from `token.Pos` resolved through the originating `*token.FileSet`, so a `.vir` diagnostic points back at Vertex source.

**Tier-2 symbol ownership.** Synthesized routines (per-type copy/deinit, monomorphized instantiations, state machines) are emitted into a *canonical owning module*, not into every module that happens to need them: the module declaring the generic function, or for a per-type routine, the module declaring the type. VIR has a strict flat namespace and declare-before-use, so two packages both copying a `[]Foo` would otherwise produce a duplicate symbol rather than a silently-merged one. Consumers reference the routine as an ordinary qualified cross-module call.

Output is unverified; `vvm.BuildModule`/`RunModule` run `ir/verify.Verify` before anything trusts it.

**One rule with teeth:** a module emitted from ordinary Vertex source contains **no `target`, no `link`, and no `extern` group**. It is pure compute. The single exception is a file carrying a `declare` block — and A.8.1 already requires such a file to carry a `BuildClause`, which is exactly the information needed to emit a triple. Interop is the only path by which user source touches the linker, and the grammar already fenced it.

---

## 5. Builtins

VIR has no built-in heap: all built-in allocation is `alloca.ptr`, and heap allocation requires an `extern fn` call. There is no GC, no exception machinery, no support library. Meanwhile Vertex has `unique`, `shared`, `[]T`, `map[K]V`, `chan T`, `thread`, `async`, and `string`. Everything in that list needs code VIR cannot express as instructions.

That code is **builtins** — deliberately not called a runtime, because it isn't a library that ships alongside the program and has no independent existence. It is whatever `vir` modules the compiler has to manufacture to make the language's own features work, and nothing more.

It lives at **`github.com/vertex-language/vertex/builtins`**, and it is **Go that constructs `*vir.Module` values through `vir`'s builder API**. No `.vir` files on disk, no `go:embed`, no text to parse. Anything the language needs — the heap, reference counting, thread creation, the poll reactor, or a raw door out to the OS for something not yet imagined — is added here, in Go, shaped however that feature wants to be shaped: any namespace, any module layout, any set of exported signatures, targeting any platform, all in memory. The attached `memory_*.vir` modules stop being source and become worked examples of what `builtins` is specified to produce — golden test fixtures, not build inputs.

### 5.1 What the builder API buys

```go
func memoryModule(t vir.Target) *vir.Module {
    m := vir.NewModule("memory").SetNamespace("builtins")
    m.SetTarget(t.Arch, t.OS, t.ABI)

    dep, alloc, free := libcNames(t.OS) // "c" / "System" / "kernel32"
    m.DeclareLink(vir.LinkShared, dep)
    g := m.DeclareExternGroup(dep)
    g.DeclareFunction(alloc, []vir.Param{{Name: "size", Type: vir.I64}}, vir.Ptr)
    g.DeclareFunction(free, []vir.Param{{Name: "p", Type: vir.Ptr}}, vir.Void)

    fb := m.DeclareFunction("allocate",
        []vir.Param{{Name: "size", Type: vir.I64}}, vir.Ptr, true)
    p := fb.Call("p", alloc, vir.Ident("size"))
    fb.Return(p)
    return m
}
```

Four things follow, and they are the whole design:

1. **The mangled symbol is identical on every platform** — `_M8builtins6memory8allocate` — because mangling reads namespace and module, and both come from Go constants that don't depend on the target. `SetTarget` is the only per-triple call.
2. **Namespace and module shape are parameters, not carved into a file's first line.** Nothing forces `"builtins"`; the constants in §5.5 are the single place that decides.
3. **Platform selection is a Go switch on the triple**, resolved inside `builtins`, below the driver. No caller can observe which arm ran except through the emitted `target`/`link`/`extern` sections.
4. **Cross-module calls are typed.** A portable module calls down with `fb.CallImported("p", "memory", "allocate", size)` after `m.DeclareImport("memory")` — the same qualified form `vvm`'s importer resolves and `Rewrite` erases into a plain extern symbol before `cpu/lower` ever runs.

The cost is real and worth naming: a body that was ten legible lines of `.vir` becomes forty lines of builder calls, and the builder validates nothing. `builtins` needs a thin helper layer over `FunctionBuilder` for the shapes that recur — early-return-on-null, refcount CAS loop, element loop over a slice — or it will be the least readable package in the repo. §5.5's `--dump-builtins` is the other half of that answer.

### 5.2 Three tiers

| Tier | Target-dependent? | Type-dependent? | Built by |
|---|---|---|---|
| **0 — platform** | yes | no | `builtins`, switching on triple |
| **1 — portable** | no | no | `builtins`, ignoring triple |
| **2 — synthesized** | no | yes | `lower/hir` |

The split is a decision procedure, not a taxonomy: **tier 0 if it's target-dependent and type-agnostic; tier 2 if it's type-dependent and target-agnostic.** Monomorphization is why the second half is forced — you cannot pre-build a copy routine for a type the user hasn't written yet.

**Tier 0 — platform.** The only modules that carry `target`, `link`, and `extern`.

| Module | Backs |
|---|---|
| `memory` | `mem_alloc`/`mem_free`/`mem_resize` — the one allocator door |
| `thread` | `thread f()` — create, join, detach, yield |
| `poll` | reactor backend — epoll / kqueue / IOCP |
| `sync` | futex/condvar primitives under blocking `chan` operations |
| `io` | files, sockets, the raw read/write the std library wraps |
| `console` | stdout/stderr writes for panic and test output |
| `time` | monotonic clock, sleep, timer registration |
| `process` | exit, argv/envp access |

**Tier 1 — portable.** Pure compute: no `target`, no `link`, no `extern`, whichever triple they were built for. They reach the platform only by `import`ing tier 0.

| Module | Backs |
|---|---|
| `rc` | `shared T` retain/release; `weak T`; `upgrade(w)`'s CAS-against-zero (A.4.8's "a race the type system cannot statically win") |
| `slice` | `[]T` growth policy — the sole implicit-allocation exception (A.3.1) |
| `map` | `map[K]V` — hashing, probing, `nil`-assignment erase (A.5.2) |
| `string` | UTF-8 length/compare/concat, `char` decode at variable stride (A.5.6), NUL-terminated marshalling at a `declare` boundary (A.8.5) |
| `chan` | send/receive/trySend/tryReceive/close, refcounted handle (A.10.1) |
| `reactor` | the poll loop that drives `hir`-synthesized state machines, through the task ABI in §3 |
| `panic` | trap-with-message; OOM path for allocation failures that panic rather than return a tuple |
| `fmt` | `%d`/`%u`/`%f`/`%s` rendering — needed by `panic` and normative for A.12.2 |

Note the allocation-failure split the spec already fixes: `new`/`resize` return `(typed_ptr T, string)` and fail politely, while channel allocation *panics* "matching native array allocation" (A.10.1) — so slice growth panics too, and `panic` is required rather than optional.

**The panic floor.** `panic` → `fmt` → `console` (tier 0) is a dependency chain that cannot itself fail, or a stubbed `console` on a new platform would make panic recurse. `panic` therefore has an allocation-free, `fmt`-free fallback path — raw trap plus `process.exit` — used when the rendering path is unavailable. Since the OOM path routes through `panic`, nothing on it may allocate.

**Tier 2 — synthesized.** Emitted by `hir` into the canonical owning module (§4): per-type copy/deinit/drop routines, per-type slice and channel element moves, monomorphized instantiations, `async` state machines with their payload enums and drop routines, `defer`/`deinit` epilogues, and the program entry shim.

The entry shim deserves naming. Vertex's `main` takes no parameters and returns nothing (A.6.1) but sets `[+Await]`, so it may suspend. VIR's entry is `export fn main() i32 entry`, and `vvm`'s `crt` synthesizes the process-entry stub around it. So `hir` synthesizes a wrapper that starts the reactor, drives user `main` to completion, tears down, and returns a status. Under `build test` the same slot holds a wrapper around one test function: start the reactor, drive it, render its result through `fmt` to `console`, tear down, return 0. One seam, two occupants — and both are tier 2.

**Build only what the program needs.** `lower/hir` already knows whether the program contains a `map`, a `chan`, an `async` function, or a `shared T`. That feature set reaches the driver, and `builtins.Modules(triple, needs)` builds only the required modules plus their closure — asking for `chan` implicitly pulls `sync`, `rc`, and `memory`. Hello-world links `memory`, `console`, `panic`, `fmt`, and `process`, and nothing else. This removes the dead-symbol-elimination question rather than answering it.

### 5.3 Builtins ≠ standard library

Two different things, and conflating them is the failure mode this section exists to prevent.

* **Builtins** are `vir` modules the compiler manufactures. Vertex source cannot import them and has no spelling for them. Calls into them are emitted by the compiler, for language constructs the user wrote in ordinary syntax.
* The **standard library** is ordinary Vertex source in ordinary packages, resolved by the ordinary `importer`. Where it needs the OS, it uses a `declare` block (A.8) — the sanctioned user-level path — never a builtin symbol.

An abstraction that shows up in Vertex source is std. An abstraction the *grammar* commits to is a builtin.

### 5.4 The "when" — one row per stage

| Stage | What it does about builtins |
|---|---|
| `parser`, `analyzer` | nothing. Neither knows builtins exist |
| `lower/hir` | **decides** every call: which primitive, with what arguments, at which program point. References symbols through `builtins`' constants, never string literals |
| `lower/vir` | **emits** them: an `import "<module>"` line and `call <module>.<fn>` operands. No decisions |
| driver (§6) | **selects**: passes the resolved triple and the feature set to `builtins.Modules` |
| `builtins` | **constructs** each module in memory through `vir`'s builder, resolving the triple internally |
| `vvm` importer | **resolves** each qualified reference against the real module and rewrites it to a mangled extern symbol |
| `vvm` linker | **links** it, with tier 0's `link` lines feeding `linkdeps.go` |

The platform is unknown to every row above the driver. That is the property worth protecting.

### 5.5 The `builtins` package

`github.com/vertex-language/vertex/builtins` holds three things:

1. **The ABI as constants** — namespace, module names, function names. `lower/hir` writes `builtins.SharedRetain`, never `"rc"`/`"retain"`. One place to grep, one place to break the build when a signature changes, and the single place that decides the namespace spelling.
2. **Module constructors** — one per module; tier 0 takes a `vir.Target`, tier 1 ignores it. This is where you add whatever the language turns out to need: a new allocator, a different refcount discipline, a platform door to something not yet imagined.
3. **`Modules(triple vir.Target, needs FeatureSet) []*vir.Module`** — resolves `needs` to its closure and returns the built set.

`lower/hir` depends on (1) only, and stays platform-blind.

Two testing obligations that a file-based design couldn't have had, and which are now cheap:

* **Signature parity.** Iterate every supported triple, build every tier-0 module, assert the exported signature sets are byte-identical. This turns invariant 5 from a review rule into an assertion.
* **Verification.** The builder validates nothing, so every module, for every triple, goes through `ir/verify.Verify` in a unit test — not only when a real build happens to link it.

Add `--dump-builtins`, routing the built set through `vvm`'s canonical text printer. It's a few lines, it round-trips losslessly, and without it the thing that used to be readable source becomes opaque.

---

## 6. The driver and the build set

### 6.1 The normal build

Nothing above ever resolves a target triple. Something must, and it is the driver's job, in this order:

1. Take the user's requested triple, or default from the host.
2. Derive `token.BuildTag` from its `os` — that is what `parser`/`analyzer` need, and it is coarser than the triple by design.
3. Run `importer.Load` → `lower/hir` → `lower/vir` over the user's package graph.
4. Collect the feature set from `lower/hir` and ask `builtins.Modules(triple, needs)` for the module set.
5. Hand `vvm` the union, with `--root` naming the module holding the entry point.

Two mismatches live at this seam and should be errors with real messages, not surprises:

* **A build tag with no `cpu/lower` backend.** `js` and `wasm` are legal A.2.2 tags, but `vvm`'s `arch` table is real silicon only — bytecode and VM targets are explicitly excluded. These tags parse and check and have nowhere to go today.
* **A triple with no tier-0 support.** Adding an architecture to `vvm` does not add builtins. `builtins.Modules` must fail loudly for an OS it has no `extern` mapping for, rather than emitting a module with a missing symbol.

### 6.2 `build test` — one function, one binary

A `test` function is not a special construct with builtin support. It is a `main` with its return value printed, compiled and run in isolation, judged from outside the process.

**The rule.** Each `test` function gets its own complete compilation and its own binary. Nothing about a test function is visible to another, and no binary contains more than one.

This is forced, not stylistic: `Expected(error)` bodies must fail to type-check. If they were checked alongside their neighbours, one intentional failure would abort the package and take every passing test with it. **Isolation starts at `analyzer`**, not at codegen.

| Form | How far it gets | Verdict comes from |
|---|---|---|
| `Expected(T, "...")` | full pipeline → native binary → exec | exact match of stdout against the literal, plus exit status 0 |
| `test` with no `Expected` | full pipeline → native binary → exec | exit status 0 only; stdout ignored |
| `Expected(error)` | stops at `analyzer` | at least one error diagnostic was produced |
| `Expected(error, "...")` | stops at `analyzer` | an error diagnostic whose `diag.Text()` contains the literal |

A compile-failure test never reaches `lower/hir` and never produces a `.vir` module, an object, or a binary.

**Per-test driver sequence, for a value test:**

1. `parser.ParseDir` the test package once. Parsing succeeds for every test function including the ill-typed ones — they are syntactically valid Vertex.
2. Form a package variant containing package-level declarations plus exactly one `test` function.
3. `importer.Load` the imports (shared across every test in the package — only the variant differs) and run `analyzer` on the variant.
4. If the test expects an error: stop here and judge the diagnostics. Otherwise a diagnostic here is a failure.
5. `lower/hir` with that one function as root, synthesizing the render wrapper into the entry slot.
6. `lower/vir` → `builtins.Modules(triple, needs)` → `vvm` → binary.
7. Exec it. Compare stdout to the literal; compare exit status to 0.

**The render wrapper.** The `Expected` type argument selects the format verb from A.12.2's table (`int32`→`%d`, `float32`→`%f`, `bool`→`%d` over 1/0, and so on), and the wrapper emits a single `fmt` render of the returned value to `console`, with no trailing newline and no decoration. This is why `fmt` is normative rather than convenient: the expected-value literals in source are defined as the exact bytes `fmt` produces, so `Expected(float32, "5.000000")` is a statement about `%f`'s six-digit default and nothing else.

**The three channels.** stdout carries the rendered value and nothing else. stderr carries panic text and any diagnostic output. Exit status carries crash detection — which is the entire mechanism behind a bare `test` function with no `Expected`. Keeping them separate lets the driver distinguish "wrong value" from "panicked" from "exited nonzero" without parsing anything.

**Cost.** N test functions means N analyzer runs, N lowerings, N links. The import graph and its checked packages are computed once and reused; only the test package's variant is re-checked. Isolation is the feature, and it also means a crashed test takes down one process rather than the suite.

The harness is Go, in the driver, outside every process it judges. There is no builtin test module and nothing inside a binary knows it is a test.

---

## 7. Invariants

The rules that keep the above from drifting. Each is checkable.

1. `lower/hir` never imports `vvm` and never sees a target triple.
2. `lower/vir` decides nothing — no type switch on ownership, transfer, or copy depth.
3. A module emitted from non-`declare` Vertex source contains no `target`, `link`, or `extern`.
4. Tier-1 modules contain no `target`, `link`, or `extern` either, whichever triple they were built for. Platform contact is tier 0 or nowhere.
5. Every tier-0 module exports an identical signature set on every supported triple — asserted by test, not by inspection. A platform that can't provide one gets a stub that panics, never a missing symbol.
6. Builtin symbol names and the namespace appear in Go exactly once, in `builtins`.
7. Nothing after `analyzer` mutates `ast` or `types`.
8. A test binary is indistinguishable from an ordinary binary at the VIR level — no test-specific module, symbol, or branch anywhere in the built set.
9. A test binary contains exactly one user `test` function. Nothing enumerates tests inside a process.
10. `lower/hir` imports `builtins`' constants, never its constructors. What decides *which* call to emit may not see *how* the callee is built.
11. Every tier-2 symbol has exactly one owning module. Two packages needing the same synthesized routine reference it; they don't each emit it.

---

## 8. Not yet implemented, and open questions

**Not implemented**

* **Device lowering.** No `lower/gvir`. `gpu`/`npu` bodies type-check with no path to a `.gvir` module, and `vvm` has no host↔device story to receive one.
* **`js` / `wasm` build tags** parse and check with no backend behind them.
* **The `builtins` package.** Tier 0's target shape is proven by the attached `memory_*.vir` fixtures; the Go constructors, the tier-1 set, the feature-set plumbing, and the parity/verify tests are specified here and unwritten.
* **The test driver.** Per-function variant construction, process exec, and comparison.

**Open**

* **`defer` vs. `deinit` epilogue ordering** at a shared exit edge. The spec doesn't resolve it; it should be pinned explicitly rather than picked silently in code.
* **String interning.** A.9.4 licenses sharing or interning `string` payloads since they're immutable. Builtins can take that license or decline it; declining is the safe default and the choice should be recorded, not left to whoever writes the `string` constructor.
* **Where the feature set is computed.** `lower/hir` knows which language features a program uses but has no output channel for it beyond the emitted module. Either it returns a `FeatureSet` alongside the `*hir.Program`, or the driver re-derives it by scanning emitted `import` lines. The second is less code and harder to get wrong; the first is more direct.
* **Diagnostic matching is text-based and therefore brittle.** `Expected(error, "cannot convert string to int32")` pins wording that will drift. `diag`'s `Code` registry is the stable thing to match against; the spec fixes the syntax as a string literal, so either the driver matches codes when the literal parses as one, or the wording becomes normative. Pick one, don't let the first broken test decide.
* **Substring vs. exact, and which diagnostic.** If the analyzer reports three errors, does any match count, or the first by position?
* **A compiler crash is not a passing error test.** `Expected(error)` must be satisfied by a real `diag` error from `parser` or `analyzer`; a panic in `lower/hir` should be reported as a driver failure, never silently counted as the expected failure.
* **stdout contamination.** If a test body writes to stdout itself, it corrupts the comparison. Either that's the user's problem and documented, or the wrapper renders on a separate channel.
* **No optimization tier.** Every `shared` retain and every string copy is a cross-module call with no inlining pass named anywhere in this document. If `vvm` has no mid-level optimizer, that's a deliberate choice worth recording here; if it does, it belongs in §1's diagram.