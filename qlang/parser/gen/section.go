package gen

type SubSection struct {
	Name    string
	Content []string
}

type Section struct {
	Decorators  []string
	SubSections []*SubSection
}
