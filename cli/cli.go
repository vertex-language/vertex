// cli/cli.go
//
// Package cli is the argument-parsing layer for the `vertex` command. It
// parses, validates, and translates flags into a driver.Options — and does
// nothing else. Every decision about what a build *means* (which target,
// which packages, which lowering, which vvm entry point) belongs to driver.
//
// The command word is optional: `vertex build main.vs` and `vertex main.vs`
// are the same invocation, so the 0.4.0 spelling keeps working while the
// subcommand form matches vvm's own CLI.
package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vertex-language/vertex/driver"
)

func Main(args []string) int {
	if len(args) == 0 {
		printUsage(os.Stderr)
		return 2
	}

	cmd, rest := splitCommand(args)
	switch cmd {
	case "build":
		return cmdBuild(rest)
	case "run":
		return cmdRun(rest)
	case "test":
		return cmdTest(rest)
	case "targets":
		return cmdTargets(rest)
	case "version":
		fmt.Printf("vertex %s (spec %s)\n", driver.Version, driver.SpecVersion)
		return 0
	case "help":
		printUsage(os.Stdout)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "vertex: unknown command %q\n\n", cmd)
		printUsage(os.Stderr)
		return 2
	}
}

// splitCommand recognizes a leading command word. Anything else — a file, a
// directory, a flag — means the implicit "build", which is what the 0.4.0
// CLI did and what most invocations still want.
func splitCommand(args []string) (cmd string, rest []string) {
	switch args[0] {
	case "build", "run", "test", "targets", "version", "help":
		return args[0], args[1:]
	case "-h", "--help":
		return "help", args[1:]
	case "-v", "-version", "--version":
		return "version", args[1:]
	}
	return "build", args
}

// splitArgs separates positional arguments from flags so a file may be
// written before or after a flag — `vertex main.vs -o app` and
// `vertex -o app main.vs` both work. Go's flag package stops at the first
// non-flag token, which is why this pre-pass exists (same shape as
// cmd/vvm's own splitArgs, for the same reason).
//
// A bare "--" ends flag parsing entirely: everything after it is a
// positional, which is how `vertex run app.vs -- --verbose` passes its own
// flags through to the compiled program.
func splitArgs(args []string) (positionals, flags, passthrough []string) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "" {
			continue
		}
		if a == "--" {
			passthrough = append(passthrough, args[i+1:]...)
			return positionals, flags, passthrough
		}
		if a[0] == '-' {
			flags = append(flags, a)
			if isValueFlag(a) && i+1 < len(args) {
				i++
				flags = append(flags, args[i])
			}
			continue
		}
		positionals = append(positionals, a)
	}
	return positionals, flags, nil
}

// isValueFlag names the flags that consume a following token. The boolean
// flags are deliberately absent: swallowing the next token would eat a
// source path.
func isValueFlag(a string) bool {
	switch strings.TrimLeft(a, "-") {
	case "o", "target", "min-os-version", "root", "packages-dir",
		"flat-base", "dir", "file", "run":
		// A flag written as -o=app already carries its value; only the
		// separated spelling reaches here, since "-o=app" doesn't match.
		return !strings.Contains(a, "=")
	}
	return false
}

// rejectRetiredFlags turns a flag the 0.4.0 CLI accepted but this toolchain
// cannot honor into a specific error naming *why*, rather than flag's
// generic "flag provided but not defined". A silently-ignored -O2 would be
// worse than either.
func rejectRetiredFlags(flags []string) error {
	retired := map[string]string{
		"emit-mir": "there is no Machine IR stage — vvm lowers vir straight to machine code, and its top-level package exposes no MIR surface",
		"emit-asm": "vvm's assembly printers are debug-only and live below its public API; there is no supported path to a .s file",
		"emit-obj": "vvm's object-file stage (toObjectBytes) is internal to its build pipeline — it emits linked images or nothing",
		"c":        "same as -emit-obj: vvm has no public relocatable-object entry point",
		"dump":     "the pipeline-dump format described in the 0.4.0 README was never wired to this pipeline",
		"g":        "vvm's --debug is device-only (amdtx .file/.loc); the host pipeline emits no debug info yet",
		"sysroot":  "link dependencies resolve through vvm's own linkdeps search paths, which take no sysroot",
		"O0":       "there is no optimizer in either tree; -O is not accepted rather than silently ignored",
		"O1":       "there is no optimizer in either tree; -O is not accepted rather than silently ignored",
		"O2":       "there is no optimizer in either tree; -O is not accepted rather than silently ignored",
		"Os":       "there is no optimizer in either tree; -O is not accepted rather than silently ignored",
	}
	for _, f := range flags {
		name := strings.SplitN(strings.TrimLeft(f, "-"), "=", 2)[0]
		if why, ok := retired[name]; ok {
			return fmt.Errorf("-%s is not supported: %s", name, why)
		}
	}
	return nil
}

func fail(err error) int {
	fmt.Fprintf(os.Stderr, "vertex: %v\n", err)
	return 1
}

// --- build -----------------------------------------------------------------

