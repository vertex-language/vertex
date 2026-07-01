package codegen

import (
	"fmt"

	"github.com/vertex-language/ir/machine"

	asmAMD64 "github.com/vertex-language/ir/machine/asm/amd64"
	asmARM64 "github.com/vertex-language/ir/machine/asm/arm64"
	asmRV64 "github.com/vertex-language/ir/machine/asm/riscv64"

	lowerAMD64 "github.com/vertex-language/ir/machine/lower/amd64"
	lowerARM64 "github.com/vertex-language/ir/machine/lower/arm64"
	lowerRISCV64 "github.com/vertex-language/ir/machine/lower/riscv64"
)

// Each function calls the real machine/lower/{arch} package, which performs
// instruction selection, naive register allocation, and frame layout,
// returning one arch-specific asm.Program per MIR function.

func iselAMD64(m *machine.Module, _ Options) ([]*asmAMD64.Program, error) {
	l, err := lowerAMD64.NewLower(m)
	if err != nil {
		return nil, fmt.Errorf("isel/amd64: %w", err)
	}
	return l.Funcs, nil
}

func iselARM64(m *machine.Module, _ Options) ([]*asmARM64.Program, error) {
	l, err := lowerARM64.NewLower(m)
	if err != nil {
		return nil, fmt.Errorf("isel/arm64: %w", err)
	}
	return l.Funcs, nil
}

func iselRISCV64(m *machine.Module, _ Options) ([]*asmRV64.Program, error) {
	l, err := lowerRISCV64.NewLower(m)
	if err != nil {
		return nil, fmt.Errorf("isel/riscv64: %w", err)
	}
	return l.Funcs, nil
}