// Package parser turns Vertex source into an ast.File or an ast.Package.
//
// Two grammar-wide mechanisms shape everything here.
//
// Statement termination is driven by a bracket-nesting depth rather than by a
// token. A line terminator is a terminator if and only if the innermost
// enclosing bracketing construct is one whose production writes `terminator`.
// `(`, `[`, and every brace that does not open such a construct push depth;
// the ones that do — a Block, a struct or class body, a constraint body, a
// switch or select body, a declare body, a foreign class body — reset depth to
// zero instead, and the source file starts there. At depth zero a token
// carrying NLBefore ends the statement and may not continue a postfix or binary
// chain. Because a run of terminators is one terminator, `[ terminator ]` at
// the head of a list needs no code at all: consecutive line breaks collapse
// into the one flag on the following token.
//
// The literal-in-header ambiguity is the only place the parser must decide
// something the grammar leaves to prose. A CompositeLit or MapLit written
// unparenthesized between a control-flow keyword and its block brace is read as
// the block, and the fix is to parenthesize; `noLit` carries that state and the
// diagnostic names the fix.
//
// Everything else the grammar forbids is parsed and left alone. Forms that
// "parse and are rejected" — a stacked ownership qualifier, `tensor` outside an
// npu body, `var` on a computed expression, a body on a foreign declaration, a
// marker repeated on one signature — reach the analyzer as real nodes so the
// diagnostic can name the construct instead of pointing at a token.
package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/scanner"
	"github.com/vertex-language/vertex/token"
)

type Mode uint

const (
	// PackageClauseOnly stops after the package and build clauses.
	// Load-bearing for the loader, not an optimization: the qualifier under
	// which an imported package's symbols are reached comes from that
	// package's own package clause, so a file's names cannot resolve until
	// every imported directory's package line has been read.
	PackageClauseOnly Mode = 1 << iota

	// ImportsOnly stops after the import declarations.
	ImportsOnly

	// ParseComments retains comments in the tree.
	ParseComments
)

type parser struct {
	file *token.File
	scan scanner.Scanner
	rep  diag.Reporter
	mode Mode

	tok token.Token

	// ahead is raw lookahead, comments included and in source order. Keeping
	// comments in the buffer rather than discarding them at peek time is what
	// lets a lookahead decision leave comment attachment untouched.
	ahead []token.Token

	comments    []*ast.CommentGroup
	leadComment *ast.CommentGroup
	lineComment *ast.CommentGroup

	// noLit suppresses a bare composite or map literal, and headerKw names the
	// construct whose header set it. Both exist only for the one ambiguity the
	// grammar resolves in prose.
	noLit    bool
	headerKw string

	// depth is bracket nesting. See the package comment.
	depth int

	errCount int
	syncPos  token.Pos
	syncCnt  int
}

func (p *parser) init(file *token.File, src []byte, rep diag.Reporter, mode Mode) {
	p.file = file
	p.rep = rep
	p.mode = mode
	p.scan.Init(file, src, rep, scanner.ScanComments)
	p.advanceToken()
}

// ------------------------------------------------------------- token flow

// scanRaw returns the next token from the lookahead buffer, or from the
// scanner when the buffer is empty.
func (p *parser) scanRaw() token.Token {
	if len(p.ahead) > 0 {
		t := p.ahead[0]
		p.ahead = p.ahead[1:]
		return t
	}
	return p.scan.Scan()
}

// advanceToken moves to the next non-comment token, collecting comment groups.
//
// A skipped comment's NLBefore is folded into the token that follows it.
// Without that, a line break landing on a COMMENT would vanish with the
// COMMENT, and the statement it terminated would run on.
func (p *parser) advanceToken() {
	p.leadComment = nil
	p.lineComment = nil

	t := p.scanRaw()

	nl := false
	for t.Kind == token.COMMENT {
		nl = nl || t.NLBefore
		g, endedLine := p.consumeCommentGroup(&t, nl)
		if g != nil {
			if !nl && !endedLine {
				p.lineComment = g
			} else {
				p.leadComment = g
			}
			if p.mode&ParseComments != 0 {
				p.comments = append(p.comments, g)
			}
		}
		nl = nl || endedLine
	}
	t.NLBefore = t.NLBefore || nl

	p.tok = t
}

