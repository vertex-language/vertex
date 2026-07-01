package driver

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vertex-language/ir/vertex/ast"

	virbinary "github.com/vertex-language/ir/vertex/encoding/binary"
	virtext "github.com/vertex-language/ir/vertex/encoding/text"

	mirtext "github.com/vertex-language/ir/machine/encoding/text/mir"

	codesign "github.com/vertex-language/linker/macho/codesign"

	"github.com/vertex-language/pkg"
	"github.com/vertex-language/pkg/importer"

	"github.com/vertex-language/vertex/codegen"
	"github.com/vertex-language/vertex/nativelibs"
	"github.com/vertex-language/vertex/pipeline"
	"github.com/vertex-language/vertex/target"
)

// Compile is the top-level entry point for a normal (non-test) build. It
// resolves cfg.Input's module graph — if it has one — into a pkg.Graph,
// runs every module in that graph through the pipeline in build order,
// and produces whatever cfg.Mode asks for. See doc.go for why no module
// here (including any project runtime) is ever treated specially.
func Compile(cfg Config, stderr io.Writer) int {
	tri, err := target.ParseTriple(cfg.Target)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 2
	}

	root, err := findModuleRoot(cfg.Input)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	if root == "" {
		p, err := parseInput(cfg.Input)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		if imp, ok := firstNonStdlibImport(p); ok {
			fmt.Fprintf(stderr, "vertex: %s imports %q, but no vs.mod was found in this directory (or any parent).\n\n", cfg.Input, imp)
			fmt.Fprintf(stderr, "Run `vertex mod init <module-path>` to create one, then:\n    vertex mod get %s\n", imp)
			return 1
		}
		units := []*pipeline.Unit{{IsRoot: true, Dir: cfg.Input, Pkg: p}}
		libDirs := nativelibs.SearchDirs(tri, effectiveSysroot(cfg, tri), nil)
		return compileUnits(cfg, tri, units, libDirs, stderr)
	}

	homeDir, err := pkg.Home(cfg.VertexHome)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	cache, err := pkg.OpenCache(homeDir, importer.NewGitFetcher())
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	graph, err := pkg.Load(root, cache, cfg.LoadMode)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	units := make([]*pipeline.Unit, 0, len(graph.Modules))
	for _, m := range graph.Modules {
		p, err := parseInput(m.Dir)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %s: %v\n", m.Path, err)
			return 1
		}
		units = append(units, &pipeline.Unit{
			ModulePath: string(m.Path),
			Dir:        m.Dir,
			IsRoot:     m == graph.Root,
			Pkg:        p,
		})
	}

	// hostRelease would let EnsureNativeLibs prefer a release-qualified
	// vs.lib target (e.g. "ubuntu-22.04") over a bare catch-all; probing
	// the host's actual release is toolchain-detection work that belongs
	// beside pkg/toolchain, not here. "" is always correct, just less
	// specific: only unqualified targets match.
	const hostRelease = ""
	libResults, err := graph.EnsureNativeLibs(cache, tri.Arch, tri.OS, hostRelease, nil)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	libDirs := nativelibs.SearchDirs(tri, effectiveSysroot(cfg, tri), libResults)

	return compileUnits(cfg, tri, units, libDirs, stderr)
}

// CompilePackage compiles a single, already-parsed, module-graph-free
// package — used directly by the test runner for its synthetic per-test
// packages.
func CompilePackage(p *ast.Package, cfg Config, stderr io.Writer) int {
	tri, err := target.ParseTriple(cfg.Target)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 2
	}
	units := []*pipeline.Unit{{IsRoot: true, Pkg: p}}
	libDirs := nativelibs.SearchDirs(tri, effectiveSysroot(cfg, tri), nil)
	return compileUnits(cfg, tri, units, libDirs, stderr)
}

// DumpPackage runs every pipeline stage against an already-parsed package
// and writes the annotated result to path, ignoring cfg.Output and
// cfg.Mode. Used by the test runner to capture a full pipeline dump for a
// failed test.
func DumpPackage(p *ast.Package, cfg Config, path string, stderr io.Writer) {
	dumpCfg := cfg
	dumpCfg.Mode = ModeDump
	dumpCfg.Output = path

	tri, err := target.ParseTriple(dumpCfg.Target)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: dump: %v\n", err)
		return
	}
	units := []*pipeline.Unit{{IsRoot: true, Pkg: p}}
	libDirs := nativelibs.SearchDirs(tri, effectiveSysroot(dumpCfg, tri), nil)
	compileUnits(dumpCfg, tri, units, libDirs, stderr)
}

