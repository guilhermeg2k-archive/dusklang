package main

import (
	"errors"
	"fmt"
	"strings"
)

var UNEXPECTED_ERROR error = errors.New("UNEXPECTED TOKEN")

func Parse(lexer Lexer) (Program, error) {
	var program Program
	var nextToken Token
	var imports []Import
	//var functions []Function
	_packageId, err := parsePackage(&lexer)
	if err != nil {
		return program, errors.New(formatError(_packageId.Line, _packageId.Name))
	}
	program.Package = _packageId.Value
	nextToken = lexer.next()
	if nextToken.Name == "import" {
		imports, err = parseImports(&lexer, &nextToken)
		program.Imports = imports
		if err != nil {
			return program, errors.New(formatError(nextToken.Line, nextToken.Value))
		}
	} else if err != nil {
		return program, errors.New(formatError(nextToken.Line, nextToken.Value))
	}
	program.Functions, err = parseFunctions(&lexer, &nextToken)
	if err != nil {
		return program, errors.New(formatError(nextToken.Line, nextToken.Value))
	}
	return program, nil
}

func parseFunctions(lexer *Lexer, nextToken *Token) ([]Function, error) {
	var functions []Function
	var function Function
	for {
		*nextToken = lexer.next()
		println(nextToken.Value)
		if nextToken.Name == "EOF" {
			return functions, nil
		}
		if nextToken.Name == "function" {
			*nextToken = lexer.next()
			if isType(nextToken.Name) {
				function.ReturnType = nextToken.Name
				*nextToken = lexer.next()
			}
			if nextToken.Name == "identifier" {
				function.Identifier = nextToken.Value
				args, err := parseFunctionArgs(lexer, nextToken)
				if err != nil {
					return functions, err
				}
				function.Args = args
				statements, err := parseStatements(lexer, nextToken)
				if err != nil {
					return functions, UNEXPECTED_ERROR
				}
				function.Statements = statements
				functions = append(functions, function)
			} else {
				return functions, UNEXPECTED_ERROR
			}
		} else {
			return functions, UNEXPECTED_ERROR
		}
	}
}

func parseFunctionArgs(lexer *Lexer, nextToken *Token) ([]Arg, error) {
	var args []Arg
	*nextToken = lexer.next()
	if nextToken.Name == "(" {
		*nextToken = lexer.next()
		if nextToken.Name == ")" {
			return args, nil
		} else {
			lexer.back()
		}
		for {
			*nextToken = lexer.next()
			if isType(nextToken.Name) {
				_type := nextToken.Value
				var arg Arg
				arg.Type = _type
				*nextToken = lexer.next()
				if nextToken.Name == "identifier" {
					arg.Identifier = nextToken.Value
					args = append(args, arg)
					*nextToken = lexer.next()
					if nextToken.Name == "," {
						continue
					} else {
						if nextToken.Name == ")" {
							return args, nil
						} else {
							return args, UNEXPECTED_ERROR
						}
					}
				} else {
					return args, UNEXPECTED_ERROR
				}
			} else {
				return args, UNEXPECTED_ERROR
			}
		}
	} else {
		return args, UNEXPECTED_ERROR
	}
}

