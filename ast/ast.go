package ast

import (
	"strings"

	"github.com/gkampitakis/monkey/token"
)

var (
	_ Node      = (*Identifier)(nil)
	_ Node      = (*LetStatement)(nil)
	_ Node      = (*ReturnStatement)(nil)
	_ Statement = (*LetStatement)(nil)
	_ Statement = (*ReturnStatement)(nil)
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	// TODO:
	s := strings.Builder{}
	for _, stmt := range p.Statements {
		s.WriteString(stmt.TokenLiteral())
	}

	return s.String()
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (*LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return string(ls.Token.Literal)
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (*Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return string(i.Token.Literal)
}

type ReturnStatement struct {
	// 'return' Token
	Token       token.Token
	ReturnValue Expression
}

func (*ReturnStatement) statementNode() {}

func (r *ReturnStatement) TokenLiteral() string {
	return string(r.Token.Literal)
}
