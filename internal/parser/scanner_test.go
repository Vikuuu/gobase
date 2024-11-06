package parser

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `package main

	type users struct {
		ID int 
		Name string
		IsActive bool
	}
	`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{PACKAGE, "package"},
		{IDENT, "main"},
		{TYPE, "type"},
		{IDENT, "users"},
		{STRUCT, "struct"},
		{LBRACE, "{"},
		{IDENT, "ID"},
		{INT, "int"},
		{IDENT, "Name"},
		{STRING, "string"},
		{IDENT, "IsActive"},
		{BOOL, "bool"},
		{RBRACE, "}"},
		{EOF, ""},
	}

	scanner := NewScanner(input)

	for i, tt := range tests {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf(
				"tests[%d] - tokentype wrong. expected=%q, got=%q",
				i,
				tt.expectedType,
				tok.Type,
			)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] - literal wrong. expected=%q, got=%q",
				i,
				tt.expectedLiteral,
				tok.Literal,
			)
		}
	}
}
