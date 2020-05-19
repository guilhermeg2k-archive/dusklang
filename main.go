package main

import (
	"fmt"
	"os"

	"github.com/guilhermeg2k/dusklang/ast"
	"github.com/guilhermeg2k/dusklang/bytecode"
	"github.com/guilhermeg2k/dusklang/lexer"
	"github.com/guilhermeg2k/dusklang/parser"
	"github.com/guilhermeg2k/dusklang/vm"
)

func main() {
	var program ast.Program
	l, err := lexer.NewLexerFromFile("lexer/tokens.lexer")
	handleError(err)
	tok, err := l.TestTokens("examples/test.dsk")
	handleError(err)
	l.TokenTable = tok
	program, err = parser.Parse(l)
	handleError(err)
	err = parser.Analyze(&program)
	if err != nil {
		handleError(err)
	}
	main := bytecode.GenerateByteCode(&program)
	VM := vm.VirtualMachine{
		Stack: &vm.Stack{},
	}
	vm.Evaluate(&VM, &main)
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
