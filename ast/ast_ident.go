package ast

import (
	"bytes"

	"github.com/nirosys/stitch/lexing"
)

// Identifier /////////////////////////////////////////////////////////////////
type Identifier struct {
	Expression

	Token      lexing.Token
	Identifier string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Text }
func (i *Identifier) String() string       { return i.Identifier }

// FunctionParameter //////////////////////////////////////////////////////////
type FunctionParameter struct {
	Token      lexing.Token
	Identifier *Identifier
	Type       *Identifier
}

func (p *FunctionParameter) statementNode()       {}
func (p *FunctionParameter) TokenLiteral() string { return p.Token.Text }
func (p *FunctionParameter) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(p.Identifier.String())
	if p.Type != nil {
		buffer.WriteByte(':')
		buffer.WriteString(p.Type.String())
	}
	return buffer.String()
}

// TagName ////////////////////////////////////////////////////////////////////
type TagName struct {
	Token      lexing.Token
	Identifier *Identifier
}

func (t *TagName) statementNode()       {}
func (t *TagName) expressionNode()      {}
func (t *TagName) TokenLiteral() string { return t.Token.Text }
func (t *TagName) String() string {
	return "@" + t.Identifier.String()
}
