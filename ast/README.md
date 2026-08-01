# ast

`ast` defines the Vertex syntax tree. It is the parser's sole output and the
analyzer's sole input.

The tree records shape, never meaning. Several Vertex constructs are resolved
by what an operand denotes rather than by syntactic form — `a[i]` versus
`Stack[int32]` (A.3.6), `&x` as address-of versus dereference (A.4.4), a lone
identifier in a constraint body (A.7.2). The parser cannot know which, and
A.0.2 says it does not have to: "Nothing else in Vertex is context-sensitive.
Every other rejection is a static rule checked over an already-parsed tree."
This package is that tree.

The corollary is that **types are `Expr`s**. `TypeName` is an `*Ident`,
`[...]` is index/slice/instantiate/type-argument all at once, and a tuple
type is a tuple literal. Conversion to a type representation happens in the
analyzer, not here.

Context parameters (`Await`, `Npu`, `Own`, `Lit`) leave no trace in the tree —
they are parser state. Source syntax *licensed* by them does survive: A.4.6's
`var` marker is the entire difference between move and deep copy, so it gets
its own node (`*TransferExpr`).

`ast` depends on `token` for `Pos`, `NoPos`, and `Kind`; it has no other
dependencies.

## Core interfaces

```go
type Node interface {
    Pos() token.Pos // first character of the node
    End() token.Pos // one past the last character
}

type Expr interface { Node; exprNode() }
type Stmt interface { Node; stmtNode() }
type Decl interface { Node; declNode() }
```

Every node in the tree satisfies `Node`. `Expr`, `Stmt`, and `Decl` are
otherwise-empty marker interfaces that partition nodes into the three
syntactic categories the grammar distinguishes; the unexported `exprNode()` /
`stmtNode()` / `declNode()` methods exist only to make that partition
exhaustive and closed to outside packages.

## Comments

```go
type Comment struct {
    Slash token.Pos
    Text  string // including the opening delimiter, excluding any trailing newline
}

type CommentGroup struct {
    List []*Comment // len(List) > 0
}
```

A `CommentGroup` is a run of comments with no other tokens and no blank lines
between them. `(*Comment) IsLineTerminator() bool` reports whether the
comment spans a `LineTerminator`, which per A.1.1 makes it one for
statement-termination purposes (A.0.6) — relevant because `//`-comments and
multi-line `/* */` comments interact differently with Vertex's
newline-sensitive statement grammar.

## Identifiers

```go
type Ident struct {
    NamePos token.Pos
    Name    string
}
```

`Ident` covers plain identifiers, the `BlankIdentifier` (`_`), every
`PredeclaredTypeName`, every `ReservedBuiltinName`, and every
`ContextualKeyword` used as a name (`init`, `deinit`, `test`, ...). None of
these are distinguished lexically (A.1.3, A.1.4); the analyzer resolves them
against the implicit outermost scope.

- `(*Ident) IsBlank() bool` — reports whether the identifier is `_`. A blank
  identifier never introduces a usable binding.
- `(*Ident) String() string` — returns the name, or `"<nil>"` for a nil
  receiver.

## Shared syntax parts

```go
type Param struct {
    Doc      *CommentGroup
    Name     *Ident    // nil inside a bare FunctionType
    Colon    token.Pos // NoPos when Name is nil
    Ellipsis token.Pos // position of `...`, else NoPos
    Type     Expr
}

type ParamList struct {
    Lparen token.Pos
    List   []*Param
    Rparen token.Pos
}
```

`Name` is nil inside a `FunctionType` (A.3.4: a function type names parameter
*types* only). `Ellipsis` marks a variadic parameter; A.6.1 requires it last
and permits at most one, both enforced statically rather than by the shape of
`ParamList` itself.

```go
type TypeParam struct {
    Name       *Ident // may be a BlankIdentifier
    Colon      token.Pos
    Constraint Expr // nil
}

type TypeParamList struct {
    Lbrack token.Pos
    List   []*TypeParam
    Rbrack token.Pos
}
```

`Constraint` is nil both for a bare name (constraint `any`) *and* for a name
in a group whose constraint appears on a later entry — A.7.1's
`[A, B: Number]` constrains both `A` and `B`. The parser does not distribute
the trailing constraint backward across the group; the analyzer does,
because distributing here would erase the written form a formatter needs to
reproduce.

