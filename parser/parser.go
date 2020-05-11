package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/guilhermeg2k/dusklang/ast"
	"github.com/guilhermeg2k/dusklang/lexer"
)

var UNEXPECTED_ERROR error = errors.New("UNEXPECTED TOKEN")

//TODO: The functions on this file are a mess that's no pattern, fix it

func Parse(l lexer.Lexer) (ast.Program, error) {
	var program ast.Program
	var nextToken lexer.Token
	var imports []ast.Import
	_packageId, err := parsePackage(&l)
	if err != nil {
		return program, errors.New(formatError(_packageId.Line, _packageId.Name))
	}
	program.Package = _packageId.Value
	nextToken = l.Next()
	if nextToken.Name == "import" {
		imports, err = parseImports(&l, &nextToken)
		program.Imports = imports
		if err != nil {
			//TODO: Better error msg pattern
			return program, errors.New(formatError(nextToken.Line, nextToken.Value))
		}
	} else if err != nil {
		return program, errors.New(formatError(nextToken.Line, nextToken.Value))
	}
	program.Functions, err = parseFunctions(&l, &nextToken)
	if err != nil {
		return program, errors.New(formatError(nextToken.Line, nextToken.Value))
	}
	return program, nil
}

func parseFunctions(l *lexer.Lexer, nextToken *lexer.Token) ([]ast.Function, error) {
	var functions []ast.Function
	var function ast.Function
	for {
		*nextToken = l.Next()
		if nextToken.Name == "EOF" {
			return functions, nil
		}
		if nextToken.Name == "function" {
			function.Line = nextToken.Line
			*nextToken = l.Next()
			if isType(nextToken.Name) {
				function.ReturnType = nextToken.Name
				*nextToken = l.Next()
			}
			if nextToken.Name == "identifier" {
				function.Identifier = nextToken.Value
				args, err := parseFunctionArgs(l, nextToken)
				if err != nil {
					return functions, err
				}
				function.Args = args
				statements, err := parseStatements(l, nextToken)
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

func parseFunctionArgs(l *lexer.Lexer, nextToken *lexer.Token) ([]ast.Variable, error) {
	var args []ast.Variable
	*nextToken = l.Next()
	if nextToken.Name == "(" {
		*nextToken = l.Next()
		if nextToken.Name == ")" {
			return args, nil
		} else {
			l.Back()
		}
		for {
			*nextToken = l.Next()
			if isType(nextToken.Name) {
				_type := nextToken.Value
				var arg ast.Variable
				arg.Type = _type
				*nextToken = l.Next()
				if nextToken.Name == "identifier" {
					arg.Identifier = nextToken.Value
					args = append(args, arg)
					*nextToken = l.Next()
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

func parseStatements(l *lexer.Lexer, nextToken *lexer.Token) ([]ast.Statement, error) {
	var statements []ast.Statement
	*nextToken = l.Next()
	if nextToken.Name == ":" {
		*nextToken = l.Next()
		if nextToken.Name == "INDENT" {
			for {
				*nextToken = l.Next()
				if nextToken.Name == "EOF" {
					return statements, nil
				}
				switch nextToken.Name {
				case "var":
					statement, err := parseVarDeclaration(l, nextToken)
					if err != nil {
						return statements, UNEXPECTED_ERROR
					}
					statements = append(statements, statement)
				case "identifier":
					identifier := nextToken.Value
					*nextToken = l.Next()
					switch nextToken.Name {
					//VAR ASSIGN:
					case "=":
						assign, err := parseAssign(l, nextToken, identifier)
						if err != nil {
							return statements, err
						}
						statement := ast.Statement{
							Statement: assign,
							Type:      "assign",
						}
						statement.Line = nextToken.Line
						statements = append(statements, statement)

					//FUNCCALL
					case "(":
						funcCall, err := parseFuncCall(l, nextToken, identifier)
						if err != nil {
							return statements, err
						}
						statement := ast.Statement{
							Statement: funcCall,
							Type:      "funcCall",
						}
						statement.Line = nextToken.Line
						statements = append(statements, statement)
					default:
						return statements, UNEXPECTED_ERROR
					}
				case "if":
					var ifBlock ast.IfBlock
					ifLine := nextToken.Line
					ifBlock.Line = ifLine
					condition, err := parseExpression(l, nextToken)
					if err != nil {
						return statements, err
					}
					ifBlock.Condition = condition
					ifStatements, err := parseStatements(l, nextToken)
					if err != nil {
						return statements, err
					}
					ifBlock.Statements = ifStatements
					*nextToken = l.Next()
					if nextToken.Name == "else" {
						elseBlock, err := parseElse(l, nextToken)
						if err != nil {
							return statements, err
						}
						ifBlock.Else = elseBlock
						statements = append(statements, ast.Statement{Type: "IFBLOCK", Statement: ifBlock, Line: ifLine})
					} else {
						statements = append(statements, ast.Statement{Type: "IFBLOCK", Statement: ifBlock, Line: ifLine})
					}

				/*case "switch":
				case "try":*/
				case "for":
					var forBlock ast.ForBlock
					condition, err := parseExpression(l, nextToken)
					if err != nil {
						return statements, err
					}
					forBlock.Condition = condition
					forLine := nextToken.Line
					forStatements, err := parseStatements(l, nextToken)
					if err != nil {
						return statements, err
					}
					forBlock.Statements = forStatements
					statement := ast.Statement{
						Type:      "For",
						Statement: forBlock,
						Line:      forLine,
					}
					statements = append(statements, statement)

				case "return":
					exp, err := parseExpression(l, nextToken)
					if err != nil {
						return statements, UNEXPECTED_ERROR
					}
					var _return ast.Return
					_return.Expression = exp
					statement := ast.Statement{
						Type:      "return",
						Statement: _return,
						Line:      nextToken.Line,
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
func parseElse(l *lexer.Lexer, nextToken *lexer.Token) (ast.IfBlock, error) {
	var elseBlock ast.IfBlock
	*nextToken = l.Next()
	elseBlock.Line = nextToken.Line
	if nextToken.Name == "if" {
		condition, err := parseExpression(l, nextToken)
		if err != nil {
			return elseBlock, err
		}
		elseBlock.Condition = condition
		ifStatements, err := parseStatements(l, nextToken)
		if err != nil {
			return elseBlock, err
		}
		elseBlock.Statements = ifStatements
		if nextToken.Name == "EOF" {
			return elseBlock, nil
		}
		*nextToken = l.Next()
		if nextToken.Name == "else" {
			elseElseBlock, err := parseElse(l, nextToken)
			if err != nil {
				return elseBlock, err
			}
			elseBlock.Else = elseElseBlock
		} else {
			*nextToken = l.Back()
		}
		return elseBlock, nil
	} else {
		*nextToken = l.Back()
		ifStatements, err := parseStatements(l, nextToken)
		if err != nil {
			return elseBlock, err
		}
		elseBlock.Statements = ifStatements
		return elseBlock, nil
	}
}
func parseAssign(l *lexer.Lexer, nextToken *lexer.Token, identifier string) (ast.Assign, error) {
	var assign ast.Assign
	assign.Identifier = identifier
	exp, err := parseExpression(l, nextToken)
	if err != nil {
		return assign, err
	}
	assign.Expression = exp
	return assign, nil
}
func parseVarDeclaration(l *lexer.Lexer, nextToken *lexer.Token) (ast.Statement, error) {
	*nextToken = l.Next()
	var statement ast.Statement
	if isType(nextToken.Name) {
		statement.Type = "FullVarDeclaration"
		statement.Line = nextToken.Line
		_type := nextToken.Name
		var vars ast.FullVarDeclaration
		for {
			*nextToken = l.Next()
			if nextToken.Name == "identifier" {
				var variable ast.AutoVarDeclaration
				variable.Type = _type
				variable.Identifier = nextToken.Value
				*nextToken = l.Next()
				if nextToken.Name == "=" {
					expression, err := parseExpression(l, nextToken)
					if err != nil {
						return statement, UNEXPECTED_ERROR
					}
					variable.Expression = expression
					*nextToken = l.Next()
				}
				vars.Variables = append(vars.Variables, variable)
				if nextToken.Name == "," {
					continue
				} else {
					*nextToken = l.Back()
					break
				}
			} else {
				return statement, UNEXPECTED_ERROR
			}
		}
		statement.Statement = vars
		return statement, nil
	} else if nextToken.Name == "identifier" {
		var variable ast.AutoVarDeclaration
		variable.Identifier = nextToken.Value
		*nextToken = l.Next()
		if nextToken.Name == "=" {
			//*nextToken = l.Next()
			//TODO: DECIDE TYPE OF VAR
			exp, err := parseExpression(l, nextToken)
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
func parseArgs(l *lexer.Lexer, nextToken *lexer.Token) ([]ast.Expression, error) {
	var expressions []ast.Expression
	for {
		exp, err := parseExpression(l, nextToken)
		if err != nil {
			return expressions, err
		}
		expressions = append(expressions, exp)
		*nextToken = l.Next()
		if nextToken.Name == "," {
			*nextToken = l.Next()
			continue
		} else {
			*nextToken = l.Back()
			break
		}
	}
	return expressions, nil
}
func parseFuncCall(l *lexer.Lexer, nextToken *lexer.Token, identifier string) (ast.FuncCall, error) {
	var funcCall ast.FuncCall
	*nextToken = l.Next()
	if nextToken.Name == ")" {
		var funcCall ast.FuncCall
		funcCall.Identifier = identifier
		return funcCall, nil
	}
	*nextToken = l.Back()
	funcCall.Identifier = identifier
	expressions, err := parseArgs(l, nextToken)
	if err != nil {
		return funcCall, err
	}
	funcCall.Expressions = expressions
	*nextToken = l.Next()
	if nextToken.Name != ")" {
		return funcCall, UNEXPECTED_ERROR
	}
	return funcCall, nil
}

//TODO: Function to big, modularize it
func parseExpression(l *lexer.Lexer, nextToken *lexer.Token) (ast.Expression, error) {
	*nextToken = l.Next()
	if nextToken.Name == "number" ||
		nextToken.Name == "string" ||
		nextToken.Name == "decimalNumber" ||
		nextToken.Name == "true" ||
		nextToken.Name == "false" {
		//BINARY OPERATION
		var literal ast.Literal
		if nextToken.Name == "true" || nextToken.Name == "false" {
			literal.Type = "boolean"
		} else {
			literal.Type = nextToken.Name
		}
		literal.Value = nextToken.Value
		left := literal
		binaryOp, err := parseBinaryOperation(left, l, nextToken)
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
		*nextToken = l.Next()
		//FUNC CALL
		if nextToken.Name == "(" {
			funcCall, err := parseFuncCall(l, nextToken, identifier)
			if err != nil {
				return funcCall, err
			}
			left := funcCall
			binaryOp, err := parseBinaryOperation(left, l, nextToken)
			if err != nil {
				if err.Error() == "NOT_BINARY_OPERATION" {
					return funcCall, nil
				}
				return nil, err
			} else {
				return binaryOp, nil
			}
		} else { //VAR
			*nextToken = l.Back()
			var variable ast.Variable
			variable.Identifier = nextToken.Value
			fmt.Println(variable.Type)
			left := variable
			binaryOp, err := parseBinaryOperation(left, l, nextToken)
			if err != nil {
				if err.Error() == "NOT_BINARY_OPERATION" {
					return variable, nil
				}
				return nil, err
			} else {
				return binaryOp, nil
			}
		}
	}
	if nextToken.Name == "!" {
		var unaryOperation ast.UnaryOperation
		exp, err := parseExpression(l, nextToken)
		if err != nil {
			return nil, UNEXPECTED_ERROR
		}
		unaryOperation.Expression = exp
		unaryOperation.Operator = nextToken.Name
		return unaryOperation, nil
	}
	if nextToken.Name == "(" {
		var parenExpression ast.ParenExpression
		exp, err := parseExpression(l, nextToken)
		if err != nil {
			return nil, UNEXPECTED_ERROR
		}
		parenExpression.Expression = exp
		*nextToken = l.Next()
		if nextToken.Name != ")" {
			return nil, UNEXPECTED_ERROR
		}
		left := parenExpression
		binaryOp, err := parseBinaryOperation(left, l, nextToken)
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

//TODO: Lexer as first parameter
func parseBinaryOperation(left ast.Expression, l *lexer.Lexer, nextToken *lexer.Token) (ast.BinaryOperation, error) {
	*nextToken = l.Next()

	var binaryOperation ast.BinaryOperation
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
		exp, err := parseExpression(l, nextToken)
		if err != nil {
			return binaryOperation, UNEXPECTED_ERROR
		}
		binaryOperation.Right = exp
		return binaryOperation, nil
	} else {
		*nextToken = l.Back()
		return binaryOperation, errors.New("NOT_BINARY_OPERATION")
	}
}

func parseImports(l *lexer.Lexer, nextToken *lexer.Token) ([]ast.Import, error) {
	//TODO: Get rid of the beggining of the function and put everything inside of the forloop
	var imports []ast.Import
	*nextToken = l.Next()
	if nextToken.Name == "identifier" {
		var _import ast.Import
		_import, err := parseImport(l, nextToken)
		if err != nil {
			return imports, err
		}
		imports = append(imports, _import)
		//Parse multiple imports
		for nextToken.Name == "import" {
			*nextToken = l.Next()
			if nextToken.Name == "identifier" {
				var _import ast.Import
				_import, err := parseImport(l, nextToken)
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

func parseImport(l *lexer.Lexer, nextToken *lexer.Token) (ast.Import, error) {
	var _import ast.Import
	var identifiers []string
	identifiers = append(identifiers, nextToken.Value)
	*nextToken = l.Next()
	//More than one import from the same source
	for ; nextToken.Name == ","; *nextToken = l.Next() {
		if *nextToken = l.Next(); nextToken.Name == "identifier" {
			identifiers = append(identifiers, nextToken.Value)
		} else {
			return _import, UNEXPECTED_ERROR
		}
	}
	if nextToken.Name == "from" {
		_import.Identifiers = identifiers
		if *nextToken = l.Next(); nextToken.Name == "identifier" {
			_import.From = nextToken.Value
			_import.Line = nextToken.Line
			return _import, nil
		} else {
			return _import, UNEXPECTED_ERROR
		}
	} else {
		return _import, UNEXPECTED_ERROR
	}
}

func parsePackage(l *lexer.Lexer) (lexer.Token, error) {
	_package := l.Next()
	value := l.Next()
	if _package.Name == "package" {
		if value.Name != "identifier" {
			return value, UNEXPECTED_ERROR
		}
		_package.Line = value.Line
		return value, nil
	}
	return value, UNEXPECTED_ERROR
}

func isType(s string) bool {
	var types string = "byte int float string bool"
	return strings.Contains(types, s)
}

func formatError(line uint, close string) string {
	errStr := fmt.Sprintf("Unexpected token close to '%s' at line %d\n", close, line)
	return errStr
}
