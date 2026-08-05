// vir.go
package vir

import (
	"fmt"

	"github.com/vertex-language/vertex/lower/hir"
	"github.com/vertex-language/vertex/token"
	vir "github.com/vertex-language/vvm/ir/vir"
)

// The mechanical translator. A *hir.Program in, one *vir.Module per
// originating Vertex package out.
//
// Every decision was already made in hir. What is left is translation,
// ordering, and naming. If anything here needs a type switch on "is this
// owning" or "was this transferred", that logic belongs upstream.
//
// Note the import alias: this package is named vir and imports vvm's vir.
// A package clause declares no identifier in its own file scope, so the
// alias is unambiguous, and it keeps the emitted-side spelling identical to
// vvm's own documentation — vir.I32, vir.Ptr, vir.OpCall.

// Target is the triple. It enters the compiler here and nowhere earlier:
// hir never saw one, and the BuildClause the packages were checked under is
// coarser by design — a tag names the OS, the arch and ABI come from the
// driver.
//
// Field-for-field identical to vvm.Target so a driver can convert rather
// than copy.
type Target struct {
	Arch  string
	OS    string
	ABI   string
	Tiers []string
}

type Config struct {
	Fset   *token.FileSet // resolves hir positions into loc lines
	Target Target         // written into every emitted module; see links()
}

// Error is this package's only error shape. There is no todo/internal split
// here as there is in hir: everything this package refuses is a malformed
// *hir.Program or a construct with no vir spelling at all, and both are
// bugs upstream rather than gaps in the user's source.
type Error struct{ What string }

func (e *Error) Error() string { return "lower/vir: " + e.What }

// runtimeLink is the dependency every Vertex module carries. vir §2.1 makes
// `target` mandatory whenever `link` is present, which is why Config.Target
// reaches every module rather than only the ones with a declare block.
const runtimeLink = "vertexrt"

// Lower translates a whole program. Modules come back in hir's topological
// order — the order vvm's importer expects.
func Lower(conf *Config, prog *hir.Program) (mods []*vir.Module, err error) {
	if conf == nil {
		return nil, &Error{"Lower: Config is required"}
	}
	if prog == nil {
		return nil, &Error{"Lower: nil program"}
	}

	// bail unwinds out of the block/instruction walk without threading an
	// error through every return. Nothing else in this package panics.
	defer func() {
		if r := recover(); r != nil {
			if b, ok := r.(bail); ok {
				mods, err = nil, b.err
				return
			}
			panic(r)
		}
	}()

	for _, m := range prog.Modules {
		ml := &moduleLowerer{
			conf: conf,
			prog: prog,
			src:  m,
			out:  vir.NewModule(m.Name),
			seen: map[string]bool{},
		}
		mods = append(mods, ml.lower())
	}
	return mods, nil
}

// Root names the module holding the entry point, for vvm's --root.
func Root(prog *hir.Program) string {
	if prog == nil || prog.Root == nil {
		return ""
	}
	return prog.Root.Name
}

type bail struct{ err error }

type moduleLowerer struct {
	conf *Config
	prog *hir.Program
	src  *hir.Module
	out  *vir.Module

	// imports is the emitted import set: hir's own list plus anything
	// discovered while translating a cross-module struct reference, which
	// hir does not record (it interns the shape, not the dependency).
	imports []string
	seen    map[string]bool
}

func (ml *moduleLowerer) bug(what string) {
	panic(bail{&Error{fmt.Sprintf("%s (module %s)", what, ml.src.Name)}})
}

// lower walks the sections in vir's mandatory §2.1 order. The Module holds
// one slice per section, so the emitted order is structural rather than
// call-order dependent — but the calls run in section order anyway, because
// a reader checking the invariant against the emitted text should be able
// to read this function top to bottom and see the same thing.
func (ml *moduleLowerer) lower() *vir.Module {
	// namespace is the import path, module is the PackageClause name:
	// semantics.md §1.3's "the path is a locator, the declared name is the
	// qualifier", carried to the linker unchanged.
	if ml.src.Path != "" {
		ml.out.SetNamespace(ml.src.Path)
	}
	t := ml.conf.Target
	ml.out.SetTarget(t.Arch, t.OS, t.ABI, t.Tiers...)

	ml.structs()
	// FunctionSignatures: nothing upstream produces one. hir todos on every
	// call through a function value, so there is no indirect call to type.
	ml.constants()
	ml.globals()
	ml.links()
	ml.externs()
	ml.declareImports()
	ml.functions()

	return ml.out
}

// needImport records a dependency discovered during translation. hir's own
// list is seeded first so first-reference order — and therefore byte
// reproducibility — survives.
func (ml *moduleLowerer) needImport(path string) {
	if path == "" || ml.seen[path] {
		return
	}
	ml.seen[path] = true
	ml.imports = append(ml.imports, path)
}