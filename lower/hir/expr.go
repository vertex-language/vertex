package hir

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// owning lowers an expression in one of the six owning positions. The
// difference between a move and a deep copy rides entirely on whether the
// analyzer accepted a transfer marker here — Info.Transfers is the record,
// and it is the one analysis result that changes generated code rather than
// merely licensing it.
func (fb *funcBuilder) owning(e ast.Expr) Value {
	if tr, ok := e.(*ast.TransferExpr); ok {
		v := fb.expr(tr.Target)
		if b := fb.bindingOf(tr.Target); b != nil {
			// A transfer emits nothing. Its entire lowering is the *absence*
			// of a drop at the source's original end of liveness.
			b.dead = true
		}
		return v
	}
	// A bare owning use is a copy, at whatever the type's copy costs.
	return fb.copyOf(fb.expr(e), fb.l.info.TypeOf(e))
}

func (fb *funcBuilder) expr(e ast.Expr) Value {
	switch x := e.(type) {
	case nil, *ast.BadExpr:
		return Value{}

	case *ast.BasicLit:
		return fb.basicLit(x)

	case *ast.Ident:
		return fb.identValue(x)

	case *ast.ParenExpr:
		return fb.expr(x.X)

	case *ast.BinaryExpr:
		return fb.binary(x)

	case *ast.UnaryExpr:
		return fb.unary(x)

	case *ast.CastExpr:
		return fb.cast(x)

	case *ast.CallExpr:
		return fb.callExpr(x)

	case *ast.SelectorExpr:
		p := fb.address(x)
		t := fb.typ(fb.l.info.TypeOf(x))
		if IsAggregate(t) {
			return p
		}
		return fb.emit(OpLoad, t, p)

	case *ast.IndexExpr:
		p := fb.address(x)
		t := fb.typ(fb.l.info.TypeOf(x))
		if IsAggregate(t) {
			return p
		}
		return fb.emit(OpLoad, t, p)

	case *ast.TupleIndexExpr:
		p := fb.address(x)
		t := fb.typ(fb.l.info.TypeOf(x))
		if IsAggregate(t) {
			return p
		}
		return fb.emit(OpLoad, t, p)

	case *ast.CompositeLit:
		return fb.compositeLit(x)

	case *ast.TupleExpr:
		return fb.tupleLit(x)

	case *ast.ArrayLit:
		return fb.arrayLit(x)

	case *ast.HeapConstructor:
		return fb.heapConstructor(x)

	case *ast.TransferExpr:
		// Reached outside an owning position: the analyzer already rejected
		// it, so nothing valid gets here.
		fb.l.bugAt(x.Pos(), "transfer marker outside an owning position")

	case *ast.FuncLit:
		// A capturing literal needs an env struct and the {code, env} pair;
		// a non-capturing one needs the thunk of lowering.md §12.2. Neither
		// is built, and nothing distinguishes them here yet.
		fb.l.todoAt(x.Pos(), "function literal")

	case *ast.AwaitExpr, *ast.LaunchExpr:
		fb.l.todoAt(x.Pos(), "await / launch prefix (see async.go)")

	case *ast.MapLit, *ast.ChanConstructor:
		fb.l.todoAt(x.Pos(), "map literal / channel constructor")

	case *ast.EnumShorthand:
		return fb.enumShorthand(x)
	}
	fb.l.bugAt(e.Pos(), "unlowerable expression")
	return Value{}
}

func (fb *funcBuilder) basicLit(x *ast.BasicLit) Value {
	t := fb.typ(fb.l.info.TypeOf(x))
	tv, ok := fb.l.info.Types[x]
	if !ok || tv.Value == nil {
		// The checker records a value for every literal it folded. Missing
		// one means a nil, whose bytes are null.
		if x.Kind == token.NIL {
			return Null()
		}
		fb.l.bugAt(x.Pos(), "literal reached lowering with no folded value")
	}
	switch {
	case IsFloat(t):
		f, _ := types.FloatVal(tv.Value)
		return Float(f, t)
	case t.Kind == KStruct:
		// A string literal: a global byte array plus a constant length,
		// built into a {ptr, len} header at the use site.
		s, _ := types.StringVal(tv.Value)
		return fb.stringHeader(s)
	default:
		n, _ := types.Int64Val(tv.Value)
		return Int(n, t)
	}
}

