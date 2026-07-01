package linker

import (
	"fmt"

	linkermacho "github.com/vertex-language/linker/macho"

	"github.com/vertex-language/vertex/target"
)

type machoLinker struct {
	l    *linkermacho.Linker
	opts Options
}

func newMachOLinker(tri target.Triple, opts Options) (Linker, error) {
	l := linkermacho.NewLinker(machoArch(tri))
	l.SetEntryPoint(opts.EntryPoint)
	l.AddSONeeded("libSystem.B.dylib")
	return &machoLinker{l: l, opts: opts}, nil
}

func (m *machoLinker) AddObject(name string, data []byte) error {
	return m.l.AddObject(name, data)
}

func (m *machoLinker) Link() ([]byte, error) {
	for _, lib := range m.opts.DynLibs {
		if lib.Bytes == nil {
			m.l.AddCachedDylib(lib.Name, m.opts.LibSymbols[lib.Name])
			continue
		}
		if err := m.l.AddDynamicLibrary(lib.Name, lib.Bytes); err != nil {
			return nil, fmt.Errorf("add dynamic library %s: %w", lib.Name, err)
		}
	}
	return m.l.Link()
}

func machoArch(tri target.Triple) linkermacho.Arch {
	if tri.Arch == "arm64" {
		return linkermacho.ArchARM64
	}
	return linkermacho.ArchAMD64
}