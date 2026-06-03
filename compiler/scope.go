package compiler

// SymbolKind classifies a named symbol.
type SymbolKind int

const (
	SymVar       SymbolKind = iota // let / var binding
	SymFunc                        // function or associated function
	SymStruct                      // struct type
	SymClass                       // class type
	SymEnum                        // enum type
	SymTypeAlias                   // type alias
	SymEnumCase                    // individual enum case
	SymParam                       // function parameter
)

// Symbol is a named entry in a Scope.
type Symbol struct {
	Name    string
	Kind    SymbolKind
	Type    VType
	Decl    Node
	IsConst bool // true for let bindings
}

// Scope is a lexical scope that chains to its parent.
type Scope struct {
	parent  *Scope
	symbols map[string]*Symbol
}

// NewScope creates a new scope with the given parent (may be nil).
func NewScope(parent *Scope) *Scope {
	return &Scope{parent: parent, symbols: make(map[string]*Symbol)}
}

// Define adds sym to this scope.  Silently replaces an existing entry with the
// same name (the resolver emits a diagnostic before calling Define in that case).
func (s *Scope) Define(sym *Symbol) {
	s.symbols[sym.Name] = sym
}

// Lookup walks the scope chain and returns the first matching symbol.
func (s *Scope) Lookup(name string) (*Symbol, bool) {
	for cur := s; cur != nil; cur = cur.parent {
		if sym, ok := cur.symbols[name]; ok {
			return sym, true
		}
	}
	return nil, false
}

// LookupLocal returns a symbol defined in this scope only.
func (s *Scope) LookupLocal(name string) (*Symbol, bool) {
	sym, ok := s.symbols[name]
	return sym, ok
}

// Symbols returns a snapshot of all symbols in this scope (not parents).
func (s *Scope) Symbols() map[string]*Symbol {
	out := make(map[string]*Symbol, len(s.symbols))
	for k, v := range s.symbols {
		out[k] = v
	}
	return out
}

// newGlobalScope builds the root scope pre-populated with built-in type names.
func newGlobalScope() *Scope {
	s := NewScope(nil)
	for name, vtype := range BuiltinTypes {
		s.Define(&Symbol{
			Name:    name,
			Kind:    SymTypeAlias,
			Type:    vtype,
			IsConst: true,
		})
	}
	return s
}