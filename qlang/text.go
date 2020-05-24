package main

import (
	"strings"
)

type (
	Sect struct {
		Decorators []string
		Label      string
		Content    []string
		SubSects   []*Sect
	}
	Prog struct {
		ConstSect, FuncSect *Sect
		LibrarySubSect      *Sect
		CurSect             *Sect
		Scope               *Scope
	}
)

func NewProg() *Prog {
	prog := &Prog{
		ConstSect: &Sect{
			Decorators: []string{"data"},
		},
		FuncSect: &Sect{
			Decorators: []string{"text", "intel_syntax noprefix", "globl _main"},
		},
	}
	prog.LibrarySubSect = new(Sect)
	prog.FuncSect.SubSects = append(prog.FuncSect.SubSects, prog.LibrarySubSect)
	prog.Scope = NewScope(nil)
	return prog
}

func (sect *Sect) ToString(builder *strings.Builder, indent int) string {
	for _, decorator := range sect.Decorators {
		builder.WriteString(strings.Repeat(" ", indent))
		builder.WriteString(".")
		builder.WriteString(decorator)
		builder.WriteString("\n")
	}
	if len(sect.Label) > 0 {
		builder.WriteString(strings.Repeat(" ", indent))
		// if sect.Label != "main" {
		// }
		builder.WriteString("_")
		builder.WriteString(sect.Label)
		builder.WriteString(":")
		builder.WriteString("\n")
		indent += 4
	}
	for _, line := range sect.Content {
		builder.WriteString(strings.Repeat(" ", indent))
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	for _, sect := range sect.SubSects {
		sect.ToString(builder, indent)
	}
	return builder.String()
}

func (prog *Prog) ToString() string {
	builder := &strings.Builder{}
	for _, sect := range []*Sect{prog.ConstSect, prog.FuncSect} {
		sect.ToString(builder, 0)
	}
	return builder.String()
}
