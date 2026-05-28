// checker.go
package compiler

import (
	"strconv"
	"strings"

	"github.com/vertex-language/vertex/parser"
)

type checker struct {
	pkg    *Package
	scope  *Scope
	tags   *BuildTags
	errors ErrorList
}

func newChecker(pkg *Package, defaultTags *BuildTags) *checker {
	return &checker{
		pkg:   pkg,
		scope: NewScope(nil),
		tags:  defaultTags,
	}
}

func (c *checker) Run() *Scope {
	for _, sf := range c.pkg.Files {
		tags := c.tags
		if len(sf.BuildTags.Tags) > 0 {
			tags = sf.BuildTags
		}
		for _, tld := range sf.Tree.AllTopLevelDecl() {
			c.collectTopLevel(tld, sf, tags)
		}
	}
	return c.scope
}

func (c *checker) collectTopLevel(ctx parser.ITopLevelDeclContext, sf *SourceFile, tags *BuildTags) {
	switch {
	case ctx.StructDecl() != nil:
		c.collectStruct(ctx.StructDecl())

	case ctx.ClassDecl() != nil:
		cd := ctx.ClassDecl()
		if cd.COLON() != nil {
			c.collectNativeClass(cd, sf, tags)
		} else {
			c.collectClass(cd)
		}

	case ctx.EnumDecl() != nil:
		c.collectEnum(ctx.EnumDecl())

	case ctx.TypeAliasDecl() != nil:
		c.collectTypeAlias(ctx.TypeAliasDecl())

	case ctx.FuncDecl() != nil:
		c.collectFuncDecl(ctx.FuncDecl())

	case ctx.VarDeclStmt() != nil:
		// Top-level var/let — treat as global (simplified)
	}
}

// ── Struct ───────────────────────────────────────────────────────────────────

func (c *checker) collectStruct(ctx parser.IStructDeclContext) {
	name := ctx.ID().GetText()
	st := &StructType{Name: name}
	for _, sf := range ctx.AllStructField() {
		// ← 2.0: fields have no let/var keyword; mutability is determined by
		// the binding at the declaration site, not the field definition.
		f := &StructField{
			Name: sf.ID().GetText(),
			Type: c.resolveType(sf.Type_()),
		}
		st.Fields = append(st.Fields, f)
	}
	st.Size = LayoutStruct(st.Fields)
	c.scope.Define(&Symbol{Name: name, Kind: SymType, Type: st})
}

// ── Class (regular) ──────────────────────────────────────────────────────────

func (c *checker) collectClass(ctx parser.IClassDeclContext) {
	name := ctx.AllID()[0].GetText()
	ct := &ClassType{Name: name}
	for _, member := range ctx.AllClassMember() {
		if cf := member.ClassField(); cf != nil {
			// ← 2.0: class fields also have no let/var keyword.
			f := &StructField{
				Name: cf.ID().GetText(),
				Type: c.resolveType(cf.Type_()),
			}
			ct.Fields = append(ct.Fields, f)
		}
	}
	ct.Size = LayoutStruct(ct.Fields)
	c.scope.Define(&Symbol{Name: name, Kind: SymType, Type: ct})
}

// ── Native class ─────────────────────────────────────────────────────────────

func (c *checker) collectNativeClass(ctx parser.IClassDeclContext, sf *SourceFile, tags *BuildTags) {
	ids := ctx.AllID()
	className := ids[0].GetText()
	parentName := ids[1].GetText()

	module := parentName
	for _, imp := range sf.Imports {
		if imp.IsNative() && imp.Namespace == parentName {
			module = imp.WasmModule(tags)
			break
		}
	}
	_ = module

	ct := &ClassType{Name: className, Native: true, Parent: parentName}
	c.scope.Define(&Symbol{Name: className, Kind: SymType, Type: ct})

	for _, member := range ctx.AllClassMember() {
		nfd := member.NativeFuncDecl()
		if nfd == nil {
			continue
		}
		methodName := nfd.ID().GetText()
		sig := c.nativeFuncSig(nfd)
		ft := &FuncType{Sig: sig}

		c.scope.Define(&Symbol{Name: methodName, Kind: SymNative, Type: ft})
		c.scope.Define(&Symbol{Name: className + "." + methodName, Kind: SymNative, Type: ft})
	}
}

