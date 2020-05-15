package bytecode

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/guilhermeg2k/dusklang/ast"
	"github.com/guilhermeg2k/dusklang/vm"
)

type VariablesOffset map[string]uint64

type Function struct {
	Consts          vm.Consts
	Labels          vm.Labels
	VariablesOffset VariablesOffset
	StorageCounter  uint64
	ConstCounter    uint64
	LabelCounter    uint64
	bytecode        []byte
}

func GenerateByteCode(program *ast.Program) ([]byte, error) {
	f := Function{
		Consts:          make(vm.Consts),
		Labels:          make(vm.Labels),
		VariablesOffset: make(VariablesOffset),
		bytecode:        []byte{},
	}
	function := program.Functions[0]
	for _, statement := range function.Statements {
		switch statement.Type {
		case "FullVarDeclaration":
			generateFullVarDeclaration(&f, statement.Statement.(ast.FullVarDeclaration))
		}
	}
	return f.bytecode, nil
}

func generateFullVarDeclaration(function *Function, fullVarDeclaration ast.FullVarDeclaration) error {
	for _, variable := range fullVarDeclaration.Variables {
		if variable.Expression != nil {
			generateExpression(function, variable.Expression)
		} else {
			switch variable.Type {
			case "int":
				function.Consts[function.ConstCounter] = GetInt(0)
				function.bytecode = append(function.bytecode, vm.ILOAD_CONST)
				function.bytecode = append(function.bytecode, GetUint(function.ConstCounter)...)
			case "float":
				function.Consts[function.ConstCounter] = GetFloat(0)
				function.bytecode = append(function.bytecode, vm.FLOAD_CONST)
				function.bytecode = append(function.bytecode, GetUint(function.ConstCounter)...)
			case "bool":
				function.Consts[function.ConstCounter] = []byte{0}
				function.bytecode = append(function.bytecode, vm.BOLOAD_CONST)
				function.bytecode = append(function.bytecode, GetUint(function.ConstCounter)...)
			}
			function.ConstCounter++
		}
		switch variable.Type {
		case "int":
			function.bytecode = append(function.bytecode, vm.ISTORE)
			function.bytecode = append(function.bytecode, GetUint(function.StorageCounter)...)
		case "float":
			function.bytecode = append(function.bytecode, vm.FSTORE)
			function.bytecode = append(function.bytecode, GetUint(function.StorageCounter)...)
		case "bool":
			function.bytecode = append(function.bytecode, vm.BOSTORE)
			function.bytecode = append(function.bytecode, GetUint(function.StorageCounter)...)
		}
		function.VariablesOffset[variable.Identifier] = function.StorageCounter
		function.StorageCounter++
	}
	return nil
}

func generateExpression(function *Function, expression ast.Expression) error {
	switch expression.GetType() {
	case "BinaryOperation":
		if !expression.(*ast.BinaryOperation).Left.(*ast.Literal).Visited {
			generateExpression(function, expression.(*ast.BinaryOperation).Left)
		}
		switch expression.(*ast.BinaryOperation).Operator {
		case "+":
			generateExpression(function, expression.(*ast.BinaryOperation).Right)
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.IADD)
			case "decimalNumber":
				function.bytecode = append(function.bytecode, vm.FADD)
			}
		case "-":
			generateExpression(function, expression.(*ast.BinaryOperation).Right)
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.ISUB)
			case "decimalNumber":
				function.bytecode = append(function.bytecode, vm.FSUB)
			}
		case "*":
			switch expression.(*ast.BinaryOperation).Right.GetType() {
			case "BinaryOperation":
				generateExpression(function, expression.(*ast.BinaryOperation).Right.(*ast.BinaryOperation).Left)
			case "Literal":
				generateExpression(function, expression.(*ast.BinaryOperation).Right)
			}
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.IMULT)
			case "decimalNumber":
				function.bytecode = append(function.bytecode, vm.FMULT)
			}
			if expression.(*ast.BinaryOperation).Right.GetType() == "BinaryOperation" {
				generateExpression(function, expression.(*ast.BinaryOperation).Right)
			}
		case "/":
			switch expression.(*ast.BinaryOperation).Right.GetType() {
			case "BinaryOperation":
				generateExpression(function, expression.(*ast.BinaryOperation).Right.(*ast.BinaryOperation).Left)
			case "Literal":
				generateExpression(function, expression.(*ast.BinaryOperation).Right)
			}
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.IDIV)
			case "DecimalNumber":
				function.bytecode = append(function.bytecode, vm.FDIV)
			}
			if expression.(*ast.BinaryOperation).Right.GetType() == "BinaryOperation" {
				generateExpression(function, expression.(*ast.BinaryOperation).Right)
			}
		case "%":
			switch expression.(*ast.BinaryOperation).Right.GetType() {
			case "BinaryOperation":
				generateExpression(function, expression.(*ast.BinaryOperation).Right.(*ast.BinaryOperation).Left)
			case "Literal":
				generateExpression(function, expression.(*ast.BinaryOperation).Right)
			}
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.IMOD)
			}
			if expression.(*ast.BinaryOperation).Right.GetType() == "BinaryOperation" {
				generateExpression(function, expression.(*ast.BinaryOperation).Right)
			}
		}
	case "Literal":
		expression.(*ast.Literal).Visited = true
		switch expression.(*ast.Literal).Type {
		case "number":
			i, err := GetIntBytes(expression.(*ast.Literal).Value)
			if err != nil {
				return err
			}
			function.Consts[function.ConstCounter] = i
			function.bytecode = append(function.bytecode, vm.ILOAD_CONST)
			function.bytecode = append(function.bytecode, GetUint(function.ConstCounter)...)
			function.ConstCounter++
		case "decimalNumber":
			f, err := GetIntBytes(expression.(*ast.Literal).Value)
			if err != nil {
				return err
			}
			function.Consts[function.ConstCounter] = f
			function.bytecode = append(function.bytecode, vm.FLOAD_CONST)
			function.bytecode = append(function.bytecode, GetUint(function.ConstCounter)...)
			function.ConstCounter++
		case "boolean":
			b := GetBoolBytes(expression.(*ast.Literal).Value)
			function.Consts[function.ConstCounter] = b
			function.bytecode = append(function.bytecode, vm.BOLOAD_CONST)
			function.bytecode = append(function.bytecode, GetUint(function.ConstCounter)...)
			function.ConstCounter++
		}
	}
	return nil
}

func GetIntBytes(str string) ([]byte, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, i)
	return buffer.Bytes(), nil
}

func GetBoolBytes(str string) []byte {
	if str == "true" {
		return []byte{1}
	} else if str == "false" {
		return []byte{0}
	}
	return []byte{}
}

func GetFloatBytes(str string) ([]byte, error) {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, f)
	return buffer.Bytes(), nil
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
func GetUint(i uint64) []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, i)
	return buffer.Bytes()
}
