package evaluator

import (
	"testing"
	"github.com/alanfoster/monkey/object"
	"github.com/alanfoster/monkey/lexer"
	"github.com/alanfoster/monkey/parser"
	"github.com/stretchr/testify/assert"
)

func assertIntegerObject(t *testing.T, o object.Object, expected int64) {
	assert.IsType(t, new(object.Integer), o)
	result, ok := o.(*object.Integer)
	if assert.True(t, ok) {
		assert.Equal(t, expected, result.Value)
	}
}

func assertBooleanObject(t *testing.T, o object.Object, expected bool) {
	assert.IsType(t, new(object.Boolean), o)
	result, ok := o.(*object.Boolean)
	if assert.True(t, ok) {
		assert.Equal(t, expected, result.Value)
	}
}

func assertNullObject(t *testing.T, o object.Object, expected bool) {
	assert.IsType(t, new(object.Null), o)
}

func eval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())
	return Eval(program)
}

func TestEvalIntegerExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"15;",
			15,
		},
		{
			"1337",
			1337,
		},
		{
			"-15;",
			-15,
		},
		{
			"-1337",
			-1337,
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		assertIntegerObject(t, evaluated, test.expected)
	}
}

func TestBooleanExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"true;",
			true,
		},
		{
			"false;",
			false,
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		assertBooleanObject(t, evaluated, test.expected)
	}
}


func TestPrefixBangBooleanExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"!true;",
			false,
		},
		{
			"!false;",
			true,
		},
		{
			"!!false;",
			false,
		},
		{
			"!!!false;",
			true,
		},
		{
			"!5;",
			false,
		},
		{
			"!!5;",
			true,
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		assertBooleanObject(t, evaluated, test.expected)
	}
}