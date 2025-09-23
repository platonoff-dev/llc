package files

import (
	"fmt"
	"os"

	"llc/lang/evaluator"
	"llc/lang/lexer"
	"llc/lang/object"
	"llc/lang/parser"
)

func ReadStd() (*object.Environment, error) {
	folder := "std"
	files := []string{
		"array.an",
	}
	env := object.NewEnvironment()

	for _, file := range files {
		path := folder + "/" + file
		_, err := ReadFile(path, env)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

func ReadFile(path string, env *object.Environment) (*object.Environment, error) {
	sourceCode, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	l := lexer.New(string(sourceCode))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, fmt.Errorf("error hapend during parsing of %s. errors=%v", path, p.Errors())
	}

	_ = evaluator.Eval(program, env)

	return env, nil
}
