package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/vertex-language/pkg"
	"github.com/vertex-language/pkg/importer"
	"github.com/vertex-language/pkg/mod"
)

func runMod(args []string, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: vertex mod init|get|tidy ...")
		return 2
	}
	switch args[0] {
	case "init":
		return modInit(args[1:], stderr)
	case "get":
		return modGet(args[1:], stderr)
	case "tidy":
		fmt.Fprintln(stderr, "vertex: mod tidy is not implemented yet")
		return 1
	default:
		fmt.Fprintf(stderr, "vertex: unknown mod subcommand %q\n", args[0])
		return 2
	}
}

func modInit(args []string, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: vertex mod init <module-path>")
		return 2
	}
	if !mod.IsValidModulePath(args[0]) {
		fmt.Fprintf(stderr, "vertex: %q is not a valid module path\n", args[0])
		return 2
	}
	if _, err := os.Stat("vs.mod"); err == nil {
		fmt.Fprintln(stderr, "vertex: vs.mod already exists in this directory")
		return 1
	}
	content := fmt.Sprintf("module %s\n\nvertex %s\n", args[0], version)
	if err := os.WriteFile("vs.mod", []byte(content), 0o644); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	fmt.Printf("initialized vs.mod for %s\n", args[0])
	return 0
}

func modGet(args []string, stderr io.Writer) int {
	fs := newFlagSet("mod get", stderr)
	var homeOverride string
	fs.StringVar(&homeOverride, "vertex-home", "", "override $VERTEX_HOME for this invocation")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: vertex mod get <module-path>[@version]")
		return 2
	}

	spec := fs.Arg(0)
	path, query := spec, "latest"
	if i := strings.LastIndexByte(spec, '@'); i >= 0 {
		path, query = spec[:i], spec[i+1:]
	}

	data, err := os.ReadFile("vs.mod")
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v (run `vertex mod init <module-path>` first)\n", err)
		return 1
	}
	mf, err := mod.ParseLax("vs.mod", data, nil)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: vs.mod: %v\n", err)
		return 1
	}

	homeDir, err := pkg.Home(homeOverride)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	fetcher := importer.NewGitFetcher()
	cache, err := pkg.OpenCache(homeDir, fetcher)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	resolved, err := fetcher.Resolve(mod.ModulePath(path), query)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: resolving %s: %v\n", spec, err)
		return 1
	}

	// This is the one place a version is chosen without one already being
	// pinned; every build/test invocation defaults to -mod=readonly
	// specifically so that never happens implicitly.
	if _, err := cache.Mod(mod.ModulePath(path), resolved, "vs.sum", pkg.ModUpdate); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	if err := addOrUpdateDependency(mf, mod.ModulePath(path), resolved); err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}

	fmt.Printf("added %s %s\n", path, resolved)
	return 0
}

// addOrUpdateDependency rewrites vs.mod's dependencies block to include
// path@version, replacing any existing requirement on the same path.
//
// This works from mf.Dependencies (the interpreted view) rather than
// editing mf.Syntax in place, so a hand-authored vs.mod's comments and
// factoring aren't guaranteed to survive a `mod get` — a real
// implementation would mutate FileSyntax directly and go through
// mod.Format so a round-trip preserves both. Good enough for a freshly
// `vertex mod init`-ed file, which is the common case.
func addOrUpdateDependency(mf *mod.File, path mod.ModulePath, version string) error {
	deps := make(map[mod.ModulePath]string, len(mf.Dependencies)+1)
	for _, d := range mf.Dependencies {
		deps[d.Mod.Path] = d.Mod.Version
	}
	deps[path] = version

	paths := make([]mod.ModulePath, 0, len(deps))
	for p := range deps {
		paths = append(paths, p)
	}
	sort.Slice(paths, func(i, j int) bool { return paths[i] < paths[j] })

	var sb strings.Builder
	fmt.Fprintf(&sb, "module %s\n\n", mf.Module.Path)
	if mf.Vertex != nil {
		fmt.Fprintf(&sb, "vertex %s\n\n", mf.Vertex.Version)
	}
	sb.WriteString("dependencies (\n")
	for _, p := range paths {
		fmt.Fprintf(&sb, "\t%s %s\n", p, deps[p])
	}
	sb.WriteString(")\n")

	return os.WriteFile("vs.mod", []byte(sb.String()), 0o644)
}