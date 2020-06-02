package vm

import (
	"fmt"
	"testing"
)

func TestJumpElse(t *testing.T) {
	virtualMachine := VirtualMachine{
		Stack: &Stack{},
	}
	function := Function{
		Storage: Storage{},
	}
	consts := make(Consts)
	consts[0] = []byte{0}
	consts[1] = GetInt(200)
	labels := make(Labels)
	labels[0] = 27

	//TEST IF FALSE
	function.Bytecode = []byte{15} //LOAD CONST POS 0
	function.Bytecode = append(function.Bytecode, GetInt(0)...)
	function.Bytecode = append(function.Bytecode, 29) //JUMP IF FALSE TO LABEL 0
	function.Bytecode = append(function.Bytecode, GetInt(0)...)
	function.Bytecode = append(function.Bytecode, 0) // LOAD CONST AT POS 1
	function.Bytecode = append(function.Bytecode, GetInt(1)...)
	function.Bytecode = append(function.Bytecode, 255)

	function.Consts = consts
	function.Labels = labels
	returned := Evaluate(&virtualMachine, &function)
	if len(returned) != 0 {
		t.Errorf(fmt.Sprintf("want %d, got %d", 0, returned[0]))
	}
	//TEST IF TRUE
	consts[0] = []byte{1}
	function.Bytecode = []byte{15} //LOAD CONST POS 0
	function.Bytecode = append(function.Bytecode, GetInt(0)...)
	function.Bytecode = append(function.Bytecode, 29) //JUMP IF FALSE TO LABEL 0
	function.Bytecode = append(function.Bytecode, GetInt(0)...)
	function.Bytecode = append(function.Bytecode, 0) // LOAD CONST AT POS 1
	function.Bytecode = append(function.Bytecode, GetInt(1)...)
	function.Bytecode = append(function.Bytecode, 255)

	function.Consts = consts
	function.Labels = labels
	returned = Evaluate(&virtualMachine, &function)
	if len(returned) == 0 {
		t.Errorf(fmt.Sprintf("want a value, got nothing"))
	}
}

func TestJumpTrue(t *testing.T) {
	virtualMachine := VirtualMachine{
		Stack: &Stack{},
	}
	function := Function{
		Storage: Storage{},
	}
	consts := make(Consts)
	consts[0] = []byte{0}
	consts[1] = GetInt(200)
	labels := make(Labels)
	labels[0] = 27

	//TEST IF FALSE
	function.Bytecode = []byte{15} //LOAD CONST POS 0
	function.Bytecode = append(function.Bytecode, GetInt(0)...)
	function.Bytecode = append(function.Bytecode, 30) //JUMP IF TRUE TO LABEL 0
	function.Bytecode = append(function.Bytecode, GetInt(0)...)
	function.Bytecode = append(function.Bytecode, 0) // LOAD CONST AT POS 1
	function.Bytecode = append(function.Bytecode, GetInt(1)...)
	function.Bytecode = append(function.Bytecode, 255)

	function.Consts = consts
	function.Labels = labels
	returned := Evaluate(&virtualMachine, &function)
	if len(returned) == 0 {
		t.Errorf(fmt.Sprintf("want a value, got nothing"))
	}

	//TEST IF TRUE
	consts[0] = []byte{1}
	function.Bytecode = []byte{15} //LOAD CONST POS 0
	function.Bytecode = append(function.Bytecode, GetInt(0)...)
	function.Bytecode = append(function.Bytecode, 30) //JUMP IF TRUE TO LABEL 0
	function.Bytecode = append(function.Bytecode, GetInt(0)...)
	function.Bytecode = append(function.Bytecode, 0) // LOAD CONST AT POS 1
	function.Bytecode = append(function.Bytecode, GetInt(1)...)
	function.Bytecode = append(function.Bytecode, 255)

	function.Consts = consts
	function.Labels = labels
	virtualMachine.Stack = &Stack{}
	returned = Evaluate(&virtualMachine, &function)
	if len(returned) != 0 {
		t.Errorf(fmt.Sprintf("want %d, got %d", 0, returned[0]))
	}
}
