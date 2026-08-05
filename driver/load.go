// driver/load.go
package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vertex-language/vertex/analyzer"
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/importer"
	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// Package is one checked compilation unit, in the shape the lowering chain
// wants. It is deliberately the driver's own type rather than
// *importer.Package: the single-file path below never builds one of those,
// and lower.go should have exactly one input shape to convert from.
//
// Path and Dir are carried for diagnostics and verbose output only. hir
// reads neither — hir.Unit takes the *types.Package, and §1.3's qualifier
// is that package's own declared Name, never something derived from a
// locator.
type Package struct {
	Path  string
	Dir   string
	Name  string
	Files []*ast.File
	Types *types.Package
	Info  *types.Info
}

// Loaded is one load's whole result: the packages in dependency-first
// order, the single FileSet their positions live in, and every diagnostic
// collected along the way.
//
// The FileSet is the reason this type exists rather than a bare
// []*Package. hir.Config takes one (for loc lines; hir resolves no
// positions itself) and so does lower/vir.Config, and a compilation has
// exactly one — so it travels with the packages instead of being rebuilt
// or threaded separately.
type Loaded struct {
	Fset     *token.FileSet
	Packages []*Package

	lc *loadContext
}

// Diagnostics is everything reported during the load, in report order.
// The test runner partitions these by test-function extent; an ordinary
// build never looks at them, because Load already turned them into an
// error.
func (l *Loaded) Diagnostics() []*diag.Diagnostic { return l.lc.list.Items() }

// Renderer builds the caret-and-excerpt renderer over this load's sources.
func (l *Loaded) Renderer(color bool) *diag.Renderer { return l.lc.renderer(color) }

// loadContext carries the pieces every load path shares: one FileSet for
// the whole compilation (so a diagnostic can span two files in two
// packages), the collected diagnostics, and the source bytes the renderer
// needs to excerpt.
type loadContext struct {
	fset    *token.FileSet
	list    *diag.List
	sources map[string][]byte
}

func newLoadContext() *loadContext {
	return &loadContext{
		fset:    token.NewFileSet(),
		list:    &diag.List{},
		sources: map[string][]byte{},
	}
}

// Report satisfies diag.Reporter. Collecting rather than streaming is right
// for a batch compile: the list is sorted and deduped once at the end, so
// parser recovery re-reporting at a resync point costs one entry, not two.
func (lc *loadContext) Report(d *diag.Diagnostic) { lc.list.Report(d) }

func (lc *loadContext) renderer(color bool) *diag.Renderer {
	return &diag.Renderer{
		Fset:  lc.fset,
		Color: color,
		Source: func(f *token.File) []byte {
			if b, ok := lc.sources[f.Name()]; ok {
				return b
			}
			// Not cached (a dependency parsed by the importer from disk):
			// read it now. A failure just means no source excerpt, which
			// the renderer already handles.
			b, err := os.ReadFile(f.Name())
			if err != nil {
				return nil
			}
			lc.sources[f.Name()] = b
			return b
		},
	}
}

// diagnosticsError renders every collected diagnostic and returns a plain
// error summarizing the count. The rendered text goes to stderr because
// that's where a caret-and-excerpt display belongs; the error value stays
// short, since cli/ prefixes it with "vertex: " and prints one line.
func (lc *loadContext) diagnosticsError(opts *Options) error {
	lc.list.Sort()
	lc.list.Dedup()
	if lc.list.NumErrors() == 0 {
		// Warnings still render; they just don't fail the build.
		if lc.list.Len() > 0 {
			_ = lc.renderer(false).RenderList(opts.Stderr, lc.list)
		}
		return nil
	}
	_ = lc.renderer(false).RenderList(opts.Stderr, lc.list)
	n := lc.list.NumErrors()
	if n == 1 {
		return fmt.Errorf("1 error")
	}
	return fmt.Errorf("%d errors", n)
}

// Load resolves, parses, and checks whatever Options.Input names, renders
// any diagnostics, and fails if any of them was an error. This is the
// entry point for a build.
//
// The load itself (see load, below) is deliberately *not* the thing that
// decides a diagnostic is fatal: importer's own contract is "partial
// results always", and the test runner needs exactly that — an
// Expected(error) test is supposed to fail to check, and a runner handed
// nil packages could never discover it. So collection and judgement are
// two steps, and only this one judges.
func Load(opts *Options, tag token.BuildTag) (*Loaded, error) {
	ld, err := load(opts, tag)
	if err != nil {
		// A load-level failure (an unresolvable import, a cycle) still
		// wants any diagnostics collected so far rendered first: they're
		// usually what explains it.
		if ld != nil {
			if derr := ld.lc.diagnosticsError(opts); derr != nil {
				return nil, fmt.Errorf("%v (%v)", err, derr)
			}
		}
		return nil, err
	}
	if derr := ld.lc.diagnosticsError(opts); derr != nil {
		return nil, derr
	}
	return ld, nil
}

