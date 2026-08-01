package hir

import (
	"sort"
	"strings"

	"github.com/vertex-language/vertex/types"
)

// worklist drives monomorphization. types.Info.Instances records type
// arguments at each instantiation site, but only as the analyzer saw them
// while checking a generic body once, generically — it is not a concrete
// call graph. So hir builds its own worklist, seeded from every externally
// reachable root, composing substitutions as it descends and memoizing by
// (Func, ConcreteTypeArgs) so a diamond of instantiations is built once.
type worklist struct {
	l     *lowerer
	queue []job
	done  map[string]*Func
}

type job struct {
	unit  *Unit
	fn    *types.Func
	args  []types.Type
	depth int
}

func newWorklist(l *lowerer) *worklist {
	return &worklist{l: l, done: map[string]*Func{}}
}

func (w *worklist) key(fn *types.Func, args []types.Type) string {
	var sb strings.Builder
	if p := fn.Pkg(); p != nil {
		sb.WriteString(p.Path())
	}
	sb.WriteString(".")
	sb.WriteString(fn.Name())
	for _, a := range args {
		sb.WriteString("[")
		sb.WriteString(types.TypeString(a))
		sb.WriteString("]")
	}
	return sb.String()
}

// enqueue schedules one instantiation and returns the *Func it will occupy,
// creating the shell immediately so a recursive or mutually recursive call
// can reference it before its body exists.
//
// The shell's Result is resolved right here, eagerly, from fn's own
// checked signature — not left nil for lowerBody to fill in later. A nil
// Type reads as "not void" to IsVoid (the type assertion nil.(VoidType)
// fails), so any call reaching this shell before its body is lowered would
// otherwise wrongly bind a result name to what may be a void call. This is
// not a hypothetical: buildEntry calls the seeded main/test function from
// inside seed(), strictly before work.run() ever executes lowerBody, so
// the entry shim is guaranteed to observe a half-built shell on every
// build unless Result is settled up front. The signature is already fully
// resolved by the checker at this point, so there is no reason to defer
// this — lowerBody's own assignment of f.Result becomes a redundant
// re-computation of the same value once this runs first.
func (w *worklist) enqueue(u *Unit, fn *types.Func, args []types.Type, depth int) *Func {
	k := w.key(fn, args)
	if f, ok := w.done[k]; ok {
		return f
	}
	max := w.l.conf.MaxDepth
	if max == 0 {
		max = 64
	}
	if depth > max {
		// A.7.6: recursive instantiation must terminate; unbounded
		// deepening is a compile error, because the stamping-out would not.
		w.l.errorf(fn.Pos(), "recursive instantiation of %s exceeds depth %d (A.7.6)", fn.Name(), max)
		return nil
	}

	mod := w.l.byUnit[u]
	shell := &Func{
		Name:   mod.uniqueName(instanceName(fn, args)),
		Module: mod,
		Kind:   FuncUser,
		Origin: fn,
		Pos:    fn.Pos(),
	}

	// Resolve Result now, under this instantiation's own substitution,
	// exactly as lowerBody would later — but before any caller can
	// possibly observe the shell in a half-built state.
	prev := w.l.cur
	w.l.cur = &instance{unit: u, mod: mod, subst: substitution(fn, args), depth: depth}
	shell.Result = w.l.result(fn.Signature())
	w.l.cur = prev

	mod.Funcs = append(mod.Funcs, shell)
	w.done[k] = shell
	w.queue = append(w.queue, job{unit: u, fn: fn, args: args, depth: depth})
	return shell
}

func (w *worklist) run() {
	for len(w.queue) > 0 {
		j := w.queue[0]
		w.queue = w.queue[1:]
		w.lowerBody(j)
	}
}

// lookup finds an already-scheduled instance without scheduling one.
func (w *worklist) lookup(fn *types.Func, args []types.Type) *Func {
	return w.done[w.key(fn, args)]
}

func (w *worklist) lowerBody(j job) {
	f := w.lookup(j.fn, j.args)
	if f == nil {
		return
	}
	decl := findFuncDecl(j.unit, j.fn)
	if decl == nil || decl.Body == nil {
		// A foreign declaration or an error-recovery shell: nothing to
		// lower, and the extern group already named the entry point.
		return
	}

	subst := substitution(j.fn, j.args)
	prev := w.l.cur
	w.l.cur = &instance{unit: j.unit, mod: f.Module, subst: subst, depth: j.depth}
	defer func() { w.l.cur = prev }()

	// The receiver and the marker both live on the Signature, not on the
	// Func — A.4.2 makes the marker part of the callee's contract, so it is
	// part of the type. Everything below therefore reads through here.
	sig := j.fn.Signature()

	// f.Result was already resolved in enqueue, under the same
	// substitution, so this is a no-op recomputation of the same value —
	// left in place only because it is harmless and keeps lowerBody
	// self-contained if enqueue's early resolution is ever reverted.
	f.Result = w.l.result(sig)
	if st, ok := f.Result.(StructType); ok {
		// An aggregate result is written through an sret destination the
		// caller supplies, so the vir function returns void.
		f.Params = append(f.Params, &Param{Name: "sret", Type: Ptr, SRet: st.Def})
		f.Result = Void
	}
	if sig != nil {
		if recv := sig.Recv(); recv != nil {
			f.Params = append(f.Params, w.l.param("self", recv.Type()))
		}
	}
	for _, p := range paramsOf(sig) {
		f.Params = append(f.Params, w.l.param(p.Name, p.Type))
	}

	// vir has one flat namespace per module and spells every cross-module
	// call `module.symbol`, so a callee that is not exported cannot be named
	// from another module.
	//
	// There is no visibility fact for hir to read. ast.FuncDecl carries no
	// modifier list and types.Func records none; the grammar admits a
	// visibility modifier only inside a declare block (A.8.3), where it is
	// banned. So everything lowered is exported. That is sound rather than
	// merely convenient: monomorphization reaches a function only from a
	// root, so an unreachable function is never lowered and there is no
	// private symbol being leaked — the dead-symbol question was removed
	// upstream rather than answered here.
	f.Export = true

	b := newFuncBuilder(w.l, f, decl)
	b.body(decl.Body)
}

// substitution composes the concrete arguments onto the declaration's own
// parameter list.
func substitution(fn *types.Func, args []types.Type) map[*types.TypeParam]types.Type {
	if len(args) == 0 {
		return nil
	}
	params := typeParamsOf(fn)
	out := make(map[*types.TypeParam]types.Type, len(args))
	for i, p := range params {
		if i < len(args) {
			out[p] = args[i]
		}
	}
	return out
}

func instanceName(fn *types.Func, args []types.Type) string {
	var sb strings.Builder
	if sig := fn.Signature(); sig != nil && sig.Recv() != nil {
		sb.WriteString(sanitize(types.TypeString(sig.Recv().Type())))
		sb.WriteString("_")
	}
	sb.WriteString(fn.Name())
	for _, a := range args {
		sb.WriteString("__")
		sb.WriteString(sanitize(types.TypeString(a)))
	}
	return sb.String()
}

// sortedTypes gives a deterministic rendering for keys built from a map. A
// TypeParam carries its own name (A.7.1's TypeParameterList entry), so there
// is no declaring object to go through.
func sortedTypes(m map[*types.TypeParam]types.Type) []string {
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k.Name()+"="+types.TypeString(v))
	}
	sort.Strings(out)
	return out
}