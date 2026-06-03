package compiler

import "os"

// osReadFile wraps os.ReadFile so that compiler.go does not need to import os
// separately (os is already imported in imports.go).
func osReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}