// stringHeader materializes a literal's header. The bytes are a global; the
// header is a value, so each use builds its own — which costs nothing, since
// both fields are constants.
func (fb *funcBuilder) stringHeader(s string) Value {
	g := fb.l.internString(fb.fn.Module, s)
	st := fb.l.header(fb.fn.Module, "_Vstr", field("p", Ptr), field("len", I64))
	slot := fb.alloca(StructOf(st), "str")
	fb.emitVoid(OpStore, Ptr, fb.fieldPtr(slot, st, 0), GlobalRef(g, Ptr))
	fb.emitVoid(OpStore, I64, fb.fieldPtr(slot, st, 1), Int(int64(len(s)), I64))
	return slot
}

func (fb *funcBuilder) identValue(x *ast.Ident) Value {
	obj := fb.l.info.ObjectOf(x)
	if obj == nil {
		fb.l.bugAt(x.Pos(), "unresolved identifier "+x.Name)
	}
	if c, ok := obj.(*types.Const); ok {
		return fb.l.constValue(c.Val(), fb.typ(c.Type()))
	}
	if b := fb.lookup(obj); b != nil {
		if b.slot && !IsAggregate(b.typ) {
			return fb.emit(OpLoad, b.typ, b.val)
		}
		return b.val
	}
	if fn, ok := obj.(*types.Func); ok {
		// A function named as a value. Only the non-capturing form crosses,
		// and it needs the thunk pair.
		_ = fn
		fb.l.todoAt(x.Pos(), "function used as a value")
	}
	// A package-scope global. globalRef is an address — the same thing
	// address() wants — so a scalar use is that address plus a load, and an
	// aggregate use is the address itself, exactly as a slot binding is.
	ref := fb.globalRef(obj)
	t := fb.typ(obj.Type())
	if IsAggregate(t) {
		return ref
	}
	return fb.emit(OpLoad, t, ref)
}

// binary lowers one binary operator, paying for the gap between vir and
// Vertex: vir wraps on integer overflow and masks shift counts, Vertex traps
// on both.
func (fb *funcBuilder) binary(x *ast.BinaryExpr) Value {
	switch x.Op {
	case token.LAND, token.LOR:
		return fb.shortCircuit(x)
	case token.DOTDOT:
		fb.l.bugAt(x.Pos(), "a range reached expression lowering")
	}

	lhs := fb.expr(x.X)
	rhs := fb.expr(x.Y)
	lt := fb.l.info.TypeOf(x.X)

	if types.IsString(lt) {
		return fb.stringBinary(x, lhs, rhs)
	}

	op := binaryOpFor(x.Op, lhs.Type)
	if op.IsComparison() {
		return fb.emit(op, I1, lhs, rhs)
	}
	return fb.checkedBinary(op, lhs, rhs, x.OpPos)
}

// checkedBinary emits the operation plus its trap. The overflow opcodes
// yield the flag, not the result, so a checked + is two instructions plus a
// branch. Signedness picks between the s and u forms from the operand type,
// exactly as it does for sdiv vs udiv.
func (fb *funcBuilder) checkedBinary(op Op, a, b Value, pos token.Pos) Value {
	var check Op
	switch op {
	case OpAdd:
		check = signed(a.Type, OpSAddO, OpUAddO)
	case OpSub:
		check = signed(a.Type, OpSSubO, OpUSubO)
	case OpMul:
		check = signed(a.Type, OpSMulO, OpUMulO)
	case OpShl, OpLShr, OpAShr:
		// A shift count at or beyond the left operand's width traps.
		width := Int(int64(a.Type.Bits), b.Type)
		ok := fb.emit(OpUlt, I1, b, width)
		fb.trapUnless(ok)
		return fb.emit(op, a.Type, a, b)
	case OpSDiv, OpUDiv, OpSRem, OpURem:
		// Zero and INT_MIN / -1 trap for free in vir.
		return fb.emit(op, a.Type, a, b)
	default:
		return fb.emit(op, a.Type, a, b)
	}

	if IsFloat(a.Type) {
		// IEEE, no trap, no fast-math.
		return fb.emit(op, a.Type, a, b)
	}
	res := fb.emit(op, a.Type, a, b)
	flag := fb.emit(check, I1, a, b)
	fb.trapIf(flag)
	_ = pos
	return res
}

