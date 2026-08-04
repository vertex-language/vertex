package parser

import (
	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
)

// open and close bracket the depth counter. Every group that suspends
// statement termination goes through these; the braces that do not — a Block, a
// record or constraint body, a switch or select body, a declare or foreign
// class body — go through enterTerminated instead.
func (p *parser) open(k token.Kind) token.Pos  { p.depth++; return p.expect(k) }
func (p *parser) close(k token.Kind) token.Pos { p.depth--; return p.expect(k) }

// withLit runs f with a suppressed literal re-enabled. Every bracketed group
// does this: parentheses are the escape hatch the grammar prescribes for a
// literal in a header, and an index bracket or argument list encloses the
// literal just as effectively.
func (p *parser) withLit(f func()) {
	saved := p.noLit
	p.noLit = false
	f()
	p.noLit = saved
}

// ------------------------------------------------------------ expressions

func (p *parser) parseExpr() ast.Expr { return p.parseBinaryExpr(1) }

// parseHeaderExpr parses a control-flow header, where an unparenthesized
// composite or map literal is ambiguous against the block brace that follows.
//
// The literal's brace is read as the block's, which is what the grammar
// prescribes; the fix is to parenthesize. A composite literal is caught here,
// after the header expression, by the two tokens that can only begin a field
// value; a map literal is caught in parsePrimaryExpr, where the brace is
// reached with an operand expected.
func (p *parser) parseHeaderExpr(kw string) ast.Expr {
	savedLit, savedKw := p.noLit, p.headerKw
	p.noLit, p.headerKw = true, kw
	x := p.parseExpr()
	p.noLit, p.headerKw = savedLit, savedKw

	if p.at(token.LBRACE) && p.isTypeLike(x) &&
		p.peekAt(0).Kind == token.IDENT && p.peekAt(1).Kind == token.COLON {
		p.reportLiteralInHeader(x.Pos(), kw)
	}
	return x
}

func (p *parser) reportLiteralInHeader(pos token.Pos, kw string) {
	d := diag.At(diag.LiteralInHeader, pos, kw)
	d.WithInsert(pos, "(", "parenthesize the literal")
	p.report(d)
}

// isTypeLike reports whether x has the shape of a LiteralType: a type name,
// possibly qualified, possibly instantiated.
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

// parseBinaryExpr is precedence climbing over the seven binary levels.
//
// `as` is deliberately absent from this loop even though it binds tighter than
// every operator in it: its right operand is a Type, not an Expression, so a
// precedence-climbing loop cannot consume it and must not try. It is folded in
// parseCastExpr, where the type parser is reachable.
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

		// `..` is non-associative: `a..b..c` is a compile error, folded
		// neither left nor right. The second range is consumed so recovery
		// continues at the statement rather than at the operator.
		if op.IsNonAssociative() && p.continues() && p.at(op) {
			p.errorHere(diag.RangeNotAssociative)
			p.advanceToken()
			p.parseBinaryExpr(prec + 1)
		}
	}
}

// parseCastExpr folds `x as T` chains, which are left-associative: two
// conversions, not one.
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
		// Parsed unconditionally. Whether the enclosing body licenses it is a
		// static rule, and a keyword needs no context to recognize.
		pos := p.tok.Pos
		p.advanceToken()
		return &ast.AwaitExpr{Await: pos, X: p.parseUnaryExpr()}

	case token.VAR:
		// The ownership marker. Its operand is written as a full unary
		// expression so that `var f(a)` and `var items[0]` parse and reach the
		// analyzer as transfers of a computed value rather than as syntax
		// errors. There is no TransferTarget production to restrict it to.
		pos := p.tok.Pos
		p.advanceToken()
		return &ast.TransferExpr{Var: pos, Target: p.parseUnaryExpr()}
	}
	return p.parsePostfixExpr()
}

// parsePostfixExpr handles a launch prefix, then the `&` chain, then the
// postfix operators.
func (p *parser) parsePostfixExpr() ast.Expr {
	if x := p.tryParseLaunch(); x != nil {
		return x
	}
	return p.parsePostfixOps(p.parsePointerPrimary())
}

// tryParseLaunch parses a launch expression, or returns nil.
//
// The one-token lookahead is the whole reason for the peek: it is what keeps
// `npu.Dot(a, b)` a namespace member access rather than a launch of a member
// call, and it applies uniformly to async, gpu, and npu. `thread` is not a
// namespace and needs no lookahead.
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

	// A launch config is written only on gpu. A config on another prefix has
	// no production, so `async (blocks: …)` reads as an ordinary call and is
	// diagnosed by whatever the call turns out to be.
	if kw == token.GPU && p.at(token.LPAREN) {
		x.Config = p.parseLaunchConfig()
	}

	call := p.parsePostfixOps(p.parsePointerPrimary())
	x.Call = p.requireCall(call, kw.Spelling())
	return x
}

