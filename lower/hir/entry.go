package hir

import (
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/types"
)

// buildEntry synthesizes the program-entry wrapper (§5.2, tier 2).
//
// Vertex's main takes no parameters and returns nothing (A.6.1) but sets
// [+Await], so it may suspend. vir's entry is `export fn main() i32 entry`,
// and vvm's crt synthesizes the process-entry stub around it. So the thing
// in between — start the reactor, drive user main to completion, tear
// down, return a status — is synthesized here, and only when main's own
// marker says it can actually suspend. A program whose main never awaits
// gets no reactor wiring at all: no ReactorStart/ReactorShutdown calls, no
// FeatReactor need, and therefore no `import "reactor"` in the emitted
// module for a build that never asked for one.
//
// Under ModeTest the same slot holds a wrapper around one test function:
// start the reactor (if the test function suspends), drive it, render its
// result through fmt to console, tear down, return 0. One seam, two
// occupants, and both are tier 2.
func (l *lowerer) buildEntry(u *Unit, user *types.Func) {
	mod := l.byUnit[u]
	f := &Func{
		Name:   mod.uniqueName("main"),
		Module: mod,
		Result: I32,
		Export: true,
		Entry:  true,
		Kind:   FuncEntryShim,
	}
	mod.Funcs = append(mod.Funcs, f)
	l.prog.Entry = f

	prev := l.cur
	l.cur = &instance{unit: u, mod: mod}
	defer func() { l.cur = prev }()

	b := newFuncBuilder(l, f, nil)

	// Only wire in the reactor if main/the test function can actually
	// suspend (A.6.1's [+Await] marker). Every program pays the feature
	// floor regardless (see features.go); the reactor is not part of that
	// floor and must not be assumed.
	needsReactor := l.suspends(user)
	if needsReactor {
		l.need(builtins.FeatReactor)
		b.callBuiltin(0, builtins.ReactorStart, Void)
	}

	target := l.work.lookup(user, nil)
	if target == nil {
		target = l.work.enqueue(u, user, nil, 0)
	}
	if target != nil {
		v := b.call(0, target, nil...)
		if l.conf.Mode == ModeTest {
			l.renderResult(b, user, v)
		}
	}

	if needsReactor {
		b.callBuiltin(0, builtins.ReactorShutdown, Void)
	}
	b.seq.add(&Return{Value: ptrTo(IntVal(I32, 0))})
}

// renderResult emits the test wrapper's single fmt render to console, with
// no trailing newline and no decoration.
//
// This is why fmt is normative rather than convenient: the expected-value
// literals in source are defined as the exact bytes fmt produces, so
// Expected(float32, "5.000000") is a statement about %f's six-digit default
// and nothing else. A.12.2's table picks the verb from the Expected type
// argument.
//
// The three channels stay separate: stdout carries the rendered value and
// nothing else, stderr carries panic text, and exit status carries crash
// detection — which is the entire mechanism behind a bare test function
// with no Expected.
func (l *lowerer) renderResult(b *funcBuilder, user *types.Func, v Value) {
	rt := l.testResultType(user)
	if rt == nil {
		return // a bare test: exit status only, stdout ignored
	}
	var sym builtins.Symbol
	switch l.classify(rt) {
	case kBool:
		sym = builtins.FmtBool
	case kFloat:
		sym = builtins.FmtFloat
	case kString:
		sym = builtins.FmtString
	case kChar:
		sym = builtins.FmtChar
	case kInt:
		if l.isSigned(rt) {
			sym = builtins.FmtInt
		} else {
			sym = builtins.FmtUint
		}
	default:
		l.todo(user.Pos(), "Expected(%s, ...) has no fmt verb in A.12.2's table", types.TypeString(rt))
		return
	}
	l.need(builtins.FeatFmt)
	rendered := b.callBuiltin(0, sym, l.hirType(types.Typ[types.String]), v)
	b.callBuiltin(0, builtins.ConsoleOut, Void, rendered)
}

func ptrTo(v Value) *Value { return &v }