func parseStatements(lexer *Lexer, nextToken *Token) ([]Statement, error) {
	var statements []Statement
	*nextToken = lexer.next()
	if nextToken.Name == ":" {
		*nextToken = lexer.next()
		if nextToken.Name == "INDENT" {
			for {
				*nextToken = lexer.next()
				if nextToken.Name == "EOF" {
					return statements, nil
				}
				switch nextToken.Name {
				case "var":
					statement, err := parseVarDeclaration(lexer, nextToken)
					if err != nil {
						return statements, UNEXPECTED_ERROR
					}
					statements = append(statements, statement)
				case "identifier":
					identifier := nextToken.Value
					*nextToken = lexer.next()
					switch nextToken.Name {
					//VAR ASSIGN:
					case "=":
						assign, err := parseAssign(lexer, nextToken, identifier)
						if err != nil {
							return statements, err
						}
						statement := Statement{
							Statement: assign,
							Type:      "assign",
						}
						statements = append(statements, statement)

					//FUNCCALL
					case "(":
						funcCall, err := parseFuncCall(lexer, nextToken, identifier)
						if err != nil {
							return statements, err
						}
						statement := Statement{
							Statement: funcCall,
							Type:      "funcCall",
						}
						statements = append(statements, statement)
					default:
						return statements, UNEXPECTED_ERROR
					}
				case "if":
					var ifBlock IfBlock
					condition, err := parseExpression(lexer, nextToken)
					if err != nil {
						return statements, err
					}
					ifBlock.Condition = condition
					ifStatements, err := parseStatements(lexer, nextToken)
					if err != nil {
						return statements, err
					}
					ifBlock.Statements = ifStatements
					*nextToken = lexer.next()
					if nextToken.Name == "else" {
						elseBlock, err := parseElse(lexer, nextToken)
						if err != nil {
							return statements, err
						}
						ifBlock.Else = elseBlock
						statements = append(statements, Statement{Type: "IFBLOCK", Statement: ifBlock})
					} else {
						statements = append(statements, Statement{Type: "IFBLOCK", Statement: ifBlock})
					}
				/*case "switch":
				case "try":*/
				case "for":
					var forBlock ForBlock
					condition, err := parseExpression(lexer, nextToken)
					if err != nil {
						return statements, err
					}
					forBlock.Condition = condition

					forStatements, err := parseStatements(lexer, nextToken)
					if err != nil {
						return statements, err
					}
					forBlock.Statements = forStatements
					statement := Statement{
						Type:      "For",
						Statement: forBlock,
					}
					statements = append(statements, statement)

				case "return":
					exp, err := parseExpression(lexer, nextToken)
					if err != nil {
						return statements, UNEXPECTED_ERROR
					}
					var _return Return
					_return.Expression = exp
					statement := Statement{
						Type:      "return",
						Statement: _return,
					}
					statements = append(statements, statement)
				case "DEDENT":
					return statements, nil
				case "EOF":
					return statements, nil
				default:
					return statements, UNEXPECTED_ERROR
				}
			}
		} else {
			return statements, UNEXPECTED_ERROR
		}
	} else {
		return statements, UNEXPECTED_ERROR
	}
}
func parseElse(lexer *Lexer, nextToken *Token) (IfBlock, error) {
	var elseBlock IfBlock
	*nextToken = lexer.next()
	if nextToken.Name == "if" {
		condition, err := parseExpression(lexer, nextToken)
		if err != nil {
			return elseBlock, err
		}
		elseBlock.Condition = condition
		ifStatements, err := parseStatements(lexer, nextToken)
		if err != nil {
			return elseBlock, err
		}
		elseBlock.Statements = ifStatements
		if nextToken.Name == "EOF" {
			return elseBlock, nil
		}
		*nextToken = lexer.next()
		if nextToken.Name == "else" {
			elseElseBlock, err := parseElse(lexer, nextToken)
			if err != nil {
				return elseBlock, err
			}
			elseBlock.Else = elseElseBlock
		} else {
			*nextToken = lexer.back()
		}
		return elseBlock, nil
	} else {
		*nextToken = lexer.back()
		ifStatements, err := parseStatements(lexer, nextToken)
		if err != nil {
			return elseBlock, err
		}
		elseBlock.Statements = ifStatements
		return elseBlock, nil
	}
}
func parseAssign(lexer *Lexer, nextToken *Token, identifier string) (Assign, error) {
	var assign Assign
	assign.Identifier = identifier
	exp, err := parseExpression(lexer, nextToken)
	if err != nil {
		return assign, err
	}
	assign.Expression = exp
	return assign, nil
}
func parseVarDeclaration(lexer *Lexer, nextToken *Token) (Statement, error) {
	*nextToken = lexer.next()
	var statement Statement
	if isType(nextToken.Name) {
		statement.Type = "FullVarDeclaration"
		_type := nextToken.Name
		var vars FullVarDeclaration
		for {
			*nextToken = lexer.next()
			if nextToken.Name == "identifier" {
				var variable AutoVarDeclaration
				variable.Type = _type
				variable.Identifier = nextToken.Value
				*nextToken = lexer.next()
				if nextToken.Name == "=" {
					expression, err := parseExpression(lexer, nextToken)
					if err != nil {
						return statement, UNEXPECTED_ERROR
					}
					variable.Expression = expression
					*nextToken = lexer.next()
				}
				vars.Variables = append(vars.Variables, variable)
				if nextToken.Name == "," {
					continue
				} else {
					*nextToken = lexer.back()
					break
				}
			} else {
				return statement, UNEXPECTED_ERROR
			}
		}
		statement.Statement = vars
		return statement, nil
	} else if nextToken.Name == "identifier" {
		var variable AutoVarDeclaration
		variable.Identifier = nextToken.Value
		*nextToken = lexer.next()
		if nextToken.Name == "=" {
			//*nextToken = lexer.next()
			//TODO: DECIDE TYPE OF VAR
			exp, err := parseExpression(lexer, nextToken)
			if err != nil {
				return statement, UNEXPECTED_ERROR
			}
			variable.Expression = exp
		} else {
			return statement, UNEXPECTED_ERROR
		}
		statement.Type = "AutoVarDeclaration"
		statement.Statement = variable
		return statement, nil
	}
	return statement, UNEXPECTED_ERROR
}
func parseArgs(lexer *Lexer, nextToken *Token) ([]Expression, error) {
	var expressions []Expression
	for {
		exp, err := parseExpression(lexer, nextToken)
		if err != nil {
			return expressions, err
		}
		expressions = append(expressions, exp)
		*nextToken = lexer.next()
		if nextToken.Name == "," {
			*nextToken = lexer.next()
			continue
		} else {
			*nextToken = lexer.back()
			break
		}
	}
	return expressions, nil
}
func parseFuncCall(lexer *Lexer, nextToken *Token, identifier string) (FuncCall, error) {
	var funcCall FuncCall
	*nextToken = lexer.next()
	if nextToken.Name == ")" {
		var funcCall FuncCall
		funcCall.Identifier = identifier
		return funcCall, nil
	}
	*nextToken = lexer.back()
	funcCall.Identifier = identifier
	expressions, err := parseArgs(lexer, nextToken)
	if err != nil {
		return funcCall, err
	}
	funcCall.Expressions = expressions
	*nextToken = lexer.next()
	if nextToken.Name != ")" {
		return funcCall, UNEXPECTED_ERROR
	}
	return funcCall, nil
}

