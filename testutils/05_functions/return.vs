package functions_test
build test

// Omitting -> Type is the void form; there is no void type name (A.3.4).
// A test function with no result type passes if it compiles and runs
// without crashing (A.12.1) — the third outcome shape, used here for the
// first time in this suite since it's exactly what a void function needs.

func doNothing() {
    return
}

func test_bare_return_from_void_function() test {
    doNothing()
}

func doNothingImplicit() {
}

func test_void_function_falls_off_end_without_return() test {
    doNothingImplicit()
}

func earlyExitVoid(flag: bool, out: mut int32) {
    if flag {
        out = 1
        return
    }
    out = 2
}

func test_bare_return_exits_early_from_void_function() test -> Expected(int32, "1") {
    var result: int32 = 0
    earlyExitVoid(true, result)
    return result
}

func test_void_function_reaches_final_statement_when_not_returned_early() test -> Expected(int32, "2") {
    var result: int32 = 0
    earlyExitVoid(false, result)
    return result
}

func absValue(x: int32) -> int32 {
    if x < 0 {
        return -x
    }
    return x
}

func test_early_return_short_circuits_remaining_body() test -> Expected(int32, "5") {
    return absValue(-5)
}

func test_return_reaches_final_statement_when_early_branch_skipped() test -> Expected(int32, "5") {
    return absValue(5)
}

// Multi-value return (A.5.3) is written bare, comma-separated, with no
// wrapping parentheses — the signature's -> (A, B) is where the
// parentheses are required, since there they describe the return type's
// shape (A.3.1). This uses a bare tuple only as the minimal vehicle to
// observe the return rule, the same way 03_bindings/let_var.vs used one to
// observe destructuring; full tuple coverage (nested tuples, .0/.1 access,
// the one-element trailing-comma rule) is 06_composite_types/ territory.
func minMax(a: int32, b: int32) -> (int32, int32) {
    if a < b {
        return a, b
    }
    return b, a
}

func test_multi_value_return_low() test -> Expected(int32, "4") {
    let lo, hi = minMax(9, 4)
    return lo
}

func test_multi_value_return_high() test -> Expected(int32, "9") {
    let lo, hi = minMax(9, 4)
    return hi
}

func spread(n: int32) -> (int32, int32, int32) {
    return n, n * 2, n * 3
}

func test_multi_value_return_three_values() test -> Expected(int32, "123") {
    let a, b, c = spread(1)
    return a * 100 + b * 10 + c
}