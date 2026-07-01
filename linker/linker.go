package linker

import (
	"fmt"

	"github.com/vertex-language/vertex/target"
)

// DynLib is one dynamic library to link against. Bytes is the library's
// raw on-disk content to embed a reference to; a nil Bytes with a
// non-empty Name marks a library the platform resolves at load time from
// its own shared cache rather than from an on-disk file (only meaningful
// on darwin — see AddCachedDylib in macho.go).
type DynLib struct {
	Name  string
	Bytes []byte
}

// CRTObjects are the system C runtime start/end objects a linux ELF
// executable needs bracketing every other object. Unused on darwin
// (libSystem provides its own startup machinery) and windows (a different
// startup model entirely).
type CRTObjects struct {
	CRT1 []byte // _start → __libc_start_main → main
	CRTI []byte // .init section prologue
	CRTN []byte // .init section epilogue
}

// OutputType selects a linker-specific output flavor. Currently only
// meaningful for windows, where PIE additionally requests DYNAMIC_BASE,
// HIGH_ENTROPY_VA, and a .reloc section.
type OutputType uint8

const (
	OutputDefault OutputType = iota
	OutputPIE
)

// Options configures a Linker before any objects are added.
type Options struct {
	EntryPoint string
	DynLibs    []DynLib
	CRT        CRTObjects
	LibSymbols map[string][]string
	OutputType OutputType
}

// Linker accumulates object files for one target and links them into a
// final executable. Every AddObject call must precede Link. This
// interface deliberately mirrors pkg/provider.Provider's own "one
// interface, N OS-specific implementations, callers don't switch on OS
// themselves" shape — this toolchain already has a working pattern for
// exactly that, in a different corner of the same repo family.
type Linker interface {
	AddObject(name string, data []byte) error
	Link() ([]byte, error)
}

// New returns the Linker implementation for tri.OS, already primed with
// opts (CRT objects bracketing every other object on ELF, entry point and
// SO-needed set for darwin, and so on).
func New(tri target.Triple, opts Options) (Linker, error) {
	switch tri.OS {
	case "linux":
		return newELFLinker(tri, opts)
	case "darwin":
		return newMachOLinker(tri, opts)
	case "windows":
		return newPELinker(tri, opts)
	}
	return nil, fmt.Errorf("linker: unsupported OS %q", tri.OS)
}