package main

const (
	Byte   int = 1
	Int16  int = 2
	Int32  int = 3
	Int64  int = 4
	Uint16 int = 5
	Uint32 int = 6
	Uint64 int = 7
	Float  int = 8
	Double int = 9
	String int = 10
)

type Type int
type Statement map[string]interface{}
type Expression interface{}

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
	ReturnType Type
	Args       []Arg
	Statements Statement
}

type Arg struct {
	Type       Type
	Identifier string
}

type AutoVarDeclaration struct {
	Identifier string
	Type       Type
	Expression Expression
}

type FullVarDeclaration struct {
	Type        Type
	Identifiers []string
}

type Assign struct {
	Identifier string
	Expression Expression
}

type If struct {
	Condition  Expression
	Statements []Statement
	Else       Else
}

type Else struct {
	Statements []Statement
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

type For struct {
	Condition  Expression
	Statements []Statement
}

type ForIn struct {
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
	Type Type
	Expression
}

type BinaryExpression struct {
	Operator   string
	ReturnType Type
	Left       Expression
	Right      Expression
}

type UnaryExpression struct {
	Operator   string
	Expression Expression
}

type ParenExpression struct {
	BinaryExpression BinaryExpression
}

type Literal struct {
	Type  Type
	Value interface{}
}
