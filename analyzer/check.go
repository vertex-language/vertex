// Package analyzer checks a parsed Vertex package against the static rules and
// produces a *types.Package plus a types.Info side table.
//
// grammar.md resolves exactly one thing by context, the literal-in-header
// ambiguity, and the parser owns that. Everywhere else grammar.md says "static
// rule", which means the form derives, parses, and is checked afterwards over
// an already-parsed tree. This package is that check. The forms grammar.md
// describes as parsing-and-then-rejected are its inventory, and diag's ranges
// outside 1xxx and 2xxx are its diagnostic surface.
//
// Nothing here writes to the syntax tree. Every result lands in a types.Info
// side table, which is what lets a printer round-trip a checked file and lets
// the checker run twice over one tree without residue. The checker also keeps
// its own defs/uses/file-scope maps rather than reading back out of Info, so a
// caller that allocated only some of Info's tables still gets a full check.
//
// Citation convention, matching types: a bare § is semantics.md, CamelCase
// names are grammar.md productions. Where neither document fixes something —
// the variant-tag set, the shape of a constant expression richer than this
// package folds — the comment says so rather than inventing a rule.
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
	// nil, in which case every import is reported as an undeclared name — the
	// correct behaviour for a single-package check, and the seam a loader fills.
	Importer Importer
}

// Importer is what the loader satisfies. §1.3 ⊢ the qualifier comes from the
// imported package's own package clause, so an implementation must have read
// the imported directory's package line before it can answer.
type Importer interface {
	Import(path string) (*types.Package, error)
}

// declInfo is a package-scope declaration awaiting phase 2.
//
// The file scope travels with the node because import qualifiers are
// file-scoped (§1.3) while the declarations using them are package-scoped: a
// qualified type in one file must not resolve through another file's imports.
type declInfo struct {
	decl      ast.Decl
	file      *ast.File
	fileScope *types.Scope

	// node is the specific member or binding when decl covers several, and
	// index is its position among them — a VarDecl's bindings line up with its
	// values positionally, and phase 2 needs to know which.
	node  ast.Node
	index int

	// family is the import family a declare block's handles were minted by. It
	// is decided in phase 1, because the block keyword refines what the build
	// tag started — under darwin a `declare module` is flat C while a
	// `declare framework` is an object graph — and by phase 2 the enclosing
	// block is no longer in hand. FamilyUnknown for everything else, which is
	// correct: nothing else mints a handle.
	family types.Family
}

// bodyCtx is the parse-adjacent context a function body establishes.
//
// grammar.md leaves npu and async to this phase: `tensor` is legal only inside
// an npu-marked function, `await` only inside an async-marked one, and a
// FuncLit "begins with all enclosing parse context cleared", so both reset at
// every function boundary rather than propagating into a closure.
type bodyCtx struct {
	npu   bool
	async bool
}

// Checker holds one package's checking state.
type Checker struct {
	conf  *Config
	pkg   *types.Package
	info  *types.Info
	sizes *types.Sizes
	files []*ast.File

	// scope is the innermost open scope.
	scope *types.Scope

	// objMap maps each package-scope object to the declaration that produced
	// it. Phase 2 consults it, which is what makes order-independence work:
	// §1.1 ⊢ a declaration may reference any other in its package regardless of
	// position or file.
	objMap map[types.Object]*declInfo

	// defs, uses, fileScopes, and funcObj are the checker's own copies of what
	// it also records into Info. Info's maps are optional by design — a nil map
	// means "do not record this" — so the checker must not depend on them.
	defs       map[*ast.Ident]types.Object
	uses       map[*ast.Ident]types.Object
	fileScopes map[*ast.File]*types.Scope
	funcObj    map[*ast.FuncDecl]*types.Func

	// resolving is the phase-2 stack. An object reached twice while still on it
	// is a cycle, to be diagnosed rather than recursed into.
	resolving []types.Object

	ctx bodyCtx

	// inSignature and inTensorElem gate the two tensor-element rules, which
	// differ by position rather than by shape. See typeFromName.
	inSignature  bool
	inTensorElem bool

	errCount int
}

