package lexer

import "github.com/alanfoster/monkey/token"

type Lexer struct {
	input        string
	position     int  // Current position in input (points to current char)
	readPosition int  // Current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()

	return lexer;
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // Ascii code for the 'NUL' character
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = newToken(token.BANG, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LESS_THAN, l.ch)
	case '>':
		tok = newToken(token.GREATER_THAN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LEFT_PAREN, l.ch)
	case ')':
		tok = newToken(token.RIGHT_PAREN, l.ch)
	case '{':
		tok = newToken(token.LEFT_BRACE, l.ch)
	case '}':
		tok = newToken(token.RIGHT_BRACE, l.ch)
	default:
		if isIdentifierLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// Reads an identifier from the input string. This will update the position of the
// lexer internally as appropriate
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isIdentifierLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Reads a number from the input string. This will update the position of the
// lexer internally a appropriate
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Skips any whitespace
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// A valid identifier in monkey is similar to Java, other than lack of dollar sign support
func isIdentifierLetter(ch byte) bool {
	return isLetter(ch) || ch == '_'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' ||
		ch == '\t' ||
		ch == '\n' ||
		ch == '\r'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
