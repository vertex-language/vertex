package ownership_test
build test

// TransferTarget's grammar (A.4.6) is Identifier | TransferTarget.Identifier
// only — never an arbitrary expression, and never a freshly constructed
// literal. That's why "the marker is meaningless on a temporary" (A.4.6) is
// a fact about *omitting* var on a fresh literal, not something a var-
// prefixed test could exercise: `var Point{...}` isn't grammatical to begin
// with.

struct Point {
    x: int32
    y: int32
}

struct Pair {
    first: Point
    second: Point
}

struct Nested {
    inner: Pair
}

struct Box {
    contents: Point
}

// --- VariableDeclaration RHS ---

func test_transfer_marker_on_plain_binding() test -> Expected(int32, "3") {
    let a = Point{x: 1, y: 2}
    let b = var a
    return b.x + b.y
}

// --- AssignmentStatement RHS (A.9.1 lists this separately from decl RHS) ---

func test_transfer_marker_on_assignment_statement_rhs() test -> Expected(int32, "3") {
    let a = Point{x: 1, y: 2}
    var b: Point = Point{x: 0, y: 0}
    b = var a
    return b.x + b.y
}

// --- Field-path TransferTarget (partial move) ---

func test_transfer_marker_on_field_path() test -> Expected(int32, "3") {
    var pr: Pair = Pair{first: Point{x: 1, y: 2}, second: Point{x: 10, y: 20}}
    let moved = var pr.first
    return moved.x + moved.y
}

func test_transfer_marker_on_field_path_leaves_sibling_field_usable() test -> Expected(int32, "30") {
    // pr.first and pr.second are distinct field paths (A.4.6); moving one
    // must not disturb the other.
    var pr: Pair = Pair{first: Point{x: 1, y: 2}, second: Point{x: 10, y: 20}}
    let moved = var pr.first
    return pr.second.x + pr.second.y
}

// TransferTarget is recursive (TransferTarget . Identifier), so an
// arbitrarily deep path is legal, not just a single hop.
func test_transfer_marker_on_two_level_field_path() test -> Expected(int32, "3") {
    var n: Nested = Nested{inner: Pair{first: Point{x: 1, y: 2}, second: Point{x: 9, y: 9}}}
    let moved = var n.inner.first
    return moved.x + moved.y
}

// --- Element of a composite / array / tuple / map literal (A.9.1, A.4.7) ---

func test_transfer_marker_as_composite_literal_field_value() test -> Expected(int32, "3") {
    let a = Point{x: 1, y: 2}
    let boxed = Box{contents: var a}
    return boxed.contents.x + boxed.contents.y
}

func test_transfer_marker_as_array_literal_element() test -> Expected(int32, "3") {
    let a = Point{x: 1, y: 2}
    let pts: [2]Point = [var a, Point{x: 0, y: 0}]
    return pts[0].x + pts[0].y
}

func test_transfer_marker_as_tuple_literal_element() test -> Expected(int32, "3") {
    let a = Point{x: 1, y: 2}
    let t = (var a, 99)
    return t.0.x + t.0.y
}

func test_transfer_marker_as_map_literal_value() test -> Expected(int32, "3") {
    let a = Point{x: 1, y: 2}
    let m: map[string]Point = {"p": var a}
    return m["p"].x + m["p"].y
}

// --- Returned expression (A.5.3, A.9.1) ---

func makeAndTransfer() -> Point {
    let a = Point{x: 1, y: 2}
    return var a
}

func test_transfer_marker_on_return_expression() test -> Expected(int32, "3") {
    let p = makeAndTransfer()
    return p.x + p.y
}

// --- Own re-established independently at a nested owning position ---
// (A.9.1: "Own does not propagate into subexpressions... re-established at
// each listed position") — the var here sits inside a composite-literal
// element, itself inside a plain (non-var) call argument.

func sumBoxed(b: Box) -> int32 {
    return b.contents.x + b.contents.y
}

func test_transfer_marker_nested_owning_position_inside_call_argument() test -> Expected(int32, "3") {
    let a = Point{x: 1, y: 2}
    return sumBoxed(Box{contents: var a})
}

// The remaining A.9.1 owning position — "the binding of a consuming for
// loop" — is already covered by 06_composite_types/arrays.vs's
// `for var v in arr` test and isn't repeated here.