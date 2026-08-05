package hir

import (
	"strings"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/types"
)

// instance is one entry of the monomorphization worklist: a function plus
// the concrete type arguments it is being built for.
type instance struct {
	obj  *types.Func
	args []types.Type
	fn   *Func // the shell, allocated at enqueue time
	pkg  *types.Package
}

// key identifies an instance for memoization. Two call sites reaching
// smaller[int32] must produce one Func, or a diamond of instantiations
// builds the same body twice.
func instanceKey(obj *types.Func, args []types.Type) string {
	var b strings.Builder
	if obj.Pkg() != nil {
		b.WriteString(obj.Pkg().Path())
		b.WriteByte('.')
	}
	if sig := obj.Signature(); sig != nil && sig.Recv() != nil {
		b.WriteString(types.TypeString(sig.Recv().Type()))
		b.WriteByte('.')
	}
	b.WriteString(obj.Name())
	for _, a := range args {
		b.WriteByte('_')
		b.WriteString(types.TypeString(a))
	}
	return b.String()
}

// enqueue appends a shell when a call site is *discovered*, which is why
// Module.Funcs runs caller-before-callee and lower/vir has to reverse it.
//
// types.Info.Instances is not consulted as a call graph: it records only
// what the analyzer saw while checking each generic body once, generically.
// The worklist is built here, composing substitutions as it descends.
func (l *lowerer) enqueue(obj *types.Func, args []types.Type) *Func {
	key := instanceKey(obj, args)
	if f, ok := l.done[key]; ok {
		return f
	}

	pkg := obj.Pkg()
	mod := l.modOf[pkg]
	if mod == nil {
		l.bug("function " + obj.Name() + " belongs to no lowered module")
	}

	f := &Func{
		Name:   l.funcName(obj, args),
		Module: mod,
		Pos:    obj.Pos(),
		// Monomorphized instances are internal: vir has no linkonce, so two
		// modules instantiating smaller[int32] would collide in the flat
		// namespace. Each module emits its own copy.
		Export: len(args) == 0,
	}
	l.done[key] = f
	mod.Funcs = append(mod.Funcs, f)
	l.work = append(l.work, &instance{obj: obj, args: args, fn: f, pkg: pkg})
	return f
}

// drain lowers every queued instance, including the ones discovered while
// lowering. The depth guard is this package's own: semantics.md §9 makes
// non-terminating instantiation a compile error, and nothing here assumes
// the analyzer enforced it.
func (l *lowerer) drain() {
	for len(l.work) > 0 {
		in := l.work[0]
		l.work = l.work[1:]

		l.depth++
		if l.depth > maxInstantiationDepth {
			l.bug("instantiation depth exceeded at " + in.fn.Name +
				"; recursive instantiation must terminate")
		}
		l.lowerInstance(in)
		l.depth--
	}
}

func (l *lowerer) lowerInstance(in *instance) {
	saveMod, saveInfo, saveUnit := l.mod, l.info, l.unit
	l.mod = in.fn.Module
	l.info = l.infoOf[in.pkg]
	l.unit = l.unitOf(in.pkg)
	defer func() { l.mod, l.info, l.unit = saveMod, saveInfo, saveUnit }()

	decl := l.declOf(in.obj)
	if decl == nil {
		// No body: a foreign declaration, or a synthesized object. Foreign
		// ones were already emitted as externs.
		return
	}
	l.lowerFuncDecl(in, decl)
}

func (l *lowerer) unitOf(p *types.Package) *Unit {
	for _, u := range l.units {
		if u.Pkg == p {
			return u
		}
	}
	return nil
}

// declOf finds the syntax a Func object came from. types.Info records
// definitions by identifier, so this walks the package's files once and
// caches — the analyzer's funcObj map is not exported.
func (l *lowerer) declOf(obj *types.Func) *ast.FuncDecl {
	u := l.unitOf(obj.Pkg())
	if u == nil {
		return nil
	}
	for _, f := range u.Files {
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}
			if u.Info.Defs[fd.Name] == obj {
				return fd
			}
		}
	}
	return nil
}

// subst maps type parameters to concrete arguments while a generic body is
// lowered. It is composed on descent: instantiating Pair[T] inside
// Stack[int32]'s body substitutes T=int32 before enqueuing Pair[int32].
type subst struct {
	params []*types.TypeParam
	args   []types.Type
	outer  *subst
}

func (s *subst) apply(t types.Type) types.Type {
	if s == nil || t == nil {
		return t
	}
	if tp, ok := t.(*types.TypeParam); ok {
		for i, p := range s.params {
			if p == tp && i < len(s.args) {
				return s.outer.apply(s.args[i])
			}
		}
		return s.outer.apply(t)
	}
	// todo: a composite type mentioning a parameter — []T, map[K]V,
	// Stack[T] — needs a structural rebuild here. Today only bare parameter
	// positions substitute, which covers every generic in the corpus and
	// silently mis-lowers `func f[T](xs: []T)`.
	return t
}

func genericParams(obj *types.Func) []*types.TypeParam {
	sig := obj.Signature()
	if sig == nil || sig.Recv() == nil {
		return nil
	}
	if n := types.AsNamed(sig.Recv().Type()); n != nil {
		return n.TypeParams()
	}
	return nil
}