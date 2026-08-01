package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

func (b *funcBuilder) stmtList(list []ast.Stmt) {
	for _, s := range list {
		b.stmt(s)
	}
}

func (b *funcBuilder) stmt(s ast.Stmt) {
	switch s := s.(type) {
	case *ast.BlockStmt:
		b.push(scopeBlock)
		b.stmtList(s.List)
		b.pop()

	case *ast.DeclStmt:
		if vd, ok := s.Decl.(*ast.VarDecl); ok {
			b.varDecl(vd)
		}

	case *ast.ExprStmt:
		b.expr(s.X)

	case *ast.AssignStmt:
		b.assign(s)

	case *ast.IfStmt:
		b.ifStmt(s)

	case *ast.WhileStmt:
		b.whileStmt(s)

	case *ast.ForStmt:
		b.forStmt(s)

	case *ast.SwitchStmt:
		b.switchStmt(s)

	case *ast.ReturnStmt:
		b.returnStmt(s)

	case *ast.DeferStmt:
		b.deferStmt(s)

	case *ast.BranchStmt:
		// There are no loop labels (A.5.9), so a multi-level exit was
		// already written as a flag or an extracted function.
		switch s.Tok.String() {
		case "break":
			b.unwind(scopeLoop)
			b.seq.add(&Break{})
		case "continue":
			b.unwind(scopeLoop)
			b.seq.add(&Continue{})
		case "fallthrough":
			// Handled structurally by switchStmt.
		}

	case *ast.SelectStmt:
		b.l.todo(s.Pos(), "select lowering — every case is a channel receive (A.10.2), so this becomes a builtins/chan select over the case set")

	case *ast.BadStmt:
		// A diagnostic was already reported at parse time.

	default:
		b.l.todo(s.Pos(), "statement %T", s)
	}
}

// varDecl lowers `let`/`var`. The two differ in representation, not in
// storage class: a let is a named value, a var is a slot.
func (b *funcBuilder) varDecl(d *ast.VarDecl) {
	isVar := d.Kw.String() == "var"

	// A multi-binding declaration with a single initializer is a tuple
	// destructure; with a matching count it is parallel declaration.
	if len(d.Bindings) > 1 && len(d.Values) == 1 {
		b.destructure(d, isVar)
		return
	}

	for i, bind := range d.Bindings {
		obj := b.info().Defs[bind.Name]
		if obj == nil || bind.Name.IsBlank() {
			if i < len(d.Values) {
				b.expr(d.Values[i]) // still evaluated for effect
			}
			continue
		}
		src := obj.Type()
		ht := b.l.hirType(src)

		var init Value
		if i < len(d.Values) {
			init = b.owningExpr(d.Values[i], src)
		} else {
			// Bare `var x: T` yields the type's zero value (A.5.1).
			init = b.zeroValue(bind.Pos(), ht)
		}

		bd := &binding{Type: ht, Src: src, Owning: b.l.classify(src).owning()}
		if isVar || IsAggregate(ht) {
			slot := b.alloca(bind.Pos(), ht)
			b.store(bind.Pos(), ht, slot, init)
			bd.Value, bd.Slot = slot, true
		} else {
			bd.Value = init
		}
		b.declareBinding(obj, bd)
	}
}

func (b *funcBuilder) destructure(d *ast.VarDecl, isVar bool) {
	tupleType := b.info().TypeOf(d.Values[0])
	v := b.expr(d.Values[0])
	st, ok := b.l.hirType(tupleType).(StructType)
	if !ok {
		b.l.errorf(d.Pos(), "internal: destructuring a non-tuple %s", types.TypeString(tupleType))
		return
	}
	for i, bind := range d.Bindings {
		if bind.Name.IsBlank() || i >= len(st.Def.Fields) {
			continue
		}
		obj := b.info().Defs[bind.Name]
		if obj == nil {
			continue
		}
		f := st.Def.Fields[i]
		val := b.load(bind.Pos(), f.Type, b.fieldPtr(bind.Pos(), st.Def, v, f.Name))
		bd := &binding{Type: f.Type, Src: obj.Type(), Owning: b.l.classify(obj.Type()).owning()}
		if isVar || IsAggregate(f.Type) {
			slot := b.alloca(bind.Pos(), f.Type)
			b.store(bind.Pos(), f.Type, slot, val)
			bd.Value, bd.Slot = slot, true
		} else {
			bd.Value = val
		}
		b.declareBinding(obj, bd)
	}
}

