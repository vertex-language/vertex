# Vertex

**Language Spec 2.2 · Compiler 0.4.0**

Vertex is a statically-typed systems and application programming language built
for explicit control, zero-overhead C interop, and first-class concurrency.
Its syntax draws from Swift and Go, while its call-site execution sigils and
layered ownership model make it uniquely suited to systems work — from GPU
kernels and AI inference to networking, file I/O, and low-level kernel
programming.

---

## Contents

- [Install](#install)
- [Quick Start](#quick-start)
- [Language Tour](#language-tour)
  - [Variables and Types](#variables-and-types)
  - [Ownership and Functions](#ownership-and-functions)
  - [Generics](#generics)
  - [Structs and Classes](#structs-and-classes)
  - [Raw Pointers](#raw-pointers)
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

import "lib/c"

class C : c {
    func printf(fmt: ...string) -> int32
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
    C.printf("--- Recursive ---\n")
    var i: int32 = 0
    while true {
        if i >= 10 { break }
        C.printf("fib(%d) = %d\n", i, fibRecursive(i))
        i = i + 1
    }

    C.printf("--- Iterative ---\n")
    var j: int32 = 0
    while true {
        if j >= 10 { break }
        C.printf("fib(%d) = %d\n", j, fibIterative(j))
        j = j + 1
    }

    return 0
}
```

`C` declares no `init func`, so it's a flat namespace — its functions are
called directly on the type (`C.printf(...)`), and `C()` would be a compile
error. See [Native Interface](#native-interface).

```sh
vertex fib.vs

or compile

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
| `uint8` / `byte` | `uint8_t` |
| `uint16` | `uint16_t` |
| `uint64` | `uint64_t` |
| `float32` | `float` |
| `float64` | `double` |
| `bool` | `bool` |
| `string` | `let` → rodata, `var` → heap |

`int` is an alias for `int32`; `uint` is an alias for `uint32`; `byte` is the
preferred spelling of `uint8`. Vertex's `char` is a 4-byte Unicode scalar, not
a C byte — C strings cross the FFI boundary as `*const char` and marshal to
`string` (see [Native Interface](#native-interface)).

The `as` operator performs explicit type conversion for numeric widening,
pointer reinterpretation, and float-to-integer truncation:

```vertex
let small: int32 = 42
let wide  = small as int64          // integer widening — sign-extended
let big   = small as uint64         // zero-extended

let f: float64 = 3.99
let i = f as int32                  // truncates toward zero → 3

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

### Ownership and Functions

```vertex
func add(a: int32, b: int32) -> int32 {
    return a + b
}

// call with positional or labeled arguments
add(1, 2)
add(a: 1, b: 2)
```

A parameter's convention is picked in the signature. The caller writes at
most one thing at the call site, and only for the owning case:

```vertex
func inspect(w: Widget)      // shared — read-only view, bare, always
func rename(w: mut Widget)   // exclusive — mutates the caller's binding
func archive(w: var Widget)  // owning — transfer or copy, chosen by the caller
```

```vertex
func increment(n: mut int32) {
    n += 1
}

var count = 0
increment(count)      // bare — never a keyword at the call site
```

`mut` is a pointer under the hood — `increment(count)` lowers to
`increment(&count)` — but the address is never written by the caller; the
callee's signature is the only thing that decides.

An owning (`var`) parameter is where the caller's choice becomes visible.
Writing `.transfer()` moves the value — the source dies, the move costs
O(1). Leaving it off copies — the source survives, the copy costs O(data):

```vertex
func archive(w: var Widget) {
    storage.push(w)
}

var w = Widget(1)
archive(w.transfer())   // TRANSFER — w is dead after this line
```

```vertex
var w = Widget(1)
archive(w)               // COPY — w survives, archive gets an independent copy
inspect(w)                // ok
```

The same bare-copy / `.transfer()`-move pair governs ordinary bindings, not
just function arguments:

```vertex
let a = w              // COPY
let b = w.transfer()   // TRANSFER
```

Two heap doors extend the same rules onto the heap — `unique(...)` for sole
ownership, `shared(...)` for a reference-counted, cheaply cloneable handle.
Both are covered in [Structs and Classes](#structs-and-classes).

---

### Generics

Type parameters use square brackets, matching type-argument position:

```vertex
func identity[T](value: T) -> T {
    return value
}

struct Box[T] {
    value: T
}

let b      = Box[int32]{value: 42}
let result = identity(value: "hello")   // T inferred from the argument
```

A bare name is constraint `any` — `[T]` means `[T: any]`. Under `any`, only
assignment and argument passing are available; no `<`, no `+`, no `==`.

**Constraints** are their own declaration — a compile-time type set,
optionally paired with required methods:

```vertex
constraint Ordered {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~string
}

func min[T: Ordered](a: T, b: T) -> T {
    if a < b { return a }
    return b
}

let m = min[float64](3.14, 2.71)
```

`~T` admits `T` and every type whose underlying type is `T`; a bare `T`
(no tilde) admits only `T` exactly. `constraints.Number`,
`constraints.Integer`, `constraints.Ordered`, and friends are predeclared in
`builtins/constraints`, alongside `any` and `comparable`.

Every instantiation is compiled as a separate concrete body with type
arguments substituted in — no runtime type information, no vtable. A
method-constraint call lowers to a direct call on the concrete type.

---

### Structs and Classes

**Structs** and **classes** are both stack-resident value types by default —
copied by value under the same rules as any other binding. A `class`
differs from a `struct` only in its member/method model; declaring
something a `class` does not, by itself, put it on the heap.

```vertex
struct Vec2 {
    x: float32
    y: float32
}

// shared receiver — read-only view
func (v: Vec2) describe() {
    // ...
}

// exclusive receiver — mutations affect the caller's binding
func (v: mut Vec2) scale(factor: float32) {
    v.x *= factor
    v.y *= factor
}

var pos = Vec2{x: 1.0, y: 2.0}
pos.scale(factor: 2.0)   // bare — no & needed, mut is a signature fact
```

Struct literals require all field labels; positional initialization is not
supported. Trailing commas are valid in multiline forms:

```vertex
let p = Vec2{
    x: 3.0,
    y: 4.0,
}
```

**Classes** add `init`/`deinit` and identity (`===`), but share a struct's
layout and lifetime rules exactly — no header, no vtable, no inheritance:

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
    // runs automatically when `a`'s liveness ends
}

let dog = Animal(name: "Rex", health: 100)
// dog.deinit() runs at the end of this scope — no defer needed
```

Because every binding's liveness is tracked statically, teardown (`deinit`)
is emitted wherever that liveness ends — no garbage collector, no manual
delete call for a stack value. The only way onto the heap is through one of
two doors:

```vertex
// sole ownership — no refcount
var u = unique(Animal(name: "Rex", health: 100))

// reference counted — cloneable handle
var s  = shared(Animal(name: "Luna", health: 100))
var s2 = s                    // cheap refcount bump, not a deep copy

// weak — observes a shared value without keeping it alive
var w = weak(s)

let observer, err = w.upgrade()
if err != "" {
    // the shared value is gone
}
```

`unique(...)` and `shared(...)` still tear down automatically — `unique`
frees at the end of its owner's liveness, `shared` frees once its strong
count reaches zero. Neither needs an explicit delete call.

`===` / `!==` (classes only) compare storage addresses — "same allocation?",
never "same bytes?"; that's `==`'s job.

---

### Raw Pointers

`typed_ptr T` is the raw, last-resort pointer — no ownership tracking, no
refcount, none of the discipline above. Reach for it only when `mut T` /
`[]T` genuinely can't express what's needed.

```vertex
var x: int32 = 42
let p = &x           // address-of — int32 -> typed_ptr int32
let v = &p            // dereference — typed_ptr int32 -> int32
&p = 99               // write through p

let p2 = p.add(1)     // arithmetic is a method, scaled by sizeof(T)
let n: int64 = p2.diff(p)

let buf, err = new[uint8](1024)     // zeroed by default
if err != "" { panic("allocation failed: " + err) }
defer delete(buf)

buf.setAt(0, 0xFF)
let b = buf.at(0)
```

`nil` is the one general value in the language, and it exists solely for
`typed_ptr T`:

```vertex
var p: typed_ptr int32 = nil
if p == nil { }
```

---

### Arrays and Maps

**Fixed arrays** are stack-allocated; the size comes first in the brackets
and is part of the type.

```vertex
var buf:  [1024]uint8             // zero-filled, no initializer needed
var mask: [64]uint8
mask.fill(0xFF)

var coords: [3]int32 = [10, 20, 30]
let flags:  [3]uint8 = [0xFF, 0x00, 0xAB]

// nested (multidimensional)
let matrix: [2][2]float32 = [
    [0.0, 1.0],
    [1.0, 0.0],
]
```

**Dynamic arrays** are heap-allocated and growable — empty brackets before
the element type. This is the container exception: the backing storage is
implicitly heap-allocated, but it's still owned and torn down the normal way
through whatever binding holds the array.

```vertex
var items: []int32 = []

items.push(42)            // add to end
items.unshift(0)          // add to front
let last  = items.pop()   // remove from end
let first = items.shift() // remove from front

items.reserve(64)         // pre-allocate capacity
```

Methods that return a new array allocate on the heap, torn down
automatically like anything else once their binding's liveness ends:

```vertex
let doubled = items.map(func(x: int32) -> int32 { return x * 2 })
let evens   = items.filter(func(x: int32) -> bool { return x % 2 == 0 })
let sub     = items.slice(1, 3)
```

The rule at a glance:

| Form | Storage | Growable |
|------|---------|----------|
| `var buf: [N]T` | stack | no |
| `let arr = [...]` | stack / rodata | no |
| `var x: []T = []` | heap | yes |
| `var x = [...]` | heap | yes |

**Maps** use brace literals. A type annotation is required for empty maps.
Assigning `nil` to a key removes it — the one other place besides
`typed_ptr T` where the grammar admits the literal.

```vertex
let scores = {"alice": 42, "bob": 7}
let val = scores["alice"]

var config: map[string]int32 = {}

config["workers"] = 4
config["verbose"] = nil    // removes key
```

---

### Tuples

**Rule of thumb: parens build, bare commas unbuild.** Parens appear when a
tuple is *constructed* — a literal or a type annotation. The moment a tuple
is being *pulled apart* — a `let` destructure, or a `return` handing
multiple values back — it is written bare, with no wrapping parens.

```vertex
let pair  = (1, true)
let point = (x: 10, y: 20)

func divmod(a: int32, b: int32) -> (int32, int32) {
    return a / b, a % b
}
let quotient, remainder = divmod(10, 3)

func minMax(values: []int32) -> (min: int32, max: int32) {
    return 0, 100
}
let lo, hi = minMax(values: [3, 1, 4])
```

`()` is the empty tuple and is an alias for `void`. Tuples are value types —
assignment copies all elements. `==`, `!=`, `<`, `>`, `<=`, `>=` work on
tuples whose elements are all comparable, up to 6 elements.

Channels can carry tuples for paired data:

```vertex
let stream = chan[(int32, bool)](64)

select {
case msg = stream.receive():
    let val, ok = msg
    if ok { print(val) }
default:
}
```

---

### Enums

Vertex enums support unit variants, tuple variants (positional associated
data), or a mix of both. The `case` keyword is used only inside `switch`
statements — not in enum declarations.

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

**Explicit discriminants** require a backing integer type and are only
valid on all-unit enums. Auto-increment applies to unspecified variants:

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

Integer-to-enum conversion is an ordinary `switch` with a fallback case —
there's no optional type to lean on:

```vertex
func statusFromInt(n: int32) -> Status {
    switch n {
    case 0: return .Inactive
    case 1: return .Active
    case 2: return .Pending
    default: return .Inactive
    }
}
```

Enums do not support `==`, `!=`, or any comparison operator — equality is
expressed via `switch`. Enums are value types; assignment copies.

---

### Error Handling

There is no optional type and no propagation operator. Every fallible or
possibly-absent value is returned as a plain tuple: the value, and a string
that is empty (`""`) on success.

```vertex
func parseInt(s: string) -> (int32, string) {
    if s == "" { return 0, "empty string" }
    return 42, ""
}

func findUser(id: int32) -> (User, string) {
    if id < 0 { return User{}, "invalid id" }
    return User(id), ""
}
```

Checking the error is the only pattern — the happy path continues directly
below it:

```vertex
let n, err = parseInt(s: "42")
if err != "" {
    log.printf("failed: %s\n", err)
    return
}
// n is usable past this point
```

Chained calls repeat the same shape at every step — it does not get shorter
as call depth grows, and every branch stays visible in the text:

```vertex
func loadModel(path: string) -> (Model, string) {
    let text, err = readFile(path)
    if err != "" {
        return Model{}, err
    }

    let config, err2 = parseConfig(text)
    if err2 != "" {
        return Model{}, err2
    }

    return Model(config), ""
}
```

A function that may simply find nothing uses the exact same shape as one
that can fail outright — absence is not a special case. On the error path
the paired value is always the type's zero value (`0`, `""`, a zeroed
struct/class); the compiler doesn't enforce checking `err` first, matching
the convention's explicit-over-automatic philosophy.

---

### Concurrency

Every function is written as an ordinary synchronous function. The caller
decides how a call runs by prefixing it with an execution sigil at the call
site — there are no `async` or `thread` function qualifiers.

| Sigil | Reality | Best for |
|---|---|---|
| `async` | Virtual thread — state machine, scheduler-driven | Millions of idle connections, massive fan-out |
| `thread` | Real OS thread — shared-memory concurrency | Heavy CPU work, blocking calls |
| `gpu` | PTX / SPIR-V kernel | Matrix math, parallel workloads |
| `tpu` | Tensor-typed kernel | AI inference |

```vertex
let a = async fetch_network(id: 1)
let b = thread heavy_compute(data: x)
```

A `thread`/`async` call whose function returns `T` evaluates to a
receive-only channel of `T` carrying exactly one value (Path A):

```vertex
let worker = thread func(seed: int32) -> float32 {
    return crunch_numbers(seed)
}(105)

let final_data = worker.receive()
```

For many values, construct a channel explicitly and hand it to the worker
as an ordinary argument (Path B):

```vertex
let out_stream = chan[float32](64)

thread func(data: []float32, ch: chan float32) {
    for chunk in data {
        ch.send(process(chunk))
    }
    ch.close()
}(dataset, out_stream)

while true {
    let chunk, err = out_stream.tryReceive()
    if err != "" { break }
    print(chunk)
}
```

**Channels** are a built-in generic type, constructed like any other
instantiable type — type argument in `[...]`, capacity as an ordinary
constructor argument:

```vertex
let ch1 = chan[float32]()          // unbuffered
let ch2 = chan[int32](64)          // buffered
let ch3: chan float32 = chan[float32]()
```

| Method | Blocking | Returns |
| --- | --- | --- |
| `.send(value)` | yes | `void` |
| `.receive()` | yes | `T` |
| `.trySend(v)` | no | `bool` |
| `.tryReceive()` | no | `(T, string)` |
| `.close()` | no | `void` |

**`select`** multiplexes channels, suspending with 0% CPU until a case is
ready. Adding a `default` case makes it non-blocking:

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

#### GPU and TPU

```vertex
let d = gpu(blocks: 16, threads: 256) matrix_mult(x, y)
```

The `gpu` sigil, with an optional `(blocks: n, threads: n)` config, compiles
the call to a PTX/SPIR-V kernel. Function bodies are ordinary Vertex — no
restricted types or constructs.

```vertex
func vecAdd(a: tensor[float32, 1024], b: tensor[float32, 1024]) -> tensor[float32, 1024] {
    return a + b
}

var ha: [1024]float32
var hb: [1024]float32

let sum = tpu vecAdd(ha, hb)   // sum: [1024]float32
```

`tpu` channels a host array call into a `tensor`-typed function body, and
channels the `tensor` result back to a plain host array — same
return-type-channeling rule as `gpu`. `tensor[ElementType, Shape...]` is
valid only inside a `tpu`-sigil function body — no subscripting, only
elementwise ops and the `tpu.` builtin namespace (`Abs`, `Dot`, `Reshape`,
`Sum`, and friends).

---

### Native Interface

Foreign interop is a structural contract: you describe the shape of the
external library using native Vertex types, and the backend emits the
correct calling convention for the target — C ABI, Objective-C message
passing, or a JS bundle's call/property-access shape.

```vertex
import "lib/sdl2"
import "darwin/framework/AppKit"
import "wasm/wasi_snapshot_preview1"
import "js/websocket"
```

A foreign resource is declared as an opaque handle:

```vertex
type SDL_Window = abstract
type NSView     = abstract
```

`abstract` says "structure exists, but Vertex declines to model it" — no
arithmetic, no dereference, no stride. Its zero value is legal only as an
error-path value paired with a non-empty error string; there's no
comparable `nil` for it.

A `class` bound to an import path (`class Name : path`) is a **flat
namespace** if it declares no `init func` — its functions are called
directly on the type, and it cannot be instantiated:

```vertex
class SDL2_API : sdl2 {
    func SDL_CreateWindow(title: string, x: int32, y: int32, w: int32, h: int32, flags: uint32) -> (SDL_Window, string)
    func SDL_DestroyWindow(window: SDL_Window)
}

let window, err = SDL2_API.SDL_CreateWindow("game", 0, 0, 800, 600, 2)
```

Declaring one or more `init func` instead makes it an **instantiable
object** — at most one unnamed constructor (`init func() -> Self`, resolved
by bare `Type(...)`) and any number of named ones (`init func someName(...)
-> Self`, resolved by `Type.someName(...)`):

```vertex
class NSWindow : AppKit {
    init func() -> NSWindow
    init func initWithContentRect(
        contentRect: Rect, styleMask: uint64, backing: uint64, defer: bool
    ) -> NSWindow

    func center()
}
```

A parameter label is matched positionally, never looked up as an
identifier — so a Vertex keyword (`defer` above) is a perfectly legal
label. Foreign calls that can fail map to the standard error tuple; a call
that doesn't synchronously fail (like a JS constructor whose errors surface
later as an event) is declared as an ordinary bare return.

| Foreign Shape | Vertex form |
| --- | --- |
| `const char*` | `string` — marshalled NUL-terminated at the boundary |
| `T*` (writable scalar out-param) | `mut T` |
| `T*` + length | `[]T` (read) / `mut []T` (write) |
| pointer strided manually | `typed_ptr T` — last resort |

An abstract interface is declaration-only — no bodies, no visibility
modifiers, no ownership keywords. Ownership lives in the **wrapper**
written on top of it, whose `deinit` releases the resource:

```vertex
class Window {
    handle: SDL_Window
}

func (w: Window) init(title: string) {
    let handle, err = SDL2_API.SDL_CreateWindow(title, 0, 0, 800, 600, 2)
    if err != "" {
        panic("Failed to create window: " + err)
    }
    w.handle = handle
}

func (w: Window) deinit() {
    SDL2_API.SDL_DestroyWindow(w.handle)
}

func (w: Window) size() -> (int32, int32) {
    var width:  int32 = 0
    var height: int32 = 0
    SDL2_API.SDL_GetWindowSize(w.handle, width, height)
    return width, height
}
```

Only a **non-capturing** function converts across an abstract boundary —
one word, no environment:

```vertex
func on_event(code: int32) -> int32 { return 0 }

SDL2_API.SDL_SetEventFilter(on_event)   // legal
```

```vertex
var count = 0

SDL2_API.SDL_SetEventFilter(func(code: int32) -> int32 {
    count += 1        // error: capturing closure cannot cross the abstract boundary
    return 0
})
```

Native class instances are zero-size — the backend removes them entirely.
No allocation, no runtime overhead.

---

### Testing

Test functions carry the `test` qualifier and are auto-discovered by the
test runner. Declare them in files tagged `build test`. The `test`
qualifier sits between the parameter list and `->`. Test functions may
declare no parameters.

```vertex
package arithmetic_test
build test

import "arithmetic"

func test_add()        test -> Expected(int32, "15") { return add(a: 10, b: 5) }
func test_comparison() test -> Expected(bool, "1")   { return 5 > 3 }
func test_no_crash()   test                          { add(a: 0, b: 0) }
```

`Expected(type, string)` declares the return type and the exact stdout
string the runner checks. Omitting `Expected` means the test passes as
long as it does not crash.

**Return value format reference:**

| Type | Format | `Expected` for value `5` |
| --- | --- | --- |
| `int32` | `%d` | `Expected(int32, "5")` |
| `int64` | `%lld` | `Expected(int64, "5")` |
| `uint32` | `%u` | `Expected(uint32, "5")` |
| `float32` | `%f` | `Expected(float32, "5.000000")` |
| `bool` | `%d` | `Expected(bool, "1")` / `"0"` |
| `string` | `%s` | `Expected(string, "hello")` |

A test can also assert that a line fails to compile:

```vertex
func test_bad_add() test -> Expected(error) {
    return add(a: 10, b: "5")
}

func test_bad_cast() test -> Expected(error, "cannot convert string to int32") {
    let x: int32 = "hello" as int32
}
```

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
  -emit-mir              emit Machine IR text (.mir)
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

- [Foundation](https://github.com/vertex-language/spec/foundation.md) — literals, control flow, functions, structs, enums, tuples, error handling
- [Ownership](https://github.com/vertex-language/spec/ownership.md) — shared/`mut`/`var` conventions, `unique`/`shared`/`weak`
- [Generics](https://github.com/vertex-language/spec/generics.md) — type parameters, constraints, monomorphization
- [Memory](https://github.com/vertex-language/spec/memory.md) — `typed_ptr`, manual allocation
- [Concurrency](https://github.com/vertex-language/spec/concurrency.md) — `thread`/`async`, channels, `select`
- [Hardware Acceleration](https://github.com/vertex-language/spec/accel.md) — `gpu`/`tpu`, tensors
- [Abstract Interfaces](https://github.com/vertex-language/spec/abstract_interfaces.md) — native interop

---

## License

MIT — see [LICENSE](https://github.com/vertex-language/vertex/LICENSE)