package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/guilhermeg2k/dusklang/vm"
)

func main() {
	/* var program ast.Program
	l, err := lexer.NewLexerFromFile("lexer/tokens.lex")
	handleError(err)
	tok, err := l.TestTokens("examples/test.dsk")
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
	frame := make(vm.Storage)
	var a []byte
	var b []byte
	c := make([]byte, 8)
	d := make([]byte, 8)

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, 1.2)
	a = buffer.Bytes()

	var bufferB bytes.Buffer
	binary.Write(&bufferB, binary.LittleEndian, 1.1)
	b = bufferB.Bytes()

	binary.PutUvarint(c, 0)
	binary.PutUvarint(d, 1)

	consts[0] = a
	consts[1] = b
	main.Consts = consts
	main.Storage = &frame
	main.Bytecode = []byte{8}
	main.Bytecode = append(main.Bytecode, c...)
	main.Bytecode = append(main.Bytecode, 8)
	main.Bytecode = append(main.Bytecode, d...)
	main.Bytecode = append(main.Bytecode, 26)
	main.Bytecode = append(main.Bytecode, 99)
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
