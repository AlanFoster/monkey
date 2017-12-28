package parser

import (
	"github.com/alanfoster/monkey/lexer"
	"github.com/alanfoster/monkey/token"
	"github.com/alanfoster/monkey/ast"
	"fmt"
)

type Precedence int

// Precedences for Monkey. The order here is important, more so than specific values.
const (
	// https://github.com/golang/go/wiki/Iota
	_               Precedence = iota
	LOWEST
	EQUALS           // ==
	LESS_OR_GREATER  // > or <
	SUM              // +
	PRODUCT          // *
	PREFIX           // -X or !X
	CALL             // myFunction(x)
)

type (
	prefixParseFn func() ast.Expression
	// The argument is the left hand side of the expression, i.e. `1 + 2` this will be `1`
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Read two tokens so curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	// Register the tokens for our expression parsing
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(tokenType token.TokenType, prefixParseFn prefixParseFn) {
	p.prefixParseFns[tokenType] = prefixParseFn
}

func (p *Parser) registerInfix(tokenType token.TokenType, infixParseFn infixParseFn) {
	p.infixParseFns[tokenType] = infixParseFn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) appendPeekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, but got %s instead", t, p.peekToken)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{

		},
	}

	for !p.isCurToken(token.EOF) {
		stmt := p.parseStatement();
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
	return nil
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// Allow optional semicolons for ease to the developer when using a REPL
	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.EQ) {
		return nil
	}

	// TODO: Skip over parsing expressions for now until we see a semicolon
	for !p.isCurToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt;
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: Skip over parsing expressions for now until we see a semicolon
	for !p.isCurToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt;
}

func (p *Parser) isCurToken(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) isPeekToken(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.isPeekToken(t) {
		p.nextToken()
		return true
	} else {
		p.appendPeekError(t)
		return false
	}
}
