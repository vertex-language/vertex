package ast

import (
	"fmt"
	"sort"

	"github.com/vertex-language/vertex/token"
)

// BuildClause is `build <tag>`. Both `build` and the tag scan as identifiers.
//
// Tag is TagNone when Name is unrecognized. An unrecognized tag is a compile
// error rather than a silently excluded file, so a caller must distinguish
// "unknown tag" from "no build clause" — the nil-ness of File.Build answers the
// latter.
type BuildClause struct {
	Build  token.Pos
	TagPos token.Pos
	Name   string
	Tag    token.BuildTag
}

func (b *BuildClause) Pos() token.Pos { return b.Build }
func (b *BuildClause) End() token.Pos { return b.TagPos + token.Pos(len(b.Name)) }

// File is one parsed source file.
//
// The package clause is the first non-comment construct and is mandatory; the
// build clause, if present, is the second; all imports precede all
// declarations. Top-level declarations are order-independent.
//
// Comments holds every comment in the file in source order, including those
// already attached as a Doc or trailing Comment elsewhere. Keeping the flat
// list alongside the attachments is what lets a printer place anything the
// attachment heuristic missed.
type File struct {
	Doc     *CommentGroup
	Package token.Pos
	Name    *Ident // the package clause name — the qualifier importers use

	Build   *BuildClause // nil if absent
	Imports []*ImportDecl
	Decls   []Decl

	Comments []*CommentGroup

	FileStart token.Pos
	FileEnd   token.Pos
}

func (f *File) Pos() token.Pos { return f.Package }
func (f *File) End() token.Pos { return f.FileEnd }

// BuildTag reports the file's target tag, or TagNone if it carries no clause.
func (f *File) BuildTag() token.BuildTag {
	if f.Build == nil {
		return token.TagNone
	}
	return f.Build.Tag
}

// Filename returns the file's name as recorded in the FileSet.
func (f *File) Filename(fset *token.FileSet) string {
	if tf := fset.File(f.Package); tf != nil {
		return tf.Name()
	}
	return ""
}

// Package is a set of files compiled as one unit.
//
// It is a validated container and nothing more: no I/O, no import resolution,
// no scopes. Resolution belongs to the loader, which is why this takes no
// importer.
type Package struct {
	Name  string // from the package clause; all files agree
	Path  string // resolved import path — a locator, not a name
	Dir   string
	Tag   token.BuildTag
	Files []*File // sorted by filename
}

func (p *Package) Pos() token.Pos {
	if len(p.Files) == 0 {
		return token.NoPos
	}
	return p.Files[0].Pos()
}

func (p *Package) End() token.Pos {
	if len(p.Files) == 0 {
		return token.NoPos
	}
	return p.Files[len(p.Files)-1].End()
}

// NewPackage groups files from one directory into a compilation unit.
//
// It is pure: no filesystem access, no imports, no diagnostics beyond its own
// well-formedness checks. It asserts only what makes the container coherent —
// at least one file, agreement on the package clause name, and agreement with
// the target tag. Everything else is a static rule over an already-parsed tree
// and belongs to the analyzer.
//
// Files are sorted by filename so the resulting object is byte-reproducible.
func NewPackage(fset *token.FileSet, path, dir string, target token.BuildTag, files []*File) (*Package, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("package %s: no files", path)
	}

	sorted := make([]*File, len(files))
	copy(sorted, files)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Filename(fset) < sorted[j].Filename(fset)
	})

	name := ""
	var nameFile string
	for _, f := range sorted {
		if f.Name == nil {
			return nil, fmt.Errorf("package %s: %s has no package clause",
				path, f.Filename(fset))
		}
		if name == "" {
			name, nameFile = f.Name.Name, f.Filename(fset)
			continue
		}
		if f.Name.Name != name {
			return nil, fmt.Errorf(
				"package %s: %s declares package %q but %s declares %q",
				path, f.Filename(fset), f.Name.Name, nameFile, name)
		}
	}

	// A file whose tag does not match the target is excluded from the build
	// outright, and filtering is the loader's job — so a mismatch reaching here
	// is a caller bug rather than a user error.
	for _, f := range sorted {
		if t := f.BuildTag(); t != token.TagNone && t != target {
			return nil, fmt.Errorf("package %s: %s carries build tag %q, not %q",
				path, f.Filename(fset), t, target)
		}
	}

	return &Package{Name: name, Path: path, Dir: dir, Tag: target, Files: sorted}, nil
}