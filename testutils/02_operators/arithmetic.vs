package operators_test
build test

func test_add_int32() test -> Expected(int32, "42") {
    var a: int32 = 20
    var b: int32 = 22
    return a + b
}

func test_sub_int32() test -> Expected(int32, "42") {
    var a: int32 = 50
    var b: int32 = 8
    return a - b
}

func test_mul_int32() test -> Expected(int32, "42") {
    var a: int32 = 6
    var b: int32 = 7
    return a * b
}

func test_div_int32() test -> Expected(int32, "42") {
    // Exact quotient only — truncation direction on a remainder is not
    // pinned by A.4.5 for negative operands, so it is left untested here.
    var a: int32 = 84
    var b: int32 = 2
    return a / b
}

func test_mod_int32() test -> Expected(int32, "2") {
    // Positive operands only; negative-operand sign convention is spec-silent.
    var a: int32 = 17
    var b: int32 = 5
    return a % b
}

func test_unary_minus_int32() test -> Expected(int32, "-42") {
    var x: int32 = 42
    return -x
}

func test_add_float32() test -> Expected(float32, "3.750000") {
    var a: float32 = 1.5
    var b: float32 = 2.25
    return a + b
}

func test_sub_float32() test -> Expected(float32, "3.500000") {
    var a: float32 = 5.5
    var b: float32 = 2.0
    return a - b
}

func test_mul_float32() test -> Expected(float32, "10.000000") {
    var a: float32 = 2.5
    var b: float32 = 4.0
    return a * b
}

func test_div_float32() test -> Expected(float32, "2.500000") {
    var a: float32 = 10.0
    var b: float32 = 4.0
    return a / b
}

func test_unary_minus_float32() test -> Expected(float32, "-3.500000") {
    var x: float32 = 3.5
    return -x
}

// Wrapping forms (A.4.5): "the plain forms trap on overflow" — &+, &-, &*
// are the escape hatch that doesn't. These also pin the longest-match rule
// from the top-level README: if &+ tokenized as `&` then `+`, this wouldn't
// even parse, since there is no unary `+` in the grammar (A.4.4).
func test_wrapping_add_overflow() test -> Expected(int8, "-128") {
    var a: int8 = 127
    var one: int8 = 1
    return a &+ one
}

func test_wrapping_sub_underflow() test -> Expected(int8, "127") {
    var a: int8 = -128
    var one: int8 = 1
    return a &- one
}

func test_wrapping_mul_overflow() test -> Expected(int8, "-56") {
    // 100 * 2 = 200, which wraps to 200 - 256 = -56 in a signed 8-bit lane.
    var a: int8 = 100
    var two: int8 = 2
    return a &* two
}

func test_wrapping_add_no_overflow() test -> Expected(int8, "15") {
    // Sanity check: the wrapping form matches the plain form when there is
    // nothing to wrap.
    var a: int8 = 10
    var b: int8 = 5
    return a &+ b
}