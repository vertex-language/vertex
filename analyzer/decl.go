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
// The nil type is the point. §1.1 ⊢ top-level declarations are order-
// independent: a declaration may reference any other in its package regardless
// of position or file. Inserting every name before resolving any type is what
// makes that true without a forward-declaration form.
func (c *Checker) collectObjects() {
	for _, f := range c.files {
		fileScope := types.NewScope(c.pkg.Scope(), "file "+f.Filename(c.conf.Fset))
		fileScope.SetExtent(f.Pos(), f.End())
		c.fileScopes[f] = fileScope
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
// §1.3 ⊢ the qualifier is the imported package's own package clause name, never
// the path — which is why this cannot be answered from the path alone, and why
// the importer must have read the imported directory's package line first. The
// same two-pass shape parser.ParseDir already uses.
//
// The qualifier is file-scoped while the declarations that use it are
// package-scoped, so one file's import never resolves a name in another's.
func (c *Checker) collectImports(f *ast.File, fileScope *types.Scope) {
	for _, imp := range f.Imports {
		for _, p := range imp.Paths {
			path, ok := decodeString(p.Value)
			if !ok {
				continue
			}
			if c.conf.Importer == nil {
				// A single-package check has no way to resolve a path. The name
				// is genuinely undeclared from this checker's position.
				c.errorExpr(p, diag.UndeclaredName, path)
				continue
			}
			pkg, err := c.conf.Importer.Import(path)
			if err != nil || pkg == nil {
				c.errorExpr(p, diag.UndeclaredName, path)
				continue
			}
			// There is no aliasing form, no dot-import, and no blank import, so
			// there is exactly one name to bind and no choice about what it is.
			// §1.3 makes two imports declaring the same name an error, which is
			// just a duplicate insert.
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
		// exists: a receiver may name a type declared later or in another file.
		if x.Recv != nil {
			return
		}
		obj := types.NewFunc(x.Name.Pos(), c.pkg, x.Name.Name, nil)
		c.funcObj[x] = obj
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
		// A constraint's TypeName has a non-nil Constraint from the start, so
		// IsConstraint answers correctly even before phase 2 fills in the set.
		// That is what rejects `var c: Ordered` no matter which order the two
		// declarations are checked in.
		obj := types.NewConstraintName(x.Name.Pos(), c.pkg, x.Name.Name, types.Any)
		c.declare(c.pkg.Scope(), x.Name, obj)
		c.objMap[obj] = &declInfo{decl: d, file: f, fileScope: fileScope}

	case *ast.VarDecl:
		// A top-level `let` is a constant: §6.1 requires its initializer to be
		// compile-time-evaluable, and the object that carries a folded value is
		// types.Const. A top-level `var` is an ordinary mutable binding under
		// the same initializer rule.
		for i, b := range x.Bindings {
			var obj types.Object
			if x.Kw == token.LET {
				obj = types.NewConst(b.Name.Pos(), c.pkg, b.Name.Name, nil, nil)
			} else {
				v := types.NewVar(b.Name.Pos(), c.pkg, b.Name.Name, nil)
				v.SetMutable(true)
				obj = v
			}
			c.declare(c.pkg.Scope(), b.Name, obj)
			c.objMap[obj] = &declInfo{
				decl: d, file: f, fileScope: fileScope, node: b, index: i,
			}
		}

	case *ast.DeclareDecl:
		c.collectForeign(x, f, fileScope)
	}
}

// collectForeign collects a declare block's members and checks the block's own
// legality.
//
// A declare block is a linkage boundary rather than a namespace: its symbols
// join the file's package. So members land in package scope and a collision
// with an ordinary declaration is an ordinary redeclaration.
func (c *Checker) collectForeign(d *ast.DeclareDecl, f *ast.File, fileScope *types.Scope) {
	if d.Kind == token.CtxFramework && d.Variant != nil {
		// A framework block takes no variant tag. The tagged form parses so
		// this can name the rule rather than failing at the bracket.
		c.errorAt(d.Variant.Pos(), d.Variant.End(), diag.FrameworkTag)
	}
	// The variant tag set is closed, but neither token nor this package holds
	// its membership — the set has no home yet, so UnknownVariantTag has no
	// site here. It belongs wherever the set eventually lives, next to
	// LookupBuildTag.

	family := c.familyForBlock(d)
	c.collectForeignMembers(d.Members, d, f, fileScope, family)
}

func (c *Checker) collectForeignMembers(
	members []ast.ForeignMember, d *ast.DeclareDecl,
	f *ast.File, fileScope *types.Scope, family types.Family,
) {
	for _, m := range members {
		switch x := m.(type) {
		case *ast.ForeignFunc:
			// `init` is a prefix modifier on func, not a function name; the
			// unnamed form is what bare `Type(...)` construction resolves to,
			// so it has no package-scope name and is reached only through its
			// enclosing class.
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

		case *ast.DeclareDecl:
			// A declare block may not contain another. It parses so the
			// diagnostic can name the construct.
			c.errorAt(x.Declare, x.KindPos, diag.NestedDeclare)

		case *ast.Field:
			// A declare block describes call shape only, never foreign-side
			// layout. Fields parse so the caret lands on the field rather than
			// on a stray colon.
			c.errorExpr(x.Name, diag.ForeignField, x.Name.Name)
		}
	}
}

// familyForBlock decides the import family a declare block's handles belong to.
//
// The family is what makes `abstract` → `typed_ptr T` legal or not: it is legal
// only for a memory-flat family, and a compile error for an object-graph family
// whose handles have no byte representation to point at.
//
// The build tag alone cannot answer it on every target. Under darwin a
// `declare module` is flat C while a `declare framework` is an object graph, so
// the block keyword refines what the tag started.
func (c *Checker) familyForBlock(d *ast.DeclareDecl) types.Family {
	if d.Kind == token.CtxFramework {
		return types.FamilyObjectGraph
	}
	if c.conf.Tag == token.TagJS {
		return types.FamilyObjectGraph
	}
	return types.FamilyMemoryFlat
}

// familyForTag maps a build tag to a family for a bare abstract alias, which is
// not inside a block. It is strictly less informed than familyForBlock — under
// darwin the tag alone cannot separate flat C from an object graph — so it
// answers Unknown there and the cast check treats Unknown conservatively.
func (c *Checker) familyForTag() types.Family {
	switch c.conf.Tag {
	case token.TagJS:
		return types.FamilyObjectGraph
	case token.TagDarwin:
		return types.FamilyUnknown
	}
	return types.FamilyMemoryFlat
}

// collectMethods attaches each method to its receiver's Named type.
//
// It runs after every type name exists, for the same reason phase 1 exists at
// all: a receiver may name a type declared later in the file or in another file
// of the package.
func (c *Checker) collectMethods() {
	for _, f := range c.files {
		fileScope := c.fileScopeOf(f)
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Recv == nil {
				continue
			}

			// §2.3 ⊢ a reserved builtin name may not be declared as a member,
			// method, field, local, or parameter.
			if types.Reserved(fd.Name.Name) {
				c.errorExpr(fd.Name, diag.ReservedAsMember, fd.Name.Name)
				continue
			}

			// A method may not declare its own type parameters; everything it
			// is generic over comes from its receiver. The parser parses the
			// list either way so the caret lands on the brackets.
			if fd.TypeParams != nil {
				c.errorAt(fd.TypeParams.Pos(), fd.TypeParams.End(),
					diag.MethodTypeParams, fd.Name.Name)
			}

			obj := types.NewFunc(fd.Name.Pos(), c.pkg, fd.Name.Name, nil)
			c.funcObj[fd] = obj
			c.objMap[obj] = &declInfo{decl: fd, file: f, fileScope: fileScope}
			c.recordDef(fd.Name, obj)

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
// It deliberately does not resolve the receiver's full type. A receiver's
// bracket list re-declares the type's existing parameter names rather than
// introducing fresh ones, so the base name is what associates the method and
// the list is handled when the signature is built.
func (c *Checker) receiverBase(r *ast.Receiver, fileScope *types.Scope) *types.Named {
	e := stripParens(r.Type)
	if own, ok := e.(*ast.OwnershipType); ok {
		e = stripParens(own.X)
	}
	if ix, ok := e.(*ast.IndexExpr); ok {
		e = stripParens(ix.X)
	}

	id, ok := e.(*ast.Ident)
	if !ok {
		c.errorExpr(r.Type, diag.NotAType, exprString(r.Type))
		return nil
	}

	// A ReceiverType is a TypeName, not a qualified one: a method on another
	// package's type is not a thing, so this looks only in package scope.
	obj := c.pkg.Scope().Lookup(id.Name)
	if obj == nil {
		c.errorExpr(id, diag.UndeclaredName, id.Name)
		return nil
	}
	tn, ok := obj.(*types.TypeName)
	if !ok || tn.IsConstraint() {
		c.errorExpr(r.Type, diag.NotAType, id.Name)
		return nil
	}
	c.recordUse(id, tn)

	saved := c.scope
	c.scope = fileScope
	c.objDecl(tn)
	c.scope = saved

	return types.AsNamed(tn.Type())
}

func (c *Checker) fileScopeOf(f *ast.File) *types.Scope {
	if s, ok := c.fileScopes[f]; ok {
		return s
	}
	return c.pkg.Scope()
}

// ------------------------------------------------- phase 2: resolve types

func (c *Checker) resolveDeclTypes() {
	// Iteration order over objMap is unspecified, which is fine: objDecl is
	// idempotent and resolves dependencies on demand, so the result does not
	// depend on it. Diagnostics are sorted by position by the caller.
	for obj := range c.objMap {
		c.objDecl(obj)
	}
}

// objDecl resolves one object's type, on demand and at most once.
//
// The resolving stack is what turns a self-referential declaration into a
// diagnostic instead of a stack overflow. Order-independence makes a cycle
// reachable from any entry point, so the guard lives here rather than in a
// pre-pass over a dependency graph.
func (c *Checker) objDecl(obj types.Object) {
	if obj == nil || obj.Type() != nil {
		return
	}
	if tn, ok := obj.(*types.TypeName); ok && tn.IsConstraint() && tn.Constraint() != types.Any {
		return // already resolved; a constraint name carries no Type
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
	// resolves through the file that wrote it and no other.
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
		c.varDecl(obj, x, d)
	case *ast.DeclareDecl:
		c.foreignDecl(obj, d)
	}
}

// recordDecl resolves a StructDecl or a ClassDecl.
//
// One path for both, matching ast.RecordDecl: a class is byte-for-byte
// identical in layout to a struct and differs only in its member and method
// model. The keyword is carried into Struct.class, and the two consumers that
// care read it there — construction syntax (a struct is built by a composite
// literal, a class by calling an initializer) and `===`, which is legal on
// classes only.
//
// The Named is bound to its TypeName before any field resolves. That two-step
// is what lets `struct Node { next: typed_ptr Node }` reach its own name
// without tripping the cycle guard: the guard fires on an object whose type is
// still nil, and this one's is not.
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
			// A field default is evaluated at each construction for each
			// omitted field, so evaluating it is a construction-site question;
			// recording that one exists is all the type needs.
			HasDefault: f.Default != nil,
		})
		c.recordDef(f.Name, types.NewField(f.Name.Pos(), c.pkg, name, ft))
	}

	// Fields are laid out in declaration order, so the slice order is the
	// layout order and nothing reorders it.
	named.SetUnderlying(types.NewStruct(fields, d.Kw == token.CLASS))
}

// enumDecl resolves an enum declaration.
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

	// Only an integer DiscriminantType makes sense: a unit-only enum is its
	// discriminant integer.
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
			// An explicit discriminant is legal only on a unit variant, and
			// only with a declared discriminant type. Both suffixes parse on
			// any variant so each rejection can name itself.
			switch {
			case len(payload) > 0:
				c.errorExpr(v.Value, diag.PayloadDiscrim)
			case discrim == nil:
				c.errorExpr(v.Value, diag.DiscrimNoType)
			default:
				cv := c.constValue(v.Value)
				n, ok := types.Int64Val(cv)
				switch {
				case !ok:
					c.errorExpr(v.Value, diag.ArrayLenNotConst)
				case !c.sizes.Representable(cv, discrim):
					c.errorExpr(v.Value, diag.NotRepresentable,
						cv.String(), types.TypeString(discrim))
				default:
					val = n
				}
			}
		}
		// An unassigned variant continues from the previous value. The sources
		// do not say what an omitted discriminant is; this is the only reading
		// under which a partially annotated list has an answer at all, and it
		// is this implementation's choice.
		next = val + 1

		variants = append(variants, &types.Variant{
			Name: name, Payload: payload, Value: val,
		})
	}

	named.SetUnderlying(types.NewEnum(variants, discrim))
}

