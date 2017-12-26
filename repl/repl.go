package repl

import (
	"io"
	"bufio"
	"fmt"
	"github.com/alanfoster/monkey/lexer"
	"github.com/alanfoster/monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Reader) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" || line == "exit()" {
			fmt.Println("Exiting...")
			break;
		}

		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
		fmt.Println()
	}
}