func (fb *funcBuilder) trapIf(cond Value) {
	n := &If{Cond: cond, Then: &Seq{}}
	n.Then.add(&TrapStmt{})
	fb.seq.add(n)
}

func (fb *funcBuilder) trapUnless(cond Value) {
	n := &If{Cond: cond, Then: &Seq{}, Else: &Seq{}}
	n.Else.add(&TrapStmt{})
	fb.seq.add(n)
}

// shortCircuit is a branch, not a select. There is no truthiness, so the
// operand is already a bool.
func (fb *funcBuilder) shortCircuit(x *ast.BinaryExpr) Value {
	slot := fb.alloca(Bool, "sc")
	lhs := fb.expr(x.X)
	fb.emitVoid(OpStore, Bool, slot, lhs)

	n := &If{Cond: fb.narrowBool(lhs), Then: &Seq{}, Else: &Seq{}}
	branch := n.Then
	if x.Op == token.LOR {
		branch = n.Else
	}
	outer := fb.seq
	fb.seq = branch
	rhs := fb.expr(x.Y)
	fb.emitVoid(OpStore, Bool, slot, rhs)
	fb.seq = outer
	fb.seq.add(n)

	return fb.emit(OpLoad, Bool, slot)
}

func (fb *funcBuilder) unary(x *ast.UnaryExpr) Value {
	switch x.Op {
	case token.SUB:
		v := fb.expr(x.X)
		if IsFloat(v.Type) {
			return fb.emit(OpNeg, v.Type, v)
		}
		// Unary minus is checked the same way subtraction is, against (0, a).
		res := fb.emit(OpNeg, v.Type, v)
		flag := fb.emit(signed(v.Type, OpSSubO, OpUSubO), I1, Int(0, v.Type), v)
		fb.trapIf(flag)
		return res

	case token.NOT:
		v := fb.expr(x.X)
		return fb.emit(OpXor, v.Type, v, Int(1, v.Type))

	case token.TILDE:
		v := fb.expr(x.X)
		return fb.emit(OpNot, v.Type, v)

	case token.AND:
		// Address-of on a value, dereference on a typed_ptr — read from the
		// operand's statically written type, which is the one fork the
		// analyzer left open for exactly this reason. It must be asked of
		// the types.Type and not of the lowered type: ownership, map, and
		// chan all lower to ptr and none of them dereferences.
		if _, isPtr := types.Underlying(fb.l.info.TypeOf(x.X)).(*types.Pointer); isPtr {
			p := fb.expr(x.X)
			t := fb.typ(fb.l.info.TypeOf(x))
			if IsAggregate(t) {
				return p
			}
			return fb.emit(OpLoad, t, p)
		}
		return fb.address(x.X)
	}
	fb.l.bugAt(x.Pos(), "unlowerable unary operator")
	return Value{}
}

// cast is `x as T`. vir conversions are destination-explicit, so Instr.Type
// is the destination and the source is read off the argument.
func (fb *funcBuilder) cast(x *ast.CastExpr) Value {
	v := fb.expr(x.X)
	dst := fb.typ(fb.l.info.TypeOf(x))
	return fb.convert(v, dst)
}

