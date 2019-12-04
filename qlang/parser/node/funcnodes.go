package node

import (
	"fmt"
	"q-lang-go/parser/util"
	"strings"
)

type DefineFuncNode struct {
	ParseNode
}

type ImplFuncNode struct {
	ParseNode
}

func (node *ImplFuncNode) Parse(body, innerBody string, tokens []string, parent Node) {
	lines := strings.Split(innerBody, ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		tokens := util.Tokenize(line, []string{
			"+", "-", "*", "/", "&&", "||",
		})
		//firstToken := tokens[0]
		fmt.Println("[" + strings.Join(tokens, "|") + "]")
	}
}
