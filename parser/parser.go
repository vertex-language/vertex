// Package parser builds a Vertex syntax tree from source.
//
// Recursive descent with precedence climbing over the binary chain in
// vertex_grammar.md B.4. It always returns a non-nil file, so a partial parse
// is still usable (compiler_frontend.md §6).
//
// One goal symbol, no pre-scan, no per-file mode (§1). The Mode bitset below
// does not select a goal symbol and never will.
package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/scanner"
	"github.com/vertex-language/vertex/token"
)

// Mode is a bitset of parse options. It does not select a goal symbol — §1
// forbids that.
type Mode uint

const (
	// ParseComments retains comments in ast.File.Comments. Without it, only
	// NLBefore survives (§4.3). The scanner emits COMMENT tokens either way;
	// this controls whether the parser keeps them.
	ParseComments Mode = 1 << iota

	// ImportsOnly stops after the import prologue, so a build driver can
	// discover the dependency graph without paying for full parses. It stops
	// the parser, not the scanner — the file is still tokenized whole (§4.1).
	ImportsOnly

	// Trace prints a production trace (§6.1), for debugging.
	Trace
)

// ParseFile scans and parses file.
//
// The signature is the whole contract: one *token.File in, one tree and one
// diagnostic slice out. No import resolver, no goal-symbol argument, no
// any-typed source parameter (§8.1).
//
// The returned tree is never nil. A caller that checks for nil is writing dead
// code; a caller that discards the tree because diags is non-empty is throwing
// away the recovery work in §6.4.
func ParseFile(file *token.File, mode Mode) (*ast.File, []token.Diagnostic) {
	toks, diags := scanner.Scan(file)
	return ParseFileTokens(file, toks, mode)
}

// ParseFileTokens parses an already-scanned buffer, for tools that scanned for
// their own reasons (§8.5) and want the tree too. The parser buffers anyway
// (§4.1), so handing it a buffer costs nothing.
//
// diags from the scan are not included in the result; the caller already has
// them. Merging them here would double-report for ParseFile, which passes its
// own scan diagnostics through separately.
func ParseFileTokens(file *token.File, toks []token.Token, mode Mode) (*ast.File, []token.Diagnostic) {
	p := newParser(file, toks, mode)
	f := p.parseFile()
	return f, p.finish()
}

// ParseExpr parses a single expression, for repls, debuggers, --eval flags,
// and hover providers (§8.6).
func ParseExpr(file *token.File) (ast.Expr, []token.Diagnostic) {
	toks, _ := scanner.Scan(file)
	p := newParser(file, toks, 0)
	x := p.parseExpr(allowIn)
	p.expectEOF("expression")
	// The arena outlives this call: the caller has no ast.File to Release, so
	// fragment parses are deliberately GC-managed rather than arena-freed.
	p.arena.detach()
	return x, p.finish()
}

// ParseTypeExpr parses a single type.
//
// Two entry points rather than one, because §5.1 gives types their own
// hierarchy: there is no node that could serve as both return types, and no
// way to guess which side of the language a fragment is on.
func ParseTypeExpr(file *token.File) (ast.TypeExpr, []token.Diagnostic) {
	toks, _ := scanner.Scan(file)
	p := newParser(file, toks, 0)
	t := p.parseType()
	p.expectEOF("type")
	p.arena.detach()
	return t, p.finish()
}

// ---------------------------------------------------------------------------

const (
	// maxDepth caps recursion. Exceeding it is a diagnostic, never a hang
	// (§6.1). Deeply nested but legal input exists — `((((...))))` in generated
	// code — so the cap is generous and the message says what to do.
	maxDepth = 1000

	// maxResync bounds repeated recovery at one position before the parser
	// bails to the enclosing synchronization set (§6.4).
	maxResync = 1
)

type parser struct {
	file *token.File
	mode Mode

	// toks excludes comments. Filtering them here rather than skipping them at
	// every lookahead is safe because the scanner leaves NLBefore set on the
	// following real token as well as on the comment itself (§4.3) — dropping a
	// comment cannot lose a line break, which is what expectSemi depends on.
	toks     []token.Token
	comments []token.Token
	i        int

	diags []token.Diagnostic
	arena *arena

	depth int

	// Recovery state (§6.4).
	lastSync  int
	syncCount int

	// Trace state.
	indent int
}

