package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/vertex-language/pkg"
)

const version = "0.4.0"

// Main is the CLI's only exported entry point; cmd/vertex/main.go just
// calls this with os.Args[1:] and os.Stderr. "mod" is the one positional
// subcommand keyword; everything else is the flag-based build/run/test
// invocation handled by runBuildOrTest.
func Main(args []string, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "mod" {
		return runMod(args[1:], stderr)
	}
	return runBuildOrTest(args, stderr)
}

func newFlagSet(name string, stderr io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	return fs
}

// parseLoadMode maps the shared -mod flag's string value to a pkg.LoadMode.
func parseLoadMode(s string) (pkg.LoadMode, error) {
	switch s {
	case "", "readonly":
		return pkg.ModReadonly, nil
	case "mod":
		return pkg.ModUpdate, nil
	default:
		return 0, fmt.Errorf("invalid -mod value %q (want %q or %q)", s, "readonly", "mod")
	}
}