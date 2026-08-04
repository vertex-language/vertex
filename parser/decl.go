package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
)

// parseTopLevelDecl parses one TopLevelDecl. Top-level declarations are
// order-independent, so there is nothing to track across calls.
func (p *parser) parseTopLevelDecl() ast.Decl {
	doc := p.leadComment

	switch p.tok.Kind {
	case token.FUNC:
		return p.parseFuncDecl(doc)
	case token.STRUCT, token.CLASS:
		return p.parseRecordDecl(doc)
	case token.ENUM:
		return p.parseEnumDecl(doc)
	case token.TYPE:
		return p.parseTypeAliasDecl(doc)
	case token.CONSTRAINT:
		return p.parseConstraintDecl(doc)
	case token.DECLARE:
		return p.parseDeclareDecl(doc)
	case token.LET, token.VAR:
		// Also a Statement. That a top-level initializer must be
		// compile-time-evaluable, and that the bare form is rejected here, are
		// static rules over this node.
		return p.parseVarDecl(doc)
	}

	p.errorHere(diag.ExpectedDecl, p.describe(p.tok))
	bad := &ast.BadDecl{From: p.tok.Pos, To: p.tok.End()}
	p.advance(declStart)
	return bad
}

// parseImportDecl parses `import "path"` or `import ( ... )`. There is no
// aliasing form, no dot-import, and no blank import, so there is nothing to
// record but paths.
func (p *parser) parseImportDecl() *ast.ImportDecl {
	d := &ast.ImportDecl{Doc: p.leadComment, Import: p.expect(token.IMPORT)}

	if p.at(token.LPAREN) {
		d.Lparen = p.open(token.LPAREN)
		for !p.at(token.RPAREN) && !p.at(token.EOF) {
			before := p.tok.Pos
			d.Paths = append(d.Paths, p.parseImportPath())
			if p.stalled(before) {
				continue
			}
		}
		d.Rparen = p.close(token.RPAREN)
	} else {
		d.Paths = append(d.Paths, p.parseImportPath())
	}
	p.expectTerminator()
	return d
}

func (p *parser) parseImportPath() *ast.BasicLit {
	if !p.at(token.STRING) {
		p.errorHere(diag.ExpectedToken, "an import path string", p.describe(p.tok))
		lit := &ast.BasicLit{ValuePos: p.tok.Pos, Kind: token.STRING, Value: `""`}
		return lit
	}
	lit := &ast.BasicLit{ValuePos: p.tok.Pos, Kind: token.STRING, Value: p.tok.Lit}
	p.advanceToken()
	return lit
}

// parseVarDecl parses both VarDecl forms: `let`/`var` with initializers, and
// bare `var Binding`.
//
// Statement-leading `var` is always a declaration. A bare `var w` would
// otherwise read as the transfer marker outside an owning position; the two
// readings collide and this one is chosen because `var Binding` is a production
// while a bare transfer statement is not a production of anything. The other
// reading is then diagnosed against a real declaration node.
func (p *parser) parseVarDecl(doc *ast.CommentGroup) *ast.VarDecl {
	d := &ast.VarDecl{Doc: doc, KwPos: p.tok.Pos, Kw: p.tok.Kind}
	p.advanceToken()

	for {
		d.Bindings = append(d.Bindings, p.parseBinding())
		if !p.got(token.COMMA) {
			break
		}
	}

	if p.at(token.ASSIGN) {
		d.Assign = p.tok.Pos
		p.advanceToken()
		for {
			d.Values = append(d.Values, p.parseExpr())
			if !p.got(token.COMMA) {
				break
			}
		}
	} else if d.Kw == token.LET || len(d.Bindings) > 1 {
		// The initializer-free alternative is `"var" Binding` and nothing
		// else: `let` always requires one, and the bare form takes a single
		// binding.
		p.errorHere(diag.ExpectedToken, "'='", p.describe(p.tok))
	}

	d.Comment = p.lineComment
	return d
}

