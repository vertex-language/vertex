package token

// Ctx is the identity of a contextual keyword. It is zero (CtxNone) unless
// Kind is IDENT.
//
// This is how contextual keywords stay contextual (§3): the token remains an
// IDENT and stays usable as a binding name, which is what A.2 requires, while
// the parser tests tok.Ctx == CtxStruct in O(1) with no string compare.
// Nothing downstream should ever compare token text against a keyword string.
//
// Scope note: the grammar's StrictReservedWord list (implements, interface,
// let, package, private, protected, public, static) has no productions of its
// own — strict-mode restrictions are early errors here, not grammar. Those
// words are therefore *not* Kinds, and they live in this table so that
// ClassElementModifier (declare static abstract override readonly public
// protected private) is a single contiguous range test rather than a mix of
// Kind and Ctx checks. `let` in particular must not be unconditionally
// reserved: ExpressionStatement's lookahead ∉ { let [ } depends on it.
type Ctx uint8

const (
	CtxNone Ctx = iota

	// ClassElementModifier (E), in grammar order where it matters.
	// Accessibility (public/protected/private) is the contiguous tail.
	modifierBeg
	CtxDeclare
	CtxStatic
	CtxAbstract
	CtxOverride
	CtxReadonly
	accessibilityBeg
	CtxPublic
	CtxProtected
	CtxPrivate
	accessibilityEnd
	modifierEnd

	// PredefinedType (G.1), minus `void` and `null`, which are ReservedWords
	// and therefore Kinds.
	predefinedBeg
	CtxAny
	CtxUnknown
	CtxNever
	CtxObject
	CtxString
	CtxNumber
	CtxBigint
	CtxBoolean
	CtxSymbol
	CtxUndefined
	CtxIntrinsic
	predefinedEnd

	// Everything else, alphabetical.
	CtxAccessor
	CtxAs
	CtxAsserts
	CtxAsync
	CtxConstructor
	CtxDefer
	CtxFrom
	CtxGet
	CtxGlobal
	CtxGraph
	CtxImplements
	CtxInfer
	CtxInterface
	CtxIs
	CtxKernel
	CtxKeyof
	CtxLet
	CtxModule
	CtxMutating
	CtxNamespace
	CtxOf
	CtxOut
	CtxPackage
	CtxRequire
	CtxSatisfies
	CtxSet
	CtxSource
	CtxStruct
	CtxType
	CtxUnique
	CtxUsing
)

var ctxNames = [...]string{
	CtxDeclare:   "declare",
	CtxStatic:    "static",
	CtxAbstract:  "abstract",
	CtxOverride:  "override",
	CtxReadonly:  "readonly",
	CtxPublic:    "public",
	CtxProtected: "protected",
	CtxPrivate:   "private",

	CtxAny:       "any",
	CtxUnknown:   "unknown",
	CtxNever:     "never",
	CtxObject:    "object",
	CtxString:    "string",
	CtxNumber:    "number",
	CtxBigint:    "bigint",
	CtxBoolean:   "boolean",
	CtxSymbol:    "symbol",
	CtxUndefined: "undefined",
	CtxIntrinsic: "intrinsic",

	CtxAccessor:    "accessor",
	CtxAs:          "as",
	CtxAsserts:     "asserts",
	CtxAsync:       "async",
	CtxConstructor: "constructor",
	CtxDefer:       "defer",
	CtxFrom:        "from",
	CtxGet:         "get",
	CtxGlobal:      "global",
	CtxGraph:       "graph",
	CtxImplements:  "implements",
	CtxInfer:       "infer",
	CtxInterface:   "interface",
	CtxIs:          "is",
	CtxKernel:      "kernel",
	CtxKeyof:       "keyof",
	CtxLet:         "let",
	CtxModule:      "module",
	CtxMutating:    "mutating",
	CtxNamespace:   "namespace",
	CtxOf:          "of",
	CtxOut:         "out",
	CtxPackage:     "package",
	CtxRequire:     "require",
	CtxSatisfies:   "satisfies",
	CtxSet:         "set",
	CtxSource:      "source",
	CtxStruct:      "struct",
	CtxType:        "type",
	CtxUnique:      "unique",
	CtxUsing:       "using",
}

// String returns the source spelling, or "" for CtxNone.
func (c Ctx) String() string {
	if int(c) < len(ctxNames) {
		return ctxNames[c]
	}
	return "Ctx(" + itoa(int(c)) + ")"
}

// IsClassElementModifier reports whether c is a ClassElementModifier (E).
func (c Ctx) IsClassElementModifier() bool { return modifierBeg < c && c < modifierEnd }

// IsAccessibilityModifier reports whether c is public, protected, or private.
func (c Ctx) IsAccessibilityModifier() bool { return accessibilityBeg < c && c < accessibilityEnd }

// IsPredefinedType reports whether c is a PredefinedType name (G.1). Note that
// `void` and `null` are ReservedWords and answer to Kind, not Ctx, so a type
// parser must test both.
func (c Ctx) IsPredefinedType() bool { return predefinedBeg < c && c < predefinedEnd }

// IsAccelModifier reports whether c is `kernel` or `graph` (D.4).
func (c Ctx) IsAccelModifier() bool { return c == CtxKernel || c == CtxGraph }

type identInfo struct {
	kind Kind
	ctx  Ctx
}

// identTable is the single lookup for every IdentifierName. It is built once
// at init from kindNames and ctxNames, so a word cannot be added to one table
// and forgotten in the other.
//
// A map today; the shape is a perfect hash (§3) and swapping the
// implementation touches nothing outside this file.
var identTable map[string]identInfo

func init() {
	identTable = make(map[string]identInfo, 96)
	for k := reservedBeg + 1; k < reservedEnd; k++ {
		identTable[kindNames[k]] = identInfo{kind: k}
	}
	for c := CtxNone + 1; int(c) < len(ctxNames); c++ {
		name := ctxNames[c]
		if name == "" {
			continue // a sentinel
		}
		if prev, dup := identTable[name]; dup {
			panic("token: duplicate identifier table entry " + name +
				" (" + prev.kind.String() + ")")
		}
		identTable[name] = identInfo{kind: IDENT, ctx: c}
	}
}

// LookupIdent classifies an IdentifierName. It returns a reserved Kind with
// CtxNone, or IDENT with a Ctx, or IDENT with CtxNone for an ordinary name.
//
// The caller must NOT call this for an identifier whose source spelling
// contained an escape: `\u0069f` is an IDENT named "if", not the IF keyword.
// That is what Flags.HasEscape marks, and skipping this call is the scanner's
// responsibility (see token.go).
//
// name is not retained, and m[string(name)] does not allocate.
func LookupIdent(name []byte) (Kind, Ctx) {
	if info, ok := identTable[string(name)]; ok {
		return info.kind, info.ctx
	}
	return IDENT, CtxNone
}