//TODO: Function to big, modularize it
func parseExpression(lexer *Lexer, nextToken *Token) (interface{}, error) {
	*nextToken = lexer.next()
	if nextToken.Name == "number" ||
		nextToken.Name == "string" ||
		nextToken.Name == "decimalNumber" ||
		nextToken.Name == "true" ||
		nextToken.Name == "false" {
		//BINARY OPERATION
		var literal Literal
		literal.Value = nextToken.Value
		literal.Type = nextToken.Name
		left := literal
		binaryOp, err := parseBinaryOperation(left, lexer, nextToken)
		if err != nil {
			if err.Error() == "NOT_BINARY_OPERATION" {
				//LITERAL
				return literal, nil
			}
			return nil, err
		} else {
			return binaryOp, nil
		}
	}
	if nextToken.Name == "identifier" {
		identifier := nextToken.Value
		*nextToken = lexer.next()
		//FUNC CALL
		if nextToken.Name == "(" {
			funcCall, err := parseFuncCall(lexer, nextToken, identifier)
			if err != nil {
				return funcCall, err
			}
			left := funcCall
			binaryOp, err := parseBinaryOperation(left, lexer, nextToken)
			if err != nil {
				if err.Error() == "NOT_BINARY_OPERATION" {
					return funcCall, nil
				}
				return nil, err
			} else {
				return binaryOp, nil
			}
		} else { //VAR
			*nextToken = lexer.back()
			var literal Literal
			literal.Value = nextToken.Value
			literal.Type = nextToken.Name
			left := literal
			binaryOp, err := parseBinaryOperation(left, lexer, nextToken)
			if err != nil {
				if err.Error() == "NOT_BINARY_OPERATION" {
					return literal, nil
				}
				return nil, err
			} else {
				return binaryOp, nil
			}
		}
	}
	if nextToken.Name == "!" {
		var unaryOperation UnaryOperation
		exp, err := parseExpression(lexer, nextToken)
		if err != nil {
			return nil, UNEXPECTED_ERROR
		}
		unaryOperation.Expression = exp
		unaryOperation.Operator = nextToken.Name
		return unaryOperation, nil
	}
	if nextToken.Name == "(" {
		var parenExpression ParenExpression
		exp, err := parseExpression(lexer, nextToken)
		if err != nil {
			return nil, UNEXPECTED_ERROR
		}
		parenExpression.Expression = exp
		*nextToken = lexer.next()
		if nextToken.Name != ")" {
			return nil, UNEXPECTED_ERROR
		}
		left := parenExpression
		binaryOp, err := parseBinaryOperation(left, lexer, nextToken)
		if err != nil {
			if err.Error() == "NOT_BINARY_OPERATION" {
				return parenExpression, nil
			}
			return nil, err
		} else {
			return binaryOp, nil
		}
	}
	return nil, UNEXPECTED_ERROR
}

