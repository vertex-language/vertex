# lower/hir

```
github.com/vertex-language/vertex/lower/hir
```

`hir` is where every decision is made. Given a checked package graph and its
`*types.Info`, it produces a `*hir.Program`: monomorphic, ownership-explicit,
control-flow-flattened, with builtin calls named and nothing left that a later
phase has to reinterpret. `lower/vir` is mechanical afterward by construction.

It depends on `ast`, `token`, `types`, and `builtins` — the last for constants
only. It does **not** depend on `vvm`, and it never sees a target triple
(invariants 1 and 10). Layout is the one target-shaped fact it consumes, and it
arrives as a `types.Sizes` on `Config`, chosen by the driver from the build tag.

## Entry point

```go
func Lower(conf *Config, units []*Unit) (*Program, error)
```

`units` are the checked packages in topological order — `importer.Result.Order`
already is one. `hir` takes this shape rather than an `*importer.Result` so it
doesn't depend on the loader, and so a test can hand it a package built by hand.

`Program.Features` is this package's answer to the overview's open question
about where the feature set is computed: `hir` returns it alongside the program
rather than making the driver re-derive it by scanning emitted `import` lines.
Every emitted builtin call records its feature at the call site, so the set can
never disagree with the calls actually emitted.

## Three representation rules

1. **Aggregates are pointers.** vir makes `struct` and `array` memory-only, so a
   `Value` whose `Type` is aggregate *is* a `ptr` at the vir level. Aggregate
   parameters carry `ByVal`, aggregate results carry `SRet`, aggregate
   assignment is a `memcopy`. Stated once, in the package doc comment;
   `lower/vir` reads it off `IsAggregate`.
2. **`let` is a value, `var` is a slot.** A.5.1 already says a `let` "may be a
   register, an SSA value, or folded away entirely" while a `var` "owns a real
   stack slot for its whole lifetime." Making that literal sidesteps every
   Join-Convention definite-assignment subtlety around mutation across branches,
   and it's why only a `var` can be passed to a `mut` parameter.
3. **One instruction representation, two control-flow shapes.** A `Func` carries
   structured `Body` until `Flatten` runs and flat `Blocks` afterward, but both
   hold the same `*Instr` values. `Flatten` moves instructions; it never
   rewrites them.

## Pass order

1. Declarations — shapes, constants, globals, `declare` blocks.
2. Monomorphization, seeded from `main` (or the one test function under
   `ModeTest`). Exported *generic* functions are not roots: there are no
   concrete type arguments to seed with.
3. `async`/`await` state-machine split — **not implemented**, see `async.go`.
4. `defer`/`deinit` epilogue expansion and ownership expansion.
5. Control-flow flattening.

**Known deviation.** The overview specifies steps 4 and 5 as standalone passes
running after step 3. This package performs step 4 *during* body lowering,
through the builder's scope stack (`scope`, `unwind`, `epilogue`). That is
equivalent for a program containing no async functions and wrong for one that
does, because a suspend edge is not a scope exit. Landing `async.go` means first
lifting epilogue expansion out of `stmt.go` into a pass over the structured tree.
The seam is deliberately narrow — three methods on `funcBuilder` — so the lift
is mechanical. This is recorded here rather than discovered later.

## What is not implemented

Each of these produces a `todo:`-prefixed error naming the construct, distinct
from an `internal:` error, which means something the analyzer should already
have caught got through. That's the same distinction vvm's lowering backends
draw.

- **`async`/`await`.** The whole state-machine split, plus `select`, plus the
  `thread`/`async` launch prefixes and their channel handshake.
- **`gpu`/`npu` bodies.** There is no `lower/gvir` counterpart, and vvm lists
  "no host↔device story" as its largest gap. The launch *call site* validates;
  the kernel body has nowhere to go.
- **Closures.** A `FuncLit` needs an env struct and the two-word `{code, env}`
  representation; only the non-capturing one-word form crosses a boundary.
- **Payload enums**, beyond the tag: variant construction, the payload view a
  `case` binding produces, and the tag-driven copy/deinit recursion.
- **`map` and `string` iteration**, `map` literals and subscripts, `resize`,
  and `fallthrough` (which needs the following clause's body duplicated, since
  vir's switch cases do not fall through).

## Open questions this package pinned rather than left

- **String interning.** A.9.4 licenses sharing or interning a `string` payload
  since it's immutable. This package *declines* the license: `own.go` emits a
  real duplication through `builtins.StringCopy`. Recorded, not left to whoever
  writes the `string` constructor.
- **`defer` vs. `deinit` ordering at a shared exit edge.** `epilogue` runs
  deferred calls first, in reverse registration order, then per-binding
  teardown in reverse declaration order. The spec doesn't resolve it; this is
  the pin.

Still open and *not* pinned: whether a cross-module call passing a fat type by
value needs the header struct to be a shared imported shape rather than one
declared per module. Structs produce no vir symbol, so duplicating the
declaration is safe for layout — but `byval`/`sret` comparison is nominal
per-origin, and that interaction wants a decision before the first
cross-package `[]T` argument.