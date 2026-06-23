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

type testCase struct {
	file          string
	funcName      string
	expectedType  string
	expectedValue string
	decl          *ast.FuncDecl
	origFile      *ast.File
}

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
				"test function %q: return type must be Expected(Type, \"value\")",
				fd.Name,
			)
		}
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

func hasBuildTest(f *ast.File) bool {
	for _, b := range f.Build {
		if b.Tag == "test" {
			return true
		}
	}
	return false
}

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
	raw := strings.Trim(nt.CtorArgs[1].String, `"`)
	return typeArg.Name, raw, true
}

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

func buildSyntheticPackage(tc testCase) *ast.Package {
	var decls []ast.Decl
	hasClassC := false
	for _, d := range tc.origFile.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.IsTest {
			continue
		}
		if cd, ok := d.(*ast.ClassDecl); ok && cd.Name == "C" && cd.Parent != "" {
			hasClassC = true
		}
		decls = append(decls, d)
	}
	if !hasClassC {
		decls = append(decls, syntheticClassC())
	}
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
					{Name: "fmt", Variadic: true, Type: charPtr},
				},
				Return: &ast.NamedType{Name: "int32"},
			},
		},
	}
}

func syntheticMain(tc testCase) *ast.FuncDecl {
	// For string results, pass result.c_str() to printf instead of result
	// directly — printf expects a null-terminated *const char for %s.
	var printfResultArg *ast.Arg
	if tc.expectedType == "string" {
		printfResultArg = &ast.Arg{
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "result"},
					Sel: "c_str",
				},
			},
		}
	} else {
		printfResultArg = &ast.Arg{Value: &ast.Ident{Name: "result"}}
	}

	return &ast.FuncDecl{
		Name:   "main",
		Return: &ast.NamedType{Name: "int32"},
		Body: &ast.Block{
			Stmts: []ast.Stmt{
				&ast.VarDecl{
					Kind:  ast.BindLet,
					Names: []string{"result"},
					Value: &ast.CallExpr{
						Fun: &ast.Ident{Name: tc.funcName},
					},
				},
				&ast.VarDecl{
					Kind:  ast.BindLet,
					Names: []string{"libc"},
					Value: &ast.CallExpr{
						Fun: &ast.Ident{Name: "C"},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "libc"},
							Sel: "printf",
						},
						Args: []*ast.Arg{
							{Value: &ast.BasicLit{
								Kind:  ast.LitString,
								Value: printfFmt(tc.expectedType),
							}},
							printfResultArg,
						},
					},
				},
				&ast.ReturnStmt{
					Value: &ast.BasicLit{Kind: ast.LitIntDec, Value: "0"},
				},
			},
		},
	}
}

func printfFmt(typeName string) string {
	switch typeName {
	case "float32", "float64":
		return `"%f\n"`
	case "int64", "uint64":
		return `"%lld\n"`
	case "string":
		return `"%s\n"`
	case "char":
		return `"%c\n"`
	default:
		return `"%d\n"`
	}
}