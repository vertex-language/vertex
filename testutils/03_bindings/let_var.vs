package bindings_test
build test

// let vs var (A.5.1): let is immutable, var is mutable and owns a real
// stack slot for its whole lifetime. Attempting to assign to a `let`
// binding is a compile error and belongs in 11_rejected/, not here — this
// file only pins the positive half: var accepts reassignment.
func test_var_reassignment() test -> Expected(int32, "99") {
    var x: int32 = 1
    x = 99
    return x
}

func test_let_single_binding() test -> Expected(int32, "42") {
    let x: int32 = 42
    return x
}

// Multi-binding with matching initializer count is parallel declaration,
// not destructuring — each Binding pairs positionally with its own
// initializer in InitializerList.
func test_let_parallel_declaration() test -> Expected(int32, "3") {
    let a, b = 1, 2
    return a + b
}

func test_var_parallel_declaration() test -> Expected(int32, "30") {
    var a, b = 10, 20
    return a + b
}

func test_let_multi_binding_explicit_types() test -> Expected(int32, "30") {
    let a: int32, b: int32 = 10, 20
    return a + b
}

// Multi-binding with a single initializer is a tuple destructure: the one
// initializer is a tuple value, and the binding list unpacks it.
func test_let_tuple_destructure() test -> Expected(int32, "3") {
    let a, b = (1, 2)
    return a + b
}

func test_var_tuple_destructure_then_mutate() test -> Expected(int32, "10") {
    // Bindings introduced by a var destructure are independently mutable,
    // same as any other var binding.
    var a, b = (3, 4)
    a = 10
    return a
}

func test_let_discard_in_parallel_declaration() test -> Expected(int32, "2") {
    let _, y = 1, 2
    return y
}