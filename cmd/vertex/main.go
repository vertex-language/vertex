// cmd/vertex/main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	cir "github.com/vertex-language/ir/c"
	"github.com/vertex-language/ir/mir"
	mirAMD64 "github.com/vertex-language/ir/mir/amd64"
	cLower "github.com/vertex-language/ir/mir/amd64/c"

	encAMD64 "github.com/vertex-language/encoder/amd64"

	objCOFF "github.com/vertex-language/objectfile/coff"
	objELF "github.com/vertex-language/objectfile/elf"
	objMachO "github.com/vertex-language/objectfile/macho"

	lnkELF "github.com/vertex-language/linker/elf"
	lnkMachO "github.com/vertex-language/linker/macho"
	lnkPE "github.com/vertex-language/linker/pe"

	"github.com/vertex-language/vertex/compiler"
)

const version = "0.2.0"

// ─────────────────────────────────────────────────────────────────────────────
// Repeatable flags
// ─────────────────────────────────────────────────────────────────────────────

type stringListFlag []string

func (f *stringListFlag) String() string {
	return strings.Join(*f, string(filepath.ListSeparator))
}
func (f *stringListFlag) Set(v string) error {
	*f = append(*f, v)
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// expandShortFlags converts POSIX-style concatenated short flags into the two-
// token form that Go's flag package understands.  For example:
//
//	-lc   →  -l  c
//	-lm   →  -l  m
//	-L/my/dir  →  -L  /my/dir
//	-Isrc      →  -I  src
//
// Only single-character flags that accept a value are expanded; boolean flags
// (-c, -v) are left untouched because their names are a single character but
// they never have a fused value.
// ─────────────────────────────────────────────────────────────────────────────

func expandShortFlags(args []string) []string {
	// Single-char flags that take a value argument.
	valueFlags := map[byte]bool{
		'l': true,
		'L': true,
		'I': true,
		'o': true,
	}

	out := make([]string, 0, len(args))
	for _, arg := range args {
		// Must look like -X<value>: dash, one letter in valueFlags, then more chars.
		if len(arg) >= 3 && arg[0] == '-' && arg[1] != '-' && valueFlags[arg[1]] {
			out = append(out, arg[:2], arg[2:]) // e.g. "-lc" → "-l", "c"
		} else {
			out = append(out, arg)
		}
	}
	return out
}

// ─────────────────────────────────────────────────────────────────────────────
// Entry point
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}

