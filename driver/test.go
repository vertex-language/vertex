// driver/test.go
//
// The `vertex test` runner.
//
// Two things shape this file. First, hir monomorphizes from a single root,
// and under ModeTest that root is one test function — so a package with
// twelve tests is twelve lowerings and twelve binaries, not one. Second,
// Expected(error) tests are *supposed* to fail to compile, which means
// their diagnostics arrive from the checker before any lowering happens.
// So the runner checks once, partitions the resulting diagnostics by which
// test function's extent they land in, and only then builds and runs the
// tests that were meant to compile.
//
// A third thing shapes the suite-discovery half at the bottom of this
// file: `loadDir` (load.go) treats a directory as one package only if it
// holds .vs files directly, the same directory-granularity importer.Load
// and parser.ParseDir already assume. A directory that holds none but has
// subdirectories that do — `testutils/` over `testutils/01_values/`,
// `testutils/02_operators/`, ... — is not itself a package, so RunTests
// alone cannot point at it. RunTestsAuto tells the two cases apart and
// RunTestSuite runs every discovered package, aggregating the result. No
// new CLI flag is needed: whether a directory is a single package or a
// suite root is a fact about its own contents, never ambiguous, so it's
// decided by inspection rather than by asking the caller to say which.
package driver

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/diag"
	"github.com/vertex-language/vertex/token"
)

// TestKind distinguishes the three shapes A.12.2 allows.
type TestKind int

const (
	// TestRuns has no result type: it passes if it compiles and runs
	// without crashing.
	TestRuns TestKind = iota
	// TestValue is Expected(Type, "rendered"): compare the program's
	// output against the string.
	TestValue
	// TestError is Expected(error) or Expected(error, "message"): the
	// function must fail to compile, optionally with that exact message.
	TestError
)

// TestCase is one discovered `test`-marked function.
type TestCase struct {
	Name     string
	Kind     TestKind
	Expected string // the rendered value, or the diagnostic text
	HasText  bool   // whether Expected was written at all
	Pos, End token.Pos
}

// DiscoverTests finds every test function in a package.
//
// Discovery is done over the AST rather than over checked objects on
// purpose: an Expected(error) test may well be the reason the package
// didn't fully check, and a runner that can only see successfully-checked
// objects would silently skip exactly the tests that matter most.
func DiscoverTests(pkg *Package) []TestCase {
	var out []TestCase
	for _, file := range pkg.Files {
		for _, d := range file.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Recv != nil || fd.Type == nil {
				continue
			}
			if fd.Type.Marker == nil || fd.Type.Marker.Name != token.CtxTest {
				continue
			}
			tc := TestCase{
				Name: fd.Name.Name,
				Pos:  fd.Pos(),
				End:  fd.End(),
			}
			classifyExpected(fd.Type.Result, &tc)
			out = append(out, tc)
		}
	}
	return out
}

// classifyExpected reads the `-> Expected(...)` clause. It is a plain
// CallExpr in the tree (A.12.2 gives it no node of its own), so this is a
// shape check over an ordinary call: first argument names the type — or
// the identifier `error` — and the optional second is the string literal.
func classifyExpected(result ast.Expr, tc *TestCase) {
	call, ok := result.(*ast.CallExpr)
	if !ok {
		tc.Kind = TestRuns
		return
	}
	fun, ok := call.Fun.(*ast.Ident)
	if !ok || fun.Name != "Expected" || len(call.Args) == 0 {
		tc.Kind = TestRuns
		return
	}

	tc.Kind = TestValue
	if id, ok := call.Args[0].(*ast.Ident); ok && id.Name == token.CtxError {
		tc.Kind = TestError
	}
	if len(call.Args) >= 2 {
		if lit, ok := call.Args[1].(*ast.BasicLit); ok {
			tc.Expected = unquote(lit.Value)
			tc.HasText = true
		}
	}
}

// unquote decodes a string literal's source spelling. The scanner keeps
// literals raw (decoding is the analyzer's job), and this runner needs the
// value, so it does the small subset of decoding an Expected string can
// contain. strconv.Unquote handles Go's escape set, which covers every
// escape A.1.5.2 shares with it; anything it rejects is passed through
// with only the quotes stripped rather than dropped.
func unquote(raw string) string {
	if s, err := strconv.Unquote(raw); err == nil {
		return s
	}
	return strings.Trim(raw, "\"`")
}

// TestResult is one test's outcome.
type TestResult struct {
	Case   TestCase
	Passed bool
	Reason string
}

// RunTests loads the input as a `build test` package, runs every test
// whose name contains filter, and reports. It returns false if any test
// failed.
//
// This is the single-package path: opts.Input must itself hold .vs files
// directly (loadDir's own definition of a package). RunTestsAuto is what
// tells this case apart from a suite root; most callers should use that
// instead of calling this directly, unless the input is already known to
// be one package.
func RunTests(opts *Options, filter string) (bool, error) {
	results, err := runPackageTests(opts, filter)
	if err != nil {
		return false, err
	}
	_, failed := report(opts, results)
	return failed == 0, nil
}