```go
type Marker struct {
    MarkerPos token.Pos
    Kind      token.Kind // ASYNC, GPU, NPU, or IDENT for `test`
    Name      string     // "async", "gpu", "npu", "test"
}
```

A `FunctionMarker` (A.6.1). `test` is a `ContextualKeyword` and scans as
`IDENT`, so `Kind` alone can't carry it — `Name` is authoritative for
distinguishing which marker this is.

## Expressions

### Primaries

| Type | Notes |
|---|---|
| `BasicLit` | A `NumericLiteral`, `StringLiteral`, `CharLiteral`, or `ReservedLiteralKeyword` (`true`/`false`/`nil`). `Value` is the raw source spelling, never unescaped. |
| `NamespaceExpr` | One of the four keyword namespaces of A.4.1: `async`, `gpu`, `npu`, `chan`. Appears only as the `X` of a `SelectorExpr`, or (for `chan`) as the head of a construction. |
| `ParenExpr` | `( X )`. |
| `TupleExpr` | `TupleType`, `TupleLiteral`, *and* `UnitType` (A.3.1, A.4.7) — one node, because the three share a shape once types are `Expr`s. A named element is a `*KeyValueExpr`. Empty `Elems` with no trailing comma is the unit type `()`. A single element requires `TrailingComma` (A.4.7); `(x)` without one is a `*ParenExpr` instead. |
| `ArrayLit` | `[ ElementList ]` (A.4.7). Elements are owning positions, so any may be a `*TransferExpr`. |
| `CompositeLit` | `TypeName LiteralBody` or `InstantiatedType LiteralBody` (A.4.7). `Type` is never nil — a bare `{...}` is a `MapLit` instead. Elems are `*KeyValueExpr`; A.4.7 requires the key to be an identifier. |
| `MapLit` | A braced literal with no type prefix. Keys are arbitrary expressions, unlike `CompositeLit`'s field names. |
| `KeyValueExpr` | Covers every `X : Y` pair: a composite-literal field, a map entry, a named tuple element, and a named call argument (A.4.3, A.4.7). |
| `EnumShorthand` | `.Name` or `.Name(args)` (A.4.1). Legal only where the enum type is inferable from context — a static rule, not a parse-time one. |
| `FuncLit` | A `FunctionExpression` (A.4.1). Per A.0.3, it begins with all four context parameters cleared and re-sets them from its own `Marker`, so an anonymous closure inside an `async` body may not `await` unless it is itself marked. |

### Postfix

| Type | Notes |
|---|---|
| `SelectorExpr` | `X . Sel`. |
| `TupleIndexExpr` | Positional tuple access, `t.0` (A.4.3). Chains compose (`t.0.0`). Because maximal munch scans `t.0.0` as `IDENT PERIOD FLOAT("0.0")`, the parser must split a `<digits>.<digits>` float immediately after a selector dot; `Index` holds the decoded value, `Text` the source spelling of this component only. |
| `IndexExpr` | `x[a]`, `x[a..b]`, and `Stack[int32]` — index, slice, generic instantiation, and type-argument list, all one node (A.3.6), Vertex's general compile-time-configuration slot. Also carries `new[T]`, `npu.Quantize[T]`, and `chan[T]`. A slice is an `IndexExpr` whose single index is a `*BinaryExpr` with `Op DOTDOT`. |
| `CallExpr` | `f(a, b)` and `f[T](a)` (the latter with `Fun` a `*IndexExpr`). Builtins get no node of their own — A.4.8 treats them as ordinary calls over reserved, unshadowable names (A.1.4), so `sizeof(Type)` works because types are `Expr`s. `ExpectedType` (A.12.2) is likewise a plain `CallExpr`. |
| `LaunchExpr` | A launch prefix applied to a call (A.4.2): `thread`, `async`, `gpu`, `npu`. Per A.4.2 these modify scheduling only, never the callee's signature. `Config` is legal only on `gpu`. The `[lookahead != .]` restriction separating `npu Dot(a,b)` from `npu.Dot(a,b)` is a parser decision that leaves no trace: a namespace access simply produces `SelectorExpr{X: NamespaceExpr}`. |
| `LaunchConfig` | `(blocks: E, threads: E)` (A.4.2) — fixed arity and fixed names, so it is not modeled as a general argument list. |
| `AwaitExpr` | `await X`. |

