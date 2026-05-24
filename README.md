# Vertex Programming Language

### Version 1.9

Vertex is a statically-typed systems and application programming language built
for explicit control, zero-overhead C interop, and first-class concurrency.
Its syntax draws from Swift and Go, while its postfix execution model and
unified pointer system make it uniquely suited to systems work — from GPU
kernels and AI inference to networking, file I/O, and low-level kernel
programming.

Vertex compiles to native binaries and WebAssembly.

---

## Installation

Install the Vertex compiler via Go:

```bash
go install github.com/vertex-language/vertex/cmd/vertex@latest
```

### Compiler Usage

```bash
# Compile a single file
vertex main.vs -o myapp

# Compile an entire package directory
vertex -o myapp ./cmd/myapp

# Emit WebAssembly
vertex -wasm main.vs -o myapp.wasm

# Cross-compile
vertex -platform linux  -o myapp_linux  main.vs
vertex -platform darwin -o myapp_darwin main.vs

# Custom entry point
vertex -entry start -wasm -o lib.wasm lib.vs

# Disable dead-code elimination
vertex -no-dce -o debug main.vs
```

**All flags:**

| Flag        | Default    | Meaning                                        |
|-------------|------------|------------------------------------------------|
| `-o`        | `a.out`    | Output file path                               |
| `-wasm`     | off        | Emit `.wasm` module instead of a native binary |
| `-platform` | `linux`    | Target platform: `linux`, `darwin`, `windows`  |
| `-tags`     | —          | Comma-separated extra build tags               |
| `-module`   | —          | Module path prefix                             |
| `-root`     | source dir | Module root directory                          |
| `-entry`    | `main`     | Entry-point function name                      |
| `-no-dce`   | off        | Disable dead-code elimination                  |
| `-v`        | —          | Print version and exit                         |

---

## Language Tour

### Variables and Types

`let` declares an immutable binding. `var` declares a mutable one.
All numeric conversions are explicit — there is no implicit coercion.

```swift
let limit: int   = 100
let ratio: float = float(limit) / 3.0
var count        = 0

// Multiline strings use backticks
let message = `
Hello,
World!
`

// Numeric literals support underscores, binary, octal, and hex
let mask   = 0xFF_00_FF
let flags  = 0b1010_0011
let big    = 1_000_000
let hfloat = 0xFp2      // hex float = 60.0
```

**Primitive types:**

| Category | Types                                         |
|----------|-----------------------------------------------|
| Signed   | `int` `int8` `int16` `int32` `int64`          |
| Unsigned | `uint` `uint8` `uint16` `uint32` `uint64`     |
| Float    | `float` `double`                              |
| Other    | `bool` `string` `char` `void`                 |

**Numeric conversion** uses `targetType(value)` — no cast keyword:

```swift
let i: int   = 42
let f: float = float(i)    // int → float, always safe
let b: int8  = int8(i)     // narrowing — wraps on overflow
let t: int   = int(3.99)   // float → int, truncates toward zero
```

---

### Structs and Classes

**Structs** are stack-allocated value types. Assignment always copies.
**Classes** are heap-allocated reference types.

```swift
struct Point {
    let x: int
    var y: int
}

class Animal {
    var name: string
}
```

Struct literals use brace syntax — all field labels are required:

```swift
let p = Point{x: 3, y: 4}
var q = Point{x: 3, y: 4}
q.y   = 10

// Multiline — trailing comma valid
let p = Point{
    x: 3,
    y: 4,
}
```

Classes support manual lifetime management or opt-in reference counting:

```swift
// Manual — call .delete() explicitly; deinit fires automatically
let a = Animal(name: "Rex")
defer a.delete()

// Reference counted — freed when all owners go out of scope
let a = Animal(name: "Rex").new()

// Non-owning weak reference — does not increment the count
weak let b = a   // b: Animal?
if let animal = b {
    // safe — animal: Animal within this scope
}
```

`init` and `deinit` are reserved associated function names. `init` is called
automatically after allocation; `deinit` fires when `.delete()` is called or
when a ref-counted instance's count reaches zero.

```swift
func init(a: mut Animal, name: string) {
    a.name = name
}

func deinit(a: mut Animal) {
    // cleanup before memory is freed
}
```

