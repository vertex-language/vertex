package analyzer

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// ------------------------------------------------- phase 1: collect names

// collectObjects inserts every package-scope name with a nil type.
//
// The nil type is the point. A.2 ⊢ "top-level declarations are order-
// independent: a declaration may refer to any other declaration in the same
// package regardless of textual position. There is no forward-declaration form
// because none is needed." Inserting every name before resolving any type is
// what makes that true without one.
func (c *Checker) collectObjects() {
	for _, f := range c.files {
		fileScope := types.NewScope(c.pkg.Scope(), "file "+f.Filename(c.conf.Fset))
		fileScope.SetExtent(f.Pos(), f.End())
		c.info.RecordScope(f, fileScope)

		c.collectImports(f, fileScope)

		for _, d := range f.Decls {
			c.collectDecl(d, f, fileScope)
		}
	}
	c.collectMethods()
}

// collectImports binds each import's qualifier in file scope.
//
// A.2.3 ⊢ "the imported package's declared name (its PackageClause) is the
// qualifier under which its symbols are reached; the import path is a locator,
// not a name." That is why this cannot be answered from the path alone, and why
// the importer must have read the imported directory's package clause before
// this runs — the same two-pass shape parser.ParseDir already uses.
//
// The qualifier is file-scoped while the declarations that use it are
// package-scoped, so one file's import never resolves a name in another's.
func (c *Checker) collectImports(f *ast.File, fileScope *types.Scope) {
	for _, imp := range f.Imports {
		for _, p := range imp.Paths {
			path, ok := unquote(p.Value)
			if !ok {
				continue
			}
			if c.conf.Importer == nil {
				// A single-package check has no way to resolve a path. The name
				// is genuinely undeclared from this checker's position, so that
				// is what it reports.
				c.errorExpr(p, diag.UndeclaredName, path)
				continue
			}
			pkg, err := c.conf.Importer.Import(path)
			if err != nil || pkg == nil {
				c.errorExpr(p, diag.UndeclaredName, path)
				continue
			}
			// A.2.3 ⊢ "there is no aliasing form, no dot-import, and no blank
			// import", so there is exactly one name to bind and no choice about
			// what it is.
			name := types.NewPkgName(p.Pos(), c.pkg, pkg, path)
			if alt := fileScope.Insert(name); alt != nil {
				d := diag.New(diag.DuplicateDeclaration, p.Pos(), p.End(), pkg.Name())
				if alt.Pos().IsValid() {
					d.WithNote(alt.Pos(), alt.Pos(),
						"%s is also the qualifier of an earlier import", pkg.Name())
				}
				c.report(d)
			}
		}
	}
}

func (c *Checker) collectDecl(d ast.Decl, f *ast.File, fileScope *types.Scope) {
	switch x := d.(type) {
	case *ast.FuncDecl:
		// A method is collected in collectMethods, after every type name
		// exists. A.6.3 ⊢ "class methods are declared outside the class body",
		// so a receiver may name a type declared later or in another file.
		if x.Recv != nil {
			return
		}
		obj := types.NewFunc(x.Name.Pos(), c.pkg, x.Name.Name, nil)
		c.declare(c.pkg.Scope(), x.Name, obj)
		c.objMap[obj] = &declInfo{decl: d, file: f, fileScope: fileScope}

	case *ast.RecordDecl:
		obj := types.NewTypeName(x.Name.Pos(), c.pkg, x.Name.Name, nil)
		c.declare(c.pkg.Scope(), x.Name, obj)
		c.objMap[obj] = &declInfo{decl: d, file: f, fileScope: fileScope}

	case *ast.EnumDecl:
		obj := types.NewTypeName(x.Name.Pos(), c.pkg, x.Name.Name, nil)
		c.declare(c.pkg.Scope(), x.Name, obj)
		c.objMap[obj] = &declInfo{decl: d, file: f, fileScope: fileScope}

	case *ast.TypeAliasDecl:
		obj := types.NewTypeName(x.Name.Pos(), c.pkg, x.Name.Name, nil)
		c.declare(c.pkg.Scope(), x.Name, obj)
		c.objMap[obj] = &declInfo{decl: d, file: f, fileScope: fileScope}

	case *ast.ConstraintDecl:
		// A constraint gets a TypeName whose Constraint is non-nil from the
		// start, so IsConstraint answers correctly even before phase 2 fills in
		// the set — which is what rejects `var c: Ordered` (A.14) no matter
		// which order the two declarations are checked in.
		obj := types.NewConstraintName(x.Name.Pos(), c.pkg, x.Name.Name, types.Any)
		c.declare(c.pkg.Scope(), x.Name, obj)
		c.objMap[obj] = &declInfo{decl: d, file: f, fileScope: fileScope}

	case *ast.VarDecl:
		// A.2 ⊢ a top-level VariableDeclaration "must have a compile-time-
		// evaluable initializer; there is no static initialization order and no
		// initialization-time code." Evaluating it needs the constant
		// evaluator, so the object is collected here and the value is
		// typecheck's.
		for _, b := range x.Bindings {
			obj := types.NewVar(b.Name.Pos(), c.pkg, b.Name.Name, nil)
			obj.SetMutable(x.Kw == token.VAR)
			c.declare(c.pkg.Scope(), b.Name, obj)
			c.objMap[obj] = &declInfo{decl: d, file: f, fileScope: fileScope, node: b}
		}

	case *ast.DeclareDecl:
		c.collectForeign(x, f, fileScope)
	}
}

