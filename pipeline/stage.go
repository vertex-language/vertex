package pipeline

import (
	"io"

	"github.com/vertex-language/ir/machine"
	objbridge "github.com/vertex-language/ir/machine/object"
	"github.com/vertex-language/ir/vertex"
	"github.com/vertex-language/ir/vertex/ast"

	"github.com/vertex-language/vertex/codegen"
	"github.com/vertex-language/vertex/target"
)

// Unit is one compiled module's progress through the pipeline: a package
// parsed from disk (or synthesized in memory, for the test runner), and
// whatever each stage has produced from it so far.
//
// A dependency module and the project's own root package are both Units,
// handled identically by every stage — including any language runtime a
// project happens to depend on. There is nothing here that special-cases
// "the runtime": it's just another entry in State.Units, resolved into
// the build by the same pkg.Graph as everything else the project imports.
type Unit struct {
	ModulePath string // mod.ModulePath as a string; "" for a standalone, module-less compile
	Dir        string // source directory or file; "" for an in-memory package
	IsRoot     bool

	Pkg    *ast.Package
	VIR    *vertex.Module
	VIRErr error // set only on the root unit: a non-nil VIRErr alongside a non-nil VIR is a "soft" failure -emit-vir/-emit-vbytes may still render
	MIR    *machine.Module
	ASM    string
	Funcs  []objbridge.AssembledFunc
	Obj    []byte
}

// State carries every stage's intermediate artifacts through one compile
// invocation. Fields are filled in as stages run; a stage may read any
// upstream field but should only write the fields it owns.
type State struct {
	Units   []*Unit // build order: dependencies before dependents; the last element is the root
	Triple  target.Triple
	Opts    codegen.Options
	LibDirs []string // from nativelibs.SearchDirs; already includes any sysroot and any pkg.LibResult cache dirs

	Exe []byte

	Sink io.Writer // non-nil in dump mode: Run writes a banner and each stage's rendered artifact here
	UpTo string    // stage Name to stop after; "" runs every stage
}

func (st *State) Root() *Unit { return st.Units[len(st.Units)-1] }

// Stage is one named step of the pipeline. Dump mode and every other emit
// mode run the exact same Stage list through Run — dump mode just also
// captures a rendered banner and artifact per stage into st.Sink. This is
// what replaces the old pipeline.go/dump.go fork: one list, one place a
// stage's behavior is defined.
type Stage struct {
	Name string
	Run  func(*State) error
}