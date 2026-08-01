package importer

import "sort"

// This file is the query surface over a completed load. Nothing here parses or
// checks; it answers questions a driver and a lowering pass ask about the
// graph they were handed.

// Sorted returns every loaded package by path, in lexical order. Load's Order
// is topological and therefore not stable against an unrelated edit; this is
// what a caller that wants reproducible output over the whole set uses.
func (r *Result) Sorted() []*Package {
	out := make([]*Package, 0, len(r.Packages))
	for _, p := range r.Packages {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// Lookup returns a loaded package by import path.
func (r *Result) Lookup(path string) *Package {
	return r.Packages[path]
}

// Deps returns every package reachable from p, excluding p, in topological
// order. This is the set a link step needs and the set lower walks.
func (r *Result) Deps(p *Package) []*Package {
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

// Importers returns every loaded package that directly imports path. It is the
// reverse edge, which nothing in a build needs and everything in tooling does.
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
// A.6.1 ⊢ "a function named main must take no parameters, return nothing, and
// acts as the program entry point." Which package that lives in is not
// something the annex fixes, so this searches the roots rather than assuming a
// package named main.
func (r *Result) Entry() *Package {
	for _, p := range r.Roots {
		if p.Types == nil {
			continue
		}
		if obj := p.Types.Scope().Lookup("main"); obj != nil {
			if fn, ok := obj.(interface{ IsEntry() bool }); ok && fn.IsEntry() {
				return p
			}
		}
	}
	return nil
}

func sortedKeys(m map[string]*Package) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}