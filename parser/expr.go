package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
)

// open and close bracket the depth counter. Every group that suspends
// statement termination goes through these; a Block does not (see block()).
func (p *parser) open(k token.Kind) token.Pos { p.depth++; return p.expect(k) }
func (p *parser) close(k token.Kind) token.Pos { p.depth--; return p.expect(k) }

// noLitOff runs f with [+Lit] restored. Every bracketed group re-enters +Lit:
// A.4.1 makes parentheses the escape hatch for a literal in a header, and
// A.4.3's index brackets carry Expression[+Lit] explicitly.
func (p *parser) withLit(f func()) {
	saved := p.noLit
	p.noLit = false
	f()
	p.noLit = saved
}

// ------------------------------------------------------------ expressions

func (p *parser) parseExpr() ast.Expr { return p.parseBinaryExpr(1) }

// parseHeaderExpr parses a control-flow header under [~Lit] (A.4.7).
func (p *parser) parseHeaderExpr(kw string) ast.Expr {
	saved := p.noLit
	p.noLit = true
	x := p.parseExpr()
	p.noLit = saved

	// The literal is ungrammatical here rather than merely unparsed, so when a
	// brace does follow the header expression the diagnostic names the fix
	// A.4.7 prescribes instead of reporting a stray token.
	if p.at(token.LBRACE) && p.isTypeLike(x) && p.looksLikeLiteralBody() {
		d := diag.At(diag.LiteralInHeader, x.Pos(), kw)
		d.WithFixit(x.Pos(), x.Pos(), "(", "wrap the literal in parentheses")
		p.errCount++
		if p.rep != nil {
			p.rep.Report(d)
		}
	}
	return x
}

// looksLikeLiteralBody distinguishes `if C {` (a block) from `if C {a: 1}` (a
// misplaced literal) by one token of lookahead past the brace.
func (p *parser) looksLikeLiteralBody() bool {
	if !p.at(token.LBRACE) {
		return false
	}
	n := p.peek()
	return n.Kind == token.IDENT
}

func (p *parser) isTypeLike(x ast.Expr) bool {
	switch t := x.(type) {
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		return p.isTypeLike(t.X)
	case *ast.IndexExpr:
		return p.isTypeLike(t.X)
	}
	return false
}

// parseBinaryExpr is precedence climbing over A.13.
//
// `as` is deliberately absent from this loop even though it has a precedence:
// its right operand is a Type, not an Expr, so it is folded in parseCastExpr
// where the type parser is reachable. That matches the grammar, where
// CastExpression sits below ShiftExpression rather than inside the cascade.
func (p *parser) parseBinaryExpr(minPrec int) ast.Expr {
	x := p.parseCastExpr()

	for {
		if !p.continues() {
			return x
		}
		op := p.tok.Kind
		prec := op.Prec()
		if op == token.AS || prec == token.LowestPrec || prec < minPrec {
			return x
		}
		opPos := p.tok.Pos
		p.advanceToken()

		y := p.parseBinaryExpr(prec + 1)
		x = &ast.BinaryExpr{X: x, OpPos: opPos, Op: op, Y: y}

		// A.4.5 ⊢ `..` is non-associative; a..b..c is a compile error. The
		// second range is consumed so recovery continues at the statement
		// rather than at the operator.
		if op == token.DOTDOT && p.continues() && p.at(token.DOTDOT) {
			p.errorHere(diag.RangeNotAssociative)
			p.advanceToken()
			p.parseBinaryExpr(prec + 1)
		}
	}
}

// parseCastExpr folds `x as T` chains. A.4.4 ⊢ left-associative, so
// `value as int32 as int64` is two conversions.
func (p *parser) parseCastExpr() ast.Expr {
	x := p.parseUnaryExpr()
	for p.continues() && p.at(token.AS) {
		as := p.tok.Pos
		p.advanceToken()
		x = &ast.CastExpr{X: x, As: as, Type: p.parseType()}
	}
	return x
}

