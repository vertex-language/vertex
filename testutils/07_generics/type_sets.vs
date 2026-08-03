package generics_test
build test

// `|` is union: the type set is every listed type (A.7.3). Exercised via
// pass-through only — see the file header note on why arithmetic isn't
// attempted on a type-set-constrained parameter.
constraint IntOrFloat {
    int32 | float32
}

func passThroughEither[T: IntOrFloat](x: T) -> T {
    return x
}

func test_type_set_union_admits_first_member() test -> Expected(int32, "9") {
    return passThroughEither(9)
}

func test_type_set_union_admits_second_member() test -> Expected(float32, "9.500000") {
    return passThroughEither(9.5)
}

// ~T admits T and every type whose underlying type is T, so an alias to
// float32 still satisfies ~float32; a bare `float32` element would admit
// only float32 exactly (A.7.3). Type aliases (A.6.6) have no folder of
// their own in testutils/README.md — this is the one place their
// semantics become observable, through their interaction with a type set.
type Meters = float32

constraint FloatLike {
    ~float32
}

func identityFloatLike[T: FloatLike](x: T) -> T {
    return x
}

func test_type_set_underlying_type_admits_alias() test -> Expected(float32, "3.500000") {
    let m: Meters = 3.5
    return identityFloatLike(m)
}

func test_type_set_underlying_type_admits_named_type_directly() test -> Expected(float32, "2.000000") {
    return identityFloatLike(2.0)
}