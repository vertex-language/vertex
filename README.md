# Vertex

**Language Spec 2.2 · Compiler 0.4.0**

Vertex is a statically-typed, compiled, multi-platform programming language built
for explicit control and zero-overhead interop. Nothing crosses a boundary
implicitly: absence is a union you spell, failure is a return type you declare,
every change of numeric representation is a call you write, and a condition is a
`bool` or it is an error.

One grammar spans seven platforms, and a single `use` line at the top of a file
states which one it expects — who owns memory, which numerics resolve, how a
foreign name is looked up. The build checks the claim. A file that makes no claim
is portable by construction and compiles everywhere unchanged. Four native targets
add machine access: the full pointer family, layout control, deterministic
destructors, and foreign declarations taken on trust rather than parsed from a
header. Two host targets hand memory to a runtime and drop the pointer family in
exchange for reach into a foreign object graph.

Its distinguishing feature is that the target is a property of the file rather
than of the build system, stated in one line and checked rather than inferred.

The normative reference is the **Grammar Summary**. This README is a tour.

---

## Contents

- [Install](#install)
- [Quick Start](#quick-start)
- [Language Tour](#language-tour)
  - [Source File Layout](#source-file-layout)
  - [The `use` Line](#the-use-line)
  - [Variables and Types](#variables-and-types)
  - [Numerics and Conversion](#numerics-and-conversion)
  - [Literals](#literals)
  - [Control Flow](#control-flow)
  - [Functions and Passing Modes](#functions-and-passing-modes)
  - [Structs and Classes](#structs-and-classes)
  - [Generics](#generics)
  - [Contiguous Storage](#contiguous-storage)
  - [Tuples](#tuples)
  - [Enums](#enums)
  - [Absence and Failure](#absence-and-failure)
  - [Pointers and Manual Memory](#pointers-and-manual-memory)
  - [Foreign Declarations](#foreign-declarations)
  - [Host Targets](#host-targets)
  - [Kernel and Graph Functions](#kernel-and-graph-functions)
- [Operator Precedence](#operator-precedence)
- [What Is Absent](#what-is-absent)
- [Open Questions](#open-questions)
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
vertex version
# vertex 0.4.0 (spec 2.2)
```

---

## Quick Start

A portable file makes no claim about its target, so it needs no `use` line at all
and compiles under every build:

```vertex
namespace math

export func fib(n: int): int {
  if n <= 1 {
    return n
  }

  var a = 0
  var b = 1
  var i = 2
  while i <= n {
    let next = a + b
    a = b
    b = next
    i += 1
  }
  return b
}
```

A file that wants machine access says so, and gets sized numerics, pointers, and
a foreign surface in exchange:

```vertex
namespace main

use native
use linux

declare module "c" {
  export func write(fd: int32, buf: const_ptr<byte>, n: usize): int64
}

func fib(n: int32): int32 {
  var a: int32 = 0
  var b: int32 = 1
  var i: int32 = 0

  while i < n {
    let next = a + b
    a = b
    b = next
    i += 1
  }
  return a
}

// There is no bridge from `string` to `const_ptr<byte>`, so anything crossing to
// a C API is bytes assembled by hand. ASCII codes are spelled numerically —
// Vertex has no character literal, since `'A'` is a string.
func putLine(v: int32): void {
  var digits: array<byte, 12> = array<byte, 12>()
  var n: usize = 0
  var x = v

  if x === 0 {
    digits[0] = 48
    n = 1
  }

  while x > 0 {
    digits[n] = uint8(48 + (x % 10))
    x /= 10
    n += 1
  }

  var out: array<byte, 13> = array<byte, 13>()
  for i in 0..n {
    out[i] = digits[n - 1 - i]
  }
  out[n] = 10                    // '\n'

  write(1, addressof(out[0]) as const_ptr<byte>, n + 1)
}

func main(): int32 {
  for i in 0..10 {
    putLine(fib(i))
  }
  return 0
}
```

`declare module` is **both the declaration and the binder** — `write` enters this
file's scope directly and is called unqualified. There is no named-binding import
form anywhere in the language, which is what makes a `declare module` block the
only route by which a foreign name can enter a file.

```sh
vertex fib.vs

# or compile
vertex -o fib fib.vs
./fib
```

---

## Language Tour

### Source File Layout

A compilation unit is a namespace header, then `use` directives, then top-level
items — in that order, fixed by the grammar rather than by a directive-region rule.

```vertex
namespace app.net

use native
use linux

import "core/io"
import (
    "app/config"
    sysio "std/io"
    _     "app/drivers/sqlite"
)
```

- `namespace` is mandatory and file-scoped. It is not a declaration: it cannot
  nest, cannot be re-opened, and a dotted name names one namespace rather than a
  chain of implicitly exported ones.
- Imports are Go-form and bind **Vertex modules by path only**. An `ImportSpec` is
  an optional alias followed by a string; `_` is the blank alias. There is no
  `from` clause, no named bindings, no default import, no `import =`.
- `export Declaration` is the only export form. `export default`, `export =`, and
  named export lists are gone, because only the deleted import forms could have
  consumed them.

**Statement termination.** There is no ASI. A line break inserts a terminator when
two things hold: the preceding token can end a statement, and the innermost
unclosed bracket is a statement-or-member container (a block, class body, struct
body, enum body, object type, ambient module body, or the top level). A `(` or `[`
is not such a container, so a line inside an argument list or an array literal
continues freely, and a *new* line beginning with `(` or `[` can never silently
glue itself to the previous statement. `;` remains legal everywhere and is
required nowhere.

---

### The `use` Line

Four slots, in order. Every one is optional; every one present is a claim the
build checks.

| Slot | Values | Optional because | Meaningful under |
|---|---|---|---|
| Memory model | `native`, `host` | implied by the platform | any platform but `any` |
| Platform | `any`, `windows`, `linux`, `darwin`, `wasm`, `android`, `js` | absent means `any` | — |
| Runtime | `nostd`, `noentry` | — | native platforms only |
| Accelerated backend | `cuda`, `msl`, `stablehlo` | — | native platforms only |

```vertex
use android                 // platform alone — host ownership is implied

use host
use android                 // identical, and says so

use native
use linux
use nostd
use noentry                 // a kernel module or bootloader

use linux
use cuda                    // gates kernel func

use any                     // states portability explicitly
```

The platform line is the load-bearing one: it fixes the numeric table and decides
how a bare `declare module` specifier resolves. Every platform name is unique
across the set — `android` is only ever host-owned, `windows` only ever manual — so
the memory model is already known by the time you've read the platform line. The
memory model line never changes what compiles; it exists so a file can *say* what
it is, and so the type system's two independent questions get two places to be
written.

| | `any` | `android` | `js` | `windows` | `linux` | `darwin` | `wasm` |
|---|---|---|---|---|---|---|---|
| Implies | — | `host` | `host` | `native` | `native` | `native` | `native` |
| `int` | only numeric | boundary-restricted¹ | only numeric | valid | valid | valid | valid |
| Sized numerics | invalid | mandatory | invalid | valid | valid | valid | valid |
| Pointer family | no | `weak_ptr` only | no | full | full | full | full |
| `destructor` | invalid | invalid | invalid | valid | valid | valid | valid |
| Foreign surface | none | classpath | host bundler | library path | library path, `dynamic:` | library + `framework:` | import namespace |

¹ Open — see [Open Questions](#open-questions).

Three errors are internal to a file and checkable before the build's target is
known: a memory model line with no platform, a memory model line contradicting its
platform, and `noentry` without `nostd`. Everything else is checked against the
build: a line that matches passes, a line that contradicts errors. `use any` is
the one line that always passes, since "this file asserts no target-specific
types" is true under every target.

---

### Variables and Types

`var` and `let` split on **mutability, not scope** — Swift's rule, not
JavaScript's. Both are block-scoped; neither hoists. `const` is reserved for
compile-time-only values, which is why three heads survive where two would do.

```vertex
var x = 10
x = 20                  // fine

let y = 10
// y = 20               // error — let is immutable

const LIMIT: int32 = 4096
```

`number` and `boolean` are renamed `int` and `bool`. Both are `PredefinedType`
entries; every sized numeric name is an ordinary type identifier, which is exactly
what lets `int32(a)` parse as a call rather than requiring a conversion production.

---

### Numerics and Conversion

Numerics are the one part of the type system that is not uniform across platforms.
The grammar for writing `int32` is identical everywhere; whether it *resolves*
depends on the platform line. The two host targets sit at opposite ends for the
same underlying reason — the target's own representation decides. A JS engine has
one numeric representation, so a sized type would assert a width nothing
distinguishes. A JVM descriptor demands an exact width, so an unsized type is a
claim it cannot encode.

```vertex
int  int8  int16  int32  int64
     uint8 uint16 uint32 uint64
     usize float32 float64
     byte          // alias for uint8 — says "raw storage", not "small number"
     bool
```

`usize` is pointer-width: 64-bit on `windows`/`linux`/`darwin`, **32-bit on
`wasm`**, and `long` on `android`. Code that treats `usize` and `uint64` as
interchangeable is correct on three platforms and wrong on two.

**No truthiness, in either direction.** `bool` is one byte with two valid patterns,
is not a numeric type, and has no conversion call. An integer becomes a boolean
through a comparison; a boolean becomes an integer through a conditional.

```vertex
let a: int32 = 0
let b: bool  = a !== 0
let c: int32 = b ? 1 : 0
```

**Every change of numeric representation is an explicit call** — no implicit
narrowing and no implicit widening either, because a rule with an exception is
harder to hold than a rule without one.

```vertex
let wide = int64(a)         // conversion call — may change the bits
let same = a as uint8       // type assertion — static only, never changes bits
```

For reinterpretation that *is* a representation change, `bit_cast<T>` is the
native-only route.

---

### Literals

```vertex
let a = 3_000_000_000
let b = 0xffff_ffff
let c = 0b1010_0101
let d = 0o755
let e = 0.0                 // `1.` is not a literal — fractional digits are mandatory
let f = -9_000_000_000
let s = "vertex"
let t = `multi
line`
```

The mandatory fractional digit is the one lexical cost of the range operator:
without it `0..10` would lex as the two literals `0.` and `.10`. Underscores are
separators with no grouping convention, so `1_000_000` and `10_00_000` are the same
literal.

There is **no character literal**. `'A'` is a single-quoted string, as in
TypeScript; byte values crossing to a C API are spelled numerically.

---

### Control Flow

Parens are dropped from every head and the body is a mandatory block. `if x doA()`
is not derivable, and the dangling-`else` ambiguity disappears with it.

```vertex
if x > 10 {
  doA()
} else if x > 0 {
  doB()
} else {
  doC()
}

while x < 10 {
  x += 1
}
```

`if let` binds one identifier — no patterns — and is the ordinary unwrap for every
nullable union in the language:

```vertex
if let value = maybeValue {
  use(value)
} else {
  return
}
```

`for` has one shape. `for`-`in`, `for`-`of`, `for await`, and the C-style triple all
collapse into it; the iterated expression may be a range, which is the whole reason
`..` exists:

```vertex
for i in 0..10 { }                 // half-open range
for item in items { }
for (name, score) in scores { }    // the binding reaches TupleBindingPattern
```

`switch` takes comma lists instead of grouped fallthrough, and `default` is
structurally last rather than positionally free:

```vertex
switch code {
  case 0 {
    doZero()
  }
  case 1, 2, 3 {
    doSmall()
  }
  default {
    doOther()
  }
}
```

Heads are `ConditionExpression` — an expression with no unparenthesized object
literal at depth 0. That is the price of paren-free heads, and it is the same price
Go and Rust pay.

---

### Functions and Passing Modes

```vertex
func add(a: int32, b: int32): int32 {
  return a + b
}

func divide(a: int, b: int): (int, int) {
  return (a / b, a % b)
}

let (quotient, remainder) = divide(17, 5)
```

A body-less function signature — a terminator where the block would be — is an
overload declaration. Decorators are admitted on function declarations as well as
the inherited class positions.

**Passing modes** are prefix type operators that pass a value type without copying
it. An unannotated parameter copies.

```vertex
func scale(v: mutating Vec2, k: float32): void {
  v.x *= k
}

func length(v: readonly Vec2): float32 {
  return v.x
}
```

`readonly` and an unannotated parameter mean the same thing, minus the memcpy;
`mutating` is the only mode that changes what the program means, and the only route
to mutating a caller's value type. Both bind in parameter, return, and
`this`-parameter position — never a local, never a field, never a type argument.

- **Value types only.** A `class` binding is already a reference. Pointer types never
  take a mode; they pass a handle by copy.
- **Shallow.** `readonly Vec` freezes the struct's own fields; a `mutable_ptr<T>`
  member stays writable through.
- **Non-escaping.** A borrow may be passed downward, never stored.
  `constructor(private x: mutating T)` is an error.
- **No exclusivity guarantee, no call-site marker.** Two `mutating` parameters may
  alias; `f(a)` does not signal that `a` is rewritten.

---

### Structs and Classes

A `struct` is a value type with fixed layout, copied by default, constructed by
direct call. Its body is a deliberate subset of a class body: no `extends`, no
static blocks, no accessors, no initializers, and a **mandatory type annotation**
on every field — a value type with a fixed layout has no inferred field types.

```vertex
struct Vec2 {
  x: float32
  y: float32

  constructor(x: float32, y: float32) {
    this.x = x
    this.y = y
  }

  scaled(this: readonly Vec2, k: float32): Vec2 {
    return Vec2(this.x * k, this.y * k)
  }
}
```

Layout control is spelled with decorators, which the language needs for function
attributes and bitfields regardless:

```vertex
@packed struct Address {
  ull: uint64
}

@align(64) struct CacheLine {
  storage: array<byte, 64>
}

struct ClassOfDevice {
  @bits(2)  format: uint32
  @bits(6)  minor: uint32
  @bits(5)  major: uint32
  @bits(11) service: uint32
  @bits(8)  reserved: uint32
}
```

Bitfields require an unsigned sized type and are valid only inside a struct; the
packed word is reached with `bit_cast` rather than by shifting at every site.

A `class` is a reference type, heap-allocated with an inline refcount header, and
ARC'd invisibly under the native model. Its distinguishing feature is deterministic
teardown:

```vertex
class Socket {
  readonly fd: int32

  constructor(fd: int32) {
    this.fd = fd
  }

  destructor() {
    close(this.fd)
  }
}
```

That destructor is why an error path is a bare `return` instead of ten lines of
cleanup: every exit past the constructor closes the descriptor at a known point.
`destructor` takes no parameters, no modifiers, no overloads, and a mandatory
block — there is nothing to overload, so there is no signature-only form.

**Classes have no inheritance.** `implements` is the sole route to polymorphism.
The one exemption is extending a *foreign* class, which exists because `Activity`,
`Service`, `View`, and `NSWindow` are subclass-or-nothing APIs.

---

### Generics

```vertex
func clamp<T extends Ordered>(a: T): T {
  return a
}

struct Buffer<T, const N: usize> {
  private storage: array<T, N>
}
```

A const generic parameter binds a compile-time value as a type parameter, usable
directly in value position inside the body. It is distinguished from TypeScript's
`<const T>` modifier by the presence of `:` — one token of lookahead.

Instantiation is written `make_shared<Sample>(0, 0)`, `sizeof<Header>()`,
`bit_cast<uint32>(x)`. Note that `>>` is not a token anywhere in Vertex, precisely
so `block<span<int32>>` and `alignof<mutable_ptr<T>>()` are ordinary.

---

### Contiguous Storage

Postfix `T[]` is removed. There are four forms, and each says something different:

```vertex
let a: span<int32>                        // non-owning view
let b: array<byte, 16>                    // inline, fixed-length
var c: vector<byte> = vector<byte>()      // heap-backed, growable, owned
let d: block<int32> = alloc<int32>(usize(256))   // heap-sized-once, move-only, owned
```

Owned forms release at scope exit of the owner. Growth and allocation come in
panicking and fallible pairs throughout:

```vertex
a.push(1)
let ok: bool = a.try_push(1)

let e: block<int32> | null = try_alloc<int32>(usize(256))
```

`vector<T>` and anything depending on `block<T>`'s heap allocation is invalid or
conditional under `nostd`.

---

### Tuples

Two spellings exist for a tuple type, `[A, B]` and `(A, B)`; the parenthesized form
is the one the corpus uses. It takes two or more elements, which keeps it disjoint
from a parenthesized type.

```vertex
func listenOn(port: uint16): (Socket | null, int32) { }

let (server, err) = listenOn(uint16(8080))
```

Declaration-position destructuring matches the type element for element.
Assignment-position tuple destructuring is not derivable — `(a, b) = f()` is the
parenthesized comma expression, and assigning to it is not a form.

---

### Enums

Enums extend TypeScript's rather than replacing them: the backing type arrives
through the same `:` annotation used everywhere else, and an associated-value case
reuses a parameter list verbatim.

```vertex
enum StatusCode: int32 {
  OK = 200
  NotFound = 404
  ServerError = 500
}

enum Shape {
  Circle(radius: float64)
  Rectangle(width: float64, height: float64)
  Point
}

enum Result<T, E> {
  Ok(value: T)
  Err(error: E)
}
```

Members are terminator-separated rather than comma-separated. `const enum` is gone —
`const` is compile-time-only and an enum is already a compile-time construct.

---

### Absence and Failure

Vertex has no unwinder and does not throw. Two rules cover every boundary in the
language, and both are *told* rather than inferred.

**Absence is an explicit union**, unwrapped with `if let`. Bindings are
non-nullable by default even though C, C++, Objective-C, Java, and JS references
are all nullable by default on their own side:

```vertex
declare module "dynamic:c" {
  export func find(a: const_ptr<byte>): mutable_ptr<Entry> | null
}

if let e = find(name) {
  use(e)
}
```

**Failure is a return union.** A foreign method that unwinds is declared with its
failure in the return type; a checked failure left out of the union that fires
anyway panics rather than silently dropping. `instanceof` narrows the union —
deliberately not `if let`, since one is a value of another type and the other is
absence.

```vertex
declare module "java.io" {
  export declare class Reader {
    read(a: string): Buffer | IOException
  }
}
```

C has no exceptions to model, so a C-shaped failure is a sentinel plus a tuple, not
a union — there is nothing unwinding for a union to describe:

```vertex
func listenOn(port: uint16): (Socket | null, int32) {
  let fd = socket(AF_INET, SOCK_STREAM, 0)
  if fd < 0 {
    return (null, errno())
  }

  let s = Socket(fd)
  // ...every return past this point closes fd via the destructor
  return (s, 0)
}
```

---

### Pointers and Manual Memory

Under `native`, ordinary construction is still ARC'd — `var a = Sample()` allocates
and refcounts invisibly. Native is *optionally* manual, not unmanaged; reaching for
the pointer family is a decision, not the price of admission.

```vertex
let a: shared_ptr<Sample>       let d: mutable_ptr<Sample>
let b: unique_ptr<Sample>       let e: const_ptr<Sample>
let c: weak_ptr<Sample>         let f: void_ptr

var g = make_unique<Sample>(0, 0)
var h = make_shared<Sample>(0, 0)
var i = weak_ptr(h)
if let live = i.lock() { }
```

`weak_ptr<T>` is only ever derived from an existing `shared_ptr<T>` and never
allocates. Raw pointers are non-nullable by default; absence is an explicit union.

| Category | Operations |
|---|---|
| Address | `addressof(a)` |
| Arithmetic | `offset` · `byte_offset` · `distance` · `byte_distance` · `align_up` · `align_down` · `is_aligned` |
| Access | `a[0]` is the dereference — pointer methods are not shadowed by pointee ones |
| Allocation | `alloc<T>` · `try_alloc<T>` · `alloc_uninit<T>` |
| Placement | `construct_at(p, …)` · `destroy_at(p)` — the only sanctioned hand-call of a destructor |
| Casts | `bit_cast<T>` · `pointer_cast<T>` · `pointer_from_address<T>` |
| Explicit access | `unaligned_load/store<T>` · `volatile_load/store<T>` |
| Layout | `sizeof<T>()` · `alignof<T>()` · `offsetof<T>("x")` |

---

### Foreign Declarations

One grammar covers C, C++, Objective-C, JVM packages, wasm imports, and JS modules.
`declare struct` introduces a layout-free type whose definition lives outside
Vertex, valid only behind a pointer; `declare module` is both the declaration and
the binder.

```vertex
declare struct llama_model

declare module "llama" {
  export func llama_model_get_vocab(m: const_ptr<llama_model>): const_ptr<llama_vocab>
  export func llama_decode(c: mutable_ptr<llama_context>, b: llama_batch): int32
}
```

**Trust, not verification.** The compiler never reads a header, class file, or
framework definition. A signature that disagrees with the real symbol fails at link
time, at instantiation, at first use, or at the call site — and *which* is a
property of the platform, not of the declaration:

| Platform | Specifier means | Bound at | Schemes |
|---|---|---|---|
| `windows` | library search path | link time | `dynamic:` |
| `linux` | library search path | link time | `dynamic:`, `cpp:` |
| `darwin` | library **or** framework path | link time | `framework:`, `dynamic:`, `cpp:` |
| `wasm` | import module name | instantiation | none |
| `android` | classpath package | first use | none |
| `js` | whatever the host resolver does | build/bundle time | none |

Scheme prefixes live inside the string literal and cost zero grammar; an unknown
scheme is a resolution error, not a parse error.

**`darwin` is the one platform where a bare specifier is genuinely ambiguous**, so
the specifier names the resolver directly. A library-resolved block may contain only
flat function declarations; a framework-resolved block may contain
`class`/`interface` for Objective-C dispatch, or boundary C structs. Mixing is
disallowed by construction rather than by a rule to remember — one block resolves
through exactly one path.

```vertex
declare module "System" {              // bare = library path, -lSystem
  func sample(): int32
}

declare module "framework:WebKit" {    // framework path, -framework WebKit
  class WKWebView { }
}
```

The `.framework` extension is never spelled; the resolver appends it.

**C++** rides in on `cpp:`, which answers *how a name is spelled* at the ABI rather
than *when* it resolves. Mangling is platform-owned — Itanium on `linux`/`darwin`,
MSVC decoration on `windows` — and a C++ namespace rides in the specifier string.
Templates monomorphize, matching Vertex generics, so there is no erasure tax;
`unique_ptr<T>`/`shared_ptr<T>` are the same layout rather than an analogy, and
`T&`/`const T&` map onto `mutating`/`readonly`.

**Rest parameters are call-shape, not a collection.** C varargs land in registers or
on the stack per the ABI; nothing inside an extern declaration can index or iterate
them.

```vertex
declare module "c" {
  export func printf(fmt: const_ptr<byte>, ...args: CVarArg): int32
}
```

**Selectors and descriptors are ordinary string literals.** An Objective-C selector
send and a JVM overload disambiguated by descriptor are the same grammar — the
string-literal arm of a property name, plus a computed member call. No production
distinguishes a selector from a descriptor from an ordinary name; the resolver does,
from the enclosing block's scheme.

```vertex
a["sample002(Ljava/lang/String;)V"](b)
a["webView:didFinishNavigation:"](c, d)
```

---

### Host Targets

`android` and `js` both imply `host`: a runtime owns memory, so there is no pointer
family, no manual allocation, no layout control, no `bit_cast`, no `sizeof`, and no
`destructor`. Teardown in Vertex is deterministic and scope-bound, while ART's
`Cleaner` and JS's `FinalizationRegistry` both run at an unspecified time relative
to collection — the two cannot be reconciled by declaring one in terms of the other.
Cleanup for foreign resources is written by hand on every exit path.

The two targets are mirror images on numerics: `android` **mandates** sized types
because descriptors are exact; `js` **forbids** them because a JS engine has one
numeric representation.

One pointer-family name survives, and only under `android`:

```vertex
var strong: Activity = Activity()
var weak: weak_ptr<Activity> = weak_ptr(strong)

if let live = weak.lock() {
  live.doSomething()
}
```

That is not a memory-model feature. A long-lived callback or static holder capturing
`this` is a real root, not a cycle, and a tracing collector is correct to keep the
object alive — so the problem is visibility, not ownership. `js` has no equivalent
binding today, because the event loop's long-lived roots are usually fixed by
unregistering rather than by weakening.

On `android`, every class is a class file, which is why there is no decorator to
mark one dispatchable from the host runtime — `use android` already says it, one
line earlier and once per file. `darwin`'s `@objc` survives for the opposite reason:
there are genuinely two lowerings there, and the platform line cannot pick.

---

### Kernel and Graph Functions

A `kernel func` or `graph func` body is valid only in a file that names an
accelerated backend, and the backend decides which intrinsics resolve inside it.

```vertex
namespace kernels

use linux
use cuda

export kernel func saxpy(a: device_ptr<float32>, n: int32): void {
  let i = threadIdx.x
}
```

- **Kernel** — SIMT, side-effecting, lowers to PTX or MSL. `device_ptr<T>` and the
  thread-context intrinsics (`threadIdx`, `blockIdx`, `blockDim`, `gridDim`) are
  valid only inside one. Kernels are nameable but not callable by host code; they go
  through `compile` and `launch`.
- **Graph** — pure dataflow over whole tensors, no thread context, lowers to a
  StableHLO string. Callable directly once compiled, no launch configuration.

```vertex
let compiled = compile(saxpy)
launch(config, compiled, a, b, n)
```

`simd<T, N>` is trivially copyable and needs no modifier; `tensor<T, ...dims>` is
shape-encoded and used exclusively on the graph route. No backend applies under
`android`, `js`, or `wasm`.

---

## Compiler Reference

`vertex` compiles through one pipeline, orchestrated by `driver` and ending
either in vvm's own IR encoders or in vvm's own native-image builder — this
package makes the decisions, vvm does the verifying, lowering, and linking:

```
<file.vs | dir>
     │  parser + analyzer (via importer for a package,
     │           directly for a single file)
     ▼
[]*driver.Package  (checked, dependency-first)
     │  lower/hir.Lower  — every decision made here
     ▼
*hir.Program
     │  lower/vir.Lower  — mechanical
     ▼
[]*vir.Module  (unverified; ir/verify is vvm's to run)
     │
     ├─ -emit-vir / -emit-vbyte ── format/vbyte/{text,binary}.Encode ─► files
     └─ (default)      ────────── vvm.BuildModule / BuildModuleGraph ─► image
```

There is no separately-exposed Machine IR, assembly, or object-file stage:
those are internal to vvm's own build pipeline, and there is no optimizer in
either tree to gate with an `-O` flag.

```text
Usage:
  vertex [build] [flags] <file.vs | package-dir>
      Compile a package to a native executable for the host or -target.
      The command word is optional: "vertex main.vs" is "vertex build main.vs".

  vertex run [flags] <file.vs | package-dir> [-- args...]
      Build for the host, execute immediately, forward the exit code.
      Anything after -- is passed to the compiled program, not to vertex.

  vertex test [-dir <path> | -file <path>] [-run <substr>]
      Discover 'test'-marked functions in a test package and run them,
      comparing each against its expected result.

  vertex targets
      List every target this toolchain can actually build for.

  vertex version | help

Build flags:
  -o <path>             output file (default: derived from the input name).
                          For -emit-vir/-emit-vbyte on a multi-package build,
                          this must be a directory — one file per package.
  -target <triple>      see "vertex targets"; defaults to the host
  -min-os-version <v>   required for darwin targets; defaults to 14.0
  -shared               build a shared library instead of an executable
  -flat-base <addr>     load address for a freestanding flat image
  -root <module>        override which module's entry function is the
                          program's entry point (default: the package
                          holding main)
  -packages-dir <path>  packages root; overrides $VERTEX_PATH
  -emit-vir             emit Vertex IR text (.vir), one file per package
  -emit-vbyte           emit Vertex IR binary (.vbyte), one file per package
  -v                    report each pipeline stage on stderr

Run flags:
  -packages-dir <path>  packages root; overrides $VERTEX_PATH
  -v                    report each pipeline stage on stderr

Test flags:
  -dir <path>           directory holding test files (default: .)
  -file <path>          a single test file
  -run <substr>         only run tests whose name contains this substring
  -packages-dir <path>  packages root; overrides $VERTEX_PATH
  -v                    print every test, not just failures
```

**Examples:**

```sh
vertex main.vs
vertex -o app ./cmd/app
vertex build -target darwin-arm64 -o app main.vs
vertex run main.vs -- --verbose
vertex build -emit-vir -o build/ ./cmd/app
vertex test -dir ./tests
vertex test -file literals_test.vs -run comparison
vertex targets
vertex version
```

| Command | Output | Use case |
| --- | --- | --- |
| `vertex` / `vertex build` | executable (or `-shared` library) | the default: fully linked native binary |
| `vertex build -emit-vir` | `.vir` per package | human-readable Vertex IR; inspect lowering |
| `vertex build -emit-vbyte` | `.vbyte` per package | binary Vertex IR |
| `vertex run` | *(runs immediately)* | build for the host and execute in one step |
| `vertex test` | *console* | discover and run test functions |
| `vertex targets` | *console* | list every buildable target, marking the host |

`$VERTEX_PATH` sets the packages root; `-packages-dir` overrides it. When
neither is set, the compiler defaults to `~/.vertex/packages`.

A `-target` must agree with every `use` line in the files it builds. A platform line
that contradicts the resolved target is a compile error, not a filter — and a file
carrying *any* `use` line cannot compile at all until the build has a resolved
target, since there is nothing to check the claim against. A file with no `use` line
made no claim and rides along under whatever profile the build resolves. An
unrecognized target name is a compile error naming the known set, not a silent
fallback.

---

## Platform Support

Every target below has a `cpu/lower`, an `object` writer, and either a
registered linker or a flat writer in vvm — `vertex targets` reports the same
list at runtime, marking the host with `*`.

| Target | `vertex build` (executable) | `-shared` | `-emit-vir` / `-emit-vbyte` |
| --- | --- | --- | --- |
| `linux-amd64` | yes | yes | yes |
| `linux-arm64` | yes | yes | yes |
| `darwin-amd64` | yes | yes | yes |
| `darwin-arm64` | yes | yes | yes |
| `windows-amd64` | yes | yes | yes |
| `windows-arm64` | yes | yes | yes |
| `freestanding-amd64` | flat image only; single package, entry must be `_start` | no — flat has no loader | yes |
| `freestanding-arm64` | flat image only; single package, entry must be `_start` | no — flat has no loader | yes |

The freestanding targets are the `use native` + `use nostd` + `use noentry` case: a
kernel module and a bootloader are both spelled that way, neither needs its own
platform name, and both inherit native's pointer and sizing rules unchanged.

`linux-riscv64` and every powerpc/mips/loongarch/s390x spelling are valid
`.vir` target triples per the language spec but have no `cpu/lower`, object,
or linker implementation in vvm, so they don't appear above. `linux-386` and
`windows-386` are similarly absent: vvm emits x86 ELF object bytes but
registers no ELF linker backend for x86. The `wasm`, `js`, and `android` platform
lines are specified in full but have no backend yet.

---

## License

MIT — see [LICENSE](https://github.com/vertex-language/vertex/LICENSE)