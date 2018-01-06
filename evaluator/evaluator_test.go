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

func assertNullObject(t *testing.T, o object.Object) {
	assert.EqualValues(t, NULL, o)
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

func TestEvalIntegerOperatorExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"1 + 1 + 1;",
			3,
		},
		{
			"1 - 1",
			0,
		},
		{
			"5 * 5;",
			25,
		},
		{
			"-5 * -5",
			25,
		},
		{
			"5 / 5",
			1,
		},
		{
			"25 / 5",
			5,
		},
		{
			"5 / 2",
			2,
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

func TestInfixIntegerBooleanOperatorExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"1 < 2;",
			true,
		},
		{
			"1 > 2;",
			false,
		},
		{
			"1 < 1;",
			false,
		},
		{
			"1 > 1;",
			false,
		},
		{
			"1 == 1;",
			true,
		},
		{
			"1 == 2;",
			false,
		},
		{
			"1 != 1;",
			false,
		},
		{
			"1 != 2;",
			true,
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

func TestInfixBooleanBooleanOperatorExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"true == true;",
			true,
		},
		{
			"true == false;",
			false,
		},
		{
			"true != true;",
			false,
		},
		{
			"true != false;",
			true,
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		assertBooleanObject(t, evaluated, test.expected)
	}
}

func TestIfStatementExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
				if (true) {
					10;
				}
			`,
			10,
		},
		{
			`
				if (false) {
					10;
				}
			`,
			nil,
		},
		{
			`
				if (true) {
					10;
				} else {
					20;
				}
			`,
			10,
		},
		{
			`
				if (false) {
					10;
				} else {
					20;
				}
			`,
			20,
		},
		{
			`
				if (1) {
					10;
					30;
				} else {
					20;
				}
			`,
			30,
		},
		{
			`
				if (0) {
					10;
					30;
				} else {
					20;
				}
			`,
			30,
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		integer, ok := test.expected.(int)

		if ok {
			assertIntegerObject(t, evaluated, int64(integer))
		} else {
			assertNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"return 5;",
			5,
		},
		{
			"return 5; 10;",
			5,
		},
		{
			"return 5; return 10;",
			5,
		},
		{
			"10; return 2; return 7",
			2,
		},
		{
			"if (10 > 5) { return 1; } return 0;",
			1,
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		assertIntegerObject(t, evaluated, test.expected)
	}
}