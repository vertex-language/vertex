package types

import (
	"math"
	"math/big"
	"strconv"
)

// Value is a compile-time constant.
//
// §4.1 ⊢ "a literal has no type until it lands... A literal whose value does not
// fit its destination is a compile error, not a wraparound." That rule is the
// entire reason this type exists: a literal must be held at unbounded precision
// until its destination is known, or the check has already been lost by the
// time it can be made.
//
// It is also the arithmetic §5.3's constant expressions are evaluated in: an
// ArrayLength, an enum discriminant, a top-level VarDecl initializer, a switch
// case pattern, and `new`'s align: argument.
type Value interface {
	Kind() ValueKind
	String() string
	valueNode()
}

type ValueKind int

const (
	UnknownKind ValueKind = iota // erroneous, or not yet computed
	BoolKind
	IntKind
	FloatKind
	CharKind
	StringKind
)

type (
	unknownVal struct{}
	boolVal    bool
	intVal     struct{ val *big.Int }
	floatVal   struct{ val *big.Rat }
	charVal    struct{ val rune }
	stringVal  struct{ val string }
)

func (unknownVal) Kind() ValueKind { return UnknownKind }
func (boolVal) Kind() ValueKind    { return BoolKind }
func (intVal) Kind() ValueKind     { return IntKind }
func (floatVal) Kind() ValueKind   { return FloatKind }
func (charVal) Kind() ValueKind    { return CharKind }
func (stringVal) Kind() ValueKind  { return StringKind }

func (unknownVal) valueNode() {}
func (boolVal) valueNode()    {}
func (intVal) valueNode()     {}
func (floatVal) valueNode()   {}
func (charVal) valueNode()    {}
func (stringVal) valueNode()  {}

func (unknownVal) String() string  { return "unknown" }
func (v boolVal) String() string   { return strconv.FormatBool(bool(v)) }
func (v intVal) String() string    { return v.val.String() }
func (v charVal) String() string   { return strconv.QuoteRune(v.val) }
func (v stringVal) String() string { return strconv.Quote(v.val) }

