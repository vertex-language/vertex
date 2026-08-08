# token

```go
import "github.com/vertex-language/vertex/token"
```

Package `token` defines the lexical vocabulary of Vertex: token kinds, contextual-keyword identity, source positions, and diagnostics. It imports nothing from `scanner`, `ast`, or `parser`, and it never will — it is the bottom of the dependency graph, and `Diagnostic` lives here rather than in `scanner` precisely because both `scanner` and `parser` need to emit them and neither can import the other.

## Token

```go
type Token struct {
	Kind  Kind
	Ctx   Ctx
	Flags Flags
	Pos   Pos
	End   Pos
}
```

A `Token` is a lexeme: a classification and a span. It carries no literal text — recover that with `File.Slice`, since the scanner allocates no strings. The struct is pointer-free and packs to twelve bytes, which matters because the parser buffers whole files as a `[]Token`; a pointer-bearing element would put that buffer under GC scan.

`Adjacent` reports whether one token begins exactly where another ends, with no intervening whitespace or comment — this is the test the expression parser uses when joining a run of `GT` tokens, which is why `a > > b` never becomes a shift. `IsContextual` reports whether a token spells a given contextual keyword; whether it *declares* anything at that position is a parse question, not a token question.

### Flags

`Flags` records lexical facts with no home in `Kind`:

- `NLBefore` marks a line terminator between the previous token and this one. A newline is never a token of its own — whether it ends a statement is answered later by `expectSemi` and the grammar's `[no LineTerminator here]` restrictions. A comment spanning a newline sets this on the next real token regardless of whether comments are being retained.
- `HasEscape` marks an identifier containing a Unicode escape, a string or template containing an escape sequence, or a numeric literal containing a separator. Keyword lookup must be skipped for an escaped identifier — `\u0069f` is an `IDENT` named `if`, never the `IF` keyword — and this flag is what tells the scanner to skip the `LookupIdent` call.
- `Unterminated` marks a string, template, comment, or regex that ran to end of input or end of line. The token still gets an exact span, so recovery can proceed and every node keeps a non-zero span.

## Kind

`Kind` classifies a token and is deliberately `uint8`. Kinds are laid out in contiguous ranges bounded by unexported sentinels, so `IsLiteral`, `IsOperator`, and `IsReserved` are range comparisons rather than switches — the ranges *are* the API, and reordering constants without updating the sentinels breaks them silently.

Only `ReservedWord` gets a `Kind`. Contextual keywords and strict reserved words scan as `IDENT` and carry their identity in `Ctx` instead (see below).

Five operator kinds — `GEQ`, `SHR`, `USHR`, `SHR_ASSIGN`, `USHR_ASSIGN` — exist in the table but are never produced by the scanner; `ScannerEmits` reports this so golden scanner tests can assert no fixture ever contains one. They exist because the scanner deliberately under-munches `>` so that `Array<Box<int32>>` tokenizes as two `GT`s in type context, and the expression parser's `JoinGT` reassembles a run of adjacent `GT`s (plus an optional trailing `=`) back into the joined form when it's looking for a binary operator. Adjacency, not just sequence, is the test — `Token.Adjacent` is what keeps `a > > b` from joining.

`Precedence` gives the binary precedence of a `Kind` for the precedence-climbing parser, including precedences for the joined forms even though the scanner can't emit them directly. `CoalesceExpression` and `LogicalORExpression` share a precedence level because the grammar forbids mixing `??` with `||`/`&&` in one chain — a fact a precedence table can't express, so the parser records the mix and rejects it later instead.

## Ctx

`Ctx` is the identity of a contextual keyword. It's zero (`CtxNone`) unless `Kind == IDENT`. This is how contextual keywords stay contextual: the token remains an `IDENT` and stays usable as a binding name, while the parser tests `tok.Ctx == CtxStruct` in O(1) with no string comparison anywhere downstream.

`Ctx` also holds the grammar's `StrictReservedWord` list (`implements`, `interface`, `let`, `package`, `private`, `protected`, `public`, `static`), which has no grammar productions of its own — strict-mode restrictions are early errors, not grammar, so those words are `Ctx` values rather than `Kind` values. Keeping them here lets `ClassElementModifier` be a single contiguous range test instead of a mix of `Kind` and `Ctx` checks. `let` in particular must stay unconditionally non-reserved, since `ExpressionStatement`'s lookahead restriction (`∉ { let [ }`) depends on being able to see it as a plain `IDENT` with `CtxLet`.

`identTable` is the single lookup for every `IdentifierName`, built once at `init` from both `kindNames` and `ctxNames` so a word can't be added to one table and forgotten in the other — a duplicate entry panics at startup rather than silently shadowing. `LookupIdent` classifies a name into a reserved `Kind` with `CtxNone`, or `IDENT` with a `Ctx`, or plain `IDENT` for an ordinary name; callers must not invoke it for an identifier whose source spelling contained an escape, which is `Flags.HasEscape`'s job to flag upstream.

## Pos and File

```go
type Pos uint32
```

`Pos` is a 1-based byte offset into one translation unit, in a per-unit address space. There is deliberately no `FileSet` — a global position space across files would be a cross-file dependency at the one layer that must not have any, so positions travel with their `File`. The offset is biased by one so `NoPos == 0` is distinguishable from a file's first byte; only `File.Slice`, `File.Between`, and `File.Position` unbias, each at its own boundary.

`File` is one translation unit: a name, its bytes, and the arithmetic to turn a `Pos` back into text or a line and column. It's safe for concurrent use once constructed, since files may be parsed in parallel — one `File` per goroutine — with diagnostics rendered afterward from any goroutine. The line index backing `Position` is built lazily on first use rather than in `NewFile`, because most files are parsed without any diagnostic ever being rendered.

`Slice` returns a token's raw bytes with no decoding — `1_024` yields exactly those five bytes, separators and all, because decoding belongs to a later phase that knows the target type. `Between` takes two `Pos` values rather than an `ast.Node`, on purpose: `token` must not import `ast`.

## Diagnostic

```go
type Diagnostic struct {
	Pos Pos
	End Pos
	Msg string
}
```

`Diagnostic` lives in `token` rather than `scanner` or `parser` because both of those packages emit them, and either choice of home would otherwise need a shared package above `token`. It carries no severity and no code — a cap on reported errors, and any grouping or coloring, are presentation decisions that belong to whatever renders the diagnostics, not to this type. `End` is `NoPos` for a point diagnostic; `IsSpan` reports whether it instead covers a range.

`SortDiagnostics` orders by `Pos`, then `End`, then `Msg`, so output stays reproducible when two diagnostics share a position. It does not define ordering across files, and can't: a bare `Pos` has no meaning outside its own unit, so a caller holding several units must key on `(unit index, Pos)` itself.