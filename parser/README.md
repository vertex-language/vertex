# parser

`github.com/vertex-language/vertex/parser`

`parser` turns Vertex source into an `ast.File` or `ast.Package`. It depends on `ast`, `diag`, `scanner`, and `token`; nothing else in the toolchain depends on `parser`.

## Entry points

### `ParseFile`

```go
func ParseFile(fset *token.FileSet, filename string, src []byte, rep diag.Reporter, mode Mode) (*ast.File, error)
```

Parses a single `.vs` source file. `src` may be `nil`, in which case `filename` is read from disk. The returned `*ast.File` is non-nil even when errors were reported — error recovery produces `Bad*` nodes so later phases and editor tooling can still walk a partial tree.

`rep` controls how diagnostics leave the function:

- If `rep` is `nil`, the parser collects into an internal `diag.List`, sorts and dedups it, and returns it as the `error` result via `list.Err()`.
- If `rep` is non-nil, diagnostics stream to it directly as they're produced, and the returned `error` is a plain `fmt.Errorf("%s: %d errors", ...)` summarizing the parser's and scanner's error counts — the caller is assumed to already have the full diagnostics from `rep`.

### `ParseDir`

```go
func ParseDir(fset *token.FileSet, dir, importPath string, target token.BuildTag, rep diag.Reporter, mode Mode) (*ast.Package, error)
```

Parses every `.vs` file in `dir` that targets `target` and groups them into one `*ast.Package`. Runs in two passes rather than opportunistically:

1. **Pass 1** — every file is parsed with `PackageClauseOnly` into a throwaway `*token.FileSet`, just far enough to read its build clause. A.2.2 excludes a non-matching file from the build outright, so the target filter has to run before any file is fully parsed.
2. **Pass 2** — the surviving files are fully parsed against the real `fset` and handed to `ast.NewPackage`.

The two-pass split also serves A.2.3: the imported package's own `PackageClause` supplies the qualifier, so a file's names can't be resolved until every surviving file's package line has been read.

### `Mode`

```go
type Mode uint

const (
	PackageClauseOnly Mode = 1 << iota
	ImportsOnly
	ParseComments
)
```

- **`PackageClauseOnly`** — stop after the package and build clauses. This isn't just an optimization; it's load-bearing for `ParseDir`'s first pass.
- **`ImportsOnly`** — stop after the import declarations.
- **`ParseComments`** — retain comments in the tree (as `ast.CommentGroup`s on `File.Comments` and on the declarations/fields they lead or trail).

## Design

### Context parameters

A.0.2 defines four context parameters for the grammar. The parser tracks exactly one of them as state. `Await`, `Npu`, and `Own` name forms that parse unconditionally and are rejected by a later phase — `await` outside an async body, a launch prefix outside its device context, `var` outside an owning position — and since all three are spelled with dedicated keywords, parsing them without tracking context is unambiguous and still lets the resulting diagnostic quote the construct rather than a bare token position.

`Lit` is different: A.4.7 requires a literal used in a control-flow header to be parenthesized, which is a decision the parser has to make while parsing, not something a later pass can retrofit. It's carried as `noLit` on the parser and toggled by `parseHeaderExpr` (sets it) and `withLit` (clears it inside any bracketed group — parentheses, index brackets, argument lists — since A.4.1 makes parentheses the escape hatch and A.4.3 gives index brackets `Expression[+Lit]` explicitly).

### Statement termination

A.0.6 is driven by a bracket-nesting depth (`p.depth`) rather than by a dedicated token. `(`, `[`, and a literal `{` push depth via `open`/`close`; a `Block`'s `{` resets depth to zero instead, because statements inside a block *do* end at line terminators while entries inside a bracketed list do not. `continues()` reports `p.depth > 0 || !p.tok.NLBefore` — at depth zero, a token carrying `NLBefore` ends the current statement and can't continue a postfix or binary chain.

