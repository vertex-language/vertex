package types

import "github.com/vertex-language/vertex/ast"

// Info records what analysis learned. It is the only file in this package that
// imports ast, and the dependency runs one way: ast never imports types.
//
// Every result is a side table keyed by syntax node, never a field written back
// onto one. ast's own doc comment says the tree "records shape, never meaning",
// and keeping that true is what lets a printer round-trip a checked file and
// lets the checker run twice over one tree without residue.
//
// A nil map means "do not record this". Each recorder checks before storing, so
// a caller that wants only Uses pays for only Uses.
type Info struct {
	// Types records the type and value of every expression, including the
	// nodes ast represents as Exprs but that denote types — A.3.6's IndexExpr,
	// a TypeName's Ident, an OwnershipType. TypeAndValue.IsType distinguishes
	// them, which is the recorded answer to the ambiguity A.3.6 calls "the one
	// syntactic overlap resolved by the operand's meaning rather than by shape."
	Types map[ast.Expr]TypeAndValue

	// Defs maps an identifier to the object it declares. The entry is nil for a
	// BlankIdentifier, which A.1.2 ⊢ "never introduces a usable binding" — the
	// key is still present so tooling can point at it.
	Defs map[*ast.Ident]Object

	// Uses maps an identifier to the object it refers to.
	Uses map[*ast.Ident]Object

	// Selections records what a SelectorExpr resolved to: a struct field, a
	// method, a package member, or a member of one of A.4.1's four keyword
	// namespaces.
	Selections map[*ast.SelectorExpr]*Selection

	// Instances records the type arguments and resulting type at each generic
	// instantiation site, keyed by the identifier naming the generic.
	//
	// A.7.5 ⊢ "constraint satisfaction is checked once per instantiation, at
	// the instantiation site, never at runtime." This is the record of those
	// sites, and it is the set monomorphization enumerates in lower — A.7.2 ⊢
	// "because every instantiation is monomorphized, the call in the generic
	// body lowers to a direct call on the concrete type."
	Instances map[*ast.Ident]Instance

	// Scopes maps each scope-introducing node to its Scope.
	Scopes map[ast.Node]*Scope

	// Transfers records every TransferExpr the ownership pass accepted, mapped
	// to the binding that dies at that point.
	//
	// A.4.6 ⊢ the marker's presence "is the entire difference between move and
	// deep copy. One marker, two meanings, read by presence." This is therefore
	// the one analysis result that changes generated code rather than merely
	// licensing it: A.9.4 puts a bare copy at O(data) and a transfer at O(1),
	// and lower reads this table to decide which to emit.
	Transfers map[*ast.TransferExpr]*Var

	// Defers records each DeferStmt's enclosing scope.
	//
	// A.5.8 ⊢ deferred calls "are collected per scope and emitted in reverse
	// registration order on every exit edge — fall-through, return, break,
	// continue. Because there is no unwinder, 'every exit edge' is a finite,
	// statically known set." Grouping by scope is what makes that set
	// enumerable in lower.
	Defers map[*ast.DeferStmt]*Scope
}

// NewInfo returns an Info with every table allocated. A checker that wants a
// subset builds its own literal instead.
func NewInfo() *Info {
	return &Info{
		Types:      make(map[ast.Expr]TypeAndValue),
		Defs:       make(map[*ast.Ident]Object),
		Uses:       make(map[*ast.Ident]Object),
		Selections: make(map[*ast.SelectorExpr]*Selection),
		Instances:  make(map[*ast.Ident]Instance),
		Scopes:     make(map[ast.Node]*Scope),
		Transfers:  make(map[*ast.TransferExpr]*Var),
		Defers:     make(map[*ast.DeferStmt]*Scope),
	}
}

// ---------------------------------------------------------- type and value

// OperandMode classifies what an expression denotes.
//
// It is exported because the analyzer is a separate package and is the only
// thing that ever sets one. The distinction between VarMode and ValueMode is
// not stylistic: A.4.8 ⊢ addr "requires an addressable operand: a var binding
// or a field path", and A.3.2 ⊢ a mut argument "must be an addressable var
// binding or field path". Both read Addressable, which reads this.
type OperandMode uint8

const (
	InvalidMode OperandMode = iota

	// NoValue is a call to a function with no result. A.3.4 ⊢ "omitting
	// -> Type is the void form. There is no void type name" — so this mode,
	// rather than a distinct type, is what marks it.
	NoValue

	// BuiltinMode is a ReservedBuiltinName that has not been called. A.4.8
	// makes each builtin's shape arity-specific, so it has no Signature to
	// carry and cannot be used as a value.
	BuiltinMode

	// TypeMode marks an expression that denotes a type rather than a value.
	TypeMode

	// ConstMode is a constant expression; Value is non-nil.
	ConstMode

	// VarMode is an addressable value: a var binding, a field path, or a
	// dereferenced typed_ptr.
	VarMode

	// ValueMode is a computed value with no storage slot.
	ValueMode

	// PkgMode is an import qualifier, legal only as the base of a selector.
	PkgMode
)

// TypeAndValue holds an expression's type and, when constant, its value.
type TypeAndValue struct {
	Mode  OperandMode
	Type  Type
	Value Value // non-nil only when Mode is ConstMode
}

func NewTypeAndValue(mode OperandMode, t Type, v Value) TypeAndValue {
	return TypeAndValue{Mode: mode, Type: t, Value: v}
}

