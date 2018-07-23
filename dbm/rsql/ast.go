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

type Integer struct{
	Token Token
	Value int64
}
func (i *Integer) expressionNode()      {}
func (i *Integer) TokenLiteral() string { return i.Token.Literal }
func (i *Integer) String() string       { return i.Token.Literal }

type Float struct{
	Token Token
	Value float64
}
func (i *Float) expressionNode()      {}
func (i *Float) TokenLiteral() string { return i.Token.Literal }
func (i *Float) String() string       { return i.Token.Literal }

type Bool struct{
	Token Token
	Value bool
}
func (i *Bool) expressionNode()      {}
func (i *Bool) TokenLiteral() string { return i.Token.Literal }
func (i *Bool) String() string       { return i.Token.Literal }

type Null struct{
	Token Token
	Value bool
}
func (i *Null) expressionNode()      {}
func (i *Null) TokenLiteral() string { return i.Token.Literal }
func (i *Null) String() string       { return i.Token.Literal }

type String struct{
	Token Token
	Value string
}
func (i *String) expressionNode()      {}
func (i *String) TokenLiteral() string { return i.Token.Literal }
func (i *String) String() string {
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
func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteRune('(')
	out.WriteString(i.Left.String())
	out.WriteString(i.Operator)
	out.WriteString(i.Right.String())
	out.WriteRune(')')

	return out.String()
}