func (p *parser) parseUnaryExpr() ast.Expr {
	switch p.tok.Kind {
	case token.SUB, token.NOT, token.TILDE:
		op, pos := p.tok.Kind, p.tok.Pos
		p.advanceToken()
		return &ast.UnaryExpr{OpPos: pos, Op: op, X: p.parseUnaryExpr()}

	case token.AWAIT:
		// Parsed unconditionally. A.14 lists `await` outside an async body as a
		// form that parses and is rejected, so [+Await] is not tracked here.
		pos := p.tok.Pos
		p.advanceToken()
		return &ast.AwaitExpr{Await: pos, X: p.parseUnaryExpr()}

	case token.VAR:
		// The ownership marker (A.4.6). Its operand is parsed as a full unary
		// expression rather than restricted to a TransferTarget, so that
		// `var pick(a, b)` reaches the analyzer as a TransferExpr over a call
		// and earns A.14's "var on a computed expression" rather than a
		// syntax error.
		pos := p.tok.Pos
		p.advanceToken()
		return &ast.TransferExpr{Var: pos, Target: p.parseUnaryExpr()}
	}
	return p.parsePostfixExpr()
}

// parsePostfixExpr handles launch prefixes, then `&` chains, then the postfix
// operators.
func (p *parser) parsePostfixExpr() ast.Expr {
	if x := p.tryParseLaunch(); x != nil {
		return x
	}
	x := p.parsePointerPrimary()
	return p.parsePostfixOps(x)
}

// tryParseLaunch parses a LaunchExpression (A.4.2), or returns nil.
//
// The [lookahead != .] restriction is the whole reason for the peek: it is what
// keeps `npu.Dot(a, b)` a namespace member access rather than a launch of a
// member call, and A.4.2 applies it uniformly to async, gpu, and npu.
func (p *parser) tryParseLaunch() ast.Expr {
	kw := p.tok.Kind
	switch kw {
	case token.THREAD:
	case token.ASYNC, token.GPU, token.NPU:
		if p.peek().Kind == token.PERIOD {
			return nil
		}
	default:
		return nil
	}

	x := &ast.LaunchExpr{KwPos: p.tok.Pos, Kw: kw}
	p.advanceToken()

	if kw == token.GPU && p.at(token.LPAREN) {
		x.Config = p.parseLaunchConfig()
	}

	call := p.parsePostfixOps(p.parsePointerPrimary())
	if _, ok := call.(*ast.CallExpr); !ok {
		p.errorAt(diag.ExpectedExpr, call.Pos(), "a call")
	}
	x.Call = call
	return x
}

// parseLaunchConfig parses `(blocks: E, threads: E)` (A.4.2). Fixed arity and
// fixed names, so it is not a general argument list.
func (p *parser) parseLaunchConfig() *ast.LaunchConfig {
	c := &ast.LaunchConfig{}
	c.Lparen = p.open(token.LPAREN)
	p.withLit(func() {
		p.expectIdent() // blocks
		p.expect(token.COLON)
		c.Blocks = p.parseExpr()
		p.expect(token.COMMA)
		p.expectIdent() // threads
		p.expect(token.COLON)
		c.Threads = p.parseExpr()
	})
	c.Rparen = p.close(token.RPAREN)
	return c
}

// parsePointerPrimary implements A.4.3's PointerPrimary.
//
// A.4.3 ⊢ `&` binds tighter than member access, so it wraps a PrimaryExpression
// and the postfix loop runs outside it: `&p.add(1)` is `(&p).add(1)`.
func (p *parser) parsePointerPrimary() ast.Expr {
	if p.at(token.AND) {
		pos := p.tok.Pos
		p.advanceToken()
		return &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: p.parsePointerPrimary()}
	}
	return p.parsePrimaryExpr()
}

