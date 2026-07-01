package pipeline

import (
	"fmt"
	"io"

	mirtext "github.com/vertex-language/ir/machine/encoding/text/mir"
	virtext "github.com/vertex-language/ir/vertex/encoding/text"
)

// Run executes stages in order against st, stopping after the stage named
// upTo ("" runs every stage). If st.Sink is set, Run writes a banner and
// that stage's rendered artifact to it before moving on — this is what
// -dump uses to capture the whole pipeline in one annotated file; every
// other emit mode leaves Sink nil and just reads st's final field values.
func Run(stages []Stage, st *State, upTo string) error {
	st.UpTo = upTo
	for _, s := range stages {
		if st.Sink != nil {
			banner(st.Sink, s.Name)
		}
		if err := s.Run(st); err != nil {
			if st.Sink != nil {
				fmt.Fprintf(st.Sink, "; ERROR: %v\n; (pipeline stopped)\n", err)
			}
			return fmt.Errorf("%s: %w", s.Name, err)
		}
		if st.Sink != nil {
			writeArtifact(st.Sink, s.Name, st)
		}
		if s.Name == upTo {
			break
		}
	}
	return nil
}

func banner(w io.Writer, name string) {
	const line = "════════════════════════════════════════════════════════════"
	fmt.Fprintf(w, "; ════ %-40s%s\n\n", name, line)
}

func writeArtifact(w io.Writer, stage string, st *State) {
	switch stage {
	case StageVIR:
		for _, u := range st.Units {
			if u.VIR == nil {
				continue
			}
			fmt.Fprintf(w, "; module: %s\n", unitLabel(u))
			if u.VIRErr != nil {
				fmt.Fprintf(w, "; WARNING: %v\n", u.VIRErr)
			}
			io.WriteString(w, virtext.Format(u.VIR))
			io.WriteString(w, "\n\n")
		}
	case StageMIR:
		for _, u := range st.Units {
			if u.MIR == nil {
				continue
			}
			fmt.Fprintf(w, "; module: %s\n", unitLabel(u))
			io.WriteString(w, mirtext.PrintModule(u.MIR))
			io.WriteString(w, "\n\n")
		}
	case StageASM:
		for _, u := range st.Units {
			if u.ASM == "" {
				continue
			}
			fmt.Fprintf(w, "; module: %s\n", unitLabel(u))
			io.WriteString(w, u.ASM)
			io.WriteString(w, "\n")
		}
	case StageObject:
		for _, u := range st.Units {
			if u.Obj == nil {
				continue
			}
			fmt.Fprintf(w, "; module: %s — object, %d bytes\n", unitLabel(u), len(u.Obj))
		}
		io.WriteString(w, "\n")
	case StageLink:
		fmt.Fprintf(w, "; linked executable, %d bytes\n", len(st.Exe))
	}
}