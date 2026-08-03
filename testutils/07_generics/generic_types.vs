package generics_test
build test

// A method receiver re-declares the type's parameter list to bring the
// name into scope: the receiver's [T] binds Box's own T, it does not
// introduce a fresh one (A.7.6). A constraint declared on the type is in
// force inside every method of that type without being restated.
struct Box[T: comparable] {
    value: T
}

func (b: Box[T]) equals(other: T) -> bool {
    return b.value == other
}

func test_generic_type_method_constraint_in_force() test -> Expected(bool, "1") {
    let box: Box[int32] = Box[int32]{value: 5}
    return box.equals(5)
}

func test_generic_type_method_different_instantiation() test -> Expected(bool, "0") {
    let box: Box[string] = Box[string]{value: "hello"}
    return box.equals("world")
}

// Classes follow the identical receiver-declaration model (A.6.3, A.7.6):
// a class differs from a struct in construction and identity, not in how
// its methods relate to the type's parameter list.
class Wrapper[T: comparable] {
    payload: T
}

func (w: Wrapper[T]) init(payload: T) {
    w.payload = payload
}

func (w: Wrapper[T]) matches(other: T) -> bool {
    return w.payload == other
}

func test_generic_class_method_constraint_in_force() test -> Expected(bool, "1") {
    let w: Wrapper[int32] = Wrapper[int32](payload: 5)
    return w.matches(5)
}

func test_generic_class_method_different_instantiation() test -> Expected(bool, "0") {
    let w: Wrapper[string] = Wrapper[string](payload: "x")
    return w.matches("y")
}

// A type argument may itself be an instantiated generic type
// (TypeArgumentList admits any Type, A.3.6), so instantiation nests as
// long as it terminates (A.7.6 forbids only the unbounded case, which is
// rejection-territory). Kept unconstrained (any): struct comparability is
// never established anywhere in this suite, so Box[T: comparable] above
// is deliberately not reused here.
struct Wrap[T] {
    value: T
}

func test_generic_nested_instantiation() test -> Expected(int32, "4") {
    let inner: Wrap[int32] = Wrap[int32]{value: 4}
    let outer: Wrap[Wrap[int32]] = Wrap[Wrap[int32]]{value: inner}
    return outer.value.value
}