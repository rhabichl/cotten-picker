package models

import (
	"fmt"
	"strings"
)

type Class struct {
	Vars      []Variable
	Function  []Function
	Name      string
	Import    []Import
	IsBetween bool
}

// takes in an array of printable objects and the classname plus package name
func GenerateClass(Name, Package string, content, imports []Printable) string {
	var sb strings.Builder

	// print the name of the package
	sb.WriteString(fmt.Sprintf("package %s.models;\n", Package))

	// write the imports
	sb.WriteString(getPrintableString(imports))

	// write the class
	sb.WriteString(fmt.Sprintf("@Entity\npublic class %s {\n", Name))
	// write content
	sb.WriteString(getPrintableString(content))
	// close the class
	sb.WriteString("\n}")

	return sb.String()
}
