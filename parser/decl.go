package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
)

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
		d := p.parseVarDecl()
		p.expectStmtEnd()
		return d
	}

	p.errorHere(diag.ExpectedDecl, p.describe(p.tok))
	bad := &ast.BadDecl{From: p.tok.Pos, To: p.tok.End()}
	p.advance(declStart)
	return bad
}

// parseImportDecl parses A.2.3. There is no aliasing form, no dot-import, and
// no blank import, so there is nothing to record but paths.
func (p *parser) parseImportDecl() *ast.ImportDecl {
	d := &ast.ImportDecl{Doc: p.leadComment, Import: p.expect(token.IMPORT)}

	if p.at(token.LPAREN) {
		d.Lparen = p.open(token.LPAREN)
		for !p.at(token.RPAREN) && !p.at(token.EOF) {
			d.Paths = append(d.Paths, p.parseImportPath())
		}
		d.Rparen = p.close(token.RPAREN)
	} else {
		d.Paths = append(d.Paths, p.parseImportPath())
	}
	p.expectStmtEnd()
	return d
}

func (p *parser) parseImportPath() *ast.BasicLit {
	if !p.at(token.STRING) {
		p.errorHere(diag.ExpectedToken, "an import path string", p.describe(p.tok))
		lit := &ast.BasicLit{ValuePos: p.tok.Pos, Kind: token.STRING, Value: `""`}
		p.advanceToken()
		return lit
	}
	lit := &ast.BasicLit{ValuePos: p.tok.Pos, Kind: token.STRING, Value: p.tok.Lit}
	p.advanceToken()
	return lit
}

// parseVarDecl parses A.5.1's three forms: `let` and `var` with initializers,
// and bare `var Binding`.
//
// Statement-leading `var` is always a declaration here. A.4.6 also lists ✗
// `var w` as a statement under "transfer outside an owning position", which
// would make it a TransferExpr instead; the two readings collide and this one
// is chosen because `var Binding` is a production of A.5.1 while a bare
// transfer statement is not a production of anything.
func (p *parser) parseVarDecl() *ast.VarDecl {
	d := &ast.VarDecl{Doc: p.leadComment, KwPos: p.tok.Pos, Kw: p.tok.Kind}
	p.advanceToken()

	for {
		b := &ast.Binding{Name: p.expectIdent()}
		if p.at(token.COLON) {
			b.Colon = p.tok.Pos
			p.advanceToken()
			b.Type = p.parseType()
		}
		d.Bindings = append(d.Bindings, b)
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
	}

	d.Comment = p.lineComment
	return d
}

// parseFuncDecl parses A.6.1, and with it A.6.4's initializer and
// deinitializer forms.
//
// Those need no separate path: A.1.3 makes `init` and `deinit` ContextualKeywords
// that are ordinary method names in a receiver declaration, so they arrive as
// IDENT and land in Name like any other method.
func (p *parser) parseFuncDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.FuncDecl{Doc: doc}
	funcPos := p.expect(token.FUNC)

	if p.at(token.LPAREN) {
		d.Recv = p.parseReceiver()
	}

	d.Name = p.expectIdent()

	if p.at(token.LBRACK) {
		// A.7.6 forbids a method declaring its own type parameter; the list is
		// parsed either way so the analyzer can say which rule was broken.
		d.TypeParams = p.parseTypeParamList()
	}

	ft := &ast.FuncType{Func: funcPos}
	ft.Params = p.parseParamList()
	ft.Marker = p.tryParseMarker()
	if p.at(token.ARROW) {
		ft.Arrow = p.tok.Pos
		p.advanceToken()
		ft.Result = p.parseType()
	}
	d.Type = ft

	if p.at(token.LBRACE) {
		d.Body = p.parseBlockStmt()
	} else {
		p.errorHere(diag.ExpectedToken, "'{'", p.describe(p.tok))
		p.advance(declStart)
	}
	p.expectStmtEnd()
	return d
}

// parseReceiver parses `( Identifier : ReceiverType )` (A.6.1).
func (p *parser) parseReceiver() *ast.Receiver {
	r := &ast.Receiver{}
	r.Lparen = p.open(token.LPAREN)
	r.Name = p.expectIdent()
	r.Colon = p.expect(token.COLON)
	r.Type = p.parseType()
	r.Rparen = p.close(token.RPAREN)
	return r
}

