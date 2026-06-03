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
// Repeatable -I flag
// ─────────────────────────────────────────────────────────────────────────────

type searchPathFlag []string

func (f *searchPathFlag) String() string {
	return strings.Join(*f, string(filepath.ListSeparator))
}

func (f *searchPathFlag) Set(v string) error {
	*f = append(*f, v)
	return nil
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
		paths       searchPathFlag
	)

	fs.BoolVar(&emitC, "emit-c", false, "emit C source code instead of native binary")
	fs.BoolVar(&compileOnly, "c", false, "compile and assemble, but do not link (outputs object file)")
	fs.StringVar(&outputFile, "o", "", "write output to `file` (default: input stem + \".o\" / \".c\" / or executable)")
	fs.StringVar(&targetStr, "target", "linux-amd64",
		"target platform: linux-amd64 (default), darwin-amd64, windows-amd64")
	fs.BoolVar(&printVer, "version", false, "print version and exit")
	fs.BoolVar(&printVer, "v", false, "shorthand for -version")
	fs.Var(&paths, "I", "add a package search `path` (repeatable)")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Vertex compiler %s\n\n", version)
		fmt.Fprintf(stderr, "Usage:\n  vertex [flags] <source.vs>\n\nFlags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(stderr, "\nExamples:\n")
		fmt.Fprintf(stderr, "  vertex -o main main.vs          (Build native executable)\n")
		fmt.Fprintf(stderr, "  vertex -c -o main.o main.vs     (Build object file)\n")
		fmt.Fprintf(stderr, "  vertex -emit-c -o main.c main.vs (Emit C source)\n")
	}

	if err := fs.Parse(args); err != nil {
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

	// Print warnings / notes
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
		linker.AddObject(filepath.Base(inputFile)+".o", objBytes)
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

	// Mark the binary as executable on Unix systems.
	if mTarget != mir.TargetWindowsAMD64 && outputFile != "-" {
		os.Chmod(outputFile, 0755)
	}

	return 0
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

// parseTarget returns both the C IR target and the MIR target.
// Note: Currently limited to AMD64 targets as they are the fully-supported 
// paths demonstrated in the provided encoding toolchain.
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