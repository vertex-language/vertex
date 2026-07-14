package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vertex-language/ir/vertex/ast"
	"github.com/vertex-language/ir/vertex/parser"
)

type TestCase struct {
	File           string
	FuncName       string
	ExpectedType   string
	ExpectedValue  string
	WantCompileErr bool // true for `test -> Expected(error)`
	Decl           *ast.FuncDecl
	OrigFile       *ast.File
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
	f, err := parser.ParseFile(path, src)
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
		typeName, value, wantErr, ok := extractExpected(fd)
		if !ok {
			return nil, fmt.Errorf(
				"test function %q: return type must be Expected(Type, \"value\") or Expected(error)",
				fd.Name.Name,
			)
		}
		// Turn the test function into an ordinary function so it
		// compiles as normal code. A compile-error test keeps its
		// declared (possibly absent) result as-is — there's no
		// synthetic return value to compare, since the point is that
		// this function shouldn't type-check at all.
		fd.Test = nil
		if !wantErr {
			fd.Result = &ast.NamedType{Name: []*ast.Ident{{Name: typeName}}}
		}

		cases = append(cases, TestCase{
			File:           path,
			FuncName:       fd.Name.Name,
			ExpectedType:   typeName,
			ExpectedValue:  value,
			WantCompileErr: wantErr,
			Decl:           fd,
			OrigFile:       f,
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

// extractExpected pulls the expectation out of a test function's
// `test -> Expected(...)` clause. Two forms are recognized:
//
//   - Expected(Type, "value")  — run the function, compare stdout.
//   - Expected(error)          — the function's body must fail to
//     compile; wantErr is true and typeName/value are unused.
func extractExpected(fd *ast.FuncDecl) (typeName, value string, wantErr bool, ok bool) {
	if fd.Test == nil || fd.Test.Expect == nil {
		return "", "", false, false
	}
	call, ok := fd.Test.Expect.(*ast.CallExpr)
	if !ok {
		return "", "", false, false
	}
	fun, ok := call.Fun.(*ast.Ident)
	if !ok || fun.Name != "Expected" {
		return "", "", false, false
	}

	if len(call.Args) == 1 {
		id, ok := call.Args[0].Value.(*ast.Ident)
		if !ok || id.Name != "error" {
			return "", "", false, false
		}
		return "", "", true, true
	}

	if len(call.Args) != 2 {
		return "", "", false, false
	}
	typeIdent, ok := call.Args[0].Value.(*ast.Ident)
	if !ok {
		return "", "", false, false
	}
	lit, ok := call.Args[1].Value.(*ast.BasicLit)
	if !ok || lit.Kind != ast.LitString {
		return "", "", false, false
	}
	raw := strings.Trim(lit.Value, `"`)
	return typeIdent.Name, raw, false, true
}

func testDumpBase(tc TestCase) string {
	base := strings.TrimSuffix(filepath.Base(tc.File), ".vs")
	return base + "_" + tc.FuncName
}