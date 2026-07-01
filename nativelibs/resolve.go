package nativelibs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vertex-language/vertex/target"
)

type ResolvedLib struct {
	Name  string
	Bytes []byte
}

type CRTObjects struct {
	CRT1 []byte
	CRTI []byte
	CRTN []byte
}

// ResolveLibs reads each of names off disk, searching dirs in order. A
// darwin library that lives only in the dyld shared cache (see
// isDarwinCacheOnly) resolves to a ResolvedLib with Bytes == nil rather
// than being read; a windows API Set virtual DLL is skipped entirely —
// neither has a real on-disk file for this function to read.
func ResolveLibs(names []string, tri target.Triple, dirs []string) ([]ResolvedLib, error) {
	out := make([]ResolvedLib, 0, len(names))
	for _, name := range names {
		if tri.OS == "darwin" && isDarwinCacheOnly(name) {
			out = append(out, ResolvedLib{Name: name, Bytes: nil})
			continue
		}
		if tri.OS == "windows" && isWindowsAPISet(name) {
			continue
		}
		path, err := findLib(name, dirs)
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("nativelibs: read %s: %w", path, err)
		}
		out = append(out, ResolvedLib{Name: name, Bytes: data})
	}
	return out, nil
}

// ResolveCRT finds and reads the system CRT objects a fully linked ELF
// executable needs bracketing every other object:
//
//	crt1.o  crti.o  <objects...>  crtn.o
//
// crt1.o provides _start, which calls __libc_start_main, which calls
// main; crti.o/crtn.o bracket the .init/.fini sections. Only linux ELF
// targets need this — darwin's libSystem provides its own startup
// machinery and windows uses a different startup model entirely — so
// ResolveCRT returns an empty CRTObjects with no error for those.
func ResolveCRT(tri target.Triple, dirs []string) (CRTObjects, error) {
	if tri.OS != "linux" {
		return CRTObjects{}, nil
	}
	read := func(name string) ([]byte, error) {
		path, err := findLib(name, dirs)
		if err != nil {
			return nil, fmt.Errorf("nativelibs: crt: %w", err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("nativelibs: read %s: %w", path, err)
		}
		return data, nil
	}
	crt1, err := read("crt1.o")
	if err != nil {
		return CRTObjects{}, err
	}
	crti, err := read("crti.o")
	if err != nil {
		return CRTObjects{}, err
	}
	crtn, err := read("crtn.o")
	if err != nil {
		return CRTObjects{}, err
	}
	return CRTObjects{CRT1: crt1, CRTI: crti, CRTN: crtn}, nil
}

func isWindowsAPISet(name string) bool {
	s := strings.ToLower(name)
	return strings.HasPrefix(s, "api-ms-win-") || strings.HasPrefix(s, "ext-ms-win-")
}

func isDarwinCacheOnly(name string) bool {
	switch name {
	case "libSystem.B.dylib", "libobjc.dylib", "libobjc.A.dylib",
		"libSystem.dylib", "libc.dylib", "libpthread.dylib",
		"libm.dylib", "libdyld.dylib", "libc++.1.dylib":
		return true
	}
	return strings.Contains(name, ".framework/")
}

func findLib(name string, dirs []string) (string, error) {
	for _, dir := range dirs {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	if len(dirs) == 0 {
		return "", fmt.Errorf("nativelibs: no search directories configured for this target; use -sysroot to specify a sysroot")
	}
	return "", fmt.Errorf("nativelibs: %q not found; searched: %s", name, strings.Join(dirs, ", "))
}