package compiler

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"

	"github.com/vertex-language/compiler/wasm"
	"github.com/vertex-language/vertex/parser"
)

const (
	dataBase  int32 = 0x0100
	frameBase int32 = 0x4000
)

type GenerateOptions struct {
	DCE       bool
	EntryName string
}

func (o GenerateOptions) entryName() string {
	if o.EntryName != "" {
		return o.EntryName
	}
	return "main"
}

// nativeInfo holds the resolved wasm import info for a native class method.
type nativeInfo struct {
	Module     string
	ImportName string
	Sig        *FuncSig
	FuncIdx    uint32
}

type funcInfo struct {
	Name    string
	Sig     *FuncSig
	FuncIdx uint32
	Ctx     parser.IFuncDeclContext
}

// CodeGen is the wasm code-generation driver.
type CodeGen struct {
	pkg        *Package
	globalTags *BuildTags
	scope      *Scope

	mod       *wasm.Module
	natives   map[string]*nativeInfo // native class methods (wasm imports)
	funcs     map[string]*funcInfo   // local Vertex functions
	funcSlice []*funcInfo
	nextFnIdx uint32

	reachable *DCEResult

	dataBuf     []byte
	strCache    map[string]int32
	frameCursor int32

	errors ErrorList
}

func newCodeGen(pkg *Package, scope *Scope, tags *BuildTags, reachable *DCEResult) *CodeGen {
	return &CodeGen{
		pkg:         pkg,
		globalTags:  tags,
		scope:       scope,
		mod:         wasm.NewModule(),
		natives:     make(map[string]*nativeInfo),
		funcs:       make(map[string]*funcInfo),
		strCache:    make(map[string]int32),
		frameCursor: frameBase,
		reachable:   reachable,
	}
}

// Generate is the top-level entry: two passes → wasm.Module.
func Generate(pkg *Package, defaultTags *BuildTags, opts GenerateOptions) (*wasm.Module, error) {
	// Pass 1: collect declarations into the global scope.
	chk := newChecker(pkg, defaultTags)
	scope := chk.Run()
	if err := chk.errors.err(); err != nil {
		return nil, err
	}

	// Optional DCE pass.
	var reachable *DCEResult
	if opts.DCE {
		reachable = ComputeReachable(pkg, opts.entryName())
	}

	cg := newCodeGen(pkg, scope, defaultTags, reachable)

	// Pass 2a: register native class imports first so their wasm function
	// indices are lower than local functions (wasm requires imports first).
	for _, sf := range pkg.Files {
		tags := defaultTags
		if len(sf.BuildTags.Tags) > 0 {
			tags = sf.BuildTags
		}
		for _, tld := range sf.Tree.AllTopLevelDecl() {
			cd := tld.ClassDecl()
			if cd != nil && cd.COLON() != nil {
				cg.collectDecl(tld, sf, tags)
			}
		}
	}

	// Pass 2b: register local function slots.
	for _, sf := range pkg.Files {
		tags := defaultTags
		if len(sf.BuildTags.Tags) > 0 {
			tags = sf.BuildTags
		}
		for _, tld := range sf.Tree.AllTopLevelDecl() {
			cd := tld.ClassDecl()
			if cd != nil && cd.COLON() != nil {
				continue // already handled in pass 2a
			}
			cg.collectDecl(tld, sf, tags)
		}
	}

	cg.mod.Memories.Add(wasm.MemoryType{Lim: wasm.Limits{Min: 2}})

	// Pass 2c: generate function bodies.
	for _, fi := range cg.funcSlice {
		cg.genFuncBody(fi)
	}

	cg.flushData()

	if err := cg.errors.err(); err != nil {
		return nil, err
	}
	return cg.mod, nil
}

// ── String / data helpers ─────────────────────────────────────────────────────

func (cg *CodeGen) internString(s string) int32 {
	if off, ok := cg.strCache[s]; ok {
		return off
	}
	off := dataBase + int32(len(cg.dataBuf))
	cg.dataBuf = append(cg.dataBuf, s...)
	cg.dataBuf = append(cg.dataBuf, 0)
	cg.strCache[s] = off
	return off
}

func (cg *CodeGen) flushData() {
	if len(cg.dataBuf) == 0 {
		return
	}
	cg.mod.Datas.Add(
		wasm.DataModeActive{MemIdx: 0, Offset: wasm.ConstI32(dataBase)},
		cg.dataBuf,
	)
}

// ── Frame helpers ─────────────────────────────────────────────────────────────

func (cg *CodeGen) allocFrame(size int) int32 {
	addr := cg.frameCursor
	cg.frameCursor += int32(size)
	if rem := cg.frameCursor % 8; rem != 0 {
		cg.frameCursor += 8 - rem
	}
	return addr
}

// ── Wasm type section ─────────────────────────────────────────────────────────

func (cg *CodeGen) addFuncType(sig *FuncSig) uint32 {
	return cg.mod.Types.AddFuncType(wasm.FuncType{
		Params:  ParamsToWasm(sig.Params),
		Results: RetToWasm(sig.Ret),
	})
}

// ── Misc ──────────────────────────────────────────────────────────────────────

func stripQuotes(s string) string {
	if unquoted, err := strconv.Unquote(s); err == nil {
		return unquoted
	}
	return strings.Trim(s, `"`)
}

func f32bits(v float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(v))
	return b
}

func f64bits(v float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(v))
	return b
}