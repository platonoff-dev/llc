package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }

	turnOnetoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	tests := []struct {
		input    Node
		expected Node
	}{
		{
			input:    one(),
			expected: two(),
		},
		{
			input: &Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: one()},
				},
			},
			expected: &Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: two()},
				},
			},
		},
		{
			input:    &InfixExpression{Left: one(), Operator: "+", Right: two()},
			expected: &InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			input:    &InfixExpression{Left: two(), Operator: "+", Right: one()},
			expected: &InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			input:    &PrefixExpression{Operator: "-", Right: one()},
			expected: &PrefixExpression{Operator: "-", Right: two()},
		},
		{
			input:    &IndexExpression{Left: one(), Index: one()},
			expected: &IndexExpression{Left: two(), Index: two()},
		},
		{
			input: &IfExpression{
				Condition: one(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			expected: &IfExpression{
				Condition: two(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			input:    &ReturnStatement{ReturnValue: one()},
			expected: &ReturnStatement{ReturnValue: two()},
		},
		{
			input:    &LetStatement{Value: one()},
			expected: &LetStatement{Value: two()},
		},
		{
			input: &FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			expected: &FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			input:    &ArrayLiteral{Elements: []Expression{one(), one()}},
			expected: &ArrayLiteral{Elements: []Expression{two(), two()}},
		},
	}

	for _, tt := range tests {
		modified := Modify(tt.input, turnOnetoTwo)

		equal := reflect.DeepEqual(tt.expected, modified)
		if !equal {
			t.Errorf("Modify(%v): expected %v, got %v", tt.input, tt.expected, modified)
		}
	}

	hashLiteral := &HashLiteral{
		Pairs: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}

	Modify(hashLiteral, turnOnetoTwo)

	for key, val := range hashLiteral.Pairs {
		key, _ := key.(*IntegerLiteral)
		if key.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, key.Value)
		}
		val, _ := val.(*IntegerLiteral)
		if val.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, val.Value)
		}
	}
}
