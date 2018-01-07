package lexer

import (
	"testing"
	"github.com/alanfoster/monkey/token"
	"github.com/stretchr/testify/assert"
)

func TestBasicPunctuation(t *testing.T) {
	input := `
		=
		==
		+
		-
		!
		!=
		*
		/
		<
		>
		,
		;
		()
		{}
		[]
	`

	expectedTokens := []struct {
		Type    token.TokenType
		Literal string
	}{
		{token.EQ, "="},
		{token.EQ_EQ, "=="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.BANG, "!"},
		{token.NOT_EQ, "!="},
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
		{token.LEFT_BRACKET, "["},
		{token.RIGHT_BRACKET, "]"},
		{token.EOF, ""},
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

		let foo = "hello world";
		let emptyString = "";
	`

	expectedTokens := []struct {
		Type    token.TokenType
		Literal string
	}{
		// let five = 5;
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.EQ, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// let ten = 10;
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.EQ, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// let add = fn(x, y) { x + y; };
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.EQ, "="},
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
		{token.EQ, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFT_PAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_PAREN, ")"},

		// let foo = "hello world";
		{token.LET, "let"},
		{token.IDENTIFIER, "foo"},
		{token.EQ, "="},
		{token.STRING, "hello world"},
		{token.SEMICOLON, ";"},

		// let emptyString = "";
		{token.LET, "let"},
		{token.IDENTIFIER, "emptyString"},
		{token.EQ, "="},
		{token.STRING, ""},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	l := New(input)

	for _, expected := range expectedTokens {
		tok := l.NextToken()

		assert.Equal(t, expected.Type, tok.Type)
		assert.Equal(t, expected.Literal, tok.Literal)
	}
}

func TestBranching(t *testing.T) {
	input := `
		if ( 5 < 10 ) {
			return true;
		} else {
			return false;
		}
	`

	expectedTokens := []struct {
		Type    token.TokenType
		Literal string
	}{
		// if ( 5 < 10 ) {
		{token.IF, "if"},
		{token.LEFT_PAREN, "("},
		{token.INT, "5"},
		{token.LESS_THAN, "<"},
		{token.INT, "10"},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},

		// return true;
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},

		// } else {
		{token.RIGHT_BRACE, "}"},
		{token.ELSE, "else"},
		{token.LEFT_BRACE, "{"},

		// return false;
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},

		// }
		{token.RIGHT_BRACE, "}"},

		{token.EOF, ""},
	}

	l := New(input)

	for _, expected := range expectedTokens {
		tok := l.NextToken()

		assert.Equal(t, expected.Type, tok.Type)
		assert.Equal(t, expected.Literal, tok.Literal)
	}
}
