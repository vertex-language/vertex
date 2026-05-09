# Vertex

> A unified computing model.
> Not a language. Not a runtime. Not a toolchain.
> Something new.

---

## The Problem With Every Platform Ever Built

Every computing platform in history asks the same thing of you: **buy in.**

- Java: buy the VM
- Node: buy the event loop
- Python: buy the interpreter
- Go: buy the runtime
- LLVM: buy the toolchain
- Docker: buy the container
- WebAssembly: buy the host

You get the platform's power, and you inherit the platform's ceiling.
The language, the memory model, the security model, the deployment model —
all of it is decided for you before you write a single line of code.

Vertex doesn't ask you to buy in. It compiles everything away.
What comes out is a native binary that speaks directly to the kernel.
No intermediary. No runtime. No platform tax.

---

## What Vertex Actually Is

Vertex is the layer between the logic a human is trying to compute
and bare metal.

It owns that space completely:

- **Language** — write in whatever is most ergonomic for how you think
- **Modules** — consume anything, from any language, ever written
- **Security** — capability model baked into the binary, not the environment
- **Platform** — one source, any target, native everywhere
- **Output** — a single binary. That's it. That's the whole thing.

The operating system provides syscalls.
Vertex provides everything above that.

The frontend is not part of the model. It is comfort.
The logic you are trying to compute is the invariant.
The frontend is just the notation you reach for to express it.

---

## The Architecture

### One Backend. Every Language.

Every frontend in the Vertex ecosystem compiles to the same place:
**Wasm IR** — WebAssembly Intermediate Representation.

```
C          →  vcc frontend   ─┐
C++        →  v++ frontend   ─┤
JavaScript →  vjs frontend   ─┤──→  Wasm IR Compiler  ──→  Native Binary
Vertex     →  vtx frontend   ─┤
FORTRAN    →  frontend        ─┤
Lisp       →  frontend        ─┤
...        →  frontend        ─┘
```

The Wasm IR compiler is the model.
The frontends are ergonomics — a translation from the notation you prefer
into a form the compiler already understands.

Wasm is a stack machine. Most imperative and functional languages map
naturally. Logic languages embed their runtime as a module.
Dataflow and actor models express their concurrency as a Wasm runtime layer.
Every computing paradigm ever designed reduces to a frontend.

The logic you are trying to compute does not change based on which
frontend you used to express it. The model is the constant.

---

### The Module System

Modules are stored as two files:

```
module.wat        ← Wasm IR text. The actual thing. Language-agnostic.
exports.json      ← Public API surface. The polyglot interface.
```

A module has no origin language. A `.wat` file compiled from C by `vcc`
is identical in format to one compiled from C++ by `v++` or from
JavaScript by `vjs`. Every frontend can consume every module.
No one knows or cares what notation produced it.

The `exports.json` is the resolved type surface — structured, versioned,
and read by every frontend through its own type lens.
It is the replacement for headers, manifests, and `.d.ts` files across
the entire ecosystem, forever.

```json
{
  "vcc_exports_version": "1",
  "module": "curl",
  "version": "8.7.1",
  "exports": [...],
  "types": {...},
  "capabilities": ["net", "fs"]
}
```

---

### The Module Registry

```bash
vcc get curl          # fetch, compile to .wat, store locally
vcc get sqlite
vcc get openssl

vcc build main.c      # resolves modules, links IR, emits one native binary
```

No flags. No `-l`. No linker invocations.
Modules are discovered from source declarations.
Everything statically absorbed into a single binary.

The registry is language-agnostic. A module published by a C author
is immediately available to every Vertex frontend without translation,
wrapping, or binding generation.

Third party authors don't choose a language anymore.
They publish `.wat` + `exports.json` once.
The entire Vertex ecosystem can use it.
The frontend they wrote it in is nobody's problem.

---

### The Capability Model

Security in Vertex is not a runtime check. It is a compile-time decision.

```bash
vcc build main.c --permissions=permissions.json
```

```json
{
  "net":     { "allow": false },
  "fs":      { "allow": true, "root": "/sandbox/app", "mode": "chroot-emulated" },
  "memory":  { "heap_max": "64mb" },
  "io":      { "stdout": true, "stderr": true, "stdin": false },
  "process": { "allow_fork": false, "allow_exec": false }
}
```

Disabling a capability means the syscall shims for that capability
are **never emitted**. The code does not exist in the binary.
There is no runtime check to bypass because there is nothing to bypass.

The result is a **self-describing binary** — its security contract is part
of the artifact, not the environment it runs in. You can hand the binary
to anyone and the contract travels with it.

---

### The wapi Syscall Layer

```
[ libc   ]    write(), malloc(), printf()         ← C programmers
     ↓
[ wapi   ]    __wapi_write(), __wapi_mmap()       ← compiler internal only
     ↓
[ kernel ]    NR_write, NR_mmap, NR_brk           ← Linux syscalls
```

`wapi` is a private namespace — never called by user code.
Inside the Go backend, syscalls are organized as `wasm::linux::write`,
`wasm::linux::mmap`, giving the codegen layer a clean registry that maps
directly to x86_64 instruction emission.

