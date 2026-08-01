package hir

import "github.com/vertex-language/vertex/types"

// subst.go applies a monomorphization substitution to a checked type.
//
// It lives here rather than in types because the substitution map is
// hir's: types records what the analyzer saw while checking a generic body
// once, generically, and mono.go composes the concrete arguments onto it as
// it descends. Nothing upstream of hir ever needs to rewrite a type.
//
// Two properties matter to callers. A type mentioning no substituted
// parameter is returned unchanged, pointer-identical, so the typeLowerer's
// cache keys stay stable. And a self-referential Named — A.2's
// order-independence lets `struct Node { next: typed_ptr Node }` reach
// itself — terminates, because the rewritten Named is bound before its
// underlying is walked, the same two-step named.go's SetUnderlying exists for.
func substitute(t types.Type, m map[*types.TypeParam]types.Type) types.Type {
	if t == nil || len(m) == 0 {
		return t
	}
	s := &substituter{m: m, seen: map[*types.Named]*types.Named{}}
	return s.typ(t)
}

type substituter struct {
	m    map[*types.TypeParam]types.Type
	seen map[*types.Named]*types.Named
}

func (s *substituter) typ(t types.Type) types.Type {
	switch x := t.(type) {
	case nil:
		return nil

	case *types.TypeParam:
		if r, ok := s.m[x]; ok {
			return r
		}
		return x

	case *types.Basic, *types.Abstract:
		return t

	case *types.Ownership:
		if e := s.typ(x.Elem()); e != x.Elem() {
			return types.NewOwnership(x.Kind(), e)
		}

	case *types.Array:
		if e := s.typ(x.Elem()); e != x.Elem() {
			return types.NewArray(e, x.Len())
		}

	case *types.Slice:
		if e := s.typ(x.Elem()); e != x.Elem() {
			return types.NewSlice(e)
		}

	case *types.Chan:
		if e := s.typ(x.Elem()); e != x.Elem() {
			return types.NewChan(e)
		}

	case *types.Pointer:
		if e := s.typ(x.Elem()); e != x.Elem() {
			return types.NewPointer(e)
		}

	case *types.Tensor:
		if e := s.typ(x.Elem()); e != x.Elem() {
			return types.NewTensor(e, x.Shape())
		}

	case *types.Map:
		k, v := s.typ(x.Key()), s.typ(x.Elem())
		if k != x.Key() || v != x.Elem() {
			return types.NewMap(k, v)
		}

	case *types.Tuple:
		return s.tuple(x)

	case *types.Signature:
		recv := x.Recv()
		if recv != nil {
			recv = s.vr(recv)
		}
		p, r := s.tuple(x.Params()), s.tuple(x.Results())
		if recv != x.Recv() || p != x.Params() || r != x.Results() {
			return types.NewSignature(recv, p, r, x.Variadic(), x.Marker())
		}

	case *types.Struct:
		changed := false
		fields := make([]*types.Field, 0, x.NumFields())
		for i := 0; i < x.NumFields(); i++ {
			f := x.Field(i)
			ft := s.typ(f.Type)
			if ft != f.Type {
				changed = true
			}
			fields = append(fields, &types.Field{Name: f.Name, Type: ft, HasDefault: f.HasDefault})
		}
		if changed {
			return types.NewStruct(fields, x.Class())
		}

	case *types.Enum:
		changed := false
		vs := make([]*types.Variant, 0, x.NumVariants())
		for i := 0; i < x.NumVariants(); i++ {
			v := x.Variant(i)
			var payload []types.Type
			for _, p := range v.Payload {
				np := s.typ(p)
				if np != p {
					changed = true
				}
				payload = append(payload, np)
			}
			vs = append(vs, &types.Variant{Name: v.Name, Payload: payload, Value: v.Value})
		}
		if changed {
			return types.NewEnum(vs, x.Discriminant())
		}

	case *types.Named:
		return s.named(x)
	}
	return t
}

func (s *substituter) tuple(t *types.Tuple) *types.Tuple {
	if t == nil {
		return nil
	}
	changed := false
	vars := make([]*types.Var, 0, t.Len())
	for i := 0; i < t.Len(); i++ {
		v := t.At(i)
		nv := s.vr(v)
		if nv != v {
			changed = true
		}
		vars = append(vars, nv)
	}
	if !changed {
		return t
	}
	return types.NewTuple(vars...)
}

// vr rewrites one Var, preserving its Mode. A.3.2's mut/var are not part of
// the Type, so a substitution that dropped them would silently turn a `mut`
// parameter into a by-value one — and Identical compares modes.
func (s *substituter) vr(v *types.Var) *types.Var {
	nt := s.typ(v.Type())
	if nt == v.Type() {
		return v
	}
	return types.NewParam(v.Pos(), v.Pkg(), v.Name(), nt, v.Mode())
}

// named rewrites a defined type. The instantiation keeps the declaring
// object, so Identical still compares by object plus type arguments and
// owningModule still finds the declaring package — only the shape moves.
func (s *substituter) named(x *types.Named) types.Type {
	if got, ok := s.seen[x]; ok {
		return got
	}

	args := x.TypeArgs()
	changed := false
	var newArgs []types.Type
	if len(args) > 0 {
		newArgs = make([]types.Type, len(args))
		for i, a := range args {
			newArgs[i] = s.typ(a)
			if newArgs[i] != a {
				changed = true
			}
		}
	}

	// Bind before descending: a field naming its own enclosing type reaches
	// this entry instead of recursing forever.
	out := types.NewNamed(x.Obj(), nil)
	s.seen[x] = out

	u := s.typ(x.Underlying())
	if !changed && u == x.Underlying() {
		delete(s.seen, x)
		return x
	}

	out.SetUnderlying(u)
	if newArgs != nil {
		out.SetTypeArgs(newArgs)
	}
	out.SetTypeParams(x.TypeParams())
	for i := 0; i < x.NumMethods(); i++ {
		out.AddMethod(x.Method(i))
	}
	return out
}