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
// It is an interface because the path is a locator rather than a name, and
// nothing in the language says what it locates against. A build system, a
// module cache, and a test fixture all answer differently, and none of those
// answers belongs in this package.
//
// from is the directory of the importing package, or "" for a root. A resolver
// supporting relative paths uses it; one that does not ignores it.
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
// That is deliberately the simplest thing that works — no module resolution, no
// version selection, no vendor directory — because the language describes none
// of those and inventing one here would fix a policy the toolchain has not
// chosen.
type DirResolver struct {
	Roots []string
}

// NewDirResolver returns a resolver over the given roots, each made absolute so
// a later working-directory change cannot move them.
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
	if filepath.IsAbs(path) {
		// An absolute path is not a locator relative to anything, so no root
		// applies and no policy here could make one apply.
		return "", fmt.Errorf("importer: import path %q must be relative", path)
	}

	// A relative path resolves against the importing package. Nothing in the
	// language says whether one is even spelled; this is here so a project of
	// one or two directories needs no root configured.
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
		if from == "" {
			return "", fmt.Errorf(
				"importer: relative path %q has no importing package", path)
		}
		dir := filepath.Clean(filepath.Join(from, filepath.FromSlash(path)))
		if hasSources(dir) {
			return dir, nil
		}
		return "", notFound(path, []string{dir})
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
	return "", notFound(path, searched)
}

// MapResolver resolves paths from a fixed table. For a driver that already
// computed the layout, and for tests wanting no filesystem at all.
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
	for i, d := range e.Searched {
		if i == 0 {
			s += "\n  searched:"
		}
		s += "\n    " + d
	}
	return s
}

// notFound distinguishes a directory that exists but holds no source from one
// that does not exist. They are different mistakes and the message should say
// which.
func notFound(path string, searched []string) *NotFoundError {
	e := &NotFoundError{Path: path, Searched: searched}
	for _, dir := range searched {
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			e.Reason = "directory exists but contains no .vs files"
			break
		}
	}
	return e
}

// hasSources reports whether dir holds at least one .vs file.
//
// The check is by extension only. Whether a given file is in the build is the
// build tag's question, which parser.ParseDir answers in its own first pass — a
// directory whose every file is tagged for another target resolves successfully
// and yields no package, and that error belongs there rather than here.
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