# types

```go
import "github.com/vertex-language/vertex/types"
```

`types` defines the Vertex type representation, the compile-time constant
representation, and the predicates over both that `analyzer` and `lower`
share. It is the package `analyzer.Checker` fills in and `lower` reads back
out — deliberately not the analyzer itself, so that `lower` can consume
sizes, alignments, and ownership discipline without pulling in the checking
machinery that produced them.

It depends on `ast` (`info.go`'s side tables are keyed by `ast` nodes) and on
`token` (positions throughout, and the one home of the reserved-builtin set,
consulted from `scope.go`). Nothing else. The dependency with `ast` runs one
way — `ast` never imports `types` — matching `ast`'s own doc comment that the
tree "records shape, never meaning."

## Design philosophy

**Citation convention.** A bare `§` cites `semantics.md`; CamelCase names
(`ArrayLength`, `ConstraintElem`, `ExpectedType`, ...) are `grammar.md`
productions. Where both documents are silent — enum layout, closure
representation, the vector ABI, narrow-float representability — the comment
says so explicitly rather than dressing an implementation choice as a rule.
One file hasn't caught up: `info.go` still cites the older `A.x.y` numbering
(`A.1.2`, `A.3.6`, `A.4.1`, `A.4.6`, `A.5.8`, `A.7.2`, `A.7.5`, `A.9.4`)
rather than the bare-`§` scheme the rest of the package uses — worth
converting the next time that file changes.

**Two things are deliberately not `Type`.**

- `Constraint` (`constraint.go`). §9 ⊢ a `ConstraintDecl` "takes no type
  parameters, is legal only in a `[`…`]` position, and is never a value
  type." Keeping it out of the `Type` hierarchy turns `var c: Ordered` into a
  shape error the checker cannot forget to raise, instead of a predicate
  someone has to remember to call.
- `Expected` (`composite.go`). §3.2 ⊢ `Expected(…)` is legal "only as a
  FunctionDecl/MethodDecl result, in a `build test` file." It hangs off
  `Signature` as its own field rather than sitting in the result `Tuple`, so
  it can never reach a field, a binding, or an ordinary `func` type.

A third case has no representation at all: §5.2 ⊢ a range expression `a..b`
"has no type and cannot be bound, returned, passed, or stored," so it never
reaches this package as anything.

## Package layout

| File | Contents |
|---|---|
| `type.go` | Package doc, the `Type` interface, `Mode` (parameter/receiver convention), `Marker` (`FunctionMarker`) |
| `basic.go` | `BasicKind`, `BasicInfo`, `Basic`, the `Typ` singleton table, predeclared type/tensor-element lookup, `Default` |
| `composite.go` | Every non-basic, non-named `Type`: `Ownership`, `Array`, `Slice`, `Map`, `Tuple`, `Chan`, `Pointer`, `Vector`, `Predicate`, `Signature`, `Expected`, `Tensor` |
| `named.go` | `Named`, `Struct`, `Field`, `Enum`, `Variant`, `Abstract`, `Family`, `TypeParam` |
| `const.go` | `Value` (the compile-time constant representation), constructors, folds (`Neg`, `Not`, `BinaryOp`, `Shift`, `Compare`), `Sizes.Representable` |
| `constraint.go` | `Constraint`, `Term`, the predeclared `Any`/`Comparable`, `Satisfies` |
| `object.go` | `Object` and its kinds: `Var`, `Const`, `Func`, `TypeName`, `Builtin`, `PkgName`, `Package` |
| `scope.go` | `Scope`, `Universe` and its `init()`, `Reserved` |
| `predicates.go` | `Underlying`, the `As*`/`Is*` family, `Identical`, `IsComparable`, `AssignableTo`, the two implicit conversions, `ConvertibleTo`, `TensorElemConvertible` |
| `sizes.go` | `Sizes`, `Alignof`/`Sizeof`/`Offsetsof`, `CopyKind`/`CopyCost` |
| `info.go` | `Info`, `TypeAndValue`, `Selection`, `Instance`, the `Record*` writers and accessors |
| `string.go` | `TypeString`, `ExpectedString`, `ConstraintString`, `ObjectString` |

## `Mode` and `Marker` (`type.go`)

```go
type Type interface {
    Underlying() Type
    String() string
}
```

`mut T`/`var T` and `unique T`/`shared T`/`weak T` wear the same syntactic
hat but get two different representations, on purpose. §3.2 ⊢ `mut`/`var`
are legal "a parameter or receiver only," so they're a `Mode` on a `Var`
(`ModeNone`, `ModeMut`, `ModeVar`) rather than part of the `Type`. §3.2 ⊢
`unique`/`shared`/`weak` are legal "anywhere a `Type` is," so they're an
`*Ownership`. The split is what makes a receiver's three-way shape fall out
for free: `(x: shared T)` is a `ModeNone` `Var` over an `*Ownership`; `(x:
mut T)` is a `ModeMut` `Var` over a bare `T`. It's also why "qualifiers do
not stack" needs no special case here — a `Var` carries exactly one `Mode`,
so `mut var T` is unrepresentable, and `mut shared T` is representable and
left for the analyzer to reject.

`Marker` (`MarkerNone`, `MarkerAsync`, `MarkerGPU`, `MarkerNPU`,
`MarkerTest`) lives on `Signature` because §7.4 ⊢ "the marker is part of the
function's type... checked at the declaration and again at every call. The
marker must agree at both ends." That's also why §4.1's `func(int32)` not
being assignable to `func(int32) async` needs no rule of its own:
`Identical` already compares markers, and assignability reduces to identity.

## Basic types and the untyped kinds

`Typ` (`basic.go`) holds exactly one `*Basic` singleton per `BasicKind`, so
identity comparison on `*Basic` **is** type identity. The kinds fall into
three families:

- The `PredeclaredTypeName`s — `bool`, `int`/`int8`…`int64`,
  `uint`/`uint8`…`uint64`, `float32`/`float64`, `char`, `string`. `byte`
  deliberately shares `Typ[Uint8]` rather than getting its own entry — §2.3
  ⊢ "`byte` is an alias for `uint8`, not a distinct type" falls straight out
  of the pointer comparison in `Identical`, at the cost that a diagnostic
  written against `byte` reports as `uint8`.
- The tensor element types — `BF16`, `FP8E4M3`, `FP8E5M2`, `Int4`. Ordinary
  identifiers in the same implicit scope as the type names, but §2.3 gives
  them their own legality rule ("only inside an `npu` body"), which is why
  they live in a separate `predeclaredTensorElems` table reached through
  `LookupTensorElem` rather than `LookupPredeclared` — not because the
  scanner tells them apart (it doesn't), but because the rule is separate.
- The untyped kinds — `UntypedBool`…`UntypedString`, plus `UntypedNil`.
  §4.1 ⊢ "a literal has no type until it lands," so these exist only
  between a literal and its destination and never appear in a declared
  signature. `Default` is the fallback an untyped constant takes when
  nothing imposes one (`let x = 1` needs an answer: `int`). `UntypedNil` has
  none — §10 gives it no type of its own, so a destination that isn't a
  `typed_ptr` is rejected rather than defaulted.

`BasicInfo` is a bit set (`InfoBoolean`, `InfoInteger`, `InfoUnsigned`,
`InfoFloat`, `InfoChar`, `InfoString`, `InfoUntyped`, `InfoTensorElem`), with
two derived masks: `InfoNumeric = InfoInteger | InfoFloat`, and
`InfoOrdered = InfoNumeric | InfoString`. `char` is pointedly absent from
`InfoOrdered` — §3.5 makes it comparable but not ordered, so `'a' < 'b'` is a
compile error. The flags are named `Info*` rather than `Is*` because the
`Is*` names in `predicates.go` are the package's actual API; the flags are
just the representation those functions read.

## Constant values (`const.go`)

`Value` is the compile-time constant representation, held at unbounded
precision (`*big.Int` for integers, `*big.Rat` for floats) for exactly one
reason: §4.1 ⊢ "a literal whose value does not fit its destination is a
compile error, not a wraparound," and that check is only possible if the
value hasn't already been narrowed by the time the destination is known.
It's also the arithmetic §5.3's constant expressions evaluate in — an
`ArrayLength`, an enum discriminant, a top-level `VarDecl` initializer, a
`switch` case pattern, and `new`'s `align:` argument.

- `MakeBool`/`MakeInt`/`MakeInt64`/`MakeFloat`/`MakeChar`/`MakeString`
  construct values; `Int64Val`/`BoolVal`/`StringVal`/`FloatVal` extract them.
- `Neg` folds unary minus — grammar.md ⊢ "there is no literal syntax for a
  negative number," so `-128` is unary minus on `128`, and this fold has to
  run *before* representability is checked, or `-128` would wrongly be
  rejected against `int8`.
- `BinaryOp` folds `+ - * / % | & ^`; a constant division or remainder by
  zero returns `Unknown` rather than trapping, since §5.5's trap is a
  runtime tier and a constant divide-by-zero is the caller's diagnostic. The
  wrapping operators `&+ &- &*` are absent on purpose — they wrap at a
  width, and an untyped constant doesn't have one yet.
- `Shift` folds `<< >>`; a count at or beyond the operand's width is §5.5's
  trap, not folded here.
- `Compare` folds `== != < <= > >=`, returning whether the operands were
  comparable as constants at all — legality of the operator on their type is
  §3.5/§5.1's question, answered elsewhere.

`Sizes.Representable` is the single enforcement point for §4.1's
no-silent-truncation rule, and it's a method on `*Sizes` rather than a free
function because whether a literal fits is a question about the target:
`int`/`uint`'s bounds come from `Sizes.intRange`, computed on demand from
`WordSize` rather than a precomputed table, since the same literal can be
representable under one build tag's `Sizes` and not another's. A float
constant with a nonzero fractional part is never representable in an integer
type — §4.2 gives `float → integer` only as an explicit `as` conversion that
truncates toward zero, so there's no implicit rounding to fall back on.

## Composite types (`composite.go`)

Every non-basic, non-named `Type` lives here. Each constructor is a thin
wrapper; what's worth knowing per type is which parts are a spec consequence
and which are this implementation filling a gap:

| Type | Spelling | Note |
|---|---|---|
| `Ownership` | `unique T`, `shared T`, `weak T` | Legal "anywhere a Type is" (§3.2) — a keyword form (`HeapConstructor`), not a builtin, so `scope.go` must not put it in `Universe`. |
| `Array` | `[N]T` | `N` is resolved to `int64` at construction, not kept as an expression (§5.3); it's part of the type's identity — `[8]int32` and `[16]int32` are unrelated (§3.1). |
| `Slice` | `[]T` | One of §3.4's one-word indirections; usable to break a recursive type. |
| `Map` | `map[K]V` | `K` must satisfy `comparable`; the analyzer checks this at construction and raises `NonComparableKey`. |
| `Tuple` | `(T, T, ...)` | Holds `*Var`s, not bare `Type`s, since a `TupleElem` may be named — an unnamed element is a `Var` with an empty name. Also serves as every `Signature`'s result list. grammar.md ⊢ "a tuple has at least one element; there is no unit type" — an empty `Tuple` appears in exactly one place, a void-form `Signature`'s results, and `Signature.IsVoid` is the question to ask instead of comparing against a unit value that doesn't exist. |
| `Chan` | `chan T` | One word (§3.4); zero value is a closed, empty channel (§3.3), so an uninitialized declaration is well-defined. Carries no direction — `thread`/`async` "hand back a `chan T`" is a property of the launch expression's result (§7.4), not a second channel type. |
| `Pointer` | `typed_ptr T` | §10's third tier — the one type "these rules do not reach"; §8.4 ⊢ two copies are "two unchecked aliases." May never be the direct base of another pointer type or a receiver type, both static rules over a declaration rather than shapes this type records (§3.2). |
| `Vector` | `vector[T, N]` | `N` is resolved to `int64` at construction, one of §5.3's bare-literal-token positions. Legal anywhere except a `gpu`/`npu` body or signature, a foreign boundary, or a map key (§3.2, again a static rule, not a shape constraint). |
| `Predicate` | *(none)* | The result of a vector comparison. §5.1 ⊢ it "has no source spelling and may not be an `if` condition, a `&&` operand, a field, or a channel element." It gets a `Type` anyway so the checker can carry it to `blend`. Because it can never be a field or channel element, it never reaches storage — `Sizeof` answers `0`. |
| `Signature` | `func(...) -> ...` | Type of every function, method, initializer, deinitializer, function literal, foreign declaration, and `MethodRequirement`. `Recv` is `nil` for a free function. `results` empty is the void form (§7.1: no `void` type, no unit type). `Expected` is a separate field rather than a result-tuple element, since putting it there would make it reachable wherever a `Type` is, when §3.2 restricts it to a `FunctionDecl`/`MethodDecl` result. |
| `Expected` | `Expected(T, "...")`, `Expected(error, "...")` | The `ExpectedType` result form of a test function. `error` is recognized only in this production and mints no `TypeName`, so `Type()` answers `nil` for it. `Msg` is normative text compared against a `Diagnostic.Text()`. |
| `Tensor` | `tensor[T, dims...]` | Legal "inside an `npu` body or that function's own signature" (§3.2). `shape` is resolved to `[]int64` at construction, since each `ShapeList` entry is another of §5.3's bare-literal-token positions. `Elems()` (the total element count) is what `Sizeof` multiplies by element width, and what §4.3's `[N]T ↔ tensor[T, N]` launch conversion matches against an array length. |

## Named types (`named.go`)

`Named` is what a struct, class, enum, or `abstract` alias declaration
mints. §3.1 ⊢ declared types "are nominal: two declarations are distinct
types even with identical fields." A transparent `TypeAliasDecl` mints
nothing — the resolver substitutes its target and no `*Named` is ever built
for it. Only `= abstract` is nominal, "because each `abstract` alias is
distinct from every other."

Construction is deliberately two-step (`NewNamed`, then `SetUnderlying`
later): §1.1's order-independence lets a field name its own enclosing type
(`struct Node { next: typed_ptr Node }`, explicitly endorsed by §3.4 as the
way to break a cycle), so the resolver must bind the `Named` to its
`TypeName` before walking any field — a premature read of `Underlying()`
gets an honest `nil`, which every structural predicate already treats as
"not a struct, not an enum" via `predicates.Underlying`'s type switches,
rather than a panic. `LookupMethod` is a flat linear scan — there's no
inheritance or embedding in Vertex, so there's no promotion to walk.

- `Struct` is the shared underlying of both a `StructDecl` and a
  `ClassDecl` — §7.2 ⊢ "a class is byte-for-byte identical in layout to a
  struct and differs only in its member and method model." `Sizeof` never
  branches on `Class()`; construction syntax (composite literal vs.
  initializer call) and `===` (legal on classes only, §3.5) do. `Field`'s
  `HasDefault` records that `= Expression` was written; §7.2 ⊢ defaults are
  "evaluated at each construction for each omitted field," while §3.3 ⊢ a
  zero value "applies none of them" — defaults belong to construction, not
  to the type.
- `Enum` keeps variants in declaration order because §3.3 ⊢ the zero value
  is "the first declared variant, with any payload zeroed" — `Variant(0)`
  is load-bearing. An explicit discriminant must be a constant expression
  (§5.3); the sources don't say what an omitted one is, so this
  implementation continues from the previous variant, stated as a choice
  rather than a rule, same as the `int32` default when `discrim` is
  omitted from `NewEnum`. `UnitOnly` is a shape query used by `enumSize`
  and by the `as`-cast rule (§4.2's one-way enum → discriminant
  conversion), not a rule in itself.
- `Abstract` is the bare `abstract` of an `AliasTarget`. It carries a
  `Family` (`FamilyMemoryFlat`/`FamilyObjectGraph`) because §4.2's
  `abstract → typed_ptr T` cast is legal only where linkage is
  memory-flat — an object-graph handle "has no byte representation to
  point at." Its zero value (§3.3) is "legal only on an error path, paired
  with a non-empty string" — a rule about where the zero may appear, not
  about this shape.
- `TypeParam` always has a non-nil `Constraint` after resolution —
  grammar.md ⊢ "a bare `TypeParamName` is constrained by `any`"; a parsed
  `nil` means "not written." Group-distributed constraints (`[A, B:
  Number]` constraining both) are resolved over the already-parsed list
  and land here, not in the tree.

## Objects (`object.go`)

`Object` is what the analyzer resolves every `*ast.Ident` to. `Type()`
returning `nil` isn't an error state — §1.1's order-independence means every
package-scope name is inserted before any type resolves, so `nil` means
"not yet resolved" and `Typ[Invalid]` means "resolution failed" (predicates
treat the latter as already-diagnosed, to avoid cascading).

- **`Var`** — variable, parameter, receiver, field, or tuple element.
  `Mode` carries `mut`/`var` (kept off `Type`, per the split above).
  `Mutable` distinguishes `let` (fixed, requires an initializer) from `var`
  (rebindable, required for exclusive access or transfer, §6.1).
  `Assignable` answers the first two entries of §6.2's `AssignTarget` list
  (a `var` binding, or a field of an assignable value) — the rest is a
  question about the base's type and lives in
  `predicates.InteriorAssignable`. `Owning` reports only that a
  parameter's *position* takes ownership (`var T`), never which of copy
  or move happens at a given call — §8.1 ⊢ "the convention lives in the
  signature; only the owning one has a choice at the call," which is a
  call-site (marker) question.
- **`Const`** — an explicit enum discriminant, an `ArrayLength`, or a
  top-level binding whose initializer §6.1 requires to be a constant.
- **`Func`** — function, method, initializer, deinitializer, foreign
  declaration, or `MethodRequirement`. `init`/`deinit` get no separate
  object kind — §7.2 ⊢ they're "ordinary method names recognized by
  spelling," so `IsInit`/`IsDeinit` are name-and-receiver questions.
  `IsEntry` checks §1.4's `main` shape (no parameters, void, no marker,
  and the one non-`async` function `await` is legal in); the "exactly
  one, in a package named `main`" parts are the loader's and checker's
  job. `IsTest` checks §7.4's shape (no parameters, an `Expected` result
  or none) — the enclosing `build test` file tag is checked separately by
  the analyzer against `token.LicensesTest`.
