package main

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	scanner := NewScanner("  a program  {qlang()&&||b } ")
	if scanner.Next(Split) != "a" {
		t.Error()
	}
	if scanner.Next(Split) != "program" {
		t.Error()
	}
	if scanner.Next(Split) != "{" {
		t.Error()
	}
	if scanner.Next(Split) != "qlang" {
		t.Error()
	}
	if scanner.Next(Split) != "(" {
		t.Error()
	}
	if scanner.Next(Split) != ")" {
		t.Error()
	}
	if scanner.Next(Split) != "&&" {
		t.Error()
	}
	if scanner.Next(Split) != "||" {
		t.Error()
	}
	if scanner.Next(Split) != "b" {
		t.Error()
	}
	if scanner.Next(Split) != "}" {
		t.Error()
	}
}
