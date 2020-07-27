package ast

import (
	"bytes"
	"strings"

	"github.com/nirosys/stitch/lexing"
)

type Literal interface {
	Expression

	Value() interface{}
}

type StringLiteral struct {
	Literal
	Expression

	Token lexing.Token
	Value string
}

func (s *StringLiteral) statementNode()       {}
func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Text }
func (s *StringLiteral) String() string {
	var buffer bytes.Buffer
	buffer.WriteByte('"')
	buffer.WriteString(s.Token.Text)
	buffer.WriteByte('"')
	return buffer.String()
}

type IntegerLiteral struct {
	Token lexing.Token
	Value int64
}

func (i *IntegerLiteral) statementNode()       {}
func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Text }
func (i *IntegerLiteral) String() string       { return i.Token.Text }

/// Node Literal //////////////////////////////////////////////////////////////

type NodeLiteral struct {
	Token lexing.Token

	Identifier  *Identifier // nil unless we are a definition.
	Block       *BlockExpression
	Arguments   []*FunctionParameter
	InputSlots  []*FunctionParameter
	OutputSlots []*FunctionParameter
}

func (n *NodeLiteral) statementNode()       {}
func (n *NodeLiteral) expressionNode()      {}
func (n *NodeLiteral) TokenLiteral() string { return n.Token.Text }
func (n *NodeLiteral) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("node ")
	if n.Identifier != nil {
		buffer.WriteString(n.Identifier.String())
	}
	buffer.WriteByte('(')
	buffer.WriteByte(')')
	buffer.WriteString(" -> ()")
	buffer.WriteString(n.Block.String())
	return buffer.String()
}

/// Function Literal //////////////////////////////////////////////////////////

type FunctionLiteral struct {
	Token lexing.Token

	Identifier *Identifier
	Parameters []*FunctionParameter
	Body       *BlockExpression
}

func (f *FunctionLiteral) statementNode()       {}
func (f *FunctionLiteral) expressionNode()      {}
func (f *FunctionLiteral) TokenLiteral() string { return f.Token.Text }
func (f *FunctionLiteral) String() string {
	return "<not implemented>"
}

/// List Literal /////////////////////////////////////////////////////////////

type ListLiteral struct {
	Token lexing.Token

	Contents []Expression
}

func (l *ListLiteral) statementNode()       {}
func (l *ListLiteral) expressionNode()      {}
func (l *ListLiteral) TokenLiteral() string { return l.Token.Text }
func (l *ListLiteral) String() string {
	var buffer bytes.Buffer
	buffer.WriteByte('[')
	vals := []string{}
	for _, v := range l.Contents {
		vals = append(vals, v.String())
	}
	buffer.WriteString(strings.Join(vals, ","))
	buffer.WriteByte(']')
	return buffer.String()
}

/// Boolean Literal ///////////////////////////////////////////////////////////

type BoolLiteral struct {
	Token lexing.Token

	Value bool
}

func (b *BoolLiteral) statementNode()       {}
func (b *BoolLiteral) expressionNode()      {}
func (b *BoolLiteral) TokenLiteral() string { return b.Token.Text }
func (b *BoolLiteral) String() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}

/// Map Literal ///////////////////////////////////////////////////////////////
type MapLiteral struct {
	Token       lexing.Token
	Assignments []*AssignmentExpression
}

func (m *MapLiteral) statementNode()       {}
func (m *MapLiteral) expressionNode()      {}
func (m *MapLiteral) TokenLiteral() string { return m.Token.Text }
func (m *MapLiteral) String() string {
	return "<not implemented"
}
