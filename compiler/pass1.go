package compiler

import (
	"fmt"

	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/wasm-compiler/wasm"
)

// collectDeclarations is Pass 1: walk the top-level and register all symbols
// before code generation begins. Order: types → imports → local functions.
func (c *Compiler) collectDeclarations(tree parser.ITopLevelContext) error {
	stmts := tree.Statements()
	if stmts == nil {
		return nil
	}
	all := stmts.AllStatement()

	// ── A: collect struct / class types ──────────────────────────────────────
	for _, stmt := range all {
		ds, ok := stmt.(*parser.DeclarationStatementContext)
		if !ok {
			continue
		}
		decl := ds.Declaration()
		if decl.StructDeclaration() != nil {
			if err := c.collectStructDecl(
				decl.StructDeclaration().(*parser.StructDeclarationContext),
			); err != nil {
				return err
			}
		}
		if decl.ClassDeclaration() != nil {
			if err := c.collectClassDecl(
				decl.ClassDeclaration().(*parser.ClassDeclarationContext),
			); err != nil {
				return err
			}
		}
	}

	// ── B: collect @extern function imports ───────────────────────────────────
	for _, stmt := range all {
		ds, ok := stmt.(*parser.DeclarationStatementContext)
		if !ok {
			continue
		}
		if ds.Declaration().FunctionDeclaration() == nil {
			continue
		}
		fd := ds.Declaration().FunctionDeclaration().(*parser.FunctionDeclarationContext)
		if isExternFunc(fd) {
			if err := c.collectExternFunc(fd); err != nil {
				return err
			}
		}
	}

	// ── C: collect local function declarations ────────────────────────────────
	for _, stmt := range all {
		ds, ok := stmt.(*parser.DeclarationStatementContext)
		if !ok {
			continue
		}
		if ds.Declaration().FunctionDeclaration() == nil {
			continue
		}
		fd := ds.Declaration().FunctionDeclaration().(*parser.FunctionDeclarationContext)
		if !isExternFunc(fd) {
			if err := c.collectLocalFunc(fd); err != nil {
				return err
			}
		}
	}

	return nil
}

// ── @extern detection ─────────────────────────────────────────────────────────

func isExternFunc(fd *parser.FunctionDeclarationContext) bool {
	// Must have no body.
	if fd.FunctionBody() != nil {
		return false
	}
	head := fd.FunctionHead()
	if head.Attributes() == nil {
		return false
	}
	for _, attr := range head.Attributes().AllAttribute() {
		if attr.Identifier().GetText() == "extern" {
			return true
		}
	}
	return false
}

// externNames extracts (module, name) from @extern("module", "symbol").
func externNames(fd *parser.FunctionDeclarationContext) (module, sym string) {
	module, sym = "env", fd.Identifier().GetText()
	for _, attr := range fd.FunctionHead().Attributes().AllAttribute() {
		if attr.Identifier().GetText() != "extern" {
			continue
		}
		if attr.AttributeArguments() == nil {
			break
		}
		al := attr.AttributeArguments().AttributeArgumentList()
		if al == nil {
			break
		}
		args := al.AllAttributeArgument()
		if len(args) >= 1 {
			module = stripQuotes(args[0].GetText())
		}
		if len(args) >= 2 {
			sym = stripQuotes(args[1].GetText())
		}
	}
	return
}

// ── Function declaration collectors ──────────────────────────────────────────

func (c *Compiler) collectExternFunc(fd *parser.FunctionDeclarationContext) error {
	name := fd.Identifier().GetText()
	if _, dup := c.funcMap[name]; dup {
		return nil // already registered (harmless re-declaration)
	}

	params, paramNames, ret, err := c.resolveFuncSignature(fd)
	if err != nil {
		return fmt.Errorf("extern func %s: %w", name, err)
	}

	ft := buildWasmFuncType(params, ret)
	typeIdx := c.internFuncType(ft)

	mod, sym := externNames(fd)
	c.mod.Imports.AddFunc(mod, sym, typeIdx)
	funcIdx := c.mod.Imports.NumFuncs() - 1

	c.funcMap[name] = &FuncInfo{
		Name:       name,
		WasmIdx:    funcIdx,
		TypeIdx:    typeIdx,
		Params:     params,
		ParamNames: paramNames,
		Ret:        ret,
		IsImport:   true,
	}
	if c.opts.Verbose {
		fmt.Printf("[pass1] extern %s → %s.%s\n", name, mod, sym)
	}
	return nil
}

func (c *Compiler) collectLocalFunc(fd *parser.FunctionDeclarationContext) error {
	name := fd.Identifier().GetText()
	if _, dup := c.funcMap[name]; dup {
		return fmt.Errorf("duplicate function %q", name)
	}

	params, paramNames, ret, err := c.resolveFuncSignature(fd)
	if err != nil {
		return fmt.Errorf("func %s: %w", name, err)
	}

	ft := buildWasmFuncType(params, ret)
	typeIdx := c.internFuncType(ft)

	funcIdx := c.mod.Imports.NumFuncs() + uint32(c.mod.Functions.Len())
	c.mod.Functions.Add(typeIdx)

	info := &FuncInfo{
		Name:       name,
		WasmIdx:    funcIdx,
		TypeIdx:    typeIdx,
		Params:     params,
		ParamNames: paramNames,
		Ret:        ret,
		IsImport:   false,
		ASTNode:    fd,
	}
	c.funcMap[name] = info
	c.funcOrder = append(c.funcOrder, name)

	// Automatically export `main`.
	if name == "main" {
		c.mod.Exports.Add("main", wasm.ExportFunc, funcIdx)
	}
	if c.opts.Verbose {
		fmt.Printf("[pass1] func %s → idx %d\n", name, funcIdx)
	}
	return nil
}

