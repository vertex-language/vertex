// decl.go
package vir

import (
	"strings"

	"github.com/vertex-language/vertex/lower/hir"
	vir "github.com/vertex-language/vvm/ir/vir"
)

// ------------------------------------------------------------------ structs

// structs emits the struct section in dependency order.
//
// hir appends a *Struct to its module *before* lowering the struct's
// fields, so a self-referential field can reach its enclosing type without
// looping. That means a struct-typed field's shape lands in the slice after
// the struct naming it — exactly the order vir forbids (§2.2, no forward
// references). orderStructs does the post-order walk that fixes it.
func (ml *moduleLowerer) structs() {
	for _, s := range ml.orderStructs() {
		fields := make([]vir.Field, 0, len(s.Fields))
		for _, f := range s.Fields {
			fields = append(fields, vir.Field{Name: f.Name, Type: ml.typ(f.Type)})
		}
		// Deliberately not Exported(): hir declares synthesized headers
		// (_Vstr, _Vvec, _Vfn, tuples) once per module, a struct produces no
		// symbol, and duplicating a declaration is safe for layout. It is
		// not obviously safe for byval/sret, which compare nominally per
		// origin — the first cross-package []T passed by value wants a
		// decision, and this is where it lands.
		ml.out.DeclareStruct(s.Name, fields...)
	}
}

const (
	unvisited = iota
	visiting
	visited
)

func (ml *moduleLowerer) orderStructs() []*hir.Struct {
	state := make(map[*hir.Struct]int, len(ml.src.Structs))
	out := make([]*hir.Struct, 0, len(ml.src.Structs))
	var stack []*hir.Struct

	var visit func(s *hir.Struct)
	visit = func(s *hir.Struct) {
		switch state[s] {
		case visited:
			return
		case visiting:
			// A by-value cycle is a struct of infinite size, which
			// semantics.md §3.4 already forbids — so this is a bug in the
			// checker or in hir, never in the user's source.
			ml.bug("struct cycle through by-value fields: " + cycleOf(stack, s))
		}
		state[s] = visiting
		stack = append(stack, s)
		for _, dep := range ml.structDeps(s) {
			visit(dep)
		}
		stack = stack[:len(stack)-1]
		state[s] = visited
		out = append(out, s)
	}

	for _, s := range ml.src.Structs {
		visit(s)
	}
	return out
}

// structDeps collects the by-value struct types a struct's fields name. A
// pointer field erases the edge — that is what makes a linked list finite —
// and a cross-module field is an import reference rather than a local
// ordering constraint.
func (ml *moduleLowerer) structDeps(s *hir.Struct) []*hir.Struct {
	var deps []*hir.Struct
	var walk func(t *hir.Type)
	walk = func(t *hir.Type) {
		if t == nil {
			return
		}
		switch t.Kind {
		case hir.KStruct:
			if t.Struct != nil && t.Struct.Module == ml.src && t.Struct != s {
				deps = append(deps, t.Struct)
			}
		case hir.KArray:
			walk(t.Elem)
		}
	}
	for _, f := range s.Fields {
		walk(f.Type)
	}
	return deps
}

func cycleOf(stack []*hir.Struct, s *hir.Struct) string {
	names := make([]string, 0, len(stack)+1)
	start := 0
	for i, x := range stack {
		if x == s {
			start = i
			break
		}
	}
	for _, x := range stack[start:] {
		names = append(names, x.Name)
	}
	return strings.Join(append(names, s.Name), " -> ")
}

// ------------------------------------------------------- constants, globals

func (ml *moduleLowerer) constants() {
	for _, c := range ml.src.Consts {
		d := ml.out.DeclareConstant(c.Name, ml.typ(c.Type), ml.operand(c.Value))
		if c.Export {
			d.Exported()
		}
	}
}

func (ml *moduleLowerer) globals() {
	for _, g := range ml.src.Globals {
		d := ml.out.DeclareGlobal(g.Name, ml.typ(g.Type), ml.constInit(g.Init))
		if g.Export {
			d.Exported()
		}
		if g.TLS {
			d.ThreadLocal()
		}
		if g.Align > 0 {
			d.Aligned(g.Align)
		}
	}
}

// --------------------------------------------------------- links and externs

// links emits the runtime dependency plus whatever a declare block asked
// for. A module lowered from ordinary Vertex source carries no hir Links at
// all, so everything past the first line came from user source — the
// invariant is checkable by reading the emitted text.
func (ml *moduleLowerer) links() {
	ml.out.DeclareLink(vir.LinkStatic, runtimeLink)
	for _, l := range ml.src.Links {
		ml.out.DeclareLink(linkKind(ml, l.Kind), l.Name)
	}
}

func linkKind(ml *moduleLowerer, k string) vir.LinkKind {
	switch k {
	case "static":
		return vir.LinkStatic
	case "shared":
		return vir.LinkShared
	case "framework":
		return vir.LinkFramework
	}
	ml.bug("link kind with no vir spelling: " + k)
	return vir.LinkStatic
}

// externs emits one group per declare block.
//
// No group is emitted for vertexrt. hir reaches every runtime symbol as a
// *qualified* call into a builtins module and records the corresponding
// import (see funcBuilder.callBuiltin), so the runtime is an import
// dependency at the IR level and a link dependency at the object level, and
// the two are declared in the two different places that say so. See the
// note at the bottom of func.go.
func (ml *moduleLowerer) externs() {
	for _, g := range ml.src.Externs {
		og := ml.out.DeclareExternGroup(g.Dependency)
		for _, f := range g.Functions {
			var attrs []vir.FunctionAttribute
			if f.NoReturn {
				attrs = append(attrs, vir.AttributeNoReturn)
			}
			ef := og.DeclareFunction(f.Name, ml.params(f.Params), ml.typ(f.Result), attrs...)
			if f.Variadic {
				ef.SetVariadic()
			}
		}
	}
}

func (ml *moduleLowerer) declareImports() {
	// hir's list first, in its own first-reference order, so a build stays
	// byte-reproducible; then anything translation discovered.
	for _, p := range ml.src.Imports {
		ml.needImport(p)
	}
	for _, p := range ml.imports {
		ml.out.DeclareImport(p)
	}
}