func run(args []string, stderr io.Writer) int {
	fs := flag.NewFlagSet("vertex", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		emitC       bool
		compileOnly bool
		outputFile  string
		targetStr   string
		printVer    bool
		paths       stringListFlag // -I  source/package search paths
		libDirs     stringListFlag // -L  library search paths (ELF only)
		libs        stringListFlag // -l  libraries to link  (ELF only)
	)

	fs.BoolVar(&emitC, "emit-c", false, "emit C source code instead of native binary")
	fs.BoolVar(&compileOnly, "c", false, "compile and assemble, but do not link (outputs object file)")
	fs.StringVar(&outputFile, "o", "", "write output to `file` (default: input stem + \".o\" / \".c\" / or executable)")
	fs.StringVar(&targetStr, "target", "linux-amd64",
		"target platform: linux-amd64 (default), darwin-amd64, windows-amd64")
	fs.BoolVar(&printVer, "version", false, "print version and exit")
	fs.BoolVar(&printVer, "v", false, "shorthand for -version")
	fs.Var(&paths, "I", "add a package search `path` (repeatable)")
	fs.Var(&libDirs, "L", "add a library search `dir` (ELF targets only, repeatable)")
	fs.Var(&libs, "l", "link against lib`name` e.g. -lc -lm (ELF targets only, repeatable)")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Vertex compiler %s\n\n", version)
		fmt.Fprintf(stderr, "Usage:\n  vertex [flags] <source.vs>\n\nFlags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(stderr, "\nExamples:\n")
		fmt.Fprintf(stderr, "  vertex -o main main.vs              (build native executable)\n")
		fmt.Fprintf(stderr, "  vertex -lc -lm -o main main.vs      (link against libc and libm)\n")
		fmt.Fprintf(stderr, "  vertex -c -o main.o main.vs         (build object file only)\n")
		fmt.Fprintf(stderr, "  vertex -emit-c -o main.c main.vs    (emit C source)\n")
	}

	// Expand POSIX-style fused short flags (e.g. -lc → -l c) before parsing.
	if err := fs.Parse(expandShortFlags(args)); err != nil {
		return 2
	}

	if printVer {
		fmt.Fprintf(os.Stdout, "vertex %s\n", version)
		return 0
	}

	if fs.NArg() != 1 {
		fmt.Fprintf(stderr, "vertex: expected exactly 1 input file, got %d\n", fs.NArg())
		fs.Usage()
		return 2
	}
	inputFile := fs.Arg(0)

	// ── Target Resolution ─────────────────────────────────────────────────────
	cTarget, mTarget, ok := parseTarget(targetStr)
	if !ok {
		fmt.Fprintf(stderr, "vertex: unsupported target %q\n", targetStr)
		return 2
	}

	// Warn if -l / -L are used on non-ELF targets (accepted but ignored).
	if mTarget != mir.TargetLinuxAMD64 && (len(libs) > 0 || len(libDirs) > 0) {
		fmt.Fprintf(stderr, "vertex: warning: -l / -L flags are only supported for Linux ELF targets; ignored\n")
	}

	// ── 1. Compile to C IR (Frontend) ─────────────────────────────────────────
	cfg := compiler.Config{
		Target:      cTarget,
		SearchPaths: []string(paths),
	}
	comp := compiler.New(cfg)

	cMod, err := comp.CompileFile(inputFile)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	for _, d := range comp.Diagnostics().All() {
		if d.Severity != compiler.SevError {
			fmt.Fprintln(stderr, d)
		}
	}

	// ── Branch: Emit C source directly ────────────────────────────────────────
	if emitC {
		if outputFile == "" {
			outputFile = replaceExt(inputFile, ".c")
		}
		cSource, err := cMod.EmitC()
		if err != nil {
			fmt.Fprintf(stderr, "vertex: c emission failed: %v\n", err)
			return 1
		}
		if err := writeOutput(outputFile, cSource); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0
	}

	// ── 2. Lower to MIR ───────────────────────────────────────────────────────
	abi := mirAMD64.ABIForTarget(mTarget)
	mirMod := mir.NewModule(mTarget)

	cFrames, err := cLower.LowerModule(cMod, mirMod, abi)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: mir lowering failed: %v\n", err)
		return 1
	}

	// ── 3. Encode to Machine Code ─────────────────────────────────────────────
	enc := encAMD64.NewEncoder(abi)
	sections, err := enc.Encode(mirMod, cFrames)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: encoding failed: %v\n", err)
		return 1
	}

	// ── 4. Assemble Object File ───────────────────────────────────────────────
	var objBytes []byte
	switch mTarget {
	case mir.TargetLinuxAMD64:
		objFile := objELF.NewObjectFile(mTarget)
		for _, s := range sections {
			objFile.AddSection(s)
		}
		objBytes, err = objFile.Serialize()

	case mir.TargetWindowsAMD64:
		objFile := objCOFF.NewObjectFile(mTarget)
		for _, s := range sections {
			objFile.AddSection(s)
		}
		objBytes, err = objFile.Serialize()

	case mir.TargetDarwinAMD64:
		objFile := objMachO.NewObjectFile(mTarget)
		for _, s := range sections {
			objFile.AddSection(s)
		}
		objBytes, err = objFile.Serialize()
	}

	if err != nil {
		fmt.Fprintf(stderr, "vertex: object file assembly failed: %v\n", err)
		return 1
	}

	// ── Branch: Compile Only (Output Object File) ─────────────────────────────
	if compileOnly {
		if outputFile == "" {
			outputFile = replaceExt(inputFile, ".o")
			if mTarget == mir.TargetWindowsAMD64 {
				outputFile = replaceExt(inputFile, ".obj")
			}
		}
		if err := writeOutput(outputFile, objBytes); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0
	}

	// ── 5. Link Executable ────────────────────────────────────────────────────
	if outputFile == "" {
		outputFile = replaceExt(inputFile, "")
		if mTarget == mir.TargetWindowsAMD64 {
			outputFile += ".exe"
		}
	}

	var exeBytes []byte
	switch mTarget {

	case mir.TargetLinuxAMD64:
		linker := lnkELF.NewLinker(lnkELF.ArchAMD64)
		if err := linker.AddObject(filepath.Base(inputFile)+".o", objBytes); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}

		// User-supplied -L paths prepended; built-in system dirs appended as fallback.
		searchDirs := append([]string(libDirs), elfLibSearchDirs()...)
		for _, p := range searchDirs {
			linker.AddLibraryPath(p)
		}

		// Resolve and load each -l<name> library.
		for _, name := range libs {
			data, soname, err := findSharedLib(name, searchDirs)
			if err != nil {
				fmt.Fprintf(stderr, "vertex: %v\n", err)
				return 1
			}
			if err := linker.AddDynamicLibrary(soname, data); err != nil {
				fmt.Fprintf(stderr, "vertex: -l%s: %v\n", name, err)
				return 1
			}
		}

		exeBytes, err = linker.Link()

	case mir.TargetWindowsAMD64:
		linker := lnkPE.NewLinker(lnkPE.ArchAMD64)
		linker.AddObject(filepath.Base(inputFile)+".obj", objBytes)
		exeBytes, err = linker.Link()

	case mir.TargetDarwinAMD64:
		linker := lnkMachO.NewLinker(lnkMachO.ArchAMD64)
		linker.AddObject(filepath.Base(inputFile)+".o", objBytes)
		exeBytes, err = linker.Link()
	}

	if err != nil {
		fmt.Fprintf(stderr, "vertex: linking failed: %v\n", err)
		return 1
	}

	if err := writeOutput(outputFile, exeBytes); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	if mTarget != mir.TargetWindowsAMD64 && outputFile != "-" {
		os.Chmod(outputFile, 0755)
	}

	return 0
}

