package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ImportKind classifies an import path by its routing prefix.
type ImportKind int

const (
	ImportLib     ImportKind = iota // "lib/..."     → linked C library
	ImportLinux                      // "linux/..."   → inline syscall
	ImportDarwin                     // "darwin/..."  → objc_msgSend dispatch
	ImportWindows                    // "windows/..." → COM vtable dispatch
	ImportGPU                        // "gpu/..."     → PTX / shader kernel
	ImportMetal                      // "metal/..."   → bare metal interrupt
	ImportModule                     // "github.com/..." → Vertex package
	ImportLocal                      // "./..." or "../..." → relative path
)

// nativePrefixes maps import path prefixes to their ImportKind.
// Order matters: longer prefixes must not shadow shorter ones here, but
// since each prefix is distinct this is just a linear scan.
var nativePrefixes = []struct {
	prefix string
	kind   ImportKind
}{
	{"lib/",     ImportLib},
	{"linux/",   ImportLinux},
	{"darwin/",  ImportDarwin},
	{"windows/", ImportWindows},
	{"gpu/",     ImportGPU},
	{"metal/",   ImportMetal},
}

// Import is a parsed import declaration.
type Import struct {
	Raw       string
	Kind      ImportKind
	Prefix    string // routing prefix without trailing slash, e.g. "lib", "linux", "gpu"
	Namespace string // last path segment — matches the class parent name, e.g. "sdl2", "syscalls"
}

// IsNative reports whether this import binds to a native class.
func (imp *Import) IsNative() bool {
	switch imp.Kind {
	case ImportLib, ImportLinux, ImportDarwin, ImportWindows, ImportGPU, ImportMetal:
		return true
	}
	return false
}

// ParseImportPath classifies and decomposes an import path string.
func ParseImportPath(path string) *Import {
	imp := &Import{Raw: path}

	for _, p := range nativePrefixes {
		if strings.HasPrefix(path, p.prefix) {
			imp.Kind = p.kind
			imp.Prefix = strings.TrimSuffix(p.prefix, "/")
			// Namespace is the last segment of everything after the prefix.
			// "lib/sdl2"            → "sdl2"
			// "darwin/objc/foundation" → "foundation"
			rest := strings.TrimPrefix(path, p.prefix)
			parts := strings.Split(rest, "/")
			imp.Namespace = parts[len(parts)-1]
			return imp
		}
	}

	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
		imp.Kind = ImportLocal
		return imp
	}

	imp.Kind = ImportModule
	return imp
}

// WasmModule returns the wasm import module string for this import.
//   lib/sdl2              → "sdl2"
//   linux/syscalls        → "linux:syscalls"
//   darwin/objc/foundation → "darwin:objc:foundation"
//   gpu/cuda              → "gpu:cuda"
func (imp *Import) WasmModule(_ *BuildTags) string {
	switch imp.Kind {
	case ImportLib:
		// For linked libraries the bare namespace is the module name.
		return imp.Namespace
	default:
		// For all other native strategies embed the full routing path.
		rest := strings.TrimPrefix(imp.Raw, imp.Prefix+"/")
		return imp.Prefix + ":" + strings.ReplaceAll(rest, "/", ":")
	}
}

// Resolver resolves import paths to .vs source files.
type Resolver struct {
	ModRoot string
	ModPath string
}

func NewResolver(root, modPath string) *Resolver {
	return &Resolver{ModRoot: root, ModPath: modPath}
}

// ResolveFiles returns the .vs files for an import path.
// Returns (nil, nil) for native imports — those are handled by native classes.
func (r *Resolver) ResolveFiles(importPath, fromDir string) ([]string, error) {
	imp := ParseImportPath(importPath)

	switch imp.Kind {
	case ImportLib, ImportLinux, ImportDarwin, ImportWindows, ImportGPU, ImportMetal:
		return nil, nil // native interface; no .vs files

	case ImportLocal:
		dir := filepath.Join(fromDir, importPath)
		return vertexFiles(dir)

	case ImportModule:
		if r.ModPath != "" && strings.HasPrefix(importPath, r.ModPath+"/") {
			rel := strings.TrimPrefix(importPath, r.ModPath+"/")
			dir := filepath.Join(r.ModRoot, rel)
			return vertexFiles(dir)
		}
		vdir := filepath.Join(r.ModRoot, "vendor", importPath)
		if fs, err := vertexFiles(vdir); err == nil {
			return fs, nil
		}
		return nil, fmt.Errorf("cannot find package %q (not in vendor/)", importPath)
	}

	return nil, fmt.Errorf("unrecognised import path %q", importPath)
}

func vertexFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".vs") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no .vs files in %s", dir)
	}
	return out, nil
}