# diag

```go
import "github.com/vertex-language/vertex/diag"
```

Package `diag` is the shared diagnostic currency for the Vertex toolchain. Scanner, parser, analyzer, and vir all produce `*Diagnostic` values through this package rather than formatting their own error strings.

`diag` depends on `token` for `Pos` and `FileSet`, and on nothing else.

## Design philosophy

A rejection has one shape everywhere: a stable `Code`, a span, a message rendered from a registered template, and optionally secondary spans and machine-applicable edits.

The template registry is what keeps one rule's message identical across every site that raises it. This matters beyond consistency: it's what makes the error form of `grammar.md`'s `ExpectedType` — whose second operand is a message string — stable enough to serve as specification. A test can assert `Expected(error, "expected a type, found ...")` and trust that text never drifts based on which call site happened to raise it.

**Rendering is a separate concern from the diagnostic itself.** The normative comparison is always against `Diagnostic.Text()`, never against anything `Renderer` produces. This is enforced structurally: `Text()` returns the message alone, stripped of position, severity, code, notes, and caret — so a source excerpt or column number can never leak into a spec-level comparison.

## Package layout

| File | Contents |
|---|---|
| `diag.go` | `Diagnostic`, `Severity`, constructors (`New`, `At`, `AtToken`, `AtNode`), and the `WithNote`/`WithFixit` builder methods |
| `code.go` | `Code`, the numeric ranges, the template registry, and the conformance list `declaredCodes` |
| `list.go` | `Reporter`, `List` (a diagnostic sink with sort/dedup/truncate), and `ReporterFunc` |
| `render.go` | `Renderer`, which produces human-facing output with source excerpts and carets |

## Core type

```go
type Diagnostic struct {
	Code Code
	Sev  Severity

	Pos token.Pos // where the diagnostic points
	End token.Pos // bounds the underlined span

	Message string // the rendered template — the normative text
	Notes   []Note
	Fixits  []Fixit
}
```

`Notes` carries secondary spans — a use-after-transfer needs two: where the binding died, and where it was subsequently read. `Fixits` carries machine-applicable edits, present from the start because grammar.md already commits to one by name — `transfer` is reserved and bound to nothing so that x.transfer() diagnoses as a misspelled ownership marker, and it's cheaper to carry an empty slice now than to thread a new field through every constructor later.

## Constructing diagnostics

```go
func New(code Code, pos, end token.Pos, args ...any) *Diagnostic
func At(code Code, pos token.Pos, args ...any) *Diagnostic       // zero-width position
func AtToken(code Code, t token.Token, args ...any) *Diagnostic  // over a token's full extent
func AtNode(code Code, pos, end token.Pos, args ...any) *Diagnostic
```

Call sites never format their own text — `New` looks up the code's registered template and formats it with `args`. `AtNode` deliberately takes `pos, end token.Pos` rather than an `ast.Node`, so that `diag` does not depend on `ast`. The dependency runs the other way; a cycle here would be structural.

Builder methods attach secondary information and return the diagnostic for chaining:

```go
d := list.Add(diag.UndeclaredName, pos, end, name)
d.WithNote(declPos, declEnd, "did you mean %s?", suggestion)
d.WithFixit(pos, end, "var "+name, "use the 'var' prefix")
```

`WithInsert` and `WithDelete` are convenience wrappers over `WithFixit` for an empty span and empty replacement text, respectively.

## Codes and the registry

`Code` is a stable, explicitly-numbered handle — never assigned by `iota`, so inserting a new code can't shift another's value. Ranges follow grammar.md's own section order: 0xxx internal, 1xxx lexical elements, 2xxx syntax errors, 3xxx types, 4xxx expressions, 5xxx statements, 6xxx declarations and names, 7xxx generics, 8xxx declare blocks, 9xxx test result types.

Only the lexical (1xxx) and syntax (2xxx) ranges are densely populated today, since the scanner and the parser are the only phases that exist; every other range holds the forms grammar.md names as parsing-and-then-rejected, waiting for the analyzer pass that raises them.

```go
func (c Code) String() string     // "V1003" — the stable identifier to pin in tests
func (c Code) Severity() Severity // defaults to Error if unregistered
func (c Code) Template() string   // the raw format string
func (c Code) Registered() bool
func Codes() []Code               // every declared code, in declaration order
```

`declaredCodes` is a hand-maintained list backing a conformance test, since Go has no reflection over untyped constants, and a code declared but never registered would surface as a silent Internal at the first call site that used it. The test checks that every declared code has a registry entry, no two codes share a numeric value, and no two share a message template.

## Collecting diagnostics

```go
type Reporter interface {
	Report(*Diagnostic)
}
```

Scanner, parser, and analyzer all accept a `Reporter` rather than a concrete `*List` — so an LSP can stream diagnostics and a batch build can collect them, without either phase knowing which.

`List` is the standard collector (its zero value is ready to use):

```go
l := &diag.List{}
l.Add(diag.ExpectedToken, pos, end, "')'", "','")

l.Sort()             // by position, then code, stably
l.Dedup()             // drop identical (code, pos, message) triples
l.Truncate(100)       // keep the first n
if err := l.Err(); err != nil { ... }
```

`Sort` groups by file in `AddFile` order, since positions are `FileSet`-global. `Dedup` exists because parser recovery routinely re-reports at a resync point; this is cheaper than making every recovery path idempotent. `Truncate` deliberately does not synthesize a "too many errors" entry — that text would be an unregistered message.

`List` also implements `error`, rendering one message per line, so it can be returned directly from a phase without dragging in the renderer or a `FileSet`.

## Rendering

`Renderer` turns a `Diagnostic` or `List` into human-facing text, with an optional source excerpt and caret run under the offending span:

```go
r := &diag.Renderer{
	Fset:   fset,
	Source: func(f *token.File) []byte { return sourceBytes[f.Name()] },
	Color:  true,
}
r.Render(os.Stderr, d)
r.RenderList(os.Stderr, list)
```

`Source` may be `nil`, in which case diagnostics render without a source line or caret. The caret run is built by echoing the line's leading whitespace — preserving tabs — so it aligns correctly under any tab width without the renderer needing to know one.