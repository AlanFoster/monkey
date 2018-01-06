package parser

import (
	"github.com/alanfoster/monkey/lexer"
	"github.com/alanfoster/monkey/token"
	"github.com/alanfoster/monkey/ast"
	"fmt"
	"strconv"
)

type Precedence int

// Available precedences for Monkey. The order here is important, more so than specific values.
const (
	_               Precedence = iota
	LOWEST
	EQUALS           // == or !=
	LESS_OR_GREATER  // > or <
	SUM              // + or -
	PRODUCT          // * or /
	PREFIX           // -X or !X
	CALL             // myFunction(x)
)

// This particular parser does not make use of a separate left/right precedence, instead they are
// the same value
var precedences = map[token.TokenType]Precedence{
	token.EQ_EQ:        EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LESS_THAN:    LESS_OR_GREATER,
	token.GREATER_THAN: LESS_OR_GREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.LEFT_PAREN:   CALL,
}

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
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)
	p.registerPrefix(token.LEFT_PAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfStatement)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ_EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)
	p.registerInfix(token.LEFT_PAREN, p.parseCallExpression)

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

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseIfStatement() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LEFT_PAREN) {
		return nil
	}

	p.nextToken()
	expression.Predicate = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}

	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}
	expression.TrueBlock = p.parseBlockStatement()

	if p.isPeekToken(token.ELSE) {
		p.expectPeek(token.ELSE)

		if !p.expectPeek(token.LEFT_BRACE) {
			return nil
		}

		expression.FalseBlock = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	// This function assumes the curToken is currently on the left brace
	p.nextToken()

	for !p.isCurToken(token.RIGHT_BRACE) && !p.isCurToken(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	integerLiteral := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}

	integerLiteral.Value = value

	return integerLiteral
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteral := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LEFT_PAREN) {
		return nil
	}

	functionLiteral.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}

	functionLiteral.Body = p.parseBlockStatement()

	return functionLiteral
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	p.expectCur(token.LEFT_PAREN)

	for !p.isCurToken(token.RIGHT_PAREN) {
		identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, identifier)
		p.expectCur(token.IDENTIFIER)

		if !p.isCurToken(token.COMMA) {
			break
		}

		if !p.expectCur(token.COMMA) {
			return nil
		}
	}

	// The convention is the calling of next function consumes the next token
	// p.expectCur(token.RIGHT_PAREN)

	return identifiers
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	expression := &ast.CallExpression{
		Token:    p.curToken,
		Function: left,
	}
	expression.Arguments = p.parseFunctionArguments()

	return expression
}

func (p *Parser) parseFunctionArguments() []ast.Expression {
	args := []ast.Expression{}
	p.expectCur(token.LEFT_PAREN)

	for !p.isCurToken(token.RIGHT_PAREN) {
		arg := p.parseExpression(LOWEST)
		args = append(args, arg)
		p.nextToken()

		if !p.isCurToken(token.COMMA) {
			break;
		}

		if !p.expectCur(token.COMMA) {
			return nil
		}
	}

	// Unlike the previous parseFunction implementation, we consume this right paren
	p.expectCur(token.RIGHT_PAREN)

	return args
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.isCurToken(token.TRUE)}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) appendCurError(t token.TokenType) {
	msg := fmt.Sprintf("expected current token to be %s, but got %s instead", t, p.peekToken)
	p.errors = append(p.errors, msg)
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
		stmt := p.parseStatement()
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

func (p *Parser) appendPrefixParseError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.appendPrefixParseError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.isPeekToken(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		// Note the explicit re-assignment to 'leftExp'
		leftExp = infix(leftExp)
	}

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

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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

func (p *Parser) expectCur(t token.TokenType) bool {
	if p.isCurToken(t) {
		p.nextToken()
		return true
	} else {
		p.appendCurError(t)
		return false
	}
}

func (p *Parser) peekPrecedence() Precedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() Precedence {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
