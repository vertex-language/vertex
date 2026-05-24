# Vertex Native Interface

## 1. Concept

Native Interface is Vertex's unified model for reaching outside the managed
runtime. Every foreign target — C libraries, OS syscalls, GPU kernels, bare
metal interrupts, platform dispatch — is expressed through the same three-part
pattern:

1. **Import path** — tells the compiler which emission strategy to use
2. **Class declaration** — typed signature manifest bound to the import target
3. **Package** — wraps the class as a typed, importable export

The import path does all the routing. The class provides the contract.
The compiler enforces it.

---

## 2. The Three-Part Pattern

```swift
package libc                 // 3. package — the export
build linux
import "lib/c"               // 1. import path — emission strategy

class C : c {                // 2. class — signature manifest
    func fopen(path: any char, mode: any char) -> any opaque?
    func fwrite(ptr: any void, size: int, count: int, stream: mut any opaque) -> int
    func fclose(stream: mut any opaque) -> int
    func malloc(size: int) -> any void?
    func free(ptr: mut any void)
    func printf(fmt: any char, ...) -> int
}
```

**Rules:**

* The class name after `:` must match the final segment of the import path.
* Methods in a native class have no bodies — the backend owns the implementation.
* The package is the public boundary. Consumers import the package, not the
  raw class.
* One import path per file. One class per import target.

---

## 3. Import Path Protocol

The first path segment is the **emission prefix** — it routes the compiler
to the correct backend strategy. Everything after it identifies the specific
target.

```
"lib/*"          → linker — validate sigs at link time, emit linked call
"linux/*"        → inline syscall instruction, no linker
"darwin/*"       → platform dispatch (objc_msgSend / selectors)
"windows/*"      → vtable / COM slot dispatch
"gpu/*"          → kernel emission (PTX / shader IR)
"metal/*"        → bare metal (interrupts, ports, MMIO)
```

The path is the single source of truth for how the compiler emits.
No build tag modifiers, no extern qualifiers, no extra annotations needed.

---

## 4. Pointer Model

`any` in type position declares a raw pointer. `mut any` declares a mutable
raw pointer. These are the only way to express pointer types in Vertex —
they appear in function signatures, struct fields, and type aliases.

| Vertex                 | C equivalent    |
|------------------------|-----------------|
| `name: any T`          | `const T*`      |
| `name: mut any T`      | `T*`            |
| `name: any void`       | `const void*`   |
| `name: mut any void`   | `void*`         |
| `name: any char`       | `const char*`   |
| `name: any opaque`     | `const struct*` |
| `name: mut any opaque` | `struct*`       |

---

## 5. `.any()` — Escape to Raw Pointer

`.any()` is a postfix escape — it steps outside Vertex's type system and
returns the raw memory address of the value. It follows the same postfix
convention as `.new()`, `.try()`, `.await()`, and `.dispatch()`.

| Value            | `.any()` produces | C type          |
|------------------|-------------------|-----------------|
| `string`         | `const char*`     | `const char*`   |
| `[T]`            | pointer to `T[0]` | `const T*`      |
| `struct`         | `const struct*`   | `const struct*` |
| `struct.field`   | pointer to field  | `const T*`      |
| `class instance` | `const struct*`   | `const struct*` |

**Rules:**

* `.any()` may be called on strings, arrays, structs, struct fields, and
  class instances.
* The returned pointer is valid only as long as the backing value is alive.
  Passing it to a native call and then freeing the backing value is undefined
  behavior.
* For ref-counted instances (`.new()`), retain the owning reference for the
  full duration of any native call that holds the pointer.
* `.any()` is zero-cost — no allocation, no copy, no runtime overhead.
* `.any()` does not appear in type position. Type-level pointer declarations
  use the `any T` keyword form (`name: any char`, `name: mut any void`, etc.).

---

## 6. Type Intrinsics

Type-level compile-time intrinsics are called on the type name, not an
instance. They resolve to constants at compile time — zero runtime cost.

```swift
SockaddrIn.sizeof()    // size of the type in bytes
SockaddrIn.alignof()   // alignment requirement in bytes
```

Instance-level intrinsic for byte length of arrays:

```swift
let vertices: [float] = [0.0, 0.5, 0.0]
vertices.byteSize()    // len * sizeof(element) — for C byte-count arguments
vertices.len           // element count
```

**Rules:**

* `.sizeof()` and `.alignof()` are only valid on type names, not instances.
* `.byteSize()` is only valid on array instances.
* All three resolve at compile time.

---

## 7. Type Aliases

```swift
type FILE    = any opaque
type size_t  = uint64
type errno_t = int
```

**Rules:**

* `type` declares an alias — the two names are interchangeable.
* Aliases may appear at package level only, not inside functions or blocks.
* Aliases are resolved at compile time — no runtime representation.

---

## 8. Variadic Functions

```swift
class C : c {
    func printf(fmt: any char, ...) -> int
    func sprintf(buf: mut any char, fmt: any char, ...) -> int
    func open(path: any char, flags: int, ...) -> int
}
```