// collectForeign collects a declare block's members and checks the block's own
// legality (A.8.1, A.8.2).
//
// A.8.1 ⊢ "a declare block is a linkage boundary, not a namespace: symbols
// declared inside it are injected into the file's current package." So its
// members land in package scope rather than a nested one, and a name collision
// with an ordinary declaration is an ordinary redeclaration.
func (c *Checker) collectForeign(d *ast.DeclareDecl, f *ast.File, fileScope *types.Scope) {
	// A.8.1 ⊢ "a DeclareDeclaration is legal only in a file carrying a
	// BuildClause. The build tag picks the ABI family; the block keyword and
	// the member shapes pick the convention within it."
	if f.Build == nil {
		c.errorAt(d.Declare, d.KindPos, diag.DeclareNoTag)
	}

	if d.Kind == token.CtxFramework {
		// A.8.1 ⊢ `declare framework` "names a platform-bundled, versioned
		// library and is legal only where the target platform has a
		// first-class notion of one."
		if !token.HasFrameworks(c.conf.Tag) {
			c.errorAt(d.KindPos, d.KindPos, diag.NoFrameworks, c.conf.Tag.String())
		}
		// A.8.2 ✗ `declare framework["windows", "com"] "SomeLib" { }` —
		// ⊢ "bundled message-passing linkage has exactly one convention by
		// design, and unlike a C++ ABI it does not fork by compiler, standard
		// library, or flag — which is precisely why it is safe to leave
		// silent." The tagged form parses so this can name the rule.
		if d.Variant != nil {
			c.errorAt(d.Variant.Pos(), d.Variant.End(), diag.FrameworkTag)
		}
	}

	family := c.familyForBlock(d)

	for _, m := range d.Members {
		switch x := m.(type) {
		case *ast.ForeignFunc:
			// A.8.3 ⊢ "init is a prefix modifier on func, not a function name.
			// The unnamed form is what bare Type(...) construction resolves
			// to" — so it has no package-scope name at all and is reached only
			// through its enclosing class.
			if x.Name == nil {
				continue
			}
			obj := types.NewFunc(x.Name.Pos(), c.pkg, x.Name.Name, nil)
			c.declare(c.pkg.Scope(), x.Name, obj)
			c.objMap[obj] = &declInfo{
				decl: d, file: f, fileScope: fileScope, node: x, family: family,
			}

		case *ast.ForeignClass:
			obj := types.NewTypeName(x.Name.Pos(), c.pkg, x.Name.Name, nil)
			c.declare(c.pkg.Scope(), x.Name, obj)
			c.objMap[obj] = &declInfo{
				decl: d, file: f, fileScope: fileScope, node: x, family: family,
			}

		case *ast.ForeignField:
			// A.8.3 ✗ "fields are banned. A declare block describes call shape
			// only, never foreign-side layout. This is what keeps the question
			// 'which C++ ABI, exactly?' out of the type system."
			c.errorExpr(x.Name, diag.ForeignField, x.Name.Name)
		}
	}
}

// familyForBlock decides the import family a declare block's handles belong to
// (A.4.4).
//
// The family is what makes `abstract` → `typed_ptr T` legal or not: ⊢ it "is
// legal only for a handle minted by a memory-flat import family (C, WASM). It
// is a compile error for an object-graph family (Objective-C, JS), whose
// handles have no byte representation to point at."
//
// The build tag alone cannot answer it on every target. Under darwin a
// `declare module` is flat C while a `declare framework` is an Objective-C
// object graph, so the block keyword refines what the tag started.
func (c *Checker) familyForBlock(d *ast.DeclareDecl) types.Family {
	if d.Kind == token.CtxFramework {
		return types.FamilyObjectGraph
	}
	if c.conf.Tag == token.TagJS {
		return types.FamilyObjectGraph
	}
	return types.FamilyMemoryFlat
}

// collectMethods attaches each method to its receiver's Named type.
//
// It runs after every type name exists, for the same reason phase 1 exists at
// all: a receiver may name a type declared later in the file or in another file
// of the package, and A.2 makes that legal.
func (c *Checker) collectMethods() {
	for _, f := range c.files {
		fileScope := c.fileScopeOf(f)
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Recv == nil {
				continue
			}

			// A.1.4 ⊢ a ReservedBuiltinName "may not be declared as a member,
			// method, or field name."
			// ✗ func (w: var Widget) new() { }
			// ✗ func (p: typed_ptr int32) addr() { }
			if types.Reserved(fd.Name.Name) {
				c.errorExpr(fd.Name, diag.ReservedAsMember, fd.Name.Name)
				continue
			}

			// A.7.6 ⊢ "a method may not declare a type parameter of its own.
			// Everything a method is generic over comes from its receiver
			// type." The parser accepts the list either way so this can name
			// the rule rather than reporting a syntax error at the bracket.
			if fd.TypeParams != nil {
				c.errorAt(fd.TypeParams.Pos(), fd.TypeParams.End(),
					diag.MethodTypeParams, fd.Name.Name)
			}

			obj := types.NewFunc(fd.Name.Pos(), c.pkg, fd.Name.Name, nil)
			c.objMap[obj] = &declInfo{decl: fd, file: f, fileScope: fileScope}
			c.info.RecordDef(fd.Name, obj)

			named := c.receiverBase(fd.Recv, fileScope)
			if named == nil {
				continue
			}
			if prev := named.LookupMethod(fd.Name.Name); prev != nil {
				d := diag.New(diag.DuplicateDeclaration,
					fd.Name.Pos(), fd.Name.End(), fd.Name.Name)
				if prev.Pos().IsValid() {
					d.WithNote(prev.Pos(), prev.Pos(),
						"previous declaration of %s", fd.Name.Name)
				}
				c.report(d)
				continue
			}
			named.AddMethod(obj)
		}
	}
}

