// Package analyzer checks a parsed Vertex package against the static rules of
// the grammar annex.
//
// A.0.2 says the parser tracks one context parameter and that "nothing else in
// Vertex is context-sensitive. Every other rejection this document records is a
// static rule checked over an already-parsed tree." This package is the phase
// that checks them. A.14's index of rejected forms is its inventory.
//
// Nothing here writes to the syntax tree. Every result lands in a types.Info
// side table, which is what lets a printer round-trip a checked file and lets
// the checker run twice over one tree without residue.
package analyzer

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// Config carries what checking a package needs beyond the files themselves.
type Config struct {
	Fset     *token.FileSet
	Tag      token.BuildTag
	Reporter diag.Reporter

	// Importer resolves an import path to an already-checked package. It may be
	// nil, in which case any import is reported as an undeclared name — which is
	// the correct behaviour for a single-package check and the seam the importer
	// package will fill.
	Importer Importer
}

// Importer is what the loader will satisfy. A.2.3 ⊢ the imported package's own
// PackageClause supplies the qualifier, so an implementation must have read the
// imported directory's package line before it can answer.
type Importer interface {
	Import(path string) (*types.Package, error)
}

// declInfo is a package-scope declaration awaiting phase 2.
//
// The file scope is recorded alongside the node because import qualifiers are
// file-scoped while the declarations that use them are package-scoped — a
// qualified type in one file's declaration must not resolve through another
// file's imports.
type declInfo struct {
	decl      ast.Decl
	file      *ast.File // the declaring file; A.8.1's BuildClause check reads it
	fileScope *types.Scope
	node      ast.Node // the specific spec/field when decl covers several

	// family is the import family a declare block's handles were minted by
	// (A.4.4). It is decided in phase 1 by familyForBlock, because the block
	// keyword refines what the build tag started — under darwin a `declare
	// module` is flat C while a `declare framework` is an Objective-C object
	// graph — and by phase 2 the enclosing block is no longer in hand.
	//
	// It is the zero value (FamilyUnknown) for every non-foreign declaration,
	// which is correct: nothing else mints a handle.
	family types.Family
}

// Checker holds one package's checking state.
type Checker struct {
	conf  *Config
	pkg   *types.Package
	info  *types.Info
	sizes *types.Sizes
	files []*ast.File

	// scope is the innermost open scope during a body walk.
	scope *types.Scope

	// objMap maps each package-scope object to the declaration that produced
	// it. Phase 2 consults it, which is what makes A.2's order-independence
	// work: "a declaration may refer to any other declaration in the same
	// package regardless of textual position."
	objMap map[types.Object]*declInfo

	// resolving is the phase-2 stack. A named type reached twice while already
	// on it is a cycle, which must be diagnosed rather than recursed into.
	resolving []types.Object

	// npu is set while checking an npu-marked body. It is the one A.0.2 context
	// parameter the parser did not track and this phase must: A.3.5 makes
	// `tensor` grammatical only under [+Npu], and A.14 lists a tensor outside
	// an npu body among the forms that parse and are rejected here.
	npu bool

	errCount int
}

// NewChecker prepares a checker for one package.
func NewChecker(conf *Config, path, name string, info *types.Info) *Checker {
	if info == nil {
		info = &types.Info{}
	}
	pkg := types.NewPackage(path, name)
	return &Checker{
		conf:   conf,
		pkg:    pkg,
		info:   info,
		sizes:  types.SizesFor(conf.Tag.String()),
		objMap: make(map[types.Object]*declInfo),
	}
}

// Files checks a package's files and returns the resulting *types.Package.
//
// The three phases are not an optimization. A.2 makes top-level declarations
// order-independent, so every package-scope name must exist before any
// declaration's type is resolved; and a declaration's type must be known before
// any body that calls it is walked. Collapsing any two would reintroduce a
// forward-declaration requirement the language does not have.
func (c *Checker) Files(files []*ast.File) (*types.Package, error) {
	c.files = files

	c.collectObjects()   // phase 1: names exist, types are nil
	c.resolveDeclTypes() // phase 2: types filled in, cycles caught
	c.checkBodies()      // phase 3: bodies walked, uses recorded

	if c.errCount > 0 {
		return c.pkg, errCount(c.errCount)
	}
	return c.pkg, nil
}

type errCount int

func (n errCount) Error() string {
	if n == 1 {
		return "1 error"
	}
	return "analysis found errors"
}

// Package returns the package under construction, valid even after errors —
// resolution produces a partial result deliberately, so editor tooling can use
// what did resolve.
func (c *Checker) Package() *types.Package { return c.pkg }

// ------------------------------------------------------------------ scopes

func (c *Checker) openScope(n ast.Node, comment string) {
	s := types.NewScope(c.scope, comment)
	if n != nil {
		s.SetExtent(n.Pos(), n.End())
		c.info.RecordScope(n, s)
	}
	c.scope = s
}

func (c *Checker) closeScope() { c.scope = c.scope.Parent() }

// declare inserts obj into scope, reporting a redeclaration or a shadowed
// builtin.
//
// A.1.4 ⊢ "a ReservedBuiltinName may not be shadowed, and may not be declared
// as a member, method, or field name." That is checked here rather than in
// Scope.Insert, because Insert cannot express the asymmetry between the
// predeclared type names (shadowable) and the builtins (not).
func (c *Checker) declare(scope *types.Scope, id *ast.Ident, obj types.Object) {
	if id == nil {
		return
	}
	if id.IsBlank() {
		// A.1.2 ⊢ `_` never introduces a usable binding. It is recorded as a
		// definition so tooling can still point at it, but nothing is inserted.
		c.info.RecordDef(id, nil)
		return
	}
	if types.Reserved(id.Name) {
		c.errorAt(id.Pos(), id.End(), diag.ShadowedBuiltin, id.Name)
		return
	}
	if alt := scope.Insert(obj); alt != nil {
		d := diag.New(diag.DuplicateDeclaration, id.Pos(), id.End(), id.Name)
		if alt.Pos().IsValid() {
			d.WithNote(alt.Pos(), alt.Pos(), "previous declaration of %s", id.Name)
		}
		c.report(d)
		return
	}
	c.info.RecordDef(id, obj)
}

// lookup resolves a name from the innermost open scope outward, recording the
// use. A.1.4's predeclared names are found because Universe is the root of
// every chain.
func (c *Checker) lookup(id *ast.Ident) types.Object {
	if id.IsBlank() {
		c.errorAt(id.Pos(), id.End(), diag.BlankNotUsable)
		return nil
	}
	scope := c.scope
	if scope == nil {
		scope = c.pkg.Scope()
	}
	_, obj := scope.LookupParent(id.Name)
	if obj == nil {
		c.errorAt(id.Pos(), id.End(), diag.UndeclaredName, id.Name)
		return nil
	}
	c.info.RecordUse(id, obj)
	return obj
}