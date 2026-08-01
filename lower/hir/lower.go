package hir

import (
	"errors"
	"fmt"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// Mode selects what the build exists for.
type Mode uint8

const (
	// ModeProgram seeds monomorphization from main.
	ModeProgram Mode = iota
	// ModeTest compiles one test function in isolation into its own
	// binary (overview §6.2). Isolation starts at analyzer, not here; by
	// the time hir runs, the unit already contains exactly one test
	// function, and hir's only extra job is the render wrapper.
	ModeTest
)

// Config is everything lowering needs that isn't the checked graph itself.
// Note what is absent: no target triple, no arch, no OS. The one
// target-shaped fact hir consumes is layout, and it arrives as Sizes.
type Config struct {
	Fset  *token.FileSet
	Sizes types.Sizes

	// PointerSize is the width Sizes assigns a pointer, in bytes. Zero
	// means 8. It is separate from Sizes only because the synthesized
	// header layouts (slice, string, closure) are hir's own shapes and
	// have no types.Type for Sizes to measure.
	PointerSize int64

	Mode Mode

	// TestFunc names the single test function under ModeTest.
	TestFunc string

	// MaxDepth bounds recursive instantiation. A.7.6 makes non-terminating
	// instantiation a compile error; the worklist carries its own guard
	// rather than assuming the analyzer enforced it upstream. Zero means 64.
	MaxDepth int
}

// Unit is one checked Vertex package, as importer.Load produced it. hir
// takes this shape rather than an *importer.Result so it does not depend on
// the loader — and so a test can hand it a package built by hand.
type Unit struct {
	Path  string
	Name  string
	Files []*ast.File
	Types *types.Package
	Info  *types.Info
}

// Lower runs the whole front-to-hir pipeline over a checked graph and
// returns one *hir.Program. units must be in topological order (a package
// after everything it imports) — importer.Result.Order already is.
//
// Errors are accumulated rather than fatal, mirroring the toolchain's habit
// everywhere else: a construct hir cannot lower yet produces a `todo:`
// error naming it, distinct from an `internal:` error, which means
// something the analyzer should already have caught got through.
func Lower(conf *Config, units []*Unit) (*Program, error) {
	if conf == nil {
		return nil, errors.New("hir: nil Config")
	}
	l := &lowerer{
		conf:    conf,
		units:   units,
		modules: map[string]*Module{},
		byUnit:  map[*Unit]*Module{},
		prog:    &Program{},
	}
	l.types = newTypeLowerer(l)
	l.work = newWorklist(l)

	for _, u := range units {
		m := newModule(u.Path, u.Name)
		l.modules[u.Path] = m
		l.byUnit[u] = m
		l.prog.Modules = append(l.prog.Modules, m)
	}

	// 0. Declarations first: structs, consts, globals, and every declare
	//    block. Bodies reference them, and vir's fixed section order plus
	//    declare-before-use means they must exist before any call site can
	//    name one.
	for _, u := range units {
		l.declarations(u)
	}

	// 1. Monomorphization, seeded from the roots. Everything below is
	//    type-dependent, so nothing can run before concrete types exist.
	root := l.seed()
	if root == nil {
		return l.prog, l.err()
	}
	l.work.run()

	// 2. async/await state-machine split. Not implemented; see async.go.
	//    When it lands it must run here, before epilogue expansion.
	l.splitAsync()

	// 3+4. defer/deinit epilogue expansion and ownership expansion are
	//    performed during body lowering, through the builder's scope stack
	//    (see builder.go's scope type and epilogue.go). This is a
	//    deliberate deviation from the overview's five-pass shape, recorded
	//    in the package README: it is equivalent for a program with no
	//    async functions, and the seam is kept narrow so it can be lifted
	//    into standalone passes when the split in step 2 arrives.

	// 5. Control-flow flattening: structured statements become the Join
	//    Convention shape, which is what makes lower/vir mechanical.
	for _, m := range l.prog.Modules {
		for _, f := range m.Funcs {
			Flatten(f)
		}
	}

	l.prog.Features = l.feats
	return l.prog, l.err()
}

// lowerer is the state of one Lower call.
type lowerer struct {
	conf  *Config
	units []*Unit

	prog    *Program
	modules map[string]*Module
	byUnit  map[*Unit]*Module

	types *typeLowerer
	work  *worklist
	own   *ownership

	feats builtins.FeatureSet

	// cur is the instance being lowered: which unit's Info to read, which
	// module to emit into, and which type-parameter substitution is live.
	cur *instance

	errs []error
}

// instance is one monomorphized function in flight.
type instance struct {
	unit  *Unit
	mod   *Module
	subst map[*types.TypeParam]types.Type
	depth int
}

func (l *lowerer) info() *types.Info {
	if l.cur == nil {
		return nil
	}
	return l.cur.Info()
}

func (i *instance) Info() *types.Info { return i.unit.Info }

// subst applies the live monomorphization substitution to t. Outside an
// instance, or for a type mentioning no parameters, it is the identity.
func (l *lowerer) subst(t types.Type) types.Type {
	if l.cur == nil || len(l.cur.subst) == 0 || t == nil {
		return t
	}
	return types.Substitute(t, l.cur.subst)
}

func (l *lowerer) hirType(t types.Type) Type {
	m := l.currentModule()
	return l.types.lower(m, t)
}

func (l *lowerer) currentModule() *Module {
	if l.cur != nil {
		return l.cur.mod
	}
	return l.prog.Modules[0]
}

// need records a builtin feature the program uses. Every emitted builtin
// call goes through here, so the feature set can never disagree with the
// calls actually emitted.
func (l *lowerer) need(f builtins.Feature) { l.feats = l.feats.With(f) }

func (l *lowerer) needSymbol(s builtins.Symbol) {
	if f, ok := builtins.FeatureFor(s); ok {
		l.need(f)
	}
}

func (l *lowerer) errorf(pos token.Pos, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if pos.IsValid() && l.conf.Fset != nil {
		msg = l.conf.Fset.Position(pos).String() + ": " + msg
	}
	l.errs = append(l.errs, errors.New("hir: "+msg))
}

// todo marks a construct that is valid Vertex this package does not lower
// yet — the same distinction vvm's lowering backends draw between a plain
// error (a bug upstream) and a todo (an unimplemented construct).
func (l *lowerer) todo(pos token.Pos, format string, args ...any) {
	l.errorf(pos, "todo: "+format, args...)
}

func (l *lowerer) err() error {
	if len(l.errs) == 0 {
		return nil
	}
	return errors.Join(l.errs...)
}

// seed picks the roots monomorphization starts from. Exported *generic*
// functions are not roots — there are no concrete type arguments to seed
// with — so only reachability from main (or the one test function)
// instantiates anything.
func (l *lowerer) seed() *types.Func {
	want := "main"
	if l.conf.Mode == ModeTest {
		want = l.conf.TestFunc
	}
	for _, u := range l.units {
		if u.Types == nil {
			continue
		}
		obj := u.Types.Scope().Lookup(want)
		fn, ok := obj.(*types.Func)
		if !ok {
			continue
		}
		l.work.enqueue(u, fn, nil, 0)
		l.buildEntry(u, fn)
		return fn
	}
	l.errorf(0, "internal: no %q function in any root package", want)
	return nil
}