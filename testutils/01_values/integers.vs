package literals_test
build test

// Decimal
func test_int_zero() test -> Expected(int32, "0") {
    var x: int32 = 0
    return x
}

func test_int_single_digit() test -> Expected(int32, "7") {
    var x: int32 = 7
    return x
}

func test_int_decimal() test -> Expected(int32, "42") {
    var x: int32 = 42
    return x
}

func test_int_negative() test -> Expected(int32, "-7") {
    // -7 is unary minus applied to the literal 7, folded at compile time (A.1.5.1)
    var x: int32 = -7
    return x
}

func test_int_underscore_separator() test -> Expected(int32, "1000000") {
    var x: int32 = 1_000_000
    return x
}

// Binary
func test_int_binary() test -> Expected(int32, "42") {
    var x: int32 = 0b101010
    return x
}

func test_int_binary_uppercase_prefix() test -> Expected(int32, "42") {
    var x: int32 = 0B101010
    return x
}

func test_int_binary_underscore() test -> Expected(int32, "170") {
    var x: int32 = 0b1010_1010
    return x
}

// Octal
func test_int_octal() test -> Expected(int32, "42") {
    var x: int32 = 0o52
    return x
}

func test_int_octal_uppercase_prefix() test -> Expected(int32, "42") {
    var x: int32 = 0O52
    return x
}

func test_int_octal_underscore() test -> Expected(int32, "1000") {
    var x: int32 = 0o1_750
    return x
}

// Hex
func test_int_hex_upper() test -> Expected(int32, "255") {
    var x: int32 = 0xFF
    return x
}

func test_int_hex_lower() test -> Expected(int32, "255") {
    var x: int32 = 0xff
    return x
}

func test_int_hex_uppercase_prefix() test -> Expected(int32, "255") {
    var x: int32 = 0XFF
    return x
}

func test_int_hex_simple() test -> Expected(int32, "42") {
    var x: int32 = 0x2A
    return x
}

func test_int_hex_underscore() test -> Expected(int32, "4660") {
    var x: int32 = 0x12_34
    return x
}

// An integer literal is untyped until it reaches a typed position, where it
// takes that position's type (A.1.5.1).
func test_int_destination_int8_max() test -> Expected(int8, "127") {
    var x: int8 = 127
    return x
}

func test_int_destination_int8_min() test -> Expected(int8, "-128") {
    var x: int8 = -128
    return x
}

func test_int_destination_uint8_max() test -> Expected(uint8, "255") {
    var x: uint8 = 255
    return x
}

func test_int_destination_byte() test -> Expected(byte, "255") {
    // byte and uint8 denote the same type (A.1.4)
    var x: byte = 255
    return x
}

func test_int_destination_int64() test -> Expected(int64, "1000000000000") {
    var x: int64 = 1_000_000_000_000
    return x
}

func test_int_destination_uint64_max() test -> Expected(uint64, "18446744073709551615") {
    var x: uint64 = 18_446_744_073_709_551_615
    return x
}