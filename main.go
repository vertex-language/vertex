// main.go
//
// vertex — the Vertex compiler.
//
// This file exists only so the module root is `package main` and
// `go install github.com/vertex-language/vertex@latest` produces a `vertex`
// binary. Argument parsing lives in cli/, and every pipeline decision lives
// in driver/; nothing here does either.
package main

import (
	"os"

	"github.com/vertex-language/vertex/cli"
)

func main() {
	os.Exit(cli.Main(os.Args[1:]))
}