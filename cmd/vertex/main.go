// cmd/main.go — Vertex language compiler driver.
//
// Usage:
//
//	vertex -o <output>            <file.vs> [file.vs ...]   # native binary
//	vertex -wasm -o <output.wasm> <file.vs> [file.vs ...]  # wasm module
//
// Flags:
//
//	-o string      Output file path (default "a.out" or "a.wasm")
//	-wasm          Emit a .wasm module instead of a native binary
//	-platform      Target platform: linux | darwin | windows (default: linux)
//	-tags          Comma-separated extra build tags (e.g. "syscalls")
//	-module        Module path prefix (e.g. github.com/acme/myapp)
//	-root          Module root directory (default: directory of first source file)
//	-entry         Entry-point function name (default: "main")
//	-no-dce        Disable dead-code elimination (DCE is on by default)
//	-v             Print version and exit
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vertex-language/vertex/compiler"
)

func main() {
	var (
		outPath  = flag.String("o", "", "output file path")
		wasmMode = flag.Bool("wasm", false, "emit .wasm module")
		platform = flag.String("platform", "linux", "target platform: linux | darwin | windows")
		tagsRaw  = flag.String("tags", "", "comma-separated extra build tags")
		modPath  = flag.String("module", "", "module path prefix")
		modRoot  = flag.String("root", "", "module root directory")
		entry    = flag.String("entry", "", `entry-point function name (default "main")`)
		noDCE    = flag.Bool("no-dce", false, "disable dead-code elimination")
		verbose  = flag.Bool("v", false, "print version and exit")
	)
	flag.Usage = usage
	flag.Parse()

	if *verbose {
		fmt.Printf("vertex compiler %s\n", compiler.Version)
		os.Exit(0)
	}

	sources := flag.Args()
	if len(sources) == 0 {
		fmt.Fprintln(os.Stderr, "error: no source files specified")
		usage()
		os.Exit(1)
	}

	// Validate source files exist.
	for _, s := range sources {
		if _, err := os.Stat(s); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	// Default output path.
	out := *outPath
	if out == "" {
		if *wasmMode {
			out = "a.wasm"
		} else {
			out = "a.out"
		}
	}

	// Extra build tags.
	var extraTags []string
	if *tagsRaw != "" {
		for _, t := range strings.Split(*tagsRaw, ",") {
			if t = strings.TrimSpace(t); t != "" {
				extraTags = append(extraTags, t)
			}
		}
	}

	opts := compiler.Options{
		Platform:   *platform,
		ExtraTags:  extraTags,
		OutputWasm: *wasmMode,
		ModulePath: *modPath,
		ModuleRoot: *modRoot,
		Entry:      *entry,
		DCE:        !*noDCE, // DCE is on by default; -no-dce turns it off
	}

	binary, err := compiler.CompileFiles(sources, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(out, binary, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error: write %s: %v\n", out, err)
		os.Exit(1)
	}

	fmt.Printf("wrote %s (%d bytes)\n", out, len(binary))
}

func usage() {
	fmt.Fprintf(os.Stderr, `Vertex language compiler %s

Usage:
  vertex [flags] <file.vs> [file.vs ...]

Flags:
`, compiler.Version)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, `
Examples:
  vertex -o server server.vs
  vertex -wasm -o main.wasm main.vs
  vertex -platform darwin -o myapp main.vs
  vertex -platform linux -tags syscalls -o net net.vs
  vertex -no-dce -o debug main.vs
  vertex -entry start -wasm -o lib.wasm lib.vs
`)
}