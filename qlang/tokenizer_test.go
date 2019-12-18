package main

import (
	"q-lang-go/parser/node"
	"q-lang-go/parser/util"
	"testing"
)

func TestTokenizer(t *testing.T) {
	scanner := util.NewScanner("  a program  {qlang()&&||b } ")
	if scanner.Next(node.Split) != "a" {
		t.Error()
	}
	if scanner.Next(node.Split) != "program" {
		t.Error()
	}
	if scanner.Next(node.Split) != "{" {
		t.Error()
	}
	if scanner.Next(node.Split) != "qlang" {
		t.Error()
	}
	if scanner.Next(node.Split) != "(" {
		t.Error()
	}
	if scanner.Next(node.Split) != ")" {
		t.Error()
	}
	if scanner.Next(node.Split) != "&&" {
		t.Error()
	}
	if scanner.Next(node.Split) != "||" {
		t.Error()
	}
	if scanner.Next(node.Split) != "b" {
		t.Error()
	}
	if scanner.Next(node.Split) != "}" {
		t.Error()
	}
}
