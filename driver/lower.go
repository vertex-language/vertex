// driver/lower.go
//
// This file is the seam between the Vertex frontend and vvm. Every
// construction of a lower/hir or lower/vir Config lives here and nowhere
// else, so a change to either package's option struct is a change to one
// file.
package driver

import (
	"fmt"

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
		Sizes: t.Sizes,
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
func virConfig(t Target) *lowervir.Config {
	return &lowervir.Config{
		Target: t.VVM.String(),
	}
}