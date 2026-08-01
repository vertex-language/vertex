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
	Sizes  types.Sizes
	VVM    vvm.Target
	Linked bool   // vvm can produce a finished image for this cell
	Note   string // shown by `vertex targets`
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
	"darwin-amd64":  {"x86_64", "macos", "none", token.TagDarwin, false, true, "Mach-O, ad-hoc signed"},
	"darwin-arm64":  {"aarch64", "macos", "none", token.TagDarwin, false, true, "Mach-O, ad-hoc signed"},
	"windows-amd64": {"x86_64", "windows", "msvc", token.TagWindows, false, true, "PE; cross-building needs the target DLLs on disk"},
	"windows-arm64": {"aarch64", "windows", "msvc", token.TagWindows, false, true, "PE; cross-building needs the target DLLs on disk"},

	// Freestanding is flat-only: no loader, no relocations, no linker.
	// vvm refuses a synthesized crt stub against flat, so a freestanding
	// program must name its entry function _start rather than main, and
	// must be a single package (flat forbids the relocations a
	// cross-module call needs).
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
		Name:   name,
		Tag:    spec.tag,
		Sizes:  types.SizesFor(spec.tag),
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