# scanner

```go
import "github.com/vertex-language/vertex/scanner"
```

Package `scanner` converts Vertex source bytes into tokens. It depends on `diag` for diagnostics and `token` for the token vocabulary and `FileSet`/`File`.

## Design philosophy

### Line structure is recorded, never interpreted

The scanner emits every line terminator and interprets none of them. Whether a line terminator actually *ends* a statement depends on the innermost enclosing bracketing construct — a question the parser answers, not the scanner. So the scanner's only contribution is `token.Token.NLBefore`: a single bool flag meaning "a line terminator preceded this token."

There is no `SEMICOLON` kind, no synthesized `NEWLINE`, and no automatic semicolon insertion. A run of consecutive line terminators collapses into that one bool, so no consumer ever needs a count.

This split matters at the edges too: a line terminator inside a raw string is part of that string's *value* and must not end a statement, though it still needs to advance the line table (for correct position reporting later). That's why `consumeLineTerminator` (updates the line table) and the `nlPending` flag (signals statement-ending) are kept as separate concerns — `scanRawString` calls the former without ever setting the latter.

### Keywords vs. everything else

The scanner recognizes only the reserved keywords and the three reserved literal keywords (`true`, `false`, `nil`), via `token.Lookup`, and nothing else. Contextual keywords, predeclared type names, tensor-element names, constraint names, and reserved builtin names are all ordinary identifiers at this layer — they're pre-bound in an implicit outermost scope that the scanner must not know about. That resolution is the parser's and analyzer's job.

## Package layout

| File | Contents |
|---|---|
| `scanner.go` | `Scanner`, `Init`, the main `Scan` dispatch loop, whitespace/line handling, and operator/punctuation recognition |
| `ident.go` | Unicode `ID_Start`/`ID_Continue` classification and `scanIdent` |
| `number.go` | `scanNumber`, `scanDigits` (with separator rules), and `scanTupleIndex` |
| `string.go` | `scanString`, `scanRawString`, `scanChar`, and `scanEscape` |

## Core type

```go
type Scanner struct {
	// ... unexported fields
	ErrorCount int
}

func (s *Scanner) Init(file *token.File, src []byte, rep diag.Reporter, mode Mode)
func (s *Scanner) Scan() token.Token
```

`Init` panics if `file.Size()` doesn't match `len(src)` — the two must already agree before scanning starts. `Scan` returns tokens one at a time; at end of input it returns `EOF` indefinitely, with `NLBefore` still reflecting any trailing line terminator.

```go
const ScanComments Mode = 1 << iota
```

By default comments are consumed and discarded (though they still participate in line structure — see below). `ScanComments` makes them come out as `COMMENT` tokens instead.

```go
func Tokenize(file *token.File, src []byte, rep diag.Reporter, mode Mode) []token.Token
```

A convenience wrapper that drains a `Scanner` to completion, including the trailing `EOF`. Intended for tests and tooling; the parser drives `Scan` directly rather than pre-materializing a slice.

## Notable rules

- **Comments and line structure.** A comment participates in line-terminator logic regardless of whether it's surfaced as a token. A line comment always acts as a terminator; a general (`/* */`) comment acts as one only if it actually spans a line break — `/* inline */` is white space, but a comment containing a newline is not.
- **The selector-dot / tuple-index restriction.** `afterSelectorDot` is the one deliberate exception to longest-match scanning: if the previous token is `.` and the one before *that* can end an operand, a following digit run scans as a plain `int_lit` via `scanTupleIndex` rather than being folded into a float. This is what lets `t.0.0` scan as a chain of tuple accesses instead of `t` `.` `0.0`.
- **`$` in identifiers.** `$` is not a legal identifier character in any position, but rather than aborting the scan, `scanIdent` consumes a run containing one and reports it once over the whole span — so a pasted foreign identifier produces one diagnostic and one token, not a cascade of fragments.
- **`_` as digit separator.** `scanDigits` enforces that `_` may appear only between digits: not leading, trailing, or doubled. An out-of-range decimal digit in a narrow base (e.g. `0b1234`) is consumed and diagnosed per-digit rather than ending the run early, so one bad literal produces one coherent diagnostic rather than a spurious second token.
- **NUL rejection.** `next()` reports `NulCharacter` the moment a `\x00` byte is read, wherever it occurs — inside a string, a comment, anywhere — because Vertex source must never contain a NUL, not merely code outside comments/strings.
- **Raw text preservation.** String, char, and number literals all keep their *raw source spelling* — quotes, escapes, digit separators, everything — in `Token.Lit`. Decoding (resolving `\n`, computing a numeric value) is explicitly deferred to a later phase; a formatter needs the original text and would otherwise have to reconstruct it.