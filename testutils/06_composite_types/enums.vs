package composite_types_test
build test

// An explicit discriminant is used on the first variant below rather than
// relying on an implicit starting value: A.6.5 states unassigned variants
// "continue from the previous value" but never states what an entirely
// unassigned enum's first variant defaults to, so that case is left
// untested (see testutils/README.md on spec-silent behavior).
enum Status: int32 {
    Active = 1,
    Inactive,
    Suspended,
}

// A unit-only enum is its discriminant integer; `as int32` is a
// reinterpretation of that integer, not a conversion (A.6.5).
func statusCode(s: Status) -> int32 {
    return s as int32
}

func test_enum_explicit_discriminant() test -> Expected(int32, "1") {
    return statusCode(Status.Active)
}

func test_enum_discriminant_continues_from_previous() test -> Expected(int32, "2") {
    return statusCode(Status.Inactive)
}

func test_enum_discriminant_continues_across_multiple_unassigned() test -> Expected(int32, "3") {
    return statusCode(Status.Suspended)
}

// EnumShorthand (.Identifier) is legal wherever the enum type is
// inferable from context (A.4.1): a typed binding, an argument position,
// a return, or a case label — each exercised below.

func test_enum_shorthand_in_typed_binding() test -> Expected(int32, "1") {
    var s: Status = .Active
    return statusCode(s)
}

func test_enum_shorthand_in_argument_position() test -> Expected(int32, "2") {
    return statusCode(.Inactive)
}

func pickStatus(useActive: bool) -> Status {
    if useActive {
        return .Active
    }
    return .Suspended
}

func test_enum_shorthand_in_return_position_true_branch() test -> Expected(int32, "1") {
    return statusCode(pickStatus(true))
}

func test_enum_shorthand_in_return_position_false_branch() test -> Expected(int32, "3") {
    return statusCode(pickStatus(false))
}

// A switch over a unit-only enum with no default must be exhaustive
// (A.5.7); all three variants are covered here.
func describe(s: Status) -> int32 {
    switch s {
        case .Active:
            return 10
        case .Inactive:
            return 20
        case .Suspended:
            return 30
    }
    return -1
}

func test_switch_enum_pattern_active() test -> Expected(int32, "10") {
    return describe(Status.Active)
}

func test_switch_enum_pattern_inactive() test -> Expected(int32, "20") {
    return describe(Status.Inactive)
}

func test_switch_enum_pattern_suspended() test -> Expected(int32, "30") {
    return describe(Status.Suspended)
}

// Payload variants (A.6.5): a payload enum is a tagged union. Payload
// bindings in an EnumPattern are matched positionally, with arity equal
// to the variant's declared arity, and `_` discards a position (A.5.7).
enum Shape {
    Circle(float32),
    Rectangle(float32, float32),
    Point,
}

func area(s: Shape) -> float32 {
    switch s {
        case .Circle(r):
            return r * r * 3.0
        case .Rectangle(w, h):
            return w * h
        case .Point:
            return 0.0
    }
    return -1.0
}

func test_enum_payload_single_value_binding() test -> Expected(float32, "12.000000") {
    return area(.Circle(2.0))
}

func test_enum_payload_two_value_binding() test -> Expected(float32, "6.000000") {
    return area(.Rectangle(2.0, 3.0))
}

func test_enum_payload_unit_variant_among_payload_variants() test -> Expected(float32, "0.000000") {
    return area(.Point)
}

func hasPayload(s: Shape) -> bool {
    switch s {
        case .Circle(_):
            return true
        case .Rectangle(_, _):
            return true
        case .Point:
            return false
    }
    return false
}

func test_enum_payload_binding_discarded_single() test -> Expected(bool, "1") {
    return hasPayload(.Circle(1.0))
}

func test_enum_payload_binding_discarded_pair() test -> Expected(bool, "1") {
    return hasPayload(.Rectangle(1.0, 2.0))
}

func test_enum_payload_unit_variant_reports_false() test -> Expected(bool, "0") {
    return hasPayload(.Point)
}

// Whether an EnumPattern payload binding is a mutable view or an
// observably-distinct copy (A.5.7 calls it "a view into the payload, not
// a copy") isn't tested here: proving that distinction needs a way to
// write through the binding, which isn't shown for switch-case payload
// bindings anywhere in the excerpt.