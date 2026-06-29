// cli.go
package driver

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const version = "0.4.0"

type emitMode uint8

const (
	modeVIR    emitMode = iota
	modeVBytes
	modeMIR
	modeASM
	modeObj
	modeExe
	modeDump
	modeTest
	modeRun // compile to temp binary and execute; triggered when no -o and no emit flag
)

type config struct {
	input       string
	output      string
	target      string
	packagesDir string
	sysroot     string
	mode        emitMode
	optLevel    int
	debugInfo   bool
	testDir     string
	testFile    string
}

func parseFlags(args []string, stderr io.Writer) (config, int) {
	fs := flag.NewFlagSet("vertex", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		cfg config

		fVIR    bool
		fVBytes bool
		fMIR    bool
		fASM    bool
		fObj    bool
		fC      bool
		fDump   bool
		fTest   bool

		fO0, fO1, fO2, fOs bool
		printVer            bool
		listTargets         bool
		printTarget         bool

		outputExplicit bool
	)

	fs.BoolVar(&fVIR,    "emit-vir",    false, "emit Vertex IR text (.vir)")
	fs.BoolVar(&fVBytes, "emit-vbytes", false, "emit Vertex IR binary (.vbytes)")
	fs.BoolVar(&fMIR,    "emit-mir",    false, "emit Machine IR text (.mir)")
	fs.BoolVar(&fASM,    "emit-asm",    false, "emit native assembly text (.s)")
	fs.BoolVar(&fObj,    "emit-obj",    false, "emit relocatable object file (.o / .obj)")
	fs.BoolVar(&fC,      "c",           false, "compile to object file (alias for -emit-obj)")
	fs.BoolVar(&fDump,   "dump",        false, "dump all pipeline stages to a single annotated file (.dump)")
	fs.BoolVar(&fDump,   "dump-all",    false, "alias for -dump")
	fs.BoolVar(&fTest,   "test",        false, "discover and run test functions")
	fs.BoolVar(&listTargets, "list-targets", false, "list all supported targets and exit")
	fs.BoolVar(&printTarget, "print-target", false, "print the effective target triple and exit")

	fs.BoolVar(&fO0, "O0", false, "disable optimisation (default)")
	fs.BoolVar(&fO1, "O1", false, "light optimisation")
	fs.BoolVar(&fO2, "O2", false, "full optimisation")
	fs.BoolVar(&fOs, "Os", false, "optimise for size")

	fs.Func("o", "write output to `file`", func(s string) error {
		cfg.output = s
		outputExplicit = true
		return nil
	})
	fs.StringVar(&cfg.target,      "target",        defaultTarget(),      "target triple (os-arch)")
	fs.StringVar(&cfg.packagesDir, "packages-dir",  defaultPackagesDir(), "Vertex packages root (overrides $VERTEX_PATH)")
	fs.StringVar(&cfg.sysroot,     "sysroot",       "",                   "sysroot for cross-compilation library search")
	fs.StringVar(&cfg.testDir,     "dir",           "",                   "directory to search for test files (recursive)")
	fs.StringVar(&cfg.testFile,    "file",          "",                   "single test file to run")
	fs.BoolVar(&cfg.debugInfo,     "g",             false,                "include debug information")
	fs.BoolVar(&printVer, "version", false, "print version and exit")
	fs.BoolVar(&printVer, "v",       false, "alias for -version")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Vertex compiler %s\n\n", version)
		fmt.Fprintf(stderr, "Usage:\n  vertex [flags] <source.vs | package/>\n\n")
		fmt.Fprintf(stderr, "Emit mode (default: compile, link, and run as a temporary executable):\n")
		fmt.Fprintf(stderr, "  -emit-vir             emit Vertex IR text (.vir)\n")
		fmt.Fprintf(stderr, "  -emit-vbytes          emit Vertex IR binary (.vbytes)\n")
		fmt.Fprintf(stderr, "  -emit-mir             emit Machine IR text (.mir)\n")
		fmt.Fprintf(stderr, "  -emit-asm             emit native assembly text (.s)\n")
		fmt.Fprintf(stderr, "  -emit-obj, -c         emit relocatable object file (.o / .obj)\n")
		fmt.Fprintf(stderr, "  -dump, -dump-all      dump all pipeline stages (.dump)\n")
		fmt.Fprintf(stderr, "  -test                 discover and run test functions\n\n")
		fmt.Fprintf(stderr, "Test options:\n")
		fmt.Fprintf(stderr, "  -dir  <path>     directory to search recursively (default: .)\n")
		fmt.Fprintf(stderr, "  -file <path>     single test file\n\n")
		fmt.Fprintf(stderr, "Options:\n")
		fmt.Fprintf(stderr, "  -o <file>        output file; when given, writes a permanent binary (no auto-run)\n")
		fmt.Fprintf(stderr, "  -target <triple> target triple (default: %s)\n", defaultTarget())
		fmt.Fprintf(stderr, "  -list-targets    list all supported targets and exit\n")
		fmt.Fprintf(stderr, "  -print-target    print the effective target triple and exit\n")
		fmt.Fprintf(stderr, "  -sysroot <path>  sysroot for cross-compilation library search\n")
		fmt.Fprintf(stderr, "  -packages-dir    Vertex packages root (overrides $VERTEX_PATH)\n")
		fmt.Fprintf(stderr, "  -O0/-O1/-O2/-Os  optimisation level (default: -O0)\n")
		fmt.Fprintf(stderr, "  -g               include debug information\n")
		fmt.Fprintf(stderr, "  -v, -version     print version and exit\n\n")
		fmt.Fprintf(stderr, "Examples:\n")
		fmt.Fprintf(stderr, "  vertex                main.vs                        (compile + run)\n")
		fmt.Fprintf(stderr, "  vertex -o main        main.vs                        (compile to binary)\n")
		fmt.Fprintf(stderr, "  vertex -o main        -target darwin-arm64 -O2 main.vs\n")
		fmt.Fprintf(stderr, "  vertex -c           -o main.o      main.vs\n")
		fmt.Fprintf(stderr, "  vertex -emit-asm    -o main.s      main.vs\n")
		fmt.Fprintf(stderr, "  vertex -emit-mir    -o main.mir    main.vs\n")
		fmt.Fprintf(stderr, "  vertex -emit-vir    -o main.vir    main.vs\n")
		fmt.Fprintf(stderr, "  vertex -emit-vbytes -o main.vbytes main.vs\n")
		fmt.Fprintf(stderr, "  vertex -dump        -o main.dump   main.vs\n")
		fmt.Fprintf(stderr, "  vertex -dump        -o -           main.vs\n")
		fmt.Fprintf(stderr, "  vertex -test\n")
		fmt.Fprintf(stderr, "  vertex -test -dir ./tests\n")
		fmt.Fprintf(stderr, "  vertex -test -file literals_test.vs\n")
	}

	if err := fs.Parse(expandShortFlag(args)); err != nil {
		return config{}, 2
	}

	if printVer {
		fmt.Fprintf(os.Stdout, "vertex %s\n", version)
		return config{}, 0
	}

	if listTargets {
		fmt.Fprintf(os.Stdout, "Supported targets:\n")
		for _, t := range supportedTargets() {
			marker := ""
			if t == defaultTarget() {
				marker = "  (default)"
			}
			fmt.Fprintf(os.Stdout, "  %s%s\n", t, marker)
		}
		return config{}, 0
	}

	if printTarget {
		fmt.Fprintf(os.Stdout, "%s\n", cfg.target)
		return config{}, 0
	}

	// ── test mode ────────────────────────────────────────────────────────────
	if fTest {
		cfg.mode = modeTest

		if cfg.testFile != "" && cfg.testDir != "" {
			fmt.Fprintf(stderr, "vertex: -test: -file and -dir are mutually exclusive\n")
			return config{}, 2
		}

		if cfg.testFile != "" || cfg.testDir != "" {
			if fs.NArg() != 0 {
				fmt.Fprintf(stderr, "vertex: -test: unexpected positional argument when -file or -dir is set\n")
				return config{}, 2
			}
		} else {
			switch fs.NArg() {
			case 0:
				cfg.testDir = "."
			case 1:
				cfg.testDir = fs.Arg(0)
			default:
				fmt.Fprintf(stderr, "vertex: -test: expected at most one positional argument\n")
				return config{}, 2
			}
		}

		switch {
		case fOs:
			cfg.optLevel = -1
		case fO2:
			cfg.optLevel = 2
		case fO1:
			cfg.optLevel = 1
		default:
			cfg.optLevel = 0
		}

		return cfg, -1
	}

	// ── emit mode ─────────────────────────────────────────────────────────────
	modes := []bool{fVIR, fVBytes, fMIR, fASM, fObj || fC, fDump}
	count := 0
	for _, b := range modes {
		if b {
			count++
		}
	}
	switch count {
	case 0:
		if outputExplicit {
			cfg.mode = modeExe
		} else {
			cfg.mode = modeRun
		}
	case 1:
	default:
		fmt.Fprintf(stderr, "vertex: emit modes are mutually exclusive\n")
		return config{}, 2
	}

	switch {
	case fVIR:
		cfg.mode = modeVIR
	case fVBytes:
		cfg.mode = modeVBytes
	case fMIR:
		cfg.mode = modeMIR
	case fASM:
		cfg.mode = modeASM
	case fObj || fC:
		cfg.mode = modeObj
	case fDump:
		cfg.mode = modeDump
	}

	switch {
	case fOs:
		cfg.optLevel = -1
	case fO2:
		cfg.optLevel = 2
	case fO1:
		cfg.optLevel = 1
	default:
		cfg.optLevel = 0
	}

	if fs.NArg() != 1 {
		fmt.Fprintf(stderr, "vertex: expected exactly one input (file or directory), got %d\n", fs.NArg())
		fs.Usage()
		return config{}, 2
	}
	cfg.input = fs.Arg(0)

	if cfg.mode != modeRun && cfg.output == "" {
		cfg.output = deriveOutput(cfg.input, cfg.mode, cfg.target)
	}

	return cfg, -1
}

