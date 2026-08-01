// Package importer resolves a Vertex import graph into checked packages.
//
// It is the seam ast.NewPackage's doc comment names: "Resolution belongs to the
// loader, which is why this does not take an importer the way Go's deprecated
// ast.NewPackage did." parser.ParseDir handles one directory; this handles the
// graph.
//
// The ordering constraint is A.2.3: "the imported package's declared name (its
// PackageClause) is the qualifier under which its symbols are reached; the
// import path is a locator, not a name." A file's names therefore cannot be
// resolved until every directory it imports has had its package clause read —
// which is what parser.PackageClauseOnly exists for, and why loading is a
// post-order walk rather than a single pass.
//
// Nothing here type-checks. It parses, orders, and calls the analyzer; the
// analyzer decides what the result means.
package importer

import (
	"fmt"
	"sync"

	"github.com/vertex-language/vertex/analyzer"
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// Config controls a load.
type Config struct {
	// Fset is the position space every loaded file shares. One per load, so a
	// diagnostic can span two files in different packages.
	Fset *token.FileSet

	// Tag is the target (A.2.2). It decides which files are in the build at
	// all, and — because `build test` "is the only build tag that changes what
	// is grammatical rather than only what is linkable" — what parses.
	Tag token.BuildTag

	// Resolver maps an import path to a directory. Nil means DirResolver
	// rooted at the current directory.
	Resolver Resolver

	// Reporter receives every diagnostic from every package. May be nil, in
	// which case diagnostics are collected into the returned *Result.
	Reporter diag.Reporter

	// Mode is passed through to the parser. ParseComments is the usual choice
	// for tooling and unnecessary for a build.
	Mode parser.Mode
}

// Package is one loaded package: its syntax, its checked types, and its
// resolved dependencies.
type Package struct {
	Path  string
	Dir   string
	Name  string // from the PackageClause (A.2.1)
	Files []*ast.File

	Types *types.Package
	Info  *types.Info

	// Imports maps each import path this package declares to the loaded
	// package it resolved to.
	Imports map[string]*Package

	// Errors is this package's own diagnostics, in position order. A package
	// with errors is still returned — partial results are what editor tooling
	// needs, and A.0.5 makes rejected forms part of the grammar rather than
	// something that aborts a load.
	Errors []*diag.Diagnostic
}

func (p *Package) String() string { return p.Path }

// Result is one load's output.
type Result struct {
	// Roots are the packages named in the Load call, in the order given.
	Roots []*Package

	// Packages holds every loaded package by import path, roots included.
	Packages map[string]*Package

	// Order is a topological order: a package appears after everything it
	// imports. This is the order the analyzer ran in and the order lower
	// should follow, since A.2.3 makes a package's qualifiers depend on its
	// imports having been read.
	Order []*Package

	// Diagnostics is every diagnostic from every package, sorted and deduped.
	// Empty when Config.Reporter is non-nil, since they streamed there instead.
	Diagnostics *diag.List
}

// HasErrors reports whether any loaded package produced an error-severity
// diagnostic.
func (r *Result) HasErrors() bool {
	for _, p := range r.Order {
		for _, d := range p.Errors {
			if d.Sev == diag.Error {
				return true
			}
		}
	}
	return r.Diagnostics != nil && r.Diagnostics.HasErrors()
}

// Load resolves and checks the packages named by paths, plus everything they
// import transitively.
func Load(conf *Config, paths ...string) (*Result, error) {
	if conf == nil {
		return nil, fmt.Errorf("importer: nil config")
	}
	if conf.Fset == nil {
		return nil, fmt.Errorf("importer: nil FileSet")
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("importer: no packages named")
	}

	res := conf.Resolver
	if res == nil {
		r, err := NewDirResolver(".")
		if err != nil {
			return nil, err
		}
		res = r
	}

	l := &loader{
		conf:     conf,
		resolver: res,
		packages: make(map[string]*Package),
		state:    make(map[string]loadState),
	}
	if conf.Reporter == nil {
		l.list = &diag.List{}
	}

	out := &Result{
		Packages: l.packages,
		Diagnostics: l.list,
	}

	for _, path := range paths {
		pkg, err := l.load(path, nil)
		if err != nil {
			return out, err
		}
		out.Roots = append(out.Roots, pkg)
	}

	out.Order = l.order
	if l.list != nil {
		l.list.Sort()
		l.list.Dedup()
	}
	return out, nil
}

// loadState tracks a path through the walk, so a cycle is diagnosed rather
// than recursed into.
type loadState uint8

const (
	unstarted loadState = iota
	loading
	loaded
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

// load resolves one import path, recursing into its imports first.
//
// The post-order is not an optimization. A.2.3 takes a package's qualifiers
// from its imports' own PackageClauses, so a package cannot be checked until
// every package it imports has been — and the analyzer's collectImports asks
// the Importer interface for a complete *types.Package, not a promise of one.
func (l *loader) load(path string, from *Package) (*Package, error) {
	switch l.state[path] {
	case loaded:
		return l.packages[path], nil

	case loading:
		// An import cycle. Vertex has no forward-declaration form and no way
		// to break one, so this is fatal for the load rather than a diagnostic
		// on some particular file.
		return nil, &CycleError{Path: path, Stack: append([]string(nil), l.stack...)}
	}

	dir, err := l.resolver.Resolve(path, dirOf(from))
	if err != nil {
		return nil, err
	}

	l.state[path] = loading
	l.stack = append(l.stack, path)
	defer func() {
		l.stack = l.stack[:len(l.stack)-1]
		l.state[path] = loaded
	}()

	pkg := &Package{
		Path:    path,
		Dir:     dir,
		Imports: make(map[string]*Package),
	}
	l.packages[path] = pkg

	if err := l.parse(pkg); err != nil {
		return nil, err
	}

	// Recurse before checking. Every import must be a complete
	// *types.Package before collectImports can bind its qualifier.
	if err := l.loadImports(pkg); err != nil {
		return nil, err
	}

	l.check(pkg)
	l.order = append(l.order, pkg)
	return pkg, nil
}

// parse reads the package's directory and parses the files that survive the
// build-tag filter.
//
// It reuses parser.ParseDir, which already implements A.2.2's two-pass shape:
// every file's package and build clauses are read into a throwaway FileSet
// first, because "a file whose tag does not match the current target is
// excluded from the build whole, never partially" — so the filter has to run
// before any file is fully parsed.
func (l *loader) parse(pkg *Package) error {
	collector := diag.ReporterFunc(func(d *diag.Diagnostic) {
		pkg.Errors = append(pkg.Errors, d)
		l.report(d)
	})

	astPkg, err := parser.ParseDir(
		l.conf.Fset, pkg.Dir, pkg.Path, l.conf.Tag, collector, l.conf.Mode)
	if err != nil && astPkg == nil {
		return fmt.Errorf("%s: %w", pkg.Path, err)
	}
	if astPkg == nil {
		return fmt.Errorf("%s: no package produced", pkg.Path)
	}

	pkg.Name = astPkg.Name
	pkg.Files = astPkg.Files
	return nil
}

// loadImports walks every import declaration in every file of pkg.
//
// A.2.3 ⊢ "there is no aliasing form, no dot-import, and no blank import", so
// an import path maps to exactly one package and there is nothing to record
// but the mapping.
func (l *loader) loadImports(pkg *Package) error {
	seen := make(map[string]bool)

	for _, f := range pkg.Files {
		for _, decl := range f.Imports {
			for _, lit := range decl.Paths {
				path, ok := unquote(lit.Value)
				if !ok || seen[path] {
					continue
				}
				seen[path] = true

				dep, err := l.load(path, pkg)
				if err != nil {
					var ce *CycleError
					if asCycle(err, &ce) {
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

	collector := diag.ReporterFunc(func(d *diag.Diagnostic) {
		pkg.Errors = append(pkg.Errors, d)
		l.report(d)
	})

	conf := &analyzer.Config{
		Fset:     l.conf.Fset,
		Tag:      l.conf.Tag,
		Reporter: collector,
		Importer: &pkgImporter{pkg: pkg},
	}

	name := pkg.Name
	if name == "" {
		name = pkg.Path
	}

	chk := analyzer.NewChecker(conf, pkg.Path, name, pkg.Info)
	tpkg, _ := chk.Files(pkg.Files)
	pkg.Types = tpkg

	// A package is complete once its own check finishes, errors included: a
	// partial scope still answers most lookups, and refusing to hand it to a
	// dependent would turn one bad file into a cascade across the graph.
	if tpkg != nil {
		tpkg.MarkComplete()
	}
}

// pkgImporter satisfies analyzer.Importer for one package, over the
// dependencies the loader already resolved.
//
// It is deliberately not the whole graph. A package may only reach what it
// declares an import of — resolving a path this package never imported would
// let a name leak in through a transitive dependency, and A.2.3's qualifier
// rule gives no spelling for that.
type pkgImporter struct{ pkg *Package }

func (i *pkgImporter) Import(path string) (*types.Package, error) {
	dep, ok := i.pkg.Imports[path]
	if !ok || dep.Types == nil {
		return nil, fmt.Errorf("package %q is not imported by %s", path, i.pkg.Path)
	}
	return dep.Types, nil
}

// ---------------------------------------------------------------- errors

// CycleError reports an import cycle. It is fatal rather than a diagnostic:
// Vertex has no forward-declaration form, so there is no file to point at as
// the one that should have broken it.
type CycleError struct {
	Path     string
	Importer string
	Stack    []string
}

func (e *CycleError) Error() string {
	s := "import cycle: "
	for _, p := range e.Stack {
		s += p + " -> "
	}
	return s + e.Path
}

func asCycle(err error, out **CycleError) bool {
	ce, ok := err.(*CycleError)
	if ok {
		*out = ce
	}
	return ok
}

func dirOf(p *Package) string {
	if p == nil {
		return ""
	}
	return p.Dir
}

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

var _ = sync.Once{} // reserved for the concurrent loader; see resolver.go