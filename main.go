package main

import (
	"fmt"
	"os"
	"os/user"

	"anubis/repl"
)

var art = `
    /|       |\
   / |       | \
  /  |       |  \
 / /||       ||\ \
| |_||_______||_| |
\                /
 |  __      __  |
 |  \ |    | /  |
  \  \|    |/  /
   \          /
   |\        /|
   | \  ▄▄  / |
   |   \__/   |
  /            \
`

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Anubis programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	fmt.Println(art)
	repl.Start(os.Stdin, os.Stdout)
}
