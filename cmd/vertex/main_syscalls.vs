package main
build linux
import "linux/kernel/syscalls"

class Syscalls : syscalls {
    func write(fd: int, buf: any void, count: uint) -> int
    func exit(status: int) -> int
}

// ── helpers ───────────────────────────────────────────────────────────────────
func sys_print(s: Syscalls, msg: string, len: int) {
    s.write(1, msg.any(), uint(len))
}

// ── entry point ───────────────────────────────────────────────────────────────
func main() -> int {
    var s = Syscalls()

    s.sys_print(msg: "=== fibonacci ===\n", len: 18)
    s.sys_print(msg: "=== factorial ===\n", len: 18)
    s.sys_print(msg: "=== power ===\n",     len: 14)
    s.sys_print(msg: "=== gcd ===\n",       len: 12)
    s.sys_print(msg: "=== abs ===\n",       len: 12)

    return 0
}