// consumeCommentGroup gathers a run of comments with no intervening token and
// no blank line, leaving *t on the first token that is neither. endedLine
// reports whether the group was followed by a line terminator.
func (p *parser) consumeCommentGroup(t *token.Token, leading bool) (g *ast.CommentGroup, endedLine bool) {
	var list []*ast.Comment
	for t.Kind == token.COMMENT {
		c := &ast.Comment{Slash: t.Pos, Text: t.Lit}
		list = append(list, c)
		if c.ActsAsTerminator() {
			endedLine = true
		}
		*t = p.scanRaw()
		if t.NLBefore {
			endedLine = true
			if t.Kind != token.COMMENT {
				break
			}
		}
	}
	if len(list) == 0 {
		return nil, endedLine
	}
	return &ast.CommentGroup{List: list}, endedLine
}

// peekAt returns the nth non-comment token after the current one, without
// consuming anything. peekAt(0) is the immediately following token.
//
// Lookahead is used only where the grammar writes an explicit restriction —
// the launch-prefix `[lookahead != "."]`, a named parameter's `identifier ":"`
// — and for the one two-token check that separates a block brace from a
// literal brace in a control-flow header. It is not a general convenience.
func (p *parser) peekAt(n int) token.Token {
	seen := 0
	for i := 0; ; i++ {
		if i == len(p.ahead) {
			p.ahead = append(p.ahead, p.scan.Scan())
		}
		t := p.ahead[i]
		if t.Kind == token.COMMENT {
			continue
		}
		if seen == n || t.Kind == token.EOF {
			return t
		}
		seen++
	}
}

func (p *parser) peek() token.Token { return p.peekAt(0) }

func (p *parser) at(k token.Kind) bool { return p.tok.Kind == k }

func (p *parser) got(k token.Kind) bool {
	if p.tok.Kind == k {
		p.advanceToken()
		return true
	}
	return false
}

// atCtx reports whether the current token is the given contextual keyword,
// which scans as an ordinary identifier.
func (p *parser) atCtx(name string) bool { return p.tok.IsCtx(name) }

// ------------------------------------------------------------ diagnostics

func (p *parser) describe(t token.Token) string {
	switch t.Kind {
	case token.EOF:
		return "end of file"
	case token.IDENT:
		return strconv.Quote(t.Lit)
	case token.INT, token.FLOAT, token.CHAR, token.STRING:
		return t.Kind.String() + " " + strconv.Quote(t.Lit)
	}
	return "'" + t.Kind.String() + "'"
}

func (p *parser) report(d *diag.Diagnostic) {
	p.errCount++
	if p.rep != nil {
		p.rep.Report(d)
	}
}

func (p *parser) errorAt(code diag.Code, pos token.Pos, args ...any) {
	p.report(diag.At(code, pos, args...))
}

func (p *parser) errorSpan(code diag.Code, pos, end token.Pos, args ...any) {
	p.report(diag.New(code, pos, end, args...))
}

func (p *parser) errorHere(code diag.Code, args ...any) {
	p.report(diag.AtToken(code, p.tok, args...))
}

func (p *parser) expect(k token.Kind) token.Pos {
	pos := p.tok.Pos
	if p.tok.Kind != k {
		p.errorHere(diag.ExpectedToken, "'"+k.String()+"'", p.describe(p.tok))
		return pos
	}
	p.advanceToken()
	return pos
}

func (p *parser) expectIdent() *ast.Ident {
	if p.tok.Kind != token.IDENT {
		p.errorHere(diag.ExpectedIdent, p.describe(p.tok))
		return &ast.Ident{NamePos: p.tok.Pos, Name: "_"}
	}
	x := &ast.Ident{NamePos: p.tok.Pos, Name: p.tok.Lit}
	p.advanceToken()
	return x
}

// ------------------------------------------------- statement termination

// continues reports whether the current token may extend the construct in
// progress. Inside a bracketed group a line terminator is ordinary white
// space; at depth zero it ends the statement.
func (p *parser) continues() bool { return p.depth > 0 || !p.tok.NLBefore }

