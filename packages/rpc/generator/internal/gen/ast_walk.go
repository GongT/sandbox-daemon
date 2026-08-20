package gen

import (
	"go/ast"
	"go/token"
)

func findTypes(file *ast.File) []string {
	names := []string{}

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok := ts.Type.(*ast.StructType); ok {
				names = append(names, ts.Name.Name) // type XXX struct{}
			}
		}
	}

	return names
}

func methodsOf(file *ast.File, typeName string) []*ast.FuncDecl {
	var out []*ast.FuncDecl

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue // not a method
		}

		recvType := validateReceiver(fn.Recv.List[0].Type)
		if recvType == typeName {
			out = append(out, fn)
		}
	}

	return out
}

func validateReceiver(expr ast.Expr) string {
	switch t := expr.(type) {
	// case *ast.Ident: // no this
	case *ast.StarExpr:
		if id, ok := t.X.(*ast.Ident); ok {
			if id.IsExported() {
				return id.Name
			}
		}
	}
	return ""
}