func (b *funcBuilder) assign(s *ast.AssignStmt) {
	if s.Op.IsCompoundAssign() {
		t := b.info().TypeOf(s.Targets[0])
		lhs := b.addressOf(s.Targets[0])
		cur := b.load(s.Pos(), b.l.hirType(t), lhs)
		rhs := b.expr(s.Values[0])
		op := b.l.binaryOpFor(compoundBase(s.Op), t)
		res := b.op(s.OpPos, op, b.l.hirType(t), cur, rhs)
		b.store(s.Pos(), b.l.hirType(t), lhs, res)
		return
	}
	for i, tgt := range s.Targets {
		if i >= len(s.Values) {
			break
		}
		if id, ok := tgt.(*ast.Ident); ok && id.IsBlank() {
			b.expr(s.Values[i])
			continue
		}
		dstType := b.info().TypeOf(tgt)
		val := b.owningExpr(s.Values[i], dstType)
		ht := b.l.hirType(dstType)

		// Assigning over a live owning binding tears down what was there.
		if bd := b.bindingFor(tgt); bd != nil && bd.Owning && !bd.Dead {
			b.l.own.emitDeinit(b, s.Pos(), bd)
		}
		b.store(s.Pos(), ht, b.addressOf(tgt), val)
	}
}

func (b *funcBuilder) ifStmt(s *ast.IfStmt) {
	cond := b.expr(s.Cond)
	n := &If{Cond: cond}
	n.Then = b.into(func() {
		b.push(scopeBlock)
		b.stmtList(s.Body.List)
		b.pop()
	})
	if s.Else != nil {
		n.Else = b.into(func() {
			b.push(scopeBlock)
			b.stmt(s.Else)
			b.pop()
		})
	}
	b.seq.add(n)
}

// whileStmt is the only loop primitive. The condition sits at the head of
// the body, so `continue` re-evaluates it without a second edge shape.
func (b *funcBuilder) whileStmt(s *ast.WhileStmt) {
	loop := &Loop{}
	loop.Body = b.into(func() {
		b.push(scopeLoop)
		cond := b.expr(s.Cond)
		not := b.op(s.Cond.Pos(), OpNot, I1, cond)
		b.seq.add(&If{Cond: not, Then: &Seq{List: []Stmt{&Break{}}}})
		b.stmtList(s.Body.List)
		b.pop()
	})
	b.seq.add(loop)
}

// forStmt lowers A.5.6's single shape. The iterable decides the desugaring;
// the marker on the binding decides whether each element is shared, mutated
// in place, or consumed.
func (b *funcBuilder) forStmt(s *ast.ForStmt) {
	src := b.info().TypeOf(s.X)
	switch b.l.classify(src) {
	case kInt: // a range: A.4.5's `..`, always half-open
		b.forRange(s)
	case kArray, kSlice:
		b.forIndexed(s, src)
	case kMap:
		b.l.todo(s.Pos(), "map iteration — builtins/map iter_init+iter_next, order unspecified (A.5.6)")
	case kString:
		b.l.todo(s.Pos(), "string iteration — builtins/string decode at variable stride (A.5.6)")
	default:
		b.l.todo(s.Pos(), "for over %s", types.TypeString(src))
	}
}

func (b *funcBuilder) forRange(s *ast.ForStmt) {
	rng, ok := s.X.(*ast.BinaryExpr)
	if !ok {
		b.l.errorf(s.Pos(), "internal: integer for-iterable is not a range")
		return
	}
	t := b.l.hirType(b.info().TypeOf(rng.X))
	lo, hi := b.expr(rng.X), b.expr(rng.Y)
	slot := b.alloca(s.Pos(), t)
	b.store(s.Pos(), t, slot, lo)

	loop := &Loop{}
	loop.Body = b.into(func() {
		b.push(scopeLoop)
		cur := b.load(s.Pos(), t, slot)
		done := b.op(s.Pos(), b.l.cmpOp(true, "ge", b.info().TypeOf(rng.X)), I1, cur, hi)
		b.seq.add(&If{Cond: done, Then: &Seq{List: []Stmt{&Break{}}}})
		if len(s.Names) > 0 && !s.Names[0].IsBlank() {
			if obj := b.info().Defs[s.Names[0]]; obj != nil {
				b.declareBinding(obj, &binding{Value: cur, Type: t, Src: obj.Type()})
			}
		}
		b.stmtList(s.Body.List)
		next := b.op(s.Pos(), OpAdd, t, cur, IntVal(t, 1))
		b.store(s.Pos(), t, slot, next)
		b.pop()
	})
	b.seq.add(loop)
}