// parseLaunchConfig parses `( "blocks" : E , "threads" : E )`. Fixed arity and
// fixed names, so it is not a general argument list and the names are not
// recorded.
func (p *parser) parseLaunchConfig() *ast.LaunchConfig {
	c := &ast.LaunchConfig{}
	c.Lparen = p.open(token.LPAREN)
	p.withLit(func() {
		if !p.atCtx(token.CtxBlocks) {
			p.errorHere(diag.ExpectedToken, "'blocks'", p.describe(p.tok))
		} else {
			p.advanceToken()
		}
		p.expect(token.COLON)
		c.Blocks = p.parseExpr()
		p.expect(token.COMMA)
		if !p.atCtx(token.CtxThreads) {
			p.errorHere(diag.ExpectedToken, "'threads'", p.describe(p.tok))
		} else {
			p.advanceToken()
		}
		p.expect(token.COLON)
		c.Threads = p.parseExpr()
	})
	c.Rparen = p.close(token.RPAREN)
	return c
}

// parsePointerPrimary implements `Operand | "&" PointerPrimary`.
//
// `&` binds tighter than member access, so it wraps an operand and the postfix
// loop runs outside it: `&p.add(1)` is `(&p).add(1)`.
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
				// Positional tuple access. The scanner has already decided
				// this is an index rather than a float, under the one
				// restriction to longest-match scanning, so the digits arrive
				// as their own token; that the spelling must be decimal and
				// free of `_` is a static rule.
				x = &ast.TupleIndexExpr{
					X: x, Dot: dot,
					IndexPos: p.tok.Pos,
					Text:     p.tok.Lit,
				}
				p.advanceToken()
				continue
			}
			x = &ast.SelectorExpr{X: x, Dot: dot, Sel: p.expectIdent()}

		case token.LBRACK:
			// Index and TypeArgs are one node. Which reading applies is
			// settled by what the operand denotes, not by shape, so the parser
			// records the brackets and stops.
			ix := &ast.IndexExpr{X: x}
			ix.Lbrack = p.open(token.LBRACK)
			p.withLit(func() {
				for !p.at(token.RBRACK) && !p.at(token.EOF) {
					before := p.tok.Pos
					ix.Indices = append(ix.Indices, p.parseExprOrType())
					if !p.got(token.COMMA) {
						break
					}
					if p.stalled(before) {
						continue
					}
				}
			})
			ix.Rbrack = p.close(token.RBRACK)
			x = ix

		case token.LPAREN:
			x = p.parseCallSuffix(x)

		case token.LBRACE:
			// A composite literal. Suppressed inside a control-flow header,
			// and never begun after a line break at depth zero, since a name
			// on one line followed by a block on the next is two constructs.
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

// parseArgumentList parses Arguments.
//
// Three reserved builtin names take a Type in argument position: sizeof and
// alignof take one, reinterpret takes one followed by an expression. The parser
// recognizes them by name, which is sound only because a reserved builtin name
// may not be shadowed or declared — without that guarantee this would be a
// hack. The type argument leaves its own trace in the tree, since `[3]int32`
// parses as an array type where an expression would have been an array literal.
func (p *parser) parseArgumentList(fun ast.Expr) []ast.Expr {
	typeFirst := false
	if id, ok := fun.(*ast.Ident); ok {
		typeFirst = token.IsTypeOperator(id.Name)
	}

	var args []ast.Expr
	for !p.at(token.RPAREN) && !p.at(token.EOF) {
		before := p.tok.Pos
		if typeFirst && len(args) == 0 {
			args = append(args, p.parseType())
		} else {
			args = append(args, p.parseArgument())
		}
		if !p.got(token.COMMA) {
			break
		}
		if p.stalled(before) {
			continue
		}
	}
	return args
}

