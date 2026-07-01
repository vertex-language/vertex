package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vertex-language/vertex/driver"
	"github.com/vertex-language/vertex/target"
)

// commonFlags are the flags shared across the build and test paths.
type commonFlags struct {
	target     string
	sysroot    string
	vertexHome string
	modFlag    string
	optLevel   int
	debugInfo  bool
}

func runBuildOrTest(args []string, stderr io.Writer) int {
	fs := newFlagSet("vertex", stderr)

	var (
		tf targetFlags
		xf testFlags

		output         string
		outputExplicit bool
		vertexHome     string
		modFlagVal     string
		sysroot        string
		debugInfo      bool

		fVIR, fVBytes, fMIR, fASM, fObj, fC, fDump bool
		fO0, fO1, fO2, fOs                         bool
		printVer                                   bool
	)

	registerTargetFlags(fs, &tf)
	registerTestFlags(fs, &xf)

	fs.BoolVar(&fVIR, "emit-vir", false, "emit Vertex IR text (.vir)")
	fs.BoolVar(&fVBytes, "emit-vbytes", false, "emit Vertex IR binary (.vbytes)")
	fs.BoolVar(&fMIR, "emit-mir", false, "emit Machine IR text (.mir)")
	fs.BoolVar(&fASM, "emit-asm", false, "emit native assembly text (.s)")
	fs.BoolVar(&fObj, "emit-obj", false, "emit relocatable object file (.o / .obj)")
	fs.BoolVar(&fC, "c", false, "compile to object file (alias for -emit-obj)")
	fs.BoolVar(&fDump, "dump", false, "dump all pipeline stages to a single annotated file (.dump)")
	fs.BoolVar(&fDump, "dump-all", false, "alias for -dump")

	fs.BoolVar(&fO0, "O0", false, "disable optimisation (default)")
	fs.BoolVar(&fO1, "O1", false, "light optimisation")
	fs.BoolVar(&fO2, "O2", false, "full optimisation")
	fs.BoolVar(&fOs, "Os", false, "optimise for size")

	fs.Func("o", "write output to `file`", func(s string) error {
		output = s
		outputExplicit = true
		return nil
	})
	fs.StringVar(&sysroot, "sysroot", "", "sysroot for cross-compilation library search")
	fs.StringVar(&vertexHome, "vertex-home", "", "override $VERTEX_HOME for this invocation (default: $VERTEX_HOME, or ~/.vertex)")
	fs.StringVar(&modFlagVal, "mod", "readonly", `dependency mode: "readonly" (default; never fetches or edits vs.sum) or "mod" (may fetch and record new vs.sum entries)`)
	fs.BoolVar(&debugInfo, "g", false, "include debug information")
	fs.BoolVar(&printVer, "version", false, "print version and exit")
	fs.BoolVar(&printVer, "v", false, "alias for -version")

	fs.Usage = func() { printUsage(stderr) }

	if err := fs.Parse(expandShortFlag(args)); err != nil {
		return 2
	}

	if printVer {
		fmt.Fprintf(os.Stdout, "vertex %s\n", version)
		return 0
	}
	if code, handled := handleTargetQueries(tf); handled {
		return code
	}

	common := commonFlags{
		target:     tf.target,
		sysroot:    sysroot,
		vertexHome: vertexHome,
		modFlag:    modFlagVal,
		debugInfo:  debugInfo,
	}
	switch {
	case fOs:
		common.optLevel = -1
	case fO2:
		common.optLevel = 2
	case fO1:
		common.optLevel = 1
	default:
		common.optLevel = 0
	}

	if xf.enabled {
		return runTestMode(fs, xf, common, stderr)
	}

	modes := []bool{fVIR, fVBytes, fMIR, fASM, fObj || fC, fDump}
	count := 0
	for _, b := range modes {
		if b {
			count++
		}
	}
	if count > 1 {
		fmt.Fprintln(stderr, "vertex: emit modes are mutually exclusive")
		return 2
	}

	var mode driver.EmitMode
	switch {
	case fVIR:
		mode = driver.ModeVIR
	case fVBytes:
		mode = driver.ModeVBytes
	case fMIR:
		mode = driver.ModeMIR
	case fASM:
		mode = driver.ModeASM
	case fObj || fC:
		mode = driver.ModeObj
	case fDump:
		mode = driver.ModeDump
	case outputExplicit:
		mode = driver.ModeExe
	default:
		mode = driver.ModeRun
	}

	if fs.NArg() != 1 {
		fmt.Fprintf(stderr, "vertex: expected exactly one input (file or directory), got %d\n", fs.NArg())
		fs.Usage()
		return 2
	}
	input := fs.Arg(0)

	if mode != driver.ModeRun && output == "" {
		output = deriveOutput(input, mode, common.target)
	}

	loadMode, err := parseLoadMode(common.modFlag)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 2
	}

	cfg := driver.Config{
		Input:      input,
		Output:     output,
		Target:     common.target,
		Sysroot:    common.sysroot,
		VertexHome: common.vertexHome,
		LoadMode:   loadMode,
		Mode:       mode,
		OptLevel:   common.optLevel,
		DebugInfo:  common.debugInfo,
	}
	return driver.Compile(cfg, stderr)
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

