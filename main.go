package main

import (
	"os"

	"github.com/vertex-language/vertex/cli"
)

func main() { os.Exit(cli.Main(os.Args[1:], os.Stderr)) }