func (p *parser) parsePostfixOps(x ast.Expr) ast.Expr {
	for {
		if !p.continues() {
			return x
		}
		switch p.tok.Kind {
		case token.PERIOD:
			dot := p.tok.Pos
			p.advanceToken()
			if p.at(token.INT) {
				// Positional tuple access (A.4.3). The scanner has already
				// decided this is an index rather than a float, because a digit
				// after `.` can mean nothing else in Vertex.
				x = &ast.TupleIndexExpr{
					X: x, Dot: dot,
					IndexPos: p.tok.Pos,
					Index:    atoi(p.tok.Lit),
					Text:     p.tok.Lit,
				}
				p.advanceToken()
				continue
			}
			x = &ast.SelectorExpr{X: x, Dot: dot, Sel: p.expectIdent()}

		case token.LBRACK:
			ix := &ast.IndexExpr{X: x}
			ix.Lbrack = p.open(token.LBRACK)
			p.withLit(func() {
				for {
					ix.Indices = append(ix.Indices, p.parseExprOrType())
					if !p.got(token.COMMA) || p.at(token.RBRACK) {
						break
					}
				}
			})
			ix.Rbrack = p.close(token.RBRACK)
			x = ix

		case token.LPAREN:
			x = p.parseCallSuffix(x)

		case token.LBRACE:
			// CompositeLiteral (A.4.7). Suppressed under [~Lit], and never
			// begun across a line break at depth zero, since `x` on one line
			// followed by a block on the next is two constructs.
			if p.noLit || !p.isTypeLike(x) {
				return x
			}
			x = p.parseCompositeLitBody(x)

		default:
			return x
		}
	}
}

func (p *parser) parseCallSuffix(fun ast.Expr) ast.Expr {
	c := &ast.CallExpr{Fun: fun}
	c.Lparen = p.open(token.LPAREN)
	p.withLit(func() {
		c.Args = p.parseArgumentList(fun)
	})
	c.Rparen = p.close(token.RPAREN)
	return c
}

// parseArgumentList parses Arguments (A.4.3).
//
// Three ReservedBuiltinNames take a Type in argument position — sizeof, alignof,
// and the first argument of reinterpret (A.4.8). The parser recognizes them by
// name, which is sound only because A.1.4 ⊢ forbids shadowing or redeclaring a
// ReservedBuiltinName; without that guarantee this would be a hack.
func (p *parser) parseArgumentList(fun ast.Expr) []ast.Expr {
	typeFirst := ""
	if id, ok := fun.(*ast.Ident); ok {
		switch id.Name {
		case "sizeof", "alignof", "reinterpret":
			typeFirst = id.Name
		}
	}

	var args []ast.Expr
	for !p.at(token.RPAREN) && !p.at(token.EOF) {
		if typeFirst != "" && len(args) == 0 {
			args = append(args, p.parseType())
		} else {
			args = append(args, p.parseArgument())
		}
		if !p.got(token.COMMA) {
			break
		}
	}
	return args
}

// parseArgument parses one Argument: an OwningExpression, or `name: value` for
// the named form. A.14 makes mixing named and positional a static rejection, so
// both shapes land in the same slice here.
func (p *parser) parseArgument() ast.Expr {
	x := p.parseExpr()
	if id, ok := x.(*ast.Ident); ok && p.at(token.COLON) {
		colon := p.tok.Pos
		p.advanceToken()
		return &ast.KeyValueExpr{Key: id, Colon: colon, Value: p.parseExpr()}
	}
	return x
}

