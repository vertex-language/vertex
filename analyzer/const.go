package analyzer

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// This file folds the constant expressions §5.3 names: an ArrayLength, a
// ShapeList entry, an enum discriminant, and a top-level VarDecl initializer.
// It is deliberately not a general expression evaluator — that needs types, and
// belongs to whatever types expressions. What it does fold, it folds in
// types.Value, at unbounded precision, because §4.1 ⊢ a literal whose value does
// not fit its destination is a compile error rather than a wraparound, and the
// check is already lost if the value was narrowed first.
//
// Anything richer returns Unknown and the caller diagnoses. Rejecting rather
// than silently accepting is what makes widening this later unable to change a
// program that compiles today.

func (c *Checker) constValue(e ast.Expr) types.Value {
	switch x := stripParens(e).(type) {
	case *ast.BasicLit:
		return basicLitValue(x)

	case *ast.Ident:
		obj := c.lookup(x)
		cn, ok := obj.(*types.Const)
		if !ok {
			return types.Unknown
		}
		c.objDecl(cn)
		return cn.Val()

	case *ast.UnaryExpr:
		v := c.constValue(x.X)
		switch x.Op {
		case token.SUB:
			// grammar.md ⊢ there is no literal syntax for a negative number;
			// `-1000` is unary minus applied to `1000`. This fold is why a
			// negative enum discriminant is spellable at all, and why it must
			// run before any representability check.
			return types.Neg(v)
		case token.NOT:
			return types.Not(v)
		}
		// `~` has no constant fold: bitwise-NOT of an unbounded integer is not
		// defined, so the operand needs a width first.
		return types.Unknown

	case *ast.BinaryExpr:
		return c.constBinary(x)
	}
	return types.Unknown
}

func (c *Checker) constBinary(x *ast.BinaryExpr) types.Value {
	a, b := c.constValue(x.X), c.constValue(x.Y)
	if isUnknown(a) || isUnknown(b) {
		return types.Unknown
	}

	switch x.Op {
	case token.SHL, token.SHR:
		// §5.1 ⊢ the left operand's type is the result type; a count at or
		// beyond its width is a runtime trap and is not folded.
		n, ok := types.Int64Val(b)
		if !ok || n < 0 {
			return types.Unknown
		}
		return types.Shift(a, x.Op == token.SHL, uint(n))

	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		res, ok := types.Compare(a, x.Op.Spelling(), b)
		if !ok {
			return types.Unknown
		}
		return types.MakeBool(res)

	case token.LAND, token.LOR:
		p, okp := types.BoolVal(a)
		q, okq := types.BoolVal(b)
		if !okp || !okq {
			return types.Unknown
		}
		if x.Op == token.LAND {
			return types.MakeBool(p && q)
		}
		return types.MakeBool(p || q)

	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM,
		token.OR, token.AND, token.XOR:
		// The wrapping operators are absent on purpose: they wrap at a width,
		// and an untyped constant has none.
		return types.BinaryOp(a, x.Op.Spelling(), b)
	}
	return types.Unknown
}

func isUnknown(v types.Value) bool { return v == nil || v.Kind() == types.UnknownKind }

// ---------------------------------------------------------------- literals

// basicLitValue decodes a literal's raw source spelling. The scanner kept the
// spelling exactly as written — separators intact, escapes unresolved — because
// a formatter needs the original; decoding is this side's job.
func basicLitValue(x *ast.BasicLit) types.Value {
	switch x.Kind {
	case token.TRUE:
		return types.MakeBool(true)
	case token.FALSE:
		return types.MakeBool(false)

	case token.INT:
		n, ok := parseIntLit(x.Value)
		if !ok {
			return types.Unknown
		}
		return types.MakeInt(n)

	case token.FLOAT:
		f, ok := parseFloatLit(x.Value)
		if !ok {
			return types.Unknown
		}
		r := new(big.Rat).SetFloat64(f)
		if r == nil { // Inf or NaN
			return types.Unknown
		}
		return types.MakeFloat(r)

	case token.CHAR:
		r, ok := decodeChar(x.Value)
		if !ok {
			return types.Unknown
		}
		return types.MakeChar(r)

	case token.STRING:
		s, ok := decodeString(x.Value)
		if !ok {
			return types.Unknown
		}
		return types.MakeString(s)
	}
	// §10 ⊢ `nil` belongs to typed_ptr T and to nothing else, so it has no
	// constant value of its own.
	return types.Unknown
}

