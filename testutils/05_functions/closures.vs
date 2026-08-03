package functions_test
build test

// FunctionExpression captures by value at creation (A.4.1). Assigning to a
// captured binding inside the body is a compile error (11_rejected/); this
// file only pins the legal half — reading a captured value, snapshotted at
// the moment the closure was created.

func test_closure_reads_captured_value() test -> Expected(int32, "5") {
    let x: int32 = 5
    let f = func() -> int32 {
        return x
    }
    return f()
}

func test_closure_captures_snapshot_not_reference() test -> Expected(int32, "5") {
    // x's later reassignment must not leak into the closure's own copy.
    var x: int32 = 5
    let f = func() -> int32 {
        return x
    }
    x = 99
    return f()
}

func test_closure_captures_multiple_bindings() test -> Expected(int32, "23") {
    let a: int32 = 20
    let b: int32 = 3
    let f = func() -> int32 {
        return a + b
    }
    return f()
}

func test_closure_combines_capture_with_parameter() test -> Expected(int32, "17") {
    let offset: int32 = 10
    let f = func(x: int32) -> int32 {
        return x + offset
    }
    return f(7)
}

func makeAdder(amount: int32) -> func(int32) -> int32 {
    return func(x: int32) -> int32 {
        return x + amount
    }
}

func test_closure_returned_from_function() test -> Expected(int32, "15") {
    let addFive = makeAdder(5)
    return addFive(10)
}

func test_closure_returned_from_function_distinct_instances() test -> Expected(int32, "18") {
    // Each call to makeAdder creates its own closure with its own captured
    // amount; the two returned closures must not share state.
    let addFive = makeAdder(5)
    let addTen = makeAdder(10)
    return addFive(3) + addTen(0)
}

func test_immediately_invoked_function_expression() test -> Expected(int32, "42") {
    // A parenthesized FunctionExpression re-enters a PrimaryExpression
    // position (A.4.1), so it can be called directly without a name.
    return (func() -> int32 {
        return 42
    })()
}