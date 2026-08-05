package hir

// The async/await state-machine split. Not implemented.
//
// The shape is fixed and the reason it is not written is ordering, not
// difficulty. An async function becomes three things:
//
//  1. struct _Vframe_f(state i32, <locals live across a suspend>,
//     <child frames>, result)
//  2. fn _Vresume_f(frame ptr) i32 — one block per resume point, dispatched
//     by a switch on state at entry. 0 = complete, 1 = suspended.
//  3. a stub fn f(…) that initializes the frame.
//
// State lives in the frame in memory, so there are no phi nodes to
// reconstruct and no registers to restore across a suspend. A localized
// liveness pass extracts only the variables surviving across an await, which
// is the whole of the memory-footprint claim.
//
// `await g()` where g is async embeds _Vframe_g inside _Vframe_f and calls
// _Vresume_g directly, so one allocation covers a whole await chain — and
// therefore recursive async is a compile error, since the frame size would
// not be computable. Vertex has no boxed-future type and should not grow
// one; spawn instead, which allocates a fresh frame.
//
// ---------------------------------------------------------------------------
//
// What has to happen first, and why this file is empty rather than partial:
//
// This package performs defer/deinit epilogue expansion *during* body
// lowering, through funcBuilder's scope stack (openScope, epilogueTo,
// closeScope). lowering.md specifies it as a pass over the still-structured
// tree, running after the split. The two are equivalent for a program
// containing no async functions and wrong for one that does, because a
// suspend edge is not a scope exit: a `return Pending` must bypass the
// epilogue entirely to preserve live state, and an epilogue already inlined
// at every exit edge cannot tell the two kinds of edge apart.
//
// So landing async means, in order:
//
//  1. Lift epilogue expansion out of stmt.go into a pass over *Seq. The seam
//     is three methods on funcBuilder — openScope, epilogueTo, closeScope —
//     which is why it was written as three methods.
//  2. Run the split before that pass, over the structured tree, so the split
//     can still see scopes.
//  3. Synthesize the drop routine each state machine needs: if a task is
//     cancelled, or its handle dies mid-flight, or the reactor is torn down
//     at exit, the payload frame still holds live owning variables whose
//     teardown must run.
//
// The reactor that drives the machines is not synthesized here. It is a
// builtins module, reached through the task ABI — {poll fnptr, drop fnptr,
// state ptr} — which is the one indirect dispatch the compiler emits and the
// one thing standing between this file and lower/vir's missing fnsig
// support.