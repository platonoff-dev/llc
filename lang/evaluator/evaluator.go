package evaluator

import (
	"fmt"

	"llc/lang/ast"
	object2 "llc/lang/object"
)

var (
	TRUE  = &object2.Boolean{Value: true}
	FALSE = &object2.Boolean{Value: false}
	NULL  = &object2.Null{}
)

func Eval(node ast.Node, env *object2.Environment) object2.Object { //nolint:gocognit,cyclop,funlen
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object2.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object2.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object2.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object2.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object2.Function{Parameters: params, Env: env, Body: body}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	}

	return nil
}

func evalIndexExpression(left, index object2.Object) object2.Object {
	switch {
	case left.Type() == object2.ArrayObj && index.Type() == object2.IntegerObj:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object2.HashObj:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalHashIndexExpression(hash, index object2.Object) object2.Object {
	hashObj, _ := hash.(*object2.Hash)

	key, ok := index.(object2.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalArrayIndexExpression(array, index object2.Object) object2.Object {
	arrayObj, _ := array.(*object2.Array)
	idxInt, _ := index.(*object2.Integer)
	maxIdx := int64(len(arrayObj.Elements) - 1)

	if idxInt.Value < 0 || idxInt.Value > maxIdx {
		return NULL
	}

	return arrayObj.Elements[idxInt.Value]
}

func evalHashLiteral(node *ast.HashLiteral, env *object2.Environment) object2.Object {
	pairs := make(map[object2.HashKey]object2.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object2.Hashable)
		if !ok {
			newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object2.HashPair{Key: key, Value: value}
	}

	return &object2.Hash{Pairs: pairs}
}

func applyFunction(fn object2.Object, args []object2.Object) object2.Object {
	switch fn := fn.(type) {
	case *object2.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object2.Builtin:
		return fn.Function(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object2.Function, args []object2.Object) *object2.Environment {
	env := object2.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object2.Object) object2.Object {
	if returnValue, ok := obj.(*object2.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalExpressions(exps []ast.Expression, env *object2.Environment) []object2.Object {
	var result []object2.Object //nolint: prealloc

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object2.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object2.Environment) object2.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) { //nolint:gocritic
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object2.Environment) object2.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("%s", "identifier not found: "+node.Value)
}

func isTruthy(obj object2.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isError(obj object2.Object) bool {
	if obj != nil {
		return obj.Type() == object2.ErrorObj
	}
	return false
}

func evalInfixExpression(operator string, left, right object2.Object) object2.Object {
	switch {
	case left.Type() == object2.IntegerObj && right.Type() == object2.IntegerObj:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object2.StringObj && right.Type() == object2.StringObj:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object2.Object) object2.Object {
	leftStmt, _ := left.(*object2.String)
	rightStmt, _ := right.(*object2.String)

	switch operator { //nolint: gocritic
	case "+":
		return &object2.String{Value: leftStmt.Value + rightStmt.Value}
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalIntegerInfixExpression(operator string, left, right object2.Object) object2.Object {
	leftVal, _ := left.(*object2.Integer)
	rightVal, _ := right.(*object2.Integer)

	switch operator {
	case "+":
		return &object2.Integer{Value: leftVal.Value + rightVal.Value}
	case "-":
		return &object2.Integer{Value: leftVal.Value - rightVal.Value}
	case "*":
		return &object2.Integer{Value: leftVal.Value * rightVal.Value}
	case "/":
		return &object2.Integer{Value: leftVal.Value / rightVal.Value}
	case "<":
		return nativeBoolToBooleanObject(leftVal.Value < rightVal.Value)
	case ">":
		return nativeBoolToBooleanObject(leftVal.Value > rightVal.Value)
	case "==":
		return nativeBoolToBooleanObject(leftVal.Value == rightVal.Value)
	case "!=":
		return nativeBoolToBooleanObject(leftVal.Value != rightVal.Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object2.Object) object2.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object2.Object) object2.Object {
	if right.Type() != object2.IntegerObj {
		return newError("unknown operator: -%s", right.Type())
	}

	value, _ := right.(*object2.Integer)
	return &object2.Integer{Value: -value.Value}
}

func evalBangOperatorExpression(right object2.Object) object2.Object {
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

func nativeBoolToBooleanObject(input bool) *object2.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalProgram(program *ast.Program, env *object2.Environment) object2.Object {
	var result object2.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object2.ReturnValue:
			return result.Value
		case *object2.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object2.Environment) object2.Object {
	var result object2.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object2.ReturnValueObj || rt == object2.ErrorObj {
				return result
			}
			return result
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object2.Error {
	return &object2.Error{Message: fmt.Sprintf(format, a...)}
}
