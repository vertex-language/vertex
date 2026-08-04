# parser

```go
import "github.com/vertex-language/vertex/parser"
```

Package `parser` turns Vertex source into an `ast.File` or an `ast.Package`. It depends on `ast` (its output), `diag` (its rejections), `scanner` (its input), and `token`.

## Design philosophy

Two grammar-wide mechanisms shape everything in this package.

### Statement termination by bracket depth, not by token

Whether a line terminator ends a statement is **not** a property of any token — it's a property of the innermost enclosing bracketing construct. `(`, `[`, and every brace that does *not* open a terminator-significant body push a depth counter; the ones that do — a block, a struct/class body, a constraint body, a switch/select body, a declare body, a foreign class body — reset depth to zero instead, and the source file itself starts at zero.

```go
func (p *parser) continues() bool { return p.depth > 0 || !p.tok.NLBefore }
```

At depth zero, a token carrying `NLBefore` ends the statement and cannot continue a postfix or binary chain. Two pairs of helpers manage this:

- `open`/`close` — increment/decrement depth around a bracket that suspends termination (parens, index brackets, argument lists).
- `enterTerminated`/`leave` — save the current `(depth, noLit)` state and reset depth to zero for a body that *is* terminator-significant, restoring it afterward.

Because a run of consecutive terminators collapses into a single `NLBefore` flag on the following token, `[ terminator ]` at the head of a list needs no code at all.

### The literal-in-header ambiguity

This is the one place the parser must resolve something the grammar leaves to prose: a `CompositeLit` or `MapLit` written unparenthesized between a control-flow keyword (`if`, `while`, `for`, `switch`) and its block brace reads as the block. `noLit` (plus `headerKw`, which names the construct for the diagnostic) suppresses the literal reading while a header is being parsed, and `withLit` locally re-enables it inside any bracketed group — parentheses are the grammar's prescribed escape hatch, and an index bracket or argument list encloses a literal just as effectively.

```go
func (p *parser) parseHeaderExpr(kw string) ast.Expr {
	savedLit, savedKw := p.noLit, p.headerKw
	p.noLit, p.headerKw = true, kw
	x := p.parseExpr()
	p.noLit, p.headerKw = savedLit, savedKw
	...
}
```

A composite literal is caught in `parseHeaderExpr` itself, by lookahead for the two tokens (`identifier`, `:`) that can only begin a field value; a map literal is caught earlier, in `parsePrimaryExpr`, where the brace is reached with an operand expected.

### Everything else parses and is rejected later

Every other construct the grammar forbids — a stacked ownership qualifier, `tensor` outside an npu body, `var` on a computed expression, a body on a foreign declaration, a repeated marker — is parsed into a real node rather than rejected at the token level. This is deliberate: it lets the diagnostic (raised later, by the analyzer) name the construct itself instead of pointing at a bare token.

## Package layout

| File | Contents |
|---|---|
| `parser.go` | The `parser` struct, token flow (`advanceToken`, `peekAt`), diagnostics helpers, termination (`continues`, `expectTerminator`, `enterTerminated`), recovery (`advance`, `stalled`), and the entry points `ParseFile`/`ParseDir` |
| `expr.go` | Expressions, operators, types (types are exprs, per `ast`), and signatures |
| `stmt.go` | Statements, `switch`/`select` clauses, and patterns |
| `decl.go` | Top-level declarations, `declare` blocks, and foreign members |

## Entry points

```go
func ParseFile(fset *token.FileSet, filename string, src []byte, rep diag.Reporter, mode Mode) (*ast.File, error)
func ParseDir(fset *token.FileSet, dir, importPath string, target token.BuildTag, rep diag.Reporter, mode Mode) (*ast.Package, error)
```

`ParseFile` always returns a non-nil `*ast.File`, even when errors were reported — recovery produces `Bad*` nodes so later phases and editor tooling can still walk a partial tree. If `rep` is `nil`, diagnostics are collected internally, sorted, deduplicated, and returned as the `error`.

`ParseDir` runs two passes over every `.vs` file in a directory:

1. **Probe pass** — each file is parsed with `PackageClauseOnly` against a throwaway `FileSet`, just far enough to read its build tag. A file whose tag doesn't match `target` is excluded from the build outright, so this filter must run before any file is fully parsed.
2. **Full pass** — the surviving files are fully parsed against the real `FileSet` and handed to `ast.NewPackage`, which checks that they agree on a package clause name.

Both passes are required rather than opportunistic: the qualifier under which an imported package's symbols are reached comes from that package's own package clause, so surviving files' names must agree before anything can resolve.

### Modes

```go
const (
	PackageClauseOnly Mode = 1 << iota // stop after package + build clauses
	ImportsOnly                        // stop after import declarations
	ParseComments                      // retain comments in the tree
)
```

`PackageClauseOnly` is load-bearing for `ParseDir`'s probe pass, not an optimization.

## Parsing structure

The parser is a straightforward recursive-descent/precedence-climbing hybrid:

- **Expressions** (`expr.go`) climb seven binary precedence levels (`parseBinaryExpr`), folding `as`-casts separately (`parseCastExpr`) since a cast's right operand is a `Type`, not an `Expr`, and a precedence loop can't consume it.
- **Types are parsed by `parseType`**, consistent with `ast`'s "types are Exprs" design. `parseExprOrType` and `parseBracketedTypeOrArray` resolve the genuinely ambiguous positions — an index bracket vs. a type-argument list, an array literal vs. an array/slice type — by shape where possible and by trailing context (`startsType`) where not.
- **Statements** (`stmt.go`) dispatch on leading keyword in `parseStmt`; anything else falls through to `parseSimpleStmt`, which parses an expression and only *then* decides, from what follows, whether it was actually an assignment target list.
- **Declarations** (`decl.go`) similarly dispatch by keyword in `parseTopLevelDecl`, with shared machinery for the constructs that reappear inside `declare` blocks (`parseForeignMember`, `parseForeignFunc`, `parseForeignClass`).

## Recovery

Unbounded body loops (import lists, field lists, argument lists, etc.) universally follow the same pattern:

```go
for !p.at(token.RPAREN) && !p.at(token.EOF) {
	before := p.tok.Pos
	// ... parse one element ...
	if p.stalled(before) {
		continue
	}
}
```

`stalled` checks whether the current position actually advanced since the loop iteration began; if not, it force-advances one token. This guarantees no such loop can spin forever on a malformed element that consumes nothing.

For coarser recovery, `advance(to map[token.Kind]bool)` skips tokens until one in a synchronization set is reached (`stmtStart`, `declStart`, `memberStart`, `clauseStart`), so one malformed construct doesn't cascade into spurious errors for everything that follows. A `syncPos`/`syncCnt` guard caps repeated resync attempts at the same position, the standard defense against a recovery loop that never actually consumes input.