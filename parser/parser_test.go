package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/alanfoster/monkey/lexer"
	"github.com/alanfoster/monkey/ast"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	assert.Equal(t, 3, len(program.Statements))

	expectedStatements := []struct {
		identifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, expected := range expectedStatements {
		stmt := program.Statements[i]

		letStmt, ok := stmt.(*ast.LetStatement)
		assert.True(t, ok)

		assert.Equal(t, "let", stmt.TokenLiteral())
		assert.Equal(t, expected.identifier, letStmt.Name.Value)
		assert.Equal(t, expected.identifier, letStmt.Name.TokenLiteral())

	}
}
