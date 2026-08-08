package token

import (
	"sort"
	"sync"
)

// Pos is a source position: a 1-based byte offset into one translation unit,
// in a per-unit address space (§3). There is deliberately no FileSet — a
// global position space across files is a cross-file dependency at the one
// layer that must not have any. Positions travel with their file (§8.2).
//
// The offset is biased by one so that NoPos == 0 is distinguishable from the
// first byte of a file. Comparison, ordering, and span subtraction are
// unaffected; only File.Slice, File.Between, and File.Position unbias, and
// each does so at its own boundary.
type Pos uint32

// NoPos means "no position."
const NoPos Pos = 0

// IsValid reports whether p denotes a location.
func (p Pos) IsValid() bool { return p != NoPos }

// MaxFileSize is the largest source a single File can hold, bounded by Pos
// being uint32 with a one-byte bias.
const MaxFileSize = 1<<32 - 2

// Position is a human-readable location, for rendering only. Line and Col are
// 1-based; Col counts bytes, not runes or UTF-16 code units. A renderer that
// needs either of those converts using the source, which it has.
type Position struct {
	Line int
	Col  int
}

func (p Position) IsValid() bool { return p.Line > 0 }

func (p Position) String() string {
	if !p.IsValid() {
		return "-"
	}
	return itoa(p.Line) + ":" + itoa(p.Col)
}

// File is one translation unit: a name, the bytes, and the arithmetic to turn
// a Pos back into text or a line and column.
//
// A File is safe for concurrent use once constructed. §8.7 parses files in
// parallel, one File per goroutine, and renders diagnostics afterward from
// any goroutine.
type File struct {
	name string
	src  []byte

	linesOnce sync.Once
	lines     []Pos // lines[i] is the Pos of the first byte of line i+1

	// Line directives are not implemented. When they are, the remapping table
	// goes here and is consulted only by Position. Pos stays a raw offset, and
	// Slice, Between, and all span arithmetic stay unaffected (§3). A design
	// where a //line-equivalent changes what Pos *is* has to be rejected,
	// because every later phase is written against these semantics.
	//
	// lineDirectives *lineTable
}

// NewFile returns a File over src. src is retained, not copied; the caller
// must not mutate it afterward.
func NewFile(name string, src []byte) *File {
	if len(src) > MaxFileSize {
		panic("token: source exceeds MaxFileSize")
	}
	return &File{name: name, src: src}
}

func (f *File) Name() string { return f.name }

// Size is the length of the source in bytes.
func (f *File) Size() int { return len(f.src) }

// PosAt converts a byte offset into a Pos. offset == Size() is valid and is
// where EOF sits. Panics on an out-of-range offset, which is always a bug in
// the caller rather than bad input.
//
// This is not listed in §8.8's surface, but the scanner lives in another
// package and has to construct positions somehow.
func (f *File) PosAt(offset int) Pos {
	if offset < 0 || offset > len(f.src) {
		panic("token: offset out of range")
	}
	return Pos(offset) + 1
}

// Offset converts a Pos back to a byte offset, or -1 for NoPos.
func (f *File) Offset(p Pos) int {
	if p == NoPos {
		return -1
	}
	off := int(p) - 1
	if off > len(f.src) {
		return len(f.src)
	}
	return off
}

// Slice returns the raw bytes of tok. Literals keep their source spelling:
// `1_024` yields the five bytes as written, with no value, no width, and no
// separator stripping. Decoding belongs to a phase that knows the target
// type (§4.6, §8.3).
func (f *File) Slice(tok Token) []byte { return f.Between(tok.Pos, tok.End) }

// Between returns the bytes in [lo, hi). It takes two Pos rather than an
// ast.Node on purpose: token must not import ast (§8.3).
//
// Returns nil if either bound is NoPos or the range is inverted.
func (f *File) Between(lo, hi Pos) []byte {
	if lo == NoPos || hi == NoPos || hi < lo {
		return nil
	}
	l, h := f.Offset(lo), f.Offset(hi)
	if l < 0 || l > len(f.src) {
		return nil
	}
	if h > len(f.src) {
		h = len(f.src)
	}
	return f.src[l:h]
}

// Position converts p to a line and column, for rendering. Returns the zero
// Position for NoPos.
//
// The line index is built on first use rather than in NewFile: most files are
// parsed without any diagnostic ever being rendered, and the scanner already
// walks the source once.
func (f *File) Position(p Pos) Position {
	if p == NoPos {
		return Position{}
	}
	f.linesOnce.Do(f.buildLines)

	off := Pos(f.Offset(p))
	// The first line start greater than off belongs to the next line.
	i := sort.Search(len(f.lines), func(i int) bool { return f.lines[i] > off })
	if i == 0 {
		return Position{Line: 1, Col: int(off) + 1}
	}
	return Position{Line: i, Col: int(off-f.lines[i-1]) + 1}
}

// IsLineTerminator reports whether r ends a line, per A.1's LineTerminator:
// LF, CR, LS (U+2028), and PS (U+2029). Exported so the scanner and the line
// index cannot disagree about what a line is.
func IsLineTerminator(r rune) bool {
	switch r {
	case '\n', '\r', 0x2028, 0x2029:
		return true
	}
	return false
}

func (f *File) buildLines() {
	// lines[0] is the offset of line 1, always 0.
	lines := make([]Pos, 1, len(f.src)/24+2)
	for i := 0; i < len(f.src); {
		switch c := f.src[i]; {
		case c == '\n':
			i++
			lines = append(lines, Pos(i))
		case c == '\r':
			i++
			if i < len(f.src) && f.src[i] == '\n' {
				i++ // CRLF is one terminator
			}
			lines = append(lines, Pos(i))
		case c == 0xE2 && i+2 < len(f.src) && f.src[i+1] == 0x80 &&
			(f.src[i+2] == 0xA8 || f.src[i+2] == 0xA9): // LS, PS
			i += 3
			lines = append(lines, Pos(i))
		default:
			i++
		}
	}
	f.lines = lines
}