// parseIntLit decodes an int_lit, prefix and `_` separators included.
//
// Base 0 parsing is deliberately not used: grammar.md has no prefix-free octal
// form, so `0600` is the decimal integer 600 and a host convention that reads a
// leading zero as octal would be wrong here.
func parseIntLit(s string) (*big.Int, bool) {
	base := 10
	if len(s) > 2 && s[0] == '0' {
		switch s[1] {
		case 'b', 'B':
			base, s = 2, s[2:]
		case 'o', 'O':
			base, s = 8, s[2:]
		case 'x', 'X':
			base, s = 16, s[2:]
		}
	}
	s = strings.ReplaceAll(s, "_", "")
	if s == "" {
		return nil, false
	}
	n, ok := new(big.Int).SetString(s, base)
	return n, ok
}

func parseFloatLit(s string) (float64, bool) {
	f, err := strconv.ParseFloat(strings.ReplaceAll(s, "_", ""), 64)
	return f, err == nil
}

// decodeString decodes a string_lit. A raw string recognizes no escape
// sequence, and every line terminator it spans is part of its value.
func decodeString(raw string) (string, bool) {
	if len(raw) < 2 {
		return "", false
	}
	if raw[0] == '`' {
		if raw[len(raw)-1] != '`' {
			return "", false
		}
		return raw[1 : len(raw)-1], true
	}
	if raw[0] != '"' || raw[len(raw)-1] != '"' {
		return "", false
	}

	var b strings.Builder
	body := raw[1 : len(raw)-1]
	for i := 0; i < len(body); {
		if body[i] != '\\' {
			b.WriteByte(body[i])
			i++
			continue
		}
		r, w, ok := decodeEscape(body[i:])
		if !ok {
			return "", false
		}
		b.WriteRune(r)
		i += w
	}
	return b.String(), true
}

// decodeChar decodes a char_lit, which denotes exactly one Unicode scalar
// value. That it holds exactly one is the scanner's check; this only reads it.
func decodeChar(raw string) (rune, bool) {
	if len(raw) < 3 || raw[0] != '\'' || raw[len(raw)-1] != '\'' {
		return 0, false
	}
	body := raw[1 : len(raw)-1]
	if body[0] == '\\' {
		r, w, ok := decodeEscape(body)
		if !ok || w != len(body) {
			return 0, false
		}
		return r, true
	}
	for _, r := range body {
		return r, true
	}
	return 0, false
}

// decodeEscape decodes one escape sequence, returning its value and the number
// of bytes consumed including the backslash.
func decodeEscape(s string) (rune, int, bool) {
	if len(s) < 2 {
		return 0, 0, false
	}
	switch s[1] {
	case '\'':
		return '\'', 2, true
	case '"':
		return '"', 2, true
	case '\\':
		return '\\', 2, true
	case 'n':
		return '\n', 2, true
	case 'r':
		return '\r', 2, true
	case 't':
		return '\t', 2, true
	case 'v':
		return '\v', 2, true
	case 'b':
		return '\b', 2, true
	case 'f':
		return '\f', 2, true
	case '0':
		return 0, 2, true

	case 'x':
		if len(s) < 4 {
			return 0, 0, false
		}
		n, err := strconv.ParseUint(s[2:4], 16, 32)
		if err != nil {
			return 0, 0, false
		}
		return rune(n), 4, true

	case 'u':
		if len(s) < 4 || s[2] != '{' {
			return 0, 0, false
		}
		end := strings.IndexByte(s, '}')
		if end < 0 {
			return 0, 0, false
		}
		n, err := strconv.ParseUint(s[3:end], 16, 32)
		if err != nil {
			return 0, 0, false
		}
		return rune(n), end + 1, true
	}
	return 0, 0, false
}

// -------------------------------------------------------------- positions

// constInt folds an expression that must be a non-negative integer constant:
// an ArrayLength, a ShapeList entry, a vector lane count.
func (c *Checker) constInt(e ast.Expr) (int64, bool) {
	v := c.constValue(e)
	n, ok := types.Int64Val(v)
	if !ok {
		c.errorExpr(e, diag.ArrayLenNotConst)
		return 0, false
	}
	if n < 0 {
		c.errorExpr(e, diag.ArrayLenNegative)
		return 0, false
	}
	return n, true
}