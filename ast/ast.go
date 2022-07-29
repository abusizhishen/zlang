package ast

import (
	"bytes"
	"zlang/token"
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

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString("=")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

type Bool struct {
	Token token.Token
	Value bool
}

func (b *Bool) expressionNode()      {}
func (b *Bool) TokenLiteral() string { return b.Token.Literal }
func (b *Bool) String() string       { return b.Token.Literal }

type GroupExpress struct {
	Token   token.Token
	Express Expression
}

func (g *GroupExpress) expressionNode()      {}
func (g *GroupExpress) TokenLiteral() string { return g.Token.Literal }
func (g *GroupExpress) String() string       { return g.Express.String() }

func (id *Identifier) expressionNode()      {}
func (id *Identifier) TokenLiteral() string { return id.Token.Literal }
func (id *Identifier) String() string       { return id.Value }

type IfStatement struct {
	Token         token.Token
	Condition     Expression
	TrueStatement Statement
	//ElseIfStatement s
	ElseStatement Statement
}

func (i *IfStatement) statementNode()       {}
func (i *IfStatement) TokenLiteral() string { return i.Token.Literal }
func (i *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString(" if ")
	out.WriteString(i.Condition.String())
	out.WriteString(i.TrueStatement.String())
	if i.ElseStatement != nil {
		out.WriteString(i.ElseStatement.String())
	}

	return out.String()
}

type IfExpress struct {
	Token         token.Token
	Condition     Expression
	TrueStatement Expression
	ElseStatement Expression
}

func (i *IfExpress) expressionNode()      {}
func (i *IfExpress) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpress) String() string {
	var out bytes.Buffer
	out.WriteString(" if ")

	out.WriteString(" " + token.LPAREN + " ")
	out.WriteString(i.Condition.String())
	out.WriteString(" " + token.RPAREN + " ")

	out.WriteString(" " + token.LBRACE + " ")
	out.WriteString(i.TrueStatement.String())
	out.WriteString(" " + token.RBRACE + " ")

	out.WriteString(" " + token.Else + " ")
	out.WriteString(" " + token.LBRACE + " ")
	out.WriteString(i.ElseStatement.String())
	out.WriteString(" " + token.RBRACE + " ")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (g *BlockStatement) statementNode()       {}
func (g *BlockStatement) TokenLiteral() string { return g.Token.Literal }
func (g *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString(" { \n")
	for _, stmt := range g.Statements {
		out.WriteString(stmt.String())
	}

	out.WriteString(" } ")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}
