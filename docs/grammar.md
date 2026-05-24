# Vertex Language Grammar

## Specification 1.9

---

## 1. Literals

```swift
// Integers — decimal (default)
42
-1000
1_000_000       // underscore separator — ignored by compiler

// Integers — alternate bases
0b101010        // binary   (base 2)
0o52            // octal    (base 8)
0x2A            // hex      (base 16)

// Hex digits are case-insensitive
0xFF
0xBadFace
0x0123_4567_89ab_cdef   // underscores valid in any base

// Floats — decimal
3.14
1_000.000_1
1.25e2          // = 1.25 × 10²  = 125.0
1.25e-2         // = 1.25 × 10⁻² = 0.0125
1.25E2          // uppercase E — equivalent

// Floats — hex (binary exponent)
0xFp2           // = 15 × 2²  = 60.0
0xFp-2          // = 15 × 2⁻² = 3.75
0xC.3p0         // fractional hex mantissa

// Boolean
true
false

// Nil — absence of a value
nil

// Other literals (unchanged)
"hello"
"A"
```

---

## 2. Variable Declarations

```swift
let x = 10
var y = 20
```

---

## 3. Type Annotations

```swift
let a: int    = 100
let b: int8   = 127
let c: int16  = 32767
let d: int32  = 2147483647
let e: int64  = 9223372036854775807
let f: uint   = 100
let g: uint8  = 255
let h: uint16 = 65535
let i: uint32 = 4294967295
let j: uint64 = 18446744073709551615
let k: float  = 3.14
let l: double = 3.14159265358979
let m: bool   = true
let n: string = "hello"
let o: string = `
multi
line
`
let p: char = "A"
let q: void = ()
```

Multiline strings are delimited by backticks. Content begins after the opening
backtick and ends before the closing backtick. No indentation stripping is
applied.

---

## 4. Numeric Type Conversion

All numeric conversions are explicit. There is no implicit coercion between
numeric types.

```swift
let i: int    = 42
let f: float  = float(i)       // int → float, always safe
let d: double = double(f)      // float → double, always safe
let i2: int   = int(3.99)      // truncates toward zero → 3
let b: int8   = int8(i)        // narrowing — wraps on overflow
```

**Rules:**

* Conversion syntax is `targetType(value)` — no cast keyword.
* No implicit numeric conversion at any point.
* Float-to-integer conversion truncates toward zero.
* Narrowing integer conversions wrap on overflow, identical to the `&+`, `&-`,
  `&*` overflow operators.
* Widening conversions (e.g. `int` → `double`) are always value-preserving.

---

## 5. Arithmetic Operators

```swift
a + b
a - b
a * b
a / b
a % b
-a
```

---

## 6. Compound Assignment

```swift
a += b
a -= b
a *= b
a /= b
a %= b
```

---

## 7. Bitwise Operators

```swift
~a        // NOT
a & b     // AND
a | b     // OR
a ^ b     // XOR
a << b    // left shift
a >> b    // right shift
```

---

## 8. Overflow Operators

```swift
a &+ b    // overflow add
a &- b    // overflow subtract
a &* b    // overflow multiply
```

---

## 9. Comparison Operators

```swift
a == b
a != b
a >  b
a <  b
a >= b
a <= b
```

---

## 10. Logical Operators

```swift
!a
a && b
a || b
```

---

## 11. Range Operators

```swift
0...5     // closed
0..<5     // half-open
```

---

## 12. Ternary Operator

```swift
condition ? a : b
```

---

## 13. Nil-Coalescing

```swift
a ?? b
```

---

## 14. Identity Operators (classes only)

```swift
a === b
a !== b
```

---

## 15. Operator Precedence (high → low)

