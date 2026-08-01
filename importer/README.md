# importer

```
github.com/vertex-language/vertex/importer
```

`importer` resolves a Vertex import graph into checked packages. It is the seam `ast.NewPackage`'s doc comment names: resolution belongs here rather than to the parser, which is why `parser.ParseDir` handles one directory and `importer` handles the graph.

It depends on `analyzer`, `ast`, `diag`, `parser`, `token`, and `types`, and calls all of them in that order for each package: parse, resolve imports, check. Nothing in the toolchain depends back on `importer` — like `analyzer`, it is a terminal stage, but for a whole graph rather than one directory.

The ordering constraint is A.2.3: an imported package's `PackageClause` supplies the qualifier its symbols are reached under, so a file's names cannot be resolved until every directory it imports has had its package clause read. That is what makes loading a post-order walk rather than a single pass, and what `parser.PackageClauseOnly` exists to support cheaply.

## Config and Load

### `Config`

```go
type Config struct {
    Fset     *token.FileSet
    Tag      token.BuildTag
    Resolver Resolver
    Reporter diag.Reporter
    Mode     parser.Mode
}
```

`Fset` is one `*token.FileSet` shared by the whole load, so a diagnostic can span two files in two different packages. `Tag` is threaded to both the parser and each package's `analyzer.Config`, since A.2.2's build tag decides what parses as well as what links. `Resolver` defaults to a `DirResolver` rooted at the current directory when nil. `Reporter` may be nil, in which case diagnostics are collected onto the returned `*Result` instead of streamed.

### `Load(conf *Config, paths ...string) (*Result, error)`

Resolves and checks the packages named by `paths`, plus everything they import transitively. Each root is loaded in turn via the internal `loader`, whose `packages` and `state` maps persist across roots, so a package imported by two different roots is loaded once.

Failure during a given root's load returns immediately with the partial `*Result` still populated — `Result.Packages` and `Result.Order` hold whatever finished before the error, which is what a caller wanting partial results (an editor) reads.

## `Package` and `Result`

### `Package`

```go
type Package struct {
    Path    string
    Dir     string
    Name    string
    Files   []*ast.File

    Types   *types.Package
    Info    *types.Info

    Imports map[string]*Package
    Errors  []*diag.Diagnostic
}
```

One loaded package: its syntax, its checked types, and its resolved dependencies. `Name` comes from the `PackageClause` (A.2.1), not from `Path` — a directory's import path and the name its own files declare are independent per A.2.3. `Imports` maps each import path this package declares to the loaded `*Package` it resolved to; per A.2.3 there is no aliasing, dot-import, or blank-import form, so the map is the whole story.

`Errors` holds this package's own diagnostics regardless of whether the load ultimately succeeds — a package with errors is still returned, since A.0.5 makes a rejected form part of the grammar rather than something that aborts a load. `String()` returns `Path`.

### `Result`

```go
type Result struct {
    Roots       []*Package
    Packages    map[string]*Package
    Order       []*Package
    Diagnostics *diag.List
}
```

`Roots` are the packages named in the `Load` call, in the order given. `Packages` holds every loaded package by import path, roots included. `Order` is topological — a package appears after everything it imports — which is both the order the analyzer ran in and the order `lower` should follow, since A.2.3 makes a package's qualifiers depend on its imports having already been read. `Diagnostics` is empty when `Config.Reporter` was non-nil, since diagnostics streamed there instead; otherwise it is sorted and deduplicated once the load finishes.

### `(*Result) HasErrors() bool`

True if any loaded package holds an error-severity diagnostic, checked across both `Order` and `Diagnostics` — the second check covers a `Reporter`-less load whose per-package `Errors` were populated but whose severities are only otherwise visible through the aggregated list.

## The loader

`loader` is unexported state for one `Load` call: `packages` and `state` (by path), the post-order accumulator `order`, the current root-to-here `stack` (for rendering a cycle), and `list`, the `*diag.List` collected into when `Config.Reporter` is nil.

### `load(path string, from *Package) (*Package, error)`

The recursive step, keyed by `loadState` (`unstarted` / `loading` / `loaded`):

- **`loaded`** — returns the cached `*Package` immediately.
- **`loading`** — an import cycle. This is fatal for the whole load rather than a diagnostic on some file: Vertex has no forward-declaration form and no way to break a cycle, so there is nothing to point at as the fix. Returned as `*CycleError`.
- **`unstarted`** — resolves a directory via `Resolver`, marks the path `loading`, and proceeds through `parse` → `loadImports` → `check` in that fixed order, restoring `state[path] = loaded` via `defer` regardless of outcome.

Imports are recursed into *before* the package itself is checked (`loadImports` runs before `check`): `collectImports` in `analyzer` asks its `Importer` for a complete `*types.Package`, not a promise of one, so every dependency must already be a finished `*Package` by the time this one's `analyzer.Checker` runs.

