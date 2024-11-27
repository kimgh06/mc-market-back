package generate

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/parser"
	"go/token"
	"maple/internal/utilities"
)

type LiteralConstantDecl struct {
	Name    *ast.Ident
	Type    *ast.Ident
	Literal *ast.BasicLit
}

func (d LiteralConstantDecl) String() string {
	if d.Literal != nil {
		return fmt.Sprintf("LiteralConstantDecl(%s %s = %s)", d.Name, d.Type, d.Literal.Value)
	} else {
		return fmt.Sprintf("LiteralConstantDecl(%s %s)", d.Name, d.Type)
	}
}

func GetLiteralConstantDecls(filename string, src any) []LiteralConstantDecl {
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filename, src, parser.AllErrors)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to parse file")
		return nil
	}

	var constants []LiteralConstantDecl

	for _, decl := range parsed.Decls {
		switch t := decl.(type) {
		case *ast.GenDecl:
			if t.Tok == token.CONST {
				specs := utilities.Map(t.Specs, func(spec ast.Spec) *ast.ValueSpec {
					return spec.(*ast.ValueSpec)
				})

				for _, spec := range specs {
					typeIdent := spec.Type.(*ast.Ident)

					for i, name := range spec.Names {
						if len(spec.Values) >= i {
							value := spec.Values[i].(*ast.BasicLit)
							constants = append(constants, LiteralConstantDecl{
								Name:    name,
								Type:    typeIdent,
								Literal: value,
							})
						} else {
							constants = append(constants, LiteralConstantDecl{
								Name:    name,
								Type:    typeIdent,
								Literal: nil,
							})
						}
					}
				}
			}
		}
	}

	return constants
}
