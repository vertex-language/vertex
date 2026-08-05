// driver/driver.go
//
// Package driver owns every decision about what a `vertex` invocation
// means. It is the one place allowed to import both halves of the
// toolchain at once — the frontend (parser, importer, analyzer, types),
// the lowering chain (lower/hir, lower/vir), and vvm — and to pick the
// right combination for a given target and output form.
//
// Nothing below this package knows it exists, exactly as vvm's own
// top-level package relates to ir/vir, cpu/lower, and linker/*. cli/
// parses flags and calls in here; it makes no pipeline choices of its own.
//
// The pipeline, once:
//
//	<file.vs | dir>
//	     │  load.go: parser + analyzer (via importer for a package,
//	     │           directly for a single file)
//	     ▼
//	*driver.Loaded  (one FileSet, checked packages in topological order,
//	     │           every diagnostic collected)
//	     │  lower.go: hir.Lower  — every decision made here
//	     ▼
//	*hir.Program
//	     │  lower.go: lower/vir.Lower — mechanical
//	     ▼
//	[]*vir.Module  (unverified; ir/verify is vvm's to run)
//	     │
//	     ├─ EmitVIR / EmitVByte ── format/vbyte/{text,binary}.Encode ─► files
//	     └─ EmitBinary ─────────── vvm.BuildModule / BuildModuleGraph ─► image
//
// One rule governs the whole file set: this package never reimplements a
// stage it can call. It does not verify (vvm does), does not link (vvm
// does), does not re-derive a target's container format (vvm.Target does).
package driver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	virmod "github.com/vertex-language/vvm/ir/vir"
)

const (
	// Version is this compiler's version; SpecVersion is the language
	// spec revision it implements. They move independently.
	Version     = "0.4.0"
	SpecVersion = "2.2"
)

// EmitKind selects what a build produces. There are exactly three, because
// vvm's public surface offers exactly three: the two Vertex IR
// serializations and a finished native image.
type EmitKind int

const (
	EmitBinary EmitKind = iota
	EmitVIR
	EmitVByte
)

func (e EmitKind) String() string {
	switch e {
	case EmitVIR:
		return "vir"
	case EmitVByte:
		return "vbyte"
	default:
		return "binary"
	}
}

// ext is the conventional file extension for an IR emission. EmitBinary
// has none — a native image's name is the program's name.
func (e EmitKind) ext() string {
	switch e {
	case EmitVIR:
		return ".vir"
	case EmitVByte:
		return ".vbyte"
	default:
		return ""
	}
}

// Options is one invocation's full configuration. Every zero value is the
// "ordinary host build" default.
type Options struct {
	// Input is a .vs file or a directory holding one package.
	Input string

	// Target names an entry in the target table (see target.go). Empty
	// means the host.
	Target string

	// MinOSVersion is required by vvm for every Mach-O target. Empty
	// means DefaultMinOSVersion for a darwin target, and is ignored
	// everywhere else.
	MinOSVersion string

	// Output is the file (or, for a multi-package IR emission, the
	// directory) to write. Empty derives a name from Input.
	Output string

	Emit            EmitKind
	Shared          bool
	FlatBaseAddress uint64

	// RootModule overrides which module's entry function becomes the
	// process entry point. Empty uses lower/vir.Root, which already knows
	// which package holds main — this exists for the case where a caller
	// wants a different one, matching vvm's own --root.
	RootModule string

	// PackagesDir overrides $VERTEX_PATH.
	PackagesDir string

	Verbose bool

	// Stdout/Stderr default to the process's own. They exist so a test
	// can capture output without touching globals.
	Stdout io.Writer
	Stderr io.Writer
}

func (o *Options) defaults() {
	if o.Stdout == nil {
		o.Stdout = os.Stdout
	}
	if o.Stderr == nil {
		o.Stderr = os.Stderr
	}
}

func (o *Options) logf(format string, args ...any) {
	if !o.Verbose {
		return
	}
	fmt.Fprintf(o.Stderr, "vertex: "+format+"\n", args...)
}

// Result is what a completed Compile hands back. The modules and the
// checked packages are both kept, not just the bytes, so a caller (the
// test runner, a future language server) can inspect what was lowered
// without re-running the pipeline.
type Result struct {
	Target   Target
	Root     string
	Packages []*Package
	Modules  []*virmod.Module
	Outputs  []string
}

// Compile runs the full pipeline and writes whatever Options.Emit selects.
func Compile(opts *Options) (*Result, error) {
	opts.defaults()

	t, err := ResolveTarget(targetRequest{
		Name:         opts.Target,
		MinOSVersion: opts.MinOSVersion,
		Shared:       opts.Shared,
		FlatBase:     opts.FlatBaseAddress,
	})
	if err != nil {
		return nil, err
	}
	opts.logf("target %s -> %s (tag %s)", t.Name, t.VVM, t.Tag)

	ld, err := Load(opts, t.Tag)
	if err != nil {
		return nil, err
	}
	opts.logf("checked %d package(s): %s",
		len(ld.Packages), strings.Join(sortedPaths(ld.Packages), ", "))

	modules, root, err := Lower(opts, t, ld, LowerOptions{})
	if err != nil {
		return nil, err
	}
	opts.logf("lowered %d module(s), root %q", len(modules), root)

	if opts.RootModule != "" {
		root = opts.RootModule
	}

	res := &Result{Target: t, Root: root, Packages: ld.Packages, Modules: modules}
	if err := emit(opts, res); err != nil {
		return nil, err
	}
	return res, nil
}

// RunProgram builds for the host and executes the result, forwarding args
// and the exit code.
//
// This deliberately does not call vvm.Run/RunModule. Those buffer the
// child's stdout and stderr and launch it with no arguments, which is the
// right shape for `vvm run` as a pipeline smoke test and the wrong one for
// a language's run command: a Vertex program should stream its output and
// receive its own argv. run.go's execute() does that instead, over the
// same bytes vvm.BuildModule produced.
func RunProgram(opts *Options, args []string) (int, error) {
	opts.defaults()

	if opts.Target != "" {
		host, err := HostTargetName()
		if err != nil {
			return 1, err
		}
		if opts.Target != host {
			return 1, fmt.Errorf(
				"run builds for the host (%s) and executes it — cannot run a %s binary here; "+
					"use `vertex build -target %s` and run it on that machine",
				host, opts.Target, opts.Target)
		}
	}

	// Force the shape run implies, whatever the caller set: a native
	// executable for this machine.
	ropts := *opts
	ropts.Target = ""
	ropts.Emit = EmitBinary
	ropts.Shared = false
	ropts.Output = ""

	t, err := ResolveTarget(targetRequest{MinOSVersion: opts.MinOSVersion})
	if err != nil {
		return 1, err
	}
	ld, err := Load(&ropts, t.Tag)
	if err != nil {
		return 1, err
	}
	modules, root, err := Lower(&ropts, t, ld, LowerOptions{})
	if err != nil {
		return 1, err
	}
	if ropts.RootModule != "" {
		root = ropts.RootModule
	}
	bin, err := buildImage(t, modules, root)
	if err != nil {
		return 1, err
	}
	return execute(&ropts, bin, args)
}

// outputBase derives the default output name from the input path: a file's
// base without .vs, or a directory's own name. A trailing separator on a
// directory is stripped first, so "./cmd/app/" and "./cmd/app" agree.
func outputBase(input string) string {
	clean := filepath.Clean(input)
	base := filepath.Base(clean)
	return strings.TrimSuffix(base, ".vs")
}