func (p *parser) parsePrimaryExpr() ast.Expr {
	pos := p.tok.Pos

	switch p.tok.Kind {
	case token.IDENT:
		x := &ast.Ident{NamePos: pos, Name: p.tok.Lit}
		p.advanceToken()
		return x

	case token.INT, token.FLOAT, token.CHAR, token.STRING,
		token.TRUE, token.FALSE, token.NIL:
		x := &ast.BasicLit{ValuePos: pos, Kind: p.tok.Kind, Value: p.tok.Lit}
		if x.Value == "" {
			x.Value = p.tok.Kind.String() // true / false / nil carry no Lit
		}
		p.advanceToken()
		return x

	case token.ASYNC, token.GPU, token.NPU, token.CHAN:
		// A NamespaceExpression (A.4.1). A launch prefix was already ruled out
		// by tryParseLaunch's lookahead, and `chan` reaches here as the head of
		// a construction such as chan[float32](64).
		x := &ast.NamespaceExpr{KwPos: pos, Kw: p.tok.Kind}
		p.advanceToken()
		return x

	case token.LPAREN:
		return p.parseParenOrTuple()

	case token.LBRACK:
		lit := &ast.ArrayLit{}
		lit.Lbrack = p.open(token.LBRACK)
		p.withLit(func() {
			for !p.at(token.RBRACK) && !p.at(token.EOF) {
				lit.Elems = append(lit.Elems, p.parseExpr())
				if !p.got(token.COMMA) {
					break
				}
			}
		})
		lit.Rbrack = p.close(token.RBRACK)
		return lit

	case token.LBRACE:
		if p.noLit {
			break
		}
		return p.parseMapLit()

	case token.PERIOD:
		return p.parseEnumShorthand()

	case token.FUNC:
		ft := p.parseFuncType()
		if p.at(token.LBRACE) {
			return &ast.FuncLit{Type: ft, Body: p.parseBlockStmt()}
		}
		return ft

	case token.MAP, token.TYPED_PTR, token.TENSOR, token.ABSTRACT,
		token.MUT, token.UNIQUE, token.SHARED, token.WEAK:
		// A type in expression position. It is grammatical in a few places
		// (a type argument, a sizeof operand) and a static error elsewhere;
		// either way the shape is a type, so parse it as one.
		return p.parseType()
	}

	p.errorHere(diag.ExpectedExpr, p.describe(p.tok))
	bad := &ast.BadExpr{From: pos, To: p.tok.End()}
	p.advance(stmtStart)
	return bad
}

// parseParenOrTuple resolves `(x)`, `(x,)`, `(a, b)`, and `()`.
//
// A.4.7 ⊢ a one-element tuple literal requires its trailing comma; `(1)` is a
// parenthesized integer. A.3.1 says the same distinction does not exist for the
// type form, where a parenthesized single type is that type — which falls out,
// since both produce the same node here.
func (p *parser) parseParenOrTuple() ast.Expr {
	lparen := p.open(token.LPAREN)

	if p.at(token.RPAREN) {
		rparen := p.close(token.RPAREN)
		return &ast.TupleExpr{Lparen: lparen, Rparen: rparen} // UnitType
	}

	var elems []ast.Expr
	trailing := false
	p.withLit(func() {
		for {
			elems = append(elems, p.parseTupleElement())
			if !p.got(token.COMMA) {
				break
			}
			trailing = true
			if p.at(token.RPAREN) {
				break
			}
			trailing = false
		}
	})
	rparen := p.close(token.RPAREN)

	if len(elems) == 1 && !trailing {
		return &ast.ParenExpr{Lparen: lparen, X: elems[0], Rparen: rparen}
	}
	return &ast.TupleExpr{
		Lparen: lparen, Elems: elems,
		TrailingComma: trailing, Rparen: rparen,
	}
}

func (p *parser) parseTupleElement() ast.Expr {
	x := p.parseExprOrType()
	if id, ok := x.(*ast.Ident); ok && p.at(token.COLON) {
		colon := p.tok.Pos
		p.advanceToken()
		return &ast.KeyValueExpr{Key: id, Colon: colon, Value: p.parseExprOrType()}
	}
	return x
}

func (p *parser) parseCompositeLitBody(typ ast.Expr) ast.Expr {
	lit := &ast.CompositeLit{Type: typ}
	lit.Lbrace = p.open(token.LBRACE)
	p.withLit(func() {
		for !p.at(token.RBRACE) && !p.at(token.EOF) {
			lit.Elems = append(lit.Elems, p.parseFieldValue())
			if !p.got(token.COMMA) {
				break
			}
		}
	})
	lit.Rbrace = p.close(token.RBRACE)
	return lit
}

// parseFieldValue parses a FieldValue (A.4.7), whose key must be an Identifier.
func (p *parser) parseFieldValue() ast.Expr {
	key := p.parseExpr()
	if !p.at(token.COLON) {
		return key
	}
	colon := p.tok.Pos
	p.advanceToken()
	return &ast.KeyValueExpr{Key: key, Colon: colon, Value: p.parseExpr()}
}

