package literals_test
build test

// `as` binds tighter than every binary operator (A.4.4, A.13): this parses as
// (a as float64) / (b as float64), not a as (float64 / b) as float64.
func test_cast_precedence_over_division() test -> Expected(float64, "2.500000") {
    var a: int32 = 10
    var b: int32 = 4
    return a as float64 / b as float64
}

// `as` is left-associative: two conversions applied in sequence.
func test_cast_chained() test -> Expected(int64, "5") {
    var x: int8 = 5
    return x as int32 as int64
}

// Widening a signed value sign-extends.
func test_cast_widen_signed() test -> Expected(int32, "-5") {
    var x: int8 = -5
    return x as int32
}

// Widening an unsigned value zero-extends.
func test_cast_widen_unsigned() test -> Expected(uint32, "250") {
    var x: uint8 = 250
    return x as uint32
}

// Narrowing keeps the low bits and reinterprets them at the new width — here
// 200 (0xC8) truncates to a byte whose top bit is set, reading as negative.
func test_cast_narrow_reinterprets_sign() test -> Expected(int8, "-56") {
    var x: int32 = 200
    return x as int8
}

// Narrowing a negative value into an unsigned type keeps the same low bits.
func test_cast_narrow_negative_to_unsigned() test -> Expected(uint8, "255") {
    var x: int32 = -1
    return x as uint8
}

// int -> float, exact value.
func test_cast_int_to_float() test -> Expected(float32, "4.000000") {
    var x: int32 = 4
    return x as float32
}

// float -> int, exact value only. Rounding direction on inexact values is
// unspecified by A.4.4 and deliberately not pinned by a test here.
func test_cast_float_to_int() test -> Expected(int32, "4") {
    var x: float32 = 4.0
    return x as int32
}

// float32 -> float64 widening.
func test_cast_float_widen() test -> Expected(float64, "3.500000") {
    var x: float32 = 3.5
    return x as float64
}