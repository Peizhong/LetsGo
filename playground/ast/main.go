package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "./service/demo.go", nil, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	for _, decl := range f.Decls {
		if gend, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range gend.Specs {
				if impspec, ok := spec.(*ast.ImportSpec); ok {
					log.Println("import", impspec.Path.Value)
					continue
				}
				if typspce, ok := spec.(*ast.TypeSpec); ok {
					log.Println("type", typspce.Name)
					continue
				}
			}
			continue
		}
		if fund, ok := decl.(*ast.FuncDecl); ok {
			log.Println("name", fund.Name)
			for _, p := range fund.Type.Params.List {
				log.Println(p.Type)
				if star, ok := p.Type.(*ast.StarExpr); ok {
					log.Println("param", star.X)
				}
			}
			for _, p := range fund.Type.Results.List {
				log.Println(p.Type)
				if star, ok := p.Type.(*ast.StarExpr); ok {
					log.Println("result", star.X)
				}
			}
			for _, recv := range fund.Recv.List {
				log.Println("recv", recv.Type.(*ast.StarExpr).X)
			}
		}
	}

}
