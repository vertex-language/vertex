package token

import (
	"fmt"
	"sort"
	"sync"
)

// Pos is a compact position: an offset into a FileSet's global address space.
// The zero value NoPos means "no position".
type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool { return p != NoPos }

// Position is a resolved, human-facing position.
type Position struct {
	Filename string
	Offset   int // byte offset, 0-based
	Line     int // 1-based
	Column   int // 1-based, in bytes
}

func (p Position) IsValid() bool { return p.Line > 0 }

func (p Position) String() string {
	if !p.IsValid() {
		return "-"
	}
	if p.Filename == "" {
		return fmt.Sprintf("%d:%d", p.Line, p.Column)
	}
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line, p.Column)
}

// File is one source file's slice of a FileSet's address space.
type File struct {
	name  string
	base  int
	size  int
	mu    sync.Mutex
	lines []int // byte offset of the first character of each line; lines[0] == 0
}

func (f *File) Name() string { return f.name }
func (f *File) Base() int    { return f.base }
func (f *File) Size() int    { return f.size }

func (f *File) LineCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.lines)
}

// AddLine records that a line begins at offset. The scanner calls this each
// time it consumes a LineTerminator (A.1: <LF>, <CR>, <CR><LF>, <LS>, <PS>).
// Offsets must be monotonically increasing.
func (f *File) AddLine(offset int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if i := len(f.lines); (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
}

// LineStart returns the Pos of the first character of line (1-based).
// vir needs this to emit debug line tables.
func (f *File) LineStart(line int) Pos {
	f.mu.Lock()
	defer f.mu.Unlock()
	if line < 1 || line > len(f.lines) {
		return NoPos
	}
	return Pos(f.base + f.lines[line-1])
}

func (f *File) Pos(offset int) Pos {
	if offset < 0 || offset > f.size {
		panic("token: offset out of bounds")
	}
	return Pos(f.base + offset)
}

func (f *File) Offset(p Pos) int {
	if int(p) < f.base || int(p) > f.base+f.size {
		panic("token: Pos out of bounds for file")
	}
	return int(p) - f.base
}

func (f *File) Position(p Pos) Position {
	if !p.IsValid() {
		return Position{}
	}
	offset := f.Offset(p)
	f.mu.Lock()
	i := sort.SearchInts(f.lines, offset+1) - 1
	lineStart := 0
	if i >= 0 {
		lineStart = f.lines[i]
	}
	f.mu.Unlock()
	return Position{
		Filename: f.name,
		Offset:   offset,
		Line:     i + 1,
		Column:   offset - lineStart + 1,
	}
}

// FileSet is a shared position space. One per compilation; a package's files
// all live in it so a diagnostic can span two files in the same directory.
type FileSet struct {
	mu    sync.RWMutex
	base  int
	files []*File
}

func NewFileSet() *FileSet {
	return &FileSet{base: 1} // 0 is NoPos
}

func (s *FileSet) AddFile(name string, size int) *File {
	s.mu.Lock()
	defer s.mu.Unlock()
	f := &File{name: name, base: s.base, size: size, lines: []int{0}}
	s.base += size + 1
	s.files = append(s.files, f)
	return f
}

func (s *FileSet) File(p Pos) *File {
	if !p.IsValid() {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	i := sort.Search(len(s.files), func(i int) bool {
		return s.files[i].base > int(p)
	}) - 1
	if i < 0 {
		return nil
	}
	f := s.files[i]
	if int(p) > f.base+f.size {
		return nil
	}
	return f
}

func (s *FileSet) Position(p Pos) Position {
	if f := s.File(p); f != nil {
		return f.Position(p)
	}
	return Position{}
}