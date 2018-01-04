#!/bin/sh
set -xe

go build main.go

expect -f -<<EOF
    set timeout 2

    spawn "./main"
    expect ">> "

    ## Testing Lexing
    send "mode=lexing\n"
    expect "Entering lexing mode"
    send "1 + 2 + 3\n"
    expect "{Type:INT Literal:1}"
    expect "{Type:+ Literal:+}"
    expect "{Type:INT Literal:2}"
    expect "{Type:+ Literal:+}"
    expect "{Type:INT Literal:3}"

    ## Testing Parsing
    send "mode=parsing\n"
    expect "Entering parsing mode"
    send "1 + 2 + 3\n"
    expect "((1 + 2) + 3)"

    ## Testing Parsing Failures
    send "let foo \n"
    expect "expected next token to be =, but got {EOF } instead"

    ## Exiting
    send "exit\n"
    expect {
        eof { puts "Tests passed successfully" }
        timeout { puts "Error:: Timed out"; exit 1 }
    }
EOF