# lower/hir

```go
import "github.com/vertex-language/vertex/lower/hir"
```

`hir` is where every decision is made. Given the checked package graph and its
`*types.Info`, it produces a `*hir.Program`: monomorphic, ownership-explicit,
control-flow-flattened, with every builtin call named and nothing left that a
later phase has to reinterpret. `lower/vir` is mechanical afterward by
construction — if it ever needs a type switch on "is this owning" or "was this
transferred", that logic belongs here.

It depends on `ast`, `builtins`, `token`, and `types`. It does **not** depend on
`vvm`, and it never sees a target triple. Layout is the one target-shaped fact it
consumes, and it arrives as a `*types.Sizes` on `Config`. It imports `builtins`
for `names.go`'s ABI constants and `features.go`'s `FeatureSet` only — never a
module constructor, so what decides *which* call to emit cannot see *how* the
callee is built.

`lowering.md` is normative for everything below. Where this README describes a
shape, that document describes the rule.

## Usage

```go
// Result.Order is already topological. hir takes []*Unit rather than an
// *importer.Result so it does not depend on the loader, and so a test can
// hand it a graph built by hand.
units := make([]*hir.Unit, 0, len(res.Order))
for _, p := range res.Order {
	units = append(units, &hir.Unit{Pkg: p.Types, Info: p.Info, Files: p.Files})
}

prog, err := hir.Lower(&hir.Config{
	Fset:  fset,
	Sizes: types.SizesFor(target),
	Mode:  hir.ModeProgram,
}, units)
if err != nil {
	var todo *hir.TodoError
	if errors.As(err, &todo) {
		return fmt.Errorf("unsupported construct: %w", err)
	}
	return err // internal: something the analyzer should have caught
}

mods, err := virlower.Lower(&virlower.Config{Target: tgt}, prog)
```

## Entry point

```go
func Lower(conf *Config, units []*Unit) (*Program, error)

type Config struct {
	Fset  *token.FileSet // for loc lines; hir resolves no positions itself
	Sizes *types.Sizes   // the only target-shaped input
	Mode  Mode           // ModeProgram or ModeTest
	Test  string         // the one test function, under ModeTest
}

type Unit struct {
	Pkg   *types.Package
	Info  *types.Info
	Files []*ast.File
}
```

`units` must be in topological order — dependencies before dependents. `Sizes`
must match the tag the packages were *checked* under: `int` and `uint` are the
target's pointer width (`semantics.md` §2.3), so a mismatch silently changes
every layout the checker already committed to.

`ModeProgram` seeds from `main`; every exported non-generic function is also a
root. `ModeTest` seeds from the single function named by `Test` and synthesizes a
render wrapper into the entry slot instead. See `overview.md` §6 for why a test
binary holds exactly one test function.

## Errors

Two shapes, deliberately distinct — the same split `vvm`'s lowering backends draw.

| Prefix | Means |
|---|---|
| `todo:` | the program is valid; this package does not lower that construct yet |
| `internal:` | something the analyzer should already have rejected got through |

An `internal:` error is a bug in `analyzer` or here, never in the user's source.
Do not report one as a compile error. Both unwind out of the recursive expression
walk by panicking with a private `bail`; nothing else in this package panics.

## Three representation rules

1. **Aggregates are pointers.** vir makes `struct` and `array` memory-only, so a
   `Value` whose `Type` is aggregate *is* a `ptr` at the vir level. Aggregate
   parameters carry `ByVal`, aggregate results carry `SRet`, aggregate assignment
   is a `memcopy`. `vector[T,N]` and the lane predicate are the exceptions —
   register-class types, passed and returned by value (`lowering.md` §4.3, §15).
2. **`let` is a value, `var` is a slot.** `semantics.md` §6.1 already says a `let`
   may be a register, an SSA value, or folded away entirely. Making that literal
   sidesteps every Join-Convention definite-assignment subtlety around mutation
   across branches, and it is why only a `var` can reach a `mut` parameter.
   `lowering.md` §5 lists the six things that force a slot anyway.
3. **One instruction representation, two control-flow shapes.** A `Func` carries
   structured `Body` until `Flatten` runs and flat `Blocks` afterward, but both
   hold the same `*Instr` values. `Flatten` moves instructions; it never rewrites
   them.

## Pass order

1. Declarations — shapes, constants, globals, `declare` blocks. Function bodies
   are deliberately absent: a body is lowered only when monomorphization reaches
   it.
2. Monomorphization, seeded from the roots, memoized by
   `(Func, ConcreteTypeArgs)` with its own depth guard. `types.Info.Instances` is
   not a concrete call graph — it records what the analyzer saw checking each
   generic body once, generically — so the worklist is built here. Exported
   *generic* functions are not roots: there are no concrete type arguments to
   seed with.
3. `async`/`await` state-machine split — **not implemented**, see `async.go`.
4. `defer`/`deinit` epilogue expansion and ownership expansion.
5. Control-flow flattening.