// receiverBase finds the Named type a receiver declares a method on.
//
// It deliberately does not resolve the receiver's full type. A.7.6 ⊢ "a method
// receiver re-declares the type's parameter list to bring the names into scope.
// The receiver's [T] binds the name; it does not introduce a fresh one" — so
// the base name is what associates the method, and the parameter list is
// handled later when the signature is built.
func (c *Checker) receiverBase(r *ast.Receiver, fileScope *types.Scope) *types.Named {
	e := stripParens(r.Type)

	// A.6.1's ReceiverType admits mut/var/shared before the type name.
	if own, ok := e.(*ast.OwnershipType); ok {
		e = stripParens(own.X)
	}
	// ...and A.7.6's parameter list arrives as an IndexExpr.
	if ix, ok := e.(*ast.IndexExpr); ok {
		e = stripParens(ix.X)
	}

	id, ok := e.(*ast.Ident)
	if !ok {
		c.errorExpr(r.Type, diag.NotAType, exprString(r.Type))
		return nil
	}

	obj := c.pkg.Scope().Lookup(id.Name)
	if obj == nil {
		// A method on a type from another package is not a thing: A.6.1's
		// ReceiverType is a TypeName, not a QualifiedTypeName.
		c.errorExpr(id, diag.UndeclaredName, id.Name)
		return nil
	}
	tn, ok := obj.(*types.TypeName)
	if !ok || tn.IsConstraint() {
		c.errorExpr(r.Type, diag.NotAType, id.Name)
		return nil
	}
	c.info.RecordUse(id, tn)

	saved := c.scope
	c.scope = fileScope
	c.objDecl(tn)
	c.scope = saved

	return types.AsNamed(tn.Type())
}

func (c *Checker) fileScopeOf(f *ast.File) *types.Scope {
	if s, ok := c.info.Scopes[f]; ok {
		return s
	}
	return c.pkg.Scope()
}

// ----------------------------------------------- phase 2: resolve types

func (c *Checker) resolveDeclTypes() {
	// Iteration order over objMap is unspecified, which is fine: objDecl is
	// idempotent and resolves dependencies on demand, so the result does not
	// depend on the order. Diagnostics are sorted by position by the caller.
	for obj := range c.objMap {
		c.objDecl(obj)
	}
}

// objDecl resolves one object's type, on demand and at most once.
//
// The resolving stack is what turns a self-referential declaration into a
// diagnostic instead of a stack overflow. A.2's order-independence means a
// cycle is reachable from any entry point, so the guard lives here rather than
// in a pre-pass over the dependency graph.
func (c *Checker) objDecl(obj types.Object) {
	if obj == nil || obj.Type() != nil {
		return
	}
	d := c.objMap[obj]
	if d == nil {
		return
	}

	for _, o := range c.resolving {
		if o != obj {
			continue
		}
		c.errorAt(obj.Pos(), obj.Pos(), diag.TypeCycle, obj.Name())
		obj.SetType(c.invalid())
		return
	}
	c.resolving = append(c.resolving, obj)
	defer func() { c.resolving = c.resolving[:len(c.resolving)-1] }()

	// Resolution runs in the declaring file's scope, so a qualified type
	// resolves through that file's imports and no other's.
	saved := c.scope
	c.scope = d.fileScope
	defer func() { c.scope = saved }()

	switch x := d.decl.(type) {
	case *ast.RecordDecl:
		c.recordDecl(obj.(*types.TypeName), x)
	case *ast.EnumDecl:
		c.enumDecl(obj.(*types.TypeName), x)
	case *ast.TypeAliasDecl:
		c.aliasDecl(obj.(*types.TypeName), x)
	case *ast.ConstraintDecl:
		c.constraintDecl(obj.(*types.TypeName), x)
	case *ast.FuncDecl:
		c.funcDecl(obj.(*types.Func), x)
	case *ast.VarDecl:
		c.varDecl(obj.(*types.Var), x, d.node)
	case *ast.DeclareDecl:
		c.foreignDecl(obj, d)
	}
}

// recordDecl resolves a StructDeclaration or ClassDeclaration (A.6.2, A.6.3).
//
// One path for both, matching ast.RecordDecl: A.6.3 ⊢ "a class is byte-for-byte
// identical in layout to a struct. It differs in its member and method model —
// initializers, teardown, receiver conventions, identity comparison — not in
// where its bytes live." The keyword is carried into Struct.class, and the two
// consumers that care read it there: construction syntax (A.4.7 ⊢ "a struct is
// built with a brace literal, a class is built by calling its init") and `===`
// identity comparison (A.4.5 ⊢ "legal on classes only").
//
// The Named is created and bound before any field resolves. That two-step is
// what lets `struct Node { next: typed_ptr Node }` reach its own name without
// tripping the cycle guard — the guard fires on an object whose *type* is still
// nil, and this one's is not.
func (c *Checker) recordDecl(obj *types.TypeName, d *ast.RecordDecl) {
	named := types.NewNamed(obj, nil)
	obj.SetType(named)

	tparams := c.typeParams(d.TypeParams)
	named.SetTypeParams(tparams)

	if len(tparams) > 0 {
		c.openScope(d, "type parameters of "+d.Name.Name)
		defer c.closeScope()
		c.declareTypeParams(d.TypeParams, tparams)
	}

	seen := make(map[string]bool, len(d.Fields))
	fields := make([]*types.Field, 0, len(d.Fields))

	for _, f := range d.Fields {
		name := f.Name.Name

		// A.1.4 ⊢ a ReservedBuiltinName may not be a field name.
		if types.Reserved(name) {
			c.errorExpr(f.Name, diag.ReservedAsMember, name)
			continue
		}
		if !f.Name.IsBlank() {
			if seen[name] {
				c.errorExpr(f.Name, diag.DuplicateField, name, d.Name.Name)
				continue
			}
			seen[name] = true
		}

		ft := c.typ(f.Type)
		fields = append(fields, &types.Field{
			Name: name,
			Type: ft,
			// A.6.2 ⊢ "a field default is evaluated at construction for any
			// field the literal omits." Evaluating it is typecheck's; recording
			// that one exists is enough for layout and for construction checks.
			HasDefault: f.Default != nil,
		})
		c.info.RecordDef(f.Name, types.NewField(f.Name.Pos(), c.pkg, name, ft))
	}

	// A.6.2 ⊢ fields are laid out "in declaration order with ABI padding", so
	// the slice order is the layout order and nothing reorders it.
	named.SetUnderlying(types.NewStruct(fields, d.Kw == token.CLASS))
}