// parseArgument parses one Argument: an owning expression, or `name: value`.
// Mixing named and positional arguments is a static rule, so both shapes land
// in the same slice.
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
			x.Value = p.tok.Kind.Spelling() // true / false / nil carry no Lit
		}
		p.advanceToken()
		return x

	case token.ASYNC, token.GPU, token.NPU:
		// A namespace name, which appears only as the operand of a selector.
		// A launch prefix was already ruled out by tryParseLaunch's lookahead.
		x := &ast.NamespaceExpr{KwPos: pos, Kw: p.tok.Kind}
		p.advanceToken()
		return x

	case token.CHAN:
		// `chan` is not a namespace: its expression form is the constructor,
		// and its type form writes no brackets, so the two never compete.
		return p.parseChanConstructor()

	case token.UNIQUE, token.SHARED, token.WEAK:
		// A heap constructor is spelled with a keyword and so cannot be an
		// ordinary call over a reserved name. Without its own node the shape
		// would collide: `unique (T)` is an ownership type over a
		// parenthesized type, and once types are expressions that is
		// indistinguishable from the constructor.
		if p.peek().Kind == token.LPAREN {
			return p.parseHeapConstructor()
		}
		return p.parseType()

	case token.LPAREN:
		return p.parseParenOrTuple()

	case token.LBRACK:
		lit := &ast.ArrayLit{}
		lit.Lbrack = p.open(token.LBRACK)
		p.withLit(func() {
			for !p.at(token.RBRACK) && !p.at(token.EOF) {
				before := p.tok.Pos
				lit.Elems = append(lit.Elems, p.parseExpr())
				if !p.got(token.COMMA) {
					break
				}
				if p.stalled(before) {
					continue
				}
			}
		})
		lit.Rbrack = p.close(token.RBRACK)
		return lit

	case token.LBRACE:
		if p.noLit {
			// A map literal reached with an operand expected inside a header.
			// It is ungrammatical here rather than merely unparsed, so the
			// diagnostic names the fix and the literal is parsed anyway,
			// leaving the following brace to open the block.
			p.reportLiteralInHeader(pos, p.headerKw)
			var x ast.Expr
			p.withLit(func() { x = p.parseMapLit() })
			return x
		}
		return p.parseMapLit()

	case token.PERIOD:
		return p.parseEnumShorthand()

	case token.FUNC:
		// A function literal begins with all enclosing parse context cleared
		// and re-establishes it from its own marker; the marker lives on the
		// signature, and the body's block resets termination state itself.
		ft := p.parseSignature(p.tok.Pos, false)
		if p.at(token.LBRACE) {
			return &ast.FuncLit{Type: ft, Body: p.parseBlockStmt()}
		}
		return ft

	case token.MAP, token.TYPED_PTR, token.TENSOR, token.VECTOR,
		token.ABSTRACT, token.MUT:
		// A type in expression position. It is grammatical in a few places — a
		// type argument, a sizeof operand, the callee of a vector call — and a
		// static error elsewhere; either way the shape is a type.
		return p.parseType()
	}

	p.errorHere(diag.ExpectedExpr, p.describe(p.tok))
	bad := &ast.BadExpr{From: pos, To: p.tok.End()}
	p.advance(stmtStart)
	return bad
}

// parseParenOrTuple resolves `(x)`, `(x,)`, and `(a, b)`.
//
// A one-element tuple requires its trailing comma; `(1)` is a parenthesized
// integer. The same distinction does not exist for the type form, where a
// parenthesized single type is that type — which falls out, since both produce
// the same node here. A tuple has at least one element and there is no unit
// type, so `()` has no production at all.
func (p *parser) parseParenOrTuple() ast.Expr {
	lparen := p.open(token.LPAREN)

	if p.at(token.RPAREN) {
		rparen := p.close(token.RPAREN)
		p.errorSpan(diag.EmptyTuple, lparen, rparen+1)
		return &ast.BadExpr{From: lparen, To: rparen + 1}
	}

	var elems []ast.Expr
	trailing := false
	p.withLit(func() {
		for {
			elems = append(elems, p.parseTupleElement())
			if !p.got(token.COMMA) {
				break
			}
			if p.at(token.RPAREN) || p.at(token.EOF) {
				trailing = true
				break
			}
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

// parseTupleElement parses `[ identifier ":" ] Type` and
// `[ identifier ":" ] OwningExpr` at once — the two are the same shape once
// types are expressions.
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
			before := p.tok.Pos
			lit.Elems = append(lit.Elems, p.parseFieldValue())
			if !p.got(token.COMMA) {
				break
			}
			if p.stalled(before) {
				continue
			}
		}
	})
	lit.Rbrace = p.close(token.RBRACE)
	return lit
}

