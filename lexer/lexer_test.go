package lexer

import (
	"interpreter/token"
	"testing"
)

func TestNestToken(t *testing.T) {
	input := "=+(){},;"

	tests := []struct {
		ExpectedType    token.TokenType
		ExpectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.ExpectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q got=%q", i, tt.ExpectedType, tok.Type)
		}
		if tok.Literal != tt.ExpectedLiteral {
			t.Fatalf("tests[%d] - tokenLiteral wrong. expected=%q got=%q", i, tt.ExpectedLiteral, tok.Literal)
		}
	}
}