func (v floatVal) String() string {
	f, _ := v.val.Float64()
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// Unknown is the constant of a subexpression that already failed. Every fold
// below propagates it rather than reporting again, so one bad literal yields
// one diagnostic.
var Unknown Value = unknownVal{}

func MakeBool(b bool) Value      { return boolVal(b) }
func MakeInt(x *big.Int) Value   { return intVal{new(big.Int).Set(x)} }
func MakeInt64(x int64) Value    { return intVal{big.NewInt(x)} }
func MakeFloat(x *big.Rat) Value { return floatVal{new(big.Rat).Set(x)} }
func MakeChar(r rune) Value      { return charVal{r} }
func MakeString(s string) Value  { return stringVal{s} }

// ------------------------------------------------------------- accessors

func Int64Val(v Value) (int64, bool) {
	switch x := v.(type) {
	case intVal:
		if x.val.IsInt64() {
			return x.val.Int64(), true
		}
	case charVal:
		return int64(x.val), true
	}
	return 0, false
}

func BoolVal(v Value) (bool, bool) {
	b, ok := v.(boolVal)
	return bool(b), ok
}

func StringVal(v Value) (string, bool) {
	s, ok := v.(stringVal)
	return s.val, ok
}

func FloatVal(v Value) (float64, bool) {
	r, ok := toRat(v)
	if !ok {
		return 0, false
	}
	f, _ := r.Float64()
	return f, true
}

// ------------------------------------------------------------ operations

// Neg is unary minus. grammar.md ⊢ "there is no literal syntax for a negative
// number. `-1000` is unary minus applied to `1000`." This is that fold, and it
// must happen before the representability check — otherwise -128 would be
// rejected against int8 for being 128.
//
// §5.3's three literal-token positions — a ShapeList, a VectorType's lane
// count, and a TupleIndex — are exactly the ones that must not accept this
// fold's output, which is why they are specified against the token rather than
// against a constant expression.
func Neg(v Value) Value {
	switch x := v.(type) {
	case intVal:
		return intVal{new(big.Int).Neg(x.val)}
	case floatVal:
		return floatVal{new(big.Rat).Neg(x.val)}
	}
	return Unknown
}

// Not is unary `!`. §5.1 ⊢ `!` takes "`bool` only".
//
// `~` has no fold here on purpose: bitwise-NOT of an unbounded integer is not
// defined, so the checker converts the operand to its destination type first
// and the result is folded at that width or not at all.
func Not(v Value) Value {
	if b, ok := v.(boolVal); ok {
		return boolVal(!b)
	}
	return Unknown
}

// BinaryOp folds a constant binary operation.
//
// A division or remainder by zero returns Unknown rather than trapping: §5.5's
// trap is a runtime tier, and a constant divide-by-zero is a compile error the
// caller diagnoses.
//
// The wrapping operators `&+ &- &*` are absent. They wrap at a width, and an
// untyped constant has none — the checker types both operands first and folds
// at that width.
func BinaryOp(x Value, op string, y Value) Value {
	if xs, ok := x.(stringVal); ok {
		// §5.1 ⊢ `+` on two `string` is concatenation.
		ys, ok := y.(stringVal)
		if !ok || op != "+" {
			return Unknown
		}
		return stringVal{xs.val + ys.val}
	}

	xi, xok := x.(intVal)
	yi, yok := y.(intVal)
	if !xok || !yok {
		return binaryOpFloat(x, op, y)
	}

	z := new(big.Int)
	switch op {
	case "+":
		z.Add(xi.val, yi.val)
	case "-":
		z.Sub(xi.val, yi.val)
	case "*":
		z.Mul(xi.val, yi.val)
	case "/":
		if yi.val.Sign() == 0 {
			return Unknown
		}
		z.Quo(xi.val, yi.val)
	case "%":
		if yi.val.Sign() == 0 {
			return Unknown
		}
		z.Rem(xi.val, yi.val)
	case "|":
		z.Or(xi.val, yi.val)
	case "&":
		z.And(xi.val, yi.val)
	case "^":
		z.Xor(xi.val, yi.val)
	default:
		return Unknown
	}
	return intVal{z}
}

func binaryOpFloat(x Value, op string, y Value) Value {
	xf, xok := toRat(x)
	yf, yok := toRat(y)
	if !xok || !yok {
		return Unknown
	}
	z := new(big.Rat)
	switch op {
	case "+":
		z.Add(xf, yf)
	case "-":
		z.Sub(xf, yf)
	case "*":
		z.Mul(xf, yf)
	case "/":
		if yf.Sign() == 0 {
			return Unknown
		}
		z.Quo(xf, yf)
	default:
		return Unknown
	}
	return floatVal{z}
}

func toRat(v Value) (*big.Rat, bool) {
	switch x := v.(type) {
	case intVal:
		return new(big.Rat).SetInt(x.val), true
	case floatVal:
		return x.val, true
	case charVal:
		return new(big.Rat).SetInt64(int64(x.val)), true
	}
	return nil, false
}

// Shift folds `<<` and `>>`. §5.1 ⊢ the left operand's type is the result type;
// a count at or beyond its width is §5.5's trap and is not folded here.
func Shift(x Value, left bool, count uint) Value {
	xi, ok := x.(intVal)
	if !ok {
		return Unknown
	}
	z := new(big.Int)
	if left {
		z.Lsh(xi.val, count)
	} else {
		z.Rsh(xi.val, count)
	}
	return intVal{z}
}

// Compare folds a relational operation. The second result reports whether the
// operands were comparable as constants at all; §3.5 and §5.1 decide whether
// the operator was legal on their type in the first place.
func Compare(x Value, op string, y Value) (bool, bool) {
	if xr, ok := toRat(x); ok {
		yr, ok := toRat(y)
		if !ok {
			return false, false
		}
		return cmpResult(xr.Cmp(yr), op)
	}

	switch a := x.(type) {
	case stringVal:
		b, ok := y.(stringVal)
		if !ok {
			return false, false
		}
		switch {
		case a.val < b.val:
			return cmpResult(-1, op)
		case a.val > b.val:
			return cmpResult(1, op)
		}
		return cmpResult(0, op)

	case boolVal:
		b, ok := y.(boolVal)
		if !ok {
			return false, false
		}
		switch op {
		case "==":
			return a == b, true
		case "!=":
			return a != b, true
		}
	}
	return false, false
}

func cmpResult(c int, op string) (bool, bool) {
	switch op {
	case "==":
		return c == 0, true
	case "!=":
		return c != 0, true
	case "<":
		return c < 0, true
	case "<=":
		return c <= 0, true
	case ">":
		return c > 0, true
	case ">=":
		return c >= 0, true
	}
	return false, false
}

// ------------------------------------------------------ representability

// Representable reports whether v fits in b, and is the single enforcement
// point for §4.1's no-wraparound rule.
//
// It is a method on Sizes rather than a free function because §2.3 makes `int`
// and `uint` "the target's pointer width" — so whether a literal fits is a
// question about the target, and a caller must have one in hand to ask it. That
// is also why a file must not compile under one build tag and fail under
// another for a reason invisible in the source: the width is the tag's, and the
// tag is the file's.
//
// A float constant with a non-zero fractional part is not representable in an
// integer type at all. §4.2 gives `float → integer` only as an `as` conversion
// that "truncates toward zero"; there is no implicit rounding, so
// `let x: int32 = 1.5` is an error rather than a truncation.
func (s *Sizes) Representable(v Value, b *Basic) bool {
	if s == nil || b == nil || v == Unknown {
		return false
	}
	switch {
	case b.is(InfoInteger):
		i, ok := toInt(v)
		if !ok {
			return false
		}
		if b.is(InfoUntyped) {
			return true
		}
		lo, hi, ok := s.intRange(b.kind)
		if !ok {
			return false
		}
		return i.Cmp(lo) >= 0 && i.Cmp(hi) <= 0

	case b.is(InfoFloat):
		r, ok := toRat(v)
		if !ok {
			return false
		}
		f, _ := r.Float64()
		if math.IsInf(f, 0) {
			return false
		}
		if b.kind == Float32 || b.is(InfoTensorElem) {
			// The narrow float formats are checked against float32's range;
			// the sources fix no bf16/fp8 exponent range, so this is the
			// conservative bound rather than a stated one.
			return !math.IsInf(float64(float32(f)), 0)
		}
		return true

	case b.is(InfoChar):
		// grammar.md ⊢ a char_lit "denotes exactly one Unicode scalar value",
		// and a unicode_escape "must denote a Unicode scalar value: not above
		// 10FFFF, and not a surrogate". A surrogate is a code point but not a
		// scalar, so it is excluded here for the same reason the scanner
		// excludes it in an escape.
		i, ok := Int64Val(v)
		if !ok {
			return false
		}
		return i >= 0 && i <= 0x10FFFF && !(i >= 0xD800 && i < 0xE000)

	case b.is(InfoBoolean):
		return v.Kind() == BoolKind

	case b.is(InfoString):
		return v.Kind() == StringKind
	}
	return false
}

// toInt converts an exact-integer constant to a big.Int. A float constant
// converts only when its denominator is 1 — see Representable.
func toInt(v Value) (*big.Int, bool) {
	switch x := v.(type) {
	case intVal:
		return x.val, true
	case charVal:
		return big.NewInt(int64(x.val)), true
	case floatVal:
		if x.val.IsInt() {
			return new(big.Int).Set(x.val.Num()), true
		}
	}
	return nil, false
}