func deriveOutput(input string, mode emitMode, target string) string {
	base := input
	if isDir(input) {
		base = filepath.Base(input)
	}
	switch mode {
	case modeVIR:
		return replaceExt(base, ".vir")
	case modeVBytes:
		return replaceExt(base, ".vbytes")
	case modeMIR:
		return replaceExt(base, ".mir")
	case modeASM:
		return replaceExt(base, ".s")
	case modeObj:
		if isWindowsTarget(target) {
			return replaceExt(base, ".obj")
		}
		return replaceExt(base, ".o")
	case modeExe:
		name := stripExt(base)
		if isWindowsTarget(target) {
			return name + ".exe"
		}
		return name
	case modeDump:
		return replaceExt(base, ".dump")
	}
	return replaceExt(base, ".out")
}

func defaultTarget() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	switch goarch {
	case "amd64", "arm64", "riscv64":
	default:
		goarch = "amd64"
	}
	switch goos {
	case "linux", "darwin", "windows":
		return goos + "-" + goarch
	default:
		return "linux-amd64"
	}
}

func expandShortFlag(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if len(a) >= 3 && a[0] == '-' && a[1] == 'o' && a[2] != '-' {
			out = append(out, "-o", a[2:])
		} else {
			out = append(out, a)
		}
	}
	return out
}

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

func supportedTargets() []string {
	return []string{
		"linux-amd64",
		"linux-arm64",
		"linux-riscv64",
		"darwin-amd64",
		"darwin-arm64",
		"windows-amd64",
		"windows-arm64",
		"freestanding-amd64",
		"freestanding-arm64",
		"freestanding-riscv64",
	}
}