### Operators

| Type | Notes |
|---|---|
| `UnaryExpr` | Covers `-`, `!`, `~`, and `&`. Two are deliberately left unresolved at this stage: `&` is address-of on a value and dereference on a `typed_ptr`, keyed on the operand's statically written type (A.4.4); `~` is bitwise-NOT in an expression and underlying-type in a type-set element (A.7.3). Both are one node; the analyzer decides which. |
| `BinaryExpr` | Collapses A.4.5's precedence cascade into one shape. Precedence comes from `token.Kind.Prec()` and A.13 — the cascade nonterminals are grammar-writing devices, not distinct tree shapes. `DOTDOT` is a `BinaryExpr` too; A.4.5 makes it non-associative, which the parser rejects outright rather than folding (`a..b..c` is a compile error). |
| `CastExpr` | `x as T` (A.4.4). Left-associative, binds tighter than every binary operator, and never touches memory. |
| `TransferExpr` | The ownership marker: `var target` in an owning position (A.4.6, A.9.1). It is a node rather than a bool because A.9.1 lists six distinct owning positions, and one node covers all of them without six flags. Its presence is the entire difference between move and deep copy, so it must survive parsing as syntax and never be normalized away. `Target` is grammatically an identifier or selector chain; anything else parses into this node and is rejected statically — A.14 lists "`var` on a computed expression" as a rejected form, meaning it must parse. |

### Types (as Exprs)

| Type | Notes |
|---|---|
| `OwnershipType` | `mut T`, `var T`, `unique T`, `shared T`, `weak T` (A.3.2). Qualifiers do not stack, but a stacked form parses and is rejected statically (A.14). |
| `ArrayType` | `[N]T` and `[]T` (A.3.1). `Len == nil` is the slice form. |
| `MapType` | `map[Key]Value`. |
| `FuncType` | A `FunctionType` (A.3.4), and also the signature half of a `FuncDecl`/`FuncLit`. In a bare `FunctionType`, every `Param.Name` is nil. `Result` carries `-> Type` or `-> ExpectedType` (the latter a `*CallExpr`); omitting it is the void form — A.3.4 has no `void` type name. |
| `ChanType` | `chan Elem`. |
| `PointerType` | `typed_ptr T` (A.3.3). Nesting requires parentheses, so the nested form has `Elem` be a `*ParenExpr`. |
| `TensorType` | `tensor[T, dims...]` (A.3.5). Grammatical only under `[+Npu]`; a tensor outside an `npu` body parses and is rejected (A.14). `Shape` holds `IntegerLiteral`s. |
| `AbstractType` | The bare `abstract` of A.3.3, legal only as an alias target. |

### Error recovery

`BadExpr` marks a span the parser could not make sense of, so recovery
produces a tree the analyzer can still walk. The analyzer skips it silently,
since a diagnostic was already reported at parse time. `BadStmt` and
`BadDecl` (below) serve the same purpose in statement and declaration
position.

## Statements

