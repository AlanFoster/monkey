package repl

import (
	"io"
	"bufio"
	"github.com/alanfoster/monkey/lexer"
	"github.com/alanfoster/monkey/parser"
	"github.com/alanfoster/monkey/token"
	"github.com/alanfoster/monkey/evaluator"
	"fmt"
)

const PROMPT = ">> "

type Mode int

const (
	_     Mode = iota
	LEX
	PARSE
	EVAL
)

const (
	LEX_MODE   = "mode=lex"
	PARSE_MODE = "mode=parse"
	EVAL_MODE  = "mode=eval"
)

type Repl struct {
	Mode Mode
	out  io.Writer
}

func (r *Repl) OutputUsage() {
	fmt.Fprintln(r.out, "To set the mode:")
	fmt.Fprintln(r.out, LEX_MODE)
	fmt.Fprintln(r.out, PARSE_MODE)
	fmt.Fprintln(r.out, EVAL_MODE)
}

func (r *Repl) Configure(line string) bool {
	if line == LEX_MODE {
		fmt.Fprintln(r.out, "Entering lex mode")
		r.Mode = LEX
		return true
	} else if line == PARSE_MODE {
		fmt.Fprintln(r.out, "Entering parse mode")
		r.Mode = PARSE
		return true
	} else if line == EVAL_MODE {
		fmt.Fprintln(r.out, "Entering eval mode")
		r.Mode = EVAL
		return true
	}

	return false
}

func (r *Repl) Handle(line string) {
	switch r.Mode {
	case LEX:
		r.lex(line)
	case PARSE:
		r.parse(line)
	case EVAL:
		r.eval(line)
	}
}

func (r *Repl) lex(line string) {
	l := lexer.New(line)
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Fprintf(r.out, "%+v\n", tok)
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

	fmt.Fprintln(r.out, program.PrettyPrint())
}

func (r *Repl) eval(line string) {
	l := lexer.New(line)
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()

	if len(errors) != 0 {
		r.printParsingErrors(errors)
		return
	}

	eval := evaluator.Eval(program)
	fmt.Fprintln(r.out, eval.Inspect())
}

func (r *Repl) printParsingErrors(errors []string) {
	io.WriteString(r.out, "Error: Parsing errors found.")
	for _, e := range errors {
		fmt.Fprintf(r.out, "%v\n", e)
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	repl := Repl{
		Mode: EVAL,
		out:  out,
	}

	repl.OutputUsage()

	for {
		io.WriteString(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" || line == "exit()" {
			fmt.Fprintln(out, "Exiting...")
			break;
		}

		if repl.Configure(line) {
			fmt.Fprintln(out, "Successfully configured\n")
			continue
		}

		repl.Handle(line)
		fmt.Fprintln(out)
	}
}
