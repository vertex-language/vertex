# importer

```go
import "github.com/vertex-language/vertex/importer"
```

`importer` resolves a Vertex import graph into checked packages. It is the seam `ast.NewPackage` names — that constructor is "a validated container and nothing more: no I/O, no import resolution, no scopes. Resolution belongs to the loader, which is why this takes no importer." `parser.ParseDir` handles one directory; `importer` handles the graph.

It depends on `analyzer`, `ast`, `diag`, `parser`, `token`, and `types`, and calls them in that order per package: parse, resolve imports, check. Nothing depends back on it.

## Design philosophy

The ordering constraint is §1.3: the qualifier under which an imported package's symbols are reached comes from that package's own package clause, and the path is a locator rather than a name. A file's names cannot resolve until every directory it imports has had its package line read — which is what makes loading a **post-order walk** rather than a single pass, and what `parser.PackageClauseOnly` exists to support cheaply.

**Two failure kinds, deliberately separate.** A `return`ed `error` is fatal to the load: an unresolvable path, or an import cycle. Neither has a file to point at as the fix, and Vertex has no forward-declaration form to break a cycle with. Everything else is a `*diag.Diagnostic` on some package's `Errors`, and the load continues — most of what this toolchain rejects is a form that parses and is diagnosed, not one that aborts anything.

**Partial results always.** The `*Result` is non-nil on every path including a fatal one, and holds whatever finished. `Load`'s `defer` populates `Order` and `Diagnostics` on exit rather than only on success, so an editor gets the graph it managed to build.

**Citation convention.** A bare `§` cites semantics.md; CamelCase names are grammar.md productions, matching `types` and `analyzer`. Where neither document fixes something — how a path locates a directory, whether a relative import is even spelled — the comment says so.

## Package layout

| File | Contents |
|---|---|
| `importer.go` | Package doc, `Config`, `Package`, `Result`, `Load` |
| `loader.go` | `loader` and the walk (`load`/`parse`/`loadImports`/`check`), `pkgImporter`, `CycleError` |
| `resolver.go` | `Resolver`, `ResolverFunc`, `DirResolver`, `MapResolver`, `NotFoundError` |
| `graph.go` | The query surface over a finished load: `Sorted`, `Lookup`, `Deps`, `Importers`, `Entry` |

## Config and Load

```go
type Config struct {
	Fset     *token.FileSet
	Tag      token.BuildTag
	Resolver Resolver      // nil → DirResolver rooted at "."
	Reporter diag.Reporter // nil → collected onto Result.Diagnostics
	Mode     parser.Mode
}

func Load(conf *Config, paths ...string) (*Result, error)
```

`Fset` is one `*token.FileSet` for the whole load, so a diagnostic can span two files in two packages. `Tag` is threaded to both the parser and each `analyzer.Config`, since the build tag decides what parses as well as what links.

The internal `loader`'s `packages` and `state` maps persist across roots, so a package imported by two roots is loaded once.

## `Package` and `Result`

`Package.Name` comes from the package clause, never from `Path` — §1.3 makes the two independent, and the name is the qualifier importers use. `Imports` maps each declared path to its loaded package; with no aliasing, dot-import, or blank-import form, the mapping is the whole story. `Errors` holds this package's diagnostics from both the parse and the check.

`Result.Order` is topological and is what a lowering pass follows. `Sorted()` is what a caller wanting output stable against unrelated edits uses instead. `Diagnostics` is nil when a `Reporter` was supplied; `HasErrors` reads the per-package lists, which are populated either way.

## The walk

`load` is keyed by `loadState`:

| State | Behaviour |
|---|---|
| `loaded` | return the cached `*Package` |
| `loading` | an import cycle → `*CycleError`, fatal |
| `failed` | a plain error; the original failure was already returned to whoever reached it first |
| `unstarted` | resolve → `parse` → `loadImports` → `check`, in that fixed order |

Imports are recursed into *before* the package is checked, because `analyzer`'s `collectImports` asks its `Importer` for a complete `*types.Package`, not a promise of one.

`parse` delegates to `parser.ParseDir`, which already runs the two passes — every file's package and build clauses read against a throwaway `FileSet` first, so the tag filter runs before any file is fully parsed and a mistagged file's exclusion stays whole. A nil `*ast.Package` is fatal (no source in this build); a non-nil one alongside an error is a syntax error already reported, and the check proceeds over the `Bad*` recovery nodes.

`check` allocates `pkg.Info`, runs the checker, and then calls `MarkComplete` **unconditionally** — deliberately overriding `analyzer.Files`, which marks complete only on a clean check. That default is right for a caller holding one package; a loader holding a graph wants the opposite, since refusing to hand a partial scope to a dependent turns one bad file into a cascade of undeclared-name errors across everything downstream.

`pkgImporter` scopes each package's `Importer` to exactly the dependencies it declared. A path some *other* package imports is refused: §1.3 gives no spelling for a name reached through a transitive dependency.

## Resolvers

```go
type Resolver interface {
	Resolve(path, from string) (dir string, err error)
}
```

An interface because the path is a locator and nothing says what it locates against. `from` is the importing package's directory, or `""` for a root.

`DirResolver` searches roots in order, first hit with a `.vs` file wins. A `./` or `../` path resolves against `from` instead and fails without one; an absolute path is rejected outright, since no root could apply to it. No module resolution, no versions, no vendoring — the language describes none, and inventing one here would fix a choice the toolchain hasn't made. `hasSources` gates by extension only; whether a file is *in* the build is the tag's question, answered in `ParseDir`.

`NotFoundError` distinguishes a directory that exists but holds no `.vs` file from one that doesn't exist, since those are different mistakes.

## Entry points

`Entry()` searches the whole graph for a package named `main` whose `main` satisfies `types.Func.IsEntry` (no parameters, void, no marker, no receiver). §1.4 requires exactly one; if the load holds more than one, `Entry` returns nil rather than picking, and the caller must say which.