// Package parser turns Vertex source into an ast.File or an ast.Package.
//
// Of A.0.2's four context parameters the parser tracks exactly one. Await, Npu,
// and Own appear in A.14's index of forms that parse and are rejected by a later
// phase — `await` outside an async body, `tensor` outside an npu body, `var`
// outside an owning position — and all three are spelled with keywords, so
// parsing them unconditionally is unambiguous and yields a diagnostic that can
// quote the construct rather than a token position. Lit is different: A.4.7
// directs the reader to parenthesize a literal used in a header, which is a
// parse decision, so it is the one parameter carried here as state.
//
// Statement termination (A.0.6) is driven by a bracket-nesting depth rather
// than by a token. `(`, `[`, and a literal `{` push depth; a Block's `{` resets
// depth to zero, because statements inside a block do end at line terminators
// while entries inside a literal do not. At depth zero a token carrying
// NLBefore ends the statement and may not continue a postfix or binary chain.
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
	// PackageClauseOnly stops after the package and build clauses. Load-bearing
	// for the loader, not an optimization: A.2.3 makes the imported package's
	// own PackageClause the qualifier, so a file's names cannot be resolved
	// until every imported directory's package line has been read.
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

	tok  token.Token
	next token.Token
	has  bool // next is buffered

	comments    []*ast.CommentGroup
	leadComment *ast.CommentGroup
	lineComment *ast.CommentGroup

	// noLit is [~Lit]: set while parsing a control-flow header, where a bare
	// composite or map literal is ungrammatical (A.4.7).
	noLit bool

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

func (p *parser) scanRaw() token.Token { return p.scan.Scan() }

// advanceToken moves to the next non-comment token, collecting comment groups.
//
// A skipped comment's NLBefore is folded into the token that follows it.
// Without that, `x = 1` <newline> `/*c*/ y` would lose the line break entirely,
// since the flag lands on the COMMENT and the COMMENT is dropped.
func (p *parser) advanceToken() {
	p.leadComment = nil
	p.lineComment = nil

	prevEnd := p.tok.Pos
	if p.tok.Kind != token.INVALID {
		prevEnd = p.tok.End()
	}

	var t token.Token
	if p.has {
		t, p.has = p.next, false
	} else {
		t = p.scanRaw()
	}

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

	_ = prevEnd
	p.tok = t
}

// consumeCommentGroup gathers a run of comments with no intervening token,
// leaving *t on the first non-comment token. endedLine reports whether the group
// was followed by a line terminator.
func (p *parser) consumeCommentGroup(t *token.Token, leading bool) (g *ast.CommentGroup, endedLine bool) {
	var list []*ast.Comment
	for t.Kind == token.COMMENT {
		c := &ast.Comment{Slash: t.Pos, Text: t.Lit}
		list = append(list, c)
		if c.IsLineTerminator() {
			endedLine = true
		}
		if p.has {
			*t, p.has = p.next, false
		} else {
			*t = p.scanRaw()
		}
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

// peek returns the token after the current one without consuming it.
// Used only where the grammar writes an explicit [lookahead] restriction.
func (p *parser) peek() token.Token {
	if !p.has {
		t := p.scanRaw()
		for t.Kind == token.COMMENT {
			t = p.scanRaw()
		}
		p.next, p.has = t, true
	}
	return p.next
}

func (p *parser) at(k token.Kind) bool { return p.tok.Kind == k }

func (p *parser) got(k token.Kind) bool {
	if p.tok.Kind == k {
		p.advanceToken()
		return true
	}
	return false
}

// atCtx reports whether the current token is the given ContextualKeyword
// (A.1.3), which scans as IDENT.
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

func (p *parser) errorAt(code diag.Code, pos token.Pos, args ...any) {
	p.errCount++
	if p.rep != nil {
		p.rep.Report(diag.At(code, pos, args...))
	}
}

func (p *parser) errorHere(code diag.Code, args ...any) {
	p.errCount++
	if p.rep != nil {
		p.rep.Report(diag.AtToken(code, p.tok, args...))
	}
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
// progress. Inside any bracketed group a line terminator is ordinary
// whitespace; at depth zero it ends the statement (A.0.6).
func (p *parser) continues() bool { return p.depth > 0 || !p.tok.NLBefore }

func (p *parser) expectStmtEnd() {
	switch p.tok.Kind {
	case token.RBRACE, token.EOF, token.CASE, token.DEFAULT:
		return
	}
	if p.tok.NLBefore {
		return
	}
	p.errorHere(diag.ExpectedStmtEnd, p.describe(p.tok))
	p.advance(stmtStart)
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

// ParseFile parses a single .vs source file.
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
// The two passes are required rather than opportunistic. A.2.2 excludes a
// non-matching file from the build whole, so the target filter must run before
// any file is fully parsed; and A.2.3 takes the qualifier from the package
// clause, so the surviving files' names must agree before anything can resolve.
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

	// Pass 1: package and build clauses only.
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

	// Pass 2: full parse of the survivors.
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

func (p *parser) parseFile() *ast.File {
	f := &ast.File{FileStart: p.file.Pos(0), FileEnd: p.file.Pos(p.file.Size())}
	f.Doc = p.leadComment

	// A.2 ⊢ PackageClause is mandatory and first.
	if !p.at(token.PACKAGE) {
		p.errorHere(diag.MissingPackageClause)
		p.advance(declStart)
		f.Comments = p.comments
		return f
	}
	f.Package = p.tok.Pos
	p.advanceToken()
	f.Name = p.expectIdent()
	p.expectStmtEnd()

	// A.2.2: the build clause is the second construct, and `build` is a
	// ContextualKeyword valid only as the second line-initial token.
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
		switch {
		case p.at(token.IMPORT):
			p.errorHere(diag.ImportAfterDecl)
			f.Imports = append(f.Imports, p.parseImportDecl())
		case p.atCtx(token.CtxBuild):
			p.errorHere(diag.MisplacedBuildClause)
			p.parseBuildClause()
		default:
			f.Decls = append(f.Decls, p.parseTopLevelDecl())
		}
	}

	f.Comments = p.comments
	return f
}

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
		// A.2.2 ⊢ an unrecognised tag is a compile error, not a silently
		// excluded file, so the clause is kept with TagNone and diagnosed.
		p.errorAt(diag.UnknownBuildTag, b.TagPos, b.Name)
	}
	b.Tag = tag
	p.expectStmtEnd()
	return b
}