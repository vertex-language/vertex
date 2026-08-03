package control_flow_test
build test

// Defer's effects need a caller-visible slot: return evaluates its
// expression before deferred calls run (see
// test_defer_return_value_captured_before_defer_runs below), so a deferred
// mutation is never observable through that same call's return value. Every
// test here instead threads a `mut` parameter out to the caller and reads
// it after the call returns — the writeback idiom A.4.1 names explicitly.
// Reading/writing that parameter as a plain identifier, with no explicit
// deref, is an assumption on my part: the excerpt names the idiom but never
// shows its in-body access syntax.

func add_to(target: mut int32, amount: int32) {
    target += amount
}

func note(target: mut int32, digit: int32) {
    target = target * 10 + digit
}

func store_squared(target: mut int32, value: int32) {
    target = value * value
}

// --- registration order ---

func deferred_sequence(order: mut int32) {
    defer note(order, 1)
    defer note(order, 2)
    defer note(order, 3)
}

func test_defer_reverse_registration_order() test -> Expected(int32, "321") {
    // Emitted in reverse registration order (A.5.8): 3, then 2, then 1.
    var order: int32 = 0
    deferred_sequence(order)
    return order
}

// --- argument evaluation timing ---

func defer_captures_registration_time_args(result: mut int32) {
    var n: int32 = 4
    defer store_squared(result, n)
    n = 99
}

func test_defer_arguments_evaluated_at_registration() test -> Expected(int32, "16") {
    // Only the call is postponed; arguments are evaluated at registration
    // (A.5.8). If n's later reassignment leaked in, this would be 9801.
    var result: int32 = 0
    defer_captures_registration_time_args(result)
    return result
}

// --- return exit edge ---

func compute_and_defer(x: mut int32) -> int32 {
    defer add_to(x, 100)
    x = 5
    return x
}

func test_defer_return_value_captured_before_defer_runs() test -> Expected(int32, "5") {
    var r: int32 = 0
    return compute_and_defer(r)
}

func test_defer_still_runs_after_return_exit_edge() test -> Expected(int32, "105") {
    var r: int32 = 0
    compute_and_defer(r)
    return r
}

// --- fall-through exit edge ---
// deferred_sequence above has no return statement, so fall-through is its
// only exit edge; test_defer_reverse_registration_order would read 0 if
// defers didn't fire on it.

// --- break exit edge ---

func loop_defer_break(order: mut int32) {
    var i: int32 = 0
    while i < 5 {
        defer add_to(order, 1)
        if i == 2 {
            break
        }
        i += 1
    }
}

func test_defer_runs_on_break() test -> Expected(int32, "3") {
    // Iterations i=0,1,2 each register a defer before the break check; the
    // loop breaks on i==2, and that iteration's defer must still fire.
    // Otherwise this reads 2, not 3.
    var order: int32 = 0
    loop_defer_break(order)
    return order
}

// --- continue exit edge ---

func loop_defer_continue(order: mut int32) {
    var i: int32 = 0
    while i < 3 {
        i += 1
        defer add_to(order, 1)
        if i == 2 {
            continue
        }
    }
}

func test_defer_runs_on_continue() test -> Expected(int32, "3") {
    // Every iteration registers a defer; continuing past i==2 must not
    // skip that iteration's deferred call. Otherwise this reads 2, not 3.
    var order: int32 = 0
    loop_defer_continue(order)
    return order
}