package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkampitakis/monkey/token"
)

var (
	_ Node       = (*Program)(nil)
	_ Node       = (*Identifier)(nil)
	_ Node       = (*IntegerLiteral)(nil)
	_ Node       = (*PrefixExpression)(nil)
	_ Node       = (*InfixExpression)(nil)
	_ Node       = (*Boolean)(nil)
	_ Node       = (*IfExpression)(nil)
	_ Expression = (*Identifier)(nil)
	_ Expression = (*Boolean)(nil)
	_ Expression = (*InfixExpression)(nil)
	_ Expression = (*PrefixExpression)(nil)
	_ Expression = (*IntegerLiteral)(nil)
	_ Expression = (*StringLiteral)(nil)
	_ Expression = (*IfExpression)(nil)
	_ Expression = (*FunctionLiteral)(nil)
	_ Expression = (*IndexExpression)(nil)
	_ Expression = (*ArrayLiteral)(nil)
	_ Expression = (*HashLiteral)(nil)
	_ Expression = (*WhileExpression)(nil)
	_ Statement  = (*LetStatement)(nil)
	_ Statement  = (*ReturnStatement)(nil)
	_ Statement  = (*ExpressionStatement)(nil)
	_ Statement  = (*BlockStatement)(nil)
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
	s := strings.Builder{}
	for _, stmt := range p.Statements {
		s.WriteString(stmt.String())
	}

	return s.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value []byte
}

func (*Identifier) expressionNode()        {}
func (i *Identifier) TokenLiteral() string { return string(i.Token.Literal) }
func (i *Identifier) String() string {
	if i == nil {
		return ""
	}

	return string(i.Value)
}

/* Statements */

type ReturnStatement struct {
	// 'return' Token
	Token       token.Token
	ReturnValue Expression
}

func (*ReturnStatement) statementNode()         {}
func (r *ReturnStatement) TokenLiteral() string { return string(r.Token.Literal) }
func (r *ReturnStatement) String() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("%s %s;", r.TokenLiteral(), r.ReturnValue.String())
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (*LetStatement) statementNode()          {}
func (ls *LetStatement) TokenLiteral() string { return string(ls.Token.Literal) }
func (ls *LetStatement) String() string {
	if ls == nil {
		return ""
	}

	return fmt.Sprintf(
		"%s %s = %s;",
		ls.TokenLiteral(),
		ls.Name.String(),
		ls.Value.String(),
	)
}

type ExpressionStatement struct {
	// the first token of the expression
	Token      token.Token
	Expression Expression
}

func (*ExpressionStatement) statementNode()          {}
func (ex *ExpressionStatement) TokenLiteral() string { return string(ex.Token.Literal) }
func (ex *ExpressionStatement) String() string {
	if ex == nil {
		return ""
	}

	return ex.Expression.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (*IntegerLiteral) expressionNode()        {}
func (i *IntegerLiteral) TokenLiteral() string { return string(i.Token.Literal) }
func (i *IntegerLiteral) String() string {
	if i == nil {
		return ""
	}

	return strconv.Itoa(i.Value)
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (*StringLiteral) expressionNode()        {}
func (s *StringLiteral) TokenLiteral() string { return string(s.Token.Literal) }
func (s *StringLiteral) String() string {
	return s.Value
}

type PrefixExpression struct {
	// the prefix token, e.g !
	Token    token.Token
	Operator string
	Right    Expression
}

func (*PrefixExpression) expressionNode()         {}
func (px *PrefixExpression) TokenLiteral() string { return string(px.Token.Literal) }
func (px *PrefixExpression) String() string {
	if px == nil {
		return ""
	}

	return fmt.Sprintf("(%s%s)", px.Operator, px.Right.String())
}

type InfixExpression struct {
	// The operator token, e.g. +
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (*InfixExpression) expressionNode()         {}
func (ix *InfixExpression) TokenLiteral() string { return string(ix.Token.Literal) }
func (ix *InfixExpression) String() string {
	if ix == nil {
		return ""
	}

	return fmt.Sprintf("(%s %s %s)", ix.Left.String(), ix.Operator, ix.Right.String())
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (*Boolean) expressionNode()        {}
func (b *Boolean) TokenLiteral() string { return string(b.Token.Literal) }
func (b *Boolean) String() string       { return string(b.Token.Literal) }

type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

func (*BlockStatement) statementNode()         {}
func (b *BlockStatement) TokenLiteral() string { return string(b.Token.Literal) }
func (b *BlockStatement) String() string {
	str := strings.Builder{}
	for _, stmt := range b.Statements {
		str.WriteString(stmt.String())
	}

	return str.String()
}

type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (*IfExpression) expressionNode()         {}
func (ix *IfExpression) TokenLiteral() string { return string(ix.Token.Literal) }
func (ix *IfExpression) String() string {
	ifStr := fmt.Sprintf("if%s %s", ix.Condition.String(), ix.Consequence.String())
	if ix.Alternative == nil {
		return ifStr
	}

	return fmt.Sprintf("%selse%s", ifStr, ix.Alternative.String())
}

type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (*FunctionLiteral) expressionNode()         {}
func (fl *FunctionLiteral) TokenLiteral() string { return string(fl.Token.Literal) }
func (fl *FunctionLiteral) String() string {
	params := make([]string, len(fl.Parameters))

	for i, param := range fl.Parameters {
		params[i] = param.String()
	}

	return fmt.Sprintf("%s(%s)%s", fl.TokenLiteral(), strings.Join(params, ","), fl.Body.String())
}

type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression
	Arguments []Expression
}

func (*CallExpression) expressionNode()         {}
func (ce *CallExpression) TokenLiteral() string { return string(ce.Token.Literal) }
func (ce *CallExpression) String() string {
	args := make([]string, len(ce.Arguments))

	for i, arg := range ce.Arguments {
		args[i] = arg.String()
	}

	return fmt.Sprintf("%s(%s)", ce.Function.String(), strings.Join(args, ", "))
}

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (*ArrayLiteral) expressionNode()         {}
func (al *ArrayLiteral) TokenLiteral() string { return string(al.Token.Literal) }
func (al *ArrayLiteral) String() string {
	args := make([]string, len(al.Elements))

	for i, arg := range al.Elements {
		args[i] = arg.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(args, ", "))
}

type IndexExpression struct {
	Token token.Token // the [ token
	Left  Expression
	Index Expression
}

func (*IndexExpression) expressionNode()         {}
func (ie *IndexExpression) TokenLiteral() string { return string(ie.Token.Literal) }

func (ie *IndexExpression) String() string { return fmt.Sprintf("(%s[%s])", ie.Left, ie.Index) }

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (*HashLiteral) expressionNode()         {}
func (hl *HashLiteral) TokenLiteral() string { return string(hl.Token.Literal) }
func (hl *HashLiteral) String() string {
	pairs := make([]string, 0, len(hl.Pairs))

	for i, p := range hl.Pairs {
		pairs = append(pairs, fmt.Sprintf("%q: %q", i.String(), p.String()))
	}

	return fmt.Sprintf("{\n	%s\n}", strings.Join(pairs, ",\n"))
}

type WhileExpression struct {
	Token       token.Token // the 'while' token
	Condition   Expression
	Consequence *BlockStatement
}

func (*WhileExpression) expressionNode()         {}
func (we *WhileExpression) TokenLiteral() string { return string(we.Token.Literal) }
func (we *WhileExpression) String() string {
	return fmt.Sprintf("while(%s) %s", we.Condition.String(), we.Consequence.String())
}
