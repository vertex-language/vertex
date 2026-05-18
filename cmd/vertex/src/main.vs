
// ── WASM host imports ─────────────────────────────────────────────────────────
@extern("env", "print_i32")
func print_i32(_ value: Int)

@extern("env", "print_f64")
func print_f64(_ value: Double)

// ── libc imports (native target) ──────────────────────────────────────────────
// "printf__vararg_i" strips to real symbol "printf" in the x86_64 backend.
// The "_i" suffix marks the single variadic argument as an int (not a pointer),
// so the backend does not apply the R15 wasm-offset translation to it.
// The format string pointer IS translated via the pointerArgs map in main.go.
@extern("c", "printf__vararg_i")
func printf_i(_ format: String, _ value: Int)

// puts(3) writes a string then a newline — simpler than printf for headings.
@extern("c", "puts")
func puts(_ s: String)

// ── fibonacci ─────────────────────────────────────────────────────────────────
func fibonacci(_ n: Int) -> Int {
    if n <= 1 {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

// ── factorial ─────────────────────────────────────────────────────────────────
func factorial(_ n: Int) -> Int {
    var result = 1
    var i = 1
    while i <= n {
        result = result * i
        i = i + 1
    }
    return result
}

// ── power ─────────────────────────────────────────────────────────────────────
func power(_ base: Int, _ exp: Int) -> Int {
    var result = 1
    var e = exp
    while e != 0 {
        result = result * base
        e = e - 1
    }
    return result
}

// ── gcd ───────────────────────────────────────────────────────────────────────
func gcd(_ a: Int, _ b: Int) -> Int {
    var x = a
    var y = b
    while y != 0 {
        let tmp = y
        y = x - (x / y) * y
        x = tmp
    }
    return x
}

// ── abs ───────────────────────────────────────────────────────────────────────
func abs(_ n: Int) -> Int {
    if n <= -1 {
        return 0 - n
    }
    return n
}

// ── entry point ───────────────────────────────────────────────────────────────
func main() -> Int {
    puts("=== fibonacci ===")
    printf_i("fib(0)  = %d\n", fibonacci(0))
    printf_i("fib(1)  = %d\n", fibonacci(1))
    printf_i("fib(5)  = %d\n", fibonacci(5))
    printf_i("fib(10) = %d\n", fibonacci(10))
    printf_i("fib(15) = %d\n", fibonacci(15))

    puts("=== factorial ===")
    printf_i("fact(1)  = %d\n", factorial(1))
    printf_i("fact(5)  = %d\n", factorial(5))
    printf_i("fact(10) = %d\n", factorial(10))

    puts("=== power ===")
    printf_i("2^8  = %d\n", power(2, 8))
    printf_i("3^5  = %d\n", power(3, 5))
    printf_i("10^3 = %d\n", power(10, 3))

    puts("=== gcd ===")
    printf_i("gcd(48,18)  = %d\n", gcd(48, 18))
    printf_i("gcd(100,75) = %d\n", gcd(100, 75))
    printf_i("gcd(17,13)  = %d\n", gcd(17, 13))

    puts("=== abs ===")
    printf_i("abs(-42) = %d\n", abs(-42))
    printf_i("abs(7)   = %d\n", abs(7))

    return 0
}