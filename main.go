package main

import (
	"os"

	"github.com/vertex-language/vertex/driver"
)

func main() { os.Exit(driver.Run(os.Args[1:], os.Stderr)) }