package ast

//TODO: Remove Type and make a interface
type Statement struct {
	Type      string
	Statement interface{}
	Line      uint
}

type Expression interface {
	GetType() string
}

type Program struct {
	Package   string
	Imports   []Import
	Functions []Function
}

type Import struct {
	Identifiers []string
	From        string
	Line        uint
}

type Variables map[string]Variable
type Function struct {
	Identifier string
	ReturnType string
	Args       []Variable
	Variables  Variables
	Statements []Statement
	Line       uint
}

type Variable struct {
	Type       string
	Identifier string
	Visited    bool
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
	Type       string
}

type IfBlock struct {
	Condition  Expression
	Statements []Statement
	Else       interface{}
	Line       uint
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

type ParenExpression struct {
	Expression Expression
	Type       string
	Visited    bool
}

type Literal struct {
	Type    string
	Value   string
	Visited bool
}

func (*BinaryOperation) GetType() string {
	return "BinaryOperation"
}

func (*UnaryOperation) GetType() string {
	return "UnaryOperation"
}
func (*ParenExpression) GetType() string {
	return "ParenExpression"
}

func (*Literal) GetType() string {
	return "Literal"
}

func (*FuncCall) GetType() string {
	return "FuncCall"
}

func (*Variable) GetType() string {
	return "Variable"
}
