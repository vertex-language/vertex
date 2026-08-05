package importer

import (
	"errors"
	"fmt"

	"github.com/vertex-language/vertex/analyzer"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/vertex/types"
)

// loadState tracks a path through the walk, so a cycle is diagnosed rather
// than recursed into and a failed package is not silently reused as if it had
// loaded cleanly.
type loadState uint8

const (
	unstarted loadState = iota
	loading
	loaded
	failed
)

type loader struct {
	conf     *Config
	resolver Resolver

	packages map[string]*Package
	state    map[string]loadState

	// order is the post-order accumulation: a package is appended after every
	// package it imports, which is exactly the topological order.
	order []*Package

	// stack is the current path from a root, used to render a cycle.
	stack []string

	list *diag.List
}

func (l *loader) report(d *diag.Diagnostic) {
	if d == nil {
		return
	}
	if l.conf.Reporter != nil {
		l.conf.Reporter.Report(d)
		return
	}
	l.list.Report(d)
}

// collector returns a Reporter that records into pkg.Errors and forwards on.
// Both the parse and the check go through one, so a package's own list holds
// its diagnostics from both phases in the order they were produced.
func (l *loader) collector(pkg *Package) diag.Reporter {
	return diag.ReporterFunc(func(d *diag.Diagnostic) {
		pkg.Errors = append(pkg.Errors, d)
		l.report(d)
	})
}

// load resolves one import path, recursing into its imports first.
//
// The post-order is not an optimization. §1.3 takes a package's qualifiers from
// its imports' own package clauses, so a package cannot be checked until every
// package it imports has been — and analyzer's collectImports asks the Importer
// interface for a complete *types.Package, not a promise of one.
func (l *loader) load(path string, from *Package) (*Package, error) {
	switch l.state[path] {
	case loaded:
		return l.packages[path], nil

	case failed:
		// Reached a second time through another importer. The first failure
		// was already returned; re-reporting it here would attribute one
		// broken directory to whichever package happened to name it second.
		return nil, fmt.Errorf("importer: %s failed to load", path)

	case loading:
		// An import cycle. Vertex has no forward-declaration form and no way
		// to break one, so this is fatal for the load rather than a diagnostic
		// on some particular file.
		return nil, &CycleError{Path: path, Stack: append([]string(nil), l.stack...)}
	}

	dir, err := l.resolver.Resolve(path, dirOf(from))
	if err != nil {
		l.state[path] = failed
		return nil, err
	}

	l.state[path] = loading
	l.stack = append(l.stack, path)
	ok := false
	defer func() {
		l.stack = l.stack[:len(l.stack)-1]
		if ok {
			l.state[path] = loaded
		} else {
			l.state[path] = failed
		}
	}()

	pkg := &Package{Path: path, Dir: dir, Imports: make(map[string]*Package)}

	// Registered before parsing, so a caller reading a partial Result can see
	// which package a fatal error was reached through.
	l.packages[path] = pkg

	if err := l.parse(pkg); err != nil {
		return nil, err
	}
	// Recurse before checking: every import must already be a complete
	// *types.Package before collectImports can bind its qualifier.
	if err := l.loadImports(pkg); err != nil {
		return nil, err
	}

	l.check(pkg)
	l.order = append(l.order, pkg)
	ok = true
	return pkg, nil
}

// parse reads the package's directory and parses the files that survive the
// build-tag filter.
//
// It reuses parser.ParseDir, which already runs the two passes: every file's
// package and build clauses are read against a throwaway FileSet first, because
// a file whose tag does not match the target is excluded from the build
// outright — so the filter has to run before any file is fully parsed.
func (l *loader) parse(pkg *Package) error {
	astPkg, err := parser.ParseDir(
		l.conf.Fset, pkg.Dir, pkg.Path, l.conf.Tag, l.collector(pkg), l.conf.Mode)

	if astPkg == nil {
		// ParseDir returns nothing at all only for a directory that holds no
		// source in this build: no .vs files, or none matching the tag. Both
		// are fatal to the load, since there is no package to hand a dependent.
		if err == nil {
			err = fmt.Errorf("no package produced")
		}
		return fmt.Errorf("%s: %w", pkg.Path, err)
	}

	// A non-nil error alongside a package is a syntax error, already reported
	// through the collector. Recovery produced Bad* nodes, so the tree is still
	// walkable and the check proceeds.
	pkg.Name = astPkg.Name
	pkg.Files = astPkg.Files
	return nil
}

