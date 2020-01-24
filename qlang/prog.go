package main

import (
	"strings"
)

type SubSect struct {
	Label   string
	Content []string
	Vars    map[string]int
	Anons   map[Node]int
}

type Sect struct {
	Decorators []string
	SubSects   []*SubSect
}

type Program struct {
	ConstSect, FuncSect *Sect
	FuncStackHead       int
	FuncSubSect         *SubSect
}

func (program *Program) ToString() string {
	builder := strings.Builder{}
	for _, sect := range []*Sect{program.ConstSect, program.FuncSect} {
		for _, decorator := range sect.Decorators {
			builder.WriteString(".")
			builder.WriteString(decorator)
			builder.WriteString("\n")
		}
		for _, subSects := range sect.SubSects {
			builder.WriteString("_")
			builder.WriteString(subSects.Label)
			builder.WriteString(":")
			builder.WriteString("\n")
			for _, line := range subSects.Content {
				builder.WriteString("    ")
				builder.WriteString(line)
				builder.WriteString("\n")
			}
			builder.WriteString("\n")
		}
	}
	return builder.String()
}