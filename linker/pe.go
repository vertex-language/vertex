package linker

import (
	"fmt"

	linkerpe "github.com/vertex-language/linker/pe"

	"github.com/vertex-language/vertex/target"
)

type peLinker struct {
	l    *linkerpe.Linker
	opts Options
}

func newPELinker(tri target.Triple, opts Options) (Linker, error) {
	l := linkerpe.NewLinker(peArch(tri))
	if opts.OutputType == OutputPIE {
		l.SetOutputType(linkerpe.OutputPIE) // DYNAMIC_BASE + HIGH_ENTROPY_VA + .reloc
	}
	l.SetEntryPoint(opts.EntryPoint)
	return &peLinker{l: l, opts: opts}, nil
}

func (p *peLinker) AddObject(name string, data []byte) error {
	return p.l.AddObject(name, data)
}

func (p *peLinker) Link() ([]byte, error) {
	for _, lib := range p.opts.DynLibs {
		if err := p.l.AddDynamicLibrary(lib.Name, lib.Bytes); err != nil {
			return nil, fmt.Errorf("add dynamic library %s: %w", lib.Name, err)
		}
	}
	return p.l.Link()
}

func peArch(tri target.Triple) linkerpe.Arch {
	if tri.Arch == "arm64" {
		return linkerpe.ArchARM64
	}
	return linkerpe.ArchAMD64
}