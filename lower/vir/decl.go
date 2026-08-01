package vir

import (
	"github.com/vertex-language/vertex/lower/hir"
	ir "github.com/vertex-language/vvm/ir/vir"
)

// decl.go emits everything above the function section: structs, consts,
// globals, links, extern groups, imports.
//
// The only non-mechanical thing here is struct ordering. hir appends a
// *Struct to its module *before* lowering the struct's fields, so that a
// self-referential field reaches its own enclosing type without looping —
// which means a struct-typed field's shape lands in the slice *after* the
// struct that names it. vir has no forward references (§2.2), so the
// declaration order has to be recomputed rather than preserved.

func (l *lowerer) structs(m *ir.Module, hm *hir.Module) {
	for _, s := range l.orderStructs(hm) {
		fields := make([]ir.Field, 0, len(s.Fields))
		for _, f := range s.Fields {
			fields = append(fields, ir.Field{Name: f.Name, Type: l.typ(hm, f.Type)})
		}
		m.DeclareStruct(s.Name, fields...)
	}
}

// orderStructs returns hm's structs in dependency order, each after every
// struct reachable from its fields by value. A pointer field erases the
// dependency (vir's ptr is untyped), so a cycle here would mean a struct
// of infinite size — which the analyzer already rejected. Reaching one
// anyway is a compiler bug, reported as such.
func (l *lowerer) orderStructs(hm *hir.Module) []*hir.Struct {
	own := make(map[*hir.Struct]bool, len(hm.Structs))
	for _, s := range hm.Structs {
		own[s] = true
	}
	const (
		open = 1
		done = 2
	)
	state := make(map[*hir.Struct]int, len(hm.Structs))
	out := make([]*hir.Struct, 0, len(hm.Structs))

	var visit func(s *hir.Struct)
	visit = func(s *hir.Struct) {
		switch state[s] {
		case done:
			return
		case open:
			l.errorf(0, "internal: struct %s participates in a by-value cycle, which has no layout", s.Name)
			return
		}
		state[s] = open
		for _, f := range s.Fields {
			for _, dep := range byValueStructs(f.Type) {
				if own[dep] {
					visit(dep)
				}
			}
		}
		state[s] = done
		out = append(out, s)
	}
	for _, s := range hm.Structs {
		visit(s)
	}
	return out
}

// byValueStructs lists the structs a type embeds without an intervening
// pointer. An array of structs embeds its element; ptr embeds nothing.
func byValueStructs(t hir.Type) []*hir.Struct {
	switch x := t.(type) {
	case hir.StructType:
		return []*hir.Struct{x.Def}
	case hir.ArrayType:
		return byValueStructs(x.Elem)
	}
	return nil
}

func (l *lowerer) consts(m *ir.Module, hm *hir.Module) {
	for _, c := range hm.Consts {
		m.DeclareConstant(c.Name, l.typ(hm, c.Type), l.value(hm, c.Value))
	}
}

func (l *lowerer) globals(m *ir.Module, hm *hir.Module) {
	for _, g := range hm.Globals {
		vg := m.DeclareGlobal(g.Name, l.typ(hm, g.Type), l.init(hm, g.Init))
		if g.Export {
			vg.Exported()
		}
		if g.TLS {
			vg.ThreadLocal()
		}
		if g.Align > 0 {
			vg.Aligned(g.Align)
		}
	}
}

// init maps hir's folded initializer forms onto vir's init grammar (§6.2).
// Both are already narrow — hir did the folding precisely so this is a
// rename.
func (l *lowerer) init(hm *hir.Module, in hir.Init) ir.ConstInit {
	switch x := in.(type) {
	case nil:
		return ir.InitZero{}
	case hir.InitZero:
		return ir.InitZero{}
	case hir.InitConst:
		return ir.InitLiteral{Value: l.value(hm, x.Value)}
	case hir.InitBytes:
		return ir.InitByteString{Data: x.Data}
	case hir.InitAddr:
		return ir.InitAddressOf{Name: x.Name}
	case hir.InitAggregate:
		elems := make([]ir.ConstInit, 0, len(x.Elems))
		for _, e := range x.Elems {
			elems = append(elems, l.init(hm, e))
		}
		return ir.InitAggregate{Elems: elems}
	}
	l.errorf(0, "internal: no vir init form for %T", in)
	return ir.InitZero{}
}

// links emits the link declarations and their extern groups together,
// since §7.2 requires every extern group's Dependency to match a Link.Name
// byte-for-byte and hir already built them as a pair.
func (l *lowerer) links(m *ir.Module, hm *hir.Module) {
	for _, lk := range hm.Links {
		m.DeclareLink(ir.LinkKind(lk.Kind), lk.Name)
	}
	for _, eg := range hm.Externs {
		g := m.DeclareExternGroup(eg.Dependency)
		for _, f := range eg.Funcs {
			ef := g.DeclareFunction(f.Name, l.params(hm, f.Params), l.typ(hm, f.Result))
			if f.Variadic {
				ef.SetVariadic()
			}
		}
	}
}

// imports emits one line per cross-module dependency. hir stores module
// *names*, which is the bare form §7.3 accepts, and every qualified call
// operand downstream is spelled against the same name.
func (l *lowerer) imports(m *ir.Module, hm *hir.Module) {
	for _, name := range hm.Imports {
		m.DeclareImport(name)
	}
}

func (l *lowerer) params(hm *hir.Module, ps []*hir.Param) []ir.Param {
	if len(ps) == 0 {
		return nil
	}
	out := make([]ir.Param, 0, len(ps))
	for _, p := range ps {
		out = append(out, l.param(hm, p))
	}
	return out
}

// param applies vir's aggregate conventions. Both byval and sret mean the
// parameter's vir type is ptr — hir states this once in its package doc
// and this is the one place that reads it off.
func (l *lowerer) param(hm *hir.Module, p *hir.Param) ir.Param {
	out := ir.Param{Name: p.Name, Type: l.typ(hm, p.Type)}
	switch {
	case p.SRet != nil:
		out.SRet, out.Type = p.SRet.Name, ir.Ptr
	case p.ByVal != nil:
		out.ByVal, out.Type = p.ByVal.Name, ir.Ptr
	}
	return out
}