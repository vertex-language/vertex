package testrunner

import (
	"github.com/vertex-language/ir/vertex/ast"
)

// buildSyntheticPackage wraps tc's test function in a tiny standalone
// "main" package: the original file's non-test declarations, a
// synthesized `class C` libc binding if the file doesn't already define
// one, and a synthesized main() that calls the test function and prints
// its result for execTest to compare against tc.ExpectedValue.
func buildSyntheticPackage(tc TestCase) *ast.Package {
	var decls []ast.Decl
	hasClassC := false
	for _, d := range tc.OrigFile.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Test != nil {
			continue
		}
		if cd, ok := d.(*ast.ClassDecl); ok && cd.Name.Name == "C" && cd.ABI != nil {
			hasClassC = true
		}
		decls = append(decls, d)
	}
	if !hasClassC {
		decls = append(decls, syntheticClassC())
	}
	decls = append(decls, syntheticMain(tc))

	f := &ast.File{
		Path:    tc.File,
		Package: &ast.Ident{Name: "main"},
		Decls:   decls,
	}
	return &ast.Package{
		Name:  "main",
		Files: []*ast.File{f},
	}
}

func syntheticClassC() *ast.ClassDecl {
	charPtr := &ast.PtrType{
		Elem: &ast.NamedType{Name: []*ast.Ident{{Name: "char"}}},
	}
	return &ast.ClassDecl{
		Name: &ast.Ident{Name: "C"},
		ABI:  &ast.NamedType{Name: []*ast.Ident{{Name: "c"}}},
		Methods: []*ast.MethodDecl{
			{
				Name: &ast.Ident{Name: "printf"},
				Params: []*ast.Param{
					{Label: "fmt", Variadic: true, Type: charPtr},
				},
				Result: &ast.NamedType{Name: []*ast.Ident{{Name: "int32"}}},
			},
		},
	}
}

func syntheticMain(tc TestCase) *ast.FuncDecl {
	var printfResultArg *ast.Arg
	if tc.ExpectedType == "string" {
		printfResultArg = &ast.Arg{
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "result"},
					Sel: &ast.Ident{Name: "c_str"},
				},
			},
		}
	} else {
		printfResultArg = &ast.Arg{Value: &ast.Ident{Name: "result"}}
	}

	return &ast.FuncDecl{
		Name:   &ast.Ident{Name: "main"},
		Result: &ast.NamedType{Name: []*ast.Ident{{Name: "int32"}}},
		Body: &ast.BlockStmt{
			Stmts: []ast.Stmt{
				&ast.BindingDecl{
					Let:   true,
					Names: []*ast.Ident{{Name: "result"}},
					Values: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: tc.FuncName},
						},
					},
				},
				&ast.BindingDecl{
					Let:   true,
					Names: []*ast.Ident{{Name: "libc"}},
					Values: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: "C"},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "libc"},
							Sel: &ast.Ident{Name: "printf"},
						},
						Args: []*ast.Arg{
							{Value: &ast.BasicLit{
								Kind:  ast.LitString,
								Value: printfFmt(tc.ExpectedType),
							}},
							printfResultArg,
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BasicLit{Kind: ast.LitInt, Value: "0"},
					},
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