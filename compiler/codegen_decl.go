package compiler

import (
	"github.com/vertex-language/compiler/wasm"
	"github.com/vertex-language/vertex/parser"
)

// collectDecl is called for each top-level declaration in pass 2a/2b.
func (cg *CodeGen) collectDecl(
	ctx parser.ITopLevelDeclContext,
	sf *SourceFile,
	tags *BuildTags,
) {
	switch {
	case ctx.ClassDecl() != nil && ctx.ClassDecl().COLON() != nil:
		cg.collectNativeClass(ctx.ClassDecl(), sf, tags)
	case ctx.FuncDecl() != nil:
		cg.collectFuncSlot(ctx.FuncDecl())
	// Structs / regular classes / enums handled by the checker; no wasm code needed.
	}
}

// ── Native class ──────────────────────────────────────────────────────────────

// collectNativeClass registers every method of a native class as a wasm import.
// The emission strategy is determined by the file's native import prefix.
func (cg *CodeGen) collectNativeClass(
	ctx parser.IClassDeclContext,
	sf *SourceFile,
	tags *BuildTags,
) {
	// Determine wasm module name from the file's native import.
	// Fall back to the parent namespace if no import is found.
	module := ctx.AllID()[1].GetText() // parent namespace as fallback
	if sf.NativeImport != nil {
		module = sf.NativeImport.WasmModule(tags)
	}

	resolve := func(tc parser.ITypeContext) Type {
		return ResolveType(tc, cg.scope)
	}

	entries := CollectNativeFuncs(ctx, module, cg.scope, resolve)

	for name, entry := range entries {
		// DCE: skip methods that are never called.
		if !cg.reachable.IsReachable(name) {
			continue
		}
		if _, already := cg.natives[name]; already {
			continue
		}

		typeIdx := cg.addFuncType(entry.Sig)
		cg.mod.Imports.AddFunc(entry.WasmModule, entry.ImportName, typeIdx)

		idx := cg.nextFnIdx
		cg.nextFnIdx++

		cg.natives[name] = &nativeInfo{
			Module:     entry.WasmModule,
			ImportName: entry.ImportName,
			Sig:        entry.Sig,
			FuncIdx:    idx,
		}

		// Patch the FuncIdx into the global scope symbol.
		if sym := cg.scope.Lookup(name); sym != nil {
			sym.FuncIdx = idx
		}
	}
}

// ── Function slot ─────────────────────────────────────────────────────────────

func (cg *CodeGen) collectFuncSlot(ctx parser.IFuncDeclContext) {
	name := ctx.ID().GetText()

	if !cg.reachable.IsReachable(name) {
		return
	}
	if _, already := cg.funcs[name]; already {
		return
	}

	sym := cg.scope.Lookup(name)
	var sig *FuncSig
	if sym != nil {
		if ft, ok := sym.Type.(*FuncType); ok {
			sig = ft.Sig
		}
	}
	if sig == nil {
		chk := newChecker(cg.pkg, cg.globalTags)
		chk.scope = cg.scope
		sig = chk.funcDeclSig(ctx)
	}

	typeIdx := cg.addFuncType(sig)

	importCount := cg.mod.Imports.NumFuncs()
	localIdx := cg.mod.Functions.Add(typeIdx)
	globalIdx := importCount + localIdx

	exportName := ""
	switch name {
	case "main":
		exportName = "main"
	default:
		if sig.Qual != QualNone {
			exportName = ExportName(name, sig.Qual, sig.Params)
		}
	}
	if exportName != "" {
		cg.mod.Exports.Add(exportName, wasm.ExportFunc, globalIdx)
	}

	fi := &funcInfo{Name: name, Sig: sig, FuncIdx: globalIdx, Ctx: ctx}
	cg.funcs[name] = fi
	cg.funcSlice = append(cg.funcSlice, fi)
	cg.nextFnIdx = globalIdx + 1

	if sym != nil {
		sym.FuncIdx = globalIdx
	}
}

// ── Function body ─────────────────────────────────────────────────────────────

func (cg *CodeGen) genFuncBody(fi *funcInfo) {
	ctx := fi.Ctx
	body := wasm.NewFunctionBody()

	fnScope := NewScope(cg.scope)
	fg := &funcGen{
		cg:          cg,
		sig:         fi.Sig,
		scope:       fnScope,
		body:        body,
		localCount:  0,
		frameAllocs: make(map[string]int32),
	}

	// ← 2.0: Process the receiver block first, treating it as parameter 0.
	if rec := ctx.Receiver(); rec != nil {
		pname := rec.ID().GetText()
		ptype := ResolveType(rec.Type_(), cg.scope)
		localIdx := fg.newLocal(ptype)
		fg.paramCount++
		fnScope.Define(&Symbol{
			Name:     pname,
			Kind:     SymParam,
			Type:     ptype,
			LocalIdx: localIdx,
			Mutable:  rec.MUT() != nil,
		})
	}

	if pl := ctx.ParamList(); pl != nil {
		for _, p := range pl.AllParam() {
			pname := p.ID().GetText()
			ptype := ResolveType(p.Type_(), cg.scope)
			localIdx := fg.newLocal(ptype)
			fg.paramCount++
			fnScope.Define(&Symbol{
				Name:     pname,
				Kind:     SymParam,
				Type:     ptype,
				LocalIdx: localIdx,
				Mutable:  p.MUT() != nil,
			})
		}
	}

	if blk := ctx.Block(); blk != nil {
		fg.preScanBlock(blk)
	}
	if blk := ctx.Block(); blk != nil {
		fg.genBlock(blk)
	}

	fg.declareLocals(body)
	body.End()
	cg.mod.Codes.Add(body)
}