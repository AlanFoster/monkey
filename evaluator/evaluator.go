package evaluator

import (
	"github.com/alanfoster/monkey/ast"
	"github.com/alanfoster/monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL = &object.Null{}
)

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

func asBoolean(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func evalBangOperatorExpression(object object.Object) object.Object {
	switch object {
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

func evalPrefixExpression(operator string, object object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(object)
	default:
		return NULL
	}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return asBoolean(node.Value)
	case *ast.PrefixExpression:
		value := Eval(node.Right)
		return evalPrefixExpression(node.Operator, value)
	}

	return nil
}
