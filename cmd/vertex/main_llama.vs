package main
import "lib/llama"

class Llama : llama {
    func llama_backend_init()
    func llama_backend_free()

    func llama_load_model_from_file(path: any char, params: any void) -> mut any opaque?
    func llama_free_model(model: mut any opaque)

    func llama_new_context_with_model(model: mut any opaque, params: any void) -> mut any opaque?
    func llama_free(ctx: mut any opaque)

    func llama_generate_to_stdout(ctx: mut any opaque, prompt: any char, max_tokens: int) -> int

    func puts(str: any char) -> int
}

// ── Engine class ──────────────────────────────────────────────────────────────
// Heap-allocated class to manage the two opaque C pointers
class LlamaEngine {
    var model: mut any opaque
    var ctx:   mut any opaque
}

// deinit runs automatically when the ref count hits zero (§28.1).
func deinit(engine: mut LlamaEngine) {
    var ll = Llama()
    ll.puts("Vertex deinit: Freeing llama.cpp C-pointers...".any())
    ll.llama_free(engine.ctx)
    ll.llama_free_model(engine.model)
}

// Factory — returns Result so the caller must handle both Ok and Err
func load_engine(path: string) -> Result(LlamaEngine, string) {
    var ll = Llama()
    var dummy_params = 0

    let model_ptr = ll.llama_load_model_from_file(path.any(), dummy_params.any())

    if let m = model_ptr {
        let ctx_ptr = ll.llama_new_context_with_model(m, dummy_params.any())

        if let c = ctx_ptr {
            return Result(Ok, LlamaEngine(model: m, ctx: c).new())
        }

        ll.llama_free_model(m)
        return Result(Err, "Failed to create context")
    }

    return Result(Err, "Failed to load model")
}

// Associated function — first parameter is the LlamaEngine receiver
func generate(engine: LlamaEngine, prompt: string) {
    var ll = Llama()
    ll.llama_generate_to_stdout(engine.ctx, prompt.any(), 128)
}

// ── entry point ───────────────────────────────────────────────────────────────
func main() -> int {
    var ll = Llama()

    ll.puts("Booting AI engine...".any())

    ll.llama_backend_init()
    defer ll.llama_backend_free()

    let path = "./models/llama-3-8b.Q4_K_M.gguf"

    // Exhaustive switch over Result — both arms are required (§37.3)
    switch load_engine(path: path) {

    case Ok(let engine):
        ll.puts("Model loaded successfully!".any())
        engine.generate(prompt: "The future of language models is")
        // 'engine' scope ends here — ref count hits 0 — deinit fires automatically

    case Err(let err_msg):
        ll.puts("Fatal Error:".any())
        ll.puts(err_msg.any())
        return 1
    }

    return 0
}