package hir

import (
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/types"
)

// The entry shim. Vertex's main takes no parameters and returns nothing but
// may await, while vir's entry is `export fn main() i32 entry` and vvm's crt
// synthesizes the process-entry stub around it. So hir synthesizes a wrapper
// that starts the reactor, drives user main to completion, tears down, and
// returns a status.
//
// The shim is named _Ventry rather than main: vir gives an `entry` export a
// bare symbol even in a namespaced module, so the emitted symbol is `main`
// either way, and user main keeps its own ident without colliding.
//
// One seam, two occupants — under ModeTest the same slot holds a wrapper
// around one test function.

func (l *lowerer) mainShim(user *types.Func) *Func {
	mod := l.modOf[user.Pkg()]
	l.prog.Root = mod

	fn := &Func{Name: "_Ventry", Module: mod, Result: I32, Export: true, Entry: true}
	mod.Funcs = append(mod.Funcs, fn)

	fb := &funcBuilder{l: l, fn: fn, seq: &Seq{}, sub: &subst{}}
	fb.openScope(scopeFunc)

	// todo: the reactor. main sets async — it is the one non-async function
	// in which await is legal — so the real shim starts the reactor, drives
	// main's state machine to completion, and tears down. Without async.go
	// there is no state machine to drive, so this calls main directly and
	// the reactor never starts.
	target := l.enqueue(user, nil)
	fb.emitCall(target.Name, moduleNameOf(target, mod), Void, nil)

	zero := Int(0, I32)
	fb.seq.add(&ReturnStmt{Value: &zero})
	fb.closeScope()
	fn.Body = fb.seq
	return fn
}

// testShim renders one test function's result and exits.
//
// The Expected type argument selects the format verb, and the wrapper emits
// a single render of the returned value to stdout with no trailing newline
// and no decoration. This is why fmt is normative rather than convenient:
// the expected-value literals in source are defined as the exact bytes it
// produces, so Expected(float32, "5.000000") is a statement about %f's
// six-digit default and nothing else.
//
// stdout carries the rendered value and nothing else; stderr carries panic
// text; the exit status carries crash detection, which is the entire
// mechanism behind a bare test function with no Expected.
func (l *lowerer) testShim(user *types.Func) *Func {
	mod := l.modOf[user.Pkg()]
	l.prog.Root = mod

	fn := &Func{Name: "_Ventry", Module: mod, Result: I32, Export: true, Entry: true}
	mod.Funcs = append(mod.Funcs, fn)

	fb := &funcBuilder{l: l, fn: fn, seq: &Seq{}, sub: &subst{}}
	fb.openScope(scopeFunc)

	target := l.enqueue(user, nil)
	sig := user.Signature()

	exp := sig.Expected()
	switch {
	case exp == nil, exp.IsError():
		// A bare test passes if it returns; an error test never reaches
		// lowering at all, since it stops at the analyzer.
		fb.emitCall(target.Name, moduleNameOf(target, mod), Void, nil)

	default:
		rt := fb.typ(exp.Type())
		v := fb.emitCall(target.Name, moduleNameOf(target, mod), rt, nil)
		fb.render(v, exp.Type())
	}

	zero := Int(0, I32)
	fb.seq.add(&ReturnStmt{Value: &zero})
	fb.closeScope()
	fn.Body = fb.seq
	return fn
}

// render emits one fmt call, chosen by the result type's verb.
func (fb *funcBuilder) render(v Value, t types.Type) {
	fb.l.need(builtins.FeatFmt)
	var sym builtins.Symbol
	switch {
	case types.IsString(t):
		sym = builtins.FmtString
	case types.IsFloat(t):
		sym = builtins.FmtFloat
	case types.IsBool(t):
		sym = builtins.FmtBool
	case types.IsChar(t):
		sym = builtins.FmtChar
	case types.IsInteger(t):
		// Signedness is a property of the lowered type, so it is read there
		// rather than asked of types twice.
		if fb.typ(t).Signed {
			sym = builtins.FmtInt
		} else {
			sym = builtins.FmtUint
		}
	default:
		fb.l.todoAt(fb.fn.Pos, "rendering "+types.TypeString(t))
	}
	fb.callBuiltin(sym, Void, v)
}