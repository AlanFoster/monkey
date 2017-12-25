package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + Literal
	IDENTIFIER = "IDENTIFIER" // add, foobar, x, y
	INT        = "INT"        // 12345...

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Deliminators
	COMMA     = ","
	SEMICOLON = ";"

	LEFT_PAREN  = "("
	RIGHT_PAREN = ")"
	LEFT_BRACE  = "{"
	RIGHT_BRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