| Level   | Operators                         |
|---------|-----------------------------------|
| Highest | `<<` `>>`                         |
|         | `*` `/` `%` `&*`                  |
|         | `+` `-` `&+` `&-`                 |
|         | `...` `..<`                       |
|         | `??`                              |
|         | `==` `!=` `<` `>` `<=` `>=`      |
|         | `&&`                              |
|         | `\|\|`                            |
|         | `? :`                             |
| Lowest  | `=` `+=` `-=` `*=` `/=` `%=`     |

---

## 16. If / Else / Else If

```swift
if x > 0 {
    // positive
} else if x < 0 {
    // negative
} else {
    // zero
}
```

---

## 17. Switch

```swift
switch x {
case 0:
    // exactly zero
case 1, 2:
    // one or two
default:
    // anything else
}
```

**String switch:**

```swift
switch s {
case "hello":
    // ...
case "world":
    // ...
default:
    // ...
}
```

**Enum switch:**

```swift
switch direction {
case .north:
    // ...
case .south:
    // ...
case .east:
    // ...
case .west:
    // ...
// no default required — all cases covered
}
```

**Explicit fallthrough:**

```swift
switch x {
case 0:
    // zero
    fallthrough
case 1:
    // zero or one — reached by fallthrough from above
default:
    // other
}
```

**Rules:**

* Cases do not fall through implicitly — each case is independent.
* `fallthrough` transfers control to the next case unconditionally, without
  re-evaluating its condition.
* Multiple values per case are separated by commas.
* `default` is required unless the compiler can statically verify exhaustiveness.
* Switching on an enum with all cases covered is exhaustive — `default` is not
  required.
* An empty case body (with no `fallthrough`) is a compile error.
* `break` may be used inside a case to exit the switch early (§18).
* `switch` may appear anywhere a statement is valid.

---

## 18. Break and Continue

```swift
for i in 0..<10 {
    if i % 2 == 0 { continue }   // skip even numbers
    if i == 7     { break }      // stop at 7
}

var n = 0
while true {
    if n >= 5 { break }
    n += 1
}
```

**Rules:**

* `break` exits the immediately enclosing `for`, `while`, or `switch` statement.
* `continue` skips the remainder of the current loop iteration and begins the
  next.
* `continue` is not valid inside `switch`.
* Neither `break` nor `continue` may appear inside a `defer` block.

---

## 19. Functions

```swift
func add(a: int, b: int) -> int {
    return a + b
}
add(1, 2)
add(a: 1, b: 2)
```

**Mutable parameters (`mut`):**

```swift
func increment(n: mut int) {
    n += 1
}

var count = 0
increment(n: &count)   // count is now 1
```

**Function qualifiers:**

A qualifier sits between the parameter list and the return arrow. All qualifiers
are mutually exclusive except where noted.

```swift
func fetchUser(id: int) async -> User { }
func crunch(data: [float]) thread -> [float] { }
func isolated(data: [float]) process -> [float] { }
func vectorAdd(a: [float], b: [float]) gpu -> [float] { }
```

**Rules:**

* Parameters are immutable by default.
* `mut` marks a parameter as mutable — the function may write back through it.
* The call site must pass a `var` binding prefixed with `&` to a `mut`
  parameter.
* `mut` always appears immediately before the parameter type, after the label.
* `mut` may be applied to any parameter, not only receiver parameters.

---

## 20. While Loop

```swift
var i = 0
while i < 5 {
    i += 1
}
```

---

## 21. For-In Loop

```swift
// Range
for i in 0..<5 {
    // i: int
}

// Closed range
for i in 0...5 {
    // i: int, includes 5
}

// Array
let nums = [1, 2, 3]
for n in nums {
    // n: int
}
```

**Rules:**

* `for i in range` binds the loop variable as the range element type (`int` for
  integer ranges).
* `for item in array` binds each element in order, from index 0 to the last.
* The loop variable is immutable — it may not be assigned inside the body.
* `break` and `continue` are valid inside any `for-in` body (§18).

---

## 22. Arrays (literal + index + mutation)

