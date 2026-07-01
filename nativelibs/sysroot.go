package nativelibs

import (
	"os"
	"runtime"

	"github.com/vertex-language/vertex/target"
)

// AutoSysroot probes well-known cross-compilation toolchain roots when
// the user hasn't specified -sysroot explicitly. On native builds (target
// OS matches host OS) it returns "" so absolute system paths are used
// directly. For Linux → Windows cross-compilation it probes MinGW-w64.
func AutoSysroot(tri target.Triple) string {
	nativeOS := runtime.GOOS
	if tri.OS == nativeOS {
		return ""
	}
	if tri.OS == "darwin" {
		return ""
	}
	if nativeOS == "linux" && tri.OS == "windows" {
		var candidates []string
		switch tri.Arch {
		case "arm64":
			candidates = []string{"/usr/aarch64-w64-mingw32", "/opt/aarch64-w64-mingw32"}
		default:
			candidates = []string{"/usr/x86_64-w64-mingw32", "/opt/x86_64-w64-mingw32", "/usr/mingw64"}
		}
		for _, p := range candidates {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	if nativeOS == "linux" && tri.OS == "linux" {
		switch tri.Arch {
		case "arm64":
			if _, err := os.Stat("/usr/aarch64-linux-gnu"); err == nil {
				return "/"
			}
		case "riscv64":
			if _, err := os.Stat("/usr/riscv64-linux-gnu"); err == nil {
				return "/"
			}
		}
	}
	return ""
}