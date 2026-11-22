// Package object defines values produced by the interpreter.
package object

import "strconv"

type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct{ Value int64 }

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }

type Boolean struct{ Value bool }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue carries a return through nested blocks to the program boundary.
type ReturnValue struct{ Value Object }

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }

type Error struct{ Message string }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
