package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/alanfoster/monkey/lexer"
	"github.com/bradleyjkemp/cupaloy"
)

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

func TestReturnStatements(t *testing.T) {
	input := `return 5;`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	assert.Empty(t, p.Errors())

	cupaloy.SnapshotT(t, program)
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())

	cupaloy.SnapshotT(t, program)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `1337;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())

	cupaloy.SnapshotT(t, program)
}

func TestParsingPrefixExpression(t *testing.T) {
	input := `!5; -15;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())

	cupaloy.SnapshotT(t, program)
}

func TestParsingInfixExpression(t *testing.T) {
	input := `
		5 + 5;
		5 - 5;
		5 * 5;
		5 / 5;
		5 > 5;
		5 < 5;
		5 == 5;
		5 != 5;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())

	cupaloy.SnapshotT(t, program)
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input              string
		expectedPrettyText string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		assert.Empty(t, p.Errors())
		assert.Equal(t, test.expectedPrettyText, program.PrettyPrint())
	}
}

func TestTrueBoolean(t *testing.T) {
	tests := []struct {
		input              string
		expectedPrettyText string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"false == true",
			"(false == true)",
		},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		assert.Empty(t, p.Errors())
		assert.Equal(t, test.expectedPrettyText, program.PrettyPrint())
	}
}