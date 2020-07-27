package ast

import (
	"bytes"
	"strings"

	"github.com/nirosys/stitch/lexing"
)

type ASTree struct {
	Statements []Statement `json:"statements"`
}

type Node interface {
	TokenLiteral() string
	String() string
}

type Expression interface {
	Statement

	expressionNode()
}

type ExpressionStatement struct {
	Token      lexing.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Text }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type CallExpression struct {
	Token     lexing.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) statementNode()       {}
func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Text }
func (c *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(c.Function.String())
	out.WriteByte('(')
	if len(c.Arguments) > 0 {
		args := []string{}
		for _, arg := range c.Arguments {
			args = append(args, arg.String())
		}
		out.WriteString(strings.Join(args, ", "))
	}
	out.WriteByte(')')
	return out.String()
}
