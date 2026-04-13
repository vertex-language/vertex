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
| `vertex-language/vertex` | New modern language frontend | `.vs` |
| `vertex-language/vcx` | C++ frontend (compiler support) | `.cpp`, `.hpp` |
| `vertex-language/vcc` | C frontend (compiler support) | `.c`, `.h` |

All three repositories are **frontends only** — clean, modern, and focused solely 
on parsing and language concerns.

---

## Compute Framework

The backbone of all three compilers is the **Compute Framework** — a shared internal 
library responsible for:

- **IR (Intermediate Representation)** — unified code representation across all frontends
- **Builder Patterns** — structured, composable code generation pipelines
- **Backend targets** — what all three frontends compile down to

This means `vertex`, `vcx`, and `vcc` don't duplicate compiler internals. 
They each speak to Compute Framework and let it do the heavy lifting.

```
vertex (.vs)  ─┐
vcx    (.cpp) ─┼──▶  Compute Framework (IR + Builder)  ──▶  Native Output
vcc    (.c)   ─┘
```

---

## All Compilers Written In Modern C++

All repos, including Compute Framework itself, are implemented in **modern C++** — 
not the legacy toolchain they aim to replace, but clean contemporary C++ 
compiled through Compute Framework itself once bootstrapped.

---

## Goals

- **No CMake** — replaced by Vertex's integrated build system
- **No vcpkg** — replaced by Vertex's native package manager
- **No legacy g++/gcc dependency** — `vcc` and `vcx` handle C and C++ compilation
- **One toolchain** — install Vertex, get everything

---

## File Extensions

- `.vs` — Vertex source files
- `.vsp` — Vertex package manifest (proposed)

---

## Status

Early development. All repositories currently being scaffolded.
