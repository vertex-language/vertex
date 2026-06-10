package compiler

import (
	"fmt"
	"path/filepath"

	"github.com/antlr4-go/antlr/v4"
	cir "github.com/vertex-language/ir/c"
	"github.com/vertex-language/vertex/parser"
)

// ─────────────────────────────────────────────────────────────────────────────
// Config
// ─────────────────────────────────────────────────────────────────────────────

// Config holds the compiler's runtime configuration.
type Config struct {
	Target      cir.Target
	SearchPaths []string
	// StdlibPath is the directory that contains the Vertex runtime packages
	// (e.g. the folder that holds runtime/arrays). When set it is prepended to
	// SearchPaths so built-in packages are always found before user paths.
	StdlibPath string
}

// ─────────────────────────────────────────────────────────────────────────────
// TestFuncInfo — metadata extracted from a test function's Expected annotation
// ─────────────────────────────────────────────────────────────────────────────

// TestFuncInfo describes one test-qualified function and what it expects.
type TestFuncInfo struct {
	Name     string // function name, e.g. "test_add"
	Channel  string // "stdout" | "exitCode"
	Expected string // expected output string, e.g. "15"
}

// ─────────────────────────────────────────────────────────────────────────────
// Compiler
// ─────────────────────────────────────────────────────────────────────────────

// Compiler orchestrates the frontend pipeline for a single source file.
//
// Pipeline:
//
//	Source
//	  └─ Parse (ANTLR)       → *File AST
//	       └─ Import loading  → package scopes
//	            └─ Resolve    → typed *File AST
//	                 └─ Lower → ir/c module
type Compiler struct {
	cfg    Config
	diags  *Diagnostics
	loader *PackageLoader
}

func New(cfg Config) *Compiler {
	if cfg.Target == cir.TargetUnknown {
		cfg.Target = cir.TargetLinuxAMD64
	}
	searchPaths := cfg.SearchPaths
	if cfg.StdlibPath != "" {
		searchPaths = append([]string{cfg.StdlibPath}, searchPaths...)
	}
	return &Compiler{
		cfg:    cfg,
		diags:  NewDiagnostics(),
		loader: NewPackageLoader(searchPaths),
	}
}

func (c *Compiler) Diagnostics() *Diagnostics { return c.diags }

// CompileFile reads path and compiles it, returning the optimised C IR module.
func (c *Compiler) CompileFile(path string) (*cir.Module, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return c.CompileSource(string(data), path)
}

// CompileSource compiles src (identified by filename for diagnostics) and
// returns the optimised C IR module.
func (c *Compiler) CompileSource(src, filename string) (*cir.Module, error) {
	c.diags.Reset()

	file, err := c.parseSource(src, filename)
	if err != nil {
		return nil, err
	}

	injectRuntimeImports(file)

	pkgs, err := c.loader.LoadImports(file.Imports)
	if err != nil {
		c.diags.Errorf(Pos{File: filename}, "import error: %v", err)
	}

	global := newGlobalScope()
	for _, pkg := range pkgs {
		pkg.ExportTo(global)
	}

	resolver := NewResolver(c.diags, global)
	resolver.ResolveFile(file)
	if c.diags.HasErrors() {
		return nil, c.diags.Error()
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
		return nil, c.diags.Error()
	}

	mod.Optimize(cir.ConstantFold, cir.DeadCodeElim, cir.StrengthReduce)
	return mod, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Test compilation
// ─────────────────────────────────────────────────────────────────────────────

// CompileTestFile parses a "build test" file and returns one *cir.Module per
// test function. Each module is compiled as a standalone program with its own
// main() that calls the single test function and prints its return value.
// The returned slices are parallel: infos[i] describes modules[i].
func (c *Compiler) CompileTestFile(path string) ([]TestFuncInfo, []*cir.Module, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return c.compileTestSource(string(data), path)
}

func (c *Compiler) compileTestSource(src, filename string) ([]TestFuncInfo, []*cir.Module, error) {
	c.diags.Reset()

	file, err := c.parseSource(src, filename)
	if err != nil {
		return nil, nil, err
	}

	injectRuntimeImports(file)

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
		return nil, nil, nil
	}

	pkgs, err := c.loader.LoadImports(file.Imports)
	if err != nil {
		c.diags.Errorf(Pos{File: filename}, "import error: %v", err)
	}
	global := newGlobalScope()
	for _, pkg := range pkgs {
		pkg.ExportTo(global)
	}

	resolver := NewResolver(c.diags, global)
	resolver.ResolveFile(file)
	if c.diags.HasErrors() {
		return nil, nil, c.diags.Error()
	}

	modBase := stripExt(filepath.Base(filename))
	var modules []*cir.Module

	for _, info := range infos {
		c.diags.Reset()

		mod := cir.NewModule(modBase + "_" + info.Name)
		mod.BindTarget(c.cfg.Target)

		lwr := NewLowerer(c.diags, mod)
		lwr.testEntryFunc = info.Name
		lwr.LowerFile(file)
		if c.diags.HasErrors() {
			return nil, nil, c.diags.Error()
		}

		mod.Optimize(cir.ConstantFold, cir.DeadCodeElim, cir.StrengthReduce)
		modules = append(modules, mod)
	}

	return infos, modules, nil
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
	b := newASTBuilder(filename, c.diags)
	return b.BuildFile(tree), nil
}

func stripExt(name string) string {
	ext := filepath.Ext(name)
	if ext != "" {
		return name[:len(name)-len(ext)]
	}
	return name
}

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

func readFile(path string) ([]byte, error) {
	return osReadFile(path)
}

// injectRuntimeImports prepends the core runtime packages that every
// compilation unit implicitly depends on. Safe to call multiple times —
// it skips paths already present in the import list.
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