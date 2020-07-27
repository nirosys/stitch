package parsing

import (
	"fmt"

	"github.com/nirosys/stitch/ast"
	"github.com/nirosys/stitch/lexing"
)

func (p *Parser) parseForeach() ast.Statement {
	stmt := &ast.ForeachStatement{Token: p.curToken}

	if !p.expectPeek(lexing.IDENT) {
		return nil
	}

	if ident := p.parseIdentifier(); ident == nil {
		return nil
	} else {
		stmt.LoopVar = ident.(*ast.Identifier)
	}

	if !p.expectPeek(lexing.K_IN) {
		return nil
	}
	p.nextToken()

	if exp := p.parseExpression(LOWEST); exp == nil {
		return nil
	} else {
		stmt.List = exp
	}

	if !p.expectPeek(lexing.D_LBRACE) {
		return nil
	}

	if block := p.parseBlockExpression(); block == nil {
		return nil
	} else {
		stmt.Block = block.(*ast.BlockExpression)
	}

	return stmt
}

func (p *Parser) parseConditionalExpression() ast.Expression {
	exp := &ast.ConditionalExpression{Token: p.curToken}

	p.nextToken()

	exp.Condition = p.parseExpression(LOWEST)

	if exp.Condition == nil {
		return nil
	}

	if !p.expectPeek(lexing.D_LBRACE) {
		return nil
	}

	if e := p.parseBlockExpression(); e != nil {
		exp.Block = e.(*ast.BlockExpression)
	} else {
		fmt.Printf("Returning block not expression: %#v", e)
		return nil
	}

	if p.peekTokenIs(lexing.K_ELSE) {
		p.nextToken()
		p.nextToken()
		exp.Else = p.parseExpression(LOWEST)
	}

	return exp
}

func (p *Parser) parseBlockExpression() ast.Expression {
	block := &ast.BlockExpression{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	allAssign := true // A block that is all assignments should be interpreted as a map.

	for p.curToken.Type != lexing.D_RBRACE && p.curToken.Type != lexing.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
			_, assign := stmt.(*ast.AssignmentExpression)
			allAssign = allAssign && assign
		}
		p.nextToken()
	}

	if allAssign {
		maplit := &ast.MapLiteral{Token: block.Token}
		maplit.Assignments = make([]*ast.AssignmentExpression, 0, len(block.Statements))
		for _, stmt := range block.Statements {
			assign := stmt.(*ast.AssignmentExpression)
			maplit.Assignments = append(maplit.Assignments, assign)
		}
		return maplit
	}

	return block
}
