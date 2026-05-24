package main
import "lib/c"

// fread writes into its buffer, so ptr is mut any void.
// fwrite only reads from its buffer, so ptr is any void (const).
class C : c {
    func fopen(path: any char, mode: any char) -> mut any opaque?
    func fwrite(ptr: any void, size: int, count: int, stream: mut any opaque) -> int
    func fread(ptr: mut any void, size: int, count: int, stream: mut any opaque) -> int
    func fclose(stream: mut any opaque) -> int
    func printf(fmt: any char, ...) -> int
}

// ── File handler class ────────────────────────────────────────────────────────
class FileHandler {
    var path: string
}

// Write `content` to the file, repeated 100 times (overwrites on open)
func write_file(f: FileHandler, content: string, len: int) -> int {
    var c = C()
    let handle_opt = c.fopen(f.path.any(), "w".any())
    if let handle = handle_opt {
        defer c.fclose(handle)
        for i in 0...99 {
            c.fwrite(content.any(), 1, len, handle)
        }
    }
    return 0
}

// Read up to `len` bytes from the file into `buf`
func read_file(f: FileHandler, buf: string, len: int) -> int {
    var c = C()
    let handle_opt = c.fopen(f.path.any(), "r".any())
    if let handle = handle_opt {
        defer c.fclose(handle)
        let bytes_read = c.fread(buf.any(), 1, len, handle)
        return bytes_read
    }
    return 0
}

// ── entry point ───────────────────────────────────────────────────────────────
func main() -> int {
    var c = C()
    let f = FileHandler(path: "hello.txt").new()

    // ── Write ─────────────────────────────────────────────────────────────────
    var msg = "Hello from FileHandler!\n"
    f.write_file(content: msg, len: 24)

    // ── Read back ─────────────────────────────────────────────────────────────
    var buf = "                        "   // 24-char blank buffer
    f.read_file(buf: buf, len: 24)
    c.printf("Read back: %s".any(), buf.any())

    return 0
}