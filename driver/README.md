# driver

`github.com/vertex-language/vertex/driver`

`driver` is the top-level orchestrator: the one package allowed to import
the frontend (`parser`, `importer`, `analyzer`, `types`), the lowering
chain (`lower/hir`, `lower/vir`), and `vvm` all at once, and to pick the
right combination for a given invocation. Nothing below it knows it
exists. `cli/` parses flags and calls in here; it makes no pipeline
decisions.

## The pipeline

```
<file.vs | dir>
     │  load.go
     ▼
[]*driver.Package         checked, dependency-first
     │  lower.go: hir.Lower       — every decision made here
     ▼
*hir.Program
     │  lower.go: lower/vir.Lower — mechanical
     ▼
[]*vir.Module             unverified; ir/verify is vvm's to run
     │
     ├─ emit.go: format/vbyte/{text,binary}.Encode ─► .vir / .vbyte
     └─ emit.go: vvm.BuildModule[Graph] ────────────► native image
```

## Two load paths, on purpose

A directory is a package, and `importer.Load` does the whole job. A single
`.vs` file is *that file* — `parser.ParseDir` is directory-granular by
construction, so `loadFile` parses one file, wraps it in an `ast.Package`,
and runs `analyzer` directly against an `Importer` that delegates each
declared import back to `importer.Load`. Compiling `main.vs` must not
silently pull in a sibling `scratch.vs`.

## What the target table does and doesn't hold

`target.go` maps a Vertex triple (`linux-arm64`) to a `vvm.Target`
(`aarch64-linux-gnu`) — two vocabularies with two owners, translated in one
place, the same discipline `vvm`'s `dispatch.go` keeps between `vvm.Target`
and each linker's own. A triple appears only if vvm has a `cpu/lower`, an
`object`, and a registered linker (or a flat writer) for it. `riscv64`,
32-bit x86 and ARM, wasm, and js are absent for that reason, not by
oversight.

## What this package deliberately does not do

- **Verify.** `ir/verify.Verify` runs inside `vvm.BuildModule`. Calling it
  here would be a second, drifting copy of the same gate.
- **Emit objects, assembly, or MIR.** `vvm`'s object and assembly stages
  are internal to its build pipeline. Reaching into `cpu/lower` and
  `object/*` directly would make this package a second dispatcher; the CLI
  rejects those flags with the reason instead.
- **Device offload.** There is no `lower/gvir`, so no `.gvir` is ever
  produced and `vvm.BuildDevice` is never called. `gpu`/`npu` launch sites
  type-check; the kernel body has nowhere to go, and `lower/hir` reports it
  as a `todo:`.

## Seams

Every construction of a `hir.Config`, `hir.Unit`, or `lower/vir.Config`
lives in `lower.go` under the "the seam" banner and nowhere else, so a
rename in either lowering package is a one-file edit here.