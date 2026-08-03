package ownership_test
build test

// var T parameters (A.3.2): "whether the callee receives the caller's
// original or a fresh deep copy is decided at the call site by the
// presence or absence of the var marker, never by the declaration." This
// was deferred from 05_functions/ because a scalar-only parameter can't
// show the difference — both copy and transfer are register moves there
// (A.9.4's cost split only exists for a "fat" owning type). A struct is
// the minimal vehicle that makes the call-site choice observable — not
// through mutation (nothing states whether a var-typed parameter binding
// is mutable in-body), but through whether the *original* survives the
// call: reusing a transferred binding is a compile error (A.9.2, tested
// only in 11_rejected/), so proving the original is still readable after
// a bare call is exactly the positive signature of "this was a copy."

struct Spot {
    x: int32
    y: int32
}

struct Segment {
    first: Spot
    second: Spot
}

func sumSpot(p: var Spot) -> int32 {
    return p.x + p.y
}

func test_var_parameter_bare_call_returns_correct_sum() test -> Expected(int32, "3") {
    let a = Spot{x: 1, y: 2}
    return sumSpot(a)
}

func test_var_parameter_bare_call_leaves_original_usable_afterward() test -> Expected(int32, "3") {
    // Bare (no var marker) means copy (A.4.6). If this had transferred a,
    // reading a.x / a.y below would be a compile error instead.
    let a = Spot{x: 1, y: 2}
    let discard = sumSpot(a)
    return a.x + a.y
}

func test_var_parameter_transfer_call_returns_correct_sum() test -> Expected(int32, "3") {
    // var-marked at the call site means transfer (A.4.6). The same
    // declared parameter (var Spot) accepts this call just as readily as
    // the bare one above — the declaration never changes, only the call
    // site does (A.3.2).
    let a = Spot{x: 1, y: 2}
    return sumSpot(var a)
}

func test_var_parameter_transfer_via_field_path_leaves_sibling_usable() test -> Expected(int32, "33") {
    // Transferring seg.first (as the var-typed argument) must not disturb
    // seg.second, a distinct field path that stays fully alive and is
    // itself passed bare (copied) into the same function afterward.
    var seg: Segment = Segment{first: Spot{x: 1, y: 2}, second: Spot{x: 10, y: 20}}
    let firstSum = sumSpot(var seg.first)
    let secondSum = sumSpot(seg.second)
    return firstSum + secondSum
}