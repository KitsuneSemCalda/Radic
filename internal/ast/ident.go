package ast

type Ident struct {
	Name string
}

func (*Ident) expr() {}
