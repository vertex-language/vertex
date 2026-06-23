// pipeline.go
package driver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vertex-language/ir/vertex/ast"
	virlower "github.com/vertex-language/ir/vertex/lower/vir"
	mirlower "github.com/vertex-language/ir/vertex/lower/mir"

	virtext   "github.com/vertex-language/ir/vertex/encoding/text"
	virbinary "github.com/vertex-language/ir/vertex/encoding/binary"

	"github.com/vertex-language/ir/machine"
	mirtext   "github.com/vertex-language/ir/machine/encoding/text/mir"
	objbridge "github.com/vertex-language/ir/machine/object"

	"github.com/vertex-language/objectfile/object"
	"github.com/vertex-language/objectfile/elf"
	"github.com/vertex-language/objectfile/coff"
	"github.com/vertex-language/objectfile/macho"

	linkerelf   "github.com/vertex-language/linker/elf"
	linkermacho "github.com/vertex-language/linker/macho"
	linkerpe    "github.com/vertex-language/linker/pe"
)

// emit is the top-level pipeline entry for normal (non-test) compilation.
// It parses the input from disk then hands off to emitPackage.
func emit(cfg config, stderr io.Writer) int {
	if cfg.mode == modeTest {
		return runTests(cfg, stderr)
	}
	if cfg.mode == modeDump {
		return dumpAll(cfg, stderr)
	}

	pkg, err := parseInput(cfg.input)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	return emitPackage(pkg, cfg, stderr)
}

// emitPackage runs pipeline stages 2–7 against an already-parsed package.
// It is also called by the test runner, which builds a synthetic ast.Package
// rather than reading one from disk.
func emitPackage(pkg *ast.Package, cfg config, stderr io.Writer) int {
	tri, err := parseTriple(cfg.target)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 2
	}

	// ── Stage 2: Lower AST → Vertex IR ───────────────────────────────────────
	virMod, virErr := virlower.NewLower(pkg, nil, cfg.target)
	if virErr != nil {
		fmt.Fprintln(stderr, virErr)
	}
	if virMod == nil {
		return 1
	}

	virMod.SetTarget(tri.virTargetString())

	if cfg.mode == modeVIR {
		if err := writeOutput(cfg.output, []byte(virtext.Format(virMod))); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return boolToCode(virErr != nil)
	}
	if cfg.mode == modeVBytes {
		data, err := virbinary.Marshal(virMod)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: vbytes encoding: %v\n", err)
			return 1
		}
		if err := writeOutput(cfg.output, data); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return boolToCode(virErr != nil)
	}

	if virErr != nil {
		return 1
	}

	// ── Stage 3: Lower Vertex IR → Machine IR ────────────────────────────────
	mirMod, err := mirlower.NewLower(virMod)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: MIR lowering: %v\n", err)
		return 1
	}
	if err := machine.Verify(mirMod); err != nil {
		fmt.Fprintf(stderr, "vertex: MIR verification: %v\n", err)
		return 1
	}

	if cfg.mode == modeMIR {
		if err := writeOutput(cfg.output, []byte(mirtext.PrintModule(mirMod))); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0
	}

	// ── Stage 4: Instruction selection, register allocation, assembly ─────────
	opts := codegenOptions{optLevel: cfg.optLevel, debugInfo: cfg.debugInfo}

	if cfg.mode == modeASM {
		text, err := compileToASM(mirMod, tri, opts)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: code generation: %v\n", err)
			return 1
		}
		if err := writeOutput(cfg.output, []byte(text)); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0
	}

	fns, err := compileToFuncs(mirMod, tri, opts)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: code generation: %v\n", err)
		return 1
	}

	// ── Stage 5: Build main object file ───────────────────────────────────────
	objTarget, err := tri.objectTarget()
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 2
	}

	sections := buildSections(fns, mirMod)
	objBytes, err := marshalObject(tri, objTarget, sections)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: object serialization: %v\n", err)
		return 1
	}

	if cfg.mode == modeObj {
		if err := writeOutput(cfg.output, objBytes); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0
	}

	// ── Stage 6: Resolve shared libraries and CRT objects ────────────────────
	if tri.os == "freestanding" {
		fmt.Fprintf(stderr, "vertex: cannot link a freestanding target; use -c/-emit-obj instead\n")
		return 2
	}

	libNames := extractDynLibs(virMod, tri)
	dynLibs, err := resolveLibs(libNames, tri, cfg.sysroot)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	crt, err := resolveCRT(tri, cfg.sysroot)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	// ── Stage 6.5: Compile & Link Reserved Runtime Directory ─────────────────
	runtimeObj, err := compileRuntime(cfg, tri)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: runtime compilation failed: %v\n", err)
		return 1
	}

	// ── Stage 7: Link to native executable ───────────────────────────────────
	exeBytes, err := linkObject(tri, objBytes, dynLibs, crt, runtimeObj)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: link: %v\n", err)
		return 1
	}
	if err := writeExe(cfg.output, exeBytes); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	return 0
}