**Rules:**

* `...` marks a variadic tail — zero or more additional arguments.
* `...` may only appear as the final parameter.
* `...` is only valid in native class declarations — Vertex functions may
  not be variadic.
* Arguments passed to `...` slots are not type-checked by the compiler.

---

## 9. Build Tags

Build tags are **pure platform selectors**. Emission strategy is fully
owned by the import path.

| Tag       | Meaning                          |
|-----------|----------------------------------|
| `linux`   | compile this file for linux      |
| `darwin`  | compile this file for darwin     |
| `windows` | compile this file for windows    |
| `metal`   | compile this file for bare metal |

---

## 10. Multi-Platform Packages

Platform-specific files use a filename suffix matching the build tag. The
compiler selects the correct file for the target platform — no explicit
`build` tag required inside the file when the suffix is present.

```
tcp.vs          // compiled everywhere — shared types, constants, public API
tcp_linux.vs    // import "linux/syscalls"    — linux native binding
tcp_darwin.vs   // import "darwin/objc/..."   — darwin native binding
tcp_windows.vs  // import "windows/com/..."   — windows native binding
```

---

## 11. Instances

Native class instances are **zero-size compile-time dispatch surfaces**.
The backend removes them entirely — no allocation, no runtime overhead.

```swift
var c = libc.C()             // zero bytes at runtime
c.printf("hello\n".any())    // emits direct linked call — instance vanishes
```

The instance exists only to give the developer a consistent, typed call
surface that matches every other Vertex object.

---

## 12. Strategies

### 12.1 `lib/` — Linked

```swift
package libc
build linux
import "lib/c"

class C : c {
    func fopen(path: any char, mode: any char) -> any opaque?
    func fwrite(ptr: any void, size: int, count: int, stream: mut any opaque) -> int
    func fread(ptr: mut any void, size: int, count: int, stream: mut any opaque) -> int
    func fclose(stream: mut any opaque) -> int
    func malloc(size: int) -> any void?
    func free(ptr: mut any void)
    func printf(fmt: any char, ...) -> int
}
```

Compiler: fires the linker, validates signatures against the library at link
time. Emits a standard linked call at each call site.

Usage:
```swift
var c = libc.C()
let f = c.fopen("/tmp/out.bin".any(), "wb".any())
defer c.fclose(f)

let data = [uint8](repeating: 0xFF, count: 256)
c.fwrite(data.any(), 1, data.byteSize(), f)
c.printf("wrote %d bytes\n".any(), data.len)
```

---

### 12.2 `linux/` — Syscall

```swift
package syscall2
build linux
import "linux/syscalls"

class Syscalls : syscalls {
    func open(path: any char, flags: int, mode: uint) -> int
    func openat(dirfd: int, path: any char, flags: int, mode: uint) -> int
    func openat2(dirfd: int, path: any char, how: any opaque, size: uint) -> int
    func close(fd: int) -> int
    func close_range(first: uint, last: uint, flags: uint) -> int
    func read(fd: int, buf: mut any void, count: uint) -> int
    func write(fd: int, buf: any void, count: uint) -> int
}
```

Compiler: resolves syscall number from method name (`open` → `NR_open`),
determines arg count, emits inline syscall instruction. No linker involved.

Usage:
```swift
var s = syscall2.Syscalls()
let fd = s.open("/tmp/out".any(), O_WRONLY, 0o644)
defer s.close(fd)
s.write(fd, data.any(), data.byteSize())
```

---

### 12.3 `darwin/` — Platform Dispatch

```swift
package foundation
build darwin
import "darwin/objc/foundation"

class NSString : foundation {
    func length(s: NSString) -> int
    func UTF8String(s: NSString) -> any char
    func isEqualToString(s: NSString, other: NSString) -> bool
    func substringFromIndex(s: NSString, from: int) -> NSString
    func stringWithUTF8String(str: any char) -> NSString?
}
```

Compiler: emits `objc_msgSend` with selector derived from method name.
First typed param is the receiver. No receiver = class-side method.

Usage:
```swift
var ns = foundation.NSString()
let str = ns.stringWithUTF8String("hello".any())
let len = ns.length(str)
```

---

### 12.4 `windows/` — COM / Vtable

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
    func CreateTexture2D(
        d: ID3D11Device, desc: mut any opaque,
        init: mut any opaque, ppTexture: mut any opaque) -> int
}
```

Compiler: emits vtable slot dispatch. Declaration order = slot order.
Parent class slots are prepended — child methods begin at `len(parent)`.
Lifetime managed via `AddRef` / `Release`.

---

### 12.5 `gpu/` — Kernel Emission

```swift
package cuda2
build linux
import "gpu/cuda"

class Cuda : cuda {
    let threadIdx: (x: uint, y: uint, z: uint)
    let blockIdx:  (x: uint, y: uint, z: uint)
    let blockDim:  (x: uint, y: uint, z: uint)
    let gridDim:   (x: uint, y: uint, z: uint)

