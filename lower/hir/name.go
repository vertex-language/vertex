package hir

import (
	"strings"

	"github.com/vertex-language/vertex/types"
)

// Naming. The `_V` prefix is reserved for compiler-synthesized names;
// semantics.md's list of what may not be declared needs it added, which is
// one line and is not yet there.
//
// vir applies its own Itanium-style mangling on top of everything here, so
// nothing pre-mangles: these are idents, not symbols.

// funcName is the ident a function instance gets.
//
//	func add                  -> add
//	func (w: Widget) rename   -> Widget_rename
//	func (w: Widget) init     -> Widget_init  (named: Widget_init_withRect)
//	smaller[float64]          -> smaller__f64
func (l *lowerer) funcName(obj *types.Func, args []types.Type) string {
	var b strings.Builder
	if sig := obj.Signature(); sig != nil && sig.Recv() != nil {
		if n := types.AsNamed(sig.Recv().Type()); n != nil && n.Obj() != nil {
			b.WriteString(n.Obj().Name())
			b.WriteByte('_')
		}
	}
	b.WriteString(obj.Name())
	for _, a := range args {
		b.WriteString("__")
		b.WriteString(mangleType(l.typ(a, l.mod)))
	}
	return b.String()
}

// qualify is the ident a package-scope global or const gets. Method names
// live in no scope, so a method `read` and a function `read` in one package
// cannot collide before mangling and do not collide after it — which means
// nothing needs qualifying beyond the receiver prefix above.
func (l *lowerer) qualify(name string) string { return name }

// mangleTypeName is the ident a declared type's struct gets.
func (l *lowerer) mangleTypeName(t types.Type) string {
	if n := types.AsNamed(t); n != nil && n.Obj() != nil {
		name := n.Obj().Name()
		for _, a := range n.TypeArgs() {
			name += "__" + mangleType(l.typ(a, l.mod))
		}
		return name
	}
	return "_Vanon_" + mangleType(l.typ(t, l.mod))
}

// mangleType is the short spelling a type contributes to an instantiated
// name. Nested instantiations recurse.
func mangleType(t *Type) string {
	switch t.Kind {
	case KVoid:
		return "void"
	case KInt:
		if t.Signed {
			return "i" + itoa(t.Bits)
		}
		return "u" + itoa(t.Bits)
	case KFloat:
		return "f" + itoa(t.Bits)
	case KPtr:
		return "ptr"
	case KStruct:
		return t.Struct.Name
	case KArray:
		return "a" + itoa(int(t.Len)) + "_" + mangleType(t.Elem)
	case KVector:
		return "v" + itoa(int(t.Len)) + "_" + mangleType(t.Elem)
	case KPredicate:
		return "mask" + itoa(int(t.Len))
	}
	return "x"
}

// internString gives a literal's bytes one global per module per distinct
// value. The name is content-derived so two occurrences of the same literal
// share storage and the emitted text stays stable.
func (l *lowerer) internString(m *Module, s string) string {
	name := "_Vstr_" + itoa(len(s)) + "_" + itoa(int(hashString(s)))
	for _, g := range m.Globals {
		if g.Name == name {
			return name
		}
	}
	m.Globals = append(m.Globals, &Global{
		Name:  name,
		Type:  ArrayOf(U8, int64(len(s))),
		Init:  InitBytes{Data: []byte(s)},
		Align: 1,
	})
	return name
}

func hashString(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h & 0x7fffffff
}