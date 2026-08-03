package composite_types_test
build test

class Pixel {
    x: int32
    y: int32
}

// A class is constructed by calling its type name with the initializer's
// arguments (A.6.4) — never with a brace literal. The receiver identifier
// (`p` below) is how the initializer body refers to the instance being
// built; the call site never names it, matching Animal(name: "Rex").
func (p: Pixel) init(x: int32, y: int32) {
    p.x = x
    p.y = y
}

func test_class_construction_via_init() test -> Expected(int32, "3") {
    let px = Pixel(x: 1, y: 2)
    return px.x + px.y
}

func test_class_field_access() test -> Expected(int32, "10") {
    let px = Pixel(x: 10, y: 20)
    return px.x
}

func test_class_instances_independent() test -> Expected(int32, "3") {
    // Two constructed instances must not share storage.
    let a = Pixel(x: 1, y: 2)
    let b = Pixel(x: 100, y: 200)
    return a.x + a.y
}

// "At most one unnamed initializer" is a rule A.8.3 states specifically
// for FOREIGN classes inside a declare block — the exception implies
// ordinary classes may declare more than one unnamed init, distinguished
// by arity.
class Vector {
    dx: int32
    dy: int32
}

func (v: Vector) init(dx: int32, dy: int32) {
    v.dx = dx
    v.dy = dy
}

func (v: Vector) init(uniform: int32) {
    v.dx = uniform
    v.dy = uniform
}

func test_class_overloaded_init_two_args() test -> Expected(int32, "7") {
    let v = Vector(dx: 3, dy: 4)
    return v.dx + v.dy
}

func test_class_overloaded_init_one_arg() test -> Expected(int32, "10") {
    let v = Vector(uniform: 5)
    return v.dx + v.dy
}