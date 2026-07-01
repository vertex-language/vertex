package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/vertex-language/vertex/target"
)

type targetFlags struct {
	target      string
	listTargets bool
	printTarget bool
}

func registerTargetFlags(fs *flag.FlagSet, tf *targetFlags) {
	fs.StringVar(&tf.target, "target", target.DefaultTarget(), "target triple (os-arch)")
	fs.BoolVar(&tf.listTargets, "list-targets", false, "list all supported targets and exit")
	fs.BoolVar(&tf.printTarget, "print-target", false, "print the effective target triple and exit")
}

// handleTargetQueries prints and returns (code, true) if -list-targets or
// -print-target was given; (0, false) otherwise, meaning the caller should
// go on to actually compile something.
func handleTargetQueries(tf targetFlags) (int, bool) {
	if tf.listTargets {
		fmt.Fprintf(os.Stdout, "Supported targets:\n")
		for _, t := range target.SupportedTargets() {
			marker := ""
			if t == target.DefaultTarget() {
				marker = "  (default)"
			}
			fmt.Fprintf(os.Stdout, "  %s%s\n", t, marker)
		}
		return 0, true
	}
	if tf.printTarget {
		fmt.Fprintf(os.Stdout, "%s\n", tf.target)
		return 0, true
	}
	return 0, false
}