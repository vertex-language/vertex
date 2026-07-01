package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/vertex-language/vertex/driver"
)

type testFlags struct {
	enabled bool
	dir     string
	file    string
}

func registerTestFlags(fs *flag.FlagSet, tf *testFlags) {
	fs.BoolVar(&tf.enabled, "test", false, "discover and run test functions")
	fs.StringVar(&tf.dir, "dir", "", "directory to search for test files (recursive)")
	fs.StringVar(&tf.file, "file", "", "single test file to run")
}

func runTestMode(fs *flag.FlagSet, tf testFlags, common commonFlags, stderr io.Writer) int {
	if tf.file != "" && tf.dir != "" {
		fmt.Fprintln(stderr, "vertex: -test: -file and -dir are mutually exclusive")
		return 2
	}

	testDir := tf.dir
	if tf.file == "" && tf.dir == "" {
		switch fs.NArg() {
		case 0:
			testDir = "."
		case 1:
			testDir = fs.Arg(0)
		default:
			fmt.Fprintln(stderr, "vertex: -test: expected at most one positional argument")
			return 2
		}
	} else if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "vertex: -test: unexpected positional argument when -file or -dir is set")
		return 2
	}

	cfg := driver.Config{
		Target:    common.target,
		Sysroot:   common.sysroot,
		Mode:      driver.ModeTest,
		OptLevel:  common.optLevel,
		DebugInfo: common.debugInfo,
		TestDir:   testDir,
		TestFile:  tf.file,
	}
	return driver.Test(cfg, stderr)
}