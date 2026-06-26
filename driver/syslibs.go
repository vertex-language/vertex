package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sort"
	"github.com/vertex-language/ir/vertex"
)

type resolvedLib struct {
	name  string
	bytes []byte
}

type crtObjects struct {
	crt1 []byte // _start → __libc_start_main → main
	crti []byte // .init section prologue
	crtn []byte // .init section epilogue
}

// extractDynLibs walks the VIR import section and returns the unique set of
// shared library names required for tri.os.
//
// Import Module strings follow the "<platform>:<libname>" convention defined
// in §13.1 of the VIR spec (e.g. "linux:libc.so.6", "darwin:libSystem.B.dylib").
// Only entries whose platform component equals tri.os are collected; imports
// for other platforms present in the same module are silently skipped.
func extractDynLibs(m *vertex.Module, tri triple) []string {
	seen := make(map[string]bool)
	var libs []string
	for _, imp := range m.Imports.Imports {
		platform, lib, ok := splitImportModule(imp.Module)
		if !ok || platform != tri.os {
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

// normalizeDarwinLib maps known aliases to their canonical dyld name so
// duplicate LC_LOAD_DYLIB entries are never emitted for the same library.
func normalizeDarwinLib(lib string) string {
	switch lib {
	case "libSystem.dylib":
		return "libSystem.B.dylib"
	case "libobjc.dylib":
		return "libobjc.A.dylib"
	}
	return lib
}

// libSortKey returns a sort priority: lower = earlier LC_LOAD_DYLIB position.
// libSystem must be first (ordinal 1), then other system dylibs, then frameworks.
func libSortKey(lib string) int {
    switch lib {
    case "libSystem.B.dylib", "libSystem.dylib":
        return 0 // ordinal 1 — always first
    case "libobjc.dylib", "libobjc.A.dylib":
        return 3 // always last — after all frameworks
    }
    if strings.Contains(lib, ".framework/") {
        return 2 // frameworks in the middle
    }
    return 1 // other system dylibs (libm, libpthread etc)
}

// splitImportModule splits a VIR import module string of the form
// "<platform>:<libname>" into its two components.
// Returns ok=false for any string that contains no colon.
func splitImportModule(module string) (platform, lib string, ok bool) {
	i := strings.IndexByte(module, ':')
	if i < 0 {
		return "", "", false
	}
	return module[:i], module[i+1:], true
}

// resolveLibs locates each named shared library on the filesystem, reads it,
// and returns a slice of resolvedLib ready for the linker.
//
// If sysroot is non-empty every search directory is prefixed with it, which
// supports cross-compilation.
func resolveLibs(names []string, tri triple, sysroot string) ([]resolvedLib, error) {
	dirs := libSearchDirs(tri, sysroot)
	out := make([]resolvedLib, 0, len(names))
	for _, name := range names {
		// These darwin libraries live only in the dyld shared cache on macOS 12+.
		// No on-disk file exists to read — only the LC_LOAD_DYLIB name is needed.
		if tri.os == "darwin" && isDarwinCacheOnly(name) {
			out = append(out, resolvedLib{name: name, bytes: nil})
			continue
		}
		path, err := findLib(name, dirs)
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("syslibs: read %s: %w", path, err)
		}
		out = append(out, resolvedLib{name: name, bytes: data})
	}
	return out, nil
}

// isDarwinCacheOnly reports whether name is a darwin library that lives only
// in the dyld shared cache and has no on-disk file to read.
func isDarwinCacheOnly(name string) bool {
	switch name {
	case "libSystem.B.dylib", "libobjc.dylib", "libobjc.A.dylib",
		"libSystem.dylib", "libc.dylib", "libpthread.dylib",
		"libm.dylib", "libdyld.dylib", "libc++.1.dylib":
		return true
	}
	// ObjC frameworks: "Foundation.framework/Foundation" etc.
	return strings.Contains(name, ".framework/")
}

// isFrameworkPath reports whether name is an Apple framework dylib path of the
// form "X.framework/X" injected by declareExterns for ObjC darwin bindings.
func isFrameworkPath(name string) bool {
	return strings.Contains(name, ".framework/")
}

// resolveCRT finds and reads the system CRT objects required for a fully
// linked ELF executable. The canonical link order is:
//
//	crt1.o  crti.o  <user object>  crtn.o
//
// crt1.o provides _start which calls __libc_start_main which calls main.
// crti.o / crtn.o bracket the .init/.fini sections.
//
// CRT injection is only needed for Linux ELF targets. Darwin's libSystem
// provides its own startup machinery and the Mach-O linker handles it
// transparently; Windows uses a different startup model entirely.
// For those targets resolveCRT returns an empty crtObjects with no error.
func resolveCRT(tri triple, sysroot string) (crtObjects, error) {
	if tri.os != "linux" {
		return crtObjects{}, nil
	}

	dirs := libSearchDirs(tri, sysroot)

	read := func(name string) ([]byte, error) {
		path, err := findLib(name, dirs)
		if err != nil {
			return nil, fmt.Errorf("syslibs: crt: %w", err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("syslibs: read %s: %w", path, err)
		}
		return data, nil
	}

	crt1, err := read("crt1.o")
	if err != nil {
		return crtObjects{}, err
	}
	crti, err := read("crti.o")
	if err != nil {
		return crtObjects{}, err
	}
	crtn, err := read("crtn.o")
	if err != nil {
		return crtObjects{}, err
	}

	return crtObjects{crt1: crt1, crti: crti, crtn: crtn}, nil
}

// findLib searches dirs in order and returns the first path where name exists.
func findLib(name string, dirs []string) (string, error) {
	for _, dir := range dirs {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	if len(dirs) == 0 {
		return "", fmt.Errorf("syslibs: no search directories configured for this target; use -sysroot to specify a sysroot")
	}
	return "", fmt.Errorf("syslibs: %q not found; searched: %s", name, strings.Join(dirs, ", "))
}

// libSearchDirs returns the ordered list of directories to probe for shared
// libraries and CRT objects for the given triple. If sysroot is non-empty
// it is prepended to every path.
func libSearchDirs(tri triple, sysroot string) []string {
	var dirs []string

	switch tri.os {
	case "linux":
		switch tri.arch {
		case "arm64":
			dirs = []string{
				"/usr/lib/aarch64-linux-gnu",
				"/lib/aarch64-linux-gnu",
				"/usr/lib",
				"/lib",
			}
		case "riscv64":
			dirs = []string{
				"/usr/lib/riscv64-linux-gnu",
				"/lib/riscv64-linux-gnu",
				"/usr/lib",
				"/lib",
			}
		default: // amd64
			dirs = []string{
				"/usr/lib/x86_64-linux-gnu",
				"/lib/x86_64-linux-gnu",
				"/usr/lib",
				"/lib",
			}
		}

	case "darwin":
		dirs = []string{
			"/usr/lib",
			"/usr/local/lib",
			"/opt/homebrew/lib",
		}

	// Windows: no native host paths; requires explicit -sysroot.
	}

	if sysroot == "" {
		return dirs
	}
	rooted := make([]string, len(dirs))
	for i, d := range dirs {
		rooted[i] = filepath.Join(sysroot, d)
	}
	return rooted
}

func extractLibFuncSymbols(m *vertex.Module, triOS string) map[string][]string {
	result := make(map[string][]string)
	for _, imp := range m.Imports.Imports {
		platform, lib, ok := splitImportModule(imp.Module)
		if !ok || platform != triOS {
			continue
		}
		lib = normalizeDarwinLib(lib) // same normalization as extractDynLibs
		if _, ok := imp.Desc.(vertex.FuncImport); ok {
			result[lib] = append(result[lib], imp.Name)
		}
	}
	return result
}