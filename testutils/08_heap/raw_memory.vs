package heap_test
build test

// new[T] allocates count × sizeof(T) bytes and returns (typed_ptr T,
// string) (A.4.8) — the same boundary-tuple shape 06_composite_types/
// tuples.vs already exercised natively. Blocks are zeroed by default.

func test_new_zeroed_by_default() test -> Expected(int32, "0") {
    let p, err = new[int32](1)
    return &p
}

func test_new_success_error_string_empty() test -> Expected(bool, "1") {
    let p, err = new[int32](1)
    return err == ""
}

func test_new_success_returns_non_nil_pointer() test -> Expected(bool, "1") {
    let p, err = new[int32](1)
    return p != nil
}

// Unary & is dereference on a typed_ptr T (A.4.4); AssignTarget's
// `& UnaryExpression` form (A.5.2) writes through that dereference — the
// grammar's own example is `&p = 99`.
func test_new_pointer_write_then_read() test -> Expected(int32, "42") {
    let p, err = new[int32](1)
    &p = 42
    return &p
}

// T may be inferred from the declared type of the binding the result
// flows into (A.4.8), rather than written explicitly.
func test_new_type_argument_inferred_from_destination() test -> Expected(int32, "0") {
    let p: typed_ptr int32, err = new(1)
    return &p
}

// align must be a power of two; a violation is an allocation failure —
// nil and a non-empty string, not a distinct diagnostic (A.4.8).
func test_new_align_not_power_of_two_fails() test -> Expected(bool, "1") {
    let p, err = new[int32](1, align: 3)
    return err != ""
}

func test_new_align_violation_returns_nil_pointer() test -> Expected(bool, "1") {
    let p, err = new[int32](1, align: 3)
    return p == nil
}

func test_new_power_of_two_align_succeeds() test -> Expected(bool, "1") {
    let p, err = new[int32](1, align: 8)
    return err == ""
}

// delete has no return value and reading through a deleted pointer is
// undefined rather than diagnosable (A.14) — this only pins that the call
// itself compiles and runs without crashing.
func test_delete_runs_without_crash() test {
    let p, err = new[int32](1)
    delete(p)
}

// resize on success invalidates its input pointer (A.4.8) — reading
// through p after a successful resize is undefined (A.14), so this only
// checks the boundary tuple's result, never p itself afterward.
func test_resize_success_reports_empty_error() test -> Expected(bool, "1") {
    let p, err1 = new[int32](1)
    let p2, err2 = resize(p, 3)
    return err2 == ""
}

func test_resize_success_returns_non_nil_pointer() test -> Expected(bool, "1") {
    let p, err1 = new[int32](1)
    let p2, err2 = resize(p, 3)
    return p2 != nil
}

// A byte extent overflowing uint64 is an allocation failure (A.4.8); on
// failure the input pointer is left untouched and valid, so p may still
// be safely read afterward — unlike the success case above.
func test_resize_failure_leaves_original_pointer_valid() test -> Expected(int32, "7") {
    let p, err1 = new[int32](1)
    &p = 7
    let p2, err2 = resize(p, 18_446_744_073_709_551_615)
    return &p
}

func test_resize_failure_reports_nonempty_error() test -> Expected(bool, "1") {
    let p, err1 = new[int32](1)
    let p2, err2 = resize(p, 18_446_744_073_709_551_615)
    return err2 != ""
}

// addr is the sole way to take the address OF a typed_ptr binding, since
// & on a typed_ptr operand is already dereference (A.4.4); it requires an
// addressable operand — a var binding, never a let (A.4.8).
func test_addr_round_trips_through_typed_ptr_binding() test -> Expected(bool, "1") {
    let p, err = new[int32](1)
    var q: typed_ptr int32 = p
    let addrOfQ: typed_ptr (typed_ptr int32) = addr(q)
    return &addrOfQ == q
}

// copy is always overlap-safe (A.4.8); tested on non-overlapping
// single-element blocks. The count argument's unit (elements vs. bytes)
// is not stated in the excerpt — assumed to match new's element-count
// convention, since both apply to a new[T]-allocated block.
func test_copy_duplicates_value_between_blocks() test -> Expected(int32, "55") {
    let src, errSrc = new[int32](1)
    &src = 55
    let dest, errDest = new[int32](1)
    copy(dest, src, 1)
    return &dest
}

func test_copy_source_left_unchanged() test -> Expected(int32, "55") {
    let src, errSrc = new[int32](1)
    &src = 55
    let dest, errDest = new[int32](1)
    copy(dest, src, 1)
    return &src
}

// zero(dest, count) clears a block. Allocated with zero: false so the
// clearing effect is actually observable, and written before it is read,
// since reading unzeroed memory is undefined (A.14).
func test_zero_builtin_clears_previously_written_memory() test -> Expected(int32, "0") {
    let p, err = new[int32](1, zero: false)
    &p = 77
    zero(p, 1)
    return &p
}

// sizeof/alignof's own return type is never stated, so no test depends on
// it directly — each result is cast to int32 first. This only pins a
// relative property that must hold under any width-respecting
// implementation, regardless of the actual byte counts involved.
func test_sizeof_wider_type_is_larger() test -> Expected(bool, "1") {
    return (sizeof(int64) as int32) > (sizeof(int8) as int32)
}

func test_sizeof_same_type_is_equal() test -> Expected(bool, "1") {
    return (sizeof(int32) as int32) == (sizeof(int32) as int32)
}

func test_alignof_is_positive() test -> Expected(bool, "1") {
    return (alignof(int32) as int32) > 0
}

// reinterpret is a bit-cast between value types of identical byte size
// (A.4.4), never a value conversion. Tested as a round trip so the
// assertion never depends on knowing float32's exact bit layout — only
// that casting out and back is lossless.
func test_reinterpret_int_to_float_round_trips() test -> Expected(bool, "1") {
    let original: int32 = 1078530011
    let asFloat: float32 = reinterpret(float32, original)
    let roundTripped: int32 = reinterpret(int32, asFloat)
    return roundTripped == original
}

// The wrapping arithmetic operators (&+, &-, &*, A.4.5) only make sense
// under a two's-complement representation, which grounds this assumption:
// reinterpreting int32 -1's bit pattern (all-ones) as uint32 should read
// as the unsigned type's maximum value.
func test_reinterpret_negative_one_bit_pattern_as_unsigned() test -> Expected(uint32, "4294967295") {
    let x: int32 = -1
    return reinterpret(uint32, x)
}