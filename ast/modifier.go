package ast

import "github.com/vertex-language/vertex/token"

// ModifierSet is a bitset over ClassElementModifier (E), for queries.
type ModifierSet uint16

const (
	ModDeclare ModifierSet = 1 << iota
	ModStatic
	ModAbstract
	ModOverride
	ModReadonly
	ModPublic
	ModProtected
	ModPrivate
)

// ModAccessibility is the public | protected | private group, which E treats
// as a unit under AccessibilityModifier.
const ModAccessibility = ModPublic | ModProtected | ModPrivate

// ModifierTok is one modifier word with its position.
type ModifierTok struct {
	Bit     ModifierSet
	Ctx     token.Ctx // the spelling, for messages
	Pos     token.Pos
	End     token.Pos
}

// Modifiers records a modifier sequence twice over (§5.2).
//
// ClassElementModifiers (E) derives *any* sequence of modifiers, including
// duplicates and orders the language rejects. Pointing at the offending word
// needs per-modifier positions; the bitset alone can't do it. So Set answers
// queries and List answers diagnostics, and List is in source order.
//
// The redundancy is deliberate and the two must agree — see checkModifiers in
// the parser, and the invariant test in this package.
type Modifiers struct {
	Set  ModifierSet
	List []ModifierTok // source order + position
}

func (m *Modifiers) Has(bit ModifierSet) bool { return m.Set&bit != 0 }

// Find returns the first occurrence of bit, for a diagnostic that has to point
// somewhere. Returns nil if absent.
func (m *Modifiers) Find(bit ModifierSet) *ModifierTok {
	for i := range m.List {
		if m.List[i].Bit == bit {
			return &m.List[i]
		}
	}
	return nil
}

// Empty reports whether any modifier is present. A Modifiers value is used
// inline rather than by pointer, so the zero value must be meaningful.
func (m *Modifiers) Empty() bool { return len(m.List) == 0 }

func (m *Modifiers) Pos() token.Pos {
	if len(m.List) == 0 {
		return token.NoPos
	}
	return m.List[0].Pos
}

func (m *Modifiers) End() token.Pos {
	if len(m.List) == 0 {
		return token.NoPos
	}
	return m.List[len(m.List)-1].End
}