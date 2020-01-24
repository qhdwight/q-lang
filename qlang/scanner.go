package main

type Scanner struct {
	str   string
	index int
}

func NewScanner(str string) *Scanner {
	return &Scanner{str: str}
}

func (scanner *Scanner) Next(tokenFunc func(str string) (string, int)) string {
	if scanner.index >= len(scanner.str) {
		return ""
	}
	token, skip := tokenFunc(scanner.str[scanner.index:])
	scanner.index += skip
	return token
}
