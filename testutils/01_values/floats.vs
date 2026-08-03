package literals_test
build test

func test_float32_simple() test -> Expected(float32, "3.140000") {
    var x: float32 = 3.14
    return x
}

func test_float32_no_fraction() test -> Expected(float32, "2000.000000") {
    // DecimalDigitsWithSeparators ExponentPart — no decimal point required (A.1.5.1)
    var x: float32 = 2e3
    return x
}

func test_float32_exponent_pos() test -> Expected(float32, "125.000000") {
    var x: float32 = 1.25e2
    return x
}

func test_float32_exponent_neg() test -> Expected(float32, "0.012500") {
    var x: float32 = 1.25e-2
    return x
}

func test_float32_exponent_uppercase() test -> Expected(float32, "0.012500") {
    var x: float32 = 1.25E-2
    return x
}

func test_float32_hex_exp_pos() test -> Expected(float32, "60.000000") {
    var x: float32 = 0xFp2
    return x
}

func test_float32_hex_exp_neg() test -> Expected(float32, "3.750000") {
    var x: float32 = 0xFp-2
    return x
}

func test_float32_hex_fraction() test -> Expected(float32, "3.000000") {
    var x: float32 = 0x1.8p1
    return x
}

func test_float32_hex_exp_uppercase_p() test -> Expected(float32, "3.000000") {
    var x: float32 = 0x1.8P1
    return x
}

func test_float64_simple() test -> Expected(float64, "3.140000") {
    var x: float64 = 3.14
    return x
}

func test_float64_exponent() test -> Expected(float64, "125.000000") {
    var x: float64 = 1.25e2
    return x
}

func test_float64_hex_fraction() test -> Expected(float64, "3.000000") {
    var x: float64 = 0x1.8p1
    return x
}