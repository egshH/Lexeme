// Package evaluator implements the tree-walking interpreter from chapter 3,
// through bindings and environments (section 3.9).
package evaluator

import (
	"fmt"

	"github.com/egshH/Lexeme/ast"
	"github.com/egshH/Lexeme/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval evaluates a successfully parsed AST in env. Declarations and empty
// programs return nil (no output); language-level absence of a value is NULL.
// Unsupported syntax returns an Error rather than silently producing null.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return booleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isControl(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		condition := Eval(node.Condition, env)
		if isControl(condition) {
			return condition
		}
		if isTruthy(condition) {
			return Eval(node.Consequence, env)
		}
		if node.Alternative != nil {
			return Eval(node.Alternative, env)
		}
		return NULL
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)
		if isControl(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if isControl(value) {
			return value
		}
		env.Set(node.Name.Value, value)
		return nil
	case *ast.Identifier:
		if value, ok := env.Get(node.Value); ok {
			return value
		}
		return newError("identifier not found: %s", node.Value)
	case nil: // A bare return; has no expression.
		return NULL
	default:
		return newError("evaluation not implemented for %T", node)
	}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if isControl(result) {
			return result
		}
	}
	if result == nil {
		return NULL
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return booleanObject(!isTruthy(right))
	case "-":
		if integer, ok := right.(*object.Integer); ok {
			return &object.Integer{Value: -integer.Value}
		}
	}
	return newError("unknown operator: %s%s", operator, right.Type())
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isControl(left) {
		return left
	}
	// Lexeme also has logical operators. Return booleans and skip the unused
	// operand, just as if/else skips the unused branch.
	if node.Operator == "&&" && !isTruthy(left) {
		return FALSE
	}
	if node.Operator == "||" && isTruthy(left) {
		return TRUE
	}
	right := Eval(node.Right, env)
	if isControl(right) {
		return right
	}
	if node.Operator == "&&" || node.Operator == "||" {
		return booleanObject(isTruthy(right))
	}
	if l, ok := left.(*object.Integer); ok {
		if r, ok := right.(*object.Integer); ok {
			return evalIntegerInfixExpression(node.Operator, l.Value, r.Value)
		}
	}
	// Like Monkey, equality between different types is false, not an error.
	if node.Operator == "==" {
		return booleanObject(equalObjects(left, right))
	}
	if node.Operator == "!=" {
		return booleanObject(!equalObjects(left, right))
	}
	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	}
	return newError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
}

func evalIntegerInfixExpression(operator string, left, right int64) object.Object {
	switch operator {
	case "+":
		return &object.Integer{Value: left + right}
	case "-":
		return &object.Integer{Value: left - right}
	case "*":
		return &object.Integer{Value: left * right}
	case "/":
		if right == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: left / right}
	case "<":
		return booleanObject(left < right)
	case ">":
		return booleanObject(left > right)
	case "<=":
		return booleanObject(left <= right)
	case ">=":
		return booleanObject(left >= right)
	case "==":
		return booleanObject(left == right)
	case "!=":
		return booleanObject(left != right)
	default:
		return newError("unknown operator: INTEGER %s INTEGER", operator)
	}
}

func equalObjects(left, right object.Object) bool {
	if l, ok := left.(*object.Boolean); ok {
		r, ok := right.(*object.Boolean)
		return ok && l.Value == r.Value
	}
	return left.Type() == object.NULL_OBJ && right.Type() == object.NULL_OBJ
}

// Only false and null are falsey; even integer zero is truthy in Monkey.
func isTruthy(value object.Object) bool {
	switch value := value.(type) {
	case *object.Null:
		return false
	case *object.Boolean:
		return value.Value
	default:
		return true
	}
}

func booleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func isControl(value object.Object) bool {
	return value != nil && (value.Type() == object.ERROR_OBJ || value.Type() == object.RETURN_VALUE_OBJ)
}

func newError(format string, args ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}
