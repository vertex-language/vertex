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
package driver

import (
	"fmt"
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
func RunTests(opts *Options, filter string) (bool, error) {
	opts.defaults()

	t, err := ResolveTarget(targetRequest{MinOSVersion: opts.MinOSVersion})
	if err != nil {
		return false, err
	}
	// `build test` is the only tag that changes what is grammatical, so a
	// test run loads under it rather than under the host's own tag.
	t.Tag = token.TagTest

	lc := newLoadContext()
	pkgs, loadErr := loadForTest(opts, lc, t)
	if pkgs == nil && loadErr != nil {
		return false, loadErr
	}

	root := pkgs[len(pkgs)-1]
	cases := DiscoverTests(root)
	if len(cases) == 0 {
		return false, fmt.Errorf(
			"%s declares no `test`-marked functions (a test file needs `build test`)", opts.Input)
	}

	results, fatal := runCases(opts, t, pkgs, cases, lc, filter)
	if fatal != nil {
		return false, fatal
	}
	return report(opts, results), nil
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

func report(opts *Options, results []TestResult) bool {
	passed, failed := 0, 0
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
	return failed == 0
}