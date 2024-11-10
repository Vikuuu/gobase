package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	input := `package main

	type users struct {
		ID int
		Name string
		IsActive bool
	}
	`

	tests := []Token{
		{Type: PACKAGE, Literal: "package"},
		{Type: IDENT, Literal: "main"},
		{Type: TYPE, Literal: "type"},
		{Type: IDENT, Literal: "users"},
		{Type: STRUCT, Literal: "struct"},
		{Type: LBRACE, Literal: "{"},
		{Type: IDENT, Literal: "ID"},
		{Type: INT, Literal: "int"},
		{Type: IDENT, Literal: "Name"},
		{Type: STRING, Literal: "string"},
		{Type: IDENT, Literal: "IsActive"},
		{Type: BOOL, Literal: "bool"},
		{Type: RBRACE, Literal: "}"},
		{Type: EOF, Literal: ""},
	}

	scanner := NewScanner(input)
	parser := NewLexer(scanner)

	for i, expected := range tests {
		tok := parser.nextToken()
		if tok.Type != expected.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, expected.Type, tok.Type)
		}
		if tok.Literal != expected.Literal {
			t.Fatalf(
				"tests[%d] - tokentype wrong. expected=%q, got=%q",
				i,
				expected.Literal,
				tok.Literal,
			)
		}
	}
}