```swift
// Literal
let nums  = [1, 2, 3]
let flags = [0xFF, 0x00, 0xAB]

// Typed literal
let typed: [int] = [1, 2, 3]

// Trailing comma — valid in multi-line literals
let bytes: [uint8] = [
    0xFF,
    0x00,
    0xAB,
]

// Empty array
var a = [int]()
var b: [int] = []

// Repeating initializer
let zeros  = [int](repeating: 0,  count: 5)
let blanks = [string](repeating: "", count: 3)

// Nested (multidimensional)
let matrix = [[1, 2], [3, 4]]
let grid: [[float]] = [
    [0.0, 1.0],
    [1.0, 0.0],
]

// Read / Write
let first = nums[0]

var items = [10, 20, 30]
items[0] = 99
```

**Rules:**

* Subscript read is valid on both `let` and `var` bindings.
* Subscript write requires the array binding to be `var`.
* Index bounds are not checked at compile time; an out-of-bounds access is a
  runtime error.

---

## 23. Dictionaries (literal + key access + mutation)

```swift
let map = ["a": 1, "b": 2]
let val = map["a"]           // val: int? — nil if key absent

var config = ["debug": 0]
config["debug"] = 1
config["verbose"] = 0
```

**Rules:**

* Key access always returns an optional (`T?`) — the key may not be present.
* Key write requires the dictionary binding to be `var`.
* Assigning `nil` to a key removes it from the dictionary.

---

## 24. Optionals (declare + if-let unwrap only)

```swift
var maybe: int? = nil
maybe = 5
if let val = maybe {
    // val: int — unwrapped, safe to use
}
```

---

## 25. Structs — Stack Allocated, Value Semantics

```swift
struct Point {
    let x: int
    var y: int
}

let p  = Point{x: 3, y: 4}
let p2 = p
let n  = p.x

var q = Point{x: 3, y: 4}
q.y = 10
```

**Multiline form:**

```swift
let p = Point{
    x: 3,
    y: 4,
}
```

**Nested field initialization:**

```swift
struct Line {
    var start: Point
    var end: Point
}

let l = Line{
    start: Point{x: 0, y: 0},
    end:   Point{x: 10, y: 10},
}
```

**Rules:**

* Struct literals use brace syntax: `TypeName{field: value, ...}`.
* All field labels are required — positional initialization is not supported.
* Fields may appear in any order inside the literal.
* Trailing commas are valid in multiline struct literals.
* All fields must be provided — partial initialization is a compile error.
* Struct literals may not appear directly as the condition of `if`, `for`, or
  `switch` statements. Wrap in parentheses to disambiguate:
  `if (Point{x: 1, y: 2} == p) { }`.
* Structs are pure data — no instance methods, no protocols, no inheritance.
* Fields declared `let` cannot be reassigned after the struct is initialized.
* Fields declared `var` can be reassigned only when the enclosing binding is
  `var`.
* A struct may contain fields whose type is another struct.
* Assignment always produces a full copy.
* Structs may not contain class fields that imply ownership semantics.
* Field access via dot notation compiles to a direct byte offset calculation.
* No heap allocation occurs at any point.
* Struct definitions may not appear inside other struct or class definitions.

---

## 26. Associated Functions (index-0 receiver)

A function whose first parameter type is a known struct or class is implicitly
an associated function of that type.

```swift
func describe(p: Point) {
    let n = p.x
}

func reset(p: mut Point) {
    p.x = 0
    p.y = 0
}

p.describe()
p.reset()
```

**Rules:**

* The first parameter at index 0 whose type is a known struct or class is the
  receiver.
* `mut` before the receiver type marks the receiver as mutable. The caller
  must pass a `var` binding prefixed with `&`.
* Associated functions are ordinary functions and follow all rules in §19.
* Any function whose first parameter type is a known struct or class is
  implicitly associated with that type — there is no free-function escape hatch
  at index 0.
* To write a utility function that takes a known type without it acting as the
  receiver, place the known type at index 1 or later.

---

