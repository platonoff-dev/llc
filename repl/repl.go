package repl

import (
	"bufio"
	"fmt"
	"io"

	"llc/compiler"
	"llc/lexer"
	"llc/parser"
	"llc/vm"
)

const PROMPT = ">>>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	// env, err := files.ReadStd()
	// if err != nil {
	// 	panic(err)
	// }

	for {
		fmt.Printf("%s ", PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// evaluated := evaluator.Eval(program, env)
		// if evaluated != nil {
		// 	_, _ = io.WriteString(out, evaluated.Inspect())
		// 	_, _ = io.WriteString(out, "\n")
		// }

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			_, _ = fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			_, _ = fmt.Fprintf(out, "Woops! Execution bytecode failed:\n %s\n", err)
		}

		stackTop := machine.LastPoppedStackElem()
		_, _ = io.WriteString(out, stackTop.Inspect())
		_, _ = io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
