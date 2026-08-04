package token

// BuildTag is the target selector a BuildClause names.
//
// It lives here rather than in the parser because the parser, the loader, and
// the driver all need it, and because it changes what is admissible: a file's
// tag is what licenses an ExpectedType result.
type BuildTag int

const (
	TagNone BuildTag = iota // no build clause in the file
	TagLinux
	TagWindows
	TagDarwin
	TagJS
	TagWasm
	TagTest
)

// buildTags holds each tag's spelling. Every spelling is a contextual keyword,
// recognized only in BuildTag position.
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

func (t BuildTag) IsValid() bool { return t != TagNone }

// LookupBuildTag reports the tag spelled s.
//
// The bool is load-bearing. An unrecognized tag is a compile error, never a
// silently excluded file, so a caller must be able to tell an unknown spelling
// from TagNone's "no clause at all" — which a zero-value return cannot do.
// Implementations may recognize tags beyond this set.
func LookupBuildTag(s string) (BuildTag, bool) {
	for i, name := range buildTags {
		if name != "" && name == s {
			return BuildTag(i), true
		}
	}
	return TagNone, false
}

// LicensesTest reports whether a file built under t admits an ExpectedType
// result. The grammar admits ExpectedType only through DeclResult; restricting
// it further to this tag is a static rule, and this predicate is what that rule
// reads.
func LicensesTest(t BuildTag) bool { return t == TagTest }