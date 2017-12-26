package main

import (
	"fmt"
	"github.com/alanfoster/monkey/repl"
	"os"
)

func main() {
	fmt.Println("This is the monkey programming language!")
	fmt.Println("Feel free to type in commands, for example: 1 + 3")
	repl.Start(os.Stdin, os.Stdout)
}
