package nativelibs

import (
	"sort"
	"strings"

	"github.com/vertex-language/ir/vertex"

	"github.com/vertex-language/vertex/target"
)

// ExtractDynLibs walks the VIR import section and returns the unique set
// of shared library names required for tri.OS.
//
// Import Module strings follow the "<platform>:<libname>" convention
// defined in §13.1 of the VIR spec (e.g. "linux:libc.so.6",
// "darwin:libSystem.B.dylib"). Only entries whose platform component
// equals tri.OS are collected.
func ExtractDynLibs(m *vertex.Module, tri target.Triple) []string {
	seen := make(map[string]bool)
	var libs []string
	for _, imp := range m.Imports.Imports {
		platform, lib, ok := splitImportModule(imp.Module)
		if !ok || platform != tri.OS {
			continue
		}
		lib = normalizeDarwinLib(lib)
		if !seen[lib] {
			seen[lib] = true
			libs = append(libs, lib)
		}
	}
	sort.Slice(libs, func(i, j int) bool {
		return libSortKey(libs[i]) < libSortKey(libs[j])
	})
	return libs
}

// ExtractLibFuncSymbols groups every imported native function by the
// shared library it comes from, for triOS. Used to give darwin's
// cached-dylib path (a library resolved by name only, with no on-disk
// bytes — see isDarwinCacheOnly in resolve.go) the symbol list it needs to
// synthesize lazy-bind stubs.
func ExtractLibFuncSymbols(m *vertex.Module, triOS string) map[string][]string {
	result := make(map[string][]string)
	for _, imp := range m.Imports.Imports {
		platform, lib, ok := splitImportModule(imp.Module)
		if !ok || platform != triOS {
			continue
		}
		lib = normalizeDarwinLib(lib)
		if _, ok := imp.Desc.(vertex.FuncImport); ok {
			result[lib] = append(result[lib], imp.Name)
		}
	}
	return result
}

func normalizeDarwinLib(lib string) string {
	switch lib {
	case "libSystem.dylib":
		return "libSystem.B.dylib"
	case "libobjc.dylib":
		return "libobjc.A.dylib"
	}
	return lib
}

// libSortKey returns a sort priority: lower = earlier LC_LOAD_DYLIB
// position. libSystem must be first (ordinal 1), then other system
// dylibs, then frameworks, then libobjc last.
func libSortKey(lib string) int {
	switch lib {
	case "libSystem.B.dylib", "libSystem.dylib":
		return 0
	case "libobjc.dylib", "libobjc.A.dylib":
		return 3
	}
	if strings.Contains(lib, ".framework/") {
		return 2
	}
	return 1
}

func splitImportModule(module string) (platform, lib string, ok bool) {
	i := strings.IndexByte(module, ':')
	if i < 0 {
		return "", "", false
	}
	return module[:i], module[i+1:], true
}