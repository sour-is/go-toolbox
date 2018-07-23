package rsql

import "bytes"

type Node interface {
	TokenLiteral() string
	String() string
}

type Expression interface {
	Node
	expressionNode()
}

type IntegerLiteral struct{
	Token Token
	Value int64
}
func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string { return i.Token.Literal }

type FloatLiteral struct{
	Token Token
	Value float64
}
func (i *FloatLiteral) expressionNode() {}
func (i *FloatLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *FloatLiteral) String() string { return i.Token.Literal }

type BoolLiteral struct{
	Token Token
	Value bool
}
func (i *BoolLiteral) expressionNode() {}
func (i *BoolLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *BoolLiteral) String() string { return i.Token.Literal }

type NullLiteral struct{
	Token Token
	Value bool
}
func (i *NullLiteral) expressionNode() {}
func (i *NullLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *NullLiteral) String() string { return i.Token.Literal }

type StringLiteral struct{
	Token Token
	Value string
}
func (i *StringLiteral) expressionNode() {}
func (i *StringLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *StringLiteral) String() string {
	var out bytes.Buffer

	out.WriteRune('"')
	out.WriteString(i.Token.Literal)
	out.WriteRune('"')

	return out.String()
}

type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}