// enumDecl resolves an EnumDeclaration (A.6.5).
func (c *Checker) enumDecl(obj *types.TypeName, d *ast.EnumDecl) {
	named := types.NewNamed(obj, nil)
	obj.SetType(named)

	tparams := c.typeParams(d.TypeParams)
	named.SetTypeParams(tparams)

	if len(tparams) > 0 {
		c.openScope(d, "type parameters of "+d.Name.Name)
		defer c.closeScope()
		c.declareTypeParams(d.TypeParams, tparams)
	}

	// A.6.5's DiscriminantType is a PredeclaredTypeName. Only an integer one
	// makes sense, since ⊢ "a unit-only enum is its discriminant integer."
	var discrim *types.Basic
	if d.Discrim != nil {
		dt := c.typ(d.Discrim)
		if b := types.AsBasic(dt); b != nil && types.IsInteger(dt) {
			discrim = b
		} else if !types.IsInvalid(dt) {
			c.errorExpr(d.Discrim, diag.NotAType, exprString(d.Discrim))
		}
	}

	seen := make(map[string]bool, len(d.Variants))
	variants := make([]*types.Variant, 0, len(d.Variants))
	next := int64(0)

	for _, v := range d.Variants {
		name := v.Name.Name
		if !v.Name.IsBlank() {
			if seen[name] {
				c.errorExpr(v.Name, diag.DuplicateVariant, name, d.Name.Name)
				continue
			}
			seen[name] = true
		}

		payload := make([]types.Type, 0, len(v.Payload))
		for _, p := range v.Payload {
			payload = append(payload, c.typ(p))
		}

		val := next
		if v.Value != nil {
			// A.6.5 ⊢ "= Expression (an explicit discriminant) is legal only on
			// a unit variant, and only when a DiscriminantType was declared."
			// A.14 lists the payload case among the forms that parse and are
			// rejected here.
			switch {
			case len(payload) > 0:
				c.errorExpr(v.Value, diag.PayloadDiscrim)
			case discrim == nil:
				c.errorExpr(v.Value, diag.DiscrimNoType)
			default:
				if n, ok := c.constIntExpr(v.Value); ok {
					val = n
				}
			}
		}
		// A.6.5 ⊢ "unassigned variants continue from the previous value."
		next = val + 1

		variants = append(variants, &types.Variant{
			Name:    name,
			Payload: payload,
			Value:   val,
		})
	}

	named.SetUnderlying(types.NewEnum(variants, discrim))
}

// aliasDecl resolves a TypeAliasDeclaration (A.6.6).
//
// The two targets behave oppositely, and that opposition is the whole content
// of this function:
//
//   - A.6.6 ⊢ "an alias to a Type is transparent: it names the same type and
//     satisfies a ~T type-set element." No Named is minted, the alias leaves no
//     trace in the type graph, and Identical sees straight through it.
//   - A.6.6 ⊢ "an alias to abstract is nominal and opaque... Two abstract
//     aliases never unify, however identical their provenance." So this one
//     mints an Abstract keyed on this object, and identity is the object.
//
// ✗ `type SDL_Window = ref` never reaches here: bare `ref` is not a type and
// the parser produces something typ() rejects as NotAType.
func (c *Checker) aliasDecl(obj *types.TypeName, d *ast.TypeAliasDecl) {
	if _, isAbstract := stripParens(d.Target).(*ast.AbstractType); isAbstract {
		// A.3.3 ⊢ an abstract handle "has a zeroed representation, legal only
		// as the value half of an error-path tuple paired with a non-empty
		// error string." Nothing about that is a shape, so it is not recorded
		// here — the boundary-tuple rule is A.8.4's, checked over a return.
		//
		// The family comes from the tag rather than a block, because a bare
		// alias is not inside one. A handle actually minted by a declare block
		// gets the block's family via foreignDecl.
		named := types.NewNamed(obj, nil)
		obj.SetType(named)
		named.SetUnderlying(types.NewAbstract(obj, c.familyForTag()))
		return
	}

	tparams := c.typeParams(d.TypeParams)
	if len(tparams) > 0 {
		c.openScope(d, "type parameters of "+d.Name.Name)
		defer c.closeScope()
		c.declareTypeParams(d.TypeParams, tparams)
	}
	obj.SetType(c.typ(d.Target))
}