func (c *checker) nativeFuncSig(ctx parser.INativeFuncDeclContext) *FuncSig {
	var params []Type
	var muts []bool
	if pl := ctx.NativeParamList(); pl != nil {
		for _, np := range pl.AllNativeParam() {
			if np.ELLIPSIS() != nil {
				continue
			}
			params = append(params, c.resolveType(np.Type_()))
			muts = append(muts, np.MUT() != nil)
		}
	}
	var ret Type = Void
	if rt := ctx.ReturnType(); rt != nil && rt.Type_() != nil {
		ret = c.resolveType(rt.Type_())
	}
	return &FuncSig{Params: params, Ret: ret, Muts: muts}
}

// ── NativeEntry / CollectNativeFuncs ─────────────────────────────────────────

type NativeEntry struct {
	WasmModule string
	ImportName string
	Sig        *FuncSig
}

func CollectNativeFuncs(
	ctx parser.IClassDeclContext,
	module string,
	scope *Scope,
	resolve func(parser.ITypeContext) Type,
) map[string]*NativeEntry {
	out := make(map[string]*NativeEntry)
	for _, member := range ctx.AllClassMember() {
		nfd := member.NativeFuncDecl()
		if nfd == nil {
			continue
		}
		name := nfd.ID().GetText()
		var params []Type
		var muts []bool

		if pl := nfd.NativeParamList(); pl != nil {
			for _, np := range pl.AllNativeParam() {
				if np.ELLIPSIS() != nil {
					continue
				}
				t := resolve(np.Type_())
				params = append(params, t)
				muts = append(muts, np.MUT() != nil)
			}
		}

		var ret Type = Void
		if rt := nfd.ReturnType(); rt != nil && rt.Type_() != nil {
			ret = resolve(rt.Type_())
		}

		sig := &FuncSig{Params: params, Ret: ret, Muts: muts}
		out[name] = &NativeEntry{
			WasmModule: module,
			ImportName: BuildImportName(name, params, ret),
			Sig:        sig,
		}
	}
	return out
}

// ── Enum ─────────────────────────────────────────────────────────────────────

func (c *checker) collectEnum(ctx parser.IEnumDeclContext) {
	name := ctx.ID().GetText()
	et := &EnumType{Name: name, RawKind: KindInt}

	if ctx.EnumRawType() != nil && ctx.EnumRawType().STRING() != nil {
		et.RawKind = KindString
	}

	var nextInt int64
	for _, ecd := range ctx.AllEnumCaseDecl() {
		for _, ec := range ecd.AllEnumCase() {
			cs := &EnumCase{Name: ec.ID().GetText()}
			if ec.Literal() != nil {
				lit := ec.Literal()
				if lit.DEC_INT_LIT() != nil {
					v, _ := strconv.ParseInt(lit.DEC_INT_LIT().GetText(), 10, 64)
					cs.IntVal = v
					nextInt = v + 1
				} else if lit.STRING_LIT() != nil {
					cs.StrVal = strings.Trim(lit.STRING_LIT().GetText(), `"`)
				}
			} else {
				cs.IntVal = nextInt
				nextInt++
			}
			et.Cases = append(et.Cases, cs)
		}
	}

	c.scope.Define(&Symbol{Name: name, Kind: SymType, Type: et})
}

// ── Type alias ───────────────────────────────────────────────────────────────

func (c *checker) collectTypeAlias(ctx parser.ITypeAliasDeclContext) {
	name := ctx.ID().GetText()
	underlying := c.resolveType(ctx.Type_())
	c.scope.Define(&Symbol{Name: name, Kind: SymType, Type: underlying})
}

// ── Function declaration ─────────────────────────────────────────────────────

// collectFuncDecl registers the function in the global scope.
//
// §26 — Associated functions / receivers: if a receiver block is present,
// the function is also registered under "TypeName.funcName" so dot-call syntax works.
func (c *checker) collectFuncDecl(ctx parser.IFuncDeclContext) {
	name := ctx.ID().GetText()
	sig := c.funcDeclSig(ctx)
	ft := &FuncType{Sig: sig}
	c.scope.Define(&Symbol{Name: name, Kind: SymFunc, Type: ft})

	if rec := ctx.Receiver(); rec != nil {
		typeName := receiverTypeName(c.resolveType(rec.Type_()))
		if typeName != "" {
			c.scope.Define(&Symbol{
				Name: typeName + "." + name,
				Kind: SymFunc,
				Type: ft,
			})
		}
	}
}

func receiverTypeName(t Type) string {
	switch v := t.(type) {
	case *StructType:
		return v.Name
	case *ClassType:
		if !v.Native {
			return v.Name
		}
	}
	return ""
}

// ── Signature helpers ─────────────────────────────────────────────────────────