func parseBinaryOperation(left Expression, lexer *Lexer, nextToken *Token) (BinaryOperation, error) {
	*nextToken = lexer.next()

	var binaryOperation BinaryOperation
	if nextToken.Name == "+" ||
		nextToken.Name == "-" ||
		nextToken.Name == "*" ||
		nextToken.Name == "/" ||
		nextToken.Name == "%" ||
		nextToken.Name == ">" ||
		nextToken.Name == "<" ||
		nextToken.Name == "==" ||
		nextToken.Name == ">=" ||
		nextToken.Name == "<=" ||
		nextToken.Name == "!=" ||
		nextToken.Name == "and" ||
		nextToken.Name == "or" {
		binaryOperation.Left = left
		binaryOperation.Operator = nextToken.Name
		exp, err := parseExpression(lexer, nextToken)
		if err != nil {
			return binaryOperation, UNEXPECTED_ERROR
		}
		binaryOperation.Right = exp
		return binaryOperation, nil
	} else {
		*nextToken = lexer.back()
		return binaryOperation, errors.New("NOT_BINARY_OPERATION")
	}
}
func parseImports(lexer *Lexer, nextToken *Token) ([]Import, error) {
	var imports []Import
	*nextToken = lexer.next()
	if nextToken.Name == "identifier" {
		var _import Import
		_import, err := parseImport(lexer, nextToken)
		if err != nil {
			return imports, err
		}
		imports = append(imports, _import)
		//Parse multiple imports
		for nextToken.Name == "import" {
			*nextToken = lexer.next()
			if nextToken.Name == "identifier" {
				var _import Import
				_import, err := parseImport(lexer, nextToken)
				if err != nil {
					return imports, err
				}
				imports = append(imports, _import)
			} else {
				return imports, UNEXPECTED_ERROR
			}
		}
	} else {
		return imports, UNEXPECTED_ERROR
	}
	return imports, nil
}

func parseImport(lexer *Lexer, nextToken *Token) (Import, error) {
	var _import Import
	var identifiers []string
	identifiers = append(identifiers, nextToken.Value)
	*nextToken = lexer.next()
	//More than one import from the same source
	for ; nextToken.Name == ","; *nextToken = lexer.next() {
		if *nextToken = lexer.next(); nextToken.Name == "identifier" {
			identifiers = append(identifiers, nextToken.Value)
		} else {
			return _import, UNEXPECTED_ERROR
		}
	}
	if nextToken.Name == "from" {
		_import.Identifiers = identifiers
		if *nextToken = lexer.next(); nextToken.Name == "identifier" {
			_import.from = nextToken.Value
			return _import, nil
		} else {
			return _import, UNEXPECTED_ERROR
		}
	} else {
		return _import, UNEXPECTED_ERROR
	}
}

func parsePackage(lexer *Lexer) (Token, error) {
	_package := lexer.next()
	value := lexer.next()
	if _package.Name == "package" {
		if value.Name != "identifier" {
			return value, UNEXPECTED_ERROR
		}
		return value, nil
	}
	return value, UNEXPECTED_ERROR
}

func isType(s string) bool {
	var types string = "byte int16 int32 int64 uint16 uint32 uint64 float double string"
	return strings.Contains(types, s)
}

func formatError(line uint, close string) string {
	errStr := fmt.Sprintf("Unexpected token close to '%s' at line %d\n", close, line)
	return errStr
}
