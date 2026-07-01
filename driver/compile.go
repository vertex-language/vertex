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
	"github.com/vertex-language/pkg/parser/mod"

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

		imports := collectNonStdlibImports(p)
		if len(imports) == 0 {
			units := []*pipeline.Unit{{IsRoot: true, Dir: cfg.Input, Pkg: p}}
			libDirs := nativelibs.SearchDirs(tri, effectiveSysroot(cfg, tri), nil)
			return compileUnits(cfg, tri, units, libDirs, stderr)
		}

		// No vs.mod anywhere above cfg.Input, but the file itself names
		// real imports: resolve them directly rather than erroring out,
		// the same way `go run somefile.go` resolves its imports against
		// the module cache without a go.mod on disk anywhere.
		homeDir, err := pkg.Home(cfg.VertexHome)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		fetcher := importer.NewGitFetcher()
		cache, err := pkg.OpenCache(homeDir, fetcher)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}
		srcDir, err := sourceDir(cfg.Input)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n", err)
			return 1
		}

		fmt.Fprintf(stderr, "vertex: %s: no vs.mod found; resolving %s directly\n", cfg.Input, strings.Join(imports, ", "))
		graph, err := loadGraphFromImports(srcDir, imports, cache, fetcher)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %v\n\nrun `vertex mod init <module-path>` and `vertex mod get` to pin these explicitly instead.\n", err)
			return 1
		}

		return compileGraph(cfg, tri, graph, cache, cfg.Input, p, stderr)
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

	return compileGraph(cfg, tri, graph, cache, "", nil, stderr)
}

// compileGraph turns a resolved *pkg.Graph into pipeline.Units, ensures
// native libraries, and runs the pipeline. rootPkg, if non-nil, is used
// directly as the graph's root unit's package instead of re-parsing
// graph.Root.Dir off disk — what a single-file compile needs, since
// cfg.Input may name one file inside a directory that also holds others
// the compile was never asked to include; rootUnitDir is that unit's Dir
// in that case. Both are ignored when rootPkg is nil (the normal,
// vs.mod-rooted case), where every unit — root included — is parsed from
// its own graph.Module.Dir.
func compileGraph(cfg Config, tri target.Triple, graph *pkg.Graph, cache *pkg.Cache, rootUnitDir string, rootPkg *ast.Package, stderr io.Writer) int {
	units := make([]*pipeline.Unit, 0, len(graph.Modules))
	for _, m := range graph.Modules {
		if m == graph.Root && rootPkg != nil {
			units = append(units, &pipeline.Unit{Dir: rootUnitDir, IsRoot: true, Pkg: rootPkg})
			continue
		}
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

// loadGraphFromImports resolves a dependency graph for a single source
// file (or directory) that has no vs.mod anywhere above it, but does
// import one or more non-stdlib packages. Each import is resolved to a
// concrete version directly through fetcher — the same call cli's `mod
// get` makes — rather than being read as a requirement off disk; the
// result is assembled into a synthetic *mod.File (module path
// "command-line-arguments", the same sentinel Go's own tooling uses for
// an ad hoc, manifest-free compile unit) and handed to pkg.LoadModule,
// which walks it exactly as it would a committed vs.mod.
//
// The scratch vs.sum this creates (and removes before returning) exists
// only to satisfy pkg.Cache.Mod's verify/record bookkeeping for this one
// call — nothing is written into rootDir, and nothing here is meant to
// be reproducible build-to-build the way a committed vs.sum is: every
// call re-resolves "latest" for itself, the same as `go run` with no
// go.mod.
func loadGraphFromImports(rootDir string, imports []string, cache *pkg.Cache, fetcher importer.Fetcher) (*pkg.Graph, error) {
	deps := make([]*mod.Dependency, 0, len(imports))
	for _, imp := range imports {
		version, err := fetcher.Resolve(mod.ModulePath(imp), "latest")
		if err != nil {
			return nil, fmt.Errorf("resolving %s: %w", imp, err)
		}
		deps = append(deps, &mod.Dependency{Mod: mod.ModuleVersion{Path: mod.ModulePath(imp), Version: version}})
	}

	rootMF := &mod.File{
		Module:       &mod.Module{Path: "command-line-arguments"},
		Dependencies: deps,
	}

	sumFile, err := os.CreateTemp("", "vertex-sum-*")
	if err != nil {
		return nil, fmt.Errorf("create scratch vs.sum: %w", err)
	}
	sumPath := sumFile.Name()
	sumFile.Close()
	defer os.Remove(sumPath)

	// ModUpdate: there is no committed vs.sum here for -mod=readonly to
	// protect, so cfg.LoadMode doesn't apply to this resolution — same
	// reasoning `go run pkg@version` fetches unconditionally rather than
	// refusing for lack of a local go.sum.
	return pkg.LoadModule(rootDir, sumPath, rootMF, cache, pkg.ModUpdate)
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