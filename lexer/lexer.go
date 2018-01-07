package lexer

import (
	"github.com/alanfoster/monkey/token"
)

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
		if l.peekChar() == '=' {
			tok = newStringToken(token.EQ_EQ, l.readTwoCharacterLiteral())
		} else {
			tok = newCharToken(token.EQ, l.ch)
		}
	case '+':
		tok = newCharToken(token.PLUS, l.ch)
	case '-':
		tok = newCharToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			tok = newStringToken(token.NOT_EQ, l.readTwoCharacterLiteral())
		} else {
			tok = newCharToken(token.BANG, l.ch)
		}
	case '*':
		tok = newCharToken(token.ASTERISK, l.ch)
	case '/':
		tok = newCharToken(token.SLASH, l.ch)
	case '<':
		tok = newCharToken(token.LESS_THAN, l.ch)
	case '>':
		tok = newCharToken(token.GREATER_THAN, l.ch)
	case ',':
		tok = newCharToken(token.COMMA, l.ch)
	case ';':
		tok = newCharToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newCharToken(token.LEFT_PAREN, l.ch)
	case ')':
		tok = newCharToken(token.RIGHT_PAREN, l.ch)
	case '{':
		tok = newCharToken(token.LEFT_BRACE, l.ch)
	case '}':
		tok = newCharToken(token.RIGHT_BRACE, l.ch)
	case '"':
		tok = newStringToken(token.STRING, l.readString())
	case 0:
		tok = newStringToken(token.EOF, "")
	default:
		if isIdentifierLetter(l.ch) {
			literal := l.readIdentifier()
			tokenType := token.LookupIdentifier(literal)
			return newStringToken(tokenType, literal)
		} else if isDigit(l.ch) {
			return newStringToken(token.INT, l.readNumber())
		} else {
			tok = newCharToken(token.ILLEGAL, l.ch)
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

// Reads a two character literal from the input string. Useful for two character operators
func (l *Lexer) readTwoCharacterLiteral() string {
	currentChar := l.ch
	l.readChar()
	nextChar := l.ch

	return string(currentChar) + string(nextChar)
}

func (l *Lexer) readString() string {
	// We don't want the character " within our lexeme
	l.readChar()
	position := l.position

	for l.ch != '"' && l.ch != 0 {
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

// Peek at the next character within the input stream without updating the position
// of the lexer internally
func (l *Lexer) peekChar() byte {
	return l.input[l.position+1]
}

func newCharToken(tokenType token.TokenType, ch byte) token.Token {
	return newStringToken(tokenType, string(ch))
}

func newStringToken(tokenType token.TokenType, string string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string,
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
