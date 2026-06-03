package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/vertex-language/vertex/parser"
)

// ─────────────────────────────────────────────────────────────────────────────
// Package — a compiled Vertex package (one directory)
// ─────────────────────────────────────────────────────────────────────────────

// Package represents a resolved Vertex package.  For this release we parse the
// files and extract symbol names for type-checking; we do NOT emit ir/c for
// imported packages (assume they are linked separately).
type Package struct {
	ImportPath string
	Name       string // last path segment
	Scope      *Scope // exported symbols
}

// ExportTo copies all symbols from p.Scope into dst under both bare name and
// "pkg.Name" qualified form.
func (p *Package) ExportTo(dst *Scope) {
	for name, sym := range p.Scope.Symbols() {
		dst.Define(sym)
		qualified := &Symbol{
			Name:    p.Name + "." + name,
			Kind:    sym.Kind,
			Type:    sym.Type,
			Decl:    sym.Decl,
			IsConst: sym.IsConst,
		}
		dst.Define(qualified)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// PackageLoader — resolves import paths → Package
// ─────────────────────────────────────────────────────────────────────────────

// PackageLoader finds and parses Vertex packages from a list of search paths.
// Results are cached by import path.
type PackageLoader struct {
	searchPaths []string
	cache       map[string]*Package
}

func NewPackageLoader(searchPaths []string) *PackageLoader {
	return &PackageLoader{
		searchPaths: searchPaths,
		cache:       make(map[string]*Package),
	}
}

// LoadImports resolves all imports declared in a file.
func (l *PackageLoader) LoadImports(imports []*ImportDecl) ([]*Package, error) {
	var pkgs []*Package
	for _, imp := range imports {
		// A grouped import block generates one ImportDecl per path already.
		paths := l.expandImport(imp.Path)
		for _, path := range paths {
			pkg, err := l.Load(path)
			if err != nil {
				return nil, fmt.Errorf("loading %q: %w", path, err)
			}
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs, nil
}

// expandImport handles the case where a single ImportDecl might contain
// a grouped import block (parsed as one path per line). For now it
// returns the path as-is.
func (l *PackageLoader) expandImport(path string) []string {
	// Future: handle multi-path grouped imports.
	return []string{path}
}

// Load returns the Package for the given import path, loading it if necessary.
func (l *PackageLoader) Load(importPath string) (*Package, error) {
	if pkg, ok := l.cache[importPath]; ok {
		return pkg, nil
	}
	dir, err := l.resolve(importPath)
	if err != nil {
		// Non-fatal: emit a warning and return an empty package.
		return &Package{
			ImportPath: importPath,
			Name:       filepath.Base(importPath),
			Scope:      NewScope(nil),
		}, nil
	}
	pkg, err := l.loadDir(importPath, dir)
	if err != nil {
		return nil, err
	}
	l.cache[importPath] = pkg
	return pkg, nil
}

// resolve searches all search paths for the given import path.
func (l *PackageLoader) resolve(importPath string) (string, error) {
	rel := filepath.FromSlash(importPath)
	for _, sp := range l.searchPaths {
		candidate := filepath.Join(sp, rel)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("package %q not found in search paths %v", importPath, l.searchPaths)
}

// loadDir parses all *.vs files in dir and extracts top-level symbols.
func (l *PackageLoader) loadDir(importPath, dir string) (*Package, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}

	pkgName := filepath.Base(importPath)
	scope := NewScope(nil)
	diags := NewDiagnostics()

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".vs") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		file := parseFileQuick(string(data), path, diags)
		if file == nil {
			continue
		}
		// Override package name from the source file declaration.
		if file.Package != "" {
			pkgName = file.Package
		}
		collectFileSymbols(file, scope)
	}

	return &Package{
		ImportPath: importPath,
		Name:       pkgName,
		Scope:      scope,
	}, nil
}

// parseFileQuick parses a single Vertex source file and returns the AST File.
// Parse errors are added to diags but do not stop processing.
func parseFileQuick(src, filename string, diags *Diagnostics) *File {
	input := antlr.NewInputStream(src)
	lexer := parser.NewVertexLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewVertexParser(stream)

	el := &antlrErrorListener{diags: diags, filename: filename}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(el)
	p.RemoveErrorListeners()
	p.AddErrorListener(el)

	tree := p.File()
	if diags.HasErrors() {
		return nil
	}
	builder := newASTBuilder(filename, diags)
	return builder.BuildFile(tree)
}

// collectFileSymbols adds the top-level symbol names from file into scope.
// Types are left as VUnknown; the resolver fills them in for the main file.
func collectFileSymbols(file *File, scope *Scope) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *FuncDecl:
			if d.Receiver == nil {
				scope.Define(&Symbol{Name: d.Name, Kind: SymFunc, Decl: d})
			}
		case *StructDecl:
			scope.Define(&Symbol{
				Name: d.Name, Kind: SymStruct,
				Type: &VStruct{Name: d.Name, Decl: d}, Decl: d,
			})
		case *ClassDecl:
			scope.Define(&Symbol{
				Name: d.Name, Kind: SymClass,
				Type: &VClass{Name: d.Name, Decl: d}, Decl: d,
			})
		case *EnumDecl:
			scope.Define(&Symbol{
				Name: d.Name, Kind: SymEnum,
				Type: &VEnum{Name: d.Name, Decl: d}, Decl: d,
			})
		case *TypeAliasDecl:
			scope.Define(&Symbol{Name: d.Name, Kind: SymTypeAlias, Decl: d})
		case *VarDecl:
			for _, name := range d.Binding.Names {
				scope.Define(&Symbol{Name: name, Kind: SymVar, Decl: d, IsConst: d.IsLet})
			}
		}
	}
}