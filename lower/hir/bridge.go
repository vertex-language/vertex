package hir

import "github.com/vertex-language/vertex/types"

// bridge.go is the only file in this package that knows types' concrete
// API. Everything else asks these helpers a classified question and
// switches on hir's own vocabulary. If a types signature changes, this file
// changes and nothing else does.

// kind is hir's classification of a checked Vertex type. It exists so that
// no other file needs a type switch over types' concrete kinds.
type kind uint8

const (
	kInvalid kind = iota
	kUnit
	kBool
	kChar
	kInt
	kFloat
	kString
	kSlice
	kArray
	kMap
	kChan
	kPointer
	kUnique
	kShared
	kWeak
	kFunc
	kTuple
	kStruct
	kEnum
	kAbstract
	kTensor
)

// owning reports whether a bare copy of this kind is more than a bit copy —
// the left column of overview §3's ownership table.
func (k kind) owning() bool {
	switch k {
	case kString, kSlice, kMap, kChan, kUnique, kShared, kWeak, kStruct, kEnum:
		return true
	}
	return false
}

func (l *lowerer) classify(t types.Type) kind {
	t = l.subst(t)
	switch u := t.(type) {
	case *types.Basic:
		switch {
		case l.isInvalid(u):
			return kInvalid
		case u.Kind() == types.Bool || u.Kind() == types.UntypedBool:
			return kBool
		case u.Kind() == types.Char || u.Kind() == types.UntypedChar:
			return kChar
		case u.Kind() == types.String || u.Kind() == types.UntypedString:
			return kString
		case l.isFloatKind(u):
			return kFloat
		}
		return kInt
	case *types.Slice:
		return kSlice
	case *types.Array:
		return kArray
	case *types.Map:
		return kMap
	case *types.Chan:
		return kChan
	case *types.Pointer:
		return kPointer
	case *types.Signature:
		return kFunc
	case *types.Tensor:
		return kTensor
	case *types.Abstract:
		return kAbstract
	case *types.Tuple:
		if u.Len() == 0 {
			return kUnit
		}
		return kTuple
	case *types.Ownership:
		switch u.Kind() {
		case types.Unique:
			return kUnique
		case types.Shared:
			return kShared
		case types.Weak:
			return kWeak
		}
		// mut/var are parameter conventions, never a stored shape: the
		// underlying type is what a value of this position holds.
		return l.classify(u.Elem())
	case *types.Named:
		return l.classify(u.Underlying())
	case *types.Struct:
		return kStruct
	case *types.Enum:
		return kEnum
	case *types.TypeParam:
		// Post-monomorphization this cannot happen; reaching it means the
		// substitution map was incomplete, which is a compiler bug.
		l.errorf(0, "internal: unsubstituted type parameter %s reached hir", types.TypeString(t))
		return kInvalid
	}
	return kInvalid
}

func (l *lowerer) isInvalid(b *types.Basic) bool { return b.Kind() == types.Invalid }

func (l *lowerer) isFloatKind(b *types.Basic) bool {
	switch b.Kind() {
	case types.Float32, types.Float64, types.UntypedFloat:
		return true
	}
	return false
}

// isSigned decides sdiv-vs-udiv, slt-vs-ult, and sext-vs-zext. Every
// opcode-selection site asks here rather than re-deriving it.
func (l *lowerer) isSigned(t types.Type) bool {
	b, ok := l.underlying(t).(*types.Basic)
	if !ok {
		return false
	}
	switch b.Kind() {
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		return false
	}
	return true
}

func (l *lowerer) underlying(t types.Type) types.Type {
	if t == nil {
		return nil
	}
	return l.subst(t).Underlying()
}

func (l *lowerer) intBits(t types.Type) int {
	if n := l.conf.Sizes.Sizeof(l.subst(t)); n > 0 {
		return int(n * 8)
	}
	return 64
}

func (l *lowerer) floatBits(t types.Type) int {
	if n := l.conf.Sizes.Sizeof(l.subst(t)); n > 0 {
		return int(n * 8)
	}
	return 64
}

func (l *lowerer) elem(t types.Type) types.Type {
	switch u := l.underlying(t).(type) {
	case *types.Slice:
		return u.Elem()
	case *types.Array:
		return u.Elem()
	case *types.Pointer:
		return u.Elem()
	case *types.Chan:
		return u.Elem()
	case *types.Ownership:
		return u.Elem()
	case *types.Map:
		return u.Elem()
	}
	return nil
}

func (l *lowerer) arrayParts(t types.Type) (types.Type, int64) {
	a, ok := l.underlying(t).(*types.Array)
	if !ok {
		return nil, 0
	}
	return a.Elem(), a.Len()
}

type fieldInfo struct {
	Name string
	Type types.Type
}

func (l *lowerer) structParts(t types.Type) (string, []fieldInfo) {
	name := "anon"
	if n, ok := l.subst(t).(*types.Named); ok {
		name = n.Obj().Name()
	}
	s, ok := l.underlying(t).(*types.Struct)
	if !ok {
		return name, nil
	}
	out := make([]fieldInfo, 0, s.NumFields())
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		out = append(out, fieldInfo{Name: f.Name(), Type: f.Type()})
	}
	return name, out
}

type tupleElem struct {
	Name string
	Type types.Type
}

func (l *lowerer) tupleElems(t types.Type) []tupleElem {
	tp, ok := l.underlying(t).(*types.Tuple)
	if !ok {
		return nil
	}
	out := make([]tupleElem, 0, tp.Len())
	for i := 0; i < tp.Len(); i++ {
		v := tp.At(i)
		out = append(out, tupleElem{Name: v.Name(), Type: v.Type()})
	}
	return out
}

type variantInfo struct {
	Name  string
	Tag   int64
	Types []types.Type
}

func (l *lowerer) enumParts(t types.Type) (disc types.Type, unitOnly bool, variants []variantInfo) {
	e, ok := l.underlying(t).(*types.Enum)
	if !ok {
		return nil, true, nil
	}
	for i := 0; i < e.NumVariants(); i++ {
		v := e.Variant(i)
		vi := variantInfo{Name: v.Name(), Tag: v.Tag()}
		for j := 0; j < v.NumTypes(); j++ {
			vi.Types = append(vi.Types, v.Type(j))
		}
		variants = append(variants, vi)
	}
	return e.Discriminant(), e.UnitOnly(), variants
}