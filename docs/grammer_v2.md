# Vertex Language Grammar Specification
# vertex_grammar.md
# v0.1 — Draft

Vertex is a statically typed, compiled language that targets WebAssembly 2.0
via wasm-compiler. It compiles to native machine code with no VM or runtime.
Types map directly to WASM's validated type system, with friendly surface
aliases layered on top.

---

## 1. Lexical Elements

### 1.1 Identifiers
    identifier  = letter { letter | digit | "_" }
    letter      = "A"..."Z" | "a"..."z" | "_"
    digit       = "0"..."9"

### 1.2 Keywords
    async       await       break       continue
    else        enum        export      false
    for         func        gpu         if
    import      interface   let         match
    return      true        as          in

### 1.3 Operators & Punctuation
    +   -   *   /   %           arithmetic
    ==  !=  <   >   <=  >=      comparison
    &&  ||  !                   logical
    =                           assignment
    :                           type annotation (let bindings and params only)
    .                           member access
    ,                           separator
    ;                           statement terminator (optional; newline implies one)
    { }   [ ]   ( )             delimiters
    =>                          match arm
    ->                          (reserved for future closure return types)
    `                           template literal delimiter
    ..                          range (used in for..in)

### 1.4 Comments
    // single line comment
    /* multi
       line comment */

---

## 2. Type System

Vertex types map directly to WebAssembly value types.
Friendly aliases are first-class — the compiler resolves them to their
underlying wasm.ValType before encoding.

### 2.1 Primitive Types

| Vertex type | WASM ValType | Notes                               |
|-------------|--------------|-------------------------------------|
| i32         | wasm.I32     | 32-bit integer                      |
| i64         | wasm.I64     | 64-bit integer                      |
| f32         | wasm.F32     | 32-bit float                        |
| f64         | wasm.F64     | 64-bit float                        |
| bool        | wasm.I32     | alias; 0 = false, 1 = true          |
| byte        | wasm.I32     | alias; unsigned 8-bit semantics     |
| int         | wasm.I64     | alias; platform-width signed int    |
| float       | wasm.F64     | alias; double-precision float       |
| string      | wasm.I32     | alias; pointer into linear memory   |
|             |              | (ptr: i32, len: i32) ABI pair       |
| void        | (none)       | empty result type; no stack value   |

### 2.2 Composite Types

    // Slice — a contiguous region of linear memory
    // Represented as (ptr: i32, len: i32) at the ABI level
    []T

    // Examples
    []f32       // slice of f32
    []byte      // byte buffer
    []string    // slice of strings

### 2.3 Interface Types (Structs)

    // Interfaces define named field layouts.
    // They are lowered to a linear memory layout at compile time.
    // Field access compiles to a calculated offset load/store.
    interface Point {
        x: f32
        y: f32
    }

    // Nested interfaces are allowed
    interface Rect {
        origin: Point
        size:   Point
    }

### 2.4 Enum Types

    // Enums are sum types. Each variant may carry a typed payload.
    // Lowered to a tagged union in linear memory: tag (i32) + payload.
    enum Direction {
        North
        South
        East
        West
    }

    enum Shape {
        Circle(f32)         // radius
        Rect(f32, f32)      // width, height
        Point               // no payload
    }

### 2.5 Generic Built-in Types

The following generic types are compiler built-ins for v0.1.
Full user-defined generics land in v0.3.

    Option<T>               // Some(T) | None
    Result<T, E>            // Ok(T)   | Err(E)
    Promise<T>              // async coroutine handle; resolves to T

### 2.6 Type Annotations

    // Colon is used for type annotations on let bindings and parameters only.
    // Return types follow the closing paren directly — no colon.

    let x: f32 = 1.0
    let y = 2.0         // inferred as f64 (default float literal type)
    let n = 42          // inferred as i64 (default integer literal type)

---

## 3. Imports

Vertex supports two import forms.

### 3.1 WASM Module Import

Imports a pre-compiled `.wasm` binary, making its exported functions
available under an alias. At compile time the compiler resolves this to
a WASM import entry:

    (import "<module>" "<name>" (func ...))

Syntax:
    import "<path>.wasm" as <alias>

Examples:
    import "math.wasm"  as mathLib      // C++ compiled via vcx
    import "simd.wasm"  as simd

Usage:
    mathLib.sqrt(x)     // compiles to a direct WASM function call
    simd.dot(a, b)

### 3.2 Standard Library Import

    import "std/<package>"

Examples:
    import "std/fmt"
    import "std/mem"
    import "std/os"
    import "std/http"
    import "std/runtime"

Standard library packages are resolved at link time to WAPI syscall stubs
or inline compiler intrinsics. No runtime linked in.

### 3.3 Import Block (multiple)

    import (
        "math.wasm"  as mathLib
        "std/fmt"
        "std/mem"
        "std/http"
    )

---

## 4. Declarations

### 4.1 Interface Declaration

    InterfaceDecl = "interface" identifier "{" { FieldDecl } "}"
    FieldDecl     = identifier ":" Type newline

    interface Point {
        x: f32
        y: f32
    }

Interfaces have no methods. They are pure data layout descriptors.
Methods are free functions that accept an interface type as a parameter.

### 4.2 Enum Declaration

    EnumDecl    = "enum" identifier "{" { EnumVariant } "}"
    EnumVariant = identifier [ "(" TypeList ")" ] newline
    TypeList    = Type { "," Type }

    // No payload — variants are i32 constants
    enum Direction {
        North
        South
        East
        West
    }

    // With payload — tagged union in linear memory
    enum Shape {
        Circle(f32)
        Rect(f32, f32)
        Point
    }

### 4.3 Function Declaration

    FuncDecl     = { FuncModifier } "func" identifier "(" ParamList ")" [ ReturnType ] "{" Body "}"
    FuncModifier = "export" | "async" | "gpu"
    ParamList    = [ Param { "," Param } ]
    Param        = identifier ":" Type
    ReturnType   = Type                         // no colon — type follows ) directly

Modifiers may be combined where valid:

    export async func fetchAndProcess(url: string) f32 { ... }

Invalid combinations are a compile error:

    gpu async func ...      // ERROR: gpu functions cannot be async
    gpu export func ...     // ERROR: gpu functions cannot be exported

Examples:

    // unexported, void return
    func greet(name: string) {
        fmt.println(`Hello ${name}`)
    }

    // exported, returns f32
    export func calculateDistance(a: Point, b: Point) f32 {
        let dx = a.x - b.x
        let dy = a.y - b.y
        return mathLib.sqrt((dx * dx) + (dy * dy))
    }

    // async, returns Promise<string> implicitly
    async func fetchData(url: string) string {
        let resp = await http.get(url)
        return resp.body
    }

    // exported async
    export async func getDistance(url: string) f32 {
        let raw   = await fetchData(url)
        let point = parsePoint(raw)
        return calculateDistance(origin, point)
    }

    // gpu kernel — no return type; operates on slices in linear memory
    gpu func normalize(x: []f32, y: []f32, out: []f32) {
        let i   = thread.idx
        let len = mathLib.sqrt((x[i] * x[i]) + (y[i] * y[i]))
        out[i]  = x[i] / len
    }

    // entry point
    func main() {
        ...
    }

    // async entry point
    async func main() {
        let dist = await getDistance("https://api.example.com/point")
        fmt.println(`Distance: ${dist}`)
    }

#### 4.3.1 export modifier
Marks the function as a WASM export entry. Adds the function to the
ExportSection of the compiled module. Callable from the host or from
other Vertex modules.

#### 4.3.2 async modifier
Marks the function as a stackless coroutine. The declared return type `T`
is implicitly wrapped as `Promise<T>`. The function body may contain
`await` expressions. See section 10 for full async/await lowering.

#### 4.3.3 gpu modifier
Designates a GPU kernel. The compiler translates the function body to
PTX (Parallel Thread Execution) IR instead of WASM opcodes and schedules
it for execution on the GPU. Parameters must be slices of numeric
primitives. GPU functions cannot call non-gpu functions and cannot be
exported. Dispatch is via a runtime intrinsic (v0.2).

---

## 5. Statements

    Statement =
        | LetStmt
        | AssignStmt
        | ReturnStmt
        | IfStmt
        | ForStmt
        | MatchStmt
        | ExprStmt
        | Block

### 5.1 Let Statement (variable binding)

    LetStmt = "let" identifier [ ":" Type ] "=" Expr

    let x: f32     = 10.0
    let name       = "vertex"           // inferred: string
    let p: Point   = { x: 1.0, y: 2.0 }
    let s: Shape   = Circle(5.0)
    let d          = Direction.North

Let bindings are immutable by default.
`let mut` is reserved for v0.2:

    let mut counter: i32 = 0    // reserved — not in v0.1

### 5.2 Assignment

    AssignStmt = Expr "=" Expr

    p.x    = 5.0
    buf[i] = 0

### 5.3 Return Statement

    ReturnStmt = "return" [ Expr ]

    return 42
    return mathLib.sqrt(x)
    return Some(value)
    return                      // void return

### 5.4 If Statement

    IfStmt = "if" Expr Block [ "else" ( IfStmt | Block ) ]

    if dist < 1.0 {
        fmt.println(`close`)
    } else if dist < 10.0 {
        fmt.println(`near`)
    } else {
        fmt.println(`far`)
    }

### 5.5 For Statement

    ForStmt =
        | "for" identifier "in" Expr Block      // range loop
        | "for" Expr Block                       // while-style loop

    for i in 0..len(xs) {
        fmt.println(xs[i])
    }

    for running {
        tick()
    }

### 5.6 Match Statement / Expression

Match can be used as a statement or as an expression that produces a value.
Every variant of the matched enum must be handled — exhaustiveness is
enforced at compile time. Missing arms are a hard error.

    MatchExpr  = "match" Expr "{" { MatchArm } [ DefaultArm ] "}"
    MatchArm   = VariantPattern "=>" ( Expr | Block )
    DefaultArm = "_" "=>" ( Expr | Block )

    VariantPattern =
        | identifier                        // no-payload variant
        | identifier "(" BindingList ")"    // payload binding

    BindingList = identifier { "," identifier }

    // As a statement
    match direction {
        North => fmt.println(`going north`)
        South => fmt.println(`going south`)
        East  => fmt.println(`going east`)
        West  => fmt.println(`going west`)
    }

    // As an expression — all arms must return the same type
    let area = match shape {
        Circle(r)    => mathLib.sqrt(r * r) * 3.14159
        Rect(w, h)   => w * h
        Point        => 0.0
    }

    // With block arms
    let msg = match result {
        Ok(body)  => body
        Err(e)    => {
            fmt.println(`Error: ${e}`)
            ""
        }
    }

    // Default arm (use sparingly — prefer exhaustive matches)
    match code {
        200 => fmt.println(`ok`)
        404 => fmt.println(`not found`)
        _   => fmt.println(`other`)
    }

---

## 6. Expressions

    Expr =
        | Literal
        | identifier
        | Expr "." identifier                   // member access
        | Expr "[" Expr "]"                     // index
        | Expr "(" [ ArgList ] ")"              // call
        | "await" Expr                          // await
        | Expr BinaryOp Expr                    // binary
        | UnaryOp Expr                          // unary
        | "(" Expr ")"                          // grouped
        | StructLiteral
        | EnumVariantExpr
        | TemplateLiteral
        | MatchExpr

    BinaryOp = "+" | "-" | "*" | "/" | "%" | "==" | "!=" | "<" | ">"
             | "<=" | ">=" | "&&" | "||"
    UnaryOp  = "-" | "!"

### 6.1 Struct Literal

    StructLiteral = "{" FieldInit { "," FieldInit } [ "," ] "}"
    FieldInit     = identifier ":" Expr

    let p: Point = { x: 10.0, y: 20.0 }

### 6.2 Enum Variant Expression

    EnumVariantExpr =
        | identifier "." identifier                         // no payload
        | identifier "." identifier "(" ArgList ")"        // with payload
        | identifier "(" ArgList ")"                       // shorthand (in-scope enum)

    let d = Direction.North
    let s = Shape.Circle(5.0)
    let r = Shape.Rect(3.0, 4.0)

    // shorthand — valid when enum type is unambiguous from context
    return Some(value)
    return None
    return Ok(result)
    return Err("something went wrong")

### 6.3 Await Expression

    AwaitExpr = "await" Expr

`await` suspends the current async coroutine until the awaited `Promise<T>`
resolves. May only appear inside an `async func`. Using `await` outside an
async function is a hard compile error:

    // ERROR: await used outside async func
    func bad() {
        let x = await fetchData(url)
    }

    // OK
    async func good() string {
        let x = await fetchData(url)
        return x
    }

### 6.4 Template Literal

    TemplateLiteral = "`" { RawChar | "${" Expr "}" } "`"

    `The distance is ${dist}`
    `Hello ${user.name}, you have ${count} messages`
    `Status: ${match code { 200 => "ok"  _ => "error" }}`

Template literals compile to a sequence of string concatenation operations
over linear memory. Interpolated expressions are coerced to string via a
compiler-generated format intrinsic.

### 6.5 Member Access

    a.x             // field access — compiles to load at computed byte offset
    mathLib.sqrt    // module-qualified function reference
    fmt.println     // stdlib function reference
    Direction.North // enum variant reference

### 6.6 Slice Index

    buf[i]          // compiles to load with bounds check (debug builds)
                    // bounds check elided in release builds

---

## 7. Literals

| Literal           | Example               | Inferred type |
|-------------------|-----------------------|---------------|
| Integer           | 42                    | i64           |
| Float             | 3.14                  | f64           |
| Float (explicit)  | 1.0f                  | f32           |
| Boolean           | true / false          | bool (i32)    |
| String            | "hello"               | string        |
| Template string   | \`hi ${name}\`        | string        |
| Struct            | { x: 1.0, y: 2.0 }    | inferred      |
| Enum variant      | Direction.North       | Direction     |
| Enum with payload | Shape.Circle(5.0)     | Shape         |

---

## 8. Scoping Rules

- File scope: imports, interface/enum declarations, top-level functions
- Function scope: parameters, let bindings
- Block scope: let bindings inside `if`, `for`, `match` arms, nested blocks
- Match arm bindings (e.g. `Circle(r)`) are scoped to their arm block only
- No closures in v0.1 — functions cannot capture outer variables
- Interfaces are nominal — two interfaces with identical fields are distinct types
- Enums are nominal — variant names are qualified by their enum type unless
  the type is unambiguous from context (e.g. function return type)

---

## 9. WASM Lowering Summary

| Vertex construct             | WASM encoding                                        |
|------------------------------|------------------------------------------------------|
| Primitive local              | local.get / local.set                               |
| Arithmetic expr              | i32.add, f32.mul, etc.                              |
| Interface field access       | i32.load / f32.load at computed byte offset         |
| Slice []T                    | (ptr: i32, len: i32) pair in locals                 |
| Slice index                  | i32.add(ptr, i32.mul(idx, sizeof(T))) + load        |
| Function call                | call $funcIdx                                       |
| Imported WASM call           | call $importIdx                                     |
| export func                  | ExportSection entry + function body                 |
| string literal               | data segment + (ptr: i32, len: i32)                 |
| Template literal             | data segments + memory.copy + format stubs          |
| bool                         | i32; branch on i32.eqz                              |
| Simple enum variant          | i32 constant                                        |
| Payload variant alloc        | linear memory alloc: i32.store(tag) + payload store |
| match on enum                | i32.load(tag ptr) + br_table over variant indices   |
| match arm binding            | payload load at known offset into local             |
| Option<T>.None               | tag = 1, no payload load                            |
| Option<T>.Some(v)            | tag = 0, payload load → local                       |
| Exhaustiveness check         | compile-time — missing arms are a hard error        |
| async func f()               | state machine + frame alloc in linear memory        |
| await expr                   | spill locals → frame, return state index            |
| Promise<T> resolve           | write result → frame, set state=1, call resume      |
| async func main()            | std/runtime event loop drives top-level promise     |
| gpu func                     | PTX codegen path (not WASM opcodes)                 |

---

## 10. Async / Await

Vertex's async model compiles to a stackless coroutine represented as a
WASM state machine function with its frame stored in linear memory.
No threads. No OS scheduler. No runtime event loop baked in —
the host or `std/runtime` drives resumption.

### 10.1 Async Functions

An `async func` with declared return type `T` implicitly returns `Promise<T>`.
The actual value `T` is delivered when the coroutine completes.

    async func fetchData(url: string) string {
        let resp = await http.get(url)
        return resp.body
    }

At the WASM level `async func` compiles to a state machine:
- Each `await` point becomes a numbered state
- The coroutine's locals are spilled to a heap-allocated frame in linear memory
- Resumption is a `call` into the function with the frame pointer + state index

### 10.2 Await Expression

`await` may only appear inside an `async func`.
It suspends the current coroutine and yields control until the awaited
`Promise<T>` resolves. The expression evaluates to the resolved value of
type `T`.

    async func getDistance(url: string) f32 {
        let raw   = await fetchData(url)
        let point = parsePoint(raw)
        return calculateDistance(origin, point)
    }

### 10.3 Promise<T>

`Promise<T>` is a compiler-managed struct in linear memory:

    // internal layout (not user-facing)
    interface _Promise {
        state:   i32    // 0=pending, 1=resolved, 2=rejected
        result:  i32    // ptr to resolved value in linear memory
        frame:   i32    // ptr to coroutine stack frame
        resume:  i32    // function table index for resumption
    }

### 10.4 Async Entry Point

`main` may be declared async. `std/runtime` drives the top-level event loop:

    async func main() {
        let dist = await getDistance("https://api.example.com/point")
        fmt.println(`Distance: ${dist}`)
    }

### 10.5 Error Handling with Async

Async functions compose naturally with `Result<T, E>`:

    async func load(url: string) Result<string, string> {
        let resp = await http.get(url)
        if resp.status != 200 {
            return Err(`bad status: ${resp.status}`)
        }
        return Ok(resp.body)
    }

    async func main() {
        let data = match await load("https://api.example.com/data") {
            Ok(body)  => body
            Err(msg)  => {
                fmt.println(`Error: ${msg}`)
                ""
            }
        }
        fmt.println(data)
    }

---

## 11. Enums

Vertex enums are sum types — each variant can carry typed payload.
They compile to a tagged union in linear memory: a tag `i32` followed
by the payload of the largest variant (padded to alignment).

### 11.1 Simple Enum (no payload)

Variants with no payload compile to manifest i32 constants:

    enum Direction {
        North       // → i32 const 0
        South       // → i32 const 1
        East        // → i32 const 2
        West        // → i32 const 3
    }

### 11.2 Enum with Payload (tagged union)

    enum Shape {
        Circle(f32)         // tag=0, payload: f32
        Rect(f32, f32)      // tag=1, payload: (f32, f32)
        Point               // tag=2, no payload
    }

In linear memory:

    [ tag: i32 ][ payload: <largest variant size, padded to alignment> ]

### 11.3 Match on Enum

Enums are consumed via `match`. Exhaustiveness is enforced at compile time.
The compiler emits a `br_table` over the tag i32.

    let area = match shape {
        Circle(r)    => mathLib.sqrt(r * r) * 3.14159
        Rect(w, h)   => w * h
        Point        => 0.0
    }

    func describe(d: Direction) string {
        return match d {
            North => "north"
            South => "south"
            East  => "east"
            West  => "west"
        }
    }

### 11.4 Option<T> — built-in enum

    enum Option<T> {
        Some(T)
        None
    }

    func safeDivide(a: f32, b: f32) Option<f32> {
        if b == 0.0 {
            return None
        }
        return Some(a / b)
    }

    let result = safeDivide(10.0, 2.0)
    let value  = match result {
        Some(v) => v
        None    => 0.0
    }

### 11.5 Result<T, E> — built-in enum

    enum Result<T, E> {
        Ok(T)
        Err(E)
    }

    async func load(url: string) Result<string, string> {
        let resp = await http.get(url)
        if resp.status != 200 {
            return Err(`bad status: ${resp.status}`)
        }
        return Ok(resp.body)
    }

---

## 12. Full Example

```vertex
import (
    "math.wasm" as mathLib
    "std/fmt"
    "std/http"
)

