# types

```
github.com/vertex-language/vertex/types
```

`types` defines the Vertex type representation, the compile-time constant representation, and the predicates over both that `analyzer` and `lower` share. It is the package `analyzer.Checker` fills in and `lower` reads back out.

It depends only on `ast` — and only for one reason (`info.go`'s side tables are keyed by `ast` nodes). Every other file is dependency-free. `ast` never imports `types`; the dependency runs one way, matching `analyzer`'s doc comment that the tree "records shape, never meaning."

Keeping the representation in its own package, separate from the checker that populates it, is what lets `lower` consume sizes, alignments, and ownership discipline without pulling in the checking machinery that produced them (A.15's invariant that "every value has a statically known layout, and every cost is decided at compile time").

## The `Type` interface

```go
type Type interface {
    Underlying() Type
    String() string
}
```

Every concrete type implements this by returning itself from `Underlying()`, except `*Named`, whose `Underlying()` returns what it was declared as. `Underlying` (in `predicates.go`) strips exactly one layer of naming, which is the operation A.7.3's `~T` type-set term is defined over.

`Constraint` deliberately does **not** implement `Type`. A.7.2 makes a constraint "never a value type, legal only in a `[...]` position," and A.14 lists `var c: Ordered` among the forms the grammar accepts but the checker must reject — making `Constraint` a separate hierarchy turns that rejection into a type error the checker cannot forget to raise, rather than a predicate someone has to remember to call.

## Kinds of `Type`

| Type | Spelling | File |
|---|---|---|
| `*Basic` | `int`, `bool`, `string`, … plus the untyped kinds | `basic.go` |
| `*Ownership` | `unique T`, `shared T`, `weak T` | `composite.go` |
| `*Array` | `[N]T` | `composite.go` |
| `*Slice` | `[]T` | `composite.go` |
| `*Map` | `map[K]V` | `composite.go` |
| `*Tuple` | `(T, T, ...)`, and every `Signature`'s result list | `composite.go` |
| `*Chan` | `chan T` | `composite.go` |
| `*Pointer` | `typed_ptr T` | `composite.go` |
| `*Signature` | a function/method/initializer's type | `composite.go` |
| `*Tensor` | `tensor[T, dims...]` | `composite.go` |
| `*Named` | a struct, class, enum, or `abstract` alias | `named.go` |
| `*Abstract` | the bare `abstract` alias target | `named.go` |
| `*TypeParam` | a generic type parameter | `named.go` |
| `*Struct` | the underlying of a struct/class | `named.go` |
| `*Enum` | the underlying of an enum | `named.go` |

`Unit` (`composite.go`) is `()`: an empty `*Tuple`, not its own kind, so a void return and a unit value are one object and nothing has to special-case one against the other.

### Basic types and the untyped kinds

`Typ` (`basic.go`) is a `[]*Basic` slice holding exactly one singleton per `BasicKind`, so identity comparison on `*Basic` **is** type identity. `byte` and `uint8` deliberately share `Typ[Uint8]` — A.1.4's rule that they "denote the same type; no conversion is required or permitted between them" falls straight out of the pointer comparison in `Identical` rather than needing a case of its own.

The untyped kinds (`UntypedBool` … `UntypedNil`) exist only between a literal and the typed position it eventually reaches (A.1.5.1): they never appear in a declared signature. `Default` (`basic.go`) is the fallback an untyped constant takes when nothing imposes a destination — `let x = 1` needs an answer, and it's `int`.

`LookupPredeclared` looks up A.1.4's `PredeclaredTypeName`s by spelling; they're ordinary identifiers pre-bound in `Universe`, not scanner keywords.

### Constant values

`Value` (`const.go`) is the compile-time constant representation — the reason it exists at all is A.1.5.1: an integer literal is untyped until it reaches a typed position, and "a literal that does not fit the destination type is a compile error, never a silent truncation." Holding every constant at unbounded precision (`*big.Int` / `*big.Rat`) until the destination is known is what makes that check possible instead of already-lost.

`MakeBool` / `MakeInt` / `MakeInt64` / `MakeFloat` / `MakeChar` / `MakeString` construct values; `Int64Val` / `BoolVal_` / `StringVal_` extract them. `Neg`, `BinaryOp`, `Shift`, and `Compare` fold constant operations — `BinaryOp` returns `Unknown` on constant division by zero rather than trapping, since a constant divide-by-zero is the caller's diagnostic, not a runtime trap (A.15).

`Representable` is the single enforcement point for the no-silent-truncation rule: it checks a `Value` against a `*Basic`'s range (`intRanges`, built in `init()` from `sizes.go`'s `intBits`), and a float constant with a nonzero fractional part is never representable in an integer type — there is no rounding rule in Vertex.

### Composite types

Each composite constructor (`NewOwnership`, `NewArray`, `NewSlice`, `NewMap`, `NewTuple`, `NewChan`, `NewPointer`, `NewSignature`, `NewTensor`) is a thin literal wrapper; the interesting content is in what each type's doc comment pins down as a spec consequence rather than an implementation choice:

- **`Tuple`** holds `*Var`s, not bare `Type`s, because A.3.1's `TupleElement` may be named. An unnamed element is a `Var` with an empty name.
- **`Chan`** carries no direction. A.4.2's `thread`/`async` launch prefixes both "evaluate to a receive-only `chan T`" — that's a property of the launch expression's result, not a second channel type, which is the entire point of A.10.1's sigil unification.
- **`Signature`** carries a `Marker` (`type.go`) because A.4.2 makes the marker part of the callee's contract, checked at both the definition and the launch site. Whether a signature is *capturing* is unknowable from the type alone (A.3.4: one word non-capturing, two words `{code, env}` capturing) — `Sizeof` is conservative and the analyzer narrows it per-expression.
- **`Tensor`** shape entries are resolved to `int64` at construction rather than kept as expressions, and is grammatical only under `[+Npu]` (A.3.5).

### Named types

`Named` (`named.go`) is what a struct, class, enum, or `abstract` alias declaration mints. A **transparent** `TypeAliasDeclaration` mints nothing — A.6.6 says it "names the same type and satisfies a `~T` type-set element," so the resolver substitutes the target and no `*Named` is ever constructed for it. Only `= abstract` is nominal.

- **`Struct`** is the shared underlying of both `StructDeclaration` and `ClassDeclaration` — A.6.3: "a class is byte-for-byte identical in layout to a struct and differs only in its member and method model." `Sizeof` never branches on `Class()`; construction syntax and `===` identity do.
- **`Enum`** resolves each `Variant`'s discriminant at construction (continuing from the previous variant when unwritten, per A.6.5) and exposes `UnitOnly`, which is what `Sizeof` and the `as`-cast rules key on: a unit-only enum *is* its discriminant integer.
- **`Abstract`** carries a `Family` (`FamilyMemoryFlat` / `FamilyObjectGraph`) because A.4.4's `abstract → typed_ptr T` cast is legal only for a memory-flat import family — an object-graph handle "has no byte representation to point at."
- **`TypeParam`** always resolves to a non-nil `Constraint` after resolution — A.7.1: "a bare name is constraint `any`."

## Objects (`object.go`)

`Object` is what `analyzer` resolves every `*ast.Ident` to: a variable, constant, function, type name, builtin, or package name.

```go
type Object interface {
    Name() string
    Type() Type
    Pos() token.Pos
    Pkg() *Package
    SetType(Type)
    objectNode()
}
```

`Type()` returning `nil` is not an error state — A.2's order-independence guarantee means every package-scope name is inserted before any type is resolved, so `nil` means "not yet resolved" and `Typ[Invalid]` means "resolution failed." `SetType` is exported because `analyzer` lives in a separate package and is the only caller.

- **`Var`** covers a variable, parameter, receiver, struct field, or tuple element. `Mode` (A.3.2's `mut`/`var`) is kept off the `Type` deliberately — see `type.go`'s `Mode` doc comment. `Addressable` answers whether it has a real storage slot: a `mut` binding or a field does, a `let` does not (A.5.1: "may be a register, an SSA value, or folded away entirely"). `Owning` reports only that a parameter's *position* is owning (`var T`); whether a given call passes the original or a copy is a call-site question (A.4.6/A.9.1), never a `Var` question.
- **`Const`** is a compile-time constant: an enum discriminant, an `ArrayLength`, or a compile-time-evaluable top-level initializer.
- **`Func`** covers functions, methods, initializers, deinitializers, foreign declarations, and constraint method requirements. `init`/`deinit` get no separate object kind — A.6.4 makes them `ContextualKeyword`s that are ordinary method names, so `IsInit`/`IsDeinit` are name-and-receiver questions answered on demand. `IsEntry` checks A.6.1's `main` shape.
- **`TypeName`** serves both a type name and a constraint name with one struct, because A.7.2's single-identifier-in-`[...]` ambiguity is resolved by what the name denotes, not by shape — a second object kind would turn that into a two-map lookup. `IsAlias` distinguishes a transparent alias (no `*Named` minted, or one minted for a different object) from a nominal one.
- **`Builtin`** covers A.1.4's `ReservedBuiltinName`s (`sizeof`, `new`, `transfer`, …). It carries `Typ[Invalid]` as its type because each builtin's shape is arity- and type-argument-specific rather than expressible as a `Signature` — the checker dispatches on `Id()` instead.
- **`PkgName`** is an imported package's qualifier, scoped to one file. Its `Name()` always comes from the imported package's own `PackageClause` (A.2.3), never from the import path — there is no aliasing, dot-import, or blank-import form to choose a different name from.
- **`Package`** is one checked compilation unit: name, path, `Scope`, and a `Complete` flag an `Importer` must check before handing the package out, since a half-checked scope would answer a lookup with a nil type.

## Scopes (`scope.go`)

`Scope` is a lexical block of named `Object`s, with `Pos`/`End` extent that a block-scope lookup consults and the package/universe scopes ignore — package scope is position-independent because A.2 makes top-level order irrelevant.

`Insert` drops a `BlankIdentifier` rather than inserting it (A.1.2: `_` "never introduces a usable binding") and returns any existing object of the same name for the caller to diagnose as a duplicate. `LookupParent` walks outward to `Universe`, A.1.4's implicit outermost scope.

`Universe` is built once in `init()`: every `PredeclaredTypeName` and `PredeclaredConstraintName` (`any`, `comparable`), plus every `ReservedBuiltinName`. Predeclared type/constraint names are shadowable; builtins are not (A.1.4) — `Insert` alone can't express that asymmetry, so `Reserved` exists as a separate check the resolver consults before inserting anything at all.

## Predicates (`predicates.go`)

- **`Identical`** is type identity. Named types compare by declaring object, not by shape — two structurally identical `*Named`s are still distinct (A.6.6's "two abstract aliases never unify, however identical their provenance"), and `byte`/`uint8` never need a special case since they already share one `*Basic`.
- **`IsComparable`** pins down what A.7.4 leaves circular ("every type supporting `==`/`!=`"): scalars and strings by value; arrays, tuples, structs, and payload enums elementwise/fieldwise; slices and maps **not** comparable (a header comparison answers a question nobody asked, and a contents comparison is an implicit O(n) walk the language never performs); `typed_ptr` by address; `unique`/`shared`/`weak`, `chan`, and `abstract` not comparable at all.
- **`AssignableTo`** reduces to identity plus one case: an untyped constant is assignable when representable at the destination (A.1.5.1). There is no implicit widening between distinct types — A.1.5.2's `'A'`/`"A"` split is exactly this rule. Ownership is deliberately not part of assignability: whether `var T` receives the original or a copy is decided by the `transfer` marker at the call site (A.4.6), not by the type.
- **`ConvertibleTo`** is A.4.4's `as`: always a static reinterpretation or width-selected instruction, never touching memory. Handles numeric↔numeric, char↔integer, unit-only-enum↔integer (a tag read), pointer↔pointer, and `abstract`→`typed_ptr` gated on `FamilyMemoryFlat` — with no direction convertible *to* `abstract` at all, since there's no runtime type information to reconstruct one from.

## Constraints (`constraint.go`)

`Constraint` holds a `TypeSet` (`[]Term`), method requirements, and embedded constraints — the intersection of all three (A.7.2). `Satisfies` checks a candidate type against all of them; `Any` (`terms`/`methods`/`embeds` all empty) admits everything and is what a bare type-parameter name means.

`Term.Tilde` is A.7.3's `~T` — admits `T` and anything whose *underlying* type is `T` — versus a bare `T`, which admits only `T` exactly, checked via `Identical` rather than `Underlying`+`Identical`.

Method requirements are checked structurally (`hasMethod`/`matchesRequirement`, ignoring the receiver, since the receiver is exactly what varies across satisfying types) — and because A.7.2 makes every instantiation monomorphized, satisfaction found here is never re-checked at runtime; there is no interface value and no vtable to build one into.

## `Info`: the side-table (`info.go`)

`Info` is the only file in this package that imports `ast`, and it's where every analysis result lands — never back onto the syntax tree itself:

```go
type Info struct {
    Types      map[ast.Expr]TypeAndValue
    Defs       map[*ast.Ident]Object
    Uses       map[*ast.Ident]Object
    Selections map[*ast.SelectorExpr]*Selection
    Instances  map[*ast.Ident]Instance
    Scopes     map[ast.Node]*Scope
    Transfers  map[*ast.TransferExpr]*Var
    Defers     map[*ast.DeferStmt]*Scope
}
```

A nil map means "do not record this" — each `Record*` method checks before storing, so a caller wanting only `Uses` pays for only `Uses`. `NewInfo` allocates every table; a caller wanting a subset builds the struct literal directly.

Keeping every result out-of-tree is what lets a printer round-trip a checked file unchanged and lets the checker run twice over one tree without residue — `ast`'s own doc comment is that the tree "records shape, never meaning."

- **`TypeAndValue`** pairs an `OperandMode` with a `Type` and, for a constant, a `Value`. `Mode` is what resolves A.3.6's `a[i]` vs. `Stack[int32]` ambiguity and similar: `IsType()`/`IsBuiltin()`/`IsVoid()`/`IsPackage()` read it directly, and `Addressable()` — true only for `VarMode` — is what A.4.8's `addr` and A.3.2's `mut` argument both gate on.
- **`Selection`** records what a `SelectorExpr` resolved to: `FieldVal`, `MethodVal`, `TupleIndex` (A.4.3's `t.0`), `PackageMember`, or `NamespaceMember` (one of A.4.1's four keyword namespaces — a closed set, per A.11.3, that never resolves to a user object).
- **`Instances`** records type arguments at each generic instantiation site. It's the set `lower` enumerates for monomorphization — A.7.5: "because every instantiation is monomorphized, the call in the generic body lowers to a direct call on the concrete type."
- **`Transfers`** is the one table that changes generated code rather than merely licensing it: A.4.6 makes the `transfer` marker's presence "the entire difference between move and deep copy," and A.9.4 prices a bare copy at O(data) versus a transfer at O(1) — `lower` reads this table to pick which to emit.
- **`Defers`** groups each `DeferStmt` by enclosing scope, which is what makes A.5.8's "reverse registration order on every exit edge" enumerable in `lower` — there is no unwinder, so the exit-edge set is finite and static.

`TypeOf` and `ObjectOf` are convenience accessors: `TypeOf` falls back to resolving through `ObjectOf` when an expression itself has no recorded entry (e.g. a bare identifier), and `ObjectOf` checks `Defs` before `Uses` so a declaring occurrence answers itself rather than falling through.

## Layout (`sizes.go`)

`Sizes` (`Sizes64`/`Sizes32`, selected by `SizesFor(tag)` off A.2.2's build tag) answers `Alignof`/`Sizeof`/`Offsetsof`. Only pointer width varies by target — `int`/`uint` are pinned to 64 bits everywhere via `intBits` (`type.go`... actually `sizes.go`), including `wasm32`, so that `sizeof(int)` and struct layout stay portable across a `build` split instead of becoming a per-target answer (A.15's static-layout invariant).

`Offsetsof` lays out fields in declaration order with ABI padding and no reordering (A.6.2) — a reader can predict layout straight from source order. `enumSize` implements A.6.5 directly: a unit-only enum is sized as its discriminant; a payload enum is a tag plus the largest variant, padded to the max of the two alignments.

## Rendering (`string.go`)

`TypeString` renders a `Type` back to source syntax and is normative: it's what a `diag` template's `%s` receives, so a change here changes what an `Expected(error, "...")` test matches (A.12.2). `ObjectString` renders a full object description (`"variable x of type int32"`, `"constraint Ordered"`, …) for diagnostics.

Both are single recursive writers (`writeType`/`writeTuple`/`writeSignature`) rather than per-type `String()` methods with independent logic — every concrete type's `String()` just calls `TypeString(t)`, so there is exactly one rendering rule to keep in sync with the grammar.