---

### Associated Functions

Any function whose first parameter is a known struct or class is automatically
an associated function of that type — no `self`, no method block.

```swift
class YourClass {
    var x: int32
    var y: int32
}

func setX(s: mut YourClass, x: int32) { s.x = x }
func getX(s: YourClass) -> int32      { return s.x }

var m = YourClass(x: 0, y: 0).new()
m.setX(x: int32(10))
let x = m.getX()
```

Mutable receivers require `mut` on the parameter and `&` at the call site:

```swift
func reset(p: mut Point) {
    p.x = 0
    p.y = 0
}

var p = Point{x: 3, y: 4}
p.reset()
```

To write a utility function that takes a known type without it acting as the
receiver, place the known type at parameter index 1 or later.

---

### Enums

```swift
enum Direction { case north, south, east, west }

enum Status: int {
    case inactive = 0
    case active   = 1
    case pending  = 2
}

let raw = Status.active.rawValue          // 1
let s   = Status(rawValue: 1)             // Status?
```

Raw value types are `int` or `string`. String raw values default to the case
name if omitted. A switch over an enum with all cases covered needs no
`default`.

---

### Error Handling

Vertex has no exceptions. Errors are values.

**Optionals** — when a value may simply be absent:

```swift
func findUser(id: int) -> User? {
    if id < 0 { return nil }
    return User(id)
}

if let user = findUser(id: 1) { }
let name = findUser(id: -1) ?? defaultUser
```

**Result** — when the caller must handle `Ok` or `Err` explicitly:

```swift
func parseInt(s: string) -> Result(int, string) {
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

**`.try()`** — propagate a Result error up the call stack:

```swift
func process(s: string) -> Result(int, string) {
    let n = parseInt(s: s).try()
    let d = divide(a: n, b: 2).try()
    return Result(Ok, d)
}
```

**Choosing the right primitive:**

| Situation                               | Use              |
|-----------------------------------------|------------------|
| Value may simply not exist              | `T?`             |
| Caller needs value and error together   | `(T, E?)` tuple  |
| Caller must handle Ok or Err explicitly | `Result(T, E)`   |

---

### Control Flow

```swift
// if / else if / else
if x > 0 {
    // positive
} else if x < 0 {
    // negative
} else {
    // zero
}

// switch — no implicit fallthrough
switch status {
case .active:
    // ...
case .inactive, .pending:
    // multiple values per case
default:
    // required unless exhaustive
}

// explicit fallthrough
switch x {
case 0:
    fallthrough
case 1:
    // reached from 0 or 1
default:
    // other
}

// for-in ranges and arrays
for i in 0..<10 { }
for i in 0...10 { }   // closed — includes 10
for item in items { }

// while
var n = 0
while n < 5 { n += 1 }
```

---

### Defer

`defer` executes at scope exit. Multiple defers run in LIFO order.

```swift
let handle = fopen(path.any(), "r".any())
defer fclose(handle)

// Multi-statement cleanup
defer func() {
    cleanup(a)
    cleanup(b)
}()
```

---

### Generics

```swift
func identity<T>(value: T) -> T {
    return value
}

struct Box<T> {
    var value: T
}

let b      = Box{value: 42}
let result = identity(value: "hello")
```

---

### Tuples

```swift
let pair  = (1, true)
let point = (x: 10, y: 20)

// Destructuring
let (a, b) = pair

// Function return
func minMax(values: [int]) -> (min: int, max: int) {
    return (0, 100)
}
let (lo, hi) = minMax(values: [3, 1, 4])
```

---

### First-Class Functions and Anonymous Functions

```swift
// Function type variable
let double: func(int) -> int

// Anonymous function
let add = func(a: int, b: int) -> int { return a + b }

// Higher-order
func apply(values: [int], f: func(int) -> int) -> [int] { }
```

Anonymous functions capture by value. Use `mut` parameters for writeback:

```swift
var total = 0
run(n: &total, f: func(n: mut int) {
    n += 10
})
```

---

### Concurrency

Execution contexts are declared in the function signature and dispatched at
the call site with a postfix. There is no separate concurrency API to learn.

```swift
// Async
func fetchUser(id: int) async -> User { }
let user = fetchUser(id: 1).await()

