package main

import (
	"fmt"
)

func main() {
	lexer, err := newLexerFromFile("tokens.lex")
	if err != nil {
		fmt.Println(err)
	}
	tok, err := lexer.testTokensOfFile("syntax.g")
	for _, t := range tok {
		fmt.Println(t)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}
