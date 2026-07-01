package codegen

// Options is forwarded to the instruction selector.
type Options struct {
	OptLevel  int  // 0=none 1=light 2=full -1=size
	DebugInfo bool
}