package types

import "github.com/vertex-language/vertex/token"

// Object is a named entity: a variable, constant, function, type name, builtin,
// or package name. Resolution maps each *ast.Ident to one of these.
//
// Type() may return nil during checking. That is not an error state — A.2 ⊢
// "top-level declarations are order-independent: a declaration may refer to any
// other declaration in the same package regardless of textual position", so
// every package-scope name is inserted before any type is resolved. A nil type
// means "not yet resolved"; Typ[Invalid] means "resolution failed".
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

func (o *object) Name() string    { return o.name }
func (o *object) Type() Type      { return o.typ }
func (o *object) Pos() token.Pos  { return o.pos }
func (o *object) Pkg() *Package   { return o.pkg }
func (o *object) SetType(t Type)  { o.typ = t }
func (o *object) objectNode()     {}

// ------------------------------------------------------------------- var

// Var is a variable, parameter, receiver, struct field, or tuple element.
//
// Mode carries A.3.2's `mut` and `var` qualifiers, which are "legal only in a
// parameter or receiver position" and are therefore not part of the Type — see
// Mode in type.go for why the split is load-bearing.
//
// Mutable distinguishes A.5.1's two binding forms: ⊢ "let is immutable and not
// guaranteed to be addressable — it may be a register, an SSA value, or folded
// away entirely. var is mutable and owns a real stack slot for its whole
// lifetime." That is what makes Addressable answerable at all.
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

// Addressable reports whether this binding has a real storage slot.
//
// A.3.2 ⊢ a `mut` argument "must be an addressable var binding or field path",
// and A.4.8 ⊢ addr "requires an addressable operand: a var binding or a field
// path". Both read this, and A.5.1 ⊢ is why a `let` answers false: it "may not
// physically exist anywhere to point at."
func (v *Var) Addressable() bool { return v.mutable || v.field }

// Owning reports whether this parameter takes ownership of its argument.
//
// A.3.2 ⊢ `var T` "denotes the owning convention; whether the callee receives
// the caller's original or a fresh deep copy is decided at the call site by the
// presence or absence of the var marker (A.4.6), never by the declaration." So
// this answers only that the position is owning — A.9.1's question — and never
// which of the two happens.
func (v *Var) Owning() bool { return v.mode == ModeVar }

// ----------------------------------------------------------------- const

// Const is a compile-time constant: an explicit enum discriminant, an
// ArrayLength identifier, or a top-level binding whose initializer A.2 requires
// to be compile-time-evaluable.
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

func (c *Const) Val() Value       { return c.val }
func (c *Const) SetVal(v Value)   { c.val = v }

// ------------------------------------------------------------------ func

// Func is a function, method, initializer, deinitializer, foreign declaration,
// or MethodRequirement.
//
// A.6.4's init and deinit get no distinct object kind, matching ast.FuncDecl:
// A.1.3 makes them ContextualKeywords that are ordinary method names in a
// receiver declaration. Whether a Func is an initializer is a question about
// its name and receiver, answered by whoever asks.
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

// IsMethod reports whether this function declares a receiver (A.6.1).
func (f *Func) IsMethod() bool {
	sig := f.Signature()
	return sig != nil && sig.Recv() != nil
}

// IsInit reports A.6.4's InitializerDeclaration shape: the ContextualKeyword
// `init` in a receiver declaration.
func (f *Func) IsInit() bool { return f.name == token.CtxInit && f.IsMethod() }

// IsDeinit reports A.6.4's DeinitializerDeclaration shape.
func (f *Func) IsDeinit() bool { return f.name == token.CtxDeinit && f.IsMethod() }

// IsEntry reports whether this is the program entry point.
//
// A.6.1 ⊢ "a function named main must take no parameters, return nothing, and
// acts as the program entry point." It also sets [+Await] in its body, which is
// the one place outside an async-marked function where await is licensed.
func (f *Func) IsEntry() bool {
	sig := f.Signature()
	return f.name == "main" && sig != nil && sig.Recv() == nil &&
		sig.Params().Len() == 0 && sig.Results().IsUnit()
}

// -------------------------------------------------------------- typename

// TypeName binds a name to a type or to a constraint: a struct, class, enum,
// abstract alias, type parameter, predeclared name, or ConstraintDeclaration.
//
// One object kind serves both because A.7.2 ⊢ "a single identifier in a
// constraint body parses as both a TypeSet of one term and a ConstraintName; it
// is resolved by what the name denotes." A separate kind would make that
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
	return &TypeName{
		object:     object{name: name, pos: pos, pkg: pkg},
		constraint: c,
	}
}

