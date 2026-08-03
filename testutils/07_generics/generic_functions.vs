package generics_test
build test

// A bare name is constraint `any` (A.7.1): [T] means [T: any].
func identity[T](x: T) -> T {
    return x
}

func test_generic_bare_name_is_constraint_any() test -> Expected(int32, "3") {
    return identity(3)
}

// Explicit `: any` is the same constraint, spelled out.
func identityExplicitAny[T: any](x: T) -> T {
    return x
}

func test_generic_explicit_any_equivalent_to_bare() test -> Expected(int32, "3") {
    return identityExplicitAny(3)
}

// Under `any`, assignment and argument passing are explicitly licensed
// (A.7.1). `var y: T = x` is a declaration whose right-hand side is an
// owning position (A.9.1), the same shape as an ordinary assignment.
func passThroughAny[T](x: T) -> T {
    var y: T = x
    return y
}

func test_generic_any_permits_assignment_int32() test -> Expected(int32, "8") {
    return passThroughAny(8)
}

func test_generic_any_permits_assignment_string() test -> Expected(string, "ok") {
    return passThroughAny("ok")
}

func test_generic_any_permits_assignment_bool() test -> Expected(bool, "1") {
    return passThroughAny(true)
}