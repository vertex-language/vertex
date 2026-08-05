// func.go
package vir

import (
	"strings"

	"github.com/vertex-language/vertex/lower/hir"
	"github.com/vertex-language/vertex/token"
	vir "github.com/vertex-language/vvm/ir/vir"
)

// functions emits the function section, callee-first.
func (ml *moduleLowerer) functions() {
	for _, f := range ml.orderFuncs() {
		ml.function(f)
	}
}

// orderFuncs reverses hir's discovery order.
//
// hir's worklist appends a shell when a call site is *discovered*, so its
// Funcs slice runs caller-before-callee — exactly backwards for vir §2.2,
// which has no forward references. A plain reversal is right for a tree of
// calls and wrong for a diamond, so this is a real post-order walk with the
// reversal used only as the visit order, which keeps the emitted text
// stable against hir's ordering.
//
// vir §2.2 exempts *direct* self-recursion and nothing else, so mutual
// recursion has no vir spelling at all. It is reported here, where the
// cycle is known, rather than reaching ir/verify as an undeclared-name
// error. That is a language-level limitation, not a gap in this package:
// vir has no fnsig for a defined function to forward-declare it with, so
// either vir grows forward declarations for fn, or Vertex does not have
// mutually recursive functions.
func (ml *moduleLowerer) orderFuncs() []*hir.Func {
	byName := make(map[string]*hir.Func, len(ml.src.Funcs))
	for _, f := range ml.src.Funcs {
		byName[f.Name] = f
	}

	state := make(map[*hir.Func]int, len(ml.src.Funcs))
	out := make([]*hir.Func, 0, len(ml.src.Funcs))
	var stack []*hir.Func

	var visit func(f *hir.Func)
	visit = func(f *hir.Func) {
		switch state[f] {
		case visited:
			return
		case visiting:
			ml.bug("mutually recursive functions have no vir spelling: " +
				funcCycle(stack, f) + " (vir §2.2 exempts direct self-recursion only)")
		}
		state[f] = visiting
		stack = append(stack, f)
		for _, c := range ml.calleesOf(f, byName) {
			if c == f {
				continue // direct self-recursion is exempt
			}
			visit(c)
		}
		stack = stack[:len(stack)-1]
		state[f] = visited
		out = append(out, f)
	}

	for i := len(ml.src.Funcs) - 1; i >= 0; i-- {
		visit(ml.src.Funcs[i])
	}
	return out
}

// calleesOf collects this module's own callees. A qualified call names
// another module and is resolved by the importer, not by section order.
func (ml *moduleLowerer) calleesOf(f *hir.Func, byName map[string]*hir.Func) []*hir.Func {
	var out []*hir.Func
	seen := map[*hir.Func]bool{}
	add := func(name string) {
		if c, ok := byName[name]; ok && !seen[c] {
			seen[c] = true
			out = append(out, c)
		}
	}
	for _, b := range f.Blocks {
		for _, in := range b.Lines {
			if in.Op == hir.OpCall && in.Module == "" && in.Callee != "" {
				add(in.Callee)
			}
		}
	}
	return out
}

func funcCycle(stack []*hir.Func, f *hir.Func) string {
	names := make([]string, 0, len(stack)+1)
	start := 0
	for i, x := range stack {
		if x == f {
			start = i
			break
		}
	}
	for _, x := range stack[start:] {
		names = append(names, x.Name)
	}
	return strings.Join(append(names, f.Name), " -> ")
}

// ---------------------------------------------------------------- one body

func (ml *moduleLowerer) function(f *hir.Func) {
	if f.Body != nil {
		// Flatten is hir's to run, and Lower runs it. A structured body
		// reaching here means the pass order changed and nothing said so.
		ml.bug("function " + f.Name + " reached lowering unflattened")
	}
	if len(f.Blocks) == 0 {
		// A shell with no body: a foreign declaration, or an object the
		// worklist enqueued and monomorphization found no syntax for.
		// Foreign ones were already emitted into an extern group; emitting
		// an empty fn here would produce a block with no terminator.
		return
	}

	var attrs []vir.FunctionAttribute
	if f.Entry {
		// entry forces a bare symbol even in a namespaced module, which is
		// why hir names the shim _Ventry and lets the symbol come out as
		// main without colliding with the user's own ident.
		attrs = append(attrs, vir.AttributeEntry)
	}
	if f.NoReturn {
		attrs = append(attrs, vir.AttributeNoReturn)
	}

	// An aggregate result left through the sret parameter, so Result is nil
	// and the vir return type is void.
	b := ml.out.DeclareFunction(f.Name, ml.params(f.Params), ml.typ(f.Result), f.Export, attrs...)
	if f.Variadic {
		b.SetVariadic()
	}

	fl := &funcLowerer{ml: ml, fn: f, b: b}
	fl.body()
}

type funcLowerer struct {
	ml *moduleLowerer
	fn *hir.Func
	b  *vir.FunctionBuilder

	// last is the most recently emitted loc, so a run of instructions from
	// one source line costs one line of text rather than twenty.
	last position
}

type position struct {
	file string
	line int
	col  int
}

func (fl *funcLowerer) body() {
	blocks := fl.fn.Blocks
	if blocks[0].Label != "" {
		// vir's entry block is implicit, unlabeled, and unbranchable-to.
		fl.ml.bug("entry block of " + fl.fn.Name + " carries a label")
	}

	fl.loc(fl.fn.Pos)
	fl.block(blocks[0])

	for _, blk := range blocks[1:] {
		if blk.Label == "" {
			fl.ml.bug("unlabeled non-entry block in " + fl.fn.Name)
		}
		fl.b.Label(blk.Label)
		fl.block(blk)
	}
}

func (fl *funcLowerer) block(b *hir.FlatBlock) {
	for _, in := range b.Lines {
		fl.loc(in.Pos)
		fl.instr(in)
	}
	fl.term(b)
}

func (fl *funcLowerer) term(b *hir.FlatBlock) {
	switch t := b.Term.(type) {
	case hir.Br:
		fl.b.Branch(t.Label)
	case hir.BrIf:
		fl.b.BranchIf(fl.ml.operand(t.Cond), t.Then, t.Else)
	case hir.SwitchTerm:
		cases := make([]vir.SwitchCase, 0, len(t.Cases))
		for _, c := range t.Cases {
			cases = append(cases, vir.SwitchCase{Value: c.Value, Label: c.Label})
		}
		fl.b.Switch(fl.ml.operand(t.Value), t.Default, cases...)
	case hir.Ret:
		if t.Value != nil {
			fl.b.Return(fl.ml.operand(*t.Value))
			return
		}
		fl.b.Return()
	case hir.TrapTerm:
		fl.b.Trap()
	case hir.UnreachTerm:
		fl.b.Unreachable()
	case nil:
		fl.ml.bug("block " + b.Label + " in " + fl.fn.Name + " has no terminator")
	default:
		fl.ml.bug("terminator with no vir spelling in " + fl.fn.Name)
	}
}

// loc emits a debug line. hir resolves no positions itself — it carries
// token.Pos and hands the FileSet down, which is the whole reason Config
// takes one.
func (fl *funcLowerer) loc(pos token.Pos) {
	if fl.ml.conf.Fset == nil || pos == token.NoPos {
		return
	}
	p := fl.ml.conf.Fset.Position(pos)
	cur := position{file: p.Filename, line: p.Line, col: p.Column}
	if cur == fl.last {
		return
	}
	fl.last = cur
	fl.b.Location(cur.file, cur.line, cur.col)
}