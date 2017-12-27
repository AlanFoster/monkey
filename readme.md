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

### Testing

Running all tests
```shell
go test ./...
```

Running a particular folder's tests:

```shell
go test ./lexer
```

### Dependencies

Dependency management is handled by [dep](https://github.com/golang/dep)
