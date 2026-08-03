package ownership_test
build test

// A.9.3's stated rules are almost entirely negative: aliasing one binding
// as two exclusive arguments, reading a binding while exclusively
// accessing it, and overlap through a field path are all compile errors,
// which belong in 11_rejected/. This file pins the positive complement —
// distinct, genuinely non-overlapping bindings and field paths passing the
// exclusivity check and behaving correctly together.

struct Coord {
    x: int32
    y: int32
}

struct Span {
    start: Coord
    end: Coord
}

func combine(a: mut Coord, b: Coord) -> int32 {
    a.x = a.x + b.x
    return a.x
}

func test_exclusivity_distinct_bindings_mut_and_read_together() test -> Expected(int32, "13") {
    var p1: Coord = Coord{x: 3, y: 0}
    let p2 = Coord{x: 10, y: 0}
    return combine(p1, p2)
}

// span.start (mutated) and span.end (read) are disjoint field paths into
// the same struct binding — not overlapping storage, so both may appear
// together in one call.
func test_exclusivity_nonoverlapping_field_paths_mut_and_read() test -> Expected(int32, "21") {
    var span: Span = Span{start: Coord{x: 1, y: 0}, end: Coord{x: 20, y: 0}}
    return combine(span.start, span.end)
}

func test_exclusivity_mut_field_path_leaves_sibling_field_untouched() test -> Expected(int32, "20") {
    var span: Span = Span{start: Coord{x: 1, y: 0}, end: Coord{x: 20, y: 0}}
    combine(span.start, span.end)
    return span.end.x
}

func swapCoords(a: mut Coord, b: mut Coord) {
    let tmp = a
    a = b
    b = tmp
}

// Two mut arguments in the same call, both field paths into the same
// struct — legal precisely because start and end don't overlap. This is
// the direct positive counterpart to "overlap through a field path is
// caught identically" (A.9.3).
func test_exclusivity_two_mut_field_paths_distinct_targets_start() test -> Expected(int32, "20") {
    var span: Span = Span{start: Coord{x: 1, y: 0}, end: Coord{x: 20, y: 0}}
    swapCoords(span.start, span.end)
    return span.start.x
}

func test_exclusivity_two_mut_field_paths_distinct_targets_end() test -> Expected(int32, "1") {
    var span: Span = Span{start: Coord{x: 1, y: 0}, end: Coord{x: 20, y: 0}}
    swapCoords(span.start, span.end)
    return span.end.x
}