// runPackageTests is RunTests minus the pass/fail reduction: it loads
// opts.Input as one `build test` package and returns every discovered
// test's outcome. Factored out so RunTestSuite can run this once per
// discovered package and aggregate the results itself, using the same
// per-line formatting report() already gives a single-package run.
func runPackageTests(opts *Options, filter string) ([]TestResult, error) {
	opts.defaults()

	t, err := ResolveTarget(targetRequest{MinOSVersion: opts.MinOSVersion})
	if err != nil {
		return nil, err
	}
	// `build test` is the only tag that changes what is grammatical, so a
	// test run loads under it rather than under the host's own tag.
	t.Tag = token.TagTest

	lc := newLoadContext()
	pkgs, loadErr := loadForTest(opts, lc, t)
	if pkgs == nil && loadErr != nil {
		return nil, loadErr
	}

	root := pkgs[len(pkgs)-1]
	cases := DiscoverTests(root)
	if len(cases) == 0 {
		return nil, fmt.Errorf(
			"%s declares no `test`-marked functions (a test file needs `build test`)", opts.Input)
	}

	return runCases(opts, t, pkgs, cases, lc, filter)
}

// loadForTest is Load, minus the "abort on any diagnostic" behavior: an
// Expected(error) test *produces* diagnostics, and the runner has to look
// at them before deciding whether the load failed.
func loadForTest(opts *Options, lc *loadContext, t Target) ([]*Package, error) {
	saved := *opts
	saved.Stderr = discardWriter{} // diagnostics are classified, not printed
	return Load(&saved, t.Tag)
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }

// runCases evaluates each test. Error tests are settled from the
// diagnostics already collected; value and run tests are lowered, built,
// and executed one at a time.
func runCases(opts *Options, t Target, pkgs []*Package, cases []TestCase, lc *loadContext, filter string) ([]TestResult, error) {
	// Partition the collected diagnostics by which test function's extent
	// they fall inside. A diagnostic outside every error test is a real
	// compile failure and stops the run: nothing downstream would be
	// trustworthy.
	inside := map[string][]*diag.Diagnostic{}
	var stray []*diag.Diagnostic

	for _, d := range lc.list.Items() {
		owner := ""
		for _, c := range cases {
			if c.Kind == TestError && d.Pos >= c.Pos && d.Pos < c.End {
				owner = c.Name
				break
			}
		}
		if owner == "" {
			if d.Sev == diag.Error {
				stray = append(stray, d)
			}
			continue
		}
		inside[owner] = append(inside[owner], d)
	}

	if len(stray) > 0 {
		for _, d := range stray {
			_ = lc.renderer(false).Render(opts.Stderr, d)
		}
		return nil, fmt.Errorf(
			"%d compile error(s) outside any Expected(error) test — fix these before the suite can run",
			len(stray))
	}

	var results []TestResult
	for _, c := range cases {
		if filter != "" && !strings.Contains(c.Name, filter) {
			continue
		}
		switch c.Kind {
		case TestError:
			results = append(results, judgeErrorCase(c, inside[c.Name]))
		default:
			results = append(results, runValueCase(opts, t, pkgs, c))
		}
	}
	return results, nil
}

func judgeErrorCase(c TestCase, ds []*diag.Diagnostic) TestResult {
	if len(ds) == 0 {
		return TestResult{Case: c, Reason: "expected a compile error, but the function compiled"}
	}
	if !c.HasText {
		return TestResult{Case: c, Passed: true}
	}
	for _, d := range ds {
		// Text() deliberately excludes position, severity, and code — it
		// is exactly the normative string A.12.2 compares against.
		if d.Text() == c.Expected {
			return TestResult{Case: c, Passed: true}
		}
	}
	return TestResult{Case: c, Reason: fmt.Sprintf(
		"expected the diagnostic %q, got %q", c.Expected, ds[0].Text())}
}

func runValueCase(opts *Options, t Target, pkgs []*Package, c TestCase) TestResult {
	modules, root, err := Lower(opts, t, pkgs, LowerOptions{TestFunc: c.Name})
	if err != nil {
		return TestResult{Case: c, Reason: fmt.Sprintf("lowering failed: %v", err)}
	}
	bin, err := buildImage(t, modules, root)
	if err != nil {
		return TestResult{Case: c, Reason: fmt.Sprintf("build failed: %v", err)}
	}
	stdout, stderr, code, err := captureOutput(bin)
	if err != nil {
		return TestResult{Case: c, Reason: fmt.Sprintf("execution failed: %v", err)}
	}
	if code != 0 {
		reason := fmt.Sprintf("exited with status %d", code)
		if len(stderr) > 0 {
			reason += ": " + strings.TrimSpace(string(stderr))
		}
		return TestResult{Case: c, Reason: reason}
	}
	if c.Kind == TestRuns {
		return TestResult{Case: c, Passed: true}
	}

	// A value test's result is compared against the auto-emitted rendered
	// form, so what's compared is the program's whole stdout, trailing
	// newline removed.
	got := strings.TrimRight(string(stdout), "\r\n")
	if got == c.Expected {
		return TestResult{Case: c, Passed: true}
	}
	return TestResult{Case: c, Reason: fmt.Sprintf("expected %q, got %q", c.Expected, got)}
}

