package driver

import (
	"fmt"
	"strings"

	"github.com/vertex-language/ir/machine"

	asmAMD64 "github.com/vertex-language/ir/machine/asm/amd64"
	asmARM64 "github.com/vertex-language/ir/machine/asm/arm64"
	asmRV64  "github.com/vertex-language/ir/machine/asm/riscv64"

	lowerAMD64   "github.com/vertex-language/ir/machine/lower/amd64"
	lowerARM64   "github.com/vertex-language/ir/machine/lower/arm64"
	lowerRISCV64 "github.com/vertex-language/ir/machine/lower/riscv64"

	txtAMD64   "github.com/vertex-language/ir/machine/encoding/text/amd64"
	txtARM64   "github.com/vertex-language/ir/machine/encoding/text/arm64"
	txtRISCV64 "github.com/vertex-language/ir/machine/encoding/text/riscv64"

	objbridge "github.com/vertex-language/ir/machine/object"
)

// codegenOptions is forwarded to the instruction selector.
type codegenOptions struct {
	optLevel  int  // 0=none 1=light 2=full -1=size
	debugInfo bool
}

// ── instruction selection ─────────────────────────────────────────────────────
//
// Each function calls the real machine/lower/{arch} package, which performs
// instruction selection, naive register allocation, and frame layout, returning
// one arch-specific asm.Program per MIR function.

func iselAMD64(m *machine.Module, _ codegenOptions) ([]*asmAMD64.Program, error) {
	l, err := lowerAMD64.NewLower(m)
	if err != nil {
		return nil, fmt.Errorf("isel/amd64: %w", err)
	}
	return l.Funcs, nil
}

func iselARM64(m *machine.Module, _ codegenOptions) ([]*asmARM64.Program, error) {
	l, err := lowerARM64.NewLower(m)
	if err != nil {
		return nil, fmt.Errorf("isel/arm64: %w", err)
	}
	return l.Funcs, nil
}

func iselRISCV64(m *machine.Module, _ codegenOptions) ([]*asmRV64.Program, error) {
	l, err := lowerRISCV64.NewLower(m)
	if err != nil {
		return nil, fmt.Errorf("isel/riscv64: %w", err)
	}
	return l.Funcs, nil
}

// ── compileToFuncs ────────────────────────────────────────────────────────────

// compileToFuncs runs isel + regalloc + binary encoding for every function in
// the MIR module, returning one AssembledFunc per function ready for the object
// builder. riscv64 object assembly is not yet wired into the object bridge, so
// -c / -lc for that arch returns a clear error; -emit-asm works fine.
func compileToFuncs(m *machine.Module, tri triple, opts codegenOptions) ([]objbridge.AssembledFunc, error) {
	switch tri.arch {
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
		// The riscv64 object bridge is not yet wired into encoding/object.
		// Use -emit-asm to get assembly text, or -emit-mir to inspect MIR.
		return nil, fmt.Errorf("riscv64: object file emission not yet supported; use -emit-asm or -emit-mir")
	}

	return nil, fmt.Errorf("unsupported arch for code generation: %s", tri.arch)
}

// ── compileToASM ──────────────────────────────────────────────────────────────

// compileToASM runs isel + regalloc and returns human-readable assembly text
// for every function (-emit-asm). amd64 output is Intel-syntax; arm64 and
// riscv64 use GNU-as syntax, matching each arch's text printer.
func compileToASM(m *machine.Module, tri triple, opts codegenOptions) (string, error) {
	var sb strings.Builder

	switch tri.arch {
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
		return "", fmt.Errorf("unsupported arch for assembly emit: %s", tri.arch)
	}

	return sb.String(), nil
}