# Core Architecture & Execution Flow

This document defines the foundational architecture and build order for a modern,
from-scratch programming language in 2026. The system is built in three discrete,
independently testable layers — each one depending only on the layer below it.

The guiding philosophy is simple: **build the bottom of the stack first.**
You should be able to run raw `.wasm` files before you write a single line of your
own language's syntax.

---

## The Three Layers

```
Layer 3 — Your Custom Frontend
         ↓
Layer 2 — WASM Byte Builder (Go)
         ↓
Layer 1 — WASM Engine (C++)
```

Build from the bottom up. Each layer has a hard contract with the one below it.
Each layer is independently testable. A bug is always traceable to exactly one layer.

---

## Layer 1 — The WASM Engine (C++)

**Build this first. Everything else depends on it.**

This is a standalone binary that takes any valid `.wasm` file and runs it at native
speed. No browser. No JavaScript runtime. No legacy baggage. The WASM specification
is just a document — your engine is a lean, systems-level implementation of it that
you own entirely.

### What it does

**Loader & Validator** reads the binary `.wasm` format and verifies it is structurally
legal before a single instruction executes. Malformed or malicious bytecode is
rejected at this stage.

**Host API Binding** is where the outside world enters the sandbox. The WASM
specification intentionally defines no I/O of its own. Your engine injects it by
binding C++ functions — filesystem access, stdin/stdout, sockets, OpenGL calls —
into the WASM import namespace. The running bytecode calls these by name and the
engine fulfills them natively.

**AOT Compiler** translates the generic, stack-based WASM instruction set into
register-based machine code tuned for the exact CPU the engine is running on. This
is where native performance comes from.

**Linear Memory Sandbox** allocates a single contiguous byte array at startup. All
bytecode memory operations are bounds-checked against it. A running program is
mathematically incapable of touching memory outside its assigned region.

### Milestone — Layer 1 is complete when

- The engine loads and validates a reference `.wasm` binary from an external tool
- A hello world program compiled by `wasi-sdk` or `emcc` runs correctly
- A program performing an illegal memory access is caught and rejected cleanly
- Host-injected functions are callable from within the WASM sandbox

---

## Layer 2 — The WASM Byte Builder (Go)

**Build this second. It sits between your language and the engine.**

This is a pure Go library that accepts an AST and emits a valid, binary-encoded
`.wasm` file. It does no parsing. It does no type checking. It knows nothing about
your language's syntax. It knows only the WASM binary specification and how to
write it to bytes.

### What it does

The builder exposes a structured API for programmatically assembling the sections
of a WASM binary. Your frontend calls this API — it never writes raw bytes or
thinks about section offsets or encoding rules.

A `.wasm` file is a sequence of typed sections. The builder is responsible for
constructing and encoding each one correctly:

- `type` — function signature definitions
- `import` — host API bindings declared to the engine
- `function` — index to signature mapping
- `memory` — linear memory size declaration
- `export` — symbols the engine can call into
- `code` — the actual instruction bodies, one per function
- `data` — static bytes for string literals and constants

### A rough sketch of the API

```go
m := wasm.NewModule()

sig := m.Types.Add(wasm.Params(wasm.I32, wasm.I32), wasm.Results(wasm.I32))

fn := m.Code.NewFunction(sig)
fn.LocalGet(0)
fn.LocalGet(1)
fn.I32Add()
fn.End()

m.Exports.Add("add", fn)

binary, err := m.Encode()
```

The output of `Encode()` is a byte slice that Layer 1 can load and run directly.

### Why Go

Go's goroutines make parallel function-body compilation trivial. Its strong typing
catches malformed instruction sequences at build time rather than producing a broken
binary silently. The emit phase is fast enough to feel instantaneous. And the layer
is easy to unit test in isolation — given a known AST node, assert the exact byte
sequence produced.

### Milestone — Layer 2 is complete when

- The builder emits a `.wasm` file that Layer 1 loads and runs correctly
- All core numeric types and instructions work (`i32`, `i64`, `f32`, `f64`)
- Function calls, locals, and control flow (`if`, `block`, `loop`, `br`) are supported
- The import section correctly wires through to Layer 1 host bindings

---

## Layer 3 — Your Custom Frontend

**Build this last. This is the only layer your users ever see.**

This layer is entirely about your language's identity — its syntax, its type system,
its rules and guarantees. It produces no output of its own. Its job is to take
source text, understand it fully, and hand a clean, verified AST to the Layer 2
builder.

### What it does

**Lexer** converts raw source text into a flat stream of tokens. Line and column
positions are tracked here. This is where error reporting starts.

**Parser** consumes the token stream and builds the AST. It encodes the grammatical
rules of your language — operator precedence, expression vs statement boundaries,
block structure.

**Semantic Analyzer** is the contract with your users. If your language promises
type safety, prove it here. If it promises a particular inference algorithm or
memory ownership model, enforce it here. By the time the AST reaches the builder,
every node must be fully resolved and verified. The builder never second-guesses
the frontend.

### Milestone — Layer 3 is complete when

- A source file tokenizes correctly and error positions are accurate
- The parser produces a well-formed AST for all valid syntax
- The semantic analyzer rejects type errors with clear, located messages
- A full source file compiles end-to-end: source → `.wasm` → runs on Layer 1

---

## The Developer Workflow

Once all three layers exist, the user sees none of this complexity.

```bash
# Compile — runs Layer 3 then Layer 2, produces a .wasm file
$ lang build main.lang -o program.wasm

# Run — runs Layer 1, AOT compiles and executes natively
$ lang run program.wasm
```

One source file in. One native execution out.

---

## Why This Build Order

Building Layer 1 first means you validate the runtime against known-good `.wasm`
binaries before writing a single line of your own language. You know your engine
works before you trust it.

Building Layer 2 as a pure library with no parser dependency means any frontend
can be swapped onto it. You can prototype multiple syntax ideas without touching
the bytecode layer.

Building Layer 3 last means your frontend only ever calls a stable, tested API.
When something goes wrong end-to-end, the layer boundaries tell you exactly where
to look — a bad `.wasm` points to Layer 2, a bad AST points to Layer 3, a crash
at runtime points to Layer 1.

Porting your language to a new platform means porting one C++ project. Everything
above it comes for free.