    func syncThreads()
    func syncWarp(mask: uint)
    func atomicAdd(addr: mut any float, val: float) -> float
    func atomicAdd(addr: mut any int,   val: int)   -> int
    func fmaf(a: float, b: float, c: float) -> float
    func rsqrtf(x: float) -> float
}
```

Compiler: emits PTX. `threadIdx` → `%tid` register reads. `syncThreads()`
→ `bar.sync`. `fmaf` → `fma.rn.f32`. Instance removed entirely.

GPU functions use the `gpu` qualifier between the parameter list and the
return arrow. Shared memory uses the `shared` storage qualifier — the backend
allocates in `.shared` PTX memory space. The call site uses `.dispatch()`.

```swift
let TILE = 16

func matMul(a: [float], b: [float], out: mut [float], n: int) gpu {
    var c = cuda2.Cuda()

    shared var tileA = [float](count: TILE * TILE)
    shared var tileB = [float](count: TILE * TILE)

    let tx  = c.threadIdx.x
    let ty  = c.threadIdx.y
    let row = c.blockIdx.y * TILE + ty
    let col = c.blockIdx.x * TILE + tx

    var sum = 0.0

    for t in 0 ..< n / TILE {
        tileA[ty * TILE + tx] = a[row * n + t * TILE + tx]
        tileB[ty * TILE + tx] = b[(t * TILE + ty) * n + col]
        c.syncThreads()

        for k in 0 ..< TILE {
            sum = c.fmaf(tileA[ty * TILE + k], tileB[k * TILE + tx], sum)
        }
        c.syncThreads()
    }

    out[row * n + col] = sum
}

// call site — out is mut, caller passes &result
var result = [float](repeating: 0.0, count: n * n)
matMul(a: a, b: b, out: &result, n: n).dispatch()
```

---

### 12.6 `metal/` — Bare Metal

```swift
package bios
build metal
import "metal/int10h"

class Int10h : int10h {
    func set_video_mode(mode: uint8)
    func set_cursor_size(start_line: uint8, end_line: uint8)
    func set_cursor_pos(page: uint8, row: uint8, col: uint8)
    func get_cursor_pos(
        page: uint8,
        row: mut any uint8, col: mut any uint8,
        start: mut any uint8, end: mut any uint8)
    func scroll_up(
        lines: uint8, attr: uint8,
        top: uint8, left: uint8, bottom: uint8, right: uint8)
    func scroll_down(
        lines: uint8, attr: uint8,
        top: uint8, left: uint8, bottom: uint8, right: uint8)
    func read_char_attr(page: uint8, char: mut any uint8, attr: mut any uint8)
    func write_char_attr(page: uint8, char: uint8, attr: uint8, count: uint16)
    func write_char(page: uint8, char: uint8, count: uint16)
    func write_tty(char: uint8, page: uint8, color: uint8)
    func get_video_mode(
        mode: mut any uint8, cols: mut any uint8, page: mut any uint8)
    func write_string(
        page: uint8, attr: uint8, flags: uint8,
        row: uint8, col: uint8, str: any char, len: uint16)
    func write_pixel(page: uint8, color: uint8, col: uint16, row: uint16)
    func read_pixel(page: uint8, col: uint16, row: uint16) -> uint8
    func set_palette(palette_id: uint8, color: uint8)
    func set_active_page(page: uint8)
}
```

Compiler: emits hardware interrupt instruction with register mapping
derived from the interrupt number in the path (`int10h` → `INT 10h`).
No linker, no runtime.

Usage:
```swift
var b = bios.Int10h()
b.set_video_mode(0x03)
b.set_cursor_pos(0, 0, 0)
b.write_tty(0x41, 0, 0x07)
```

---

## 13. Rules Summary

| Rule | Detail |
|------|--------|
| Import prefix routes emission | `lib/` links, `linux/` syscalls, `gpu/` PTX, etc. |
| Class name matches import tail | `import "lib/c"` → `class C : c` |
| Methods have no bodies | backend owns the implementation |
| Instances are zero-size | removed entirely by the backend |
| Build tags select platform only | emission strategy lives in the import path |
| Package is the export boundary | consumers import the package, not the class |
| Sigs enforced at compile time | mismatch = error before anything links or runs |
| Variadic `...` in native only | Vertex functions may not be variadic |
| `mut` after label before type | `out: mut [float]` not `mut out: [float]` |
| GPU call site uses `.dispatch()` | consistent with postfix execution model |

---

## 14. Postfix Summary

| Postfix        | Meaning                                 |
|----------------|-----------------------------------------|
| `.any()`       | escape to raw pointer                   |
| `.sizeof()`    | compile-time size of type in bytes      |
| `.alignof()`   | compile-time alignment of type in bytes |
| `.byteSize()`  | runtime byte length of array instance   |
| `.dispatch()`  | execute gpu kernel, return result       |