package parser

import (
	"fmt"
	"strconv"
	"zlang/ast"
	"zlang/lexer"
	"zlang/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      //==
	LESSGREATER // > OR <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X OR !X
	CALL        //myFunction(x)
)

type Parser struct {
	curToken, peekToken token.Token
	l                   *lexer.Lexer
	errors              []string
	prefixParseFns      map[token.TokenType]prefixParseFns
	infixParseFns       map[token.TokenType]infixParseFns
}

type (
	prefixParseFns func() ast.Expression
	infixParseFns  func(expression ast.Expression) ast.Expression
)

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()
	p.prefixParseFns = make(map[token.TokenType]prefixParseFns)
	p.infixParseFns = make(map[token.TokenType]infixParseFns)

	p.registerPrefixParseFns(token.Identifier, p.parseIdentifier)
	p.registerPrefixParseFns(token.Integer, p.parseIntegerLiteral)
	p.registerPrefixParseFns(token.BANG, p.parsePrefixExpression)
	p.registerPrefixParseFns(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixParseFns(token.True, p.parseBoolLiteral)
	p.registerPrefixParseFns(token.False, p.parseBoolLiteral)
	p.registerPrefixParseFns(token.LPAREN, p.parseGroupedExpress)
	p.registerPrefixParseFns(token.IF, p.parseIfExpress)

	p.registerInfixParseFns(token.PLUS, p.parseInfixExpression)
	p.registerInfixParseFns(token.MINUS, p.parseInfixExpression)
	p.registerInfixParseFns(token.SLASH, p.parseInfixExpression)
	p.registerInfixParseFns(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixParseFns(token.EQ, p.parseInfixExpression)
	p.registerInfixParseFns(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixParseFns(token.LT, p.parseInfixExpression)
	p.registerInfixParseFns(token.GT, p.parseInfixExpression)
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	var program = &ast.Program{}
	for p.curToken.Type != token.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) nextStatement() ast.Statement {
	return nil
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if p.peekToken.Type != token.Identifier {
		p.peekError(token.Identifier, p.curToken.Type)
		return nil
	}

	p.nextToken()
	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.peekTokenIs(token.ASSIGN) {
		p.peekError(token.ASSIGN, p.peekToken.Type)
		return nil
	}

	p.nextToken()
	//p.nextToken()

	if p.peekTokenIs(token.IF) {
		p.nextToken()
		stmt.Value = p.parseIfExpress()
		return stmt
	}

	//todo express handle until meet semi
	for p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.EOF {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectToken(tokenType token.TokenType) bool {
	if p.peekToken.Type == tokenType {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}
func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) peekError(should, cur token.TokenType) {
	msg := fmt.Sprintf("expected next token type to be %s, got: %s", should, cur)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) registerPrefixParseFns(tokenType token.TokenType, fn prefixParseFns) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixParseFns(tokenType token.TokenType, fn infixParseFns) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseGroupedStatement() *ast.BlockStatement {
	group := &ast.BlockStatement{Token: p.curToken}
	p.nextToken()
	if !p.curTokenIs(token.LBRACE) {
		return nil
	}
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.ParseStatement()
		if stmt != nil {
			group.Statements = append(group.Statements, stmt)
		}

		p.nextToken()
	}

	return group
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	p.parseGroupedStatement()

	if !p.peekTokenIs(token.Else) {
		return stmt
	}

	stmt.ElseStatement = p.parseGroupedStatement()

	return stmt
}

func (p *Parser) parseIfExpress() ast.Expression {
	stmt := &ast.IfExpress{Token: p.curToken}
	//p.nextToken()
	if !p.peekTokenIs(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if !p.peekTokenIs(token.RPAREN) {
		return nil
	}

	p.nextToken()
	if !p.peekTokenIs(token.LBRACE) {
		return nil
	}

	p.nextToken()
	stmt.TrueStatement = p.parseExpression(LOWEST)

	if !p.peekTokenIs(token.RBRACE) {
		return nil
	}

	if !p.peekTokenIs(token.RBRACE) {
		return nil
	}

	p.nextToken()
	if !p.peekTokenIs(token.Else) {
		return nil
	}

	p.nextToken()
	if !p.peekTokenIs(token.Else) {
		return nil
	}
	p.nextToken()
	if !p.peekTokenIs(token.LBRACE) {
		return nil
	}

	p.nextToken()
	stmt.ElseStatement = p.parseExpression(LOWEST)
	if !p.peekTokenIs(token.RBRACE) {
		return nil
	}

	p.nextToken()
	return stmt
}

func (p *Parser) parseIfExpressBlockStatement() {

}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(lit.Token.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %s as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %q", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseBoolLiteral() ast.Expression {
	return &ast.Bool{
		Token: p.curToken,
		Value: p.curTokenIs(token.True),
	}
}

func (p *Parser) parseGroupedExpress() ast.Expression {
	ge := &ast.GroupExpress{Token: p.curToken}
	p.nextToken()

	ge.Express = p.parseExpression(LOWEST)
	if !p.peekTokenIs(token.RPAREN) {
		return nil
	}

	p.nextToken()

	return ge
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
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
