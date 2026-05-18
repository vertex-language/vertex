package compiler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/wasm-compiler/wasm"
)

// Options configures the Vertex compiler.
type Options struct {
	Verbose bool
}

// Compiler translates Vertex source into a WebAssembly module.
type Compiler struct {
	opts         Options
	mod          *wasm.Module
	funcMap      map[string]*FuncInfo
	funcOrder    []string // insertion order for codegen
	typeDedup    map[string]uint32
	stringPool   map[string]uint32
	nextDataAddr uint32
	structMap    map[string]*StructType
}

// NewCompiler creates a ready-to-use Compiler.
func NewCompiler(opts Options) (*Compiler, error) {
	c := &Compiler{
		opts:         opts,
		mod:          wasm.NewModule(),
		funcMap:      make(map[string]*FuncInfo),
		typeDedup:    make(map[string]uint32),
		stringPool:   make(map[string]uint32),
		nextDataAddr: 1024,
		structMap:    make(map[string]*StructType),
	}
	c.setupModule()
	return c, nil
}

func (c *Compiler) setupModule() {
	// 16 pages = 1 MiB of linear memory.
	c.mod.Memories.Add(wasm.MemoryType{Lim: wasm.Limits{Min: 16}})
	// Global 0 = stack pointer (grows downward from 64 KiB).
	c.mod.Globals.Add(
		wasm.GlobalType{Val: wasm.I32, Mutable: true},
		wasm.ConstI32(65536),
	)
}

// CompileToModule parses and compiles source into a *wasm.Module.
func (c *Compiler) CompileToModule(source, filename string) (*wasm.Module, error) {
	tree, err := c.parse(source, filename)
	if err != nil {
		return nil, err
	}
	if err := c.collectDeclarations(tree); err != nil {
		return nil, fmt.Errorf("declaration pass: %w", err)
	}
	if err := c.generateCode(tree); err != nil {
		return nil, fmt.Errorf("codegen: %w", err)
	}
	return c.mod, nil
}

// BuildPointerMap returns GC root information for native backends.
func (c *Compiler) BuildPointerMap() map[uint32][]uint32 {
	return make(map[uint32][]uint32)
}

// ── Internal helpers ──────────────────────────────────────────────────────────

func (c *Compiler) parse(source, filename string) (parser.ITopLevelContext, error) {
	input := antlr.NewInputStream(source)
	lexer := parser.NewVertexLexer(input)
	lexer.RemoveErrorListeners()
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewVertexParser(stream)
	p.RemoveErrorListeners()

	el := &errorListener{filename: filename}
	p.AddErrorListener(el)
	tree := p.TopLevel()

	if len(el.errors) > 0 {
		return nil, fmt.Errorf("syntax errors:\n  %s", strings.Join(el.errors, "\n  "))
	}
	return tree, nil
}

// internFuncType deduplicates WASM function type entries.
func (c *Compiler) internFuncType(ft wasm.FuncType) uint32 {
	key := fmt.Sprintf("%v|%v", ft.Params, ft.Results)
	if idx, ok := c.typeDedup[key]; ok {
		return idx
	}
	idx := c.mod.Types.AddFuncType(ft)
	c.typeDedup[key] = idx
	return idx
}

// internString stores a NUL-terminated string in the data section and returns
// its linear-memory offset.
func (c *Compiler) internString(raw string) uint32 {
	text := raw
	
	// Strip the literal quotes parsed by ANTLR
	if len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"' {
		text = text[1 : len(text)-1]
	}
	
	// Decode escape sequences (e.g. \n into a real newline byte)
	if unescaped, err := strconv.Unquote(`"` + text + `"`); err == nil {
		text = unescaped
	}

	if off, ok := c.stringPool[text]; ok {
		return off
	}
	
	off := c.nextDataAddr
	c.stringPool[text] = off
	data := append([]byte(text), 0) // NUL-terminated
	c.mod.Datas.Add(
		wasm.DataModeActive{MemIdx: 0, Offset: wasm.ConstI32(int32(off))},
		data,
	)
	c.nextDataAddr += uint32(len(data))
	return off
}

// ── Error listener ────────────────────────────────────────────────────────────

type errorListener struct {
	*antlr.DefaultErrorListener
	filename string
	errors   []string
}

func (e *errorListener) SyntaxError(
	_ antlr.Recognizer, _ interface{},
	line, col int, msg string, _ antlr.RecognitionException,
) {
	e.errors = append(e.errors,
		fmt.Sprintf("%s:%d:%d: %s", e.filename, line, col, msg))
}