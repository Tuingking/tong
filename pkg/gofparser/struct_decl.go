// Package contain function to parse go file
// Feature:
// 1. Extract struct from go file

package gofparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GoStruct struct {
	Name   string
	Fields []GoStructField
}

type GoStructField struct {
	Name string
	Type string
	Tag  string
}

func (g GoStructField) ParseTag() (*structtag.Tags, error) {
	tags, err := structtag.Parse(g.Tag)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (g GoStructField) ParseTagAndGet(tagName string) string {
	tags, err := structtag.Parse(g.Tag)
	if err != nil {
		logrus.Error(errors.Wrap(err, "parse struct tag"))
		return ""
	}

	tag, err := tags.Get(tagName)
	if err != nil {
		logrus.Error(errors.Wrap(err, "get tag `orm`"))
		return ""
	}

	return tag.Name
}

// GetStructDecl return list of struct declaration on a *.go file
func GetStructDecl(filepath string) ([]GoStruct, error) {
	src, err := os.Open(filepath)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "open go file"))
	}

	fs := token.NewFileSet()

	f, err := parser.ParseFile(fs, "", src, parser.AllErrors)
	if err != nil {
		return []GoStruct{}, errors.Wrap(err, "parse file")
	}

	v := &structVisitor{
		structs: []GoStruct{},
	}
	ast.Walk(v, f)

	return v.structs, nil
}

type structVisitor struct {
	depth   int        // depth of the tree
	visited []ast.Node // store visited node
	structs []GoStruct
}

// walk through tree, if find struct then store inside `structs` field
// example tree of struct declaration:
// *ast.GenDecl
// 		*ast.TypeSpec						<-- get struct name
// 			*ast.Ident
// 			*ast.StructType 				<-- get struct fields
// 				*ast.FieldList
// 					*ast.Field
// 						*ast.Ident
// 						*ast.Ident			<-- *ast.Ident if datatype of the field is `primitive` (e.g. int, int64, float, string, etc)
// 						*ast.BasicLit
// 					*ast.Field
// 						*ast.Ident
// 						*ast.MapType		<-- *ast.MapType if datatype of the field is `map` (e.g. map[string]int)
// 							*ast.Ident
// 							*ast.Ident
// 					*ast.Field
// 						*ast.Ident
// 						*ast.SelectorExpr	<-- *ast.SelectorExpr if datatype of the field is `struct` (e.g time.Time, sql.NullTime)
// 							*ast.Ident
// 							*ast.Ident
// 						*ast.BasicLit
func (v *structVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return nil
	}

	// since we want to extract the struct,
	// so we only care if visitor visit `*ast.TypeSpec` and `*ast.StructType` as explain above
	switch d := node.(type) {
	case *ast.TypeSpec:
		v.visited = append(v.visited, node)
	case *ast.StructType:
		lastVisited := v.visited[len(v.visited)-1]

		// if lastVisited is TypeSpec means this node belongs to that struct
		if ts, ok := lastVisited.(*ast.TypeSpec); ok {
			structName := ts.Name.Name
			// logrus.Info("%s🐶[depth=%d][%s]\n", strings.Repeat("\t", int(v.depth)), v.depth, structName)

			fields := []GoStructField{}
			for _, field := range d.Fields.List {
				var fName string
				if len(field.Names) > 0 {
					fName = field.Names[0].Name
				}

				var fTyp string
				switch x := field.Type.(type) {
				case *ast.SelectorExpr:
					fTyp = fmt.Sprintf("%s.%s", x.X, x.Sel)
				case *ast.MapType:
					fTyp = fmt.Sprintf("map[%v]%v", x.Key, x.Value)
				case *ast.Ident:
					fTyp = x.Name
				}

				var fTag string
				if field.Tag != nil {
					fTag = strings.ReplaceAll(field.Tag.Value, "`", "")
				}

				fields = append(fields, GoStructField{
					Name: fName,
					Type: fTyp,
					Tag:  fTag,
				})
			}

			v.structs = append(v.structs, GoStruct{
				Name:   structName,
				Fields: fields,
			})
		}
	}

	v.depth += 1
	return v
}
