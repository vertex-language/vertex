package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
)

// Modules — vertex_grammar.md section J.

// parseModuleItem is ModuleItem (J).
func (p *parser) parseModuleItem() ast.Stmt {
	switch p.kind() {
	case token.IMPORT:
		return p.parseImportDecl()
	case token.EXPORT:
		return p.parseExportDecl()
	}
	return p.parseStmt()
}

func (p *parser) parseModuleItems(stop token.Kind) []ast.Stmt {
	var out []ast.Stmt
	for !p.at(stop) && !p.atEOF() {
		before := p.i
		s := p.parseModuleItem()
		if s != nil {
			out = append(out, s)
		}
		if !p.advanced(before) {
			if !p.advanceTo(syncDecl) {
				break
			}
		}
	}
	return out
}

// --- imports (J.1) ----------------------------------------------------------

func (p *parser) parseImportDecl() ast.Stmt {
	importPos := p.next().Pos

	// ImportCall and ImportMeta are expressions, not declarations. `import(` and
	// `import.` in statement position belong to ExpressionStatement.
	if p.at(token.LPAREN) || p.at(token.PERIOD) {
		p.i-- // give the `import` back; parsePrimary re-reads it
		return p.parseExprStmt()
	}

	// ImportEqualsDeclaration (J.1): `import X = Y;`.
	if p.at(token.IDENT) && p.peek(1).Kind == token.ASSIGN {
		d := &ast.ImportEqualsDecl{ImportPos: importPos, Name: p.parseIdent()}
		d.Assign = p.next().Pos
		if p.atCtx(token.CtxRequire) && p.peek(1).Kind == token.LPAREN {
			d.RequirePos = p.next().Pos
			p.expect(token.LPAREN)
			if p.at(token.STRING) {
				t := p.next()
				d.Path = &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
			} else {
				p.errorf(p.cur(), "expected a module specifier string")
			}
			d.Rparen = p.expect(token.RPAREN)
		} else {
			d.Entity = p.parseQualifiedName()
		}
		d.Semi = p.expectSemi()
		return d
	}

	d := &ast.ImportDecl{ImportPos: importPos, Phase: ast.PhaseEval}

	// `import ModuleSpecifier;` — the side-effect form, no clause.
	if p.at(token.STRING) {
		t := p.next()
		d.Path = &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
		d.With = p.tryWithClause()
		d.Semi = p.expectSemi()
		return d
	}

	switch {
	case p.atCtx(token.CtxDefer) && p.peek(1).Kind == token.MUL:
		// `import defer NameSpaceImport FromClause` — namespace form only.
		d.Phase, d.PhasePos = ast.PhaseDefer, p.next().Pos
	case p.atCtx(token.CtxSource) && p.peek(1).Kind == token.IDENT:
		d.Phase, d.PhasePos = ast.PhaseSource, p.next().Pos
	case p.atCtx(token.CtxType) && p.importClauseFollows(1):
		// `import type ImportClause FromClause` (J.1).
		d.TypePos = p.next().Pos
	}

	p.parseImportClause(d)

	d.FromPos = p.expectCtx(token.CtxFrom)
	if p.at(token.STRING) {
		t := p.next()
		d.Path = &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
	} else {
		p.errorf(p.cur(), "expected a module specifier string after `from`")
	}
	d.With = p.tryWithClause()
	d.Semi = p.expectSemi()
	return d
}

// importClauseFollows distinguishes `import type X from "m"` from
// `import type from "m"`, where `type` is the imported binding.
func (p *parser) importClauseFollows(n int) bool {
	switch t := p.peek(n); t.Kind {
	case token.MUL, token.LBRACE:
		return true
	case token.IDENT:
		return t.Ctx != token.CtxFrom
	}
	return false
}

func (p *parser) parseImportClause(d *ast.ImportDecl) {
	switch {
	case p.at(token.MUL):
		d.Namespace = p.parseNamespaceImport()
		return
	case p.at(token.LBRACE):
		d.NamedPos, d.Named, d.NamedEnd = p.parseNamedImports()
		return
	}

	d.Default = p.parseIdent()
	if !p.got(token.COMMA) {
		return
	}
	switch {
	case p.at(token.MUL):
		d.Namespace = p.parseNamespaceImport()
	case p.at(token.LBRACE):
		d.NamedPos, d.Named, d.NamedEnd = p.parseNamedImports()
	default:
		p.errorf(p.cur(), "expected `*` or `{` after the default import")
	}
}

