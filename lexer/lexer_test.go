package lexer

import (
	"testing"
	"github.com/alanfoster/monkey/token"
	"github.com/stretchr/testify/assert"
)

func TestBasicPunctuation(t *testing.T) {
	input := `=+-!*/<>,;(){}`

	expectedTokens := []struct {
		Type    token.TokenType
		Literal string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.BANG, "!"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.LESS_THAN, "<"},
		{token.GREATER_THAN, ">"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.LEFT_PAREN, "("},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
	}

	l := New(input)

	for _, expected := range expectedTokens {
		tok := l.NextToken()

		assert.Equal(t, expected.Type, tok.Type)
		assert.Equal(t, expected.Literal, tok.Literal)
	}
}

func TestAssignmentAndFunctions(t *testing.T) {
	input := `
		let five = 5;
		let ten = 10;

		let add = fn(x, y) {
			x + y;
		};

		let result = add(five, ten)
	`

	expectedTokens := []struct {
		Type    token.TokenType
		Literal string
	}{
		// let five = 5;
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// let ten = 10;
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// let add = fn(x, y) { x + y; };
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LEFT_PAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},

		// let result = add(five, ten)
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFT_PAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_PAREN, ")"},
	}

	l := New(input)

	for _, expected := range expectedTokens {
		tok := l.NextToken()

		assert.Equal(t, expected.Type, tok.Type)
		assert.Equal(t, expected.Literal, tok.Literal)
	}
}
