package main

import (
	"errors"
	"fmt"
)

var UNEXPECTED_ERROR error = errors.New("UNEXPECTED TOKEN")

func Parse(lexer Lexer) (Program, error) {
	var program Program
	var nextToken Token
	_packageId, err := parsePackage(&lexer)
	if err != nil {
		return program, errors.New(formatError(_packageId.Line, _packageId.Name))
	}
	program.Package = _packageId.Value
	nextToken = lexer.next()
	if nextToken.Name == "NEWLINE" {
		jumpNewLines(&lexer)
		nextToken = lexer.next()
		if nextToken.Name == "import" {
			imports, err := parseImports(&lexer, &nextToken)
			program.Imports = imports
			if err != nil {
				return program, errors.New(formatError(nextToken.Line, nextToken.Value))
			}
		} else if err != nil {
			return program, errors.New(formatError(nextToken.Line, nextToken.Value))
		}
	} else {
		return program, errors.New(formatError(nextToken.Line, nextToken.Value))
	}
	if nextToken.Name == "function" {

	} else {
		return program, errors.New(formatError(nextToken.Line, nextToken.Value))
	}
	return program, nil
}

func parseFunction(lxer *Lexer, nextToken *Token) ([]Function, error) {

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
			if *nextToken = lexer.next(); nextToken.Name != "NEWLINE" {
				return _import, UNEXPECTED_ERROR
			} else {
				jumpNewLines(lexer)
				*nextToken = lexer.next()
			}
		} else {
			return _import, UNEXPECTED_ERROR
		}
	} else {
		return _import, UNEXPECTED_ERROR
	}
	return _import, nil
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

func formatError(line uint, close string) string {
	errStr := fmt.Sprintf("Unexpected token close to '%s' at line %d\n", close, line)
	return errStr
}

func jumpNewLines(lexer *Lexer) {
	nextToken := lexer.next()
	for ; nextToken.Name == "NEWLINE"; nextToken = lexer.next() {
		continue
	}
	lexer.CurrentToken--
}
