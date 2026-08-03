package operators_test
build test

// `as` vs `/` is already pinned in 01_values/casts.vs
// (test_cast_precedence_over_division). Not repeated here.

// A.4.5: shift sits ABOVE multiplication, a deliberate departure from C.
// Naive C-family precedence would compute 1 << (2*3) = 64; Vertex computes
// (1 << 2) * 3 = 12.
func test_shift_over_multiplicative() test -> Expected(int32, "12") {
    var a: int32 = 1
    var b: int32 = 2
    var c: int32 = 3
    return a << b * c
}

func test_multiplicative_over_additive() test -> Expected(int32, "14") {
    var a: int32 = 2
    var b: int32 = 3
    var c: int32 = 4
    return a + b * c
}

// `&` sits at the multiplicative level, `|` at the additive level (A.13),
// so `&` binds tighter. (a|b)&c would give 1; a|(b&c) gives 5.
func test_bitand_over_bitor() test -> Expected(int32, "5") {
    var a: int32 = 4
    var b: int32 = 3
    var c: int32 = 1
    return a | b & c
}

func test_additive_over_comparison() test -> Expected(bool, "1") {
    var a: int32 = 1
    var b: int32 = 2
    var c: int32 = 3
    return a + b == c
}

func test_comparison_over_logical_and() test -> Expected(bool, "1") {
    var a: int32 = 1
    var b: int32 = 2
    var c: int32 = 3
    var d: int32 = 4
    return a < b && c < d
}

// && binds tighter than ||. Left-to-right without precedence would compute
// (true || false) && false = false; correct grouping gives true || (false
// && false) = true.
func test_and_over_or() test -> Expected(bool, "1") {
    var a: bool = true
    var b: bool = false
    var c: bool = false
    return a || b && c
}

// Every level except `..` is left-associative (A.13). Right-associative
// would give 8 >> (1 >> 1) = 8 >> 0 = 8; left-associative gives 2.
func test_left_associative_shift() test -> Expected(int32, "2") {
    var a: int32 = 8
    return a >> 1 >> 1
}

// `!` only ever takes a UnaryExpression (A.4.4), so it binds to the
// immediate operand, not the whole && expression. !(a && b) would give
// true here; (!a) && b gives false.
func test_not_binds_to_immediate_operand() test -> Expected(bool, "0") {
    var a: bool = false
    var b: bool = false
    return !a && b
}