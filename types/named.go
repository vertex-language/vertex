package types

// Named is a defined type: a struct, a class, an enum, or an `abstract` alias.
//
// §3.1 ⊢ declared types "are nominal: two declarations are distinct types even
// with identical fields". A transparent TypeAliasDecl produces no Named at all
// — ⊢ it "introduces a second name for one type, interchangeable with the first
// in both directions and at every depth of composition" — so the resolver
// substitutes the target and the alias leaves no trace here. Only
// `= abstract` mints a Named, because ⊢ "each `abstract` alias is distinct from
// every other".
type Named struct {
	obj        *TypeName // the declaration this was minted from
	underlying Type      // *Struct, *Enum, or *Abstract
	methods    []*Func
	typeParams []*TypeParam // nil if not generic
	typeArgs   []Type       // non-nil on an instantiation
}

func NewNamed(obj *TypeName, underlying Type) *Named {
	return &Named{obj: obj, underlying: underlying}
}

func (n *Named) Obj() *TypeName           { return n.obj }
func (n *Named) Underlying() Type         { return n.underlying }
func (n *Named) NumMethods() int          { return len(n.methods) }
func (n *Named) Method(i int) *Func       { return n.methods[i] }
func (n *Named) TypeParams() []*TypeParam { return n.typeParams }
func (n *Named) TypeArgs() []Type         { return n.typeArgs }

// AddMethod attaches a declared method.
//
// §2.1 ⊢ "method names are not in any of these scopes — they are reached only
// through a receiver, so a method `read` and a function `read` in one package
// do not collide." That is why the method set lives on the type rather than in
// the package scope, and why AddMethod is not an Insert.
func (n *Named) AddMethod(f *Func) { n.methods = append(n.methods, f) }

func (n *Named) SetTypeParams(tp []*TypeParam) { n.typeParams = tp }
func (n *Named) SetTypeArgs(ta []Type)         { n.typeArgs = ta }

// SetUnderlying fills in what this Named was declared as, after the object
// itself exists.
//
// The two-step construction is not a convenience. §1.1's order-independence
// lets a field name its own enclosing type — `struct Node { next: typed_ptr
// Node }`, which §3.4 explicitly endorses as the way to break a cycle — so the
// resolver must bind the Named to its TypeName before walking any field, or the
// field's lookup trips on an object whose type is still nil.
//
// Underlying answers nil inside that window. That is the honest answer, and
// every structural predicate reaches it through predicates.Underlying, whose
// type switches simply do not match — so a premature read degrades to "not a
// struct, not an enum" rather than panicking.
func (n *Named) SetUnderlying(t Type) { n.underlying = t }

// LookupMethod finds a declared method by name. There is no inheritance and no
// embedding in Vertex, so there is no promotion to walk: the set is flat.
func (n *Named) LookupMethod(name string) *Func {
	for _, m := range n.methods {
		if m.name == name {
			return m
		}
	}
	return nil
}

// ------------------------------------------------------------------ struct

// Field is one FieldDecl.
//
// HasDefault records that the declaration wrote `= Expression`. §7.2 ⊢ "field
// defaults are evaluated at each construction for each omitted field, and may
// not reference other fields or the value under construction" — and §3.3 ⊢ a
// zero value applies none of them, since "field defaults are not applied; they
// belong to construction". So the default is a construction-site question and
// only its presence is a property of the type.
type Field struct {
	Name       string
	Type       Type
	HasDefault bool
}

// Struct is the underlying of both a StructDecl and a ClassDecl.
//
// One type, because §7.2 and grammar.md agree that "a class is byte-for-byte
// identical in layout to a struct and differs only in its member and method
// model". Sizeof does not branch on Class(); two things do, and both are about
// meaning rather than layout: §7.2's construction rule ("a struct is built by a
// composite literal; a class is built by calling an initializer"), and §3.5's
// `===`, which "is legal on classes only".
type Struct struct {
	fields []*Field
	class  bool
}

func NewStruct(fields []*Field, class bool) *Struct { return &Struct{fields, class} }

func (s *Struct) NumFields() int     { return len(s.fields) }
func (s *Struct) Field(i int) *Field { return s.fields[i] }
func (s *Struct) Class() bool        { return s.class }
func (s *Struct) Underlying() Type   { return s }