// IsConstraint reports whether this name denotes a constraint.
//
// It is what rejects A.14's `var c: Ordered`: A.7.2 ⊢ a constraint "is never a
// value type and is legal only in a [...] position." Because Constraint does
// not implement Type, a checker that forgets to call this cannot accidentally
// let one through — it has no Type to return.
func (t *TypeName) IsConstraint() bool { return t.constraint != nil }

func (t *TypeName) Constraint() *Constraint { return t.constraint }

// SetConstraint fills in a constraint after collection. A name minted by
// NewConstraintName remains a constraint name even before this is called, so
// IsConstraint's answer never changes mid-check.
func (t *TypeName) SetConstraint(c *Constraint) {
	if c != nil {
		t.constraint = c
	}
}

// IsAlias reports whether this name was introduced by a transparent
// TypeAliasDeclaration.
//
// A.6.6 ⊢ "an alias to a Type is transparent: it names the same type and
// satisfies a ~T type-set element", so no Named is minted and the name points
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

// Builtin is a ReservedBuiltinName (A.1.4).
//
// It has no meaningful Type: A.4.8 makes each builtin's shape arity- and
// type-argument-specific rather than expressible as a Signature — `new[T]`
// takes a count plus optional named arguments, `sizeof` takes a Type, and
// `resize` has two arities. The checker special-cases each by Id.
//
// A.1.4 ⊢ "a ReservedBuiltinName may not be shadowed, and may not be declared
// as a member, method, or field name." That guarantee is what makes the
// parser's recognition of sizeof/alignof/reinterpret by name sound rather than
// a hack, and it is enforced by Reserved (scope.go) rather than by Insert.
type Builtin struct {
	object
	id BuiltinId
}

func NewBuiltin(name string, id BuiltinId) *Builtin {
	return &Builtin{object: object{name: name, typ: Typ[Invalid]}, id: id}
}

func (b *Builtin) Id() BuiltinId { return b.id }

type BuiltinId int

const (
	_Sizeof BuiltinId = iota
	_Alignof
	_Reinterpret
	_Addr
	_New
	_Delete
	_Resize
	_Copy
	_Zero
	_Unique
	_Shared
	_Weak
	_Upgrade
	_Drop

	// _Transfer exists only so A.1.4's `.transfer()` diagnostic can carry its
	// fix-it rather than degrading into "no such method". ⊢ "the name stays
	// reserved so the diagnostic can carry a fix-it." It is never callable.
	_Transfer
)

// Exported ids, for the analyzer's builtin dispatch.
const (
	SizeofId      = _Sizeof
	AlignofId     = _Alignof
	ReinterpretId = _Reinterpret
	AddrId        = _Addr
	NewId         = _New
	DeleteId      = _Delete
	ResizeId      = _Resize
	CopyId        = _Copy
	ZeroId        = _Zero
	UniqueId      = _Unique
	SharedId      = _Shared
	WeakId        = _Weak
	UpgradeId     = _Upgrade
	DropId        = _Drop
	TransferId    = _Transfer
)

// -------------------------------------------------------------- pkgname

// PkgName is an imported package's qualifier in one file's scope.
//
// A.2.3 ⊢ "the imported package's declared name (its PackageClause) is the
// qualifier under which its symbols are reached; the import path is a locator,
// not a name." So Name() is the imported package's own name and never derived
// from Path — and since A.2.3 also says "there is no aliasing form, no
// dot-import, and no blank import", there is exactly one name to bind and no
// choice about what it is.
type PkgName struct {
	object
	imported *Package
	path     string
}

func NewPkgName(pos token.Pos, pkg *Package, imported *Package, path string) *PkgName {
	return &PkgName{
		object: object{
			name: imported.name,
			typ:  Typ[Invalid],
			pos:  pos,
			pkg:  pkg,
		},
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
// One package is one compilation unit, matching ast.Package. The scope's parent
// is Universe, so A.1.4's predeclared names are found by any lookup that walks
// out far enough — which is what makes them "ordinary identifiers pre-bound in
// an implicit outermost scope" rather than keywords.
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
// only hand out complete packages, since A.2.3 makes the qualifier come from
// the imported package's own clause and a half-checked scope would answer a
// lookup with a nil type.
func (p *Package) Complete() bool     { return p.complete }
func (p *Package) MarkComplete()      { p.complete = true }

func (p *Package) String() string { return "package " + p.path }