package ast

import "github.com/vertex-language/vertex/token"

// File is one translation unit — SourceFile in §1's single goal symbol.
//
// There is no goal-symbol field and no mode: every file is parsed the same
// way, `[+Await]` is unconditional at file scope, and nothing about the tree
// records which of Script or Module it "is" (§1).
//
// The file's name and source are not here. Positions are byte offsets in a
// per-unit address space and mean nothing without their token.File, so the two
// travel together as a pair (§8.2). Anything that outlives a single parse
// keys on that pair, not on this node.
type File struct {
	Items    []Stmt          // ModuleItemList, source order
	Comments []*CommentGroup // nil unless parser.ParseComments
	FileEnd  token.Pos       // end of input, so an empty file still has a span

	// Arena owns the nodes in this tree. It is set by whoever allocated them
	// and may be nil for a hand-built tree.
	//
	// This field is how ast avoids importing an arena package: §2's diagram
	// has four packages and the allocator is a fifth, below everything. An
	// interface keeps the edge from existing at all.
	Arena Releaser
}

// Releaser frees a tree's backing storage.
type Releaser interface{ Release() }

func (f *File) Pos() token.Pos {
	if len(f.Items) > 0 {
		return f.Items[0].Pos()
	}
	return token.Pos(1)
}

func (f *File) End() token.Pos {
	if len(f.Items) > 0 {
		return f.Items[len(f.Items)-1].End()
	}
	return f.FileEnd
}

// Release frees the arena. Nodes reachable from this tree are invalid
// afterward, including any Pos values a caller copied out of them — those are
// plain integers and stay readable, which is exactly the trap: read the text
// you need before releasing (§8.1, §8.7).
//
// Idempotent, so `defer tree.Release()` alongside an explicit release is safe.
// Arenas are per-file and not shared, so this can happen as soon as one file's
// tree is consumed; a whole-program build does not have to hold every tree at
// once.
func (f *File) Release() {
	if f.Arena != nil {
		f.Arena.Release()
		f.Arena = nil
	}
}