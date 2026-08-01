# diag

`github.com/vertex-language/vertex/diag`

`diag` is the shared diagnostic currency for the Vertex toolchain — scanner, parser, analyzer, and `vir` all produce `*diag.Diagnostic` values through this package rather than formatting their own error strings. It depends on `token` (`github.com/vertex-language/vertex/token`) for `Pos` and `FileSet`; nothing in the toolchain depends on `diag` for anything but diagnostics.

The package splits into four concerns: stable rule identifiers (`Code`), the diagnostic value itself, collection (`List`), and human-facing rendering (`Renderer`). The first two are what a phase produces; the last two are what a driver does with the result.

## Codes

### `Code`

```go
type Code int
```

`Code` is the stable handle for a diagnostic rule. Renumbering one is a breaking change to the language's test corpus, because A.12.2's `Expected(error, "...")` compares rendered message text as specification. Values are assigned explicitly rather than via `iota` so that inserting a new code can never shift another's number.

Codes fall into ranges keyed to A.14's categories:

```
0xxx  internal
1xxx  lexical                       (A.1)
2xxx  syntactic                     (A.2–A.5)
3xxx  declarations and names        (A.6)
4xxx  types                         (A.3)
5xxx  ownership and exclusivity     (A.9)
6xxx  generics                      (A.7)
7xxx  pointers and memory           (A.4.8)
8xxx  interop                       (A.8)
9xxx  concurrency, devices, testing (A.10–A.12)
```

Gaps within a range are deliberate and reserved for analyzer passes that don't exist yet; only 1xxx (scanner) and 2xxx (parser) are densely populated today, plus two early entries in 5xxx and 7xxx that exist because their spec sections already commit to a named fix-it (`TransferMethodRemoved`, `AddrOnNonPointer`/`AddrOnTemporary`).

### `Internal`

```go
const Internal Code = 1
```

Marks a broken compiler invariant — never produced by valid or invalid source. If a caller sees one, it's a bug in the compiler, not in the user's program. `New` also falls back to `Internal` automatically if asked to build a diagnostic from an unregistered code.

### Registry

Each `Code` maps to a message template and default `Severity` in an internal registry. The template is the normative text: changing one changes what `Expected(error, "...")` tests match, so it's treated as a language change rather than an implementation detail. Templates avoid deictic phrasing ("this literal", "here") in favor of naming the construct, since `Diagnostic.Text()` strips position information and the message has to read the same with or without a span shown.

- **`(Code) String() string`** — renders the stable identifier, e.g. `"V1003"`. This is what a test that cares which rule fired should pin, rather than the message text.
- **`(Code) Severity() Severity`** — the code's registered severity, defaulting to `Error` for an unregistered code so a registry gap can't silently downgrade a rejection into a warning.
- **`(Code) Template() string`** — the raw format string, intended for the test runner and tooling that enumerates the normative surface, not for call sites.
- **`(Code) Registered() bool`** — whether the code has a registry entry at all.
- **`Codes() []Code`** — every declared code, in declaration order, from a hand-maintained list. It exists so a conformance test can assert that every declared code has a registry entry, that no two codes share a numeric value, and that no two share a template — properties the compiler can't check itself since Go has no reflection over untyped constants.

## Diagnostics

### `Severity`

```go
type Severity uint8

const (
	Error Severity = iota
	Warning
)
```

Stringifies to `"error"` or `"warning"`.

### `Diagnostic`

```go
type Diagnostic struct {
	Code Code
	Sev  Severity

	Pos token.Pos
	End token.Pos

	Message string

	Notes  []Note
	Fixits []Fixit
}
```

The single currency for a rejection, produced identically by every phase. `Pos` is where the diagnostic points; `End` bounds the underlined span, and if it's `NoPos` or doesn't exceed `Pos`, only one column is underlined. `Message` is the already-rendered template — the normative text `Expected(error, "...")` compares against.

`Notes` carries secondary spans for diagnostics that need to point at more than one place — a use-after-transfer (A.9.2) needs both where the binding died and where it was subsequently read. `Fixits` carries machine-applicable edits; the field exists from day one because the spec already commits to two fix-its by name (A.1.4's `.transfer()` removal and A.4.8's `addr(x)` → `&x`), and threading a new field through every constructor later would be worse than carrying an empty slice now.

### `Note` and `Fixit`

```go
type Note struct {
	Pos     token.Pos
	End     token.Pos
	Message string
}

type Fixit struct {
	Message string
	Pos     token.Pos
	End     token.Pos
	Text    string
}
```

A `Fixit` is a single replacement over `[Pos, End)`: an empty `Text` is a deletion, an empty span is an insertion.

### Constructors