func newParser(file *token.File, toks []token.Token, mode Mode) *parser {
	p := &parser{file: file, mode: mode, arena: newArena()}

	p.toks = make([]token.Token, 0, len(toks))
	for _, t := range toks {
		if t.Kind == token.COMMENT {
			if mode&ParseComments != 0 {
				p.comments = append(p.comments, t)
			}
			continue
		}
		p.toks = append(p.toks, t)
	}
	if len(p.toks) == 0 {
		// scanner.Scan always ends in EOF, so this only fires for a
		// hand-constructed empty buffer. Synthesize one rather than
		// bounds-check every access.
		end := file.PosAt(file.Size())
		p.toks = append(p.toks, token.Token{Kind: token.EOF, Pos: end, End: end})
	}
	return p
}

func (p *parser) finish() []token.Diagnostic {
	token.SortDiagnostics(p.diags)
	return p.diags
}

// --- token access -----------------------------------------------------------

func (p *parser) cur() token.Token  { return p.toks[p.i] }
func (p *parser) kind() token.Kind  { return p.toks[p.i].Kind }
func (p *parser) pos() token.Pos    { return p.toks[p.i].Pos }
func (p *parser) end() token.Pos    { return p.toks[p.i].End }
func (p *parser) nlBefore() bool    { return p.toks[p.i].NLBefore() }
func (p *parser) atEOF() bool       { return p.toks[p.i].Kind == token.EOF }

// peek returns the token n ahead, clamped to EOF.
func (p *parser) peek(n int) token.Token {
	j := p.i + n
	if j >= len(p.toks) {
		return p.toks[len(p.toks)-1]
	}
	return p.toks[j]
}

func (p *parser) at(k token.Kind) bool { return p.kind() == k }

// atCtx reports whether the current token is an identifier spelling a
// contextual keyword. Never a string compare — that is what the perfect hash in
// §3 exists to prevent.
func (p *parser) atCtx(c token.Ctx) bool { return p.cur().IsContextual(c) }

func (p *parser) next() token.Token {
	t := p.toks[p.i]
	if p.i < len(p.toks)-1 {
		p.i++
	}
	return t
}

func (p *parser) got(k token.Kind) bool {
	if p.at(k) {
		p.next()
		return true
	}
	return false
}

func (p *parser) gotCtx(c token.Ctx) bool {
	if p.atCtx(c) {
		p.next()
		return true
	}
	return false
}

// expect consumes k or reports a diagnostic without consuming. Returning NoPos
// on failure lets callers keep building a node with a short span rather than
// bailing, which is what keeps recovery productive.
func (p *parser) expect(k token.Kind) token.Pos {
	if p.at(k) {
		return p.next().Pos
	}
	p.errorf(p.cur(), "expected %s, found %s", k, p.describe(p.cur()))
	return token.NoPos
}

func (p *parser) expectCtx(c token.Ctx) token.Pos {
	if p.atCtx(c) {
		return p.next().Pos
	}
	p.errorf(p.cur(), "expected %s, found %s", c, p.describe(p.cur()))
	return token.NoPos
}

func (p *parser) expectEOF(what string) {
	if !p.atEOF() {
		p.errorf(p.cur(), "unexpected %s after %s", p.describe(p.cur()), what)
	}
}

// describe names a token for a message. IDENT renders as its source text,
// which is the only place the parser reads token text at all — and it is for
// human output, never for a decision.
func (p *parser) describe(t token.Token) string {
	switch {
	case t.Kind == token.EOF:
		return "end of input"
	case t.Kind == token.IDENT:
		return fmt.Sprintf("%q", p.file.Slice(t))
	case t.Kind.IsLiteral():
		return strings.ToLower(t.Kind.String()) + " literal"
	default:
		return "`" + t.Kind.String() + "`"
	}
}

// --- diagnostics ------------------------------------------------------------

func (p *parser) errorf(t token.Token, format string, args ...any) {
	p.diags = append(p.diags, token.Diagnostic{
		Pos: t.Pos,
		End: t.End,
		Msg: fmt.Sprintf(format, args...),
	})
}

func (p *parser) errorAt(pos, end token.Pos, format string, args ...any) {
	p.diags = append(p.diags, token.Diagnostic{
		Pos: pos, End: end, Msg: fmt.Sprintf(format, args...),
	})
}

// --- speculation (§6.1) -----------------------------------------------------

