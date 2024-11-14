/*
* The aim of this file is to parse the given go file that stores
* the schema data into parsed file that will be changed into the
* SQL equivalent syntax.
 */

package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

type Schema struct {
	SchemaName   string
	SchemaFields []struct {
		Name     string
		DataType string
	}
}

func Parse(fileName string) Schema {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, fileName, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Errors: %s", err)
	}

	table := Schema{}

	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.GenDecl); ok {
			if funcDecl.Tok == token.TYPE {
				table.SchemaName = funcDecl.Specs[0].(*ast.TypeSpec).Name.Name
				for _, field := range funcDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
					fields := struct {
						Name     string
						DataType string
					}{}
					fields.Name = field.Names[0].Name
					switch t := field.Type.(type) {
					case *ast.Ident:
						fields.DataType = field.Type.(*ast.Ident).Name
					case *ast.SelectorExpr:
						fields.DataType = t.Sel.Name
					default:
						log.Printf("Unkown data type: %s", t)
						fields.DataType = "unkown"
					}

					table.SchemaFields = append(table.SchemaFields, fields)
				}
			}
		}
		return true
	})

	return table
}
