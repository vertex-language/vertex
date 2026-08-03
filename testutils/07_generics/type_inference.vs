package generics_test
build test

func firstValue[T](x: T, y: T) -> T {
    return x
}

// Type arguments may be omitted when every type parameter is determined
// by a value argument (A.7.5).
func test_generic_explicit_type_argument() test -> Expected(int32, "5") {
    return firstValue[int32](5, 9)
}

func test_generic_inferred_type_argument() test -> Expected(int32, "5") {
    return firstValue(5, 9)
}

func test_generic_inferred_type_argument_string() test -> Expected(string, "a") {
    return firstValue("a", "b")
}

// A type parameter appearing only in the return type cannot be inferred
// and must be supplied explicitly (A.7.5) — there is no value argument
// here for T to be read from at all.
func zeroValue[T]() -> T {
    var x: T
    return x
}

func test_generic_return_only_type_parameter_requires_explicit_argument() test -> Expected(int32, "0") {
    return zeroValue[int32]()
}

func test_generic_return_only_type_parameter_different_instantiation() test -> Expected(float32, "0.000000") {
    return zeroValue[float32]()
}

// Inference reaches through composite arguments (A.7.5); the excerpt's own
// example is a []T argument fixing T, which needs a slice (08_heap/). This
// extends the same principle to a generic struct argument instead: Container[T]
// fixes T from its own instantiation, no explicit type argument required.
struct Container[T] {
    value: T
}

func unwrap[T](c: Container[T]) -> T {
    return c.value
}

func test_generic_inference_through_composite_argument() test -> Expected(int32, "7") {
    let c: Container[int32] = Container[int32]{value: 7}
    return unwrap(c)
}

func test_generic_inference_through_composite_argument_different_type() test -> Expected(string, "hi") {
    let c: Container[string] = Container[string]{value: "hi"}
    return unwrap(c)
}