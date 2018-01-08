Monkey
======

A Go implementation of the Monkey language as described in https://interpreterbook.com/

### Interpreter

You can run a simple monkey example file with:

```shell
> go run ./main.go --entry-file ./examples/hello-world.monkey
```

### REPL

There is a REPL (Read Eval Print Loop) available via:

```shell
> go run ./main.go
This is the monkey programming language!
Feel free to type in commands, for example: 1 + 2 + 3
To set the mode:
mode=lex
mode=parse
mode=eval
>> mode=lex
Entering lex mode
Successfully configured

>> 1 + 2 + 3
{Type:INT Literal:1}
{Type:+ Literal:+}
{Type:INT Literal:2}
{Type:+ Literal:+}
{Type:INT Literal:3}

>> mode=parse
Entering parse mode
Successfully configured

>> 1 + 2 + 3
((1 + 2) + 3)

>> mode=eval
Entering eval mode
Successfully configured

>> 1 + 2 + 3
6

>> exit
Exiting...
```

### Testing

Running all tests

```shell
go test ./...
```

Running a particular folder's tests:

```shell
go test ./lexer
```

To run a particular test:

```shell
go test ./parser -run TestIfStatement
```

The testing approach chosen is similar to [Golden Master Testing](https://en.wikipedia.org/wiki/Characterization_test).
A test's output is recorded, and used as a reference to compare future test runs against. If the output has changed, the
test will fail. This is similar to Jest's [snapshot testing](https://facebook.github.io/jest/docs/en/snapshot-testing.html).

If you wish to re-record snapshot tests:


```shell
UPDATE_SNAPSHOTS=true go test ./...
```

There is also REPL tests available:

```shell
> ./test-repl.sh
...
Tests passed successfully
```

### Parsing

The parsing algorithm used is based on the "Top Down operator Precedence" work of Vaughan Pratt. Although this parser
implementation doesn't follow the specific terminology of nud/led.

In particular, the prefix parsing functions are "nuds" for null denotations, and infix parse functions are "leds" for
left denotations. Leds care about _left_ expressions, whilst a nud does not mind what came before it.

Related work:
- https://web.archive.org/web/20151223215421/http://hall.org.ua/halls/wizzard/pdf/Vaughan.Pratt.TDOP.pdf
- http://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
- https://eli.thegreenplace.net/2010/01/02/top-down-operator-precedence-parsing
- https://crockford.com/javascript/tdop/tdop.html

### Take aways

Go supports "enums" via https://github.com/golang/go/wiki/Iota, however if you want a Stringer implementation this can
be generated for you via a comment:

```go
//go:generate stringer -type=ObjectType
```

However you require a manual call to generate these files:

```shell
> go generate
```

Further details: https://blog.golang.org/generate

### Dependencies

Dependency management is handled by [dep](https://github.com/golang/dep)
