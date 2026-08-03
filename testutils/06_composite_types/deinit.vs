package composite_types_test
build test

// Observing deinit needs a way for a dying binding to report out, and
// deinit's own signature takes no parameters (A.6.4). These tests use a
// top-level `var` as a package-level mutable log. A.2 only constrains a
// top-level declaration's *initializer* to be compile-time-evaluable — it
// doesn't say the binding can't be reassigned later from a function body —
// but this exact usage isn't shown verbatim in the excerpt, so treat it
// as inference.
var deinitLog: int32 = 0

class Tagged {
    tag: int32
}

func (t: Tagged) init(tag: int32) {
    t.tag = tag
}

// deinit is emitted where the binding's liveness ends (A.6.4). A bare
// `{ }` block is itself a Statement (A.5), giving a lexical scope shorter
// than the enclosing function.
func (t: Tagged) deinit() {
    deinitLog = deinitLog * 10 + t.tag
}

func test_deinit_fires_at_block_scope_end() test -> Expected(int32, "7") {
    deinitLog = 0
    {
        let a = Tagged(tag: 7)
    }
    return deinitLog
}

// Locals torn down in reverse declaration order (A.6.4): b was declared
// after a, so b's deinit fires first.
func test_deinit_locals_fire_in_reverse_declaration_order() test -> Expected(int32, "21") {
    deinitLog = 0
    {
        let a = Tagged(tag: 1)
        let b = Tagged(tag: 2)
    }
    return deinitLog
}

func makeAndDropTagged(tag: int32) {
    let t = Tagged(tag: tag)
}

func test_deinit_fires_on_function_fallthrough_exit() test -> Expected(int32, "9") {
    deinitLog = 0
    makeAndDropTagged(9)
    return deinitLog
}

// Field-level teardown order ("fields in reverse declaration order",
// A.6.4) would need a class whose fields are themselves owning class
// instances — a pattern 09_ownership/ is responsible for establishing
// first. Left untested here for that reason.