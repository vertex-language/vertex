package functions_test
build test

func square(x: int32) -> int32 {
    return x * x
}

func test_single_parameter() test -> Expected(int32, "49") {
    return square(7)
}

func addThree(a: int32, b: int32, c: int32) -> int32 {
    return a + b + c
}

func test_multiple_parameters_same_type() test -> Expected(int32, "6") {
    return addThree(1, 2, 3)
}

func scale(base: float32, factor: int32) -> float32 {
    return base * (factor as float32)
}

func test_parameters_of_mixed_types() test -> Expected(float32, "7.500000") {
    return scale(2.5, 3)
}

func answer() -> int32 {
    return 42
}

func test_zero_parameter_function() test -> Expected(int32, "42") {
    return answer()
}

// Named arguments resolve to positional order at compile time (A.4.3): a
// call binds by declared parameter name, not by the textual order the
// arguments are written in. Using asymmetric operators (subtract, divide)
// makes a wrong binding produce a visibly wrong answer rather than a
// coincidentally correct one.
func subtract(a: int32, b: int32) -> int32 {
    return a - b
}

func test_named_arguments_resolve_to_declared_position() test -> Expected(int32, "7") {
    return subtract(b: 3, a: 10)
}

func test_named_arguments_in_declared_order() test -> Expected(int32, "7") {
    return subtract(a: 10, b: 3)
}

func divideNamed(numerator: int32, denominator: int32) -> int32 {
    return numerator / denominator
}

func test_named_arguments_three_reordered() test -> Expected(int32, "4") {
    return divideNamed(denominator: 2, numerator: 8)
}