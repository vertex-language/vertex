# token

```go
import "github.com/vertex-language/vertex/token"
```

Package `token` defines the lexical vocabulary and source-position machinery shared by every stage of the Vertex toolchain. It depends on nothing else — every other package (`scanner`, `parser`, `ast`, `diag`) sits above it.

## Design philosophy

### No statement terminator token

There is no `SEMICOLON` kind and no automatic semicolon insertion. A terminator is a *run* of line terminators, and a run of any length is one terminator — so line structure reaches the parser as a single flag (`Token.NLBefore`) on the token that follows, rather than as a token of its own.

A dedicated `NEWLINE` token wouldn't work here, because whether a line terminator actually terminates a statement depends on the innermost enclosing bracketing construct — significant inside a `Block`, ordinary white space inside a `LiteralValue` — and whether a given `{` opens one or the other is a *parsing* question, not a lexical one. So the scanner emits every line terminator and interprets none; the parser is the only phase equipped to decide what one means.

### Contextual keywords are just identifiers

Contextual keywords (`build`, `test`, `init`, `deinit`, `framework`, `module`, `blocks`, `threads`, `Expected`, `error`, and every `BuildTag` spelling) mint no `Kind` of their own. Each scans as plain `IDENT` and is recognized only at the one production that names it, via `Token.IsCtx` and the `Ctx*` string constants. This is why they're absent from `Lookup`: baking one into keyword lookup would make it reserved unconditionally, when in fact each is an ordinary name everywhere except its own recognition site.

## Package layout

| File | Contents |
|---|---|
| `token.go` | `Token`, contextual-keyword constants, type-operator names, and the reserved-builtin set |
| `kind.go` | `Kind`, its classification predicates, the keyword table, and the precedence ladder |
| `pos.go` | `Pos`, `Position`, `File`, and `FileSet` — the position machinery |
| `buildtag.go` | `BuildTag` and its lookup |

## Core types

### `Token`

```go
type Token struct {
	Kind     Kind
	Pos      Pos
	Lit      string // raw source text; empty when Kind has a fixed spelling
	NLBefore bool   // one or more line terminators precede this token
}
```

`Lit` is deliberately raw, not decoded — escape sequences and digit separators survive exactly as written, because a formatter needs the original spelling and decoding is a later phase's job. `Text()` returns `Lit` where present and falls back to `Kind.Spelling()` otherwise; `End()` derives a token's extent from `Text()`, which is why keywords and punctuation get correct extents without needing their own stored length.

### `Kind`

`Kind` values are assigned by `iota` in declaration order, with related kinds grouped into contiguous ranges bounded by unexported sentinels (`literalBeg`/`literalEnd`, `operatorBeg`/`operatorEnd`, `keywordBeg`/`keywordEnd`, etc.). Every classification predicate — `IsLiteral`, `IsOperator`, `IsKeyword`, `IsAssign`, `IsCompoundAssign` — is therefore a pair of comparisons rather than a switch. Renumbering is safe as long as each group's sentinels move together with it.

Notable predicates:

- **`Prec()`** returns one of seven binary-precedence levels. `as` (cast) is *not* among them despite `CastPrec` being defined above `UnaryPrec` — `as`'s right operand is a `Type`, not an `Expr`, so a precedence-climbing loop can't consume it and `Prec` never returns `CastPrec` for it. Casts are folded by the parser as their own production instead.
- **`EndsOperand()`** backs the *one* deliberate exception to longest-match scanning: a float literal may not begin immediately after a `.` whose preceding token also satisfies this predicate. This is what lets `t.0.0` scan as two tuple-index accesses instead of `t` `.` `0.0`. The set is narrow by design (`IDENT`, `)`, `]`, `}`, `INT`, `FLOAT`, `STRING`) and must not be generalized — `char_lit` and `true`/`false`/`nil` are deliberately excluded, so `1.5` still scans as a plain float.
- **`IsNonAssociative()`** is true only for `DOTDOT`: `a..b..c` is a compile error that a precedence table alone can't express, so the parser checks this explicitly.

### `Pos` / `Position` / `File` / `FileSet`

`Pos` is a compact integer offset into a `FileSet`'s shared global address space — `NoPos` (zero) means "no position." `Position` is what you get after resolving a `Pos` against its file: filename, byte offset, 1-based line, 1-based column.

```go
fset := token.NewFileSet()
f := fset.AddFile("main.vs", len(src))
pos := f.Pos(offset)
```

`FileSet` exists so that every file in one compilation shares one position space — a diagnostic can span two files without needing to carry a file reference alongside each position. `AddFile` leaves a one-`Pos` gap after each file so a one-past-the-end position can never collide with the next file's first character.

`File.lines` is mutex-guarded because `AddLine` runs on the scanner's goroutine while `Position`/`LineCount` run on whatever goroutine renders a diagnostic — a driver may stream diagnostics as they're produced, concurrently with scanning.

### `BuildTag`

`BuildTag` lives here — not in `parser` — because the parser, the loader, and the driver all need it, and because it changes what's admissible: a file's tag is what licenses an `ExpectedType` result (see `LicensesTest`). `LookupBuildTag`'s second return value is load-bearing: an unrecognized tag is a compile error, never a silently excluded file, so a caller must be able to distinguish "unknown spelling" from `TagNone`'s "no clause at all" — something a bare zero-value return couldn't express.

## Reserved names

Two disjoint sets of "special" identifiers live in `token.go`:

- **Type operators** (`sizeof`, `alignof`, `reinterpret`) — the only call forms that take a `Type` in argument position. The parser recognizes them by name via `IsTypeOperator`, which is sound only because reserved builtin names may not be shadowed.
- **Reserved builtins** (`new`, `delete`, `resize`, `copy`, `zero`, `addr`, `sizeof`, `alignof`, `reinterpret`, `upgrade`, `drop`, `panic`, `blend`, `min`, `max`, `clamp`, `transfer`) — pre-bound in the implicit outermost scope and checked via `IsReservedBuiltin`. Notably, `transfer` is reserved and bound to *nothing*, purely so that `x.transfer()` is diagnosed as a misspelled ownership marker (the real syntax is the `var` prefix) rather than as an ordinary unknown-name error.