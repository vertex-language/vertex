package main
import "linux/libc"

// Native class: zero-size at runtime. The compiler turns every method call
// into a direct linked C call and removes all instances from the output.
class C : libc {
    func printf(fmt: string, ...) -> int32
    func puts(s: string) -> int32
    func exit(code: int32)
}

// greet takes a Vertex string — which is already const char* at the C level —
// and passes it straight through to printf. No conversion needed.
func greet(name: string) {
    var c = C()
    c.printf("Hello, %s!\n", name)
}

func main() -> int {
    greet("Vertex")

    var c = C()
    var x: int32 = 100

    if x > 50 {
        x = x + 10
        if x > 100 {
            x = x + 10
        } else {
            x = x - 10
        }
    } else {
        x = x - 10
    }

    c.printf("x = %d\n", x)
    return 0
}