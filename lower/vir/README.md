# lower/vir

```go
import "github.com/vertex-language/vertex/lower/vir"
```

The mechanical translator. A `*hir.Program` in, one `*vir.Module` per originating Vertex package out, emitted through `vvm`'s append-only builder API.

Every decision was already made in `lower/hir`. If this package ever needs a type switch on "is this an owning type" or "was this transferred," that logic belongs upstream. What is left is three kinds of work — translation, ordering, and naming — and nothing else.

It depends on `hir`, `token` (for `loc` lines), and `vvm/ir/vir`. Output is unverified: the builder API validates nothing, and `ir/verify.Verify` is `vvm`'s to run.

This package is named `vir` and imports `vvm/ir/vir` under the same name. A package clause declares no identifier in its own file scope, so the alias is unambiguous, and it keeps the emitted-side spelling identical to `vvm`'s own documentation — `vir.I32`, `vir.OpCall`, `vir.Ptr`.

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/vertex-language/vertex/lower/hir"
	virlower "github.com/vertex-language/vertex/lower/vir"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
	"github.com/vertex-language/vvm"
)

func emit(fset *token.FileSet, units []*hir.Unit, tag token.BuildTag) error {
	prog, err := hir.Lower(&hir.Config{
		Fset:  fset,
		Sizes: types.SizesFor(tag),
		Mode:  hir.ModeProgram,
	}, units)
	if err != nil {
		return err
	}

	// The triple enters the compiler here and nowhere earlier. hir never
	// saw one; the tag it was checked under is coarser by design.
	target := virlower.Target{Arch: "x86_64", OS: "linux", ABI: "gnu"}

	mods, err := virlower.Lower(&virlower.Config{
		Fset:   fset,
		Target: target,
	}, prog)
	if err != nil {
		return err
	}

	// Modules come back in hir's topological order — the order vvm's
	// importer expects. Root names the module holding the entry point.
	bin, err := vvm.BuildModules(mods, vvm.Options{
		Target: vvm.Target(target),
		Root:   virlower.Root(prog),
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "vertex: wrote a.out (%s-%s-%s)\n",
		target.Arch, target.OS, target.ABI)
	return os.WriteFile("a.out", bin, 0o755)
}
```

To read the emitted text instead of linking it — which is how the invariants below are checked:

```go
for _, m := range mods {
	if err := virtext.Encode(os.Stdout, m); err != nil {
		return err
	}
}
```

## Entry points

```go
func Lower(conf *Config, prog *hir.Program) ([]*vir.Module, error)
func Root(prog *hir.Program) string   // the module holding the entry point, for --root

type Config struct {
	Fset   *token.FileSet // resolves hir positions into loc lines
	Target Target         // written into every emitted module; see "The one target-shaped fact"
}

type Target struct {
	Arch, OS, ABI string
	Tiers         []string
}
```

`Target` is field-for-field identical to `vvm.Target` so a driver can convert rather than copy.

## Errors

One shape, `*Error`, spelled `lower/vir: …`. There is no `todo:`/`internal:` split as there is in `hir`: everything this package refuses is either a malformed `*hir.Program` or a construct with no vir spelling at all, and both are bugs upstream rather than gaps in the user's source. Never report one as a compile error.

Errors unwind out of the block and instruction walk by panicking with a private `bail`; nothing else in this package panics.

## Ordering, concretely

This is the only part of the package with any teeth. vir fixes a mandatory section order (spec §2.1) and has no forward references (§2.2). `hir`'s slices are in construction order, which is neither.

- **Structs.** `hir` appends a `*Struct` to its module *before* lowering the struct's fields, so a self-referential field can reach its enclosing type without looping. That means a struct-typed field's shape lands in the slice after the struct naming it. `orderStructs` does a post-order walk over by-value field types; a pointer field erases the edge, so a cycle would be a struct of infinite size and is reported as a bug — `semantics.md` §3.4 already forbids it. A field naming a struct owned by another module is an import reference, not a local ordering constraint, and contributes no edge.
- **Functions.** `hir`'s worklist appends a shell when a call site is *discovered*, so its `Funcs` slice runs caller-before-callee — exactly backwards. `orderFuncs` is a post-order walk over each function's own callees, visiting in reverse slice order. A plain reversal would be right for a tree of calls and wrong for a diamond, which is why the reversal is only the visit order and not the result. A qualified call names another module and is resolved by the importer, so it contributes no edge either.

vir §2.2 exempts *direct* self-recursion and nothing else, so **mutual recursion has no vir spelling at all**. It is reported here, where the cycle is known, rather than reaching `ir/verify` as an undeclared-name error.

That last one is a real language-level limitation, not a gap in this package. `lowering.md` §2.1 claims mutual recursion is broken by the fact that only bodies reference each other — that is true of a *signature*, but vir has no `fnsig` for a defined function to forward-declare it with, so the claim does not survive contact. Either vir grows forward declarations for `fn`, or Vertex does not have mutually recursive functions.

## Translation

`hir.Op` is a subset of vir's §4 opcodes spelled identically, so `instr.go` is a table. It is an explicit table rather than `vir.ParseOpcode(op.String())` because `hir` spells its pointer arithmetic `field.ptr` / `index.ptr` — the instruction form, suffix included — while vir's closed vocabulary spells the opcodes `field` / `index`. A name-based bridge works for sixty opcodes and fails silently for two.

Five rules cover everything the table does not.

- **Comparisons name their operand type, not their result type.** vir's `eq`/`slt`/`ult` family yields `i1` from a suffix naming what is being compared (`eq.ptr`, `slt.i32`), while `hir.Instr.Type` holds the result. The suffix is read off `Args[0]` for these.
- **The float comparison family loses its `f`.** `hir` disambiguates by mnemonic — `OpFlt` against `OpSlt` — and vir disambiguates by the operand type in the suffix, so `OpFlt` becomes `lt` and the suffix rule above carries the distinction.
- **A void-typed instruction has no suffix.** `memcopy dst, src, n` takes none; `store.i32 p, v` does, and `hir` passes the real type there. The test is on `hir`'s `Type` being void, not on the opcode being void-resulting, because `store` is both.
- **`call` takes no suffix at all.** vir derives a call's result type from the callee's declaration rather than from the site.
- **`index` is the one arity mismatch.** vir's `index` takes `(base, elem type, index)` and scales by the element's stride; `hir` scales in the front end and emits `(base, byte offset)`, because the bounds check and the multiply are already separate instructions above it whether or not the scale folds. This package passes `i8` as the element type, so a stride of one leaves the arithmetic `hir` already did untouched. It is the only place here that translates a shape rather than a name — if the scaling should move down, that is a change to `hir`'s `elementPtr`, not to this table.

Two smaller notes. An `i1`-typed integer operand is emitted as `true`/`false`: an `i1` only ever arrives as a comparison result or a branch condition, so that is the only place vir's bool operand form can come from. And `loc` lines are deduplicated against the last one emitted, so a run of instructions from one source line costs one line of text rather than twenty.

## Naming

`namespace` is the import path, `module` is the `PackageClause` name: `semantics.md` §1.3's "the path is a locator, the declared name is the qualifier," carried through to the linker unchanged. Idents follow `lowering.md` §3.1; `vvm` applies its own Itanium-style mangling on top and nothing here pre-mangles.

`export` goes on every package-level declaration, since Vertex has no visibility modifier. Monomorphized instances are the exception — they are internal, because vir has no `linkonce` and two modules instantiating `smaller[int32]` would collide in the flat namespace (`lowering.md` §3.3, §22.4).

`entry` forces a bare symbol even in a namespaced module, which is why `hir` names the shim `_Ventry` and lets the emitted symbol come out as `main` without colliding with the user's own ident.

## The one target-shaped fact

`hir` never sees a triple, and neither does most of this package. vir §2.1 makes `target` mandatory whenever `link` is present, and every Vertex module links `vertexrt` — so `Config.Target` is written into every emitted module, alongside a `link static "vertexrt"` line.

**The invariant is checkable by reading the emitted text:** a module lowered from ordinary Vertex source carries exactly one `link` line. A `declare` block adds its own `link` line and its own `extern` group beside it; nothing else does.

**No `extern` group is emitted for `vertexrt`, and this is unresolved.** `hir` reaches every runtime symbol as a *qualified* call into a builtins module — `callBuiltin` sets `Instr.Module` from the `builtins.Symbol` and records `Module.Import(sym.ImportPath())` — so at the IR level the runtime is an import dependency and a `module.name` call operand, while at the object level it is a link dependency. The two spellings are mutually exclusive at a call site, and this package emits the one `hir`'s data actually supports.

If the `extern` spelling is the intended one, the change is upstream: an `extern` group needs parameter and result types, and `hir`'s call sites carry only a symbol. Landing it means `builtins` exposing the §2.2 signature table and `hir` recording which symbols a module called. Recorded here rather than papered over, because choosing between them is a decision and this package makes none.

## Not emitted

- **`fnsig` declarations.** Nothing upstream produces one — `hir`'s `callTarget` `todo`s on a call through a function value — so the section is always empty. The indirect-call *spelling* is present in `instr.go` (`Instruction.Sig`), so landing `fnsig` upstream does not also mean landing it here.
- **Function addresses in operand position.** vir has no address-of instruction and no `addr` operand form — only `InitAddressOf`, inside a global initializer. `hir.VFuncAddr` in instruction position therefore has no spelling at all and is reported as a bug. It is the second upstream problem below, seen from the other side.
- **`export struct`.** `hir` declares synthesized headers (`_Vstr`, `_Vvec`, `_Vfn`, boundary tuples) per module, and a struct produces no symbol, so duplicating the declaration is safe for layout. It is *not* obviously safe for `byval`/`sret`, which compare nominally per origin. A cross-module field already emits `StructType.Import` and declares the import, which is correct for layout and untested for the byval/sret case — so the first cross-package `[]T` passed by value is what forces the decision, and it wants an `Exported()` call here and an origin claim upstream.
- **Bodyless shells.** A `hir.Func` with no `Blocks` is a foreign declaration or an object monomorphization found no syntax for. Foreign ones were already emitted into an `extern` group; emitting an empty `fn` here would produce a block with no terminator.
- **`extern_c`, `tls`, `inline`/`noinline`/`cold`, `readonly`.** No `hir` field carries them yet. `entry` and `noreturn` are wired.

## Two upstream problems this package will not paper over

Both are emitted or reported as written, because papering over either would be a decision, and this package makes none.

1. **Float division has no vir opcode.** `hir.binaryOpFor` returns `OpSDiv` for a float operand, commented "vir spells float division with the same mnemonic family" — but vir's `opTable` gives `sdiv` `ConstraintInt`, and §4.1 lists no float division at all. `sdiv.f32` is emitted faithfully and `ir/verify` will reject it. Either vir needs `fdiv` or `hir` needs to stop claiming one exists. Same shape as `lowering.md` §22.10's `frem` item, and the same fix does not apply — `%` can be restricted to integers, `/` cannot.

2. **`addr` in a global initializer names a later function.** vir §6.2 permits `addr ident` for "earlier functions/globals", but §2.1 puts the whole `global` section before the whole `fn` section, so a global holding a function address — `hir`'s only spelling for one, since vir has no address-of instruction — can never reference an *earlier* function. The declaration is emitted as `hir` built it. Either the spec means "resolved by relocation, order irrelevant" and §6.2's wording is loose, or function addresses need a form that does not exist yet.

## File organization

| File | Purpose |
|---|---|
| `vir.go` | `Target`, `Config`, `Error`, `Lower`, `Root`, and the §2.1 section walk — the one guarantee the top level makes. |
| `decl.go` | Structs (with the dependency sort), consts, globals and their init forms, the `vertexrt` link, declare-block links and extern groups, imports. |
| `func.go` | The callee-first function sort, signatures, blocks, terminators, and `loc` lines. |
| `instr.go` | The opcode table, the five suffix and arity rules, and the instructions with dedicated builder entry points. |
| `value.go` | `hir.Type` → `vir.Type`, `hir.Value` → `vir.Operand`, `hir.ConstInit` → `vir.ConstInit`, and params. |