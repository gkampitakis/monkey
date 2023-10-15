package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/gkampitakis/monkey/ast"
)

var (
	_ Object   = (*Integer)(nil)
	_ Object   = (*Boolean)(nil)
	_ Object   = (*Null)(nil)
	_ Object   = (*ReturnValue)(nil)
	_ Object   = (*ErrorValue)(nil)
	_ Object   = (*Function)(nil)
	_ Object   = (*String)(nil)
	_ Object   = (*Array)(nil)
	_ Object   = (*Builtin)(nil)
	_ Object   = (*Hash)(nil)
	_ Hashable = (*Boolean)(nil)
	_ Hashable = (*String)(nil)
	_ Hashable = (*Integer)(nil)
)

type BuiltinFunction func(args ...Object) Object

type Hashable interface {
	HashKey() HashKey
}

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
	BUILTIN
	ARRAY
	HASH
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int
}

func (*Integer) Type() ObjectType   { return INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprint(i.Value) }
func (i *Integer) HashKey() HashKey { return HashKey{Type: i.Type(), Value: i.Value} }

type Boolean struct {
	Value bool
}

func (*Boolean) Type() ObjectType  { return BOOLEAN }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	if b.Value {
		return HashKey{Type: b.Type(), Value: 1}
	}

	return HashKey{Type: b.Type(), Value: 0}
}

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
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: int(h.Sum64())}
}

type Builtin struct {
	Fn BuiltinFunction
}

func (*Builtin) Type() ObjectType { return BUILTIN }
func (*Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (*Array) Type() ObjectType { return ARRAY }
func (a *Array) Inspect() string {
	elements := make([]string, len(a.Elements))

	for i, e := range a.Elements {
		elements[i] = e.Inspect()
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ","))
}

type HashKey struct {
	Type  ObjectType
	Value int
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH }
func (h *Hash) Inspect() string {
	pairs := make([]string, 0, len(h.Pairs))

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%q: %q", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(pairs, ",\n  "))
}