- **`New(code Code, pos, end token.Pos, args ...any) *Diagnostic`** — builds a diagnostic by formatting the code's registered template with `args`. Call sites never format their own text; that's what keeps one rule's message identical across every site that raises it. If `code` isn't registered, `New` returns an `Internal` diagnostic instead of panicking.
- **`At(code Code, pos token.Pos, args ...any) *Diagnostic`** — `New` for a zero-width position (`end` is `token.NoPos`).
- **`AtToken(code Code, t token.Token, args ...any) *Diagnostic`** — `New` over a token's full extent, using `t.Pos` and `t.End()`.

### Mutators

- **`(*Diagnostic) WithNote(pos, end token.Pos, format string, args ...any) *Diagnostic`** — appends a `Note`, formatted like `fmt.Sprintf`. Returns the receiver so it chains.
- **`(*Diagnostic) WithFixit(pos, end token.Pos, text, format string, args ...any) *Diagnostic`** — appends a `Fixit`. Typical use:

  ```go
  d.WithFixit(call.Pos(), call.End(), "var "+name, "use the 'var' prefix")
  ```

### Accessors

- **`(*Diagnostic) Text() string`** — the message alone: no position, severity, code, notes, or caret. This is what an `Expected(error, "...")` test compares against; keeping it separate from rendering is what stops a source excerpt or column number from leaking into a normative comparison.
- **`(*Diagnostic) Error() string`** — same as `Text()`, so a `*Diagnostic` satisfies `error`.
- **`(*Diagnostic) Span() (pos, end token.Pos)`** — the underlined extent, normalizing an absent or inverted `End` to `Pos + 1`.

## Collection

### `Reporter`

```go
type Reporter interface {
	Report(*Diagnostic)
}
```

What a phase accepts. Scanner, parser, and analyzer all take a `Reporter` rather than a concrete `*List`, so an LSP can stream diagnostics one at a time and a batch build can collect them into a list, without either phase knowing which. `ReporterFunc` adapts a plain function to the interface.

### `List`

```go
type List struct { /* ... */ }
```

Collects diagnostics; the zero value is ready to use.

- **`(*List) Report(d *Diagnostic)`** — the `Reporter` implementation. A `nil` diagnostic is silently ignored, which lets a phase call `Report(maybeNilCheck())` without an extra branch.
- **`(*List) Add(code Code, pos, end token.Pos, args ...any) *Diagnostic`** — `New` plus `Report` in one call, returning the diagnostic so a caller can chain `WithNote`/`WithFixit` onto an already-recorded entry.
- **`(*List) Items() []*Diagnostic`**, **`Len() int`**, **`NumErrors() int`**, **`HasErrors() bool`** — basic accessors; `NumErrors` counts only `Error`-severity entries.
- **`(*List) Sort()`** — orders by position, then by code, stably. Positions are `FileSet`-global, so this also groups by file in the order files were added to the set.
- **`(*List) Dedup()`** — removes diagnostics identical in code, position, and message. Parser error recovery routinely re-reports at a resync point; deduping after the fact is cheaper than making every recovery path idempotent.
- **`(*List) Truncate(n int)`** — keeps the first `n` diagnostics and recomputes the error count. It deliberately does not synthesize a "too many errors" entry, since that text would be an unregistered message outside the `Code`/registry system; a caller that truncates is expected to say so itself.
- **`(*List) Err() error`** — returns the list itself as an `error` if it holds any `Error`-severity diagnostic, `nil` otherwise. Idiomatic use is `if err := list.Err(); err != nil { ... }`.
- **`(*List) Error() string`** — renders messages only, one per line, so a `*List` can be returned as a plain `error` without dragging in `Renderer` or a `FileSet`.

## Rendering

### `Renderer`

```go
type Renderer struct {
	Fset   *token.FileSet
	Source func(*token.File) []byte
	Color  bool
}
```

Produces human-facing output, deliberately kept separate from `Diagnostic` itself — the normative comparison in A.12.2 is always against `Text()`, never against anything `Renderer` produces. `Source` supplies a file's bytes for excerpting and may be `nil`, in which case diagnostics render without a source line or caret. `Color` enables ANSI escapes.

- **`(*Renderer) Render(w io.Writer, d *Diagnostic) error`** — writes a diagnostic as a severity/code header line, a `-->` file position, a source excerpt with carets under the span, then each note and fix-it in turn.
- **`(*Renderer) RenderList(w io.Writer, l *List) error`** — calls `Render` for every item in a `*List`, separated by a blank line.

Excerpting is internal: given a start/end span, the renderer locates the containing line via `Fset.File`, blanks everything before the span while preserving tabs (so the caret run aligns under any tab width without the renderer needing to know one), and clips a multi-line span's caret run to the first line.