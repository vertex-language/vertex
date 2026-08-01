package hir

// async.go holds the async/await state-machine split. It is not
// implemented, and this file exists to hold the shape and the constraint
// rather than to pretend otherwise.
//
// What it must do (overview §3):
//
//   - Rewrite each async-marked function into a poll loop plus a
//     synthesized payload enum of suspended states. Each await point
//     becomes a state boundary yielding Pending.
//   - Run a localized liveness pass so only variables surviving across an
//     await enter the enum, keeping the payload small.
//   - Synthesize a drop routine per state machine: if a task is cancelled,
//     its handle dies mid-flight, or the reactor is torn down at exit, the
//     payload enum still holds live owning variables whose deinits must run.
//   - Drive a bare `await someAsyncFn(x)` inline in the caller's poll loop,
//     merging state machines without touching the reactor or allocating a
//     channel. A launch-prefixed `async f(a, b)` bypasses the
//     transformation entirely and generates the channel handshake instead.
//
// Two ordering constraints bind whoever writes it:
//
//  1. It runs over still-structured control flow, before flattening, so the
//     split can see scopes.
//  2. It runs *before* defer/deinit epilogue expansion. Expanding epilogues
//     first would duplicate them across exit edges that the split then has
//     to cut apart and reassign. A suspend edge is not a scope exit:
//     `return Pending` explicitly bypasses these points, to preserve live
//     state.
//
// Constraint 2 is why the current arrangement is temporary. This package
// expands epilogues during body lowering, through the builder's scope
// stack, which is equivalent for a program containing no async functions
// and wrong for one that does. Landing this pass means first lifting
// epilogue expansion out of stmt.go into a standalone pass over the
// structured tree.
func (l *lowerer) splitAsync() {
	for _, m := range l.prog.Modules {
		for _, f := range m.Funcs {
			if !l.isAsync(f) {
				continue
			}
			l.todo(f.Pos, "async state-machine split for %s — see async.go", f.Name)
		}
	}
}

func (l *lowerer) isAsync(f *Func) bool {
	return l.markerOf(f) == "async"
}