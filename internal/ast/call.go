package ast

type Call struct {
	Fn   Expr
	Args []Expr
}

func (*Call) expr() {}
