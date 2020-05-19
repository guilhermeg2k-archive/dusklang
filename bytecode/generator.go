package bytecode

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/guilhermeg2k/dusklang/ast"
	"github.com/guilhermeg2k/dusklang/vm"
)

type VariablesOffset map[string]uint64
type LabelOffset map[uint64]uint64
type Function struct {
	Consts          vm.Consts
	Labels          vm.Labels
	Variables       []ast.Variable
	VariablesOffset VariablesOffset
	StorageCounter  uint64
	ConstCounter    uint64
	LabelCounter    uint64
	LabelOffset     LabelOffset
	bytecode        []byte
}

func GenerateByteCode(program *ast.Program) vm.Function {
	main := program.Functions[0]
	function := generateFunctionByteCode(&main)
	f := vm.Function{
		Labels:   function.Labels,
		Consts:   function.Consts,
		Bytecode: function.bytecode,
		Storage:  &vm.Storage{},
	}
	return f
}

func generateFunctionByteCode(function *ast.Function) Function {
	f := Function{
		Consts:          make(vm.Consts),
		Labels:          make(vm.Labels),
		VariablesOffset: make(VariablesOffset),
		LabelOffset:     make(LabelOffset),
		bytecode:        []byte{},
	}
	for _, statement := range function.Statements {
		generateStatement(&f, statement)
	}
	f.bytecode = append(f.bytecode, 255)
	return f
}

func generateStatement(function *Function, statement ast.Statement) {
	switch statement.Type {
	case "FullVarDeclaration":
		generateFullVarDeclaration(function, statement.Statement.(ast.FullVarDeclaration))
	case "AutoVarDeclaration":
		generateAutoVarDeclaration(function, statement.Statement.(ast.AutoVarDeclaration))
	case "IFBLOCK":
		generateIf(function, statement.Statement.(ast.IfBlock))
	case "For":
		generateFor(function, statement.Statement.(ast.ForBlock))
	}
}

func generateFor(function *Function, forBlock ast.ForBlock) {
	generateExpression(function, forBlock.Condition)
	function.bytecode = append(function.bytecode, vm.JUMP_IF_ELSE)
	function.bytecode = append(function.bytecode, GetUint(function.LabelCounter)...)

	forLabel := len(function.bytecode)
	for _, statement := range forBlock.Statements {
		generateStatement(function, statement)
	}

	generateExpression(function, forBlock.Condition)
	function.bytecode = append(function.bytecode, vm.JUMP_IF_TRUE)
	function.bytecode = append(function.bytecode, GetUint(uint64(forLabel))...)

	function.LabelOffset[function.LabelCounter] = uint64(len(function.bytecode))
	function.LabelCounter++

	function.LabelOffset[function.LabelCounter] = uint64(forLabel)
	function.LabelCounter++
}

func generateIf(function *Function, ifBlock ast.IfBlock) {
	if ifBlock.Condition != nil {
		generateExpression(function, ifBlock.Condition)
		function.bytecode = append(function.bytecode, vm.JUMP_IF_ELSE)
		function.bytecode = append(function.bytecode, GetUint(function.LabelCounter)...)
		for _, statement := range ifBlock.Statements {
			generateStatement(function, statement)
		}
		function.LabelOffset[function.LabelCounter] = uint64(len(function.bytecode))
		function.LabelCounter++
		if ifBlock.Else != nil {
			generateIf(function, ifBlock.Else.(ast.IfBlock))
		}
	} else {
		for _, statement := range ifBlock.Statements {
			generateStatement(function, statement)
		}
	}
}

func generateAssign(function *Function, assign ast.Assign) {
	generateExpression(function, assign.Expression)
	switch assign.Type {
	case "int":
		storeInt(function, function.VariablesOffset[assign.Identifier])
	case "float":
		storeFloat(function, function.VariablesOffset[assign.Identifier])
	case "bool":
		storeBool(function, function.VariablesOffset[assign.Identifier])
	}
}

func generateAutoVarDeclaration(function *Function, variable ast.AutoVarDeclaration) {
	generateExpression(function, variable.Expression)
	switch variable.Type {
	case "int":
		storeInt(function, function.ConstCounter)
	case "float":
		storeFloat(function, function.ConstCounter)
	case "bool":
		storeBool(function, function.ConstCounter)
	}
	function.VariablesOffset[variable.Identifier] = function.StorageCounter
	function.StorageCounter++
}

func generateFullVarDeclaration(function *Function, fullVarDeclaration ast.FullVarDeclaration) {
	for _, variable := range fullVarDeclaration.Variables {
		if variable.Expression != nil {
			generateExpression(function, variable.Expression)
		} else {
			switch variable.Type {
			case "int":
				initiateInt(function)
			case "float":
				iniateFloat(function)
			case "bool":
				iniateBool(function)
			}
			function.ConstCounter++
		}
		switch variable.Type {
		case "int":
			storeInt(function, function.ConstCounter)
		case "float":
			storeFloat(function, function.ConstCounter)
		case "bool":
			storeBool(function, function.ConstCounter)
		}
		function.VariablesOffset[variable.Identifier] = function.StorageCounter
		function.StorageCounter++
	}
}

