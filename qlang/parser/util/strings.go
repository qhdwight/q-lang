package util

import (
	"bufio"
	"bytes"
	"strings"
)

func Tokenize(str string, kws []string) []string {
	scanner := bufio.NewScanner(strings.NewReader(str))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Split keywords even if not separated by white space
		for _, kw := range kws {
			kwData := []byte(kw) // TODO inefficient conversion each time
			if bytes.Equal(kwData, data[:len(kw)]) {
				return len(kw), kwData, nil
			}
		}
		return bufio.ScanWords(data, atEOF)
	}
	scanner.Split(split)
	var tokens []string
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}
	return tokens
}
