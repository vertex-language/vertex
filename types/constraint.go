package types

// Constraint is a ConstraintDecl, and it is not a Type.
//
// §9 ⊢ "a ConstraintDecl takes no type parameters, is legal only in a `[`…`]`
// position, and is never a value type." Because Constraint does not implement
// Type, `var c: Ordered` is a shape error the checker cannot avoid noticing,
// rather than a predicate someone has to remember to call.
//
// §9 ⊢ "multiple elements in its body are an intersection; `|` within one
// element is a union." Embedded constraints contribute their own sets, so
// satisfaction is a walk rather than a single membership test.
//
// There are no interfaces: grammar.md ⊢ "a constraint is its own declaration
// form", and §9.1 ⊢ the declaration "is checked once, on its own terms" against
// what the constraint permits, whatever the instantiations happen to supply.
type Constraint struct {
	obj     *TypeName // nil for an inline TypeSet written at the use site
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

// Term is one TypeSetTerm.
//
// grammar.md ⊢ "`~T` admits every type whose underlying type is `Type`; a bare
// `Type` admits only that type exactly." §3.1 makes this the one place an alias
// and its target are distinguishable at all, which is why the check below is
// Identical over Underlying rather than Identical alone.
type Term struct {
	Tilde bool
	Type  Type
}

// The predeclared constraints. §2.3 lists exactly two, `any` and `comparable`,
// both "only in a `[`…`]` position, never as a `Type`".
//
// Any is the empty intersection — every type satisfies it — and is what a bare
// TypeParamName means. Comparable is identified by pointer rather than by
// content, because §3.5 pins its membership as a predicate (IsComparable) and
// not as a term list.
var (
	Any        = &Constraint{}
	Comparable = &Constraint{}
)

// IsAny reports whether c admits every type.
func (c *Constraint) IsAny() bool {
	return c == Any || (c != nil && len(c.terms) == 0 && len(c.methods) == 0 && len(c.embeds) == 0)
}

// Satisfies reports whether t is admitted by c.
//
// §9.2 ⊢ "constraint satisfaction is checked per instantiation, at the
// instantiation site, and a failure is a compile error there." This is that
// check; there is no runtime counterpart to it.
func (c *Constraint) Satisfies(t Type) bool {
	if c == nil || c.IsAny() {
		return true
	}
	if c == Comparable {
		return IsComparable(t)
	}

	if len(c.terms) > 0 && !termsAdmit(c.terms, t) {
		return false
	}
	for _, e := range c.embeds {
		if !e.Satisfies(t) {
			return false
		}
	}
	// A MethodRequirement is satisfied by any type declaring a matching
	// receiver method. §2.1 keeps method names out of every scope, so this is a
	// lookup on the type and not a lookup by name.
	for _, req := range c.methods {
		if !hasMethod(t, req) {
			return false
		}
	}
	return true
}

// termsAdmit is the union half: one matching term is enough.
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
// ignoring the receiver — the requirement names none, and the receiver is
// exactly what varies across satisfying types.
//
// The marker is compared, because grammar.md ⊢ "because MethodRequirement takes
// a full Signature, a constraint can require a marked method", and §7.4 makes
// the marker part of the type at both ends.
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