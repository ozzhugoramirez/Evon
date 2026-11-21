package compiler

import "fmt"

type ObjectType string

const (
	STRING_OBJ  = "STRING"
	INTEGER_OBJ = "INTEGER"
	NULL_OBJ    = "NULL"
	BUILTIN_OBJ = "BUILTIN" // Nuevo objeto de la librería estándar
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// --- LIBRERÍA ESTÁNDAR ---
type Builtin struct {
	Name string
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "Función Nativa <" + b.Name + ">" }
