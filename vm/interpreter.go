package vm

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	"github.com/guilhermeg2k/glang/glang"
)

func exit() {
	//fmt.Println("Eval completed")
	os.Exit(1)
}

func Evaluate(vm *VirtualMachine) {
	functions := *vm.Functions
	main := functions[0]
	reader := bytes.NewReader(main.Bytecode)
	for {
		opCode, err := reader.ReadByte()
		if err == io.EOF {
			exit()
		}
		switch opCode {
		case ILOAD_CONST:
			iLoadConst(reader, vm.Stack, main.Consts)
		case ISTORE:
			iStore(reader, vm.Stack, main.Frame)
		case ILOAD:
			iLoad(reader, vm.Stack, main.Frame)
		case IADD:
			iAdd(vm.Stack)
		case ISUB:
			iSub(vm.Stack)
		case IMULT:
			iMult(vm.Stack)
		case IDIV:
			iDiv(vm.Stack)
		case IMOD:
			iMod(vm.Stack)
		case PRINT:
			Print(vm.Stack)
		}
	}
}

func iLoadConst(reader *bytes.Reader, stack *Stack, consts Consts) {
	var offset []byte
	offset = make([]byte, 8)
	reader.Read(offset)
	offsetValue, _ := binary.Uvarint(offset)
	push(stack, consts[offsetValue])
}

func iStore(reader *bytes.Reader, stack *Stack, frame *Frame) {
	var offset []byte
	offset = make([]byte, 8)
	reader.Read(offset)
	offsetValue, _ := binary.Uvarint(offset)
	bytes := pop(stack, 8)
	store(frame, offsetValue, bytes)
}

func iLoad(reader *bytes.Reader, stack *Stack, frame *Frame) {
	var offset []byte
	offset = make([]byte, 8)
	reader.Read(offset)
	offsetValue, _ := binary.Uvarint(offset)
	bytes := load(frame, offsetValue)
	push(stack, bytes)
}

func iAdd(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := glang.IAdd(left, right)
	push(stack, res)
}

func iSub(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := glang.ISub(left, right)
	push(stack, res)
}

func iMult(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := glang.IMult(left, right)
	push(stack, res)
}

func iDiv(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := glang.IDiv(left, right)
	push(stack, res)
}

func iMod(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := glang.IMod(left, right)
	push(stack, res)
}
func iPop(stack *Stack) {

}

func Print(stack *Stack) {
	value := pop(stack, 8)
	glang.Print(value)
}
