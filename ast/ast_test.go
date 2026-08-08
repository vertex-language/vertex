package ast

import (
	"reflect"
	"testing"

	"github.com/vertex-language/vertex/token"
)

// allNodes lists one value of every exported node type. The Walk and span
// tests iterate it, so a new node type that isn't added here shows up as a
// Walk panic in the very next run — which is §5.4's tripwire working.
func allNodes() []Node {
	id := func(lo, hi token.Pos) *Ident { return &Ident{NamePos: lo, NameEnd: hi} }
	lit := &BasicLit{Kind: token.NUMBER, ValuePos: 10, ValueEnd: 11}

	return []Node{
		id(1, 2),
		&PrivateIdent{HashPos: 1, NameEnd: 5},
		&QualifiedName{X: id(1, 2), Sel: id(3, 6)},
		lit,
		&BadExpr{From: 1, To: 3},
		&BadStmt{From: 1, To: 3},
		&BadDecl{From: 1, To: 3},
		&BadType{From: 1, To: 3},
		&ParenExpr{Lparen: 1, X: lit, Rparen: 12},
		&BinaryExpr{X: lit, OpPos: 12, OpEnd: 13, Op: token.ADD, Y: lit},
		&UnionType{Types: []TypeExpr{&ThisType{ThisPos: 1, ThisEnd: 5}}},
		&EmptyStmt{Semi: 1},
		// ... one per type; the grammar-coverage harness in §7 keeps this honest
	}
}

// §1: every node has a non-zero span, and §5.4: End() is exact.
func TestNonZeroSpans(t *testing.T) {
	for _, n := range allNodes() {
		if !n.Pos().IsValid() {
			t.Errorf("%T: Pos() is NoPos", n)
		}
		if !n.End().IsValid() {
			t.Errorf("%T: End() is NoPos", n)
		}
		if n.End() <= n.Pos() {
			t.Errorf("%T: End() %d <= Pos() %d", n, n.End(), n.Pos())
		}
	}
}

// §5.4: Walk panics on unknown node types.
type unknownNode struct{}

func (unknownNode) Pos() token.Pos { return 1 }
func (unknownNode) End() token.Pos { return 2 }

func TestWalkPanicsOnUnknown(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Walk did not panic on an unknown node type")
		}
	}()
	Inspect(unknownNode{}, func(Node) bool { return true })
}

// Walk must reach every child. A node with an unvisited field is a traversal
// bug that no other test catches.
func TestWalkVisitsEveryChild(t *testing.T) {
	for _, n := range allNodes() {
		want := childCount(n)
		got := 0
		Inspect(n, func(Node) bool { got++; return true })
		if got != want+1 { // +1 for n itself
			t.Errorf("%T: Walk visited %d nodes, %d children reachable by reflection", n, got-1, want)
		}
	}
}

// childCount counts reachable Node-typed fields by reflection, independently of
// Walk's hand-written switch.
func childCount(n Node) int {
	var count func(v reflect.Value) int
	count = func(v reflect.Value) int {
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr:
			if v.IsNil() {
				return 0
			}
			if _, ok := v.Interface().(Node); ok && v.Kind() == reflect.Ptr {
				return 1 + count(v.Elem())
			}
			return count(v.Elem())
		case reflect.Struct:
			total := 0
			for i := 0; i < v.NumField(); i++ {
				if v.Type().Field(i).PkgPath == "" {
					total += count(v.Field(i))
				}
			}
			return total
		case reflect.Slice:
			total := 0
			for i := 0; i < v.Len(); i++ {
				total += count(v.Index(i))
			}
			return total
		}
		return 0
	}
	return count(reflect.ValueOf(n)) - 1
}

func TestUnparen(t *testing.T) {
	inner := &BasicLit{Kind: token.NUMBER, ValuePos: 3, ValueEnd: 4}
	x := Expr(&ParenExpr{Lparen: 1, X: &ParenExpr{Lparen: 2, X: inner, Rparen: 4}, Rparen: 5})
	if Unparen(x) != Expr(inner) {
		t.Error("Unparen did not strip nested parentheses")
	}
	if Unparen(inner) != Expr(inner) {
		t.Error("Unparen altered an unparenthesized expression")
	}
}

// §5.2: Set and List must agree, or a diagnostic points at the wrong word.
func TestModifiersConsistent(t *testing.T) {
	m := Modifiers{}
	for _, tok := range []ModifierTok{
		{Bit: ModStatic, Pos: 1, End: 7},
		{Bit: ModReadonly, Pos: 8, End: 16},
	} {
		m.List = append(m.List, tok)
		m.Set |= tok.Bit
	}
	var union ModifierSet
	for _, tok := range m.List {
		union |= tok.Bit
	}
	if union != m.Set {
		t.Errorf("Set = %b, List implies %b", m.Set, union)
	}
	if got := m.Find(ModReadonly); got == nil || got.Pos != 8 {
		t.Error("Find(ModReadonly) did not return the right token")
	}
	if m.Find(ModAbstract) != nil {
		t.Error("Find returned a modifier that isn't present")
	}
}

func TestReleaseIdempotent(t *testing.T) {
	var count int
	f := &File{FileEnd: 1, Arena: releaserFunc(func() { count++ })}
	f.Release()
	f.Release()
	if count != 1 {
		t.Errorf("Release ran %d times, want 1", count)
	}
}

type releaserFunc func()

func (r releaserFunc) Release() { r() }