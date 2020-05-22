package main

type Scanner struct {
	str   string
	index int
}

func NewScanner(str string) *Scanner {
	return &Scanner{str: str}
}

func (scanner *Scanner) peek(tokenFunc func(str string) (string, int)) (string, int) {
	if scanner.index >= len(scanner.str) {
		return "", 0
	}
	token, skip := tokenFunc(scanner.str[scanner.index:])
	return token, skip
}

func (scanner *Scanner) Skip(tokenFunc func(str string) (string, int)) {
	scanner.Next(tokenFunc)
}

func (scanner *Scanner) Next(tokenFunc func(str string) (string, int)) string {
	token, skip := scanner.peek(tokenFunc)
	scanner.index += skip
	return token
}

func (scanner *Scanner) Peek(tokenFunc func(str string) (string, int)) (string, int) {
	token, skip := scanner.peek(tokenFunc)
	return token, skip
}

func (scanner *Scanner) PeekAdvanceIf(tokenFunc func(str string) (string, int), skipCond func(str string) bool) bool {
	peek, skip := scanner.Peek(tokenFunc)
	if skipCond(peek) {
		scanner.index += skip
		return true
	}
	return false
}