Every capability group maps to a set of wapi shims.
The capability model controls shim emission.
Tracing instruments shim boundaries.
Platform translation rewrites shim targets.

The wapi layer is the hinge of the entire system.

---

### Cross-Platform Native Compilation

```bash
vcc build main.c --target=linux-x86_64
vcc build main.c --target=linux-arm64
vcc build main.c --target=macos-arm64
vcc build main.c --target=windows-x86_64
```

The Wasm IR backend retargets per platform:

| Target          | Binary Format | Syscall Routing          |
|-----------------|---------------|--------------------------|
| linux-x86_64    | ELF           | direct kernel            |
| linux-arm64     | ELF           | direct kernel            |
| macos-arm64     | Mach-O        | libSystem (Apple's moat) |
| windows-x86_64  | PE/COFF       | wapi::windows            |
| browser         | .wasm         | host API                 |

One source. Every platform. No cross-compilation toolchain required.
The `.wat` module format is already WebAssembly text —
the browser target ships modules as-is.

---

### Runtime Visibility

```bash
vcc run main.c               # compile and execute
vcc run main.c --trace       # compile with instrumented wapi shims
vcc trace ./main             # inspect any compiled Vertex binary
vcc trace ./main --filter=fs # filter by capability group
```

```
[vcc trace] running ./main
────────────────────────────────────────────────────────────────
  0.000ms   wapi::io::write      fd=1  len=31        → ok
  0.031ms   wapi::fs::open       path="/etc/passwd"  → BLOCKED
  4.201ms   wapi::proc::exit     code=0              → ok
────────────────────────────────────────────────────────────────
47 syscalls  |  4.201ms  |  1 blocked  |  exit 0
```

Tracing operates at the wapi semantic layer — not the raw kernel.
You see intent. `wapi::fs::open`, not `openat(AT_FDCWD, ...)`.
The output shows what the binary attempted and what the
capability model enforced, inline, in order.

---

## The Frontends

Frontends are ergonomics. They are the notation you are comfortable with.
They do not define what you can compute. They do not define what modules
you can use. They do not define what platforms you can target.
They are the on-ramp. The model is the road.

### vcc — C
The origin. Zero-dependency ELF binaries. Direct syscalls.
No libc required. No LLVM. No GCC.

### v++ — C++
Full C++ frontend targeting the same Wasm IR backend.
Shares every module vcc can consume.

### vjs — JavaScript
JavaScript with memory model adjustments for systems-level output.
Compiles to Wasm IR. Shares the module ecosystem.
Where scripting meets bare metal.

### Vertex — The Native Language
Designed from the ground up to feel native to the Wasm IR model.
Imports C, C++, and JS modules without bindings or wrappers.
The language that knows what it's running on.

---

## What Gets Absorbed

Any computing paradigm reduces to a frontend.
The logic being computed does not change. Only the notation does.

| Paradigm             | Map to Wasm IR via                          |
|----------------------|---------------------------------------------|
| Imperative           | direct — natural fit                        |
| Functional           | direct — closures as table references       |
| Object-oriented      | vtables in linear memory                    |
| Logic / Prolog       | runtime embedded as a module                |
| Stack / Forth        | trivial — Wasm is a stack machine           |
| Dataflow             | runtime scheduler as a module               |
| Actor model          | message queues + scheduler as modules       |
| Hardware description | simulation layer as a module                |
| SQL                  | query engine as a module                    |
| APL / array langs    | SIMD-mapped IR operations                   |

There is no ceiling. A frontend author permanently expands
what every other frontend can consume.

---

## What Vertex Is Not

- **Not a VM.** Nothing is interpreted at runtime.
- **Not a container.** The binary enforces its own security contract.
- **Not a managed runtime.** GC is opt-in via module, not baked in.
- **Not an abstraction layer.** It compiles away. Nothing is left over.
- **Not a framework.** It has no opinions about your application structure.
- **Not a language ecosystem.** It is the unified model all ecosystems run on.

---

## The Endgame

Every language ever written becomes a frontend.
Every library ever built becomes a module.
Every platform becomes a backend target.
Every binary is self-describing, self-securing, and dependency-free.

The source language is a comfort, not a constraint.
The target platform is a flag, not a rebuild.
The security model is in the artifact, not the environment.
The module ecosystem is universal, not siloed.

Vertex is not trying to replace any language.
It is trying to make the choice of language irrelevant
to everything that happens after you write the code.

The logic you are trying to compute was always the point.
Vertex just refuses to lose it in translation.

---

## Status

| Component              | Status         |
|------------------------|----------------|
| vcc — C frontend       | v1 active      |
| Wasm IR emission       | v1 active      |
| x86_64 ELF backend     | v1 active      |
| wapi syscall layer     | v1 active      |
| libc layer             | v2 planned     |
| Capability model       | v2 planned     |
| vcc run / vcc trace    | v2 planned     |
| Module system (.wat)   | v2 planned     |
| vcc get registry       | v2 planned     |
| Cross-platform targets | v3 proposed    |
| v++ frontend           | roadmap        |
| vjs frontend           | roadmap        |
| Vertex language        | roadmap        |

---

*The layer between the logic a human is trying to compute and bare metal.*
*The frontend is just ergonomics. Everything else compiles away.*