| Type | Notes |
|---|---|
| `BlockStmt` | `{ StatementList }` (A.5). Per A.0.6 this brace is newline-significant, unlike the braces of a `CompositeLit`, `MapLit`, `FieldList`, or `DeclareBody`; the scanner can't tell them apart lexically, which is why it records line breaks as `Token.NLBefore` rather than suppressing them. |
| `DeclStmt` | Wraps a declaration in statement position — in practice a `*VarDecl`, since A.5 lists `VariableDeclaration` among the statements and A.2 lists it among top-level declarations. |
| `ExprStmt` | A.5.9's `ExpressionStatement`. The grammar excludes a bare `CompositeLiteral`/`MapLiteral` here, which is what keeps `{...}` unambiguous against a `BlockStmt`. |
| `AssignStmt` | Both forms of A.5.2. `Op` is `ASSIGN` for the list form and a compound-assignment kind otherwise, in which case there is exactly one target and one value. A target may be `&x` (dereference-on-write, an `*UnaryExpr`) or the blank identifier. Assignment is always a statement — there is no `=` inside a condition anywhere in this tree. |
| `IfStmt` | Has no initializer clause; A.5.4 is explicit that the two-statement error-checking idiom is intentional. `Cond` is parsed under `[~Lit]`. `Else` is a `*BlockStmt` or `*IfStmt`. |
| `WhileStmt` | The only loop primitive (A.5.5). |
| `ForStmt` | `for IterationBinding in Expr Block` (A.5.6). `Mode` is `ILLEGAL` for the bare (shared-access) form, `MUT` for exclusive access, or `VAR` for the consuming form — carried on the binding rather than the iterable, since what transfers is each element, one per iteration. `Names` holds one name, or two for the index/value and key/value forms. |
| `SwitchStmt` / `CaseClause` | `Patterns == nil` marks the default clause; A.5.7 permits at most one. |
| `EnumPattern` | `.Name` or `.Name(bindings)` in case position (A.5.7). Distinct from `EnumShorthand` because the payload entries are binding names, not expressions — A.5.7 makes them views into the payload rather than copies. |
| `ReturnStmt` | A bare comma list, never parenthesized (A.5.3). |
| `DeferStmt` | Takes a call and nothing else (A.5.8); arguments are evaluated at registration, only the call is postponed. |
| `BranchStmt` | `break`, `continue`, `fallthrough`. There are no loop labels (A.5.9), so it carries none. |
| `SelectStmt` / `SelectClause` | One case of A.10.2. `Targets` is non-nil for the `x = ch.receive()` form. `Op` is the channel operation, optionally wrapped in an `*AwaitExpr` — A.10.2 requires a single `select` to be entirely bare or entirely awaited, checked statically. `Op == nil` marks the default clause. |
| `BadStmt` | Error recovery; see above. |

## Declarations

```go
type FuncDecl struct {
    Doc        *CommentGroup
    Recv       *Receiver // nil for a free function
    Name       *Ident
    TypeParams *TypeParamList // nil; A.7.6 forbids one on a method
    Type       *FuncType
    Body       *BlockStmt // nil only in error recovery
}
```

`FuncDecl` is also an `InitializerDeclaration` and `DeinitializerDeclaration`
(A.6.1, A.6.4). Those get no node of their own — `init` and `deinit` are
`ContextualKeyword`s that are ordinary method names in a receiver
declaration, so they scan as `IDENT` and parse as `FuncDecl`. Whether a given
`FuncDecl` is an initializer is a question about its `Name` and `Recv`,
answered by the analyzer.

```go
type Receiver struct {
    Lparen token.Pos
    Name   *Ident
    Colon  token.Pos
    Type   Expr
    Rparen token.Pos
}
```

`( Identifier : ReceiverType )` (A.6.1). `Type` may be wrapped in an
`*OwnershipType` for the `mut`/`var`/`shared` forms, and carries a
`TypeParameterList` as an `*IndexExpr` for a method on a generic type — A.7.6
makes that list bind the type's existing names rather than introduce fresh
ones.

```go
type RecordDecl struct {
    Doc        *CommentGroup
    KwPos      token.Pos
    Kw         token.Kind // STRUCT or CLASS
    Name       *Ident
    TypeParams *TypeParamList
    Lbrace     token.Pos
    Fields     []*Field
    Rbrace     token.Pos
}
```

Both `StructDeclaration` and `ClassDeclaration` (A.6.2, A.6.3) — one node,
because A.6.3 says a class is byte-for-byte identical in layout to a struct
and differs only in its member and method model. `Kw` carries the
distinction; every consumer that cares reads it. `Field.Default` is
evaluated at construction for any omitted field (A.6.2).

```go
type EnumDecl struct {
    Doc        *CommentGroup
    Enum       token.Pos
    Name       *Ident
    TypeParams *TypeParamList
    Colon      token.Pos
    Discrim    Expr // DiscriminantType; nil if absent
    Lbrace     token.Pos
    Variants   []*Variant
    Rbrace     token.Pos
}
```

`Variant` is one entry of the enum body (A.6.5): a unit variant, a payload
variant `Name(T, U)`, or a unit variant with an explicit discriminant. An
explicit discriminant on a *payload* variant parses and is rejected (A.14).