// parseFieldValue parses `identifier ":" OwningExpr`. That the key must be an
// identifier is what separates a composite literal from a map literal; a
// non-identifier key parses and is rejected.
func (p *parser) parseFieldValue() ast.Expr {
	key := p.parseExpr()
	if !p.at(token.COLON) {
		return key
	}
	colon := p.tok.Pos
	p.advanceToken()
	return &ast.KeyValueExpr{Key: key, Colon: colon, Value: p.parseExpr()}
}

// parseMapLit parses a braced literal with no type prefix. Its keys are
// arbitrary expressions, unlike a composite literal's field names.
func (p *parser) parseMapLit() ast.Expr {
	lit := &ast.MapLit{}
	lit.Lbrace = p.open(token.LBRACE)
	p.withLit(func() {
		for !p.at(token.RBRACE) && !p.at(token.EOF) {
			before := p.tok.Pos
			key := p.parseExpr()
			colon := p.expect(token.COLON)
			val := p.parseExpr()
			lit.Elems = append(lit.Elems, &ast.KeyValueExpr{Key: key, Colon: colon, Value: val})
			if !p.got(token.COMMA) {
				break
			}
			if p.stalled(before) {
				continue
			}
		}
	})
	lit.Rbrace = p.close(token.RBRACE)
	return lit
}

// parseEnumShorthand parses `.identifier` or `.identifier(args)` in expression
// position. That the enum type must be fixed by context is a static rule.
func (p *parser) parseEnumShorthand() ast.Expr {
	x := &ast.EnumShorthand{Dot: p.tok.Pos}
	p.advanceToken()
	x.Name = p.expectIdent()
	if p.at(token.LPAREN) {
		x.Lparen = p.open(token.LPAREN)
		p.withLit(func() {
			for !p.at(token.RPAREN) && !p.at(token.EOF) {
				before := p.tok.Pos
				x.Args = append(x.Args, p.parseArgument())
				if !p.got(token.COMMA) {
					break
				}
				if p.stalled(before) {
					continue
				}
			}
		})
		x.Rparen = p.close(token.RPAREN)
	}
	return x
}

// parseChanConstructor parses `chan [ Type ] ( [ Expression ] )`, the only
// expression form of chan. The optional argument is the capacity.
func (p *parser) parseChanConstructor() ast.Expr {
	c := &ast.ChanConstructor{Chan: p.expect(token.CHAN)}
	c.Lbrack = p.open(token.LBRACK)
	c.Elem = p.parseType()
	c.Rbrack = p.close(token.RBRACK)
	c.Lparen = p.open(token.LPAREN)
	p.withLit(func() {
		if !p.at(token.RPAREN) && !p.at(token.EOF) {
			c.Cap = p.parseExpr()
		}
	})
	c.Rparen = p.close(token.RPAREN)
	return c
}

// parseHeapConstructor parses `unique(x)`, `shared(x)`, or `weak(x)`.
func (p *parser) parseHeapConstructor() ast.Expr {
	h := &ast.HeapConstructor{KwPos: p.tok.Pos, Kw: p.tok.Kind}
	p.advanceToken()
	h.Lparen = p.open(token.LPAREN)
	p.withLit(func() { h.X = p.parseExpr() })
	h.Rparen = p.close(token.RPAREN)
	return h
}

// ------------------------------------------------------------------ types

// startsTypeOnly reports whether the current token can begin a type and cannot
// begin an expression.
func (p *parser) startsTypeOnly() bool {
	switch p.tok.Kind {
	case token.MAP, token.TYPED_PTR, token.TENSOR, token.VECTOR,
		token.ABSTRACT, token.MUT, token.UNIQUE, token.SHARED, token.WEAK,
		token.CHAN, token.FUNC:
		return true
	}
	return false
}

// startsType reports whether the current token can begin a type at all. It is
// the one-token decision that separates `[3]int32` from `[3]`.
func (p *parser) startsType() bool {
	if p.startsTypeOnly() {
		return true
	}
	switch p.tok.Kind {
	case token.IDENT, token.LPAREN, token.LBRACK:
		return true
	}
	return false
}

// parseExprOrType parses an index operand or a type argument, which the grammar
// makes deliberately ambiguous: `Stack[int32]` and `a[i]` share bracket syntax
// and are resolved by what the operand denotes, not by shape.
//
// `chan` resolves toward the type reading here. Its constructor form writes an
// argument list after the brackets, but a channel of array type writes brackets
// too, and a type argument is the far likelier reading in this position.
func (p *parser) parseExprOrType() ast.Expr {
	if p.startsTypeOnly() {
		return p.parseType()
	}
	if p.at(token.LBRACK) {
		return p.parseBracketedTypeOrArray()
	}
	if p.at(token.VAR) {
		// In a type-argument position `var T` is the ownership qualifier, not
		// the transfer marker: a type argument is not an owning position.
		return p.parseType()
	}
	return p.parseExpr()
}