## 27. Enums — Value Semantics, Exhaustive Switch

```swift
enum Direction {
    case north
    case south
    case east
    case west
}

enum Permission {
    case read, write, execute
}
```

**Raw values — int:**

```swift
enum Status: int {
    case inactive = 0
    case active   = 1
    case pending  = 2
}

let s   = Status.active
let raw = Status.active.rawValue    // 1

let fromRaw: Status? = Status(rawValue: 1)
```

**Raw values — string:**

```swift
enum Color: string {
    case red   = "red"
    case green = "green"
    case blue  = "blue"
}

enum Planet: string {
    case mercury   // rawValue = "mercury"
    case venus     // rawValue = "venus"
    case earth     // rawValue = "earth"
}
```

**Rules:**

* Cases are declared with the `case` keyword, one or more per line,
  comma-separated.
* Enum values are accessed via dot notation: `EnumType.caseName`.
* When the type is known from context, the type name may be omitted:
  `.caseName`.
* Raw value types must be `int` or `string`.
* `int` raw values auto-increment from the previous value if omitted; the first
  case defaults to `0` if no value is given.
* `string` raw values default to the case name as a string literal if omitted.
* `.rawValue` accesses the underlying raw value on a raw-value enum.
* `EnumType(rawValue:)` constructs from a raw value and returns `EnumType?`.
* Enums support `==` and `!=`. Raw-value enums also support `<`, `>`, `<=`,
  `>=`.
* A `switch` over an enum with all cases covered is exhaustive — `default` is
  not required.
* Enums are value types — assignment copies.
* Enums may not be nested inside structs or classes.
* Associated values are not supported in 1.9 (deferred).

---

## 28. Classes — Heap Allocated, Programmer Manages Lifetime

```swift
class Animal {
    var name: string
}

func init(a: mut Animal, name: string) {
    a.name = name
}

func deinit(a: mut Animal) {
    // runs before memory is freed
}

let a = Animal(name: "Rex")
a.delete()
```

**Rules:**

* Classes are heap allocated — the runtime cost is exactly what the programmer
  pays.
* Assignment passes a reference — two variables may point to the same object.
* Identity operators `===` and `!==` compare references, not values.
* Inheritance is not supported — classes are standalone types.
* A class may contain fields whose type is a struct.
* `init` is a reserved associated function name called automatically after
  allocation. Its first parameter must be the instance, typed `mut`.
* `deinit` is a reserved associated function name. It runs automatically when
  `.delete()` is called. Its first parameter must be the instance, typed `mut`.
* Neither `init` nor `deinit` may be called directly.
* If no `func init` is declared, the compiler provides a default memberwise
  initializer.
* The programmer is responsible for calling `.delete()` on class instances.
* Failing to call `.delete()` on a class instance is a memory leak.
* Class definitions may not appear inside other class or struct definitions.

---

## 28.1 Reference Counting — `.new()`

```swift
let a = Animal(name: "Rex").new()
let b = a                           // count = 2
// b scope ends — count = 1
// a scope ends — count = 0, deinit called, freed
```

**Weak references:**

```swift
let a = Animal(name: "Rex").new()
weak let b = a                   // b: Animal? — non-owning, count stays 1

if let animal = b {
    // safe — animal is Animal within this scope
}
```

**Rules:**

* `.new()` is postfix on any class instantiation expression.
* `.new()` and `.delete()` are mutually exclusive.
* `weak let` declares a non-owning reference. It does not increment the count.
* `weak let` produces a value of type `T?`. Use `if let` to safely unwrap
  before use.
* After the owning references reach zero, all `weak` references become `nil`.
* `weak` is only valid on ref counted instances (`.new()`).

---

## 29. Defer

```swift
let a = Animal(name: "Rex")
defer a.delete()
```

**Anonymous function form:**

```swift
defer func() { cleanup(a) }()
```

**Multiple defers (LIFO):**

