// Package vir is the mechanical translator: a *hir.Program in, one
// *ir.Module per originating Vertex package out, emitted through vvm's
// append-only builder API.
//
// Every decision was already made in lower/hir. If this package ever needs
// a type switch on "is this an owning type" or "was this transferred," that
// logic belongs upstream (overview §4, invariant 2). What is left here is
// three kinds of work, and nothing else:
//
//   - Translation. hir's Op vocabulary is a subset of vir's §4 opcodes,
//     spelled identically, so instruction selection is the table in
//     instr.go rather than a decision.
//   - Ordering. vir fixes a mandatory section order (§2.1) and forbids
//     forward references (§2.2). hir's slices are in construction order,
//     which is neither. decl.go and func.go sort them.
//   - Naming. A vir module's `namespace` is the import path and its
//     `module` is the PackageClause name — A.2.3's "the path is a locator,
//     the declared name is the qualifier," carried through unchanged.
//
// Output is unverified by construction: the builder API validates nothing,
// and ir/verify.Verify is vvm's to run. This package's own errors are
// therefore about shapes it could not *spell*, never about shapes it
// judged wrong.
//
// # The one target-shaped fact
//
// hir never sees a triple (invariant 1) and neither does most of this
// package. The exception is a module lowered from a file carrying a
// `declare` block: §2.1 makes `target` mandatory whenever `link` is
// present, so Config.Target is emitted into exactly those modules and no
// others. Invariant 3 is therefore checkable by reading the emitted text —
// a module with no links has no target line.
package vir

import (
	"errors"
	"fmt"

	"github.com/vertex-language/vertex/lower/hir"
	"github.com/vertex-language/vertex/token"
	ir "github.com/vertex-language/vvm/ir/vir"
)

// Config is everything emission needs that isn't the program itself.
type Config struct {
	// Fset resolves token.Pos for `loc` lines. Nil disables them.
	Fset *token.FileSet

	// Target is written into modules carrying `link` declarations, and
	// only those. Zero is legal for a program with no declare block.
	Target ir.Target

	// Debug emits `loc` body-lines, so a .vir diagnostic points back at
	// Vertex source. One per source line, not one per instruction.
	Debug bool
}

// Lower emits one *ir.Module per hir Module, preserving hir's topological
// order — a module appears after everything it imports, which is what
// vvm's importer expects and what declare-before-use requires across the
// module graph as well as within one.
//
// Errors accumulate rather than aborting, matching the toolchain's habit
// everywhere else. An `internal:` error means hir handed over something
// that should not exist; a `todo:` error means a valid hir shape this
// package cannot spell in vir yet.
func Lower(conf *Config, prog *hir.Program) ([]*ir.Module, error) {
	if conf == nil {
		return nil, errors.New("lower/vir: nil Config")
	}
	if prog == nil {
		return nil, errors.New("lower/vir: nil Program")
	}
	l := &lowerer{conf: conf, prog: prog}
	out := make([]*ir.Module, 0, len(prog.Modules))
	for _, hm := range prog.Modules {
		out = append(out, l.module(hm))
	}
	return out, l.err()
}

// Root names the module holding the program's entry point, for vvm's
// --root. It is the driver's other half of Lower's return value.
func Root(prog *hir.Program) string {
	if prog == nil || prog.Entry == nil || prog.Entry.Module == nil {
		return ""
	}
	return prog.Entry.Module.Name
}

type lowerer struct {
	conf *Config
	prog *hir.Program
	errs []error
}

// module emits one module in §2.1's mandatory section order. The order of
// these calls is the whole guarantee this function makes; everything they
// call is free to be as mechanical as it likes.
func (l *lowerer) module(hm *hir.Module) *ir.Module {
	m := ir.NewModule(hm.Name)
	if hm.Path != "" {
		m.SetNamespace(hm.Path)
	}
	// §2.1 step 3: optional, but required if `link` is present. Invariant
	// 3 says only a declare-block module has links.
	if len(hm.Links) > 0 {
		t := l.conf.Target
		m.SetTarget(t.Arch, t.OS, t.ABI, t.Tiers...)
	}

	l.structs(m, hm)
	// fnsig: nothing emits one yet — see the README's "Not emitted".
	l.consts(m, hm)
	l.globals(m, hm)
	l.links(m, hm)
	l.imports(m, hm)
	l.funcs(m, hm)
	return m
}

func (l *lowerer) errorf(pos token.Pos, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if pos.IsValid() && l.conf.Fset != nil {
		msg = l.conf.Fset.Position(pos).String() + ": " + msg
	}
	l.errs = append(l.errs, errors.New("lower/vir: "+msg))
}

func (l *lowerer) todo(pos token.Pos, format string, args ...any) {
	l.errorf(pos, "todo: "+format, args...)
}

func (l *lowerer) err() error {
	if len(l.errs) == 0 {
		return nil
	}
	return errors.Join(l.errs...)
}