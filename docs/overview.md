# Vertex Compiler — Architecture Overview

How the packages fit together, and what remains before `lower/` can hand a module to [`vvm`](https://github.com/vertex-language/vvm).

The organizing principle is A.0.2: nothing in Vertex is context-sensitive except one parse-time parameter; every other rejection is a static rule checked over an already-parsed tree. Everything before that line is shape; everything after is meaning. The packages split on exactly that line.

---

## Pipeline

```
.vs source
  │ scanner.Scan
token.Token stream
  │ parser.ParseFile / ParseDir
*ast.File / *ast.Package        ← shape only
  │ importer.Load
import graph, topological order
  │ analyzer.Checker.Files
*types.Package + *types.Info    ← meaning, in side tables
  │ lower.Lower                                   [NOT WRITTEN]
*vir.Module
  │ vvm.BuildModule
native binary
```

`token` and `diag` sit under every stage. Nothing depends on `parser`; only `parser` depends on `scanner`.

---

## Dependency rules

| Package | May import | Must never import |
| --- | --- | --- |
| `token` | — | anything in this repo |
| `diag` | `token` | `ast`, `scanner`, `parser` |
| `scanner` | `token`, `diag` | `ast`, `parser` |
| `ast` | `token` | `types`, `diag`, `parser` |
| `types` | `token`, `ast` | `analyzer`, `diag` |
| `parser` | `token`, `diag`, `scanner`, `ast` | `types`, `analyzer` |
| `analyzer` | all of the above | `importer`, `lower` |
| `importer` | all of the above | `lower` |
| `lower` | all of the above | — |

Two need justifying:

- **`ast` must never import `types`.** The tree records shape, never meaning. A node holding a `types.Type` breaks that and the printer's round-trip property with it. `types` imports `ast` for `Info`'s map keys; the arrow runs one way.
- **`types` must never import `diag`.** `types` answers questions, it doesn't report. `IsComparable` returns a bool; the analyzer decides whether that bool is an error *here*. A reporter inside `types` would make the predicate unaskable speculatively.

---

## Packages

### `token`

Lexical vocabulary and source positions.

| File | Contents |
| --- | --- |
| `pos.go` | `Pos`, `Position`, `File`, `FileSet` |
| `token.go` | `Token`, `Ctx*` spellings, `IsCtx`, `IsBlank` |
| `kind.go` | `Kind`, spelling table, `Lookup`, `Prec` |
| `buildtag.go` | `BuildTag`, `LookupBuildTag`, `LicensesTest`, `HasFrameworks` |

- **`Pos` is a flat global integer**, not a `(file, offset)` pair. Each `File` claims a disjoint range of a shared `FileSet`, so a bare `Pos` carries its file implicitly — which is what lets `ast` store positions on punctuation-only fields without threading a file reference everywhere.
- **No `SEMICOLON`, no synthesized `NEWLINE`.** A.0.6 gives Vertex no statement terminator, so line structure arrives as `NLBefore` on the *following* token. The scanner can't know whether `{` opens a `Block` or a `CompositeLiteral` — that's the parser's `[+Lit]` question — so it records and never suppresses.
- **Contextual keywords are absent from `Kind` entirely.** `init`, `deinit`, `test` scan as `IDENT`; disambiguation is deferred past `Lookup` to the one production that names each.

### `diag`

The single currency for a rejection.

| File | Contents |
| --- | --- |
| `code.go` | `Code`, range scheme, message registry, `declaredCodes` |
| `diag.go` | `Diagnostic`, `Note`, `Fixit`, `New`/`At`/`AtToken` |
| `list.go` | `Reporter`, `List`, sort and dedup |
| `render.go` | `Renderer` — carets, notes, fix-its, colour |

The registry is normative, not an implementation detail: A.12.2's `Expected(error, "...")` compares rendered text as specification, so **changing a template is a language change**. Codes are numbered explicitly rather than by `iota` so inserting one can't shift another.

```
0xxx internal      3xxx declarations/names   7xxx pointers and memory
1xxx lexical       4xxx types                8xxx interop
2xxx syntactic     5xxx ownership            9xxx concurrency, devices
                   6xxx generics
```

`1xxx`/`2xxx` are dense; `3xxx`/`4xxx` fill in with the resolve pass; `5xxx`–`9xxx` await passes that don't exist. A.14 is their inventory.

### `scanner`

Bytes to tokens, recovering from everything: an illegal character is consumed and returned as `INVALID`, an unterminated string returns what it scanned. `Scan` never signals failure out of band, so one bad file still yields a full token stream.

| File | Contents |
| --- | --- |
| `scanner.go` | `Scanner`, `Init`, dispatch, punctuators, comments, `Tokenize` |
| `ident.go` | `IdentifierStart`/`Part`, Unicode ID properties |
| `number.go` | `scanNumber`, `scanDigits`, `scanTupleIndex` |
| `string.go` | `scanString`, `scanRawString`, `scanChar`, `scanEscape` |

Two decisions reach beyond this package:

- **A fractional part must be non-empty**, so `1.` isn't a literal — which makes `1..5` scan as `INT DOTDOT INT` with no lookahead surgery.
- **`s.prev == PERIOD` routes a digit to `scanTupleIndex`.** A.4.3's `t.0.0` falls out: the second dot ends the run and the rule fires again. The alternative was splitting a `FLOAT` apart in the parser.

Literals keep their raw spelling — `vfmt` needs the original text, and decoding is the analyzer's job.

### `ast`

The syntax tree. Records shape, never meaning.

| File | Contents |
| --- | --- |
| `ast.go` | `Node`/`Expr`/`Stmt`/`Decl`, comments, `Ident`, `Param`, `Marker` |
| `expr.go` | Expressions *and* types — the same interface |
| `stmt.go` | Statements, `CaseClause`, `SelectClause` |
| `decl.go` | Declarations, `declare` blocks, foreign members |
| `file.go` | `File`, `Package`, `NewPackage`, `BuildClause` |
| `walk.go` | `Walk`, `Inspect` |

**Types are `Expr`s.** A.3.6 distinguishes `Stack[int32]` from `a[i]` by whether the operand names a generic declaration — something the parser cannot know. Conversion to a type representation happens in the analyzer.

Several fields exist only for a consumer that doesn't exist yet — currently untested claims:

| Field | Why |
| --- | --- |
| `BasicLit.Value` | Raw spelling; a printer needs the original |
| `TupleExpr.TrailingComma` | `(x,)` and `(x)` differ |
| `TypeParam.Constraint` nil-for-earlier-names | Distributing would erase the written form |
| `TupleIndexExpr.Text` | Alongside the decoded `Index` |
| `File.Comments` | Flat list *and* per-node attachment |
| `ForeignFunc.Modifiers`, `.Body` | A.0.5 makes rejected forms part of the grammar |

### `parser`

Source to tree. Tracks exactly one of A.0.2's four context parameters.

| File | Contents |
| --- | --- |
| `parser.go` | Struct, token flow, diagnostics, recovery, `ParseFile`/`ParseDir` |
| `decl.go` | Top-level declarations and imports |
| `stmt.go` | Statements; assignment vs. expression-statement |
| `expr.go` | Expressions and types — precedence climbing, postfix chains |

**Only `Lit` is tracked.** `Await`, `Npu`, and `Own` name constructs that parse unconditionally and are rejected later; all three use dedicated keywords, so parsing them without context is unambiguous and still lets the diagnostic quote the construct. `Lit` is different — A.4.7's parenthesization rule is a decision made *while* parsing.

**Statement termination is a depth counter, not a token.** `(`, `[`, and a literal `{` push depth; a `Block`'s `{` resets it to zero, because statements inside a block do end at line terminators while entries inside a literal don't. `continues()` is `p.depth > 0 || !p.tok.NLBefore`.

**`ParseDir` is two-pass.** Every file is read with `PackageClauseOnly` into a throwaway `FileSet` first, since A.2.2 excludes a mistagged file from the build outright — the filter has to run before any file is fully parsed.

Where productions collide, the parser picks the reading with an actual grammar production and hands the other to the analyzer:

| Collision | Parser picks | Because |
| --- | --- | --- |
| `var w` as statement | `VarDecl` | `var Binding` is a production; a bare transfer statement isn't |
| `mut a, b` in a `for` | Both, combined | Neither alternative combines; A.14 rejects the pair |
| One identifier in a constraint | `ConstraintElem.Set` | A.7.2: resolved by what the name denotes |
| `Stack[int32]` vs `a[i]` | `IndexExpr` | Same |

### `types`

The type representation and the predicates over it. Depends only on `token`, except `info.go`.

| File | Contents |
| --- | --- |
| `type.go` | `Type` interface, `Mode`, `Marker` |
| `basic.go` | `Basic`, `Typ[]`, predeclared names, `Default` |
| `const.go` | `Value` over `math/big`, folding, **`Representable`** |
| `composite.go` | `Ownership`, `Array`, `Slice`, `Map`, `Tuple`, `Chan`, `Pointer`, `Signature`, `Tensor` |
| `named.go` | `Named`, `Struct`, `Enum`, `Abstract`, `TypeParam` |
| `constraint.go` | `Constraint`, `Term`, `Satisfies` |
| `predicates.go` | `Identical`, `IsComparable`, `AssignableTo`, `ConvertibleTo` |
| `sizes.go` | `Sizeof`, `Alignof`, `Offsetsof` |
| `object.go` | `Var`, `Const`, `Func`, `TypeName`, `Builtin`, `PkgName`, `Package` |
| `scope.go` | `Scope`, `Universe`, `Reserved` |
| `string.go` | Type and object rendering |
| `info.go` | The side tables — the only file importing `ast` |

**`Constraint` is not a `Type`.** A.7.2 makes a constraint legal only in a `[...]` position, and A.14 lists `var c: Ordered` among the rejected forms. Because `Constraint` doesn't implement `Type`, that rejection is structural: a checker that forgets to test has no `Type` to return.

**`mut`/`var` are a `Mode`; `unique`/`shared`/`weak` are `Type`s.** A.3.2 lets the latter three appear anywhere a `Type` may, while the former two are legal only in parameter or receiver position. The stacking rules fall out: `mut shared T` is a `ModeMut` `Var` over an `*Ownership`; `mut var T` is unrepresentable because a `Var` carries one `Mode`.

**`int` is 64-bit on every target**, including `js` and `wasm`. A target-varying width would make `sizeof(int)` and A.1.5.1's representability error both target-varying — a file could compile under one tag and fail under another for a reason invisible in the source. Only `WordSize` varies.

**`comparable` is pinned in `IsComparable`**, since A.7.4's definition is circular. Included: scalars, `string`, `typed_ptr`, fixed arrays, tuples, structs/classes, enums, recursively. Excluded: slices, maps, `chan`, `unique`/`shared`/`weak`, `abstract`. Consequence: a class with a slice field can't be a map key even though `===` still works — `===` asks about the allocation, `==` about the bytes.

**Constants hold unbounded precision** (`*big.Int`/`*big.Rat`) until a destination is known. That's what makes A.1.5.1's "never a silent truncation" checkable instead of already-lost.

**`TypeString` is normative** — it's what a `diag` template's `%s` receives, so changing it changes what an `Expected(...)` test matches.

### `analyzer`

Checks a parsed package against the static rules. Three phases, not one walk.

| File | Contents |
| --- | --- |
| `check.go` | `Checker`, `Config`, `Importer`, scopes, `declare`, `lookup` |
| `errors.go` | `report`, `errorAt`, `errorExpr`, `invalid` |
| `typexpr.go` | `ast.Expr` → `types.Type`, **and** → `*types.Constraint` |
| `decl.go` | Phases 1 and 2 |
| `resolve.go` | Phase 3 |

**Why three phases.** A.2 makes top-level declarations order-independent with no forward-declaration form, so:

1. **`collectObjects`** — every package-scope name exists, with a `nil` type. Imports bind into a *file*-scoped scope (A.2.3).
2. **`resolveDeclTypes`** — types filled in on demand via `objDecl`; the `resolving` stack turns a self-reference into `TypeCycle` instead of a stack overflow.
3. **`checkBodies`** — bodies walked, `Uses` recorded.

Collapsing any two reintroduces a forward-declaration requirement the language doesn't have.

**`typexpr.go` has two entry points, and that's the point.** `typ()` for a type position, `constraintExpr()` for a constraint position. A single identifier goes one way or the other by what the name denotes — a scope lookup, not a shape test.

Parser punts closed here:

| Ambiguity | Resolved by |
| --- | --- |
| `Stack[int32]` vs `a[i]` (A.3.6) | `denotesType(x.X)` |
| `TypeSet` vs `ConstraintName` (A.7.2) | `TypeName.IsConstraint()` |
| `~T` as bitwise-NOT vs underlying-type (A.7.3) | Position: `typ` rejects, `constraintExpr` accepts |
| `init`/`deinit` as method names (A.1.3) | `Func.IsInit()` — name plus receiver |
| Constraint distribution over `[A, B: Number]` (A.7.1) | `typeParams`, walking backwards |

**`npu` is the one context parameter tracked here.** A.3.5 makes `tensor` grammatical only under `[+Npu]`; the parser accepts it everywhere on purpose. `funcLit` clears it, per A.0.3: propagation stops at a function boundary.

**Nothing writes to the tree.** Every result lands in `types.Info`, which is what lets a printer round-trip a checked file and the checker run twice without residue. `invalid()` returns `Typ[Invalid]` rather than `nil` so one bad type yields one diagnostic, not a cascade — the same discipline `ast.BadExpr` follows on the parse side.

### `importer`

Resolves the import graph and drives the analyzer over it.

| File | Contents |
| --- | --- |
| `importer.go` | `Config`, `Package`, `Result`, `Load`, the `loader` |
| `resolver.go` | `Resolver`, `DirResolver`, `MapResolver`, `NotFoundError` |
| `graph.go` | `Sorted`, `Deps`, `Importers`, `Entry` |

**Loading is post-order.** A.2.3: the imported package's `PackageClause` is the qualifier its symbols are reached under; the import path is a locator, not a name. A package's qualifiers therefore depend on its imports having been read — which is what `PackageClauseOnly` exists for. Per package: `parse` → `loadImports` → `check`, in that order, because `collectImports` asks its `Importer` for a *complete* `*types.Package`, not a promise of one.

**A cycle is fatal, not a diagnostic.** Vertex has no forward-declaration form and no import-breaking construct, so there's no file to underline as the one that should have broken it.

**An erroring package is still marked complete** and still handed to dependents. Refusing would turn one bad file into a cascade of `UndeclaredName` across the graph.

**`pkgImporter` scopes to declared imports only.** A transitive dependency isn't reachable, because A.2.3 gives no spelling for one.

`Result.Order` is topological (what the analyzer ran in, what `lower` should follow); `Sorted()` is the stable-by-path alternative for reproducible output.

---

## How a rejection travels

`x.transfer()` — A.1.4 removes the spelling and keeps the name reserved so the diagnostic can carry a fix-it rather than degrading into "no such method."

1. **`scanner`** sees `IDENT("transfer")`. `token.Lookup` returns `IDENT` — `ReservedBuiltinName`s are absent from the keyword table because they're ordinary identifiers pre-bound in an implicit scope.
2. **`parser`** builds `CallExpr{Fun: SelectorExpr{Sel: Ident("transfer")}}`. No special case; ordinary call syntax.
3. **`types.Universe`** binds `transfer` to a `Builtin`, and `Reserved("transfer")` is true.
4. **`analyzer`** resolves the selector, finds no method, and — because the name is reserved — raises `TransferMethodRemoved` (`V5001`) with a `Fixit` rewriting to a `var` prefix.
5. **`diag.Renderer`** prints message, caret, and `help:` line.

Five packages, one rule, written down in exactly one place: the registry template. That's the intended shape for every entry in A.14.

---

## What's missing

### Before `lower/`

- **`analyzer/typecheck`** — the 4xxx pass: expression types, `Info.Types` populated, representability enforced at every typed position, constant folding, `let` inference, A.7.5 instantiation. **`lower` cannot start before this**; `Info.Types` is empty until it runs.
- **`analyzer/ownership`** — the 5xxx pass, and the one that changes generated code: A.9.2 liveness, A.9.3 exclusivity, and populating `Info.Transfers`. A.4.6 makes the marker "the entire difference between move and deep copy"; A.9.4 prices a copy at O(data) against a transfer at O(1).
- **`analyzer/generics`** — the 6xxx pass. `Info.Instances` is the monomorphization worklist; A.7.2 is explicit that a method requirement introduces no interface value and no vtable.

### Probably needed between analyzer and lower

A **`hir`/`core`** package. Straight from AST to `vir` means doing monomorphization, closure capture, `defer` ordering, deinit insertion, `async`/`await`, and `select` at once, against a tree still carrying source shape. A desugared middle IR — generics instantiated, drops explicit, still typed — is worth deciding on deliberately rather than discovering when `defer` inside a loop inside an `async` function arrives all at once.

### Deferred, not forgotten

- **`printer` / `vfmt`** — the `ast` fields above stay untested claims until something round-trips. `parse → print → parse` reaching a fixed point with structurally equal trees is the only end-to-end check available on the parser.
- **`cmd/vertex`** — owns the `FileSet`, picks a `BuildTag`, wires `diag.Renderer` to stderr. ~150 lines. `diag/render.go` has never rendered anything.
- **Tests** — notably the conformance test `diag/code.go` names by construction: every declared code has a registry entry, no two share a number, no two share a template. Go has no reflection over untyped constants, which is why `declaredCodes` is hand-maintained and why the test is the only thing that can check it.

---

## What `lower/` will need

| From | What |
| --- | --- |
| `importer.Result.Order` | Packages in dependency order |
| `Info.Types` | Every expression's type and constant value |
| `Info.Uses` / `Defs` | Which object each identifier means |
| `Info.Selections` | Field index, method target, tuple position |
| `Info.Instances` | The monomorphization worklist |
| `Info.Transfers` | Move vs. deep copy at each owning position |
| `Info.Defers` | Per-scope teardown, reverse registration order |
| `types.Sizes` | Layout, for offsets and allocation sizes |

And the guarantees from A.15 — each a runtime question the front end already turned into a compile-time proof:

| Absent from every binary | Because |
| --- | --- |
| Garbage collector | Static liveness; refcounts only where `shared` is written |
| Exception unwinder | The `(T, string)` tuple and ordinary control flow |
| Vtables, dynamic dispatch | No inheritance; generics monomorphized |
| Drop flags | Conditional transfer is a compile error |
| Runtime type information | Every cast resolved statically |
| Hidden allocation | The heap is reachable only through `unique`, `shared`, `[]T` |

The last row is why `lower` emits allocation only where the source spelled it; the fourth is why teardown never needs a runtime flag — a transferred binding simply has its teardown not emitted, which is exactly why conditional transfer is a compile error.