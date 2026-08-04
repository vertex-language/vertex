package types

import "github.com/vertex-language/vertex/token"

// Object is a named entity: a variable, constant, function, type name, builtin,
// or package qualifier. Resolution maps each *ast.Ident to one of these.
//
// Type may return nil during checking. That is not an error state — §1.1 ⊢
// "top-level declarations are order-independent: a declaration may reference
// any other in its package regardless of position or file" — so every
// package-scope name is inserted before any type is resolved. A nil type means
// "not yet resolved"; Typ[Invalid] means "resolution failed", and predicates
// treat the second as already-diagnosed.
type Object interface {
	Name() string
	Type() Type
	Pos() token.Pos
	Pkg() *Package

	// SetType fills in a type after collection. Exported because the analyzer
	// is a separate package and is the only thing that ever calls it.
	SetType(Type)

	objectNode()
}

type object struct {
	name string
	typ  Type
	pos  token.Pos
	pkg  *Package
}

func (o *object) Name() string   { return o.name }
func (o *object) Type() Type     { return o.typ }
func (o *object) Pos() token.Pos { return o.pos }
func (o *object) Pkg() *Package  { return o.pkg }
func (o *object) SetType(t Type) { o.typ = t }
func (o *object) objectNode()    {}

// ------------------------------------------------------------------- var

// Var is a variable, parameter, receiver, struct field, or tuple element.
//
// Mode carries `mut` and `var`, which §3.2 makes "a parameter or receiver only"
// and therefore not part of the Type — see type.go for why the split is
// load-bearing.
//
// Mutable distinguishes §6.1's two binding forms: ⊢ "`let` requires an
// initializer and fixes the binding; `var` may be rebound and is required by
// anything taking exclusive access or transferring."
type Var struct {
	object
	mode    Mode
	mutable bool
	field   bool
}

func NewVar(pos token.Pos, pkg *Package, name string, typ Type) *Var {
	return &Var{object: object{name, typ, pos, pkg}}
}

func NewParam(pos token.Pos, pkg *Package, name string, typ Type, mode Mode) *Var {
	return &Var{object: object{name, typ, pos, pkg}, mode: mode}
}

func NewField(pos token.Pos, pkg *Package, name string, typ Type) *Var {
	return &Var{object: object{name, typ, pos, pkg}, field: true}
}

func (v *Var) Mode() Mode        { return v.mode }
func (v *Var) Mutable() bool     { return v.mutable }
func (v *Var) IsField() bool     { return v.field }
func (v *Var) SetMutable(b bool) { v.mutable = b }
func (v *Var) SetMode(m Mode)    { v.mode = m }

// Assignable reports the first two entries of §6.2's list: ⊢ an AssignTarget is
// assignable when it is "a `var` binding" or "a field of an assignable value".
// The remaining entries are questions about the base expression's type rather
// than about one object, and live in predicates.go as InteriorAssignable.
func (v *Var) Assignable() bool { return v.mutable || v.field }

// Owning reports whether this parameter takes ownership of its argument.
//
// §8.1 ⊢ "the convention lives in the signature; only the owning one has a
// choice at the call." So this answers that the position is owning, and never
// which of copy or move happens — §8.2 makes that the marker's question at the
// call site.
func (v *Var) Owning() bool { return v.mode == ModeVar }

// ----------------------------------------------------------------- const

// Const is a compile-time constant: an explicit enum discriminant, a constant
// used as an ArrayLength, or a top-level binding, whose initializer §6.1
// requires to be one.
type Const struct {
	object
	val Value
}

func NewConst(pos token.Pos, pkg *Package, name string, typ Type, v Value) *Const {
	if v == nil {
		v = Unknown
	}
	return &Const{object{name, typ, pos, pkg}, v}
}

func (c *Const) Val() Value     { return c.val }
func (c *Const) SetVal(v Value) { c.val = v }

// ------------------------------------------------------------------ func

// Func is a function, method, initializer, deinitializer, foreign declaration,
// or MethodRequirement.
//
// `init` and `deinit` get no distinct object kind, matching ast.FuncDecl: they
// are contextual keywords that are ordinary method names in a receiver
// declaration. §7.2 ⊢ they are "ordinary method names recognized by spelling",
// so whether a Func is one is a question about its name and receiver, answered
// by whoever asks.
type Func struct {
	object
}