// compileRuntime targets the reserved "runtime/" package directory.
// To fix slow compile times, it ONLY scans this specific directory and caches
// the compiled object file based on the target architecture.
func compileRuntime(cfg config, tri triple) ([]byte, error) {
	if cfg.packagesDir == "" {
		return nil, nil // No packages dir configured
	}

	runtimeDir := filepath.Join(cfg.packagesDir, "runtime")
	if !isDir(runtimeDir) {
		return nil, nil // No runtime directory exists, skip.
	}

	// The cached object file name includes the target so cross-compiling works safely.
	cacheFile := filepath.Join(runtimeDir, "runtime_"+tri.virTargetString()+".o")

	var files []*ast.File
	var newestSrc time.Time

	err := filepath.WalkDir(runtimeDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".vs") {
			return nil
		}
		
		// Track the latest modification time among the source files
		info, err := d.Info()
		if err == nil && info.ModTime().After(newestSrc) {
			newestSrc = info.ModTime()
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		f, err := ast.NewFile(path, src)
		if err != nil {
			return err
		}
		files = append(files, f)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed scanning runtime directory: %w", err)
	}
	if len(files) == 0 {
		return nil, nil // No .vs files in runtime/
	}

	// Fast path: Check if the cached object file is up-to-date
	if cacheInfo, err := os.Stat(cacheFile); err == nil {
		if !cacheInfo.ModTime().Before(newestSrc) {
			// Cache hit! Return the bytes immediately.
			return os.ReadFile(cacheFile)
		}
	}

	// Cache miss: We need to compile the runtime package
	pkg, err := ast.NewPackage(files)
	if err != nil {
		return nil, fmt.Errorf("failed to build runtime package: %w", err)
	}

	// Sub-pipeline for runtime
	virMod, errs := virlower.NewLower(pkg, nil, cfg.target)
	if errs != nil {
		return nil, fmt.Errorf("VIR lowering: %v", errs)
	}
	virMod.SetTarget(tri.virTargetString())

	mirMod, err := mirlower.NewLower(virMod)
	if err != nil {
		return nil, fmt.Errorf("MIR lowering: %w", err)
	}
	if err := machine.Verify(mirMod); err != nil {
		return nil, fmt.Errorf("MIR verify: %w", err)
	}

	opts := codegenOptions{optLevel: cfg.optLevel, debugInfo: cfg.debugInfo}
	fns, err := compileToFuncs(mirMod, tri, opts)
	if err != nil {
		return nil, fmt.Errorf("codegen: %w", err)
	}

	objTarget, err := tri.objectTarget()
	if err != nil {
		return nil, err
	}
	secs := buildSections(fns, mirMod)
	objBytes, err := marshalObject(tri, objTarget, secs)
	if err != nil {
		return nil, err
	}

	// Save to cache for the next run (ignore errors if folder is read-only)
	_ = os.WriteFile(cacheFile, objBytes, 0644)

	return objBytes, nil
}

