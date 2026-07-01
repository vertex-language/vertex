package codegen

import (
	"fmt"

	"github.com/vertex-language/ir/machine"
	objbridge "github.com/vertex-language/ir/machine/object"

	"github.com/vertex-language/vertex/target"
)

// ToFuncs runs isel + regalloc + binary encoding for every function in the
// MIR module, returning one AssembledFunc per function ready for the object
// builder. riscv64 object assembly is not yet wired into the object
// bridge, so -c/-emit-obj (and therefore link modes) return a clear error
// for that arch; -emit-asm works fine.
func ToFuncs(m *machine.Module, tri target.Triple, opts Options) ([]objbridge.AssembledFunc, error) {
	switch tri.Arch {
	case "amd64":
		progs, err := iselAMD64(m, opts)
		if err != nil {
			return nil, err
		}
		fns := make([]objbridge.AssembledFunc, 0, len(progs))
		for _, p := range progs {
			fn, err := objbridge.AssembleAMD64(p)
			if err != nil {
				return nil, fmt.Errorf("amd64 encode %s: %w", p.Name, err)
			}
			fns = append(fns, fn)
		}
		return fns, nil

	case "arm64":
		progs, err := iselARM64(m, opts)
		if err != nil {
			return nil, err
		}
		fns := make([]objbridge.AssembledFunc, 0, len(progs))
		for _, p := range progs {
			fn, err := objbridge.AssembleARM64(p)
			if err != nil {
				return nil, fmt.Errorf("arm64 encode %s: %w", p.Name, err)
			}
			fns = append(fns, fn)
		}
		return fns, nil

	case "riscv64":
		return nil, fmt.Errorf("riscv64: object file emission not yet supported; use -emit-asm or -emit-mir")
	}
	return nil, fmt.Errorf("unsupported arch for code generation: %s", tri.Arch)
}