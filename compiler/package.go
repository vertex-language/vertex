package compiler

import (
	"path/filepath"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/vertex-language/vertex/parser"
)

type BuildTags struct {
	Tags []string
}

func (b *BuildTags) Has(tag string) bool {
	for _, t := range b.Tags {
		if strings.TrimSpace(t) == tag {
			return true
		}
	}
	return false
}

func DefaultBuildTags(platform string) *BuildTags {
	if platform == "" {
		platform = "linux"
	}
	return &BuildTags{Tags: []string{platform}}
}

// SourceFile is a parsed .vs file.
type SourceFile struct {
	Path         string
	PackName     string
	BuildTags    *BuildTags
	Imports      []*Import // all import paths in this file
	NativeImport *Import   // first native import found (lib/, linux/, gpu/, …)
	Tree         parser.IFileContext
	ParseErrs    []error
}

// Package is a group of SourceFiles that form one Vertex package.
type Package struct {
	Name       string
	Dir        string
	ImportPath string
	Files      []*SourceFile
	BuildTags  *BuildTags
}

// ParseFile parses a single .vs source file using the ANTLR parser.
func ParseFile(path string) (*SourceFile, error) {
	input, err := antlr.NewFileStream(path)
	if err != nil {
		return nil, err
	}

	lexer := parser.NewVertexLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewVertexParser(stream)

	errL := &syntaxErrListener{file: path}
	p.RemoveErrorListeners()
	p.AddErrorListener(errL)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errL)

	tree := p.File()

	sf := &SourceFile{
		Path:      path,
		Tree:      tree,
		BuildTags: &BuildTags{},
		ParseErrs: errL.errs,
	}

	if pd := tree.PackageDecl(); pd != nil && pd.ID() != nil {
		sf.PackName = pd.ID().GetText()
	}

	if bd := tree.BuildDecl(); bd != nil {
		for _, tag := range bd.AllBuildTag() {
			sf.BuildTags.Tags = append(sf.BuildTags.Tags, tag.GetText())
		}
	}

	for _, id := range tree.AllImportDecl() {
		for _, lit := range id.AllSTRING_LIT() {
			raw := strings.Trim(lit.GetText(), `"`)
			imp := ParseImportPath(raw)
			sf.Imports = append(sf.Imports, imp)
			if imp.IsNative() && sf.NativeImport == nil {
				sf.NativeImport = imp
			}
		}
	}

	return sf, nil
}

// PlatformMatch reports whether filename should be compiled for these tags.
func PlatformMatch(filename string, tags *BuildTags) bool {
	base := strings.TrimSuffix(filepath.Base(filename), ".vs")
	for _, p := range []string{"linux", "darwin", "windows"} {
		if strings.HasSuffix(base, "_"+p) {
			return tags.Has(p)
		}
	}
	return true
}

type syntaxErrListener struct {
	*antlr.DefaultErrorListener
	file string
	errs []error
}

func (l *syntaxErrListener) SyntaxError(
	_ antlr.Recognizer, _ interface{},
	line, col int, msg string, _ antlr.RecognitionException,
) {
	l.errs = append(l.errs, &CompileError{
		Pos:     Pos{File: l.file, Line: line, Col: col},
		Message: msg,
	})
}