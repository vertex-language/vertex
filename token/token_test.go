package token

import (
	"testing"
	"unsafe"
)

// Token must stay pointer-free and small: the parser buffers whole files, and
// a []Token that the GC has to scan defeats the point (§3).
func TestTokenLayout(t *testing.T) {
	if got := unsafe.Sizeof(Token{}); got != 12 {
		t.Errorf("sizeof(Token) = %d, want 12", got)
	}
}

func TestRanges(t *testing.T) {
	for _, tc := range []struct {
		k                   Kind
		lit, op, reservedOK bool
	}{
		{IDENT, true, false, false},
		{TEMPLATE_TAIL, true, false, false},
		{LBRACE, false, true, false},
		{COALESCE_ASSIGN, false, true, false},
		{AWAIT, false, false, true},
		{YIELD, false, false, true},
		{EOF, false, false, false},
	} {
		if tc.k.IsLiteral() != tc.lit || tc.k.IsOperator() != tc.op || tc.k.IsReserved() != tc.reservedOK {
			t.Errorf("%s: ranges = (%v,%v,%v)", tc.k,
				tc.k.IsLiteral(), tc.k.IsOperator(), tc.k.IsReserved())
		}
	}
}

func TestJoinGT(t *testing.T) {
	for _, tc := range []struct {
		gts  int
		eq   bool
		want Kind
	}{
		{1, false, GT}, {1, true, GEQ},
		{2, false, SHR}, {2, true, SHR_ASSIGN},
		{3, false, USHR}, {3, true, USHR_ASSIGN},
		{4, false, INVALID},
	} {
		if got := JoinGT(tc.gts, tc.eq); got != tc.want {
			t.Errorf("JoinGT(%d,%v) = %s, want %s", tc.gts, tc.eq, got, tc.want)
		}
		if tc.want != INVALID && tc.want != GT && ScannerEmits(tc.want) {
			t.Errorf("%s must not be scanner-emittable", tc.want)
		}
		if tc.want.IsBinaryOperator() != (tc.want != INVALID) {
			t.Errorf("%s has no precedence; the climb in §6 needs one", tc.want)
		}
	}
}

func TestLookupIdent(t *testing.T) {
	for _, tc := range []struct {
		name string
		kind Kind
		ctx  Ctx
	}{
		{"if", IF, CtxNone},
		{"const", CONST, CtxNone},
		{"struct", IDENT, CtxStruct},
		{"kernel", IDENT, CtxKernel},
		{"graph", IDENT, CtxGraph},
		{"let", IDENT, CtxLet},     // StrictReservedWord: IDENT, not a Kind
		{"static", IDENT, CtxStatic},
		{"int32", IDENT, CtxNone},  // predeclared, not lexical (§4.6)
		{"foo", IDENT, CtxNone},
	} {
		k, c := LookupIdent([]byte(tc.name))
		if k != tc.kind || c != tc.ctx {
			t.Errorf("LookupIdent(%q) = (%s,%s), want (%s,%s)", tc.name, k, c, tc.kind, tc.ctx)
		}
	}
}

func TestCtxGroups(t *testing.T) {
	for _, c := range []Ctx{CtxDeclare, CtxStatic, CtxAbstract, CtxOverride, CtxReadonly, CtxPublic, CtxProtected, CtxPrivate} {
		if !c.IsClassElementModifier() {
			t.Errorf("%s: not a class element modifier", c)
		}
	}
	if CtxStruct.IsClassElementModifier() {
		t.Error("struct must not be a class element modifier")
	}
	for _, c := range []Ctx{CtxPublic, CtxProtected, CtxPrivate} {
		if !c.IsAccessibilityModifier() {
			t.Errorf("%s: not an accessibility modifier", c)
		}
	}
	if CtxStatic.IsAccessibilityModifier() {
		t.Error("static is not an accessibility modifier")
	}
}

func TestPosition(t *testing.T) {
	src := []byte("a\r\nbb\nccc\u2028d")
	f := NewFile("t.vs", src)

	for _, tc := range []struct {
		off       int
		line, col int
	}{
		{0, 1, 1},  // a
		{3, 2, 1},  // b
		{4, 2, 2},  // b
		{6, 3, 1},  // c
		{12, 4, 1}, // d, after the 3-byte U+2028
	} {
		got := f.Position(f.PosAt(tc.off))
		if got.Line != tc.line || got.Col != tc.col {
			t.Errorf("offset %d = %s, want %d:%d", tc.off, got, tc.line, tc.col)
		}
	}
	if p := f.Position(NoPos); p.IsValid() {
		t.Error("NoPos must render invalid")
	}
	// The first byte of a file must be distinguishable from NoPos.
	if f.PosAt(0) == NoPos {
		t.Error("PosAt(0) collides with NoPos")
	}
}

func TestBetween(t *testing.T) {
	f := NewFile("t.vs", []byte("let x = 1_024;"))
	if got := string(f.Between(f.PosAt(8), f.PosAt(13))); got != "1_024" {
		t.Errorf("Between = %q, want raw spelling %q", got, "1_024")
	}
	if f.Between(NoPos, f.PosAt(3)) != nil {
		t.Error("Between with NoPos must be nil")
	}
}