// familyForTag maps a build tag to an import family for a bare abstract alias
// (A.4.4).
//
// A block-minted handle uses familyForBlock instead, which is strictly better
// informed — under darwin the tag alone cannot separate a flat C module from an
// Objective-C framework, so this answers Unknown there and the cast check
// treats Unknown conservatively.
func (c *Checker) familyForTag() types.Family {
	switch c.conf.Tag {
	case token.TagJS:
		return types.FamilyObjectGraph
	case token.TagDarwin:
		return types.FamilyUnknown
	}
	return types.FamilyMemoryFlat
}

// constraintDecl resolves a ConstraintDeclaration (A.7.2).
//
// A.7.2 ⊢ "multiple elements in a constraint body form an intersection: a type
// argument must satisfy all of them." So each element contributes to the same
// constraint rather than replacing what came before, and a bare ConstraintName
// element "embeds that constraint's set."
func (c *Checker) constraintDecl(obj *types.TypeName, d *ast.ConstraintDecl) {
	var terms []types.Term
	var embeds []*types.Constraint
	var methods []*types.Func

	seen := make(map[string]bool)

	for _, e := range d.Elems {
		switch {
		case e.Method != nil:
			// A.7.2 ⊢ a MethodRequirement "is satisfied by any type declaring a
			// matching receiver method. Because every instantiation is
			// monomorphized, the call in the generic body lowers to a direct
			// call on the concrete type. This is not an interface value and
			// introduces no vtable."
			name := e.Method.Name.Name
			if seen[name] {
				c.errorExpr(e.Method.Name, diag.DuplicateDeclaration, name)
				continue
			}
			seen[name] = true

			sig := c.methodReqSignature(e.Method)
			methods = append(methods, types.NewFunc(
				e.Method.Name.Pos(), c.pkg, name, sig))

		case e.Set != nil:
			sub := c.constraintExpr(e.Set)
			if sub == nil {
				continue
			}
			terms = append(terms, sub.Terms()...)
			embeds = append(embeds, sub.Embeds()...)
			methods = append(methods, sub.Methods()...)
		}
	}

	obj.SetConstraint(types.NewConstraint(obj, terms, methods, embeds))
}

// methodReqSignature builds the signature of a MethodRequirement (A.7.2).
//
// The requirement names no receiver — the receiver is exactly what varies
// across satisfying types — so Recv is nil and Constraint.Satisfies compares
// the rest.
func (c *Checker) methodReqSignature(m *ast.MethodReq) *types.Signature {
	params, variadic := c.paramList(m.Params)

	var results *types.Tuple
	if m.Result != nil {
		rt := c.typ(m.Result)
		if tup, ok := rt.(*types.Tuple); ok {
			results = tup
		} else {
			results = types.NewTuple(types.NewVar(m.Result.Pos(), c.pkg, "", rt))
		}
	}
	return types.NewSignature(nil, params, results, variadic, types.MarkerNone)
}

// funcDecl resolves a FunctionDeclaration, and with it A.6.4's initializer and
// deinitializer forms.
//
// Those need no separate path: A.1.3 makes `init` and `deinit`
// ContextualKeywords that are ordinary method names in a receiver declaration,
// so they arrive as IDENT and land in Name like any other method — the same
// reason ast.FuncDecl covers all three.
func (c *Checker) funcDecl(obj *types.Func, d *ast.FuncDecl) {
	// The receiver and the type parameters share one scope, because A.7.6 ⊢ the
	// receiver's [T] "binds the name" — the names it brings in must be visible
	// to the receiver's own type expression.
	if d.Recv != nil || d.TypeParams != nil {
		c.openScope(d, "signature of "+d.Name.Name)
		defer c.closeScope()
	}

	var recv *types.Var
	if d.Recv != nil {
		c.declareReceiverParams(d.Recv)
		recv = c.receiverVar(d.Recv)
	}

	if d.TypeParams != nil {
		tparams := c.typeParams(d.TypeParams)
		c.declareTypeParams(d.TypeParams, tparams)
	}

	sig := c.signature(recv, d.Type)
	obj.SetType(sig)

	c.checkDeclKind(obj, d, sig)
}

// checkDeclKind applies the rules that depend on what kind of function this
// turned out to be.
func (c *Checker) checkDeclKind(obj *types.Func, d *ast.FuncDecl, sig *types.Signature) {
	// A.12.1 ⊢ `test` "is a FunctionMarker, legal only under build test." It is
	// also, per A.2.2, "the only build tag that changes what is grammatical
	// rather than only what is linkable."
	if sig.Marker() == types.MarkerTest && !token.LicensesTest(c.conf.Tag) {
		c.errorExpr(d.Name, diag.MutOutsidePosition, "test")
	}

	// A.6.4 ⊢ "deinit takes no parameters and returns nothing."
	if obj.IsDeinit() {
		if sig.Params().Len() > 0 || !sig.Results().IsUnit() {
			c.errorExpr(d.Name, diag.InitializerResult, obj.Name())
		}
	}

	// A.6.1 ⊢ "a function named main must take no parameters, return nothing,
	// and acts as the program entry point." It also sets [+Await] in its body,
	// which is why the shape matters beyond convention.
	if obj.Name() == "main" && sig.Recv() == nil && !obj.IsEntry() {
		c.errorExpr(d.Name, diag.InitializerResult, "main")
	}
}

