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

// emitMode controls how far through the pipeline to run and what to write out.
type emitMode uint8

const (
	modeVIR    emitMode = iota // -emit-vir    → Vertex IR text      (.vir)
	modeVBytes                  // -emit-vbytes → Vertex IR binary    (.vbytes)
	modeMIR                     // -emit-mir    → Machine IR text      (.mir)
	modeASM                     // -emit-asm    → native assembly text (.s)
	modeObj                     // -c           → relocatable object   (.o / .obj)
	modeExe                     // -lc          → native executable
)

type config struct {
	input       string
	output      string
	target      string   // e.g. "linux-amd64"
	packagesDir string
	mode        emitMode
	optLevel    int  // 0=none 1=light 2=full -1=size
	debugInfo   bool
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
		fObj    bool // -emit-obj
		fC      bool // -c (alias)
		fLC     bool // -lc

		fO0, fO1, fO2, fOs bool
		printVer            bool
	)

	// ── Emit modes ────────────────────────────────────────────────────────────
	fs.BoolVar(&fVIR,    "emit-vir",    false, "emit Vertex IR text (.vir)")
	fs.BoolVar(&fVBytes, "emit-vbytes", false, "emit Vertex IR binary (.vbytes)")
	fs.BoolVar(&fMIR,    "emit-mir",    false, "emit Machine IR text (.mir)")
	fs.BoolVar(&fASM,    "emit-asm",    false, "emit native assembly text (.s)")
	fs.BoolVar(&fObj,    "emit-obj",    false, "emit relocatable object file (.o / .obj)")
	fs.BoolVar(&fC,      "c",           false, "compile to object file (alias for -emit-obj)")
	fs.BoolVar(&fLC,     "lc",          false, "compile and link to a native executable")

	// ── Optimisation ──────────────────────────────────────────────────────────
	fs.BoolVar(&fO0, "O0", false, "disable optimisation (default)")
	fs.BoolVar(&fO1, "O1", false, "light optimisation")
	fs.BoolVar(&fO2, "O2", false, "full optimisation")
	fs.BoolVar(&fOs, "Os", false, "optimise for size")

	// ── Common options ────────────────────────────────────────────────────────
	fs.StringVar(&cfg.output,      "o",            "",                   "write output to `file`")
	fs.StringVar(&cfg.target,      "target",        defaultTarget(),      "target triple (os-arch)")
	fs.StringVar(&cfg.packagesDir, "packages-dir",  defaultPackagesDir(), "Vertex packages root (overrides $VERTEX_PATH)")
	fs.BoolVar(&cfg.debugInfo,     "g",             false,                "include debug information")
	fs.BoolVar(&printVer, "version", false, "print version and exit")
	fs.BoolVar(&printVer, "v",       false, "alias for -version")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Vertex compiler %s\n\n", version)
		fmt.Fprintf(stderr, "Usage:\n  vertex [flags] <source.vs | package/>\n\n")
		fmt.Fprintf(stderr, "Emit mode (exactly one required):\n")
		fmt.Fprintf(stderr, "  -emit-vir       emit Vertex IR text (.vir)\n")
		fmt.Fprintf(stderr, "  -emit-vbytes    emit Vertex IR binary (.vbytes)\n")
		fmt.Fprintf(stderr, "  -emit-mir       emit Machine IR text (.mir)\n")
		fmt.Fprintf(stderr, "  -emit-asm       emit native assembly text (.s)\n")
		fmt.Fprintf(stderr, "  -emit-obj, -c   emit relocatable object file (.o / .obj)\n")
		fmt.Fprintf(stderr, "  -lc             compile and link to native executable\n\n")
		fmt.Fprintf(stderr, "Options:\n")
		fmt.Fprintf(stderr, "  -o <file>        output file (default: derived from input)\n")
		fmt.Fprintf(stderr, "  -target <triple> linux-amd64, linux-arm64, linux-riscv64,\n")
		fmt.Fprintf(stderr, "                   darwin-amd64, darwin-arm64,\n")
		fmt.Fprintf(stderr, "                   windows-amd64, windows-arm64,\n")
		fmt.Fprintf(stderr, "                   freestanding-amd64, freestanding-arm64,\n")
		fmt.Fprintf(stderr, "                   freestanding-riscv64  (default: %s)\n", defaultTarget())
		fmt.Fprintf(stderr, "  -packages-dir    Vertex packages root (overrides $VERTEX_PATH)\n")
		fmt.Fprintf(stderr, "  -O0/-O1/-O2/-Os  optimisation level (default: -O0)\n")
		fmt.Fprintf(stderr, "  -g               include debug information\n")
		fmt.Fprintf(stderr, "  -v, -version     print version and exit\n\n")
		fmt.Fprintf(stderr, "Examples:\n")
		fmt.Fprintf(stderr, "  vertex -emit-vir    -o main.vir    main.vs\n")
		fmt.Fprintf(stderr, "  vertex -emit-vbytes -o main.vbytes main.vs\n")
		fmt.Fprintf(stderr, "  vertex -emit-mir    -o main.mir    main.vs\n")
		fmt.Fprintf(stderr, "  vertex -emit-asm    -o main.s      main.vs\n")
		fmt.Fprintf(stderr, "  vertex -c           -o main.o      main.vs\n")
		fmt.Fprintf(stderr, "  vertex -lc          -o main        main.vs\n")
		fmt.Fprintf(stderr, "  vertex -lc -target darwin-arm64 -O2 -o main main.vs\n")
	}

	if err := fs.Parse(expandShortFlag(args)); err != nil {
		return config{}, 2
	}
	if printVer {
		fmt.Fprintf(os.Stdout, "vertex %s\n", version)
		return config{}, 0
	}

	// Exactly one emit mode must be set.
	modes := []bool{fVIR, fVBytes, fMIR, fASM, fObj || fC, fLC}
	count := 0
	for _, b := range modes {
		if b {
			count++
		}
	}
	switch count {
	case 0:
		fmt.Fprintf(stderr, "vertex: one emit mode is required (-emit-vir, -emit-vbytes, -emit-mir, -emit-asm, -c/-emit-obj, -lc)\n\n")
		fs.Usage()
		return config{}, 2
	case 1:
		// good
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
	case fLC:
		cfg.mode = modeExe
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

	if cfg.output == "" {
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
	}
	return replaceExt(base, ".out")
}

// defaultTarget returns the host OS and arch as a Vertex triple.
func defaultTarget() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	switch goarch {
	case "amd64", "arm64", "riscv64":
		// already matching Vertex arch names
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

// expandShortFlag rewrites -oFILE → -o FILE so the stdlib flag parser is happy.
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