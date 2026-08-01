# Vertex

**Language Spec 2.2 · Compiler 0.4.0**

Vertex is a statically-typed systems and application programming language built
for explicit control, zero-overhead interop, and first-class concurrency. Every
value has a statically known layout and every cost is decided at compile time:
no garbage collector, no unwinder, no vtables, no runtime type information, and
no hidden allocation.

Its syntax draws from Swift and Go. Its distinguishing features are a layered
ownership model with a single visible copy/move marker, and execution substrates
(`async`, `gpu`, `npu`) declared as function markers and spelled again at the
call site.

The normative reference is **Annex A — Grammar Summary**. This README is a tour.

---

## Contents

- [Install](#install)
- [Quick Start](#quick-start)
- [Language Tour](#language-tour)
  - [Source File Layout](#source-file-layout)
  - [Variables and Types](#variables-and-types)
  - [Literals](#literals)
  - [Control Flow](#control-flow)
  - [Ownership and Functions](#ownership-and-functions)
  - [Structs and Classes](#structs-and-classes)
  - [Generics](#generics)
  - [Arrays, Slices, and Maps](#arrays-slices-and-maps)
  - [Tuples](#tuples)
  - [Enums](#enums)
  - [Error Handling](#error-handling)
  - [Raw Pointers and Builtins](#raw-pointers-and-builtins)
  - [Concurrency](#concurrency)
  - [Device Offload](#device-offload)
  - [Abstract Interfaces](#abstract-interfaces)
  - [Testing](#testing)
- [Operator Precedence](#operator-precedence)
- [What Is Absent](#what-is-absent)
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

```vertex
package main
build linux

declare module "c" {
    func printf(fmt: string, args: ... int32) -> int32
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
    while i <= n {
        var tmp: int32 = a + b
        a = b
        b = tmp
        i += 1
    }
    return b
}

func main() {
    printf("--- recursive ---\n")
    for i in 0..10 {
        printf("fib(%d) = %d\n", i, fibRecursive(i))
    }

    printf("--- iterative ---\n")
    for i in 0..10 {
        printf("fib(%d) = %d\n", i, fibIterative(i))
    }
}
```

`declare module` is a **linkage boundary, not a namespace** — `printf` is
injected into this file's package and called unqualified. See
[Abstract Interfaces](#abstract-interfaces).

`main` takes no parameters and returns nothing.

```sh
vertex fib.vs

# or compile
vertex -o fib fib.vs
./fib
```

---

## Language Tour

### Source File Layout

A source file is a package clause, an optional build clause, imports, then
declarations — in that order.

```vertex
package net

build linux

import "core/io"
import (
    "core/time"
    "app/config"
)
```

- `package` is mandatory and must be the first non-comment token sequence.
- `build` selects the target platform and, with it, the ABI family used by every
  `declare` block in the file. It is not a preprocessor directive: a file whose
  tag doesn't match the target is excluded whole, never partially. Tags are
  `linux`, `windows`, `darwin`, `js`, `wasm`, `test`.
- Imports name a *locator*. The qualifier you write in code is the imported
  package's own declared name — `time.now()`, not something derived from the
  path. There is no aliasing, no dot-import, and no blank import.
- Top-level declarations are order-independent. There is no forward-declaration
  form because none is needed.
- A top-level `var` must have a compile-time-evaluable initializer. There is no
  static initialization order and no initialization-time code.

**Statement termination.** There is no terminator token and no semicolon. A
statement ends at a line break or at the `}` closing its block. There is no
automatic-semicolon-insertion machinery and no continuation rule: a line break
inside `(…)`, `[…]`, or `{…}` is ordinary whitespace, and a line break outside
one ends the statement. Every multi-line construct is bracket-delimited by
construction.

---

### Variables and Types

`let` declares an immutable binding; `var` declares a mutable one. Neither says
anything about heap versus stack.

```vertex
let x: int32     = 42
var y: float32   = 1.5
let name: string = "vertex"
let flag: bool   = true

var count: int32          // no initializer — zero value; annotation required
```

`let` is immutable and **not guaranteed to be addressable** — it may live in a
register, be an SSA value, or fold away entirely. `var` is mutable and owns a
real stack slot for its whole lifetime. This is why only a `var` binding can be
passed to a `mut` parameter: a `let` may not physically exist anywhere to point
at.

Multi-binding forms destructure or declare in parallel:

```vertex
let q, r = divmod(10, 3)        // one initializer → tuple destructure
var a, b = 0, 1                 // matching counts → parallel declaration
let _, err = parseInt(s: "42")  // `_` discards
```

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
| `string` | UTF-8 bytes + length, **no NUL terminator** |

`int` is `int32`; `uint` is `uint32`; `byte` and `uint8` are the *same type* —
no conversion is required or permitted between them in either direction.
`char` is a 4-byte Unicode scalar, not a C byte.

Every scalar has a statically known width, copies by register move, and has an
all-zero zero value. There is no boxed form of any of them.

**Conversions are explicit and written with `as`:**

```vertex
let wide = x as int64          // sign-extended
let big  = x as uint64         // zero-extended

let f: float64 = 3.99
let i = f as int32             // truncates toward zero → 3

let chain = value as int32 as int64   // left-associative, two conversions
```

`as` never touches memory. Between pointer types it is a static
reinterpretation; between numeric types it is a width-selected
truncate/extend/int↔float instruction; on an enum it is a tag read. There is no
dynamic cast, because there is no runtime type information for one to consult.

`as` binds tighter than every binary operator, so
`count as float64 / total as float64` divides two already-converted values.

Predeclared type names are ordinary identifiers, not keywords — they can be
shadowed. Reserved builtin names (`new`, `addr`, `delete`, `copy`, `zero`,
`sizeof`, `alignof`, `reinterpret`, `resize`, `drop`, `upgrade`, `transfer`)
cannot be, and may not be declared as a member, method, or field name.

`_` is legal only as a type-parameter name, an enum-payload binding, or a
discarded destructuring target. It never introduces a usable binding. `$` is not
an identifier character.

---

### Literals

```vertex
let dec = 1_000_000
let hex = 0xDEAD_BEEF
let oct = 0o755
let bin = 0b1010_0101

let f1 = 1.5e-3
let f2 = 0x1.8p3        // hex float — the binary exponent is required
```

`_` is a digit separator with no value; it may not lead or trail a digit run and
may not be doubled. There is **no negative literal** — `-1000` is unary minus
applied to `1000`, folded at compile time.

An integer literal is untyped until it reaches a typed position, where it takes
that position's type. A literal that does not fit the destination is a compile
error, never a silent truncation.

```vertex
let s   = "hello\tworld\u{1F600}"
let raw = `
  Vertex 2.2
  raw and multi-line: no escapes, every line break is part of the value
`
let c: char = 'A'
```

`'A'` and `"A"` are different types and never interconvert implicitly. A string
carries no NUL terminator; one is manufactured only where a `string` crosses an
abstract-interface boundary.

Arithmetic traps on overflow. The `&`-prefixed forms wrap:

```vertex
let t = a + b     // traps
let w = a &+ b    // wraps; also &- and &*
```

---

### Control Flow

```vertex
if err != "" {
    return
} else if n > 0 {
    // ...
} else {
    // ...
}
```

The condition must be `bool` — there is no truthiness conversion anywhere in the
language. There is **no initializer clause** in the header; the error-checking
idiom is two statements, and its verbosity is intentional.

`while` is the only loop primitive besides `for`-in. There is no C-style `for`
and no do/while:

```vertex
while i < n {
    i += 1
}
```

`for` has exactly one shape and consumes an iterable — ranges, fixed arrays,
dynamic arrays, maps, and strings:

```vertex
for x in 0..10 { }              // ranges are half-open; there is no inclusive form
for item in items { }           // shared access
for i, v in items { }           // index and value
for k, v in scores { }          // key and value; map order is unspecified
```

The consuming marker sits on the **binding**, not the iterable: what moves is
each element, one per iteration, and the marker names what moves.

String iteration decodes UTF-8 into `char` scalars at variable stride; byte-level
iteration is a separate method and strides raw `uint8`. Neither allocates.

```vertex
switch d {
case .North, .South:
    // ...
case .East:
    fallthrough
case .West:
    // ...
default:
}
```

Cases do not fall through implicitly; `fallthrough` is explicit and must be the
last statement in its clause. At most one `default` per switch. A switch over a
unit-only enum with no `default` must be exhaustive. The discriminant is read
once — dense tags lower to a jump table, sparse ones to a compare chain.

```vertex
defer delete(buf)
```

`defer` takes a call and nothing else. Its arguments are evaluated at
registration; only the call is postponed. Deferred calls run in reverse
registration order on **every** exit edge — fall-through, `return`, `break`,
`continue`. Because there is no unwinder, "every exit edge" is a finite,
statically known set, and a `defer` costs exactly the call it defers.

There are no loop labels. A multi-level exit is written with an explicit flag or
an extracted function.

**Control-flow headers reject bare composite and map literals.** Parenthesize:

```vertex
if (Vec2{x: 1.0, y: 2.0}).isUnit() { }
```

Assignment is a statement, never an expression — there is no assignment inside a
condition and therefore no `=`/`==` confusion class.

---

### Ownership and Functions

```vertex
func add(a: int32, b: int32) -> int32 {
    return a + b
}

add(1, 2)
add(a: 1, b: 2)
```

Named and positional arguments may not be mixed in one call. Named arguments
resolve to positional order at compile time and leave no trace in the binary.

A parameter's convention is picked in the **signature**. The caller writes at
most one thing, and only for the owning case:

```vertex
func inspect(w: Widget)      // shared   — read-only view, always bare
func rename(w: mut Widget)   // exclusive — mutates the caller's binding
func archive(w: var Widget)  // owning    — copy or move, chosen at the call site
```

`mut` lowers to a pointer to the caller's slot, but the address is never written
by the caller:

```vertex
func increment(n: mut int32) {
    n += 1
}

var count: int32 = 0
increment(count)     // bare — the argument must be an addressable `var` binding
```

**The ownership marker.** In an owning position, a `var` prefix means *move* and
its absence means *copy*. One marker, two meanings, read by presence:

```vertex
var w = Widget{id: 1}

archive(w)         // COPY — deep copy, O(data); w survives
inspect(w)         // fine

archive(var w)     // TRANSFER — header only, O(1); w is dead after this line
```

The same pair governs ordinary bindings:

```vertex
let a = w          // COPY
let b = var w      // TRANSFER
```

There is no `.clone()` and no copy operator. Copying is what happens when you do
not write `var`. The marker takes a binding or a field path and nothing else —
it does not compose through arbitrary expressions:

```vertex
var w                       // ✗ transfer outside an owning position
if var w { }                // ✗ a control-flow header is not an owning position
let y = var pick(a, b)      // ✗ transfer requires a binding or field path
```

`mut` is unrelated to this rule: `mut` never takes ownership, so the
copy/transfer question never arises and its call sites are always bare.

**Owning positions** — where `var` is legal — are exactly:

- the right-hand side of a variable declaration or assignment;
- an argument to a `var`-typed parameter;
- an element of a tuple, array, map, or composite literal;
- a returned expression;
- the binding of a consuming `for` loop.

Anywhere else a `var` prefix is a compile error naming the position.

**Liveness** is tracked statically through control flow. Use after transfer is a
compile error — and so is use after a transfer that *may* have happened on some
path. A conditional transfer is rejected outright rather than resolved at
runtime, because the moment "was it transferred?" becomes a runtime question the
language would need drop flags. It forbids the question instead. A transfer
inside a loop body is rejected for the same reason.

**The Law of Exclusivity** — aliasing *or* mutation, never both — is enforced at
every call site by reading the callee's signature. Passing one binding as two
exclusive arguments is a compile error, as is reading a binding in the same call
that exclusively accesses it, as is overlap through a field path.

**Cost, at a glance:**

| Operation | Copies | Cost |
| --- | --- | --- |
| bare copy of an owning fat type | header + payload | O(data) |
| `var` transfer | header only | O(1) |
| non-owning view | two words | O(1) |
| bare copy of `unique T` | **deep** — walks the pointee | O(data) |
| copy of `shared T` | atomic increment | O(1) |

`unique T` is one word but its bare copy is deep. This is the language's one
cost cliff hidden behind a thin type, and it's exactly why the marker is visible
rather than inferred.

**Closures** capture by value at creation. Assigning to a captured binding inside
the body is a compile error — the write would land on a private copy, and the
language declines to compile the lie. Writeback is spelled by taking a `mut`
parameter and letting the caller thread the pointer.

```vertex
let double = func(x: int32) -> int32 { return x * 2 }
```

A non-capturing function value is one word — a bare code pointer. A capturing
closure is two, `{code, env}`.

---

### Structs and Classes

Structs and classes are **byte-for-byte identical in layout**. Both are
stack-resident value types by default. A `class` differs only in its
member/method model — initializers, teardown, receiver conventions, identity
comparison. Declaring something a `class` does not, by itself, put it on the
heap.

```vertex
struct Vec2 {
    x: float32,
    y: float32,
}

func (v: Vec2) length() -> float32 { }             // shared receiver
func (v: mut Vec2) scale(factor: float32) {        // exclusive receiver
    v.x *= factor
    v.y *= factor
}

var pos = Vec2{x: 1.0, y: 2.0}
pos.scale(factor: 2.0)      // bare — no `&`; `mut` is a signature fact
```

Fields are separated by commas; a line break between them is conventional but
not required. Trailing commas are valid. Struct literals require field labels.

Field defaults are evaluated at construction for any field the literal omits:

```vertex
struct Config {
    workers: int32 = 4,
    verbose: bool  = false,
}

let c = Config{}
```

Classes are constructed by calling the type name with the initializer's
arguments — **never** with a brace literal. The punctuation alone tells the
reader which storage discipline is in play:

```vertex
class Animal {
    name:   string,
    health: int32,
}

func (a: Animal) init(name: string, health: int32) {
    a.name   = name
    a.health = health
}

func (a: Animal) deinit() { }

let dog = Animal(name: "Rex", health: 100)
// dog.deinit() is emitted where dog's liveness ends — no defer needed
```

Class methods are declared outside the class body, as ordinary functions with a
receiver. The body holds fields only. There is no inheritance, no vtable, no
dynamic dispatch, and therefore no class header of any kind — every call is
direct.

A **transferred** binding simply has its teardown not emitted. No flag is set and
none is checked.

A receiver typed `var` consumes its receiver **unconditionally**: the receiver
position has no argument slot to carry a marker, so there is no bare form that
copies. Copy first into a fresh binding to call one non-destructively. This is
the single exception to the bare-means-copy rule.

`===` and `!==` compare storage identity and are legal on classes only. They
answer "same allocation?", never "same bytes?" — that's `==`'s question.

**The heap has exactly two doors:**

```vertex
var u  = unique(Animal(name: "Rex", health: 100))    // sole ownership, no refcount
var s  = shared(Animal(name: "Luna", health: 100))   // refcounted handle
var s2 = s                                            // cheap increment, not a deep copy
var wk = weak(s)                                      // observes without keeping alive

let observer, err = upgrade(wk)
if err != "" {
    // the shared value is gone
}
```

Both still tear down automatically — `unique` at the end of its owner's liveness,
`shared` once the strong count reaches zero. Neither needs a delete call.

`weak T` observes only a `shared` allocation. There is no `unique T` → `weak T`
path, because a `unique` block carries no control word for a weak reference to
inspect.

Qualifiers do not stack: `mut var T`, `mut mut T`, and `shared unique T` are
errors. `mut shared T` and `var shared T` are fine — the qualifier applies to the
handle, which is itself an ordinary value.

---

### Generics

Type parameters use square brackets, matching type-argument position:

```vertex
func identity[T](value: T) -> T {
    return value
}

struct Box[T] {
    value: T,
}

func (b: Box[T]) get() -> T {
    return b.value
}

let b = Box[int32]{value: 42}
let r = identity(value: "hello")     // T inferred from the argument
```

A bare name is constraint `any` — `[T]` means `[T: any]`. Under `any`, only
assignment, argument passing, and the ownership operations are available: no
comparison, no arithmetic, no field access.

A constraint written after a name applies to that name **and to every
immediately preceding unconstrained name**, so `[A, B: Number]` constrains both.
A type parameter's scope begins after its own name, so a later parameter may be
constrained by an earlier one.

**Constraints are their own declaration form** — a compile-time type set,
optionally paired with required methods. Vertex has no interfaces. A constraint
is never a value type and is legal only in a `[...]` position.

```vertex
constraint Ordered {
    ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~string
}

constraint Renderable {
    func render() -> string
}

func min[T: Ordered](a: T, b: T) -> T {
    if a < b { return a }
    return b
}

let m = min[float64](3.14, 2.71)
```

`~T` admits `T` and every type whose underlying type is `T`, so an alias to
`float32` still satisfies `~float32`. A bare `T` admits only `T` exactly. `~`
here is underlying-type, never bitwise-NOT — the two never collide, because a
type-set element is not an expression position.

Multiple elements in a constraint body form an **intersection**: a type argument
must satisfy all of them.

Every instantiation is monomorphized into a separate concrete body, so a
method-constraint call lowers to a direct call on the concrete type. No
interface value, no vtable.

A method may **not** declare a type parameter of its own — everything a method is
generic over comes from its receiver type. The receiver's `[T]` *binds* the
name; it does not introduce a fresh one.

Type arguments may be omitted when every parameter is determined by a value
argument; inference reaches through composite arguments. A parameter appearing
only in the return type cannot be inferred and must be supplied explicitly.
Inference either succeeds or fails — on failure the compiler asks for explicit
arguments rather than guessing.

Predeclared constraints are `any` (every type) and `comparable` (every type
supporting `==`/`!=`).

---

### Arrays, Slices, and Maps

**Fixed arrays** are inline storage — `N × sizeof(T)` bytes, no header, no
pointer. They live wherever the binding lives. The length is part of the type and
must be a compile-time constant.

```vertex
var buf:  [1024]uint8            // zero-filled, no initializer needed
let coords: [3]int32 = [10, 20, 30]

let matrix: [2][2]float32 = [
    [0.0, 1.0],
    [1.0, 0.0],
]
```

**Dynamic arrays** are a three-word `{ptr, len, cap}` header over an implicitly
heap-allocated block — the sole implicit-allocation exception in the language,
justified by the impossibility of fitting a growable buffer into a fixed frame.
The block is still owned and torn down normally through whatever binding holds
the array.

```vertex
var items: []int32 = []

items.push(42)
items.reserve(64)
let last = items.pop()
```

Subscripting with a range produces a **slice view**: two words, `{ptr, len}`,
owning nothing. While a view is live it counts as a shared borrow, so the buffer
may be neither mutated nor transferred:

```vertex
let view = items[1..3]
```

Interior pointers into a `[]T` do not exist, because `push` may reallocate. Only
lifetime-checked views do.

| Form | Storage | Growable |
| --- | --- | --- |
| `var buf: [N]T` | inline | no |
| `let arr = [...]` | inline / rodata | no |
| `var x: []T = []` | heap | yes |

**Maps** use brace literals and require `comparable` keys. Assigning `nil` to a
key erases it — the load-bearing appearance of `nil` outside `typed_ptr`.

```vertex
let scores = {"alice": 42, "bob": 7}
let val    = scores["alice"]

var config: map[string]int32 = {}
config["workers"] = 4
config["workers"] = nil          // erase
```

`nil` is not a general value and has no type of its own.

---

### Tuples

**Parens build, bare commas unbuild.** Parentheses appear when a tuple is
*constructed* — a literal or a type annotation, where they are part of the type's
shape and never optional. The moment a tuple is *pulled apart* — a `let`
destructure or a `return` handing several values back — it is written bare.

```vertex
let pair   = (1, true)
let point  = (x: 10, y: 20)
let single = (42,)          // one-element literal: the trailing comma is required
let plain  = (42)           // just a parenthesized integer

func divmod(a: int32, b: int32) -> (int32, int32) {
    return a / b, a % b
}
let quotient, remainder = divmod(10, 3)

func minMax(values: []int32) -> (min: int32, max: int32) {
    return 0, 100
}
let lo, hi = minMax(values: [3, 1, 4])
```

Positional access uses `.0`, and chains compose: `t.0.0`.

`()` is the unit type: zero bytes, one value, used where a fallible function has
nothing to hand back but an error string. It is not spelled `void`, and there is
no `void` type name — omitting `-> Type` *is* the void form.

Channels can carry tuples:

```vertex
let stream = chan[(int32, bool)](64)
```

---

### Enums

Enums support unit variants, tuple variants, or a mix. `case` appears only inside
`switch`, never in the declaration.

```vertex
enum Direction {
    North,
    South,
    East,
    West,
}

let d: Direction = .South

switch d {
case .North:
case .South:
case .East:
case .West:
}
```

The `.Name` shorthand is legal only where the enum type is inferable from
context: a typed binding, an argument position, a `return`, or a `case` label.

A unit-only enum **is** its discriminant integer — a switch over one with no
`default` must be exhaustive.

```vertex
enum Shape {
    Point,
    Circle(float32),
    Rectangle(float32, float32),
}

let s = Shape.Circle(1.5)

switch s {
case .Point:
case .Circle(r):
case .Rectangle(w, _):
}
```

A payload enum is a tagged union sized to the largest variant plus the tag. A
pattern's payload binding is a **view** into the payload, not a copy; its arity
must match the variant's declared arity exactly, and `_` discards a position
without naming it.

When two positional fields share a type and the distinction matters, carry a
named struct as the payload:

```vertex
struct MousePos {
    x: int32,
    y: int32,
}

enum Event {
    Quit,
    KeyPress(uint8),
    MouseClick(MousePos),
}
```

**Explicit discriminants** require a backing integer type and are legal only on
unit variants. Unassigned variants continue from the previous value:

```vertex
enum Status : int32 {
    Inactive = 0,
    Active,      // 1
    Pending,     // 2
}

let raw = Status.Active as int32    // a tag read, not a conversion
```

Integer-to-enum conversion is an ordinary `switch` with a fallback — there's no
optional type to lean on:

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

---

### Error Handling

There is no optional type, no exception unwinder, and no propagation operator.
Every fallible or possibly-absent value is a plain tuple: the value, and a string
that is empty on success.

```vertex
func parseInt(s: string) -> (int32, string) {
    if s == "" { return 0, "empty string" }
    return 42, ""
}

func findUser(id: int32) -> (User, string) {
    if id < 0 { return User{}, "invalid id" }
    return User(id), ""
}

func removeFile(path: string) -> ((), string) { }
```

Checking the error is the only pattern; the happy path continues directly below:

```vertex
let n, err = parseInt(s: "42")
if err != "" {
    log.printf("failed: %s\n", err)
    return
}
// n is usable past this point
```

The shape repeats at every step and does not get shorter as call depth grows.
Every branch stays visible in the text:

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

A function that may simply find nothing uses the same shape as one that can fail
outright — absence is not a special case. On the error path the value half must
be the type's zero value, never a partially constructed one. The compiler does
not enforce this, matching the model's explicit-over-automatic philosophy.

---

### Raw Pointers and Builtins

`typed_ptr T` is the raw, last-resort pointer: no ownership tracking, no
refcount, no teardown ever emitted. It is the one type in the language where
exclusivity is a convention rather than a proof. Reach for it only when `mut T`
or `[]T` genuinely can't express what's needed.

```vertex
var x: int32 = 42
var p = &x            // address-of  — int32 → typed_ptr int32
let v = &p            // dereference — typed_ptr int32 → int32
&p = 99               // write through p

var q: typed_ptr int32 = nil
if q == nil { }
```

Unary `&` is address-of on an ordinary value and dereference on a `typed_ptr T`.
The direction keys on the operand's *statically written* type, so the meaning of
a source line never flips between instantiations of a generic.

Taking the address *of* a `typed_ptr` binding is therefore unspellable with `&`,
and is the sole purpose of `addr`:

```vertex
let pp = addr(p)      // typed_ptr (typed_ptr int32) — nesting needs parentheses
```

`&` binds **tighter** than member access. `&p.add(1)` parses as `(&p).add(1)`;
write `&(p.add(1))` when a dereferenced read of a computed address is meant. The
parenthesis is the visible mark that an address was computed before it was read.

```vertex
let buf, err = new[uint8](1024)         // zeroed by default
if err != "" { return }
defer delete(buf)

let raw, err2 = new[float32](n, align: 64, zero: false)
```

`new[T]` allocates `count × sizeof(T)` bytes and returns `(typed_ptr T, string)`.
`zero: false` is a *claim* that every byte is written before it's read — nothing
checks the claim. `align` must be a power of two; a violation is an allocation
failure, not a distinct diagnostic. A count whose byte extent overflows `uint64`
is likewise an allocation failure.

| Builtin | Shape |
| --- | --- |
| `sizeof(T)` / `alignof(T)` | compile-time constants |
| `new[T](n, align:, zero:)` | `(typed_ptr T, string)` |
| `delete(p)` | free |
| `resize(p, n)` / `resize(p, n, zero:)` | invalidates `p` on success, leaves it valid on failure |
| `copy(dst, src, n)` | always overlap-safe; there is deliberately no unsafe variant |
| `zero(p, n)` | bulk zero |
| `addr(p)` | address of an addressable `typed_ptr` binding |
| `reinterpret(T, x)` | bit-cast between value types of identical size; never casts pointers |
| `unique(v)` / `shared(v)` | the two heap doors |
| `weak(s)` / `upgrade(w)` | `upgrade` returns `(shared T, string)` |
| `drop(x)` | ends a transferred binding's lifetime without emitting teardown |

`unique(...)` and `shared(...)` *construct* a wrapper, so the copy/transfer rule
does not apply to the operand — it is moved in unconditionally, exactly as into
any constructor.

**Undefined rather than rejected** — this is the whole tradeoff of reaching for
the raw tier: out-of-bounds pointer arithmetic beyond one-past-the-end,
dereferencing an out-of-bounds pointer, ordering pointers into unrelated blocks,
deleting a non-`new` address, reading an unzeroed block before writing it, a bulk
`copy`/`zero` past either extent, and using a `typed_ptr` after a successful
`resize`.

---

### Concurrency

`chan T` is the single currency for moving values between execution contexts.
It's an implicitly heap-resident refcounted handle — copying it bumps the count,
never deep-copies buffered contents.

```vertex
let ch1 = chan[float32]()          // unbuffered rendezvous
let ch2 = chan[int32](64)          // buffered
let ch3: chan float32 = chan[float32]()
```

Construction is ordinary generic instantiation with an optional capacity
argument. Allocation failure panics rather than returning a boundary tuple,
matching native array allocation.

| Method | Waits? | Result |
| --- | --- | --- |
| `.send(v)` | yes | *(void)* |
| `.receive()` | yes | `T` |
| `.trySend(v)` | no | `bool` |
| `.tryReceive()` | no | `(T, string)` |
| `.close()` | no | *(void)* |

**Two spawn prefixes, one handle.** Both `thread` and `async` are call-expression
prefixes — they modify how a call is scheduled, never the callee's signature —
and both evaluate to a receive-only `chan T` carrying the callee's single return
value. Because both terminate in the identical handle, a value produced on an OS
thread and consumed by a reactor task needs no adapter.

```vertex
let worker = thread crunchNumbers(seed: 105)
let result = worker.receive()
```

`thread` runs a call on a real OS thread. The callee is an ordinary function —
nothing about a declaration changes because some call site spawns it.

**`async` is a function marker**, not just a sigil. It declares that the body
contains a real poll point — a place the kernel may answer "not yet" — and sets
`[+Await]` inside:

```vertex
func fetchBody(id: int32) async -> string {
    let conn = await dial(id)
    return await conn.readAll()
}

func main() {
    let body = await fetchBody(1)     // main is [+Await]
    let ch   = async fetchBody(2)     // spawn → chan string
    let other = ch.receive()
}
```

`await` is licensed only under `[+Await]`, which is set by an `async`-marked body
and by `main`. Propagation stops at a function boundary: an anonymous closure
written inside an `async` body may not `await` unless it is itself marked
`async`.

`.receive()` is the one method whose waiting mechanism is not fixed. Called bare
it blocks the calling OS thread; called as `await ch.receive()` it suspends the
task on the reactor. There is no third form and no per-call-site ambiguity — the
mode is fully determined by whether `await` is written. A bare `.receive()`
inside an `async` body would block the underlying thread and starve the reactor,
so it is rejected.

For many values, construct a channel and hand it to the worker as an ordinary
argument:

```vertex
func produce(data: []float32, out: chan float32) {
    for chunk in data {
        out.send(process(chunk))
    }
    out.close()
}

let stream = chan[float32](64)
thread produce(dataset, stream)

while true {
    let chunk, err = stream.tryReceive()
    if err != "" { break }
}
```

**`select`** multiplexes channels. Every case must be a channel receive
operation — not a bare async call, not an arbitrary function, nothing else. To
race a standalone async call, spawn it with the `async` prefix first, which hands
back a `chan T`, and put the receive on that.

```vertex
select {
case a = task1.receive():
    // ...
case b = task2.receive():
    // ...
default:
    // makes the whole statement non-blocking
}
```

`select` introduces **no waiting behavior of its own** — each case waits exactly
the way its receive would wait in that context. A single `select` must therefore
be entirely bare or entirely `await`-prefixed: one mode blocks a thread and the
other suspends a task, and there is no "first ready wins" across two different
wait primitives.

```vertex
select {
case a = await task1.receive():
case b = await task2.receive():
}
```

Function coloring is accepted deliberately. Keeping the state machine explicit is
what lets a custom reactor be written against platform primitives, and what keeps
foreign blocking calls from silently breaking a hidden scheduler. Putting a
blocking call inside an `async` body starves the event loop; the language does
not detect this, and the marker discipline is what makes it visible in the source.

---

### Device Offload

`gpu` and `npu` are device-offload markers built into the core language. They
differ by **programming model, not by vendor** — no vendor name appears anywhere
in Vertex source, and the specific device is selected by the toolchain.

|  | `gpu` | `npu` |
| --- | --- | --- |
| Model | per-thread execution over an index space | whole-array operations over tensors |
| Launch shape | optional `(blocks:, threads:)` | none — shape is carried by the types |
| Body language | unrestricted Vertex | restricted |
| Element access | ordinary subscripting | elementwise operators and namespace calls only |
| Divergent branching | permitted | rejected — the selector must be scalar |

A marker sitting at the vendor level would force source targeting one vendor's
silicon to be written using another vendor's product name, which is why no such
spelling exists.

```vertex
func matmul(a: []float32, b: []float32) gpu -> []float32 {
    // ordinary Vertex
}

let d = gpu(blocks: 16, threads: 256) matmul(x, y)
```

```vertex
func vecAdd(a: tensor[float32, 1024], b: tensor[float32, 1024]) npu
        -> tensor[float32, 1024] {
    return a + b
}

var ha: [1024]float32
var hb: [1024]float32

let sum = npu vecAdd(ha, hb)     // sum: [1024]float32
```

`gpu` and `npu` launches are ordinary **synchronous** calls: host-typed arguments
in, a host-typed result directly out. They do not produce channels. A launch is
legal only when the callee carries the matching marker, and a marked function is
not callable without its prefix — the marker must agree at both ends. A function
carries at most one marker; the markers name mutually exclusive execution
substrates and no combination of two is meaningful.

`LaunchConfig` is legal only on `gpu`; omitting it dispatches with a
compiler-chosen shape.

**Inside an `npu` body:**

- `tensor[...]` types and the `npu.` namespace become grammatical, and are
  grammatical nowhere else.
- Subscripting a tensor is a compile error. Element access is available only
  through elementwise operators and namespace calls.
- Elementwise `+ - * /` and unary `-` require operands sharing element type and
  shape. Comparisons yield a `bool` tensor of the same shape.
- A branch selector must be scalar `bool` or `int32`. Per-element branching is
  expressed with `npu.Select`.
- Loop-carried bindings must keep identical type, shape, and element type across
  iterations. `break` and `continue` are compile errors.
- Plain casts saturate on overflow into the narrow integer types.

Signature-eligible tensor element types are `float32` and `int8`. `bf16`,
`fp8e4m3`, `fp8e5m2`, and `int4` are body-only — legal on a local binding, never
on the function's own signature.

The `npu.` member set is **closed**: not declarable, not shadowable, not
extensible.

| Category | Members |
| --- | --- |
| Math | `Abs Sign Floor Ceil Round Sqrt Rsqrt Exp Expm1 Log Log1p Sin Cos Tan Tanh Sigmoid IsFinite Max Min Mod Pow Atan2` |
| Contraction | `Dot` — accumulates in `float32` regardless of input precision |
| Selection | `Select` |
| Shape | `Reshape Transpose Broadcast Concat Slice Reverse Pad` |
| Reduction | `Sum MaxReduce MinReduce Product` |
| Constants | `Splat Iota` |
| Quantization | `Quantize Dequantize` |

`npu.Quantize[T]` and `npu.Dequantize[T]` take a type argument and a scalar
scale, and are the only members that do.

---

### Abstract Interfaces

Foreign interop is a structural contract: you describe the *call shape* of the
external library using native Vertex types, and the backend emits the correct
calling convention for the target.

A foreign resource is an opaque handle:

```vertex
type SDL_Window = abstract
type NSView     = abstract
```

`abstract` says "structure exists, but Vertex declines to model it" — no
arithmetic, no dereference, no stride. Each alias is a distinct **nominal** type;
two abstract aliases never unify however identically they were minted. An
abstract handle has no `nil` and never participates in a null comparison; its
zeroed representation is legal only as the value half of an error-path tuple.
Copy does not exist for one — it may be accessed or moved.

**Two block forms.** Both require the file to carry a build tag, which picks the
ABI family:

```vertex
package app
build linux

type SDL_Window = abstract

declare module "sdl2" {
    func SDL_CreateWindow(title: string, x: int32, y: int32,
                          w: int32, h: int32, flags: uint32) -> (SDL_Window, string)
    func SDL_DestroyWindow(window: SDL_Window)
    func SDL_GetWindowSize(window: SDL_Window, w: mut int32, h: mut int32)
    func SDL_SetEventFilter(filter: func(int32) -> int32)
}

let window, err = SDL_CreateWindow("game", 0, 0, 800, 600, 2)
```

A declare block is a **linkage boundary, not a namespace** — symbols declared
inside are injected into the file's current package and called unqualified. The
string names the module or framework the linker or bundler resolves; it is never
a path and never contains slashes.

`declare framework` names a platform-bundled, versioned library and is legal only
where the target platform has a first-class notion of one:

```vertex
build darwin

declare framework "AppKit" {
    class NSWindow {
        init func() -> NSWindow
        init func initWithContentRect(contentRect: Rect, styleMask: uint64,
                                      backing: uint64, deferred: bool) -> NSWindow
        func center()
    }
}
```

`init` is a **prefix modifier** on `func`, not a function name. The unnamed form
is what bare `Type(...)` construction resolves to; the named form is what
`Type.someName(...)` resolves to. At most one unnamed initializer per foreign
class, and an initializer must return its enclosing type.

`declare module` is the everything-else bucket — flat C libraries, C++ shared
objects, Windows DLLs, JS modules. It optionally takes a **variant tag**, a
fixed, closed set of strings checked against the file's build tag:

```vertex
declare module ["<variant>"] "engine" { }
```

The bracket exists to *narrow* a default convention, never to introduce a
capability unavailable without it. Omitting it means "use the default for this
build tag". `declare framework` never takes one — bundled message-passing linkage
has exactly one convention by design, and unlike a C++ ABI it does not fork by
compiler, standard library, or flag.

**Exactly what is written is what is linked.** A declare block contains only
declarations corresponding to real entry points: no bodies, no visibility
modifiers, no marker declarations, no remapping clauses, and **no fields** — it
describes call shape only, never foreign-side layout. That is what keeps "which
C++ ABI, exactly?" out of the type system and confined to the linker, where the
variant tag answers it. Ownership keywords are banned too: ownership is a fact
about a wrapper's field, decided in the wrapper.

**The boundary tuple.** Foreign functions do not throw into Vertex. A call that
can fail — a nullable pointer return, a JS call that may throw or yield
`undefined` — is declared as returning the standard error tuple `-> (T, string)`.
A status-plus-out-param shape is `-> (int32, T, string)`. Interop adopts the
native convention, not the reverse.

| Foreign shape | Vertex form |
| --- | --- |
| `const char*` | `string` — marshalled NUL-terminated at the boundary |
| writable scalar out-param | `mut T` — literally the pointer parameter |
| pointer plus length | `[]T` for read, `mut []T` for write |
| pointer held and strided manually | `typed_ptr T` — the last resort |
| property read or foreign static field | ordinary bodyless `func` returning the field's type |

**Ownership lives in the wrapper**, whose `deinit` releases the resource:

```vertex
class Window {
    handle: SDL_Window,
}

func (w: Window) init(title: string) {
    let handle, err = SDL_CreateWindow(title, 0, 0, 800, 600, 2)
    if err != "" {
        // report and bail
    }
    w.handle = handle
}

func (w: Window) deinit() {
    SDL_DestroyWindow(w.handle)
}

func (w: Window) size() -> (int32, int32) {
    var width:  int32 = 0
    var height: int32 = 0
    SDL_GetWindowSize(w.handle, width, height)
    return width, height
}
```

**Callbacks.** A boundary `func(...)` parameter is a bare function pointer: one
word, no environment. Only a non-capturing function converts:

```vertex
func onEvent(code: int32) -> int32 { return 0 }

SDL_SetEventFilter(onEvent)     // legal
```

```vertex
var count: int32 = 0

SDL_SetEventFilter(func(code: int32) -> int32 {
    count += 1     // rejected twice over: assignment to a captured binding,
    return 0       // and a capturing closure cannot cross the boundary
})
```

The rejection is arithmetic rather than stylistic: the closure is two words, the
foreign slot holds one, and nothing on the foreign side will own the environment.

If a foreign signature requires a non-opaque, layout-dependent foreign type to
cross directly, that is out of scope for this layer — wrap it behind an opaque
handle and expose accessors.

---

### Testing

`test` is a function marker, legal only in a file tagged `build test`. It sits
between the parameter list and `->`. A test function takes no parameters and is
auto-discovered by the runner.

```vertex
package arithmetic_test
build test

import "arithmetic"

func test_add()        test -> Expected(int32, "15") { return arithmetic.add(a: 10, b: 5) }
func test_comparison() test -> Expected(bool, "1")   { return 5 > 3 }
func test_no_crash()   test                          { arithmetic.add(a: 0, b: 0) }
```

`Expected(Type, "string")` names the result type and the exact **rendered** value,
compared against the auto-emitted format. Omitting the result type means the test
passes as long as it compiles and runs without crashing.

| Type | Format | `Expected` for value `5` |
| --- | --- | --- |
| signed integers | `%d` | `Expected(int32, "5")` |
| unsigned integers | `%u` | `Expected(uint32, "5")` |
| floats | `%f` | `Expected(float32, "5.000000")` |
| `bool` | `%d` over `1`/`0` | `Expected(bool, "1")` |
| `string` | `%s` | `Expected(string, "hello")` |

A test can also assert that a line **fails to compile**. The two-argument form
additionally requires the diagnostic text to match, which is how the language
pins its own error messages as part of its specification rather than as an
implementation detail:

```vertex
func test_bad_add() test -> Expected(error) {
    return arithmetic.add(a: 10, b: "5")
}

func test_bad_cast() test -> Expected(error, "cannot convert string to int32") {
    let x: int32 = "hello" as int32
}
```

`build test` is the only build tag that changes what is *grammatical* rather than
only what is linkable.

---

## Operator Precedence

Highest binding first. Every level except `..` is left-associative.

| Level | Operators |
| --- | --- |
| 0 | `&` address-of / dereference (prefix — binds tighter than `.`) |
| 1 | `.` member and tuple access, `(...)` call, `[...]` index/slice/instantiate, launch prefixes |
| 2 | `-` `!` `~` (prefix), `await` |
| 3 | `as` |
| 4 | `<<` `>>` |
| 5 | `*` `/` `%` `&` `&*` |
| 6 | `+` `-` `\|` `^` `&+` `&-` |
| 7 | `..` (non-associative) |
| 8 | `==` `!=` `<` `>` `<=` `>=` `===` `!==` |
| 9 | `&&` |
| 10 | `\|\|` |
| — | `=` `+=` `-=` `*=` `/=` `%=` `&=` `\|=` `^=` `<<=` `>>=` — statements, not operators |

**Shifts sit above multiplication.** This is a deliberate departure from C: a
shift is a scaling operation and reads as one, and the C precedence is the single
most common source of parenthesis-omission bugs in bit-manipulation code.

`..` is non-associative — `a..b..c` is a compile error. Every range is half-open;
to cover the full domain of a narrow integer type, iterate a wider one and
convert.

`&&` and `||` short-circuit and their operands must be `bool`.

---

## What Is Absent

Every construction in the language serves one invariant: **every value has a
statically known layout, and every cost is decided at compile time.**

| Absent from every Vertex binary | Replaced by |
| --- | --- |
| Garbage collector | static liveness and scope teardown; refcounts only where `shared` is written |
| Exception unwinder | the `(T, string)` tuple and ordinary control flow |
| Vtables and dynamic dispatch | no inheritance; every call direct; generics monomorphized |
| Drop flags | conditional transfer is a compile error |
| Null-pointer discipline | no general `nil`; absence is an error tuple |
| Runtime type information | every cast resolved statically |
| Hidden allocation | the heap is reachable only through `unique`, `shared`, and the container exception — all spelled in source |

Each row is the same trade in the same direction: a runtime question converted
into a compile-time proof or a visible piece of syntax.

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
      Discover 'test'-marked functions in a `build test` package and run
      them, comparing each against its Expected(...) result.

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
  -v                     report each pipeline stage on stderr

Test flags:
  -dir <path>            directory holding `build test` files (default: .)
  -file <path>           a single test file
  -run <substr>          only run tests whose name contains this substring
  -packages-dir <path>   packages root; overrides $VERTEX_PATH
  -v                      print every test, not just failures
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
| `vertex test` | *console* | discover and run `test`-marked functions |
| `vertex targets` | *console* | list every buildable target, marking the host |

`$VERTEX_PATH` sets the packages root; `-packages-dir` overrides it. When
neither is set, the compiler defaults to `~/.vertex/packages`.

A `-target` must be compatible with a file's `build` tag; a file whose tag does
not match is excluded from the build whole, never partially. An unrecognized
target name is a compile error naming the known set, not a silent fallback.

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

`linux-riscv64` and every powerpc/mips/loongarch/s390x spelling are valid
`.vir` target triples per the language spec but have no `cpu/lower`, object,
or linker implementation in vvm, so they don't appear above. `linux-386` and
`windows-386` are similarly absent: vvm emits x86 ELF object bytes but
registers no ELF linker backend for x86. `browser/wasm`, `browser/js`, and
`android` have no backend at all yet.

---

## Documentation

- [Annex A — Grammar Summary](https://github.com/vertex-language/spec/README.md) — normative grammar, static rules, and the index of rejected forms

---

## License

MIT — see [LICENSE](https://github.com/vertex-language/vertex/LICENSE)