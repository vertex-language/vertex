package builtins

import "sort"

// Feature names one language capability a program either uses or does not.
// lower/hir already knows whether a program contains a map, a chan, an async
// function, or a shared T; this is the channel that fact reaches the driver
// through, so builtins.Modules can build only what is needed plus its
// closure (§5.2, "Build only what the program needs").
type Feature uint32

const (
	FeatMemory Feature = 1 << iota
	FeatPanic
	FeatFmt
	FeatConsole
	FeatProcess
	FeatRC
	FeatSlice
	FeatMap
	FeatString
	FeatChan
	FeatReactor
	FeatSync
	FeatPoll
	FeatThread
	FeatIO
	FeatTime
)

// FeatureSet is a set of Features. The zero value is the empty set; Closure
// is what turns it into a buildable list.
type FeatureSet uint32

func (s FeatureSet) Has(f Feature) bool     { return s&FeatureSet(f) != 0 }
func (s FeatureSet) With(f Feature) FeatureSet { return s | FeatureSet(f) }
func (s FeatureSet) Union(o FeatureSet) FeatureSet { return s | o }

// deps is the closure relation: asking for chan implicitly pulls sync, rc,
// and memory. Written as one table so the relation has a single home.
var deps = map[Feature]FeatureSet{
	FeatSlice:   FeatureSet(FeatMemory | FeatPanic),
	FeatMap:     FeatureSet(FeatMemory | FeatString),
	FeatString:  FeatureSet(FeatMemory),
	FeatRC:      FeatureSet(FeatMemory),
	FeatChan:    FeatureSet(FeatSync | FeatRC | FeatMemory),
	FeatReactor: FeatureSet(FeatPoll | FeatChan | FeatMemory),
	FeatThread:  FeatureSet(FeatMemory),
	FeatSync:    FeatureSet(FeatMemory),
	FeatPoll:    FeatureSet(FeatMemory),
	FeatIO:      FeatureSet(FeatMemory | FeatString),
	FeatFmt:     FeatureSet(FeatConsole | FeatString),
	FeatPanic:   FeatureSet(FeatFmt | FeatConsole | FeatProcess),
}

// floor is what every program links regardless: hello-world gets memory,
// console, panic, fmt, and process, and nothing else (§5.2).
const floor = FeatureSet(FeatMemory | FeatConsole | FeatPanic | FeatFmt | FeatProcess)

// Closure returns s plus everything reachable from it, plus the floor.
func (s FeatureSet) Closure() FeatureSet {
	out := s | floor
	for {
		next := out
		for f, d := range deps {
			if out.Has(f) {
				next |= d
			}
		}
		if next == out {
			return out
		}
		out = next
	}
}

var featureModule = map[Feature]string{
	FeatMemory: ModuleMemory, FeatPanic: ModulePanic, FeatFmt: ModuleFmt,
	FeatConsole: ModuleConsole, FeatProcess: ModuleProcess, FeatRC: ModuleRC,
	FeatSlice: ModuleSlice, FeatMap: ModuleMap, FeatString: ModuleString,
	FeatChan: ModuleChan, FeatReactor: ModuleReactor, FeatSync: ModuleSync,
	FeatPoll: ModulePoll, FeatThread: ModuleThread, FeatIO: ModuleIO,
	FeatTime: ModuleTime,
}

// Modules returns the module names s's closure requires, sorted so a build
// is byte-reproducible.
func (s FeatureSet) Modules() []string {
	c := s.Closure()
	var out []string
	for f, m := range featureModule {
		if c.Has(f) {
			out = append(out, m)
		}
	}
	sort.Strings(out)
	return out
}

// FeatureFor reports which Feature owns a Symbol's module, so a call site
// that emits a builtin call can record its need without naming the feature
// twice.
func FeatureFor(s Symbol) (Feature, bool) {
	for f, m := range featureModule {
		if m == s.Module {
			return f, true
		}
	}
	return 0, false
}