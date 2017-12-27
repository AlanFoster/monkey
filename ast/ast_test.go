package ast

import (
	"testing"
	"github.com/alanfoster/monkey/token"
	"github.com/stretchr/testify/assert"
)

func TestPrettyPrint(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "anotherVar",
					},
					Value: "anotherVar",
				},
			},
		},
	}

	expectedString := "let myVar = anotherVar;"
	assert.Equal(t, expectedString, program.PrettyPrint())
}
