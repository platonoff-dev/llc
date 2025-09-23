package main

import (
	"fmt"
	"os"

	"llc/repl"
)

var art = `
 _      _      _____
| |    | |    / ____|
| |    | |   | |
| |____| |___| |____
|______|_____|______|
`

func main() {
	fmt.Printf("Hello! This is the llc programming language!\n")
	fmt.Printf("Feel free to type in commands\n")
	fmt.Println(art)
	repl.Start(os.Stdin, os.Stdout)
}