// expectTerminator consumes nothing — a terminator is a property of the next
// token, not a token of its own. A terminator may be omitted immediately
// before `}`, immediately before `case` or `default`, and at end of file.
func (p *parser) expectTerminator() {
	switch p.tok.Kind {
	case token.RBRACE, token.EOF, token.CASE, token.DEFAULT:
		return
	}
	if p.tok.NLBefore {
		return
	}
	p.errorHere(diag.ExpectedTerminator, p.describe(p.tok))
	p.advance(stmtStart)
}

// ctx is the parse state a terminator-significant body replaces.
type ctx struct {
	depth int
	noLit bool
}

// enterTerminated opens a body whose production writes `terminator`: depth
// resets so line breaks inside it end statements, and any suppressed literal
// becomes legal again because the header that suppressed it is over.
func (p *parser) enterTerminated() ctx {
	saved := ctx{p.depth, p.noLit}
	p.depth, p.noLit = 0, false
	return saved
}

func (p *parser) leave(c ctx) { p.depth, p.noLit = c.depth, c.noLit }

// stalled reports whether no token has been consumed since at, advancing one
// token when so. Every unbounded body loop pairs it with a saved position, so
// that an error path which reports without consuming cannot spin.
func (p *parser) stalled(at token.Pos) bool {
	if p.tok.Pos != at {
		return false
	}
	p.advanceToken()
	return true
}

// --------------------------------------------------------------- recovery

var stmtStart = map[token.Kind]bool{
	token.LET: true, token.VAR: true, token.IF: true, token.WHILE: true,
	token.FOR: true, token.SWITCH: true, token.SELECT: true, token.RETURN: true,
	token.DEFER: true, token.BREAK: true, token.CONTINUE: true,
	token.FALLTHROUGH: true, token.RBRACE: true,
}

var declStart = map[token.Kind]bool{
	token.FUNC: true, token.STRUCT: true, token.CLASS: true, token.ENUM: true,
	token.TYPE: true, token.CONSTRAINT: true, token.DECLARE: true,
	token.IMPORT: true, token.LET: true, token.VAR: true, token.EOF: true,
}

var memberStart = map[token.Kind]bool{
	token.FUNC: true, token.CLASS: true, token.DECLARE: true,
	token.RBRACE: true, token.EOF: true,
}

var clauseStart = map[token.Kind]bool{
	token.CASE: true, token.DEFAULT: true, token.RBRACE: true, token.EOF: true,
}

// advance skips tokens until one in `to` is reached, so a single malformed
// construct does not cascade. The syncPos/syncCnt guard is the standard
// protection against a recovery loop that never consumes.
func (p *parser) advance(to map[token.Kind]bool) {
	for ; p.tok.Kind != token.EOF; p.advanceToken() {
		if to[p.tok.Kind] {
			if p.tok.Pos == p.syncPos && p.syncCnt < 10 {
				p.syncCnt++
				return
			}
			if p.tok.Pos > p.syncPos {
				p.syncPos, p.syncCnt = p.tok.Pos, 0
				return
			}
		}
	}
}

// ------------------------------------------------------------ entry points

// ParseFile parses a single Vertex source file.
//
// src may be nil, in which case filename is read from disk. The returned File
// is non-nil even when errors were reported: recovery produces Bad* nodes so
// later phases and editor tooling can still walk a partial tree.
func ParseFile(fset *token.FileSet, filename string, src []byte, rep diag.Reporter, mode Mode) (*ast.File, error) {
	text, err := readSource(filename, src)
	if err != nil {
		return nil, err
	}

	var list diag.List
	report := rep
	if report == nil {
		report = &list
	}

	file := fset.AddFile(filename, len(text))

	var p parser
	p.init(file, text, report, mode)
	f := p.parseFile()

	if rep == nil && list.HasErrors() {
		list.Sort()
		list.Dedup()
		return f, list.Err()
	}
	if p.errCount > 0 || p.scan.ErrorCount > 0 {
		return f, fmt.Errorf("%s: %d errors", filename, p.errCount+p.scan.ErrorCount)
	}
	return f, nil
}

