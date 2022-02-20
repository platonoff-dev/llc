package main

import (
	"anubis/repl"
	"fmt"
	"os"
	"os/user"
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
