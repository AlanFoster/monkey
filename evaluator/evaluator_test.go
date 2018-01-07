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

func assertErrorObject(t *testing.T, o object.Object, expected string) {
	assert.IsType(t, new(object.Error), o)
	result, ok := o.(*object.Error)
	if assert.True(t, ok) {
		assert.Equal(t, expected, result.Message)
	}
}

func assertStringObject(t *testing.T, o object.Object, expected string) {
	assert.IsType(t, new(object.String), o)
	result, ok := o.(*object.String)
	if assert.True(t, ok) {
		assert.Equal(t, expected, result.Value)
	}
}

func assertNullObject(t *testing.T, o object.Object) {
	assert.Equal(t, NULL, o)
}

func eval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	environment := object.NewEnvironment()
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())
	return Eval(program, environment)
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
		{
			`"hello" == "hello";`,
			true,
		},
		{
			`"hello" == " hello ";`,
			false,
		},
		{
			`"hello" != "hello";`,
			false,
		},
		{
			`"hello" != " hello ";`,
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

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true;",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { if (10 > 1) { true + false; } }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 + true) { 10 }",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"nonExistentVariable;",
			"identifier not found: nonExistentVariable",
		},
		{
			"5(1, 2, 3);",
			"not a function: INTEGER",
		},
		{
			"let add = 5; add(1, 2, 3);",
			"not a function: INTEGER",
		},
		{
			"let add = fn(x, y) { x + y }; add(5, true);",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"let add = fn(x, y) { x + y }; add(true, 5 + true);",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			`"hello" - " world"`,
			"unknown operator: STRING - STRING",
		},
		{
			"[1, 2, 3] + [4, 5]",
			"unknown operator: ARRAY + ARRAY",
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		assertErrorObject(t, evaluated, test.expected)
	}
}

func TestAssignmentHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let a = 5; a",
			5,
		},
		{
			"let a = 5; let b = 10; a + b;",
			15,
		},
		{
			"let a = 5; let b = 10; let c = 20; a + b + c + 5;",
			40,
		},
		{
			"let a = 5; let b = a; a + b",
			10,
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		assertIntegerObject(t, evaluated, test.expected)
	}
}

func TestFunctionHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let identity = fn(x) { x }; identity(5)",
			5,
		},
		{
			"let identity = fn(x) { return x }; identity(5)",
			5,
		},
		{
			"let add = fn(x, y) { x + y }; add(3, 5)",
			8,
		},
		{
			"fn(x, y) { x + y }(5, 9)",
			14,
		},
		{
			"let x = 5; let identity = fn(x) { x }; identity(8)",
			8,
		},
		{
			`
				let newAdder = fn(x) { fn(y) { x + y }; };
				let addTwo = newAdder(2);
				addTwo(6)
			`,
			8,
		},
		{
			`
				let earlyReturn = fn(x) { if (x > 5) { return x } }
				earlyReturn(11)
			`,
			11,
		},
		{
			`
				let earlyReturn = fn(x) { if (x > 5) { return x } }
				let multipleCalls = fn(x) { earlyReturn(x); earlyReturn(x); earlyReturn(55); }
				multipleCalls(11)
			`,
			55,
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		assertIntegerObject(t, evaluated, test.expected)
	}
}

func TestStringHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`"hello world"`,
			"hello world",
		},
		{
			`let a = "hello world";`,
			"hello world",
		},
		{
			`let a = "hello "; let b = "world"; a + b + "!" + "!";`,
			"hello world!!",
		},
		{
			`let echo = fn(x) { x }; echo("hello") + echo(" world");`,
			"hello world",
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		assertStringObject(t, evaluated, test.expected)
	}
}

func TestLenFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`len("")`,
			0,
		},
		{
			`len("hello world")`,
			11,
		},
		{
			`len(1)`,
			"argument to `len` not supported, got INTEGER",
		},
		{
			`len("one", "two")`,
			"wrong number of arguments. got=2, want=1",
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		switch expected := test.expected.(type) {
		case int:
			assertIntegerObject(t, evaluated, int64(expected))
		case string:
			assertErrorObject(t, evaluated, expected)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 4]"

	evaluated := eval(t, input)
	result, ok := evaluated.(*object.Array)
	assert.True(t, ok)
	if ok {
		assertIntegerObject(t, result.Elements[0], 1)
		assertIntegerObject(t, result.Elements[1], 4)
		assertIntegerObject(t, result.Elements[2], 7)
	}
}

func TestArrayExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`[1, 2, 3][0]`,
			1,
		},
		{
			`[1, 2, 3][1]`,
			2,
		},
		{
			`[1, 2, 3][2]`,
			3,
		},
		{
			`[1, 2, 3][4]`,
			nil,
		},
		{
			`[1, 2, 3][-1]`,
			nil,
		},
		{
			`let array = [1, 2, 3]; let index = 1; array[index]`,
			2,
		},
		{
			`let array = [[1, 2, [3]], 2, 3]; array[0][2][0]`,
			3,
		},
	}

	for _, test := range tests {
		evaluated := eval(t, test.input)
		switch expected := test.expected.(type) {
		case int:
			assertIntegerObject(t, evaluated, int64(expected))
		default:
			assertNullObject(t, evaluated)
		}
	}
}