// mark is a checkpoint: a token index, a diagnostic count, and an arena mark.
// Rollback truncates all three, which makes speculation side-effect-free by
// construction — the token buffer is immutable, no diagnostic escapes, and no
// node leaks.
type mark struct {
	tok   int
	diags int
	arena arenaMark
}

func (p *parser) mark() mark {
	return mark{tok: p.i, diags: len(p.diags), arena: p.arena.mark()}
}

func (p *parser) reset(m mark) {
	p.i = m.tok
	p.diags = p.diags[:m.diags]
	p.arena.reset(m.arena)
}

// speculate runs f from a checkpoint, keeping its effects only if it returns
// true. The bool is the commit signal, not a success signal: instantiation
// expressions commit on a lookahead test that succeeds *after* the inner parse
// already did (§6.1).
func (p *parser) speculate(f func() bool) bool {
	m := p.mark()
	if f() {
		return true
	}
	p.reset(m)
	return false
}

// enter guards recursion depth. Exceeding the cap is a diagnostic, never a
// hang (§6.1).
func (p *parser) enter() bool {
	p.depth++
	if p.depth > maxDepth {
		p.errorf(p.cur(), "expression nests too deeply (limit %d)", maxDepth)
		return false
	}
	return true
}

func (p *parser) leave() { p.depth-- }

// --- semicolon insertion (§6.2) --------------------------------------------

// expectSemi consults, in order: an explicit `;`, NLBefore on the current
// token, `}`, and end of file. It applies to exactly the nonterminals listed as
// TerminatedByASI in grammar L.
//
// Returns the position of a written `;`, or NoPos when the terminator was
// inserted — so the node's span stops at the last real token rather than
// covering a character that is not there.
func (p *parser) expectSemi() token.Pos {
	if p.at(token.SEMI) {
		return p.next().Pos
	}
	if p.nlBefore() || p.at(token.RBRACE) || p.atEOF() {
		return token.NoPos
	}
	p.errorf(p.cur(), "expected `;` or a line break, found %s", p.describe(p.cur()))
	return token.NoPos
}

// expectMemberSep is expectSemi for the seven type-member signatures.
//
// TypeBody (G.3) admits `;` or `,` as a TypeMemberSeparator, so newline-
// separated interface bodies are entirely this rule's job (§6.2). Omitting
// PropertySignature from the ASI list breaks `interface U { id: string ⏎ name:
// string }` while methods keep working, which is exactly the bug this function
// exists to prevent.
func (p *parser) expectMemberSep() token.Pos {
	if p.at(token.SEMI) || p.at(token.COMMA) {
		return p.next().Pos
	}
	if p.nlBefore() || p.at(token.RBRACE) || p.atEOF() {
		return token.NoPos
	}
	p.errorf(p.cur(), "expected `;`, `,`, or a line break between members, found %s",
		p.describe(p.cur()))
	return token.NoPos
}

// noLineTerminator enforces one entry in NoLineTerminatorRestriction (grammar
// L). Each site calls it by name so the table in §7's tests maps one to one.
func (p *parser) noLineTerminator(what string) bool {
	if p.nlBefore() {
		p.errorf(p.cur(), "no line break allowed before %s", what)
		return false
	}
	return true
}

// --- recovery (§6.4) --------------------------------------------------------

// syncSet is a per-context set of tokens the parser can resume at.
type syncSet uint8

const (
	syncDecl   syncSet = 1 << iota // declaration starts and `}`
	syncStmt                       // statement starts, `;`, `}`
	syncMember                     // class or struct member starts, `}`
	syncTypeMember
)

func (p *parser) inSync(s syncSet) bool {
	switch p.kind() {
	case token.EOF, token.RBRACE:
		return true
	case token.SEMI:
		return s&(syncStmt|syncMember|syncTypeMember) != 0
	case token.VAR, token.CONST, token.CLASS, token.FUNCTION, token.IMPORT, token.EXPORT:
		return true
	case token.IF, token.FOR, token.WHILE, token.DO, token.SWITCH, token.TRY,
		token.THROW, token.RETURN, token.BREAK, token.CONTINUE, token.DEBUGGER,
		token.LBRACE:
		return s&syncStmt != 0
	case token.AT:
		return true // a decorator starts a declaration or a member
	case token.IDENT:
		switch p.cur().Ctx {
		case token.CtxStruct, token.CtxInterface, token.CtxType, token.CtxNamespace,
			token.CtxDeclare, token.CtxModule, token.CtxGlobal, token.CtxLet,
			token.CtxKernel, token.CtxGraph, token.CtxAbstract, token.CtxUsing:
			return true
		case token.CtxStatic, token.CtxReadonly, token.CtxPublic, token.CtxProtected,
			token.CtxPrivate, token.CtxOverride, token.CtxAccessor, token.CtxGet,
			token.CtxSet:
			return s&(syncMember|syncTypeMember) != 0
		}
	}
	return false
}

