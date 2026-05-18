// main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	x86    "github.com/vertex-language/wasm-compiler/compiler/x86_64"
	"github.com/vertex-language/wasm-compiler/encoder"
	"github.com/vertex-language/wasm-compiler/linker"
	"github.com/vertex-language/wasm-compiler/object"
	"github.com/vertex-language/vertex/compiler"
)

func main() {
	outPath  := flag.String("o", "", "output file (default: derived from input)")
	emitWasm := flag.Bool("wasm", false, "emit .wasm binary instead of native ELF")
	verbose  := flag.Bool("v", false, "verbose compilation output")
	entry    := flag.String("entry", "main", "ELF entry-point export symbol (native only)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr,
			"usage: vertex [-o output] [-wasm] [-v] [-entry sym] <source.vs>")
		os.Exit(1)
	}

	if *outPath == "" {
		base := strings.TrimSuffix(filepath.Base(args[0]), filepath.Ext(args[0]))
		if *emitWasm {
			*outPath = base + ".wasm"
		} else {
			*outPath = base
		}
	}

	src, err := os.ReadFile(args[0])
	if err != nil {
		fatalf("cannot read %s: %v", args[0], err)
	}

	comp, err := compiler.NewCompiler(compiler.Options{Verbose: *verbose})
	if err != nil {
		fatalf("compiler init: %v", err)
	}

	mod, err := comp.CompileToModule(string(src), args[0])
	if err != nil {
		fatalf("%v", err)
	}

	var (
		out  []byte
		perm = os.FileMode(0o644)
	)

	if *emitWasm {
		out, err = encoder.Encode(mod)
		if err != nil {
			fatalf("wasm encode: %v", err)
		}
	} else {
		// pointerArgs tells the x86_64 backend which i32 parameters of each
		// imported function are wasm-relative memory offsets that must be
		// translated to native pointers by adding R15 before the call.
		// Key = real symbol name (after __vararg stripping).
		// Value = one bool per fixed parameter; true means "add R15".
		pointerArgs := map[string][]bool{
			"printf":  {true},         // char *format
			"puts":    {true},         // char *s
			"putchar": {false},        // int c  — not a pointer
			"fputs":   {true, false},  // char *s, FILE *stream
			"fprintf": {false, true},  // FILE *stream, char *format
			"strlen":  {true},         // char *s
			"strcpy":  {true, true},   // char *dst, char *src
			"strcmp":  {true, true},   // char *s1, char *s2
			"malloc":  {false},        // size_t size — not a pointer in
			"free":    {true},         // void *ptr
		}

		obj, err := x86.Compile(mod, false, pointerArgs)
		if err != nil {
			fatalf("x86_64 compile: %v", err)
		}
		out, err = linker.Link([]*object.WasmObj{obj}, linker.Options{
			Output: linker.ELF,
			Entry:  *entry,
		})
		if err != nil {
			fatalf("link: %v", err)
		}
		perm = 0o755
	}

	// Remove any pre-existing file before writing.
	// os.WriteFile preserves the mode of an existing file and ignores the
	// perm argument, so a previously-built 0644 binary stays non-executable
	// across rebuilds unless we create a fresh file every time.
	_ = os.Remove(*outPath)

	if err := os.WriteFile(*outPath, out, perm); err != nil {
		fatalf("write %s: %v", *outPath, err)
	}

	// Belt-and-suspenders: explicitly set executable bits for native output
	// in case the process umask stripped them from the WriteFile call above.
	if !*emitWasm {
		if err := os.Chmod(*outPath, 0o755); err != nil {
			fatalf("chmod %s: %v", *outPath, err)
		}
	}

	fmt.Printf("wrote %s (%d bytes)\n", *outPath, len(out))
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "vertex: "+format+"\n", args...)
	os.Exit(1)
}