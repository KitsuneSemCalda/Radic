package ast

// Ident is a reference to a named identifier, such as a variable or
// function name.
type Ident struct {
	Name string
}

func (*Ident) expr() {}
