package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vertex-language/ir/vertex/ast"

	"github.com/vertex-language/vertex/pipeline"
)

// sourceDir returns the absolute directory containing input, whether
// input is itself a directory or a single file.
func sourceDir(input string) (string, error) {
	dir := input
	if !isDir(dir) {
		dir = filepath.Dir(dir)
	}
	return filepath.Abs(dir)
}

// findModuleRoot walks upward from input (a file or directory) looking
// for vs.mod, the same way `go build` locates a module root from any file
// inside it. Returns "" (not an error) if no vs.mod is found anywhere up
// the tree — the caller decides what that means: no non-stdlib imports
// makes it a non-issue, a real import makes it a case for resolving a
// graph some other way (see driver.Compile's use of loadGraphFromImports).
func findModuleRoot(input string) (string, error) {
	dir, err := sourceDir(input)
	if err != nil {
		return "", fmt.Errorf("resolving %s: %w", input, err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "vs.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}
		dir = parent
	}
}

// collectNonStdlibImports returns the distinct set of import paths in p
// that look like module paths rather than standard-library packages —
// i.e. whose first path segment contains a ".", the same convention
// vs.mod's own module paths use throughout this toolchain — in the
// order they're first seen. Each file's Imports is a list of
// *ast.ImportDecl, and each ImportDecl carries one or more Paths (more
// than one for the grouped `import ( ... )` form), so this walks both
// levels.
func collectNonStdlibImports(p *ast.Package) []string {
	seen := make(map[string]bool)
	var out []string
	for _, f := range p.Files {
		for _, imp := range f.Imports {
			for _, path := range imp.Paths {
				seg, _, _ := strings.Cut(path.Path, "/")
				if !strings.Contains(seg, ".") {
					continue
				}
				if seen[path.Path] {
					continue
				}
				seen[path.Path] = true
				out = append(out, path.Path)
			}
		}
	}
	return out
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func parseInput(input string) (*ast.Package, error) {
	var files []*ast.File
	if isDir(input) {
		entries, err := os.ReadDir(input)
		if err != nil {
			return nil, fmt.Errorf("cannot read directory %s: %w", input, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".vs") {
				continue
			}
			f, err := parseFile(filepath.Join(input, e.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, f)
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("no .vs files found in %s", input)
		}
	} else {
		f, err := parseFile(input)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return ast.NewPackage(files)
}

func parseFile(path string) (*ast.File, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return ast.NewFile(path, src)
}

// appendSourceFiles writes the raw source of every .vs file under u's
// directory into sb, each prefixed with a file path comment. Used by dump
// mode ahead of the pipeline's own stage banners. A unit with no Dir (a
// synthetic in-memory package, e.g. a test case) has nothing to echo.
func appendSourceFiles(sb *strings.Builder, u *pipeline.Unit) error {
	if u.Dir == "" {
		return nil
	}
	var paths []string
	if isDir(u.Dir) {
		entries, err := os.ReadDir(u.Dir)
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", u.Dir, err)
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".vs") {
				paths = append(paths, filepath.Join(u.Dir, e.Name()))
			}
		}
	} else {
		paths = []string{u.Dir}
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

func writeOutput(path string, data []byte) error { return writeFile(path, data, 0o644) }
func writeExe(path string, data []byte) error    { return writeFile(path, data, 0o755) }

func writeFile(path string, data []byte, perm os.FileMode) error {
	if path == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("cannot create output directory %s: %w", dir, err)
		}
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("cannot write %s: %w", path, err)
	}
	return nil
}

func stripExt(path string) string {
	if ext := filepath.Ext(path); ext != "" {
		return path[:len(path)-len(ext)]
	}
	return path
}

func boolToCode(failed bool) int {
	if failed {
		return 1
	}
	return 0
}