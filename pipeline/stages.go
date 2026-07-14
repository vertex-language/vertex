package pipeline

import (
	"fmt"
	"strings"

	"github.com/vertex-language/ir/machine"
	"github.com/vertex-language/ir/vertex/analyzer"
	mirlower "github.com/vertex-language/ir/vertex/lower/mir"
	virlower "github.com/vertex-language/ir/vertex/lower/vir"

	"github.com/vertex-language/vertex/codegen"
	"github.com/vertex-language/vertex/linker"
	"github.com/vertex-language/vertex/nativelibs"
	"github.com/vertex-language/vertex/objectfmt"
	"github.com/vertex-language/vertex/target"
)

const (
	StageVIR    = "Vertex IR (.vir)"
	StageMIR    = "Machine IR (.mir)"
	StageASM    = "Assembly (.s)"
	StageObject = "Object"
	StageLink   = "Link"
)

// Stages is the full pipeline, in order. Every emit mode runs a prefix of
// this list (see Run's upTo parameter); nothing here is specific to any
// one output format.
var Stages = []Stage{
	{StageVIR, runVIR},
	{StageMIR, runMIR},
	{StageASM, runASM},
	{StageObject, runObject},
	{StageLink, runLink},
}

// runVIR analyzes and lowers every unit to Vertex IR, in build order.
//
// virlower.NewLower takes the analyzer.Info for a single package plus a
// flat list of host Imports (ir/vertex/lower/vir.Import) — it has no
// notion of threading a dependency's already-lowered *vertex.Module or its
// exported names into a dependent unit's lowering. Cross-package symbol
// resolution therefore isn't something this stage can do yet; each unit is
// analyzed and lowered independently.
//
// A dependency unit must lower cleanly; there's no partial output to fall
// back to for a module nothing asked to see directly. The root unit is
// different: a "soft" lowering error (virErr != nil but vmod != nil) is
// tolerated when the caller only wants VIR/VBytes text (st.UpTo ==
// StageVIR) or is in dump mode (st.Sink != nil, which wants to show as
// much of the pipeline as possible); otherwise it's fatal.
func runVIR(st *State) error {
	for _, u := range st.Units {
		info, err := analyzer.Analyze(u.Pkg)
		if err != nil {
			return fmt.Errorf("analyzing %s: %w", unitLabel(u), err)
		}

		// customLibs (host Imports) is nil until vs.lib parsing/resolution
		// is threaded through here (per-module, likely off pkg.Graph) —
		// see resolveImportLib's doc comment for what plugs in.
		vmod, virErr := virlower.NewLower(u.Pkg, info, nil)

		tolerate := st.Sink != nil
		if virErr != nil && !u.IsRoot && !tolerate {
			return fmt.Errorf("resolving dependency %s: %w", unitLabel(u), virErr)
		}
		if u.IsRoot {
			u.VIRErr = virErr
		}
		if vmod == nil {
			return fmt.Errorf("vertex IR: %s: %w", unitLabel(u), virErr)
		}
		if u.IsRoot && virErr != nil && !tolerate && st.UpTo != StageVIR {
			return fmt.Errorf("vertex IR: %w (use -emit-vir/-emit-vbytes to inspect partial output)", virErr)
		}

		vmod.SetTarget(st.Triple.VirTargetString())
		u.VIR = vmod
	}
	return nil
}

func runMIR(st *State) error {
	for _, u := range st.Units {
		mirMod, err := mirlower.NewLower(u.VIR)
		if err != nil {
			return fmt.Errorf("MIR lowering %s: %w", unitLabel(u), err)
		}
		mirMod.OS = st.Triple.OS
		if err := machine.Verify(mirMod); err != nil {
			return fmt.Errorf("MIR verification %s: %w", unitLabel(u), err)
		}
		u.MIR = mirMod
	}
	return nil
}