- **`TypeName`** — binds a name to either a type or a `Constraint`. One
  kind serves both because grammar.md ⊢ a single-identifier
  `ConstraintElem` "parses as both a one-term `TypeSet` and a constraint
  name; resolution is by what the name denotes" — a second kind would
  turn that into a two-map lookup. `IsAlias` distinguishes a transparent
  alias (mints no `Named`, or mints one for a different object) from a
  nominal one; an alias to `abstract` *is* nominal and mints its own, so
  it answers `false` here.
- **`Builtin`** — a reserved builtin name. It carries `Typ[Invalid]`
  rather than a real `Signature`, since each builtin's shape is arity- and
  type-argument-specific (`sizeof` takes a `Type`, `new` takes a count
  plus named arguments, `resize` has two arities) — the checker dispatches
  on `Id()` instead. `unique`/`shared`/`weak` are absent from `BuiltinId`
  on purpose: they're keywords with their own `HeapConstructor`
  production and never resolve to an object. `TransferId` is bound to
  nothing — so `x.transfer()` diagnoses against ownership §3.3 (the real
  syntax is the `var` prefix) rather than as an unknown name; the object
  exists only to carry a fix-it.
- **`PkgName`** — an imported package's qualifier, scoped to one file
  (§1.3: an `import` in one file doesn't bring the qualifier into
  others). `Name` always comes from the imported package's own
  `PackageName`, never derived from the import path — there's no
  aliasing, dot-import, or blank-import form to pick something else from.
