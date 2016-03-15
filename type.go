package gengo

import (
	"encoding/json"
	"fmt"
	"go/printer"
	"go/token"
	"io"
	"strings"
)

type File struct {
	PackageName string    `json:"package_name,omitempty"`
	Imports     []*Import `json:"imports,omitempty"`
	TypeDecls   TypeDecls `json:"type_decls,omitempty"`
	Doc         string    `json:"doc,omitempty"`
}
type TypeDecls []*TypeDecl

func (f *File) RemoveDecl(name string) *File {
	file := *f
	file.TypeDecls = nil
	for _, decl := range f.TypeDecls {
		if decl.Name != name {
			file.TypeDecls = append(file.TypeDecls, decl)
		}
	}
	return &file
}

type Import struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path"`
	Doc  string `json:"doc,omitempty"`
}

type TypeDecl struct {
	Name string `json:"name,omitempty"`
	Type Type   `json:"type,omitempty"`
	Doc  string `json:"doc,omitempty"`
}

type Kind int

const (
	IdentKind Kind = iota
	StructKind
	ArrayKind
)

func (k Kind) String() string {
	switch k {
	case IdentKind:
		return ""
	case StructKind:
		return "struct"
	case ArrayKind:
		return "[]"
	}
	return ""
}

func (k Kind) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

type Type struct {
	Kind   Kind                   `json:"kind"`
	Ident  string                 `json:"ident,omitemtpy"`
	Fields Fields                 `json:"fields,omitempty"`
	Attrs  map[string]interface{} `json:"attr,omitempty"` // Additional data that can be attached to a type
}

type Fields []*Field

func (d *Type) Set(key string, value interface{}) {
	if d.Attrs == nil {
		d.Attrs = make(map[string]interface{})
	}
	d.Attrs[key] = value
}

func (d *Type) Get(key string) interface{} {
	if d.Attrs == nil {
		return nil
	}
	return d.Attrs[key]
}

type IdentType string

type Field struct {
	Name string `json:"name"`
	Type Type   `json:"type"`
	Tag  Tag    `json:"tag,omitempty"`
	Doc  string `json:"doc,omitempty"`
}

type Tag []*TagPart

func (t *Tag) String() string {
	if t == nil {
		return ""
	}
	ss := make([]string, len(*t))
	for i := range ss {
		ss[i] = (*t)[i].String()
	}
	return strings.Join(ss, " ")
}

type TagPart struct {
	Encoding  string
	Name      string
	Type      string
	Omitted   bool
	OmitEmpty bool
}

func (p *TagPart) String() string {
	if p.Encoding == "" {
		return ""
	}
	segments := []string{p.Name}
	if p.Type != "" {
		segments = append(segments, p.Type)
	}
	if p.OmitEmpty {
		segments = append(segments, "omitempty")
	}
	return fmt.Sprintf(`%s:"%s"`, p.Encoding, strings.Join(segments, ","))
}

func (d *TypeDecl) Marshal(w io.Writer) error {
	return printer.Fprint(w, token.NewFileSet(), d.AST())
}

func (f *File) Marshal(w io.Writer) error {
	return printer.Fprint(w, token.NewFileSet(), f.AST())
}

func (f *File) JSON() string {
	buf, _ := json.MarshalIndent(f, "", "    ")
	return string(buf)
}

func (a TypeDecls) Len() int           { return len(a) }
func (a TypeDecls) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TypeDecls) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (a Fields) Len() int           { return len(a) }
func (a Fields) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Fields) Less(i, j int) bool { return a[i].Name < a[j].Name }
