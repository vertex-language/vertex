# analyzer

```
github.com/vertex-language/vertex/analyzer
```

`analyzer` checks a parsed Vertex package (`ast.File`) against the static rules of the grammar annex and produces a `types.Package` plus a `types.Info` side table. It depends on `ast`, `token`, `diag`, and `types`; nothing depends on it back — it is the toolchain's terminal analysis stage, sitting after the parser and before anything that consumes checked types (a printer, a backend).

The package's governing constraint, stated in its own doc comment, is A.0.2: the parser tracks exactly one context-sensitive parameter, and "nothing else in Vertex is context-sensitive." Every other rejected form in A.14's index is a static rule enforced here, over an already-parsed tree that the parser accepted unconditionally.

## Entry points

### `Config`

```go
type Config struct {
    Fset     *token.FileSet
    Tag      token.BuildTag
    Reporter diag.Reporter
    Importer Importer
}
```

What checking a package needs beyond the files themselves. `Importer` may be `nil`, in which case every import is reported as an undeclared name — correct behavior for a single-package check, and the seam a loader fills in later.

### `Importer`

```go
type Importer interface {
    Import(path string) (*types.Package, error)
}
```

What a loader satisfies. Per A.2.3, the imported package's own `PackageClause` supplies the qualifier a caller imports it under, so an implementation must have read the imported directory's package line before it can answer.

### `Checker`

```go
type Checker struct { /* ... */ }
```

Holds one package's checking state: the package under construction, the `types.Info` side table, per-declaration bookkeeping (`objMap`), the innermost open scope during a body walk, and a couple of pieces of transient phase state (see [Context tracking](#context-tracking-npu) below).

### `NewChecker(conf *Config, path, name string, info *types.Info) *Checker`

Prepares a checker for one package. If `info` is `nil` a fresh `*types.Info` is allocated, so a caller that doesn't need the side table needn't construct one.

### `(*Checker) Files(files []*ast.File) (*types.Package, error)`

Runs all three phases (below) over a package's files and returns the resulting `*types.Package`. Errors accumulate rather than aborting the walk; the returned error is non-nil only if `errCount > 0`, and its `Error()` text does not enumerate — diagnostics themselves come through `Config.Reporter`.

### `(*Checker) Package() *types.Package`

Returns the package under construction, valid even after errors. Resolution deliberately produces a partial result on failure, so a caller like an editor's tooling can use whatever did resolve.

## The three phases

`Files` runs three passes in a fixed order, and the order is load-bearing rather than an optimization:

```go
c.collectObjects()   // phase 1: names exist, types are nil
c.resolveDeclTypes()  // phase 2: types filled in, cycles caught
c.checkBodies()       // phase 3: bodies walked, uses recorded
```

A.2 makes top-level declarations order-independent — any declaration may refer to any other regardless of textual position — and that guarantee is what forces the split: every package-scope name must exist before any declaration's type is resolved, and every declaration's type must be known before any body that calls it is walked. Collapsing two of these phases would reintroduce a forward-declaration requirement the language doesn't have.

### Phase 1 — `collectObjects`

Walks every file, inserting each package-scope name into `pkg.Scope()` with a **nil type**. The nil type is the point: it lets every name be visible everywhere before any of them commit to a shape. Imports are bound into a **file-scoped** scope (`collectImports`), since A.2.3 gives an import qualifier file scope while the declarations that use it are package-scoped — one file's imports never resolve a name in another file of the same package.

Methods are collected separately, in `collectMethods`, after every type name exists — a receiver may name a type declared later in the file or in another file, which A.6.3's "methods are declared outside the class body" permits.

### Phase 2 — `resolveDeclTypes`

Iterates `objMap` and calls `objDecl` on each object. Iteration order is unspecified and doesn't matter: `objDecl` is idempotent and resolves dependencies on demand (`obj.Type() != nil` short-circuits), so whichever object is reached first pulls in whatever it needs.

A `resolving` stack guards against self-reference: an object reached a second time while still on the stack is a cycle, reported as `diag.TypeCycle` rather than overflowing. Because A.2's order-independence makes a cycle reachable from any entry point, the guard lives in `objDecl` itself rather than in a separate pre-pass over the dependency graph.

Resolution for one object runs in *that object's declaring file's scope* (`d.fileScope`), so a qualified type resolves through the file that wrote it and no other.

### Phase 3 — `checkBodies`

Walks every function body, resolving identifiers to objects via `expr`/`stmt`/`typ`. This pass computes no expression types — resolution alone (what a name denotes) is enough to settle several of the parser's deferred ambiguities:

- **A.3.6** — `a[i]` vs. `Stack[int32]`: an `IndexExpr` is a type instantiation if its operand denotes a type (`denotesType`), otherwise an index/slice.
- **A.4.4** — `&` as address-of vs. dereference, resolved the same way, downstream of what the operand names.

Scopes opened here (`openScope`/`block`) are recorded into `types.Info` via `RecordScope`, keyed by the `ast.Node` that owns the scope's extent.

## What `types.Info` is for

> Nothing here writes to the syntax tree. Every result lands in a `types.Info` side table.

This is what lets a printer round-trip a checked file unchanged, and what lets the checker run twice over one tree without residue: `Defs`, `Uses`, `Scopes`, `RecordType`, `RecordSelection`, `RecordInstance`, and `Defers` are all populated by side effect as the tree is walked, never by mutating `ast` nodes.

