# Vertex

**Language Spec 2.2 · Compiler 0.4.0**

Vertex is a statically-typed systems and application programming language built
for explicit control, zero-overhead C interop, and first-class concurrency.
Its syntax draws from Swift and Go, while its call-site execution sigils and
unified pointer system make it uniquely suited to systems work — from GPU
kernels and AI inference to networking, file I/O, and low-level kernel
programming.

---

## Contents

- [Install](#install)
- [Quick Start](#quick-start)
- [Language Tour](#language-tour)
  - [Variables and Types](#variables-and-types)
  - [Functions and Pointers](#functions-and-pointers)
  - [Structs and Classes](#structs-and-classes)
  - [Arrays and Maps](#arrays-and-maps)
  - [Concurrency](#concurrency)
  - [Error Handling](#error-handling)
  - [Native Interface](#native-interface)
  - [Testing](#testing)
- [Compiler Reference](#compiler-reference)
- [Platform Support](#platform-support)

---

## Install

Requires Go 1.23 or later.

```sh
GOPROXY=direct go install github.com/vertex-language/vertex@latest
```

Verify:

```sh
vertex -version
# vertex 0.4.0
```

---

## Quick Start

```vertex
package main
build linux

import "linux/lib/c"

class C : c {
    func printf(fmt: ...*const char) -> int32
}

func fibRecursive(n: int32) -> int32 {
    if n <= 1 {
        return n
    }
    return fibRecursive(n - 1) + fibRecursive(n - 2)
}

func fibIterative(n: int32) -> int32 {
    if n <= 1 {
        return n
    }
    var a: int32 = 0
    var b: int32 = 1
    var i: int32 = 2
    while true {
        if i > n { break }
        var tmp: int32 = a + b
        a = b
        b = tmp
        i = i + 1
    }
    return b
}

func main() -> int32 {
    var libc = C()

    libc.printf("--- Recursive ---\n")
    var i: int32 = 0
    while true {
        if i >= 10 { break }
        libc.printf("fib(%d) = %d\n", i, fibRecursive(i))
        i = i + 1
    }

    libc.printf("--- Iterative ---\n")
    var j: int32 = 0
    while true {
        if j >= 10 { break }
        libc.printf("fib(%d) = %d\n", j, fibIterative(j))
        j = j + 1
    }

    return 0
}
```

```sh
vertex -lc -o fib fib.vs
./fib
```

---

## Language Tour

### Variables and Types

`let` declares an immutable binding; `var` declares a mutable one. All numeric
conversions are explicit — there is no implicit coercion.

```vertex
let x: int32   = 42
var y: float32 = float32(x)
let name: string = "vertex"
let flag: bool   = true

let banner: string = `
  Vertex 2.2
  systems · concurrency · zero-overhead interop
`
```

Scalar types map directly to C:

| Vertex            | C type     |
|-------------------|------------|
| `int` / `int32`   | `int32_t`  |
| `int64`           | `int64_t`  |
| `uint` / `uint32` | `uint32_t` |
| `uint8`           | `uint8_t`  |
| `float32`         | `float`    |
| `float64`         | `double`   |
| `bool`            | `bool`     |
| `string`          | see docs   |

---

### Functions and Pointers

```vertex
func add(a: int32, b: int32) -> int32 {
    return a + b
}

func increment(n: *int32) {
    n += 1    // auto-dereferenced
}

var count = 0
increment(n: &count)   // count is now 1
```

`let`/`var` and `*const` are orthogonal: `let` locks the binding; `*const` locks
the pointed-to data. All four combinations are valid.

---

### Structs and Classes

**Structs** are value types — stack-allocated, always copied on assignment.

```vertex
struct Vec2 {
    x: float32
    y: float32
}

func (v: *Vec2) scale(factor: float32) {
    v.x *= factor
    v.y *= factor
}

var pos = Vec2{x: 1.0, y: 2.0}
pos.scale(factor: 2.0)
```

**Classes** are heap-allocated. The programmer controls lifetime explicitly via
`.delete()`, or opts into reference counting with `.new()`.

```vertex
class Animal {
    name:   string
    health: int32
}

func (a: Animal) init(name: string, health: int32) {
    a.name   = name
    a.health = health
}

func (a: Animal) deinit() { }

// manual lifetime
let dog = Animal(name: "Rex", health: 100)
defer dog.delete()

// reference counted
let cat = Animal(name: "Luna", health: 100).new()
weak let observer = cat    // Animal? — non-owning
```

---

### Arrays and Maps

**Fixed arrays** are stack-allocated; the size is part of the type.

```vertex
var buf:  [uint8; 1024]
var mask: [uint8; 64]
mask.fill(0xFF)

let coords: [int32; 3] = [10, 20, 30]
```

**Dynamic arrays** are heap-allocated and growable.

```vertex
var items: [int32] = []
defer items.delete()

items.push(10)
items.push(20)

var doubled = items.map(func(x: int32) -> int32 { return x * 2 })
defer doubled.delete()
```

**Maps** use brace literals. A type annotation is required for empty maps.

```vertex
let scores = {"alice": 42, "bob": 7}

var config: map[string]int32 = {}
defer config.delete()

config["workers"] = 4
config["verbose"] = nil    // removes key

let w = config["workers"]  // int32? — nil if absent
```

---

### Concurrency

Vertex decouples business logic from execution strategy. Every function is
written as an ordinary synchronous function — the caller decides how a call runs
by prefixing it with an execution sigil: `async`, `thread`, or `gpu`.

```vertex
// async — cooperative virtual thread, zero OS threads spawned
let pending = async fetchUser(id: 1)
let user    = pending.receive()

// thread — real OS thread, shared memory
let result = thread crunch(data: input)
let output = result.receive()

// gpu — PTX / SPIR-V kernel
let out = gpu(blocks: 16, threads: 256) vectorAdd(a: x, b: y)
let ans = out.receive()
```

**Channels** carry values across execution boundaries.

```vertex
let stream: chan float32 = {cap: 64}

thread func(data: [float32], ch: chan float32) {
    for chunk in data {
        ch.send(process(chunk))
    }
    ch.close()
}(dataset, stream)

while let val = stream.tryReceive() {
    print(val)
}
```

**`select`** suspends with 0% CPU until a channel is ready.

```vertex
select {
case a = task1.receive():
    print("task 1 done")
case b = task2.receive():
    print("task 2 done")
default:
    print("neither ready yet")
}
```

**Reactive state** broadcasts a value to all subscribers.

```vertex
struct AppState {
    count: int32
    done:  bool
}

let app: state AppState = { AppState{count: 0, done: false} }

thread func(st: state AppState) {
    st.set(AppState{count: 1, done: true})
}(app)

async func(s: state AppState) {
    if s.get().done { runtime.exit(0) }
}(app)

runtime.loop()
```

---

### Error Handling

Three primitives — choose based on what the caller needs:

```vertex
// T? — value may simply not exist
func findUser(id: int32) -> User? {
    if id < 0 { return nil }
    return User(id)
}

let u = findUser(id: -1) ?? defaultUser

// (T, E?) tuple — value and error together
func divide(a: int32, b: int32) -> (int32, string?) {
    if b == 0 { return (0, "division by zero") }
    return (a / b, nil)
}

let (result, err) = divide(a: 10, b: 0)

// Result(T, E) — explicit Ok / Err
func parseInt(s: string) -> Result(int32, string) {
    if s == "" { return Result(Err, "empty string") }
    return Result(Ok, 42)
}

switch parseInt(s: input) {
case Ok(let value):  // use value
case Err(let msg):   // handle error
}
```

---

### Native Interface

Every foreign target is expressed through an import path, a class declaration,
and a package.

| Import prefix | Strategy                             |
|---------------|--------------------------------------|
| `lib/`        | linked call (validated at link time)  |
| `linux/`      | inline syscall instruction            |
| `darwin/`     | `objc_msgSend` / selector dispatch    |
| `windows/`    | COM vtable slot dispatch              |
| `gpu/`        | PTX / shader kernel emission          |

```vertex
package libc
build linux
import "lib/c"

class C : c {
    func fopen(path: *const char, mode: *const char) -> *void?
    func fwrite(ptr: *const void, size: uint64, count: uint64, stream: *void) -> uint64
    func fclose(stream: *void) -> int32
    func printf(fmt: ...*const char) -> int32
}
```

Native class instances are zero-size — the backend removes them entirely. No
allocation, no runtime overhead.

---

### Testing

Test functions carry the `test` qualifier and are auto-discovered by the test
runner. Declare them in files tagged `build test`.

```vertex
package arithmetic_test
build test

import "arithmetic"

func test_add()        test -> Expected(int32, "15") { return add(a: 10, b: 5) }
func test_comparison() test -> Expected(bool, "1")   { return 5 > 3 }
func test_no_crash()   test                          { add(a: 0, b: 0) }
```

`Expected(type, string)` declares the return type and the exact stdout string
the runner checks. Omitting `Expected` means the test passes as long as it does
not crash.

**Return value format reference:**

| Type      | Format | `Expected` for value `5`        |
|-----------|--------|---------------------------------|
| `int32`   | `%d`   | `Expected(int32, "5")`          |
| `int64`   | `%lld` | `Expected(int64, "5")`          |
| `uint32`  | `%u`   | `Expected(uint32, "5")`         |
| `float32` | `%f`   | `Expected(float32, "5.000000")` |
| `bool`    | `%d`   | `Expected(bool, "1")` / `"0"`   |
| `string`  | `%s`   | `Expected(string, "hello")`     |

---

## Compiler Reference

The compiler transforms `.vs` source through a four-stage pipeline:

```
.vs source → AST → Vertex IR (.vir / .vbytes) → Machine IR (.mir) → native code
```

Each intermediate form can be captured independently with an emit flag,
which is useful for debugging, tooling, and build caches.

```
Usage:
  vertex [flags] <source.vs | package/>

Emit mode (exactly one required):
  -emit-vir         emit Vertex IR text (.vir)
  -emit-vbytes      emit Vertex IR binary (.vbytes)
  -emit-mir         emit Machine IR text (.mir)
  -emit-asm         emit native assembly text (.s)
  -emit-obj, -c     emit relocatable object file (.o / .obj)
  -lc               compile and link to native executable

Options:
  -o <file>         output file (default: derived from input name)
  -target <triple>  target triple — see Platform Support (default: host)
  -packages-dir     Vertex packages root (overrides $VERTEX_PATH)
  -O0               disable optimisation (default)
  -O1               light optimisation
  -O2               full optimisation
  -Os               optimise for size
  -g                include debug information
  -v, -version      print version and exit
```

**Examples:**

```sh
vertex -emit-vir    -o main.vir    main.vs
vertex -emit-vbytes -o main.vbytes main.vs
vertex -emit-mir    -o main.mir    main.vs
vertex -emit-asm    -o main.s      main.vs
vertex -c           -o main.o      main.vs
vertex -lc          -o main        main.vs
vertex -lc -target darwin-arm64 -O2 -o main main.vs
```

Output extension and name are derived from the input automatically when `-o` is
omitted. On Windows targets, object files use `.obj` and executables gain `.exe`.

**Intermediate formats:**

| Flag            | Output      | Use case                                    |
|-----------------|-------------|---------------------------------------------|
| `-emit-vir`     | `.vir`      | Human-readable Vertex IR; inspect lowering  |
| `-emit-vbytes`  | `.vbytes`   | Binary Vertex IR; incremental build cache   |
| `-emit-mir`     | `.mir`      | SSA Machine IR; inspect register allocation |
| `-emit-asm`     | `.s`        | Native assembly; inspect code generation    |
| `-c`/`-emit-obj`| `.o`/`.obj` | Relocatable object; link separately         |
| `-lc`           | executable  | Fully linked native binary                  |

The `$VERTEX_PATH` environment variable sets the packages root; `-packages-dir`
overrides it. When neither is set, the compiler defaults to
`~/.vertex/packages`.

---

## Platform Support

| Target                 | Object file | Executable | Assembly |
|------------------------|-------------|------------|----------|
| `linux-amd64`          | yes         | yes        | yes      |
| `linux-arm64`          | yes         | yes        | yes      |
| `linux-riscv64`        | yes         | yes        | yes (`-emit-asm` only — `-c`/`-lc` not yet supported) |
| `darwin-amd64`         | yes         | yes        | yes      |
| `darwin-arm64`         | yes         | yes        | yes      |
| `windows-amd64`        | yes         | yes        | yes      |
| `windows-arm64`        | yes         | yes        | yes      |
| `freestanding-amd64`   | yes         | —          | yes      |
| `freestanding-arm64`   | yes         | —          | yes      |
| `freestanding-riscv64` | yes         | —          | yes (`-emit-asm` only) |

Freestanding targets produce object files only; `-lc` is not supported.

Upcoming targets: `browser/wasm`, `android`, `browser/js`.

---

## Documentation

- [Grammar Specification 2.2](https://github.com/vertex-language/specs/README.md)

---

## License

MIT — see [LICENSE](LICENSE).