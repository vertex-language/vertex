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
		Filename: tc.File,
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

func syntheticMain(tc TestCase) *ast.FuncDecl {
	var printfResultArg *ast.Arg
	if tc.ExpectedType == "string" {
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
						Fun: &ast.Ident{Name: tc.FuncName},
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
								Value: printfFmt(tc.ExpectedType),
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