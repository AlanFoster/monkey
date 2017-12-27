package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/alanfoster/monkey/lexer"
	"github.com/bradleyjkemp/cupaloy"
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
	assert.Empty(t, p.Errors())

	cupaloy.SnapshotT(t, program)
}

func TestInvalidParsing(t *testing.T) {
	input := `
		let x 10;
		let x
	`
	l := lexer.New(input)
	p := New(l)

	p.ParseProgram()
	expectedErrors := []string{
		"expected next token to be =, but got {INT 10} instead",
		"expected next token to be =, but got {EOF } instead",
	}
	assert.Equal(t, expectedErrors, p.Errors())
}
