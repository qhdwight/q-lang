package gen

type SubSection struct {
	Label     string
	Content   []string
	Variables map[string]int
}

type Section struct {
	Decorators  []string
	SubSections []*SubSection
}