// parseBracketedTypeOrArray resolves a `[` in a position admitting both an
// array literal and an array or slice type.
//
// `[]T` is unambiguous. Otherwise the decision is made after the closing
// bracket: a type following it means the brackets held a length, and anything
// else means they held elements.
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
			before := p.tok.Pos
			elems = append(elems, p.parseExpr())
			if !p.got(token.COMMA) {
				break
			}
			trailing = true
			if p.stalled(before) {
				continue
			}
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
		// A type name, qualified or instantiated. Predeclared type names,
		// tensor element names, and constraint names all arrive here as
		// ordinary identifiers.
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
			for !p.at(token.RBRACK) && !p.at(token.EOF) {
				before := p.tok.Pos
				ix.Indices = append(ix.Indices, p.parseType())
				if !p.got(token.COMMA) {
					break
				}
				if p.stalled(before) {
					continue
				}
			}
			ix.Rbrack = p.close(token.RBRACK)
			x = ix
		}
		return x

	case token.MUT, token.VAR, token.UNIQUE, token.SHARED, token.WEAK:
		// Qualifiers do not stack, but a stacked form parses and is rejected,
		// so the recursion is unguarded here.
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
		// A channel type carries no direction and writes no brackets.
		p.advanceToken()
		return &ast.ChanType{Chan: pos, Elem: p.parseType()}

	case token.TYPED_PTR:
		// One may not be the direct base of another, so a nested form is
		// written with parentheses and arrives with a parenthesized element.
		// The recursion is unguarded so the unparenthesized form still parses.
		p.advanceToken()
		return &ast.PointerType{Kw: pos, Elem: p.parseType()}

	case token.TENSOR:
		return p.parseTensorType()

	case token.VECTOR:
		return p.parseVectorType()

	case token.ABSTRACT:
		p.advanceToken()
		return &ast.AbstractType{Abstract: pos}

	case token.FUNC:
		return p.parseSignature(pos, false)

	case token.LPAREN:
		return p.parseParenOrTuple()
	}

	p.errorHere(diag.ExpectedType, p.describe(p.tok))
	bad := &ast.BadExpr{From: pos, To: p.tok.End()}
	p.advanceToken()
	return bad
}

// parseTensorType parses `tensor [ ElementType , ShapeList ]`. Legal only
// inside an npu-marked function; elsewhere it parses and is rejected.
func (p *parser) parseTensorType() ast.Expr {
	t := &ast.TensorType{Tensor: p.expect(token.TENSOR)}
	t.Lbrack = p.open(token.LBRACK)
	t.Elem = p.parseType()
	for p.got(token.COMMA) {
		if p.at(token.RBRACK) || p.at(token.EOF) {
			break
		}
		var dim ast.Expr
		p.withLit(func() { dim = p.parseExpr() })
		t.Shape = append(t.Shape, dim)
	}
	if len(t.Shape) == 0 {
		p.errorHere(diag.ExpectedToken, "a shape list", p.describe(p.tok))
	}
	t.Rbrack = p.close(token.RBRACK)
	return t
}

// parseVectorType parses `vector [ ElementType , int_lit ]`. As the callee of a
// call it makes that call a vector call; no ordinary call reading applies.
func (p *parser) parseVectorType() ast.Expr {
	t := &ast.VectorType{Vector: p.expect(token.VECTOR)}
	t.Lbrack = p.open(token.LBRACK)
	t.Elem = p.parseType()
	t.Comma = p.expect(token.COMMA)
	p.withLit(func() { t.Len = p.parseExpr() })
	t.Rbrack = p.close(token.RBRACK)
	return t
}

// ------------------------------------------------------------- signatures

// parseSignature parses `Parameters { FunctionMarker } [ Result ]`, sharing one
// routine across every construct that names a function shape.
//
// declResult admits the Expected result form, which reaches the grammar only
// through a function or method declaration. That is what keeps an Expected
// result out of a function type or a function literal syntactically.
//
// A signature carries at most one marker, but the repetition is written so that
// more than one parses; all of them are kept and the extras are rejected later.
func (p *parser) parseSignature(funcPos token.Pos, declResult bool) *ast.FuncType {
	ft := &ast.FuncType{Func: funcPos}
	if p.at(token.FUNC) {
		ft.Func = p.expect(token.FUNC)
	}
	ft.Params = p.parseParamList()
	ft.Markers = p.parseMarkers()

	if p.at(token.ARROW) {
		ft.Arrow = p.tok.Pos
		p.advanceToken()
		if declResult && p.atCtx(token.CtxExpected) && p.peek().Kind == token.LPAREN {
			ft.Result = p.parseExpectedType()
		} else {
			ft.Result = p.parseType()
		}
	}
	return ft
}

