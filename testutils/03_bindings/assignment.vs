package bindings_test
build test

// AssignTargetList = InitializerList (A.5.2): plain multi-target
// assignment, evaluated as a unit — this is what makes a swap possible
// without a temporary.
func test_assign_swap_via_multi_target() test -> Expected(int32, "21") {
    var a: int32 = 1
    var b: int32 = 2
    a, b = b, a
    return a * 10 + b
}

// BlankIdentifier is legal as an AssignTarget, discarding that position.
func test_assign_blank_identifier_target() test -> Expected(int32, "5") {
    var a: int32 = 0
    _, a = 1, 5
    return a
}

// Compound assignment operators (A.5.2)
func test_compound_add_assign() test -> Expected(int32, "15") {
    var x: int32 = 10
    x += 5
    return x
}

func test_compound_sub_assign() test -> Expected(int32, "5") {
    var x: int32 = 10
    x -= 5
    return x
}

func test_compound_mul_assign() test -> Expected(int32, "50") {
    var x: int32 = 10
    x *= 5
    return x
}

func test_compound_div_assign() test -> Expected(int32, "2") {
    var x: int32 = 10
    x /= 5
    return x
}

func test_compound_mod_assign() test -> Expected(int32, "1") {
    var x: int32 = 11
    x %= 5
    return x
}

func test_compound_and_assign() test -> Expected(int32, "8") {
    var x: int32 = 12  // 0b1100
    x &= 10             // 0b1010
    return x             // 0b1000
}

func test_compound_or_assign() test -> Expected(int32, "14") {
    var x: int32 = 12  // 0b1100
    x |= 10             // 0b1010
    return x             // 0b1110
}

func test_compound_xor_assign() test -> Expected(int32, "6") {
    var x: int32 = 12  // 0b1100
    x ^= 10             // 0b1010
    return x             // 0b0110
}

func test_compound_shl_assign() test -> Expected(int32, "16") {
    var x: int32 = 1
    x <<= 4
    return x
}

func test_compound_shr_assign() test -> Expected(int32, "8") {
    var x: int32 = 64
    x >>= 3
    return x
}

func test_compound_add_assign_float32() test -> Expected(float32, "7.500000") {
    var x: float32 = 5.0
    x += 2.5
    return x
}