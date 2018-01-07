package object

import (
	"fmt"
	"github.com/alanfoster/monkey/ast"
	"bytes"
	"strings"
)

type ObjectType int

//go:generate stringer -type=ObjectType
const (
	_            ObjectType = iota
	INTEGER
	BOOLEAN
	NULL
	RETURN_VALUE
	ERROR
	FUNCTION
	STRING
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}

func (n *Null) Inspect() string {
	return "NULL"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params []string
	for _, identifier := range f.Parameters {
		params = append(params, identifier.Value)
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.PrettyPrint())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING
}

func (s *String) Inspect() string {
	return s.Value
}