// parseMapLit parses a MapLiteral (A.4.7). Its keys are arbitrary Expressions,
// unlike a CompositeLiteral's field names.
func (p *parser) parseMapLit() ast.Expr {
	lit := &ast.MapLit{}
	lit.Lbrace = p.open(token.LBRACE)
	p.withLit(func() {
		for !p.at(token.RBRACE) && !p.at(token.EOF) {
			key := p.parseExpr()
			colon := p.expect(token.COLON)
			val := p.parseExpr()
			lit.Elems = append(lit.Elems, &ast.KeyValueExpr{Key: key, Colon: colon, Value: val})
			if !p.got(token.COMMA) {
				break
			}
		}
	})
	lit.Rbrace = p.close(token.RBRACE)
	return lit
}

func (p *parser) parseEnumShorthand() ast.Expr {
	x := &ast.EnumShorthand{Dot: p.tok.Pos}
	p.advanceToken()
	x.Name = p.expectIdent()
	if p.at(token.LPAREN) {
		x.Lparen = p.open(token.LPAREN)
		p.withLit(func() {
			for !p.at(token.RPAREN) && !p.at(token.EOF) {
				x.Args = append(x.Args, p.parseArgument())
				if !p.got(token.COMMA) {
					break
				}
			}
		})
		x.Rparen = p.close(token.RPAREN)
	}
	return x
}

// ------------------------------------------------------------------ types

// startsTypeOnly reports whether the current token can begin a type and cannot
// begin an expression.
func (p *parser) startsTypeOnly() bool {
	switch p.tok.Kind {
	case token.MAP, token.TYPED_PTR, token.TENSOR, token.ABSTRACT,
		token.MUT, token.UNIQUE, token.SHARED, token.WEAK, token.CHAN:
		return true
	}
	return false
}

// startsType reports whether the current token can begin a type at all. Used
// for the one-token decision that separates `[3]int32` from `[3]`.
func (p *parser) startsType() bool {
	if p.startsTypeOnly() {
		return true
	}
	switch p.tok.Kind {
	case token.IDENT, token.LPAREN, token.LBRACK, token.FUNC:
		return true
	}
	return false
}

// parseExprOrType parses a type-argument or index operand, which A.3.6 makes
// deliberately ambiguous: `Stack[int32]` and `a[i]` share bracket syntax and are
// resolved by what the operand denotes, not by shape.
func (p *parser) parseExprOrType() ast.Expr {
	if p.startsTypeOnly() {
		return p.parseType()
	}
	if p.at(token.LBRACK) {
		return p.parseBracketedTypeOrArray()
	}
	if p.at(token.VAR) {
		// In a type-argument position `var T` is the ownership qualifier, not
		// the transfer marker: A.9.1 does not list type arguments among the
		// owning positions.
		return p.parseType()
	}
	return p.parseExpr()
}

// parseBracketedTypeOrArray resolves `[` in a position that admits both an
// array literal and an array or slice type.
//
// `[]T` is unambiguous. Otherwise the decision is made after the closing
// bracket: a type following it means the brackets were a length, and anything
// else means they were elements.
func (p *parser) parseBracketedTypeOrArray() ast.Expr {
	lbrack := p.open(token.LBRACK)

	if p.at(token.RBRACK) {
		rbrack := p.close(token.RBRACK)
		return &ast.ArrayType{Lbrack: lbrack, Rbrack: rbrack, Elem: p.parseType()}
	}

	var elems []ast.Expr
	trailing := false
	p.withLit(func() {
		for !p.at(token.RBRACK) && !p.at(token.EOF) {
			elems = append(elems, p.parseExpr())
			if !p.got(token.COMMA) {
				break
			}
			trailing = true
		}
	})
	rbrack := p.close(token.RBRACK)

	if len(elems) == 1 && !trailing && p.startsType() {
		return &ast.ArrayType{Lbrack: lbrack, Len: elems[0], Rbrack: rbrack, Elem: p.parseType()}
	}
	return &ast.ArrayLit{Lbrack: lbrack, Elems: elems, Rbrack: rbrack}
}

