package vir

import (
	"strings"

	"github.com/vertex-language/vertex/lower/hir"
	"github.com/vertex-language/vertex/token"
	ir "github.com/vertex-language/vvm/ir/vir"
)

// func.go emits the function section: signatures, blocks, terminators.
//
// Like the struct section, the order is recomputed rather than preserved.
// hir's worklist appends a shell when a call site is *discovered*, so its
// Funcs slice runs caller-before-callee — exactly backwards for §2.2's
// declare-before-use.

func (l *lowerer) funcs(m *ir.Module, hm *hir.Module) {
	for _, f := range l.orderFuncs(hm) {
		l.function(m, hm, f)
	}
}

// orderFuncs returns hm's functions with every callee before its caller.
// §2.2 exempts direct self-recursion and nothing else, so a mutual
// recursion cycle has no vir spelling at all; it is reported here, where
// the cycle is actually known, rather than surfacing as an undeclared-name
// error out of ir/verify.
func (l *lowerer) orderFuncs(hm *hir.Module) []*hir.Func {
	byName := make(map[string]*hir.Func, len(hm.Funcs))
	for _, f := range hm.Funcs {
		byName[f.Name] = f
	}
	const (
		open = 1
		done = 2
	)
	state := make(map[*hir.Func]int, len(hm.Funcs))
	out := make([]*hir.Func, 0, len(hm.Funcs))

	var visit func(f *hir.Func, stack []*hir.Func)
	visit = func(f *hir.Func, stack []*hir.Func) {
		switch state[f] {
		case done:
			return
		case open:
			l.errorf(f.Pos, "mutual recursion (%s) has no vir spelling — §2.2 exempts only direct self-recursion from declare-before-use", cycleText(stack, f))
			return
		}
		state[f] = open
		for _, callee := range localCallees(f, byName) {
			if callee != f {
				visit(callee, append(stack, f))
			}
		}
		state[f] = done
		out = append(out, f)
	}
	for _, f := range hm.Funcs {
		visit(f, nil)
	}
	return out
}

// localCallees lists the same-module functions f calls, in first-call
// order so the emitted order is byte-reproducible. A qualified call names
// another module and is resolved by vvm's importer, not by declaration
// order here.
func localCallees(f *hir.Func, byName map[string]*hir.Func) []*hir.Func {
	var out []*hir.Func
	seen := map[*hir.Func]bool{}
	for _, b := range f.Blocks {
		for _, in := range b.Instr {
			if in.Op != hir.OpCall || in.Call == nil || in.Call.Module != "" {
				continue
			}
			g, ok := byName[in.Call.Name]
			if !ok || seen[g] {
				continue
			}
			seen[g] = true
			out = append(out, g)
		}
	}
	return out
}

func cycleText(stack []*hir.Func, f *hir.Func) string {
	names := make([]string, 0, len(stack)+1)
	for _, s := range stack {
		names = append(names, s.Name)
	}
	return strings.Join(append(names, f.Name), " -> ")
}

func (l *lowerer) function(m *ir.Module, hm *hir.Module, f *hir.Func) {
	var attrs []ir.FunctionAttribute
	if f.Entry {
		attrs = append(attrs, ir.AttributeEntry)
	}
	if f.NoReturn {
		attrs = append(attrs, ir.AttributeNoReturn)
	}
	// §2.2: both `entry` and `extern_c` require `export`.
	export := f.Export || f.Entry

	fb := m.DeclareFunction(f.Name, l.params(hm, f.Params), l.typ(hm, f.Result), export, attrs...)

	if f.Body != nil {
		l.errorf(f.Pos, "internal: %s still carries structured Body — hir.Flatten did not run", f.Name)
		return
	}
	if len(f.Blocks) == 0 {
		// A foreign declaration or an error-recovery shell reaches here
		// with nothing to emit; the extern group already named it.
		return
	}

	e := &emitter{l: l, hm: hm, fb: fb, line: -1}
	for i, b := range f.Blocks {
		switch {
		case i == 0 && b.Label != "":
			l.errorf(f.Pos, "internal: %s's first block is labeled %q — vir's entry block is implicit and unbranchable-to", f.Name, b.Label)
		case i > 0 && b.Label == "":
			l.errorf(f.Pos, "internal: %s has an unlabeled non-entry block", f.Name)
		case i > 0:
			fb.Label(b.Label)
			e.line = -1 // a new block, so re-anchor the loc run
		}
		for _, in := range b.Instr {
			e.instr(in)
		}
		if b.Term == nil {
			l.errorf(f.Pos, "internal: block %q of %s has no terminator", b.Label, f.Name)
			fb.Unreachable()
			continue
		}
		e.terminator(b.Term)
	}
}

// emitter carries the per-function state instruction emission needs: the
// builder, and the last source line a `loc` was written for.
type emitter struct {
	l    *lowerer
	hm   *hir.Module
	fb   *ir.FunctionBuilder
	line int
}

func (e *emitter) terminator(t hir.Terminator) {
	switch x := t.(type) {
	case hir.TermBranch:
		e.fb.Branch(x.Label)
	case hir.TermBranchIf:
		e.fb.BranchIf(e.l.value(e.hm, x.Cond), x.Then, x.Else)
	case hir.TermSwitch:
		cases := make([]ir.SwitchCase, 0, len(x.Cases))
		for _, c := range x.Cases {
			cases = append(cases, ir.SwitchCase{Value: c.Value, Label: c.Label})
		}
		// vir's switch takes a uniform operand/label list regardless of
		// density: jump-table-vs-compare-chain is cpu/lower/<arch>'s call.
		e.fb.Switch(e.l.value(e.hm, x.Value), x.Default, cases...)
	case hir.TermReturn:
		if x.Value == nil {
			e.fb.Return()
			return
		}
		e.fb.Return(e.l.value(e.hm, *x.Value))
	case hir.TermTrap:
		e.fb.Trap()
	case hir.TermUnreachable:
		e.fb.Unreachable()
	default:
		e.l.errorf(0, "internal: no vir terminator for %T", t)
		e.fb.Unreachable()
	}
}

// loc emits one debug line per source line, not per instruction. §2.3
// makes loc a body-line like any other, so it participates in ordering and
// costs nothing at runtime.
func (e *emitter) loc(pos token.Pos) {
	if !e.l.conf.Debug || e.l.conf.Fset == nil || !pos.IsValid() {
		return
	}
	p := e.l.conf.Fset.Position(pos)
	if p.Line == e.line {
		return
	}
	e.line = p.Line
	e.fb.Location(p.Filename, p.Line, p.Column)
}