- **`Package`** — name, path, and scope, whose parent is `Universe` so
  predeclared names resolve by walking out far enough. §1.1 ⊢ Vertex has
  no visibility modifier, so a package scope has no exported subset.
  `Complete`/`MarkComplete` exist because a half-checked scope would
  answer a lookup with a `nil` type — an importer must only hand out
  complete packages.

## Scopes (`scope.go`)

`Scope` records `pos`/`end` extent so `LookupParent` can tell a
position-independent scope (package, universe — §1.1) from one where
visibility runs "from its declaration to the end of its scope" (§2.1,
local). `NewFuncScope` marks a scope that `defer` registers against — the
enclosing *function*, not block, since §6.6 ⊢ deferred calls run "in
reverse order of registration" at every exit edge of the function; a
`FunctionLit` opens its own because §7.3 ⊢ it "begins with all enclosing
parse context cleared." `Insert` drops the blank identifier rather than
inserting it (§2.4: `_` "introduces no binding and may be repeated
freely") and returns any existing object of the same name for the caller to
diagnose as a duplicate (§2.2: no overloading, one name per scope).

`Universe` is §2.3's implicit outermost scope, built once in `init()`:
every `PredeclaredTypeName`, every predeclared tensor element name, the two
constraint names (`any`, `comparable`), and every reserved builtin. None of
these families expresses its legality rule by membership alone — a
constraint name is legal only in a bracket position, a tensor element name
only inside an `npu` body, and a builtin may not be shadowed at all, which
is `Reserved`'s job rather than `Insert`'s. `Reserved` delegates to
`token.IsReservedBuiltin` rather than keeping its own list, "so the
non-shadowing guarantee has one home."

## Predicates (`predicates.go`)

`Underlying` strips one layer of naming and is what §3.1's `~T` is defined
over; every structural predicate below reaches for it first
(`AsBasic`/`AsNamed`/`AsStruct`/`AsEnum`/`AsSignature`/`AsVector`/`AsTensor`).

- **`Is*` family** — `IsBool`/`IsInteger`/`IsFloat`/`IsNumeric`/
  `IsString`/`IsChar` read `BasicInfo` directly. `IsOrdered` is numerics
  plus `string`; `char` is excluded (§3.5). `IsUntyped` deliberately
  doesn't go through `Underlying` — an untyped constant never reaches a
  `Named`, so a non-`*Basic` answers `false` via `is()`'s nil-safe guard.
  `IsIndirection` is exactly §3.4's seven one-word spellings that break a
  recursive type (`unique T`/`shared T`/`weak T`, `typed_ptr T`, `[]T`,
  `map[K]V`, `chan T`) — the same set `sizes.go` gives one word.
  `IsIdentityComparable` gates `===`/`!==` to classes (§3.5): `==` on a
  `typed_ptr` already compares addresses, so `===` doesn't apply to one.
- **`Identical`** — type identity, and per §4.1 the *whole* of
  assignability ("no subtyping, no coercion, no promotion"). Named types
  compare by declaring object, never by shape; two instantiations of one
  generic are identical only when their type arguments are. `Basic`
  compares `kind` alone — `byte`/`uint8` need no special case since they
  already share one `*Basic`. `Signature` compares `variadic`, `marker`,
  params, and results, but not `Expected` (not a `Type`; a `func` type can
  never carry one).
- **`IsComparable`** — §3.5's table, read literally: comparable is
  numerics, `bool`, `char`, `string`, `typed_ptr T`, enums (flat,
  regardless of payload contents), and any struct/class/tuple/`[N]T` whose
  every component is comparable; *not* comparable is `[]T`, `map[K]V`,
  `chan T`, `func` types, `vector`, `tensor`, `abstract`, and the three
  heap handles — a header compare "answers a question nobody asked," and a
  contents compare would be an implicit `O(n)` walk the language never
  performs. A `TypeParam` is comparable when its constraint is exactly
  `comparable`, embeds something that implies it, or is a type set whose
  every term is itself comparable — the last is what lets `[K: Ordered]`
  serve as a map key without separately writing `comparable`.
- **`AssignableTo`** — `Identical`, relaxed by exactly two things: an
  untyped literal (kind-by-kind — bool→bool, int→any numeric, float→float
  only, char→char only since `'A'` and `"A"` "never interconvert
  implicitly," string→string), and `UntypedNil` (assignable only to a
  `Pointer` — §10: `nil` belongs to `typed_ptr T` and nothing else).
  Ownership is explicitly *not* an assignability question — whether a
  `var T` parameter receives the original or a copy is the `transfer`
  marker's call-site question (§8.1), so a bare `T` argument is always
  assignable to a `var T` parameter.
- **The two implicit conversions (§4.3)** are deliberately outside
  `AssignableTo` — "neither reaches inside a body, and no third case is
  added anywhere." `LaunchConvertible` is `[N]T ↔ tensor[T, N]` at a
  `gpu`/`npu` launch site, element type and shape matching exactly.
  `PointerCastElidable` is `typed_ptr T → typed_ptr U` with a written
  destination, `as` elided, both sides pointer types.
- **`ConvertibleTo`** — §4.2's `as` table: numeric↔numeric, enum→its
  discriminant type (one-way — "there is no `n as Status`"),
  `typed_ptr T`→`typed_ptr U`, `typed_ptr T`↔integer, and
  `abstract`→`typed_ptr T` for a memory-flat family (nothing converts *to*
  `abstract`; there's no runtime type information to reconstruct one
  from). `char`↔integer is deliberately absent from the table — read as a
  closed enumeration, so a code point is reached through a builtin
  instead.
- **`TensorElemConvertible`** — the one place a type name is callable:
  `bf16(val)` is the constructor spelling §4.2 gives the tensor element
  types instead of `as`. Kept separate from `ConvertibleTo` since the two
  spellings aren't interchangeable in either direction; the reverse (bf16
  back out as float32) is spelled by neither and isn't licensed here.

