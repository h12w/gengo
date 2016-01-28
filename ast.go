package gengo

import (
	"go/ast"
	"go/token"
)

func (f *File) AST() *ast.File {
	imports := make([]ast.Spec, len(f.Imports))
	for i := range f.Imports {
		imports[i] = f.Imports[i].AST()
	}
	importDecl := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: imports,
	}
	decls := []ast.Decl{importDecl}
	for _, d := range f.TypeDecls {
		decls = append(decls, d.AST())
	}
	return &ast.File{
		Name:  &ast.Ident{Name: f.PackageName},
		Decls: decls,
	}
}

func (im *Import) AST() *ast.ImportSpec {
	spec := &ast.ImportSpec{
		Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + im.Path + `"`},
	}
	if im.Name != "" {
		spec.Name = &ast.Ident{Name: im.Name}
	}
	if im.Doc != "" {
		spec.Comment = &ast.CommentGroup{List: []*ast.Comment{{Text: im.Doc}}}
	}
	return spec
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

func (f *Field) names() []*ast.Ident {
	if f.Name == f.Type.Ident {
		return nil
	}
	return []*ast.Ident{{Name: f.Name}}
}

func (f *Field) AST() *ast.Field {
	return &ast.Field{
		Names: f.names(),
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