// declareReceiverParams brings A.7.6's receiver type parameters into scope.
//
// A.7.6 ⊢ "a method receiver re-declares the type's parameter list to bring the
// names into scope. The receiver's [T] binds the name; it does not introduce a
// fresh one." So each name is bound to the *receiver type's own* TypeParam
// rather than to a newly minted one — that is what makes a constraint declared
// on the type "in force inside every method of that type."
func (c *Checker) declareReceiverParams(r *ast.Receiver) {
	e := stripParens(r.Type)
	if own, ok := e.(*ast.OwnershipType); ok {
		e = stripParens(own.X)
	}
	ix, ok := e.(*ast.IndexExpr)
	if !ok {
		return
	}
	base, ok := stripParens(ix.X).(*ast.Ident)
	if !ok {
		return
	}
	obj := c.pkg.Scope().Lookup(base.Name)
	tn, ok := obj.(*types.TypeName)
	if !ok {
		return
	}
	c.objDecl(tn)

	named := types.AsNamed(tn.Type())
	if named == nil {
		return
	}
	tps := named.TypeParams()

	for i, arg := range ix.Indices {
		id, ok := stripParens(arg).(*ast.Ident)
		if !ok || i >= len(tps) {
			continue
		}
		c.declare(c.scope, id, types.NewTypeName(id.Pos(), c.pkg, id.Name, tps[i]))
	}
}

// receiverVar resolves `( Identifier : ReceiverType )` (A.6.1).
//
// A.6.1 ⊢ "only the receiver position may carry shared; it is the spelling for
// a method that needs a strong handle to itself in order to hand out weak
// back-references." And ⊢ "a receiver typed var consumes its receiver
// unconditionally: the receiver position has no argument slot to carry a var
// marker, so there is no bare form that copies. This is the single exception to
// the bare-means-copy rule."
//
// That exception is why the mode lives on the Var rather than being decided at
// each call site the way A.4.6 decides it for an argument.
func (c *Checker) receiverVar(r *ast.Receiver) *types.Var {
	mode, inner := splitMode(r.Type)
	t := c.typ(inner)

	v := types.NewParam(r.Name.Pos(), c.pkg, r.Name.Name, t, mode)
	// A `mut` receiver lowers to a pointer to the caller's slot and a `var` one
	// owns its copy; both give the body something addressable to work with.
	v.SetMutable(mode == types.ModeMut || mode == types.ModeVar)

	c.declare(c.scope, r.Name, v)
	return v
}

func (c *Checker) signature(recv *types.Var, ft *ast.FuncType) *types.Signature {
	params, variadic := c.paramList(ft.Params)

	var results *types.Tuple
	if ft.Result != nil {
		// A.3.4 ⊢ "omitting -> Type is the void form. There is no void type
		// name." A nil Result therefore yields Unit, which NewSignature
		// substitutes — a void function and a unit-returning one are the same
		// thing and nothing has to special-case one against the other.
		rt := c.typ(ft.Result)
		if tup, ok := rt.(*types.Tuple); ok {
			results = tup
		} else {
			results = types.NewTuple(types.NewVar(ft.Result.Pos(), c.pkg, "", rt))
		}
	}

	marker := types.MarkerNone
	if ft.Marker != nil {
		switch ft.Marker.Name {
		case "async":
			marker = types.MarkerAsync
		case "gpu":
			marker = types.MarkerGPU
		case "npu":
			marker = types.MarkerNPU
		case token.CtxTest:
			marker = types.MarkerTest
		}
	}
	return types.NewSignature(recv, params, results, variadic, marker)
}

// paramList resolves A.6.1's ParameterList.
func (c *Checker) paramList(l *ast.ParamList) (*types.Tuple, bool) {
	if l == nil {
		return types.NewTuple(), false
	}

	vars := make([]*types.Var, 0, len(l.List))
	variadic := false
	seen := make(map[string]bool, len(l.List))

	for i, p := range l.List {
		if p.Ellipsis.IsValid() {
			// A.6.1 ⊢ "a variadic parameter must be last, and there may be at
			// most one. It lowers to a stack-local fixed array plus a two-word
			// view over it." The parser accepts any arrangement so both
			// violations get their own rule rather than a syntax error.
			if variadic {
				c.errorAt(p.Ellipsis, p.Ellipsis, diag.MultipleVariadic)
			} else if i != len(l.List)-1 {
				c.errorAt(p.Ellipsis, p.Ellipsis, diag.VariadicNotLast)
			}
			variadic = true
		}

		mode, inner := splitMode(p.Type)
		t := c.typ(inner)

		name := ""
		if p.Name != nil {
			name = p.Name.Name
			if !p.Name.IsBlank() {
				if seen[name] {
					c.errorExpr(p.Name, diag.DuplicateDeclaration, name)
				}
				seen[name] = true
			}
			if types.Reserved(name) {
				c.errorExpr(p.Name, diag.ShadowedBuiltin, name)
			}
		}

		v := types.NewParam(p.Pos(), c.pkg, name, t, mode)
		// A.3.2 ⊢ `mut T` "lowers to a pointer to the caller's slot — which is
		// why its argument must be an addressable var binding or field path."
		// Inside the body the parameter itself is addressable either way.
		v.SetMutable(mode == types.ModeMut || mode == types.ModeVar)

		if p.Name != nil {
			c.info.RecordDef(p.Name, v)
		}
		vars = append(vars, v)
	}
	return types.NewTuple(vars...), variadic
}

