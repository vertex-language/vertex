// driver/lower.go
//
// This file is the seam between the Vertex frontend and vvm. Every
// construction of a lower/hir or lower/vir Config lives here and nowhere
// else, so a change to either package's option struct is a change to one
// file.
package driver

import (
	"errors"
	"fmt"

	"github.com/vertex-language/vertex/lower/hir"
	lowervir "github.com/vertex-language/vertex/lower/vir"
	"github.com/vertex-language/vertex/token"
	virmod "github.com/vertex-language/vvm/ir/vir"
)

// LowerOptions carries the driver-level knobs the lowering chain honors.
type LowerOptions struct {
	// TestFunc names the single test function to seed monomorphization
	// from, instead of main. Empty means an ordinary program.
	//
	// hir monomorphizes from one root, and under ModeTest that root is one
	// test function — which is why the test runner lowers and builds once
	// per test rather than emitting one binary holding all of them.
	TestFunc string
}

// Lower runs hir.Lower then lower/vir.Lower, returning the modules in
// hir's topological order and the name of the module holding the entry
// point.
//
// Both stages are called and neither is second-guessed: hir makes every
// decision (ownership, monomorphization, control-flow flattening) and
// lower/vir is mechanical afterward. If this function ever needs to
// inspect what came out, that inspection belongs upstream.
//
// A third step follows both: hir.Program.Features names which builtin
// modules (memory, panic, ...) the program's emitted calls actually need
// — recorded at every builtin call site, so this set can never disagree
// with what got emitted. FeatureSet.BuildModules resolves that set into
// real *vir.Module values, constructed directly in Go against builder.go's
// API rather than parsed from a hand-maintained .vir text source per
// target. Appending them here, after lower/vir's own modules and before
// this function returns, is what makes an `import "memory"` line in hir's
// output resolve to something real once vvm's importer walks the returned
// module list.
func Lower(opts *Options, t Target, ld *Loaded, lo LowerOptions) ([]*virmod.Module, string, error) {
	units := hirUnits(ld.Packages)

	prog, err := hir.Lower(hirConfig(ld.Fset, t, lo), units)
	if err != nil {
		// hir draws two error shapes and they mean different things to a
		// user: `todo:` is a valid program this compiler can't lower yet,
		// `internal:` is a bug in analyzer or hir. Neither is a compile
		// error against the source, so neither goes through diag — but
		// only the first is worth phrasing as a limitation.
		var todo *hir.TodoError
		if errors.As(err, &todo) {
			return nil, "", fmt.Errorf("unsupported construct: %w", err)
		}
		return nil, "", fmt.Errorf("lower/hir: %w", err)
	}
	opts.logf("hir: program lowered from %d unit(s)", len(units))

	vt := virTarget(t)

	modules, err := lowervir.Lower(virConfig(ld, vt), prog)
	if err != nil {
		return nil, "", fmt.Errorf("lower/vir: %w", err)
	}
	if len(modules) == 0 {
		return nil, "", fmt.Errorf("lower/vir produced no modules")
	}

	builtinModules := prog.Features.BuildModules(vt)
	if len(builtinModules) > 0 {
		opts.logf("builtins: %d module(s) resolved from feature set", len(builtinModules))
		modules = append(modules, builtinModules...)
	}

	return modules, lowervir.Root(prog), nil
}

// --- the seam ---------------------------------------------------------------
//
// The four functions below are the only place this repository names a
// field of hir.Config, hir.Unit, or lower/vir.Config. They exist as
// separate functions rather than inline literals precisely so that
// tracking a rename in either package means editing here and only here.

func hirConfig(fset *token.FileSet, t Target, lo LowerOptions) *hir.Config {
	c := &hir.Config{
		// hir resolves no positions itself; the FileSet is only what its
		// loc lines are rendered against.
		Fset: fset,

		// Layout is the one target-shaped fact hir consumes, and the
		// driver is what chooses it — hir never sees a triple. It must
		// match the tag the packages were *checked* under, since `int`
		// and `uint` are the target's pointer width (§2.3) and a mismatch
		// silently changes every layout the checker already committed to.
		// target.go derives both from the same targetSpec.tag, which is
		// what makes that guarantee structural rather than a convention.
		Sizes: t.Sizes,

		Mode: hir.ModeProgram,
	}
	if lo.TestFunc != "" {
		c.Mode = hir.ModeTest
		c.Test = lo.TestFunc
	}
	return c
}

// hirUnits converts checked packages into hir's input shape.
//
// A Unit is the *types.Package, its Info side table, and its files —
// nothing else. Path and Name deliberately don't appear: §1.3 makes the
// declared package name the qualifier and the import path a mere locator,
// and hir reads the former off the types.Package itself rather than
// trusting a second copy that could drift.
func hirUnits(pkgs []*Package) []*hir.Unit {
	units := make([]*hir.Unit, 0, len(pkgs))
	for _, p := range pkgs {
		units = append(units, &hir.Unit{
			Pkg:   p.Types,
			Info:  p.Info,
			Files: p.Files,
		})
	}
	return units
}

// virTarget builds the vvm-side triple: what ends up written into any
// module carrying a `link` section (§2.1 makes `target` mandatory exactly
// there, and every Vertex module links vertexrt). The OS half already came
// from each file's build clause; the arch and ABI are the driver's to
// supply.
//
// This is also what builtins.FeatureSet.BuildModules is handed, so the
// builtin modules declare the identical triple lower/vir wrote into the
// program's own modules — otherwise vvm's importer sees a target mismatch
// across the module graph.
func virTarget(t Target) virmod.Target {
	return virmod.Target{
		Arch:  t.VVM.Arch,
		OS:    t.VVM.OS,
		ABI:   t.VVM.ABI,
		Tiers: t.VVM.Tier,
	}
}

// virConfig hands lower/vir its FileSet and its triple.
//
// lower/vir declares its own Target rather than importing vvm's, "field
// for field identical so a driver can convert rather than copy" — so the
// conversion is spelled here, once, and is the single place this package
// asserts the two shapes still agree. If either grows a field, this line
// stops compiling, which is the intended failure mode.
func virConfig(ld *Loaded, vt virmod.Target) *lowervir.Config {
	return &lowervir.Config{
		Fset:   ld.Fset,
		Target: lowervir.Target(vt),
	}
}