package operators_test
build test

// === and !== (storage identity, class-only, A.4.5) are deferred to
// 06_composite_types/: constructing a class means calling a declared init
// (A.6.4), which that folder is responsible for pinning first.

func test_eq_int_true() test -> Expected(bool, "1") {
    var a: int32 = 5
    var b: int32 = 5
    return a == b
}

func test_eq_int_false() test -> Expected(bool, "0") {
    var a: int32 = 5
    var b: int32 = 6
    return a == b
}

func test_neq_int() test -> Expected(bool, "1") {
    var a: int32 = 5
    var b: int32 = 6
    return a != b
}

func test_lt_int() test -> Expected(bool, "1") {
    var a: int32 = 3
    var b: int32 = 5
    return a < b
}

func test_gt_int() test -> Expected(bool, "1") {
    var a: int32 = 5
    var b: int32 = 3
    return a > b
}

func test_le_int_equal() test -> Expected(bool, "1") {
    var a: int32 = 5
    var b: int32 = 5
    return a <= b
}

func test_ge_int_equal() test -> Expected(bool, "1") {
    var a: int32 = 5
    var b: int32 = 5
    return a >= b
}

func test_eq_float() test -> Expected(bool, "1") {
    var a: float32 = 3.14
    var b: float32 = 3.14
    return a == b
}

func test_lt_float() test -> Expected(bool, "1") {
    var a: float32 = 1.5
    var b: float32 = 2.5
    return a < b
}

func test_eq_bool_true() test -> Expected(bool, "1") {
    var a: bool = true
    var b: bool = true
    return a == b
}

func test_neq_bool() test -> Expected(bool, "1") {
    var a: bool = true
    var b: bool = false
    return a != b
}

func test_eq_char() test -> Expected(bool, "1") {
    var a: char = 'a'
    var b: char = 'a'
    return a == b
}

func test_lt_char() test -> Expected(bool, "1") {
    // A CharLiteral is a single Unicode scalar value (A.1.5.2); ordering by
    // codepoint is unambiguous, unlike multi-byte string ordering below.
    var a: char = 'a'
    var b: char = 'b'
    return a < b
}

func test_eq_string() test -> Expected(bool, "1") {
    var a: string = "abc"
    var b: string = "abc"
    return a == b
}

func test_neq_string() test -> Expected(bool, "1") {
    var a: string = "abc"
    var b: string = "abd"
    return a != b
}