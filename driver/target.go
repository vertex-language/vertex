// driver/target.go
package driver

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/vertex-language/vertex/token"
	"github.com/vertex-language/vertex/types"
	"github.com/vertex-language/vvm"
)

// DefaultMinOSVersion is what a darwin target gets when the caller passes
// none. vvm rejects an empty MinOSVersion for every Mach-O target — Apple's
// triple grammar has no unversioned form — so something has to be chosen,
// and this matches what vvm's own host detection picks.
const DefaultMinOSVersion = "14.0"

// Target is the driver's own resolved target: a Vertex-facing name, the
// build tag that decides what is even grammatical, the layout model, and
// the vvm.Target the backend actually wants.
//
// This is deliberately not vvm.Target with extra fields bolted on. A
// vertex triple (`linux-amd64`) and a vvm triple (`x86_64-linux-gnu`)
// are different vocabularies with different owners, and this struct is the
// one place that translates — the same discipline vvm's dispatch.go keeps
// between vvm.Target and each linker's own Target.
type Target struct {
	Name   string
	Tag    token.BuildTag
	Sizes  *types.Sizes
	VVM    vvm.Target
	Linked bool   // vvm can produce a finished image for this cell
	Note   string // shown by `vertex targets`
}

// forTag re-derives the layout model for a different build tag.
//
// Tag and Sizes must move together: types.SizesFor answers from the tag,
// and the packages are checked under that same tag, so changing one
// without the other would hand hir a Sizes the checker never used. The
// test runner is the only caller — `build test` is the one tag that
// changes what is grammatical.
func (t Target) forTag(tag token.BuildTag) Target {
	t.Tag = tag
	t.Sizes = types.SizesFor(tag)
	return t
}

type targetSpec struct {
	arch, os, abi string
	tag           token.BuildTag
	flat          bool
	linked        bool
	note          string
}

// targetTable is every (vertex triple -> vvm triple) pair this toolchain
// can honor. It is intentionally much shorter than the 0.4.0 README's
// list: a triple only appears here if vvm has a cpu/lower, an object, and
// a registered linker (or a flat writer) for it.
//
// ABI note: darwin cells use vir §7.1's canonical "macho" ABI token — the
// row the spec provides for exactly this case ("Apple convention for
// targets without an OS-specific ABI above") — rather than the English
// word "none", which is not a member of vir.CanonicalABI and fails
// ir/verify's checkTarget with "unknown abi \"none\" (§7.1)". An empty
// string would also pass verification (ABI is optional in the grammar),
// but "macho" is the literal, named answer the spec already gives for
// this exact case, so it's used instead of relying on emptiness meaning
// the same thing by accident.
//
// Deliberately absent, and why:
//
//   - linux-riscv64, and every powerpc/mips/loongarch/s390x spelling:
//     valid .vir target triples per the spec, with no cpu/lower, object,
//     or linker implementation in vvm.
//   - linux-386 / windows-386: vvm can produce x86 ELF *object* bytes but
//     registers no ELF linker backend for x86, so no image comes out.
//   - 32-bit arm: flat only in vvm, and lower/hir emits a hosted main,
//     which flat can't carry an entry stub for — the pair is unusable
//     end to end, so it isn't offered.
//   - browser/wasm, browser/js, android: no backend at all.
var targetTable = map[string]targetSpec{
	"linux-amd64":   {"x86_64", "linux", "gnu", token.TagLinux, false, true, "ELF, fully linked"},
	"linux-arm64":   {"aarch64", "linux", "gnu", token.TagLinux, false, true, "ELF, fully linked"},
	"darwin-amd64":  {"x86_64", "macos", "macho", token.TagDarwin, false, true, "Mach-O, ad-hoc signed"},
	"darwin-arm64":  {"aarch64", "macos", "macho", token.TagDarwin, false, true, "Mach-O, ad-hoc signed"},
	"windows-amd64": {"x86_64", "windows", "msvc", token.TagWindows, false, true, "PE; cross-building needs the target DLLs on disk"},
	"windows-arm64": {"aarch64", "windows", "msvc", token.TagWindows, false, true, "PE; cross-building needs the target DLLs on disk"},

	// Freestanding is flat-only: no loader, no relocations, no linker.
	// vvm refuses a synthesized crt stub against flat, so a freestanding
	// program must name its entry function _start rather than main, and
	// must be a single package (flat forbids the relocations a
	// cross-module call needs). ABI stays "none" here, unlike the darwin
	// rows above: os is also "none" (bare metal), there's no Mach-O
	// convention in play, and this hasn't been confirmed against
	// vir.CanonicalABI's actual contents — if it hits the same "unknown
	// abi" error, it needs the same kind of fix, but with whatever token
	// the spec actually names for a bare-metal target, not a guess.
	"freestanding-amd64": {"x86_64", "none", "none", token.TagNone, true, true, "flat image; single package, entry must be _start"},
	"freestanding-arm64": {"aarch64", "none", "none", token.TagNone, true, true, "flat image; single package, entry must be _start"},
}