// parseExpectedType parses the test result form. `Expected` and `error` are
// ordinary identifiers, so this is a call node with an identifier callee; its
// arity and argument shape are static rules.
func (p *parser) parseExpectedType() ast.Expr {
	fun := p.expectIdent()
	c := &ast.CallExpr{Fun: fun}
	c.Lparen = p.open(token.LPAREN)
	p.withLit(func() {
		for !p.at(token.RPAREN) && !p.at(token.EOF) {
			before := p.tok.Pos
			if p.at(token.STRING) {
				c.Args = append(c.Args, p.parseImportPath())
			} else {
				c.Args = append(c.Args, p.parseType())
			}
			if !p.got(token.COMMA) {
				break
			}
			if p.stalled(before) {
				continue
			}
		}
	})
	c.Rparen = p.close(token.RPAREN)
	return c
}

// parseMarkers collects every FunctionMarker written. `test` is a contextual
// keyword and arrives as an identifier, which is why the node records a name
// alongside a kind.
func (p *parser) parseMarkers() []*ast.Marker {
	var out []*ast.Marker
	for {
		switch {
		case p.at(token.ASYNC), p.at(token.GPU), p.at(token.NPU):
			out = append(out, &ast.Marker{
				MarkerPos: p.tok.Pos, Kind: p.tok.Kind, Name: p.tok.Kind.Spelling(),
			})
			p.advanceToken()
		case p.atCtx(token.CtxTest):
			out = append(out, &ast.Marker{
				MarkerPos: p.tok.Pos, Kind: token.IDENT, Name: token.CtxTest,
			})
			p.advanceToken()
		default:
			return out
		}
	}
}

func (p *parser) parseParamList() *ast.ParamList {
	l := &ast.ParamList{}
	l.Lparen = p.open(token.LPAREN)
	for !p.at(token.RPAREN) && !p.at(token.EOF) {
		before := p.tok.Pos
		l.List = append(l.List, p.parseParam())
		if !p.got(token.COMMA) {
			break
		}
		if p.stalled(before) {
			continue
		}
	}
	l.Rparen = p.close(token.RPAREN)
	return l
}

// parseParam parses `[ identifier ":" ] [ "..." ] Type`. The bare-type form is
// what a function type's parameter list produces; that names must be either all
// present or all absent within one list is a static rule, so a mixed list
// parses.
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

// parseTypeParamList parses TypeParameters.
//
// A constraint is attached only to the entry it follows. The rule that it also
// applies to every immediately preceding unconstrained entry is distribution
// over an already-parsed list; performing it here would erase the written form
// a formatter needs to reproduce.
func (p *parser) parseTypeParamList() *ast.TypeParamList {
	l := &ast.TypeParamList{}
	l.Lbrack = p.open(token.LBRACK)
	for !p.at(token.RBRACK) && !p.at(token.EOF) {
		before := p.tok.Pos
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
		if p.stalled(before) {
			continue
		}
	}
	l.Rbrack = p.close(token.RBRACK)
	return l
}

// parseConstraintExpr parses a TypeSet: one or more terms joined by `|`.
func (p *parser) parseConstraintExpr() ast.Expr {
	x := p.parseTypeSetTerm()
	for p.at(token.OR) {
		opPos := p.tok.Pos
		p.advanceToken()
		x = &ast.BinaryExpr{X: x, OpPos: opPos, Op: token.OR, Y: p.parseTypeSetTerm()}
	}
	return x
}

// parseTypeSetTerm parses `Type` or `"~" Type`.
//
// `~` here is underlying-type, never bitwise-NOT. It is the same node as the
// operator, because a type-set element is not an expression position and the
// two never collide; `~` outside a type set is a static rejection.
func (p *parser) parseTypeSetTerm() ast.Expr {
	if p.at(token.TILDE) {
		pos := p.tok.Pos
		p.advanceToken()
		return &ast.UnaryExpr{OpPos: pos, Op: token.TILDE, X: p.parseType()}
	}
	return p.parseType()
}