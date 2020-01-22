package gen

import "strings"

type Program struct {
	ConstSections, FuncSection *Section
	FuncStackHead              int
	FuncSubSection             *SubSection
}

func (program *Program) ToString() string {
	builder := strings.Builder{}
	for _, section := range []*Section{program.ConstSections, program.FuncSection} {
		for _, decorator := range section.Decorators {
			builder.WriteString(".")
			builder.WriteString(decorator)
			builder.WriteString("\n")
		}
		for _, subSections := range section.SubSections {
			builder.WriteString("_")
			builder.WriteString(subSections.Label)
			builder.WriteString(":")
			builder.WriteString("\n")
			for _, line := range subSections.Content {
				builder.WriteString("    ")
				builder.WriteString(line)
				builder.WriteString("\n")
			}
			builder.WriteString("\n")
		}
	}
	return builder.String()
}
