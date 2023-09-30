package object

import "fmt"

var (
	_ Object = (*Integer)(nil)
	_ Object = (*Boolean)(nil)
	_ Object = (*Null)(nil)
	_ Object = (*ReturnValue)(nil)
	_ Object = (*ErrorValue)(nil)
)

//go:generate stringer -type=ObjectType
type ObjectType uint8

const (
	INTEGER ObjectType = iota
	BOOLEAN
	NULL
	RETURN_VALUE
	ERROR_VALUE
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int
}

func (i *Integer) Type() ObjectType { return INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprint(i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }

type ErrorValue struct {
	Message string
}

func (r *ErrorValue) Type() ObjectType { return ERROR_VALUE }
func (r *ErrorValue) Inspect() string  { return "[error]: " + r.Message }
