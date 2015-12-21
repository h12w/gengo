package gengo

import (
	"strings"
)

func SnakeToUpperCamel(s string) string {
	ss := strings.Split(s, "_")
	for i := range ss {
		ss[i] = upperOrTitle(ss[i])
	}
	return strings.Join(ss, "")
}

func SnakeToLowerCamel(s string) string {
	ss := strings.Split(s, "_")
	for i := 1; i < len(ss); i++ {
		ss[i] = upperOrTitle(ss[i])
	}
	return strings.Join(ss, "")
}

var Abbreviations = map[string]struct{}{
	"cpi": struct{}{},
	"gif": struct{}{},
	"id":  struct{}{},
}

func upperOrTitle(s string) string {
	if _, ok := Abbreviations[s]; ok {
		return strings.ToUpper(s)
	}
	return strings.Title(s)
}