func (b *funcBuilder) forIndexed(s *ast.ForStmt, src types.Type) {
	elemT := b.l.hirType(b.l.elem(src))
	base, length := b.sequenceParts(s.X, src)

	slot := b.alloca(s.Pos(), I64)
	b.store(s.Pos(), I64, slot, IntVal(I64, 0))

	loop := &Loop{}
	loop.Body = b.into(func() {
		b.push(scopeLoop)
		i := b.load(s.Pos(), I64, slot)
		done := b.op(s.Pos(), OpUge, I1, i, length)
		b.seq.add(&If{Cond: done, Then: &Seq{List: []Stmt{&Break{}}}})

		ep := b.indexPtr(s.Pos(), elemT, base, i)
		names := s.Names
		if len(names) == 2 {
			if obj := b.info().Defs[names[0]]; obj != nil && !names[0].IsBlank() {
				b.declareBinding(obj, &binding{Value: i, Type: I64, Src: obj.Type()})
			}
			names = names[1:]
		}
		if len(names) == 1 && !names[0].IsBlank() {
			if obj := b.info().Defs[names[0]]; obj != nil {
				// The bare form iterates by shared access; `mut` iterates
				// by exclusive access; `var` consumes each element, and
				// the container is dead after the loop.
				bd := &binding{Type: elemT, Src: obj.Type()}
				switch s.Mode.String() {
				case "mut":
					bd.Value, bd.Slot = ep, true
				case "var":
					bd.Value, bd.Slot = ep, true
					bd.Owning = true
				default:
					bd.Value = b.load(s.Pos(), elemT, ep)
					bd.Slot = IsAggregate(elemT)
					if bd.Slot {
						bd.Value = ep
					}
				}
				b.declareBinding(obj, bd)
			}
		}
		b.stmtList(s.Body.List)
		b.store(s.Pos(), I64, slot, b.op(s.Pos(), OpAdd, I64, i, IntVal(I64, 1)))
		b.pop()
	})
	b.seq.add(loop)
}

// sequenceParts yields (base pointer, length) for an array or slice.
func (b *funcBuilder) sequenceParts(x ast.Expr, src types.Type) (Value, Value) {
	v := b.expr(x)
	switch b.l.classify(src) {
	case kArray:
		_, n := b.l.arrayParts(src)
		return v, IntVal(I64, n)
	case kSlice:
		st := b.l.hirType(src).(StructType)
		return b.loadField(x.Pos(), st.Def, v, "ptr"), b.loadField(x.Pos(), st.Def, v, "len")
	}
	return v, IntVal(I64, 0)
}

func (b *funcBuilder) switchStmt(s *ast.SwitchStmt) {
	tagType := b.info().TypeOf(s.Tag)
	tag := b.expr(s.Tag)

	// A payload enum's discriminant is its tag field; a unit-only enum
	// *is* its discriminant integer, so nothing is loaded.
	if b.l.classify(tagType) == kEnum {
		if st, ok := b.l.hirType(tagType).(StructType); ok {
			tag = b.loadField(s.Pos(), st.Def, tag, "tag")
		}
	}

	n := &Switch{Tag: tag}
	for _, c := range s.Cases {
		body := b.into(func() {
			b.push(scopeBlock)
			b.caseBody(c, tagType, tag)
			b.pop()
		})
		if c.Patterns == nil {
			n.Default = body
			continue
		}
		var vals []int64
		for _, p := range c.Patterns {
			v, ok := b.patternTag(p, tagType)
			if !ok {
				b.l.todo(p.Pos(), "non-constant switch pattern — vir's switch takes int-literal cases")
				continue
			}
			vals = append(vals, v)
		}
		n.Cases = append(n.Cases, SwitchCase{Values: vals, Body: body})
	}
	b.seq.add(n)
}

