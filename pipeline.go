package main

import (
	"fmt"
	"io"

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
	tri, err := parseTriple(cfg.target)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 2
	}

	// ── Stage 1: Parse Vertex source (.vs) → AST ─────────────────────────────
	pkg, err := ast.NewPackageFromPath(cfg.input, ast.ParseOptions{
		PackagesDir: cfg.packagesDir,
		BuildTags:   tri.buildTags(),
	})
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	hadFrontendError := false
	for _, d := range pkg.Diagnostics() {
		fmt.Fprintln(stderr, d)
		if d.Severity == ast.SevError {
			hadFrontendError = true
		}
	}
	if hadFrontendError {
		return 1
	}

	// ── Stage 2: Lower AST → Vertex IR ───────────────────────────────────────
	virMod, virErr := virlower.NewLower(pkg, nil)
	if virErr != nil {
		// NewLower always returns a non-nil *vertex.Module; the error carries
		// per-node diagnostics. Report and continue so partial VIR artifacts
		// remain useful (e.g. for a language server).
		fmt.Fprintln(stderr, virErr)
	}

	// Set the target on the module before handing it downstream; lowering
	// and encoding both inspect it to resolve arch-specific behaviour.
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

	// All stages below require clean VIR.
	if virErr != nil {
		return 1
	}

	// ── Stage 3: Lower Vertex IR → Machine IR (SSA) ───────────────────────────
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

	// ── Stage 6: Link to native executable ───────────────────────────────────
	if tri.os == "freestanding" {
		fmt.Fprintf(stderr, "vertex: cannot link a freestanding target; use -c/-emit-obj instead\n")
		return 2
	}

	exeBytes, err := linkObject(tri, objBytes)
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

// ── helpers ───────────────────────────────────────────────────────────────────

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

func linkObject(tri triple, objBytes []byte) ([]byte, error) {
	objName := "main.o"
	if tri.os == "windows" {
		objName = "main.obj"
	}

	switch tri.os {
	case "linux":
		l := linkerelf.NewLinker(tri.elfArch())
		if err := l.AddObject(objName, objBytes); err != nil {
			return nil, err
		}
		return l.Link()

	case "darwin":
		l := linkermacho.NewLinker(tri.machoArch())
		if err := l.AddObject(objName, objBytes); err != nil {
			return nil, err
		}
		return l.Link()

	case "windows":
		l := linkerpe.NewLinker(tri.peArch())
		if err := l.AddObject(objName, objBytes); err != nil {
			return nil, err
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