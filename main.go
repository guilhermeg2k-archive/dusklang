package main

import (
	"fmt"
	"os"
)

func main() {
	var program Program
	lexer, err := newLexerFromFile("tokens.lex")
	handleError(err)
	tok, err := lexer.testTokensOfFile("ex.g")
	fmt.Println(tok)
	handleError(err)
	lexer.TokenTable = tok
	program, err = Parse(lexer)
	fmt.Println(program)
	handleError(err)
	//b, err := json.Marshal(program)
	//fmt.Println(string(b))
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
