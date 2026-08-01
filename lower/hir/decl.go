package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/types"
)

// declarations lowers everything in a unit that is not a function body:
// struct and enum shapes, compile-time constants, top-level variables, and
// declare blocks.
//
// Function *bodies* are deliberately not walked here. They are reached only
// through monomorphization from a root, so a package's unreachable
// functions never produce code and the dead-symbol-elimination question is
// removed rather than answered.
func (l *lowerer) declarations(u *Unit) {
	mod := l.byUnit[u]
	prev := l.cur
	l.cur = &instance{unit: u, mod: mod}
	defer func() { l.cur = prev }()

	for _, f := range u.Files {
		for _, d := range f.Decls {
			switch d := d.(type) {
			case *ast.RecordDecl, *ast.EnumDecl, *ast.TypeAliasDecl,
				*ast.ConstraintDecl, *ast.FuncDecl:
				// Shapes are minted on demand by typeLowerer as bodies
				// reference them; a declaration nobody reaches produces
				// nothing. Constraints and aliases produce nothing at all.
			case *ast.VarDecl:
				l.globalDecl(u, mod, d)
			case *ast.DeclareDecl:
				l.declareBlock(u, mod, d)
			}
		}
	}
}

// globalDecl lowers a top-level VariableDeclaration. A.2 already requires a
// compile-time-evaluable initializer, and vir's global init form is
// narrower still — literal, zero, addr ident, or an aggregate of those,
// with no arithmetic and no const references. There is no
// static-initialization-order problem to have because there is no
// initialization-time code.
func (l *lowerer) globalDecl(u *Unit, mod *Module, d *ast.VarDecl) {
	for i, b := range d.Bindings {
		if b.Name.IsBlank() {
			continue
		}
		obj := u.Info.Defs[b.Name]
		if obj == nil {
			continue
		}
		ht := l.hirType(obj.Type())

		var init Init = InitZero{}
		if i < len(d.Values) {
			if v, ok := l.constInit(u, d.Values[i], obj.Type()); ok {
				init = v
			} else {
				l.todo(d.Values[i].Pos(), "top-level initializer is not reducible to a vir global init form")
			}
		}

		// A.5.1: `let` is immutable and may be folded away entirely, so a
		// scalar `let` becomes a vir const (no runtime storage) rather
		// than a global.
		if d.Kw == 0 /* LET */ && !IsAggregate(ht) {
			if c, ok := init.(InitConst); ok {
				mod.Consts = append(mod.Consts, &Const{
					Name: mod.uniqueName(obj.Name()), Type: ht, Value: c.Value,
				})
				l.bindGlobal(obj, mod.Consts[len(mod.Consts)-1].Name, ht, true)
				continue
			}
		}

		g := &Global{Name: mod.uniqueName(obj.Name()), Type: ht, Init: init}
		mod.Globals = append(mod.Globals, g)
		l.bindGlobal(obj, g.Name, ht, false)
	}
}

// constInit folds an initializer expression down to vir's init grammar. It
// is deliberately narrow: anything it cannot fold is reported rather than
// smuggled into an initialization-time code path that does not exist.
func (l *lowerer) constInit(u *Unit, x ast.Expr, want types.Type) (Init, bool) {
	tv, ok := u.Info.Types[x]
	if ok && tv.Value != nil {
		if s, isStr := types.StringVal_(tv.Value); isStr {
			return InitBytes{Data: []byte(s)}, true
		}
		if v, isInt := types.Int64Val(tv.Value); isInt {
			return InitConst{Value: IntVal(l.hirType(want), v)}, true
		}
		if b, isBool := types.BoolVal_(tv.Value); isBool {
			return InitConst{Value: BoolVal(b)}, true
		}
	}
	switch x := x.(type) {
	case *ast.ParenExpr:
		return l.constInit(u, x.X, want)
	case *ast.ArrayLit:
		var elems []Init
		et := l.elem(want)
		for _, e := range x.Elems {
			ei, ok := l.constInit(u, e, et)
			if !ok {
				return nil, false
			}
			elems = append(elems, ei)
		}
		return InitAggregate{Elems: elems}, true
	case *ast.CompositeLit:
		var elems []Init
		for _, e := range x.Elems {
			kv, ok := e.(*ast.KeyValueExpr)
			if !ok {
				return nil, false
			}
			ei, ok := l.constInit(u, kv.Value, u.Info.TypeOf(kv.Value))
			if !ok {
				return nil, false
			}
			elems = append(elems, ei)
		}
		return InitAggregate{Elems: elems}, true
	}
	return nil, false
}

