// driver/lower.go
//
// This file is the seam between the Vertex frontend and vvm. Every
// construction of a lower/hir or lower/vir Config lives here and nowhere
// else, so a change to either package's option struct is a change to one
// file.
package driver

import (
	"fmt"

	"github.com/vertex-language/vertex/builtins"
	"github.com/vertex-language/vertex/lower/hir"
	lowervir "github.com/vertex-language/vertex/lower/vir"
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
// — recorded via hir's own l.need at every builtin call site, so this set
// can never disagree with what got emitted. builtins.Modules resolves
// that set into real *vir.Module values, constructed directly in Go
// against builder.go's API rather than parsed from a hand-maintained .vir
// text source per target. Appending them here, after lower/vir's own
// modules and before this function returns, is what makes an `import
// "memory"` line in hir's output resolve to something real once vvm's
// importer walks the returned module list.
func Lower(opts *Options, t Target, pkgs []*Package, lo LowerOptions) ([]*virmod.Module, string, error) {
	units := hirUnits(pkgs)

	prog, err := hir.Lower(hirConfig(t, lo), units)
	if err != nil {
		return nil, "", fmt.Errorf("lower/hir: %w", err)
	}
	opts.logf("hir: program lowered from %d unit(s)", len(units))

	modules, err := lowervir.Lower(virConfig(t), prog)
	if err != nil {
		return nil, "", fmt.Errorf("lower/vir: %w", err)
	}
	if len(modules) == 0 {
		return nil, "", fmt.Errorf("lower/vir produced no modules")
	}

	builtinModules := builtins.Modules(prog.Features, virConfig(t).Target)
	if len(builtinModules) > 0 {
		opts.logf("builtins: %d module(s) resolved from feature set", len(builtinModules))
		modules = append(modules, builtinModules...)
	}

	return modules, lowervir.Root(prog), nil
}

// --- the seam ---------------------------------------------------------------
//
// The three functions below are the only place this repository names a
// field of hir.Config, hir.Unit, or lower/vir.Config. They exist as
// separate functions rather than inline literals precisely so that
// tracking a rename in either package means editing here and only here.

func hirConfig(t Target, lo LowerOptions) *hir.Config {
	c := &hir.Config{
		// Layout is the one target-shaped fact hir consumes, and the
		// driver is what chooses it — hir never sees a triple.
		Sizes: *t.Sizes,
	}
	if lo.TestFunc != "" {
		c.Mode = hir.ModeTest
		c.TestFunc = lo.TestFunc
	}
	return c
}

func hirUnits(pkgs []*Package) []*hir.Unit {
	units := make([]*hir.Unit, 0, len(pkgs))
	for _, p := range pkgs {
		units = append(units, &hir.Unit{
			Path:  p.Path,
			Name:  p.Name,
			Files: p.Files,
			Types: p.Types,
			Info:  p.Info,
		})
	}
	return units
}

// virConfig hands lower/vir the triple it writes into any module carrying
// a `link` section (§2.1 makes `target` mandatory exactly there). The OS
// half already came from each file's build clause; the arch and ABI are
// the driver's to supply, which is the whole reason Config has this field.
//
// Also the source builtins.Modules reads its vir.Target from, via
// virConfig(t).Target in Lower above — the builtin modules must declare
// the identical triple lower/vir wrote into the program's own modules, or
// vvm's importer would see a target mismatch across the module graph.
func virConfig(t Target) *lowervir.Config {
	return &lowervir.Config{
		Target: virmod.Target{
			Arch:  t.VVM.Arch,
			OS:    t.VVM.OS,
			ABI:   t.VVM.ABI,
			Tiers: t.VVM.Tier,
		},
	}
}