// advanceTo skips to the next token in s. A guard caps repeated resync at one
// position, after which the parser bails to the enclosing set by returning
// false — the caller must then propagate rather than loop.
func (p *parser) advanceTo(s syncSet) bool {
	if p.i == p.lastSync {
		p.syncCount++
		if p.syncCount > maxResync {
			p.syncCount = 0
			p.lastSync = -1
			return false
		}
	} else {
		p.lastSync = p.i
		p.syncCount = 0
	}

	for !p.atEOF() && !p.inSync(s) {
		p.skipBalanced()
	}
	return true
}

// skipBalanced advances one token, or past a whole bracketed group, so that
// recovery does not stop inside a nested construct at a token that only looks
// like a synchronization point.
func (p *parser) skipBalanced() {
	switch p.kind() {
	case token.LPAREN, token.LBRACK, token.LBRACE:
		open := p.kind()
		close := token.RPAREN
		if open == token.LBRACK {
			close = token.RBRACK
		} else if open == token.LBRACE {
			close = token.RBRACE
		}
		depth := 0
		for !p.atEOF() {
			switch p.kind() {
			case open:
				depth++
			case close:
				depth--
				if depth == 0 {
					p.next()
					return
				}
			}
			p.next()
		}
	default:
		p.next()
	}
}

// advanced asserts that an unbounded loop made progress and force-advances if
// it did not, so no input can spin (§6.4). Every loop over a member or element
// list calls this.
func (p *parser) advanced(before int) bool {
	if p.i > before {
		return true
	}
	p.next()
	return false
}

// --- trace ------------------------------------------------------------------

func (p *parser) trace(name string) func() {
	if p.mode&Trace == 0 {
		return func() {}
	}
	pos := p.file.Position(p.pos())
	fmt.Fprintf(os.Stderr, "%s%s (%s, %s)\n",
		strings.Repeat(". ", p.indent), name, pos, p.describe(p.cur()))
	p.indent++
	return func() { p.indent-- }
}

// --- file -------------------------------------------------------------------

// parseFile parses SourceFile: ModuleItemList_opt (§1). Every file is parsed
// the same way; [+Await] is unconditional at file scope.
func (p *parser) parseFile() *ast.File {
	defer p.trace("File")()

	f := &ast.File{FileEnd: p.toks[len(p.toks)-1].End, Arena: p.arena}

	for !p.atEOF() {
		before := p.i

		if p.mode&ImportsOnly != 0 && !p.atPrologueItem() {
			break
		}

		item := p.parseModuleItem()
		if item != nil {
			f.Items = append(f.Items, item)
		}
		if !p.advanced(before) {
			continue
		}
	}

	if p.mode&ParseComments != 0 {
		f.Comments = p.groupComments()
	}
	return f
}

// atPrologueItem reports whether the current token can begin an item that
// contributes to the dependency graph, for ImportsOnly.
//
// This covers `export ... from` as well as `import`, because a re-export is a
// dependency edge and a driver that missed them would build an incomplete
// graph. See the note in the review.
func (p *parser) atPrologueItem() bool {
	switch p.kind() {
	case token.IMPORT, token.EXPORT:
		return true
	}
	return false
}

func (p *parser) groupComments() []*ast.CommentGroup {
	var out []*ast.CommentGroup
	var cur *ast.CommentGroup
	for _, t := range p.comments {
		c := &ast.Comment{Slash: t.Pos, Tail: t.End}
		if cur == nil || t.NLBefore() && len(cur.List) > 0 && p.blankLineBetween(cur.List[len(cur.List)-1], c) {
			cur = &ast.CommentGroup{}
			out = append(out, cur)
		}
		cur.List = append(cur.List, c)
	}
	return out
}

func (p *parser) blankLineBetween(a, b *ast.Comment) bool {
	return p.file.Position(b.Pos()).Line > p.file.Position(a.End()).Line+1
}