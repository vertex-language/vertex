// Package hir is where every decision is made.
//
// Given a checked package graph and its *types.Info, hir produces a
// *hir.Program: monomorphic, ownership-explicit, control-flow-flattened,
// with builtin calls named and nothing left that a later phase has to
// reinterpret. lower/vir is mechanical afterward by construction — if it
// ever needs a type switch on "is this an owning type" or "was this
// transferred," that logic belongs here.
//
// Two invariants govern the package and are checkable from outside it:
// hir never imports vvm, and hir never sees a target triple (overview §7,
// invariants 1 and 10). Layout is the one target-shaped fact it consumes,
// and it arrives as a types.Sizes on Config, chosen by the driver from the
// build tag.
//
// # Representation
//
// A Program is a list of Modules, one per originating Vertex package plus
// one per canonical owning module for synthesized symbols. A Module maps
// onto vir's namespace/module split exactly: Path is the import path and
// becomes vir's `namespace`, Name is the PackageClause name and becomes
// vir's `module` (A.2.3's "the path is a locator, the declared name is the
// qualifier," carried through to the linker unchanged).
//
// A Func holds structured statements (Body) until Flatten runs, and flat
// Blocks afterward. Both shapes hold the same *Instr values, so there is
// one instruction representation across the whole pipeline.
//
// # Aggregates are pointers
//
// Vertex IR makes struct and array memory-only — never held in a named
// value. So a hir Value whose Type is aggregate (IsAggregate) *is* a ptr at
// the vir level, pointing at storage the emitting code owns. Aggregate
// parameters carry ByVal, aggregate results carry SRet, and aggregate
// assignment is a memcopy of Size bytes. This is stated once, here, and
// lower/vir reads it off IsAggregate rather than re-deriving it.
package hir

