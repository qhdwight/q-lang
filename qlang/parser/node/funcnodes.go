package node

type DefineFuncNode struct {
	ParseNode
}

type ImplFuncNode struct {
	ParseNode
}

//func (node *ImplFuncNode) Parse(body, innerBody string, tokens []string, parent Node) {
//	lines := strings.Split(innerBody, ";")
//	for _, line := range lines {
//		tokens := tokenRegex.FindAllString(line, -1)
//		firstToken := tokens[0]
//		fmt.Println(firstToken)
//		if varNodeFunc, isVar := varNodeFactory[firstToken]; isVar {
//			varNode := varNodeFunc()
//			node.Add(varNode)
//		}
//	}
//}
