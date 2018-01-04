package repl

import (
	"io"
	"bufio"
	"fmt"
	"github.com/alanfoster/monkey/lexer"
	"github.com/alanfoster/monkey/parser"
	"github.com/alanfoster/monkey/token"
)

const PROMPT = ">> "

type Mode int

const (
	_       Mode = iota
	LEXING
	PARSING
)

type Repl struct {
	Mode Mode
}

func (r *Repl) Configure(line string) bool {
	if line == "mode=lexing" {
		fmt.Println("Entering lexing mode")
		r.Mode = LEXING
		return true
	} else if line == "mode=parsing" {
		fmt.Println("Entering parsing mode")
		r.Mode = PARSING
		return true
	}

	return false
}

func (r *Repl) Handle(line string) {
	switch r.Mode {
	case LEXING:
		r.lex(line)
	case PARSING:
		r.parse(line)
	}
}

func (r *Repl) lex(line string) {
	l := lexer.New(line)
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
}

func (r *Repl) parse(line string) {
	l := lexer.New(line)
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()

	if len(errors) != 0 {
		r.printParsingErrors(errors)
		return
	}

	fmt.Println(program.PrettyPrint())
}

func (r *Repl) printParsingErrors(errors []string) {
	fmt.Println("Error: Parsing errors found.")
	for _, error := range errors {
		fmt.Printf("%v\n", error)
	}
}

func Start(in io.Reader, out io.Reader) {
	scanner := bufio.NewScanner(in)
	repl := Repl{
		Mode: PARSING,
	}

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

		if repl.Configure(line) {
			fmt.Println("Successfully configured\n")
			continue
		}

		repl.Handle(line)
		fmt.Println()
	}
}
