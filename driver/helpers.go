package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// writeOutput writes data to path, creating any missing parent directories.
// Passing "-" writes to stdout.
func writeOutput(path string, data []byte) error {
	if path == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("cannot create output directory %s: %w", dir, err)
		}
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("cannot write %s: %w", path, err)
	}
	return nil
}

// replaceExt swaps the file extension of path for newExt (which must include
// the leading dot), or appends newExt if path has no extension.
func replaceExt(path, newExt string) string {
	if ext := filepath.Ext(path); ext != "" {
		return path[:len(path)-len(ext)] + newExt
	}
	return path + newExt
}

// stripExt removes the file extension from path, if any.
func stripExt(path string) string {
	if ext := filepath.Ext(path); ext != "" {
		return path[:len(path)-len(ext)]
	}
	return path
}

// appendSourceFiles writes the raw source of every .vs file that was compiled
// into sb, each prefixed with a file path comment. Used by dumpAll.
func appendSourceFiles(sb *strings.Builder, input string) error {
	var paths []string
	if isDir(input) {
		entries, err := os.ReadDir(input)
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", input, err)
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".vs") {
				paths = append(paths, filepath.Join(input, e.Name()))
			}
		}
	} else {
		paths = []string{input}
	}
	for _, p := range paths {
		src, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", p, err)
		}
		fmt.Fprintf(sb, "; file: %s\n", p)
		sb.Write(src)
		sb.WriteByte('\n')
	}
	return nil
}