package gen

import "strings"

type Program struct {
	ConstantsSection, FuncSection *Section
	CurrentSubSection             *SubSection
}

func (program *Program) ToString() string {
	builder := strings.Builder{}
	for _, section := range []*Section{program.ConstantsSection, program.FuncSection} {
		for _, decorator := range section.Decorators {
			builder.WriteString(decorator)
			builder.WriteString("\n")
		}
		for _, subSections := range section.SubSections {
			builder.WriteString(subSections.Name)
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
