// Package driver implements the Vertex compiler pipeline and is the public
// entry point for the vertex command and any tooling that embeds the compiler
// (language servers, build systems, incremental cache layers, …).
//
// Cross-package compilation — discovering the import DAG, topologically
// compiling dependencies, and threading the resulting vertex.Module map into
// virlower.NewLower — will live here once the sema pass exists. For now the
// pipeline compiles a single package and records the gap explicitly.
package driver

import "io"

// Run parses args, runs the full compiler pipeline, and returns a POSIX exit
// code (0 success, 1 compile error, 2 usage/flag error).
func Run(args []string, stderr io.Writer) int {
	cfg, code := parseFlags(args, stderr)
	if code >= 0 {
		return code
	}
	return emit(cfg, stderr)
}