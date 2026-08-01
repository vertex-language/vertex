package types

import "github.com/vertex-language/vertex/token"

// Scope is a lexical block holding named objects.
//
// A.2 ⊢ "top-level declarations are order-independent: a declaration may refer
// to any other declaration in the same package regardless of textual position."
// That makes package scope position-independent, unlike a block scope where a
// name is visible only after its declaration — hence Pos/End, which a lookup
// consults for block scopes and ignores for the package and universe scopes.
type Scope struct {
	parent   *Scope
	children []*Scope
	elems    map[string]Object
	comment  string

	pos, end token.Pos
	isFunc   bool
}

func NewScope(parent *Scope, comment string) *Scope {
	s := &Scope{parent: parent, comment: comment}
	if parent != nil {
		parent.children = append(parent.children, s)
	}
	return s
}

func (s *Scope) Parent() *Scope  { return s.parent }
func (s *Scope) Len() int        { return len(s.elems) }
func (s *Scope) Comment() string { return s.comment }

func (s *Scope) SetExtent(pos, end token.Pos) { s.pos, s.end = pos, end }

// Lookup returns the object with the given name declared directly in s.
func (s *Scope) Lookup(name string) Object {
	if s.elems == nil {
		return nil
	}
	return s.elems[name]
}

// LookupParent walks outward. It is what resolves an Identifier against A.1.4's
// implicit outermost scope when nothing nearer declares it.
func (s *Scope) LookupParent(name string) (*Scope, Object) {
	for ; s != nil; s = s.parent {
		if obj := s.Lookup(name); obj != nil {
			return s, obj
		}
	}
	return nil, nil
}

// Insert adds obj to s, returning the existing object if the name is taken.
// A non-nil return is a duplicate declaration for the caller to diagnose.
//
// A BlankIdentifier is dropped rather than inserted: A.1.2 ⊢ `_` "never
// introduces a usable binding."
func (s *Scope) Insert(obj Object) Object {
	name := obj.Name()
	if name == "_" {
		return nil
	}
	if alt := s.Lookup(name); alt != nil {
		return alt
	}
	if s.elems == nil {
		s.elems = make(map[string]Object)
	}
	s.elems[name] = obj
	return nil
}

// Names returns the declared names, unsorted. Callers that need determinism
// sort; this does not, because most callers are lookups.
func (s *Scope) Names() []string {
	out := make([]string, 0, len(s.elems))
	for n := range s.elems {
		out = append(out, n)
	}
	return out
}

// Universe is A.1.4's implicit outermost scope.
//
// It holds every PredeclaredTypeName, PredeclaredConstraintName, and
// ReservedBuiltinName, plus the three ReservedLiteralKeywords. The predeclared
// type and constraint names may be shadowed by a user declaration; the builtins
// may not — A.1.4 ⊢ "a ReservedBuiltinName may not be shadowed, and may not be
// declared as a member, method, or field name." Insert cannot express that
// asymmetry, so the resolver checks Reserved before inserting anything.
var Universe *Scope

// reserved is the unshadowable subset of Universe.
var reserved = map[string]bool{}

// Reserved reports whether name is a ReservedBuiltinName (A.1.4). A declaration
// of this name at any scope, or as a member, method, or field, is an error.
func Reserved(name string) bool { return reserved[name] }

func init() {
	Universe = NewScope(nil, "universe")

	for name, t := range predeclared {
		Universe.Insert(NewTypeName(token.NoPos, nil, name, t))
	}

	Universe.Insert(NewConstraintName(token.NoPos, nil, "any", Any))
	Universe.Insert(NewConstraintName(token.NoPos, nil, "comparable", Comparable))

	// true/false/nil are ReservedLiteralKeywords (A.1.3): Literals
	// syntactically, reserved lexically. The scanner already gives them their
	// own Kinds, so they need no scope entry — they can never be an Identifier.

	builtins := []struct {
		name string
		id   BuiltinId
	}{
		{"sizeof", _Sizeof}, {"alignof", _Alignof}, {"reinterpret", _Reinterpret},
		{"addr", _Addr}, {"new", _New}, {"delete", _Delete}, {"resize", _Resize},
		{"copy", _Copy}, {"zero", _Zero}, {"unique", _Unique}, {"shared", _Shared},
		{"weak", _Weak}, {"upgrade", _Upgrade}, {"drop", _Drop},
		{"transfer", _Transfer},
	}
	for _, b := range builtins {
		Universe.Insert(NewBuiltin(b.name, b.id))
		reserved[b.name] = true
	}
}