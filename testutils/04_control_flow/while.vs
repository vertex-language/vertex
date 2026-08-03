package control_flow_test
build test

func test_while_basic_count() test -> Expected(int32, "5") {
    var i: int32 = 0
    while i < 5 {
        i += 1
    }
    return i
}

func test_while_false_never_executes() test -> Expected(int32, "0") {
    var i: int32 = 0
    while false {
        i = 99
    }
    return i
}

func test_while_accumulates() test -> Expected(int32, "10") {
    var i: int32 = 0
    var sum: int32 = 0
    while i < 4 {
        i += 1
        sum += i
    }
    return sum
}

func test_while_break_stops_early() test -> Expected(int32, "3") {
    var i: int32 = 0
    while i < 100 {
        if i == 3 {
            break
        }
        i += 1
    }
    return i
}

func test_while_continue_skips_rest_of_body() test -> Expected(int32, "6") {
    // Sum 1..5 but skip adding when i == 2 — continue jumps straight back
    // to the condition check, past the sum += below it.
    var i: int32 = 0
    var sum: int32 = 0
    while i < 5 {
        i += 1
        if i == 2 {
            continue
        }
        sum += i
    }
    return sum
}

// There are no loop labels (A.5.9): a nested break only ever exits its own
// innermost while.
func test_while_nested_break_exits_only_inner_loop() test -> Expected(int32, "30") {
    var outer: int32 = 0
    var total: int32 = 0
    while outer < 3 {
        var inner: int32 = 0
        while inner < 100 {
            if inner == 10 {
                break
            }
            inner += 1
        }
        total += inner
        outer += 1
    }
    return total
}

func test_while_multi_level_exit_via_flag() test -> Expected(int32, "1") {
    // The explicit-flag idiom A.5.9 points to in place of loop labels.
    var done: bool = false
    var outer: int32 = 0
    while outer < 5 && !done {
        var inner: int32 = 0
        while inner < 5 {
            if outer == 1 && inner == 1 {
                done = true
                break
            }
            inner += 1
        }
        if !done {
            outer += 1
        }
    }
    return outer
}