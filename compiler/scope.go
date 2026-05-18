package compiler

import "fmt"

// Symbol is a named value visible in a lexical scope.
type Symbol struct {
	Name    string
	Type    VType
	WasmIdx uint32 // WASM local index
	IsMut   bool   // false = let, true = var
}

// Scope is a lexical scope (block, function, …).
type Scope struct {
	symbols map[string]Symbol
	parent  *Scope
}

func newScope(parent *Scope) *Scope {
	return &Scope{symbols: make(map[string]Symbol), parent: parent}
}

func (s *Scope) define(name string, sym Symbol) error {
	if _, exists := s.symbols[name]; exists {
		return fmt.Errorf("invalid redeclaration of %q", name)
	}
	s.symbols[name] = sym
	return nil
}

// set updates an existing mutable binding anywhere in the scope chain.
func (s *Scope) set(name string, sym Symbol) error {
	for cur := s; cur != nil; cur = cur.parent {
		if existing, ok := cur.symbols[name]; ok {
			if !existing.IsMut {
				return fmt.Errorf("cannot assign to let constant %q", name)
			}
			cur.symbols[name] = sym
			return nil
		}
	}
	return fmt.Errorf("use of unresolved identifier %q", name)
}

func (s *Scope) lookup(name string) (Symbol, bool) {
	for cur := s; cur != nil; cur = cur.parent {
		if sym, ok := cur.symbols[name]; ok {
			return sym, true
		}
	}
	return Symbol{}, false
}

// FuncInfo describes a hoisted function (local or imported).
type FuncInfo struct {
	Name       string
	WasmIdx    uint32 // absolute WASM function index
	TypeIdx    uint32
	Params     []VType
	ParamNames []string // internal (local) parameter names
	Ret        VType
	IsImport   bool
	ASTNode    interface{} // *parser.FunctionDeclarationContext when !IsImport
}