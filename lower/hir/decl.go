package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// collect is pass 1: everything whose shape is known without walking a body.
//
// Function bodies are deliberately absent. A body is lowered only when
// monomorphization reaches it, so a generic never instantiated emits
// nothing and a package whose functions are all unreachable emits only its
// types and globals.
func (l *lowerer) collect(u *Unit) {
	for _, f := range u.Files {
		for _, d := range f.Decls {
			switch x := d.(type) {
			case *ast.RecordDecl:
				l.collectRecord(x)
			case *ast.EnumDecl:
				l.collectEnum(x)
			case *ast.VarDecl:
				l.collectTopLevelVar(x)
			case *ast.DeclareDecl:
				l.collectDeclare(x)
			case *ast.TypeAliasDecl, *ast.ConstraintDecl:
				// Neither has a lowering. A transparent alias mints no type
				// at all, an abstract one is a ptr, and a constraint is not a
				// value type — all already settled in types.
			case *ast.FuncDecl:
				// Deferred to the worklist.
			}
		}
	}
}

// collectRecord forces a generic-free record's shape into existence eagerly,
// so a struct that is only ever named in a signature still gets declared.
// A generic one has no shape until it is instantiated.
func (l *lowerer) collectRecord(d *ast.RecordDecl) {
	if d.TypeParams != nil {
		return
	}
	obj := l.defOf(d.Name)
	tn, ok := obj.(*types.TypeName)
	if !ok || tn.Type() == nil {
		return
	}
	l.typ(tn.Type(), l.mod)
}

func (l *lowerer) collectEnum(d *ast.EnumDecl) {
	if d.TypeParams != nil {
		return
	}
	tn, ok := l.defOf(d.Name).(*types.TypeName)
	if !ok || tn.Type() == nil {
		return
	}
	l.typ(tn.Type(), l.mod)
}

// collectTopLevelVar lowers a top-level binding.
//
// A `let` whose type is scalar becomes a vir const — no runtime storage. A
// `var`, or a `let` of aggregate type, becomes a global. vir's global init
// grammar is narrower than semantics.md §5.3's constant expressions, so the
// folded types.Value is projected onto literal/zero/bytes/aggregate here and
// anything that will not fit is a bug rather than a diagnostic: the analyzer
// already required the initializer to be compile-time-evaluable.
func (l *lowerer) collectTopLevelVar(d *ast.VarDecl) {
	for i, b := range d.Bindings {
		obj := l.defOf(b.Name)
		if obj == nil || b.Name.IsBlank() {
			continue
		}
		t := l.typ(obj.Type(), l.mod)
		name := l.qualify(obj.Name())

		if c, ok := obj.(*types.Const); ok && !IsAggregate(t) {
			l.mod.Consts = append(l.mod.Consts, &Const{
				Name:   name,
				Type:   t,
				Value:  l.constValue(c.Val(), t),
				Export: true,
			})
			continue
		}

		init := ConstInit(InitZero{})
		if i < len(d.Values) {
			init = l.constInit(obj, t)
		}
		l.mod.Globals = append(l.mod.Globals, &Global{
			Name:   name,
			Type:   t,
			Init:   init,
			Export: true,
			Align:  int(l.alignOf(t)),
		})
	}
}

// constInit projects a folded constant onto vir's global-init grammar.
func (l *lowerer) constInit(obj types.Object, t *Type) ConstInit {
	c, ok := obj.(*types.Const)
	if !ok {
		// A top-level `var` still has a constant initializer, but types.Var
		// does not carry the folded value — the checker discarded it after
		// the representability check.
		//
		// todo: analyzer should record a top-level var's folded value in
		// Info (a new table, or a *types.Const-shaped side entry), so this
		// stops being a zero.
		return InitZero{}
	}
	v := c.Val()
	if s, ok := types.StringVal(v); ok {
		// A string global is its bytes; the {ptr,len} header is built at each
		// use, since the header is a value and the bytes are storage.
		return InitBytes{Data: []byte(s)}
	}
	return InitScalar{Value: l.constValue(v, t)}
}

func (l *lowerer) constValue(v types.Value, t *Type) Value {
	switch {
	case v == nil:
		return Int(0, t)
	case IsFloat(t):
		f, _ := types.FloatVal(v)
		return Float(f, t)
	default:
		n, _ := types.Int64Val(v)
		return Int(n, t)
	}
}

// ------------------------------------------------------------ declare blocks

// collectDeclare lowers a declare block into a link line and an extern
// group. This is the one path by which user source touches the linker, and
// grammar.md already fenced it: a file carrying a declare block must carry a
// BuildClause, which is exactly the information needed to emit a triple.
func (l *lowerer) collectDeclare(d *ast.DeclareDecl) {
	kind := "shared"
	if d.Kind == token.CtxFramework {
		kind = "framework"
	}
	lib, ok := stringLit(d.Path)
	if !ok {
		l.bugAt(d.Pos(), "declare block with an undecodable path")
	}
	l.mod.Links = append(l.mod.Links, &Link{Kind: kind, Name: lib})
	g := &ExternGroup{Dependency: lib}
	l.mod.Externs = append(l.mod.Externs, g)

	for _, m := range d.Members {
		switch x := m.(type) {
		case *ast.ForeignFunc:
			if x.Name == nil {
				// The unnamed initializer form. It has no package-scope name
				// and is reached only through its enclosing class.
				continue
			}
			obj, _ := l.useOrDefOf(x.Name).(*types.Func)
			if obj == nil {
				continue
			}
			g.Functions = append(g.Functions, l.externFunc(obj.Name(), obj.Signature()))
		case *ast.ForeignClass:
			// todo: foreign class members are methods on an opaque handle,
			// and the flat-C case is a direct call with `this` first.
			// Objective-C (declare framework, darwin) needs the selector
			// cache of lowering.md §20.4, and C++/COM need fnsig-typed
			// indirect calls, which nothing upstream produces yet.
			l.todoAt(x.Pos(), "foreign class members")
		}
	}
}

func (l *lowerer) externFunc(name string, sig *types.Signature) *ExternFunc {
	if sig == nil {
		return &ExternFunc{Name: name, Result: Void}
	}
	f := &ExternFunc{Name: name, Variadic: sig.Variadic()}
	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		t := l.typ(p.Type(), l.mod)
		if IsAggregate(t) {
			// A foreign boundary takes a pointer, never a Vertex aggregate
			// by value: the foreign side has its own layout expectations and
			// a declare block describes call shape only.
			t = Ptr
		}
		f.Params = append(f.Params, &Param{Name: p.Name(), Type: t})
	}
	if sig.Results().Len() == 1 {
		f.Result = l.typ(sig.Results().At(0).Type(), l.mod)
	} else {
		f.Result = Void
	}
	return f
}

// ------------------------------------------------------------------ helpers

func (l *lowerer) defOf(id *ast.Ident) types.Object {
	if l.info == nil {
		return nil
	}
	return l.info.Defs[id]
}

func (l *lowerer) useOrDefOf(id *ast.Ident) types.Object {
	return l.info.ObjectOf(id)
}

func stringLit(b *ast.BasicLit) (string, bool) {
	if b == nil || len(b.Value) < 2 {
		return "", false
	}
	// The scanner kept the raw spelling; a declare path admits no escape
	// worth decoding, so the quotes come off and nothing else happens.
	return b.Value[1 : len(b.Value)-1], true
}