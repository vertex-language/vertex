# testutils/

Compiler conformance tests for the Vertex grammar (Annex A). This suite tests
the **compiler** against the spec — it is not example code, not a library, and
not documentation of how to use Vertex. Every file compiles under `build test`
and exists to pin one rule to one checkable result.

## What "testable" means here

Every `test` function under `build test` has exactly two legal result shapes
(A.12.1, A.12.2):

- `Expected(Type, "rendered value")` — the function must compile, run, and
  return a value whose auto-emitted rendering matches the string exactly.
- `Expected(error)` / `Expected(error, "diagnostic text")` — the function body
  must **fail to compile**, optionally with a pinned diagnostic.
- no result type at all — passes if it compiles and runs without crashing.

**A rule only gets a test file if it can be pinned to one of these three
outcomes.** If a rule has no observable effect on a returned value, no
observable runtime behavior, and produces no compile failure, it cannot be
tested here — full stop.

## How folders are ordered

Folders are ordered by **what a test in them requires in order to run**, not by
Annex A section number. A `build test` file has a compiler, a heap, and a
scheduler available to it, and nothing else:

| Tier | Requires | Folders |
| --- | --- | --- |
| 1 | the compiler only — static layout, no allocation | `01`–`07` |
| 2 | an allocator | `08`, `09` |
| 3 | a scheduler | `10` |
| 4 | a linker with real foreign symbols | none — unreachable |
| 5 | a device target | none — unreachable |

The tier-1/tier-2 line is the one that matters most: everything through `07`
has statically known layout and allocates nothing. `[N]T` is inline storage
(A.3.1) and lives in `06_composite_types/`; `[]T` is the language's sole
implicit heap allocation and lives in `08_heap/`. They are not neighbours,
deliberately.

Read order is the ladder. A failure in an early folder means every later
folder's results are suspect — a broken `as` cast will produce noise in
`08_heap/` that has nothing to do with allocation.

## Layout

```
testutils/
    01_values/          A.1.5, A.4.4 — literals, scalars, char, bool, strings, casts
    02_operators/       A.1.6, A.4.5, A.13 — arithmetic, comparison, logic, precedence
    03_bindings/        A.1.2–A.1.4, A.5.1, A.5.2 — let/var, zero values, assignment
    04_control_flow/    A.5.4–A.5.9 — if, while, for-in, switch, defer, break/continue
    05_functions/       A.4.1, A.5.3, A.6.1 — parameters, returns, closures
    06_composite_types/ A.3.1, A.6.2–A.6.5 — arrays, tuples, structs, classes, enums
    07_generics/        A.7 — type parameters, constraints, type sets, instantiation
    08_heap/            A.3.3, A.4.8 — slices, maps, new/delete/resize, handles
    09_ownership/       A.9 — copy vs transfer, liveness, exclusivity
    10_concurrency/     A.3.5, A.10 — channels, select, thread, async
    11_rejected/        A.14 — every Expected(error) test
```

No subfolders anywhere — one flat file set per folder.

Rejection tests live in `11_rejected/` only, never as ad-hoc `errors/` folders
inside other sections; that keeps the rejection taxonomy in one place instead of
duplicating it per section. It sits last rather than first — by the tier ladder
it requires the least, since it never runs — because its files quote constructs
from every folder below it and are unreadable before them.

There is no `lexical/` folder. Identifier forms, contextual keywords used as
ordinary identifiers, and predeclared-name shadowing are observable only *as a
binding*, so they live in `03_bindings/`; punctuator longest-match (`&&` is
never two `&`, `&+` is never `&` then `+`) is observable only *as an operator*,
so it lives in `02_operators/`. There was never a lexical thing to test on its
own, only lexical rules showing through somewhere else.

## When to add a folder

Split a folder when a file inside it can't be named without an "and." None of
the eleven hit that today. The one to watch is `08_heap/`, which currently
carries containers, raw allocation, and the ownership handles; if it grows past
roughly a dozen files it splits into `08_containers/` and `09_raw_memory/` and
everything below renumbers.

## What has no folder, and why

Not an oversight list — these sit at tiers a `build test` file cannot reach.

- **`PackageClause` / `BuildClause` / import ordering (A.2).** Every `.vs` file
  already needs a valid package clause and `build test` clause just to compile
  at all — there is no version of a test file *without* one to compare against.
  These are preconditions for the harness to run, not properties it can
  observe. No `Expected(...)` shape captures "the package clause was first."

- **Abstract interfaces / `declare` blocks (A.8) — tier 4.** FFI linkage against
  a platform ABI isn't something a `build test` file can exercise in isolation:
  it requires a real foreign symbol on the other end of the linkage boundary.
  Testing it means testing the linker and the platform, not the grammar.

- **Device offload, `gpu`/`npu` (A.11) — tier 5.** Requires an actual device
  target to produce a checkable result; a `build test` file has no device to
  launch against.

- **`QualifiedTypeName` (A.3).** Needs a genuine second imported package, which
  is infrastructure this suite does not have. Same failure mode as import
  ordering.

- **Purely structural rules with no return-value or compile-failure
  signature** — e.g. "top-level declarations are order-independent," "a
  `ContextualKeyword` is an `Identifier` everywhere except its special
  position." Either implied by every other test passing, or a non-event.

If a rule can't be expressed as "this compiles and returns X" or "this fails to
compile, optionally with diagnostic Y," it belongs in the spec's prose, not in a
test file pretending to check something it can't.

## Where the spec is silent

Some rules have a shape Annex A does not pin down — float→int rounding
direction, for instance, where A.4.4 says truncate/extend/int↔float without
naming a mode. Do **not** write a test that guesses. Either restrict the test to
what the spec does state (`3.0 as int32`, not `3.9 as int32`), or leave the case
untested and raise it against the spec. A test file asserting an unspecified
behavior silently promotes one implementation's choice into the conformance
suite.

## Adding a test

1. Decide the outcome shape first — `Expected(Type, "...")` or
   `Expected(error[, "..."])`. If neither fits, the rule doesn't get a file
   here.
2. Ask what the test requires in order to run. That picks the tier, and the tier
   picks the folder. If it allocates, it does not belong in a tier-1 folder.
3. One rule, or one tight cluster of directly related rules, per file.
4. Keep bodies minimal. A test that needs elaborate setup is testing the setup;
   find the smallest program that distinguishes the rule from its absence.