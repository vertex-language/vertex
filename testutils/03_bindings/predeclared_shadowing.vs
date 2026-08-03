package bindings_test
build test

// Predeclared names (A.1.4) are ordinary identifiers pre-bound in an
// implicit outermost scope — not keywords — so a user declaration shadows
// them freely, except a ReservedBuiltinName (untestable here: shadowing one
// is a compile error, pinned in 11_rejected/).

func test_shadow_predeclared_type_name() test -> Expected(int32, "8") {
    var string: int32 = 8
    return string
}

func test_shadow_predeclared_type_name_int32() test -> Expected(bool, "1") {
    var int32: bool = true
    return int32
}

func test_shadow_predeclared_constraint_any() test -> Expected(int32, "9") {
    var any: int32 = 9
    return any
}

func test_shadow_predeclared_constraint_comparable() test -> Expected(int32, "10") {
    var comparable: int32 = 10
    return comparable
}

func test_predeclared_type_name_unshadowed_elsewhere() test -> Expected(string, "hi") {
    // The shadow of `string` above is local to its own function; this
    // function's scope never saw it.
    var s: string = "hi"
    return s
}