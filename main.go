package main

import (
	"fmt"
	"os"
	"os/user"

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
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the llc programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	fmt.Println(art)
	repl.Start(os.Stdin, os.Stdout)
}
