// cmd/vertex/main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
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
func (f *stringListFlag) Set(v string) error { *f = append(*f, v); return nil }

func expandShortFlags(args []string) []string {
	valueFlags := map[byte]bool{'l': true, 'L': true, 'I': true, 'o': true}
	out := make([]string, 0, len(args))
	for _, arg := range args {
		if len(arg) >= 3 && arg[0] == '-' && arg[1] != '-' && valueFlags[arg[1]] {
			out = append(out, arg[:2], arg[2:])
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
		testMode    bool
		testDir     string
		outputFile  string
		targetStr   string
		printVer    bool
		paths       stringListFlag
		libDirs     stringListFlag
		libs        stringListFlag
	)

	fs.BoolVar(&emitC, "emit-c", false, "emit C source code instead of native binary")
	fs.BoolVar(&compileOnly, "c", false, "compile and assemble, but do not link (outputs object file)")
	fs.BoolVar(&testMode, "test", false, "compile and run 'build test' functions, checking Expected(...) output")
	fs.StringVar(&testDir, "dir", "", "run tests recursively in `directory` (used with -test)")
	fs.StringVar(&outputFile, "o", "", "write output to `file`")
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
		fmt.Fprintf(stderr, "  vertex -o main main.vs                (build native executable)\n")
		fmt.Fprintf(stderr, "  vertex -lc -lm -o main main.vs        (link against libc and libm)\n")
		fmt.Fprintf(stderr, "  vertex -c -o main.o main.vs           (build object file only)\n")
		fmt.Fprintf(stderr, "  vertex -emit-c -o main.c main.vs      (emit C source)\n")
		fmt.Fprintf(stderr, "  vertex -test arithmetic_test.vs       (run test functions in file)\n")
		fmt.Fprintf(stderr, "  vertex -test -dir .                   (run all test functions recursively)\n")
	}

	if err := fs.Parse(expandShortFlags(args)); err != nil {
		return 2
	}
	if printVer {
		fmt.Fprintf(os.Stdout, "vertex %s\n", version)
		return 0
	}

	var inputFiles []string

	// Determine input file(s) based on recursive directory flag
	if testMode && testDir != "" {
		if fs.NArg() > 0 {
			fmt.Fprintf(stderr, "vertex: cannot specify both -dir and individual input files\n")
			return 2
		}
		err := filepath.WalkDir(testDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(d.Name(), ".vs") {
				inputFiles = append(inputFiles, path)
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(stderr, "vertex: error reading directory: %v\n", err)
			return 2
		}
		if len(inputFiles) == 0 {
			fmt.Fprintf(stderr, "vertex: no .vs files found in %s\n", testDir)
			return 1
		}
	} else {
		if fs.NArg() != 1 {
			fmt.Fprintf(stderr, "vertex: expected exactly 1 input file, got %d\n", fs.NArg())
			fs.Usage()
			return 2
		}
		inputFiles = append(inputFiles, fs.Arg(0))
	}

	cTarget, mTarget, ok := parseTarget(targetStr)
	if !ok {
		fmt.Fprintf(stderr, "vertex: unsupported target %q\n", targetStr)
		return 2
	}
	if mTarget != mir.TargetLinuxAMD64 && (len(libs) > 0 || len(libDirs) > 0) {
		fmt.Fprintf(stderr, "vertex: warning: -l / -L flags are only supported for Linux ELF targets; ignored\n")
	}

	cfg := compiler.Config{
		Target:      cTarget,
		SearchPaths: []string(paths),
	}

	// ── Test mode ─────────────────────────────────────────────────────────────
	if testMode {
		return runTests(inputFiles, cfg, mTarget, []string(libDirs), []string(libs), stderr)
	}

	// ── Normal compile ────────────────────────────────────────────────────────
	inputFile := inputFiles[0] // Only ever 1 element during normal compilation
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

	if compileOnly {
		if outputFile == "" {
			outputFile = replaceExt(inputFile, ".o")
			if mTarget == mir.TargetWindowsAMD64 {
				outputFile = replaceExt(inputFile, ".obj")
			}
		}
		objBytes, err := moduleToObject(cMod, mTarget)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		if err := writeOutput(outputFile, objBytes); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0
	}

	if outputFile == "" {
		outputFile = replaceExt(inputFile, "")
		if mTarget == mir.TargetWindowsAMD64 {
			outputFile += ".exe"
		}
	}
	if err := buildBinaryFromModule(cMod, outputFile, mTarget, []string(libDirs), []string(libs)); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	if mTarget != mir.TargetWindowsAMD64 && outputFile != "-" {
		os.Chmod(outputFile, 0755)
	}
	return 0
}

// ─────────────────────────────────────────────────────────────────────────────
// Test runner
// ─────────────────────────────────────────────────────────────────────────────

// runTests compiles each test function as its own binary across multiple files,
// executes it, and diffs stdout against the Expected(...) annotation.
func runTests(
	inputFiles []string,
	cfg compiler.Config,
	mTarget mir.Target,
	libDirs, libs []string,
	stderr io.Writer,
) int {
	comp := compiler.New(cfg)

	tmpDir, err := os.MkdirTemp("", "vertex-tests-*")
	if err != nil {
		fmt.Fprintf(stderr, "vertex: cannot create temp dir: %v\n", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	// On Linux, auto-prepend libc so printf resolves without the caller
	// needing to pass -lc explicitly for test mode.
	testLibs := append([]string(nil), libs...)
	if mTarget == mir.TargetLinuxAMD64 {
		hasC := false
		for _, l := range libs {
			if l == "c" {
				hasC = true
				break
			}
		}
		if !hasC {
			testLibs = append([]string{"c"}, testLibs...)
		}
	}

	passed, failed := 0, 0
	binCounter := 0 // Ensures unique binary names if tests across files share names

	for _, file := range inputFiles {
		infos, modules, err := comp.CompileTestFile(file)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %s: %v\n", file, err)
			failed++
			continue
		}
		if len(infos) == 0 {
			continue // Gracefully skip files that don't have testing functions
		}

		for i, info := range infos {
			binName := fmt.Sprintf("%s_%d", info.Name, binCounter)
			binCounter++

			binPath := filepath.Join(tmpDir, binName)
			if mTarget == mir.TargetWindowsAMD64 {
				binPath += ".exe"
			}

			if buildErr := buildBinaryFromModule(modules[i], binPath, mTarget, libDirs, testLibs); buildErr != nil {
				fmt.Printf("FAIL\t%s\t(%s)\n\t\t[build: %v]\n", info.Name, file, buildErr)
				failed++
				continue
			}

			out, runErr := exec.Command(binPath).Output()
			if runErr != nil {
				fmt.Printf("FAIL\t%s\t(%s)\n\t\t[run: %v]\n", info.Name, file, runErr)
				failed++
				continue
			}

			got := strings.TrimRight(string(out), "\r\n")
			if got == info.Expected {
				fmt.Printf("ok  \t%s\t(%s)\n", info.Name, file)
				passed++
			} else {
				fmt.Printf("FAIL\t%s\t(%s)\n\t\twant: %q\n\t\t got: %q\n", info.Name, file, info.Expected, got)
				failed++
			}
		}
	}

	if passed == 0 && failed == 0 {
		fmt.Fprintln(stderr, "vertex: no test functions found (functions need 'test' qualifier and Expected(...) return type)")
		return 1
	}

	fmt.Printf("\n%d passed, %d failed\n", passed, failed)
	if failed > 0 {
		return 1
	}
	return 0
}

// ─────────────────────────────────────────────────────────────────────────────
// Backend pipeline helpers
// ─────────────────────────────────────────────────────────────────────────────

// moduleToObject runs MIR lowering + encoding and returns raw object-file bytes.
func moduleToObject(mod *cir.Module, mTarget mir.Target) ([]byte, error) {
	abi := mirAMD64.ABIForTarget(mTarget)
	mirMod := mir.NewModule(mTarget)

	cFrames, err := cLower.LowerModule(mod, mirMod, abi)
	if err != nil {
		return nil, fmt.Errorf("mir lowering failed: %w", err)
	}
	enc := encAMD64.NewEncoder(abi)
	sections, err := enc.Encode(mirMod, cFrames)
	if err != nil {
		return nil, fmt.Errorf("encoding failed: %w", err)
	}

	switch mTarget {
	case mir.TargetLinuxAMD64:
		obj := objELF.NewObjectFile(mTarget)
		for _, s := range sections {
			obj.AddSection(s)
		}
		return obj.Serialize()
	case mir.TargetWindowsAMD64:
		obj := objCOFF.NewObjectFile(mTarget)
		for _, s := range sections {
			obj.AddSection(s)
		}
		return obj.Serialize()
	case mir.TargetDarwinAMD64:
		obj := objMachO.NewObjectFile(mTarget)
		for _, s := range sections {
			obj.AddSection(s)
		}
		return obj.Serialize()
	}
	return nil, fmt.Errorf("unsupported target for object file")
}

// buildBinaryFromModule runs the full backend pipeline (object file → link)
// and writes an executable to outputPath.
func buildBinaryFromModule(
	mod *cir.Module,
	outputPath string,
	mTarget mir.Target,
	libDirs, libs []string,
) error {
	objBytes, err := moduleToObject(mod, mTarget)
	if err != nil {
		return err
	}

	objName := filepath.Base(outputPath) + ".o"
	if mTarget == mir.TargetWindowsAMD64 {
		objName = filepath.Base(outputPath) + ".obj"
	}

	var exeBytes []byte
	switch mTarget {
	case mir.TargetLinuxAMD64:
		linker := lnkELF.NewLinker(lnkELF.ArchAMD64)
		if addErr := linker.AddObject(objName, objBytes); addErr != nil {
			return fmt.Errorf("add object: %w", addErr)
		}
		searchDirs := append(append([]string(nil), libDirs...), elfLibSearchDirs()...)
		for _, p := range searchDirs {
			linker.AddLibraryPath(p)
		}
		for _, name := range libs {
			data, soname, ferr := findSharedLib(name, searchDirs)
			if ferr != nil {
				return ferr
			}
			if lerr := linker.AddDynamicLibrary(soname, data); lerr != nil {
				return fmt.Errorf("-l%s: %w", name, lerr)
			}
		}
		exeBytes, err = linker.Link()

	case mir.TargetWindowsAMD64:
		linker := lnkPE.NewLinker(lnkPE.ArchAMD64)
		linker.AddObject(objName, objBytes)
		exeBytes, err = linker.Link()

	case mir.TargetDarwinAMD64:
		linker := lnkMachO.NewLinker(lnkMachO.ArchAMD64)
		linker.AddObject(objName, objBytes)
		exeBytes, err = linker.Link()

	default:
		return fmt.Errorf("unsupported target for linking")
	}
	if err != nil {
		return fmt.Errorf("linking failed: %w", err)
	}

	if err := writeOutput(outputPath, exeBytes); err != nil {
		return err
	}
	if mTarget != mir.TargetWindowsAMD64 {
		os.Chmod(outputPath, 0755)
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// ELF shared-library resolution (unchanged)
// ─────────────────────────────────────────────────────────────────────────────

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

func findSharedLib(name string, searchDirs []string) ([]byte, string, error) {
	base := "lib" + name + ".so"
	for _, dir := range searchDirs {
		matches, _ := filepath.Glob(filepath.Join(dir, base+".[0-9]*"))
		for _, m := range matches {
			data, err := os.ReadFile(m)
			if err != nil || !isELF(data) {
				continue
			}
			return data, filepath.Base(m), nil
		}
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

func isELF(data []byte) bool {
	return len(data) >= 4 &&
		data[0] == 0x7f && data[1] == 'E' && data[2] == 'L' && data[3] == 'F'
}

// ─────────────────────────────────────────────────────────────────────────────
// Misc helpers (unchanged)
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