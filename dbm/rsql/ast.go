package rsql

import (
	"bytes"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}
type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}

type Identifier struct{
	Token Token
	Value string
}
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

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
	out.WriteString(i.Value)
	out.WriteRune('"')

	return out.String()
}

type Array struct{
	Token Token
	Elements []Expression
}
func (a *Array) expressionNode() {}
func (a *Array) TokenLiteral() string { return a.Token.Literal }
func (a *Array) String() string {
	var out bytes.Buffer

	var elements []string
	for _, el := range a.Elements {
		elements = append(elements, el.String())
	}

	out.WriteRune('(')
	out.WriteString(strings.Join(elements, ","))
	out.WriteRune(')')

	return out.String()
}


type Program struct {
	Statements []Statement
}
func (p *Program) expressionNode() {}
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type ExpressionStatement struct {
	Token Token
	Expression Expression
}
func (ExpressionStatement) statementNode()         {}
func (e ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e ExpressionStatement) String()       string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}
func (p *PrefixExpression) expressionNode()        {}
func (p *PrefixExpression) TokenLiteral()   string { return p.Token.Literal }
func (p *PrefixExpression) String()         string {
	var out bytes.Buffer

	out.WriteRune('(')
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteRune(')')

	return out.String()
}

type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}
func (i *InfixExpression) expressionNode()        {}
func (i *InfixExpression) TokenLiteral()   string { return i.Token.Literal }
func (i *InfixExpression) String()         string {
	var out bytes.Buffer

	out.WriteRune('(')
	if i.Left != nil {
		out.WriteString(i.Left.String())
	} else {
		out.WriteString("nil")
	}
	out.WriteString(i.Operator)
	if i.Right != nil {
		out.WriteString(i.Right.String())
	} else {
		out.WriteString("nil")
	}
	out.WriteRune(')')

	return out.String()
}

