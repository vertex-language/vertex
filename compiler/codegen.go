package compiler

import (
	"fmt"

	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/wasm-compiler/wasm"
)

// FuncCompiler generates WASM for a single Vertex function body.
type FuncCompiler struct {
	c      *Compiler
	info   *FuncInfo
	body   *wasm.FunctionBody
	scope  *Scope
	locals []wasm.ValType // extra locals beyond params
}

// generateCode is Pass 2: emit WASM for every local function.
func (c *Compiler) generateCode(tree parser.ITopLevelContext) error {
	for _, name := range c.funcOrder {
		info := c.funcMap[name]
		if info.IsImport {
			continue
		}
		fd, ok := info.ASTNode.(*parser.FunctionDeclarationContext)
		if !ok || fd == nil {
			continue
		}
		if err := c.compileFunctionDecl(fd, info); err != nil {
			return fmt.Errorf("func %s: %w", name, err)
		}
	}
	return nil
}

func (c *Compiler) compileFunctionDecl(
	fd *parser.FunctionDeclarationContext,
	info *FuncInfo,
) error {
	body := wasm.NewFunctionBody()
	fc := &FuncCompiler{
		c:    c,
		info: info,
		body: body,
		scope: newScope(nil),
	}

	// Bind parameters as locals (WASM params are the first locals).
	for i, vt := range info.Params {
		wasmIdx := uint32(i)
		name := "_"
		if i < len(info.ParamNames) {
			name = info.ParamNames[i]
		}
		if name == "_" {
			continue
		}
		if wt, ok := vt.WasmType(); ok {
			_ = wt
		}
		_ = fc.scope.define(name, Symbol{
			Name:    name,
			Type:    vt,
			WasmIdx: wasmIdx,
			IsMut:   false,
		})
	}

	// Compile body if present.
	if fb := fd.FunctionBody(); fb != nil {
		cb := fb.CodeBlock()
		if err := fc.compileCodeBlock(cb); err != nil {
			return err
		}
	}

	// Append trailing `end` if needed.
	body.End()

	// Declare additional locals in the body.
	for _, lt := range fc.locals {
		body.AddLocals(1, lt)
	}

	c.mod.Codes.Add(body)

	if c.opts.Verbose {
		fmt.Printf("[codegen] func %s: %d params, %d extra locals\n",
			info.Name, len(info.Params), len(fc.locals))
	}
	return nil
}

// ── Local allocation ──────────────────────────────────────────────────────────

// allocLocal appends a new WASM local and returns its index.
func (fc *FuncCompiler) allocLocal(vt wasm.ValType) uint32 {
	idx := uint32(len(fc.info.Params)) + uint32(len(fc.locals))
	fc.locals = append(fc.locals, vt)
	return idx
}

// ── Scope helpers ─────────────────────────────────────────────────────────────

func (fc *FuncCompiler) pushScope() { fc.scope = newScope(fc.scope) }
func (fc *FuncCompiler) popScope()  { fc.scope = fc.scope.parent }

// ── Code block ────────────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileCodeBlock(cb parser.ICodeBlockContext) error {
	if stmts := cb.Statements(); stmts != nil {
		for _, stmt := range stmts.AllStatement() {
			if err := fc.compileStatement(stmt); err != nil {
				return err
			}
		}
	}
	return nil
}