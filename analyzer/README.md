# analyzer

```go
import "github.com/vertex-language/vertex/analyzer"
```

Package `analyzer` checks a parsed Vertex package against the static rules and produces a `*types.Package` plus a `types.Info` side table. It depends on `ast`, `token`, `diag`, and `types`; nothing depends on it back.

## Design philosophy

The parser resolves exactly one thing by context — the literal-in-header ambiguity — and accepts everything else unconditionally. Every place grammar.md says **"static rule"** means the form derives, parses, and is checked afterwards over an already-parsed tree. This package is that check, and `diag`'s ranges outside 1xxx and 2xxx are its diagnostic surface.

**Nothing here writes to the syntax tree.** Every result lands in a `types.Info` side table, which is what lets a printer round-trip a checked file and lets the checker run twice over one tree without residue.

**The checker does not read back out of `Info`.** `Info`'s maps are optional by design (a nil map means "do not record this"), so the checker keeps its own `defs`/`uses`/`fileScopes`/`funcObj` maps and mirrors into `Info`. A caller that allocated only `Uses` still gets a full check.

**Citation convention.** A bare `§` cites semantics.md; CamelCase names are grammar.md productions — the same scheme `types` uses. Where neither document fixes something (the variant-tag set, what an omitted enum discriminant continues from), the comment says so rather than dressing a choice as a rule.

## Package layout

| File | Contents |
|---|---|
| `check.go` | Package doc, `Config`, `Importer`, `Checker`, `Files`, scope helpers, `declare`/`lookup` |
| `errors.go` | `report`/`errorAt`/`errorExpr`, and `invalid()` |
| `decl.go` | Phase 1 (`collectObjects`, `collectImports`, `collectDecl`, `collectForeign`, `collectMethods`) and phase 2 (`objDecl` and one resolver per declaration kind) |
| `typexpr.go` | `typ` and `constraintExpr`, the two entry points from an `ast.Expr` in type position |
| `resolve.go` | Phase 3: the body walk, `expr`/`stmt`/`owning` |
| `const.go` | The bounded constant folder for §5.3's constant-expression positions, plus literal decoding |

## Entry points

```go
type Config struct {
	Fset     *token.FileSet
	Tag      token.BuildTag
	Reporter diag.Reporter
	Importer Importer // may be nil
}

type Importer interface {
	Import(path string) (*types.Package, error)
}

func NewChecker(conf *Config, path, name string, info *types.Info) *Checker
func (*Checker) Files(files []*ast.File) (*types.Package, error)
func (*Checker) Package() *types.Package
```

A nil `Importer` reports every import as an undeclared name — correct for a single-package check, and the seam a loader fills. A nil `info` gets a fresh `types.NewInfo()`. Errors accumulate rather than aborting; `Package()` is valid even after them, so editor tooling can use whatever resolved. `Files` calls `MarkComplete` only on a clean check, since §1.3 requires an importer to hand out only complete packages.

## The three phases

```go
c.collectObjects()   // phase 1: names exist, types are nil
c.resolveDeclTypes() // phase 2: types filled in, cycles caught
c.checkBodies()      // phase 3: bodies walked, uses recorded
```

§1.1 makes top-level declarations order-independent, and that is what forces the split: every package-scope name must exist before any declaration's type resolves, and every type must be known before any body that calls it is walked. Collapsing two would reintroduce a forward-declaration form the language does not have.

- **Phase 1** inserts each package-scope name with a **nil type**. Imports bind into a *file* scope, since §1.3 gives a qualifier file scope while the declarations using it are package-scoped. Methods are collected separately in `collectMethods`, after every type name exists, because a receiver may name a type declared later or in another file. A top-level `let` mints a `*types.Const` and a top-level `var` a `*types.Var`.
- **Phase 2** iterates `objMap` and calls `objDecl`. Order is unspecified and irrelevant: `objDecl` is idempotent (`obj.Type() != nil` short-circuits) and pulls in dependencies on demand. A `resolving` stack turns a self-reference into `diag.TypeCycle` instead of an overflow; because order-independence makes a cycle reachable from any entry point, the guard lives in `objDecl` rather than in a pre-pass. Resolution runs in *that object's declaring file's scope*.
- **Phase 3** walks bodies, resolving names to objects. It computes no expression types — resolution alone settles the parser's deferred forks.

## What resolution alone settles

| Deferred by the parser | Settled by |
|---|---|
| `a[i]` vs. `Stack[int32]` | `denotesType` — whether the operand denotes a generic declaration |
| `&x` as address-of vs. dereference | left open: the operand's *type* decides, not its name |
| `~T` as bitwise-NOT vs. underlying-type | position — `typ` rejects it, `constraintExpr` accepts it |
| single-identifier constraint element | a scope lookup: a constraint name embeds, anything else is a one-term set |

## Owning positions and the transfer marker

`ast` marks exactly six positions as owning: a `VarDecl`'s values, an `AssignStmt`'s values, an `ArrayLit`'s elements, a `CallExpr`'s arguments, a composite literal's field values, and a tuple element. `resolve.go` mirrors that with two functions — `owning(e)` accepts a `*ast.TransferExpr`, and `expr(e)` reports one as `TransferOutsideOwning`. An accepted marker is recorded into `Info.Transfers` against the binding that dies, which is the one analysis result that changes generated code rather than merely licensing it.

## Context tracking

`bodyCtx{npu, async}` is established at every function boundary from the signature's marker and cleared for a `FuncLit`, because a literal "begins with all enclosing parse context cleared" — a closure inside an async body may not `await` unless itself marked. `main` sets `async` too, per §1.4. `npu` gates `TensorOutsideNpu` and the tensor-element rules; `async` gates `AwaitOutsideAsync`.

`inSignature` and `inTensorElem` separate the two tensor-element rules, which differ by position rather than shape: a bare `bf16` is illegal in a signature and outside an npu body, while the element slot of a `tensor[...]` is exempt from both. That exemption is this implementation's reading and is marked as such in the code.

## The constant folder (`const.go`)

Bounded on purpose. It folds literals, named constants, unary `-`/`!`, and binary arithmetic/shift/comparison in `types.Value` at unbounded precision, because §4.1 makes a non-fitting literal a compile error rather than a wraparound and the check is lost once the value is narrowed. It serves an `ArrayLength`, a `ShapeList` entry, a vector lane count, an enum discriminant, and a top-level initializer — and `Sizes.Representable` is applied to the last two, since whether a literal fits is a question about the target.

Anything richer returns `Unknown` and the caller diagnoses, so widening this later cannot change a program that compiles today. Integer literals are decoded with an explicit base rather than base-0, since there is no prefix-free octal form — `0600` is decimal 600.

## Codes not raised here

Deliberate gaps, each for a stated reason:

- `InfiniteType` — needs a value-containment walk over resolved types; not yet written.
- `UnknownVariantTag` — the variant-tag set is closed, but nothing in `token` holds its membership. The check belongs wherever that set eventually lives, beside `LookupBuildTag`.
- `SelectCaseNotChannelOp`, `EnumShorthandNoType`, `NotRepresentable` in bodies, and most of 4xxx — these need expression types, which this pass deliberately does not compute.