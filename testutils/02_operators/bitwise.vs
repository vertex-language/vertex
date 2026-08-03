package operators_test
build test

func test_bitwise_and() test -> Expected(int32, "8") {
    var a: int32 = 12  // 0b1100
    var b: int32 = 10  // 0b1010
    return a & b       // 0b1000
}

func test_bitwise_or() test -> Expected(int32, "14") {
    var a: int32 = 12  // 0b1100
    var b: int32 = 10  // 0b1010
    return a | b       // 0b1110
}

func test_bitwise_xor() test -> Expected(int32, "6") {
    var a: int32 = 12  // 0b1100
    var b: int32 = 10  // 0b1010
    return a ^ b       // 0b0110
}

func test_bitwise_not_zero() test -> Expected(int32, "-1") {
    var x: int32 = 0
    return ~x
}

func test_bitwise_not_nonzero() test -> Expected(int32, "-6") {
    var x: int32 = 5
    return ~x
}

func test_shift_left() test -> Expected(int32, "16") {
    var x: int32 = 1
    return x << 4
}

func test_shift_right() test -> Expected(int32, "8") {
    var x: int32 = 64
    return x >> 3
}

func test_shift_left_unsigned() test -> Expected(uint8, "128") {
    var x: uint8 = 1
    return x << 7
}