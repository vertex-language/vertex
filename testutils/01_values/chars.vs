package literals_test
build test

func test_char_literal() test -> Expected(char, "A") {
    let c: char = 'A'
    return c
}

func test_char_lowercase() test -> Expected(char, "z") {
    let c: char = 'z'
    return c
}

func test_char_digit() test -> Expected(char, "0") {
    let c: char = '0'
    return c
}

func test_char_unicode_scalar() test -> Expected(char, "日") {
    // A CharLiteral denotes one Unicode scalar value, not one byte (A.1.5.2)
    let c: char = '日'
    return c
}

func test_char_escape_newline() test -> Expected(char, "\n") {
    let c: char = '\n'
    return c
}

func test_char_escape_tab() test -> Expected(char, "\t") {
    let c: char = '\t'
    return c
}

func test_char_escape_carriage_return() test -> Expected(char, "\r") {
    let c: char = '\r'
    return c
}

func test_char_escape_backslash() test -> Expected(char, "\\") {
    let c: char = '\\'
    return c
}

func test_char_escape_single_quote() test -> Expected(char, "'") {
    let c: char = '\''
    return c
}

func test_char_hex_escape() test -> Expected(char, "A") {
    // \x41 is ASCII 0x41 = 'A'
    let c: char = '\x41'
    return c
}

func test_char_unicode_escape() test -> Expected(char, "H") {
    let c: char = '\u{48}'
    return c
}

func test_char_unicode_escape_non_ascii() test -> Expected(char, "é") {
    let c: char = '\u{E9}'
    return c
}