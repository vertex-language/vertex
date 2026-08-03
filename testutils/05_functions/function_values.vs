package functions_test
build test

// A non-capturing function value is one word — a bare code pointer
// (A.3.4). Referencing a declared function by name without a trailing call
// evaluates to that value, same as referencing any other typed binding.

func double(x: int32) -> int32 {
    return x * 2
}

func increment(x: int32) -> int32 {
    return x + 1
}

func applyTwice(f: func(int32) -> int32, x: int32) -> int32 {
    return f(f(x))
}

func test_named_function_passed_as_value() test -> Expected(int32, "20") {
    return applyTwice(double, 5)
}

func test_different_function_value_changes_result() test -> Expected(int32, "7") {
    return applyTwice(increment, 5)
}

func apply(f: func(int32) -> int32, x: int32) -> int32 {
    return f(x)
}

func test_function_value_stored_in_binding() test -> Expected(int32, "10") {
    let op: func(int32) -> int32 = double
    return apply(op, 5)
}

func choose(useDouble: bool) -> func(int32) -> int32 {
    if useDouble {
        return double
    }
    return increment
}

func test_function_value_returned_from_function() test -> Expected(int32, "10") {
    let op = choose(true)
    return op(5)
}

func test_function_value_returned_from_function_other_branch() test -> Expected(int32, "6") {
    let op = choose(false)
    return op(5)
}

func combine(f: func(int32) -> int32, g: func(int32) -> int32, x: int32) -> int32 {
    return f(x) + g(x)
}

func test_two_function_values_in_one_call() test -> Expected(int32, "16") {
    // double(5) + increment(5) = 10 + 6
    return combine(double, increment, 5)
}