package util

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

//func (scanner *Scanner) NextWord() string {
//	return scanner.Next(SplitWord)
//}
//
//func SplitWord(rest string) (string, int) {
//	blankLength := 0
//	for ; blankLength < len(rest); blankLength++ {
//		if !unicode.IsSpace(rune(rest[blankLength])) {
//			break
//		}
//	}
//	word := rest[blankLength:]
//	wordLength := 0
//	for ; wordLength < len(word); wordLength++ {
//		if unicode.IsSpace(rune(word[wordLength])) {
//			break
//		}
//	}
//	return word[:wordLength], blankLength + wordLength
//}
