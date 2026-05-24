// Package compiler is the Vertex language front-end compiler.
//
// It parses .vs source files using the ANTLR-based parser at
// github.com/vertex-language/vertex/parser, performs a two-pass
// compilation (declaration collection + wasm code generation), and
// produces either a native ELF binary or a raw .wasm file.
package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	backendCompiler "github.com/vertex-language/compiler"
	"github.com/vertex-language/compiler/encoder"
	"github.com/vertex-language/compiler/linker"
	"github.com/vertex-language/compiler/object"
)

// Options controls the compilation.
type Options struct {
	// Target platform: "linux", "darwin", "windows". Default: "linux".
	Platform string

	// Additional build tags (e.g. "syscalls").
	ExtraTags []string

	// OutputWasm: if true, produce a .wasm binary instead of a native ELF.
	OutputWasm bool

	// Entry is the export name of the entry-point function. Default: "main".
	Entry string

	// ModuleRoot is the directory that contains go.mod / vertex.mod.
	// Defaults to the directory of the first source file.
	ModuleRoot string

	// ModulePath is the Go module path prefix (e.g. "github.com/acme/myapp").
	ModulePath string

	// DCE enables dead-code elimination. When true, only functions reachable
	// from the entry point are emitted into the wasm module. This reduces
	// binary size but adds a full call-graph analysis pass.
	//
	// Extern (C) functions are always kept when they are actually called from
	// reachable Vertex code; unreferenced externs are dropped.
	DCE bool
}

func (o *Options) platform() string {
	if o.Platform != "" {
		return o.Platform
	}
	return "linux"
}

func (o *Options) entry() string {
	if o.Entry != "" {
		return o.Entry
	}
	return "main"
}

func (o *Options) buildTags() *BuildTags {
	tags := []string{o.platform()}
	tags = append(tags, o.ExtraTags...)
	return &BuildTags{Tags: tags}
}

// CompileFiles compiles one or more .vs source files and returns the binary.
// If opts.OutputWasm is true the output is a .wasm module; otherwise it is a
// native ELF binary (only linux/ELF is supported by the linker at this time).
func CompileFiles(paths []string, opts Options) ([]byte, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no source files")
	}

	// Determine module root.
	root := opts.ModuleRoot
	if root == "" {
		root = filepath.Dir(paths[0])
	}

	// Parse all source files.
	pkg := &Package{
		Dir:        root,
		ImportPath: opts.ModulePath,
		BuildTags:  opts.buildTags(),
	}

	for _, p := range paths {
		sf, err := ParseFile(p)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", p, err)
		}
		for _, pe := range sf.ParseErrs {
			_, _ = fmt.Fprintln(os.Stderr, pe)
		}
		if len(sf.ParseErrs) > 0 {
			return nil, fmt.Errorf("%d parse error(s) in %s", len(sf.ParseErrs), p)
		}
		if pkg.Name == "" && sf.PackName != "" {
			pkg.Name = sf.PackName
		}
		pkg.Files = append(pkg.Files, sf)
	}

	// Resolve and load imports that are Vertex packages (not lib/).
	resolver := NewResolver(root, opts.ModulePath)
	if err := loadImports(pkg, resolver, opts.buildTags()); err != nil {
		return nil, err
	}

	// Generate wasm module.
	genOpts := GenerateOptions{
		DCE:       opts.DCE,
		EntryName: opts.entry(),
	}
	mod, err := Generate(pkg, opts.buildTags(), genOpts)
	if err != nil {
		return nil, fmt.Errorf("codegen: %w", err)
	}

	if opts.OutputWasm {
		// Encode to .wasm binary.
		data, err := encoder.Encode(mod)
		if err != nil {
			return nil, fmt.Errorf("wasm encode: %w", err)
		}
		return data, nil
	}

	// Native binary output: only ELF is supported by the linker right now.
	if plat := opts.platform(); plat != "linux" {
		return nil, fmt.Errorf(
			"native output is not yet supported for platform %q (use -wasm, or target linux)",
			plat,
		)
	}

	// Compile wasm → native object.
	obj, err := backendCompiler.CompileWith(mod, backendCompiler.Options{})
	if err != nil {
		return nil, fmt.Errorf("backend compile: %w", err)
	}

	// Link → ELF binary.
	bin, err := linker.Link([]*object.WasmObj{obj}, linker.Options{
		Output: linker.ELF,
		Entry:  opts.entry(),
	})
	if err != nil {
		return nil, fmt.Errorf("link: %w", err)
	}
	return bin, nil
}

// loadImports recursively loads Vertex-package imports.
// lib/ imports are skipped (handled by extern declarations in the source).
func loadImports(pkg *Package, r *Resolver, tags *BuildTags) error {
	seen := map[string]bool{pkg.ImportPath: true}
	return loadImportsRec(pkg, r, tags, seen)
}

func loadImportsRec(pkg *Package, r *Resolver, tags *BuildTags, seen map[string]bool) error {
	for _, sf := range pkg.Files {
		for _, imp := range sf.Imports {
			if imp.IsNative() {
				continue // native interface binding; no .vs files to load
			}
			if seen[imp.Raw] {
				continue
			}
			seen[imp.Raw] = true

			fromDir := filepath.Dir(sf.Path)
			files, err := r.ResolveFiles(imp.Raw, fromDir)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "warn: cannot resolve import %q: %v\n", imp.Raw, err)
				continue
			}

			subPkg := &Package{
				ImportPath: imp.Raw,
				BuildTags:  tags,
				Dir:        filepath.Dir(files[0]),
			}
			for _, fp := range files {
				if !PlatformMatch(fp, tags) {
					continue
				}
				parsedSf, err := ParseFile(fp)
				if err != nil {
					return err
				}
				if subPkg.Name == "" && parsedSf.PackName != "" {
					subPkg.Name = parsedSf.PackName
				}
				subPkg.Files = append(subPkg.Files, parsedSf)
			}

			pkg.Files = append(pkg.Files, subPkg.Files...)

			if err := loadImportsRec(subPkg, r, tags, seen); err != nil {
				return err
			}
		}
	}
	return nil
}

// Version is the compiler version string.
const Version = "0.1.0"

// Suppress unused import warning.
var _ = strings.TrimSpace