// aliasDecl resolves a type alias. The two targets behave oppositely, and that
// opposition is the whole content of this function:
//
//   - An alias to a Type is transparent: it names the same type at every depth
//     of composition. No Named is minted, so it leaves no trace in the type
//     graph and Identical sees straight through it.
//   - An alias to `abstract` is nominal and opaque, and each such alias is
//     distinct from every other. So this one mints an Abstract keyed on this
//     object, and identity is the object.
func (c *Checker) aliasDecl(obj *types.TypeName, d *ast.TypeAliasDecl) {
	if _, isAbstract := stripParens(d.Target).(*ast.AbstractType); isAbstract {
		named := types.NewNamed(obj, nil)
		obj.SetType(named)
		// The family comes from the tag, because a bare alias is not inside a
		// declare block. A handle actually minted by one gets the block's
		// family through foreignDecl.
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

// constraintDecl resolves a constraint declaration.
//
// §9 ⊢ multiple elements in the body are an intersection, so each element
// contributes to the same constraint rather than replacing what came before,
// and a bare constraint name element embeds that constraint's set.
func (c *Checker) constraintDecl(obj *types.TypeName, d *ast.ConstraintDecl) {
	var (
		terms   []types.Term
		embeds  []*types.Constraint
		methods []*types.Func
	)
	seen := make(map[string]bool)

	for _, e := range d.Elems {
		switch {
		case e.Method != nil:
			// A MethodRequirement is satisfied by any type declaring a matching
			// receiver method. Because every instantiation is monomorphized,
			// the call in the generic body lowers to a direct call on the
			// concrete type — this is not an interface and introduces no vtable.
			name := e.Method.Name.Name
			if seen[name] {
				c.errorExpr(e.Method.Name, diag.DuplicateDeclaration, name)
				continue
			}
			seen[name] = true

			// A MethodRequirement takes a full Signature, so a constraint can
			// require a marked method; Recv stays nil, since the receiver is
			// exactly what varies across satisfying types.
			sig := c.signature(nil, e.Method.Type, false)
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

// funcDecl resolves a function declaration, and with it the initializer and
// deinitializer forms.
//
// Those need no separate path: `init` and `deinit` are contextual keywords that
// are ordinary method names in a receiver declaration, so they arrive as
// identifiers and land in Name like any other — the same reason ast.FuncDecl
// covers all three.
func (c *Checker) funcDecl(obj *types.Func, d *ast.FuncDecl) {
	// The receiver and the type parameters share one scope, because a
	// receiver's bracket list binds names its own type expression must see.
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
		c.declareTypeParams(d.TypeParams, c.typeParams(d.TypeParams))
	}

	sig := c.signature(recv, d.Type, true)
	obj.SetType(sig)

	// An Expected result reaches the grammar only through a declaration, and is
	// restricted further to a file carrying `build test`. That restriction is
	// the file's tag, and this is the one place it can be read.
	if sig.Expected() != nil && !token.LicensesTest(c.conf.Tag) {
		c.errorExpr(d.Name, diag.ExpectedOutsideTest)
	}
}

// declareReceiverParams brings a receiver's type parameters into scope.
//
// A receiver re-declares the type's parameter list to bring the names in; the
// bracket list binds the existing names rather than introducing fresh ones. So
// each name is bound to the receiver type's own TypeParam, which is what keeps
// a constraint declared on the type in force inside every method of it.
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
	tn, ok := c.pkg.Scope().Lookup(base.Name).(*types.TypeName)
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

// receiverVar resolves `( identifier : ReceiverType )`.
//
// A receiver typed `var` consumes its receiver unconditionally: the receiver
// position has no argument slot to carry a transfer marker, so there is no bare
// form that copies. That is the single exception to bare-means-copy, and it is
// why the mode lives on the Var here rather than being decided per call site.
func (c *Checker) receiverVar(r *ast.Receiver) *types.Var {
	mode, inner := splitMode(r.Type)
	t := c.typ(inner)

	v := types.NewParam(r.Name.Pos(), c.pkg, r.Name.Name, t, mode)
	// A `mut` receiver lowers to a pointer to the caller's slot and a `var` one
	// owns its copy; both give the body something addressable.
	v.SetMutable(mode == types.ModeMut || mode == types.ModeVar)

	c.declare(c.scope, r.Name, v)
	return v
}

// signature resolves a Signature. declResult admits the Expected result form,
// which reaches the grammar only through a function or method declaration —
// that is what keeps it out of a function type and a function literal.
func (c *Checker) signature(recv *types.Var, ft *ast.FuncType, declResult bool) *types.Signature {
	if ft == nil {
		return types.NewSignature(recv, nil, nil, false, types.MarkerNone)
	}

	savedSig := c.inSignature
	c.inSignature = true
	defer func() { c.inSignature = savedSig }()

	params, variadic := c.paramList(ft.Params)

	// A signature carries at most one marker, but the repetition is written so
	// more than one parses; the parser keeps them all and the extras are
	// rejected here.
	marker := types.MarkerNone
	if n := len(ft.Markers); n > 0 {
		marker = markerFor(ft.Markers[0])
		if n > 1 {
			c.errorAt(ft.Markers[1].Pos(), ft.Markers[n-1].End(), diag.MultipleMarkers)
		}
	}

	// A tensor type is legal in an npu-marked function's own signature as well
	// as in its body, so the result and parameter types are resolved with that
	// context already set.
	savedCtx := c.ctx
	c.ctx.npu = marker == types.MarkerNPU
	defer func() { c.ctx = savedCtx }()

	var (
		results  *types.Tuple
		expected *types.Expected
	)
	if ft.Result != nil {
		if call, ok := ft.Result.(*ast.CallExpr); ok && declResult && isExpectedCall(call) {
			expected = c.expectedResult(call)
		} else {
			rt := c.typ(ft.Result)
			if tup, ok := rt.(*types.Tuple); ok {
				results = tup
			} else {
				results = types.NewTuple(types.NewVar(ft.Result.Pos(), c.pkg, "", rt))
			}
		}
	}
	// Omitting the result is the void form; there is no `void` type name, so a
	// nil results tuple is exactly what NewSignature turns into the empty one.

	sig := types.NewSignature(recv, params, results, variadic, marker)
	if expected != nil {
		sig.SetExpected(expected)
	}
	return sig
}

func markerFor(m *ast.Marker) types.Marker {
	switch m.Kind {
	case token.ASYNC:
		return types.MarkerAsync
	case token.GPU:
		return types.MarkerGPU
	case token.NPU:
		return types.MarkerNPU
	}
	// `test` is a contextual keyword and scans as an identifier, which is why
	// the node records a name alongside a kind.
	if m.Name == token.CtxTest {
		return types.MarkerTest
	}
	return types.MarkerNone
}

func isExpectedCall(x *ast.CallExpr) bool {
	id, ok := x.Fun.(*ast.Ident)
	return ok && id.Name == token.CtxExpected
}

// expectedResult resolves `Expected(TypeName, string_lit)` and
// `Expected(error [, string_lit])`.
//
// `Expected` and `error` are ordinary identifiers recognized only in this
// production, which is why the parser leaves a call node behind and `error`
// mints no TypeName. The message is normative text: it is compared against a
// Diagnostic's Text(), which is what diag's template registry exists to keep
// stable.
func (c *Checker) expectedResult(x *ast.CallExpr) *types.Expected {
	if len(x.Args) == 0 {
		c.errorExpr(x, diag.ExpectedToken, "a result type or 'error'", "no argument")
		return nil
	}

	msg, hasMsg := "", false
	if len(x.Args) > 1 {
		lit, ok := x.Args[1].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			c.errorExpr(x.Args[1], diag.ExpectedToken, "a message string",
				exprString(x.Args[1]))
		} else if s, ok := decodeString(lit.Value); ok {
			msg, hasMsg = s, true
		}
	}
	if len(x.Args) > 2 {
		c.errorExpr(x.Args[2], diag.ExpectedToken, "')'", exprString(x.Args[2]))
	}

	if id, ok := x.Args[0].(*ast.Ident); ok && id.Name == token.CtxError {
		return types.NewExpectedError(msg, hasMsg)
	}

	t := c.typ(x.Args[0])
	if !hasMsg {
		// There is no message-free value form; only the error form permits an
		// omitted message.
		c.errorExpr(x, diag.ExpectedToken, "a message string", "no argument")
	}
	return types.NewExpectedValue(t, msg)
}

// paramList resolves Parameters.
func (c *Checker) paramList(l *ast.ParamList) (*types.Tuple, bool) {
	if l == nil {
		return types.NewTuple(), false
	}

	vars := make([]*types.Var, 0, len(l.List))
	variadic := false
	named, bare := 0, 0
	seen := make(map[string]bool, len(l.List))

	for i, p := range l.List {
		if p.Ellipsis.IsValid() {
			// A variadic parameter must be last and there may be at most one.
			// The parser accepts any arrangement so each violation gets its own
			// rule rather than a syntax error at the ellipsis.
			if variadic {
				c.errorAt(p.Ellipsis, p.Ellipsis+3, diag.MultipleVariadic)
			} else if i != len(l.List)-1 {
				c.errorAt(p.Ellipsis, p.Ellipsis+3, diag.VariadicNotLast)
			}
			variadic = true
		}

		mode, inner := splitMode(p.Type)
		t := c.typ(inner)

		name := ""
		if p.Name != nil {
			named++
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
		} else {
			bare++
		}

		v := types.NewParam(p.Pos(), c.pkg, name, t, mode)
		// A `mut T` parameter lowers to a pointer to the caller's slot, which
		// is why its argument must be an addressable var binding or field path.
		// Inside the body the parameter itself is addressable either way.
		v.SetMutable(mode == types.ModeMut || mode == types.ModeVar)

		if p.Name != nil {
			c.recordDef(p.Name, v)
		}
		vars = append(vars, v)
	}

	// Names must be either all present or all absent within one list. A mixed
	// list parses so the diagnostic can point at the list rather than at
	// whichever colon happened to be missing.
	if named > 0 && bare > 0 {
		c.errorAt(l.Pos(), l.End(), diag.MixedParamNames)
	}
	return types.NewTuple(vars...), variadic
}

// splitMode peels the parameter-position qualifiers off a type expression.
//
// `mut` and `var` become a types.Mode rather than a Type — see types.Mode for
// why. Splitting here is what makes the stacking rules fall out instead of
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

// varDecl resolves a top-level declaration's type and folds its initializer.
//
// A top-level initializer must be compile-time-evaluable — there is no static
// initialization order and no initialization-time code — and the bare `var`
// form is rejected here. Both are static rules over this node, and both are
// this function.
func (c *Checker) varDecl(obj types.Object, d *ast.VarDecl, di *declInfo) {
	b, _ := di.node.(*ast.Binding)
	if b == nil {
		obj.SetType(c.invalid())
		return
	}

	if !d.Assign.IsValid() {
		// The initializer-free form covers `var w` and nothing else, and it has
		// no meaning at the top level.
		c.errorExpr(b.Name, diag.TopLevelBareVar)
		obj.SetType(c.invalid())
		return
	}

	var declared types.Type
	if b.Type != nil {
		declared = c.typ(b.Type)
		obj.SetType(declared)
	}

	// A binding list and a value list line up positionally. Anything else — one
	// call unbuilding a tuple across several bindings — needs the expression
	// typer, so the object's type is left for it rather than guessed at.
	if di.index >= len(d.Values) {
		if declared == nil {
			obj.SetType(c.invalid())
		}
		return
	}
	val := d.Values[di.index]

	// The transfer marker is not a constant, and a top-level initializer has
	// nothing to move from.
	if tr, ok := val.(*ast.TransferExpr); ok {
		c.errorExpr(tr, diag.TransferOutsideOwning)
		if declared == nil {
			obj.SetType(c.invalid())
		}
		return
	}

	v := c.constValue(val)
	if isUnknown(v) {
		c.errorExpr(val, diag.TopLevelVarNotConst)
		if declared == nil {
			obj.SetType(c.invalid())
		}
		return
	}

	if cn, ok := obj.(*types.Const); ok {
		cn.SetVal(v)
	}

	switch {
	case declared == nil:
		// Nothing imposes a type, so the untyped constant takes its default.
		obj.SetType(types.Default(untypedFor(v)))
	case types.AsBasic(declared) != nil && !types.IsInvalid(declared):
		// §4.1 ⊢ a literal whose value does not fit its destination is a
		// compile error, not a wraparound. Whether it fits is a question about
		// the target, which is why Representable is a method on Sizes.
		if !c.sizes.Representable(v, types.AsBasic(declared)) {
			c.errorExpr(val, diag.NotRepresentable, v.String(), types.TypeString(declared))
		}
	}
}

// untypedFor maps a folded constant back to the untyped kind it was carried as.
func untypedFor(v types.Value) types.Type {
	switch v.Kind() {
	case types.BoolKind:
		return types.Typ[types.UntypedBool]
	case types.IntKind:
		return types.Typ[types.UntypedInt]
	case types.FloatKind:
		return types.Typ[types.UntypedFloat]
	case types.CharKind:
		return types.Typ[types.UntypedChar]
	case types.StringKind:
		return types.Typ[types.UntypedString]
	}
	return types.Typ[types.Invalid]
}

// --------------------------------------------------------- declare blocks

// foreignDecl resolves one member of a declare block.
//
// Everything here is shape-checking a linkage boundary: exactly what is written
// is what is linked, so a declare block holds only declarations corresponding
// to real entry points — no bodies, no markers, no ownership decorations, and
// no foreign-side layout.
func (c *Checker) foreignDecl(obj types.Object, d *declInfo) {
	switch x := d.node.(type) {
	case *ast.ForeignFunc:
		c.foreignFunc(obj, x)

	case *ast.ForeignClass:
		tn, ok := obj.(*types.TypeName)
		if !ok {
			return
		}
		named := types.NewNamed(tn, nil)
		tn.SetType(named)
		// A foreign class is an opaque handle: the block describes call shape
		// only, so there is nothing to give it but an Abstract carrying the
		// block's family.
		named.SetUnderlying(types.NewAbstract(tn, d.family))
		c.foreignClassMembers(named, x)
	}
}

// foreignSignature resolves a foreign declaration's signature and raises the
// rejections that are about the declaration rather than the type.
func (c *Checker) foreignSignature(x *ast.ForeignFunc) *types.Signature {
	// A declare block describes call shapes only, so a body must parse in order
	// to be diagnosed as itself.
	if x.Body != nil {
		c.errorAt(x.Body.Pos(), x.Body.End(), diag.ForeignBody)
	}
	// A marker needs no field on the node: the signature already carries every
	// marker written.
	if x.Type != nil && len(x.Type.Markers) > 0 {
		m := x.Type.Markers
		c.errorAt(m[0].Pos(), m[len(m)-1].End(), diag.ForeignMarker)
	}

	sig := c.signature(nil, x.Type, false)

	// Ownership is a fact about a wrapper's field, decided in the wrapper — not
	// a decoration on an external stub. So `mut` and `var` are banned here.
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		if p.Mode() == types.ModeNone {
			continue
		}
		at := ast.Node(x)
		if x.Type != nil && x.Type.Params != nil && i < len(x.Type.Params.List) {
			at = x.Type.Params.List[i]
		}
		c.errorExpr(at, diag.MutOutsidePosition, p.Mode().String())
	}
	return sig
}

func (c *Checker) foreignFunc(obj types.Object, x *ast.ForeignFunc) {
	sig := c.foreignSignature(x)
	if f, ok := obj.(*types.Func); ok {
		f.SetType(sig)
	}
}

// foreignClassMembers resolves a foreign class's methods and initializers.
//
// `init` is a prefix modifier on func rather than a function name: the unnamed
// form is what bare `Type(...)` construction resolves to, the named form what
// `Type.name(...)` resolves to. Both attach to the class rather than to package
// scope, which is why they are resolved here and not in collectForeign.
func (c *Checker) foreignClassMembers(named *types.Named, x *ast.ForeignClass) {
	unnamedInit := false

	for _, m := range x.Members {
		switch mem := m.(type) {
		case *ast.Field:
			c.errorExpr(mem.Name, diag.ForeignField, mem.Name.Name)
			continue

		case *ast.DeclareDecl:
			c.errorAt(mem.Declare, mem.KindPos, diag.NestedDeclare)
			continue

		case *ast.ForeignFunc:
			sig := c.foreignSignature(mem)

			name := token.CtxInit
			pos := mem.Func
			if mem.Name != nil {
				name = mem.Name.Name
				pos = mem.Name.Pos()
			}

			if mem.Init.IsValid() {
				if mem.Name == nil {
					if unnamedInit {
						c.errorAt(mem.Init, mem.Func, diag.DuplicateDeclaration,
							token.CtxInit)
						continue
					}
					unnamedInit = true
				}
				// A foreign initializer returns its enclosing type.
				res := sig.Results()
				if res.Len() != 1 || !types.Identical(res.At(0).Type(), named) {
					c.errorAt(pos, pos, diag.ForeignInitResult, x.Name.Name)
				}
			}

			// The receiver is the class itself; a foreign method is reached
			// only through a handle to it.
			recv := types.NewParam(pos, c.pkg, "", named, types.ModeNone)
			full := types.NewSignature(recv, sig.Params(), sig.Results(),
				sig.Variadic(), types.MarkerNone)

			if named.LookupMethod(name) != nil {
				c.errorAt(pos, pos, diag.DuplicateDeclaration, name)
				continue
			}
			fn := types.NewFunc(pos, c.pkg, name, full)
			named.AddMethod(fn)
			if mem.Name != nil {
				c.recordDef(mem.Name, fn)
			}
		}
	}
}

// ------------------------------------------------------------ type params

// typeParams builds TypeParameters, performing the distribution the parser
// deliberately did not.
//
// A constraint written after a name applies to that name and to every
// immediately preceding unconstrained name in the same list — `[A, B: Number]`
// constrains both. ast.TypeParam leaves Constraint nil for the earlier names on
// purpose, because distributing in the tree would erase the written form a
// formatter needs to reproduce. This is where the distribution belongs.
//
// The walk runs backwards, carrying each written constraint left over the run
// of unconstrained names preceding it. A trailing run with no constraint at all
// gets `any`, since a bare name is constrained by `any`.
func (c *Checker) typeParams(l *ast.TypeParamList) []*types.TypeParam {
	if l == nil {
		return nil
	}

	seen := make(map[string]bool, len(l.List))
	for _, p := range l.List {
		// Names must be unique within a list. The blank identifier is exempt:
		// it introduces no binding to collide with.
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
		c.declare(c.scope, p.Name,
			types.NewTypeName(p.Name.Pos(), c.pkg, p.Name.Name, tps[i]))
	}
}