func deriveOutput(input string, mode driver.EmitMode, tri string) string {
	base := input
	if isDir(input) {
		base = filepath.Base(input)
	}
	switch mode {
	case driver.ModeVIR:
		return replaceExt(base, ".vir")
	case driver.ModeVBytes:
		return replaceExt(base, ".vbytes")
	case driver.ModeMIR:
		return replaceExt(base, ".mir")
	case driver.ModeASM:
		return replaceExt(base, ".s")
	case driver.ModeObj:
		if target.IsWindowsTarget(tri) {
			return replaceExt(base, ".obj")
		}
		return replaceExt(base, ".o")
	case driver.ModeExe:
		name := stripExt(base)
		if target.IsWindowsTarget(tri) {
			return name + ".exe"
		}
		return name
	case driver.ModeDump:
		return replaceExt(base, ".dump")
	}
	return replaceExt(base, ".out")
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func replaceExt(path, newExt string) string {
	if ext := filepath.Ext(path); ext != "" {
		return path[:len(path)-len(ext)] + newExt
	}
	return path + newExt
}

func stripExt(path string) string {
	if ext := filepath.Ext(path); ext != "" {
		return path[:len(path)-len(ext)]
	}
	return path
}

func printUsage(stderr io.Writer) {
	fmt.Fprintf(stderr, "Vertex compiler %s\n\n", version)
	fmt.Fprintf(stderr, "Usage:\n  vertex [flags] <source.vs | package/>\n  vertex mod init|get <module-path>[@version]\n\n")
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
	fmt.Fprintf(stderr, "  -o <file>           output file; when given, writes a permanent binary (no auto-run)\n")
	fmt.Fprintf(stderr, "  -target <triple>    target triple (default: %s)\n", target.DefaultTarget())
	fmt.Fprintf(stderr, "  -list-targets       list all supported targets and exit\n")
	fmt.Fprintf(stderr, "  -print-target       print the effective target triple and exit\n")
	fmt.Fprintf(stderr, "  -sysroot <path>     sysroot for cross-compilation library search\n")
	fmt.Fprintf(stderr, "  -vertex-home <dir>  override $VERTEX_HOME for this invocation\n")
	fmt.Fprintf(stderr, "  -mod readonly|mod   dependency mode (default: readonly)\n")
	fmt.Fprintf(stderr, "  -O0/-O1/-O2/-Os     optimisation level (default: -O0)\n")
	fmt.Fprintf(stderr, "  -g                  include debug information\n")
	fmt.Fprintf(stderr, "  -v, -version        print version and exit\n\n")
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
	fmt.Fprintf(stderr, "  vertex mod init example.com/widget\n")
	fmt.Fprintf(stderr, "  vertex mod get github.com/someuser/yourpackage\n")
}