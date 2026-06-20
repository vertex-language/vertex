package main

import (
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}

func run(args []string, stderr io.Writer) int {
	cfg, code := parseFlags(args, stderr)
	if code >= 0 {
		return code
	}
	return emit(cfg, stderr)
}