package bindings_test
build test

// Ordinary identifier forms (A.1.2)
func test_identifier_leading_underscore() test -> Expected(int32, "11") {
    var _foo: int32 = 11
    return _foo
}

func test_identifier_underscore_middle() test -> Expected(int32, "12") {
    var foo_bar: int32 = 12
    return foo_bar
}

func test_identifier_trailing_digits() test -> Expected(int32, "13") {
    var foo123: int32 = 13
    return foo123
}

func test_identifier_unicode_id_start() test -> Expected(int32, "3") {
    // UnicodeIDStart admits any code point with Unicode property ID_Start,
    // not just ASCII letters.
    var π: int32 = 3
    return π
}

// BlankIdentifier as a discarded destructuring target (A.1.2)
func test_blank_identifier_discards_destructure_target() test -> Expected(int32, "4") {
    let _, b = (3, 4)
    return b
}

// Contextual keywords (A.1.3): each is an Identifier everywhere except its
// one special position, so a local binding of the same name is unambiguous.
func test_contextual_keyword_module_as_identifier() test -> Expected(int32, "1") {
    var module: int32 = 1
    return module
}

func test_contextual_keyword_framework_as_identifier() test -> Expected(int32, "2") {
    var framework: int32 = 2
    return framework
}

func test_contextual_keyword_init_as_identifier() test -> Expected(int32, "3") {
    // `init` is a modifier only in a receiver declaration or a declare
    // block (A.6.4, A.8.3); elsewhere it is an ordinary identifier.
    var init: int32 = 3
    return init
}

func test_contextual_keyword_deinit_as_identifier() test -> Expected(int32, "4") {
    var deinit: int32 = 4
    return deinit
}

func test_contextual_keyword_error_as_identifier() test -> Expected(int32, "5") {
    // `error` is meaningful only inside Expected(...) (A.12.2).
    var error: int32 = 5
    return error
}

func test_contextual_keyword_build_as_identifier() test -> Expected(int32, "6") {
    // `build` is special only as the second line-initial token of a source
    // file (A.2.2).
    var build: int32 = 6
    return build
}

func test_contextual_keyword_test_as_identifier() test -> Expected(int32, "7") {
    // `test` is special only in FunctionMarker position (A.6.1); a local
    // variable named `test` inside a test-marked function is unambiguous.
    var test: int32 = 7
    return test
}