package heap_test
build test

// unique(...) and shared(...) are "the only two heap doors" (A.4.8): each
// constructs a wrapper by moving its operand in unconditionally. Their
// cost distinction — unique's bare copy is a deep O(data) walk, shared's
// is a cheap atomic increment (A.9.4) — belongs to 09_ownership/, which
// owns the copy-vs-transfer story generally. This file only pins that the
// constructors, and weak/upgrade alongside them, work mechanically.
//
// Reading a field back out through a constructed unique/shared handle
// would need a dot-access-through-handle convention the excerpt never
// shows (unary & is documented as dereference specifically for
// typed_ptr, A.4.4 — not for unique/shared/weak), so field-level content
// checks are left for whichever folder can ground that mechanism, rather
// than guessed at here.

struct Payload {
    value: int32
}

func test_unique_construction_runs_without_crash() test {
    let u: unique Payload = unique(Payload{value: 9})
}

func test_shared_construction_runs_without_crash() test {
    let s: shared Payload = shared(Payload{value: 9})
}

// weak T observes only a shared allocation (A.3.2); upgrade(w) returns
// (shared T, string) — the boundary-tuple convention applied to a race
// the type system cannot statically win (A.4.8). While the original
// shared handle is still alive, the upgrade must succeed.
func test_weak_upgrade_succeeds_while_shared_handle_alive() test -> Expected(bool, "1") {
    let s: shared Payload = shared(Payload{value: 1})
    let w: weak Payload = weak(s)
    let s2, err = upgrade(w)
    return err == ""
}

// Once every shared handle to the allocation has gone out of scope — its
// teardown emitted at block-scope end, same as any binding (A.6.4) — the
// race is lost: upgrade must report failure.
func test_weak_upgrade_fails_after_shared_handle_deallocated() test -> Expected(bool, "1") {
    var w: weak Payload
    {
        let s: shared Payload = shared(Payload{value: 1})
        w = weak(s)
    }
    let s2, err = upgrade(w)
    return err != ""
}

// drop(x) explicitly ends a transferred binding's lifetime without
// emitting its teardown (A.4.8). Observing "teardown was skipped" needs
// deinit (06_composite_types/) combined with the var-transfer convention
// (09_ownership/), so it's named here but not tested — this folder only
// establishes that construction and weak/upgrade resolve correctly.