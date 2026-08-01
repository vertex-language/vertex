// builtins/registry.go
package builtins

import vir "github.com/vertex-language/vvm/ir/vir"

// Feature identifies one builtin module's worth of functionality a
// program's emitted calls may require. hir.need records these as it
// lowers; Modules resolves the recorded set into the actual *vir.Module
// values a build must link in.
//
// NOTE: this assumes Feature/FeatureSet already exist elsewhere in this
// package with roughly this shape (hir references FeatureFor(Symbol) and
// FeatureSet.With(Feature)). Reconcile the two if the real definitions
// differ — this file only adds Modules and the two new constructors.
const (
	FeatMemory Feature = iota
	FeatPanic
	// existing features (FeatString, FeatSlice, FeatReactor, FeatFmt, ...)
	// stay wherever they're already declared.
)

// Modules resolves feats into the concrete builtin modules a build needs,
// for target t. Order is deterministic (declaration order below) so a
// build's module list — and therefore any diagnostic naming a module
// index — doesn't depend on map iteration.
func Modules(feats FeatureSet, t vir.Target) []*vir.Module {
	var out []*vir.Module
	if feats.Has(FeatMemory) {
		out = append(out, Memory(t))
	}
	if feats.Has(FeatPanic) {
		out = append(out, Panic(t))
	}
	return out
}