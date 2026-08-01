package types

// Named is a defined type: a struct, a class, an enum, or an abstract alias.
//
// A transparent TypeAliasDeclaration produces no Named at all. A.6.6 ⊢ "an
// alias to a Type is transparent: it names the same type and satisfies a ~T
// type-set element", so the resolver substitutes the target and the alias
// leaves no trace here. Only `= abstract` mints a Named, because A.6.6 makes
// that one nominal.
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

func (n *Named) AddMethod(f *Func) { n.methods = append(n.methods, f) }

func (n *Named) SetTypeParams(tp []*TypeParam) { n.typeParams = tp }
func (n *Named) SetTypeArgs(ta []Type)         { n.typeArgs = ta }

// SetUnderlying fills in what this Named was declared as, after the object
// itself exists.
//
// The two-step construction is not a convenience. A.2's order-independence lets
// a field name its own enclosing type — `struct Node { next: typed_ptr Node }` —
// so the resolver must bind the Named to its TypeName before walking any field,
// or the field's lookup trips the cycle guard on an object whose type is still
// nil. NewNamed(obj, nil) opens that window and this closes it.
//
// Underlying() answers nil inside the window. That is the honest answer — the
// declaration has no shape yet — and every structural predicate reaches it
// through predicates.Underlying, whose type switches simply do not match, so a
// premature read degrades to "not a struct, not an enum" rather than panicking.
func (n *Named) SetUnderlying(t Type) { n.underlying = t }

// LookupMethod finds a declared method by name. Vertex has no inheritance and
// no embedding (A.6.3), so there is no promotion to walk — the set is flat.
func (n *Named) LookupMethod(name string) *Func {
	for _, m := range n.methods {
		if m.name == name {
			return m
		}
	}
	return nil
}

// ------------------------------------------------------------------ struct

// Field is one entry of a struct or class body (A.6.2).
type Field struct {
	Name       string
	Type       Type
	HasDefault bool // A.6.2: evaluated at construction for any omitted field
}

// Struct is the underlying of both a StructDeclaration and a ClassDeclaration.
//
// One type, because A.6.3 ⊢ "a class is byte-for-byte identical in layout to a
// struct and differs only in its member and method model." Sizeof does not
// branch on it; the construction syntax (A.4.7's brace literal vs. calling
// init) and `===` identity comparison do, and both read Class().
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

// Variant is one entry of an enum body (A.6.5).
//
// Payload is nil for a unit variant. Value is the discriminant, which A.6.5
// says continues from the previous variant when unassigned — so it is resolved
// here rather than left optional.
type Variant struct {
	Name    string
	Payload []Type // nil for a unit variant
	Value   int64
}

func (v *Variant) IsUnit() bool { return len(v.Payload) == 0 }

// Enum is an EnumDeclaration (A.6.5).
//
// A.6.5 ⊢ a unit-only enum *is* its discriminant integer, so `Status.Active as
// int32` is a reinterpretation and not a conversion. A payload enum is a tagged
// union sized to the largest variant plus the tag. UnitOnly is what Sizeof and
// the cast rules branch on.
type Enum struct {
	variants []*Variant
	discrim  *Basic // the DiscriminantType; defaults to Int32
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

// UnitOnly reports whether every variant is a unit variant (A.6.5).
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

// Family is the import family an abstract handle was minted by (A.4.4).
//
// It exists because the cast rules differ by family and not by declaration:
// A.4.4 ⊢ `abstract` → `typed_ptr T` is legal only for a memory-flat family
// (C, WASM) and is a compile error for an object-graph family (Objective-C,
// JS), "whose handles have no byte representation to point at".
type Family uint8

const (
	FamilyUnknown Family = iota
	FamilyMemoryFlat
	FamilyObjectGraph
)

// Abstract is the bare `abstract` of A.3.3, legal only as an alias target.
//
// Each such alias is a distinct nominal type: A.6.6 ⊢ "two abstract aliases
// never unify, however identical their provenance." That is why Abstract holds
// its minting object — identity is the object, not the shape.
type Abstract struct {
	obj    *TypeName
	family Family
}

func NewAbstract(obj *TypeName, f Family) *Abstract { return &Abstract{obj, f} }

func (a *Abstract) Obj() *TypeName   { return a.obj }
func (a *Abstract) Family() Family   { return a.family }
func (a *Abstract) Underlying() Type { return a }

// ------------------------------------------------------------ type params

// TypeParam is one entry of a TypeParameterList (A.7.1).
//
// A.7.1 ⊢ "a bare name is constraint `any`", so Constraint is never nil after
// resolution — the parser's nil means "not written", and distribution across a
// group happens here rather than in the tree.
type TypeParam struct {
	name       string
	index      int
	constraint *Constraint
}

func NewTypeParam(name string, index int, c *Constraint) *TypeParam {
	return &TypeParam{name, index, c}
}

func (t *TypeParam) Name() string            { return t.name }
func (t *TypeParam) Index() int              { return t.index }
func (t *TypeParam) Constraint() *Constraint { return t.constraint }
func (t *TypeParam) Underlying() Type        { return t }