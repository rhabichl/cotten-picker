package models

import (
	"fmt"
	"strings"
)

type Printable interface {
	Print() string
}

type Import struct {
	Name string
}

type Args struct {
	Name string
	Type string
}

type Function struct {
	Security     string
	ReturnType   string
	Name         string
	FunctionArgs []Args
	Annotations  []string
	Content      string
}

// implement Printable interface
func (f Function) Print() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s%v %v %v", f.getAnno(), f.Security, f.ReturnType, f.Name))
	sb.WriteString(fmt.Sprintf("(%s) {\n", f.getArgs()))
	sb.WriteString(f.Content)
	sb.WriteString("\n}")
	return sb.String()
}

func (i Import) Print() string {
	return fmt.Sprintf("import %s;", i.Name)
}

// helper functions
func getPrintableString(content []Printable) string {
	var sb strings.Builder

	for _, c := range content {
		sb.WriteString(fmt.Sprintf("\n%s\n", c.Print()))
	}
	return sb.String()
}

func (f Function) getArgs() string {
	if len(f.FunctionArgs) == 0 {
		return ""
	}
	var sb strings.Builder
	// get all but the last item
	for i := 0; i < len(f.FunctionArgs)-1; i++ {
		sb.WriteString(fmt.Sprintf("%s %s, ", f.FunctionArgs[i].Type, f.FunctionArgs[i].Name))
	}

	sb.WriteString(fmt.Sprintf("%s %s", f.FunctionArgs[len(f.FunctionArgs)-1].Type, f.FunctionArgs[len(f.FunctionArgs)-1].Name))

	return sb.String()
}

func (f Function) getAnno() string {
	if len(f.Annotations) == 0 {
		return ""
	}
	var sb strings.Builder
	for i := 0; i < len(f.Annotations); i++ {
		sb.WriteString(f.Annotations[i] + "\n")
	}
	return sb.String()
}