// load is Load minus the judgement: it returns whatever finished, with
// every diagnostic on the returned *Loaded, and errors only on what has no
// file to point at as the fix — a missing input, an unresolvable import
// path, a cycle, a build tag that excludes everything.
//
// Two paths, because a file and a directory are genuinely different
// requests:
//
//   - A directory is a package. importer.Load does the whole job: resolve,
//     parse (with §1.3's two-pass build-tag filter), check, in
//     dependency-first order.
//   - A single .vs file is *that file*, not its directory. importer has no
//     spelling for that — parser.ParseDir is directory-granular by
//     construction — so this path parses the one file, wraps it in an
//     ast.Package, and runs analyzer directly, with an Importer that
//     delegates each declared import back to importer.Load. Compiling
//     `main.vs` must not silently pull in a sibling `scratch.vs`, which
//     is exactly what reusing the directory path would do.
func load(opts *Options, tag token.BuildTag) (*Loaded, error) {
	opts.defaults()
	lc := newLoadContext()
	ld := &Loaded{Fset: lc.fset, lc: lc}

	info, err := os.Stat(opts.Input)
	if err != nil {
		return ld, fmt.Errorf("reading %s: %w", opts.Input, err)
	}

	var pkgs []*Package
	if info.IsDir() {
		pkgs, err = loadDir(opts, lc, tag, opts.Input)
	} else {
		if !strings.HasSuffix(opts.Input, ".vs") {
			return ld, fmt.Errorf("%s is not a .vs source file or a package directory", opts.Input)
		}
		pkgs, err = loadFile(opts, lc, tag, opts.Input)
	}
	ld.Packages = pkgs
	if err != nil {
		return ld, err
	}
	if len(pkgs) == 0 {
		return ld, fmt.Errorf(
			"%s holds no package for build tag %q — every file's `build` clause selects another target",
			opts.Input, tag)
	}
	return ld, nil
}

