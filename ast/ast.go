package ast

import (
	"github.com/alanfoster/monkey/token"
	"bytes"
)

type Node interface {
	// Temporary function only used for debugging and testing
	TokenLiteral() string
	PrettyPrint() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) PrettyPrint() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.PrettyPrint())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) PrettyPrint() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.PrettyPrint())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.PrettyPrint())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) PrettyPrint() string {
	return i.Value
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) PrettyPrint() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.PrettyPrint())
	}

	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) PrettyPrint() string {
	if es.Expression != nil {
		return es.Expression.PrettyPrint()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) PrettyPrint() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) PrettyPrint() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.PrettyPrint())
	out.WriteString(")")

	return out.String()
}

// AST For `left [operand] right`, i.e. `5 + 5`
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) PrettyPrint() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.PrettyPrint())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.PrettyPrint())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) PrettyPrint() string {
	return b.Token.Literal
}