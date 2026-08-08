package parser

import (
	"strings"
	"testing"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
)

func parse(t *testing.T, src string) (*ast.File, []token.Diagnostic) {
	t.Helper()
	f := token.NewFile("t.vs", []byte(src))
	tree, diags := ParseFile(f, 0)
	if tree == nil {
		t.Fatal("ParseFile returned a nil tree; §6 forbids it")
	}
	return tree, diags
}

// §6.2 and grammar L: one fixture per TerminatedByASI entry. The seven
// type-member signatures are the ones that break silently, so they lead.
func TestSemicolonInsertion(t *testing.T) {
	for _, src := range []string{
		"let a = 1\nlet b = 2\n",
		"const x = 1\n",
		"using r = open()\n",
		"async function f() { await using r = open()\n }\n",
		"a++\nb--\n",
		"do x() ; while (c)\ny()\n",
		"function f() { return\n }\n",
		"outer: while (c) { break outer\n }\n",
		"function f() { throw e\n }\n",
		"debugger\n",
		"import \"m\"\n",
		"import x = y\n",
		"export { }\n",
		"class C { x = 1\n y = 2\n }\n",
		"type T = number\n",
		"declare var g: number\n",
		// The seven type-member signatures.
		"interface U { id: string\n name: string\n }\n",
		"interface U { f(): void\n g(): void\n }\n",
		"interface U { (): void\n (x: number): void\n }\n",
		"interface U { new (): U\n new (x: number): U\n }\n",
		"interface U { [k: string]: number\n }\n",
		"interface U { get x(): number\n get y(): number\n }\n",
		"interface U { set x(v: number)\n set y(v: number)\n }\n",
	} {
		if _, diags := parse(t, src); len(diags) != 0 {
			t.Errorf("%q: %v", src, diags[0].Msg)
		}
	}
}

// §4.2 and token.JoinGT: the scanner under-munches `>` and the expression
// parser joins. Both sides must work in the same file.
func TestGreaterThanJoining(t *testing.T) {
	for _, src := range []string{
		"let a: Array<FixedArray<int32, 4>> = init;\n",
		"let a: Box<T>= init;\n",
		"let x = a >> b;\n",
		"let x = a >>> b;\n",
		"let x = a >= b;\n",
		"x >>= 2;\n",
		"x >>>= 2;\n",
		"let x = a > b;\n",
	} {
		if _, diags := parse(t, src); len(diags) != 0 {
			t.Errorf("%q: %v", src, diags[0].Msg)
		}
	}
	// Adjacency is the whitespace test (§4.2), so this must not join.
	if _, diags := parse(t, "let x = a > > b;\n"); len(diags) == 0 {
		t.Error("`a > > b` must not parse as a shift")
	}
}

// §4.3: a line break after `struct`, `kernel`, or `graph` forces the
// identifier reading, which is what keeps them usable as names.
func TestContextualKeywordsStayIdentifiers(t *testing.T) {
	for _, src := range []string{
		"let struct = 1;\n",
		"let kernel = 1;\n",
		"let graph = 1;\n",
		"struct\n= 1;\n",
		"kernel\n= 1;\n",
		"struct S { x: int32; }\n",
		"kernel function k() {}\n",
		"graph function g() {}\n",
	} {
		if _, diags := parse(t, src); len(diags) != 0 {
			t.Errorf("%q: %v", src, diags[0].Msg)
		}
	}
}

// §5.3 and §8.4: kernel and graph functions are FuncDecl with Accel set —
// there is no KernelDecl to match on.
func TestAccelIsAField(t *testing.T) {
	tree, _ := parse(t, "kernel function k() {}\ngraph function g() {}\nfunction p() {}\n")
	var got []ast.AccelKind
	ast.Inspect(tree, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			got = append(got, fn.Accel)
		}
		return true
	})
	want := []ast.AccelKind{ast.AccelKernel, ast.AccelGraph, ast.AccelNone}
	if len(got) != len(want) {
		t.Fatalf("found %d FuncDecls, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("FuncDecl %d: Accel = %v, want %v", i, got[i], want[i])
		}
	}
}

// §6.3: forms the grammar rules out by early error parse into real nodes, so
// the message can name the construct.
func TestParseFirstRejectLater(t *testing.T) {
	tree, diags := parse(t, "struct S extends B {}\n")
	if len(diags) == 0 {
		t.Fatal("expected a diagnostic for `struct ... extends`")
	}
	if strings.Contains(diags[0].Msg, "unexpected") {
		t.Errorf("message should name the construct, got %q", diags[0].Msg)
	}
	var sd *ast.StructDecl
	ast.Inspect(tree, func(n ast.Node) bool {
		if d, ok := n.(*ast.StructDecl); ok {
			sd = d
		}
		return true
	})
	if sd == nil {
		t.Fatal("no StructDecl; the heritage clause should not have prevented the parse")
	}
	if sd.Extends == nil {
		t.Error("StructDecl.Extends is nil; the clause was dropped rather than carried")
	}

	// `kernel async function` — rejected by name, not by parse failure (§5.3).
	if _, diags := parse(t, "kernel async function f() {}\n"); len(diags) == 0 {
		t.Error("expected a diagnostic for an async accelerated function")
	}
}

