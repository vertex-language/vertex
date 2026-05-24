package main
build linux
import "linux/syscalls"

class Syscalls : syscalls {
    func write(fd: int, buf: any void, count: uint) -> int
    func exit(status: int) -> int
}

// ── helpers ───────────────────────────────────────────────────────────────────
func sys_print(s: Syscalls, msg: string, len: int) {
    s.write(1, msg.any(), uint(len))
}

// ── Class definition ──────────────────────────────────────────────────────────
class YourClass {
    var x: int32
    var y: int32
}

// ── Associated functions ──────────────────────────────────────────────────────
// Receiver is mut so the function can write back through it (§26)
func setX(self: mut YourClass, x: int32) {
    self.x = x
}

func setY(self: mut YourClass, y: int32) {
    self.y = y
}

func getX(self: YourClass) -> int32 {
    return self.x
}

func getY(self: YourClass) -> int32 {
    return self.y
}

// ── entry point ───────────────────────────────────────────────────────────────
func main() -> int {
    var s = Syscalls()

    // .new() opts the instance into ref counting — no manual .delete() needed
    var m = YourClass(x: 0, y: 0).new()

    // mut receiver requires a var binding prefixed with & at the call site (§26)
    m.setX(x: int32(10))
    m.setY(y: int32(20))

    s.sys_print(msg: "x = 10\n", len: 7)
    s.sys_print(msg: "y = 20\n", len: 7)

    return 0
}