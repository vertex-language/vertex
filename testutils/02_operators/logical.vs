package operators_test
build test

func test_and_both_true() test -> Expected(bool, "1") {
    var a: bool = true
    var b: bool = true
    return a && b
}

func test_and_one_false() test -> Expected(bool, "0") {
    var a: bool = true
    var b: bool = false
    return a && b
}

// Short-circuit, proven without needing functions or observable side
// effects: the plain `+` traps on overflow (A.4.5). If && evaluated both
// operands unconditionally, this test would crash instead of returning "0".
// Also pins the longest-match rule: mistokenized as two `&`, this wouldn't
// type-check (bitwise `&` doesn't accept bool operands).
func test_and_short_circuits() test -> Expected(bool, "0") {
    var a: int8 = 127
    var one: int8 = 1
    return false && (a + one == 0)
}

func test_or_both_false() test -> Expected(bool, "0") {
    var a: bool = false
    var b: bool = false
    return a || b
}

func test_or_one_true() test -> Expected(bool, "1") {
    var a: bool = false
    var b: bool = true
    return a || b
}

func test_or_short_circuits() test -> Expected(bool, "1") {
    var a: int8 = 127
    var one: int8 = 1
    return true || (a + one == 0)
}

func test_not_true() test -> Expected(bool, "0") {
    var a: bool = true
    return !a
}

func test_not_false() test -> Expected(bool, "1") {
    var a: bool = false
    return !a
}