// parseRecordDecl parses both StructDeclaration and ClassDeclaration (A.6.2,
// A.6.3), which differ only in the keyword.
func (p *parser) parseRecordDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.RecordDecl{Doc: doc, KwPos: p.tok.Pos, Kw: p.tok.Kind}
	p.advanceToken()

	d.Name = p.expectIdent()
	if p.at(token.LBRACK) {
		d.TypeParams = p.parseTypeParamList()
	}

	// A.6.2 ⊢ fields are comma-separated and a line terminator between them is
	// conventional but not required, so this brace pushes depth rather than
	// resetting it.
	d.Lbrace = p.open(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		d.Fields = append(d.Fields, p.parseField())
		if !p.got(token.COMMA) {
			break
		}
	}
	d.Rbrace = p.close(token.RBRACE)
	p.expectStmtEnd()
	return d
}

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
		d.Variants = append(d.Variants, p.parseVariant())
		if !p.got(token.COMMA) {
			break
		}
	}
	d.Rbrace = p.close(token.RBRACE)
	p.expectStmtEnd()
	return d
}

// parseVariant parses A.6.5's three variant forms. An explicit discriminant on
// a payload variant parses and is rejected (A.14), so both suffixes are
// accepted here.
func (p *parser) parseVariant() *ast.Variant {
	v := &ast.Variant{Doc: p.leadComment}
	v.Name = p.expectIdent()

	if p.at(token.LPAREN) {
		v.Lparen = p.open(token.LPAREN)
		for !p.at(token.RPAREN) && !p.at(token.EOF) {
			v.Payload = append(v.Payload, p.parseType())
			if !p.got(token.COMMA) {
				break
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

func (p *parser) parseTypeAliasDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.TypeAliasDecl{Doc: doc, Type: p.expect(token.TYPE)}
	d.Name = p.expectIdent()
	if p.at(token.LBRACK) {
		d.TypeParams = p.parseTypeParamList()
	}
	d.Assign = p.expect(token.ASSIGN)
	d.Target = p.parseType()
	p.expectStmtEnd()
	return d
}

// parseConstraintDecl parses A.7.2. Vertex has no interfaces; a constraint is
// its own declaration form and is legal only in a `[...]` position, which is a
// static rule rather than anything this function can check.
func (p *parser) parseConstraintDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.ConstraintDecl{Doc: doc, Constraint: p.expect(token.CONSTRAINT)}
	d.Name = p.expectIdent()

	savedDepth := p.depth
	p.depth = 0

	d.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		e := &ast.ConstraintElem{}
		if p.at(token.FUNC) {
			e.Method = p.parseMethodReq()
		} else {
			// A.7.2 ⊢ a single identifier parses as both a one-term TypeSet and
			// a ConstraintName, and is resolved by what the name denotes. One
			// field holds both readings.
			e.Set = p.parseConstraintExpr()
		}
		d.Elems = append(d.Elems, e)
		p.expectStmtEnd()
	}
	d.Rbrace = p.expect(token.RBRACE)

	p.depth = savedDepth
	p.expectStmtEnd()
	return d
}

func (p *parser) parseMethodReq() *ast.MethodReq {
	m := &ast.MethodReq{Doc: p.leadComment, Func: p.expect(token.FUNC)}
	m.Name = p.expectIdent()
	m.Params = p.parseParamList()
	if p.at(token.ARROW) {
		m.Arrow = p.tok.Pos
		p.advanceToken()
		m.Result = p.parseType()
	}
	return m
}

// ---------------------------------------------------------- declare blocks

// parseDeclareDecl parses A.8.1's two block forms. `framework` and `module` are
// ContextualKeywords meaningful only immediately after `declare`.
func (p *parser) parseDeclareDecl(doc *ast.CommentGroup) ast.Decl {
	d := &ast.DeclareDecl{Doc: doc, Declare: p.expect(token.DECLARE)}

	d.KindPos = p.tok.Pos
	switch {
	case p.atCtx("framework"), p.atCtx("module"):
		d.Kind = p.tok.Lit
		p.advanceToken()
	default:
		p.errorHere(diag.ExpectedToken, "'framework' or 'module'", p.describe(p.tok))
		p.advance(declStart)
		return d
	}

	// A.8.2 ⊢ a framework block never takes a variant tag, but the tagged form
	// parses so the diagnostic can name the rule instead of the bracket.
	if p.at(token.LBRACK) {
		d.Variant = p.parseVariantTag()
	}

	d.Path = p.parseImportPath()

	savedDepth := p.depth
	p.depth = 0
	d.Lbrace = p.expect(token.LBRACE)
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		d.Members = append(d.Members, p.parseForeignMember())
	}
	d.Rbrace = p.expect(token.RBRACE)
	p.depth = savedDepth

	p.expectStmtEnd()
	return d
}