// NewChecker prepares a checker for one package.
func NewChecker(conf *Config, path, name string, info *types.Info) *Checker {
	if info == nil {
		info = types.NewInfo()
	}
	return &Checker{
		conf:       conf,
		pkg:        types.NewPackage(path, name),
		info:       info,
		sizes:      types.SizesFor(conf.Tag),
		objMap:     make(map[types.Object]*declInfo),
		defs:       make(map[*ast.Ident]types.Object),
		uses:       make(map[*ast.Ident]types.Object),
		fileScopes: make(map[*ast.File]*types.Scope),
		funcObj:    make(map[*ast.FuncDecl]*types.Func),
	}
}

// Files checks a package's files and returns the resulting *types.Package.
//
// The three phases are not an optimization. §1.1 makes top-level declarations
// order-independent, so every package-scope name must exist before any
// declaration's type is resolved, and every declaration's type must be known
// before any body that calls it is walked. Collapsing any two would reintroduce
// a forward-declaration requirement the language does not have.
func (c *Checker) Files(files []*ast.File) (*types.Package, error) {
	c.files = files

	c.collectObjects()   // phase 1: names exist, types are nil
	c.resolveDeclTypes() // phase 2: types filled in, cycles caught
	c.checkBodies()      // phase 3: bodies walked, uses recorded

	if c.errCount > 0 {
		return c.pkg, errCount(c.errCount)
	}
	// §1.3 ⊢ an importer must only hand out complete packages: a half-checked
	// scope would answer a lookup with a nil type.
	c.pkg.MarkComplete()
	return c.pkg, nil
}

type errCount int

func (n errCount) Error() string {
	if n == 1 {
		return "1 error"
	}
	return "analysis found errors"
}

// Package returns the package under construction, valid even after errors.
// Resolution produces a partial result deliberately, so editor tooling can use
// whatever did resolve.
func (c *Checker) Package() *types.Package { return c.pkg }

// ------------------------------------------------------------------ scopes

func (c *Checker) openScope(n ast.Node, comment string) {
	c.pushScope(types.NewScope(c.scope, comment), n)
}

// openFuncScope opens the scope a `defer` registers against. §6.6 ⊢ a deferred
// call runs when the enclosing *function* returns, so Info.Defers groups by
// this scope and not by the innermost block.
func (c *Checker) openFuncScope(n ast.Node, comment string) {
	c.pushScope(types.NewFuncScope(c.scope, comment), n)
}

func (c *Checker) pushScope(s *types.Scope, n ast.Node) {
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
// §2.3 ⊢ a reserved builtin name may not be shadowed — not as a local,
// parameter, type parameter, field, method, or top-level declaration. That is
// checked here rather than in Scope.Insert, because Insert cannot express the
// asymmetry between the predeclared type names (shadowable) and the builtins
// (not).
func (c *Checker) declare(scope *types.Scope, id *ast.Ident, obj types.Object) {
	if id == nil || scope == nil {
		return
	}
	if id.IsBlank() {
		// §2.4 ⊢ `_` introduces no binding and may be repeated freely. It is
		// recorded as a definition anyway so tooling can point at it.
		c.recordDef(id, nil)
		return
	}
	if types.Reserved(id.Name) {
		c.errorExpr(id, diag.ShadowedBuiltin, id.Name)
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
	c.recordDef(id, obj)
}

// lookup resolves a name from the innermost open scope outward, recording the
// use. §2.3's predeclared names are found because Universe is the root of every
// chain — which is what makes them ordinary identifiers in an implicit scope
// rather than keywords.
func (c *Checker) lookup(id *ast.Ident) types.Object {
	if id == nil {
		return nil
	}
	if id.IsBlank() {
		// §2.4 ⊢ reading `_` is always an error: it names nothing.
		c.errorExpr(id, diag.BlankNotUsable)
		return nil
	}
	scope := c.scope
	if scope == nil {
		scope = c.pkg.Scope()
	}
	_, obj := scope.LookupParent(id.Name, id.Pos())
	if obj == nil {
		c.errorExpr(id, diag.UndeclaredName, id.Name)
		return nil
	}
	c.recordUse(id, obj)
	return obj
}

func (c *Checker) recordDef(id *ast.Ident, obj types.Object) {
	c.defs[id] = obj
	c.info.RecordDef(id, obj)
}

func (c *Checker) recordUse(id *ast.Ident, obj types.Object) {
	if obj == nil {
		return
	}
	c.uses[id] = obj
	c.info.RecordUse(id, obj)
}

func (c *Checker) objectOf(id *ast.Ident) types.Object {
	if obj, ok := c.defs[id]; ok && obj != nil {
		return obj
	}
	return c.uses[id]
}