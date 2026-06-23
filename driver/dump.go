package driver

import (
	"fmt"
	"io"
	"strings"

	"github.com/vertex-language/ir/vertex/ast"
	virlower "github.com/vertex-language/ir/vertex/lower/vir"
	mirlower "github.com/vertex-language/ir/vertex/lower/mir"

	virtext "github.com/vertex-language/ir/vertex/encoding/text"

	"github.com/vertex-language/ir/machine"
	mirtext "github.com/vertex-language/ir/machine/encoding/text/mir"
)

// dumpAll runs every pipeline stage and writes an annotated multi-section file
// to cfg.output (use "-" for stdout). Stages are separated by banner comments:
//
//	; ════ Stage N: <name> ════════════════════════════════════════════
//
// Soft VIR errors (non-nil virErr with non-nil virMod) are noted inline and
// the dump continues so the programmer sees as much of the pipeline as
// possible. Hard failures stop the dump at the failing stage and return 1.
func dumpAll(cfg config, stderr io.Writer) int {
	tri, err := parseTriple(cfg.target)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 2
	}

	var sb strings.Builder
	exitCode := 0

	banner := func(n int, name string) {
		const line = "════════════════════════════════════════════════════════════"
		fmt.Fprintf(&sb, "; ════ Stage %d: %-36s%s\n\n", n, name, line)
	}

	// ── Stage 1: Source (.vs) ─────────────────────────────────────────────────
	banner(1, "Source (.vs)")
	if err := appendSourceFiles(&sb, cfg.input); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	sb.WriteString("\n")

	// ── Stage 2: Vertex IR ────────────────────────────────────────────────────
	banner(2, "Vertex IR (.vir)")

	pkg, err := parseInput(cfg.input)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	virMod, virErr := virlower.NewLower(pkg, nil, cfg.target)
	if virErr != nil {
		fmt.Fprintf(&sb, "; WARNING: %v\n\n", virErr)
		exitCode = 1
	}
	if virMod == nil {
		sb.WriteString("; (no VIR produced — pipeline stopped)\n")
		return writeAndReturn(cfg.output, sb.String(), stderr, 1)
	}
	virMod.SetTarget(tri.virTargetString())
	sb.WriteString(virtext.Format(virMod))
	sb.WriteString("\n\n")

	// ── Stage 3: Machine IR ───────────────────────────────────────────────────
	banner(3, "Machine IR (.mir)")

	mirMod, err := mirlower.NewLower(virMod)
	if err != nil {
		fmt.Fprintf(&sb, "; ERROR: MIR lowering: %v\n", err)
		sb.WriteString("; (pipeline stopped)\n")
		return writeAndReturn(cfg.output, sb.String(), stderr, 1)
	}
	if err := machine.Verify(mirMod); err != nil {
		fmt.Fprintf(&sb, "; ERROR: MIR verification: %v\n", err)
		sb.WriteString("; (pipeline stopped)\n")
		return writeAndReturn(cfg.output, sb.String(), stderr, 1)
	}
	sb.WriteString(mirtext.PrintModule(mirMod))
	sb.WriteString("\n\n")

	// ── Stage 4: Assembly ─────────────────────────────────────────────────────
	banner(4, "Assembly (.s)")

	opts := codegenOptions{optLevel: cfg.optLevel, debugInfo: cfg.debugInfo}
	asmText, err := compileToASM(mirMod, tri, opts)
	if err != nil {
		fmt.Fprintf(&sb, "; ERROR: code generation: %v\n", err)
		sb.WriteString("; (pipeline stopped)\n")
		return writeAndReturn(cfg.output, sb.String(), stderr, 1)
	}
	sb.WriteString(asmText)
	sb.WriteString("\n")

	return writeAndReturn(cfg.output, sb.String(), stderr, exitCode)
}

// dumpPackage runs every pipeline stage against an already-parsed synthetic
// package and writes the result to path. Used by the test runner to capture
// a full pipeline dump for failed tests.
func dumpPackage(pkg *ast.Package, cfg config, path string, stderr io.Writer) {
	tri, err := parseTriple(cfg.target)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: dump: %v\n", err)
		return
	}

	var sb strings.Builder

	banner := func(n int, name string) {
		const line = "════════════════════════════════════════════════════════════"
		fmt.Fprintf(&sb, "; ════ Stage %d: %-36s%s\n\n", n, name, line)
	}

	// ── Stage 2: Vertex IR ────────────────────────────────────────────────────
	banner(2, "Vertex IR (.vir)")
	virMod, virErr := virlower.NewLower(pkg, nil, cfg.target)
	if virErr != nil {
		fmt.Fprintf(&sb, "; WARNING: %v\n\n", virErr)
	}
	if virMod == nil {
		sb.WriteString("; (no VIR produced — pipeline stopped)\n")
		_ = writeOutput(path, []byte(sb.String()))
		return
	}
	virMod.SetTarget(tri.virTargetString())
	sb.WriteString(virtext.Format(virMod))
	sb.WriteString("\n\n")

	// ── Stage 3: Machine IR ───────────────────────────────────────────────────
	banner(3, "Machine IR (.mir)")
	mirMod, err := mirlower.NewLower(virMod)
	if err != nil {
		fmt.Fprintf(&sb, "; ERROR: MIR lowering: %v\n", err)
		sb.WriteString("; (pipeline stopped)\n")
		_ = writeOutput(path, []byte(sb.String()))
		return
	}
	if err := machine.Verify(mirMod); err != nil {
		fmt.Fprintf(&sb, "; ERROR: MIR verification: %v\n", err)
		sb.WriteString("; (pipeline stopped)\n")
		_ = writeOutput(path, []byte(sb.String()))
		return
	}
	sb.WriteString(mirtext.PrintModule(mirMod))
	sb.WriteString("\n\n")

	// ── Stage 4: Assembly ─────────────────────────────────────────────────────
	banner(4, "Assembly (.s)")
	opts := codegenOptions{optLevel: cfg.optLevel, debugInfo: cfg.debugInfo}
	asmText, err := compileToASM(mirMod, tri, opts)
	if err != nil {
		fmt.Fprintf(&sb, "; ERROR: code generation: %v\n", err)
		sb.WriteString("; (pipeline stopped)\n")
		_ = writeOutput(path, []byte(sb.String()))
		return
	}
	sb.WriteString(asmText)
	sb.WriteString("\n")

	_ = writeOutput(path, []byte(sb.String()))
}

// writeAndReturn writes the accumulated dump to path and returns code.
// It is the final action of dumpAll regardless of which stage stopped last.
func writeAndReturn(path, content string, stderr io.Writer, code int) int {
	if err := writeOutput(path, []byte(content)); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	return code
}