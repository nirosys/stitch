package ast

import (
	"bytes"

	"github.com/nirosys/stitch/lexing"
)

type ArrowExpression struct {
	Expression

	Token lexing.Token
	Left  Expression
	Right Expression
}

func (a *ArrowExpression) statementNode()       {}
func (a *ArrowExpression) TokenLiteral() string { return a.Token.Text }
func (a *ArrowExpression) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(a.Left.String())
	buffer.WriteString(" -> ")
	buffer.WriteString(a.Right.String())
	return buffer.String()
}

type BlockExpression struct {
	Expression

	Token      lexing.Token
	Statements []Statement
}

func (b *BlockExpression) statementNode()       {}
func (b *BlockExpression) TokenLiteral() string { return b.Token.Text }
func (b *BlockExpression) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{\n")
	for _, stmt := range b.Statements {
		buffer.WriteString("  ")
		buffer.WriteString(stmt.String())
		buffer.WriteString(";\n")
	}
	buffer.WriteString("}\n")
	return buffer.String()
}

// Assignment Expression //////////////////////////////////////////////////////

type AssignmentExpression struct {
	Token      lexing.Token
	Identifier *Identifier
	Value      Expression
}

func (a *AssignmentExpression) statementNode()       {}
func (a *AssignmentExpression) expressionNode()      {}
func (a *AssignmentExpression) TokenLiteral() string { return a.Token.Text }
func (a *AssignmentExpression) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(a.Identifier.String())
	buffer.WriteString(" = ")
	buffer.WriteString(a.Value.String())
	return buffer.String()
}

// Infix Expression ///////////////////////////////////////////////////////////
type InfixExpression struct {
	Token    lexing.Token
	Left     Expression
	Right    Expression
	Operator string
}

func (i *InfixExpression) statementNode()       {}
func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Text }
func (i *InfixExpression) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(i.Left.String())
	buffer.WriteString(i.Operator)
	buffer.WriteString(i.Right.String())
	return buffer.String()
}

// FieldName //////////////////////////////////////////////////////////////////
type Tag struct {
	Token      lexing.Token
	Expression Expression
}

func (t *Tag) statementNode()       {}
func (t *Tag) expressionNode()      {}
func (t *Tag) TokenLiteral() string { return t.Token.Text }
func (t *Tag) String() string       { return t.Expression.String() }

// InternalExpression /////////////////////////////////////////////////////////

type InternalExpression struct {
	Token lexing.Token

	Name *StringLiteral
}

func (i *InternalExpression) statementNode()       {}
func (i *InternalExpression) expressionNode()      {}
func (i *InternalExpression) TokenLiteral() string { return i.Token.Text }
func (i *InternalExpression) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("internal \"")
	buffer.WriteString(i.Name.String())
	buffer.WriteByte('"')
	return buffer.String()
}

// ConditionalExpression /////////////////////////////////////////////////////
type ConditionalExpression struct {
	Token lexing.Token

	Condition Expression
	Block     *BlockExpression
	Else      Expression
}

func (i *ConditionalExpression) statementNode()       {}
func (i *ConditionalExpression) expressionNode()      {}
func (i *ConditionalExpression) TokenLiteral() string { return i.Token.Text }
func (i *ConditionalExpression) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("if ")
	buffer.WriteString(i.Condition.String())
	buffer.WriteString(i.Block.String())
	if i.Else != nil {
		buffer.WriteString(" else ")
		buffer.WriteString(i.Else.String())
	}
	return buffer.String()
}

// NamedNodeExpression ////////////////////////////////////////////////////////
type NamedNodeExpression struct {
	Token lexing.Token

	FieldName  *Identifier
	TagName    *Identifier
	Expression Expression
}

func (n *NamedNodeExpression) statementNode()       {}
func (n *NamedNodeExpression) expressionNode()      {}
func (n *NamedNodeExpression) TokenLiteral() string { return n.Token.Text }
func (n *NamedNodeExpression) String() string {
	return "<not implemented>"
}

// NotExpression //////////////////////////////////////////////////////////////
type NotExpression struct {
	Token lexing.Token

	Expression Expression
}

func (n *NotExpression) statementNode()       {}
func (n *NotExpression) expressionNode()      {}
func (n *NotExpression) TokenLiteral() string { return n.Token.Text }
func (n *NotExpression) String() string {
	return "!" + n.Expression.String()
}