```swift
defer a.delete()           // runs second
defer b.delete()           // runs first
```

**Rules:**

* `defer` takes a direct function call — no surrounding braces.
* For multi-statement cleanup, use `defer func() { ... }()` — the trailing `()`
  invokes the anonymous function immediately, deferring its execution to scope
  exit.
* `defer` executes when the immediately enclosing scope exits.
* Multiple `defer` statements in the same scope run in reverse declaration order
  (LIFO).
* `defer` may appear anywhere in a function body — not only at the top.
* The deferred call may not contain `return`, `break`, or `continue`.
* `defer` is not valid at the top level — only inside function bodies.

---

## 30. Generics (unconstrained)

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

## 31. Import Declarations

```swift
import "github.com/something"

import (
    "github.com/something"
    "github.com/something/else"
)
```

**Rules:**

* Import paths are double-quoted string literals.
* The grouped form parenthesizes one or more newline-separated paths — no
  commas.
* Imports must appear at the top of a file, after any `package` and `build`
  declarations.

---

## 32. First-Class Function Types

A function type describes the signature of a callable value — its parameter
types and return type.

```swift
// variable holding a function
let double:    func(int) -> int
let predicate: func(int) -> bool
let transform: func(string, int) -> string

// void return — arrow omitted
let onFire: func(int)

// function type as a parameter
func apply(values: [int], f: func(int) -> int) -> [int] { }

// function type as a return type
func makeAdder(n: int) -> func(int) -> int { }

// mutable parameter in a function type
func run(n: mut int, f: func(mut int)) { }

// calling a function value — standard call syntax
let result = double(21)
```

**Rules:**

* Function type syntax is `func(ParamTypes) -> ReturnType`.
* When the return type is `void`, the arrow and return type are omitted:
  `func(int)`.
* `mut` before a type in a function type signature indicates the caller must
  pass a `var` binding prefixed with `&` to that position — identical to `mut`
  in named function declarations (§19).
* Function types may appear anywhere a type is valid: `let`, `var`, parameter
  annotations, and return type annotations.
* A function value is called with standard call syntax: `f(42)`.
* Function types are value types — assignment copies the callable reference.
* Parameter names are not part of the type — only the types and their `mut`
  modifiers are significant.

---

## 33. Anonymous Functions

An anonymous function is a function literal without a name. It follows the
same syntax as a named function declaration, minus the name. It may capture
values from the enclosing scope.

```swift
// stored in a variable
let double = func(n: int) -> int { return n * 2 }

// void return — arrow omitted
let log = func(n: int) { print(n) }

// passed inline — higher-order function pattern
let doubled = process(nums, func(n: int) -> int {
    return n * 2
})

// passed inline — callback registration pattern
emitter.on(func(n: int) -> int {
    return n * 2
})
```

**Capture — value semantics:**

```swift
let factor = 3
let multiply = func(n: int) -> int {
    return n * factor    // factor captured by value at creation
}

var count = 0
let increment = func() {
    count += 1           // compile error — captured copy, not the original
}
```

**Explicit writeback via mut parameter:**

```swift
func run(n: mut int, f: func(mut int)) {
    f(&n)
}

var total = 0
run(n: &total, f: func(n: mut int) {
    n += 10              // writes back through n — total is now 10
})
```

**Rules:**

* Anonymous function syntax is `func(params) -> ReturnType { body }` —
  identical to a named function declaration minus the name.
* Anonymous functions capture variables from the enclosing scope by value
  at the point of creation.
* Captured values are copied — mutations inside the anonymous function do not
  affect the original binding.
* To write back through a variable, pass it explicitly as a `mut` parameter
  (§19) — capture alone cannot produce writeback.
* `mut` parameters inside an anonymous function follow the same rules as in
  named functions (§19) — the call site passes a `var` binding prefixed
  with `&`.
* `return` inside an anonymous function returns from the anonymous function,
  not the enclosing function.
