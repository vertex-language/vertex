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
	cirlower "github.com/vertex-language/ir/c/lower"
	mirir "github.com/vertex-language/ir/mir/ir"
	mircompiler "github.com/vertex-language/ir/mir/compiler"
	mirprofile "github.com/vertex-language/ir/mir/profile"
	mirtext "github.com/vertex-language/ir/mir/encoding/text"
	mirvalidate "github.com/vertex-language/ir/mir/validate"

	"github.com/vertex-language/objectfile/object"
	objcoff "github.com/vertex-language/objectfile/coff"
	objelf "github.com/vertex-language/objectfile/elf"
	objmacho "github.com/vertex-language/objectfile/macho"

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

// expandShortFlags lets the user write -lc instead of -l c (and -Ldir, -ofile).
func expandShortFlags(args []string) []string {
	valueFlags := map[byte]bool{'l': true, 'L': true, 'o': true}
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
		emitMIR     bool
		compileOnly bool
		testMode    bool
		testDir     string
		outputFile  string
		targetStr   string
		printVer    bool
		packagesDir string
		rebuild     bool
		libDirs     stringListFlag
		libs        stringListFlag
	)

	fs.BoolVar(&emitC, "emit-c", false, "emit C source instead of a native binary (accepts file or directory)")
	fs.BoolVar(&emitMIR, "emit-mir", false, "emit MIR text (S-expression) instead of a native binary (accepts file or directory)")
	fs.BoolVar(&compileOnly, "c", false, "compile and assemble but do not link (outputs object file)")
	fs.BoolVar(&testMode, "test", false, "compile and run test functions, checking Expected(...) output")
	fs.StringVar(&testDir, "dir", "", "run tests recursively in `directory` (used with -test)")
	fs.StringVar(&outputFile, "o", "", "write output to `file`")
	fs.StringVar(&targetStr, "target", "linux-amd64",
		"target platform: linux-amd64 (default), darwin-amd64, windows-amd64")
	fs.BoolVar(&printVer, "version", false, "print version and exit")
	fs.BoolVar(&printVer, "v", false, "shorthand for -version")
	fs.StringVar(&packagesDir, "packages-dir", defaultPackagesDir(),
		"Vertex packages directory (overrides $VERTEX_PATH)")
	fs.BoolVar(&rebuild, "rebuild", false, "force rebuild of cached packages")
	fs.Var(&libDirs, "L", "add a library search `dir` (ELF targets only, repeatable)")
	fs.Var(&libs, "l", "link against lib`name` e.g. -lc -lm (ELF targets only, repeatable)")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Vertex compiler %s\n\n", version)
		fmt.Fprintf(stderr, "Usage:\n  vertex [flags] <source.vs|package/>\n\nFlags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(stderr, "\nExamples:\n")
		fmt.Fprintf(stderr, "  vertex -o main main.vs                        (build native executable)\n")
		fmt.Fprintf(stderr, "  vertex -lc -o main main.vs                    (link against libc)\n")
		fmt.Fprintf(stderr, "  vertex -c -o main.o main.vs                   (build object file only)\n")
		fmt.Fprintf(stderr, "  vertex -emit-c -o main.c main.vs              (emit C source for a file)\n")
		fmt.Fprintf(stderr, "  vertex -emit-c -o arrays.c ./packages/arrays/ (emit C source for a package)\n")
		fmt.Fprintf(stderr, "  vertex -emit-mir -o main.mir main.vs          (emit MIR text for a file)\n")
		fmt.Fprintf(stderr, "  vertex -emit-mir -o arrays.mir ./packages/arrays/ (emit MIR text for a package)\n")
		fmt.Fprintf(stderr, "  vertex -test arithmetic_test.vs               (run test functions in file)\n")
		fmt.Fprintf(stderr, "  vertex -test arithmetic_test.vs -o ./out      (run tests and save binaries)\n")
		fmt.Fprintf(stderr, "  vertex -test -dir .                           (run all tests recursively)\n")
		fmt.Fprintf(stderr, "  vertex -rebuild -o main main.vs               (force rebuild of cached packages)\n")
	}

	if err := fs.Parse(expandShortFlags(args)); err != nil {
		return 2
	}
	if printVer {
		fmt.Fprintf(os.Stdout, "vertex %s\n", version)
		return 0
	}

	if emitC && emitMIR {
		fmt.Fprintf(stderr, "vertex: -emit-c and -emit-mir are mutually exclusive\n")
		return 2
	}

	var inputFiles []string

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
			fmt.Fprintf(stderr, "vertex: expected exactly 1 input file or directory, got %d\n", fs.NArg())
			fs.Usage()
			return 2
		}
		if testMode {
			inputFiles = append(inputFiles, fs.Arg(0))
		}
	}

	cTarget, objTarget, ok := parseTarget(targetStr)
	if !ok {
		fmt.Fprintf(stderr, "vertex: unsupported target %q\n", targetStr)
		return 2
	}
	if objTarget.OS != object.OSLinux && (len(libs) > 0 || len(libDirs) > 0) {
		fmt.Fprintf(stderr, "vertex: warning: -l / -L flags are only supported for Linux ELF targets; ignored\n")
	}

	cfg := compiler.Config{
		Target:      cTarget,
		PackagesDir: packagesDir,
		ObjectFunc: func(mod *cir.Module) ([]byte, error) {
			return moduleToObject(mod, objTarget)
		},
		Rebuild: rebuild,
	}

	if testMode {
		return runTests(inputFiles, cfg, objTarget, []string(libDirs), []string(libs), outputFile, stderr)
	}

	inputPath := fs.Arg(0)
	inputIsDir := isDir(inputPath)

	comp := compiler.New(cfg)
	cMod, pkgObjs, err := comp.CompileFile(inputPath)
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
			if inputIsDir {
				outputFile = replaceExt(filepath.Base(inputPath), ".c")
			} else {
				outputFile = replaceExt(inputPath, ".c")
			}
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

	if emitMIR {
		if outputFile == "" {
			if inputIsDir {
				outputFile = replaceExt(filepath.Base(inputPath), ".mir")
			} else {
				outputFile = replaceExt(inputPath, ".mir")
			}
		}
		mirSource, err := moduleToMIRText(cMod)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: MIR emission failed: %v\n", err)
			return 1
		}
		if err := writeOutput(outputFile, mirSource); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0
	}

	if compileOnly {
		if outputFile == "" {
			base := inputPath
			if inputIsDir {
				base = filepath.Base(inputPath)
			}
			if objTarget.OS == object.OSWindows {
				outputFile = replaceExt(base, ".obj")
			} else {
				outputFile = replaceExt(base, ".o")
			}
		}
		objBytes, err := moduleToObject(cMod, objTarget)
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
		if inputIsDir {
			outputFile = filepath.Base(inputPath)
		} else {
			outputFile = replaceExt(inputPath, "")
		}
		if objTarget.OS == object.OSWindows {
			outputFile += ".exe"
		}
	}
	if err := buildBinaryFromModule(cMod, pkgObjs, outputFile, objTarget, []string(libDirs), []string(libs)); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	if objTarget.OS != object.OSWindows && outputFile != "-" {
		os.Chmod(outputFile, 0755)
	}
	return 0
}

