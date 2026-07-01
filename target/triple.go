package target

import (
	"fmt"
	"strings"
)

// Triple is the parsed form of an "os-arch" target string. It is a pure
// value type: no OS-specific object-format or linker-arch mapping lives
// here — those belong to objectfmt and linker respectively, each of which
// does its own tiny switch over OS/Arch rather than this package knowing
// about every downstream consumer's vocabulary.
type Triple struct {
	OS   string // linux | darwin | windows | freestanding
	Arch string // amd64 | arm64 | riscv64
}

func ParseTriple(s string) (Triple, error) {
	parts := strings.SplitN(strings.ToLower(strings.TrimSpace(s)), "-", 2)
	if len(parts) != 2 {
		return Triple{}, fmt.Errorf("invalid target %q: expected <os>-<arch>", s)
	}
	t := Triple{OS: parts[0], Arch: parts[1]}

	switch t.OS {
	case "linux", "darwin", "windows", "freestanding":
	default:
		return Triple{}, fmt.Errorf("unknown OS %q in target %q", t.OS, s)
	}
	switch t.Arch {
	case "amd64", "arm64", "riscv64":
	default:
		return Triple{}, fmt.Errorf("unknown arch %q in target %q", t.Arch, s)
	}
	return t, nil
}

// VirTargetString returns the target string understood by
// vertex.Module.SetTarget. The MIR lowerer infers the target arch by
// substring matching against this value.
func (t Triple) VirTargetString() string { return t.OS + "-" + t.Arch }

// BuildTags returns the frontend build-tag slice (e.g. ["linux", "amd64"]).
func (t Triple) BuildTags() []string { return []string{t.OS, t.Arch} }

// IsWindowsTarget reports whether target (an unparsed "-target" flag value)
// names a windows target. Exposed as a string check, rather than requiring
// a parsed Triple, because it's useful before a target string can be (or
// needs to be) fully parsed — e.g. picking ".exe" vs no extension while
// deriving a default output name.
func IsWindowsTarget(target string) bool {
	return strings.HasPrefix(strings.ToLower(target), "windows-")
}