func runASM(st *State) error {
	for _, u := range st.Units {
		text, err := codegen.ToASM(u.MIR, st.Triple, st.Opts)
		if err != nil {
			return fmt.Errorf("code generation %s: %w", unitLabel(u), err)
		}
		u.ASM = text
	}
	return nil
}

func runObject(st *State) error {
	for _, u := range st.Units {
		fns, err := codegen.ToFuncs(u.MIR, st.Triple, st.Opts)
		if err != nil {
			return fmt.Errorf("code generation %s: %w", unitLabel(u), err)
		}
		u.Funcs = fns
		sections := objectfmt.BuildSections(fns, u.MIR)
		objBytes, err := objectfmt.Marshal(st.Triple, sections)
		if err != nil {
			return fmt.Errorf("object serialization %s: %w", unitLabel(u), err)
		}
		u.Obj = objBytes
	}
	return nil
}

// runLink resolves every native dynamic library and CRT object the root
// unit's VIR imports declare, then adds every unit's object — in
// State.Units' own build order (dependencies before the root) — to the
// linker before finalizing. This linker resolves symbols across every
// added object in one pass rather than doing traditional archive-style
// resolution, so add order doesn't affect correctness; it's kept
// dependency-first purely because that's the order pkg.Graph already
// produces.
//
// No unit is treated specially here: whatever a project's own vs.mod
// happened to pull in — a runtime module included — contributes its
// object the same way every other dependency does.
func runLink(st *State) error {
	if st.Triple.OS == "freestanding" {
		return fmt.Errorf("cannot link a freestanding target; use -c/-emit-obj instead")
	}

	root := st.Root()
	libNames := nativelibs.ExtractDynLibs(root.VIR, st.Triple)
	dynLibs, err := nativelibs.ResolveLibs(libNames, st.Triple, st.LibDirs)
	if err != nil {
		return err
	}
	crt, err := nativelibs.ResolveCRT(st.Triple, st.LibDirs)
	if err != nil {
		return err
	}
	libSymbols := nativelibs.ExtractLibFuncSymbols(root.VIR, st.Triple.OS)

	l, err := linker.New(st.Triple, linker.Options{
		EntryPoint: entryPointFor(st.Triple),
		DynLibs:    toLinkerDynLibs(dynLibs),
		CRT:        linker.CRTObjects(crt),
		LibSymbols: libSymbols,
		OutputType: outputTypeFor(st.Triple),
	})
	if err != nil {
		return err
	}

	for _, u := range st.Units {
		if err := l.AddObject(objectFileName(u, st.Triple), u.Obj); err != nil {
			return fmt.Errorf("add object %s: %w", unitLabel(u), err)
		}
	}

	exe, err := l.Link()
	if err != nil {
		return fmt.Errorf("link: %w", err)
	}
	st.Exe = exe
	return nil
}

func entryPointFor(tri target.Triple) string {
	switch tri.OS {
	case "darwin":
		return "_main"
	case "windows":
		return "main"
	default:
		return "" // ELF's crt1.o already defines _start; no explicit entry symbol needed
	}
}

func outputTypeFor(tri target.Triple) linker.OutputType {
	if tri.OS == "windows" {
		return linker.OutputPIE
	}
	return linker.OutputDefault
}

func objectFileName(u *Unit, tri target.Triple) string {
	ext := ".o"
	if tri.OS == "windows" {
		ext = ".obj"
	}
	if u.IsRoot {
		return "main" + ext
	}
	name := u.ModulePath
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		name = name[i+1:]
	}
	if name == "" {
		name = "dep"
	}
	return name + ext
}

func toLinkerDynLibs(rs []nativelibs.ResolvedLib) []linker.DynLib {
	out := make([]linker.DynLib, len(rs))
	for i, r := range rs {
		out[i] = linker.DynLib{Name: r.Name, Bytes: r.Bytes}
	}
	return out
}

func unitLabel(u *Unit) string {
	if u.ModulePath != "" {
		return u.ModulePath
	}
	return "(root)"
}