func generateExpression(function *Function, expression ast.Expression) error {
	//TODO: 'AND' AND 'OR' OPERATIONS
	switch expression.GetType() {
	case "ParenExpression":
		generateExpression(function, expression.(*ast.ParenExpression).Expression)
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
			case "ParenExpression":
				generateExpression(function, expression.(*ast.BinaryOperation).Right)
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
			case "ParenExpression":
				generateExpression(function, expression.(*ast.BinaryOperation).Right)
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
			case "ParenExpression":
				generateExpression(function, expression.(*ast.BinaryOperation).Right)
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
		case "==":
			generateExpression(function, expression.(*ast.BinaryOperation).Left)
			generateExpression(function, expression.(*ast.BinaryOperation).Right)
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.ICMP_EQUALS)
			case "decimalNumber":
				function.bytecode = append(function.bytecode, vm.FCMP_EQUALS)
			}
		case "<=":
			generateExpression(function, expression.(*ast.BinaryOperation).Left)
			generateExpression(function, expression.(*ast.BinaryOperation).Right)
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.ICMP_LESS_EQUALS)
			case "decimalNumber":
				function.bytecode = append(function.bytecode, vm.FCMP_LESS_EQUALS)
			}
		case ">=":
			generateExpression(function, expression.(*ast.BinaryOperation).Left)
			generateExpression(function, expression.(*ast.BinaryOperation).Right)
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.ICMP_GREATER_EQUALS)
			case "decimalNumber":
				function.bytecode = append(function.bytecode, vm.FCMP_GREATER_EQUALS)
			}
		case "<":
			generateExpression(function, expression.(*ast.BinaryOperation).Left)
			generateExpression(function, expression.(*ast.BinaryOperation).Right)
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.ICMP_LESS_THEN)
			case "decimalNumber":
				function.bytecode = append(function.bytecode, vm.FCMP_LESS_THEN)
			}
		case ">":
			generateExpression(function, expression.(*ast.BinaryOperation).Left)
			generateExpression(function, expression.(*ast.BinaryOperation).Right)
			switch expression.(*ast.BinaryOperation).Left.(*ast.Literal).Type {
			case "number":
				function.bytecode = append(function.bytecode, vm.ICMP_GREATER_THEN)
			case "decimalNumber":
				function.bytecode = append(function.bytecode, vm.FCMP_GREATER_THEN)
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
	case "Variable":
		switch expression.(*ast.Variable).Type {
		case "int":
			function.bytecode = append(function.bytecode, vm.ILOAD)
			function.bytecode = append(function.bytecode, GetUint(function.VariablesOffset[expression.(*ast.Variable).Identifier])...)
		case "float":
			function.bytecode = append(function.bytecode, vm.FLOAD)
			function.bytecode = append(function.bytecode, GetUint(function.VariablesOffset[expression.(*ast.Variable).Identifier])...)
		case "boolean":
			function.bytecode = append(function.bytecode, vm.BOLOAD)
			function.bytecode = append(function.bytecode, GetUint(function.VariablesOffset[expression.(*ast.Variable).Identifier])...)
		}
	}
	return nil
}

func initiateInt(function *Function) {
	function.Consts[function.ConstCounter] = GetInt(0)
	function.bytecode = append(function.bytecode, vm.ILOAD_CONST)
	function.bytecode = append(function.bytecode, GetUint(function.ConstCounter)...)
}

func iniateFloat(function *Function) {
	function.Consts[function.ConstCounter] = GetFloat(0)
	function.bytecode = append(function.bytecode, vm.FLOAD_CONST)
	function.bytecode = append(function.bytecode, GetUint(function.ConstCounter)...)
}

func iniateBool(function *Function) {
	function.Consts[function.ConstCounter] = []byte{0}
	function.bytecode = append(function.bytecode, vm.BOLOAD_CONST)
	function.bytecode = append(function.bytecode, GetUint(function.ConstCounter)...)
}
func storeBool(function *Function, pos uint64) {
	function.bytecode = append(function.bytecode, vm.BOSTORE)
	function.bytecode = append(function.bytecode, GetUint(pos)...)
}

func storeFloat(function *Function, pos uint64) {
	function.bytecode = append(function.bytecode, vm.FSTORE)
	function.bytecode = append(function.bytecode, GetUint(pos)...)
}
func storeInt(function *Function, pos uint64) {
	function.bytecode = append(function.bytecode, vm.ISTORE)
	function.bytecode = append(function.bytecode, GetUint(pos)...)
}

//TODO: Modularize those functions, are same of them on other packages
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