// declareBlock lowers A.8's linkage boundary. This is the single exception
// to invariant 3 — and A.8.1 already requires such a file to carry a
// BuildClause, so the grammar fenced it before hir had to.
func (l *lowerer) declareBlock(u *Unit, mod *Module, d *ast.DeclareDecl) {
	kind := LinkShared
	switch {
	case d.Kind == "framework":
		kind = LinkFramework
	case hasTag(d.Variant, "static"):
		kind = LinkStatic
	}
	lib := unquote(d.Path.Value)
	mod.Links = append(mod.Links, Link{Kind: kind, Name: lib})

	g := &ExternGroup{Dependency: lib}
	mod.Externs = append(mod.Externs, g)
	l.foreignMembers(u, mod, g, d.Members)
}

func (l *lowerer) foreignMembers(u *Unit, mod *Module, g *ExternGroup, members []ast.ForeignMember) {
	for _, mem := range members {
		switch mem := mem.(type) {
		case *ast.ForeignFunc:
			if mem.Name == nil {
				// An unnamed initializer resolves through bare Type(...)
				// construction; it still needs a real entry point name,
				// which the analyzer recorded on the object.
				l.todo(mem.Pos(), "unnamed foreign initializer needs a symbol name from the checked object")
				continue
			}
			obj := u.Info.Defs[mem.Name]
			if obj == nil {
				continue
			}
			sig, _ := obj.Type().(*types.Signature)
			ef := &ExternFunc{Name: mem.Name.Name, Result: l.result(sig)}
			for _, p := range paramsOf(sig) {
				ef.Params = append(ef.Params, l.param(p.Name, p.Type))
			}
			ef.Variadic = sig != nil && sig.Variadic()
			g.Funcs = append(g.Funcs, ef)
			mod.names[ef.Name] = true
			// Record the checked object -> ExternFunc mapping so a call
			// naming this identifier (expr.go's callExpr) resolves here
			// directly instead of falling through to the monomorphization
			// worklist, which has no notion of a foreign body — see
			// lower.go's externs field doc for what that misrouting used
			// to produce.
			l.bindExtern(obj, ef)
		case *ast.ForeignClass:
			l.foreignMembers(u, mod, g, mem.Members)
		case *ast.ForeignField:
			// A.8.3 bans these and the analyzer already reported it.
		}
	}
}

// param applies vir's aggregate conventions: an aggregate parameter is
// passed byval, which means its vir type is ptr and the callee may not
// mutate the caller's object.
func (l *lowerer) param(name string, t types.Type) *Param {
	ht := l.hirType(t)
	p := &Param{Name: sanitize(name), Type: ht}
	if st, ok := ht.(StructType); ok {
		p.ByVal = st.Def
	}
	return p
}

func (l *lowerer) result(sig *types.Signature) Type {
	if sig == nil {
		return Void
	}
	res := sig.Results()
	if res == nil || res.Len() == 0 {
		return Void
	}
	if res.Len() == 1 {
		return l.hirType(res.At(0).Type())
	}
	// A multi-value return is a tuple; the tuple *is* the type, and an
	// aggregate result travels through sret.
	return l.hirType(res)
}

type paramInfo struct {
	Name string
	Type types.Type
	Mut  bool
}

func paramsOf(sig *types.Signature) []paramInfo {
	if sig == nil {
		return nil
	}
	ps := sig.Params()
	out := make([]paramInfo, 0, ps.Len())
	for i := 0; i < ps.Len(); i++ {
		v := ps.At(i)
		out = append(out, paramInfo{Name: v.Name(), Type: v.Type(), Mut: v.Mode() == types.ModeMut})
	}
	return out
}

func hasTag(v *ast.VariantTag, want string) bool {
	if v == nil {
		return false
	}
	for _, t := range v.Tags {
		if unquote(t.Value) == want {
			return true
		}
	}
	return false
}

func unquote(s string) string {
	if len(s) >= 2 && (s[0] == '"' || s[0] == '`') {
		return s[1 : len(s)-1]
	}
	return s
}