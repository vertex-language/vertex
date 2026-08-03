package composite_types_test
build test

func test_array_literal_first_element() test -> Expected(int32, "1") {
    var arr: [3]int32 = [1, 2, 3]
    return arr[0]
}

func test_array_literal_element_access() test -> Expected(int32, "2") {
    var arr: [3]int32 = [1, 2, 3]
    return arr[1]
}

func test_array_literal_last_element() test -> Expected(int32, "3") {
    var arr: [3]int32 = [1, 2, 3]
    return arr[2]
}

func test_array_element_assignment() test -> Expected(int32, "99") {
    var arr: [3]int32 = [1, 2, 3]
    arr[0] = 99
    return arr[0]
}

func test_array_element_assignment_leaves_others_untouched() test -> Expected(int32, "2") {
    var arr: [3]int32 = [1, 2, 3]
    arr[0] = 99
    return arr[1]
}

// [N]T's Type may itself be an ArrayType (A.3.1): inline storage nests.
func test_array_nested() test -> Expected(int32, "6") {
    var grid: [2][3]int32 = [[1, 2, 3], [4, 5, 6]]
    return grid[1][2]
}

// ArrayLength may be an Identifier denoting a compile-time constant
// (A.3.1). The only compile-time-constant binding form this suite has
// established is a top-level `let` with a compile-time-evaluable
// initializer (A.2). This is the most defensible reading available, but
// the excerpt never shows this exact spelling, so treat it as inference.
let ArraySize: int32 = 4

func test_array_length_via_top_level_constant_identifier() test -> Expected(int32, "40") {
    var arr: [ArraySize]int32 = [10, 20, 30, 40]
    return arr[3]
}

// The two-name for-in form binds index-and-value over arrays (A.5.6).
func test_array_for_in_two_name_form_sums_values() test -> Expected(int32, "60") {
    var arr: [3]int32 = [10, 20, 30]
    var sum: int32 = 0
    for _, v in arr {
        sum += v
    }
    return sum
}

// The index binding is assumed int32-compatible here — not stated
// explicitly anywhere in the excerpt, but consistent with int32 being the
// default integer type used throughout this suite.
func test_array_for_in_two_name_form_uses_index() test -> Expected(int32, "3") {
    var arr: [3]int32 = [10, 20, 30]
    var sum: int32 = 0
    for i, _ in arr {
        sum += i
    }
    return sum
}

// `mut` iterates by exclusive access, permitting in-place mutation
// (A.5.6) — the same writeback mechanism as a mut parameter (A.4.1).
func test_array_for_in_mut_mutates_in_place() test -> Expected(int32, "10") {
    var arr: [3]int32 = [1, 2, 3]
    for mut v in arr {
        v = v * 10
    }
    return arr[0]
}

func test_array_for_in_mut_mutates_every_element() test -> Expected(int32, "60") {
    var arr: [3]int32 = [1, 2, 3]
    for mut v in arr {
        v = v * 10
    }
    return arr[0] + arr[1] + arr[2]
}

// `var` is the consuming form: each element moves into the loop binding
// (A.5.6). For a scalar element type this is behaviorally identical to the
// bare form — true divergence needs an owning fat element type, which is
// 09_ownership/ territory — so this only pins that the form compiles and
// reads correctly.
func test_array_for_in_var_consuming_form() test -> Expected(int32, "6") {
    var arr: [3]int32 = [1, 2, 3]
    var sum: int32 = 0
    for var v in arr {
        sum += v
    }
    return sum
}

// The bare form iterates by shared access and needs no addressable
// container.
func test_array_for_in_bare_form_over_let_binding() test -> Expected(int32, "6") {
    let arr: [3]int32 = [1, 2, 3]
    var sum: int32 = 0
    for v in arr {
        sum += v
    }
    return sum
}