// ── Signature resolution ──────────────────────────────────────────────────────

func (c *Compiler) resolveFuncSignature(fd *parser.FunctionDeclarationContext) (
	params []VType, paramNames []string, ret VType, err error,
) {
	sig := fd.FunctionSignature()

	if pl := sig.ParameterClause().ParameterList(); pl != nil {
		for _, p := range pl.AllParameter() {
			ta := p.TypeAnnotation()
			if ta == nil {
				err = fmt.Errorf("parameter missing type annotation")
				return
			}
			var vt VType
			vt, err = c.resolveType(ta.Type_())
			if err != nil {
				return
			}
			params = append(params, vt)

			pname := "_"
			if lpn := p.LocalParameterName(); lpn != nil {
				if id := lpn.Identifier(); id != nil {
					pname = id.GetText()
				}
			}
			paramNames = append(paramNames, pname)
		}
	}

	ret = &PrimitiveType{Kind: KindVoid}
	if fr := sig.FunctionResult(); fr != nil {
		ret, err = c.resolveType(fr.Type_())
	}
	return
}

// ── Struct / class layout ─────────────────────────────────────────────────────

func (c *Compiler) collectStructDecl(sd *parser.StructDeclarationContext) error {
	name := sd.Identifier().GetText()
	st := &StructType{Name: name}
	var offset uint32
	for _, m := range sd.StructBody().AllStructMember() {
		decl := m.Declaration()
		if decl == nil {
			continue
		}
		fname, ftype, ok, err := c.extractFieldDecl(decl)
		if err != nil {
			return fmt.Errorf("struct %s: %w", name, err)
		}
		if !ok {
			continue
		}
		st.Fields = append(st.Fields, StructField{Name: fname, Type: ftype, Offset: offset})
		offset += ftype.Size()
	}
	c.structMap[name] = st
	return nil
}

func (c *Compiler) collectClassDecl(cd *parser.ClassDeclarationContext) error {
	name := cd.Identifier().GetText()
	st := &StructType{Name: name, IsRef: true}
	var offset uint32
	for _, m := range cd.ClassBody().AllClassMember() {
		decl := m.Declaration()
		if decl == nil {
			continue
		}
		fname, ftype, ok, err := c.extractFieldDecl(decl)
		if err != nil {
			return fmt.Errorf("class %s: %w", name, err)
		}
		if !ok {
			continue
		}
		st.Fields = append(st.Fields, StructField{Name: fname, Type: ftype, Offset: offset})
		offset += ftype.Size()
	}
	c.structMap[name] = st
	return nil
}

// extractFieldDecl pulls (name, type) from a stored-property declaration.
func (c *Compiler) extractFieldDecl(decl parser.IDeclarationContext) (
	name string, vtype VType, ok bool, err error,
) {
	// `var fieldName: Type` with explicit variable name
	if vd := decl.VariableDeclaration(); vd != nil {
		if vd.VariableName() != nil && vd.TypeAnnotation() != nil {
			name = vd.VariableName().Identifier().GetText()
			vtype, err = c.resolveType(vd.TypeAnnotation().Type_())
			ok = err == nil
			return
		}
		// `var fieldName: Type` via patternInitializerList
		if vd.PatternInitializerList() != nil {
			for _, pi := range vd.PatternInitializerList().AllPatternInitializer() {
				if ip, cast := pi.Pattern().(*parser.IdentPatContext); cast {
					name = ip.IdentifierPattern().Identifier().GetText()
					if ip.TypeAnnotation() != nil {
						vtype, err = c.resolveType(ip.TypeAnnotation().Type_())
						ok = err == nil
						return
					}
				}
			}
		}
	}
	// `let fieldName: Type`
	if cd := decl.ConstantDeclaration(); cd != nil {
		if cd.PatternInitializerList() != nil {
			for _, pi := range cd.PatternInitializerList().AllPatternInitializer() {
				if ip, cast := pi.Pattern().(*parser.IdentPatContext); cast {
					name = ip.IdentifierPattern().Identifier().GetText()
					if ip.TypeAnnotation() != nil {
						vtype, err = c.resolveType(ip.TypeAnnotation().Type_())
						ok = err == nil
						return
					}
				}
			}
		}
	}
	return "", nil, false, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func buildWasmFuncType(params []VType, ret VType) wasm.FuncType {
	ft := wasm.FuncType{}
	for _, p := range params {
		if wt, ok := p.WasmType(); ok {
			ft.Params = append(ft.Params, wt)
		}
	}
	if ret != nil && !ret.IsVoid() {
		if wt, ok := ret.WasmType(); ok {
			ft.Results = append(ft.Results, wt)
		}
	}
	return ft
}