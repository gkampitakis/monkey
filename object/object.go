package object

import (
	"fmt"
	"strings"

	"github.com/gkampitakis/monkey/ast"
)

var (
	_ Object = (*Integer)(nil)
	_ Object = (*Boolean)(nil)
	_ Object = (*Null)(nil)
	_ Object = (*ReturnValue)(nil)
	_ Object = (*ErrorValue)(nil)
	_ Object = (*Function)(nil)
)

//go:generate stringer -type=ObjectType
type ObjectType uint8

const (
	INTEGER ObjectType = iota
	BOOLEAN
	NULL
	RETURN_VALUE
	ERROR_VALUE
	FUNCTION
	STRING
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int
}

func (*Integer) Type() ObjectType  { return INTEGER }
func (i *Integer) Inspect() string { return fmt.Sprint(i.Value) }

type Boolean struct {
	Value bool
}

func (*Boolean) Type() ObjectType  { return BOOLEAN }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (*Null) Type() ObjectType  { return NULL }
func (n *Null) Inspect() string { return "null" }

type ReturnValue struct {
	Value Object
}

func (*ReturnValue) Type() ObjectType  { return RETURN_VALUE }
func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }

type ErrorValue struct {
	Message string
}

func (*ErrorValue) Type() ObjectType  { return ERROR_VALUE }
func (r *ErrorValue) Inspect() string { return "[error]: " + r.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (*Function) Type() ObjectType { return FUNCTION }
func (f *Function) Inspect() string {
	params := make([]string, len(f.Parameters))

	for i, param := range f.Parameters {
		params[i] = param.String()
	}

	return fmt.Sprintf("fn(%s){\n%s\n}", strings.Join(params, ","), f.Body.String())
}

type String struct {
	Value string
}

func (*String) Type() ObjectType  { return STRING }
func (s *String) Inspect() string { return s.Value }