func (fb *funcBuilder) convert(v Value, dst *Type) Value {
	src := v.Type
	switch {
	case src == dst:
		return v
	case IsInt(src) && IsInt(dst):
		switch {
		case dst.Bits < src.Bits:
			return fb.emit(OpTrunc, dst, v)
		case dst.Bits > src.Bits && src.Signed:
			return fb.emit(OpSext, dst, v)
		case dst.Bits > src.Bits:
			return fb.emit(OpZext, dst, v)
		}
		return fb.emit(OpBitcast, dst, v)
	case IsInt(src) && IsFloat(dst):
		return fb.emit(signed(src, OpSfromint, OpUfromint), dst, v)
	case IsFloat(src) && IsInt(dst):
		// Traps out of range including NaN and ±Inf, which discharges
		// semantics.md §5.5's last row for free.
		return fb.emit(signed(dst, OpStoint, OpUtoint), dst, v)
	case IsFloat(src) && IsFloat(dst):
		if dst.Bits < src.Bits {
			return fb.emit(OpFdemote, dst, v)
		}
		return fb.emit(OpFpromote, dst, v)
	case IsPtr(src) && IsPtr(dst):
		// typed_ptr T as typed_ptr U, and abstract as typed_ptr T on a
		// memory-flat linkage: both emit nothing at all.
		return v
	case IsPtr(src) || IsPtr(dst):
		return fb.emit(OpBitcast, dst, v)
	}
	// enum -> its discriminant type is a no-op reinterpretation on a unit
	// enum, which already lowered to the discriminant integer.
	return fb.emit(OpBitcast, dst, v)
}

func (fb *funcBuilder) callExpr(x *ast.CallExpr) Value {
	// A reserved builtin is not an ordinary call: each one's shape is
	// arity- and type-argument-specific.
	if id, ok := x.Fun.(*ast.Ident); ok {
		if b, ok := fb.l.info.ObjectOf(id).(*types.Builtin); ok {
			return fb.builtinCall(b, x)
		}
	}

	callee, mod, args := fb.callTarget(x)
	sig := fb.calleeSignature(x)
	if sig == nil {
		fb.l.todoAt(x.Pos(), "call through a function value")
	}

	// An aggregate result is an out-parameter, allocated by the caller.
	if sig.Results().Len() > 0 {
		rt := fb.resultType(sig)
		if IsAggregate(rt) {
			out := fb.alloca(rt, "out")
			args = append([]Value{out}, args...)
			fb.emitCall(callee, mod, Void, args)
			return out
		}
		return fb.emitCall(callee, mod, rt, args)
	}
	fb.emitCall(callee, mod, Void, args)
	return Value{}
}

func (fb *funcBuilder) emitCall(callee, mod string, result *Type, args []Value) Value {
	name := ""
	if !IsVoid(result) {
		name = fb.fn.fresh("r")
	}
	fb.seq.add(&Instr{
		Result: name, Op: OpCall, Type: result,
		Callee: callee, Module: mod, Args: args,
	})
	if name == "" {
		return Value{}
	}
	return Name(name, result)
}

// callTarget resolves a callee to a name plus its arguments, enqueuing the
// instance if this is the first call site to reach it.
//
// Named arguments resolve to positional order here and leave no trace; they
// are not mixed, which the analyzer already guaranteed. Arguments evaluate
// left to right, and vir is sequential, so emission order *is* evaluation
// order.
func (fb *funcBuilder) callTarget(x *ast.CallExpr) (callee, mod string, args []Value) {
	obj := fb.calleeObject(x)
	if obj == nil {
		fb.l.todoAt(x.Pos(), "indirect call")
	}

	targs := fb.typeArgsFor(x, obj)
	target := fb.l.enqueue(obj, targs)
	callee = target.Name
	if target.Module != fb.fn.Module {
		fb.fn.Module.Import(target.Module.Path)
		mod = target.Module.Name
	}

	sig := obj.Signature()
	// A method's receiver is the first argument.
	if sig != nil && sig.Recv() != nil {
		if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
			args = append(args, fb.receiverArg(sel.X, sig.Recv()))
		}
	}
	ordered := fb.orderArgs(x, sig)
	for i, a := range ordered {
		args = append(args, fb.argument(a, sig, i))
	}
	return callee, mod, args
}

// argument applies the calling convention at the site: a mut argument passes
// the address of the caller's slot, an owning one copies or moves, a shared
// one passes the value or a bare pointer.
func (fb *funcBuilder) argument(e ast.Expr, sig *types.Signature, i int) Value {
	if sig == nil || i >= sig.Params().Len() {
		return fb.owning(e)
	}
	p := sig.Params().At(i)
	switch p.Mode() {
	case types.ModeMut:
		return fb.address(e)
	case types.ModeVar:
		return fb.owning(e)
	default:
		v := fb.expr(e)
		return v
	}
}

