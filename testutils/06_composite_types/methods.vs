package composite_types_test
build test

// `shared` receivers are not exercised here: A.6.1 ties `shared` to
// handing out weak back-references, which needs weak()/upgrade() from
// 08_heap/.

struct Coord {
    x: int32
    y: int32
}

// A bare (shared) receiver reads fields without needing an addressable
// var binding.
func (c: Coord) total() -> int32 {
    return c.x + c.y
}

func test_method_bare_receiver_reads_fields() test -> Expected(int32, "7") {
    let c = Coord{x: 3, y: 4}
    return c.total()
}

func test_method_call_on_var_binding() test -> Expected(int32, "7") {
    var c: Coord = Coord{x: 3, y: 4}
    return c.total()
}

// `mut T` in receiver position lowers to a pointer to the caller's slot
// (A.3.2), the same writeback mechanism as a mut parameter (A.4.1); its
// argument must be an addressable var binding.
func (c: mut Coord) scale(factor: int32) {
    c.x = c.x * factor
    c.y = c.y * factor
}

func test_method_mut_receiver_writes_through() test -> Expected(int32, "6") {
    var c: Coord = Coord{x: 1, y: 2}
    c.scale(3)
    return c.x
}

func test_method_mut_receiver_writes_through_both_fields() test -> Expected(int32, "9") {
    var c: Coord = Coord{x: 1, y: 2}
    c.scale(3)
    return c.x + c.y
}

// A method may take ordinary parameters alongside the receiver.
func (c: Coord) distanceSquaredTo(other: Coord) -> int32 {
    let dx: int32 = c.x - other.x
    let dy: int32 = c.y - other.y
    return dx * dx + dy * dy
}

func test_method_with_additional_parameters() test -> Expected(int32, "25") {
    let a = Coord{x: 0, y: 0}
    let b = Coord{x: 3, y: 4}
    return a.distanceSquaredTo(b)
}

// A `var` receiver consumes its receiver unconditionally — there is no
// bare form that copies (A.6.1), the single exception to the
// bare-means-copy rule (A.4.6). Calling one non-destructively means
// copying first into a fresh binding (A.6.4).
func (c: var Coord) consume() -> int32 {
    return c.x + c.y
}

func test_method_var_receiver_consumes_copy_leaves_original_usable() test -> Expected(int32, "7") {
    let original = Coord{x: 3, y: 4}
    var copy: Coord = original
    copy.consume()
    return original.x + original.y
}

func test_method_var_receiver_return_value() test -> Expected(int32, "7") {
    let original = Coord{x: 3, y: 4}
    var copy: Coord = original
    return copy.consume()
}

// Methods on classes use identical receiver-declaration syntax to structs
// (A.6.3): a class differs only in construction and identity, not in its
// method model.
class Counter {
    count: int32
}

func (c: Counter) init(start: int32) {
    c.count = start
}

func (c: mut Counter) increment() {
    c.count = c.count + 1
}

func (c: Counter) value() -> int32 {
    return c.count
}

func test_method_on_class_mut_receiver() test -> Expected(int32, "6") {
    var c: Counter = Counter(start: 5)
    c.increment()
    return c.value()
}