package compiler

import (
	"fmt"
	"strings"

	"github.com/vertex-language/vertex/parser"
)

// resolveType maps an AST type node to a VType.
func (c *Compiler) resolveType(ctx parser.ITypeContext) (VType, error) {
	if ctx == nil {
		return &PrimitiveType{Kind: KindVoid}, nil
	}
	switch t := ctx.(type) {
	case *parser.NamedTypeContext:
		name := typeIdentName(t.TypeIdentifier())
		return c.resolveTypeByName(name)

	case *parser.OptTypeContext:
		inner, err := c.resolveType(t.Type_())
		if err != nil {
			return nil, err
		}
		return &OptionalType{Inner: inner}, nil

	case *parser.ArrTypeContext:
		// Simplified: [T] is an i32 pointer into linear memory.
		return &PrimitiveType{Kind: KindInt}, nil

	case *parser.DictTypeContext:
		return &PrimitiveType{Kind: KindInt}, nil

	case *parser.TupTypeContext:
		return &PrimitiveType{Kind: KindInt}, nil

	case *parser.FuncTypeContext:
		return &PrimitiveType{Kind: KindInt}, nil

	case *parser.SelfType_Context:
		return &PrimitiveType{Kind: KindInt}, nil

	case *parser.ExistTypeContext:
		return &PrimitiveType{Kind: KindInt}, nil

	case *parser.OpaqueType_Context:
		return c.resolveType(t.OpaqueType().Type_())
	}

	// Fallback: strip qualifiers and look up by name.
	text := strings.TrimRight(ctx.GetText(), "?")
	return c.resolveTypeByName(text)
}

func (c *Compiler) resolveTypeByName(name string) (VType, error) {
	if bt, ok := builtinTypes[name]; ok {
		return bt, nil
	}
	if st, ok := c.structMap[name]; ok {
		return st, nil
	}
	return nil, fmt.Errorf("unknown type %q", name)
}

// typeIdentName collapses a dotted TypeIdentifier to a string.
func typeIdentName(ctx parser.ITypeIdentifierContext) string {
	if ctx == nil {
		return ""
	}
	name := ctx.Identifier().GetText()
	if inner := ctx.TypeIdentifier(); inner != nil {
		return name + "." + typeIdentName(inner)
	}
	return name
}

// stripQuotes removes surrounding double-quotes from a string literal token.
func stripQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}