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

func evalProgram(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)

		// Unwrap the final result
		if result.Type() == object.RETURN_VALUE {
			return result.(*object.ReturnValue).Value
		}
	}

	return result
}

func evalBlockStatement(statements []ast.Statement) object.Object {
	var result object.Object = NULL

	for _, statement := range statements {
		result = Eval(statement)

		// Return the wrapped return value, as we may be nested in multiple block statements
		if result.Type() == object.RETURN_VALUE {
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
		return NULL
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
		return NULL
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
		return NULL
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		first := *left.(*object.Integer)
		second := *right.(*object.Integer)
		return evalIntegerInfixExpression(operator, first, second)
	case operator == "==":
		return asBoolean(left == right)
	case operator == "!=":
		return asBoolean(left != right)
	default:
		return NULL
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

func evalIfExpression(node *ast.IfExpression) object.Object {
	result := Eval(node.Predicate)

	if isTruthy(result) {
		return Eval(node.TrueBlock)
	} else if node.FalseBlock != nil {
		return Eval(node.FalseBlock)
	} else {
		return NULL
	}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return asBoolean(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements)
	case *ast.ReturnStatement:
		value := Eval(node.Value)
		return &object.ReturnValue{Value: value}
	}

	panic(fmt.Sprintf("Unexpected value %#v", node))
}
