package main
build linux
import "lib/sdl2"

// Opaque C pointers map to mut any opaque (mutable struct*).
// SDL functions with no return value are declared without a return type.
class Sdl2 : sdl2 {
    func SDL_Init(flags: uint32) -> int
    func SDL_Quit()

    func SDL_CreateWindow(title: any char, x: int, y: int, w: int, h: int, flags: uint32) -> mut any opaque?
    func SDL_DestroyWindow(window: mut any opaque)

    func SDL_CreateRenderer(window: mut any opaque, index: int, flags: uint32) -> mut any opaque?
    func SDL_DestroyRenderer(renderer: mut any opaque)

    func SDL_SetRenderDrawColor(renderer: mut any opaque, r: uint8, g: uint8, b: uint8, a: uint8) -> int
    func SDL_RenderClear(renderer: mut any opaque) -> int
    func SDL_RenderPresent(renderer: mut any opaque)

    func SDL_Delay(ms: uint32)

    func puts(str: any char) -> int
}

func main() -> int {
    var sdl = Sdl2()

    // 1. Initialize SDL (SDL_INIT_VIDEO = 0x20)
    if sdl.SDL_Init(uint32(0x20)) != 0 {
        sdl.puts("Failed to initialize SDL2".any())
        return 1
    }
    defer sdl.SDL_Quit()

    // 2. Create window (SDL_WINDOWPOS_CENTERED = 0x2FFF0000, SHOWN = 0x04)
    let window = sdl.SDL_CreateWindow(
        "Vertex SDL2 Window".any(),
        0x2FFF0000, 0x2FFF0000,
        800, 600,
        uint32(0x04)
    )

    // 3. Unwrap the optional window pointer
    if let win = window {
        defer sdl.SDL_DestroyWindow(win)

        // 4. Create hardware-accelerated renderer (ACCELERATED = 2)
        let renderer = sdl.SDL_CreateRenderer(win, -1, uint32(2))

        if let ren = renderer {
            defer sdl.SDL_DestroyRenderer(ren)

            sdl.puts("Window created! Drawing red screen...".any())

            // 5. Single-frame render — explicit uint8 casts per numeric rules
            sdl.SDL_SetRenderDrawColor(ren, uint8(255), uint8(50), uint8(50), uint8(255))
            sdl.SDL_RenderClear(ren)
            sdl.SDL_RenderPresent(ren)

            sdl.SDL_Delay(uint32(3000))

            sdl.puts("Shutting down cleanly...".any())
            return 0
        }
    }

    sdl.puts("Failed to create window or renderer".any())
    return 1
}