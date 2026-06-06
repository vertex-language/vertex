# Vertex Programming Language

**Version 2.1** · [Grammar Spec](docs/v2/grammar.md) · [Native Interface](docs/v2/native_interface.md)

Vertex is a statically-typed systems and application programming language built
for explicit control, zero-overhead C interop, and first-class concurrency.
Its syntax draws from Swift and Go, while its postfix execution model and
unified pointer system make it uniquely suited to systems work — from GPU
kernels and AI inference to networking, file I/O, and low-level kernel
programming.

Vertex compiles to native binaries on Linux, Windows, and Darwin.
Upcoming targets: `browser/wasm`, `android`, and `browser/js`.

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
- [Compiler](#compiler)
- [Platform Support](#platform-support)

---

## Install

Requires Go 1.21 or later.

```sh
go install github.com/vertex-language/vertex/cmd/vertex@latest
```

Verify:

```sh
vertex -version
# Vertex compiler 0.2.0
```

---

## Quick Start

```vertex
package main
build linux

import "linux/lib/c"

class C : c {
    func printf(fmt:  ...*const char) -> int
}

// fibRecursive computes the nth Fibonacci number using classic recursion.
func fibRecursive(n: int32) -> int32 {
    if n <= 1 {
        return n
    }
    return fibRecursive(n - 1) + fibRecursive(n - 2)
}

// fibIterative computes the nth Fibonacci number in O(n) time and O(1) space.
func fibIterative(n: int32) -> int32 {
    if n <= 1 {
        return n
    }
    var a: int32 = 0
    var b: int32 = 1
    var i: int32 = 2
    while true {
        if i > n {
            break
        }
        var tmp: int32 = a + b
        a = b
        b = tmp
        i = i + 1
    }
    return b
}

func main() -> int {

    var libc = C()

    libc.printf("--- Recursive ---\n")
    var i: int32 = 0
    while true {
        if i >= 10 {
            break
        }
        libc.printf("fib(%d) = %d\n", i, fibRecursive(i))
        i = i + 1
    }

    libc.printf("--- Iterative ---\n")
    var j: int32 = 0
    while true {
        if j >= 10 {
            break
        }
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

`let` declares an immutable binding; `var` declares a mutable one.
All numeric conversions are explicit — there is no implicit coercion.

```vertex
let x: int32  = 42
var y: float32  = float32(x)   // explicit int → float32
let name: string = "vertex"
let flag: bool = true

// multiline string
let banner: string = `
  Vertex 2.1
  systems · concurrency · zero-overhead interop
`
```

Scalar types map directly to C:

| Vertex          | C type     |
|-----------------|------------|
| `int` / `int32` | `int32_t`  |
| `int64`         | `int64_t`  |
| `uint8`         | `uint8_t`  |
| `float32`       | `float`    |
| `float64`       | `double`   |
| `bool`          | `bool`     |
| `string`        | see docs   |

---

### Functions and Pointers

```vertex
func add(a: int32, b: int32) -> int32 {
    return a + b
}

// pointer parameter — mutations affect the caller
func increment(n: *int32) {
    n += 1    // auto-dereferenced
}

var count = 0
increment(n: &count)   // count is now 1
```

`let`/`var` and `*const` are orthogonal: `let` locks the binding;
`*const` locks the data. All four combinations are valid.

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
pos.scale(factor: 2.0)    // pos.x == 2.0, pos.y == 4.0
```

**Classes** are heap-allocated. The programmer controls lifetime explicitly
via `.delete()`, or opts into reference counting with `.new()`.

```vertex
class Animal {
    name: string
    health: int32
}

func (a: *Animal) init(name: string, health: int32) {
    a.name   = name
    a.health = health
}

func (a: *Animal) deinit() {
    // cleanup before free
}

// manual lifetime
let dog = Animal(name: "Rex", health: 100)
defer dog.delete()

// reference counted
let cat = Animal(name: "Luna", health: 100).new()
weak let observer = cat    // observer: Animal? — non-owning
```

---

### Arrays and Maps

**Fixed arrays** are stack-allocated:

```vertex
var buf  = [uint8](1024)             // 1024 zero bytes, stack
var mask = [uint8](repeating: 0xFF, count: 64)
```

**Growable arrays** are heap-allocated and must be freed:

```vertex
var items = [int32]()
defer items.delete()

items.push(10)
items.push(20)

var doubled = items.map(func(x: int32) -> int32 { return x * 2 })
defer doubled.delete()
```

**Maps:**

```vertex
var config = map[string]int32()
defer config.delete()

config["workers"] = 4
config["verbose"] = 1
config["verbose"] = nil    // removes key

let w = config["workers"]  // int32? — nil if absent
```

---

### Concurrency

Vertex exposes four execution models through a unified postfix syntax.
The qualifier sits between the parameter list and the return arrow.

**Async / Await:**

```vertex
func fetchUser(id: int32) async -> User {
    // non-blocking I/O
}

let user = fetchUser(id: 1).await()
```

**Threads** — shared memory, lightweight:

```vertex
func crunch(data: [float32]) thread -> [float32] {
    // parallel compute over shared memory
}

let result = crunch(data: input).spawn(threads: 4)
```

**Processes** — fully isolated memory:

```vertex
func isolatedWork(data: [float32]) process -> [float32] {
    // separate address space
}

let result = isolatedWork(data: input).fork(processes: 4)
```

**GPU Kernels:**

```vertex
func vectorAdd(a: [float32], b: [float32]) gpu -> [float32] {
    // emits PTX / shader IR
}

let output = vectorAdd(a: x, b: y).dispatch()
```

**Channels** wire concurrent functions together:

```vertex
let ch = float32.channel(size: 256)

func(data: [float32], out: chan float32) thread {
    for chunk in data {
        out.send(process(chunk))
    }
    out.close()
}(dataset, ch).spawn()

while true {
    let val = ch.tryReceive() ?? break
    // consume val
}
```

---

### Error Handling

Three primitives — choose based on what the caller needs to know:

```vertex
// T? — value may simply not exist
func findUser(id: int32) -> User? {
    if id < 0 { return nil }
    return User(id)
}

if let user = findUser(id: 42) { }
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
case Ok(let value):
    // use value
case Err(let msg):
    // handle error
}
```

---

### Native Interface

Every foreign target — C libraries, OS syscalls, COM, CUDA, bare metal
interrupts — is expressed through the same three-part pattern:
an **import path**, a **class declaration**, and a **package**.

The import prefix is the single source of truth for how the compiler emits:

| Prefix       | Strategy                              |
|--------------|---------------------------------------|
| `lib/`       | linked call (validated at link time)  |
| `linux/`     | inline syscall instruction            |
| `darwin/`    | `objc_msgSend` / selector dispatch    |
| `windows/`   | COM vtable slot dispatch              |
| `gpu/`       | PTX / shader kernel emission          |
| `metal/`     | hardware interrupt instruction        |

**C standard library:**

```vertex
package libc
build linux
import "lib/c"

class C : c {
    func fopen(path: *const char, mode: *const char) -> *void?
    func fwrite(ptr: *const void, size: uint64, count: uint64, stream: *void) -> uint64
    func fclose(stream: *void) -> int32
    func printf(fmt: *const char, ...) -> int32
}
```

```vertex
var c = libc.C()
let f = c.fopen("/tmp/out.bin", "wb")
defer c.fclose(f)

var data = [uint8](repeating: 0xFF, count: 256)
c.fwrite(&data, 1, data.byteSize(), f)
c.printf("wrote %d bytes\n", data.length)
```

**Linux syscalls:**

```vertex
package syscall
build linux
import "linux/syscalls"

class Syscalls : syscalls {
    func write(fd: int32, buf: *const void, count: uint64) -> int32
    func read(fd: int32, buf: *void, count: uint64) -> int32
}
```

Native class instances are zero-size — the backend removes them entirely.
No allocation, no runtime overhead.

---

### Testing

Test functions use the `test` qualifier and are auto-discovered by the
test runner. Declare them in files tagged `build test`.

```vertex
package arithmetic_test
build test

import "arithmetic"

func test_add()        test -> Expected("15")   { return add(a: 10, b: 5) }
func test_comparison() test -> Expected("true") { return 5 > 3 }
func test_no_crash()   test                     { add(a: 0, b: 0) }
```

```sh
vertex -test arithmetic_test.vs        # run tests in one file
vertex -test -dir .                    # run all tests recursively
```

A test passes when its return value — auto-formatted to stdout — matches
the `Expected` string. Omitting `Expected` means the test passes if it
completes without crashing.

---

## Compiler

```
Usage:
  vertex [flags] <source.vs>

Flags:
  -o file           write output to file
  -target string    linux-amd64 (default), darwin-amd64, windows-amd64
  -emit-c           emit C source instead of native binary
  -c                compile and assemble, do not link (outputs .o)
  -I path           add a package search path (repeatable)
  -L dir            add a library search dir — ELF targets only (repeatable)
  -l name           link against libname e.g. -lc -lm (repeatable)
  -test             compile and run build-test functions
  -dir directory    run tests recursively in directory (with -test)
  -version          print version and exit

Examples:
  vertex -o main main.vs                    # native executable
  vertex -lc -lm -o main main.vs            # link against libc and libm
  vertex -c -o main.o main.vs               # object file only
  vertex -emit-c -o main.c main.vs          # emit C source
  vertex -test arithmetic_test.vs           # run tests in file
  vertex -test -dir .                       # run all tests recursively
  vertex -target darwin-amd64 -o app main.vs
```

---

## Platform Support

| Target           | Status     |
|------------------|------------|
| `linux-amd64`    | Supported  |
| `darwin-amd64`   | Supported  |
| `windows-amd64`  | Supported  |
| `browser/wasm`   | Upcoming   |
| `android`        | Upcoming   |
| `browser/js`     | Upcoming   |

---

## Documentation

- [Grammar Specification 2.1](docs/grammar.md)
- [Native Interface](docs/native_interface.md)

---

## License

See [LICENSE](LICENSE).