## Layout and cost (`sizes.go`)

`Sizes` (`WordSize`, `MaxAlign`) answers every layout question, and answers
them as a method rather than a free function precisely because §2.3 ties
`int`/`uint` to "the target's pointer width" — genuinely different from
`int64`/`uint64` even where the widths agree. `SizesFor` picks `Sizes32`
(`WordSize: 4`) for `TagJS`/`TagWasm` and `Sizes64` (`WordSize: 8`)
otherwise, so `int`, `uint`, and any struct that embeds one really do differ
in size across that build split — which is exactly the distinction §2.3 is
protecting, and exactly why a file's width answer must come from its own
tag rather than a global constant.

`Offsetsof` lays fields out in declaration order with padding and no
reordering, so a reader can predict layout straight from the source (§7.2
already makes order observable, through destruction order). `Sizeof` pins
one fact directly from §3.4 — every type has a compile-time size, and the
seven indirections are all one word — and leaves everything else
(`string`'s `{ptr, len}` header, a closure's conservative two-word
`{code, env}` estimate, enum layout as tag-plus-largest-variant, a vector
aligning to its own rounded-up-power-of-two size rather than clamping to
`MaxAlign`) as an implementation choice the comments call out as such.
`Predicate`'s `Sizeof` is `0` on purpose — §5.1 keeps it out of every
storage position, so it never has a layout to report. `Tensor` sizing goes
through `elemBits` rather than `Sizeof` per element, since `int4` is
sub-byte and only ever exists packed inside one.

