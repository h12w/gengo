package gengo

import (
	"go/ast"
	"go/token"
)

func (f *File) AST() *ast.File {
	var decls []ast.Decl
	for _, d := range f.TypeDecls {
		decls = append(decls, d.AST())
	}
	return &ast.File{
		Name:  &ast.Ident{Name: f.PackageName},
		Decls: decls,
	}
}

func (d *TypeDecl) AST() ast.Decl {
	return &ast.GenDecl{Tok: token.TYPE, Specs: []ast.Spec{
		&ast.TypeSpec{
			Name: &ast.Ident{Name: d.Name},
			Type: d.Type.AST(),
		},
	}}
}

func (t *Type) AST() ast.Expr {
	switch t.Kind {
	case IdentKind:
		return &ast.Ident{Name: t.Ident}
	case ArrayKind:
		return &ast.ArrayType{
			Elt: &ast.Ident{Name: t.Ident},
		}
	case StructKind:
		var fields []*ast.Field
		for _, f := range t.Fields {
			fields = append(fields, f.AST())
		}
		return &ast.StructType{
			Fields: &ast.FieldList{List: fields},
		}
	}
	return nil
}

func (f *Field) AST() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{{Name: f.Name}},
		Type:  f.Type.AST(),
		Tag:   tag(f.Tag.String()),
	}
}

func tag(value string) *ast.BasicLit {
	if value == "" {
		return nil
	}
	return &ast.BasicLit{Kind: token.STRING, Value: "`" + value + "`"}
}