```go
type TypeAliasDecl struct {
    Doc        *CommentGroup
    Type       token.Pos
    Name       *Ident
    TypeParams *TypeParamList
    Assign     token.Pos
    Target     Expr
}
```

`type Name[params] = Target` (A.6.6). A `Target` of `*AbstractType` makes the
alias nominal and opaque; anything else makes it transparent.

```go
type ConstraintDecl struct {
    Doc        *CommentGroup
    Constraint token.Pos
    Name       *Ident
    Lbrace     token.Pos
    Elems      []*ConstraintElem
    Rbrace     token.Pos
}

type ConstraintElem struct {
    Set    Expr // exactly one of Set/Method is non-nil
    Method *MethodReq
}
```

`Set` holds a `TypeSet` or a `ConstraintName` undifferentiated, because A.7.2
says a single identifier parses as both and is resolved by what the name
denotes. A union is a `*BinaryExpr` with `Op OR`; a `~T` term is a
`*UnaryExpr` with `TILDE`. A `MethodReq` (`MethodRequirement`, A.7.2) is
satisfied by any type declaring a matching receiver method; monomorphization
lowers the call directly, so it introduces no interface value and no vtable.

```go
type VarDecl struct {
    Doc      *CommentGroup
    KwPos    token.Pos
    Kw       token.Kind // LET or VAR
    Bindings []*Binding
    Assign   token.Pos // NoPos for bare `var x: T`
    Values   []Expr    // owning positions: any may be a *TransferExpr
    Comment  *CommentGroup
}
```

`let`/`var` with initializers, and bare `var Binding`, from A.5.1. Also a
`TopLevelDeclaration`, where A.2 requires a compile-time-evaluable
initializer. `Binding.Type` is nil when inferred, required for bare
`var x: T`.

```go
type ImportDecl struct {
    Doc    *CommentGroup
    Import token.Pos
    Lparen token.Pos // NoPos for the single-path form
    Paths  []*BasicLit
    Rparen token.Pos
}
```

`import "path"` or `import ( ... )` (A.2.3). There is no aliasing form, no
dot-import, and no blank import, so there is nothing to record but
paths — the qualifier importers use comes from the imported package's own
`PackageClause`.

### Declare blocks

```go
type DeclareDecl struct {
    Doc     *CommentGroup
    Declare token.Pos
    KindPos token.Pos
    Kind    string // "framework" or "module"
    Variant *VariantTag
    Path    *BasicLit
    Lbrace  token.Pos
    Members []ForeignMember
    Rbrace  token.Pos
}
```

`declare framework "S" { }` or `declare module ["tag"] "S" { }` (A.8.1).
`Kind` is the `ContextualKeyword`, scanning as `IDENT`. A variant tag on a
`framework` block parses and is rejected (A.8.2). `VariantTag` reuses the
same compile-time-configuration bracket as generic instantiation and array
length.

`ForeignMember` is implemented by:

- **`ForeignFunc`** — a `ForeignFunctionDeclaration`/`ForeignInitializerDeclaration`
  (A.8.3). `Init` marks the `init` prefix modifier, which A.8.3 is explicit
  is a modifier on `func`, not a function name. `Name` is nil for the
  unnamed initializer form that bare `Type(...)` construction resolves to.
  `Modifiers` and `Body` exist only for the error forms — A.0.5 makes
  rejected forms part of the grammar, so a banned visibility modifier or a
  foreign declaration with a body must parse in order to be diagnosed as
  themselves rather than as a syntax error.
- **`ForeignClass`** — a `ForeignClassDeclaration` (A.8.3).
- **`ForeignField`** — exists only to be rejected (A.8.3 ✗ fields describe
  layout, which a `declare` block must not).

### Error recovery

`BadDecl` marks an unparseable declaration span, same role as `BadExpr` /
`BadStmt`.

## Files and packages

```go
type File struct {
    Doc     *CommentGroup
    Package token.Pos
    Name    *Ident // the PackageClause name — the qualifier importers use

    Build   *BuildClause // nil if absent
    Imports []*ImportDecl
    Decls   []Decl

    Comments []*CommentGroup

    FileStart token.Pos
    FileEnd   token.Pos
}
```

