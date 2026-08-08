package parser

import (
	"github.com/vertex-language/vertex/ast"
)

// The file is the arena boundary: nodes for one file allocate together and are
// freed together (§1). Arenas are per-file and not shared, so Release is
// per-file and can happen as soon as that file's tree is consumed (§8.7).
//
// This lives inside parser rather than in its own package on purpose. §2's
// diagram has exactly four packages, and an allocator package would be a fifth
// sitting below all of them. ast depends on it only through the one-method
// ast.Releaser interface, so no import edge exists.
//
// The implementation is per-type slabs rather than a byte arena, because a byte
// arena in Go needs unsafe to hand out typed pointers and the GC has to be
// told about them anyway. Slabs give the two properties that actually matter:
// bulk free, and O(1) rollback by truncation (§6.1).
type arena struct {
	idents []ast.Ident
	slabs  []func()  // reset closures, one per slab kind
	freed  bool
	detached bool
}

// arenaMark records slab lengths. Rollback truncates them; pointers into the
// truncated region belong to nodes the parser just discarded, so nothing leaks
// and nothing dangles.
type arenaMark struct {
	idents int
}

const identSlab = 256

func newArena() *arena {
	return &arena{idents: make([]ast.Ident, 0, identSlab)}
}

func (a *arena) mark() arenaMark { return arenaMark{idents: len(a.idents)} }

func (a *arena) reset(m arenaMark) {
	if m.idents <= len(a.idents) {
		a.idents = a.idents[:m.idents]
	}
}

// ident allocates from the slab while there is room, and falls back to the
// heap once the slab is full.
//
// Growing the slab with append would be wrong: append reallocates, and every
// pointer already handed out would point into the old backing array. Nodes
// would silently stop being part of the tree the caller walks. So the slab is
// fixed-capacity and overflow goes to the GC.
func (a *arena) ident(pos, end interface{ }) *ast.Ident { panic("unused") }

func (a *arena) newIdent() *ast.Ident {
	if len(a.idents) == cap(a.idents) {
		return new(ast.Ident)
	}
	a.idents = a.idents[:len(a.idents)+1]
	id := &a.idents[len(a.idents)-1]
	*id = ast.Ident{}
	return id
}

// Release frees the arena. Idempotent, so `defer tree.Release()` alongside an
// explicit release is safe (§8.1).
func (a *arena) Release() {
	if a.freed || a.detached {
		return
	}
	a.freed = true
	a.idents = nil
}

// detach gives up bulk freeing for the fragment entry points (§8.6), which
// return a bare node with no ast.File for the caller to Release.
func (a *arena) detach() { a.detached = true }

var _ ast.Releaser = (*arena)(nil)