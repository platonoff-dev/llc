package evaluator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"llc/lang/ast"
	"llc/lang/lexer"
	"llc/lang/object"
	"llc/lang/parser"
)

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1;
	let function = fn(x, y) { x + y };
	let mymacro = macro(x, y) { x + y };
	`

	env := object.NewEnvironment()
	program := testParseProgram(input)

	DefineMacros(program, env)

	assert.Len(t, program.Statements, 2)

	_, ok := env.Get("number")
	assert.False(t, ok, "number should not be defined")

	_, ok = env.Get("function")
	assert.False(t, ok, "function should not be defined")

	obj, ok := env.Get("mymacro")
	assert.True(t, ok, "mymacro should be defined")

	macro, ok := obj.(*object.Macro)
	assert.True(t, ok, "macro should be object")
	assert.Len(t, macro.Parameters, 2, "macro should have 2 parameters")
	assert.Equal(t, macro.Parameters[0].String(), "x")
	assert.Equal(t, macro.Parameters[1].String(), "y")

	expectedBody := "(x + y)"
	assert.Equal(t, macro.Body.String(), expectedBody)
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: `
				let infixExpression = macro() { quote(1 + 2); };
				infixExpression();
			`,
			expected: "(1 + 2)",
		},
		{
			input: `
				let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
				reverse(2 + 2, 10 - 5);
			`,
			expected: "(10 - 5) - (2 + 2)",
		},
		{
			`
            let unless = macro(condition, consequence, alternative) {
                quote(if (!(unquote(condition))) {
                    unquote(consequence);
                } else {
                    unquote(alternative);
                });
            };

            unless(10 > 5, puts("not greater"), puts("greater"));
            `,
			`if (!(10 > 5)) { puts("not greater") } else { puts("greater") }`,
		},
	}

	for _, tt := range tests {
		expected := testParseProgram(tt.expected)
		program := testParseProgram(tt.input)

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)
		assert.Equal(t, expected.String(), expanded.String())
	}
}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return program
}