One parsed `.vs` source file. `Comments` holds every comment in the file in
source order, including any already attached as a `Doc` or trailing
`Comment` elsewhere — the same arrangement `gofmt` relies on: the flat list
is what a printer walks to place anything the attachment heuristic missed.

- `(*File) BuildTag() token.BuildTag` — the file's target tag, or
  `token.TagNone` if it carries no `build` clause.
- `(*File) Filename(fset *token.FileSet) string` — resolves the file's name
  via the `FileSet`.

```go
type BuildClause struct {
    Build  token.Pos
    TagPos token.Pos
    Name   string
    Tag    token.BuildTag
}
```

`build <tag>` (A.2.2). `Tag` is `token.TagNone` when `Name` is unrecognized.
A.2.2 makes an unrecognized tag a compile error rather than a silently
excluded file, so callers must distinguish "unknown tag" from "no build
clause" — `Build.IsValid()` answers the latter question.

```go
type Package struct {
    Name  string // from the PackageClause; all files agree
    Path  string // resolved import path — a locator, not a name (A.2.3)
    Dir   string
    Tag   token.BuildTag
    Files []*File // sorted by filename
}

func NewPackage(fset *token.FileSet, path, dir string, target token.BuildTag, files []*File) (*Package, error)
```

A `Package` groups files from one directory into a compilation unit — one
package is one `.o`/`.vbyte`, so its contents fix the compilation unit
exactly. It is a validated container and nothing more: no I/O, no import
resolution, no scopes. Resolution belongs to the loader, which is why
`NewPackage` does not take an importer.

`NewPackage` is pure — no filesystem access, no diagnostics beyond its own
well-formedness checks:

1. At least one file.
2. Every file agrees on the `PackageClause` name.
3. Every file's build tag (if any) matches `target`; filtering mismatched
   files out of the build is the *loader's* job, so `NewPackage` treats a
   mismatch reaching it as a caller bug, returning an error rather than
   silently dropping the file.

Files are sorted by filename so the resulting `*Package` is byte-reproducible.

## Tree traversal

```go
type Visitor interface {
    // Visit is called for each node. Returning a non-nil visitor walks the
    // node's children with it, then calls Visit(nil) on it to signal the end
    // of that subtree.
    Visit(node Node) Visitor
}

func Walk(v Visitor, node Node)
func Inspect(node Node, f func(Node) bool)
```

`Walk` performs a depth-first traversal, dispatching on the dynamic type of
`node` to visit exactly its declared children (e.g. `*IfStmt` visits `Cond`,
`Body`, and `Else` only if non-nil). Every exported node type is handled;
`Walk` panics on an unrecognized type, which in practice means the switch
must be extended whenever a node type is added.

`Inspect` adapts a plain `func(Node) bool` into a `Visitor` via the
unexported `inspector` type, for callers who don't need the `Visit(nil)`
end-of-subtree signal.

## Design notes

- **The tree is untyped on purpose.** Because types are `Expr`s, there is no
  separate `TypeExpr` hierarchy — `CompositeLit.Type`, `Field.Type`,
  `Param.Type`, and every type-position field are all just `Expr`. This is
  what lets `sizeof((int32, bool))` reach the ordinary expression parser
  instead of a separate type parser.
- **Rejected forms still have to parse.** A.0.5 and A.14 mean this package
  models several constructs that are always compile errors — a stacked
  `OwnershipType`, a `TensorType` outside `npu`, a `ForeignFunc` with a body,
  a discriminant on a payload `Variant` — because the analyzer needs a real
  node to attach each diagnostic to, not a parse failure.
- **One node per shape, not per grammar production.** `RecordDecl` unifies
  `struct`/`class`; `TupleExpr` unifies tuple type, tuple literal, and unit
  type; `FuncDecl` unifies function/initializer/deinitializer declarations.
  In each case the annex's grammar has multiple named productions but only
  one tree shape, and the tree follows the shape.
- **Context parameters and static rules leave no trace.** Neither the four
  parser context parameters (`Await`, `Npu`, `Own`, `Lit`) nor any rule A.0.2
  calls "static" appears as a field anywhere in this package — only their
  *syntactic residue* (e.g. `*TransferExpr` for `var`) does.