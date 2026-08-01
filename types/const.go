package types

import (
	"math"
	"math/big"
	"strconv"
)

// Value is a compile-time constant.
//
// A.1.5.1 ⊢ "an integer literal is untyped until it reaches a typed position,
// where it takes that position's type. A literal that does not fit the
// destination type is a compile error, never a silent truncation." That rule is
// the entire reason this type exists: a literal must be held at unbounded
// precision until its destination is known, or the check has already been lost
// by the time it can be made.
//
// It is also what A.3.1's ArrayLength, A.6.5's enum discriminants, and A.2's
// compile-time-evaluable top-level initializers are computed in.
type Value interface {
	Kind() ValueKind
	String() string
	valueNode()
}

type ValueKind int

const (
	UnknownVal ValueKind = iota // an erroneous or not-yet-computed constant
	BoolVal
	IntVal
	FloatVal
	CharVal
	StringVal
)

type (
	unknownVal struct{}
	boolVal    bool
	intVal     struct{ val *big.Int }
	floatVal   struct{ val *big.Rat }
	charVal    struct{ val rune }
	stringVal  struct{ val string }
)

func (unknownVal) Kind() ValueKind { return UnknownVal }
func (boolVal) Kind() ValueKind    { return BoolVal }
func (intVal) Kind() ValueKind     { return IntVal }
func (floatVal) Kind() ValueKind   { return FloatVal }
func (charVal) Kind() ValueKind    { return CharVal }
func (stringVal) Kind() ValueKind  { return StringVal }

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

func BoolVal_(v Value) (bool, bool) {
	b, ok := v.(boolVal)
	return bool(b), ok
}

func StringVal_(v Value) (string, bool) {
	s, ok := v.(stringVal)
	return s.val, ok
}

// ------------------------------------------------------------ operations

// Neg is unary minus. A.1.5.1 ⊢ "there is no literal syntax for a negative
// number. -1000 is unary minus applied to 1000, folded at compile time." This
// is that fold, and it must happen before the representability check — otherwise
// -128 would be rejected against int8 for being 128.
func Neg(v Value) Value {
	switch x := v.(type) {
	case intVal:
		return intVal{new(big.Int).Neg(x.val)}
	case floatVal:
		return floatVal{new(big.Rat).Neg(x.val)}
	}
	return Unknown
}

// BinaryOp folds a constant binary operation. It returns Unknown for a
// division by zero rather than trapping — the caller diagnoses, because a
// constant divide-by-zero is a compile error and not A.15's runtime trap.
func BinaryOp(x Value, op string, y Value) Value {
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
	}
	return nil, false
}

// Shift folds << and >>. A.4.5 puts the shift operators at the top of the
// cascade, above multiplication, but precedence does not reach constant
// folding — only the operator does.
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

// Compare folds a relational operation over two constants.
func Compare(x Value, op string, y Value) (bool, bool) {
	var c int
	switch a := x.(type) {
	case intVal:
		b, ok := y.(intVal)
		if !ok {
			return false, false
		}
		c = a.val.Cmp(b.val)
	case floatVal:
		b, ok := y.(floatVal)
		if !ok {
			return false, false
		}
		c = a.val.Cmp(b.val)
	case stringVal:
		b, ok := y.(stringVal)
		if !ok {
			return false, false
		}
		switch {
		case a.val < b.val:
			c = -1
		case a.val > b.val:
			c = 1
		}
	case charVal:
		b, ok := y.(charVal)
		if !ok {
			return false, false
		}
		switch {
		case a.val < b.val:
			c = -1
		case a.val > b.val:
			c = 1
		}
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
		return false, false
	default:
		return false, false
	}

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

// intRange holds the inclusive bounds of a sized integer kind.
type intRange struct{ lo, hi *big.Int }

var intRanges = map[BasicKind]intRange{}

func init() {
	signed := func(bits uint) intRange {
		hi := new(big.Int).Lsh(big.NewInt(1), bits-1)
		lo := new(big.Int).Neg(hi)
		return intRange{lo, hi.Sub(hi, big.NewInt(1))}
	}
	unsigned := func(bits uint) intRange {
		hi := new(big.Int).Lsh(big.NewInt(1), bits)
		return intRange{big.NewInt(0), hi.Sub(hi, big.NewInt(1))}
	}

	intRanges[Int8] = signed(8)
	intRanges[Int16] = signed(16)
	intRanges[Int32] = signed(32)
	intRanges[Int64] = signed(64)
	intRanges[Uint8] = unsigned(8)
	intRanges[Uint16] = unsigned(16)
	intRanges[Uint32] = unsigned(32)
	intRanges[Uint64] = unsigned(64)

	// int and uint follow IntSize. See sizes.go — the width is fixed across
	// every build tag, so these are settled here rather than per-target.
	intRanges[Int] = signed(intBits)
	intRanges[Uint] = unsigned(intBits)
}

// Representable reports whether v fits in b, and is the single enforcement
// point for A.1.5.1's no-silent-truncation rule.
//
// A float constant with a non-zero fractional part is not representable in an
// integer type at all — there is no rounding rule in Vertex, so `let x: int32 =
// 1.5` is an error rather than a truncation.
func Representable(v Value, b *Basic) bool {
	if b == nil || v == Unknown {
		return false
	}
	switch {
	case b.is(InfoInteger):
		i, ok := toInt(v)
		if !ok {
			return false
		}
		r, ok := intRanges[b.kind]
		if !ok {
			return false
		}
		return i.Cmp(r.lo) >= 0 && i.Cmp(r.hi) <= 0

	case b.is(InfoFloat):
		r, ok := toRat(v)
		if !ok {
			return false
		}
		f, _ := r.Float64()
		if math.IsInf(f, 0) {
			return false
		}
		if b.kind == Float32 {
			return !math.IsInf(float64(float32(f)), 0)
		}
		return true

	case b.is(InfoChar):
		// A.1.5.2 ⊢ a CharLiteral denotes exactly one Unicode scalar value. A
		// surrogate is a code point but not a scalar, so it is excluded here
		// for the same reason the scanner excludes it in an escape.
		i, ok := Int64Val(v)
		if !ok {
			return false
		}
		return i >= 0 && i <= 0x10FFFF && !(i >= 0xD800 && i < 0xE000)

	case b.is(InfoBoolean):
		return v.Kind() == BoolVal

	case b.is(InfoString):
		return v.Kind() == StringVal
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