func (p *parser) parseVariantTag() *ast.VariantTag {
	v := &ast.VariantTag{}
	v.Lbrack = p.open(token.LBRACK)
	for !p.at(token.RBRACK) && !p.at(token.EOF) {
		v.Tags = append(v.Tags, p.parseImportPath())
		if !p.got(token.COMMA) {
			break
		}
	}
	v.Rbrack = p.close(token.RBRACK)
	return v
}

// parseForeignMember parses one member of a declare block (A.8.3).
//
// Leading identifiers are collected without judgement. `init` is a prefix
// modifier on func rather than a name, and a visibility modifier is banned but
// listed as a ✗ form — so both arrive as identifiers and are separated after
// the fact, letting `private init func() -> Bad` be diagnosed as itself rather
// than as a syntax error.
func (p *parser) parseForeignMember() ast.ForeignMember {
	doc := p.leadComment

	var mods []*ast.Ident
	for p.at(token.IDENT) {
		mods = append(mods, p.expectIdent())
	}

	var initPos token.Pos
	if n := len(mods); n > 0 && mods[n-1].Name == token.CtxInit {
		initPos = mods[n-1].Pos()
		mods = mods[:n-1]
	}

	switch {
	case p.at(token.FUNC):
		f := &ast.ForeignFunc{Doc: doc, Modifiers: mods, Init: initPos}
		f.Func = p.expect(token.FUNC)
		// A.8.3: the unnamed initializer form is what bare Type(...)
		// construction resolves to, so a missing name is grammatical here and
		// only here.
		if p.at(token.IDENT) {
			f.Name = p.expectIdent()
		} else if !initPos.IsValid() {
			f.Name = p.expectIdent()
		}
		f.Params = p.parseParamList()
		if p.at(token.ARROW) {
			f.Arrow = p.tok.Pos
			p.advanceToken()
			f.Result = p.parseType()
		}
		// A.8.3 ✗ declarations cannot have bodies. Parsed, then rejected.
		if p.at(token.LBRACE) {
			f.Body = p.parseBlockStmt()
		}
		p.expectStmtEnd()
		return f

	case p.at(token.CLASS):
		c := &ast.ForeignClass{Doc: doc, Modifiers: mods}
		c.Class = p.expect(token.CLASS)
		c.Name = p.expectIdent()
		c.Lbrace = p.expect(token.LBRACE)
		for !p.at(token.RBRACE) && !p.at(token.EOF) {
			c.Members = append(c.Members, p.parseForeignMember())
		}
		c.Rbrace = p.expect(token.RBRACE)
		p.expectStmtEnd()
		return c

	case p.at(token.COLON) && len(mods) > 0:
		// A.8.3 ✗ fields describe foreign-side layout and are banned. Parsed so
		// the diagnostic can point at the field rather than at a stray colon.
		f := &ast.ForeignField{Doc: doc, Name: mods[len(mods)-1]}
		f.Colon = p.expect(token.COLON)
		f.Type = p.parseType()
		p.expectStmtEnd()
		return f
	}

	p.errorHere(diag.ExpectedDecl, p.describe(p.tok))
	p.advance(map[token.Kind]bool{token.FUNC: true, token.CLASS: true, token.RBRACE: true})
	return &ast.ForeignFunc{Doc: doc, Modifiers: mods, Func: p.tok.Pos, Params: &ast.ParamList{}}
}