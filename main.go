package main

import (
	"fmt"
	"github.com/alanfoster/monkey/repl"
	"os"
)

func main() {
	fmt.Println("This is the monkey programming language!")
	fmt.Println("Feel free to type in commands, for example: 1 + 2 + 3")
	fmt.Println("To set the mode:")
	fmt.Println("mode=lexing")
	fmt.Println("mode=parsing")
	repl.Start(os.Stdin, os.Stdout)
}