// resolverFor builds the import-path resolver: the packages root (from
// -packages-dir, $VERTEX_PATH, or ~/.vertex/packages, in that order),
// plus the current directory so a relative import resolves against the
// tree being built.
func resolverFor(opts *Options) (importer.Resolver, error) {
	roots := []string{}
	switch {
	case opts.PackagesDir != "":
		roots = append(roots, opts.PackagesDir)
	case os.Getenv("VERTEX_PATH") != "":
		roots = append(roots, os.Getenv("VERTEX_PATH"))
	default:
		home, err := os.UserHomeDir()
		if err == nil {
			roots = append(roots, filepath.Join(home, ".vertex", "packages"))
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("resolving the current directory: %w", err)
	}
	roots = append(roots, cwd)
	return importer.NewDirResolver(roots...)
}

func importerConfig(opts *Options, lc *loadContext, tag token.BuildTag) (*importer.Config, error) {
	res, err := resolverFor(opts)
	if err != nil {
		return nil, err
	}
	return &importer.Config{
		Fset:     lc.fset,
		Tag:      tag,
		Resolver: res,
		Reporter: lc,
	}, nil
}

// loadDir loads a directory as a package via importer, mapping the given
// directory to a synthetic root path so a build of `./cmd/app` doesn't
// require that directory to already be reachable from a packages root.
//
// A check error is not a failure here. importer marks every package
// complete unconditionally and keeps partial results on purpose; returning
// them is what lets `Load` render one summarizing error and lets the test
// runner see the tree an Expected(error) test lives in.
func loadDir(opts *Options, lc *loadContext, tag token.BuildTag, dir string) ([]*Package, error) {
	conf, err := importerConfig(opts, lc, tag)
	if err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	rootPath := filepath.ToSlash(filepath.Clean(dir))

	// The root's own path resolves to the directory the caller named; every
	// other path falls through to the ordinary resolver. A locator is not a
	// name (§1.3), so inventing one for the root is exactly the kind of
	// thing a build system is expected to decide.
	base := conf.Resolver
	conf.Resolver = importer.ResolverFunc(func(path, from string) (string, error) {
		if path == rootPath {
			return abs, nil
		}
		return base.Resolve(path, from)
	})

	res, err := importer.Load(conf, rootPath)
	if res == nil {
		return nil, err
	}
	return fromImporterResult(res), err
}

func fromImporterResult(res *importer.Result) []*Package {
	out := make([]*Package, 0, len(res.Order))
	for _, p := range res.Order {
		out = append(out, &Package{
			Path:  p.Path,
			Dir:   p.Dir,
			Name:  p.Name,
			Files: p.Files,
			Types: p.Types,
			Info:  p.Info,
		})
	}
	return out
}

// loadFile is the single-file path described in load's doc comment.
func loadFile(opts *Options, lc *loadContext, tag token.BuildTag, path string) ([]*Package, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	lc.sources[abs] = src

	// ParseFile always returns a non-nil *ast.File — recovery produces
	// Bad* nodes rather than nothing — so a nil here means something
	// structurally impossible happened, not a syntax error.
	file, _ := parser.ParseFile(lc.fset, abs, src, lc, 0)
	if file == nil {
		return nil, fmt.Errorf("%s: could not be parsed", path)
	}

	// §1.3: a file whose build clause names another target is excluded
	// whole. Naming the file explicitly doesn't override that — it makes it
	// an error, since the caller clearly meant to build this file.
	if bt := file.BuildTag(); bt != token.TagNone && bt != tag {
		return nil, fmt.Errorf(
			"%s carries `build %s` but the target selects %q — a file whose tag doesn't match "+
				"is excluded from the build entirely", path, bt, tag)
	}

	// ast.NewPackage is a validated container and nothing more: it checks
	// one file minimum, package-clause agreement, and tag agreement. The
	// locator is synthesized from the file's own directory, the same shape
	// loadDir invents for its root.
	dir := filepath.Dir(abs)
	importPath := filepath.ToSlash(filepath.Clean(dir))
	astPkg, err := ast.NewPackage(lc.fset, dir, importPath, tag, []*ast.File{file})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	conf, err := importerConfig(opts, lc, tag)
	if err != nil {
		return nil, err
	}
	imp := &fileImporter{conf: conf, cache: map[string]*types.Package{}, seen: map[string]bool{}}

	acheck := &analyzer.Config{
		Fset:     lc.fset,
		Tag:      tag,
		Reporter: lc,
		Importer: imp,
	}
	info := types.NewInfo()
	checker := analyzer.NewChecker(acheck, astPkg.Path, astPkg.Name, info)

	// Errors accumulate rather than abort, and Package() is valid after
	// them — so a failed check still yields a scope worth handing on. The
	// diagnostics are already on lc; Load summarizes them.
	tpkg, _ := checker.Files(astPkg.Files)
	if tpkg == nil {
		tpkg = checker.Package()
	}

	// An import failure is the other kind: no file to point at, nothing
	// downstream can be trusted.
	if imp.err != nil {
		return nil, imp.err
	}

	// Dependencies first — the same order importer.Result.Order guarantees,
	// which is what lower/hir and vvm's importer both expect.
	out := append([]*Package(nil), imp.deps...)
	return append(out, &Package{
		Path:  astPkg.Path,
		Dir:   astPkg.Dir,
		Name:  astPkg.Name,
		Files: astPkg.Files,
		Types: tpkg,
		Info:  info,
	}), nil
}

// fileImporter satisfies analyzer.Importer for the single-file path by
// running a full importer.Load per declared import path, and remembering
// every package that came back so the lowering chain gets the whole graph,
// not just the one file's own package.
//
// Loads are cached by path, so a diamond costs one load, and the recorded
// dependency list stays topologically ordered because each Load's own
// Result.Order already is.
type fileImporter struct {
	conf  *importer.Config
	cache map[string]*types.Package
	seen  map[string]bool
	deps  []*Package
	err   error
}

func (f *fileImporter) Import(path string) (*types.Package, error) {
	if p, ok := f.cache[path]; ok {
		return p, nil
	}
	res, err := importer.Load(f.conf, path)
	if err != nil {
		f.err = err
		return nil, err
	}
	for _, p := range res.Order {
		if p.Types == nil {
			continue
		}
		f.cache[p.Path] = p.Types
		if f.seen[p.Path] {
			continue
		}
		f.seen[p.Path] = true
		f.deps = append(f.deps, &Package{
			Path: p.Path, Dir: p.Dir, Name: p.Name,
			Files: p.Files, Types: p.Types, Info: p.Info,
		})
	}
	p, ok := f.cache[path]
	if !ok {
		err := fmt.Errorf("import %q resolved but produced no checked package", path)
		f.err = err
		return nil, err
	}
	return p, nil
}

// sortedPaths is a small helper for deterministic verbose output.
func sortedPaths(pkgs []*Package) []string {
	out := make([]string, 0, len(pkgs))
	for _, p := range pkgs {
		out = append(out, p.Path)
	}
	sort.Strings(out)
	return out
}