`CopyKind`/`CopyCost` price §8.5's bare-copy column (transfer is uniformly
`O(1)` and needs no classification): register move for scalars,
`typed_ptr`, `Signature`, `Abstract`, `Predicate`; recursive fieldwise copy
for `[N]T`, struct/class, tuple, enum, vector, tensor; deep copy for
`string`, `[]T`, `map[K]V`; alloc-and-deep-copy for `unique T`; refcount
increment for `shared T`/`chan T`; weak-count increment for `weak T`.
`CopyCost` takes a `Type`, not a declaration, because §8.5 ⊢ "under generics
the cost is fixed by the concrete type at instantiation, so a lint on large
owned types fires per instantiation, not per declaration."

## `Info`: the side table (`info.go`)

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

`info.go` is the only file that imports `ast`, and everything analysis
learns lands here — never written back onto the syntax tree. A `nil` map
means "do not record this," and each `Record*` method checks first, so a
caller wanting only `Uses` pays for only `Uses`.

- **`Types`** records the type and value of every expression, including
  nodes `ast` represents as `Expr` but that denote types (an `IndexExpr`
  standing in for a type-argument list, a `TypeName`'s `Ident`, an
  `OwnershipType`) — `TypeAndValue.IsType` is the recorded answer to that
  ambiguity.