* Anonymous functions may not refer to themselves by name — they are not
  recursive. Recursion requires a named function.
* Anonymous functions are valid anywhere an expression is valid.
* The inferred type of an anonymous function is `func(ParamTypes) -> ReturnType`
  (§32).

---

## 35. Async / Await

```swift
func fetchUser(id: int) async -> User {
    // body
}

let user = fetchUser(id: 1).await()
let result = fetchUser(id: 1).await().name
```

**Rules:**

* `async` is written between the parameter list and the return arrow.
* A function with no return value may omit the return arrow: `func f(...) async { }`.
* `.await()` may only appear inside an `async` function body.
* Multiple `.await()` calls may appear in the same function body.

---

## 36. Tuples

```swift
let pair  = (1, true)
let point = (x: 10, y: 20)
let nothing: () = ()
```

**Destructuring:**

```swift
let (a, b) = pair
let (x, y): (int, int) = (14, 17)
```

**Function return:**

```swift
func minMax(values: [int]) -> (min: int, max: int) {
    return (0, 100)
}

let (lo, hi) = minMax(values: [3, 1, 4])
```

**Rules:**

* `()` is the empty tuple and is an alias for `void`.
* A single-element parenthesised expression like `(x)` has the type of `x`,
  not a tuple.
* Element labels are optional. Unlabelled elements are only accessible via
  destructuring.
* Two tuple types are identical if they share the same element types and labels
  in order.
* `==`, `!=`, `<`, `>`, `<=`, `>=` work on tuples whose elements are all
  comparable, up to 6 elements. Labels are ignored during comparison.
* Tuples are value types — assignment copies all elements.

---

## 37. Error Handling

### 37.1 Optionals — absence without context

```swift
func findUser(id: int) -> User? {
    if id < 0 { return nil }
    return User(id)
}

if let user = findUser(id: 1) { }
let name = findUser(id: -1) ?? defaultUser
```

### 37.2 Tuples — multiple returns, caller decides

```swift
func divide(a: int, b: int) -> (int, string?) {
    if b == 0 { return (0, "division by zero") }
    return (a / b, nil)
}

let (result, err) = divide(a: 10, b: 0)
if err != nil { }
```

### 37.3 Result — explicit Ok/Err

```swift
func parseInt(s: string) -> Result(int, string) {
    if s == "" { return Result(Err, "empty string") }
    return Result(Ok, 42)
}
```

**Consuming with `if let`:**

```swift
if let value = divide(a: 10, b: 2) {
    // value: int
}
```

**Consuming with `switch`:**

```swift
switch divide(a: 10, b: 0) {
case Ok(let value):
    // value: int
case Err(let err):
    // err: string
}
```

**Propagating with `.try()`:**

```swift
func process(s: string) -> Result(int, string) {
    let n = parseInt(s).try()
    let d = divide(a: n, b: 2).try()
    return Result(Ok, d)
}
```

**Rules:**

* `Result(Ok, value)` and `Result(Err, error)` are the only valid construction
  forms.
* `if let` on a `Result` binds the `Ok` value only — use `switch` to inspect
  `Err`.
* `.try()` may only appear inside a function whose return type is
  `Result(T, E)`.

### 37.4 Choosing the Right Primitive

| Situation                               | Use           |
|-----------------------------------------|---------------|
| Value may simply not exist              | `T?`          |
| Caller needs value and error together   | `(T, E?)` tuple |
| Caller must handle Ok or Err explicitly | `Result(T, E)` |

---

## 38. GPU Kernels

```swift
func vectorAdd(a: [float], b: [float]) gpu -> [float] {
    // normal logic in here
    return result
}

let result = vectorAdd(a: x, b: y).dispatch()
let result = vectorAdd(a: x, b: y).dispatch(gpu: 0, mem: 256)
```

**Rules:**

* `gpu` is written between the parameter list and the return arrow.
* `.dispatch()` returns the kernel result directly.
* `.dispatch(gpu: n, mem: n)` selects a specific device and memory allocation.

