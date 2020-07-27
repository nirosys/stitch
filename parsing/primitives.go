package parsing

import (
	"fmt"
	"strconv"

	"github.com/nirosys/stitch/ast"
	"github.com/nirosys/stitch/lexing"
)

func (p *Parser) parseListExpression() ast.Expression {
	expr := &ast.ListLiteral{Token: p.curToken}
	if e := p.parseExpressionList(lexing.D_RBRACKET); e != nil {
		expr.Contents = e
	}
	return expr
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	if !p.curTokenIs(lexing.K_TRUE) && !p.curTokenIs(lexing.K_FALSE) {
		return nil
	}
	var value bool = p.curTokenIs(lexing.K_TRUE)
	return &ast.BoolLiteral{
		Token: p.curToken,
		Value: value,
	}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	if !p.curTokenIs(lexing.L_STRING) {
		return nil
	}
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Text,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	if !p.curTokenIs(lexing.L_INTEGER) {
		return nil
	}

	i, err := strconv.ParseInt(p.curToken.Text, 10, 64)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return nil
	}
	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: i,
	}
}
