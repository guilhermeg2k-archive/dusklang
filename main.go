package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/guilhermeg2k/glang/vm"
)

func main() {
	/* var program ast.Program
	l, err := lexer.NewLexerFromFile("lexer/tokens.lex")
	handleError(err)
	tok, err := l.TestTokens("examples/test.g")
	handleError(err)
	l.TokenTable = tok
	program, err = parser.Parse(l)
	handleError(err)
	err = parser.Analyze(&program)
	if err != nil {
		handleError(err)
	} */
	virtualMachine := vm.VirtualMachine{
		Stack: &vm.Stack{},
	}
	main := vm.Function{}
	consts := make(vm.Consts)
	frame := make(vm.Frame)
	var a []byte
	var b []byte
	a = make([]byte, 8)
	b = make([]byte, 8)
	c := make([]byte, 8)
	d := make([]byte, 8)
	binary.PutVarint(a, 400)
	binary.PutVarint(b, 3)
	binary.PutUvarint(c, 0)
	binary.PutUvarint(d, 1)
	consts[0] = a
	consts[1] = b
	main.Consts = consts
	main.Frame = &frame
	main.Bytecode = []byte{0}
	main.Bytecode = append(main.Bytecode, c...)
	main.Bytecode = append(main.Bytecode, 0)
	main.Bytecode = append(main.Bytecode, d...)
	main.Bytecode = append(main.Bytecode, 4)
	main.Bytecode = append(main.Bytecode, 28)
	funcs := []vm.Function{}
	funcs = append(funcs, main)
	virtualMachine.Functions = &funcs
	vm.Evaluate(&virtualMachine)
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