// ─────────────────────────────────────────────────────────────────────────────
// Test runner
// ─────────────────────────────────────────────────────────────────────────────

func runTests(
	inputFiles []string,
	cfg compiler.Config,
	objTarget object.Target,
	libDirs, libs []string,
	outputFile string,
	stderr io.Writer,
) int {
	comp := compiler.New(cfg)

	tmpDir, err := os.MkdirTemp("", "vertex-tests-*")
	if err != nil {
		fmt.Fprintf(stderr, "vertex: cannot create temp dir: %v\n", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	testLibs := append([]string(nil), libs...)
	if objTarget == object.TargetLinuxAMD64 {
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

	// Count total tests across all files so we know whether -o needs suffixing.
	totalTests := 0
	type pendingFile struct {
		file    string
		infos   []compiler.TestFuncInfo
		modules []*cir.Module
		pkgObjs [][]byte
	}
	var pending []pendingFile

	for _, file := range inputFiles {
		infos, modules, pkgObjs, err := comp.CompileTestFile(file)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %s: %v\n", file, err)
			continue
		}
		if len(infos) == 0 {
			continue
		}
		totalTests += len(infos)
		pending = append(pending, pendingFile{file, infos, modules, pkgObjs})
	}

	if totalTests == 0 {
		fmt.Fprintln(stderr, "vertex: no test functions found (functions need 'test' qualifier and Expected(...) return type)")
		return 1
	}

	passed, failed := 0, 0
	binCounter := 0

	for _, pf := range pending {
		for i, info := range pf.infos {
			binCounter++

			// Determine binary output path.
			var binPath string
			if outputFile != "" {
				if totalTests == 1 {
					binPath = outputFile
				} else {
					binPath = fmt.Sprintf("%s_%s", outputFile, info.Name)
				}
			} else {
				binPath = filepath.Join(tmpDir, fmt.Sprintf("%s_%d", info.Name, binCounter))
			}
			if objTarget.OS == object.OSWindows && !strings.HasSuffix(binPath, ".exe") {
				binPath += ".exe"
			}

			if buildErr := buildBinaryFromModule(pf.modules[i], pf.pkgObjs, binPath, objTarget, libDirs, testLibs); buildErr != nil {
				fmt.Printf("FAIL\t%s\t(%s)\n\t\t[build: %v]\n", info.Name, pf.file, buildErr)
				failed++
				continue
			}

			if objTarget.OS != object.OSWindows {
				os.Chmod(binPath, 0755)
			}

			out, runErr := exec.Command(binPath).Output()
			if runErr != nil {
				fmt.Printf("FAIL\t%s\t(%s)\n\t\t[run: %v]\n", info.Name, pf.file, runErr)
				failed++
				continue
			}

			got := strings.TrimRight(string(out), "\r\n")
			if got == info.Expected {
				fmt.Printf("ok  \t%s\t(%s)\n", info.Name, pf.file)
				passed++
			} else {
				fmt.Printf("FAIL\t%s\t(%s)\n\t\twant: %q\n\t\t got: %q\n", info.Name, pf.file, info.Expected, got)
				failed++
			}
		}
	}

	fmt.Printf("\n%d passed, %d failed\n", passed, failed)
	if failed > 0 {
		return 1
	}
	return 0
}

// ─────────────────────────────────────────────────────────────────────────────
// Backend pipeline
// ─────────────────────────────────────────────────────────────────────────────

// lowerToMIR is the shared first half of the backend pipeline.  It lowers a
// CIR module to Machine IR and validates the result against the resolved target
// profile.  Both moduleToObject and moduleToMIRText call this so that the
// lowering and validation logic lives in exactly one place.
func lowerToMIR(cmod *cir.Module) (*mirir.Module, mirprofile.Profile, error) {
	mirMod, err := cirlower.NewLowerMIR(cmod)
	if err != nil {
		return nil, mirprofile.Profile{}, fmt.Errorf("c→mir lowering failed: %w", err)
	}

	prof, err := mirprofile.Resolve(mirMod.Target)
	if err != nil {
		return nil, mirprofile.Profile{}, fmt.Errorf("profile resolution failed: %w", err)
	}

	if diags := mirvalidate.Validate(mirMod, prof); len(diags) > 0 {
		var sb strings.Builder
		for _, d := range diags {
			fmt.Fprintln(&sb, d)
		}
		return nil, mirprofile.Profile{}, fmt.Errorf("MIR validation failed:\n%s", sb.String())
	}

	return mirMod, prof, nil
}

// moduleToMIRText lowers cmod to MIR and returns the S-expression text
// representation as a UTF-8 byte slice suitable for writing to a .mir file.
func moduleToMIRText(cmod *cir.Module) ([]byte, error) {
	mirMod, _, err := lowerToMIR(cmod)
	if err != nil {
		return nil, err
	}
	return []byte(mirtext.Format(mirMod)), nil
}

// moduleToObject lowers cmod all the way to a relocatable object file and
// returns its raw bytes.
//
// Full pipeline:
//  1. lowerToMIR       — CIR → MIR + profile resolution + validation
//  2. newObjectBuilder — pick ELF / COFF / Mach-O for objTarget
//  3. mircompiler.Compile — emit sections / symbols / relocs into the builder
//  4. builder.Serialize   — produce the on-disk bytes
func moduleToObject(cmod *cir.Module, objTarget object.Target) ([]byte, error) {
	mirMod, prof, err := lowerToMIR(cmod)
	if err != nil {
		return nil, err
	}

	out := newObjectBuilder(objTarget)
	if err := mircompiler.Compile(mirMod, prof, out); err != nil {
		return nil, fmt.Errorf("MIR compilation failed: %w", err)
	}

	return out.Serialize()
}

// newObjectBuilder selects the correct object.Builder for the given target.
func newObjectBuilder(t object.Target) object.Builder {
	switch t.OS {
	case object.OSLinux, object.OSFreestanding:
		return objelf.NewFile(t)
	case object.OSDarwin:
		return objmacho.NewFile(t)
	case object.OSWindows:
		return objcoff.NewFile(t)
	default:
		panic(fmt.Sprintf("newObjectBuilder: unsupported OS %v", t.OS))
	}
}

// buildBinaryFromModule compiles cmod to an object file, links it together
// with all package objects and any requested shared libraries, and writes a
// native executable to outputPath.
func buildBinaryFromModule(
	cmod *cir.Module,
	pkgObjs [][]byte,
	outputPath string,
	objTarget object.Target,
	libDirs, libs []string,
) error {
	objBytes, err := moduleToObject(cmod, objTarget)
	if err != nil {
		return err
	}

	objExt := ".o"
	if objTarget.OS == object.OSWindows {
		objExt = ".obj"
	}
	mainObjName := filepath.Base(outputPath) + objExt

	var exeBytes []byte

	switch objTarget {
	case object.TargetLinuxAMD64:
		linker := lnkELF.NewLinker(lnkELF.ArchAMD64)
		
		// 1. Setup search paths
		searchDirs := append(append([]string(nil), libDirs...), elfLibSearchDirs()...)
		for _, p := range searchDirs {
			linker.AddLibraryPath(p)
		}

		// Check if we are hosted (linking against libc)
		isHosted := false
		for _, l := range libs {
			if l == "c" {
				isHosted = true
				break
			}
		}

		// 2. Add CRT Prologue (provides _start and initialization)
		if isHosted {
			for _, crt := range []string{"crt1.o", "crti.o"} {
				data, err := findCrtObj(crt, searchDirs)
				if err != nil {
					return fmt.Errorf("hosted link requires %s: %w", crt, err)
				}
				if err := linker.AddObject(crt, data); err != nil {
					return err
				}
			}
		}

		// 3. Add your compiled Vertex objects
		if addErr := linker.AddObject(mainObjName, objBytes); addErr != nil {
			return fmt.Errorf("add object: %w", addErr)
		}
		for i, pkgObj := range pkgObjs {
			if addErr := linker.AddObject(fmt.Sprintf("pkg%d.o", i), pkgObj); addErr != nil {
				return fmt.Errorf("add package object %d: %w", i, addErr)
			}
		}

		// 4. Add Dynamic Libraries (-lc, etc.)
		for _, name := range libs {
			data, soname, ferr := findSharedLib(name, searchDirs)
			if ferr != nil {
				return ferr
			}
			if lerr := linker.AddDynamicLibrary(soname, data); lerr != nil {
				return fmt.Errorf("-l%s: %w", name, lerr)
			}
		}

		// 5. Add CRT Epilogue (provides cleanup/destructors)
		if isHosted {
			data, err := findCrtObj("crtn.o", searchDirs)
			if err != nil {
				return fmt.Errorf("hosted link requires crtn.o: %w", err)
			}
			if err := linker.AddObject("crtn.o", data); err != nil {
				return err
			}
		}

		exeBytes, err = linker.Link()

	case object.TargetWindowsAMD64:
		linker := lnkPE.NewLinker(lnkPE.ArchAMD64)
		linker.AddObject(mainObjName, objBytes)
		for i, pkgObj := range pkgObjs {
			linker.AddObject(fmt.Sprintf("pkg%d.obj", i), pkgObj)
		}
		exeBytes, err = linker.Link()

	case object.TargetDarwinAMD64:
		linker := lnkMachO.NewLinker(lnkMachO.ArchAMD64)
		linker.AddObject(mainObjName, objBytes)
		for i, pkgObj := range pkgObjs {
			linker.AddObject(fmt.Sprintf("pkg%d.o", i), pkgObj)
		}
		exeBytes, err = linker.Link()

	default:
		return fmt.Errorf("unsupported target for linking: %v", objTarget)
	}
	if err != nil {
		return fmt.Errorf("linking failed: %w", err)
	}

	if err := writeOutput(outputPath, exeBytes); err != nil {
		return err
	}
	if objTarget.OS != object.OSWindows {
		os.Chmod(outputPath, 0755)
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// ELF shared-library resolution
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

func findCrtObj(name string, searchDirs []string) ([]byte, error) {
	for _, dir := range searchDirs {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err == nil && isELF(data) {
			return data, nil
		}
	}
	return nil, fmt.Errorf("c runtime object %s not found (searched: %s)", name, strings.Join(searchDirs, ", "))
}

func findSharedLib(name string, searchDirs []string) ([]byte, string, error) {
	base := "lib" + name + ".so"
	for _, dir := range searchDirs {
		// Prefer versioned sonames (libfoo.so.6) over the bare link name.
		matches, _ := filepath.Glob(filepath.Join(dir, base+".[0-9]*"))
		for _, m := range matches {
			data, err := os.ReadFile(m)
			if err != nil || !isELF(data) {
				continue
			}
			return data, filepath.Base(m), nil
		}
		// Fall back to the unversioned name.
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
// Misc helpers
// ─────────────────────────────────────────────────────────────────────────────

// parseTarget maps a user-supplied target string to the pair of type values
// that the vertex compiler and object-file library each need.
func parseTarget(s string) (cir.Target, object.Target, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "linux-amd64", "linux-x86_64", "linux/amd64":
		return cir.TargetLinuxAMD64, object.TargetLinuxAMD64, true
	case "darwin-amd64", "darwin-x86_64", "macos-amd64":
		return cir.TargetDarwinAMD64, object.TargetDarwinAMD64, true
	case "windows-amd64", "windows-x86_64", "windows/amd64":
		return cir.TargetWindowsAMD64, object.TargetWindowsAMD64, true
	default:
		return cir.TargetUnknown, object.Target{}, false
	}
}

// defaultPackagesDir returns $VERTEX_PATH if set, otherwise ~/.vertex/packages.
func defaultPackagesDir() string {
	if p := os.Getenv("VERTEX_PATH"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".vertex", "packages")
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

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
	if ext := filepath.Ext(path); ext != "" {
		return path[:len(path)-len(ext)] + newExt
	}
	return path + newExt
}