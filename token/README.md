# token

`github.com/vertex-language/vertex/token`

`token` defines the lexical vocabulary and source-position machinery shared by every stage of the Vertex toolchain. It has no dependency on any other package in the toolchain; `scanner`, `ast`, `diag`, the analyzer, and the printer all depend on it instead.

## Positions

### `Pos`

```go
type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool { return p != NoPos }
```

An opaque, comparable offset into a `FileSet`'s address space. `NoPos` is the zero value, reserved as the sentinel for "this optional position is absent" — no valid `Pos` is ever `0`, since `FileSet.base` starts at `1`.

### `Position`

```go
type Position struct {
	Filename string
	Offset   int // byte offset, 0-based
	Line     int // 1-based
	Column   int // 1-based, in bytes
}
```

The resolved, human-facing counterpart to `Pos` — what a `Pos` means once you know which file it's in. `IsValid` tests `Line > 0`. `String` renders `file:line:col`, or bare `line:col` if `Filename` is empty, or `"-"` if invalid.

### `File`

```go
type File struct { /* unexported */ }
```

One source file's slice of a `FileSet`'s address space: a `name`, a `base` offset, a `size`, and a `lines` table.

- **`Name`**, **`Base`**, **`Size`** — basic accessors.
- **`AddLine(offset int)`** — records that a line begins at `offset`. The scanner calls this each time it consumes a `LineTerminator` (A.1). Offsets must be monotonically increasing; a non-increasing or out-of-range offset is silently dropped rather than corrupting the table.
- **`LineCount() int`** — number of recorded line starts.
- **`LineStart(line int) Pos`** — the `Pos` of a line's first character, or `NoPos` if `line` is out of range. Used to emit debug line tables.
- **`Pos(offset int) Pos`** / **`Offset(p Pos) int`** — convert between a file-local byte offset and a global `Pos`, each panicking if the argument falls outside the file.
- **`Position(p Pos) Position`** — resolves `p` to line/column via `sort.SearchInts` over `lines`, under the file's own mutex.

`File` guards `lines` with a `sync.Mutex` since `AddLine` (write, from the scanner) and `Position`/`LineCount` (read, from anything rendering a diagnostic) can run concurrently once a driver starts streaming diagnostics through a `Reporter`.

### `FileSet`

```go
type FileSet struct { /* unexported */ }
```

Owns the mapping from `Pos` back to `File`. One per compilation — every file a package loads shares the same set, so a diagnostic can span two files without extra plumbing.

