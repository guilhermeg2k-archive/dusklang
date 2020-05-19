package parser

import (
	"errors"
	"fmt"

	"github.com/guilhermeg2k/dusklang/ast"
)

func Analyze(p *ast.Program) error {
	err := analyzeFunctions(p.Functions)
	if err != nil {
		return err
	}
	return nil
}

//TODO: Make it all more procedural
func analyzeFunctions(functions []ast.Function) error {
	var funcList []string
	//Check function names
	for _, f := range functions {
		for _, alreadyAddedFunc := range funcList {
			if f.Identifier == alreadyAddedFunc {
				//TODO: Make a Better error msg pattern
				return errors.New(fmt.Sprintf("Function '%s' in line %d already declared", f.Identifier, f.Line))
			}
			funcList = append(funcList, f.Identifier)
		}
	}

	for _, f := range functions {
		var functionVars []ast.Variable
		functionVars = append(functionVars, f.Args...)
		//Check function statements
		for _, stament := range f.Statements {
			//TODO: Change the name of types to match with ast.go
			switch stament.Type {
			case "FullVarDeclaration":
				for _, variable := range stament.Statement.(ast.FullVarDeclaration).Variables {
					if !isVarDeclared(variable.Identifier, functionVars) {
						functionVars = append(functionVars,
							ast.Variable{
								Identifier: variable.Identifier,
								Type:       variable.Type,
							})
						if variable.Expression != nil {
							expType, err := getExpressionType(functions, variable.Expression, functionVars)
							if err != nil {
								return errors.New(formatError2(stament.Line, err.Error()))
							}
							if !(expType == variable.Type) {
								return errors.New(fmt.Sprintf("Invalid type of value on variable '%s' assignment on line %d", variable.Identifier, stament.Line))
							}
						}
					} else {
						return errors.New(fmt.Sprintf("ast.Variable '%s' on line %d is already declared ", variable.Identifier, stament.Line))
					}
				}
			//Defines the type of auto vars
			case "AutoVarDeclaration":
				expType, err := getExpressionType(functions, stament.Statement.(ast.AutoVarDeclaration).Expression, functionVars)
				if err != nil {

					return errors.New(formatError2(stament.Line, err.Error()))
				}
				autoVarDeclaration := stament.Statement.(ast.AutoVarDeclaration)
				autoVarDeclaration.Type = expType
			case "Assign":
				expType, err := getExpressionType(functions, stament.Statement.(ast.Assign).Expression, functionVars)
				if err != nil {
					return errors.New(formatError2(stament.Line, err.Error()))
				}
				assign := stament.Statement.(ast.Assign)
				assign.Type = expType
			case "funcCall":
				funcCAll := stament.Statement.(ast.FuncCall)
				if !funcExists(functions, funcCAll.Identifier) {
					return errors.New(fmt.Sprintf("Function '%s' not declared, line %d", funcCAll.Identifier, stament.Line))
				}
			case "IFBLOCK":
				err := analyzeIfBlock(functions, stament.Statement.(ast.IfBlock), functionVars)
				if err != nil {
					return errors.New(formatError2(stament.Line, err.Error()))
				}
			case "For":
				//Todo: add support for all for style, it's just supporting the basic for-while style
				expType, err := getExpressionType(functions, stament.Statement.(ast.ForBlock).Condition, functionVars)
				if err != nil {
					return errors.New(formatError2(stament.Line, err.Error()))
				}
				if !(expType == "bool") {
					return errors.New(fmt.Sprintf("Invalid condition expression, line %d", stament.Line))
				}
			case "return":
				expType, err := getExpressionType(functions, stament.Statement.(ast.Return).Expression, functionVars)
				if err != nil {
					return errors.New(formatError2(stament.Line, err.Error()))
				}
				if !(expType == f.ReturnType) {
					return errors.New(fmt.Sprintf("Invalid return type, line %d", stament.Line))
				}
			}
		}
	}
	return nil
}

func analyzeIfBlock(functions []ast.Function, ifBlock ast.IfBlock, functionVars []ast.Variable) error {
	if ifBlock.Condition == nil {
		return nil
	}
	expType, err := getExpressionType(functions, ifBlock.Condition, functionVars)
	if err != nil {
		return err
	}
	if !(expType == "bool") {
		return errors.New("Invalid condition return type")
	}
	if ifBlock.Else != nil {
		fmt.Println(ifBlock.Else)
	}
	if ifBlock.Else != nil {
		err = analyzeIfBlock(functions, ifBlock.Else.(ast.IfBlock), functionVars)
		if err != nil {
			return err
		}
	}
	return nil
}

func getExpressionType(functions []ast.Function, expression ast.Expression, functionVars []ast.Variable) (string, error) {
	switch expression.GetType() {
	case "Literal":
		switch expression.(*ast.Literal).Type {
		case "number":
			return "int", nil
		case "decimalNumber":
			return "float", nil
		case "string":
			return "string", nil
		case "boolean":
			return "bool", nil
		}
	case "UnaryOperation":
		//TODO: It's just working for "!" operator, make it work for all unary operators
		return "bool", nil
	case "BinaryOperation":
		typeOfLeft, err := getExpressionType(functions, expression.(*ast.BinaryOperation).Left, functionVars)
		if err != nil {
			return "", err
		}
		typeOfRight, err := getExpressionType(functions, expression.(*ast.BinaryOperation).Right, functionVars)
		if err != nil {
			return "", err
		}
		fmt.Println(typeOfLeft)
		if typeOfLeft == typeOfRight {
			if isComparison(expression.(*ast.BinaryOperation).Operator) {
				return "bool", nil
			}
			return typeOfLeft, nil
		}
		return "", errors.New("Invalid type operation")
	case "ParenExpression":
		parenType, err := getExpressionType(functions, expression.(*ast.ParenExpression).Expression, functionVars)
		if err != nil {
			return "", err
		}
		return parenType, nil
	case "FuncCall":
		funcReturnType, err := getFuncReturnType(functions, expression.(*ast.FuncCall).Identifier)
		if err != nil {
			return "", err
		}
		return funcReturnType, nil
	case "Variable":
		varType, err := getVariableType(functionVars, expression.(*ast.Variable).Identifier)
		if err != nil {
			return "", err
		}
		return varType, nil
	}
	return "", errors.New("Invalid expression type")
}
func getVariableType(vars []ast.Variable, identifier string) (string, error) {
	for _, v := range vars {
		if v.Identifier == identifier {
			return v.Type, nil
		}
	}
	return "", errors.New("ast.Variable not declared")
}

func funcExists(functions []ast.Function, functionName string) bool {
	for _, f := range functions {
		if f.Identifier == functionName {
			return true
		}
	}
	return false
}
func getFuncReturnType(functions []ast.Function, functionName string) (string, error) {
	for _, f := range functions {
		if f.Identifier == functionName {
			return f.ReturnType, nil
		}
	}
	return "", errors.New("Undeclared function call")
}

func isVarDeclared(variable string, varsList []ast.Variable) bool {
	for _, v := range varsList {
		if v.Identifier == variable {
			return true
		}
	}
	return false
}

func formatError2(line uint, msg string) string {
	errStr := fmt.Sprintf("%s on line %d", msg, line)
	return errStr
}
