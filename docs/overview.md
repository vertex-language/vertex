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
