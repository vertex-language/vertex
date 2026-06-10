// imports.go
package compiler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	cir "github.com/vertex-language/ir/c"
	"github.com/vertex-language/vertex/parser"
)

// ObjectFunc compiles a CIR module to target-specific object-file bytes.
// Provided by cmd/vertex/main.go so the compiler package stays backend-agnostic.
type ObjectFunc func(*cir.Module) ([]byte, error)

// ─────────────────────────────────────────────────────────────────────────────
// CompiledPackage
// ─────────────────────────────────────────────────────────────────────────────

// CompiledPackage holds the type-checking scope and the compiled object bytes
// for one Vertex package.
type CompiledPackage struct {
	ImportPath string
	Name       string
	Scope      *Scope // exported symbols for type-checking
	ObjBytes   []byte // compiled .o; nil if objectFunc was nil or compilation failed
}

// ExportTo copies all symbols into dst under both bare and qualified forms.
func (p *CompiledPackage) ExportTo(dst *Scope) {
	for name, sym := range p.Scope.Symbols() {
		dst.Define(sym)
		dst.Define(&Symbol{
			Name:    p.Name + "." + name,
			Kind:    sym.Kind,
			Type:    sym.Type,
			Decl:    sym.Decl,
			IsConst: sym.IsConst,
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// PackageLoader
// ─────────────────────────────────────────────────────────────────────────────

// PackageLoader resolves, compiles, and caches Vertex packages from the
// central packages directory (VERTEX_PATH equivalent).
type PackageLoader struct {
	packagesDir string
	target      cir.Target
	objectFunc  ObjectFunc                  // nil → symbol extraction only
	cache       map[string]*CompiledPackage // keyed by import path
	cacheDir    string                      // packagesDir/.cache
}

func NewPackageLoader(packagesDir string, target cir.Target, objectFunc ObjectFunc) *PackageLoader {
	cacheDir := ""
	if packagesDir != "" {
		cacheDir = filepath.Join(packagesDir, ".cache")
	}
	return &PackageLoader{
		packagesDir: packagesDir,
		target:      target,
		objectFunc:  objectFunc,
		cache:       make(map[string]*CompiledPackage),
		cacheDir:    cacheDir,
	}
}

// LoadImports resolves and compiles every import declared in a file.
func (l *PackageLoader) LoadImports(imports []*ImportDecl) ([]*CompiledPackage, error) {
	var pkgs []*CompiledPackage
	for _, imp := range imports {
		pkg, err := l.Load(imp.Path)
		if err != nil {
			return nil, fmt.Errorf("loading %q: %w", imp.Path, err)
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

// Load returns the CompiledPackage for importPath, compiling it if needed.
func (l *PackageLoader) Load(importPath string) (*CompiledPackage, error) {
	if pkg, ok := l.cache[importPath]; ok {
		return pkg, nil
	}

	dir, err := l.resolve(importPath)
	if err != nil {
		// Non-fatal: return an empty package so type-checking can continue.
		// The linker will surface any missing symbols later.
		return &CompiledPackage{
			ImportPath: importPath,
			Name:       filepath.Base(importPath),
			Scope:      NewScope(nil),
		}, nil
	}

	pkg, err := l.buildPackage(importPath, dir)
	if err != nil {
		return nil, err
	}
	l.cache[importPath] = pkg
	return pkg, nil
}

// resolve maps importPath to a directory under packagesDir.
func (l *PackageLoader) resolve(importPath string) (string, error) {
	if l.packagesDir == "" {
		return "", fmt.Errorf("no packages directory configured (set VERTEX_PATH or --packages-dir)")
	}
	candidate := filepath.Join(l.packagesDir, filepath.FromSlash(importPath))
	info, err := os.Stat(candidate)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("package %q not found in %s", importPath, l.packagesDir)
	}
	return candidate, nil
}

// buildPackage parses symbols and compiles all .vs files in dir.
func (l *PackageLoader) buildPackage(importPath, dir string) (*CompiledPackage, error) {
	srcPaths, sources, err := collectSources(dir)
	if err != nil {
		return nil, err
	}

	diags := NewDiagnostics()
	pkgName := filepath.Base(importPath)
	scope := NewScope(nil)

	for i, src := range sources {
		file := parseFileQuick(src, srcPaths[i], diags)
		if file == nil {
			continue
		}
		if file.Package != "" {
			pkgName = file.Package
		}
		collectFileSymbols(file, scope)
	}

	pkg := &CompiledPackage{
		ImportPath: importPath,
		Name:       pkgName,
		Scope:      scope,
	}

	if l.objectFunc != nil && len(sources) > 0 {
		obj, err := l.compileToObject(importPath, srcPaths, sources, pkgName)
		if err != nil {
			return nil, fmt.Errorf("package %q failed to compile:\n%w", importPath, err)
		}
		pkg.ObjBytes = obj
	}

	return pkg, nil
}

// compileToObject drives the parse → resolve → lower → encode pipeline for a
// package, using a content-hash disk cache to skip redundant work.
func (l *PackageLoader) compileToObject(importPath string, paths, sources []string, pkgName string) ([]byte, error) {
	// contentHash includes the target so that cross-compiling for different
	// platforms never returns a cached object built for a different ABI, and
	// so that a compiler update that changes code generation (e.g. the SizeOf
	// struct fix) is reflected in a new hash once the cache directory is cleared.
	hash := contentHash(sources, l.target)
	if cached, err := l.readCache(hash); err == nil {
		return cached, nil
	}

	mod, err := l.buildModule(importPath, paths, sources, pkgName)
	if err != nil {
		return nil, err
	}

	obj, err := l.objectFunc(mod)
	if err != nil {
		return nil, fmt.Errorf("encoding package %q: %w", importPath, err)
	}

	l.writeCache(hash, obj)
	return obj, nil
}

// buildModule is the shared parse → resolve → lower pipeline.
func (l *PackageLoader) buildModule(importPath string, paths, sources []string, pkgName string) (*cir.Module, error) {
	combined := strings.Join(sources, "\n\n")
	diags := NewDiagnostics()

	file := parseFileQuick(combined, paths[0], diags)
	if file == nil || diags.HasErrors() {
		return nil, diags.Error()
	}
	if file.Package == "" {
		file.Package = pkgName
	}

	global := newGlobalScope()
	resolver := NewResolver(diags, global)
	resolver.ResolveFile(file)
	if diags.HasErrors() {
		return nil, diags.Error()
	}

	mod := cir.NewModule(pkgName)
	mod.BindTarget(l.target)

	lowerer := NewLowerer(diags, mod)
	lowerer.LowerFile(file)
	if diags.HasErrors() {
		return nil, diags.Error()
	}

	mod.Optimize(cir.ConstantFold, cir.DeadCodeElim, cir.StrengthReduce)
	return mod, nil
}

// ── Cache ─────────────────────────────────────────────────────────────────────

func (l *PackageLoader) readCache(hash string) ([]byte, error) {
	if l.cacheDir == "" {
		return nil, fmt.Errorf("no cache dir")
	}
	return os.ReadFile(filepath.Join(l.cacheDir, hash+".o"))
}

func (l *PackageLoader) writeCache(hash string, data []byte) {
	if l.cacheDir == "" {
		return
	}
	_ = os.MkdirAll(l.cacheDir, 0o755)
	_ = os.WriteFile(filepath.Join(l.cacheDir, hash+".o"), data, 0o644)
}

// contentHash returns a short fingerprint of the package source files combined
// with the compilation target. Including the target ensures that object files
// cached for linux-amd64 are never served to a windows-amd64 build, and that
// clearing the cache after a compiler update forces a clean rebuild.
func contentHash(sources []string, target cir.Target) string {
	h := sha256.New()
	// Target string first so a target change always produces a different prefix.
	h.Write([]byte(target.String()))
	for _, s := range sources {
		h.Write([]byte(s))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// ── Source helpers ────────────────────────────────────────────────────────────

// collectSources returns all .vs paths and their contents from dir.
func collectSources(dir string) (paths, sources []string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".vs") {
			continue
		}
		p := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, nil, fmt.Errorf("reading %s: %w", p, err)
		}
		paths = append(paths, p)
		sources = append(sources, string(data))
	}
	return paths, sources, nil
}

// parseFileQuick parses a single Vertex source file and returns the AST.
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
	return newASTBuilder(filename, diags).BuildFile(tree)
}

// collectFileSymbols adds top-level symbol names from file into scope.
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