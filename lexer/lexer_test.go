package lexer

import (
	"fmt"
	"testing"
)

func TestLexer_Tokens(t *testing.T) {
	var input = `
	let a= 1>=4
	let b=2
	return a*b
`

	l := New(input)
	tokens := l.Tokens()
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}
