package driver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

func emit(cfg config, stderr io.Writer) int {
	if cfg.mode == modeDump {
		return dumpAll(cfg, stderr)
	}

	tri, err := parseTriple(cfg.target)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 2
	}

	// ── Stage 1: Parse Vertex source (.vs) → AST ─────────────────────────────
	pkg, err := parseInput(cfg.input)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
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

	// ── Stage 5: Build object file ────────────────────────────────────────────
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

	// ── Stage 7: Link to native executable ───────────────────────────────────
	exeBytes, err := linkObject(tri, objBytes, dynLibs, crt)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: link: %v\n", err)
		return 1
	}
	if err := writeOutput(cfg.output, exeBytes); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	return 0
}

// parseInput reads one .vs file or all .vs files in a directory.
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

// linkObject links objBytes against the provided dynamic libraries and CRT
// objects, returning the final executable bytes.
//
// The canonical ELF link order is: crt1.o, crti.o, <user object>, crtn.o,
// followed by dynamic libraries. crtObjects is empty for non-Linux targets
// where the system linker or runtime handles startup transparently.
func linkObject(tri triple, objBytes []byte, dynLibs []resolvedLib, crt crtObjects) ([]byte, error) {
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