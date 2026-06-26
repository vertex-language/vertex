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
  - [Generics](#generics)
  - [Structs and Classes](#structs-and-classes)
  - [Arrays and Maps](#arrays-and-maps)
  - [Tuples](#tuples)
  - [Enums](#enums)
  - [Error Handling](#error-handling)
  - [Concurrency](#concurrency)
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
vertex -o fib fib.vs
./fib
```

---

## Language Tour

### Variables and Types

`let` declares an immutable binding; `var` declares a mutable one. All numeric
conversions are explicit — there is no implicit coercion between numeric types.

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

Multiline strings are delimited by backticks. No indentation stripping is applied.

Scalar types map directly to C:

| Vertex | C type |
| --- | --- |
| `int` / `int32` | `int32_t` |
| `int8` | `int8_t` |
| `int16` | `int16_t` |
| `int64` | `int64_t` |
| `uint` / `uint32` | `uint32_t` |
| `uint8` | `uint8_t` |
| `uint16` | `uint16_t` |
| `uint64` | `uint64_t` |
| `float32` | `float` |
| `float64` | `double` |
| `bool` | `bool` |
| `char` | `char` |
| `string` | `let` → rodata, `var` → heap |

`int` is an alias for `int32`; `uint` is an alias for `uint32`.

The `as` operator performs explicit type conversion for numeric widening,
pointer reinterpretation, and float-to-integer truncation:

```vertex
let small: int32 = 42
let wide  = small as int64          // integer widening — sign-extended
let big   = small as uint64         // zero-extended

let f: float64 = 3.99
let i = f as int32                  // truncates toward zero → 3

var buf: [uint8; 256]
libc.recv(fd, &buf as *char, 256, 0)  // pointer reinterpret

// chaining — left-associative
let x = value as int32 as int64
```

Conversion syntax `targetType(value)` is also valid for narrowing, widening, and
float-to-integer truncation:

```vertex
let b: int8 = int8(i)        // narrowing — wraps on overflow
let f2: float32 = float32(i) // int → float32, always safe
let i2: int   = int(3.99)    // truncates toward zero → 3
```

---

### Functions and Pointers

```vertex
func add(a: int32, b: int32) -> int32 {
    return a + b
}

// call with positional or labeled arguments
add(1, 2)
add(a: 1, b: 2)
```

Pointer parameters let a function mutate the caller's value. Reads and writes
through a pointer parameter are auto-dereferenced by the compiler:

```vertex
func increment(n: *int32) {
    n += 1    // auto-dereferenced — lowers to *n += 1
}

var count = 0
increment(n: &count)   // count is now 1
```

`let`/`var` and `*const` are orthogonal. `let` locks the binding; `*const` locks
the pointed-to data. All four combinations are valid:

| Vertex               | C                | Binding   | Data      |
|----------------------|------------------|-----------|-----------|
| `let name: *const T` | `const T* const` | fixed     | read-only |
| `let name: *T`       | `T* const`       | fixed     | mutable   |
| `var name: *const T` | `const T*`       | rebind OK | read-only |
| `var name: *T`       | `T*`             | rebind OK | mutable   |

Every function is written as an ordinary synchronous function. The caller decides
how a call runs by prefixing it with an execution sigil at the call site —
`async`, `thread`, or `gpu`. There are no `async` or `thread` function qualifiers.

---

### Generics

Unconstrained generic functions and structs use angle-bracket type parameters:

```vertex
func identity<T>(value: T) -> T {
    return value
}

struct Box<T> {
    value: T
}

let b      = Box<T>{value: 42}
let result = identity<T>(value: "hello")
```

---

### Structs and Classes

**Structs** are value types — stack-allocated, always copied on assignment.
Mutability is determined entirely by the binding at the declaration site.

```vertex
struct Vec2 {
    x: float32
    y: float32
}

// value receiver — copy; mutations do not affect the caller
func (v: Vec2) describe() {
    // ...
}

// pointer receiver — mutations affect the caller's binding
func (v: *Vec2) scale(factor: float32) {
    v.x *= factor   // auto-dereferenced — lowers to v->x
    v.y *= factor
}

var pos = Vec2{x: 1.0, y: 2.0}
pos.scale(factor: 2.0)   // compiler inserts & automatically
```

Struct literals require all field labels; positional initialization is not
supported. Trailing commas are valid in multiline forms:

```vertex
let p = Vec2{
    x: 3.0,
    y: 4.0,
}
```

**Classes** are heap-allocated references. The programmer controls lifetime
explicitly via `.delete()`, or opts into reference counting with `.new()`.

```vertex
class Animal {
    name:   string
    health: int32
}

func (a: Animal) init(name: string, health: int32) {
    a.name   = name
    a.health = health
}

func (a: Animal) deinit() {
    // runs before memory is freed
}

// manual lifetime
let dog = Animal(name: "Rex", health: 100)
defer dog.delete()

// reference counted
let cat = Animal(name: "Luna", health: 100).new()
weak let observer = cat    // Animal? — non-owning, count stays 1

if let animal = observer {
    // safe — animal is Animal within this scope
}
```

`.new()` and `.delete()` are mutually exclusive. After all owning references
reach zero, all `weak` references become `nil`.

The `auto` binding modifier instructs the compiler to automatically inject
`.delete()` at scope exit, eliminating `defer` boilerplate for the common case.
Cleanup order follows LIFO — identical to explicit `defer`:

```vertex
// before
let log = Logger(path: "job.log")
defer log.delete()
let buf = Buffer(capacity: 4096)
defer buf.delete()

// after — identical semantics, same LIFO teardown
auto let log = Logger(path: "job.log")
auto let buf = Buffer(capacity: 4096)
```

`auto` is only valid on class bindings. Drop it when you need explicit control —
early release, conditional cleanup, or ownership transfer.

| Binding | Lifetime |
|---|---|
| `let` / `var` | manual — you call `.delete()` |
| `auto let` / `auto var` | automatic — compiler injects `.delete()` at scope exit |
| `weak let` | non-owning — no cleanup, never owned |

---

### Arrays and Maps

**Fixed arrays** are stack-allocated; the size is part of the type.
`push`, `pop`, `shift`, `unshift`, `.reserve()`, and `.delete()` are compile
errors on any fixed array.

```vertex
var buf:  [uint8; 1024]           // zero-filled, no initializer needed
var mask: [uint8; 64]
mask.fill(0xFF)

var coords: [int32; 3] = [10, 20, 30]
let flags:  [uint8; 3] = [0xFF, 0x00, 0xAB]

// nested (multidimensional)
let matrix: [[float32; 2]; 2] = [
    [0.0, 1.0],
    [1.0, 0.0],
]
```

**Dynamic arrays** are heap-allocated and growable. `var` with no size annotation
(`[T]`) always produces a dynamic array.

```vertex
var items: [int32] = []
defer items.delete()

items.push(42)           // add to end
items.unshift(0)         // add to front
let last  = items.pop()  // remove from end  — returns T?
let first = items.shift()// remove from front — returns T?

items.reserve(64)        // pre-allocate capacity
```

Methods that return a new array allocate on the heap — call `.delete()` on the
result:

```vertex
var doubled = items.map(func(x: int32) -> int32 { return x * 2 })
defer doubled.delete()

var evens = items.filter(func(x: int32) -> bool { return x % 2 == 0 })
defer evens.delete()

var sub = items.slice(1, 3)
defer sub.delete()
```

In-place mutation methods (`sort`, `reverse`, `fill`, `reserve`) do not
allocate — no `.delete()` needed on the result.

The rule at a glance:

| Form | Storage | Growable |
|------|---------|----------|
| `var buf: [T; N]` | stack | no |
| `let arr = [...]` | stack / rodata | no |
| `var x: [T] = []` | heap | yes |
| `var x = [...]` | heap | yes |

**Maps** use brace literals. A type annotation is required for empty maps.
Key access always returns an optional (`T?`). Assigning `nil` to a key removes it.

```vertex
let scores = {"alice": 42, "bob": 7}
let val = scores["alice"]          // int32?

var config: map[string]int32 = {}
defer config.delete()

config["workers"] = 4
config["verbose"] = nil    // removes key
```

---

### Tuples

Tuples are stack-allocated value types for multi-value returns and paired data —
zero overhead, no heap allocation, no named struct required.

```vertex
let pair  = (1, true)
let point = (x: 10, y: 20)

// unlabeled return
func divmod(a: int32, b: int32) -> (int32, int32) {
    return (a / b, a % b)
}
let (quotient, remainder) = divmod(10, 3)

// labeled return
func minMax(values: [int32]) -> (min: int32, max: int32) {
    return (0, 100)
}
let (lo, hi) = minMax(values: [3, 1, 4])
```

`()` is the empty tuple and is an alias for `void`. Tuples are value types —
assignment copies all elements. `==`, `!=`, `<`, `>`, `<=`, `>=` work on tuples
whose elements are all comparable, up to 6 elements.

Channels can carry tuples for paired data:

```vertex
let stream: chan (int32, bool) = {cap: 64}

select {
case (val, ok) = stream.receive():
    if ok { print(val) }
}
```

---

### Enums

Vertex enums support unit variants, tuple variants (positional associated data),
or a mix of both. The `case` keyword is used only inside `switch` statements —
not in enum declarations.

**Unit variants:**

```vertex
enum Direction {
    North,
    South,
    East,
    West,
}

let d: Direction = .South

switch d {
case .North: // ...
case .South: // ...
case .East:  // ...
case .West:  // ...
}
```

**Tuple variants:**

```vertex
enum Shape {
    Point,
    Circle(float32),
    Rectangle(float32, float32),
}

let s = Shape.Circle(1.5)

switch s {
case .Point:
    // no data
case .Circle(r):
    // r: float32
case .Rectangle(w, h):
    // w: float32, h: float32
}
```

Unwanted fields may be discarded with `_`:

```vertex
switch s {
case .Rectangle(w, _):
    // only care about width
default:
    // ...
}
```

When two or more positional fields share the same type and the distinction
matters, carry a named struct as the payload:

```vertex
struct MousePos { x: int32; y: int32 }

enum Event {
    Quit,
    KeyPress(uint8),
    MouseClick(MousePos),
}
```

**Explicit discriminants** require a backing integer type and are only valid on
all-unit enums. Auto-increment applies to unspecified variants:

```vertex
enum Status : int32 {
    Inactive = 0,
    Active   = 1,
    Pending  = 2,
}

enum HttpMethod : uint8 {
    Get = 0,
    Post,    // 1
    Put,     // 2
    Delete,  // 3
}

let raw = Status.Active as int32   // cast to backing type with `as`
```

Integer-to-enum conversion requires an explicit `switch` returning `EnumType?`:

```vertex
func statusFromInt(n: int32) -> Status? {
    switch n {
    case 0: return .Inactive
    case 1: return .Active
    case 2: return .Pending
    default: return nil
    }
}
```

Enums do not support `==`, `!=`, or any comparison operator — equality is
expressed via `switch`. Enums are value types; assignment copies.

---

### Error Handling

Vertex error handling is built on plain tuples and `?` propagation. A function
that can fail returns a tuple where the last element signals the outcome. The
zero value for its type signals success — `""` for string, `false` for bool,
`nil` for optionals, `0` for integers.

```vertex
// T? — value may simply not exist
func findUser(id: int32) -> User? {
    if id < 0 { return nil }
    return User(id)
}

let u = findUser(id: -1) ?? defaultUser
```

```vertex
// (T, E) tuple — value and error together
func divide(a: int32, b: int32) -> (int32, string) {
    if b == 0 { return (0, "division by zero") }
    return (a / b, "")
}

// 1 — plain destructuring
let (result, err) = divide(a: 10, b: 0)
if err != "" { /* handle */ }
```

```vertex
// 2 — ? propagation (only valid inside a function that itself returns a tuple)
func process(s: string) -> (int32, string) {
    let n = parseInt(s: s)?   // returns early if err != ""
    return (n * 2, "")
}
```

```vertex
// 3 — happy path only
if let n = parseInt(s: "42") { /* use n */ }

// 4 — both paths with else ->
if let n = parseInt(s: input) {
    // use n
} else -> err {
    // use err
}

// 5 — full control via switch
let (n, err) = parseInt(s: "42")
switch err {
case "":      // use n
default:      // use err
}
```

| Situation | Use |
|---|---|
| Value may simply not exist | `T?` |
| Handle it yourself | `let (val, err) = f()` |
| Bubble the error up | `?` |
| Happy path only | `if let` |
| Inspect both paths | `else ->` on `if let` |
| Full destructuring | `switch` |

---

### Concurrency

Vertex completely decouples business logic from execution strategy. Every
function is written as an ordinary synchronous function — the caller decides how
a call runs by prefixing it with an execution sigil at the call site:

| Sigil | Reality | Best for |
|---|---|---|
| `async` | Virtual thread — kilobyte-scale stack, userspace context-switch | Millions of idle connections, massive fan-out |
| `thread` | Real OS thread — full 2 MB+ stack | Heavy CPU work, blocking C calls |
| `gpu` | PTX / SPIR-V kernel | Matrix math, AI inference |

```vertex
// async — cooperative virtual thread
let pending = async fetchUser(id: 1)
let user    = pending.receive()

// thread — real OS thread
let result = thread crunch(data: input)
let output = result.receive()

// gpu — hardware kernel
let out = gpu(blocks: 16, threads: 256) vectorAdd(a: x, b: y)
let ans = out.receive()
```

If a function returns `-> T`, the compiler auto-channels the result through a
1-capacity channel (Path A). If a function returns void, the developer passes
channels explicitly for full stream control (Path B):

```vertex
// Path B — explicit stream
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

**Channel API:**

| Method | Blocking | Returns | Behaviour |
|---|---|---|---|
| `.send(value)` | yes | `void` | waits until value is accepted |
| `.receive()` | yes | `T` | waits until value is available |
| `.trySend(v)` | no | `bool` | false if full or no receiver ready |
| `.tryReceive()` | no | `T?` | nil if channel is empty |
| `.close()` | no | `void` | always completes immediately |

**`select`** multiplexes channels, suspending with 0% CPU until a case is ready.
Adding a `default` case makes it non-blocking:

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

**Reactive state** broadcasts a value to all subscribers. `state T` wraps any
value type in a pub/sub primitive — unlike channels, delivery is lossy-latest
and many subscribers may receive each update:

```vertex
struct AppState { count: int32; done: bool }

let app: state AppState = { AppState{count: 0, done: false} }

// thread broadcasts state changes
thread func(st: state AppState) {
    st.set(AppState{count: 1, done: true})
}(app)

// async effect — wakes on every broadcast
async func(s: state AppState) {
    if s.get().done { runtime.exit(0) }
}(app)

runtime.loop()
```

When an `async` function declares a `state T` parameter, the compiler
automatically generates the subscriber endpoint, the loop, and the receive
machinery. Multiple `state T` parameters subscribe to all of them simultaneously.

```
chan T    →  point-to-point, FIFO, one consumer per message
state T  →  broadcast (pub/sub), lossy-latest, many subscribers
```

---

### Native Interface

Every foreign target is expressed through an import path, a class declaration,
and a package:

| Import prefix | Strategy |
| --- | --- |
| `lib/` | linked call (validated at link time) |
| `dynamic/lib/` | runtime `dlopen` / `LoadLibrary` |
| `linux/` | inline syscall instruction |
| `darwin/` | `objc_msgSend` / selector dispatch |
| `windows/` | COM vtable slot dispatch |
| `gpu/` | PTX / shader kernel emission |

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

For libraries that resolve at runtime, prefix the import with `dynamic/lib/`.
All declared symbols are resolved eagerly at construction. Use a nullable binding
to handle absence gracefully:

```vertex
import "dynamic/lib/cuda"

class Cuda : cuda {
    func cuInit(flags: int32) -> int32
    func cuMemAlloc(dptr: *CUdevptr, size: int32) -> int32
    func cuMemFree(dptr: CUdevptr) -> int32
}

// nil if library not found or any symbol missing
var cuda: Cuda? = Cuda()
if let c = cuda {
    c.cuInit(0)
}
```

Individual functions resolve to `nil` when a symbol is absent from the loaded
library — useful for version compatibility:

```vertex
if cuda.cuMemAllocAsync != nil {
    cuda.cuMemAllocAsync(&ptr, size, stream)
} else {
    cuda.cuMemAlloc(&ptr, size)
}
```

---

### Testing

Test functions carry the `test` qualifier and are auto-discovered by the test
runner. Declare them in files tagged `build test`. The `test` qualifier sits
between the parameter list and `->`. Test functions may declare no parameters.

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

| Type | Format | `Expected` for value `5` |
| --- | --- | --- |
| `int32` | `%d` | `Expected(int32, "5")` |
| `int64` | `%lld` | `Expected(int64, "5")` |
| `uint32` | `%u` | `Expected(uint32, "5")` |
| `float32` | `%f` | `Expected(float32, "5.000000")` |
| `bool` | `%d` | `Expected(bool, "1")` / `"0"` |
| `string` | `%s` | `Expected(string, "hello")` |

---

## Compiler Reference

The compiler transforms `.vs` source through a four-stage pipeline:

```
.vs source → AST → Vertex IR (.vir / .vbytes) → Machine IR (.mir) → native code
```

```text
Usage:
  vertex [flags] <source.vs | package/>

Emit mode (default: compile and link to native executable):
  -emit-vir             emit Vertex IR text (.vir)
  -emit-vbytes          emit Vertex IR binary (.vbytes)
  -emit-mir             emit Machine IR text (.mir)
  -emit-asm             emit native assembly text (.s)
  -emit-obj, -c         emit relocatable object file (.o / .obj)
  -dump, -dump-all      dump all pipeline stages (.dump)
  -test                 discover and run test functions

Test options:
  -dir  <path>     directory to search recursively (default: .)
  -file <path>     single test file

Options:
  -o <file>        output file (default: derived from input)
  -target <triple> linux-amd64, linux-arm64, linux-riscv64,
                   darwin-amd64, darwin-arm64,
                   windows-amd64, windows-arm64,
                   freestanding-amd64, freestanding-arm64,
                   freestanding-riscv64  (default: host OS/Arch)
  -sysroot <path>  sysroot for cross-compilation library search
  -packages-dir    Vertex packages root (overrides $VERTEX_PATH)
  -O0/-O1/-O2/-Os  optimisation level (default: -O0)
  -g               include debug information
  -v, -version     print version and exit
```

**Examples:**

```sh
vertex -o main        main.vs
vertex -o main        -target darwin-arm64 -O2 main.vs
vertex -c           -o main.o      main.vs
vertex -emit-asm    -o main.s      main.vs
vertex -emit-mir    -o main.mir    main.vs
vertex -emit-vir    -o main.vir    main.vs
vertex -emit-vbytes -o main.vbytes main.vs
vertex -dump        -o main.dump   main.vs
vertex -dump        -o -           main.vs
vertex -test
vertex -test -dir ./tests
vertex -test -file literals_test.vs
```

**Modes & Formats:**

| Mode | Output | Use case |
| --- | --- | --- |
| *(default)* | executable | Fully linked native binary |
| `-emit-vir` | `.vir` | Human-readable Vertex IR; inspect lowering |
| `-emit-vbytes` | `.vbytes` | Binary Vertex IR; incremental build cache |
| `-emit-mir` | `.mir` | SSA Machine IR; inspect register allocation |
| `-emit-asm` | `.s` | Native assembly; inspect code generation |
| `-c` / `-emit-obj` | `.o` / `.obj` | Relocatable object; link separately |
| `-dump` | `.dump` | Output all annotated pipeline stages to one file |
| `-test` | *console* | Execute `test` functions directly |

The `$VERTEX_PATH` environment variable sets the packages root; `-packages-dir`
overrides it. When neither is set, the compiler defaults to `~/.vertex/packages`.

---

## Platform Support

| Target | Object file (`-c`) | Executable (default) | Assembly (`-emit-asm`) |
| --- | --- | --- | --- |
| `linux-amd64` | yes | yes | yes |
| `linux-arm64` | yes | yes | yes |
| `linux-riscv64` | yes | yes | yes (`-emit-asm` only — object/executable emission not yet supported) |
| `darwin-amd64` | yes | yes | yes |
| `darwin-arm64` | yes | yes | yes |
| `windows-amd64` | yes | yes | yes |
| `windows-arm64` | yes | yes | yes |
| `freestanding-amd64` | yes | — | yes |
| `freestanding-arm64` | yes | — | yes |
| `freestanding-riscv64` | yes | — | yes (`-emit-asm` only) |

Freestanding targets produce object files only; executable linking is not
supported.

Upcoming targets: `browser/wasm`, `android`, `browser/js`.

---

## Documentation

- [Grammar Specification 2.2](https://github.com/vertex-language/spec/README.md)

---

## License

MIT — see [LICENSE](https://github.com/vertex-language/vertex/LICENSE)