## Errors (`errors.go`)

Every diagnostic funnels through three functions:

- **`report(d *diag.Diagnostic)`** — the only place `errCount` is incremented and the only place `Config.Reporter` is invoked.
- **`errorAt(pos, end token.Pos, code diag.Code, args ...any)`** — constructs and reports in one call.
- **`errorExpr(x ast.Node, code diag.Code, args ...any)`** — reports over a whole node's `Pos()`/`End()` extent, since a type error should underline the construct, not just its first token.

Call sites never format their own message text — `diag`'s registry keeps one rule's wording identical across every site that raises it, which is what makes A.12.2's `Expected(error, "...")` stable enough to appear in the spec itself.

`invalid() types.Type` returns `types.Typ[types.Invalid]` as the universal recovery value: it's returned instead of `nil` so every consumer can keep going without a nil check, and `types`'s predicates treat it as compatible with everything, so one bad type produces exactly one diagnostic rather than a cascade. This mirrors the discipline `ast.BadExpr` follows on the parse side.

## Type expressions (`typexpr.go`)

Two entry points convert an `ast.Expr` written in type position to a `types.Type` or `*types.Constraint`, and the split is deliberate:

- **`typ(e ast.Expr) types.Type`** — ordinary type positions. Also records the result into `types.Info` via `RecordType`.
- **`constraintExpr(e ast.Expr) *types.Constraint`** — `[...]` constraint positions only. A.7.2 makes a constraint "never a value type," and `types.Constraint` deliberately does not implement `types.Type`, so a bare `var c: Ordered` falls out of `typ` as an ordinary `diag.ConstraintAsType` diagnostic rather than needing every caller to remember a predicate check.

A single identifier in a `[...]` position is the fork A.7.2 names explicitly: it "parses as both a `TypeSet` of one term and a `ConstraintName`," and `constraintExpr` resolves it by what the name actually denotes (a lookup), not by shape.

Several checks live here specifically because the parser was deliberately permissive and left the rejection to this phase (A.14's "forms that parse and are rejected here"):

| Form | Rejected because |
|---|---|
| `TensorType` outside an npu body | A.3.5 — grammatical only under `[+Npu]` |
| `abstract` written inline | A.3.3 — legal only as a `TypeAliasDeclaration` target |
| `~T` outside a type set | A.7.3 — `~` only valid inside a constraint |
| stacked ownership (`mut var T`) | A.3.2 — qualifiers don't compose |
| generic name used bare, no type args | A.7.5 — inference works from value arguments, and a type position has none |

## Declarations (`decl.go`)

Holds phase 1 and phase 2 proper: `collectObjects`/`collectImports`/`collectDecl`/`collectMethods` (phase 1), and `objDecl` plus one resolver per declaration kind — `recordDecl`, `enumDecl`, `aliasDecl`, `constraintDecl`, `funcDecl`, `varDecl`, `foreignDecl` (phase 2).

Notable shapes:

- **`recordDecl`** binds the `*types.Named` to its `TypeName` *before* resolving any field, which is what lets a field like `next: typed_ptr Node` reach its own enclosing type without tripping the cycle guard — the guard only fires on an object whose type is still nil, and by the time fields are walked it no longer is.
- **`aliasDecl`** branches on whether the alias target is `abstract`: a transparent alias mints no `*types.Named` and is invisible to `Identical`, while an `abstract` alias mints one keyed on the alias object itself, so two abstract aliases never unify "however identical their provenance."
- **`declare` blocks** (`collectForeign`, `foreignDecl`, `foreignFunc`, `foreignClassMembers`) check A.8's linkage-boundary rules: no bodies, no visibility modifiers, no fields, and family-dependent rules like `abstract → typed_ptr T` being legal only for memory-flat import families (`familyForBlock`).
- **`typeParams`** performs constraint distribution the parser deliberately punted: A.7.1's "a constraint written after a name applies to that name and every immediately preceding unconstrained name" is computed here, walking the list backwards, rather than in the tree — doing it in `ast` would erase the written form a formatter needs to reprint.

## Scopes and resolution (`check.go`)

- **`openScope`/`closeScope`** push and pop `c.scope`, recording the new scope's extent and its owning node into `types.Info` when the node is non-nil.
- **`declare(scope, id, obj)`** inserts a binding, handling three special cases before an ordinary `Scope.Insert`: a blank identifier (`_`) is recorded as a definition but never inserted (A.1.2); a `types.Reserved` name is rejected as `diag.ShadowedBuiltin` (A.1.4's builtins, as opposed to the shadowable predeclared type names); and a genuine collision is reported as `diag.DuplicateDeclaration` with a note pointing at the earlier declaration.
- **`lookup(id)`** resolves outward from the innermost open scope via `Scope.LookupParent`, falling back to package scope if none is open, and records the resulting use.

## Context tracking (`npu`)

A.0.2 names exactly one context parameter the parser itself does not track: `npu`. `Checker.npu` is the field that closes that gap — set on entry to an npu-marked function body or literal (`funcBody`, `funcLit`) and restored on exit via `defer`, so it resets at every function boundary the way A.0.3 requires ("an anonymous closure written inside an async body may not await unless it is itself marked async"). It's consulted exactly once, in `typInternal`'s `TensorType` case, to reject a tensor type reached outside an npu body.