func (p *parser) parseBinding() *ast.Binding {
	b := &ast.Binding{Name: p.expectIdent()}
	if p.at(token.COLON) {
		b.Colon = p.tok.Pos
		p.advanceToken()
		b.Type = p.parseType()
	}
	return b
}

// parseFuncDecl parses FunctionDecl and MethodDecl, and with them the
// initializer and deinitializer forms.
//
// Those need no separate path: `init` and `deinit` are contextual keywords that
// are ordinary method names in a receiver declaration, so they arrive as
// identifiers and land in Name like any other. Whether a given declaration is
// one is a question about its name and receiver.
func (p *parser) parseFuncDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.FuncDecl{Doc: doc}
	funcPos := p.expect(token.FUNC)

	if p.at(token.LPAREN) {
		d.Recv = p.parseReceiver()
	}

	d.Name = p.expectIdent()

	if p.at(token.LBRACK) {
		// A method may not declare its own type parameters. The slot is
		// parsed either way, so the diagnostic can point a caret at the
		// bracket list rather than report a syntax error.
		d.TypeParams = p.parseTypeParamList()
	}

	d.Type = p.parseSignature(funcPos, true)

	if p.at(token.LBRACE) {
		d.Body = p.parseBlockStmt()
	} else {
		p.errorHere(diag.ExpectedToken, "'{'", p.describe(p.tok))
		p.advance(declStart)
	}
	return d
}

// parseReceiver parses `( identifier : ReceiverType )`.
//
// A ReceiverType's own bracket list re-declares the receiver type's existing
// names rather than introducing fresh ones; it arrives as an IndexExpr, which
// is the same brackets in another position.
func (p *parser) parseReceiver() *ast.Receiver {
	r := &ast.Receiver{}
	r.Lparen = p.open(token.LPAREN)
	r.Name = p.expectIdent()
	r.Colon = p.expect(token.COLON)
	r.Type = p.parseType()
	r.Rparen = p.close(token.RPAREN)
	return r
}

// parseRecordDecl parses StructDecl and ClassDecl, which differ only in the
// keyword. A class is byte-for-byte identical in layout to a struct.
//
// A field list is newline-separated juxtaposition rather than a comma list, so
// this brace is terminator-significant and two fields on one line do not parse.
func (p *parser) parseRecordDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.RecordDecl{Doc: doc, KwPos: p.tok.Pos, Kw: p.tok.Kind}
	p.advanceToken()

	d.Name = p.expectIdent()
	if p.at(token.LBRACK) {
		d.TypeParams = p.parseTypeParamList()
	}

	saved := p.enterTerminated()
	d.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		before := p.tok.Pos
		d.Fields = append(d.Fields, p.parseField())
		p.expectTerminator()
		if p.stalled(before) {
			continue
		}
	}
	d.Rbrace = p.expect(token.RBRACE)
	p.leave(saved)
	return d
}

// parseField parses a FieldDecl. The default is evaluated at construction for
// any omitted field.
func (p *parser) parseField() *ast.Field {
	f := &ast.Field{Doc: p.leadComment}
	f.Name = p.expectIdent()
	f.Colon = p.expect(token.COLON)
	f.Type = p.parseType()
	if p.at(token.ASSIGN) {
		f.Assign = p.tok.Pos
		p.advanceToken()
		f.Default = p.parseExpr()
	}
	f.Comment = p.lineComment
	return f
}

// parseEnumDecl parses an enum declaration. Its body is a comma-separated
// variant list and is not terminator-significant, which is what lets a variant
// list span lines — so the brace pushes depth rather than resetting it.
func (p *parser) parseEnumDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.EnumDecl{Doc: doc, Enum: p.expect(token.ENUM)}
	d.Name = p.expectIdent()

	if p.at(token.LBRACK) {
		d.TypeParams = p.parseTypeParamList()
	}
	if p.at(token.COLON) {
		d.Colon = p.tok.Pos
		p.advanceToken()
		d.Discrim = p.parseType()
	}

	d.Lbrace = p.open(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		before := p.tok.Pos
		d.Variants = append(d.Variants, p.parseVariant())
		if !p.got(token.COMMA) {
			break
		}
		if p.stalled(before) {
			continue
		}
	}
	d.Rbrace = p.close(token.RBRACE)
	return d
}

