# Vertex Language

Vertex is a statically typed programming language designed to compile directly to WebAssembly (`.wasm`) or natively to `x86_64` ELF binaries. It features a clean, Swift-inspired syntax, support for value-type structs and reference-type classes, and seamless interoperability with native C/WASM host environments via `@extern` functions.

---

## Installation

You can install the Vertex compiler using Go:

```bash
go install github.com/vertex-language/vertex/cmd/vertex@latest

```

---

## Usage

Use the `vertex` CLI to compile your `.vs` source files:

```bash
vertex [-o output] [-wasm] [-v] [-entry sym] <source.vs>

```

### Options

* **`-wasm`**: Emits a WebAssembly (`.wasm`) binary. If omitted, the compiler defaults to generating a native `x86_64` ELF binary.
* **`-o <output>`**: Specifies the output file name. By default, it derives the name from the input file (e.g., `main.vs` becomes `main` or `main.wasm`).
* **`-entry <sym>`**: Sets the ELF entry-point export symbol for native binaries (defaults to `main`).
* **`-v`**: Enables verbose compilation output, detailing AST passes, WebAssembly function indexes, and generated locals.

---

## Language Features

### Type System

Vertex is strictly typed and includes the following built-in primitives:

* **Integers:** `Int` (default 32-bit), `Int64`, `UInt`
* **Floating Point:** `Float` (32-bit), `Double` (64-bit)
* **Other:** `Bool`, `String` (NUL-terminated byte pointers), and `Void`
* **Composite:** `struct` (pass-by-value) and `class` (pass-by-reference/heap-allocated).

### Control Flow

Vertex supports standard control flow mechanisms:

* `if` / `else if` / `else`
* `while` loops
* `repeat { ... } while` loops
* `for i in lo..<hi` half-open range loops

### Foreign Function Interface (FFI)

You can easily bind external WASM host functions or `libc` functions (when targeting native) using the `@extern` attribute:

```swift
// WASM host import
@extern("env", "print_i32")
func print_i32(_ value: Int)

// libc import (Native target)
@extern("c", "puts")
func puts(_ s: String)

```

---

## Example: `main.vs`

Here is a quick example of Vertex code demonstrating recursion, loops, and external C bindings:

```swift
@extern("c", "printf__vararg_i")
func printf_i(_ format: String, _ value: Int)

@extern("c", "puts")
func puts(_ s: String)

// Recursive Fibonacci
func fibonacci(_ n: Int) -> Int {
    if n <= 1 {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

// While loop Factorial
func factorial(_ n: Int) -> Int {
    var result = 1
    var i = 1
    while i <= n {
        result = result * i
        i = i + 1
    }
    return result
}

func main() -> Int {
    puts("=== fibonacci ===")
    printf_i("fib(10) = %d\n", fibonacci(10))

    puts("=== factorial ===")
    printf_i("fact(5)  = %d\n", factorial(5))

    return 0
}

```

To compile and run this natively:

```bash
$ vertex main.vs
$ ./main
=== fibonacci ===
fib(10) = 55
=== factorial ===
fact(5)  = 120

```