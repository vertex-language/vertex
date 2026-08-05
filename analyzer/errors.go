package analyzer

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// Every error path in this package goes through these. Call sites never format
// their own text: diag's registry is what keeps one rule's message identical
// across every site that raises it, which is what makes the error form of
// ExpectedType stable enough to serve as specification.

func (c *Checker) report(d *diag.Diagnostic) {
	if d == nil {
		return
	}
	if d.Sev == diag.Error {
		c.errCount++
	}
	if c.conf.Reporter != nil {
		c.conf.Reporter.Report(d)
	}
}

func (c *Checker) errorAt(pos, end token.Pos, code diag.Code, args ...any) {
	c.report(diag.New(code, pos, end, args...))
}

// errorExpr reports over a whole node's extent, which is what a reader wants
// underlined: the construct, not its first token.
func (c *Checker) errorExpr(x ast.Node, code diag.Code, args ...any) {
	if x == nil {
		return
	}
	c.errorAt(x.Pos(), x.End(), code, args...)
}

// invalid is the recovery value for a failed type resolution.
//
// It is returned rather than nil so every consumer can keep going without a nil
// check, and types' predicates treat it as compatible with everything — one bad
// type produces one diagnostic, never a cascade. Same discipline ast.BadExpr
// follows on the parse side.
func (c *Checker) invalid() types.Type { return types.Typ[types.Invalid] }