import (
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// Program is the whole compilation: every module, in dependency order.
type Program struct {
	// Modules in topological order — a module appears after everything it
	// imports, which is the order lower/vir must emit in and the order
	// vvm's importer expects.
	Modules []*Module

	// Entry is the synthesized program-entry wrapper (§5.2, tier 2): the
	// thing that starts the reactor, drives user main to completion, tears
	// down, and returns a status. Under ModeTest the same slot holds the
	// wrapper around one test function. Never a user Func.
	Entry *Func

	// Features is what the program actually uses. The driver hands it to
	// builtins.Modules, which resolves the closure. This is the answer to
	// the overview's open question about where the feature set is
	// computed: hir returns it alongside the Program rather than making
	// the driver re-derive it by scanning emitted import lines.
	Features builtins.FeatureSet
}

// Lookup returns the module for an import path, or nil.
func (p *Program) Lookup(path string) *Module {
	for _, m := range p.Modules {
		if m.Path == path {
			return m
		}
	}
	return nil
}

// Module is one vir module's worth of declarations.
type Module struct {
	Path string // import path   -> vir `namespace`
	Name string // PackageClause -> vir `module`

	Structs []*Struct
	Consts  []*Const
	Globals []*Global
	Funcs   []*Func

	// Imports holds module *names* (not paths): lower/vir emits one
	// `import "<name>"` line each, and every cross-module call operand is
	// spelled `<name>.<symbol>`.
	Imports []string

	// Links and Externs are non-empty only for a module lowered from a
	// file carrying a declare block (A.8.1), which is the single exception
	// to invariant 3. Interop is the only path by which user source
	// touches the linker.
	Links   []Link
	Externs []*ExternGroup

	names map[string]bool // vir enforces one flat namespace per module
	seen  map[string]bool // Imports dedup
}

func newModule(path, name string) *Module {
	return &Module{Path: path, Name: name, names: map[string]bool{}, seen: map[string]bool{}}
}

// AddImport records a cross-module dependency by module name, once.
func (m *Module) AddImport(name string) {
	if name == "" || name == m.Name || m.seen[name] {
		return
	}
	m.seen[name] = true
	m.Imports = append(m.Imports, name)
}

// uniqueName reserves a name in the module's flat namespace (§2.2), adding
// a numeric suffix on collision. Identifiers are [A-Za-z_][A-Za-z0-9_]*.
func (m *Module) uniqueName(base string) string {
	base = sanitize(base)
	if !m.names[base] {
		m.names[base] = true
		return base
	}
	for i := 1; ; i++ {
		n := base + "_" + itoa(i)
		if !m.names[n] {
			m.names[n] = true
			return n
		}
	}
}

func (m *Module) lookupFunc(name string) *Func {
	for _, f := range m.Funcs {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// LinkKind mirrors A.8's linkage boundary as vir spells it (§7.2).
type LinkKind string

const (
	LinkStatic    LinkKind = "static"
	LinkShared    LinkKind = "shared"
	LinkFramework LinkKind = "framework"
)

type Link struct {
	Kind LinkKind
	Name string
}

// ExternGroup attributes imported symbols to the dependency that provides
// them. Dependency matches a Link.Name byte-for-byte.
type ExternGroup struct {
	Dependency string
	Funcs      []*ExternFunc
}

type ExternFunc struct {
	Name     string
	Params   []*Param
	Variadic bool
	Result   Type
}

// Struct is a memory-only aggregate. It produces no symbol (§2.2: an
// exported struct is a shape-visibility flag), so declaring the same
// synthesized header in two modules is harmless — unlike a tier-2 function,
// which needs exactly one owning module.
type Struct struct {
	Name   string
	Module *Module
	Fields []Field
	Size   int64
	Align  int64

	// Origin is the types.Type this shape was derived from, or nil for a
	// synthesized header (a slice, a string, a tuple, a payload enum).
	Origin types.Type
}

func (s *Struct) Field(name string) (Field, bool) {
	for _, f := range s.Fields {
		if f.Name == name {
			return f, true
		}
	}
	return Field{}, false
}

type Field struct {
	Name   string
	Type   Type
	Offset int64
}

// Const is a compile-time scalar. It yields a direct value and occupies no
// runtime storage (§6.2).
type Const struct {
	Name  string
	Type  Type
	Value Value
}

// Global is module-level storage. A.2 already requires a top-level
// VariableDeclaration to have a compile-time-evaluable initializer, and
// vir's init grammar is narrower still, so hir folds initializers down to
// these forms and there is no initialization-time code to order.
type Global struct {
	Name   string
	Type   Type
	Init   Init
	Export bool
	TLS    bool
	Align  int
}

type Init interface{ initNode() }

type InitZero struct{}
type InitConst struct{ Value Value }
type InitBytes struct{ Data []byte }
type InitAggregate struct{ Elems []Init }

// InitAddr is a relocated pointer to an earlier function or global. It is
// the only way to name a function's address, since vir has no address-of
// instruction — see funcAddr in expr.go.
type InitAddr struct{ Name string }

func (InitZero) initNode()      {}
func (InitConst) initNode()     {}
func (InitBytes) initNode()     {}
func (InitAggregate) initNode() {}
func (InitAddr) initNode()      {}

// FuncKind records why a Func exists. Every kind but FuncUser is tier 2
// (§5.2): type-dependent, target-agnostic, synthesized here.
type FuncKind uint8

const (
	FuncUser FuncKind = iota
	FuncCopy         // per-type deep copy
	FuncDeinit       // per-type teardown
	FuncDrop         // state-machine payload drop
	FuncStateMachine // an async body's poll function
	FuncEntryShim    // the program entry / test wrapper
)

type Func struct {
	Name   string
	Module *Module

	Params []*Param
	Result Type

	Export   bool
	Entry    bool // vir's `entry` attribute; at most one per module
	NoReturn bool

	Kind   FuncKind
	Origin types.Object // nil for synthesized
	Pos    token.Pos

	// Body is the structured form, valid until Flatten. Blocks is the flat
	// form, valid after. Exactly one is non-nil in a well-formed Func.
	Body   *Seq
	Blocks []*Block

	names map[string]int
}

// Param is one declared parameter. ByVal and SRet carry vir's aggregate
// conventions; both imply the parameter's vir type is ptr.
type Param struct {
	Name  string
	Type  Type
	ByVal *Struct
	SRet  *Struct
}

func (f *Func) fresh(base string) string {
	if f.names == nil {
		f.names = map[string]int{}
	}
	base = sanitize(base)
	n := f.names[base]
	f.names[base] = n + 1
	if n == 0 {
		return base
	}
	return base + "_" + itoa(n)
}

func sanitize(s string) string {
	if s == "" {
		return "_"
	}
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c == '_':
			out = append(out, c)
		case c >= '0' && c <= '9':
			if len(out) == 0 {
				out = append(out, '_')
			}
			out = append(out, c)
		default:
			out = append(out, '_')
		}
	}
	return string(out)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}