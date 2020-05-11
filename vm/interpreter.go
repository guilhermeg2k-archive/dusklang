package vm

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	"github.com/guilhermeg2k/dusklang/dusk"
)

func exit() {
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
		//INT ARITHIMETICS OPERATIONS
		case ILOAD_CONST:
			iLoadConst(reader, vm.Stack, main.Consts)
		case ISTORE:
			iStore(reader, vm.Stack, main.Storage)
		case ILOAD:
			iLoad(reader, vm.Stack, main.Storage)
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
		//FLOAT ARITHIMETICS OPERATIONS
		case FLOAD_CONST:
			fLoadConst(reader, vm.Stack, main.Consts)
		case FSTORE:
			fStore(reader, vm.Stack, main.Storage)
		case FLOAD:
			fLoad(reader, vm.Stack, main.Storage)
		case FADD:
			fAdd(vm.Stack)
		case FSUB:
			fSub(vm.Stack)
		case FMULT:
			fMult(vm.Stack)
		case FDIV:
			fDiv(vm.Stack)
		//INT COMPARISONS
		case ICMP_EQUALS:
			iCmpEquals(vm.Stack)
		case ICMP_LESS_EQUALS:
			iCmpLessEquals(vm.Stack)
		case ICMP_GREATER_EQUALS:
			iCmpGreaterEquals(vm.Stack)
		case ICMP_LESS_THEN:
			iCmpLessThen(vm.Stack)
		case ICMP_GREATER_THEN:
			iCmpLessThen(vm.Stack)
		//FLOAT COMPARISONS
		case FCMP_EQUALS:
			fCmpEquals(vm.Stack)
		case FCMP_LESS_EQUALS:
			fCmpLessEquals(vm.Stack)
		case FCMP_GREATER_EQUALS:
			fCmpGreaterEquals(vm.Stack)
		case FCMP_LESS_THEN:
			fCmpLessThen(vm.Stack)
		case FCMP_GREATER_THEN:
			fCmpGreaterThen(vm.Stack)
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

func iStore(reader *bytes.Reader, stack *Stack, frame *Storage) {
	var offset []byte
	offset = make([]byte, 8)
	reader.Read(offset)
	offsetValue, _ := binary.Uvarint(offset)
	bytes := pop(stack, 8)
	store(frame, offsetValue, bytes)
}

func iLoad(reader *bytes.Reader, stack *Stack, frame *Storage) {
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
	res := dusk.IAdd(left, right)
	push(stack, res)
}

func iSub(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ISub(left, right)
	push(stack, res)
}

func iMult(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.IMult(left, right)
	push(stack, res)
}

func iDiv(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.IDiv(left, right)
	push(stack, res)
}

func iMod(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.IMod(left, right)
	push(stack, res)
}

func fLoadConst(reader *bytes.Reader, stack *Stack, consts Consts) {
	var offset []byte
	offset = make([]byte, 8)
	reader.Read(offset)
	offsetValue, _ := binary.Uvarint(offset)
	push(stack, consts[offsetValue])
}

func fStore(reader *bytes.Reader, stack *Stack, frame *Storage) {
	var offset []byte
	offset = make([]byte, 8)
	reader.Read(offset)
	offsetValue, _ := binary.Uvarint(offset)
	bytes := pop(stack, 8)
	store(frame, offsetValue, bytes)
}

func fLoad(reader *bytes.Reader, stack *Stack, frame *Storage) {
	var offset []byte
	offset = make([]byte, 8)
	reader.Read(offset)
	offsetValue, _ := binary.Uvarint(offset)
	bytes := load(frame, offsetValue)
	push(stack, bytes)
}

func fAdd(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FAdd(left, right)
	push(stack, res)
}

func fSub(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FSub(left, right)
	push(stack, res)
}

func fMult(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FMult(left, right)
	push(stack, res)
}

func fDiv(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FDiv(left, right)
	push(stack, res)
}

func iCmpEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpEquals(left, right)
	push(stack, res)
}

func iCmpLessEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpLessEquals(left, right)
	push(stack, res)
}

func iCmpGreaterEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpGreaterEquals(left, right)
	push(stack, res)
}

func iCmpLessThen(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpLessThen(left, right)
	push(stack, res)
}
func iCmpGreaterThen(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpGreaterThen(left, right)
	push(stack, res)
}

func fCmpEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpEquals(left, right)
	push(stack, res)
}

func fCmpLessEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpLessEquals(left, right)
	push(stack, res)
}

func fCmpGreaterEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpGreaterEquals(left, right)
	push(stack, res)
}

func fCmpLessThen(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpLessThen(left, right)
	push(stack, res)
}
func fCmpGreaterThen(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpGreaterThen(left, right)
	push(stack, res)
}

//Temporary
func Print(stack *Stack) {
	value := pop(stack, 1)
	dusk.Print(value)
}
