package compiler

import (
	"github.com/vertex-language/vertex/parser"
)

// ── Public API ────────────────────────────────────────────────────────────────

// DCEResult is returned by ComputeReachable.
type DCEResult struct {
	// Reachable is the set of every function/extern name that is live.
	Reachable map[string]bool
}

// IsReachable reports whether name is live (always true when DCE is disabled).
func (r *DCEResult) IsReachable(name string) bool {
	if r == nil {
		return true
	}
	return r.Reachable[name]
}

// ComputeReachable performs call-graph–based dead-code elimination starting
// from entryName. It walks the Vertex AST of every reachable function body,
// collecting call sites accurately (not just token text). Extern functions are
// treated as opaque leaves — they are marked reachable when called but their
// (C) bodies are never scanned.
// ComputeReachable performs call-graph–based dead-code elimination starting
// from entryName. It walks the Vertex AST of every reachable function body,
// collecting call sites accurately (not just token text). Extern functions are
// treated as opaque leaves — they are marked reachable when called but their
// (C) bodies are never scanned.
func ComputeReachable(pkg *Package, entryName string) *DCEResult {
	w := &dceWalker{
		defs:      make(map[string]*dceEntry),
		reachable: make(map[string]bool),
	}

	// Pass 1: catalog all top-level definitions.
	for _, sf := range pkg.Files {
		for _, tld := range sf.Tree.AllTopLevelDecl() {
			if fn := tld.FuncDecl(); fn != nil {
				name := fn.ID().GetText()
				w.defs[name] = &dceEntry{funcCtx: fn}
			}
			if cd := tld.ClassDecl(); cd != nil && cd.COLON() != nil {
				for _, member := range cd.AllClassMember() {
					if nfd := member.NativeFuncDecl(); nfd != nil {
						name := nfd.ID().GetText()
						w.defs[name] = &dceEntry{isExtern: true}
					}
				}
			}
		}
	}

	// Pass 2: worklist traversal from the entry point.
	w.mark(entryName)
	for len(w.worklist) > 0 {
		curr := w.worklist[len(w.worklist)-1]
		w.worklist = w.worklist[:len(w.worklist)-1]

		entry, ok := w.defs[curr]
		if !ok || entry.isExtern {
			continue
		}
		w.scanFunc(entry.funcCtx)
	}

	return &DCEResult{Reachable: w.reachable}
}

// ── Internal types ────────────────────────────────────────────────────────────

type dceEntry struct {
	funcCtx  parser.IFuncDeclContext
	isExtern bool
}

type dceWalker struct {
	defs      map[string]*dceEntry
	reachable map[string]bool
	worklist  []string
}

func (w *dceWalker) mark(name string) {
	if w.reachable[name] {
		return
	}
	w.reachable[name] = true
	// Only enqueue if we have a definition to scan.
	if _, ok := w.defs[name]; ok {
		w.worklist = append(w.worklist, name)
	}
}

// ── Function / block / statement scanners ────────────────────────────────────

func (w *dceWalker) scanFunc(ctx parser.IFuncDeclContext) {
	if ctx == nil || ctx.Block() == nil {
		return
	}
	w.scanBlock(ctx.Block())
}

func (w *dceWalker) scanBlock(ctx parser.IBlockContext) {
	if ctx == nil {
		return
	}
	for _, stmt := range ctx.AllStmt() {
		w.scanStmt(stmt)
	}
}

func (w *dceWalker) scanStmt(ctx parser.IStmtContext) {
	if ctx == nil {
		return
	}
	switch {
	case ctx.VarDeclStmt() != nil:
		if e := ctx.VarDeclStmt().Expr(); e != nil {
			w.scanExpr(e)
		}

	case ctx.AssignStmt() != nil:
		w.scanExpr(ctx.AssignStmt().Expr())

	case ctx.CompoundAssignStmt() != nil:
		w.scanExpr(ctx.CompoundAssignStmt().Expr())

	case ctx.ExprStmt() != nil:
		w.scanExpr(ctx.ExprStmt().Expr())

	case ctx.ReturnStmt() != nil:
		if e := ctx.ReturnStmt().Expr(); e != nil {
			w.scanExpr(e)
		}

	case ctx.DeferStmt() != nil:
		if e := ctx.DeferStmt().Expr(); e != nil {
			w.scanExpr(e)
		}

	case ctx.IfStmt() != nil:
		w.scanIf(ctx.IfStmt())

	case ctx.WhileStmt() != nil:
		w.scanWhile(ctx.WhileStmt())

	case ctx.ForInStmt() != nil:
		w.scanForIn(ctx.ForInStmt())

	case ctx.SwitchStmt() != nil:
		w.scanSwitch(ctx.SwitchStmt())

	// break / continue carry no call sites.
	}
}

func (w *dceWalker) scanIf(ctx parser.IIfStmtContext) {
	if ctx == nil {
		return
	}
	if cond := ctx.IfCondition(); cond != nil && cond.Expr() != nil {
		w.scanExpr(cond.Expr())
	}
	w.scanBlock(ctx.Block())

	for _, eic := range ctx.AllElseIfClause() {
		if cond := eic.IfCondition(); cond != nil && cond.Expr() != nil {
			w.scanExpr(cond.Expr())
		}
		w.scanBlock(eic.Block())
	}
	if ec := ctx.ElseClause(); ec != nil {
		w.scanBlock(ec.Block())
	}
}

