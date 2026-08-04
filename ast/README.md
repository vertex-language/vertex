# ast

```go
import "github.com/vertex-language/vertex/ast"
```

Package `ast` defines the Vertex syntax tree. It is the parser's sole output and the analyzer's sole input.

## Design philosophy

The tree records **shape, never meaning**. Several constructs in Vertex are ambiguous at the syntax level and are resolved only by what an operand denotes — something only the analyzer knows:

| Construct | Reading A | Reading B | Resolved by |
|---|---|---|---|
| `Index` node | `Index` expression | `TypeArgs` list | whether the operand denotes a generic declaration |
| `&x` | address-of | dereference | whether `x` is a `typed_ptr` |
| `~x` | bitwise-NOT | underlying-type | expression vs. `TypeSetTerm` position |
| single-identifier constraint element | one-term `TypeSet` | constraint name | what the name denotes |

The parser doesn't attempt to disambiguate these — it can't, and doesn't need to. Each is a static rule checked over an already-parsed tree, which is exactly what the analyzer does.

A corollary of this is that **types are `Expr`s**. There is no separate type-tree: a `TypeName` is an `*Ident`, one bracket node (`*IndexExpr`) serves both `Index` and `TypeArgs`, and a tuple type is a tuple literal (`*TupleExpr`). Conversion to a proper type representation happens downstream, in the analyzer.

## What the tree deliberately omits

Two things leave no trace in the tree:

- **Terminator significance.** Whether a line terminator ends a statement depends on the innermost enclosing bracketing construct, which the parser resolves as it goes. No node records this.
- **Trailing commas**, with one exception: `TupleExpr.TrailingComma`. A one-element tuple is distinguished from a parenthesized expression by nothing else, so it must survive. Everywhere else a trailing comma is optional and inert (`TypeArgs`, `Parameters`, `ArrayLit`, `LiteralValue`) and is not recorded.

Conversely, syntax that is *licensed by context* survives, because it is syntax: the `var` transfer marker (`*TransferExpr`) is the entire difference between a move and a deep copy, so it gets a node even though only some positions accept it.

## Package layout

| File | Contents |
|---|---|
| `ast.go` | Core interfaces (`Node`, `Expr`, `Stmt`, `Decl`), comments, identifiers, and shared parts (`Param`, `TypeParam`, `Marker`) shared across declarations and types |
| `expr.go` | Expressions, operators, and types (types are exprs) |
| `stmt.go` | Statements, including `switch`/`select` clauses and patterns |
| `decl.go` | Declarations: functions, records (struct/class), enums, aliases, constraints, vars, imports, and `declare` blocks for foreign interop |
| `file.go` | `File` and `Package` containers, plus `NewPackage` construction/validation |
| `walk.go` | `Visitor`, `Walk`, and `Inspect` for tree traversal |

## Core interfaces

```go
type Node interface {
	Pos() token.Pos // first character of the node
	End() token.Pos // one past the last character
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Decl interface {
	Node
	declNode()
}
```

Every node reports its own `Pos`/`End`, computed from its children rather than stored, except where a leading token (a keyword, a marker) makes storing a start position necessary.

## Notable node shapes

- **`RecordDecl`** is both `StructDecl` and `ClassDecl`. A class is byte-for-byte identical in layout to a struct and differs only in its member and method model; `Kw` (`STRUCT` or `CLASS`) carries the distinction.
- **`FuncDecl`** is `FunctionDecl`, `MethodDecl`, and the initializer/deinitializer forms. `init`/`deinit` are contextual keywords parsed as ordinary identifiers into `Name`; whether a given `FuncDecl` is one is a question for the analyzer, not the parser.
- **`ConstraintElem`** has exactly one of `Set` or `Method` non-nil. `Set` holds a `TypeSet` and a constraint name undifferentiated — again, resolved by what the name denotes.
- **`ForeignMember`** is implemented by `*ForeignFunc`, `*ForeignClass`, `*DeclareDecl`, and `*Field` — some of those only so a rejected construct (a nested `declare`, a field in a `declare` body) can parse and be diagnosed as itself rather than surfacing as a raw syntax error.
- **`BadExpr` / `BadStmt` / `BadDecl`** mark unparseable spans so recovery still yields a walkable tree; the analyzer skips them silently since a diagnostic was already reported at parse time.

## Traversal

```go
func Walk(v Visitor, node Node)
func Inspect(node Node, f func(Node) bool)
```

`Walk` visits exactly a node's declared children in depth-first order and panics on an unrecognized type — a deliberate tripwire meaning the switch in `walk.go` must be extended whenever a node type is added. `Inspect` wraps `Walk` with a plain `func(Node) bool`, called again with `nil` after each subtree finishes.

## Files and packages

`File` holds one parsed source file: an optional `Build` clause, its `Imports`, its `Decls`, and a flat `Comments` list (kept alongside `Doc`/trailing-comment attachments so a printer can recover anything the attachment heuristic missed).

`Package` is a validated container of `Files` — no I/O, no import resolution, no scopes. `NewPackage` checks only what makes the container internally coherent: at least one file, agreement on the package clause name, and agreement with the target build tag. Files are sorted by filename so the result is byte-reproducible. Everything else (import resolution, scoping) belongs to the loader and analyzer.