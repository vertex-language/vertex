# ast

```go
import "github.com/vertex-language/vertex/ast"
```

Package `ast` defines the syntax tree for Vertex. It records shape, never meaning: nodes carry spans, not strings, and text is recovered from the `token.File` that produced them. Nothing in this package resolves a name, decodes a literal, folds a constant, or knows what a type means.

`ast` imports `token` and nothing else. It doesn't know who built the tree — no import of `scanner` or `parser` is legal here.

## Node

Every node satisfies:

```go
type Node interface {
	Pos() token.Pos
	End() token.Pos
}
```

`Pos()`/`End()` are computed from a node's children, not stored redundantly, except where a leading keyword or contextual identifier makes storage necessary. `End()` is exact — not "roughly the closing brace" — because runtime traps depend on it. Every node has a non-zero span; a node with no children of its own (an unrecoverable parse) is a `Bad*` node that stores its own bounds.

Four hierarchies sit on top of `Node`: `Decl`, `Stmt`, `Expr`, and `TypeExpr`. Declaration nodes implement both `Decl` and `Stmt` — there is no synthesized `DeclStmt` wrapper, since the parser synthesizes nothing.

Types are a separate hierarchy from expressions, entered at type annotations, type arguments, heritage clauses, and the `< Type >` assertion. The parser always knows which side of that boundary it's on.

There is no `Ident.Obj` field and there never will be — identifier resolution isn't computable without types, so it doesn't belong on the node that exists before types are known.

## Walking a tree

```go
ast.Inspect(n, func(node ast.Node) bool {
	if id, ok := node.(*ast.Ident); ok {
		fmt.Println(id.NamePos)
	}
	return true // descend
})
```

`Walk`/`Inspect` traverse in source order, depth first. `Walk` **panics** on a node type it doesn't recognize. That's a tripwire, not a defect: it's what keeps the traversal switch from falling behind a new node type, and it fires just as usefully in your code if you hand `Walk` a node from an `ast` newer than the one you compiled against.

For write-your-own traversal, implement `Visitor`:

```go
type Visitor interface{ Visit(Node) Visitor }
```

`Inspect` is a thin adapter over `Visitor` for the common case of a single closure.

## Dumping a tree

```go
if err := ast.Fdump(os.Stdout, file); err != nil {
	log.Fatal(err)
}
```

`Fdump` writes a deterministic debug rendering — struct fields in declaration order, slices in source order, no map iteration anywhere. The format is **not** part of this package's contract and can change between commits. Golden tests that consume it are expected to be regenerated, not hand-edited, so a format change shows up as a diff across every fixture rather than as a silent behavior change elsewhere.

## Files and arenas

```go
tree, err := parser.ParseFile(fset, "main.vx", src, 0)
if err != nil {
	// ...
}
defer tree.Release()

// use tree...
```

A `*File` pairs with the `token.File` that produced it; positions are byte offsets in a per-unit address space and are meaningless without that pairing. The file's name and source text are deliberately not stored on the node.

`File.Release()` frees the arena backing the tree (via the `Releaser` interface, so this package doesn't import an allocator package). After `Release`, every node reachable from the tree is invalid — including any `token.Pos` a caller copied out earlier, which stay readable as plain integers and are exactly the trap: read the text you need *before* releasing. `Release` is idempotent, so pairing an explicit release with `defer tree.Release()` is safe, and because arenas are per-file, one file's tree can be released as soon as it's consumed — a whole-program build never has to hold every tree in memory at once.

## Optional children and nil

Optional children are typed pointers stored in interface-typed fields (e.g. an `Expr`-typed field holding a nil `*Ident`), so a plain `== nil` check on the interface doesn't work. Use the package's own nil-safe span helpers (`spanOf`, `endOf`) when writing new span methods, and don't compare an optional `Node`-typed field to `nil` directly.

## Modifiers

`Modifiers` records a modifier sequence twice: `Set` (a `ModifierSet` bitset) for membership queries, and `List` (`[]ModifierTok`, source order) for diagnostics that need to point at a specific word. The two must always agree — `List` is the source of truth and `Set` is derived from it during parsing; a test in this package checks the invariant. Use `Has` for "is this modifier present" and `Find` when you need the token to blame in an error.

## `StructDecl.Extends`

`StructDecl` carries an `Extends *HeritageClause` field even though the grammar has no `ClassExtendsClause` production for structs. This is deliberate: `struct S extends B {}` must still parse into a real `StructDecl` carrying its heritage clause, so that a later, name-based check can report *"structs don't support `extends`"* instead of the parser failing early with a bare "unexpected token `extends`". In any valid program this field is `nil`; a non-nil value only ever appears on a tree that a later pass is expected to reject.

## Stripping parentheses

```go
switch x := ast.Unparen(expr).(type) {
case *ast.CallExpr:
	// ...
}
```

`ParenExpr`/`ParenType` are retained in the tree rather than folded away — `(makeBox<boolean>)(true)` doesn't mean the same thing without the parens present. Most consumers don't care about them, so `Unparen`/`UnparenType` strip them on demand rather than every call site reimplementing the loop.

## Recovering literal text

`BasicLit` and friends store only spans, never a decoded value — `1_024` and `0b1010` need the target width to decode, which is a later phase's job. Get the raw spelling from the originating `token.File`:

```go
text := file.Between(lit.Pos(), lit.End())
```