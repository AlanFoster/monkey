package evaluator

import (
	"github.com/alanfoster/monkey/ast"
	"github.com/alanfoster/monkey/object"
	"fmt"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(o object.Object) bool {
	return o != nil && o.Type() == object.ERROR
}

func evalProgram(statements []ast.Statement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, environment)

		switch result := result.(type) {
		case *object.ReturnValue:
			// Unwrap the final intended result
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(statements []ast.Statement, environment *object.Environment) object.Object {
	var result object.Object = NULL

	for _, statement := range statements {
		result = Eval(statement, environment)

		// Return the wrapped return value, as we may be nested in multiple block statements
		if rt := result.Type(); rt == object.RETURN_VALUE || rt == object.ERROR {
			return result
		}
	}

	return result
}

func asBoolean(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func evalBangPrefixOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangPrefixOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Integer, right object.Integer) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return asBoolean(leftVal > rightVal)
	case "<":
		return asBoolean(leftVal < rightVal)
	case "==":
		return asBoolean(leftVal == rightVal)
	case "!=":
		return asBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		first := *left.(*object.Integer)
		second := *right.(*object.Integer)
		return evalIntegerInfixExpression(operator, first, second)
	case operator == "==":
		return asBoolean(left == right)
	case operator == "!=":
		return asBoolean(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func isTruthy(o object.Object) bool {
	switch o {
	case FALSE:
		return false
	case NULL:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func evalIfExpression(node *ast.IfExpression, environment *object.Environment) object.Object {
	result := Eval(node.Predicate, environment)

	if isError(result) {
		return result
	}

	if isTruthy(result) {
		return Eval(node.TrueBlock, environment)
	} else if node.FalseBlock != nil {
		return Eval(node.FalseBlock, environment)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, environment *object.Environment) object.Object {
	if value, ok := environment.Get(node.Value); ok {
		return value
	}

	return newError("identifier not found: %s", node.Value)
}

func evalExpressions(expressions []ast.Expression, environment *object.Environment) []object.Object {
	var args []object.Object
	for _, argument := range expressions {
		argumentValue := Eval(argument, environment)
		if isError(argumentValue) {
			return []object.Object{argumentValue}
		}
		args = append(args, argumentValue)
	}
	return args
}

func extendFunctionEnvironment(function *object.Function, args []object.Object) *object.Environment {
	newEnvironment := object.NewClosedEnvironment(function.Environment)
	for index, identifier := range function.Parameters {
		newEnvironment.Add(identifier.Value, args[index])
	}
	return newEnvironment
}

func unwrapResult(o object.Object) object.Object {
	if returnValue, ok := o.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return o
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	scopedEnvironment := extendFunctionEnvironment(function, args)
	evaluated := Eval(function.Body, scopedEnvironment)
	return unwrapResult(evaluated)
}

func evalCallExpression(node *ast.CallExpression, environment *object.Environment) object.Object {
	function := Eval(node.Function, environment)
	if isError(function) {
		return function
	}

	argumentValues := evalExpressions(node.Arguments, environment)
	if len(argumentValues) == 1 && isError(argumentValues[0]) {
		return argumentValues[0]
	}

	return applyFunction(function, argumentValues)
}

func Eval(node ast.Node, environment *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, environment)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, environment)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return asBoolean(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, environment)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, environment)
	case *ast.ReturnStatement:
		value := Eval(node.Value, environment)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.LetStatement:
		value := Eval(node.Value, environment)
		if isError(value) {
			return value
		}

		environment.Add(node.Name.Value, value)
		return value
	case *ast.Identifier:
		return evalIdentifier(node, environment)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Environment: environment}
	case *ast.CallExpression:
		return evalCallExpression(node, environment)
	}

	panic(fmt.Sprintf("Unexpected value %#v", node))
}