// loadImports walks every import declaration in every file of pkg.
//
// There is no aliasing form, no dot-import, and no blank import, so a path maps
// to exactly one package and there is nothing to record but the mapping. Two
// files of one package naming the same path, or one file naming it twice, is
// one edge.
func (l *loader) loadImports(pkg *Package) error {
	for _, f := range pkg.Files {
		for _, decl := range f.Imports {
			for _, lit := range decl.Paths {
				path, ok := unquote(lit.Value)
				if !ok {
					continue // the parser already reported the malformed path
				}
				if _, done := pkg.Imports[path]; done {
					continue
				}

				dep, err := l.load(path, pkg)
				if err != nil {
					var ce *CycleError
					if errors.As(err, &ce) && ce.Importer == "" {
						// The innermost frame to see the cycle is the package
						// whose import list actually closed it. Frames above it
						// must not overwrite that with their own path.
						ce.Importer = pkg.Path
					}
					return err
				}
				pkg.Imports[path] = dep
			}
		}
	}
	return nil
}

// check runs the analyzer over one package.
func (l *loader) check(pkg *Package) {
	pkg.Info = types.NewInfo()

	conf := &analyzer.Config{
		Fset:     l.conf.Fset,
		Tag:      l.conf.Tag,
		Reporter: l.collector(pkg),
		Importer: &pkgImporter{pkg: pkg},
	}

	// A package clause is mandatory, but a file too broken to have one still
	// reaches here. The path is the only name left to check under.
	name := pkg.Name
	if name == "" {
		name = pkg.Path
	}

	tpkg, _ := analyzer.NewChecker(conf, pkg.Path, name, pkg.Info).Files(pkg.Files)
	pkg.Types = tpkg
	if tpkg == nil {
		return
	}

	// The analyzer marks a package complete only on a clean check, which is the
	// right default for a caller holding one package. A loader holding a graph
	// overrides it: a partial scope still answers most lookups, and refusing to
	// hand it to a dependent would turn one bad file into a cascade of
	// undeclared-name errors across everything downstream. The dependent's own
	// diagnostics are the honest report of what actually failed.
	tpkg.MarkComplete()
}

// pkgImporter satisfies analyzer.Importer for one package, over the
// dependencies the loader already resolved.
//
// It is deliberately not the whole graph. A package may reach only what it
// declares an import of: resolving a path this package never imported would let
// a name in through a transitive dependency, and §1.3's qualifier rule gives no
// spelling for that.
type pkgImporter struct{ pkg *Package }

func (i *pkgImporter) Import(path string) (*types.Package, error) {
	dep, ok := i.pkg.Imports[path]
	if !ok {
		return nil, fmt.Errorf("package %q is not imported by %s", path, i.pkg.Path)
	}
	if dep.Types == nil {
		return nil, fmt.Errorf("package %q produced no types", path)
	}
	return dep.Types, nil
}

// ---------------------------------------------------------------- errors

// CycleError reports an import cycle. It is fatal rather than a diagnostic:
// Vertex has no forward-declaration form, so there is no file to point at as
// the one that should have broken it.
type CycleError struct {
	// Path is the package reached a second time — the cycle's join point.
	Path string

	// Importer is the package whose own import list closed the cycle.
	Importer string

	// Stack is the root-to-here path at the point of detection.
	Stack []string
}

func (e *CycleError) Error() string {
	s := "import cycle: "
	for _, p := range e.Stack {
		s += p + " -> "
	}
	return s + e.Path
}

func dirOf(p *Package) string {
	if p == nil {
		return ""
	}
	return p.Dir
}

// unquote strips an import path's quotes. The scanner kept the raw spelling,
// escapes included, because a formatter needs the original — but an import path
// is a locator, and no escape has any meaning in one, so this only unwraps.
func unquote(s string) (string, bool) {
	if len(s) < 2 {
		return "", false
	}
	q := s[0]
	if (q != '"' && q != '`') || s[len(s)-1] != q {
		return "", false
	}
	return s[1 : len(s)-1], true
}