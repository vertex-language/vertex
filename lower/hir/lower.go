package hir

import (
	"fmt"
	"sort"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
)

// Mode selects what the program is built around.
type Mode uint8

const (
	// ModeProgram seeds monomorphization from main. Exported *generic*
	// functions are not roots — there are no concrete type arguments to seed
	// with — so only reachability from main instantiates anything.
	ModeProgram Mode = iota

	// ModeTest seeds from the single test function named by Config.Test, and
	// synthesizes a render wrapper into the entry slot instead of a main
	// wrapper. A test binary holds exactly one test function; nothing
	// enumerates tests inside a process.
	ModeTest
)

type Config struct {
	Fset  *token.FileSet // resolves positions into loc lines; hir resolves none itself
	Sizes *types.Sizes   // the only target-shaped input
	Mode  Mode
	Test  string // the one test function, under ModeTest
}

// Unit is one checked package. hir takes []*Unit rather than an
// *importer.Result so it does not depend on the loader, and so a test can
// hand it a graph built by hand.
type Unit struct {
	Pkg   *types.Package
	Info  *types.Info
	Files []*ast.File
}

// TodoError means the program is valid and this package does not lower that
// construct yet. Distinct from InternalError, which means something the
// analyzer should already have rejected got through — a bug here or in
// analyzer, never in the user's source. A driver reports the first as
// "unsupported" and the second as a compiler crash.
type TodoError struct{ What string }

func (e *TodoError) Error() string { return "todo: " + e.What }

type InternalError struct{ What string }

func (e *InternalError) Error() string { return "internal: " + e.What }

// Lower turns the checked graph into a Program. units must be in topological
// order — dependencies before dependents.
//
// Sizes must match the build tag the packages were *checked* under: `int`
// and `uint` are the target's pointer width, so a mismatch silently changes
// every layout the checker already committed to.
func Lower(conf *Config, units []*Unit) (prog *Program, err error) {
	if conf == nil || conf.Sizes == nil {
		return nil, &InternalError{"Lower: Config.Sizes is required"}
	}

	l := &lowerer{
		conf:     conf,
		units:    units,
		prog:     &Program{},
		modOf:    map[*types.Package]*Module{},
		infoOf:   map[*types.Package]*types.Info{},
		typeName: nil, // set below; needs l
		done:     map[string]*Func{},
	}
	l.typeName = l.mangleTypeName

	// bail is how todo/bug unwind out of the deeply recursive expression
	// walk without threading an error through every return. Nothing else in
	// this package panics.
	defer func() {
		if r := recover(); r != nil {
			if b, ok := r.(bail); ok {
				prog, err = nil, b.err
				return
			}
			panic(r)
		}
	}()

	for _, u := range units {
		m := newModule(u.Pkg.Name(), u.Pkg.Path())
		l.prog.Modules = append(l.prog.Modules, m)
		l.modOf[u.Pkg] = m
		l.infoOf[u.Pkg] = u.Info
	}

	// Pass 1: declarations. Shapes, constants, globals, declare blocks, and
	// a shell for every non-generic function.
	for _, u := range units {
		l.unit = u
		l.mod = l.modOf[u.Pkg]
		l.info = u.Info
		l.collect(u)
	}

	// Pass 2: monomorphization, seeded from the roots.
	if err := l.seed(); err != nil {
		return nil, err
	}
	l.drain()

	// Pass 3 is async.go's state-machine split and is not implemented.
	//
	// Pass 4 — defer/deinit epilogue expansion — happens *during* body
	// lowering, through the builder's scope stack. That is a known deviation
	// from lowering.md, equivalent for a program with no async functions and
	// wrong for one that has them, because a suspend edge is not a scope
	// exit. Landing async.go means lifting it out first; see async.go.

	// Pass 5: flatten.
	for _, m := range l.prog.Modules {
		for _, f := range m.Funcs {
			if f.Body != nil {
				Flatten(f)
			}
		}
	}

	l.prog.Features = l.features
	return l.prog, nil
}

type bail struct{ err error }

type lowerer struct {
	conf  *Config
	units []*Unit
	prog  *Program

	// current unit, set by every pass that walks per-package
	unit *Unit
	mod  *Module
	info *types.Info

	modOf  map[*types.Package]*Module
	infoOf map[*types.Package]*types.Info

	typeName func(types.Type) string

	// work is the monomorphization worklist; done memoizes by instance key
	// so a diamond of instantiations is built once.
	work []*instance
	done map[string]*Func

	// depth guards non-terminating instantiation. semantics.md §9 makes it a
	// compile error; the guard lives here rather than assuming the analyzer
	// enforced it upstream.
	depth int

	features builtins.FeatureSet
}

const maxInstantiationDepth = 64

func (l *lowerer) todo(what string) *Type {
	panic(bail{&TodoError{what}})
}

func (l *lowerer) bug(what string) *Type {
	panic(bail{&InternalError{what}})
}

func (l *lowerer) todoAt(pos token.Pos, what string) {
	panic(bail{&TodoError{fmt.Sprintf("%s (%s)", what, l.conf.Fset.Position(pos))}})
}

func (l *lowerer) bugAt(pos token.Pos, what string) {
	panic(bail{&InternalError{fmt.Sprintf("%s (%s)", what, l.conf.Fset.Position(pos))}})
}

// need records a feature at the site that requires it, which is what makes
// Program.Features unable to disagree with the calls actually emitted.
func (l *lowerer) need(f builtins.Feature) {
	l.features = l.features.With(f)
}

// ownerOf answers which module a named type's shape belongs to. A type
// declared in another package lands in that package's module and is
// referenced across; anything anonymous lands where it was first needed.
func (l *lowerer) ownerOf(t types.Type, fallback *Module) *Module {
	if n := types.AsNamed(t); n != nil && n.Obj() != nil && n.Obj().Pkg() != nil {
		if m, ok := l.modOf[n.Obj().Pkg()]; ok {
			return m
		}
	}
	return fallback
}

// seed enumerates the roots. Exported generic functions are not among them.
func (l *lowerer) seed() error {
	switch l.conf.Mode {
	case ModeTest:
		fn := l.findFunc(l.conf.Test, func(f *types.Func) bool { return f.IsTest() })
		if fn == nil {
			return &InternalError{"no test function named " + l.conf.Test}
		}
		l.enqueue(fn, nil)
		l.prog.Entry = l.testShim(fn)
	default:
		fn := l.findFunc("main", func(f *types.Func) bool { return f.IsEntry() })
		if fn == nil {
			return &InternalError{"no func main() in package main"}
		}
		l.enqueue(fn, nil)
		l.prog.Entry = l.mainShim(fn)
	}

	// Every exported non-generic function is also a root in ModeProgram:
	// something outside the program may call it. That is the only reason to
	// walk beyond main.
	if l.conf.Mode == ModeProgram {
		for _, u := range l.units {
			s := u.Pkg.Scope()
			names := s.Names()
			sort.Strings(names) // map iteration must not decide emission order
			for _, n := range names {
				fn, ok := s.Lookup(n).(*types.Func)
				if !ok || fn.IsEntry() {
					continue
				}
				if sig := fn.Signature(); sig != nil && len(genericParams(fn)) == 0 {
					l.enqueue(fn, nil)
				}
			}
		}
	}
	return nil
}

func (l *lowerer) findFunc(name string, ok func(*types.Func) bool) *types.Func {
	for _, u := range l.units {
		if obj, is := u.Pkg.Scope().Lookup(name).(*types.Func); is && ok(obj) {
			return obj
		}
	}
	return nil
}