package ast

type Member struct {
	Object Expr
	Name   string
	Deref  bool
}

func (*Member) expr() {}
