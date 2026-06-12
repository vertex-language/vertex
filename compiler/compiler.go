// compiler.go
package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	cir "github.com/vertex-language/ir/c"
	"github.com/vertex-language/vertex/parser"
)

// ─────────────────────────────────────────────────────────────────────────────
// Config
// ─────────────────────────────────────────────────────────────────────────────

// ObjectFunc is the backend hook: given a compiled *cir.Module, produce a
// relocatable object file as raw bytes.  A nil ObjectFunc disables package
// compilation (packages are loaded but not re-compiled from source).
type ObjectFunc func(mod *cir.Module) ([]byte, error)

// Config holds the compiler's runtime configuration.
type Config struct {
	Target      cir.Target // platform the CIR module is bound to
	PackagesDir string     // central packages root ($VERTEX_PATH / --packages-dir)
	ObjectFunc  ObjectFunc // nil disables package compilation
	Rebuild     bool       // force cache wipe on package load
}

// ─────────────────────────────────────────────────────────────────────────────
// TestFuncInfo
// ─────────────────────────────────────────────────────────────────────────────

// TestFuncInfo captures the metadata the test runner needs for one test
// function: its name, which output channel to capture, and the expected value.
type TestFuncInfo struct {
	Name     string
	Channel  string // "stdout" | "stderr" | …
	Expected string
}

// ─────────────────────────────────────────────────────────────────────────────
// Compiler
// ─────────────────────────────────────────────────────────────────────────────

type Compiler struct {
	cfg    Config
	diags  *Diagnostics
	loader *PackageLoader
}

func New(cfg Config) *Compiler {
	if cfg.Target == cir.TargetUnknown {
		cfg.Target = cir.TargetLinuxAMD64
	}
	return &Compiler{
		cfg:    cfg,
		diags:  NewDiagnostics(),
		loader: NewPackageLoader(cfg.PackagesDir, cfg.Target, cfg.ObjectFunc, cfg.Rebuild),
	}
}

func (c *Compiler) Diagnostics() *Diagnostics { return c.diags }

// ─────────────────────────────────────────────────────────────────────────────
// CompileFile — single .vs file or package directory
// ─────────────────────────────────────────────────────────────────────────────

// CompileFile compiles path, which may be a .vs source file or a package
// directory.  It returns the CIR module and the compiled object bytes for all
// imported packages that the caller must link.
func (c *Compiler) CompileFile(path string) (*cir.Module, [][]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot access %s: %w", path, err)
	}
	if info.IsDir() {
		return c.CompileDir(path)
	}
	data, err := readFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return c.CompileSource(string(data), path)
}

// CompileDir compiles all .vs files in dir as a single package module.
// Unlike CompileSource, runtime imports are not injected automatically —
// package directories are expected to declare their own dependencies
// explicitly.
func (c *Compiler) CompileDir(dir string) (*cir.Module, [][]byte, error) {
	c.diags.Reset()

	_, sources, err := collectSources(dir)
	if err != nil {
		return nil, nil, err
	}
	if len(sources) == 0 {
		return nil, nil, fmt.Errorf("no .vs files found in %s", dir)
	}

	pkgName := filepath.Base(dir)
	combined := strings.Join(sources, "\n\n")
	virtualFile := filepath.Join(dir, pkgName+".vs")

	file, err := c.parseSource(combined, virtualFile)
	if err != nil {
		return nil, nil, err
	}

	pkgs, loadErr := c.loader.LoadImports(file.Imports)
	if loadErr != nil {
		c.diags.Errorf(Pos{File: virtualFile}, "import error: %v", loadErr)
	}

	pkgObjs := collectObjBytes(pkgs)

	global := newGlobalScope()
	for _, pkg := range pkgs {
		pkg.ExportTo(global)
	}

	resolver := NewResolver(c.diags, global)
	resolver.ResolveFile(file)
	if c.diags.HasErrors() {
		return nil, nil, c.diags.Error()
	}

	mod := cir.NewModule(pkgName)
	mod.BindTarget(c.cfg.Target)

	lowerer := NewLowerer(c.diags, mod)
	lowerer.LowerFile(file)
	if c.diags.HasErrors() {
		return nil, nil, c.diags.Error()
	}

	mod.Optimize(cir.ConstantFold, cir.DeadCodeElim, cir.StrengthReduce)
	return mod, pkgObjs, nil
}