// report prints one package's results — ok lines only under -v, every FAIL
// line always — followed by its own pass/fail/total line, and returns the
// counts so a caller aggregating several packages (RunTestSuite) doesn't
// have to re-walk results to total them.
func report(opts *Options, results []TestResult) (passed, failed int) {
	for _, r := range results {
		if r.Passed {
			passed++
			if opts.Verbose {
				fmt.Fprintf(opts.Stdout, "ok    %s\n", r.Case.Name)
			}
			continue
		}
		failed++
		fmt.Fprintf(opts.Stdout, "FAIL  %s — %s\n", r.Case.Name, r.Reason)
	}
	fmt.Fprintf(opts.Stdout, "\n%d passed, %d failed, %d total\n", passed, failed, len(results))
	return passed, failed
}

// ---------------------------------------------------------------------------
// suite discovery
// ---------------------------------------------------------------------------

// RunTestsAuto decides between RunTests and RunTestSuite from opts.Input's
// own shape, so neither the CLI nor a caller has to say up front whether
// it's pointing at one package or a tree of them:
//
//   - A file always goes through RunTests, matching Load's own two-path
//     split (loadFile never treats a file as a directory).
//   - A directory holding .vs files directly is one package (loadDir's own
//     definition) and runs through RunTests unchanged.
//   - A directory holding none is a suite root only if some subdirectory
//     does; RunTestSuite runs every one of those and aggregates.
func RunTestsAuto(opts *Options, filter string) (bool, error) {
	opts.defaults()

	info, err := os.Stat(opts.Input)
	if err != nil {
		return false, fmt.Errorf("reading %s: %w", opts.Input, err)
	}
	if !info.IsDir() {
		return RunTests(opts, filter)
	}

	has, err := dirHasVsFiles(opts.Input)
	if err != nil {
		return false, err
	}
	if has {
		return RunTests(opts, filter)
	}
	return RunTestSuite(opts, filter)
}

// RunTestSuite runs every test package discovered under opts.Input —
// every directory in the tree that holds .vs files directly — as an
// independent RunTests-style run, and aggregates pass/fail counts across
// all of them. It returns false if any package failed to run or any test
// within any package failed.
//
// Each package is fully isolated: its own Load, its own per-test-function
// Lower/build/run cycle (runPackageTests), so a compile error in one
// package's tests cannot mask or block another's.
func RunTestSuite(opts *Options, filter string) (bool, error) {
	opts.defaults()

	dirs, err := discoverTestPackageDirs(opts.Input)
	if err != nil {
		return false, fmt.Errorf("walking %s: %w", opts.Input, err)
	}
	if len(dirs) == 0 {
		return false, fmt.Errorf(
			"%s holds no test package — no directory under it carries .vs files directly", opts.Input)
	}

	allOK := true
	var totalPassed, totalFailed int
	for _, dir := range dirs {
		sub := *opts
		sub.Input = dir

		fmt.Fprintf(opts.Stdout, "%s:\n", dir)
		results, err := runPackageTests(&sub, filter)
		if err != nil {
			fmt.Fprintf(opts.Stdout, "  error: %v\n\n", err)
			allOK = false
			continue
		}
		passed, failed := report(opts, results)
		totalPassed += passed
		totalFailed += failed
		if failed > 0 {
			allOK = false
		}
		fmt.Fprintln(opts.Stdout)
	}

	fmt.Fprintf(opts.Stdout, "=== %d passed, %d failed, %d total across %d package(s) ===\n",
		totalPassed, totalFailed, totalPassed+totalFailed, len(dirs))
	return allOK, nil
}

// discoverTestPackageDirs walks root and returns every directory — root
// included — that holds at least one .vs file directly, sorted so a suite
// run is byte-reproducible in the order it reports packages. A directory
// whose name starts with '.' is skipped entirely (its own name, not its
// path, so "./.git" is skipped but a root the caller literally names
// "./." is still walked): the same convention every other tree-walking
// tool assumes for VCS and metadata directories, and nothing here needs a
// flag to opt out of it.
func discoverTestPackageDirs(root string) ([]string, error) {
	root = filepath.Clean(root)
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if path != root && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		has, err := dirHasVsFiles(path)
		if err != nil {
			return err
		}
		if has {
			out = append(out, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

// dirHasVsFiles reports whether dir holds at least one .vs file directly —
// the same test loadDir itself effectively applies by handing every such
// file to parser.ParseDir. Subdirectories don't count: a package is
// directory-granular, never recursive (the same rule ast.NewPackage and
// importer.Load already enforce).
func dirHasVsFiles(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".vs") {
			return true, nil
		}
	}
	return false, nil
}