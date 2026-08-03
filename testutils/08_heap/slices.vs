package heap_test
build test

// Two ways a slice value is produced by this grammar:
//   1. PostfixExpression[Expression] "forms a slice when passed a range"
//      (A.4.3) — slicing an existing array. This is the mechanism the
//      grammar states explicitly, so it anchors this file.
//   2. Assigning an ArrayLiteral ([ElementListopt], A.3.1/A.4.7) directly
//      into a []T-typed position. There's no dedicated slice-literal
//      production in the grammar — reading the same bracket-list literal
//      as slice-producing when the destination type is []T rather than
//      [N]T is the most defensible inference available, but it is
//      inference, flagged individually below.

func sumSlice(s: []int32) -> int32 {
    var total: int32 = 0
    for v in s {
        total += v
    }
    return total
}

func test_slice_view_from_array_range_empty() test -> Expected(int32, "0") {
    var arr: [4]int32 = [10, 20, 30, 40]
    let view: []int32 = arr[0..0]
    return sumSlice(view)
}

func test_slice_view_is_half_open_range() test -> Expected(int32, "50") {
    // arr[1..3] views indices 1 and 2 only (half-open, A.4.5): 20 + 30.
    var arr: [4]int32 = [10, 20, 30, 40]
    let view: []int32 = arr[1..3]
    return sumSlice(view)
}

func test_slice_view_full_range() test -> Expected(int32, "100") {
    var arr: [4]int32 = [10, 20, 30, 40]
    let view: []int32 = arr[0..4]
    return sumSlice(view)
}

func test_slice_view_element_access() test -> Expected(int32, "30") {
    var arr: [4]int32 = [10, 20, 30, 40]
    let view: []int32 = arr[1..4]
    return view[1]
}

// []T is mutable and must genuinely duplicate on a bare copy (A.9.4) —
// so element assignment through a slice view writing back into the
// backing array is expected to be legal.
func test_slice_element_assignment_writes_through_to_backing_array() test -> Expected(int32, "99") {
    var arr: [3]int32 = [1, 2, 3]
    var view: []int32 = arr[0..3]
    view[0] = 99
    return arr[0]
}

func test_slice_passed_as_parameter() test -> Expected(int32, "60") {
    var arr: [3]int32 = [10, 20, 30]
    return sumSlice(arr[0..3])
}

// The two-name for-in form binds index-and-value; "dynamic arrays" are
// explicitly named among for-in's iterables alongside fixed arrays
// (A.5.6), and the form is not restated separately for either.
func test_slice_for_in_two_name_form() test -> Expected(int32, "2") {
    var arr: [3]int32 = [10, 20, 30]
    var view: []int32 = arr[0..3]
    var lastIndex: int32 = 0
    for i, _ in view {
        lastIndex = i
    }
    return lastIndex
}

// INFERENCE, unverified against any explicitly-stated grammar mechanism:
// an ArrayLiteral flowing directly into a []T-typed position, with no
// prior array to slice from. Delete this test if the real compiler
// rejects it.
func test_slice_literal_into_slice_typed_binding() test -> Expected(int32, "6") {
    let s: []int32 = [1, 2, 3]
    return sumSlice(s)
}