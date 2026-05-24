package compiler

// SymKind classifies what a symbol refers to.
type SymKind int

const (
	SymVar    SymKind = iota // let / var binding
	SymFunc                   // local Vertex function
	SymNative                 // native class method (imported from C / syscall / GPU / …)
	SymType                   // struct / class / enum / alias
	SymParam                  // function parameter
)

// Symbol is a single entry in a scope.
type Symbol struct {
	Name    string
	Kind    SymKind
	Type    Type
	Mutable bool

	// Wasm code-gen fields
	FuncIdx  uint32 // SymFunc / SymNative: wasm function index
	LocalIdx uint32 // SymVar / SymParam: wasm local index
	FrameOff int32  // SymVar struct: offset inside function's linear-memory frame
	IsFrame  bool   // true → variable lives in linear memory (struct local)
}

// Scope is a name→Symbol map with parent chain.
type Scope struct {
	parent  *Scope
	entries map[string]*Symbol
}

func NewScope(parent *Scope) *Scope {
	return &Scope{parent: parent, entries: make(map[string]*Symbol)}
}

// Define adds sym. Returns false if the name is already taken in this scope.
func (s *Scope) Define(sym *Symbol) bool {
	if _, ok := s.entries[sym.Name]; ok {
		return false
	}
	s.entries[sym.Name] = sym
	return true
}

// Lookup searches this scope and all ancestors.
func (s *Scope) Lookup(name string) *Symbol {
	for sc := s; sc != nil; sc = sc.parent {
		if sym, ok := sc.entries[name]; ok {
			return sym
		}
	}
	return nil
}