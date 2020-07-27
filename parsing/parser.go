package parsing

import (
	"fmt"
	"io"

	"github.com/nirosys/stitch/ast"
	"github.com/nirosys/stitch/lexing"
)

const (
	_ int = iota
	LOWEST
	OR
	AND
	EQUAL
	SUM
	PRODUCT
	PREFIX
	CALL
	DEREFERENCE
)

type infixParseFunc func(ast.Expression) ast.Expression
type prefixParseFunc func() ast.Expression
type postfixParseFunc infixParseFunc

var precedences = map[lexing.TokenType]int{
	lexing.O_COMMA:    OR,
	lexing.O_ASSIGN:   OR,
	lexing.K_OR:       OR,
	lexing.K_AND:      AND,
	lexing.O_PLUS:     SUM,
	lexing.O_MINUS:    SUM,
	lexing.O_STAR:     PRODUCT,
	lexing.O_SLASH:    PRODUCT,
	lexing.O_MODULUS:  PRODUCT,
	lexing.D_LPARENTH: CALL,
	lexing.O_DOT:      DEREFERENCE,
	lexing.O_COLON:    DEREFERENCE,
	lexing.O_EQ:       EQUAL,
	lexing.O_NEQ:      EQUAL,
	lexing.O_LT:       EQUAL,
	lexing.O_LTEQ:     EQUAL,
	lexing.O_GT:       EQUAL,
	lexing.O_GTEQ:     EQUAL,
	lexing.O_ARROW:    OR,
}

type Parser struct {
	lex *lexing.Lexer

	curToken  lexing.Token
	peekToken lexing.Token

	prefixParseFns  map[lexing.TokenType]prefixParseFunc
	infixParseFns   map[lexing.TokenType]infixParseFunc
	postFixParseFns map[lexing.TokenType]postfixParseFunc

	errors []string
}

func NewParser(r io.Reader) *Parser {
	p := &Parser{
		lex:    lexing.NewLexer(r),
		errors: []string{},
	}

	p.prefixParseFns = map[lexing.TokenType]prefixParseFunc{
		lexing.L_STRING:  p.parseStringLiteral,
		lexing.L_INTEGER: p.parseIntegerLiteral,
		lexing.IDENT:     p.parseIdentifier,
		lexing.D_LBRACE:  p.parseBlockExpression,
		// TODO: anonymous node
		//K_NODE:     p.parseNodeLiteral,
		lexing.D_LPARENTH:  p.parseGroupedExpression,
		lexing.K_FUNCTION:  p.parseFunctionExpression,
		lexing.K_INTERNAL:  p.parseInternalExpression,
		lexing.K_TRUE:      p.parseBooleanLiteral,
		lexing.K_FALSE:     p.parseBooleanLiteral,
		lexing.O_BANG:      p.parseNegatedExpression,
		lexing.K_IF:        p.parseConditionalExpression,
		lexing.D_LBRACKET:  p.parseListExpression,
		lexing.O_TAGMARKER: p.parseTagName,
	}
	p.infixParseFns = map[lexing.TokenType]infixParseFunc{
		lexing.D_LPARENTH: p.parseCallExpression,
		lexing.O_ARROW:    p.parseArrowExpression,
		lexing.O_ASSIGN:   p.parseInfixExpression,
		lexing.O_PLUS:     p.parseInfixExpression,
		lexing.O_MINUS:    p.parseInfixExpression,
		lexing.O_STAR:     p.parseInfixExpression,
		lexing.O_SLASH:    p.parseInfixExpression,
		lexing.O_MODULUS:  p.parseInfixExpression,
		lexing.O_DOT:      p.parseInfixExpression,
		lexing.O_EQ:       p.parseInfixExpression,
		lexing.O_GT:       p.parseInfixExpression,
		lexing.O_GTEQ:     p.parseInfixExpression,
		lexing.O_LT:       p.parseInfixExpression,
		lexing.O_LTEQ:     p.parseInfixExpression,
		lexing.O_NEQ:      p.parseInfixExpression,
		lexing.O_LESSTHAN: p.parseInfixExpression,
		lexing.K_AND:      p.parseInfixExpression,
		lexing.K_OR:       p.parseInfixExpression,
		lexing.O_COLON:    p.parseNamedNode,
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) CurrentToken() lexing.Token {
	return p.curToken
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Parse() *ast.ASTree {
	tree := &ast.ASTree{Statements: []ast.Statement{}}

	for p.curToken.Type != lexing.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			tree.Statements = append(tree.Statements, stmt)
		}
		p.nextToken()
	}
	if len(p.errors) > 0 {
		return nil
	}

	return tree
}

