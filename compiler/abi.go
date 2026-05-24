// abi.go
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
		return "ptr" // const char* on the C side
	case KindArray, KindStruct, KindClass:
		return "ptr"
	default:
		return "i32"
	}
}

// needsABISuffix reports whether the @-suffix must be appended.
// Only ptr and hptr tokens require the suffix; plain integer/float params do not.
func needsABISuffix(params []Type, ret Type) bool {
	for _, p := range params {
		tok := abiToken(p)
		if tok == "ptr" || tok == "hptr" {
			return true
		}
	}
	return ret != nil && abiToken(ret) == "hptr"
}

// BuildImportName constructs the decorated wasm import name.
//
// Examples matching the ABI reference:
//
//	write(fd i32, buf ptr, count i32) i32  → "write@i32.ptr.i32"
//	fopen(path ptr, mode ptr) hptr         → "fopen@ptr.ptr:hptr"
//	getpid() i32                           → "getpid"   (no suffix needed)
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

// ExportName builds the export name for a qualified (non-CPU) function.
//
// Concurrency exports:
//
//	"worker@thread:ptr.i32"
//	"handler@async"
//	"task@process:i32"
//
// GPU kernel exports:
//
//	"vectorAdd@cuda:ptr.ptr.i32"
//	"tileConv@msl:ptr.ptr.i32"
//	"histogram@vulkan:ptr.i32"
func ExportName(name string, qual FuncQual, params []Type) string {
	suffix := qualSuffix(qual)
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

// qualSuffix maps a FuncQual to its ABI export suffix string.
func qualSuffix(q FuncQual) string {
	switch q {
	case QualAsync:
		return "@async"
	case QualThread:
		return "@thread"
	case QualProcess:
		return "@process"
	case QualCUDA:
		return "@cuda"
	case QualVulkan:
		return "@vulkan"
	case QualMSL:
		return "@msl"
	default:
		return ""
	}
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

// ParamsToWasm converts Vertex param types to a wasm.ValType slice.
func ParamsToWasm(types []Type) []wasm.ValType {
	var out []wasm.ValType
	for _, t := range types {
		if t.Kind() != KindVoid {
			out = append(out, ToWasmVal(t.Wasm()))
		}
	}
	return out
}

// RetToWasm converts a Vertex return type to a wasm results slice.
func RetToWasm(t Type) []wasm.ValType {
	if t == nil || t.Kind() == KindVoid {
		return nil
	}
	return []wasm.ValType{ToWasmVal(t.Wasm())}
}