// §6.1 site 1: instantiation expressions commit only when TypeArguments closes
// and the next token is in InstantiationFollowSet.
func TestInstantiationExpressions(t *testing.T) {
	tree, diags := parse(t, "let f = makeBox<boolean>;\nlet g = a < b > c;\n")
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostic: %v", diags[0].Msg)
	}
	var insts, binaries int
	ast.Inspect(tree, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.InstantiationExpr:
			insts++
		case *ast.BinaryExpr:
			binaries++
		}
		return true
	})
	if insts != 1 {
		t.Errorf("found %d instantiation expressions, want 1", insts)
	}
	if binaries != 2 {
		t.Errorf("found %d binary expressions, want 2 for `a < b > c`", binaries)
	}
}

// §5.4: ParenExpr is retained, because `(makeBox<boolean>)(true)` does not read
// the same without it.
func TestParenRetained(t *testing.T) {
	tree, _ := parse(t, "let x = (makeBox<boolean>)(true);\n")
	found := false
	ast.Inspect(tree, func(n ast.Node) bool {
		if _, ok := n.(*ast.ParenExpr); ok {
			found = true
		}
		return true
	})
	if !found {
		t.Error("ParenExpr was folded away")
	}
}

// §5.1 and grammar K: the only cover, and the only in-parser reinterpretation.
func TestCoverInitializedName(t *testing.T) {
	tree, diags := parse(t, "({ a = 1 } = obj);\n")
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostic: %v", diags[0].Msg)
	}
	found := false
	ast.Inspect(tree, func(n ast.Node) bool {
		if _, ok := n.(*ast.ObjectPattern); ok {
			found = true
		}
		return true
	})
	if !found {
		t.Error("the object literal was not reinterpreted as a pattern")
	}
	// Without the assignment, the cover form is illegal and must be reported.
	if _, diags := parse(t, "let o = { a = 1 };\n"); len(diags) == 0 {
		t.Error("expected a diagnostic for an unreinterpreted CoverInitializedName")
	}
}

// §6.1: depth is capped, and exceeding the cap is a diagnostic, never a hang.
func TestDepthCapped(t *testing.T) {
	src := "let x = " + strings.Repeat("(", 5000) + "1" + strings.Repeat(")", 5000) + ";\n"
	done := make(chan struct{})
	go func() {
		defer close(done)
		if _, diags := parse(t, src); len(diags) == 0 {
			t.Error("expected a depth-limit diagnostic")
		}
	}()
	<-done
}

// §6.4 and §7: every corpus file, truncated at every token boundary, must
// parse without panic, without hang, and yield a walkable tree.
func TestTruncationRecovery(t *testing.T) {
	corpus := []string{
		"struct Point { x: float64; y: float64; }\n",
		"kernel function saxpy(a: float64, x: Array<float64>) { }\n",
		"interface U<T extends object = {}> { readonly id?: string; }\n",
		"type M<T> = { [K in keyof T as `get${K}`]?: () => T[K] };\n",
		"export default class extends Base implements I { #n = 0; static { init(); } }\n",
		"import defer * as m from \"m\" with { type: \"json\" };\n",
		"let f = <T,>(x: T): T => x;\n",
		"declare module \"m\" { export function f(): void; }\n",
	}
	for _, src := range corpus {
		f := token.NewFile("t.vs", []byte(src))
		full, _ := ParseFile(f, 0)
		if full == nil {
			t.Fatalf("%q: nil tree", src)
		}
		for n := 1; n <= len(src); n++ {
			trunc := token.NewFile("t.vs", []byte(src[:n]))
			tree, _ := ParseFile(trunc, 0)
			if tree == nil {
				t.Fatalf("%q truncated at %d: nil tree", src, n)
			}
			ast.Inspect(tree, func(ast.Node) bool { return true })
		}
	}
}

// §8.6: two fragment entry points, because types have their own hierarchy.
func TestFragments(t *testing.T) {
	if x, diags := ParseExpr(token.NewFile("<eval>", []byte("a?.b(1)"))); x == nil || len(diags) != 0 {
		t.Errorf("ParseExpr failed: %v", diags)
	}
	if ty, diags := ParseTypeExpr(token.NewFile("<hover>", []byte("keyof T[]"))); ty == nil || len(diags) != 0 {
		t.Errorf("ParseTypeExpr failed: %v", diags)
	}
}

// §8.1: ImportsOnly stops the parser, not the scanner.
func TestImportsOnly(t *testing.T) {
	src := "import a from \"x\";\nexport { b } from \"y\";\nfunction huge() { /* ... */ }\n"
	f := token.NewFile("t.vs", []byte(src))
	tree, _ := ParseFile(f, ImportsOnly)
	if len(tree.Items) != 2 {
		t.Errorf("got %d prologue items, want 2 (the re-export is a dependency edge)", len(tree.Items))
	}
}

// §8.3: literals keep raw spelling; nothing is decoded here.
func TestLiteralsAreRaw(t *testing.T) {
	f := token.NewFile("t.vs", []byte("let n = 1_024;\n"))
	tree, _ := ParseFile(f, 0)
	var raw string
	ast.Inspect(tree, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.NUMBER {
			raw = string(f.Between(lit.Pos(), lit.End()))
		}
		return true
	})
	if raw != "1_024" {
		t.Errorf("literal text = %q, want %q", raw, "1_024")
	}
}