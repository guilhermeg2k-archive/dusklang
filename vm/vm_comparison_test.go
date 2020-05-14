package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

type FloatCase struct {
	a      float64
	b      float64
	opCode byte
	want   interface{}
}

type IntCase struct {
	a      int64
	b      int64
	opCode byte
	want   interface{}
}

func TestFloatComparison(t *testing.T) {
	floatCases := []FloatCase{
		//EQUALS
		FloatCase{
			a:      1.2,
			b:      2.3,
			opCode: 24,
			want:   byte(0),
		},
		FloatCase{
			a:      1.2,
			b:      1.2,
			opCode: 24,
			want:   byte(1),
		},
		FloatCase{
			a:      2.3,
			b:      1.2,
			opCode: 24,
			want:   byte(0),
		},
		//LESS QUALS
		FloatCase{
			a:      1.1,
			b:      1.2,
			opCode: 25,
			want:   byte(1),
		},
		FloatCase{
			a:      1.3,
			b:      1.2,
			opCode: 25,
			want:   byte(0),
		},
		FloatCase{
			a:      1.2,
			b:      1.2,
			opCode: 25,
			want:   byte(1),
		},
		//GREATER EQUALS
		FloatCase{
			a:      1.1,
			b:      1.2,
			opCode: 26,
			want:   byte(0),
		},
		FloatCase{
			a:      1.3,
			b:      1.2,
			opCode: 26,
			want:   byte(1),
		},
		FloatCase{
			a:      1.2,
			b:      1.2,
			opCode: 26,
			want:   byte(1),
		},
		//LESS THEN
		FloatCase{
			a:      1.1,
			b:      1.2,
			opCode: 27,
			want:   byte(1),
		},
		FloatCase{
			a:      1.3,
			b:      1.2,
			opCode: 27,
			want:   byte(0),
		},
		FloatCase{
			a:      1.2,
			b:      1.2,
			opCode: 27,
			want:   byte(0),
		},
		//GREATER THEN
		FloatCase{
			a:      1.1,
			b:      1.2,
			opCode: 28,
			want:   byte(0),
		},
		FloatCase{
			a:      1.3,
			b:      1.2,
			opCode: 28,
			want:   byte(1),
		},
		FloatCase{
			a:      1.2,
			b:      1.2,
			opCode: 28,
			want:   byte(0),
		},
	}
	for i, fc := range floatCases {
		virtualMachine := VirtualMachine{
			Stack: &Stack{},
		}
		function := Function{
			Storage: &Storage{},
		}
		consts := make(Consts)
		consts[0] = GetFloat(fc.a)
		consts[1] = GetFloat(fc.b)

		function.Bytecode = []byte{8}
		function.Bytecode = append(function.Bytecode, GetInt(0)...)
		function.Bytecode = append(function.Bytecode, 8)
		function.Bytecode = append(function.Bytecode, GetInt(1)...)
		function.Bytecode = append(function.Bytecode, fc.opCode) // OPERATION
		function.Bytecode = append(function.Bytecode, 255)

		function.Consts = consts
		returned := Evaluate(&virtualMachine, &function)
		if returned[0] != fc.want.(byte) {
			t.Error(fmt.Sprintf("case: %d, op: %d, want %d, got %b", i, fc.opCode, fc.want.(byte), returned[0]))
		}
	}
}

func TestIntComparison(t *testing.T) {
	intCases := []IntCase{
		//EQUALS
		IntCase{
			a:      1,
			b:      2,
			opCode: 24,
			want:   byte(0),
		},
		IntCase{
			a:      1,
			b:      1,
			opCode: 24,
			want:   byte(1),
		},
		IntCase{
			a:      2,
			b:      1,
			opCode: 24,
			want:   byte(0),
		},
		//LESS QUALS
		IntCase{
			a:      1,
			b:      2,
			opCode: 25,
			want:   byte(1),
		},
		IntCase{
			a:      1,
			b:      1,
			opCode: 25,
			want:   byte(1),
		},
		IntCase{
			a:      2,
			b:      1,
			opCode: 25,
			want:   byte(0),
		},
		//GREATER EQUALS
		IntCase{
			a:      1,
			b:      1,
			opCode: 26,
			want:   byte(1),
		},
		IntCase{
			a:      2,
			b:      1,
			opCode: 26,
			want:   byte(1),
		},
		IntCase{
			a:      1,
			b:      2,
			opCode: 26,
			want:   byte(0),
		},
		//LESS THEN
		IntCase{
			a:      1,
			b:      2,
			opCode: 27,
			want:   byte(1),
		},
		IntCase{
			a:      1,
			b:      1,
			opCode: 27,
			want:   byte(0),
		},
		IntCase{
			a:      2,
			b:      1,
			opCode: 27,
			want:   byte(0),
		},
		//GREATER THEN
		IntCase{
			a:      1,
			b:      1,
			opCode: 28,
			want:   byte(0),
		},
		IntCase{
			a:      2,
			b:      1,
			opCode: 28,
			want:   byte(1),
		},
		IntCase{
			a:      1,
			b:      2,
			opCode: 28,
			want:   byte(0),
		},
	}
	for i, ic := range intCases {
		virtualMachine := VirtualMachine{
			Stack: &Stack{},
		}
		function := Function{
			Storage: &Storage{},
		}
		consts := make(Consts)
		consts[0] = GetInt(ic.a)
		consts[1] = GetInt(ic.b)

		function.Bytecode = []byte{8}
		function.Bytecode = append(function.Bytecode, GetInt(0)...)
		function.Bytecode = append(function.Bytecode, 8)
		function.Bytecode = append(function.Bytecode, GetInt(1)...)
		function.Bytecode = append(function.Bytecode, ic.opCode) // OPERATION
		function.Bytecode = append(function.Bytecode, 255)

		function.Consts = consts
		returned := Evaluate(&virtualMachine, &function)
		if returned[0] != ic.want.(byte) {
			t.Error(fmt.Sprintf("case: %d, op: %d, want %d, got %b", i, ic.opCode, ic.want.(byte), returned[0]))
		}
	}

}

func GetInt(i int64) []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, i)
	return buffer.Bytes()
}

func GetFloat(f float64) []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, f)
	return buffer.Bytes()
}
