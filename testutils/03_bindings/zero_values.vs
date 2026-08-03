package bindings_test
build test

// `var Binding` with no initializer yields the type's zero value and
// requires an explicit TypeAnnotation (A.5.1). Note the grammar only
// permits a single Binding in this form, never a BindingList — one
// statement per zero-valued declaration.

func test_zero_value_int8() test -> Expected(int8, "0") {
    var x: int8
    return x
}

func test_zero_value_int16() test -> Expected(int16, "0") {
    var x: int16
    return x
}

func test_zero_value_int32() test -> Expected(int32, "0") {
    var x: int32
    return x
}

func test_zero_value_int64() test -> Expected(int64, "0") {
    var x: int64
    return x
}

func test_zero_value_int() test -> Expected(int, "0") {
    var x: int
    return x
}

func test_zero_value_uint8() test -> Expected(uint8, "0") {
    var x: uint8
    return x
}

func test_zero_value_uint16() test -> Expected(uint16, "0") {
    var x: uint16
    return x
}

func test_zero_value_uint32() test -> Expected(uint32, "0") {
    var x: uint32
    return x
}

func test_zero_value_uint64() test -> Expected(uint64, "0") {
    var x: uint64
    return x
}

func test_zero_value_uint() test -> Expected(uint, "0") {
    var x: uint
    return x
}

func test_zero_value_byte() test -> Expected(byte, "0") {
    // byte and uint8 denote the same type (A.1.4)
    var x: byte
    return x
}

func test_zero_value_float32() test -> Expected(float32, "0.000000") {
    var x: float32
    return x
}

func test_zero_value_float64() test -> Expected(float64, "0.000000") {
    var x: float64
    return x
}

func test_zero_value_bool() test -> Expected(bool, "0") {
    var x: bool
    return x
}

func test_zero_value_string() test -> Expected(string, "") {
    var x: string
    return x
}

func test_zero_value_char() test -> Expected(bool, "1") {
    // A zero-valued char has no agreed-on printable rendering, so this pins
    // the rule via comparison instead: the zero value is the NUL scalar.
    var c: char
    return c == '\0'
}