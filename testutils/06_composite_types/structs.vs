package composite_types_test
build test

struct Point {
    x: int32
    y: int32
}

func test_struct_literal_field_access() test -> Expected(int32, "3") {
    let p = Point{x: 1, y: 2}
    return p.x + p.y
}

// FieldValue is named (Identifier : OwningExpression), not positional
// (A.4.7), so field order in the literal doesn't need to match
// declaration order.
func test_struct_literal_fields_out_of_order() test -> Expected(int32, "21") {
    let p = Point{y: 1, x: 20}
    return p.x + p.y
}

struct Config {
    width: int32
    height: int32 = 50
}

// A field default is evaluated at construction for any field the literal
// omits (A.6.2).
func test_struct_field_default_applied_when_omitted() test -> Expected(int32, "50") {
    let c = Config{width: 10}
    return c.height
}

func test_struct_field_default_overridden_when_present() test -> Expected(int32, "99") {
    let c = Config{width: 10, height: 99}
    return c.height
}

func test_struct_explicit_field_not_affected_by_default() test -> Expected(int32, "10") {
    let c = Config{width: 10}
    return c.width
}

// A struct is inline data copied by value (A.6.2): bare assignment (no
// var marker) copies, per A.4.6, so mutating the copy must not affect the
// original.
func test_struct_bare_assignment_copies_independently() test -> Expected(int32, "1") {
    var a: Point = Point{x: 1, y: 1}
    var b: Point = a
    b.x = 99
    return a.x
}

func test_struct_copy_mutation_visible_on_copy() test -> Expected(int32, "99") {
    var a: Point = Point{x: 1, y: 1}
    var b: Point = a
    b.x = 99
    return b.x
}

// CompositeLiteral is grammatical only under [+Lit], which every
// control-flow header clears; parenthesizing re-enters [+Lit] (A.4.7).
func test_struct_literal_in_if_header_requires_parens() test -> Expected(int32, "1") {
    var result: int32 = 0
    if (Point{x: 1, y: 1}).x == 1 {
        result = 1
    }
    return result
}

struct Line {
    start: Point
    end: Point
}

// A struct-typed FieldValue is itself under an owning position (A.9.1),
// and Lit propagates unchanged into it (A.0.3), so no extra parentheses
// are needed around the inner literal.
func test_struct_nested_composite_literal() test -> Expected(int32, "10") {
    let l = Line{start: Point{x: 1, y: 2}, end: Point{x: 3, y: 4}}
    return l.start.x + l.start.y + l.end.x + l.end.y
}

func distanceSquaredFromOrigin(p: Point) -> int32 {
    return p.x * p.x + p.y * p.y
}

func test_struct_passed_as_plain_parameter() test -> Expected(int32, "25") {
    return distanceSquaredFromOrigin(Point{x: 3, y: 4})
}