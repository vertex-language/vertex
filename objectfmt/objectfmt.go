package objectfmt

import (
	"fmt"

	"github.com/vertex-language/ir/machine"
	objbridge "github.com/vertex-language/ir/machine/object"

	"github.com/vertex-language/objectfile/coff"
	"github.com/vertex-language/objectfile/elf"
	"github.com/vertex-language/objectfile/macho"
	"github.com/vertex-language/objectfile/object"

	"github.com/vertex-language/vertex/target"
)

// Marshal serializes sections into tri's native relocatable object format:
// ELF for linux/freestanding, Mach-O for darwin, COFF for windows.
func Marshal(tri target.Triple, sections []object.Section) ([]byte, error) {
	tgt, err := objectTargetFor(tri)
	if err != nil {
		return nil, err
	}

	type objectFile interface {
		AddSection(object.Section)
		Serialize() ([]byte, error)
	}
	var f objectFile
	switch tri.OS {
	case "linux", "freestanding":
		f = elf.NewFile(tgt)
	case "darwin":
		f = macho.NewFile(tgt)
	case "windows":
		f = coff.NewFile(tgt)
	default:
		return nil, fmt.Errorf("objectfmt: unsupported OS %q", tri.OS)
	}
	for _, s := range sections {
		f.AddSection(s)
	}
	return f.Serialize()
}

func objectTargetFor(tri target.Triple) (object.Target, error) {
	switch tri.OS + "-" + tri.Arch {
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
	return object.Target{}, fmt.Errorf("no objectfile target for %s-%s", tri.OS, tri.Arch)
}

// BuildSections assembles the section list for one compiled unit's object
// file: a single .text built from fns, plus m's data sections, plus a
// windows/amd64 .pdata/.xdata exception directory when applicable.
//
// Windows x64 requires a .pdata exception directory for the loader to
// accept the binary and for stack unwinding to work. ARM64 uses a
// different packed format and will be added separately.
func BuildSections(fns []objbridge.AssembledFunc, m *machine.Module) []object.Section {
	secs := make([]object.Section, 0, 4)
	secs = append(secs, objbridge.BuildText(fns))
	secs = append(secs, objbridge.DataSections(m)...)

	if m.OS == "windows" && m.Arch == machine.AMD64 {
		pdata, xdata := objbridge.BuildUnwind(fns)
		if len(pdata.Code) > 0 {
			secs = append(secs, pdata)
		}
		if len(xdata.Code) > 0 {
			secs = append(secs, xdata)
		}
	}
	return secs
}