// parseVariant parses one enum variant. Both suffixes are accepted on any
// variant, so an explicit discriminant on a payload variant parses and can be
// diagnosed as itself.
func (p *parser) parseVariant() *ast.Variant {
	v := &ast.Variant{Doc: p.leadComment}
	v.Name = p.expectIdent()

	if p.at(token.LPAREN) {
		v.Lparen = p.open(token.LPAREN)
		for !p.at(token.RPAREN) && !p.at(token.EOF) {
			before := p.tok.Pos
			v.Payload = append(v.Payload, p.parseType())
			if !p.got(token.COMMA) {
				break
			}
			if p.stalled(before) {
				continue
			}
		}
		v.Rparen = p.close(token.RPAREN)
	}
	if p.at(token.ASSIGN) {
		v.Assign = p.tok.Pos
		p.advanceToken()
		v.Value = p.parseExpr()
	}
	v.Comment = p.lineComment
	return v
}

// parseTypeAliasDecl parses `type Name[params] = AliasTarget`. A target of
// `abstract` makes the alias nominal and opaque.
func (p *parser) parseTypeAliasDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.TypeAliasDecl{Doc: doc, Type: p.expect(token.TYPE)}
	d.Name = p.expectIdent()
	if p.at(token.LBRACK) {
		d.TypeParams = p.parseTypeParamList()
	}
	d.Assign = p.expect(token.ASSIGN)
	d.Target = p.parseType()
	return d
}

// parseConstraintDecl parses a constraint declaration. There are no
// interfaces; a constraint is its own declaration form, and that it is legal
// only in a bracket position is a static rule. Elements are one per line and
// form an intersection, so this body is terminator-significant.
func (p *parser) parseConstraintDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.ConstraintDecl{Doc: doc, Constraint: p.expect(token.CONSTRAINT)}
	d.Name = p.expectIdent()

	saved := p.enterTerminated()
	d.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		before := p.tok.Pos
		e := &ast.ConstraintElem{}
		if p.at(token.FUNC) {
			e.Method = p.parseMethodReq()
		} else {
			// A single identifier parses as both a one-term type set and a
			// constraint name, and is resolved by what the name denotes. One
			// field holds both readings.
			e.Set = p.parseConstraintExpr()
		}
		d.Elems = append(d.Elems, e)
		p.expectTerminator()
		if p.stalled(before) {
			continue
		}
	}
	d.Rbrace = p.expect(token.RBRACE)
	p.leave(saved)
	return d
}

// parseMethodReq parses a MethodRequirement. It takes a full Signature, so a
// constraint can require a marked method.
func (p *parser) parseMethodReq() *ast.MethodReq {
	m := &ast.MethodReq{Doc: p.leadComment}
	funcPos := p.expect(token.FUNC)
	m.Func = funcPos
	m.Name = p.expectIdent()
	m.Type = p.parseSignature(funcPos, false)
	return m
}

// ---------------------------------------------------------- declare blocks

// parseDeclareDecl parses both declare block forms. `framework` and `module`
// are contextual keywords meaningful only immediately after `declare`.
func (p *parser) parseDeclareDecl(doc *ast.CommentGroup) *ast.DeclareDecl {
	d := &ast.DeclareDecl{Doc: doc, Declare: p.expect(token.DECLARE)}

	d.KindPos = p.tok.Pos
	switch {
	case p.atCtx(token.CtxFramework), p.atCtx(token.CtxModule):
		d.Kind = p.tok.Lit
		p.advanceToken()
	default:
		p.errorHere(diag.ExpectedToken, "'framework' or 'module'", p.describe(p.tok))
		p.advance(declStart)
		return d
	}

	// The variant tag is hoisted out of the module form, so a tagged framework
	// block parses and is rejected with a message about `declare framework`
	// rather than as a syntax error at the bracket.
	if p.at(token.LBRACK) {
		d.Variant = p.parseVariantTag()
	}

	d.Path = p.parseImportPath()

	saved := p.enterTerminated()
	d.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		before := p.tok.Pos
		d.Members = append(d.Members, p.parseForeignMember())
		p.expectTerminator()
		if p.stalled(before) {
			continue
		}
	}
	d.Rbrace = p.expect(token.RBRACE)
	p.leave(saved)
	return d
}

