package generics_test
build test

// A MethodRequirement is satisfied by any type declaring a matching
// receiver method; every instantiation is monomorphized, so the call
// inside the generic body lowers to a direct call on the concrete type
// (A.7.2) — no vtable, no interface value.

constraint Describable {
    func describe() -> int32
}

struct Circle {
    radius: int32
}

func (c: Circle) describe() -> int32 {
    return c.radius * 100
}

struct Square {
    side: int32
}

func (s: Square) describe() -> int32 {
    return s.side * 10
}

func describeValue[T: Describable](x: T) -> int32 {
    return x.describe()
}

func test_constraint_method_requirement_satisfied_by_struct() test -> Expected(int32, "200") {
    return describeValue(Circle{radius: 2})
}

func test_constraint_method_requirement_different_concrete_type() test -> Expected(int32, "30") {
    return describeValue(Square{side: 3})
}

// A bare ConstraintName element embeds that constraint's set (A.7.2);
// multiple elements in a constraint body form an intersection, so
// FullyDescribable requires both describe() and label().
constraint Nameable {
    func label() -> string
}

constraint FullyDescribable {
    Describable
    Nameable
}

struct Widget {
    id: int32
}

func (w: Widget) describe() -> int32 {
    return w.id
}

func (w: Widget) label() -> string {
    return "widget"
}

func combinedScore[T: FullyDescribable](x: T) -> int32 {
    let n = x.label()
    if n == "widget" {
        return x.describe() + 1
    }
    return x.describe()
}

func test_constraint_embedding_grants_both_requirements() test -> Expected(int32, "6") {
    return combinedScore(Widget{id: 5})
}

// comparable admits every type supporting ==/!= (A.7.4). Under bare `any`
// this operator would be unavailable (A.7.1); `comparable` is what
// licenses it.
func areEqual[T: comparable](a: T, b: T) -> bool {
    return a == b
}

func test_constraint_comparable_licenses_equality_int32() test -> Expected(bool, "1") {
    return areEqual(5, 5)
}

func test_constraint_comparable_licenses_equality_string() test -> Expected(bool, "0") {
    return areEqual("abc", "xyz")
}

func test_constraint_comparable_licenses_equality_char() test -> Expected(bool, "1") {
    return areEqual('a', 'a')
}

// A constraint written after a name applies to that name and to every
// immediately preceding unconstrained name in the same list (A.7.1):
// [A, B: comparable] constrains both A and B — independently of each other.
func bothComparableEqual[A, B: comparable](a: A, b: A, c: B, d: B) -> bool {
    return a == b && c == d
}

func test_constraint_group_applies_to_both_names() test -> Expected(bool, "1") {
    return bothComparableEqual(1, 1, "x", "x")
}

func test_constraint_group_names_vary_independently() test -> Expected(bool, "0") {
    // A matches, B doesn't — proves the two constrained slots are checked
    // independently rather than collapsing into one shared type.
    return bothComparableEqual(1, 1, "x", "y")
}