- **`Defs`**/**`Uses`** map an identifier to what it declares/refers to. A
  blank identifier still gets a (nil-object) `Defs` entry, since `_` "never
  introduces a usable binding" but tooling still needs to know one was
  seen.
- **`Selections`** records what a `SelectorExpr` resolved to — a field, a
  method, a package member, or a member of one of the four keyword
  namespaces (`NamespaceMember`, a closed set that never resolves to a
  user object).
- **`Instances`** is the set `lower` enumerates for monomorphization: every
  call inside a generic body lowers to a direct call on the concrete type
  because every instantiation is monomorphized and checked once, at its
  own site.
- **`Transfers`** is the one table that changes generated code rather than
  merely licensing it — the `var` marker's presence is the entire
  difference between a move and a deep copy, and `lower` reads this table
  to decide which to emit.
- **`Defers`** groups each `DeferStmt` by enclosing (function) scope, which
  is what makes "reverse registration order on every exit edge" enumerable
  in `lower` — there's no unwinder, so the exit-edge set is finite and
  static.

`TypeOf` falls back to resolving through `ObjectOf` when an expression has
no direct entry (a bare identifier, say); `ObjectOf` checks `Defs` before
`Uses`, so a declaring occurrence answers itself.

## Rendering (`string.go`)

`TypeString` renders a `Type` back to source syntax and is normative text —
it's what a `diag` template's `%s` receives, so a change here changes what
an `Expected(error, "...")` result matches against. `ExpectedString` and
`ConstraintString` are separate functions rather than methods on those
types, since neither `Expected` nor `Constraint` is a `Type`.
`ObjectString` renders a full description (`"variable x of type int32"`,
`"constraint Ordered"`) for diagnostics.

Every concrete type's own `String()` just calls `TypeString(t)` — one
recursive writer (`writeType`/`writeTuple`/`writeSignature`), not
independent per-type logic, so there is exactly one rendering rule to keep
in sync with the grammar. A couple of details worth knowing: a nested
`typed_ptr` is written parenthesized (`typed_ptr (typed_ptr T)`) even
though that shape is never legal to construct in the first place, purely so
rendering round-trips cleanly; and `Predicate`, which §5.1 gives no source
spelling at all, renders as a prose description ("lane predicate of N
lanes") specifically so a reader can't mistake it for writable syntax.