func (s *Struct) LookupField(name string) (*Field, int) {
	for i, f := range s.fields {
		if f.Name == name {
			return f, i
		}
	}
	return nil, -1
}

// -------------------------------------------------------------------- enum

// Variant is one entry of an enum body. Payload is nil for a unit variant.
//
// Value is the resolved discriminant. §5.3 requires an explicit one to be a
// constant expression; the sources do not say what an omitted one is, and this
// implementation continues from the previous variant, since that is the only
// reading under which a partially annotated list has an answer at all. Treat it
// as this compiler's choice, not as a stated rule.
type Variant struct {
	Name    string
	Payload []Type // nil for a unit variant
	Value   int64
}

func (v *Variant) IsUnit() bool { return len(v.Payload) == 0 }

// Enum is an EnumDecl.
//
// §3.3 ⊢ the zero value is "the first declared variant, with any payload
// zeroed", so Variant(0) is load-bearing and the variant list keeps declaration
// order.
//
// discrim is the DiscriminantType. §4.2 ⊢ the one conversion is "enum → its
// discriminant type... one-way only; there is no `n as Status`" — so this field
// is what ConvertibleTo compares against, and the direction is not symmetric.
// The sources fix no default when the clause is omitted; int32 is this
// implementation's choice.
type Enum struct {
	variants []*Variant
	discrim  *Basic
}

func NewEnum(variants []*Variant, discrim *Basic) *Enum {
	if discrim == nil {
		discrim = Typ[Int32]
	}
	return &Enum{variants, discrim}
}

func (e *Enum) NumVariants() int       { return len(e.variants) }
func (e *Enum) Variant(i int) *Variant { return e.variants[i] }
func (e *Enum) Discriminant() *Basic   { return e.discrim }
func (e *Enum) Underlying() Type       { return e }

// UnitOnly reports whether every variant is a unit variant. It is a shape
// query, used by enumSize to pick a layout — the sources fix no enum layout, so
// it licenses nothing about conversion, which §4.2 states independently of it.
func (e *Enum) UnitOnly() bool {
	for _, v := range e.variants {
		if !v.IsUnit() {
			return false
		}
	}
	return true
}

func (e *Enum) LookupVariant(name string) *Variant {
	for _, v := range e.variants {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// ---------------------------------------------------------------- abstract

// Family is the import family an abstract handle was minted by.
//
// It exists because §4.2's cast rule differs by family and not by declaration:
// ⊢ "`abstract` → `typed_ptr T` | only where linkage is memory-flat; never the
// reverse." A family whose handles have no byte representation has nothing to
// point at, so the same written cast is legal in one declare block and not in
// another.
type Family uint8

const (
	FamilyUnknown Family = iota
	FamilyMemoryFlat
	FamilyObjectGraph
)

// Abstract is the bare `abstract` of an AliasTarget, legal only there.
//
// Each such alias is a distinct nominal type (§3.1), which is why Abstract
// holds its minting object: identity is the object, not the shape.
//
// §3.3 ⊢ its zero value is "the zeroed representation — legal only on an error
// path, paired with a non-empty string", which is a rule about where the zero
// may appear and not about this shape.
type Abstract struct {
	obj    *TypeName
	family Family
}

func NewAbstract(obj *TypeName, f Family) *Abstract { return &Abstract{obj, f} }

func (a *Abstract) Obj() *TypeName   { return a.obj }
func (a *Abstract) Family() Family   { return a.family }
func (a *Abstract) Underlying() Type { return a }

// ------------------------------------------------------------ type params

// TypeParam is one TypeParamDecl.
//
// grammar.md ⊢ "a bare TypeParamName is constrained by `any`", so Constraint is
// never nil after resolution: the parser's nil means "not written". The same
// applies to grammar.md's group distribution — `[A, B: Number]` constrains both
// — which is "performed over an already-parsed list, not by the grammar", and
// therefore lands here rather than in the tree.
type TypeParam struct {
	name       string
	index      int
	constraint *Constraint
}

func NewTypeParam(name string, index int, c *Constraint) *TypeParam {
	if c == nil {
		c = Any
	}
	return &TypeParam{name, index, c}
}

func (t *TypeParam) Name() string            { return t.name }
func (t *TypeParam) Index() int              { return t.index }
func (t *TypeParam) Constraint() *Constraint { return t.constraint }
func (t *TypeParam) Underlying() Type        { return t }