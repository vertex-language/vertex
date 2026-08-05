// driver/emit.go
package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	virbinary "github.com/vertex-language/vvm/format/vbyte/binary"
	virtext "github.com/vertex-language/vvm/format/vbyte/text"
	virmod "github.com/vertex-language/vvm/ir/vir"

	"github.com/vertex-language/vvm"
)

// emit writes whatever Options.Emit selected and records the paths on res.
func emit(opts *Options, res *Result) error {
	switch opts.Emit {
	case EmitVIR, EmitVByte:
		return emitIR(opts, res)
	default:
		return emitBinary(opts, res)
	}
}

// emitBinary builds one native image through vvm and writes it.
func emitBinary(opts *Options, res *Result) error {
	if !res.Target.Linked {
		return fmt.Errorf("%s cannot be linked into an image by this toolchain", res.Target.Name)
	}
	if res.Target.VVM.Flat && len(res.Modules) > 1 {
		return fmt.Errorf(
			"%s produces a flat image, which forbids relocations — but this build has %d modules "+
				"and a cross-module call is exactly a relocation; build a single-package program "+
				"for a freestanding target",
			res.Target.Name, len(res.Modules))
	}

	bin, err := buildImage(res.Target, res.Modules, res.Root)
	if err != nil {
		return err
	}

	out := opts.Output
	if out == "" {
		out = outputBase(opts.Input)
		if res.Target.VVM.OS == "windows" {
			out += ".exe"
		}
	}
	if err := os.WriteFile(out, bin, 0o755); err != nil {
		return fmt.Errorf("writing %s: %w", out, err)
	}
	res.Outputs = []string{out}
	return nil
}

// buildImage is the one call into vvm's host pipeline.
//
// The single-module fast path matters: vvm.BuildModule refuses a module
// with imports (bare verify.Verify has no notion of a cross-module
// reference), and BuildModuleGraph refuses flat targets. Choosing by the
// module's own shape rather than by package count is what keeps both
// constraints satisfied without the driver duplicating either rule — and
// it stays correct now that a build's module list can hold builtin modules
// the package count never mentions.
func buildImage(t Target, modules []*virmod.Module, root string) ([]byte, error) {
	if len(modules) == 1 && len(modules[0].Imports) == 0 {
		bin, err := vvm.BuildModule(modules[0], t.VVM)
		if err != nil {
			return nil, err
		}
		return bin, nil
	}
	if root == "" {
		return nil, fmt.Errorf(
			"a multi-module build needs a root module to resolve the entry point, and none was " +
				"derived — pass -root <module>")
	}
	bin, err := vvm.BuildModuleGraph(modules, root, t.VVM)
	if err != nil {
		return nil, err
	}
	return bin, nil
}

// emitIR encodes each module through vvm's own codec and writes it.
//
// Both encoders come from format/vbyte, not from anything here: the point
// of `-emit-vir` is to show what vvm will actually read back, and a second
// printer could drift from the grammar model. That's the same "one
// canonical printer per format" rule vvm keeps internally.
func emitIR(opts *Options, res *Result) error {
	encoded := make([][]byte, len(res.Modules))
	for i, m := range res.Modules {
		var (
			b   []byte
			err error
		)
		if opts.Emit == EmitVIR {
			b, err = virtext.Encode(m)
		} else {
			b, err = virbinary.Encode(m)
		}
		if err != nil {
			return fmt.Errorf("encoding module %q as %s: %w", m.Name, opts.Emit, err)
		}
		encoded[i] = b
	}

	dests, err := irOutputPaths(opts, res)
	if err != nil {
		return err
	}
	for i, dest := range dests {
		if err := os.WriteFile(dest, encoded[i], 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", dest, err)
		}
	}
	res.Outputs = dests
	return nil
}

// irOutputPaths resolves -o against an emission that may fan out.
//
// The rule is vvm's own device-build rule, for the same reason: one Vertex
// program is one .vir per module, so -o may name a single file only when
// there's exactly one. Silently picking one module out of four to write,
// or overwriting the same path four times, are both worse than refusing.
func irOutputPaths(opts *Options, res *Result) ([]string, error) {
	names := make([]string, len(res.Modules))
	for i, m := range res.Modules {
		names[i] = m.Name + opts.Emit.ext()
	}

	o := opts.Output
	if o == "" {
		if len(res.Modules) == 1 {
			return []string{outputBase(opts.Input) + opts.Emit.ext()}, nil
		}
		return names, nil
	}

	if looksLikeDir(o) {
		if err := os.MkdirAll(o, 0o755); err != nil {
			return nil, fmt.Errorf("creating %s: %w", o, err)
		}
		dests := make([]string, len(names))
		for i, n := range names {
			dests[i] = filepath.Join(o, n)
		}
		return dests, nil
	}

	if len(res.Modules) == 1 {
		return []string{o}, nil
	}
	return nil, fmt.Errorf(
		"-o names a single file, but this build emits %d modules (%s) — pass a directory",
		len(names), strings.Join(names, ", "))
}

// looksLikeDir treats an existing directory, or any path written with a
// trailing separator, as a directory — "-o build/" creates it, "-o build"
// with no existing directory is a filename.
func looksLikeDir(p string) bool {
	if strings.HasSuffix(p, "/") || strings.HasSuffix(p, string(os.PathSeparator)) {
		return true
	}
	if fi, err := os.Stat(p); err == nil && fi.IsDir() {
		return true
	}
	return false
}