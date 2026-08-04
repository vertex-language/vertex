package types

import "github.com/vertex-language/vertex/token"

// Scope is a lexical block holding named objects.
//
// §2.1 lists four scopes — universe, package, file, and local — and gives each
// its extent. The distinction this type has to carry is that package and
// universe scope are position-independent (§1.1 ⊢ top-level declarations "are
// order-independent") while a local scope is not (§2.1 ⊢ "a binding is visible
// from its declaration to the end of its scope").
//
// That is what pos/end are for, and LookupParent consults them: a scope with no
// recorded extent is position-independent, and a scope with one is not. §2.1
// also fixes where a local scope opens — "at every Block, at a Signature's
// parameter list, at a TypeParameters list, and at each CaseClause/SelectClause
// statement list" — which is the set the analyzer calls NewScope for.
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

// NewFuncScope opens a scope that a `defer` registers against.
//
// §6.6 ⊢ a deferred "call runs when the enclosing function returns, in reverse
// order of registration" — the enclosing *function*, not the enclosing block —
// so Info.Defers groups by the scope this marks. A FunctionLit opens one of its
// own, since §7.3 ⊢ it "begins with all enclosing parse context cleared".
func NewFuncScope(parent *Scope, comment string) *Scope {
	s := NewScope(parent, comment)
	s.isFunc = true
	return s
}

func (s *Scope) Parent() *Scope   { return s.parent }
func (s *Scope) Len() int         { return len(s.elems) }
func (s *Scope) Comment() string  { return s.comment }
func (s *Scope) NumChildren() int { return len(s.children) }
func (s *Scope) Child(i int) *Scope { return s.children[i] }
func (s *Scope) IsFunc() bool     { return s.isFunc }
func (s *Scope) Pos() token.Pos   { return s.pos }
func (s *Scope) End() token.Pos   { return s.end }

func (s *Scope) SetExtent(pos, end token.Pos) { s.pos, s.end = pos, end }

// Contains reports whether p falls inside s's recorded extent. A scope with no
// extent contains everything, which is the correct answer for package and
// universe scope.
func (s *Scope) Contains(p token.Pos) bool {
	if !s.pos.IsValid() {
		return true
	}
	return s.pos <= p && p < s.end
}

// FuncScope returns the innermost enclosing function scope, or nil. This is
// what a DeferStmt is recorded against.
func (s *Scope) FuncScope() *Scope {
	for ; s != nil; s = s.parent {
		if s.isFunc {
			return s
		}
	}
	return nil
}

// Lookup returns the object with the given name declared directly in s.
func (s *Scope) Lookup(name string) Object {
	if s.elems == nil {
		return nil
	}
	return s.elems[name]
}

// LookupParent walks outward, ending at Universe — §2.3's implicit outermost
// scope, which is how a predeclared name resolves when nothing nearer declares
// it.
//
// p is the position of the reference. In a scope with a recorded extent, an
// object declared after p is skipped, because §2.1 makes a local binding
// visible only "from its declaration to the end of its scope". Package and
// universe scope record no extent and therefore ignore p entirely, which is
// §1.1's order-independence. Pass token.NoPos to disable the check — what
// tooling asking "what is in scope here" wants.
func (s *Scope) LookupParent(name string, p token.Pos) (*Scope, Object) {
	for ; s != nil; s = s.parent {
		obj := s.Lookup(name)
		if obj == nil {
			continue
		}
		if p.IsValid() && s.pos.IsValid() && obj.Pos().IsValid() && p < obj.Pos() {
			continue
		}
		return s, obj
	}
	return nil, nil
}

// Innermost returns the innermost scope containing p, for tooling.
func (s *Scope) Innermost(p token.Pos) *Scope {
	if !s.Contains(p) {
		return nil
	}
	for _, c := range s.children {
		if in := c.Innermost(p); in != nil {
			return in
		}
	}
	return s
}

// Insert adds obj to s, returning the existing object if the name is taken.
// A non-nil return is a duplicate declaration for the caller to diagnose —
// §2.2 ⊢ "within a single scope a name is declared once", and ⊢ "Vertex has no
// overloading. One name denotes one declaration."
//
// The blank identifier is dropped rather than inserted: §2.4 ⊢ `_` "introduces
// no binding and may be repeated freely", and ⊢ "reading `_` is always an
// error: it names nothing."
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

// Universe is §2.3's implicit outermost scope.
//
// It holds all four predeclared families: the type names, the two constraint
// names, the tensor element type names, and the reserved builtins. §2.3 ⊢ "all
// are ordinary identifiers — no scanner recognizes them", which is exactly why
// they are here and not in token.Lookup.
//
// Each family has its own legality rule, and none of them is expressible by
// membership alone: a constraint name is legal only in a bracket position, a
// tensor element name only inside an npu body, and a builtin may not be
// shadowed. The first two are the analyzer's checks; the third is Reserved.
var Universe *Scope

// Reserved reports whether name may not be shadowed or declared.
//
// It delegates to token, which is the one home of the set — token.go ⊢ "the set
// lives here so the non-shadowing guarantee has one home". Keeping a second
// list here is how the two drift, and a drift would silently re-legalize
// shadowing a builtin.
func Reserved(name string) bool { return token.IsReservedBuiltin(name) }

// builtinIds maps each reserved builtin spelling to its id. Every name in
// token's set appears here; init panics if one does not, since a builtin
// reachable by name but with no object would resolve to nothing at its first
// call site.
var builtinIds = map[string]BuiltinId{
	"new": NewId, "delete": DeleteId, "resize": ResizeId, "copy": CopyId,
	"zero": ZeroId, "addr": AddrId, "sizeof": SizeofId, "alignof": AlignofId,
	"reinterpret": ReinterpretId, "upgrade": UpgradeId, "drop": DropId,
	"panic": PanicId, "blend": BlendId, "min": MinId, "max": MaxId,
	"clamp": ClampId, "transfer": TransferId,
}

func init() {
	Universe = NewScope(nil, "universe")

	for name, t := range predeclaredTypes {
		Universe.Insert(NewTypeName(token.NoPos, nil, name, t))
	}
	for name, t := range predeclaredTensorElems {
		Universe.Insert(NewTypeName(token.NoPos, nil, name, t))
	}

	Universe.Insert(NewConstraintName(token.NoPos, nil, "any", Any))
	Universe.Insert(NewConstraintName(token.NoPos, nil, "comparable", Comparable))

	// true/false/nil are reserved literal keywords: literals syntactically,
	// reserved lexically. The scanner gives them their own Kinds, so they can
	// never be an identifier and need no entry here.

	for name, id := range builtinIds {
		if !token.IsReservedBuiltin(name) {
			panic("types: " + name + " is not a reserved builtin in token")
		}
		Universe.Insert(NewBuiltin(name, id))
	}
}