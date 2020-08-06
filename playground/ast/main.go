package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"strings"
)

// 生成测试模板
// 每个文件只有一个struct
// 为struct的方法生成测试文档
type fn struct {
	Name    string
	Params  []string
	Results []string
}
type document struct {
	Imports    []string
	StructName string
	Func       []*fn
}

const tmp = `
package test

import "testing"
import "github.com/stretchr/testify/assert"
{{ range $v := .Imports}}import "{{$v}}"
{{ end }}
{{ $structName := .StructName }}
func New{{$structName}}() *{{$structName}} {
	return &{{$structName}}{

	}
}

{{ range $f := .Func}}
func Test{{$structName}}{{$f.Name}}(t *testing.T){
	service := New{{$structName}}()
	{{ $l := len $f.Results }}
	{{if eq $l 2}} _, {{ else }} _,_, {{end}}err := service.{{$f.Name}}(context.Background(),)
	assert.NoError(t, err)
}
{{ end }}
`

func main() {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "./service/demo.go", nil, 0)
	if err != nil {
		panic(err)
	}
	doc := document{}
	// Print the AST.
	for _, decl := range f.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			gend := decl.(*ast.GenDecl)
			for _, spec := range gend.Specs {
				if impspec, ok := spec.(*ast.ImportSpec); ok {
					log.Println("import", impspec.Path.Value)
					doc.Imports = append(doc.Imports, strings.Trim(impspec.Path.Value, `"`))
					continue
				}
				if typspce, ok := spec.(*ast.TypeSpec); ok {
					log.Println("type", typspce.Name)
					doc.StructName = typspce.Name.String()
					continue
				}
			}
		case *ast.FuncDecl:
			f := &fn{}
			fund := decl.(*ast.FuncDecl)
			log.Println("func name", fund.Name)
			for _, recv := range fund.Recv.List {
				log.Println("recv", recv.Type.(*ast.StarExpr).X.(*ast.Ident))
				if recv.Type.(*ast.StarExpr).X.(*ast.Ident).String() == doc.StructName {
					f.Name = fund.Name.Name
				}
			}
			if f.Name == "" {
				log.Println("not", doc.StructName, "func")
				continue
			}
			for _, p := range fund.Type.Params.List {
				switch p.Type.(type) {
				case *ast.StarExpr:
					star := p.Type.(*ast.StarExpr)
					if sel, ok := star.X.(*ast.SelectorExpr); ok {
						log.Println("param", sel.X, sel.Sel)
						f.Params = append(f.Params, fmt.Sprintf("%v.%v", sel.X, sel.Sel))
					}
				case *ast.Ident:
					ident := p.Type.(*ast.Ident)
					log.Println("param", ident)
					f.Params = append(f.Params, fmt.Sprintf("%v", ident))
				case *ast.SelectorExpr:
					sel := p.Type.(*ast.SelectorExpr)
					log.Println("param", sel.X, sel.Sel)
					f.Params = append(f.Params, fmt.Sprintf("%v.%v", sel.X, sel.Sel))
				}
			}
			for _, p := range fund.Type.Results.List {
				switch p.Type.(type) {
				case *ast.StarExpr:
					star := p.Type.(*ast.StarExpr)
					if sel, ok := star.X.(*ast.SelectorExpr); ok {
						log.Println("result", sel.X, sel.Sel)
						f.Results = append(f.Results, fmt.Sprintf("%v.%v", sel.X, sel.Sel))
					}
				case *ast.Ident:
					ident := p.Type.(*ast.Ident)
					log.Println("result", ident)
					f.Results = append(f.Results, fmt.Sprintf("%v", ident))
				case *ast.SelectorExpr:
					sel := p.Type.(*ast.SelectorExpr)
					log.Println("result", sel.X, sel.Sel)
					f.Results = append(f.Results, fmt.Sprintf("%v.%v", sel.X, sel.Sel))
				}
			}
			doc.Func = append(doc.Func, f)
		}
	}
	log.Println("###")
	log.Println("doc of", doc.StructName)
	log.Println("import", doc.Imports)
	for _, f := range doc.Func {
		log.Println("func", f.Name)
		log.Println("params", f.Params)
		log.Println("results", f.Results)
	}
	log.Println("###")
	tmpl, err := template.New("test").Parse(tmp)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, doc)
	if err != nil {
		panic(err)
	}
}
