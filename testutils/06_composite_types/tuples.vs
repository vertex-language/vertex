package composite_types_test
build test

func test_tuple_positional_access() test -> Expected(int32, "6") {
    let t = (1, 2, 3)
    return t.0 + t.1 + t.2
}

// Chained positional access into a nested tuple (A.4.3).
func test_tuple_nested_positional_access() test -> Expected(int32, "5") {
    let t = ((2, 3), 10)
    return t.0.0 + t.0.1
}

// A one-element tuple literal requires its trailing comma (A.4.7);
// (5) alone would just be a parenthesized integer, not a tuple.
func test_tuple_one_element_requires_trailing_comma() test -> Expected(int32, "5") {
    let t = (5,)
    return t.0
}

func sumPair(pair: (int32, int32)) -> int32 {
    return pair.0 + pair.1
}

func test_tuple_type_as_parameter() test -> Expected(int32, "9") {
    return sumPair((4, 5))
}

// TupleElement may name its element (Identifier : Type), but the only
// access syntax the grammar gives is positional (.DecimalDigit, A.4.3);
// the name is documentation on the type, not an alternate accessor.
func namedTupleSum(point: (x: int32, y: int32)) -> int32 {
    return point.0 + point.1
}

func test_tuple_named_element_type_still_accessed_positionally() test -> Expected(int32, "12") {
    return namedTupleSum((x: 5, y: 7))
}

// Multi-value return is written bare, comma-separated, with no wrapping
// parentheses (A.5.3); the signature's own -> (A, B) is where parens are
// required, since there they describe the return type's shape.
func namedReturnPair() -> (x: int32, y: int32) {
    return 3, 4
}

func test_tuple_named_return_type_bare_return_statement() test -> Expected(int32, "7") {
    let p = namedReturnPair()
    return p.0 + p.1
}

// Every element of a tuple literal is an owning position (A.9.1);
// PostfixExpression . DecimalDigit is a legal AssignTarget when the tuple
// binding is addressable (var, not let).
func test_tuple_field_assignment_through_var_binding() test -> Expected(int32, "99") {
    var t: (int32, int32) = (1, 2)
    t.0 = 99
    return t.0
}

func test_tuple_field_assignment_leaves_other_field_untouched() test -> Expected(int32, "2") {
    var t: (int32, int32) = (1, 2)
    t.0 = 99
    return t.1
}

// The (T, string) boundary-tuple shape (A.8.4's convention) used natively,
// with no FFI boundary involved — A.8.4 itself notes this is "the same
// tuple shape as any native fallible function." On success the value half
// is meaningful and the string is empty; on failure the value half is the
// zero value and the string carries a message.
func safeDivide(a: int32, b: int32) -> (int32, string) {
    if b == 0 {
        return 0, "division by zero"
    }
    return a / b, ""
}

func test_tuple_boundary_shape_success_value() test -> Expected(int32, "5") {
    let quotient, _ = safeDivide(10, 2)
    return quotient
}

func test_tuple_boundary_shape_success_empty_error() test -> Expected(bool, "1") {
    let _, err = safeDivide(10, 2)
    return err == ""
}

func test_tuple_boundary_shape_failure_zero_value() test -> Expected(int32, "0") {
    let quotient, _ = safeDivide(10, 0)
    return quotient
}

func test_tuple_boundary_shape_failure_nonempty_error() test -> Expected(bool, "1") {
    let _, err = safeDivide(10, 0)
    return err != ""
}