package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"llc/lang/ast"
)

type TypeObject string

const (
	IntegerObj     = "INTEGER"
	BooleanObj     = "BOOLEAN"
	NullObj        = "NULL"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"
	FunctionObj    = "FUNCTION"
	StringObj      = "STRING"
	BuiltinObj     = "BUILTIN"
	ArrayObj       = "ARRAY"
	HashObj        = "HASH"
	QuoteObj       = "QUOTE_OBJ"
	MacroObj       = "MACRO"
)

type HashKey struct {
	Type  TypeObject
	Value uint64
}

type Object interface {
	Type() TypeObject
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
func (h *Hash) Type() TypeObject { return HashObj }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() TypeObject { return IntegerObj }
func (i *Integer) HashKey() HashKey { return HashKey{Type: i.Type(), Value: uint64(i.Value)} } //nolint:gosec

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() TypeObject { return BooleanObj }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Type() TypeObject { return NullObj }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() TypeObject { return ReturnValueObj }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() TypeObject { return ErrorObj }
func (e *Error) Inspect() string  { return "Error: " + e.Message }

type Function struct {
	Body       *ast.BlockStatement
	Env        *Environment
	Parameters []*ast.Identifier
}

func (f *Function) Type() TypeObject { return FunctionObj }
func (f *Function) Inspect() string {
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}

	return fmt.Sprintf("fn(%s) {\n%s\n}", strings.Join(params, ", "), f.Body.String())
}

type String struct {
	Value string
}

func (s *String) Type() TypeObject { return StringObj }
func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	hash := fnv.New64a()
	hash.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: hash.Sum64()}
}

type BuiltinFunction = func(args ...Object) Object

type Builtin struct {
	Function BuiltinFunction
}

func (s *Builtin) Type() TypeObject { return BuiltinObj }
func (s *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (s *Array) Type() TypeObject { return ArrayObj }
func (s *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range s.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() TypeObject { return QuoteObj }
func (q *Quote) Inspect() string {
	return fmt.Sprintf("QUOTE(%s)", q.Node.String())
}

type Macro struct {
	Body       *ast.BlockStatement
	Env        *Environment
	Parameters []*ast.Identifier
}

func (m *Macro) Type() TypeObject { return MacroObj }
func (m *Macro) Inspect() string {
	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf(
		"macro(%s) {\n%s\n}",
		strings.Join(params, ", "),
		m.Body.String(),
	)
}
