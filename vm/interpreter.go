package vm

import (
	"fmt"
	"io"
	"os"
)

func exit() {
	os.Exit(1)
}

func Evaluate(vm *VirtualMachine, function *Function) []byte {
	function.CurrentOffset = 0
	for {
		opCode, err := function.readByte()
		fmt.Println(opCode)
		if err == io.EOF {
			return nil
		}
		switch opCode {

		//Int arithmetics
		case ILOAD_CONST:
			iLoadConst(vm.Stack, function)
		case ISTORE:
			iStore(vm.Stack, function)
		case ILOAD:
			iLoad(vm.Stack, function)
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
		//Float arithmetics
		case FLOAD_CONST:
			fLoadConst(vm.Stack, function)
		case FSTORE:
			fStore(vm.Stack, function)
		case FLOAD:
			fLoad(vm.Stack, function)
		case FADD:
			fAdd(vm.Stack)
		case FSUB:
			fSub(vm.Stack)
		case FMULT:
			fMult(vm.Stack)
		case FDIV:
			fDiv(vm.Stack)
		case BOLOAD_CONST:
			bLoadConst(vm.Stack, function)
		case BOSTORE:
			bStore(vm.Stack, function)
		case BOLOAD:
			bLoad(vm.Stack, function)
		//Int comparisons
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
		//Float comparisons
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
		//Jumps
		case JUMP_IF_ELSE:
			jumpIfElse(vm.Stack, function)
		case JUMP_IF_TRUE:
			jumpIfTrue(vm.Stack, function)
		case PRINT:
			Print(vm.Stack)
		case 255:
			return *vm.Stack
		}
	}
}

//TODO: HANDLE ERRORS

func handleError(err error) {
	fmt.Println("PANIC: ", err.Error())
	os.Exit(0)
}
