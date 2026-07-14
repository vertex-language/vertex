package driver

import "github.com/vertex-language/pkg"

type EmitMode uint8

const (
	ModeVIR EmitMode = iota
	ModeMIR
	ModeASM
	ModeObj
	ModeExe
	ModeDump
	ModeTest
	ModeRun // compile to a temp binary and execute; the default when no -o and no emit flag was given
)

type Config struct {
	Input      string
	Output     string
	Target     string
	Sysroot    string
	VertexHome string       // "-vertex-home" override; "" defers to $VERTEX_HOME, then ~/.vertex
	LoadMode   pkg.LoadMode // governs whether an unrecorded dependency may be fetched; see pkg.LoadMode
	Mode       EmitMode
	OptLevel   int
	DebugInfo  bool
	TestDir    string
	TestFile   string
}