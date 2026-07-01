package nativelibs

import (
	"path/filepath"

	"github.com/vertex-language/pkg"

	"github.com/vertex-language/vertex/target"
)

// SearchDirs returns the ordered list of directories to search for shared
// libraries and CRT objects: installed's cache directories first (from
// pkg.Graph.EnsureNativeLibs), then sysroot-relative system paths as
// fallback — never the reverse, so a project's own resolved vs.lib
// version always wins over whatever happens to be on the host. installed
// is nil for a standalone, module-less compile (no vs.lib resolution is
// possible without a pkg.Graph), which just means only the system
// fallback applies.
func SearchDirs(tri target.Triple, sysroot string, installed []pkg.LibResult) []string {
	dirs := make([]string, 0, len(installed)+4)
	for _, r := range installed {
		dirs = append(dirs, r.Dir)
	}
	dirs = append(dirs, systemDirs(tri, sysroot)...)
	return dirs
}

func systemDirs(tri target.Triple, sysroot string) []string {
	var dirs []string
	switch tri.OS {
	case "linux":
		switch tri.Arch {
		case "arm64":
			dirs = []string{"/usr/lib/aarch64-linux-gnu", "/lib/aarch64-linux-gnu", "/usr/lib", "/lib"}
		case "riscv64":
			dirs = []string{"/usr/lib/riscv64-linux-gnu", "/lib/riscv64-linux-gnu", "/usr/lib", "/lib"}
		default:
			dirs = []string{"/usr/lib/x86_64-linux-gnu", "/lib/x86_64-linux-gnu", "/usr/lib", "/lib"}
		}
	case "darwin":
		dirs = []string{"/usr/lib", "/usr/local/lib", "/opt/homebrew/lib"}
	case "windows":
		dirs = []string{`C:\Windows\System32`, `C:\Windows\SysWOW64`, `C:\Windows\System`}
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