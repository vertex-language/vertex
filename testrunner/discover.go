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

		cases = append(cases, TestCase{
			File:          path,
			FuncName:      fd.Name,
			ExpectedType:  typeName,
			ExpectedValue: value,
			Decl:          fd,
			OrigFile:      f,
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

func testDumpBase(tc TestCase) string {
	base := strings.TrimSuffix(filepath.Base(tc.File), ".vs")
	return base + "_" + tc.FuncName
}