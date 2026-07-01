package linker

import (
	"fmt"

	linkerelf "github.com/vertex-language/linker/elf"

	"github.com/vertex-language/vertex/target"
)

type elfLinker struct {
	l    *linkerelf.Linker
	opts Options
}

func newELFLinker(tri target.Triple, opts Options) (Linker, error) {
	e := &elfLinker{l: linkerelf.NewLinker(elfArch(tri)), opts: opts}
	if err := e.l.AddObject("crt1.o", opts.CRT.CRT1); err != nil {
		return nil, fmt.Errorf("add crt1.o: %w", err)
	}
	if err := e.l.AddObject("crti.o", opts.CRT.CRTI); err != nil {
		return nil, fmt.Errorf("add crti.o: %w", err)
	}
	return e, nil
}

func (e *elfLinker) AddObject(name string, data []byte) error {
	return e.l.AddObject(name, data)
}

func (e *elfLinker) Link() ([]byte, error) {
	if err := e.l.AddObject("crtn.o", e.opts.CRT.CRTN); err != nil {
		return nil, fmt.Errorf("add crtn.o: %w", err)
	}
	for _, lib := range e.opts.DynLibs {
		if err := e.l.AddDynamicLibrary(lib.Name, lib.Bytes); err != nil {
			return nil, fmt.Errorf("add dynamic library %s: %w", lib.Name, err)
		}
	}
	return e.l.Link()
}

func elfArch(tri target.Triple) linkerelf.Arch {
	switch tri.Arch {
	case "arm64":
		return linkerelf.ArchARM64
	case "riscv64":
		return linkerelf.ArchRISCV64
	default:
		return linkerelf.ArchAMD64
	}
}