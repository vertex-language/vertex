package composite_types_test
build test

// === and !== compare storage identity and are legal on classes only
// (A.4.5) — "same allocation?", never "same bytes?". Comparing a binding
// to itself needs no parameter-passing indirection (the exclusivity rule
// in A.9.3 governs argument passing, not referencing a binding twice
// inside an ordinary expression), so it's tested directly here.

class Box {
    value: int32
}

func (b: Box) init(value: int32) {
    b.value = value
}

func test_identity_binding_compared_to_itself() test -> Expected(bool, "1") {
    let a = Box(value: 5)
    return a === a
}

func test_identity_different_instances_not_identical() test -> Expected(bool, "1") {
    // Two separately constructed instances, even with identical field
    // values, are different storage.
    let a = Box(value: 5)
    let b = Box(value: 5)
    return a !== b
}

func test_identity_operator_false_for_different_instances() test -> Expected(bool, "0") {
    let a = Box(value: 5)
    let b = Box(value: 5)
    return a === b
}

func test_identity_negation_true_for_same_binding() test -> Expected(bool, "0") {
    let a = Box(value: 5)
    return a !== a
}

// == on a class is not exercised here: nothing in the excerpt states that
// classes derive a default byte-equality operator (A.7.4 only says
// `comparable` admits "every type supporting ==/!="), so whether == is
// even legal on Box is left unasserted rather than guessed.