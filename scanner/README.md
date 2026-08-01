# scanner

`github.com/vertex-language/vertex/scanner`

`scanner` converts Vertex source bytes into `token.Token`s. It depends on `token` and `diag`; nothing else in the toolchain depends on `scanner` except `parser`.

## Entry points

### `Scanner`

```go
type Scanner struct { /* ... */ }
```

The scanner is a value type driven by `Init` then repeated `Scan` calls; there is no constructor.

### `(*Scanner) Init`

```go
func (s *Scanner) Init(file *token.File, src []byte, rep diag.Reporter, mode Mode)
```

Prepares `s` to scan `src`, which must be the contents of `file` — `Init` panics if `file.Size()` doesn't match `len(src)`. Resets all scanning state, including `ErrorCount`, and primes `s.ch` with the first character via an initial `next()`.

### `(*Scanner) Scan`

```go
func (s *Scanner) Scan() token.Token
```

Returns the next token. At end of input it returns `EOF` indefinitely, with `NLBefore` still reflecting any trailing line terminator. A stray character that matches no production is reported, consumed, and returned as `token.INVALID` rather than aborting the scan — recovery is the scanner's default, not an exception.

### `Mode`

```go
type Mode uint

const (
	ScanComments Mode = 1 << iota
)
```

`ScanComments` makes `Scan` return `COMMENT` tokens instead of discarding them. A discarded comment still affects `NLBefore` either way (A.1.1) — the mode controls whether the token is surfaced, not whether it participates in line structure.

### `Tokenize`

```go
func Tokenize(file *token.File, src []byte, rep diag.Reporter, mode Mode) []token.Token
```

Convenience for tests and tooling that scans `src` to completion, including the trailing `EOF`. The parser does not use this — it drives `Scan` directly, token by token.

## Design

### Line structure is recorded, not interpreted

The scanner never suppresses a line terminator and never synthesizes one. A.0.6 makes a line break inside `{ }` ordinary whitespace, but statements also live inside `{ }` blocks and end at a line terminator — whether a given brace opens a `Block` or a `CompositeLiteral` is exactly the `[+Lit]` question the parser tracks, not something the scanner can resolve. So line structure is only ever recorded, on `token.Token.NLBefore`, and the parser decides what it means.