`parseBlockStmt`, `parseSwitchStmt`, and `parseSelectStmt` all save and zero `p.depth` (and `p.noLit`) around their body for this reason; `parseConstraintDecl` and `parseDeclareDecl` do the same around their member lists, since those are also statement-like line-terminated bodies rather than comma-separated ones.

### Token flow

`advanceToken` moves to the next non-comment token, buffering at most one token of lookahead (`p.next`/`p.has`) for `peek()`. Comments are consumed by `consumeCommentGroup` and attached as `leadComment` or `lineComment` depending on whether a newline preceded or followed the group; a skipped comment's `NLBefore` is folded into the token that follows it, so `x = 1` / `/*c*/ y` on the next line doesn't lose the line break just because it landed on a dropped `COMMENT` token.

`peek()` exists only for the grammar's explicit `[lookahead]` restrictions — currently the launch-prefix disambiguation (`tryParseLaunch`) and the named-parameter check in `parseParam`. It is not a general parsing convenience.

### Diagnostics

Every error path goes through `errorAt`/`errorHere`, which increment `p.errCount` and, if a `Reporter` is set, call `diag.At`/`diag.AtToken`. `describe(t)` renders a token for message args (quoted literal text for idents/literals, `'kind'` for everything else, `"end of file"` for `EOF`). `expect` and `expectIdent` are the two workhorses: on success they consume and return a position/node; on failure they report and synthesize a placeholder (a zero-width `Pos`, or an `Ident{Name: "_"}`) so the caller can keep building a tree instead of branching on error.

### Recovery

`advance(to)` skips tokens until one in the given set is reached, bounded by a `syncPos`/`syncCnt` guard: if the scanner is stuck reporting the same position more than 10 times, `advance` gives up and returns anyway rather than looping forever. Two sync sets are shared across files — `stmtStart` (statement-leading keywords plus `RBRACE`) and `declStart` (declaration-leading keywords plus `IMPORT`/`EOF`) — matched to whichever grammar level the failing construct belongs to.

### File organization

- **`parser.go`** — the `parser` struct, token flow, diagnostic helpers, recovery, and the two entry points plus `parseFile`/`parseBuildClause`.
- **`decl.go`** — top-level declarations (`func`, `struct`/`class`, `enum`, `type`, `constraint`, `declare`, `let`/`var`) and imports.
- **`stmt.go`** — statements: blocks, `if`/`while`/`for`/`switch`/`select`, `return`/`defer`, assignment vs. expression-statement disambiguation.
- **`expr.go`** — expressions and types: precedence climbing, postfix chains, literals, and the type grammar, which shares enough surface with expressions (index brackets, parenthesized groups) that the two are parsed together.

A few grammar-level decisions worth knowing before extending any of these:

- **Ambiguities are resolved locally, then handed to the analyzer.** Where two productions collide — `var w` as a `VarDecl` vs. a bare `TransferExpr` (`parseVarDecl`), `mut a, b` in a `for` binding (`parseForStmt`), a single identifier as both a one-term `TypeSet` and a `ConstraintName` (`parseConstraintDecl`) — the parser picks the reading with an actual grammar production and leaves rejection of the other reading to a later phase, rather than trying to encode the static rule itself.
- **Grammatically-illegal-but-parseable forms are parsed anyway.** `declare` block bodies reject fields and bodied functions (A.8.3); a payload enum variant with an explicit discriminant is grammatical-but-static-rejected (A.6.5). Both parse fully so the diagnostic that eventually fires can point at the construct instead of the parser bailing at a syntax error.
- **`sizeof`/`alignof`/`reinterpret` are recognized by name**, not by keyword, in `parseArgumentList`, to take a `Type` as their first argument. This is sound only because A.1.4 forbids shadowing a `ReservedBuiltinName` — without that guarantee it would be unsound to special-case an identifier this way.
- **Type and expression parsing overlap by necessity**, not accident: `Stack[int32]` vs. `a[i]` (`parseExprOrType`), and a bracketed group that turns out to be an array *type*'s length vs. an array *literal*'s elements, decided only after the closing bracket is seen (`parseBracketedTypeOrArray`).