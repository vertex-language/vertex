package literals_test
build test

func test_string_let() test -> Expected(string, "hello") {
    let s: string = "hello"
    return s
}

func test_string_var() test -> Expected(string, "world") {
    var s: string = "world"
    return s
}

func test_string_empty() test -> Expected(string, "") {
    let s: string = ""
    return s
}

func test_string_escape_newline() test -> Expected(string, "line1\nline2") {
    let s: string = "line1\nline2"
    return s
}

func test_string_escape_tab() test -> Expected(string, "a\tb") {
    let s: string = "a\tb"
    return s
}

func test_string_escape_backslash() test -> Expected(string, "a\\b") {
    let s: string = "a\\b"
    return s
}

func test_string_escape_double_quote() test -> Expected(string, "say \"hi\"") {
    let s: string = "say \"hi\""
    return s
}

func test_string_hex_escape() test -> Expected(string, "AB") {
    let s: string = "\x41\x42"
    return s
}

func test_string_unicode_escape() test -> Expected(string, "Hi") {
    let s: string = "\u{48}\u{69}"
    return s
}

func test_string_unicode_direct() test -> Expected(string, "日本語") {
    let s: string = "日本語"
    return s
}

func test_string_multiline() test -> Expected(string, "hi") {
    let s: string = `hi`
    return s
}

func test_string_raw_no_escapes() test -> Expected(string, "a\\nb") {
    // A backtick string recognises no escape sequence: this is a literal
    // backslash followed by 'n', not a newline (A.1.5.2)
    let s: string = `a\nb`
    return s
}

func test_string_raw_embedded_quotes() test -> Expected(string, "she said \"hi\"") {
    // Double quotes need no escaping inside a raw string
    let s: string = `she said "hi"`
    return s
}

func test_string_raw_spans_lines() test -> Expected(string, "line1\nline2") {
    // Every LineTerminator a backtick string spans is part of its value
    let s: string = `line1
line2`
    return s
}