// splitMode peels A.3.2's parameter-position qualifiers off a type expression.
//
// `mut` and `var` become a types.Mode rather than a Type — see types.Mode for
// why. Splitting here is what makes A.3.2's stacking rules fall out instead of
// needing a special case: `mut shared T` yields ModeMut over an *Ownership,
// while `mut var T` is unrepresentable because a Var carries one Mode.
//
// A stacked *type* qualifier (`shared unique T`) is a different error and is
// caught in ownershipType, since it survives as a nested node.
func splitMode(e ast.Expr) (types.Mode, ast.Expr) {
	own, ok := stripParens(e).(*ast.OwnershipType)
	if !ok {
		return types.ModeNone, e
	}
	switch own.Kw {
	case token.MUT:
		return types.ModeMut, own.X
	case token.VAR:
		return types.ModeVar, own.X
	}
	return types.ModeNone, e
}

// varDecl resolves a top-level VariableDeclaration's type (A.2, A.5.1).
//
// An annotated binding resolves here. An inferred one cannot: its type comes
// from the initializer, and A.2 ⊢ that initializer "must be compile-time-
// evaluable" — which needs the constant evaluator, and that is typecheck's. The
// object's type is left nil rather than marked invalid so typecheck can fill it
// without a cascade of errors from this pass.
func (c *Checker) varDecl(obj *types.Var, d *ast.VarDecl, node ast.Node) {
	b, _ := node.(*ast.Binding)
	if b == nil {
		obj.SetType(c.invalid())
		return
	}
	if b.Type != nil {
		obj.SetType(c.typ(b.Type))
		return
	}
	// A.5.1 ⊢ bare `var x: T` "requires a TypeAnnotation, since there is
	// nothing to infer from." No annotation and no initializer is therefore a
	// parse-level shape the parser already rejected; reaching here with both
	// nil means an initializer exists and typecheck owns it.
}

// foreignDecl resolves one member of a declare block (A.8.3).
//
// Everything here is shape-checking a linkage boundary, and A.8.3's governing
// sentence is: "exactly what is written is what is linked. A declare block
// contains only declarations corresponding to real entry points: no marker
// declarations, no visibility modifiers, no remapping clauses, and no bodies."
func (c *Checker) foreignDecl(obj types.Object, d *declInfo) {
	switch x := d.node.(type) {
	case *ast.ForeignFunc:
		c.foreignFunc(obj, x, d)

	case *ast.ForeignClass:
		// A.8.3 ✗ fields are banned. The parser keeps a *ForeignField in
		// Members so the diagnostic can point at the field rather than at a
		// stray colon; collectForeign reported the top-level ones, and this
		// covers the nested case.
		for _, m := range x.Members {
			if ff, ok := m.(*ast.ForeignField); ok {
				c.errorExpr(ff.Name, diag.ForeignField, ff.Name.Name)
			}
		}
		for _, mod := range x.Modifiers {
			c.errorExpr(mod, diag.ForeignModifier)
		}

		if tn, ok := obj.(*types.TypeName); ok {
			named := types.NewNamed(tn, nil)
			tn.SetType(named)
			// A foreign class is an opaque handle: A.8.3 ⊢ a declare block
			// "describes call shape only, never foreign-side layout", so there
			// is nothing to give it but an Abstract.
			named.SetUnderlying(types.NewAbstract(tn, d.family))
			c.foreignClassMembers(named, x, d)
		}
	}
}

func (c *Checker) foreignFunc(obj types.Object, x *ast.ForeignFunc, d *declInfo) {
	// A.8.3 ✗ `private init func() -> Bad` — visibility modifiers are banned.
	// They parse so the form can be diagnosed as itself rather than as a
	// syntax error.
	for _, mod := range x.Modifiers {
		c.errorExpr(mod, diag.ForeignModifier)
	}
	// A.8.3 ✗ `func SDL_Init() -> int32 { return 0 }` — declarations cannot
	// have bodies.
	if x.Body != nil {
		c.errorAt(x.Body.Pos(), x.Body.End(), diag.ForeignBody)
	}

	params, variadic := c.paramList(x.Params)

	// A.8.3 ⊢ "var — and any consume or transfer marking — is banned from a
	// foreign declaration. Ownership is a fact about a wrapper's field, decided
	// in the wrapper, not a decoration on an external stub."
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		if p.Mode() != types.ModeNone {
			c.errorAt(x.Params.List[i].Pos(), x.Params.List[i].End(),
				diag.MutOutsidePosition, p.Mode().String())
		}
	}

	var results *types.Tuple
	if x.Result != nil {
		rt := c.typ(x.Result)
		if tup, ok := rt.(*types.Tuple); ok {
			results = tup
		} else {
			results = types.NewTuple(types.NewVar(x.Result.Pos(), c.pkg, "", rt))
		}
	}

	if f, ok := obj.(*types.Func); ok {
		f.SetType(types.NewSignature(nil, params, results, variadic, types.MarkerNone))
	}
}