`nlPending` accumulates this: `skipSpace` sets it when it consumes a `LineTerminator`, a comment that spans a line sets it directly (a multi-line comment is itself a line terminator for A.0.6's purposes per A.1.1), and `token()` consumes it into the next-built token's `NLBefore` before clearing it. `COMMENT` tokens are the one exception — `token()` deliberately does not update `s.prev` for them, so a comment sitting between `.` and a digit doesn't disturb the tuple-index rule below, but a comment's own `NLBefore` is still folded correctly into whatever follows.

A line terminator inside a raw string is part of that string's value (A.1.5.2) and must not set `nlPending`. That's why `consumeLineTerminator` (advances the line table) and `skipSpace`'s call to it (also sets `nlPending`) are kept as separate steps: `scanRawString` calls the former directly and skips the latter.

### `prev` and the tuple-index rule

The scanner tracks only one piece of token-kind state across calls: `s.prev`, the kind of the last non-comment token. Its sole use is A.4.3 — a `.` immediately followed by a digit is always positional tuple access, never a float, since `EnumShorthand` takes an identifier and Vertex has no leading-dot float literals. `Scan` checks `s.prev == token.PERIOD` before falling into digit scanning and routes to `scanTupleIndex` instead of `scanNumber` when it matches. Splitting a `FLOAT` apart in the parser was rejected in favor of this: scanning the tuple index directly is what makes a chain like `t.0.0` fall out naturally, since the second `.` simply ends the run and the rule fires again on the next token.

### Identifiers (`ident.go`)

`isIdentStart`/`isIdentPart` implement `IdentifierStart`/`IdentifierPart` (A.1.2): ASCII fast-path, falling back for non-ASCII to `isIDStart`/`isIDContinue`, which assemble Unicode's `ID_Start`/`ID_Continue` properties from the pieces Go's `unicode` package exposes (`L`, `Nl`, `Other_ID_Start`, etc., minus `Pattern_Syntax`/`Pattern_White_Space`) since the derived properties themselves aren't exposed directly. Both are marked as candidates for a generated table pinned to a specific Unicode version, since deriving identifier rules from the host toolchain's stdlib means the grammar drifts with it.

`$` is deliberately excluded from `isIdentPart` — A.1.2 makes it a diagnosed error (`DollarInIdent`), not a silently-unmatched character, so `Scan` gives it its own case that consumes a whole would-be identifier and reports once.

`scanIdent` just consumes an `IdentifierName`; keyword classification happens in `Scan` via `token.Lookup`, and the blank identifier `_` falls out as an ordinary `IDENT`, which `ast.Ident.IsBlank` tests later.

### Numbers (`number.go`)

`scanNumber` implements A.1.5.1's `NumericLiteral` grammar directly rather than through a generic digit-run helper composed blindly: base prefixes (`0b`/`0o`/`0x`) are detected up front, then each base has its own tail handling, because the three bases don't share a shape — only hex has a binary-exponent-only float form, only decimal has a bare exponent, and both hex and decimal define "float" differently.

Two decisions here are load-bearing beyond `scanNumber` itself:

- A fractional part must be non-empty, so `1.` is not a literal — that's what makes `1..5` scan as `INT DOTDOT INT` without lookahead surgery in the scanner.
- A hexadecimal float requires its binary exponent (A.1.5.1), so `0xC.3` is diagnosed as an incomplete float (`HexFloatNoExponent`) rather than silently accepted or split into three tokens.

`scanDigits` enforces A.1.5.1's separator rules for a single digit run in a given base: `_` may not lead a run, trail it, or double up, each reported with its own code (`SeparatorLeads`/`SeparatorDoubled`/`SeparatorTrails`). A decimal digit out of range for a narrow base (e.g. `9` in `0b1234`) is consumed rather than ending the run, so one bad literal produces one diagnostic per bad digit instead of a literal followed by a spurious second one.

Both `scanNumber` and `scanTupleIndex` end by consuming any identifier characters that immediately follow the literal and reporting `NumberJoinedToIdent` over that span — `1abc` is one diagnosed token, not `INT` followed by `IDENT`.

### Strings, chars, and escapes (`string.go`)

`scanString` and `scanChar` return the raw source spelling, quotes and escapes included exactly as written — `vfmt` needs the original text, and decoding is the analyzer's job, not the scanner's. A `LineTerminator` is excluded from `DoubleStringCharacter` and is left unconsumed on error so it still ends the statement (A.0.6) and recovery resumes on the next line rather than swallowing it.

`scanRawString` is the opposite on both counts: no escape sequence is recognized, and every line terminator it spans is part of its value (A.1.5.2), so it calls `consumeLineTerminator` directly and continues rather than erroring or setting `nlPending`.

`scanChar` counts scalar values in source units (an escape counts once, a multi-byte rune counts once) and, once a literal terminates, reports `EmptyChar` or `CharTooLong` depending on that count — a char denotes exactly one Unicode scalar value (A.1.5.2).

`scanEscape` is entered with the backslash already consumed. It recognizes the fixed single-character escapes, `\xNN` (exactly two hex digits, erroring via `InvalidHexEscape` on the first bad one), and `\u{...}` (erroring via `InvalidUnicodeEscape` on a missing brace or empty digit run, and `UnicodeEscapeRange` separately for a value that's a valid code point but not a scalar — out of range or a surrogate). An unrecognized escape reports `InvalidEscape` and consumes the character, except at `eof` or a line terminator, where it returns without reporting since the enclosing literal will report its own unterminated error.

### Errors never abort a scan

`error`/`errorSpan` both increment `s.ErrorCount` and, if a `Reporter` is set, report through `diag`. Every scanning function that hits an error keeps going and returns a token anyway — an unterminated string still returns the text scanned so far, an illegal character is still consumed and returned as `token.INVALID` — so `Scan` never needs to signal failure out-of-band and a single bad file always yields a full token stream.

### File organization

- **`scanner.go`** — the `Scanner` struct, `Init`, `Scan`'s main dispatch, punctuator/operator recognition, comment scanning, line-terminator and whitespace handling, and `Tokenize`.
- **`ident.go`** — `IdentifierStart`/`IdentifierPart` classification and `scanIdent`.
- **`number.go`** — `scanNumber`, `scanDigits`, and `scanTupleIndex` for A.1.5.1's numeric literals and A.4.3's tuple-index digits.
- **`string.go`** — `scanString`, `scanRawString`, `scanChar`, and `scanEscape` for A.1.5.2's textual literals.