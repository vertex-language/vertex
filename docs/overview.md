# Vertex Compiler Architecture Overview

`vertex` is the compiler front end: source bytes in, an unverified `*vir.Module` out. Verification, machine codegen, linking, and device codegen belong to `vvm`, a separate module this repo imports only at the last step.

This document is a map, not a specification. Where it touches a rule, the owning document wins:

| Question | Owned by |
|---|---|
| What parses | `grammar.md` |
| What is accepted, and what it means | `semantics.md` |
| What each construct becomes in VIR | `lowering.md` |
| What VIR itself means | `vvm/spec` |

The governing constraint across every boundary is the one `ast`'s doc comment states: keep **shape** and **meaning** apart, and once meaning is decided, keep **decision** and **emission** apart.

---

## 1. The pipeline

```text
  .vs source
     │ scanner.Scan
     ▼
  token.Token stream                shape only; line structure recorded, never interpreted
     │ parser.ParseFile / ParseDir
     ▼
  *ast.File / *ast.Package          shape only; ast records no meaning
     │ importer.Load
     │   ├─ parser.ParseDir per package
     │   ├─ recurse into imports, post-order
     │   └─ analyzer.Checker.Files
     ▼
  *types.Package + *types.Info      meaning, in side tables
     │ lower/hir
     ▼
  *hir.Program                      monomorphic, ownership-explicit, CFG-flattened
     │ lower/vir
     ▼
  *vir.Module per package  ──────►  vvm: verify → resolve imports → lower → link
     ▼
  native binary  +  libvertexrt
```

Each stage produces a *new* representation rather than annotating an old one. `types.Info` is a side table, not a mutation of `ast`; `hir.Program` is a fresh tree; `vir.Module` is a fresh tree.

---

## 2. Front end

Each package is documented in its own README; this is only their role in the flow.

| Package | Produces | Depends on |
|---|---|---|
| `token` | `Pos`, `FileSet`, `Kind`, `BuildTag` | — |
| `diag` | `*Diagnostic`, the `Code` registry | `token` |
| `ast` | `*ast.File` — shape, never meaning | `token` |
| `scanner` | `token.Token` stream | `token`, `diag` |
| `parser` | `*ast.File` / `*ast.Package` | `ast`, `diag`, `scanner`, `token` |
| `types` | `Type`, `Object`, `Info` | `ast`, `token` |
| `analyzer` | `*types.Package` + `*types.Info` | `ast`, `token`, `diag`, `types` |
| `importer` | the checked package graph (`*Result`) | the above |

`importer.Load` is the single entry point a driver calls: root import paths in, every transitively-needed checked package out in dependency order (`Result.Order`).

Two facts from this half matter downstream and are easy to lose:

* **A build tag is not a target triple.** `grammar.md` gives `linux | windows | darwin | js | wasm | test` — no architecture. It decides what is *selected* and what *parses*. The `(arch, os, abi)` triple is a driver concern (§5).
* **`build test` changes the unit of compilation**, not just the grammar. See §6.

---

## 3. `lower/hir` — where every decision is made

Given the checked graph and its `*types.Info`, `hir` produces a `*hir.Program`: monomorphic, ownership-explicit, control-flow-flattened, with every runtime call named. It never imports `vvm` and never sees a target triple. Layout is the one target-shaped fact it consumes, and it arrives as a `types.Sizes` on `Config`, chosen by the driver from the build tag.

Whole-program, not per-package: monomorphization is seeded from `main` (or the one test function) and walks the call graph, so `hir` consumes the entire graph at once. `types.Info.Instances` is *not* a concrete call graph — it records what the analyzer saw while checking each generic body once, generically — so `hir` builds its own worklist, memoized by `(Func, ConcreteTypeArgs)`, with its own depth guard.

What it decides, each specified in `lowering.md`:

| Decision | Read from | Spec |
|---|---|---|
| move vs. copy, and what "copy" costs | `Info.Transfers` | §7.1–7.2 |
| teardown, `defer`, `deinit` on every exit edge | `Info.Defers`, scopes | §7.3 |
| field and element addressing | `Info.Selections`, `Sizes.Offsetsof` | §10.4, §16 |
| `async` frames and resume points | markers, liveness | §18 |
| which `rt_*` call, with what arguments | the construct | §2.2 |

Pass order, and the reasons it is an order rather than a set:

1. **Monomorphization** — everything below is type-dependent.
2. **`async`/`await` state-machine split** — over still-structured control flow, so it can see scopes.
3. **`defer`/`deinit` epilogue expansion** — after the split, so each epilogue copy lands in a known state. A suspend edge is not a scope exit; epilogues bypass those points.
4. **Ownership expansion** — copy/drop calls and `_V*` routine synthesis.
5. **Control-flow flattening** — into VIR's Join Convention shape.

What it deliberately does not do: register allocation, instruction selection, ABI classification beyond `byval[S]`/`sret[S]`, object files, or anything that would need a triple.

---

## 4. `lower/vir` — the mechanical translator

