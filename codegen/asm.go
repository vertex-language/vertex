package codegen

import (
	"fmt"
	"strings"

	"github.com/vertex-language/ir/machine"

	txtAMD64 "github.com/vertex-language/ir/machine/encoding/text/amd64"
	txtARM64 "github.com/vertex-language/ir/machine/encoding/text/arm64"
	txtRISCV64 "github.com/vertex-language/ir/machine/encoding/text/riscv64"

	"github.com/vertex-language/vertex/target"
)

// ToASM runs isel + regalloc and returns human-readable assembly text for
// every function (-emit-asm). amd64 output is Intel-syntax; arm64 and
// riscv64 use GNU-as syntax, matching each arch's text printer.
func ToASM(m *machine.Module, tri target.Triple, opts Options) (string, error) {
	var sb strings.Builder

	switch tri.Arch {
	case "amd64":
		progs, err := iselAMD64(m, opts)
		if err != nil {
			return "", err
		}
		for _, p := range progs {
			sb.WriteString(txtAMD64.Print(p))
			sb.WriteByte('\n')
		}

	case "arm64":
		progs, err := iselARM64(m, opts)
		if err != nil {
			return "", err
		}
		for _, p := range progs {
			sb.WriteString(txtARM64.Print(p))
			sb.WriteByte('\n')
		}

	case "riscv64":
		progs, err := iselRISCV64(m, opts)
		if err != nil {
			return "", err
		}
		for _, p := range progs {
			sb.WriteString(txtRISCV64.Print(p))
			sb.WriteByte('\n')
		}

	default:
		return "", fmt.Errorf("unsupported arch for assembly emit: %s", tri.Arch)
	}

	return sb.String(), nil
}