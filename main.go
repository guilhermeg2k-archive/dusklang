package main

import (
	"fmt"
	"os"
)

func main() {
	lexer, err := newLexerFromFile("tokens.lex")
	handleError(err)
	tok, err := lexer.testTokensOfFile("ex.g")
	handleError(err)
	lexer.TokenTable = tok
	program, err := Parse(lexer)
	handleError(err)
	fmt.Println(program)
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
