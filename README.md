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
vertex -os linux -arch arm64 main.vs -o myapp_linux_arm64
```

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
let mask  = 0xFF_00_FF
let flags = 0b1010_0011
let big   = 1_000_000
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

class FileHandler {
    var path: string
}
```

Classes are managed either manually or via opt-in reference counting:

```swift
// Manual — you call .delete(), deinit fires automatically
let a = Animal(name: "Rex")
defer a.delete()

// Reference counted — freed automatically when all owners go out of scope
let a = Animal(name: "Rex").new()

// Non-owning weak reference — does not increment the count
weak let b = a   // b: Animal?
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

var p = Point(3, 4)
p.reset()
```

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

**Result** — when the caller must handle Ok or Err explicitly:

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
    // required unless switch is exhaustive
}

// for-in ranges and arrays
for i in 0..<10 { }
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
```

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

A switch over an enum with all cases covered needs no `default`.

---

### First-Class Functions and Anonymous Functions

```swift
// Function types
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

---

### C Interop (FFI v2)

Vertex binds to C libraries with zero overhead. `.any()` escapes a value to
its raw memory address — no allocation, no copy.

```swift
extern "C" {
    func fopen(path: any char, mode: any char) -> mut any opaque?
    func fwrite(ptr: any void, size: int, count: int, stream: mut any opaque) -> int
    func fread(ptr: mut any void, size: int, count: int, stream: mut any opaque) -> int
    func fclose(stream: mut any opaque) -> int
    func printf(fmt: any char, ...) -> int
}

let f = fopen("/tmp/out.bin".any(), "wb".any())
defer fclose(f)

let data = [uint8](repeating: 0xFF, count: 256)
fwrite(data.any(), 1, data.byteSize(), f)
printf("wrote %d bytes\n".any(), data.len)
```

**Pointer type mapping:**

| Vertex annotation      | C equivalent    |
|------------------------|-----------------|
| `name: any char`       | `const char*`   |
| `name: mut any char`   | `char*`         |
| `name: any void`       | `const void*`   |
| `name: mut any void`   | `void*`         |
| `name: any opaque`     | `const struct*` |
| `name: mut any opaque` | `struct*`       |

**Build tags** control linking and ABI:

| Tag                 | Effect                                 |
|---------------------|----------------------------------------|
| `linux`             | linked call, C++ Itanium mangling      |
| `linux, syscalls`   | inline syscall instruction, no linker  |
| `darwin`            | linked call, C++ Itanium mangling      |
| `darwin, objc`      | `objc_msgSend` + selector refs         |
| `windows`           | linked call, C++ MSVC mangling         |
| `windows, com`      | vtable slot offsets, positional        |

Platform-specific files use filename suffixes — no explicit build tag needed:

```
tcp.vs          # all platforms
tcp_linux.vs    # linux only
tcp_darwin.vs   # darwin only
tcp_windows.vs  # windows only
```

**Objective-C interop** (`build darwin, objc`):

```swift
extern class NSString {
    func length(s: NSString) -> int
    func UTF8String(s: NSString) -> any char
    func stringWithUTF8String(str: any char) -> NSString?
}

let ns  = NSString.stringWithUTF8String("hello".any())
let len = ns.length()
```

**COM interop** (`build windows, com`):

```swift
extern class IUnknown {
    func QueryInterface(obj: IUnknown, riid: any opaque, ppv: mut any opaque) -> int
    func AddRef(obj: IUnknown) -> uint
    func Release(obj: IUnknown) -> uint
}

extern class ID3D11Device : IUnknown {
    func CreateBuffer(d: ID3D11Device, desc: mut any opaque, init: mut any opaque, ppBuffer: mut any opaque) -> int
}
```

---

### Postfix Reference

| Postfix                                 | Meaning                               |
|-----------------------------------------|---------------------------------------|
| `.new()`                                | opt class instance into ref counting  |
| `.delete()`                             | manually free a class instance        |
| `.any()`                                | escape value to raw C pointer         |
| `.try()`                                | propagate Result error upward         |
| `.await()`                              | suspend until async function resolves |
| `.spawn()` / `.spawn(threads:)`         | run on OS thread(s)                   |
| `.fork()` / `.fork(processes:)`         | run in isolated process(es)           |
| `.dispatch()` / `.dispatch(gpu:, mem:)` | run on GPU                            |
| `.sizeof()`                             | compile-time size of type in bytes    |
| `.alignof()`                            | compile-time alignment in bytes       |
| `.byteSize()`                           | runtime byte length of an array       |
| `.len`                                  | element count of an array             |

---

## License

MIT