func (p *parser) parseNamespaceImport() *ast.Ident {
	p.expect(token.MUL)
	p.expectCtx(token.CtxAs)
	return p.parseIdent()
}

func (p *parser) parseNamedImports() (token.Pos, []*ast.ImportSpec, token.Pos) {
	lb := p.expect(token.LBRACE)
	var specs []*ast.ImportSpec
	for !p.at(token.RBRACE) && !p.atEOF() {
		before := p.i
		s := &ast.ImportSpec{}
		if p.atCtx(token.CtxType) && p.peek(1).Kind != token.COMMA &&
			p.peek(1).Kind != token.RBRACE && p.peek(1).Ctx != token.CtxAs {
			s.TypePos = p.next().Pos
		}
		name := p.parseModuleExportName()
		if p.atCtx(token.CtxAs) {
			s.Name = name
			s.AsPos = p.next().Pos
			s.Local = p.parseIdent()
		} else {
			id, ok := name.(*ast.Ident)
			if !ok {
				p.errorAt(name.Pos(), name.End(), "a string import name requires `as`")
				id = &ast.Ident{NamePos: name.Pos(), NameEnd: name.End()}
			}
			s.Local = id
		}
		specs = append(specs, s)
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	return lb, specs, p.expect(token.RBRACE)
}

// parseModuleExportName is ModuleExportName (J.1): an IdentifierName or a
// StringLiteral.
func (p *parser) parseModuleExportName() ast.Node {
	if p.at(token.STRING) {
		t := p.next()
		return &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
	}
	return p.parseIdentName()
}

func (p *parser) tryWithClause() *ast.WithClause {
	// `with` is a ReservedWord (A.2, see the note on parseWithClause below),
	// so it always scans as token.WITH and never as an IDENT carrying a Ctx.
	// There is no token.CtxWith — checking atCtx here would be dead code and
	// referencing the constant doesn't compile, since it was never defined.
	if !p.at(token.WITH) {
		return nil
	}
	return p.parseWithClause()
}

// parseWithClause is WithClause (J.1).
//
// `with` is a ReservedWord (A.2) with no statement production anywhere in the
// grammar, so it arrives here as a Kind rather than a Ctx.
func (p *parser) parseWithClause() *ast.WithClause {
	w := &ast.WithClause{WithPos: p.next().Pos}
	w.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.atEOF() {
		before := p.i
		a := &ast.ImportAttr{Key: p.parseAttributeKey()}
		a.Colon = p.expect(token.COLON)
		if p.at(token.STRING) {
			t := p.next()
			a.Value = &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
		} else {
			p.errorf(p.cur(), "an import attribute value must be a string literal")
			a.Value = &ast.BasicLit{Kind: token.STRING, ValuePos: p.pos(), ValueEnd: p.end()}
		}
		w.List = append(w.List, a)
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	w.Rbrace = p.expect(token.RBRACE)
	return w
}

func (p *parser) parseAttributeKey() ast.Node {
	if p.at(token.STRING) {
		t := p.next()
		return &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
	}
	return p.parseIdentName()
}

// --- exports (J.2) ----------------------------------------------------------

func (p *parser) parseExportDecl() ast.Stmt {
	d := &ast.ExportDecl{ExportPos: p.next().Pos}

	switch {
	case p.at(token.DEFAULT):
		d.DefaultPos = p.next().Pos
		switch {
		case p.at(token.FUNCTION):
			d.Decl = p.parseFunction(nil, token.NoPos, ast.AccelNone, false)
		case p.atCtx(token.CtxAsync) && p.peek(1).Kind == token.FUNCTION && !p.peek(1).NLBefore():
			asyncPos := p.next().Pos
			d.Decl = p.parseFunction(nil, asyncPos, ast.AccelNone, true)
		case p.at(token.CLASS):
			d.Decl = p.parseClass(nil, token.NoPos)
		case p.at(token.AT):
			// `export default @dec class C {}` — decorators reach the class
			// through Declaration, since J.2 has no DecoratorList of its own.
			decorators := p.parseDecorators()
			if p.at(token.CLASS) {
				d.Decl = p.parseClass(decorators, token.NoPos)
			} else {
				p.errorf(p.cur(), "expected a class after these decorators")
				d.Decl = &ast.BadDecl{From: decorators[0].Pos(), To: p.pos()}
			}
		default:
			// `export default [lookahead ∉ { function, async function, class }]
			// AssignmentExpression ;`
			d.Value = p.parseAssign(allowIn)
			d.Semi = p.expectSemi()
		}
		return d

	case p.at(token.MUL):
		d.StarPos = p.next().Pos
		if p.atCtx(token.CtxAs) {
			p.next()
			d.StarAs = p.parseModuleExportName()
		}
		d.FromPos = p.expectCtx(token.CtxFrom)
		d.Path = p.parseSpecifier()
		d.With = p.tryWithClause()
		d.Semi = p.expectSemi()
		return d

	case p.at(token.LBRACE):
		d.NamedPos, d.Named, d.NamedEnd = p.parseNamedExports()
		if p.atCtx(token.CtxFrom) {
			d.FromPos = p.next().Pos
			d.Path = p.parseSpecifier()
			d.With = p.tryWithClause()
		}
		d.Semi = p.expectSemi()
		return d

	case p.at(token.ASSIGN):
		// `export = EntityName ;` (J.2).
		d.AssignPos = p.next().Pos
		d.Entity = p.parseQualifiedName()
		d.Semi = p.expectSemi()
		return d

	case p.atCtx(token.CtxAs) && p.peek(1).Ctx == token.CtxNamespace:
		// `export as namespace Identifier ;` (J.2).
		p.next()
		d.NamespacePos = p.next().Pos
		d.Namespace = p.parseIdent()
		d.Semi = p.expectSemi()
		return d

	case p.atCtx(token.CtxType):
		// `export type NamedExports ;` and `export type * from ...` (J.2).
		if p.peek(1).Kind == token.LBRACE || p.peek(1).Kind == token.MUL {
			d.TypePos = p.next().Pos
			switch {
			case p.at(token.MUL):
				d.StarPos = p.next().Pos
				if p.atCtx(token.CtxAs) {
					p.next()
					d.StarAs = p.parseModuleExportName()
				}
			default:
				d.NamedPos, d.Named, d.NamedEnd = p.parseNamedExports()
			}
			if p.atCtx(token.CtxFrom) {
				d.FromPos = p.next().Pos
				d.Path = p.parseSpecifier()
				d.With = p.tryWithClause()
			}
			d.Semi = p.expectSemi()
			return d
		}
	}

	// `export VariableStatement` and `export Declaration` (J.2).
	d.Decl = p.parseStmt()
	return d
}

func (p *parser) parseSpecifier() *ast.BasicLit {
	if !p.at(token.STRING) {
		p.errorf(p.cur(), "expected a module specifier string")
		return &ast.BasicLit{Kind: token.STRING, ValuePos: p.pos(), ValueEnd: p.end()}
	}
	t := p.next()
	return &ast.BasicLit{Kind: t.Kind, ValuePos: t.Pos, ValueEnd: t.End, HasEscape: t.HasEscape()}
}

func (p *parser) parseNamedExports() (token.Pos, []*ast.ExportSpec, token.Pos) {
	lb := p.expect(token.LBRACE)
	var specs []*ast.ExportSpec
	for !p.at(token.RBRACE) && !p.atEOF() {
		before := p.i
		s := &ast.ExportSpec{}
		if p.atCtx(token.CtxType) && p.peek(1).Kind != token.COMMA &&
			p.peek(1).Kind != token.RBRACE && p.peek(1).Ctx != token.CtxAs {
			s.TypePos = p.next().Pos
		}
		s.Name = p.parseModuleExportName()
		if p.atCtx(token.CtxAs) {
			s.AsPos = p.next().Pos
			s.As = p.parseModuleExportName()
		}
		specs = append(specs, s)
		if !p.got(token.COMMA) {
			break
		}
		p.advanced(before)
	}
	return lb, specs, p.expect(token.RBRACE)
}