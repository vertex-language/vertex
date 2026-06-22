package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vertex-language/ir/vertex"
)

// resolvedLib is a shared library name paired with its raw bytes, ready to
// hand directly to the linker's AddDynamicLibrary method.
type resolvedLib struct {
	name  string
	bytes []byte
}

// extractDynLibs walks the VIR import section and returns the unique set of
// shared library names required for tri.os.
//
// Import Module strings follow the "<platform>:<libname>" convention defined
// in §13.1 of the VIR spec (e.g. "linux:libc.so.6", "darwin:libSystem.B.dylib").
// Only entries whose platform component equals tri.os are collected; imports
// for other platforms present in the same module (cross-compiled packages,
// conditional FFI) are silently skipped. Only FuncImport descriptors can
// actually reference a shared library symbol, but we filter on the platform
// prefix alone so table/global/memory imports to platform libs are also covered
// if they ever appear.
func extractDynLibs(m *vertex.Module, tri triple) []string {
	seen := make(map[string]bool)
	var libs []string
	for _, imp := range m.Imports.Imports {
		platform, lib, ok := splitImportModule(imp.Module)
		if !ok || platform != tri.os {
			continue
		}
		if !seen[lib] {
			seen[lib] = true
			libs = append(libs, lib)
		}
	}
	return libs
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
// supports cross-compilation (e.g. sysroot="/opt/sysroot/aarch64-linux-gnu").
// An empty sysroot searches the native host layout.
//
// Libraries that cannot be found produce an error that names both the library
// and the directories that were searched, to help diagnose missing dev packages.
func resolveLibs(names []string, tri triple, sysroot string) ([]resolvedLib, error) {
	dirs := libSearchDirs(tri, sysroot)
	out := make([]resolvedLib, 0, len(names))
	for _, name := range names {
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
// libraries for the given triple. If sysroot is non-empty it is prepended to
// every path.
//
// Linux multiarch directories are listed first (Debian/Ubuntu layout), followed
// by the legacy single-arch fallbacks so older distributions still work.
// Darwin lists the standard SDK paths plus Homebrew's prefix.
// Windows cross-compilation requires an explicit sysroot since there are no
// standard Windows library paths on Linux/macOS hosts.
func libSearchDirs(tri triple, sysroot string) []string {
	var dirs []string

	switch tri.os {
	case "linux":
		switch tri.arch {
		case "arm64":
			dirs = []string{
				"/lib/aarch64-linux-gnu",
				"/usr/lib/aarch64-linux-gnu",
				"/lib",
				"/usr/lib",
			}
		case "riscv64":
			dirs = []string{
				"/lib/riscv64-linux-gnu",
				"/usr/lib/riscv64-linux-gnu",
				"/lib",
				"/usr/lib",
			}
		default: // amd64
			dirs = []string{
				"/lib/x86_64-linux-gnu",
				"/usr/lib/x86_64-linux-gnu",
				"/lib",
				"/usr/lib",
			}
		}

	case "darwin":
		dirs = []string{
			"/usr/lib",
			"/usr/local/lib",
			"/opt/homebrew/lib",
		}

	// Windows: no native host paths are probed. An explicit -sysroot pointing
	// at a MinGW or MSVC sysroot is required for cross-compilation.
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