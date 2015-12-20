package gengo

import (
	"strings"
)

func GoUpperName(s string) string {
	s = SnakeToUpperCamel(s)
	if strings.HasSuffix(s, "Id") {
		s = strings.TrimSuffix(s, "Id") + "ID"
	}
	return s
}

func SnakeToUpperCamel(s string) string {
	ss := strings.Split(s, "_")
	for i := range ss {
		ss[i] = strings.Title(ss[i])
	}
	return strings.Join(ss, "")
}

func SnakeToLowerCamel(s string) string {
	ss := strings.Split(s, "_")
	for i := 1; i < len(ss); i++ {
		ss[i] = strings.Title(ss[i])
	}
	return strings.Join(ss, "")
}
