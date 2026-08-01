package token

// BuildTag is the A.2.2 target selector. It lives in token rather than parser
// because parser, loader, and the driver all need it, and because it changes
// what is grammatical (build test licenses the `test` marker and Expected).
type BuildTag int

const (
	TagNone BuildTag = iota // no BuildClause in the file
	TagLinux
	TagWindows
	TagDarwin
	TagJS
	TagWasm
	TagTest
)

var buildTags = [...]string{
	TagNone:    "",
	TagLinux:   "linux",
	TagWindows: "windows",
	TagDarwin:  "darwin",
	TagJS:      "js",
	TagWasm:    "wasm",
	TagTest:    "test",
}

func (t BuildTag) String() string {
	if t >= 0 && int(t) < len(buildTags) {
		return buildTags[t]
	}
	return "BuildTag(?)"
}

// LookupBuildTag reports the tag for s. The bool is load-bearing: A.2.2 makes
// an unrecognised tag a compile error, never a silently-excluded file, so the
// caller must be able to tell "unknown" from "not this target".
func LookupBuildTag(s string) (BuildTag, bool) {
	for i, name := range buildTags {
		if name != "" && name == s {
			return BuildTag(i), true
		}
	}
	return TagNone, false
}

// LicensesTest reports whether this tag makes the `test` FunctionMarker and
// Expected(...) result types grammatical (A.2.2, A.12.1).
func LicensesTest(t BuildTag) bool { return t == TagTest }

// HasFrameworks reports whether the target platform has a first-class notion of
// a bundled versioned library, i.e. whether `declare framework` is legal
// (A.8.1). Checked by the analyzer; the predicate lives here so the answer has
// one home.
func HasFrameworks(t BuildTag) bool { return t == TagDarwin }