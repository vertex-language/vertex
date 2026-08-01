// cli/usage.go
package cli

import (
	"fmt"
	"io"

	"github.com/vertex-language/vertex/driver"
)

func printUsage(w io.Writer) {
	fmt.Fprintf(w, `vertex %s — the Vertex compiler (language spec %s)

Usage:
  vertex [build] [flags] <file.vs | package-dir>
      Compile a package to a native executable for the host or --target.
      The command word is optional: "vertex main.vs" is "vertex build main.vs".

  vertex run [flags] <file.vs | package-dir> [-- args...]
      Build for the host, execute immediately, forward the exit code.
      Anything after -- is passed to the compiled program, not to vertex.

  vertex test [-dir <path> | -file <path>] [-run <substr>]
      Discover 'test'-marked functions in a `+"`build test`"+` package and run
      them, comparing each against its Expected(...) result.

  vertex targets
      List every target this toolchain can actually build for.

  vertex version | help

Build flags:
  -o <path>             output file (default: derived from the input name).
                          For -emit-vir/-emit-vbyte on a multi-package build,
                          this must be a directory — one file per package.
  -target <triple>      see "vertex targets"; defaults to the host
  -min-os-version <v>   required for darwin targets (Apple's triple grammar
                          has no unversioned form); defaults to %s
  -shared               build a shared library instead of an executable
  -flat-base <addr>     load address for a freestanding flat image
  -root <module>        override which module's entry function is the
                          program's entry point (default: the package
                          holding main)
  -packages-dir <path>  packages root; overrides $VERTEX_PATH
  -emit-vir             emit Vertex IR text (.vir), one file per package
  -emit-vbyte           emit Vertex IR binary (.vbyte), one file per package
  -v                    report each pipeline stage on stderr

Examples:
  vertex main.vs
  vertex -o app ./cmd/app
  vertex build -target darwin-arm64 -o app main.vs
  vertex run main.vs -- --verbose
  vertex build -emit-vir -o build/ ./cmd/app
  vertex test -dir ./tests
`, driver.Version, driver.SpecVersion, driver.DefaultMinOSVersion)
}