func NewFunc(pos token.Pos, pkg *Package, name string, sig *Signature) *Func {
	f := &Func{object{name: name, pos: pos, pkg: pkg}}
	if sig != nil {
		f.typ = sig
	}
	return f
}

// Signature returns the function's type, or nil if it is not yet resolved or
// resolution failed.
func (f *Func) Signature() *Signature {
	sig, _ := f.typ.(*Signature)
	return sig
}

func (f *Func) IsMethod() bool {
	sig := f.Signature()
	return sig != nil && sig.Recv() != nil
}

func (f *Func) IsInit() bool   { return f.name == token.CtxInit && f.IsMethod() }
func (f *Func) IsDeinit() bool { return f.name == token.CtxDeinit && f.IsMethod() }

// IsEntry reports whether this is the program's entry point.
//
// §1.4 ⊢ "a program has exactly one package named `main` declaring exactly one
// `func main()` — no parameters, no result, no marker." All four conditions are
// checked here; that the package is named `main`, and that there is exactly
// one, are the loader's and the checker's questions rather than this object's.
//
// §1.4 also makes it "the one non-`async` function in which `await` is legal",
// which is why the marker check matters beyond tidiness.
func (f *Func) IsEntry() bool {
	sig := f.Signature()
	return f.name == "main" && sig != nil &&
		sig.Recv() == nil &&
		sig.Params().Len() == 0 &&
		sig.IsVoid() &&
		sig.Marker() == MarkerNone
}

// IsTest reports §7.4's test-function shape: ⊢ "a `test` function takes no
// parameters, carries an `Expected` result or none, and exists only in a
// `build test` file." The build tag is the file's and is checked by the
// analyzer against token.LicensesTest.
func (f *Func) IsTest() bool {
	sig := f.Signature()
	return sig != nil && sig.Marker() == MarkerTest &&
		sig.Params().Len() == 0 &&
		(sig.Expected() != nil || sig.Results().Len() == 0)
}

// -------------------------------------------------------------- typename

// TypeName binds a name to a type or to a constraint: a struct, class, enum,
// abstract alias, transparent alias, type parameter, predeclared name, or
// ConstraintDecl.
//
// One object kind serves both because grammar.md ⊢ "a ConstraintElem that is a
// single identifier parses as both a one-term TypeSet and a constraint name;
// resolution is by what the name denotes." A separate kind would make that
// resolution a two-map lookup instead of one field test.
type TypeName struct {
	object
	constraint *Constraint // non-nil when this name denotes a constraint
}

func NewTypeName(pos token.Pos, pkg *Package, name string, typ Type) *TypeName {
	return &TypeName{object: object{name, typ, pos, pkg}}
}

// NewConstraintName mints a name that denotes a constraint rather than a type.
// c may be nil during collection and filled in by SetConstraint, matching the
// nil-type convention Object documents.
func NewConstraintName(pos token.Pos, pkg *Package, name string, c *Constraint) *TypeName {
	return &TypeName{object: object{name: name, pos: pos, pkg: pkg}, constraint: c}
}

// IsConstraint reports whether this name denotes a constraint.
//
// It is what rejects `var c: Ordered`. Because Constraint does not implement
// Type, a checker that forgets to call this cannot let one through anyway — it
// has no Type to return.
func (t *TypeName) IsConstraint() bool { return t.constraint != nil }

func (t *TypeName) Constraint() *Constraint { return t.constraint }

// SetConstraint fills in a constraint after collection. A name minted by
// NewConstraintName is already a constraint name before this is called, so
// IsConstraint's answer never changes mid-check.
func (t *TypeName) SetConstraint(c *Constraint) {
	if c != nil {
		t.constraint = c
	}
}

// IsAlias reports whether this name was introduced by a transparent
// TypeAliasDecl.
//
// §3.1 ⊢ such an alias "introduces a second name for one type, interchangeable
// with the first in both directions", so no Named is minted and the name points
// straight at its target. An alias to `abstract` is nominal and does mint one,
// which is why that case answers false.
func (t *TypeName) IsAlias() bool {
	if t.typ == nil {
		return false
	}
	n, ok := t.typ.(*Named)
	return !ok || n.Obj() != t
}

// --------------------------------------------------------------- builtin