func (p *parser) parseType() ast.Expr {
	pos := p.tok.Pos

	switch p.tok.Kind {
	case token.IDENT:
		// TypeName, QualifiedTypeName, InstantiatedType. PredeclaredTypeNames
		// arrive here as ordinary identifiers (A.1.4).
		var x ast.Expr = &ast.Ident{NamePos: pos, Name: p.tok.Lit}
		p.advanceToken()
		if p.at(token.PERIOD) {
			dot := p.tok.Pos
			p.advanceToken()
			x = &ast.SelectorExpr{X: x, Dot: dot, Sel: p.expectIdent()}
		}
		if p.at(token.LBRACK) {
			ix := &ast.IndexExpr{X: x}
			ix.Lbrack = p.open(token.LBRACK)
			for {
				ix.Indices = append(ix.Indices, p.parseType())
				if !p.got(token.COMMA) || p.at(token.RBRACK) {
					break
				}
			}
			ix.Rbrack = p.close(token.RBRACK)
			x = ix
		}
		return x

	case token.MUT, token.VAR, token.UNIQUE, token.SHARED, token.WEAK:
		// A.3.2. Qualifiers do not stack, but a stacked form parses and is
		// rejected (A.14), so the recursion is unguarded here.
		kw := p.tok.Kind
		p.advanceToken()
		return &ast.OwnershipType{KwPos: pos, Kw: kw, X: p.parseType()}

	case token.LBRACK:
		lbrack := p.open(token.LBRACK)
		if p.at(token.RBRACK) {
			rbrack := p.close(token.RBRACK)
			return &ast.ArrayType{Lbrack: lbrack, Rbrack: rbrack, Elem: p.parseType()}
		}
		var length ast.Expr
		p.withLit(func() { length = p.parseExpr() })
		rbrack := p.close(token.RBRACK)
		return &ast.ArrayType{Lbrack: lbrack, Len: length, Rbrack: rbrack, Elem: p.parseType()}

	case token.MAP:
		t := &ast.MapType{Map: pos}
		p.advanceToken()
		t.Lbrack = p.open(token.LBRACK)
		t.Key = p.parseType()
		t.Rbrack = p.close(token.RBRACK)
		t.Value = p.parseType()
		return t

	case token.CHAN:
		p.advanceToken()
		return &ast.ChanType{Chan: pos, Elem: p.parseType()}

	case token.TYPED_PTR:
		// A.3.3 ⊢ nesting requires parentheses, so a nested form arrives as a
		// ParenExpr and needs no special case.
		p.advanceToken()
		return &ast.PointerType{Kw: pos, Elem: p.parseType()}

	case token.TENSOR:
		t := &ast.TensorType{Tensor: pos}
		p.advanceToken()
		t.Lbrack = p.open(token.LBRACK)
		t.Elem = p.parseType()
		for p.got(token.COMMA) {
			if p.at(token.RBRACK) {
				break
			}
			var dim ast.Expr
			p.withLit(func() { dim = p.parseExpr() })
			t.Shape = append(t.Shape, dim)
		}
		t.Rbrack = p.close(token.RBRACK)
		return t

	case token.ABSTRACT:
		p.advanceToken()
		return &ast.AbstractType{Abstract: pos}

	case token.FUNC:
		return p.parseFuncType()

	case token.LPAREN:
		return p.parseParenOrTuple()
	}

	p.errorHere(diag.ExpectedType, p.describe(p.tok))
	bad := &ast.BadExpr{From: pos, To: p.tok.End()}
	p.advanceToken()
	return bad
}

// parseFuncType parses a FunctionType and the signature half of a declaration
// or literal. In a bare FunctionType every Param has a nil Name (A.3.4).
func (p *parser) parseFuncType() *ast.FuncType {
	ft := &ast.FuncType{Func: p.expect(token.FUNC)}
	ft.Params = p.parseParamList()
	ft.Marker = p.tryParseMarker()
	if p.at(token.ARROW) {
		ft.Arrow = p.tok.Pos
		p.advanceToken()
		ft.Result = p.parseType()
	}
	return ft
}

