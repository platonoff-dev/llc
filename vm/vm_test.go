package vm

import (
	"fmt"
	"testing"

	"anubis/ast"
	"anubis/compiler"
	"anubis/lexer"
	"anubis/object"
	"anubis/parser"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{input: "1", expected: 1},
		{input: "2", expected: 2},
		{input: "1 + 2", expected: 3},
		{input: "1 - 2", expected: -1},
		{input: "4 / 2", expected: 2},
		{input: "50 / 2 * 2 + 10 - 5", expected: 55},
		{input: "5 + 5 + 5 + 5 - 10", expected: 10},
		{input: "2 * 2 * 2 * 2 * 2", expected: 32},
		{input: "5 * 2 + 10", expected: 20},
		{input: "5 + 2 * 10", expected: 25},
		{input: "5 * (2 + 10)", expected: 60},
		{input: "-5", expected: -5},
		{input: "-10", expected: -10},
		{input: "-50 + 100 + -50", expected: 0},
		{input: "(5 + 10 * 2 + 15 / 3) * 2 + -10", expected: 50},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: "true", expected: true},
		{input: "false", expected: false},
		{input: "1 < 2", expected: true},
		{input: "1 > 2", expected: false},
		{input: "1 < 1", expected: false},
		{input: "1 > 1", expected: false},
		{input: "1 == 1", expected: true},
		{input: "1 != 1", expected: false},
		{input: "1 == 2", expected: false},
		{input: "1 != 2", expected: true},
		{input: "true == true", expected: true},
		{input: "false == false", expected: true},
		{input: "true == false", expected: false},
		{input: "true != false", expected: true},
		{input: "false != true", expected: true},
		{input: "(1 < 2) == true", expected: true},
		{input: "(1 < 2) == false", expected: false},
		{input: "(1 > 2) == true", expected: false},
		{input: "(1 > 2) == false", expected: true},
		{input: "!true", expected: false},
		{input: "!false", expected: true},
		{input: "!5", expected: false},
		{input: "!!true", expected: true},
		{input: "!!false", expected: false},
		{input: "!!5", expected: true},
	}

	runVmTests(t, tests)
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	reuslt, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T, (%+v)", actual, actual)
	}

	if reuslt.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", reuslt.Value, expected)
	}

	return nil
}

type vmTestCase struct {
	expected interface{}
	input    string
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for i, tt := range tests {
		name := fmt.Sprintf("[%d]", i)
		t.Run(name, func(t *testing.T) {
			program := parse(tt.input)
			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			vm := New(comp.Bytecode())
			err = vm.Run()
			if err != nil {
				t.Fatalf("vm error: %s", err)
			}

			stackElem := vm.LastPoppedStackElem()
			testExpectedObject(t, tt.expected, stackElem)
		})
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	}
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}