// ParseDir parses every .vs file in dir that targets `target` and groups them
// into one compilation unit.
//
// The two passes are required rather than opportunistic. A file whose build
// tag does not match is excluded from the build whole, so the filter must run
// before any file is fully parsed; and the qualifier for an import comes from
// the imported package's own package clause, so the surviving files' names must
// agree before anything can resolve.
func ParseDir(fset *token.FileSet, dir, importPath string, target token.BuildTag, rep diag.Reporter, mode Mode) (*ast.Package, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".vs") {
			continue
		}
		names = append(names, filepath.Join(dir, e.Name()))
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("%s: no .vs files", dir)
	}

	// Pass 1: package and build clauses only, against a throwaway set.
	type candidate struct {
		path string
		src  []byte
	}
	var keep []candidate
	probe := token.NewFileSet()
	for _, name := range names {
		text, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}
		head, err := ParseFile(probe, name, text, nil, PackageClauseOnly)
		if err != nil {
			return nil, err
		}
		switch t := head.BuildTag(); t {
		case token.TagNone, target:
			keep = append(keep, candidate{name, text})
		}
	}
	if len(keep) == 0 {
		return nil, fmt.Errorf("%s: no files match build tag %q", dir, target)
	}

	// Pass 2: full parse of the survivors, against the real set.
	files := make([]*ast.File, 0, len(keep))
	for _, c := range keep {
		f, err := ParseFile(fset, c.path, c.src, rep, mode)
		if f != nil {
			files = append(files, f)
		}
		if err != nil && rep == nil {
			return nil, err
		}
	}

	return ast.NewPackage(fset, importPath, dir, target, files)
}

func readSource(filename string, src []byte) ([]byte, error) {
	if src != nil {
		return src, nil
	}
	return os.ReadFile(filename)
}

// --------------------------------------------------------------- the file

// parseFile parses a SourceFile. A file may open with blank or comment-only
// lines; the package clause is the first non-comment construct and is
// mandatory, the build clause if present is the second, and all imports precede
// all declarations.
func (p *parser) parseFile() *ast.File {
	f := &ast.File{FileStart: p.file.Pos(0), FileEnd: p.file.Pos(p.file.Size())}
	f.Doc = p.leadComment

	if !p.at(token.PACKAGE) {
		p.errorHere(diag.MissingPackageClause)
		p.advance(declStart)
		f.Comments = p.comments
		return f
	}
	f.Package = p.tok.Pos
	p.advanceToken()
	f.Name = p.expectIdent()
	p.expectTerminator()

	if p.atCtx(token.CtxBuild) {
		f.Build = p.parseBuildClause()
	}

	if p.mode&PackageClauseOnly != 0 {
		f.Comments = p.comments
		return f
	}

	for p.at(token.IMPORT) {
		f.Imports = append(f.Imports, p.parseImportDecl())
	}

	if p.mode&ImportsOnly != 0 {
		f.Comments = p.comments
		return f
	}

	for !p.at(token.EOF) {
		before := p.tok.Pos
		switch {
		case p.at(token.IMPORT):
			p.errorHere(diag.ImportAfterDecl)
			f.Imports = append(f.Imports, p.parseImportDecl())
		case p.atCtx(token.CtxBuild):
			p.errorHere(diag.MisplacedBuildClause)
			p.parseBuildClause()
		default:
			f.Decls = append(f.Decls, p.parseTopLevelDecl())
			p.expectTerminator()
		}
		if p.stalled(before) {
			continue
		}
	}

	f.Comments = p.comments
	return f
}

// parseBuildClause parses `build BuildTag`. Both `build` and the tag are
// contextual keywords scanning as identifiers.
func (p *parser) parseBuildClause() *ast.BuildClause {
	b := &ast.BuildClause{Build: p.tok.Pos}
	p.advanceToken()

	b.TagPos = p.tok.Pos
	if p.tok.Kind != token.IDENT {
		p.errorHere(diag.ExpectedIdent, p.describe(p.tok))
		p.advance(declStart)
		return b
	}
	b.Name = p.tok.Lit
	p.advanceToken()

	tag, ok := token.LookupBuildTag(b.Name)
	if !ok {
		// An unrecognized tag is a compile error, never a silently excluded
		// file, so the clause is kept with TagNone and diagnosed.
		p.errorAt(diag.UnknownBuildTag, b.TagPos, b.Name)
	}
	b.Tag = tag
	p.expectTerminator()
	return b
}