// foreignClassMembers resolves a foreign class's methods and initializers
// (A.8.3).
//
// A.8.3 ⊢ "init is a prefix modifier on func, not a function name. The unnamed
// form is what bare Type(...) construction resolves to; the named form is what
// Type.someName(...) resolves to." Both attach to the class rather than to
// package scope, which is why they are resolved here and not in collectForeign.
func (c *Checker) foreignClassMembers(named *types.Named, x *ast.ForeignClass, d *declInfo) {
	unnamedInit := false

	for _, m := range x.Members {
		ff, ok := m.(*ast.ForeignFunc)
		if !ok {
			continue
		}

		for _, mod := range ff.Modifiers {
			c.errorExpr(mod, diag.ForeignModifier)
		}
		if ff.Body != nil {
			c.errorAt(ff.Body.Pos(), ff.Body.End(), diag.ForeignBody)
		}

		params, variadic := c.paramList(ff.Params)
		var results *types.Tuple
		if ff.Result != nil {
			rt := c.typ(ff.Result)
			results = types.NewTuple(types.NewVar(ff.Result.Pos(), c.pkg, "", rt))
		}

		name := token.CtxInit
		pos := ff.Func
		if ff.Name != nil {
			name = ff.Name.Name
			pos = ff.Name.Pos()
		}

		if ff.Init.IsValid() {
			// A.8.3 ⊢ "at most one unnamed initializer per foreign class."
			if ff.Name == nil {
				if unnamedInit {
					c.errorAt(ff.Init, ff.Init, diag.DuplicateDeclaration, token.CtxInit)
					continue
				}
				unnamedInit = true
			}
			// A.8.3 ⊢ "an initializer must return its enclosing type."
			if results == nil || results.Len() != 1 ||
				!types.Identical(results.At(0).Type(), named) {
				c.errorAt(pos, pos, diag.InitializerResult, x.Name.Name)
			}
		}

		sig := types.NewSignature(
			types.NewParam(pos, c.pkg, "", named, types.ModeNone),
			params, results, variadic, types.MarkerNone)

		fn := types.NewFunc(pos, c.pkg, name, sig)
		if named.LookupMethod(name) != nil && ff.Name != nil {
			c.errorExpr(ff.Name, diag.DuplicateDeclaration, name)
			continue
		}
		named.AddMethod(fn)
		if ff.Name != nil {
			c.info.RecordDef(ff.Name, fn)
		}
	}
}

// ------------------------------------------------------------ type params

// typeParams builds A.7.1's list, performing the distribution the parser
// deliberately did not.
//
// A.7.1 ⊢ "a constraint written after a name applies to that name and to every
// immediately preceding unconstrained name in the same list — [A, B: Number]
// constrains both." ast.TypeParam leaves Constraint nil for the earlier names
// on purpose, because distributing in the tree "would erase the written form
// vfmt needs to reprint." This is where the distribution belongs.
//
// The walk runs backwards, carrying each written constraint left over the run
// of unconstrained names preceding it. A trailing run with no constraint at all
// gets `any`, per ⊢ "a bare name is constraint any: [T] means [T: any]."
func (c *Checker) typeParams(l *ast.TypeParamList) []*types.TypeParam {
	if l == nil {
		return nil
	}

	seen := make(map[string]bool, len(l.List))
	for _, p := range l.List {
		// A.7.1 ⊢ "names must be unique within a list." A BlankIdentifier is
		// exempt: A.1.2 makes `_` legal as a type-parameter name precisely
		// because it introduces no binding to collide.
		if p.Name.IsBlank() {
			continue
		}
		if seen[p.Name.Name] {
			c.errorExpr(p.Name, diag.DuplicateTypeParam, p.Name.Name)
		}
		seen[p.Name.Name] = true
	}

	out := make([]*types.TypeParam, len(l.List))
	var pending *types.Constraint

	for i := len(l.List) - 1; i >= 0; i-- {
		p := l.List[i]
		if p.Constraint != nil {
			pending = c.constraintExpr(p.Constraint)
		}
		cst := pending
		if cst == nil {
			cst = types.Any
		}
		out[i] = types.NewTypeParam(p.Name.Name, i, cst)

		// The written constraint reaches this name and every unconstrained name
		// to its left, and stops at the name that carried it.
		if p.Constraint != nil {
			pending = nil
		}
	}

	// A.7.1 ⊢ "a type parameter's scope begins after its own name and runs to
	// the end of the declaration's body, so a later parameter may be
	// constrained by an earlier one." That forward reference resolves because
	// declareTypeParams runs before any constraint expression is evaluated in a
	// body — the constraint here sees only what the enclosing scope already has.
	return out
}

func (c *Checker) declareTypeParams(l *ast.TypeParamList, tps []*types.TypeParam) {
	if l == nil {
		return
	}
	for i, p := range l.List {
		if i >= len(tps) {
			break
		}
		obj := types.NewTypeName(p.Name.Pos(), c.pkg, p.Name.Name, tps[i])
		c.declare(c.scope, p.Name, obj)
	}
}

// constIntExpr evaluates an integer constant expression in a declaration.
//
// It handles a literal, a named constant, and unary minus — the last because
// A.1.5.1 ⊢ "there is no literal syntax for a negative number. -1000 is unary
// minus applied to 1000, folded at compile time", so a negative enum
// discriminant is unspellable without it.
//
// Anything richer needs the full evaluator and belongs to typecheck. Rejecting
// it here rather than silently accepting means widening this later cannot
// change a program that currently compiles.
func (c *Checker) constIntExpr(e ast.Expr) (int64, bool) {
	switch x := stripParens(e).(type) {
	case *ast.UnaryExpr:
		if x.Op != token.SUB {
			break
		}
		n, ok := c.constIntExpr(x.X)
		return -n, ok

	case *ast.BasicLit:
		if x.Kind != token.INT {
			break
		}
		if v, ok := parseIntLit(x.Value); ok {
			return v, true
		}

	case *ast.Ident:
		obj := c.lookup(x)
		if obj == nil {
			return 0, false
		}
		cn, ok := obj.(*types.Const)
		if !ok {
			break
		}
		if n, ok := types.Int64Val(cn.Val()); ok {
			return n, true
		}
	}
	c.errorExpr(e, diag.ArrayLenNotConst)
	return 0, false
}

func unquote(s string) (string, bool) {
	if len(s) < 2 {
		return "", false
	}
	q := s[0]
	if (q != '"' && q != '`') || s[len(s)-1] != q {
		return "", false
	}
	return s[1 : len(s)-1], true
}