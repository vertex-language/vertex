# lower/vir

```
github.com/vertex-language/vertex/lower/vir
```

The mechanical translator. A `*hir.Program` in, one `*vir.Module` per
originating Vertex package out, emitted through `vvm`'s append-only builder
API. Every decision was already made in `lower/hir`; if this package ever
needs a type switch on "is this an owning type" or "was this transferred,"
that logic belongs upstream (invariant 2).

```go
func Lower(conf *Config, prog *hir.Program) ([]*vir.Module, error)
func Root(prog *hir.Program) string   // the module holding the entry point, for --root
```

Modules come back in `hir`'s topological order, which is the order the
`vvm` importer expects. Output is unverified: the builder API validates
nothing and `ir/verify.Verify` is `vvm`'s to run.

## What "mechanical" turned out to mean

Three kinds of work, and nothing else.

**Translation.** `hir.Op` is a subset of vir's §4 opcodes spelled
identically, so `instr.go` is a table. It is an explicit table rather than
`vir.ParseOpcode(op.String())` because `hir` spells its pointer arithmetic
`field.ptr` / `index.ptr` — the instruction form, suffix included — while
vir's closed vocabulary spells the opcodes `field` / `index`. A name-based
bridge works for sixty opcodes and fails silently for two.

**Ordering.** vir fixes a mandatory section order (§2.1) and has no forward
references (§2.2). `hir`'s slices are in construction order, which is
neither, so both the struct section and the function section are
topologically re-sorted here. See below — this is the only part of the
package with any teeth.

**Naming.** `namespace` is the import path, `module` is the
`PackageClause` name: A.2.3's "the path is a locator, the declared name is
the qualifier," carried through to the linker unchanged.

## Ordering, concretely

- **Structs.** `hir` appends a `*Struct` to its module *before* lowering
  the struct's fields, so a self-referential field can reach its enclosing
  type without looping. That means a struct-typed field's shape lands in
  the slice after the struct naming it. `orderStructs` does a post-order
  walk over by-value field types; a pointer field erases the edge, so a
  cycle would be a struct of infinite size and is reported as a bug.
- **Functions.** `hir`'s worklist appends a shell when a call site is
  *discovered*, so its `Funcs` slice runs caller-before-callee — exactly
  backwards. `orderFuncs` reverses it. §2.2 exempts *direct*
  self-recursion and nothing else, so **mutual recursion has no vir
  spelling at all**; it is reported here, where the cycle is known, rather
  than reaching `ir/verify` as an undeclared-name error. This is a real
  language-level limitation, not a gap in this package.

## The one target-shaped fact

`hir` never sees a triple, and neither does most of this package. §2.1
makes `target` mandatory whenever `link` is present, and invariant 3 says
only a `declare`-block module has links — so `Config.Target` is written
into exactly those modules and no others. **The invariant is checkable by
reading the emitted text:** a module with no `link` line has no `target`
line.

Note that a `BuildClause` is a tag (`linux | windows | darwin | ...`), not
a triple, so the overview's "A.8.1 already requires such a file to carry a
BuildClause, which is exactly the information needed to emit a triple" is
half true. The OS comes from the file; the arch and ABI come from the
driver, through `Config`.

## Two suffix rules worth knowing

- **Comparisons name their operand type, not their result type.** vir's
  `eq`/`slt`/`lt` family yields `i1` from a suffix naming what is being
  compared (`eq.ptr`, `slt.i32`), while `hir.Instr.Type` holds the result.
  The suffix is read off `Args[0]` for these.
- **A void-typed instruction has no suffix.** `memcopy dst, src, n` takes
  none; `store.i32 p, v` does, and `hir` passes the real type there.

## Not emitted

- **`fnsig` declarations, and therefore indirect calls.** Nothing upstream
  produces one — `hir`'s `resolveCallee` `todo`s on a call through a
  function value. The task ABI's `{poll, drop}` function pointers will need
  this: they are the one indirect dispatch the compiler emits, and landing
  them means `hir` naming a signature this package can declare.
- **`export struct`.** `hir` declares synthesized headers per-module, and a
  struct produces no symbol, so nothing is shared today. The open question
  in `lower/hir`'s README — whether a cross-module `byval` fat type needs a
  single imported shape, since `byval`/`sret` comparison is nominal
  per-origin — is decided *here* when it is decided, and this package will
  need `Exported()` and an `Import` field on the type at that point.
- **`extern_c`, `tls`, `inline`/`noinline`/`cold`, `readonly`.** No `hir`
  field carries them yet. `entry` and `noreturn` are wired.

## Two upstream problems this package will not paper over

Both are reported or emitted as written, because papering over either would
be a decision, and this package makes none.

1. **Float division has no vir opcode.** `hir.binaryOp` returns `OpSDiv`
   for a float operand, commented "vir spells float division with the same
   mnemonic family" — but vir's `opTable` gives `sdiv` `ConstraintInt`, and
   §4.1 lists no float division at all. `sdiv.f32` is emitted faithfully
   and `ir/verify` will reject it under §9.18. Either vir needs an `fdiv`
   opcode or `hir` needs to stop claiming one exists; the fix is not here.

2. **`addr` in a global initializer names a later function.** §6.2 permits
   `addr ident` for "earlier functions/globals", but §2.1 puts the whole
   `global` section before the whole `fn` section, so a global holding a
   function address — which is `hir`'s only spelling for one, since vir has
   no address-of instruction — can never reference an *earlier* function.
   The declaration is emitted as `hir` built it. Either the spec means
   "resolved by relocation, order irrelevant" and §6.2's wording is loose,
   or function addresses need a form that doesn't exist yet.

## File organization

| File | Purpose |
| --- | --- |
| `vir.go` | `Config`, `Lower`, `Root`, and the §2.1 section order — the one guarantee the top level makes. |
| `decl.go` | Structs (with the dependency sort), consts, globals and their init forms, links/externs, imports, params. |
| `func.go` | The callee-first function sort, signatures, blocks, terminators, and `loc` lines. |
| `instr.go` | The opcode table, the suffix rules, and the four instructions with dedicated builder entry points. |
| `value.go` | `hir.Value` → `vir.Operand`, `hir.Type` → `vir.Type`. |