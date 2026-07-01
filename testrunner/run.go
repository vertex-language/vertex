// Package testrunner discovers `build test`-tagged functions and runs each
// as its own tiny compiled program.
//
// It depends only on driver's public entry points (via the Compiler
// interface below), never on driver's internals — driver.Test supplies an
// implementation as a closure. This keeps driver -> testrunner strictly
// one-way: testrunner has no import of driver at all, so there's no
// import cycle even though driver is the only caller of Run.
package testrunner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vertex-language/ir/vertex/ast"
)

type Options struct {
	Target    string
	Sysroot   string
	OptLevel  int
	DebugInfo bool
	TestDir   string
	TestFile  string
}

// Compiler is the subset of driver's public surface testrunner needs:
// compile a synthetic package to an executable at output, or dump its
// full pipeline to path for a failed test's post-mortem.
type Compiler interface {
	Compile(pkg *ast.Package, output string, stderr io.Writer) int
	Dump(pkg *ast.Package, path string, stderr io.Writer)
}

func Run(opts Options, compiler Compiler, stderr io.Writer) int {
	files, err := collectTestFiles(opts.TestDir, opts.TestFile)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	if len(files) == 0 {
		fmt.Fprintf(stderr, "vertex: no test files found\n")
		return 1
	}

	var cases []TestCase
	for _, path := range files {
		tc, err := parseTestCases(path)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %s: %v\n", path, err)
			return 1
		}
		cases = append(cases, tc...)
	}
	if len(cases) == 0 {
		fmt.Fprintf(stderr, "vertex: no test functions found\n")
		return 1
	}

	pass, fail := 0, 0
	for _, tc := range cases {
		if execTest(tc, opts, compiler, stderr) {
			pass++
		} else {
			fail++
		}
	}

	fmt.Fprintf(os.Stdout, "\n%d passed, %d failed\n", pass, fail)
	if fail > 0 {
		return 1
	}
	return 0
}

func execTest(tc TestCase, opts Options, compiler Compiler, stderr io.Writer) bool {
	label := fmt.Sprintf("%s::%s", filepath.Base(tc.File), tc.FuncName)

	tmpDir, err := os.MkdirTemp("", "vertex-test-*")
	if err != nil {
		fmt.Fprintf(os.Stdout, "FAIL  %s  (tmpdir: %v)\n", label, err)
		return false
	}
	defer os.RemoveAll(tmpDir)

	binPath := filepath.Join(tmpDir, tc.FuncName)
	if strings.HasPrefix(strings.ToLower(opts.Target), "windows-") {
		binPath += ".exe"
	}

	pkg := buildSyntheticPackage(tc)

	var compileErr strings.Builder
	if code := compiler.Compile(pkg, binPath, &compileErr); code != 0 {
		fmt.Fprintf(os.Stdout, "FAIL  %s  (compile error)\n", label)
		fmt.Fprint(os.Stdout, compileErr.String())
		saveTestArtifacts(pkg, compiler, tc, "", stderr)
		return false
	}

	out, err := exec.Command(binPath).Output()
	if err != nil {
		fmt.Fprintf(os.Stdout, "FAIL  %s  (exec: %v)\n", label, err)
		saveTestArtifacts(pkg, compiler, tc, binPath, stderr)
		return false
	}

	got := strings.TrimRight(string(out), "\r\n")
	if got == tc.ExpectedValue {
		fmt.Fprintf(os.Stdout, "ok    %s\n", label)
		return true
	}
	fmt.Fprintf(os.Stdout, "FAIL  %s\n", label)
	fmt.Fprintf(os.Stdout, "      want: %q\n", tc.ExpectedValue)
	fmt.Fprintf(os.Stdout, "      got:  %q\n", got)
	saveTestArtifacts(pkg, compiler, tc, binPath, stderr)
	return false
}

// saveTestArtifacts writes a pipeline dump and, if the binary was built
// successfully, copies it into ./dumps for a failed test's post-mortem.
func saveTestArtifacts(pkg *ast.Package, compiler Compiler, tc TestCase, binPath string, stderr io.Writer) {
	if err := os.MkdirAll("dumps", 0o755); err != nil {
		fmt.Fprintf(stderr, "vertex: cannot create dumps dir: %v\n", err)
		return
	}
	base := testDumpBase(tc)
	compiler.Dump(pkg, filepath.Join("dumps", base+".dump"), stderr)

	if binPath != "" {
		data, err := os.ReadFile(binPath)
		if err == nil {
			binDst := filepath.Join("dumps", base+".bin")
			if err := os.WriteFile(binDst, data, 0o755); err != nil {
				fmt.Fprintf(stderr, "vertex: cannot write %s: %v\n", binDst, err)
			}
		}
	}
}