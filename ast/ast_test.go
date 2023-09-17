package ast

import (
	"testing"

	"github.com/gkampitakis/monkey/token"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: []byte("let")},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: []byte("myVar")},
					Value: []byte("myVar"),
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: []byte("anotherVar")},
					Value: []byte("anotherVar"),
				},
			},
		},
	}

	require.Equal(t, program.String(), "let myVar = anotherVar;")
}
