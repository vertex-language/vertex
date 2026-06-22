// testrun.go
package driver

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vertex-language/ir/vertex/ast"
)

// testCase is one discovered test function with its extracted metadata.
type testCase struct {
	file          string        // source file the function came from
	funcName      string        // e.g. "test_bool_true"
	expectedType  string        // e.g. "bool", "int32", "float64"
	expectedValue string        // e.g. "1", "42", "3.140000"
	decl          *ast.FuncDecl // mutated in place: IsTest=false, Return=plain type
	origFile      *ast.File     // the parsed source file (for helper decls)
}

// runTests is the entry point called by driver.Run for modeTest.
func runTests(cfg config, stderr io.Writer) int {
	files, err := collectTestFiles(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "vertex: %v\n", err)
		return 1
	}
	if len(files) == 0 {
		fmt.Fprintf(stderr, "vertex: no test files found\n")
		return 1
	}

	var cases []testCase
	for _, path := range files {
		tc, err := parseTestCases(path)
		if err != nil {
			fmt.Fprintf(stderr, "vertex: %s: %v\n", path, err)
			return 1
		}
		cases = append(cases, tc...)
	}
	if len(cases) == 0 {
		fmt.Fprintf(stderr, "vertex: no test functions found\n")
		return 1
	}

	pass, fail := 0, 0
	for _, tc := range cases {
		if execTest(tc, cfg, stderr) {
			pass++
		} else {
			fail++
		}
	}

	fmt.Fprintf(os.Stdout, "\n%d passed, %d failed\n", pass, fail)
	if fail > 0 {
		return 1
	}
	return 0
}

// collectTestFiles finds all .vs files that carry a `build test` tag.
// It does a quick byte scan before committing to a full ANTLR parse.
func collectTestFiles(cfg config) ([]string, error) {
	if cfg.testFile != "" {
		return []string{cfg.testFile}, nil
	}
	var out []string
	err := filepath.WalkDir(cfg.testDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".vs") {
			return nil
		}
		if quickHasBuildTest(path) {
			out = append(out, path)
		}
		return nil
	})
	return out, err
}

// quickHasBuildTest scans the first 512 bytes of a file looking for a line
// that is exactly "build test". This avoids a full ANTLR parse during discovery.
func quickHasBuildTest(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	end := len(data)
	if end > 512 {
		end = 512
	}
	for _, line := range strings.Split(string(data[:end]), "\n") {
		if strings.TrimSpace(line) == "build test" {
			return true
		}
	}
	return false
}

// parseTestCases fully parses one test file and returns a testCase for every
// function decorated with the `test` qualifier. The FuncDecl is mutated in
// place: IsTest is cleared and the Expected(T,"v") return type is unwrapped
// to plain T so the existing VIR lowerer sees a normal function.
func parseTestCases(path string) ([]testCase, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := ast.NewFile(path, src)
	if err != nil {
		return nil, err
	}
	if !hasBuildTest(f) {
		return nil, nil
	}

	var cases []testCase
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || !fd.IsTest {
			continue
		}
		typeName, value, ok := extractExpected(fd)
		if !ok {
			return nil, fmt.Errorf(
				"test function %q: return type must be Expected(Type, \"value\"), got %T",
				fd.Name, fd.Return,
			)
		}
		// Mutate in place so the VIR lowerer sees a plain function.
		fd.IsTest = false
		fd.Return = &ast.NamedType{Name: typeName}

		cases = append(cases, testCase{
			file:          path,
			funcName:      fd.Name,
			expectedType:  typeName,
			expectedValue: value,
			decl:          fd,
			origFile:      f,
		})
	}
	return cases, nil
}

// hasBuildTest returns true when the parsed file has a `build test` tag.
func hasBuildTest(f *ast.File) bool {
	for _, b := range f.Build {
		if b.Tag == "test" {
			return true
		}
	}
	return false
}

// extractExpected pulls the type name and expected stdout string out of
// Expected(TypeName, "value") in the return clause of a test FuncDecl.
func extractExpected(fd *ast.FuncDecl) (typeName, value string, ok bool) {
	nt, ok := fd.Return.(*ast.NamedType)
	if !ok || nt.Name != "Expected" || len(nt.CtorArgs) != 2 {
		return "", "", false
	}
	typeArg, ok := nt.CtorArgs[0].Type.(*ast.NamedType)
	if !ok {
		return "", "", false
	}
	if !nt.CtorArgs[1].IsString {
		return "", "", false
	}
	// CtorArgs[1].String is the raw lexeme including surrounding quotes — strip them.
	raw := strings.Trim(nt.CtorArgs[1].String, `"`)
	return typeArg.Name, raw, true
}