// CompileSource compiles src (identified by filename for diagnostics).
// It returns the CIR module together with compiled object bytes for all
// imports.
func (c *Compiler) CompileSource(src, filename string) (*cir.Module, [][]byte, error) {
	c.diags.Reset()

	file, err := c.parseSource(src, filename)
	if err != nil {
		return nil, nil, err
	}

	injectRuntimeImports(file)

	pkgs, loadErr := c.loader.LoadImports(file.Imports)
	if loadErr != nil {
		c.diags.Errorf(Pos{File: filename}, "import error: %v", loadErr)
	}

	pkgObjs := collectObjBytes(pkgs)

	global := newGlobalScope()
	for _, pkg := range pkgs {
		pkg.ExportTo(global)
	}

	resolver := NewResolver(c.diags, global)
	resolver.ResolveFile(file)
	if c.diags.HasErrors() {
		return nil, nil, c.diags.Error()
	}

	modName := file.Package
	if modName == "" {
		modName = stripExt(filepath.Base(filename))
	}
	mod := cir.NewModule(modName)
	mod.BindTarget(c.cfg.Target)

	lowerer := NewLowerer(c.diags, mod)
	lowerer.LowerFile(file)
	if c.diags.HasErrors() {
		return nil, nil, c.diags.Error()
	}

	mod.Optimize(cir.ConstantFold, cir.DeadCodeElim, cir.StrengthReduce)
	return mod, pkgObjs, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Test compilation
// ─────────────────────────────────────────────────────────────────────────────

// CompileTestFile parses a test file and returns one *cir.Module per test
// function, the shared package object bytes for linking, and test metadata.
func (c *Compiler) CompileTestFile(path string) ([]TestFuncInfo, []*cir.Module, [][]byte, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return c.compileTestSource(string(data), path)
}

func (c *Compiler) compileTestSource(src, filename string) ([]TestFuncInfo, []*cir.Module, [][]byte, error) {
	c.diags.Reset()

	file, err := c.parseSource(src, filename)
	if err != nil {
		return nil, nil, nil, err
	}

	injectRuntimeImports(file)

	// Collect test function metadata before loading imports so we can bail
	// early if there are no tests.
	var infos []TestFuncInfo
	for _, decl := range file.Decls {
		fn, ok := decl.(*FuncDecl)
		if !ok || fn.Qualifier != FuncQualTest {
			continue
		}
		info := TestFuncInfo{Name: fn.Name, Channel: "stdout"}
		if exp, ok2 := fn.RetType.(*ExpectedTypeExpr); ok2 {
			info.Channel = exp.Channel
			info.Expected = exp.Value
		}
		infos = append(infos, info)
	}
	if len(infos) == 0 {
		return nil, nil, nil, nil
	}

	pkgs, loadErr := c.loader.LoadImports(file.Imports)
	if loadErr != nil {
		c.diags.Errorf(Pos{File: filename}, "import error: %v", loadErr)
	}

	pkgObjs := collectObjBytes(pkgs)

	global := newGlobalScope()
	for _, pkg := range pkgs {
		pkg.ExportTo(global)
	}

	resolver := NewResolver(c.diags, global)
	resolver.ResolveFile(file)
	if c.diags.HasErrors() {
		return nil, nil, nil, c.diags.Error()
	}

	modBase := stripExt(filepath.Base(filename))
	modules := make([]*cir.Module, 0, len(infos))

	for _, info := range infos {
		c.diags.Reset()

		mod := cir.NewModule(modBase + "_" + info.Name)
		mod.BindTarget(c.cfg.Target)

		lwr := NewLowerer(c.diags, mod)
		lwr.testEntryFunc = info.Name
		lwr.LowerFile(file)
		if c.diags.HasErrors() {
			return nil, nil, nil, c.diags.Error()
		}

		mod.Optimize(cir.ConstantFold, cir.DeadCodeElim, cir.StrengthReduce)
		modules = append(modules, mod)
	}

	return infos, modules, pkgObjs, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────────────────────

func (c *Compiler) parseSource(src, filename string) (*File, error) {
	input := antlr.NewInputStream(src)
	lexer := parser.NewVertexLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewVertexParser(stream)

	el := &antlrErrorListener{diags: c.diags, filename: filename}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(el)
	p.RemoveErrorListeners()
	p.AddErrorListener(el)

	tree := p.File()
	if c.diags.HasErrors() {
		return nil, c.diags.Error()
	}
	return newASTBuilder(filename, c.diags).BuildFile(tree), nil
}

// collectObjBytes extracts the compiled object bytes from a loaded package
// slice, skipping packages that were loaded from cache without re-compilation.
func collectObjBytes(pkgs []*CompiledPackage) [][]byte {
	var out [][]byte
	for _, pkg := range pkgs {
		if pkg.ObjBytes != nil {
			out = append(out, pkg.ObjBytes)
		}
	}
	return out
}

func stripExt(name string) string {
	if ext := filepath.Ext(name); ext != "" {
		return name[:len(name)-len(ext)]
	}
	return name
}

// injectRuntimeImports prepends the core runtime packages that every user
// compilation unit implicitly depends on.  Safe to call multiple times.
func injectRuntimeImports(file *File) {
	for _, path := range []string{"arrays"} {
		found := false
		for _, imp := range file.Imports {
			if imp.Path == path {
				found = true
				break
			}
		}
		if !found {
			file.Imports = append([]*ImportDecl{{Path: path}}, file.Imports...)
		}
	}
}

func readFile(path string) ([]byte, error) { return osReadFile(path) }

// ─────────────────────────────────────────────────────────────────────────────
// ANTLR error listener
// ─────────────────────────────────────────────────────────────────────────────

type antlrErrorListener struct {
	*antlr.DefaultErrorListener
	diags    *Diagnostics
	filename string
}

func (l *antlrErrorListener) SyntaxError(
	_ antlr.Recognizer,
	_ interface{},
	line, column int,
	msg string,
	_ antlr.RecognitionException,
) {
	l.diags.Errorf(
		Pos{File: l.filename, Line: line, Column: column + 1},
		"syntax: %s", msg,
	)
}