func cmdBuild(args []string) int {
	positionals, flags, _ := splitArgs(args)
	if err := rejectRetiredFlags(flags); err != nil {
		return fail(err)
	}

	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	var (
		output       string
		target       string
		minOS        string
		root         string
		pkgDir       string
		flatBase     uint64
		emitVIR      bool
		emitVByte    bool
		shared       bool
		verbose      bool
	)
	fs.StringVar(&output, "o", "", "output path")
	fs.StringVar(&target, "target", "", "target triple (default: host)")
	fs.StringVar(&minOS, "min-os-version", "", "minimum OS version; required for darwin targets")
	fs.StringVar(&root, "root", "", "root module name override for entry-point resolution")
	fs.StringVar(&pkgDir, "packages-dir", "", "Vertex packages root (overrides $VERTEX_PATH)")
	fs.Uint64Var(&flatBase, "flat-base", 0, "base address for a freestanding flat image")
	fs.BoolVar(&emitVIR, "emit-vir", false, "emit Vertex IR text (.vir) instead of a binary")
	fs.BoolVar(&emitVByte, "emit-vbyte", false, "emit Vertex IR binary (.vbyte) instead of a binary")
	fs.BoolVar(&shared, "shared", false, "build a shared library instead of an executable")
	fs.BoolVar(&verbose, "v", false, "report each pipeline stage")
	fs.Usage = func() { fmt.Fprintln(os.Stderr, "usage: vertex build [flags] <file.vs | package-dir>") }

	if err := fs.Parse(flags); err != nil {
		return 2
	}
	if len(positionals) != 1 {
		fs.Usage()
		return 2
	}
	if emitVIR && emitVByte {
		return fail(fmt.Errorf("-emit-vir and -emit-vbyte both name an output form; pick one"))
	}

	emit := driver.EmitBinary
	switch {
	case emitVIR:
		emit = driver.EmitVIR
	case emitVByte:
		emit = driver.EmitVByte
	}

	opts := &driver.Options{
		Input:           positionals[0],
		Target:          target,
		MinOSVersion:    minOS,
		Output:          output,
		Emit:            emit,
		Shared:          shared,
		FlatBaseAddress: flatBase,
		RootModule:      root,
		PackagesDir:     pkgDir,
		Verbose:         verbose,
	}

	res, err := driver.Compile(opts)
	if err != nil {
		return fail(err)
	}
	for _, out := range res.Outputs {
		fmt.Fprintf(os.Stderr, "vertex: wrote %s (%s)\n", out, res.Target.Name)
	}
	return 0
}

// --- run -------------------------------------------------------------------

func cmdRun(args []string) int {
	positionals, flags, passthrough := splitArgs(args)
	if err := rejectRetiredFlags(flags); err != nil {
		return fail(err)
	}

	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	var (
		pkgDir  string
		verbose bool
	)
	fs.StringVar(&pkgDir, "packages-dir", "", "Vertex packages root (overrides $VERTEX_PATH)")
	fs.BoolVar(&verbose, "v", false, "report each pipeline stage")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: vertex run [flags] <file.vs | package-dir> [-- args...]")
	}
	if err := fs.Parse(flags); err != nil {
		return 2
	}
	if len(positionals) != 1 {
		fs.Usage()
		return 2
	}

	opts := &driver.Options{
		Input:       positionals[0],
		PackagesDir: pkgDir,
		Verbose:     verbose,
	}
	code, err := driver.RunProgram(opts, passthrough)
	if err != nil {
		return fail(err)
	}
	return code
}

// --- test ------------------------------------------------------------------

func cmdTest(args []string) int {
	positionals, flags, _ := splitArgs(args)
	if err := rejectRetiredFlags(flags); err != nil {
		return fail(err)
	}

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var (
		dir     string
		file    string
		filter  string
		pkgDir  string
		verbose bool
	)
	fs.StringVar(&dir, "dir", "", "directory holding `build test` files (default: .)")
	fs.StringVar(&file, "file", "", "a single test file")
	fs.StringVar(&filter, "run", "", "only run tests whose name contains this substring")
	fs.StringVar(&pkgDir, "packages-dir", "", "Vertex packages root (overrides $VERTEX_PATH)")
	fs.BoolVar(&verbose, "v", false, "print every test, not just failures")
	fs.Usage = func() { fmt.Fprintln(os.Stderr, "usage: vertex test [-dir <path> | -file <path>] [-run <substr>]") }
	if err := fs.Parse(flags); err != nil {
		return 2
	}
	if dir != "" && file != "" {
		return fail(fmt.Errorf("-dir and -file both name what to test; pass one"))
	}

	input := dir
	if file != "" {
		input = file
	}
	if input == "" && len(positionals) == 1 {
		input = positionals[0]
	}
	if input == "" {
		input = "."
	}

	opts := &driver.Options{
		Input:       input,
		PackagesDir: pkgDir,
		Verbose:     verbose,
	}
	ok, err := driver.RunTests(opts, filter)
	if err != nil {
		return fail(err)
	}
	if !ok {
		return 1
	}
	return 0
}

// --- targets ---------------------------------------------------------------

func cmdTargets(args []string) int {
	_, flags, _ := splitArgs(args)
	fs := flag.NewFlagSet("targets", flag.ContinueOnError)
	if err := fs.Parse(flags); err != nil {
		return 2
	}
	host, hostErr := driver.HostTargetName()
	for _, t := range driver.Targets() {
		marker := "  "
		if hostErr == nil && t.Name == host {
			marker = "* "
		}
		fmt.Printf("%s%-22s %-24s %s\n", marker, t.Name, t.VVM.String(), t.Note)
	}
	if hostErr == nil {
		fmt.Printf("\n* = this host\n")
	} else {
		fmt.Printf("\nno entry matches this host: %v\n", hostErr)
	}
	return 0
}