**Known deviation.** `lowering.md` specifies steps 4 and 5 as passes running
after step 3. This package performs step 4 *during* body lowering, through
`funcBuilder`'s scope stack (`openScope`, `epilogueTo`, `closeScope`). That is
equivalent for a program containing no async functions and wrong for one that
does, because a suspend edge is not a scope exit (`lowering.md` §18.2): a
`return Pending` must bypass the epilogue entirely, and an epilogue already
inlined at every exit edge cannot tell the two kinds of edge apart. Landing
`async.go` means first lifting epilogue expansion out of `stmt.go` into a pass
over the structured tree. The seam is deliberately narrow — three methods on
`funcBuilder` — so the lift is mechanical. Recorded here rather than discovered
later.

## Builtin symbols and features

Every runtime name is a `builtins.Symbol`, and `hir` references it that way —
never a `"module"`/`"func"` string pair. `lowering.md` §2.2 is the complete list;
a symbol not on it does not exist. `lower/vir` emits the `link` line and the
`extern` group from what was actually called, so the emitted text and the calls
can never disagree.

`Program.Features` is computed the same way: each emitted builtin call records
its feature at the call site, so the set cannot disagree with the calls actually
emitted.

## Naming

The `_V` prefix is reserved for compiler-synthesized names — the entry shim
(`_Ventry`), synthesized headers (`_Vstr`, `_Vvec`, `_Vfn`), tuple structs, and
interned string globals. `semantics.md`'s list of what may not be declared needs
`_V` added; that is one line and is not yet there.

Nothing here pre-mangles. vir applies its own Itanium-style mangling on top, so
these are idents, not symbols.

## What is not implemented

Each produces a `todo:` error naming the construct.

- **`async`/`await`.** The state-machine split, `select`, and the `thread`/`async`
  launch prefixes with their channel handshake.
- **`gpu`/`npu` bodies.** There is no `lower/gvir` counterpart and `vvm` has no
  host↔device story. The launch *call site* validates; the kernel body has
  nowhere to go.
- **Closures**, a function used as a value, and every indirect call. All three
  want the `{code, env}` pair and an `fnsig`, which nothing upstream produces —
  the same missing piece the task ABI needs.
- **Payload enums** beyond the tag: variant construction, the payload view a
  `case` binding produces, and the tag-driven copy/deinit recursion.
- **Iteration** over an array or slice (including the consuming form), a map, or
  a string. Only the range form lowers.
- **`map` literals, channel constructors, slice literals**, and `fallthrough`
  (which needs the following clause's body duplicated, since vir's switch cases
  do not fall through).
- **Builtins** `min`/`max`/`clamp`, `resize`, `blend`.
- **Foreign class members.** Objective-C needs `lowering.md` §20.4's selector
  cache; C++/COM need fnsig-typed indirect calls.
- **Tensor element types** (`bf16`, `fp8e4m3`, `fp8e5m2`, `int4`), which exist
  only packed inside a tensor, which exists only in an npu body.

## Known gaps that produce no error

These lower to something. That something is wrong, and they are recorded here
because nothing in the build will say so.

- **Generic substitution is bare-parameter-only.** `subst.apply` substitutes a
  `*types.TypeParam` in operand position and returns any composite mentioning one
  unchanged, so `func f[T](xs: []T)` mis-lowers. Covers every generic in the
  corpus today; wants a structural rebuild.
- **A slice's deep copy duplicates only the header.** Right for `[]int32`, wrong
  for `[]string` — the per-element loop at each element's own cost is not there.
- **`RCRelease` is passed a null drop routine**, so a shared payload's teardown
  does not run when the count reaches zero. It wants the per-type `_Vdrop_T`
  below.
- **A top-level `var` initializes to zero.** `types.Var` does not carry the
  folded value; the checker discards it after the representability check.
- **An omitted defaulted struct field gets the zero value**, not the default
  expression, which sits unevaluated on the `ast.Field`.
- **A failed `new` leaves its message slot empty**, pending field indices the
  boundary-tuple shape does not name yet.
- **Float division emits `sdiv`/`udiv`.** vir §4.1 lists no float division at
  all, so `sdiv.f32` is emitted faithfully and `ir/verify` rejects it. Either vir
  needs `fdiv` or this needs to stop claiming one exists; papering over it is not
  this package's decision.

Correct but unoptimized, for completeness: every `var` gets a slot (narrowing
`forcesSlot` wants a use-collecting pre-pass), and `copyInto`/`dropInPlace`
inline the walk at every site instead of calling one `_Vcopy_T`/`_Vdrop_T` per
type — which is mechanical from here, except that the owning-module rule (which
module holds `_Vdrop_Foo` when two packages both drop a `Foo`) wants deciding
first.

## Pinned here rather than left open

- **String interning.** `semantics.md` §8.5 licenses sharing or interning a
  `string` payload since it is immutable. This package declines the license:
  `own.go` emits a real duplication through the `StringCopy` builtin, matching
  `lowering.md` §16.1. A refcount header would put atomics on a type whose
  spelling contains no `shared` — a cost invisible in the source.
- **`defer` vs. `deinit` ordering at a shared exit edge.** `epilogueTo` runs
  deferred calls first, in reverse registration order, then per-binding teardown
  in reverse declaration order. Neither source document resolves it; this is the
  pin.
- **Enum layout.** Tag plus opaque payload bytes, matching
  `types.Sizes.enumSize` so the two cannot disagree. The sources fix no layout.
- **Struct field order is never reordered.** `semantics.md` §7.2 makes order
  observable through destruction order, and interop assumes it.