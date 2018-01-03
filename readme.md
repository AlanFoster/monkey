Monkey
======

A Go implementation of the Monkey language as described in https://interpreterbook.com/

### REPL

There is a REPL (Read Eval Print Loop) available via:

```shell
go run ./main.go

This is the monkey programming language!
Feel free to type in commands, for example: 1 + 3
>> (1 + 3) > 2
{Type:( Literal:(}
{Type:INT Literal:1}
{Type:+ Literal:+}
{Type:INT Literal:3}
{Type:) Literal:)}
{Type:> Literal:>}
{Type:INT Literal:2
>>
```

### Parsing

The parsing algorithm used is based on the "Top Down operator Precedence" work of Vaughan Pratt. Although this parser
implementation doesn't follow the specific terminology of nud/led/std, etc.

In particular, the prefix parsing functions are "nuds" for null denotations, and infix parse functions are "leds" for
left denotations.

Related work:
- http://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
- https://eli.thegreenplace.net/2010/01/02/top-down-operator-precedence-parsing
- https://crockford.com/javascript/tdop/tdop.html

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

### Dependencies

Dependency management is handled by [dep](https://github.com/golang/dep)
