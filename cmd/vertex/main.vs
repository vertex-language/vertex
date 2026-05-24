package main
build linux
import "lib/c"

class C : c {
    func printf(fmt: any char, ...) -> int
    func puts(s: any char) -> int
}

// ── fibonacci ─────────────────────────────────────────────────────────────────
func fibonacci(n: int) -> int {
    if n <= 1 {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

// ── factorial ─────────────────────────────────────────────────────────────────
func factorial(n: int) -> int {
    var result = 1
    var i = 1
    while i <= n {
        result = result * i
        i = i + 1
    }
    return result
}

// ── power ─────────────────────────────────────────────────────────────────────
func power(base: int, exp: int) -> int {
    var result = 1
    var e = exp
    while e != 0 {
        result = result * base
        e = e - 1
    }
    return result
}

// ── gcd ───────────────────────────────────────────────────────────────────────
func gcd(a: int, b: int) -> int {
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
func abs(n: int) -> int {
    if n <= -1 {
        return 0 - n
    }
    return n
}

// ── entry point ───────────────────────────────────────────────────────────────
func main() -> int {
    var c = C()

    c.printf("fib(0)  = %d\n".any(), fibonacci(0))
    c.printf("fib(1)  = %d\n".any(), fibonacci(1))
    c.printf("fib(5)  = %d\n".any(), fibonacci(5))
    c.printf("fib(10) = %d\n".any(), fibonacci(10))
    c.printf("fib(15) = %d\n".any(), fibonacci(15))

    c.printf("fact(1)  = %d\n".any(), factorial(1))
    c.printf("fact(5)  = %d\n".any(), factorial(5))
    c.printf("fact(10) = %d\n".any(), factorial(10))

    c.printf("2^8  = %d\n".any(), power(2, 8))
    c.printf("3^5  = %d\n".any(), power(3, 5))
    c.printf("10^3 = %d\n".any(), power(10, 3))

    c.printf("gcd(48,18)  = %d\n".any(), gcd(48, 18))
    c.printf("gcd(100,75) = %d\n".any(), gcd(100, 75))
    c.printf("gcd(17,13)  = %d\n".any(), gcd(17, 13))

    c.printf("abs(-42) = %d\n".any(), abs(-42))
    c.printf("abs(7)   = %d\n".any(), abs(7))

    return 0
}