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
	PackageName string `json:"package_name,omitempty"`
	TypeDecls   []*TypeDecl
	Doc         string `json:"doc,omitempty"`
}

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
	Fields []*Field               `json:"fields,omitempty"`
	Attrs  map[string]interface{} `json:"attr,omitempty"`
}

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
	Name string
	Type Type   `json:"type"`
	Tag  *Tag   `json:"tag,omitempty"`
	Doc  string `json:"doc,omitempty"`
}

type Tag struct {
	Parts []*TagPart
}

func (t *Tag) String() string {
	if t == nil {
		return ""
	}
	ss := make([]string, len(t.Parts))
	for i := range ss {
		ss[i] = t.Parts[i].String()
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

func (f *File) Fprint(w io.Writer) error {
	astFile := f.AST()
	return printer.Fprint(w, token.NewFileSet(), astFile)
}

func (f *File) JSON() string {
	buf, _ := json.MarshalIndent(f, "", "    ")
	return string(buf)
}
