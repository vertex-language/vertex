package control_flow_test
build test

// EnumPattern (A.5.7) and exhaustiveness-on-a-unit-enum are deferred to
// 06_composite_types/, where `enum` is declared. This file exercises
// switch over the scalar and string types already available.

func test_switch_matches_case() test -> Expected(int32, "1") {
    var x: int32 = 2
    var result: int32 = 0
    switch x {
        case 1:
            result = 1
        case 2:
            result = 2
        default:
            result = 99
    }
    return result
}

func test_switch_falls_to_default() test -> Expected(int32, "99") {
    var x: int32 = 42
    var result: int32 = 0
    switch x {
        case 1:
            result = 1
        case 2:
            result = 2
        default:
            result = 99
    }
    return result
}

func test_switch_no_default_no_match_executes_nothing() test -> Expected(int32, "0") {
    var x: int32 = 42
    var result: int32 = 0
    switch x {
        case 1:
            result = 1
        case 2:
            result = 2
    }
    return result
}

// A single CaseClause may list several patterns (PatternList, A.5.7); any
// one of them matches.
func test_switch_multiple_patterns_matches_first() test -> Expected(int32, "1") {
    var x: int32 = 1
    var result: int32 = 0
    switch x {
        case 1, 2, 3:
            result = 1
        default:
            result = 0
    }
    return result
}

func test_switch_multiple_patterns_matches_later() test -> Expected(int32, "1") {
    var x: int32 = 3
    var result: int32 = 0
    switch x {
        case 1, 2, 3:
            result = 1
        default:
            result = 0
    }
    return result
}

// Cases do not fall through implicitly (A.5.7).
func test_switch_no_implicit_fallthrough() test -> Expected(int32, "1") {
    var x: int32 = 1
    var result: int32 = 0
    switch x {
        case 1:
            result = 1
        case 2:
            result = 2
    }
    return result
}

// fallthrough is explicit and must be the last statement in its clause.
func test_switch_explicit_fallthrough() test -> Expected(int32, "12") {
    var x: int32 = 1
    var result: int32 = 0
    switch x {
        case 1:
            result = 1
            fallthrough
        case 2:
            result = result * 10 + 2
        default:
            result = 99
    }
    return result
}

func test_switch_fallthrough_chains_through_default() test -> Expected(int32, "129") {
    var x: int32 = 1
    var result: int32 = 0
    switch x {
        case 1:
            result = 1
            fallthrough
        case 2:
            result = result * 10 + 2
            fallthrough
        default:
            result = result * 10 + 9
    }
    return result
}

func test_switch_over_char() test -> Expected(int32, "1") {
    var c: char = 'b'
    var result: int32 = 0
    switch c {
        case 'a':
            result = 0
        case 'b':
            result = 1
        default:
            result = 99
    }
    return result
}

func test_switch_over_string() test -> Expected(int32, "1") {
    var s: string = "beta"
    var result: int32 = 0
    switch s {
        case "alpha":
            result = 0
        case "beta":
            result = 1
        default:
            result = 99
    }
    return result
}

func test_switch_over_bool() test -> Expected(int32, "1") {
    var b: bool = true
    var result: int32 = 0
    switch b {
        case true:
            result = 1
        case false:
            result = 0
    }
    return result
}

func test_switch_pattern_is_general_expression() test -> Expected(int32, "1") {
    // Pattern is Expression[~Lit] (A.5.7), not restricted to bare literals.
    var one: int32 = 1
    var x: int32 = 1
    var result: int32 = 0
    switch x {
        case one:
            result = 1
        default:
            result = 0
    }
    return result
}

func test_switch_nested() test -> Expected(int32, "1") {
    var x: int32 = 1
    var y: int32 = 2
    var result: int32 = 0
    switch x {
        case 1:
            switch y {
                case 2:
                    result = 1
                default:
                    result = 0
            }
        default:
            result = 0
    }
    return result
}