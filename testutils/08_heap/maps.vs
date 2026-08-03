package heap_test
build test

// map[K]V requires K to satisfy comparable (A.3.1); string and int32 keys
// both support ==/!= and so both qualify (A.7.4).

func test_map_literal_read_by_key() test -> Expected(int32, "1") {
    var m: map[string]int32 = {"a": 1, "b": 2}
    return m["a"]
}

func test_map_literal_read_second_key() test -> Expected(int32, "2") {
    var m: map[string]int32 = {"a": 1, "b": 2}
    return m["b"]
}

func test_map_subscript_write() test -> Expected(int32, "99") {
    var m: map[string]int32 = {"a": 1}
    m["a"] = 99
    return m["a"]
}

func test_map_subscript_write_new_key() test -> Expected(int32, "42") {
    var m: map[string]int32 = {"a": 1}
    m["z"] = 42
    return m["z"]
}

func test_map_int32_keys() test -> Expected(int32, "7") {
    var m: map[int32]int32 = {1: 7, 2: 8}
    return m[1]
}

// Map iteration order is unspecified (A.5.6); only order-independent
// properties (a sum, a count) are asserted.
func test_map_for_in_key_value_sums_values() test -> Expected(int32, "6") {
    var m: map[string]int32 = {"a": 1, "b": 2, "c": 3}
    var total: int32 = 0
    for _, v in m {
        total += v
    }
    return total
}

func test_map_for_in_counts_entries() test -> Expected(int32, "3") {
    var m: map[string]int32 = {"a": 1, "b": 2, "c": 3}
    var count: int32 = 0
    for _, _ in m {
        count += 1
    }
    return count
}

// Assigning nil into a map subscript is the erase operation (A.5.2). What
// a subsequent read of the erased key itself returns is not stated
// anywhere in the excerpt, so that read is never attempted here; erase is
// instead verified indirectly, through the entry count dropping.
func test_map_erase_via_nil_reduces_entry_count() test -> Expected(int32, "2") {
    var m: map[string]int32 = {"a": 1, "b": 2, "c": 3}
    m["a"] = nil
    var count: int32 = 0
    for _, _ in m {
        count += 1
    }
    return count
}

func test_map_erase_does_not_affect_other_keys() test -> Expected(int32, "2") {
    var m: map[string]int32 = {"a": 1, "b": 2}
    m["a"] = nil
    return m["b"]
}