// OS threads — shared memory, zero-copy
func crunchData(data: [float]) thread -> [float] { }
let result = crunchData(data: x).spawn(threads: 4)

// Processes — fully isolated memory
func isolatedWork(data: [float]) process -> [float] { }
let result = isolatedWork(data: x).fork(processes: 4)

// GPU kernels
func vectorAdd(a: [float], b: [float]) gpu -> [float] {
    let i = gpu.threadId
    return a[i] + b[i]
}
let result = vectorAdd(a: x, b: y).dispatch(gpu: 0, mem: 256)
```

**Channels** work across all contexts:

```swift
let ch: channel float = Channel(size: 8)
ch.send(42.0)
let val = ch.receive()
ch.close()
```

| Context   | Transport                       |
|-----------|---------------------------------|
| `async`   | shared memory, non-blocking     |
| `thread`  | shared memory, lightweight      |
| `process` | ring buffer, high-speed IPC     |

---

### Native Interface

Vertex reaches outside the managed runtime through a unified three-part
pattern: an **import path**, a **class declaration**, and a **package**.

The import path is the single source of truth for how the compiler emits.
Build tags are pure platform selectors — they never affect emission strategy.

```swift
package libc
build linux
import "lib/c"

class C : c {
    func fopen(path: any char, mode: any char) -> any opaque?
    func fwrite(ptr: any void, size: int, count: int, stream: mut any opaque) -> int
    func fclose(stream: mut any opaque) -> int
    func printf(fmt: any char, ...) -> int
}
```

**Import path → emission strategy:**

| Prefix     | Strategy                                        |
|------------|-------------------------------------------------|
| `lib/`     | linked call — signatures validated at link time |
| `linux/`   | inline syscall instruction, no linker           |
| `darwin/`  | `objc_msgSend` + selector dispatch              |
| `windows/` | vtable / COM slot dispatch                      |
| `gpu/`     | PTX / shader kernel emission                    |
| `metal/`   | hardware interrupt instruction, no linker       |

**Build tags — platform selectors only:**

| Tag       | Meaning                          |
|-----------|----------------------------------|
| `linux`   | compile this file for Linux      |
| `darwin`  | compile this file for Darwin     |
| `windows` | compile this file for Windows    |
| `metal`   | compile this file for bare metal |

Native class instances are zero-size compile-time dispatch surfaces — the
backend removes them entirely at compile time, no allocation, no runtime
overhead.

**Multi-platform packages** use filename suffixes instead of explicit build
tags:

```
tcp.vs          # all platforms — shared types and public API
tcp_linux.vs    # import "linux/syscalls"
tcp_darwin.vs   # import "darwin/objc/..."
tcp_windows.vs  # import "windows/com/..."
```

---

### Pointer Types and `.any()`

`any` in type position declares a read-only raw pointer. `mut any` declares a
mutable one. `.any()` is the postfix escape that produces a raw pointer from
any value — zero cost, no allocation, no copy.

**Type mapping:**

| Vertex annotation      | C equivalent    |
|------------------------|-----------------|
| `name: any char`       | `const char*`   |
| `name: mut any char`   | `char*`         |
| `name: any void`       | `const void*`   |
| `name: mut any void`   | `void*`         |
| `name: any opaque`     | `const struct*` |
| `name: mut any opaque` | `struct*`       |

**`.any()` — what it produces:**

| Value            | `.any()` type   |
|------------------|-----------------|
| `string`         | `const char*`   |
| `[T]`            | `const T*`      |
| `struct`         | `const struct*` |
| `class instance` | `const struct*` |

The returned pointer is valid only as long as the backing value is alive. For
ref-counted instances, retain the owning reference for the full duration of any
native call that holds the pointer.

**C interop example:**

```swift
extern "C" {
    func fopen(path: any char, mode: any char) -> mut any opaque?
    func fwrite(ptr: any void, size: int, count: int, stream: mut any opaque) -> int
    func fclose(stream: mut any opaque) -> int
    func printf(fmt: any char, ...) -> int
}

let f = fopen("/tmp/out.bin".any(), "wb".any())
defer fclose(f)