interface Point {
    x: f32
    y: f32
}

enum Shape {
    Circle(f32)
    Rect(f32, f32)
    Point
}

export func calculateDistance(a: Point, b: Point) f32 {
    let dx = a.x - b.x
    let dy = a.y - b.y
    return mathLib.sqrt((dx * dx) + (dy * dy))
}

func area(s: Shape) f32 {
    return match s {
        Circle(r)    => mathLib.sqrt(r * r) * 3.14159
        Rect(w, h)   => w * h
        Point        => 0.0
    }
}

func safeDivide(a: f32, b: f32) Option<f32> {
    if b == 0.0 {
        return None
    }
    return Some(a / b)
}

async func fetchPoint(url: string) Result<Point, string> {
    let resp = await http.get(url)
    if resp.status != 200 {
        return Err(`bad status: ${resp.status}`)
    }
    return Ok(parsePoint(resp.body))
}

gpu func normalize(x: []f32, y: []f32, out: []f32) {
    let i   = thread.idx
    let len = mathLib.sqrt((x[i] * x[i]) + (y[i] * y[i]))
    out[i]  = x[i] / len
}

async func main() {
    let p1: Point = { x: 10.0, y: 20.0 }
    let p2: Point = { x: 15.0, y: 25.0 }

    let dist = calculateDistance(p1, p2)
    fmt.println(`Local distance: ${dist}`)

    let remote = match await fetchPoint("https://api.example.com/point") {
        Ok(p)    => p
        Err(msg) => {
            fmt.println(`Error: ${msg}`)
            p1
        }
    }

    let remoteDist = calculateDistance(p1, remote)
    fmt.println(`Remote distance: ${remoteDist}`)

    let s: Shape = Circle(5.0)
    fmt.println(`Area: ${area(s)}`)
}
```

---

## 13. Reserved for Future Versions

- `let mut` — mutable bindings (v0.2)
- `enum` with methods (v0.2)
- Closures / first-class functions (v0.3)
- Generics — `Option<T>` and `Result<T,E>` are compiler built-ins until then (v0.3)
- `async` GPU dispatch (v0.3)
- ARM64 Linux target (tracked in wasm-compiler TODO)
- ARM64 Mach-O Darwin target (tracked in wasm-compiler TODO)
- x86_64 Windows PE target (tracked in wasm-compiler TODO)