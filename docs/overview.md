# Vertex Language — Project Overview

## What is Vertex?

Vertex is a modern systems language designed for high-performance graphics and AI workloads.
It replaces the need for CMake, vcpkg, and aging compiler toolchains with a clean,
integrated build and package system.

The lineage is simple:

**C → C++ → Vertex**

---

## Repository Structure

| Repo | Role | File Ext |
|------|------|----------|
| `vertex-language/vertex` | Vertex language frontend | `.vs` |
| `vertex-language/vcx` | C++ frontend | `.cpp`, `.hpp` |
| `vertex-language/vcc` | C frontend | `.c`, `.h` |

All three repositories are **frontends only** — clean, modern, and focused solely
on parsing and language concerns.

---

## WebAssembly as the Universal IR

The entire toolchain is unified by a single contract: **everything compiles to `.wasm`**.

`vcc` compiles C to a `.wasm` module.
`vcx` compiles C++ to a `.wasm` module.
`vertex` compiles `.vs` source to a `.wasm` module.

A `.wasm` binary is a fully validated, typed, self-describing module. It carries
its exported function signatures, its memory layout, and its code — everything
the consumer needs. There is no separate header file. No ABI negotiation.
No platform-specific linking step.

```
C   (.c)    ──▶  vcc  ──▶  module.wasm  ─┐
C++ (.cpp)  ──▶  vcx  ──▶  module.wasm  ─┼──▶  Vertex imports it directly
Vertex (.vs)──▶  vertex ▶  module.wasm  ─┘
```

Because the output format is always identical `.wasm`, modules produced by any
of the three frontends are **fully interchangeable and plug-and-play**. Vertex
does not know or care which language produced a module it imports — it sees
a validated WASM binary and calls into it directly.

---

## What This Means in Practice

A C++ math library compiled with `vcx`:

```cpp
// math.cpp
extern "C" float sqrt(float x) { ... }
```

Becomes a `.wasm` module with a typed, validated export. Vertex consumes it
with zero friction, zero overhead, and no knowledge of the source language:

```vertex
import "math.wasm" as mathLib

export func calculateDistance(a: Point, b: Point) f32 {
    let dx = a.x - b.x
    let dy = a.y - b.y
    return mathLib.sqrt((dx * dx) + (dy * dy))
}
```

The call to `mathLib.sqrt` compiles to a direct WASM function call. No FFI.
No marshalling. No runtime bridge. The module boundary is transparent.

---

## Native Output

Vertex code — including any imported `.wasm` modules — is compiled to a
native binary by `wasm-compiler`. The output is a standalone ELF (or future
PE / Mach-O) binary that runs directly on the CPU with no VM, no interpreter,
and no runtime overhead.

```
your .vs source
      │
      ▼
  vertex frontend    ──▶  module.wasm
                                │
      math.wasm  ──────────────▶│
      simd.wasm  ──────────────▶│
                                ▼
                        wasm-compiler
                                │
                                ▼
                        native binary (ELF)
                        talks to OS via WAPI
```

---

## All Compilers Written In Modern C++

All repos, including `wasm-compiler` itself, are implemented in modern C++
and Go respectively — not the legacy toolchains they aim to replace.
Once bootstrapped, `vcc` and `vcx` compile themselves.

---

## Goals

- **No CMake** — replaced by Vertex's integrated build system
- **No vcpkg** — replaced by Vertex's native package manager
- **No legacy g++/gcc dependency** — `vcc` and `vcx` handle C and C++ compilation
- **No ABI negotiation** — `.wasm` modules are the universal, self-describing contract
- **One toolchain** — install Vertex, get everything

---

## File Extensions

- `.vs`   — Vertex source files
- `.wasm` — compiled module output from any frontend; universal import format