func parseInput(input string) (*ast.Package, error) {
	var files []*ast.File
	if isDir(input) {
		entries, err := os.ReadDir(input)
		if err != nil {
			return nil, fmt.Errorf("cannot read directory %s: %w", input, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".vs") {
				continue
			}
			f, err := parseFile(filepath.Join(input, e.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, f)
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("no .vs files found in %s", input)
		}
	} else {
		f, err := parseFile(input)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return ast.NewPackage(files)
}

func parseFile(path string) (*ast.File, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return ast.NewFile(path, src)
}

func buildSections(fns []objbridge.AssembledFunc, m *machine.Module) []object.Section {
	secs := make([]object.Section, 0, 2)
	secs = append(secs, objbridge.BuildText(fns))
	secs = append(secs, objbridge.DataSections(m)...)
	return secs
}

func marshalObject(tri triple, tgt object.Target, sections []object.Section) ([]byte, error) {
	type objectFile interface {
		AddSection(object.Section)
		Serialize() ([]byte, error)
	}
	var f objectFile
	switch tri.os {
	case "linux", "freestanding":
		f = elf.NewFile(tgt)
	case "darwin":
		f = macho.NewFile(tgt)
	case "windows":
		f = coff.NewFile(tgt)
	default:
		return nil, fmt.Errorf("unsupported OS: %s", tri.os)
	}
	for _, s := range sections {
		f.AddSection(s)
	}
	return f.Serialize()
}

func linkObject(tri triple, objBytes []byte, dynLibs []resolvedLib, crt crtObjects, runtimeObj []byte) ([]byte, error) {
	objName := "main.o"
	if tri.os == "windows" {
		objName = "main.obj"
	}

	switch tri.os {
	case "linux":
		l := linkerelf.NewLinker(tri.elfArch())
		if err := l.AddObject("crt1.o", crt.crt1); err != nil {
			return nil, fmt.Errorf("add crt1.o: %w", err)
		}
		if err := l.AddObject("crti.o", crt.crti); err != nil {
			return nil, fmt.Errorf("add crti.o: %w", err)
		}
		if err := l.AddObject(objName, objBytes); err != nil {
			return nil, err
		}
		// Inject the runtime object right after the main object
		if len(runtimeObj) > 0 {
			if err := l.AddObject("runtime.o", runtimeObj); err != nil {
				return nil, fmt.Errorf("add runtime.o: %w", err)
			}
		}
		if err := l.AddObject("crtn.o", crt.crtn); err != nil {
			return nil, fmt.Errorf("add crtn.o: %w", err)
		}
		for _, lib := range dynLibs {
			if err := l.AddDynamicLibrary(lib.name, lib.bytes); err != nil {
				return nil, fmt.Errorf("add dynamic library %s: %w", lib.name, err)
			}
		}
		return l.Link()

	case "darwin":
		l := linkermacho.NewLinker(tri.machoArch())
		if err := l.AddObject(objName, objBytes); err != nil {
			return nil, err
		}
		if len(runtimeObj) > 0 {
			if err := l.AddObject("runtime.o", runtimeObj); err != nil {
				return nil, fmt.Errorf("add runtime.o: %w", err)
			}
		}
		for _, lib := range dynLibs {
			if err := l.AddDynamicLibrary(lib.name, lib.bytes); err != nil {
				return nil, fmt.Errorf("add dynamic library %s: %w", lib.name, err)
			}
		}
		return l.Link()

	case "windows":
		l := linkerpe.NewLinker(tri.peArch())
		if err := l.AddObject(objName, objBytes); err != nil {
			return nil, err
		}
		if len(runtimeObj) > 0 {
			if err := l.AddObject("runtime.obj", runtimeObj); err != nil {
				return nil, fmt.Errorf("add runtime.obj: %w", err)
			}
		}
		for _, lib := range dynLibs {
			if err := l.AddDynamicLibrary(lib.name, lib.bytes); err != nil {
				return nil, fmt.Errorf("add dynamic library %s: %w", lib.name, err)
			}
		}
		return l.Link()
	}

	return nil, fmt.Errorf("linking not supported for OS: %s", tri.os)
}

func boolToCode(failed bool) int {
	if failed {
		return 1
	}
	return 0
}