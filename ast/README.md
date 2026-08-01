# token

`token` defines the lexical vocabulary and source-position machinery shared by every stage of the Vertex toolchain — scanner, parser (`ast`), analyzer, and printer. It has no dependency on `ast`; `ast` depends on it.

This document describes the package's public surface as inferred from its usage across `github.com/vertex-language/vertex/ast`.

## Positions

### `Pos`

```go
type Pos int
```

`Pos` is an opaque, comparable, offset-like handle into a `FileSet`. Every node in `ast` records its extent as a `Pos` pair (`Pos()`/`End()`), computed either directly from a stored field or arithmetically from a stored `Pos` plus the byte length of some spelled text, e.g.:

```go
func (x *Ident) End() token.Pos { return x.NamePos + token.Pos(len(x.Name)) }
```

This means `Pos` supports integer-like addition against `token.Pos(int)`, and `End()` is conventionally "one past the last character."

### `NoPos`

```go
const NoPos Pos = 0
```

The zero value of `Pos`, used as a sentinel for "this optional position is absent" — e.g. `Colon` on a `Param` is `NoPos` when `Name` is nil, and `Ellipsis` is `NoPos` when the parameter isn't variadic.

### `(Pos) IsValid() bool`

Reports whether a `Pos` is anything other than `NoPos`. Used throughout `ast` to test optional/error-recovery fields before consulting them, e.g.:

```go
if b.Ellipsis.IsValid() {
    return b.Ellipsis
}
```

## Files and file sets

### `FileSet`

```go
type FileSet struct { /* ... */ }
```

Owns the mapping from `Pos` values back to source files and filenames. Constructed and threaded through the toolchain by callers (e.g. `ast.NewPackage(fset, ...)`); `ast` never constructs one itself.

### `(*FileSet) File(p Pos) *File`

Resolves a `Pos` to the `*File` (see below) that contains it, or `nil` if `p` doesn't belong to any file registered in the set. `ast.File.Filename` uses this to turn its `Package` position back into a name:

```go
func (f *File) Filename(fset *token.FileSet) string {
    if tf := fset.File(f.Package); tf != nil {
        return tf.Name()
    }
    return ""
}
```

### `File` (token package)

```go
type File struct { /* ... */ }
```

Not to be confused with `ast.File` (the parsed syntax tree) — this is the `FileSet`'s bookkeeping record for one source file.

#### `(*File) Name() string`

Returns the filename as registered in the `FileSet`.

## Build tags

### `BuildTag`

```go
type BuildTag string // or similar underlying type
```

Identifies the target a file's `build` clause selects for (A.2.2). Formats via `%q`, as seen in error construction:

```go
fmt.Errorf("package %s: %s carries build tag %q, not %q", path, ..., t, target)
```

### `TagNone`

```go
const TagNone BuildTag = ""
```

The tag value meaning "no build clause" — distinct from an unrecognized tag name, which `ast.BuildClause` represents the same way but which callers must reject as a compile error rather than silently treat as absent (per `ast.BuildClause`'s doc comment).

## Token kinds

### `Kind`

```go
type Kind int // or similar small integer type
```

Represents both punctuation/operator tokens and reserved keywords produced by the scanner. `ast` stores `Kind` directly on nodes that need to distinguish between syntactically-identical alternatives (e.g. `RecordDecl.Kw` is `STRUCT` or `CLASS`; `VarDecl.Kw` is `LET` or `VAR`; `AssignStmt.Op` is `ASSIGN` or a compound-assignment kind).

Constants referenced across `ast` include (non-exhaustive):

- **Literals:** `INT`, `FLOAT`, `CHAR`, `STRING`, `TRUE`, `FALSE`, `NIL`
- **Keyword namespaces / launch prefixes:** `ASYNC`, `GPU`, `NPU`, `CHAN`, `THREAD`
- **Ownership qualifiers:** `MUT`, `VAR`, `UNIQUE`, `SHARED`, `WEAK`
- **Declarations:** `STRUCT`, `CLASS`, `LET`
- **Statements:** `BREAK`, `CONTINUE`, `FALLTHROUGH`, `ASSIGN`, `ILLEGAL`
- **Operators:** `DOTDOT` (the non-associative range operator, A.4.5), `TILDE`
- **Sentinel:** `IDENT` (used as `Marker.Kind` for the `test` `ContextualKeyword`, since it scans indistinguishably from an identifier)

#### `(Kind) String() string`

Renders a `Kind` back to its source spelling. Used both for diagnostics and for position arithmetic where a node's `End()` is derived from a keyword's length rather than a stored token:

```go
func (x *NamespaceExpr) End() token.Pos {
    return x.KwPos + token.Pos(len(x.Kw.String()))
}
```

#### `(Kind) Prec() int`

Returns binary operator precedence. `ast.BinaryExpr` deliberately stores no precedence itself — per its doc comment, *"Precedence comes from `token.Kind.Prec()` and A.13; the cascade nonterminals are grammar-writing devices, not shapes"* — so the parser (and any consumer re-deriving structure, such as a formatter) calls this rather than encoding precedence in the tree.

## Design notes carried over from `ast`

A few properties of `token` are implied by how carefully `ast` avoids duplicating information it can express instead:

- **`Kind` is the single source of truth for spelling and precedence.** `ast` nodes hold a `Kind` plus a position, not redundant spelling or precedence fields, wherever `Kind` alone determines them (`BinaryExpr`, `NamespaceExpr`, `LaunchExpr`).
- **`Pos` is cheap enough to store liberally.** Every `ast` node — including punctuation-only ones like `Rparen token.Pos` — keeps its own positions rather than recomputing them, implying `Pos` is a small value type (consistent with an `int` offset).
- **Contextual keywords are not a `Kind` concern.** Names like `init`, `deinit`, and `test` scan as plain `IDENT` (per `ast.FuncDecl` and `ast.Marker`'s doc comments) — `token` does not mint separate `Kind` values for them. Disambiguation is deferred to the analyzer, which is consistent with `ast`'s stated philosophy that the tree "records shape, never meaning."