/*
func (p *Parser) ParsePartial() *ast.Program {
	prog := &ast.Program{Statements: []ast.Statement{}}

	for p.curToken.Type != lexing.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.nextToken()
	}

	return prog
}
*/

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexing.K_LET:
		return p.parseLetStatement()
	case lexing.COMMENT:
		return p.parseComment()
	case lexing.K_IMPORT:
		return p.parseImportStatement()
	case lexing.K_FUNCTION:
		return p.parseFunctionStatement()
	case lexing.K_NODE:
		return p.parseNodeStatement()
	case lexing.K_MODIFIER:
		return p.parseModifier()
	case lexing.K_FOREACH:
		return p.parseForeach()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() ast.Expression {
	exp := p.parseExpression(LOWEST)

	return exp
}

// The syntax for internals is pretty simple right now..
// we do not define in stitch what the data type is,
// only what the name is.
//
// let Agent = internal "std:Agent"
// let get_snmp = internal "snmp:get"
//
// etc..
func (p *Parser) parseInternalExpression() ast.Expression {
	exp := &ast.InternalExpression{Token: p.curToken}

	if !p.expectPeek(lexing.L_STRING) {
		return nil
	}

	if name := p.parseStringLiteral(); name != nil {
		exp.Name = name.(*ast.StringLiteral)
	}
	return exp
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(lexing.IDENT) {
		return nil // TODO: Do errors..
	}

	stmt.Name = &ast.Identifier{
		Token:      p.curToken,
		Identifier: p.curToken.Text,
	}

	if !p.expectPeek(lexing.O_ASSIGN) {
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexing.D_SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.peekTokenIs(lexing.L_STRING) {
		// TODO: ERROR
		return nil
	}
	p.nextToken()

	stmt.Path = p.curToken.Text
	return stmt
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	stmt := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(lexing.D_LPARENTH) {
		return nil
	}

	stmt.Parameters = p.parseParameterList(lexing.D_RPARENTH)

	if !p.expectPeek(lexing.O_COLON) {
		return nil
	}
	p.nextToken()

	tok := p.curToken
	expr := p.parseExpression(LOWEST)
	stmt.Body = &ast.BlockExpression{
		Token:      tok,
		Statements: []ast.Statement{expr},
	}

	return stmt
}

// Modifiers
//
// mod <identifier>(<ident>: <ident>, <dent>,...) {
// }

func (p *Parser) parseModifier() ast.Statement {
	mod := &ast.ModifierStatement{Token: p.curToken}

	if !p.expectPeek(lexing.IDENT) {
		return nil
	}
	if ident := p.parseIdentifier(); ident == nil {
		return nil
	} else {
		mod.Identifier = ident.(*ast.Identifier)
	}

	if !p.expectPeek(lexing.D_LPARENTH) {
		return nil
	}
	// Now we need to parse an ident with a type '<ident> : <ident>'
	mod.Parameters = p.parseParameterList(lexing.D_RPARENTH)

	if !p.expectPeek(lexing.D_LBRACE) {
		return nil
	}

	if block := p.parseBlockExpression(); block == nil {
		return nil
	} else {
		mod.Block = block.(*ast.BlockExpression)
	}

	return mod
}

// Node statements/definitions follow the following syntax:
//    node <name>(<input slot list>) -> (<output slot list>) <block expression>
//
//    node[Input] Test(arg) -> [Output, Error] { }
//
func (p *Parser) parseNodeStatement() ast.Statement {
	stmt := &ast.NodeStatement{Token: p.curToken}
	node := &ast.NodeLiteral{Token: p.curToken}

	if !p.expectPeek(lexing.D_LBRACKET) { // Input slots
		return nil
	}

	node.InputSlots = p.parseParameterList(lexing.D_RBRACKET)
	if node.InputSlots == nil {
		// TODO: errors?
		return nil
	}

	if !p.expectPeek(lexing.IDENT) {
		return nil
	}

	if i := p.parseIdentifier(); i == nil {
		return nil
	} else {
		stmt.Identifier = i.(*ast.Identifier)
	}

	if !p.expectPeek(lexing.D_LPARENTH) {
		return nil
	}

	args := p.parseParameterList(lexing.D_RPARENTH)
	if args == nil {
		return nil
	}
	node.Arguments = args

	if !p.expectPeek(lexing.O_ARROW) {
		return nil
	}

	if !p.expectPeek(lexing.D_LBRACKET) {
		return nil
	}

	outputs := p.parseParameterList(lexing.D_RBRACKET)
	if outputs == nil {
		return nil
	}
	node.OutputSlots = outputs
	p.nextToken()

	if b := p.parseBlockExpression(); b != nil {
		node.Block = b.(*ast.BlockExpression)
	} else {
		return nil
	}
	stmt.Literal = node

	return stmt
}

func (p *Parser) parseFunctionStatement() ast.Statement {
	// Before we create a statement.. we need to know what type of

	stmt := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(lexing.IDENT) {
		return nil
	}

	ident := p.parseIdentifier()

	if !p.expectPeek(lexing.D_LPARENTH) {
		return nil
	}

	stmt.Identifier = ident.(*ast.Identifier)
	stmt.Parameters = p.parseParameterList(lexing.D_RPARENTH)

	if !p.expectPeek(lexing.D_LBRACE) {
		return nil
	}
	body := p.parseBlockExpression()
	if body == nil {
		return nil
	}
	if b, ok := body.(*ast.BlockExpression); !ok {
		return nil
	} else {
		stmt.Body = b
	}

	return stmt
}

func (p *Parser) parseExpression(prec int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()
	for !p.peekTokenIs(lexing.D_SEMICOLON) && !p.peekTokenIs(lexing.EOF) && prec < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseNegatedExpression() ast.Expression {
	if !p.curTokenIs(lexing.O_BANG) {
		return nil
	}
	p.nextToken()
	not := &ast.NotExpression{Token: p.curToken}
	if exp := p.parseExpression(LOWEST); exp == nil {
		return nil
	} else {
		not.Expression = exp
	}
	return not
}

func (p *Parser) parseTagName() ast.Expression {
	tag := &ast.TagName{Token: p.curToken}

	if !p.curTokenIs(lexing.O_TAGMARKER) {
		// TODO: Error
		return nil
	}

	if !p.expectPeek(lexing.IDENT) {
		return nil
	}

	if ident := p.parseIdentifier(); ident == nil {
		fmt.Printf("Invalid ident\n")
		return nil
	} else {
		tag.Identifier = ident.(*ast.Identifier)
	}
	return tag
}

func (p *Parser) parseIdentifier() ast.Expression {
	if !p.curTokenIs(lexing.IDENT) {
		fmt.Printf("ERROR: %+v is not an IDENT\n", p.curToken)
		return nil
	}

	return &ast.Identifier{
		Token:      p.curToken,
		Identifier: p.curToken.Text,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // eat the '('

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexing.D_RPARENTH) {
		return nil
	}

	return exp
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: left,
	}
	exp.Arguments = p.parseExpressionList(lexing.D_RPARENTH)
	return exp
}

func (p *Parser) parseNamedNode(left ast.Expression) ast.Expression {
	curToken := p.curToken

	named := &ast.NamedNodeExpression{Token: curToken}

	p.nextToken()
	switch t := left.(type) {
	case *ast.Identifier:
		named.FieldName = t
	case *ast.TagName:
		named.TagName = t.Identifier
	}
	if exp := p.parseExpression(LOWEST); exp == nil {
		return nil
	} else {
		named.Expression = exp
	}
	return named
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	prec := p.curPrecedence()

	curToken := p.curToken

	switch curToken.Type {
	case lexing.O_ASSIGN:
		if ident, ok := left.(*ast.Identifier); ok {
			stmt := &ast.AssignmentExpression{Token: p.curToken}
			stmt.Identifier = ident
			p.nextToken()
			stmt.Value = p.parseExpression(prec)
			return stmt
		} else {
			// TODO: Error
		}
	default:
		stmt := &ast.InfixExpression{
			Token:    p.curToken,
			Left:     left,
			Operator: p.curToken.Text,
		}
		p.nextToken()
		stmt.Right = p.parseExpression(prec)
		return stmt
	}
	return nil
}

func (p *Parser) parseArrowExpression(left ast.Expression) ast.Expression {
	exp := &ast.ArrowExpression{
		Token: p.curToken,
		Left:  left,
	}
	prec := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(prec - 1) // Make us lean right so we can ultimately return the left expr

	return exp
}

// Could generalize this using a parse func as argument

func (p *Parser) parseParameterList(end lexing.TokenType) []*ast.FunctionParameter {
	list := []*ast.FunctionParameter{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	fp := &ast.FunctionParameter{Token: p.curToken}
	fp.Identifier = p.parseIdentifier().(*ast.Identifier)
	list = append(list, fp)

	for p.peekTokenIs(lexing.O_COMMA) {
		p.nextToken()
		p.nextToken()
		fp := &ast.FunctionParameter{Token: p.curToken}
		fp.Identifier = p.parseIdentifier().(*ast.Identifier)
		list = append(list, fp)
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseExpressionList(end lexing.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexing.O_COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseComment() *ast.CommentStatement {
	return &ast.CommentStatement{
		Token: p.curToken,
		Text:  p.curToken.Text,
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	//fmt.Printf("TOKEN: %+v\n", p.curToken)
	if peek, err := p.lex.NextToken(); err == nil {
		p.peekToken = peek
	} else {
		// TODO: This should bubble up
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

func (p *Parser) expectPeek(t lexing.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekTokenIs(t lexing.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t lexing.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekError(t lexing.TokenType) {
	pos := p.peekToken.Position
	exp := lexing.TokenStrings[t]
	have := lexing.TokenStrings[p.peekToken.Type]
	msg := fmt.Sprintf("line %d column %d: expected %s; have %s", pos.Line, pos.Column, exp, have)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
