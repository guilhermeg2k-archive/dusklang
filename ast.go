package main

type Statement struct {
	Type      string
	Statement interface{}
}

type Expression interface {
}

type Program struct {
	Package   string
	Imports   []Import
	Functions []Function
}

type Import struct {
	Identifiers []string
	from        string
}

type Function struct {
	Identifier string
	ReturnType string
	Args       []Arg
	Statements []Statement
}

type Arg struct {
	Type       string
	Identifier string
}

type AutoVarDeclaration struct {
	Identifier string
	Type       string
	Expression Expression
}

type FullVarDeclaration struct {
	Variables []AutoVarDeclaration
}

type Assign struct {
	Identifier string
	Expression Expression
}

type IfBlock struct {
	Condition  Expression
	Statements []Statement
	Else       interface{}
}

type Case struct {
	Condition  Expression
	Statements []Statement
}
type Switch struct {
	Identifier string
	Cases      []Case
	Default    Case
}

type ForBlock struct {
	Condition  Expression
	Statements []Statement
}

type ForIn struct {
	Inicialize interface{}
	Identifier string
	Range      Expression
	Statement  []Statement
}

type FuncCall struct {
	Identifier  string
	Expressions []Expression
}

type Try struct {
	Statement []Statement
	Catch     []Catch
	Finally   []Statement
}

type Catch struct {
	Identifer string
	Statement []Statement
}

type Return struct {
	Expression
}

type BinaryOperation struct {
	Operator string
	Left     Expression
	Right    Expression
}

type UnaryOperation struct {
	Operator   string
	Expression Expression
}

type Literal struct {
	Type  string
	Value interface{}
}

type ParenExpression struct {
	Expression Expression
}
