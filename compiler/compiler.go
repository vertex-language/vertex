package compiler

import (
	"fmt"
	"path/filepath"

	"github.com/antlr4-go/antlr/v4"
	"github.com/vertex-language/ir/c"
	"github.com/vertex-language/vertex/parser"
)

// ─────────────────────────────────────────────────────────────────────────────
// Config
// ─────────────────────────────────────────────────────────────────────────────

// Config holds the compiler's runtime configuration.
type Config struct {
	// Target is the platform the emitted C or MIR will be compiled for.
	// Defaults to TargetLinuxAMD64.
	Target cir.Target

	// SearchPaths is an ordered list of root directories searched when
	// resolving import paths.  The compiler checks each directory for a
	// sub-directory whose path matches the import string.
	//
	// Example: if SearchPaths = ["/usr/share/vertex/lib", "./vendor"] and the
	// source imports "core/io", the loader checks
	//   /usr/share/vertex/lib/core/io/
	//   ./vendor/core/io/
	SearchPaths []string
}

// ─────────────────────────────────────────────────────────────────────────────
// Compiler
// ─────────────────────────────────────────────────────────────────────────────

// Compiler orchestrates the frontend pipeline for a single source file or string.
//
// Pipeline:
//
//	Source
//	  └─ Parse (ANTLR)       → *File AST
//	       └─ Import loading  → package scopes
//	            └─ Resolve    → typed *File AST
//	                 └─ Lower → ir/c module (returned for MIR/C emission)
type Compiler struct {
	cfg    Config
	diags  *Diagnostics
	loader *PackageLoader
}

// New creates a Compiler with the given Config.
func New(cfg Config) *Compiler {
	if cfg.Target == cir.TargetUnknown {
		cfg.Target = cir.TargetLinuxAMD64
	}
	return &Compiler{
		cfg:    cfg,
		diags:  NewDiagnostics(),
		loader: NewPackageLoader(cfg.SearchPaths),
	}
}

// Diagnostics exposes the compiler's accumulated messages.
func (c *Compiler) Diagnostics() *Diagnostics { return c.diags }

// CompileFile reads path and compiles it, returning the optimized C IR module.
func (c *Compiler) CompileFile(path string) (*cir.Module, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return c.CompileSource(string(data), path)
}

// CompileSource compiles the Vertex source string src (with the given filename
// for diagnostic messages) and returns the optimized C IR module.
func (c *Compiler) CompileSource(src, filename string) (*cir.Module, error) {
	c.diags.Reset()

	// ── 1. Parse ─────────────────────────────────────────────────────────────
	file, err := c.parseSource(src, filename)
	if err != nil {
		return nil, err
	}

	// ── 2. Load imports ───────────────────────────────────────────────────────
	pkgs, err := c.loader.LoadImports(file.Imports)
	if err != nil {
		c.diags.Errorf(Pos{File: filename}, "import error: %v", err)
	}

	// ── 3. Build package scope (global + imported) ────────────────────────────
	global := newGlobalScope()
	for _, pkg := range pkgs {
		pkg.ExportTo(global)
	}

	// ── 4. Resolve (type-check) ───────────────────────────────────────────────
	resolver := NewResolver(c.diags, global)
	resolver.ResolveFile(file)
	if c.diags.HasErrors() {
		return nil, c.diags.Error()
	}

	// ── 5. Lower to ir/c ──────────────────────────────────────────────────────
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

	// ── 6. Optimise ───────────────────────────────────────────────────────────
	mod.Optimize(cir.ConstantFold, cir.DeadCodeElim, cir.StrengthReduce)

	// ── 7. Return Module ──────────────────────────────────────────────────────
	// C emission is delegated to the caller so the MIR pipeline can use the IR.
	return mod, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────────────────────

// parseSource runs the ANTLR lexer+parser over src and builds the raw AST.
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

// antlrErrorListener forwards ANTLR parse errors into our Diagnostics.
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

// readFile is a thin wrapper so imports.go can call it without importing os.
func readFile(path string) ([]byte, error) {
	import_os_ReadFile := func() ([]byte, error) {
		// Inline to avoid a circular import when imports.go also uses os.
		return nil, nil
	}
	_ = import_os_ReadFile

	// Use os directly here.
	return osReadFile(path)
}