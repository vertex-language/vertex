# parser

```go
import "github.com/vertex-language/vertex/parser"
```

Package `parser` builds a Vertex [`ast`](../ast) tree from source. It is recursive descent with precedence climbing over the binary chain (`vertex_grammar.md` B.4), and it always returns a non-nil tree — a partial parse is still a usable one (`compiler_frontend.md` §6). There is one goal symbol, no pre-scan, and no per-file mode that changes what gets parsed; `Mode` is a bitset of *options*, never a goal-symbol selector.

`parser` imports `ast`, `scanner`, and `token`. Nothing sits below it except an unexported `arena`, which lives inside this package on purpose — an allocator package would be a fifth package under a diagram that only has four, and `ast` depends on it only through the one-method `ast.Releaser` interface, so no import edge is added.

## Entry points

```go
tree, diags := parser.ParseFile(file, 0)
```

`ParseFile` scans and parses in one call. `ParseFileTokens` takes an already-scanned buffer, for callers that scanned for their own reasons and want the tree too — the parser buffers its input anyway, so handing it a pre-scanned slice costs nothing extra.

Two fragment entry points exist for REPLs, debuggers, and hover providers:

```go
x, diags := parser.ParseExpr(file)
t, diags := parser.ParseTypeExpr(file)
```

There are two rather than one because expressions and types are separate node hierarchies with no common return type to guess at — see `ast`'s note on `Decl`/`Stmt`/`Expr`/`TypeExpr`. Both fragment parses detach their arena instead of freeing it: the caller has no `ast.File` to call `Release` on, so the nodes are left for the GC rather than handed a lifetime they can't discharge.

`Mode` is a bitset:

- `ParseComments` retains comments on `ast.File.Comments`. The scanner always emits `COMMENT` tokens; this only controls whether the parser keeps them instead of dropping them at consumption.
- `ImportsOnly` stops after the import prologue — including re-exports, since `export ... from` is a dependency edge too — so a build driver can discover the dependency graph without paying for a full parse. It stops the *parser*, not the scanner: the file is still tokenized whole.
- `Trace` prints a production trace to stderr for debugging.

## Speculation

Several constructs are prefixes of something else: an arrow-function head is a prefix of a parenthesized expression, and `< Type >` competes with a generic arrow head. Rather than a cover grammar, ambiguous sites checkpoint and retry:

```go
ok := p.speculate(func() bool {
	// try a reading; return true to commit
})
```

A `mark` captures a token index, a diagnostic count, and an arena mark. Because the token buffer is immutable, resetting the diagnostic slice and truncating the arena back to the mark makes speculation side-effect-free by construction — no diagnostic escapes a failed attempt, and no node from a discarded reading leaks into the tree the caller walks. The bool `speculate` receives is a *commit* signal, not a success signal: instantiation expressions, for instance, commit only after the inner parse already succeeded, on a lookahead test that runs afterward.

Recursion depth is capped (`maxDepth`, currently 1000). Exceeding it is a diagnostic, never a hang — deeply nested but legal input exists (generated code full of parenthesization), so the cap is generous, and every unbounded recursive entry point (`parseAssign`, `parseBinary`, `parseType`, `parseStmt`) calls `enter`/`leave` around itself.

## Recovery

A parse never aborts. When a list-parsing loop fails to make progress, `advanced` force-advances by one token so nothing can spin, and `advanceTo` skips forward to the next token in a per-context `syncSet` (declaration starts, statement starts, member starts, or type-member starts). `skipBalanced` advances past a whole bracketed group in one step, rather than one token at a time, so recovery can't stop on a token that only *looks* like a synchronization point because it's nested inside `(...)`, `[...]`, or `{...}`.

Repeated resync at the same position is capped at one retry (`maxResync`); past that, `advanceTo` returns `false` and the caller must propagate the failure up to its own enclosing sync set rather than loop forever at a position recovery can't get past.

Member and statement lists that skip a broken slot still record a `Bad*` node for it, so downstream offsets don't silently shift.

## Semicolon insertion

`expectSemi` implements automatic semicolon insertion for the `TerminatedByASI` nonterminals in grammar section L: an explicit `;`, a line break before the next token, a following `}`, or EOF all terminate; anything else is a diagnostic. It returns the position of a written `;`, or `NoPos` when the terminator was inserted, so a node's span stops at its last real token instead of covering a semicolon that was never there.

`expectMemberSep` is the same idea for the seven `TypeMember` signature forms, which additionally accept `,` as a separator. Newline-separated interface bodies — `interface U { id: string\n name: string }` — are entirely this function's job; property signatures are deliberately included in the ASI-terminated set precisely so that case doesn't regress while method signatures keep working.

## Arenas

Nodes for one file are allocated together and freed together. `File.Release()` (via the single-method `ast.Releaser` interface) truncates or drops the arena's slabs in one step, so a whole-program build never has to hold every file's tree in memory just because one file was parsed early.

The arena hands out `*ast.Ident` from a fixed-capacity slab and falls back to the heap once the slab fills. Growing the slab with `append` would be wrong — `append` can reallocate, which would silently invalidate every pointer already handed out to the tree the caller is building — so the slab has fixed capacity and overflow is deliberately routed to the GC instead.

`speculate`'s rollback truncates the arena back to its mark as well as the token cursor: pointers into the truncated region belong to nodes the parser just discarded mid-speculation, so nothing dangles and nothing leaks.

## Contextual keywords and lookahead

Words like `struct`, `kernel`, `type`, `interface`, `namespace`, `declare`, `abstract`, and `using` are ordinary identifiers that only sometimes introduce a declaration. `ExpressionStatement`'s grammar lookahead restriction names just `{`, `function`, `async function`, `class`, and `let [` — it says nothing about these words, so the decision of when they start a declaration versus an expression is this package's own policy (`tryContextualDecl`), not a rule the grammar states directly.

The general shape of that policy is a `[no LineTerminator here]` check before a following token that could plausibly complete a declaration head: a line break forces the identifier reading, which is what keeps `struct` and `kernel` usable as ordinary variable names.

`>` is a related case at the token level rather than the identifier level: the scanner never merges adjacent `>` characters, so `Array<Box<int32>>` tokenizes as two `GT`s and the *expression* parser is what joins a run of them back into `>>`, `>>>`, or one of their `=`-suffixed forms when it's looking for a binary operator. Adjacency — no whitespace between the tokens — is the test, which is why `a > > b` never joins into a shift.

## Grammar forms that parse before they're rejected

A few constructs the grammar rules out are still parsed into real nodes rather than failing early, so that a later, name-based pass can give a specific diagnostic instead of a bare "unexpected token." `struct S extends B {}` is the canonical example: structs have no `ClassExtendsClause` production, but `parseStruct` still builds a `StructDecl` carrying a heritage clause and reports *"structs don't support `extends`"* at the clause's own span. `kernel async function` is the same idea one level up — accepted by the parser, rejected by name once `Accel` and `Async` are both set on the resulting `FuncDecl`.