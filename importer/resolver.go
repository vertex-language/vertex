package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Resolver maps an import path to the directory holding that package's .vs
// files.
//
// It is an interface because A.2.3 ⊢ "the import path is a locator, not a
// name" and says nothing about what it locates against. A build system, a
// module cache, and a test fixture all answer differently, and none of them
// belongs in this package.
//
// from is the directory of the importing package, or "" for a root. A resolver
// that supports relative paths uses it; one that does not ignores it.
type Resolver interface {
	Resolve(path, from string) (dir string, err error)
}

// ResolverFunc adapts a function to Resolver.
type ResolverFunc func(path, from string) (string, error)

func (f ResolverFunc) Resolve(path, from string) (string, error) { return f(path, from) }

// DirResolver resolves an import path against one or more root directories.
//
// A path is a slash-separated locator interpreted relative to each root in
// order; the first root holding a directory with at least one .vs file wins.
// That is deliberately the simplest thing that works — there is no module
// resolution, no version selection, and no vendor directory, because the
// language annex does not describe any and inventing one here would fix a
// policy the toolchain has not chosen yet.
type DirResolver struct {
	Roots []string
}

// NewDirResolver returns a resolver over the given roots, each made absolute.
func NewDirResolver(roots ...string) (*DirResolver, error) {
	if len(roots) == 0 {
		roots = []string{"."}
	}
	abs := make([]string, 0, len(roots))
	for _, r := range roots {
		a, err := filepath.Abs(r)
		if err != nil {
			return nil, err
		}
		abs = append(abs, a)
	}
	return &DirResolver{Roots: abs}, nil
}

func (r *DirResolver) Resolve(path, from string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("importer: empty import path")
	}

	// A relative path resolves against the importing package. This is not in
	// the annex either way; it is here because a single-directory or
	// two-directory project should not need a root configured.
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
		if from == "" {
			return "", fmt.Errorf("importer: relative path %q has no importing package", path)
		}
		dir := filepath.Join(from, filepath.FromSlash(path))
		if hasSources(dir) {
			return dir, nil
		}
		return "", &NotFoundError{Path: path, Searched: []string{dir}}
	}

	rel := filepath.FromSlash(path)
	searched := make([]string, 0, len(r.Roots))

	for _, root := range r.Roots {
		dir := filepath.Join(root, rel)
		if hasSources(dir) {
			return dir, nil
		}
		searched = append(searched, dir)
	}

	// A directory that exists but holds no .vs file is a different mistake
	// from one that does not exist, and the message should say which.
	for _, dir := range searched {
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			return "", &NotFoundError{
				Path:     path,
				Searched: searched,
				Reason:   "directory exists but contains no .vs files",
			}
		}
	}
	return "", &NotFoundError{Path: path, Searched: searched}
}

// MapResolver resolves paths from a fixed table. Useful for a driver that has
// already computed the layout, and for tests.
type MapResolver map[string]string

func (m MapResolver) Resolve(path, from string) (string, error) {
	if dir, ok := m[path]; ok {
		return dir, nil
	}
	return "", &NotFoundError{Path: path}
}

// NotFoundError reports an unresolvable import path.
type NotFoundError struct {
	Path     string
	Searched []string
	Reason   string
}

func (e *NotFoundError) Error() string {
	s := fmt.Sprintf("cannot resolve import %q", e.Path)
	if e.Reason != "" {
		s += ": " + e.Reason
	}
	if len(e.Searched) > 0 {
		s += "\n  searched:"
		for _, d := range e.Searched {
			s += "\n    " + d
		}
	}
	return s
}

// hasSources reports whether dir holds at least one .vs file.
//
// The check is by extension only. Whether a given file is actually in the
// build is A.2.2's build-tag question, which parser.ParseDir answers in its
// first pass — a directory whose every file is tagged for another target is
// resolvable but yields no package, and that error belongs there rather than
// here.
func hasSources(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".vs") {
			return true
		}
	}
	return false
}