// ─────────────────────────────────────────────────────────────────────────────
// ELF shared-library resolution
// ─────────────────────────────────────────────────────────────────────────────

// elfLibSearchDirs returns the standard AMD64 Linux shared-library directories
// in the order they are typically searched.
func elfLibSearchDirs() []string {
	return []string{
		"/lib/x86_64-linux-gnu",
		"/usr/lib/x86_64-linux-gnu",
		"/lib64",
		"/usr/lib64",
		"/lib",
		"/usr/lib",
		"/usr/local/lib",
	}
}

// findSharedLib locates the real ELF binary for -l<name>, trying versioned
// filenames before the bare .so name (which is often a GNU linker script).
// Returns the file bytes, the soname to register, and any error.
func findSharedLib(name string, searchDirs []string) ([]byte, string, error) {
	base := "lib" + name + ".so"
	for _, dir := range searchDirs {
		// Prefer versioned files (e.g. libc.so.6) — they are real ELF binaries.
		// The unversioned .so is frequently a GNU ld linker script, not an ELF.
		matches, _ := filepath.Glob(filepath.Join(dir, base+".[0-9]*"))
		for _, m := range matches {
			data, err := os.ReadFile(m)
			if err != nil || !isELF(data) {
				continue
			}
			return data, filepath.Base(m), nil
		}

		// Fall back to the unversioned name if it really is an ELF.
		path := filepath.Join(dir, base)
		data, err := os.ReadFile(path)
		if err != nil || !isELF(data) {
			continue
		}
		return data, base, nil
	}
	return nil, "", fmt.Errorf("-l%s: library not found (searched: %s)",
		name, strings.Join(searchDirs, ", "))
}

// isELF reports whether data begins with the ELF magic bytes.
func isELF(data []byte) bool {
	return len(data) >= 4 &&
		data[0] == 0x7f && data[1] == 'E' && data[2] == 'L' && data[3] == 'F'
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func writeOutput(path string, data []byte) error {
	if path == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("cannot create output directory %s: %w", dir, err)
		}
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("cannot write %s: %w", path, err)
	}
	return nil
}

func replaceExt(path, newExt string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return path + newExt
	}
	return path[:len(path)-len(ext)] + newExt
}

func parseTarget(s string) (cir.Target, mir.Target, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "linux-amd64", "linux-x86_64", "linux/amd64":
		return cir.TargetLinuxAMD64, mir.TargetLinuxAMD64, true
	case "darwin-amd64", "darwin-x86_64", "macos-amd64":
		return cir.TargetDarwinAMD64, mir.TargetDarwinAMD64, true
	case "windows-amd64", "windows-x86_64", "windows/amd64":
		return cir.TargetWindowsAMD64, mir.TargetWindowsAMD64, true
	default:
		return cir.TargetUnknown, 0, false
	}
}