// parseVariantTag parses the bracketed tag set. The set is closed; membership
// is a static rule, so any string list parses here.
func (p *parser) parseVariantTag() *ast.VariantTag {
	v := &ast.VariantTag{}
	v.Lbrack = p.open(token.LBRACK)
	for !p.at(token.RBRACK) && !p.at(token.EOF) {
		before := p.tok.Pos
		v.Tags = append(v.Tags, p.parseImportPath())
		if !p.got(token.COMMA) {
			break
		}
		if p.stalled(before) {
			continue
		}
	}
	v.Rbrack = p.close(token.RBRACK)
	return v
}

// parseForeignMember parses one member of a declare body or a foreign class
// body.
//
// A declare body admits a foreign function, a foreign class, and a nested
// declare; a foreign class body admits a foreign function, a foreign
// initializer, and a field. A nested declare and a field each parse and are
// rejected, so both land here — a declare block describes call shapes only, and
// the diagnostic should name the construct rather than fail at a token.
func (p *parser) parseForeignMember() ast.ForeignMember {
	doc := p.leadComment

	// `init` is a prefix modifier on func, not a function name.
	var initPos token.Pos
	if p.atCtx(token.CtxInit) && p.peek().Kind == token.FUNC {
		initPos = p.tok.Pos
		p.advanceToken()
	}

	switch {
	case p.at(token.FUNC):
		return p.parseForeignFunc(doc, initPos)

	case p.at(token.CLASS):
		return p.parseForeignClass(doc)

	case p.at(token.DECLARE):
		return p.parseDeclareDecl(doc)

	case p.at(token.IDENT) && p.peek().Kind == token.COLON:
		// A field describes foreign-side layout and is banned. Parsed so the
		// diagnostic can point at the field rather than at a stray colon.
		return p.parseField()
	}

	p.errorHere(diag.ExpectedDecl, p.describe(p.tok))
	p.advance(memberStart)
	return &ast.ForeignFunc{
		Doc:  doc,
		Func: p.tok.Pos,
		Type: &ast.FuncType{Func: p.tok.Pos, Params: &ast.ParamList{}},
	}
}

// parseForeignFunc parses a foreign function or a foreign initializer.
//
// Name is nil only for the unnamed initializer form that bare `Type(...)`
// construction resolves to. A body and a marker are both rejected forms, and
// both are parsed: the body needs a node to hang the diagnostic on, and the
// marker already has one, since the signature keeps every marker written.
func (p *parser) parseForeignFunc(doc *ast.CommentGroup, initPos token.Pos) *ast.ForeignFunc {
	f := &ast.ForeignFunc{Doc: doc, Init: initPos}
	f.Func = p.expect(token.FUNC)

	if p.at(token.IDENT) || !initPos.IsValid() {
		f.Name = p.expectIdent()
	}

	f.Type = p.parseSignature(f.Func, false)

	if p.at(token.LBRACE) {
		f.Body = p.parseBlockStmt()
	}
	return f
}

func (p *parser) parseForeignClass(doc *ast.CommentGroup) *ast.ForeignClass {
	c := &ast.ForeignClass{Doc: doc, Class: p.expect(token.CLASS)}
	c.Name = p.expectIdent()

	saved := p.enterTerminated()
	c.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		before := p.tok.Pos
		c.Members = append(c.Members, p.parseForeignMember())
		p.expectTerminator()
		if p.stalled(before) {
			continue
		}
	}
	c.Rbrace = p.expect(token.RBRACE)
	p.leave(saved)
	return c
}