let data = [uint8](repeating: 0xFF, count: 256)
fwrite(data.any(), 1, data.byteSize(), f)
printf("wrote %d bytes\n".any(), data.len)
```

**Syscall example:**

```swift
package syscalls
build linux
import "linux/syscalls"

class Syscalls : syscalls {
    func write(fd: int, buf: any void, count: uint) -> int
    func exit(status: int) -> int
}
```

**Objective-C example:**

```swift
package foundation
build darwin
import "darwin/objc/foundation"

class NSString : foundation {
    func length(s: NSString) -> int
    func UTF8String(s: NSString) -> any char
    func stringWithUTF8String(str: any char) -> NSString?
}

var ns  = foundation.NSString()
let str = ns.stringWithUTF8String("hello".any())
let len = ns.length(str)
```

**COM / vtable example:**

```swift
package d3d11
build windows
import "windows/com/d3d11"

class IUnknown : d3d11 {
    func QueryInterface(obj: IUnknown, riid: any opaque, ppv: mut any opaque) -> int
    func AddRef(obj: IUnknown) -> uint
    func Release(obj: IUnknown) -> uint
}

class ID3D11Device : IUnknown {
    func CreateBuffer(
        d: ID3D11Device, desc: mut any opaque,
        init: mut any opaque, ppBuffer: mut any opaque) -> int
}
```

**Bare metal example:**

```swift
package bios
build metal
import "metal/int10h"

class Int10h : int10h {
    func set_video_mode(mode: uint8)
    func set_cursor_pos(page: uint8, row: uint8, col: uint8)
    func write_tty(char: uint8, page: uint8, color: uint8)
}

var b = bios.Int10h()
b.set_video_mode(0x03)
b.set_cursor_pos(0, 0, 0)
b.write_tty(0x41, 0, 0x07)
```

---

### Type Intrinsics

```swift
SockAddrIn.sizeof()     // compile-time size in bytes — called on type name
SockAddrIn.alignof()    // compile-time alignment in bytes

let data: [float] = [0.0, 1.0, 2.0]
data.byteSize()         // runtime byte length of array instance
data.len                // element count
```

---

### Type Aliases

```swift
type FILE    = any opaque
type size_t  = uint64
type errno_t = int
```

Aliases are resolved at compile time and may only appear at package level.

---

## Postfix Reference

| Postfix                                  | Meaning                               |
|------------------------------------------|---------------------------------------|
| `.new()`                                 | opt class instance into ref counting  |
| `.delete()`                              | manually free a class instance        |
| `.any()`                                 | escape value to raw C pointer         |
| `.try()`                                 | propagate Result error upward         |
| `.await()`                               | suspend until async function resolves |
| `.spawn()` / `.spawn(threads: n)`        | run on OS thread(s)                   |
| `.fork()` / `.fork(processes: n)`        | run in isolated process(es)           |
| `.dispatch()` / `.dispatch(gpu: n, mem:n)` | run on GPU                          |
| `.sizeof()`                              | compile-time size of type in bytes    |
| `.alignof()`                             | compile-time alignment in bytes       |
| `.byteSize()`                            | runtime byte length of an array       |
| `.len`                                   | element count of an array             |

---

## Out of Scope in 1.9

| Feature                                       | Status   |
|-----------------------------------------------|----------|
| Inheritance                                   | Removed  |
| String interpolation `\()`                    | Removed  |
| `_` parameter labels                          | Removed  |
| `self` keyword                                | Removed  |
| `static` keyword                              | Removed  |
| Methods inside structs or classes             | Removed  |
| `mutating` keyword                            | Removed  |
| Protocols                                     | Removed  |
| Extensions                                    | Removed  |
| `try` / `throws` / `do-catch`                | Removed  |
| Nested structs or classes                     | Removed  |
| Generic constraints (`where T:`)              | Deferred |
| Enums with associated values                  | Deferred |
| Access control                                | Deferred |
| `async let` / `TaskGroup` concurrency         | Deferred |
| `actor` keyword                               | Deferred |
| `select` over multiple channels               | Deferred |
| Labeled `break` / `continue`                  | Deferred |
| GPU grid/block control                        | Deferred |

---

## License

MIT