func (fb *funcBuilder) receiverArg(e ast.Expr, recv *types.Var) Value {
	if recv.Mode() == types.ModeMut || IsAggregate(fb.typ(recv.Type())) {
		return fb.address(e)
	}
	return fb.expr(e)
}

// address computes an lvalue's address, forcing a slot where one is needed.
func (fb *funcBuilder) address(e ast.Expr) Value {
	switch x := e.(type) {
	case *ast.ParenExpr:
		return fb.address(x.X)

	case *ast.Ident:
		obj := fb.l.info.ObjectOf(x)
		if b := fb.lookup(obj); b != nil {
			if b.slot {
				return b.val
			}
			// A `let` has no address by construction. Reaching here means
			// forcesSlot missed a case.
			fb.l.bugAt(x.Pos(), "address of a non-addressable binding "+x.Name)
		}
		return fb.globalRef(obj)

	case *ast.SelectorExpr:
		sel := fb.l.info.SelectionOf(x)
		if sel == nil || sel.Kind != types.FieldVal {
			fb.l.bugAt(x.Pos(), "selector with no recorded field selection")
		}
		base := fb.address(x.X)
		st := fb.typ(fb.l.info.TypeOf(x.X))
		if st.Kind != KStruct {
			fb.l.bugAt(x.Pos(), "field access on a non-struct")
		}
		return fb.fieldPtr(base, st.Struct, sel.Index)

	case *ast.TupleIndexExpr:
		base := fb.address(x.X)
		st := fb.typ(fb.l.info.TypeOf(x.X))
		idx := 0
		for i := 0; i < len(x.Text); i++ {
			idx = idx*10 + int(x.Text[i]-'0')
		}
		return fb.fieldPtr(base, st.Struct, idx)

	case *ast.IndexExpr:
		return fb.elementPtr(x)

	case *ast.UnaryExpr:
		if x.Op == token.AND {
			// A dereference-write: the pointer already is the address.
			return fb.expr(x.X)
		}
	}
	// A temporary. Materializing it into a slot is correct for a shared
	// aggregate argument and wrong for anything the callee mutates, which
	// the exclusivity rules already forbid.
	v := fb.expr(e)
	slot := fb.alloca(v.Type, "tmp")
	fb.emitVoid(OpStore, v.Type, slot, v)
	return slot
}

func (fb *funcBuilder) fieldPtr(base Value, s *Struct, i int) Value {
	return fb.emit(OpFieldPtr, Ptr, base, Name(s.Name, Ptr), Name(s.Fields[i].Name, Ptr))
}

// elementPtr is a bounds check plus index.ptr. The offset is in bytes; the
// front end scales by the element's stride.
func (fb *funcBuilder) elementPtr(x *ast.IndexExpr) Value {
	ct := fb.l.info.TypeOf(x.X)
	et := fb.typ(fb.l.info.TypeOf(x))
	stride := align(fb.l.sizeOf(et), fb.l.alignOf(et))

	base := fb.address(x.X)
	idx := fb.convert(fb.expr(x.Indices[0]), I64)

	switch u := types.Underlying(ct).(type) {
	case *types.Array:
		// todo: a constant subscript provably out of range is a compile
		// error and emits no check; that needs the folded index here.
		ok := fb.emit(OpUlt, I1, idx, Int(u.Len(), I64))
		fb.trapUnless(ok)
	case *types.Slice:
		hdr := fb.l.header(fb.fn.Module, "_Vvec", field("p", Ptr), field("len", I64), field("cap", I64))
		lenp := fb.fieldPtr(base, hdr, 1)
		n := fb.emit(OpLoad, I64, lenp)
		ok := fb.emit(OpUlt, I1, idx, n)
		fb.trapUnless(ok)
		base = fb.emit(OpLoad, Ptr, fb.fieldPtr(base, hdr, 0))
	default:
		fb.l.todoAt(x.Pos(), "subscript of "+types.TypeString(ct))
	}

	off := fb.emit(OpMul, I64, idx, Int(stride, I64))
	return fb.emit(OpIndexPtr, Ptr, base, off)
}