type targetRequest struct {
	Name         string
	MinOSVersion string
	Shared       bool
	FlatBase     uint64
}

// ResolveTarget turns a vertex triple into a fully populated Target,
// failing loudly on anything the table doesn't name rather than passing an
// unknown triple down to vvm and letting it surface as an unrecognized-os
// error three layers deeper.
func ResolveTarget(req targetRequest) (Target, error) {
	name := req.Name
	if name == "" {
		host, err := HostTargetName()
		if err != nil {
			return Target{}, err
		}
		name = host
	}

	spec, ok := targetTable[name]
	if !ok {
		return Target{}, fmt.Errorf(
			"unknown target %q (known: %s) — run `vertex targets` for what each one produces",
			name, joinNames())
	}

	t := Target{
		Name: name,
		Tag:  spec.tag,

		// Sizes comes from the build tag, not from the arch. types.SizesFor
		// splits on TagJS/TagWasm (32-bit word) and answers 64-bit for
		// everything else, because §2.3 ties `int`/`uint` to the target's
		// pointer width and a file's width answer must come from its own
		// tag. That the split happens to be invisible across this table —
		// every cell here is 64-bit, and the 32-bit cells are the absent
		// ones listed above — is a property of today's table, not a reason
		// to derive layout from spec.arch: the tag is also what the
		// packages are checked under, and hir.Config.Sizes has to match
		// that, not the machine word.
		Sizes: types.SizesFor(spec.tag),

		Linked: spec.linked,
		Note:   spec.note,
		VVM: vvm.Target{
			Arch:            spec.arch,
			OS:              spec.os,
			ABI:             spec.abi,
			Flat:            spec.flat,
			FlatBaseAddress: req.FlatBase,
		},
	}

	if spec.os == "macos" {
		v := req.MinOSVersion
		if v == "" {
			v = DefaultMinOSVersion
		}
		t.VVM.MinOSVersion = v
	} else if req.MinOSVersion != "" {
		return Target{}, fmt.Errorf(
			"-min-os-version applies to darwin targets only; %s has no versioned triple form", name)
	}

	if req.Shared {
		if spec.flat {
			return Target{}, fmt.Errorf(
				"%s produces a flat image, which has no loader and therefore no shared-library form",
				name)
		}
		t.VVM.Kind = vvm.OutputSharedLibrary
	}

	if req.FlatBase != 0 && !spec.flat {
		return Target{}, fmt.Errorf(
			"-flat-base applies to freestanding targets only; %s links a real container format", name)
	}

	return t, nil
}

// HostTargetName reports the table entry matching the machine vertex is
// running on, so an unspecified -target needs no configuration.
func HostTargetName() (string, error) {
	var arch string
	switch runtime.GOARCH {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	default:
		return "", fmt.Errorf(
			"no supported target for this host's architecture (%s); pass -target explicitly",
			runtime.GOARCH)
	}

	var os string
	switch runtime.GOOS {
	case "linux":
		os = "linux"
	case "darwin":
		os = "darwin"
	case "windows":
		os = "windows"
	default:
		return "", fmt.Errorf(
			"no supported target for this host's OS (%s); pass -target explicitly", runtime.GOOS)
	}

	name := os + "-" + arch
	if _, ok := targetTable[name]; !ok {
		return "", fmt.Errorf("this host (%s) has no entry in the target table", name)
	}
	return name, nil
}

// Targets lists every known target, sorted, for `vertex targets`.
func Targets() []Target {
	out := make([]Target, 0, len(targetTable))
	for name := range targetTable {
		t, err := ResolveTarget(targetRequest{Name: name})
		if err != nil {
			continue
		}
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func joinNames() string {
	names := make([]string, 0, len(targetTable))
	for n := range targetTable {
		names = append(names, n)
	}
	sort.Strings(names)
	s := ""
	for i, n := range names {
		if i > 0 {
			s += ", "
		}
		s += n
	}
	return s
}