func (c *checker) funcDeclSig(ctx parser.IFuncDeclContext) *FuncSig {
	var params []Type
	var muts []bool

	// ← 2.0: Check for an isolated receiver block and prepend it to the parameter list.
	if rec := ctx.Receiver(); rec != nil {
		params = append(params, c.resolveType(rec.Type_()))
		muts = append(muts, rec.MUT() != nil)
	}

	if pl := ctx.ParamList(); pl != nil {
		for _, p := range pl.AllParam() {
			params = append(params, c.resolveType(p.Type_()))
			muts = append(muts, p.MUT() != nil)
		}
	}
	var ret Type = Void
	if rt := ctx.ReturnType(); rt != nil && rt.Type_() != nil {
		ret = c.resolveType(rt.Type_())
	}
	qual := c.resolveFuncQual(ctx.FuncQualifier())
	return &FuncSig{Params: params, Ret: ret, Muts: muts, Qual: qual}
}

func (c *checker) resolveFuncQual(fq parser.IFuncQualifierContext) FuncQual {
	if fq == nil {
		return QualNone
	}
	switch {
	case fq.ASYNC() != nil:
		return QualAsync
	case fq.THREAD() != nil:
		return QualThread
	case fq.PROCESS() != nil:
		return QualProcess
	case fq.GPU() != nil:
		return QualCUDA
	}
	return QualNone
}

// ── Type resolution ───────────────────────────────────────────────────────────

func (c *checker) resolveType(ctx parser.ITypeContext) Type {
	return ResolveType(ctx, c.scope)
}

// ResolveType converts a parse-tree TypeContext into a Type.
func ResolveType(ctx parser.ITypeContext, scope *Scope) Type {
	if ctx == nil {
		return Void
	}
	if pt := ctx.PrimitiveType(); pt != nil {
		return resolvePrimitive(pt)
	}
	if ctx.ANY() != nil {
		subTypes := ctx.AllType_()
		var elem Type = Void
		if len(subTypes) > 0 {
			elem = ResolveType(subTypes[0], scope)
		}
		if ctx.OPAQUE() != nil {
			elem = Opaque
		}
		return &PointerType{Elem: elem, Mutable: ctx.MUT() != nil}
	}
	if ctx.OPAQUE() != nil {
		return Opaque
	}
	if ctx.LBRACKET() != nil && ctx.COLON() == nil {
		subs := ctx.AllType_()
		if len(subs) > 0 {
			return &ArrayType{Elem: ResolveType(subs[0], scope)}
		}
		return &ArrayType{Elem: Void}
	}
	if ctx.LBRACKET() != nil && ctx.COLON() != nil {
		subs := ctx.AllType_()
		if len(subs) >= 2 {
			return &ArrayType{Elem: ResolveType(subs[1], scope)}
		}
	}
	if ctx.QUESTION() != nil {
		subs := ctx.AllType_()
		if len(subs) > 0 {
			return &OptionalType{Elem: ResolveType(subs[0], scope)}
		}
	}
	if ctx.RESULT() != nil {
		subs := ctx.AllType_()
		ok := Void
		errT := Void
		if len(subs) > 0 {
			ok = ResolveType(subs[0], scope)
		}
		if len(subs) > 1 {
			errT = ResolveType(subs[1], scope)
		}
		return &ResultType{Ok: ok, Err: errT}
	}
	if ctx.CHAN() != nil {
		subs := ctx.AllType_()
		if len(subs) > 0 {
			return &ChannelType{Elem: ResolveType(subs[0], scope)}
		}
		return &ChannelType{Elem: Void}
	}
	if ctx.ID() != nil {
		name := ctx.ID().GetText()
		if scope != nil {
			if sym := scope.Lookup(name); sym != nil && sym.Kind == SymType {
				return sym.Type
			}
		}
		return &NamedType{Name: name}
	}
	return Void
}

func resolvePrimitive(ctx parser.IPrimitiveTypeContext) Type {
	switch {
	case ctx.INT() != nil:
		return Int
	case ctx.INT8() != nil:
		return Int8
	case ctx.INT16() != nil:
		return Int16
	case ctx.INT32() != nil:
		return Int32
	case ctx.INT64() != nil:
		return Int64
	case ctx.UINT() != nil:
		return Uint
	case ctx.UINT8() != nil:
		return Uint8
	case ctx.UINT16() != nil:
		return Uint16
	case ctx.UINT32() != nil:
		return Uint32
	case ctx.UINT64() != nil:
		return Uint64
	case ctx.FLOAT() != nil:
		return Float
	case ctx.DOUBLE() != nil:
		return Double
	case ctx.BOOL() != nil:
		return Bool
	case ctx.STRING() != nil:
		return StrType
	case ctx.CHAR() != nil:
		return Char
	case ctx.VOID() != nil:
		return Void
	}
	return Void
}