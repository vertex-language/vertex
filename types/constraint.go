package types

// Constraint is a ConstraintDeclaration (A.7.2). It is not a Type, and that is
// the point: A.7.2 ⊢ "Vertex has no interfaces. A constraint is its own
// declaration form: a compile-time type set, optionally paired with required
// methods. It is never a value type and is legal only in a [...] position."
//
// A.14 lists `var c: Ordered` among the rejected forms. Because Constraint does
// not implement Type, that rejection is a type error the checker cannot avoid
// noticing, rather than a predicate someone has to remember to call.
//
// A.7.2 ⊢ multiple elements form an *intersection*: a type argument must
// satisfy all of them. Embedded constraints contribute their own sets, so
// satisfaction is a walk rather than a single membership test.
type Constraint struct {
	obj     *TypeName // nil for an inline ConstraintExpression
	terms   []Term    // the TypeSet; empty means `any`
	methods []*Func   // MethodRequirements
	embeds  []*Constraint
}

func NewConstraint(obj *TypeName, terms []Term, methods []*Func, embeds []*Constraint) *Constraint {
	return &Constraint{obj, terms, methods, embeds}
}

func (c *Constraint) Obj() *TypeName        { return c.obj }
func (c *Constraint) Terms() []Term         { return c.terms }
func (c *Constraint) Methods() []*Func      { return c.methods }
func (c *Constraint) Embeds() []*Constraint { return c.embeds }

// Term is one TypeSetTerm (A.7.3).
//
// A.7.3 ⊢ `~T` admits T and every type whose underlying type is T, so an alias
// to float32 still satisfies ~float32. A bare T admits only T exactly.
type Term struct {
	Tilde bool
	Type  Type
}

// Predeclared constraints (A.7.4). `any` is the empty intersection — every type
// satisfies it — and is what a bare type-parameter name means.
var (
	Any        = &Constraint{}
	Comparable = &Constraint{ /* populated by the universe scope */ }
)

// IsAny reports whether c admits every type.
func (c *Constraint) IsAny() bool {
	return c == Any || (c != nil && len(c.terms) == 0 && len(c.methods) == 0 && len(c.embeds) == 0)
}

// Satisfies reports whether t is admitted by c.
//
// A.7.5 ⊢ "constraint satisfaction is checked once per instantiation, at the
// instantiation site, never at runtime." This is that check.
func (c *Constraint) Satisfies(t Type) bool {
	if c == nil || c.IsAny() {
		return true
	}
	if c == Comparable {
		return IsComparable(t)
	}

	// A.7.2 ⊢ intersection: every element must admit t.
	if len(c.terms) > 0 && !termsAdmit(c.terms, t) {
		return false
	}
	for _, e := range c.embeds {
		if !e.Satisfies(t) {
			return false
		}
	}
	// A.7.2 ⊢ a MethodRequirement "is satisfied by any type declaring a
	// matching receiver method." Monomorphization lowers the call directly, so
	// this introduces no interface value and no vtable.
	for _, req := range c.methods {
		if !hasMethod(t, req) {
			return false
		}
	}
	return true
}

// termsAdmit is the union half: A.7.3 ⊢ `|` is union, so one matching term is
// enough.
func termsAdmit(terms []Term, t Type) bool {
	for _, term := range terms {
		if term.Tilde {
			if Identical(Underlying(t), Underlying(term.Type)) {
				return true
			}
			continue
		}
		if Identical(t, term.Type) {
			return true
		}
	}
	return false
}

func hasMethod(t Type, req *Func) bool {
	n, ok := t.(*Named)
	if !ok {
		return false
	}
	m := n.LookupMethod(req.name)
	if m == nil {
		return false
	}
	return matchesRequirement(m.Signature(), req.Signature())
}

// matchesRequirement compares a declared method against a MethodRequirement,
// ignoring the receiver — A.7.2's requirement names no receiver, and the
// receiver is exactly what varies across satisfying types.
func matchesRequirement(got, want *Signature) bool {
	if got == nil || want == nil {
		return false
	}
	if got.variadic != want.variadic || got.marker != want.marker {
		return false
	}
	return identicalTuple(got.params, want.params) &&
		identicalTuple(got.results, want.results)
}