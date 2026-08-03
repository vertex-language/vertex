package control_flow_test
build test

// Arrays and their index-value form, and maps and their key-value form,
// need [N]T and map[K]V, which live in 06_composite_types/ and 08_heap/
// respectively. This file exercises for-in over ranges and strings only —
// the two iterables available at this tier.
//
// mut/var IterationBinding forms are deferred alongside arrays: in-place
// mutation needs an addressable element, and neither a Range value nor an
// immutable string gives one.

func test_for_in_range_sums() test -> Expected(int32, "10") {
    var sum: int32 = 0
    for i in 1..5 {
        sum += i
    }
    return sum
}

func test_for_in_range_is_half_open() test -> Expected(int32, "4") {
    // 1..5 visits 1, 2, 3, 4 and stops before 5 (A.4.5).
    var last: int32 = 0
    for i in 1..5 {
        last = i
    }
    return last
}

func test_for_in_range_empty_when_bounds_equal() test -> Expected(int32, "0") {
    var count: int32 = 0
    for i in 3..3 {
        count += 1
    }
    return count
}

func test_for_in_range_counts_iterations() test -> Expected(int32, "5") {
    var count: int32 = 0
    for i in 0..5 {
        count += 1
    }
    return count
}

func test_for_in_string_decodes_char_scalars() test -> Expected(char, "c") {
    // String iteration decodes UTF-8 into char scalars (A.5.6).
    var last: char = '\0'
    for c in "abc" {
        last = c
    }
    return last
}

func test_for_in_string_counts_scalars() test -> Expected(int32, "3") {
    var count: int32 = 0
    for c in "abc" {
        count += 1
    }
    return count
}

func test_for_in_string_counts_scalars_not_bytes() test -> Expected(int32, "1") {
    // One multi-byte UTF-8 scalar is still one iteration, not one per
    // encoded byte.
    var count: int32 = 0
    for c in "日" {
        count += 1
    }
    return count
}

func test_for_in_break() test -> Expected(int32, "3") {
    var last: int32 = 0
    for i in 0..100 {
        if i == 3 {
            break
        }
        last = i
    }
    return last
}

func test_for_in_continue() test -> Expected(int32, "6") {
    // Sum 1..5 but skip 2 via continue.
    var sum: int32 = 0
    for i in 1..6 {
        if i == 2 {
            continue
        }
        sum += i
    }
    return sum
}

func test_for_in_blank_binding_discards_value() test -> Expected(int32, "3") {
    // BlankIdentifier is a legal IterationName (A.5.6): only the iteration
    // count is observable, not the value.
    var count: int32 = 0
    for _ in 0..3 {
        count += 1
    }
    return count
}