Walks a `*hir.Program` and emits one `*vir.Module` per originating Vertex package through `vvm`'s append-only builder API. If this package ever needs a type switch on "is this owning" or "was this transferred," that logic belongs in `hir`.

Three kinds of work and nothing else: **translation** (an explicit opcode table, not name-based — `hir` spells `field.ptr`, vir spells `field`), **ordering** (vir fixes a section order and forbids forward references; `hir`'s slices are in construction order, so structs and functions are topologically re-sorted here), and **naming** (`namespace` is the import path, `module` is the `PackageClause` name — `semantics.md` §1.3's "the path is a locator, the declared name is the qualifier," carried to the linker unchanged).

Output is unverified. `ir/verify.Verify` is `vvm`'s to run.

---

## 5. The runtime and the driver

VIR has no heap, no unwinder, and no support library, while Vertex has `unique`, `shared`, `[]T`, `map[K]V`, `chan T`, `string`, threads, and async. The gap is **`vertexrt`**, a static library built from Vertex and VIR, reached through exactly one `link` line and one `extern` group. `lowering.md` §2.2 is the complete symbol list; nothing else in this repo may add to it silently.

Three consequences of that list are worth naming here because they constrain the source language:

* **Free is unsized.** `delete[T](p)` carries no count, so `rt` must recover a block's size class from its address. Segment metadata, not a per-block header.
* **`rt_alloc` returns null.** Container paths panic on null; explicit `new[T]` builds a boundary tuple.
* **Refcounting is not in `rt`.** `shared`/`weak` retain and release are inline VIR atomics.

The driver resolves the triple — nothing above it does:

1. Take the requested triple, or default from the host.
2. Derive `token.BuildTag` from its `os`. Coarser than the triple, by design.
3. `importer.Load` → `lower/hir` → `lower/vir`.
4. Hand `vvm` the modules with `--root` naming the entry module, and `vertexrt` on the link line.

Two mismatches live at this seam and should be real errors: a build tag with no `cpu/lower` backend (`js`, `wasm` parse and check with nowhere to go — `lowering.md` §22.5), and a triple `vertexrt` has no build for.

---

## 6. `build test` — one function, one binary

A `test` function is a `main` with its result printed, compiled and run in isolation, judged from outside the process. Each gets its own complete compilation and its own binary.

This is forced rather than stylistic: `Expected(error)` bodies must *fail* to type-check, so checking them alongside their neighbours would let one intentional failure take every passing test with it. **Isolation starts at `analyzer`**, not at codegen.

| Form | How far it gets | Verdict from |
|---|---|---|
| `Expected(T, "…")` | full pipeline → binary → exec | stdout matches the literal exactly, exit 0 |
| `test`, no `Expected` | full pipeline → binary → exec | exit status only |
| `Expected(error)` | stops at `analyzer` | at least one error diagnostic |
| `Expected(error, "…")` | stops at `analyzer` | a diagnostic whose `diag.Text()` matches |

The import graph and its checked packages are computed once and reused; only the test package's single-function variant is re-checked per test. The harness is Go, in the driver, outside every process it judges. Nothing inside a binary knows it is a test.

---

## 7. Invariants

Each is checkable, most by reading emitted text.

1. `lower/hir` never imports `vvm` and never sees a target triple.
2. `lower/vir` decides nothing — no type switch on ownership, transfer, or copy depth.
3. Every emitted module carries `link static "vertexrt"` and therefore a `target` line. A `declare` block adds its own library; nothing else adds a `link`, `extern`, or `target`.
4. `rt` symbol names appear in Go exactly once, in one constants file. `lower/hir` names them; it never spells them as string literals.
5. Nothing after `analyzer` mutates `ast` or `types`.
6. Monomorphized instances are internal, never exported — VIR has no `linkonce`, so two modules instantiating the same generic each emit their own copy.
7. A test binary is indistinguishable from an ordinary binary at the VIR level, and contains exactly one user `test` function.
8. The `_V` prefix is reserved and may not be declared in source.

---

## 8. Status

**Not implemented:** device lowering (`gpu`/`npu` bodies type-check with no `lower/gvir` and no host↔device story in `vvm`); `js`/`wasm` backends; closures beyond the non-capturing form; payload enums beyond the tag; `map`/`string` iteration, `map` literals, `resize`, `fallthrough`; the whole `async`/`await` split, `select`, and the launch prefixes; the `vertexrt` library itself; the test driver.

**Known deviation:** `lower/hir` performs epilogue expansion during body lowering rather than as a pass over the structured tree. That is equivalent for a program with no async functions and wrong for one that has them, because a suspend edge is not a scope exit. Landing async means lifting it out first — three methods on `funcBuilder`.

**Open, and deliberately not decided in code:** `defer` vs. `deinit` ordering at a shared exit edge; whether diagnostic matching in `Expected(error, "…")` pins wording or `diag.Code`; substring vs. exact, and which diagnostic when several are reported; whether a compiler crash can ever satisfy an error test (it must not); stdout contamination from a test body. `lowering.md` §22 carries the items that need a change in `vvm` or in the language documents rather than here.