package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vertex-language/ir/vertex/ast"
)

type TestCase struct {
	File          string
	FuncName      string
	ExpectedType  string
	ExpectedValue string
	Decl          *ast.FuncDecl
	OrigFile      *ast.File
}

func collectTestFiles(dir, file string) ([]string, error) {
	if file != "" {
		return []string{file}, nil
	}
	var out []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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

func parseTestCases(path string) ([]TestCase, error) {
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

	var cases []TestCase
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Test == nil {
			continue
		}
		typeName, value, ok := extractExpected(fd)
		if !ok {
			return nil, fmt.Errorf(
				"test function %q: return type must be Expected(Type, \"value\")",
				fd.Name.Name,
			)
		}
		// Turn the test function into an ordinary function with a
		// concrete return type so it compiles as normal code.
		fd.Test = nil
		fd.Result = &ast.NamedType{Name: []*ast.Ident{{Name: typeName}}}

		cases = append(cases, TestCase{
			File:          path,
			FuncName:      fd.Name.Name,
			ExpectedType:  typeName,
			ExpectedValue: value,
			Decl:          fd,
			OrigFile:      f,
		})
	}
	return cases, nil
}

func hasBuildTest(f *ast.File) bool {
	for _, b := range f.Builds {
		if b.Name == "test" {
			return true
		}
	}
	return false
}

// extractExpected pulls the type name and expected string value out of
// a test function's `test -> Expected(Type, "value")` clause. The
// clause's Expect field is an ordinary expression: a call to Expected
// with the type name spelled as a bare identifier and the value as a
// string literal.
func extractExpected(fd *ast.FuncDecl) (typeName, value string, ok bool) {
	if fd.Test == nil || fd.Test.Expect == nil {
		return "", "", false
	}
	call, ok := fd.Test.Expect.(*ast.CallExpr)
	if !ok {
		return "", "", false
	}
	fun, ok := call.Fun.(*ast.Ident)
	if !ok || fun.Name != "Expected" || len(call.Args) != 2 {
		return "", "", false
	}
	typeIdent, ok := call.Args[0].Value.(*ast.Ident)
	if !ok {
		return "", "", false
	}
	lit, ok := call.Args[1].Value.(*ast.BasicLit)
	if !ok || lit.Kind != ast.LitString {
		return "", "", false
	}
	raw := strings.Trim(lit.Value, `"`)
	return typeIdent.Name, raw, true
}

func testDumpBase(tc TestCase) string {
	base := strings.TrimSuffix(filepath.Base(tc.File), ".vs")
	return base + "_" + tc.FuncName
}