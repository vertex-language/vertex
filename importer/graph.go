package importer

import (
	"sort"

	"github.com/vertex-language/vertex/types"
)

// This file is the query surface over a completed load. Nothing here parses or
// checks; it answers the questions a driver and a lowering pass ask about the
// graph they were handed.

// Sorted returns every loaded package by path, in lexical order. Order is
// topological and therefore not stable against an unrelated edit; this is what
// a caller wanting reproducible output over the whole set uses.
func (r *Result) Sorted() []*Package {
	out := make([]*Package, 0, len(r.Packages))
	for _, p := range r.Packages {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// Lookup returns a loaded package by import path, or nil.
func (r *Result) Lookup(path string) *Package { return r.Packages[path] }

// Deps returns every package reachable from p, excluding p, in topological
// order: the set a link step needs and the set a lowering pass walks.
//
// Each level is walked in path order, so the result is deterministic even
// though the graph is a map. A cycle cannot reach here — one is fatal to the
// load — but the seen set makes the walk terminate regardless.
func (r *Result) Deps(p *Package) []*Package {
	if p == nil {
		return nil
	}
	seen := make(map[string]bool)
	var out []*Package

	var walk func(*Package)
	walk = func(cur *Package) {
		for _, path := range sortedKeys(cur.Imports) {
			dep := cur.Imports[path]
			if seen[dep.Path] {
				continue
			}
			seen[dep.Path] = true
			walk(dep)
			out = append(out, dep)
		}
	}
	walk(p)
	return out
}

// Importers returns every loaded package that directly imports path, in path
// order. It is the reverse edge, which nothing in a build needs and everything
// in tooling does.
func (r *Result) Importers(path string) []*Package {
	var out []*Package
	for _, p := range r.Sorted() {
		if _, ok := p.Imports[path]; ok {
			out = append(out, p)
		}
	}
	return out
}

// Entry returns the package holding the program's entry point, or nil.
//
// §1.4 ⊢ a program has exactly one package named `main` declaring exactly one
// `func main()` — no parameters, no result, no marker. Func.IsEntry answers all
// four; that there is exactly one such package across the load is this
// function's half, and it answers nil rather than picking when the load holds
// more than one.
//
// The search covers the whole graph rather than only the roots: a driver may
// name a library path and still want to know whether an entry point is reachable
// from it.
func (r *Result) Entry() *Package {
	var found *Package
	for _, p := range r.Sorted() {
		if p.Types == nil || p.Types.Name() != "main" {
			continue
		}
		fn, ok := p.Types.Scope().Lookup("main").(*types.Func)
		if !ok || !fn.IsEntry() {
			continue
		}
		if found != nil {
			return nil // more than one; the caller must say which
		}
		found = p
	}
	return found
}

func sortedKeys(m map[string]*Package) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}