// expr.go
package hir

import (
	"strconv"
	"fmt"
	"os"
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// expr lowers one expression to a Value. A Value of aggregate type is a
// pointer to storage (see the package doc comment).
func (b *funcBuilder) expr(x ast.Expr) Value {
	switch x := x.(type) {
	case *ast.ParenExpr:
		return b.expr(x.X)
	case *ast.BasicLit:
		return b.basicLit(x)
	case *ast.Ident:
		return b.ident(x)
	case *ast.UnaryExpr:
		return b.unary(x)
	case *ast.BinaryExpr:
		return b.binary(x)
	case *ast.CastExpr:
		return b.cast(x)
	case *ast.CallExpr:
		return b.callExpr(x)
	case *ast.SelectorExpr:
		return b.selector(x)
	case *ast.TupleIndexExpr:
		return b.tupleIndex(x)
	case *ast.IndexExpr:
		return b.index(x)
	case *ast.CompositeLit:
		return b.compositeLit(x)
	case *ast.ArrayLit:
		return b.arrayLit(x)
	case *ast.TupleExpr:
		return b.tupleLit(x)
	case *ast.TransferExpr:
		// A transfer is a move: the source dies here and the destination
		// becomes sole owner. Its *absence* is what synthesizes a copy, so
		// the marker itself lowers to nothing but a liveness edit.
		v := b.expr(x.Target)
		if bd := b.bindingFor(x.Target); bd != nil {
			bd.Dead = true
		}
		return v
	case *ast.LaunchExpr:
		return b.launch(x)
	case *ast.AwaitExpr:
		b.l.todo(x.Pos(), "await — a state boundary yielding Pending, produced by the async split (async.go)")
		return Value{}
	case *ast.FuncLit:
		b.l.todo(x.Pos(), "function literal — captures by value at creation, so this needs an env struct and a two-word closure")
		return Value{}
	case *ast.EnumShorthand:
		return b.enumShorthand(x)
	case *ast.MapLit:
		b.l.todo(x.Pos(), "map literal — builtins/map new + set per entry")
		return Value{}
	case *ast.BadExpr:
		return Value{}
	}
	b.l.todo(x.Pos(), "expression %T", x)
	return Value{}
}

func (b *funcBuilder) basicLit(x *ast.BasicLit) Value {
	t := b.info().TypeOf(x)
	ht := b.l.hirType(t)
	tv, ok := b.info().Types[x]

	// Fallback: If the upstream expression typechecker is missing or bypassed,
	// we evaluate the raw AST literal locally so HIR lowering can proceed.
	if (!ok || tv.Value == nil) && x.Kind.String() != "nil" {
		switch x.Kind.String() {
		case "INT":
			if val, err := strconv.ParseInt(x.Value, 0, 64); err == nil {
				if t == nil {
					return IntVal(I32, val)
				}
				return IntVal(ht, val)
			}
		case "FLOAT":
			if val, err := strconv.ParseFloat(x.Value, 64); err == nil {
				if t == nil {
					return FloatVal(F64, val)
				}
				return FloatVal(ht, val)
			}
		case "TRUE":
			return BoolVal(true)
		case "FALSE":
			return BoolVal(false)
		case "STRING", "CHAR":
			s := x.Value
			if len(s) >= 2 && (s[0] == '"' || s[0] == '\'' || s[0] == '`') {
				s = s[1 : len(s)-1]
			}
			if x.Kind.String() == "CHAR" {
				if len(s) > 0 {
					return IntVal(I32, int64(s[0]))
				}
				return IntVal(I32, 0)
			}
			return b.stringConstant(x.Pos(), s)
		}
	}

	if ok && tv.Value != nil {
		if s, isStr := types.StringVal_(tv.Value); isStr {
			return b.stringConstant(x.Pos(), s)
		}
		if v, isInt := types.Int64Val(tv.Value); isInt {
			if IsFloat(ht) {
				return FloatVal(ht, float64(v))
			}
			return IntVal(ht, v)
		}
		if bv, isBool := types.BoolVal_(tv.Value); isBool {
			return BoolVal(bv)
		}
	}
	if x.Kind.String() == "nil" {
		return NullVal()
	}
	b.l.errorf(x.Pos(), "internal: literal %q has no constant value in types.Info", x.Value)
	return Value{}
}

// stringConstant materializes a string literal: the bytes go into a module
// global, and the value is a {ptr, len} header. A.1.5.2 is explicit that a
// string carries no NUL terminator; one is manufactured only at a declare
// boundary — see cStringArg below, which is where that manufacturing
// actually happens.
func (b *funcBuilder) stringConstant(pos token.Pos, s string) Value {
	g := &Global{
		Name: b.mod().uniqueName("str"),
		Type: ArrayType{Elem: I8, Len: int64(len(s))},
		Init: InitBytes{Data: []byte(s)},
	}
	b.mod().Globals = append(b.mod().Globals, g)

	st := b.l.hirType(types.Typ[types.String]).(StructType)
	slot := b.alloca(pos, st)
	b.storeField(pos, st.Def, slot, "ptr", GlobalVal(g.Name, Ptr))
	b.storeField(pos, st.Def, slot, "len", IntVal(I64, int64(len(s))))
	return Value{Kind: ValRef, Name: slot.Name, Type: st}
}

func (b *funcBuilder) ident(x *ast.Ident) Value {
	obj := b.info().ObjectOf(x)
	if obj == nil {
		return Value{}
	}
	if bd := b.lookup(obj); bd != nil {
		if bd.Slot && !IsAggregate(bd.Type) {
			return b.load(x.Pos(), bd.Type, bd.Value)
		}
		if bd.Slot {
			return Value{Kind: ValRef, Name: bd.Value.Name, Type: bd.Type}
		}
		return bd.Value
	}
	if g := b.l.globalFor(obj); g != nil {
		if g.isConst {
			return Value{Kind: ValConst, Name: g.name, Module: b.qualify(g.mod), Type: g.typ}
		}
		addr := GlobalVal(g.name, Ptr)
		if IsAggregate(g.typ) {
			return Value{Kind: ValRef, Name: g.name, Type: g.typ}
		}
		return b.load(x.Pos(), g.typ, addr)
	}
	if fn, ok := obj.(*types.Func); ok {
		return b.funcAddr(x.Pos(), fn)
	}
	if tv, ok := b.info().Types[x]; ok && tv.Value != nil {
		if v, isInt := types.Int64Val(tv.Value); isInt {
			return IntVal(b.l.hirType(obj.Type()), v)
		}
	}
	b.l.errorf(x.Pos(), "internal: identifier %q resolved to no lowered binding", x.Name)
	return Value{}
}

// funcAddr yields a function's address. vir has no address-of instruction,
// so the address is spelled the one way the grammar allows — a global
// initialized `addr <fn>` — and loaded. This is what the task ABI's
// {poll, drop} function pointers travel through.
func (b *funcBuilder) funcAddr(pos token.Pos, fn *types.Func) Value {
	target := b.l.work.lookup(fn, nil)
	if target == nil {
		target = b.l.work.enqueue(b.l.cur.unit, fn, nil, b.l.cur.depth+1)
	}
	if target == nil {
		return Value{}
	}
	name := "addr_" + target.Name
	if g := findGlobal(b.mod(), name); g == nil {
		b.mod().Globals = append(b.mod().Globals, &Global{
			Name: b.mod().uniqueName(name), Type: Ptr, Init: InitAddr{Name: target.Name},
		})
	}
	return b.load(pos, Ptr, GlobalVal(name, Ptr))
}

// unary covers -, !, ~, and &. `&` is address-of on an ordinary value and
// dereference on a typed_ptr, keyed on the operand's statically written
// type — the analyzer already decided which, and recorded it.
func (b *funcBuilder) unary(x *ast.UnaryExpr) Value {
	t := b.info().TypeOf(x)
	ht := b.l.hirType(t)
	switch x.Op.String() {
	case "-":
		return b.op(x.OpPos, OpNeg, ht, b.expr(x.X))
	case "!":
		return b.op(x.OpPos, OpNot, I1, b.expr(x.X))
	case "~":
		return b.op(x.OpPos, OpNot, ht, b.expr(x.X))
	case "&":
		if b.l.classify(b.info().TypeOf(x.X)) == kPointer {
			return b.load(x.Pos(), ht, b.expr(x.X)) // dereference
		}
		return b.addressOf(x.X) // address-of
	}
	b.l.todo(x.Pos(), "unary %s", x.Op)
	return Value{}
}

func (b *funcBuilder) binary(x *ast.BinaryExpr) Value {
	op := x.Op.String()

	// && and || short-circuit, so they are control flow, not an opcode.
	if op == "&&" || op == "||" {
		slot := b.alloca(x.Pos(), I1)
		lhs := b.expr(x.X)
		b.store(x.Pos(), I1, slot, lhs)
		cond := lhs
		if op == "||" {
			cond = b.op(x.OpPos, OpNot, I1, lhs)
		}
		n := &If{Cond: cond}
		n.Then = b.into(func() { b.store(x.Pos(), I1, slot, b.expr(x.Y)) })
		b.seq.add(n)
		return b.load(x.Pos(), I1, slot)
	}

	srcT := b.info().TypeOf(x.X)

	// String equality is a builtin call, not an opcode: a string is a
	// {ptr,len} header and == asks about bytes, never about identity.
	if b.l.classify(srcT) == kString && (op == "==" || op == "!=") {
		b.l.need(builtinFeatureString)
		r := b.callBuiltin(x.OpPos, symStringCompare, I32, b.expr(x.X), b.expr(x.Y))
		want := OpEq
		if op == "!=" {
			want = OpNe
		}
		return b.op(x.OpPos, want, I1, r, IntVal(I32, 0))
	}

	lhs, rhs := b.expr(x.X), b.expr(x.Y)
	if hop := b.l.binaryOp(x.Op, srcT); hop != OpInvalid {
		rt := b.l.hirType(b.info().TypeOf(x))
		return b.op(x.OpPos, hop, rt, lhs, rhs)
	}
	b.l.todo(x.Pos(), "binary operator %s", op)
	return Value{}
}

// binaryOp picks the opcode for a written operator.
func (l *lowerer) binaryOp(k token.Kind, t types.Type) Op {
	return l.binaryOpFor(k.String(), t)
}

// binaryOpFor is the spelling-keyed table. Signedness lives here and nowhere
// else: the plain forms trap on overflow while &+ &- &* wrap, and vir's
// add/sub/mul already wrap modulo 2^N, so the trapping forms need an
// explicit overflow check.
//
// It takes a spelling rather than a token.Kind because a compound assignment
// must ask the same question about its base operator and token offers no way
// back from a spelling to a Kind — see stmt.go's compoundBase. One table,
// two callers, no second place for signedness to be decided.
func (l *lowerer) binaryOpFor(op string, t types.Type) Op {
	signed := l.isSigned(t)
	float := l.classify(t) == kFloat
	switch op {
	case "+", "&+":
		return OpAdd
	case "-", "&-":
		return OpSub
	case "*", "&*":
		return OpMul
	case "/":
		switch {
		case float:
			return OpSDiv // vir spells float division with the same mnemonic family
		case signed:
			return OpSDiv
		}
		return OpUDiv
	case "%":
		if signed {
			return OpSRem
		}
		return OpURem
	case "&":
		return OpAnd
	case "|":
		return OpOr
	case "^":
		return OpXor
	case "<<":
		return OpShl
	case ">>":
		if signed {
			return OpAShr
		}
		return OpLShr
	case "==", "===":
		return OpEq
	case "!=", "!==":
		return OpNe
	}
	return l.cmpOp(signed, op, t)
}

func (l *lowerer) cmpOp(signed bool, k string, t types.Type) Op {
	if l.classify(t) == kFloat {
		switch k {
		case "<", "lt":
			return OpLt
		case ">", "gt":
			return OpGt
		case "<=", "le":
			return OpLe
		case ">=", "ge":
			return OpGe
		}
		return OpInvalid
	}
	switch k {
	case "<", "lt":
		if signed {
			return OpSlt
		}
		return OpUlt
	case ">", "gt":
		if signed {
			return OpSgt
		}
		return OpUgt
	case "<=", "le":
		if signed {
			return OpSle
		}
		return OpUle
	case ">=", "ge":
		if signed {
			return OpSge
		}
		return OpUge
	}
	return OpInvalid
}

// cast lowers `as`. It never touches memory: between pointer types it is a
// static reinterpretation, between numeric types a width-selected
// truncate/extend/int-float instruction, on a unit-only enum a tag read.
// There is no dynamic cast because there is no runtime type information.
func (b *funcBuilder) cast(x *ast.CastExpr) Value {
	from := b.info().TypeOf(x.X)
	to := b.info().TypeOf(x)
	v := b.expr(x.X)
	ft, tt := b.l.classify(from), b.l.classify(to)
	hto := b.l.hirType(to)

	switch {
	case ft == kPointer || ft == kAbstract, tt == kPointer:
		return v // reinterpretation; nothing to emit
	case (ft == kInt || ft == kChar || ft == kEnum) && (tt == kInt || tt == kChar || tt == kEnum):
		fb, tb := b.l.intBits(from), b.l.intBits(to)
		switch {
		case tb == fb:
			return v
		case tb < fb:
			return b.op(x.As, OpTrunc, hto, v)
		case b.l.isSigned(from):
			return b.op(x.As, OpSext, hto, v)
		default:
			return b.op(x.As, OpZext, hto, v)
		}
	case ft == kInt && tt == kFloat:
		if b.l.isSigned(from) {
			return b.op(x.As, OpSfromint, hto, v)
		}
		return b.op(x.As, OpUfromint, hto, v)
	case ft == kFloat && tt == kInt:
		// A.11.2's saturating rule is the npu one; on the host, an
		// out-of-range float-to-int conversion is a deterministic trap, so
		// the saturating opcode is used only where the source asked for it.
		if b.l.isSigned(to) {
			return b.op(x.As, OpStointSat, hto, v)
		}
		return b.op(x.As, OpUtointSat, hto, v)
	case ft == kFloat && tt == kFloat:
		if b.l.floatBits(to) < b.l.floatBits(from) {
			return b.op(x.As, OpFdemote, hto, v)
		}
		return b.op(x.As, OpFpromote, hto, v)
	}
	b.l.todo(x.Pos(), "cast %s as %s", types.TypeString(from), types.TypeString(to))
	return v
}

// callExpr routes a call four ways: a reserved builtin name, a
// module-local extern function (a declare-block member, A.8), a method
// through a recorded Selection, or an ordinary direct call. There are no
// vtables to consult — every call, including a generic one post-
// monomorphization, is direct by construction.
func (b *funcBuilder) callExpr(x *ast.CallExpr) Value {
	fun := unparen(x.Fun)

	if id, ok := fun.(*ast.Ident); ok {
		if id.Name == "printf" {
			o := b.info().ObjectOf(id)
			fmt.Fprintf(os.Stderr, "CALL %s: obj=%p externFound=%v\n",
				id.Name, o, b.l.externFor(o) != nil)
		}
		if bi, isBuiltin := b.info().ObjectOf(id).(*types.Builtin); isBuiltin {
			return b.builtinCall(x, bi)
		}
		// An extern has no Vertex body to instantiate, and must never
		// reach resolveCallee/the monomorphization worklist below — see
		// lower.go's externs field doc for what that misrouting used to
		// produce.
		if ef := b.l.externFor(b.info().ObjectOf(id)); ef != nil {
			return b.externCallExpr(x, ef)
		}
	}
	if sel, ok := fun.(*ast.SelectorExpr); ok {
		if s := b.info().Selections[sel]; s != nil && s.Kind == types.MethodVal {
			return b.methodCall(x, sel, s)
		}
	}

	target, args := b.resolveCallee(fun)
	if target == nil {
		return Value{}
	}
	vals := b.arguments(x, target)
	_ = args
	return b.callWithSRet(x.Pos(), target, vals)
}

// externCallExpr lowers a call to a module-local extern function (a
// declare-block member, A.8). There is no *hir.Func to instantiate here —
// an extern has no Vertex body to monomorphize, so it must never reach the
// worklist. Routing it there is exactly the bug this replaces: the
// worklist would schedule a same-named shell (colliding with the extern's
// own reserved name and getting suffixed, e.g. "printf" -> "printf_1"),
// findFuncDecl would find no ast.FuncDecl behind it, and lower/vir would
// emit a declared-but-unterminated function.
func (b *funcBuilder) externCallExpr(x *ast.CallExpr, ef *ExternFunc) Value {
	pos := x.Pos()
	args := make([]Value, 0, len(x.Args))
	var temps []Value
	for _, a := range x.Args {
		if kv, ok := a.(*ast.KeyValueExpr); ok {
			a = kv.Value
		}
		// A byval param is a pointer at the vir level (decl.go's param);
		// an aggregate hir.Value already *is* that pointer (the package
		// doc comment), so a plain b.expr is correct for both scalar and
		// byval-aggregate positions, including a trailing variadic
		// argument (A.4.4's `...` on the declare-block signature) — vir's
		// call site takes a flat operand list regardless of arity, so a
		// variadic argument needs no different lowering than a fixed one.
		// There is no owning-marker system to consult here the way
		// owningExpr has for a Vertex callee: ownership markers are a
		// Vertex-source convention (A.9.1) with nothing on the other side
		// of a declare boundary to honor it.
		//
		// One case is not "pass through unchanged", though: a Vertex
		// `string` has no C-ABI shape of its own (decl.go's foreignParam
		// doc explains why), so it is bridged here rather than at the
		// signature alone — this covers a string reaching a fixed
		// parameter and one reaching the `...` tail alike, since a
		// variadic position has no declared Param for foreignParam to
		// have marked in the first place.
		//
		// b.externArgType, not b.info().TypeOf, decides this: the
		// analyzer's phase 3 (resolve.go's checkBodies) resolves
		// identifiers only and "computes no types for expressions", so a
		// literal argument reaching a call directly — never bound to a
		// name, never assigned — has no Types[e] entry to read. TypeOf's
		// own fallback recovers a bare identifier through ObjectOf, but
		// nothing recovers a bare literal; externArgType is that recovery,
		// scoped to exactly the classification this call needs.
		if b.l.classify(b.externArgType(a)) == kString {
			buf := b.cStringArg(a)
			args = append(args, buf)
			temps = append(temps, buf)
			continue
		}
		args = append(args, b.expr(a))
	}
	res := b.callExtern(pos, ef.Name, ef.Result, args...)
	// Each marshaled buffer is call-scoped: built fresh for this call and
	// freed immediately after, never bound to anything a mut/var/transfer
	// marker could reach and never live past the call that used it.
	for _, t := range temps {
		b.callBuiltin(pos, symMemFree, Void, t)
	}
	return res
}

// externArgType is TypeOf with one extra fallback, needed only at an
// extern call boundary: b.info().TypeOf falls back through ObjectOf for a
// bare identifier, but has nothing to fall back to for a bare literal
// argument — analyzer's phase 3 never records a Types[e] entry for one,
// since resolving identifiers is that phase's whole job (resolve.go's
// checkBodies doc comment is explicit that it "computes no types for
// expressions"). A literal reached directly as a call argument, with no
// let/var binding in between, is exactly that case, and cStringArg's
// caller needs to know it's looking at a string before it can marshal one.
//
// This mirrors basicLit's own recovery for the identical gap, deliberately
// narrow: it classifies a literal for this one decision and is not a
// general type-inference fallback. An arbitrary non-literal expression
// with no recorded type (e.g. one produced by a construct the checker
// doesn't yet type) still returns nil here, exactly as TypeOf would.
func (b *funcBuilder) externArgType(a ast.Expr) types.Type {
	if t := b.info().TypeOf(a); t != nil {
		return t
	}
	lit, ok := unparen(a).(*ast.BasicLit)
	if !ok {
		return nil
	}
	switch lit.Kind.String() {
	case "STRING":
		return types.Typ[types.String]
	case "CHAR":
		return types.Typ[types.Char]
	case "TRUE", "FALSE":
		return types.Typ[types.Bool]
	case "INT":
		return types.Typ[types.Int]
	case "FLOAT":
		return types.Typ[types.Float64]
	}
	return nil
}

// cStringArg marshals a Vertex string argument into a heap buffer holding
// its bytes plus a manufactured NUL terminator, and yields the pointer —
// the conversion A.1.5.2 names for exactly this seam, and the one
// externCallExpr performs for every string-typed argument crossing a
// declare boundary. The buffer is temporary: externCallExpr frees it
// immediately after the call it was built for, so nothing here needs an
// owning binding or a deinit routine of its own.
func (b *funcBuilder) cStringArg(x ast.Expr) Value {
	pos := x.Pos()
	sv := b.expr(x)
	st, ok := b.l.hirType(b.info().TypeOf(x)).(StructType)
	if !ok {
		b.l.errorf(pos, "internal: cStringArg on a non-string operand")
		return Value{}
	}
	src := b.loadField(pos, st.Def, sv, "ptr")
	n := b.loadField(pos, st.Def, sv, "len")
	size := b.op(pos, OpAdd, I64, n, IntVal(I64, 1)) // +1 for the manufactured NUL

	buf := b.callBuiltin(pos, symMemAllocate, Ptr, size)
	oom := &If{Cond: b.op(pos, OpEq, I1, buf, NullVal())}
	oom.Then = b.into(func() {
		// Matches builtinBox's convention (builtin.go): an allocation this
		// package makes on the program's behalf, rather than one the user
		// wrote a `new`/`resize` call for, fails loudly (A.10.1's split).
		b.callBuiltin(pos, symPanicOOM, Void)
		b.seq.add(&Unreachable{})
	})
	b.seq.add(oom)

	b.opVoid(pos, OpMemcopy, Void, buf, src, n)
	b.store(pos, I8, b.indexPtr(pos, I8, buf, n), IntVal(I8, 0))
	return buf
}

// callWithSRet supplies the destination for an aggregate result and hands
// the caller back a pointer to it.
func (b *funcBuilder) callWithSRet(pos token.Pos, f *Func, args []Value) Value {
	if len(f.Params) > 0 && f.Params[0].SRet != nil {
		st := StructType{f.Params[0].SRet}
		dst := b.alloca(pos, st)
		b.call(pos, f, append([]Value{dst}, args...)...)
		return Value{Kind: ValRef, Name: dst.Name, Type: st}
	}
	return b.call(pos, f, args...)
}

// arguments lowers an argument list. Named arguments resolve to positional
// order at compile time and leave no trace in the binary, and every
// argument is an owning position, so any of them may carry the marker.
func (b *funcBuilder) arguments(x *ast.CallExpr, target *Func) []Value {
	out := make([]Value, 0, len(x.Args))
	params := target.Params
	if len(params) > 0 && params[0].SRet != nil {
		params = params[1:]
	}
	for i, a := range x.Args {
		if kv, ok := a.(*ast.KeyValueExpr); ok {
			a = kv.Value // already resolved to position by the analyzer
		}
		var want types.Type
		if i < len(params) {
			want = b.info().TypeOf(a)
		}
		out = append(out, b.owningExpr(a, want))
	}
	return out
}

func (b *funcBuilder) methodCall(x *ast.CallExpr, sel *ast.SelectorExpr, s *types.Selection) Value {
	fn, ok := s.Obj.(*types.Func)
	if !ok {
		return Value{}
	}
	target := b.instantiate(fn, sel.Sel)
	if target == nil {
		return Value{}
	}
	// A receiver typed `var` consumes its receiver unconditionally: the
	// receiver position has no argument slot to carry a marker, so there
	// is no bare form that copies (A.6.1's single exception).
	recv := b.receiverValue(sel.X, fn)
	if b.l.consumesReceiver(fn) {
		if bd := b.bindingFor(sel.X); bd != nil {
			bd.Dead = true
		}
	}
	return b.callWithSRet(x.Pos(), target, append([]Value{recv}, b.arguments(x, target)...))
}

func (b *funcBuilder) receiverValue(x ast.Expr, fn *types.Func) Value {
	if b.l.needsAddressableReceiver(fn) {
		return b.addressOf(x)
	}
	return b.expr(x)
}

// resolveCallee maps a call target expression to a monomorphized *Func,
// scheduling the instantiation if this is the first reference to it.
func (b *funcBuilder) resolveCallee(fun ast.Expr) (*Func, []Value) {
	switch f := fun.(type) {
	case *ast.Ident:
		obj := b.info().ObjectOf(f)
		if fn, ok := obj.(*types.Func); ok {
			return b.instantiate(fn, f), nil
		}
	case *ast.IndexExpr:
		// Explicit type arguments: f[T](a).
		if id, ok := unparen(f.X).(*ast.Ident); ok {
			if fn, ok := b.info().ObjectOf(id).(*types.Func); ok {
				return b.instantiate(fn, id), nil
			}
		}
	case *ast.SelectorExpr:
		if s := b.info().Selections[f]; s != nil && s.Kind == types.PackageMember {
			if fn, ok := s.Obj.(*types.Func); ok {
				return b.instantiate(fn, f.Sel), nil
			}
		}
	}
	b.l.todo(fun.Pos(), "indirect call through a function value — needs a vir fnsig declaration")
	return nil, nil
}

// instantiate resolves the concrete type arguments at this site and
// schedules the instance, memoized by (Func, ConcreteTypeArgs).
func (b *funcBuilder) instantiate(fn *types.Func, at *ast.Ident) *Func {
	var args []types.Type
	if inst, ok := b.info().Instances[at]; ok {
		for _, a := range inst.TypeArgs {
			args = append(args, b.l.subst(a))
		}
	}
	unit := b.l.unitFor(fn)
	if unit == nil {
		b.l.errorf(fn.Pos(), "internal: no loaded unit declares %s", fn.Name())
		return nil
	}
	if f := b.l.work.lookup(fn, args); f != nil {
		return f
	}
	return b.l.work.enqueue(unit, fn, args, b.l.cur.depth+1)
}

func (b *funcBuilder) selector(x *ast.SelectorExpr) Value {
	s := b.info().Selections[x]
	if s == nil {
		return b.expr(x.Sel)
	}
	switch s.Kind {
	case types.FieldVal:
		base := b.addressOf(x.X)
		st, ok := b.l.hirType(b.info().TypeOf(x.X)).(StructType)
		if !ok {
			b.l.errorf(x.Pos(), "internal: field access on a non-struct")
			return Value{}
		}
		f, _ := st.Def.Field(sanitize(x.Sel.Name))
		p := b.fieldPtr(x.Pos(), st.Def, base, f.Name)
		return b.load(x.Pos(), f.Type, p)
	case types.PackageMember:
		return b.ident(x.Sel)
	case types.NamespaceMember:
		b.l.todo(x.Pos(), "namespace member %s — async/gpu/npu/chan namespaces", x.Sel.Name)
		return Value{}
	}
	b.l.todo(x.Pos(), "selector kind %v", s.Kind)
	return Value{}
}

func (b *funcBuilder) tupleIndex(x *ast.TupleIndexExpr) Value {
	st, ok := b.l.hirType(b.info().TypeOf(x.X)).(StructType)
	if !ok || x.Index >= len(st.Def.Fields) {
		b.l.errorf(x.Pos(), "internal: tuple index out of range")
		return Value{}
	}
	f := st.Def.Fields[x.Index]
	return b.load(x.Pos(), f.Type, b.fieldPtr(x.Pos(), st.Def, b.addressOf(x.X), f.Name))
}

func (b *funcBuilder) index(x *ast.IndexExpr) Value {
	src := b.info().TypeOf(x.X)
	switch b.l.classify(src) {
	case kArray, kSlice:
		elem := b.l.hirType(b.l.elem(src))
		base, length := b.sequenceParts(x.X, src)
		idx := b.expr(x.Indices[0])
		b.boundsCheck(x.Pos(), idx, length)
		return b.load(x.Pos(), elem, b.indexPtr(x.Pos(), elem, base, idx))
	case kMap:
		b.l.todo(x.Pos(), "map subscript — builtins/map get")
	case kPointer:
		elem := b.l.hirType(b.l.elem(src))
		return b.load(x.Pos(), elem, b.indexPtr(x.Pos(), elem, b.expr(x.X), b.expr(x.Indices[0])))
	}
	b.l.todo(x.Pos(), "index into %s", types.TypeString(src))
	return Value{}
}

// boundsCheck is not optional: A.15 lists out-of-bounds access among the
// undefined behaviors vir permits, and Vertex declines to inherit it for a
// checked container.
func (b *funcBuilder) boundsCheck(pos token.Pos, idx, length Value) {
	bad := b.op(pos, OpUge, I1, idx, length)
	n := &If{Cond: bad}
	n.Then = b.into(func() {
		b.callBuiltin(pos, symSliceOOB, Void, idx, length)
		b.seq.add(&Unreachable{})
	})
	b.seq.add(n)
}

func (b *funcBuilder) compositeLit(x *ast.CompositeLit) Value {
	t := b.info().TypeOf(x)
	st, ok := b.l.hirType(t).(StructType)
	if !ok {
		b.l.errorf(x.Pos(), "internal: composite literal of non-struct type")
		return Value{}
	}
	slot := b.alloca(x.Pos(), st)
	b.opVoid(x.Pos(), OpMemset, Void, slot, IntVal(I8, 0), IntVal(I64, st.Def.Size))
	written := map[string]bool{}
	for _, e := range x.Elems {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, _ := kv.Key.(*ast.Ident)
		if key == nil {
			continue
		}
		name := sanitize(key.Name)
		f, _ := st.Def.Field(name)
		b.store(kv.Value.Pos(), f.Type, b.fieldPtr(x.Pos(), st.Def, slot, name),
			b.owningExpr(kv.Value, b.info().TypeOf(kv.Value)))
		written[name] = true
	}
	// A.6.2: a field default is evaluated at construction for any field
	// the literal omits. Zeroing above covers the no-default case.
	b.l.emitFieldDefaults(b, x.Pos(), t, st, slot, written)
	return Value{Kind: ValRef, Name: slot.Name, Type: st}
}

func (b *funcBuilder) arrayLit(x *ast.ArrayLit) Value {
	t := b.info().TypeOf(x)
	at, ok := b.l.hirType(t).(ArrayType)
	if !ok {
		b.l.todo(x.Pos(), "array literal of slice type — needs a builtins/slice allocation")
		return Value{}
	}
	slot := b.alloca(x.Pos(), at)
	for i, e := range x.Elems {
		p := b.indexPtr(e.Pos(), at.Elem, slot, IntVal(I64, int64(i)))
		b.store(e.Pos(), at.Elem, p, b.owningExpr(e, b.info().TypeOf(e)))
	}
	return Value{Kind: ValRef, Name: slot.Name, Type: at}
}

func (b *funcBuilder) tupleLit(x *ast.TupleExpr) Value {
	return b.tupleFromValues(x.Pos(), x.Elems)
}

func (b *funcBuilder) tupleFromValues(pos token.Pos, elems []ast.Expr) Value {
	var ts []types.Type
	for _, e := range elems {
		ts = append(ts, b.info().TypeOf(e))
	}
	st, ok := b.l.hirType(types.NewTuple(varsOf(ts)...)).(StructType)
	if !ok {
		return Value{}
	}
	slot := b.alloca(pos, st)
	for i, e := range elems {
		if i >= len(st.Def.Fields) {
			break
		}
		f := st.Def.Fields[i]
		b.store(e.Pos(), f.Type, b.fieldPtr(pos, st.Def, slot, f.Name),
			b.owningExpr(e, ts[i]))
	}
	return Value{Kind: ValRef, Name: slot.Name, Type: st}
}

func (b *funcBuilder) enumShorthand(x *ast.EnumShorthand) Value {
	t := b.info().TypeOf(x)
	tag, ok := b.l.variantTag(t, x.Name.Name)
	if !ok {
		b.l.errorf(x.Pos(), "internal: no discriminant for variant %s", x.Name.Name)
		return Value{}
	}
	ht := b.l.hirType(t)
	if IsInt(ht) {
		return IntVal(ht, tag) // a unit-only enum *is* its discriminant
	}
	st := ht.(StructType)
	slot := b.alloca(x.Pos(), st)
	b.storeField(x.Pos(), st.Def, slot, "tag", IntVal(st.Def.Fields[0].Type, tag))
	if len(x.Args) > 0 {
		b.l.todo(x.Pos(), "payload variant construction — write args into the payload bytes")
	}
	return Value{Kind: ValRef, Name: slot.Name, Type: st}
}

// launch lowers A.4.2's call-expression prefixes. Each modifies how a call
// is scheduled, never the callee's signature.
func (b *funcBuilder) launch(x *ast.LaunchExpr) Value {
	switch x.Kw.String() {
	case "thread":
		b.l.todo(x.Pos(), "thread launch — builtins/thread spawn, result delivered on a chan T")
	case "async":
		b.l.todo(x.Pos(), "async launch — channel handshake dispatching an ordinary call to the reactor")
	case "gpu", "npu":
		// A device-marked function is out of scope for VIR emission: the
		// launch call site validates, but the kernel body has nowhere to
		// go. vvm lists "no host<->device story" as its largest gap.
		b.l.todo(x.Pos(), "%s launch — no lower/gvir counterpart exists", x.Kw)
	}
	return Value{}
}

// addressOf yields a pointer to an assignable location.
func (b *funcBuilder) addressOf(x ast.Expr) Value {
	switch x := unparen(x).(type) {
	case *ast.Ident:
		obj := b.info().ObjectOf(x)
		if bd := b.lookup(obj); bd != nil {
			if bd.Slot {
				return Value{Kind: ValRef, Name: bd.Value.Name, Type: Ptr}
			}
			// A let has no storage slot; anything needing its address gets
			// a fresh one, which is sound because a let is immutable.
			slot := b.alloca(x.Pos(), bd.Type)
			b.store(x.Pos(), bd.Type, slot, bd.Value)
			return slot
		}
		if g := b.l.globalFor(obj); g != nil {
			return GlobalVal(g.name, Ptr)
		}
	case *ast.SelectorExpr:
		st, ok := b.l.hirType(b.info().TypeOf(x.X)).(StructType)
		if ok {
			return b.fieldPtr(x.Pos(), st.Def, b.addressOf(x.X), sanitize(x.Sel.Name))
		}
	case *ast.TupleIndexExpr:
		if st, ok := b.l.hirType(b.info().TypeOf(x.X)).(StructType); ok && x.Index < len(st.Def.Fields) {
			return b.fieldPtr(x.Pos(), st.Def, b.addressOf(x.X), st.Def.Fields[x.Index].Name)
		}
	case *ast.IndexExpr:
		src := b.info().TypeOf(x.X)
		elem := b.l.hirType(b.l.elem(src))
		base, length := b.sequenceParts(x.X, src)
		idx := b.expr(x.Indices[0])
		b.boundsCheck(x.Pos(), idx, length)
		return b.indexPtr(x.Pos(), elem, base, idx)
	case *ast.UnaryExpr:
		if x.Op.String() == "&" {
			return b.expr(x.X) // &p = v writes through p
		}
	}
	b.l.todo(x.Pos(), "address of %T", x)
	return Value{}
}

// bindingFor returns the binding an expression names, if it names one.
// Only a binding or a field path can be transferred, so this is where the
// liveness edit for `var` lands.
func (b *funcBuilder) bindingFor(x ast.Expr) *binding {
	switch x := unparen(x).(type) {
	case *ast.Ident:
		return b.lookup(b.info().ObjectOf(x))
	case *ast.SelectorExpr:
		return b.bindingFor(x.X)
	case *ast.TransferExpr:
		return b.bindingFor(x.Target)
	}
	return nil
}

// owningExpr lowers an expression in an owning position (A.9.1). This is
// the single place move-versus-copy is decided, and it is decided by asking
// types.Info.Transfers — the one side table that changes generated code
// rather than merely licensing it.
func (b *funcBuilder) owningExpr(x ast.Expr, want types.Type) Value {
	if tr, ok := unparen(x).(*ast.TransferExpr); ok {
		// Transfer: header only, O(1). The source dies here.
		v := b.expr(tr.Target)
		if bd := b.bindingFor(tr.Target); bd != nil {
			bd.Dead = true
		}
		return v
	}
	v := b.expr(x)
	if want == nil || !b.l.classify(want).owning() {
		return v
	}
	// A temporary has no prior binding to preserve, so it is simply
	// consumed into its destination — the marker would be meaningless.
	if b.bindingFor(x) == nil {
		return v
	}
	return b.l.own.emitCopy(b, x.Pos(), want, v)
}

func (b *funcBuilder) zeroValue(pos token.Pos, t Type) Value {
	if IsAggregate(t) {
		slot := b.alloca(pos, t)
		b.opVoid(pos, OpMemset, Void, slot, IntVal(I8, 0), IntVal(I64, b.l.types.Sizeof(t)))
		return Value{Kind: ValRef, Name: slot.Name, Type: t}
	}
	switch t.(type) {
	case FloatType:
		return FloatVal(t, 0)
	case PtrType:
		return NullVal()
	}
	return IntVal(t, 0)
}

func (b *funcBuilder) qualify(m *Module) string {
	if m == b.mod() {
		return ""
	}
	b.mod().AddImport(m.Name)
	return m.Name
}

func unparen(x ast.Expr) ast.Expr {
	for {
		p, ok := x.(*ast.ParenExpr)
		if !ok {
			return x
		}
		x = p.X
	}
}

func findGlobal(m *Module, name string) *Global {
	for _, g := range m.Globals {
		if g.Name == name {
			return g
		}
	}
	return nil
}