// compositeLit constructs a struct. A class is constructed by calling an
// initializer instead, and that is an ordinary call — the punctuation is
// load-bearing and the tree kept the two distinct.
func (fb *funcBuilder) compositeLit(x *ast.CompositeLit) Value {
	t := fb.typ(fb.l.info.TypeOf(x))
	if t.Kind != KStruct {
		fb.l.bugAt(x.Pos(), "composite literal of non-struct type")
	}
	slot := fb.alloca(t, "lit")

	// Zero first, then write each named field. A field with a default that
	// the literal omits is evaluated at *this* construction — todo: the
	// default expression is on the ast.Field and is not evaluated here, so
	// an omitted defaulted field currently gets the zero value instead.
	fb.zero(slot, t)

	for _, el := range x.Elems {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, _ := kv.Key.(*ast.Ident)
		if key == nil {
			continue
		}
		i := t.Struct.FieldIndex(key.Name)
		if i < 0 {
			fb.l.bugAt(kv.Pos(), "composite literal names an absent field "+key.Name)
		}
		v := fb.owning(kv.Value)
		fb.storeInto(fb.fieldPtr(slot, t.Struct, i), v, t.Struct.Fields[i].Type)
	}
	return slot
}

func (fb *funcBuilder) tupleLit(x *ast.TupleExpr) Value {
	t := fb.typ(fb.l.info.TypeOf(x))
	slot := fb.alloca(t, "tup")
	for i, el := range x.Elems {
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			el = kv.Value
		}
		v := fb.owning(el)
		fb.storeInto(fb.fieldPtr(slot, t.Struct, i), v, t.Struct.Fields[i].Type)
	}
	return slot
}

func (fb *funcBuilder) arrayLit(x *ast.ArrayLit) Value {
	t := fb.typ(fb.l.info.TypeOf(x))
	if t.Kind != KArray {
		fb.l.todoAt(x.Pos(), "slice literal")
	}
	slot := fb.alloca(t, "arr")
	stride := align(fb.l.sizeOf(t.Elem), fb.l.alignOf(t.Elem))
	for i, el := range x.Elems {
		v := fb.owning(el)
		p := fb.emit(OpIndexPtr, Ptr, slot, Int(int64(i)*stride, I64))
		fb.storeInto(p, v, t.Elem)
	}
	return slot
}

func (fb *funcBuilder) enumShorthand(x *ast.EnumShorthand) Value {
	e := types.AsEnum(fb.l.info.TypeOf(x))
	if e == nil {
		fb.l.bugAt(x.Pos(), "enum shorthand with no enum type")
	}
	v := e.LookupVariant(x.Name.Name)
	if v == nil {
		fb.l.bugAt(x.Pos(), "enum shorthand names an absent variant")
	}
	if len(x.Args) > 0 {
		// Payload construction writes the tag plus the variant's fields into
		// the payload bytes.
		fb.l.todoAt(x.Pos(), "payload enum construction")
	}
	if e.UnitOnly() {
		return Int(v.Value, fb.l.basic(e.Discriminant()))
	}
	t := fb.typ(fb.l.info.TypeOf(x))
	slot := fb.alloca(t, "enum")
	fb.zero(slot, t)
	fb.emitVoid(OpStore, t.Struct.Fields[0].Type,
		fb.fieldPtr(slot, t.Struct, 0), Int(v.Value, t.Struct.Fields[0].Type))
	return slot
}

// ------------------------------------------------------------------ helpers

// narrowBool turns an i8 bool into the i1 a branch wants. Where the value
// came from a comparison it is already i1 and this is a no-op.
func (fb *funcBuilder) narrowBool(v Value) Value {
	if v.Type == I1 {
		return v
	}
	return fb.emit(OpNe, I1, v, Int(0, v.Type))
}