// execTest compiles and runs one test case, returning true on pass.
func execTest(tc testCase, cfg config, stderr io.Writer) bool {
	label := fmt.Sprintf("%s::%s", filepath.Base(tc.file), tc.funcName)

	tmpDir, err := os.MkdirTemp("", "vertex-test-*")
	if err != nil {
		fmt.Fprintf(os.Stdout, "FAIL  %s  (tmpdir: %v)\n", label, err)
		return false
	}
	defer os.RemoveAll(tmpDir)

	binPath := filepath.Join(tmpDir, tc.funcName)
	if isWindowsTarget(cfg.target) {
		binPath += ".exe"
	}

	pkg := buildSyntheticPackage(tc)

	compileCfg := config{
		output:      binPath,
		target:      cfg.target,
		packagesDir: cfg.packagesDir,
		sysroot:     cfg.sysroot,
		mode:        modeExe,
		optLevel:    cfg.optLevel,
	}
	var compileErr strings.Builder
	if code := emitPackage(pkg, compileCfg, &compileErr); code != 0 {
		fmt.Fprintf(os.Stdout, "FAIL  %s  (compile error)\n", label)
		fmt.Fprint(os.Stdout, compileErr.String())
		return false
	}

	out, err := exec.Command(binPath).Output()
	if err != nil {
		fmt.Fprintf(os.Stdout, "FAIL  %s  (exec: %v)\n", label, err)
		return false
	}

	got := strings.TrimRight(string(out), "\r\n")
	if got == tc.expectedValue {
		fmt.Fprintf(os.Stdout, "ok    %s\n", label)
		return true
	}
	fmt.Fprintf(os.Stdout, "FAIL  %s\n", label)
	fmt.Fprintf(os.Stdout, "      want: %q\n", tc.expectedValue)
	fmt.Fprintf(os.Stdout, "      got:  %q\n", got)
	return false
}

// buildSyntheticPackage constructs an ast.Package for one test binary.
// It contains:
//   - all non-test decls from the original file (helpers, types, globals)
//   - the mutated test FuncDecl (already has IsTest=false, plain return type)
//   - a synthetic `class C : c` with printf
//   - a synthetic `main()` that calls the test function and prints the result
func buildSyntheticPackage(tc testCase) *ast.Package {
	var decls []ast.Decl

	// Include every decl from the original file except other (still-unmutated)
	// test functions. The target test function has already been mutated to
	// IsTest=false so it passes through the loop naturally.
	for _, d := range tc.origFile.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.IsTest {
			continue // skip test functions that belong to other test cases
		}
		decls = append(decls, d)
	}

	decls = append(decls, syntheticClassC())
	decls = append(decls, syntheticMain(tc))

	f := &ast.File{
		Filename: tc.file,
		Package:  &ast.PackageClause{Name: "main"},
		Decls:    decls,
	}
	return &ast.Package{
		Name:  "main",
		Files: []*ast.File{f},
		Decls: decls,
	}
}

// syntheticClassC returns the AST for:
//
//	class C : c {
//	    func printf(fmt: *const char, args: ...*const char) -> int32
//	}
func syntheticClassC() *ast.ClassDecl {
	charPtr := &ast.PointerType{
		Const: true,
		Elem:  &ast.NamedType{Name: "char"},
	}
	return &ast.ClassDecl{
		Name:   "C",
		Parent: "c",
		Methods: []*ast.MethodSig{
			{
				Name: "printf",
				Params: []*ast.Param{
					{Name: "fmt", Type: charPtr},
					{Name: "args", Variadic: true, Type: charPtr},
				},
				Return: &ast.NamedType{Name: "int32"},
			},
		},
	}
}

// syntheticMain returns the AST for:
//
//	func main() -> int32 {
//	    let result = <funcName>()
//	    let libc   = C()
//	    libc.printf("<fmt>\n", result)
//	    return 0
//	}
func syntheticMain(tc testCase) *ast.FuncDecl {
	fmtLit := &ast.BasicLit{
		Kind:  ast.LitString,
		Value: printfFmt(tc.expectedType),
	}
	return &ast.FuncDecl{
		Name:   "main",
		Return: &ast.NamedType{Name: "int32"},
		Body: &ast.Block{
			Stmts: []ast.Stmt{
				// let result = <funcName>()
				&ast.VarDecl{
					Kind:  ast.BindLet,
					Names: []string{"result"},
					Value: &ast.CallExpr{
						Fun: &ast.Ident{Name: tc.funcName},
					},
				},
				// let libc = C()
				&ast.VarDecl{
					Kind:  ast.BindLet,
					Names: []string{"libc"},
					Value: &ast.CallExpr{
						Fun: &ast.Ident{Name: "C"},
					},
				},
				// libc.printf("<fmt>\n", result)
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "libc"},
							Sel: "printf",
						},
						Args: []*ast.Arg{
							{Value: fmtLit},
							{Value: &ast.Ident{Name: "result"}},
						},
					},
				},
				// return 0
				&ast.ReturnStmt{
					Value: &ast.BasicLit{Kind: ast.LitIntDec, Value: "0"},
				},
			},
		},
	}
}

// printfFmt returns the raw Vertex string literal (quotes included) for
// printing a value of the given type name via printf.
func printfFmt(typeName string) string {
	switch typeName {
	case "float32", "float64":
		return `"%f\n"`
	case "int64", "uint64":
		return `"%lld\n"`
	default: // bool, int, int32, uint32, int8, uint8, int16, uint16
		return `"%d\n"`
	}
}