func (tv TypeAndValue) IsType() bool    { return tv.Mode == TypeMode }
func (tv TypeAndValue) IsBuiltin() bool { return tv.Mode == BuiltinMode }
func (tv TypeAndValue) IsVoid() bool    { return tv.Mode == NoValue }
func (tv TypeAndValue) IsPackage() bool { return tv.Mode == PkgMode }

func (tv TypeAndValue) IsValue() bool {
	switch tv.Mode {
	case ConstMode, VarMode, ValueMode:
		return true
	}
	return false
}

// Addressable reports whether the expression has a storage slot.
//
// A.5.1 ⊢ a `let` "is not guaranteed to be addressable — it may be a register,
// an SSA value, or folded away entirely", which is why only a `var` binding or
// a field path answers true.
func (tv TypeAndValue) Addressable() bool { return tv.Mode == VarMode }

// --------------------------------------------------------------- selection

// Selection is what a SelectorExpr resolved to.
type Selection struct {
	Kind  SelectionKind
	Recv  Type   // the base's type; nil for a package or namespace member
	Obj   Object // the field, method, or member
	Index int    // field index for FieldVal, element index for TupleIndex
}

type SelectionKind uint8

const (
	// FieldVal is a struct or class field access.
	FieldVal SelectionKind = iota

	// MethodVal is a method on a Named type. A.6.3 ⊢ there is "no inheritance,
	// no vtable, and no dynamic dispatch... Every call is direct", so this
	// always names one concrete function.
	MethodVal

	// TupleIndex is A.4.3's positional access, `t.0`. It is a Selection rather
	// than a field because ast gives it its own node (TupleIndexExpr) and the
	// index is decoded, not named.
	TupleIndex

	// PackageMember is a member reached through an import qualifier (A.2.3).
	PackageMember

	// NamespaceMember is a member of one of A.4.1's four keyword namespaces.
	// A.11.3 ⊢ the npu. member set "is closed. Its members are not declarable,
	// shadowable, or extensible", so this never resolves to a user object.
	NamespaceMember
)

// ---------------------------------------------------------------- instance

// Instance is one generic instantiation.
type Instance struct {
	TypeArgs []Type
	Type     Type
}

// --------------------------------------------------------------- recorders

func (info *Info) RecordType(x ast.Expr, tv TypeAndValue) {
	if info == nil || info.Types == nil || x == nil {
		return
	}
	info.Types[x] = tv
}

// RecordDef records a declaring identifier. A BlankIdentifier is recorded with
// a nil object rather than skipped: A.1.2 ⊢ `_` "never introduces a usable
// binding", but tooling still needs to know the parser saw one there.
func (info *Info) RecordDef(id *ast.Ident, obj Object) {
	if info == nil || info.Defs == nil || id == nil {
		return
	}
	info.Defs[id] = obj
}

func (info *Info) RecordUse(id *ast.Ident, obj Object) {
	if info == nil || info.Uses == nil || id == nil || obj == nil {
		return
	}
	info.Uses[id] = obj
}

func (info *Info) RecordSelection(x *ast.SelectorExpr, sel *Selection) {
	if info == nil || info.Selections == nil || x == nil {
		return
	}
	info.Selections[x] = sel
}

func (info *Info) RecordInstance(id *ast.Ident, inst Instance) {
	if info == nil || info.Instances == nil || id == nil {
		return
	}
	info.Instances[id] = inst
}

func (info *Info) RecordScope(n ast.Node, s *Scope) {
	if info == nil || info.Scopes == nil || n == nil {
		return
	}
	info.Scopes[n] = s
}

func (info *Info) RecordTransfer(x *ast.TransferExpr, v *Var) {
	if info == nil || info.Transfers == nil || x == nil {
		return
	}
	info.Transfers[x] = v
}

func (info *Info) RecordDefer(x *ast.DeferStmt, s *Scope) {
	if info == nil || info.Defers == nil || x == nil {
		return
	}
	info.Defers[x] = s
}

// --------------------------------------------------------------- accessors

// TypeOf returns an expression's type, or nil if none was recorded.
func (info *Info) TypeOf(e ast.Expr) Type {
	if info == nil {
		return nil
	}
	if tv, ok := info.Types[e]; ok {
		return tv.Type
	}
	if id, _ := e.(*ast.Ident); id != nil {
		if obj := info.ObjectOf(id); obj != nil {
			return obj.Type()
		}
	}
	return nil
}

// ObjectOf returns the object an identifier denotes, whether it declares or
// uses it. A BlankIdentifier answers nil in both directions.
func (info *Info) ObjectOf(id *ast.Ident) Object {
	if info == nil || id == nil {
		return nil
	}
	if obj, ok := info.Defs[id]; ok && obj != nil {
		return obj
	}
	return info.Uses[id]
}

// SelectionOf returns what a selector resolved to, or nil.
func (info *Info) SelectionOf(x *ast.SelectorExpr) *Selection {
	if info == nil {
		return nil
	}
	return info.Selections[x]
}

// IsTransfer reports whether a TransferExpr was accepted as a move.
//
// A.4.6 ⊢ the marker's presence is what separates move from deep copy, so lower
// asks this at every owning position rather than re-deriving it from the tree.
func (info *Info) IsTransfer(x *ast.TransferExpr) bool {
	if info == nil || info.Transfers == nil {
		return false
	}
	_, ok := info.Transfers[x]
	return ok
}