func (fb *funcBuilder) zeroValue(t *Type) Value {
	if IsFloat(t) {
		return Float(0, t)
	}
	if IsPtr(t) {
		return Null()
	}
	return Int(0, t)
}

func (fb *funcBuilder) zero(p Value, t *Type) {
	fb.emitVoid(OpMemset, Void, p, Int(0, U8), Int(fb.l.sizeOf(t), I64))
}

// storeInto writes v at p. Copying an aggregate is a memcopy, never an
// assignment — vir aggregates are memory-only.
func (fb *funcBuilder) storeInto(p Value, v Value, t *Type) {
	if v.IsZero() {
		fb.zero(p, t)
		return
	}
	if IsAggregate(t) {
		fb.emitVoid(OpMemcopy, Void, p, v, Int(fb.l.sizeOf(t), I64))
		return
	}
	fb.emitVoid(OpStore, t, p, v)
}

func (fb *funcBuilder) bindingOf(e ast.Expr) *binding {
	switch x := e.(type) {
	case *ast.ParenExpr:
		return fb.bindingOf(x.X)
	case *ast.Ident:
		return fb.lookup(fb.l.info.ObjectOf(x))
	case *ast.SelectorExpr:
		// A field path kills its root binding: the whole value is what
		// leaves, and there is nowhere for a hole to go otherwise.
		return fb.bindingOf(x.X)
	}
	return nil
}

func signed(t *Type, s, u Op) Op {
	if t != nil && t.Signed {
		return s
	}
	return u
}

func cmpOp(base Op, t *Type) Op {
	if IsFloat(t) {
		switch base {
		case OpSlt:
			return OpFlt
		case OpSgt:
			return OpFgt
		case OpSle:
			return OpFle
		case OpSge:
			return OpFge
		}
	}
	if t != nil && !t.Signed {
		switch base {
		case OpSlt:
			return OpUlt
		case OpSgt:
			return OpUgt
		case OpSle:
			return OpUle
		case OpSge:
			return OpUge
		}
	}
	return base
}

func binaryOpFor(k token.Kind, t *Type) Op {
	switch k {
	case token.ADD, token.WRAP_ADD:
		return OpAdd
	case token.SUB, token.WRAP_SUB:
		return OpSub
	case token.MUL, token.WRAP_MUL:
		return OpMul
	case token.QUO:
		// Float division has no vir opcode. vir's opTable gives sdiv
		// ConstraintInt and §4.1 lists no float division at all, so sdiv.f32
		// is emitted faithfully and ir/verify rejects it. Either vir needs
		// fdiv or this needs to stop claiming one exists; the fix is not
		// here, and papering over it would be a decision this package does
		// not get to make.
		return signed(t, OpSDiv, OpUDiv)
	case token.REM:
		return signed(t, OpSRem, OpURem)
	case token.AND:
		return OpAnd
	case token.OR:
		return OpOr
	case token.XOR:
		return OpXor
	case token.SHL:
		return OpShl
	case token.SHR:
		if t != nil && t.Signed {
			return OpAShr
		}
		return OpLShr
	case token.EQL, token.IDENTICAL:
		return OpEq
	case token.NEQ, token.NOT_IDENTICAL:
		return OpNe
	case token.LSS:
		return cmpOp(OpSlt, t)
	case token.GTR:
		return cmpOp(OpSgt, t)
	case token.LEQ:
		return cmpOp(OpSle, t)
	case token.GEQ:
		return cmpOp(OpSge, t)
	}
	return OpInvalid
}

func compoundBase(k token.Kind) token.Kind {
	switch k {
	case token.ADD_ASSIGN:
		return token.ADD
	case token.SUB_ASSIGN:
		return token.SUB
	case token.MUL_ASSIGN:
		return token.MUL
	case token.QUO_ASSIGN:
		return token.QUO
	case token.REM_ASSIGN:
		return token.REM
	case token.AND_ASSIGN:
		return token.AND
	case token.OR_ASSIGN:
		return token.OR
	case token.XOR_ASSIGN:
		return token.XOR
	case token.SHL_ASSIGN:
		return token.SHL
	case token.SHR_ASSIGN:
		return token.SHR
	}
	return token.INVALID
}