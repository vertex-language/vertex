package main

import (
	"fmt"
	"strings"

	"github.com/vertex-language/objectfile/object"

	linkerelf   "github.com/vertex-language/linker/elf"
	linkermacho "github.com/vertex-language/linker/macho"
	linkerpe    "github.com/vertex-language/linker/pe"
)

// triple is the parsed form of a "os-arch" target string.
type triple struct {
	os   string // linux | darwin | windows | freestanding
	arch string // amd64 | arm64 | riscv64
}

func parseTriple(s string) (triple, error) {
	parts := strings.SplitN(strings.ToLower(strings.TrimSpace(s)), "-", 2)
	if len(parts) != 2 {
		return triple{}, fmt.Errorf("invalid target %q: expected <os>-<arch>", s)
	}
	t := triple{os: parts[0], arch: parts[1]}

	switch t.os {
	case "linux", "darwin", "windows", "freestanding":
	default:
		return triple{}, fmt.Errorf("unknown OS %q in target %q", t.os, s)
	}
	switch t.arch {
	case "amd64", "arm64", "riscv64":
	default:
		return triple{}, fmt.Errorf("unknown arch %q in target %q", t.arch, s)
	}
	return t, nil
}

// virTargetString returns the target string understood by vertex.Module.SetTarget.
// The MIR lowerer infers the target arch by substring matching against this value.
func (t triple) virTargetString() string {
	return t.os + "-" + t.arch
}

// buildTags returns the frontend build-tag slice (e.g. ["linux", "amd64"]).
func (t triple) buildTags() []string { return []string{t.os, t.arch} }

// objectTarget maps the triple to the objectfile/object.Target constant.
func (t triple) objectTarget() (object.Target, error) {
	switch t.os + "-" + t.arch {
	case "linux-amd64":
		return object.TargetLinuxAMD64, nil
	case "linux-arm64":
		return object.TargetLinuxARM64, nil
	case "linux-riscv64":
		return object.TargetLinuxRISCV64, nil
	case "darwin-amd64":
		return object.TargetDarwinAMD64, nil
	case "darwin-arm64":
		return object.TargetDarwinARM64, nil
	case "windows-amd64":
		return object.TargetWindowsAMD64, nil
	case "windows-arm64":
		return object.TargetWindowsARM64, nil
	case "freestanding-amd64":
		return object.TargetFreestandingAMD64, nil
	case "freestanding-arm64":
		return object.TargetFreestandingARM64, nil
	case "freestanding-riscv64":
		return object.TargetFreestandingRISCV64, nil
	}
	return object.Target{}, fmt.Errorf("no objectfile target for %s-%s", t.os, t.arch)
}

// ── Linker arch helpers ───────────────────────────────────────────────────────

func (t triple) elfArch() linkerelf.Arch {
	switch t.arch {
	case "arm64":
		return linkerelf.ArchARM64
	case "riscv64":
		return linkerelf.ArchRISCV64
	default:
		return linkerelf.ArchAMD64
	}
}

func (t triple) machoArch() linkermacho.Arch {
	if t.arch == "arm64" {
		return linkermacho.ArchARM64
	}
	return linkermacho.ArchAMD64
}

func (t triple) peArch() linkerpe.Arch {
	if t.arch == "arm64" {
		return linkerpe.ArchARM64
	}
	return linkerpe.ArchAMD64
}

func isWindowsTarget(target string) bool {
	return strings.HasPrefix(strings.ToLower(target), "windows-")
}