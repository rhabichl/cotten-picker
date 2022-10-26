package models

import (
	"fmt"
	"strings"
)

type Variable struct {
	Name        string
	Type        string
	Annotations []string
	Security    string
}

func (v Variable) Print() string {
	var sb strings.Builder

	if len(v.Annotations) != 0 {
		for _, a := range v.Annotations {
			sb.WriteString(a + "\n")
		}
	}

	// check if the security is set
	var tmp string

	if len(v.Security) != 0 {
		tmp = v.Security + " "
	}

	sb.WriteString(fmt.Sprintf("%s%s %s;", tmp, v.Type, v.Name))

	if v.Security == "private" || v.Security == "protected" {
		sb.WriteString("\n" + makeGetterAndSetter(v.Name, v.Type))
	}

	return sb.String()
}

func makeGetterAndSetter(Name, Type string) string {
	var sb strings.Builder

	// make setter
	sb.WriteString(fmt.Sprintf("\npublic void set%s(%s o) {\n", strings.Title(Name), Type))
	sb.WriteString(fmt.Sprintf("this.%s = o;\n", Name))
	sb.WriteString("}\n")

	// make getter
	sb.WriteString(fmt.Sprintf("\npublic %s get%s() {\n", Type, strings.Title(Name)))
	sb.WriteString(fmt.Sprintf("return this.%s;\n", Name))
	sb.WriteString("}")

	return sb.String()
}