---

## 39. Threads

```swift
func crunchData(data: [float]) thread -> [float] {
    // runs in a thread, shared memory
}

let result = crunchData(data: x).spawn()
let result = crunchData(data: x).spawn(threads: 4)
```

**Rules:**

* `thread` is written between the parameter list and the return arrow.
* Threads share memory with the caller — passing data is zero-copy.
* `.spawn(threads: n)` spawns exactly n threads.

---

## 40. Processes

```swift
func isolatedWork(data: [float]) process -> [float] {
    // runs in a separate process, full memory isolation
}

let result = isolatedWork(data: x).fork()
let result = isolatedWork(data: x).fork(processes: 4)
```

**Rules:**

* `process` is written between the parameter list and the return arrow.
* Processes have fully isolated memory — data is copied across the boundary.
* `.fork(processes: n)` forks exactly n processes.

---

## 41. Channels

```swift
let ch: channel float = Channel()
let ch: channel float = Channel(size: 8)

ch.send(42.0)
let val = ch.receive()
ch.close()
```

**Runtime transport by context:**

| context   | transport                        |
|-----------|----------------------------------|
| `async`   | shared memory, non-blocking      |
| `thread`  | shared memory, lightweight       |
| `process` | ring buffer, high-speed IPC      |

**Rules:**

* `channel T` is the type — `Channel()` is the constructor.
* Unbuffered channels block the sender until the receiver is ready.
* Buffered channels block the sender only when the buffer is full.
* `.send()`, `.receive()`, and `.close()` are the only valid operations.
* Calling `.send()` on a closed channel is a runtime error.
* `.receive()` on a closed, empty channel is a runtime error.

---

## 42. Postfix Execution Model — Summary

| keyword   | postfix                                     | meaning                     |
|-----------|---------------------------------------------|-----------------------------|
| `async`   | `.await()`                                  | suspend until done          |
| `thread`  | `.spawn()` / `.spawn(threads: n)`           | shared memory parallel      |
| `process` | `.fork()` / `.fork(processes: n)`           | isolated memory parallel    |
| `gpu`     | `.dispatch()` / `.dispatch(gpu: n, mem: n)` | gpu device execution        |
| —         | `.new()`                                    | opt class into ref counting |
| —         | `.delete()`                                 | manually free class instance|
| —         | `.try()`                                    | propagate Result error      |

---

## Explicitly Out of Scope in 1.9

| Feature                                      | Status   |
|----------------------------------------------|----------|
| Inheritance                                  | Removed  |
| String interpolation `\()`                   | Removed  |
| `_` parameter labels                         | Removed  |
| `self` keyword                               | Removed  |
| `static` keyword                             | Removed  |
| Methods inside structs or classes            | Removed  |
| `mutating` keyword                           | Removed  |
| Protocols                                    | Removed  |
| Extensions                                   | Removed  |
| `try` / `throws` / `do-catch`               | Removed  |
| Nested structs or classes                    | Removed  |
| Generic constraints (`where T:`)             | Deferred |
| Closures                                     | §33      |
| Enums with associated values                 | Deferred |
| Access control                               | Deferred |
| Pattern matching beyond `if let` and `switch`| Deferred |
| Custom operators                             | Deferred |
| `async let` / `TaskGroup` concurrency        | Deferred |
| `actor` keyword                              | Deferred |
| Conditional build expressions                | Deferred |
| Import aliasing                              | Deferred |
| Mixed tuple element labels                   | Deferred |
| Tuple `for`-loop destructuring               | Deferred |
| `inout` tuple parameters                     | Deferred |
| Tuple splat into function arguments          | Deferred |
| `Err` value binding in `if let` else branch  | Deferred |
| `select` over multiple channels              | Deferred |
| `weak` in manual class instances             | Deferred |
| GPU grid/block control                       | Deferred |
| Labeled `break` / `continue`                 | Deferred |