// caseBody binds an enum pattern's payload before the clause runs. A
// payload binding is a *view* into the payload, not a copy (A.5.7).
func (b *funcBuilder) caseBody(c *ast.CaseClause, tagType types.Type, tag Value) {
	for _, p := range c.Patterns {
		ep, ok := p.(*ast.EnumPattern)
		if !ok || len(ep.Binds) == 0 {
			continue
		}
		b.l.todo(ep.Pos(), "enum payload binding — a view into the tagged union's payload bytes")
	}
	b.stmtList(c.Body)
	for _, st := range c.Body {
		if br, ok := st.(*ast.BranchStmt); ok && br.Tok.String() == "fallthrough" {
			b.l.todo(br.Pos(), "fallthrough — duplicate the following clause's body, since vir's switch cases do not fall through")
		}
	}
}

func (b *funcBuilder) patternTag(p ast.Expr, tagType types.Type) (int64, bool) {
	if tv, ok := b.info().Types[p]; ok && tv.Value != nil {
		if v, isInt := types.Int64Val(tv.Value); isInt {
			return v, true
		}
	}
	switch p := p.(type) {
	case *ast.EnumPattern:
		return b.l.variantTag(tagType, p.Name.Name)
	case *ast.EnumShorthand:
		return b.l.variantTag(tagType, p.Name.Name)
	}
	return 0, false
}

func (b *funcBuilder) returnStmt(s *ast.ReturnStmt) {
	var out *Value
	switch len(s.Results) {
	case 0:
	case 1:
		v := b.owningExpr(s.Results[0], b.resultType())
		out = &v
	default:
		// A multi-value return is a bare comma list that unbuilds a tuple
		// (A.5.3); the tuple is what actually travels, through sret.
		v := b.tupleFromValues(s.Pos(), s.Results)
		out = &v
	}

	// If the result is an aggregate, it is written into the caller's sret
	// destination and the vir function returns void.
	if len(b.fn.Params) > 0 && b.fn.Params[0].SRet != nil && out != nil {
		dst := Ref(b.fn.Params[0].Name, StructType{b.fn.Params[0].SRet})
		b.store(s.Pos(), StructType{b.fn.Params[0].SRet}, dst, *out)
		out = nil
	}

	b.unwind(scopeFunc)
	b.seq.add(&Return{Value: out})
}

func (b *funcBuilder) resultType() types.Type {
	sig, _ := b.fn.Origin.(*types.Func)
	if sig == nil {
		return nil
	}
	s, _ := sig.Type().(*types.Signature)
	if s == nil || s.Results() == nil || s.Results().Len() == 0 {
		return nil
	}
	if s.Results().Len() == 1 {
		return s.Results().At(0).Type()
	}
	return s.Results()
}

// deferStmt records a call whose arguments are evaluated now and whose call
// is postponed to every exit edge of the enclosing scope.
func (b *funcBuilder) deferStmt(s *ast.DeferStmt) {
	saved := b.seq
	b.seq = &Seq{}
	call := b.expr(s.Call) // arguments evaluated at registration
	pending := b.seq
	b.seq = saved

	// Everything but the call itself stays where it was written.
	var last *Instr
	for _, st := range pending.List {
		run, ok := st.(*Instrs)
		if !ok {
			continue
		}
		for i, in := range run.List {
			if i == len(run.List)-1 && in.Op == OpCall {
				last = in
				continue
			}
			b.emit(in)
		}
	}
	_ = call
	if last == nil {
		b.l.errorf(s.Pos(), "internal: defer's operand did not lower to a call")
		return
	}
	sc := b.scopes[len(b.scopes)-1]
	sc.defers = append(sc.defers, last)
}

// compoundBase strips a compound assignment down to its base operator:
// += -> +, <<= -> <<, and so on.
//
// It yields a spelling rather than a token.Kind because token offers no
// reverse lookup from a spelling to a Kind, and reconstructing one by
// arithmetic on the enum would be a guess about an ordering token never
// promised. The spelling is the source of truth either way — binaryOp
// switches on it — so the base travels as a string and binaryOpFor answers
// the same question expr.go asks. The analyzer already validated the pair.
func compoundBase(k token.Kind) string {
	s := k.String()
	if len(s) > 1 && s[len(s)-1] == '=' {
		return s[:len(s)-1]
	}
	return s
}