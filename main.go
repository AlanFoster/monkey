package main

import (
	"flag"
	"github.com/alanfoster/monkey/repl"
	"os"
	"io"
	"io/ioutil"
	"github.com/alanfoster/monkey/lexer"
	"github.com/alanfoster/monkey/evaluator"
	"fmt"
	"github.com/alanfoster/monkey/parser"
	"github.com/alanfoster/monkey/object"
)

func printParsingErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Error: Parsing errors found.")
	for _, e := range errors {
		fmt.Fprintf(out, "%v\n", e)
	}
}

func interpretFile(path string, out io.Writer) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("oops")
		panic(err)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()

	if len(errors) != 0 {
		printParsingErrors(out, errors)
		return
	}

	environment := object.NewEnvironment()
	evaluator.Eval(program, environment)
}

func main() {
	var entryFile string
	flag.StringVar(&entryFile, "entry-file", "", "File to run as a monkey file program")
	flag.Parse()

	if entryFile != "" {
		interpretFile(entryFile, os.Stdout)
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
