package gengo

import (
	"encoding/json"
	"go/printer"
	"go/token"
	"io"
)

type File struct {
	PackageName string
	TypeDecls   []*TypeDecl
	Doc         string `json:"doc,omitempty"`
}

type TypeDecl struct {
	Name string
	Type Type
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
	Kind   Kind
	Ident  string   `json:"ident,omitemtpy"`
	Fields []*Field `json:"fields,omitempty"`
}

type IdentType string

type Field struct {
	Name string
	Type Type   `json:"type"`
	Tag  *Tag   `json:"tag,omitempty"`
	Doc  string `json:"doc,omitempty"`
}

type Tag struct {
}

func (f *File) Fprint(w io.Writer) error {
	astFile := f.AST()
	return printer.Fprint(w, token.NewFileSet(), astFile)
}

func (f *File) JSON() string {
	buf, _ := json.MarshalIndent(f, "", "    ")
	return string(buf)
}
