# scanner

```go
import "github.com/vertex-language/vertex/scanner"
```

Package `scanner` turns Vertex source into a token buffer. It records structure and interprets none of it: contextual keywords stay identifiers, literals keep their raw spelling, and a newline is a flag rather than a token of its own. The whole file is tokenized up front, because speculation in the parser needs an immutable buffer and O(1) rollback — there is no streaming or callback-driven mode.

`scanner` imports only `token`.

## Entry point

```go
toks, diags := scanner.Scan(file)
```

`Scan` always returns a token slice ending in EOF, and diagnostics sorted by position. Comments are always emitted as `COMMENT` tokens regardless of whether anything downstream wants them; whether they survive into a tree is the parser's `ParseComments` decision, and a scanner that dropped them could not serve highlighters, while `ParseFileTokens` must be able to accept a buffer from either kind of caller.

## The one context-dependent production

Everything else in the lexical grammar is context-free. `/` is not: it opens a line comment, a block comment, a regex, or a division/assignment operator depending on what came before it. `scanSlash` handles the comment cases directly by peeking the next byte, then falls back to `regexAllowed`, which implements the previous-significant-token rule: the question is whether a value has just been produced. After a value, `/` divides; otherwise it opens a RegularExpressionLiteral.

That rule needs to know, for a closing `)` or `}`, what kind of opener it matched — `if (x) /re/` is a regex while `f(x) / 2` is division — which is what the frame stack (below) is for.

## The frame stack

One stack does three jobs that could be described separately but can't be *implemented* separately: each `(` records whether it followed a control keyword, each `{` records block vs. object literal, and each `${` opens a template substitution. They share a stack because the `}` in `` `a${b}c` `` must not be read as closing a block, and that decision needs the same depth counter as the other two.

`controlParen` decides whether a `(` opens a control-flow head rather than a call or parameter list, by checking whether the previous token was `if`, `while`, `for`, `switch`, `catch`, or `with`. `braceOpensBlock` classifies a `{` as a `Block` or an `ObjectLiteral`; this is explicitly a heuristic — misclassification costs one recoverable diagnostic. The blast radius is narrow: the classification is read back only when a `/` immediately follows the matching `}`, so getting it wrong turns one regex into one division and nothing else.

A mismatched closer — `popIf` failing to find the frame kind it expected — doesn't unwind the stack: the closer is still emitted, the parser's recovery handles the structure, and "one recoverable diagnostic, never a cascade" is preserved because the stack does not unwind past a real opener.

## Numbers

`scanNumber` finds the extent of a numeric literal and validates its shape in the same walk, since a malformed separator changes where the literal ends. It never decodes — `1_024` yields the bytes as written, with no value, no width, and no separator stripping; decoding belongs to a phase that knows the target type. There's deliberately no legacy octal production: `0123` is not a literal, and leading zeros are diagnosed rather than silently accepted. A literal immediately followed by an identifier character is also a diagnostic rather than two adjacent tokens — `3in` and `0x1p` would otherwise scan cleanly and fail later with a confusing parse error instead of a lexical one.

## Identifiers

`scanIdent` scans an `IdentifierName` and classifies it in one call to `token.LookupIdent`. An identifier containing a Unicode escape is never looked up as a keyword — an escaped spelling never matches a keyword: `\u0069f` is an identifier named "if," not the IF token; this is the whole reason `HasEscape` exists on identifiers. `scanUnicodeEscape` consumes but does not validate the escape's code point: that is a later phase's job, and rejecting here would mean the scanner knows what an identifier means rather than where it ends. Predeclared type names and every `PredefinedType` member scan as ordinary identifiers here — the scanner has no notion of which identifiers denote types.

## Strings, templates, and regexes

`scanString` records escapes without decoding them, same as numbers, and gives an unterminated string a real token with an exact span rather than failing outright, so recovery has something to work with.

`scanTemplateBody` scans `TemplateCharacters` up to the next delimiter and reports which of `TEMPLATE`, `TEMPLATE_HEAD`, `TEMPLATE_MIDDLE`, or `TEMPLATE_TAIL` it produced; a `${` pushes a `frameSubst` frame, and the matching `}` is routed back into this function by `scanPunct` rather than through a separate lexer mode. The same tokens are reused for type-position templates: the grammar has no TemplateTypeHead precisely so there is one way to tokenize `` `hello ${ ``  — this function does not know or care which side of the language it is on.

`scanRegex` delimits rather than parses the body — Pattern and everything under it is a separate grammar that a later phase compiles; all this needs to know is where the closing `/` is, which means tracking character classes, since a `/` inside `[...]` is literal. A raw line terminator ends a regex (or a string) as unterminated rather than being consumed, keeping the misclassification cost to one diagnostic instead of a cascade.

## Whitespace and newlines

`skipSpace` recognizes ASCII space, tab, and the usual Unicode whitespace and line-terminator code points, including NBSP, ZWNBSP, the ideographic space, and `LS`/`PS`. A line terminator doesn't become a token — it sets `nlPending`, which `emit` folds into the next real token's `NLBefore` flag. That flag deliberately survives an intervening comment: `x ⏎ //c ⏎ y` must still leave `NLBefore` set on `y`, since the parser skips comments and ASI would otherwise never see the break. A block comment spanning multiple lines sets the same flag on whatever follows, whether or not comments are being retained downstream.

## Diagnostics

Scanning never stops early. Every scan function that can fail — an unterminated string, an unterminated regex, a malformed numeric separator, an unrecognized character — records a `token.Diagnostic` and still emits a token with a real, non-empty span, so recovery in the parser always has a token stream to work with rather than a truncated one. `Scan` sorts diagnostics by position before returning them, per `token.SortDiagnostics`.

## Directives

```go
diags := scanner.Directives(file, re)
```

`Directives` collects comments matching re and returns them as expected diagnostics, keyed by the 1-based line of the comment. It lives in `scanner` rather than a test-only package so that fixtures outside the repo can use it. Each match is deliberately positioned at the preceding significant token rather than at the comment itself — that is what makes `let x: = 1; // ERROR "expected type"` point at the `=` rather than at the comment. The result is a map rather than a sorted slice, so callers that render or diff output must sort the keys themselves.