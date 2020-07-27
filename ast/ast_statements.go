package ast

import (
	"bytes"

	"github.com/nirosys/stitch/lexing"
)

type Statement interface {
	Node

	statementNode()
}

type LetStatement struct {
	Token lexing.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Text }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	return out.String()
}

type CommentStatement struct {
	Token lexing.Token
	Text  string
}

func (c *CommentStatement) statementNode()       {}
func (c *CommentStatement) TokenLiteral() string { return c.Token.Text }
func (c *CommentStatement) String() string {
	var out bytes.Buffer
	out.WriteString("# ")
	out.WriteString(c.Text)
	return out.String()
}

type ImportStatement struct {
	Token lexing.Token

	Path string
}

func (i *ImportStatement) statementNode()       {}
func (i *ImportStatement) TokenLiteral() string { return i.Token.Text }
func (i *ImportStatement) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("import \"")
	buffer.WriteString(i.Path)
	buffer.WriteString("\"\n")
	return buffer.String()
}

type NodeStatement struct {
	Token lexing.Token

	Identifier *Identifier
	Literal    *NodeLiteral
}

func (n *NodeStatement) statementNode()       {}
func (n *NodeStatement) TokenLiteral() string { return n.Token.Text }
func (n *NodeStatement) String() string {
	return "<not implemented>"
}

// ModifierStatement //////////////////////////////////////////////////////////
type ModifierStatement struct {
	Token lexing.Token

	Identifier *Identifier
	Parameters []*FunctionParameter
	Block      *BlockExpression
}

func (m *ModifierStatement) statementNode()       {}
func (m *ModifierStatement) TokenLiteral() string { return m.Token.Text }
func (m *ModifierStatement) String() string {
	return "<not implemented>"
}

// ForeachStatement ///////////////////////////////////////////////////////////
type ForeachStatement struct {
	Token lexing.Token

	LoopVar *Identifier
	List    Expression
	Block   *BlockExpression
}

func (f *ForeachStatement) statementNode()       {}
func (f *ForeachStatement) TokenLiteral() string { return f.Token.Text }
func (f *ForeachStatement) String() string {
	return "<not implemented>"
}