func effectiveSysroot(cfg Config, tri target.Triple) string {
	if cfg.Sysroot != "" {
		return cfg.Sysroot
	}
	return nativelibs.AutoSysroot(tri)
}

func compileUnits(cfg Config, tri target.Triple, units []*pipeline.Unit, libDirs []string, stderr io.Writer) int {
	if cfg.Mode == ModeRun {
		return runExe(units, cfg, tri, libDirs, stderr)
	}

	st := &pipeline.State{
		Units:   units,
		Triple:  tri,
		Opts:    codegen.Options{OptLevel: cfg.OptLevel, DebugInfo: cfg.DebugInfo},
		LibDirs: libDirs,
	}
	root := st.Root()

	if cfg.Mode == ModeDump {
		var sb strings.Builder
		if err := appendSourceFiles(&sb, root); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		st.Sink = &sb
		runErr := pipeline.Run(pipeline.Stages, st, "") // run every stage; failures are already visible in the banners
		if err := writeOutput(cfg.Output, []byte(sb.String())); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		if runErr != nil || root.VIRErr != nil {
			return 1
		}
		return 0
	}

	if err := pipeline.Run(pipeline.Stages, st, stageFor(cfg.Mode)); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	switch cfg.Mode {
	case ModeVIR:
		if err := writeOutput(cfg.Output, []byte(virtext.Format(root.VIR))); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return boolToCode(root.VIRErr != nil)

	case ModeVBytes:
		data, err := virbinary.Marshal(root.VIR)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: vbytes encoding: %v\n", err)
			return 1
		}
		if err := writeOutput(cfg.Output, data); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return boolToCode(root.VIRErr != nil)

	case ModeMIR:
		if err := writeOutput(cfg.Output, []byte(mirtext.PrintModule(root.MIR))); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0

	case ModeASM:
		if err := writeOutput(cfg.Output, []byte(root.ASM)); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0

	case ModeObj:
		if err := writeOutput(cfg.Output, root.Obj); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0

	case ModeExe:
		exeBytes := st.Exe
		if tri.OS == "darwin" {
			id := stripExt(filepath.Base(cfg.Output))
			var signErr error
			exeBytes, signErr = codesign.SignImage(exeBytes, codesign.Options{Identifier: id})
			if signErr != nil {
				fmt.Fprintf(stderr, "vertex: codesign: %v\n", signErr)
				return 1
			}
		}
		if err := writeExe(cfg.Output, exeBytes); err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		return 0
	}

	fmt.Fprintf(stderr, "vertex: internal error: unhandled mode %v\n", cfg.Mode)
	return 1
}

func stageFor(mode EmitMode) string {
	switch mode {
	case ModeVIR, ModeVBytes:
		return pipeline.StageVIR
	case ModeMIR:
		return pipeline.StageMIR
	case ModeASM:
		return pipeline.StageASM
	case ModeObj:
		return pipeline.StageObject
	default:
		return pipeline.StageLink
	}
}

// runExe compiles units to a temporary executable, runs it with the
// process's own stdin/stdout/stderr attached, cleans up the temp
// directory, and returns the child's exit code. It is the implementation
// of ModeRun.
func runExe(units []*pipeline.Unit, cfg Config, tri target.Triple, libDirs []string, stderr io.Writer) int {
	tmpDir, err := os.MkdirTemp("", "vertex-run-*")
	if err != nil {
		fmt.Fprintf(stderr, "vertex: cannot create temp dir: %v\n", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	baseName := "main"
	if cfg.Input != "" {
		baseName = stripExt(filepath.Base(cfg.Input))
	}
	binPath := filepath.Join(tmpDir, baseName)
	if target.IsWindowsTarget(cfg.Target) {
		binPath += ".exe"
	}

	runCfg := cfg
	runCfg.Mode = ModeExe
	runCfg.Output = binPath

	if code := compileUnits(runCfg, tri, units, libDirs, stderr); code != 0 {
		return code
	}

	cmd := exec.Command(binPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(stderr, "vertex: run: %v\n", err)
		return 1
	}
	return 0
}