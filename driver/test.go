package driver

import (
	"io"
	"strings"

	"github.com/vertex-language/ir/vertex/ast"

	"github.com/vertex-language/vertex/testrunner"
)

// Test discovers and runs every `build test`-tagged function reachable
// from cfg.TestDir (or cfg.TestFile), each compiled as its own tiny
// synthetic program via CompilePackage.
func Test(cfg Config, stderr io.Writer) int {
	opts := testrunner.Options{
		Target:    cfg.Target,
		Sysroot:   cfg.Sysroot,
		OptLevel:  cfg.OptLevel,
		DebugInfo: cfg.DebugInfo,
		TestDir:   cfg.TestDir,
		TestFile:  cfg.TestFile,
	}
	return testrunner.Run(opts, compilerAdapter{cfg}, stderr)
}

// compilerAdapter implements testrunner.Compiler over this package's own
// public surface (CompilePackage, DumpPackage). testrunner itself never
// imports driver — this closure is what keeps driver -> testrunner a
// one-way dependency instead of a cycle, since testrunner only needs
// driver's public entry points, never its internals.
type compilerAdapter struct{ base Config }

func (a compilerAdapter) Compile(p *ast.Package, output string, stderr io.Writer) int {
	cfg := a.base
	cfg.Mode = ModeExe
	cfg.Output = output
	var sink strings.Builder
	code := CompilePackage(p, cfg, &sink)
	io.WriteString(stderr, sink.String())
	return code
}

func (a compilerAdapter) Dump(p *ast.Package, path string, stderr io.Writer) {
	DumpPackage(p, a.base, path, stderr)
}