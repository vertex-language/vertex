package control_flow_test
build test

// The [~Lit] restriction on the header only matters once a composite or
// map literal exists to parenthesize there (A.4.7); both are deferred to
// 06_composite_types/ and 08_heap/, so that escape hatch isn't testable yet.

func test_if_true_executes_body() test -> Expected(int32, "1") {
    var result: int32 = 0
    if true {
        result = 1
    }
    return result
}

func test_if_false_skips_body() test -> Expected(int32, "0") {
    var result: int32 = 0
    if false {
        result = 1
    }
    return result
}

func test_if_condition_from_comparison() test -> Expected(int32, "1") {
    var result: int32 = 0
    var a: int32 = 5
    var b: int32 = 3
    if a > b {
        result = 1
    }
    return result
}

func test_if_else_takes_if_branch() test -> Expected(int32, "1") {
    var result: int32 = 0
    if true {
        result = 1
    } else {
        result = 2
    }
    return result
}

func test_if_else_takes_else_branch() test -> Expected(int32, "2") {
    var result: int32 = 0
    if false {
        result = 1
    } else {
        result = 2
    }
    return result
}

func test_if_else_if_chain_first_match() test -> Expected(int32, "1") {
    var x: int32 = 1
    var result: int32 = 0
    if x == 1 {
        result = 1
    } else if x == 2 {
        result = 2
    } else {
        result = 3
    }
    return result
}

func test_if_else_if_chain_middle_match() test -> Expected(int32, "2") {
    var x: int32 = 2
    var result: int32 = 0
    if x == 1 {
        result = 1
    } else if x == 2 {
        result = 2
    } else {
        result = 3
    }
    return result
}

func test_if_else_if_chain_falls_to_final_else() test -> Expected(int32, "3") {
    var x: int32 = 99
    var result: int32 = 0
    if x == 1 {
        result = 1
    } else if x == 2 {
        result = 2
    } else {
        result = 3
    }
    return result
}

func test_if_nested() test -> Expected(int32, "1") {
    var result: int32 = 0
    if true {
        if true {
            result = 1
        }
    }
    return result
}

func test_if_condition_from_logical_and() test -> Expected(int32, "1") {
    var result: int32 = 0
    if true && true {
        result = 1
    }
    return result
}