// Builtin is a reserved builtin name.
//
// It has no meaningful Type: each builtin's shape is arity- and
// type-argument-specific rather than expressible as a Signature — `sizeof`
// takes a Type, `new` takes a count plus named arguments, `resize` has two
// arities — so the checker special-cases each by Id.
//
// §2.3 ⊢ "reserved builtin names may not be shadowed — not as a local,
// parameter, type parameter, field, method, or top-level declaration, and not
// as a parameter label." That guarantee is what makes the parser's recognition
// of sizeof/alignof/reinterpret by name sound, and it is enforced by Reserved
// (scope.go) rather than by Insert.
type Builtin struct {
	object
	id BuiltinId
}

func NewBuiltin(name string, id BuiltinId) *Builtin {
	return &Builtin{object: object{name: name, typ: Typ[Invalid]}, id: id}
}

func (b *Builtin) Id() BuiltinId { return b.id }

// BuiltinId enumerates §2.3's reserved builtin names, and nothing else.
// `unique`, `shared`, and `weak` are absent on purpose: grammar.md ⊢ they "are
// keywords, not reserved builtin names; they get the HeapConstructor
// production", so they never resolve to an object at all.
type BuiltinId int

const (
	NewId BuiltinId = iota
	DeleteId
	ResizeId
	CopyId
	ZeroId
	AddrId
	SizeofId
	AlignofId
	ReinterpretId
	UpgradeId
	DropId
	PanicId
	BlendId
	MinId
	MaxId
	ClampId

	// TransferId is bound to nothing. §2.3 ⊢ "`transfer` is reserved and bound
	// to nothing, so `x.transfer()` and `transfer(x)` diagnose against
	// ownership §3.3 rather than as unknown names" — the object exists only so
	// the diagnostic can carry diag's fix-it to the `var` prefix. Calling it is
	// always an error.
	TransferId
)

// -------------------------------------------------------------- pkgname

// PkgName is an imported package's qualifier in one file's scope.
//
// §1.3 ⊢ "imports are file-scoped: an `import` in one file of a package does
// not bring the qualifier into the others. The qualifier is the imported
// package's own PackageName, never the path." So Name is the imported
// package's own name and is never derived from Path — and since there is no
// aliasing form, no dot-import, and no blank import, there is exactly one name
// to bind and no choice about what it is.
//
// §1.3 also makes two imports whose packages declare the same name a compile
// error, "there is no aliasing form to resolve the clash with" — which is a
// duplicate Insert into the file scope and needs nothing here.
type PkgName struct {
	object
	imported *Package
	path     string
}

func NewPkgName(pos token.Pos, pkg *Package, imported *Package, path string) *PkgName {
	return &PkgName{
		object:   object{name: imported.name, typ: Typ[Invalid], pos: pos, pkg: pkg},
		imported: imported,
		path:     path,
	}
}

func (p *PkgName) Imported() *Package { return p.imported }
func (p *PkgName) Path() string       { return p.path }

// -------------------------------------------------------------- package

// Package is a checked package: its declared name, its resolved path, and its
// package scope.
//
// §1.1 ⊢ "a package is the set of selected files carrying the same
// PackageName", which matches ast.Package. The scope's parent is Universe, so
// §2.3's predeclared names are found by any lookup that walks out far enough —
// which is what makes them ordinary identifiers pre-bound in an implicit scope
// rather than keywords.
//
// There is no visibility modifier in Vertex (§1.1), so a package scope has no
// exported subset and every top-level declaration is reachable by any importer.
type Package struct {
	name     string
	path     string
	scope    *Scope
	imports  []*Package
	complete bool
}

func NewPackage(path, name string) *Package {
	p := &Package{name: name, path: path}
	p.scope = NewScope(Universe, "package "+path)
	return p
}

func (p *Package) Name() string        { return p.name }
func (p *Package) Path() string        { return p.path }
func (p *Package) Scope() *Scope       { return p.scope }
func (p *Package) Imports() []*Package { return p.imports }

func (p *Package) SetImports(l []*Package) { p.imports = l }

// Complete reports whether this package finished checking. An importer must
// only hand out complete packages: §1.3 makes the qualifier come from the
// imported package's own clause, and a half-checked scope would answer a lookup
// with a nil type.
func (p *Package) Complete() bool { return p.complete }
func (p *Package) MarkComplete()  { p.complete = true }

func (p *Package) String() string { return "package " + p.path }