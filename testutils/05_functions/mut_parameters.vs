package functions_test
build test

// mut T lowers to a pointer to the caller's slot (A.3.2); its argument must
// be an addressable var binding. Every test below passes a var — the
// negative case (passing a let) is a compile error and belongs in
// 11_rejected/, not here.

func bump(x: mut int32) {
    x += 1
}

func test_mut_parameter_writes_through() test -> Expected(int32, "6") {
    var n: int32 = 5
    bump(n)
    return n
}

// Distinct point from plain writeback: the parameter starts out holding the
// caller's current value, not a zero value waiting to be filled in.
func addToSelf(x: mut int32) {
    x = x + x
}

func test_mut_parameter_reads_initial_value() test -> Expected(int32, "14") {
    var n: int32 = 7
    addToSelf(n)
    return n
}

func swapValues(a: mut int32, b: mut int32) {
    let tmp = a
    a = b
    b = tmp
}

func test_mut_parameters_swap_caller_bindings() test -> Expected(int32, "21") {
    var x: int32 = 1
    var y: int32 = 2
    swapValues(x, y)
    return x * 10 + y
}

func setDoubled(target: mut int32, source: int32) {
    target = source * 2
}

func test_mut_parameter_combined_with_plain_parameter() test -> Expected(int32, "20") {
    var out: int32 = 0
    setDoubled(out, 10)
    return out
}

func accumulate(x: mut int32, amount: int32) {
    x += amount
}

func test_mut_parameter_sequential_calls_accumulate() test -> Expected(int32, "15") {
    var n: int32 = 0
    accumulate(n, 5)
    accumulate(n, 10)
    return n
}

// A function may mutate through a mut parameter and separately return a
// value in the same call — the two mechanisms don't interfere.
func incrementAndReport(x: mut int32) -> int32 {
    x += 1
    return x
}

func test_mut_parameter_combined_with_return_value() test -> Expected(int32, "12") {
    var n: int32 = 11
    var reported: int32 = incrementAndReport(n)
    return reported
}

func test_mut_parameter_write_through_alongside_return() test -> Expected(int32, "12") {
    var n: int32 = 11
    incrementAndReport(n)
    return n
}