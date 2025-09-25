package evaluator

import (
	"fmt"

	"llc/lang/ast"
	"llc/lang/object"
	"llc/lang/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCall(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCall(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		uquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(uquoted)
	})
}

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.Int,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.True, Literal: "true"}
		} else {
			t = token.Token{Type: token.False, Literal: "false"}
		}

		return &ast.Boolean{Token: t, Value: obj.Value}
	case *object.Quote:
		return obj.Node
	default:
		return nil
	}
}

func isUnquoteCall(node ast.Node) bool {
	call, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return call.Function.TokenLiteral() == "unquote"
}
