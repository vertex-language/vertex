// Package importer resolves a Vertex import graph into checked packages.
//
// It is the seam ast.NewPackage names: that constructor "is a validated
// container and nothing more: no I/O, no import resolution, no scopes.
// Resolution belongs to the loader, which is why this takes no importer."
// parser.ParseDir handles one directory; this handles the graph.
//
// The ordering constraint is §1.3: the qualifier under which an imported
// package's symbols are reached comes from that package's own package clause,
// and the path is a locator rather than a name. A file's names therefore cannot
// resolve until every directory it imports has had its package line read —
// which is what parser.PackageClauseOnly exists for, and why loading is a
// post-order walk rather than a single pass.
//
// Nothing here type-checks. It parses, orders, and calls the analyzer; the
// analyzer decides what the result means.
//
// Citation convention, matching types and analyzer: a bare § is semantics.md,
// CamelCase names are grammar.md productions. Where neither document fixes
// something — how a path locates a directory, whether a relative path is even
// spelled — the comment says so rather than presenting a choice as a rule.
package importer

import (
	"fmt"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// Config controls a load.
type Config struct {
	// Fset is the position space every loaded file shares. One per load, so a
	// diagnostic can span two files in two different packages without carrying
	// a file reference alongside each position.
	Fset *token.FileSet

	// Tag is the target. It decides which files are in the build at all, and —
	// because a `build test` file is the one whose tag changes what is
	// grammatical rather than only what is linkable — what an ExpectedType
	// result is licensed by.
	Tag token.BuildTag

	// Resolver maps an import path to a directory. Nil means a DirResolver
	// rooted at the current directory.
	Resolver Resolver

	// Reporter receives every diagnostic from every package as it is produced.
	// May be nil, in which case diagnostics are collected onto the returned
	// *Result instead.
	Reporter diag.Reporter

	// Mode is passed through to the parser. ParseComments is the usual choice
	// for tooling and unnecessary for a build.
	Mode parser.Mode
}

// Package is one loaded package: its syntax, its checked types, and its
// resolved dependencies.
type Package struct {
	Path string
	Dir  string

	// Name is the package clause name, which every file in the directory
	// agrees on — ast.NewPackage already checked that. It is not derived from
	// Path: §1.3 makes the two independent, and this is the qualifier an
	// importer reaches this package's symbols under.
	Name  string
	Files []*ast.File

	Types *types.Package
	Info  *types.Info

	// Imports maps each import path this package declares to the loaded
	// package it resolved to. There is no aliasing form, no dot-import, and no
	// blank import, so the mapping is the whole story.
	Imports map[string]*Package

	// Errors is this package's own diagnostics, in the order produced. A
	// package with errors is still returned: partial results are what editor
	// tooling needs, and most of what this toolchain rejects is a form that
	// parses and is diagnosed rather than one that aborts a load.
	Errors []*diag.Diagnostic
}

func (p *Package) String() string { return p.Path }

// HasErrors reports whether this package produced an error-severity
// diagnostic. A package can hold warnings and still be clean.
func (p *Package) HasErrors() bool {
	for _, d := range p.Errors {
		if d.Sev == diag.Error {
			return true
		}
	}
	return false
}

// Result is one load's output.
type Result struct {
	// Roots are the packages named in the Load call, in the order given. A
	// path named twice appears twice and is the same *Package both times.
	Roots []*Package

	// Packages holds every loaded package by import path, roots included.
	Packages map[string]*Package

	// Order is topological: a package appears after everything it imports.
	// This is the order the analyzer ran in and the order a lowering pass
	// should follow, since §1.3 makes a package's qualifiers depend on its
	// imports having been read.
	//
	// It is a partial order in one direction only — a package that failed to
	// parse or resolve is absent, because it was never checked.
	Order []*Package

	// Diagnostics is every diagnostic from every package, sorted and deduped.
	// Nil when Config.Reporter was non-nil, since they streamed there instead;
	// per-package Errors are populated either way.
	Diagnostics *diag.List
}

// HasErrors reports whether any loaded package produced an error-severity
// diagnostic. It reads the per-package lists, which are populated whether or
// not a Reporter was supplied — the aggregated list is not.
func (r *Result) HasErrors() bool {
	for _, p := range r.Packages {
		if p.HasErrors() {
			return true
		}
	}
	return false
}

// Load resolves and checks the packages named by paths, plus everything they
// import transitively.
//
// A returned error is fatal to the load — an unresolvable path or an import
// cycle, neither of which any single file could be blamed for. Ordinary
// rejections are diagnostics and do not produce one. The *Result is non-nil
// either way and holds whatever finished before the failure, which is what a
// caller wanting partial results reads.
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

	out := &Result{Packages: l.packages}
	defer func() {
		// Assigned on every exit, so a failed load still hands back the
		// packages that did finish.
		out.Order = l.order
		if l.list != nil {
			l.list.Sort()
			l.list.Dedup()
			out.Diagnostics = l.list
		}
	}()

	for _, path := range paths {
		pkg, err := l.load(path, nil)
		if err != nil {
			return out, err
		}
		out.Roots = append(out.Roots, pkg)
	}
	return out, nil
}