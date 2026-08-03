package functions_test
build test

// Recursion and mutual recursion exercise ordinary call/return mechanics
// already pinned elsewhere in this folder — they aren't a standalone rule
// with their own citation. isOdd being called before its own declaration
// below relies on top-level declarations being order-independent (A.2),
// but that broader claim has no return-value signature of its own and is
// not separately asserted here (see testutils/README.md, "What has no
// folder, and why"); it's just along for the ride.

func factorial(n: int32) -> int32 {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func test_recursive_function_computes_factorial() test -> Expected(int32, "120") {
    return factorial(5)
}

func test_recursive_function_base_case() test -> Expected(int32, "1") {
    return factorial(1)
}

func fibonacci(n: int32) -> int32 {
    if n < 2 {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

func test_recursive_function_computes_fibonacci() test -> Expected(int32, "13") {
    return fibonacci(7)
}

func isEven(n: int32) -> bool {
    if n == 0 {
        return true
    }
    return isOdd(n - 1)
}

func isOdd(n: int32) -> bool {
    if n == 0 {
        return false
    }
    return isEven(n - 1)
}

func test_mutual_recursion_even() test -> Expected(bool, "1") {
    return isEven(4)
}

func test_mutual_recursion_odd() test -> Expected(bool, "1") {
    return isOdd(7)
}