func (w *dceWalker) scanWhile(ctx parser.IWhileStmtContext) {
	if ctx == nil {
		return
	}
	w.scanExpr(ctx.Expr())
	w.scanBlock(ctx.Block())
}

func (w *dceWalker) scanForIn(ctx parser.IForInStmtContext) {
	if ctx == nil {
		return
	}
	w.scanExpr(ctx.Expr())
	w.scanBlock(ctx.Block())
}

func (w *dceWalker) scanSwitch(ctx parser.ISwitchStmtContext) {
	if ctx == nil {
		return
	}
	w.scanExpr(ctx.Expr())
	for _, sc := range ctx.AllSwitchCase() {
		for _, stmt := range sc.AllStmt() {
			w.scanStmt(stmt)
		}
	}
	if dc := ctx.DefaultCase(); dc != nil {
		for _, stmt := range dc.AllStmt() {
			w.scanStmt(stmt)
		}
	}
}

// ── Expression scanner ────────────────────────────────────────────────────────

// scanExpr recursively walks an expression, marking every function name that
// appears in a genuine call position.
func (w *dceWalker) scanExpr(ctx parser.IExprContext) {
	if ctx == nil {
		return
	}

	exprs := ctx.AllExpr()

	// ── Direct call: f(args)  ─────────────────────────────────────────────
	// Grammar shape: one sub-expr, an ArgList, no DOT.
	if ctx.ArgList() != nil && ctx.DOT() == nil && len(exprs) == 1 {
		if name := dceIdentOf(exprs[0]); name != "" {
			w.mark(name)
		}
		w.scanExpr(exprs[0])
		w.scanArgList(ctx.ArgList())
		return
	}

	// ── Method / package-qualified call: recv.method(args) ───────────────
	// Grammar shape: one sub-expr, DOT, PostfixName, ArgList.
	if ctx.DOT() != nil && ctx.PostfixName() != nil &&
		ctx.ArgList() != nil && len(exprs) == 1 {

		method := ctx.PostfixName().GetText()
		// Mark both "TypeName.method" (static dispatch) and bare "method"
		// so either lookup strategy in codegen finds it.
		if typeName := dceIdentOf(exprs[0]); typeName != "" {
			w.mark(typeName + "." + method)
		}
		w.mark(method)
		w.scanExpr(exprs[0])
		w.scanArgList(ctx.ArgList())
		return
	}

	// ── Postfix field access / chained member (no call) ──────────────────
	if ctx.DOT() != nil && ctx.PostfixName() != nil && len(exprs) == 1 {
		w.scanExpr(exprs[0])
		return
	}

	// ── Subscript: arr[idx] ───────────────────────────────────────────────
	if ctx.LBRACKET() != nil && len(exprs) == 2 {
		w.scanExpr(exprs[0])
		w.scanExpr(exprs[1])
		return
	}

	// ── Ternary: cond ? a : b ─────────────────────────────────────────────
	if ctx.QUESTION() != nil && ctx.COLON() != nil && len(exprs) == 3 {
		for _, e := range exprs {
			w.scanExpr(e)
		}
		return
	}

	// ── Nil-coalesce: a ?? b ──────────────────────────────────────────────
	if ctx.NIL_COALESCE() != nil && len(exprs) == 2 {
		w.scanExpr(exprs[0])
		w.scanExpr(exprs[1])
		return
	}

	// ── Type cast: PrimitiveType(expr) ───────────────────────────────────
	if ctx.PrimitiveType() != nil && len(exprs) == 1 {
		w.scanExpr(exprs[0])
		return
	}

	// ── Primary leaf ──────────────────────────────────────────────────────
	if p := ctx.Primary(); p != nil && len(exprs) == 0 {
		w.scanPrimary(p)
		return
	}

	// ── Unary / binary fallthrough ────────────────────────────────────────
	for _, e := range exprs {
		w.scanExpr(e)
	}
	// Trailing ArgList on something we didn't pattern-match above.
	if ctx.ArgList() != nil {
		w.scanArgList(ctx.ArgList())
	}
}

func (w *dceWalker) scanPrimary(ctx parser.IPrimaryContext) {
	if ctx == nil {
		return
	}
	// Parenthesised expression.
	if ctx.Expr() != nil {
		w.scanExpr(ctx.Expr())
	}
	// Array literal: [a, b, c]
	if al := ctx.ArrayLiteralExpr(); al != nil {
		if el := al.ExprList(); el != nil {
			for _, e := range el.AllExpr() {
				w.scanExpr(e)
			}
		}
	}
	// Anonymous / lambda function — walk its body too so closures that call
	// other functions are not accidentally eliminated.
	if af := ctx.AnonFuncExpr(); af != nil {
		if blk := af.Block(); blk != nil {
			w.scanBlock(blk)
		}
	}
}

func (w *dceWalker) scanArgList(ctx parser.IArgListContext) {
	if ctx == nil {
		return
	}
	for _, arg := range ctx.AllArg() {
		w.scanExpr(arg.Expr())
	}
}

// ── Helper ────────────────────────────────────────────────────────────────────

// dceIdentOf returns the bare identifier text if ctx is a simple identifier
// primary (e.g. `foo`), otherwise "". This is the crucial guard that prevents
// the old token-matching false positives: only expressions that are in callee
// position are examined.
func dceIdentOf(ctx parser.IExprContext) string {
	if ctx == nil {
		return ""
	}
	// Must be a primary with no sub-expressions.
	p := ctx.Primary()
	if p == nil || p.ID() == nil || len(ctx.AllExpr()) != 0 {
		return ""
	}
	return p.ID().GetText()
}