- **`NewFileSet() *FileSet`** — the only constructor; `ast` and every other consumer take a `*FileSet` as an argument rather than building one.
- **`(*FileSet) AddFile(name string, size int) *File`** — registers a new file at the current `base` and advances `base` past it (with a one-`Pos` gap, so `End()` positions one-past-the-last-character never collide with the next file's start).
- **`(*FileSet) File(p Pos) *File`** — resolves `p` to its containing `*File` via binary search over `files`, or `nil` if `p` doesn't belong to any registered file.
- **`(*FileSet) Position(p Pos) Position`** — convenience for `File(p).Position(p)`, returning a zero `Position` if `p` resolves to no file.

`FileSet` guards `files` with a `sync.RWMutex`, separately from each `File`'s own lock — adding a file and resolving a position in a different, already-registered file don't contend.

## Tokens

### `Token`

```go
type Token struct {
	Kind     Kind
	Pos      Pos
	Lit      string
	NLBefore bool
}
```

One lexeme. `Lit` is the raw source spelling for `IDENT`, `INT`, `FLOAT`, `CHAR`, `STRING`, and `COMMENT` — unescaped, since decoding is the analyzer's job, not the scanner's — and empty otherwise. `NLBefore` reports whether a `LineTerminator` (or a multi-line comment, per A.1.1) separates this token from the previous one; it's also what a `[no LineTerminator here]` grammar restriction reads.

There is deliberately no `SEMICOLON` kind and no synthesized `NEWLINE` kind. Vertex has no statement terminator (A.0.6), so line structure reaches the parser as a flag on the *following* token rather than as a token of its own — see Design below.

- **`(Token) Is(k Kind) bool`** — `Kind == k`.
- **`(Token) End() Pos`** — `Pos + len(Lit)`; the conventional "one past the last character" used throughout `ast` for node extents.
- **`(Token) IsBlank() bool`** — reports the `BlankIdentifier` `_` (`Kind == IDENT && Lit == "_"`).

### Contextual keywords

```go
const (
	CtxBuild     = "build"
	CtxDeinit    = "deinit"
	CtxError     = "error"
	CtxFramework = "framework"
	CtxInit      = "init"
	CtxModule    = "module"
	CtxTest      = "test"
)
```

Spellings for A.1.3's `ContextualKeyword`s. Each scans as an ordinary `IDENT` — `token` mints no separate `Kind` for any of them — and is disambiguated by the parser at the one production that names it (e.g. `CtxBuild` only means something as the second line-initial token of a file; `CtxTest` only as a function marker).

- **`(Token) IsCtx(name string) bool`** — `Kind == IDENT && Lit == name`. Call sites read `tok.IsCtx(token.CtxFramework)` rather than comparing `Lit` directly.

## Token kinds

### `Kind`

```go
type Kind int
```

Represents punctuation, operators, reserved keywords, and literal categories. The constant block is organized into contiguous ranges bounded by unexported sentinels (`literalBeg`/`literalEnd`, `reservedLitBeg`/`reservedLitEnd`, `operatorBeg`/`operatorEnd`, `keywordBeg`/`keywordEnd`), so membership tests are a single pair of comparisons rather than a switch:

- **`(Kind) IsLiteral() bool`** — `INT`, `FLOAT`, `CHAR`, `STRING`, `IDENT` (the scanned literal kinds) or `TRUE`/`FALSE`/`NIL` (A.1.3's `ReservedLiteralKeyword`s — literals syntactically, but reserved lexically and carrying no `Lit` text, which is why they sit in their own range rather than the scanned one).
- **`(Kind) IsOperator() bool`**, **`(Kind) IsKeyword() bool`** — same pattern for the operator and keyword ranges.

### `(Kind) String() string`

Renders a `Kind` back to its source spelling via a parallel `[...]string` array indexed by `Kind`, e.g. `ADD.String() == "+"`. Used both for diagnostics and for position arithmetic where a node's `End()` is derived from a keyword's length rather than a stored token.

### `Lookup`

```go
func Lookup(ident string) Kind
```

Maps an identifier spelling to its keyword `Kind`, or `IDENT` if it isn't one. The backing `keywords` map is built once in `init()` from the keyword range of `kinds`, plus `true`/`false`/`nil` added explicitly since `ReservedLiteralKeyword`s are reserved lexically and so belong in the same lookup table despite living in a separate `Kind` range.

`ContextualKeyword`s are deliberately absent from this map — A.1.3 makes them identifiers everywhere except the single production that names each, so baking one into `Lookup` would make it a keyword unconditionally. `PredeclaredTypeName` and `ReservedBuiltinName` (A.1.4) are absent for a different reason: they're ordinary identifiers pre-bound in an implicit scope, and the scanner must not know them at all.

### Precedence

```go
const (
	LowestPrec  = 0
	UnaryPrec   = 9
	HighestPrec = 10
)

func (k Kind) Prec() int
```

`Prec` returns a binary operator's precedence level, or `LowestPrec` for anything that isn't one — the single source of truth the parser (and any consumer re-deriving structure, such as a formatter) uses instead of encoding precedence in the tree.

`DOTDOT` is listed at level 4 but is **non-associative** (A.4.5): `Prec` reports a level for it like any other operator, but the parser's precedence-climbing loop must special-case it to reject `a..b..c` rather than fold it left- or right-associatively. The table entry alone can't express non-associativity.

### `(Kind) IsCompoundAssign() bool`

Reports whether `k` is one of the ten `CompoundAssignOperator`s (`+=` through `>>=`, A.5.2). `AssignStmt.Op` stores whichever `Kind` matched — `ASSIGN` or a compound kind — and this predicate is how the parser and analyzer tell the two cases apart without a second switch.

## Build tags

### `BuildTag`

```go
type BuildTag int

const (
	TagNone BuildTag = iota
	TagLinux
	TagWindows
	TagDarwin
	TagJS
	TagWasm
	TagTest
)
```

The A.2.2 target selector. It lives in `token` rather than `parser` because the parser, loader, and driver all need it, and because it changes what's grammatical — `LicensesTest` gates whether the `test` marker and `Expected(...)` result types are legal at all.

- **`(BuildTag) String() string`** — the tag's spelling (`"linux"`, `"darwin"`, …), or `""` for `TagNone`. Used in `%q` form when constructing build-mismatch errors.
- **`LookupBuildTag(s string) (BuildTag, bool)`** — the inverse of `String`. The `bool` is load-bearing: A.2.2 makes an unrecognized tag a compile error, never a silently-excluded file, so callers must distinguish "unknown spelling" from `TagNone`'s "no clause at all" — a plain zero-value return can't do that.
- **`LicensesTest(t BuildTag) bool`** — `t == TagTest`.
- **`HasFrameworks(t BuildTag) bool`** — `t == TagDarwin`; whether `declare framework` (A.8.1) is legal on this target. The predicate lives here, next to the tag values it switches on, so the answer has one home instead of being re-derived at each call site.

## Design

**`Pos` is a flat, global integer space, not a `(file, offset)` pair.** Every `File` claims a disjoint range of a shared `FileSet`, so a bare `Pos` carries its file implicitly and `FileSet.File` recovers it by binary search. This is what lets `ast` store `Pos` liberally — even on punctuation-only fields like `Rparen token.Pos` — without also threading a file reference everywhere: `Pos` alone is enough to resolve a full `Position` later.

**Kind ranges are sentinel-delimited, not tag-checked.** `IsLiteral`/`IsOperator`/`IsKeyword` all compile down to `lo < k && k < hi` because the constant block groups related kinds contiguously by construction. This only works because `Kind` values are assigned by `iota` in declaration order — unlike `diag.Code`, where explicit numbering matters because renumbering a diagnostic code breaks the spec's normative text, renumbering a `Kind` is safe as long as the group boundaries move with it.

**Contextual keywords are absent from `Kind` entirely**, not merely unregistered in `Lookup`. Names like `init`, `deinit`, and `test` scan as plain `IDENT`; disambiguation is deferred past the scanner and even past `Lookup`, to the individual parser productions and analyzer passes that care — consistent with the toolchain's broader habit of having the tree "record shape, never meaning."

## File organization

- **`pos.go`** — `Pos`, `Position`, `File`, `FileSet`.
- **`token.go`** — `Token`, the contextual-keyword constants, `IsCtx`, `IsBlank`.
- **`kind.go`** — `Kind`, the `kinds` string table, `Lookup`, the classification predicates, and precedence.
- **`buildtag.go`** — `BuildTag` and its lookup/predicate functions.