// tryParseMarker parses a FunctionMarker (A.6.1), or returns nil.
//
// A.6.1 permits at most one, but more than one is listed in A.14 as a static
// rejection, so only the first is attached and any extras are consumed for
// recovery. `test` is a ContextualKeyword and arrives as IDENT.
func (p *parser) tryParseMarker() *ast.Marker {
	var m *ast.Marker
	for {
		switch {
		case p.at(token.ASYNC), p.at(token.GPU), p.at(token.NPU):
			cur := &ast.Marker{MarkerPos: p.tok.Pos, Kind: p.tok.Kind, Name: p.tok.Kind.String()}
			p.advanceToken()
			if m == nil {
				m = cur
			}
		case p.atCtx(token.CtxTest):
			cur := &ast.Marker{MarkerPos: p.tok.Pos, Kind: token.IDENT, Name: token.CtxTest}
			p.advanceToken()
			if m == nil {
				m = cur
			}
		default:
			return m
		}
	}
}

func (p *parser) parseParamList() *ast.ParamList {
	l := &ast.ParamList{}
	l.Lparen = p.open(token.LPAREN)
	for !p.at(token.RPAREN) && !p.at(token.EOF) {
		l.List = append(l.List, p.parseParam())
		if !p.got(token.COMMA) {
			break
		}
	}
	l.Rparen = p.close(token.RPAREN)
	return l
}

// parseParam parses `name: T`, `name: ... T`, or a bare type. The bare form is
// what a FunctionType's TypeList produces.
func (p *parser) parseParam() *ast.Param {
	prm := &ast.Param{Doc: p.leadComment}

	if p.at(token.IDENT) && p.peek().Kind == token.COLON {
		prm.Name = p.expectIdent()
		prm.Colon = p.expect(token.COLON)
	}
	if p.at(token.ELLIPSIS) {
		prm.Ellipsis = p.tok.Pos
		p.advanceToken()
	}
	prm.Type = p.parseType()
	return prm
}

// parseTypeParamList parses a TypeParameterList (A.7.1).
//
// A constraint is attached only to the name it follows. A.7.1's rule that it
// also applies to every immediately preceding unconstrained name is
// distribution over an already-parsed list, and doing it here would erase the
// written form vfmt needs to reprint.
func (p *parser) parseTypeParamList() *ast.TypeParamList {
	l := &ast.TypeParamList{}
	l.Lbrack = p.open(token.LBRACK)
	for !p.at(token.RBRACK) && !p.at(token.EOF) {
		tp := &ast.TypeParam{Name: p.expectIdent()}
		if p.at(token.COLON) {
			tp.Colon = p.tok.Pos
			p.advanceToken()
			tp.Constraint = p.parseConstraintExpr()
		}
		l.List = append(l.List, tp)
		if !p.got(token.COMMA) {
			break
		}
	}
	l.Rbrack = p.close(token.RBRACK)
	return l
}

// parseConstraintExpr parses a ConstraintExpression: a type set, possibly a
// union, possibly with `~` terms (A.7.3).
func (p *parser) parseConstraintExpr() ast.Expr {
	x := p.parseTypeSetTerm()
	for p.at(token.OR) {
		opPos := p.tok.Pos
		p.advanceToken()
		x = &ast.BinaryExpr{X: x, OpPos: opPos, Op: token.OR, Y: p.parseTypeSetTerm()}
	}
	return x
}

func (p *parser) parseTypeSetTerm() ast.Expr {
	if p.at(token.TILDE) {
		// A.7.3 ⊢ `~` here is underlying-type, never bitwise-NOT. Same node as
		// the operator; a type-set element is not an expression position, so
		// the two never collide.
		pos := p.tok.Pos
		p.advanceToken()
		return &ast.UnaryExpr{OpPos: pos, Op: token.TILDE, X: p.parseType()}
	}
	return p.parseType()
}

func atoi(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			continue
		}
		n = n*10 + int(s[i]-'0')
	}
	return n
}