### `parse(pkg *Package) error`

Delegates to `parser.ParseDir`, which already implements A.2.2's two-pass shape — every file's package and build clauses are read into a throwaway `FileSet` first, so the build-tag filter runs before any file is fully parsed, keeping a mistagged file's exclusion whole rather than partial. Diagnostics are forwarded through a `diag.ReporterFunc` that both appends to `pkg.Errors` and calls `l.report`.

### `loadImports(pkg *Package) error`

Walks every `ImportDecl` in every file of `pkg`, recursing via `load` and populating `pkg.Imports`. A `seen` set by path means a path imported by two files of the same package, or twice in one file, recurses once. A cycle surfaced from a nested `load` call has its `Importer` field set here to `pkg.Path`, identifying which package's import list actually closed the cycle.

### `check(pkg *Package)`

Allocates `pkg.Info`, builds an `analyzer.Config` around a `pkgImporter` scoped to this one package, and runs `analyzer.NewChecker(...).Files(pkg.Files)`. The resulting `*types.Package` is stored on `pkg.Types` and, if non-nil, marked complete — per `types`'s own `Package.Complete` doc comment, a half-checked scope must never be handed to a dependent, since it would answer a lookup with a nil type. A package is handed onward as complete once its *own* check finishes, errors included, so that one bad file does not cascade into every package downstream of it.

### `pkgImporter`

Satisfies `analyzer.Importer` over exactly the dependencies `loadImports` already resolved for one package — deliberately not the whole graph. `Import(path)` fails for any path the package did not itself declare an import of, even if some other loaded package imports it: A.2.3 gives no spelling for a name reached through a transitive dependency, so `pkgImporter` refuses to let one leak in that way.

## Resolvers (`resolver.go`)

### `Resolver`

```go
type Resolver interface {
    Resolve(path, from string) (dir string, err error)
}
```

Maps an import path to the directory holding that package's `.vs` files. It is an interface because A.2.3 says only that "the import path is a locator, not a name" — a build system, a module cache, and a test fixture all answer the locate question differently, and `importer` doesn't pick one. `from` is the importing package's directory, or `""` for a root; a resolver that supports relative imports uses it, one that doesn't ignores it. `ResolverFunc` adapts a plain function to the interface.

### `DirResolver`

Resolves a path against one or more root directories, in order, made absolute by `NewDirResolver`. A `./` or `../`-prefixed path resolves relative to the importing package's directory instead of against the roots, and fails if there is no importing package (`from == ""`) to resolve against. This is deliberately the simplest policy that works: no module resolution, no version selection, no vendor directory, since the annex describes none of those and inventing one here would fix a choice the toolchain hasn't made.

`hasSources` gates a candidate directory by extension alone (`.vs`) — whether a given file is actually *in* the build is A.2.2's build-tag question, which `parser.ParseDir`'s own first pass answers, so a directory whose files are all tagged for another target resolves successfully but yields no package there.

### `MapResolver`

`map[string]string`, resolving from a fixed table. Exists for a driver that has already computed a layout, and for tests that want no filesystem dependency at all.

### `NotFoundError`

Reports an unresolvable path, carrying every directory searched and, when a searched directory exists but holds no `.vs` file, a `Reason` distinguishing that from a directory that doesn't exist at all — the two are different mistakes and the message says which.

## Errors (`importer.go`)

### `CycleError`

```go
type CycleError struct {
    Path     string
    Importer string
    Stack    []string
}
```

Reports an import cycle as fatal rather than as a `diag.Diagnostic` on some file, since Vertex has no forward-declaration form and therefore no single file whose fix would break the cycle. `Stack` is the root-to-here path at the point the cycle was detected; `Importer` (set by `loadImports`) identifies which package's own import list closed it.

## The query surface (`graph.go`)

Everything in `graph.go` answers questions about an already-completed `Result` — nothing here parses or checks.

- **`Sorted() []*Package`** — every loaded package by path, in lexical order. `Result.Order` is topological and therefore not stable against an unrelated edit; `Sorted` is what a caller wanting reproducible output over the whole set uses instead.
- **`Lookup(path string) *Package`** — a loaded package by import path, or nil.
- **`Deps(p *Package) []*Package`** — every package reachable from `p`, excluding `p`, in topological order: the set a link step needs and the set `lower` walks. Computed by a `seen`-guarded post-order walk over each package's `Imports`, sorted by path at each level for determinism.
- **`Importers(path string) []*Package`** — the reverse edge: every loaded package that directly imports `path`. Nothing in a build needs this; tooling (e.g. "what breaks if I change this package") does.
- **`Entry() *Package`** — the package holding the program's entry point, or nil, searching the roots rather than assuming a package literally named main since A.6.1 fixes what shape a main function has but not which package it lives in. A root whose `Types` is still nil (a failed check) is skipped rather than causing a panic.