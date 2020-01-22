package gen

type SubSection struct {
	Label   string
	Content []string
}

type Section struct {
	Decorators  []string
	SubSections []*SubSection
}
