package compiler

import (
	"strings"

	"github.com/vertex-language/compiler/wasm"
)

// abiToken returns the ABI type token for a Vertex type.
func abiToken(t Type) string {
	if t == nil {
		return "i32"
	}
	switch t.Kind() {
	case KindInt64, KindUint64:
		return "i64"
	case KindFloat:
		return "f32"
	case KindDouble:
		return "f64"
	case KindPointer:
		pt := t.(*PointerType)
		if pt.Elem.Kind() == KindOpaque {
			return "hptr"
		}
		return "ptr"
	case KindString:
		// strings passed to C as const char* → ptr
		return "ptr"
	case KindArray, KindStruct, KindClass:
		return "ptr"
	default:
		return "i32"
	}
}

// needsABISuffix returns true when at least one token is ptr or hptr,
// or the return is hptr.
func needsABISuffix(params []Type, ret Type) bool {
	for _, p := range params {
		tok := abiToken(p)
		if tok == "ptr" || tok == "hptr" {
			return true
		}
	}
	if ret != nil && abiToken(ret) == "hptr" {
		return true
	}
	return false
}

// BuildImportName builds the decorated import name, e.g. "write@i32.ptr.i32"
// or "fopen@ptr.ptr:hptr".
func BuildImportName(name string, params []Type, ret Type) string {
	if !needsABISuffix(params, ret) {
		return name
	}
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteByte('@')
	for i, p := range params {
		if i > 0 {
			sb.WriteByte('.')
		}
		sb.WriteString(abiToken(p))
	}
	if ret != nil && ret.Kind() != KindVoid && abiToken(ret) == "hptr" {
		sb.WriteByte(':')
		sb.WriteString("hptr")
	}
	return sb.String()
}

// ExportName builds the export name for a qualified function.
// e.g. "vectorAdd@cuda:ptr.ptr.i32"
func ExportName(name string, qual FuncQual, params []Type) string {
	suffix := ""
	switch qual {
	case QualAsync:
		suffix = "@async"
	case QualThread:
		suffix = "@thread"
	case QualProcess:
		suffix = "@process"
	case QualGPU:
		suffix = "@cuda" // default GPU target
	}
	if suffix == "" {
		return name
	}
	if len(params) == 0 {
		return name + suffix
	}
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString(suffix)
	sb.WriteByte(':')
	for i, p := range params {
		if i > 0 {
			sb.WriteByte('.')
		}
		sb.WriteString(abiToken(p))
	}
	return sb.String()
}

// ToWasmVal converts a WasmKind to wasm.ValType.
func ToWasmVal(w WasmKind) wasm.ValType {
	switch w {
	case WasmI64:
		return wasm.I64
	case WasmF32:
		return wasm.F32
	case WasmF64:
		return wasm.F64
	default:
		return wasm.I32
	}
}

// ParamsToWasm converts Vertex param types to wasm.ValType slice.
func ParamsToWasm(types []Type) []wasm.ValType {
	var out []wasm.ValType
	for _, t := range types {
		if t.Kind() != KindVoid {
			out = append(out, ToWasmVal(t.Wasm()))
		}
	}
	return out
}

// RetToWasm converts a Vertex return type to wasm results slice.
func RetToWasm(t Type) []wasm.ValType {
	if t == nil || t.Kind() == KindVoid {
		return nil
	}
	return []wasm.ValType{ToWasmVal(t.Wasm())}
}